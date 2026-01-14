package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// OutcomesService handles outcome-related API calls
type OutcomesService struct {
	client *Client
}

// NewOutcomesService creates a new outcomes service
func NewOutcomesService(client *Client) *OutcomesService {
	return &OutcomesService{client: client}
}

// Outcome represents a Canvas learning outcome
type Outcome struct {
	ID                   int64           `json:"id"`
	URL                  string          `json:"url,omitempty"`
	ContextID            int64           `json:"context_id,omitempty"`
	ContextType          string          `json:"context_type,omitempty"`
	Title                string          `json:"title"`
	DisplayName          string          `json:"display_name,omitempty"`
	Description          string          `json:"description,omitempty"`
	VendorGUID           string          `json:"vendor_guid,omitempty"`
	PointsPossible       float64         `json:"points_possible,omitempty"`
	MasteryPoints        float64         `json:"mastery_points,omitempty"`
	CalculationMethod    string          `json:"calculation_method,omitempty"`
	CalculationInt       int             `json:"calculation_int,omitempty"`
	Ratings              []OutcomeRating `json:"ratings,omitempty"`
	CanEdit              bool            `json:"can_edit,omitempty"`
	CanUnlink            bool            `json:"can_unlink,omitempty"`
	Assessed             bool            `json:"assessed,omitempty"`
	HasUpdateableRubrics bool            `json:"has_updateable_rubrics,omitempty"`
}

// OutcomeRating represents a rating level for an outcome
type OutcomeRating struct {
	Description string  `json:"description"`
	Points      float64 `json:"points"`
}

// OutcomeGroup represents an outcome group
type OutcomeGroup struct {
	ID                 int64            `json:"id"`
	URL                string           `json:"url,omitempty"`
	ParentOutcomeGroup *OutcomeGroupRef `json:"parent_outcome_group,omitempty"`
	ContextID          int64            `json:"context_id,omitempty"`
	ContextType        string           `json:"context_type,omitempty"`
	Title              string           `json:"title"`
	Description        string           `json:"description,omitempty"`
	VendorGUID         string           `json:"vendor_guid,omitempty"`
	SubgroupsURL       string           `json:"subgroups_url,omitempty"`
	OutcomesURL        string           `json:"outcomes_url,omitempty"`
	ImportURL          string           `json:"import_url,omitempty"`
	CanEdit            bool             `json:"can_edit,omitempty"`
}

// OutcomeGroupRef represents a reference to a parent outcome group
type OutcomeGroupRef struct {
	ID   int64  `json:"id"`
	Type string `json:"type,omitempty"`
}

// OutcomeLink represents a link between an outcome and a group
type OutcomeLink struct {
	URL          string        `json:"url,omitempty"`
	ContextID    int64         `json:"context_id,omitempty"`
	ContextType  string        `json:"context_type,omitempty"`
	OutcomeGroup *OutcomeGroup `json:"outcome_group,omitempty"`
	Outcome      *Outcome      `json:"outcome,omitempty"`
	Assessed     bool          `json:"assessed,omitempty"`
	CanUnlink    bool          `json:"can_unlink,omitempty"`
}

// OutcomeResult represents a student's result for an outcome
type OutcomeResult struct {
	ID            int64               `json:"id"`
	Score         float64             `json:"score,omitempty"`
	Submitted     bool                `json:"submitted,omitempty"`
	Possible      float64             `json:"possible,omitempty"`
	Mastery       bool                `json:"mastery,omitempty"`
	Percent       float64             `json:"percent,omitempty"`
	HidePoints    bool                `json:"hide_points,omitempty"`
	HiddenOutcome bool                `json:"hidden_outcome,omitempty"`
	Links         *OutcomeResultLinks `json:"links,omitempty"`
}

// OutcomeResultLinks contains references for an outcome result
type OutcomeResultLinks struct {
	User            string `json:"user,omitempty"`
	LearningOutcome string `json:"learning_outcome,omitempty"`
	Alignment       string `json:"alignment,omitempty"`
}

