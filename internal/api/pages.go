package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// Page represents a Canvas wiki page
type Page struct {
	PageID                int64                  `json:"page_id"`
	URL                   string                 `json:"url"`
	Title                 string                 `json:"title"`
	CreatedAt             time.Time              `json:"created_at"`
	UpdatedAt             time.Time              `json:"updated_at"`
	HideFromStudents      bool                   `json:"hide_from_students"`
	EditingRoles          string                 `json:"editing_roles"`
	LastEditedBy          *User                  `json:"last_edited_by,omitempty"`
	Body                  string                 `json:"body,omitempty"`
	Published             bool                   `json:"published"`
	PublishAt             *time.Time             `json:"publish_at,omitempty"`
	FrontPage             bool                   `json:"front_page"`
	LockedForUser         bool                   `json:"locked_for_user"`
	LockInfo              *LockInfo              `json:"lock_info,omitempty"`
	LockExplanation       string                 `json:"lock_explanation,omitempty"`
	Editor                string                 `json:"editor,omitempty"`
	BlockEditorAttributes map[string]interface{} `json:"block_editor_attributes,omitempty"`
}

// PageRevision represents a revision of a wiki page
type PageRevision struct {
	RevisionID int64     `json:"revision_id"`
	UpdatedAt  time.Time `json:"updated_at"`
	Latest     bool      `json:"latest"`
	EditedBy   *User     `json:"edited_by,omitempty"`
	URL        string    `json:"url,omitempty"`
	Title      string    `json:"title,omitempty"`
	Body       string    `json:"body,omitempty"`
}

// PagesService handles page-related API calls
type PagesService struct {
	client *Client
}

// NewPagesService creates a new pages service
func NewPagesService(client *Client) *PagesService {
	return &PagesService{client: client}
}

// ListPagesOptions holds options for listing pages
type ListPagesOptions struct {
	Sort       string // title, created_at, updated_at
	Order      string // asc, desc
	SearchTerm string
	Published  *bool
	Include    []string // body
	Page       int
	PerPage    int
}

