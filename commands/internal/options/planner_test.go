package options

import "testing"

func TestPlannerItemsOptions_Validate(t *testing.T) {
	opts := &PlannerItemsOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("PlannerItemsOptions.Validate() error = %v, want nil", err)
	}
}

func TestPlannerNotesListOptions_Validate(t *testing.T) {
	opts := &PlannerNotesListOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("PlannerNotesListOptions.Validate() error = %v, want nil", err)
	}
}

func TestPlannerNotesGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *PlannerNotesGetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &PlannerNotesGetOptions{NoteID: 1},
			wantErr: false,
		},
		{
			name:    "missing note ID",
			opts:    &PlannerNotesGetOptions{NoteID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PlannerNotesGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPlannerNotesCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *PlannerNotesCreateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &PlannerNotesCreateOptions{Title: "Study session"},
			wantErr: false,
		},
		{
			name:    "missing title",
			opts:    &PlannerNotesCreateOptions{Title: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PlannerNotesCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPlannerNotesUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *PlannerNotesUpdateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &PlannerNotesUpdateOptions{NoteID: 1},
			wantErr: false,
		},
		{
			name:    "missing note ID",
			opts:    &PlannerNotesUpdateOptions{NoteID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PlannerNotesUpdateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPlannerNotesDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *PlannerNotesDeleteOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &PlannerNotesDeleteOptions{NoteID: 1},
			wantErr: false,
		},
		{
			name:    "missing note ID",
			opts:    &PlannerNotesDeleteOptions{NoteID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PlannerNotesDeleteOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPlannerCompleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *PlannerCompleteOptions
		wantErr bool
	}{
		{
			name:    "valid Assignment",
			opts:    &PlannerCompleteOptions{PlannableType: "Assignment", PlannableID: 1},
			wantErr: false,
		},
		{
			name:    "valid Quiz",
			opts:    &PlannerCompleteOptions{PlannableType: "Quiz", PlannableID: 1},
			wantErr: false,
		},
		{
			name:    "valid DiscussionTopic",
			opts:    &PlannerCompleteOptions{PlannableType: "DiscussionTopic", PlannableID: 1},
			wantErr: false,
		},
		{
			name:    "valid WikiPage",
			opts:    &PlannerCompleteOptions{PlannableType: "WikiPage", PlannableID: 1},
			wantErr: false,
		},
		{
			name:    "valid CalendarEvent",
			opts:    &PlannerCompleteOptions{PlannableType: "CalendarEvent", PlannableID: 1},
			wantErr: false,
		},
		{
			name:    "valid PlannerNote",
			opts:    &PlannerCompleteOptions{PlannableType: "PlannerNote", PlannableID: 1},
			wantErr: false,
		},
		{
			name:    "valid Announcement",
			opts:    &PlannerCompleteOptions{PlannableType: "Announcement", PlannableID: 1},
			wantErr: false,
		},
		{
			name:    "missing plannable type",
			opts:    &PlannerCompleteOptions{PlannableType: "", PlannableID: 1},
			wantErr: true,
		},
		{
			name:    "invalid plannable type",
			opts:    &PlannerCompleteOptions{PlannableType: "InvalidType", PlannableID: 1},
			wantErr: true,
		},
		{
			name:    "missing plannable ID",
			opts:    &PlannerCompleteOptions{PlannableType: "Assignment", PlannableID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PlannerCompleteOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPlannerDismissOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *PlannerDismissOptions
		wantErr bool
	}{
		{
			name:    "valid Assignment",
			opts:    &PlannerDismissOptions{PlannableType: "Assignment", PlannableID: 1},
			wantErr: false,
		},
		{
			name:    "missing plannable type",
			opts:    &PlannerDismissOptions{PlannableType: "", PlannableID: 1},
			wantErr: true,
		},
		{
			name:    "invalid plannable type",
			opts:    &PlannerDismissOptions{PlannableType: "Invalid", PlannableID: 1},
			wantErr: true,
		},
		{
			name:    "missing plannable ID",
			opts:    &PlannerDismissOptions{PlannableType: "Assignment", PlannableID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PlannerDismissOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPlannerOverridesOptions_Validate(t *testing.T) {
	opts := &PlannerOverridesOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("PlannerOverridesOptions.Validate() error = %v, want nil", err)
	}
}
