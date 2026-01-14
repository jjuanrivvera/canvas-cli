package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOutcomesService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/outcomes/123" {
			t.Errorf("expected /api/v1/outcomes/123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Outcome{
			ID:                123,
			Title:             "Test Outcome",
			Description:       "Test description",
			MasteryPoints:     3.0,
			CalculationMethod: "decaying_average",
			CalculationInt:    65,
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

	service := NewOutcomesService(client)
	outcome, err := service.Get(context.Background(), 123)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if outcome.ID != 123 {
		t.Errorf("expected ID 123, got %d", outcome.ID)
	}

	if outcome.Title != "Test Outcome" {
		t.Errorf("expected 'Test Outcome', got %s", outcome.Title)
	}

	if outcome.MasteryPoints != 3.0 {
		t.Errorf("expected mastery points 3.0, got %f", outcome.MasteryPoints)
	}
}

func TestOutcomesService_Update(t *testing.T) {
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

		if body["title"] != "Updated Outcome" {
			t.Errorf("expected title 'Updated Outcome', got %v", body["title"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Outcome{
			ID:    123,
			Title: "Updated Outcome",
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

	service := NewOutcomesService(client)
	title := "Updated Outcome"
	params := &UpdateOutcomeParams{
		Title: &title,
	}

	outcome, err := service.Update(context.Background(), 123, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if outcome.Title != "Updated Outcome" {
		t.Errorf("expected 'Updated Outcome', got %s", outcome.Title)
	}
}

func TestOutcomesService_ListGroupsAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/outcome_groups" {
			t.Errorf("expected /api/v1/accounts/1/outcome_groups, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]OutcomeGroup{
			{ID: 1, Title: "Group 1"},
			{ID: 2, Title: "Group 2"},
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

	service := NewOutcomesService(client)
	groups, err := service.ListGroupsAccount(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(groups))
	}

	if groups[0].Title != "Group 1" {
		t.Errorf("expected 'Group 1', got %s", groups[0].Title)
	}
}

func TestOutcomesService_ListGroupsCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/123/outcome_groups" {
			t.Errorf("expected /api/v1/courses/123/outcome_groups, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]OutcomeGroup{
			{ID: 1, Title: "Course Group 1"},
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

	service := NewOutcomesService(client)
	groups, err := service.ListGroupsCourse(context.Background(), 123, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(groups))
	}
}

func TestOutcomesService_GetGroupAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/outcome_groups/456" {
			t.Errorf("expected /api/v1/accounts/1/outcome_groups/456, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(OutcomeGroup{
			ID:          456,
			Title:       "Test Group",
			Description: "Test description",
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

	service := NewOutcomesService(client)
	group, err := service.GetGroupAccount(context.Background(), 1, 456)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if group.ID != 456 {
		t.Errorf("expected ID 456, got %d", group.ID)
	}

	if group.Title != "Test Group" {
		t.Errorf("expected 'Test Group', got %s", group.Title)
	}
}

func TestOutcomesService_ListOutcomesInGroupAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/outcome_groups/456/outcomes" {
			t.Errorf("expected /api/v1/accounts/1/outcome_groups/456/outcomes, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]OutcomeLink{
			{
				Outcome: &Outcome{ID: 100, Title: "Outcome 1"},
			},
			{
				Outcome: &Outcome{ID: 101, Title: "Outcome 2"},
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

	service := NewOutcomesService(client)
	links, err := service.ListOutcomesInGroupAccount(context.Background(), 1, 456, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(links) != 2 {
		t.Errorf("expected 2 links, got %d", len(links))
	}
}

func TestOutcomesService_LinkOutcomeAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/accounts/1/outcome_groups/456/outcomes/789" {
			t.Errorf("expected /api/v1/accounts/1/outcome_groups/456/outcomes/789, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(OutcomeLink{
			Outcome: &Outcome{ID: 789, Title: "Linked Outcome"},
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

	service := NewOutcomesService(client)
	link, err := service.LinkOutcomeAccount(context.Background(), 1, 456, 789)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if link.Outcome == nil {
		t.Fatal("expected outcome in link")
	}

	if link.Outcome.ID != 789 {
		t.Errorf("expected outcome ID 789, got %d", link.Outcome.ID)
	}
}

func TestOutcomesService_UnlinkOutcomeAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/accounts/1/outcome_groups/456/outcomes/789" {
			t.Errorf("expected /api/v1/accounts/1/outcome_groups/456/outcomes/789, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(OutcomeLink{
			Outcome: &Outcome{ID: 789, Title: "Unlinked Outcome"},
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

	service := NewOutcomesService(client)
	link, err := service.UnlinkOutcomeAccount(context.Background(), 1, 456, 789)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if link.Outcome == nil {
		t.Fatal("expected outcome in link")
	}

	if link.Outcome.ID != 789 {
		t.Errorf("expected outcome ID 789, got %d", link.Outcome.ID)
	}
}

func TestOutcomesService_GetResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/123/outcome_results" {
			t.Errorf("expected /api/v1/courses/123/outcome_results, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(OutcomeResultsResponse{
			OutcomeResults: []OutcomeResult{
				{ID: 1, Score: 4.0, Mastery: true},
				{ID: 2, Score: 2.5, Mastery: false},
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

	service := NewOutcomesService(client)
	response, err := service.GetResults(context.Background(), 123, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(response.OutcomeResults) != 2 {
		t.Errorf("expected 2 results, got %d", len(response.OutcomeResults))
	}

	if !response.OutcomeResults[0].Mastery {
		t.Error("expected first result to have mastery")
	}
}

func TestOutcomesService_GetAlignments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/123/outcome_alignments" {
			t.Errorf("expected /api/v1/courses/123/outcome_alignments, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]OutcomeAlignment{
			{ID: "1", Name: "Assignment 1"},
			{ID: "2", Name: "Quiz 1"},
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

	service := NewOutcomesService(client)
	alignments, err := service.GetAlignments(context.Background(), 123, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(alignments) != 2 {
		t.Errorf("expected 2 alignments, got %d", len(alignments))
	}

	if alignments[0].Name != "Assignment 1" {
		t.Errorf("expected 'Assignment 1', got %s", alignments[0].Name)
	}
}

func TestNewOutcomesService(t *testing.T) {
	client, err := NewClient(ClientConfig{
		BaseURL: "https://canvas.example.com",
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewOutcomesService(client)
	if service == nil {
		t.Fatal("expected non-nil service")
	}
	if service.client != client {
		t.Error("expected client to be set")
	}
}
