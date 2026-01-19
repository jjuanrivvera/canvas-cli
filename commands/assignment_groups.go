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

func init() {
	rootCmd.AddCommand(assignmentGroupsCmd)
	assignmentGroupsCmd.AddCommand(newAssignmentGroupsListCmd())
	assignmentGroupsCmd.AddCommand(newAssignmentGroupsGetCmd())
	assignmentGroupsCmd.AddCommand(newAssignmentGroupsCreateCmd())
	assignmentGroupsCmd.AddCommand(newAssignmentGroupsUpdateCmd())
	assignmentGroupsCmd.AddCommand(newAssignmentGroupsDeleteCmd())
}

func newAssignmentGroupsListCmd() *cobra.Command {
	opts := &options.AssignmentGroupsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List assignment groups in a course",
		Long: `List all assignment groups in a course.

Examples:
  canvas assignment-groups list --course-id 123
  canvas assignment-groups list --course-id 123 --include assignments
  canvas assignment-groups list --course-id 123 --include rules,assignments`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAssignmentGroupsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Include additional data (assignments, discussion_topic, rules)")

	return cmd
}

func newAssignmentGroupsGetCmd() *cobra.Command {
	opts := &options.AssignmentGroupsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <group-id>",
		Short: "Get assignment group details",
		Long: `Get details of a specific assignment group.

Examples:
  canvas assignment-groups get 456 --course-id 123
  canvas assignment-groups get 456 --course-id 123 --include assignments`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAssignmentGroupsGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Include additional data (assignments, discussion_topic, rules)")

	return cmd
}

