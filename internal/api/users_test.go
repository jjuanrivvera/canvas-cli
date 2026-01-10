package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// handleVersionDetection handles the version detection request made during client initialization
func handleVersionDetection(w http.ResponseWriter) {
	w.Header().Set("X-Canvas-Meta", `{"primaryCollection":"accounts"}`)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`[]`))
}

func TestUsersService_GetCurrentUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle version detection
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/users/self" {
			t.Errorf("Expected path /api/v1/users/self, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 1,
			"name": "Test User",
			"login_id": "testuser",
			"email": "test@example.com"
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

	service := NewUsersService(client)
	ctx := context.Background()

	user, err := service.GetCurrentUser(ctx)
	if err != nil {
		t.Fatalf("GetCurrentUser failed: %v", err)
	}

	if user.ID != 1 {
		t.Errorf("Expected user ID 1, got %d", user.ID)
	}
	if user.Name != "Test User" {
		t.Errorf("Expected user name 'Test User', got %s", user.Name)
	}
	if user.LoginID != "testuser" {
		t.Errorf("Expected login ID 'testuser', got %s", user.LoginID)
	}
}

func TestUsersService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/users/123" {
			t.Errorf("Expected path /api/v1/users/123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 123,
			"name": "Specific User",
			"email": "specific@example.com"
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

	service := NewUsersService(client)
	ctx := context.Background()

	user, err := service.Get(ctx, 123, nil)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if user.ID != 123 {
		t.Errorf("Expected user ID 123, got %d", user.ID)
	}
	if user.Name != "Specific User" {
		t.Errorf("Expected user name 'Specific User', got %s", user.Name)
	}
}

func TestUsersService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/users" {
			t.Errorf("Expected path /api/v1/accounts/1/users, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"id": 1,
				"name": "User 1",
				"email": "user1@example.com"
			},
			{
				"id": 2,
				"name": "User 2",
				"email": "user2@example.com"
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

	service := NewUsersService(client)
	ctx := context.Background()

	users, err := service.List(ctx, 1, nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
	if users[0].ID != 1 {
		t.Errorf("Expected first user ID 1, got %d", users[0].ID)
	}
	if users[1].ID != 2 {
		t.Errorf("Expected second user ID 2, got %d", users[1].ID)
	}
}

func TestUsersService_Search(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/search/recipients" {
			t.Errorf("Expected path /api/v1/search/recipients, got %s", r.URL.Path)
		}

		searchTerm := r.URL.Query().Get("search")
		if searchTerm != "john" {
			t.Errorf("Expected search term 'john', got %s", searchTerm)
		}

		typeParam := r.URL.Query().Get("type")
		if typeParam != "user" {
			t.Errorf("Expected type 'user', got %s", typeParam)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"id": 1,
				"name": "John Doe",
				"email": "john@example.com"
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

	service := NewUsersService(client)
	ctx := context.Background()

	users, err := service.Search(ctx, "john")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users))
	}
	if users[0].Name != "John Doe" {
		t.Errorf("Expected user name 'John Doe', got %s", users[0].Name)
	}
}

