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

	// Create/Update flags
	coursesName        string
	coursesCode        string
	coursesStartAt     string
	coursesEndAt       string
	coursesTermID      int64
	coursesLicense     string
	coursesPublic      bool
	coursesSISCourseID string
	coursesDefaultView string
	coursesOffer       bool

	// Delete flags
	coursesDeleteEvent string
	coursesForce       bool
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

By default, lists courses you are enrolled in. Use --account-id to list all courses
in an account (requires admin permissions).

User context (default):
  canvas courses list                    # Your enrolled courses
  canvas courses list --enrollment-type teacher

Account context (admin):
  canvas courses list --account-id 1        # All courses in account 1
  canvas courses list --account-id 1 --search "Biology"
  canvas courses list --account-id 1 --sort course_name --order asc

Examples:
  canvas courses list
  canvas courses list --enrollment-type student
  canvas courses list --enrollment-state active
  canvas courses list --state available
  canvas courses list --include syllabus_body,term
  canvas courses list --account-id 1 --search "2024"`,
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
	Args: ExactArgsWithUsage(1, "course-id"),
	RunE: runCoursesGet,
}

// coursesCreateCmd represents the courses create command
var coursesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new course",
	Long: `Create a new course in an account.

Examples:
  canvas courses create --account-id 1 --name "Introduction to Programming"
  canvas courses create --account-id 1 --name "Biology 101" --code "BIO101" --term 5
  canvas courses create --account-id 1 --name "Math 201" --start-at "2024-09-01" --end-at "2024-12-15"
  canvas courses create --account-id 1 --name "Public Course" --public --offer`,
	RunE: runCoursesCreate,
}

// coursesUpdateCmd represents the courses update command
var coursesUpdateCmd = &cobra.Command{
	Use:   "update <course-id>",
	Short: "Update a course",
	Long: `Update an existing course.

Examples:
  canvas courses update 123 --name "Updated Course Name"
  canvas courses update 123 --code "NEW101" --start-at "2024-10-01"
  canvas courses update 123 --public
  canvas courses update 123 --offer`,
	Args: ExactArgsWithUsage(1, "course-id"),
	RunE: runCoursesUpdate,
}

// coursesDeleteCmd represents the courses delete command
var coursesDeleteCmd = &cobra.Command{
	Use:   "delete <course-id>",
	Short: "Delete a course",
	Long: `Delete (conclude or completely remove) a course.

By default, courses are concluded (soft delete). Use --event to specify the action.

Events:
  conclude - Marks the course as concluded (default)
  delete   - Permanently deletes the course and all its data

Examples:
  canvas courses delete 123                   # Concludes the course
  canvas courses delete 123 --event conclude  # Same as above
  canvas courses delete 123 --event delete    # Permanently deletes
  canvas courses delete 123 --event delete --force  # Skip confirmation`,
	Args: ExactArgsWithUsage(1, "course-id"),
	RunE: runCoursesDelete,
}

func init() {
	rootCmd.AddCommand(coursesCmd)
	coursesCmd.AddCommand(coursesListCmd)
	coursesCmd.AddCommand(coursesGetCmd)
	coursesCmd.AddCommand(coursesCreateCmd)
	coursesCmd.AddCommand(coursesUpdateCmd)
	coursesCmd.AddCommand(coursesDeleteCmd)

	// List flags - User context
	coursesListCmd.Flags().StringVar(&coursesEnrollmentType, "enrollment-type", "", "Filter by enrollment type (student, teacher, ta, observer, designer)")
	coursesListCmd.Flags().StringVar(&coursesEnrollmentState, "enrollment-state", "", "Filter by enrollment state (active, invited_or_pending, completed)")
	coursesListCmd.Flags().StringSliceVar(&coursesInclude, "include", []string{}, "Additional data to include (comma-separated)")
	coursesListCmd.Flags().StringSliceVar(&coursesState, "state", []string{}, "Filter by course state (comma-separated: available, completed, unpublished, deleted)")

	// List flags - Account context (admin)
	coursesListCmd.Flags().Int64Var(&coursesAccountID, "account-id", 0, "Account ID to list courses from (admin mode)")
	coursesListCmd.Flags().Int64Var(&coursesAccountID, "account", 0, "Alias for --account-id")
	_ = coursesListCmd.Flags().MarkHidden("account")
	coursesListCmd.Flags().StringVar(&coursesSearchTerm, "search", "", "Search by course name or code (account context only)")
	coursesListCmd.Flags().StringVar(&coursesSort, "sort", "", "Sort by: course_name, sis_course_id, teacher, account_name (account context only)")
	coursesListCmd.Flags().StringVar(&coursesOrder, "order", "", "Sort order: asc, desc (account context only)")

	// Get flags
	coursesGetCmd.Flags().StringSliceVar(&coursesInclude, "include", []string{}, "Additional data to include (comma-separated)")

	// Create flags
	coursesCreateCmd.Flags().Int64Var(&coursesAccountID, "account-id", 0, "Account ID to create course in (required)")
	coursesCreateCmd.Flags().Int64Var(&coursesAccountID, "account", 0, "Alias for --account-id")
	_ = coursesCreateCmd.Flags().MarkHidden("account")
	coursesCreateCmd.Flags().StringVar(&coursesName, "name", "", "Course name (required)")
	coursesCreateCmd.Flags().StringVar(&coursesCode, "code", "", "Course code")
	coursesCreateCmd.Flags().StringVar(&coursesStartAt, "start-at", "", "Start date (ISO 8601)")
	coursesCreateCmd.Flags().StringVar(&coursesEndAt, "end-at", "", "End date (ISO 8601)")
	coursesCreateCmd.Flags().Int64Var(&coursesTermID, "term", 0, "Enrollment term ID")
	coursesCreateCmd.Flags().StringVar(&coursesLicense, "license", "", "Course license")
	coursesCreateCmd.Flags().BoolVar(&coursesPublic, "public", false, "Make course public")
	coursesCreateCmd.Flags().StringVar(&coursesSISCourseID, "sis-course-id", "", "SIS course ID")
	coursesCreateCmd.Flags().StringVar(&coursesDefaultView, "default-view", "", "Default view (feed, wiki, modules, syllabus, assignments)")
	coursesCreateCmd.Flags().BoolVar(&coursesOffer, "offer", false, "Publish course immediately")
	coursesCreateCmd.MarkFlagRequired("account-id")
	coursesCreateCmd.MarkFlagRequired("name")

	// Update flags
	coursesUpdateCmd.Flags().StringVar(&coursesName, "name", "", "Course name")
	coursesUpdateCmd.Flags().StringVar(&coursesCode, "code", "", "Course code")
	coursesUpdateCmd.Flags().StringVar(&coursesStartAt, "start-at", "", "Start date (ISO 8601)")
	coursesUpdateCmd.Flags().StringVar(&coursesEndAt, "end-at", "", "End date (ISO 8601)")
	coursesUpdateCmd.Flags().StringVar(&coursesLicense, "license", "", "Course license")
	coursesUpdateCmd.Flags().BoolVar(&coursesPublic, "public", false, "Make course public")
	coursesUpdateCmd.Flags().StringVar(&coursesDefaultView, "default-view", "", "Default view (feed, wiki, modules, syllabus, assignments)")

	// Delete flags
	coursesDeleteCmd.Flags().StringVar(&coursesDeleteEvent, "event", "conclude", "Delete event: conclude, delete")
	coursesDeleteCmd.Flags().BoolVar(&coursesForce, "force", false, "Skip confirmation prompt")
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

		printVerbose("Found %d enrolled courses:\n\n", len(courses))
	}

	// Format and display courses
	return formatEmptyOrOutput(courses, "No courses found")
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

