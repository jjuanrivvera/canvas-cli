package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// AccountsService handles account-related API operations
type AccountsService struct {
	client *Client
}

// NewAccountsService creates a new accounts service
func NewAccountsService(client *Client) *AccountsService {
	return &AccountsService{client: client}
}

// ListAccountsOptions holds options for listing accounts
type ListAccountsOptions struct {
	Include []string `url:"include[],omitempty"` // lti_guid, registration_settings, services
	PerPage int      `url:"per_page,omitempty"`
}

// List returns accounts the current user can view
// This typically returns accounts where the user has admin permissions
func (s *AccountsService) List(ctx context.Context, opts *ListAccountsOptions) ([]Account, error) {
	var accounts []Account

	path := "/api/v1/accounts"
	if opts != nil {
		query := url.Values{}
		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}
		if opts.PerPage > 0 {
			query.Add("per_page", strconv.Itoa(opts.PerPage))
		}
		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	err := s.client.GetAllPages(ctx, path, &accounts)
	if err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}

	return accounts, nil
}

// Get returns a single account by ID
func (s *AccountsService) Get(ctx context.Context, accountID int64) (*Account, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d", accountID)

	var account Account
	err := s.client.GetJSON(ctx, path, &account)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return &account, nil
}

// ListSubAccounts returns sub-accounts for a given account
type ListSubAccountsOptions struct {
	Recursive bool `url:"recursive,omitempty"` // If true, returns the entire account tree
	PerPage   int  `url:"per_page,omitempty"`
}

// ListSubAccounts returns sub-accounts for a given account
func (s *AccountsService) ListSubAccounts(ctx context.Context, accountID int64, opts *ListSubAccountsOptions) ([]Account, error) {
	var accounts []Account

	path := fmt.Sprintf("/api/v1/accounts/%d/sub_accounts", accountID)
	if opts != nil {
		query := url.Values{}
		if opts.Recursive {
			query.Add("recursive", "true")
		}
		if opts.PerPage > 0 {
			query.Add("per_page", strconv.Itoa(opts.PerPage))
		}
		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	err := s.client.GetAllPages(ctx, path, &accounts)
	if err != nil {
		return nil, fmt.Errorf("failed to list sub-accounts: %w", err)
	}

	return accounts, nil
}

// ListAccountCoursesOptions holds options for listing courses in an account
type ListAccountCoursesOptions struct {
	// Filter options
	WithEnrollments     bool     `url:"with_enrollments,omitempty"`     // Only courses with at least one enrollment
	EnrollmentType      []string `url:"enrollment_type[],omitempty"`    // teacher, student, ta, observer, designer
	Published           bool     `url:"published,omitempty"`            // Only published courses
	Completed           bool     `url:"completed,omitempty"`            // Only completed courses
	Blueprint           bool     `url:"blueprint,omitempty"`            // Only blueprint courses
	BlueprintAssociated bool     `url:"blueprint_associated,omitempty"` // Only courses associated with blueprints

	// Search options
	SearchTerm    string  `url:"search_term,omitempty"`      // Search by name/code
	ByTeachers    []int64 `url:"by_teachers[],omitempty"`    // Filter by teacher IDs
	BySubaccounts []int64 `url:"by_subaccounts[],omitempty"` // Filter by sub-account IDs

	// State options
	State            []string `url:"state[],omitempty"`            // created, claimed, available, completed, deleted, all
	EnrollmentTermID int64    `url:"enrollment_term_id,omitempty"` // Filter by term

	// Sorting options
	Sort  string `url:"sort,omitempty"`  // course_name, sis_course_id, teacher, account_name
	Order string `url:"order,omitempty"` // asc, desc

	// Include options
	Include []string `url:"include[],omitempty"` // syllabus_body, term, course_progress, etc.

	// Pagination
	PerPage int `url:"per_page,omitempty"`
}

// ListCourses returns courses for a given account
// This requires admin permissions on the account
func (s *AccountsService) ListCourses(ctx context.Context, accountID int64, opts *ListAccountCoursesOptions) ([]Course, error) {
	var courses []Course

	path := fmt.Sprintf("/api/v1/accounts/%d/courses", accountID)
	if opts != nil {
		query := url.Values{}
		if opts.SearchTerm != "" {
			query.Add("search_term", opts.SearchTerm)
		}
		if opts.WithEnrollments {
			query.Add("with_enrollments", "true")
		}
		if opts.Published {
			query.Add("published", "true")
		}
		if opts.Completed {
			query.Add("completed", "true")
		}
		if opts.Blueprint {
			query.Add("blueprint", "true")
		}
		if opts.BlueprintAssociated {
			query.Add("blueprint_associated", "true")
		}
		if opts.EnrollmentTermID > 0 {
			query.Add("enrollment_term_id", strconv.FormatInt(opts.EnrollmentTermID, 10))
		}
		for _, state := range opts.State {
			query.Add("state[]", state)
		}
		for _, et := range opts.EnrollmentType {
			query.Add("enrollment_type[]", et)
		}
		for _, teacher := range opts.ByTeachers {
			query.Add("by_teachers[]", strconv.FormatInt(teacher, 10))
		}
		for _, subaccount := range opts.BySubaccounts {
			query.Add("by_subaccounts[]", strconv.FormatInt(subaccount, 10))
		}
		if opts.Sort != "" {
			query.Add("sort", opts.Sort)
		}
		if opts.Order != "" {
			query.Add("order", opts.Order)
		}
		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}
		if opts.PerPage > 0 {
			query.Add("per_page", strconv.Itoa(opts.PerPage))
		}
		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	err := s.client.GetAllPages(ctx, path, &courses)
	if err != nil {
		return nil, fmt.Errorf("failed to list account courses: %w", err)
	}

	return courses, nil
}

// ListAccountUsersOptions holds options for listing users in an account
type ListAccountUsersOptions struct {
	SearchTerm     string   `url:"search_term,omitempty"`     // Search by name/email
	EnrollmentType string   `url:"enrollment_type,omitempty"` // Filter by enrollment type
	Sort           string   `url:"sort,omitempty"`            // username, email, sis_id, last_login
	Order          string   `url:"order,omitempty"`           // asc, desc
	Include        []string `url:"include[],omitempty"`       // avatar_url, email, last_login, etc.
	PerPage        int      `url:"per_page,omitempty"`
}

// ListUsers returns users for a given account
// This requires admin permissions on the account
func (s *AccountsService) ListUsers(ctx context.Context, accountID int64, opts *ListAccountUsersOptions) ([]User, error) {
	var users []User

	path := fmt.Sprintf("/api/v1/accounts/%d/users", accountID)
	if opts != nil {
		query := url.Values{}
		if opts.SearchTerm != "" {
			query.Add("search_term", opts.SearchTerm)
		}
		if opts.EnrollmentType != "" {
			query.Add("enrollment_type", opts.EnrollmentType)
		}
		if opts.Sort != "" {
			query.Add("sort", opts.Sort)
		}
		if opts.Order != "" {
			query.Add("order", opts.Order)
		}
		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}
		if opts.PerPage > 0 {
			query.Add("per_page", strconv.Itoa(opts.PerPage))
		}
		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	err := s.client.GetAllPages(ctx, path, &users)
	if err != nil {
		return nil, fmt.Errorf("failed to list account users: %w", err)
	}

	return users, nil
}
