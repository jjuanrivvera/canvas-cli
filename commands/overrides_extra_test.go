package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestOverridesCreateCmd_WithStudentIDs(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create override with student IDs",
		Args: []string{"--course-id", "1", "--assignment-id", "100", "--student-ids", "10,20,30"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/assignments/100/overrides": cmdtest.NewMockResponse(`{
				"id": 21,
				"assignment_id": 100,
				"student_ids": [10, 20, 30]
			}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Override created") {
				t.Error("Expected 'Override created' in output")
			}
		},
	}
	cmd := newOverridesCreateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestOverridesCreateCmd_InvalidStudentIDs(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create override - invalid student ID format",
		Args: []string{"--course-id", "1", "--assignment-id", "100", "--student-ids", "10,abc,30"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
		},
		ExpectError: true,
	}
	cmd := newOverridesCreateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestOverridesCreateCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create override - API error",
		Args: []string{"--course-id", "1", "--assignment-id", "100", "--section-id", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                           courseMock,
			"/api/v1/courses/1/assignments/100/overrides": cmdtest.NewErrorResponse(422, "invalid override"),
		},
		ExpectError: true,
	}
	cmd := newOverridesCreateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestOverridesUpdateCmd_WithStudentIDs(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update override with student IDs",
		Args: []string{"5", "--course-id", "1", "--assignment-id", "100", "--student-ids", "10,20"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/assignments/100/overrides/5": cmdtest.NewMockResponse(`{
				"id": 5,
				"assignment_id": 100,
				"student_ids": [10, 20]
			}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Override updated") {
				t.Error("Expected 'Override updated' in output")
			}
		},
	}
	cmd := newOverridesUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestOverridesUpdateCmd_InvalidStudentIDs(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update override - invalid student ID format",
		Args: []string{"5", "--course-id", "1", "--assignment-id", "100", "--student-ids", "bad,ids"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
		},
		ExpectError: true,
	}
	cmd := newOverridesUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestOverridesUpdateCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update override - API error",
		Args: []string{"5", "--course-id", "1", "--assignment-id", "100", "--title", "Updated"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                             courseMock,
			"/api/v1/courses/1/assignments/100/overrides/5": cmdtest.NewErrorResponse(404, "override not found"),
		},
		ExpectError: true,
	}
	cmd := newOverridesUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestOverridesDeleteCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete override - API error",
		Args: []string{"5", "--course-id", "1", "--assignment-id", "100", "--force"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                             courseMock,
			"/api/v1/courses/1/assignments/100/overrides/5": cmdtest.NewErrorResponse(404, "override not found"),
		},
		ExpectError: true,
	}
	cmd := newOverridesDeleteCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestOverridesGetCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get override - API error",
		Args: []string{"99", "--course-id", "1", "--assignment-id", "100"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/assignments/100/overrides/99": cmdtest.NewErrorResponse(404, "override not found"),
		},
		ExpectError: true,
	}
	cmd := newOverridesGetCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestOverridesListCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list overrides - API error",
		Args: []string{"--course-id", "1", "--assignment-id", "100"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                           courseMock,
			"/api/v1/courses/1/assignments/100/overrides": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmd := newOverridesListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}
