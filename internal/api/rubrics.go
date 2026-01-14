package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// RubricsService handles rubric-related API calls
type RubricsService struct {
	client *Client
}

// NewRubricsService creates a new rubrics service
func NewRubricsService(client *Client) *RubricsService {
	return &RubricsService{client: client}
}

// Rubric represents a Canvas rubric
type Rubric struct {
	ID                        int64               `json:"id"`
	Title                     string              `json:"title"`
	ContextType               string              `json:"context_type"`
	ContextID                 int64               `json:"context_id"`
	PointsPossible            float64             `json:"points_possible"`
	Reusable                  bool                `json:"reusable"`
	ReadOnly                  bool                `json:"read_only"`
	FreeFormCriterionComments bool                `json:"free_form_criterion_comments"`
	HideScoreTotal            bool                `json:"hide_score_total"`
	Data                      []RubricCriterion   `json:"data,omitempty"`
	Assessments               []RubricAssessment  `json:"assessments,omitempty"`
	Associations              []RubricAssociation `json:"associations,omitempty"`
}

// RubricAssociation represents a rubric association with an assignment
type RubricAssociation struct {
	ID                 int64  `json:"id"`
	RubricID           int64  `json:"rubric_id"`
	AssociationID      int64  `json:"association_id"`
	AssociationType    string `json:"association_type"`
	UseForGrading      bool   `json:"use_for_grading"`
	SummaryData        string `json:"summary_data,omitempty"`
	Purpose            string `json:"purpose"`
	HideScoreTotal     bool   `json:"hide_score_total"`
	HidePoints         bool   `json:"hide_points"`
	HideOutcomeResults bool   `json:"hide_outcome_results"`
}

// rubricResponse is a wrapper for API responses that include a rubric
type rubricResponse struct {
	Rubric *Rubric `json:"rubric"`
}

// rubricAssociationResponse is a wrapper for API responses that include a rubric association
type rubricAssociationResponse struct {
	RubricAssociation *RubricAssociation `json:"rubric_association"`
}

// ListRubricsOptions holds options for listing rubrics
type ListRubricsOptions struct {
	Include []string // assessments, associations, assignment_associations
	Page    int
	PerPage int
}

// ListCourse retrieves all rubrics for a course
func (s *RubricsService) ListCourse(ctx context.Context, courseID int64, opts *ListRubricsOptions) ([]Rubric, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/rubrics", courseID)

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

	var rubrics []Rubric
	if err := s.client.GetAllPages(ctx, path, &rubrics); err != nil {
		return nil, err
	}

	return rubrics, nil
}

// ListAccount retrieves all rubrics for an account
func (s *RubricsService) ListAccount(ctx context.Context, accountID int64, opts *ListRubricsOptions) ([]Rubric, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/rubrics", accountID)

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

	var rubrics []Rubric
	if err := s.client.GetAllPages(ctx, path, &rubrics); err != nil {
		return nil, err
	}

	return rubrics, nil
}

// GetCourse retrieves a single rubric from a course
func (s *RubricsService) GetCourse(ctx context.Context, courseID, rubricID int64, include []string) (*Rubric, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/rubrics/%d", courseID, rubricID)

	if len(include) > 0 {
		query := url.Values{}
		for _, inc := range include {
			query.Add("include[]", inc)
		}
		path += "?" + query.Encode()
	}

	var rubric Rubric
	if err := s.client.GetJSON(ctx, path, &rubric); err != nil {
		return nil, err
	}

	return &rubric, nil
}

// GetAccount retrieves a single rubric from an account
func (s *RubricsService) GetAccount(ctx context.Context, accountID, rubricID int64, include []string) (*Rubric, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/rubrics/%d", accountID, rubricID)

	if len(include) > 0 {
		query := url.Values{}
		for _, inc := range include {
			query.Add("include[]", inc)
		}
		path += "?" + query.Encode()
	}

	var rubric Rubric
	if err := s.client.GetJSON(ctx, path, &rubric); err != nil {
		return nil, err
	}

	return &rubric, nil
}

// CreateRubricParams holds parameters for creating a rubric
type CreateRubricParams struct {
	Title                     string
	PointsPossible            float64
	FreeFormCriterionComments bool
	HideScoreTotal            bool
	Criteria                  []RubricCriterion
}

