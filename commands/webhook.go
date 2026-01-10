package commands

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jjuanrivvera/canvas-cli/internal/webhook"
	"github.com/spf13/cobra"
)

var (
	webhookAddr   string
	webhookSecret string
	webhookEvents []string
	webhookLog    bool
)

var webhookCmd = &cobra.Command{
	Use:   "webhook",
	Short: "Manage Canvas webhook listeners",
	Long: `Start and manage webhook listeners for Canvas LMS events.

The webhook listener receives real-time notifications from Canvas for
various events like assignment creation, submission updates, grade changes,
enrollment changes, and more.

The listener provides:
  - HMAC signature verification for security
  - Event routing to specific handlers
  - Graceful shutdown
  - Health check endpoint
  - Request logging

Examples:
  # Start webhook listener on default port
  canvas webhook listen

  # Start with custom address and secret
  canvas webhook listen --addr :8080 --secret your-secret-key

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
webhook events. It verifies HMAC signatures if a secret is provided
and routes events to registered handlers.

Endpoints:
  POST /webhook - Receive webhook events
  GET  /health  - Health check endpoint

Examples:
  # Start listener on default address
  canvas webhook listen

  # Start with custom configuration
  canvas webhook listen --addr :9000 --secret my-secret

  # Enable request logging
  canvas webhook listen --log

  # Listen for specific events
  canvas webhook listen --events submission_created,assignment_updated`,
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
