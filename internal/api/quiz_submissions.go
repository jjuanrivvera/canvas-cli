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

// UpdateQuizSubmissionParams holds parameters for updating a quiz submission
type UpdateQuizSubmissionParams struct {
	ExtraAttempts    *int     `json:"extra_attempts,omitempty"`
	ExtraTime        *int     `json:"extra_time,omitempty"`
	ManuallyUnlocked *bool    `json:"manually_unlocked,omitempty"`
	FudgePoints      *float64 `json:"fudge_points,omitempty"`
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
