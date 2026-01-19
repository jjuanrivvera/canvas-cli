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

// conversationsCmd represents the conversations command group
var conversationsCmd = &cobra.Command{
	Use:   "conversations",
	Short: "Manage Canvas conversations (inbox)",
	Long: `Manage Canvas conversations and messages.

Conversations are Canvas's internal messaging system for communication
between users within courses and the institution.

Examples:
  canvas conversations list
  canvas conversations list --scope unread
  canvas conversations get 123
  canvas conversations create --recipients 456,789 --subject "Hello" --body "Message content"`,
}

func init() {
	rootCmd.AddCommand(conversationsCmd)
	conversationsCmd.AddCommand(newConversationsListCmd())
	conversationsCmd.AddCommand(newConversationsGetCmd())
	conversationsCmd.AddCommand(newConversationsCreateCmd())
	conversationsCmd.AddCommand(newConversationsReplyCmd())
	conversationsCmd.AddCommand(newConversationsAddRecipientsCmd())
	conversationsCmd.AddCommand(newConversationsArchiveCmd())
	conversationsCmd.AddCommand(newConversationsUnarchiveCmd())
	conversationsCmd.AddCommand(newConversationsStarCmd())
	conversationsCmd.AddCommand(newConversationsUnstarCmd())
	conversationsCmd.AddCommand(newConversationsMarkReadCmd())
	conversationsCmd.AddCommand(newConversationsMarkAllReadCmd())
	conversationsCmd.AddCommand(newConversationsDeleteCmd())
	conversationsCmd.AddCommand(newConversationsUnreadCountCmd())
}

