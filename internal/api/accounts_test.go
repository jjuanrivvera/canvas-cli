package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewAccountsService(t *testing.T) {
	client, err := NewClient(ClientConfig{BaseURL: "https://canvas.example.com", Token: "t"})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	svc := NewAccountsService(client)
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
	if svc.client != client {
		t.Error("expected client to be set")
	}
}

func TestAccountsService_List(t *testing.T) {
	// /api/v1/accounts is used for both version detection and list. We serve
	// a valid account array for every request to that path.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]Account{
				{ID: 1, Name: "Root Account"},
				{ID: 2, Name: "Sub Account"},
			})
			return
		}
		t.Errorf("unexpected path: %s", r.URL.Path)
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	svc := NewAccountsService(client)
	accounts, err := svc.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(accounts) < 1 {
		t.Errorf("expected at least 1 account, got %d", len(accounts))
	}
}

func TestAccountsService_List_WithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			q := r.URL.Query()
			// Version detection probe will arrive without include[] param on first call;
			// subsequent call from List will have per_page.
			_ = q
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]Account{{ID: 1, Name: "Root"}})
			return
		}
		t.Errorf("unexpected path: %s", r.URL.Path)
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	svc := NewAccountsService(client)
	opts := &ListAccountsOptions{Include: []string{"lti_guid"}, PerPage: 10}
	accounts, err := svc.List(context.Background(), opts)
	if err != nil {
		t.Fatalf("List with options: %v", err)
	}
	if len(accounts) < 1 {
		t.Errorf("expected accounts, got %d", len(accounts))
	}
}

func TestAccountsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/42" {
			t.Errorf("expected /api/v1/accounts/42, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Account{ID: 42, Name: "My Account"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	svc := NewAccountsService(client)
	acct, err := svc.Get(context.Background(), 42)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if acct.ID != 42 {
		t.Errorf("expected ID 42, got %d", acct.ID)
	}
	if acct.Name != "My Account" {
		t.Errorf("expected name 'My Account', got %q", acct.Name)
	}
}

func TestAccountsService_Get_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":[{"message":"not found"}]}`))
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	svc := NewAccountsService(client)
	_, err = svc.Get(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAccountsService_ListSubAccounts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/sub_accounts" {
			t.Errorf("expected /api/v1/accounts/1/sub_accounts, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Account{
			{ID: 10, Name: "Sub A"},
			{ID: 11, Name: "Sub B"},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	svc := NewAccountsService(client)
	accounts, err := svc.ListSubAccounts(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("ListSubAccounts: %v", err)
	}
	if len(accounts) != 2 {
		t.Errorf("expected 2, got %d", len(accounts))
	}
}

func TestAccountsService_ListSubAccounts_WithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Query().Get("recursive") != "true" {
			t.Errorf("expected recursive=true, got %q", r.URL.Query().Get("recursive"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Account{{ID: 20, Name: "Deep Sub"}})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	svc := NewAccountsService(client)
	opts := &ListSubAccountsOptions{Recursive: true, PerPage: 50}
	accounts, err := svc.ListSubAccounts(context.Background(), 1, opts)
	if err != nil {
		t.Fatalf("ListSubAccounts with opts: %v", err)
	}
	if len(accounts) != 1 {
		t.Errorf("expected 1, got %d", len(accounts))
	}
}

func TestAccountsService_ListCourses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/courses" {
			t.Errorf("expected /api/v1/accounts/1/courses, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Course{
			{ID: 100, Name: "Course A"},
			{ID: 101, Name: "Course B"},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	svc := NewAccountsService(client)
	courses, err := svc.ListCourses(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("ListCourses: %v", err)
	}
	if len(courses) != 2 {
		t.Errorf("expected 2 courses, got %d", len(courses))
	}
}

func TestAccountsService_ListCourses_WithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		q := r.URL.Query()
		if q.Get("search_term") != "math" {
			t.Errorf("expected search_term=math, got %q", q.Get("search_term"))
		}
		if q.Get("published") != "true" {
			t.Errorf("expected published=true, got %q", q.Get("published"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Course{{ID: 200, Name: "Math 101"}})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	svc := NewAccountsService(client)
	opts := &ListAccountCoursesOptions{
		SearchTerm:          "math",
		Published:           true,
		Completed:           true,
		Blueprint:           true,
		BlueprintAssociated: true,
		WithEnrollments:     true,
		EnrollmentTermID:    5,
		State:               []string{"available"},
		EnrollmentType:      []string{"teacher"},
		ByTeachers:          []int64{1, 2},
		BySubaccounts:       []int64{3},
		Sort:                "course_name",
		Order:               "asc",
		Include:             []string{"term"},
		PerPage:             20,
	}
	courses, err := svc.ListCourses(context.Background(), 1, opts)
	if err != nil {
		t.Fatalf("ListCourses with opts: %v", err)
	}
	if len(courses) != 1 {
		t.Errorf("expected 1 course, got %d", len(courses))
	}
}

func TestAccountsService_ListUsers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/1/users" {
			t.Errorf("expected /api/v1/accounts/1/users, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]User{
			{ID: 10, Name: "Alice"},
			{ID: 11, Name: "Bob"},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	svc := NewAccountsService(client)
	users, err := svc.ListUsers(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("ListUsers: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestAccountsService_ListUsers_WithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		q := r.URL.Query()
		if q.Get("search_term") != "alice" {
			t.Errorf("expected search_term=alice, got %q", q.Get("search_term"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]User{{ID: 10, Name: "Alice"}})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	svc := NewAccountsService(client)
	opts := &ListAccountUsersOptions{
		SearchTerm:     "alice",
		EnrollmentType: "student",
		Sort:           "username",
		Order:          "asc",
		Include:        []string{"email"},
		PerPage:        10,
	}
	users, err := svc.ListUsers(context.Background(), 1, opts)
	if err != nil {
		t.Fatalf("ListUsers with opts: %v", err)
	}
	if len(users) != 1 {
		t.Errorf("expected 1 user, got %d", len(users))
	}
}
