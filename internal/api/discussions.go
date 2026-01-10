package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// DiscussionTopic represents a Canvas discussion topic or announcement
type DiscussionTopic struct {
	ID                      int64                  `json:"id"`
	Title                   string                 `json:"title"`
	Message                 string                 `json:"message"`
	HTMLURL                 string                 `json:"html_url"`
	PostedAt                *time.Time             `json:"posted_at,omitempty"`
	LastReplyAt             *time.Time             `json:"last_reply_at,omitempty"`
	RequireInitialPost      bool                   `json:"require_initial_post"`
	UserCanSeePosts         bool                   `json:"user_can_see_posts"`
	DiscussionSubentryCount int                    `json:"discussion_subentry_count"`
	ReadState               string                 `json:"read_state"`
	UnreadCount             int                    `json:"unread_count"`
	Subscribed              bool                   `json:"subscribed"`
	SubscriptionHold        string                 `json:"subscription_hold,omitempty"`
	AssignmentID            *int64                 `json:"assignment_id,omitempty"`
	DelayedPostAt           *time.Time             `json:"delayed_post_at,omitempty"`
	Published               bool                   `json:"published"`
	LockAt                  *time.Time             `json:"lock_at,omitempty"`
	Locked                  bool                   `json:"locked"`
	Pinned                  bool                   `json:"pinned"`
	LockedForUser           bool                   `json:"locked_for_user"`
	LockInfo                *LockInfo              `json:"lock_info,omitempty"`
	LockExplanation         string                 `json:"lock_explanation,omitempty"`
	UserName                string                 `json:"user_name,omitempty"`
	RootTopicID             *int64                 `json:"root_topic_id,omitempty"`
	PodcastURL              string                 `json:"podcast_url,omitempty"`
	DiscussionType          string                 `json:"discussion_type"`
	GroupCategoryID         *int64                 `json:"group_category_id,omitempty"`
	Attachments             []FileAttachment       `json:"attachments,omitempty"`
	Permissions             map[string]bool        `json:"permissions,omitempty"`
	AllowRating             bool                   `json:"allow_rating"`
	OnlyGradersCanRate      bool                   `json:"only_graders_can_rate"`
	SortByRating            bool                   `json:"sort_by_rating"`
	ContextCode             string                 `json:"context_code,omitempty"`
	Author                  *User                  `json:"author,omitempty"`
	IsAnnouncement          bool                   `json:"is_announcement,omitempty"`
}

// FileAttachment represents a file attachment
type FileAttachment struct {
	ContentType string `json:"content-type"`
	URL         string `json:"url"`
	Filename    string `json:"filename"`
	DisplayName string `json:"display_name"`
}

