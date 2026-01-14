package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSubmissionsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/assignments/456/submissions/789" {
			t.Errorf("Expected path /api/v1/courses/123/assignments/456/submissions/789, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 1,
			"user_id": 789,
			"assignment_id": 456,
			"workflow_state": "graded",
			"score": 95.0,
			"grade": "A",
			"submission_type": "online_text_entry",
			"submitted_at": "2024-01-15T10:00:00Z"
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

	service := NewSubmissionsService(client)
	ctx := context.Background()

	submission, err := service.Get(ctx, 123, 456, 789, nil)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if submission.ID != 1 {
		t.Errorf("Expected submission ID 1, got %d", submission.ID)
	}
	if submission.UserID != 789 {
		t.Errorf("Expected user ID 789, got %d", submission.UserID)
	}
	if submission.Score != 95.0 {
		t.Errorf("Expected score 95.0, got %.1f", submission.Score)
	}
	if submission.Grade != "A" {
		t.Errorf("Expected grade 'A', got %s", submission.Grade)
	}
}

func TestSubmissionsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/assignments/456/submissions" {
			t.Errorf("Expected path /api/v1/courses/123/assignments/456/submissions, got %s", r.URL.Path)
		}

		// Check query parameters
		workflowState := r.URL.Query().Get("workflow_state")
		if workflowState != "graded" {
			t.Errorf("Expected workflow_state 'graded', got %s", workflowState)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"id": 1,
				"user_id": 100,
				"assignment_id": 456,
				"workflow_state": "graded",
				"score": 85.0
			},
			{
				"id": 2,
				"user_id": 101,
				"assignment_id": 456,
				"workflow_state": "graded",
				"score": 92.0
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

	service := NewSubmissionsService(client)
	ctx := context.Background()

	opts := &ListSubmissionsOptions{
		WorkflowState: "graded",
	}

	submissions, err := service.List(ctx, 123, 456, opts)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(submissions) != 2 {
		t.Errorf("Expected 2 submissions, got %d", len(submissions))
	}
	if submissions[0].UserID != 100 {
		t.Errorf("Expected first submission user ID 100, got %d", submissions[0].UserID)
	}
	if submissions[1].UserID != 101 {
		t.Errorf("Expected second submission user ID 101, got %d", submissions[1].UserID)
	}
}

func TestSubmissionsService_ListMultiple(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/students/submissions" {
			t.Errorf("Expected path /api/v1/courses/123/students/submissions, got %s", r.URL.Path)
		}

		// Check query parameters
		studentIDs := r.URL.Query()["student_ids[]"]
		if len(studentIDs) != 2 {
			t.Errorf("Expected 2 student IDs, got %d", len(studentIDs))
		}

		assignmentIDs := r.URL.Query()["assignment_ids[]"]
		if len(assignmentIDs) != 2 {
			t.Errorf("Expected 2 assignment IDs, got %d", len(assignmentIDs))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"id": 1,
				"user_id": 100,
				"assignment_id": 200
			},
			{
				"id": 2,
				"user_id": 100,
				"assignment_id": 201
			},
			{
				"id": 3,
				"user_id": 101,
				"assignment_id": 200
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

	service := NewSubmissionsService(client)
	ctx := context.Background()

	studentIDs := []int64{100, 101}
	assignmentIDs := []int64{200, 201}

	submissions, err := service.ListMultiple(ctx, 123, studentIDs, assignmentIDs, nil)
	if err != nil {
		t.Fatalf("ListMultiple failed: %v", err)
	}

	if len(submissions) != 3 {
		t.Errorf("Expected 3 submissions, got %d", len(submissions))
	}
}

