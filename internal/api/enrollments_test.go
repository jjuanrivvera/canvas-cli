package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEnrollmentsService_ListCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/enrollments" {
			t.Errorf("Expected path /api/v1/courses/123/enrollments, got %s", r.URL.Path)
		}

		// Check query parameters
		enrollmentType := r.URL.Query()["type[]"]
		if len(enrollmentType) != 1 || enrollmentType[0] != "StudentEnrollment" {
			t.Errorf("Expected type[] 'StudentEnrollment', got %v", enrollmentType)
		}

		state := r.URL.Query()["state[]"]
		if len(state) != 1 || state[0] != "active" {
			t.Errorf("Expected state[] 'active', got %v", state)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"id": 1,
				"user_id": 100,
				"course_id": 123,
				"type": "StudentEnrollment",
				"enrollment_state": "active",
				"role": "StudentEnrollment",
				"role_id": 3
			},
			{
				"id": 2,
				"user_id": 101,
				"course_id": 123,
				"type": "StudentEnrollment",
				"enrollment_state": "active",
				"role": "StudentEnrollment",
				"role_id": 3
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

	service := NewEnrollmentsService(client)
	ctx := context.Background()

	opts := &ListEnrollmentsOptions{
		Type:  []string{"StudentEnrollment"},
		State: []string{"active"},
	}

	enrollments, err := service.ListCourse(ctx, 123, opts)
	if err != nil {
		t.Fatalf("ListCourse failed: %v", err)
	}

	if len(enrollments) != 2 {
		t.Errorf("Expected 2 enrollments, got %d", len(enrollments))
	}
	if enrollments[0].UserID != 100 {
		t.Errorf("Expected first enrollment user ID 100, got %d", enrollments[0].UserID)
	}
	if enrollments[1].UserID != 101 {
		t.Errorf("Expected second enrollment user ID 101, got %d", enrollments[1].UserID)
	}
}

func TestEnrollmentsService_ListSection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/sections/456/enrollments" {
			t.Errorf("Expected path /api/v1/sections/456/enrollments, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"id": 1,
				"user_id": 100,
				"course_section_id": 456,
				"type": "StudentEnrollment",
				"enrollment_state": "active"
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

	service := NewEnrollmentsService(client)
	ctx := context.Background()

	enrollments, err := service.ListSection(ctx, 456, nil)
	if err != nil {
		t.Fatalf("ListSection failed: %v", err)
	}

	if len(enrollments) != 1 {
		t.Errorf("Expected 1 enrollment, got %d", len(enrollments))
	}
	if enrollments[0].CourseSectionID != 456 {
		t.Errorf("Expected section ID 456, got %d", enrollments[0].CourseSectionID)
	}
}

func TestEnrollmentsService_ListUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/users/789/enrollments" {
			t.Errorf("Expected path /api/v1/users/789/enrollments, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"id": 1,
				"user_id": 789,
				"course_id": 100,
				"type": "StudentEnrollment",
				"enrollment_state": "active"
			},
			{
				"id": 2,
				"user_id": 789,
				"course_id": 101,
				"type": "StudentEnrollment",
				"enrollment_state": "active"
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

	service := NewEnrollmentsService(client)
	ctx := context.Background()

	enrollments, err := service.ListUser(ctx, 789, nil)
	if err != nil {
		t.Fatalf("ListUser failed: %v", err)
	}

	if len(enrollments) != 2 {
		t.Errorf("Expected 2 enrollments, got %d", len(enrollments))
	}
	if enrollments[0].UserID != 789 {
		t.Errorf("Expected user ID 789, got %d", enrollments[0].UserID)
	}
}

