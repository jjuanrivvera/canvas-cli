package options

import "fmt"

// GradesHistoryOptions contains options for gradebook history
type GradesHistoryOptions struct {
	CourseID  int64
	StartDate string
	EndDate   string
}

// Validate validates the options
func (o *GradesHistoryOptions) Validate() error {
	return ValidateRequired("course-id", o.CourseID)
}

// GradesFeedOptions contains options for gradebook feed
type GradesFeedOptions struct {
	CourseID     int64
	UserID       int64
	AssignmentID int64
	StartDate    string
	EndDate      string
}

// Validate validates the options
func (o *GradesFeedOptions) Validate() error {
	return ValidateRequired("course-id", o.CourseID)
}

// GradesColumnsListOptions contains options for listing custom columns
type GradesColumnsListOptions struct {
	CourseID      int64
	IncludeHidden bool
}

// Validate validates the options
func (o *GradesColumnsListOptions) Validate() error {
	return ValidateRequired("course-id", o.CourseID)
}

// GradesColumnsGetOptions contains options for getting a custom column
type GradesColumnsGetOptions struct {
	CourseID int64
	ColumnID int64
}

// Validate validates the options
func (o *GradesColumnsGetOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return ValidateRequired("column-id", o.ColumnID)
}

// GradesColumnsCreateOptions contains options for creating a custom column
type GradesColumnsCreateOptions struct {
	CourseID     int64
	Title        string
	Position     int
	Hidden       bool
	TeacherNotes bool
	ReadOnly     bool
}

// Validate validates the options
func (o *GradesColumnsCreateOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	if o.Title == "" {
		return fmt.Errorf("title is required")
	}
	return nil
}

// GradesColumnsUpdateOptions contains options for updating a custom column
type GradesColumnsUpdateOptions struct {
	CourseID     int64
	ColumnID     int64
	Title        string
	Position     int
	Hidden       bool
	TeacherNotes bool
	ReadOnly     bool
	// Track which fields were set
	TitleSet        bool
	PositionSet     bool
	HiddenSet       bool
	TeacherNotesSet bool
	ReadOnlySet     bool
}

// Validate validates the options
func (o *GradesColumnsUpdateOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return ValidateRequired("column-id", o.ColumnID)
}

// GradesColumnsDeleteOptions contains options for deleting a custom column
type GradesColumnsDeleteOptions struct {
	CourseID int64
	ColumnID int64
	Force    bool
}

// Validate validates the options
func (o *GradesColumnsDeleteOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return ValidateRequired("column-id", o.ColumnID)
}

// GradesColumnsDataListOptions contains options for listing custom column data
type GradesColumnsDataListOptions struct {
	CourseID int64
	ColumnID int64
}

// Validate validates the options
func (o *GradesColumnsDataListOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return ValidateRequired("column-id", o.ColumnID)
}

// GradesColumnsDataSetOptions contains options for setting custom column data
type GradesColumnsDataSetOptions struct {
	CourseID int64
	ColumnID int64
	UserID   int64
	Content  string
}

// Validate validates the options
func (o *GradesColumnsDataSetOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	if err := ValidateRequired("column-id", o.ColumnID); err != nil {
		return err
	}
	if err := ValidateRequired("user-id", o.UserID); err != nil {
		return err
	}
	if o.Content == "" {
		return fmt.Errorf("content is required")
	}
	return nil
}
