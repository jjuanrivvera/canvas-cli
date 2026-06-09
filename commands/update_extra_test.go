package commands

import (
	"strings"
	"testing"

	"github.com/jjuanrivvera/canvas-cli/internal/config"
)

// setupUpdateTestHome isolates config for update command tests.
func setupUpdateTestHome(t *testing.T) {
	t.Helper()
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	t.Setenv("USERPROFILE", tmpDir) // Windows: os.UserHomeDir() reads %USERPROFILE%
	config.ResetCache()
	t.Cleanup(config.ResetCache)
}

func TestUpdateEnableCmd_Success(t *testing.T) {
	setupUpdateTestHome(t)

	out := captureRunOutput(func() {
		cmd := newUpdateEnableCmd()
		cmd.SetArgs([]string{"--interval", "30"})
		if err := cmd.Execute(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "enabled") {
		t.Errorf("expected 'enabled' in output, got: %q", out)
	}
	if !strings.Contains(out, "30") {
		t.Errorf("expected interval '30' in output, got: %q", out)
	}
}

func TestUpdateDisableCmd_Success(t *testing.T) {
	setupUpdateTestHome(t)

	out := captureRunOutput(func() {
		cmd := newUpdateDisableCmd()
		cmd.SetArgs([]string{})
		if err := cmd.Execute(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "disabled") {
		t.Errorf("expected 'disabled' in output, got: %q", out)
	}
}

func TestUpdateStatusCmd_DefaultSettings(t *testing.T) {
	setupUpdateTestHome(t)

	out := captureRunOutput(func() {
		cmd := newUpdateStatusCmd()
		cmd.SetArgs([]string{})
		if err := cmd.Execute(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "Auto-Update Status") {
		t.Errorf("expected 'Auto-Update Status' in output, got: %q", out)
	}
}

func TestUpdateStatusCmd_AfterEnable(t *testing.T) {
	setupUpdateTestHome(t)

	// Enable first
	enableCmd := newUpdateEnableCmd()
	enableCmd.SetArgs([]string{"--interval", "120"})
	_ = enableCmd.Execute()
	config.ResetCache()

	out := captureRunOutput(func() {
		cmd := newUpdateStatusCmd()
		cmd.SetArgs([]string{})
		if err := cmd.Execute(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "enabled") {
		t.Errorf("expected 'enabled' status in output, got: %q", out)
	}
}

func TestUpdateStatusCmd_AfterDisable(t *testing.T) {
	setupUpdateTestHome(t)

	// Disable first
	disableCmd := newUpdateDisableCmd()
	disableCmd.SetArgs([]string{})
	_ = disableCmd.Execute()
	config.ResetCache()

	out := captureRunOutput(func() {
		cmd := newUpdateStatusCmd()
		cmd.SetArgs([]string{})
		if err := cmd.Execute(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, "disabled") {
		t.Errorf("expected 'disabled' status in output, got: %q", out)
	}
}

// TestUpdateCheckCmd_NoNetwork covers runUpdateCheck path when network is unavailable
// (the GitHub API will time out or fail). We simply verify the command returns an
// error (network unavailable in most CI environments) — or succeeds if network
// is up. The test validates that the command is reachable at all.
func TestUpdateCheckCmd_Reachable(t *testing.T) {
	setupUpdateTestHome(t)

	origVersion := version
	version = "1.0.0"
	defer func() { version = origVersion }()

	cmd := newUpdateCheckCmd()
	cmd.SetArgs([]string{})
	// We don't assert on the error here — a network call may succeed or fail.
	// The goal is to exercise the code path, not validate the network call.
	_ = cmd.Execute()
}
