package options

import "testing"

func TestUsersListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *UsersListOptions
		wantErr bool
	}{
		{
			name:    "valid - no context",
			opts:    &UsersListOptions{},
			wantErr: false,
		},
		{
			name:    "valid - account ID only",
			opts:    &UsersListOptions{AccountID: 1},
			wantErr: false,
		},
		{
			name:    "valid - course ID only",
			opts:    &UsersListOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "both account and course ID",
			opts:    &UsersListOptions{AccountID: 1, CourseID: 2},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("UsersListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsersGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *UsersGetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &UsersGetOptions{UserID: 1},
			wantErr: false,
		},
		{
			name:    "zero user ID",
			opts:    &UsersGetOptions{UserID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("UsersGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsersMeOptions_Validate(t *testing.T) {
	opts := &UsersMeOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("UsersMeOptions.Validate() error = %v, want nil", err)
	}
}

func TestUsersSearchOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *UsersSearchOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &UsersSearchOptions{SearchTerm: "john"},
			wantErr: false,
		},
		{
			name:    "missing search term",
			opts:    &UsersSearchOptions{SearchTerm: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("UsersSearchOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsersCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *UsersCreateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &UsersCreateOptions{AccountID: 1},
			wantErr: false,
		},
		{
			name:    "zero account ID",
			opts:    &UsersCreateOptions{AccountID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("UsersCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsersUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *UsersUpdateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &UsersUpdateOptions{UserID: 1},
			wantErr: false,
		},
		{
			name:    "zero user ID",
			opts:    &UsersUpdateOptions{UserID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("UsersUpdateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
