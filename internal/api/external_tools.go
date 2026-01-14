package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// ExternalToolsService handles external tool-related API calls
type ExternalToolsService struct {
	client *Client
}

// NewExternalToolsService creates a new external tools service
func NewExternalToolsService(client *Client) *ExternalToolsService {
	return &ExternalToolsService{client: client}
}

// ExternalTool represents a Canvas external tool (LTI)
type ExternalTool struct {
	ID                      int64             `json:"id"`
	Name                    string            `json:"name"`
	Description             string            `json:"description,omitempty"`
	URL                     string            `json:"url,omitempty"`
	Domain                  string            `json:"domain,omitempty"`
	ConsumerKey             string            `json:"consumer_key,omitempty"`
	CreatedAt               string            `json:"created_at,omitempty"`
	UpdatedAt               string            `json:"updated_at,omitempty"`
	PrivacyLevel            string            `json:"privacy_level,omitempty"`
	CustomFields            map[string]string `json:"custom_fields,omitempty"`
	WorkflowState           string            `json:"workflow_state,omitempty"`
	VendorHelpLink          string            `json:"vendor_help_link,omitempty"`
	IconURL                 string            `json:"icon_url,omitempty"`
	VersionNumber           string            `json:"version_number,omitempty"`
	ResourceSelectionWidth  int               `json:"resource_selection_width,omitempty"`
	ResourceSelectionHeight int               `json:"resource_selection_height,omitempty"`
	DeploymentID            string            `json:"deployment_id,omitempty"`
	NotSelectable           bool              `json:"not_selectable,omitempty"`
	Selectable              bool              `json:"selectable,omitempty"`
	IsRCEFavorite           bool              `json:"is_rce_favorite,omitempty"`

	// Placement settings
	CourseNavigation    *ToolPlacement `json:"course_navigation,omitempty"`
	AccountNavigation   *ToolPlacement `json:"account_navigation,omitempty"`
	UserNavigation      *ToolPlacement `json:"user_navigation,omitempty"`
	EditorButton        *ToolPlacement `json:"editor_button,omitempty"`
	Homework            *ToolPlacement `json:"homework_submission,omitempty"`
	LinkSelection       *ToolPlacement `json:"link_selection,omitempty"`
	MigrationSelection  *ToolPlacement `json:"migration_selection,omitempty"`
	ResourceSelection   *ToolPlacement `json:"resource_selection,omitempty"`
	ToolConfiguration   *ToolPlacement `json:"tool_configuration,omitempty"`
	AssignmentSelection *ToolPlacement `json:"assignment_selection,omitempty"`
	ModuleItemSelection *ToolPlacement `json:"module_item_selection,omitempty"`
}

// ToolPlacement represents a tool placement configuration
type ToolPlacement struct {
	Enabled             bool              `json:"enabled,omitempty"`
	URL                 string            `json:"url,omitempty"`
	MessageType         string            `json:"message_type,omitempty"`
	Text                string            `json:"text,omitempty"`
	IconURL             string            `json:"icon_url,omitempty"`
	SelectionWidth      int               `json:"selection_width,omitempty"`
	SelectionHeight     int               `json:"selection_height,omitempty"`
	Labels              map[string]string `json:"labels,omitempty"`
	Visibility          string            `json:"visibility,omitempty"`
	DisplayType         string            `json:"display_type,omitempty"`
	LaunchWidth         int               `json:"launch_width,omitempty"`
	LaunchHeight        int               `json:"launch_height,omitempty"`
	WindowTarget        string            `json:"windowTarget,omitempty"`
	Default             string            `json:"default,omitempty"`
	AcceptMediaTypes    string            `json:"accept_media_types,omitempty"`
	CanvasIconClass     string            `json:"canvas_icon_class,omitempty"`
	RequiredPermissions []string          `json:"required_permissions,omitempty"`
}

// SessionlessLaunchURL represents the response from sessionless launch
type SessionlessLaunchURL struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

// ListExternalToolsOptions holds options for listing external tools
type ListExternalToolsOptions struct {
	Search         string
	Selectable     *bool
	IncludeParents bool
	Placements     []string
	Page           int
	PerPage        int
}

// ListByCourse retrieves external tools for a course
func (s *ExternalToolsService) ListByCourse(ctx context.Context, courseID int64, opts *ListExternalToolsOptions) ([]ExternalTool, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/external_tools", courseID)
	return s.list(ctx, path, opts)
}

