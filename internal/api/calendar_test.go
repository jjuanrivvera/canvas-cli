package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalendarService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/calendar_events" {
			t.Errorf("Expected path /api/v1/calendar_events, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"id": 1,
				"title": "Team Meeting",
				"context_code": "course_123",
				"workflow_state": "active",
				"all_day": false
			},
			{
				"id": 2,
				"title": "Project Deadline",
				"context_code": "course_123",
				"workflow_state": "active",
				"all_day": true
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

	service := NewCalendarService(client)
	ctx := context.Background()

	events, err := service.List(ctx, nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(events) != 2 {
		t.Errorf("Expected 2 events, got %d", len(events))
	}
	if events[0].Title != "Team Meeting" {
		t.Errorf("Expected first event title 'Team Meeting', got %s", events[0].Title)
	}
}

func TestCalendarService_ListWithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		// Check for context_codes parameter
		if !r.URL.Query().Has("context_codes[]") {
			t.Error("Expected context_codes[] parameter")
		}

		if r.URL.Query().Get("type") != "event" {
			t.Errorf("Expected type=event, got %s", r.URL.Query().Get("type"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id": 1, "title": "Event"}]`))
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

	service := NewCalendarService(client)
	ctx := context.Background()

	opts := &ListCalendarEventsOptions{
		Type:         "event",
		ContextCodes: []string{"course_123"},
		StartDate:    "2024-01-01",
		EndDate:      "2024-12-31",
	}

	events, err := service.List(ctx, opts)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}
}

func TestCalendarService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/calendar_events/1" {
			t.Errorf("Expected path /api/v1/calendar_events/1, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 1,
			"title": "Team Meeting",
			"description": "<p>Weekly sync</p>",
			"context_code": "course_123",
			"workflow_state": "active",
			"location_name": "Room 101"
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

	service := NewCalendarService(client)
	ctx := context.Background()

	event, err := service.Get(ctx, 1)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if event.Title != "Team Meeting" {
		t.Errorf("Expected event title 'Team Meeting', got %s", event.Title)
	}
	if event.LocationName != "Room 101" {
		t.Errorf("Expected location 'Room 101', got %s", event.LocationName)
	}
}

func TestCalendarService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/calendar_events" {
			t.Errorf("Expected path /api/v1/calendar_events, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		eventData, ok := body["calendar_event"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected 'calendar_event' key in request body")
		}

		if eventData["context_code"] != "course_123" {
			t.Errorf("Expected context_code 'course_123', got %v", eventData["context_code"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 3,
			"title": "New Event",
			"context_code": "course_123",
			"workflow_state": "active"
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

	service := NewCalendarService(client)
	ctx := context.Background()

	params := &CreateCalendarEventParams{
		ContextCode: "course_123",
		Title:       "New Event",
		StartAt:     "2024-07-19T15:00:00Z",
		EndAt:       "2024-07-19T16:00:00Z",
	}

	event, err := service.Create(ctx, params)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if event.Title != "New Event" {
		t.Errorf("Expected event title 'New Event', got %s", event.Title)
	}
}

func TestCalendarService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/calendar_events/1" {
			t.Errorf("Expected path /api/v1/calendar_events/1, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 1,
			"title": "Updated Event",
			"context_code": "course_123",
			"workflow_state": "active"
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

	service := NewCalendarService(client)
	ctx := context.Background()

	title := "Updated Event"
	params := &UpdateCalendarEventParams{
		Title: &title,
	}

	event, err := service.Update(ctx, 1, params)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if event.Title != "Updated Event" {
		t.Errorf("Expected event title 'Updated Event', got %s", event.Title)
	}
}

func TestCalendarService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/calendar_events/1" {
			t.Errorf("Expected path /api/v1/calendar_events/1, got %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
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

	service := NewCalendarService(client)
	ctx := context.Background()

	err = service.Delete(ctx, 1, nil)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestCalendarService_DeleteWithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Query().Get("which") != "all" {
			t.Errorf("Expected which=all, got %s", r.URL.Query().Get("which"))
		}

		w.WriteHeader(http.StatusOK)
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

	service := NewCalendarService(client)
	ctx := context.Background()

	opts := &DeleteOptions{
		Which:        "all",
		CancelReason: "Series cancelled",
	}

	err = service.Delete(ctx, 1, opts)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestCalendarService_Reserve(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/calendar_events/1/reservations" {
			t.Errorf("Expected path /api/v1/calendar_events/1/reservations, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 10,
			"title": "Reserved Slot",
			"context_code": "user_123",
			"workflow_state": "active",
			"own_reservation": true
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

	service := NewCalendarService(client)
	ctx := context.Background()

	event, err := service.Reserve(ctx, 1, nil, "Looking forward to it", false)
	if err != nil {
		t.Fatalf("Reserve failed: %v", err)
	}

	if !event.OwnReservation {
		t.Error("Expected own_reservation to be true")
	}
}

func TestNewCalendarService(t *testing.T) {
	client := &Client{}
	service := NewCalendarService(client)
	if service == nil {
		t.Fatal("Expected non-nil service")
	}
	if service.client != client {
		t.Error("Expected service client to match input client")
	}
}
