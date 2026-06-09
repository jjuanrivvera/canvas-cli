package options

import "testing"

func TestAnalyticsActivityOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AnalyticsActivityOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &AnalyticsActivityOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &AnalyticsActivityOptions{CourseID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AnalyticsActivityOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAnalyticsAssignmentsOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AnalyticsAssignmentsOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &AnalyticsAssignmentsOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &AnalyticsAssignmentsOptions{CourseID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AnalyticsAssignmentsOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAnalyticsStudentsOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AnalyticsStudentsOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &AnalyticsStudentsOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &AnalyticsStudentsOptions{CourseID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AnalyticsStudentsOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAnalyticsUserOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AnalyticsUserOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &AnalyticsUserOptions{CourseID: 1, UserID: 2},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &AnalyticsUserOptions{CourseID: 0, UserID: 2},
			wantErr: true,
		},
		{
			name:    "missing user ID",
			opts:    &AnalyticsUserOptions{CourseID: 1, UserID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AnalyticsUserOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAnalyticsDepartmentOptions_Validate(t *testing.T) {
	opts := &AnalyticsDepartmentOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("AnalyticsDepartmentOptions.Validate() error = %v, want nil", err)
	}
}
