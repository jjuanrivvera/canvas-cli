package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
	"github.com/jjuanrivvera/canvas-cli/internal/progress"
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

func init() {
	rootCmd.AddCommand(usersCmd)
	usersCmd.AddCommand(newUsersListCmd())
	usersCmd.AddCommand(newUsersGetCmd())
	usersCmd.AddCommand(newUsersMeCmd())
	usersCmd.AddCommand(newUsersSearchCmd())
	usersCmd.AddCommand(newUsersCreateCmd())
	usersCmd.AddCommand(newUsersUpdateCmd())
}

func newUsersListCmd() *cobra.Command {
	opts := &options.UsersListOptions{}

	cmd := &cobra.Command{
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
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runUsersList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID (for account users)")
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (for course enrollees)")
	cmd.Flags().StringVar(&opts.SearchTerm, "search", "", "Search by name, login ID, or email")
	cmd.Flags().StringVar(&opts.EnrollmentType, "enrollment-type", "", "Filter by enrollment type (student, teacher, ta, observer, designer)")
	cmd.Flags().StringVar(&opts.EnrollmentState, "enrollment-state", "", "Filter by enrollment state (active, invited, rejected, completed, inactive)")
	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Additional data to include (comma-separated)")

	return cmd
}

func newUsersGetCmd() *cobra.Command {
	opts := &options.UsersGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <user-id>",
		Short: "Get details of a specific user",
		Long: `Get details of a specific user by ID.

Examples:
  canvas users get 123
  canvas users get 123 --include email,enrollments,avatar_url`,
		Args: ExactArgsWithUsage(1, "user-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			userID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid user ID: %s", args[0])
			}
			opts.UserID = userID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runUsersGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Additional data to include (comma-separated)")

	return cmd
}

