package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCoursesService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses" {
			t.Errorf("Expected path /api/v1/courses, got %s", r.URL.Path)
		}

		// Check query parameters
		enrollmentType := r.URL.Query().Get("enrollment_type")
		if enrollmentType != "teacher" {
			t.Errorf("Expected enrollment_type 'teacher', got %s", enrollmentType)
		}

		enrollmentState := r.URL.Query().Get("enrollment_state")
		if enrollmentState != "active" {
			t.Errorf("Expected enrollment_state 'active', got %s", enrollmentState)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"id": 1,
				"name": "Course 1",
				"course_code": "CS101",
				"workflow_state": "available"
			},
			{
				"id": 2,
				"name": "Course 2",
				"course_code": "CS102",
				"workflow_state": "available"
			}
		]`))
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

	service := NewCoursesService(client)
	ctx := context.Background()

	opts := &ListCoursesOptions{
		EnrollmentType:  "teacher",
		EnrollmentState: "active",
	}

	courses, err := service.List(ctx, opts)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(courses) != 2 {
		t.Errorf("Expected 2 courses, got %d", len(courses))
	}
	if courses[0].ID != 1 {
		t.Errorf("Expected first course ID 1, got %d", courses[0].ID)
	}
	if courses[1].ID != 2 {
		t.Errorf("Expected second course ID 2, got %d", courses[1].ID)
	}
}

func TestCoursesService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123" {
			t.Errorf("Expected path /api/v1/courses/123, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 123,
			"name": "Introduction to Computer Science",
			"course_code": "CS101",
			"workflow_state": "available",
			"account_id": 1,
			"enrollment_term_id": 5
		}`))
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

	service := NewCoursesService(client)
	ctx := context.Background()

	course, err := service.Get(ctx, 123, nil)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if course.ID != 123 {
		t.Errorf("Expected course ID 123, got %d", course.ID)
	}
	if course.Name != "Introduction to Computer Science" {
		t.Errorf("Expected course name 'Introduction to Computer Science', got %s", course.Name)
	}
	if course.CourseCode != "CS101" {
		t.Errorf("Expected course code 'CS101', got %s", course.CourseCode)
	}
}

func TestCoursesService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/courses" {
			t.Errorf("Expected path /api/v1/accounts/1/courses, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{
			"id": 999,
			"name": "New Course",
			"course_code": "NEW101",
			"workflow_state": "unpublished",
			"account_id": 1
		}`))
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

	service := NewCoursesService(client)
	ctx := context.Background()

	params := &CreateCourseParams{
		AccountID:  1,
		Name:       "New Course",
		CourseCode: "NEW101",
	}

	course, err := service.Create(ctx, params)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if course.ID != 999 {
		t.Errorf("Expected course ID 999, got %d", course.ID)
	}
	if course.Name != "New Course" {
		t.Errorf("Expected course name 'New Course', got %s", course.Name)
	}
	if course.WorkflowState != "unpublished" {
		t.Errorf("Expected workflow state 'unpublished', got %s", course.WorkflowState)
	}
}

func TestCoursesService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123" {
			t.Errorf("Expected path /api/v1/courses/123, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 123,
			"name": "Updated Course Name",
			"course_code": "CS101-UPDATED",
			"workflow_state": "available"
		}`))
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

	service := NewCoursesService(client)
	ctx := context.Background()

	params := &UpdateCourseParams{
		Name:       "Updated Course Name",
		CourseCode: "CS101-UPDATED",
	}

	course, err := service.Update(ctx, 123, params)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if course.ID != 123 {
		t.Errorf("Expected course ID 123, got %d", course.ID)
	}
	if course.Name != "Updated Course Name" {
		t.Errorf("Expected course name 'Updated Course Name', got %s", course.Name)
	}
}

