package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// SubmissionsService handles submission-related API calls
type SubmissionsService struct {
	client *Client
}

// NewSubmissionsService creates a new submissions service
func NewSubmissionsService(client *Client) *SubmissionsService {
	return &SubmissionsService{client: client}
}

// ListSubmissionsOptions holds options for listing submissions
type ListSubmissionsOptions struct {
	Include          []string // Additional data to include (submission_history, submission_comments, rubric_assessment, assignment, visibility, course, user, etc.)
	Grouped          bool     // Group submissions by student
	PostToSIS        *bool    // Filter by post_to_sis
	SubmittedSince   string   // ISO8601 timestamp
	GradedSince      string   // ISO8601 timestamp
	GradingPeriodID  int64    // Filter by grading period
	WorkflowState    string   // Filter by workflow state (submitted, unsubmitted, graded, pending_review)
	EnrollmentState  string   // Filter by enrollment state (active, concluded)
	StateBasedOnDate bool     // If true, filter by state based on assignment due date
	Order            string   // Order by (id, graded_at)
	OrderDirection   string   // Order direction (ascending, descending)
	Page             int
	PerPage          int
}

// Get retrieves a single submission
func (s *SubmissionsService) Get(ctx context.Context, courseID, assignmentID, userID int64, include []string) (*Submission, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/assignments/%d/submissions/%d", courseID, assignmentID, userID)

	if len(include) > 0 {
		query := url.Values{}
		for _, inc := range include {
			query.Add("include[]", inc)
		}
		path += "?" + query.Encode()
	}

	var submission Submission
	if err := s.client.GetJSON(ctx, path, &submission); err != nil {
		return nil, err
	}

	return NormalizeSubmission(&submission), nil
}

// List retrieves all submissions for an assignment
func (s *SubmissionsService) List(ctx context.Context, courseID, assignmentID int64, opts *ListSubmissionsOptions) ([]Submission, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/assignments/%d/submissions", courseID, assignmentID)

	if opts != nil {
		query := url.Values{}

		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}

		if opts.Grouped {
			query.Add("grouped", "true")
		}

		if opts.PostToSIS != nil {
			query.Add("post_to_sis", strconv.FormatBool(*opts.PostToSIS))
		}

		if opts.SubmittedSince != "" {
			query.Add("submitted_since", opts.SubmittedSince)
		}

		if opts.GradedSince != "" {
			query.Add("graded_since", opts.GradedSince)
		}

		if opts.GradingPeriodID > 0 {
			query.Add("grading_period_id", strconv.FormatInt(opts.GradingPeriodID, 10))
		}

		if opts.WorkflowState != "" {
			query.Add("workflow_state", opts.WorkflowState)
		}

		if opts.EnrollmentState != "" {
			query.Add("enrollment_state", opts.EnrollmentState)
		}

		if opts.StateBasedOnDate {
			query.Add("state_based_on_date", "true")
		}

		if opts.Order != "" {
			query.Add("order", opts.Order)
		}

		if opts.OrderDirection != "" {
			query.Add("order_direction", opts.OrderDirection)
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

	var submissions []Submission
	if err := s.client.GetAllPages(ctx, path, &submissions); err != nil {
		return nil, err
	}

	return NormalizeSubmissions(submissions), nil
}

// ListMultiple retrieves submissions for multiple assignments and users
func (s *SubmissionsService) ListMultiple(ctx context.Context, courseID int64, studentIDs, assignmentIDs []int64, opts *ListSubmissionsOptions) ([]Submission, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/students/submissions", courseID)

	query := url.Values{}

	for _, id := range studentIDs {
		query.Add("student_ids[]", strconv.FormatInt(id, 10))
	}

	for _, id := range assignmentIDs {
		query.Add("assignment_ids[]", strconv.FormatInt(id, 10))
	}

	if opts != nil {
		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}

		if opts.Grouped {
			query.Add("grouped", "true")
		}

		if opts.PostToSIS != nil {
			query.Add("post_to_sis", strconv.FormatBool(*opts.PostToSIS))
		}

		if opts.SubmittedSince != "" {
			query.Add("submitted_since", opts.SubmittedSince)
		}

		if opts.GradedSince != "" {
			query.Add("graded_since", opts.GradedSince)
		}

		if opts.GradingPeriodID > 0 {
			query.Add("grading_period_id", strconv.FormatInt(opts.GradingPeriodID, 10))
		}

		if opts.WorkflowState != "" {
			query.Add("workflow_state", opts.WorkflowState)
		}

		if opts.EnrollmentState != "" {
			query.Add("enrollment_state", opts.EnrollmentState)
		}

		if opts.StateBasedOnDate {
			query.Add("state_based_on_date", "true")
		}

		if opts.Order != "" {
			query.Add("order", opts.Order)
		}

		if opts.OrderDirection != "" {
			query.Add("order_direction", opts.OrderDirection)
		}

		if opts.Page > 0 {
			query.Add("page", strconv.Itoa(opts.Page))
		}

		if opts.PerPage > 0 {
			query.Add("per_page", strconv.Itoa(opts.PerPage))
		}
	}

	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	var submissions []Submission
	if err := s.client.GetAllPages(ctx, path, &submissions); err != nil {
		return nil, err
	}

	return NormalizeSubmissions(submissions), nil
}

