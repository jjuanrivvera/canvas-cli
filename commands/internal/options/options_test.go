package options

import (
	"testing"
)

func TestValidateRequired_String(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		value     string
		wantErr   bool
	}{
		{
			name:      "non-empty string",
			fieldName: "name",
			value:     "test",
			wantErr:   false,
		},
		{
			name:      "empty string",
			fieldName: "name",
			value:     "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRequired(tt.fieldName, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRequired() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateRequired_Int64(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		value     int64
		wantErr   bool
	}{
		{
			name:      "non-zero value",
			fieldName: "course-id",
			value:     123,
			wantErr:   false,
		},
		{
			name:      "zero value",
			fieldName: "course-id",
			value:     0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRequired(tt.fieldName, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRequired() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCoursesListOptions_Validate(t *testing.T) {
	opts := &CoursesListOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("CoursesListOptions.Validate() error = %v, want nil", err)
	}
}

func TestCoursesGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *CoursesGetOptions
		wantErr bool
	}{
		{
			name: "valid options",
			opts: &CoursesGetOptions{
				CourseID: 123,
			},
			wantErr: false,
		},
		{
			name: "missing course ID",
			opts: &CoursesGetOptions{
				CourseID: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CoursesGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCoursesCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *CoursesCreateOptions
		wantErr bool
	}{
		{
			name: "valid options",
			opts: &CoursesCreateOptions{
				AccountID: 1,
				Name:      "Test Course",
			},
			wantErr: false,
		},
		{
			name: "missing account ID",
			opts: &CoursesCreateOptions{
				AccountID: 0,
				Name:      "Test Course",
			},
			wantErr: true,
		},
		{
			name: "missing name",
			opts: &CoursesCreateOptions{
				AccountID: 1,
				Name:      "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CoursesCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