// OutcomeAlignment represents an alignment between an outcome and content
type OutcomeAlignment struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name"`
	HTMLURL string `json:"html_url,omitempty"`
}

// Get retrieves a single outcome by ID
func (s *OutcomesService) Get(ctx context.Context, outcomeID int64) (*Outcome, error) {
	path := fmt.Sprintf("/api/v1/outcomes/%d", outcomeID)

	var outcome Outcome
	if err := s.client.GetJSON(ctx, path, &outcome); err != nil {
		return nil, err
	}

	return &outcome, nil
}

// UpdateOutcomeParams holds parameters for updating an outcome
type UpdateOutcomeParams struct {
	Title             *string
	DisplayName       *string
	Description       *string
	VendorGUID        *string
	MasteryPoints     *float64
	Ratings           []OutcomeRating
	CalculationMethod *string
	CalculationInt    *int
}

// Update updates an existing outcome
func (s *OutcomesService) Update(ctx context.Context, outcomeID int64, params *UpdateOutcomeParams) (*Outcome, error) {
	path := fmt.Sprintf("/api/v1/outcomes/%d", outcomeID)

	body := make(map[string]interface{})

	if params.Title != nil {
		body["title"] = *params.Title
	}

	if params.DisplayName != nil {
		body["display_name"] = *params.DisplayName
	}

	if params.Description != nil {
		body["description"] = *params.Description
	}

	if params.VendorGUID != nil {
		body["vendor_guid"] = *params.VendorGUID
	}

	if params.MasteryPoints != nil {
		body["mastery_points"] = *params.MasteryPoints
	}

	if len(params.Ratings) > 0 {
		body["ratings"] = params.Ratings
	}

	if params.CalculationMethod != nil {
		body["calculation_method"] = *params.CalculationMethod
	}

	if params.CalculationInt != nil {
		body["calculation_int"] = *params.CalculationInt
	}

	var outcome Outcome
	if err := s.client.PutJSON(ctx, path, body, &outcome); err != nil {
		return nil, err
	}

	return &outcome, nil
}

// ListOutcomeGroupsOptions holds options for listing outcome groups
type ListOutcomeGroupsOptions struct {
	Page    int
	PerPage int
}

// ListGroupsAccount retrieves outcome groups for an account
func (s *OutcomesService) ListGroupsAccount(ctx context.Context, accountID int64, opts *ListOutcomeGroupsOptions) ([]OutcomeGroup, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/outcome_groups", accountID)

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

	var groups []OutcomeGroup
	if err := s.client.GetAllPages(ctx, path, &groups); err != nil {
		return nil, err
	}

	return groups, nil
}

// ListGroupsCourse retrieves outcome groups for a course
func (s *OutcomesService) ListGroupsCourse(ctx context.Context, courseID int64, opts *ListOutcomeGroupsOptions) ([]OutcomeGroup, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/outcome_groups", courseID)

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

	var groups []OutcomeGroup
	if err := s.client.GetAllPages(ctx, path, &groups); err != nil {
		return nil, err
	}

	return groups, nil
}

// GetGroupAccount retrieves a single outcome group from an account
func (s *OutcomesService) GetGroupAccount(ctx context.Context, accountID, groupID int64) (*OutcomeGroup, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/outcome_groups/%d", accountID, groupID)

	var group OutcomeGroup
	if err := s.client.GetJSON(ctx, path, &group); err != nil {
		return nil, err
	}

	return &group, nil
}

// GetGroupCourse retrieves a single outcome group from a course
func (s *OutcomesService) GetGroupCourse(ctx context.Context, courseID, groupID int64) (*OutcomeGroup, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/outcome_groups/%d", courseID, groupID)

	var group OutcomeGroup
	if err := s.client.GetJSON(ctx, path, &group); err != nil {
		return nil, err
	}

	return &group, nil
}

