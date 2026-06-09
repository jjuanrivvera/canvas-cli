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

func TestPagesService_UpdateFrontPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/front_page" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Page{Title: "Front Page", URL: "front-page"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewPagesService(client)
	title := "Front Page"
	page, err := svc.UpdateFrontPage(context.Background(), 1, &UpdatePageParams{Title: &title})
	if err != nil {
		t.Fatalf("UpdateFrontPage: %v", err)
	}
	if page.Title != "Front Page" {
		t.Errorf("expected 'Front Page', got %s", page.Title)
	}
}

func TestPagesService_Duplicate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/pages/intro/duplicate" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Page{Title: "Copy of intro", URL: "copy-of-intro"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewPagesService(client)
	page, err := svc.Duplicate(context.Background(), 1, "intro")
	if err != nil {
		t.Fatalf("Duplicate: %v", err)
	}
	if page.URL != "copy-of-intro" {
		t.Errorf("expected copy-of-intro, got %s", page.URL)
	}
}

func TestPagesService_GetRevision(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/pages/intro/revisions/3" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(PageRevision{RevisionID: 3})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewPagesService(client)
	rev, err := svc.GetRevision(context.Background(), 1, "intro", 3, false)
	if err != nil {
		t.Fatalf("GetRevision: %v", err)
	}
	if rev.RevisionID != 3 {
		t.Errorf("expected RevisionID 3, got %d", rev.RevisionID)
	}
}

func TestPagesService_GetRevision_Summary(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Query().Get("summary") != "1" {
			t.Errorf("expected summary=1")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(PageRevision{RevisionID: 4})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewPagesService(client)
	rev, err := svc.GetRevision(context.Background(), 1, "intro", 4, true)
	if err != nil {
		t.Fatalf("GetRevision summary: %v", err)
	}
	if rev.RevisionID != 4 {
		t.Errorf("expected RevisionID 4, got %d", rev.RevisionID)
	}
}

func TestPagesService_GetLatestRevision(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/pages/intro/revisions/latest" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(PageRevision{RevisionID: 10})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewPagesService(client)
	rev, err := svc.GetLatestRevision(context.Background(), 1, "intro", false)
	if err != nil {
		t.Fatalf("GetLatestRevision: %v", err)
	}
	if rev.RevisionID != 10 {
		t.Errorf("expected RevisionID 10, got %d", rev.RevisionID)
	}
}

func TestPagesService_GetLatestRevision_Summary(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Query().Get("summary") != "1" {
			t.Errorf("expected summary=1")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(PageRevision{RevisionID: 11})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewPagesService(client)
	rev, err := svc.GetLatestRevision(context.Background(), 1, "intro", true)
	if err != nil {
		t.Fatalf("GetLatestRevision summary: %v", err)
	}
	if rev.RevisionID != 11 {
		t.Errorf("expected RevisionID 11, got %d", rev.RevisionID)
	}
}

func TestPagesService_RevertToRevision(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/pages/intro/revisions/5" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(PageRevision{RevisionID: 5})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewPagesService(client)
	rev, err := svc.RevertToRevision(context.Background(), 1, "intro", 5)
	if err != nil {
		t.Fatalf("RevertToRevision: %v", err)
	}
	if rev.RevisionID != 5 {
		t.Errorf("expected RevisionID 5, got %d", rev.RevisionID)
	}
}

func TestPagesService_List_WithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		q := r.URL.Query()
		if q.Get("sort") != "title" {
			t.Errorf("expected sort=title, got %q", q.Get("sort"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Page{{Title: "Page 1", URL: "page-1"}})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewPagesService(client)
	published := true
	opts := &ListPagesOptions{
		Sort:       "title",
		Order:      "asc",
		SearchTerm: "page",
		Published:  &published,
		Page:       1,
		PerPage:    10,
	}
	pages, err := svc.List(context.Background(), 1, opts)
	if err != nil {
		t.Fatalf("List with opts: %v", err)
	}
	if len(pages) != 1 {
		t.Errorf("expected 1, got %d", len(pages))
	}
}

func TestNewPagesService(t *testing.T) {
	client := &Client{}
	service := NewPagesService(client)
	if service == nil {
		t.Fatal("Expected non-nil service")
		return
	}
	if service.client != client {
		t.Error("Expected service client to match input client")
	}
}
