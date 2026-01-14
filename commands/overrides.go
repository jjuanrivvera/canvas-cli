package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var (
	// Common flags
	overridesCourseID     int64
	overridesAssignmentID int64

	// Create/Update flags
	overridesStudentIDs string // Comma-separated
	overridesSectionID  int64
	overridesGroupID    int64
	overridesTitle      string
	overridesDueAt      string
	overridesUnlockAt   string
	overridesLockAt     string

	// Delete flags
	overridesForce bool
)

// overridesCmd represents the overrides command group
var overridesCmd = &cobra.Command{
	Use:   "overrides",
	Short: "Manage Canvas assignment overrides",
	Long: `Manage Canvas assignment overrides for extending or modifying due dates.

Assignment overrides allow you to give specific students, sections, or groups
different due dates, availability dates, or other assignment settings.

Examples:
  canvas overrides list --course-id 123 --assignment-id 456
  canvas overrides get 789 --course-id 123 --assignment-id 456
  canvas overrides create --course-id 123 --assignment-id 456 --section-id 100 --due-at "2024-03-15T23:59:00Z"`,
}

// overridesListCmd represents the overrides list command
var overridesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List overrides for an assignment",
	Long: `List all overrides for an assignment.

Examples:
  canvas overrides list --course-id 123 --assignment-id 456`,
	RunE: runOverridesList,
}

// overridesGetCmd represents the overrides get command
var overridesGetCmd = &cobra.Command{
	Use:   "get <override-id>",
	Short: "Get override details",
	Long: `Get details of a specific assignment override.

Examples:
  canvas overrides get 789 --course-id 123 --assignment-id 456`,
	Args: ExactArgsWithUsage(1, "override-id"),
	RunE: runOverridesGet,
}

// overridesCreateCmd represents the overrides create command
var overridesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new override",
	Long: `Create a new assignment override.

You must specify one of: --student-ids, --section-id, or --group-id.

Examples:
  canvas overrides create --course-id 123 --assignment-id 456 --section-id 100 --due-at "2024-03-15T23:59:00Z"
  canvas overrides create --course-id 123 --assignment-id 456 --student-ids "200,201" --title "Extended deadline" --due-at "2024-03-20T23:59:00Z"
  canvas overrides create --course-id 123 --assignment-id 456 --group-id 50 --unlock-at "2024-03-01" --lock-at "2024-03-30"`,
	RunE: runOverridesCreate,
}

// overridesUpdateCmd represents the overrides update command
var overridesUpdateCmd = &cobra.Command{
	Use:   "update <override-id>",
	Short: "Update an override",
	Long: `Update an existing assignment override.

Examples:
  canvas overrides update 789 --course-id 123 --assignment-id 456 --due-at "2024-03-18T23:59:00Z"
  canvas overrides update 789 --course-id 123 --assignment-id 456 --title "New title"`,
	Args: ExactArgsWithUsage(1, "override-id"),
	RunE: runOverridesUpdate,
}

// overridesDeleteCmd represents the overrides delete command
var overridesDeleteCmd = &cobra.Command{
	Use:   "delete <override-id>",
	Short: "Delete an override",
	Long: `Delete an assignment override.

Examples:
  canvas overrides delete 789 --course-id 123 --assignment-id 456
  canvas overrides delete 789 --course-id 123 --assignment-id 456 --force`,
	Args: ExactArgsWithUsage(1, "override-id"),
	RunE: runOverridesDelete,
}

