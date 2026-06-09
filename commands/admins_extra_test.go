package commands

import (
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestRolesUpdateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "update role label successfully",
			Args: []string{"123", "--account-id", "1", "--label", "New Label"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/roles/123": cmdtest.NewMockResponse(`{
					"id": 123,
					"label": "New Label",
					"role": "custom_role",
					"workflow_state": "active"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "update role - missing account ID",
			Args:        []string{"123", "--label", "New Label"},
			ExpectError: true,
		},
		{
			Name:        "update role - missing role ID arg",
			Args:        []string{"--account-id", "1", "--label", "New Label"},
			ExpectError: true,
		},
		{
			Name: "update role - API error",
			Args: []string{"123", "--account-id", "1", "--label", "New Label"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/roles/123": cmdtest.NewErrorResponse(404, "role not found"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newRolesUpdateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestRolesDeactivateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "deactivate role successfully",
			Args: []string{"123", "--account-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/roles/123": cmdtest.NewMockResponse(`{
					"id": 123,
					"label": "Custom Role",
					"role": "custom_role",
					"workflow_state": "inactive"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "deactivate role - missing account ID",
			Args:        []string{"123"},
			ExpectError: true,
		},
		{
			Name: "deactivate role - API error",
			Args: []string{"123", "--account-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/roles/123": cmdtest.NewErrorResponse(404, "role not found"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newRolesDeactivateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestRolesActivateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "activate role successfully",
			Args: []string{"123", "--account-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/roles/123": cmdtest.NewMockResponse(`{
					"id": 123,
					"label": "Custom Role",
					"role": "custom_role",
					"workflow_state": "active"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "activate role - missing account ID",
			Args:        []string{"123"},
			ExpectError: true,
		},
		{
			Name: "activate role - API error",
			Args: []string{"456", "--account-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/roles/456": cmdtest.NewErrorResponse(404, "role not found"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newRolesActivateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAdminsRemoveCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "remove admin - API error",
		Args: []string{"99", "--account-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/admins/99": cmdtest.NewErrorResponse(404, "admin not found"),
		},
		ExpectError: true,
	}
	cmd := newAdminsRemoveCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestRolesGetCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get role - API error",
		Args: []string{"99", "--account-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/roles/99": cmdtest.NewErrorResponse(404, "role not found"),
		},
		ExpectError: true,
	}
	cmd := newRolesGetCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestRolesListCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list roles - API error",
		Args: []string{"--account-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/roles": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmd := newRolesListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}
