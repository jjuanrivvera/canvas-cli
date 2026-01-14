package webhook

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	jwkSet     *JWKSet
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
	JWKSetURL  string // URL to fetch JWKs for JWT verification (Canvas Data Services)
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

	// Initialize JWK set if URL is provided
	if cfg.JWKSetURL != "" {
		l.jwkSet = NewJWKSet(cfg.JWKSetURL)
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

	var event Event
	var verified bool

	// Try JWT verification first if JWK set is configured
	if l.jwkSet != nil {
		parsedEvent, err := l.verifyJWT(body)
		if err == nil {
			event = *parsedEvent
			verified = true
			l.logger.Println("JWT verification successful")
		} else {
			l.logger.Printf("JWT verification failed: %v\n", err)
			// Fall through to try HMAC or raw JSON
		}
	}

	// Try HMAC verification if secret is configured and JWT didn't work
	if !verified && l.secret != "" {
		signature := r.Header.Get("X-Canvas-Signature")
		if !l.verifySignature(body, signature) {
			l.logger.Println("Invalid webhook signature")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		verified = true
	}

	// If no verification method succeeded and one was configured, reject
	if !verified && (l.jwkSet != nil || l.secret != "") {
		l.logger.Println("No valid verification method succeeded")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse raw JSON event if not already parsed from JWT
	if event.EventType == "" {
		// Log raw body for debugging
		bodyStr := string(body)
		if len(bodyStr) > 200 {
			l.logger.Printf("Raw body (truncated): %s...\n", bodyStr[:200])
		} else {
			l.logger.Printf("Raw body: %s\n", bodyStr)
		}

		// Canvas Data Services sends JWT wrapped in quotes (e.g., "eyJhbGci...")
		// Strip the quotes to get the raw JWT
		if strings.HasPrefix(bodyStr, "\"") && strings.HasSuffix(bodyStr, "\"") {
			bodyStr = strings.Trim(bodyStr, "\"")
			body = []byte(bodyStr)
			l.logger.Println("Stripped quotes from JWT")
		}

		// Check if body looks like a JWT (3 base64 segments separated by dots)
		if strings.Count(bodyStr, ".") == 2 && !strings.HasPrefix(bodyStr, "{") {
			l.logger.Println("Body appears to be a JWT, attempting to parse claims without verification")
			parsedEvent, err := l.parseJWTClaims(body)
			if err != nil {
				l.logger.Printf("Failed to parse JWT claims: %v\n", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			event = *parsedEvent
		} else if err := json.Unmarshal(body, &event); err != nil {
			l.logger.Printf("Failed to parse webhook event: %v\n", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
	}

	l.logger.Printf("Received event: %s (ID: %s)\n", event.EventType, event.ID)
	l.logger.Printf("Event body: %v\n", event.Body)

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

// verifySignature verifies the webhook signature using HMAC-SHA256
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

// verifyJWT verifies a JWT-signed payload from Canvas Data Services
// and extracts the event data from the claims
func (l *Listener) verifyJWT(body []byte) (*Event, error) {
	if l.jwkSet == nil {
		return nil, errors.New("JWK set not configured")
	}

	tokenString := string(body)

	// Parse the JWT with key function to fetch the appropriate key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method is RSA
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Get key ID from header
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("missing kid in JWT header")
		}

		// Fetch the public key
		return l.jwkSet.GetKey(kid)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse/verify JWT: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid JWT token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to extract JWT claims")
	}

	// Build event from claims
	// Canvas Data Services embeds the event data in the JWT claims
	event := &Event{}

	// Canvas Data Services uses metadata.event_name for the event type
	if metadata, ok := claims["metadata"].(map[string]interface{}); ok {
		if eventName, ok := metadata["event_name"].(string); ok {
			event.EventType = eventName
		}
	}
	// Fallback to event_type if metadata.event_name not found
	if event.EventType == "" {
		if eventType, ok := claims["event_type"].(string); ok {
			event.EventType = eventType
		}
	}

	if id, ok := claims["id"].(string); ok {
		event.ID = id
	}

	if eventTime, ok := claims["event_time"].(string); ok {
		if t, err := time.Parse(time.RFC3339, eventTime); err == nil {
			event.EventTime = t
		}
	}

	// The body/data might be nested in different ways
	if body, ok := claims["body"].(map[string]interface{}); ok {
		event.Body = body
	} else if data, ok := claims["data"].(map[string]interface{}); ok {
		event.Body = data
	} else {
		// Use all claims as the body (excluding standard JWT claims)
		event.Body = make(map[string]interface{})
		for k, v := range claims {
			if k != "iss" && k != "sub" && k != "aud" && k != "exp" && k != "nbf" && k != "iat" && k != "jti" {
				event.Body[k] = v
			}
		}
	}

	return event, nil
}

// parseJWTClaims parses JWT claims WITHOUT verification
// This is useful for debugging or when verification is not required
func (l *Listener) parseJWTClaims(body []byte) (*Event, error) {
	tokenString := string(body)

	// Parse the JWT without verification
	parser := jwt.NewParser()
	token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT: %w", err)
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to extract JWT claims")
	}

	// Build event from claims (same logic as verifyJWT)
	event := &Event{}

	// Canvas Data Services uses metadata.event_name for the event type
	if metadata, ok := claims["metadata"].(map[string]interface{}); ok {
		if eventName, ok := metadata["event_name"].(string); ok {
			event.EventType = eventName
		}
	}
	// Fallback to event_type if metadata.event_name not found
	if event.EventType == "" {
		if eventType, ok := claims["event_type"].(string); ok {
			event.EventType = eventType
		}
	}

	if id, ok := claims["id"].(string); ok {
		event.ID = id
	}

	if eventTime, ok := claims["event_time"].(string); ok {
		if t, err := time.Parse(time.RFC3339, eventTime); err == nil {
			event.EventTime = t
		}
	}

	// The body/data might be nested in different ways
	if body, ok := claims["body"].(map[string]interface{}); ok {
		event.Body = body
	} else if data, ok := claims["data"].(map[string]interface{}); ok {
		event.Body = data
	} else {
		// Use all claims as the body (excluding standard JWT claims)
		event.Body = make(map[string]interface{})
		for k, v := range claims {
			if k != "iss" && k != "sub" && k != "aud" && k != "exp" && k != "nbf" && k != "iat" && k != "jti" {
				event.Body[k] = v
			}
		}
	}

	return event, nil
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
