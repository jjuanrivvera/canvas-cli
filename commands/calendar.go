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

func init() {
	rootCmd.AddCommand(calendarCmd)
	calendarCmd.AddCommand(newCalendarListCmd())
	calendarCmd.AddCommand(newCalendarGetCmd())
	calendarCmd.AddCommand(newCalendarCreateCmd())
	calendarCmd.AddCommand(newCalendarUpdateCmd())
	calendarCmd.AddCommand(newCalendarDeleteCmd())
	calendarCmd.AddCommand(newCalendarReserveCmd())
}

func newCalendarListCmd() *cobra.Command {
	opts := &options.CalendarListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List calendar events",
		Long: `List calendar events for courses, groups, or users.

Examples:
  canvas calendar list --course-id 123
  canvas calendar list --context course_123,course_456
  canvas calendar list --start-date 2024-01-01 --end-date 2024-01-31
  canvas calendar list --type assignment
  canvas calendar list --all-events`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runCalendarList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (adds course context)")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (list events for specific user)")
	cmd.Flags().StringVar(&opts.Type, "type", "", "Event type: event, assignment, sub_assignment")
	cmd.Flags().StringVar(&opts.StartDate, "start-date", "", "Start date (YYYY-MM-DD or ISO 8601)")
	cmd.Flags().StringVar(&opts.EndDate, "end-date", "", "End date (YYYY-MM-DD or ISO 8601)")
	cmd.Flags().BoolVar(&opts.Undated, "undated", false, "Only return undated events")
	cmd.Flags().BoolVar(&opts.AllEvents, "all-events", false, "Return all events (ignore date filters)")
	cmd.Flags().StringSliceVar(&opts.ContextCodes, "context", []string{}, "Context codes (course_123, user_456)")
	cmd.Flags().BoolVar(&opts.ImportantDates, "important-dates", false, "Only important dates")

	return cmd
}