// List retrieves all pages for a course
func (s *PagesService) List(ctx context.Context, courseID int64, opts *ListPagesOptions) ([]Page, error) {
	if err := ValidatePositiveID(courseID, "course_id"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/api/v1/courses/%d/pages", courseID)

	if opts != nil {
		query := url.Values{}

		if opts.Sort != "" {
			query.Add("sort", opts.Sort)
		}

		if opts.Order != "" {
			query.Add("order", opts.Order)
		}

		if opts.SearchTerm != "" {
			query.Add("search_term", opts.SearchTerm)
		}

		if opts.Published != nil {
			query.Add("published", strconv.FormatBool(*opts.Published))
		}

		for _, inc := range opts.Include {
			query.Add("include[]", inc)
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

	var pages []Page
	if err := s.client.GetAllPages(ctx, path, &pages); err != nil {
		return nil, err
	}

	return pages, nil
}

// Get retrieves a single page by URL or ID
func (s *PagesService) Get(ctx context.Context, courseID int64, urlOrID string) (*Page, error) {
	if err := ValidatePositiveID(courseID, "course_id"); err != nil {
		return nil, err
	}
	if err := ValidateNonEmpty(urlOrID, "url_or_id"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/api/v1/courses/%d/pages/%s", courseID, url.PathEscape(urlOrID))

	var page Page
	if err := s.client.GetJSON(ctx, path, &page); err != nil {
		return nil, err
	}

	return &page, nil
}

// GetFrontPage retrieves the front page for a course
func (s *PagesService) GetFrontPage(ctx context.Context, courseID int64) (*Page, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/front_page", courseID)

	var page Page
	if err := s.client.GetJSON(ctx, path, &page); err != nil {
		return nil, err
	}

	return &page, nil
}

// CreatePageParams holds parameters for creating a page
type CreatePageParams struct {
	Title          string
	Body           string
	EditingRoles   string
	NotifyOfUpdate bool
	Published      bool
	FrontPage      bool
	PublishAt      string
}

// Create creates a new page in a course
func (s *PagesService) Create(ctx context.Context, courseID int64, params *CreatePageParams) (*Page, error) {
	if err := ValidatePositiveID(courseID, "course_id"); err != nil {
		return nil, err
	}
	if params == nil {
		return nil, ErrNilParams
	}
	if err := ValidateNonEmpty(params.Title, "title"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/api/v1/courses/%d/pages", courseID)

	body := map[string]interface{}{
		"wiki_page": make(map[string]interface{}),
	}

	pageData := body["wiki_page"].(map[string]interface{})
	pageData["title"] = params.Title

	if params.Body != "" {
		pageData["body"] = params.Body
	}

	if params.EditingRoles != "" {
		pageData["editing_roles"] = params.EditingRoles
	}

	if params.NotifyOfUpdate {
		pageData["notify_of_update"] = true
	}

	pageData["published"] = params.Published

	if params.FrontPage {
		pageData["front_page"] = true
	}

	if params.PublishAt != "" {
		pageData["publish_at"] = params.PublishAt
	}

	var page Page
	if err := s.client.PostJSON(ctx, path, body, &page); err != nil {
		return nil, err
	}

	return &page, nil
}

// UpdatePageParams holds parameters for updating a page
type UpdatePageParams struct {
	Title          *string
	Body           *string
	EditingRoles   *string
	NotifyOfUpdate *bool
	Published      *bool
	FrontPage      *bool
	PublishAt      *string
}

// Update updates an existing page
func (s *PagesService) Update(ctx context.Context, courseID int64, urlOrID string, params *UpdatePageParams) (*Page, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/pages/%s", courseID, url.PathEscape(urlOrID))

	body := map[string]interface{}{
		"wiki_page": make(map[string]interface{}),
	}

	pageData := body["wiki_page"].(map[string]interface{})

	if params.Title != nil {
		pageData["title"] = *params.Title
	}

	if params.Body != nil {
		pageData["body"] = *params.Body
	}

	if params.EditingRoles != nil {
		pageData["editing_roles"] = *params.EditingRoles
	}

	if params.NotifyOfUpdate != nil {
		pageData["notify_of_update"] = *params.NotifyOfUpdate
	}

	if params.Published != nil {
		pageData["published"] = *params.Published
	}

	if params.FrontPage != nil {
		pageData["front_page"] = *params.FrontPage
	}

	if params.PublishAt != nil {
		pageData["publish_at"] = *params.PublishAt
	}

	var page Page
	if err := s.client.PutJSON(ctx, path, body, &page); err != nil {
		return nil, err
	}

	return &page, nil
}

// UpdateFrontPage updates the front page
func (s *PagesService) UpdateFrontPage(ctx context.Context, courseID int64, params *UpdatePageParams) (*Page, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/front_page", courseID)

	body := map[string]interface{}{
		"wiki_page": make(map[string]interface{}),
	}

	pageData := body["wiki_page"].(map[string]interface{})

	if params.Title != nil {
		pageData["title"] = *params.Title
	}

	if params.Body != nil {
		pageData["body"] = *params.Body
	}

	if params.EditingRoles != nil {
		pageData["editing_roles"] = *params.EditingRoles
	}

	if params.NotifyOfUpdate != nil {
		pageData["notify_of_update"] = *params.NotifyOfUpdate
	}

	if params.Published != nil {
		pageData["published"] = *params.Published
	}

	var page Page
	if err := s.client.PutJSON(ctx, path, body, &page); err != nil {
		return nil, err
	}

	return &page, nil
}

// Delete deletes a page
func (s *PagesService) Delete(ctx context.Context, courseID int64, urlOrID string) error {
	path := fmt.Sprintf("/api/v1/courses/%d/pages/%s", courseID, url.PathEscape(urlOrID))
	_, err := s.client.Delete(ctx, path)
	return err
}

// Duplicate duplicates a page
func (s *PagesService) Duplicate(ctx context.Context, courseID int64, urlOrID string) (*Page, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/pages/%s/duplicate", courseID, url.PathEscape(urlOrID))

	var page Page
	if err := s.client.PostJSON(ctx, path, nil, &page); err != nil {
		return nil, err
	}

	return &page, nil
}

// ListRevisions retrieves all revisions for a page
func (s *PagesService) ListRevisions(ctx context.Context, courseID int64, urlOrID string) ([]PageRevision, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/pages/%s/revisions", courseID, url.PathEscape(urlOrID))

	var revisions []PageRevision
	if err := s.client.GetAllPages(ctx, path, &revisions); err != nil {
		return nil, err
	}

	return revisions, nil
}

// GetRevision retrieves a specific revision
func (s *PagesService) GetRevision(ctx context.Context, courseID int64, urlOrID string, revisionID int64, summary bool) (*PageRevision, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/pages/%s/revisions/%d", courseID, url.PathEscape(urlOrID), revisionID)

	if summary {
		path += "?summary=1"
	}

	var revision PageRevision
	if err := s.client.GetJSON(ctx, path, &revision); err != nil {
		return nil, err
	}

	return &revision, nil
}

// GetLatestRevision retrieves the latest revision
func (s *PagesService) GetLatestRevision(ctx context.Context, courseID int64, urlOrID string, summary bool) (*PageRevision, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/pages/%s/revisions/latest", courseID, url.PathEscape(urlOrID))

	if summary {
		path += "?summary=1"
	}

	var revision PageRevision
	if err := s.client.GetJSON(ctx, path, &revision); err != nil {
		return nil, err
	}

	return &revision, nil
}

// RevertToRevision reverts a page to a specific revision
func (s *PagesService) RevertToRevision(ctx context.Context, courseID int64, urlOrID string, revisionID int64) (*PageRevision, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/pages/%s/revisions/%d", courseID, url.PathEscape(urlOrID), revisionID)

	var revision PageRevision
	if err := s.client.PostJSON(ctx, path, nil, &revision); err != nil {
		return nil, err
	}

	return &revision, nil
}
