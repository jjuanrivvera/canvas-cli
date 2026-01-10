package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	listener := New(&Config{
		Addr:   "localhost:8080",
		Secret: "my-secret",
	})

	if listener == nil {
		t.Fatal("expected non-nil listener")
	}

	if listener.addr != "localhost:8080" {
		t.Errorf("expected addr 'localhost:8080', got '%s'", listener.addr)
	}

	if listener.secret != "my-secret" {
		t.Errorf("expected secret 'my-secret', got '%s'", listener.secret)
	}

	if listener.handlers == nil {
		t.Error("expected handlers map to be initialized")
	}
}

func TestListener_On(t *testing.T) {
	listener := New(&Config{
		Addr: "localhost:8080",
	})

	called := false
	handler := func(ctx context.Context, event *Event) error {
		called = true
		return nil
	}

	listener.On("assignment_created", handler)

	// Verify handler was registered
	if len(listener.handlers["assignment_created"]) != 1 {
		t.Error("expected handler to be registered")
	}

	// Simulate calling the handler
	err := listener.handlers["assignment_created"][0](context.Background(), &Event{})
	if err != nil {
		t.Errorf("handler returned error: %v", err)
	}

	if !called {
		t.Error("expected handler to be called")
	}
}

func TestListener_OnMultipleHandlers(t *testing.T) {
	listener := New(&Config{
		Addr: "localhost:8080",
	})

	count := 0
	handler1 := func(ctx context.Context, event *Event) error {
		count++
		return nil
	}
	handler2 := func(ctx context.Context, event *Event) error {
		count++
		return nil
	}

	listener.On("assignment_created", handler1)
	listener.On("assignment_created", handler2)

	// Both handlers should be registered
	if len(listener.handlers["assignment_created"]) != 2 {
		t.Errorf("expected 2 handlers, got %d", len(listener.handlers["assignment_created"]))
	}

	// Call both handlers
	for _, h := range listener.handlers["assignment_created"] {
		h(context.Background(), &Event{})
	}

	if count != 2 {
		t.Errorf("expected count 2, got %d", count)
	}
}

func TestListener_verifySignature(t *testing.T) {
	secret := "test-secret"
	listener := New(&Config{
		Addr:   "localhost:8080",
		Secret: secret,
	})

	body := []byte(`{"event_type":"assignment_created"}`)

	// Generate valid signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	validSignature := hex.EncodeToString(mac.Sum(nil))

	if !listener.verifySignature(body, validSignature) {
		t.Error("expected valid signature to pass verification")
	}

	// Test invalid signature
	if listener.verifySignature(body, "invalid-signature") {
		t.Error("expected invalid signature to fail verification")
	}

	// Test empty signature
	if listener.verifySignature(body, "") {
		t.Error("expected empty signature to fail verification")
	}
}

func TestListener_verifySignature_NoSecret(t *testing.T) {
	listener := New(&Config{
		Addr: "localhost:8080",
	})

	body := []byte(`{"event_type":"assignment_created"}`)

	// Should always return true when no secret is configured
	if !listener.verifySignature(body, "") {
		t.Error("expected verification to pass when no secret is configured")
	}

	if !listener.verifySignature(body, "any-signature") {
		t.Error("expected verification to pass when no secret is configured")
	}
}

