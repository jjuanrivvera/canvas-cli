package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRolesService_List_WithOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		q := r.URL.Query()
		if q.Get("state") != "active" {
			t.Errorf("expected state=active, got %q", q.Get("state"))
		}
		if q.Get("show_inherited") != "true" {
			t.Errorf("expected show_inherited=true, got %q", q.Get("show_inherited"))
		}
		if q.Get("page") != "2" {
			t.Errorf("expected page=2, got %q", q.Get("page"))
		}
		if q.Get("per_page") != "25" {
			t.Errorf("expected per_page=25, got %q", q.Get("per_page"))
		}
		roles := []Role{
			{ID: 10, Label: "Teacher", BaseRoleType: "TeacherEnrollment", WorkflowState: "active"},
		}
		json.NewEncoder(w).Encode(roles)
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewRolesService(client)
	opts := &ListRolesOptions{State: "active", ShowInherited: true, Page: 2, PerPage: 25}
	roles, err := svc.List(context.Background(), 1, opts)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(roles) != 1 {
		t.Errorf("expected 1 role, got %d", len(roles))
	}
}

func TestRolesService_Create_WithPermissions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body["label"] != "Custom Admin" {
			t.Errorf("expected label 'Custom Admin', got %v", body["label"])
		}
		if body["base_role_type"] != "AccountMembership" {
			t.Errorf("expected base_role_type AccountMembership, got %v", body["base_role_type"])
		}
		perms, ok := body["permissions"].(map[string]interface{})
		if !ok {
			t.Fatal("expected permissions map")
		}
		if _, ok := perms["manage_courses"]; !ok {
			t.Error("expected manage_courses in permissions")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Role{ID: 50, Label: "Custom Admin", WorkflowState: "active"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewRolesService(client)
	trueVal := true
	falseVal := false
	params := &CreateRoleParams{
		Label:        "Custom Admin",
		BaseRoleType: "AccountMembership",
		Permissions: map[string]PermissionOverride{
			"manage_courses": {
				Enabled:              &trueVal,
				Locked:               &falseVal,
				AppliesToSelf:        &trueVal,
				AppliesToDescendants: &trueVal,
			},
		},
	}
	role, err := svc.Create(context.Background(), 1, params)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if role.Label != "Custom Admin" {
		t.Errorf("expected 'Custom Admin', got %s", role.Label)
	}
}

func TestRolesService_Update_WithPermissions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body["label"] != "Revised Role" {
			t.Errorf("expected label 'Revised Role', got %v", body["label"])
		}
		perms, ok := body["permissions"].(map[string]interface{})
		if !ok {
			t.Fatal("expected permissions map")
		}
		if _, ok := perms["read_course_list"]; !ok {
			t.Error("expected read_course_list in permissions")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Role{ID: 50, Label: "Revised Role", WorkflowState: "active"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewRolesService(client)
	label := "Revised Role"
	enabled := true
	params := &UpdateRoleParams{
		Label: &label,
		Permissions: map[string]PermissionOverride{
			"read_course_list": {Enabled: &enabled},
		},
	}
	role, err := svc.Update(context.Background(), 1, 50, params)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if role.Label != "Revised Role" {
		t.Errorf("expected 'Revised Role', got %s", role.Label)
	}
}

func TestRolesService_Get_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":[{"message":"not found"}]}`))
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewRolesService(client)
	_, err = svc.Get(context.Background(), 1, 9999)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestRolesService_Deactivate_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"errors":[{"message":"forbidden"}]}`))
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewRolesService(client)
	_, err = svc.Deactivate(context.Background(), 1, 50)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestNewRolesService(t *testing.T) {
	client := &Client{}
	svc := NewRolesService(client)
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
	if svc.client != client {
		t.Error("expected client to be set")
	}
}
