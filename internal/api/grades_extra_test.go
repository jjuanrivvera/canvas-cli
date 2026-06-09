package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGradesService_GetFeed_WithAllOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		q := r.URL.Query()
		if q.Get("user_id") != "77" {
			t.Errorf("expected user_id=77, got %q", q.Get("user_id"))
		}
		if q.Get("assignment_id") != "88" {
			t.Errorf("expected assignment_id=88, got %q", q.Get("assignment_id"))
		}
		if q.Get("start_date") != "2024-01-01" {
			t.Errorf("expected start_date=2024-01-01, got %q", q.Get("start_date"))
		}
		if q.Get("end_date") != "2024-06-30" {
			t.Errorf("expected end_date=2024-06-30, got %q", q.Get("end_date"))
		}
		if q.Get("page") != "1" {
			t.Errorf("expected page=1, got %q", q.Get("page"))
		}
		if q.Get("per_page") != "20" {
			t.Errorf("expected per_page=20, got %q", q.Get("per_page"))
		}
		entries := []GradebookHistoryEntry{
			{ID: 1, UserID: 77, AssignmentID: 88, CurrentGrade: "A", NewGrade: "B"},
		}
		json.NewEncoder(w).Encode(entries)
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewGradesService(client)
	opts := &ListGradebookFeedOptions{
		UserID:       77,
		AssignmentID: 88,
		StartDate:    "2024-01-01",
		EndDate:      "2024-06-30",
		Page:         1,
		PerPage:      20,
	}
	entries, err := svc.GetFeed(context.Background(), 10, opts)
	if err != nil {
		t.Fatalf("GetFeed: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1, got %d", len(entries))
	}
}

func TestGradesService_GetFeed_NilOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		json.NewEncoder(w).Encode([]GradebookHistoryEntry{})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewGradesService(client)
	_, err = svc.GetFeed(context.Background(), 10, nil)
	if err != nil {
		t.Fatalf("GetFeed with nil opts: %v", err)
	}
}

func TestGradesService_ListCustomColumns_WithOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		q := r.URL.Query()
		if q.Get("include_hidden") != "true" {
			t.Errorf("expected include_hidden=true, got %q", q.Get("include_hidden"))
		}
		if q.Get("page") != "3" {
			t.Errorf("expected page=3, got %q", q.Get("page"))
		}
		if q.Get("per_page") != "50" {
			t.Errorf("expected per_page=50, got %q", q.Get("per_page"))
		}
		columns := []CustomGradebookColumn{
			{ID: 1, Title: "Notes", Hidden: true},
		}
		json.NewEncoder(w).Encode(columns)
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewGradesService(client)
	opts := &ListCustomColumnsOptions{IncludeHidden: true, Page: 3, PerPage: 50}
	columns, err := svc.ListCustomColumns(context.Background(), 10, opts)
	if err != nil {
		t.Fatalf("ListCustomColumns: %v", err)
	}
	if len(columns) != 1 {
		t.Errorf("expected 1, got %d", len(columns))
	}
}

func TestGradesService_ListCustomColumns_NilOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		json.NewEncoder(w).Encode([]CustomGradebookColumn{})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewGradesService(client)
	_, err = svc.ListCustomColumns(context.Background(), 10, nil)
	if err != nil {
		t.Fatalf("ListCustomColumns with nil opts: %v", err)
	}
}

