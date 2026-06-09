package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGroupsService_ListCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/123/groups" {
			t.Errorf("expected /api/v1/courses/123/groups, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Group{
			{ID: 1, Name: "Group 1", MembersCount: 5},
			{ID: 2, Name: "Group 2", MembersCount: 3},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	groups, err := service.ListCourse(context.Background(), 123, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(groups))
	}

	if groups[0].Name != "Group 1" {
		t.Errorf("expected 'Group 1', got %s", groups[0].Name)
	}
}

func TestGroupsService_ListAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/groups" {
			t.Errorf("expected /api/v1/accounts/1/groups, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Group{
			{ID: 1, Name: "Account Group"},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	groups, err := service.ListAccount(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(groups))
	}
}

func TestGroupsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/groups/456" {
			t.Errorf("expected /api/v1/groups/456, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Group{
			ID:           456,
			Name:         "Test Group",
			Description:  "Test description",
			MembersCount: 10,
			JoinLevel:    "invitation_only",
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	group, err := service.Get(context.Background(), 456, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if group.ID != 456 {
		t.Errorf("expected ID 456, got %d", group.ID)
	}

	if group.Name != "Test Group" {
		t.Errorf("expected 'Test Group', got %s", group.Name)
	}
}

func TestGroupsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if body["name"] != "New Group" {
			t.Errorf("expected name 'New Group', got %v", body["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Group{
			ID:   789,
			Name: "New Group",
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	params := &CreateGroupParams{
		Name:      "New Group",
		JoinLevel: "invitation_only",
	}

	group, err := service.Create(context.Background(), 123, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if group.ID != 789 {
		t.Errorf("expected ID 789, got %d", group.ID)
	}
}

func TestGroupsService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if body["name"] != "Updated Group" {
			t.Errorf("expected name 'Updated Group', got %v", body["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Group{
			ID:   456,
			Name: "Updated Group",
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	name := "Updated Group"
	params := &UpdateGroupParams{
		Name: &name,
	}

	group, err := service.Update(context.Background(), 456, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if group.Name != "Updated Group" {
		t.Errorf("expected 'Updated Group', got %s", group.Name)
	}
}

func TestGroupsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/groups/456" {
			t.Errorf("expected /api/v1/groups/456, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Group{
			ID:   456,
			Name: "Deleted Group",
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	group, err := service.Delete(context.Background(), 456)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if group.ID != 456 {
		t.Errorf("expected ID 456, got %d", group.ID)
	}
}

func TestGroupsService_ListMembers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/groups/456/users" {
			t.Errorf("expected /api/v1/groups/456/users, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]User{
			{ID: 100, Name: "User 1"},
			{ID: 101, Name: "User 2"},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	users, err := service.ListMembers(context.Background(), 456)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestGroupsService_AddMember(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if body["user_id"].(float64) != 100 {
			t.Errorf("expected user_id 100, got %v", body["user_id"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(GroupMembership{
			ID:      999,
			GroupID: 456,
			UserID:  100,
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	membership, err := service.AddMember(context.Background(), 456, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if membership.UserID != 100 {
		t.Errorf("expected user ID 100, got %d", membership.UserID)
	}
}

func TestGroupsService_ListCategoriesCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/123/group_categories" {
			t.Errorf("expected /api/v1/courses/123/group_categories, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]GroupCategory{
			{ID: 1, Name: "Category 1"},
			{ID: 2, Name: "Category 2"},
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	categories, err := service.ListCategoriesCourse(context.Background(), 123, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(categories) != 2 {
		t.Errorf("expected 2 categories, got %d", len(categories))
	}
}

func TestGroupsService_GetCategory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/group_categories/456" {
			t.Errorf("expected /api/v1/group_categories/456, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(GroupCategory{
			ID:         456,
			Name:       "Test Category",
			SelfSignup: "enabled",
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	category, err := service.GetCategory(context.Background(), 456)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if category.ID != 456 {
		t.Errorf("expected ID 456, got %d", category.ID)
	}

	if category.Name != "Test Category" {
		t.Errorf("expected 'Test Category', got %s", category.Name)
	}
}

func TestGroupsService_CreateCategoryCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if body["name"] != "New Category" {
			t.Errorf("expected name 'New Category', got %v", body["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(GroupCategory{
			ID:   789,
			Name: "New Category",
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	params := &CreateCategoryParams{
		Name:       "New Category",
		SelfSignup: "enabled",
	}

	category, err := service.CreateCategoryCourse(context.Background(), 123, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if category.ID != 789 {
		t.Errorf("expected ID 789, got %d", category.ID)
	}
}

func TestGroupsService_UpdateCategory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if body["name"] != "Updated Category" {
			t.Errorf("expected name 'Updated Category', got %v", body["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(GroupCategory{
			ID:   456,
			Name: "Updated Category",
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	name := "Updated Category"
	params := &UpdateCategoryParams{
		Name: &name,
	}

	category, err := service.UpdateCategory(context.Background(), 456, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if category.Name != "Updated Category" {
		t.Errorf("expected 'Updated Category', got %s", category.Name)
	}
}

func TestGroupsService_DeleteCategory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/group_categories/456" {
			t.Errorf("expected /api/v1/group_categories/456, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(GroupCategory{
			ID:   456,
			Name: "Deleted Category",
		})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	category, err := service.DeleteCategory(context.Background(), 456)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if category.ID != 456 {
		t.Errorf("expected ID 456, got %d", category.ID)
	}
}

func TestGroupsService_ListUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/users/42/groups" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Group{{ID: 10, Name: "Group A"}})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewGroupsService(client)
	groups, err := svc.ListUser(context.Background(), 42, nil)
	if err != nil {
		t.Fatalf("ListUser: %v", err)
	}
	if len(groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(groups))
	}
}

func TestGroupsService_ListUser_Self(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/users/self/groups" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Group{{ID: 11}})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewGroupsService(client)
	opts := &ListGroupsOptions{Include: []string{"context_info"}, Page: 1, PerPage: 5}
	groups, err := svc.ListUser(context.Background(), 0, opts)
	if err != nil {
		t.Fatalf("ListUser self: %v", err)
	}
	if len(groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(groups))
	}
}

func TestGroupsService_RemoveMember(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/memberships/99" {
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
	svc := NewGroupsService(client)
	if err := svc.RemoveMember(context.Background(), 5, 99); err != nil {
		t.Fatalf("RemoveMember: %v", err)
	}
}

func TestGroupsService_ListCategoriesAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/3/group_categories" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]GroupCategory{{ID: 20, Name: "Cat A"}})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewGroupsService(client)
	cats, err := svc.ListCategoriesAccount(context.Background(), 3, nil)
	if err != nil {
		t.Fatalf("ListCategoriesAccount: %v", err)
	}
	if len(cats) != 1 {
		t.Errorf("expected 1 category, got %d", len(cats))
	}
}

func TestGroupsService_ListCategoriesAccount_WithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		q := r.URL.Query()
		if q.Get("per_page") != "10" {
			t.Errorf("expected per_page=10, got %q", q.Get("per_page"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]GroupCategory{{ID: 21}})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewGroupsService(client)
	cats, err := svc.ListCategoriesAccount(context.Background(), 3, &ListCategoriesOptions{Page: 1, PerPage: 10})
	if err != nil {
		t.Fatalf("ListCategoriesAccount with opts: %v", err)
	}
	if len(cats) != 1 {
		t.Errorf("expected 1 category, got %d", len(cats))
	}
}

func TestGroupsService_CreateCategoryAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/accounts/3/group_categories" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(GroupCategory{ID: 30, Name: "New Category"})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewGroupsService(client)
	params := &CreateCategoryParams{
		Name:               "New Category",
		SelfSignup:         "enabled",
		AutoLeader:         "first",
		GroupLimit:         5,
		SISGroupCategoryID: "sis-123",
	}
	cat, err := svc.CreateCategoryAccount(context.Background(), 3, params)
	if err != nil {
		t.Fatalf("CreateCategoryAccount: %v", err)
	}
	if cat.ID != 30 {
		t.Errorf("expected ID 30, got %d", cat.ID)
	}
}

func TestGroupsService_ListGroupsInCategory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/group_categories/5/groups" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Group{{ID: 100, Name: "Group X"}, {ID: 101, Name: "Group Y"}})
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 10})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	svc := NewGroupsService(client)
	groups, err := svc.ListGroupsInCategory(context.Background(), 5)
	if err != nil {
		t.Fatalf("ListGroupsInCategory: %v", err)
	}
	if len(groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(groups))
	}
}

func TestNewGroupsService(t *testing.T) {
	client, err := NewClient(ClientConfig{
		BaseURL: "https://canvas.example.com",
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	service := NewGroupsService(client)
	if service == nil {
		t.Fatal("expected non-nil service")
		return
	}
	if service.client != client {
		t.Error("expected client to be set")
	}
}
