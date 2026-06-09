package options

import "testing"

func TestExternalToolsListOptions_Validate(t *testing.T) {
	opts := &ExternalToolsListOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("ExternalToolsListOptions.Validate() error = %v, want nil", err)
	}
}

func TestExternalToolsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ExternalToolsGetOptions
		wantErr bool
	}{
		{
			name:    "valid with course ID",
			opts:    &ExternalToolsGetOptions{CourseID: 1, ToolID: 2},
			wantErr: false,
		},
		{
			name:    "valid with account ID",
			opts:    &ExternalToolsGetOptions{AccountID: 1, ToolID: 2},
			wantErr: false,
		},
		{
			name:    "neither course nor account ID",
			opts:    &ExternalToolsGetOptions{ToolID: 2},
			wantErr: true,
		},
		{
			name:    "missing tool ID",
			opts:    &ExternalToolsGetOptions{CourseID: 1, ToolID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ExternalToolsGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExternalToolsCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ExternalToolsCreateOptions
		wantErr bool
	}{
		{
			name:    "valid with course ID",
			opts:    &ExternalToolsCreateOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "valid with account ID",
			opts:    &ExternalToolsCreateOptions{AccountID: 1},
			wantErr: false,
		},
		{
			name:    "neither course nor account ID",
			opts:    &ExternalToolsCreateOptions{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ExternalToolsCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExternalToolsUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ExternalToolsUpdateOptions
		wantErr bool
	}{
		{
			name:    "valid with course ID",
			opts:    &ExternalToolsUpdateOptions{CourseID: 1, ToolID: 2},
			wantErr: false,
		},
		{
			name:    "valid with account ID",
			opts:    &ExternalToolsUpdateOptions{AccountID: 1, ToolID: 2},
			wantErr: false,
		},
		{
			name:    "neither course nor account ID",
			opts:    &ExternalToolsUpdateOptions{ToolID: 2},
			wantErr: true,
		},
		{
			name:    "missing tool ID",
			opts:    &ExternalToolsUpdateOptions{CourseID: 1, ToolID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ExternalToolsUpdateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExternalToolsDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ExternalToolsDeleteOptions
		wantErr bool
	}{
		{
			name:    "valid with course ID",
			opts:    &ExternalToolsDeleteOptions{CourseID: 1, ToolID: 2},
			wantErr: false,
		},
		{
			name:    "neither course nor account ID",
			opts:    &ExternalToolsDeleteOptions{ToolID: 2},
			wantErr: true,
		},
		{
			name:    "missing tool ID",
			opts:    &ExternalToolsDeleteOptions{CourseID: 1, ToolID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ExternalToolsDeleteOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExternalToolsLaunchOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ExternalToolsLaunchOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ExternalToolsLaunchOptions{CourseID: 1, ToolID: 2},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &ExternalToolsLaunchOptions{CourseID: 0, ToolID: 2},
			wantErr: true,
		},
		{
			name:    "missing tool ID",
			opts:    &ExternalToolsLaunchOptions{CourseID: 1, ToolID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ExternalToolsLaunchOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
