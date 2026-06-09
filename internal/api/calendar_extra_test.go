package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalendarService_ListForUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/users/42/calendar_events" {
			t.Errorf("expected /api/v1/users/42/calendar_events, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]CalendarEvent{
			{ID: 1, Title: "User Event", ContextCode: "course_10", WorkflowState: "active"},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewCalendarService(client)
	events, err := svc.ListForUser(context.Background(), 42, nil)
	if err != nil {
		t.Fatalf("ListForUser: %v", err)
	}
	if len(events) != 1 {
		t.Errorf("expected 1 event, got %d", len(events))
	}
	if events[0].Title != "User Event" {
		t.Errorf("expected 'User Event', got %s", events[0].Title)
	}
}

func TestCalendarService_ListForUser_WithOpts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		q := r.URL.Query()
		if q.Get("type") != "event" {
			t.Errorf("expected type=event, got %q", q.Get("type"))
		}
		if q.Get("start_date") != "2024-01-01" {
			t.Errorf("expected start_date=2024-01-01, got %q", q.Get("start_date"))
		}
		if q.Get("end_date") != "2024-12-31" {
			t.Errorf("expected end_date=2024-12-31, got %q", q.Get("end_date"))
		}
		if q.Get("all_events") != "true" {
			t.Errorf("expected all_events=true, got %q", q.Get("all_events"))
		}
		if q.Get("undated") != "true" {
			t.Errorf("expected undated=true, got %q", q.Get("undated"))
		}
		if q.Get("important_dates") != "true" {
			t.Errorf("expected important_dates=true, got %q", q.Get("important_dates"))
		}
		if q.Get("blackout_date") != "true" {
			t.Errorf("expected blackout_date=true, got %q", q.Get("blackout_date"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]CalendarEvent{{ID: 2, Title: "Filtered Event"}})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewCalendarService(client)
	opts := &ListCalendarEventsOptions{
		Type:           "event",
		StartDate:      "2024-01-01",
		EndDate:        "2024-12-31",
		Undated:        true,
		AllEvents:      true,
		ContextCodes:   []string{"course_10"},
		Excludes:       []string{"description"},
		Includes:       []string{"web_conference"},
		ImportantDates: true,
		BlackoutDate:   true,
		Page:           1,
		PerPage:        20,
	}
	events, err := svc.ListForUser(context.Background(), 42, opts)
	if err != nil {
		t.Fatalf("ListForUser with opts: %v", err)
	}
	if len(events) != 1 {
		t.Errorf("expected 1 event, got %d", len(events))
	}
}

