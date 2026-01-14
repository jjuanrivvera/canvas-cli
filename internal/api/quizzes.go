package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// QuizzesService handles quiz-related API calls
type QuizzesService struct {
	client *Client
}

// NewQuizzesService creates a new quizzes service
func NewQuizzesService(client *Client) *QuizzesService {
	return &QuizzesService{client: client}
}

// Quiz represents a Canvas quiz
type Quiz struct {
	ID                            int64            `json:"id"`
	Title                         string           `json:"title"`
	HTMLURL                       string           `json:"html_url"`
	MobileURL                     string           `json:"mobile_url"`
	Description                   string           `json:"description"`
	QuizType                      string           `json:"quiz_type"`
	AssignmentGroupID             int64            `json:"assignment_group_id,omitempty"`
	TimeLimit                     int              `json:"time_limit,omitempty"`
	ShuffleAnswers                bool             `json:"shuffle_answers"`
	HideResults                   string           `json:"hide_results,omitempty"`
	ShowCorrectAnswers            bool             `json:"show_correct_answers"`
	ShowCorrectAnswersLastAttempt bool             `json:"show_correct_answers_last_attempt"`
	ShowCorrectAnswersAt          *time.Time       `json:"show_correct_answers_at,omitempty"`
	HideCorrectAnswersAt          *time.Time       `json:"hide_correct_answers_at,omitempty"`
	OneTimeResults                bool             `json:"one_time_results"`
	ScoringPolicy                 string           `json:"scoring_policy,omitempty"`
	AllowedAttempts               int              `json:"allowed_attempts"`
	OneQuestionAtATime            bool             `json:"one_question_at_a_time"`
	QuestionCount                 int              `json:"question_count"`
	PointsPossible                float64          `json:"points_possible"`
	CantGoBack                    bool             `json:"cant_go_back"`
	AccessCode                    string           `json:"access_code,omitempty"`
	IPFilter                      string           `json:"ip_filter,omitempty"`
	DueAt                         *time.Time       `json:"due_at,omitempty"`
	LockAt                        *time.Time       `json:"lock_at,omitempty"`
	UnlockAt                      *time.Time       `json:"unlock_at,omitempty"`
	Published                     bool             `json:"published"`
	Unpublishable                 bool             `json:"unpublishable"`
	LockedForUser                 bool             `json:"locked_for_user"`
	LockInfo                      *LockInfo        `json:"lock_info,omitempty"`
	LockExplanation               string           `json:"lock_explanation,omitempty"`
	SpeedGraderURL                string           `json:"speedgrader_url,omitempty"`
	QuizExtensionsURL             string           `json:"quiz_extensions_url,omitempty"`
	Permissions                   *QuizPermissions `json:"permissions,omitempty"`
	AllDates                      []DateSet        `json:"all_dates,omitempty"`
	VersionNumber                 int              `json:"version_number"`
	QuestionTypes                 []string         `json:"question_types,omitempty"`
	AnonymousSubmissions          bool             `json:"anonymous_submissions"`
}

// QuizPermissions represents quiz permissions
type QuizPermissions struct {
	Read           bool `json:"read"`
	Submit         bool `json:"submit"`
	Create         bool `json:"create"`
	Manage         bool `json:"manage"`
	ReadStatistics bool `json:"read_statistics"`
	ReviewGrades   bool `json:"review_grades"`
	Update         bool `json:"update"`
	Delete         bool `json:"delete"`
}

// DateSet represents availability dates
type DateSet struct {
	ID       int64      `json:"id,omitempty"`
	DueAt    *time.Time `json:"due_at,omitempty"`
	UnlockAt *time.Time `json:"unlock_at,omitempty"`
	LockAt   *time.Time `json:"lock_at,omitempty"`
	Base     bool       `json:"base,omitempty"`
	Title    string     `json:"title,omitempty"`
}

// ListQuizzesOptions holds options for listing quizzes
type ListQuizzesOptions struct {
	SearchTerm string
	Page       int
	PerPage    int
}

