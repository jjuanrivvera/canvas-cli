package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestConversationsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list conversations successfully",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/conversations": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"subject": "Question about assignment",
						"workflow_state": "unread",
						"last_message": "Can you help?"
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Question about assignment") {
					t.Error("Expected 'Question about assignment' in output")
				}
			},
		},
		{
			Name: "list conversations - empty response",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/conversations": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No conversations found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newConversationsListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestConversationsGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get conversation successfully",
			Args: []string{"10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/conversations/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"subject": "Grade inquiry",
					"workflow_state": "read",
					"messages": [
						{
							"id": 1,
							"body": "Message content"
						}
					]
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Grade inquiry") {
					t.Error("Expected 'Grade inquiry' in output")
				}
			},
		},
		{
			Name:        "get conversation - missing conversation ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newConversationsGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestConversationsDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete conversation successfully",
			Args: []string{"10", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/conversations/10": cmdtest.NewMockResponse(`{
					"id": 10
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete conversation - missing conversation ID",
			Args:        []string{"--force"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newConversationsDeleteCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
