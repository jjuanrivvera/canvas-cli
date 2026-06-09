package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPlannerService_ListItems_WithOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]PlannerItem{
			{PlannableType: "Assignment", PlannableID: 1},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "test-token"})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewPlannerService(client)
	opts := &ListPlannerItemsOptions{
		StartDate:    "2025-01-01",
		EndDate:      "2025-12-31",
		ContextCodes: []string{"course_1", "course_2"},
		Filter:       "all_assignments",
		Page:         1,
		PerPage:      20,
	}
	items, err := service.ListItems(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}
}

func TestPlannerService_ListNotes_WithOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]PlannerNote{
			{ID: 5, Title: "Study note"},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "test-token"})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewPlannerService(client)
	opts := &ListPlannerNotesOptions{
		StartDate:    "2025-01-01",
		EndDate:      "2025-12-31",
		ContextCodes: []string{"course_1"},
		CourseID:     10,
		Page:         1,
		PerPage:      20,
	}
	notes, err := service.ListNotes(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notes) != 1 {
		t.Errorf("expected 1 note, got %d", len(notes))
	}
	if notes[0].Title != "Study note" {
		t.Errorf("expected 'Study note', got %s", notes[0].Title)
	}
}

func TestPlannerService_ListOverrides_WithOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]PlannerOverride{
			{ID: 1, PlannableType: "Assignment"},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "test-token"})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewPlannerService(client)
	opts := &ListOverridesOptions{
		PlannableType: "Assignment",
		PlannableID:   5,
	}
	overrides, err := service.ListOverrides(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(overrides) != 1 {
		t.Errorf("expected 1 override, got %d", len(overrides))
	}
}
