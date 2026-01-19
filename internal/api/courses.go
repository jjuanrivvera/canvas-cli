package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// CoursesService handles course-related API calls
type CoursesService struct {
	client *Client
}

// NewCoursesService creates a new courses service
func NewCoursesService(client *Client) *CoursesService {
	return &CoursesService{client: client}
}

// ListCoursesOptions holds options for listing courses
type ListCoursesOptions struct {
	EnrollmentType  string   // student, teacher, ta, observer, designer
	EnrollmentState string   // active, invited_or_pending, completed
	Include         []string // needs_grading_count, syllabus_body, total_scores, term, etc.
	State           []string // unpublished, available, completed, deleted
	Page            int
	PerPage         int
}

// List retrieves all courses for the current user
func (s *CoursesService) List(ctx context.Context, opts *ListCoursesOptions) ([]Course, error) {
	path := "/api/v1/courses"

	if opts != nil {
		query := url.Values{}

		if opts.EnrollmentType != "" {
			query.Add("enrollment_type", opts.EnrollmentType)
		}

		if opts.EnrollmentState != "" {
			query.Add("enrollment_state", opts.EnrollmentState)
		}

		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}

		for _, state := range opts.State {
			query.Add("state[]", state)
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

	courses, err := GetAllPagesGeneric[Course](s.client, ctx, path)
	if err != nil {
		return nil, err
	}

	return NormalizeCourses(courses), nil
}

// Get retrieves a single course by ID
func (s *CoursesService) Get(ctx context.Context, courseID int64, include []string) (*Course, error) {
	path := fmt.Sprintf("/api/v1/courses/%d", courseID)

	if len(include) > 0 {
		query := url.Values{}
		for _, inc := range include {
			query.Add("include[]", inc)
		}
		path += "?" + query.Encode()
	}

	var course Course
	if err := s.client.GetJSON(ctx, path, &course); err != nil {
		return nil, err
	}

	return NormalizeCourse(&course), nil
}

// CreateCourseParams holds parameters for creating a course
type CreateCourseParams struct {
	AccountID                        int64
	Name                             string
	CourseCode                       string
	StartAt                          string
	EndAt                            string
	License                          string
	IsPublic                         bool
	IsPublicToAuthUsers              bool
	PublicSyllabus                   bool
	PublicSyllabusToAuth             bool
	PublicDescription                string
	AllowStudentWikiEdits            bool
	AllowWikiComments                bool
	AllowStudentForumAttachments     bool
	OpenEnrollment                   bool
	SelfEnrollment                   bool
	RestrictEnrollmentsToCourseDates bool
	TermID                           int64
	SISCourseID                      string
	IntegrationID                    string
	HideFinalGrades                  bool
	ApplyAssignmentGroupWeights      bool
	TimeZone                         string
	Offer                            bool
	EnrollMe                         bool
	DefaultView                      string
	SyllabusBody                     string
	GradingStandardID                int64
	CourseFormat                     string
}

// addCourseStringField adds a string field to courseData if non-empty
func addCourseStringField(courseData map[string]interface{}, key, value string) {
	if value != "" {
		courseData[key] = value
	}
}

// addCourseBoolField adds a bool field to courseData if true
func addCourseBoolField(courseData map[string]interface{}, key string, value bool) {
	if value {
		courseData[key] = true
	}
}

// addCourseBoolPtrField adds a bool pointer field to courseData if not nil
func addCourseBoolPtrField(courseData map[string]interface{}, key string, value *bool) {
	if value != nil {
		courseData[key] = *value
	}
}

// addCourseInt64Field adds an int64 field to courseData if positive
func addCourseInt64Field(courseData map[string]interface{}, key string, value int64) {
	if value > 0 {
		courseData[key] = value
	}
}

// Create creates a new course
func (s *CoursesService) Create(ctx context.Context, params *CreateCourseParams) (*Course, error) {
	if params.AccountID == 0 {
		return nil, fmt.Errorf("account_id is required")
	}

	path := fmt.Sprintf("/api/v1/accounts/%d/courses", params.AccountID)

	body := map[string]interface{}{
		"course": make(map[string]interface{}),
	}

	courseData := body["course"].(map[string]interface{})

	// String fields
	addCourseStringField(courseData, "name", params.Name)
	addCourseStringField(courseData, "course_code", params.CourseCode)
	addCourseStringField(courseData, "start_at", params.StartAt)
	addCourseStringField(courseData, "end_at", params.EndAt)
	addCourseStringField(courseData, "license", params.License)
	addCourseStringField(courseData, "public_description", params.PublicDescription)
	addCourseStringField(courseData, "sis_course_id", params.SISCourseID)
	addCourseStringField(courseData, "integration_id", params.IntegrationID)
	addCourseStringField(courseData, "time_zone", params.TimeZone)
	addCourseStringField(courseData, "default_view", params.DefaultView)
	addCourseStringField(courseData, "syllabus_body", params.SyllabusBody)
	addCourseStringField(courseData, "course_format", params.CourseFormat)

	// Boolean fields
	addCourseBoolField(courseData, "is_public", params.IsPublic)
	addCourseBoolField(courseData, "is_public_to_auth_users", params.IsPublicToAuthUsers)
	addCourseBoolField(courseData, "public_syllabus", params.PublicSyllabus)
	addCourseBoolField(courseData, "public_syllabus_to_auth", params.PublicSyllabusToAuth)
	addCourseBoolField(courseData, "allow_student_wiki_edits", params.AllowStudentWikiEdits)
	addCourseBoolField(courseData, "allow_wiki_comments", params.AllowWikiComments)
	addCourseBoolField(courseData, "allow_student_forum_attachments", params.AllowStudentForumAttachments)
	addCourseBoolField(courseData, "open_enrollment", params.OpenEnrollment)
	addCourseBoolField(courseData, "self_enrollment", params.SelfEnrollment)
	addCourseBoolField(courseData, "restrict_enrollments_to_course_dates", params.RestrictEnrollmentsToCourseDates)
	addCourseBoolField(courseData, "hide_final_grades", params.HideFinalGrades)
	addCourseBoolField(courseData, "apply_assignment_group_weights", params.ApplyAssignmentGroupWeights)
	addCourseBoolField(courseData, "offer", params.Offer)
	addCourseBoolField(courseData, "enroll_me", params.EnrollMe)

	// Integer fields
	addCourseInt64Field(courseData, "term_id", params.TermID)
	addCourseInt64Field(courseData, "grading_standard_id", params.GradingStandardID)

	var course Course
	if err := s.client.PostJSON(ctx, path, body, &course); err != nil {
		return nil, err
	}

	return NormalizeCourse(&course), nil
}

// UpdateCourseParams holds parameters for updating a course
type UpdateCourseParams struct {
	Name                             string
	CourseCode                       string
	StartAt                          string
	EndAt                            string
	License                          string
	IsPublic                         *bool
	IsPublicToAuthUsers              *bool
	PublicSyllabus                   *bool
	PublicSyllabusToAuth             *bool
	PublicDescription                string
	AllowStudentWikiEdits            *bool
	AllowWikiComments                *bool
	AllowStudentForumAttachments     *bool
	OpenEnrollment                   *bool
	SelfEnrollment                   *bool
	RestrictEnrollmentsToCourseDates *bool
	HideFinalGrades                  *bool
	ApplyAssignmentGroupWeights      *bool
	TimeZone                         string
	DefaultView                      string
	SyllabusBody                     string
	GradingStandardID                int64
	CourseFormat                     string
	ImageID                          int64
	ImageURL                         string
	RemoveImage                      bool
}

// Update updates an existing course
func (s *CoursesService) Update(ctx context.Context, courseID int64, params *UpdateCourseParams) (*Course, error) {
	path := fmt.Sprintf("/api/v1/courses/%d", courseID)

	body := map[string]interface{}{
		"course": make(map[string]interface{}),
	}

	courseData := body["course"].(map[string]interface{})

	// String fields
	addCourseStringField(courseData, "name", params.Name)
	addCourseStringField(courseData, "course_code", params.CourseCode)
	addCourseStringField(courseData, "start_at", params.StartAt)
	addCourseStringField(courseData, "end_at", params.EndAt)
	addCourseStringField(courseData, "license", params.License)
	addCourseStringField(courseData, "public_description", params.PublicDescription)
	addCourseStringField(courseData, "time_zone", params.TimeZone)
	addCourseStringField(courseData, "default_view", params.DefaultView)
	addCourseStringField(courseData, "syllabus_body", params.SyllabusBody)
	addCourseStringField(courseData, "course_format", params.CourseFormat)
	addCourseStringField(courseData, "image_url", params.ImageURL)

	// Boolean pointer fields (allow explicit false)
	addCourseBoolPtrField(courseData, "is_public", params.IsPublic)
	addCourseBoolPtrField(courseData, "is_public_to_auth_users", params.IsPublicToAuthUsers)
	addCourseBoolPtrField(courseData, "public_syllabus", params.PublicSyllabus)
	addCourseBoolPtrField(courseData, "public_syllabus_to_auth", params.PublicSyllabusToAuth)
	addCourseBoolPtrField(courseData, "allow_student_wiki_edits", params.AllowStudentWikiEdits)
	addCourseBoolPtrField(courseData, "allow_wiki_comments", params.AllowWikiComments)
	addCourseBoolPtrField(courseData, "allow_student_forum_attachments", params.AllowStudentForumAttachments)
	addCourseBoolPtrField(courseData, "open_enrollment", params.OpenEnrollment)
	addCourseBoolPtrField(courseData, "self_enrollment", params.SelfEnrollment)
	addCourseBoolPtrField(courseData, "restrict_enrollments_to_course_dates", params.RestrictEnrollmentsToCourseDates)
	addCourseBoolPtrField(courseData, "hide_final_grades", params.HideFinalGrades)
	addCourseBoolPtrField(courseData, "apply_assignment_group_weights", params.ApplyAssignmentGroupWeights)

	// Boolean field
	addCourseBoolField(courseData, "remove_image", params.RemoveImage)

	// Integer fields
	addCourseInt64Field(courseData, "grading_standard_id", params.GradingStandardID)
	addCourseInt64Field(courseData, "image_id", params.ImageID)

	var course Course
	if err := s.client.PutJSON(ctx, path, body, &course); err != nil {
		return nil, err
	}

	return NormalizeCourse(&course), nil
}

// Delete deletes a course (sets to deleted state)
func (s *CoursesService) Delete(ctx context.Context, courseID int64, event string) error {
	path := fmt.Sprintf("/api/v1/courses/%d", courseID)

	if event != "" {
		path += "?event=" + url.QueryEscape(event)
	}

	_, err := s.client.Delete(ctx, path)
	return err
}
