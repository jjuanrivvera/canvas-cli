package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var (
	// List flags
	assignmentGroupsCourseID int64
	assignmentGroupsInclude  []string

	// Create/Update flags
	assignmentGroupsName        string
	assignmentGroupsPosition    int
	assignmentGroupsWeight      float64
	assignmentGroupsDropLowest  int
	assignmentGroupsDropHighest int

	// Delete flags
	assignmentGroupsForce  bool
	assignmentGroupsMoveTo int64
)

// assignmentGroupsCmd represents the assignment-groups command group
var assignmentGroupsCmd = &cobra.Command{
	Use:     "assignment-groups",
	Aliases: []string{"ag"},
	Short:   "Manage Canvas assignment groups",
	Long: `Manage Canvas assignment groups for organizing and weighting assignments.

Assignment groups allow you to organize assignments into categories (like Homework,
Quizzes, Exams) and optionally weight them for grade calculation.

Examples:
  canvas assignment-groups list --course-id 123
  canvas assignment-groups get 456 --course-id 123
  canvas assignment-groups create --course-id 123 --name "Homework" --weight 25`,
}

// assignmentGroupsListCmd represents the assignment-groups list command
var assignmentGroupsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List assignment groups in a course",
	Long: `List all assignment groups in a course.

Examples:
  canvas assignment-groups list --course-id 123
  canvas assignment-groups list --course-id 123 --include assignments
  canvas assignment-groups list --course-id 123 --include rules,assignments`,
	RunE: runAssignmentGroupsList,
}

// assignmentGroupsGetCmd represents the assignment-groups get command
var assignmentGroupsGetCmd = &cobra.Command{
	Use:   "get <group-id>",
	Short: "Get assignment group details",
	Long: `Get details of a specific assignment group.

Examples:
  canvas assignment-groups get 456 --course-id 123
  canvas assignment-groups get 456 --course-id 123 --include assignments`,
	Args: ExactArgsWithUsage(1, "group-id"),
	RunE: runAssignmentGroupsGet,
}

// assignmentGroupsCreateCmd represents the assignment-groups create command
var assignmentGroupsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new assignment group",
	Long: `Create a new assignment group in a course.

Examples:
  canvas assignment-groups create --course-id 123 --name "Homework"
  canvas assignment-groups create --course-id 123 --name "Exams" --weight 40 --position 2
  canvas assignment-groups create --course-id 123 --name "Quizzes" --weight 20 --drop-lowest 1`,
	RunE: runAssignmentGroupsCreate,
}

// assignmentGroupsUpdateCmd represents the assignment-groups update command
var assignmentGroupsUpdateCmd = &cobra.Command{
	Use:   "update <group-id>",
	Short: "Update an assignment group",
	Long: `Update an existing assignment group.

Examples:
  canvas assignment-groups update 456 --course-id 123 --name "Updated Name"
  canvas assignment-groups update 456 --course-id 123 --weight 30
  canvas assignment-groups update 456 --course-id 123 --drop-lowest 2`,
	Args: ExactArgsWithUsage(1, "group-id"),
	RunE: runAssignmentGroupsUpdate,
}

// assignmentGroupsDeleteCmd represents the assignment-groups delete command
var assignmentGroupsDeleteCmd = &cobra.Command{
	Use:   "delete <group-id>",
	Short: "Delete an assignment group",
	Long: `Delete an assignment group.

You can optionally move assignments to another group before deleting.

Examples:
  canvas assignment-groups delete 456 --course-id 123
  canvas assignment-groups delete 456 --course-id 123 --force
  canvas assignment-groups delete 456 --course-id 123 --move-to 789`,
	Args: ExactArgsWithUsage(1, "group-id"),
	RunE: runAssignmentGroupsDelete,
}

