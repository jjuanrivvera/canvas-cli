package options

import (
	"encoding/json"
	"fmt"
)

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

// QuizzesQuestionsUpdateOptions contains options for updating a quiz question.
// The *Set fields record which flags were explicitly passed so that only those
// fields are sent to Canvas (an unset flag must not clear the existing value).
type QuizzesQuestionsUpdateOptions struct {
	CourseID   int64
	QuizID     int64
	QuestionID int64

	QuestionName      string
	QuestionText      string
	QuestionType      string
	PointsPossible    float64
	CorrectComments   string
	IncorrectComments string
	Position          int
	AnswersJSON       string // JSON array of answer objects (same shape Canvas returns)

	QuestionNameSet      bool
	QuestionTextSet      bool
	QuestionTypeSet      bool
	PointsPossibleSet    bool
	CorrectCommentsSet   bool
	IncorrectCommentsSet bool
	PositionSet          bool
	AnswersJSONSet       bool
}

// Validate validates the options
func (o *QuizzesQuestionsUpdateOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.QuizID <= 0 {
		return fmt.Errorf("quiz-id is required and must be greater than 0")
	}
	if o.QuestionID <= 0 {
		return fmt.Errorf("question-id is required and must be greater than 0")
	}
	if o.AnswersJSONSet {
		var probe []json.RawMessage
		if err := json.Unmarshal([]byte(o.AnswersJSON), &probe); err != nil {
			return fmt.Errorf("answers-json must be a JSON array of answer objects: %w", err)
		}
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

// QuizzesSubmissionsCreateOptions contains options for starting a quiz submission
type QuizzesSubmissionsCreateOptions struct {
	CourseID   int64
	QuizID     int64
	AccessCode string
}

// Validate validates the options
func (o *QuizzesSubmissionsCreateOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.QuizID <= 0 {
		return fmt.Errorf("quiz-id is required and must be greater than 0")
	}
	return nil
}

// QuizzesGroupsGetOptions contains options for getting a quiz question group
type QuizzesGroupsGetOptions struct {
	CourseID int64
	QuizID   int64
	GroupID  int64
}

// Validate validates the options
func (o *QuizzesGroupsGetOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.QuizID <= 0 {
		return fmt.Errorf("quiz-id is required and must be greater than 0")
	}
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	return nil
}

// QuizzesGroupsCreateOptions contains options for creating a quiz question group
type QuizzesGroupsCreateOptions struct {
	CourseID       int64
	QuizID         int64
	Name           string
	PickCount      int
	QuestionPoints float64
}

// Validate validates the options
func (o *QuizzesGroupsCreateOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.QuizID <= 0 {
		return fmt.Errorf("quiz-id is required and must be greater than 0")
	}
	return nil
}

// QuizzesGroupsUpdateOptions contains options for updating a quiz question group
type QuizzesGroupsUpdateOptions struct {
	CourseID       int64
	QuizID         int64
	GroupID        int64
	Name           string
	PickCount      int
	QuestionPoints float64
	NameSet        bool
	PickCountSet   bool
	PointsSet      bool
}

// Validate validates the options
func (o *QuizzesGroupsUpdateOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.QuizID <= 0 {
		return fmt.Errorf("quiz-id is required and must be greater than 0")
	}
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	if !o.NameSet && !o.PickCountSet && !o.PointsSet {
		return fmt.Errorf("at least one field must be specified for update")
	}
	return nil
}

// QuizzesGroupsDeleteOptions contains options for deleting a quiz question group
type QuizzesGroupsDeleteOptions struct {
	CourseID int64
	QuizID   int64
	GroupID  int64
	Force    bool
}

// Validate validates the options
func (o *QuizzesGroupsDeleteOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.QuizID <= 0 {
		return fmt.Errorf("quiz-id is required and must be greater than 0")
	}
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	return nil
}

// QuizzesReportsListOptions contains options for listing quiz reports
type QuizzesReportsListOptions struct {
	CourseID            int64
	QuizID              int64
	IncludesAllVersions bool
}

// Validate validates the options
func (o *QuizzesReportsListOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.QuizID <= 0 {
		return fmt.Errorf("quiz-id is required and must be greater than 0")
	}
	return nil
}

// QuizzesReportsGetOptions contains options for getting a quiz report
type QuizzesReportsGetOptions struct {
	CourseID int64
	QuizID   int64
	ReportID int64
}

// Validate validates the options
func (o *QuizzesReportsGetOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.QuizID <= 0 {
		return fmt.Errorf("quiz-id is required and must be greater than 0")
	}
	if o.ReportID <= 0 {
		return fmt.Errorf("report-id is required and must be greater than 0")
	}
	return nil
}

// QuizzesReportsCreateOptions contains options for creating a quiz report
type QuizzesReportsCreateOptions struct {
	CourseID            int64
	QuizID              int64
	ReportType          string
	IncludesAllVersions bool
}

// Validate validates the options
func (o *QuizzesReportsCreateOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.QuizID <= 0 {
		return fmt.Errorf("quiz-id is required and must be greater than 0")
	}
	if o.ReportType == "" {
		return fmt.Errorf("report-type is required")
	}
	return nil
}

// QuizzesReportsDeleteOptions contains options for deleting a quiz report
type QuizzesReportsDeleteOptions struct {
	CourseID int64
	QuizID   int64
	ReportID int64
	Force    bool
}

// Validate validates the options
func (o *QuizzesReportsDeleteOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.QuizID <= 0 {
		return fmt.Errorf("quiz-id is required and must be greater than 0")
	}
	if o.ReportID <= 0 {
		return fmt.Errorf("report-id is required and must be greater than 0")
	}
	return nil
}

// QuizzesStatisticsListOptions contains options for listing quiz statistics
type QuizzesStatisticsListOptions struct {
	CourseID int64
	QuizID   int64
}

// Validate validates the options
func (o *QuizzesStatisticsListOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.QuizID <= 0 {
		return fmt.Errorf("quiz-id is required and must be greater than 0")
	}
	return nil
}

// QuizzesExtensionsCreateOptions contains options for creating quiz extensions
type QuizzesExtensionsCreateOptions struct {
	CourseID         int64
	QuizID           int64
	UserID           int64
	ExtraAttempts    int
	ExtraTime        int
	ManuallyUnlocked bool
	ExtendFromNow    int
}

// Validate validates the options
func (o *QuizzesExtensionsCreateOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.QuizID <= 0 {
		return fmt.Errorf("quiz-id is required and must be greater than 0")
	}
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}

// QuizzesIPFiltersListOptions contains options for listing quiz IP filters
type QuizzesIPFiltersListOptions struct {
	CourseID int64
	QuizID   int64
}

// Validate validates the options
func (o *QuizzesIPFiltersListOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.QuizID <= 0 {
		return fmt.Errorf("quiz-id is required and must be greater than 0")
	}
	return nil
}

// QuizzesAssignmentOverridesListOptions contains options for listing quiz assignment overrides
type QuizzesAssignmentOverridesListOptions struct {
	CourseID int64
	QuizIDs  []int64
}

// Validate validates the options
func (o *QuizzesAssignmentOverridesListOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	return nil
}

// QuizzesAssignmentOverridesSetOptions contains options for setting quiz assignment overrides
type QuizzesAssignmentOverridesSetOptions struct {
	CourseID        int64
	QuizID          int64
	CourseSectionID int64
	DueAt           string
	UnlockAt        string
	LockAt          string
}

// Validate validates the options
func (o *QuizzesAssignmentOverridesSetOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.QuizID <= 0 {
		return fmt.Errorf("quiz-id is required and must be greater than 0")
	}
	return nil
}
