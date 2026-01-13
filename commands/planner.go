package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var (
	plannerCourseID      int64
	plannerStartDate     string
	plannerEndDate       string
	plannerContextCodes  []string
	plannerFilter        string
	plannerTitle         string
	plannerDetails       string
	plannerTodoDate      string
	plannerPlannableType string
	plannerForce         bool
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

// plannerItemsCmd represents the planner items command
var plannerItemsCmd = &cobra.Command{
	Use:   "items",
	Short: "List planner items",
	Long: `List planner items including assignments, quizzes, and calendar events.

Examples:
  canvas planner items
  canvas planner items --course-id 123
  canvas planner items --start-date 2024-01-01 --end-date 2024-01-31
  canvas planner items --filter all_assignments`,
	RunE: runPlannerItems,
}

// plannerNotesCmd represents the planner notes command group
var plannerNotesCmd = &cobra.Command{
	Use:   "notes",
	Short: "Manage planner notes",
	Long:  `Manage personal planner notes.`,
}

// plannerNotesListCmd lists planner notes
var plannerNotesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List planner notes",
	Long: `List all planner notes.

Examples:
  canvas planner notes list
  canvas planner notes list --course-id 123
  canvas planner notes list --start-date 2024-01-01`,
	RunE: runPlannerNotesList,
}

// plannerNotesGetCmd gets a specific note
var plannerNotesGetCmd = &cobra.Command{
	Use:   "get <note-id>",
	Short: "Get a specific planner note",
	Long: `Get details of a specific planner note.

Examples:
  canvas planner notes get 123`,
	Args: cobra.ExactArgs(1),
	RunE: runPlannerNotesGet,
}

// plannerNotesCreateCmd creates a new note
var plannerNotesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new planner note",
	Long: `Create a new planner note.

Examples:
  canvas planner notes create --title "Study Session"
  canvas planner notes create --title "Project Work" --details "Work on final project" --todo-date 2024-12-15
  canvas planner notes create --title "Review" --course-id 123`,
	RunE: runPlannerNotesCreate,
}

// plannerNotesUpdateCmd updates a note
var plannerNotesUpdateCmd = &cobra.Command{
	Use:   "update <note-id>",
	Short: "Update a planner note",
	Long: `Update an existing planner note.

Examples:
  canvas planner notes update 123 --title "Updated Title"
  canvas planner notes update 123 --todo-date 2024-12-20`,
	Args: cobra.ExactArgs(1),
	RunE: runPlannerNotesUpdate,
}

// plannerNotesDeleteCmd deletes a note
var plannerNotesDeleteCmd = &cobra.Command{
	Use:   "delete <note-id>",
	Short: "Delete a planner note",
	Long: `Delete a planner note.

Examples:
  canvas planner notes delete 123`,
	Args: cobra.ExactArgs(1),
	RunE: runPlannerNotesDelete,
}

// plannerCompleteCmd marks an item as complete
var plannerCompleteCmd = &cobra.Command{
	Use:   "complete <type> <id>",
	Short: "Mark an item as complete",
	Long: `Mark a planner item as complete by creating an override.

Type can be: Assignment, Quiz, DiscussionTopic, WikiPage, CalendarEvent

Examples:
  canvas planner complete Assignment 123
  canvas planner complete Quiz 456`,
	Args: cobra.ExactArgs(2),
	RunE: runPlannerComplete,
}

// plannerDismissCmd dismisses an item from the planner
var plannerDismissCmd = &cobra.Command{
	Use:   "dismiss <type> <id>",
	Short: "Dismiss an item from the planner",
	Long: `Dismiss a planner item so it no longer appears.

Type can be: Assignment, Quiz, DiscussionTopic, WikiPage, CalendarEvent

Examples:
  canvas planner dismiss CalendarEvent 789`,
	Args: cobra.ExactArgs(2),
	RunE: runPlannerDismiss,
}

// plannerOverridesCmd lists overrides
var plannerOverridesCmd = &cobra.Command{
	Use:   "overrides",
	Short: "List planner overrides",
	Long: `List all planner overrides (completed/dismissed items).

Examples:
  canvas planner overrides
  canvas planner overrides --type Assignment`,
	RunE: runPlannerOverrides,
}

