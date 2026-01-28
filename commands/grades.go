package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// gradesCmd represents the grades command group
var gradesCmd = &cobra.Command{
	Use:   "grades",
	Short: "Manage Canvas gradebook",
	Long: `Manage Canvas gradebook history and custom columns.

View gradebook history, manage custom gradebook columns, and update grades.

Examples:
  canvas grades history --course-id 123
  canvas grades feed --course-id 123 --user-id 456
  canvas grades columns list --course-id 123`,
}

// gradesColumnsCmd represents the grades columns command group
var gradesColumnsCmd = &cobra.Command{
	Use:   "columns",
	Short: "Manage custom gradebook columns",
	Long: `Manage custom gradebook columns.

Custom columns allow instructors to add additional data columns to the gradebook.

Examples:
  canvas grades columns list --course-id 123
  canvas grades columns create --course-id 123 --title "Notes"`,
}

// gradesColumnsDataCmd represents the grades columns data command group
var gradesColumnsDataCmd = &cobra.Command{
	Use:   "data",
	Short: "Manage custom column data",
	Long: `Manage data in custom gradebook columns.

Examples:
  canvas grades columns data list 456 --course-id 123
  canvas grades columns data set 456 --course-id 123 --user-id 789 --content "Note"`,
}

func init() {
	rootCmd.AddCommand(gradesCmd)
	gradesCmd.AddCommand(newGradesHistoryCmd())
	gradesCmd.AddCommand(newGradesFeedCmd())
	gradesCmd.AddCommand(gradesColumnsCmd)

	gradesColumnsCmd.AddCommand(newGradesColumnsListCmd())
	gradesColumnsCmd.AddCommand(newGradesColumnsGetCmd())
	gradesColumnsCmd.AddCommand(newGradesColumnsCreateCmd())
	gradesColumnsCmd.AddCommand(newGradesColumnsUpdateCmd())
	gradesColumnsCmd.AddCommand(newGradesColumnsDeleteCmd())
	gradesColumnsCmd.AddCommand(gradesColumnsDataCmd)

	gradesColumnsDataCmd.AddCommand(newGradesColumnsDataListCmd())
	gradesColumnsDataCmd.AddCommand(newGradesColumnsDataSetCmd())
}

