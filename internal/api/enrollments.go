package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// EnrollmentsService handles enrollment-related API calls
type EnrollmentsService struct {
	client *Client
}

// NewEnrollmentsService creates a new enrollments service
func NewEnrollmentsService(client *Client) *EnrollmentsService {
	return &EnrollmentsService{client: client}
}

// ListEnrollmentsOptions holds options for listing enrollments
type ListEnrollmentsOptions struct {
	Type             []string // Enrollment types to include (StudentEnrollment, TeacherEnrollment, TaEnrollment, DesignerEnrollment, ObserverEnrollment)
	Role             []string // Deprecated, use RoleID instead
	RoleID           []int64  // Filter by role ID
	State            []string // Filter by enrollment state (active, invited, creation_pending, deleted, rejected, completed, inactive, current_and_future, etc.)
	Include          []string // Additional data to include (avatar_url, group_ids, locked, observed_users, can_be_removed, uuid, current_points)
	UserID           int64    // Filter by user ID
	GradingPeriodID  int64    // Filter by grading period
	EnrollmentTermID int64    // Filter by enrollment term
	SISAccountID     []string // Filter by SIS account ID
	SISCourseID      []string // Filter by SIS course ID
	SISSectionID     []string // Filter by SIS section ID
	SISUserID        []string // Filter by SIS user ID
	CreatedForSISID  []string // Filter by created_for_sis_id
	Page             int
	PerPage          int
}

