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

// FilesService handles file-related API calls
type FilesService struct {
	client *Client
}

// isCanvasDomain checks if the given URL belongs to the same Canvas instance
// This is used to prevent leaking the Authorization header to third-party storage providers
func isCanvasDomain(redirectURL, baseURL string) bool {
	redirectParsed, err := url.Parse(redirectURL)
	if err != nil {
		return false
	}
	baseParsed, err := url.Parse(baseURL)
	if err != nil {
		return false
	}
	return redirectParsed.Host == baseParsed.Host
}

// NewFilesService creates a new files service
func NewFilesService(client *Client) *FilesService {
	return &FilesService{client: client}
}

// ListFilesOptions holds options for listing files
type ListFilesOptions struct {
	ContentTypes []string // Filter by MIME type
	SearchTerm   string   // Search by file name
	Include      []string // Additional data to include (user)
	Only         []string // Filter by type (names, folders)
	Sort         string   // Sort by (name, size, created_at, updated_at, content_type)
	Order        string   // Order direction (asc, desc)
	Page         int
	PerPage      int
}

// ListCourseFiles retrieves files for a course
func (s *FilesService) ListCourseFiles(ctx context.Context, courseID int64, opts *ListFilesOptions) ([]Attachment, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/files", courseID)
	return s.listFiles(ctx, path, opts)
}

// ListFolderFiles retrieves files in a folder
func (s *FilesService) ListFolderFiles(ctx context.Context, folderID int64, opts *ListFilesOptions) ([]Attachment, error) {
	path := fmt.Sprintf("/api/v1/folders/%d/files", folderID)
	return s.listFiles(ctx, path, opts)
}

// ListUserFiles retrieves files for a user
func (s *FilesService) ListUserFiles(ctx context.Context, userID int64, opts *ListFilesOptions) ([]Attachment, error) {
	path := fmt.Sprintf("/api/v1/users/%d/files", userID)
	return s.listFiles(ctx, path, opts)
}

