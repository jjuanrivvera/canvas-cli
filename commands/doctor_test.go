package commands

import (
	"os"
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestDoctorCmd(t *testing.T) {
	// Skip in CI as doctor performs real system diagnostics that may fail in CI environments
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skip("Skipping doctor tests in CI environment")
	}

	tests := []cmdtest.CommandTestCase{
		{
			Name: "doctor check passes",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/users/self": cmdtest.NewMockResponse(`{
					"id": 1,
					"name": "Test User",
					"short_name": "Test",
					"sortable_name": "User, Test",
					"email": "test@example.com"
				}`),
				"/api/v1/courses": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Running diagnostics") {
					t.Error("Expected 'Running diagnostics' in output")
				}
				if !strings.Contains(output, "Environment") {
					t.Error("Expected 'Environment' check in output")
				}
			},
		},
		{
			Name: "doctor check with verbose",
			Args: []string{"--verbose"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/users/self": cmdtest.NewMockResponse(`{
					"id": 1,
					"name": "Test User",
					"short_name": "Test",
					"sortable_name": "User, Test",
					"email": "test@example.com"
				}`),
				"/api/v1/courses": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Description") {
					t.Error("Expected 'Description' in verbose output")
				}
			},
		},
		{
			Name: "doctor check with JSON output",
			Args: []string{"--json"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/users/self": cmdtest.NewMockResponse(`{
					"id": 1,
					"name": "Test User"
				}`),
				"/api/v1/courses": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, `"duration"`) {
					t.Error("Expected JSON output with duration field")
				}
				if !strings.Contains(output, `"healthy"`) {
					t.Error("Expected JSON output with healthy field")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newDoctorCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
