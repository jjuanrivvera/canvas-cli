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