func TestUsersService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.URL.Path != "/api/v1/accounts/1/users/123" {
			t.Errorf("Expected path /api/v1/accounts/1/users/123, got %s", r.URL.Path)
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

	service := NewUsersService(client)
	ctx := context.Background()

	err = service.Delete(ctx, 1, 123)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestNewUsersService(t *testing.T) {
	client := &Client{}
	service := NewUsersService(client)

	if service == nil {
		t.Fatal("NewUsersService returned nil")
	}
	if service.client != client {
		t.Error("NewUsersService did not set client correctly")
	}
}

func TestUsersService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/accounts/1/users" {
			t.Errorf("Expected path /api/v1/accounts/1/users, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 123,
			"name": "John Doe",
			"short_name": "John",
			"sortable_name": "Doe, John"
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

	service := NewUsersService(client)
	ctx := context.Background()

	params := &CreateUserParams{
		Name:         "John Doe",
		ShortName:    "John",
		SortableName: "Doe, John",
	}

	user, err := service.Create(ctx, 1, params)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if user.ID != 123 {
		t.Errorf("Expected user ID 123, got %d", user.ID)
	}

	if user.Name != "John Doe" {
		t.Errorf("Expected user name 'John Doe', got %s", user.Name)
	}
}

func TestUsersService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/users/123" {
			t.Errorf("Expected path /api/v1/users/123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 123,
			"name": "Jane Smith",
			"short_name": "Jane"
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

	service := NewUsersService(client)
	ctx := context.Background()

	params := &UpdateUserParams{
		Name:      "Jane Smith",
		ShortName: "Jane",
	}

	user, err := service.Update(ctx, 123, params)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if user.ID != 123 {
		t.Errorf("Expected user ID 123, got %d", user.ID)
	}

	if user.Name != "Jane Smith" {
		t.Errorf("Expected user name 'Jane Smith', got %s", user.Name)
	}
}

func TestUsersService_ListCourseUsers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/courses/100/users" {
			t.Errorf("Expected path /api/v1/courses/100/users, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"id": 1,
				"name": "Student One"
			},
			{
				"id": 2,
				"name": "Student Two"
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

	service := NewUsersService(client)
	ctx := context.Background()

	users, err := service.ListCourseUsers(ctx, 100, nil)
	if err != nil {
		t.Fatalf("ListCourseUsers failed: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}

	if users[0].ID != 1 {
		t.Errorf("Expected first user ID 1, got %d", users[0].ID)
	}
}

func TestUsersService_ListCourseUsers_WithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		searchTerm := r.URL.Query().Get("search_term")
		if searchTerm != "alice" {
			t.Errorf("Expected search term 'alice', got %s", searchTerm)
		}

		enrollmentType := r.URL.Query().Get("enrollment_type[]")
		if enrollmentType != "student" {
			t.Errorf("Expected enrollment type 'student', got %s", enrollmentType)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"id": 1,
				"name": "Alice"
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

	service := NewUsersService(client)
	ctx := context.Background()

	opts := &ListUsersOptions{
		SearchTerm:     "alice",
		EnrollmentType: "student",
	}

	users, err := service.ListCourseUsers(ctx, 100, opts)
	if err != nil {
		t.Fatalf("ListCourseUsers failed: %v", err)
	}

	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users))
	}

	if users[0].Name != "Alice" {
		t.Errorf("Expected user name 'Alice', got %s", users[0].Name)
	}
}

func TestUsersService_List_WithAllOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		// Check query parameters
		query := r.URL.Query()
		if query.Get("search_term") == "" {
			t.Error("Expected search_term parameter")
		}
		if query.Get("enrollment_type") == "" {
			t.Error("Expected enrollment_type parameter")
		}
		if query.Get("enrollment_state") == "" {
			t.Error("Expected enrollment_state parameter")
		}
		if len(query["include[]"]) == 0 {
			t.Error("Expected include[] parameters")
		}
		if query.Get("page") == "" {
			t.Error("Expected page parameter")
		}
		if query.Get("per_page") == "" {
			t.Error("Expected per_page parameter")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id": 1, "name": "Test User"}]`))
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

	service := NewUsersService(client)
	ctx := context.Background()

	opts := &ListUsersOptions{
		SearchTerm:      "test",
		EnrollmentType:  "teacher",
		EnrollmentState: "active",
		Include:         []string{"email", "enrollments", "avatar_url"},
		Page:            2,
		PerPage:         50,
	}

	users, err := service.List(ctx, 1, opts)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users))
	}
}

