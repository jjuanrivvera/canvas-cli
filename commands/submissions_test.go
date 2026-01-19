package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestSubmissionsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list submissions successfully",
			Args: []string{"--course-id", "1", "--assignment-id", "100"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/assignments/100/submissions": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"assignment_id": 100,
						"user_id": 10,
						"score": 95,
						"grade": "A",
						"workflow_state": "graded"
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "graded") {
					t.Error("Expected 'graded' in output")
				}
			},
		},
		{
			Name: "list submissions - empty response",
			Args: []string{"--course-id", "1", "--assignment-id", "100"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                             courseMock,
				"/api/v1/courses/1/assignments/100/submissions": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No submissions found",
		},
		{
			Name:        "list submissions - missing assignment ID",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newSubmissionsListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestSubmissionsGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get submission successfully",
			Args: []string{"--course-id", "1", "--assignment-id", "100", "--user-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/assignments/100/submissions/10": cmdtest.NewMockResponse(`{
					"id": 1,
					"assignment_id": 100,
					"user_id": 10,
					"score": 95,
					"grade": "A",
					"workflow_state": "graded"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "95") {
					t.Error("Expected '95' in output")
				}
			},
		},
		{
			Name:        "get submission - missing user ID",
			Args:        []string{"--course-id", "1", "--assignment-id", "100"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newSubmissionsGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestSubmissionsGradeCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "grade submission successfully",
			Args: []string{"--course-id", "1", "--assignment-id", "100", "--user-id", "10", "--score", "95"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/assignments/100/submissions/10": cmdtest.NewMockResponse(`{
					"id": 1,
					"assignment_id": 100,
					"user_id": 10,
					"score": 95,
					"grade": "95.00",
					"workflow_state": "graded"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Successfully graded") {
					t.Error("Expected 'Successfully graded' in output")
				}
			},
		},
		{
			Name:        "grade submission - missing user ID",
			Args:        []string{"--course-id", "1", "--assignment-id", "100", "--score", "95"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newSubmissionsGradeCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestSubmissionsCommentsCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list submission comments successfully",
			Args: []string{"--course-id", "1", "--assignment-id", "100", "--user-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/assignments/100/submissions/10": cmdtest.NewMockResponse(`{
					"id": 1,
					"assignment_id": 100,
					"user_id": 10,
					"submission_comments": [
						{
							"id": 1,
							"comment": "Great work!",
							"author_name": "Teacher"
						}
					]
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Great work") {
					t.Error("Expected 'Great work' in output")
				}
			},
		},
		{
			Name: "list submission comments - empty response",
			Args: []string{"--course-id", "1", "--assignment-id", "100", "--user-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/assignments/100/submissions/10": cmdtest.NewMockResponse(`{
					"id": 1,
					"submission_comments": []
				}`),
			},
			ExpectError:  false,
			ExpectOutput: "No comments found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newSubmissionsCommentsCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestSubmissionsAddCommentCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "add comment successfully",
			Args: []string{"--course-id", "1", "--assignment-id", "100", "--user-id", "10", "--text", "Well done!"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/assignments/100/submissions/10": cmdtest.NewMockResponse(`{
					"id": 1,
					"assignment_id": 100,
					"user_id": 10
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Comment added successfully") {
					t.Error("Expected 'Comment added successfully' in output")
				}
			},
		},
		{
			Name:        "add comment - missing text",
			Args:        []string{"--course-id", "1", "--assignment-id", "100", "--user-id", "10"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newSubmissionsAddCommentCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestSubmissionsDeleteCommentCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete comment successfully",
			Args: []string{"--course-id", "1", "--assignment-id", "100", "--user-id", "10", "--comment-id", "5"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/assignments/100/submissions/10/comments/5": cmdtest.NewMockResponse(`{
					"id": 5
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "deleted successfully") {
					t.Error("Expected 'deleted successfully' in output")
				}
			},
		},
		{
			Name:        "delete comment - missing comment ID",
			Args:        []string{"--course-id", "1", "--assignment-id", "100", "--user-id", "10"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newSubmissionsDeleteCommentCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
