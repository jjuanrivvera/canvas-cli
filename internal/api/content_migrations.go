package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

// ContentMigrationsService handles content migration-related API calls
type ContentMigrationsService struct {
	client *Client
}

// NewContentMigrationsService creates a new content migrations service
func NewContentMigrationsService(client *Client) *ContentMigrationsService {
	return &ContentMigrationsService{client: client}
}

// ContentMigration represents a Canvas content migration
type ContentMigration struct {
	ID                   int64                `json:"id"`
	MigrationType        string               `json:"migration_type,omitempty"`
	MigrationTypeTitle   string               `json:"migration_type_title,omitempty"`
	MigrationIssuesURL   string               `json:"migration_issues_url,omitempty"`
	MigrationIssuesCount int                  `json:"migration_issues_count,omitempty"`
	Attachment           *MigrationAttachment `json:"attachment,omitempty"`
	ProgressURL          string               `json:"progress_url,omitempty"`
	UserID               int64                `json:"user_id,omitempty"`
	WorkflowState        string               `json:"workflow_state,omitempty"`
	StartedAt            string               `json:"started_at,omitempty"`
	FinishedAt           string               `json:"finished_at,omitempty"`
	CreatedAt            string               `json:"created_at,omitempty"`
	PreAttachment        *PreAttachment       `json:"pre_attachment,omitempty"`
	Settings             *MigrationSettings   `json:"settings,omitempty"`
}

// MigrationAttachment represents an attached file for migration
type MigrationAttachment struct {
	ID          int64  `json:"id"`
	DisplayName string `json:"display_name,omitempty"`
	Filename    string `json:"filename,omitempty"`
	ContentType string `json:"content-type,omitempty"`
	URL         string `json:"url,omitempty"`
	Size        int64  `json:"size,omitempty"`
}

// PreAttachment represents a pre-attachment for file upload
type PreAttachment struct {
	UploadURL    string            `json:"upload_url,omitempty"`
	UploadParams map[string]string `json:"upload_params,omitempty"`
}

// MigrationSettings represents migration settings
type MigrationSettings struct {
	SourceCourseID      int64                  `json:"source_course_id,omitempty"`
	FileURL             string                 `json:"file_url,omitempty"`
	ContentExportID     int64                  `json:"content_export_id,omitempty"`
	QuestionBankID      int64                  `json:"question_bank_id,omitempty"`
	QuestionBankName    string                 `json:"question_bank_name,omitempty"`
	FolderID            int64                  `json:"folder_id,omitempty"`
	OverwriteQuizzes    bool                   `json:"overwrite_quizzes,omitempty"`
	QuestionBankMapping map[string]interface{} `json:"question_bank_mapping,omitempty"`
}

// Migrator represents a content migrator type
type Migrator struct {
	Type               string   `json:"type"`
	RequiresFileUpload bool     `json:"requires_file_upload"`
	Name               string   `json:"name"`
	RequiredSettings   []string `json:"required_settings,omitempty"`
}

// MigrationIssue represents an issue encountered during migration
type MigrationIssue struct {
	ID                 int64  `json:"id"`
	ContentMigrationID int64  `json:"content_migration_id"`
	Description        string `json:"description"`
	WorkflowState      string `json:"workflow_state"`
	FixIssueHTMLURL    string `json:"fix_issue_html_url,omitempty"`
	IssueType          string `json:"issue_type"`
	ErrorReportHTMLURL string `json:"error_report_html_url,omitempty"`
	ErrorMessage       string `json:"error_message,omitempty"`
	CreatedAt          string `json:"created_at,omitempty"`
	UpdatedAt          string `json:"updated_at,omitempty"`
}

// ContentListItem represents an item in a migration content list
type ContentListItem struct {
	Type     string            `json:"type"`
	Property string            `json:"property"`
	Title    string            `json:"title"`
	SubItems []ContentListItem `json:"sub_items,omitempty"`
}

// ListContentMigrationsOptions holds options for listing content migrations
type ListContentMigrationsOptions struct {
	Page    int
	PerPage int
}

// List retrieves content migrations for a course
func (s *ContentMigrationsService) List(ctx context.Context, courseID int64, opts *ListContentMigrationsOptions) ([]ContentMigration, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/content_migrations", courseID)

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

	var migrations []ContentMigration
	if err := s.client.GetAllPages(ctx, path, &migrations); err != nil {
		return nil, err
	}

	return migrations, nil
}

// Get retrieves a single content migration
func (s *ContentMigrationsService) Get(ctx context.Context, courseID, migrationID int64) (*ContentMigration, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/content_migrations/%d", courseID, migrationID)

	var migration ContentMigration
	if err := s.client.GetJSON(ctx, path, &migration); err != nil {
		return nil, err
	}

	return &migration, nil
}

// CreateContentMigrationParams holds parameters for creating a content migration
type CreateContentMigrationParams struct {
	MigrationType    string
	SourceCourseID   *int64
	FilePath         string // Local file path for upload
	FileURL          string // Remote file URL
	ContentExportID  *int64
	QuestionBankID   *int64
	QuestionBankName string
	FolderID         *int64
	OverwriteQuizzes *bool
	SelectiveImport  *bool
	CopyOptions      map[string]interface{}
	DateShiftOptions *DateShiftOptions
}

