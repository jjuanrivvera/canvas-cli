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

// SISImportsService handles SIS import-related API calls
type SISImportsService struct {
	client *Client
}

// NewSISImportsService creates a new SIS imports service
func NewSISImportsService(client *Client) *SISImportsService {
	return &SISImportsService{client: client}
}

// SISImport represents a Canvas SIS import
type SISImport struct {
	ID                       int64                `json:"id"`
	CreatedAt                string               `json:"created_at,omitempty"`
	StartedAt                string               `json:"started_at,omitempty"`
	EndedAt                  string               `json:"ended_at,omitempty"`
	UpdatedAt                string               `json:"updated_at,omitempty"`
	Progress                 float64              `json:"progress,omitempty"`
	WorkflowState            string               `json:"workflow_state,omitempty"`
	Data                     *SISImportData       `json:"data,omitempty"`
	Statistics               *SISImportStats      `json:"statistics,omitempty"`
	BatchMode                *bool                `json:"batch_mode,omitempty"`
	BatchModeTermID          *int64               `json:"batch_mode_term_id,omitempty"`
	OverrideSISStickiness    *bool                `json:"override_sis_stickiness,omitempty"`
	AddSISStickiness         *bool                `json:"add_sis_stickiness,omitempty"`
	ClearSISStickiness       *bool                `json:"clear_sis_stickiness,omitempty"`
	DiffingDataSetIdentifier *string              `json:"diffing_data_set_identifier,omitempty"`
	DiffedAgainstImportID    *int64               `json:"diffed_against_import_id,omitempty"`
	CSVAttachments           []SISImportCSV       `json:"csv_attachments,omitempty"`
	ErrorsAttachment         *SISImportAttachment `json:"errors_attachment,omitempty"`
	ProcessingWarnings       [][]string           `json:"processing_warnings,omitempty"`
	ProcessingErrors         [][]string           `json:"processing_errors,omitempty"`
}

// SISImportData represents import data summary
type SISImportData struct {
	ImportType      string           `json:"import_type,omitempty"`
	SuppliedBatches []string         `json:"supplied_batches,omitempty"`
	Counts          *SISImportCounts `json:"counts,omitempty"`
}

// SISImportCounts represents counts of imported items
type SISImportCounts struct {
	Accounts                int `json:"accounts,omitempty"`
	Terms                   int `json:"terms,omitempty"`
	Courses                 int `json:"courses,omitempty"`
	Sections                int `json:"sections,omitempty"`
	Xlists                  int `json:"xlists,omitempty"`
	Users                   int `json:"users,omitempty"`
	Enrollments             int `json:"enrollments,omitempty"`
	Groups                  int `json:"groups,omitempty"`
	GroupMemberships        int `json:"group_memberships,omitempty"`
	GradePublishing         int `json:"grade_publishing_results,omitempty"`
	BatchCoursesDeleted     int `json:"batch_courses_deleted,omitempty"`
	BatchSectionsDeleted    int `json:"batch_sections_deleted,omitempty"`
	BatchEnrollmentsDeleted int `json:"batch_enrollments_deleted,omitempty"`
	ErrorCount              int `json:"error_count,omitempty"`
	WarningCount            int `json:"warning_count,omitempty"`
}

// SISImportStats represents import statistics
type SISImportStats struct {
	TotalStateChanges           int `json:"total_state_changes,omitempty"`
	AccountStateChanges         int `json:"Account,omitempty"`
	CourseStateChanges          int `json:"Course,omitempty"`
	SectionStateChanges         int `json:"CourseSection,omitempty"`
	EnrollmentStateChanges      int `json:"Enrollment,omitempty"`
	UserStateChanges            int `json:"AbstractCourse,omitempty"`
	GroupStateChanges           int `json:"Group,omitempty"`
	GroupMembershipStateChanges int `json:"GroupMembership,omitempty"`
}

// SISImportCSV represents a CSV file in an import
type SISImportCSV struct {
	Filename    string `json:"filename,omitempty"`
	ContentType string `json:"content-type,omitempty"`
	Size        int64  `json:"size,omitempty"`
	URL         string `json:"url,omitempty"`
}