func init() {
	rootCmd.AddCommand(assignmentGroupsCmd)
	assignmentGroupsCmd.AddCommand(assignmentGroupsListCmd)
	assignmentGroupsCmd.AddCommand(assignmentGroupsGetCmd)
	assignmentGroupsCmd.AddCommand(assignmentGroupsCreateCmd)
	assignmentGroupsCmd.AddCommand(assignmentGroupsUpdateCmd)
	assignmentGroupsCmd.AddCommand(assignmentGroupsDeleteCmd)

	// List flags
	assignmentGroupsListCmd.Flags().Int64Var(&assignmentGroupsCourseID, "course-id", 0, "Course ID (required)")
	assignmentGroupsListCmd.MarkFlagRequired("course-id")
	assignmentGroupsListCmd.Flags().StringSliceVar(&assignmentGroupsInclude, "include", []string{}, "Include additional data (assignments, discussion_topic, rules)")

	// Get flags
	assignmentGroupsGetCmd.Flags().Int64Var(&assignmentGroupsCourseID, "course-id", 0, "Course ID (required)")
	assignmentGroupsGetCmd.MarkFlagRequired("course-id")
	assignmentGroupsGetCmd.Flags().StringSliceVar(&assignmentGroupsInclude, "include", []string{}, "Include additional data (assignments, discussion_topic, rules)")

	// Create flags
	assignmentGroupsCreateCmd.Flags().Int64Var(&assignmentGroupsCourseID, "course-id", 0, "Course ID (required)")
	assignmentGroupsCreateCmd.MarkFlagRequired("course-id")
	assignmentGroupsCreateCmd.Flags().StringVar(&assignmentGroupsName, "name", "", "Group name (required)")
	assignmentGroupsCreateCmd.MarkFlagRequired("name")
	assignmentGroupsCreateCmd.Flags().IntVar(&assignmentGroupsPosition, "position", 0, "Position in course")
	assignmentGroupsCreateCmd.Flags().Float64Var(&assignmentGroupsWeight, "weight", 0, "Group weight percentage (0-100)")
	assignmentGroupsCreateCmd.Flags().IntVar(&assignmentGroupsDropLowest, "drop-lowest", 0, "Number of lowest scores to drop")
	assignmentGroupsCreateCmd.Flags().IntVar(&assignmentGroupsDropHighest, "drop-highest", 0, "Number of highest scores to drop")

	// Update flags
	assignmentGroupsUpdateCmd.Flags().Int64Var(&assignmentGroupsCourseID, "course-id", 0, "Course ID (required)")
	assignmentGroupsUpdateCmd.MarkFlagRequired("course-id")
	assignmentGroupsUpdateCmd.Flags().StringVar(&assignmentGroupsName, "name", "", "Group name")
	assignmentGroupsUpdateCmd.Flags().IntVar(&assignmentGroupsPosition, "position", 0, "Position in course")
	assignmentGroupsUpdateCmd.Flags().Float64Var(&assignmentGroupsWeight, "weight", 0, "Group weight percentage (0-100)")
	assignmentGroupsUpdateCmd.Flags().IntVar(&assignmentGroupsDropLowest, "drop-lowest", 0, "Number of lowest scores to drop")
	assignmentGroupsUpdateCmd.Flags().IntVar(&assignmentGroupsDropHighest, "drop-highest", 0, "Number of highest scores to drop")

	// Delete flags
	assignmentGroupsDeleteCmd.Flags().Int64Var(&assignmentGroupsCourseID, "course-id", 0, "Course ID (required)")
	assignmentGroupsDeleteCmd.MarkFlagRequired("course-id")
	assignmentGroupsDeleteCmd.Flags().BoolVar(&assignmentGroupsForce, "force", false, "Skip confirmation prompt")
	assignmentGroupsDeleteCmd.Flags().Int64Var(&assignmentGroupsMoveTo, "move-to", 0, "Move assignments to another group before deleting")
}

func runAssignmentGroupsList(cmd *cobra.Command, args []string) error {
	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create assignment groups service
	service := api.NewAssignmentGroupsService(client)

	// Build options
	opts := &api.ListAssignmentGroupsOptions{
		Include: assignmentGroupsInclude,
	}

	// List groups
	ctx := context.Background()
	groups, err := service.List(ctx, assignmentGroupsCourseID, opts)
	if err != nil {
		return fmt.Errorf("failed to list assignment groups: %w", err)
	}

	if len(groups) == 0 {
		fmt.Printf("No assignment groups found in course %d\n", assignmentGroupsCourseID)
		return nil
	}

	printVerbose("Found %d assignment groups in course %d:\n\n", len(groups), assignmentGroupsCourseID)
	return formatOutput(groups, nil)
}

