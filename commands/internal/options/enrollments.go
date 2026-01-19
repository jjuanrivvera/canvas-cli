package options

import "fmt"

// EnrollmentsListOptions contains options for listing enrollments
type EnrollmentsListOptions struct {
	CourseID int64
	UserID   int64
	Type     []string
	State    []string
	Include  []string
}

// Validate validates the options
func (o *EnrollmentsListOptions) Validate() error {
	// Must specify exactly one context
	contextsSpecified := 0
	if o.CourseID > 0 {
		contextsSpecified++
	}
	if o.UserID > 0 {
		contextsSpecified++
	}

	if contextsSpecified == 0 {
		return fmt.Errorf("must specify one of course-id or user-id")
	}
	if contextsSpecified > 1 {
		return fmt.Errorf("can only specify one of course-id or user-id")
	}

	return nil
}

// EnrollmentsGetOptions contains options for getting an enrollment
type EnrollmentsGetOptions struct {
	CourseID     int64
	EnrollmentID int64
}

// Validate validates the options
func (o *EnrollmentsGetOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.EnrollmentID <= 0 {
		return fmt.Errorf("enrollment-id is required and must be greater than 0")
	}
	return nil
}

// EnrollmentsCreateOptions contains options for creating an enrollment
type EnrollmentsCreateOptions struct {
	CourseID        int64
	UserID          int64
	EnrollmentType  string
	EnrollmentState string
	SectionID       int64
	Notify          bool
	Role            string
}

// Validate validates the options
func (o *EnrollmentsCreateOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}

// EnrollmentsConcludeOptions contains options for concluding an enrollment
type EnrollmentsConcludeOptions struct {
	CourseID     int64
	EnrollmentID int64
	Task         string
	Force        bool
}

// Validate validates the options
func (o *EnrollmentsConcludeOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.EnrollmentID <= 0 {
		return fmt.Errorf("enrollment-id is required and must be greater than 0")
	}
	// Validate task
	switch o.Task {
	case "conclude", "deactivate", "delete":
		// Valid
	default:
		return fmt.Errorf("invalid task: %s (use 'conclude', 'deactivate', or 'delete')", o.Task)
	}
	return nil
}

// EnrollmentsReactivateOptions contains options for reactivating an enrollment
type EnrollmentsReactivateOptions struct {
	CourseID     int64
	EnrollmentID int64
}

// Validate validates the options
func (o *EnrollmentsReactivateOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.EnrollmentID <= 0 {
		return fmt.Errorf("enrollment-id is required and must be greater than 0")
	}
	return nil
}

// EnrollmentsAcceptOptions contains options for accepting an enrollment
type EnrollmentsAcceptOptions struct {
	CourseID     int64
	EnrollmentID int64
}

// Validate validates the options
func (o *EnrollmentsAcceptOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.EnrollmentID <= 0 {
		return fmt.Errorf("enrollment-id is required and must be greater than 0")
	}
	return nil
}

// EnrollmentsRejectOptions contains options for rejecting an enrollment
type EnrollmentsRejectOptions struct {
	CourseID     int64
	EnrollmentID int64
}

// Validate validates the options
func (o *EnrollmentsRejectOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.EnrollmentID <= 0 {
		return fmt.Errorf("enrollment-id is required and must be greater than 0")
	}
	return nil
}
