package testing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// CommandTestCase defines a single command test scenario
type CommandTestCase struct {
	Name           string
	Args           []string
	MockResponses  map[string]MockResponse
	ExpectError    bool
	ExpectOutput   string
	ValidateOutput func(t *testing.T, output string)
	SetupClient    func(client *api.Client)
}

// MockResponse defines an HTTP mock response
type MockResponse struct {
	StatusCode int
	Body       string
	Headers    map[string]string
}

// RunCommandTest executes a command test with mocked API
func RunCommandTest(t *testing.T, cmd *cobra.Command, tc CommandTestCase) {
	t.Helper()

	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Find matching mock response - check longer patterns first for more specific matches
		var mockResp MockResponse
		found := false
		var patterns []string
		for pattern := range tc.MockResponses {
			patterns = append(patterns, pattern)
		}
		// Sort by length descending so longer/more specific patterns match first
		sort.Slice(patterns, func(i, j int) bool {
			return len(patterns[i]) > len(patterns[j])
		})

		for _, pattern := range patterns {
			if strings.Contains(path, pattern) || path == pattern {
				mockResp = tc.MockResponses[pattern]
				found = true
				break
			}
		}

		if !found {
			t.Logf("No mock response for path: %s", path)
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"errors": []map[string]string{
					{"message": "not found"},
				},
			})
			return
		}

		// Set headers
		for key, value := range mockResp.Headers {
			w.Header().Set(key, value)
		}

		// Set status code
		if mockResp.StatusCode != 0 {
			w.WriteHeader(mockResp.StatusCode)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		// Write body
		_, _ = io.WriteString(w, mockResp.Body)
	}))
	defer server.Close()

	// Create test API client pointing to mock server
	client, err := api.NewClient(api.ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 100, // High rate for tests
		UserAgent:      "canvas-cli-test",
	})
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}

	// Setup client if callback provided
	if tc.SetupClient != nil {
		tc.SetupClient(client)
	}

	// Set environment variables to use mock server
	// This allows getAPIClient() to return our test client
	t.Setenv("CANVAS_URL", server.URL)
	t.Setenv("CANVAS_TOKEN", "test-token")

	// Setup command with test args
	cmd.SetArgs(tc.Args)

	// Capture output - Need to redirect os.Stdout since formatOutput writes directly to it
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()

	os.Stdout = wOut
	os.Stderr = wErr

	// Capture cobra's output as well
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)

	// Execute command
	err = cmd.ExecuteContext(context.Background())

	// Restore stdout/stderr and close pipes
	wOut.Close()
	wErr.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	// Read captured output
	outputBytes, _ := io.ReadAll(rOut)
	stderrBytes, _ := io.ReadAll(rErr)

	output := string(outputBytes) + stdout.String()
	stderrOutput := string(stderrBytes) + stderr.String()

	// Validate results
	if tc.ExpectError && err == nil {
		t.Errorf("%s: expected error but got none", tc.Name)
	}

	if !tc.ExpectError && err != nil {
		t.Errorf("%s: unexpected error: %v", tc.Name, err)
		t.Logf("stderr: %s", stderrOutput)
	}

	// Validate output
	if tc.ExpectOutput != "" && !strings.Contains(output, tc.ExpectOutput) {
		t.Errorf("%s: expected output to contain %q, got: %s", tc.Name, tc.ExpectOutput, output)
	}

	if tc.ValidateOutput != nil {
		tc.ValidateOutput(t, output)
	}
}

// NewMockResponse creates a simple 200 OK JSON response
func NewMockResponse(body string) MockResponse {
	return MockResponse{
		StatusCode: http.StatusOK,
		Body:       body,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

// NewErrorResponse creates an error response
func NewErrorResponse(statusCode int, message string) MockResponse {
	body := fmt.Sprintf(`{"errors":[{"message":"%s"}]}`, message)
	return MockResponse{
		StatusCode: statusCode,
		Body:       body,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

// NewPaginatedResponse creates a paginated response with Link header
func NewPaginatedResponse(body string, nextPage string) MockResponse {
	resp := NewMockResponse(body)
	if nextPage != "" {
		resp.Headers["Link"] = fmt.Sprintf(`<%s>; rel="next"`, nextPage)
	}
	return resp
}
