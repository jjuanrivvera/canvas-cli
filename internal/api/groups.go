package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// GroupsService handles group-related API calls
type GroupsService struct {
	client *Client
}

// NewGroupsService creates a new groups service
func NewGroupsService(client *Client) *GroupsService {
	return &GroupsService{client: client}
}

// Group represents a Canvas group
type Group struct {
	ID              int64             `json:"id"`
	Name            string            `json:"name"`
	Description     string            `json:"description,omitempty"`
	IsPublic        bool              `json:"is_public"`
	FollowedByUser  bool              `json:"followed_by_user,omitempty"`
	JoinLevel       string            `json:"join_level,omitempty"`
	MembersCount    int               `json:"members_count,omitempty"`
	AvatarURL       string            `json:"avatar_url,omitempty"`
	ContextType     string            `json:"context_type,omitempty"`
	CourseID        int64             `json:"course_id,omitempty"`
	AccountID       int64             `json:"account_id,omitempty"`
	Role            string            `json:"role,omitempty"`
	GroupCategoryID int64             `json:"group_category_id,omitempty"`
	SISGroupID      string            `json:"sis_group_id,omitempty"`
	SISImportID     int64             `json:"sis_import_id,omitempty"`
	StorageQuotaMb  int64             `json:"storage_quota_mb,omitempty"`
	Permissions     *GroupPermissions `json:"permissions,omitempty"`
	Users           []User            `json:"users,omitempty"`
}

// GroupPermissions represents permissions on a group
type GroupPermissions struct {
	CreateDiscussionTopic bool `json:"create_discussion_topic"`
	CreateAnnouncement    bool `json:"create_announcement"`
}

// GroupCategory represents a Canvas group category
type GroupCategory struct {
	ID                 int64       `json:"id"`
	Name               string      `json:"name"`
	Role               string      `json:"role,omitempty"`
	SelfSignup         string      `json:"self_signup,omitempty"`
	AutoLeader         string      `json:"auto_leader,omitempty"`
	ContextType        string      `json:"context_type,omitempty"`
	AccountID          int64       `json:"account_id,omitempty"`
	CourseID           int64       `json:"course_id,omitempty"`
	GroupLimit         int         `json:"group_limit,omitempty"`
	SISGroupCategoryID string      `json:"sis_group_category_id,omitempty"`
	SISImportID        int64       `json:"sis_import_id,omitempty"`
	Progress           interface{} `json:"progress,omitempty"`
}

// GroupMembership represents a user's membership in a group
type GroupMembership struct {
	ID            int64  `json:"id"`
	GroupID       int64  `json:"group_id"`
	UserID        int64  `json:"user_id"`
	WorkflowState string `json:"workflow_state"`
	Moderator     bool   `json:"moderator"`
	JustCreated   bool   `json:"just_created,omitempty"`
	SISImportID   int64  `json:"sis_import_id,omitempty"`
}

// ListGroupsOptions holds options for listing groups
type ListGroupsOptions struct {
	Include []string // users, permissions, tabs
	Page    int
	PerPage int
}

