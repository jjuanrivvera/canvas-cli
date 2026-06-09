package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// --- OverridesService.BatchUpdate ---

func TestOverridesService_BatchUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/10/assignments/overrides" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]AssignmentOverride{
			{ID: 1, AssignmentID: 100, Title: "Updated Override"},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "test-token"})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewOverridesService(client)
	overrides := []AssignmentOverrideBatchParams{
		{
			AssignmentID: 100,
			StudentIDs:   []int64{1, 2},
			Title:        "Override One",
			DueAt:        "2025-01-15T23:59:59Z",
			UnlockAt:     "2025-01-01T00:00:00Z",
			LockAt:       "2025-01-20T00:00:00Z",
		},
	}
	result, err := service.BatchUpdate(context.Background(), 10, overrides)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 override, got %d", len(result))
	}
	if result[0].Title != "Updated Override" {
		t.Errorf("expected title 'Updated Override', got %s", result[0].Title)
	}
}

// --- PeerReviewsService.ListSections ---

func TestPeerReviewsService_ListSections(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/5/assignments/20/peer_reviews" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]PeerReview{
			{ID: 1, AssessorID: 10, UserID: 20},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "test-token"})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewPeerReviewsService(client)
	reviews, err := service.ListSections(context.Background(), 5, 20, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reviews) != 1 {
		t.Errorf("expected 1 review, got %d", len(reviews))
	}
}

func TestPeerReviewsService_ListSections_WithInclude(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]PeerReview{})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "test-token"})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewPeerReviewsService(client)
	// Passing include values exercises the path-building branch.
	_, err = service.ListSections(context.Background(), 5, 20, []string{"submission_comments", "user"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- PlannerService.GetOverride ---

func TestPlannerService_GetOverride(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/planner/overrides/42" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(PlannerOverride{
			ID:             42,
			PlannableType:  "Assignment",
			PlannableID:    99,
			MarkedComplete: true,
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "test-token"})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewPlannerService(client)
	override, err := service.GetOverride(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if override.ID != 42 {
		t.Errorf("expected override ID 42, got %d", override.ID)
	}
	if override.PlannableType != "Assignment" {
		t.Errorf("expected plannable type 'Assignment', got %s", override.PlannableType)
	}
	if !override.MarkedComplete {
		t.Error("expected MarkedComplete to be true")
	}
}

// --- ExternalToolsService.GetSessionlessLaunchURLForAccount ---

func TestExternalToolsService_GetSessionlessLaunchURLForAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/external_tools/sessionless_launch" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SessionlessLaunchURL{
			ID:   7,
			Name: "Test Tool",
			URL:  "https://lti.example.com/launch?token=abc",
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "test-token"})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewExternalToolsService(client)
	params := &SessionlessLaunchParams{
		LaunchType:   "course_navigation",
		ID:           7,
		URL:          "https://lti.example.com",
		AssignmentID: 99,
	}
	result, err := service.GetSessionlessLaunchURLForAccount(context.Background(), 1, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != 7 {
		t.Errorf("expected ID 7, got %d", result.ID)
	}
	if result.Name != "Test Tool" {
		t.Errorf("expected name 'Test Tool', got %s", result.Name)
	}
}

// --- URLParamsBuilder.Add and Encode ---

func TestURLParamsBuilder_Add(t *testing.T) {
	b := NewURLParamsBuilder()
	b.Add("key", "value1")
	b.Add("key", "value2")
	b.Add("other", "") // empty value should be skipped

	vals := b.Build()
	if got := vals["key"]; len(got) != 2 {
		t.Errorf("expected 2 values for 'key', got %d: %v", len(got), got)
	}
	if other, ok := vals["other"]; ok {
		t.Errorf("empty value should be skipped, but got: %v", other)
	}
}

func TestURLParamsBuilder_Encode(t *testing.T) {
	b := NewURLParamsBuilder()
	b.Set("course_id", "123")
	encoded := b.Encode()
	if encoded == "" {
		t.Error("expected non-empty encoded string")
	}
	if encoded != "course_id=123" {
		t.Errorf("expected 'course_id=123', got %q", encoded)
	}
}

func TestURLParamsBuilder_Encode_Empty(t *testing.T) {
	b := NewURLParamsBuilder()
	if encoded := b.Encode(); encoded != "" {
		t.Errorf("expected empty string for empty builder, got %q", encoded)
	}
}
