package options

import "testing"

func TestSectionsListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *SectionsListOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &SectionsListOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &SectionsListOptions{CourseID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SectionsListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSectionsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *SectionsGetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &SectionsGetOptions{SectionID: 1},
			wantErr: false,
		},
		{
			name:    "missing section ID",
			opts:    &SectionsGetOptions{SectionID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SectionsGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSectionsCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *SectionsCreateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &SectionsCreateOptions{CourseID: 1, Name: "Section A"},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &SectionsCreateOptions{CourseID: 0, Name: "Section A"},
			wantErr: true,
		},
		{
			name:    "missing name",
			opts:    &SectionsCreateOptions{CourseID: 1, Name: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SectionsCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSectionsUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *SectionsUpdateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &SectionsUpdateOptions{SectionID: 1},
			wantErr: false,
		},
		{
			name:    "missing section ID",
			opts:    &SectionsUpdateOptions{SectionID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SectionsUpdateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSectionsDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *SectionsDeleteOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &SectionsDeleteOptions{SectionID: 1},
			wantErr: false,
		},
		{
			name:    "missing section ID",
			opts:    &SectionsDeleteOptions{SectionID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SectionsDeleteOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSectionsCrosslistOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *SectionsCrosslistOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &SectionsCrosslistOptions{SectionID: 1, NewCourseID: 2},
			wantErr: false,
		},
		{
			name:    "missing section ID",
			opts:    &SectionsCrosslistOptions{SectionID: 0, NewCourseID: 2},
			wantErr: true,
		},
		{
			name:    "missing new course ID",
			opts:    &SectionsCrosslistOptions{SectionID: 1, NewCourseID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SectionsCrosslistOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSectionsUncrosslistOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *SectionsUncrosslistOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &SectionsUncrosslistOptions{SectionID: 1},
			wantErr: false,
		},
		{
			name:    "missing section ID",
			opts:    &SectionsUncrosslistOptions{SectionID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SectionsUncrosslistOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
