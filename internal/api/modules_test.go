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

func TestModulesService_Relock(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/modules/2/relock" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Module{ID: 2, Name: "M1"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewModulesService(client)
	mod, err := svc.Relock(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("Relock: %v", err)
	}
	if mod.ID != 2 {
		t.Errorf("expected ID 2, got %d", mod.ID)
	}
}

func TestModulesService_UpdateItem(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/modules/2/items/3" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		itemData, ok := body["module_item"].(map[string]interface{})
		if !ok {
			t.Fatalf("expected module_item in body")
		}
		if itemData["title"] != "Updated Item" {
			t.Errorf("expected title 'Updated Item', got %v", itemData["title"])
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ModuleItem{ID: 3, Title: "Updated Item"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewModulesService(client)
	title := "Updated Item"
	published := true
	pos := 1
	indent := 0
	extURL := "https://example.com"
	newTab := false
	moveID := int64(5)
	params := &UpdateModuleItemParams{
		Title:                 &title,
		Published:             &published,
		Position:              &pos,
		Indent:                &indent,
		ExternalURL:           &extURL,
		NewTab:                &newTab,
		MoveToModuleID:        &moveID,
		CompletionRequirement: &CompletionRequirementParams{Type: "must_view"},
	}
	item, err := svc.UpdateItem(context.Background(), 1, 2, 3, params)
	if err != nil {
		t.Fatalf("UpdateItem: %v", err)
	}
	if item.Title != "Updated Item" {
		t.Errorf("expected 'Updated Item', got %s", item.Title)
	}
}

func TestModulesService_MarkItemDone(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/modules/2/items/3/done" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewModulesService(client)
	if err := svc.MarkItemDone(context.Background(), 1, 2, 3); err != nil {
		t.Fatalf("MarkItemDone: %v", err)
	}
}

func TestModulesService_MarkItemNotDone(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/modules/2/items/3/done" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewModulesService(client)
	if err := svc.MarkItemNotDone(context.Background(), 1, 2, 3); err != nil {
		t.Fatalf("MarkItemNotDone: %v", err)
	}
}

func TestModulesService_MarkItemRead(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/modules/2/items/3/mark_read" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewModulesService(client)
	if err := svc.MarkItemRead(context.Background(), 1, 2, 3); err != nil {
		t.Fatalf("MarkItemRead: %v", err)
	}
}

func TestModulesService_GetItemSequence(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/1/module_item_sequence" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("asset_type") != "Assignment" {
			t.Errorf("expected asset_type=Assignment, got %q", r.URL.Query().Get("asset_type"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ModuleItemSequence{
			Items:   []ModuleItemSequenceNode{},
			Modules: []ModuleReference{{ID: 2, Name: "Week 1"}},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewModulesService(client)
	seq, err := svc.GetItemSequence(context.Background(), 1, "Assignment", 10)
	if err != nil {
		t.Fatalf("GetItemSequence: %v", err)
	}
	if len(seq.Modules) != 1 {
		t.Errorf("expected 1 module, got %d", len(seq.Modules))
	}
}

func TestModulesService_ListItems_WithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		q := r.URL.Query()
		if q.Get("search_term") != "quiz" {
			t.Errorf("expected search_term=quiz, got %q", q.Get("search_term"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]ModuleItem{{ID: 5, Title: "Quiz Item"}})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewModulesService(client)
	opts := &ListModuleItemsOptions{
		Include:    []string{"content_details"},
		SearchTerm: "quiz",
		StudentID:  "42",
		Page:       1,
		PerPage:    10,
	}
	items, err := svc.ListItems(context.Background(), 1, 2, opts)
	if err != nil {
		t.Fatalf("ListItems: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}
}

func TestNewModulesService(t *testing.T) {
	client := &Client{}
	service := NewModulesService(client)
	if service == nil {
		t.Fatal("Expected non-nil service")
		return
	}
	if service.client != client {
		t.Error("Expected service client to match input client")
	}
}
