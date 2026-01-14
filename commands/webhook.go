package commands

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/webhook"
)

var (
	webhookAddr          string
	webhookSecret        string
	webhookJWKsURL       string
	webhookCanvasDataSvc bool
	webhookEvents        []string
	webhookLog           bool
)

var webhookCmd = &cobra.Command{
	Use:   "webhook",
	Short: "Manage Canvas webhook listeners",
	Long: `Start and manage webhook listeners for Canvas LMS events.

The webhook listener receives real-time notifications from Canvas for
various events like assignment creation, submission updates, grade changes,
enrollment changes, and more.

The listener provides:
  - JWT verification for Canvas Data Services (recommended)
  - HMAC signature verification for custom integrations
  - Event routing to specific handlers
  - Graceful shutdown
  - Health check endpoint
  - Request logging

Examples:
  # Start with Canvas Data Services JWT verification (recommended)
  canvas webhook listen --canvas-data-services

  # Start with HMAC signature verification
  canvas webhook listen --secret your-secret-key

  # Start with custom JWK URL
  canvas webhook listen --jwks-url https://your-jwks-endpoint.com/jwks

  # Listen for specific event types
  canvas webhook listen --events submission_created,grade_change

  # List available event types
  canvas webhook events`,
}

var webhookListenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Start webhook listener server",
	Long: `Start a webhook listener server to receive Canvas events.

The server listens on the specified address and processes incoming
webhook events. It supports two verification methods:

1. JWT verification (Canvas Data Services):
   Use --canvas-data-services or --jwks-url to verify JWT-signed payloads
   from Canvas Data Services. JWKs are fetched from Canvas and cached.

2. HMAC verification (custom integrations):
   Use --secret to verify HMAC-SHA256 signatures for custom webhooks.

Endpoints:
  POST /webhook - Receive webhook events
  GET  /health  - Health check endpoint

Examples:
  # Canvas Data Services mode (recommended for production)
  canvas webhook listen --canvas-data-services --log

  # Custom JWK endpoint
  canvas webhook listen --jwks-url https://your-jwks.com/jwks

  # HMAC signature verification
  canvas webhook listen --secret my-secret --log

  # Both modes (fallback)
  canvas webhook listen --canvas-data-services --secret backup-secret

  # Listen for specific events
  canvas webhook listen --events grade_change,submission_created`,
	RunE: runWebhookListen,
}

var webhookEventsCmd = &cobra.Command{
	Use:   "events",
	Short: "List available Canvas webhook event types",
	Long:  `Display all supported Canvas webhook event types with their descriptions.`,
	Run:   runWebhookEvents,
}

func init() {
	rootCmd.AddCommand(webhookCmd)
	webhookCmd.AddCommand(webhookListenCmd)
	webhookCmd.AddCommand(webhookEventsCmd)

	// Listen command flags
	webhookListenCmd.Flags().StringVar(&webhookAddr, "addr", ":8080", "Server address to listen on")
	webhookListenCmd.Flags().StringVar(&webhookSecret, "secret", "", "Webhook secret for HMAC verification")
	webhookListenCmd.Flags().StringVar(&webhookJWKsURL, "jwks-url", "", "JWK Set URL for JWT verification")
	webhookListenCmd.Flags().BoolVar(&webhookCanvasDataSvc, "canvas-data-services", false, "Use Canvas Data Services JWK URL for JWT verification")
	webhookListenCmd.Flags().StringSliceVar(&webhookEvents, "events", []string{}, "Event types to handle (comma-separated)")
	webhookListenCmd.Flags().BoolVar(&webhookLog, "log", false, "Enable request logging")
}

func runWebhookListen(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Create logger
	var logger *log.Logger
	if webhookLog {
		logger = log.New(os.Stdout, "[webhook] ", log.LstdFlags)
	} else {
		logger = log.New(os.Stderr, "", log.LstdFlags)
	}

	// Determine JWK URL
	jwksURL := webhookJWKsURL
	if webhookCanvasDataSvc {
		jwksURL = webhook.CanvasDataServicesJWKURL
		logger.Printf("Using Canvas Data Services JWK URL: %s\n", jwksURL)
	}

	// Create middleware
	middleware := []webhook.Middleware{
		webhook.RecoveryMiddleware(logger),
	}
	if webhookLog {
		middleware = append(middleware, webhook.EventLogger(logger))
	}

	// Create listener
	listener := webhook.New(&webhook.Config{
		Addr:       webhookAddr,
		Secret:     webhookSecret,
		JWKSetURL:  jwksURL,
		Logger:     logger,
		Middleware: middleware,
	})

	// Register handlers
	if len(webhookEvents) > 0 {
		// Register specific event types
		for _, eventType := range webhookEvents {
			registerHandler(listener, eventType, logger)
		}
	} else {
		// Register all event types
		for _, eventType := range webhook.AllEventTypes() {
			registerHandler(listener, eventType, logger)
		}
	}

	// Print stats
	listener.PrintStats()

	// Setup signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Println("Received shutdown signal")
		cancel()
	}()

	// Start listener
	if err := listener.Start(ctx); err != nil {
		return fmt.Errorf("webhook listener error: %w", err)
	}

	logger.Println("Webhook listener stopped")
	return nil
}

func runWebhookEvents(cmd *cobra.Command, args []string) {
	fmt.Println("Available Canvas Webhook Event Types:")

	events := webhook.AllEventTypes()
	for _, eventType := range events {
		name := webhook.GetEventName(eventType)
		fmt.Printf("  %-35s %s\n", eventType, name)
	}

	fmt.Printf("\nTotal: %d event types\n", len(events))
}

func registerHandler(listener *webhook.Listener, eventType string, logger *log.Logger) {
	listener.On(eventType, func(ctx context.Context, event *webhook.Event) error {
		logger.Printf("Received event: %s (ID: %s)\n", event.EventType, event.ID)

		// Print event body
		if webhookLog {
			logger.Printf("Event body: %+v\n", event.Body)
		}

		return nil
	})
}