func newGradesHistoryCmd() *cobra.Command {
	opts := &options.GradesHistoryOptions{}

	cmd := &cobra.Command{
		Use:   "history",
		Short: "Get gradebook history",
		Long: `Get gradebook history days showing grading activity.

Examples:
  canvas grades history --course-id 123
  canvas grades history --course-id 123 --start-date 2024-01-01 --end-date 2024-01-31`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGradesHistory(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.Flags().StringVar(&opts.StartDate, "start-date", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.EndDate, "end-date", "", "End date (YYYY-MM-DD)")

	return cmd
}

func newGradesFeedCmd() *cobra.Command {
	opts := &options.GradesFeedOptions{}

	cmd := &cobra.Command{
		Use:   "feed",
		Short: "Get gradebook history feed",
		Long: `Get gradebook history feed showing grade changes.

Examples:
  canvas grades feed --course-id 123
  canvas grades feed --course-id 123 --user-id 456
  canvas grades feed --course-id 123 --assignment-id 789`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGradesFeed(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "Filter by user ID")
	cmd.Flags().Int64Var(&opts.AssignmentID, "assignment-id", 0, "Filter by assignment ID")
	cmd.Flags().StringVar(&opts.StartDate, "start-date", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.EndDate, "end-date", "", "End date (YYYY-MM-DD)")

	return cmd
}

func newGradesColumnsListCmd() *cobra.Command {
	opts := &options.GradesColumnsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List custom gradebook columns",
		Long: `List all custom gradebook columns in a course.

Examples:
  canvas grades columns list --course-id 123
  canvas grades columns list --course-id 123 --include-hidden`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGradesColumnsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.Flags().BoolVar(&opts.IncludeHidden, "include-hidden", false, "Include hidden columns")

	return cmd
}

func newGradesColumnsGetCmd() *cobra.Command {
	opts := &options.GradesColumnsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <column-id>",
		Short: "Get custom column details",
		Long: `Get details of a specific custom gradebook column.

Examples:
  canvas grades columns get 456 --course-id 123`,
		Args: ExactArgsWithUsage(1, "column-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			columnID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid column ID: %s", args[0])
			}
			opts.ColumnID = columnID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGradesColumnsGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newGradesColumnsCreateCmd() *cobra.Command {
	opts := &options.GradesColumnsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a custom column",
		Long: `Create a new custom gradebook column.

Examples:
  canvas grades columns create --course-id 123 --title "Notes"
  canvas grades columns create --course-id 123 --title "Attendance" --teacher-notes`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGradesColumnsCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.Flags().StringVar(&opts.Title, "title", "", "Column title (required)")
	cmd.MarkFlagRequired("title")
	cmd.Flags().IntVar(&opts.Position, "position", 0, "Column position")
	cmd.Flags().BoolVar(&opts.Hidden, "hidden", false, "Hide column")
	cmd.Flags().BoolVar(&opts.TeacherNotes, "teacher-notes", false, "Teacher notes column")
	cmd.Flags().BoolVar(&opts.ReadOnly, "read-only", false, "Read-only column")

	return cmd
}

func newGradesColumnsUpdateCmd() *cobra.Command {
	opts := &options.GradesColumnsUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <column-id>",
		Short: "Update a custom column",
		Long: `Update an existing custom gradebook column.

Examples:
  canvas grades columns update 456 --course-id 123 --title "Updated Title"
  canvas grades columns update 456 --course-id 123 --hidden`,
		Args: ExactArgsWithUsage(1, "column-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			columnID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid column ID: %s", args[0])
			}
			opts.ColumnID = columnID

			// Track which fields were set
			opts.TitleSet = cmd.Flags().Changed("title")
			opts.PositionSet = cmd.Flags().Changed("position")
			opts.HiddenSet = cmd.Flags().Changed("hidden")
			opts.TeacherNotesSet = cmd.Flags().Changed("teacher-notes")
			opts.ReadOnlySet = cmd.Flags().Changed("read-only")

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGradesColumnsUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.Flags().StringVar(&opts.Title, "title", "", "Column title")
	cmd.Flags().IntVar(&opts.Position, "position", 0, "Column position")
	cmd.Flags().BoolVar(&opts.Hidden, "hidden", false, "Hide column")
	cmd.Flags().BoolVar(&opts.TeacherNotes, "teacher-notes", false, "Teacher notes column")
	cmd.Flags().BoolVar(&opts.ReadOnly, "read-only", false, "Read-only column")

	return cmd
}

func newGradesColumnsDeleteCmd() *cobra.Command {
	opts := &options.GradesColumnsDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <column-id>",
		Short: "Delete a custom column",
		Long: `Delete a custom gradebook column.

Examples:
  canvas grades columns delete 456 --course-id 123
  canvas grades columns delete 456 --course-id 123 --force`,
		Args: ExactArgsWithUsage(1, "column-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			columnID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid column ID: %s", args[0])
			}
			opts.ColumnID = columnID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGradesColumnsDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.Flags().BoolVar(&opts.Force, "force", false, "Skip confirmation prompt")

	return cmd
}

func newGradesColumnsDataListCmd() *cobra.Command {
	opts := &options.GradesColumnsDataListOptions{}

	cmd := &cobra.Command{
		Use:   "list <column-id>",
		Short: "List custom column data",
		Long: `List all data entries for a custom gradebook column.

Examples:
  canvas grades columns data list 456 --course-id 123`,
		Args: ExactArgsWithUsage(1, "column-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			columnID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid column ID: %s", args[0])
			}
			opts.ColumnID = columnID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGradesColumnsDataList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newGradesColumnsDataSetCmd() *cobra.Command {
	opts := &options.GradesColumnsDataSetOptions{}

	cmd := &cobra.Command{
		Use:   "set <column-id>",
		Short: "Set custom column data for a user",
		Long: `Set data for a user in a custom gradebook column.

Examples:
  canvas grades columns data set 456 --course-id 123 --user-id 789 --content "Student note"`,
		Args: ExactArgsWithUsage(1, "column-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			columnID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid column ID: %s", args[0])
			}
			opts.ColumnID = columnID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGradesColumnsDataSet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (required)")
	cmd.MarkFlagRequired("user-id")
	cmd.Flags().StringVar(&opts.Content, "content", "", "Column content (required)")
	cmd.MarkFlagRequired("content")

	return cmd
}

func runGradesHistory(ctx context.Context, client *api.Client, opts *options.GradesHistoryOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "grades.history", map[string]interface{}{
		"course_id":  opts.CourseID,
		"start_date": opts.StartDate,
		"end_date":   opts.EndDate,
	})

	service := api.NewGradesService(client)

	apiOpts := &api.ListGradebookHistoryOptions{
		StartDate: opts.StartDate,
		EndDate:   opts.EndDate,
	}

	days, err := service.GetHistory(ctx, opts.CourseID, apiOpts)
	if err != nil {
		logger.LogCommandError(ctx, "grades.history", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to get gradebook history: %w", err)
	}

	if len(days) == 0 {
		fmt.Println("No gradebook history found")
		logger.LogCommandComplete(ctx, "grades.history", 0)
		return nil
	}

	printVerbose("Found %d history days:\n\n", len(days))

	logger.LogCommandComplete(ctx, "grades.history", len(days))
	return formatOutput(days, nil)
}

func runGradesFeed(ctx context.Context, client *api.Client, opts *options.GradesFeedOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "grades.feed", map[string]interface{}{
		"course_id":     opts.CourseID,
		"user_id":       opts.UserID,
		"assignment_id": opts.AssignmentID,
	})

	service := api.NewGradesService(client)

	apiOpts := &api.ListGradebookFeedOptions{
		UserID:       opts.UserID,
		AssignmentID: opts.AssignmentID,
		StartDate:    opts.StartDate,
		EndDate:      opts.EndDate,
	}

	entries, err := service.GetFeed(ctx, opts.CourseID, apiOpts)
	if err != nil {
		logger.LogCommandError(ctx, "grades.feed", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to get gradebook feed: %w", err)
	}

	if len(entries) == 0 {
		fmt.Println("No gradebook entries found")
		logger.LogCommandComplete(ctx, "grades.feed", 0)
		return nil
	}

	printVerbose("Found %d feed entries:\n\n", len(entries))

	logger.LogCommandComplete(ctx, "grades.feed", len(entries))
	return formatOutput(entries, nil)
}

func runGradesColumnsList(ctx context.Context, client *api.Client, opts *options.GradesColumnsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "grades.columns.list", map[string]interface{}{
		"course_id":      opts.CourseID,
		"include_hidden": opts.IncludeHidden,
	})

	service := api.NewGradesService(client)

	apiOpts := &api.ListCustomColumnsOptions{
		IncludeHidden: opts.IncludeHidden,
	}

	columns, err := service.ListCustomColumns(ctx, opts.CourseID, apiOpts)
	if err != nil {
		logger.LogCommandError(ctx, "grades.columns.list", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to list custom columns: %w", err)
	}

	if len(columns) == 0 {
		fmt.Println("No custom columns found")
		logger.LogCommandComplete(ctx, "grades.columns.list", 0)
		return nil
	}

	printVerbose("Found %d custom columns:\n\n", len(columns))

	logger.LogCommandComplete(ctx, "grades.columns.list", len(columns))
	return formatOutput(columns, nil)
}

func runGradesColumnsGet(ctx context.Context, client *api.Client, opts *options.GradesColumnsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "grades.columns.get", map[string]interface{}{
		"course_id": opts.CourseID,
		"column_id": opts.ColumnID,
	})

	service := api.NewGradesService(client)

	column, err := service.GetCustomColumn(ctx, opts.CourseID, opts.ColumnID)
	if err != nil {
		logger.LogCommandError(ctx, "grades.columns.get", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"column_id": opts.ColumnID,
		})
		return fmt.Errorf("failed to get custom column: %w", err)
	}

	logger.LogCommandComplete(ctx, "grades.columns.get", 1)
	return formatOutput(column, nil)
}