// listFiles is a helper for listing files with options
func (s *FilesService) listFiles(ctx context.Context, basePath string, opts *ListFilesOptions) ([]Attachment, error) {
	path := basePath

	if opts != nil {
		query := url.Values{}

		for _, ct := range opts.ContentTypes {
			query.Add("content_types[]", ct)
		}

		if opts.SearchTerm != "" {
			query.Add("search_term", opts.SearchTerm)
		}

		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}

		for _, only := range opts.Only {
			query.Add("only[]", only)
		}

		if opts.Sort != "" {
			query.Add("sort", opts.Sort)
		}

		if opts.Order != "" {
			query.Add("order", opts.Order)
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

	var files []Attachment
	if err := s.client.GetAllPages(ctx, path, &files); err != nil {
		return nil, err
	}

	return files, nil
}

// Get retrieves a single file by ID
func (s *FilesService) Get(ctx context.Context, fileID int64, include []string) (*Attachment, error) {
	path := fmt.Sprintf("/api/v1/files/%d", fileID)

	if len(include) > 0 {
		query := url.Values{}
		for _, inc := range include {
			query.Add("include[]", inc)
		}
		path += "?" + query.Encode()
	}

	var file Attachment
	if err := s.client.GetJSON(ctx, path, &file); err != nil {
		return nil, err
	}

	return &file, nil
}

// UploadParams holds parameters for uploading a file
type UploadParams struct {
	Name           string // File name
	Size           int64  // File size in bytes (required for Canvas)
	ContentType    string // MIME type
	ParentFolderID int64  // Folder to upload to
	OnDuplicate    string // How to handle duplicates: overwrite, rename
	LockAt         string // ISO8601 date
	UnlockAt       string // ISO8601 date
	Locked         bool   // Lock the file
	Hidden         bool   // Hide from students
}

// UploadToCourse uploads a file to a course
func (s *FilesService) UploadToCourse(ctx context.Context, courseID int64, filePath string, params *UploadParams) (*Attachment, error) {
	uploadPath := fmt.Sprintf("/api/v1/courses/%d/files", courseID)
	return s.upload(ctx, uploadPath, filePath, params)
}

// UploadToFolder uploads a file to a specific folder
func (s *FilesService) UploadToFolder(ctx context.Context, folderID int64, filePath string, params *UploadParams) (*Attachment, error) {
	uploadPath := fmt.Sprintf("/api/v1/folders/%d/files", folderID)
	return s.upload(ctx, uploadPath, filePath, params)
}

// UploadToUser uploads a file to a user's files
func (s *FilesService) UploadToUser(ctx context.Context, userID int64, filePath string, params *UploadParams) (*Attachment, error) {
	uploadPath := fmt.Sprintf("/api/v1/users/%d/files", userID)
	return s.upload(ctx, uploadPath, filePath, params)
}

// upload is a helper that handles the Canvas three-step upload process
// Step 1: Request upload parameters from Canvas
// Step 2: Upload file to the provided URL using multipart/form-data with upload_params
// Step 3: Confirm upload (follow redirect or parse response)
func (s *FilesService) upload(ctx context.Context, uploadPath, filePath string, params *UploadParams) (*Attachment, error) {
	// Open the file to upload
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	// Use provided name or default to filename
	fileName := params.Name
	if fileName == "" {
		fileName = filepath.Base(filePath)
	}

	// Build upload request parameters
	uploadBody := map[string]interface{}{
		"name": fileName,
		"size": fileInfo.Size(),
	}

	if params.ContentType != "" {
		uploadBody["content_type"] = params.ContentType
	}
	if params.ParentFolderID > 0 {
		uploadBody["parent_folder_id"] = params.ParentFolderID
	}
	if params.OnDuplicate != "" {
		uploadBody["on_duplicate"] = params.OnDuplicate
	}
	if params.LockAt != "" {
		uploadBody["lock_at"] = params.LockAt
	}
	if params.UnlockAt != "" {
		uploadBody["unlock_at"] = params.UnlockAt
	}
	if params.Locked {
		uploadBody["locked"] = true
	}
	if params.Hidden {
		uploadBody["hidden"] = true
	}

	// Step 1: Tell Canvas we want to upload
	var uploadResponse struct {
		UploadURL    string                 `json:"upload_url"`
		UploadParams map[string]interface{} `json:"upload_params"`
		FileParam    string                 `json:"file_param"`
	}

	if err := s.client.PostJSON(ctx, uploadPath, uploadBody, &uploadResponse); err != nil {
		return nil, fmt.Errorf("failed to initialize upload: %w", err)
	}

	// Step 2: Upload the actual file to the provided URL using multipart/form-data
	// Create a buffer to write our multipart form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add all upload params from Canvas first (order matters for some storage providers)
	for key, value := range uploadResponse.UploadParams {
		var strValue string
		switch v := value.(type) {
		case string:
			strValue = v
		case float64:
			strValue = strconv.FormatFloat(v, 'f', -1, 64)
		case bool:
			strValue = strconv.FormatBool(v)
		default:
			strValue = fmt.Sprintf("%v", v)
		}
		if err := writer.WriteField(key, strValue); err != nil {
			return nil, fmt.Errorf("failed to write form field %s: %w", key, err)
		}
	}

	// Determine the file field name (usually "file" but Canvas tells us)
	fileFieldName := "file"
	if uploadResponse.FileParam != "" {
		fileFieldName = uploadResponse.FileParam
	}

	// Add the file itself
	part, err := writer.CreateFormFile(fileFieldName, fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	// Reset file position and copy content
	if _, err := file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to seek file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	// Close the writer to finalize the multipart message
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Create the upload request
	req, err := http.NewRequestWithContext(ctx, "POST", uploadResponse.UploadURL, &requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create upload request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create a custom client that doesn't follow redirects automatically
	// We need to handle the redirect ourselves to get the file confirmation
	noRedirectClient := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Execute the upload
	resp, err := noRedirectClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}
	defer resp.Body.Close()

	// Handle different response scenarios
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		// Direct response with file info
		var uploadedFile Attachment
		if err := json.NewDecoder(resp.Body).Decode(&uploadedFile); err != nil {
			return nil, fmt.Errorf("failed to parse upload response: %w", err)
		}
		return &uploadedFile, nil

	case http.StatusFound, http.StatusSeeOther, http.StatusMovedPermanently, http.StatusTemporaryRedirect:
		// Step 3: Follow redirect to confirm upload
		location := resp.Header.Get("Location")
		if location == "" {
			return nil, fmt.Errorf("redirect response missing Location header")
		}

		// Make a GET request to the redirect URL to confirm the upload
		confirmReq, err := http.NewRequestWithContext(ctx, "GET", location, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create confirmation request: %w", err)
		}

		// Add authorization header only if the redirect is to Canvas domain
		// This prevents leaking the bearer token to third-party storage providers
		if isCanvasDomain(location, s.client.baseURL) {
			confirmReq.Header.Set("Authorization", "Bearer "+s.client.token)
		}

		confirmResp, err := s.client.httpClient.Do(confirmReq)
		if err != nil {
			return nil, fmt.Errorf("failed to confirm upload: %w", err)
		}
		defer confirmResp.Body.Close()

		if confirmResp.StatusCode != http.StatusOK && confirmResp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(confirmResp.Body)
			return nil, fmt.Errorf("upload confirmation failed with status %d: %s", confirmResp.StatusCode, string(bodyBytes))
		}

		var uploadedFile Attachment
		if err := json.NewDecoder(confirmResp.Body).Decode(&uploadedFile); err != nil {
			return nil, fmt.Errorf("failed to parse confirmation response: %w", err)
		}
		return &uploadedFile, nil

	default:
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}
}

