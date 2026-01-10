package commands

import (
	"os"
	"testing"
)

func TestGetAPIClient_EnvironmentVariables(t *testing.T) {
	// Set environment variables
	os.Setenv("CANVAS_URL", "https://test.instructure.com")
	os.Setenv("CANVAS_TOKEN", "test-token-123")
	defer os.Unsetenv("CANVAS_URL")
	defer os.Unsetenv("CANVAS_TOKEN")

	// Create API client
	client, err := getAPIClient()

	// Should create client from environment variables
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if client == nil {
		t.Fatal("Expected client to be created")
	}
}

func TestGetAPIClient_EnvironmentVariables_WithRPS(t *testing.T) {
	// Set environment variables including requests per second
	os.Setenv("CANVAS_URL", "https://test.instructure.com")
	os.Setenv("CANVAS_TOKEN", "test-token-123")
	os.Setenv("CANVAS_REQUESTS_PER_SEC", "10.5")
	defer os.Unsetenv("CANVAS_URL")
	defer os.Unsetenv("CANVAS_TOKEN")
	defer os.Unsetenv("CANVAS_REQUESTS_PER_SEC")

	// Create API client
	client, err := getAPIClient()

	// Should create client from environment variables
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if client == nil {
		t.Fatal("Expected client to be created")
	}
}

func TestGetAPIClient_EnvironmentVariables_Partial(t *testing.T) {
	// Set only URL, not token - should fall through to config-based auth
	os.Setenv("CANVAS_URL", "https://test.instructure.com")
	defer os.Unsetenv("CANVAS_URL")

	// Create API client - will either work (if config exists) or fail (if no config)
	// This test just verifies partial env vars don't cause panic
	client, _ := getAPIClient()

	// If we have environment variables set, this should work with config
	// The important thing is that CANVAS_URL alone doesn't bypass auth entirely
	if client != nil {
		t.Log("Client created from config (expected if config exists)")
	}
}
