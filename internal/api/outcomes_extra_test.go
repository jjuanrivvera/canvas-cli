package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOutcomesService_GetGroupCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/10/outcome_groups/5" {
			t.Errorf("expected /api/v1/courses/10/outcome_groups/5, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(OutcomeGroup{ID: 5, Title: "Course Group"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "test-token"})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewOutcomesService(client)
	group, err := service.GetGroupCourse(context.Background(), 10, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.ID != 5 {
		t.Errorf("expected group ID 5, got %d", group.ID)
	}
	if group.Title != "Course Group" {
		t.Errorf("expected 'Course Group', got %s", group.Title)
	}
}

func TestOutcomesService_ListOutcomesInGroupCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/10/outcome_groups/5/outcomes" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]OutcomeLink{
			{URL: "/outcomes/1", ContextType: "Course"},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "test-token"})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewOutcomesService(client)
	links, err := service.ListOutcomesInGroupCourse(context.Background(), 10, 5, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(links) != 1 {
		t.Errorf("expected 1 link, got %d", len(links))
	}
}

func TestOutcomesService_ListOutcomesInGroupCourse_WithOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]OutcomeLink{})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "test-token"})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewOutcomesService(client)
	opts := &ListOutcomesInGroupOptions{Page: 1, PerPage: 10}
	links, err := service.ListOutcomesInGroupCourse(context.Background(), 10, 5, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// GetAllPages may return nil for an empty result; just ensure no error.
	_ = links
}

func TestOutcomesService_CreateOutcomeAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/outcome_groups/5/outcomes" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body["title"] != "New Outcome" {
			t.Errorf("expected title 'New Outcome', got %v", body["title"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(OutcomeLink{
			URL:         "/accounts/1/outcome_groups/5/outcomes/99",
			ContextType: "Account",
			Outcome:     &Outcome{ID: 99, Title: "New Outcome"},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "test-token"})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewOutcomesService(client)
	params := &CreateOutcomeParams{
		Title:             "New Outcome",
		DisplayName:       "NO",
		Description:       "desc",
		VendorGUID:        "guid-1",
		MasteryPoints:     3.0,
		CalculationMethod: "decaying_average",
		CalculationInt:    65,
		Ratings: []OutcomeRating{
			{Points: 3.0, Description: "Mastery"},
		},
	}
	link, err := service.CreateOutcomeAccount(context.Background(), 1, 5, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if link.Outcome == nil || link.Outcome.ID != 99 {
		t.Errorf("expected outcome ID 99, got %v", link.Outcome)
	}
}

func TestOutcomesService_CreateOutcomeCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/10/outcome_groups/5/outcomes" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(OutcomeLink{
			URL:         "/courses/10/outcome_groups/5/outcomes/88",
			ContextType: "Course",
			Outcome:     &Outcome{ID: 88, Title: "Course Outcome"},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "test-token"})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewOutcomesService(client)
	params := &CreateOutcomeParams{
		Title:             "Course Outcome",
		DisplayName:       "CO",
		Description:       "course desc",
		VendorGUID:        "guid-2",
		MasteryPoints:     4.0,
		CalculationMethod: "n_mastery",
		CalculationInt:    3,
		Ratings: []OutcomeRating{
			{Points: 4.0, Description: "Excellent"},
		},
	}
	link, err := service.CreateOutcomeCourse(context.Background(), 10, 5, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if link.Outcome == nil || link.Outcome.ID != 88 {
		t.Errorf("expected outcome ID 88, got %v", link.Outcome)
	}
}

func TestOutcomesService_LinkOutcomeCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/10/outcome_groups/5/outcomes/88" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(OutcomeLink{
			URL:         "/courses/10/outcome_groups/5/outcomes/88",
			ContextType: "Course",
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "test-token"})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewOutcomesService(client)
	link, err := service.LinkOutcomeCourse(context.Background(), 10, 5, 88)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if link.ContextType != "Course" {
		t.Errorf("expected context type 'Course', got %s", link.ContextType)
	}
}

func TestOutcomesService_UnlinkOutcomeCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/10/outcome_groups/5/outcomes/88" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(OutcomeLink{
			URL:         "/courses/10/outcome_groups/5/outcomes/88",
			ContextType: "Course",
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "test-token"})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewOutcomesService(client)
	link, err := service.UnlinkOutcomeCourse(context.Background(), 10, 5, 88)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if link.ContextType != "Course" {
		t.Errorf("expected context type 'Course', got %s", link.ContextType)
	}
}
