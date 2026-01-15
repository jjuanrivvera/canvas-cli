package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var (
	usersAccountID       int64
	usersCourseID        int64
	usersSearchTerm      string
	usersInclude         []string
	usersEnrollmentType  string
	usersEnrollmentState string

	// Create/Update flags
	userName             string
	userShortName        string
	userSortableName     string
	userEmail            string
	userLoginID          string
	userPassword         string
	userSISUserID        string
	userTimeZone         string
	userLocale           string
	userSkipRegistration bool
	userSkipConfirmation bool
	userJSONFile         string
	userStdin            bool
)

// usersCmd represents the users command group
var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Manage Canvas users",
	Long: `Manage Canvas users including listing, viewing, searching, and managing user information.

Examples:
  canvas users list --account-id 1
  canvas users get 123
  canvas users search "john"
  canvas users me`,
}

// usersListCmd represents the users list command
var usersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List users in an account or course",
	Long: `List users in a Canvas account or course.

Specify --account-id or --course-id, or uses default account if configured.

WARNING: Account-level user lists can be very large. Use --limit or --search
to avoid long wait times.

Account context (admin):
  canvas users list --limit 100           # First 100 users (recommended)
  canvas users list --search "john"       # Search in default account

Course context:
  canvas users list --course-id 123       # All users enrolled in course 123
  canvas users list --course-id 123 --enrollment-type teacher

Examples:
  canvas users list --limit 50
  canvas users list --account-id 1 --limit 100
  canvas users list --course-id 123
  canvas users list --search "john"
  canvas users list --include email,enrollments`,
	RunE: runUsersList,
}

// usersGetCmd represents the users get command
var usersGetCmd = &cobra.Command{
	Use:   "get <user-id>",
	Short: "Get details of a specific user",
	Long: `Get details of a specific user by ID.

Examples:
  canvas users get 123
  canvas users get 123 --include email,enrollments,avatar_url`,
	Args: ExactArgsWithUsage(1, "user-id"),
	RunE: runUsersGet,
}

// usersMeCmd represents the users me command
var usersMeCmd = &cobra.Command{
	Use:   "me",
	Short: "Get details of the current authenticated user",
	Long: `Get details of the current authenticated user.

Examples:
  canvas users me`,
	RunE: runUsersMe,
}

// usersSearchCmd represents the users search command
var usersSearchCmd = &cobra.Command{
	Use:   "search <search-term>",
	Short: "Search for users",
	Long: `Search for users across the Canvas instance.

Examples:
  canvas users search "john doe"
  canvas users search "john@example.com"`,
	Args: ExactArgsWithUsage(1, "search-term"),
	RunE: runUsersSearch,
}

// usersCreateCmd represents the users create command
var usersCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new user (admin)",
	Long: `Create a new user in a Canvas account. Requires admin privileges.

You can provide user data via flags or JSON file/stdin.

Examples:
  # Using flags
  canvas users create --account-id 1 --name "John Doe" --email "john@example.com"
  canvas users create --account-id 1 --name "Jane" --login-id "jane123" --password "secret"

  # Using JSON file
  canvas users create --account-id 1 --json user.json

  # Using stdin
  echo '{"name":"John Doe","email":"john@example.com"}' | canvas users create --account-id 1 --stdin`,
	RunE: runUsersCreate,
}

// usersUpdateCmd represents the users update command
var usersUpdateCmd = &cobra.Command{
	Use:   "update <user-id>",
	Short: "Update an existing user",
	Long: `Update an existing user's information.

You can provide user data via flags or JSON file/stdin.
Only specified fields will be updated.

Examples:
  # Using flags
  canvas users update 123 --name "John Smith"
  canvas users update 123 --email "newemail@example.com" --timezone "America/New_York"

  # Using JSON file
  canvas users update 123 --json updates.json

  # Using stdin
  echo '{"name":"Updated Name"}' | canvas users update 123 --stdin`,
	Args: ExactArgsWithUsage(1, "user-id"),
	RunE: runUsersUpdate,
}

