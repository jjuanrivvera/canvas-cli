package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestAssignmentsGetCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get assignment - API error",
		Args: []string{"--course-id", "1", "99"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":             courseMock,
			"/api/v1/courses/1/assignments": cmdtest.NewErrorResponse(404, "assignment not found"),
		},
		ExpectError: true,
	}
	cmd := newAssignmentsGetCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestAssignmentsCreateCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create assignment - API error",
		Args: []string{"--course-id", "1", "--name", "Failing Assignment"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":             courseMock,
			"/api/v1/courses/1/assignments": cmdtest.NewErrorResponse(422, "invalid assignment"),
		},
		ExpectError: true,
	}
	cmd := newAssignmentsCreateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestAssignmentsCreateCmd_WithPoints(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create assignment with points and grading type",
		Args: []string{
			"--course-id", "1",
			"--name", "Graded Assignment",
			"--points", "50",
			"--grading-type", "points",
			"--published",
		},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/assignments": cmdtest.NewMockResponse(`{
				"id": 30,
				"name": "Graded Assignment",
				"points_possible": 50,
				"grading_type": "points",
				"published": true
			}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Graded Assignment") {
				t.Error("Expected 'Graded Assignment' in output")
			}
		},
	}
	cmd := newAssignmentsCreateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestAssignmentsUpdateCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update assignment - API error",
		Args: []string{"--course-id", "1", "10", "--name", "Updated"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                courseMock,
			"/api/v1/courses/1/assignments/10": cmdtest.NewErrorResponse(404, "assignment not found"),
		},
		ExpectError: true,
	}
	cmd := newAssignmentsUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestAssignmentsUpdateCmd_WithAllFlags(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update assignment with multiple flags",
		Args: []string{
			"--course-id", "1", "10",
			"--name", "Updated Name",
			"--points", "75",
			"--grading-type", "percent",
			"--published",
		},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/assignments/10": cmdtest.NewMockResponse(`{
				"id": 10,
				"name": "Updated Name",
				"points_possible": 75,
				"published": true
			}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Updated Name") {
				t.Error("Expected 'Updated Name' in output")
			}
		},
	}
	cmd := newAssignmentsUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestAssignmentsDeleteCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete assignment - API error",
		Args: []string{"--course-id", "1", "10", "--force"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                courseMock,
			"/api/v1/courses/1/assignments/10": cmdtest.NewErrorResponse(404, "assignment not found"),
		},
		ExpectError: true,
	}
	cmd := newAssignmentsDeleteCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestAssignmentsListCmd_WithFilters(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list assignments with search term",
		Args: []string{"--course-id", "1", "--search", "quiz"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/assignments": cmdtest.NewMockResponse(`[
				{
					"id": 5,
					"name": "Midterm Quiz",
					"points_possible": 100,
					"published": true
				}
			]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Midterm Quiz") {
				t.Error("Expected 'Midterm Quiz' in output")
			}
		},
	}
	cmd := newAssignmentsListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestAssignmentsCreateCmd_FromJSONFile(t *testing.T) {
	tmpDir := t.TempDir()
	jsonPath := filepath.Join(tmpDir, "assignment.json")
	jsonContent := `{"name":"JSON Assignment","points_possible":75}`
	if err := os.WriteFile(jsonPath, []byte(jsonContent), 0644); err != nil {
		t.Fatal(err)
	}

	tc := cmdtest.CommandTestCase{
		Name: "create assignment from JSON file",
		Args: []string{"--course-id", "1", "--json", jsonPath},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/assignments": cmdtest.NewMockResponse(`{
				"id": 50,
				"name": "JSON Assignment",
				"points_possible": 75
			}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "JSON Assignment") {
				t.Error("Expected 'JSON Assignment' in output")
			}
		},
	}
	cmd := newAssignmentsCreateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestAssignmentsCreateCmd_InvalidJSONFile(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create assignment from nonexistent JSON file",
		Args: []string{"--course-id", "1", "--json", "/nonexistent/file.json"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
		},
		ExpectError: true,
	}
	cmd := newAssignmentsCreateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestAssignmentsUpdateCmd_FromJSONFile(t *testing.T) {
	tmpDir := t.TempDir()
	jsonPath := filepath.Join(tmpDir, "update.json")
	jsonContent := `{"name":"Updated via JSON","points_possible":80}`
	if err := os.WriteFile(jsonPath, []byte(jsonContent), 0644); err != nil {
		t.Fatal(err)
	}

	tc := cmdtest.CommandTestCase{
		Name: "update assignment from JSON file",
		Args: []string{"--course-id", "1", "10", "--json", jsonPath},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/assignments/10": cmdtest.NewMockResponse(`{
				"id": 10,
				"name": "Updated via JSON",
				"points_possible": 80
			}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Updated via JSON") {
				t.Error("Expected 'Updated via JSON' in output")
			}
		},
	}
	cmd := newAssignmentsUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestAssignmentsUpdateCmd_InvalidJSONFile(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update assignment from nonexistent JSON file",
		Args: []string{"--course-id", "1", "10", "--json", "/nonexistent/file.json"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
		},
		ExpectError: true,
	}
	cmd := newAssignmentsUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}
