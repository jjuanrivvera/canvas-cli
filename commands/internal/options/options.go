// Package options provides option structs for commands to eliminate global state.
// Each command should have its own options struct that encapsulates all flags.
package options

import "fmt"

// Validator is an interface for option structs that can validate themselves
type Validator interface {
	Validate() error
}

// ValidateRequired checks if a required field is set
func ValidateRequired(fieldName string, value interface{}) error {
	switch v := value.(type) {
	case string:
		if v == "" {
			return fmt.Errorf("%s is required", fieldName)
		}
	case int64:
		if v == 0 {
			return fmt.Errorf("%s is required", fieldName)
		}
	case int:
		if v == 0 {
			return fmt.Errorf("%s is required", fieldName)
		}
	}
	return nil
}

// ErrInvalidValue returns an error for invalid field values
func ErrInvalidValue(fieldName string, value string, validOptions ...string) error {
	return fmt.Errorf("invalid %s: %s (valid options: %s)",
		fieldName, value, joinWithOr(validOptions))
}

// joinWithOr joins strings with "or" for the last element
func joinWithOr(items []string) string {
	if len(items) == 0 {
		return ""
	}
	if len(items) == 1 {
		return items[0]
	}
	if len(items) == 2 {
		return items[0] + " or " + items[1]
	}
	return joinWithCommas(items[:len(items)-1]) + " or " + items[len(items)-1]
}

// joinWithCommas joins strings with commas
func joinWithCommas(items []string) string {
	result := ""
	for i, item := range items {
		if i > 0 {
			result += ", "
		}
		result += item
	}
	return result
}