func init() {
	rootCmd.AddCommand(usersCmd)
	usersCmd.AddCommand(usersListCmd)
	usersCmd.AddCommand(usersGetCmd)
	usersCmd.AddCommand(usersMeCmd)
	usersCmd.AddCommand(usersSearchCmd)
	usersCmd.AddCommand(usersCreateCmd)
	usersCmd.AddCommand(usersUpdateCmd)

	// List flags
	usersListCmd.Flags().Int64Var(&usersAccountID, "account-id", 0, "Account ID (for account users)")
	usersListCmd.Flags().Int64Var(&usersCourseID, "course-id", 0, "Course ID (for course enrollees)")
	usersListCmd.Flags().StringVar(&usersSearchTerm, "search", "", "Search by name, login ID, or email")
	usersListCmd.Flags().StringVar(&usersEnrollmentType, "enrollment-type", "", "Filter by enrollment type (student, teacher, ta, observer, designer)")
	usersListCmd.Flags().StringVar(&usersEnrollmentState, "enrollment-state", "", "Filter by enrollment state (active, invited, rejected, completed, inactive)")
	usersListCmd.Flags().StringSliceVar(&usersInclude, "include", []string{}, "Additional data to include (comma-separated)")

	// Get flags
	usersGetCmd.Flags().StringSliceVar(&usersInclude, "include", []string{}, "Additional data to include (comma-separated)")

	// Create flags
	usersCreateCmd.Flags().Int64Var(&usersAccountID, "account-id", 0, "Account ID (required)")
	usersCreateCmd.Flags().StringVar(&userName, "name", "", "User's full name")
	usersCreateCmd.Flags().StringVar(&userShortName, "short-name", "", "User's display name")
	usersCreateCmd.Flags().StringVar(&userSortableName, "sortable-name", "", "User's sortable name (e.g., 'Doe, John')")
	usersCreateCmd.Flags().StringVar(&userEmail, "email", "", "User's email address")
	usersCreateCmd.Flags().StringVar(&userLoginID, "login-id", "", "Login ID (unique identifier)")
	usersCreateCmd.Flags().StringVar(&userPassword, "password", "", "User's password")
	usersCreateCmd.Flags().StringVar(&userSISUserID, "sis-user-id", "", "SIS User ID")
	usersCreateCmd.Flags().StringVar(&userTimeZone, "timezone", "", "User's timezone")
	usersCreateCmd.Flags().StringVar(&userLocale, "locale", "", "User's locale (e.g., 'en')")
	usersCreateCmd.Flags().BoolVar(&userSkipRegistration, "skip-registration", false, "Skip registration email")
	usersCreateCmd.Flags().BoolVar(&userSkipConfirmation, "skip-confirmation", false, "Skip email confirmation")
	usersCreateCmd.Flags().StringVar(&userJSONFile, "json", "", "JSON file with user data")
	usersCreateCmd.Flags().BoolVar(&userStdin, "stdin", false, "Read JSON from stdin")
	usersCreateCmd.MarkFlagRequired("account-id")

	// Update flags
	usersUpdateCmd.Flags().StringVar(&userName, "name", "", "User's full name")
	usersUpdateCmd.Flags().StringVar(&userShortName, "short-name", "", "User's display name")
	usersUpdateCmd.Flags().StringVar(&userSortableName, "sortable-name", "", "User's sortable name (e.g., 'Doe, John')")
	usersUpdateCmd.Flags().StringVar(&userEmail, "email", "", "User's email address")
	usersUpdateCmd.Flags().StringVar(&userTimeZone, "timezone", "", "User's timezone")
	usersUpdateCmd.Flags().StringVar(&userLocale, "locale", "", "User's locale (e.g., 'en')")
	usersUpdateCmd.Flags().StringVar(&userJSONFile, "json", "", "JSON file with user data")
	usersUpdateCmd.Flags().BoolVar(&userStdin, "stdin", false, "Read JSON from stdin")
}

