package options

import "fmt"

// UsersListOptions contains options for listing users
type UsersListOptions struct {
	AccountID       int64
	CourseID        int64
	SearchTerm      string
	EnrollmentType  string
	EnrollmentState string
	Include         []string
}

// Validate validates the options
func (o *UsersListOptions) Validate() error {
	if o.AccountID > 0 && o.CourseID > 0 {
		return fmt.Errorf("can only specify one of --account-id or --course-id")
	}
	return nil
}

// UsersGetOptions contains options for getting a user
type UsersGetOptions struct {
	UserID  int64
	Include []string
}

// Validate validates the options
func (o *UsersGetOptions) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}

// UsersMeOptions contains options for getting current user
type UsersMeOptions struct {
	// No options needed
}

// Validate validates the options
func (o *UsersMeOptions) Validate() error {
	return nil
}

// UsersSearchOptions contains options for searching users
type UsersSearchOptions struct {
	SearchTerm string
}

// Validate validates the options
func (o *UsersSearchOptions) Validate() error {
	if o.SearchTerm == "" {
		return fmt.Errorf("search-term is required")
	}
	return nil
}

// UsersCreateOptions contains options for creating a user
type UsersCreateOptions struct {
	AccountID        int64
	Name             string
	ShortName        string
	SortableName     string
	Email            string
	LoginID          string
	Password         string
	SISUserID        string
	TimeZone         string
	Locale           string
	SkipRegistration bool
	SkipConfirmation bool
	JSONFile         string
	Stdin            bool
}

// Validate validates the options
func (o *UsersCreateOptions) Validate() error {
	if o.AccountID <= 0 {
		return fmt.Errorf("account-id is required and must be greater than 0")
	}
	return nil
}

// UsersUpdateOptions contains options for updating a user
type UsersUpdateOptions struct {
	UserID       int64
	Name         string
	ShortName    string
	SortableName string
	Email        string
	TimeZone     string
	Locale       string
	JSONFile     string
	Stdin        bool
}

// Validate validates the options
func (o *UsersUpdateOptions) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}
