package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// PlannerItem represents a planner item (assignment, quiz, calendar event, etc.)
type PlannerItem struct {
	CourseID       *int64      `json:"course_id,omitempty"`
	GroupID        *int64      `json:"group_id,omitempty"`
	UserID         *int64      `json:"user_id,omitempty"`
	ContextType    string      `json:"context_type,omitempty"`
	ContextName    string      `json:"context_name,omitempty"`
	PlannableType  string      `json:"plannable_type"`
	PlannableID    int64       `json:"plannable_id"`
	PlannableDate  *time.Time  `json:"plannable_date,omitempty"`
	Submissions    interface{} `json:"submissions,omitempty"`
	Plannable      interface{} `json:"plannable,omitempty"`
	HTMLURL        string      `json:"html_url,omitempty"`
	NewActivity    bool        `json:"new_activity"`
	ContextImage   string      `json:"context_image,omitempty"`
}

// PlannerNote represents a planner note
type PlannerNote struct {
	ID                  int64      `json:"id"`
	Title               string     `json:"title"`
	Description         string     `json:"description,omitempty"`
	UserID              int64      `json:"user_id"`
	CourseID            *int64     `json:"course_id,omitempty"`
	TodoDate            *time.Time `json:"todo_date,omitempty"`
	Details             string     `json:"details,omitempty"`
	WorkflowState       string     `json:"workflow_state"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	LinkedObjectType    string     `json:"linked_object_type,omitempty"`
	LinkedObjectID      *int64     `json:"linked_object_id,omitempty"`
	LinkedObjectHTMLURL string     `json:"linked_object_html_url,omitempty"`
	LinkedObjectURL     string     `json:"linked_object_url,omitempty"`
}

// PlannerOverride represents a planner override
type PlannerOverride struct {
	ID             int64      `json:"id"`
	PlannableType  string     `json:"plannable_type"`
	PlannableID    int64      `json:"plannable_id"`
	UserID         int64      `json:"user_id"`
	AssignmentID   *int64     `json:"assignment_id,omitempty"`
	WorkflowState  string     `json:"workflow_state"`
	MarkedComplete bool       `json:"marked_complete"`
	Dismissed      bool       `json:"dismissed"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// PlannerService handles planner-related API calls
type PlannerService struct {
	client *Client
}

// NewPlannerService creates a new planner service
func NewPlannerService(client *Client) *PlannerService {
	return &PlannerService{client: client}
}

// ListPlannerItemsOptions holds options for listing planner items
type ListPlannerItemsOptions struct {
	StartDate    string   // ISO 8601
	EndDate      string   // ISO 8601
	ContextCodes []string // course_123, group_456
	Filter       string   // all_ungraded_todo_items, all_assignments, all_quizzes, all_calendar_events, all_planner_notes
	Page         int
	PerPage      int
}

// ListItems retrieves planner items for the current user
func (s *PlannerService) ListItems(ctx context.Context, opts *ListPlannerItemsOptions) ([]PlannerItem, error) {
	path := "/api/v1/planner/items"

	if opts != nil {
		query := url.Values{}

		if opts.StartDate != "" {
			query.Add("start_date", opts.StartDate)
		}

		if opts.EndDate != "" {
			query.Add("end_date", opts.EndDate)
		}

		for _, code := range opts.ContextCodes {
			query.Add("context_codes[]", code)
		}

		if opts.Filter != "" {
			query.Add("filter", opts.Filter)
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

	var items []PlannerItem
	if err := s.client.GetAllPages(ctx, path, &items); err != nil {
		return nil, err
	}

	return items, nil
}

// ListPlannerNotesOptions holds options for listing planner notes
type ListPlannerNotesOptions struct {
	StartDate    string
	EndDate      string
	ContextCodes []string
	CourseID     int64
	Page         int
	PerPage      int
}

// ListNotes retrieves planner notes for the current user
func (s *PlannerService) ListNotes(ctx context.Context, opts *ListPlannerNotesOptions) ([]PlannerNote, error) {
	path := "/api/v1/planner_notes"

	if opts != nil {
		query := url.Values{}

		if opts.StartDate != "" {
			query.Add("start_date", opts.StartDate)
		}

		if opts.EndDate != "" {
			query.Add("end_date", opts.EndDate)
		}

		for _, code := range opts.ContextCodes {
			query.Add("context_codes[]", code)
		}

		if opts.CourseID > 0 {
			query.Add("course_id", strconv.FormatInt(opts.CourseID, 10))
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

	var notes []PlannerNote
	if err := s.client.GetAllPages(ctx, path, &notes); err != nil {
		return nil, err
	}

	return notes, nil
}

// GetNote retrieves a single planner note
func (s *PlannerService) GetNote(ctx context.Context, noteID int64) (*PlannerNote, error) {
	path := fmt.Sprintf("/api/v1/planner_notes/%d", noteID)

	var note PlannerNote
	if err := s.client.GetJSON(ctx, path, &note); err != nil {
		return nil, err
	}

	return &note, nil
}

// CreateNoteParams holds parameters for creating a planner note
type CreateNoteParams struct {
	Title            string
	Details          string
	TodoDate         string // ISO 8601
	CourseID         int64
	LinkedObjectType string
	LinkedObjectID   int64
}

// CreateNote creates a new planner note
func (s *PlannerService) CreateNote(ctx context.Context, params *CreateNoteParams) (*PlannerNote, error) {
	path := "/api/v1/planner_notes"

	body := make(map[string]interface{})

	body["title"] = params.Title

	if params.Details != "" {
		body["details"] = params.Details
	}

	if params.TodoDate != "" {
		body["todo_date"] = params.TodoDate
	}

	if params.CourseID > 0 {
		body["course_id"] = params.CourseID
	}

	if params.LinkedObjectType != "" {
		body["linked_object_type"] = params.LinkedObjectType
	}

	if params.LinkedObjectID > 0 {
		body["linked_object_id"] = params.LinkedObjectID
	}

	var note PlannerNote
	if err := s.client.PostJSON(ctx, path, body, &note); err != nil {
		return nil, err
	}

	return &note, nil
}

// UpdateNoteParams holds parameters for updating a planner note
type UpdateNoteParams struct {
	Title    *string
	Details  *string
	TodoDate *string
	CourseID *int64
}

// UpdateNote updates an existing planner note
func (s *PlannerService) UpdateNote(ctx context.Context, noteID int64, params *UpdateNoteParams) (*PlannerNote, error) {
	path := fmt.Sprintf("/api/v1/planner_notes/%d", noteID)

	body := make(map[string]interface{})

	if params.Title != nil {
		body["title"] = *params.Title
	}

	if params.Details != nil {
		body["details"] = *params.Details
	}

	if params.TodoDate != nil {
		body["todo_date"] = *params.TodoDate
	}

	if params.CourseID != nil {
		body["course_id"] = *params.CourseID
	}

	var note PlannerNote
	if err := s.client.PutJSON(ctx, path, body, &note); err != nil {
		return nil, err
	}

	return &note, nil
}

// DeleteNote deletes a planner note
func (s *PlannerService) DeleteNote(ctx context.Context, noteID int64) error {
	path := fmt.Sprintf("/api/v1/planner_notes/%d", noteID)
	_, err := s.client.Delete(ctx, path)
	return err
}

// ListOverridesOptions holds options for listing planner overrides
type ListOverridesOptions struct {
	PlannableType string
	PlannableID   int64
}

// ListOverrides retrieves planner overrides for the current user
func (s *PlannerService) ListOverrides(ctx context.Context, opts *ListOverridesOptions) ([]PlannerOverride, error) {
	path := "/api/v1/planner/overrides"

	if opts != nil {
		query := url.Values{}

		if opts.PlannableType != "" {
			query.Add("plannable_type", opts.PlannableType)
		}

		if opts.PlannableID > 0 {
			query.Add("plannable_id", strconv.FormatInt(opts.PlannableID, 10))
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var overrides []PlannerOverride
	if err := s.client.GetAllPages(ctx, path, &overrides); err != nil {
		return nil, err
	}

	return overrides, nil
}

// GetOverride retrieves a single planner override
func (s *PlannerService) GetOverride(ctx context.Context, overrideID int64) (*PlannerOverride, error) {
	path := fmt.Sprintf("/api/v1/planner/overrides/%d", overrideID)

	var override PlannerOverride
	if err := s.client.GetJSON(ctx, path, &override); err != nil {
		return nil, err
	}

	return &override, nil
}

// CreateOverrideParams holds parameters for creating a planner override
type CreateOverrideParams struct {
	PlannableType  string // Required: Assignment, Quiz, CalendarEvent, etc.
	PlannableID    int64  // Required
	MarkedComplete bool
	Dismissed      bool
}

// CreateOverride creates a new planner override
func (s *PlannerService) CreateOverride(ctx context.Context, params *CreateOverrideParams) (*PlannerOverride, error) {
	path := "/api/v1/planner/overrides"

	body := map[string]interface{}{
		"plannable_type": params.PlannableType,
		"plannable_id":   params.PlannableID,
	}

	if params.MarkedComplete {
		body["marked_complete"] = true
	}

	if params.Dismissed {
		body["dismissed"] = true
	}

	var override PlannerOverride
	if err := s.client.PostJSON(ctx, path, body, &override); err != nil {
		return nil, err
	}

	return &override, nil
}

// UpdateOverrideParams holds parameters for updating a planner override
type UpdateOverrideParams struct {
	MarkedComplete *bool
	Dismissed      *bool
}

// UpdateOverride updates an existing planner override
func (s *PlannerService) UpdateOverride(ctx context.Context, overrideID int64, params *UpdateOverrideParams) (*PlannerOverride, error) {
	path := fmt.Sprintf("/api/v1/planner/overrides/%d", overrideID)

	body := make(map[string]interface{})

	if params.MarkedComplete != nil {
		body["marked_complete"] = *params.MarkedComplete
	}

	if params.Dismissed != nil {
		body["dismissed"] = *params.Dismissed
	}

	var override PlannerOverride
	if err := s.client.PutJSON(ctx, path, body, &override); err != nil {
		return nil, err
	}

	return &override, nil
}

// DeleteOverride deletes a planner override
func (s *PlannerService) DeleteOverride(ctx context.Context, overrideID int64) error {
	path := fmt.Sprintf("/api/v1/planner/overrides/%d", overrideID)
	_, err := s.client.Delete(ctx, path)
	return err
}
