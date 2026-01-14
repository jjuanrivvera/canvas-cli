package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var (
	// Common flags
	gradesCourseID int64

	// History flags
	gradesStartDate string
	gradesEndDate   string

	// Feed flags
	gradesUserID       int64
	gradesAssignmentID int64

	// Custom column flags
	gradesColumnTitle        string
	gradesColumnPosition     int
	gradesColumnHidden       bool
	gradesColumnTeacherNotes bool
	gradesColumnReadOnly     bool
	gradesIncludeHidden      bool

	// Column data flags
	gradesColumnContent string

	// Delete flags
	gradesForce bool
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

// gradesHistoryCmd represents the grades history command
var gradesHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Get gradebook history",
	Long: `Get gradebook history days showing grading activity.

Examples:
  canvas grades history --course-id 123
  canvas grades history --course-id 123 --start-date 2024-01-01 --end-date 2024-01-31`,
	RunE: runGradesHistory,
}

// gradesFeedCmd represents the grades feed command
var gradesFeedCmd = &cobra.Command{
	Use:   "feed",
	Short: "Get gradebook history feed",
	Long: `Get gradebook history feed showing grade changes.

Examples:
  canvas grades feed --course-id 123
  canvas grades feed --course-id 123 --user-id 456
  canvas grades feed --course-id 123 --assignment-id 789`,
	RunE: runGradesFeed,
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

// gradesColumnsListCmd represents the grades columns list command
var gradesColumnsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List custom gradebook columns",
	Long: `List all custom gradebook columns in a course.

Examples:
  canvas grades columns list --course-id 123
  canvas grades columns list --course-id 123 --include-hidden`,
	RunE: runGradesColumnsList,
}

// gradesColumnsGetCmd represents the grades columns get command
var gradesColumnsGetCmd = &cobra.Command{
	Use:   "get <column-id>",
	Short: "Get custom column details",
	Long: `Get details of a specific custom gradebook column.

Examples:
  canvas grades columns get 456 --course-id 123`,
	Args: ExactArgsWithUsage(1, "column-id"),
	RunE: runGradesColumnsGet,
}

// gradesColumnsCreateCmd represents the grades columns create command
var gradesColumnsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a custom column",
	Long: `Create a new custom gradebook column.

Examples:
  canvas grades columns create --course-id 123 --title "Notes"
  canvas grades columns create --course-id 123 --title "Attendance" --teacher-notes`,
	RunE: runGradesColumnsCreate,
}

// gradesColumnsUpdateCmd represents the grades columns update command
var gradesColumnsUpdateCmd = &cobra.Command{
	Use:   "update <column-id>",
	Short: "Update a custom column",
	Long: `Update an existing custom gradebook column.

Examples:
  canvas grades columns update 456 --course-id 123 --title "Updated Title"
  canvas grades columns update 456 --course-id 123 --hidden`,
	Args: ExactArgsWithUsage(1, "column-id"),
	RunE: runGradesColumnsUpdate,
}

