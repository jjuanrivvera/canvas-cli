package options

import "testing"

func TestAssignmentsListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AssignmentsListOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &AssignmentsListOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &AssignmentsListOptions{CourseID: 0},
			wantErr: true,
		},
		{
			name:    "negative course ID",
			opts:    &AssignmentsListOptions{CourseID: -1},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AssignmentsListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAssignmentsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AssignmentsGetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &AssignmentsGetOptions{CourseID: 1, AssignmentID: 2},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &AssignmentsGetOptions{CourseID: 0, AssignmentID: 2},
			wantErr: true,
		},
		{
			name:    "zero assignment ID",
			opts:    &AssignmentsGetOptions{CourseID: 1, AssignmentID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AssignmentsGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAssignmentsCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AssignmentsCreateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &AssignmentsCreateOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &AssignmentsCreateOptions{CourseID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AssignmentsCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAssignmentsUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AssignmentsUpdateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &AssignmentsUpdateOptions{CourseID: 1, AssignmentID: 2},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &AssignmentsUpdateOptions{CourseID: 0, AssignmentID: 2},
			wantErr: true,
		},
		{
			name:    "zero assignment ID",
			opts:    &AssignmentsUpdateOptions{CourseID: 1, AssignmentID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AssignmentsUpdateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAssignmentsDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AssignmentsDeleteOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &AssignmentsDeleteOptions{CourseID: 1, AssignmentID: 2},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &AssignmentsDeleteOptions{CourseID: 0, AssignmentID: 2},
			wantErr: true,
		},
		{
			name:    "zero assignment ID",
			opts:    &AssignmentsDeleteOptions{CourseID: 1, AssignmentID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AssignmentsDeleteOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
