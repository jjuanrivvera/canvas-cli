package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// studentSummaryFixture is the documented example response for "Get
// course-level student summary data" from
// https://canvas.instructure.com/doc/api/analytics.html, copied verbatim
// (not derived from the Go struct) so the test fails if the struct's JSON
// tags or field types drift from what Canvas actually sends.
const studentSummaryFixture = `{"id":2346,"page_views":351,"page_views_level":"1","max_page_views":415,"participations":1,"participations_level":"3","max_participations":3,"tardiness_breakdown":{"total":5,"on_time":3,"late":0,"missing":2,"floating":0}}`

// TestStudentSummary_UnmarshalJSON_DocumentedShape decodes the documented
// example response directly and checks every field, including
// tardiness_breakdown and the two level fields Canvas sends as quoted
// integers.
func TestStudentSummary_UnmarshalJSON_DocumentedShape(t *testing.T) {
	var s StudentSummary
	if err := json.Unmarshal([]byte(studentSummaryFixture), &s); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if s.ID != 2346 {
		t.Errorf("ID = %d, want 2346", s.ID)
	}
	if s.PageViews != 351 {
		t.Errorf("PageViews = %d, want 351", s.PageViews)
	}
	if s.PageViewsLevel != 1 {
		t.Errorf("PageViewsLevel = %d, want 1", s.PageViewsLevel)
	}
	if s.MaxPageViews != 415 {
		t.Errorf("MaxPageViews = %d, want 415", s.MaxPageViews)
	}
	if s.Participations != 1 {
		t.Errorf("Participations = %d, want 1", s.Participations)
	}
	if s.ParticipationsLevel != 3 {
		t.Errorf("ParticipationsLevel = %d, want 3", s.ParticipationsLevel)
	}
	if s.MaxParticipations != 3 {
		t.Errorf("MaxParticipations = %d, want 3", s.MaxParticipations)
	}

	if s.Tardiness == nil {
		t.Fatal("Tardiness = nil, want populated (tardiness_breakdown field)")
	}
	if s.Tardiness.Total != 5 {
		t.Errorf("Tardiness.Total = %d, want 5", s.Tardiness.Total)
	}
	if s.Tardiness.OnTime != 3 {
		t.Errorf("Tardiness.OnTime = %d, want 3", s.Tardiness.OnTime)
	}
	if s.Tardiness.Late != 0 {
		t.Errorf("Tardiness.Late = %d, want 0", s.Tardiness.Late)
	}
	if s.Tardiness.Missing != 2 {
		t.Errorf("Tardiness.Missing = %d, want 2", s.Tardiness.Missing)
	}
	if s.Tardiness.Floating != 0 {
		t.Errorf("Tardiness.Floating = %d, want 0", s.Tardiness.Floating)
	}
}

// TestAnalyticsLevel_UnmarshalJSON_BareNumber covers a deployment that sends
// the level fields as bare JSON numbers instead of quoted strings.
func TestAnalyticsLevel_UnmarshalJSON_BareNumber(t *testing.T) {
	var s StudentSummary
	if err := json.Unmarshal([]byte(`{"page_views_level":1,"participations_level":3}`), &s); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if s.PageViewsLevel != 1 {
		t.Errorf("PageViewsLevel = %d, want 1", s.PageViewsLevel)
	}
	if s.ParticipationsLevel != 3 {
		t.Errorf("ParticipationsLevel = %d, want 3", s.ParticipationsLevel)
	}
}

// TestAnalyticsLevel_UnmarshalJSON_Null covers a level field sent as JSON
// null, which should decode to the zero value rather than error.
func TestAnalyticsLevel_UnmarshalJSON_Null(t *testing.T) {
	var s StudentSummary
	if err := json.Unmarshal([]byte(`{"page_views_level":null}`), &s); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if s.PageViewsLevel != 0 {
		t.Errorf("PageViewsLevel = %d, want 0", s.PageViewsLevel)
	}
}

// TestAnalyticsLevel_UnmarshalJSON_Malformed asserts that a value which is
// neither a quoted integer, a bare integer, nor null produces an error
// instead of being silently coerced to zero.
func TestAnalyticsLevel_UnmarshalJSON_Malformed(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{"non-numeric string", `{"page_views_level":"abc"}`},
		{"boolean", `{"page_views_level":true}`},
		{"object", `{"page_views_level":{}}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s StudentSummary
			if err := json.Unmarshal([]byte(tt.json), &s); err == nil {
				t.Fatalf("Unmarshal(%s): expected error, got none (PageViewsLevel=%d)", tt.json, s.PageViewsLevel)
			}
		})
	}
}

// TestAnalyticsService_GetStudentSummaries_DocumentedShape decodes the
// documented example through GetStudentSummaries's actual list/collection
// path (Client.GetAllPages), which drops any element that fails to decode
// rather than failing the whole request. Before the fix, the mismatched
// json tags and level field types made every real student summary fail to
// decode, so it was silently dropped and this list came back empty.
func TestAnalyticsService_GetStudentSummaries_DocumentedShape(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/analytics/student_summaries" {
			t.Errorf("Expected path /api/v1/courses/123/analytics/student_summaries, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[" + studentSummaryFixture + "]"))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	service := NewAnalyticsService(client)

	summaries, err := service.GetStudentSummaries(context.Background(), 123, nil)
	if err != nil {
		t.Fatalf("GetStudentSummaries failed: %v", err)
	}

	if len(summaries) != 1 {
		t.Fatalf("expected 1 summary, got %d (element was dropped as undecodable)", len(summaries))
	}

	got := summaries[0]
	if got.ID != 2346 {
		t.Errorf("ID = %d, want 2346", got.ID)
	}
	if got.PageViewsLevel != 1 {
		t.Errorf("PageViewsLevel = %d, want 1", got.PageViewsLevel)
	}
	if got.ParticipationsLevel != 3 {
		t.Errorf("ParticipationsLevel = %d, want 3", got.ParticipationsLevel)
	}
	if got.Tardiness == nil || got.Tardiness.Total != 5 {
		t.Errorf("Tardiness = %+v, want Total=5", got.Tardiness)
	}
}