// DiscussionEntry represents an entry in a discussion
type DiscussionEntry struct {
	ID             int64             `json:"id"`
	UserID         int64             `json:"user_id"`
	ParentID       *int64            `json:"parent_id,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	Message        string            `json:"message"`
	Rating         int               `json:"rating"`
	RatingSum      int               `json:"rating_sum"`
	ReadState      string            `json:"read_state"`
	ForcedReadState bool             `json:"forced_read_state"`
	User           *User             `json:"user,omitempty"`
	Replies        []DiscussionEntry `json:"replies,omitempty"`
}

// DiscussionsService handles discussion-related API calls
type DiscussionsService struct {
	client *Client
}

// NewDiscussionsService creates a new discussions service
func NewDiscussionsService(client *Client) *DiscussionsService {
	return &DiscussionsService{client: client}
}

// ListDiscussionsOptions holds options for listing discussions
type ListDiscussionsOptions struct {
	Include           []string // all_dates, sections, sections_user_count, overrides
	OrderBy           string   // position, recent_activity, title
	Scope             string   // locked, unlocked, pinned, unpinned
	OnlyAnnouncements bool
	FilterBy          string // all, unread
	SearchTerm        string
	Page              int
	PerPage           int
}

// List retrieves all discussion topics for a course
func (s *DiscussionsService) List(ctx context.Context, courseID int64, opts *ListDiscussionsOptions) ([]DiscussionTopic, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/discussion_topics", courseID)

	if opts != nil {
		query := url.Values{}

		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}

		if opts.OrderBy != "" {
			query.Add("order_by", opts.OrderBy)
		}

		if opts.Scope != "" {
			query.Add("scope", opts.Scope)
		}

		if opts.OnlyAnnouncements {
			query.Add("only_announcements", "true")
		}

		if opts.FilterBy != "" {
			query.Add("filter_by", opts.FilterBy)
		}

		if opts.SearchTerm != "" {
			query.Add("search_term", opts.SearchTerm)
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

	var topics []DiscussionTopic
	if err := s.client.GetAllPages(ctx, path, &topics); err != nil {
		return nil, err
	}

	return topics, nil
}

// Get retrieves a single discussion topic
func (s *DiscussionsService) Get(ctx context.Context, courseID, topicID int64, include []string) (*DiscussionTopic, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/discussion_topics/%d", courseID, topicID)

	if len(include) > 0 {
		query := url.Values{}
		for _, inc := range include {
			query.Add("include[]", inc)
		}
		path += "?" + query.Encode()
	}

	var topic DiscussionTopic
	if err := s.client.GetJSON(ctx, path, &topic); err != nil {
		return nil, err
	}

	return &topic, nil
}

// CreateDiscussionParams holds parameters for creating a discussion
type CreateDiscussionParams struct {
	Title                 string
	Message               string
	DiscussionType        string // side_comment, threaded, not_threaded
	Published             bool
	DelayedPostAt         string
	AllowRating           bool
	LockAt                string
	PodcastEnabled        bool
	PodcastHasStudentPosts bool
	RequireInitialPost    bool
	IsAnnouncement        bool
	Pinned                bool
	PositionAfter         string
	GroupCategoryID       int64
	OnlyGradersCanRate    bool
	SpecificSections      string
}

// Create creates a new discussion topic
func (s *DiscussionsService) Create(ctx context.Context, courseID int64, params *CreateDiscussionParams) (*DiscussionTopic, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/discussion_topics", courseID)

	body := make(map[string]interface{})

	if params.Title != "" {
		body["title"] = params.Title
	}

	if params.Message != "" {
		body["message"] = params.Message
	}

	if params.DiscussionType != "" {
		body["discussion_type"] = params.DiscussionType
	}

	body["published"] = params.Published

	if params.DelayedPostAt != "" {
		body["delayed_post_at"] = params.DelayedPostAt
	}

	if params.AllowRating {
		body["allow_rating"] = true
	}

	if params.LockAt != "" {
		body["lock_at"] = params.LockAt
	}

	if params.PodcastEnabled {
		body["podcast_enabled"] = true
	}

	if params.PodcastHasStudentPosts {
		body["podcast_has_student_posts"] = true
	}

	if params.RequireInitialPost {
		body["require_initial_post"] = true
	}

	if params.IsAnnouncement {
		body["is_announcement"] = true
	}

	if params.Pinned {
		body["pinned"] = true
	}

	if params.PositionAfter != "" {
		body["position_after"] = params.PositionAfter
	}

	if params.GroupCategoryID > 0 {
		body["group_category_id"] = params.GroupCategoryID
	}

	if params.OnlyGradersCanRate {
		body["only_graders_can_rate"] = true
	}

	if params.SpecificSections != "" {
		body["specific_sections"] = params.SpecificSections
	}

	var topic DiscussionTopic
	if err := s.client.PostJSON(ctx, path, body, &topic); err != nil {
		return nil, err
	}

	return &topic, nil
}

// UpdateDiscussionParams holds parameters for updating a discussion
type UpdateDiscussionParams struct {
	Title              *string
	Message            *string
	DiscussionType     *string
	Published          *bool
	DelayedPostAt      *string
	AllowRating        *bool
	LockAt             *string
	PodcastEnabled     *bool
	RequireInitialPost *bool
	Pinned             *bool
	Locked             *bool
}

// Update updates an existing discussion topic
func (s *DiscussionsService) Update(ctx context.Context, courseID, topicID int64, params *UpdateDiscussionParams) (*DiscussionTopic, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/discussion_topics/%d", courseID, topicID)

	body := make(map[string]interface{})

	if params.Title != nil {
		body["title"] = *params.Title
	}

	if params.Message != nil {
		body["message"] = *params.Message
	}

	if params.DiscussionType != nil {
		body["discussion_type"] = *params.DiscussionType
	}

	if params.Published != nil {
		body["published"] = *params.Published
	}

	if params.DelayedPostAt != nil {
		body["delayed_post_at"] = *params.DelayedPostAt
	}

	if params.AllowRating != nil {
		body["allow_rating"] = *params.AllowRating
	}

	if params.LockAt != nil {
		body["lock_at"] = *params.LockAt
	}

	if params.PodcastEnabled != nil {
		body["podcast_enabled"] = *params.PodcastEnabled
	}

	if params.RequireInitialPost != nil {
		body["require_initial_post"] = *params.RequireInitialPost
	}

	if params.Pinned != nil {
		body["pinned"] = *params.Pinned
	}

	if params.Locked != nil {
		body["locked"] = *params.Locked
	}

	var topic DiscussionTopic
	if err := s.client.PutJSON(ctx, path, body, &topic); err != nil {
		return nil, err
	}

	return &topic, nil
}

// Delete deletes a discussion topic
func (s *DiscussionsService) Delete(ctx context.Context, courseID, topicID int64) error {
	path := fmt.Sprintf("/api/v1/courses/%d/discussion_topics/%d", courseID, topicID)
	_, err := s.client.Delete(ctx, path)
	return err
}

// ListEntries retrieves all entries for a discussion topic
func (s *DiscussionsService) ListEntries(ctx context.Context, courseID, topicID int64) ([]DiscussionEntry, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/discussion_topics/%d/entries", courseID, topicID)

	var entries []DiscussionEntry
	if err := s.client.GetAllPages(ctx, path, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

// PostEntry posts a new entry to a discussion topic
func (s *DiscussionsService) PostEntry(ctx context.Context, courseID, topicID int64, message string) (*DiscussionEntry, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/discussion_topics/%d/entries", courseID, topicID)

	body := map[string]interface{}{
		"message": message,
	}

	var entry DiscussionEntry
	if err := s.client.PostJSON(ctx, path, body, &entry); err != nil {
		return nil, err
	}

	return &entry, nil
}

// PostReply posts a reply to an entry
func (s *DiscussionsService) PostReply(ctx context.Context, courseID, topicID, entryID int64, message string) (*DiscussionEntry, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/discussion_topics/%d/entries/%d/replies", courseID, topicID, entryID)

	body := map[string]interface{}{
		"message": message,
	}

	var entry DiscussionEntry
	if err := s.client.PostJSON(ctx, path, body, &entry); err != nil {
		return nil, err
	}

	return &entry, nil
}

// MarkTopicRead marks a topic as read
func (s *DiscussionsService) MarkTopicRead(ctx context.Context, courseID, topicID int64) error {
	path := fmt.Sprintf("/api/v1/courses/%d/discussion_topics/%d/read", courseID, topicID)
	return s.client.PutJSON(ctx, path, nil, nil)
}

// MarkTopicUnread marks a topic as unread
func (s *DiscussionsService) MarkTopicUnread(ctx context.Context, courseID, topicID int64) error {
	path := fmt.Sprintf("/api/v1/courses/%d/discussion_topics/%d/read", courseID, topicID)
	_, err := s.client.Delete(ctx, path)
	return err
}

// Subscribe subscribes to a topic
func (s *DiscussionsService) Subscribe(ctx context.Context, courseID, topicID int64) error {
	path := fmt.Sprintf("/api/v1/courses/%d/discussion_topics/%d/subscribed", courseID, topicID)
	return s.client.PutJSON(ctx, path, nil, nil)
}

// Unsubscribe unsubscribes from a topic
func (s *DiscussionsService) Unsubscribe(ctx context.Context, courseID, topicID int64) error {
	path := fmt.Sprintf("/api/v1/courses/%d/discussion_topics/%d/subscribed", courseID, topicID)
	_, err := s.client.Delete(ctx, path)
	return err
}

// AnnouncementsService handles announcement-specific API calls
type AnnouncementsService struct {
	client *Client
}

// NewAnnouncementsService creates a new announcements service
func NewAnnouncementsService(client *Client) *AnnouncementsService {
	return &AnnouncementsService{client: client}
}

// ListAnnouncementsOptions holds options for listing announcements
type ListAnnouncementsOptions struct {
	ContextCodes   []string // course_123, course_456
	StartDate      string   // yyyy-mm-dd or ISO 8601
	EndDate        string
	ActiveOnly     bool
	LatestOnly     bool
	Include        []string // sections, sections_user_count
}

// List retrieves announcements for the given contexts
func (s *AnnouncementsService) List(ctx context.Context, opts *ListAnnouncementsOptions) ([]DiscussionTopic, error) {
	path := "/api/v1/announcements"

	if opts != nil {
		query := url.Values{}

		for _, code := range opts.ContextCodes {
			query.Add("context_codes[]", code)
		}

		if opts.StartDate != "" {
			query.Add("start_date", opts.StartDate)
		}

		if opts.EndDate != "" {
			query.Add("end_date", opts.EndDate)
		}

		if opts.ActiveOnly {
			query.Add("active_only", "true")
		}

		if opts.LatestOnly {
			query.Add("latest_only", "true")
		}

		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var announcements []DiscussionTopic
	if err := s.client.GetAllPages(ctx, path, &announcements); err != nil {
		return nil, err
	}

	return announcements, nil
}