func newConversationsListCmd() *cobra.Command {
	opts := &options.ConversationsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List conversations",
		Long: `List conversations for the current user.

Examples:
  canvas conversations list
  canvas conversations list --scope unread
  canvas conversations list --scope starred
  canvas conversations list --filter course_123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runConversationsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Scope, "scope", "", "Scope: inbox, unread, archived, starred, sent")
	cmd.Flags().StringSliceVar(&opts.Filter, "filter", nil, "Filter by course/group (e.g., course_123, group_456)")
	cmd.Flags().StringVar(&opts.FilterMode, "filter-mode", "", "Filter mode: and, or")
	cmd.Flags().BoolVar(&opts.InterleaveSubmissions, "interleave", false, "Interleave submissions")

	return cmd
}

func newConversationsGetCmd() *cobra.Command {
	opts := &options.ConversationsGetOptions{
		AutoMarkRead: true, // default value
	}

	cmd := &cobra.Command{
		Use:   "get <conversation-id>",
		Short: "Get conversation details",
		Long: `Get details of a specific conversation including messages.

Examples:
  canvas conversations get 123
  canvas conversations get 123 --auto-mark-read=false`,
		Args: ExactArgsWithUsage(1, "conversation-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			conversationID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid conversation ID: %s", args[0])
			}
			opts.ConversationID = conversationID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runConversationsGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().BoolVar(&opts.AutoMarkRead, "auto-mark-read", true, "Auto mark as read")

	return cmd
}

func newConversationsCreateCmd() *cobra.Command {
	opts := &options.ConversationsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new conversation",
		Long: `Create a new conversation with one or more recipients.

Examples:
  canvas conversations create --recipients 123,456 --subject "Hello" --body "Message text"
  canvas conversations create --recipients 123 --body "Quick message"
  canvas conversations create --recipients course_123 --subject "Announcement" --body "Message" --group`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runConversationsCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Recipients, "recipients", "", "Recipient IDs (comma-separated, required)")
	cmd.Flags().StringVar(&opts.Subject, "subject", "", "Subject line")
	cmd.Flags().StringVar(&opts.Body, "body", "", "Message body (required)")
	cmd.Flags().BoolVar(&opts.ForceNew, "force-new", false, "Force new conversation")
	cmd.Flags().BoolVar(&opts.GroupConversation, "group", false, "Send as group message")
	cmd.Flags().Int64SliceVar(&opts.AttachmentIDs, "attachment-ids", nil, "Attachment IDs")
	cmd.Flags().StringVar(&opts.MediaCommentID, "media-comment-id", "", "Media comment ID")
	cmd.Flags().StringVar(&opts.ContextCode, "context-code", "", "Context code (e.g., course_123)")
	cmd.MarkFlagRequired("recipients")
	cmd.MarkFlagRequired("body")

	return cmd
}

func newConversationsReplyCmd() *cobra.Command {
	opts := &options.ConversationsReplyOptions{}

	cmd := &cobra.Command{
		Use:   "reply <conversation-id>",
		Short: "Reply to a conversation",
		Long: `Add a reply message to an existing conversation.

Examples:
  canvas conversations reply 123 --body "Thank you for your message"
  canvas conversations reply 123 --body "Here's the attachment" --attachment-ids 456`,
		Args: ExactArgsWithUsage(1, "conversation-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			conversationID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid conversation ID: %s", args[0])
			}
			opts.ConversationID = conversationID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runConversationsReply(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Body, "body", "", "Reply body (required)")
	cmd.Flags().Int64SliceVar(&opts.AttachmentIDs, "attachment-ids", nil, "Attachment IDs")
	cmd.Flags().StringVar(&opts.MediaCommentID, "media-comment-id", "", "Media comment ID")
	cmd.Flags().Int64SliceVar(&opts.IncludedMessages, "included-messages", nil, "Include previous message IDs")
	cmd.MarkFlagRequired("body")

	return cmd
}

func newConversationsAddRecipientsCmd() *cobra.Command {
	opts := &options.ConversationsAddRecipientsOptions{}

	cmd := &cobra.Command{
		Use:   "add-recipients <conversation-id>",
		Short: "Add recipients to a conversation",
		Long: `Add additional recipients to an existing conversation.

Examples:
  canvas conversations add-recipients 123 --recipients 456,789`,
		Args: ExactArgsWithUsage(1, "conversation-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			conversationID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid conversation ID: %s", args[0])
			}
			opts.ConversationID = conversationID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runConversationsAddRecipients(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Recipients, "recipients", "", "Recipient IDs to add (comma-separated, required)")
	cmd.MarkFlagRequired("recipients")

	return cmd
}

func newConversationsArchiveCmd() *cobra.Command {
	opts := &options.ConversationsArchiveOptions{}

	cmd := &cobra.Command{
		Use:   "archive <conversation-id>",
		Short: "Archive a conversation",
		Long: `Archive a conversation to remove it from the inbox.

Examples:
  canvas conversations archive 123`,
		Args: ExactArgsWithUsage(1, "conversation-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			conversationID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid conversation ID: %s", args[0])
			}
			opts.ConversationID = conversationID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runConversationsArchive(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newConversationsUnarchiveCmd() *cobra.Command {
	opts := &options.ConversationsUnarchiveOptions{}

	cmd := &cobra.Command{
		Use:   "unarchive <conversation-id>",
		Short: "Unarchive a conversation",
		Long: `Move an archived conversation back to the inbox.

Examples:
  canvas conversations unarchive 123`,
		Args: ExactArgsWithUsage(1, "conversation-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			conversationID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid conversation ID: %s", args[0])
			}
			opts.ConversationID = conversationID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runConversationsUnarchive(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newConversationsStarCmd() *cobra.Command {
	opts := &options.ConversationsStarOptions{}

	cmd := &cobra.Command{
		Use:   "star <conversation-id>",
		Short: "Star a conversation",
		Long: `Star a conversation to mark it as important.

Examples:
  canvas conversations star 123`,
		Args: ExactArgsWithUsage(1, "conversation-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			conversationID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid conversation ID: %s", args[0])
			}
			opts.ConversationID = conversationID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runConversationsStar(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newConversationsUnstarCmd() *cobra.Command {
	opts := &options.ConversationsUnstarOptions{}

	cmd := &cobra.Command{
		Use:   "unstar <conversation-id>",
		Short: "Unstar a conversation",
		Long: `Remove the star from a conversation.

Examples:
  canvas conversations unstar 123`,
		Args: ExactArgsWithUsage(1, "conversation-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			conversationID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid conversation ID: %s", args[0])
			}
			opts.ConversationID = conversationID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runConversationsUnstar(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newConversationsMarkReadCmd() *cobra.Command {
	opts := &options.ConversationsMarkReadOptions{}

	cmd := &cobra.Command{
		Use:   "mark-read <conversation-id>",
		Short: "Mark a conversation as read",
		Long: `Mark a conversation as read.

Examples:
  canvas conversations mark-read 123`,
		Args: ExactArgsWithUsage(1, "conversation-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			conversationID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid conversation ID: %s", args[0])
			}
			opts.ConversationID = conversationID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runConversationsMarkRead(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newConversationsMarkAllReadCmd() *cobra.Command {
	opts := &options.ConversationsMarkAllReadOptions{}

	cmd := &cobra.Command{
		Use:   "mark-all-read",
		Short: "Mark all conversations as read",
		Long: `Mark all conversations as read.

Examples:
  canvas conversations mark-all-read`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runConversationsMarkAllRead(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newConversationsDeleteCmd() *cobra.Command {
	opts := &options.ConversationsDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <conversation-id>",
		Short: "Delete a conversation",
		Long: `Delete a conversation permanently.

Examples:
  canvas conversations delete 123
  canvas conversations delete 123 --force`,
		Args: ExactArgsWithUsage(1, "conversation-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			conversationID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid conversation ID: %s", args[0])
			}
			opts.ConversationID = conversationID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runConversationsDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Force, "force", false, "Skip confirmation prompt")

	return cmd
}

func newConversationsUnreadCountCmd() *cobra.Command {
	opts := &options.ConversationsUnreadCountOptions{}

	cmd := &cobra.Command{
		Use:   "unread-count",
		Short: "Get unread conversation count",
		Long: `Get the number of unread conversations.

Examples:
  canvas conversations unread-count`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runConversationsUnreadCount(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func runConversationsList(ctx context.Context, client *api.Client, opts *options.ConversationsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "conversations.list", map[string]interface{}{
		"scope":       opts.Scope,
		"filter":      opts.Filter,
		"filter_mode": opts.FilterMode,
	})

	service := api.NewConversationsService(client)

	listOpts := &api.ListConversationsOptions{
		Scope:                 opts.Scope,
		Filter:                opts.Filter,
		FilterMode:            opts.FilterMode,
		InterleaveSubmissions: opts.InterleaveSubmissions,
	}

	conversations, err := service.List(ctx, listOpts)
	if err != nil {
		logger.LogCommandError(ctx, "conversations.list", err, map[string]interface{}{})
		return fmt.Errorf("failed to list conversations: %w", err)
	}

	if len(conversations) == 0 {
		fmt.Println("No conversations found")
		logger.LogCommandComplete(ctx, "conversations.list", 0)
		return nil
	}

	printVerbose("Found %d conversations:\n\n", len(conversations))
	logger.LogCommandComplete(ctx, "conversations.list", len(conversations))
	return formatOutput(conversations, nil)
}

func runConversationsGet(ctx context.Context, client *api.Client, opts *options.ConversationsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "conversations.get", map[string]interface{}{
		"conversation_id": opts.ConversationID,
		"auto_mark_read":  opts.AutoMarkRead,
	})

	service := api.NewConversationsService(client)

	conversation, err := service.Get(ctx, opts.ConversationID, opts.AutoMarkRead)
	if err != nil {
		logger.LogCommandError(ctx, "conversations.get", err, map[string]interface{}{
			"conversation_id": opts.ConversationID,
		})
		return fmt.Errorf("failed to get conversation: %w", err)
	}

	logger.LogCommandComplete(ctx, "conversations.get", 1)
	return formatOutput(conversation, nil)
}

func runConversationsCreate(ctx context.Context, client *api.Client, opts *options.ConversationsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "conversations.create", map[string]interface{}{
		"recipients": opts.Recipients,
		"subject":    opts.Subject,
	})

	service := api.NewConversationsService(client)

	recipients := api.ParseRecipients(opts.Recipients)
	if len(recipients) == 0 {
		err := fmt.Errorf("at least one recipient is required")
		logger.LogCommandError(ctx, "conversations.create", err, map[string]interface{}{})
		return err
	}

	params := &api.CreateConversationParams{
		Recipients:        recipients,
		Subject:           opts.Subject,
		Body:              opts.Body,
		ForceNew:          opts.ForceNew,
		GroupConversation: opts.GroupConversation,
		AttachmentIDs:     opts.AttachmentIDs,
		MediaCommentID:    opts.MediaCommentID,
		ContextCode:       opts.ContextCode,
	}

	conversations, err := service.Create(ctx, params)
	if err != nil {
		logger.LogCommandError(ctx, "conversations.create", err, map[string]interface{}{})
		return fmt.Errorf("failed to create conversation: %w", err)
	}

	logger.LogCommandComplete(ctx, "conversations.create", len(conversations))

	if len(conversations) > 0 {
		return formatSuccessOutput(conversations[0], fmt.Sprintf("Conversation created successfully (ID: %d)", conversations[0].ID))
	}

	return formatSuccessOutput(conversations, "Conversation created successfully")
}

func runConversationsReply(ctx context.Context, client *api.Client, opts *options.ConversationsReplyOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "conversations.reply", map[string]interface{}{
		"conversation_id": opts.ConversationID,
	})

	service := api.NewConversationsService(client)

	params := &api.AddMessageParams{
		Body:             opts.Body,
		AttachmentIDs:    opts.AttachmentIDs,
		MediaCommentID:   opts.MediaCommentID,
		IncludedMessages: opts.IncludedMessages,
	}

	conversation, err := service.AddMessage(ctx, opts.ConversationID, params)
	if err != nil {
		logger.LogCommandError(ctx, "conversations.reply", err, map[string]interface{}{
			"conversation_id": opts.ConversationID,
		})
		return fmt.Errorf("failed to reply: %w", err)
	}

	logger.LogCommandComplete(ctx, "conversations.reply", 1)

	fmt.Printf("Reply added successfully (conversation ID: %d)\n", conversation.ID)
	return formatOutput(conversation, nil)
}

func runConversationsAddRecipients(ctx context.Context, client *api.Client, opts *options.ConversationsAddRecipientsOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "conversations.add_recipients", map[string]interface{}{
		"conversation_id": opts.ConversationID,
		"recipients":      opts.Recipients,
	})

	service := api.NewConversationsService(client)

	recipients := api.ParseRecipients(opts.Recipients)
	if len(recipients) == 0 {
		err := fmt.Errorf("at least one recipient is required")
		logger.LogCommandError(ctx, "conversations.add_recipients", err, map[string]interface{}{})
		return err
	}

	conversation, err := service.AddRecipients(ctx, opts.ConversationID, recipients)
	if err != nil {
		logger.LogCommandError(ctx, "conversations.add_recipients", err, map[string]interface{}{
			"conversation_id": opts.ConversationID,
		})
		return fmt.Errorf("failed to add recipients: %w", err)
	}

	logger.LogCommandComplete(ctx, "conversations.add_recipients", 1)

	fmt.Printf("Recipients added to conversation %d\n", conversation.ID)
	return formatOutput(conversation, nil)
}

func runConversationsArchive(ctx context.Context, client *api.Client, opts *options.ConversationsArchiveOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "conversations.archive", map[string]interface{}{
		"conversation_id": opts.ConversationID,
	})

	service := api.NewConversationsService(client)

	conversation, err := service.Archive(ctx, opts.ConversationID)
	if err != nil {
		logger.LogCommandError(ctx, "conversations.archive", err, map[string]interface{}{
			"conversation_id": opts.ConversationID,
		})
		return fmt.Errorf("failed to archive conversation: %w", err)
	}

	logger.LogCommandComplete(ctx, "conversations.archive", 1)

	fmt.Printf("Conversation %d archived\n", conversation.ID)
	return nil
}

func runConversationsUnarchive(ctx context.Context, client *api.Client, opts *options.ConversationsUnarchiveOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "conversations.unarchive", map[string]interface{}{
		"conversation_id": opts.ConversationID,
	})

	service := api.NewConversationsService(client)

	conversation, err := service.Unarchive(ctx, opts.ConversationID)
	if err != nil {
		logger.LogCommandError(ctx, "conversations.unarchive", err, map[string]interface{}{
			"conversation_id": opts.ConversationID,
		})
		return fmt.Errorf("failed to unarchive conversation: %w", err)
	}

	logger.LogCommandComplete(ctx, "conversations.unarchive", 1)

	fmt.Printf("Conversation %d moved to inbox\n", conversation.ID)
	return nil
}

func runConversationsStar(ctx context.Context, client *api.Client, opts *options.ConversationsStarOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "conversations.star", map[string]interface{}{
		"conversation_id": opts.ConversationID,
	})

	service := api.NewConversationsService(client)

	conversation, err := service.Star(ctx, opts.ConversationID)
	if err != nil {
		logger.LogCommandError(ctx, "conversations.star", err, map[string]interface{}{
			"conversation_id": opts.ConversationID,
		})
		return fmt.Errorf("failed to star conversation: %w", err)
	}

	logger.LogCommandComplete(ctx, "conversations.star", 1)

	fmt.Printf("Conversation %d starred\n", conversation.ID)
	return nil
}

func runConversationsUnstar(ctx context.Context, client *api.Client, opts *options.ConversationsUnstarOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "conversations.unstar", map[string]interface{}{
		"conversation_id": opts.ConversationID,
	})

	service := api.NewConversationsService(client)

	conversation, err := service.Unstar(ctx, opts.ConversationID)
	if err != nil {
		logger.LogCommandError(ctx, "conversations.unstar", err, map[string]interface{}{
			"conversation_id": opts.ConversationID,
		})
		return fmt.Errorf("failed to unstar conversation: %w", err)
	}

	logger.LogCommandComplete(ctx, "conversations.unstar", 1)

	fmt.Printf("Conversation %d unstarred\n", conversation.ID)
	return nil
}

func runConversationsMarkRead(ctx context.Context, client *api.Client, opts *options.ConversationsMarkReadOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "conversations.mark_read", map[string]interface{}{
		"conversation_id": opts.ConversationID,
	})

	service := api.NewConversationsService(client)

	conversation, err := service.MarkAsRead(ctx, opts.ConversationID)
	if err != nil {
		logger.LogCommandError(ctx, "conversations.mark_read", err, map[string]interface{}{
			"conversation_id": opts.ConversationID,
		})
		return fmt.Errorf("failed to mark as read: %w", err)
	}

	logger.LogCommandComplete(ctx, "conversations.mark_read", 1)

	fmt.Printf("Conversation %d marked as read\n", conversation.ID)
	return nil
}

func runConversationsMarkAllRead(ctx context.Context, client *api.Client, opts *options.ConversationsMarkAllReadOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "conversations.mark_all_read", map[string]interface{}{})

	service := api.NewConversationsService(client)

	err := service.MarkAllAsRead(ctx)
	if err != nil {
		logger.LogCommandError(ctx, "conversations.mark_all_read", err, map[string]interface{}{})
		return fmt.Errorf("failed to mark all as read: %w", err)
	}

	logger.LogCommandComplete(ctx, "conversations.mark_all_read", 0)

	fmt.Println("All conversations marked as read")
	return nil
}

func runConversationsDelete(ctx context.Context, client *api.Client, opts *options.ConversationsDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "conversations.delete", map[string]interface{}{
		"conversation_id": opts.ConversationID,
		"force":           opts.Force,
	})

	if !opts.Force {
		fmt.Printf("WARNING: This will permanently delete conversation %d.\n", opts.ConversationID)
		fmt.Print("Type 'yes' to confirm: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Delete cancelled")
			return nil
		}
	}

	service := api.NewConversationsService(client)

	_, err := service.Delete(ctx, opts.ConversationID)
	if err != nil {
		logger.LogCommandError(ctx, "conversations.delete", err, map[string]interface{}{
			"conversation_id": opts.ConversationID,
		})
		return fmt.Errorf("failed to delete conversation: %w", err)
	}

	logger.LogCommandComplete(ctx, "conversations.delete", 1)

	fmt.Printf("Conversation %d deleted\n", opts.ConversationID)
	return nil
}

func runConversationsUnreadCount(ctx context.Context, client *api.Client, opts *options.ConversationsUnreadCountOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "conversations.unread_count", map[string]interface{}{})

	service := api.NewConversationsService(client)

	count, err := service.GetUnreadCount(ctx)
	if err != nil {
		logger.LogCommandError(ctx, "conversations.unread_count", err, map[string]interface{}{})
		return fmt.Errorf("failed to get unread count: %w", err)
	}

	logger.LogCommandComplete(ctx, "conversations.unread_count", 1)

	fmt.Printf("Unread conversations: %s\n", count.UnreadCount)
	return nil
}
