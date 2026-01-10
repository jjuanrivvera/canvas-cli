package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPagesService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/pages" {
			t.Errorf("Expected path /api/v1/courses/123/pages, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"page_id": 1,
				"url": "welcome",
				"title": "Welcome",
				"published": true,
				"front_page": true
			},
			{
				"page_id": 2,
				"url": "syllabus",
				"title": "Syllabus",
				"published": true,
				"front_page": false
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

	service := NewPagesService(client)
	ctx := context.Background()

	pages, err := service.List(ctx, 123, nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(pages) != 2 {
		t.Errorf("Expected 2 pages, got %d", len(pages))
	}
	if pages[0].Title != "Welcome" {
		t.Errorf("Expected first page title 'Welcome', got %s", pages[0].Title)
	}
	if !pages[0].FrontPage {
		t.Error("Expected first page to be front page")
	}
}

func TestPagesService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/pages/welcome" {
			t.Errorf("Expected path /api/v1/courses/123/pages/welcome, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"page_id": 1,
			"url": "welcome",
			"title": "Welcome",
			"body": "<p>Welcome to the course!</p>",
			"published": true,
			"front_page": true
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

	service := NewPagesService(client)
	ctx := context.Background()

	page, err := service.Get(ctx, 123, "welcome")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if page.Title != "Welcome" {
		t.Errorf("Expected page title 'Welcome', got %s", page.Title)
	}
	if page.Body != "<p>Welcome to the course!</p>" {
		t.Errorf("Expected page body '<p>Welcome to the course!</p>', got %s", page.Body)
	}
}

func TestPagesService_GetFrontPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/front_page" {
			t.Errorf("Expected path /api/v1/courses/123/front_page, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"page_id": 1,
			"url": "home",
			"title": "Home",
			"published": true,
			"front_page": true
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

	service := NewPagesService(client)
	ctx := context.Background()

	page, err := service.GetFrontPage(ctx, 123)
	if err != nil {
		t.Fatalf("GetFrontPage failed: %v", err)
	}

	if !page.FrontPage {
		t.Error("Expected front page to be true")
	}
}

func TestPagesService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/pages" {
			t.Errorf("Expected path /api/v1/courses/123/pages, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		pageData, ok := body["wiki_page"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected 'wiki_page' key in request body")
		}

		if pageData["title"] != "New Page" {
			t.Errorf("Expected page title 'New Page', got %v", pageData["title"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"page_id": 3,
			"url": "new-page",
			"title": "New Page",
			"published": false
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

	service := NewPagesService(client)
	ctx := context.Background()

	params := &CreatePageParams{
		Title: "New Page",
	}

	page, err := service.Create(ctx, 123, params)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if page.Title != "New Page" {
		t.Errorf("Expected page title 'New Page', got %s", page.Title)
	}
}

func TestPagesService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/pages/my-page" {
			t.Errorf("Expected path /api/v1/courses/123/pages/my-page, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"page_id": 1,
			"url": "my-page",
			"title": "Updated Title",
			"published": true
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

	service := NewPagesService(client)
	ctx := context.Background()

	title := "Updated Title"
	params := &UpdatePageParams{
		Title: &title,
	}

	page, err := service.Update(ctx, 123, "my-page", params)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if page.Title != "Updated Title" {
		t.Errorf("Expected page title 'Updated Title', got %s", page.Title)
	}
}

func TestPagesService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/pages/my-page" {
			t.Errorf("Expected path /api/v1/courses/123/pages/my-page, got %s", r.URL.Path)
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

	service := NewPagesService(client)
	ctx := context.Background()

	err = service.Delete(ctx, 123, "my-page")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestPagesService_ListRevisions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/pages/my-page/revisions" {
			t.Errorf("Expected path /api/v1/courses/123/pages/my-page/revisions, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"revision_id": 3,
				"latest": true
			},
			{
				"revision_id": 2,
				"latest": false
			},
			{
				"revision_id": 1,
				"latest": false
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

	service := NewPagesService(client)
	ctx := context.Background()

	revisions, err := service.ListRevisions(ctx, 123, "my-page")
	if err != nil {
		t.Fatalf("ListRevisions failed: %v", err)
	}

	if len(revisions) != 3 {
		t.Errorf("Expected 3 revisions, got %d", len(revisions))
	}
	if !revisions[0].Latest {
		t.Error("Expected first revision to be latest")
	}
}

func TestNewPagesService(t *testing.T) {
	client := &Client{}
	service := NewPagesService(client)
	if service == nil {
		t.Fatal("Expected non-nil service")
	}
	if service.client != client {
		t.Error("Expected service client to match input client")
	}
}
