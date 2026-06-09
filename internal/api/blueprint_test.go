package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBlueprintService_GetTemplate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/1/blueprint_templates/default" {
			t.Errorf("Expected path /api/v1/courses/1/blueprint_templates/default, got %s", r.URL.Path)
		}

		template := BlueprintTemplate{
			ID:                    1,
			CourseID:              1,
			AssociatedCourseCount: 5,
		}
		json.NewEncoder(w).Encode(template)
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

	service := NewBlueprintService(client)
	template, err := service.GetTemplate(context.Background(), 1, "")
	if err != nil {
		t.Fatalf("GetTemplate failed: %v", err)
	}

	if template.CourseID != 1 {
		t.Errorf("Expected course ID 1, got %d", template.CourseID)
	}

	if template.AssociatedCourseCount != 5 {
		t.Errorf("Expected 5 associated courses, got %d", template.AssociatedCourseCount)
	}
}

func TestBlueprintService_ListAssociatedCourses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/1/blueprint_templates/default/associated_courses" {
			t.Errorf("Expected path /api/v1/courses/1/blueprint_templates/default/associated_courses, got %s", r.URL.Path)
		}

		courses := []AssociatedCourse{
			{ID: 100, Name: "Course A", CourseCode: "COURSE-A"},
			{ID: 101, Name: "Course B", CourseCode: "COURSE-B"},
		}
		json.NewEncoder(w).Encode(courses)
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

	service := NewBlueprintService(client)
	courses, err := service.ListAssociatedCourses(context.Background(), 1, "", nil)
	if err != nil {
		t.Fatalf("ListAssociatedCourses failed: %v", err)
	}

	if len(courses) != 2 {
		t.Errorf("Expected 2 courses, got %d", len(courses))
	}

	if courses[0].Name != "Course A" {
		t.Errorf("Expected name 'Course A', got '%s'", courses[0].Name)
	}
}

func TestBlueprintService_UpdateAssociations(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/1/blueprint_templates/default/update_associations" {
			t.Errorf("Expected path /api/v1/courses/1/blueprint_templates/default/update_associations, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("Expected method PUT, got %s", r.Method)
		}

		json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
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

	service := NewBlueprintService(client)
	params := &UpdateAssociationsParams{
		CourseIDsToAdd: []int64{100, 101},
	}

	err = service.UpdateAssociations(context.Background(), 1, "", params)
	if err != nil {
		t.Fatalf("UpdateAssociations failed: %v", err)
	}
}

func TestBlueprintService_BeginSync(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/1/blueprint_templates/default/migrations" {
			t.Errorf("Expected path /api/v1/courses/1/blueprint_templates/default/migrations, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		migration := BlueprintMigration{
			ID:            456,
			TemplateID:    1,
			WorkflowState: "queued",
			Comment:       "Test sync",
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

	service := NewBlueprintService(client)
	params := &SyncParams{
		Comment: "Test sync",
	}

	migration, err := service.BeginSync(context.Background(), 1, "", params)
	if err != nil {
		t.Fatalf("BeginSync failed: %v", err)
	}

	if migration.Comment != "Test sync" {
		t.Errorf("Expected comment 'Test sync', got '%s'", migration.Comment)
	}
}

func TestBlueprintService_ListUnsyncedChanges(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/1/blueprint_templates/default/unsynced_changes" {
			t.Errorf("Expected path /api/v1/courses/1/blueprint_templates/default/unsynced_changes, got %s", r.URL.Path)
		}

		changes := []UnsyncedChange{
			{AssetID: 1, AssetType: "assignment", AssetName: "Homework 1", ChangeType: "updated"},
			{AssetID: 2, AssetType: "wiki_page", AssetName: "Welcome Page", ChangeType: "created"},
		}
		json.NewEncoder(w).Encode(changes)
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

	service := NewBlueprintService(client)
	changes, err := service.ListUnsyncedChanges(context.Background(), 1, "")
	if err != nil {
		t.Fatalf("ListUnsyncedChanges failed: %v", err)
	}

	if len(changes) != 2 {
		t.Errorf("Expected 2 changes, got %d", len(changes))
	}

	if changes[0].AssetType != "assignment" {
		t.Errorf("Expected asset type 'assignment', got '%s'", changes[0].AssetType)
	}
}