func runUsersList(cmd *cobra.Command, args []string) error {
	// Use default account ID if neither account nor course is specified
	accountID := usersAccountID
	if accountID == 0 && usersCourseID == 0 {
		defaultID, err := getDefaultAccountID()
		if err != nil || defaultID == 0 {
			return fmt.Errorf("must specify --account-id or --course-id (no default account configured). Use 'canvas config account --detect' to set one")
		}
		accountID = defaultID
		printVerbose("Using default account ID: %d\n", accountID)
	}

	// Validate that only one context is specified
	if accountID > 0 && usersCourseID > 0 {
		return fmt.Errorf("can only specify one of --account-id or --course-id")
	}

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create users service
	usersService := api.NewUsersService(client)

	// Build options
	opts := &api.ListUsersOptions{
		SearchTerm:      usersSearchTerm,
		EnrollmentType:  usersEnrollmentType,
		EnrollmentState: usersEnrollmentState,
		Include:         usersInclude,
	}

	// List users based on context
	ctx := context.Background()
	var users []api.User
	var contextName string

	if accountID > 0 {
		// Account context - list all users in the account
		users, err = usersService.List(ctx, accountID, opts)
		contextName = fmt.Sprintf("account %d", accountID)
	} else {
		// Course context - list users enrolled in the course
		users, err = usersService.ListCourseUsers(ctx, usersCourseID, opts)
		contextName = fmt.Sprintf("course %d", usersCourseID)
	}

	if err != nil {
		return fmt.Errorf("failed to list users: %w", err)
	}

	if len(users) == 0 {
		fmt.Printf("No users found in %s\n", contextName)
		return nil
	}

	// Format and display users
	printVerbose("Found %d users in %s:\n\n", len(users), contextName)

	return formatOutput(users, nil)
}

func runUsersGet(cmd *cobra.Command, args []string) error {
	// Parse user ID
	userID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create users service
	usersService := api.NewUsersService(client)

	// Get user
	ctx := context.Background()
	user, err := usersService.Get(ctx, userID, usersInclude)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Format and display user details
	return formatOutput(user, nil)
}