func init() {
	rootCmd.AddCommand(plannerCmd)
	plannerCmd.AddCommand(plannerItemsCmd)
	plannerCmd.AddCommand(plannerNotesCmd)
	plannerCmd.AddCommand(plannerCompleteCmd)
	plannerCmd.AddCommand(plannerDismissCmd)
	plannerCmd.AddCommand(plannerOverridesCmd)

	plannerNotesCmd.AddCommand(plannerNotesListCmd)
	plannerNotesCmd.AddCommand(plannerNotesGetCmd)
	plannerNotesCmd.AddCommand(plannerNotesCreateCmd)
	plannerNotesCmd.AddCommand(plannerNotesUpdateCmd)
	plannerNotesCmd.AddCommand(plannerNotesDeleteCmd)

	// Items flags
	plannerItemsCmd.Flags().Int64Var(&plannerCourseID, "course-id", 0, "Filter by course ID")
	plannerItemsCmd.Flags().StringVar(&plannerStartDate, "start-date", "", "Start date (YYYY-MM-DD or ISO 8601)")
	plannerItemsCmd.Flags().StringVar(&plannerEndDate, "end-date", "", "End date (YYYY-MM-DD or ISO 8601)")
	plannerItemsCmd.Flags().StringSliceVar(&plannerContextCodes, "context", []string{}, "Context codes (course_123)")
	plannerItemsCmd.Flags().StringVar(&plannerFilter, "filter", "", "Filter: all_assignments, all_quizzes, all_calendar_events, all_planner_notes")

	// Notes list flags
	plannerNotesListCmd.Flags().Int64Var(&plannerCourseID, "course-id", 0, "Filter by course ID")
	plannerNotesListCmd.Flags().StringVar(&plannerStartDate, "start-date", "", "Start date")
	plannerNotesListCmd.Flags().StringVar(&plannerEndDate, "end-date", "", "End date")

	// Notes create flags
	plannerNotesCreateCmd.Flags().StringVar(&plannerTitle, "title", "", "Note title (required)")
	plannerNotesCreateCmd.Flags().StringVar(&plannerDetails, "details", "", "Note details")
	plannerNotesCreateCmd.Flags().StringVar(&plannerTodoDate, "todo-date", "", "Todo date (ISO 8601)")
	plannerNotesCreateCmd.Flags().Int64Var(&plannerCourseID, "course-id", 0, "Associate with course")
	plannerNotesCreateCmd.MarkFlagRequired("title")

	// Notes update flags
	plannerNotesUpdateCmd.Flags().StringVar(&plannerTitle, "title", "", "New title")
	plannerNotesUpdateCmd.Flags().StringVar(&plannerDetails, "details", "", "New details")
	plannerNotesUpdateCmd.Flags().StringVar(&plannerTodoDate, "todo-date", "", "New todo date")
	plannerNotesUpdateCmd.Flags().Int64Var(&plannerCourseID, "course-id", 0, "New course association")

	// Notes delete flags
	plannerNotesDeleteCmd.Flags().BoolVarP(&plannerForce, "force", "f", false, "Skip confirmation prompt")

	// Overrides flags
	plannerOverridesCmd.Flags().StringVar(&plannerPlannableType, "type", "", "Filter by type (Assignment, Quiz, etc.)")
}

func runPlannerItems(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	plannerService := api.NewPlannerService(client)

	contextCodes := plannerContextCodes
	if plannerCourseID > 0 {
		contextCodes = append(contextCodes, fmt.Sprintf("course_%d", plannerCourseID))
	}

	opts := &api.ListPlannerItemsOptions{
		StartDate:    plannerStartDate,
		EndDate:      plannerEndDate,
		ContextCodes: contextCodes,
		Filter:       plannerFilter,
	}

	ctx := context.Background()
	items, err := plannerService.ListItems(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to list planner items: %w", err)
	}

	if len(items) == 0 {
		fmt.Println("No planner items found")
		return nil
	}

	fmt.Printf("Found %d planner items:\n\n", len(items))

	for _, item := range items {
		displayPlannerItem(&item)
	}

	return nil
}

func runPlannerNotesList(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	plannerService := api.NewPlannerService(client)

	opts := &api.ListPlannerNotesOptions{
		StartDate: plannerStartDate,
		EndDate:   plannerEndDate,
		CourseID:  plannerCourseID,
	}

	ctx := context.Background()
	notes, err := plannerService.ListNotes(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to list planner notes: %w", err)
	}

	if len(notes) == 0 {
		fmt.Println("No planner notes found")
		return nil
	}

	fmt.Printf("Found %d planner notes:\n\n", len(notes))

	for _, note := range notes {
		displayPlannerNote(&note)
	}

	return nil
}