func runGradesColumnsCreate(ctx context.Context, client *api.Client, opts *options.GradesColumnsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "grades.columns.create", map[string]interface{}{
		"course_id": opts.CourseID,
		"title":     opts.Title,
	})

	service := api.NewGradesService(client)

	params := &api.CreateCustomColumnParams{
		Title:        opts.Title,
		Position:     opts.Position,
		Hidden:       opts.Hidden,
		TeacherNotes: opts.TeacherNotes,
		ReadOnly:     opts.ReadOnly,
	}

	column, err := service.CreateCustomColumn(ctx, opts.CourseID, params)
	if err != nil {
		logger.LogCommandError(ctx, "grades.columns.create", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"title":     opts.Title,
		})
		return fmt.Errorf("failed to create custom column: %w", err)
	}

	printInfo("Custom column created successfully (ID: %d)\n", column.ID)

	logger.LogCommandComplete(ctx, "grades.columns.create", 1)
	return formatOutput(column, nil)
}

func runGradesColumnsUpdate(ctx context.Context, client *api.Client, opts *options.GradesColumnsUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "grades.columns.update", map[string]interface{}{
		"course_id": opts.CourseID,
		"column_id": opts.ColumnID,
	})

	service := api.NewGradesService(client)

	params := &api.UpdateCustomColumnParams{}

	if opts.TitleSet {
		params.Title = &opts.Title
	}
	if opts.PositionSet {
		params.Position = &opts.Position
	}
	if opts.HiddenSet {
		params.Hidden = &opts.Hidden
	}
	if opts.TeacherNotesSet {
		params.TeacherNotes = &opts.TeacherNotes
	}
	if opts.ReadOnlySet {
		params.ReadOnly = &opts.ReadOnly
	}

	column, err := service.UpdateCustomColumn(ctx, opts.CourseID, opts.ColumnID, params)
	if err != nil {
		logger.LogCommandError(ctx, "grades.columns.update", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"column_id": opts.ColumnID,
		})
		return fmt.Errorf("failed to update custom column: %w", err)
	}

	printInfo("Custom column updated successfully (ID: %d)\n", column.ID)

	logger.LogCommandComplete(ctx, "grades.columns.update", 1)
	return formatOutput(column, nil)
}