// SISImportAttachment represents an attachment (like errors file)
type SISImportAttachment struct {
	ID          int64  `json:"id,omitempty"`
	UUID        string `json:"uuid,omitempty"`
	FolderID    int64  `json:"folder_id,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Filename    string `json:"filename,omitempty"`
	ContentType string `json:"content-type,omitempty"`
	URL         string `json:"url,omitempty"`
	Size        int64  `json:"size,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

// SISImportError represents an error from a SIS import
type SISImportError struct {
	SISImportID int64  `json:"sis_import_id"`
	File        string `json:"file,omitempty"`
	Message     string `json:"message"`
	Row         int    `json:"row,omitempty"`
	RowInfo     string `json:"row_info,omitempty"`
}

// SISRestoreProgress represents progress of a restore operation
type SISRestoreProgress struct {
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

// ListSISImportsOptions holds options for listing SIS imports
type ListSISImportsOptions struct {
	WorkflowState string
	CreatedSince  string // ISO8601 timestamp
	CreatedBefore string // ISO8601 timestamp
	Page          int
	PerPage       int
}

// List retrieves SIS imports for an account
func (s *SISImportsService) List(ctx context.Context, accountID int64, opts *ListSISImportsOptions) ([]SISImport, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/sis_imports", accountID)

	if opts != nil {
		query := url.Values{}

		if opts.WorkflowState != "" {
			query.Add("workflow_state", opts.WorkflowState)
		}

		if opts.CreatedSince != "" {
			query.Add("created_since", opts.CreatedSince)
		}

		if opts.CreatedBefore != "" {
			query.Add("created_before", opts.CreatedBefore)
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

	var response struct {
		SISImports []SISImport `json:"sis_imports"`
	}
	if err := s.client.GetJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	// Respect global limit if set
	results := response.SISImports
	if maxResults := s.client.GetMaxResults(); maxResults > 0 && len(results) > maxResults {
		results = results[:maxResults]
	}

	return results, nil
}

// Get retrieves a single SIS import
func (s *SISImportsService) Get(ctx context.Context, accountID, importID int64) (*SISImport, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/sis_imports/%d", accountID, importID)

	var sisImport SISImport
	if err := s.client.GetJSON(ctx, path, &sisImport); err != nil {
		return nil, err
	}

	return &sisImport, nil
}

// CreateSISImportParams holds parameters for creating a SIS import
type CreateSISImportParams struct {
	FilePath                 string
	ImportType               string // instructure_csv
	Extension                string // csv, zip
	BatchMode                *bool
	BatchModeTermID          *int64
	OverrideSISStickiness    *bool
	AddSISStickiness         *bool
	ClearSISStickiness       *bool
	DiffingDataSetIdentifier string
	DiffingRemasterDataSet   *bool
	ChangeThreshold          *float64
}

// Create creates a new SIS import
func (s *SISImportsService) Create(ctx context.Context, accountID int64, params *CreateSISImportParams) (*SISImport, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/sis_imports", accountID)

	file, err := os.Open(params.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("attachment", filepath.Base(params.FilePath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	if params.ImportType != "" {
		if err := writer.WriteField("import_type", params.ImportType); err != nil {
			return nil, fmt.Errorf("failed to write import_type: %w", err)
		}
	}

	if params.Extension != "" {
		if err := writer.WriteField("extension", params.Extension); err != nil {
			return nil, fmt.Errorf("failed to write extension: %w", err)
		}
	}

	if params.BatchMode != nil {
		if err := writer.WriteField("batch_mode", strconv.FormatBool(*params.BatchMode)); err != nil {
			return nil, fmt.Errorf("failed to write batch_mode: %w", err)
		}
	}

	if params.BatchModeTermID != nil {
		if err := writer.WriteField("batch_mode_term_id", strconv.FormatInt(*params.BatchModeTermID, 10)); err != nil {
			return nil, fmt.Errorf("failed to write batch_mode_term_id: %w", err)
		}
	}

	if params.OverrideSISStickiness != nil {
		if err := writer.WriteField("override_sis_stickiness", strconv.FormatBool(*params.OverrideSISStickiness)); err != nil {
			return nil, fmt.Errorf("failed to write override_sis_stickiness: %w", err)
		}
	}

	if params.AddSISStickiness != nil {
		if err := writer.WriteField("add_sis_stickiness", strconv.FormatBool(*params.AddSISStickiness)); err != nil {
			return nil, fmt.Errorf("failed to write add_sis_stickiness: %w", err)
		}
	}

	if params.ClearSISStickiness != nil {
		if err := writer.WriteField("clear_sis_stickiness", strconv.FormatBool(*params.ClearSISStickiness)); err != nil {
			return nil, fmt.Errorf("failed to write clear_sis_stickiness: %w", err)
		}
	}

	if params.DiffingDataSetIdentifier != "" {
		if err := writer.WriteField("diffing_data_set_identifier", params.DiffingDataSetIdentifier); err != nil {
			return nil, fmt.Errorf("failed to write diffing_data_set_identifier: %w", err)
		}
	}

	if params.DiffingRemasterDataSet != nil {
		if err := writer.WriteField("diffing_remaster_data_set", strconv.FormatBool(*params.DiffingRemasterDataSet)); err != nil {
			return nil, fmt.Errorf("failed to write diffing_remaster_data_set: %w", err)
		}
	}

	if params.ChangeThreshold != nil {
		if err := writer.WriteField("change_threshold", strconv.FormatFloat(*params.ChangeThreshold, 'f', -1, 64)); err != nil {
			return nil, fmt.Errorf("failed to write change_threshold: %w", err)
		}
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

	var sisImport SISImport
	if err := json.NewDecoder(resp.Body).Decode(&sisImport); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &sisImport, nil
}

// Abort aborts a pending SIS import
func (s *SISImportsService) Abort(ctx context.Context, accountID, importID int64) (*SISImport, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/sis_imports/%d/abort", accountID, importID)

	var sisImport SISImport
	if err := s.client.PutJSON(ctx, path, nil, &sisImport); err != nil {
		return nil, err
	}

	return &sisImport, nil
}

// RestoreStates restores workflow states for an import
func (s *SISImportsService) RestoreStates(ctx context.Context, accountID, importID int64, batchMode bool) (*SISRestoreProgress, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/sis_imports/%d/restore_states", accountID, importID)

	body := map[string]interface{}{}
	if batchMode {
		body["batch_mode"] = true
	}

	var progress SISRestoreProgress
	if err := s.client.PutJSON(ctx, path, body, &progress); err != nil {
		return nil, err
	}

	return &progress, nil
}

// ListErrors retrieves errors for a SIS import
func (s *SISImportsService) ListErrors(ctx context.Context, accountID, importID int64) ([]SISImportError, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/sis_imports/%d/errors", accountID, importID)

	var errors []SISImportError
	if err := s.client.GetAllPages(ctx, path, &errors); err != nil {
		return nil, err
	}

	return errors, nil
}