func newUsersMeCmd() *cobra.Command {
	opts := &options.UsersMeOptions{}

	cmd := &cobra.Command{
		Use:   "me",
		Short: "Get details of the current authenticated user",
		Long: `Get details of the current authenticated user.

Examples:
  canvas users me`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runUsersMe(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newUsersSearchCmd() *cobra.Command {
	opts := &options.UsersSearchOptions{}

	cmd := &cobra.Command{
		Use:   "search <search-term>",
		Short: "Search for users",
		Long: `Search for users across the Canvas instance.

Examples:
  canvas users search "john doe"
  canvas users search "john@example.com"`,
		Args: ExactArgsWithUsage(1, "search-term"),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.SearchTerm = args[0]
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runUsersSearch(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newUsersCreateCmd() *cobra.Command {
	opts := &options.UsersCreateOptions{}

	cmd := &cobra.Command{
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
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runUsersCreate(cmd.Context(), client, cmd, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID (required)")
	cmd.Flags().StringVar(&opts.Name, "name", "", "User's full name")
	cmd.Flags().StringVar(&opts.ShortName, "short-name", "", "User's display name")
	cmd.Flags().StringVar(&opts.SortableName, "sortable-name", "", "User's sortable name (e.g., 'Doe, John')")
	cmd.Flags().StringVar(&opts.Email, "email", "", "User's email address")
	cmd.Flags().StringVar(&opts.LoginID, "login-id", "", "Login ID (unique identifier)")
	cmd.Flags().StringVar(&opts.Password, "password", "", "User's password")
	cmd.Flags().StringVar(&opts.SISUserID, "sis-user-id", "", "SIS User ID")
	cmd.Flags().StringVar(&opts.TimeZone, "timezone", "", "User's timezone")
	cmd.Flags().StringVar(&opts.Locale, "locale", "", "User's locale (e.g., 'en')")
	cmd.Flags().BoolVar(&opts.SkipRegistration, "skip-registration", false, "Skip registration email")
	cmd.Flags().BoolVar(&opts.SkipConfirmation, "skip-confirmation", false, "Skip email confirmation")
	cmd.Flags().StringVar(&opts.JSONFile, "json", "", "JSON file with user data")
	cmd.Flags().BoolVar(&opts.Stdin, "stdin", false, "Read JSON from stdin")
	cmd.MarkFlagRequired("account-id")

	return cmd
}

func newUsersUpdateCmd() *cobra.Command {
	opts := &options.UsersUpdateOptions{}

	cmd := &cobra.Command{
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
		RunE: func(cmd *cobra.Command, args []string) error {
			userID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid user ID: %s", args[0])
			}
			opts.UserID = userID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runUsersUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "User's full name")
	cmd.Flags().StringVar(&opts.ShortName, "short-name", "", "User's display name")
	cmd.Flags().StringVar(&opts.SortableName, "sortable-name", "", "User's sortable name (e.g., 'Doe, John')")
	cmd.Flags().StringVar(&opts.Email, "email", "", "User's email address")
	cmd.Flags().StringVar(&opts.TimeZone, "timezone", "", "User's timezone")
	cmd.Flags().StringVar(&opts.Locale, "locale", "", "User's locale (e.g., 'en')")
	cmd.Flags().StringVar(&opts.JSONFile, "json", "", "JSON file with user data")
	cmd.Flags().BoolVar(&opts.Stdin, "stdin", false, "Read JSON from stdin")

	return cmd
}

func runUsersList(ctx context.Context, client *api.Client, opts *options.UsersListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "users.list", map[string]interface{}{
		"account_id":       opts.AccountID,
		"course_id":        opts.CourseID,
		"search_term":      opts.SearchTerm,
		"enrollment_type":  opts.EnrollmentType,
		"enrollment_state": opts.EnrollmentState,
	})

	// Use default account ID if neither account nor course is specified
	accountID := opts.AccountID
	if accountID == 0 && opts.CourseID == 0 {
		defaultID, err := getDefaultAccountID()
		if err != nil || defaultID == 0 {
			logger.LogCommandError(ctx, "users.list", fmt.Errorf("no account or course specified"), map[string]interface{}{
				"error": "must specify --account-id or --course-id (no default account configured)",
			})
			return fmt.Errorf("must specify --account-id or --course-id (no default account configured). Use 'canvas config account --detect' to set one")
		}
		accountID = defaultID
		printVerbose("Using default account ID: %d\n", accountID)
	}

	// Create users service
	usersService := api.NewUsersService(client)

	// Build options
	listOpts := &api.ListUsersOptions{
		SearchTerm:      opts.SearchTerm,
		EnrollmentType:  opts.EnrollmentType,
		EnrollmentState: opts.EnrollmentState,
		Include:         opts.Include,
	}

	// List users based on context
	spin := progress.New("Fetching users...")
	if !quiet {
		spin.Start()
	}

	var users []api.User
	var contextName string
	var err error

	if accountID > 0 {
		users, err = usersService.List(ctx, accountID, listOpts)
		contextName = fmt.Sprintf("account %d", accountID)
	} else {
		users, err = usersService.ListCourseUsers(ctx, opts.CourseID, listOpts)
		contextName = fmt.Sprintf("course %d", opts.CourseID)
	}
	spin.Stop()

	if err != nil {
		logger.LogCommandError(ctx, "users.list", err, map[string]interface{}{
			"context": contextName,
		})
		return fmt.Errorf("failed to list users: %w", err)
	}

	if len(users) == 0 {
		fmt.Printf("No users found in %s\n", contextName)
		logger.LogCommandComplete(ctx, "users.list", 0)
		return nil
	}

	// Format and display users
	printVerbose("Found %d users in %s:\n\n", len(users), contextName)
	logger.LogCommandComplete(ctx, "users.list", len(users))

	return formatOutput(users, nil)
}

func runUsersGet(ctx context.Context, client *api.Client, opts *options.UsersGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "users.get", map[string]interface{}{
		"user_id": opts.UserID,
		"include": opts.Include,
	})

	// Create users service
	usersService := api.NewUsersService(client)

	// Get user
	user, err := usersService.Get(ctx, opts.UserID, opts.Include)
	if err != nil {
		logger.LogCommandError(ctx, "users.get", err, map[string]interface{}{
			"user_id": opts.UserID,
		})
		return fmt.Errorf("failed to get user: %w", err)
	}

	logger.LogCommandComplete(ctx, "users.get", 1)

	// Format and display user details
	return formatOutput(user, nil)
}

func runUsersMe(ctx context.Context, client *api.Client, opts *options.UsersMeOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "users.me", map[string]interface{}{})

	// Create users service
	usersService := api.NewUsersService(client)

	// Get current user
	user, err := usersService.GetCurrentUser(ctx)
	if err != nil {
		logger.LogCommandError(ctx, "users.me", err, map[string]interface{}{})
		return fmt.Errorf("failed to get current user: %w", err)
	}

	logger.LogCommandComplete(ctx, "users.me", 1)

	// Format and display user details
	return formatOutput(user, nil)
}

func runUsersSearch(ctx context.Context, client *api.Client, opts *options.UsersSearchOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "users.search", map[string]interface{}{
		"search_term": opts.SearchTerm,
	})

	// Create users service
	usersService := api.NewUsersService(client)

	// Search users
	users, err := usersService.Search(ctx, opts.SearchTerm)
	if err != nil {
		logger.LogCommandError(ctx, "users.search", err, map[string]interface{}{
			"search_term": opts.SearchTerm,
		})
		return fmt.Errorf("failed to search users: %w", err)
	}

	if len(users) == 0 {
		fmt.Printf("No users found matching '%s'\n", opts.SearchTerm)
		logger.LogCommandComplete(ctx, "users.search", 0)
		return nil
	}

	// Format and display users
	printVerbose("Found %d users matching '%s':\n\n", len(users), opts.SearchTerm)
	logger.LogCommandComplete(ctx, "users.search", len(users))

	return formatOutput(users, nil)
}

func runUsersCreate(ctx context.Context, client *api.Client, cmd *cobra.Command, opts *options.UsersCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "users.create", map[string]interface{}{
		"account_id": opts.AccountID,
		"has_json":   opts.JSONFile != "" || opts.Stdin,
	})

	// Create users service
	usersService := api.NewUsersService(client)

	// Build params from flags or JSON
	params := &api.CreateUserParams{}

	// Check for JSON input
	if opts.JSONFile != "" || opts.Stdin {
		jsonData, err := readUserJSON(opts.JSONFile, opts.Stdin)
		if err != nil {
			logger.LogCommandError(ctx, "users.create", err, map[string]interface{}{
				"json_file": opts.JSONFile,
				"stdin":     opts.Stdin,
			})
			return fmt.Errorf("failed to read JSON: %w", err)
		}
		if err := parseUserCreateJSON(jsonData, params); err != nil {
			logger.LogCommandError(ctx, "users.create", err, map[string]interface{}{
				"error": "failed to parse JSON",
			})
			return fmt.Errorf("failed to parse JSON: %w", err)
		}
	}

	// Override with flags if provided
	if opts.Name != "" {
		params.Name = opts.Name
	}
	if opts.ShortName != "" {
		params.ShortName = opts.ShortName
	}
	if opts.SortableName != "" {
		params.SortableName = opts.SortableName
	}
	if opts.Email != "" {
		params.Email = opts.Email
	}
	if opts.LoginID != "" {
		params.UniqueID = opts.LoginID
	}
	if opts.Password != "" {
		params.Password = opts.Password
	}
	if opts.SISUserID != "" {
		params.SISUserID = opts.SISUserID
	}
	if opts.TimeZone != "" {
		params.TimeZone = opts.TimeZone
	}
	if opts.Locale != "" {
		params.Locale = opts.Locale
	}
	if cmd.Flags().Changed("skip-registration") {
		params.SkipRegistration = opts.SkipRegistration
	}
	if cmd.Flags().Changed("skip-confirmation") {
		params.SkipConfirmation = opts.SkipConfirmation
	}

	// Validate required fields
	if params.Name == "" {
		err := fmt.Errorf("user name is required (use --name or provide in JSON)")
		logger.LogCommandError(ctx, "users.create", err, map[string]interface{}{})
		return err
	}

	// Create user
	user, err := usersService.Create(ctx, opts.AccountID, params)
	if err != nil {
		logger.LogCommandError(ctx, "users.create", err, map[string]interface{}{
			"account_id": opts.AccountID,
			"user_name":  params.Name,
		})
		return fmt.Errorf("failed to create user: %w", err)
	}

	logger.LogCommandComplete(ctx, "users.create", 1)

	printInfo("User created successfully!\n")
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

func runUsersUpdate(ctx context.Context, client *api.Client, opts *options.UsersUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "users.update", map[string]interface{}{
		"user_id":  opts.UserID,
		"has_json": opts.JSONFile != "" || opts.Stdin,
	})

	// Create users service
	usersService := api.NewUsersService(client)

	// Build params from flags or JSON
	params := &api.UpdateUserParams{}

	// Check for JSON input
	if opts.JSONFile != "" || opts.Stdin {
		jsonData, err := readUserJSON(opts.JSONFile, opts.Stdin)
		if err != nil {
			logger.LogCommandError(ctx, "users.update", err, map[string]interface{}{
				"json_file": opts.JSONFile,
				"stdin":     opts.Stdin,
			})
			return fmt.Errorf("failed to read JSON: %w", err)
		}
		if err := parseUserUpdateJSON(jsonData, params); err != nil {
			logger.LogCommandError(ctx, "users.update", err, map[string]interface{}{
				"error": "failed to parse JSON",
			})
			return fmt.Errorf("failed to parse JSON: %w", err)
		}
	}

	// Override with flags if provided
	if opts.Name != "" {
		params.Name = opts.Name
	}
	if opts.ShortName != "" {
		params.ShortName = opts.ShortName
	}
	if opts.SortableName != "" {
		params.SortableName = opts.SortableName
	}
	if opts.Email != "" {
		params.Email = opts.Email
	}
	if opts.TimeZone != "" {
		params.TimeZone = opts.TimeZone
	}
	if opts.Locale != "" {
		params.Locale = opts.Locale
	}

	// Update user
	user, err := usersService.Update(ctx, opts.UserID, params)
	if err != nil {
		logger.LogCommandError(ctx, "users.update", err, map[string]interface{}{
			"user_id": opts.UserID,
		})
		return fmt.Errorf("failed to update user: %w", err)
	}

	logger.LogCommandComplete(ctx, "users.update", 1)

	printInfo("User updated successfully!\n")
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
