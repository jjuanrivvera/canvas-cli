package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestExternalToolsLaunchCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "launch external tool successfully",
			Args: []string{"10", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/external_tools/sessionless_launch": cmdtest.NewMockResponse(`{
					"id": 10,
					"name": "My Tool",
					"url": "https://lti.example.com/launch?token=abc"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "lti.example.com") {
					t.Error("Expected launch URL in output")
				}
			},
		},
		{
			Name:        "launch - missing course ID",
			Args:        []string{"10"},
			ExpectError: true,
		},
		{
			Name:        "launch - missing tool ID",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
		{
			Name: "launch - API error",
			Args: []string{"10", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/external_tools/sessionless_launch": cmdtest.NewErrorResponse(404, "tool not found"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newExtToolsLaunchCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestExternalToolsCreateCmd_AccountContext(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create external tool in account context",
		Args: []string{"--account-id", "1", "--name", "Account Tool", "--url", "https://tool.example.com"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/external_tools": cmdtest.NewMockResponse(`{
				"id": 51,
				"name": "Account Tool",
				"url": "https://tool.example.com"
			}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Account Tool") {
				t.Error("Expected 'Account Tool' in output")
			}
		},
	}
	cmd := newExtToolsCreateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestExternalToolsCreateCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create external tool - API error",
		Args: []string{"--course-id", "1", "--name", "Bad Tool"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                courseMock,
			"/api/v1/courses/1/external_tools": cmdtest.NewErrorResponse(422, "invalid tool"),
		},
		ExpectError: true,
	}
	cmd := newExtToolsCreateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestExternalToolsUpdateCmd_AccountContext(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update external tool in account context",
		Args: []string{"10", "--account-id", "1", "--name", "Updated Account Tool"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/external_tools/10": cmdtest.NewMockResponse(`{
				"id": 10,
				"name": "Updated Account Tool"
			}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Updated Account Tool") {
				t.Error("Expected 'Updated Account Tool' in output")
			}
		},
	}
	cmd := newExtToolsUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestExternalToolsUpdateCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update external tool - API error",
		Args: []string{"10", "--course-id", "1", "--name", "Broken"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                   courseMock,
			"/api/v1/courses/1/external_tools/10": cmdtest.NewErrorResponse(404, "tool not found"),
		},
		ExpectError: true,
	}
	cmd := newExtToolsUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestExternalToolsDeleteCmd_AccountContext(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete external tool in account context",
		Args: []string{"10", "--account-id", "1", "--force"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/external_tools/10": cmdtest.NewMockResponse(`{
				"id": 10,
				"name": "Old Account Tool"
			}`),
		},
		ExpectError: false,
	}
	cmd := newExtToolsDeleteCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestExternalToolsDeleteCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete external tool - API error",
		Args: []string{"10", "--course-id", "1", "--force"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                   courseMock,
			"/api/v1/courses/1/external_tools/10": cmdtest.NewErrorResponse(404, "tool not found"),
		},
		ExpectError: true,
	}
	cmd := newExtToolsDeleteCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestExternalToolsListCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list external tools - API error",
		Args: []string{"--course-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                courseMock,
			"/api/v1/courses/1/external_tools": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmd := newExtToolsListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}
