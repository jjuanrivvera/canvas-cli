package api

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
)

// assignmentAnalyticsFixture mirrors the documented example response for "Get
// course-level assignment data" from
// https://canvas.instructure.com/doc/api/analytics.html. Unlike the student
// summary, this endpoint reports tardiness_breakdown as fractions of the class
// (and sends 0.0, not 0, for the empty buckets), which is why the two levels
// cannot share one Go type.
const assignmentAnalyticsFixture = `{"assignment_id":1234,"title":"Assignment 1","points_possible":10,"due_at":"2013-01-01T12:00:00Z","unlock_at":"2012-12-31T12:00:00Z","muted":false,"min_score":2,"max_score":10,"median":7,"first_quartile":4,"third_quartile":8,"tardiness_breakdown":{"on_time":0.6666666666666666,"late":0.3333333333333333,"missing":0.0,"floating":0.0,"total":3}}`

func closeEnough(got, want float64) bool {
	return math.Abs(got-want) < 1e-9
}

// TestAssignmentAnalytics_UnmarshalJSON_FractionalTardiness decodes the
// documented example directly. With the assignment-level breakdown typed as
// integer counts, the fractional on_time value fails to decode and takes the
// whole response with it.
func TestAssignmentAnalytics_UnmarshalJSON_FractionalTardiness(t *testing.T) {
	var a AssignmentAnalytics
	if err := json.Unmarshal([]byte(assignmentAnalyticsFixture), &a); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if a.AssignmentID != 1234 {
		t.Errorf("AssignmentID = %d, want 1234", a.AssignmentID)
	}
	if a.Title != "Assignment 1" {
		t.Errorf("Title = %q, want %q", a.Title, "Assignment 1")
	}
	if a.MedianScore != 7 {
		t.Errorf("MedianScore = %v, want 7", a.MedianScore)
	}

	if a.Tardiness == nil {
		t.Fatal("Tardiness = nil, want populated (tardiness_breakdown field)")
	}
	if !closeEnough(a.Tardiness.OnTime, 0.6666666666666666) {
		t.Errorf("Tardiness.OnTime = %v, want 0.6666666666666666", a.Tardiness.OnTime)
	}
	if !closeEnough(a.Tardiness.Late, 0.3333333333333333) {
		t.Errorf("Tardiness.Late = %v, want 0.3333333333333333", a.Tardiness.Late)
	}
	if a.Tardiness.Missing != 0 {
		t.Errorf("Tardiness.Missing = %v, want 0", a.Tardiness.Missing)
	}
	if a.Tardiness.Total != 3 {
		t.Errorf("Tardiness.Total = %v, want 3", a.Tardiness.Total)
	}
}

// TestAnalyticsService_GetCourseAssignments_FractionalTardiness drives the
// same fixture through the real call path. GetCourseAssignments decodes with
// GetJSON, which surfaces a decode failure to the caller, so an int-typed
// breakdown made `canvas analytics assignments` error out for any course whose
// assignment analytics carry a fractional on-time share.
func TestAnalyticsService_GetCourseAssignments_FractionalTardiness(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/analytics/assignments" {
			t.Errorf("Expected path /api/v1/courses/123/analytics/assignments, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[" + assignmentAnalyticsFixture + "]"))
	}))
	defer server.Close()

	service := NewAnalyticsService(newTestClient(t, server.URL))

	assignments, err := service.GetCourseAssignments(context.Background(), 123)
	if err != nil {
		t.Fatalf("GetCourseAssignments failed: %v", err)
	}

	if len(assignments) != 1 {
		t.Fatalf("expected 1 assignment, got %d", len(assignments))
	}
	if got := assignments[0].Tardiness; got == nil || !closeEnough(got.OnTime, 0.6666666666666666) {
		t.Errorf("Tardiness = %+v, want OnTime=0.6666666666666666", got)
	}
}
