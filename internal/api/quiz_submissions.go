package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// QuizSubmissionsService handles quiz submission-related API calls
type QuizSubmissionsService struct {
	client *Client
}

// NewQuizSubmissionsService creates a new quiz submissions service
func NewQuizSubmissionsService(client *Client) *QuizSubmissionsService {
	return &QuizSubmissionsService{client: client}
}

// QuizSubmission represents a Canvas quiz submission
type QuizSubmission struct {
	ID                     int64      `json:"id"`
	QuizID                 int64      `json:"quiz_id"`
	UserID                 int64      `json:"user_id"`
	SubmissionID           int64      `json:"submission_id"`
	StartedAt              *time.Time `json:"started_at,omitempty"`
	FinishedAt             *time.Time `json:"finished_at,omitempty"`
	EndAt                  *time.Time `json:"end_at,omitempty"`
	Attempt                int        `json:"attempt"`
	ExtraAttempts          int        `json:"extra_attempts"`
	ExtraTime              int        `json:"extra_time"`
	ManuallyUnlocked       bool       `json:"manually_unlocked"`
	TimeSpent              int        `json:"time_spent"`
	Score                  float64    `json:"score"`
	ScoreBeforeRegrade     float64    `json:"score_before_regrade,omitempty"`
	KeptScore              float64    `json:"kept_score"`
	FudgePoints            float64    `json:"fudge_points"`
	HasSeenResults         bool       `json:"has_seen_results"`
	WorkflowState          string     `json:"workflow_state"`
	Overdue                bool       `json:"overdue_and_needs_submission"`
	HTMLURL                string     `json:"html_url"`
	ValidationToken        string     `json:"validation_token,omitempty"`
	QuizPointsPossible     float64    `json:"quiz_points_possible"`
	QuestionsRegraded      int        `json:"questions_regraded_count"`
	QuestionsRegradePoints float64    `json:"questions_regraded_since_last_attempt"`
}

// QuizSubmissionsResponse wraps quiz submissions response
type QuizSubmissionsResponse struct {
	QuizSubmissions []QuizSubmission `json:"quiz_submissions"`
}

// ListQuizSubmissionsOptions holds options for listing quiz submissions
type ListQuizSubmissionsOptions struct {
	Include []string // "submission", "quiz", "user"
	Page    int
	PerPage int
}

