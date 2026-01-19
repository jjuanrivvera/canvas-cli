package options

import "fmt"

// PlannerItemsOptions contains options for listing planner items
type PlannerItemsOptions struct {
	CourseID     int64
	StartDate    string
	EndDate      string
	ContextCodes []string
	Filter       string
}

// Validate validates the options
func (o *PlannerItemsOptions) Validate() error {
	return nil
}

// PlannerNotesListOptions contains options for listing planner notes
type PlannerNotesListOptions struct {
	CourseID  int64
	StartDate string
	EndDate   string
}

// Validate validates the options
func (o *PlannerNotesListOptions) Validate() error {
	return nil
}

// PlannerNotesGetOptions contains options for getting a planner note
type PlannerNotesGetOptions struct {
	NoteID int64
}

// Validate validates the options
func (o *PlannerNotesGetOptions) Validate() error {
	return ValidateRequired("note-id", o.NoteID)
}

// PlannerNotesCreateOptions contains options for creating a planner note
type PlannerNotesCreateOptions struct {
	Title    string
	Details  string
	TodoDate string
	CourseID int64
}

// Validate validates the options
func (o *PlannerNotesCreateOptions) Validate() error {
	if o.Title == "" {
		return fmt.Errorf("title is required")
	}
	return nil
}

// PlannerNotesUpdateOptions contains options for updating a planner note
type PlannerNotesUpdateOptions struct {
	NoteID   int64
	Title    string
	Details  string
	TodoDate string
	CourseID int64
	// Track which fields were set
	TitleSet    bool
	DetailsSet  bool
	TodoDateSet bool
	CourseIDSet bool
}

// Validate validates the options
func (o *PlannerNotesUpdateOptions) Validate() error {
	return ValidateRequired("note-id", o.NoteID)
}

// PlannerNotesDeleteOptions contains options for deleting a planner note
type PlannerNotesDeleteOptions struct {
	NoteID int64
	Force  bool
}

// Validate validates the options
func (o *PlannerNotesDeleteOptions) Validate() error {
	return ValidateRequired("note-id", o.NoteID)
}

// PlannerCompleteOptions contains options for completing a planner item
type PlannerCompleteOptions struct {
	PlannableType string
	PlannableID   int64
}

// Validate validates the options
func (o *PlannerCompleteOptions) Validate() error {
	if o.PlannableType == "" {
		return fmt.Errorf("plannable type is required")
	}

	// Validate plannable type
	validTypes := []string{"Assignment", "Quiz", "DiscussionTopic", "WikiPage", "CalendarEvent", "PlannerNote", "Announcement"}
	isValidType := false
	for _, t := range validTypes {
		if t == o.PlannableType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return fmt.Errorf("invalid plannable type: %s\nValid types: Assignment, Quiz, DiscussionTopic, WikiPage, CalendarEvent, PlannerNote, Announcement", o.PlannableType)
	}

	return ValidateRequired("plannable-id", o.PlannableID)
}

// PlannerDismissOptions contains options for dismissing a planner item
type PlannerDismissOptions struct {
	PlannableType string
	PlannableID   int64
}

// Validate validates the options
func (o *PlannerDismissOptions) Validate() error {
	if o.PlannableType == "" {
		return fmt.Errorf("plannable type is required")
	}

	// Validate plannable type
	validTypes := []string{"Assignment", "Quiz", "DiscussionTopic", "WikiPage", "CalendarEvent", "PlannerNote", "Announcement"}
	isValidType := false
	for _, t := range validTypes {
		if t == o.PlannableType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return fmt.Errorf("invalid plannable type: %s\nValid types: Assignment, Quiz, DiscussionTopic, WikiPage, CalendarEvent, PlannerNote, Announcement", o.PlannableType)
	}

	return ValidateRequired("plannable-id", o.PlannableID)
}

// PlannerOverridesOptions contains options for listing planner overrides
type PlannerOverridesOptions struct {
	PlannableType string
}

// Validate validates the options
func (o *PlannerOverridesOptions) Validate() error {
	return nil
}
