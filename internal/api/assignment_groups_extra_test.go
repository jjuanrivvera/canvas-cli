package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAssignmentGroupsService_List_WithAllOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		q := r.URL.Query()
		includes := q["include[]"]
		if len(includes) != 2 {
			t.Errorf("expected 2 includes, got %d", len(includes))
		}
		ids := q["assignment_ids[]"]
		if len(ids) != 2 {
			t.Errorf("expected 2 assignment_ids, got %d", len(ids))
		}
		exclude := q["exclude_assignment_submission_types[]"]
		if len(exclude) != 1 || exclude[0] != "online_quiz" {
			t.Errorf("expected exclude online_quiz, got %v", exclude)
		}
		if q.Get("override_assignment_dates") != "true" {
			t.Errorf("expected override_assignment_dates=true, got %q", q.Get("override_assignment_dates"))
		}
		if q.Get("grading_period_id") != "3" {
			t.Errorf("expected grading_period_id=3, got %q", q.Get("grading_period_id"))
		}
		if q.Get("scope_assignments_to_student") != "false" {
			t.Errorf("expected scope_assignments_to_student=false, got %q", q.Get("scope_assignments_to_student"))
		}
		if q.Get("page") != "1" {
			t.Errorf("expected page=1, got %q", q.Get("page"))
		}
		if q.Get("per_page") != "10" {
			t.Errorf("expected per_page=10, got %q", q.Get("per_page"))
		}
		groups := []AssignmentGroup{
			{ID: 1, Name: "Homework", GroupWeight: 30.0},
		}
		json.NewEncoder(w).Encode(groups)
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewAssignmentGroupsService(client)
	override := true
	scope := false
	opts := &ListAssignmentGroupsOptions{
		Include:                          []string{"assignments", "submission"},
		AssignmentIDs:                    []int64{10, 20},
		ExcludeAssignmentSubmissionTypes: []string{"online_quiz"},
		OverrideAssignmentDates:          &override,
		GradingPeriodID:                  3,
		ScopeAssignmentsToStudent:        &scope,
		Page:                             1,
		PerPage:                          10,
	}
	groups, err := svc.List(context.Background(), 5, opts)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(groups) != 1 {
		t.Errorf("expected 1, got %d", len(groups))
	}
}

func TestAssignmentGroupsService_List_NilOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		json.NewEncoder(w).Encode([]AssignmentGroup{})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewAssignmentGroupsService(client)
	_, err = svc.List(context.Background(), 5, nil)
	if err != nil {
		t.Fatalf("List nil opts: %v", err)
	}
}

func TestAssignmentGroupsService_Get_WithInclude(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		includes := r.URL.Query()["include[]"]
		if len(includes) == 0 {
			t.Errorf("expected includes to be set")
		}
		json.NewEncoder(w).Encode(AssignmentGroup{ID: 7, Name: "Quizzes", GroupWeight: 20.0})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewAssignmentGroupsService(client)
	group, err := svc.Get(context.Background(), 5, 7, []string{"assignments", "discussion_topic"})
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if group.ID != 7 {
		t.Errorf("expected ID 7, got %d", group.ID)
	}
}

func TestAssignmentGroupsService_Update_WithAllParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body["name"] != "Revised Group" {
			t.Errorf("expected name 'Revised Group', got %v", body["name"])
		}
		if body["group_weight"] == nil {
			t.Error("expected group_weight to be set")
		}
		if body["sis_source_id"] != "SIS42" {
			t.Errorf("expected sis_source_id 'SIS42', got %v", body["sis_source_id"])
		}
		rules, ok := body["rules"].(map[string]interface{})
		if !ok {
			t.Fatal("expected rules map")
		}
		if rules["drop_highest"] == nil {
			t.Error("expected drop_highest to be set")
		}
		if rules["never_drop"] == nil {
			t.Error("expected never_drop to be set")
		}
		json.NewEncoder(w).Encode(AssignmentGroup{ID: 7, Name: "Revised Group"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewAssignmentGroupsService(client)
	name := "Revised Group"
	pos := 2
	weight := 25.5
	sisID := "SIS42"
	params := &UpdateAssignmentGroupParams{
		Name:        &name,
		Position:    &pos,
		GroupWeight: &weight,
		SISSourceID: &sisID,
		IntegrationData: map[string]interface{}{
			"key": "value",
		},
		Rules: &GradingRules{
			DropHighest: 1,
			NeverDrop:   []int64{100, 200},
		},
	}
	group, err := svc.Update(context.Background(), 5, 7, params)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if group.Name != "Revised Group" {
		t.Errorf("expected 'Revised Group', got %s", group.Name)
	}
}

func TestAssignmentGroupsService_Update_EmptyRules(t *testing.T) {
	// Test branch where rules.DropLowest/DropHighest are both 0 and NeverDrop is empty
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		// rules should still be present (empty map) since Rules != nil
		if _, ok := body["rules"]; !ok {
			t.Error("expected rules key to be present in body")
		}
		json.NewEncoder(w).Encode(AssignmentGroup{ID: 7, Name: "Group"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewAssignmentGroupsService(client)
	name := "Group"
	params := &UpdateAssignmentGroupParams{
		Name:  &name,
		Rules: &GradingRules{}, // empty rules, no drop values
	}
	_, err = svc.Update(context.Background(), 5, 7, params)
	if err != nil {
		t.Fatalf("Update with empty rules: %v", err)
	}
}
