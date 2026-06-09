package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestSubmissionsListCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list submissions - API error",
		Args: []string{"--course-id", "1", "--assignment-id", "100"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                             courseMock,
			"/api/v1/courses/1/assignments/100/submissions": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmd := newSubmissionsListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestSubmissionsGetCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get submission - API error",
		Args: []string{"--course-id", "1", "--assignment-id", "100", "--user-id", "99"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/assignments/100/submissions/99": cmdtest.NewErrorResponse(404, "submission not found"),
		},
		ExpectError: true,
	}
	cmd := newSubmissionsGetCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestSubmissionsGradeCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "grade submission - API error",
		Args: []string{"--course-id", "1", "--assignment-id", "100", "--user-id", "10", "--score", "95"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/assignments/100/submissions/10": cmdtest.NewErrorResponse(422, "invalid grade"),
		},
		ExpectError: true,
	}
	cmd := newSubmissionsGradeCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestSubmissionsGradeCmd_WithComment(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "grade submission with comment",
		Args: []string{"--course-id", "1", "--assignment-id", "100", "--user-id", "10", "--score", "85", "--comment", "Good work"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/assignments/100/submissions/10": cmdtest.NewMockResponse(`{
				"id": 1,
				"assignment_id": 100,
				"user_id": 10,
				"score": 85,
				"grade": "85.00"
			}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Successfully graded") {
				t.Error("Expected 'Successfully graded' in output")
			}
		},
	}
	cmd := newSubmissionsGradeCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestSubmissionsGradeCmd_WithExcuse(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "grade submission with excuse",
		Args: []string{"--course-id", "1", "--assignment-id", "100", "--user-id", "10", "--excuse"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/assignments/100/submissions/10": cmdtest.NewMockResponse(`{
				"id": 1,
				"assignment_id": 100,
				"user_id": 10,
				"excused": true
			}`),
		},
		ExpectError: false,
	}
	cmd := newSubmissionsGradeCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestSubmissionsCommentsCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list submission comments - API error",
		Args: []string{"--course-id", "1", "--assignment-id", "100", "--user-id", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/assignments/100/submissions/10": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmd := newSubmissionsCommentsCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestSubmissionsAddCommentCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "add comment - API error",
		Args: []string{"--course-id", "1", "--assignment-id", "100", "--user-id", "10", "--text", "Hi"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/assignments/100/submissions/10": cmdtest.NewErrorResponse(422, "invalid comment"),
		},
		ExpectError: true,
	}
	cmd := newSubmissionsAddCommentCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestSubmissionsDeleteCommentCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete comment - API error",
		Args: []string{"--course-id", "1", "--assignment-id", "100", "--user-id", "10", "--comment-id", "5"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/assignments/100/submissions/10/comments/5": cmdtest.NewErrorResponse(404, "comment not found"),
		},
		ExpectError: true,
	}
	cmd := newSubmissionsDeleteCommentCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestSubmissionsBulkGradeCmd_MissingCSV(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "bulk grade - missing csv flag",
		Args:        []string{"--course-id", "1"},
		ExpectError: true,
	}
	cmd := newSubmissionsBulkGradeCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestSubmissionsBulkGradeCmd_NonexistentCSV(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "bulk grade - nonexistent CSV file",
		Args: []string{"--course-id", "1", "--csv", "/nonexistent/path/grades.csv"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
		},
		ExpectError: true,
	}
	cmd := newSubmissionsBulkGradeCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestSubmissionsBulkGradeCmd_EmptyCSV(t *testing.T) {
	// Create a temporary CSV with only a header (no data rows)
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "empty_grades.csv")
	if err := os.WriteFile(csvPath, []byte("user_id,assignment_id,grade\n"), 0644); err != nil {
		t.Fatal(err)
	}

	tc := cmdtest.CommandTestCase{
		Name: "bulk grade - empty CSV file",
		Args: []string{"--course-id", "1", "--csv", csvPath},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
		},
		ExpectError: true,
	}
	cmd := newSubmissionsBulkGradeCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestSubmissionsBulkGradeCmd_DryRun(t *testing.T) {
	// Create a temporary CSV with grade data
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "grades.csv")
	csvContent := "user_id,assignment_id,grade\n10,100,95\n20,100,87\n"
	if err := os.WriteFile(csvPath, []byte(csvContent), 0644); err != nil {
		t.Fatal(err)
	}

	tc := cmdtest.CommandTestCase{
		Name: "bulk grade - dry run",
		Args: []string{"--course-id", "1", "--csv", csvPath, "--dry-run"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "DRY RUN") {
				t.Error("Expected 'DRY RUN' in output")
			}
		},
	}
	cmd := newSubmissionsBulkGradeCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestSubmissionsBulkGradeCmd_Success(t *testing.T) {
	// Create a temporary CSV with grade data
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "grades.csv")
	csvContent := "user_id,assignment_id,grade\n10,100,95\n"
	if err := os.WriteFile(csvPath, []byte(csvContent), 0644); err != nil {
		t.Fatal(err)
	}

	tc := cmdtest.CommandTestCase{
		Name: "bulk grade - success",
		Args: []string{"--course-id", "1", "--csv", csvPath},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/assignments/100/submissions/10": cmdtest.NewMockResponse(`{
				"id": 1,
				"assignment_id": 100,
				"user_id": 10,
				"score": 95,
				"grade": "95"
			}`),
		},
		ExpectError: false,
	}
	cmd := newSubmissionsBulkGradeCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}
