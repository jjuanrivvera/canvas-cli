package options

import "testing"

func TestFilesListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *FilesListOptions
		wantErr bool
	}{
		{
			name:    "valid with course ID",
			opts:    &FilesListOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "valid with folder ID",
			opts:    &FilesListOptions{FolderID: 1},
			wantErr: false,
		},
		{
			name:    "valid with user ID",
			opts:    &FilesListOptions{UserID: 1},
			wantErr: false,
		},
		{
			name:    "no context specified",
			opts:    &FilesListOptions{},
			wantErr: true,
		},
		{
			name:    "multiple contexts specified",
			opts:    &FilesListOptions{CourseID: 1, FolderID: 2},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("FilesListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFilesGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *FilesGetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &FilesGetOptions{FileID: 1},
			wantErr: false,
		},
		{
			name:    "missing file ID",
			opts:    &FilesGetOptions{FileID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("FilesGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFilesUploadOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *FilesUploadOptions
		wantErr bool
	}{
		{
			name:    "valid with course ID",
			opts:    &FilesUploadOptions{FilePath: "/tmp/file.txt", CourseID: 1},
			wantErr: false,
		},
		{
			name:    "missing file path",
			opts:    &FilesUploadOptions{FilePath: "", CourseID: 1},
			wantErr: true,
		},
		{
			name:    "no context specified",
			opts:    &FilesUploadOptions{FilePath: "/tmp/file.txt"},
			wantErr: true,
		},
		{
			name:    "valid with folder ID",
			opts:    &FilesUploadOptions{FilePath: "/tmp/file.txt", FolderID: 1},
			wantErr: false,
		},
		{
			name:    "multiple contexts",
			opts:    &FilesUploadOptions{FilePath: "/tmp/file.txt", CourseID: 1, UserID: 2},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("FilesUploadOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFilesDownloadOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *FilesDownloadOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &FilesDownloadOptions{FileID: 1},
			wantErr: false,
		},
		{
			name:    "missing file ID",
			opts:    &FilesDownloadOptions{FileID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("FilesDownloadOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFilesDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *FilesDeleteOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &FilesDeleteOptions{FileID: 1},
			wantErr: false,
		},
		{
			name:    "missing file ID",
			opts:    &FilesDeleteOptions{FileID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("FilesDeleteOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFilesQuotaOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *FilesQuotaOptions
		wantErr bool
	}{
		{
			name:    "valid with course ID",
			opts:    &FilesQuotaOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "valid with user ID",
			opts:    &FilesQuotaOptions{UserID: 1},
			wantErr: false,
		},
		{
			name:    "neither course nor user ID",
			opts:    &FilesQuotaOptions{},
			wantErr: true,
		},
		{
			name:    "both course and user ID",
			opts:    &FilesQuotaOptions{CourseID: 1, UserID: 2},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("FilesQuotaOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
