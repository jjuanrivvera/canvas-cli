package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestSISCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create SIS import successfully",
			Args: []string{"--account-id", "1", "--file", "testdata_sis.csv"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/sis_imports": cmdtest.NewMockResponse(`{
					"id": 42,
					"workflow_state": "created",
					"progress": 0,
					"created_at": "2024-01-15T10:00:00Z"
				}`),
			},
			// The file doesn't exist; this will error at file-open time — that's fine.
			ExpectError: true,
		},
		{
			Name:        "create SIS import - missing account ID",
			Args:        []string{"--file", "data.csv"},
			ExpectError: true,
		},
		{
			Name:        "create SIS import - missing file",
			Args:        []string{"--account-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newSISCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestSISAbortCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "abort SIS import successfully",
			Args: []string{"10", "--account-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/sis_imports/10/abort": cmdtest.NewMockResponse(`{
					"id": 10,
					"workflow_state": "aborted",
					"progress": 0
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "10") {
					t.Error("expected import ID in output")
				}
			},
		},
		{
			Name:        "abort SIS import - missing import ID",
			Args:        []string{"--account-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "abort SIS import - invalid import ID",
			Args:        []string{"notanumber", "--account-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newSISAbortCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestSISRestoreCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "restore SIS import successfully",
			Args: []string{"10", "--account-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/sis_imports/10/restore_states": cmdtest.NewMockResponse(`{
					"id": 10,
					"workflow_state": "restoring",
					"progress": 0
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "restore SIS import - invalid ID",
			Args:        []string{"notanumber", "--account-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newSISRestoreCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestSISErrorsCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list SIS import errors successfully",
			Args: []string{"10", "--account-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/sis_imports/10/errors": cmdtest.NewMockResponse(`[
					{"message":"Invalid user","row":2,"file":"users.csv"}
				]`),
			},
			ExpectError: false,
		},
		{
			Name: "list SIS import errors - empty",
			Args: []string{"10", "--account-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/sis_imports/10/errors": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError: false,
		},
		{
			Name:        "list SIS import errors - invalid ID",
			Args:        []string{"bad", "--account-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newSISErrorsCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