func init() {
	rootCmd.AddCommand(overridesCmd)
	overridesCmd.AddCommand(overridesListCmd)
	overridesCmd.AddCommand(overridesGetCmd)
	overridesCmd.AddCommand(overridesCreateCmd)
	overridesCmd.AddCommand(overridesUpdateCmd)
	overridesCmd.AddCommand(overridesDeleteCmd)

	// List flags
	overridesListCmd.Flags().Int64Var(&overridesCourseID, "course-id", 0, "Course ID (required)")
	overridesListCmd.MarkFlagRequired("course-id")
	overridesListCmd.Flags().Int64Var(&overridesAssignmentID, "assignment-id", 0, "Assignment ID (required)")
	overridesListCmd.MarkFlagRequired("assignment-id")

	// Get flags
	overridesGetCmd.Flags().Int64Var(&overridesCourseID, "course-id", 0, "Course ID (required)")
	overridesGetCmd.MarkFlagRequired("course-id")
	overridesGetCmd.Flags().Int64Var(&overridesAssignmentID, "assignment-id", 0, "Assignment ID (required)")
	overridesGetCmd.MarkFlagRequired("assignment-id")

	// Create flags
	overridesCreateCmd.Flags().Int64Var(&overridesCourseID, "course-id", 0, "Course ID (required)")
	overridesCreateCmd.MarkFlagRequired("course-id")
	overridesCreateCmd.Flags().Int64Var(&overridesAssignmentID, "assignment-id", 0, "Assignment ID (required)")
	overridesCreateCmd.MarkFlagRequired("assignment-id")
	overridesCreateCmd.Flags().StringVar(&overridesStudentIDs, "student-ids", "", "Comma-separated student IDs")
	overridesCreateCmd.Flags().Int64Var(&overridesSectionID, "section-id", 0, "Section ID")
	overridesCreateCmd.Flags().Int64Var(&overridesGroupID, "group-id", 0, "Group ID")
	overridesCreateCmd.Flags().StringVar(&overridesTitle, "title", "", "Override title")
	overridesCreateCmd.Flags().StringVar(&overridesDueAt, "due-at", "", "Due date (ISO 8601)")
	overridesCreateCmd.Flags().StringVar(&overridesUnlockAt, "unlock-at", "", "Unlock date (ISO 8601)")
	overridesCreateCmd.Flags().StringVar(&overridesLockAt, "lock-at", "", "Lock date (ISO 8601)")

	// Update flags
	overridesUpdateCmd.Flags().Int64Var(&overridesCourseID, "course-id", 0, "Course ID (required)")
	overridesUpdateCmd.MarkFlagRequired("course-id")
	overridesUpdateCmd.Flags().Int64Var(&overridesAssignmentID, "assignment-id", 0, "Assignment ID (required)")
	overridesUpdateCmd.MarkFlagRequired("assignment-id")
	overridesUpdateCmd.Flags().StringVar(&overridesStudentIDs, "student-ids", "", "Comma-separated student IDs")
	overridesUpdateCmd.Flags().StringVar(&overridesTitle, "title", "", "Override title")
	overridesUpdateCmd.Flags().StringVar(&overridesDueAt, "due-at", "", "Due date (ISO 8601)")
	overridesUpdateCmd.Flags().StringVar(&overridesUnlockAt, "unlock-at", "", "Unlock date (ISO 8601)")
	overridesUpdateCmd.Flags().StringVar(&overridesLockAt, "lock-at", "", "Lock date (ISO 8601)")

	// Delete flags
	overridesDeleteCmd.Flags().Int64Var(&overridesCourseID, "course-id", 0, "Course ID (required)")
	overridesDeleteCmd.MarkFlagRequired("course-id")
	overridesDeleteCmd.Flags().Int64Var(&overridesAssignmentID, "assignment-id", 0, "Assignment ID (required)")
	overridesDeleteCmd.MarkFlagRequired("assignment-id")
	overridesDeleteCmd.Flags().BoolVar(&overridesForce, "force", false, "Skip confirmation prompt")
}

func runOverridesList(cmd *cobra.Command, args []string) error {
	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create overrides service
	service := api.NewOverridesService(client)

	// List overrides
	ctx := context.Background()
	overrides, err := service.List(ctx, overridesCourseID, overridesAssignmentID, nil)
	if err != nil {
		return fmt.Errorf("failed to list overrides: %w", err)
	}

	if len(overrides) == 0 {
		fmt.Printf("No overrides found for assignment %d in course %d\n", overridesAssignmentID, overridesCourseID)
		return nil
	}

	printVerbose("Found %d overrides for assignment %d:\n\n", len(overrides), overridesAssignmentID)
	return formatOutput(overrides, nil)
}

func runOverridesGet(cmd *cobra.Command, args []string) error {
	// Parse override ID
	overrideID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid override ID: %w", err)
	}

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create overrides service
	service := api.NewOverridesService(client)

	// Get override
	ctx := context.Background()
	override, err := service.Get(ctx, overridesCourseID, overridesAssignmentID, overrideID)
	if err != nil {
		return fmt.Errorf("failed to get override: %w", err)
	}

	return formatOutput(override, nil)
}

