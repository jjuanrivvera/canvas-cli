// Package api - additional branch-coverage tests for functions where the include/options
// path was not previously exercised.
package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// --- DiscussionsService.Get with include ---

func TestDiscussionsService_Get_WithInclude(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/5/discussion_topics/3" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DiscussionTopic{ID: 3, Title: "Test Topic"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "test-token"})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewDiscussionsService(client)
	topic, err := service.Get(context.Background(), 5, 3, []string{"full_topic"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if topic.ID != 3 {
		t.Errorf("expected ID 3, got %d", topic.ID)
	}
}

// --- GroupsService.Get with include ---

func TestGroupsService_Get_WithInclude(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Group{ID: 10, Name: "Group With Include"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "test-token"})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	group, err := service.Get(context.Background(), 10, []string{"permissions", "tabs"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.ID != 10 {
		t.Errorf("expected ID 10, got %d", group.ID)
	}
}

// --- QuizSubmissionsService.Get with include ---

func TestQuizSubmissionsService_Get_WithInclude(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"quiz_submissions": []QuizSubmission{
				{ID: 77, QuizID: 10, UserID: 5},
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "test-token"})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewQuizSubmissionsService(client)
	sub, err := service.Get(context.Background(), 1, 10, 77, []string{"submission", "user"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sub.ID != 77 {
		t.Errorf("expected ID 77, got %d", sub.ID)
	}
}

// --- DiscussionsService.Update full-body coverage ---

func TestDiscussionsService_Update_AllFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DiscussionTopic{ID: 3, Title: "Updated"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "test-token"})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewDiscussionsService(client)
	boolTrue := true
	title := "Updated"
	message := "new message"
	dtype := "threaded"
	lockAt := "2025-06-01T00:00:00Z"
	delayedAt := "2025-01-01T00:00:00Z"
	params := &UpdateDiscussionParams{
		Title:              &title,
		Message:            &message,
		DiscussionType:     &dtype,
		Published:          &boolTrue,
		AllowRating:        &boolTrue,
		RequireInitialPost: &boolTrue,
		PodcastEnabled:     &boolTrue,
		Pinned:             &boolTrue,
		LockAt:             &lockAt,
		DelayedPostAt:      &delayedAt,
	}
	topic, err := service.Update(context.Background(), 1, 3, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if topic.ID != 3 {
		t.Errorf("expected ID 3, got %d", topic.ID)
	}
}

// --- GroupsService.Create full-body coverage ---

func TestGroupsService_Create_AllFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Group{ID: 20, Name: "Full Group"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "test-token"})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	params := &CreateGroupParams{
		Name:           "Full Group",
		Description:    "desc",
		IsPublic:       true,
		JoinLevel:      "parent_context_auto_join",
		StorageQuotaMb: 50,
		SISGroupID:     "SIS-123",
	}
	group, err := service.Create(context.Background(), 7, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.ID != 20 {
		t.Errorf("expected ID 20, got %d", group.ID)
	}
}

// --- ModulesService.Get with include ---

func TestModulesService_Get_WithInclude(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Module{ID: 9, Name: "Module 9"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "test-token"})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewModulesService(client)
	mod, err := service.Get(context.Background(), 1, 9, []string{"items"}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mod.ID != 9 {
		t.Errorf("expected ID 9, got %d", mod.ID)
	}
}

// --- ModulesService.Create full-body coverage ---

func TestModulesService_Create_AllFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Module{ID: 15, Name: "New Module"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "test-token"})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewModulesService(client)
	params := &CreateModuleParams{
		Name:                      "New Module",
		UnlockAt:                  "2025-01-01T00:00:00Z",
		Position:                  1,
		RequireSequentialProgress: true,
		PrerequisiteModuleIDs:     []int64{1, 2},
		PublishFinalGrade:         true,
	}
	mod, err := service.Create(context.Background(), 1, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mod.ID != 15 {
		t.Errorf("expected ID 15, got %d", mod.ID)
	}
}
