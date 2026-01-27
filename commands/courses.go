package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
	"github.com/jjuanrivvera/canvas-cli/internal/progress"
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

// newCoursesListCmd creates the courses list command
func newCoursesListCmd() *cobra.Command {
	opts := &options.CoursesListOptions{}

	cmd := &cobra.Command{
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
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runCoursesList(cmd.Context(), client, opts)
		},
	}

	// User context flags
	cmd.Flags().StringVar(&opts.EnrollmentType, "enrollment-type", "", "Filter by enrollment type (student, teacher, ta, observer, designer)")
	cmd.Flags().StringVar(&opts.EnrollmentState, "enrollment-state", "", "Filter by enrollment state (active, invited_or_pending, completed)")
	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Additional data to include (comma-separated)")
	cmd.Flags().StringSliceVar(&opts.State, "state", []string{}, "Filter by course state (comma-separated: available, completed, unpublished, deleted)")

	// Account context flags (admin)
	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID to list courses from (admin mode)")
	cmd.Flags().Int64Var(&opts.AccountID, "account", 0, "Alias for --account-id")
	_ = cmd.Flags().MarkHidden("account")
	cmd.Flags().StringVar(&opts.SearchTerm, "search", "", "Search by course name or code (account context only)")
	cmd.Flags().StringVar(&opts.Sort, "sort", "", "Sort by: course_name, sis_course_id, teacher, account_name (account context only)")
	cmd.Flags().StringVar(&opts.Order, "order", "", "Sort order: asc, desc (account context only)")

	return cmd
}

func runCoursesList(ctx context.Context, client *api.Client, opts *options.CoursesListOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "courses.list", map[string]interface{}{
		"enrollment_type":  opts.EnrollmentType,
		"enrollment_state": opts.EnrollmentState,
		"account_id":       opts.AccountID,
		"search_term":      opts.SearchTerm,
	})

	spin := progress.New("Fetching courses...")
	if !quiet {
		spin.Start()
	}

	var courses []api.Course
	var err error

	// Check if account context is being used
	if opts.AccountID > 0 {
		// Account context - list all courses in the account (admin mode)
		accountsService := api.NewAccountsService(client)

		reqOpts := &api.ListAccountCoursesOptions{
			SearchTerm: opts.SearchTerm,
			State:      opts.State,
			Include:    opts.Include,
			Sort:       opts.Sort,
			Order:      opts.Order,
		}

		courses, err = accountsService.ListCourses(ctx, opts.AccountID, reqOpts)
		spin.Stop()
		if err != nil {
			logger.LogCommandError(ctx, "courses.list", err, map[string]interface{}{
				"account_id": opts.AccountID,
			})
			return fmt.Errorf("failed to list account courses: %w", err)
		}

		printVerbose("Found %d courses in account %d:\n\n", len(courses), opts.AccountID)
	} else {
		// User context (default) - list enrolled courses
		coursesService := api.NewCoursesService(client)

		reqOpts := &api.ListCoursesOptions{
			EnrollmentType:  opts.EnrollmentType,
			EnrollmentState: opts.EnrollmentState,
			Include:         opts.Include,
			State:           opts.State,
		}

		courses, err = coursesService.List(ctx, reqOpts)
		spin.Stop()
		if err != nil {
			logger.LogCommandError(ctx, "courses.list", err, map[string]interface{}{
				"enrollment_type": opts.EnrollmentType,
			})
			return fmt.Errorf("failed to list courses: %w", err)
		}

		printVerbose("Found %d enrolled courses:\n\n", len(courses))
	}

	// Format and display courses
	if err := formatEmptyOrOutput(courses, "No courses found"); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "courses.list", len(courses))
	return nil
}