// Create creates a new rubric in a course
func (s *RubricsService) Create(ctx context.Context, courseID int64, params *CreateRubricParams) (*Rubric, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/rubrics", courseID)

	body := map[string]interface{}{
		"rubric": make(map[string]interface{}),
	}

	rubricData, ok := body["rubric"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("internal error: invalid rubric data structure")
	}

	if params.Title != "" {
		rubricData["title"] = params.Title
	}

	if params.PointsPossible > 0 {
		rubricData["points_possible"] = params.PointsPossible
	}

	if params.FreeFormCriterionComments {
		rubricData["free_form_criterion_comments"] = params.FreeFormCriterionComments
	}

	if params.HideScoreTotal {
		rubricData["hide_score_total"] = params.HideScoreTotal
	}

	if len(params.Criteria) > 0 {
		criteria := make(map[string]interface{})
		for i, c := range params.Criteria {
			key := strconv.Itoa(i)
			criterionData := map[string]interface{}{
				"description":      c.Description,
				"long_description": c.LongDescription,
				"points":           c.Points,
			}

			if len(c.Ratings) > 0 {
				ratings := make(map[string]interface{})
				for j, r := range c.Ratings {
					ratingKey := strconv.Itoa(j)
					ratings[ratingKey] = map[string]interface{}{
						"description":      r.Description,
						"long_description": r.LongDescription,
						"points":           r.Points,
					}
				}
				criterionData["ratings"] = ratings
			}

			criteria[key] = criterionData
		}
		body["rubric"].(map[string]interface{})["criteria"] = criteria
	}

	var response rubricResponse
	if err := s.client.PostJSON(ctx, path, body, &response); err != nil {
		return nil, err
	}

	if response.Rubric == nil {
		return nil, fmt.Errorf("rubric not returned in response")
	}

	return response.Rubric, nil
}

// UpdateRubricParams holds parameters for updating a rubric
type UpdateRubricParams struct {
	Title                     *string
	PointsPossible            *float64
	FreeFormCriterionComments *bool
	HideScoreTotal            *bool
}

// Update updates an existing rubric
func (s *RubricsService) Update(ctx context.Context, courseID, rubricID int64, params *UpdateRubricParams) (*Rubric, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/rubrics/%d", courseID, rubricID)

	body := map[string]interface{}{
		"rubric": make(map[string]interface{}),
	}

	rubricData, ok := body["rubric"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("internal error: invalid rubric data structure")
	}

	if params.Title != nil {
		rubricData["title"] = *params.Title
	}

	if params.PointsPossible != nil {
		rubricData["points_possible"] = *params.PointsPossible
	}

	if params.FreeFormCriterionComments != nil {
		rubricData["free_form_criterion_comments"] = *params.FreeFormCriterionComments
	}

	if params.HideScoreTotal != nil {
		rubricData["hide_score_total"] = *params.HideScoreTotal
	}

	var response rubricResponse
	if err := s.client.PutJSON(ctx, path, body, &response); err != nil {
		return nil, err
	}

	if response.Rubric == nil {
		return nil, fmt.Errorf("rubric not returned in response")
	}

	return response.Rubric, nil
}

// Delete deletes a rubric
func (s *RubricsService) Delete(ctx context.Context, courseID, rubricID int64) (*Rubric, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/rubrics/%d", courseID, rubricID)

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response rubricResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if response.Rubric == nil {
		return nil, fmt.Errorf("rubric not returned in response")
	}

	return response.Rubric, nil
}

// AssociateParams holds parameters for associating a rubric
type AssociateParams struct {
	AssociationType string // "Assignment"
	AssociationID   int64
	UseForGrading   bool
	HideScoreTotal  bool
	HidePoints      bool
	Purpose         string // "grading", "bookmark"
}

// Associate associates a rubric with an assignment
func (s *RubricsService) Associate(ctx context.Context, courseID, rubricID int64, params *AssociateParams) (*RubricAssociation, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/rubric_associations", courseID)

	body := map[string]interface{}{
		"rubric_association": map[string]interface{}{
			"rubric_id":        rubricID,
			"association_type": params.AssociationType,
			"association_id":   params.AssociationID,
			"use_for_grading":  params.UseForGrading,
			"hide_score_total": params.HideScoreTotal,
			"hide_points":      params.HidePoints,
			"purpose":          params.Purpose,
		},
	}

	var response rubricAssociationResponse
	if err := s.client.PostJSON(ctx, path, body, &response); err != nil {
		return nil, err
	}

	if response.RubricAssociation == nil {
		return nil, fmt.Errorf("rubric association not returned in response")
	}

	return response.RubricAssociation, nil
}
