package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var (
	coursesEnrollmentType  string
	coursesEnrollmentState string
	coursesInclude         []string
	coursesState           []string
	// Account context flags
	coursesAccountID  int64
	coursesSearchTerm string
	coursesSort       string
	coursesOrder      string
)

// coursesCmd represents the courses command group
var coursesCmd = &cobra.Command{
	Use:   "courses",
	Short: "Manage Canvas courses",
	Long: `Manage Canvas courses including listing, viewing, creating, and updating courses.

Examples:
  canvas courses list
  canvas courses get 123
  canvas courses list --enrollment-type teacher
  canvas courses list --state available`,
}

// coursesListCmd represents the courses list command
var coursesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List courses",
	Long: `List courses for the authenticated user or for an account (admin).

By default, lists courses you are enrolled in. Use --account to list all courses
in an account (requires admin permissions).

User context (default):
  canvas courses list                    # Your enrolled courses
  canvas courses list --enrollment-type teacher

Account context (admin):
  canvas courses list --account 1        # All courses in account 1
  canvas courses list --account 1 --search "Biology"
  canvas courses list --account 1 --sort course_name --order asc

Examples:
  canvas courses list
  canvas courses list --enrollment-type student
  canvas courses list --enrollment-state active
  canvas courses list --state available
  canvas courses list --include syllabus_body,term
  canvas courses list --account 1 --search "2024"`,
	RunE: runCoursesList,
}

// coursesGetCmd represents the courses get command
var coursesGetCmd = &cobra.Command{
	Use:   "get <course-id>",
	Short: "Get details of a specific course",
	Long: `Get details of a specific course by ID.

Examples:
  canvas courses get 123
  canvas courses get 123 --include syllabus_body,term`,
	Args: cobra.ExactArgs(1),
	RunE: runCoursesGet,
}

func init() {
	rootCmd.AddCommand(coursesCmd)
	coursesCmd.AddCommand(coursesListCmd)
	coursesCmd.AddCommand(coursesGetCmd)

	// List flags - User context
	coursesListCmd.Flags().StringVar(&coursesEnrollmentType, "enrollment-type", "", "Filter by enrollment type (student, teacher, ta, observer, designer)")
	coursesListCmd.Flags().StringVar(&coursesEnrollmentState, "enrollment-state", "", "Filter by enrollment state (active, invited_or_pending, completed)")
	coursesListCmd.Flags().StringSliceVar(&coursesInclude, "include", []string{}, "Additional data to include (comma-separated)")
	coursesListCmd.Flags().StringSliceVar(&coursesState, "state", []string{}, "Filter by course state (comma-separated: available, completed, unpublished, deleted)")

	// List flags - Account context (admin)
	coursesListCmd.Flags().Int64Var(&coursesAccountID, "account", 0, "Account ID to list courses from (admin mode)")
	coursesListCmd.Flags().StringVar(&coursesSearchTerm, "search", "", "Search by course name or code (account context only)")
	coursesListCmd.Flags().StringVar(&coursesSort, "sort", "", "Sort by: course_name, sis_course_id, teacher, account_name (account context only)")
	coursesListCmd.Flags().StringVar(&coursesOrder, "order", "", "Sort order: asc, desc (account context only)")

	// Get flags
	coursesGetCmd.Flags().StringSliceVar(&coursesInclude, "include", []string{}, "Additional data to include (comma-separated)")
}

func runCoursesList(cmd *cobra.Command, args []string) error {
	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	var courses []api.Course

	// Check if account context is being used
	if coursesAccountID > 0 {
		// Account context - list all courses in the account (admin mode)
		accountsService := api.NewAccountsService(client)

		opts := &api.ListAccountCoursesOptions{
			SearchTerm: coursesSearchTerm,
			State:      coursesState,
			Include:    coursesInclude,
			Sort:       coursesSort,
			Order:      coursesOrder,
		}

		courses, err = accountsService.ListCourses(ctx, coursesAccountID, opts)
		if err != nil {
			return fmt.Errorf("failed to list account courses: %w", err)
		}

		if len(courses) == 0 {
			fmt.Printf("No courses found in account %d\n", coursesAccountID)
			fmt.Println("\nTip: Make sure you have admin permissions on this account.")
			return nil
		}

		printVerbose("Found %d courses in account %d:\n\n", len(courses), coursesAccountID)
	} else {
		// User context (default) - list enrolled courses
		coursesService := api.NewCoursesService(client)

		opts := &api.ListCoursesOptions{
			EnrollmentType:  coursesEnrollmentType,
			EnrollmentState: coursesEnrollmentState,
			Include:         coursesInclude,
			State:           coursesState,
		}

		courses, err = coursesService.List(ctx, opts)
		if err != nil {
			return fmt.Errorf("failed to list courses: %w", err)
		}

		if len(courses) == 0 {
			fmt.Println("No courses found")
			fmt.Println("\nShowing your enrolled courses. For all account courses, use:")
			fmt.Println("  canvas courses list --account <account-id>")
			fmt.Println("\nTip: Run 'canvas accounts list' to see available accounts.")
			return nil
		}

		printVerbose("Found %d enrolled courses:\n\n", len(courses))
	}

	// Format and display courses
	return formatOutput(courses, nil)
}

func runCoursesGet(cmd *cobra.Command, args []string) error {
	// Parse course ID
	var courseID int64
	if _, err := fmt.Sscanf(args[0], "%d", &courseID); err != nil {
		return fmt.Errorf("invalid course ID: %s", args[0])
	}

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create courses service
	coursesService := api.NewCoursesService(client)

	// Get course
	ctx := context.Background()
	course, err := coursesService.Get(ctx, courseID, coursesInclude)
	if err != nil {
		return fmt.Errorf("failed to get course: %w", err)
	}

	// Format and display course details
	return formatOutput(course, nil)
}