func runAssignmentGroupsGet(cmd *cobra.Command, args []string) error {
	// Parse group ID
	groupID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid group ID: %w", err)
	}

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create assignment groups service
	service := api.NewAssignmentGroupsService(client)

	// Get group
	ctx := context.Background()
	group, err := service.Get(ctx, assignmentGroupsCourseID, groupID, assignmentGroupsInclude)
	if err != nil {
		return fmt.Errorf("failed to get assignment group: %w", err)
	}

	return formatOutput(group, nil)
}

func runAssignmentGroupsCreate(cmd *cobra.Command, args []string) error {
	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create assignment groups service
	service := api.NewAssignmentGroupsService(client)

	// Build params
	params := &api.CreateAssignmentGroupParams{
		Name:        assignmentGroupsName,
		Position:    assignmentGroupsPosition,
		GroupWeight: assignmentGroupsWeight,
	}

	// Add rules if specified
	if assignmentGroupsDropLowest > 0 || assignmentGroupsDropHighest > 0 {
		params.Rules = &api.GradingRules{
			DropLowest:  assignmentGroupsDropLowest,
			DropHighest: assignmentGroupsDropHighest,
		}
	}

	// Create group
	ctx := context.Background()
	group, err := service.Create(ctx, assignmentGroupsCourseID, params)
	if err != nil {
		return fmt.Errorf("failed to create assignment group: %w", err)
	}

	fmt.Printf("Assignment group created successfully (ID: %d)\n", group.ID)
	return formatOutput(group, nil)
}

func runAssignmentGroupsUpdate(cmd *cobra.Command, args []string) error {
	// Parse group ID
	groupID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid group ID: %w", err)
	}

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create assignment groups service
	service := api.NewAssignmentGroupsService(client)

	// Build params - only include changed flags
	params := &api.UpdateAssignmentGroupParams{}

	if cmd.Flags().Changed("name") {
		params.Name = &assignmentGroupsName
	}
	if cmd.Flags().Changed("position") {
		params.Position = &assignmentGroupsPosition
	}
	if cmd.Flags().Changed("weight") {
		params.GroupWeight = &assignmentGroupsWeight
	}

	// Add rules if specified
	if cmd.Flags().Changed("drop-lowest") || cmd.Flags().Changed("drop-highest") {
		params.Rules = &api.GradingRules{
			DropLowest:  assignmentGroupsDropLowest,
			DropHighest: assignmentGroupsDropHighest,
		}
	}

	// Update group
	ctx := context.Background()
	group, err := service.Update(ctx, assignmentGroupsCourseID, groupID, params)
	if err != nil {
		return fmt.Errorf("failed to update assignment group: %w", err)
	}

	fmt.Printf("Assignment group updated successfully (ID: %d)\n", group.ID)
	return formatOutput(group, nil)
}

func runAssignmentGroupsDelete(cmd *cobra.Command, args []string) error {
	// Parse group ID
	groupID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid group ID: %w", err)
	}

	// Confirmation
	if !assignmentGroupsForce {
		msg := fmt.Sprintf("WARNING: This will delete assignment group %d", groupID)
		if assignmentGroupsMoveTo > 0 {
			msg += fmt.Sprintf(". Assignments will be moved to group %d", assignmentGroupsMoveTo)
		} else {
			msg += ". Any assignments in this group will also be deleted"
		}
		fmt.Println(msg)
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

	// Create assignment groups service
	service := api.NewAssignmentGroupsService(client)

	// Build options
	var opts *api.DeleteAssignmentGroupOptions
	if assignmentGroupsMoveTo > 0 {
		opts = &api.DeleteAssignmentGroupOptions{
			MoveAssignmentsTo: assignmentGroupsMoveTo,
		}
	}

	// Delete group
	ctx := context.Background()
	group, err := service.Delete(ctx, assignmentGroupsCourseID, groupID, opts)
	if err != nil {
		return fmt.Errorf("failed to delete assignment group: %w", err)
	}

	fmt.Printf("Assignment group %d deleted\n", group.ID)
	return nil
}
