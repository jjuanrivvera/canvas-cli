package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAssignmentsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/assignments/456" {
			t.Errorf("Expected path /api/v1/courses/123/assignments/456, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 456,
			"name": "Test Assignment",
			"course_id": 123,
			"points_possible": 100.0,
			"grading_type": "points",
			"submission_types": ["online_text_entry"],
			"published": true
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

	service := NewAssignmentsService(client)
	ctx := context.Background()

	assignment, err := service.Get(ctx, 123, 456, nil)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if assignment.ID != 456 {
		t.Errorf("Expected assignment ID 456, got %d", assignment.ID)
	}
	if assignment.Name != "Test Assignment" {
		t.Errorf("Expected assignment name 'Test Assignment', got %s", assignment.Name)
	}
	if assignment.PointsPossible != 100.0 {
		t.Errorf("Expected points possible 100.0, got %.1f", assignment.PointsPossible)
	}
}

func TestAssignmentsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/assignments" {
			t.Errorf("Expected path /api/v1/courses/123/assignments, got %s", r.URL.Path)
		}

		// Check query parameters
		searchTerm := r.URL.Query().Get("search_term")
		if searchTerm != "quiz" {
			t.Errorf("Expected search_term 'quiz', got %s", searchTerm)
		}

		bucket := r.URL.Query().Get("bucket")
		if bucket != "upcoming" {
			t.Errorf("Expected bucket 'upcoming', got %s", bucket)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"id": 1,
				"name": "Quiz 1",
				"points_possible": 50.0
			},
			{
				"id": 2,
				"name": "Quiz 2",
				"points_possible": 50.0
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

	service := NewAssignmentsService(client)
	ctx := context.Background()

	opts := &ListAssignmentsOptions{
		SearchTerm: "quiz",
		Bucket:     "upcoming",
	}

	assignments, err := service.List(ctx, 123, opts)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(assignments) != 2 {
		t.Errorf("Expected 2 assignments, got %d", len(assignments))
	}
	if assignments[0].ID != 1 {
		t.Errorf("Expected first assignment ID 1, got %d", assignments[0].ID)
	}
	if assignments[1].ID != 2 {
		t.Errorf("Expected second assignment ID 2, got %d", assignments[1].ID)
	}
}

func TestAssignmentsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/assignments" {
			t.Errorf("Expected path /api/v1/courses/123/assignments, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{
			"id": 789,
			"name": "New Assignment",
			"course_id": 123,
			"points_possible": 100.0,
			"grading_type": "points",
			"published": false
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

	service := NewAssignmentsService(client)
	ctx := context.Background()

	params := &CreateAssignmentParams{
		Name:           "New Assignment",
		PointsPossible: 100.0,
		GradingType:    "points",
		Published:      false,
	}

	assignment, err := service.Create(ctx, 123, params)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if assignment.ID != 789 {
		t.Errorf("Expected assignment ID 789, got %d", assignment.ID)
	}
	if assignment.Name != "New Assignment" {
		t.Errorf("Expected assignment name 'New Assignment', got %s", assignment.Name)
	}
	if assignment.Published {
		t.Error("Expected assignment to be unpublished")
	}
}

func TestAssignmentsService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/assignments/456" {
			t.Errorf("Expected path /api/v1/courses/123/assignments/456, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 456,
			"name": "Updated Assignment",
			"course_id": 123,
			"points_possible": 150.0,
			"published": true
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

	service := NewAssignmentsService(client)
	ctx := context.Background()

	newPoints := 150.0
	published := true
	params := &UpdateAssignmentParams{
		Name:           "Updated Assignment",
		PointsPossible: &newPoints,
		Published:      &published,
	}

	assignment, err := service.Update(ctx, 123, 456, params)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if assignment.ID != 456 {
		t.Errorf("Expected assignment ID 456, got %d", assignment.ID)
	}
	if assignment.Name != "Updated Assignment" {
		t.Errorf("Expected assignment name 'Updated Assignment', got %s", assignment.Name)
	}
	if assignment.PointsPossible != 150.0 {
		t.Errorf("Expected points possible 150.0, got %.1f", assignment.PointsPossible)
	}
	if !assignment.Published {
		t.Error("Expected assignment to be published")
	}
}

func TestAssignmentsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/assignments/456" {
			t.Errorf("Expected path /api/v1/courses/123/assignments/456, got %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE method, got %s", r.Method)
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

	service := NewAssignmentsService(client)
	ctx := context.Background()

	err = service.Delete(ctx, 123, 456)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestAssignmentsService_BulkUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/assignments/bulk_update" {
			t.Errorf("Expected path /api/v1/courses/123/assignments/bulk_update, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"progress": {
				"id": 999
			}
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

	service := NewAssignmentsService(client)
	ctx := context.Background()

	params := &BulkUpdateParams{
		AssignmentIDs: []int64{1, 2, 3},
		DueAt:         "2024-12-31T23:59:59Z",
	}

	err = service.BulkUpdate(ctx, 123, params)
	if err != nil {
		t.Fatalf("BulkUpdate failed: %v", err)
	}
}

func TestNewAssignmentsService(t *testing.T) {
	client := &Client{}
	service := NewAssignmentsService(client)

	if service == nil {
		t.Fatal("NewAssignmentsService returned nil")
	}
	if service.client != client {
		t.Error("NewAssignmentsService did not set client correctly")
	}
}

func TestAssignmentsService_ListUserAssignments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path == "/api/v1/users/456/courses/456/assignments" {
			if r.Method != http.MethodGet {
				t.Errorf("Expected GET method, got %s", r.Method)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{
					"id": 1,
					"name": "Assignment 1",
					"course_id": 123
				},
				{
					"id": 2,
					"name": "Assignment 2",
					"course_id": 123
				}
			]`))
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

	service := NewAssignmentsService(client)
	ctx := context.Background()

	assignments, err := service.ListUserAssignments(ctx, 456, nil)
	if err != nil {
		t.Fatalf("ListUserAssignments failed: %v", err)
	}

	if len(assignments) != 2 {
		t.Errorf("Expected 2 assignments, got %d", len(assignments))
	}

	if assignments[0].ID != 1 {
		t.Errorf("Expected first assignment ID 1, got %d", assignments[0].ID)
	}
}

func TestAssignmentsService_ListUserAssignments_WithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path == "/api/v1/users/456/courses/456/assignments" {
			// Verify query parameters
			query := r.URL.Query()
			if query.Get("bucket") != "upcoming" {
				t.Errorf("Expected bucket=upcoming, got %s", query.Get("bucket"))
			}

			includes := query["include[]"]
			if len(includes) == 0 {
				t.Error("Expected include[] parameters")
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[]`))
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

	service := NewAssignmentsService(client)
	ctx := context.Background()

	opts := &ListAssignmentsOptions{
		Include: []string{"submission"},
		Bucket:  "upcoming",
	}

	assignments, err := service.ListUserAssignments(ctx, 456, opts)
	if err != nil {
		t.Fatalf("ListUserAssignments failed: %v", err)
	}

	if assignments == nil {
		t.Error("Expected non-nil assignments slice")
	}
}

func TestAssignmentsService_Create_WithAllOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{
			"id": 999,
			"name": "Comprehensive Assignment",
			"course_id": 123,
			"points_possible": 100.0,
			"published": true,
			"moderated_grading": true
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

	service := NewAssignmentsService(client)
	ctx := context.Background()

	params := &CreateAssignmentParams{
		Name:                           "Comprehensive Assignment",
		Position:                       1,
		SubmissionTypes:                []string{"online_text_entry", "online_upload"},
		AllowedExtensions:              []string{"pdf", "doc", "docx"},
		TurnitinEnabled:                true,
		VericiteEnabled:                false,
		TurnitinSettings:               map[string]interface{}{"originality_report_visibility": "immediate"},
		IntegrationData:                map[string]interface{}{"external_id": "ext_123"},
		IntegrationID:                  "integration_123",
		PeerReviews:                    true,
		AutomaticPeerReviews:           true,
		NotifyOfUpdate:                 true,
		GroupCategoryID:                10,
		GradeGroupStudentsIndividually: false,
		ExternalToolTagAttributes:      map[string]interface{}{"url": "https://example.com/tool"},
		PointsPossible:                 100.0,
		GradingType:                    "points",
		DueAt:                          "2024-12-31T23:59:59Z",
		LockAt:                         "2025-01-07T23:59:59Z",
		UnlockAt:                       "2024-12-01T00:00:00Z",
		Description:                    "This is a comprehensive test assignment",
		AssignmentGroupID:              5,
		AssignmentOverrides: []AssignmentOverrideParams{
			{
				StudentIDs: []int64{100, 101},
				Title:      "Extended Time",
				DueAt:      "2025-01-02T23:59:59Z",
				UnlockAt:   "2024-12-01T00:00:00Z",
				LockAt:     "2025-01-09T23:59:59Z",
			},
		},
		OnlyVisibleToOverrides:          true,
		Published:                       true,
		GradingStandardID:               3,
		OmitFromFinalGrade:              false,
		ModeratedGrading:                true,
		GraderCount:                     2,
		FinalGraderID:                   50,
		GraderCommentsVisibleToGraders:  true,
		GradersAnonymousToGraders:       false,
		GraderNamesVisibleToFinalGrader: true,
		AnonymousInstructorAnnotations:  false,
		AnonymousGrading:                false,
		AllowedAttempts:                 3,
		AnnotatableAttachmentID:         200,
		HideInGradebook:                 false,
		PostToSIS:                       true,
		ImportantDates:                  true,
	}

	assignment, err := service.Create(ctx, 123, params)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if assignment.ID != 999 {
		t.Errorf("Expected assignment ID 999, got %d", assignment.ID)
	}
	if assignment.Name != "Comprehensive Assignment" {
		t.Errorf("Expected assignment name 'Comprehensive Assignment', got %s", assignment.Name)
	}
}

