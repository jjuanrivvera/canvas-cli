package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var (
	discussionsCourseID        int64
	discussionsOrderBy         string
	discussionsScope           string
	discussionsOnlyAnnounce    bool
	discussionsFilterBy        string
	discussionsSearchTerm      string
	discussionsInclude         []string
	discussionsTitle           string
	discussionsMessage         string
	discussionsType            string
	discussionsPublished       bool
	discussionsDelayedPostAt   string
	discussionsAllowRating     bool
	discussionsLockAt          string
	discussionsRequireInitPost bool
	discussionsPinned          bool
	discussionsLocked          bool
)

// discussionsCmd represents the discussions command group
var discussionsCmd = &cobra.Command{
	Use:   "discussions",
	Short: "Manage Canvas discussion topics",
	Long: `Manage Canvas discussion topics including listing, viewing, creating, and updating discussions.

Discussion topics are threaded conversations associated with Courses in Canvas.
They can be used for class discussions, Q&A, and collaborative learning.

Examples:
  canvas discussions list --course-id 123
  canvas discussions get --course-id 123 456
  canvas discussions create --course-id 123 --title "Week 1 Discussion"
  canvas discussions entries --course-id 123 456`,
}

// discussionsListCmd represents the discussions list command
var discussionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List discussion topics in a course",
	Long: `List all discussion topics in a Canvas course.

Examples:
  canvas discussions list --course-id 123
  canvas discussions list --course-id 123 --order-by recent_activity
  canvas discussions list --course-id 123 --scope pinned
  canvas discussions list --course-id 123 --filter unread`,
	RunE: runDiscussionsList,
}

// discussionsGetCmd represents the discussions get command
var discussionsGetCmd = &cobra.Command{
	Use:   "get <topic-id>",
	Short: "Get a specific discussion topic",
	Long: `Get details of a specific discussion topic.

Examples:
  canvas discussions get --course-id 123 456`,
	Args: cobra.ExactArgs(1),
	RunE: runDiscussionsGet,
}

// discussionsCreateCmd represents the discussions create command
var discussionsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new discussion topic",
	Long: `Create a new discussion topic in a course.

Examples:
  canvas discussions create --course-id 123 --title "Week 1 Discussion"
  canvas discussions create --course-id 123 --title "Q&A" --message "<p>Ask questions here</p>" --type threaded
  canvas discussions create --course-id 123 --title "Pinned" --pinned --published`,
	RunE: runDiscussionsCreate,
}

// discussionsUpdateCmd represents the discussions update command
var discussionsUpdateCmd = &cobra.Command{
	Use:   "update <topic-id>",
	Short: "Update an existing discussion topic",
	Long: `Update an existing discussion topic.

Examples:
  canvas discussions update --course-id 123 456 --title "New Title"
  canvas discussions update --course-id 123 456 --pinned
  canvas discussions update --course-id 123 456 --locked`,
	Args: cobra.ExactArgs(1),
	RunE: runDiscussionsUpdate,
}

// discussionsDeleteCmd represents the discussions delete command
var discussionsDeleteCmd = &cobra.Command{
	Use:   "delete <topic-id>",
	Short: "Delete a discussion topic",
	Long: `Delete a discussion topic from a course.

Examples:
  canvas discussions delete --course-id 123 456`,
	Args: cobra.ExactArgs(1),
	RunE: runDiscussionsDelete,
}

// discussionsEntriesCmd represents the discussions entries command
var discussionsEntriesCmd = &cobra.Command{
	Use:   "entries <topic-id>",
	Short: "List entries in a discussion",
	Long: `List all entries (posts) in a discussion topic.

Examples:
  canvas discussions entries --course-id 123 456`,
	Args: cobra.ExactArgs(1),
	RunE: runDiscussionsEntries,
}

// discussionsPostCmd represents the discussions post command
var discussionsPostCmd = &cobra.Command{
	Use:   "post <topic-id> <message>",
	Short: "Post a new entry to a discussion",
	Long: `Post a new entry to a discussion topic.

Examples:
  canvas discussions post --course-id 123 456 "My response to the discussion"`,
	Args: cobra.ExactArgs(2),
	RunE: runDiscussionsPost,
}

