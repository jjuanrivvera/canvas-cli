package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var (
	announcementsCourseID   int64
	announcementsStartDate  string
	announcementsEndDate    string
	announcementsActiveOnly bool
	announcementsLatestOnly bool
	announcementsInclude    []string
	announcementsTitle      string
	announcementsMessage    string
	announcementsDelayedAt  string
	announcementsPublished  bool
	announcementsForce      bool
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

// announcementsListCmd represents the announcements list command
var announcementsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List announcements",
	Long: `List announcements for a course or multiple courses.

Examples:
  canvas announcements list --course-id 123
  canvas announcements list --course-id 123 --active-only
  canvas announcements list --course-id 123 --start-date 2024-01-01 --end-date 2024-12-31`,
	RunE: runAnnouncementsList,
}

// announcementsGetCmd represents the announcements get command
var announcementsGetCmd = &cobra.Command{
	Use:   "get <announcement-id>",
	Short: "Get a specific announcement",
	Long: `Get details of a specific announcement.

Examples:
  canvas announcements get --course-id 123 456`,
	Args: cobra.ExactArgs(1),
	RunE: runAnnouncementsGet,
}

// announcementsCreateCmd represents the announcements create command
var announcementsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new announcement",
	Long: `Create a new announcement in a course.

Examples:
  canvas announcements create --course-id 123 --title "Welcome to the Course!"
  canvas announcements create --course-id 123 --title "Important Update" --message "<p>Please read...</p>"
  canvas announcements create --course-id 123 --title "Scheduled" --delayed-at "2024-12-01T09:00:00Z"`,
	RunE: runAnnouncementsCreate,
}

// announcementsUpdateCmd represents the announcements update command
var announcementsUpdateCmd = &cobra.Command{
	Use:   "update <announcement-id>",
	Short: "Update an existing announcement",
	Long: `Update an existing announcement.

Examples:
  canvas announcements update --course-id 123 456 --title "Updated Title"
  canvas announcements update --course-id 123 456 --message "<p>Updated content</p>"`,
	Args: cobra.ExactArgs(1),
	RunE: runAnnouncementsUpdate,
}

// announcementsDeleteCmd represents the announcements delete command
var announcementsDeleteCmd = &cobra.Command{
	Use:   "delete <announcement-id>",
	Short: "Delete an announcement",
	Long: `Delete an announcement from a course.

Examples:
  canvas announcements delete --course-id 123 456`,
	Args: cobra.ExactArgs(1),
	RunE: runAnnouncementsDelete,
}

func init() {
	rootCmd.AddCommand(announcementsCmd)
	announcementsCmd.AddCommand(announcementsListCmd)
	announcementsCmd.AddCommand(announcementsGetCmd)
	announcementsCmd.AddCommand(announcementsCreateCmd)
	announcementsCmd.AddCommand(announcementsUpdateCmd)
	announcementsCmd.AddCommand(announcementsDeleteCmd)

	// List flags
	announcementsListCmd.Flags().Int64Var(&announcementsCourseID, "course-id", 0, "Course ID (required)")
	announcementsListCmd.Flags().StringVar(&announcementsStartDate, "start-date", "", "Start date (YYYY-MM-DD or ISO 8601)")
	announcementsListCmd.Flags().StringVar(&announcementsEndDate, "end-date", "", "End date (YYYY-MM-DD or ISO 8601)")
	announcementsListCmd.Flags().BoolVar(&announcementsActiveOnly, "active-only", false, "Only return active announcements")
	announcementsListCmd.Flags().BoolVar(&announcementsLatestOnly, "latest-only", false, "Only return the latest announcement per context")
	announcementsListCmd.Flags().StringSliceVar(&announcementsInclude, "include", []string{}, "Include: sections, sections_user_count")
	announcementsListCmd.MarkFlagRequired("course-id")

	// Get flags
	announcementsGetCmd.Flags().Int64Var(&announcementsCourseID, "course-id", 0, "Course ID (required)")
	announcementsGetCmd.MarkFlagRequired("course-id")

	// Create flags
	announcementsCreateCmd.Flags().Int64Var(&announcementsCourseID, "course-id", 0, "Course ID (required)")
	announcementsCreateCmd.Flags().StringVar(&announcementsTitle, "title", "", "Announcement title (required)")
	announcementsCreateCmd.Flags().StringVar(&announcementsMessage, "message", "", "Announcement message (HTML)")
	announcementsCreateCmd.Flags().StringVar(&announcementsDelayedAt, "delayed-at", "", "Delay posting until (ISO 8601)")
	announcementsCreateCmd.Flags().BoolVar(&announcementsPublished, "published", true, "Publish the announcement (default: true)")
	announcementsCreateCmd.MarkFlagRequired("course-id")
	announcementsCreateCmd.MarkFlagRequired("title")

	// Update flags
	announcementsUpdateCmd.Flags().Int64Var(&announcementsCourseID, "course-id", 0, "Course ID (required)")
	announcementsUpdateCmd.Flags().StringVar(&announcementsTitle, "title", "", "New announcement title")
	announcementsUpdateCmd.Flags().StringVar(&announcementsMessage, "message", "", "New announcement message")
	announcementsUpdateCmd.Flags().StringVar(&announcementsDelayedAt, "delayed-at", "", "Delay posting until")
	announcementsUpdateCmd.MarkFlagRequired("course-id")

	// Delete flags
	announcementsDeleteCmd.Flags().Int64Var(&announcementsCourseID, "course-id", 0, "Course ID (required)")
	announcementsDeleteCmd.Flags().BoolVarP(&announcementsForce, "force", "f", false, "Skip confirmation prompt")
	announcementsDeleteCmd.MarkFlagRequired("course-id")
}

