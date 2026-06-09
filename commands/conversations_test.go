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

func TestConversationsCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create conversation successfully",
			Args: []string{"--recipients", "456,789", "--subject", "Hello", "--body", "Message content"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/conversations": cmdtest.NewMockResponse(`[
					{
						"id": 50,
						"subject": "Hello",
						"workflow_state": "unread"
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "50") {
					t.Error("Expected conversation ID 50 in output")
				}
			},
		},
		{
			Name:        "create conversation - missing recipients",
			Args:        []string{"--body", "Message"},
			ExpectError: true,
		},
		{
			Name:        "create conversation - missing body",
			Args:        []string{"--recipients", "456"},
			ExpectError: true,
		},
		{
			Name: "create conversation - API error",
			Args: []string{"--recipients", "456", "--body", "Hi"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/conversations": cmdtest.NewErrorResponse(500, "internal server error"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newConversationsCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestConversationsReplyCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "reply to conversation successfully",
			Args: []string{"10", "--body", "Thank you"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/conversations/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"subject": "Original Subject",
					"workflow_state": "read"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "reply - missing conversation ID",
			Args:        []string{"--body", "reply"},
			ExpectError: true,
		},
		{
			Name: "reply - API error",
			Args: []string{"10", "--body", "Hi"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/conversations/10/add_message": cmdtest.NewErrorResponse(404, "not found"),
				"/api/v1/conversations/10":             cmdtest.NewErrorResponse(404, "not found"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newConversationsReplyCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestConversationsArchiveCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "archive conversation successfully",
			Args: []string{"10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/conversations/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"workflow_state": "archived"
				}`),
			},
			ExpectError:  false,
			ExpectOutput: "10",
		},
		{
			Name:        "archive - missing conversation ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newConversationsArchiveCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestConversationsUnarchiveCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "unarchive conversation successfully",
			Args: []string{"10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/conversations/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"workflow_state": "read"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "unarchive - missing conversation ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newConversationsUnarchiveCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestConversationsStarCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "star conversation successfully",
			Args: []string{"10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/conversations/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"starred": true
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "star - missing conversation ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newConversationsStarCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestConversationsUnstarCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "unstar conversation successfully",
			Args: []string{"10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/conversations/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"starred": false
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "unstar - missing conversation ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newConversationsUnstarCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestConversationsMarkReadCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "mark conversation as read successfully",
			Args: []string{"10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/conversations/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"workflow_state": "read"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "mark-read - missing conversation ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newConversationsMarkReadCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestConversationsMarkAllReadCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "mark all conversations as read successfully",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/conversations/mark_all_as_read": cmdtest.NewMockResponse(`{}`),
			},
			ExpectError:  false,
			ExpectOutput: "All conversations marked as read",
		},
		{
			Name: "mark-all-read - API error",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/conversations/mark_all_as_read": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newConversationsMarkAllReadCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestConversationsUnreadCountCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get unread count successfully",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/conversations/unread_count": cmdtest.NewMockResponse(`{
					"unread_count": "5"
				}`),
			},
			ExpectError:  false,
			ExpectOutput: "5",
		},
		{
			Name: "get unread count - API error",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/conversations/unread_count": cmdtest.NewErrorResponse(401, "unauthorized"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newConversationsUnreadCountCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestConversationsAddRecipientsCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "add recipients successfully",
			Args: []string{"10", "--recipients", "456,789"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/conversations/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"subject": "Test"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "add-recipients - missing conversation ID",
			Args:        []string{"--recipients", "456"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newConversationsAddRecipientsCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestConversationsListCmd_WithScope(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list conversations with scope",
		Args: []string{"--scope", "unread"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/conversations": cmdtest.NewMockResponse(`[
				{
					"id": 2,
					"subject": "Unread message",
					"workflow_state": "unread"
				}
			]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Unread message") {
				t.Error("Expected 'Unread message' in output")
			}
		},
	}

	cmd := newConversationsListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}
