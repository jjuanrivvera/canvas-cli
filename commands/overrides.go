package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
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

func init() {
	rootCmd.AddCommand(overridesCmd)
	overridesCmd.AddCommand(newOverridesListCmd())
	overridesCmd.AddCommand(newOverridesGetCmd())
	overridesCmd.AddCommand(newOverridesCreateCmd())
	overridesCmd.AddCommand(newOverridesUpdateCmd())
	overridesCmd.AddCommand(newOverridesDeleteCmd())
}

func newOverridesListCmd() *cobra.Command {
	opts := &options.OverridesListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List overrides for an assignment",
		Long: `List all overrides for an assignment.

Examples:
  canvas overrides list --course-id 123 --assignment-id 456`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runOverridesList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.Flags().Int64Var(&opts.AssignmentID, "assignment-id", 0, "Assignment ID (required)")
	cmd.MarkFlagRequired("assignment-id")

	return cmd
}

func newOverridesGetCmd() *cobra.Command {
	opts := &options.OverridesGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <override-id>",
		Short: "Get override details",
		Long: `Get details of a specific assignment override.

Examples:
  canvas overrides get 789 --course-id 123 --assignment-id 456`,
		Args: ExactArgsWithUsage(1, "override-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			overrideID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid override ID: %w", err)
			}
			opts.OverrideID = overrideID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runOverridesGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.Flags().Int64Var(&opts.AssignmentID, "assignment-id", 0, "Assignment ID (required)")
	cmd.MarkFlagRequired("assignment-id")

	return cmd
}

