package options

import "fmt"

// ContextSetOptions contains options for context set command
type ContextSetOptions struct {
	Type string
	ID   int64
}

// Validate validates the options
func (o *ContextSetOptions) Validate() error {
	if o.Type == "" {
		return fmt.Errorf("context type is required")
	}
	if o.ID <= 0 {
		return fmt.Errorf("ID must be a positive number")
	}
	validTypes := []string{"course", "course_id", "course-id", "assignment", "assignment_id", "assignment-id", "user", "user_id", "user-id", "account", "account_id", "account-id"}
	isValid := false
	for _, t := range validTypes {
		if o.Type == t {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("unknown context type %q. Valid types: course, assignment, user, account", o.Type)
	}
	return nil
}

// ContextShowOptions contains options for context show command
type ContextShowOptions struct {
	// No options needed for show
}

// Validate validates the options
func (o *ContextShowOptions) Validate() error {
	return nil
}

// ContextClearOptions contains options for context clear command
type ContextClearOptions struct {
	Type string // Optional: if empty, clears all context
}

// Validate validates the options
func (o *ContextClearOptions) Validate() error {
	if o.Type == "" {
		return nil // Empty is valid, means clear all
	}
	validTypes := []string{"course", "course_id", "course-id", "assignment", "assignment_id", "assignment-id", "user", "user_id", "user-id", "account", "account_id", "account-id"}
	isValid := false
	for _, t := range validTypes {
		if o.Type == t {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("unknown context type %q. Valid types: course, assignment, user, account", o.Type)
	}
	return nil
}