// DateShiftOptions represents date shift options for migration
type DateShiftOptions struct {
	ShiftDates       bool           `json:"shift_dates,omitempty"`
	OldStartDate     string         `json:"old_start_date,omitempty"`
	OldEndDate       string         `json:"old_end_date,omitempty"`
	NewStartDate     string         `json:"new_start_date,omitempty"`
	NewEndDate       string         `json:"new_end_date,omitempty"`
	DaySubstitutions map[string]int `json:"day_substitutions,omitempty"`
	RemoveDates      bool           `json:"remove_dates,omitempty"`
}

// Create creates a new content migration
func (s *ContentMigrationsService) Create(ctx context.Context, courseID int64, params *CreateContentMigrationParams) (*ContentMigration, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/content_migrations", courseID)

	body := make(map[string]interface{})
	body["migration_type"] = params.MigrationType

	settings := make(map[string]interface{})

	if params.SourceCourseID != nil {
		settings["source_course_id"] = *params.SourceCourseID
	}

	if params.FileURL != "" {
		settings["file_url"] = params.FileURL
	}

	if params.ContentExportID != nil {
		settings["content_export_id"] = *params.ContentExportID
	}

	if params.QuestionBankID != nil {
		settings["question_bank_id"] = *params.QuestionBankID
	}

	if params.QuestionBankName != "" {
		settings["question_bank_name"] = params.QuestionBankName
	}

	if params.FolderID != nil {
		settings["folder_id"] = *params.FolderID
	}

	if params.OverwriteQuizzes != nil {
		settings["overwrite_quizzes"] = *params.OverwriteQuizzes
	}

	if len(settings) > 0 {
		body["settings"] = settings
	}

	if params.SelectiveImport != nil {
		body["selective_import"] = *params.SelectiveImport
	}

	if len(params.CopyOptions) > 0 {
		body["copy"] = params.CopyOptions
	}

	if params.DateShiftOptions != nil {
		body["date_shift_options"] = params.DateShiftOptions
	}

	// If there's a file to upload, we need to handle it differently
	if params.FilePath != "" {
		return s.createWithFile(ctx, courseID, params, body)
	}

	var migration ContentMigration
	if err := s.client.PostJSON(ctx, path, body, &migration); err != nil {
		return nil, err
	}

	return &migration, nil
}

// createWithFile creates a migration with a file upload
func (s *ContentMigrationsService) createWithFile(ctx context.Context, courseID int64, params *CreateContentMigrationParams, bodyData map[string]interface{}) (*ContentMigration, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/content_migrations", courseID)

	file, err := os.Open(params.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Write form fields
	if err := writer.WriteField("migration_type", params.MigrationType); err != nil {
		return nil, fmt.Errorf("failed to write migration_type: %w", err)
	}

	if params.SourceCourseID != nil {
		if err := writer.WriteField("settings[source_course_id]", strconv.FormatInt(*params.SourceCourseID, 10)); err != nil {
			return nil, fmt.Errorf("failed to write source_course_id: %w", err)
		}
	}

	if params.FolderID != nil {
		if err := writer.WriteField("settings[folder_id]", strconv.FormatInt(*params.FolderID, 10)); err != nil {
			return nil, fmt.Errorf("failed to write folder_id: %w", err)
		}
	}

	if params.SelectiveImport != nil {
		if err := writer.WriteField("selective_import", strconv.FormatBool(*params.SelectiveImport)); err != nil {
			return nil, fmt.Errorf("failed to write selective_import: %w", err)
		}
	}

	// Write file
	part, err := writer.CreateFormFile("pre_attachment[name]", filepath.Base(params.FilePath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	fullURL := s.client.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+s.client.token)

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s", string(bodyBytes))
	}

	var migration ContentMigration
	if err := json.NewDecoder(resp.Body).Decode(&migration); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &migration, nil
}

// UpdateContentMigrationParams holds parameters for updating a content migration
type UpdateContentMigrationParams struct {
	WorkflowState string
}

// Update updates a content migration
func (s *ContentMigrationsService) Update(ctx context.Context, courseID, migrationID int64, params *UpdateContentMigrationParams) (*ContentMigration, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/content_migrations/%d", courseID, migrationID)

	body := make(map[string]interface{})
	if params.WorkflowState != "" {
		body["workflow_state"] = params.WorkflowState
	}

	var migration ContentMigration
	if err := s.client.PutJSON(ctx, path, body, &migration); err != nil {
		return nil, err
	}

	return &migration, nil
}

// ListContentList retrieves available content for selective import
func (s *ContentMigrationsService) ListContentList(ctx context.Context, courseID, migrationID int64, contentType string) ([]ContentListItem, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/content_migrations/%d/content_list", courseID, migrationID)

	if contentType != "" {
		path += "?type=" + url.QueryEscape(contentType)
	}

	var items []ContentListItem
	if err := s.client.GetAllPages(ctx, path, &items); err != nil {
		return nil, err
	}

	return items, nil
}

// ListMigrators retrieves available migrator types for a course
func (s *ContentMigrationsService) ListMigrators(ctx context.Context, courseID int64) ([]Migrator, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/content_migrations/migrators", courseID)

	var migrators []Migrator
	if err := s.client.GetJSON(ctx, path, &migrators); err != nil {
		return nil, err
	}

	return migrators, nil
}

// ListMigrationIssues retrieves issues for a content migration
func (s *ContentMigrationsService) ListMigrationIssues(ctx context.Context, courseID, migrationID int64) ([]MigrationIssue, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/content_migrations/%d/migration_issues", courseID, migrationID)

	var issues []MigrationIssue
	if err := s.client.GetAllPages(ctx, path, &issues); err != nil {
		return nil, err
	}

	return issues, nil
}
