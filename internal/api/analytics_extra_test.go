package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAnalyticsService_GetStudentSummaries_WithOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		q := r.URL.Query()
		if q.Get("sort_column") != "score" {
			t.Errorf("expected sort_column=score, got %q", q.Get("sort_column"))
		}
		if q.Get("student_id") != "42" {
			t.Errorf("expected student_id=42, got %q", q.Get("student_id"))
		}
		if q.Get("page") != "2" {
			t.Errorf("expected page=2, got %q", q.Get("page"))
		}
		if q.Get("per_page") != "25" {
			t.Errorf("expected per_page=25, got %q", q.Get("per_page"))
		}
		summaries := []StudentSummary{
			{ID: 42, PageViews: 50, Participations: 10, CurrentScore: 88.0},
		}
		json.NewEncoder(w).Encode(summaries)
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewAnalyticsService(client)
	opts := &ListStudentSummariesOptions{
		SortColumn: "score",
		StudentID:  42,
		Page:       2,
		PerPage:    25,
	}
	summaries, err := svc.GetStudentSummaries(context.Background(), 99, opts)
	if err != nil {
		t.Fatalf("GetStudentSummaries: %v", err)
	}
	if len(summaries) != 1 {
		t.Errorf("expected 1, got %d", len(summaries))
	}
	if summaries[0].CurrentScore != 88.0 {
		t.Errorf("expected score 88.0, got %f", summaries[0].CurrentScore)
	}
}

func TestAnalyticsService_GetUserActivity_PageViews(t *testing.T) {
	// Test the page_views/participations format branch of GetUserActivity.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		// Return page_views format
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"page_views":{"2024-01-01":5,"2024-01-02":3},"participations":[]}`))
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewAnalyticsService(client)
	activity, err := svc.GetUserActivity(context.Background(), 10, 20)
	if err != nil {
		t.Fatalf("GetUserActivity (page_views format): %v", err)
	}
	if len(activity) != 2 {
		t.Errorf("expected 2 activity records, got %d", len(activity))
	}
}

func TestAnalyticsService_GetUserActivity_MapFormat(t *testing.T) {
	// Test the date->activity map format branch.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"2024-01-01":{"views":10,"participations":3},"2024-01-02":{"views":7,"participations":1}}`))
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewAnalyticsService(client)
	activity, err := svc.GetUserActivity(context.Background(), 10, 20)
	if err != nil {
		t.Fatalf("GetUserActivity (map format): %v", err)
	}
	if len(activity) != 2 {
		t.Errorf("expected 2 activity records, got %d", len(activity))
	}
}

func TestAnalyticsService_GetUserActivity_InvalidFormat(t *testing.T) {
	// Return something that cannot be decoded as any of the known formats.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// A plain string is not a valid activity response.
		w.Write([]byte(`"unexpected"`))
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewAnalyticsService(client)
	_, err = svc.GetUserActivity(context.Background(), 10, 20)
	if err == nil {
		t.Error("expected error for invalid format, got nil")
	}
}

func TestAnalyticsService_GetDepartmentStatistics_WithOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		q := r.URL.Query()
		if q.Get("term_id") != "5" {
			t.Errorf("expected term_id=5, got %q", q.Get("term_id"))
		}
		if q.Get("start_date") != "2024-01-01" {
			t.Errorf("expected start_date=2024-01-01, got %q", q.Get("start_date"))
		}
		if q.Get("end_date") != "2024-12-31" {
			t.Errorf("expected end_date=2024-12-31, got %q", q.Get("end_date"))
		}
		stats := DepartmentStatistics{Subaccounts: 3, Teachers: 10, Students: 200}
		json.NewEncoder(w).Encode(stats)
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewAnalyticsService(client)
	opts := &DepartmentAnalyticsOptions{TermID: 5, StartDate: "2024-01-01", EndDate: "2024-12-31"}
	stats, err := svc.GetDepartmentStatistics(context.Background(), 1, opts)
	if err != nil {
		t.Fatalf("GetDepartmentStatistics: %v", err)
	}
	if stats.Students != 200 {
		t.Errorf("expected 200 students, got %d", stats.Students)
	}
}

func TestAnalyticsService_GetDepartmentActivity_WithOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		q := r.URL.Query()
		if q.Get("term_id") != "3" {
			t.Errorf("expected term_id=3, got %q", q.Get("term_id"))
		}
		json.NewEncoder(w).Encode([]DepartmentActivity{{Date: "2024-06-01", Views: 50, Participations: 20}})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewAnalyticsService(client)
	opts := &DepartmentAnalyticsOptions{TermID: 3}
	activity, err := svc.GetDepartmentActivity(context.Background(), 1, opts)
	if err != nil {
		t.Fatalf("GetDepartmentActivity: %v", err)
	}
	if len(activity) != 1 {
		t.Errorf("expected 1 record, got %d", len(activity))
	}
}

func TestAnalyticsService_GetDepartmentGrades_WithOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		q := r.URL.Query()
		if q.Get("start_date") != "2024-01-01" {
			t.Errorf("expected start_date, got %q", q.Get("start_date"))
		}
		json.NewEncoder(w).Encode([]DepartmentGrades{{Score: 95.0, Count: 30}})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewAnalyticsService(client)
	opts := &DepartmentAnalyticsOptions{StartDate: "2024-01-01"}
	grades, err := svc.GetDepartmentGrades(context.Background(), 1, opts)
	if err != nil {
		t.Fatalf("GetDepartmentGrades: %v", err)
	}
	if len(grades) != 1 {
		t.Errorf("expected 1 record, got %d", len(grades))
	}
}

func TestAnalyticsService_GetCourseActivity_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":[{"message":"not found"}]}`))
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewAnalyticsService(client)
	_, err = svc.GetCourseActivity(context.Background(), 999)
	if err == nil {
		t.Error("expected error, got nil")
	}
}
