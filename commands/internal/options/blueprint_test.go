package options

import "testing"

func TestBlueprintGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *BlueprintGetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &BlueprintGetOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &BlueprintGetOptions{CourseID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("BlueprintGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBlueprintAssociationsListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *BlueprintAssociationsListOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &BlueprintAssociationsListOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &BlueprintAssociationsListOptions{CourseID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("BlueprintAssociationsListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBlueprintAssociationsAddOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *BlueprintAssociationsAddOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &BlueprintAssociationsAddOptions{CourseID: 1, CourseIDsStr: "2,3"},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &BlueprintAssociationsAddOptions{CourseID: 0, CourseIDsStr: "2,3"},
			wantErr: true,
		},
		{
			name:    "missing course IDs to add",
			opts:    &BlueprintAssociationsAddOptions{CourseID: 1, CourseIDsStr: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("BlueprintAssociationsAddOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBlueprintAssociationsRemoveOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *BlueprintAssociationsRemoveOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &BlueprintAssociationsRemoveOptions{CourseID: 1, CourseIDsStr: "2,3"},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &BlueprintAssociationsRemoveOptions{CourseID: 0, CourseIDsStr: "2,3"},
			wantErr: true,
		},
		{
			name:    "missing course IDs to remove",
			opts:    &BlueprintAssociationsRemoveOptions{CourseID: 1, CourseIDsStr: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("BlueprintAssociationsRemoveOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBlueprintSyncOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *BlueprintSyncOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &BlueprintSyncOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &BlueprintSyncOptions{CourseID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("BlueprintSyncOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBlueprintChangesOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *BlueprintChangesOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &BlueprintChangesOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &BlueprintChangesOptions{CourseID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("BlueprintChangesOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBlueprintMigrationsListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *BlueprintMigrationsListOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &BlueprintMigrationsListOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &BlueprintMigrationsListOptions{CourseID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("BlueprintMigrationsListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBlueprintMigrationsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *BlueprintMigrationsGetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &BlueprintMigrationsGetOptions{CourseID: 1, MigrationID: 2},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &BlueprintMigrationsGetOptions{CourseID: 0, MigrationID: 2},
			wantErr: true,
		},
		{
			name:    "missing migration ID",
			opts:    &BlueprintMigrationsGetOptions{CourseID: 1, MigrationID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("BlueprintMigrationsGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
