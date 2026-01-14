package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var (
	// List flags
	conversationsScope      string
	conversationsFilter     []string
	conversationsFilterMode string
	conversationsInterleave bool

	// Create flags
	conversationsRecipients     string
	conversationsSubject        string
	conversationsBody           string
	conversationsForceNew       bool
	conversationsGroup          bool
	conversationsAttachmentIDs  []int64
	conversationsMediaCommentID string
	conversationsContextCode    string

	// Reply flags
	conversationsIncludedMessages []int64

	// General flags
	conversationsAutoMarkRead bool
	conversationsForce        bool
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

// conversationsListCmd lists conversations
var conversationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List conversations",
	Long: `List conversations for the current user.

Examples:
  canvas conversations list
  canvas conversations list --scope unread
  canvas conversations list --scope starred
  canvas conversations list --filter course_123`,
	RunE: runConversationsList,
}

// conversationsGetCmd gets a single conversation
var conversationsGetCmd = &cobra.Command{
	Use:   "get <conversation-id>",
	Short: "Get conversation details",
	Long: `Get details of a specific conversation including messages.

Examples:
  canvas conversations get 123
  canvas conversations get 123 --auto-mark-read=false`,
	Args: ExactArgsWithUsage(1, "conversation-id"),
	RunE: runConversationsGet,
}

// conversationsCreateCmd creates a new conversation
var conversationsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new conversation",
	Long: `Create a new conversation with one or more recipients.

Examples:
  canvas conversations create --recipients 123,456 --subject "Hello" --body "Message text"
  canvas conversations create --recipients 123 --body "Quick message"
  canvas conversations create --recipients course_123 --subject "Announcement" --body "Message" --group`,
	RunE: runConversationsCreate,
}

// conversationsReplyCmd adds a reply to a conversation
var conversationsReplyCmd = &cobra.Command{
	Use:   "reply <conversation-id>",
	Short: "Reply to a conversation",
	Long: `Add a reply message to an existing conversation.

Examples:
  canvas conversations reply 123 --body "Thank you for your message"
  canvas conversations reply 123 --body "Here's the attachment" --attachment-ids 456`,
	Args: ExactArgsWithUsage(1, "conversation-id"),
	RunE: runConversationsReply,
}

// conversationsAddRecipientsCmd adds recipients to a conversation
var conversationsAddRecipientsCmd = &cobra.Command{
	Use:   "add-recipients <conversation-id>",
	Short: "Add recipients to a conversation",
	Long: `Add additional recipients to an existing conversation.

Examples:
  canvas conversations add-recipients 123 --recipients 456,789`,
	Args: ExactArgsWithUsage(1, "conversation-id"),
	RunE: runConversationsAddRecipients,
}

// conversationsArchiveCmd archives a conversation
var conversationsArchiveCmd = &cobra.Command{
	Use:   "archive <conversation-id>",
	Short: "Archive a conversation",
	Long: `Archive a conversation to remove it from the inbox.

Examples:
  canvas conversations archive 123`,
	Args: ExactArgsWithUsage(1, "conversation-id"),
	RunE: runConversationsArchive,
}

// conversationsUnarchiveCmd unarchives a conversation
var conversationsUnarchiveCmd = &cobra.Command{
	Use:   "unarchive <conversation-id>",
	Short: "Unarchive a conversation",
	Long: `Move an archived conversation back to the inbox.

Examples:
  canvas conversations unarchive 123`,
	Args: ExactArgsWithUsage(1, "conversation-id"),
	RunE: runConversationsUnarchive,
}

// conversationsStarCmd stars a conversation
var conversationsStarCmd = &cobra.Command{
	Use:   "star <conversation-id>",
	Short: "Star a conversation",
	Long: `Star a conversation to mark it as important.

Examples:
  canvas conversations star 123`,
	Args: ExactArgsWithUsage(1, "conversation-id"),
	RunE: runConversationsStar,
}

