package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// UsersService handles user-related API calls
type UsersService struct {
	client *Client
}

// NewUsersService creates a new users service
func NewUsersService(client *Client) *UsersService {
	return &UsersService{client: client}
}

// ListUsersOptions holds options for listing users
type ListUsersOptions struct {
	SearchTerm      string   // Search by name, login ID, or email
	EnrollmentType  string   // Filter by enrollment type
	EnrollmentState string   // Filter by enrollment state
	Include         []string // Additional data to include
	Page            int
	PerPage         int
}

// GetCurrentUser retrieves the current authenticated user
func (s *UsersService) GetCurrentUser(ctx context.Context) (*User, error) {
	path := "/api/v1/users/self"

	var user User
	if err := s.client.GetJSON(ctx, path, &user); err != nil {
		return nil, err
	}

	return NormalizeUser(&user), nil
}

// Get retrieves a single user by ID
func (s *UsersService) Get(ctx context.Context, userID int64, include []string) (*User, error) {
	path := fmt.Sprintf("/api/v1/users/%d", userID)

	if len(include) > 0 {
		query := url.Values{}
		for _, inc := range include {
			query.Add("include[]", inc)
		}
		path += "?" + query.Encode()
	}

	var user User
	if err := s.client.GetJSON(ctx, path, &user); err != nil {
		return nil, err
	}

	return NormalizeUser(&user), nil
}

