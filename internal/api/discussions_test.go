package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDiscussionsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/discussion_topics" {
			t.Errorf("Expected path /api/v1/courses/123/discussion_topics, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"id": 1,
				"title": "Welcome Discussion",
				"message": "<p>Welcome to the course!</p>",
				"published": true,
				"discussion_type": "threaded"
			},
			{
				"id": 2,
				"title": "Week 1 Discussion",
				"message": "<p>Discuss week 1 topics</p>",
				"published": true,
				"discussion_type": "side_comment"
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

	service := NewDiscussionsService(client)
	ctx := context.Background()

	topics, err := service.List(ctx, 123, nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(topics) != 2 {
		t.Errorf("Expected 2 topics, got %d", len(topics))
	}
	if topics[0].Title != "Welcome Discussion" {
		t.Errorf("Expected first topic title 'Welcome Discussion', got %s", topics[0].Title)
	}
}

func TestDiscussionsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/discussion_topics/1" {
			t.Errorf("Expected path /api/v1/courses/123/discussion_topics/1, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 1,
			"title": "Welcome Discussion",
			"message": "<p>Welcome to the course!</p>",
			"published": true,
			"discussion_type": "threaded",
			"discussion_subentry_count": 5
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

	service := NewDiscussionsService(client)
	ctx := context.Background()

	topic, err := service.Get(ctx, 123, 1, nil)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if topic.Title != "Welcome Discussion" {
		t.Errorf("Expected topic title 'Welcome Discussion', got %s", topic.Title)
	}
	if topic.DiscussionSubentryCount != 5 {
		t.Errorf("Expected 5 entries, got %d", topic.DiscussionSubentryCount)
	}
}

func TestDiscussionsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/discussion_topics" {
			t.Errorf("Expected path /api/v1/courses/123/discussion_topics, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if body["title"] != "New Discussion" {
			t.Errorf("Expected title 'New Discussion', got %v", body["title"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 3,
			"title": "New Discussion",
			"message": "<p>Discussion content</p>",
			"published": false,
			"discussion_type": "threaded"
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

	service := NewDiscussionsService(client)
	ctx := context.Background()

	params := &CreateDiscussionParams{
		Title:          "New Discussion",
		Message:        "<p>Discussion content</p>",
		DiscussionType: "threaded",
	}

	topic, err := service.Create(ctx, 123, params)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if topic.Title != "New Discussion" {
		t.Errorf("Expected topic title 'New Discussion', got %s", topic.Title)
	}
}

func TestDiscussionsService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/discussion_topics/1" {
			t.Errorf("Expected path /api/v1/courses/123/discussion_topics/1, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 1,
			"title": "Updated Discussion",
			"message": "<p>Updated content</p>",
			"published": true,
			"discussion_type": "threaded"
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

	service := NewDiscussionsService(client)
	ctx := context.Background()

	title := "Updated Discussion"
	params := &UpdateDiscussionParams{
		Title: &title,
	}

	topic, err := service.Update(ctx, 123, 1, params)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if topic.Title != "Updated Discussion" {
		t.Errorf("Expected topic title 'Updated Discussion', got %s", topic.Title)
	}
}

func TestDiscussionsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/discussion_topics/1" {
			t.Errorf("Expected path /api/v1/courses/123/discussion_topics/1, got %s", r.URL.Path)
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

	service := NewDiscussionsService(client)
	ctx := context.Background()

	err = service.Delete(ctx, 123, 1)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestDiscussionsService_ListEntries(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/discussion_topics/1/entries" {
			t.Errorf("Expected path /api/v1/courses/123/discussion_topics/1/entries, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"id": 1,
				"user_id": 100,
				"message": "<p>First entry</p>",
				"read_state": "read"
			},
			{
				"id": 2,
				"user_id": 101,
				"message": "<p>Second entry</p>",
				"read_state": "unread"
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

	service := NewDiscussionsService(client)
	ctx := context.Background()

	entries, err := service.ListEntries(ctx, 123, 1)
	if err != nil {
		t.Fatalf("ListEntries failed: %v", err)
	}

	if len(entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(entries))
	}
}

func TestDiscussionsService_PostEntry(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/discussion_topics/1/entries" {
			t.Errorf("Expected path /api/v1/courses/123/discussion_topics/1/entries, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 3,
			"user_id": 100,
			"message": "<p>New entry</p>",
			"read_state": "read"
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

	service := NewDiscussionsService(client)
	ctx := context.Background()

	entry, err := service.PostEntry(ctx, 123, 1, "<p>New entry</p>")
	if err != nil {
		t.Fatalf("PostEntry failed: %v", err)
	}

	if entry.Message != "<p>New entry</p>" {
		t.Errorf("Expected message '<p>New entry</p>', got %s", entry.Message)
	}
}

func TestDiscussionsService_PostReply(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/discussion_topics/1/entries/2/replies" {
			t.Errorf("Expected path /api/v1/courses/123/discussion_topics/1/entries/2/replies, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 4,
			"user_id": 100,
			"parent_id": 2,
			"message": "<p>Reply message</p>",
			"read_state": "read"
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

	service := NewDiscussionsService(client)
	ctx := context.Background()

	entry, err := service.PostReply(ctx, 123, 1, 2, "<p>Reply message</p>")
	if err != nil {
		t.Fatalf("PostReply failed: %v", err)
	}

	if entry.ParentID == nil || *entry.ParentID != 2 {
		t.Error("Expected parent_id to be 2")
	}
}

func TestAnnouncementsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/announcements" {
			t.Errorf("Expected path /api/v1/announcements, got %s", r.URL.Path)
		}

		// Check for context_codes parameter
		if !r.URL.Query().Has("context_codes[]") {
			t.Error("Expected context_codes[] parameter")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"id": 1,
				"title": "Course Announcement",
				"message": "<p>Important announcement!</p>",
				"is_announcement": true
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

	service := NewAnnouncementsService(client)
	ctx := context.Background()

	opts := &ListAnnouncementsOptions{
		ContextCodes: []string{"course_123"},
	}

	announcements, err := service.List(ctx, opts)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(announcements) != 1 {
		t.Errorf("Expected 1 announcement, got %d", len(announcements))
	}
	if announcements[0].Title != "Course Announcement" {
		t.Errorf("Expected title 'Course Announcement', got %s", announcements[0].Title)
	}
}

func TestNewDiscussionsService(t *testing.T) {
	client := &Client{}
	service := NewDiscussionsService(client)
	if service == nil {
		t.Error("Expected non-nil service")
	}
	if service.client != client {
		t.Error("Expected service client to match input client")
	}
}

func TestNewAnnouncementsService(t *testing.T) {
	client := &Client{}
	service := NewAnnouncementsService(client)
	if service == nil {
		t.Error("Expected non-nil service")
	}
	if service.client != client {
		t.Error("Expected service client to match input client")
	}
}
