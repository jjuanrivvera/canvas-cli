package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestPagesListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list pages successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/pages": cmdtest.NewMockResponse(`[
					{
						"page_id": 1,
						"url": "home-page",
						"title": "Home Page",
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-02T00:00:00Z",
						"published": true
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Home Page") {
					t.Error("Expected 'Home Page' in output")
				}
			},
		},
		{
			Name: "list pages - empty response",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":       courseMock,
				"/api/v1/courses/1/pages": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No pages found",
		},
		{
			Name:        "list pages - missing course ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPagesListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPagesGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get page successfully",
			Args: []string{"--course-id", "1", "home-page"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/pages/home-page": cmdtest.NewMockResponse(`{
					"page_id": 1,
					"url": "home-page",
					"title": "Home Page",
					"body": "<p>Welcome</p>",
					"published": true
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Home Page") {
					t.Error("Expected 'Home Page' in output")
				}
			},
		},
		{
			Name:        "get page - missing course ID",
			Args:        []string{"home-page"},
			ExpectError: true,
		},
		{
			Name:        "get page - missing page URL",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPagesGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPagesCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create page successfully",
			Args: []string{"--course-id", "1", "--title", "New Page"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/pages": cmdtest.NewMockResponse(`{
					"page_id": 20,
					"url": "new-page",
					"title": "New Page",
					"body": "",
					"published": false
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "New Page") {
					t.Error("Expected 'New Page' in output")
				}
			},
		},
		{
			Name:        "create page - missing course ID",
			Args:        []string{"--title", "New Page"},
			ExpectError: true,
		},
		{
			Name:        "create page - missing title",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPagesCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPagesDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete page with confirmation",
			Args: []string{"--course-id", "1", "old-page", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                courseMock,
				"/api/v1/courses/1/pages/old-page": cmdtest.NewMockResponse(`{}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete page - missing course ID",
			Args:        []string{"old-page", "--force"},
			ExpectError: true,
		},
		{
			Name:        "delete page - missing page URL",
			Args:        []string{"--course-id", "1", "--force"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPagesDeleteCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPagesFrontCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get front page successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/front_page": cmdtest.NewMockResponse(`{
					"page_id": 1,
					"url": "home",
					"title": "Course Home",
					"published": true
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Course Home") {
					t.Error("Expected 'Course Home' in output")
				}
			},
		},
		{
			Name:        "get front page - missing course ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPagesFrontCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPagesUpdateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "update page title successfully",
			Args: []string{"--course-id", "1", "home-page", "--title", "Updated Home"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/pages/home-page": cmdtest.NewMockResponse(`{
					"page_id": 1,
					"url": "home-page",
					"title": "Updated Home",
					"published": true
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Updated Home") {
					t.Error("Expected 'Updated Home' in output")
				}
			},
		},
		{
			Name:        "update page - missing course ID",
			Args:        []string{"home-page", "--title", "Updated"},
			ExpectError: true,
		},
		{
			Name:        "update page - missing page URL",
			Args:        []string{"--course-id", "1", "--title", "Updated"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPagesUpdateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPagesDuplicateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "duplicate page successfully",
			Args: []string{"--course-id", "1", "home-page"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/pages/home-page/duplicate": cmdtest.NewMockResponse(`{
					"page_id": 5,
					"url": "home-page-copy",
					"title": "Home Page Copy",
					"published": false
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Home Page Copy") {
					t.Error("Expected 'Home Page Copy' in output")
				}
			},
		},
		{
			Name:        "duplicate page - missing course ID",
			Args:        []string{"home-page"},
			ExpectError: true,
		},
		{
			Name:        "duplicate page - missing page URL",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPagesDuplicateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPagesRevisionsCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get page revisions successfully",
			Args: []string{"--course-id", "1", "home-page"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/pages/home-page/revisions": cmdtest.NewMockResponse(`[
					{
						"revision_id": 1,
						"updated_at": "2024-01-01T00:00:00Z",
						"url": "home-page"
					}
				]`),
			},
			ExpectError: false,
		},
		{
			Name:        "get page revisions - missing course ID",
			Args:        []string{"home-page"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPagesRevisionsCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPagesRevertCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "revert page to revision successfully",
			Args: []string{"--course-id", "1", "home-page", "2"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/pages/home-page/revisions/2": cmdtest.NewMockResponse(`{
					"page_id": 1,
					"url": "home-page",
					"title": "Reverted Home",
					"published": true
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "revert page - missing course ID",
			Args:        []string{"home-page", "2"},
			ExpectError: true,
		},
		{
			Name:        "revert page - missing arguments",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPagesRevertCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPagesListCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list pages - API error",
		Args: []string{"--course-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":       courseMock,
			"/api/v1/courses/1/pages": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}

	cmd := newPagesListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}
