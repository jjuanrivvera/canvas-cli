package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// BlueprintService handles blueprint course-related API calls
type BlueprintService struct {
	client *Client
}

// NewBlueprintService creates a new blueprint service
func NewBlueprintService(client *Client) *BlueprintService {
	return &BlueprintService{client: client}
}

// BlueprintTemplate represents a blueprint template
type BlueprintTemplate struct {
	ID                    int64               `json:"id"`
	CourseID              int64               `json:"course_id"`
	LastExportStartedAt   string              `json:"last_export_started_at,omitempty"`
	LastExportCompletedAt string              `json:"last_export_completed_at,omitempty"`
	AssociatedCourseCount int                 `json:"associated_course_count,omitempty"`
	LatestMigration       *BlueprintMigration `json:"latest_migration,omitempty"`
}

// BlueprintMigration represents a blueprint sync/migration
type BlueprintMigration struct {
	ID                 int64  `json:"id"`
	TemplateID         int64  `json:"template_id,omitempty"`
	SubscriptionID     int64  `json:"subscription_id,omitempty"`
	UserID             int64  `json:"user_id,omitempty"`
	User               *User  `json:"user,omitempty"`
	WorkflowState      string `json:"workflow_state,omitempty"`
	CreatedAt          string `json:"created_at,omitempty"`
	ExportsStartedAt   string `json:"exports_started_at,omitempty"`
	ImportsQueuedAt    string `json:"imports_queued_at,omitempty"`
	ImportsCompletedAt string `json:"imports_completed_at,omitempty"`
	Comment            string `json:"comment,omitempty"`
}

// AssociatedCourse represents a course associated with a blueprint
type AssociatedCourse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name,omitempty"`
	CourseCode  string `json:"course_code,omitempty"`
	TermName    string `json:"term_name,omitempty"`
	SISCourseID string `json:"sis_course_id,omitempty"`
	Teachers    []User `json:"teachers,omitempty"`
}

// UnsyncedChange represents a change not yet synced
type UnsyncedChange struct {
	AssetID         int64  `json:"asset_id"`
	AssetType       string `json:"asset_type"`
	AssetName       string `json:"asset_name,omitempty"`
	ChangeType      string `json:"change_type"`
	HTMLUrl         string `json:"html_url,omitempty"`
	Locked          bool   `json:"locked"`
	ExceptionsCount int    `json:"exceptions_count,omitempty"`
}

// BlueprintRestriction represents content lock restrictions
type BlueprintRestriction struct {
	Content           bool `json:"content,omitempty"`
	Points            bool `json:"points,omitempty"`
	DueDates          bool `json:"due_dates,omitempty"`
	AvailabilityDates bool `json:"availability_dates,omitempty"`
}

// GetTemplate retrieves blueprint template details
func (s *BlueprintService) GetTemplate(ctx context.Context, courseID int64, templateID string) (*BlueprintTemplate, error) {
	if templateID == "" {
		templateID = "default"
	}
	path := fmt.Sprintf("/api/v1/courses/%d/blueprint_templates/%s", courseID, templateID)

	var template BlueprintTemplate
	if err := s.client.GetJSON(ctx, path, &template); err != nil {
		return nil, err
	}

	return &template, nil
}

// ListAssociatedCoursesOptions holds options for listing associated courses
type ListAssociatedCoursesOptions struct {
	Page    int
	PerPage int
}

