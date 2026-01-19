package options

// CalendarListOptions contains options for listing calendar events
type CalendarListOptions struct {
	CourseID       int64
	UserID         int64
	Type           string
	StartDate      string
	EndDate        string
	Undated        bool
	AllEvents      bool
	ContextCodes   []string
	ImportantDates bool
}

// Validate validates the options
func (o *CalendarListOptions) Validate() error {
	// No required fields
	return nil
}

// CalendarGetOptions contains options for getting a calendar event
type CalendarGetOptions struct {
	EventID int64
}

// Validate validates the options
func (o *CalendarGetOptions) Validate() error {
	return ValidateRequired("event-id", o.EventID)
}

// CalendarCreateOptions contains options for creating a calendar event
type CalendarCreateOptions struct {
	ContextCode     string
	Title           string
	Description     string
	StartAt         string
	EndAt           string
	LocationName    string
	LocationAddress string
	AllDay          bool
}

// Validate validates the options
func (o *CalendarCreateOptions) Validate() error {
	if err := ValidateRequired("context-code", o.ContextCode); err != nil {
		return err
	}
	return ValidateRequired("title", o.Title)
}

// CalendarUpdateOptions contains options for updating a calendar event
type CalendarUpdateOptions struct {
	EventID         int64
	Title           string
	Description     string
	StartAt         string
	EndAt           string
	LocationName    string
	LocationAddress string
	AllDay          bool
	// Track which fields were set
	TitleSet           bool
	DescriptionSet     bool
	StartAtSet         bool
	EndAtSet           bool
	LocationNameSet    bool
	LocationAddressSet bool
	AllDaySet          bool
}

// Validate validates the options
func (o *CalendarUpdateOptions) Validate() error {
	return ValidateRequired("event-id", o.EventID)
}

// CalendarDeleteOptions contains options for deleting a calendar event
type CalendarDeleteOptions struct {
	EventID      int64
	Which        string
	CancelReason string
	Force        bool
}

// Validate validates the options
func (o *CalendarDeleteOptions) Validate() error {
	return ValidateRequired("event-id", o.EventID)
}

// CalendarReserveOptions contains options for reserving a calendar event
type CalendarReserveOptions struct {
	EventID        int64
	ParticipantID  string
	Comments       string
	CancelExisting bool
}

// Validate validates the options
func (o *CalendarReserveOptions) Validate() error {
	return ValidateRequired("event-id", o.EventID)
}
