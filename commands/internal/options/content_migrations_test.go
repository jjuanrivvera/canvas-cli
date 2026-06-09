package options

import "testing"

func TestContentMigrationsListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ContentMigrationsListOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ContentMigrationsListOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &ContentMigrationsListOptions{CourseID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ContentMigrationsListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContentMigrationsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ContentMigrationsGetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ContentMigrationsGetOptions{CourseID: 1, MigrationID: 2},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &ContentMigrationsGetOptions{CourseID: 0, MigrationID: 2},
			wantErr: true,
		},
		{
			name:    "missing migration ID",
			opts:    &ContentMigrationsGetOptions{CourseID: 1, MigrationID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ContentMigrationsGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContentMigrationsCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ContentMigrationsCreateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ContentMigrationsCreateOptions{CourseID: 1, Type: "course_copy_importer"},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &ContentMigrationsCreateOptions{CourseID: 0, Type: "course_copy_importer"},
			wantErr: true,
		},
		{
			name:    "missing type",
			opts:    &ContentMigrationsCreateOptions{CourseID: 1, Type: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ContentMigrationsCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContentMigrationsMigratorsOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ContentMigrationsMigratorsOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ContentMigrationsMigratorsOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &ContentMigrationsMigratorsOptions{CourseID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ContentMigrationsMigratorsOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContentMigrationsContentOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ContentMigrationsContentOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ContentMigrationsContentOptions{CourseID: 1, MigrationID: 2},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &ContentMigrationsContentOptions{CourseID: 0, MigrationID: 2},
			wantErr: true,
		},
		{
			name:    "missing migration ID",
			opts:    &ContentMigrationsContentOptions{CourseID: 1, MigrationID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ContentMigrationsContentOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContentMigrationsIssuesOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ContentMigrationsIssuesOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ContentMigrationsIssuesOptions{CourseID: 1, MigrationID: 2},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &ContentMigrationsIssuesOptions{CourseID: 0, MigrationID: 2},
			wantErr: true,
		},
		{
			name:    "missing migration ID",
			opts:    &ContentMigrationsIssuesOptions{CourseID: 1, MigrationID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ContentMigrationsIssuesOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
