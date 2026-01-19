package options

import "fmt"

// SubmissionsListOptions contains options for listing submissions
type SubmissionsListOptions struct {
	CourseID      int64
	AssignmentID  int64
	WorkflowState string
	GradedSince   string
	Include       []string
}

// Validate validates the options
func (o *SubmissionsListOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.AssignmentID <= 0 {
		return fmt.Errorf("assignment-id is required and must be greater than 0")
	}
	return nil
}

// SubmissionsGetOptions contains options for getting a submission
type SubmissionsGetOptions struct {
	CourseID     int64
	AssignmentID int64
	UserID       int64
	Include      []string
}

// Validate validates the options
func (o *SubmissionsGetOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.AssignmentID <= 0 {
		return fmt.Errorf("assignment-id is required and must be greater than 0")
	}
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}

// SubmissionsGradeOptions contains options for grading a submission
type SubmissionsGradeOptions struct {
	CourseID     int64
	AssignmentID int64
	UserID       int64
	Score        float64
	Comment      string
	Excuse       bool
	PostedGrade  string
}

// Validate validates the options
func (o *SubmissionsGradeOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.AssignmentID <= 0 {
		return fmt.Errorf("assignment-id is required and must be greater than 0")
	}
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	// At least one grading parameter is required
	if o.Score == 0 && o.Comment == "" && !o.Excuse && o.PostedGrade == "" {
		return fmt.Errorf("at least one grading parameter is required: score, comment, excuse, or posted-grade")
	}
	return nil
}

// SubmissionsBulkGradeOptions contains options for bulk grading submissions
type SubmissionsBulkGradeOptions struct {
	CourseID int64
	CSV      string
	DryRun   bool
}

// Validate validates the options
func (o *SubmissionsBulkGradeOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.CSV == "" {
		return fmt.Errorf("csv file is required")
	}
	return nil
}

// SubmissionsCommentsOptions contains options for listing submission comments
type SubmissionsCommentsOptions struct {
	CourseID     int64
	AssignmentID int64
	UserID       int64
}

// Validate validates the options
func (o *SubmissionsCommentsOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.AssignmentID <= 0 {
		return fmt.Errorf("assignment-id is required and must be greater than 0")
	}
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}

// SubmissionsAddCommentOptions contains options for adding a comment to a submission
type SubmissionsAddCommentOptions struct {
	CourseID     int64
	AssignmentID int64
	UserID       int64
	Text         string
	GroupShare   bool
}

// Validate validates the options
func (o *SubmissionsAddCommentOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.AssignmentID <= 0 {
		return fmt.Errorf("assignment-id is required and must be greater than 0")
	}
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	if o.Text == "" {
		return fmt.Errorf("comment text is required")
	}
	return nil
}

// SubmissionsDeleteCommentOptions contains options for deleting a submission comment
type SubmissionsDeleteCommentOptions struct {
	CourseID     int64
	AssignmentID int64
	UserID       int64
	CommentID    int64
}

// Validate validates the options
func (o *SubmissionsDeleteCommentOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.AssignmentID <= 0 {
		return fmt.Errorf("assignment-id is required and must be greater than 0")
	}
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	if o.CommentID <= 0 {
		return fmt.Errorf("comment-id is required and must be greater than 0")
	}
	return nil
}
