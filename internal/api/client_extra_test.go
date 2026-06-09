package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/oauth2"
)

// mockTokenSource implements oauth2.TokenSource for testing
type mockTokenSource struct {
	token *oauth2.Token
	err   error
}

func (m *mockTokenSource) Token() (*oauth2.Token, error) {
	return m.token, m.err
}

func TestClient_GetToken_FromTokenSource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		// Verify the Authorization header uses the token from the source
		auth := r.Header.Get("Authorization")
		if auth != "Bearer oauth-token-123" {
			t.Errorf("expected bearer oauth-token-123, got %q", auth)
		}
		json.NewEncoder(w).Encode([]Course{{ID: 1}})
	}))
	defer server.Close()

	tokenSource := &mockTokenSource{
		token: &oauth2.Token{AccessToken: "oauth-token-123"},
	}

	client, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		TokenSource:    tokenSource,
		RequestsPerSec: 10,
	})
	if err != nil {
		t.Fatalf("NewClient with TokenSource: %v", err)
	}

	var courses []Course
	if err := client.GetAllPages(context.Background(), "/api/v1/courses", &courses); err != nil {
		t.Fatalf("GetAllPages with TokenSource: %v", err)
	}
}

func TestClient_GetToken_TokenSourceError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		json.NewEncoder(w).Encode([]Course{})
	}))
	defer server.Close()

	tokenSource := &mockTokenSource{
		err: fmt.Errorf("token refresh failed"),
	}

	client, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		TokenSource:    tokenSource,
		RequestsPerSec: 10,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	var courses []Course
	err = client.GetAllPages(context.Background(), "/api/v1/courses", &courses)
	if err == nil {
		t.Error("expected error when token source fails, got nil")
	}
}

func TestClient_AsUserID_MasqueradeParam(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Query().Get("as_user_id") != "777" {
			t.Errorf("expected as_user_id=777, got %q", r.URL.Query().Get("as_user_id"))
		}
		json.NewEncoder(w).Encode([]Course{{ID: 1}})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "t",
		AsUserID:       777,
		RequestsPerSec: 10,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	var courses []Course
	if err := client.GetAllPages(context.Background(), "/api/v1/courses", &courses); err != nil {
		t.Fatalf("GetAllPages with masquerade: %v", err)
	}
}

func TestClient_GetAllPagesGeneric_Pagination(t *testing.T) {
	page1Called := false
	page2Called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		q := r.URL.Query()
		if q.Get("page") == "2" {
			page2Called = true
			json.NewEncoder(w).Encode([]Course{{ID: 3}, {ID: 4}})
		} else {
			page1Called = true
			// Return Link header with next page
			w.Header().Set("Link", fmt.Sprintf(`<http://%s/api/v1/courses?page=2>; rel="next"`, r.Host))
			json.NewEncoder(w).Encode([]Course{{ID: 1}, {ID: 2}})
		}
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 100})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	courses, err := GetAllPagesGeneric[Course](client, context.Background(), "/api/v1/courses")
	if err != nil {
		t.Fatalf("GetAllPagesGeneric: %v", err)
	}
	if !page1Called || !page2Called {
		t.Error("expected both pages to be fetched")
	}
	if len(courses) != 4 {
		t.Errorf("expected 4 courses, got %d", len(courses))
	}
}

func TestClient_GetAllPagesGeneric_WithMaxResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.Header().Set("Link", fmt.Sprintf(`<http://%s/api/v1/courses?page=2>; rel="next"`, r.Host))
		json.NewEncoder(w).Encode([]Course{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "t",
		RequestsPerSec: 100,
		MaxResults:     3,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	courses, err := GetAllPagesGeneric[Course](client, context.Background(), "/api/v1/courses")
	if err != nil {
		t.Fatalf("GetAllPagesGeneric with max: %v", err)
	}
	if len(courses) != 3 {
		t.Errorf("expected 3 courses (MaxResults=3), got %d", len(courses))
	}
}

func TestClient_DryRun_WithBody(t *testing.T) {
	// In dry-run mode, POST with a JSON body should still print the curl command
	// and return a mock response without panicking
	client, err := NewClient(ClientConfig{
		BaseURL:        "https://canvas.example.com",
		Token:          "secret",
		RequestsPerSec: 10,
		DryRun:         true,
	})
	if err != nil {
		t.Fatalf("NewClient dry-run: %v", err)
	}

	// PostJSON in dry-run mode returns "[]" which cannot unmarshal into *Course
	// but we call it to cover the body-reading code path
	var result Course
	// This will fail to decode ([] is not a Course object) but that's expected
	// The important thing is the curl printing code path runs without panic
	_ = client.PostJSON(context.Background(), "/api/v1/courses/1", map[string]string{"name": "test"}, &result)
}

func TestClient_NewClient_EmptyBaseURL(t *testing.T) {
	_, err := NewClient(ClientConfig{Token: "t"})
	if err == nil {
		t.Error("expected error for empty base URL")
	}
}

func TestClient_NewClient_NoTokenOrSource(t *testing.T) {
	_, err := NewClient(ClientConfig{BaseURL: "https://example.com"})
	if err == nil {
		t.Error("expected error for no token or token source")
	}
}

func TestClient_GetMaxResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.Write([]byte(`[]`))
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "t",
		RequestsPerSec: 10,
		MaxResults:     50,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if client.GetMaxResults() != 50 {
		t.Errorf("expected MaxResults=50, got %d", client.GetMaxResults())
	}
}
