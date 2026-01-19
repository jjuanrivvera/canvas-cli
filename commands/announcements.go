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

// announcementsCmd represents the announcements command group
var announcementsCmd = &cobra.Command{
	Use:   "announcements",
	Short: "Manage Canvas announcements",
	Long: `Manage Canvas announcements including listing, viewing, and creating announcements.

Announcements are a special type of discussion topic that appear in the
announcements section of a course. They are used for important course updates.

Examples:
  canvas announcements list --course-id 123
  canvas announcements get --course-id 123 456
  canvas announcements create --course-id 123 --title "Welcome!"`,
}

func init() {
	rootCmd.AddCommand(announcementsCmd)
	announcementsCmd.AddCommand(newAnnouncementsListCmd())
	announcementsCmd.AddCommand(newAnnouncementsGetCmd())
	announcementsCmd.AddCommand(newAnnouncementsCreateCmd())
	announcementsCmd.AddCommand(newAnnouncementsUpdateCmd())
	announcementsCmd.AddCommand(newAnnouncementsDeleteCmd())
}

func newAnnouncementsListCmd() *cobra.Command {
	opts := &options.AnnouncementsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List announcements",
		Long: `List announcements for a course or multiple courses.

Examples:
  canvas announcements list --course-id 123
  canvas announcements list --course-id 123 --active-only
  canvas announcements list --course-id 123 --start-date 2024-01-01 --end-date 2024-12-31`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAnnouncementsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.StartDate, "start-date", "", "Start date (YYYY-MM-DD or ISO 8601)")
	cmd.Flags().StringVar(&opts.EndDate, "end-date", "", "End date (YYYY-MM-DD or ISO 8601)")
	cmd.Flags().BoolVar(&opts.ActiveOnly, "active-only", false, "Only return active announcements")
	cmd.Flags().BoolVar(&opts.LatestOnly, "latest-only", false, "Only return the latest announcement per context")
	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Include: sections, sections_user_count")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newAnnouncementsGetCmd() *cobra.Command {
	opts := &options.AnnouncementsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <announcement-id>",
		Short: "Get a specific announcement",
		Long: `Get details of a specific announcement.

Examples:
  canvas announcements get --course-id 123 456`,
		Args: ExactArgsWithUsage(1, "announcement-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			announcementID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid announcement ID: %s", args[0])
			}
			opts.AnnouncementID = announcementID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAnnouncementsGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newAnnouncementsCreateCmd() *cobra.Command {
	opts := &options.AnnouncementsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new announcement",
		Long: `Create a new announcement in a course.

Examples:
  canvas announcements create --course-id 123 --title "Welcome to the Course!"
  canvas announcements create --course-id 123 --title "Important Update" --message "<p>Please read...</p>"
  canvas announcements create --course-id 123 --title "Scheduled" --delayed-at "2024-12-01T09:00:00Z"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAnnouncementsCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.Title, "title", "", "Announcement title (required)")
	cmd.Flags().StringVar(&opts.Message, "message", "", "Announcement message (HTML)")
	cmd.Flags().StringVar(&opts.DelayedAt, "delayed-at", "", "Delay posting until (ISO 8601)")
	cmd.Flags().BoolVar(&opts.Published, "published", true, "Publish the announcement (default: true)")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("title")

	return cmd
}

func newAnnouncementsUpdateCmd() *cobra.Command {
	opts := &options.AnnouncementsUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <announcement-id>",
		Short: "Update an existing announcement",
		Long: `Update an existing announcement.

Examples:
  canvas announcements update --course-id 123 456 --title "Updated Title"
  canvas announcements update --course-id 123 456 --message "<p>Updated content</p>"`,
		Args: ExactArgsWithUsage(1, "announcement-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			announcementID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid announcement ID: %s", args[0])
			}
			opts.AnnouncementID = announcementID

			// Track which fields were set
			opts.TitleSet = cmd.Flags().Changed("title")
			opts.MessageSet = cmd.Flags().Changed("message")
			opts.DelayedAtSet = cmd.Flags().Changed("delayed-at")

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAnnouncementsUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.Title, "title", "", "New announcement title")
	cmd.Flags().StringVar(&opts.Message, "message", "", "New announcement message")
	cmd.Flags().StringVar(&opts.DelayedAt, "delayed-at", "", "Delay posting until")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newAnnouncementsDeleteCmd() *cobra.Command {
	opts := &options.AnnouncementsDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <announcement-id>",
		Short: "Delete an announcement",
		Long: `Delete an announcement from a course.

Examples:
  canvas announcements delete --course-id 123 456`,
		Args: ExactArgsWithUsage(1, "announcement-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			announcementID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid announcement ID: %s", args[0])
			}
			opts.AnnouncementID = announcementID

			if err := opts.Validate(); err != nil {
				return err
			}

			// Confirm deletion
			confirmed, err := confirmDelete("announcement", opts.AnnouncementID, opts.Force)
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

			return runAnnouncementsDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Skip confirmation prompt")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func runAnnouncementsList(ctx context.Context, client *api.Client, opts *options.AnnouncementsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "announcements.list", map[string]interface{}{
		"course_id": opts.CourseID,
	})

	announcementsService := api.NewAnnouncementsService(client)

	apiOpts := &api.ListAnnouncementsOptions{
		ContextCodes: []string{fmt.Sprintf("course_%d", opts.CourseID)},
		StartDate:    opts.StartDate,
		EndDate:      opts.EndDate,
		ActiveOnly:   opts.ActiveOnly,
		LatestOnly:   opts.LatestOnly,
		Include:      opts.Include,
	}

	announcements, err := announcementsService.List(ctx, apiOpts)
	if err != nil {
		logger.LogCommandError(ctx, "announcements.list", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to list announcements: %w", err)
	}

	if len(announcements) == 0 {
		fmt.Println("No announcements found")
		logger.LogCommandComplete(ctx, "announcements.list", 0)
		return nil
	}

	printVerbose("Found %d announcements:\n\n", len(announcements))
	logger.LogCommandComplete(ctx, "announcements.list", len(announcements))
	return formatOutput(announcements, nil)
}

func runAnnouncementsGet(ctx context.Context, client *api.Client, opts *options.AnnouncementsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "announcements.get", map[string]interface{}{
		"course_id":       opts.CourseID,
		"announcement_id": opts.AnnouncementID,
	})

	// Announcements are discussion topics, so we use the discussions service
	discussionsService := api.NewDiscussionsService(client)

	announcement, err := discussionsService.Get(ctx, opts.CourseID, opts.AnnouncementID, nil)
	if err != nil {
		logger.LogCommandError(ctx, "announcements.get", err, map[string]interface{}{
			"course_id":       opts.CourseID,
			"announcement_id": opts.AnnouncementID,
		})
		return fmt.Errorf("failed to get announcement: %w", err)
	}

	logger.LogCommandComplete(ctx, "announcements.get", 1)
	return formatOutput(announcement, nil)
}

func runAnnouncementsCreate(ctx context.Context, client *api.Client, opts *options.AnnouncementsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "announcements.create", map[string]interface{}{
		"course_id": opts.CourseID,
		"title":     opts.Title,
	})

	discussionsService := api.NewDiscussionsService(client)

	params := &api.CreateDiscussionParams{
		Title:          opts.Title,
		Message:        opts.Message,
		Published:      opts.Published,
		DelayedPostAt:  opts.DelayedAt,
		IsAnnouncement: true,
	}

	announcement, err := discussionsService.Create(ctx, opts.CourseID, params)
	if err != nil {
		logger.LogCommandError(ctx, "announcements.create", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"title":     opts.Title,
		})
		return fmt.Errorf("failed to create announcement: %w", err)
	}

	logger.LogCommandComplete(ctx, "announcements.create", 1)
	return formatSuccessOutput(announcement, "Announcement created successfully!")
}

func runAnnouncementsUpdate(ctx context.Context, client *api.Client, opts *options.AnnouncementsUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "announcements.update", map[string]interface{}{
		"course_id":       opts.CourseID,
		"announcement_id": opts.AnnouncementID,
	})

	discussionsService := api.NewDiscussionsService(client)

	params := &api.UpdateDiscussionParams{}

	if opts.TitleSet {
		params.Title = &opts.Title
	}
	if opts.MessageSet {
		params.Message = &opts.Message
	}
	if opts.DelayedAtSet {
		params.DelayedPostAt = &opts.DelayedAt
	}

	announcement, err := discussionsService.Update(ctx, opts.CourseID, opts.AnnouncementID, params)
	if err != nil {
		logger.LogCommandError(ctx, "announcements.update", err, map[string]interface{}{
			"course_id":       opts.CourseID,
			"announcement_id": opts.AnnouncementID,
		})
		return fmt.Errorf("failed to update announcement: %w", err)
	}

	logger.LogCommandComplete(ctx, "announcements.update", 1)
	return formatSuccessOutput(announcement, "Announcement updated successfully!")
}

func runAnnouncementsDelete(ctx context.Context, client *api.Client, opts *options.AnnouncementsDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "announcements.delete", map[string]interface{}{
		"course_id":       opts.CourseID,
		"announcement_id": opts.AnnouncementID,
	})

	discussionsService := api.NewDiscussionsService(client)

	if err := discussionsService.Delete(ctx, opts.CourseID, opts.AnnouncementID); err != nil {
		logger.LogCommandError(ctx, "announcements.delete", err, map[string]interface{}{
			"course_id":       opts.CourseID,
			"announcement_id": opts.AnnouncementID,
		})
		return fmt.Errorf("failed to delete announcement: %w", err)
	}

	fmt.Printf("Announcement %d deleted successfully\n", opts.AnnouncementID)
	logger.LogCommandComplete(ctx, "announcements.delete", 1)
	return nil
}