// ListCourse retrieves enrollments for a course
func (s *EnrollmentsService) ListCourse(ctx context.Context, courseID int64, opts *ListEnrollmentsOptions) ([]Enrollment, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/enrollments", courseID)

	if opts != nil {
		query := url.Values{}

		for _, t := range opts.Type {
			query.Add("type[]", t)
		}

		for _, r := range opts.Role {
			query.Add("role[]", r)
		}

		for _, id := range opts.RoleID {
			query.Add("role_id[]", strconv.FormatInt(id, 10))
		}

		for _, st := range opts.State {
			query.Add("state[]", st)
		}

		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}

		if opts.UserID > 0 {
			query.Add("user_id", strconv.FormatInt(opts.UserID, 10))
		}

		if opts.GradingPeriodID > 0 {
			query.Add("grading_period_id", strconv.FormatInt(opts.GradingPeriodID, 10))
		}

		if opts.EnrollmentTermID > 0 {
			query.Add("enrollment_term_id", strconv.FormatInt(opts.EnrollmentTermID, 10))
		}

		for _, id := range opts.SISAccountID {
			query.Add("sis_account_id[]", id)
		}

		for _, id := range opts.SISCourseID {
			query.Add("sis_course_id[]", id)
		}

		for _, id := range opts.SISSectionID {
			query.Add("sis_section_id[]", id)
		}

		for _, id := range opts.SISUserID {
			query.Add("sis_user_id[]", id)
		}

		for _, id := range opts.CreatedForSISID {
			query.Add("created_for_sis_id[]", id)
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

	var enrollments []Enrollment
	if err := s.client.GetAllPages(ctx, path, &enrollments); err != nil {
		return nil, err
	}

	return NormalizeEnrollments(enrollments), nil
}

// ListSection retrieves enrollments for a section
func (s *EnrollmentsService) ListSection(ctx context.Context, sectionID int64, opts *ListEnrollmentsOptions) ([]Enrollment, error) {
	path := fmt.Sprintf("/api/v1/sections/%d/enrollments", sectionID)

	if opts != nil {
		query := url.Values{}

		for _, t := range opts.Type {
			query.Add("type[]", t)
		}

		for _, r := range opts.Role {
			query.Add("role[]", r)
		}

		for _, id := range opts.RoleID {
			query.Add("role_id[]", strconv.FormatInt(id, 10))
		}

		for _, st := range opts.State {
			query.Add("state[]", st)
		}

		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}

		if opts.UserID > 0 {
			query.Add("user_id", strconv.FormatInt(opts.UserID, 10))
		}

		if opts.GradingPeriodID > 0 {
			query.Add("grading_period_id", strconv.FormatInt(opts.GradingPeriodID, 10))
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

	var enrollments []Enrollment
	if err := s.client.GetAllPages(ctx, path, &enrollments); err != nil {
		return nil, err
	}

	return NormalizeEnrollments(enrollments), nil
}

// ListUser retrieves enrollments for a user
func (s *EnrollmentsService) ListUser(ctx context.Context, userID int64, opts *ListEnrollmentsOptions) ([]Enrollment, error) {
	path := fmt.Sprintf("/api/v1/users/%d/enrollments", userID)

	if opts != nil {
		query := url.Values{}

		for _, t := range opts.Type {
			query.Add("type[]", t)
		}

		for _, r := range opts.Role {
			query.Add("role[]", r)
		}

		for _, id := range opts.RoleID {
			query.Add("role_id[]", strconv.FormatInt(id, 10))
		}

		for _, st := range opts.State {
			query.Add("state[]", st)
		}

		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}

		if opts.GradingPeriodID > 0 {
			query.Add("grading_period_id", strconv.FormatInt(opts.GradingPeriodID, 10))
		}

		if opts.EnrollmentTermID > 0 {
			query.Add("enrollment_term_id", strconv.FormatInt(opts.EnrollmentTermID, 10))
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

	var enrollments []Enrollment
	if err := s.client.GetAllPages(ctx, path, &enrollments); err != nil {
		return nil, err
	}

	return NormalizeEnrollments(enrollments), nil
}

// EnrollUserParams holds parameters for enrolling a user
type EnrollUserParams struct {
	UserID                         int64
	Type                           string // StudentEnrollment, TeacherEnrollment, TaEnrollment, ObserverEnrollment, DesignerEnrollment
	RoleID                         int64
	EnrollmentState                string // active, invited, inactive
	CourseSectionID                int64
	LimitPrivilegesToCourseSection bool
	Notify                         bool
	SelfEnrollmentCode             string
	SelfEnrolled                   bool
	AssociatedUserID               int64 // For observer enrollments
}

// EnrollUser enrolls a user in a course
func (s *EnrollmentsService) EnrollUser(ctx context.Context, courseID int64, params *EnrollUserParams) (*Enrollment, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/enrollments", courseID)

	body := map[string]interface{}{
		"enrollment": make(map[string]interface{}),
	}

	enrollment := body["enrollment"].(map[string]interface{})

	if params.UserID > 0 {
		enrollment["user_id"] = params.UserID
	}

	if params.Type != "" {
		enrollment["type"] = params.Type
	}

	if params.RoleID > 0 {
		enrollment["role_id"] = params.RoleID
	}

	if params.EnrollmentState != "" {
		enrollment["enrollment_state"] = params.EnrollmentState
	}

	if params.CourseSectionID > 0 {
		enrollment["course_section_id"] = params.CourseSectionID
	}

	if params.LimitPrivilegesToCourseSection {
		enrollment["limit_privileges_to_course_section"] = true
	}

	if params.Notify {
		enrollment["notify"] = true
	}

	if params.SelfEnrollmentCode != "" {
		enrollment["self_enrollment_code"] = params.SelfEnrollmentCode
	}

	if params.SelfEnrolled {
		enrollment["self_enrolled"] = true
	}

	if params.AssociatedUserID > 0 {
		enrollment["associated_user_id"] = params.AssociatedUserID
	}

	var result Enrollment
	if err := s.client.PostJSON(ctx, path, body, &result); err != nil {
		return nil, err
	}

	return NormalizeEnrollment(&result), nil
}

// ConcludeParams holds parameters for concluding an enrollment
type ConcludeParams struct {
	Task string // conclude, delete, inactivate, deactivate
}

// Conclude concludes, deletes, or deactivates an enrollment
func (s *EnrollmentsService) Conclude(ctx context.Context, courseID, enrollmentID int64, task string) (*Enrollment, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/enrollments/%d", courseID, enrollmentID)

	query := url.Values{}
	query.Add("task", task)
	path += "?" + query.Encode()

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result Enrollment
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return NormalizeEnrollment(&result), nil
}

// ReactivateParams holds parameters for reactivating an enrollment
type ReactivateParams struct {
	// Empty for now, may add fields in future
}

// Reactivate reactivates a deactivated enrollment
func (s *EnrollmentsService) Reactivate(ctx context.Context, courseID, enrollmentID int64) (*Enrollment, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/enrollments/%d/reactivate", courseID, enrollmentID)

	var result Enrollment
	if err := s.client.PutJSON(ctx, path, nil, &result); err != nil {
		return nil, err
	}

	return NormalizeEnrollment(&result), nil
}

// AcceptParams holds parameters for accepting an enrollment invitation
type AcceptParams struct {
	// Empty for now, may add fields in future
}

// Accept accepts a pending enrollment invitation
func (s *EnrollmentsService) Accept(ctx context.Context, courseID, enrollmentID int64) error {
	path := fmt.Sprintf("/api/v1/courses/%d/enrollments/%d/accept", courseID, enrollmentID)

	return s.client.PostJSON(ctx, path, nil, nil)
}

// RejectParams holds parameters for rejecting an enrollment invitation
type RejectParams struct {
	// Empty for now, may add fields in future
}

// Reject rejects a pending enrollment invitation
func (s *EnrollmentsService) Reject(ctx context.Context, courseID, enrollmentID int64) error {
	path := fmt.Sprintf("/api/v1/courses/%d/enrollments/%d/reject", courseID, enrollmentID)

	return s.client.PostJSON(ctx, path, nil, nil)
}

// LastAttendedParams holds parameters for updating last attended date
type LastAttendedParams struct {
	Date string // ISO8601 format
}

// UpdateLastAttended updates the last attended date for an enrollment
func (s *EnrollmentsService) UpdateLastAttended(ctx context.Context, courseID, enrollmentID int64, date string) (*Enrollment, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/enrollments/%d/last_attended", courseID, enrollmentID)

	body := map[string]interface{}{
		"date": date,
	}

	var result Enrollment
	if err := s.client.PutJSON(ctx, path, body, &result); err != nil {
		return nil, err
	}

	return NormalizeEnrollment(&result), nil
}
