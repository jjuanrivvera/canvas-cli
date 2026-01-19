package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestOutcomesGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get outcome successfully",
			Args: []string{"10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/outcomes/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"title": "Critical Thinking",
					"description": "Students will demonstrate critical thinking skills"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Critical Thinking") {
					t.Error("Expected 'Critical Thinking' in output")
				}
			},
		},
		{
			Name:        "get outcome - missing outcome ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newOutcomesGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestOutcomesListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list outcomes successfully",
			Args: []string{"--account-id", "1", "--group-id", "456"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts": cmdtest.NewMockResponse(`[]`),
				"/api/v1/accounts/1/outcome_groups/456/outcomes": cmdtest.NewMockResponse(`[
					{
						"outcome": {
							"id": 1,
							"title": "Communication",
							"context_id": 1,
							"context_type": "Account"
						}
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Communication") {
					t.Error("Expected 'Communication' in output")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newOutcomesListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestOutcomesGroupsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list outcome groups successfully",
			Args: []string{"--account-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/outcome_groups": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"title": "Core Competencies",
						"description": "Essential skills"
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Core Competencies") {
					t.Error("Expected 'Core Competencies' in output")
				}
			},
		},
		{
			Name: "list outcome groups - empty response",
			Args: []string{"--account-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/outcome_groups": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No outcome groups found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newOutcomesGroupsListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestOutcomesGroupsGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get outcome group successfully",
			Args: []string{"5", "--account-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/outcome_groups/5": cmdtest.NewMockResponse(`{
					"id": 5,
					"title": "Program Outcomes",
					"description": "Program-level learning outcomes"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Program Outcomes") {
					t.Error("Expected 'Program Outcomes' in output")
				}
			},
		},
		{
			Name:        "get outcome group - missing group ID",
			Args:        []string{"--account-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newOutcomesGroupsGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