func newOverridesCreateCmd() *cobra.Command {
	opts := &options.OverridesCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new override",
		Long: `Create a new assignment override.

You must specify one of: --student-ids, --section-id, or --group-id.

Examples:
  canvas overrides create --course-id 123 --assignment-id 456 --section-id 100 --due-at "2024-03-15T23:59:00Z"
  canvas overrides create --course-id 123 --assignment-id 456 --student-ids "200,201" --title "Extended deadline" --due-at "2024-03-20T23:59:00Z"
  canvas overrides create --course-id 123 --assignment-id 456 --group-id 50 --unlock-at "2024-03-01" --lock-at "2024-03-30"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runOverridesCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.Flags().Int64Var(&opts.AssignmentID, "assignment-id", 0, "Assignment ID (required)")
	cmd.MarkFlagRequired("assignment-id")
	cmd.Flags().StringVar(&opts.StudentIDs, "student-ids", "", "Comma-separated student IDs")
	cmd.Flags().Int64Var(&opts.SectionID, "section-id", 0, "Section ID")
	cmd.Flags().Int64Var(&opts.GroupID, "group-id", 0, "Group ID")
	cmd.Flags().StringVar(&opts.Title, "title", "", "Override title")
	cmd.Flags().StringVar(&opts.DueAt, "due-at", "", "Due date (ISO 8601)")
	cmd.Flags().StringVar(&opts.UnlockAt, "unlock-at", "", "Unlock date (ISO 8601)")
	cmd.Flags().StringVar(&opts.LockAt, "lock-at", "", "Lock date (ISO 8601)")

	return cmd
}

func newOverridesUpdateCmd() *cobra.Command {
	opts := &options.OverridesUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <override-id>",
		Short: "Update an override",
		Long: `Update an existing assignment override.

Examples:
  canvas overrides update 789 --course-id 123 --assignment-id 456 --due-at "2024-03-18T23:59:00Z"
  canvas overrides update 789 --course-id 123 --assignment-id 456 --title "New title"`,
		Args: ExactArgsWithUsage(1, "override-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			overrideID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid override ID: %w", err)
			}
			opts.OverrideID = overrideID

			// Track which fields were set
			opts.StudentIDsSet = cmd.Flags().Changed("student-ids")
			opts.TitleSet = cmd.Flags().Changed("title")
			opts.DueAtSet = cmd.Flags().Changed("due-at")
			opts.UnlockAtSet = cmd.Flags().Changed("unlock-at")
			opts.LockAtSet = cmd.Flags().Changed("lock-at")

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runOverridesUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.Flags().Int64Var(&opts.AssignmentID, "assignment-id", 0, "Assignment ID (required)")
	cmd.MarkFlagRequired("assignment-id")
	cmd.Flags().StringVar(&opts.StudentIDs, "student-ids", "", "Comma-separated student IDs")
	cmd.Flags().StringVar(&opts.Title, "title", "", "Override title")
	cmd.Flags().StringVar(&opts.DueAt, "due-at", "", "Due date (ISO 8601)")
	cmd.Flags().StringVar(&opts.UnlockAt, "unlock-at", "", "Unlock date (ISO 8601)")
	cmd.Flags().StringVar(&opts.LockAt, "lock-at", "", "Lock date (ISO 8601)")

	return cmd
}

func newOverridesDeleteCmd() *cobra.Command {
	opts := &options.OverridesDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <override-id>",
		Short: "Delete an override",
		Long: `Delete an assignment override.

Examples:
  canvas overrides delete 789 --course-id 123 --assignment-id 456
  canvas overrides delete 789 --course-id 123 --assignment-id 456 --force`,
		Args: ExactArgsWithUsage(1, "override-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			overrideID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid override ID: %w", err)
			}
			opts.OverrideID = overrideID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runOverridesDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.Flags().Int64Var(&opts.AssignmentID, "assignment-id", 0, "Assignment ID (required)")
	cmd.MarkFlagRequired("assignment-id")
	cmd.Flags().BoolVar(&opts.Force, "force", false, "Skip confirmation prompt")

	return cmd
}

func runOverridesList(ctx context.Context, client *api.Client, opts *options.OverridesListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "overrides.list", map[string]interface{}{
		"course_id":     opts.CourseID,
		"assignment_id": opts.AssignmentID,
	})

	service := api.NewOverridesService(client)

	overrides, err := service.List(ctx, opts.CourseID, opts.AssignmentID, nil)
	if err != nil {
		logger.LogCommandError(ctx, "overrides.list", err, map[string]interface{}{
			"course_id":     opts.CourseID,
			"assignment_id": opts.AssignmentID,
		})
		return fmt.Errorf("failed to list overrides: %w", err)
	}

	if len(overrides) == 0 {
		fmt.Printf("No overrides found for assignment %d in course %d\n", opts.AssignmentID, opts.CourseID)
		logger.LogCommandComplete(ctx, "overrides.list", 0)
		return nil
	}

	printVerbose("Found %d overrides for assignment %d:\n\n", len(overrides), opts.AssignmentID)
	logger.LogCommandComplete(ctx, "overrides.list", len(overrides))
	return formatOutput(overrides, nil)
}

func runOverridesGet(ctx context.Context, client *api.Client, opts *options.OverridesGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "overrides.get", map[string]interface{}{
		"course_id":     opts.CourseID,
		"assignment_id": opts.AssignmentID,
		"override_id":   opts.OverrideID,
	})

	service := api.NewOverridesService(client)

	override, err := service.Get(ctx, opts.CourseID, opts.AssignmentID, opts.OverrideID)
	if err != nil {
		logger.LogCommandError(ctx, "overrides.get", err, map[string]interface{}{
			"course_id":     opts.CourseID,
			"assignment_id": opts.AssignmentID,
			"override_id":   opts.OverrideID,
		})
		return fmt.Errorf("failed to get override: %w", err)
	}

	logger.LogCommandComplete(ctx, "overrides.get", 1)
	return formatOutput(override, nil)
}

func runOverridesCreate(ctx context.Context, client *api.Client, opts *options.OverridesCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "overrides.create", map[string]interface{}{
		"course_id":     opts.CourseID,
		"assignment_id": opts.AssignmentID,
		"student_ids":   opts.StudentIDs,
		"section_id":    opts.SectionID,
		"group_id":      opts.GroupID,
	})

	service := api.NewOverridesService(client)

	// Build params
	params := &api.AssignmentOverrideCreateParams{
		CourseSectionID: opts.SectionID,
		GroupID:         opts.GroupID,
		Title:           opts.Title,
		DueAt:           opts.DueAt,
		UnlockAt:        opts.UnlockAt,
		LockAt:          opts.LockAt,
	}

	// Parse student IDs
	if opts.StudentIDs != "" {
		parts := strings.Split(opts.StudentIDs, ",")
		studentIDs := make([]int64, 0, len(parts))
		for _, p := range parts {
			id, err := strconv.ParseInt(strings.TrimSpace(p), 10, 64)
			if err != nil {
				logger.LogCommandError(ctx, "overrides.create", err, map[string]interface{}{
					"course_id":     opts.CourseID,
					"assignment_id": opts.AssignmentID,
					"student_ids":   opts.StudentIDs,
				})
				return fmt.Errorf("invalid student ID '%s': %w", p, err)
			}
			studentIDs = append(studentIDs, id)
		}
		params.StudentIDs = studentIDs
	}

	override, err := service.Create(ctx, opts.CourseID, opts.AssignmentID, params)
	if err != nil {
		logger.LogCommandError(ctx, "overrides.create", err, map[string]interface{}{
			"course_id":     opts.CourseID,
			"assignment_id": opts.AssignmentID,
		})
		return fmt.Errorf("failed to create override: %w", err)
	}

	fmt.Printf("Override created successfully (ID: %d)\n", override.ID)
	logger.LogCommandComplete(ctx, "overrides.create", 1)
	return formatOutput(override, nil)
}

func runOverridesUpdate(ctx context.Context, client *api.Client, opts *options.OverridesUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "overrides.update", map[string]interface{}{
		"course_id":     opts.CourseID,
		"assignment_id": opts.AssignmentID,
		"override_id":   opts.OverrideID,
	})

	service := api.NewOverridesService(client)

	// Build params - only include changed flags
	params := &api.AssignmentOverrideUpdateParams{}

	if opts.StudentIDsSet {
		parts := strings.Split(opts.StudentIDs, ",")
		studentIDs := make([]int64, 0, len(parts))
		for _, p := range parts {
			id, err := strconv.ParseInt(strings.TrimSpace(p), 10, 64)
			if err != nil {
				logger.LogCommandError(ctx, "overrides.update", err, map[string]interface{}{
					"course_id":     opts.CourseID,
					"assignment_id": opts.AssignmentID,
					"override_id":   opts.OverrideID,
					"student_ids":   opts.StudentIDs,
				})
				return fmt.Errorf("invalid student ID '%s': %w", p, err)
			}
			studentIDs = append(studentIDs, id)
		}
		params.StudentIDs = &studentIDs
	}

	if opts.TitleSet {
		params.Title = &opts.Title
	}
	if opts.DueAtSet {
		params.DueAt = &opts.DueAt
	}
	if opts.UnlockAtSet {
		params.UnlockAt = &opts.UnlockAt
	}
	if opts.LockAtSet {
		params.LockAt = &opts.LockAt
	}

	override, err := service.Update(ctx, opts.CourseID, opts.AssignmentID, opts.OverrideID, params)
	if err != nil {
		logger.LogCommandError(ctx, "overrides.update", err, map[string]interface{}{
			"course_id":     opts.CourseID,
			"assignment_id": opts.AssignmentID,
			"override_id":   opts.OverrideID,
		})
		return fmt.Errorf("failed to update override: %w", err)
	}

	fmt.Printf("Override updated successfully (ID: %d)\n", override.ID)
	logger.LogCommandComplete(ctx, "overrides.update", 1)
	return formatOutput(override, nil)
}

func runOverridesDelete(ctx context.Context, client *api.Client, opts *options.OverridesDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "overrides.delete", map[string]interface{}{
		"course_id":     opts.CourseID,
		"assignment_id": opts.AssignmentID,
		"override_id":   opts.OverrideID,
		"force":         opts.Force,
	})

	// Confirmation
	if !opts.Force {
		fmt.Printf("WARNING: This will delete override %d for assignment %d.\n", opts.OverrideID, opts.AssignmentID)
		fmt.Print("Type 'yes' to confirm: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Delete cancelled")
			logger.LogCommandComplete(ctx, "overrides.delete", 0)
			return nil
		}
	}

	service := api.NewOverridesService(client)

	override, err := service.Delete(ctx, opts.CourseID, opts.AssignmentID, opts.OverrideID)
	if err != nil {
		logger.LogCommandError(ctx, "overrides.delete", err, map[string]interface{}{
			"course_id":     opts.CourseID,
			"assignment_id": opts.AssignmentID,
			"override_id":   opts.OverrideID,
		})
		return fmt.Errorf("failed to delete override: %w", err)
	}

	fmt.Printf("Override %d deleted\n", override.ID)
	logger.LogCommandComplete(ctx, "overrides.delete", 1)
	return nil
}
