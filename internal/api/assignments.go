package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// AssignmentsService handles assignment-related API calls
type AssignmentsService struct {
	client *Client
}

// NewAssignmentsService creates a new assignments service
func NewAssignmentsService(client *Client) *AssignmentsService {
	return &AssignmentsService{client: client}
}

// ListAssignmentsOptions holds options for listing assignments
type ListAssignmentsOptions struct {
	Include                    []string // Additional data to include (submission, assignment_visibility, overrides, observed_users, etc.)
	SearchTerm                 string   // Search by assignment name
	OverrideAssignmentDates    bool     // Apply assignment overrides for each assignment
	NeedsGradingCountBySection bool     // Include needs_grading_count split by section
	Bucket                     string   // Filter by assignment bucket (past, overdue, undated, ungraded, unsubmitted, upcoming, future)
	AssignmentIDs              []int64  // Filter by assignment IDs
	OrderBy                    string   // Order by (position, name, due_at)
	PostToSIS                  *bool    // Filter by post_to_sis
	Page                       int
	PerPage                    int
}

// Get retrieves a single assignment by ID
func (s *AssignmentsService) Get(ctx context.Context, courseID, assignmentID int64, include []string) (*Assignment, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/assignments/%d", courseID, assignmentID)

	if len(include) > 0 {
		query := url.Values{}
		for _, inc := range include {
			query.Add("include[]", inc)
		}
		path += "?" + query.Encode()
	}

	var assignment Assignment
	if err := s.client.GetJSON(ctx, path, &assignment); err != nil {
		return nil, err
	}

	return NormalizeAssignment(&assignment), nil
}

