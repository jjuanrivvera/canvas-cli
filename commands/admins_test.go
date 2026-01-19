package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestAdminsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list account admins successfully",
			Args: []string{"--account-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/accounts/1/admins": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"role": "AccountAdmin",
						"user": {
							"id": 100,
							"name": "Admin User",
							"email": "admin@example.com"
						}
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Admin User") {
					t.Error("Expected 'Admin User' in output")
				}
			},
		},
		{
			Name: "list account admins - empty response",
			Args: []string{"--account-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/accounts/1/admins": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No admins found",
		},
		{
			Name:        "list account admins - missing account ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAdminsListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAdminsAddCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "add admin successfully",
			Args: []string{"--account-id", "1", "--user-id", "100", "--role", "AccountAdmin"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/admins": cmdtest.NewMockResponse(`{
					"id": 10,
					"role": "AccountAdmin",
					"user": {
						"id": 100,
						"name": "New Admin"
					}
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Admin added successfully") {
					t.Error("Expected 'Admin added successfully' in output")
				}
			},
		},
		{
			Name:        "create admin - missing user ID",
			Args:        []string{"--account-id", "1", "--role", "AccountAdmin"},
			ExpectError: true,
		},
		{
			Name:        "create admin - missing role",
			Args:        []string{"--account-id", "1", "--user-id", "100"},
			ExpectError: true,
		},
	}
	//
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAdminsAddCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAdminsRemoveCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "remove admin successfully",
			Args: []string{"--account-id", "1", "--user-id", "100"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/admins/100": cmdtest.NewMockResponse(`{
					"id": 100
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "remove admin - missing user ID",
			Args:        []string{"--account-id", "1"},
			ExpectError: true,
		},
	}
	//
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAdminsRemoveCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAdminsRolesListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list account roles successfully",
			Args: []string{"--account-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/roles": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"label": "Account Admin",
						"base_role_type": "AccountAdmin",
						"workflow_state": "active"
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Account Admin") {
					t.Error("Expected 'Account Admin' in output")
				}
			},
		},
		{
			Name: "list account roles - empty response",
			Args: []string{"--account-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/roles": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No roles found",
		},
	}
	//
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newRolesListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
