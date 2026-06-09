package options

import "testing"

func TestRubricsListOptions_Validate(t *testing.T) {
	opts := &RubricsListOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("RubricsListOptions.Validate() error = %v, want nil", err)
	}
}

func TestRubricsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *RubricsGetOptions
		wantErr bool
	}{
		{
			name:    "valid with course ID",
			opts:    &RubricsGetOptions{CourseID: 1, RubricID: 2},
			wantErr: false,
		},
		{
			name:    "valid with account ID",
			opts:    &RubricsGetOptions{AccountID: 1, RubricID: 2},
			wantErr: false,
		},
		{
			name:    "missing rubric ID",
			opts:    &RubricsGetOptions{CourseID: 1, RubricID: 0},
			wantErr: true,
		},
		{
			name:    "neither course nor account ID",
			opts:    &RubricsGetOptions{RubricID: 2},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("RubricsGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRubricsCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *RubricsCreateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &RubricsCreateOptions{CourseID: 1, Title: "Essay Rubric"},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &RubricsCreateOptions{CourseID: 0, Title: "Essay Rubric"},
			wantErr: true,
		},
		{
			name:    "missing title",
			opts:    &RubricsCreateOptions{CourseID: 1, Title: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("RubricsCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRubricsUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *RubricsUpdateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &RubricsUpdateOptions{CourseID: 1, RubricID: 2},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &RubricsUpdateOptions{CourseID: 0, RubricID: 2},
			wantErr: true,
		},
		{
			name:    "missing rubric ID",
			opts:    &RubricsUpdateOptions{CourseID: 1, RubricID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("RubricsUpdateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRubricsDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *RubricsDeleteOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &RubricsDeleteOptions{CourseID: 1, RubricID: 2},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &RubricsDeleteOptions{CourseID: 0, RubricID: 2},
			wantErr: true,
		},
		{
			name:    "missing rubric ID",
			opts:    &RubricsDeleteOptions{CourseID: 1, RubricID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("RubricsDeleteOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRubricsAssociateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *RubricsAssociateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &RubricsAssociateOptions{CourseID: 1, RubricID: 2, AssignmentID: 3},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &RubricsAssociateOptions{CourseID: 0, RubricID: 2, AssignmentID: 3},
			wantErr: true,
		},
		{
			name:    "missing rubric ID",
			opts:    &RubricsAssociateOptions{CourseID: 1, RubricID: 0, AssignmentID: 3},
			wantErr: true,
		},
		{
			name:    "missing assignment ID",
			opts:    &RubricsAssociateOptions{CourseID: 1, RubricID: 2, AssignmentID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("RubricsAssociateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
