package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestOverridesService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/123/assignments/456/overrides" {
			t.Errorf("expected /api/v1/courses/123/assignments/456/overrides, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]AssignmentOverride{
			{ID: 1, AssignmentID: 456, Title: "Section A Override", CourseSectionID: 100},
			{ID: 2, AssignmentID: 456, Title: "Student Override", StudentIDs: []int64{200, 201}},
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

	service := NewOverridesService(client)
	overrides, err := service.List(context.Background(), 123, 456, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(overrides) != 2 {
		t.Errorf("expected 2 overrides, got %d", len(overrides))
	}

	if overrides[0].Title != "Section A Override" {
		t.Errorf("expected 'Section A Override', got %s", overrides[0].Title)
	}
}

func TestOverridesService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/123/assignments/456/overrides/789" {
			t.Errorf("expected /api/v1/courses/123/assignments/456/overrides/789, got %s", r.URL.Path)
		}

		dueAt := time.Date(2024, 3, 15, 23, 59, 0, 0, time.UTC)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(AssignmentOverride{
			ID:           789,
			AssignmentID: 456,
			Title:        "Test Override",
			DueAt:        &dueAt,
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

	service := NewOverridesService(client)
	override, err := service.Get(context.Background(), 123, 456, 789)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if override.ID != 789 {
		t.Errorf("expected ID 789, got %d", override.ID)
	}

	if override.Title != "Test Override" {
		t.Errorf("expected 'Test Override', got %s", override.Title)
	}
}

func TestOverridesService_Create(t *testing.T) {
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

		overrideData, ok := body["assignment_override"].(map[string]interface{})
		if !ok {
			t.Error("expected assignment_override in body")
		}

		if overrideData["title"] != "New Override" {
			t.Errorf("expected title 'New Override', got %v", overrideData["title"])
		}

		if overrideData["course_section_id"].(float64) != 100 {
			t.Errorf("expected course_section_id 100, got %v", overrideData["course_section_id"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(AssignmentOverride{
			ID:              999,
			AssignmentID:    456,
			Title:           "New Override",
			CourseSectionID: 100,
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

	service := NewOverridesService(client)
	params := &AssignmentOverrideCreateParams{
		Title:           "New Override",
		CourseSectionID: 100,
	}

	override, err := service.Create(context.Background(), 123, 456, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if override.ID != 999 {
		t.Errorf("expected ID 999, got %d", override.ID)
	}
}

func TestOverridesService_Create_WithStudents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		overrideData, ok := body["assignment_override"].(map[string]interface{})
		if !ok {
			t.Error("expected assignment_override in body")
		}

		studentIDs, ok := overrideData["student_ids"].([]interface{})
		if !ok {
			t.Error("expected student_ids in body")
		}

		if len(studentIDs) != 2 {
			t.Errorf("expected 2 student IDs, got %d", len(studentIDs))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(AssignmentOverride{
			ID:           999,
			AssignmentID: 456,
			Title:        "Student Override",
			StudentIDs:   []int64{200, 201},
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

	service := NewOverridesService(client)
	params := &AssignmentOverrideCreateParams{
		Title:      "Student Override",
		StudentIDs: []int64{200, 201},
		DueAt:      "2024-03-15T23:59:00Z",
	}

	override, err := service.Create(context.Background(), 123, 456, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(override.StudentIDs) != 2 {
		t.Errorf("expected 2 student IDs, got %d", len(override.StudentIDs))
	}
}

func TestOverridesService_Update(t *testing.T) {
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

		overrideData, ok := body["assignment_override"].(map[string]interface{})
		if !ok {
			t.Error("expected assignment_override in body")
		}

		if overrideData["title"] != "Updated Override" {
			t.Errorf("expected title 'Updated Override', got %v", overrideData["title"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(AssignmentOverride{
			ID:           789,
			AssignmentID: 456,
			Title:        "Updated Override",
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

	service := NewOverridesService(client)
	title := "Updated Override"
	params := &AssignmentOverrideUpdateParams{
		Title: &title,
	}

	override, err := service.Update(context.Background(), 123, 456, 789, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if override.Title != "Updated Override" {
		t.Errorf("expected 'Updated Override', got %s", override.Title)
	}
}

func TestOverridesService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/courses/123/assignments/456/overrides/789" {
			t.Errorf("expected /api/v1/courses/123/assignments/456/overrides/789, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(AssignmentOverride{
			ID:           789,
			AssignmentID: 456,
			Title:        "Deleted Override",
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

	service := NewOverridesService(client)
	override, err := service.Delete(context.Background(), 123, 456, 789)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if override.ID != 789 {
		t.Errorf("expected ID 789, got %d", override.ID)
	}
}

func TestOverridesService_BatchCreate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/courses/123/assignments/overrides" {
			t.Errorf("expected /api/v1/courses/123/assignments/overrides, got %s", r.URL.Path)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		overrides, ok := body["assignment_overrides"].([]interface{})
		if !ok {
			t.Error("expected assignment_overrides array in body")
		}

		if len(overrides) != 2 {
			t.Errorf("expected 2 overrides, got %d", len(overrides))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode([]AssignmentOverride{
			{ID: 1, AssignmentID: 100, Title: "Override 1"},
			{ID: 2, AssignmentID: 101, Title: "Override 2"},
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

	service := NewOverridesService(client)
	params := []AssignmentOverrideBatchParams{
		{AssignmentID: 100, Title: "Override 1", SectionID: 50},
		{AssignmentID: 101, Title: "Override 2", SectionID: 50},
	}

	overrides, err := service.BatchCreate(context.Background(), 123, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(overrides) != 2 {
		t.Errorf("expected 2 overrides, got %d", len(overrides))
	}
}

func TestNewOverridesService(t *testing.T) {
	client, err := NewClient(ClientConfig{
		BaseURL: "https://canvas.example.com",
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewOverridesService(client)
	if service == nil {
		t.Fatal("expected non-nil service")
	}
	if service.client != client {
		t.Error("expected client to be set")
	}
}