func runAnnouncementsList(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	announcementsService := api.NewAnnouncementsService(client)

	opts := &api.ListAnnouncementsOptions{
		ContextCodes: []string{fmt.Sprintf("course_%d", announcementsCourseID)},
		StartDate:    announcementsStartDate,
		EndDate:      announcementsEndDate,
		ActiveOnly:   announcementsActiveOnly,
		LatestOnly:   announcementsLatestOnly,
		Include:      announcementsInclude,
	}

	ctx := context.Background()
	announcements, err := announcementsService.List(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to list announcements: %w", err)
	}

	if len(announcements) == 0 {
		fmt.Println("No announcements found")
		return nil
	}

	printVerbose("Found %d announcements:\n\n", len(announcements))
	return formatOutput(announcements, nil)
}

func runAnnouncementsGet(cmd *cobra.Command, args []string) error {
	announcementID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid announcement ID: %s", args[0])
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Announcements are discussion topics, so we use the discussions service
	discussionsService := api.NewDiscussionsService(client)

	ctx := context.Background()
	announcement, err := discussionsService.Get(ctx, announcementsCourseID, announcementID, nil)
	if err != nil {
		return fmt.Errorf("failed to get announcement: %w", err)
	}

	return formatOutput(announcement, nil)
}

func runAnnouncementsCreate(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	discussionsService := api.NewDiscussionsService(client)

	params := &api.CreateDiscussionParams{
		Title:          announcementsTitle,
		Message:        announcementsMessage,
		Published:      announcementsPublished,
		DelayedPostAt:  announcementsDelayedAt,
		IsAnnouncement: true,
	}

	ctx := context.Background()
	announcement, err := discussionsService.Create(ctx, announcementsCourseID, params)
	if err != nil {
		return fmt.Errorf("failed to create announcement: %w", err)
	}

	fmt.Println("Announcement created successfully!")
	displayAnnouncement(announcement)

	return nil
}

func runAnnouncementsUpdate(cmd *cobra.Command, args []string) error {
	announcementID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid announcement ID: %s", args[0])
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	discussionsService := api.NewDiscussionsService(client)

	params := &api.UpdateDiscussionParams{}

	if cmd.Flags().Changed("title") {
		params.Title = &announcementsTitle
	}
	if cmd.Flags().Changed("message") {
		params.Message = &announcementsMessage
	}
	if cmd.Flags().Changed("delayed-at") {
		params.DelayedPostAt = &announcementsDelayedAt
	}

	ctx := context.Background()
	announcement, err := discussionsService.Update(ctx, announcementsCourseID, announcementID, params)
	if err != nil {
		return fmt.Errorf("failed to update announcement: %w", err)
	}

	fmt.Println("Announcement updated successfully!")
	displayAnnouncement(announcement)

	return nil
}

func runAnnouncementsDelete(cmd *cobra.Command, args []string) error {
	announcementID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid announcement ID: %s", args[0])
	}

	// Confirm deletion
	confirmed, err := confirmDelete("announcement", announcementID, announcementsForce)
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

	discussionsService := api.NewDiscussionsService(client)

	ctx := context.Background()
	if err := discussionsService.Delete(ctx, announcementsCourseID, announcementID); err != nil {
		return fmt.Errorf("failed to delete announcement: %w", err)
	}

	fmt.Printf("Announcement %d deleted successfully\n", announcementID)
	return nil
}

func displayAnnouncement(announcement *api.DiscussionTopic) {
	fmt.Printf("📢 [%d] %s\n", announcement.ID, announcement.Title)

	if announcement.PostedAt != nil {
		fmt.Printf("   Posted: %s\n", announcement.PostedAt.Format("2006-01-02 15:04"))
	}

	if announcement.DelayedPostAt != nil {
		fmt.Printf("   Scheduled for: %s\n", announcement.DelayedPostAt.Format("2006-01-02 15:04"))
	}

	if announcement.Author != nil {
		fmt.Printf("   By: %s\n", announcement.Author.Name)
	}

	fmt.Println()
}