func newAssignmentGroupsCreateCmd() *cobra.Command {
	opts := &options.AssignmentGroupsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new assignment group",
		Long: `Create a new assignment group in a course.

Examples:
  canvas assignment-groups create --course-id 123 --name "Homework"
  canvas assignment-groups create --course-id 123 --name "Exams" --weight 40 --position 2
  canvas assignment-groups create --course-id 123 --name "Quizzes" --weight 20 --drop-lowest 1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAssignmentGroupsCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Group name (required)")
	cmd.MarkFlagRequired("name")
	cmd.Flags().IntVar(&opts.Position, "position", 0, "Position in course")
	cmd.Flags().Float64Var(&opts.Weight, "weight", 0, "Group weight percentage (0-100)")
	cmd.Flags().IntVar(&opts.DropLowest, "drop-lowest", 0, "Number of lowest scores to drop")
	cmd.Flags().IntVar(&opts.DropHighest, "drop-highest", 0, "Number of highest scores to drop")

	return cmd
}

func newAssignmentGroupsUpdateCmd() *cobra.Command {
	opts := &options.AssignmentGroupsUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <group-id>",
		Short: "Update an assignment group",
		Long: `Update an existing assignment group.

Examples:
  canvas assignment-groups update 456 --course-id 123 --name "Updated Name"
  canvas assignment-groups update 456 --course-id 123 --weight 30
  canvas assignment-groups update 456 --course-id 123 --drop-lowest 2`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID

			// Track which fields were set
			opts.NameSet = cmd.Flags().Changed("name")
			opts.PositionSet = cmd.Flags().Changed("position")
			opts.WeightSet = cmd.Flags().Changed("weight")
			opts.DropLowestSet = cmd.Flags().Changed("drop-lowest")
			opts.DropHighestSet = cmd.Flags().Changed("drop-highest")

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAssignmentGroupsUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Group name")
	cmd.Flags().IntVar(&opts.Position, "position", 0, "Position in course")
	cmd.Flags().Float64Var(&opts.Weight, "weight", 0, "Group weight percentage (0-100)")
	cmd.Flags().IntVar(&opts.DropLowest, "drop-lowest", 0, "Number of lowest scores to drop")
	cmd.Flags().IntVar(&opts.DropHighest, "drop-highest", 0, "Number of highest scores to drop")

	return cmd
}

func newAssignmentGroupsDeleteCmd() *cobra.Command {
	opts := &options.AssignmentGroupsDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <group-id>",
		Short: "Delete an assignment group",
		Long: `Delete an assignment group.

You can optionally move assignments to another group before deleting.

Examples:
  canvas assignment-groups delete 456 --course-id 123
  canvas assignment-groups delete 456 --course-id 123 --force
  canvas assignment-groups delete 456 --course-id 123 --move-to 789`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAssignmentGroupsDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.Flags().BoolVar(&opts.Force, "force", false, "Skip confirmation prompt")
	cmd.Flags().Int64Var(&opts.MoveTo, "move-to", 0, "Move assignments to another group before deleting")

	return cmd
}

func runAssignmentGroupsList(ctx context.Context, client *api.Client, opts *options.AssignmentGroupsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "assignment_groups.list", map[string]interface{}{
		"course_id": opts.CourseID,
		"include":   opts.Include,
	})

	service := api.NewAssignmentGroupsService(client)

	apiOpts := &api.ListAssignmentGroupsOptions{
		Include: opts.Include,
	}

	groups, err := service.List(ctx, opts.CourseID, apiOpts)
	if err != nil {
		logger.LogCommandError(ctx, "assignment_groups.list", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to list assignment groups: %w", err)
	}

	if len(groups) == 0 {
		fmt.Printf("No assignment groups found in course %d\n", opts.CourseID)
		logger.LogCommandComplete(ctx, "assignment_groups.list", 0)
		return nil
	}

	printVerbose("Found %d assignment groups in course %d:\n\n", len(groups), opts.CourseID)
	logger.LogCommandComplete(ctx, "assignment_groups.list", len(groups))
	return formatOutput(groups, nil)
}

func runAssignmentGroupsGet(ctx context.Context, client *api.Client, opts *options.AssignmentGroupsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "assignment_groups.get", map[string]interface{}{
		"course_id": opts.CourseID,
		"group_id":  opts.GroupID,
		"include":   opts.Include,
	})

	service := api.NewAssignmentGroupsService(client)

	group, err := service.Get(ctx, opts.CourseID, opts.GroupID, opts.Include)
	if err != nil {
		logger.LogCommandError(ctx, "assignment_groups.get", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"group_id":  opts.GroupID,
		})
		return fmt.Errorf("failed to get assignment group: %w", err)
	}

	logger.LogCommandComplete(ctx, "assignment_groups.get", 1)
	return formatOutput(group, nil)
}

func runAssignmentGroupsCreate(ctx context.Context, client *api.Client, opts *options.AssignmentGroupsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "assignment_groups.create", map[string]interface{}{
		"course_id": opts.CourseID,
		"name":      opts.Name,
		"weight":    opts.Weight,
	})

	service := api.NewAssignmentGroupsService(client)

	params := &api.CreateAssignmentGroupParams{
		Name:        opts.Name,
		Position:    opts.Position,
		GroupWeight: opts.Weight,
	}

	// Add rules if specified
	if opts.DropLowest > 0 || opts.DropHighest > 0 {
		params.Rules = &api.GradingRules{
			DropLowest:  opts.DropLowest,
			DropHighest: opts.DropHighest,
		}
	}

	group, err := service.Create(ctx, opts.CourseID, params)
	if err != nil {
		logger.LogCommandError(ctx, "assignment_groups.create", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"name":      opts.Name,
		})
		return fmt.Errorf("failed to create assignment group: %w", err)
	}

	fmt.Printf("Assignment group created successfully (ID: %d)\n", group.ID)
	logger.LogCommandComplete(ctx, "assignment_groups.create", 1)
	return formatOutput(group, nil)
}

func runAssignmentGroupsUpdate(ctx context.Context, client *api.Client, opts *options.AssignmentGroupsUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "assignment_groups.update", map[string]interface{}{
		"course_id": opts.CourseID,
		"group_id":  opts.GroupID,
	})

	service := api.NewAssignmentGroupsService(client)

	// Build params - only include changed flags
	params := &api.UpdateAssignmentGroupParams{}

	if opts.NameSet {
		params.Name = &opts.Name
	}
	if opts.PositionSet {
		params.Position = &opts.Position
	}
	if opts.WeightSet {
		params.GroupWeight = &opts.Weight
	}

	// Add rules if specified
	if opts.DropLowestSet || opts.DropHighestSet {
		params.Rules = &api.GradingRules{
			DropLowest:  opts.DropLowest,
			DropHighest: opts.DropHighest,
		}
	}

	group, err := service.Update(ctx, opts.CourseID, opts.GroupID, params)
	if err != nil {
		logger.LogCommandError(ctx, "assignment_groups.update", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"group_id":  opts.GroupID,
		})
		return fmt.Errorf("failed to update assignment group: %w", err)
	}

	fmt.Printf("Assignment group updated successfully (ID: %d)\n", group.ID)
	logger.LogCommandComplete(ctx, "assignment_groups.update", 1)
	return formatOutput(group, nil)
}

func runAssignmentGroupsDelete(ctx context.Context, client *api.Client, opts *options.AssignmentGroupsDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "assignment_groups.delete", map[string]interface{}{
		"course_id": opts.CourseID,
		"group_id":  opts.GroupID,
		"move_to":   opts.MoveTo,
		"force":     opts.Force,
	})

	// Confirmation
	if !opts.Force {
		msg := fmt.Sprintf("WARNING: This will delete assignment group %d", opts.GroupID)
		if opts.MoveTo > 0 {
			msg += fmt.Sprintf(". Assignments will be moved to group %d", opts.MoveTo)
		} else {
			msg += ". Any assignments in this group will also be deleted"
		}
		fmt.Println(msg)
		fmt.Print("Type 'yes' to confirm: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Delete cancelled")
			logger.LogCommandComplete(ctx, "assignment_groups.delete", 0)
			return nil
		}
	}

	service := api.NewAssignmentGroupsService(client)

	// Build options
	var deleteOpts *api.DeleteAssignmentGroupOptions
	if opts.MoveTo > 0 {
		deleteOpts = &api.DeleteAssignmentGroupOptions{
			MoveAssignmentsTo: opts.MoveTo,
		}
	}

	group, err := service.Delete(ctx, opts.CourseID, opts.GroupID, deleteOpts)
	if err != nil {
		logger.LogCommandError(ctx, "assignment_groups.delete", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"group_id":  opts.GroupID,
		})
		return fmt.Errorf("failed to delete assignment group: %w", err)
	}

	fmt.Printf("Assignment group %d deleted\n", group.ID)
	logger.LogCommandComplete(ctx, "assignment_groups.delete", 1)
	return nil
}
