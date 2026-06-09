package options

import "testing"

func TestGroupsListOptions_Validate(t *testing.T) {
	opts := &GroupsListOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("GroupsListOptions.Validate() error = %v, want nil", err)
	}
}

func TestGroupsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GroupsGetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &GroupsGetOptions{GroupID: 1},
			wantErr: false,
		},
		{
			name:    "zero group ID",
			opts:    &GroupsGetOptions{GroupID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GroupsGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGroupsCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GroupsCreateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &GroupsCreateOptions{CategoryID: 1, Name: "Team A"},
			wantErr: false,
		},
		{
			name:    "zero category ID",
			opts:    &GroupsCreateOptions{CategoryID: 0, Name: "Team A"},
			wantErr: true,
		},
		{
			name:    "missing name",
			opts:    &GroupsCreateOptions{CategoryID: 1, Name: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GroupsCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGroupsUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GroupsUpdateOptions
		wantErr bool
	}{
		{
			name:    "valid with one field set",
			opts:    &GroupsUpdateOptions{GroupID: 1, NameSet: true, Name: "New Name"},
			wantErr: false,
		},
		{
			name:    "zero group ID",
			opts:    &GroupsUpdateOptions{GroupID: 0, NameSet: true},
			wantErr: true,
		},
		{
			name:    "no fields set",
			opts:    &GroupsUpdateOptions{GroupID: 1},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GroupsUpdateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGroupsDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GroupsDeleteOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &GroupsDeleteOptions{GroupID: 1},
			wantErr: false,
		},
		{
			name:    "zero group ID",
			opts:    &GroupsDeleteOptions{GroupID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GroupsDeleteOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGroupsMembersListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GroupsMembersListOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &GroupsMembersListOptions{GroupID: 1},
			wantErr: false,
		},
		{
			name:    "zero group ID",
			opts:    &GroupsMembersListOptions{GroupID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GroupsMembersListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGroupsMembersAddOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GroupsMembersAddOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &GroupsMembersAddOptions{GroupID: 1, UserID: 2},
			wantErr: false,
		},
		{
			name:    "zero group ID",
			opts:    &GroupsMembersAddOptions{GroupID: 0, UserID: 2},
			wantErr: true,
		},
		{
			name:    "zero user ID",
			opts:    &GroupsMembersAddOptions{GroupID: 1, UserID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GroupsMembersAddOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGroupsMembersRemoveOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GroupsMembersRemoveOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &GroupsMembersRemoveOptions{GroupID: 1, MembershipID: 2},
			wantErr: false,
		},
		{
			name:    "zero group ID",
			opts:    &GroupsMembersRemoveOptions{GroupID: 0, MembershipID: 2},
			wantErr: true,
		},
		{
			name:    "zero membership ID",
			opts:    &GroupsMembersRemoveOptions{GroupID: 1, MembershipID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GroupsMembersRemoveOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGroupsCategoriesListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GroupsCategoriesListOptions
		wantErr bool
	}{
		{
			name:    "valid with course ID",
			opts:    &GroupsCategoriesListOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "valid with account ID",
			opts:    &GroupsCategoriesListOptions{AccountID: 1},
			wantErr: false,
		},
		{
			name:    "neither course nor account ID",
			opts:    &GroupsCategoriesListOptions{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GroupsCategoriesListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGroupsCategoriesGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GroupsCategoriesGetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &GroupsCategoriesGetOptions{CategoryID: 1},
			wantErr: false,
		},
		{
			name:    "zero category ID",
			opts:    &GroupsCategoriesGetOptions{CategoryID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GroupsCategoriesGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGroupsCategoriesCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GroupsCategoriesCreateOptions
		wantErr bool
	}{
		{
			name:    "valid with course ID",
			opts:    &GroupsCategoriesCreateOptions{CourseID: 1, Name: "Project Groups"},
			wantErr: false,
		},
		{
			name:    "valid with account ID",
			opts:    &GroupsCategoriesCreateOptions{AccountID: 1, Name: "Project Groups"},
			wantErr: false,
		},
		{
			name:    "neither course nor account ID",
			opts:    &GroupsCategoriesCreateOptions{Name: "Project Groups"},
			wantErr: true,
		},
		{
			name:    "missing name",
			opts:    &GroupsCategoriesCreateOptions{CourseID: 1, Name: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GroupsCategoriesCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGroupsCategoriesUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GroupsCategoriesUpdateOptions
		wantErr bool
	}{
		{
			name:    "valid with one field set",
			opts:    &GroupsCategoriesUpdateOptions{CategoryID: 1, NameSet: true, Name: "New Name"},
			wantErr: false,
		},
		{
			name:    "zero category ID",
			opts:    &GroupsCategoriesUpdateOptions{CategoryID: 0, NameSet: true},
			wantErr: true,
		},
		{
			name:    "no fields set",
			opts:    &GroupsCategoriesUpdateOptions{CategoryID: 1},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GroupsCategoriesUpdateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGroupsCategoriesDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GroupsCategoriesDeleteOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &GroupsCategoriesDeleteOptions{CategoryID: 1},
			wantErr: false,
		},
		{
			name:    "zero category ID",
			opts:    &GroupsCategoriesDeleteOptions{CategoryID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GroupsCategoriesDeleteOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGroupsCategoriesGroupsOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GroupsCategoriesGroupsOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &GroupsCategoriesGroupsOptions{CategoryID: 1},
			wantErr: false,
		},
		{
			name:    "zero category ID",
			opts:    &GroupsCategoriesGroupsOptions{CategoryID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GroupsCategoriesGroupsOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
