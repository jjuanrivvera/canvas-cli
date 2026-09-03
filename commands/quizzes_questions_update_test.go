package commands

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

func TestQuizzesQuestionsUpdateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "update question successfully",
			Args: []string{"789", "--course-id", "1", "--quiz-id", "10", "--name", "Renamed", "--points", "5"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/quizzes/10/questions/789": cmdtest.NewMockResponse(`{
					"id": 789,
					"quiz_id": 10,
					"question_name": "Renamed",
					"question_type": "multiple_choice_question",
					"question_text": "What is 2+2?",
					"points_possible": 5
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Renamed") {
					t.Errorf("Expected 'Renamed' in output, got: %s", output)
				}
				if !strings.Contains(output, "Question updated successfully (ID: 789)") {
					t.Errorf("Expected success message in output, got: %s", output)
				}
			},
		},
		{
			Name: "update question with answers JSON",
			Args: []string{"789", "--course-id", "1", "--quiz-id", "10",
				"--answers-json", `[{"id":1001,"text":"4","weight":100},{"id":1002,"text":"5","weight":0}]`},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/quizzes/10/questions/789": cmdtest.NewMockResponse(`{
					"id": 789,
					"quiz_id": 10,
					"question_name": "Arithmetic",
					"question_type": "multiple_choice_question",
					"answers": [{"id":1001,"text":"4","weight":100},{"id":1002,"text":"5","weight":0}]
				}`),
			},
			ExpectError:  false,
			ExpectOutput: "Arithmetic",
		},
		{
			Name:        "update question - malformed answers JSON",
			Args:        []string{"789", "--course-id", "1", "--quiz-id", "10", "--answers-json", `[{"id":`},
			ExpectError: true,
		},
		{
			Name:        "update question - missing question ID",
			Args:        []string{"--course-id", "1", "--quiz-id", "10"},
			ExpectError: true,
		},
		{
			Name:        "update question - invalid question ID",
			Args:        []string{"abc", "--course-id", "1", "--quiz-id", "10"},
			ExpectError: true,
		},
		{
			Name:        "update question - missing quiz ID",
			Args:        []string{"789", "--course-id", "1"},
			ExpectError: true,
		},
		{
			Name: "update question - API error",
			Args: []string{"789", "--course-id", "1", "--quiz-id", "10", "--points", "5"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                          courseMock,
				"/api/v1/courses/1/quizzes/10/questions/789": cmdtest.NewErrorResponse(http.StatusNotFound, "not found"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newQuizzesQuestionsUpdateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

// TestRunQuizzesQuestionsUpdate_SendsOnlySetFields captures the PUT body and
// checks that explicitly-set flags are sent (including a zero value) while
// unset flags are omitted so Canvas keeps their current values.
func TestRunQuizzesQuestionsUpdate_SendsOnlySetFields(t *testing.T) {
	var gotMethod string
	var gotBody map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[]`))
			return
		}
		if r.URL.Path != "/api/v1/courses/1/quizzes/10/questions/789" {
			t.Errorf("unexpected path %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		gotMethod = r.Method
		raw, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(raw, &gotBody); err != nil {
			t.Errorf("request body is not JSON: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":789,"quiz_id":10,"question_name":"Q","points_possible":0}`))
	}))
	defer server.Close()

	client, err := api.NewClient(api.ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 100,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	opts := &options.QuizzesQuestionsUpdateOptions{
		CourseID: 1, QuizID: 10, QuestionID: 789,
		PointsPossible: 0, PointsPossibleSet: true, // explicit zero must be sent
		Position: 3, PositionSet: true,
		AnswersJSON:    `[{"id":1001,"text":"4","weight":100},{"id":1002,"text":"5"}]`,
		AnswersJSONSet: true,
		// QuestionText deliberately not set
	}

	out := captureStdout(func() {
		if err := runQuizzesQuestionsUpdate(context.Background(), client, opts); err != nil {
			t.Errorf("runQuizzesQuestionsUpdate: %v", err)
		}
	})

	if gotMethod != http.MethodPut {
		t.Errorf("method = %s, want PUT", gotMethod)
	}
	q, ok := gotBody["question"].(map[string]interface{})
	if !ok {
		t.Fatalf("body has no question object: %v", gotBody)
	}
	if v, ok := q["points_possible"]; !ok || v != float64(0) {
		t.Errorf("points_possible = %v (present=%v), want explicit 0", v, ok)
	}
	if v := q["position"]; v != float64(3) {
		t.Errorf("position = %v, want 3", v)
	}
	if _, present := q["question_text"]; present {
		t.Errorf("question_text was sent although the flag was not set: %v", q)
	}
	answers, ok := q["answers"].([]interface{})
	if !ok || len(answers) != 2 {
		t.Fatalf("answers = %v, want 2 entries", q["answers"])
	}
	first, _ := answers[0].(map[string]interface{})
	if first["id"] != float64(1001) || first["weight"] != float64(100) || first["text"] != "4" {
		t.Errorf("first answer = %v, want id 1001 / weight 100 / text 4", first)
	}
	second, _ := answers[1].(map[string]interface{})
	if w, hasWeight := second["weight"]; !hasWeight || w != float64(0) {
		t.Errorf("second answer must carry an explicit weight 0, got %v", second)
	}
	if !strings.Contains(out, "Question updated successfully (ID: 789)") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestBuildUpdateQuizQuestionParams_InvalidAnswers(t *testing.T) {
	_, err := buildUpdateQuizQuestionParams(&options.QuizzesQuestionsUpdateOptions{
		AnswersJSON: `[{"id":"not-a-number"}]`, AnswersJSONSet: true,
	})
	if err == nil {
		t.Fatal("expected an error for answers JSON that does not match the answer shape")
	}
	if !strings.Contains(err.Error(), "answers-json") {
		t.Errorf("error should name the flag, got: %v", err)
	}
}

// TestQuizzesQuestionsUpdateCmd_DryRun checks that --dry-run prints the PUT
// as a curl command and nothing else (no success line, no empty table).
func TestQuizzesQuestionsUpdateCmd_DryRun(t *testing.T) {
	withGlobalDryRun(t, true)

	tc := cmdtest.CommandTestCase{
		Name: "dry run",
		Args: []string{"789", "--course-id", "1", "--quiz-id", "10", "--points", "5",
			"--answers-json", `[{"id":1001,"weight":100},{"id":1002,"weight":0}]`},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			for _, want := range []string{
				"curl -X PUT",
				"/api/v1/courses/1/quizzes/10/questions/789",
				`"points_possible":5`,
				`{"id":1002,"weight":0}`,
			} {
				if !strings.Contains(output, want) {
					t.Errorf("expected %q in dry-run output, got:\n%s", want, output)
				}
			}
			for _, unwanted := range []string{"updated successfully", "ID: 0"} {
				if strings.Contains(output, unwanted) {
					t.Errorf("dry run must print the curl only, got %q in:\n%s", unwanted, output)
				}
			}
			if lines := strings.Count(strings.TrimSpace(output), "\n") + 1; lines != 6 {
				t.Errorf("expected only the 6-line curl command, got %d lines:\n%s", lines, output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newQuizzesQuestionsUpdateCmd(), tc)
}
