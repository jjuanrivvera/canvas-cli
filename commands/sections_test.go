package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestSectionsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list sections successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/sections": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"name": "Section A",
						"course_id": 1,
						"start_at": null,
						"end_at": null
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Section A") {
					t.Error("Expected 'Section A' in output")
				}
			},
		},
		{
			Name: "list sections - empty response",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":          courseMock,
				"/api/v1/courses/1/sections": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No sections found",
		},
		{
			Name:        "list sections - missing course ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newSectionsListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestSectionsGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get section successfully",
			Args: []string{"10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/sections/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"name": "Section B",
					"course_id": 1,
					"start_at": null,
					"end_at": null
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Section B") {
					t.Error("Expected 'Section B' in output")
				}
			},
		},
		{
			Name:        "get section - missing section ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newSectionsGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestSectionsCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create section successfully",
			Args: []string{"--course-id", "1", "--name", "New Section"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/sections": cmdtest.NewMockResponse(`{
					"id": 20,
					"name": "New Section",
					"course_id": 1
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "New Section") {
					t.Error("Expected 'New Section' in output")
				}
			},
		},
		{
			Name:        "create section - missing course ID",
			Args:        []string{"--name", "New Section"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newSectionsCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestSectionsDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete section with confirmation",
			Args: []string{"10", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/sections/10": cmdtest.NewMockResponse(`{}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete section - missing section ID",
			Args:        []string{"--force"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newSectionsDeleteCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestSectionsUpdateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "update section name successfully",
			Args: []string{"10", "--name", "Updated Section"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/sections/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"name": "Updated Section",
					"course_id": 1
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Updated Section") {
					t.Error("Expected 'Updated Section' in output")
				}
			},
		},
		{
			Name:        "update section - missing section ID",
			Args:        []string{"--name", "Updated"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newSectionsUpdateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestSectionsCrosslistCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "crosslist section successfully",
			Args: []string{"10", "--new-course-id", "99"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/sections/10/crosslist/99": cmdtest.NewMockResponse(`{
					"id": 10,
					"name": "Section A",
					"course_id": 99
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "crosslist section - missing new course ID",
			Args:        []string{"10"},
			ExpectError: true,
		},
		{
			Name:        "crosslist section - missing section ID",
			Args:        []string{"--new-course-id", "99"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newSectionsCrosslistCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestSectionsUncrosslistCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "uncrosslist section successfully",
			Args: []string{"10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/sections/10/crosslist": cmdtest.NewMockResponse(`{
					"id": 10,
					"name": "Section A",
					"course_id": 1
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "uncrosslist section - missing section ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newSectionsUncrosslistCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestSectionsCreateCmd_MissingName(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "create section - missing name",
		Args:        []string{"--course-id", "1"},
		ExpectError: true,
	}

	cmd := newSectionsCreateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestSectionsListCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list sections - API error",
		Args: []string{"--course-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":          courseMock,
			"/api/v1/courses/1/sections": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}

	cmd := newSectionsListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}
