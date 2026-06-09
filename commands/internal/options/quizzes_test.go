package options

import "testing"

func TestQuizzesListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *QuizzesListOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &QuizzesListOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &QuizzesListOptions{CourseID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("QuizzesListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuizzesGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *QuizzesGetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &QuizzesGetOptions{CourseID: 1, QuizID: 2},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &QuizzesGetOptions{CourseID: 0, QuizID: 2},
			wantErr: true,
		},
		{
			name:    "zero quiz ID",
			opts:    &QuizzesGetOptions{CourseID: 1, QuizID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("QuizzesGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuizzesCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *QuizzesCreateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &QuizzesCreateOptions{CourseID: 1, Title: "Midterm"},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &QuizzesCreateOptions{CourseID: 0, Title: "Midterm"},
			wantErr: true,
		},
		{
			name:    "missing title",
			opts:    &QuizzesCreateOptions{CourseID: 1, Title: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("QuizzesCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuizzesUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *QuizzesUpdateOptions
		wantErr bool
	}{
		{
			name:    "valid with one field set",
			opts:    &QuizzesUpdateOptions{CourseID: 1, QuizID: 2, TitleSet: true, Title: "New Title"},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &QuizzesUpdateOptions{CourseID: 0, QuizID: 2, TitleSet: true},
			wantErr: true,
		},
		{
			name:    "zero quiz ID",
			opts:    &QuizzesUpdateOptions{CourseID: 1, QuizID: 0, TitleSet: true},
			wantErr: true,
		},
		{
			name:    "no fields set",
			opts:    &QuizzesUpdateOptions{CourseID: 1, QuizID: 2},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("QuizzesUpdateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuizzesDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *QuizzesDeleteOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &QuizzesDeleteOptions{CourseID: 1, QuizID: 2},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &QuizzesDeleteOptions{CourseID: 0, QuizID: 2},
			wantErr: true,
		},
		{
			name:    "zero quiz ID",
			opts:    &QuizzesDeleteOptions{CourseID: 1, QuizID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("QuizzesDeleteOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuizzesQuestionsListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *QuizzesQuestionsListOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &QuizzesQuestionsListOptions{CourseID: 1, QuizID: 2},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &QuizzesQuestionsListOptions{CourseID: 0, QuizID: 2},
			wantErr: true,
		},
		{
			name:    "zero quiz ID",
			opts:    &QuizzesQuestionsListOptions{CourseID: 1, QuizID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("QuizzesQuestionsListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuizzesQuestionsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *QuizzesQuestionsGetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &QuizzesQuestionsGetOptions{CourseID: 1, QuizID: 2, QuestionID: 3},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &QuizzesQuestionsGetOptions{CourseID: 0, QuizID: 2, QuestionID: 3},
			wantErr: true,
		},
		{
			name:    "zero quiz ID",
			opts:    &QuizzesQuestionsGetOptions{CourseID: 1, QuizID: 0, QuestionID: 3},
			wantErr: true,
		},
		{
			name:    "zero question ID",
			opts:    &QuizzesQuestionsGetOptions{CourseID: 1, QuizID: 2, QuestionID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("QuizzesQuestionsGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuizzesQuestionsCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *QuizzesQuestionsCreateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &QuizzesQuestionsCreateOptions{CourseID: 1, QuizID: 2, QuestionText: "What is 2+2?"},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &QuizzesQuestionsCreateOptions{CourseID: 0, QuizID: 2, QuestionText: "What is 2+2?"},
			wantErr: true,
		},
		{
			name:    "zero quiz ID",
			opts:    &QuizzesQuestionsCreateOptions{CourseID: 1, QuizID: 0, QuestionText: "What is 2+2?"},
			wantErr: true,
		},
		{
			name:    "missing question text",
			opts:    &QuizzesQuestionsCreateOptions{CourseID: 1, QuizID: 2, QuestionText: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("QuizzesQuestionsCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuizzesQuestionsDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *QuizzesQuestionsDeleteOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &QuizzesQuestionsDeleteOptions{CourseID: 1, QuizID: 2, QuestionID: 3},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &QuizzesQuestionsDeleteOptions{CourseID: 0, QuizID: 2, QuestionID: 3},
			wantErr: true,
		},
		{
			name:    "zero quiz ID",
			opts:    &QuizzesQuestionsDeleteOptions{CourseID: 1, QuizID: 0, QuestionID: 3},
			wantErr: true,
		},
		{
			name:    "zero question ID",
			opts:    &QuizzesQuestionsDeleteOptions{CourseID: 1, QuizID: 2, QuestionID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("QuizzesQuestionsDeleteOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuizzesSubmissionsListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *QuizzesSubmissionsListOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &QuizzesSubmissionsListOptions{CourseID: 1, QuizID: 2},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &QuizzesSubmissionsListOptions{CourseID: 0, QuizID: 2},
			wantErr: true,
		},
		{
			name:    "zero quiz ID",
			opts:    &QuizzesSubmissionsListOptions{CourseID: 1, QuizID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("QuizzesSubmissionsListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuizzesSubmissionsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *QuizzesSubmissionsGetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &QuizzesSubmissionsGetOptions{CourseID: 1, QuizID: 2, SubmissionID: 3},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &QuizzesSubmissionsGetOptions{CourseID: 0, QuizID: 2, SubmissionID: 3},
			wantErr: true,
		},
		{
			name:    "zero quiz ID",
			opts:    &QuizzesSubmissionsGetOptions{CourseID: 1, QuizID: 0, SubmissionID: 3},
			wantErr: true,
		},
		{
			name:    "zero submission ID",
			opts:    &QuizzesSubmissionsGetOptions{CourseID: 1, QuizID: 2, SubmissionID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("QuizzesSubmissionsGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
