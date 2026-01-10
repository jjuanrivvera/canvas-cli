package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewFilesService(t *testing.T) {
	client := &Client{}
	service := NewFilesService(client)

	if service == nil {
		t.Fatal("NewFilesService returned nil")
	}
	if service.client != client {
		t.Error("NewFilesService did not set client correctly")
	}
}

func TestFilesService_ListCourseFiles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if !strings.HasPrefix(r.URL.Path, "/api/v1/courses/100/files") {
			t.Errorf("Expected path to start with /api/v1/courses/100/files, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"id": 1,
				"display_name": "file1.pdf",
				"filename": "file1.pdf",
				"size": 1024
			},
			{
				"id": 2,
				"display_name": "file2.doc",
				"filename": "file2.doc",
				"size": 2048
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

	service := NewFilesService(client)
	ctx := context.Background()

	files, err := service.ListCourseFiles(ctx, 100, nil)
	if err != nil {
		t.Fatalf("ListCourseFiles failed: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(files))
	}
	if files[0].ID != 1 {
		t.Errorf("Expected first file ID 1, got %d", files[0].ID)
	}
}

func TestFilesService_ListCourseFiles_WithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		// Check query parameters
		searchTerm := r.URL.Query().Get("search_term")
		if searchTerm != "test" {
			t.Errorf("Expected search_term 'test', got %s", searchTerm)
		}

		sort := r.URL.Query().Get("sort")
		if sort != "name" {
			t.Errorf("Expected sort 'name', got %s", sort)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id": 1, "display_name": "test.pdf"}]`))
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

	service := NewFilesService(client)
	ctx := context.Background()

	opts := &ListFilesOptions{
		SearchTerm: "test",
		Sort:       "name",
		Order:      "asc",
		Page:       1,
		PerPage:    10,
	}

	files, err := service.ListCourseFiles(ctx, 100, opts)
	if err != nil {
		t.Fatalf("ListCourseFiles failed: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(files))
	}
}

func TestFilesService_ListFolderFiles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if !strings.HasPrefix(r.URL.Path, "/api/v1/folders/50/files") {
			t.Errorf("Expected path to start with /api/v1/folders/50/files, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id": 1, "display_name": "file.pdf"}]`))
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

	service := NewFilesService(client)
	ctx := context.Background()

	files, err := service.ListFolderFiles(ctx, 50, nil)
	if err != nil {
		t.Fatalf("ListFolderFiles failed: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(files))
	}
}

func TestFilesService_ListUserFiles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if !strings.HasPrefix(r.URL.Path, "/api/v1/users/75/files") {
			t.Errorf("Expected path to start with /api/v1/users/75/files, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id": 1, "display_name": "user_file.doc"}]`))
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

	service := NewFilesService(client)
	ctx := context.Background()

	files, err := service.ListUserFiles(ctx, 75, nil)
	if err != nil {
		t.Fatalf("ListUserFiles failed: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(files))
	}
}

func TestFilesService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/files/123" {
			t.Errorf("Expected path /api/v1/files/123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 123,
			"display_name": "test.pdf",
			"filename": "test.pdf",
			"size": 1024,
			"url": "https://example.com/download/test.pdf"
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

	service := NewFilesService(client)
	ctx := context.Background()

	file, err := service.Get(ctx, 123, nil)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if file.ID != 123 {
		t.Errorf("Expected file ID 123, got %d", file.ID)
	}
	if file.DisplayName != "test.pdf" {
		t.Errorf("Expected display name 'test.pdf', got %s", file.DisplayName)
	}
}

func TestFilesService_Get_WithInclude(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		// Check include parameter
		includes := r.URL.Query()["include[]"]
		if len(includes) != 1 || includes[0] != "user" {
			t.Errorf("Expected include[] 'user', got %v", includes)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": 123, "display_name": "test.pdf"}`))
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

	service := NewFilesService(client)
	ctx := context.Background()

	_, err = service.Get(ctx, 123, []string{"user"})
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
}

func TestFilesService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/files/123" {
			t.Errorf("Expected path /api/v1/files/123, got %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
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

	service := NewFilesService(client)
	ctx := context.Background()

	err = service.Delete(ctx, 123)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestFilesService_GetCourseQuota(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/courses/100/files/quota" {
			t.Errorf("Expected path /api/v1/courses/100/files/quota, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"quota": 1073741824, "quota_used": 536870912}`))
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

	service := NewFilesService(client)
	ctx := context.Background()

	quota, err := service.GetCourseQuota(ctx, 100)
	if err != nil {
		t.Fatalf("GetCourseQuota failed: %v", err)
	}

	if quota.Quota != 1073741824 {
		t.Errorf("Expected quota 1073741824, got %d", quota.Quota)
	}
	if quota.QuotaUsed != 536870912 {
		t.Errorf("Expected quota_used 536870912, got %d", quota.QuotaUsed)
	}
}

func TestFilesService_GetUserQuota(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/users/50/files/quota" {
			t.Errorf("Expected path /api/v1/users/50/files/quota, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"quota": 2147483648, "quota_used": 1073741824}`))
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

	service := NewFilesService(client)
	ctx := context.Background()

	quota, err := service.GetUserQuota(ctx, 50)
	if err != nil {
		t.Fatalf("GetUserQuota failed: %v", err)
	}

	if quota.Quota != 2147483648 {
		t.Errorf("Expected quota 2147483648, got %d", quota.Quota)
	}
	if quota.QuotaUsed != 1073741824 {
		t.Errorf("Expected quota_used 1073741824, got %d", quota.QuotaUsed)
	}
}

func TestFilesService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/files/123" {
			t.Errorf("Expected path /api/v1/files/123, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}

		// Parse request body
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if body["name"] != "new_name.pdf" {
			t.Errorf("Expected name 'new_name.pdf', got %v", body["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": 123, "display_name": "new_name.pdf"}`))
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

	service := NewFilesService(client)
	ctx := context.Background()

	params := &UpdateParams{
		Name: "new_name.pdf",
	}

	file, err := service.Update(ctx, 123, params)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if file.DisplayName != "new_name.pdf" {
		t.Errorf("Expected display name 'new_name.pdf', got %s", file.DisplayName)
	}
}

func TestFilesService_UploadToCourse(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	var uploadURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path == "/api/v1/courses/100/files" {
			// Step 1: Return upload URL
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"upload_url": "` + uploadURL + `",
				"upload_params": {}
			}`))
			return
		}

		if r.URL.Path == "/upload" {
			// Step 2: Handle actual upload
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": 456, "display_name": "test.txt", "size": 12}`))
			return
		}

		t.Errorf("Unexpected path: %s", r.URL.Path)
	}))
	defer server.Close()
	uploadURL = server.URL + "/upload"

	client, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 10,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	service := NewFilesService(client)
	ctx := context.Background()

	params := &UploadParams{
		Name: "test.txt",
	}

	file, err := service.UploadToCourse(ctx, 100, testFile, params)
	if err != nil {
		t.Fatalf("UploadToCourse failed: %v", err)
	}

	if file.ID != 456 {
		t.Errorf("Expected file ID 456, got %d", file.ID)
	}
}

func TestFilesService_Download(t *testing.T) {
	var downloadURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path == "/api/v1/files/123" {
			// Return file info with download URL
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"id": 123,
				"display_name": "test.pdf",
				"url": "` + downloadURL + `"
			}`))
			return
		}

		if r.URL.Path == "/download/test.pdf" {
			// Return file content
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("file content here"))
			return
		}

		t.Errorf("Unexpected path: %s", r.URL.Path)
	}))
	defer server.Close()
	downloadURL = server.URL + "/download/test.pdf"

	client, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 10,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	service := NewFilesService(client)
	ctx := context.Background()

	// Create temp dir for download
	tempDir := t.TempDir()
	destPath := filepath.Join(tempDir, "downloaded.pdf")

	err = service.Download(ctx, 123, destPath)
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}

	// Verify file was downloaded
	content, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}

	if string(content) != "file content here" {
		t.Errorf("Expected 'file content here', got %s", string(content))
	}
}

func TestFilesService_UploadToFolder(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	var uploadURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path == "/api/v1/folders/789/files" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"upload_url": "` + uploadURL + `",
				"upload_params": {}
			}`))
			return
		}

		if r.URL.Path == "/upload" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": 999, "display_name": "test.txt", "size": 12}`))
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()
	uploadURL = server.URL + "/upload"

	client, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 10,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	service := NewFilesService(client)
	ctx := context.Background()

	params := &UploadParams{
		Name: "test.txt",
	}

	file, err := service.UploadToFolder(ctx, 789, testFile, params)
	if err != nil {
		t.Fatalf("UploadToFolder failed: %v", err)
	}

	if file.ID != 999 {
		t.Errorf("Expected file ID 999, got %d", file.ID)
	}
}

func TestFilesService_UploadToUser(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	var uploadURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path == "/api/v1/users/456/files" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"upload_url": "` + uploadURL + `",
				"upload_params": {}
			}`))
			return
		}

		if r.URL.Path == "/upload" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": 888, "display_name": "test.txt", "size": 12}`))
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()
	uploadURL = server.URL + "/upload"

	client, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "test-token",
		RequestsPerSec: 10,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	service := NewFilesService(client)
	ctx := context.Background()

	params := &UploadParams{
		Name: "test.txt",
	}

	file, err := service.UploadToUser(ctx, 456, testFile, params)
	if err != nil {
		t.Fatalf("UploadToUser failed: %v", err)
	}

	if file.ID != 888 {
		t.Errorf("Expected file ID 888, got %d", file.ID)
	}
}