func TestGradesService_CreateCustomColumn_AllParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		col, ok := body["column"].(map[string]interface{})
		if !ok {
			t.Fatal("expected column map")
		}
		if col["title"] != "Grade Notes" {
			t.Errorf("expected title 'Grade Notes', got %v", col["title"])
		}
		if col["hidden"] != true {
			t.Errorf("expected hidden=true")
		}
		if col["teacher_notes"] != true {
			t.Errorf("expected teacher_notes=true")
		}
		if col["read_only"] != true {
			t.Errorf("expected read_only=true")
		}
		// position is 5 -> > 0 so it should be set
		if col["position"] == nil {
			t.Errorf("expected position to be set")
		}
		json.NewEncoder(w).Encode(CustomGradebookColumn{ID: 10, Title: "Grade Notes"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewGradesService(client)
	params := &CreateCustomColumnParams{
		Title:        "Grade Notes",
		Position:     5,
		Hidden:       true,
		TeacherNotes: true,
		ReadOnly:     true,
	}
	col, err := svc.CreateCustomColumn(context.Background(), 10, params)
	if err != nil {
		t.Fatalf("CreateCustomColumn: %v", err)
	}
	if col.ID != 10 {
		t.Errorf("expected ID 10, got %d", col.ID)
	}
}

func TestGradesService_UpdateCustomColumn_AllParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		col, ok := body["column"].(map[string]interface{})
		if !ok {
			t.Fatal("expected column map")
		}
		if col["title"] != "Updated Notes" {
			t.Errorf("expected title 'Updated Notes', got %v", col["title"])
		}
		if col["hidden"] != false {
			t.Errorf("expected hidden=false")
		}
		if col["teacher_notes"] != false {
			t.Errorf("expected teacher_notes=false")
		}
		if col["read_only"] != false {
			t.Errorf("expected read_only=false")
		}
		json.NewEncoder(w).Encode(CustomGradebookColumn{ID: 10, Title: "Updated Notes"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewGradesService(client)
	title := "Updated Notes"
	pos := 2
	hidden := false
	teacherNotes := false
	readOnly := false
	params := &UpdateCustomColumnParams{
		Title:        &title,
		Position:     &pos,
		Hidden:       &hidden,
		TeacherNotes: &teacherNotes,
		ReadOnly:     &readOnly,
	}
	col, err := svc.UpdateCustomColumn(context.Background(), 10, 5, params)
	if err != nil {
		t.Fatalf("UpdateCustomColumn: %v", err)
	}
	if col.Title != "Updated Notes" {
		t.Errorf("expected 'Updated Notes', got %s", col.Title)
	}
}

func TestGradesService_GetHistory_WithAllOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		q := r.URL.Query()
		if q.Get("start_date") != "2024-01-01" {
			t.Errorf("expected start_date=2024-01-01, got %q", q.Get("start_date"))
		}
		if q.Get("end_date") != "2024-12-31" {
			t.Errorf("expected end_date=2024-12-31, got %q", q.Get("end_date"))
		}
		if q.Get("page") != "2" {
			t.Errorf("expected page=2, got %q", q.Get("page"))
		}
		if q.Get("per_page") != "30" {
			t.Errorf("expected per_page=30, got %q", q.Get("per_page"))
		}
		days := []GradebookHistoryDay{
			{Date: "2024-06-01", Graders: []GradebookHistoryGrader{{ID: 1, Name: "Prof Jones"}}},
		}
		json.NewEncoder(w).Encode(days)
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewGradesService(client)
	opts := &ListGradebookHistoryOptions{
		StartDate: "2024-01-01",
		EndDate:   "2024-12-31",
		Page:      2,
		PerPage:   30,
	}
	days, err := svc.GetHistory(context.Background(), 10, opts)
	if err != nil {
		t.Fatalf("GetHistory: %v", err)
	}
	if len(days) != 1 {
		t.Errorf("expected 1, got %d", len(days))
	}
}

func TestGradesService_GetHistory_NilOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		json.NewEncoder(w).Encode([]GradebookHistoryDay{})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewGradesService(client)
	_, err = svc.GetHistory(context.Background(), 10, nil)
	if err != nil {
		t.Fatalf("GetHistory with nil opts: %v", err)
	}
}

func TestGradesService_BulkUpdateGrades_WithExcused(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		gradeData, ok := body["grade_data"].(map[string]interface{})
		if !ok {
			t.Fatal("expected grade_data map")
		}
		if len(gradeData) == 0 {
			t.Error("expected non-empty grade_data")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewGradesService(client)
	grades := []BulkUpdateGrade{
		{StudentID: 10, AssignmentID: 5, Grade: "A+", Excused: false},
		{StudentID: 20, AssignmentID: 5, Grade: "", Excused: true},
	}
	if err := svc.BulkUpdateGrades(context.Background(), 99, grades); err != nil {
		t.Fatalf("BulkUpdateGrades: %v", err)
	}
}

func TestGradesService_DeleteCustomColumn_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		json.NewEncoder(w).Encode(CustomGradebookColumn{ID: 5, Title: "Deleted Col"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewGradesService(client)
	col, err := svc.DeleteCustomColumn(context.Background(), 10, 5)
	if err != nil {
		t.Fatalf("DeleteCustomColumn: %v", err)
	}
	if col.ID != 5 {
		t.Errorf("expected ID 5, got %d", col.ID)
	}
}

func TestGradesService_SetCustomColumnData_SetsColumnID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		// Return a datum without column_id to test the fallback assignment
		json.NewEncoder(w).Encode(CustomColumnDatum{UserID: 7, Content: "Good progress"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewGradesService(client)
	datum, err := svc.SetCustomColumnData(context.Background(), 10, 42, 7, "Good progress")
	if err != nil {
		t.Fatalf("SetCustomColumnData: %v", err)
	}
	// Column ID should be set from the columnID argument since API doesn't return it
	if datum.ColumnID != 42 {
		t.Errorf("expected ColumnID 42, got %d", datum.ColumnID)
	}
}

func TestGradesService_GetCustomColumnData_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		data := []CustomColumnDatum{
			{ColumnID: 5, UserID: 1, Content: "note A"},
			{ColumnID: 5, UserID: 2, Content: "note B"},
		}
		json.NewEncoder(w).Encode(data)
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewGradesService(client)
	data, err := svc.GetCustomColumnData(context.Background(), 10, 5)
	if err != nil {
		t.Fatalf("GetCustomColumnData: %v", err)
	}
	if len(data) != 2 {
		t.Errorf("expected 2, got %d", len(data))
	}
}
