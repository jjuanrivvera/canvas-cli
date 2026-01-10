package commands

import (
	"context"
	"fmt"
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

You must specify one of --account-id or --course-id to indicate the context.

Account context (admin):
  canvas users list --account-id 1        # All users in account 1
  canvas users list --account-id 1 --search "john"
  canvas users list --account-id 1 --enrollment-type student

Course context:
  canvas users list --course-id 123       # All users enrolled in course 123
  canvas users list --course-id 123 --enrollment-type teacher
  canvas users list --course-id 123 --include enrollments,email

Examples:
  canvas users list --account-id 1
  canvas users list --course-id 123
  canvas users list --account-id 1 --search "john"
  canvas users list --course-id 123 --enrollment-type student
  canvas users list --account-id 1 --include email,enrollments`,
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
	Args: cobra.ExactArgs(1),
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
	Args: cobra.ExactArgs(1),
	RunE: runUsersSearch,
}

func init() {
	rootCmd.AddCommand(usersCmd)
	usersCmd.AddCommand(usersListCmd)
	usersCmd.AddCommand(usersGetCmd)
	usersCmd.AddCommand(usersMeCmd)
	usersCmd.AddCommand(usersSearchCmd)

	// List flags
	usersListCmd.Flags().Int64Var(&usersAccountID, "account-id", 0, "Account ID (for account users)")
	usersListCmd.Flags().Int64Var(&usersCourseID, "course-id", 0, "Course ID (for course enrollees)")
	usersListCmd.Flags().StringVar(&usersSearchTerm, "search", "", "Search by name, login ID, or email")
	usersListCmd.Flags().StringVar(&usersEnrollmentType, "enrollment-type", "", "Filter by enrollment type (student, teacher, ta, observer, designer)")
	usersListCmd.Flags().StringVar(&usersEnrollmentState, "enrollment-state", "", "Filter by enrollment state (active, invited, rejected, completed, inactive)")
	usersListCmd.Flags().StringSliceVar(&usersInclude, "include", []string{}, "Additional data to include (comma-separated)")

	// Get flags
	usersGetCmd.Flags().StringSliceVar(&usersInclude, "include", []string{}, "Additional data to include (comma-separated)")
}

func runUsersList(cmd *cobra.Command, args []string) error {
	// Validate that exactly one context is specified
	contextsSpecified := 0
	if usersAccountID > 0 {
		contextsSpecified++
	}
	if usersCourseID > 0 {
		contextsSpecified++
	}

	if contextsSpecified == 0 {
		return fmt.Errorf("must specify one of --account-id or --course-id")
	}
	if contextsSpecified > 1 {
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

	if usersAccountID > 0 {
		// Account context - list all users in the account
		users, err = usersService.List(ctx, usersAccountID, opts)
		contextName = fmt.Sprintf("account %d", usersAccountID)
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

	// Display users
	fmt.Printf("Found %d users in %s:\n\n", len(users), contextName)

	for _, user := range users {
		fmt.Printf("👤 %s\n", user.Name)
		fmt.Printf("   ID: %d\n", user.ID)
		if user.LoginID != "" {
			fmt.Printf("   Login: %s\n", user.LoginID)
		}
		if user.Email != "" {
			fmt.Printf("   Email: %s\n", user.Email)
		}
		fmt.Println()
	}

	return nil
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

	// Display user details
	fmt.Printf("👤 %s\n", user.Name)
	fmt.Printf("   ID: %d\n", user.ID)
	if user.LoginID != "" {
		fmt.Printf("   Login: %s\n", user.LoginID)
	}
	if user.Email != "" {
		fmt.Printf("   Email: %s\n", user.Email)
	}
	if user.SisUserID != "" {
		fmt.Printf("   SIS ID: %s\n", user.SisUserID)
	}
	if user.Bio != "" {
		fmt.Printf("   Bio: %s\n", user.Bio)
	}

	return nil
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

	// Display user details
	fmt.Printf("👤 %s (You)\n", user.Name)
	fmt.Printf("   ID: %d\n", user.ID)
	if user.LoginID != "" {
		fmt.Printf("   Login: %s\n", user.LoginID)
	}
	if user.Email != "" {
		fmt.Printf("   Email: %s\n", user.Email)
	}
	if user.SisUserID != "" {
		fmt.Printf("   SIS ID: %s\n", user.SisUserID)
	}
	if user.Bio != "" {
		fmt.Printf("   Bio: %s\n", user.Bio)
	}

	return nil
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

	// Display users
	fmt.Printf("Found %d users matching '%s':\n\n", len(users), searchTerm)

	for _, user := range users {
		fmt.Printf("👤 %s\n", user.Name)
		fmt.Printf("   ID: %d\n", user.ID)
		if user.LoginID != "" {
			fmt.Printf("   Login: %s\n", user.LoginID)
		}
		if user.Email != "" {
			fmt.Printf("   Email: %s\n", user.Email)
		}
		fmt.Println()
	}

	return nil
}
