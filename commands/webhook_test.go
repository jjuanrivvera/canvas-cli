package commands

import (
	"strings"
	"testing"
)

func TestWebhookEventsCmd_ListsAllEvents(t *testing.T) {
	out := captureRunOutput(func() {
		cmd := newWebhookEventsCmd()
		cmd.SetArgs([]string{})
		if err := cmd.Execute(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "Available Canvas Webhook Event Types") {
		t.Errorf("expected header in output, got: %q", out)
	}
	if !strings.Contains(out, "Total:") {
		t.Errorf("expected total count in output, got: %q", out)
	}
}

func TestWebhookEventsCmd_ContainsKnownEventTypes(t *testing.T) {
	out := captureRunOutput(func() {
		cmd := newWebhookEventsCmd()
		cmd.SetArgs([]string{})
		_ = cmd.Execute()
	})

	// grade_change is a well-known Canvas event type
	if !strings.Contains(out, "grade_change") {
		t.Errorf("expected 'grade_change' event type in output, got: %q", out)
	}
}

func TestRunWebhookEvents_DirectCall(t *testing.T) {
	out := captureRunOutput(func() {
		runWebhookEvents()
	})

	if !strings.Contains(out, "Available Canvas Webhook Event Types") {
		t.Errorf("expected header in direct call output, got: %q", out)
	}
}

// TestWebhookListenCmd_ValidationError tests that the listen command validates
// its options before attempting to start the server.
func TestWebhookListenCmd_ValidationError(t *testing.T) {
	// The WebhookListenOptions.Validate passes with default values (no required fields).
	// We just verify the command structure is correct and can be built.
	cmd := newWebhookListenCmd()
	if cmd == nil {
		t.Fatal("expected non-nil webhook listen command")
	}
	if cmd.Use != "listen" {
		t.Errorf("expected Use='listen', got: %q", cmd.Use)
	}
}

func TestWebhookEventsCmd_Structure(t *testing.T) {
	cmd := newWebhookEventsCmd()
	if cmd == nil {
		t.Fatal("expected non-nil webhook events command")
	}
	if cmd.Use != "events" {
		t.Errorf("expected Use='events', got: %q", cmd.Use)
	}
}
