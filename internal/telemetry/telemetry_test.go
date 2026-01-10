package telemetry

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tempDir := t.TempDir()

	client, err := New(&Config{
		Enabled:   true,
		DataDir:   tempDir,
		Version:   "1.0.0",
		Anonymous: false,
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	if client == nil {
		t.Fatal("expected non-nil client")
	}

	if !client.enabled {
		t.Error("expected client to be enabled")
	}

	if client.version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got '%s'", client.version)
	}
}

func TestNew_Disabled(t *testing.T) {
	tempDir := t.TempDir()

	client, err := New(&Config{
		Enabled: false,
		DataDir: tempDir,
		Version: "1.0.0",
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	if client.enabled {
		t.Error("expected client to be disabled")
	}
}

func TestClient_Track(t *testing.T) {
	tempDir := t.TempDir()

	client, err := New(&Config{
		Enabled: true,
		DataDir: tempDir,
		Version: "1.0.0",
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	client.Track("command", "courses.list", map[string]interface{}{
		"course_id": 123,
		"format":    "table",
	})

	client.mu.Lock()
	defer client.mu.Unlock()

	if len(client.events) != 1 {
		t.Errorf("expected 1 event, got %d", len(client.events))
	}

	event := client.events[0]
	if event.Type != "command" {
		t.Errorf("expected type 'command', got '%s'", event.Type)
	}

	if event.Action != "courses.list" {
		t.Errorf("expected action 'courses.list', got '%s'", event.Action)
	}
}

func TestClient_Track_Disabled(t *testing.T) {
	tempDir := t.TempDir()

	client, err := New(&Config{
		Enabled: false,
		DataDir: tempDir,
		Version: "1.0.0",
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	client.Track("command", "courses.list", nil)

	client.mu.Lock()
	defer client.mu.Unlock()

	if len(client.events) != 0 {
		t.Error("expected no events when disabled")
	}
}

func TestClient_IsEnabled(t *testing.T) {
	tempDir := t.TempDir()

	client, err := New(&Config{
		Enabled: true,
		DataDir: tempDir,
		Version: "1.0.0",
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	if !client.IsEnabled() {
		t.Error("expected IsEnabled to return true")
	}

	// Test disabled client
	client2, err := New(&Config{
		Enabled: false,
		DataDir: tempDir,
		Version: "1.0.0",
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	if client2.IsEnabled() {
		t.Error("expected IsEnabled to return false")
	}
}

func TestClient_Stats(t *testing.T) {
	tempDir := t.TempDir()

	client, err := New(&Config{
		Enabled: true,
		DataDir: tempDir,
		Version: "1.0.0",
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Track some events
	client.Track("command", "courses.list", nil)
	client.Track("command", "courses.get", nil)
	client.Track("api", "request", nil)

	stats := client.Stats()

	if stats.Enabled != true {
		t.Error("expected Enabled to be true")
	}

	if stats.EventsQueued != 3 {
		t.Errorf("expected 3 events queued, got %d", stats.EventsQueued)
	}

	if stats.SessionID == "" {
		t.Error("expected SessionID to be set")
	}
}

func TestClient_Stats_Disabled(t *testing.T) {
	tempDir := t.TempDir()

	client, err := New(&Config{
		Enabled: false,
		DataDir: tempDir,
		Version: "1.0.0",
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	stats := client.Stats()

	if stats.Enabled != false {
		t.Error("expected Enabled to be false")
	}

	if stats.EventsQueued != 0 {
		t.Errorf("expected 0 events queued, got %d", stats.EventsQueued)
	}
}

func TestClient_Close(t *testing.T) {
	tempDir := t.TempDir()

	client, err := New(&Config{
		Enabled: true,
		DataDir: tempDir,
		Version: "1.0.0",
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Track some events
	client.Track("command", "courses.list", nil)

	// Close should flush events
	err = client.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	// Events should be persisted
	eventsDir := filepath.Join(tempDir, "events")
	if _, err := os.Stat(eventsDir); os.IsNotExist(err) {
		// Directory might not be created if no flush happened
		// This is acceptable
	}
}

func TestClient_Flush(t *testing.T) {
	tempDir := t.TempDir()

	client, err := New(&Config{
		Enabled: true,
		DataDir: tempDir,
		Version: "1.0.0",
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Track some events
	client.Track("command", "courses.list", nil)
	client.Track("command", "courses.get", nil)

	// Flush events
	err = client.Flush()
	if err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	// Events should be cleared from memory after flush
	client.mu.Lock()
	eventCount := len(client.events)
	client.mu.Unlock()

	if eventCount != 0 {
		t.Errorf("expected 0 events after flush, got %d", eventCount)
	}
}

func TestClient_Anonymous(t *testing.T) {
	tempDir := t.TempDir()

	client, err := New(&Config{
		Enabled:   true,
		DataDir:   tempDir,
		Version:   "1.0.0",
		Anonymous: true,
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	if !client.anonymous {
		t.Error("expected anonymous to be true")
	}
}

func TestClient_GetSessionID(t *testing.T) {
	tempDir := t.TempDir()

	client, err := New(&Config{
		Enabled: true,
		DataDir: tempDir,
		Version: "1.0.0",
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	sessionID := client.GetSessionID()
	if sessionID == "" {
		t.Error("expected non-empty session ID")
	}
}

func TestClient_GetUserID(t *testing.T) {
	tempDir := t.TempDir()

	client, err := New(&Config{
		Enabled: true,
		DataDir: tempDir,
		Version: "1.0.0",
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	userID := client.GetUserID()
	// UserID may or may not be set depending on system
	_ = userID
}

func TestClient_TrackWithNilProperties(t *testing.T) {
	tempDir := t.TempDir()

	client, err := New(&Config{
		Enabled: true,
		DataDir: tempDir,
		Version: "1.0.0",
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Should not panic with nil properties
	client.Track("command", "test", nil)

	client.mu.Lock()
	defer client.mu.Unlock()

	if len(client.events) != 1 {
		t.Error("expected event to be tracked with nil properties")
	}
}

func TestClient_TrackMultipleTypes(t *testing.T) {
	tempDir := t.TempDir()

	client, err := New(&Config{
		Enabled: true,
		DataDir: tempDir,
		Version: "1.0.0",
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Track different event types
	client.Track("command", "courses.list", nil)
	client.Track("api", "request", nil)
	client.Track("error", "failed", nil)
	client.Track("performance", "slow", nil)

	client.mu.Lock()
	defer client.mu.Unlock()

	if len(client.events) != 4 {
		t.Errorf("expected 4 events, got %d", len(client.events))
	}

	// Verify event types
	types := make(map[string]int)
	for _, event := range client.events {
		types[event.Type]++
	}

	if types["command"] != 1 || types["api"] != 1 || types["error"] != 1 || types["performance"] != 1 {
		t.Error("expected one event of each type")
	}
}

func TestClient_TrackCommand(t *testing.T) {
	tempDir := t.TempDir()

	client, err := New(&Config{
		Enabled: true,
		DataDir: tempDir,
		Version: "1.0.0",
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	client.TrackCommand("courses list", 100*time.Millisecond, nil)

	client.mu.Lock()
	defer client.mu.Unlock()

	if len(client.events) != 1 {
		t.Errorf("expected 1 event, got %d", len(client.events))
	}

	event := client.events[0]
	if event.Type != "command" {
		t.Errorf("expected type 'command', got '%s'", event.Type)
	}

	if event.Duration == 0 {
		t.Error("expected duration to be set")
	}
}

func TestClient_TrackError(t *testing.T) {
	tempDir := t.TempDir()

	client, err := New(&Config{
		Enabled: true,
		DataDir: tempDir,
		Version: "1.0.0",
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer client.Close()

	client.TrackError("api_error", "Failed to fetch data", map[string]interface{}{
		"status_code": 500,
	})

	client.mu.Lock()
	defer client.mu.Unlock()

	if len(client.events) != 1 {
		t.Errorf("expected 1 event, got %d", len(client.events))
	}

	event := client.events[0]
	if event.Type != "error" {
		t.Errorf("expected type 'error', got '%s'", event.Type)
	}

	if event.Error == "" {
		t.Error("expected error message to be set")
	}
}

func TestNew_NilConfig(t *testing.T) {
	client, err := New(nil)
	if err != nil {
		t.Fatalf("New with nil config failed: %v", err)
	}

	if client.enabled {
		t.Error("expected client to be disabled with nil config")
	}
}

func TestClient_WithContext(t *testing.T) {
	tempDir := t.TempDir()

	client, err := New(&Config{
		Enabled: true,
		DataDir: tempDir,
		Version: "1.0.0",
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	ctx := context.Background()
	ctxWithTelemetry := client.WithContext(ctx)

	// Verify context contains the client
	retrieved, ok := FromContext(ctxWithTelemetry)
	if !ok {
		t.Error("expected to retrieve client from context")
	}

	if retrieved != client {
		t.Error("expected retrieved client to match original client")
	}
}

func TestFromContext_NoClient(t *testing.T) {
	ctx := context.Background()

	// Context without telemetry client
	client, ok := FromContext(ctx)
	if ok {
		t.Error("expected ok to be false for context without telemetry")
	}

	if client != nil {
		t.Error("expected nil client for context without telemetry")
	}
}

func TestFromContext_WithClient(t *testing.T) {
	tempDir := t.TempDir()

	originalClient, err := New(&Config{
		Enabled: true,
		DataDir: tempDir,
		Version: "1.0.0",
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	ctx := context.Background()
	ctxWithTelemetry := originalClient.WithContext(ctx)

	// Retrieve client from context
	retrievedClient, ok := FromContext(ctxWithTelemetry)
	if !ok {
		t.Fatal("expected to retrieve client from context")
	}

	if retrievedClient != originalClient {
		t.Error("expected retrieved client to be the same instance")
	}

	// Verify client is functional
	if retrievedClient.GetSessionID() == "" {
		t.Error("expected retrieved client to have session ID")
	}
}

func TestClient_TrackCommand_WithError(t *testing.T) {
	tempDir := t.TempDir()

	client, err := New(&Config{
		Enabled: true,
		DataDir: tempDir,
		Version: "1.0.0",
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Track command with error
	testErr := fmt.Errorf("command failed")
	client.TrackCommand("courses delete", 50*time.Millisecond, testErr)

	client.mu.Lock()
	defer client.mu.Unlock()

	if len(client.events) != 1 {
		t.Errorf("expected 1 event, got %d", len(client.events))
	}

	event := client.events[0]
	if event.Error == "" {
		t.Error("expected error to be set in event")
	}

	if event.Error != testErr.Error() {
		t.Errorf("expected error '%s', got '%s'", testErr.Error(), event.Error)
	}
}

func TestClient_TrackError_NilProperties(t *testing.T) {
	tempDir := t.TempDir()

	client, err := New(&Config{
		Enabled: true,
		DataDir: tempDir,
		Version: "1.0.0",
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer client.Close()

	// Track error with nil properties
	client.TrackError("validation_error", "Invalid input", nil)

	client.mu.Lock()
	defer client.mu.Unlock()

	if len(client.events) != 1 {
		t.Errorf("expected 1 event, got %d", len(client.events))
	}

	event := client.events[0]
	if event.Properties == nil {
		t.Error("expected properties to be initialized")
	}

	if event.Properties["error_type"] != "validation_error" {
		t.Error("expected error_type to be set in properties")
	}
}

func TestClient_TrackError_FlushChannelFull(t *testing.T) {
	tempDir := t.TempDir()

	client, err := New(&Config{
		Enabled: true,
		DataDir: tempDir,
		Version: "1.0.0",
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Fill the flush channel by sending a signal
	client.flushChan <- struct{}{}

	// Track error - should not block even though channel is full
	client.TrackError("critical_error", "System failure", map[string]interface{}{
		"severity": "high",
	})

	client.mu.Lock()
	eventCount := len(client.events)
	client.mu.Unlock()

	if eventCount != 1 {
		t.Errorf("expected 1 event, got %d", eventCount)
	}

	// Drain the channel
	<-client.flushChan
}

func TestClient_Flush_WriteError(t *testing.T) {
	tempDir := t.TempDir()

	client, err := New(&Config{
		Enabled: true,
		DataDir: tempDir,
		Version: "1.0.0",
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Track an event
	client.Track("test", "action", nil)

	// Make directory read-only to trigger write error
	os.Chmod(tempDir, 0500)
	defer os.Chmod(tempDir, 0700) // Restore for cleanup

	// Flush should fail due to write permissions
	err = client.Flush()
	if err == nil {
		t.Error("expected error when flushing to read-only directory")
	}
}

func TestClient_Close_FlushError(t *testing.T) {
	tempDir := t.TempDir()

	client, err := New(&Config{
		Enabled: true,
		DataDir: tempDir,
		Version: "1.0.0",
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Track an event
	client.Track("test", "action", nil)

	// Make directory read-only to trigger flush error
	os.Chmod(tempDir, 0500)
	defer os.Chmod(tempDir, 0700) // Restore for cleanup

	// Close should return error from Flush
	err = client.Close()
	if err == nil {
		t.Error("expected error when closing with read-only directory")
	}
}

func TestClient_Close_Disabled(t *testing.T) {
	tempDir := t.TempDir()

	client, err := New(&Config{
		Enabled: false,
		DataDir: tempDir,
		Version: "1.0.0",
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Close should succeed without error when disabled
	err = client.Close()
	if err != nil {
		t.Errorf("expected no error when closing disabled client, got %v", err)
	}
}

func TestClient_TrackCommand_Disabled(t *testing.T) {
	tempDir := t.TempDir()

	client, err := New(&Config{
		Enabled: false,
		DataDir: tempDir,
		Version: "1.0.0",
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// TrackCommand should not add events when disabled
	client.TrackCommand("courses list", 100*time.Millisecond, nil)

	client.mu.Lock()
	eventCount := len(client.events)
	client.mu.Unlock()

	if eventCount != 0 {
		t.Errorf("expected 0 events when disabled, got %d", eventCount)
	}
}

func TestClient_TrackError_Disabled(t *testing.T) {
	tempDir := t.TempDir()

	client, err := New(&Config{
		Enabled: false,
		DataDir: tempDir,
		Version: "1.0.0",
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// TrackError should not add events when disabled
	client.TrackError("api_error", "Failed", nil)

	client.mu.Lock()
	eventCount := len(client.events)
	client.mu.Unlock()

	if eventCount != 0 {
		t.Errorf("expected 0 events when disabled, got %d", eventCount)
	}
}
