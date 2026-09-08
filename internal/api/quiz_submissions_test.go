package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestQuizSubmissionsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/123/quizzes/456/submissions" {
			t.Errorf("expected /api/v1/courses/123/quizzes/456/submissions, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(QuizSubmissionsResponse{
			QuizSubmissions: []QuizSubmission{
				{ID: 1, QuizID: 456, UserID: 100, Score: 85.0, Attempt: 1, WorkflowState: "complete"},
				{ID: 2, QuizID: 456, UserID: 101, Score: 92.0, Attempt: 1, WorkflowState: "complete"},
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewQuizSubmissionsService(client)
	submissions, err := service.List(context.Background(), 123, 456, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(submissions) != 2 {
		t.Errorf("expected 2 submissions, got %d", len(submissions))
	}

	if submissions[0].Score != 85.0 {
		t.Errorf("expected Score 85.0, got %f", submissions[0].Score)
	}
}

func TestQuizSubmissionsService_List_WithInclude(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		includes := r.URL.Query()["include[]"]
		if len(includes) != 2 {
			t.Errorf("expected 2 include params, got %d", len(includes))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(QuizSubmissionsResponse{
			QuizSubmissions: []QuizSubmission{
				{ID: 1, QuizID: 456, UserID: 100, Score: 85.0},
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewQuizSubmissionsService(client)
	opts := &ListQuizSubmissionsOptions{
		Include: []string{"submission", "user"},
	}

	submissions, err := service.List(context.Background(), 123, 456, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(submissions) != 1 {
		t.Errorf("expected 1 submission, got %d", len(submissions))
	}
}

func TestQuizSubmissionsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/123/quizzes/456/submissions/789" {
			t.Errorf("expected /api/v1/courses/123/quizzes/456/submissions/789, got %s", r.URL.Path)
		}

		startedAt := time.Date(2024, 3, 15, 10, 0, 0, 0, time.UTC)
		finishedAt := time.Date(2024, 3, 15, 10, 45, 0, 0, time.UTC)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(QuizSubmissionsResponse{
			QuizSubmissions: []QuizSubmission{
				{
					ID:            789,
					QuizID:        456,
					UserID:        100,
					Score:         88.5,
					Attempt:       1,
					TimeSpent:     2700,
					StartedAt:     &startedAt,
					FinishedAt:    &finishedAt,
					WorkflowState: "complete",
				},
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewQuizSubmissionsService(client)
	submission, err := service.Get(context.Background(), 123, 456, 789, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if submission.ID != 789 {
		t.Errorf("expected ID 789, got %d", submission.ID)
	}

	if submission.Score != 88.5 {
		t.Errorf("expected Score 88.5, got %f", submission.Score)
	}

	if submission.TimeSpent != 2700 {
		t.Errorf("expected TimeSpent 2700, got %d", submission.TimeSpent)
	}
}

func TestQuizSubmissionsService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		submissions, ok := body["quiz_submissions"].([]interface{})
		if !ok {
			t.Error("expected quiz_submissions array in body")
		}

		if len(submissions) != 1 {
			t.Errorf("expected 1 submission in body, got %d", len(submissions))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(QuizSubmissionsResponse{
			QuizSubmissions: []QuizSubmission{
				{
					ID:            789,
					QuizID:        456,
					ExtraTime:     30,
					FudgePoints:   5.0,
					WorkflowState: "complete",
				},
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewQuizSubmissionsService(client)
	extraTime := 30
	fudgePoints := 5.0
	params := &UpdateQuizSubmissionParams{
		ExtraTime:   &extraTime,
		FudgePoints: &fudgePoints,
	}

	submission, err := service.Update(context.Background(), 123, 456, 789, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if submission.ExtraTime != 30 {
		t.Errorf("expected ExtraTime 30, got %d", submission.ExtraTime)
	}

	if submission.FudgePoints != 5.0 {
		t.Errorf("expected FudgePoints 5.0, got %f", submission.FudgePoints)
	}
}

func TestQuizSubmissionsService_Complete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/courses/123/quizzes/456/submissions/789/complete" {
			t.Errorf("expected complete path, got %s", r.URL.Path)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if body["attempt"].(float64) != 1 {
			t.Errorf("expected attempt 1, got %v", body["attempt"])
		}

		if body["validation_token"] != "abc123" {
			t.Errorf("expected validation_token 'abc123', got %v", body["validation_token"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(QuizSubmissionsResponse{
			QuizSubmissions: []QuizSubmission{
				{
					ID:            789,
					WorkflowState: "complete",
				},
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewQuizSubmissionsService(client)
	submission, err := service.Complete(context.Background(), 123, 456, 789, 1, "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if submission.WorkflowState != "complete" {
		t.Errorf("expected WorkflowState 'complete', got %s", submission.WorkflowState)
	}
}

func TestQuizSubmissionsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/1/quizzes/10/submissions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(QuizSubmissionsResponse{
			QuizSubmissions: []QuizSubmission{
				{ID: 99, QuizID: 10, UserID: 5, WorkflowState: "untaken"},
			},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewQuizSubmissionsService(client)

	sub, err := svc.Create(context.Background(), 1, 10, &StartQuizSubmissionParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sub.ID != 99 {
		t.Errorf("expected ID 99, got %d", sub.ID)
	}
	if sub.WorkflowState != "untaken" {
		t.Errorf("expected WorkflowState 'untaken', got %s", sub.WorkflowState)
	}
}

func TestNewQuizSubmissionsService(t *testing.T) {
	client, err := NewClient(ClientConfig{
		BaseURL: "https://canvas.example.com",
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewQuizSubmissionsService(client)
	if service == nil {
		t.Fatal("expected non-nil service")
		return
	}
	if service.client != client {
		t.Error("expected client to be set")
	}
}

// TestQuizSubmissionsService_Update_QuestionScores checks the documented
// "Update student question scores and comments" body shape:
// quiz_submissions[][attempt], [fudge_points], [questions][<qid>][score|comment].
func TestQuizSubmissionsService_Update_QuestionScores(t *testing.T) {
	var got map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/123/quizzes/456/submissions/789" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(QuizSubmissionsResponse{
			QuizSubmissions: []QuizSubmission{{ID: 789, QuizID: 456, Attempt: 2, Score: 7.5, WorkflowState: "complete"}},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "test-token"})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	attempt := 2
	fudge := -1.5
	score := 2.5
	zero := 0.0
	comment := "Regraded: answer key corrected"
	params := &UpdateQuizSubmissionParams{
		Attempt:     &attempt,
		FudgePoints: &fudge,
		Questions: map[int64]QuizSubmissionQuestionScore{
			11: {Score: &score, Comment: &comment},
			12: {Score: &zero},
		},
	}

	sub, err := NewQuizSubmissionsService(client).Update(context.Background(), 123, 456, 789, params)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if sub.ID != 789 || sub.Score != 7.5 {
		t.Errorf("unexpected submission returned: %+v", sub)
	}

	subs, ok := got["quiz_submissions"].([]interface{})
	if !ok || len(subs) != 1 {
		t.Fatalf("expected one quiz_submissions entry, got %v", got["quiz_submissions"])
	}
	entry, _ := subs[0].(map[string]interface{})
	if entry["attempt"] != float64(2) {
		t.Errorf("attempt = %v, want 2", entry["attempt"])
	}
	if entry["fudge_points"] != float64(-1.5) {
		t.Errorf("fudge_points = %v, want -1.5", entry["fudge_points"])
	}
	questions, ok := entry["questions"].(map[string]interface{})
	if !ok {
		t.Fatalf("questions missing or not an object: %v", entry["questions"])
	}
	q11, _ := questions["11"].(map[string]interface{})
	if q11["score"] != float64(2.5) || q11["comment"] != comment {
		t.Errorf("questions[11] = %v, want score 2.5 and comment", q11)
	}
	q12, _ := questions["12"].(map[string]interface{})
	if q12["score"] != float64(0) {
		t.Errorf("questions[12].score = %v, want explicit 0", q12["score"])
	}
	if _, present := q12["comment"]; present {
		t.Errorf("questions[12] should have no comment key, got %v", q12)
	}
}
