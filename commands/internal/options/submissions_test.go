package options

import "testing"

func TestSubmissionsListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *SubmissionsListOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &SubmissionsListOptions{CourseID: 1, AssignmentID: 2},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &SubmissionsListOptions{CourseID: 0, AssignmentID: 2},
			wantErr: true,
		},
		{
			name:    "zero assignment ID",
			opts:    &SubmissionsListOptions{CourseID: 1, AssignmentID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SubmissionsListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSubmissionsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *SubmissionsGetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &SubmissionsGetOptions{CourseID: 1, AssignmentID: 2, UserID: 3},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &SubmissionsGetOptions{CourseID: 0, AssignmentID: 2, UserID: 3},
			wantErr: true,
		},
		{
			name:    "zero assignment ID",
			opts:    &SubmissionsGetOptions{CourseID: 1, AssignmentID: 0, UserID: 3},
			wantErr: true,
		},
		{
			name:    "zero user ID",
			opts:    &SubmissionsGetOptions{CourseID: 1, AssignmentID: 2, UserID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SubmissionsGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSubmissionsGradeOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *SubmissionsGradeOptions
		wantErr bool
	}{
		{
			name:    "valid with score",
			opts:    &SubmissionsGradeOptions{CourseID: 1, AssignmentID: 2, UserID: 3, Score: 95.0},
			wantErr: false,
		},
		{
			name:    "valid with comment",
			opts:    &SubmissionsGradeOptions{CourseID: 1, AssignmentID: 2, UserID: 3, Comment: "Great work"},
			wantErr: false,
		},
		{
			name:    "valid with excuse",
			opts:    &SubmissionsGradeOptions{CourseID: 1, AssignmentID: 2, UserID: 3, Excuse: true},
			wantErr: false,
		},
		{
			name:    "valid with posted grade",
			opts:    &SubmissionsGradeOptions{CourseID: 1, AssignmentID: 2, UserID: 3, PostedGrade: "A"},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &SubmissionsGradeOptions{CourseID: 0, AssignmentID: 2, UserID: 3, Score: 90},
			wantErr: true,
		},
		{
			name:    "zero assignment ID",
			opts:    &SubmissionsGradeOptions{CourseID: 1, AssignmentID: 0, UserID: 3, Score: 90},
			wantErr: true,
		},
		{
			name:    "zero user ID",
			opts:    &SubmissionsGradeOptions{CourseID: 1, AssignmentID: 2, UserID: 0, Score: 90},
			wantErr: true,
		},
		{
			name:    "no grading parameter",
			opts:    &SubmissionsGradeOptions{CourseID: 1, AssignmentID: 2, UserID: 3},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SubmissionsGradeOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSubmissionsBulkGradeOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *SubmissionsBulkGradeOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &SubmissionsBulkGradeOptions{CourseID: 1, CSV: "grades.csv"},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &SubmissionsBulkGradeOptions{CourseID: 0, CSV: "grades.csv"},
			wantErr: true,
		},
		{
			name:    "missing CSV",
			opts:    &SubmissionsBulkGradeOptions{CourseID: 1, CSV: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SubmissionsBulkGradeOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSubmissionsCommentsOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *SubmissionsCommentsOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &SubmissionsCommentsOptions{CourseID: 1, AssignmentID: 2, UserID: 3},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &SubmissionsCommentsOptions{CourseID: 0, AssignmentID: 2, UserID: 3},
			wantErr: true,
		},
		{
			name:    "zero assignment ID",
			opts:    &SubmissionsCommentsOptions{CourseID: 1, AssignmentID: 0, UserID: 3},
			wantErr: true,
		},
		{
			name:    "zero user ID",
			opts:    &SubmissionsCommentsOptions{CourseID: 1, AssignmentID: 2, UserID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SubmissionsCommentsOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSubmissionsAddCommentOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *SubmissionsAddCommentOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &SubmissionsAddCommentOptions{CourseID: 1, AssignmentID: 2, UserID: 3, Text: "Good job"},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &SubmissionsAddCommentOptions{CourseID: 0, AssignmentID: 2, UserID: 3, Text: "Good job"},
			wantErr: true,
		},
		{
			name:    "zero assignment ID",
			opts:    &SubmissionsAddCommentOptions{CourseID: 1, AssignmentID: 0, UserID: 3, Text: "Good job"},
			wantErr: true,
		},
		{
			name:    "zero user ID",
			opts:    &SubmissionsAddCommentOptions{CourseID: 1, AssignmentID: 2, UserID: 0, Text: "Good job"},
			wantErr: true,
		},
		{
			name:    "missing text",
			opts:    &SubmissionsAddCommentOptions{CourseID: 1, AssignmentID: 2, UserID: 3, Text: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SubmissionsAddCommentOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSubmissionsDeleteCommentOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *SubmissionsDeleteCommentOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &SubmissionsDeleteCommentOptions{CourseID: 1, AssignmentID: 2, UserID: 3, CommentID: 4},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &SubmissionsDeleteCommentOptions{CourseID: 0, AssignmentID: 2, UserID: 3, CommentID: 4},
			wantErr: true,
		},
		{
			name:    "zero assignment ID",
			opts:    &SubmissionsDeleteCommentOptions{CourseID: 1, AssignmentID: 0, UserID: 3, CommentID: 4},
			wantErr: true,
		},
		{
			name:    "zero user ID",
			opts:    &SubmissionsDeleteCommentOptions{CourseID: 1, AssignmentID: 2, UserID: 0, CommentID: 4},
			wantErr: true,
		},
		{
			name:    "zero comment ID",
			opts:    &SubmissionsDeleteCommentOptions{CourseID: 1, AssignmentID: 2, UserID: 3, CommentID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SubmissionsDeleteCommentOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
