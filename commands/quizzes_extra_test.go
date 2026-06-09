package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestQuizzesQuestionsCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create quiz question successfully",
			Args: []string{"--course-id", "1", "--quiz-id", "10", "--type", "multiple_choice", "--text", "What is 2+2?"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/quizzes/10/questions": cmdtest.NewMockResponse(`{
					"id": 30,
					"question_name": "Question",
					"question_type": "multiple_choice",
					"question_text": "What is 2+2?",
					"points_possible": 1
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Question created") {
					t.Error("Expected 'Question created' in output")
				}
			},
		},
		{
			Name:        "create quiz question - missing course ID",
			Args:        []string{"--quiz-id", "10", "--type", "multiple_choice", "--text", "Q?"},
			ExpectError: true,
		},
		{
			Name:        "create quiz question - missing quiz ID",
			Args:        []string{"--course-id", "1", "--type", "multiple_choice", "--text", "Q?"},
			ExpectError: true,
		},
		{
			Name: "create quiz question - API error",
			Args: []string{"--course-id", "1", "--quiz-id", "10", "--type", "multiple_choice", "--text", "Q?"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                      courseMock,
				"/api/v1/courses/1/quizzes/10/questions": cmdtest.NewErrorResponse(422, "invalid question"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newQuizzesQuestionsCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestQuizzesQuestionsDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete quiz question with force",
			Args: []string{"5", "--course-id", "1", "--quiz-id", "10", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                        courseMock,
				"/api/v1/courses/1/quizzes/10/questions/5": cmdtest.NewMockResponse(`{}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete quiz question - missing question ID",
			Args:        []string{"--course-id", "1", "--quiz-id", "10", "--force"},
			ExpectError: true,
		},
		{
			Name:        "delete quiz question - missing course ID",
			Args:        []string{"5", "--quiz-id", "10", "--force"},
			ExpectError: true,
		},
		{
			Name: "delete quiz question - API error",
			Args: []string{"5", "--course-id", "1", "--quiz-id", "10", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                        courseMock,
				"/api/v1/courses/1/quizzes/10/questions/5": cmdtest.NewErrorResponse(404, "question not found"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newQuizzesQuestionsDeleteCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestQuizzesListCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list quizzes - API error",
		Args: []string{"--course-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":         courseMock,
			"/api/v1/courses/1/quizzes": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmd := newQuizzesListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestQuizzesGetCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get quiz - API error",
		Args: []string{"99", "--course-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":            courseMock,
			"/api/v1/courses/1/quizzes/99": cmdtest.NewErrorResponse(404, "quiz not found"),
		},
		ExpectError: true,
	}
	cmd := newQuizzesGetCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestQuizzesCreateCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create quiz - API error",
		Args: []string{"--course-id", "1", "--title", "Bad Quiz"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":         courseMock,
			"/api/v1/courses/1/quizzes": cmdtest.NewErrorResponse(422, "invalid quiz"),
		},
		ExpectError: true,
	}
	cmd := newQuizzesCreateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestQuizzesUpdateCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update quiz - API error",
		Args: []string{"10", "--course-id", "1", "--title", "Updated Quiz"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":            courseMock,
			"/api/v1/courses/1/quizzes/10": cmdtest.NewErrorResponse(404, "quiz not found"),
		},
		ExpectError: true,
	}
	cmd := newQuizzesUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestQuizzesDeleteCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete quiz - API error",
		Args: []string{"10", "--course-id", "1", "--force"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":            courseMock,
			"/api/v1/courses/1/quizzes/10": cmdtest.NewErrorResponse(404, "quiz not found"),
		},
		ExpectError: true,
	}
	cmd := newQuizzesDeleteCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestQuizzesQuestionsListCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list quiz questions - API error",
		Args: []string{"--course-id", "1", "--quiz-id", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                      courseMock,
			"/api/v1/courses/1/quizzes/10/questions": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmd := newQuizzesQuestionsListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestQuizzesQuestionsGetCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get quiz question - API error",
		Args: []string{"99", "--course-id", "1", "--quiz-id", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                         courseMock,
			"/api/v1/courses/1/quizzes/10/questions/99": cmdtest.NewErrorResponse(404, "question not found"),
		},
		ExpectError: true,
	}
	cmd := newQuizzesQuestionsGetCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestQuizzesSubmissionsListCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list quiz submissions - API error",
		Args: []string{"--course-id", "1", "--quiz-id", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                        courseMock,
			"/api/v1/courses/1/quizzes/10/submissions": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmd := newQuizzesSubmissionsListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestQuizzesSubmissionsGetCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get quiz submission - API error",
		Args: []string{"5", "--course-id", "1", "--quiz-id", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                          courseMock,
			"/api/v1/courses/1/quizzes/10/submissions/5": cmdtest.NewErrorResponse(404, "submission not found"),
		},
		ExpectError: true,
	}
	cmd := newQuizzesSubmissionsGetCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}