// GradeSubmissionParams holds parameters for grading a submission
type GradeSubmissionParams struct {
	PostedGrade         string // The grade to assign (letter grade, percentage, points, etc.)
	Excuse              bool   // Excuse the submission
	LatePolicyStatus    string // late, missing, none
	SecondsLateOverride *int   // Override seconds late calculation
	Comment             *SubmissionCommentParams
	RubricAssessment    map[string]RubricAssessmentParams
}

// SubmissionCommentParams holds parameters for adding a submission comment
type SubmissionCommentParams struct {
	TextComment      string
	GroupComment     bool
	MediaCommentID   string
	MediaCommentType string
	FileIDs          []int64
}

// RubricAssessmentParams holds parameters for rubric assessment
type RubricAssessmentParams struct {
	Points   float64
	Rating   string
	Comments string
}

// Grade grades a submission
func (s *SubmissionsService) Grade(ctx context.Context, courseID, assignmentID, userID int64, params *GradeSubmissionParams) (*Submission, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/assignments/%d/submissions/%d", courseID, assignmentID, userID)

	body := map[string]interface{}{
		"submission": make(map[string]interface{}),
	}

	submission := body["submission"].(map[string]interface{})

	if params.PostedGrade != "" {
		submission["posted_grade"] = params.PostedGrade
	}

	if params.Excuse {
		submission["excuse"] = true
	}

	if params.LatePolicyStatus != "" {
		submission["late_policy_status"] = params.LatePolicyStatus
	}

	if params.SecondsLateOverride != nil {
		submission["seconds_late_override"] = *params.SecondsLateOverride
	}

	if params.Comment != nil {
		comment := make(map[string]interface{})
		if params.Comment.TextComment != "" {
			comment["text_comment"] = params.Comment.TextComment
		}
		if params.Comment.GroupComment {
			comment["group_comment"] = true
		}
		if params.Comment.MediaCommentID != "" {
			comment["media_comment_id"] = params.Comment.MediaCommentID
		}
		if params.Comment.MediaCommentType != "" {
			comment["media_comment_type"] = params.Comment.MediaCommentType
		}
		if len(params.Comment.FileIDs) > 0 {
			comment["file_ids"] = params.Comment.FileIDs
		}
		body["comment"] = comment
	}

	if len(params.RubricAssessment) > 0 {
		assessment := make(map[string]interface{})
		for criterionID, criterion := range params.RubricAssessment {
			criterionData := make(map[string]interface{})
			if criterion.Points > 0 {
				criterionData["points"] = criterion.Points
			}
			if criterion.Rating != "" {
				criterionData["rating_id"] = criterion.Rating
			}
			if criterion.Comments != "" {
				criterionData["comments"] = criterion.Comments
			}
			assessment[criterionID] = criterionData
		}
		body["rubric_assessment"] = assessment
	}

	var result Submission
	if err := s.client.PutJSON(ctx, path, body, &result); err != nil {
		return nil, err
	}

	return NormalizeSubmission(&result), nil
}

// BulkGradeParams holds parameters for bulk grading
type BulkGradeParams struct {
	GradeData map[int64]GradeData // Map of user ID to grade data
}

// GradeData holds grade information for a single student
type GradeData struct {
	PostedGrade      string
	Excuse           bool
	LatePolicyStatus string
	RubricAssessment map[string]RubricAssessmentParams
}

// BulkGrade grades multiple submissions at once
func (s *SubmissionsService) BulkGrade(ctx context.Context, courseID, assignmentID int64, params *BulkGradeParams) ([]Submission, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/assignments/%d/submissions/update_grades", courseID, assignmentID)

	body := map[string]interface{}{
		"grade_data": make(map[string]interface{}),
	}

	gradeData := body["grade_data"].(map[string]interface{})

	for userID, data := range params.GradeData {
		userKey := strconv.FormatInt(userID, 10)
		userData := make(map[string]interface{})

		if data.PostedGrade != "" {
			userData["posted_grade"] = data.PostedGrade
		}
		if data.Excuse {
			userData["excuse"] = true
		}
		if data.LatePolicyStatus != "" {
			userData["late_policy_status"] = data.LatePolicyStatus
		}

		if len(data.RubricAssessment) > 0 {
			assessment := make(map[string]interface{})
			for criterionID, criterion := range data.RubricAssessment {
				criterionData := make(map[string]interface{})
				if criterion.Points > 0 {
					criterionData["points"] = criterion.Points
				}
				if criterion.Rating != "" {
					criterionData["rating_id"] = criterion.Rating
				}
				if criterion.Comments != "" {
					criterionData["comments"] = criterion.Comments
				}
				assessment[criterionID] = criterionData
			}
			userData["rubric_assessment"] = assessment
		}

		gradeData[userKey] = userData
	}

	var result struct {
		Progress struct {
			ID int64 `json:"id"`
		} `json:"progress"`
	}
	if err := s.client.PostJSON(ctx, path, body, &result); err != nil {
		return nil, err
	}

	// Return empty slice for now, actual progress would need to be polled
	return []Submission{}, nil
}

