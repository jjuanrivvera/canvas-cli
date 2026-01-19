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

// plannerCmd represents the planner command group
var plannerCmd = &cobra.Command{
	Use:   "planner",
	Short: "Manage Canvas planner items and notes",
	Long: `Manage Canvas planner items, notes, and overrides.

The planner provides a unified view of assignments, quizzes, calendar events,
and personal notes. Use this command to view upcoming items or manage notes.

Examples:
  canvas planner items
  canvas planner notes list
  canvas planner notes create --title "Study Session"
  canvas planner complete Assignment 123`,
}

// plannerNotesCmd represents the planner notes command group
var plannerNotesCmd = &cobra.Command{
	Use:   "notes",
	Short: "Manage planner notes",
	Long:  `Manage personal planner notes.`,
}

func init() {
	rootCmd.AddCommand(plannerCmd)
	plannerCmd.AddCommand(newPlannerItemsCmd())
	plannerCmd.AddCommand(plannerNotesCmd)
	plannerCmd.AddCommand(newPlannerCompleteCmd())
	plannerCmd.AddCommand(newPlannerDismissCmd())
	plannerCmd.AddCommand(newPlannerOverridesCmd())

	plannerNotesCmd.AddCommand(newPlannerNotesListCmd())
	plannerNotesCmd.AddCommand(newPlannerNotesGetCmd())
	plannerNotesCmd.AddCommand(newPlannerNotesCreateCmd())
	plannerNotesCmd.AddCommand(newPlannerNotesUpdateCmd())
	plannerNotesCmd.AddCommand(newPlannerNotesDeleteCmd())
}

