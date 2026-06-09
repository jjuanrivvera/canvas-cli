package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestDiscussionsReplyCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			// topic-id=10, entry-id=55, message as positional arg
			Name: "reply to discussion entry successfully",
			Args: []string{"--course-id", "1", "10", "55", "My reply"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/discussion_topics/10/entries/55/replies": cmdtest.NewMockResponse(`{
					"id": 60,
					"message": "My reply",
					"user_id": 100
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Reply posted successfully") {
					t.Error("Expected success message in output")
				}
			},
		},
		{
			// topic-id=10, entry-id=55, message via --message flag
			Name: "reply to discussion entry via --message flag",
			Args: []string{"--course-id", "1", "10", "55", "--message", "Flag reply"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/discussion_topics/10/entries/55/replies": cmdtest.NewMockResponse(`{
					"id": 61,
					"message": "Flag reply",
					"user_id": 100
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "reply - missing course ID",
			Args:        []string{"10", "55", "Hi"},
			ExpectError: true,
		},
		{
			// Missing entry-id positional arg means only 1 positional arg (below min 2)
			Name:        "reply - missing entry ID arg",
			Args:        []string{"--course-id", "1", "10"},
			ExpectError: true,
		},
		{
			Name: "reply - API error",
			Args: []string{"--course-id", "1", "10", "55", "My reply"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/discussion_topics/10/entries/55/replies": cmdtest.NewErrorResponse(404, "entry not found"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newDiscussionsReplyCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestDiscussionsListCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list discussions - API error",
		Args: []string{"--course-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                   courseMock,
			"/api/v1/courses/1/discussion_topics": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmd := newDiscussionsListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestDiscussionsGetCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get discussion - API error",
		Args: []string{"--course-id", "1", "99"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                      courseMock,
			"/api/v1/courses/1/discussion_topics/99": cmdtest.NewErrorResponse(404, "topic not found"),
		},
		ExpectError: true,
	}
	cmd := newDiscussionsGetCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestDiscussionsCreateCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create discussion - API error",
		Args: []string{"--course-id", "1", "--title", "Bad Discussion"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                   courseMock,
			"/api/v1/courses/1/discussion_topics": cmdtest.NewErrorResponse(422, "invalid topic"),
		},
		ExpectError: true,
	}
	cmd := newDiscussionsCreateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestDiscussionsUpdateCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update discussion - API error",
		Args: []string{"--course-id", "1", "10", "--title", "Bad Update"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                      courseMock,
			"/api/v1/courses/1/discussion_topics/10": cmdtest.NewErrorResponse(404, "topic not found"),
		},
		ExpectError: true,
	}
	cmd := newDiscussionsUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestDiscussionsDeleteCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete discussion - API error",
		Args: []string{"--course-id", "1", "10", "--force"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                      courseMock,
			"/api/v1/courses/1/discussion_topics/10": cmdtest.NewErrorResponse(404, "topic not found"),
		},
		ExpectError: true,
	}
	cmd := newDiscussionsDeleteCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestDiscussionsEntriesCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list discussion entries - API error",
		Args: []string{"--course-id", "1", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/discussion_topics/10/entries": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmd := newDiscussionsEntriesCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestDiscussionsPostCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "post entry - API error",
		Args: []string{"--course-id", "1", "10", "--message", "Hi"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/discussion_topics/10/entries": cmdtest.NewErrorResponse(422, "invalid entry"),
		},
		ExpectError: true,
	}
	cmd := newDiscussionsPostCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestDiscussionsSubscribeCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "subscribe - API error",
		Args: []string{"--course-id", "1", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/discussion_topics/10/subscribed": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmd := newDiscussionsSubscribeCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestDiscussionsUnsubscribeCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "unsubscribe - API error",
		Args: []string{"--course-id", "1", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/discussion_topics/10/subscribed": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmd := newDiscussionsUnsubscribeCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestDiscussionsUpdateCmd_WithAllFlags(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update discussion with message and pinned",
		Args: []string{"--course-id", "1", "10", "--title", "Pinned Discussion", "--message", "Important!", "--pinned"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/discussion_topics/10": cmdtest.NewMockResponse(`{
				"id": 10,
				"title": "Pinned Discussion",
				"message": "Important!",
				"pinned": true
			}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Pinned Discussion") {
				t.Error("Expected 'Pinned Discussion' in output")
			}
		},
	}
	cmd := newDiscussionsUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}