// SubmitParams holds parameters for submitting an assignment
type SubmitParams struct {
	SubmissionType   string  // online_text_entry, online_url, online_upload, media_recording
	Body             string  // For online_text_entry
	URL              string  // For online_url
	FileIDs          []int64 // For online_upload
	MediaCommentID   string  // For media_recording
	MediaCommentType string  // audio or video
	UserID           int64   // Submit on behalf of user (requires permission)
	Comment          *SubmissionCommentParams
}

// Submit submits an assignment
func (s *SubmissionsService) Submit(ctx context.Context, courseID, assignmentID int64, params *SubmitParams) (*Submission, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/assignments/%d/submissions", courseID, assignmentID)

	body := map[string]interface{}{
		"submission": make(map[string]interface{}),
	}

	submission := body["submission"].(map[string]interface{})

	if params.SubmissionType != "" {
		submission["submission_type"] = params.SubmissionType
	}

	if params.Body != "" {
		submission["body"] = params.Body
	}

	if params.URL != "" {
		submission["url"] = params.URL
	}

	if len(params.FileIDs) > 0 {
		submission["file_ids"] = params.FileIDs
	}

	if params.MediaCommentID != "" {
		submission["media_comment_id"] = params.MediaCommentID
	}

	if params.MediaCommentType != "" {
		submission["media_comment_type"] = params.MediaCommentType
	}

	if params.UserID > 0 {
		submission["user_id"] = params.UserID
	}

	if params.Comment != nil {
		comment := make(map[string]interface{})
		if params.Comment.TextComment != "" {
			comment["text_comment"] = params.Comment.TextComment
		}
		if params.Comment.GroupComment {
			comment["group_comment"] = true
		}
		if params.Comment.MediaCommentID != "" {
			comment["media_comment_id"] = params.Comment.MediaCommentID
		}
		if params.Comment.MediaCommentType != "" {
			comment["media_comment_type"] = params.Comment.MediaCommentType
		}
		if len(params.Comment.FileIDs) > 0 {
			comment["file_ids"] = params.Comment.FileIDs
		}
		body["comment"] = comment
	}

	var result Submission
	if err := s.client.PostJSON(ctx, path, body, &result); err != nil {
		return nil, err
	}

	return NormalizeSubmission(&result), nil
}

// MarkAsRead marks a submission as read
func (s *SubmissionsService) MarkAsRead(ctx context.Context, courseID, assignmentID, userID int64) error {
	path := fmt.Sprintf("/api/v1/courses/%d/assignments/%d/submissions/%d/read", courseID, assignmentID, userID)

	return s.client.PutJSON(ctx, path, nil, nil)
}

// MarkAsUnread marks a submission as unread
func (s *SubmissionsService) MarkAsUnread(ctx context.Context, courseID, assignmentID, userID int64) error {
	path := fmt.Sprintf("/api/v1/courses/%d/assignments/%d/submissions/%d/read", courseID, assignmentID, userID)

	_, err := s.client.Delete(ctx, path)
	return err
}

// UploadFileParams holds parameters for uploading a file
type UploadFileParams struct {
	Name        string
	Size        int64
	ContentType string
	OnDuplicate string // overwrite, rename
}

// InitiateFileUpload initiates a file upload for submission
func (s *SubmissionsService) InitiateFileUpload(ctx context.Context, courseID, assignmentID int64, params *UploadFileParams) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/assignments/%d/submissions/self/files", courseID, assignmentID)

	body := make(map[string]interface{})

	if params.Name != "" {
		body["name"] = params.Name
	}
	if params.Size > 0 {
		body["size"] = params.Size
	}
	if params.ContentType != "" {
		body["content_type"] = params.ContentType
	}
	if params.OnDuplicate != "" {
		body["on_duplicate"] = params.OnDuplicate
	}

	var result map[string]interface{}
	if err := s.client.PostJSON(ctx, path, body, &result); err != nil {
		return nil, err
	}

	return result, nil
}
