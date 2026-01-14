package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConversationsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/conversations" {
			t.Errorf("expected /api/v1/conversations, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Conversation{
			{ID: 1, Subject: "Hello", MessageCount: 3},
			{ID: 2, Subject: "Question", MessageCount: 1},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewConversationsService(client)
	conversations, err := service.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(conversations) != 2 {
		t.Errorf("expected 2 conversations, got %d", len(conversations))
	}

	if conversations[0].Subject != "Hello" {
		t.Errorf("expected 'Hello', got %s", conversations[0].Subject)
	}
}

func TestConversationsService_List_WithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Query().Get("scope") != "unread" {
			t.Errorf("expected scope=unread, got %s", r.URL.Query().Get("scope"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Conversation{
			{ID: 1, Subject: "Unread Message"},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewConversationsService(client)
	opts := &ListConversationsOptions{
		Scope: "unread",
	}
	conversations, err := service.List(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(conversations) != 1 {
		t.Errorf("expected 1 conversation, got %d", len(conversations))
	}
}

func TestConversationsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/conversations/123" {
			t.Errorf("expected /api/v1/conversations/123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Conversation{
			ID:           123,
			Subject:      "Test Conversation",
			MessageCount: 5,
			Messages: []ConversationMessage{
				{ID: 1, Body: "Hello"},
				{ID: 2, Body: "Hi there"},
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewConversationsService(client)
	conversation, err := service.Get(context.Background(), 123, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if conversation.ID != 123 {
		t.Errorf("expected ID 123, got %d", conversation.ID)
	}

	if len(conversation.Messages) != 2 {
		t.Errorf("expected 2 messages, got %d", len(conversation.Messages))
	}
}

func TestConversationsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if body["subject"] != "New Conversation" {
			t.Errorf("expected subject 'New Conversation', got %v", body["subject"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Conversation{
			{ID: 456, Subject: "New Conversation"},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewConversationsService(client)
	params := &CreateConversationParams{
		Recipients: []string{"123", "456"},
		Subject:    "New Conversation",
		Body:       "Hello, this is a test message.",
	}

	conversations, err := service.Create(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(conversations) != 1 {
		t.Errorf("expected 1 conversation, got %d", len(conversations))
	}

	if conversations[0].ID != 456 {
		t.Errorf("expected ID 456, got %d", conversations[0].ID)
	}
}

func TestConversationsService_AddMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/conversations/123/add_message" {
			t.Errorf("expected /api/v1/conversations/123/add_message, got %s", r.URL.Path)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if body["body"] != "Reply message" {
			t.Errorf("expected body 'Reply message', got %v", body["body"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Conversation{
			ID:           123,
			MessageCount: 6,
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewConversationsService(client)
	params := &AddMessageParams{
		Body: "Reply message",
	}

	conversation, err := service.AddMessage(context.Background(), 123, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if conversation.MessageCount != 6 {
		t.Errorf("expected 6 messages, got %d", conversation.MessageCount)
	}
}

func TestConversationsService_AddRecipients(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/conversations/123/add_recipients" {
			t.Errorf("expected /api/v1/conversations/123/add_recipients, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Conversation{
			ID:       123,
			Audience: []int64{1, 2, 3, 789},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewConversationsService(client)
	conversation, err := service.AddRecipients(context.Background(), 123, []string{"789"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(conversation.Audience) != 4 {
		t.Errorf("expected 4 audience members, got %d", len(conversation.Audience))
	}
}

func TestConversationsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Conversation{
			ID: 123,
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewConversationsService(client)
	conversation, err := service.Delete(context.Background(), 123)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if conversation.ID != 123 {
		t.Errorf("expected ID 123, got %d", conversation.ID)
	}
}

func TestConversationsService_GetUnreadCount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/conversations/unread_count" {
			t.Errorf("expected /api/v1/conversations/unread_count, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(UnreadCount{
			UnreadCount: "5",
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewConversationsService(client)
	count, err := service.GetUnreadCount(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if count.UnreadCount != "5" {
		t.Errorf("expected unread count '5', got %s", count.UnreadCount)
	}
}

func TestConversationsService_Archive(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		conv, ok := body["conversation"].(map[string]interface{})
		if !ok {
			t.Fatal("expected conversation in body")
		}

		if conv["workflow_state"] != "archived" {
			t.Errorf("expected workflow_state 'archived', got %v", conv["workflow_state"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Conversation{
			ID:            123,
			WorkflowState: "archived",
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewConversationsService(client)
	conversation, err := service.Archive(context.Background(), 123)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if conversation.WorkflowState != "archived" {
		t.Errorf("expected workflow_state 'archived', got %s", conversation.WorkflowState)
	}
}

func TestConversationsService_Star(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		conv, ok := body["conversation"].(map[string]interface{})
		if !ok {
			t.Fatal("expected conversation in body")
		}

		if conv["starred"] != true {
			t.Errorf("expected starred true, got %v", conv["starred"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Conversation{
			ID:      123,
			Starred: true,
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewConversationsService(client)
	conversation, err := service.Star(context.Background(), 123)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !conversation.Starred {
		t.Error("expected starred to be true")
	}
}

func TestConversationsService_BatchUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if body["event"] != "mark_as_read" {
			t.Errorf("expected event 'mark_as_read', got %v", body["event"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(BatchUpdateProgress{
			Progress: &ConversationProgress{
				ID:            1,
				WorkflowState: "completed",
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewConversationsService(client)
	params := &BatchUpdateParams{
		ConversationIDs: []int64{1, 2, 3},
		Event:           "mark_as_read",
	}

	progress, err := service.BatchUpdate(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if progress.Progress == nil {
		t.Fatal("expected progress in response")
	}

	if progress.Progress.WorkflowState != "completed" {
		t.Errorf("expected workflow_state 'completed', got %s", progress.Progress.WorkflowState)
	}
}

func TestNewConversationsService(t *testing.T) {
	client, err := NewClient(ClientConfig{
		BaseURL: "https://canvas.example.com",
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewConversationsService(client)
	if service == nil {
		t.Fatal("expected non-nil service")
	}
	if service.client != client {
		t.Error("expected client to be set")
	}
}

func TestParseRecipients(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"", nil},
		{"123", []string{"123"}},
		{"123,456", []string{"123", "456"}},
		{"123, 456, 789", []string{"123", "456", "789"}},
		{"course_123,456", []string{"course_123", "456"}},
	}

	for _, test := range tests {
		result := ParseRecipients(test.input)
		if len(result) != len(test.expected) {
			t.Errorf("ParseRecipients(%q): expected %d items, got %d", test.input, len(test.expected), len(result))
			continue
		}
		for i, v := range result {
			if v != test.expected[i] {
				t.Errorf("ParseRecipients(%q)[%d]: expected %q, got %q", test.input, i, test.expected[i], v)
			}
		}
	}
}

func TestFormatRecipients(t *testing.T) {
	ids := []int64{123, 456, 789}
	result := FormatRecipients(ids)
	expected := []string{"123", "456", "789"}

	if len(result) != len(expected) {
		t.Errorf("expected %d items, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("expected %q at index %d, got %q", expected[i], i, v)
		}
	}
}