func TestEnrollmentsService_EnrollUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/enrollments" {
			t.Errorf("Expected path /api/v1/courses/123/enrollments, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{
			"id": 999,
			"user_id": 789,
			"course_id": 123,
			"type": "StudentEnrollment",
			"enrollment_state": "active",
			"role": "StudentEnrollment",
			"role_id": 3
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

	service := NewEnrollmentsService(client)
	ctx := context.Background()

	params := &EnrollUserParams{
		UserID:         789,
		Type:           "StudentEnrollment",
		EnrollmentState: "active",
	}

	enrollment, err := service.EnrollUser(ctx, 123, params)
	if err != nil {
		t.Fatalf("EnrollUser failed: %v", err)
	}

	if enrollment.ID != 999 {
		t.Errorf("Expected enrollment ID 999, got %d", enrollment.ID)
	}
	if enrollment.UserID != 789 {
		t.Errorf("Expected user ID 789, got %d", enrollment.UserID)
	}
	if enrollment.Type != "StudentEnrollment" {
		t.Errorf("Expected type 'StudentEnrollment', got %s", enrollment.Type)
	}
}

func TestEnrollmentsService_Conclude(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/enrollments/456" {
			t.Errorf("Expected path /api/v1/courses/123/enrollments/456, got %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}

		// Check task query parameter
		task := r.URL.Query().Get("task")
		if task != "conclude" {
			t.Errorf("Expected task 'conclude', got %s", task)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 456,
			"user_id": 789,
			"course_id": 123,
			"type": "StudentEnrollment",
			"enrollment_state": "completed"
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

	service := NewEnrollmentsService(client)
	ctx := context.Background()

	enrollment, err := service.Conclude(ctx, 123, 456, "conclude")
	if err != nil {
		t.Fatalf("Conclude failed: %v", err)
	}

	if enrollment.ID != 456 {
		t.Errorf("Expected enrollment ID 456, got %d", enrollment.ID)
	}
	if enrollment.EnrollmentState != "completed" {
		t.Errorf("Expected state 'completed', got %s", enrollment.EnrollmentState)
	}
}

func TestEnrollmentsService_Reactivate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/enrollments/456/reactivate" {
			t.Errorf("Expected path /api/v1/courses/123/enrollments/456/reactivate, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 456,
			"user_id": 789,
			"course_id": 123,
			"type": "StudentEnrollment",
			"enrollment_state": "active"
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

	service := NewEnrollmentsService(client)
	ctx := context.Background()

	enrollment, err := service.Reactivate(ctx, 123, 456)
	if err != nil {
		t.Fatalf("Reactivate failed: %v", err)
	}

	if enrollment.EnrollmentState != "active" {
		t.Errorf("Expected state 'active', got %s", enrollment.EnrollmentState)
	}
}

func TestNewEnrollmentsService(t *testing.T) {
	client := &Client{}
	service := NewEnrollmentsService(client)

	if service == nil {
		t.Fatal("NewEnrollmentsService returned nil")
	}
	if service.client != client {
		t.Error("NewEnrollmentsService did not set client correctly")
	}
}

func TestEnrollmentsService_Accept(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/enrollments/456/accept" {
			t.Errorf("Expected path /api/v1/courses/123/enrollments/456/accept, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
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

	service := NewEnrollmentsService(client)
	ctx := context.Background()

	err = service.Accept(ctx, 123, 456)
	if err != nil {
		t.Fatalf("Accept failed: %v", err)
	}
}

func TestEnrollmentsService_Reject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/enrollments/456/reject" {
			t.Errorf("Expected path /api/v1/courses/123/enrollments/456/reject, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
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

	service := NewEnrollmentsService(client)
	ctx := context.Background()

	err = service.Reject(ctx, 123, 456)
	if err != nil {
		t.Fatalf("Reject failed: %v", err)
	}
}

func TestEnrollmentsService_UpdateLastAttended(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/enrollments/456/last_attended" {
			t.Errorf("Expected path /api/v1/courses/123/enrollments/456/last_attended, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 456,
			"user_id": 789,
			"course_id": 123,
			"type": "StudentEnrollment",
			"enrollment_state": "active",
			"last_activity_at": "2024-01-15T10:00:00Z"
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

	service := NewEnrollmentsService(client)
	ctx := context.Background()

	enrollment, err := service.UpdateLastAttended(ctx, 123, 456, "2024-01-15T10:00:00Z")
	if err != nil {
		t.Fatalf("UpdateLastAttended failed: %v", err)
	}

	if enrollment.ID != 456 {
		t.Errorf("Expected enrollment ID 456, got %d", enrollment.ID)
	}
}

func TestEnrollmentsService_ListSection_WithAllOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/sections/456/enrollments" {
			t.Errorf("Expected path /api/v1/sections/456/enrollments, got %s", r.URL.Path)
		}

		// Check all query parameters
		query := r.URL.Query()
		types := query["type[]"]
		if len(types) != 2 {
			t.Errorf("Expected 2 type parameters, got %d", len(types))
		}
		roles := query["role[]"]
		if len(roles) != 1 {
			t.Errorf("Expected 1 role parameter, got %d", len(roles))
		}
		roleIDs := query["role_id[]"]
		if len(roleIDs) != 1 {
			t.Errorf("Expected 1 role_id parameter, got %d", len(roleIDs))
		}
		states := query["state[]"]
		if len(states) != 2 {
			t.Errorf("Expected 2 state parameters, got %d", len(states))
		}
		includes := query["include[]"]
		if len(includes) != 2 {
			t.Errorf("Expected 2 include parameters, got %d", len(includes))
		}
		if query.Get("user_id") != "789" {
			t.Error("Expected user_id parameter")
		}
		if query.Get("grading_period_id") != "5" {
			t.Error("Expected grading_period_id parameter")
		}
		if query.Get("page") != "2" {
			t.Error("Expected page parameter")
		}
		if query.Get("per_page") != "50" {
			t.Error("Expected per_page parameter")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"id": 1,
				"user_id": 789,
				"course_section_id": 456,
				"type": "StudentEnrollment",
				"enrollment_state": "active"
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

	service := NewEnrollmentsService(client)
	ctx := context.Background()

	opts := &ListEnrollmentsOptions{
		Type:            []string{"StudentEnrollment", "TeacherEnrollment"},
		Role:            []string{"StudentEnrollment"},
		RoleID:          []int64{3},
		State:           []string{"active", "invited"},
		Include:         []string{"avatar_url", "user"},
		UserID:          789,
		GradingPeriodID: 5,
		Page:            2,
		PerPage:         50,
	}

	enrollments, err := service.ListSection(ctx, 456, opts)
	if err != nil {
		t.Fatalf("ListSection failed: %v", err)
	}

	if len(enrollments) != 1 {
		t.Errorf("Expected 1 enrollment, got %d", len(enrollments))
	}
}

func TestEnrollmentsService_ListUser_WithAllOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/users/789/enrollments" {
			t.Errorf("Expected path /api/v1/users/789/enrollments, got %s", r.URL.Path)
		}

		// Check all query parameters
		query := r.URL.Query()
		types := query["type[]"]
		if len(types) != 1 {
			t.Errorf("Expected 1 type parameter, got %d", len(types))
		}
		roles := query["role[]"]
		if len(roles) != 1 {
			t.Errorf("Expected 1 role parameter, got %d", len(roles))
		}
		roleIDs := query["role_id[]"]
		if len(roleIDs) != 2 {
			t.Errorf("Expected 2 role_id parameters, got %d", len(roleIDs))
		}
		states := query["state[]"]
		if len(states) != 1 {
			t.Errorf("Expected 1 state parameter, got %d", len(states))
		}
		includes := query["include[]"]
		if len(includes) != 3 {
			t.Errorf("Expected 3 include parameters, got %d", len(includes))
		}
		if query.Get("grading_period_id") != "10" {
			t.Error("Expected grading_period_id parameter")
		}
		if query.Get("enrollment_term_id") != "3" {
			t.Error("Expected enrollment_term_id parameter")
		}
		if query.Get("page") != "1" {
			t.Error("Expected page parameter")
		}
		if query.Get("per_page") != "100" {
			t.Error("Expected per_page parameter")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"id": 1,
				"user_id": 789,
				"course_id": 100,
				"type": "StudentEnrollment",
				"enrollment_state": "active"
			},
			{
				"id": 2,
				"user_id": 789,
				"course_id": 101,
				"type": "StudentEnrollment",
				"enrollment_state": "active"
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

	service := NewEnrollmentsService(client)
	ctx := context.Background()

	opts := &ListEnrollmentsOptions{
		Type:             []string{"StudentEnrollment"},
		Role:             []string{"StudentEnrollment"},
		RoleID:           []int64{3, 4},
		State:            []string{"active"},
		Include:          []string{"avatar_url", "user", "locked"},
		GradingPeriodID:  10,
		EnrollmentTermID: 3,
		Page:             1,
		PerPage:          100,
	}

	enrollments, err := service.ListUser(ctx, 789, opts)
	if err != nil {
		t.Fatalf("ListUser failed: %v", err)
	}

	if len(enrollments) != 2 {
		t.Errorf("Expected 2 enrollments, got %d", len(enrollments))
	}
}

func TestEnrollmentsService_ListCourse_WithAllOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/enrollments" {
			t.Errorf("Expected path /api/v1/courses/123/enrollments, got %s", r.URL.Path)
		}

		// Check all query parameters
		query := r.URL.Query()
		if len(query["type[]"]) == 0 {
			t.Error("Expected type parameters")
		}
		if len(query["role[]"]) == 0 {
			t.Error("Expected role parameters")
		}
		if len(query["state[]"]) == 0 {
			t.Error("Expected state parameters")
		}
		if len(query["include[]"]) == 0 {
			t.Error("Expected include parameters")
		}
		if query.Get("user_id") != "100" {
			t.Error("Expected user_id parameter")
		}
		if query.Get("grading_period_id") != "7" {
			t.Error("Expected grading_period_id parameter")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"id": 1,
				"user_id": 100,
				"course_id": 123,
				"type": "StudentEnrollment",
				"enrollment_state": "active"
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

	service := NewEnrollmentsService(client)
	ctx := context.Background()

	opts := &ListEnrollmentsOptions{
		Type:            []string{"StudentEnrollment", "TeacherEnrollment"},
		Role:            []string{"StudentEnrollment", "TeacherEnrollment"},
		RoleID:          []int64{3, 4},
		State:           []string{"active", "completed"},
		Include:         []string{"avatar_url", "user", "locked", "observed_users"},
		UserID:          100,
		GradingPeriodID: 7,
		Page:            1,
		PerPage:         25,
	}

	enrollments, err := service.ListCourse(ctx, 123, opts)
	if err != nil {
		t.Fatalf("ListCourse failed: %v", err)
	}

	if len(enrollments) != 1 {
		t.Errorf("Expected 1 enrollment, got %d", len(enrollments))
	}
}

func TestEnrollmentsService_EnrollUser_WithAllOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{
			"id": 999,
			"user_id": 789,
			"course_id": 123,
			"type": "ObserverEnrollment",
			"enrollment_state": "invited",
			"role_id": 5
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

	service := NewEnrollmentsService(client)
	ctx := context.Background()

	params := &EnrollUserParams{
		UserID:                         789,
		Type:                           "ObserverEnrollment",
		RoleID:                         5,
		EnrollmentState:                "invited",
		CourseSectionID:                456,
		LimitPrivilegesToCourseSection: true,
		Notify:                         true,
		SelfEnrollmentCode:             "CODE123",
		SelfEnrolled:                   true,
		AssociatedUserID:               100,
	}

	enrollment, err := service.EnrollUser(ctx, 123, params)
	if err != nil {
		t.Fatalf("EnrollUser failed: %v", err)
	}

	if enrollment.ID != 999 {
		t.Errorf("Expected enrollment ID 999, got %d", enrollment.ID)
	}
}
