package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPeerReviewsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/1/assignments/10/peer_reviews" {
			t.Errorf("Expected path /api/v1/courses/1/assignments/10/peer_reviews, got %s", r.URL.Path)
		}

		reviews := []PeerReview{
			{ID: 1, AssessorID: 100, UserID: 200, WorkflowState: "assigned"},
			{ID: 2, AssessorID: 101, UserID: 201, WorkflowState: "completed"},
		}
		json.NewEncoder(w).Encode(reviews)
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

	service := NewPeerReviewsService(client)
	reviews, err := service.List(context.Background(), 1, 10, nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(reviews) != 2 {
		t.Errorf("Expected 2 reviews, got %d", len(reviews))
	}

	if reviews[0].WorkflowState != "assigned" {
		t.Errorf("Expected state 'assigned', got '%s'", reviews[0].WorkflowState)
	}
}

func TestPeerReviewsService_ListWithInclude(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/1/assignments/10/peer_reviews" {
			t.Errorf("Expected path /api/v1/courses/1/assignments/10/peer_reviews, got %s", r.URL.Path)
		}

		// Check include params
		includes := r.URL.Query()["include[]"]
		if len(includes) != 2 {
			t.Errorf("Expected 2 include params, got %d", len(includes))
		}

		reviews := []PeerReview{
			{
				ID:            1,
				AssessorID:    100,
				UserID:        200,
				WorkflowState: "assigned",
				User:          &User{ID: 200, Name: "Student A"},
				Assessor:      &User{ID: 100, Name: "Reviewer A"},
			},
		}
		json.NewEncoder(w).Encode(reviews)
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

	service := NewPeerReviewsService(client)
	opts := &ListPeerReviewsOptions{
		Include: []string{"user", "submission_comments"},
	}
	reviews, err := service.List(context.Background(), 1, 10, opts)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(reviews) != 1 {
		t.Errorf("Expected 1 review, got %d", len(reviews))
	}

	if reviews[0].User == nil || reviews[0].User.Name != "Student A" {
		t.Error("Expected user to be included")
	}
}

func TestPeerReviewsService_ListForSubmission(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/1/assignments/10/submissions/500/peer_reviews" {
			t.Errorf("Expected path /api/v1/courses/1/assignments/10/submissions/500/peer_reviews, got %s", r.URL.Path)
		}

		reviews := []PeerReview{
			{ID: 1, AssessorID: 100, UserID: 200, AssetID: 500, AssetType: "Submission"},
		}
		json.NewEncoder(w).Encode(reviews)
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

	service := NewPeerReviewsService(client)
	reviews, err := service.ListForSubmission(context.Background(), 1, 10, 500, nil)
	if err != nil {
		t.Fatalf("ListForSubmission failed: %v", err)
	}

	if len(reviews) != 1 {
		t.Errorf("Expected 1 review, got %d", len(reviews))
	}

	if reviews[0].AssetID != 500 {
		t.Errorf("Expected asset ID 500, got %d", reviews[0].AssetID)
	}
}

func TestPeerReviewsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/1/assignments/10/submissions/500/peer_reviews" {
			t.Errorf("Expected path /api/v1/courses/1/assignments/10/submissions/500/peer_reviews, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		review := PeerReview{
			ID:            789,
			AssessorID:    300,
			UserID:        200,
			WorkflowState: "assigned",
		}
		json.NewEncoder(w).Encode(review)
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

	service := NewPeerReviewsService(client)
	params := &CreatePeerReviewParams{
		UserID: 300,
	}

	review, err := service.Create(context.Background(), 1, 10, 500, params)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if review.AssessorID != 300 {
		t.Errorf("Expected assessor ID 300, got %d", review.AssessorID)
	}
}

func TestPeerReviewsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/1/assignments/10/submissions/500/peer_reviews" {
			t.Errorf("Expected path /api/v1/courses/1/assignments/10/submissions/500/peer_reviews, got %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("Expected method DELETE, got %s", r.Method)
		}

		// Check user_id query param
		userID := r.URL.Query().Get("user_id")
		if userID != "300" {
			t.Errorf("Expected user_id 300, got %s", userID)
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

	service := NewPeerReviewsService(client)
	err = service.Delete(context.Background(), 1, 10, 500, 300)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}
