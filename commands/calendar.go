package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var (
	calendarCourseID        int64
	calendarUserID          int64
	calendarType            string
	calendarStartDate       string
	calendarEndDate         string
	calendarUndated         bool
	calendarAllEvents       bool
	calendarContextCodes    []string
	calendarImportantDates  bool
	calendarTitle           string
	calendarDescription     string
	calendarStartAt         string
	calendarEndAt           string
	calendarLocationName    string
	calendarLocationAddress string
	calendarAllDay          bool
	calendarWhich           string
	calendarCancelReason    string
	calendarForce           bool
)

// calendarCmd represents the calendar command group
var calendarCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Manage Canvas calendar events",
	Long: `Manage Canvas calendar events including listing, viewing, creating, and updating events.

Calendar events can be associated with courses, groups, users, or accounts.
Use context codes to specify which calendars to query.

Examples:
  canvas calendar list --course-id 123
  canvas calendar list --start-date 2024-01-01 --end-date 2024-12-31
  canvas calendar get 456
  canvas calendar create --course-id 123 --title "Team Meeting"`,
}

// calendarListCmd represents the calendar list command
var calendarListCmd = &cobra.Command{
	Use:   "list",
	Short: "List calendar events",
	Long: `List calendar events for courses, groups, or users.

Examples:
  canvas calendar list --course-id 123
  canvas calendar list --context course_123,course_456
  canvas calendar list --start-date 2024-01-01 --end-date 2024-01-31
  canvas calendar list --type assignment
  canvas calendar list --all-events`,
	RunE: runCalendarList,
}

// calendarGetCmd represents the calendar get command
var calendarGetCmd = &cobra.Command{
	Use:   "get <event-id>",
	Short: "Get a specific calendar event",
	Long: `Get details of a specific calendar event.

Examples:
  canvas calendar get 456`,
	Args: cobra.ExactArgs(1),
	RunE: runCalendarGet,
}

// calendarCreateCmd represents the calendar create command
var calendarCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new calendar event",
	Long: `Create a new calendar event.

Examples:
  canvas calendar create --course-id 123 --title "Team Meeting"
  canvas calendar create --course-id 123 --title "Deadline" --start-at "2024-12-01T09:00:00Z" --all-day
  canvas calendar create --course-id 123 --title "Workshop" --start-at "2024-12-01T14:00:00Z" --end-at "2024-12-01T16:00:00Z" --location "Room 101"`,
	RunE: runCalendarCreate,
}

// calendarUpdateCmd represents the calendar update command
var calendarUpdateCmd = &cobra.Command{
	Use:   "update <event-id>",
	Short: "Update an existing calendar event",
	Long: `Update an existing calendar event.

Examples:
  canvas calendar update 456 --title "Updated Meeting"
  canvas calendar update 456 --start-at "2024-12-02T10:00:00Z"
  canvas calendar update 456 --which all --title "Updated Series"`,
	Args: cobra.ExactArgs(1),
	RunE: runCalendarUpdate,
}

// calendarDeleteCmd represents the calendar delete command
var calendarDeleteCmd = &cobra.Command{
	Use:   "delete <event-id>",
	Short: "Delete a calendar event",
	Long: `Delete a calendar event.

Examples:
  canvas calendar delete 456
  canvas calendar delete 456 --reason "Event cancelled"
  canvas calendar delete 456 --which all`,
	Args: cobra.ExactArgs(1),
	RunE: runCalendarDelete,
}

// calendarReserveCmd represents the calendar reserve command
var calendarReserveCmd = &cobra.Command{
	Use:   "reserve <event-id>",
	Short: "Reserve a time slot",
	Long: `Reserve a time slot in an appointment group.

Examples:
  canvas calendar reserve 456`,
	Args: cobra.ExactArgs(1),
	RunE: runCalendarReserve,
}