// List retrieves assignments for a course
func (s *AssignmentsService) List(ctx context.Context, courseID int64, opts *ListAssignmentsOptions) ([]Assignment, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/assignments", courseID)

	if opts != nil {
		query := url.Values{}

		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}

		if opts.SearchTerm != "" {
			query.Add("search_term", opts.SearchTerm)
		}

		if opts.OverrideAssignmentDates {
			query.Add("override_assignment_dates", "true")
		}

		if opts.NeedsGradingCountBySection {
			query.Add("needs_grading_count_by_section", "true")
		}

		if opts.Bucket != "" {
			query.Add("bucket", opts.Bucket)
		}

		for _, id := range opts.AssignmentIDs {
			query.Add("assignment_ids[]", strconv.FormatInt(id, 10))
		}

		if opts.OrderBy != "" {
			query.Add("order_by", opts.OrderBy)
		}

		if opts.PostToSIS != nil {
			query.Add("post_to_sis", strconv.FormatBool(*opts.PostToSIS))
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

	var assignments []Assignment
	if err := s.client.GetAllPages(ctx, path, &assignments); err != nil {
		return nil, err
	}

	return NormalizeAssignments(assignments), nil
}

// CreateAssignmentParams holds parameters for creating an assignment
type CreateAssignmentParams struct {
	Name                            string
	Position                        int
	SubmissionTypes                 []string // online_text_entry, online_url, online_upload, media_recording, etc.
	AllowedExtensions               []string // For online_upload submission type
	TurnitinEnabled                 bool
	VericiteEnabled                 bool
	TurnitinSettings                map[string]interface{}
	IntegrationData                 map[string]interface{}
	IntegrationID                   string
	PeerReviews                     bool
	AutomaticPeerReviews            bool
	NotifyOfUpdate                  bool
	GroupCategoryID                 int64
	GradeGroupStudentsIndividually  bool
	ExternalToolTagAttributes       map[string]interface{}
	PointsPossible                  float64
	GradingType                     string // pass_fail, percent, letter_grade, gpa_scale, points
	DueAt                           string // ISO8601 format
	LockAt                          string
	UnlockAt                        string
	Description                     string
	AssignmentGroupID               int64
	AssignmentOverrides             []AssignmentOverrideParams
	OnlyVisibleToOverrides          bool
	Published                       bool
	GradingStandardID               int64
	OmitFromFinalGrade              bool
	ModeratedGrading                bool
	GraderCount                     int
	FinalGraderID                   int64
	GraderCommentsVisibleToGraders  bool
	GradersAnonymousToGraders       bool
	GraderNamesVisibleToFinalGrader bool
	AnonymousInstructorAnnotations  bool
	AnonymousGrading                bool
	AllowedAttempts                 int
	AnnotatableAttachmentID         int64
	HideInGradebook                 bool
	PostToSIS                       bool
	ImportantDates                  bool
}

// AssignmentOverrideParams holds parameters for assignment overrides
type AssignmentOverrideParams struct {
	StudentIDs []int64
	Title      string
	DueAt      string
	UnlockAt   string
	LockAt     string
}

// Create creates a new assignment
func (s *AssignmentsService) Create(ctx context.Context, courseID int64, params *CreateAssignmentParams) (*Assignment, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/assignments", courseID)

	body := map[string]interface{}{
		"assignment": make(map[string]interface{}),
	}

	assignment := body["assignment"].(map[string]interface{})

	if params.Name != "" {
		assignment["name"] = params.Name
	}
	if params.Position > 0 {
		assignment["position"] = params.Position
	}
	if len(params.SubmissionTypes) > 0 {
		assignment["submission_types"] = params.SubmissionTypes
	}
	if len(params.AllowedExtensions) > 0 {
		assignment["allowed_extensions"] = params.AllowedExtensions
	}
	if params.TurnitinEnabled {
		assignment["turnitin_enabled"] = true
	}
	if params.VericiteEnabled {
		assignment["vericite_enabled"] = true
	}
	if len(params.TurnitinSettings) > 0 {
		assignment["turnitin_settings"] = params.TurnitinSettings
	}
	if len(params.IntegrationData) > 0 {
		assignment["integration_data"] = params.IntegrationData
	}
	if params.IntegrationID != "" {
		assignment["integration_id"] = params.IntegrationID
	}
	if params.PeerReviews {
		assignment["peer_reviews"] = true
	}
	if params.AutomaticPeerReviews {
		assignment["automatic_peer_reviews"] = true
	}
	if params.NotifyOfUpdate {
		assignment["notify_of_update"] = true
	}
	if params.GroupCategoryID > 0 {
		assignment["group_category_id"] = params.GroupCategoryID
	}
	if params.GradeGroupStudentsIndividually {
		assignment["grade_group_students_individually"] = true
	}
	if len(params.ExternalToolTagAttributes) > 0 {
		assignment["external_tool_tag_attributes"] = params.ExternalToolTagAttributes
	}
	if params.PointsPossible > 0 {
		assignment["points_possible"] = params.PointsPossible
	}
	if params.GradingType != "" {
		assignment["grading_type"] = params.GradingType
	}
	if params.DueAt != "" {
		assignment["due_at"] = params.DueAt
	}
	if params.LockAt != "" {
		assignment["lock_at"] = params.LockAt
	}
	if params.UnlockAt != "" {
		assignment["unlock_at"] = params.UnlockAt
	}
	if params.Description != "" {
		assignment["description"] = params.Description
	}
	if params.AssignmentGroupID > 0 {
		assignment["assignment_group_id"] = params.AssignmentGroupID
	}
	if len(params.AssignmentOverrides) > 0 {
		overrides := make([]map[string]interface{}, 0, len(params.AssignmentOverrides))
		for _, override := range params.AssignmentOverrides {
			o := make(map[string]interface{})
			if len(override.StudentIDs) > 0 {
				o["student_ids"] = override.StudentIDs
			}
			if override.Title != "" {
				o["title"] = override.Title
			}
			if override.DueAt != "" {
				o["due_at"] = override.DueAt
			}
			if override.UnlockAt != "" {
				o["unlock_at"] = override.UnlockAt
			}
			if override.LockAt != "" {
				o["lock_at"] = override.LockAt
			}
			overrides = append(overrides, o)
		}
		assignment["assignment_overrides"] = overrides
	}
	if params.OnlyVisibleToOverrides {
		assignment["only_visible_to_overrides"] = true
	}
	if params.Published {
		assignment["published"] = true
	}
	if params.GradingStandardID > 0 {
		assignment["grading_standard_id"] = params.GradingStandardID
	}
	if params.OmitFromFinalGrade {
		assignment["omit_from_final_grade"] = true
	}
	if params.ModeratedGrading {
		assignment["moderated_grading"] = true
	}
	if params.GraderCount > 0 {
		assignment["grader_count"] = params.GraderCount
	}
	if params.FinalGraderID > 0 {
		assignment["final_grader_id"] = params.FinalGraderID
	}
	if params.GraderCommentsVisibleToGraders {
		assignment["grader_comments_visible_to_graders"] = true
	}
	if params.GradersAnonymousToGraders {
		assignment["graders_anonymous_to_graders"] = true
	}
	if params.GraderNamesVisibleToFinalGrader {
		assignment["grader_names_visible_to_final_grader"] = true
	}
	if params.AnonymousInstructorAnnotations {
		assignment["anonymous_instructor_annotations"] = true
	}
	if params.AnonymousGrading {
		assignment["anonymous_grading"] = true
	}
	if params.AllowedAttempts > 0 {
		assignment["allowed_attempts"] = params.AllowedAttempts
	}
	if params.AnnotatableAttachmentID > 0 {
		assignment["annotatable_attachment_id"] = params.AnnotatableAttachmentID
	}
	if params.HideInGradebook {
		assignment["hide_in_gradebook"] = true
	}
	if params.PostToSIS {
		assignment["post_to_sis"] = true
	}
	if params.ImportantDates {
		assignment["important_dates"] = true
	}

	var result Assignment
	if err := s.client.PostJSON(ctx, path, body, &result); err != nil {
		return nil, err
	}

	return NormalizeAssignment(&result), nil
}

// UpdateAssignmentParams holds parameters for updating an assignment
type UpdateAssignmentParams struct {
	Name                            string
	Position                        *int
	SubmissionTypes                 []string
	AllowedExtensions               []string
	TurnitinEnabled                 *bool
	VericiteEnabled                 *bool
	TurnitinSettings                map[string]interface{}
	IntegrationData                 map[string]interface{}
	IntegrationID                   string
	PeerReviews                     *bool
	AutomaticPeerReviews            *bool
	NotifyOfUpdate                  *bool
	GroupCategoryID                 *int64
	GradeGroupStudentsIndividually  *bool
	ExternalToolTagAttributes       map[string]interface{}
	PointsPossible                  *float64
	GradingType                     string
	DueAt                           *string
	LockAt                          *string
	UnlockAt                        *string
	Description                     string
	AssignmentGroupID               *int64
	AssignmentOverrides             []AssignmentOverrideParams
	OnlyVisibleToOverrides          *bool
	Published                       *bool
	GradingStandardID               *int64
	OmitFromFinalGrade              *bool
	ModeratedGrading                *bool
	GraderCount                     *int
	FinalGraderID                   *int64
	GraderCommentsVisibleToGraders  *bool
	GradersAnonymousToGraders       *bool
	GraderNamesVisibleToFinalGrader *bool
	AnonymousInstructorAnnotations  *bool
	AnonymousGrading                *bool
	AllowedAttempts                 *int
	AnnotatableAttachmentID         *int64
	HideInGradebook                 *bool
	PostToSIS                       *bool
	ImportantDates                  *bool
}

// Update updates an existing assignment
func (s *AssignmentsService) Update(ctx context.Context, courseID, assignmentID int64, params *UpdateAssignmentParams) (*Assignment, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/assignments/%d", courseID, assignmentID)

	body := map[string]interface{}{
		"assignment": make(map[string]interface{}),
	}

	assignment := body["assignment"].(map[string]interface{})

	if params.Name != "" {
		assignment["name"] = params.Name
	}
	if params.Position != nil {
		assignment["position"] = *params.Position
	}
	if len(params.SubmissionTypes) > 0 {
		assignment["submission_types"] = params.SubmissionTypes
	}
	if len(params.AllowedExtensions) > 0 {
		assignment["allowed_extensions"] = params.AllowedExtensions
	}
	if params.TurnitinEnabled != nil {
		assignment["turnitin_enabled"] = *params.TurnitinEnabled
	}
	if params.VericiteEnabled != nil {
		assignment["vericite_enabled"] = *params.VericiteEnabled
	}
	if len(params.TurnitinSettings) > 0 {
		assignment["turnitin_settings"] = params.TurnitinSettings
	}
	if len(params.IntegrationData) > 0 {
		assignment["integration_data"] = params.IntegrationData
	}
	if params.IntegrationID != "" {
		assignment["integration_id"] = params.IntegrationID
	}
	if params.PeerReviews != nil {
		assignment["peer_reviews"] = *params.PeerReviews
	}
	if params.AutomaticPeerReviews != nil {
		assignment["automatic_peer_reviews"] = *params.AutomaticPeerReviews
	}
	if params.NotifyOfUpdate != nil {
		assignment["notify_of_update"] = *params.NotifyOfUpdate
	}
	if params.GroupCategoryID != nil {
		assignment["group_category_id"] = *params.GroupCategoryID
	}
	if params.GradeGroupStudentsIndividually != nil {
		assignment["grade_group_students_individually"] = *params.GradeGroupStudentsIndividually
	}
	if len(params.ExternalToolTagAttributes) > 0 {
		assignment["external_tool_tag_attributes"] = params.ExternalToolTagAttributes
	}
	if params.PointsPossible != nil {
		assignment["points_possible"] = *params.PointsPossible
	}
	if params.GradingType != "" {
		assignment["grading_type"] = params.GradingType
	}
	if params.DueAt != nil {
		assignment["due_at"] = *params.DueAt
	}
	if params.LockAt != nil {
		assignment["lock_at"] = *params.LockAt
	}
	if params.UnlockAt != nil {
		assignment["unlock_at"] = *params.UnlockAt
	}
	if params.Description != "" {
		assignment["description"] = params.Description
	}
	if params.AssignmentGroupID != nil {
		assignment["assignment_group_id"] = *params.AssignmentGroupID
	}
	if len(params.AssignmentOverrides) > 0 {
		overrides := make([]map[string]interface{}, 0, len(params.AssignmentOverrides))
		for _, override := range params.AssignmentOverrides {
			o := make(map[string]interface{})
			if len(override.StudentIDs) > 0 {
				o["student_ids"] = override.StudentIDs
			}
			if override.Title != "" {
				o["title"] = override.Title
			}
			if override.DueAt != "" {
				o["due_at"] = override.DueAt
			}
			if override.UnlockAt != "" {
				o["unlock_at"] = override.UnlockAt
			}
			if override.LockAt != "" {
				o["lock_at"] = override.LockAt
			}
			overrides = append(overrides, o)
		}
		assignment["assignment_overrides"] = overrides
	}
	if params.OnlyVisibleToOverrides != nil {
		assignment["only_visible_to_overrides"] = *params.OnlyVisibleToOverrides
	}
	if params.Published != nil {
		assignment["published"] = *params.Published
	}
	if params.GradingStandardID != nil {
		assignment["grading_standard_id"] = *params.GradingStandardID
	}
	if params.OmitFromFinalGrade != nil {
		assignment["omit_from_final_grade"] = *params.OmitFromFinalGrade
	}
	if params.ModeratedGrading != nil {
		assignment["moderated_grading"] = *params.ModeratedGrading
	}
	if params.GraderCount != nil {
		assignment["grader_count"] = *params.GraderCount
	}
	if params.FinalGraderID != nil {
		assignment["final_grader_id"] = *params.FinalGraderID
	}
	if params.GraderCommentsVisibleToGraders != nil {
		assignment["grader_comments_visible_to_graders"] = *params.GraderCommentsVisibleToGraders
	}
	if params.GradersAnonymousToGraders != nil {
		assignment["graders_anonymous_to_graders"] = *params.GradersAnonymousToGraders
	}
	if params.GraderNamesVisibleToFinalGrader != nil {
		assignment["grader_names_visible_to_final_grader"] = *params.GraderNamesVisibleToFinalGrader
	}
	if params.AnonymousInstructorAnnotations != nil {
		assignment["anonymous_instructor_annotations"] = *params.AnonymousInstructorAnnotations
	}
	if params.AnonymousGrading != nil {
		assignment["anonymous_grading"] = *params.AnonymousGrading
	}
	if params.AllowedAttempts != nil {
		assignment["allowed_attempts"] = *params.AllowedAttempts
	}
	if params.AnnotatableAttachmentID != nil {
		assignment["annotatable_attachment_id"] = *params.AnnotatableAttachmentID
	}
	if params.HideInGradebook != nil {
		assignment["hide_in_gradebook"] = *params.HideInGradebook
	}
	if params.PostToSIS != nil {
		assignment["post_to_sis"] = *params.PostToSIS
	}
	if params.ImportantDates != nil {
		assignment["important_dates"] = *params.ImportantDates
	}

	var result Assignment
	if err := s.client.PutJSON(ctx, path, body, &result); err != nil {
		return nil, err
	}

	return NormalizeAssignment(&result), nil
}

// Delete deletes an assignment
func (s *AssignmentsService) Delete(ctx context.Context, courseID, assignmentID int64) error {
	path := fmt.Sprintf("/api/v1/courses/%d/assignments/%d", courseID, assignmentID)

	_, err := s.client.Delete(ctx, path)
	return err
}

// BulkUpdateParams holds parameters for bulk updating assignment dates
type BulkUpdateParams struct {
	AssignmentIDs []int64
	DueAt         string
	UnlockAt      string
	LockAt        string
}

// BulkUpdate updates dates for multiple assignments at once
func (s *AssignmentsService) BulkUpdate(ctx context.Context, courseID int64, params *BulkUpdateParams) error {
	path := fmt.Sprintf("/api/v1/courses/%d/assignments/bulk_update", courseID)

	body := make(map[string]interface{})

	if len(params.AssignmentIDs) > 0 {
		ids := make([]string, len(params.AssignmentIDs))
		for i, id := range params.AssignmentIDs {
			ids[i] = strconv.FormatInt(id, 10)
		}
		body["assignment_ids[]"] = strings.Join(ids, ",")
	}

	if params.DueAt != "" {
		body["due_at"] = params.DueAt
	}
	if params.UnlockAt != "" {
		body["unlock_at"] = params.UnlockAt
	}
	if params.LockAt != "" {
		body["lock_at"] = params.LockAt
	}

	var result struct {
		Progress struct {
			ID int64 `json:"id"`
		} `json:"progress"`
	}
	if err := s.client.PutJSON(ctx, path, body, &result); err != nil {
		return err
	}

	return nil
}

// ListUserAssignments retrieves assignments for a specific user across all courses
func (s *AssignmentsService) ListUserAssignments(ctx context.Context, userID int64, opts *ListAssignmentsOptions) ([]Assignment, error) {
	path := fmt.Sprintf("/api/v1/users/%d/courses/%d/assignments", userID, userID)

	if opts != nil {
		query := url.Values{}

		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}

		if opts.Bucket != "" {
			query.Add("bucket", opts.Bucket)
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var assignments []Assignment
	if err := s.client.GetJSON(ctx, path, &assignments); err != nil {
		return nil, err
	}

	return NormalizeAssignments(assignments), nil
}
