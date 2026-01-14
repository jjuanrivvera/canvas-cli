package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRubricsService_ListCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/123/rubrics" {
			t.Errorf("expected /api/v1/courses/123/rubrics, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Rubric{
			{ID: 1, Title: "Rubric 1", PointsPossible: 100},
			{ID: 2, Title: "Rubric 2", PointsPossible: 50},
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

	service := NewRubricsService(client)
	rubrics, err := service.ListCourse(context.Background(), 123, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(rubrics) != 2 {
		t.Errorf("expected 2 rubrics, got %d", len(rubrics))
	}

	if rubrics[0].Title != "Rubric 1" {
		t.Errorf("expected 'Rubric 1', got %s", rubrics[0].Title)
	}
}

func TestRubricsService_ListAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/rubrics" {
			t.Errorf("expected /api/v1/accounts/1/rubrics, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Rubric{
			{ID: 1, Title: "Account Rubric", PointsPossible: 100, Reusable: true},
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

	service := NewRubricsService(client)
	rubrics, err := service.ListAccount(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(rubrics) != 1 {
		t.Errorf("expected 1 rubric, got %d", len(rubrics))
	}
}

func TestRubricsService_GetCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/123/rubrics/456" {
			t.Errorf("expected /api/v1/courses/123/rubrics/456, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Rubric{
			ID:             456,
			Title:          "Test Rubric",
			PointsPossible: 100,
			Data: []RubricCriterion{
				{
					ID:          "1",
					Description: "Quality",
					Points:      50,
					Ratings: []RubricRating{
						{ID: "1", Description: "Excellent", Points: 50},
						{ID: "2", Description: "Good", Points: 40},
					},
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

	service := NewRubricsService(client)
	rubric, err := service.GetCourse(context.Background(), 123, 456, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rubric.ID != 456 {
		t.Errorf("expected ID 456, got %d", rubric.ID)
	}

	if rubric.Title != "Test Rubric" {
		t.Errorf("expected 'Test Rubric', got %s", rubric.Title)
	}

	if len(rubric.Data) != 1 {
		t.Errorf("expected 1 criterion, got %d", len(rubric.Data))
	}
}

func TestRubricsService_Create(t *testing.T) {
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

		rubricData, ok := body["rubric"].(map[string]interface{})
		if !ok {
			t.Error("expected rubric in body")
		}

		if rubricData["title"] != "New Rubric" {
			t.Errorf("expected title 'New Rubric', got %v", rubricData["title"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		// Canvas API returns wrapped response
		json.NewEncoder(w).Encode(map[string]interface{}{
			"rubric": Rubric{
				ID:             789,
				Title:          "New Rubric",
				PointsPossible: 100,
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

	service := NewRubricsService(client)
	params := &CreateRubricParams{
		Title:          "New Rubric",
		PointsPossible: 100,
	}

	rubric, err := service.Create(context.Background(), 123, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rubric.ID != 789 {
		t.Errorf("expected ID 789, got %d", rubric.ID)
	}
}

func TestRubricsService_Update(t *testing.T) {
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

		rubricData, ok := body["rubric"].(map[string]interface{})
		if !ok {
			t.Error("expected rubric in body")
		}

		if rubricData["title"] != "Updated Rubric" {
			t.Errorf("expected title 'Updated Rubric', got %v", rubricData["title"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Canvas API returns wrapped response
		json.NewEncoder(w).Encode(map[string]interface{}{
			"rubric": Rubric{
				ID:    456,
				Title: "Updated Rubric",
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

	service := NewRubricsService(client)
	title := "Updated Rubric"
	params := &UpdateRubricParams{
		Title: &title,
	}

	rubric, err := service.Update(context.Background(), 123, 456, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rubric.Title != "Updated Rubric" {
		t.Errorf("expected 'Updated Rubric', got %s", rubric.Title)
	}
}

func TestRubricsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/courses/123/rubrics/456" {
			t.Errorf("expected /api/v1/courses/123/rubrics/456, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Canvas API returns wrapped response
		json.NewEncoder(w).Encode(map[string]interface{}{
			"rubric": Rubric{
				ID:    456,
				Title: "Deleted Rubric",
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

	service := NewRubricsService(client)
	rubric, err := service.Delete(context.Background(), 123, 456)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rubric.ID != 456 {
		t.Errorf("expected ID 456, got %d", rubric.ID)
	}
}

func TestRubricsService_Associate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/courses/123/rubric_associations" {
			t.Errorf("expected /api/v1/courses/123/rubric_associations, got %s", r.URL.Path)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		assocData, ok := body["rubric_association"].(map[string]interface{})
		if !ok {
			t.Error("expected rubric_association in body")
		}

		if assocData["rubric_id"].(float64) != 456 {
			t.Errorf("expected rubric_id 456, got %v", assocData["rubric_id"])
		}

		if assocData["association_id"].(float64) != 789 {
			t.Errorf("expected association_id 789, got %v", assocData["association_id"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		// Canvas API returns wrapped response
		json.NewEncoder(w).Encode(map[string]interface{}{
			"rubric_association": RubricAssociation{
				ID:              999,
				RubricID:        456,
				AssociationID:   789,
				AssociationType: "Assignment",
				UseForGrading:   true,
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

	service := NewRubricsService(client)
	params := &AssociateParams{
		AssociationType: "Assignment",
		AssociationID:   789,
		UseForGrading:   true,
		Purpose:         "grading",
	}

	association, err := service.Associate(context.Background(), 123, 456, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if association.ID != 999 {
		t.Errorf("expected ID 999, got %d", association.ID)
	}

	if !association.UseForGrading {
		t.Error("expected UseForGrading to be true")
	}
}

func TestNewRubricsService(t *testing.T) {
	client, err := NewClient(ClientConfig{
		BaseURL: "https://canvas.example.com",
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewRubricsService(client)
	if service == nil {
		t.Fatal("expected non-nil service")
	}
	if service.client != client {
		t.Error("expected client to be set")
	}
}