func TestAssignmentsService_Update_WithAllOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 456,
			"name": "Updated Comprehensive Assignment",
			"course_id": 123,
			"points_possible": 150.0,
			"published": true
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

	service := NewAssignmentsService(client)
	ctx := context.Background()

	position := 2
	points := 150.0
	turnitin := false
	vericite := true
	peerReviews := false
	automaticPeerReviews := false
	notifyOfUpdate := true
	groupCategoryID := int64(15)
	gradeGroupIndividually := true
	dueAt := "2025-01-15T23:59:59Z"
	lockAt := "2025-01-22T23:59:59Z"
	unlockAt := "2024-12-15T00:00:00Z"
	assignmentGroupID := int64(7)
	onlyVisibleToOverrides := false
	published := true
	gradingStandardID := int64(4)
	omitFromFinalGrade := true
	moderatedGrading := false
	graderCount := 3
	finalGraderID := int64(60)
	graderCommentsVisible := false
	gradersAnonymous := true
	graderNamesVisible := false
	anonymousInstructorAnnotations := true
	anonymousGrading := true
	allowedAttempts := 5
	annotatableAttachmentID := int64(300)
	hideInGradebook := true
	postToSIS := false
	importantDates := false

	params := &UpdateAssignmentParams{
		Name:                           "Updated Comprehensive Assignment",
		Position:                       &position,
		SubmissionTypes:                []string{"online_url", "media_recording"},
		AllowedExtensions:              []string{"txt", "md"},
		TurnitinEnabled:                &turnitin,
		VericiteEnabled:                &vericite,
		TurnitinSettings:               map[string]interface{}{"exclude_quoted": true},
		IntegrationData:                map[string]interface{}{"updated_id": "ext_456"},
		IntegrationID:                  "integration_456",
		PeerReviews:                    &peerReviews,
		AutomaticPeerReviews:           &automaticPeerReviews,
		NotifyOfUpdate:                 &notifyOfUpdate,
		GroupCategoryID:                &groupCategoryID,
		GradeGroupStudentsIndividually: &gradeGroupIndividually,
		ExternalToolTagAttributes:      map[string]interface{}{"url": "https://example.com/updated"},
		PointsPossible:                 &points,
		GradingType:                    "letter_grade",
		DueAt:                          &dueAt,
		LockAt:                         &lockAt,
		UnlockAt:                       &unlockAt,
		Description:                    "Updated comprehensive description",
		AssignmentGroupID:              &assignmentGroupID,
		AssignmentOverrides: []AssignmentOverrideParams{
			{
				StudentIDs: []int64{102, 103, 104},
				Title:      "Group Extension",
				DueAt:      "2025-01-20T23:59:59Z",
				UnlockAt:   "2024-12-15T00:00:00Z",
				LockAt:     "2025-01-27T23:59:59Z",
			},
		},
		OnlyVisibleToOverrides:          &onlyVisibleToOverrides,
		Published:                       &published,
		GradingStandardID:               &gradingStandardID,
		OmitFromFinalGrade:              &omitFromFinalGrade,
		ModeratedGrading:                &moderatedGrading,
		GraderCount:                     &graderCount,
		FinalGraderID:                   &finalGraderID,
		GraderCommentsVisibleToGraders:  &graderCommentsVisible,
		GradersAnonymousToGraders:       &gradersAnonymous,
		GraderNamesVisibleToFinalGrader: &graderNamesVisible,
		AnonymousInstructorAnnotations:  &anonymousInstructorAnnotations,
		AnonymousGrading:                &anonymousGrading,
		AllowedAttempts:                 &allowedAttempts,
		AnnotatableAttachmentID:         &annotatableAttachmentID,
		HideInGradebook:                 &hideInGradebook,
		PostToSIS:                       &postToSIS,
		ImportantDates:                  &importantDates,
	}

	assignment, err := service.Update(ctx, 123, 456, params)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if assignment.ID != 456 {
		t.Errorf("Expected assignment ID 456, got %d", assignment.ID)
	}
	if assignment.PointsPossible != 150.0 {
		t.Errorf("Expected points possible 150.0, got %.1f", assignment.PointsPossible)
	}
}