// List retrieves all quizzes for a course
func (s *QuizzesService) List(ctx context.Context, courseID int64, opts *ListQuizzesOptions) ([]Quiz, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes", courseID)

	if opts != nil {
		query := url.Values{}

		if opts.SearchTerm != "" {
			query.Add("search_term", opts.SearchTerm)
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

	var quizzes []Quiz
	if err := s.client.GetAllPages(ctx, path, &quizzes); err != nil {
		return nil, err
	}

	return quizzes, nil
}

// Get retrieves a single quiz
func (s *QuizzesService) Get(ctx context.Context, courseID, quizID int64) (*Quiz, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/%d", courseID, quizID)

	var quiz Quiz
	if err := s.client.GetJSON(ctx, path, &quiz); err != nil {
		return nil, err
	}

	return &quiz, nil
}

// CreateQuizParams holds parameters for creating a quiz
type CreateQuizParams struct {
	Title                         string
	Description                   string
	QuizType                      string
	AssignmentGroupID             int64
	TimeLimit                     int
	ShuffleAnswers                bool
	HideResults                   string
	ShowCorrectAnswers            bool
	ShowCorrectAnswersLastAttempt bool
	ShowCorrectAnswersAt          string
	HideCorrectAnswersAt          string
	OneTimeResults                bool
	ScoringPolicy                 string
	AllowedAttempts               int
	OneQuestionAtATime            bool
	CantGoBack                    bool
	AccessCode                    string
	IPFilter                      string
	DueAt                         string
	LockAt                        string
	UnlockAt                      string
	Published                     bool
	OnlyVisibleToOverrides        bool
	AnonymousSubmissions          bool
}

// Create creates a new quiz
func (s *QuizzesService) Create(ctx context.Context, courseID int64, params *CreateQuizParams) (*Quiz, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes", courseID)

	body := map[string]interface{}{
		"quiz": make(map[string]interface{}),
	}

	quizData, ok := body["quiz"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("internal error: invalid quiz data structure")
	}

	if params.Title != "" {
		quizData["title"] = params.Title
	}

	if params.Description != "" {
		quizData["description"] = params.Description
	}

	if params.QuizType != "" {
		quizData["quiz_type"] = params.QuizType
	}

	if params.AssignmentGroupID > 0 {
		quizData["assignment_group_id"] = params.AssignmentGroupID
	}

	if params.TimeLimit > 0 {
		quizData["time_limit"] = params.TimeLimit
	}

	if params.ShuffleAnswers {
		quizData["shuffle_answers"] = params.ShuffleAnswers
	}

	if params.HideResults != "" {
		quizData["hide_results"] = params.HideResults
	}

	if params.ShowCorrectAnswers {
		quizData["show_correct_answers"] = params.ShowCorrectAnswers
	}

	if params.ShowCorrectAnswersLastAttempt {
		quizData["show_correct_answers_last_attempt"] = params.ShowCorrectAnswersLastAttempt
	}

	if params.ShowCorrectAnswersAt != "" {
		quizData["show_correct_answers_at"] = params.ShowCorrectAnswersAt
	}

	if params.HideCorrectAnswersAt != "" {
		quizData["hide_correct_answers_at"] = params.HideCorrectAnswersAt
	}

	if params.OneTimeResults {
		quizData["one_time_results"] = params.OneTimeResults
	}

	if params.ScoringPolicy != "" {
		quizData["scoring_policy"] = params.ScoringPolicy
	}

	if params.AllowedAttempts != 0 {
		quizData["allowed_attempts"] = params.AllowedAttempts
	}

	if params.OneQuestionAtATime {
		quizData["one_question_at_a_time"] = params.OneQuestionAtATime
	}

	if params.CantGoBack {
		quizData["cant_go_back"] = params.CantGoBack
	}

	if params.AccessCode != "" {
		quizData["access_code"] = params.AccessCode
	}

	if params.IPFilter != "" {
		quizData["ip_filter"] = params.IPFilter
	}

	if params.DueAt != "" {
		quizData["due_at"] = params.DueAt
	}

	if params.LockAt != "" {
		quizData["lock_at"] = params.LockAt
	}

	if params.UnlockAt != "" {
		quizData["unlock_at"] = params.UnlockAt
	}

	if params.Published {
		quizData["published"] = params.Published
	}

	if params.OnlyVisibleToOverrides {
		quizData["only_visible_to_overrides"] = params.OnlyVisibleToOverrides
	}

	if params.AnonymousSubmissions {
		quizData["anonymous_submissions"] = params.AnonymousSubmissions
	}

	var quiz Quiz
	if err := s.client.PostJSON(ctx, path, body, &quiz); err != nil {
		return nil, err
	}

	return &quiz, nil
}

// UpdateQuizParams holds parameters for updating a quiz
type UpdateQuizParams struct {
	Title                         *string
	Description                   *string
	QuizType                      *string
	AssignmentGroupID             *int64
	TimeLimit                     *int
	ShuffleAnswers                *bool
	HideResults                   *string
	ShowCorrectAnswers            *bool
	ShowCorrectAnswersLastAttempt *bool
	ShowCorrectAnswersAt          *string
	HideCorrectAnswersAt          *string
	OneTimeResults                *bool
	ScoringPolicy                 *string
	AllowedAttempts               *int
	OneQuestionAtATime            *bool
	CantGoBack                    *bool
	AccessCode                    *string
	IPFilter                      *string
	DueAt                         *string
	LockAt                        *string
	UnlockAt                      *string
	Published                     *bool
	OnlyVisibleToOverrides        *bool
	AnonymousSubmissions          *bool
}

// Update updates an existing quiz
func (s *QuizzesService) Update(ctx context.Context, courseID, quizID int64, params *UpdateQuizParams) (*Quiz, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/%d", courseID, quizID)

	body := map[string]interface{}{
		"quiz": make(map[string]interface{}),
	}

	quizData, ok := body["quiz"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("internal error: invalid quiz data structure")
	}

	if params.Title != nil {
		quizData["title"] = *params.Title
	}

	if params.Description != nil {
		quizData["description"] = *params.Description
	}

	if params.QuizType != nil {
		quizData["quiz_type"] = *params.QuizType
	}

	if params.AssignmentGroupID != nil {
		quizData["assignment_group_id"] = *params.AssignmentGroupID
	}

	if params.TimeLimit != nil {
		quizData["time_limit"] = *params.TimeLimit
	}

	if params.ShuffleAnswers != nil {
		quizData["shuffle_answers"] = *params.ShuffleAnswers
	}

	if params.HideResults != nil {
		quizData["hide_results"] = *params.HideResults
	}

	if params.ShowCorrectAnswers != nil {
		quizData["show_correct_answers"] = *params.ShowCorrectAnswers
	}

	if params.ShowCorrectAnswersLastAttempt != nil {
		quizData["show_correct_answers_last_attempt"] = *params.ShowCorrectAnswersLastAttempt
	}

	if params.ShowCorrectAnswersAt != nil {
		quizData["show_correct_answers_at"] = *params.ShowCorrectAnswersAt
	}

	if params.HideCorrectAnswersAt != nil {
		quizData["hide_correct_answers_at"] = *params.HideCorrectAnswersAt
	}

	if params.OneTimeResults != nil {
		quizData["one_time_results"] = *params.OneTimeResults
	}

	if params.ScoringPolicy != nil {
		quizData["scoring_policy"] = *params.ScoringPolicy
	}

	if params.AllowedAttempts != nil {
		quizData["allowed_attempts"] = *params.AllowedAttempts
	}

	if params.OneQuestionAtATime != nil {
		quizData["one_question_at_a_time"] = *params.OneQuestionAtATime
	}

	if params.CantGoBack != nil {
		quizData["cant_go_back"] = *params.CantGoBack
	}

	if params.AccessCode != nil {
		quizData["access_code"] = *params.AccessCode
	}

	if params.IPFilter != nil {
		quizData["ip_filter"] = *params.IPFilter
	}

	if params.DueAt != nil {
		quizData["due_at"] = *params.DueAt
	}

	if params.LockAt != nil {
		quizData["lock_at"] = *params.LockAt
	}

	if params.UnlockAt != nil {
		quizData["unlock_at"] = *params.UnlockAt
	}

	if params.Published != nil {
		quizData["published"] = *params.Published
	}

	if params.OnlyVisibleToOverrides != nil {
		quizData["only_visible_to_overrides"] = *params.OnlyVisibleToOverrides
	}

	if params.AnonymousSubmissions != nil {
		quizData["anonymous_submissions"] = *params.AnonymousSubmissions
	}

	var quiz Quiz
	if err := s.client.PutJSON(ctx, path, body, &quiz); err != nil {
		return nil, err
	}

	return &quiz, nil
}

// Delete deletes a quiz
func (s *QuizzesService) Delete(ctx context.Context, courseID, quizID int64) (*Quiz, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/%d", courseID, quizID)

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var quiz Quiz
	if err := json.NewDecoder(resp.Body).Decode(&quiz); err != nil {
		return nil, err
	}

	return &quiz, nil
}

// ReorderQuizItems reorders quiz questions
func (s *QuizzesService) ReorderQuizItems(ctx context.Context, courseID, quizID int64, order []ReorderItem) error {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/%d/reorder", courseID, quizID)

	body := map[string]interface{}{
		"order": order,
	}

	return s.client.PostJSON(ctx, path, body, nil)
}

// ReorderItem represents an item in a reorder request
type ReorderItem struct {
	ID   int64  `json:"id"`
	Type string `json:"type"` // "question" or "group"
}
