package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestDiscussionsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list discussions successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/discussion_topics": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"title": "Welcome Discussion",
						"message": "Welcome to the course",
						"posted_at": "2024-01-01T00:00:00Z",
						"discussion_type": "threaded"
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Welcome Discussion") {
					t.Error("Expected 'Welcome Discussion' in output")
				}
			},
		},
		{
			Name: "list discussions - empty response",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts":                    cmdtest.NewMockResponse(`[]`),
				"/api/v1/courses/1/discussion_topics": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No discussion topics found",
		},
		{
			Name:        "list discussions - missing course ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newDiscussionsListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestDiscussionsGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get discussion successfully",
			Args: []string{"--course-id", "1", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/discussion_topics/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"title": "Week 1 Discussion",
					"message": "Discuss the readings",
					"posted_at": "2024-01-01T00:00:00Z",
					"discussion_type": "threaded"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Week 1 Discussion") {
					t.Error("Expected 'Week 1 Discussion' in output")
				}
			},
		},
		{
			Name:        "get discussion - missing course ID",
			Args:        []string{"10"},
			ExpectError: true,
		},
		{
			Name:        "get discussion - missing discussion ID",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newDiscussionsGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestDiscussionsCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create discussion successfully",
			Args: []string{"--course-id", "1", "--title", "New Discussion"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/discussion_topics": cmdtest.NewMockResponse(`{
					"id": 20,
					"title": "New Discussion",
					"message": null,
					"posted_at": "2024-01-01T00:00:00Z"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "New Discussion") {
					t.Error("Expected 'New Discussion' in output")
				}
			},
		},
		{
			Name:        "create discussion - missing course ID",
			Args:        []string{"--title", "New Discussion"},
			ExpectError: true,
		},
		{
			Name:        "create discussion - missing title",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newDiscussionsCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestDiscussionsDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete discussion with confirmation",
			Args: []string{"--course-id", "1", "10", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                      courseMock,
				"/api/v1/courses/1/discussion_topics/10": cmdtest.NewMockResponse(`{}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete discussion - missing course ID",
			Args:        []string{"10", "--force"},
			ExpectError: true,
		},
		{
			Name:        "delete discussion - missing discussion ID",
			Args:        []string{"--course-id", "1", "--force"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newDiscussionsDeleteCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
