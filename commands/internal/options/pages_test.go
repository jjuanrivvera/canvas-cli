package options

import "testing"

func TestPagesListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *PagesListOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &PagesListOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &PagesListOptions{CourseID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PagesListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPagesGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *PagesGetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &PagesGetOptions{CourseID: 1, URLOrID: "my-page"},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &PagesGetOptions{CourseID: 0, URLOrID: "my-page"},
			wantErr: true,
		},
		{
			name:    "missing URL or ID",
			opts:    &PagesGetOptions{CourseID: 1, URLOrID: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PagesGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPagesFrontOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *PagesFrontOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &PagesFrontOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &PagesFrontOptions{CourseID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PagesFrontOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPagesCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *PagesCreateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &PagesCreateOptions{CourseID: 1, Title: "My Page"},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &PagesCreateOptions{CourseID: 0, Title: "My Page"},
			wantErr: true,
		},
		{
			name:    "missing title",
			opts:    &PagesCreateOptions{CourseID: 1, Title: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PagesCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPagesUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *PagesUpdateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &PagesUpdateOptions{CourseID: 1, URLOrID: "my-page"},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &PagesUpdateOptions{CourseID: 0, URLOrID: "my-page"},
			wantErr: true,
		},
		{
			name:    "missing URL or ID",
			opts:    &PagesUpdateOptions{CourseID: 1, URLOrID: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PagesUpdateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPagesDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *PagesDeleteOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &PagesDeleteOptions{CourseID: 1, URLOrID: "my-page"},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &PagesDeleteOptions{CourseID: 0, URLOrID: "my-page"},
			wantErr: true,
		},
		{
			name:    "missing URL or ID",
			opts:    &PagesDeleteOptions{CourseID: 1, URLOrID: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PagesDeleteOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPagesDuplicateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *PagesDuplicateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &PagesDuplicateOptions{CourseID: 1, URLOrID: "my-page"},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &PagesDuplicateOptions{CourseID: 0, URLOrID: "my-page"},
			wantErr: true,
		},
		{
			name:    "missing URL or ID",
			opts:    &PagesDuplicateOptions{CourseID: 1, URLOrID: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PagesDuplicateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPagesRevisionsOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *PagesRevisionsOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &PagesRevisionsOptions{CourseID: 1, URLOrID: "my-page"},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &PagesRevisionsOptions{CourseID: 0, URLOrID: "my-page"},
			wantErr: true,
		},
		{
			name:    "missing URL or ID",
			opts:    &PagesRevisionsOptions{CourseID: 1, URLOrID: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PagesRevisionsOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPagesRevertOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *PagesRevertOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &PagesRevertOptions{CourseID: 1, URLOrID: "my-page", RevisionID: 3},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &PagesRevertOptions{CourseID: 0, URLOrID: "my-page", RevisionID: 3},
			wantErr: true,
		},
		{
			name:    "missing URL or ID",
			opts:    &PagesRevertOptions{CourseID: 1, URLOrID: "", RevisionID: 3},
			wantErr: true,
		},
		{
			name:    "zero revision ID",
			opts:    &PagesRevertOptions{CourseID: 1, URLOrID: "my-page", RevisionID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PagesRevertOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
