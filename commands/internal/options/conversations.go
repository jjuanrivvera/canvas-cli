package options

import "fmt"

// ConversationsListOptions contains options for listing conversations
type ConversationsListOptions struct {
	Scope                 string
	Filter                []string
	FilterMode            string
	InterleaveSubmissions bool
}

// Validate validates the options
func (o *ConversationsListOptions) Validate() error {
	return nil
}

// ConversationsGetOptions contains options for getting a conversation
type ConversationsGetOptions struct {
	ConversationID int64
	AutoMarkRead   bool
}

// Validate validates the options
func (o *ConversationsGetOptions) Validate() error {
	if o.ConversationID <= 0 {
		return fmt.Errorf("conversation-id is required and must be greater than 0")
	}
	return nil
}

// ConversationsCreateOptions contains options for creating a conversation
type ConversationsCreateOptions struct {
	Recipients        string
	Subject           string
	Body              string
	ForceNew          bool
	GroupConversation bool
	AttachmentIDs     []int64
	MediaCommentID    string
	ContextCode       string
}

// Validate validates the options
func (o *ConversationsCreateOptions) Validate() error {
	if o.Recipients == "" {
		return fmt.Errorf("recipients are required")
	}
	if o.Body == "" {
		return fmt.Errorf("body is required")
	}
	return nil
}

// ConversationsReplyOptions contains options for replying to a conversation
type ConversationsReplyOptions struct {
	ConversationID   int64
	Body             string
	AttachmentIDs    []int64
	MediaCommentID   string
	IncludedMessages []int64
}

// Validate validates the options
func (o *ConversationsReplyOptions) Validate() error {
	if o.ConversationID <= 0 {
		return fmt.Errorf("conversation-id is required and must be greater than 0")
	}
	if o.Body == "" {
		return fmt.Errorf("body is required")
	}
	return nil
}

// ConversationsAddRecipientsOptions contains options for adding recipients
type ConversationsAddRecipientsOptions struct {
	ConversationID int64
	Recipients     string
}

// Validate validates the options
func (o *ConversationsAddRecipientsOptions) Validate() error {
	if o.ConversationID <= 0 {
		return fmt.Errorf("conversation-id is required and must be greater than 0")
	}
	if o.Recipients == "" {
		return fmt.Errorf("recipients are required")
	}
	return nil
}

// ConversationsArchiveOptions contains options for archiving a conversation
type ConversationsArchiveOptions struct {
	ConversationID int64
}

// Validate validates the options
func (o *ConversationsArchiveOptions) Validate() error {
	if o.ConversationID <= 0 {
		return fmt.Errorf("conversation-id is required and must be greater than 0")
	}
	return nil
}

// ConversationsUnarchiveOptions contains options for unarchiving a conversation
type ConversationsUnarchiveOptions struct {
	ConversationID int64
}

// Validate validates the options
func (o *ConversationsUnarchiveOptions) Validate() error {
	if o.ConversationID <= 0 {
		return fmt.Errorf("conversation-id is required and must be greater than 0")
	}
	return nil
}

// ConversationsStarOptions contains options for starring a conversation
type ConversationsStarOptions struct {
	ConversationID int64
}

// Validate validates the options
func (o *ConversationsStarOptions) Validate() error {
	if o.ConversationID <= 0 {
		return fmt.Errorf("conversation-id is required and must be greater than 0")
	}
	return nil
}

// ConversationsUnstarOptions contains options for unstarring a conversation
type ConversationsUnstarOptions struct {
	ConversationID int64
}

// Validate validates the options
func (o *ConversationsUnstarOptions) Validate() error {
	if o.ConversationID <= 0 {
		return fmt.Errorf("conversation-id is required and must be greater than 0")
	}
	return nil
}

// ConversationsMarkReadOptions contains options for marking a conversation as read
type ConversationsMarkReadOptions struct {
	ConversationID int64
}

// Validate validates the options
func (o *ConversationsMarkReadOptions) Validate() error {
	if o.ConversationID <= 0 {
		return fmt.Errorf("conversation-id is required and must be greater than 0")
	}
	return nil
}

// ConversationsMarkAllReadOptions contains options for marking all conversations as read
type ConversationsMarkAllReadOptions struct {
	// No options needed
}

// Validate validates the options
func (o *ConversationsMarkAllReadOptions) Validate() error {
	return nil
}

// ConversationsDeleteOptions contains options for deleting a conversation
type ConversationsDeleteOptions struct {
	ConversationID int64
	Force          bool
}

// Validate validates the options
func (o *ConversationsDeleteOptions) Validate() error {
	if o.ConversationID <= 0 {
		return fmt.Errorf("conversation-id is required and must be greater than 0")
	}
	return nil
}

// ConversationsUnreadCountOptions contains options for getting unread count
type ConversationsUnreadCountOptions struct {
	// No options needed
}

// Validate validates the options
func (o *ConversationsUnreadCountOptions) Validate() error {
	return nil
}