// newCoursesGetCmd creates the courses get command
func newCoursesGetCmd() *cobra.Command {
	opts := &options.CoursesGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <course-id>",
		Short: "Get details of a specific course",
		Long: `Get details of a specific course by ID.

Examples:
  canvas courses get 123
  canvas courses get 123 --include syllabus_body,term`,
		Args: ExactArgsWithUsage(1, "course-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse course ID
			var courseID int64
			if _, err := fmt.Sscanf(args[0], "%d", &courseID); err != nil {
				return fmt.Errorf("invalid course ID: %s", args[0])
			}
			opts.CourseID = courseID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runCoursesGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Additional data to include (comma-separated)")

	return cmd
}

func runCoursesGet(ctx context.Context, client *api.Client, opts *options.CoursesGetOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "courses.get", map[string]interface{}{
		"course_id": opts.CourseID,
	})

	// Create courses service
	coursesService := api.NewCoursesService(client)

	// Get course
	course, err := coursesService.Get(ctx, opts.CourseID, opts.Include)
	if err != nil {
		logger.LogCommandError(ctx, "courses.get", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to get course: %w", err)
	}

	// Format and display course details
	if err := formatOutput(course, nil); err != nil {
		return fmt.Errorf("failed to print result: %w", err)
	}

	logger.LogCommandComplete(ctx, "courses.get", 1)
	return nil
}

// newCoursesCreateCmd creates the courses create command
func newCoursesCreateCmd() *cobra.Command {
	opts := &options.CoursesCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new course",
		Long: `Create a new course in an account.

Examples:
  canvas courses create --account-id 1 --name "Introduction to Programming"
  canvas courses create --account-id 1 --name "Biology 101" --code "BIO101" --term 5
  canvas courses create --account-id 1 --name "Math 201" --start-at "2024-09-01" --end-at "2024-12-15"
  canvas courses create --account-id 1 --name "Public Course" --public --offer`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runCoursesCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID to create course in (required)")
	cmd.Flags().Int64Var(&opts.AccountID, "account", 0, "Alias for --account-id")
	_ = cmd.Flags().MarkHidden("account")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Course name (required)")
	cmd.Flags().StringVar(&opts.CourseCode, "code", "", "Course code")
	cmd.Flags().StringVar(&opts.StartAt, "start-at", "", "Start date (ISO 8601)")
	cmd.Flags().StringVar(&opts.EndAt, "end-at", "", "End date (ISO 8601)")
	cmd.Flags().Int64Var(&opts.TermID, "term", 0, "Enrollment term ID")
	cmd.Flags().StringVar(&opts.License, "license", "", "Course license")
	cmd.Flags().BoolVar(&opts.IsPublic, "public", false, "Make course public")
	cmd.Flags().StringVar(&opts.SISCourseID, "sis-course-id", "", "SIS course ID")
	cmd.Flags().StringVar(&opts.DefaultView, "default-view", "", "Default view (feed, wiki, modules, syllabus, assignments)")
	cmd.Flags().BoolVar(&opts.Offer, "offer", false, "Publish course immediately")

	cmd.MarkFlagRequired("account-id")
	cmd.MarkFlagRequired("name")

	return cmd
}

func runCoursesCreate(ctx context.Context, client *api.Client, opts *options.CoursesCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "courses.create", map[string]interface{}{
		"account_id":  opts.AccountID,
		"name":        opts.Name,
		"course_code": opts.CourseCode,
	})

	// Create courses service
	coursesService := api.NewCoursesService(client)

	// Build params
	params := &api.CreateCourseParams{
		AccountID:   opts.AccountID,
		Name:        opts.Name,
		CourseCode:  opts.CourseCode,
		StartAt:     opts.StartAt,
		EndAt:       opts.EndAt,
		TermID:      opts.TermID,
		License:     opts.License,
		IsPublic:    opts.IsPublic,
		SISCourseID: opts.SISCourseID,
		DefaultView: opts.DefaultView,
		Offer:       opts.Offer,
	}

	// Create course
	course, err := coursesService.Create(ctx, params)
	if err != nil {
		logger.LogCommandError(ctx, "courses.create", err, map[string]interface{}{
			"account_id": opts.AccountID,
			"name":       opts.Name,
		})
		return fmt.Errorf("failed to create course: %w", err)
	}

	printInfo("Course created successfully (ID: %d)\n", course.ID)
	if err := formatOutput(course, nil); err != nil {
		return fmt.Errorf("failed to print result: %w", err)
	}

	logger.LogCommandComplete(ctx, "courses.create", 1)
	return nil
}

