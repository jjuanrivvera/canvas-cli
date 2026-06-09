package commands

// auth_extra_test.go — tests for auth.go
//
// Not covered:
//   - runAuthLogin: starts an interactive OAuth browser flow. Intentionally
//     skipped because it requires a live browser callback and network access.
//   - runAuthLogout: prompts for stdin confirmation (fmt.Scanln). The
//     non-interactive path (confirm != "y") is not reachable without stdin
//     injection; skipped to avoid flaky tests.
//   - runAuthTokenRemove: also prompts via fmt.Scanln; skipped for same reason.

import (
	"strings"
	"testing"

	"github.com/jjuanrivvera/canvas-cli/internal/config"
)

// setupAuthTestHome creates an isolated HOME and resets config cache.
func setupAuthTestHome(t *testing.T) {
	t.Helper()
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	config.ResetCache()
	t.Cleanup(config.ResetCache)
}

// --- auth status ---

func TestAuthStatusCmd_NoInstances(t *testing.T) {
	setupAuthTestHome(t)

	out := captureRunOutput(func() {
		cmd := newAuthStatusCmd()
		cmd.SetArgs([]string{})
		if err := cmd.Execute(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "No Canvas instances") {
		t.Errorf("expected 'No Canvas instances' in output, got: %q", out)
	}
}

// validTestToken is a token string that passes Canvas token validation (≥20 chars).
const validTestToken = "7~aaabbbbccccddddeeee1234567890"

func TestAuthStatusCmd_WithTokenInstance(t *testing.T) {
	setupAuthTestHome(t)

	cfg, _ := config.Load()
	_ = cfg.AddInstance(&config.Instance{
		Name:  "demo",
		URL:   "https://demo.canvas.com",
		Token: validTestToken,
	})
	config.ResetCache()

	out := captureRunOutput(func() {
		cmd := newAuthStatusCmd()
		cmd.SetArgs([]string{})
		if err := cmd.Execute(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "demo") {
		t.Errorf("expected instance name 'demo' in output, got: %q", out)
	}
}

func TestAuthStatusCmd_SpecificInstance(t *testing.T) {
	setupAuthTestHome(t)

	cfg, _ := config.Load()
	_ = cfg.AddInstance(&config.Instance{
		Name:  "specific",
		URL:   "https://specific.canvas.com",
		Token: validTestToken,
	})
	config.ResetCache()

	out := captureRunOutput(func() {
		cmd := newAuthStatusCmd()
		cmd.SetArgs([]string{"specific"})
		if err := cmd.Execute(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "specific") {
		t.Errorf("expected 'specific' in output, got: %q", out)
	}
}

func TestAuthStatusCmd_InstanceNotFound(t *testing.T) {
	setupAuthTestHome(t)

	// Need at least one instance in config so the status command
	// proceeds past the "no instances" check.
	cfg, _ := config.Load()
	_ = cfg.AddInstance(&config.Instance{
		Name:  "existing",
		URL:   "https://existing.canvas.com",
		Token: validTestToken,
	})
	config.ResetCache()

	cmd := newAuthStatusCmd()
	cmd.SetArgs([]string{"nonexistent"})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for non-existent instance")
	}
}

// --- auth token set ---

func TestAuthTokenSetCmd_NewInstanceMissingURL(t *testing.T) {
	setupAuthTestHome(t)

	cmd := newAuthTokenSetCmd()
	cmd.SetArgs([]string{"newinstance", "--token", "mytoken"})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error when instance does not exist and --url is not provided")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' in error, got: %v", err)
	}
}

func TestAuthTokenSetCmd_UpdateExistingInstance(t *testing.T) {
	setupAuthTestHome(t)

	cfg, _ := config.Load()
	_ = cfg.AddInstance(&config.Instance{
		Name: "updateme",
		URL:  "https://updateme.canvas.com",
	})
	config.ResetCache()

	out := captureRunOutput(func() {
		cmd := newAuthTokenSetCmd()
		cmd.SetArgs([]string{"updateme", "--token", validTestToken})
		if err := cmd.Execute(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "updateme") {
		t.Errorf("expected 'updateme' in output, got: %q", out)
	}
}

func TestAuthTokenSetCmd_CreateNewInstanceWithURL(t *testing.T) {
	setupAuthTestHome(t)

	out := captureRunOutput(func() {
		cmd := newAuthTokenSetCmd()
		cmd.SetArgs([]string{"brand-new", "--url", "https://brandnew.canvas.com", "--token", validTestToken})
		if err := cmd.Execute(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "brand-new") {
		t.Errorf("expected 'brand-new' in output, got: %q", out)
	}
}

// --- auth login validation (non-interactive paths only) ---

func TestAuthLoginCmd_MissingURL(t *testing.T) {
	setupAuthTestHome(t)

	cmd := newAuthLoginCmd()
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error when no instance URL or name provided")
	}
}

func TestAuthLoginCmd_InstanceNotFound(t *testing.T) {
	setupAuthTestHome(t)

	cmd := newAuthLoginCmd()
	cmd.SetArgs([]string{"--instance", "nonexistent"})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for non-existent instance name")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' in error, got: %v", err)
	}
}

func TestAuthLoginCmd_InvalidOAuthMode(t *testing.T) {
	setupAuthTestHome(t)

	// Provide a valid URL so we get past the URL check, but invalid mode
	cfg, _ := config.Load()
	_ = cfg.AddInstance(&config.Instance{
		Name:         "myinst",
		URL:          "https://myinst.canvas.com",
		ClientID:     "clientid",
		ClientSecret: "secret",
	})
	config.ResetCache()

	cmd := newAuthLoginCmd()
	cmd.SetArgs([]string{"--instance", "myinst", "--mode", "invalidmode"})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for invalid OAuth mode")
	}
	if !strings.Contains(err.Error(), "invalid OAuth mode") {
		t.Errorf("unexpected error: %v", err)
	}
}

// --- auth logout validation ---

func TestAuthLogoutCmd_NoDefaultInstance(t *testing.T) {
	setupAuthTestHome(t)
	// No config — no default instance

	// logout prompts for stdin, but before that it checks for a default instance.
	// With no default, it returns error without reaching the prompt.
	cmd := newAuthLogoutCmd()
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error when no default instance configured")
	}
}
