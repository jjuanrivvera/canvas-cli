package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// QuizQuestionsService handles quiz question-related API calls
type QuizQuestionsService struct {
	client *Client
}

// NewQuizQuestionsService creates a new quiz questions service
func NewQuizQuestionsService(client *Client) *QuizQuestionsService {
	return &QuizQuestionsService{client: client}
}

// QuizQuestion represents a Canvas quiz question
type QuizQuestion struct {
	ID                    int64        `json:"id"`
	QuizID                int64        `json:"quiz_id"`
	Position              int          `json:"position"`
	QuestionName          string       `json:"question_name"`
	QuestionType          string       `json:"question_type"`
	QuestionText          string       `json:"question_text"`
	PointsPossible        float64      `json:"points_possible"`
	CorrectComments       string       `json:"correct_comments,omitempty"`
	IncorrectComments     string       `json:"incorrect_comments,omitempty"`
	NeutralComments       string       `json:"neutral_comments,omitempty"`
	CorrectCommentsHTML   string       `json:"correct_comments_html,omitempty"`
	IncorrectCommentsHTML string       `json:"incorrect_comments_html,omitempty"`
	NeutralCommentsHTML   string       `json:"neutral_comments_html,omitempty"`
	Answers               []QuizAnswer `json:"answers,omitempty"`
	Variables             interface{}  `json:"variables,omitempty"`
	Formulas              interface{}  `json:"formulas,omitempty"`
	AnswerToleranceMargin float64      `json:"answer_tolerance,omitempty"`
	FormulaTolerance      float64      `json:"formula_decimal_places,omitempty"`
	Matches               []QuizMatch  `json:"matches,omitempty"`
	MatchingAnswerIn      string       `json:"matching_answer_incorrect_matches,omitempty"`
}

// QuizAnswer represents an answer choice for a quiz question
type QuizAnswer struct {
	ID            int64   `json:"id,omitempty"`
	Text          string  `json:"text,omitempty"`
	HTML          string  `json:"html,omitempty"`
	Comments      string  `json:"comments,omitempty"`
	CommentsHTML  string  `json:"comments_html,omitempty"`
	Weight        float64 `json:"weight,omitempty"`
	BlankID       string  `json:"blank_id,omitempty"`
	MatchID       int64   `json:"match_id,omitempty"`
	Left          string  `json:"left,omitempty"`
	Right         string  `json:"right,omitempty"`
	NumericalMin  float64 `json:"numerical_answer_type,omitempty"`
	NumericalMax  float64 `json:"end,omitempty"`
	ExactAnswer   float64 `json:"exact,omitempty"`
	ErrorMargin   float64 `json:"margin,omitempty"`
	ApproxAnswer  float64 `json:"approximate,omitempty"`
	ApproxPrecise float64 `json:"precision,omitempty"`
}

// QuizMatch represents a matching pair
type QuizMatch struct {
	Text    string `json:"text"`
	MatchID int64  `json:"match_id"`
}

// ListQuizQuestionsOptions holds options for listing quiz questions
type ListQuizQuestionsOptions struct {
	Page    int
	PerPage int
}

