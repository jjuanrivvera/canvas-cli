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

func init() {
	rootCmd.AddCommand(discussionsCmd)
	discussionsCmd.AddCommand(newDiscussionsListCmd())
	discussionsCmd.AddCommand(newDiscussionsGetCmd())
	discussionsCmd.AddCommand(newDiscussionsCreateCmd())
	discussionsCmd.AddCommand(newDiscussionsUpdateCmd())
	discussionsCmd.AddCommand(newDiscussionsDeleteCmd())
	discussionsCmd.AddCommand(newDiscussionsEntriesCmd())
	discussionsCmd.AddCommand(newDiscussionsPostCmd())
	discussionsCmd.AddCommand(newDiscussionsReplyCmd())
	discussionsCmd.AddCommand(newDiscussionsSubscribeCmd())
	discussionsCmd.AddCommand(newDiscussionsUnsubscribeCmd())
}

func newDiscussionsListCmd() *cobra.Command {
	opts := &options.DiscussionsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List discussion topics in a course",
		Long: `List all discussion topics in a Canvas course.

Examples:
  canvas discussions list --course-id 123
  canvas discussions list --course-id 123 --order-by recent_activity
  canvas discussions list --course-id 123 --scope pinned
  canvas discussions list --course-id 123 --filter unread`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.OrderBy, "order-by", "", "Order by: position, recent_activity, title")
	cmd.Flags().StringVar(&opts.Scope, "scope", "", "Scope: locked, unlocked, pinned, unpinned")
	cmd.Flags().BoolVar(&opts.OnlyAnnouncements, "announcements", false, "Only show announcements")
	cmd.Flags().StringVar(&opts.FilterBy, "filter", "", "Filter by: all, unread")
	cmd.Flags().StringVar(&opts.SearchTerm, "search", "", "Search term")
	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Include: all_dates, sections, sections_user_count, overrides")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newDiscussionsGetCmd() *cobra.Command {
	opts := &options.DiscussionsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <topic-id>",
		Short: "Get a specific discussion topic",
		Long: `Get details of a specific discussion topic.

Examples:
  canvas discussions get --course-id 123 456`,
		Args: ExactArgsWithUsage(1, "topic-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid topic ID: %s", args[0])
			}
			opts.TopicID = topicID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Include additional data")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newDiscussionsCreateCmd() *cobra.Command {
	opts := &options.DiscussionsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new discussion topic",
		Long: `Create a new discussion topic in a course.

Examples:
  canvas discussions create --course-id 123 --title "Week 1 Discussion"
  canvas discussions create --course-id 123 --title "Q&A" --message "<p>Ask questions here</p>" --type threaded
  canvas discussions create --course-id 123 --title "Pinned" --pinned --published`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.Title, "title", "", "Discussion title (required)")
	cmd.Flags().StringVar(&opts.Message, "message", "", "Discussion message (HTML)")
	cmd.Flags().StringVar(&opts.DiscussionType, "type", "", "Discussion type: side_comment, threaded, not_threaded")
	cmd.Flags().BoolVar(&opts.Published, "published", false, "Publish the discussion")
	cmd.Flags().StringVar(&opts.DelayedPostAt, "delayed-post-at", "", "Delay posting until (ISO 8601)")
	cmd.Flags().BoolVar(&opts.AllowRating, "allow-rating", false, "Allow rating of entries")
	cmd.Flags().StringVar(&opts.LockAt, "lock-at", "", "Lock at date (ISO 8601)")
	cmd.Flags().BoolVar(&opts.RequireInitialPost, "require-initial-post", false, "Require initial post before viewing")
	cmd.Flags().BoolVar(&opts.Pinned, "pinned", false, "Pin the discussion")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("title")

	return cmd
}

