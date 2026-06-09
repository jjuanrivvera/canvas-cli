package commands

import (
	"strings"
	"testing"

	"github.com/jjuanrivvera/canvas-cli/internal/config"
)

// setupAliasTestHome isolates config for alias tests.
func setupAliasTestHome(t *testing.T) {
	t.Helper()
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	t.Setenv("USERPROFILE", tmpDir) // Windows: os.UserHomeDir() reads %USERPROFILE%
	config.ResetCache()
	t.Cleanup(config.ResetCache)
}

func TestAliasSetCmd_Basic(t *testing.T) {
	setupAliasTestHome(t)

	out := captureRunOutput(func() {
		cmd := newAliasSetCmd()
		cmd.SetArgs([]string{"myalias", "courses list"})
		_ = cmd.Execute()
	})

	if !strings.Contains(out, "myalias") {
		t.Errorf("expected alias name in output, got: %q", out)
	}
}

func TestAliasSetCmd_ConflictsWithBuiltin(t *testing.T) {
	setupAliasTestHome(t)

	cmd := newAliasSetCmd()
	cmd.SetArgs([]string{"courses", "something else"})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error when alias conflicts with built-in command")
	}
}

func TestAliasListCmd_Empty(t *testing.T) {
	setupAliasTestHome(t)

	out := captureRunOutput(func() {
		cmd := newAliasListCmd()
		cmd.SetArgs([]string{})
		_ = cmd.Execute()
	})

	if !strings.Contains(out, "No aliases") {
		t.Errorf("expected 'No aliases' message, got: %q", out)
	}
}

func TestAliasListCmd_WithAliases(t *testing.T) {
	setupAliasTestHome(t)

	// Pre-populate an alias via config
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	if err := cfg.SetAlias("mytest", "courses list"); err != nil {
		t.Fatalf("failed to set alias: %v", err)
	}
	config.ResetCache()

	out := captureRunOutput(func() {
		cmd := newAliasListCmd()
		cmd.SetArgs([]string{})
		_ = cmd.Execute()
	})

	if !strings.Contains(out, "mytest") {
		t.Errorf("expected alias name 'mytest' in output, got: %q", out)
	}
}

func TestAliasDeleteCmd_Existing(t *testing.T) {
	setupAliasTestHome(t)

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	if err := cfg.SetAlias("todelete", "courses list"); err != nil {
		t.Fatalf("failed to set alias: %v", err)
	}
	config.ResetCache()

	out := captureRunOutput(func() {
		cmd := newAliasDeleteCmd()
		cmd.SetArgs([]string{"todelete"})
		_ = cmd.Execute()
	})

	if !strings.Contains(out, "todelete") {
		t.Errorf("expected alias name in output, got: %q", out)
	}
}

func TestAliasDeleteCmd_NonExistent(t *testing.T) {
	setupAliasTestHome(t)

	cmd := newAliasDeleteCmd()
	cmd.SetArgs([]string{"doesnotexist"})
	err := cmd.Execute()
	// Deleting non-existent alias should return an error
	if err == nil {
		t.Error("expected error when deleting non-existent alias")
	}
}

func TestAliasSetAndListRoundtrip(t *testing.T) {
	setupAliasTestHome(t)

	// Set an alias
	setCmd := newAliasSetCmd()
	setCmd.SetArgs([]string{"roundtrip", "assignments list --course-id 1"})
	if err := setCmd.Execute(); err != nil {
		t.Fatalf("set alias failed: %v", err)
	}
	config.ResetCache()

	// List should include it
	out := captureRunOutput(func() {
		listCmd := newAliasListCmd()
		listCmd.SetArgs([]string{})
		_ = listCmd.Execute()
	})
	if !strings.Contains(out, "roundtrip") {
		t.Errorf("expected 'roundtrip' in list output, got: %q", out)
	}
}
