package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestContentMigrationsService_List_WithOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		q := r.URL.Query()
		if q.Get("page") != "2" {
			t.Errorf("expected page=2, got %q", q.Get("page"))
		}
		if q.Get("per_page") != "10" {
			t.Errorf("expected per_page=10, got %q", q.Get("per_page"))
		}
		json.NewEncoder(w).Encode([]ContentMigration{
			{ID: 1, MigrationType: "course_copy_importer", WorkflowState: "running"},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewContentMigrationsService(client)
	opts := &ListContentMigrationsOptions{Page: 2, PerPage: 10}
	migrations, err := svc.List(context.Background(), 5, opts)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(migrations) != 1 {
		t.Errorf("expected 1, got %d", len(migrations))
	}
}

func TestContentMigrationsService_List_NilOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		json.NewEncoder(w).Encode([]ContentMigration{})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewContentMigrationsService(client)
	_, err = svc.List(context.Background(), 5, nil)
	if err != nil {
		t.Fatalf("List nil opts: %v", err)
	}
}

func TestContentMigrationsService_Create_WithFileURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body["migration_type"] != "common_cartridge_importer" {
			t.Errorf("expected migration_type common_cartridge_importer, got %v", body["migration_type"])
		}
		settings, ok := body["settings"].(map[string]interface{})
		if !ok {
			t.Fatal("expected settings map")
		}
		if settings["file_url"] != "https://example.com/export.imscc" {
			t.Errorf("expected file_url, got %v", settings["file_url"])
		}
		json.NewEncoder(w).Encode(ContentMigration{ID: 42, MigrationType: "common_cartridge_importer"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewContentMigrationsService(client)
	params := &CreateContentMigrationParams{
		MigrationType: "common_cartridge_importer",
		FileURL:       "https://example.com/export.imscc",
	}
	migration, err := svc.Create(context.Background(), 5, params)
	if err != nil {
		t.Fatalf("Create with FileURL: %v", err)
	}
	if migration.ID != 42 {
		t.Errorf("expected ID 42, got %d", migration.ID)
	}
}

func TestContentMigrationsService_Create_WithAllSettings(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		settings, ok := body["settings"].(map[string]interface{})
		if !ok {
			t.Fatal("expected settings map")
		}
		if settings["question_bank_name"] != "My Bank" {
			t.Errorf("expected question_bank_name 'My Bank', got %v", settings["question_bank_name"])
		}
		if settings["overwrite_quizzes"] != true {
			t.Errorf("expected overwrite_quizzes=true")
		}
		if body["selective_import"] != true {
			t.Errorf("expected selective_import=true")
		}
		if body["copy"] == nil {
			t.Error("expected copy options in body")
		}
		if body["date_shift_options"] == nil {
			t.Error("expected date_shift_options in body")
		}
		json.NewEncoder(w).Encode(ContentMigration{ID: 10, MigrationType: "qti_converter"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewContentMigrationsService(client)
	srcCourseID := int64(100)
	qbID := int64(55)
	folderID := int64(20)
	overwrite := true
	selective := true
	contentExportID := int64(77)
	params := &CreateContentMigrationParams{
		MigrationType:    "qti_converter",
		SourceCourseID:   &srcCourseID,
		ContentExportID:  &contentExportID,
		QuestionBankID:   &qbID,
		QuestionBankName: "My Bank",
		FolderID:         &folderID,
		OverwriteQuizzes: &overwrite,
		SelectiveImport:  &selective,
		CopyOptions:      map[string]interface{}{"all_quizzes": true},
		DateShiftOptions: &DateShiftOptions{
			ShiftDates:   true,
			NewStartDate: "2025-01-01",
			NewEndDate:   "2025-06-30",
		},
	}
	migration, err := svc.Create(context.Background(), 5, params)
	if err != nil {
		t.Fatalf("Create with all settings: %v", err)
	}
	if migration.ID != 10 {
		t.Errorf("expected ID 10, got %d", migration.ID)
	}
}

func TestContentMigrationsService_CreateWithFile_Success(t *testing.T) {
	// Create a temporary file to simulate a course export
	tmpFile, err := os.CreateTemp("", "canvas_test_*.imscc")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString("fake imscc content for testing")
	tmpFile.Close()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		if r.FormValue("migration_type") != "common_cartridge_importer" {
			t.Errorf("expected migration_type, got %q", r.FormValue("migration_type"))
		}
		json.NewEncoder(w).Encode(ContentMigration{ID: 99, MigrationType: "common_cartridge_importer"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewContentMigrationsService(client)
	params := &CreateContentMigrationParams{
		MigrationType: "common_cartridge_importer",
		FilePath:      tmpFile.Name(),
	}
	migration, err := svc.Create(context.Background(), 5, params)
	if err != nil {
		t.Fatalf("Create with file: %v", err)
	}
	if migration.ID != 99 {
		t.Errorf("expected ID 99, got %d", migration.ID)
	}
}

func TestContentMigrationsService_CreateWithFile_WithOptionalFields(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "canvas_test_opts_*.imscc")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString("test content")
	tmpFile.Close()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		if r.FormValue("settings[source_course_id]") != "123" {
			t.Errorf("expected source_course_id=123, got %q", r.FormValue("settings[source_course_id]"))
		}
		if r.FormValue("settings[folder_id]") != "456" {
			t.Errorf("expected folder_id=456, got %q", r.FormValue("settings[folder_id]"))
		}
		if r.FormValue("selective_import") != "true" {
			t.Errorf("expected selective_import=true, got %q", r.FormValue("selective_import"))
		}
		json.NewEncoder(w).Encode(ContentMigration{ID: 77, MigrationType: "course_copy_importer"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewContentMigrationsService(client)
	srcID := int64(123)
	folderID := int64(456)
	selective := true
	params := &CreateContentMigrationParams{
		MigrationType:   "course_copy_importer",
		FilePath:        tmpFile.Name(),
		SourceCourseID:  &srcID,
		FolderID:        &folderID,
		SelectiveImport: &selective,
	}
	migration, err := svc.Create(context.Background(), 5, params)
	if err != nil {
		t.Fatalf("Create with file and optional fields: %v", err)
	}
	if migration.ID != 77 {
		t.Errorf("expected ID 77, got %d", migration.ID)
	}
}

func TestContentMigrationsService_CreateWithFile_FileNotFound(t *testing.T) {
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
	svc := NewContentMigrationsService(client)
	params := &CreateContentMigrationParams{
		MigrationType: "common_cartridge_importer",
		FilePath:      "/nonexistent/path/to/file.imscc",
	}
	_, err = svc.Create(context.Background(), 5, params)
	if err == nil {
		t.Error("expected error for nonexistent file, got nil")
	}
}

func TestContentMigrationsService_CreateWithFile_APIError(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "canvas_err_test_*.imscc")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString("test content")
	tmpFile.Close()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"errors":[{"message":"bad request"}]}`))
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewContentMigrationsService(client)
	params := &CreateContentMigrationParams{
		MigrationType: "common_cartridge_importer",
		FilePath:      tmpFile.Name(),
	}
	_, err = svc.Create(context.Background(), 5, params)
	if err == nil {
		t.Error("expected error for API error response, got nil")
	}
}

func TestContentMigrationsService_ListContentList_WithType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Query().Get("type") != "assignments" {
			t.Errorf("expected type=assignments, got %q", r.URL.Query().Get("type"))
		}
		items := []ContentListItem{
			{Type: "assignment", Property: "copy[all_assignments]", Title: "Homework 1"},
		}
		json.NewEncoder(w).Encode(items)
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewContentMigrationsService(client)
	items, err := svc.ListContentList(context.Background(), 5, 42, "assignments")
	if err != nil {
		t.Fatalf("ListContentList: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("expected 1, got %d", len(items))
	}
}

func TestContentMigrationsService_ListContentList_NoType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Query().Get("type") != "" {
			t.Errorf("expected no type param, got %q", r.URL.Query().Get("type"))
		}
		json.NewEncoder(w).Encode([]ContentListItem{})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewContentMigrationsService(client)
	_, err = svc.ListContentList(context.Background(), 5, 42, "")
	if err != nil {
		t.Fatalf("ListContentList with no type: %v", err)
	}
}