// discussionsReplyCmd represents the discussions reply command
var discussionsReplyCmd = &cobra.Command{
	Use:   "reply <topic-id> <entry-id> <message>",
	Short: "Reply to an entry in a discussion",
	Long: `Reply to a specific entry in a discussion topic.

Examples:
  canvas discussions reply --course-id 123 456 789 "My reply to this entry"`,
	Args: cobra.ExactArgs(3),
	RunE: runDiscussionsReply,
}

// discussionsSubscribeCmd represents the discussions subscribe command
var discussionsSubscribeCmd = &cobra.Command{
	Use:   "subscribe <topic-id>",
	Short: "Subscribe to a discussion topic",
	Long: `Subscribe to receive notifications for a discussion topic.

Examples:
  canvas discussions subscribe --course-id 123 456`,
	Args: cobra.ExactArgs(1),
	RunE: runDiscussionsSubscribe,
}

// discussionsUnsubscribeCmd represents the discussions unsubscribe command
var discussionsUnsubscribeCmd = &cobra.Command{
	Use:   "unsubscribe <topic-id>",
	Short: "Unsubscribe from a discussion topic",
	Long: `Unsubscribe from a discussion topic to stop receiving notifications.

Examples:
  canvas discussions unsubscribe --course-id 123 456`,
	Args: cobra.ExactArgs(1),
	RunE: runDiscussionsUnsubscribe,
}