func TestAssignmentsService_Get_WithIncludes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		// Check include parameters
		includes := r.URL.Query()["include[]"]
		if len(includes) != 3 {
			t.Errorf("Expected 3 include parameters, got %d", len(includes))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 456,
			"name": "Test Assignment",
			"course_id": 123,
			"points_possible": 100.0
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

	service := NewAssignmentsService(client)
	ctx := context.Background()

	include := []string{"submission", "rubric", "rubric_assessment"}
	assignment, err := service.Get(ctx, 123, 456, include)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if assignment.ID != 456 {
		t.Errorf("Expected assignment ID 456, got %d", assignment.ID)
	}
}

func TestAssignmentsService_List_WithAllOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		// Check all query parameters
		query := r.URL.Query()
		if len(query["include[]"]) == 0 {
			t.Error("Expected include parameters")
		}
		if query.Get("search_term") == "" {
			t.Error("Expected search_term parameter")
		}
		if query.Get("override_assignment_dates") != "true" {
			t.Error("Expected override_assignment_dates parameter")
		}
		if query.Get("needs_grading_count_by_section") != "true" {
			t.Error("Expected needs_grading_count_by_section parameter")
		}
		if query.Get("bucket") == "" {
			t.Error("Expected bucket parameter")
		}
		if len(query["assignment_ids[]"]) == 0 {
			t.Error("Expected assignment_ids[] parameter")
		}
		if query.Get("order_by") == "" {
			t.Error("Expected order_by parameter")
		}
		if query.Get("post_to_sis") != "true" {
			t.Error("Expected post_to_sis parameter")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id": 1, "name": "Test Assignment"}]`))
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

	service := NewAssignmentsService(client)
	ctx := context.Background()

	postToSIS := true
	opts := &ListAssignmentsOptions{
		Include:                    []string{"submission", "rubric", "score_statistics"},
		SearchTerm:                 "test",
		OverrideAssignmentDates:    true,
		NeedsGradingCountBySection: true,
		Bucket:                     "past",
		AssignmentIDs:              []int64{1, 2, 3},
		OrderBy:                    "due_at",
		PostToSIS:                  &postToSIS,
	}

	assignments, err := service.List(ctx, 123, opts)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(assignments) != 1 {
		t.Errorf("Expected 1 assignment, got %d", len(assignments))
	}
}

func TestAssignmentsService_List_NilOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		// Verify no query parameters when opts is nil
		if len(r.URL.Query()) > 0 {
			t.Error("Expected no query parameters when opts is nil")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id": 1, "name": "Assignment 1"}]`))
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

	service := NewAssignmentsService(client)
	ctx := context.Background()

	assignments, err := service.List(ctx, 123, nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(assignments) != 1 {
		t.Errorf("Expected 1 assignment, got %d", len(assignments))
	}
}

func TestAssignmentsService_BulkUpdate_WithAllOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/assignments/bulk_update" {
			t.Errorf("Expected path /api/v1/courses/123/assignments/bulk_update, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}

		// Verify request body contains all parameters
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		// Check that all date parameters are present
		if _, ok := body["due_at"]; !ok {
			t.Error("Expected due_at in body")
		}
		if _, ok := body["unlock_at"]; !ok {
			t.Error("Expected unlock_at in body")
		}
		if _, ok := body["lock_at"]; !ok {
			t.Error("Expected lock_at in body")
		}
		if _, ok := body["assignment_ids[]"]; !ok {
			t.Error("Expected assignment_ids[] in body")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"progress": {
				"id": 999
			}
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

	service := NewAssignmentsService(client)
	ctx := context.Background()

	params := &BulkUpdateParams{
		AssignmentIDs: []int64{1, 2, 3, 4, 5},
		DueAt:         "2024-12-31T23:59:59Z",
		UnlockAt:      "2024-12-01T00:00:00Z",
		LockAt:        "2025-01-07T23:59:59Z",
	}

	err = service.BulkUpdate(ctx, 123, params)
	if err != nil {
		t.Fatalf("BulkUpdate failed: %v", err)
	}
}
