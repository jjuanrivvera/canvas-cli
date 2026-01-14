package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGradesService_GetHistory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/123/gradebook_history/days" {
			t.Errorf("expected /api/v1/courses/123/gradebook_history/days, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]GradebookHistoryDay{
			{Date: "2024-01-15", Graders: 5},
			{Date: "2024-01-16", Graders: 3},
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

	service := NewGradesService(client)
	days, err := service.GetHistory(context.Background(), 123, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(days) != 2 {
		t.Errorf("expected 2 days, got %d", len(days))
	}

	if days[0].Date != "2024-01-15" {
		t.Errorf("expected '2024-01-15', got %s", days[0].Date)
	}
}

func TestGradesService_GetFeed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/123/gradebook_history/feed" {
			t.Errorf("expected /api/v1/courses/123/gradebook_history/feed, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]GradebookHistoryEntry{
			{
				ID:             1,
				UserID:         100,
				UserName:       "Student One",
				AssignmentID:   200,
				AssignmentName: "Assignment 1",
				NewGrade:       "A",
				GraderID:       50,
				GraderName:     "Teacher",
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

	service := NewGradesService(client)
	entries, err := service.GetFeed(context.Background(), 123, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}

	if entries[0].UserName != "Student One" {
		t.Errorf("expected 'Student One', got %s", entries[0].UserName)
	}
}

func TestGradesService_ListCustomColumns(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/123/custom_gradebook_columns" {
			t.Errorf("expected /api/v1/courses/123/custom_gradebook_columns, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]CustomGradebookColumn{
			{ID: 1, Title: "Notes", Position: 1, TeacherNotes: true},
			{ID: 2, Title: "Attendance", Position: 2, TeacherNotes: false},
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

	service := NewGradesService(client)
	columns, err := service.ListCustomColumns(context.Background(), 123, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(columns) != 2 {
		t.Errorf("expected 2 columns, got %d", len(columns))
	}

	if columns[0].Title != "Notes" {
		t.Errorf("expected 'Notes', got %s", columns[0].Title)
	}
}

func TestGradesService_GetCustomColumn(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/123/custom_gradebook_columns/456" {
			t.Errorf("expected /api/v1/courses/123/custom_gradebook_columns/456, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(CustomGradebookColumn{
			ID:           456,
			Title:        "Custom Column",
			Position:     1,
			TeacherNotes: false,
			ReadOnly:     false,
			Hidden:       false,
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

	service := NewGradesService(client)
	column, err := service.GetCustomColumn(context.Background(), 123, 456)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if column.ID != 456 {
		t.Errorf("expected ID 456, got %d", column.ID)
	}

	if column.Title != "Custom Column" {
		t.Errorf("expected 'Custom Column', got %s", column.Title)
	}
}

func TestGradesService_CreateCustomColumn(t *testing.T) {
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

		columnData, ok := body["column"].(map[string]interface{})
		if !ok {
			t.Error("expected column in body")
		}

		if columnData["title"] != "New Column" {
			t.Errorf("expected title 'New Column', got %v", columnData["title"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(CustomGradebookColumn{
			ID:       789,
			Title:    "New Column",
			Position: 1,
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

	service := NewGradesService(client)
	params := &CreateCustomColumnParams{
		Title:    "New Column",
		Position: 1,
	}

	column, err := service.CreateCustomColumn(context.Background(), 123, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if column.ID != 789 {
		t.Errorf("expected ID 789, got %d", column.ID)
	}
}

func TestGradesService_UpdateCustomColumn(t *testing.T) {
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

		columnData, ok := body["column"].(map[string]interface{})
		if !ok {
			t.Error("expected column in body")
		}

		if columnData["title"] != "Updated Column" {
			t.Errorf("expected title 'Updated Column', got %v", columnData["title"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(CustomGradebookColumn{
			ID:    456,
			Title: "Updated Column",
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

	service := NewGradesService(client)
	title := "Updated Column"
	params := &UpdateCustomColumnParams{
		Title: &title,
	}

	column, err := service.UpdateCustomColumn(context.Background(), 123, 456, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if column.Title != "Updated Column" {
		t.Errorf("expected 'Updated Column', got %s", column.Title)
	}
}

func TestGradesService_DeleteCustomColumn(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/courses/123/custom_gradebook_columns/456" {
			t.Errorf("expected /api/v1/courses/123/custom_gradebook_columns/456, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(CustomGradebookColumn{
			ID:    456,
			Title: "Deleted Column",
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

	service := NewGradesService(client)
	column, err := service.DeleteCustomColumn(context.Background(), 123, 456)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if column.ID != 456 {
		t.Errorf("expected ID 456, got %d", column.ID)
	}
}

func TestGradesService_GetCustomColumnData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/123/custom_gradebook_columns/456/data" {
			t.Errorf("expected /api/v1/courses/123/custom_gradebook_columns/456/data, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]CustomColumnDatum{
			{UserID: 100, Content: "Note 1"},
			{UserID: 101, Content: "Note 2"},
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

	service := NewGradesService(client)
	data, err := service.GetCustomColumnData(context.Background(), 123, 456)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(data) != 2 {
		t.Errorf("expected 2 data entries, got %d", len(data))
	}

	if data[0].Content != "Note 1" {
		t.Errorf("expected 'Note 1', got %s", data[0].Content)
	}
}

func TestGradesService_SetCustomColumnData(t *testing.T) {
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

		columnData, ok := body["column_data"].(map[string]interface{})
		if !ok {
			t.Error("expected column_data in body")
		}

		if columnData["content"] != "Updated note" {
			t.Errorf("expected content 'Updated note', got %v", columnData["content"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(CustomColumnDatum{
			UserID:  100,
			Content: "Updated note",
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

	service := NewGradesService(client)
	datum, err := service.SetCustomColumnData(context.Background(), 123, 456, 100, "Updated note")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if datum.Content != "Updated note" {
		t.Errorf("expected 'Updated note', got %s", datum.Content)
	}
}

func TestGradesService_BulkUpdateGrades(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/courses/123/submissions/update_grades" {
			t.Errorf("expected /api/v1/courses/123/submissions/update_grades, got %s", r.URL.Path)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if _, ok := body["grade_data"]; !ok {
			t.Error("expected grade_data in body")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGradesService(client)
	grades := []BulkUpdateGrade{
		{StudentID: 100, AssignmentID: 200, Grade: "A"},
		{StudentID: 101, AssignmentID: 200, Grade: "B"},
	}

	err = service.BulkUpdateGrades(context.Background(), 123, grades)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewGradesService(t *testing.T) {
	client, err := NewClient(ClientConfig{
		BaseURL: "https://canvas.example.com",
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGradesService(client)
	if service == nil {
		t.Fatal("expected non-nil service")
	}
	if service.client != client {
		t.Error("expected client to be set")
	}
}
