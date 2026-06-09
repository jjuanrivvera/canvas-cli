package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestExternalToolsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list external tools for course successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/external_tools": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"name": "Google Drive",
						"description": "LTI tool for Google Drive",
						"consumer_key": "key123"
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Google Drive") {
					t.Error("Expected 'Google Drive' in output")
				}
			},
		},
		{
			Name: "list external tools - empty response",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                courseMock,
				"/api/v1/courses/1/external_tools": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No external tools found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newExtToolsListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestExternalToolsGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get external tool successfully",
			Args: []string{"10", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/external_tools/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"name": "Turnitin",
					"description": "Plagiarism detection tool",
					"url": "https://lti.turnitin.com"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Turnitin") {
					t.Error("Expected 'Turnitin' in output")
				}
			},
		},
		{
			Name:        "get external tool - missing tool ID",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newExtToolsGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestExternalToolsDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete external tool with confirmation",
			Args: []string{"10", "--course-id", "1", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/external_tools/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"name": "Old Tool"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete external tool - missing tool ID",
			Args:        []string{"--course-id", "1", "--force"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newExtToolsDeleteCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestExternalToolsCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create external tool successfully",
			Args: []string{"--course-id", "1", "--name", "My LTI Tool", "--url", "https://tool.example.com"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/external_tools": cmdtest.NewMockResponse(`{
					"id": 50,
					"name": "My LTI Tool",
					"url": "https://tool.example.com"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "My LTI Tool") {
					t.Error("Expected 'My LTI Tool' in output")
				}
			},
		},
		{
			Name:        "create external tool - missing name",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newExtToolsCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestExternalToolsUpdateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "update external tool successfully",
			Args: []string{"10", "--course-id", "1", "--name", "Updated Tool"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/external_tools/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"name": "Updated Tool",
					"url": "https://tool.example.com"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Updated Tool") {
					t.Error("Expected 'Updated Tool' in output")
				}
			},
		},
		{
			Name:        "update external tool - missing tool ID",
			Args:        []string{"--course-id", "1", "--name", "Updated"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newExtToolsUpdateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestExternalToolsListCmd_AccountContext(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list external tools by account",
		Args: []string{"--account-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/external_tools": cmdtest.NewMockResponse(`[
				{
					"id": 3,
					"name": "Account Tool",
					"description": "Account-level LTI tool"
				}
			]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Account Tool") {
				t.Error("Expected 'Account Tool' in output")
			}
		},
	}

	cmd := newExtToolsListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestExternalToolsGetCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get external tool - API error",
		Args: []string{"99", "--course-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                courseMock,
			"/api/v1/courses/1/external_tools": cmdtest.NewErrorResponse(404, "tool not found"),
		},
		ExpectError: true,
	}

	cmd := newExtToolsGetCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}
