package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSectionsService_ListCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/123/sections" {
			t.Errorf("expected /api/v1/courses/123/sections, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Section{
			{ID: 1, Name: "Section 1", CourseID: 123},
			{ID: 2, Name: "Section 2", CourseID: 123},
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

	service := NewSectionsService(client)
	sections, err := service.ListCourse(context.Background(), 123, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(sections) != 2 {
		t.Errorf("expected 2 sections, got %d", len(sections))
	}

	if sections[0].Name != "Section 1" {
		t.Errorf("expected 'Section 1', got %s", sections[0].Name)
	}
}

func TestSectionsService_ListCourse_WithOptions(t *testing.T) {
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
		json.NewEncoder(w).Encode([]Section{
			{ID: 1, Name: "Section 1", CourseID: 123, TotalStudents: 25},
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

	service := NewSectionsService(client)
	opts := &ListSectionsOptions{
		Include: []string{"students", "total_students"},
	}

	sections, err := service.ListCourse(context.Background(), 123, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(sections) != 1 {
		t.Errorf("expected 1 section, got %d", len(sections))
	}

	if sections[0].TotalStudents != 25 {
		t.Errorf("expected 25 students, got %d", sections[0].TotalStudents)
	}
}

func TestSectionsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/sections/456" {
			t.Errorf("expected /api/v1/sections/456, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Section{
			ID:       456,
			Name:     "Test Section",
			CourseID: 123,
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

	service := NewSectionsService(client)
	section, err := service.Get(context.Background(), 456, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if section.ID != 456 {
		t.Errorf("expected ID 456, got %d", section.ID)
	}

	if section.Name != "Test Section" {
		t.Errorf("expected 'Test Section', got %s", section.Name)
	}
}

func TestSectionsService_Create(t *testing.T) {
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

		sectionData, ok := body["course_section"].(map[string]interface{})
		if !ok {
			t.Error("expected course_section in body")
		}

		if sectionData["name"] != "New Section" {
			t.Errorf("expected name 'New Section', got %v", sectionData["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Section{
			ID:       789,
			Name:     "New Section",
			CourseID: 123,
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

	service := NewSectionsService(client)
	params := &CreateSectionParams{
		Name: "New Section",
	}

	section, err := service.Create(context.Background(), 123, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if section.ID != 789 {
		t.Errorf("expected ID 789, got %d", section.ID)
	}
}

func TestSectionsService_Create_WithAllParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		sectionData, ok := body["course_section"].(map[string]interface{})
		if !ok {
			t.Error("expected course_section in body")
		}

		if sectionData["name"] != "Full Section" {
			t.Errorf("expected name 'Full Section', got %v", sectionData["name"])
		}

		if sectionData["sis_section_id"] != "SIS123" {
			t.Errorf("expected sis_section_id 'SIS123', got %v", sectionData["sis_section_id"])
		}

		if sectionData["start_at"] != "2024-01-01T00:00:00Z" {
			t.Errorf("expected start_at '2024-01-01T00:00:00Z', got %v", sectionData["start_at"])
		}

		if sectionData["restrict_enrollments_to_section_dates"] != true {
			t.Errorf("expected restrict_enrollments_to_section_dates true, got %v", sectionData["restrict_enrollments_to_section_dates"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Section{
			ID:           789,
			Name:         "Full Section",
			SISSectionID: "SIS123",
			CourseID:     123,
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

	service := NewSectionsService(client)
	params := &CreateSectionParams{
		Name:                              "Full Section",
		SISSectionID:                      "SIS123",
		StartAt:                           "2024-01-01T00:00:00Z",
		RestrictEnrollmentsToSectionDates: true,
	}

	section, err := service.Create(context.Background(), 123, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if section.SISSectionID != "SIS123" {
		t.Errorf("expected SISSectionID 'SIS123', got %s", section.SISSectionID)
	}
}

func TestSectionsService_Update(t *testing.T) {
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

		sectionData, ok := body["course_section"].(map[string]interface{})
		if !ok {
			t.Error("expected course_section in body")
		}

		if sectionData["name"] != "Updated Section" {
			t.Errorf("expected name 'Updated Section', got %v", sectionData["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Section{
			ID:       456,
			Name:     "Updated Section",
			CourseID: 123,
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

	service := NewSectionsService(client)
	name := "Updated Section"
	params := &UpdateSectionParams{
		Name: &name,
	}

	section, err := service.Update(context.Background(), 456, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if section.Name != "Updated Section" {
		t.Errorf("expected 'Updated Section', got %s", section.Name)
	}
}

func TestSectionsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/sections/456" {
			t.Errorf("expected /api/v1/sections/456, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Section{
			ID:       456,
			Name:     "Deleted Section",
			CourseID: 123,
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

	service := NewSectionsService(client)
	section, err := service.Delete(context.Background(), 456)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if section.ID != 456 {
		t.Errorf("expected ID 456, got %d", section.ID)
	}
}

func TestSectionsService_Crosslist(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/sections/456/crosslist/789" {
			t.Errorf("expected /api/v1/sections/456/crosslist/789, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Section{
			ID:               456,
			Name:             "Crosslisted Section",
			CourseID:         789,
			NonXlistCourseID: ptrInt64(123),
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

	service := NewSectionsService(client)
	section, err := service.Crosslist(context.Background(), 456, 789, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if section.CourseID != 789 {
		t.Errorf("expected CourseID 789, got %d", section.CourseID)
	}
}

func TestSectionsService_Uncrosslist(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/sections/456/crosslist" {
			t.Errorf("expected /api/v1/sections/456/crosslist, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Section{
			ID:       456,
			Name:     "Uncrosslisted Section",
			CourseID: 123,
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

	service := NewSectionsService(client)
	section, err := service.Uncrosslist(context.Background(), 456, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if section.CourseID != 123 {
		t.Errorf("expected CourseID 123, got %d", section.CourseID)
	}
}

func TestNewSectionsService(t *testing.T) {
	client, err := NewClient(ClientConfig{
		BaseURL: "https://canvas.example.com",
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewSectionsService(client)
	if service == nil {
		t.Fatal("expected non-nil service")
	}
	if service.client != client {
		t.Error("expected client to be set")
	}
}

// Helper function
func ptrInt64(i int64) *int64 {
	return &i
}