// gradesColumnsDeleteCmd represents the grades columns delete command
var gradesColumnsDeleteCmd = &cobra.Command{
	Use:   "delete <column-id>",
	Short: "Delete a custom column",
	Long: `Delete a custom gradebook column.

Examples:
  canvas grades columns delete 456 --course-id 123
  canvas grades columns delete 456 --course-id 123 --force`,
	Args: ExactArgsWithUsage(1, "column-id"),
	RunE: runGradesColumnsDelete,
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

// gradesColumnsDataListCmd represents the grades columns data list command
var gradesColumnsDataListCmd = &cobra.Command{
	Use:   "list <column-id>",
	Short: "List custom column data",
	Long: `List all data entries for a custom gradebook column.

Examples:
  canvas grades columns data list 456 --course-id 123`,
	Args: ExactArgsWithUsage(1, "column-id"),
	RunE: runGradesColumnsDataList,
}

// gradesColumnsDataSetCmd represents the grades columns data set command
var gradesColumnsDataSetCmd = &cobra.Command{
	Use:   "set <column-id>",
	Short: "Set custom column data for a user",
	Long: `Set data for a user in a custom gradebook column.

Examples:
  canvas grades columns data set 456 --course-id 123 --user-id 789 --content "Student note"`,
	Args: ExactArgsWithUsage(1, "column-id"),
	RunE: runGradesColumnsDataSet,
}

func init() {
	rootCmd.AddCommand(gradesCmd)
	gradesCmd.AddCommand(gradesHistoryCmd)
	gradesCmd.AddCommand(gradesFeedCmd)
	gradesCmd.AddCommand(gradesColumnsCmd)

	gradesColumnsCmd.AddCommand(gradesColumnsListCmd)
	gradesColumnsCmd.AddCommand(gradesColumnsGetCmd)
	gradesColumnsCmd.AddCommand(gradesColumnsCreateCmd)
	gradesColumnsCmd.AddCommand(gradesColumnsUpdateCmd)
	gradesColumnsCmd.AddCommand(gradesColumnsDeleteCmd)
	gradesColumnsCmd.AddCommand(gradesColumnsDataCmd)

	gradesColumnsDataCmd.AddCommand(gradesColumnsDataListCmd)
	gradesColumnsDataCmd.AddCommand(gradesColumnsDataSetCmd)

	// History flags
	gradesHistoryCmd.Flags().Int64Var(&gradesCourseID, "course-id", 0, "Course ID (required)")
	gradesHistoryCmd.MarkFlagRequired("course-id")
	gradesHistoryCmd.Flags().StringVar(&gradesStartDate, "start-date", "", "Start date (YYYY-MM-DD)")
	gradesHistoryCmd.Flags().StringVar(&gradesEndDate, "end-date", "", "End date (YYYY-MM-DD)")

	// Feed flags
	gradesFeedCmd.Flags().Int64Var(&gradesCourseID, "course-id", 0, "Course ID (required)")
	gradesFeedCmd.MarkFlagRequired("course-id")
	gradesFeedCmd.Flags().Int64Var(&gradesUserID, "user-id", 0, "Filter by user ID")
	gradesFeedCmd.Flags().Int64Var(&gradesAssignmentID, "assignment-id", 0, "Filter by assignment ID")
	gradesFeedCmd.Flags().StringVar(&gradesStartDate, "start-date", "", "Start date (YYYY-MM-DD)")
	gradesFeedCmd.Flags().StringVar(&gradesEndDate, "end-date", "", "End date (YYYY-MM-DD)")

	// Columns list flags
	gradesColumnsListCmd.Flags().Int64Var(&gradesCourseID, "course-id", 0, "Course ID (required)")
	gradesColumnsListCmd.MarkFlagRequired("course-id")
	gradesColumnsListCmd.Flags().BoolVar(&gradesIncludeHidden, "include-hidden", false, "Include hidden columns")

	// Columns get flags
	gradesColumnsGetCmd.Flags().Int64Var(&gradesCourseID, "course-id", 0, "Course ID (required)")
	gradesColumnsGetCmd.MarkFlagRequired("course-id")

	// Columns create flags
	gradesColumnsCreateCmd.Flags().Int64Var(&gradesCourseID, "course-id", 0, "Course ID (required)")
	gradesColumnsCreateCmd.MarkFlagRequired("course-id")
	gradesColumnsCreateCmd.Flags().StringVar(&gradesColumnTitle, "title", "", "Column title (required)")
	gradesColumnsCreateCmd.MarkFlagRequired("title")
	gradesColumnsCreateCmd.Flags().IntVar(&gradesColumnPosition, "position", 0, "Column position")
	gradesColumnsCreateCmd.Flags().BoolVar(&gradesColumnHidden, "hidden", false, "Hide column")
	gradesColumnsCreateCmd.Flags().BoolVar(&gradesColumnTeacherNotes, "teacher-notes", false, "Teacher notes column")
	gradesColumnsCreateCmd.Flags().BoolVar(&gradesColumnReadOnly, "read-only", false, "Read-only column")

	// Columns update flags
	gradesColumnsUpdateCmd.Flags().Int64Var(&gradesCourseID, "course-id", 0, "Course ID (required)")
	gradesColumnsUpdateCmd.MarkFlagRequired("course-id")
	gradesColumnsUpdateCmd.Flags().StringVar(&gradesColumnTitle, "title", "", "Column title")
	gradesColumnsUpdateCmd.Flags().IntVar(&gradesColumnPosition, "position", 0, "Column position")
	gradesColumnsUpdateCmd.Flags().BoolVar(&gradesColumnHidden, "hidden", false, "Hide column")
	gradesColumnsUpdateCmd.Flags().BoolVar(&gradesColumnTeacherNotes, "teacher-notes", false, "Teacher notes column")
	gradesColumnsUpdateCmd.Flags().BoolVar(&gradesColumnReadOnly, "read-only", false, "Read-only column")

	// Columns delete flags
	gradesColumnsDeleteCmd.Flags().Int64Var(&gradesCourseID, "course-id", 0, "Course ID (required)")
	gradesColumnsDeleteCmd.MarkFlagRequired("course-id")
	gradesColumnsDeleteCmd.Flags().BoolVar(&gradesForce, "force", false, "Skip confirmation prompt")

	// Columns data list flags
	gradesColumnsDataListCmd.Flags().Int64Var(&gradesCourseID, "course-id", 0, "Course ID (required)")
	gradesColumnsDataListCmd.MarkFlagRequired("course-id")

	// Columns data set flags
	gradesColumnsDataSetCmd.Flags().Int64Var(&gradesCourseID, "course-id", 0, "Course ID (required)")
	gradesColumnsDataSetCmd.MarkFlagRequired("course-id")
	gradesColumnsDataSetCmd.Flags().Int64Var(&gradesUserID, "user-id", 0, "User ID (required)")
	gradesColumnsDataSetCmd.MarkFlagRequired("user-id")
	gradesColumnsDataSetCmd.Flags().StringVar(&gradesColumnContent, "content", "", "Column content (required)")
	gradesColumnsDataSetCmd.MarkFlagRequired("content")
}

func runGradesHistory(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewGradesService(client)

	opts := &api.ListGradebookHistoryOptions{
		StartDate: gradesStartDate,
		EndDate:   gradesEndDate,
	}

	ctx := context.Background()
	days, err := service.GetHistory(ctx, gradesCourseID, opts)
	if err != nil {
		return fmt.Errorf("failed to get gradebook history: %w", err)
	}

	if len(days) == 0 {
		fmt.Println("No gradebook history found")
		return nil
	}

	printVerbose("Found %d history days:\n\n", len(days))
	return formatOutput(days, nil)
}

func runGradesFeed(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewGradesService(client)

	opts := &api.ListGradebookFeedOptions{
		UserID:       gradesUserID,
		AssignmentID: gradesAssignmentID,
		StartDate:    gradesStartDate,
		EndDate:      gradesEndDate,
	}

	ctx := context.Background()
	entries, err := service.GetFeed(ctx, gradesCourseID, opts)
	if err != nil {
		return fmt.Errorf("failed to get gradebook feed: %w", err)
	}

	if len(entries) == 0 {
		fmt.Println("No gradebook entries found")
		return nil
	}

	printVerbose("Found %d feed entries:\n\n", len(entries))
	return formatOutput(entries, nil)
}

func runGradesColumnsList(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewGradesService(client)

	opts := &api.ListCustomColumnsOptions{
		IncludeHidden: gradesIncludeHidden,
	}

	ctx := context.Background()
	columns, err := service.ListCustomColumns(ctx, gradesCourseID, opts)
	if err != nil {
		return fmt.Errorf("failed to list custom columns: %w", err)
	}

	if len(columns) == 0 {
		fmt.Println("No custom columns found")
		return nil
	}

	printVerbose("Found %d custom columns:\n\n", len(columns))
	return formatOutput(columns, nil)
}

func runGradesColumnsGet(cmd *cobra.Command, args []string) error {
	columnID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid column ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewGradesService(client)

	ctx := context.Background()
	column, err := service.GetCustomColumn(ctx, gradesCourseID, columnID)
	if err != nil {
		return fmt.Errorf("failed to get custom column: %w", err)
	}

	return formatOutput(column, nil)
}

func runGradesColumnsCreate(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewGradesService(client)

	params := &api.CreateCustomColumnParams{
		Title:        gradesColumnTitle,
		Position:     gradesColumnPosition,
		Hidden:       gradesColumnHidden,
		TeacherNotes: gradesColumnTeacherNotes,
		ReadOnly:     gradesColumnReadOnly,
	}

	ctx := context.Background()
	column, err := service.CreateCustomColumn(ctx, gradesCourseID, params)
	if err != nil {
		return fmt.Errorf("failed to create custom column: %w", err)
	}

	fmt.Printf("Custom column created successfully (ID: %d)\n", column.ID)
	return formatOutput(column, nil)
}

func runGradesColumnsUpdate(cmd *cobra.Command, args []string) error {
	columnID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid column ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewGradesService(client)

	params := &api.UpdateCustomColumnParams{}

	if cmd.Flags().Changed("title") {
		params.Title = &gradesColumnTitle
	}
	if cmd.Flags().Changed("position") {
		params.Position = &gradesColumnPosition
	}
	if cmd.Flags().Changed("hidden") {
		params.Hidden = &gradesColumnHidden
	}
	if cmd.Flags().Changed("teacher-notes") {
		params.TeacherNotes = &gradesColumnTeacherNotes
	}
	if cmd.Flags().Changed("read-only") {
		params.ReadOnly = &gradesColumnReadOnly
	}

	ctx := context.Background()
	column, err := service.UpdateCustomColumn(ctx, gradesCourseID, columnID, params)
	if err != nil {
		return fmt.Errorf("failed to update custom column: %w", err)
	}

	fmt.Printf("Custom column updated successfully (ID: %d)\n", column.ID)
	return formatOutput(column, nil)
}

func runGradesColumnsDelete(cmd *cobra.Command, args []string) error {
	columnID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid column ID: %w", err)
	}

	if !gradesForce {
		fmt.Printf("WARNING: This will delete custom column %d.\n", columnID)
		fmt.Print("Type 'yes' to confirm: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Delete cancelled")
			return nil
		}
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewGradesService(client)

	ctx := context.Background()
	column, err := service.DeleteCustomColumn(ctx, gradesCourseID, columnID)
	if err != nil {
		return fmt.Errorf("failed to delete custom column: %w", err)
	}

	fmt.Printf("Custom column %d deleted\n", column.ID)
	return nil
}

func runGradesColumnsDataList(cmd *cobra.Command, args []string) error {
	columnID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid column ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewGradesService(client)

	ctx := context.Background()
	data, err := service.GetCustomColumnData(ctx, gradesCourseID, columnID)
	if err != nil {
		return fmt.Errorf("failed to get column data: %w", err)
	}

	if len(data) == 0 {
		fmt.Println("No column data found")
		return nil
	}

	printVerbose("Found %d data entries:\n\n", len(data))
	return formatOutput(data, nil)
}

func runGradesColumnsDataSet(cmd *cobra.Command, args []string) error {
	columnID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid column ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewGradesService(client)

	ctx := context.Background()
	datum, err := service.SetCustomColumnData(ctx, gradesCourseID, columnID, gradesUserID, gradesColumnContent)
	if err != nil {
		return fmt.Errorf("failed to set column data: %w", err)
	}

	fmt.Printf("Column data set for user %d\n", datum.UserID)
	return formatOutput(datum, nil)
}
