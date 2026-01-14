package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQuizQuestionsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/123/quizzes/456/questions" {
			t.Errorf("expected /api/v1/courses/123/quizzes/456/questions, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]QuizQuestion{
			{ID: 1, QuestionName: "Question 1", QuestionType: "multiple_choice_question", PointsPossible: 10},
			{ID: 2, QuestionName: "Question 2", QuestionType: "true_false_question", PointsPossible: 5},
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

	service := NewQuizQuestionsService(client)
	questions, err := service.List(context.Background(), 123, 456, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(questions) != 2 {
		t.Errorf("expected 2 questions, got %d", len(questions))
	}

	if questions[0].QuestionName != "Question 1" {
		t.Errorf("expected 'Question 1', got %s", questions[0].QuestionName)
	}
}

func TestQuizQuestionsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/123/quizzes/456/questions/789" {
			t.Errorf("expected /api/v1/courses/123/quizzes/456/questions/789, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(QuizQuestion{
			ID:             789,
			QuizID:         456,
			QuestionName:   "Test Question",
			QuestionType:   "multiple_choice_question",
			QuestionText:   "What is 2+2?",
			PointsPossible: 10,
			Answers: []QuizAnswer{
				{ID: 1, Text: "3", Weight: 0},
				{ID: 2, Text: "4", Weight: 100},
				{ID: 3, Text: "5", Weight: 0},
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

	service := NewQuizQuestionsService(client)
	question, err := service.Get(context.Background(), 123, 456, 789)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if question.ID != 789 {
		t.Errorf("expected ID 789, got %d", question.ID)
	}

	if question.QuestionName != "Test Question" {
		t.Errorf("expected 'Test Question', got %s", question.QuestionName)
	}

	if len(question.Answers) != 3 {
		t.Errorf("expected 3 answers, got %d", len(question.Answers))
	}
}

func TestQuizQuestionsService_Create(t *testing.T) {
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

		questionData, ok := body["question"].(map[string]interface{})
		if !ok {
			t.Error("expected question in body")
		}

		if questionData["question_text"] != "What is the capital of France?" {
			t.Errorf("expected question_text, got %v", questionData["question_text"])
		}

		if questionData["question_type"] != "multiple_choice_question" {
			t.Errorf("expected question_type 'multiple_choice_question', got %v", questionData["question_type"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(QuizQuestion{
			ID:             999,
			QuestionName:   "Capital Question",
			QuestionText:   "What is the capital of France?",
			QuestionType:   "multiple_choice_question",
			PointsPossible: 5,
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

	service := NewQuizQuestionsService(client)
	params := &CreateQuizQuestionParams{
		QuestionName:   "Capital Question",
		QuestionText:   "What is the capital of France?",
		QuestionType:   "multiple_choice_question",
		PointsPossible: 5,
		Answers: []QuizAnswer{
			{Text: "London", Weight: 0},
			{Text: "Paris", Weight: 100},
			{Text: "Berlin", Weight: 0},
		},
	}

	question, err := service.Create(context.Background(), 123, 456, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if question.ID != 999 {
		t.Errorf("expected ID 999, got %d", question.ID)
	}
}

func TestQuizQuestionsService_Update(t *testing.T) {
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

		questionData, ok := body["question"].(map[string]interface{})
		if !ok {
			t.Error("expected question in body")
		}

		if questionData["points_possible"].(float64) != 20 {
			t.Errorf("expected points_possible 20, got %v", questionData["points_possible"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(QuizQuestion{
			ID:             789,
			QuestionName:   "Updated Question",
			PointsPossible: 20,
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

	service := NewQuizQuestionsService(client)
	points := 20.0
	params := &UpdateQuizQuestionParams{
		PointsPossible: &points,
	}

	question, err := service.Update(context.Background(), 123, 456, 789, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if question.PointsPossible != 20 {
		t.Errorf("expected PointsPossible 20, got %f", question.PointsPossible)
	}
}

func TestQuizQuestionsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/courses/123/quizzes/456/questions/789" {
			t.Errorf("expected /api/v1/courses/123/quizzes/456/questions/789, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewQuizQuestionsService(client)
	err = service.Delete(context.Background(), 123, 456, 789)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewQuizQuestionsService(t *testing.T) {
	client, err := NewClient(ClientConfig{
		BaseURL: "https://canvas.example.com",
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewQuizQuestionsService(client)
	if service == nil {
		t.Fatal("expected non-nil service")
	}
	if service.client != client {
		t.Error("expected client to be set")
	}
}
