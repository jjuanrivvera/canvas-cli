package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// ConversationsService handles conversation-related API calls
type ConversationsService struct {
	client *Client
}

// NewConversationsService creates a new conversations service
func NewConversationsService(client *Client) *ConversationsService {
	return &ConversationsService{client: client}
}

// Conversation represents a Canvas conversation
type Conversation struct {
	ID               int64                 `json:"id"`
	Subject          string                `json:"subject"`
	WorkflowState    string                `json:"workflow_state"`
	LastMessage      string                `json:"last_message,omitempty"`
	StartAt          string                `json:"start_at,omitempty"`
	LastMessageAt    string                `json:"last_message_at,omitempty"`
	MessageCount     int                   `json:"message_count"`
	Subscribed       bool                  `json:"subscribed"`
	Private          bool                  `json:"private"`
	Starred          bool                  `json:"starred"`
	Properties       []string              `json:"properties,omitempty"`
	Audience         []int64               `json:"audience,omitempty"`
	AudienceContexts []AudienceContext     `json:"audience_contexts,omitempty"`
	AvatarURL        string                `json:"avatar_url,omitempty"`
	Participants     []ConversationUser    `json:"participants,omitempty"`
	Visible          bool                  `json:"visible"`
	ContextName      string                `json:"context_name,omitempty"`
	Messages         []ConversationMessage `json:"messages,omitempty"`
}

// AudienceContext represents context information for conversation audience
type AudienceContext struct {
	Courses map[string][]string `json:"courses,omitempty"`
	Groups  map[string][]string `json:"groups,omitempty"`
}

// ConversationUser represents a participant in a conversation
type ConversationUser struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

// ConversationMessage represents a single message in a conversation
type ConversationMessage struct {
	ID                   int64                 `json:"id"`
	CreatedAt            string                `json:"created_at"`
	Body                 string                `json:"body"`
	AuthorID             int64                 `json:"author_id"`
	Generated            bool                  `json:"generated"`
	MediaComment         *MediaComment         `json:"media_comment,omitempty"`
	ForwardedMessages    []ConversationMessage `json:"forwarded_messages,omitempty"`
	Attachments          []Attachment          `json:"attachments,omitempty"`
	ParticipatingUserIds []int64               `json:"participating_user_ids,omitempty"`
}

// UnreadCount represents the unread conversation count
type UnreadCount struct {
	UnreadCount string `json:"unread_count"`
}

// ListConversationsOptions holds options for listing conversations
type ListConversationsOptions struct {
	Scope                 string   // inbox, unread, archived, starred, sent
	Filter                []string // course_123, group_456
	FilterMode            string   // and, or, default
	InterleaveSubmissions bool
	Include               []string
	Page                  int
	PerPage               int
}

