package telemetry

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Client represents a telemetry client
type Client struct {
	enabled    bool
	userID     string
	sessionID  string
	events     []Event
	mu         sync.Mutex
	flushChan  chan struct{}
	stopChan   chan struct{}
	dataDir    string
	version    string
	anonymous  bool
}

// Event represents a telemetry event
type Event struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"user_id,omitempty"`
	SessionID string                 `json:"session_id"`
	Type      string                 `json:"type"`
	Action    string                 `json:"action"`
	Timestamp time.Time              `json:"timestamp"`
	Duration  time.Duration          `json:"duration,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Error     string                 `json:"error,omitempty"`
	OS        string                 `json:"os"`
	Arch      string                 `json:"arch"`
	Version   string                 `json:"version"`
}

// Config represents telemetry configuration
type Config struct {
	Enabled   bool
	DataDir   string
	Version   string
	Anonymous bool // If true, don't track user ID
}

// New creates a new telemetry client
func New(cfg *Config) (*Client, error) {
	if cfg == nil {
		cfg = &Config{
			Enabled: false,
		}
	}

	// Get data directory
	dataDir := cfg.DataDir
	if dataDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		dataDir = filepath.Join(home, ".canvas-cli", "telemetry")
	}

	// Create data directory
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create telemetry directory: %w", err)
	}

	// Generate or load user ID
	var userID string
	if !cfg.Anonymous {
		var err error
		userID, err = getUserID(dataDir)
		if err != nil {
			// Fall back to anonymous mode if we can't get/create user ID
			userID = ""
		}
	}

	client := &Client{
		enabled:    cfg.Enabled,
		userID:     userID,
		sessionID:  uuid.New().String(),
		events:     make([]Event, 0),
		flushChan:  make(chan struct{}, 1),
		stopChan:   make(chan struct{}),
		dataDir:    dataDir,
		version:    cfg.Version,
		anonymous:  cfg.Anonymous,
	}

	// Start background flush worker if enabled
	if cfg.Enabled {
		go client.flushWorker()
	}

	return client, nil
}

// getUserID gets or creates a persistent user ID
func getUserID(dataDir string) (string, error) {
	userIDPath := filepath.Join(dataDir, "user_id")

	// Try to read existing user ID
	data, err := os.ReadFile(userIDPath)
	if err == nil {
		return string(data), nil
	}

	// Generate new user ID
	userID := uuid.New().String()

	// Save user ID
	if err := os.WriteFile(userIDPath, []byte(userID), 0600); err != nil {
		return "", fmt.Errorf("failed to save user ID: %w", err)
	}

	return userID, nil
}

// Track records a telemetry event
func (c *Client) Track(eventType, action string, properties map[string]interface{}) {
	if !c.enabled {
		return
	}

	event := Event{
		ID:         uuid.New().String(),
		UserID:     c.userID,
		SessionID:  c.sessionID,
		Type:       eventType,
		Action:     action,
		Timestamp:  time.Now(),
		Properties: properties,
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
		Version:    c.version,
	}

	c.mu.Lock()
	c.events = append(c.events, event)
	c.mu.Unlock()

	// Trigger flush if we have enough events
	if len(c.events) >= 10 {
		select {
		case c.flushChan <- struct{}{}:
		default:
		}
	}
}

// TrackCommand tracks a command execution
func (c *Client) TrackCommand(command string, duration time.Duration, err error) {
	if !c.enabled {
		return
	}

	properties := map[string]interface{}{
		"command": command,
	}

	event := Event{
		ID:         uuid.New().String(),
		UserID:     c.userID,
		SessionID:  c.sessionID,
		Type:       "command",
		Action:     "execute",
		Timestamp:  time.Now(),
		Duration:   duration,
		Properties: properties,
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
		Version:    c.version,
	}

	if err != nil {
		event.Error = err.Error()
	}

	c.mu.Lock()
	c.events = append(c.events, event)
	c.mu.Unlock()
}

// TrackError tracks an error event
func (c *Client) TrackError(errorType, message string, properties map[string]interface{}) {
	if !c.enabled {
		return
	}

	if properties == nil {
		properties = make(map[string]interface{})
	}
	properties["error_type"] = errorType

	event := Event{
		ID:         uuid.New().String(),
		UserID:     c.userID,
		SessionID:  c.sessionID,
		Type:       "error",
		Action:     "occurred",
		Timestamp:  time.Now(),
		Properties: properties,
		Error:      message,
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
		Version:    c.version,
	}

	c.mu.Lock()
	c.events = append(c.events, event)
	c.mu.Unlock()

	// Flush errors immediately
	select {
	case c.flushChan <- struct{}{}:
	default:
	}
}

// Flush writes events to disk
func (c *Client) Flush() error {
	if !c.enabled {
		return nil
	}

	c.mu.Lock()
	events := make([]Event, len(c.events))
	copy(events, c.events)
	c.events = c.events[:0] // Clear events
	c.mu.Unlock()

	if len(events) == 0 {
		return nil
	}

	// Write events to file
	filename := filepath.Join(c.dataDir, fmt.Sprintf("events_%s.json", time.Now().Format("20060102_150405")))

	data, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal events: %w", err)
	}

	if err := os.WriteFile(filename, data, 0600); err != nil {
		return fmt.Errorf("failed to write events: %w", err)
	}

	return nil
}

// flushWorker periodically flushes events
func (c *Client) flushWorker() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.Flush()
		case <-c.flushChan:
			c.Flush()
		case <-c.stopChan:
			c.Flush() // Final flush
			return
		}
	}
}

// Close stops the telemetry client
func (c *Client) Close() error {
	if !c.enabled {
		return nil
	}

	close(c.stopChan)
	return c.Flush()
}

// IsEnabled returns whether telemetry is enabled
func (c *Client) IsEnabled() bool {
	return c.enabled
}

// GetSessionID returns the current session ID
func (c *Client) GetSessionID() string {
	return c.sessionID
}

// GetUserID returns the user ID (empty if anonymous)
func (c *Client) GetUserID() string {
	return c.userID
}

// Stats returns telemetry statistics
type Stats struct {
	Enabled       bool
	EventsQueued  int
	SessionID     string
	UserID        string
	Anonymous     bool
	DataDirectory string
}

// Stats returns current telemetry statistics
func (c *Client) Stats() Stats {
	c.mu.Lock()
	defer c.mu.Unlock()

	return Stats{
		Enabled:       c.enabled,
		EventsQueued:  len(c.events),
		SessionID:     c.sessionID,
		UserID:        c.userID,
		Anonymous:     c.anonymous,
		DataDirectory: c.dataDir,
	}
}

// WithContext wraps a context with telemetry tracking
func (c *Client) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, telemetryKey, c)
}

type contextKey string

const telemetryKey contextKey = "telemetry"

// FromContext retrieves the telemetry client from context
func FromContext(ctx context.Context) (*Client, bool) {
	client, ok := ctx.Value(telemetryKey).(*Client)
	return client, ok
}
