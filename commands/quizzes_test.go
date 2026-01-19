package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestQuizzesListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list quizzes successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts":  cmdtest.NewMockResponse(`[]`),
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/quizzes": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"title": "Midterm Exam",
						"quiz_type": "assignment",
						"points_possible": 100,
						"published": true
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Midterm Exam") {
					t.Error("Expected 'Midterm Exam' in output")
				}
			},
		},
		{
			Name: "list quizzes - empty response",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":         courseMock,
				"/api/v1/courses/1/quizzes": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No quizzes found",
		},
		{
			Name:        "list quizzes - missing course ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newQuizzesListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestQuizzesGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get quiz successfully",
			Args: []string{"10", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/quizzes/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"title": "Final Exam",
					"quiz_type": "assignment",
					"points_possible": 150,
					"published": true
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Final Exam") {
					t.Error("Expected 'Final Exam' in output")
				}
			},
		},
		{
			Name:        "get quiz - missing quiz ID",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "get quiz - missing course ID",
			Args:        []string{"10"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newQuizzesGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestQuizzesCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create quiz successfully",
			Args: []string{"--course-id", "1", "--title", "New Quiz", "--quiz-type", "assignment"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/quizzes": cmdtest.NewMockResponse(`{
					"id": 20,
					"title": "New Quiz",
					"quiz_type": "assignment",
					"published": false
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "New Quiz") {
					t.Error("Expected 'New Quiz' in output")
				}
			},
		},
		{
			Name:        "create quiz - missing title",
			Args:        []string{"--course-id", "1", "--quiz-type", "assignment"},
			ExpectError: true,
		},
		{
			Name:        "create quiz - missing course ID",
			Args:        []string{"--title", "New Quiz"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newQuizzesCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestQuizzesUpdateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "update quiz successfully",
			Args: []string{"10", "--course-id", "1", "--title", "Updated Quiz"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/quizzes/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"title": "Updated Quiz",
					"quiz_type": "assignment"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Updated Quiz") {
					t.Error("Expected 'Updated Quiz' in output")
				}
			},
		},
		{
			Name:        "update quiz - missing quiz ID",
			Args:        []string{"--course-id", "1", "--title", "Updated"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newQuizzesUpdateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestQuizzesDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete quiz with confirmation",
			Args: []string{"10", "--course-id", "1", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/quizzes/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"title": "Old Quiz"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete quiz - missing quiz ID",
			Args:        []string{"--course-id", "1", "--force"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newQuizzesDeleteCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestQuizzesQuestionsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list quiz questions successfully",
			Args: []string{"--course-id", "1", "--quiz-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/quizzes/10/questions": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"question_name": "Question 1",
						"question_type": "multiple_choice",
						"points_possible": 10
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Question 1") {
					t.Error("Expected 'Question 1' in output")
				}
			},
		},
		{
			Name: "list quiz questions - empty response",
			Args: []string{"--course-id", "1", "--quiz-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                      courseMock,
				"/api/v1/courses/1/quizzes/10/questions": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No questions found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newQuizzesQuestionsListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestQuizzesQuestionsGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get quiz question successfully",
			Args: []string{"5", "--course-id", "1", "--quiz-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/quizzes/10/questions/5": cmdtest.NewMockResponse(`{
					"id": 5,
					"question_name": "Question 5",
					"question_type": "essay",
					"points_possible": 20
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Question 5") {
					t.Error("Expected 'Question 5' in output")
				}
			},
		},
		{
			Name:        "get quiz question - missing question ID",
			Args:        []string{"--course-id", "1", "--quiz-id", "10"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newQuizzesQuestionsGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestQuizzesSubmissionsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list quiz submissions successfully",
			Args: []string{"--course-id", "1", "--quiz-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/quizzes/10/submissions": cmdtest.NewMockResponse(`{
					"quiz_submissions": [
						{
							"id": 1,
							"quiz_id": 10,
							"user_id": 100,
							"score": 85,
							"workflow_state": "complete"
						}
					]
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "1") {
					t.Error("Expected submission ID in output")
				}
			},
		},
		{
			Name:        "list quiz submissions - missing quiz ID",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newQuizzesSubmissionsListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestQuizzesSubmissionsGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get quiz submission successfully",
			Args: []string{"5", "--course-id", "1", "--quiz-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/quizzes/10/submissions/5": cmdtest.NewMockResponse(`{
					"quiz_submissions": [
						{
							"id": 5,
							"quiz_id": 10,
							"user_id": 100,
							"score": 90,
							"workflow_state": "complete"
						}
					]
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "100") {
					t.Error("Expected submission ID in output")
				}
			},
		},
		{
			Name:        "get quiz submission - missing submission ID",
			Args:        []string{"--course-id", "1", "--quiz-id", "10"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newQuizzesSubmissionsGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
