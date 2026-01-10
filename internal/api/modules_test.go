package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestModulesService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/modules" {
			t.Errorf("Expected path /api/v1/courses/123/modules, got %s", r.URL.Path)
		}

		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"id": 1,
				"name": "Week 1",
				"position": 1,
				"workflow_state": "active",
				"items_count": 5,
				"published": true
			},
			{
				"id": 2,
				"name": "Week 2",
				"position": 2,
				"workflow_state": "active",
				"items_count": 3,
				"published": true
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

	service := NewModulesService(client)
	ctx := context.Background()

	modules, err := service.List(ctx, 123, nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(modules) != 2 {
		t.Errorf("Expected 2 modules, got %d", len(modules))
	}
	if modules[0].ID != 1 {
		t.Errorf("Expected first module ID 1, got %d", modules[0].ID)
	}
	if modules[0].Name != "Week 1" {
		t.Errorf("Expected first module name 'Week 1', got %s", modules[0].Name)
	}
	if modules[1].ID != 2 {
		t.Errorf("Expected second module ID 2, got %d", modules[1].ID)
	}
}

func TestModulesService_List_WithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		// Check query parameters
		includes := r.URL.Query()["include[]"]
		hasItems := false
		for _, inc := range includes {
			if inc == "items" {
				hasItems = true
			}
		}
		if !hasItems {
			t.Error("Expected include[]=items parameter")
		}

		searchTerm := r.URL.Query().Get("search_term")
		if searchTerm != "Week" {
			t.Errorf("Expected search_term 'Week', got %s", searchTerm)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id": 1, "name": "Week 1", "position": 1}]`))
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

	service := NewModulesService(client)
	ctx := context.Background()

	opts := &ListModulesOptions{
		Include:    []string{"items"},
		SearchTerm: "Week",
	}

	modules, err := service.List(ctx, 123, opts)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(modules) != 1 {
		t.Errorf("Expected 1 module, got %d", len(modules))
	}
}

func TestModulesService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/modules/456" {
			t.Errorf("Expected path /api/v1/courses/123/modules/456, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 456,
			"name": "Introduction Module",
			"position": 1,
			"workflow_state": "active",
			"items_count": 10,
			"require_sequential_progress": true,
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

	service := NewModulesService(client)
	ctx := context.Background()

	module, err := service.Get(ctx, 123, 456, nil, "")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if module.ID != 456 {
		t.Errorf("Expected module ID 456, got %d", module.ID)
	}
	if module.Name != "Introduction Module" {
		t.Errorf("Expected module name 'Introduction Module', got %s", module.Name)
	}
	if !module.RequireSequentialProgress {
		t.Error("Expected require_sequential_progress to be true")
	}
}

func TestModulesService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/modules" {
			t.Errorf("Expected path /api/v1/courses/123/modules, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		moduleData, ok := body["module"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected 'module' key in request body")
		}

		if moduleData["name"] != "Week 1" {
			t.Errorf("Expected module name 'Week 1', got %v", moduleData["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 789,
			"name": "Week 1",
			"position": 1,
			"workflow_state": "active",
			"items_count": 0,
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

	service := NewModulesService(client)
	ctx := context.Background()

	params := &CreateModuleParams{
		Name: "Week 1",
	}

	module, err := service.Create(ctx, 123, params)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if module.ID != 789 {
		t.Errorf("Expected module ID 789, got %d", module.ID)
	}
	if module.Name != "Week 1" {
		t.Errorf("Expected module name 'Week 1', got %s", module.Name)
	}
}

func TestModulesService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/modules/456" {
			t.Errorf("Expected path /api/v1/courses/123/modules/456, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		moduleData, ok := body["module"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected 'module' key in request body")
		}

		if moduleData["name"] != "Updated Name" {
			t.Errorf("Expected module name 'Updated Name', got %v", moduleData["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 456,
			"name": "Updated Name",
			"position": 1,
			"workflow_state": "active",
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

	service := NewModulesService(client)
	ctx := context.Background()

	name := "Updated Name"
	params := &UpdateModuleParams{
		Name: &name,
	}

	module, err := service.Update(ctx, 123, 456, params)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if module.Name != "Updated Name" {
		t.Errorf("Expected module name 'Updated Name', got %s", module.Name)
	}
}

func TestModulesService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/modules/456" {
			t.Errorf("Expected path /api/v1/courses/123/modules/456, got %s", r.URL.Path)
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

	service := NewModulesService(client)
	ctx := context.Background()

	err = service.Delete(ctx, 123, 456)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestModulesService_ListItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/modules/456/items" {
			t.Errorf("Expected path /api/v1/courses/123/modules/456/items, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"id": 1,
				"module_id": 456,
				"title": "Introduction",
				"type": "Page",
				"position": 1,
				"published": true
			},
			{
				"id": 2,
				"module_id": 456,
				"title": "First Assignment",
				"type": "Assignment",
				"position": 2,
				"content_id": 999,
				"published": true
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

	service := NewModulesService(client)
	ctx := context.Background()

	items, err := service.ListItems(ctx, 123, 456, nil)
	if err != nil {
		t.Fatalf("ListItems failed: %v", err)
	}

	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}
	if items[0].Title != "Introduction" {
		t.Errorf("Expected first item title 'Introduction', got %s", items[0].Title)
	}
	if items[0].Type != "Page" {
		t.Errorf("Expected first item type 'Page', got %s", items[0].Type)
	}
	if items[1].Type != "Assignment" {
		t.Errorf("Expected second item type 'Assignment', got %s", items[1].Type)
	}
}

func TestModulesService_GetItem(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/modules/456/items/789" {
			t.Errorf("Expected path /api/v1/courses/123/modules/456/items/789, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 789,
			"module_id": 456,
			"title": "Quiz 1",
			"type": "Quiz",
			"position": 3,
			"content_id": 111,
			"published": true,
			"completion_requirement": {
				"type": "min_score",
				"min_score": 80
			}
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

	service := NewModulesService(client)
	ctx := context.Background()

	item, err := service.GetItem(ctx, 123, 456, 789, nil, "")
	if err != nil {
		t.Fatalf("GetItem failed: %v", err)
	}

	if item.ID != 789 {
		t.Errorf("Expected item ID 789, got %d", item.ID)
	}
	if item.Type != "Quiz" {
		t.Errorf("Expected item type 'Quiz', got %s", item.Type)
	}
	if item.CompletionRequirement == nil {
		t.Fatal("Expected completion requirement to be present")
	}
	if item.CompletionRequirement.Type != "min_score" {
		t.Errorf("Expected completion type 'min_score', got %s", item.CompletionRequirement.Type)
	}
	if item.CompletionRequirement.MinScore != 80 {
		t.Errorf("Expected min_score 80, got %f", item.CompletionRequirement.MinScore)
	}
}

func TestModulesService_CreateItem(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/modules/456/items" {
			t.Errorf("Expected path /api/v1/courses/123/modules/456/items, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		itemData, ok := body["module_item"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected 'module_item' key in request body")
		}

		if itemData["type"] != "Assignment" {
			t.Errorf("Expected item type 'Assignment', got %v", itemData["type"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 999,
			"module_id": 456,
			"title": "New Assignment",
			"type": "Assignment",
			"position": 1,
			"content_id": 888,
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

	service := NewModulesService(client)
	ctx := context.Background()

	params := &CreateModuleItemParams{
		Type:      "Assignment",
		Title:     "New Assignment",
		ContentID: 888,
	}

	item, err := service.CreateItem(ctx, 123, 456, params)
	if err != nil {
		t.Fatalf("CreateItem failed: %v", err)
	}

	if item.ID != 999 {
		t.Errorf("Expected item ID 999, got %d", item.ID)
	}
	if item.Type != "Assignment" {
		t.Errorf("Expected item type 'Assignment', got %s", item.Type)
	}
}

func TestModulesService_DeleteItem(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/modules/456/items/789" {
			t.Errorf("Expected path /api/v1/courses/123/modules/456/items/789, got %s", r.URL.Path)
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

	service := NewModulesService(client)
	ctx := context.Background()

	err = service.DeleteItem(ctx, 123, 456, 789)
	if err != nil {
		t.Fatalf("DeleteItem failed: %v", err)
	}
}

func TestNewModulesService(t *testing.T) {
	client := &Client{}
	service := NewModulesService(client)
	if service == nil {
		t.Fatal("Expected non-nil service")
	}
	if service.client != client {
		t.Error("Expected service client to match input client")
	}
}
