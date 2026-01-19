package options

// SectionsListOptions contains options for listing sections
type SectionsListOptions struct {
	CourseID int64
	Include  []string
}

// Validate validates the options
func (o *SectionsListOptions) Validate() error {
	return ValidateRequired("course-id", o.CourseID)
}

// SectionsGetOptions contains options for getting a section
type SectionsGetOptions struct {
	SectionID int64
	Include   []string
}

// Validate validates the options
func (o *SectionsGetOptions) Validate() error {
	return ValidateRequired("section-id", o.SectionID)
}

// SectionsCreateOptions contains options for creating a section
type SectionsCreateOptions struct {
	CourseID      int64
	Name          string
	SISSectionID  string
	IntegrationID string
	StartAt       string
	EndAt         string
	RestrictDates bool
}

// Validate validates the options
func (o *SectionsCreateOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return ValidateRequired("name", o.Name)
}

// SectionsUpdateOptions contains options for updating a section
type SectionsUpdateOptions struct {
	SectionID             int64
	Name                  string
	SISSectionID          string
	IntegrationID         string
	StartAt               string
	EndAt                 string
	RestrictDates         bool
	OverrideSISStickiness bool
	// Track which fields were set
	NameSet          bool
	SISSectionIDSet  bool
	IntegrationIDSet bool
	StartAtSet       bool
	EndAtSet         bool
	RestrictDatesSet bool
}

// Validate validates the options
func (o *SectionsUpdateOptions) Validate() error {
	return ValidateRequired("section-id", o.SectionID)
}

// SectionsDeleteOptions contains options for deleting a section
type SectionsDeleteOptions struct {
	SectionID int64
	Force     bool
}

// Validate validates the options
func (o *SectionsDeleteOptions) Validate() error {
	return ValidateRequired("section-id", o.SectionID)
}

// SectionsCrosslistOptions contains options for crosslisting a section
type SectionsCrosslistOptions struct {
	SectionID             int64
	NewCourseID           int64
	OverrideSISStickiness bool
}

// Validate validates the options
func (o *SectionsCrosslistOptions) Validate() error {
	if err := ValidateRequired("section-id", o.SectionID); err != nil {
		return err
	}
	return ValidateRequired("new-course-id", o.NewCourseID)
}

// SectionsUncrosslistOptions contains options for uncrosslisting a section
type SectionsUncrosslistOptions struct {
	SectionID             int64
	OverrideSISStickiness bool
}

// Validate validates the options
func (o *SectionsUncrosslistOptions) Validate() error {
	return ValidateRequired("section-id", o.SectionID)
}
