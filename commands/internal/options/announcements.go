package options

// AnnouncementsListOptions contains options for listing announcements
type AnnouncementsListOptions struct {
	CourseID   int64
	StartDate  string
	EndDate    string
	ActiveOnly bool
	LatestOnly bool
	Include    []string
}

// Validate validates the options
func (o *AnnouncementsListOptions) Validate() error {
	return ValidateRequired("course-id", o.CourseID)
}

// AnnouncementsGetOptions contains options for getting an announcement
type AnnouncementsGetOptions struct {
	CourseID       int64
	AnnouncementID int64
}

// Validate validates the options
func (o *AnnouncementsGetOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return ValidateRequired("announcement-id", o.AnnouncementID)
}

// AnnouncementsCreateOptions contains options for creating an announcement
type AnnouncementsCreateOptions struct {
	CourseID  int64
	Title     string
	Message   string
	DelayedAt string
	Published bool
}

// Validate validates the options
func (o *AnnouncementsCreateOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return ValidateRequired("title", o.Title)
}

// AnnouncementsUpdateOptions contains options for updating an announcement
type AnnouncementsUpdateOptions struct {
	CourseID       int64
	AnnouncementID int64
	Title          string
	Message        string
	DelayedAt      string
	// Track which fields were set
	TitleSet     bool
	MessageSet   bool
	DelayedAtSet bool
}

// Validate validates the options
func (o *AnnouncementsUpdateOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return ValidateRequired("announcement-id", o.AnnouncementID)
}

// AnnouncementsDeleteOptions contains options for deleting an announcement
type AnnouncementsDeleteOptions struct {
	CourseID       int64
	AnnouncementID int64
	Force          bool
}

// Validate validates the options
func (o *AnnouncementsDeleteOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return ValidateRequired("announcement-id", o.AnnouncementID)
}