func runGradesColumnsDelete(ctx context.Context, client *api.Client, opts *options.GradesColumnsDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "grades.columns.delete", map[string]interface{}{
		"course_id": opts.CourseID,
		"column_id": opts.ColumnID,
		"force":     opts.Force,
	})

	if !opts.Force {
		fmt.Printf("WARNING: This will delete custom column %d.\n", opts.ColumnID)
		fmt.Print("Type 'yes' to confirm: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Delete cancelled")
			logger.LogCommandComplete(ctx, "grades.columns.delete", 0)
			return nil
		}
	}

	service := api.NewGradesService(client)

	column, err := service.DeleteCustomColumn(ctx, opts.CourseID, opts.ColumnID)
	if err != nil {
		logger.LogCommandError(ctx, "grades.columns.delete", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"column_id": opts.ColumnID,
		})
		return fmt.Errorf("failed to delete custom column: %w", err)
	}

	fmt.Printf("Custom column %d deleted\n", column.ID)

	logger.LogCommandComplete(ctx, "grades.columns.delete", 1)
	return nil
}

func runGradesColumnsDataList(ctx context.Context, client *api.Client, opts *options.GradesColumnsDataListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "grades.columns.data.list", map[string]interface{}{
		"course_id": opts.CourseID,
		"column_id": opts.ColumnID,
	})

	service := api.NewGradesService(client)

	data, err := service.GetCustomColumnData(ctx, opts.CourseID, opts.ColumnID)
	if err != nil {
		logger.LogCommandError(ctx, "grades.columns.data.list", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"column_id": opts.ColumnID,
		})
		return fmt.Errorf("failed to get column data: %w", err)
	}

	if len(data) == 0 {
		fmt.Println("No column data found")
		logger.LogCommandComplete(ctx, "grades.columns.data.list", 0)
		return nil
	}

	printVerbose("Found %d data entries:\n\n", len(data))

	logger.LogCommandComplete(ctx, "grades.columns.data.list", len(data))
	return formatOutput(data, nil)
}

func runGradesColumnsDataSet(ctx context.Context, client *api.Client, opts *options.GradesColumnsDataSetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "grades.columns.data.set", map[string]interface{}{
		"course_id": opts.CourseID,
		"column_id": opts.ColumnID,
		"user_id":   opts.UserID,
	})

	service := api.NewGradesService(client)

	datum, err := service.SetCustomColumnData(ctx, opts.CourseID, opts.ColumnID, opts.UserID, opts.Content)
	if err != nil {
		logger.LogCommandError(ctx, "grades.columns.data.set", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"column_id": opts.ColumnID,
			"user_id":   opts.UserID,
		})
		return fmt.Errorf("failed to set column data: %w", err)
	}

	fmt.Printf("Column data set for user %d\n", datum.UserID)

	logger.LogCommandComplete(ctx, "grades.columns.data.set", 1)
	return formatOutput(datum, nil)
}