func init() {
	rootCmd.AddCommand(calendarCmd)
	calendarCmd.AddCommand(calendarListCmd)
	calendarCmd.AddCommand(calendarGetCmd)
	calendarCmd.AddCommand(calendarCreateCmd)
	calendarCmd.AddCommand(calendarUpdateCmd)
	calendarCmd.AddCommand(calendarDeleteCmd)
	calendarCmd.AddCommand(calendarReserveCmd)

	// List flags
	calendarListCmd.Flags().Int64Var(&calendarCourseID, "course-id", 0, "Course ID (adds course context)")
	calendarListCmd.Flags().Int64Var(&calendarUserID, "user-id", 0, "User ID (list events for specific user)")
	calendarListCmd.Flags().StringVar(&calendarType, "type", "", "Event type: event, assignment, sub_assignment")
	calendarListCmd.Flags().StringVar(&calendarStartDate, "start-date", "", "Start date (YYYY-MM-DD or ISO 8601)")
	calendarListCmd.Flags().StringVar(&calendarEndDate, "end-date", "", "End date (YYYY-MM-DD or ISO 8601)")
	calendarListCmd.Flags().BoolVar(&calendarUndated, "undated", false, "Only return undated events")
	calendarListCmd.Flags().BoolVar(&calendarAllEvents, "all-events", false, "Return all events (ignore date filters)")
	calendarListCmd.Flags().StringSliceVar(&calendarContextCodes, "context", []string{}, "Context codes (course_123, user_456)")
	calendarListCmd.Flags().BoolVar(&calendarImportantDates, "important-dates", false, "Only important dates")

	// Get flags (none needed beyond event ID)

	// Create flags
	calendarCreateCmd.Flags().Int64Var(&calendarCourseID, "course-id", 0, "Course ID (required if no --context)")
	calendarCreateCmd.Flags().StringSliceVar(&calendarContextCodes, "context", []string{}, "Context code (course_123, user_456)")
	calendarCreateCmd.Flags().StringVar(&calendarTitle, "title", "", "Event title")
	calendarCreateCmd.Flags().StringVar(&calendarDescription, "description", "", "Event description (HTML)")
	calendarCreateCmd.Flags().StringVar(&calendarStartAt, "start-at", "", "Start time (ISO 8601)")
	calendarCreateCmd.Flags().StringVar(&calendarEndAt, "end-at", "", "End time (ISO 8601)")
	calendarCreateCmd.Flags().StringVar(&calendarLocationName, "location", "", "Location name")
	calendarCreateCmd.Flags().StringVar(&calendarLocationAddress, "address", "", "Location address")
	calendarCreateCmd.Flags().BoolVar(&calendarAllDay, "all-day", false, "All day event")

	// Update flags
	calendarUpdateCmd.Flags().StringVar(&calendarTitle, "title", "", "New event title")
	calendarUpdateCmd.Flags().StringVar(&calendarDescription, "description", "", "New event description")
	calendarUpdateCmd.Flags().StringVar(&calendarStartAt, "start-at", "", "New start time")
	calendarUpdateCmd.Flags().StringVar(&calendarEndAt, "end-at", "", "New end time")
	calendarUpdateCmd.Flags().StringVar(&calendarLocationName, "location", "", "New location name")
	calendarUpdateCmd.Flags().StringVar(&calendarLocationAddress, "address", "", "New location address")
	calendarUpdateCmd.Flags().BoolVar(&calendarAllDay, "all-day", false, "Set as all day event")
	calendarUpdateCmd.Flags().StringVar(&calendarWhich, "which", "", "For series: one, all, following")

	// Delete flags
	calendarDeleteCmd.Flags().StringVar(&calendarCancelReason, "reason", "", "Cancellation reason")
	calendarDeleteCmd.Flags().StringVar(&calendarWhich, "which", "", "For series: one, all, following")
	calendarDeleteCmd.Flags().BoolVarP(&calendarForce, "force", "f", false, "Skip confirmation prompt")

	// Reserve flags (none needed)
}

func runCalendarList(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	calendarService := api.NewCalendarService(client)

	// Build context codes
	contextCodes := calendarContextCodes
	if calendarCourseID > 0 {
		contextCodes = append(contextCodes, fmt.Sprintf("course_%d", calendarCourseID))
	}

	opts := &api.ListCalendarEventsOptions{
		Type:           calendarType,
		StartDate:      calendarStartDate,
		EndDate:        calendarEndDate,
		Undated:        calendarUndated,
		AllEvents:      calendarAllEvents,
		ContextCodes:   contextCodes,
		ImportantDates: calendarImportantDates,
	}

	ctx := context.Background()

	var events []api.CalendarEvent
	if calendarUserID > 0 {
		events, err = calendarService.ListForUser(ctx, calendarUserID, opts)
	} else {
		events, err = calendarService.List(ctx, opts)
	}

	if err != nil {
		return fmt.Errorf("failed to list calendar events: %w", err)
	}

	if len(events) == 0 {
		fmt.Println("No calendar events found")
		return nil
	}

	fmt.Printf("Found %d calendar events:\n\n", len(events))

	for _, event := range events {
		displayCalendarEvent(&event)
	}

	return nil
}

func runCalendarGet(cmd *cobra.Command, args []string) error {
	eventID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid event ID: %s", args[0])
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	calendarService := api.NewCalendarService(client)

	ctx := context.Background()
	event, err := calendarService.Get(ctx, eventID)
	if err != nil {
		return fmt.Errorf("failed to get calendar event: %w", err)
	}

	displayCalendarEventFull(event)

	return nil
}

