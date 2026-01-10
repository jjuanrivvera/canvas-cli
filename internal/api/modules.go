package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// Module represents a Canvas course module
type Module struct {
	ID                        int64        `json:"id"`
	WorkflowState             string       `json:"workflow_state"`
	Position                  int          `json:"position"`
	Name                      string       `json:"name"`
	UnlockAt                  *time.Time   `json:"unlock_at,omitempty"`
	RequireSequentialProgress bool         `json:"require_sequential_progress"`
	RequirementType           string       `json:"requirement_type"`
	PrerequisiteModuleIDs     []int64      `json:"prerequisite_module_ids"`
	ItemsCount                int          `json:"items_count"`
	ItemsURL                  string       `json:"items_url"`
	Items                     []ModuleItem `json:"items,omitempty"`
	State                     string       `json:"state,omitempty"`
	CompletedAt               *time.Time   `json:"completed_at,omitempty"`
	PublishFinalGrade         bool         `json:"publish_final_grade"`
	Published                 bool         `json:"published"`
}

// ModuleItem represents an item within a module
type ModuleItem struct {
	ID                    int64                  `json:"id"`
	ModuleID              int64                  `json:"module_id"`
	Position              int                    `json:"position"`
	Title                 string                 `json:"title"`
	Indent                int                    `json:"indent"`
	Type                  string                 `json:"type"`
	ContentID             int64                  `json:"content_id,omitempty"`
	HTMLURL               string                 `json:"html_url"`
	URL                   string                 `json:"url,omitempty"`
	PageURL               string                 `json:"page_url,omitempty"`
	ExternalURL           string                 `json:"external_url,omitempty"`
	NewTab                bool                   `json:"new_tab,omitempty"`
	CompletionRequirement *CompletionRequirement `json:"completion_requirement,omitempty"`
	ContentDetails        *ContentDetails        `json:"content_details,omitempty"`
	Published             bool                   `json:"published"`
}

// CompletionRequirement represents how a module item must be completed
type CompletionRequirement struct {
	Type          string  `json:"type"`
	MinScore      float64 `json:"min_score,omitempty"`
	MinPercentage float64 `json:"min_percentage,omitempty"`
	Completed     bool    `json:"completed,omitempty"`
}

// ContentDetails represents additional details for module item content
type ContentDetails struct {
	PointsPossible  float64    `json:"points_possible,omitempty"`
	DueAt           *time.Time `json:"due_at,omitempty"`
	UnlockAt        *time.Time `json:"unlock_at,omitempty"`
	LockAt          *time.Time `json:"lock_at,omitempty"`
	LockedForUser   bool       `json:"locked_for_user"`
	LockExplanation string     `json:"lock_explanation,omitempty"`
}

// ModulesService handles module-related API calls
type ModulesService struct {
	client *Client
}

// NewModulesService creates a new modules service
func NewModulesService(client *Client) *ModulesService {
	return &ModulesService{client: client}
}

// ListModulesOptions holds options for listing modules
type ListModulesOptions struct {
	Include    []string // items, content_details
	SearchTerm string
	StudentID  string
	Page       int
	PerPage    int
}

