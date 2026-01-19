package options

// BlueprintGetOptions contains options for getting a blueprint
type BlueprintGetOptions struct {
	CourseID   int64
	TemplateID string
}

// Validate validates the options
func (o *BlueprintGetOptions) Validate() error {
	return ValidateRequired("course-id", o.CourseID)
}

// BlueprintAssociationsListOptions contains options for listing blueprint associations
type BlueprintAssociationsListOptions struct {
	CourseID   int64
	TemplateID string
}

// Validate validates the options
func (o *BlueprintAssociationsListOptions) Validate() error {
	return ValidateRequired("course-id", o.CourseID)
}

// BlueprintAssociationsAddOptions contains options for adding blueprint associations
type BlueprintAssociationsAddOptions struct {
	CourseID     int64
	TemplateID   string
	CourseIDsStr string
}

// Validate validates the options
func (o *BlueprintAssociationsAddOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return ValidateRequired("course-ids-to-add", o.CourseIDsStr)
}

// BlueprintAssociationsRemoveOptions contains options for removing blueprint associations
type BlueprintAssociationsRemoveOptions struct {
	CourseID     int64
	TemplateID   string
	CourseIDsStr string
}

// Validate validates the options
func (o *BlueprintAssociationsRemoveOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return ValidateRequired("course-ids-to-remove", o.CourseIDsStr)
}

// BlueprintSyncOptions contains options for syncing a blueprint
type BlueprintSyncOptions struct {
	CourseID     int64
	TemplateID   string
	Comment      string
	Notify       bool
	CopySettings bool
	Publish      bool
	// Track which fields were set
	NotifySet       bool
	CopySettingsSet bool
	PublishSet      bool
}

// Validate validates the options
func (o *BlueprintSyncOptions) Validate() error {
	return ValidateRequired("course-id", o.CourseID)
}

// BlueprintChangesOptions contains options for listing blueprint changes
type BlueprintChangesOptions struct {
	CourseID   int64
	TemplateID string
}

// Validate validates the options
func (o *BlueprintChangesOptions) Validate() error {
	return ValidateRequired("course-id", o.CourseID)
}

// BlueprintMigrationsListOptions contains options for listing blueprint migrations
type BlueprintMigrationsListOptions struct {
	CourseID   int64
	TemplateID string
}

// Validate validates the options
func (o *BlueprintMigrationsListOptions) Validate() error {
	return ValidateRequired("course-id", o.CourseID)
}

// BlueprintMigrationsGetOptions contains options for getting a blueprint migration
type BlueprintMigrationsGetOptions struct {
	CourseID    int64
	TemplateID  string
	MigrationID int64
	Include     []string
}

// Validate validates the options
func (o *BlueprintMigrationsGetOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return ValidateRequired("migration-id", o.MigrationID)
}
