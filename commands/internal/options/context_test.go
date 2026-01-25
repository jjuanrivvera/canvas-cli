package options

import (
	"testing"
)

func TestContextSetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ContextSetOptions
		wantErr bool
	}{
		{
			name: "valid course type",
			opts: &ContextSetOptions{
				Type: "course",
				ID:   123,
			},
			wantErr: false,
		},
		{
			name: "valid assignment type",
			opts: &ContextSetOptions{
				Type: "assignment",
				ID:   456,
			},
			wantErr: false,
		},
		{
			name: "valid user type",
			opts: &ContextSetOptions{
				Type: "user",
				ID:   789,
			},
			wantErr: false,
		},
		{
			name: "valid account type",
			opts: &ContextSetOptions{
				Type: "account",
				ID:   1,
			},
			wantErr: false,
		},
		{
			name: "valid course_id alias",
			opts: &ContextSetOptions{
				Type: "course_id",
				ID:   123,
			},
			wantErr: false,
		},
		{
			name: "valid course-id alias",
			opts: &ContextSetOptions{
				Type: "course-id",
				ID:   123,
			},
			wantErr: false,
		},
		{
			name: "empty type",
			opts: &ContextSetOptions{
				Type: "",
				ID:   123,
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			opts: &ContextSetOptions{
				Type: "invalid",
				ID:   123,
			},
			wantErr: true,
		},
		{
			name: "zero ID",
			opts: &ContextSetOptions{
				Type: "course",
				ID:   0,
			},
			wantErr: true,
		},
		{
			name: "negative ID",
			opts: &ContextSetOptions{
				Type: "course",
				ID:   -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ContextSetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContextShowOptions_Validate(t *testing.T) {
	opts := &ContextShowOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("ContextShowOptions.Validate() error = %v, want nil", err)
	}
}

func TestContextClearOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ContextClearOptions
		wantErr bool
	}{
		{
			name: "empty type clears all",
			opts: &ContextClearOptions{
				Type: "",
			},
			wantErr: false,
		},
		{
			name: "valid course type",
			opts: &ContextClearOptions{
				Type: "course",
			},
			wantErr: false,
		},
		{
			name: "valid assignment type",
			opts: &ContextClearOptions{
				Type: "assignment",
			},
			wantErr: false,
		},
		{
			name: "valid user type",
			opts: &ContextClearOptions{
				Type: "user",
			},
			wantErr: false,
		},
		{
			name: "valid account type",
			opts: &ContextClearOptions{
				Type: "account",
			},
			wantErr: false,
		},
		{
			name: "invalid type",
			opts: &ContextClearOptions{
				Type: "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ContextClearOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
