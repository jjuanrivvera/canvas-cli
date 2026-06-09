package commands

import (
	"strings"
	"testing"

	"github.com/jjuanrivvera/canvas-cli/internal/config"
)

func setupSyncTestHome(t *testing.T) {
	t.Helper()
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	t.Setenv("USERPROFILE", tmpDir) // Windows: os.UserHomeDir() reads %USERPROFILE%
	config.ResetCache()
	t.Cleanup(config.ResetCache)
}

func TestSyncAssignmentsCmd_InvalidSourceID(t *testing.T) {
	setupSyncTestHome(t)

	// Call RunE directly with invalid args
	err := runSyncAssignments(syncAssignmentsCmd, []string{"source-inst", "notanumber", "target-inst", "789"})
	if err == nil {
		t.Error("expected error for invalid source course ID")
		return
	}
	if !strings.Contains(err.Error(), "invalid source course ID") {
		t.Errorf("expected 'invalid source course ID' error, got: %v", err)
	}
}

func TestSyncAssignmentsCmd_InvalidTargetID(t *testing.T) {
	setupSyncTestHome(t)

	err := runSyncAssignments(syncAssignmentsCmd, []string{"source-inst", "123", "target-inst", "notanumber"})
	if err == nil {
		t.Error("expected error for invalid target course ID")
		return
	}
	if !strings.Contains(err.Error(), "invalid target course ID") {
		t.Errorf("expected 'invalid target course ID' error, got: %v", err)
	}
}

func TestSyncAssignmentsCmd_MissingInstance(t *testing.T) {
	setupSyncTestHome(t)

	err := runSyncAssignments(syncAssignmentsCmd, []string{"nonexistent-inst", "123", "another-nonexistent", "456"})
	if err == nil {
		t.Error("expected error when instance is not configured")
		return
	}
	if !strings.Contains(err.Error(), "failed to create source client") &&
		!strings.Contains(err.Error(), "instance not found") {
		t.Errorf("expected source client error, got: %v", err)
	}
}

func TestSyncCourseCmd_InvalidSourceID(t *testing.T) {
	setupSyncTestHome(t)

	err := runSyncCourse(syncCourseCmd, []string{"source-inst", "notanumber", "target-inst", "789"})
	if err == nil {
		t.Error("expected error for invalid source course ID")
		return
	}
	if !strings.Contains(err.Error(), "invalid source course ID") {
		t.Errorf("expected 'invalid source course ID' error, got: %v", err)
	}
}

func TestSyncCourseCmd_InvalidTargetID(t *testing.T) {
	setupSyncTestHome(t)

	err := runSyncCourse(syncCourseCmd, []string{"source-inst", "123", "target-inst", "notanumber"})
	if err == nil {
		t.Error("expected error for invalid target course ID")
		return
	}
	if !strings.Contains(err.Error(), "invalid target course ID") {
		t.Errorf("expected 'invalid target course ID' error, got: %v", err)
	}
}

func TestSyncCourseCmd_MissingInstance(t *testing.T) {
	setupSyncTestHome(t)

	err := runSyncCourse(syncCourseCmd, []string{"nonexistent-inst", "123", "another-nonexistent", "456"})
	if err == nil {
		t.Error("expected error when instance is not configured")
		return
	}
	if !strings.Contains(err.Error(), "failed to create source client") &&
		!strings.Contains(err.Error(), "instance not found") {
		t.Errorf("expected source client error, got: %v", err)
	}
}

func TestSyncCmdStructure(t *testing.T) {
	if syncCmd == nil {
		t.Fatal("expected non-nil sync command")
	}
	if syncCmd.Use != "sync" {
		t.Errorf("expected Use='sync', got: %q", syncCmd.Use)
	}

	subcommands := map[string]bool{
		"assignments": false,
		"course":      false,
	}
	for _, sub := range syncCmd.Commands() {
		subcommands[sub.Name()] = true
	}
	for name, found := range subcommands {
		if !found {
			t.Errorf("expected subcommand %q not found in sync", name)
		}
	}
}

func TestGetAPIClientForInstance_UnknownInstance(t *testing.T) {
	setupSyncTestHome(t)

	_, err := getAPIClientForInstance("does-not-exist")
	if err == nil {
		t.Error("expected error for unknown instance")
		return
	}
	// Could be config load error or instance not found
	if !strings.Contains(err.Error(), "instance not found") &&
		!strings.Contains(err.Error(), "failed to load") {
		t.Errorf("expected instance-not-found error, got: %v", err)
	}
}