func newCalendarGetCmd() *cobra.Command {
	opts := &options.CalendarGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <event-id>",
		Short: "Get a specific calendar event",
		Long: `Get details of a specific calendar event.

Examples:
  canvas calendar get 456`,
		Args: ExactArgsWithUsage(1, "event-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			eventID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid event ID: %s", args[0])
			}
			opts.EventID = eventID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runCalendarGet(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newCalendarCreateCmd() *cobra.Command {
	opts := &options.CalendarCreateOptions{}
	var courseID int64
	var contextCodes []string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new calendar event",
		Long: `Create a new calendar event.

Examples:
  canvas calendar create --course-id 123 --title "Team Meeting"
  canvas calendar create --course-id 123 --title "Deadline" --start-at "2024-12-01T09:00:00Z" --all-day
  canvas calendar create --course-id 123 --title "Workshop" --start-at "2024-12-01T14:00:00Z" --end-at "2024-12-01T16:00:00Z" --location "Room 101"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Determine context code
			if len(contextCodes) > 0 {
				opts.ContextCode = contextCodes[0]
			} else if courseID > 0 {
				opts.ContextCode = fmt.Sprintf("course_%d", courseID)
			} else {
				return fmt.Errorf("either --course-id or --context is required")
			}

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runCalendarCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&courseID, "course-id", 0, "Course ID (required if no --context)")
	cmd.Flags().StringSliceVar(&contextCodes, "context", []string{}, "Context code (course_123, user_456)")
	cmd.Flags().StringVar(&opts.Title, "title", "", "Event title")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Event description (HTML)")
	cmd.Flags().StringVar(&opts.StartAt, "start-at", "", "Start time (ISO 8601)")
	cmd.Flags().StringVar(&opts.EndAt, "end-at", "", "End time (ISO 8601)")
	cmd.Flags().StringVar(&opts.LocationName, "location", "", "Location name")
	cmd.Flags().StringVar(&opts.LocationAddress, "address", "", "Location address")
	cmd.Flags().BoolVar(&opts.AllDay, "all-day", false, "All day event")

	return cmd
}

func newCalendarUpdateCmd() *cobra.Command {
	opts := &options.CalendarUpdateOptions{}
	var which string

	cmd := &cobra.Command{
		Use:   "update <event-id>",
		Short: "Update an existing calendar event",
		Long: `Update an existing calendar event.

Examples:
  canvas calendar update 456 --title "Updated Meeting"
  canvas calendar update 456 --start-at "2024-12-02T10:00:00Z"
  canvas calendar update 456 --which all --title "Updated Series"`,
		Args: ExactArgsWithUsage(1, "event-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			eventID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid event ID: %s", args[0])
			}
			opts.EventID = eventID

			// Track which fields were set
			opts.TitleSet = cmd.Flags().Changed("title")
			opts.DescriptionSet = cmd.Flags().Changed("description")
			opts.StartAtSet = cmd.Flags().Changed("start-at")
			opts.EndAtSet = cmd.Flags().Changed("end-at")
			opts.LocationNameSet = cmd.Flags().Changed("location")
			opts.LocationAddressSet = cmd.Flags().Changed("address")
			opts.AllDaySet = cmd.Flags().Changed("all-day")

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runCalendarUpdate(cmd.Context(), client, opts, which)
		},
	}

	cmd.Flags().StringVar(&opts.Title, "title", "", "New event title")
	cmd.Flags().StringVar(&opts.Description, "description", "", "New event description")
	cmd.Flags().StringVar(&opts.StartAt, "start-at", "", "New start time")
	cmd.Flags().StringVar(&opts.EndAt, "end-at", "", "New end time")
	cmd.Flags().StringVar(&opts.LocationName, "location", "", "New location name")
	cmd.Flags().StringVar(&opts.LocationAddress, "address", "", "New location address")
	cmd.Flags().BoolVar(&opts.AllDay, "all-day", false, "Set as all day event")
	cmd.Flags().StringVar(&which, "which", "", "For series: one, all, following")

	return cmd
}

func newCalendarDeleteCmd() *cobra.Command {
	opts := &options.CalendarDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <event-id>",
		Short: "Delete a calendar event",
		Long: `Delete a calendar event.

Examples:
  canvas calendar delete 456
  canvas calendar delete 456 --reason "Event cancelled"
  canvas calendar delete 456 --which all`,
		Args: ExactArgsWithUsage(1, "event-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			eventID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid event ID: %s", args[0])
			}
			opts.EventID = eventID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runCalendarDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.CancelReason, "reason", "", "Cancellation reason")
	cmd.Flags().StringVar(&opts.Which, "which", "", "For series: one, all, following")
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Skip confirmation prompt")

	return cmd
}

func newCalendarReserveCmd() *cobra.Command {
	opts := &options.CalendarReserveOptions{}

	cmd := &cobra.Command{
		Use:   "reserve <event-id>",
		Short: "Reserve a time slot",
		Long: `Reserve a time slot in an appointment group.

Examples:
  canvas calendar reserve 456`,
		Args: ExactArgsWithUsage(1, "event-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			eventID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid event ID: %s", args[0])
			}
			opts.EventID = eventID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runCalendarReserve(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func runCalendarList(ctx context.Context, client *api.Client, opts *options.CalendarListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "calendar.list", map[string]interface{}{
		"course_id":       opts.CourseID,
		"user_id":         opts.UserID,
		"type":            opts.Type,
		"start_date":      opts.StartDate,
		"end_date":        opts.EndDate,
		"undated":         opts.Undated,
		"all_events":      opts.AllEvents,
		"important_dates": opts.ImportantDates,
	})

	calendarService := api.NewCalendarService(client)

	// Build context codes
	contextCodes := opts.ContextCodes
	if opts.CourseID > 0 {
		contextCodes = append(contextCodes, fmt.Sprintf("course_%d", opts.CourseID))
	}

	apiOpts := &api.ListCalendarEventsOptions{
		Type:           opts.Type,
		StartDate:      opts.StartDate,
		EndDate:        opts.EndDate,
		Undated:        opts.Undated,
		AllEvents:      opts.AllEvents,
		ContextCodes:   contextCodes,
		ImportantDates: opts.ImportantDates,
	}

	var events []api.CalendarEvent
	var err error
	if opts.UserID > 0 {
		events, err = calendarService.ListForUser(ctx, opts.UserID, apiOpts)
	} else {
		events, err = calendarService.List(ctx, apiOpts)
	}

	if err != nil {
		logger.LogCommandError(ctx, "calendar.list", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"user_id":   opts.UserID,
		})
		return fmt.Errorf("failed to list calendar events: %w", err)
	}

	if len(events) == 0 {
		fmt.Println("No calendar events found")
		logger.LogCommandComplete(ctx, "calendar.list", 0)
		return nil
	}

	printVerbose("Found %d calendar events:\n\n", len(events))
	logger.LogCommandComplete(ctx, "calendar.list", len(events))
	return formatOutput(events, nil)
}

func runCalendarGet(ctx context.Context, client *api.Client, opts *options.CalendarGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "calendar.get", map[string]interface{}{
		"event_id": opts.EventID,
	})

	calendarService := api.NewCalendarService(client)

	event, err := calendarService.Get(ctx, opts.EventID)
	if err != nil {
		logger.LogCommandError(ctx, "calendar.get", err, map[string]interface{}{
			"event_id": opts.EventID,
		})
		return fmt.Errorf("failed to get calendar event: %w", err)
	}

	logger.LogCommandComplete(ctx, "calendar.get", 1)
	return formatOutput(event, nil)
}

func runCalendarCreate(ctx context.Context, client *api.Client, opts *options.CalendarCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "calendar.create", map[string]interface{}{
		"context_code": opts.ContextCode,
		"title":        opts.Title,
		"start_at":     opts.StartAt,
		"end_at":       opts.EndAt,
		"all_day":      opts.AllDay,
	})

	calendarService := api.NewCalendarService(client)

	params := &api.CreateCalendarEventParams{
		ContextCode:     opts.ContextCode,
		Title:           opts.Title,
		Description:     opts.Description,
		StartAt:         opts.StartAt,
		EndAt:           opts.EndAt,
		LocationName:    opts.LocationName,
		LocationAddress: opts.LocationAddress,
		AllDay:          opts.AllDay,
	}

	event, err := calendarService.Create(ctx, params)
	if err != nil {
		logger.LogCommandError(ctx, "calendar.create", err, map[string]interface{}{
			"context_code": opts.ContextCode,
			"title":        opts.Title,
		})
		return fmt.Errorf("failed to create calendar event: %w", err)
	}

	logger.LogCommandComplete(ctx, "calendar.create", 1)
	return formatSuccessOutput(event, "Calendar event created successfully!")
}

func runCalendarUpdate(ctx context.Context, client *api.Client, opts *options.CalendarUpdateOptions, which string) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "calendar.update", map[string]interface{}{
		"event_id": opts.EventID,
		"which":    which,
	})

	calendarService := api.NewCalendarService(client)

	params := &api.UpdateCalendarEventParams{
		Which: which,
	}

	if opts.TitleSet {
		params.Title = &opts.Title
	}
	if opts.DescriptionSet {
		params.Description = &opts.Description
	}
	if opts.StartAtSet {
		params.StartAt = &opts.StartAt
	}
	if opts.EndAtSet {
		params.EndAt = &opts.EndAt
	}
	if opts.LocationNameSet {
		params.LocationName = &opts.LocationName
	}
	if opts.LocationAddressSet {
		params.LocationAddress = &opts.LocationAddress
	}
	if opts.AllDaySet {
		params.AllDay = &opts.AllDay
	}

	event, err := calendarService.Update(ctx, opts.EventID, params)
	if err != nil {
		logger.LogCommandError(ctx, "calendar.update", err, map[string]interface{}{
			"event_id": opts.EventID,
		})
		return fmt.Errorf("failed to update calendar event: %w", err)
	}

	logger.LogCommandComplete(ctx, "calendar.update", 1)
	return formatSuccessOutput(event, "Calendar event updated successfully!")
}

func runCalendarDelete(ctx context.Context, client *api.Client, opts *options.CalendarDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "calendar.delete", map[string]interface{}{
		"event_id": opts.EventID,
		"which":    opts.Which,
		"force":    opts.Force,
	})

	// Confirm deletion
	confirmed, err := confirmDelete("calendar event", opts.EventID, opts.Force)
	if err != nil {
		logger.LogCommandError(ctx, "calendar.delete", err, map[string]interface{}{
			"event_id": opts.EventID,
		})
		return err
	}
	if !confirmed {
		fmt.Println("Delete cancelled")
		logger.LogCommandComplete(ctx, "calendar.delete", 0)
		return nil
	}

	calendarService := api.NewCalendarService(client)

	var deleteOpts *api.DeleteOptions
	if opts.CancelReason != "" || opts.Which != "" {
		deleteOpts = &api.DeleteOptions{
			CancelReason: opts.CancelReason,
			Which:        opts.Which,
		}
	}

	if err := calendarService.Delete(ctx, opts.EventID, deleteOpts); err != nil {
		logger.LogCommandError(ctx, "calendar.delete", err, map[string]interface{}{
			"event_id": opts.EventID,
		})
		return fmt.Errorf("failed to delete calendar event: %w", err)
	}

	fmt.Printf("Calendar event %d deleted successfully\n", opts.EventID)
	logger.LogCommandComplete(ctx, "calendar.delete", 1)
	return nil
}

func runCalendarReserve(ctx context.Context, client *api.Client, opts *options.CalendarReserveOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "calendar.reserve", map[string]interface{}{
		"event_id": opts.EventID,
	})

	calendarService := api.NewCalendarService(client)

	event, err := calendarService.Reserve(ctx, opts.EventID, nil, "", false)
	if err != nil {
		logger.LogCommandError(ctx, "calendar.reserve", err, map[string]interface{}{
			"event_id": opts.EventID,
		})
		return fmt.Errorf("failed to reserve time slot: %w", err)
	}

	logger.LogCommandComplete(ctx, "calendar.reserve", 1)
	return formatSuccessOutput(event, "Time slot reserved successfully!")
}
