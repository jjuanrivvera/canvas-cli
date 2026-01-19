package options

// ContentMigrationsListOptions contains options for listing content migrations
type ContentMigrationsListOptions struct {
	CourseID int64
}

// Validate validates the options
func (o *ContentMigrationsListOptions) Validate() error {
	return ValidateRequired("course-id", o.CourseID)
}

// ContentMigrationsGetOptions contains options for getting a content migration
type ContentMigrationsGetOptions struct {
	CourseID    int64
	MigrationID int64
}

// Validate validates the options
func (o *ContentMigrationsGetOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return ValidateRequired("migration-id", o.MigrationID)
}

// ContentMigrationsCreateOptions contains options for creating a content migration
type ContentMigrationsCreateOptions struct {
	CourseID       int64
	Type           string
	SourceCourseID int64
	File           string
	FileURL        string
	FolderID       int64
	Selective      bool
	CopyOptions    string // JSON string
	DateShift      string // JSON string
	// Track which fields were set
	SourceCourseIDSet bool
	FolderIDSet       bool
	SelectiveSet      bool
}

// Validate validates the options
func (o *ContentMigrationsCreateOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return ValidateRequired("type", o.Type)
}

// ContentMigrationsMigratorsOptions contains options for listing available migration types
type ContentMigrationsMigratorsOptions struct {
	CourseID int64
}

// Validate validates the options
func (o *ContentMigrationsMigratorsOptions) Validate() error {
	return ValidateRequired("course-id", o.CourseID)
}

// ContentMigrationsContentOptions contains options for listing migration content
type ContentMigrationsContentOptions struct {
	CourseID    int64
	MigrationID int64
	ContentType string
}

// Validate validates the options
func (o *ContentMigrationsContentOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return ValidateRequired("migration-id", o.MigrationID)
}

// ContentMigrationsIssuesOptions contains options for listing migration issues
type ContentMigrationsIssuesOptions struct {
	CourseID    int64
	MigrationID int64
}

// Validate validates the options
func (o *ContentMigrationsIssuesOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return ValidateRequired("migration-id", o.MigrationID)
}