// conversationsUnstarCmd unstars a conversation
var conversationsUnstarCmd = &cobra.Command{
	Use:   "unstar <conversation-id>",
	Short: "Unstar a conversation",
	Long: `Remove the star from a conversation.

Examples:
  canvas conversations unstar 123`,
	Args: ExactArgsWithUsage(1, "conversation-id"),
	RunE: runConversationsUnstar,
}

// conversationsMarkReadCmd marks a conversation as read
var conversationsMarkReadCmd = &cobra.Command{
	Use:   "mark-read <conversation-id>",
	Short: "Mark a conversation as read",
	Long: `Mark a conversation as read.

Examples:
  canvas conversations mark-read 123`,
	Args: ExactArgsWithUsage(1, "conversation-id"),
	RunE: runConversationsMarkRead,
}

// conversationsMarkAllReadCmd marks all conversations as read
var conversationsMarkAllReadCmd = &cobra.Command{
	Use:   "mark-all-read",
	Short: "Mark all conversations as read",
	Long: `Mark all conversations as read.

Examples:
  canvas conversations mark-all-read`,
	RunE: runConversationsMarkAllRead,
}

// conversationsDeleteCmd deletes a conversation
var conversationsDeleteCmd = &cobra.Command{
	Use:   "delete <conversation-id>",
	Short: "Delete a conversation",
	Long: `Delete a conversation permanently.

Examples:
  canvas conversations delete 123
  canvas conversations delete 123 --force`,
	Args: ExactArgsWithUsage(1, "conversation-id"),
	RunE: runConversationsDelete,
}

// conversationsUnreadCountCmd gets the unread count
var conversationsUnreadCountCmd = &cobra.Command{
	Use:   "unread-count",
	Short: "Get unread conversation count",
	Long: `Get the number of unread conversations.

Examples:
  canvas conversations unread-count`,
	RunE: runConversationsUnreadCount,
}

func init() {
	rootCmd.AddCommand(conversationsCmd)
	conversationsCmd.AddCommand(conversationsListCmd)
	conversationsCmd.AddCommand(conversationsGetCmd)
	conversationsCmd.AddCommand(conversationsCreateCmd)
	conversationsCmd.AddCommand(conversationsReplyCmd)
	conversationsCmd.AddCommand(conversationsAddRecipientsCmd)
	conversationsCmd.AddCommand(conversationsArchiveCmd)
	conversationsCmd.AddCommand(conversationsUnarchiveCmd)
	conversationsCmd.AddCommand(conversationsStarCmd)
	conversationsCmd.AddCommand(conversationsUnstarCmd)
	conversationsCmd.AddCommand(conversationsMarkReadCmd)
	conversationsCmd.AddCommand(conversationsMarkAllReadCmd)
	conversationsCmd.AddCommand(conversationsDeleteCmd)
	conversationsCmd.AddCommand(conversationsUnreadCountCmd)

	// List flags
	conversationsListCmd.Flags().StringVar(&conversationsScope, "scope", "", "Scope: inbox, unread, archived, starred, sent")
	conversationsListCmd.Flags().StringSliceVar(&conversationsFilter, "filter", nil, "Filter by course/group (e.g., course_123, group_456)")
	conversationsListCmd.Flags().StringVar(&conversationsFilterMode, "filter-mode", "", "Filter mode: and, or")
	conversationsListCmd.Flags().BoolVar(&conversationsInterleave, "interleave", false, "Interleave submissions")

	// Get flags
	conversationsGetCmd.Flags().BoolVar(&conversationsAutoMarkRead, "auto-mark-read", true, "Auto mark as read")

	// Create flags
	conversationsCreateCmd.Flags().StringVar(&conversationsRecipients, "recipients", "", "Recipient IDs (comma-separated, required)")
	conversationsCreateCmd.Flags().StringVar(&conversationsSubject, "subject", "", "Subject line")
	conversationsCreateCmd.Flags().StringVar(&conversationsBody, "body", "", "Message body (required)")
	conversationsCreateCmd.Flags().BoolVar(&conversationsForceNew, "force-new", false, "Force new conversation")
	conversationsCreateCmd.Flags().BoolVar(&conversationsGroup, "group", false, "Send as group message")
	conversationsCreateCmd.Flags().Int64SliceVar(&conversationsAttachmentIDs, "attachment-ids", nil, "Attachment IDs")
	conversationsCreateCmd.Flags().StringVar(&conversationsMediaCommentID, "media-comment-id", "", "Media comment ID")
	conversationsCreateCmd.Flags().StringVar(&conversationsContextCode, "context-code", "", "Context code (e.g., course_123)")
	conversationsCreateCmd.MarkFlagRequired("recipients")
	conversationsCreateCmd.MarkFlagRequired("body")

	// Reply flags
	conversationsReplyCmd.Flags().StringVar(&conversationsBody, "body", "", "Reply body (required)")
	conversationsReplyCmd.Flags().Int64SliceVar(&conversationsAttachmentIDs, "attachment-ids", nil, "Attachment IDs")
	conversationsReplyCmd.Flags().StringVar(&conversationsMediaCommentID, "media-comment-id", "", "Media comment ID")
	conversationsReplyCmd.Flags().Int64SliceVar(&conversationsIncludedMessages, "included-messages", nil, "Include previous message IDs")
	conversationsReplyCmd.MarkFlagRequired("body")

	// Add recipients flags
	conversationsAddRecipientsCmd.Flags().StringVar(&conversationsRecipients, "recipients", "", "Recipient IDs to add (comma-separated, required)")
	conversationsAddRecipientsCmd.MarkFlagRequired("recipients")

	// Delete flags
	conversationsDeleteCmd.Flags().BoolVar(&conversationsForce, "force", false, "Skip confirmation prompt")
}