// Download downloads a file to the specified destination
func (s *FilesService) Download(ctx context.Context, fileID int64, destPath string) error {
	// Get file info first to get the download URL
	file, err := s.Get(ctx, fileID, nil)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	if file.URL == "" {
		return fmt.Errorf("file has no download URL")
	}

	// Create the destination file
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer out.Close()

	// Download the file content
	req, err := http.NewRequestWithContext(ctx, "GET", file.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to create download request: %w", err)
	}

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Copy the content to the destination file
	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to write file content: %w", err)
	}

	return nil
}

// UpdateParams holds parameters for updating a file
type UpdateParams struct {
	Name           string  // New file name
	ParentFolderID *int64  // Move to different folder
	LockAt         *string // ISO8601 date
	UnlockAt       *string // ISO8601 date
	Locked         *bool
	Hidden         *bool
}

// Update updates file metadata
func (s *FilesService) Update(ctx context.Context, fileID int64, params *UpdateParams) (*Attachment, error) {
	path := fmt.Sprintf("/api/v1/files/%d", fileID)

	body := make(map[string]interface{})

	if params.Name != "" {
		body["name"] = params.Name
	}
	if params.ParentFolderID != nil {
		body["parent_folder_id"] = *params.ParentFolderID
	}
	if params.LockAt != nil {
		body["lock_at"] = *params.LockAt
	}
	if params.UnlockAt != nil {
		body["unlock_at"] = *params.UnlockAt
	}
	if params.Locked != nil {
		body["locked"] = *params.Locked
	}
	if params.Hidden != nil {
		body["hidden"] = *params.Hidden
	}

	var updatedFile Attachment
	if err := s.client.PutJSON(ctx, path, body, &updatedFile); err != nil {
		return nil, err
	}

	return &updatedFile, nil
}

// Delete deletes a file
func (s *FilesService) Delete(ctx context.Context, fileID int64) error {
	path := fmt.Sprintf("/api/v1/files/%d", fileID)
	_, err := s.client.Delete(ctx, path)
	return err
}

// GetQuota retrieves quota information for a course or user
func (s *FilesService) GetCourseQuota(ctx context.Context, courseID int64) (*QuotaInfo, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/files/quota", courseID)
	var quota QuotaInfo
	if err := s.client.GetJSON(ctx, path, &quota); err != nil {
		return nil, err
	}
	return &quota, nil
}

// GetUserQuota retrieves quota information for a user
func (s *FilesService) GetUserQuota(ctx context.Context, userID int64) (*QuotaInfo, error) {
	path := fmt.Sprintf("/api/v1/users/%d/files/quota", userID)
	var quota QuotaInfo
	if err := s.client.GetJSON(ctx, path, &quota); err != nil {
		return nil, err
	}
	return &quota, nil
}

// QuotaInfo represents storage quota information
type QuotaInfo struct {
	QuotaUsed int64 `json:"quota_used"`
	Quota     int64 `json:"quota"`
}