func init() {
	rootCmd.AddCommand(discussionsCmd)
	discussionsCmd.AddCommand(discussionsListCmd)
	discussionsCmd.AddCommand(discussionsGetCmd)
	discussionsCmd.AddCommand(discussionsCreateCmd)
	discussionsCmd.AddCommand(discussionsUpdateCmd)
	discussionsCmd.AddCommand(discussionsDeleteCmd)
	discussionsCmd.AddCommand(discussionsEntriesCmd)
	discussionsCmd.AddCommand(discussionsPostCmd)
	discussionsCmd.AddCommand(discussionsReplyCmd)
	discussionsCmd.AddCommand(discussionsSubscribeCmd)
	discussionsCmd.AddCommand(discussionsUnsubscribeCmd)

	// List flags
	discussionsListCmd.Flags().Int64Var(&discussionsCourseID, "course-id", 0, "Course ID (required)")
	discussionsListCmd.Flags().StringVar(&discussionsOrderBy, "order-by", "", "Order by: position, recent_activity, title")
	discussionsListCmd.Flags().StringVar(&discussionsScope, "scope", "", "Scope: locked, unlocked, pinned, unpinned")
	discussionsListCmd.Flags().BoolVar(&discussionsOnlyAnnounce, "announcements", false, "Only show announcements")
	discussionsListCmd.Flags().StringVar(&discussionsFilterBy, "filter", "", "Filter by: all, unread")
	discussionsListCmd.Flags().StringVar(&discussionsSearchTerm, "search", "", "Search term")
	discussionsListCmd.Flags().StringSliceVar(&discussionsInclude, "include", []string{}, "Include: all_dates, sections, sections_user_count, overrides")
	discussionsListCmd.MarkFlagRequired("course-id")

	// Get flags
	discussionsGetCmd.Flags().Int64Var(&discussionsCourseID, "course-id", 0, "Course ID (required)")
	discussionsGetCmd.Flags().StringSliceVar(&discussionsInclude, "include", []string{}, "Include additional data")
	discussionsGetCmd.MarkFlagRequired("course-id")

	// Create flags
	discussionsCreateCmd.Flags().Int64Var(&discussionsCourseID, "course-id", 0, "Course ID (required)")
	discussionsCreateCmd.Flags().StringVar(&discussionsTitle, "title", "", "Discussion title (required)")
	discussionsCreateCmd.Flags().StringVar(&discussionsMessage, "message", "", "Discussion message (HTML)")
	discussionsCreateCmd.Flags().StringVar(&discussionsType, "type", "", "Discussion type: side_comment, threaded, not_threaded")
	discussionsCreateCmd.Flags().BoolVar(&discussionsPublished, "published", false, "Publish the discussion")
	discussionsCreateCmd.Flags().StringVar(&discussionsDelayedPostAt, "delayed-post-at", "", "Delay posting until (ISO 8601)")
	discussionsCreateCmd.Flags().BoolVar(&discussionsAllowRating, "allow-rating", false, "Allow rating of entries")
	discussionsCreateCmd.Flags().StringVar(&discussionsLockAt, "lock-at", "", "Lock at date (ISO 8601)")
	discussionsCreateCmd.Flags().BoolVar(&discussionsRequireInitPost, "require-initial-post", false, "Require initial post before viewing")
	discussionsCreateCmd.Flags().BoolVar(&discussionsPinned, "pinned", false, "Pin the discussion")
	discussionsCreateCmd.MarkFlagRequired("course-id")
	discussionsCreateCmd.MarkFlagRequired("title")

	// Update flags
	discussionsUpdateCmd.Flags().Int64Var(&discussionsCourseID, "course-id", 0, "Course ID (required)")
	discussionsUpdateCmd.Flags().StringVar(&discussionsTitle, "title", "", "New discussion title")
	discussionsUpdateCmd.Flags().StringVar(&discussionsMessage, "message", "", "New discussion message")
	discussionsUpdateCmd.Flags().StringVar(&discussionsType, "type", "", "Discussion type")
	discussionsUpdateCmd.Flags().BoolVar(&discussionsPublished, "published", false, "Publish the discussion")
	discussionsUpdateCmd.Flags().StringVar(&discussionsDelayedPostAt, "delayed-post-at", "", "Delay posting until")
	discussionsUpdateCmd.Flags().BoolVar(&discussionsAllowRating, "allow-rating", false, "Allow rating")
	discussionsUpdateCmd.Flags().StringVar(&discussionsLockAt, "lock-at", "", "Lock at date")
	discussionsUpdateCmd.Flags().BoolVar(&discussionsRequireInitPost, "require-initial-post", false, "Require initial post")
	discussionsUpdateCmd.Flags().BoolVar(&discussionsPinned, "pinned", false, "Pin the discussion")
	discussionsUpdateCmd.Flags().BoolVar(&discussionsLocked, "locked", false, "Lock the discussion")
	discussionsUpdateCmd.MarkFlagRequired("course-id")

	// Delete flags
	discussionsDeleteCmd.Flags().Int64Var(&discussionsCourseID, "course-id", 0, "Course ID (required)")
	discussionsDeleteCmd.MarkFlagRequired("course-id")

	// Entries flags
	discussionsEntriesCmd.Flags().Int64Var(&discussionsCourseID, "course-id", 0, "Course ID (required)")
	discussionsEntriesCmd.MarkFlagRequired("course-id")

	// Post flags
	discussionsPostCmd.Flags().Int64Var(&discussionsCourseID, "course-id", 0, "Course ID (required)")
	discussionsPostCmd.MarkFlagRequired("course-id")

	// Reply flags
	discussionsReplyCmd.Flags().Int64Var(&discussionsCourseID, "course-id", 0, "Course ID (required)")
	discussionsReplyCmd.MarkFlagRequired("course-id")

	// Subscribe flags
	discussionsSubscribeCmd.Flags().Int64Var(&discussionsCourseID, "course-id", 0, "Course ID (required)")
	discussionsSubscribeCmd.MarkFlagRequired("course-id")

	// Unsubscribe flags
	discussionsUnsubscribeCmd.Flags().Int64Var(&discussionsCourseID, "course-id", 0, "Course ID (required)")
	discussionsUnsubscribeCmd.MarkFlagRequired("course-id")
}

func runDiscussionsList(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	discussionsService := api.NewDiscussionsService(client)

	opts := &api.ListDiscussionsOptions{
		Include:           discussionsInclude,
		OrderBy:           discussionsOrderBy,
		Scope:             discussionsScope,
		OnlyAnnouncements: discussionsOnlyAnnounce,
		FilterBy:          discussionsFilterBy,
		SearchTerm:        discussionsSearchTerm,
	}

	ctx := context.Background()
	topics, err := discussionsService.List(ctx, discussionsCourseID, opts)
	if err != nil {
		return fmt.Errorf("failed to list discussions: %w", err)
	}

	if len(topics) == 0 {
		fmt.Println("No discussion topics found")
		return nil
	}

	fmt.Printf("Found %d discussion topics:\n\n", len(topics))

	for _, topic := range topics {
		displayDiscussionTopic(&topic)
	}

	return nil
}