func runCalendarCreate(cmd *cobra.Command, args []string) error {
	// Determine context code
	var contextCode string
	if len(calendarContextCodes) > 0 {
		contextCode = calendarContextCodes[0]
	} else if calendarCourseID > 0 {
		contextCode = fmt.Sprintf("course_%d", calendarCourseID)
	} else {
		return fmt.Errorf("either --course-id or --context is required")
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	calendarService := api.NewCalendarService(client)

	params := &api.CreateCalendarEventParams{
		ContextCode:     contextCode,
		Title:           calendarTitle,
		Description:     calendarDescription,
		StartAt:         calendarStartAt,
		EndAt:           calendarEndAt,
		LocationName:    calendarLocationName,
		LocationAddress: calendarLocationAddress,
		AllDay:          calendarAllDay,
	}

	ctx := context.Background()
	event, err := calendarService.Create(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create calendar event: %w", err)
	}

	fmt.Println("Calendar event created successfully!")
	displayCalendarEvent(event)

	return nil
}

func runCalendarUpdate(cmd *cobra.Command, args []string) error {
	eventID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid event ID: %s", args[0])
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	calendarService := api.NewCalendarService(client)

	params := &api.UpdateCalendarEventParams{
		Which: calendarWhich,
	}

	if cmd.Flags().Changed("title") {
		params.Title = &calendarTitle
	}
	if cmd.Flags().Changed("description") {
		params.Description = &calendarDescription
	}
	if cmd.Flags().Changed("start-at") {
		params.StartAt = &calendarStartAt
	}
	if cmd.Flags().Changed("end-at") {
		params.EndAt = &calendarEndAt
	}
	if cmd.Flags().Changed("location") {
		params.LocationName = &calendarLocationName
	}
	if cmd.Flags().Changed("address") {
		params.LocationAddress = &calendarLocationAddress
	}
	if cmd.Flags().Changed("all-day") {
		params.AllDay = &calendarAllDay
	}

	ctx := context.Background()
	event, err := calendarService.Update(ctx, eventID, params)
	if err != nil {
		return fmt.Errorf("failed to update calendar event: %w", err)
	}

	fmt.Println("Calendar event updated successfully!")
	displayCalendarEvent(event)

	return nil
}

func runCalendarDelete(cmd *cobra.Command, args []string) error {
	eventID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid event ID: %s", args[0])
	}

	// Confirm deletion
	confirmed, err := confirmDelete("calendar event", eventID, calendarForce)
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

	calendarService := api.NewCalendarService(client)

	var opts *api.DeleteOptions
	if calendarCancelReason != "" || calendarWhich != "" {
		opts = &api.DeleteOptions{
			CancelReason: calendarCancelReason,
			Which:        calendarWhich,
		}
	}

	ctx := context.Background()
	if err := calendarService.Delete(ctx, eventID, opts); err != nil {
		return fmt.Errorf("failed to delete calendar event: %w", err)
	}

	fmt.Printf("Calendar event %d deleted successfully\n", eventID)
	return nil
}

func runCalendarReserve(cmd *cobra.Command, args []string) error {
	eventID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid event ID: %s", args[0])
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	calendarService := api.NewCalendarService(client)

	ctx := context.Background()
	event, err := calendarService.Reserve(ctx, eventID, nil, "", false)
	if err != nil {
		return fmt.Errorf("failed to reserve time slot: %w", err)
	}

	fmt.Println("Time slot reserved successfully!")
	displayCalendarEvent(event)

	return nil
}

func displayCalendarEvent(event *api.CalendarEvent) {
	stateIcon := "📅"
	if event.AllDay {
		stateIcon = "📆"
	}
	if event.WorkflowState == "locked" {
		stateIcon = "🔒"
	}

	fmt.Printf("%s [%d] %s\n", stateIcon, event.ID, event.Title)

	if event.StartAt != nil {
		if event.AllDay {
			fmt.Printf("   Date: %s (All Day)\n", event.StartAt.Format("2006-01-02"))
		} else if event.EndAt != nil {
			fmt.Printf("   Time: %s - %s\n",
				event.StartAt.Format("2006-01-02 15:04"),
				event.EndAt.Format("15:04"))
		} else {
			fmt.Printf("   Time: %s\n", event.StartAt.Format("2006-01-02 15:04"))
		}
	}

	if event.LocationName != "" {
		fmt.Printf("   Location: %s\n", event.LocationName)
	}

	fmt.Printf("   Context: %s\n", event.ContextCode)

	fmt.Println()
}

func displayCalendarEventFull(event *api.CalendarEvent) {
	displayCalendarEvent(event)

	fmt.Printf("   State: %s\n", event.WorkflowState)

	if event.LocationAddress != "" {
		fmt.Printf("   Address: %s\n", event.LocationAddress)
	}

	if event.ContextName != "" {
		fmt.Printf("   Context Name: %s\n", event.ContextName)
	}

	if event.SeriesNaturalLang != "" {
		fmt.Printf("   Recurrence: %s\n", event.SeriesNaturalLang)
	}

	if event.Description != "" {
		fmt.Printf("\nDescription:\n")
		description := event.Description
		if len(description) > 500 {
			description = description[:500] + "..."
		}
		description = stripHTMLTags(description)
		fmt.Println(description)
	}

	fmt.Println()
}