// List retrieves all submissions for a quiz
func (s *QuizSubmissionsService) List(ctx context.Context, courseID, quizID int64, opts *ListQuizSubmissionsOptions) ([]QuizSubmission, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/%d/submissions", courseID, quizID)

	if opts != nil {
		query := url.Values{}

		for _, include := range opts.Include {
			query.Add("include[]", include)
		}

		if opts.Page > 0 {
			query.Add("page", strconv.Itoa(opts.Page))
		}

		if opts.PerPage > 0 {
			query.Add("per_page", strconv.Itoa(opts.PerPage))
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var response QuizSubmissionsResponse
	if err := s.client.GetJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	return response.QuizSubmissions, nil
}

// Get retrieves a single quiz submission
func (s *QuizSubmissionsService) Get(ctx context.Context, courseID, quizID, submissionID int64, include []string) (*QuizSubmission, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/%d/submissions/%d", courseID, quizID, submissionID)

	if len(include) > 0 {
		query := url.Values{}
		for _, inc := range include {
			query.Add("include[]", inc)
		}
		path += "?" + query.Encode()
	}

	var response QuizSubmissionsResponse
	if err := s.client.GetJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	if len(response.QuizSubmissions) == 0 {
		return nil, fmt.Errorf("quiz submission not found")
	}

	return &response.QuizSubmissions[0], nil
}

// QuizSubmissionQuestionScore is one entry of quiz_submissions[][questions]:
// the score and/or comment for a single answered question, keyed by the quiz
// question ID in UpdateQuizSubmissionParams.Questions.
//
// Canvas API, "Update student question scores and comments"
// (https://canvas.instructure.com/doc/api/quiz_submissions.html#method.quizzes/quiz_submissions_api.update):
// "A set of scores and comments for each question answered by the student.
// The keys are the question IDs, and the values are hashes of score and
// comment entries."
type QuizSubmissionQuestionScore struct {
	Score   *float64 `json:"score,omitempty"`
	Comment *string  `json:"comment,omitempty"`
}

// UpdateQuizSubmissionParams holds parameters for updating a quiz submission.
//
// Attempt is required by Canvas for the documented score/comment update
// ("The attempt number of the quiz submission that should be updated. This
// attempt MUST be already completed.") and is sent whenever set.
type UpdateQuizSubmissionParams struct {
	Attempt          *int                                  `json:"attempt,omitempty"`
	ExtraAttempts    *int                                  `json:"extra_attempts,omitempty"`
	ExtraTime        *int                                  `json:"extra_time,omitempty"`
	ManuallyUnlocked *bool                                 `json:"manually_unlocked,omitempty"`
	FudgePoints      *float64                              `json:"fudge_points,omitempty"`
	Questions        map[int64]QuizSubmissionQuestionScore `json:"questions,omitempty"`
}

// Update updates a quiz submission (for grading adjustments)
func (s *QuizSubmissionsService) Update(ctx context.Context, courseID, quizID, submissionID int64, params *UpdateQuizSubmissionParams) (*QuizSubmission, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/%d/submissions/%d", courseID, quizID, submissionID)

	body := map[string]interface{}{
		"quiz_submissions": []map[string]interface{}{
			{},
		},
	}

	submissions, ok := body["quiz_submissions"].([]map[string]interface{})
	if !ok || len(submissions) == 0 {
		return nil, fmt.Errorf("internal error: invalid submission data structure")
	}

	submissionData := submissions[0]

	if params.Attempt != nil {
		submissionData["attempt"] = *params.Attempt
	}

	if params.ExtraAttempts != nil {
		submissionData["extra_attempts"] = *params.ExtraAttempts
	}

	if params.ExtraTime != nil {
		submissionData["extra_time"] = *params.ExtraTime
	}

	if params.ManuallyUnlocked != nil {
		submissionData["manually_unlocked"] = *params.ManuallyUnlocked
	}

	if params.FudgePoints != nil {
		submissionData["fudge_points"] = *params.FudgePoints
	}

	if len(params.Questions) > 0 {
		// encoding/json renders int64 map keys as decimal strings, which is
		// the {"<question_id>": {...}} shape Canvas documents.
		submissionData["questions"] = params.Questions
	}

	var response QuizSubmissionsResponse
	if err := s.client.PutJSON(ctx, path, body, &response); err != nil {
		return nil, err
	}

	if len(response.QuizSubmissions) == 0 {
		return nil, fmt.Errorf("no submission returned after update")
	}

	return &response.QuizSubmissions[0], nil
}

// Complete marks a quiz submission as complete
func (s *QuizSubmissionsService) Complete(ctx context.Context, courseID, quizID, submissionID int64, attempt int, validationToken string) (*QuizSubmission, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/%d/submissions/%d/complete", courseID, quizID, submissionID)

	body := map[string]interface{}{
		"attempt":          attempt,
		"validation_token": validationToken,
	}

	var response QuizSubmissionsResponse
	if err := s.client.PostJSON(ctx, path, body, &response); err != nil {
		return nil, err
	}

	if len(response.QuizSubmissions) == 0 {
		return nil, fmt.Errorf("no submission returned after complete")
	}

	return &response.QuizSubmissions[0], nil
}

// StartQuizSubmissionParams holds parameters for starting a quiz submission.
type StartQuizSubmissionParams struct {
	AccessCode  string `json:"access_code,omitempty"`
	PreviewMode bool   `json:"preview,omitempty"`
}

// Create starts a new quiz submission (taking a quiz).
// Canvas API: POST /api/v1/courses/:course_id/quizzes/:quiz_id/submissions
func (s *QuizSubmissionsService) Create(ctx context.Context, courseID, quizID int64, params *StartQuizSubmissionParams) (*QuizSubmission, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/%d/submissions", courseID, quizID)

	var body interface{}
	if params != nil {
		body = params
	}

	var response QuizSubmissionsResponse
	if err := s.client.PostJSON(ctx, path, body, &response); err != nil {
		return nil, err
	}

	if len(response.QuizSubmissions) == 0 {
		return nil, fmt.Errorf("no submission returned after create")
	}

	return &response.QuizSubmissions[0], nil
}

// QuizSubmissionEvent represents an event in a quiz submission
type QuizSubmissionEvent struct {
	CreatedAt string      `json:"created_at"`
	EventType string      `json:"event_type"`
	EventData interface{} `json:"event_data,omitempty"`
}

// QuizSubmissionEventsResponse wraps quiz submission events
type QuizSubmissionEventsResponse struct {
	QuizSubmissionEvents []QuizSubmissionEvent `json:"quiz_submission_events"`
}

// QuizSubmissionTime holds timing info for a quiz submission
type QuizSubmissionTime struct {
	EndAt    *time.Time `json:"end_at,omitempty"`
	TimeLeft int        `json:"time_left"`
}

// ListEvents retrieves events for a quiz submission
func (s *QuizSubmissionsService) ListEvents(ctx context.Context, courseID, quizID, submissionID int64) ([]QuizSubmissionEvent, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/%d/submissions/%d/events", courseID, quizID, submissionID)
	var response QuizSubmissionEventsResponse
	if err := s.client.GetJSON(ctx, path, &response); err != nil {
		return nil, err
	}
	return response.QuizSubmissionEvents, nil
}

// CreateEvents submits events for a quiz submission
func (s *QuizSubmissionsService) CreateEvents(ctx context.Context, courseID, quizID, submissionID int64, events []QuizSubmissionEvent) error {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/%d/submissions/%d/events", courseID, quizID, submissionID)
	body := QuizSubmissionEventsResponse{QuizSubmissionEvents: events}
	return s.client.PostJSON(ctx, path, body, nil)
}

// GetTime retrieves timing information for a quiz submission
func (s *QuizSubmissionsService) GetTime(ctx context.Context, courseID, quizID, submissionID int64) (*QuizSubmissionTime, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/%d/submissions/%d/time", courseID, quizID, submissionID)
	var result QuizSubmissionTime
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
