package options

import "fmt"

// AssignmentsListOptions contains options for listing assignments
type AssignmentsListOptions struct {
	CourseID   int64
	SearchTerm string
	Bucket     string
	OrderBy    string
	Include    []string
}

// Validate validates the options
func (o *AssignmentsListOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	return nil
}

// AssignmentsGetOptions contains options for getting an assignment
type AssignmentsGetOptions struct {
	CourseID     int64
	AssignmentID int64
	Include      []string
}

// Validate validates the options
func (o *AssignmentsGetOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.AssignmentID <= 0 {
		return fmt.Errorf("assignment-id is required and must be greater than 0")
	}
	return nil
}

// AssignmentsCreateOptions contains options for creating an assignment
type AssignmentsCreateOptions struct {
	CourseID        int64
	Name            string
	Points          float64
	GradingType     string
	DueAt           string
	UnlockAt        string
	LockAt          string
	Description     string
	Published       bool
	SubmissionTypes []string
	GroupID         int64
	Position        int
	JSONFile        string
	Stdin           bool
}

// Validate validates the options
func (o *AssignmentsCreateOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	return nil
}

// AssignmentsUpdateOptions contains options for updating an assignment
type AssignmentsUpdateOptions struct {
	CourseID        int64
	AssignmentID    int64
	Name            string
	Points          float64
	GradingType     string
	DueAt           string
	UnlockAt        string
	LockAt          string
	Description     string
	Published       bool
	SubmissionTypes []string
	GroupID         int64
	Position        int
	JSONFile        string
	Stdin           bool
}

// Validate validates the options
func (o *AssignmentsUpdateOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.AssignmentID <= 0 {
		return fmt.Errorf("assignment-id is required and must be greater than 0")
	}
	return nil
}

// AssignmentsDeleteOptions contains options for deleting an assignment
type AssignmentsDeleteOptions struct {
	CourseID     int64
	AssignmentID int64
	Force        bool
}

// Validate validates the options
func (o *AssignmentsDeleteOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.AssignmentID <= 0 {
		return fmt.Errorf("assignment-id is required and must be greater than 0")
	}
	return nil
}
