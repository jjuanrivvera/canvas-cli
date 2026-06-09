package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"
)

// newTestLogger creates a CommandLogger that writes JSON to the supplied buffer
// so tests can inspect what was logged without capturing os.Stderr.
func newTestLogger(buf *bytes.Buffer, debug bool) *CommandLogger {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}
	handler := slog.NewJSONHandler(buf, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "token" || a.Key == "password" || a.Key == "secret" {
				return slog.Attr{Key: a.Key, Value: slog.StringValue("***")}
			}
			return a
		},
	})
	return &CommandLogger{logger: slog.New(handler)}
}

func TestCommandLogger_Debug(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf, true) // must be debug-level to emit Debug messages

	logger.Debug("debug message", "key", "value")

	output := buf.String()
	if output == "" {
		t.Fatal("expected log output but got nothing")
	}
	if !strings.Contains(output, "debug message") {
		t.Errorf("expected 'debug message' in output, got: %s", output)
	}
}

func TestCommandLogger_Debug_BelowLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf, false) // info level — Debug should be suppressed

	logger.Debug("should be suppressed", "key", "value")

	if buf.Len() != 0 {
		t.Errorf("expected no output at info level, got: %s", buf.String())
	}
}

func TestCommandLogger_Warn(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf, false)

	logger.Warn("something is off", "code", 42)

	output := buf.String()
	if output == "" {
		t.Fatal("expected log output but got nothing")
	}
	if !strings.Contains(output, "something is off") {
		t.Errorf("expected warning message in output, got: %s", output)
	}
}

func TestCommandLogger_Error(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf, false)

	logger.Error("something broke", "err", "timeout")

	output := buf.String()
	if output == "" {
		t.Fatal("expected log output but got nothing")
	}
	if !strings.Contains(output, "something broke") {
		t.Errorf("expected error message in output, got: %s", output)
	}
}

func TestCommandLogger_Info(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf, false)

	logger.Info("all good", "count", 7)

	output := buf.String()
	if !strings.Contains(output, "all good") {
		t.Errorf("expected 'all good' in output, got: %s", output)
	}
}

func TestCommandLogger_SensitiveKeys_Redacted(t *testing.T) {
	sensitive := []string{"token", "password", "secret"}

	for _, key := range sensitive {
		t.Run(key, func(t *testing.T) {
			var buf bytes.Buffer
			logger := newTestLogger(&buf, false)

			logger.Info("check redaction", key, "super-secret-value")

			// Parse JSON so we can inspect the structured field directly.
			var record map[string]interface{}
			if err := json.Unmarshal(buf.Bytes(), &record); err != nil {
				t.Fatalf("failed to parse log JSON: %v (output: %s)", err, buf.String())
			}
			if val, ok := record[key]; !ok {
				t.Errorf("expected key %q in log record", key)
			} else if val != "***" {
				t.Errorf("expected key %q to be redacted, got: %v", key, val)
			}
		})
	}
}

func TestCommandLogger_LogCommandStart_Debug(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf, true) // debug so the DebugContext line fires
	ctx := context.Background()

	logger.LogCommandStart(ctx, "courses.list", map[string]interface{}{
		"course_id": int64(42),
		"include":   []string{"assignments"},
	})

	output := buf.String()
	if !strings.Contains(output, "Command started") {
		t.Errorf("expected 'Command started' in output, got: %s", output)
	}
	if !strings.Contains(output, "courses.list") {
		t.Errorf("expected command name in output, got: %s", output)
	}
}

func TestCommandLogger_LogCommandComplete_Debug(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf, true)
	ctx := context.Background()

	logger.LogCommandComplete(ctx, "courses.list", 100)

	output := buf.String()
	if !strings.Contains(output, "Command completed") {
		t.Errorf("expected 'Command completed' in output, got: %s", output)
	}
}

func TestCommandLogger_LogAPICall_Debug(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf, true)
	ctx := context.Background()

	logger.LogAPICall(ctx, "GET", "/api/v1/courses", 200, 50*time.Millisecond)

	output := buf.String()
	if !strings.Contains(output, "API call") {
		t.Errorf("expected 'API call' in output, got: %s", output)
	}
}

func TestCommandLogger_LogCommandError_WithContext(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf, false)
	ctx := context.Background()

	logger.LogCommandError(ctx, "courses.list", errors.New("connection refused"), map[string]interface{}{
		"retry": 3,
	})

	output := buf.String()
	if !strings.Contains(output, "Command failed") {
		t.Errorf("expected 'Command failed' in output, got: %s", output)
	}
	if !strings.Contains(output, "connection refused") {
		t.Errorf("expected error message in output, got: %s", output)
	}
}
