package options

import "fmt"

// FilesListOptions contains options for listing files
type FilesListOptions struct {
	CourseID     int64
	FolderID     int64
	UserID       int64
	ContentTypes []string
	SearchTerm   string
	Include      []string
	Sort         string
	Order        string
}

// Validate validates the options
func (o *FilesListOptions) Validate() error {
	// Must specify exactly one context
	contextsSpecified := 0
	if o.CourseID > 0 {
		contextsSpecified++
	}
	if o.FolderID > 0 {
		contextsSpecified++
	}
	if o.UserID > 0 {
		contextsSpecified++
	}

	if contextsSpecified == 0 {
		return fmt.Errorf("must specify one of course-id, folder-id, or user-id")
	}
	if contextsSpecified > 1 {
		return fmt.Errorf("can only specify one of course-id, folder-id, or user-id")
	}

	return nil
}

// FilesGetOptions contains options for getting a file
type FilesGetOptions struct {
	FileID  int64
	Include []string
}

// Validate validates the options
func (o *FilesGetOptions) Validate() error {
	return ValidateRequired("file-id", o.FileID)
}

// FilesUploadOptions contains options for uploading a file
type FilesUploadOptions struct {
	FilePath       string
	CourseID       int64
	FolderID       int64
	UserID         int64
	OnDuplicate    string
	ParentFolderID int64
	Hidden         bool
	Locked         bool
}

// Validate validates the options
func (o *FilesUploadOptions) Validate() error {
	if o.FilePath == "" {
		return fmt.Errorf("file-path is required")
	}

	// Must specify exactly one context
	contextsSpecified := 0
	if o.CourseID > 0 {
		contextsSpecified++
	}
	if o.FolderID > 0 {
		contextsSpecified++
	}
	if o.UserID > 0 {
		contextsSpecified++
	}

	if contextsSpecified == 0 {
		return fmt.Errorf("must specify one of course-id, folder-id, or user-id")
	}
	if contextsSpecified > 1 {
		return fmt.Errorf("can only specify one of course-id, folder-id, or user-id")
	}

	return nil
}

// FilesDownloadOptions contains options for downloading a file
type FilesDownloadOptions struct {
	FileID      int64
	Destination string
}

// Validate validates the options
func (o *FilesDownloadOptions) Validate() error {
	return ValidateRequired("file-id", o.FileID)
}

// FilesDeleteOptions contains options for deleting a file
type FilesDeleteOptions struct {
	FileID int64
	Force  bool
}

// Validate validates the options
func (o *FilesDeleteOptions) Validate() error {
	return ValidateRequired("file-id", o.FileID)
}

// FilesQuotaOptions contains options for getting quota information
type FilesQuotaOptions struct {
	CourseID int64
	UserID   int64
}

// Validate validates the options
func (o *FilesQuotaOptions) Validate() error {
	// Must specify exactly one context
	if o.CourseID == 0 && o.UserID == 0 {
		return fmt.Errorf("must specify either course-id or user-id")
	}
	if o.CourseID > 0 && o.UserID > 0 {
		return fmt.Errorf("can only specify one of course-id or user-id")
	}

	return nil
}
