package options

import "testing"

func TestGradesHistoryOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GradesHistoryOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &GradesHistoryOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &GradesHistoryOptions{CourseID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GradesHistoryOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGradesFeedOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GradesFeedOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &GradesFeedOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &GradesFeedOptions{CourseID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GradesFeedOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGradesColumnsListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GradesColumnsListOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &GradesColumnsListOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &GradesColumnsListOptions{CourseID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GradesColumnsListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGradesColumnsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GradesColumnsGetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &GradesColumnsGetOptions{CourseID: 1, ColumnID: 2},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &GradesColumnsGetOptions{CourseID: 0, ColumnID: 2},
			wantErr: true,
		},
		{
			name:    "missing column ID",
			opts:    &GradesColumnsGetOptions{CourseID: 1, ColumnID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GradesColumnsGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGradesColumnsCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GradesColumnsCreateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &GradesColumnsCreateOptions{CourseID: 1, Title: "Notes"},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &GradesColumnsCreateOptions{CourseID: 0, Title: "Notes"},
			wantErr: true,
		},
		{
			name:    "missing title",
			opts:    &GradesColumnsCreateOptions{CourseID: 1, Title: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GradesColumnsCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGradesColumnsUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GradesColumnsUpdateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &GradesColumnsUpdateOptions{CourseID: 1, ColumnID: 2},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &GradesColumnsUpdateOptions{CourseID: 0, ColumnID: 2},
			wantErr: true,
		},
		{
			name:    "missing column ID",
			opts:    &GradesColumnsUpdateOptions{CourseID: 1, ColumnID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GradesColumnsUpdateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGradesColumnsDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GradesColumnsDeleteOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &GradesColumnsDeleteOptions{CourseID: 1, ColumnID: 2},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &GradesColumnsDeleteOptions{CourseID: 0, ColumnID: 2},
			wantErr: true,
		},
		{
			name:    "missing column ID",
			opts:    &GradesColumnsDeleteOptions{CourseID: 1, ColumnID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GradesColumnsDeleteOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGradesColumnsDataListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GradesColumnsDataListOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &GradesColumnsDataListOptions{CourseID: 1, ColumnID: 2},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &GradesColumnsDataListOptions{CourseID: 0, ColumnID: 2},
			wantErr: true,
		},
		{
			name:    "missing column ID",
			opts:    &GradesColumnsDataListOptions{CourseID: 1, ColumnID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GradesColumnsDataListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGradesColumnsDataSetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GradesColumnsDataSetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &GradesColumnsDataSetOptions{CourseID: 1, ColumnID: 2, UserID: 3, Content: "note"},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &GradesColumnsDataSetOptions{CourseID: 0, ColumnID: 2, UserID: 3, Content: "note"},
			wantErr: true,
		},
		{
			name:    "missing column ID",
			opts:    &GradesColumnsDataSetOptions{CourseID: 1, ColumnID: 0, UserID: 3, Content: "note"},
			wantErr: true,
		},
		{
			name:    "missing user ID",
			opts:    &GradesColumnsDataSetOptions{CourseID: 1, ColumnID: 2, UserID: 0, Content: "note"},
			wantErr: true,
		},
		{
			name:    "missing content",
			opts:    &GradesColumnsDataSetOptions{CourseID: 1, ColumnID: 2, UserID: 3, Content: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GradesColumnsDataSetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
