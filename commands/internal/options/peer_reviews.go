package options

// PeerReviewsListOptions contains options for listing peer reviews
type PeerReviewsListOptions struct {
	CourseID     int64
	AssignmentID int64
	Include      []string
}

// Validate validates the options
func (o *PeerReviewsListOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return ValidateRequired("assignment-id", o.AssignmentID)
}

// PeerReviewsCreateOptions contains options for creating a peer review
type PeerReviewsCreateOptions struct {
	CourseID     int64
	AssignmentID int64
	SubmissionID int64
	UserID       int64
}

// Validate validates the options
func (o *PeerReviewsCreateOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	if err := ValidateRequired("assignment-id", o.AssignmentID); err != nil {
		return err
	}
	if err := ValidateRequired("submission-id", o.SubmissionID); err != nil {
		return err
	}
	return ValidateRequired("user-id", o.UserID)
}

// PeerReviewsDeleteOptions contains options for deleting a peer review
type PeerReviewsDeleteOptions struct {
	CourseID     int64
	AssignmentID int64
	SubmissionID int64
	UserID       int64
	Force        bool
}

// Validate validates the options
func (o *PeerReviewsDeleteOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	if err := ValidateRequired("assignment-id", o.AssignmentID); err != nil {
		return err
	}
	if err := ValidateRequired("submission-id", o.SubmissionID); err != nil {
		return err
	}
	return ValidateRequired("user-id", o.UserID)
}
