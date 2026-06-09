package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestAnnouncementsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list announcements successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/announcements": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"title": "Important Announcement",
						"message": "Please read this",
						"posted_at": "2024-01-01T00:00:00Z",
						"is_announcement": true
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Important Announcement") {
					t.Error("Expected 'Important Announcement' in output")
				}
			},
		},
		{
			Name: "list announcements - empty response",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/announcements": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No announcements found",
		},
		{
			Name:        "list announcements - missing course ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAnnouncementsListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAnnouncementsCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create announcement successfully",
			Args: []string{"--course-id", "1", "--title", "New Announcement"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/discussion_topics": cmdtest.NewMockResponse(`{
					"id": 20,
					"title": "New Announcement",
					"message": null,
					"posted_at": "2024-01-01T00:00:00Z",
					"is_announcement": true
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "New Announcement") {
					t.Error("Expected 'New Announcement' in output")
				}
			},
		},
		{
			Name:        "create announcement - missing course ID",
			Args:        []string{"--title", "New Announcement"},
			ExpectError: true,
		},
		{
			Name:        "create announcement - missing title",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAnnouncementsCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAnnouncementsDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete announcement with confirmation",
			Args: []string{"--course-id", "1", "10", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                      courseMock,
				"/api/v1/courses/1/discussion_topics/10": cmdtest.NewMockResponse(`{}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete announcement - missing course ID",
			Args:        []string{"10", "--force"},
			ExpectError: true,
		},
		{
			Name:        "delete announcement - missing announcement ID",
			Args:        []string{"--course-id", "1", "--force"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAnnouncementsDeleteCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAnnouncementsGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get announcement successfully",
			Args: []string{"--course-id", "1", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/discussion_topics/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"title": "Class Update",
					"message": "Important information",
					"posted_at": "2024-01-01T00:00:00Z",
					"is_announcement": true
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Class Update") {
					t.Error("Expected 'Class Update' in output")
				}
			},
		},
		{
			Name:        "get announcement - missing course ID",
			Args:        []string{"10"},
			ExpectError: true,
		},
		{
			Name:        "get announcement - missing announcement ID",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAnnouncementsGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAnnouncementsUpdateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "update announcement title successfully",
			Args: []string{"--course-id", "1", "10", "--title", "Updated Title"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/discussion_topics/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"title": "Updated Title",
					"message": "Original message",
					"posted_at": "2024-01-01T00:00:00Z",
					"is_announcement": true
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Updated Title") {
					t.Error("Expected 'Updated Title' in output")
				}
			},
		},
		{
			Name:        "update announcement - missing announcement ID",
			Args:        []string{"--course-id", "1", "--title", "Updated"},
			ExpectError: true,
		},
		{
			Name:        "update announcement - missing course ID",
			Args:        []string{"10", "--title", "Updated"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAnnouncementsUpdateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAnnouncementsListCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list announcements - API error",
		Args: []string{"--course-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/announcements": cmdtest.NewErrorResponse(500, "internal server error"),
		},
		ExpectError: true,
	}

	cmd := newAnnouncementsListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}
