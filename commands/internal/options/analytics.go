package options

// AnalyticsActivityOptions contains options for viewing course activity analytics
type AnalyticsActivityOptions struct {
	CourseID int64
}

// Validate validates the options
func (o *AnalyticsActivityOptions) Validate() error {
	return ValidateRequired("course-id", o.CourseID)
}

// AnalyticsAssignmentsOptions contains options for viewing assignment analytics
type AnalyticsAssignmentsOptions struct {
	CourseID int64
}

// Validate validates the options
func (o *AnalyticsAssignmentsOptions) Validate() error {
	return ValidateRequired("course-id", o.CourseID)
}

// AnalyticsStudentsOptions contains options for viewing student analytics
type AnalyticsStudentsOptions struct {
	CourseID   int64
	SortColumn string
}

// Validate validates the options
func (o *AnalyticsStudentsOptions) Validate() error {
	return ValidateRequired("course-id", o.CourseID)
}

// AnalyticsUserOptions contains options for viewing user analytics
type AnalyticsUserOptions struct {
	CourseID int64
	UserID   int64
	Type     string
}

// Validate validates the options
func (o *AnalyticsUserOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return ValidateRequired("user-id", o.UserID)
}

// AnalyticsDepartmentOptions contains options for viewing department analytics
type AnalyticsDepartmentOptions struct {
	AccountID int64
	Type      string
	TermID    int64
}

// Validate validates the options
func (o *AnalyticsDepartmentOptions) Validate() error {
	// Account ID is optional - will be resolved from config if not specified
	return nil
}