func TestUsersService_Create_WithAllOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		// Verify request body contains all parameters
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		user, ok := body["user"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected user object in body")
		}

		// Check user fields
		expectedUserFields := []string{"name", "short_name", "sortable_name", "time_zone", "locale", "terms_of_use", "skip_registration"}
		for _, field := range expectedUserFields {
			if _, ok := user[field]; !ok {
				t.Errorf("Expected user field %s", field)
			}
		}

		// Check pseudonym
		if _, ok := body["pseudonym"]; !ok {
			t.Error("Expected pseudonym in body")
		}

		// Check communication_channel
		if _, ok := body["communication_channel"]; !ok {
			t.Error("Expected communication_channel in body")
		}

		// Check additional flags
		if _, ok := body["force_validations"]; !ok {
			t.Error("Expected force_validations in body")
		}
		if _, ok := body["enable_sis_reactivation"]; !ok {
			t.Error("Expected enable_sis_reactivation in body")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": 999, "name": "New User"}`))
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

	service := NewUsersService(client)
	ctx := context.Background()

	params := &CreateUserParams{
		Name:                     "New User",
		ShortName:                "NewU",
		SortableName:             "User, New",
		TimeZone:                 "America/New_York",
		Locale:                   "en",
		TermsOfUse:               true,
		SkipRegistration:         true,
		ForceValidations:         true,
		EnableSISReactivation:    true,
		UniqueID:                 "newuser123",
		SISUserID:                "SIS999",
		IntegrationID:            "INT999",
		AuthenticationProviderID: "auth123",
		Password:                 "P@ssw0rd!",
		Email:                    "newuser@example.com",
		SkipConfirmation:         true,
	}

	user, err := service.Create(ctx, 1, params)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if user.ID != 999 {
		t.Errorf("Expected user ID 999, got %d", user.ID)
	}
}

func TestUsersService_Get_WithIncludes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		// Check include parameters
		query := r.URL.Query()
		includes := query["include[]"]
		if len(includes) == 0 {
			t.Error("Expected include[] parameters")
		}

		expectedIncludes := map[string]bool{
			"email":       false,
			"enrollments": false,
			"avatar_url":  false,
		}

		for _, inc := range includes {
			expectedIncludes[inc] = true
		}

		for inc, found := range expectedIncludes {
			if !found {
				t.Errorf("Expected include parameter %s", inc)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": 123, "name": "Test User", "email": "test@example.com"}`))
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

	service := NewUsersService(client)
	ctx := context.Background()

	include := []string{"email", "enrollments", "avatar_url"}
	user, err := service.Get(ctx, 123, include)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if user.ID != 123 {
		t.Errorf("Expected user ID 123, got %d", user.ID)
	}
}

func TestUsersService_Update_WithAllOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		// Verify request body contains all parameters
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		user, ok := body["user"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected user object in body")
		}

		// Check all user fields
		expectedFields := []string{"name", "short_name", "sortable_name", "time_zone", "locale", "email", "avatar"}
		for _, field := range expectedFields {
			if _, ok := user[field]; !ok {
				t.Errorf("Expected user field %s", field)
			}
		}

		// Check avatar subfields
		avatar, ok := user["avatar"].(map[string]interface{})
		if !ok {
			t.Error("Expected avatar object")
		} else {
			if _, ok := avatar["token"]; !ok {
				t.Error("Expected avatar token")
			}
			if _, ok := avatar["url"]; !ok {
				t.Error("Expected avatar url")
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": 123, "name": "Updated User"}`))
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

	service := NewUsersService(client)
	ctx := context.Background()

	params := &UpdateUserParams{
		Name:         "Updated User",
		ShortName:    "UpUser",
		SortableName: "User, Updated",
		TimeZone:     "America/Los_Angeles",
		Locale:       "es",
		Email:        "updated@example.com",
		Avatar: &AvatarParams{
			Token: "avatar_token_123",
			URL:   "https://example.com/avatar.jpg",
		},
	}

	user, err := service.Update(ctx, 123, params)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if user.ID != 123 {
		t.Errorf("Expected user ID 123, got %d", user.ID)
	}
}
