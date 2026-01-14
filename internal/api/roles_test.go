package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRolesService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/roles" {
			t.Errorf("Expected path /api/v1/accounts/1/roles, got %s", r.URL.Path)
		}

		roles := []Role{
			{ID: 1, Label: "Admin", BaseRoleType: "AccountMembership", WorkflowState: "active"},
			{ID: 2, Label: "Teacher", BaseRoleType: "TeacherEnrollment", WorkflowState: "active"},
		}
		json.NewEncoder(w).Encode(roles)
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

	service := NewRolesService(client)
	roles, err := service.List(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(roles) != 2 {
		t.Errorf("Expected 2 roles, got %d", len(roles))
	}

	if roles[0].Label != "Admin" {
		t.Errorf("Expected label 'Admin', got '%s'", roles[0].Label)
	}
}

func TestRolesService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/roles/123" {
			t.Errorf("Expected path /api/v1/accounts/1/roles/123, got %s", r.URL.Path)
		}

		role := Role{
			ID:            123,
			Label:         "Custom Teacher",
			BaseRoleType:  "TeacherEnrollment",
			WorkflowState: "active",
		}
		json.NewEncoder(w).Encode(role)
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

	service := NewRolesService(client)
	role, err := service.Get(context.Background(), 1, 123)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if role.ID != 123 {
		t.Errorf("Expected ID 123, got %d", role.ID)
	}

	if role.Label != "Custom Teacher" {
		t.Errorf("Expected label 'Custom Teacher', got '%s'", role.Label)
	}
}

func TestRolesService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/roles" {
			t.Errorf("Expected path /api/v1/accounts/1/roles, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		role := Role{
			ID:            456,
			Label:         "New Role",
			BaseRoleType:  "TeacherEnrollment",
			WorkflowState: "active",
		}
		json.NewEncoder(w).Encode(role)
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

	service := NewRolesService(client)
	params := &CreateRoleParams{
		Label:        "New Role",
		BaseRoleType: "TeacherEnrollment",
	}

	role, err := service.Create(context.Background(), 1, params)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if role.Label != "New Role" {
		t.Errorf("Expected label 'New Role', got '%s'", role.Label)
	}
}

func TestRolesService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/roles/123" {
			t.Errorf("Expected path /api/v1/accounts/1/roles/123, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("Expected method PUT, got %s", r.Method)
		}

		role := Role{
			ID:            123,
			Label:         "Updated Role",
			WorkflowState: "active",
		}
		json.NewEncoder(w).Encode(role)
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

	service := NewRolesService(client)
	label := "Updated Role"
	params := &UpdateRoleParams{
		Label: &label,
	}

	role, err := service.Update(context.Background(), 1, 123, params)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if role.Label != "Updated Role" {
		t.Errorf("Expected label 'Updated Role', got '%s'", role.Label)
	}
}

func TestRolesService_Deactivate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/roles/123" {
			t.Errorf("Expected path /api/v1/accounts/1/roles/123, got %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("Expected method DELETE, got %s", r.Method)
		}

		role := Role{
			ID:            123,
			WorkflowState: "inactive",
		}
		json.NewEncoder(w).Encode(role)
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

	service := NewRolesService(client)
	role, err := service.Deactivate(context.Background(), 1, 123)
	if err != nil {
		t.Fatalf("Deactivate failed: %v", err)
	}

	if role.WorkflowState != "inactive" {
		t.Errorf("Expected state 'inactive', got '%s'", role.WorkflowState)
	}
}

func TestRolesService_Activate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/roles/123/activate" {
			t.Errorf("Expected path /api/v1/accounts/1/roles/123/activate, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		role := Role{
			ID:            123,
			WorkflowState: "active",
		}
		json.NewEncoder(w).Encode(role)
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

	service := NewRolesService(client)
	role, err := service.Activate(context.Background(), 1, 123)
	if err != nil {
		t.Fatalf("Activate failed: %v", err)
	}

	if role.WorkflowState != "active" {
		t.Errorf("Expected state 'active', got '%s'", role.WorkflowState)
	}
}
