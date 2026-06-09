package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRubricsService_ListCourse_WithOpts(t *testing.T) {
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
		if q.Get("page") != "1" {
			t.Errorf("expected page=1, got %q", q.Get("page"))
		}
		if q.Get("per_page") != "10" {
			t.Errorf("expected per_page=10, got %q", q.Get("per_page"))
		}
		json.NewEncoder(w).Encode([]Rubric{{ID: 1, Title: "Rubric With Opts"}})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewRubricsService(client)
	opts := &ListRubricsOptions{Include: []string{"assessments", "associations"}, Page: 1, PerPage: 10}
	rubrics, err := svc.ListCourse(context.Background(), 123, opts)
	if err != nil {
		t.Fatalf("ListCourse: %v", err)
	}
	if len(rubrics) != 1 {
		t.Errorf("expected 1, got %d", len(rubrics))
	}
}

func TestRubricsService_ListAccount_WithOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		q := r.URL.Query()
		if q.Get("page") != "2" {
			t.Errorf("expected page=2, got %q", q.Get("page"))
		}
		if q.Get("per_page") != "5" {
			t.Errorf("expected per_page=5, got %q", q.Get("per_page"))
		}
		json.NewEncoder(w).Encode([]Rubric{{ID: 2, Title: "Account Rubric Paged"}})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewRubricsService(client)
	opts := &ListRubricsOptions{Include: []string{"assessments"}, Page: 2, PerPage: 5}
	rubrics, err := svc.ListAccount(context.Background(), 1, opts)
	if err != nil {
		t.Fatalf("ListAccount: %v", err)
	}
	if len(rubrics) != 1 {
		t.Errorf("expected 1, got %d", len(rubrics))
	}
}

func TestRubricsService_GetAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/rubrics/99" {
			t.Errorf("expected /api/v1/accounts/1/rubrics/99, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(Rubric{ID: 99, Title: "Account Rubric", PointsPossible: 50.0})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewRubricsService(client)
	rubric, err := svc.GetAccount(context.Background(), 1, 99, nil)
	if err != nil {
		t.Fatalf("GetAccount: %v", err)
	}
	if rubric.ID != 99 {
		t.Errorf("expected ID 99, got %d", rubric.ID)
	}
}

func TestRubricsService_GetAccount_WithInclude(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		q := r.URL.Query()
		includes := q["include[]"]
		if len(includes) == 0 {
			t.Error("expected include[] params")
		}
		json.NewEncoder(w).Encode(Rubric{ID: 99, Title: "Account Rubric With Include"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewRubricsService(client)
	rubric, err := svc.GetAccount(context.Background(), 1, 99, []string{"assessments"})
	if err != nil {
		t.Fatalf("GetAccount with include: %v", err)
	}
	if rubric.ID != 99 {
		t.Errorf("expected ID 99, got %d", rubric.ID)
	}
}

func TestRubricsService_GetCourse_WithInclude(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		q := r.URL.Query()
		includes := q["include[]"]
		if len(includes) != 1 {
			t.Errorf("expected 1 include, got %d", len(includes))
		}
		json.NewEncoder(w).Encode(Rubric{ID: 10, Title: "Course Rubric With Include"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewRubricsService(client)
	rubric, err := svc.GetCourse(context.Background(), 123, 10, []string{"associations"})
	if err != nil {
		t.Fatalf("GetCourse with include: %v", err)
	}
	if rubric.ID != 10 {
		t.Errorf("expected ID 10, got %d", rubric.ID)
	}
}

func TestRubricsService_Create_WithCriteria(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		rubricData, ok := body["rubric"].(map[string]interface{})
		if !ok {
			t.Fatal("expected rubric in body")
		}
		if rubricData["free_form_criterion_comments"] != true {
			t.Errorf("expected free_form_criterion_comments=true, got %v", rubricData["free_form_criterion_comments"])
		}
		if rubricData["hide_score_total"] != true {
			t.Errorf("expected hide_score_total=true, got %v", rubricData["hide_score_total"])
		}
		if _, ok := rubricData["criteria"]; !ok {
			t.Error("expected criteria in rubric")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"rubric": Rubric{ID: 200, Title: "Rich Rubric"},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewRubricsService(client)
	params := &CreateRubricParams{
		Title:                     "Rich Rubric",
		PointsPossible:            100.0,
		FreeFormCriterionComments: true,
		HideScoreTotal:            true,
		Criteria: []RubricCriterion{
			{
				Description: "Quality",
				Points:      50,
				Ratings: []RubricRating{
					{Description: "Excellent", Points: 50},
					{Description: "Good", Points: 40},
				},
			},
		},
	}
	rubric, err := svc.Create(context.Background(), 123, params)
	if err != nil {
		t.Fatalf("Create with criteria: %v", err)
	}
	if rubric.ID != 200 {
		t.Errorf("expected ID 200, got %d", rubric.ID)
	}
}

func TestRubricsService_Update_AllParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		rubricData, ok := body["rubric"].(map[string]interface{})
		if !ok {
			t.Fatal("expected rubric in body")
		}
		if rubricData["title"] != "All Updated" {
			t.Errorf("expected title 'All Updated', got %v", rubricData["title"])
		}
		if rubricData["points_possible"].(float64) != 75.0 {
			t.Errorf("expected points_possible 75.0, got %v", rubricData["points_possible"])
		}
		if rubricData["free_form_criterion_comments"] != true {
			t.Errorf("expected free_form_criterion_comments=true")
		}
		if rubricData["hide_score_total"] != true {
			t.Errorf("expected hide_score_total=true")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"rubric": Rubric{ID: 456, Title: "All Updated"},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewRubricsService(client)
	title := "All Updated"
	pts := 75.0
	freeForm := true
	hideTotal := true
	params := &UpdateRubricParams{
		Title:                     &title,
		PointsPossible:            &pts,
		FreeFormCriterionComments: &freeForm,
		HideScoreTotal:            &hideTotal,
	}
	rubric, err := svc.Update(context.Background(), 123, 456, params)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if rubric.Title != "All Updated" {
		t.Errorf("expected 'All Updated', got %s", rubric.Title)
	}
}

func TestRubricsService_Create_NilResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		// Return empty rubric wrapper (rubric field nil)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{"rubric": nil})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewRubricsService(client)
	params := &CreateRubricParams{Title: "Test"}
	_, err = svc.Create(context.Background(), 123, params)
	if err == nil {
		t.Error("expected error when rubric is nil in response")
	}
}

func TestRubricsService_Update_NilResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"rubric": nil})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewRubricsService(client)
	title := "T"
	params := &UpdateRubricParams{Title: &title}
	_, err = svc.Update(context.Background(), 123, 456, params)
	if err == nil {
		t.Error("expected error when rubric is nil in response")
	}
}

func TestRubricsService_Delete_NilResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"rubric": nil})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewRubricsService(client)
	_, err = svc.Delete(context.Background(), 123, 456)
	if err == nil {
		t.Error("expected error when rubric is nil in response")
	}
}

func TestRubricsService_Associate_NilResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"rubric_association": nil})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewRubricsService(client)
	params := &AssociateParams{AssociationType: "Assignment", AssociationID: 1}
	_, err = svc.Associate(context.Background(), 123, 456, params)
	if err == nil {
		t.Error("expected error when rubric_association is nil in response")
	}
}
