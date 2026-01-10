package api

import (
	"errors"
	"fmt"
)

// Common validation errors
var (
	ErrInvalidCourseID     = errors.New("course_id must be a positive integer")
	ErrInvalidModuleID     = errors.New("module_id must be a positive integer")
	ErrInvalidItemID       = errors.New("item_id must be a positive integer")
	ErrInvalidAccountID    = errors.New("account_id must be a positive integer")
	ErrInvalidUserID       = errors.New("user_id must be a positive integer")
	ErrInvalidAssignmentID = errors.New("assignment_id must be a positive integer")
	ErrInvalidFileID       = errors.New("file_id must be a positive integer")
	ErrInvalidFolderID     = errors.New("folder_id must be a positive integer")
	ErrMissingTitle        = errors.New("title is required")
	ErrMissingName         = errors.New("name is required")
	ErrMissingType         = errors.New("type is required")
	ErrMissingURLOrID      = errors.New("url or id is required")
	ErrNilParams           = errors.New("params cannot be nil")
)

// ValidatePositiveID validates that an ID is positive
func ValidatePositiveID(id int64, name string) error {
	if id <= 0 {
		return fmt.Errorf("%s must be a positive integer", name)
	}
	return nil
}

// ValidateNonEmpty validates that a string is not empty
func ValidateNonEmpty(value, name string) error {
	if value == "" {
		return fmt.Errorf("%s is required", name)
	}
	return nil
}

// ValidateNotNil validates that a pointer is not nil
func ValidateNotNil(ptr interface{}, name string) error {
	if ptr == nil {
		return fmt.Errorf("%s cannot be nil", name)
	}
	return nil
}