// ListOutcomesInGroupOptions holds options for listing outcomes in a group
type ListOutcomesInGroupOptions struct {
	Page    int
	PerPage int
}

// ListOutcomesInGroupAccount retrieves outcomes in a group from an account
func (s *OutcomesService) ListOutcomesInGroupAccount(ctx context.Context, accountID, groupID int64, opts *ListOutcomesInGroupOptions) ([]OutcomeLink, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/outcome_groups/%d/outcomes", accountID, groupID)

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

	var links []OutcomeLink
	if err := s.client.GetAllPages(ctx, path, &links); err != nil {
		return nil, err
	}

	return links, nil
}

// ListOutcomesInGroupCourse retrieves outcomes in a group from a course
func (s *OutcomesService) ListOutcomesInGroupCourse(ctx context.Context, courseID, groupID int64, opts *ListOutcomesInGroupOptions) ([]OutcomeLink, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/outcome_groups/%d/outcomes", courseID, groupID)

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

	var links []OutcomeLink
	if err := s.client.GetAllPages(ctx, path, &links); err != nil {
		return nil, err
	}

	return links, nil
}

// CreateOutcomeParams holds parameters for creating a new outcome
type CreateOutcomeParams struct {
	Title             string
	DisplayName       string
	Description       string
	VendorGUID        string
	MasteryPoints     float64
	Ratings           []OutcomeRating
	CalculationMethod string
	CalculationInt    int
}

// CreateOutcomeAccount creates a new outcome in an account group
func (s *OutcomesService) CreateOutcomeAccount(ctx context.Context, accountID, groupID int64, params *CreateOutcomeParams) (*OutcomeLink, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/outcome_groups/%d/outcomes", accountID, groupID)

	body := map[string]interface{}{
		"title": params.Title,
	}

	if params.DisplayName != "" {
		body["display_name"] = params.DisplayName
	}

	if params.Description != "" {
		body["description"] = params.Description
	}

	if params.VendorGUID != "" {
		body["vendor_guid"] = params.VendorGUID
	}

	if params.MasteryPoints > 0 {
		body["mastery_points"] = params.MasteryPoints
	}

	if len(params.Ratings) > 0 {
		body["ratings"] = params.Ratings
	}

	if params.CalculationMethod != "" {
		body["calculation_method"] = params.CalculationMethod
	}

	if params.CalculationInt > 0 {
		body["calculation_int"] = params.CalculationInt
	}

	var link OutcomeLink
	if err := s.client.PostJSON(ctx, path, body, &link); err != nil {
		return nil, err
	}

	return &link, nil
}

// CreateOutcomeCourse creates a new outcome in a course group
func (s *OutcomesService) CreateOutcomeCourse(ctx context.Context, courseID, groupID int64, params *CreateOutcomeParams) (*OutcomeLink, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/outcome_groups/%d/outcomes", courseID, groupID)

	body := map[string]interface{}{
		"title": params.Title,
	}

	if params.DisplayName != "" {
		body["display_name"] = params.DisplayName
	}

	if params.Description != "" {
		body["description"] = params.Description
	}

	if params.VendorGUID != "" {
		body["vendor_guid"] = params.VendorGUID
	}

	if params.MasteryPoints > 0 {
		body["mastery_points"] = params.MasteryPoints
	}

	if len(params.Ratings) > 0 {
		body["ratings"] = params.Ratings
	}

	if params.CalculationMethod != "" {
		body["calculation_method"] = params.CalculationMethod
	}

	if params.CalculationInt > 0 {
		body["calculation_int"] = params.CalculationInt
	}

	var link OutcomeLink
	if err := s.client.PostJSON(ctx, path, body, &link); err != nil {
		return nil, err
	}

	return &link, nil
}