// ListCourse retrieves all groups for a course
func (s *GroupsService) ListCourse(ctx context.Context, courseID int64, opts *ListGroupsOptions) ([]Group, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/groups", courseID)

	if opts != nil {
		query := url.Values{}

		for _, include := range opts.Include {
			query.Add("include[]", include)
		}

		if opts.Page > 0 {
			query.Add("page", strconv.Itoa(opts.Page))
		}

		if opts.PerPage > 0 {
			query.Add("per_page", strconv.Itoa(opts.PerPage))
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var groups []Group
	if err := s.client.GetAllPages(ctx, path, &groups); err != nil {
		return nil, err
	}

	return groups, nil
}

// ListAccount retrieves all groups for an account
func (s *GroupsService) ListAccount(ctx context.Context, accountID int64, opts *ListGroupsOptions) ([]Group, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/groups", accountID)

	if opts != nil {
		query := url.Values{}

		for _, include := range opts.Include {
			query.Add("include[]", include)
		}

		if opts.Page > 0 {
			query.Add("page", strconv.Itoa(opts.Page))
		}

		if opts.PerPage > 0 {
			query.Add("per_page", strconv.Itoa(opts.PerPage))
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var groups []Group
	if err := s.client.GetAllPages(ctx, path, &groups); err != nil {
		return nil, err
	}

	return groups, nil
}

// ListUser retrieves all groups for a user
func (s *GroupsService) ListUser(ctx context.Context, userID int64, opts *ListGroupsOptions) ([]Group, error) {
	var path string
	if userID > 0 {
		path = fmt.Sprintf("/api/v1/users/%d/groups", userID)
	} else {
		path = "/api/v1/users/self/groups"
	}

	if opts != nil {
		query := url.Values{}

		for _, include := range opts.Include {
			query.Add("include[]", include)
		}

		if opts.Page > 0 {
			query.Add("page", strconv.Itoa(opts.Page))
		}

		if opts.PerPage > 0 {
			query.Add("per_page", strconv.Itoa(opts.PerPage))
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var groups []Group
	if err := s.client.GetAllPages(ctx, path, &groups); err != nil {
		return nil, err
	}

	return groups, nil
}

// Get retrieves a single group
func (s *GroupsService) Get(ctx context.Context, groupID int64, include []string) (*Group, error) {
	path := fmt.Sprintf("/api/v1/groups/%d", groupID)

	if len(include) > 0 {
		query := url.Values{}
		for _, inc := range include {
			query.Add("include[]", inc)
		}
		path += "?" + query.Encode()
	}

	var group Group
	if err := s.client.GetJSON(ctx, path, &group); err != nil {
		return nil, err
	}

	return &group, nil
}

// CreateGroupParams holds parameters for creating a group
type CreateGroupParams struct {
	Name           string
	Description    string
	IsPublic       bool
	JoinLevel      string // parent_context_auto_join, parent_context_request, invitation_only
	StorageQuotaMb int64
	SISGroupID     string
}

// Create creates a new group
func (s *GroupsService) Create(ctx context.Context, categoryID int64, params *CreateGroupParams) (*Group, error) {
	path := fmt.Sprintf("/api/v1/group_categories/%d/groups", categoryID)

	body := make(map[string]interface{})

	if params.Name != "" {
		body["name"] = params.Name
	}

	if params.Description != "" {
		body["description"] = params.Description
	}

	if params.IsPublic {
		body["is_public"] = params.IsPublic
	}

	if params.JoinLevel != "" {
		body["join_level"] = params.JoinLevel
	}

	if params.StorageQuotaMb > 0 {
		body["storage_quota_mb"] = params.StorageQuotaMb
	}

	if params.SISGroupID != "" {
		body["sis_group_id"] = params.SISGroupID
	}

	var group Group
	if err := s.client.PostJSON(ctx, path, body, &group); err != nil {
		return nil, err
	}

	return &group, nil
}

// UpdateGroupParams holds parameters for updating a group
type UpdateGroupParams struct {
	Name           *string
	Description    *string
	IsPublic       *bool
	JoinLevel      *string
	AvatarID       *int64
	StorageQuotaMb *int64
	SISGroupID     *string
}

// Update updates an existing group
func (s *GroupsService) Update(ctx context.Context, groupID int64, params *UpdateGroupParams) (*Group, error) {
	path := fmt.Sprintf("/api/v1/groups/%d", groupID)

	body := make(map[string]interface{})

	if params.Name != nil {
		body["name"] = *params.Name
	}

	if params.Description != nil {
		body["description"] = *params.Description
	}

	if params.IsPublic != nil {
		body["is_public"] = *params.IsPublic
	}

	if params.JoinLevel != nil {
		body["join_level"] = *params.JoinLevel
	}

	if params.AvatarID != nil {
		body["avatar_id"] = *params.AvatarID
	}

	if params.StorageQuotaMb != nil {
		body["storage_quota_mb"] = *params.StorageQuotaMb
	}

	if params.SISGroupID != nil {
		body["sis_group_id"] = *params.SISGroupID
	}

	var group Group
	if err := s.client.PutJSON(ctx, path, body, &group); err != nil {
		return nil, err
	}

	return &group, nil
}

// Delete deletes a group
func (s *GroupsService) Delete(ctx context.Context, groupID int64) (*Group, error) {
	path := fmt.Sprintf("/api/v1/groups/%d", groupID)

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var group Group
	if err := json.NewDecoder(resp.Body).Decode(&group); err != nil {
		return nil, err
	}

	return &group, nil
}

// ListMembers retrieves all members of a group
func (s *GroupsService) ListMembers(ctx context.Context, groupID int64) ([]User, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/users", groupID)

	var users []User
	if err := s.client.GetAllPages(ctx, path, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// AddMember adds a user to a group
func (s *GroupsService) AddMember(ctx context.Context, groupID, userID int64) (*GroupMembership, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/memberships", groupID)

	body := map[string]interface{}{
		"user_id": userID,
	}

	var membership GroupMembership
	if err := s.client.PostJSON(ctx, path, body, &membership); err != nil {
		return nil, err
	}

	return &membership, nil
}

// RemoveMember removes a user from a group
func (s *GroupsService) RemoveMember(ctx context.Context, groupID, membershipID int64) error {
	path := fmt.Sprintf("/api/v1/groups/%d/memberships/%d", groupID, membershipID)

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// ListCategoriesOptions holds options for listing group categories
type ListCategoriesOptions struct {
	Page    int
	PerPage int
}

// ListCategoriesCourse retrieves group categories for a course
func (s *GroupsService) ListCategoriesCourse(ctx context.Context, courseID int64, opts *ListCategoriesOptions) ([]GroupCategory, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/group_categories", courseID)

	if opts != nil {
		query := url.Values{}

		if opts.Page > 0 {
			query.Add("page", strconv.Itoa(opts.Page))
		}

		if opts.PerPage > 0 {
			query.Add("per_page", strconv.Itoa(opts.PerPage))
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var categories []GroupCategory
	if err := s.client.GetAllPages(ctx, path, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

// ListCategoriesAccount retrieves group categories for an account
func (s *GroupsService) ListCategoriesAccount(ctx context.Context, accountID int64, opts *ListCategoriesOptions) ([]GroupCategory, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/group_categories", accountID)

	if opts != nil {
		query := url.Values{}

		if opts.Page > 0 {
			query.Add("page", strconv.Itoa(opts.Page))
		}

		if opts.PerPage > 0 {
			query.Add("per_page", strconv.Itoa(opts.PerPage))
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var categories []GroupCategory
	if err := s.client.GetAllPages(ctx, path, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

// GetCategory retrieves a single group category
func (s *GroupsService) GetCategory(ctx context.Context, categoryID int64) (*GroupCategory, error) {
	path := fmt.Sprintf("/api/v1/group_categories/%d", categoryID)

	var category GroupCategory
	if err := s.client.GetJSON(ctx, path, &category); err != nil {
		return nil, err
	}

	return &category, nil
}

// CreateCategoryParams holds parameters for creating a group category
type CreateCategoryParams struct {
	Name               string
	SelfSignup         string // enabled, restricted
	AutoLeader         string // first, random
	GroupLimit         int
	CreateGroupCount   int
	SplitGroupCount    int
	SISGroupCategoryID string
}

// CreateCategoryCourse creates a new group category in a course
func (s *GroupsService) CreateCategoryCourse(ctx context.Context, courseID int64, params *CreateCategoryParams) (*GroupCategory, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/group_categories", courseID)

	body := make(map[string]interface{})

	if params.Name != "" {
		body["name"] = params.Name
	}

	if params.SelfSignup != "" {
		body["self_signup"] = params.SelfSignup
	}

	if params.AutoLeader != "" {
		body["auto_leader"] = params.AutoLeader
	}

	if params.GroupLimit > 0 {
		body["group_limit"] = params.GroupLimit
	}

	if params.CreateGroupCount > 0 {
		body["create_group_count"] = params.CreateGroupCount
	}

	if params.SplitGroupCount > 0 {
		body["split_group_count"] = params.SplitGroupCount
	}

	if params.SISGroupCategoryID != "" {
		body["sis_group_category_id"] = params.SISGroupCategoryID
	}

	var category GroupCategory
	if err := s.client.PostJSON(ctx, path, body, &category); err != nil {
		return nil, err
	}

	return &category, nil
}

// CreateCategoryAccount creates a new group category in an account
func (s *GroupsService) CreateCategoryAccount(ctx context.Context, accountID int64, params *CreateCategoryParams) (*GroupCategory, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/group_categories", accountID)

	body := make(map[string]interface{})

	if params.Name != "" {
		body["name"] = params.Name
	}

	if params.SelfSignup != "" {
		body["self_signup"] = params.SelfSignup
	}

	if params.AutoLeader != "" {
		body["auto_leader"] = params.AutoLeader
	}

	if params.GroupLimit > 0 {
		body["group_limit"] = params.GroupLimit
	}

	if params.SISGroupCategoryID != "" {
		body["sis_group_category_id"] = params.SISGroupCategoryID
	}

	var category GroupCategory
	if err := s.client.PostJSON(ctx, path, body, &category); err != nil {
		return nil, err
	}

	return &category, nil
}

// UpdateCategoryParams holds parameters for updating a group category
type UpdateCategoryParams struct {
	Name               *string
	SelfSignup         *string
	AutoLeader         *string
	GroupLimit         *int
	SISGroupCategoryID *string
}

// UpdateCategory updates an existing group category
func (s *GroupsService) UpdateCategory(ctx context.Context, categoryID int64, params *UpdateCategoryParams) (*GroupCategory, error) {
	path := fmt.Sprintf("/api/v1/group_categories/%d", categoryID)

	body := make(map[string]interface{})

	if params.Name != nil {
		body["name"] = *params.Name
	}

	if params.SelfSignup != nil {
		body["self_signup"] = *params.SelfSignup
	}

	if params.AutoLeader != nil {
		body["auto_leader"] = *params.AutoLeader
	}

	if params.GroupLimit != nil {
		body["group_limit"] = *params.GroupLimit
	}

	if params.SISGroupCategoryID != nil {
		body["sis_group_category_id"] = *params.SISGroupCategoryID
	}

	var category GroupCategory
	if err := s.client.PutJSON(ctx, path, body, &category); err != nil {
		return nil, err
	}

	return &category, nil
}

// DeleteCategory deletes a group category
func (s *GroupsService) DeleteCategory(ctx context.Context, categoryID int64) (*GroupCategory, error) {
	path := fmt.Sprintf("/api/v1/group_categories/%d", categoryID)

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var category GroupCategory
	if err := json.NewDecoder(resp.Body).Decode(&category); err != nil {
		return nil, err
	}

	return &category, nil
}

// ListGroupsInCategory retrieves all groups in a category
func (s *GroupsService) ListGroupsInCategory(ctx context.Context, categoryID int64) ([]Group, error) {
	path := fmt.Sprintf("/api/v1/group_categories/%d/groups", categoryID)

	var groups []Group
	if err := s.client.GetAllPages(ctx, path, &groups); err != nil {
		return nil, err
	}

	return groups, nil
}