func TestListener_handleWebhook_Success(t *testing.T) {
	listener := New(&Config{
		Addr: "localhost:8080",
	})

	eventReceived := false
	listener.On("assignment_created", func(ctx context.Context, event *Event) error {
		eventReceived = true
		if event.EventType != "assignment_created" {
			t.Errorf("expected event type 'assignment_created', got '%s'", event.EventType)
		}
		return nil
	})

	// Create test event
	eventData := map[string]interface{}{
		"event_type": "assignment_created",
		"assignment": map[string]interface{}{
			"id":   123,
			"name": "Test Assignment",
		},
	}
	body, _ := json.Marshal(eventData)

	// Create HTTP request
	req := httptest.NewRequest("POST", "/webhook", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Handle webhook
	listener.handleWebhook(w, req)

	// Verify response
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Give handler time to execute
	time.Sleep(50 * time.Millisecond)

	if !eventReceived {
		t.Error("expected event handler to be called")
	}
}

func TestListener_handleWebhook_WithSignature(t *testing.T) {
	secret := "test-secret"
	listener := New(&Config{
		Addr:   "localhost:8080",
		Secret: secret,
	})

	eventReceived := false
	listener.On("assignment_created", func(ctx context.Context, event *Event) error {
		eventReceived = true
		return nil
	})

	eventData := map[string]interface{}{
		"event_type": "assignment_created",
	}
	body, _ := json.Marshal(eventData)

	// Generate signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	signature := hex.EncodeToString(mac.Sum(nil))

	// Create request with signature
	req := httptest.NewRequest("POST", "/webhook", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Canvas-Signature", signature)

	w := httptest.NewRecorder()
	listener.handleWebhook(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	time.Sleep(50 * time.Millisecond)

	if !eventReceived {
		t.Error("expected event handler to be called")
	}
}

func TestListener_handleWebhook_InvalidSignature(t *testing.T) {
	secret := "test-secret"
	listener := New(&Config{
		Addr:   "localhost:8080",
		Secret: secret,
	})

	eventReceived := false
	listener.On("assignment_created", func(ctx context.Context, event *Event) error {
		eventReceived = true
		return nil
	})

	eventData := map[string]interface{}{
		"event_type": "assignment_created",
	}
	body, _ := json.Marshal(eventData)

	// Create request with invalid signature
	req := httptest.NewRequest("POST", "/webhook", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Canvas-Signature", "invalid-signature")

	w := httptest.NewRecorder()
	listener.handleWebhook(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}

	time.Sleep(50 * time.Millisecond)

	if eventReceived {
		t.Error("expected event handler NOT to be called with invalid signature")
	}
}

func TestListener_handleWebhook_NoHandler(t *testing.T) {
	listener := New(&Config{
		Addr: "localhost:8080",
	})

	// Don't register any handlers

	eventData := map[string]interface{}{
		"event_type": "assignment_created",
	}
	body, _ := json.Marshal(eventData)

	req := httptest.NewRequest("POST", "/webhook", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	listener.handleWebhook(w, req)

	// Should still return OK even if no handler
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestListener_handleWebhook_InvalidJSON(t *testing.T) {
	listener := New(&Config{
		Addr: "localhost:8080",
	})

	req := httptest.NewRequest("POST", "/webhook", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	listener.handleWebhook(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestListener_handleWebhook_WrongMethod(t *testing.T) {
	listener := New(&Config{
		Addr: "localhost:8080",
	})

	req := httptest.NewRequest("GET", "/webhook", nil)

	w := httptest.NewRecorder()
	listener.handleWebhook(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
	}
}

func TestListener_Start_Shutdown(t *testing.T) {
	listener := New(&Config{
		Addr: "localhost:18080",
	})

	ctx := context.Background()

	// Start listener in background
	errChan := make(chan error, 1)
	go func() {
		err := listener.Start(ctx)
		if err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Verify server is running
	resp, err := http.Get("http://localhost:18080/health")
	if err != nil {
		t.Fatalf("Failed to reach server: %v", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	// Shutdown listener
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = listener.Shutdown(shutdownCtx)
	if err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}

	// Verify server is stopped
	time.Sleep(100 * time.Millisecond)
	_, err = http.Get("http://localhost:18080/health")
	if err == nil {
		t.Error("expected server to be stopped")
	}

	// Check for startup errors
	select {
	case err := <-errChan:
		t.Fatalf("Server error: %v", err)
	default:
	}
}

func TestListener_handleHealth(t *testing.T) {
	listener := New(&Config{
		Addr: "localhost:8080",
	})

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	listener.handleHealth(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	body, _ := io.ReadAll(w.Body)
	var response map[string]string
	json.Unmarshal(body, &response)

	if response["status"] != "ok" {
		t.Errorf("expected status 'ok', got '%s'", response["status"])
	}
}

func TestListener_Stats(t *testing.T) {
	listener := New(&Config{
		Addr: "localhost:8080",
	})

	// Track some events
	listener.On("test_event", func(ctx context.Context, event *Event) error {
		return nil
	})

	stats := listener.Stats()

	if stats.HandlerCount == 0 {
		t.Error("expected handlers to be counted in stats")
	}
}

func TestListener_PrintStats(t *testing.T) {
	listener := New(&Config{
		Addr: "localhost:8080",
	})

	listener.On("test_event", func(ctx context.Context, event *Event) error {
		return nil
	})

	// Should not panic
	listener.PrintStats()
}

func TestEvent_Marshal(t *testing.T) {
	event := Event{
		EventType: "assignment_created",
		EventTime: time.Now(),
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal event: %v", err)
	}

	var unmarshaled Event
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal event: %v", err)
	}

	if unmarshaled.EventType != event.EventType {
		t.Errorf("expected event type '%s', got '%s'", event.EventType, unmarshaled.EventType)
	}
}

func TestListener_ConcurrentRequests(t *testing.T) {
	listener := New(&Config{
		Addr: "localhost:18081",
	})

	var count int64
	listener.On("test_event", func(ctx context.Context, event *Event) error {
		atomic.AddInt64(&count, 1)
		time.Sleep(10 * time.Millisecond)
		return nil
	})

	ctx := context.Background()

	// Start listener
	go listener.Start(ctx)
	time.Sleep(100 * time.Millisecond)

	// Send multiple concurrent requests
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			eventData := map[string]interface{}{
				"event_type": "test_event",
			}
			body, _ := json.Marshal(eventData)

			resp, err := http.Post("http://localhost:18081/webhook", "application/json", bytes.NewReader(body))
			if err == nil {
				resp.Body.Close()
			}
			done <- true
		}()
	}

	// Wait for all requests
	for i := 0; i < 10; i++ {
		<-done
	}

	// Shutdown listener
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	listener.Shutdown(shutdownCtx)

	// Give handlers time to complete
	time.Sleep(200 * time.Millisecond)

	finalCount := atomic.LoadInt64(&count)
	if finalCount != 10 {
		t.Errorf("expected count 10, got %d", finalCount)
	}
}

func TestGetEventName(t *testing.T) {
	tests := []struct {
		name      string
		eventType string
		expected  string
	}{
		{"assignment_created", "assignment_created", "Assignment Created"},
		{"submission_created", "submission_created", "Submission Created"},
		{"grade_change", "grade_change", "Grade Change"},
		{"unknown", "unknown_event", "unknown_event"}, // Returns the input for unknown events
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetEventName(tt.eventType)
			if result != tt.expected {
				t.Errorf("GetEventName(%s) = %s, want %s", tt.eventType, result, tt.expected)
			}
		})
	}
}

func TestIsValidEventType(t *testing.T) {
	tests := []struct {
		name      string
		eventType string
		expected  bool
	}{
		{"valid_assignment", "assignment_created", true},
		{"valid_submission", "submission_created", true},
		{"valid_grade", "grade_change", true},
		{"invalid", "invalid_event", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidEventType(tt.eventType)
			if result != tt.expected {
				t.Errorf("IsValidEventType(%s) = %v, want %v", tt.eventType, result, tt.expected)
			}
		})
	}
}

func TestAllEventTypes(t *testing.T) {
	eventTypes := AllEventTypes()

	if len(eventTypes) == 0 {
		t.Error("expected non-empty event types slice")
	}

	// Check for some known event types
	expectedTypes := []string{
		"assignment_created",
		"submission_created",
		"grade_change",
	}

	for _, expected := range expectedTypes {
		found := false
		for _, eventType := range eventTypes {
			if eventType == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected event type %s not found in AllEventTypes()", expected)
		}
	}
}

func TestEventLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	middleware := EventLogger(logger)

	called := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	wrappedHandler := middleware(nextHandler)

	req := httptest.NewRequest("POST", "/webhook", nil)
	rec := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rec, req)

	if !called {
		t.Error("expected handler to be called")
	}

	logOutput := buf.String()
	if logOutput == "" {
		t.Error("expected log output")
	}
}

func TestRecoveryMiddleware(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	middleware := RecoveryMiddleware(logger)

	// Test handler that panics
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	wrappedHandler := middleware(panicHandler)

	req := httptest.NewRequest("POST", "/webhook", nil)
	rec := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rec.Code)
	}

	logOutput := buf.String()
	if logOutput == "" {
		t.Error("expected log output from panic recovery")
	}
}