func runCoursesCreate(cmd *cobra.Command, args []string) error {
	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create courses service
	coursesService := api.NewCoursesService(client)

	// Build params
	params := &api.CreateCourseParams{
		AccountID:   coursesAccountID,
		Name:        coursesName,
		CourseCode:  coursesCode,
		StartAt:     coursesStartAt,
		EndAt:       coursesEndAt,
		TermID:      coursesTermID,
		License:     coursesLicense,
		IsPublic:    coursesPublic,
		SISCourseID: coursesSISCourseID,
		DefaultView: coursesDefaultView,
		Offer:       coursesOffer,
	}

	// Create course
	ctx := context.Background()
	course, err := coursesService.Create(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create course: %w", err)
	}

	fmt.Printf("Course created successfully (ID: %d)\n", course.ID)
	return formatOutput(course, nil)
}

func runCoursesUpdate(cmd *cobra.Command, args []string) error {
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

	// Build params - only include changed values
	params := &api.UpdateCourseParams{
		Name:        coursesName,
		CourseCode:  coursesCode,
		StartAt:     coursesStartAt,
		EndAt:       coursesEndAt,
		License:     coursesLicense,
		DefaultView: coursesDefaultView,
	}

	// Handle boolean flags that were explicitly set
	if cmd.Flags().Changed("public") {
		params.IsPublic = &coursesPublic
	}

	// Update course
	ctx := context.Background()
	course, err := coursesService.Update(ctx, courseID, params)
	if err != nil {
		return fmt.Errorf("failed to update course: %w", err)
	}

	fmt.Printf("Course updated successfully (ID: %d)\n", course.ID)
	return formatOutput(course, nil)
}

func runCoursesDelete(cmd *cobra.Command, args []string) error {
	// Parse course ID
	var courseID int64
	if _, err := fmt.Sscanf(args[0], "%d", &courseID); err != nil {
		return fmt.Errorf("invalid course ID: %s", args[0])
	}

	// Validate event
	switch coursesDeleteEvent {
	case "conclude", "delete":
		// Valid
	default:
		return fmt.Errorf("invalid event: %s (use 'conclude' or 'delete')", coursesDeleteEvent)
	}

	// Confirmation for delete
	if coursesDeleteEvent == "delete" && !coursesForce {
		fmt.Printf("WARNING: This will permanently delete course %d and all its data.\n", courseID)
		fmt.Print("Type 'yes' to confirm: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Delete cancelled")
			return nil
		}
	}

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create courses service
	coursesService := api.NewCoursesService(client)

	// Delete course
	ctx := context.Background()
	err = coursesService.Delete(ctx, courseID, coursesDeleteEvent)
	if err != nil {
		return fmt.Errorf("failed to delete course: %w", err)
	}

	if coursesDeleteEvent == "delete" {
		fmt.Printf("Course %d permanently deleted\n", courseID)
	} else {
		fmt.Printf("Course %d concluded\n", courseID)
	}

	return nil
}