// List retrieves all modules for a course
func (s *ModulesService) List(ctx context.Context, courseID int64, opts *ListModulesOptions) ([]Module, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/modules", courseID)

	if opts != nil {
		query := url.Values{}

		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}

		if opts.SearchTerm != "" {
			query.Add("search_term", opts.SearchTerm)
		}

		if opts.StudentID != "" {
			query.Add("student_id", opts.StudentID)
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

	var modules []Module
	if err := s.client.GetAllPages(ctx, path, &modules); err != nil {
		return nil, err
	}

	return modules, nil
}

// Get retrieves a single module by ID
func (s *ModulesService) Get(ctx context.Context, courseID, moduleID int64, include []string, studentID string) (*Module, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/modules/%d", courseID, moduleID)

	query := url.Values{}
	for _, inc := range include {
		query.Add("include[]", inc)
	}
	if studentID != "" {
		query.Add("student_id", studentID)
	}
	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	var module Module
	if err := s.client.GetJSON(ctx, path, &module); err != nil {
		return nil, err
	}

	return &module, nil
}

// CreateModuleParams holds parameters for creating a module
type CreateModuleParams struct {
	Name                      string
	UnlockAt                  string
	Position                  int
	RequireSequentialProgress bool
	PrerequisiteModuleIDs     []int64
	PublishFinalGrade         bool
}

// Create creates a new module in a course
func (s *ModulesService) Create(ctx context.Context, courseID int64, params *CreateModuleParams) (*Module, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/modules", courseID)

	body := map[string]interface{}{
		"module": make(map[string]interface{}),
	}

	moduleData := body["module"].(map[string]interface{})
	moduleData["name"] = params.Name

	if params.UnlockAt != "" {
		moduleData["unlock_at"] = params.UnlockAt
	}

	if params.Position > 0 {
		moduleData["position"] = params.Position
	}

	if params.RequireSequentialProgress {
		moduleData["require_sequential_progress"] = true
	}

	if len(params.PrerequisiteModuleIDs) > 0 {
		moduleData["prerequisite_module_ids"] = params.PrerequisiteModuleIDs
	}

	if params.PublishFinalGrade {
		moduleData["publish_final_grade"] = true
	}

	var module Module
	if err := s.client.PostJSON(ctx, path, body, &module); err != nil {
		return nil, err
	}

	return &module, nil
}

// UpdateModuleParams holds parameters for updating a module
type UpdateModuleParams struct {
	Name                      *string
	UnlockAt                  *string
	Position                  *int
	RequireSequentialProgress *bool
	PrerequisiteModuleIDs     []int64
	PublishFinalGrade         *bool
	Published                 *bool
}

// Update updates an existing module
func (s *ModulesService) Update(ctx context.Context, courseID, moduleID int64, params *UpdateModuleParams) (*Module, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/modules/%d", courseID, moduleID)

	body := map[string]interface{}{
		"module": make(map[string]interface{}),
	}

	moduleData := body["module"].(map[string]interface{})

	if params.Name != nil {
		moduleData["name"] = *params.Name
	}

	if params.UnlockAt != nil {
		moduleData["unlock_at"] = *params.UnlockAt
	}

	if params.Position != nil {
		moduleData["position"] = *params.Position
	}

	if params.RequireSequentialProgress != nil {
		moduleData["require_sequential_progress"] = *params.RequireSequentialProgress
	}

	if params.PrerequisiteModuleIDs != nil {
		moduleData["prerequisite_module_ids"] = params.PrerequisiteModuleIDs
	}

	if params.PublishFinalGrade != nil {
		moduleData["publish_final_grade"] = *params.PublishFinalGrade
	}

	if params.Published != nil {
		moduleData["published"] = *params.Published
	}

	var module Module
	if err := s.client.PutJSON(ctx, path, body, &module); err != nil {
		return nil, err
	}

	return &module, nil
}

// Delete deletes a module
func (s *ModulesService) Delete(ctx context.Context, courseID, moduleID int64) error {
	path := fmt.Sprintf("/api/v1/courses/%d/modules/%d", courseID, moduleID)
	_, err := s.client.Delete(ctx, path)
	return err
}

// Relock re-locks module progressions
func (s *ModulesService) Relock(ctx context.Context, courseID, moduleID int64) (*Module, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/modules/%d/relock", courseID, moduleID)

	var module Module
	if err := s.client.PutJSON(ctx, path, nil, &module); err != nil {
		return nil, err
	}

	return &module, nil
}

// ListModuleItemsOptions holds options for listing module items
type ListModuleItemsOptions struct {
	Include    []string // content_details
	SearchTerm string
	StudentID  string
	Page       int
	PerPage    int
}

// ListItems retrieves all items for a module
func (s *ModulesService) ListItems(ctx context.Context, courseID, moduleID int64, opts *ListModuleItemsOptions) ([]ModuleItem, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/modules/%d/items", courseID, moduleID)

	if opts != nil {
		query := url.Values{}

		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}

		if opts.SearchTerm != "" {
			query.Add("search_term", opts.SearchTerm)
		}

		if opts.StudentID != "" {
			query.Add("student_id", opts.StudentID)
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

	var items []ModuleItem
	if err := s.client.GetAllPages(ctx, path, &items); err != nil {
		return nil, err
	}

	return items, nil
}

// GetItem retrieves a single module item
func (s *ModulesService) GetItem(ctx context.Context, courseID, moduleID, itemID int64, include []string, studentID string) (*ModuleItem, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/modules/%d/items/%d", courseID, moduleID, itemID)

	query := url.Values{}
	for _, inc := range include {
		query.Add("include[]", inc)
	}
	if studentID != "" {
		query.Add("student_id", studentID)
	}
	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	var item ModuleItem
	if err := s.client.GetJSON(ctx, path, &item); err != nil {
		return nil, err
	}

	return &item, nil
}

// CreateModuleItemParams holds parameters for creating a module item
type CreateModuleItemParams struct {
	Title                 string
	Type                  string // File, Page, Discussion, Assignment, Quiz, SubHeader, ExternalUrl, ExternalTool
	ContentID             int64
	Position              int
	Indent                int
	PageURL               string
	ExternalURL           string
	NewTab                bool
	CompletionRequirement *CompletionRequirementParams
	IframeWidth           int
	IframeHeight          int
}

// CompletionRequirementParams holds completion requirement parameters
type CompletionRequirementParams struct {
	Type     string // must_view, must_contribute, must_submit, must_mark_done, min_score
	MinScore float64
}

// CreateItem creates a new module item
func (s *ModulesService) CreateItem(ctx context.Context, courseID, moduleID int64, params *CreateModuleItemParams) (*ModuleItem, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/modules/%d/items", courseID, moduleID)

	body := map[string]interface{}{
		"module_item": make(map[string]interface{}),
	}

	itemData := body["module_item"].(map[string]interface{})
	itemData["type"] = params.Type

	if params.Title != "" {
		itemData["title"] = params.Title
	}

	if params.ContentID > 0 {
		itemData["content_id"] = params.ContentID
	}

	if params.Position > 0 {
		itemData["position"] = params.Position
	}

	if params.Indent > 0 {
		itemData["indent"] = params.Indent
	}

	if params.PageURL != "" {
		itemData["page_url"] = params.PageURL
	}

	if params.ExternalURL != "" {
		itemData["external_url"] = params.ExternalURL
	}

	if params.NewTab {
		itemData["new_tab"] = true
	}

	if params.CompletionRequirement != nil {
		reqData := map[string]interface{}{
			"type": params.CompletionRequirement.Type,
		}
		if params.CompletionRequirement.MinScore > 0 {
			reqData["min_score"] = params.CompletionRequirement.MinScore
		}
		itemData["completion_requirement"] = reqData
	}

	if params.IframeWidth > 0 || params.IframeHeight > 0 {
		iframeData := make(map[string]interface{})
		if params.IframeWidth > 0 {
			iframeData["width"] = params.IframeWidth
		}
		if params.IframeHeight > 0 {
			iframeData["height"] = params.IframeHeight
		}
		itemData["iframe"] = iframeData
	}

	var item ModuleItem
	if err := s.client.PostJSON(ctx, path, body, &item); err != nil {
		return nil, err
	}

	return &item, nil
}

// UpdateModuleItemParams holds parameters for updating a module item
type UpdateModuleItemParams struct {
	Title                 *string
	Position              *int
	Indent                *int
	ExternalURL           *string
	NewTab                *bool
	CompletionRequirement *CompletionRequirementParams
	Published             *bool
	MoveToModuleID        *int64
}

// UpdateItem updates an existing module item
func (s *ModulesService) UpdateItem(ctx context.Context, courseID, moduleID, itemID int64, params *UpdateModuleItemParams) (*ModuleItem, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/modules/%d/items/%d", courseID, moduleID, itemID)

	body := map[string]interface{}{
		"module_item": make(map[string]interface{}),
	}

	itemData := body["module_item"].(map[string]interface{})

	if params.Title != nil {
		itemData["title"] = *params.Title
	}

	if params.Position != nil {
		itemData["position"] = *params.Position
	}

	if params.Indent != nil {
		itemData["indent"] = *params.Indent
	}

	if params.ExternalURL != nil {
		itemData["external_url"] = *params.ExternalURL
	}

	if params.NewTab != nil {
		itemData["new_tab"] = *params.NewTab
	}

	if params.CompletionRequirement != nil {
		reqData := map[string]interface{}{
			"type": params.CompletionRequirement.Type,
		}
		if params.CompletionRequirement.MinScore > 0 {
			reqData["min_score"] = params.CompletionRequirement.MinScore
		}
		itemData["completion_requirement"] = reqData
	}

	if params.Published != nil {
		itemData["published"] = *params.Published
	}

	if params.MoveToModuleID != nil {
		itemData["module_id"] = *params.MoveToModuleID
	}

	var item ModuleItem
	if err := s.client.PutJSON(ctx, path, body, &item); err != nil {
		return nil, err
	}

	return &item, nil
}

// DeleteItem deletes a module item
func (s *ModulesService) DeleteItem(ctx context.Context, courseID, moduleID, itemID int64) error {
	path := fmt.Sprintf("/api/v1/courses/%d/modules/%d/items/%d", courseID, moduleID, itemID)
	_, err := s.client.Delete(ctx, path)
	return err
}

// MarkItemDone marks a module item as done
func (s *ModulesService) MarkItemDone(ctx context.Context, courseID, moduleID, itemID int64) error {
	path := fmt.Sprintf("/api/v1/courses/%d/modules/%d/items/%d/done", courseID, moduleID, itemID)
	return s.client.PutJSON(ctx, path, nil, nil)
}

// MarkItemNotDone marks a module item as not done
func (s *ModulesService) MarkItemNotDone(ctx context.Context, courseID, moduleID, itemID int64) error {
	path := fmt.Sprintf("/api/v1/courses/%d/modules/%d/items/%d/done", courseID, moduleID, itemID)
	_, err := s.client.Delete(ctx, path)
	return err
}

// MarkItemRead marks a module item as read
func (s *ModulesService) MarkItemRead(ctx context.Context, courseID, moduleID, itemID int64) error {
	path := fmt.Sprintf("/api/v1/courses/%d/modules/%d/items/%d/mark_read", courseID, moduleID, itemID)
	return s.client.PostJSON(ctx, path, nil, nil)
}

// ModuleItemSequence represents a sequence of module items
type ModuleItemSequence struct {
	Items   []ModuleItemSequenceNode `json:"items"`
	Modules []ModuleReference        `json:"modules"`
}

// ModuleItemSequenceNode represents a node in the module item sequence
type ModuleItemSequenceNode struct {
	Prev        *ModuleItemRef `json:"prev,omitempty"`
	Current     *ModuleItemRef `json:"current,omitempty"`
	Next        *ModuleItemRef `json:"next,omitempty"`
	MasteryPath interface{}    `json:"mastery_path,omitempty"`
}

// ModuleItemRef is a reference to a module item in a sequence
type ModuleItemRef struct {
	ID       int64  `json:"id"`
	ModuleID int64  `json:"module_id"`
	Title    string `json:"title"`
	Type     string `json:"type"`
}

// ModuleReference is a reference to a module
type ModuleReference struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// GetItemSequence gets the module item sequence for an asset
func (s *ModulesService) GetItemSequence(ctx context.Context, courseID int64, assetType string, assetID int64) (*ModuleItemSequence, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/module_item_sequence", courseID)

	query := url.Values{}
	query.Add("asset_type", assetType)
	query.Add("asset_id", strconv.FormatInt(assetID, 10))
	path += "?" + query.Encode()

	var sequence ModuleItemSequence
	if err := s.client.GetJSON(ctx, path, &sequence); err != nil {
		return nil, err
	}

	return &sequence, nil
}
