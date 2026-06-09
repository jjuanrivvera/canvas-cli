package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestGradesColumnsUpdateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "update custom column title successfully",
			Args: []string{"10", "--course-id", "1", "--title", "New Title"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/custom_gradebook_columns/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"title": "New Title",
					"position": 1,
					"hidden": false
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "New Title") {
					t.Error("Expected 'New Title' in output")
				}
			},
		},
		{
			Name:        "update custom column - missing column ID",
			Args:        []string{"--course-id", "1", "--title", "New Title"},
			ExpectError: true,
		},
		{
			Name:        "update custom column - missing course ID",
			Args:        []string{"10", "--title", "New Title"},
			ExpectError: true,
		},
		{
			Name: "update custom column - API error",
			Args: []string{"10", "--course-id", "1", "--title", "New Title"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                             courseMock,
				"/api/v1/courses/1/custom_gradebook_columns/10": cmdtest.NewErrorResponse(404, "column not found"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newGradesColumnsUpdateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestGradesHistoryCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "grades history - API error",
		Args: []string{"--course-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts":                         cmdtest.NewMockResponse(`[]`),
			"/api/v1/courses/1/gradebook_history/days": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmd := newGradesHistoryCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestGradesFeedCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "grades feed - API error",
		Args: []string{"--course-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                        courseMock,
			"/api/v1/courses/1/gradebook_history/feed": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmd := newGradesFeedCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestGradesColumnsListCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list custom columns - API error",
		Args: []string{"--course-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                          courseMock,
			"/api/v1/courses/1/custom_gradebook_columns": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmd := newGradesColumnsListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestGradesColumnsGetCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get custom column - API error",
		Args: []string{"10", "--course-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                             courseMock,
			"/api/v1/courses/1/custom_gradebook_columns/10": cmdtest.NewErrorResponse(404, "column not found"),
		},
		ExpectError: true,
	}
	cmd := newGradesColumnsGetCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestGradesColumnsCreateCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create custom column - API error",
		Args: []string{"--course-id", "1", "--title", "Bad Column"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                          courseMock,
			"/api/v1/courses/1/custom_gradebook_columns": cmdtest.NewErrorResponse(422, "invalid column"),
		},
		ExpectError: true,
	}
	cmd := newGradesColumnsCreateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestGradesColumnsDeleteCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete custom column - API error",
		Args: []string{"10", "--course-id", "1", "--force"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                             courseMock,
			"/api/v1/courses/1/custom_gradebook_columns/10": cmdtest.NewErrorResponse(404, "column not found"),
		},
		ExpectError: true,
	}
	cmd := newGradesColumnsDeleteCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestGradesColumnsDataListCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list column data - API error",
		Args: []string{"10", "--course-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/custom_gradebook_columns/10/data": cmdtest.NewErrorResponse(404, "column not found"),
		},
		ExpectError: true,
	}
	cmd := newGradesColumnsDataListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestGradesColumnsDataSetCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "set column data - API error",
		Args: []string{"10", "--course-id", "1", "--user-id", "100", "--content", "Present"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/custom_gradebook_columns/10/data/100": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmd := newGradesColumnsDataSetCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestGradesColumnsUpdateCmd_WithPosition(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update custom column position",
		Args: []string{"10", "--course-id", "1", "--position", "3"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/custom_gradebook_columns/10": cmdtest.NewMockResponse(`{
				"id": 10,
				"title": "Attendance",
				"position": 3
			}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Attendance") {
				t.Error("Expected 'Attendance' in output")
			}
		},
	}
	cmd := newGradesColumnsUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}
