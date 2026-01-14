package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAssignmentGroupsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/123/assignment_groups" {
			t.Errorf("expected /api/v1/courses/123/assignment_groups, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]AssignmentGroup{
			{ID: 1, Name: "Group 1", Position: 1, GroupWeight: 50.0},
			{ID: 2, Name: "Group 2", Position: 2, GroupWeight: 50.0},
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

	service := NewAssignmentGroupsService(client)
	groups, err := service.List(context.Background(), 123, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(groups))
	}

	if groups[0].Name != "Group 1" {
		t.Errorf("expected 'Group 1', got %s", groups[0].Name)
	}
}

func TestAssignmentGroupsService_List_WithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		// Check include parameters
		includes := r.URL.Query()["include[]"]
		if len(includes) != 2 {
			t.Errorf("expected 2 include params, got %d", len(includes))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]AssignmentGroup{
			{ID: 1, Name: "Group 1", Position: 1},
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

	service := NewAssignmentGroupsService(client)
	opts := &ListAssignmentGroupsOptions{
		Include: []string{"assignments", "rules"},
	}

	groups, err := service.List(context.Background(), 123, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(groups))
	}
}

func TestAssignmentGroupsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/123/assignment_groups/456" {
			t.Errorf("expected /api/v1/courses/123/assignment_groups/456, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(AssignmentGroup{
			ID:          456,
			Name:        "Test Group",
			Position:    1,
			GroupWeight: 100.0,
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

	service := NewAssignmentGroupsService(client)
	group, err := service.Get(context.Background(), 123, 456, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if group.ID != 456 {
		t.Errorf("expected ID 456, got %d", group.ID)
	}

	if group.Name != "Test Group" {
		t.Errorf("expected 'Test Group', got %s", group.Name)
	}
}

func TestAssignmentGroupsService_Create(t *testing.T) {
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

		if body["name"] != "New Group" {
			t.Errorf("expected name 'New Group', got %v", body["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(AssignmentGroup{
			ID:       789,
			Name:     "New Group",
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

	service := NewAssignmentGroupsService(client)
	params := &CreateAssignmentGroupParams{
		Name: "New Group",
	}

	group, err := service.Create(context.Background(), 123, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if group.ID != 789 {
		t.Errorf("expected ID 789, got %d", group.ID)
	}
}

func TestAssignmentGroupsService_Create_WithRules(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if body["name"] != "Graded Group" {
			t.Errorf("expected name 'Graded Group', got %v", body["name"])
		}

		if body["group_weight"].(float64) != 25.0 {
			t.Errorf("expected group_weight 25.0, got %v", body["group_weight"])
		}

		rules, ok := body["rules"].(map[string]interface{})
		if !ok {
			t.Error("expected rules in body")
		}

		if rules["drop_lowest"].(float64) != 1 {
			t.Errorf("expected drop_lowest 1, got %v", rules["drop_lowest"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(AssignmentGroup{
			ID:          789,
			Name:        "Graded Group",
			GroupWeight: 25.0,
			Rules: &GradingRules{
				DropLowest: 1,
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

	service := NewAssignmentGroupsService(client)
	params := &CreateAssignmentGroupParams{
		Name:        "Graded Group",
		GroupWeight: 25.0,
		Rules: &GradingRules{
			DropLowest: 1,
		},
	}

	group, err := service.Create(context.Background(), 123, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if group.GroupWeight != 25.0 {
		t.Errorf("expected GroupWeight 25.0, got %f", group.GroupWeight)
	}

	if group.Rules == nil || group.Rules.DropLowest != 1 {
		t.Error("expected rules with drop_lowest 1")
	}
}

func TestAssignmentGroupsService_Update(t *testing.T) {
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

		if body["name"] != "Updated Group" {
			t.Errorf("expected name 'Updated Group', got %v", body["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(AssignmentGroup{
			ID:   456,
			Name: "Updated Group",
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

	service := NewAssignmentGroupsService(client)
	name := "Updated Group"
	params := &UpdateAssignmentGroupParams{
		Name: &name,
	}

	group, err := service.Update(context.Background(), 123, 456, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if group.Name != "Updated Group" {
		t.Errorf("expected 'Updated Group', got %s", group.Name)
	}
}

func TestAssignmentGroupsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/courses/123/assignment_groups/456" {
			t.Errorf("expected /api/v1/courses/123/assignment_groups/456, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(AssignmentGroup{
			ID:   456,
			Name: "Deleted Group",
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

	service := NewAssignmentGroupsService(client)
	group, err := service.Delete(context.Background(), 123, 456, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if group.ID != 456 {
		t.Errorf("expected ID 456, got %d", group.ID)
	}
}

func TestAssignmentGroupsService_Delete_WithMoveAssignments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		moveAssignmentsTo := r.URL.Query().Get("move_assignments_to")
		if moveAssignmentsTo != "789" {
			t.Errorf("expected move_assignments_to '789', got %s", moveAssignmentsTo)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(AssignmentGroup{
			ID:   456,
			Name: "Deleted Group",
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

	service := NewAssignmentGroupsService(client)
	opts := &DeleteAssignmentGroupOptions{
		MoveAssignmentsTo: 789,
	}

	group, err := service.Delete(context.Background(), 123, 456, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if group.ID != 456 {
		t.Errorf("expected ID 456, got %d", group.ID)
	}
}

func TestNewAssignmentGroupsService(t *testing.T) {
	client, err := NewClient(ClientConfig{
		BaseURL: "https://canvas.example.com",
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewAssignmentGroupsService(client)
	if service == nil {
		t.Fatal("expected non-nil service")
	}
	if service.client != client {
		t.Error("expected client to be set")
	}
}
