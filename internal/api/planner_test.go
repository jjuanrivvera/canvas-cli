package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPlannerService_ListItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/planner/items" {
			t.Errorf("Expected path /api/v1/planner/items, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"plannable_type": "Assignment",
				"plannable_id": 123,
				"context_type": "Course",
				"context_name": "Physics 101"
			},
			{
				"plannable_type": "Quiz",
				"plannable_id": 456,
				"context_type": "Course",
				"context_name": "Chemistry 101"
			}
		]`))
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 10,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	service := NewPlannerService(client)
	ctx := context.Background()

	items, err := service.ListItems(ctx, nil)
	if err != nil {
		t.Fatalf("ListItems failed: %v", err)
	}

	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}
	if items[0].PlannableType != "Assignment" {
		t.Errorf("Expected first item type 'Assignment', got %s", items[0].PlannableType)
	}
}

func TestPlannerService_ListNotes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/planner_notes" {
			t.Errorf("Expected path /api/v1/planner_notes, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"id": 1,
				"title": "Study for Exam",
				"details": "Review chapters 1-5",
				"workflow_state": "active"
			}
		]`))
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 10,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	service := NewPlannerService(client)
	ctx := context.Background()

	notes, err := service.ListNotes(ctx, nil)
	if err != nil {
		t.Fatalf("ListNotes failed: %v", err)
	}

	if len(notes) != 1 {
		t.Errorf("Expected 1 note, got %d", len(notes))
	}
	if notes[0].Title != "Study for Exam" {
		t.Errorf("Expected title 'Study for Exam', got %s", notes[0].Title)
	}
}

func TestPlannerService_GetNote(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/planner_notes/1" {
			t.Errorf("Expected path /api/v1/planner_notes/1, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 1,
			"title": "Study for Exam",
			"details": "Review chapters 1-5",
			"workflow_state": "active"
		}`))
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 10,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	service := NewPlannerService(client)
	ctx := context.Background()

	note, err := service.GetNote(ctx, 1)
	if err != nil {
		t.Fatalf("GetNote failed: %v", err)
	}

	if note.Title != "Study for Exam" {
		t.Errorf("Expected title 'Study for Exam', got %s", note.Title)
	}
}

func TestPlannerService_CreateNote(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/planner_notes" {
			t.Errorf("Expected path /api/v1/planner_notes, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if body["title"] != "New Note" {
			t.Errorf("Expected title 'New Note', got %v", body["title"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 2,
			"title": "New Note",
			"details": "Note details",
			"workflow_state": "active"
		}`))
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 10,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	service := NewPlannerService(client)
	ctx := context.Background()

	params := &CreateNoteParams{
		Title:   "New Note",
		Details: "Note details",
	}

	note, err := service.CreateNote(ctx, params)
	if err != nil {
		t.Fatalf("CreateNote failed: %v", err)
	}

	if note.Title != "New Note" {
		t.Errorf("Expected title 'New Note', got %s", note.Title)
	}
}

func TestPlannerService_UpdateNote(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/planner_notes/1" {
			t.Errorf("Expected path /api/v1/planner_notes/1, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 1,
			"title": "Updated Note",
			"workflow_state": "active"
		}`))
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 10,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	service := NewPlannerService(client)
	ctx := context.Background()

	title := "Updated Note"
	params := &UpdateNoteParams{
		Title: &title,
	}

	note, err := service.UpdateNote(ctx, 1, params)
	if err != nil {
		t.Fatalf("UpdateNote failed: %v", err)
	}

	if note.Title != "Updated Note" {
		t.Errorf("Expected title 'Updated Note', got %s", note.Title)
	}
}

func TestPlannerService_DeleteNote(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/planner_notes/1" {
			t.Errorf("Expected path /api/v1/planner_notes/1, got %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 10,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	service := NewPlannerService(client)
	ctx := context.Background()

	err = service.DeleteNote(ctx, 1)
	if err != nil {
		t.Fatalf("DeleteNote failed: %v", err)
	}
}

func TestPlannerService_ListOverrides(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/planner/overrides" {
			t.Errorf("Expected path /api/v1/planner/overrides, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"id": 1,
				"plannable_type": "Assignment",
				"plannable_id": 123,
				"marked_complete": true,
				"dismissed": false
			}
		]`))
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 10,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	service := NewPlannerService(client)
	ctx := context.Background()

	overrides, err := service.ListOverrides(ctx, nil)
	if err != nil {
		t.Fatalf("ListOverrides failed: %v", err)
	}

	if len(overrides) != 1 {
		t.Errorf("Expected 1 override, got %d", len(overrides))
	}
	if !overrides[0].MarkedComplete {
		t.Error("Expected marked_complete to be true")
	}
}

func TestPlannerService_CreateOverride(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/planner/overrides" {
			t.Errorf("Expected path /api/v1/planner/overrides, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if body["plannable_type"] != "Assignment" {
			t.Errorf("Expected plannable_type 'Assignment', got %v", body["plannable_type"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 2,
			"plannable_type": "Assignment",
			"plannable_id": 456,
			"marked_complete": true,
			"dismissed": false
		}`))
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 10,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	service := NewPlannerService(client)
	ctx := context.Background()

	params := &CreateOverrideParams{
		PlannableType:  "Assignment",
		PlannableID:    456,
		MarkedComplete: true,
	}

	override, err := service.CreateOverride(ctx, params)
	if err != nil {
		t.Fatalf("CreateOverride failed: %v", err)
	}

	if !override.MarkedComplete {
		t.Error("Expected marked_complete to be true")
	}
}

func TestPlannerService_UpdateOverride(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/planner/overrides/1" {
			t.Errorf("Expected path /api/v1/planner/overrides/1, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 1,
			"plannable_type": "Assignment",
			"plannable_id": 123,
			"marked_complete": false,
			"dismissed": true
		}`))
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 10,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	service := NewPlannerService(client)
	ctx := context.Background()

	dismissed := true
	params := &UpdateOverrideParams{
		Dismissed: &dismissed,
	}

	override, err := service.UpdateOverride(ctx, 1, params)
	if err != nil {
		t.Fatalf("UpdateOverride failed: %v", err)
	}

	if !override.Dismissed {
		t.Error("Expected dismissed to be true")
	}
}

func TestPlannerService_DeleteOverride(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/planner/overrides/1" {
			t.Errorf("Expected path /api/v1/planner/overrides/1, got %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 10,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	service := NewPlannerService(client)
	ctx := context.Background()

	err = service.DeleteOverride(ctx, 1)
	if err != nil {
		t.Fatalf("DeleteOverride failed: %v", err)
	}
}

func TestNewPlannerService(t *testing.T) {
	client := &Client{}
	service := NewPlannerService(client)
	if service == nil {
		t.Fatal("Expected non-nil service")
	}
	if service.client != client {
		t.Error("Expected service client to match input client")
	}
}