// List retrieves users for an account
func (s *UsersService) List(ctx context.Context, accountID int64, opts *ListUsersOptions) ([]User, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/users", accountID)

	if opts != nil {
		query := url.Values{}

		if opts.SearchTerm != "" {
			query.Add("search_term", opts.SearchTerm)
		}

		if opts.EnrollmentType != "" {
			query.Add("enrollment_type", opts.EnrollmentType)
		}

		if opts.EnrollmentState != "" {
			query.Add("enrollment_state", opts.EnrollmentState)
		}

		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}

		if opts.Page > 0 {
			query.Add("page", strconv.Itoa(opts.Page))
		}

		if opts.PerPage > 0 {
			query.Add("per_page", strconv.Itoa(opts.PerPage))
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var users []User
	if err := s.client.GetAllPages(ctx, path, &users); err != nil {
		return nil, err
	}

	return NormalizeUsers(users), nil
}

// CreateUserParams holds parameters for creating a user
type CreateUserParams struct {
	Name               string
	ShortName          string
	SortableName       string
	TimeZone           string
	Locale             string
	TermsOfUse         bool
	SkipRegistration   bool
	ForceValidations   bool
	EnableSISReactivation bool
	// Pseudonym (login) information
	UniqueID           string
	SISUserID          string
	IntegrationID      string
	AuthenticationProviderID string
	Password           string
	// Communication channel (email)
	Email              string
	SkipConfirmation   bool
}

// Create creates a new user
func (s *UsersService) Create(ctx context.Context, accountID int64, params *CreateUserParams) (*User, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/users", accountID)

	body := map[string]interface{}{
		"user": make(map[string]interface{}),
	}

	userData := body["user"].(map[string]interface{})

	if params.Name != "" {
		userData["name"] = params.Name
	}
	if params.ShortName != "" {
		userData["short_name"] = params.ShortName
	}
	if params.SortableName != "" {
		userData["sortable_name"] = params.SortableName
	}
	if params.TimeZone != "" {
		userData["time_zone"] = params.TimeZone
	}
	if params.Locale != "" {
		userData["locale"] = params.Locale
	}
	if params.TermsOfUse {
		userData["terms_of_use"] = true
	}
	if params.SkipRegistration {
		userData["skip_registration"] = true
	}

	// Pseudonym (login credentials)
	if params.UniqueID != "" || params.Password != "" {
		pseudonym := make(map[string]interface{})

		if params.UniqueID != "" {
			pseudonym["unique_id"] = params.UniqueID
		}
		if params.SISUserID != "" {
			pseudonym["sis_user_id"] = params.SISUserID
		}
		if params.IntegrationID != "" {
			pseudonym["integration_id"] = params.IntegrationID
		}
		if params.AuthenticationProviderID != "" {
			pseudonym["authentication_provider_id"] = params.AuthenticationProviderID
		}
		if params.Password != "" {
			pseudonym["password"] = params.Password
		}

		body["pseudonym"] = pseudonym
	}

	// Communication channel (email)
	if params.Email != "" {
		communication := map[string]interface{}{
			"type":    "email",
			"address": params.Email,
		}
		if params.SkipConfirmation {
			communication["skip_confirmation"] = true
		}
		body["communication_channel"] = communication
	}

	// Additional flags
	if params.ForceValidations {
		body["force_validations"] = true
	}
	if params.EnableSISReactivation {
		body["enable_sis_reactivation"] = true
	}

	var user User
	if err := s.client.PostJSON(ctx, path, body, &user); err != nil {
		return nil, err
	}

	return NormalizeUser(&user), nil
}

// UpdateUserParams holds parameters for updating a user
type UpdateUserParams struct {
	Name         string
	ShortName    string
	SortableName string
	TimeZone     string
	Locale       string
	Email        string
	Avatar       *AvatarParams
}

// AvatarParams holds avatar upload parameters
type AvatarParams struct {
	Token string // Avatar upload token
	URL   string // Avatar URL
}

// Update updates an existing user
func (s *UsersService) Update(ctx context.Context, userID int64, params *UpdateUserParams) (*User, error) {
	path := fmt.Sprintf("/api/v1/users/%d", userID)

	body := map[string]interface{}{
		"user": make(map[string]interface{}),
	}

	userData := body["user"].(map[string]interface{})

	if params.Name != "" {
		userData["name"] = params.Name
	}
	if params.ShortName != "" {
		userData["short_name"] = params.ShortName
	}
	if params.SortableName != "" {
		userData["sortable_name"] = params.SortableName
	}
	if params.TimeZone != "" {
		userData["time_zone"] = params.TimeZone
	}
	if params.Locale != "" {
		userData["locale"] = params.Locale
	}
	if params.Email != "" {
		userData["email"] = params.Email
	}

	if params.Avatar != nil {
		avatar := make(map[string]interface{})
		if params.Avatar.Token != "" {
			avatar["token"] = params.Avatar.Token
		}
		if params.Avatar.URL != "" {
			avatar["url"] = params.Avatar.URL
		}
		userData["avatar"] = avatar
	}

	var user User
	if err := s.client.PutJSON(ctx, path, body, &user); err != nil {
		return nil, err
	}

	return NormalizeUser(&user), nil
}

// ListCourseUsers retrieves users enrolled in a course
func (s *UsersService) ListCourseUsers(ctx context.Context, courseID int64, opts *ListUsersOptions) ([]User, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/users", courseID)

	if opts != nil {
		query := url.Values{}

		if opts.SearchTerm != "" {
			query.Add("search_term", opts.SearchTerm)
		}

		if opts.EnrollmentType != "" {
			query.Add("enrollment_type[]", opts.EnrollmentType)
		}

		if opts.EnrollmentState != "" {
			query.Add("enrollment_state[]", opts.EnrollmentState)
		}

		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}

		if opts.Page > 0 {
			query.Add("page", strconv.Itoa(opts.Page))
		}

		if opts.PerPage > 0 {
			query.Add("per_page", strconv.Itoa(opts.PerPage))
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var users []User
	if err := s.client.GetAllPages(ctx, path, &users); err != nil {
		return nil, err
	}

	return NormalizeUsers(users), nil
}

// Search searches for users across the entire Canvas instance
func (s *UsersService) Search(ctx context.Context, searchTerm string) ([]User, error) {
	path := "/api/v1/search/recipients"

	query := url.Values{}
	query.Add("search", searchTerm)
	query.Add("type", "user")
	path += "?" + query.Encode()

	var users []User
	if err := s.client.GetJSON(ctx, path, &users); err != nil {
		return nil, err
	}

	return NormalizeUsers(users), nil
}

// Delete deletes/suspends a user
func (s *UsersService) Delete(ctx context.Context, accountID, userID int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/users/%d", accountID, userID)

	_, err := s.client.Delete(ctx, path)
	return err
}
