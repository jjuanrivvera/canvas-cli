package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// RolesService handles role-related API calls
type RolesService struct {
	client *Client
}

// NewRolesService creates a new roles service
func NewRolesService(client *Client) *RolesService {
	return &RolesService{client: client}
}

// Role represents a Canvas role
type Role struct {
	ID            int64                 `json:"id"`
	Label         string                `json:"label"`
	BaseRoleType  string                `json:"base_role_type"`
	Account       *Account              `json:"account,omitempty"`
	WorkflowState string                `json:"workflow_state"`
	CreatedAt     string                `json:"created_at,omitempty"`
	LastUpdatedAt string                `json:"last_updated_at,omitempty"`
	Permissions   map[string]Permission `json:"permissions,omitempty"`
}

// Permission represents a role permission
type Permission struct {
	Enabled              bool  `json:"enabled"`
	Locked               bool  `json:"locked"`
	AppliesToSelf        bool  `json:"applies_to_self,omitempty"`
	AppliesToDescendants bool  `json:"applies_to_descendants,omitempty"`
	Explicit             bool  `json:"explicit,omitempty"`
	Prior                *bool `json:"prior_default,omitempty"`
}

// ListRolesOptions holds options for listing roles
type ListRolesOptions struct {
	State         string // active, inactive
	ShowInherited bool
	Page          int
	PerPage       int
}

// List retrieves roles for an account
func (s *RolesService) List(ctx context.Context, accountID int64, opts *ListRolesOptions) ([]Role, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/roles", accountID)

	if opts != nil {
		query := url.Values{}

		if opts.State != "" {
			query.Add("state", opts.State)
		}

		if opts.ShowInherited {
			query.Add("show_inherited", "true")
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

	var roles []Role
	if err := s.client.GetAllPages(ctx, path, &roles); err != nil {
		return nil, err
	}

	return roles, nil
}

// Get retrieves a single role
func (s *RolesService) Get(ctx context.Context, accountID, roleID int64) (*Role, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/roles/%d", accountID, roleID)

	var role Role
	if err := s.client.GetJSON(ctx, path, &role); err != nil {
		return nil, err
	}

	return &role, nil
}

// CreateRoleParams holds parameters for creating a role
type CreateRoleParams struct {
	Label        string
	BaseRoleType string // AccountMembership, StudentEnrollment, TeacherEnrollment, TaEnrollment, ObserverEnrollment, DesignerEnrollment
	Permissions  map[string]PermissionOverride
}

// PermissionOverride represents a permission override
type PermissionOverride struct {
	Enabled              *bool
	Locked               *bool
	AppliesToSelf        *bool
	AppliesToDescendants *bool
}

// Create creates a new role
func (s *RolesService) Create(ctx context.Context, accountID int64, params *CreateRoleParams) (*Role, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/roles", accountID)

	body := make(map[string]interface{})
	body["label"] = params.Label

	if params.BaseRoleType != "" {
		body["base_role_type"] = params.BaseRoleType
	}

	if len(params.Permissions) > 0 {
		perms := make(map[string]interface{})
		for name, override := range params.Permissions {
			permData := make(map[string]interface{})
			if override.Enabled != nil {
				permData["enabled"] = *override.Enabled
			}
			if override.Locked != nil {
				permData["locked"] = *override.Locked
			}
			if override.AppliesToSelf != nil {
				permData["applies_to_self"] = *override.AppliesToSelf
			}
			if override.AppliesToDescendants != nil {
				permData["applies_to_descendants"] = *override.AppliesToDescendants
			}
			perms[name] = permData
		}
		body["permissions"] = perms
	}

	var role Role
	if err := s.client.PostJSON(ctx, path, body, &role); err != nil {
		return nil, err
	}

	return &role, nil
}

// UpdateRoleParams holds parameters for updating a role
type UpdateRoleParams struct {
	Label       *string
	Permissions map[string]PermissionOverride
}

// Update updates a role
func (s *RolesService) Update(ctx context.Context, accountID, roleID int64, params *UpdateRoleParams) (*Role, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/roles/%d", accountID, roleID)

	body := make(map[string]interface{})

	if params.Label != nil {
		body["label"] = *params.Label
	}

	if len(params.Permissions) > 0 {
		perms := make(map[string]interface{})
		for name, override := range params.Permissions {
			permData := make(map[string]interface{})
			if override.Enabled != nil {
				permData["enabled"] = *override.Enabled
			}
			if override.Locked != nil {
				permData["locked"] = *override.Locked
			}
			if override.AppliesToSelf != nil {
				permData["applies_to_self"] = *override.AppliesToSelf
			}
			if override.AppliesToDescendants != nil {
				permData["applies_to_descendants"] = *override.AppliesToDescendants
			}
			perms[name] = permData
		}
		body["permissions"] = perms
	}

	var role Role
	if err := s.client.PutJSON(ctx, path, body, &role); err != nil {
		return nil, err
	}

	return &role, nil
}

// Deactivate deactivates a role
func (s *RolesService) Deactivate(ctx context.Context, accountID, roleID int64) (*Role, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/roles/%d", accountID, roleID)

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var role Role
	if err := json.NewDecoder(resp.Body).Decode(&role); err != nil {
		return nil, err
	}

	return &role, nil
}

// Activate activates a role
func (s *RolesService) Activate(ctx context.Context, accountID, roleID int64) (*Role, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/roles/%d/activate", accountID, roleID)

	var role Role
	if err := s.client.PostJSON(ctx, path, nil, &role); err != nil {
		return nil, err
	}

	return &role, nil
}
