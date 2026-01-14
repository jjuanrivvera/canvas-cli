package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAnalyticsService_GetCourseActivity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/analytics/activity" {
			t.Errorf("Expected path /api/v1/courses/123/analytics/activity, got %s", r.URL.Path)
		}

		activity := []CourseActivity{
			{Date: "2024-01-01", Views: 100, Participations: 50},
			{Date: "2024-01-02", Views: 120, Participations: 60},
		}
		json.NewEncoder(w).Encode(activity)
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

	service := NewAnalyticsService(client)
	activity, err := service.GetCourseActivity(context.Background(), 123)
	if err != nil {
		t.Fatalf("GetCourseActivity failed: %v", err)
	}

	if len(activity) != 2 {
		t.Errorf("Expected 2 activity records, got %d", len(activity))
	}

	if activity[0].Views != 100 {
		t.Errorf("Expected 100 views, got %d", activity[0].Views)
	}
}

func TestAnalyticsService_GetCourseAssignments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/analytics/assignments" {
			t.Errorf("Expected path /api/v1/courses/123/analytics/assignments, got %s", r.URL.Path)
		}

		assignments := []AssignmentAnalytics{
			{AssignmentID: 1, Title: "Assignment 1", PointsPossible: 100, MinScore: 50, MaxScore: 100, MedianScore: 85},
			{AssignmentID: 2, Title: "Assignment 2", PointsPossible: 50, MinScore: 20, MaxScore: 50, MedianScore: 40},
		}
		json.NewEncoder(w).Encode(assignments)
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

	service := NewAnalyticsService(client)
	assignments, err := service.GetCourseAssignments(context.Background(), 123)
	if err != nil {
		t.Fatalf("GetCourseAssignments failed: %v", err)
	}

	if len(assignments) != 2 {
		t.Errorf("Expected 2 assignments, got %d", len(assignments))
	}

	if assignments[0].Title != "Assignment 1" {
		t.Errorf("Expected 'Assignment 1', got '%s'", assignments[0].Title)
	}
}

func TestAnalyticsService_GetStudentSummaries(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/analytics/student_summaries" {
			t.Errorf("Expected path /api/v1/courses/123/analytics/student_summaries, got %s", r.URL.Path)
		}

		summaries := []StudentSummary{
			{ID: 1, PageViews: 100, Participations: 50, CurrentScore: 95.5, CurrentGrade: "A"},
			{ID: 2, PageViews: 80, Participations: 40, CurrentScore: 85.0, CurrentGrade: "B"},
		}
		json.NewEncoder(w).Encode(summaries)
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

	service := NewAnalyticsService(client)
	summaries, err := service.GetStudentSummaries(context.Background(), 123, nil)
	if err != nil {
		t.Fatalf("GetStudentSummaries failed: %v", err)
	}

	if len(summaries) != 2 {
		t.Errorf("Expected 2 summaries, got %d", len(summaries))
	}

	if summaries[0].CurrentGrade != "A" {
		t.Errorf("Expected grade 'A', got '%s'", summaries[0].CurrentGrade)
	}
}

func TestAnalyticsService_GetUserActivity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/analytics/users/456/activity" {
			t.Errorf("Expected path /api/v1/courses/123/analytics/users/456/activity, got %s", r.URL.Path)
		}

		activity := []UserActivity{
			{Date: "2024-01-01", Views: 10, Participations: 5},
		}
		json.NewEncoder(w).Encode(activity)
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

	service := NewAnalyticsService(client)
	activity, err := service.GetUserActivity(context.Background(), 123, 456)
	if err != nil {
		t.Fatalf("GetUserActivity failed: %v", err)
	}

	if len(activity) != 1 {
		t.Errorf("Expected 1 activity record, got %d", len(activity))
	}
}

func TestAnalyticsService_GetUserAssignments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/analytics/users/456/assignments" {
			t.Errorf("Expected path /api/v1/courses/123/analytics/users/456/assignments, got %s", r.URL.Path)
		}

		assignments := []UserAssignmentAnalytics{
			{AssignmentID: 1, Title: "Test Assignment", PointsPossible: 100, Score: 95},
		}
		json.NewEncoder(w).Encode(assignments)
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

	service := NewAnalyticsService(client)
	assignments, err := service.GetUserAssignments(context.Background(), 123, 456)
	if err != nil {
		t.Fatalf("GetUserAssignments failed: %v", err)
	}

	if len(assignments) != 1 {
		t.Errorf("Expected 1 assignment, got %d", len(assignments))
	}

	if assignments[0].Score != 95 {
		t.Errorf("Expected score 95, got %.1f", assignments[0].Score)
	}
}

func TestAnalyticsService_GetUserCommunication(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/analytics/users/456/communication" {
			t.Errorf("Expected path /api/v1/courses/123/analytics/users/456/communication, got %s", r.URL.Path)
		}

		communication := UserCommunication{
			InstructorMessages: 5,
			StudentMessages:    10,
		}
		json.NewEncoder(w).Encode(communication)
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

	service := NewAnalyticsService(client)
	communication, err := service.GetUserCommunication(context.Background(), 123, 456)
	if err != nil {
		t.Fatalf("GetUserCommunication failed: %v", err)
	}

	if communication.InstructorMessages != 5 {
		t.Errorf("Expected 5 instructor messages, got %d", communication.InstructorMessages)
	}
}

func TestAnalyticsService_GetDepartmentStatistics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/analytics/current/statistics" {
			t.Errorf("Expected path /api/v1/accounts/1/analytics/current/statistics, got %s", r.URL.Path)
		}

		stats := DepartmentStatistics{
			Subaccounts:      5,
			Teachers:         20,
			Students:         500,
			Assignments:      100,
			DiscussionTopics: 50,
		}
		json.NewEncoder(w).Encode(stats)
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

	service := NewAnalyticsService(client)
	stats, err := service.GetDepartmentStatistics(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("GetDepartmentStatistics failed: %v", err)
	}

	if stats.Students != 500 {
		t.Errorf("Expected 500 students, got %d", stats.Students)
	}
}

func TestAnalyticsService_GetDepartmentActivity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/analytics/current/activity" {
			t.Errorf("Expected path /api/v1/accounts/1/analytics/current/activity, got %s", r.URL.Path)
		}

		activity := []DepartmentActivity{
			{Date: "2024-01-01", Views: 1000, Participations: 500},
		}
		json.NewEncoder(w).Encode(activity)
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

	service := NewAnalyticsService(client)
	activity, err := service.GetDepartmentActivity(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("GetDepartmentActivity failed: %v", err)
	}

	if len(activity) != 1 {
		t.Errorf("Expected 1 activity record, got %d", len(activity))
	}
}

func TestAnalyticsService_GetDepartmentGrades(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/analytics/current/grades" {
			t.Errorf("Expected path /api/v1/accounts/1/analytics/current/grades, got %s", r.URL.Path)
		}

		grades := []DepartmentGrades{
			{Score: 90, Count: 100},
			{Score: 80, Count: 150},
		}
		json.NewEncoder(w).Encode(grades)
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

	service := NewAnalyticsService(client)
	grades, err := service.GetDepartmentGrades(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("GetDepartmentGrades failed: %v", err)
	}

	if len(grades) != 2 {
		t.Errorf("Expected 2 grade records, got %d", len(grades))
	}
}