func runConversationsList(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewConversationsService(client)

	opts := &api.ListConversationsOptions{
		Scope:                 conversationsScope,
		Filter:                conversationsFilter,
		FilterMode:            conversationsFilterMode,
		InterleaveSubmissions: conversationsInterleave,
	}

	conversations, err := service.List(context.Background(), opts)
	if err != nil {
		return fmt.Errorf("failed to list conversations: %w", err)
	}

	if len(conversations) == 0 {
		fmt.Println("No conversations found")
		return nil
	}

	printVerbose("Found %d conversations:\n\n", len(conversations))
	return formatOutput(conversations, nil)
}

func runConversationsGet(cmd *cobra.Command, args []string) error {
	conversationID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid conversation ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewConversationsService(client)

	conversation, err := service.Get(context.Background(), conversationID, conversationsAutoMarkRead)
	if err != nil {
		return fmt.Errorf("failed to get conversation: %w", err)
	}

	return formatOutput(conversation, nil)
}

func runConversationsCreate(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewConversationsService(client)

	recipients := api.ParseRecipients(conversationsRecipients)
	if len(recipients) == 0 {
		return fmt.Errorf("at least one recipient is required")
	}

	params := &api.CreateConversationParams{
		Recipients:        recipients,
		Subject:           conversationsSubject,
		Body:              conversationsBody,
		ForceNew:          conversationsForceNew,
		GroupConversation: conversationsGroup,
		AttachmentIDs:     conversationsAttachmentIDs,
		MediaCommentID:    conversationsMediaCommentID,
		ContextCode:       conversationsContextCode,
	}

	conversations, err := service.Create(context.Background(), params)
	if err != nil {
		return fmt.Errorf("failed to create conversation: %w", err)
	}

	if len(conversations) > 0 {
		return formatSuccessOutput(conversations[0], fmt.Sprintf("Conversation created successfully (ID: %d)", conversations[0].ID))
	}

	return formatSuccessOutput(conversations, "Conversation created successfully")
}

