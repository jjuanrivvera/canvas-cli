package commands

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// captureStdout captures os.Stdout during fn execution.
func captureStdout(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

// TestConfirmDelete_Force verifies that force=true skips the prompt and returns true.
func TestConfirmDelete_Force(t *testing.T) {
	confirmed, err := confirmDelete("assignment", 42, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !confirmed {
		t.Error("expected confirmed=true when force=true")
	}
}

// TestConfirmDelete_DryRun verifies dry-run mode outputs a message and returns false.
func TestConfirmDelete_DryRun(t *testing.T) {
	// Enable dry-run mode for this test
	old := dryRun
	dryRun = true
	defer func() { dryRun = old }()

	out := captureStdout(func() {
		confirmed, err := confirmDelete("course", 99, false)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if confirmed {
			t.Error("expected confirmed=false in dry-run mode")
		}
	})

	if !strings.Contains(out, "DRY RUN") {
		t.Errorf("expected 'DRY RUN' in output, got: %q", out)
	}
	if !strings.Contains(out, "99") {
		t.Errorf("expected resource ID in output, got: %q", out)
	}
}

// TestConfirmDeleteWithDetails_Force verifies force bypasses prompt and returns true.
func TestConfirmDeleteWithDetails_Force(t *testing.T) {
	details := map[string]interface{}{"name": "Test Course"}
	confirmed, err := confirmDeleteWithDetails("course", 1, details, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !confirmed {
		t.Error("expected confirmed=true when force=true")
	}
}

// TestConfirmDeleteWithDetails_DryRun verifies dry-run with details prints full preview.
func TestConfirmDeleteWithDetails_DryRun(t *testing.T) {
	old := dryRun
	dryRun = true
	defer func() { dryRun = old }()

	details := map[string]interface{}{"name": "My Course", "id": 7}

	out := captureStdout(func() {
		confirmed, err := confirmDeleteWithDetails("course", 7, details, false)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if confirmed {
			t.Error("expected confirmed=false in dry-run mode")
		}
	})

	if !strings.Contains(out, "DRY RUN") {
		t.Errorf("expected 'DRY RUN' in output, got: %q", out)
	}
}

// TestConfirmUpdateDryRun_NotDryRun verifies that when dryRun=false the function returns false.
func TestConfirmUpdateDryRun_NotDryRun(t *testing.T) {
	old := dryRun
	dryRun = false
	defer func() { dryRun = old }()

	changes := map[string]interface{}{"title": "new title"}
	result := confirmUpdateDryRun("assignment", 5, changes)
	if result {
		t.Error("expected false when dryRun=false")
	}
}

// TestConfirmUpdateDryRun_DryRun verifies that dry-run mode prints changes.
func TestConfirmUpdateDryRun_DryRun(t *testing.T) {
	old := dryRun
	dryRun = true
	defer func() { dryRun = old }()

	changes := map[string]interface{}{"title": "Updated Title"}

	out := captureStdout(func() {
		result := confirmUpdateDryRun("assignment", 10, changes)
		if !result {
			t.Error("expected true when dryRun=true")
		}
	})

	if !strings.Contains(out, "DRY RUN") {
		t.Errorf("expected 'DRY RUN' in output, got: %q", out)
	}
}
