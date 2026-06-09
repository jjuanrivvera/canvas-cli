package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAdminsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/admins" {
			t.Errorf("Expected path /api/v1/accounts/1/admins, got %s", r.URL.Path)
		}

		admins := []Admin{
			{ID: 1, UserID: 123, Role: "AccountAdmin"},
			{ID: 2, UserID: 456, Role: "AccountAdmin"},
		}
		json.NewEncoder(w).Encode(admins)
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

	service := NewAdminsService(client)
	admins, err := service.List(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(admins) != 2 {
		t.Errorf("Expected 2 admins, got %d", len(admins))
	}

	if admins[0].Role != "AccountAdmin" {
		t.Errorf("Expected role 'AccountAdmin', got '%s'", admins[0].Role)
	}
}

func TestAdminsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/admins" {
			t.Errorf("Expected path /api/v1/accounts/1/admins, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		admin := Admin{
			ID:     1,
			UserID: 123,
			Role:   "AccountAdmin",
			Status: "active",
		}
		json.NewEncoder(w).Encode(admin)
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

	service := NewAdminsService(client)
	params := &CreateAdminParams{
		UserID: 123,
		Role:   "AccountAdmin",
	}

	admin, err := service.Create(context.Background(), 1, params)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if admin.UserID != 123 {
		t.Errorf("Expected user ID 123, got %d", admin.UserID)
	}
}

func TestAdminsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/admins/123" {
			t.Errorf("Expected path /api/v1/accounts/1/admins/123, got %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("Expected method DELETE, got %s", r.Method)
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

	service := NewAdminsService(client)
	admin, err := service.Delete(context.Background(), 1, 123, nil)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if admin.UserID != 123 {
		t.Errorf("Expected user ID 123, got %d", admin.UserID)
	}
}

func TestAdminsService_List_WithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		q := r.URL.Query()
		uids := q["user_id[]"]
		if len(uids) != 2 {
			t.Errorf("expected 2 user_id[], got %d", len(uids))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Admin{{ID: 5, UserID: 1, Role: "AccountAdmin"}})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewAdminsService(client)
	opts := &ListAdminsOptions{UserID: []int64{1, 2}, Page: 1, PerPage: 10}
	admins, err := svc.List(context.Background(), 1, opts)
	if err != nil {
		t.Fatalf("List with opts: %v", err)
	}
	if len(admins) != 1 {
		t.Errorf("expected 1, got %d", len(admins))
	}
}

func TestAdminsService_Create_WithAllOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["role_id"] == nil {
			t.Error("expected role_id in body")
		}
		if body["send_confirmation"] == nil {
			t.Error("expected send_confirmation in body")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Admin{ID: 99, UserID: 50, Role: "Custom"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewAdminsService(client)
	sendConf := true
	params := &CreateAdminParams{UserID: 50, Role: "Custom", RoleID: 7, SendConfirmation: &sendConf}
	admin, err := svc.Create(context.Background(), 1, params)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if admin.ID != 99 {
		t.Errorf("expected ID 99, got %d", admin.ID)
	}
}

func TestAdminsService_Delete_WithRoleID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Query().Get("role_id") != "5" {
			t.Errorf("expected role_id=5, got %q", r.URL.Query().Get("role_id"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewAdminsService(client)
	roleID := int64(5)
	admin, err := svc.Delete(context.Background(), 1, 123, &roleID)
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if admin.UserID != 123 {
		t.Errorf("expected UserID 123, got %d", admin.UserID)
	}
}
