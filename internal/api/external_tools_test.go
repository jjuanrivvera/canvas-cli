package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExternalToolsService_ListByCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/external_tools" {
			t.Errorf("Expected path /api/v1/courses/123/external_tools, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected method GET, got %s", r.Method)
		}

		tools := []ExternalTool{
			{ID: 1, Name: "Tool One", URL: "https://tool1.example.com"},
			{ID: 2, Name: "Tool Two", URL: "https://tool2.example.com"},
		}
		json.NewEncoder(w).Encode(tools)
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

	service := NewExternalToolsService(client)
	tools, err := service.ListByCourse(context.Background(), 123, nil)
	if err != nil {
		t.Fatalf("ListByCourse failed: %v", err)
	}

	if len(tools) != 2 {
		t.Errorf("Expected 2 tools, got %d", len(tools))
	}

	if tools[0].Name != "Tool One" {
		t.Errorf("Expected tool name 'Tool One', got '%s'", tools[0].Name)
	}
}

func TestExternalToolsService_ListByAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/external_tools" {
			t.Errorf("Expected path /api/v1/accounts/1/external_tools, got %s", r.URL.Path)
		}

		tools := []ExternalTool{
			{ID: 1, Name: "Account Tool", URL: "https://tool.example.com"},
		}
		json.NewEncoder(w).Encode(tools)
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

	service := NewExternalToolsService(client)
	tools, err := service.ListByAccount(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("ListByAccount failed: %v", err)
	}

	if len(tools) != 1 {
		t.Errorf("Expected 1 tool, got %d", len(tools))
	}
}

func TestExternalToolsService_GetByCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/external_tools/456" {
			t.Errorf("Expected path /api/v1/courses/123/external_tools/456, got %s", r.URL.Path)
		}

		tool := ExternalTool{
			ID:          456,
			Name:        "Test Tool",
			Description: "A test tool",
			URL:         "https://tool.example.com",
		}
		json.NewEncoder(w).Encode(tool)
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

	service := NewExternalToolsService(client)
	tool, err := service.GetByCourse(context.Background(), 123, 456)
	if err != nil {
		t.Fatalf("GetByCourse failed: %v", err)
	}

	if tool.ID != 456 {
		t.Errorf("Expected tool ID 456, got %d", tool.ID)
	}

	if tool.Name != "Test Tool" {
		t.Errorf("Expected tool name 'Test Tool', got '%s'", tool.Name)
	}
}

func TestExternalToolsService_CreateInCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/external_tools" {
			t.Errorf("Expected path /api/v1/courses/123/external_tools, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		tool := ExternalTool{
			ID:           789,
			Name:         "New Tool",
			URL:          "https://newtool.example.com",
			PrivacyLevel: "public",
		}
		json.NewEncoder(w).Encode(tool)
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

	service := NewExternalToolsService(client)
	params := &CreateExternalToolParams{
		Name:         "New Tool",
		URL:          "https://newtool.example.com",
		ConsumerKey:  "key123",
		SharedSecret: "secret123",
		PrivacyLevel: "public",
	}

	tool, err := service.CreateInCourse(context.Background(), 123, params)
	if err != nil {
		t.Fatalf("CreateInCourse failed: %v", err)
	}

	if tool.ID != 789 {
		t.Errorf("Expected tool ID 789, got %d", tool.ID)
	}

	if tool.Name != "New Tool" {
		t.Errorf("Expected tool name 'New Tool', got '%s'", tool.Name)
	}
}

func TestExternalToolsService_UpdateInCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/external_tools/456" {
			t.Errorf("Expected path /api/v1/courses/123/external_tools/456, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("Expected method PUT, got %s", r.Method)
		}

		tool := ExternalTool{
			ID:          456,
			Name:        "Updated Tool",
			Description: "Updated description",
		}
		json.NewEncoder(w).Encode(tool)
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

	service := NewExternalToolsService(client)
	name := "Updated Tool"
	desc := "Updated description"
	params := &UpdateExternalToolParams{
		Name:        &name,
		Description: &desc,
	}

	tool, err := service.UpdateInCourse(context.Background(), 123, 456, params)
	if err != nil {
		t.Fatalf("UpdateInCourse failed: %v", err)
	}

	if tool.Name != "Updated Tool" {
		t.Errorf("Expected tool name 'Updated Tool', got '%s'", tool.Name)
	}
}

func TestExternalToolsService_DeleteFromCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/external_tools/456" {
			t.Errorf("Expected path /api/v1/courses/123/external_tools/456, got %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("Expected method DELETE, got %s", r.Method)
		}

		tool := ExternalTool{ID: 456, Name: "Deleted Tool"}
		json.NewEncoder(w).Encode(tool)
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

	service := NewExternalToolsService(client)
	tool, err := service.DeleteFromCourse(context.Background(), 123, 456)
	if err != nil {
		t.Fatalf("DeleteFromCourse failed: %v", err)
	}

	if tool == nil {
		t.Fatal("Expected tool, got nil")
	}

	if tool.ID != 456 {
		t.Errorf("Expected tool ID 456, got %d", tool.ID)
	}
}

func TestExternalToolsService_GetSessionlessLaunchURLForCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/123/external_tools/sessionless_launch" {
			t.Errorf("Expected path /api/v1/courses/123/external_tools/sessionless_launch, got %s", r.URL.Path)
		}

		if r.URL.Query().Get("id") != "456" {
			t.Errorf("Expected id=456, got %s", r.URL.Query().Get("id"))
		}

		if r.URL.Query().Get("launch_type") != "course_navigation" {
			t.Errorf("Expected launch_type=course_navigation, got %s", r.URL.Query().Get("launch_type"))
		}

		result := SessionlessLaunchURL{
			ID:   456,
			Name: "Tool",
			URL:  "https://canvas.example.com/launch/xyz123",
		}
		json.NewEncoder(w).Encode(result)
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

	service := NewExternalToolsService(client)
	params := &SessionlessLaunchParams{
		ID:         456,
		LaunchType: "course_navigation",
	}

	result, err := service.GetSessionlessLaunchURLForCourse(context.Background(), 123, params)
	if err != nil {
		t.Fatalf("GetSessionlessLaunchURLForCourse failed: %v", err)
	}

	if result.URL == "" {
		t.Error("Expected launch URL, got empty string")
	}

	if result.ID != 456 {
		t.Errorf("Expected ID 456, got %d", result.ID)
	}
}

func TestExternalToolsService_ListWithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Query().Get("search_term") != "test" {
			t.Errorf("Expected search_term=test, got %s", r.URL.Query().Get("search_term"))
		}

		if r.URL.Query().Get("selectable") != "true" {
			t.Errorf("Expected selectable=true, got %s", r.URL.Query().Get("selectable"))
		}

		if r.URL.Query().Get("include_parents") != "true" {
			t.Errorf("Expected include_parents=true, got %s", r.URL.Query().Get("include_parents"))
		}

		tools := []ExternalTool{
			{ID: 1, Name: "Test Tool"},
		}
		json.NewEncoder(w).Encode(tools)
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

	service := NewExternalToolsService(client)
	selectable := true
	opts := &ListExternalToolsOptions{
		Search:         "test",
		Selectable:     &selectable,
		IncludeParents: true,
	}

	tools, err := service.ListByCourse(context.Background(), 123, opts)
	if err != nil {
		t.Fatalf("ListByCourse with options failed: %v", err)
	}

	if len(tools) != 1 {
		t.Errorf("Expected 1 tool, got %d", len(tools))
	}
}