func runOverridesCreate(cmd *cobra.Command, args []string) error {
	// Validate that at least one target is specified
	hasStudents := overridesStudentIDs != ""
	hasSection := overridesSectionID > 0
	hasGroup := overridesGroupID > 0

	if !hasStudents && !hasSection && !hasGroup {
		return fmt.Errorf("must specify one of --student-ids, --section-id, or --group-id")
	}

	targetsCount := 0
	if hasStudents {
		targetsCount++
	}
	if hasSection {
		targetsCount++
	}
	if hasGroup {
		targetsCount++
	}
	if targetsCount > 1 {
		return fmt.Errorf("can only specify one of --student-ids, --section-id, or --group-id")
	}

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create overrides service
	service := api.NewOverridesService(client)

	// Build params
	params := &api.AssignmentOverrideCreateParams{
		CourseSectionID: overridesSectionID,
		GroupID:         overridesGroupID,
		Title:           overridesTitle,
		DueAt:           overridesDueAt,
		UnlockAt:        overridesUnlockAt,
		LockAt:          overridesLockAt,
	}

	// Parse student IDs
	if overridesStudentIDs != "" {
		parts := strings.Split(overridesStudentIDs, ",")
		studentIDs := make([]int64, 0, len(parts))
		for _, p := range parts {
			id, err := strconv.ParseInt(strings.TrimSpace(p), 10, 64)
			if err != nil {
				return fmt.Errorf("invalid student ID '%s': %w", p, err)
			}
			studentIDs = append(studentIDs, id)
		}
		params.StudentIDs = studentIDs
	}

	// Create override
	ctx := context.Background()
	override, err := service.Create(ctx, overridesCourseID, overridesAssignmentID, params)
	if err != nil {
		return fmt.Errorf("failed to create override: %w", err)
	}

	fmt.Printf("Override created successfully (ID: %d)\n", override.ID)
	return formatOutput(override, nil)
}

func runOverridesUpdate(cmd *cobra.Command, args []string) error {
	// Parse override ID
	overrideID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid override ID: %w", err)
	}

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create overrides service
	service := api.NewOverridesService(client)

	// Build params - only include changed flags
	params := &api.AssignmentOverrideUpdateParams{}

	if cmd.Flags().Changed("student-ids") {
		parts := strings.Split(overridesStudentIDs, ",")
		studentIDs := make([]int64, 0, len(parts))
		for _, p := range parts {
			id, err := strconv.ParseInt(strings.TrimSpace(p), 10, 64)
			if err != nil {
				return fmt.Errorf("invalid student ID '%s': %w", p, err)
			}
			studentIDs = append(studentIDs, id)
		}
		params.StudentIDs = &studentIDs
	}

	if cmd.Flags().Changed("title") {
		params.Title = &overridesTitle
	}
	if cmd.Flags().Changed("due-at") {
		params.DueAt = &overridesDueAt
	}
	if cmd.Flags().Changed("unlock-at") {
		params.UnlockAt = &overridesUnlockAt
	}
	if cmd.Flags().Changed("lock-at") {
		params.LockAt = &overridesLockAt
	}

	// Update override
	ctx := context.Background()
	override, err := service.Update(ctx, overridesCourseID, overridesAssignmentID, overrideID, params)
	if err != nil {
		return fmt.Errorf("failed to update override: %w", err)
	}

	fmt.Printf("Override updated successfully (ID: %d)\n", override.ID)
	return formatOutput(override, nil)
}

func runOverridesDelete(cmd *cobra.Command, args []string) error {
	// Parse override ID
	overrideID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid override ID: %w", err)
	}

	// Confirmation
	if !overridesForce {
		fmt.Printf("WARNING: This will delete override %d for assignment %d.\n", overrideID, overridesAssignmentID)
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

	// Create overrides service
	service := api.NewOverridesService(client)

	// Delete override
	ctx := context.Background()
	override, err := service.Delete(ctx, overridesCourseID, overridesAssignmentID, overrideID)
	if err != nil {
		return fmt.Errorf("failed to delete override: %w", err)
	}

	fmt.Printf("Override %d deleted\n", override.ID)
	return nil
}
