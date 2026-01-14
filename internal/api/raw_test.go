package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// handleRawVersionDetection handles version detection request during client init
func handleRawVersionDetection(w http.ResponseWriter) {
	w.Header().Set("X-Canvas-Meta", `{"primaryCollection":"accounts"}`)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`[]`))
}

func TestRawService_Request_GET(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle version detection
		if r.URL.Path == "/api/v1/accounts" {
			handleRawVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses" {
			t.Errorf("expected /api/v1/courses, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"id": 1, "name": "Course 1"},
			{"id": 2, "name": "Course 2"},
		})
	}))
	defer server.Close()

	client := createTestClient(server.URL)
	service := NewRawService(client)

	resp, err := service.Request(context.Background(), "GET", "/api/v1/courses", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var courses []map[string]interface{}
	if err := json.Unmarshal(resp.Body, &courses); err != nil {
		t.Fatalf("failed to unmarshal body: %v", err)
	}

	if len(courses) != 2 {
		t.Errorf("expected 2 courses, got %d", len(courses))
	}
}

func TestRawService_Request_POST(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleRawVersionDetection(w)
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}

		if body["name"] != "Test Course" {
			t.Errorf("expected name 'Test Course', got %v", body["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":   123,
			"name": "Test Course",
		})
	}))
	defer server.Close()

	client := createTestClient(server.URL)
	service := NewRawService(client)

	resp, err := service.Request(context.Background(), "POST", "/api/v1/courses", &RawRequestOptions{
		Body: map[string]interface{}{"name": "Test Course"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}
}

func TestRawService_Request_DELETE(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleRawVersionDetection(w)
			return
		}

		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := createTestClient(server.URL)
	service := NewRawService(client)

	resp, err := service.Request(context.Background(), "DELETE", "/api/v1/courses/123", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestRawService_Request_WithQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleRawVersionDetection(w)
			return
		}

		if r.URL.Query().Get("search_term") != "test" {
			t.Errorf("expected search_term=test, got %s", r.URL.Query().Get("search_term"))
		}
		if r.URL.Query().Get("per_page") != "50" {
			t.Errorf("expected per_page=50, got %s", r.URL.Query().Get("per_page"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
	}))
	defer server.Close()

	client := createTestClient(server.URL)
	service := NewRawService(client)

	resp, err := service.Request(context.Background(), "GET", "/api/v1/users", &RawRequestOptions{
		Query: map[string][]string{
			"search_term": {"test"},
			"per_page":    {"50"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestRawService_Request_InvalidMethod(t *testing.T) {
	client := createTestClient("http://localhost:8080")
	service := NewRawService(client)

	_, err := service.Request(context.Background(), "INVALID", "/api/v1/test", nil)
	if err == nil {
		t.Error("expected error for invalid method")
	}
}

func TestRawService_Request_WithPagination(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleRawVersionDetection(w)
			return
		}

		requestCount++
		w.Header().Set("Content-Type", "application/json")

		if requestCount == 1 {
			// First page with link header
			w.Header().Set("Link", `<http://`+r.Host+`/api/v1/courses?page=2>; rel="next"`)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]map[string]interface{}{
				{"id": 1, "name": "Course 1"},
			})
		} else {
			// Second page (last)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]map[string]interface{}{
				{"id": 2, "name": "Course 2"},
			})
		}
	}))
	defer server.Close()

	client := createTestClient(server.URL)
	service := NewRawService(client)

	resp, err := service.Request(context.Background(), "GET", "/api/v1/courses", &RawRequestOptions{
		Paginate: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var courses []map[string]interface{}
	if err := json.Unmarshal(resp.Body, &courses); err != nil {
		t.Fatalf("failed to unmarshal body: %v", err)
	}

	if len(courses) != 2 {
		t.Errorf("expected 2 courses from pagination, got %d", len(courses))
	}

	if requestCount != 2 {
		t.Errorf("expected 2 requests for pagination, got %d", requestCount)
	}
}

func TestRawService_GetRequestStats(t *testing.T) {
	client := createTestClient("https://canvas.example.com")
	service := NewRawService(client)

	stats := service.GetRequestStats()

	if stats["base_url"] != "https://canvas.example.com" {
		t.Errorf("expected base_url 'https://canvas.example.com', got %v", stats["base_url"])
	}

	if _, ok := stats["rate_limit"]; !ok {
		t.Error("expected rate_limit in stats")
	}
}

// Helper function to create a test client
func createTestClient(baseURL string) *Client {
	client, _ := NewClient(ClientConfig{
		BaseURL: baseURL,
		Token:   "test-token",
	})
	return client
}
