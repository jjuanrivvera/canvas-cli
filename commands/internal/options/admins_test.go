package options

import "testing"

func TestAdminsListOptions_Validate(t *testing.T) {
	opts := &AdminsListOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("AdminsListOptions.Validate() error = %v, want nil", err)
	}
}

func TestAdminsAddOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AdminsAddOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &AdminsAddOptions{AccountID: 1, UserID: 2},
			wantErr: false,
		},
		{
			name:    "missing account ID",
			opts:    &AdminsAddOptions{AccountID: 0, UserID: 2},
			wantErr: true,
		},
		{
			name:    "missing user ID",
			opts:    &AdminsAddOptions{AccountID: 1, UserID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminsAddOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAdminsRemoveOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AdminsRemoveOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &AdminsRemoveOptions{AccountID: 1, UserID: 2},
			wantErr: false,
		},
		{
			name:    "missing account ID",
			opts:    &AdminsRemoveOptions{AccountID: 0, UserID: 2},
			wantErr: true,
		},
		{
			name:    "missing user ID",
			opts:    &AdminsRemoveOptions{AccountID: 1, UserID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminsRemoveOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRolesListOptions_Validate(t *testing.T) {
	opts := &RolesListOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("RolesListOptions.Validate() error = %v, want nil", err)
	}
}

func TestRolesGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *RolesGetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &RolesGetOptions{AccountID: 1, RoleID: 2},
			wantErr: false,
		},
		{
			name:    "missing account ID",
			opts:    &RolesGetOptions{AccountID: 0, RoleID: 2},
			wantErr: true,
		},
		{
			name:    "missing role ID",
			opts:    &RolesGetOptions{AccountID: 1, RoleID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("RolesGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRolesCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *RolesCreateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &RolesCreateOptions{AccountID: 1, Label: "My Role"},
			wantErr: false,
		},
		{
			name:    "missing account ID",
			opts:    &RolesCreateOptions{AccountID: 0, Label: "My Role"},
			wantErr: true,
		},
		{
			name:    "missing label",
			opts:    &RolesCreateOptions{AccountID: 1, Label: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("RolesCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRolesUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *RolesUpdateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &RolesUpdateOptions{AccountID: 1, RoleID: 2},
			wantErr: false,
		},
		{
			name:    "missing account ID",
			opts:    &RolesUpdateOptions{AccountID: 0, RoleID: 2},
			wantErr: true,
		},
		{
			name:    "missing role ID",
			opts:    &RolesUpdateOptions{AccountID: 1, RoleID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("RolesUpdateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRolesDeactivateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *RolesDeactivateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &RolesDeactivateOptions{AccountID: 1, RoleID: 2},
			wantErr: false,
		},
		{
			name:    "missing account ID",
			opts:    &RolesDeactivateOptions{AccountID: 0, RoleID: 2},
			wantErr: true,
		},
		{
			name:    "missing role ID",
			opts:    &RolesDeactivateOptions{AccountID: 1, RoleID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("RolesDeactivateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRolesActivateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *RolesActivateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &RolesActivateOptions{AccountID: 1, RoleID: 2},
			wantErr: false,
		},
		{
			name:    "missing account ID",
			opts:    &RolesActivateOptions{AccountID: 0, RoleID: 2},
			wantErr: true,
		},
		{
			name:    "missing role ID",
			opts:    &RolesActivateOptions{AccountID: 1, RoleID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("RolesActivateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
