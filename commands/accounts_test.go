package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestAccountsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list accounts successfully",
			Args: []string{"list"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"name": "ACUE",
						"workflow_state": "active",
						"parent_account_id": null,
						"root_account_id": null,
						"uuid": "test-uuid"
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "ACUE") {
					t.Error("Expected 'ACUE' in output")
				}
				if !strings.Contains(output, "active") {
					t.Error("Expected 'active' in output")
				}
			},
		},
		{
			Name: "list accounts - empty response",
			Args: []string{"list"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No accounts found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAccountsListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAccountsGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get account successfully",
			Args: []string{"1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1": cmdtest.NewMockResponse(`{
					"id": 1,
					"name": "ACUE",
					"workflow_state": "active",
					"parent_account_id": null,
					"root_account_id": null,
					"uuid": "test-uuid",
					"default_time_zone": "America/New_York"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "ACUE") {
					t.Error("Expected 'ACUE' in output")
				}
				if !strings.Contains(output, "America/New_York") {
					t.Error("Expected 'America/New_York' in output")
				}
			},
		},
		{
			Name: "get account - not found",
			Args: []string{"999"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/999": cmdtest.NewErrorResponse(404, "not found"),
			},
			ExpectError: true,
		},
		{
			Name:        "get account - invalid ID",
			Args:        []string{"invalid"},
			ExpectError: true,
		},
		{
			Name:        "get account - missing ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAccountsGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAccountsSubAccountsCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list sub-accounts successfully",
			Args: []string{"1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/sub_accounts": cmdtest.NewMockResponse(`[
					{
						"id": 2,
						"name": "Sub Account 1",
						"workflow_state": "active",
						"parent_account_id": 1,
						"root_account_id": 1
					},
					{
						"id": 3,
						"name": "Sub Account 2",
						"workflow_state": "active",
						"parent_account_id": 1,
						"root_account_id": 1
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Sub Account 1") {
					t.Error("Expected 'Sub Account 1' in output")
				}
				if !strings.Contains(output, "Sub Account 2") {
					t.Error("Expected 'Sub Account 2' in output")
				}
			},
		},
		{
			Name: "list sub-accounts - recursive",
			Args: []string{"1", "--recursive"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/sub_accounts": cmdtest.NewMockResponse(`[
					{
						"id": 2,
						"name": "Sub Account 1",
						"workflow_state": "active",
						"parent_account_id": 1,
						"root_account_id": 1
					}
				]`),
			},
			ExpectError:  false,
			ExpectOutput: "Sub Account 1",
		},
		{
			Name: "list sub-accounts - empty response",
			Args: []string{"1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/sub_accounts": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No sub-accounts found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAccountsSubAccountsCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
