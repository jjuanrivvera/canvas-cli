package commands

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
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

func TestValidateCourseID_InvalidID(t *testing.T) {
	// Test with an invalid course ID (0 or negative)
	apiCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow account detection call, but fail if course API is called
		if r.URL.Path == "/api/v1/accounts" {
			w.Write([]byte("[]"))
			return
		}
		if r.URL.Path == "/api/v1/courses/0" || r.URL.Path == "/api/v1/courses/-1" {
			apiCalled = true
			t.Error("should not make API call for invalid ID")
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	client, err := api.NewClient(api.ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 10,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test with 0
	_, err = validateCourseID(client, 0)
	if err == nil {
		t.Error("Expected error for course ID 0")
	}

	// Test with negative
	_, err = validateCourseID(client, -1)
	if err == nil {
		t.Error("Expected error for negative course ID")
	}

	if apiCalled {
		t.Error("API was called for invalid IDs")
	}
}

func TestValidateCourseID_ValidCourse(t *testing.T) {
	// Create a mock server that returns a valid course
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/courses/123" {
			course := api.Course{
				ID:   123,
				Name: "Test Course",
			}
			json.NewEncoder(w).Encode(course)
			return
		}
		// Return empty array for account detection
		if r.URL.Path == "/api/v1/accounts" {
			w.Write([]byte("[]"))
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	client, err := api.NewClient(api.ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 10,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test with valid course ID
	course, err := validateCourseID(client, 123)
	if err != nil {
		t.Errorf("Expected no error for valid course, got: %v", err)
	}
	if course == nil {
		t.Error("Expected course to be returned")
	}
	if course != nil && course.ID != 123 {
		t.Errorf("Expected course ID 123, got %d", course.ID)
	}
}

func TestValidateCourseID_NotFound(t *testing.T) {
	// Create a mock server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/courses/999" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"errors":[{"message":"The specified resource does not exist"}]}`))
			return
		}
		// Return empty array for account detection
		if r.URL.Path == "/api/v1/accounts" {
			w.Write([]byte("[]"))
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	client, err := api.NewClient(api.ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 10,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test with non-existent course ID
	_, err = validateCourseID(client, 999)
	if err == nil {
		t.Error("Expected error for non-existent course")
	}
	// Error message should be user-friendly
	if err != nil {
		errMsg := err.Error()
		if !containsAny(errMsg, "not found", "999") {
			t.Errorf("Expected user-friendly error message, got: %s", errMsg)
		}
	}
}

func TestValidateCourseID_Unauthorized(t *testing.T) {
	// Create a mock server that returns 401
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/courses/456" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"errors":[{"message":"Invalid access token"}]}`))
			return
		}
		// Return empty array for account detection
		if r.URL.Path == "/api/v1/accounts" {
			w.Write([]byte("[]"))
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	client, err := api.NewClient(api.ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 10,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test with unauthorized access
	_, err = validateCourseID(client, 456)
	if err == nil {
		t.Error("Expected error for unauthorized access")
	}
}

// containsAny checks if the string contains any of the given substrings
func containsAny(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if len(s) >= len(substr) {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
		}
	}
	return false
}

func TestCreateCache(t *testing.T) {
	// Test that createCache returns a valid cache
	c := createCache()
	if c == nil {
		t.Fatal("Expected createCache to return non-nil cache")
	}

	// Test basic cache operations
	key := "test-key"
	value := []byte(`"test-value"`)

	c.Set(key, value)

	if !c.Has(key) {
		t.Error("Expected key to exist after Set")
	}

	retrieved := c.Get(key)
	if retrieved == nil {
		t.Error("Expected to retrieve value")
	}

	if string(retrieved) != string(value) {
		t.Errorf("Expected '%s', got '%s'", value, retrieved)
	}
}

func TestGetAPIClient_WithCache(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("CANVAS_URL", "https://test.instructure.com")
	os.Setenv("CANVAS_TOKEN", "test-token-123")
	defer os.Unsetenv("CANVAS_URL")
	defer os.Unsetenv("CANVAS_TOKEN")

	// Reset noCache flag
	originalNoCache := noCache
	noCache = false
	defer func() { noCache = originalNoCache }()

	// Create API client with cache enabled
	client, err := getAPIClient()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if client == nil {
		t.Fatal("Expected client to be created")
	}

	// Verify cache is enabled
	if !client.IsCacheEnabled() {
		t.Error("Expected cache to be enabled when noCache=false")
	}
}

func TestGetAPIClient_NoCache(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("CANVAS_URL", "https://test.instructure.com")
	os.Setenv("CANVAS_TOKEN", "test-token-123")
	defer os.Unsetenv("CANVAS_URL")
	defer os.Unsetenv("CANVAS_TOKEN")

	// Set noCache flag
	originalNoCache := noCache
	noCache = true
	defer func() { noCache = originalNoCache }()

	// Create API client with cache disabled
	client, err := getAPIClient()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if client == nil {
		t.Fatal("Expected client to be created")
	}

	// Verify cache is disabled
	if client.IsCacheEnabled() {
		t.Error("Expected cache to be disabled when noCache=true")
	}
}

func TestGetUserAgent_DevVersion(t *testing.T) {
	// Save original version
	originalVersion := version
	defer func() { version = originalVersion }()

	// Test with dev version
	version = "dev"
	ua := getUserAgent()
	if ua != "canvas-cli/dev" {
		t.Errorf("Expected 'canvas-cli/dev', got '%s'", ua)
	}

	// Test with empty version
	version = ""
	ua = getUserAgent()
	if ua != "canvas-cli/dev" {
		t.Errorf("Expected 'canvas-cli/dev' for empty version, got '%s'", ua)
	}
}

func TestGetUserAgent_ReleaseVersion(t *testing.T) {
	// Save original version
	originalVersion := version
	defer func() { version = originalVersion }()

	// Test with release version
	version = "v1.5.0"
	ua := getUserAgent()
	if ua != "canvas-cli/v1.5.0" {
		t.Errorf("Expected 'canvas-cli/v1.5.0', got '%s'", ua)
	}

	// Test with another version format
	version = "1.6.0-beta"
	ua = getUserAgent()
	if ua != "canvas-cli/1.6.0-beta" {
		t.Errorf("Expected 'canvas-cli/1.6.0-beta', got '%s'", ua)
	}
}

func TestGetAPIClient_UserAgentSet(t *testing.T) {
	// This test verifies that the User-Agent is set when creating a client via env vars
	var receivedUserAgent string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedUserAgent = r.Header.Get("User-Agent")
		if r.URL.Path == "/api/v1/accounts" {
			w.Write([]byte("[]"))
			return
		}
		if r.URL.Path == "/api/v1/courses" {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("[]"))
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	// Set environment variables
	os.Setenv("CANVAS_URL", server.URL)
	os.Setenv("CANVAS_TOKEN", "test-token-123")
	defer os.Unsetenv("CANVAS_URL")
	defer os.Unsetenv("CANVAS_TOKEN")

	// Save and set version for test
	originalVersion := version
	version = "v1.5.0"
	defer func() { version = originalVersion }()

	// Create API client
	client, err := getAPIClient()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Make a request to capture User-Agent header
	ctx := context.Background()
	_, _ = client.Get(ctx, "/api/v1/courses")

	// Verify User-Agent was set correctly
	expectedUA := "canvas-cli/v1.5.0"
	if receivedUserAgent != expectedUA {
		t.Errorf("Expected User-Agent '%s', got '%s'", expectedUA, receivedUserAgent)
	}
}