func TestBlueprintService_ListMigrations(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/1/blueprint_templates/default/migrations" {
			t.Errorf("Expected path /api/v1/courses/1/blueprint_templates/default/migrations, got %s", r.URL.Path)
		}

		migrations := []BlueprintMigration{
			{ID: 1, WorkflowState: "completed"},
			{ID: 2, WorkflowState: "running"},
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

	service := NewBlueprintService(client)
	migrations, err := service.ListMigrations(context.Background(), 1, "", nil)
	if err != nil {
		t.Fatalf("ListMigrations failed: %v", err)
	}

	if len(migrations) != 2 {
		t.Errorf("Expected 2 migrations, got %d", len(migrations))
	}
}

func TestBlueprintService_SetRestriction(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/blueprint_templates/default/restrict_item" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["content_type"] != "assignment" {
			t.Errorf("expected content_type=assignment, got %v", body["content_type"])
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewBlueprintService(client)
	restricted := true
	err = svc.SetRestriction(context.Background(), 10, "", &SetRestrictionParams{
		ContentType: "assignment",
		ContentID:   99,
		Restricted:  &restricted,
	})
	if err != nil {
		t.Fatalf("SetRestriction: %v", err)
	}
}

func TestBlueprintService_SetRestriction_WithRestrictions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["restrictions"] == nil {
			t.Error("expected restrictions in body")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewBlueprintService(client)
	err = svc.SetRestriction(context.Background(), 10, "default", &SetRestrictionParams{
		ContentType: "quiz",
		ContentID:   50,
		Restrictions: &BlueprintRestriction{
			Content:           true,
			Points:            true,
			DueDates:          false,
			AvailabilityDates: false,
		},
	})
	if err != nil {
		t.Fatalf("SetRestriction with restrictions: %v", err)
	}
}

func TestBlueprintService_ListMigrations_WithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		q := r.URL.Query()
		if q.Get("per_page") != "5" {
			t.Errorf("expected per_page=5, got %q", q.Get("per_page"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]BlueprintMigration{{ID: 1}})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewBlueprintService(client)
	migrations, err := svc.ListMigrations(context.Background(), 10, "default", &ListMigrationsOptions{Page: 1, PerPage: 5})
	if err != nil {
		t.Fatalf("ListMigrations: %v", err)
	}
	if len(migrations) != 1 {
		t.Errorf("expected 1, got %d", len(migrations))
	}
}

func TestBlueprintService_ListAssociatedCourses_WithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		q := r.URL.Query()
		if q.Get("per_page") != "20" {
			t.Errorf("expected per_page=20, got %q", q.Get("per_page"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]AssociatedCourse{{ID: 5, Name: "Art"}})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewBlueprintService(client)
	courses, err := svc.ListAssociatedCourses(context.Background(), 10, "default", &ListAssociatedCoursesOptions{Page: 1, PerPage: 20})
	if err != nil {
		t.Fatalf("ListAssociatedCourses: %v", err)
	}
	if len(courses) != 1 {
		t.Errorf("expected 1, got %d", len(courses))
	}
}

func TestBlueprintService_GetMigration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/1/blueprint_templates/default/migrations/123" {
			t.Errorf("Expected path /api/v1/courses/1/blueprint_templates/default/migrations/123, got %s", r.URL.Path)
		}

		migration := BlueprintMigration{
			ID:            123,
			WorkflowState: "completed",
			Comment:       "Weekly sync",
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

	service := NewBlueprintService(client)
	migration, err := service.GetMigration(context.Background(), 1, "", 123, nil)
	if err != nil {
		t.Fatalf("GetMigration failed: %v", err)
	}

	if migration.ID != 123 {
		t.Errorf("Expected ID 123, got %d", migration.ID)
	}
}