func runDiscussionsGet(cmd *cobra.Command, args []string) error {
	topicID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid topic ID: %s", args[0])
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	discussionsService := api.NewDiscussionsService(client)

	ctx := context.Background()
	topic, err := discussionsService.Get(ctx, discussionsCourseID, topicID, discussionsInclude)
	if err != nil {
		return fmt.Errorf("failed to get discussion: %w", err)
	}

	displayDiscussionTopicFull(topic)

	return nil
}

func runDiscussionsCreate(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	discussionsService := api.NewDiscussionsService(client)

	params := &api.CreateDiscussionParams{
		Title:              discussionsTitle,
		Message:            discussionsMessage,
		DiscussionType:     discussionsType,
		Published:          discussionsPublished,
		DelayedPostAt:      discussionsDelayedPostAt,
		AllowRating:        discussionsAllowRating,
		LockAt:             discussionsLockAt,
		RequireInitialPost: discussionsRequireInitPost,
		Pinned:             discussionsPinned,
	}

	ctx := context.Background()
	topic, err := discussionsService.Create(ctx, discussionsCourseID, params)
	if err != nil {
		return fmt.Errorf("failed to create discussion: %w", err)
	}

	fmt.Println("Discussion created successfully!")
	displayDiscussionTopic(topic)

	return nil
}

func runDiscussionsUpdate(cmd *cobra.Command, args []string) error {
	topicID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid topic ID: %s", args[0])
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	discussionsService := api.NewDiscussionsService(client)

	params := &api.UpdateDiscussionParams{}

	if cmd.Flags().Changed("title") {
		params.Title = &discussionsTitle
	}
	if cmd.Flags().Changed("message") {
		params.Message = &discussionsMessage
	}
	if cmd.Flags().Changed("type") {
		params.DiscussionType = &discussionsType
	}
	if cmd.Flags().Changed("published") {
		params.Published = &discussionsPublished
	}
	if cmd.Flags().Changed("delayed-post-at") {
		params.DelayedPostAt = &discussionsDelayedPostAt
	}
	if cmd.Flags().Changed("allow-rating") {
		params.AllowRating = &discussionsAllowRating
	}
	if cmd.Flags().Changed("lock-at") {
		params.LockAt = &discussionsLockAt
	}
	if cmd.Flags().Changed("require-initial-post") {
		params.RequireInitialPost = &discussionsRequireInitPost
	}
	if cmd.Flags().Changed("pinned") {
		params.Pinned = &discussionsPinned
	}
	if cmd.Flags().Changed("locked") {
		params.Locked = &discussionsLocked
	}

	ctx := context.Background()
	topic, err := discussionsService.Update(ctx, discussionsCourseID, topicID, params)
	if err != nil {
		return fmt.Errorf("failed to update discussion: %w", err)
	}

	fmt.Println("Discussion updated successfully!")
	displayDiscussionTopic(topic)

	return nil
}

func runDiscussionsDelete(cmd *cobra.Command, args []string) error {
	topicID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid topic ID: %s", args[0])
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	discussionsService := api.NewDiscussionsService(client)

	ctx := context.Background()
	if err := discussionsService.Delete(ctx, discussionsCourseID, topicID); err != nil {
		return fmt.Errorf("failed to delete discussion: %w", err)
	}

	fmt.Printf("Discussion %d deleted successfully\n", topicID)
	return nil
}

func runDiscussionsEntries(cmd *cobra.Command, args []string) error {
	topicID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid topic ID: %s", args[0])
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	discussionsService := api.NewDiscussionsService(client)

	ctx := context.Background()
	entries, err := discussionsService.ListEntries(ctx, discussionsCourseID, topicID)
	if err != nil {
		return fmt.Errorf("failed to list entries: %w", err)
	}

	if len(entries) == 0 {
		fmt.Println("No entries found")
		return nil
	}

	fmt.Printf("Found %d entries:\n\n", len(entries))

	for _, entry := range entries {
		displayDiscussionEntry(&entry, 0)
	}

	return nil
}

func runDiscussionsPost(cmd *cobra.Command, args []string) error {
	topicID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid topic ID: %s", args[0])
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	discussionsService := api.NewDiscussionsService(client)

	ctx := context.Background()
	entry, err := discussionsService.PostEntry(ctx, discussionsCourseID, topicID, args[1])
	if err != nil {
		return fmt.Errorf("failed to post entry: %w", err)
	}

	fmt.Println("Entry posted successfully!")
	displayDiscussionEntry(entry, 0)

	return nil
}