// LinkOutcomeAccount links an outcome to a group in an account
func (s *OutcomesService) LinkOutcomeAccount(ctx context.Context, accountID, groupID, outcomeID int64) (*OutcomeLink, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/outcome_groups/%d/outcomes/%d", accountID, groupID, outcomeID)

	var link OutcomeLink
	if err := s.client.PutJSON(ctx, path, nil, &link); err != nil {
		return nil, err
	}

	return &link, nil
}

// LinkOutcomeCourse links an outcome to a group in a course
func (s *OutcomesService) LinkOutcomeCourse(ctx context.Context, courseID, groupID, outcomeID int64) (*OutcomeLink, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/outcome_groups/%d/outcomes/%d", courseID, groupID, outcomeID)

	var link OutcomeLink
	if err := s.client.PutJSON(ctx, path, nil, &link); err != nil {
		return nil, err
	}

	return &link, nil
}

// UnlinkOutcomeAccount removes an outcome link from a group in an account
func (s *OutcomesService) UnlinkOutcomeAccount(ctx context.Context, accountID, groupID, outcomeID int64) (*OutcomeLink, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/outcome_groups/%d/outcomes/%d", accountID, groupID, outcomeID)

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// The unlink endpoint returns the deleted link
	var link OutcomeLink
	if err := json.NewDecoder(resp.Body).Decode(&link); err != nil {
		return nil, err
	}

	return &link, nil
}

// UnlinkOutcomeCourse removes an outcome link from a group in a course
func (s *OutcomesService) UnlinkOutcomeCourse(ctx context.Context, courseID, groupID, outcomeID int64) (*OutcomeLink, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/outcome_groups/%d/outcomes/%d", courseID, groupID, outcomeID)

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var link OutcomeLink
	if err := json.NewDecoder(resp.Body).Decode(&link); err != nil {
		return nil, err
	}

	return &link, nil
}

// OutcomeResultsOptions holds options for fetching outcome results
type OutcomeResultsOptions struct {
	UserIDs       []int64
	OutcomeIDs    []int64
	Include       []string // alignments, outcomes, outcomes.alignments, outcome_groups, outcome_links, outcome_paths, users
	IncludeHidden bool
	Page          int
	PerPage       int
}

// OutcomeResultsResponse wraps the outcome results response
type OutcomeResultsResponse struct {
	OutcomeResults []OutcomeResult       `json:"outcome_results"`
	Linked         *OutcomeResultsLinked `json:"linked,omitempty"`
}

// OutcomeResultsLinked represents linked resources in outcome results
type OutcomeResultsLinked struct {
	Users    []User    `json:"users,omitempty"`
	Outcomes []Outcome `json:"outcomes,omitempty"`
}

// GetResults retrieves outcome results for a course
func (s *OutcomesService) GetResults(ctx context.Context, courseID int64, opts *OutcomeResultsOptions) (*OutcomeResultsResponse, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/outcome_results", courseID)

	if opts != nil {
		query := url.Values{}

		for _, userID := range opts.UserIDs {
			query.Add("user_ids[]", strconv.FormatInt(userID, 10))
		}

		for _, outcomeID := range opts.OutcomeIDs {
			query.Add("outcome_ids[]", strconv.FormatInt(outcomeID, 10))
		}

		for _, include := range opts.Include {
			query.Add("include[]", include)
		}

		if opts.IncludeHidden {
			query.Add("include_hidden", "true")
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

	var response OutcomeResultsResponse
	if err := s.client.GetJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetAlignments retrieves outcome alignments for a course
func (s *OutcomesService) GetAlignments(ctx context.Context, courseID int64, studentID int64) ([]OutcomeAlignment, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/outcome_alignments", courseID)

	if studentID > 0 {
		path += fmt.Sprintf("?student_id=%d", studentID)
	}

	var alignments []OutcomeAlignment
	if err := s.client.GetJSON(ctx, path, &alignments); err != nil {
		return nil, err
	}

	return alignments, nil
}