func runPlannerNotesGet(cmd *cobra.Command, args []string) error {
	noteID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid note ID: %s", args[0])
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	plannerService := api.NewPlannerService(client)

	ctx := context.Background()
	note, err := plannerService.GetNote(ctx, noteID)
	if err != nil {
		return fmt.Errorf("failed to get planner note: %w", err)
	}

	displayPlannerNoteFull(note)

	return nil
}

func runPlannerNotesCreate(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	plannerService := api.NewPlannerService(client)

	params := &api.CreateNoteParams{
		Title:    plannerTitle,
		Details:  plannerDetails,
		TodoDate: plannerTodoDate,
		CourseID: plannerCourseID,
	}

	ctx := context.Background()
	note, err := plannerService.CreateNote(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create planner note: %w", err)
	}

	fmt.Println("Planner note created successfully!")
	displayPlannerNote(note)

	return nil
}

func runPlannerNotesUpdate(cmd *cobra.Command, args []string) error {
	noteID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid note ID: %s", args[0])
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	plannerService := api.NewPlannerService(client)

	params := &api.UpdateNoteParams{}

	if cmd.Flags().Changed("title") {
		params.Title = &plannerTitle
	}
	if cmd.Flags().Changed("details") {
		params.Details = &plannerDetails
	}
	if cmd.Flags().Changed("todo-date") {
		params.TodoDate = &plannerTodoDate
	}
	if cmd.Flags().Changed("course-id") {
		params.CourseID = &plannerCourseID
	}

	ctx := context.Background()
	note, err := plannerService.UpdateNote(ctx, noteID, params)
	if err != nil {
		return fmt.Errorf("failed to update planner note: %w", err)
	}

	fmt.Println("Planner note updated successfully!")
	displayPlannerNote(note)

	return nil
}

func runPlannerNotesDelete(cmd *cobra.Command, args []string) error {
	noteID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid note ID: %s", args[0])
	}

	// Confirm deletion
	confirmed, err := confirmDelete("planner note", noteID, plannerForce)
	if err != nil {
		return err
	}
	if !confirmed {
		fmt.Println("Delete cancelled")
		return nil
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	plannerService := api.NewPlannerService(client)

	ctx := context.Background()
	if err := plannerService.DeleteNote(ctx, noteID); err != nil {
		return fmt.Errorf("failed to delete planner note: %w", err)
	}

	fmt.Printf("Planner note %d deleted successfully\n", noteID)
	return nil
}

func runPlannerComplete(cmd *cobra.Command, args []string) error {
	plannableType := args[0]
	plannableID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid ID: %s", args[1])
	}

	// Validate plannable type
	validTypes := []string{"Assignment", "Quiz", "DiscussionTopic", "WikiPage", "CalendarEvent", "PlannerNote", "Announcement"}
	isValidType := false
	for _, t := range validTypes {
		if t == plannableType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return fmt.Errorf("invalid plannable type: %s\nValid types: Assignment, Quiz, DiscussionTopic, WikiPage, CalendarEvent, PlannerNote, Announcement", plannableType)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	plannerService := api.NewPlannerService(client)

	params := &api.CreateOverrideParams{
		PlannableType:  plannableType,
		PlannableID:    plannableID,
		MarkedComplete: true,
	}

	ctx := context.Background()
	override, err := plannerService.CreateOverride(ctx, params)
	if err != nil {
		// Check if this is a server error, which often indicates the item doesn't exist in the planner
		if api.IsServerError(err) {
			return fmt.Errorf("failed to mark as complete: %w\n\nThis may occur if:\n  1. The item doesn't appear in your planner (check 'canvas planner items')\n  2. You don't have student enrollment in the course\n  3. The Canvas server is experiencing issues", err)
		}
		return fmt.Errorf("failed to mark as complete: %w", err)
	}

	fmt.Printf("Marked %s %d as complete!\n", override.PlannableType, override.PlannableID)
	return nil
}