// ListAssociatedCourses retrieves courses associated with a blueprint
func (s *BlueprintService) ListAssociatedCourses(ctx context.Context, courseID int64, templateID string, opts *ListAssociatedCoursesOptions) ([]AssociatedCourse, error) {
	if templateID == "" {
		templateID = "default"
	}
	path := fmt.Sprintf("/api/v1/courses/%d/blueprint_templates/%s/associated_courses", courseID, templateID)

	if opts != nil {
		query := url.Values{}

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

	var courses []AssociatedCourse
	if err := s.client.GetAllPages(ctx, path, &courses); err != nil {
		return nil, err
	}

	return courses, nil
}

// UpdateAssociationsParams holds parameters for updating associations
type UpdateAssociationsParams struct {
	CourseIDsToAdd    []int64
	CourseIDsToRemove []int64
}

// UpdateAssociations updates the courses associated with a blueprint
func (s *BlueprintService) UpdateAssociations(ctx context.Context, courseID int64, templateID string, params *UpdateAssociationsParams) error {
	if templateID == "" {
		templateID = "default"
	}
	path := fmt.Sprintf("/api/v1/courses/%d/blueprint_templates/%s/update_associations", courseID, templateID)

	body := make(map[string]interface{})

	if len(params.CourseIDsToAdd) > 0 {
		body["course_ids_to_add"] = params.CourseIDsToAdd
	}

	if len(params.CourseIDsToRemove) > 0 {
		body["course_ids_to_remove"] = params.CourseIDsToRemove
	}

	var result map[string]interface{}
	if err := s.client.PutJSON(ctx, path, body, &result); err != nil {
		return err
	}

	return nil
}

// SyncParams holds parameters for beginning a sync
type SyncParams struct {
	Comment          string
	SendNotification *bool
	CopySettings     *bool
	PublishAfterSync *bool
}

// BeginSync begins a blueprint sync to associated courses
func (s *BlueprintService) BeginSync(ctx context.Context, courseID int64, templateID string, params *SyncParams) (*BlueprintMigration, error) {
	if templateID == "" {
		templateID = "default"
	}
	path := fmt.Sprintf("/api/v1/courses/%d/blueprint_templates/%s/migrations", courseID, templateID)

	body := make(map[string]interface{})

	if params != nil {
		if params.Comment != "" {
			body["comment"] = params.Comment
		}

		if params.SendNotification != nil {
			body["send_notification"] = *params.SendNotification
		}

		if params.CopySettings != nil {
			body["copy_settings"] = *params.CopySettings
		}

		if params.PublishAfterSync != nil {
			body["publish_after_initial_sync"] = *params.PublishAfterSync
		}
	}

	var migration BlueprintMigration
	if err := s.client.PostJSON(ctx, path, body, &migration); err != nil {
		return nil, err
	}

	return &migration, nil
}

// SetRestrictionParams holds parameters for setting restrictions
type SetRestrictionParams struct {
	ContentType  string // assignment, attachment, discussion_topic, external_tool, lti-quiz, quiz, wiki_page
	ContentID    int64
	Restricted   *bool
	Restrictions *BlueprintRestriction
}

// SetRestriction sets restrictions on a content item
func (s *BlueprintService) SetRestriction(ctx context.Context, courseID int64, templateID string, params *SetRestrictionParams) error {
	if templateID == "" {
		templateID = "default"
	}
	path := fmt.Sprintf("/api/v1/courses/%d/blueprint_templates/%s/restrict_item", courseID, templateID)

	body := make(map[string]interface{})
	body["content_type"] = params.ContentType
	body["content_id"] = params.ContentID

	if params.Restricted != nil {
		body["restricted"] = *params.Restricted
	}

	if params.Restrictions != nil {
		body["restrictions"] = params.Restrictions
	}

	var result map[string]interface{}
	if err := s.client.PutJSON(ctx, path, body, &result); err != nil {
		return err
	}

	return nil
}

// ListUnsyncedChanges retrieves changes not yet synced
func (s *BlueprintService) ListUnsyncedChanges(ctx context.Context, courseID int64, templateID string) ([]UnsyncedChange, error) {
	if templateID == "" {
		templateID = "default"
	}
	path := fmt.Sprintf("/api/v1/courses/%d/blueprint_templates/%s/unsynced_changes", courseID, templateID)

	var changes []UnsyncedChange
	if err := s.client.GetJSON(ctx, path, &changes); err != nil {
		return nil, err
	}

	return changes, nil
}

// ListMigrationsOptions holds options for listing migrations
type ListMigrationsOptions struct {
	Page    int
	PerPage int
}

// ListMigrations retrieves blueprint migrations
func (s *BlueprintService) ListMigrations(ctx context.Context, courseID int64, templateID string, opts *ListMigrationsOptions) ([]BlueprintMigration, error) {
	if templateID == "" {
		templateID = "default"
	}
	path := fmt.Sprintf("/api/v1/courses/%d/blueprint_templates/%s/migrations", courseID, templateID)

	if opts != nil {
		query := url.Values{}

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

	var migrations []BlueprintMigration
	if err := s.client.GetAllPages(ctx, path, &migrations); err != nil {
		return nil, err
	}

	return migrations, nil
}

// GetMigration retrieves a specific migration
func (s *BlueprintService) GetMigration(ctx context.Context, courseID int64, templateID string, migrationID int64, include []string) (*BlueprintMigration, error) {
	if templateID == "" {
		templateID = "default"
	}
	path := fmt.Sprintf("/api/v1/courses/%d/blueprint_templates/%s/migrations/%d", courseID, templateID, migrationID)

	if len(include) > 0 {
		path += "?include[]=" + strings.Join(include, "&include[]=")
	}

	var migration BlueprintMigration
	if err := s.client.GetJSON(ctx, path, &migration); err != nil {
		return nil, err
	}

	return &migration, nil
}
