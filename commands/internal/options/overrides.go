package options

import "fmt"

// OverridesListOptions contains options for listing assignment overrides
type OverridesListOptions struct {
	CourseID     int64
	AssignmentID int64
}

// Validate validates the options
func (o *OverridesListOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return ValidateRequired("assignment-id", o.AssignmentID)
}

// OverridesGetOptions contains options for getting an assignment override
type OverridesGetOptions struct {
	CourseID     int64
	AssignmentID int64
	OverrideID   int64
}

// Validate validates the options
func (o *OverridesGetOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	if err := ValidateRequired("assignment-id", o.AssignmentID); err != nil {
		return err
	}
	return ValidateRequired("override-id", o.OverrideID)
}

// OverridesCreateOptions contains options for creating an assignment override
type OverridesCreateOptions struct {
	CourseID     int64
	AssignmentID int64
	StudentIDs   string // Comma-separated
	SectionID    int64
	GroupID      int64
	Title        string
	DueAt        string
	UnlockAt     string
	LockAt       string
}

// Validate validates the options
func (o *OverridesCreateOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	if err := ValidateRequired("assignment-id", o.AssignmentID); err != nil {
		return err
	}

	// Validate that exactly one target is specified
	hasStudents := o.StudentIDs != ""
	hasSection := o.SectionID > 0
	hasGroup := o.GroupID > 0

	if !hasStudents && !hasSection && !hasGroup {
		return fmt.Errorf("must specify one of --student-ids, --section-id, or --group-id")
	}

	targetsCount := 0
	if hasStudents {
		targetsCount++
	}
	if hasSection {
		targetsCount++
	}
	if hasGroup {
		targetsCount++
	}
	if targetsCount > 1 {
		return fmt.Errorf("can only specify one of --student-ids, --section-id, or --group-id")
	}

	return nil
}

// OverridesUpdateOptions contains options for updating an assignment override
type OverridesUpdateOptions struct {
	CourseID     int64
	AssignmentID int64
	OverrideID   int64
	StudentIDs   string
	Title        string
	DueAt        string
	UnlockAt     string
	LockAt       string
	// Track which fields were set
	StudentIDsSet bool
	TitleSet      bool
	DueAtSet      bool
	UnlockAtSet   bool
	LockAtSet     bool
}

// Validate validates the options
func (o *OverridesUpdateOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	if err := ValidateRequired("assignment-id", o.AssignmentID); err != nil {
		return err
	}
	return ValidateRequired("override-id", o.OverrideID)
}

// OverridesDeleteOptions contains options for deleting an assignment override
type OverridesDeleteOptions struct {
	CourseID     int64
	AssignmentID int64
	OverrideID   int64
	Force        bool
}

// Validate validates the options
func (o *OverridesDeleteOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	if err := ValidateRequired("assignment-id", o.AssignmentID); err != nil {
		return err
	}
	return ValidateRequired("override-id", o.OverrideID)
}
