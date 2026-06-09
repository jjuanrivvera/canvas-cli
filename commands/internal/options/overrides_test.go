package options

import "testing"

func TestOverridesListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *OverridesListOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &OverridesListOptions{CourseID: 1, AssignmentID: 2},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &OverridesListOptions{CourseID: 0, AssignmentID: 2},
			wantErr: true,
		},
		{
			name:    "missing assignment ID",
			opts:    &OverridesListOptions{CourseID: 1, AssignmentID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("OverridesListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOverridesGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *OverridesGetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &OverridesGetOptions{CourseID: 1, AssignmentID: 2, OverrideID: 3},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &OverridesGetOptions{CourseID: 0, AssignmentID: 2, OverrideID: 3},
			wantErr: true,
		},
		{
			name:    "missing assignment ID",
			opts:    &OverridesGetOptions{CourseID: 1, AssignmentID: 0, OverrideID: 3},
			wantErr: true,
		},
		{
			name:    "missing override ID",
			opts:    &OverridesGetOptions{CourseID: 1, AssignmentID: 2, OverrideID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("OverridesGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOverridesCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *OverridesCreateOptions
		wantErr bool
	}{
		{
			name:    "valid with student IDs",
			opts:    &OverridesCreateOptions{CourseID: 1, AssignmentID: 2, StudentIDs: "3,4"},
			wantErr: false,
		},
		{
			name:    "valid with section ID",
			opts:    &OverridesCreateOptions{CourseID: 1, AssignmentID: 2, SectionID: 5},
			wantErr: false,
		},
		{
			name:    "valid with group ID",
			opts:    &OverridesCreateOptions{CourseID: 1, AssignmentID: 2, GroupID: 6},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &OverridesCreateOptions{CourseID: 0, AssignmentID: 2, StudentIDs: "3"},
			wantErr: true,
		},
		{
			name:    "missing assignment ID",
			opts:    &OverridesCreateOptions{CourseID: 1, AssignmentID: 0, StudentIDs: "3"},
			wantErr: true,
		},
		{
			name:    "no target specified",
			opts:    &OverridesCreateOptions{CourseID: 1, AssignmentID: 2},
			wantErr: true,
		},
		{
			name:    "multiple targets specified",
			opts:    &OverridesCreateOptions{CourseID: 1, AssignmentID: 2, StudentIDs: "3", SectionID: 5},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("OverridesCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOverridesUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *OverridesUpdateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &OverridesUpdateOptions{CourseID: 1, AssignmentID: 2, OverrideID: 3},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &OverridesUpdateOptions{CourseID: 0, AssignmentID: 2, OverrideID: 3},
			wantErr: true,
		},
		{
			name:    "missing assignment ID",
			opts:    &OverridesUpdateOptions{CourseID: 1, AssignmentID: 0, OverrideID: 3},
			wantErr: true,
		},
		{
			name:    "missing override ID",
			opts:    &OverridesUpdateOptions{CourseID: 1, AssignmentID: 2, OverrideID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("OverridesUpdateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOverridesDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *OverridesDeleteOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &OverridesDeleteOptions{CourseID: 1, AssignmentID: 2, OverrideID: 3},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &OverridesDeleteOptions{CourseID: 0, AssignmentID: 2, OverrideID: 3},
			wantErr: true,
		},
		{
			name:    "missing assignment ID",
			opts:    &OverridesDeleteOptions{CourseID: 1, AssignmentID: 0, OverrideID: 3},
			wantErr: true,
		},
		{
			name:    "missing override ID",
			opts:    &OverridesDeleteOptions{CourseID: 1, AssignmentID: 2, OverrideID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("OverridesDeleteOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
