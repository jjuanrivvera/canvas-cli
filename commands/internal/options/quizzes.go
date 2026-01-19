package options

import "fmt"

// QuizzesListOptions contains options for listing quizzes
type QuizzesListOptions struct {
	CourseID   int64
	SearchTerm string
}

// Validate validates the options
func (o *QuizzesListOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	return nil
}

// QuizzesGetOptions contains options for getting a quiz
type QuizzesGetOptions struct {
	CourseID int64
	QuizID   int64
}

// Validate validates the options
func (o *QuizzesGetOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.QuizID <= 0 {
		return fmt.Errorf("quiz-id is required and must be greater than 0")
	}
	return nil
}

// QuizzesCreateOptions contains options for creating a quiz
type QuizzesCreateOptions struct {
	CourseID             int64
	Title                string
	Description          string
	QuizType             string
	AssignmentGroupID    int64
	TimeLimit            int
	ShuffleAnswers       bool
	HideResults          string
	ShowCorrectAnswers   bool
	ScoringPolicy        string
	AllowedAttempts      int
	OneQuestionAtATime   bool
	CantGoBack           bool
	AccessCode           string
	IPFilter             string
	DueAt                string
	LockAt               string
	UnlockAt             string
	Published            bool
	AnonymousSubmissions bool
}

// Validate validates the options
func (o *QuizzesCreateOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.Title == "" {
		return fmt.Errorf("title is required")
	}
	return nil
}

// QuizzesUpdateOptions contains options for updating a quiz
type QuizzesUpdateOptions struct {
	CourseID             int64
	QuizID               int64
	Title                string
	Description          string
	QuizType             string
	AssignmentGroupID    int64
	TimeLimit            int
	ShuffleAnswers       bool
	HideResults          string
	ShowCorrectAnswers   bool
	ScoringPolicy        string
	AllowedAttempts      int
	OneQuestionAtATime   bool
	CantGoBack           bool
	AccessCode           string
	IPFilter             string
	DueAt                string
	LockAt               string
	UnlockAt             string
	Published            bool
	AnonymousSubmissions bool
	// Track which fields were actually set
	TitleSet                bool
	DescriptionSet          bool
	QuizTypeSet             bool
	AssignmentGroupIDSet    bool
	TimeLimitSet            bool
	ShuffleAnswersSet       bool
	HideResultsSet          bool
	ShowCorrectAnswersSet   bool
	ScoringPolicySet        bool
	AllowedAttemptsSet      bool
	OneQuestionAtATimeSet   bool
	CantGoBackSet           bool
	AccessCodeSet           bool
	IPFilterSet             bool
	DueAtSet                bool
	LockAtSet               bool
	UnlockAtSet             bool
	PublishedSet            bool
	AnonymousSubmissionsSet bool
}

// Validate validates the options
func (o *QuizzesUpdateOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.QuizID <= 0 {
		return fmt.Errorf("quiz-id is required and must be greater than 0")
	}
	// At least one field must be set for update
	if !o.TitleSet && !o.DescriptionSet && !o.QuizTypeSet && !o.AssignmentGroupIDSet &&
		!o.TimeLimitSet && !o.ShuffleAnswersSet && !o.HideResultsSet && !o.ShowCorrectAnswersSet &&
		!o.ScoringPolicySet && !o.AllowedAttemptsSet && !o.OneQuestionAtATimeSet &&
		!o.CantGoBackSet && !o.AccessCodeSet && !o.IPFilterSet && !o.DueAtSet &&
		!o.LockAtSet && !o.UnlockAtSet && !o.PublishedSet && !o.AnonymousSubmissionsSet {
		return fmt.Errorf("at least one field must be specified for update")
	}
	return nil
}

// QuizzesDeleteOptions contains options for deleting a quiz
type QuizzesDeleteOptions struct {
	CourseID int64
	QuizID   int64
	Force    bool
}

// Validate validates the options
func (o *QuizzesDeleteOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.QuizID <= 0 {
		return fmt.Errorf("quiz-id is required and must be greater than 0")
	}
	return nil
}

// QuizzesQuestionsListOptions contains options for listing quiz questions
type QuizzesQuestionsListOptions struct {
	CourseID int64
	QuizID   int64
}

// Validate validates the options
func (o *QuizzesQuestionsListOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.QuizID <= 0 {
		return fmt.Errorf("quiz-id is required and must be greater than 0")
	}
	return nil
}

// QuizzesQuestionsGetOptions contains options for getting a quiz question
type QuizzesQuestionsGetOptions struct {
	CourseID   int64
	QuizID     int64
	QuestionID int64
}

// Validate validates the options
func (o *QuizzesQuestionsGetOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.QuizID <= 0 {
		return fmt.Errorf("quiz-id is required and must be greater than 0")
	}
	if o.QuestionID <= 0 {
		return fmt.Errorf("question-id is required and must be greater than 0")
	}
	return nil
}

// QuizzesQuestionsCreateOptions contains options for creating a quiz question
type QuizzesQuestionsCreateOptions struct {
	CourseID          int64
	QuizID            int64
	QuestionName      string
	QuestionText      string
	QuestionType      string
	PointsPossible    float64
	CorrectComments   string
	IncorrectComments string
}

// Validate validates the options
func (o *QuizzesQuestionsCreateOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.QuizID <= 0 {
		return fmt.Errorf("quiz-id is required and must be greater than 0")
	}
	if o.QuestionText == "" {
		return fmt.Errorf("text is required")
	}
	return nil
}

// QuizzesQuestionsDeleteOptions contains options for deleting a quiz question
type QuizzesQuestionsDeleteOptions struct {
	CourseID   int64
	QuizID     int64
	QuestionID int64
	Force      bool
}

// Validate validates the options
func (o *QuizzesQuestionsDeleteOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.QuizID <= 0 {
		return fmt.Errorf("quiz-id is required and must be greater than 0")
	}
	if o.QuestionID <= 0 {
		return fmt.Errorf("question-id is required and must be greater than 0")
	}
	return nil
}

// QuizzesSubmissionsListOptions contains options for listing quiz submissions
type QuizzesSubmissionsListOptions struct {
	CourseID int64
	QuizID   int64
}

// Validate validates the options
func (o *QuizzesSubmissionsListOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.QuizID <= 0 {
		return fmt.Errorf("quiz-id is required and must be greater than 0")
	}
	return nil
}

// QuizzesSubmissionsGetOptions contains options for getting a quiz submission
type QuizzesSubmissionsGetOptions struct {
	CourseID     int64
	QuizID       int64
	SubmissionID int64
}

// Validate validates the options
func (o *QuizzesSubmissionsGetOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.QuizID <= 0 {
		return fmt.Errorf("quiz-id is required and must be greater than 0")
	}
	if o.SubmissionID <= 0 {
		return fmt.Errorf("submission-id is required and must be greater than 0")
	}
	return nil
}
