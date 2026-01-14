package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// OverridesService handles assignment override-related API calls
type OverridesService struct {
	client *Client
}

// NewOverridesService creates a new overrides service
func NewOverridesService(client *Client) *OverridesService {
	return &OverridesService{client: client}
}

// AssignmentOverrideListOptions holds options for listing assignment overrides
type AssignmentOverrideListOptions struct {
	Page    int
	PerPage int
}

// List retrieves overrides for an assignment
func (s *OverridesService) List(ctx context.Context, courseID, assignmentID int64, opts *AssignmentOverrideListOptions) ([]AssignmentOverride, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/assignments/%d/overrides", courseID, assignmentID)

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

	var overrides []AssignmentOverride
	if err := s.client.GetAllPages(ctx, path, &overrides); err != nil {
		return nil, err
	}

	return overrides, nil
}

// Get retrieves a single override
func (s *OverridesService) Get(ctx context.Context, courseID, assignmentID, overrideID int64) (*AssignmentOverride, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/assignments/%d/overrides/%d", courseID, assignmentID, overrideID)

	var override AssignmentOverride
	if err := s.client.GetJSON(ctx, path, &override); err != nil {
		return nil, err
	}

	return &override, nil
}

// AssignmentOverrideCreateParams holds parameters for creating an assignment override
type AssignmentOverrideCreateParams struct {
	StudentIDs      []int64
	GroupID         int64
	CourseSectionID int64
	Title           string
	DueAt           string
	UnlockAt        string
	LockAt          string
}

// Create creates a new override for an assignment
func (s *OverridesService) Create(ctx context.Context, courseID, assignmentID int64, params *AssignmentOverrideCreateParams) (*AssignmentOverride, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/assignments/%d/overrides", courseID, assignmentID)

	body := map[string]interface{}{
		"assignment_override": make(map[string]interface{}),
	}

	overrideData, ok := body["assignment_override"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("internal error: invalid override data structure")
	}

	if len(params.StudentIDs) > 0 {
		overrideData["student_ids"] = params.StudentIDs
	}

	if params.GroupID > 0 {
		overrideData["group_id"] = params.GroupID
	}

	if params.CourseSectionID > 0 {
		overrideData["course_section_id"] = params.CourseSectionID
	}

	if params.Title != "" {
		overrideData["title"] = params.Title
	}

	if params.DueAt != "" {
		overrideData["due_at"] = params.DueAt
	}

	if params.UnlockAt != "" {
		overrideData["unlock_at"] = params.UnlockAt
	}

	if params.LockAt != "" {
		overrideData["lock_at"] = params.LockAt
	}

	var override AssignmentOverride
	if err := s.client.PostJSON(ctx, path, body, &override); err != nil {
		return nil, err
	}

	return &override, nil
}

// AssignmentOverrideUpdateParams holds parameters for updating an assignment override
type AssignmentOverrideUpdateParams struct {
	StudentIDs *[]int64
	Title      *string
	DueAt      *string
	UnlockAt   *string
	LockAt     *string
}

// Update updates an existing override
func (s *OverridesService) Update(ctx context.Context, courseID, assignmentID, overrideID int64, params *AssignmentOverrideUpdateParams) (*AssignmentOverride, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/assignments/%d/overrides/%d", courseID, assignmentID, overrideID)

	body := map[string]interface{}{
		"assignment_override": make(map[string]interface{}),
	}

	overrideData, ok := body["assignment_override"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("internal error: invalid override data structure")
	}

	if params.StudentIDs != nil {
		overrideData["student_ids"] = *params.StudentIDs
	}

	if params.Title != nil {
		overrideData["title"] = *params.Title
	}

	if params.DueAt != nil {
		overrideData["due_at"] = *params.DueAt
	}

	if params.UnlockAt != nil {
		overrideData["unlock_at"] = *params.UnlockAt
	}

	if params.LockAt != nil {
		overrideData["lock_at"] = *params.LockAt
	}

	var override AssignmentOverride
	if err := s.client.PutJSON(ctx, path, body, &override); err != nil {
		return nil, err
	}

	return &override, nil
}

// Delete deletes an override
func (s *OverridesService) Delete(ctx context.Context, courseID, assignmentID, overrideID int64) (*AssignmentOverride, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/assignments/%d/overrides/%d", courseID, assignmentID, overrideID)

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var override AssignmentOverride
	if err := json.NewDecoder(resp.Body).Decode(&override); err != nil {
		return nil, err
	}

	return &override, nil
}

// AssignmentOverrideBatchParams holds parameters for batch operations
type AssignmentOverrideBatchParams struct {
	AssignmentID int64
	StudentIDs   []int64
	GroupID      int64
	SectionID    int64
	Title        string
	DueAt        string
	UnlockAt     string
	LockAt       string
}

// BatchCreate creates multiple overrides across assignments
func (s *OverridesService) BatchCreate(ctx context.Context, courseID int64, overrides []AssignmentOverrideBatchParams) ([]AssignmentOverride, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/assignments/overrides", courseID)

	body := map[string]interface{}{
		"assignment_overrides": make([]map[string]interface{}, len(overrides)),
	}

	overridesData := make([]map[string]interface{}, len(overrides))
	for i, o := range overrides {
		overrideData := make(map[string]interface{})

		if o.AssignmentID > 0 {
			overrideData["assignment_id"] = o.AssignmentID
		}

		if len(o.StudentIDs) > 0 {
			overrideData["student_ids"] = o.StudentIDs
		}

		if o.GroupID > 0 {
			overrideData["group_id"] = o.GroupID
		}

		if o.SectionID > 0 {
			overrideData["course_section_id"] = o.SectionID
		}

		if o.Title != "" {
			overrideData["title"] = o.Title
		}

		if o.DueAt != "" {
			overrideData["due_at"] = o.DueAt
		}

		if o.UnlockAt != "" {
			overrideData["unlock_at"] = o.UnlockAt
		}

		if o.LockAt != "" {
			overrideData["lock_at"] = o.LockAt
		}

		overridesData[i] = overrideData
	}
	body["assignment_overrides"] = overridesData

	var result []AssignmentOverride
	if err := s.client.PostJSON(ctx, path, body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// BatchUpdate updates multiple overrides
func (s *OverridesService) BatchUpdate(ctx context.Context, courseID int64, overrides []AssignmentOverrideBatchParams) ([]AssignmentOverride, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/assignments/overrides", courseID)

	body := map[string]interface{}{
		"assignment_overrides": make([]map[string]interface{}, len(overrides)),
	}

	overridesData := make([]map[string]interface{}, len(overrides))
	for i, o := range overrides {
		overrideData := make(map[string]interface{})

		if o.AssignmentID > 0 {
			overrideData["assignment_id"] = o.AssignmentID
		}

		if len(o.StudentIDs) > 0 {
			overrideData["student_ids"] = o.StudentIDs
		}

		if o.Title != "" {
			overrideData["title"] = o.Title
		}

		if o.DueAt != "" {
			overrideData["due_at"] = o.DueAt
		}

		if o.UnlockAt != "" {
			overrideData["unlock_at"] = o.UnlockAt
		}

		if o.LockAt != "" {
			overrideData["lock_at"] = o.LockAt
		}

		overridesData[i] = overrideData
	}
	body["assignment_overrides"] = overridesData

	var result []AssignmentOverride
	if err := s.client.PutJSON(ctx, path, body, &result); err != nil {
		return nil, err
	}

	return result, nil
}
