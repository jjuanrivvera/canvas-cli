package commands

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/webhook"
)

// webhookCmd represents the webhook command group
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

func init() {
	rootCmd.AddCommand(webhookCmd)
	webhookCmd.AddCommand(newWebhookListenCmd())
	webhookCmd.AddCommand(newWebhookEventsCmd())
}

func newWebhookListenCmd() *cobra.Command {
	opts := &options.WebhookListenOptions{}

	cmd := &cobra.Command{
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
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			return runWebhookListen(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.Addr, "addr", ":8080", "Server address to listen on")
	cmd.Flags().StringVar(&opts.Secret, "secret", "", "Webhook secret for HMAC verification")
	cmd.Flags().StringVar(&opts.JWKsURL, "jwks-url", "", "JWK Set URL for JWT verification")
	cmd.Flags().BoolVar(&opts.CanvasDataSvc, "canvas-data-services", false, "Use Canvas Data Services JWK URL for JWT verification")
	cmd.Flags().StringSliceVar(&opts.Events, "events", []string{}, "Event types to handle (comma-separated)")
	cmd.Flags().BoolVar(&opts.Log, "log", false, "Enable request logging")

	return cmd
}

func newWebhookEventsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "events",
		Short: "List available Canvas webhook event types",
		Long:  `Display all supported Canvas webhook event types with their descriptions.`,
		Run: func(cmd *cobra.Command, args []string) {
			runWebhookEvents()
		},
	}

	return cmd
}

func runWebhookListen(ctx context.Context, opts *options.WebhookListenOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "webhook.listen", map[string]interface{}{
		"addr":            opts.Addr,
		"canvas_data_svc": opts.CanvasDataSvc,
		"events_count":    len(opts.Events),
		"log_enabled":     opts.Log,
	})

	// Create webhook logger
	var webhookLogger *log.Logger
	if opts.Log {
		webhookLogger = log.New(os.Stdout, "[webhook] ", log.LstdFlags)
	} else {
		webhookLogger = log.New(os.Stderr, "", log.LstdFlags)
	}

	// Determine JWK URL
	jwksURL := opts.JWKsURL
	if opts.CanvasDataSvc {
		jwksURL = webhook.CanvasDataServicesJWKURL
		webhookLogger.Printf("Using Canvas Data Services JWK URL: %s\n", jwksURL)
	}

	// Create middleware
	middleware := []webhook.Middleware{
		webhook.RecoveryMiddleware(webhookLogger),
	}
	if opts.Log {
		middleware = append(middleware, webhook.EventLogger(webhookLogger))
	}

	// Create listener
	listener := webhook.New(&webhook.Config{
		Addr:       opts.Addr,
		Secret:     opts.Secret,
		JWKSetURL:  jwksURL,
		Logger:     webhookLogger,
		Middleware: middleware,
	})

	// Register handlers
	if len(opts.Events) > 0 {
		// Register specific event types
		for _, eventType := range opts.Events {
			registerHandler(listener, eventType, webhookLogger, opts.Log)
		}
	} else {
		// Register all event types
		for _, eventType := range webhook.AllEventTypes() {
			registerHandler(listener, eventType, webhookLogger, opts.Log)
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
		webhookLogger.Println("Received shutdown signal")
		cancel()
	}()

	// Start listener
	if err := listener.Start(ctx); err != nil {
		logger.LogCommandError(ctx, "webhook.listen", err, map[string]interface{}{
			"addr": opts.Addr,
		})
		return fmt.Errorf("webhook listener error: %w", err)
	}

	webhookLogger.Println("Webhook listener stopped")
	logger.LogCommandComplete(ctx, "webhook.listen", 0)
	return nil
}

func runWebhookEvents() {
	fmt.Println("Available Canvas Webhook Event Types:")

	events := webhook.AllEventTypes()
	for _, eventType := range events {
		name := webhook.GetEventName(eventType)
		fmt.Printf("  %-35s %s\n", eventType, name)
	}

	fmt.Printf("\nTotal: %d event types\n", len(events))
}

func registerHandler(listener *webhook.Listener, eventType string, logger *log.Logger, logEnabled bool) {
	listener.On(eventType, func(ctx context.Context, event *webhook.Event) error {
		logger.Printf("Received event: %s (ID: %s)\n", event.EventType, event.ID)

		// Print event body
		if logEnabled {
			logger.Printf("Event body: %+v\n", event.Body)
		}

		return nil
	})
}
