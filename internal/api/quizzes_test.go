package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestQuizzesService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/123/quizzes" {
			t.Errorf("expected /api/v1/courses/123/quizzes, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Quiz{
			{ID: 1, Title: "Quiz 1", QuizType: "assignment", Published: true},
			{ID: 2, Title: "Quiz 2", QuizType: "practice_quiz", Published: false},
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

	service := NewQuizzesService(client)
	quizzes, err := service.List(context.Background(), 123, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(quizzes) != 2 {
		t.Errorf("expected 2 quizzes, got %d", len(quizzes))
	}

	if quizzes[0].Title != "Quiz 1" {
		t.Errorf("expected 'Quiz 1', got %s", quizzes[0].Title)
	}
}

func TestQuizzesService_List_WithSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		searchTerm := r.URL.Query().Get("search_term")
		if searchTerm != "midterm" {
			t.Errorf("expected search_term 'midterm', got %s", searchTerm)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Quiz{
			{ID: 1, Title: "Midterm Exam", QuizType: "assignment"},
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

	service := NewQuizzesService(client)
	opts := &ListQuizzesOptions{
		SearchTerm: "midterm",
	}

	quizzes, err := service.List(context.Background(), 123, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(quizzes) != 1 {
		t.Errorf("expected 1 quiz, got %d", len(quizzes))
	}
}

func TestQuizzesService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/123/quizzes/456" {
			t.Errorf("expected /api/v1/courses/123/quizzes/456, got %s", r.URL.Path)
		}

		dueAt := time.Date(2024, 3, 15, 23, 59, 0, 0, time.UTC)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Quiz{
			ID:              456,
			Title:           "Test Quiz",
			QuizType:        "assignment",
			TimeLimit:       60,
			AllowedAttempts: 2,
			DueAt:           &dueAt,
			Published:       true,
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

	service := NewQuizzesService(client)
	quiz, err := service.Get(context.Background(), 123, 456)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if quiz.ID != 456 {
		t.Errorf("expected ID 456, got %d", quiz.ID)
	}

	if quiz.Title != "Test Quiz" {
		t.Errorf("expected 'Test Quiz', got %s", quiz.Title)
	}

	if quiz.TimeLimit != 60 {
		t.Errorf("expected time_limit 60, got %d", quiz.TimeLimit)
	}
}

func TestQuizzesService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		quizData, ok := body["quiz"].(map[string]interface{})
		if !ok {
			t.Error("expected quiz in body")
		}

		if quizData["title"] != "New Quiz" {
			t.Errorf("expected title 'New Quiz', got %v", quizData["title"])
		}

		if quizData["quiz_type"] != "assignment" {
			t.Errorf("expected quiz_type 'assignment', got %v", quizData["quiz_type"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Quiz{
			ID:       789,
			Title:    "New Quiz",
			QuizType: "assignment",
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

	service := NewQuizzesService(client)
	params := &CreateQuizParams{
		Title:    "New Quiz",
		QuizType: "assignment",
	}

	quiz, err := service.Create(context.Background(), 123, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if quiz.ID != 789 {
		t.Errorf("expected ID 789, got %d", quiz.ID)
	}
}

func TestQuizzesService_Create_WithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		quizData, ok := body["quiz"].(map[string]interface{})
		if !ok {
			t.Error("expected quiz in body")
		}

		if quizData["time_limit"].(float64) != 30 {
			t.Errorf("expected time_limit 30, got %v", quizData["time_limit"])
		}

		if quizData["shuffle_answers"] != true {
			t.Errorf("expected shuffle_answers true, got %v", quizData["shuffle_answers"])
		}

		if quizData["allowed_attempts"].(float64) != 3 {
			t.Errorf("expected allowed_attempts 3, got %v", quizData["allowed_attempts"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Quiz{
			ID:              789,
			Title:           "Timed Quiz",
			TimeLimit:       30,
			ShuffleAnswers:  true,
			AllowedAttempts: 3,
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

	service := NewQuizzesService(client)
	params := &CreateQuizParams{
		Title:           "Timed Quiz",
		QuizType:        "assignment",
		TimeLimit:       30,
		ShuffleAnswers:  true,
		AllowedAttempts: 3,
	}

	quiz, err := service.Create(context.Background(), 123, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if quiz.TimeLimit != 30 {
		t.Errorf("expected TimeLimit 30, got %d", quiz.TimeLimit)
	}
}

func TestQuizzesService_Update(t *testing.T) {
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

		quizData, ok := body["quiz"].(map[string]interface{})
		if !ok {
			t.Error("expected quiz in body")
		}

		if quizData["title"] != "Updated Quiz" {
			t.Errorf("expected title 'Updated Quiz', got %v", quizData["title"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Quiz{
			ID:    456,
			Title: "Updated Quiz",
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

	service := NewQuizzesService(client)
	title := "Updated Quiz"
	params := &UpdateQuizParams{
		Title: &title,
	}

	quiz, err := service.Update(context.Background(), 123, 456, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if quiz.Title != "Updated Quiz" {
		t.Errorf("expected 'Updated Quiz', got %s", quiz.Title)
	}
}

func TestQuizzesService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/courses/123/quizzes/456" {
			t.Errorf("expected /api/v1/courses/123/quizzes/456, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Quiz{
			ID:    456,
			Title: "Deleted Quiz",
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

	service := NewQuizzesService(client)
	quiz, err := service.Delete(context.Background(), 123, 456)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if quiz.ID != 456 {
		t.Errorf("expected ID 456, got %d", quiz.ID)
	}
}

func TestNewQuizzesService(t *testing.T) {
	client, err := NewClient(ClientConfig{
		BaseURL: "https://canvas.example.com",
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewQuizzesService(client)
	if service == nil {
		t.Fatal("expected non-nil service")
	}
	if service.client != client {
		t.Error("expected client to be set")
	}
}
