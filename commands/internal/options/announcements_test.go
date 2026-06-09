package options

import "testing"

func TestAnnouncementsListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AnnouncementsListOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &AnnouncementsListOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &AnnouncementsListOptions{CourseID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AnnouncementsListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAnnouncementsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AnnouncementsGetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &AnnouncementsGetOptions{CourseID: 1, AnnouncementID: 2},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &AnnouncementsGetOptions{CourseID: 0, AnnouncementID: 2},
			wantErr: true,
		},
		{
			name:    "missing announcement ID",
			opts:    &AnnouncementsGetOptions{CourseID: 1, AnnouncementID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AnnouncementsGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAnnouncementsCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AnnouncementsCreateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &AnnouncementsCreateOptions{CourseID: 1, Title: "Test"},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &AnnouncementsCreateOptions{CourseID: 0, Title: "Test"},
			wantErr: true,
		},
		{
			name:    "missing title",
			opts:    &AnnouncementsCreateOptions{CourseID: 1, Title: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AnnouncementsCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAnnouncementsUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AnnouncementsUpdateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &AnnouncementsUpdateOptions{CourseID: 1, AnnouncementID: 2},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &AnnouncementsUpdateOptions{CourseID: 0, AnnouncementID: 2},
			wantErr: true,
		},
		{
			name:    "missing announcement ID",
			opts:    &AnnouncementsUpdateOptions{CourseID: 1, AnnouncementID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AnnouncementsUpdateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAnnouncementsDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AnnouncementsDeleteOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &AnnouncementsDeleteOptions{CourseID: 1, AnnouncementID: 2},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &AnnouncementsDeleteOptions{CourseID: 0, AnnouncementID: 2},
			wantErr: true,
		},
		{
			name:    "missing announcement ID",
			opts:    &AnnouncementsDeleteOptions{CourseID: 1, AnnouncementID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AnnouncementsDeleteOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