func runUsersMe(cmd *cobra.Command, args []string) error {
	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create users service
	usersService := api.NewUsersService(client)

	// Get current user
	ctx := context.Background()
	user, err := usersService.GetCurrentUser(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	// Format and display user details
	return formatOutput(user, nil)
}

func runUsersSearch(cmd *cobra.Command, args []string) error {
	searchTerm := args[0]

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create users service
	usersService := api.NewUsersService(client)

	// Search users
	ctx := context.Background()
	users, err := usersService.Search(ctx, searchTerm)
	if err != nil {
		return fmt.Errorf("failed to search users: %w", err)
	}

	if len(users) == 0 {
		fmt.Printf("No users found matching '%s'\n", searchTerm)
		return nil
	}

	// Format and display users
	printVerbose("Found %d users matching '%s':\n\n", len(users), searchTerm)

	return formatOutput(users, nil)
}

func runUsersCreate(cmd *cobra.Command, args []string) error {
	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create users service
	usersService := api.NewUsersService(client)

	// Build params from flags or JSON
	params := &api.CreateUserParams{}

	// Check for JSON input
	if userJSONFile != "" || userStdin {
		jsonData, err := readUserJSON(userJSONFile, userStdin)
		if err != nil {
			return fmt.Errorf("failed to read JSON: %w", err)
		}
		if err := parseUserCreateJSON(jsonData, params); err != nil {
			return fmt.Errorf("failed to parse JSON: %w", err)
		}
	}

	// Override with flags if provided
	if userName != "" {
		params.Name = userName
	}
	if userShortName != "" {
		params.ShortName = userShortName
	}
	if userSortableName != "" {
		params.SortableName = userSortableName
	}
	if userEmail != "" {
		params.Email = userEmail
	}
	if userLoginID != "" {
		params.UniqueID = userLoginID
	}
	if userPassword != "" {
		params.Password = userPassword
	}
	if userSISUserID != "" {
		params.SISUserID = userSISUserID
	}
	if userTimeZone != "" {
		params.TimeZone = userTimeZone
	}
	if userLocale != "" {
		params.Locale = userLocale
	}
	if cmd.Flags().Changed("skip-registration") {
		params.SkipRegistration = userSkipRegistration
	}
	if cmd.Flags().Changed("skip-confirmation") {
		params.SkipConfirmation = userSkipConfirmation
	}

	// Validate required fields
	if params.Name == "" {
		return fmt.Errorf("user name is required (use --name or provide in JSON)")
	}

	// Create user
	ctx := context.Background()
	user, err := usersService.Create(ctx, usersAccountID, params)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	fmt.Printf("User created successfully!\n")
	fmt.Printf("  ID: %d\n", user.ID)
	fmt.Printf("  Name: %s\n", user.Name)
	if user.LoginID != "" {
		fmt.Printf("  Login: %s\n", user.LoginID)
	}
	if user.Email != "" {
		fmt.Printf("  Email: %s\n", user.Email)
	}

	return nil
}

func runUsersUpdate(cmd *cobra.Command, args []string) error {
	// Parse user ID
	userID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create users service
	usersService := api.NewUsersService(client)

	// Build params from flags or JSON
	params := &api.UpdateUserParams{}

	// Check for JSON input
	if userJSONFile != "" || userStdin {
		jsonData, err := readUserJSON(userJSONFile, userStdin)
		if err != nil {
			return fmt.Errorf("failed to read JSON: %w", err)
		}
		if err := parseUserUpdateJSON(jsonData, params); err != nil {
			return fmt.Errorf("failed to parse JSON: %w", err)
		}
	}

	// Override with flags if provided
	if userName != "" {
		params.Name = userName
	}
	if userShortName != "" {
		params.ShortName = userShortName
	}
	if userSortableName != "" {
		params.SortableName = userSortableName
	}
	if userEmail != "" {
		params.Email = userEmail
	}
	if userTimeZone != "" {
		params.TimeZone = userTimeZone
	}
	if userLocale != "" {
		params.Locale = userLocale
	}

	// Update user
	ctx := context.Background()
	user, err := usersService.Update(ctx, userID, params)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	fmt.Printf("User updated successfully!\n")
	fmt.Printf("  ID: %d\n", user.ID)
	fmt.Printf("  Name: %s\n", user.Name)
	if user.LoginID != "" {
		fmt.Printf("  Login: %s\n", user.LoginID)
	}
	if user.Email != "" {
		fmt.Printf("  Email: %s\n", user.Email)
	}

	return nil
}

// Helper functions for JSON input

func readUserJSON(filePath string, useStdin bool) ([]byte, error) {
	if filePath != "" {
		return os.ReadFile(filePath)
	}
	if useStdin {
		return io.ReadAll(os.Stdin)
	}
	return nil, nil
}

type userJSONInput struct {
	Name             string `json:"name"`
	ShortName        string `json:"short_name"`
	SortableName     string `json:"sortable_name"`
	Email            string `json:"email"`
	LoginID          string `json:"login_id"`
	Password         string `json:"password"`
	SISUserID        string `json:"sis_user_id"`
	TimeZone         string `json:"time_zone"`
	Locale           string `json:"locale"`
	SkipRegistration bool   `json:"skip_registration"`
	SkipConfirmation bool   `json:"skip_confirmation"`
}

func parseUserCreateJSON(data []byte, params *api.CreateUserParams) error {
	var input userJSONInput
	if err := json.Unmarshal(data, &input); err != nil {
		return err
	}

	if input.Name != "" {
		params.Name = input.Name
	}
	if input.ShortName != "" {
		params.ShortName = input.ShortName
	}
	if input.SortableName != "" {
		params.SortableName = input.SortableName
	}
	if input.Email != "" {
		params.Email = input.Email
	}
	if input.LoginID != "" {
		params.UniqueID = input.LoginID
	}
	if input.Password != "" {
		params.Password = input.Password
	}
	if input.SISUserID != "" {
		params.SISUserID = input.SISUserID
	}
	if input.TimeZone != "" {
		params.TimeZone = input.TimeZone
	}
	if input.Locale != "" {
		params.Locale = input.Locale
	}
	params.SkipRegistration = input.SkipRegistration
	params.SkipConfirmation = input.SkipConfirmation

	return nil
}

func parseUserUpdateJSON(data []byte, params *api.UpdateUserParams) error {
	var input userJSONInput
	if err := json.Unmarshal(data, &input); err != nil {
		return err
	}

	if input.Name != "" {
		params.Name = input.Name
	}
	if input.ShortName != "" {
		params.ShortName = input.ShortName
	}
	if input.SortableName != "" {
		params.SortableName = input.SortableName
	}
	if input.Email != "" {
		params.Email = input.Email
	}
	if input.TimeZone != "" {
		params.TimeZone = input.TimeZone
	}
	if input.Locale != "" {
		params.Locale = input.Locale
	}

	return nil
}