func TestCalendarService_Create_AllParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		evData := body["calendar_event"].(map[string]interface{})
		if evData["title"] != "All Params Event" {
			t.Errorf("expected title 'All Params Event', got %v", evData["title"])
		}
		if evData["description"] != "desc" {
			t.Errorf("expected description 'desc', got %v", evData["description"])
		}
		if evData["location_name"] != "Room A" {
			t.Errorf("expected location_name 'Room A', got %v", evData["location_name"])
		}
		if evData["location_address"] != "123 Main St" {
			t.Errorf("expected location_address, got %v", evData["location_address"])
		}
		if evData["time_zone_edited"] != "US/Eastern" {
			t.Errorf("expected time_zone_edited, got %v", evData["time_zone_edited"])
		}
		if evData["all_day"] != true {
			t.Errorf("expected all_day=true, got %v", evData["all_day"])
		}
		if evData["rrule"] != "FREQ=WEEKLY" {
			t.Errorf("expected rrule, got %v", evData["rrule"])
		}
		if evData["blackout_date"] != true {
			t.Errorf("expected blackout_date=true, got %v", evData["blackout_date"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CalendarEvent{ID: 99, Title: "All Params Event"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewCalendarService(client)
	params := &CreateCalendarEventParams{
		ContextCode:     "course_10",
		Title:           "All Params Event",
		Description:     "desc",
		StartAt:         "2024-06-01T09:00:00Z",
		EndAt:           "2024-06-01T10:00:00Z",
		LocationName:    "Room A",
		LocationAddress: "123 Main St",
		TimeZoneEdited:  "US/Eastern",
		AllDay:          true,
		RRule:           "FREQ=WEEKLY",
		BlackoutDate:    true,
	}
	event, err := svc.Create(context.Background(), params)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if event.ID != 99 {
		t.Errorf("expected ID 99, got %d", event.ID)
	}
}

func TestCalendarService_Update_AllParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		// which=all should be in query
		if r.URL.Query().Get("which") != "all" {
			t.Errorf("expected which=all, got %q", r.URL.Query().Get("which"))
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		evData := body["calendar_event"].(map[string]interface{})
		if evData["title"] != "Updated All" {
			t.Errorf("expected title, got %v", evData["title"])
		}
		if evData["description"] != "new desc" {
			t.Errorf("expected description, got %v", evData["description"])
		}
		if evData["context_code"] != "course_5" {
			t.Errorf("expected context_code, got %v", evData["context_code"])
		}
		if evData["location_name"] != "New Room" {
			t.Errorf("expected location_name, got %v", evData["location_name"])
		}
		if evData["location_address"] != "456 Main St" {
			t.Errorf("expected location_address, got %v", evData["location_address"])
		}
		if evData["time_zone_edited"] != "US/Pacific" {
			t.Errorf("expected time_zone_edited, got %v", evData["time_zone_edited"])
		}
		if evData["all_day"] != true {
			t.Errorf("expected all_day=true, got %v", evData["all_day"])
		}
		if evData["rrule"] != "FREQ=DAILY" {
			t.Errorf("expected rrule, got %v", evData["rrule"])
		}
		if evData["blackout_date"] != true {
			t.Errorf("expected blackout_date=true, got %v", evData["blackout_date"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CalendarEvent{ID: 5, Title: "Updated All"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewCalendarService(client)

	title := "Updated All"
	desc := "new desc"
	ctx := "course_5"
	loc := "New Room"
	addr := "456 Main St"
	tz := "US/Pacific"
	allDay := true
	rrule := "FREQ=DAILY"
	bd := true
	startAt := "2024-07-01T09:00:00Z"
	endAt := "2024-07-01T10:00:00Z"
	which := "all"

	params := &UpdateCalendarEventParams{
		Title:           &title,
		Description:     &desc,
		ContextCode:     &ctx,
		LocationName:    &loc,
		LocationAddress: &addr,
		TimeZoneEdited:  &tz,
		AllDay:          &allDay,
		RRule:           &rrule,
		BlackoutDate:    &bd,
		StartAt:         &startAt,
		EndAt:           &endAt,
		Which:           which,
	}
	event, err := svc.Update(context.Background(), 5, params)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if event.Title != "Updated All" {
		t.Errorf("expected 'Updated All', got %s", event.Title)
	}
}

func TestCalendarService_Reserve_WithParticipant(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/calendar_events/1/reservations/99" {
			t.Errorf("expected path with participant, got %s", r.URL.Path)
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["comments"] != "see you there" {
			t.Errorf("expected comments, got %v", body["comments"])
		}
		if body["cancel_existing"] != true {
			t.Errorf("expected cancel_existing=true, got %v", body["cancel_existing"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CalendarEvent{ID: 1, OwnReservation: true})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewCalendarService(client)
	participantID := int64(99)
	event, err := svc.Reserve(context.Background(), 1, &participantID, "see you there", true)
	if err != nil {
		t.Fatalf("Reserve: %v", err)
	}
	if !event.OwnReservation {
		t.Error("expected OwnReservation=true")
	}
}

func TestCalendarService_ListWithAllOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		q := r.URL.Query()
		if q.Get("undated") != "true" {
			t.Errorf("expected undated=true, got %q", q.Get("undated"))
		}
		if q.Get("all_events") != "true" {
			t.Errorf("expected all_events=true, got %q", q.Get("all_events"))
		}
		if q.Get("important_dates") != "true" {
			t.Errorf("expected important_dates=true, got %q", q.Get("important_dates"))
		}
		if q.Get("blackout_date") != "true" {
			t.Errorf("expected blackout_date=true, got %q", q.Get("blackout_date"))
		}
		if q.Get("page") != "1" {
			t.Errorf("expected page=1, got %q", q.Get("page"))
		}
		if q.Get("per_page") != "10" {
			t.Errorf("expected per_page=10, got %q", q.Get("per_page"))
		}
		json.NewEncoder(w).Encode([]CalendarEvent{{ID: 1, Title: "Test"}})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewCalendarService(client)
	opts := &ListCalendarEventsOptions{
		Undated:        true,
		AllEvents:      true,
		ImportantDates: true,
		BlackoutDate:   true,
		Excludes:       []string{"description"},
		Includes:       []string{"web_conference"},
		Page:           1,
		PerPage:        10,
	}
	events, err := svc.List(context.Background(), opts)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(events) != 1 {
		t.Errorf("expected 1, got %d", len(events))
	}
}

func TestCalendarService_Get_Error(t *testing.T) {
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
	svc := NewCalendarService(client)
	_, err = svc.Get(context.Background(), 9999)
	if err == nil {
		t.Error("expected error, got nil")
	}
}
