package options

import "testing"

func TestOutcomesGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *OutcomesGetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &OutcomesGetOptions{OutcomeID: 1},
			wantErr: false,
		},
		{
			name:    "zero outcome ID",
			opts:    &OutcomesGetOptions{OutcomeID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("OutcomesGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOutcomesCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *OutcomesCreateOptions
		wantErr bool
	}{
		{
			name:    "valid with course ID",
			opts:    &OutcomesCreateOptions{CourseID: 1, GroupID: 2, Title: "Outcome 1"},
			wantErr: false,
		},
		{
			name:    "valid with account ID",
			opts:    &OutcomesCreateOptions{AccountID: 1, GroupID: 2, Title: "Outcome 1"},
			wantErr: false,
		},
		{
			name:    "neither course nor account ID",
			opts:    &OutcomesCreateOptions{GroupID: 2, Title: "Outcome 1"},
			wantErr: true,
		},
		{
			name:    "zero group ID",
			opts:    &OutcomesCreateOptions{CourseID: 1, GroupID: 0, Title: "Outcome 1"},
			wantErr: true,
		},
		{
			name:    "missing title",
			opts:    &OutcomesCreateOptions{CourseID: 1, GroupID: 2, Title: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("OutcomesCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOutcomesUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *OutcomesUpdateOptions
		wantErr bool
	}{
		{
			name:    "valid with one field set",
			opts:    &OutcomesUpdateOptions{OutcomeID: 1, TitleSet: true, Title: "New Title"},
			wantErr: false,
		},
		{
			name:    "zero outcome ID",
			opts:    &OutcomesUpdateOptions{OutcomeID: 0, TitleSet: true},
			wantErr: true,
		},
		{
			name:    "no fields set",
			opts:    &OutcomesUpdateOptions{OutcomeID: 1},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("OutcomesUpdateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOutcomesListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *OutcomesListOptions
		wantErr bool
	}{
		{
			name:    "valid with course ID",
			opts:    &OutcomesListOptions{CourseID: 1, GroupID: 2},
			wantErr: false,
		},
		{
			name:    "valid with account ID",
			opts:    &OutcomesListOptions{AccountID: 1, GroupID: 2},
			wantErr: false,
		},
		{
			name:    "neither course nor account ID",
			opts:    &OutcomesListOptions{GroupID: 2},
			wantErr: true,
		},
		{
			name:    "zero group ID",
			opts:    &OutcomesListOptions{CourseID: 1, GroupID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("OutcomesListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOutcomesLinkOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *OutcomesLinkOptions
		wantErr bool
	}{
		{
			name:    "valid with course ID",
			opts:    &OutcomesLinkOptions{CourseID: 1, GroupID: 2, OutcomeID: 3},
			wantErr: false,
		},
		{
			name:    "valid with account ID",
			opts:    &OutcomesLinkOptions{AccountID: 1, GroupID: 2, OutcomeID: 3},
			wantErr: false,
		},
		{
			name:    "neither course nor account ID",
			opts:    &OutcomesLinkOptions{GroupID: 2, OutcomeID: 3},
			wantErr: true,
		},
		{
			name:    "zero group ID",
			opts:    &OutcomesLinkOptions{CourseID: 1, GroupID: 0, OutcomeID: 3},
			wantErr: true,
		},
		{
			name:    "zero outcome ID",
			opts:    &OutcomesLinkOptions{CourseID: 1, GroupID: 2, OutcomeID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("OutcomesLinkOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOutcomesUnlinkOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *OutcomesUnlinkOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &OutcomesUnlinkOptions{CourseID: 1, GroupID: 2, OutcomeID: 3},
			wantErr: false,
		},
		{
			name:    "neither course nor account ID",
			opts:    &OutcomesUnlinkOptions{GroupID: 2, OutcomeID: 3},
			wantErr: true,
		},
		{
			name:    "zero group ID",
			opts:    &OutcomesUnlinkOptions{CourseID: 1, GroupID: 0, OutcomeID: 3},
			wantErr: true,
		},
		{
			name:    "zero outcome ID",
			opts:    &OutcomesUnlinkOptions{CourseID: 1, GroupID: 2, OutcomeID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("OutcomesUnlinkOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOutcomesGroupsListOptions_Validate(t *testing.T) {
	opts := &OutcomesGroupsListOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("OutcomesGroupsListOptions.Validate() error = %v, want nil", err)
	}
}

func TestOutcomesGroupsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *OutcomesGroupsGetOptions
		wantErr bool
	}{
		{
			name:    "valid with course ID",
			opts:    &OutcomesGroupsGetOptions{CourseID: 1, GroupID: 2},
			wantErr: false,
		},
		{
			name:    "valid with account ID",
			opts:    &OutcomesGroupsGetOptions{AccountID: 1, GroupID: 2},
			wantErr: false,
		},
		{
			name:    "neither course nor account ID",
			opts:    &OutcomesGroupsGetOptions{GroupID: 2},
			wantErr: true,
		},
		{
			name:    "zero group ID",
			opts:    &OutcomesGroupsGetOptions{CourseID: 1, GroupID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("OutcomesGroupsGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOutcomesResultsOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *OutcomesResultsOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &OutcomesResultsOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &OutcomesResultsOptions{CourseID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("OutcomesResultsOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOutcomesAlignmentsOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *OutcomesAlignmentsOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &OutcomesAlignmentsOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &OutcomesAlignmentsOptions{CourseID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("OutcomesAlignmentsOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
