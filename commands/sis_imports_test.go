package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestSISListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list SIS imports successfully",
			Args: []string{"--account-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/sis_imports": cmdtest.NewMockResponse(`{
					"sis_imports": [
						{
							"id": 1,
							"workflow_state": "imported",
							"progress": 100,
							"created_at": "2024-01-15T10:00:00Z"
						}
					]
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "1") && !strings.Contains(output, "sis_import") {
					t.Error("Expected import data in output")
				}
			},
		},
		{
			Name:        "list SIS imports - missing account ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newSISListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestSISGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get SIS import successfully",
			Args: []string{"10", "--account-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/sis_imports/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"workflow_state": "imported",
					"progress": 100,
					"created_at": "2024-01-15T10:00:00Z"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "1") && !strings.Contains(output, "sis_import") {
					t.Error("Expected import data in output")
				}
			},
		},
		{
			Name:        "get SIS import - missing import ID",
			Args:        []string{"--account-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newSISGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
