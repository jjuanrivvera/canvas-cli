package webhook

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// Event represents a Canvas webhook event
type Event struct {
	ID        string                 `json:"id"`
	EventType string                 `json:"event_type"`
	EventTime time.Time              `json:"event_time"`
	Body      map[string]interface{} `json:"body"`
}

// Handler is a function that processes webhook events
type Handler func(ctx context.Context, event *Event) error

// Listener represents a webhook listener server
type Listener struct {
	addr       string
	secret     string
	handlers   map[string][]Handler
	mu         sync.RWMutex
	server     *http.Server
	logger     *log.Logger
	middleware []Middleware
}

// Middleware is a function that wraps an HTTP handler
type Middleware func(http.Handler) http.Handler

// Config represents webhook listener configuration
type Config struct {
	Addr       string
	Secret     string
	Logger     *log.Logger
	Middleware []Middleware
}

// New creates a new webhook listener
func New(cfg *Config) *Listener {
	if cfg.Logger == nil {
		cfg.Logger = log.New(io.Discard, "", 0)
	}

	l := &Listener{
		addr:       cfg.Addr,
		secret:     cfg.Secret,
		handlers:   make(map[string][]Handler),
		logger:     cfg.Logger,
		middleware: cfg.Middleware,
	}

	return l
}

// On registers a handler for a specific event type
func (l *Listener) On(eventType string, handler Handler) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.handlers[eventType] = append(l.handlers[eventType], handler)
}

// Start starts the webhook listener server
func (l *Listener) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", l.handleWebhook)
	mux.HandleFunc("/health", l.handleHealth)

	// Apply middleware
	var handler http.Handler = mux
	for i := len(l.middleware) - 1; i >= 0; i-- {
		handler = l.middleware[i](handler)
	}

	l.server = &http.Server{
		Addr:    l.addr,
		Handler: handler,
	}

	l.logger.Printf("Starting webhook listener on %s\n", l.addr)

	// Start server in goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := l.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Wait for context cancellation or error
	select {
	case <-ctx.Done():
		return l.Shutdown(context.Background())
	case err := <-errChan:
		return err
	}
}

// Shutdown gracefully shuts down the webhook listener
func (l *Listener) Shutdown(ctx context.Context) error {
	if l.server == nil {
		return nil
	}

	l.logger.Println("Shutting down webhook listener...")
	return l.server.Shutdown(ctx)
}

// handleWebhook handles incoming webhook requests
func (l *Listener) handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		l.logger.Printf("Failed to read request body: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Verify signature if secret is configured
	if l.secret != "" {
		signature := r.Header.Get("X-Canvas-Signature")
		if !l.verifySignature(body, signature) {
			l.logger.Println("Invalid webhook signature")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	// Parse event
	var event Event
	if err := json.Unmarshal(body, &event); err != nil {
		l.logger.Printf("Failed to parse webhook event: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Process event
	l.mu.RLock()
	handlers := l.handlers[event.EventType]
	l.mu.RUnlock()

	if len(handlers) == 0 {
		l.logger.Printf("No handlers registered for event type: %s\n", event.EventType)
		w.WriteHeader(http.StatusOK)
		return
	}

	// Execute handlers
	ctx := r.Context()
	for _, handler := range handlers {
		if err := handler(ctx, &event); err != nil {
			l.logger.Printf("Handler error for event %s: %v\n", event.EventType, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// handleHealth handles health check requests
func (l *Listener) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// verifySignature verifies the webhook signature
func (l *Listener) verifySignature(body []byte, signature string) bool {
	// If no secret is configured, always pass verification
	if l.secret == "" {
		return true
	}

	// If secret is configured but no signature provided, fail
	if signature == "" {
		return false
	}

	// Calculate HMAC
	mac := hmac.New(sha256.New, []byte(l.secret))
	mac.Write(body)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}

// EventLogger is a middleware that logs webhook events
func EventLogger(logger *log.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Call next handler
			next.ServeHTTP(w, r)

			// Log request
			logger.Printf("%s %s %s\n", r.Method, r.URL.Path, time.Since(start))
		})
	}
}

// RecoveryMiddleware recovers from panics and returns 500 error
func RecoveryMiddleware(logger *log.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Printf("Panic recovered: %v\n", err)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// Stats represents webhook listener statistics
type Stats struct {
	EventsReceived  int64
	EventsProcessed int64
	EventsFailed    int64
	HandlerCount    int
}

// Stats returns current listener statistics
func (l *Listener) Stats() Stats {
	l.mu.RLock()
	defer l.mu.RUnlock()

	handlerCount := 0
	for _, handlers := range l.handlers {
		handlerCount += len(handlers)
	}

	return Stats{
		HandlerCount: handlerCount,
	}
}

// PrintStats prints listener statistics
func (l *Listener) PrintStats() {
	stats := l.Stats()
	fmt.Printf("Webhook Listener Statistics:\n")
	fmt.Printf("  Handlers Registered: %d\n", stats.HandlerCount)
}
