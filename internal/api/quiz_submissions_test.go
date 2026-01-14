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
	}
	if service.client != client {
		t.Error("expected client to be set")
	}
}
