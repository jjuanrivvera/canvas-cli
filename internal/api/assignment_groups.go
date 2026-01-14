package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// AssignmentGroupsService handles assignment group-related API calls
type AssignmentGroupsService struct {
	client *Client
}

// NewAssignmentGroupsService creates a new assignment groups service
func NewAssignmentGroupsService(client *Client) *AssignmentGroupsService {
	return &AssignmentGroupsService{client: client}
}

// AssignmentGroup represents a Canvas assignment group
type AssignmentGroup struct {
	ID              int64         `json:"id"`
	Name            string        `json:"name"`
	Position        int           `json:"position"`
	GroupWeight     float64       `json:"group_weight"`
	SISSourceID     string        `json:"sis_source_id,omitempty"`
	IntegrationData interface{}   `json:"integration_data,omitempty"`
	Assignments     []Assignment  `json:"assignments,omitempty"`
	Rules           *GradingRules `json:"rules,omitempty"`
}

// GradingRules represents rules for an assignment group
type GradingRules struct {
	DropLowest  int     `json:"drop_lowest,omitempty"`
	DropHighest int     `json:"drop_highest,omitempty"`
	NeverDrop   []int64 `json:"never_drop,omitempty"`
}

// ListAssignmentGroupsOptions holds options for listing assignment groups
type ListAssignmentGroupsOptions struct {
	Include                          []string // assignments, discussion_topic, all_dates, assignment_visibility, overrides, submission, observed_users, can_edit, score_statistics
	AssignmentIDs                    []int64
	ExcludeAssignmentSubmissionTypes []string
	OverrideAssignmentDates          *bool
	GradingPeriodID                  int64
	ScopeAssignmentsToStudent        *bool
	Page                             int
	PerPage                          int
}

// List retrieves assignment groups for a course
func (s *AssignmentGroupsService) List(ctx context.Context, courseID int64, opts *ListAssignmentGroupsOptions) ([]AssignmentGroup, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/assignment_groups", courseID)

	if opts != nil {
		query := url.Values{}

		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}

		for _, id := range opts.AssignmentIDs {
			query.Add("assignment_ids[]", strconv.FormatInt(id, 10))
		}

		for _, t := range opts.ExcludeAssignmentSubmissionTypes {
			query.Add("exclude_assignment_submission_types[]", t)
		}

		if opts.OverrideAssignmentDates != nil {
			query.Add("override_assignment_dates", strconv.FormatBool(*opts.OverrideAssignmentDates))
		}

		if opts.GradingPeriodID > 0 {
			query.Add("grading_period_id", strconv.FormatInt(opts.GradingPeriodID, 10))
		}

		if opts.ScopeAssignmentsToStudent != nil {
			query.Add("scope_assignments_to_student", strconv.FormatBool(*opts.ScopeAssignmentsToStudent))
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

	var groups []AssignmentGroup
	if err := s.client.GetAllPages(ctx, path, &groups); err != nil {
		return nil, err
	}

	return groups, nil
}

// Get retrieves a single assignment group
func (s *AssignmentGroupsService) Get(ctx context.Context, courseID, groupID int64, include []string) (*AssignmentGroup, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/assignment_groups/%d", courseID, groupID)

	if len(include) > 0 {
		query := url.Values{}
		for _, inc := range include {
			query.Add("include[]", inc)
		}
		path += "?" + query.Encode()
	}

	var group AssignmentGroup
	if err := s.client.GetJSON(ctx, path, &group); err != nil {
		return nil, err
	}

	return &group, nil
}

// CreateAssignmentGroupParams holds parameters for creating an assignment group
type CreateAssignmentGroupParams struct {
	Name            string
	Position        int
	GroupWeight     float64
	SISSourceID     string
	IntegrationData map[string]interface{}
	Rules           *GradingRules
}

// Create creates a new assignment group in a course
func (s *AssignmentGroupsService) Create(ctx context.Context, courseID int64, params *CreateAssignmentGroupParams) (*AssignmentGroup, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/assignment_groups", courseID)

	body := make(map[string]interface{})

	if params.Name != "" {
		body["name"] = params.Name
	}

	if params.Position > 0 {
		body["position"] = params.Position
	}

	if params.GroupWeight > 0 {
		body["group_weight"] = params.GroupWeight
	}

	if params.SISSourceID != "" {
		body["sis_source_id"] = params.SISSourceID
	}

	if params.IntegrationData != nil {
		body["integration_data"] = params.IntegrationData
	}

	if params.Rules != nil {
		rules := make(map[string]interface{})
		if params.Rules.DropLowest > 0 {
			rules["drop_lowest"] = params.Rules.DropLowest
		}
		if params.Rules.DropHighest > 0 {
			rules["drop_highest"] = params.Rules.DropHighest
		}
		if len(params.Rules.NeverDrop) > 0 {
			rules["never_drop"] = params.Rules.NeverDrop
		}
		if len(rules) > 0 {
			body["rules"] = rules
		}
	}

	var group AssignmentGroup
	if err := s.client.PostJSON(ctx, path, body, &group); err != nil {
		return nil, err
	}

	return &group, nil
}

// UpdateAssignmentGroupParams holds parameters for updating an assignment group
type UpdateAssignmentGroupParams struct {
	Name            *string
	Position        *int
	GroupWeight     *float64
	SISSourceID     *string
	IntegrationData map[string]interface{}
	Rules           *GradingRules
}

// Update updates an existing assignment group
func (s *AssignmentGroupsService) Update(ctx context.Context, courseID, groupID int64, params *UpdateAssignmentGroupParams) (*AssignmentGroup, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/assignment_groups/%d", courseID, groupID)

	body := make(map[string]interface{})

	if params.Name != nil {
		body["name"] = *params.Name
	}

	if params.Position != nil {
		body["position"] = *params.Position
	}

	if params.GroupWeight != nil {
		body["group_weight"] = *params.GroupWeight
	}

	if params.SISSourceID != nil {
		body["sis_source_id"] = *params.SISSourceID
	}

	if params.IntegrationData != nil {
		body["integration_data"] = params.IntegrationData
	}

	if params.Rules != nil {
		rules := make(map[string]interface{})
		if params.Rules.DropLowest > 0 {
			rules["drop_lowest"] = params.Rules.DropLowest
		}
		if params.Rules.DropHighest > 0 {
			rules["drop_highest"] = params.Rules.DropHighest
		}
		if len(params.Rules.NeverDrop) > 0 {
			rules["never_drop"] = params.Rules.NeverDrop
		}
		body["rules"] = rules
	}

	var group AssignmentGroup
	if err := s.client.PutJSON(ctx, path, body, &group); err != nil {
		return nil, err
	}

	return &group, nil
}

// DeleteOptions holds options for deleting an assignment group
type DeleteAssignmentGroupOptions struct {
	MoveAssignmentsTo int64 // Assignment group to move assignments to before deleting
}

// Delete deletes an assignment group
func (s *AssignmentGroupsService) Delete(ctx context.Context, courseID, groupID int64, opts *DeleteAssignmentGroupOptions) (*AssignmentGroup, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/assignment_groups/%d", courseID, groupID)

	if opts != nil && opts.MoveAssignmentsTo > 0 {
		query := url.Values{}
		query.Add("move_assignments_to", strconv.FormatInt(opts.MoveAssignmentsTo, 10))
		path += "?" + query.Encode()
	}

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var group AssignmentGroup
	if err := json.NewDecoder(resp.Body).Decode(&group); err != nil {
		return nil, err
	}

	return &group, nil
}