func TestSubmissionsService_Grade(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/assignments/456/submissions/789" {
			t.Errorf("Expected path /api/v1/courses/123/assignments/456/submissions/789, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 1,
			"user_id": 789,
			"assignment_id": 456,
			"workflow_state": "graded",
			"score": 88.0,
			"grade": "B+"
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

	service := NewSubmissionsService(client)
	ctx := context.Background()

	params := &GradeSubmissionParams{
		PostedGrade: "B+",
	}

	submission, err := service.Grade(ctx, 123, 456, 789, params)
	if err != nil {
		t.Fatalf("Grade failed: %v", err)
	}

	if submission.Score != 88.0 {
		t.Errorf("Expected score 88.0, got %.1f", submission.Score)
	}
	if submission.Grade != "B+" {
		t.Errorf("Expected grade 'B+', got %s", submission.Grade)
	}
}

func TestSubmissionsService_Submit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/assignments/456/submissions" {
			t.Errorf("Expected path /api/v1/courses/123/assignments/456/submissions, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{
			"id": 1,
			"user_id": 789,
			"assignment_id": 456,
			"workflow_state": "submitted",
			"submission_type": "online_text_entry",
			"body": "My submission text"
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

	service := NewSubmissionsService(client)
	ctx := context.Background()

	params := &SubmitParams{
		SubmissionType: "online_text_entry",
		Body:           "My submission text",
	}

	submission, err := service.Submit(ctx, 123, 456, params)
	if err != nil {
		t.Fatalf("Submit failed: %v", err)
	}

	if submission.SubmissionType != "online_text_entry" {
		t.Errorf("Expected submission type 'online_text_entry', got %s", submission.SubmissionType)
	}
	if submission.Body != "My submission text" {
		t.Errorf("Expected body 'My submission text', got %s", submission.Body)
	}
}

func TestSubmissionsService_MarkAsRead(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/assignments/456/submissions/789/read" {
			t.Errorf("Expected path /api/v1/courses/123/assignments/456/submissions/789/read, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}

		w.WriteHeader(http.StatusNoContent)
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

	service := NewSubmissionsService(client)
	ctx := context.Background()

	err = service.MarkAsRead(ctx, 123, 456, 789)
	if err != nil {
		t.Fatalf("MarkAsRead failed: %v", err)
	}
}

func TestNewSubmissionsService(t *testing.T) {
	client := &Client{}
	service := NewSubmissionsService(client)

	if service == nil {
		t.Fatal("NewSubmissionsService returned nil")
	}
	if service.client != client {
		t.Error("NewSubmissionsService did not set client correctly")
	}
}

func TestSubmissionsService_MarkAsUnread(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path == "/api/v1/courses/123/assignments/456/submissions/789/read" {
			if r.Method != http.MethodDelete {
				t.Errorf("Expected DELETE method, got %s", r.Method)
			}

			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.WriteHeader(http.StatusNotFound)
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

	service := NewSubmissionsService(client)
	ctx := context.Background()

	err = service.MarkAsUnread(ctx, 123, 456, 789)
	if err != nil {
		t.Fatalf("MarkAsUnread failed: %v", err)
	}
}

func TestSubmissionsService_InitiateFileUpload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path == "/api/v1/courses/123/assignments/456/submissions/self/files" {
			if r.Method != http.MethodPost {
				t.Errorf("Expected POST method, got %s", r.Method)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"upload_url": "https://example.com/upload",
				"upload_params": {
					"key": "value"
				}
			}`))
			return
		}

		w.WriteHeader(http.StatusNotFound)
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

	service := NewSubmissionsService(client)
	ctx := context.Background()

	params := &UploadFileParams{
		Name:        "test.pdf",
		Size:        1024,
		ContentType: "application/pdf",
		OnDuplicate: "rename",
	}

	result, err := service.InitiateFileUpload(ctx, 123, 456, params)
	if err != nil {
		t.Fatalf("InitiateFileUpload failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	uploadURL, ok := result["upload_url"].(string)
	if !ok || uploadURL != "https://example.com/upload" {
		t.Errorf("Expected upload_url to be 'https://example.com/upload', got %v", uploadURL)
	}
}

func TestSubmissionsService_BulkGrade(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path == "/api/v1/courses/123/assignments/456/submissions/update_grades" {
			if r.Method != http.MethodPost {
				t.Errorf("Expected POST method, got %s", r.Method)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": 1, "workflow_state": "complete"}`))
			return
		}

		w.WriteHeader(http.StatusNotFound)
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

	service := NewSubmissionsService(client)
	ctx := context.Background()

	params := &BulkGradeParams{
		GradeData: map[int64]GradeData{
			101: {
				PostedGrade: "A",
			},
			102: {
				PostedGrade: "B",
			},
		},
	}

	_, err = service.BulkGrade(ctx, 123, 456, params)
	if err != nil {
		t.Fatalf("BulkGrade failed: %v", err)
	}
}

