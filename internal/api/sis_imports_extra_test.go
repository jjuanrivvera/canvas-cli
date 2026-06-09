package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSISImportsService_Create_Success(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "sis_test_*.csv")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString("user_id,login_id,first_name,last_name,email,status\n")
	tmpFile.WriteString("sis001,jdoe,John,Doe,jdoe@example.com,active\n")
	tmpFile.Close()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		if r.FormValue("import_type") != "instructure_csv" {
			t.Errorf("expected import_type=instructure_csv, got %q", r.FormValue("import_type"))
		}
		json.NewEncoder(w).Encode(SISImport{ID: 42, WorkflowState: "created"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewSISImportsService(client)
	params := &CreateSISImportParams{
		FilePath:   tmpFile.Name(),
		ImportType: "instructure_csv",
	}
	sisImport, err := svc.Create(context.Background(), 1, params)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if sisImport.ID != 42 {
		t.Errorf("expected ID 42, got %d", sisImport.ID)
	}
}

func TestSISImportsService_Create_AllParams(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "sis_all_*.csv")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString("course_id,short_name,long_name,status\n")
	tmpFile.Close()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		// Verify all optional fields
		if r.FormValue("extension") != "csv" {
			t.Errorf("expected extension=csv, got %q", r.FormValue("extension"))
		}
		if r.FormValue("batch_mode") != "true" {
			t.Errorf("expected batch_mode=true, got %q", r.FormValue("batch_mode"))
		}
		if r.FormValue("batch_mode_term_id") != "5" {
			t.Errorf("expected batch_mode_term_id=5, got %q", r.FormValue("batch_mode_term_id"))
		}
		if r.FormValue("override_sis_stickiness") != "true" {
			t.Errorf("expected override_sis_stickiness=true, got %q", r.FormValue("override_sis_stickiness"))
		}
		if r.FormValue("add_sis_stickiness") != "true" {
			t.Errorf("expected add_sis_stickiness=true, got %q", r.FormValue("add_sis_stickiness"))
		}
		if r.FormValue("clear_sis_stickiness") != "false" {
			t.Errorf("expected clear_sis_stickiness=false, got %q", r.FormValue("clear_sis_stickiness"))
		}
		if r.FormValue("diffing_data_set_identifier") != "fall2024" {
			t.Errorf("expected diffing_data_set_identifier=fall2024, got %q", r.FormValue("diffing_data_set_identifier"))
		}
		if r.FormValue("diffing_remaster_data_set") != "true" {
			t.Errorf("expected diffing_remaster_data_set=true, got %q", r.FormValue("diffing_remaster_data_set"))
		}
		if r.FormValue("change_threshold") != "5" {
			t.Errorf("expected change_threshold=5, got %q", r.FormValue("change_threshold"))
		}
		json.NewEncoder(w).Encode(SISImport{ID: 55, WorkflowState: "created"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewSISImportsService(client)
	batchMode := true
	termID := int64(5)
	overrideStick := true
	addStick := true
	clearStick := false
	diffRemaster := true
	threshold := float64(5)
	params := &CreateSISImportParams{
		FilePath:                 tmpFile.Name(),
		ImportType:               "instructure_csv",
		Extension:                "csv",
		BatchMode:                &batchMode,
		BatchModeTermID:          &termID,
		OverrideSISStickiness:    &overrideStick,
		AddSISStickiness:         &addStick,
		ClearSISStickiness:       &clearStick,
		DiffingDataSetIdentifier: "fall2024",
		DiffingRemasterDataSet:   &diffRemaster,
		ChangeThreshold:          &threshold,
	}
	sisImport, err := svc.Create(context.Background(), 1, params)
	if err != nil {
		t.Fatalf("Create with all params: %v", err)
	}
	if sisImport.ID != 55 {
		t.Errorf("expected ID 55, got %d", sisImport.ID)
	}
}

func TestSISImportsService_Create_FileNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewSISImportsService(client)
	_, err = svc.Create(context.Background(), 1, &CreateSISImportParams{
		FilePath: "/nonexistent/file.csv",
	})
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestSISImportsService_Create_APIError(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "sis_err_*.csv")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"errors":[{"message":"invalid file format"}]}`))
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewSISImportsService(client)
	_, err = svc.Create(context.Background(), 1, &CreateSISImportParams{
		FilePath: tmpFile.Name(),
	})
	if err == nil {
		t.Error("expected error for API error response, got nil")
	}
}

func TestSISImportsService_List_WithAllOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		q := r.URL.Query()
		if q.Get("created_since") != "2024-01-01T00:00:00Z" {
			t.Errorf("expected created_since, got %q", q.Get("created_since"))
		}
		if q.Get("created_before") != "2024-12-31T23:59:59Z" {
			t.Errorf("expected created_before, got %q", q.Get("created_before"))
		}
		if q.Get("page") != "2" {
			t.Errorf("expected page=2, got %q", q.Get("page"))
		}
		if q.Get("per_page") != "25" {
			t.Errorf("expected per_page=25, got %q", q.Get("per_page"))
		}
		response := struct {
			SISImports []SISImport `json:"sis_imports"`
		}{
			SISImports: []SISImport{{ID: 10, WorkflowState: "imported"}},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewSISImportsService(client)
	opts := &ListSISImportsOptions{
		CreatedSince:  "2024-01-01T00:00:00Z",
		CreatedBefore: "2024-12-31T23:59:59Z",
		Page:          2,
		PerPage:       25,
	}
	imports, err := svc.List(context.Background(), 1, opts)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(imports) != 1 {
		t.Errorf("expected 1, got %d", len(imports))
	}
}