func runPlannerDismiss(cmd *cobra.Command, args []string) error {
	plannableType := args[0]
	plannableID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid ID: %s", args[1])
	}

	// Validate plannable type
	validTypes := []string{"Assignment", "Quiz", "DiscussionTopic", "WikiPage", "CalendarEvent", "PlannerNote", "Announcement"}
	isValidType := false
	for _, t := range validTypes {
		if t == plannableType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return fmt.Errorf("invalid plannable type: %s\nValid types: Assignment, Quiz, DiscussionTopic, WikiPage, CalendarEvent, PlannerNote, Announcement", plannableType)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	plannerService := api.NewPlannerService(client)

	params := &api.CreateOverrideParams{
		PlannableType: plannableType,
		PlannableID:   plannableID,
		Dismissed:     true,
	}

	ctx := context.Background()
	override, err := plannerService.CreateOverride(ctx, params)
	if err != nil {
		// Check if this is a server error, which often indicates the item doesn't exist in the planner
		if api.IsServerError(err) {
			return fmt.Errorf("failed to dismiss: %w\n\nThis may occur if:\n  1. The item doesn't appear in your planner (check 'canvas planner items')\n  2. You don't have student enrollment in the course\n  3. The Canvas server is experiencing issues", err)
		}
		return fmt.Errorf("failed to dismiss: %w", err)
	}

	fmt.Printf("Dismissed %s %d from planner\n", override.PlannableType, override.PlannableID)
	return nil
}

func runPlannerOverrides(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	plannerService := api.NewPlannerService(client)

	var opts *api.ListOverridesOptions
	if plannerPlannableType != "" {
		opts = &api.ListOverridesOptions{
			PlannableType: plannerPlannableType,
		}
	}

	ctx := context.Background()
	overrides, err := plannerService.ListOverrides(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to list overrides: %w", err)
	}

	if len(overrides) == 0 {
		fmt.Println("No planner overrides found")
		return nil
	}

	fmt.Printf("Found %d planner overrides:\n\n", len(overrides))

	for _, override := range overrides {
		displayPlannerOverride(&override)
	}

	return nil
}

func displayPlannerItem(item *api.PlannerItem) {
	typeIcon := "📋"
	switch item.PlannableType {
	case "Assignment":
		typeIcon = "📝"
	case "Quiz":
		typeIcon = "❓"
	case "DiscussionTopic":
		typeIcon = "💬"
	case "CalendarEvent":
		typeIcon = "📅"
	case "PlannerNote":
		typeIcon = "📌"
	}

	fmt.Printf("%s [%s] ID: %d\n", typeIcon, item.PlannableType, item.PlannableID)

	if item.ContextName != "" {
		fmt.Printf("   Context: %s\n", item.ContextName)
	}

	if item.PlannableDate != nil {
		fmt.Printf("   Date: %s\n", item.PlannableDate.Format("2006-01-02 15:04"))
	}

	if item.HTMLURL != "" {
		fmt.Printf("   URL: %s\n", item.HTMLURL)
	}

	fmt.Println()
}

func displayPlannerNote(note *api.PlannerNote) {
	fmt.Printf("📌 [%d] %s\n", note.ID, note.Title)

	if note.TodoDate != nil {
		fmt.Printf("   Due: %s\n", note.TodoDate.Format("2006-01-02 15:04"))
	}

	if note.CourseID != nil {
		fmt.Printf("   Course ID: %d\n", *note.CourseID)
	}

	fmt.Printf("   State: %s\n", note.WorkflowState)

	fmt.Println()
}

func displayPlannerNoteFull(note *api.PlannerNote) {
	displayPlannerNote(note)

	if note.Details != "" {
		fmt.Printf("   Details: %s\n", note.Details)
	}

	fmt.Printf("   Created: %s\n", note.CreatedAt.Format("2006-01-02 15:04"))
	fmt.Printf("   Updated: %s\n", note.UpdatedAt.Format("2006-01-02 15:04"))

	fmt.Println()
}

func displayPlannerOverride(override *api.PlannerOverride) {
	stateIcon := "📋"
	if override.MarkedComplete {
		stateIcon = "✅"
	} else if override.Dismissed {
		stateIcon = "🚫"
	}

	fmt.Printf("%s [%d] %s %d\n", stateIcon, override.ID, override.PlannableType, override.PlannableID)

	if override.MarkedComplete {
		fmt.Printf("   Status: Complete\n")
	}
	if override.Dismissed {
		fmt.Printf("   Status: Dismissed\n")
	}

	fmt.Printf("   Updated: %s\n", override.UpdatedAt.Format("2006-01-02 15:04"))

	fmt.Println()
}
