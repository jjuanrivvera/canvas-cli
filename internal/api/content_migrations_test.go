package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContentMigrationsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/1/content_migrations" {
			t.Errorf("Expected path /api/v1/courses/1/content_migrations, got %s", r.URL.Path)
		}

		migrations := []ContentMigration{
			{ID: 1, MigrationType: "course_copy_importer", WorkflowState: "completed"},
			{ID: 2, MigrationType: "common_cartridge_importer", WorkflowState: "running"},
		}
		json.NewEncoder(w).Encode(migrations)
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

	service := NewContentMigrationsService(client)
	migrations, err := service.List(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(migrations) != 2 {
		t.Errorf("Expected 2 migrations, got %d", len(migrations))
	}

	if migrations[0].MigrationType != "course_copy_importer" {
		t.Errorf("Expected type 'course_copy_importer', got '%s'", migrations[0].MigrationType)
	}
}

func TestContentMigrationsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/1/content_migrations/123" {
			t.Errorf("Expected path /api/v1/courses/1/content_migrations/123, got %s", r.URL.Path)
		}

		migration := ContentMigration{
			ID:                 123,
			MigrationType:      "course_copy_importer",
			MigrationTypeTitle: "Course Copy",
			WorkflowState:      "completed",
		}
		json.NewEncoder(w).Encode(migration)
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

	service := NewContentMigrationsService(client)
	migration, err := service.Get(context.Background(), 1, 123)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if migration.ID != 123 {
		t.Errorf("Expected ID 123, got %d", migration.ID)
	}

	if migration.MigrationTypeTitle != "Course Copy" {
		t.Errorf("Expected title 'Course Copy', got '%s'", migration.MigrationTypeTitle)
	}
}

func TestContentMigrationsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/1/content_migrations" {
			t.Errorf("Expected path /api/v1/courses/1/content_migrations, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		migration := ContentMigration{
			ID:            456,
			MigrationType: "course_copy_importer",
			WorkflowState: "queued",
		}
		json.NewEncoder(w).Encode(migration)
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

	service := NewContentMigrationsService(client)
	sourceCourseID := int64(100)
	params := &CreateContentMigrationParams{
		MigrationType:  "course_copy_importer",
		SourceCourseID: &sourceCourseID,
	}

	migration, err := service.Create(context.Background(), 1, params)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if migration.MigrationType != "course_copy_importer" {
		t.Errorf("Expected type 'course_copy_importer', got '%s'", migration.MigrationType)
	}
}

func TestContentMigrationsService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/1/content_migrations/123" {
			t.Errorf("Expected path /api/v1/courses/1/content_migrations/123, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("Expected method PUT, got %s", r.Method)
		}

		migration := ContentMigration{
			ID:            123,
			WorkflowState: "aborted",
		}
		json.NewEncoder(w).Encode(migration)
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

	service := NewContentMigrationsService(client)
	params := &UpdateContentMigrationParams{
		WorkflowState: "aborted",
	}

	migration, err := service.Update(context.Background(), 1, 123, params)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if migration.WorkflowState != "aborted" {
		t.Errorf("Expected state 'aborted', got '%s'", migration.WorkflowState)
	}
}

func TestContentMigrationsService_ListMigrators(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/1/content_migrations/migrators" {
			t.Errorf("Expected path /api/v1/courses/1/content_migrations/migrators, got %s", r.URL.Path)
		}

		migrators := []Migrator{
			{Type: "course_copy_importer", Name: "Course Copy", RequiresFileUpload: false},
			{Type: "common_cartridge_importer", Name: "Common Cartridge", RequiresFileUpload: true},
		}
		json.NewEncoder(w).Encode(migrators)
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

	service := NewContentMigrationsService(client)
	migrators, err := service.ListMigrators(context.Background(), 1)
	if err != nil {
		t.Fatalf("ListMigrators failed: %v", err)
	}

	if len(migrators) != 2 {
		t.Errorf("Expected 2 migrators, got %d", len(migrators))
	}

	if migrators[0].Type != "course_copy_importer" {
		t.Errorf("Expected type 'course_copy_importer', got '%s'", migrators[0].Type)
	}
}

func TestContentMigrationsService_ListContentList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/1/content_migrations/123/content_list" {
			t.Errorf("Expected path /api/v1/courses/1/content_migrations/123/content_list, got %s", r.URL.Path)
		}

		items := []ContentListItem{
			{Type: "assignments", Property: "copy[assignments]", Title: "Assignments"},
			{Type: "modules", Property: "copy[modules]", Title: "Modules"},
		}
		json.NewEncoder(w).Encode(items)
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

	service := NewContentMigrationsService(client)
	items, err := service.ListContentList(context.Background(), 1, 123, "")
	if err != nil {
		t.Fatalf("ListContentList failed: %v", err)
	}

	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}

	if items[0].Type != "assignments" {
		t.Errorf("Expected type 'assignments', got '%s'", items[0].Type)
	}
}

func TestContentMigrationsService_ListMigrationIssues(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/1/content_migrations/123/migration_issues" {
			t.Errorf("Expected path /api/v1/courses/1/content_migrations/123/migration_issues, got %s", r.URL.Path)
		}

		issues := []MigrationIssue{
			{ID: 1, ContentMigrationID: 123, Description: "Missing file", IssueType: "warning"},
			{ID: 2, ContentMigrationID: 123, Description: "Invalid link", IssueType: "error"},
		}
		json.NewEncoder(w).Encode(issues)
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

	service := NewContentMigrationsService(client)
	issues, err := service.ListMigrationIssues(context.Background(), 1, 123)
	if err != nil {
		t.Fatalf("ListMigrationIssues failed: %v", err)
	}

	if len(issues) != 2 {
		t.Errorf("Expected 2 issues, got %d", len(issues))
	}

	if issues[0].IssueType != "warning" {
		t.Errorf("Expected type 'warning', got '%s'", issues[0].IssueType)
	}
}
