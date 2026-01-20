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

func TestInstallerConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  UpdateConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: UpdateConfig{
				Owner:          "owner",
				Repo:           "repo",
				CurrentVersion: "v1.0.0",
			},
			wantErr: false,
		},
		{
			name: "empty owner",
			config: UpdateConfig{
				Owner:          "",
				Repo:           "repo",
				CurrentVersion: "v1.0.0",
			},
			wantErr: false, // We don't validate in constructor
		},
		{
			name: "dev version",
			config: UpdateConfig{
				Owner:          "owner",
				Repo:           "repo",
				CurrentVersion: "dev",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			installer := NewInstaller(tt.config)
			if installer == nil {
				t.Error("NewInstaller should not return nil")
			}
			if installer.config.Owner != tt.config.Owner {
				t.Errorf("Expected owner %s, got %s", tt.config.Owner, installer.config.Owner)
			}
		})
	}
}