// newCoursesUpdateCmd creates the courses update command
func newCoursesUpdateCmd() *cobra.Command {
	opts := &options.CoursesUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <course-id>",
		Short: "Update a course",
		Long: `Update an existing course.

Examples:
  canvas courses update 123 --name "Updated Course Name"
  canvas courses update 123 --code "NEW101" --start-at "2024-10-01"
  canvas courses update 123 --public
  canvas courses update 123 --offer`,
		Args: ExactArgsWithUsage(1, "course-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse course ID
			var courseID int64
			if _, err := fmt.Sscanf(args[0], "%d", &courseID); err != nil {
				return fmt.Errorf("invalid course ID: %s", args[0])
			}
			opts.CourseID = courseID

			// Handle boolean flags that were explicitly set
			if cmd.Flags().Changed("public") {
				public := cmd.Flags().Lookup("public").Value.String() == "true"
				opts.IsPublic = &public
			}

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runCoursesUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "Course name")
	cmd.Flags().StringVar(&opts.CourseCode, "code", "", "Course code")
	cmd.Flags().StringVar(&opts.StartAt, "start-at", "", "Start date (ISO 8601)")
	cmd.Flags().StringVar(&opts.EndAt, "end-at", "", "End date (ISO 8601)")
	cmd.Flags().StringVar(&opts.License, "license", "", "Course license")
	cmd.Flags().Bool("public", false, "Make course public")
	cmd.Flags().StringVar(&opts.DefaultView, "default-view", "", "Default view (feed, wiki, modules, syllabus, assignments)")

	return cmd
}

func runCoursesUpdate(ctx context.Context, client *api.Client, opts *options.CoursesUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "courses.update", map[string]interface{}{
		"course_id": opts.CourseID,
	})

	// Create courses service
	coursesService := api.NewCoursesService(client)

	// Build params - only include changed values
	params := &api.UpdateCourseParams{
		Name:        opts.Name,
		CourseCode:  opts.CourseCode,
		StartAt:     opts.StartAt,
		EndAt:       opts.EndAt,
		License:     opts.License,
		DefaultView: opts.DefaultView,
		IsPublic:    opts.IsPublic,
	}

	// Update course
	course, err := coursesService.Update(ctx, opts.CourseID, params)
	if err != nil {
		logger.LogCommandError(ctx, "courses.update", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to update course: %w", err)
	}

	printInfo("Course updated successfully (ID: %d)\n", course.ID)
	if err := formatOutput(course, nil); err != nil {
		return fmt.Errorf("failed to print result: %w", err)
	}

	logger.LogCommandComplete(ctx, "courses.update", 1)
	return nil
}

// newCoursesDeleteCmd creates the courses delete command
func newCoursesDeleteCmd() *cobra.Command {
	opts := &options.CoursesDeleteOptions{
		Event: "conclude", // Default value
	}

	cmd := &cobra.Command{
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
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse course ID
			var courseID int64
			if _, err := fmt.Sscanf(args[0], "%d", &courseID); err != nil {
				return fmt.Errorf("invalid course ID: %s", args[0])
			}
			opts.CourseID = courseID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runCoursesDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Event, "event", "conclude", "Delete event: conclude, delete")
	cmd.Flags().BoolVar(&opts.Force, "force", false, "Skip confirmation prompt")

	return cmd
}

func runCoursesDelete(ctx context.Context, client *api.Client, opts *options.CoursesDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "courses.delete", map[string]interface{}{
		"course_id": opts.CourseID,
		"event":     opts.Event,
	})

	// Confirmation for delete
	if opts.Event == "delete" && !opts.Force {
		fmt.Printf("WARNING: This will permanently delete course %d and all its data.\n", opts.CourseID)
		fmt.Print("Type 'yes' to confirm: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Delete cancelled")
			return nil
		}
	}

	// Create courses service
	coursesService := api.NewCoursesService(client)

	// Delete course
	err := coursesService.Delete(ctx, opts.CourseID, opts.Event)
	if err != nil {
		logger.LogCommandError(ctx, "courses.delete", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"event":     opts.Event,
		})
		return fmt.Errorf("failed to delete course: %w", err)
	}

	if opts.Event == "delete" {
		fmt.Printf("Course %d permanently deleted\n", opts.CourseID)
	} else {
		fmt.Printf("Course %d concluded\n", opts.CourseID)
	}

	logger.LogCommandComplete(ctx, "courses.delete", 1)
	return nil
}

func init() {
	rootCmd.AddCommand(coursesCmd)
	coursesCmd.AddCommand(newCoursesListCmd())
	coursesCmd.AddCommand(newCoursesGetCmd())
	coursesCmd.AddCommand(newCoursesCreateCmd())
	coursesCmd.AddCommand(newCoursesUpdateCmd())
	coursesCmd.AddCommand(newCoursesDeleteCmd())
}
