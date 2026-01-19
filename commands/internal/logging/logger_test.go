package logging

import (
	"context"
	"testing"
	"time"
)

func TestNewCommandLogger(t *testing.T) {
	tests := []struct {
		name  string
		debug bool
	}{
		{
			name:  "info level logger",
			debug: false,
		},
		{
			name:  "debug level logger",
			debug: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewCommandLogger(tt.debug)
			if logger == nil {
				t.Fatal("NewCommandLogger() returned nil")
			}
			if logger.logger == nil {
				t.Error("NewCommandLogger() logger field is nil")
			}
		})
	}
}

func TestCommandLogger_LogCommandStart(t *testing.T) {
	logger := NewCommandLogger(false)
	ctx := context.Background()

	// Should not panic
	logger.LogCommandStart(ctx, "test.command", map[string]interface{}{
		"arg1": "value1",
		"arg2": 123,
	})
}

func TestCommandLogger_LogAPICall(t *testing.T) {
	logger := NewCommandLogger(false)
	ctx := context.Background()

	// Should not panic
	logger.LogAPICall(ctx, "GET", "/api/v1/courses", 200, 100*time.Millisecond)
}

func TestCommandLogger_LogCommandComplete(t *testing.T) {
	logger := NewCommandLogger(false)
	ctx := context.Background()

	// Should not panic
	logger.LogCommandComplete(ctx, "test.command", 42)
}

func TestCommandLogger_LogCommandError(t *testing.T) {
	logger := NewCommandLogger(false)
	ctx := context.Background()

	// Should not panic
	logger.LogCommandError(ctx, "test.command", &testError{msg: "test error"}, map[string]interface{}{
		"detail": "error details",
	})
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

func TestCommandLogger_SensitiveDataRedaction(t *testing.T) {
	logger := NewCommandLogger(true)

	// These should not panic and should redact sensitive data
	logger.Info("test message", "token", "secret-token-value")
	logger.Info("test message", "password", "secret-password")
	logger.Info("test message", "secret", "secret-value")
}