// List retrieves conversations for the current user
func (s *ConversationsService) List(ctx context.Context, opts *ListConversationsOptions) ([]Conversation, error) {
	path := "/api/v1/conversations"

	if opts != nil {
		query := url.Values{}

		if opts.Scope != "" {
			query.Add("scope", opts.Scope)
		}

		for _, filter := range opts.Filter {
			query.Add("filter[]", filter)
		}

		if opts.FilterMode != "" {
			query.Add("filter_mode", opts.FilterMode)
		}

		if opts.InterleaveSubmissions {
			query.Add("interleave_submissions", "true")
		}

		for _, include := range opts.Include {
			query.Add("include[]", include)
		}

		if opts.Page > 0 {
			query.Add("page", strconv.Itoa(opts.Page))
		}

		if opts.PerPage > 0 {
			query.Add("per_page", strconv.Itoa(opts.PerPage))
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var conversations []Conversation
	if err := s.client.GetAllPages(ctx, path, &conversations); err != nil {
		return nil, err
	}

	return conversations, nil
}

// Get retrieves a single conversation
func (s *ConversationsService) Get(ctx context.Context, conversationID int64, autoMarkRead bool) (*Conversation, error) {
	path := fmt.Sprintf("/api/v1/conversations/%d", conversationID)

	query := url.Values{}
	if !autoMarkRead {
		query.Add("auto_mark_as_read", "false")
	}

	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	var conversation Conversation
	if err := s.client.GetJSON(ctx, path, &conversation); err != nil {
		return nil, err
	}

	return &conversation, nil
}

// CreateConversationParams holds parameters for creating a conversation
type CreateConversationParams struct {
	Recipients        []string // User IDs or special strings like "course_123"
	Subject           string
	Body              string
	ForceNew          bool
	GroupConversation bool
	AttachmentIDs     []int64
	MediaCommentID    string
	MediaCommentType  string
	ContextCode       string // e.g., "course_123"
	Mode              string // sync or async
}

// Create creates a new conversation
func (s *ConversationsService) Create(ctx context.Context, params *CreateConversationParams) ([]Conversation, error) {
	path := "/api/v1/conversations"

	body := make(map[string]interface{})

	if len(params.Recipients) > 0 {
		body["recipients"] = params.Recipients
	}

	if params.Subject != "" {
		body["subject"] = params.Subject
	}

	if params.Body != "" {
		body["body"] = params.Body
	}

	if params.ForceNew {
		body["force_new"] = true
	}

	if params.GroupConversation {
		body["group_conversation"] = true
	}

	if len(params.AttachmentIDs) > 0 {
		ids := make([]string, len(params.AttachmentIDs))
		for i, id := range params.AttachmentIDs {
			ids[i] = strconv.FormatInt(id, 10)
		}
		body["attachment_ids"] = ids
	}

	if params.MediaCommentID != "" {
		body["media_comment_id"] = params.MediaCommentID
	}

	if params.MediaCommentType != "" {
		body["media_comment_type"] = params.MediaCommentType
	}

	if params.ContextCode != "" {
		body["context_code"] = params.ContextCode
	}

	if params.Mode != "" {
		body["mode"] = params.Mode
	}

	var conversations []Conversation
	if err := s.client.PostJSON(ctx, path, body, &conversations); err != nil {
		return nil, err
	}

	return conversations, nil
}

// AddMessageParams holds parameters for adding a message to a conversation
type AddMessageParams struct {
	Body             string
	AttachmentIDs    []int64
	MediaCommentID   string
	MediaCommentType string
	IncludedMessages []int64
	Recipients       []string
}

// AddMessage adds a message to an existing conversation
func (s *ConversationsService) AddMessage(ctx context.Context, conversationID int64, params *AddMessageParams) (*Conversation, error) {
	path := fmt.Sprintf("/api/v1/conversations/%d/add_message", conversationID)

	body := make(map[string]interface{})

	if params.Body != "" {
		body["body"] = params.Body
	}

	if len(params.AttachmentIDs) > 0 {
		ids := make([]string, len(params.AttachmentIDs))
		for i, id := range params.AttachmentIDs {
			ids[i] = strconv.FormatInt(id, 10)
		}
		body["attachment_ids"] = ids
	}

	if params.MediaCommentID != "" {
		body["media_comment_id"] = params.MediaCommentID
	}

	if params.MediaCommentType != "" {
		body["media_comment_type"] = params.MediaCommentType
	}

	if len(params.IncludedMessages) > 0 {
		ids := make([]string, len(params.IncludedMessages))
		for i, id := range params.IncludedMessages {
			ids[i] = strconv.FormatInt(id, 10)
		}
		body["included_messages"] = ids
	}

	if len(params.Recipients) > 0 {
		body["recipients"] = params.Recipients
	}

	var conversation Conversation
	if err := s.client.PostJSON(ctx, path, body, &conversation); err != nil {
		return nil, err
	}

	return &conversation, nil
}

// AddRecipients adds recipients to a conversation
func (s *ConversationsService) AddRecipients(ctx context.Context, conversationID int64, recipients []string) (*Conversation, error) {
	path := fmt.Sprintf("/api/v1/conversations/%d/add_recipients", conversationID)

	body := map[string]interface{}{
		"recipients": recipients,
	}

	var conversation Conversation
	if err := s.client.PostJSON(ctx, path, body, &conversation); err != nil {
		return nil, err
	}

	return &conversation, nil
}

// RemoveMessages removes messages from a conversation
func (s *ConversationsService) RemoveMessages(ctx context.Context, conversationID int64, messageIDs []int64) (*Conversation, error) {
	path := fmt.Sprintf("/api/v1/conversations/%d/remove_messages", conversationID)

	ids := make([]string, len(messageIDs))
	for i, id := range messageIDs {
		ids[i] = strconv.FormatInt(id, 10)
	}

	body := map[string]interface{}{
		"remove": ids,
	}

	var conversation Conversation
	if err := s.client.PostJSON(ctx, path, body, &conversation); err != nil {
		return nil, err
	}

	return &conversation, nil
}

// MarkAllAsRead marks all conversations as read
func (s *ConversationsService) MarkAllAsRead(ctx context.Context) error {
	path := "/api/v1/conversations/mark_all_as_read"

	return s.client.PostJSON(ctx, path, nil, nil)
}

// Delete deletes a conversation
func (s *ConversationsService) Delete(ctx context.Context, conversationID int64) (*Conversation, error) {
	path := fmt.Sprintf("/api/v1/conversations/%d", conversationID)

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var conversation Conversation
	if err := json.NewDecoder(resp.Body).Decode(&conversation); err != nil {
		return nil, err
	}

	return &conversation, nil
}

// BatchUpdateParams holds parameters for batch updating conversations
type BatchUpdateParams struct {
	ConversationIDs []int64
	Event           string // mark_as_read, mark_as_unread, star, unstar, archive, destroy
}

// BatchUpdate updates multiple conversations at once
func (s *ConversationsService) BatchUpdate(ctx context.Context, params *BatchUpdateParams) (*BatchUpdateProgress, error) {
	path := "/api/v1/conversations"

	ids := make([]string, len(params.ConversationIDs))
	for i, id := range params.ConversationIDs {
		ids[i] = strconv.FormatInt(id, 10)
	}

	body := map[string]interface{}{
		"conversation_ids": ids,
		"event":            params.Event,
	}

	var progress BatchUpdateProgress
	if err := s.client.PutJSON(ctx, path, body, &progress); err != nil {
		return nil, err
	}

	return &progress, nil
}

// BatchUpdateProgress represents the progress of a batch update operation
type BatchUpdateProgress struct {
	Progress *ConversationProgress `json:"progress,omitempty"`
}

// ConversationProgress represents progress for a conversation batch operation
type ConversationProgress struct {
	ID            int64   `json:"id"`
	ContextID     int64   `json:"context_id"`
	ContextType   string  `json:"context_type"`
	UserID        int64   `json:"user_id"`
	Tag           string  `json:"tag"`
	Completion    float64 `json:"completion"`
	WorkflowState string  `json:"workflow_state"`
	URL           string  `json:"url"`
	Message       string  `json:"message,omitempty"`
}

// GetUnreadCount retrieves the unread conversation count
func (s *ConversationsService) GetUnreadCount(ctx context.Context) (*UnreadCount, error) {
	path := "/api/v1/conversations/unread_count"

	var count UnreadCount
	if err := s.client.GetJSON(ctx, path, &count); err != nil {
		return nil, err
	}

	return &count, nil
}

// UpdateConversation updates a single conversation (for star/archive operations)
func (s *ConversationsService) UpdateConversation(ctx context.Context, conversationID int64, workflowState string, starred *bool, subscribed *bool) (*Conversation, error) {
	path := fmt.Sprintf("/api/v1/conversations/%d", conversationID)

	body := make(map[string]interface{})

	innerBody := make(map[string]interface{})
	if workflowState != "" {
		innerBody["workflow_state"] = workflowState
	}
	if starred != nil {
		innerBody["starred"] = *starred
	}
	if subscribed != nil {
		innerBody["subscribed"] = *subscribed
	}

	body["conversation"] = innerBody

	var conversation Conversation
	if err := s.client.PutJSON(ctx, path, body, &conversation); err != nil {
		return nil, err
	}

	return &conversation, nil
}

// Archive archives a conversation
func (s *ConversationsService) Archive(ctx context.Context, conversationID int64) (*Conversation, error) {
	return s.UpdateConversation(ctx, conversationID, "archived", nil, nil)
}

// Unarchive unarchives a conversation (moves to inbox)
func (s *ConversationsService) Unarchive(ctx context.Context, conversationID int64) (*Conversation, error) {
	return s.UpdateConversation(ctx, conversationID, "read", nil, nil)
}

// Star stars a conversation
func (s *ConversationsService) Star(ctx context.Context, conversationID int64) (*Conversation, error) {
	starred := true
	return s.UpdateConversation(ctx, conversationID, "", &starred, nil)
}

// Unstar unstars a conversation
func (s *ConversationsService) Unstar(ctx context.Context, conversationID int64) (*Conversation, error) {
	starred := false
	return s.UpdateConversation(ctx, conversationID, "", &starred, nil)
}

// MarkAsRead marks a conversation as read
func (s *ConversationsService) MarkAsRead(ctx context.Context, conversationID int64) (*Conversation, error) {
	return s.UpdateConversation(ctx, conversationID, "read", nil, nil)
}

// FindRecipients searches for valid recipients
type RecipientSearchResult struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url,omitempty"`
	Type      string `json:"type"` // user, context, group
	UserCount int    `json:"user_count,omitempty"`
	Pronouns  string `json:"pronouns,omitempty"`
}

// SearchRecipients searches for valid message recipients
func (s *ConversationsService) SearchRecipients(ctx context.Context, search string, contextCode string, perPage int) ([]RecipientSearchResult, error) {
	path := "/api/v1/search/recipients"

	query := url.Values{}
	if search != "" {
		query.Add("search", search)
	}
	if contextCode != "" {
		query.Add("context", contextCode)
	}
	if perPage > 0 {
		query.Add("per_page", strconv.Itoa(perPage))
	}

	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	var results []RecipientSearchResult
	if err := s.client.GetJSON(ctx, path, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// Helper function to format recipients as comma-separated string
func FormatRecipients(ids []int64) []string {
	result := make([]string, len(ids))
	for i, id := range ids {
		result[i] = strconv.FormatInt(id, 10)
	}
	return result
}

// Helper function to parse recipients from comma-separated string
func ParseRecipients(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
