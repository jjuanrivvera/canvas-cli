package commands

import (
	"strings"
	"testing"

	"github.com/jjuanrivvera/canvas-cli/internal/config"
)

// setupContextTestHome isolates config for context tests.
func setupContextTestHome(t *testing.T) {
	t.Helper()
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	config.ResetCache()
	t.Cleanup(config.ResetCache)
}

func TestContextSetCmd_Course(t *testing.T) {
	setupContextTestHome(t)

	out := captureRunOutput(func() {
		cmd := newContextSetCmd()
		cmd.SetArgs([]string{"course", "123"})
		if err := cmd.Execute(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "123") {
		t.Errorf("expected course ID in output, got: %q", out)
	}
}

func TestContextSetCmd_Assignment(t *testing.T) {
	setupContextTestHome(t)

	out := captureRunOutput(func() {
		cmd := newContextSetCmd()
		cmd.SetArgs([]string{"assignment", "456"})
		if err := cmd.Execute(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "456") {
		t.Errorf("expected assignment ID in output, got: %q", out)
	}
}

func TestContextSetCmd_User(t *testing.T) {
	setupContextTestHome(t)

	out := captureRunOutput(func() {
		cmd := newContextSetCmd()
		cmd.SetArgs([]string{"user", "789"})
		if err := cmd.Execute(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "789") {
		t.Errorf("expected user ID in output, got: %q", out)
	}
}

func TestContextSetCmd_Account(t *testing.T) {
	setupContextTestHome(t)

	out := captureRunOutput(func() {
		cmd := newContextSetCmd()
		cmd.SetArgs([]string{"account", "1"})
		if err := cmd.Execute(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "1") {
		t.Errorf("expected account ID in output, got: %q", out)
	}
}

func TestContextSetCmd_InvalidID(t *testing.T) {
	setupContextTestHome(t)

	cmd := newContextSetCmd()
	cmd.SetArgs([]string{"course", "not-a-number"})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for non-numeric ID")
	}
}

func TestContextSetCmd_InvalidType(t *testing.T) {
	setupContextTestHome(t)

	// Invalid type is silently accepted (no switch match) but stored
	cmd := newContextSetCmd()
	cmd.SetArgs([]string{"invalidtype", "99"})
	// Should not error at the command level; type validation is lenient
	_ = cmd.Execute()
}

func TestContextShowCmd_Empty(t *testing.T) {
	setupContextTestHome(t)

	out := captureRunOutput(func() {
		cmd := newContextShowCmd()
		cmd.SetArgs([]string{})
		if err := cmd.Execute(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "No context") {
		t.Errorf("expected 'No context' in output, got: %q", out)
	}
}

func TestContextShowCmd_WithValues(t *testing.T) {
	setupContextTestHome(t)

	// Set a course first
	setCmd := newContextSetCmd()
	setCmd.SetArgs([]string{"course", "321"})
	if err := setCmd.Execute(); err != nil {
		t.Fatalf("set failed: %v", err)
	}
	config.ResetCache()

	out := captureRunOutput(func() {
		cmd := newContextShowCmd()
		cmd.SetArgs([]string{})
		if err := cmd.Execute(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "321") {
		t.Errorf("expected course ID 321 in output, got: %q", out)
	}
}

func TestContextClearCmd_All(t *testing.T) {
	setupContextTestHome(t)

	// Set some context first
	setCmd := newContextSetCmd()
	setCmd.SetArgs([]string{"course", "100"})
	_ = setCmd.Execute()
	config.ResetCache()

	out := captureRunOutput(func() {
		cmd := newContextClearCmd()
		cmd.SetArgs([]string{})
		if err := cmd.Execute(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "cleared") {
		t.Errorf("expected 'cleared' in output, got: %q", out)
	}
}

func TestContextClearCmd_Specific(t *testing.T) {
	setupContextTestHome(t)

	// Set context
	setCmd := newContextSetCmd()
	setCmd.SetArgs([]string{"course", "200"})
	_ = setCmd.Execute()
	config.ResetCache()

	out := captureRunOutput(func() {
		cmd := newContextClearCmd()
		cmd.SetArgs([]string{"course"})
		if err := cmd.Execute(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "course") {
		t.Errorf("expected 'course' in output, got: %q", out)
	}
}

func TestGetContextCourseID_WithFlag(t *testing.T) {
	// Non-zero flag value should be returned as-is
	got := GetContextCourseID(42)
	if got != 42 {
		t.Errorf("expected 42, got %d", got)
	}
}

func TestGetContextCourseID_FromConfig(t *testing.T) {
	setupContextTestHome(t)

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config load failed: %v", err)
	}
	if err := cfg.SetContext(&config.Context{CourseID: 77}); err != nil {
		t.Fatalf("set context failed: %v", err)
	}
	config.ResetCache()

	got := GetContextCourseID(0) // zero means "use context"
	if got != 77 {
		t.Errorf("expected 77 from context, got %d", got)
	}
}

func TestGetContextAssignmentID_WithFlag(t *testing.T) {
	got := GetContextAssignmentID(55)
	if got != 55 {
		t.Errorf("expected 55, got %d", got)
	}
}

func TestGetContextAssignmentID_FromConfig(t *testing.T) {
	setupContextTestHome(t)

	cfg, _ := config.Load()
	_ = cfg.SetContext(&config.Context{AssignmentID: 88})
	config.ResetCache()

	got := GetContextAssignmentID(0)
	if got != 88 {
		t.Errorf("expected 88 from context, got %d", got)
	}
}

func TestGetContextUserID_WithFlag(t *testing.T) {
	got := GetContextUserID(11)
	if got != 11 {
		t.Errorf("expected 11, got %d", got)
	}
}

func TestGetContextUserID_FromConfig(t *testing.T) {
	setupContextTestHome(t)

	cfg, _ := config.Load()
	_ = cfg.SetContext(&config.Context{UserID: 22})
	config.ResetCache()

	got := GetContextUserID(0)
	if got != 22 {
		t.Errorf("expected 22 from context, got %d", got)
	}
}

func TestGetContextAccountID_WithFlag(t *testing.T) {
	got := GetContextAccountID(33)
	if got != 33 {
		t.Errorf("expected 33, got %d", got)
	}
}

func TestGetContextAccountID_FromConfig(t *testing.T) {
	setupContextTestHome(t)

	cfg, _ := config.Load()
	_ = cfg.SetContext(&config.Context{AccountID: 44})
	config.ResetCache()

	got := GetContextAccountID(0)
	if got != 44 {
		t.Errorf("expected 44 from context, got %d", got)
	}
}