func newPlannerItemsCmd() *cobra.Command {
	opts := &options.PlannerItemsOptions{}

	cmd := &cobra.Command{
		Use:   "items",
		Short: "List planner items",
		Long: `List planner items including assignments, quizzes, and calendar events.

Examples:
  canvas planner items
  canvas planner items --course-id 123
  canvas planner items --start-date 2024-01-01 --end-date 2024-01-31
  canvas planner items --filter all_assignments`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPlannerItems(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Filter by course ID")
	cmd.Flags().StringVar(&opts.StartDate, "start-date", "", "Start date (YYYY-MM-DD or ISO 8601)")
	cmd.Flags().StringVar(&opts.EndDate, "end-date", "", "End date (YYYY-MM-DD or ISO 8601)")
	cmd.Flags().StringSliceVar(&opts.ContextCodes, "context", []string{}, "Context codes (course_123)")
	cmd.Flags().StringVar(&opts.Filter, "filter", "", "Filter: all_assignments, all_quizzes, all_calendar_events, all_planner_notes")

	return cmd
}

func newPlannerNotesListCmd() *cobra.Command {
	opts := &options.PlannerNotesListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List planner notes",
		Long: `List all planner notes.

Examples:
  canvas planner notes list
  canvas planner notes list --course-id 123
  canvas planner notes list --start-date 2024-01-01`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPlannerNotesList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Filter by course ID")
	cmd.Flags().StringVar(&opts.StartDate, "start-date", "", "Start date")
	cmd.Flags().StringVar(&opts.EndDate, "end-date", "", "End date")

	return cmd
}

func newPlannerNotesGetCmd() *cobra.Command {
	opts := &options.PlannerNotesGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <note-id>",
		Short: "Get a specific planner note",
		Long: `Get details of a specific planner note.

Examples:
  canvas planner notes get 123`,
		Args: ExactArgsWithUsage(1, "note-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			noteID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid note ID: %s", args[0])
			}
			opts.NoteID = noteID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPlannerNotesGet(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newPlannerNotesCreateCmd() *cobra.Command {
	opts := &options.PlannerNotesCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new planner note",
		Long: `Create a new planner note.

Examples:
  canvas planner notes create --title "Study Session"
  canvas planner notes create --title "Project Work" --details "Work on final project" --todo-date 2024-12-15
  canvas planner notes create --title "Review" --course-id 123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPlannerNotesCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Title, "title", "", "Note title (required)")
	cmd.Flags().StringVar(&opts.Details, "details", "", "Note details")
	cmd.Flags().StringVar(&opts.TodoDate, "todo-date", "", "Todo date (ISO 8601)")
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Associate with course")
	cmd.MarkFlagRequired("title")

	return cmd
}

func newPlannerNotesUpdateCmd() *cobra.Command {
	opts := &options.PlannerNotesUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <note-id>",
		Short: "Update a planner note",
		Long: `Update an existing planner note.

Examples:
  canvas planner notes update 123 --title "Updated Title"
  canvas planner notes update 123 --todo-date 2024-12-20`,
		Args: ExactArgsWithUsage(1, "note-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			noteID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid note ID: %s", args[0])
			}
			opts.NoteID = noteID

			// Track which fields were set
			opts.TitleSet = cmd.Flags().Changed("title")
			opts.DetailsSet = cmd.Flags().Changed("details")
			opts.TodoDateSet = cmd.Flags().Changed("todo-date")
			opts.CourseIDSet = cmd.Flags().Changed("course-id")

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPlannerNotesUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Title, "title", "", "New title")
	cmd.Flags().StringVar(&opts.Details, "details", "", "New details")
	cmd.Flags().StringVar(&opts.TodoDate, "todo-date", "", "New todo date")
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "New course association")

	return cmd
}

func newPlannerNotesDeleteCmd() *cobra.Command {
	opts := &options.PlannerNotesDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <note-id>",
		Short: "Delete a planner note",
		Long: `Delete a planner note.

Examples:
  canvas planner notes delete 123`,
		Args: ExactArgsWithUsage(1, "note-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			noteID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid note ID: %s", args[0])
			}
			opts.NoteID = noteID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPlannerNotesDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Skip confirmation prompt")

	return cmd
}

func newPlannerCompleteCmd() *cobra.Command {
	opts := &options.PlannerCompleteOptions{}

	cmd := &cobra.Command{
		Use:   "complete <type> <id>",
		Short: "Mark an item as complete",
		Long: `Mark a planner item as complete by creating an override.

Type can be: Assignment, Quiz, DiscussionTopic, WikiPage, CalendarEvent

Examples:
  canvas planner complete Assignment 123
  canvas planner complete Quiz 456`,
		Args: ExactArgsWithUsage(2, "type", "id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.PlannableType = args[0]

			plannableID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid ID: %s", args[1])
			}
			opts.PlannableID = plannableID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPlannerComplete(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newPlannerDismissCmd() *cobra.Command {
	opts := &options.PlannerDismissOptions{}

	cmd := &cobra.Command{
		Use:   "dismiss <type> <id>",
		Short: "Dismiss an item from the planner",
		Long: `Dismiss a planner item so it no longer appears.

Type can be: Assignment, Quiz, DiscussionTopic, WikiPage, CalendarEvent

Examples:
  canvas planner dismiss CalendarEvent 789`,
		Args: ExactArgsWithUsage(2, "type", "id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.PlannableType = args[0]

			plannableID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid ID: %s", args[1])
			}
			opts.PlannableID = plannableID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPlannerDismiss(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newPlannerOverridesCmd() *cobra.Command {
	opts := &options.PlannerOverridesOptions{}

	cmd := &cobra.Command{
		Use:   "overrides",
		Short: "List planner overrides",
		Long: `List all planner overrides (completed/dismissed items).

Examples:
  canvas planner overrides
  canvas planner overrides --type Assignment`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPlannerOverrides(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.PlannableType, "type", "", "Filter by type (Assignment, Quiz, etc.)")

	return cmd
}

func runPlannerItems(ctx context.Context, client *api.Client, opts *options.PlannerItemsOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "planner.items", map[string]interface{}{
		"course_id": opts.CourseID,
	})

	plannerService := api.NewPlannerService(client)

	contextCodes := opts.ContextCodes
	if opts.CourseID > 0 {
		contextCodes = append(contextCodes, fmt.Sprintf("course_%d", opts.CourseID))
	}

	apiOpts := &api.ListPlannerItemsOptions{
		StartDate:    opts.StartDate,
		EndDate:      opts.EndDate,
		ContextCodes: contextCodes,
		Filter:       opts.Filter,
	}

	items, err := plannerService.ListItems(ctx, apiOpts)
	if err != nil {
		logger.LogCommandError(ctx, "planner.items", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to list planner items: %w", err)
	}

	if len(items) == 0 {
		fmt.Println("No planner items found")
		logger.LogCommandComplete(ctx, "planner.items", 0)
		return nil
	}

	printVerbose("Found %d planner items:\n\n", len(items))

	logger.LogCommandComplete(ctx, "planner.items", len(items))
	return formatOutput(items, nil)
}

func runPlannerNotesList(ctx context.Context, client *api.Client, opts *options.PlannerNotesListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "planner.notes.list", map[string]interface{}{
		"course_id": opts.CourseID,
	})

	plannerService := api.NewPlannerService(client)

	apiOpts := &api.ListPlannerNotesOptions{
		StartDate: opts.StartDate,
		EndDate:   opts.EndDate,
		CourseID:  opts.CourseID,
	}

	notes, err := plannerService.ListNotes(ctx, apiOpts)
	if err != nil {
		logger.LogCommandError(ctx, "planner.notes.list", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to list planner notes: %w", err)
	}

	if len(notes) == 0 {
		fmt.Println("No planner notes found")
		logger.LogCommandComplete(ctx, "planner.notes.list", 0)
		return nil
	}

	printVerbose("Found %d planner notes:\n\n", len(notes))

	logger.LogCommandComplete(ctx, "planner.notes.list", len(notes))
	return formatOutput(notes, nil)
}

func runPlannerNotesGet(ctx context.Context, client *api.Client, opts *options.PlannerNotesGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "planner.notes.get", map[string]interface{}{
		"note_id": opts.NoteID,
	})

	plannerService := api.NewPlannerService(client)

	note, err := plannerService.GetNote(ctx, opts.NoteID)
	if err != nil {
		logger.LogCommandError(ctx, "planner.notes.get", err, map[string]interface{}{
			"note_id": opts.NoteID,
		})
		return fmt.Errorf("failed to get planner note: %w", err)
	}

	logger.LogCommandComplete(ctx, "planner.notes.get", 1)
	return formatOutput(note, nil)
}

func runPlannerNotesCreate(ctx context.Context, client *api.Client, opts *options.PlannerNotesCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "planner.notes.create", map[string]interface{}{
		"title":     opts.Title,
		"course_id": opts.CourseID,
	})

	plannerService := api.NewPlannerService(client)

	params := &api.CreateNoteParams{
		Title:    opts.Title,
		Details:  opts.Details,
		TodoDate: opts.TodoDate,
		CourseID: opts.CourseID,
	}

	note, err := plannerService.CreateNote(ctx, params)
	if err != nil {
		logger.LogCommandError(ctx, "planner.notes.create", err, map[string]interface{}{
			"title": opts.Title,
		})
		return fmt.Errorf("failed to create planner note: %w", err)
	}

	logger.LogCommandComplete(ctx, "planner.notes.create", 1)
	return formatSuccessOutput(note, "Planner note created successfully!")
}

func runPlannerNotesUpdate(ctx context.Context, client *api.Client, opts *options.PlannerNotesUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "planner.notes.update", map[string]interface{}{
		"note_id": opts.NoteID,
	})

	plannerService := api.NewPlannerService(client)

	params := &api.UpdateNoteParams{}

	if opts.TitleSet {
		params.Title = &opts.Title
	}
	if opts.DetailsSet {
		params.Details = &opts.Details
	}
	if opts.TodoDateSet {
		params.TodoDate = &opts.TodoDate
	}
	if opts.CourseIDSet {
		params.CourseID = &opts.CourseID
	}

	note, err := plannerService.UpdateNote(ctx, opts.NoteID, params)
	if err != nil {
		logger.LogCommandError(ctx, "planner.notes.update", err, map[string]interface{}{
			"note_id": opts.NoteID,
		})
		return fmt.Errorf("failed to update planner note: %w", err)
	}

	logger.LogCommandComplete(ctx, "planner.notes.update", 1)
	return formatSuccessOutput(note, "Planner note updated successfully!")
}

func runPlannerNotesDelete(ctx context.Context, client *api.Client, opts *options.PlannerNotesDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "planner.notes.delete", map[string]interface{}{
		"note_id": opts.NoteID,
		"force":   opts.Force,
	})

	// Confirm deletion
	confirmed, err := confirmDelete("planner note", opts.NoteID, opts.Force)
	if err != nil {
		logger.LogCommandError(ctx, "planner.notes.delete", err, map[string]interface{}{
			"note_id": opts.NoteID,
		})
		return err
	}
	if !confirmed {
		fmt.Println("Delete cancelled")
		logger.LogCommandComplete(ctx, "planner.notes.delete", 0)
		return nil
	}

	plannerService := api.NewPlannerService(client)

	if err := plannerService.DeleteNote(ctx, opts.NoteID); err != nil {
		logger.LogCommandError(ctx, "planner.notes.delete", err, map[string]interface{}{
			"note_id": opts.NoteID,
		})
		return fmt.Errorf("failed to delete planner note: %w", err)
	}

	fmt.Printf("Planner note %d deleted successfully\n", opts.NoteID)

	logger.LogCommandComplete(ctx, "planner.notes.delete", 1)
	return nil
}

func runPlannerComplete(ctx context.Context, client *api.Client, opts *options.PlannerCompleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "planner.complete", map[string]interface{}{
		"plannable_type": opts.PlannableType,
		"plannable_id":   opts.PlannableID,
	})

	plannerService := api.NewPlannerService(client)

	params := &api.CreateOverrideParams{
		PlannableType:  opts.PlannableType,
		PlannableID:    opts.PlannableID,
		MarkedComplete: true,
	}

	override, err := plannerService.CreateOverride(ctx, params)
	if err != nil {
		logger.LogCommandError(ctx, "planner.complete", err, map[string]interface{}{
			"plannable_type": opts.PlannableType,
			"plannable_id":   opts.PlannableID,
		})
		// Check if this is a server error, which often indicates the item doesn't exist in the planner
		if api.IsServerError(err) {
			return fmt.Errorf("failed to mark as complete: %w\n\nThis may occur if:\n  1. The item doesn't appear in your planner (check 'canvas planner items')\n  2. You don't have student enrollment in the course\n  3. The Canvas server is experiencing issues", err)
		}
		return fmt.Errorf("failed to mark as complete: %w", err)
	}

	fmt.Printf("Marked %s %d as complete!\n", override.PlannableType, override.PlannableID)

	logger.LogCommandComplete(ctx, "planner.complete", 1)
	return nil
}

func runPlannerDismiss(ctx context.Context, client *api.Client, opts *options.PlannerDismissOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "planner.dismiss", map[string]interface{}{
		"plannable_type": opts.PlannableType,
		"plannable_id":   opts.PlannableID,
	})

	plannerService := api.NewPlannerService(client)

	params := &api.CreateOverrideParams{
		PlannableType: opts.PlannableType,
		PlannableID:   opts.PlannableID,
		Dismissed:     true,
	}

	override, err := plannerService.CreateOverride(ctx, params)
	if err != nil {
		logger.LogCommandError(ctx, "planner.dismiss", err, map[string]interface{}{
			"plannable_type": opts.PlannableType,
			"plannable_id":   opts.PlannableID,
		})
		// Check if this is a server error, which often indicates the item doesn't exist in the planner
		if api.IsServerError(err) {
			return fmt.Errorf("failed to dismiss: %w\n\nThis may occur if:\n  1. The item doesn't appear in your planner (check 'canvas planner items')\n  2. You don't have student enrollment in the course\n  3. The Canvas server is experiencing issues", err)
		}
		return fmt.Errorf("failed to dismiss: %w", err)
	}

	fmt.Printf("Dismissed %s %d from planner\n", override.PlannableType, override.PlannableID)

	logger.LogCommandComplete(ctx, "planner.dismiss", 1)
	return nil
}

func runPlannerOverrides(ctx context.Context, client *api.Client, opts *options.PlannerOverridesOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "planner.overrides", map[string]interface{}{
		"plannable_type": opts.PlannableType,
	})

	plannerService := api.NewPlannerService(client)

	var apiOpts *api.ListOverridesOptions
	if opts.PlannableType != "" {
		apiOpts = &api.ListOverridesOptions{
			PlannableType: opts.PlannableType,
		}
	}

	overrides, err := plannerService.ListOverrides(ctx, apiOpts)
	if err != nil {
		logger.LogCommandError(ctx, "planner.overrides", err, map[string]interface{}{
			"plannable_type": opts.PlannableType,
		})
		return fmt.Errorf("failed to list overrides: %w", err)
	}

	if len(overrides) == 0 {
		fmt.Println("No planner overrides found")
		logger.LogCommandComplete(ctx, "planner.overrides", 0)
		return nil
	}

	printVerbose("Found %d planner overrides:\n\n", len(overrides))

	logger.LogCommandComplete(ctx, "planner.overrides", len(overrides))
	return formatOutput(overrides, nil)
}
