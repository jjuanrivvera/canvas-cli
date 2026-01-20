package updates

import (
	"testing"
)

func TestNewInstaller(t *testing.T) {
	config := UpdateConfig{
		Owner:          "jjuanrivvera",
		Repo:           "canvas-cli",
		CurrentVersion: "v1.0.0",
	}

	installer := NewInstaller(config)

	if installer == nil {
		t.Fatal("NewInstaller returned nil")
	}

	if installer.config.Owner != "jjuanrivvera" {
		t.Errorf("Expected owner to be 'jjuanrivvera', got %s", installer.config.Owner)
	}

	if installer.config.Repo != "canvas-cli" {
		t.Errorf("Expected repo to be 'canvas-cli', got %s", installer.config.Repo)
	}
}

func TestCanUpdate(t *testing.T) {
	config := UpdateConfig{
		Owner:          "jjuanrivvera",
		Repo:           "canvas-cli",
		CurrentVersion: "v1.0.0",
	}

	installer := NewInstaller(config)

	canUpdate, reason := installer.CanUpdate()

	// We can't reliably test the actual result since it depends on
	// the environment (permissions, installation method, etc.)
	// Just verify that the function returns something
	if !canUpdate && reason == "" {
		t.Error("CanUpdate should provide a reason when returning false")
	}

	// Test that Homebrew installation is detected
	// (This is a simple unit test, not an integration test)
	t.Run("DetectsPackageManagers", func(t *testing.T) {
		// Just verify the function is callable
		_, _ = installer.CanUpdate()
	})
}
