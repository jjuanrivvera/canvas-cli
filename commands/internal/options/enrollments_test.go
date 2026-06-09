package options

import "testing"

func TestEnrollmentsListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *EnrollmentsListOptions
		wantErr bool
	}{
		{
			name:    "valid with course ID only",
			opts:    &EnrollmentsListOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "valid with user ID only",
			opts:    &EnrollmentsListOptions{UserID: 2},
			wantErr: false,
		},
		{
			name:    "neither course ID nor user ID",
			opts:    &EnrollmentsListOptions{},
			wantErr: true,
		},
		{
			name:    "both course ID and user ID",
			opts:    &EnrollmentsListOptions{CourseID: 1, UserID: 2},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("EnrollmentsListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnrollmentsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *EnrollmentsGetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &EnrollmentsGetOptions{CourseID: 1, EnrollmentID: 2},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &EnrollmentsGetOptions{CourseID: 0, EnrollmentID: 2},
			wantErr: true,
		},
		{
			name:    "zero enrollment ID",
			opts:    &EnrollmentsGetOptions{CourseID: 1, EnrollmentID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("EnrollmentsGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnrollmentsCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *EnrollmentsCreateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &EnrollmentsCreateOptions{CourseID: 1, UserID: 2},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &EnrollmentsCreateOptions{CourseID: 0, UserID: 2},
			wantErr: true,
		},
		{
			name:    "zero user ID",
			opts:    &EnrollmentsCreateOptions{CourseID: 1, UserID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("EnrollmentsCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnrollmentsConcludeOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *EnrollmentsConcludeOptions
		wantErr bool
	}{
		{
			name:    "valid conclude",
			opts:    &EnrollmentsConcludeOptions{CourseID: 1, EnrollmentID: 2, Task: "conclude"},
			wantErr: false,
		},
		{
			name:    "valid deactivate",
			opts:    &EnrollmentsConcludeOptions{CourseID: 1, EnrollmentID: 2, Task: "deactivate"},
			wantErr: false,
		},
		{
			name:    "valid delete",
			opts:    &EnrollmentsConcludeOptions{CourseID: 1, EnrollmentID: 2, Task: "delete"},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &EnrollmentsConcludeOptions{CourseID: 0, EnrollmentID: 2, Task: "conclude"},
			wantErr: true,
		},
		{
			name:    "zero enrollment ID",
			opts:    &EnrollmentsConcludeOptions{CourseID: 1, EnrollmentID: 0, Task: "conclude"},
			wantErr: true,
		},
		{
			name:    "invalid task",
			opts:    &EnrollmentsConcludeOptions{CourseID: 1, EnrollmentID: 2, Task: "archive"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("EnrollmentsConcludeOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnrollmentsReactivateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *EnrollmentsReactivateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &EnrollmentsReactivateOptions{CourseID: 1, EnrollmentID: 2},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &EnrollmentsReactivateOptions{CourseID: 0, EnrollmentID: 2},
			wantErr: true,
		},
		{
			name:    "zero enrollment ID",
			opts:    &EnrollmentsReactivateOptions{CourseID: 1, EnrollmentID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("EnrollmentsReactivateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnrollmentsAcceptOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *EnrollmentsAcceptOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &EnrollmentsAcceptOptions{CourseID: 1, EnrollmentID: 2},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &EnrollmentsAcceptOptions{CourseID: 0, EnrollmentID: 2},
			wantErr: true,
		},
		{
			name:    "zero enrollment ID",
			opts:    &EnrollmentsAcceptOptions{CourseID: 1, EnrollmentID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("EnrollmentsAcceptOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnrollmentsRejectOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *EnrollmentsRejectOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &EnrollmentsRejectOptions{CourseID: 1, EnrollmentID: 2},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &EnrollmentsRejectOptions{CourseID: 0, EnrollmentID: 2},
			wantErr: true,
		},
		{
			name:    "zero enrollment ID",
			opts:    &EnrollmentsRejectOptions{CourseID: 1, EnrollmentID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("EnrollmentsRejectOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