func TestSubmissionsService_Grade_WithAllOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/assignments/456/submissions/789" {
			t.Errorf("Expected path /api/v1/courses/123/assignments/456/submissions/789, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 1,
			"user_id": 789,
			"assignment_id": 456,
			"workflow_state": "graded",
			"score": 88.0,
			"grade": "B+",
			"excused": true,
			"late_policy_status": "missing"
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

	service := NewSubmissionsService(client)
	ctx := context.Background()

	secondsLate := 3600
	params := &GradeSubmissionParams{
		PostedGrade:         "B+",
		Excuse:              true,
		LatePolicyStatus:    "missing",
		SecondsLateOverride: &secondsLate,
		Comment: &SubmissionCommentParams{
			TextComment:  "Good work!",
			GroupComment: true,
		},
		RubricAssessment: map[string]RubricAssessmentParams{
			"criterion_1": {
				Points:   8.5,
				Rating:   "rating_1",
				Comments: "Excellent",
			},
		},
	}

	submission, err := service.Grade(ctx, 123, 456, 789, params)
	if err != nil {
		t.Fatalf("Grade failed: %v", err)
	}

	if submission.Score != 88.0 {
		t.Errorf("Expected score 88.0, got %.1f", submission.Score)
	}
	if submission.Grade != "B+" {
		t.Errorf("Expected grade 'B+', got %s", submission.Grade)
	}
}

func TestSubmissionsService_Grade_WithMediaComment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": 1, "user_id": 789, "assignment_id": 456}`))
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

	service := NewSubmissionsService(client)
	ctx := context.Background()

	params := &GradeSubmissionParams{
		PostedGrade: "A",
		Comment: &SubmissionCommentParams{
			TextComment:      "Great job",
			MediaCommentID:   "media_123",
			MediaCommentType: "audio",
			FileIDs:          []int64{111, 222},
		},
	}

	_, err = service.Grade(ctx, 123, 456, 789, params)
	if err != nil {
		t.Fatalf("Grade failed: %v", err)
	}
}

func TestSubmissionsService_ListMultiple_WithAllOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/students/submissions" {
			t.Errorf("Expected path /api/v1/courses/123/students/submissions, got %s", r.URL.Path)
		}

		// Check various query parameters
		if r.URL.Query().Get("workflow_state") != "graded" {
			t.Error("Expected workflow_state parameter")
		}
		if r.URL.Query().Get("enrollment_state") != "active" {
			t.Error("Expected enrollment_state parameter")
		}
		if r.URL.Query().Get("order") != "graded_at" {
			t.Error("Expected order parameter")
		}
		if r.URL.Query().Get("order_direction") != "descending" {
			t.Error("Expected order_direction parameter")
		}
		if r.URL.Query().Get("grouped") != "true" {
			t.Error("Expected grouped parameter")
		}
		if r.URL.Query().Get("state_based_on_date") != "true" {
			t.Error("Expected state_based_on_date parameter")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{"id": 1, "user_id": 100, "assignment_id": 200},
			{"id": 2, "user_id": 101, "assignment_id": 201}
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

	service := NewSubmissionsService(client)
	ctx := context.Background()

	postToSIS := true
	opts := &ListSubmissionsOptions{
		Include:          []string{"submission_history", "user"},
		Grouped:          true,
		PostToSIS:        &postToSIS,
		SubmittedSince:   "2024-01-01T00:00:00Z",
		GradedSince:      "2024-01-02T00:00:00Z",
		GradingPeriodID:  1,
		WorkflowState:    "graded",
		EnrollmentState:  "active",
		StateBasedOnDate: true,
		Order:            "graded_at",
		OrderDirection:   "descending",
		Page:             2,
		PerPage:          50,
	}

	studentIDs := []int64{100, 101}
	assignmentIDs := []int64{200, 201}

	submissions, err := service.ListMultiple(ctx, 123, studentIDs, assignmentIDs, opts)
	if err != nil {
		t.Fatalf("ListMultiple failed: %v", err)
	}

	if len(submissions) != 2 {
		t.Errorf("Expected 2 submissions, got %d", len(submissions))
	}
}

func TestSubmissionsService_Get_WithIncludes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/assignments/456/submissions/789" {
			t.Errorf("Expected path /api/v1/courses/123/assignments/456/submissions/789, got %s", r.URL.Path)
		}

		// Check include parameters
		includes := r.URL.Query()["include[]"]
		if len(includes) != 2 {
			t.Errorf("Expected 2 include parameters, got %d", len(includes))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 1,
			"user_id": 789,
			"assignment_id": 456,
			"workflow_state": "graded",
			"score": 95.0
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

	service := NewSubmissionsService(client)
	ctx := context.Background()

	include := []string{"submission_history", "rubric_assessment"}
	submission, err := service.Get(ctx, 123, 456, 789, include)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if submission.ID != 1 {
		t.Errorf("Expected submission ID 1, got %d", submission.ID)
	}
}

