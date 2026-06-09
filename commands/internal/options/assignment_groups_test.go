package options

import "testing"

func TestAssignmentGroupsListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AssignmentGroupsListOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &AssignmentGroupsListOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &AssignmentGroupsListOptions{CourseID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AssignmentGroupsListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAssignmentGroupsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AssignmentGroupsGetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &AssignmentGroupsGetOptions{CourseID: 1, GroupID: 2},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &AssignmentGroupsGetOptions{CourseID: 0, GroupID: 2},
			wantErr: true,
		},
		{
			name:    "missing group ID",
			opts:    &AssignmentGroupsGetOptions{CourseID: 1, GroupID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AssignmentGroupsGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAssignmentGroupsCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AssignmentGroupsCreateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &AssignmentGroupsCreateOptions{CourseID: 1, Name: "Homework"},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &AssignmentGroupsCreateOptions{CourseID: 0, Name: "Homework"},
			wantErr: true,
		},
		{
			name:    "missing name",
			opts:    &AssignmentGroupsCreateOptions{CourseID: 1, Name: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AssignmentGroupsCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAssignmentGroupsUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AssignmentGroupsUpdateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &AssignmentGroupsUpdateOptions{CourseID: 1, GroupID: 2},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &AssignmentGroupsUpdateOptions{CourseID: 0, GroupID: 2},
			wantErr: true,
		},
		{
			name:    "missing group ID",
			opts:    &AssignmentGroupsUpdateOptions{CourseID: 1, GroupID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AssignmentGroupsUpdateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAssignmentGroupsDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AssignmentGroupsDeleteOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &AssignmentGroupsDeleteOptions{CourseID: 1, GroupID: 2},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &AssignmentGroupsDeleteOptions{CourseID: 0, GroupID: 2},
			wantErr: true,
		},
		{
			name:    "missing group ID",
			opts:    &AssignmentGroupsDeleteOptions{CourseID: 1, GroupID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AssignmentGroupsDeleteOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