func runDiscussionsReply(cmd *cobra.Command, args []string) error {
	topicID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid topic ID: %s", args[0])
	}

	entryID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid entry ID: %s", args[1])
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	discussionsService := api.NewDiscussionsService(client)

	ctx := context.Background()
	entry, err := discussionsService.PostReply(ctx, discussionsCourseID, topicID, entryID, args[2])
	if err != nil {
		return fmt.Errorf("failed to post reply: %w", err)
	}

	fmt.Println("Reply posted successfully!")
	displayDiscussionEntry(entry, 0)

	return nil
}

func runDiscussionsSubscribe(cmd *cobra.Command, args []string) error {
	topicID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid topic ID: %s", args[0])
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	discussionsService := api.NewDiscussionsService(client)

	ctx := context.Background()
	if err := discussionsService.Subscribe(ctx, discussionsCourseID, topicID); err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	fmt.Printf("Subscribed to discussion %d\n", topicID)
	return nil
}

func runDiscussionsUnsubscribe(cmd *cobra.Command, args []string) error {
	topicID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid topic ID: %s", args[0])
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	discussionsService := api.NewDiscussionsService(client)

	ctx := context.Background()
	if err := discussionsService.Unsubscribe(ctx, discussionsCourseID, topicID); err != nil {
		return fmt.Errorf("failed to unsubscribe: %w", err)
	}

	fmt.Printf("Unsubscribed from discussion %d\n", topicID)
	return nil
}

func displayDiscussionTopic(topic *api.DiscussionTopic) {
	stateIcon := "💬"
	if topic.Pinned {
		stateIcon = "📌"
	} else if topic.Locked {
		stateIcon = "🔒"
	}

	fmt.Printf("%s [%d] %s\n", stateIcon, topic.ID, topic.Title)

	if topic.Published {
		fmt.Printf("   Published: Yes\n")
	} else {
		fmt.Printf("   Published: No (Draft)\n")
	}

	fmt.Printf("   Replies: %d", topic.DiscussionSubentryCount)
	if topic.UnreadCount > 0 {
		fmt.Printf(" (%d unread)", topic.UnreadCount)
	}
	fmt.Println()

	if topic.DiscussionType != "" {
		fmt.Printf("   Type: %s\n", topic.DiscussionType)
	}

	if topic.PostedAt != nil {
		fmt.Printf("   Posted: %s\n", topic.PostedAt.Format("2006-01-02 15:04"))
	}

	fmt.Println()
}

func displayDiscussionTopicFull(topic *api.DiscussionTopic) {
	displayDiscussionTopic(topic)

	if topic.Author != nil {
		fmt.Printf("   Author: %s\n", topic.Author.Name)
	}

	if topic.RequireInitialPost {
		fmt.Printf("   Requires Initial Post: Yes\n")
	}

	if topic.AllowRating {
		fmt.Printf("   Allow Rating: Yes\n")
	}

	if topic.Subscribed {
		fmt.Printf("   Subscribed: Yes\n")
	}

	if topic.Message != "" {
		fmt.Printf("\nMessage:\n")
		message := topic.Message
		if len(message) > 500 {
			message = message[:500] + "..."
		}
		message = stripHTMLTags(message)
		fmt.Println(message)
	}

	fmt.Println()
}

func displayDiscussionEntry(entry *api.DiscussionEntry, indent int) {
	prefix := ""
	for i := 0; i < indent; i++ {
		prefix += "  "
	}

	stateIcon := "💬"
	if entry.ReadState == "unread" {
		stateIcon = "🆕"
	}

	fmt.Printf("%s%s Entry %d (User %d)\n", prefix, stateIcon, entry.ID, entry.UserID)

	if entry.User != nil {
		fmt.Printf("%s   By: %s\n", prefix, entry.User.Name)
	}

	fmt.Printf("%s   Posted: %s\n", prefix, entry.CreatedAt.Format("2006-01-02 15:04"))

	if entry.Message != "" {
		message := entry.Message
		if len(message) > 200 {
			message = message[:200] + "..."
		}
		message = stripHTMLTags(message)
		fmt.Printf("%s   %s\n", prefix, message)
	}

	fmt.Println()

	// Display nested replies
	for _, reply := range entry.Replies {
		displayDiscussionEntry(&reply, indent+1)
	}
}
