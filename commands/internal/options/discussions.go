package options

import "fmt"

// DiscussionsListOptions contains options for listing discussions
type DiscussionsListOptions struct {
	CourseID          int64
	OrderBy           string
	Scope             string
	OnlyAnnouncements bool
	FilterBy          string
	SearchTerm        string
	Include           []string
}

// Validate validates the options
func (o *DiscussionsListOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	return nil
}

// DiscussionsGetOptions contains options for getting a discussion
type DiscussionsGetOptions struct {
	CourseID int64
	TopicID  int64
	Include  []string
}

// Validate validates the options
func (o *DiscussionsGetOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.TopicID <= 0 {
		return fmt.Errorf("topic-id is required and must be greater than 0")
	}
	return nil
}

// DiscussionsCreateOptions contains options for creating a discussion
type DiscussionsCreateOptions struct {
	CourseID           int64
	Title              string
	Message            string
	DiscussionType     string
	Published          bool
	DelayedPostAt      string
	AllowRating        bool
	LockAt             string
	RequireInitialPost bool
	Pinned             bool
}

// Validate validates the options
func (o *DiscussionsCreateOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.Title == "" {
		return fmt.Errorf("title is required")
	}
	return nil
}

// DiscussionsUpdateOptions contains options for updating a discussion
type DiscussionsUpdateOptions struct {
	CourseID           int64
	TopicID            int64
	Title              string
	Message            string
	DiscussionType     string
	Published          bool
	DelayedPostAt      string
	AllowRating        bool
	LockAt             string
	RequireInitialPost bool
	Pinned             bool
	Locked             bool
	// Track which fields were actually set
	TitleSet              bool
	MessageSet            bool
	DiscussionTypeSet     bool
	PublishedSet          bool
	DelayedPostAtSet      bool
	AllowRatingSet        bool
	LockAtSet             bool
	RequireInitialPostSet bool
	PinnedSet             bool
	LockedSet             bool
}

// Validate validates the options
func (o *DiscussionsUpdateOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.TopicID <= 0 {
		return fmt.Errorf("topic-id is required and must be greater than 0")
	}
	return nil
}

// DiscussionsDeleteOptions contains options for deleting a discussion
type DiscussionsDeleteOptions struct {
	CourseID int64
	TopicID  int64
	Force    bool
}

// Validate validates the options
func (o *DiscussionsDeleteOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.TopicID <= 0 {
		return fmt.Errorf("topic-id is required and must be greater than 0")
	}
	return nil
}

// DiscussionsEntriesOptions contains options for listing discussion entries
type DiscussionsEntriesOptions struct {
	CourseID int64
	TopicID  int64
}

// Validate validates the options
func (o *DiscussionsEntriesOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.TopicID <= 0 {
		return fmt.Errorf("topic-id is required and must be greater than 0")
	}
	return nil
}

// DiscussionsPostOptions contains options for posting an entry
type DiscussionsPostOptions struct {
	CourseID int64
	TopicID  int64
	Message  string
}

// Validate validates the options
func (o *DiscussionsPostOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.TopicID <= 0 {
		return fmt.Errorf("topic-id is required and must be greater than 0")
	}
	if o.Message == "" {
		return fmt.Errorf("message is required")
	}
	return nil
}

// DiscussionsReplyOptions contains options for replying to an entry
type DiscussionsReplyOptions struct {
	CourseID int64
	TopicID  int64
	EntryID  int64
	Message  string
}

// Validate validates the options
func (o *DiscussionsReplyOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.TopicID <= 0 {
		return fmt.Errorf("topic-id is required and must be greater than 0")
	}
	if o.EntryID <= 0 {
		return fmt.Errorf("entry-id is required and must be greater than 0")
	}
	if o.Message == "" {
		return fmt.Errorf("message is required")
	}
	return nil
}

// DiscussionsSubscribeOptions contains options for subscribing to a discussion
type DiscussionsSubscribeOptions struct {
	CourseID int64
	TopicID  int64
}

// Validate validates the options
func (o *DiscussionsSubscribeOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.TopicID <= 0 {
		return fmt.Errorf("topic-id is required and must be greater than 0")
	}
	return nil
}

// DiscussionsUnsubscribeOptions contains options for unsubscribing from a discussion
type DiscussionsUnsubscribeOptions struct {
	CourseID int64
	TopicID  int64
}

// Validate validates the options
func (o *DiscussionsUnsubscribeOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.TopicID <= 0 {
		return fmt.Errorf("topic-id is required and must be greater than 0")
	}
	return nil
}
