package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSISImportsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/sis_imports" {
			t.Errorf("Expected path /api/v1/accounts/1/sis_imports, got %s", r.URL.Path)
		}

		response := struct {
			SISImports []SISImport `json:"sis_imports"`
		}{
			SISImports: []SISImport{
				{ID: 1, WorkflowState: "imported", Progress: 100},
				{ID: 2, WorkflowState: "importing", Progress: 50},
			},
		}
		json.NewEncoder(w).Encode(response)
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

	service := NewSISImportsService(client)
	imports, err := service.List(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(imports) != 2 {
		t.Errorf("Expected 2 imports, got %d", len(imports))
	}

	if imports[0].WorkflowState != "imported" {
		t.Errorf("Expected state 'imported', got '%s'", imports[0].WorkflowState)
	}
}

func TestSISImportsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/sis_imports/123" {
			t.Errorf("Expected path /api/v1/accounts/1/sis_imports/123, got %s", r.URL.Path)
		}

		sisImport := SISImport{
			ID:            123,
			WorkflowState: "imported",
			Progress:      100,
			Data: &SISImportData{
				ImportType: "instructure_csv",
				Counts: &SISImportCounts{
					Users:       100,
					Enrollments: 500,
				},
			},
		}
		json.NewEncoder(w).Encode(sisImport)
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

	service := NewSISImportsService(client)
	sisImport, err := service.Get(context.Background(), 1, 123)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if sisImport.ID != 123 {
		t.Errorf("Expected ID 123, got %d", sisImport.ID)
	}

	if sisImport.Data.Counts.Users != 100 {
		t.Errorf("Expected 100 users, got %d", sisImport.Data.Counts.Users)
	}
}

func TestSISImportsService_Abort(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/sis_imports/123/abort" {
			t.Errorf("Expected path /api/v1/accounts/1/sis_imports/123/abort, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("Expected method PUT, got %s", r.Method)
		}

		sisImport := SISImport{
			ID:            123,
			WorkflowState: "aborted",
		}
		json.NewEncoder(w).Encode(sisImport)
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

	service := NewSISImportsService(client)
	sisImport, err := service.Abort(context.Background(), 1, 123)
	if err != nil {
		t.Fatalf("Abort failed: %v", err)
	}

	if sisImport.WorkflowState != "aborted" {
		t.Errorf("Expected state 'aborted', got '%s'", sisImport.WorkflowState)
	}
}

func TestSISImportsService_ListErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/sis_imports/123/errors" {
			t.Errorf("Expected path /api/v1/accounts/1/sis_imports/123/errors, got %s", r.URL.Path)
		}

		errors := []SISImportError{
			{SISImportID: 123, File: "users.csv", Message: "Invalid email", Row: 5},
			{SISImportID: 123, File: "users.csv", Message: "Missing required field", Row: 10},
		}
		json.NewEncoder(w).Encode(errors)
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

	service := NewSISImportsService(client)
	errors, err := service.ListErrors(context.Background(), 1, 123)
	if err != nil {
		t.Fatalf("ListErrors failed: %v", err)
	}

	if len(errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(errors))
	}

	if errors[0].Row != 5 {
		t.Errorf("Expected row 5, got %d", errors[0].Row)
	}
}

func TestSISImportsService_ListWithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Query().Get("workflow_state") != "imported" {
			t.Errorf("Expected workflow_state=imported, got %s", r.URL.Query().Get("workflow_state"))
		}

		response := struct {
			SISImports []SISImport `json:"sis_imports"`
		}{
			SISImports: []SISImport{
				{ID: 1, WorkflowState: "imported"},
			},
		}
		json.NewEncoder(w).Encode(response)
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

	service := NewSISImportsService(client)
	opts := &ListSISImportsOptions{
		WorkflowState: "imported",
	}

	imports, err := service.List(context.Background(), 1, opts)
	if err != nil {
		t.Fatalf("List with options failed: %v", err)
	}

	if len(imports) != 1 {
		t.Errorf("Expected 1 import, got %d", len(imports))
	}
}

func TestSISImportsService_RestoreStates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/sis_imports/42/restore_states" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SISRestoreProgress{ID: 5, WorkflowState: "queued"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewSISImportsService(client)
	prog, err := svc.RestoreStates(context.Background(), 1, 42, false)
	if err != nil {
		t.Fatalf("RestoreStates: %v", err)
	}
	if prog.ID != 5 {
		t.Errorf("expected ID 5, got %d", prog.ID)
	}
}

func TestSISImportsService_RestoreStates_BatchMode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["batch_mode"] != true {
			t.Errorf("expected batch_mode=true, got %v", body["batch_mode"])
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SISRestoreProgress{ID: 6, WorkflowState: "queued"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewSISImportsService(client)
	prog, err := svc.RestoreStates(context.Background(), 1, 42, true)
	if err != nil {
		t.Fatalf("RestoreStates batch: %v", err)
	}
	if prog.ID != 6 {
		t.Errorf("expected ID 6, got %d", prog.ID)
	}
}