func runConversationsReply(cmd *cobra.Command, args []string) error {
	conversationID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid conversation ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewConversationsService(client)

	params := &api.AddMessageParams{
		Body:             conversationsBody,
		AttachmentIDs:    conversationsAttachmentIDs,
		MediaCommentID:   conversationsMediaCommentID,
		IncludedMessages: conversationsIncludedMessages,
	}

	conversation, err := service.AddMessage(context.Background(), conversationID, params)
	if err != nil {
		return fmt.Errorf("failed to reply: %w", err)
	}

	fmt.Printf("Reply added successfully (conversation ID: %d)\n", conversation.ID)
	return formatOutput(conversation, nil)
}

func runConversationsAddRecipients(cmd *cobra.Command, args []string) error {
	conversationID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid conversation ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewConversationsService(client)

	recipients := api.ParseRecipients(conversationsRecipients)
	if len(recipients) == 0 {
		return fmt.Errorf("at least one recipient is required")
	}

	conversation, err := service.AddRecipients(context.Background(), conversationID, recipients)
	if err != nil {
		return fmt.Errorf("failed to add recipients: %w", err)
	}

	fmt.Printf("Recipients added to conversation %d\n", conversation.ID)
	return formatOutput(conversation, nil)
}

func runConversationsArchive(cmd *cobra.Command, args []string) error {
	conversationID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid conversation ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewConversationsService(client)

	conversation, err := service.Archive(context.Background(), conversationID)
	if err != nil {
		return fmt.Errorf("failed to archive conversation: %w", err)
	}

	fmt.Printf("Conversation %d archived\n", conversation.ID)
	return nil
}

func runConversationsUnarchive(cmd *cobra.Command, args []string) error {
	conversationID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid conversation ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewConversationsService(client)

	conversation, err := service.Unarchive(context.Background(), conversationID)
	if err != nil {
		return fmt.Errorf("failed to unarchive conversation: %w", err)
	}

	fmt.Printf("Conversation %d moved to inbox\n", conversation.ID)
	return nil
}

func runConversationsStar(cmd *cobra.Command, args []string) error {
	conversationID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid conversation ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewConversationsService(client)

	conversation, err := service.Star(context.Background(), conversationID)
	if err != nil {
		return fmt.Errorf("failed to star conversation: %w", err)
	}

	fmt.Printf("Conversation %d starred\n", conversation.ID)
	return nil
}

func runConversationsUnstar(cmd *cobra.Command, args []string) error {
	conversationID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid conversation ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewConversationsService(client)

	conversation, err := service.Unstar(context.Background(), conversationID)
	if err != nil {
		return fmt.Errorf("failed to unstar conversation: %w", err)
	}

	fmt.Printf("Conversation %d unstarred\n", conversation.ID)
	return nil
}

func runConversationsMarkRead(cmd *cobra.Command, args []string) error {
	conversationID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid conversation ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewConversationsService(client)

	conversation, err := service.MarkAsRead(context.Background(), conversationID)
	if err != nil {
		return fmt.Errorf("failed to mark as read: %w", err)
	}

	fmt.Printf("Conversation %d marked as read\n", conversation.ID)
	return nil
}

func runConversationsMarkAllRead(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewConversationsService(client)

	err = service.MarkAllAsRead(context.Background())
	if err != nil {
		return fmt.Errorf("failed to mark all as read: %w", err)
	}

	fmt.Println("All conversations marked as read")
	return nil
}

func runConversationsDelete(cmd *cobra.Command, args []string) error {
	conversationID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid conversation ID: %w", err)
	}

	if !conversationsForce {
		fmt.Printf("WARNING: This will permanently delete conversation %d.\n", conversationID)
		fmt.Print("Type 'yes' to confirm: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Delete cancelled")
			return nil
		}
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewConversationsService(client)

	_, err = service.Delete(context.Background(), conversationID)
	if err != nil {
		return fmt.Errorf("failed to delete conversation: %w", err)
	}

	fmt.Printf("Conversation %d deleted\n", conversationID)
	return nil
}

func runConversationsUnreadCount(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewConversationsService(client)

	count, err := service.GetUnreadCount(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get unread count: %w", err)
	}

	fmt.Printf("Unread conversations: %s\n", count.UnreadCount)
	return nil
}