// ListByAccount retrieves external tools for an account
func (s *ExternalToolsService) ListByAccount(ctx context.Context, accountID int64, opts *ListExternalToolsOptions) ([]ExternalTool, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/external_tools", accountID)
	return s.list(ctx, path, opts)
}

func (s *ExternalToolsService) list(ctx context.Context, path string, opts *ListExternalToolsOptions) ([]ExternalTool, error) {
	if opts != nil {
		query := url.Values{}

		if opts.Search != "" {
			query.Add("search_term", opts.Search)
		}

		if opts.Selectable != nil {
			query.Add("selectable", strconv.FormatBool(*opts.Selectable))
		}

		if opts.IncludeParents {
			query.Add("include_parents", "true")
		}

		for _, placement := range opts.Placements {
			query.Add("placement", placement)
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

	var tools []ExternalTool
	if err := s.client.GetAllPages(ctx, path, &tools); err != nil {
		return nil, err
	}

	return tools, nil
}

// GetByCourse retrieves a single external tool by ID in a course
func (s *ExternalToolsService) GetByCourse(ctx context.Context, courseID, toolID int64) (*ExternalTool, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/external_tools/%d", courseID, toolID)
	return s.get(ctx, path)
}

// GetByAccount retrieves a single external tool by ID in an account
func (s *ExternalToolsService) GetByAccount(ctx context.Context, accountID, toolID int64) (*ExternalTool, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/external_tools/%d", accountID, toolID)
	return s.get(ctx, path)
}

func (s *ExternalToolsService) get(ctx context.Context, path string) (*ExternalTool, error) {
	var tool ExternalTool
	if err := s.client.GetJSON(ctx, path, &tool); err != nil {
		return nil, err
	}

	return &tool, nil
}

// CreateExternalToolParams holds parameters for creating an external tool
type CreateExternalToolParams struct {
	Name           string            `json:"name"`
	PrivacyLevel   string            `json:"privacy_level,omitempty"`
	ConsumerKey    string            `json:"consumer_key,omitempty"`
	SharedSecret   string            `json:"shared_secret,omitempty"`
	URL            string            `json:"url,omitempty"`
	Domain         string            `json:"domain,omitempty"`
	Description    string            `json:"description,omitempty"`
	CustomFields   map[string]string `json:"custom_fields,omitempty"`
	ConfigType     string            `json:"config_type,omitempty"`
	ConfigURL      string            `json:"config_url,omitempty"`
	ConfigXML      string            `json:"config_xml,omitempty"`
	NotSelectable  *bool             `json:"not_selectable,omitempty"`
	IsRCEFavorite  *bool             `json:"is_rce_favorite,omitempty"`
	IconURL        string            `json:"icon_url,omitempty"`
	VendorHelpLink string            `json:"vendor_help_link,omitempty"`
	Text           string            `json:"text,omitempty"`
	LTIVersion     string            `json:"lti_version,omitempty"`
	ClientID       string            `json:"client_id,omitempty"`
}

// CreateInCourse creates a new external tool in a course
func (s *ExternalToolsService) CreateInCourse(ctx context.Context, courseID int64, params *CreateExternalToolParams) (*ExternalTool, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/external_tools", courseID)
	return s.create(ctx, path, params)
}

// CreateInAccount creates a new external tool in an account
func (s *ExternalToolsService) CreateInAccount(ctx context.Context, accountID int64, params *CreateExternalToolParams) (*ExternalTool, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/external_tools", accountID)
	return s.create(ctx, path, params)
}

func (s *ExternalToolsService) create(ctx context.Context, path string, params *CreateExternalToolParams) (*ExternalTool, error) {
	var tool ExternalTool
	if err := s.client.PostJSON(ctx, path, params, &tool); err != nil {
		return nil, err
	}

	return &tool, nil
}

// UpdateExternalToolParams holds parameters for updating an external tool
type UpdateExternalToolParams struct {
	Name           *string           `json:"name,omitempty"`
	PrivacyLevel   *string           `json:"privacy_level,omitempty"`
	ConsumerKey    *string           `json:"consumer_key,omitempty"`
	SharedSecret   *string           `json:"shared_secret,omitempty"`
	URL            *string           `json:"url,omitempty"`
	Domain         *string           `json:"domain,omitempty"`
	Description    *string           `json:"description,omitempty"`
	CustomFields   map[string]string `json:"custom_fields,omitempty"`
	NotSelectable  *bool             `json:"not_selectable,omitempty"`
	IsRCEFavorite  *bool             `json:"is_rce_favorite,omitempty"`
	IconURL        *string           `json:"icon_url,omitempty"`
	VendorHelpLink *string           `json:"vendor_help_link,omitempty"`
	Text           *string           `json:"text,omitempty"`
}

// UpdateInCourse updates an external tool in a course
func (s *ExternalToolsService) UpdateInCourse(ctx context.Context, courseID, toolID int64, params *UpdateExternalToolParams) (*ExternalTool, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/external_tools/%d", courseID, toolID)
	return s.update(ctx, path, params)
}

// UpdateInAccount updates an external tool in an account
func (s *ExternalToolsService) UpdateInAccount(ctx context.Context, accountID, toolID int64, params *UpdateExternalToolParams) (*ExternalTool, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/external_tools/%d", accountID, toolID)
	return s.update(ctx, path, params)
}

func (s *ExternalToolsService) update(ctx context.Context, path string, params *UpdateExternalToolParams) (*ExternalTool, error) {
	var tool ExternalTool
	if err := s.client.PutJSON(ctx, path, params, &tool); err != nil {
		return nil, err
	}

	return &tool, nil
}

// DeleteFromCourse deletes an external tool from a course
func (s *ExternalToolsService) DeleteFromCourse(ctx context.Context, courseID, toolID int64) (*ExternalTool, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/external_tools/%d", courseID, toolID)
	return s.delete(ctx, path)
}

// DeleteFromAccount deletes an external tool from an account
func (s *ExternalToolsService) DeleteFromAccount(ctx context.Context, accountID, toolID int64) (*ExternalTool, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/external_tools/%d", accountID, toolID)
	return s.delete(ctx, path)
}

func (s *ExternalToolsService) delete(ctx context.Context, path string) (*ExternalTool, error) {
	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tool ExternalTool
	if err := json.NewDecoder(resp.Body).Decode(&tool); err != nil {
		return nil, nil
	}

	return &tool, nil
}

// SessionlessLaunchParams holds parameters for sessionless launch
type SessionlessLaunchParams struct {
	LaunchType             string // assessment, module_item, course_navigation, account_navigation
	ID                     int64  // Tool ID (optional if using URL)
	URL                    string // Tool URL (optional if using ID)
	AssignmentID           int64  // For assessment launch type
	ModuleItemID           int64  // For module_item launch type
	ResourceLinkLookupUUID string
}

// GetSessionlessLaunchURLForCourse gets a sessionless launch URL for a tool in a course
func (s *ExternalToolsService) GetSessionlessLaunchURLForCourse(ctx context.Context, courseID int64, params *SessionlessLaunchParams) (*SessionlessLaunchURL, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/external_tools/sessionless_launch", courseID)
	return s.getSessionlessLaunch(ctx, path, params)
}

// GetSessionlessLaunchURLForAccount gets a sessionless launch URL for a tool in an account
func (s *ExternalToolsService) GetSessionlessLaunchURLForAccount(ctx context.Context, accountID int64, params *SessionlessLaunchParams) (*SessionlessLaunchURL, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/external_tools/sessionless_launch", accountID)
	return s.getSessionlessLaunch(ctx, path, params)
}

func (s *ExternalToolsService) getSessionlessLaunch(ctx context.Context, path string, params *SessionlessLaunchParams) (*SessionlessLaunchURL, error) {
	query := url.Values{}

	if params != nil {
		if params.LaunchType != "" {
			query.Add("launch_type", params.LaunchType)
		}

		if params.ID > 0 {
			query.Add("id", strconv.FormatInt(params.ID, 10))
		}

		if params.URL != "" {
			query.Add("url", params.URL)
		}

		if params.AssignmentID > 0 {
			query.Add("assignment_id", strconv.FormatInt(params.AssignmentID, 10))
		}

		if params.ModuleItemID > 0 {
			query.Add("module_item_id", strconv.FormatInt(params.ModuleItemID, 10))
		}

		if params.ResourceLinkLookupUUID != "" {
			query.Add("resource_link_lookup_uuid", params.ResourceLinkLookupUUID)
		}
	}

	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	var result SessionlessLaunchURL
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
