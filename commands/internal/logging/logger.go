// Package logging provides structured logging for commands
package logging

import (
	"context"
	"log/slog"
	"os"
	"time"
)

// CommandLogger wraps slog with command-specific functionality
type CommandLogger struct {
	logger *slog.Logger
}

// NewCommandLogger creates a logger for command execution
func NewCommandLogger(debug bool) *CommandLogger {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Don't log sensitive data
			if a.Key == "token" || a.Key == "password" || a.Key == "secret" {
				return slog.Attr{Key: a.Key, Value: slog.StringValue("***")}
			}
			return a
		},
	})

	return &CommandLogger{
		logger: slog.New(handler),
	}
}

// LogCommandStart logs command execution start
func (l *CommandLogger) LogCommandStart(ctx context.Context, cmd string, args map[string]interface{}) {
	l.logger.InfoContext(ctx, "Command started",
		"command", cmd,
		"args", args,
	)
}

// LogAPICall logs an API request
func (l *CommandLogger) LogAPICall(ctx context.Context, method, path string, statusCode int, duration time.Duration) {
	l.logger.InfoContext(ctx, "API call",
		"method", method,
		"path", path,
		"status_code", statusCode,
		"duration_ms", duration.Milliseconds(),
	)
}

// LogCommandComplete logs successful command completion
func (l *CommandLogger) LogCommandComplete(ctx context.Context, cmd string, recordsProcessed int) {
	l.logger.InfoContext(ctx, "Command completed",
		"command", cmd,
		"records_processed", recordsProcessed,
	)
}

// LogCommandError logs command errors with context
func (l *CommandLogger) LogCommandError(ctx context.Context, cmd string, err error, context map[string]interface{}) {
	l.logger.ErrorContext(ctx, "Command failed",
		"command", cmd,
		"error", err.Error(),
		"context", context,
	)
}

// Debug logs a debug message
func (l *CommandLogger) Debug(msg string, args ...interface{}) {
	l.logger.Debug(msg, args...)
}

// Info logs an info message
func (l *CommandLogger) Info(msg string, args ...interface{}) {
	l.logger.Info(msg, args...)
}

// Warn logs a warning message
func (l *CommandLogger) Warn(msg string, args ...interface{}) {
	l.logger.Warn(msg, args...)
}

// Error logs an error message
func (l *CommandLogger) Error(msg string, args ...interface{}) {
	l.logger.Error(msg, args...)
}