func newDiscussionsUpdateCmd() *cobra.Command {
	opts := &options.DiscussionsUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <topic-id>",
		Short: "Update an existing discussion topic",
		Long: `Update an existing discussion topic.

Examples:
  canvas discussions update --course-id 123 456 --title "New Title"
  canvas discussions update --course-id 123 456 --pinned
  canvas discussions update --course-id 123 456 --locked`,
		Args: ExactArgsWithUsage(1, "topic-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid topic ID: %s", args[0])
			}
			opts.TopicID = topicID
			// Track which fields were changed
			opts.TitleSet = cmd.Flags().Changed("title")
			opts.MessageSet = cmd.Flags().Changed("message")
			opts.DiscussionTypeSet = cmd.Flags().Changed("type")
			opts.PublishedSet = cmd.Flags().Changed("published")
			opts.DelayedPostAtSet = cmd.Flags().Changed("delayed-post-at")
			opts.AllowRatingSet = cmd.Flags().Changed("allow-rating")
			opts.LockAtSet = cmd.Flags().Changed("lock-at")
			opts.RequireInitialPostSet = cmd.Flags().Changed("require-initial-post")
			opts.PinnedSet = cmd.Flags().Changed("pinned")
			opts.LockedSet = cmd.Flags().Changed("locked")
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.Title, "title", "", "New discussion title")
	cmd.Flags().StringVar(&opts.Message, "message", "", "New discussion message")
	cmd.Flags().StringVar(&opts.DiscussionType, "type", "", "Discussion type")
	cmd.Flags().BoolVar(&opts.Published, "published", false, "Publish the discussion")
	cmd.Flags().StringVar(&opts.DelayedPostAt, "delayed-post-at", "", "Delay posting until")
	cmd.Flags().BoolVar(&opts.AllowRating, "allow-rating", false, "Allow rating")
	cmd.Flags().StringVar(&opts.LockAt, "lock-at", "", "Lock at date")
	cmd.Flags().BoolVar(&opts.RequireInitialPost, "require-initial-post", false, "Require initial post")
	cmd.Flags().BoolVar(&opts.Pinned, "pinned", false, "Pin the discussion")
	cmd.Flags().BoolVar(&opts.Locked, "locked", false, "Lock the discussion")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newDiscussionsDeleteCmd() *cobra.Command {
	opts := &options.DiscussionsDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <topic-id>",
		Short: "Delete a discussion topic",
		Long: `Delete a discussion topic from a course.

Examples:
  canvas discussions delete --course-id 123 456`,
		Args: ExactArgsWithUsage(1, "topic-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid topic ID: %s", args[0])
			}
			opts.TopicID = topicID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Skip confirmation prompt")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newDiscussionsEntriesCmd() *cobra.Command {
	opts := &options.DiscussionsEntriesOptions{}

	cmd := &cobra.Command{
		Use:   "entries <topic-id>",
		Short: "List entries in a discussion",
		Long: `List all entries (posts) in a discussion topic.

Examples:
  canvas discussions entries --course-id 123 456`,
		Args: ExactArgsWithUsage(1, "topic-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid topic ID: %s", args[0])
			}
			opts.TopicID = topicID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsEntries(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newDiscussionsPostCmd() *cobra.Command {
	opts := &options.DiscussionsPostOptions{}

	cmd := &cobra.Command{
		Use:   "post <topic-id> [message]",
		Short: "Post a new entry to a discussion",
		Long: `Post a new entry to a discussion topic.

The message can be provided as a positional argument or using the --message flag.

Examples:
  canvas discussions post --course-id 123 456 "My response to the discussion"
  canvas discussions post --course-id 123 456 --message "My response to the discussion"`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid topic ID: %s", args[0])
			}
			opts.TopicID = topicID
			// Get message from positional arg or --message flag
			if len(args) > 1 {
				opts.Message = args[1]
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsPost(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVarP(&opts.Message, "message", "m", "", "Message content (alternative to positional argument)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newDiscussionsReplyCmd() *cobra.Command {
	opts := &options.DiscussionsReplyOptions{}

	cmd := &cobra.Command{
		Use:   "reply <topic-id> <entry-id> [message]",
		Short: "Reply to an entry in a discussion",
		Long: `Reply to a specific entry in a discussion topic.

The message can be provided as a positional argument or using the --message flag.

Examples:
  canvas discussions reply --course-id 123 456 789 "My reply to this entry"
  canvas discussions reply --course-id 123 456 789 --message "My reply to this entry"`,
		Args: cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid topic ID: %s", args[0])
			}
			opts.TopicID = topicID
			entryID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid entry ID: %s", args[1])
			}
			opts.EntryID = entryID
			// Get message from positional arg or --message flag
			if len(args) > 2 {
				opts.Message = args[2]
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsReply(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVarP(&opts.Message, "message", "m", "", "Message content (alternative to positional argument)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newDiscussionsSubscribeCmd() *cobra.Command {
	opts := &options.DiscussionsSubscribeOptions{}

	cmd := &cobra.Command{
		Use:   "subscribe <topic-id>",
		Short: "Subscribe to a discussion topic",
		Long: `Subscribe to receive notifications for a discussion topic.

Examples:
  canvas discussions subscribe --course-id 123 456`,
		Args: ExactArgsWithUsage(1, "topic-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid topic ID: %s", args[0])
			}
			opts.TopicID = topicID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsSubscribe(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newDiscussionsUnsubscribeCmd() *cobra.Command {
	opts := &options.DiscussionsUnsubscribeOptions{}

	cmd := &cobra.Command{
		Use:   "unsubscribe <topic-id>",
		Short: "Unsubscribe from a discussion topic",
		Long: `Unsubscribe from a discussion topic to stop receiving notifications.

Examples:
  canvas discussions unsubscribe --course-id 123 456`,
		Args: ExactArgsWithUsage(1, "topic-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid topic ID: %s", args[0])
			}
			opts.TopicID = topicID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsUnsubscribe(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func runDiscussionsList(ctx context.Context, client *api.Client, opts *options.DiscussionsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "discussions.list", map[string]interface{}{
		"course_id":   opts.CourseID,
		"order_by":    opts.OrderBy,
		"scope":       opts.Scope,
		"filter_by":   opts.FilterBy,
		"search_term": opts.SearchTerm,
	})

	discussionsService := api.NewDiscussionsService(client)

	listOpts := &api.ListDiscussionsOptions{
		Include:           opts.Include,
		OrderBy:           opts.OrderBy,
		Scope:             opts.Scope,
		OnlyAnnouncements: opts.OnlyAnnouncements,
		FilterBy:          opts.FilterBy,
		SearchTerm:        opts.SearchTerm,
	}

	topics, err := discussionsService.List(ctx, opts.CourseID, listOpts)
	if err != nil {
		logger.LogCommandError(ctx, "discussions.list", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to list discussions: %w", err)
	}

	if len(topics) == 0 {
		fmt.Println("No discussion topics found")
		logger.LogCommandComplete(ctx, "discussions.list", 0)
		return nil
	}

	printVerbose("Found %d discussion topics:\n\n", len(topics))
	logger.LogCommandComplete(ctx, "discussions.list", len(topics))
	return formatOutput(topics, nil)
}

func runDiscussionsGet(ctx context.Context, client *api.Client, opts *options.DiscussionsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "discussions.get", map[string]interface{}{
		"course_id": opts.CourseID,
		"topic_id":  opts.TopicID,
	})

	discussionsService := api.NewDiscussionsService(client)

	topic, err := discussionsService.Get(ctx, opts.CourseID, opts.TopicID, opts.Include)
	if err != nil {
		logger.LogCommandError(ctx, "discussions.get", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"topic_id":  opts.TopicID,
		})
		return fmt.Errorf("failed to get discussion: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.get", 1)
	return formatOutput(topic, nil)
}

func runDiscussionsCreate(ctx context.Context, client *api.Client, opts *options.DiscussionsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "discussions.create", map[string]interface{}{
		"course_id": opts.CourseID,
		"title":     opts.Title,
	})

	discussionsService := api.NewDiscussionsService(client)

	params := &api.CreateDiscussionParams{
		Title:              opts.Title,
		Message:            opts.Message,
		DiscussionType:     opts.DiscussionType,
		Published:          opts.Published,
		DelayedPostAt:      opts.DelayedPostAt,
		AllowRating:        opts.AllowRating,
		LockAt:             opts.LockAt,
		RequireInitialPost: opts.RequireInitialPost,
		Pinned:             opts.Pinned,
	}

	topic, err := discussionsService.Create(ctx, opts.CourseID, params)
	if err != nil {
		logger.LogCommandError(ctx, "discussions.create", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"title":     opts.Title,
		})
		return fmt.Errorf("failed to create discussion: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.create", 1)
	return formatSuccessOutput(topic, "Discussion created successfully!")
}

func runDiscussionsUpdate(ctx context.Context, client *api.Client, opts *options.DiscussionsUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "discussions.update", map[string]interface{}{
		"course_id": opts.CourseID,
		"topic_id":  opts.TopicID,
	})

	discussionsService := api.NewDiscussionsService(client)

	params := &api.UpdateDiscussionParams{}

	if opts.TitleSet {
		params.Title = &opts.Title
	}
	if opts.MessageSet {
		params.Message = &opts.Message
	}
	if opts.DiscussionTypeSet {
		params.DiscussionType = &opts.DiscussionType
	}
	if opts.PublishedSet {
		params.Published = &opts.Published
	}
	if opts.DelayedPostAtSet {
		params.DelayedPostAt = &opts.DelayedPostAt
	}
	if opts.AllowRatingSet {
		params.AllowRating = &opts.AllowRating
	}
	if opts.LockAtSet {
		params.LockAt = &opts.LockAt
	}
	if opts.RequireInitialPostSet {
		params.RequireInitialPost = &opts.RequireInitialPost
	}
	if opts.PinnedSet {
		params.Pinned = &opts.Pinned
	}
	if opts.LockedSet {
		params.Locked = &opts.Locked
	}

	topic, err := discussionsService.Update(ctx, opts.CourseID, opts.TopicID, params)
	if err != nil {
		logger.LogCommandError(ctx, "discussions.update", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"topic_id":  opts.TopicID,
		})
		return fmt.Errorf("failed to update discussion: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.update", 1)
	return formatSuccessOutput(topic, "Discussion updated successfully!")
}

func runDiscussionsDelete(ctx context.Context, client *api.Client, opts *options.DiscussionsDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "discussions.delete", map[string]interface{}{
		"course_id": opts.CourseID,
		"topic_id":  opts.TopicID,
		"force":     opts.Force,
	})

	// Confirm deletion
	confirmed, err := confirmDelete("discussion", opts.TopicID, opts.Force)
	if err != nil {
		return err
	}
	if !confirmed {
		fmt.Println("Delete cancelled")
		return nil
	}

	discussionsService := api.NewDiscussionsService(client)

	if err := discussionsService.Delete(ctx, opts.CourseID, opts.TopicID); err != nil {
		logger.LogCommandError(ctx, "discussions.delete", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"topic_id":  opts.TopicID,
		})
		return fmt.Errorf("failed to delete discussion: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.delete", 1)
	fmt.Printf("Discussion %d deleted successfully\n", opts.TopicID)
	return nil
}

func runDiscussionsEntries(ctx context.Context, client *api.Client, opts *options.DiscussionsEntriesOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "discussions.entries", map[string]interface{}{
		"course_id": opts.CourseID,
		"topic_id":  opts.TopicID,
	})

	discussionsService := api.NewDiscussionsService(client)

	entries, err := discussionsService.ListEntries(ctx, opts.CourseID, opts.TopicID)
	if err != nil {
		logger.LogCommandError(ctx, "discussions.entries", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"topic_id":  opts.TopicID,
		})
		return fmt.Errorf("failed to list entries: %w", err)
	}

	if len(entries) == 0 {
		fmt.Println("No entries found")
		logger.LogCommandComplete(ctx, "discussions.entries", 0)
		return nil
	}

	printVerbose("Found %d entries:\n\n", len(entries))
	logger.LogCommandComplete(ctx, "discussions.entries", len(entries))
	return formatOutput(entries, nil)
}

func runDiscussionsPost(ctx context.Context, client *api.Client, opts *options.DiscussionsPostOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "discussions.post", map[string]interface{}{
		"course_id": opts.CourseID,
		"topic_id":  opts.TopicID,
	})

	discussionsService := api.NewDiscussionsService(client)

	entry, err := discussionsService.PostEntry(ctx, opts.CourseID, opts.TopicID, opts.Message)
	if err != nil {
		logger.LogCommandError(ctx, "discussions.post", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"topic_id":  opts.TopicID,
		})
		return fmt.Errorf("failed to post entry: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.post", 1)
	return formatSuccessOutput(entry, "Entry posted successfully!")
}

func runDiscussionsReply(ctx context.Context, client *api.Client, opts *options.DiscussionsReplyOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "discussions.reply", map[string]interface{}{
		"course_id": opts.CourseID,
		"topic_id":  opts.TopicID,
		"entry_id":  opts.EntryID,
	})

	discussionsService := api.NewDiscussionsService(client)

	entry, err := discussionsService.PostReply(ctx, opts.CourseID, opts.TopicID, opts.EntryID, opts.Message)
	if err != nil {
		logger.LogCommandError(ctx, "discussions.reply", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"topic_id":  opts.TopicID,
			"entry_id":  opts.EntryID,
		})
		return fmt.Errorf("failed to post reply: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.reply", 1)
	return formatSuccessOutput(entry, "Reply posted successfully!")
}

func runDiscussionsSubscribe(ctx context.Context, client *api.Client, opts *options.DiscussionsSubscribeOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "discussions.subscribe", map[string]interface{}{
		"course_id": opts.CourseID,
		"topic_id":  opts.TopicID,
	})

	discussionsService := api.NewDiscussionsService(client)

	if err := discussionsService.Subscribe(ctx, opts.CourseID, opts.TopicID); err != nil {
		logger.LogCommandError(ctx, "discussions.subscribe", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"topic_id":  opts.TopicID,
		})
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.subscribe", 1)
	fmt.Printf("Subscribed to discussion %d\n", opts.TopicID)
	return nil
}

func runDiscussionsUnsubscribe(ctx context.Context, client *api.Client, opts *options.DiscussionsUnsubscribeOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "discussions.unsubscribe", map[string]interface{}{
		"course_id": opts.CourseID,
		"topic_id":  opts.TopicID,
	})

	discussionsService := api.NewDiscussionsService(client)

	if err := discussionsService.Unsubscribe(ctx, opts.CourseID, opts.TopicID); err != nil {
		logger.LogCommandError(ctx, "discussions.unsubscribe", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"topic_id":  opts.TopicID,
		})
		return fmt.Errorf("failed to unsubscribe: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.unsubscribe", 1)
	fmt.Printf("Unsubscribed from discussion %d\n", opts.TopicID)
	return nil
}