func TestCoursesService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123" {
			t.Errorf("Expected path /api/v1/courses/123, got %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}

		// Check event query parameter
		event := r.URL.Query().Get("event")
		if event != "delete" {
			t.Errorf("Expected event 'delete', got %s", event)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
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

	service := NewCoursesService(client)
	ctx := context.Background()

	err = service.Delete(ctx, 123, "delete")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestNewCoursesService(t *testing.T) {
	client := &Client{}
	service := NewCoursesService(client)

	if service == nil {
		t.Fatal("NewCoursesService returned nil")
	}
	if service.client != client {
		t.Error("NewCoursesService did not set client correctly")
	}
}

func TestCoursesService_Create_WithAllOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/courses" {
			t.Errorf("Expected path /api/v1/accounts/1/courses, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Verify request body contains all parameters
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		course, ok := body["course"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected course object in body")
		}

		// Check that all parameters are present
		expectedFields := []string{
			"name", "course_code", "start_at", "end_at", "license",
			"is_public", "is_public_to_auth_users", "public_syllabus",
			"public_syllabus_to_auth", "public_description",
			"allow_student_wiki_edits", "allow_wiki_comments",
			"allow_student_forum_attachments", "open_enrollment",
			"self_enrollment", "restrict_enrollments_to_course_dates",
			"term_id", "sis_course_id", "integration_id",
			"hide_final_grades", "apply_assignment_group_weights",
			"time_zone", "offer", "enroll_me", "default_view",
			"syllabus_body", "grading_standard_id", "course_format",
		}

		for _, field := range expectedFields {
			if _, exists := course[field]; !exists {
				t.Errorf("Expected field %s in course data", field)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{
			"id": 999,
			"name": "Comprehensive Course",
			"course_code": "COMP101",
			"workflow_state": "available",
			"account_id": 1
		}`))
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

	service := NewCoursesService(client)
	ctx := context.Background()

	params := &CreateCourseParams{
		AccountID:                        1,
		Name:                             "Comprehensive Course",
		CourseCode:                       "COMP101",
		StartAt:                          "2024-01-01T00:00:00Z",
		EndAt:                            "2024-12-31T23:59:59Z",
		License:                          "cc_by_sa",
		IsPublic:                         true,
		IsPublicToAuthUsers:              true,
		PublicSyllabus:                   true,
		PublicSyllabusToAuth:             true,
		PublicDescription:                "A comprehensive test course",
		AllowStudentWikiEdits:            true,
		AllowWikiComments:                true,
		AllowStudentForumAttachments:     true,
		OpenEnrollment:                   true,
		SelfEnrollment:                   true,
		RestrictEnrollmentsToCourseDates: true,
		TermID:                           5,
		SISCourseID:                      "SIS123",
		IntegrationID:                    "INT456",
		HideFinalGrades:                  true,
		ApplyAssignmentGroupWeights:      true,
		TimeZone:                         "America/New_York",
		Offer:                            true,
		EnrollMe:                         true,
		DefaultView:                      "modules",
		SyllabusBody:                     "<p>Course syllabus content</p>",
		GradingStandardID:                10,
		CourseFormat:                     "on_campus",
	}

	course, err := service.Create(ctx, params)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if course.ID != 999 {
		t.Errorf("Expected course ID 999, got %d", course.ID)
	}
	if course.Name != "Comprehensive Course" {
		t.Errorf("Expected course name 'Comprehensive Course', got %s", course.Name)
	}
}

func TestCoursesService_Update_WithAllOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123" {
			t.Errorf("Expected path /api/v1/courses/123, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}

		// Verify request body contains all parameters
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		course, ok := body["course"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected course object in body")
		}

		// Check that all parameters are present
		expectedFields := []string{
			"name", "course_code", "start_at", "end_at", "license",
			"is_public", "is_public_to_auth_users", "public_syllabus",
			"public_syllabus_to_auth", "public_description",
			"allow_student_wiki_edits", "allow_wiki_comments",
			"allow_student_forum_attachments", "open_enrollment",
			"self_enrollment", "restrict_enrollments_to_course_dates",
			"hide_final_grades", "apply_assignment_group_weights",
			"time_zone", "default_view", "syllabus_body",
			"grading_standard_id", "course_format", "image_id",
			"image_url", "remove_image",
		}

		for _, field := range expectedFields {
			if _, exists := course[field]; !exists {
				t.Errorf("Expected field %s in course data", field)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 123,
			"name": "Updated Comprehensive Course",
			"course_code": "COMP101-UPD",
			"workflow_state": "available"
		}`))
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

	service := NewCoursesService(client)
	ctx := context.Background()

	// Use pointer types for all boolean fields
	isPublic := true
	isPublicToAuth := false
	publicSyllabus := true
	publicSyllabusToAuth := false
	allowWikiEdits := true
	allowWikiComments := false
	allowForumAttachments := true
	openEnrollment := false
	selfEnrollment := true
	restrictEnrollments := false
	hideFinalGrades := true
	applyWeights := false

	params := &UpdateCourseParams{
		Name:                             "Updated Comprehensive Course",
		CourseCode:                       "COMP101-UPD",
		StartAt:                          "2024-02-01T00:00:00Z",
		EndAt:                            "2024-11-30T23:59:59Z",
		License:                          "cc_by_nc_sa",
		IsPublic:                         &isPublic,
		IsPublicToAuthUsers:              &isPublicToAuth,
		PublicSyllabus:                   &publicSyllabus,
		PublicSyllabusToAuth:             &publicSyllabusToAuth,
		PublicDescription:                "Updated course description",
		AllowStudentWikiEdits:            &allowWikiEdits,
		AllowWikiComments:                &allowWikiComments,
		AllowStudentForumAttachments:     &allowForumAttachments,
		OpenEnrollment:                   &openEnrollment,
		SelfEnrollment:                   &selfEnrollment,
		RestrictEnrollmentsToCourseDates: &restrictEnrollments,
		HideFinalGrades:                  &hideFinalGrades,
		ApplyAssignmentGroupWeights:      &applyWeights,
		TimeZone:                         "America/Los_Angeles",
		DefaultView:                      "wiki",
		SyllabusBody:                     "<p>Updated syllabus</p>",
		GradingStandardID:                15,
		CourseFormat:                     "online",
		ImageID:                          500,
		ImageURL:                         "https://example.com/image.jpg",
		RemoveImage:                      true,
	}

	course, err := service.Update(ctx, 123, params)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if course.ID != 123 {
		t.Errorf("Expected course ID 123, got %d", course.ID)
	}
	if course.Name != "Updated Comprehensive Course" {
		t.Errorf("Expected course name 'Updated Comprehensive Course', got %s", course.Name)
	}
}

func TestCoursesService_Get_WithIncludes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123" {
			t.Errorf("Expected path /api/v1/courses/123, got %s", r.URL.Path)
		}

		// Check include parameters
		query := r.URL.Query()
		includes := query["include[]"]
		if len(includes) == 0 {
			t.Error("Expected include[] parameters")
		}

		expectedIncludes := map[string]bool{
			"syllabus_body": false,
			"term":          false,
			"account":       false,
			"permissions":   false,
		}

		for _, inc := range includes {
			expectedIncludes[inc] = true
		}

		for inc, found := range expectedIncludes {
			if !found {
				t.Errorf("Expected include parameter %s", inc)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 123,
			"name": "Test Course",
			"course_code": "TEST101",
			"workflow_state": "available",
			"syllabus_body": "<p>Syllabus content</p>",
			"term": {"id": 1, "name": "Fall 2024"},
			"account": {"id": 1, "name": "Test Account"},
			"permissions": {"create_announcement": true}
		}`))
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

	service := NewCoursesService(client)
	ctx := context.Background()

	include := []string{"syllabus_body", "term", "account", "permissions"}
	course, err := service.Get(ctx, 123, include)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if course.ID != 123 {
		t.Errorf("Expected course ID 123, got %d", course.ID)
	}
	if course.Name != "Test Course" {
		t.Errorf("Expected course name 'Test Course', got %s", course.Name)
	}
}

func TestCoursesService_List_WithAllOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		// Check query parameters
		query := r.URL.Query()
		if query.Get("enrollment_type") == "" {
			t.Error("Expected enrollment_type parameter")
		}
		if query.Get("enrollment_state") == "" {
			t.Error("Expected enrollment_state parameter")
		}
		if len(query["include[]"]) == 0 {
			t.Error("Expected include[] parameters")
		}
		if len(query["state[]"]) == 0 {
			t.Error("Expected state[] parameter")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{
			"id": 1,
			"name": "Course 1",
			"course_code": "C1"
		}]`))
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

	service := NewCoursesService(client)
	ctx := context.Background()

	opts := &ListCoursesOptions{
		EnrollmentType:  "teacher",
		EnrollmentState: "active",
		Include:         []string{"term", "account", "total_students"},
		State:           []string{"available", "completed"},
		Page:            1,
		PerPage:         50,
	}

	courses, err := service.List(ctx, opts)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(courses) != 1 {
		t.Errorf("Expected 1 course, got %d", len(courses))
	}
}
