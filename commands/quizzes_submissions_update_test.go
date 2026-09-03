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

const quizSubmissionUpdatedMock = `{
	"quiz_submissions": [{
		"id": 789,
		"quiz_id": 10,
		"user_id": 55,
		"attempt": 1,
		"score": 7.5,
		"kept_score": 7.5,
		"fudge_points": 2,
		"workflow_state": "complete"
	}]
}`

func TestQuizzesSubmissionsUpdateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "update fudge points successfully",
			Args: []string{"789", "--course-id", "1", "--quiz-id", "10", "--attempt", "1", "--fudge-points", "2"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                            courseMock,
				"/api/v1/courses/1/quizzes/10/submissions/789": cmdtest.NewMockResponse(quizSubmissionUpdatedMock),
			},
			ExpectError:  false,
			ExpectOutput: "Quiz submission updated successfully (ID: 789, score: 7.5)",
		},
		{
			Name: "update question scores and comments successfully",
			Args: []string{"789", "--course-id", "1", "--quiz-id", "10", "--attempt", "1",
				"--question-score", "1001=2.5", "--question-score", "1002=0",
				"--question-comment", "1001=Partial credit"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                            courseMock,
				"/api/v1/courses/1/quizzes/10/submissions/789": cmdtest.NewMockResponse(quizSubmissionUpdatedMock),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "789") {
					t.Errorf("Expected submission ID in output, got: %s", output)
				}
			},
		},
		{
			Name:        "update - nothing to update",
			Args:        []string{"789", "--course-id", "1", "--quiz-id", "10", "--attempt", "1"},
			ExpectError: true,
		},
		{
			Name:        "update - missing attempt",
			Args:        []string{"789", "--course-id", "1", "--quiz-id", "10", "--fudge-points", "1"},
			ExpectError: true,
		},
		{
			Name:        "update - bad question score",
			Args:        []string{"789", "--course-id", "1", "--quiz-id", "10", "--attempt", "1", "--question-score", "1001=full"},
			ExpectError: true,
		},
		{
			Name:        "update - missing submission ID",
			Args:        []string{"--course-id", "1", "--quiz-id", "10", "--attempt", "1", "--fudge-points", "1"},
			ExpectError: true,
		},
		{
			Name:        "update - invalid submission ID",
			Args:        []string{"abc", "--course-id", "1", "--quiz-id", "10", "--attempt", "1", "--fudge-points", "1"},
			ExpectError: true,
		},
		{
			Name: "update - API error",
			Args: []string{"789", "--course-id", "1", "--quiz-id", "10", "--attempt", "1", "--fudge-points", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                            courseMock,
				"/api/v1/courses/1/quizzes/10/submissions/789": cmdtest.NewErrorResponse(http.StatusForbidden, "unauthorized"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newQuizzesSubmissionsUpdateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

// TestQuizzesSubmissionsUpdateCmd_DryRun checks that --dry-run prints the PUT
// as a curl command with the documented body and does not fail on the
// client's empty mock response.
func TestQuizzesSubmissionsUpdateCmd_DryRun(t *testing.T) {
	withGlobalDryRun(t, true)

	tc := cmdtest.CommandTestCase{
		Name: "dry run",
		Args: []string{"789", "--course-id", "1", "--quiz-id", "10", "--attempt", "2",
			"--fudge-points", "-1", "--question-score", "1001=3"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			for _, want := range []string{
				"curl -X PUT",
				"/api/v1/courses/1/quizzes/10/submissions/789",
				`"attempt":2`,
				`"fudge_points":-1`,
				`"questions":{"1001":{"score":3}}`,
			} {
				if !strings.Contains(output, want) {
					t.Errorf("expected %q in dry-run output, got:\n%s", want, output)
				}
			}
			if strings.Contains(output, "updated successfully") {
				t.Errorf("dry run must not claim success, got:\n%s", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newQuizzesSubmissionsUpdateCmd(), tc)
}

// TestRunQuizzesSubmissionsUpdate_Body captures the PUT body produced from
// the CLI flags and checks the Canvas "Update student question scores and
// comments" shape.
func TestRunQuizzesSubmissionsUpdate_Body(t *testing.T) {
	var gotMethod string
	var gotBody map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[]`))
			return
		}
		if r.URL.Path != "/api/v1/courses/1/quizzes/10/submissions/789" {
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
		_, _ = w.Write([]byte(quizSubmissionUpdatedMock))
	}))
	defer server.Close()

	client, err := api.NewClient(api.ClientConfig{BaseURL: server.URL, Token: "test-token", RequestsPerSec: 100})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	opts := &options.QuizzesSubmissionsUpdateOptions{
		CourseID: 1, QuizID: 10, SubmissionID: 789, Attempt: 1,
		QuestionScores:   []string{"1001=2.5", "1002=0"},
		QuestionComments: []string{"1001=Partial credit", "1003=Comment only"},
		// FudgePoints deliberately not set
	}

	out := captureStdout(func() {
		if err := runQuizzesSubmissionsUpdate(context.Background(), client, opts); err != nil {
			t.Errorf("runQuizzesSubmissionsUpdate: %v", err)
		}
	})

	if gotMethod != http.MethodPut {
		t.Errorf("method = %s, want PUT", gotMethod)
	}
	subs, ok := gotBody["quiz_submissions"].([]interface{})
	if !ok || len(subs) != 1 {
		t.Fatalf("quiz_submissions = %v, want one entry", gotBody["quiz_submissions"])
	}
	entry, _ := subs[0].(map[string]interface{})
	if entry["attempt"] != float64(1) {
		t.Errorf("attempt = %v, want 1", entry["attempt"])
	}
	if _, present := entry["fudge_points"]; present {
		t.Errorf("fudge_points sent although flag not set: %v", entry)
	}
	questions, _ := entry["questions"].(map[string]interface{})
	q1001, _ := questions["1001"].(map[string]interface{})
	if q1001["score"] != float64(2.5) || q1001["comment"] != "Partial credit" {
		t.Errorf("questions[1001] = %v", q1001)
	}
	q1002, _ := questions["1002"].(map[string]interface{})
	if q1002["score"] != float64(0) {
		t.Errorf("questions[1002].score = %v, want explicit 0", q1002["score"])
	}
	if _, present := q1002["comment"]; present {
		t.Errorf("questions[1002] should have no comment, got %v", q1002)
	}
	q1003, _ := questions["1003"].(map[string]interface{})
	if q1003["comment"] != "Comment only" {
		t.Errorf("questions[1003] = %v", q1003)
	}
	if _, present := q1003["score"]; present {
		t.Errorf("questions[1003] should have no score, got %v", q1003)
	}
	if !strings.Contains(out, "Quiz submission updated successfully (ID: 789, score: 7.5)") {
		t.Errorf("unexpected output: %s", out)
	}
}