func TestSubmissionsService_List_WithAllOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		// Check all query parameters
		query := r.URL.Query()
		if query.Get("grouped") != "true" {
			t.Error("Expected grouped parameter")
		}
		if query.Get("post_to_sis") != "true" {
			t.Error("Expected post_to_sis parameter")
		}
		if query.Get("submitted_since") == "" {
			t.Error("Expected submitted_since parameter")
		}
		if query.Get("graded_since") == "" {
			t.Error("Expected graded_since parameter")
		}
		if query.Get("grading_period_id") != "5" {
			t.Error("Expected grading_period_id parameter")
		}
		if query.Get("page") != "3" {
			t.Error("Expected page parameter")
		}
		if query.Get("per_page") != "25" {
			t.Error("Expected per_page parameter")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id": 1, "user_id": 100, "assignment_id": 456}]`))
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

	service := NewSubmissionsService(client)
	ctx := context.Background()

	postToSIS := true
	opts := &ListSubmissionsOptions{
		Include:          []string{"user", "assignment"},
		Grouped:          true,
		PostToSIS:        &postToSIS,
		SubmittedSince:   "2024-01-01T00:00:00Z",
		GradedSince:      "2024-01-15T00:00:00Z",
		GradingPeriodID:  5,
		WorkflowState:    "submitted",
		EnrollmentState:  "active",
		StateBasedOnDate: true,
		Order:            "id",
		OrderDirection:   "ascending",
		Page:             3,
		PerPage:          25,
	}

	submissions, err := service.List(ctx, 123, 456, opts)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(submissions) != 1 {
		t.Errorf("Expected 1 submission, got %d", len(submissions))
	}
}

func TestSubmissionsService_Submit_WithAllOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{
			"id": 1,
			"user_id": 789,
			"assignment_id": 456,
			"workflow_state": "submitted"
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

	service := NewSubmissionsService(client)
	ctx := context.Background()

	params := &SubmitParams{
		SubmissionType:   "online_upload",
		FileIDs:          []int64{111, 222, 333},
		UserID:           789,
		MediaCommentID:   "media_456",
		MediaCommentType: "video",
		Comment: &SubmissionCommentParams{
			TextComment:      "Here's my submission",
			GroupComment:     true,
			MediaCommentID:   "media_789",
			MediaCommentType: "audio",
			FileIDs:          []int64{444},
		},
	}

	submission, err := service.Submit(ctx, 123, 456, params)
	if err != nil {
		t.Fatalf("Submit failed: %v", err)
	}

	if submission.ID != 1 {
		t.Errorf("Expected submission ID 1, got %d", submission.ID)
	}
}

func TestSubmissionsService_Submit_OnlineURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{
			"id": 1,
			"user_id": 789,
			"assignment_id": 456,
			"submission_type": "online_url",
			"url": "https://example.com"
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

	service := NewSubmissionsService(client)
	ctx := context.Background()

	params := &SubmitParams{
		SubmissionType: "online_url",
		URL:            "https://example.com",
	}

	submission, err := service.Submit(ctx, 123, 456, params)
	if err != nil {
		t.Fatalf("Submit failed: %v", err)
	}

	if submission.SubmissionType != "online_url" {
		t.Errorf("Expected submission type 'online_url', got %s", submission.SubmissionType)
	}
}

func TestSubmissionsService_BulkGrade_WithRubric(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"progress": {"id": 123}}`))
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

	service := NewSubmissionsService(client)
	ctx := context.Background()

	params := &BulkGradeParams{
		GradeData: map[int64]GradeData{
			101: {
				PostedGrade:      "A",
				Excuse:           true,
				LatePolicyStatus: "late",
				RubricAssessment: map[string]RubricAssessmentParams{
					"criterion_1": {
						Points:   10.0,
						Rating:   "excellent",
						Comments: "Outstanding work",
					},
				},
			},
		},
	}

	_, err = service.BulkGrade(ctx, 123, 456, params)
	if err != nil {
		t.Fatalf("BulkGrade failed: %v", err)
	}
}

func TestSubmissionsService_DeleteComment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		expectedPath := "/api/v1/courses/123/assignments/456/submissions/789/comments/999"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SubmissionComment{
			ID:         999,
			AuthorID:   100,
			AuthorName: "Test User",
			Comment:    "Deleted comment",
		})
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

	service := NewSubmissionsService(client)
	ctx := context.Background()

	comment, err := service.DeleteComment(ctx, 123, 456, 789, 999)
	if err != nil {
		t.Fatalf("DeleteComment failed: %v", err)
	}

	if comment == nil {
		t.Fatal("expected comment to be returned")
	}

	if comment.ID != 999 {
		t.Errorf("expected comment ID 999, got %d", comment.ID)
	}
}
