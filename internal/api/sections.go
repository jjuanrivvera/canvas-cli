package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// SectionsService handles section-related API calls
type SectionsService struct {
	client *Client
}

// NewSectionsService creates a new sections service
func NewSectionsService(client *Client) *SectionsService {
	return &SectionsService{client: client}
}

// Section represents a Canvas course section
type Section struct {
	ID                                int64      `json:"id"`
	Name                              string     `json:"name"`
	SISSectionID                      string     `json:"sis_section_id,omitempty"`
	IntegrationID                     string     `json:"integration_id,omitempty"`
	SISImportID                       int64      `json:"sis_import_id,omitempty"`
	CourseID                          int64      `json:"course_id"`
	SISCourseID                       string     `json:"sis_course_id,omitempty"`
	StartAt                           *time.Time `json:"start_at,omitempty"`
	EndAt                             *time.Time `json:"end_at,omitempty"`
	RestrictEnrollmentsToSectionDates bool       `json:"restrict_enrollments_to_section_dates"`
	NonXlistCourseID                  *int64     `json:"nonxlist_course_id,omitempty"`
	TotalStudents                     int        `json:"total_students,omitempty"`
	CreatedAt                         *time.Time `json:"created_at,omitempty"`
}

// ListSectionsOptions holds options for listing sections
type ListSectionsOptions struct {
	Include []string // students, total_students, passback_status, permissions
	Page    int
	PerPage int
}

// ListCourse retrieves sections for a course
func (s *SectionsService) ListCourse(ctx context.Context, courseID int64, opts *ListSectionsOptions) ([]Section, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/sections", courseID)

	if opts != nil {
		query := url.Values{}

		for _, inc := range opts.Include {
			query.Add("include[]", inc)
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

	var sections []Section
	if err := s.client.GetAllPages(ctx, path, &sections); err != nil {
		return nil, err
	}

	return sections, nil
}

// Get retrieves a single section by ID
func (s *SectionsService) Get(ctx context.Context, sectionID int64, include []string) (*Section, error) {
	path := fmt.Sprintf("/api/v1/sections/%d", sectionID)

	if len(include) > 0 {
		query := url.Values{}
		for _, inc := range include {
			query.Add("include[]", inc)
		}
		path += "?" + query.Encode()
	}

	var section Section
	if err := s.client.GetJSON(ctx, path, &section); err != nil {
		return nil, err
	}

	return &section, nil
}

// CreateSectionParams holds parameters for creating a section
type CreateSectionParams struct {
	Name                              string
	SISSectionID                      string
	IntegrationID                     string
	StartAt                           string
	EndAt                             string
	RestrictEnrollmentsToSectionDates bool
	EnableSISReactivation             bool
}

// Create creates a new section in a course
func (s *SectionsService) Create(ctx context.Context, courseID int64, params *CreateSectionParams) (*Section, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/sections", courseID)

	body := map[string]interface{}{
		"course_section": make(map[string]interface{}),
	}

	sectionData, ok := body["course_section"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("internal error: invalid section data structure")
	}

	if params.Name != "" {
		sectionData["name"] = params.Name
	}

	if params.SISSectionID != "" {
		sectionData["sis_section_id"] = params.SISSectionID
	}

	if params.IntegrationID != "" {
		sectionData["integration_id"] = params.IntegrationID
	}

	if params.StartAt != "" {
		sectionData["start_at"] = params.StartAt
	}

	if params.EndAt != "" {
		sectionData["end_at"] = params.EndAt
	}

	if params.RestrictEnrollmentsToSectionDates {
		sectionData["restrict_enrollments_to_section_dates"] = true
	}

	if params.EnableSISReactivation {
		sectionData["enable_sis_reactivation"] = true
	}

	var section Section
	if err := s.client.PostJSON(ctx, path, body, &section); err != nil {
		return nil, err
	}

	return &section, nil
}

// UpdateSectionParams holds parameters for updating a section
type UpdateSectionParams struct {
	Name                              *string
	SISSectionID                      *string
	IntegrationID                     *string
	StartAt                           *string
	EndAt                             *string
	RestrictEnrollmentsToSectionDates *bool
	OverrideSISStickiness             bool
}

// Update updates an existing section
func (s *SectionsService) Update(ctx context.Context, sectionID int64, params *UpdateSectionParams) (*Section, error) {
	path := fmt.Sprintf("/api/v1/sections/%d", sectionID)

	body := map[string]interface{}{
		"course_section": make(map[string]interface{}),
	}

	sectionData, ok := body["course_section"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("internal error: invalid section data structure")
	}

	if params.Name != nil {
		sectionData["name"] = *params.Name
	}

	if params.SISSectionID != nil {
		sectionData["sis_section_id"] = *params.SISSectionID
	}

	if params.IntegrationID != nil {
		sectionData["integration_id"] = *params.IntegrationID
	}

	if params.StartAt != nil {
		sectionData["start_at"] = *params.StartAt
	}

	if params.EndAt != nil {
		sectionData["end_at"] = *params.EndAt
	}

	if params.RestrictEnrollmentsToSectionDates != nil {
		sectionData["restrict_enrollments_to_section_dates"] = *params.RestrictEnrollmentsToSectionDates
	}

	if params.OverrideSISStickiness {
		body["override_sis_stickiness"] = true
	}

	var section Section
	if err := s.client.PutJSON(ctx, path, body, &section); err != nil {
		return nil, err
	}

	return &section, nil
}

// Delete deletes a section
func (s *SectionsService) Delete(ctx context.Context, sectionID int64) (*Section, error) {
	path := fmt.Sprintf("/api/v1/sections/%d", sectionID)

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var section Section
	if err := json.NewDecoder(resp.Body).Decode(&section); err != nil {
		return nil, err
	}

	return &section, nil
}

// Crosslist moves a section to a different course
func (s *SectionsService) Crosslist(ctx context.Context, sectionID, newCourseID int64, overrideSISStickiness bool) (*Section, error) {
	path := fmt.Sprintf("/api/v1/sections/%d/crosslist/%d", sectionID, newCourseID)

	if overrideSISStickiness {
		path += "?override_sis_stickiness=true"
	}

	var section Section
	if err := s.client.PostJSON(ctx, path, nil, &section); err != nil {
		return nil, err
	}

	return &section, nil
}

// Uncrosslist returns a crosslisted section to its original course
func (s *SectionsService) Uncrosslist(ctx context.Context, sectionID int64, overrideSISStickiness bool) (*Section, error) {
	path := fmt.Sprintf("/api/v1/sections/%d/crosslist", sectionID)

	if overrideSISStickiness {
		path += "?override_sis_stickiness=true"
	}

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var section Section
	if err := json.NewDecoder(resp.Body).Decode(&section); err != nil {
		return nil, err
	}

	return &section, nil
}
