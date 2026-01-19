// Package testing provides a test framework for command testing
package testing

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// CommandTestCase defines a single command test scenario
type CommandTestCase struct {
	Name           string
	Args           []string
	MockResponses  map[string]string // path -> JSON response
	MockStatus     map[string]int    // path -> HTTP status code (default: 200)
	ExpectError    bool
	ExpectOutput   string
	ValidateOutput func(t *testing.T, output string)
}

// RunCommandTest executes a command test with mocked API
func RunCommandTest(t *testing.T, cmd *cobra.Command, tc CommandTestCase) {
	t.Helper()

	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Find matching path
		path := strings.TrimPrefix(r.URL.Path, "/api/v1")

		// Check if we have a mock response
		for mockPath, response := range tc.MockResponses {
			if strings.Contains(path, mockPath) || path == mockPath {
				status := http.StatusOK
				if tc.MockStatus != nil {
					if s, ok := tc.MockStatus[mockPath]; ok {
						status = s
					}
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(status)
				_, _ = w.Write([]byte(response))
				return
			}
		}

		// Default: return empty array
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("[]"))
	}))
	defer server.Close()

	// Setup command with test args
	cmd.SetArgs(tc.Args)

	// Capture output
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)

	// Execute command
	err := cmd.ExecuteContext(context.Background())

	// Validate error expectation
	if tc.ExpectError && err == nil {
		t.Errorf("%s: expected error but got none", tc.Name)
		t.Logf("stdout: %s", stdout.String())
		t.Logf("stderr: %s", stderr.String())
	}

	if !tc.ExpectError && err != nil {
		t.Errorf("%s: unexpected error: %v", tc.Name, err)
		t.Logf("stdout: %s", stdout.String())
		t.Logf("stderr: %s", stderr.String())
	}

	// Validate output if provided
	output := stdout.String()
	if tc.ExpectOutput != "" && !strings.Contains(output, tc.ExpectOutput) {
		t.Errorf("%s: expected output to contain %q, got %q", tc.Name, tc.ExpectOutput, output)
	}

	// Custom validation
	if tc.ValidateOutput != nil {
		tc.ValidateOutput(t, output)
	}
}

// MockJSONResponse creates a JSON response string
func MockJSONResponse(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// MockErrorResponse creates an error response string
func MockErrorResponse(message string) string {
	return `{"errors":[{"message":"` + message + `"}]}`
}