// List retrieves all questions for a quiz
func (s *QuizQuestionsService) List(ctx context.Context, courseID, quizID int64, opts *ListQuizQuestionsOptions) ([]QuizQuestion, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/%d/questions", courseID, quizID)

	if opts != nil {
		query := url.Values{}

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

	var questions []QuizQuestion
	if err := s.client.GetAllPages(ctx, path, &questions); err != nil {
		return nil, err
	}

	return questions, nil
}

// Get retrieves a single quiz question
func (s *QuizQuestionsService) Get(ctx context.Context, courseID, quizID, questionID int64) (*QuizQuestion, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/%d/questions/%d", courseID, quizID, questionID)

	var question QuizQuestion
	if err := s.client.GetJSON(ctx, path, &question); err != nil {
		return nil, err
	}

	return &question, nil
}

// CreateQuizQuestionParams holds parameters for creating a quiz question
type CreateQuizQuestionParams struct {
	QuestionName      string       `json:"question_name,omitempty"`
	QuestionText      string       `json:"question_text"`
	QuestionType      string       `json:"question_type"`
	PointsPossible    float64      `json:"points_possible,omitempty"`
	CorrectComments   string       `json:"correct_comments,omitempty"`
	IncorrectComments string       `json:"incorrect_comments,omitempty"`
	NeutralComments   string       `json:"neutral_comments,omitempty"`
	Position          int          `json:"position,omitempty"`
	Answers           []QuizAnswer `json:"answers,omitempty"`
}

// Create creates a new quiz question
func (s *QuizQuestionsService) Create(ctx context.Context, courseID, quizID int64, params *CreateQuizQuestionParams) (*QuizQuestion, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/%d/questions", courseID, quizID)

	body := map[string]interface{}{
		"question": make(map[string]interface{}),
	}

	questionData, ok := body["question"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("internal error: invalid question data structure")
	}

	if params.QuestionName != "" {
		questionData["question_name"] = params.QuestionName
	}

	if params.QuestionText != "" {
		questionData["question_text"] = params.QuestionText
	}

	if params.QuestionType != "" {
		questionData["question_type"] = params.QuestionType
	}

	if params.PointsPossible > 0 {
		questionData["points_possible"] = params.PointsPossible
	}

	if params.CorrectComments != "" {
		questionData["correct_comments"] = params.CorrectComments
	}

	if params.IncorrectComments != "" {
		questionData["incorrect_comments"] = params.IncorrectComments
	}

	if params.NeutralComments != "" {
		questionData["neutral_comments"] = params.NeutralComments
	}

	if params.Position > 0 {
		questionData["position"] = params.Position
	}

	if len(params.Answers) > 0 {
		questionData["answers"] = params.Answers
	}

	var question QuizQuestion
	if err := s.client.PostJSON(ctx, path, body, &question); err != nil {
		return nil, err
	}

	return &question, nil
}

// UpdateQuizQuestionParams holds parameters for updating a quiz question
type UpdateQuizQuestionParams struct {
	QuestionName      *string       `json:"question_name,omitempty"`
	QuestionText      *string       `json:"question_text,omitempty"`
	QuestionType      *string       `json:"question_type,omitempty"`
	PointsPossible    *float64      `json:"points_possible,omitempty"`
	CorrectComments   *string       `json:"correct_comments,omitempty"`
	IncorrectComments *string       `json:"incorrect_comments,omitempty"`
	NeutralComments   *string       `json:"neutral_comments,omitempty"`
	Position          *int          `json:"position,omitempty"`
	Answers           *[]QuizAnswer `json:"answers,omitempty"`
}

// Update updates an existing quiz question
func (s *QuizQuestionsService) Update(ctx context.Context, courseID, quizID, questionID int64, params *UpdateQuizQuestionParams) (*QuizQuestion, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/%d/questions/%d", courseID, quizID, questionID)

	body := map[string]interface{}{
		"question": make(map[string]interface{}),
	}

	questionData, ok := body["question"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("internal error: invalid question data structure")
	}

	if params.QuestionName != nil {
		questionData["question_name"] = *params.QuestionName
	}

	if params.QuestionText != nil {
		questionData["question_text"] = *params.QuestionText
	}

	if params.QuestionType != nil {
		questionData["question_type"] = *params.QuestionType
	}

	if params.PointsPossible != nil {
		questionData["points_possible"] = *params.PointsPossible
	}

	if params.CorrectComments != nil {
		questionData["correct_comments"] = *params.CorrectComments
	}

	if params.IncorrectComments != nil {
		questionData["incorrect_comments"] = *params.IncorrectComments
	}

	if params.NeutralComments != nil {
		questionData["neutral_comments"] = *params.NeutralComments
	}

	if params.Position != nil {
		questionData["position"] = *params.Position
	}

	if params.Answers != nil {
		questionData["answers"] = *params.Answers
	}

	var question QuizQuestion
	if err := s.client.PutJSON(ctx, path, body, &question); err != nil {
		return nil, err
	}

	return &question, nil
}

// Delete deletes a quiz question
func (s *QuizQuestionsService) Delete(ctx context.Context, courseID, quizID, questionID int64) error {
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/%d/questions/%d", courseID, quizID, questionID)

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
