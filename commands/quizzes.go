package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var (
	// Common flags
	quizzesCourseID int64
	quizzesQuizID   int64

	// List flags
	quizzesSearchTerm string

	// Create/Update flags
	quizzesTitle                string
	quizzesDescription          string
	quizzesQuizType             string
	quizzesAssignmentGroupID    int64
	quizzesTimeLimit            int
	quizzesShuffleAnswers       bool
	quizzesHideResults          string
	quizzesShowCorrectAnswers   bool
	quizzesScoringPolicy        string
	quizzesAllowedAttempts      int
	quizzesOneQuestionAtATime   bool
	quizzesCantGoBack           bool
	quizzesAccessCode           string
	quizzesIPFilter             string
	quizzesDueAt                string
	quizzesLockAt               string
	quizzesUnlockAt             string
	quizzesPublished            bool
	quizzesAnonymousSubmissions bool

	// Question flags
	quizzesQuestionName      string
	quizzesQuestionText      string
	quizzesQuestionType      string
	quizzesPointsPossible    float64
	quizzesCorrectComments   string
	quizzesIncorrectComments string

	// Delete flags
	quizzesForce bool
)

// quizzesCmd represents the quizzes command group
var quizzesCmd = &cobra.Command{
	Use:     "quizzes",
	Aliases: []string{"quiz"},
	Short:   "Manage Canvas quizzes",
	Long: `Manage Canvas quizzes for courses.

Quizzes allow you to create assessments with various question types including
multiple choice, true/false, short answer, and more.

Examples:
  canvas quizzes list --course-id 123
  canvas quizzes get 456 --course-id 123
  canvas quizzes create --course-id 123 --title "Midterm Exam" --quiz-type assignment`,
}

// quizzesListCmd represents the quizzes list command
var quizzesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List quizzes in a course",
	Long: `List all quizzes in a course.

Examples:
  canvas quizzes list --course-id 123
  canvas quizzes list --course-id 123 --search "midterm"`,
	RunE: runQuizzesList,
}

// quizzesGetCmd represents the quizzes get command
var quizzesGetCmd = &cobra.Command{
	Use:   "get <quiz-id>",
	Short: "Get quiz details",
	Long: `Get details of a specific quiz.

Examples:
  canvas quizzes get 456 --course-id 123`,
	Args: ExactArgsWithUsage(1, "quiz-id"),
	RunE: runQuizzesGet,
}

// quizzesCreateCmd represents the quizzes create command
var quizzesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new quiz",
	Long: `Create a new quiz in a course.

Examples:
  canvas quizzes create --course-id 123 --title "Midterm Exam" --quiz-type assignment
  canvas quizzes create --course-id 123 --title "Practice Quiz" --quiz-type practice_quiz --time-limit 30
  canvas quizzes create --course-id 123 --title "Survey" --quiz-type survey --anonymous`,
	RunE: runQuizzesCreate,
}

// quizzesUpdateCmd represents the quizzes update command
var quizzesUpdateCmd = &cobra.Command{
	Use:   "update <quiz-id>",
	Short: "Update a quiz",
	Long: `Update an existing quiz.

Examples:
  canvas quizzes update 456 --course-id 123 --title "Updated Title"
  canvas quizzes update 456 --course-id 123 --time-limit 60
  canvas quizzes update 456 --course-id 123 --published`,
	Args: ExactArgsWithUsage(1, "quiz-id"),
	RunE: runQuizzesUpdate,
}

// quizzesDeleteCmd represents the quizzes delete command
var quizzesDeleteCmd = &cobra.Command{
	Use:   "delete <quiz-id>",
	Short: "Delete a quiz",
	Long: `Delete a quiz.

Examples:
  canvas quizzes delete 456 --course-id 123
  canvas quizzes delete 456 --course-id 123 --force`,
	Args: ExactArgsWithUsage(1, "quiz-id"),
	RunE: runQuizzesDelete,
}

// Quiz Questions Commands

// quizzesQuestionsCmd represents the quizzes questions command group
var quizzesQuestionsCmd = &cobra.Command{
	Use:     "questions",
	Aliases: []string{"q"},
	Short:   "Manage quiz questions",
	Long:    `Manage questions within a quiz.`,
}

// quizzesQuestionsListCmd represents the quizzes questions list command
var quizzesQuestionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List questions in a quiz",
	Long: `List all questions in a quiz.

Examples:
  canvas quizzes questions list --course-id 123 --quiz-id 456`,
	RunE: runQuizzesQuestionsList,
}

// quizzesQuestionsGetCmd represents the quizzes questions get command
var quizzesQuestionsGetCmd = &cobra.Command{
	Use:   "get <question-id>",
	Short: "Get question details",
	Long: `Get details of a specific question.

Examples:
  canvas quizzes questions get 789 --course-id 123 --quiz-id 456`,
	Args: ExactArgsWithUsage(1, "question-id"),
	RunE: runQuizzesQuestionsGet,
}

// quizzesQuestionsCreateCmd represents the quizzes questions create command
var quizzesQuestionsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new question",
	Long: `Create a new question in a quiz.

Examples:
  canvas quizzes questions create --course-id 123 --quiz-id 456 --text "What is 2+2?" --type multiple_choice_question --points 10`,
	RunE: runQuizzesQuestionsCreate,
}

// quizzesQuestionsDeleteCmd represents the quizzes questions delete command
var quizzesQuestionsDeleteCmd = &cobra.Command{
	Use:   "delete <question-id>",
	Short: "Delete a question",
	Long: `Delete a question from a quiz.

Examples:
  canvas quizzes questions delete 789 --course-id 123 --quiz-id 456 --force`,
	Args: ExactArgsWithUsage(1, "question-id"),
	RunE: runQuizzesQuestionsDelete,
}

// Quiz Submissions Commands

// quizzesSubmissionsCmd represents the quizzes submissions command group
var quizzesSubmissionsCmd = &cobra.Command{
	Use:     "submissions",
	Aliases: []string{"sub"},
	Short:   "Manage quiz submissions",
	Long:    `View and manage quiz submissions.`,
}

// quizzesSubmissionsListCmd represents the quizzes submissions list command
var quizzesSubmissionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List quiz submissions",
	Long: `List all submissions for a quiz.

Examples:
  canvas quizzes submissions list --course-id 123 --quiz-id 456`,
	RunE: runQuizzesSubmissionsList,
}

// quizzesSubmissionsGetCmd represents the quizzes submissions get command
var quizzesSubmissionsGetCmd = &cobra.Command{
	Use:   "get <submission-id>",
	Short: "Get submission details",
	Long: `Get details of a specific quiz submission.

Examples:
  canvas quizzes submissions get 789 --course-id 123 --quiz-id 456`,
	Args: ExactArgsWithUsage(1, "submission-id"),
	RunE: runQuizzesSubmissionsGet,
}

func init() {
	rootCmd.AddCommand(quizzesCmd)
	quizzesCmd.AddCommand(quizzesListCmd)
	quizzesCmd.AddCommand(quizzesGetCmd)
	quizzesCmd.AddCommand(quizzesCreateCmd)
	quizzesCmd.AddCommand(quizzesUpdateCmd)
	quizzesCmd.AddCommand(quizzesDeleteCmd)
	quizzesCmd.AddCommand(quizzesQuestionsCmd)
	quizzesCmd.AddCommand(quizzesSubmissionsCmd)

	// Questions subcommands
	quizzesQuestionsCmd.AddCommand(quizzesQuestionsListCmd)
	quizzesQuestionsCmd.AddCommand(quizzesQuestionsGetCmd)
	quizzesQuestionsCmd.AddCommand(quizzesQuestionsCreateCmd)
	quizzesQuestionsCmd.AddCommand(quizzesQuestionsDeleteCmd)

	// Submissions subcommands
	quizzesSubmissionsCmd.AddCommand(quizzesSubmissionsListCmd)
	quizzesSubmissionsCmd.AddCommand(quizzesSubmissionsGetCmd)

	// List flags
	quizzesListCmd.Flags().Int64Var(&quizzesCourseID, "course-id", 0, "Course ID (required)")
	quizzesListCmd.MarkFlagRequired("course-id")
	quizzesListCmd.Flags().StringVar(&quizzesSearchTerm, "search", "", "Search term")

	// Get flags
	quizzesGetCmd.Flags().Int64Var(&quizzesCourseID, "course-id", 0, "Course ID (required)")
	quizzesGetCmd.MarkFlagRequired("course-id")

	// Create flags
	quizzesCreateCmd.Flags().Int64Var(&quizzesCourseID, "course-id", 0, "Course ID (required)")
	quizzesCreateCmd.MarkFlagRequired("course-id")
	quizzesCreateCmd.Flags().StringVar(&quizzesTitle, "title", "", "Quiz title (required)")
	quizzesCreateCmd.MarkFlagRequired("title")
	quizzesCreateCmd.Flags().StringVar(&quizzesDescription, "description", "", "Quiz description")
	quizzesCreateCmd.Flags().StringVar(&quizzesQuizType, "quiz-type", "assignment", "Quiz type: practice_quiz, assignment, graded_survey, survey")
	quizzesCreateCmd.Flags().Int64Var(&quizzesAssignmentGroupID, "assignment-group-id", 0, "Assignment group ID")
	quizzesCreateCmd.Flags().IntVar(&quizzesTimeLimit, "time-limit", 0, "Time limit in minutes")
	quizzesCreateCmd.Flags().BoolVar(&quizzesShuffleAnswers, "shuffle-answers", false, "Shuffle answer choices")
	quizzesCreateCmd.Flags().StringVar(&quizzesHideResults, "hide-results", "", "When to hide results: always, until_after_last_attempt")
	quizzesCreateCmd.Flags().BoolVar(&quizzesShowCorrectAnswers, "show-correct", false, "Show correct answers")
	quizzesCreateCmd.Flags().StringVar(&quizzesScoringPolicy, "scoring-policy", "", "Scoring policy: keep_highest, keep_latest")
	quizzesCreateCmd.Flags().IntVar(&quizzesAllowedAttempts, "attempts", 1, "Number of allowed attempts (-1 = unlimited)")
	quizzesCreateCmd.Flags().BoolVar(&quizzesOneQuestionAtATime, "one-at-a-time", false, "Show one question at a time")
	quizzesCreateCmd.Flags().BoolVar(&quizzesCantGoBack, "cant-go-back", false, "Prevent going back to previous questions")
	quizzesCreateCmd.Flags().StringVar(&quizzesAccessCode, "access-code", "", "Quiz access code")
	quizzesCreateCmd.Flags().StringVar(&quizzesIPFilter, "ip-filter", "", "IP address filter")
	quizzesCreateCmd.Flags().StringVar(&quizzesDueAt, "due-at", "", "Due date (ISO 8601)")
	quizzesCreateCmd.Flags().StringVar(&quizzesLockAt, "lock-at", "", "Lock date (ISO 8601)")
	quizzesCreateCmd.Flags().StringVar(&quizzesUnlockAt, "unlock-at", "", "Unlock date (ISO 8601)")
	quizzesCreateCmd.Flags().BoolVar(&quizzesPublished, "published", false, "Publish immediately")
	quizzesCreateCmd.Flags().BoolVar(&quizzesAnonymousSubmissions, "anonymous", false, "Anonymous submissions")

	// Update flags
	quizzesUpdateCmd.Flags().Int64Var(&quizzesCourseID, "course-id", 0, "Course ID (required)")
	quizzesUpdateCmd.MarkFlagRequired("course-id")
	quizzesUpdateCmd.Flags().StringVar(&quizzesTitle, "title", "", "Quiz title")
	quizzesUpdateCmd.Flags().StringVar(&quizzesDescription, "description", "", "Quiz description")
	quizzesUpdateCmd.Flags().StringVar(&quizzesQuizType, "quiz-type", "", "Quiz type")
	quizzesUpdateCmd.Flags().Int64Var(&quizzesAssignmentGroupID, "assignment-group-id", 0, "Assignment group ID")
	quizzesUpdateCmd.Flags().IntVar(&quizzesTimeLimit, "time-limit", 0, "Time limit in minutes")
	quizzesUpdateCmd.Flags().BoolVar(&quizzesShuffleAnswers, "shuffle-answers", false, "Shuffle answer choices")
	quizzesUpdateCmd.Flags().StringVar(&quizzesHideResults, "hide-results", "", "When to hide results")
	quizzesUpdateCmd.Flags().BoolVar(&quizzesShowCorrectAnswers, "show-correct", false, "Show correct answers")
	quizzesUpdateCmd.Flags().StringVar(&quizzesScoringPolicy, "scoring-policy", "", "Scoring policy")
	quizzesUpdateCmd.Flags().IntVar(&quizzesAllowedAttempts, "attempts", 0, "Number of allowed attempts")
	quizzesUpdateCmd.Flags().BoolVar(&quizzesOneQuestionAtATime, "one-at-a-time", false, "Show one question at a time")
	quizzesUpdateCmd.Flags().BoolVar(&quizzesCantGoBack, "cant-go-back", false, "Prevent going back")
	quizzesUpdateCmd.Flags().StringVar(&quizzesAccessCode, "access-code", "", "Quiz access code")
	quizzesUpdateCmd.Flags().StringVar(&quizzesIPFilter, "ip-filter", "", "IP address filter")
	quizzesUpdateCmd.Flags().StringVar(&quizzesDueAt, "due-at", "", "Due date (ISO 8601)")
	quizzesUpdateCmd.Flags().StringVar(&quizzesLockAt, "lock-at", "", "Lock date (ISO 8601)")
	quizzesUpdateCmd.Flags().StringVar(&quizzesUnlockAt, "unlock-at", "", "Unlock date (ISO 8601)")
	quizzesUpdateCmd.Flags().BoolVar(&quizzesPublished, "published", false, "Publish quiz")
	quizzesUpdateCmd.Flags().BoolVar(&quizzesAnonymousSubmissions, "anonymous", false, "Anonymous submissions")

	// Delete flags
	quizzesDeleteCmd.Flags().Int64Var(&quizzesCourseID, "course-id", 0, "Course ID (required)")
	quizzesDeleteCmd.MarkFlagRequired("course-id")
	quizzesDeleteCmd.Flags().BoolVar(&quizzesForce, "force", false, "Skip confirmation prompt")

	// Questions List flags
	quizzesQuestionsListCmd.Flags().Int64Var(&quizzesCourseID, "course-id", 0, "Course ID (required)")
	quizzesQuestionsListCmd.MarkFlagRequired("course-id")
	quizzesQuestionsListCmd.Flags().Int64Var(&quizzesQuizID, "quiz-id", 0, "Quiz ID (required)")
	quizzesQuestionsListCmd.MarkFlagRequired("quiz-id")

	// Questions Get flags
	quizzesQuestionsGetCmd.Flags().Int64Var(&quizzesCourseID, "course-id", 0, "Course ID (required)")
	quizzesQuestionsGetCmd.MarkFlagRequired("course-id")
	quizzesQuestionsGetCmd.Flags().Int64Var(&quizzesQuizID, "quiz-id", 0, "Quiz ID (required)")
	quizzesQuestionsGetCmd.MarkFlagRequired("quiz-id")

	// Questions Create flags
	quizzesQuestionsCreateCmd.Flags().Int64Var(&quizzesCourseID, "course-id", 0, "Course ID (required)")
	quizzesQuestionsCreateCmd.MarkFlagRequired("course-id")
	quizzesQuestionsCreateCmd.Flags().Int64Var(&quizzesQuizID, "quiz-id", 0, "Quiz ID (required)")
	quizzesQuestionsCreateCmd.MarkFlagRequired("quiz-id")
	quizzesQuestionsCreateCmd.Flags().StringVar(&quizzesQuestionName, "name", "", "Question name")
	quizzesQuestionsCreateCmd.Flags().StringVar(&quizzesQuestionText, "text", "", "Question text (required)")
	quizzesQuestionsCreateCmd.MarkFlagRequired("text")
	quizzesQuestionsCreateCmd.Flags().StringVar(&quizzesQuestionType, "type", "multiple_choice_question", "Question type")
	quizzesQuestionsCreateCmd.Flags().Float64Var(&quizzesPointsPossible, "points", 0, "Points possible")
	quizzesQuestionsCreateCmd.Flags().StringVar(&quizzesCorrectComments, "correct-comments", "", "Comments for correct answer")
	quizzesQuestionsCreateCmd.Flags().StringVar(&quizzesIncorrectComments, "incorrect-comments", "", "Comments for incorrect answer")

	// Questions Delete flags
	quizzesQuestionsDeleteCmd.Flags().Int64Var(&quizzesCourseID, "course-id", 0, "Course ID (required)")
	quizzesQuestionsDeleteCmd.MarkFlagRequired("course-id")
	quizzesQuestionsDeleteCmd.Flags().Int64Var(&quizzesQuizID, "quiz-id", 0, "Quiz ID (required)")
	quizzesQuestionsDeleteCmd.MarkFlagRequired("quiz-id")
	quizzesQuestionsDeleteCmd.Flags().BoolVar(&quizzesForce, "force", false, "Skip confirmation prompt")

	// Submissions List flags
	quizzesSubmissionsListCmd.Flags().Int64Var(&quizzesCourseID, "course-id", 0, "Course ID (required)")
	quizzesSubmissionsListCmd.MarkFlagRequired("course-id")
	quizzesSubmissionsListCmd.Flags().Int64Var(&quizzesQuizID, "quiz-id", 0, "Quiz ID (required)")
	quizzesSubmissionsListCmd.MarkFlagRequired("quiz-id")

	// Submissions Get flags
	quizzesSubmissionsGetCmd.Flags().Int64Var(&quizzesCourseID, "course-id", 0, "Course ID (required)")
	quizzesSubmissionsGetCmd.MarkFlagRequired("course-id")
	quizzesSubmissionsGetCmd.Flags().Int64Var(&quizzesQuizID, "quiz-id", 0, "Quiz ID (required)")
	quizzesSubmissionsGetCmd.MarkFlagRequired("quiz-id")
}

func runQuizzesList(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewQuizzesService(client)

	opts := &api.ListQuizzesOptions{}
	if quizzesSearchTerm != "" {
		opts.SearchTerm = quizzesSearchTerm
	}

	ctx := context.Background()
	quizzes, err := service.List(ctx, quizzesCourseID, opts)
	if err != nil {
		return fmt.Errorf("failed to list quizzes: %w", err)
	}

	if len(quizzes) == 0 {
		fmt.Printf("No quizzes found in course %d\n", quizzesCourseID)
		return nil
	}

	printVerbose("Found %d quizzes in course %d:\n\n", len(quizzes), quizzesCourseID)
	return formatOutput(quizzes, nil)
}

func runQuizzesGet(cmd *cobra.Command, args []string) error {
	quizID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid quiz ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewQuizzesService(client)

	ctx := context.Background()
	quiz, err := service.Get(ctx, quizzesCourseID, quizID)
	if err != nil {
		return fmt.Errorf("failed to get quiz: %w", err)
	}

	return formatOutput(quiz, nil)
}

func runQuizzesCreate(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewQuizzesService(client)

	params := &api.CreateQuizParams{
		Title:                quizzesTitle,
		Description:          quizzesDescription,
		QuizType:             quizzesQuizType,
		AssignmentGroupID:    quizzesAssignmentGroupID,
		TimeLimit:            quizzesTimeLimit,
		ShuffleAnswers:       quizzesShuffleAnswers,
		HideResults:          quizzesHideResults,
		ShowCorrectAnswers:   quizzesShowCorrectAnswers,
		ScoringPolicy:        quizzesScoringPolicy,
		AllowedAttempts:      quizzesAllowedAttempts,
		OneQuestionAtATime:   quizzesOneQuestionAtATime,
		CantGoBack:           quizzesCantGoBack,
		AccessCode:           quizzesAccessCode,
		IPFilter:             quizzesIPFilter,
		DueAt:                quizzesDueAt,
		LockAt:               quizzesLockAt,
		UnlockAt:             quizzesUnlockAt,
		Published:            quizzesPublished,
		AnonymousSubmissions: quizzesAnonymousSubmissions,
	}

	ctx := context.Background()
	quiz, err := service.Create(ctx, quizzesCourseID, params)
	if err != nil {
		return fmt.Errorf("failed to create quiz: %w", err)
	}

	fmt.Printf("Quiz created successfully (ID: %d)\n", quiz.ID)
	return formatOutput(quiz, nil)
}

func runQuizzesUpdate(cmd *cobra.Command, args []string) error {
	quizID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid quiz ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewQuizzesService(client)

	params := &api.UpdateQuizParams{}

	if cmd.Flags().Changed("title") {
		params.Title = &quizzesTitle
	}
	if cmd.Flags().Changed("description") {
		params.Description = &quizzesDescription
	}
	if cmd.Flags().Changed("quiz-type") {
		params.QuizType = &quizzesQuizType
	}
	if cmd.Flags().Changed("assignment-group-id") {
		params.AssignmentGroupID = &quizzesAssignmentGroupID
	}
	if cmd.Flags().Changed("time-limit") {
		params.TimeLimit = &quizzesTimeLimit
	}
	if cmd.Flags().Changed("shuffle-answers") {
		params.ShuffleAnswers = &quizzesShuffleAnswers
	}
	if cmd.Flags().Changed("hide-results") {
		params.HideResults = &quizzesHideResults
	}
	if cmd.Flags().Changed("show-correct") {
		params.ShowCorrectAnswers = &quizzesShowCorrectAnswers
	}
	if cmd.Flags().Changed("scoring-policy") {
		params.ScoringPolicy = &quizzesScoringPolicy
	}
	if cmd.Flags().Changed("attempts") {
		params.AllowedAttempts = &quizzesAllowedAttempts
	}
	if cmd.Flags().Changed("one-at-a-time") {
		params.OneQuestionAtATime = &quizzesOneQuestionAtATime
	}
	if cmd.Flags().Changed("cant-go-back") {
		params.CantGoBack = &quizzesCantGoBack
	}
	if cmd.Flags().Changed("access-code") {
		params.AccessCode = &quizzesAccessCode
	}
	if cmd.Flags().Changed("ip-filter") {
		params.IPFilter = &quizzesIPFilter
	}
	if cmd.Flags().Changed("due-at") {
		params.DueAt = &quizzesDueAt
	}
	if cmd.Flags().Changed("lock-at") {
		params.LockAt = &quizzesLockAt
	}
	if cmd.Flags().Changed("unlock-at") {
		params.UnlockAt = &quizzesUnlockAt
	}
	if cmd.Flags().Changed("published") {
		params.Published = &quizzesPublished
	}
	if cmd.Flags().Changed("anonymous") {
		params.AnonymousSubmissions = &quizzesAnonymousSubmissions
	}

	ctx := context.Background()
	quiz, err := service.Update(ctx, quizzesCourseID, quizID, params)
	if err != nil {
		return fmt.Errorf("failed to update quiz: %w", err)
	}

	fmt.Printf("Quiz updated successfully (ID: %d)\n", quiz.ID)
	return formatOutput(quiz, nil)
}

func runQuizzesDelete(cmd *cobra.Command, args []string) error {
	quizID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid quiz ID: %w", err)
	}

	if !quizzesForce {
		fmt.Printf("WARNING: This will delete quiz %d.\n", quizID)
		fmt.Print("Type 'yes' to confirm: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Delete cancelled")
			return nil
		}
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewQuizzesService(client)

	ctx := context.Background()
	quiz, err := service.Delete(ctx, quizzesCourseID, quizID)
	if err != nil {
		return fmt.Errorf("failed to delete quiz: %w", err)
	}

	fmt.Printf("Quiz %d deleted\n", quiz.ID)
	return nil
}

func runQuizzesQuestionsList(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewQuizQuestionsService(client)

	ctx := context.Background()
	questions, err := service.List(ctx, quizzesCourseID, quizzesQuizID, nil)
	if err != nil {
		return fmt.Errorf("failed to list questions: %w", err)
	}

	if len(questions) == 0 {
		fmt.Printf("No questions found in quiz %d\n", quizzesQuizID)
		return nil
	}

	printVerbose("Found %d questions in quiz %d:\n\n", len(questions), quizzesQuizID)
	return formatOutput(questions, nil)
}

func runQuizzesQuestionsGet(cmd *cobra.Command, args []string) error {
	questionID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid question ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewQuizQuestionsService(client)

	ctx := context.Background()
	question, err := service.Get(ctx, quizzesCourseID, quizzesQuizID, questionID)
	if err != nil {
		return fmt.Errorf("failed to get question: %w", err)
	}

	return formatOutput(question, nil)
}

func runQuizzesQuestionsCreate(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewQuizQuestionsService(client)

	params := &api.CreateQuizQuestionParams{
		QuestionName:      quizzesQuestionName,
		QuestionText:      quizzesQuestionText,
		QuestionType:      quizzesQuestionType,
		PointsPossible:    quizzesPointsPossible,
		CorrectComments:   quizzesCorrectComments,
		IncorrectComments: quizzesIncorrectComments,
	}

	ctx := context.Background()
	question, err := service.Create(ctx, quizzesCourseID, quizzesQuizID, params)
	if err != nil {
		return fmt.Errorf("failed to create question: %w", err)
	}

	fmt.Printf("Question created successfully (ID: %d)\n", question.ID)
	return formatOutput(question, nil)
}

func runQuizzesQuestionsDelete(cmd *cobra.Command, args []string) error {
	questionID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid question ID: %w", err)
	}

	if !quizzesForce {
		fmt.Printf("WARNING: This will delete question %d from quiz %d.\n", questionID, quizzesQuizID)
		fmt.Print("Type 'yes' to confirm: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Delete cancelled")
			return nil
		}
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewQuizQuestionsService(client)

	ctx := context.Background()
	if err := service.Delete(ctx, quizzesCourseID, quizzesQuizID, questionID); err != nil {
		return fmt.Errorf("failed to delete question: %w", err)
	}

	fmt.Printf("Question %d deleted\n", questionID)
	return nil
}

func runQuizzesSubmissionsList(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewQuizSubmissionsService(client)

	ctx := context.Background()
	submissions, err := service.List(ctx, quizzesCourseID, quizzesQuizID, nil)
	if err != nil {
		return fmt.Errorf("failed to list submissions: %w", err)
	}

	if len(submissions) == 0 {
		fmt.Printf("No submissions found for quiz %d\n", quizzesQuizID)
		return nil
	}

	printVerbose("Found %d submissions for quiz %d:\n\n", len(submissions), quizzesQuizID)
	return formatOutput(submissions, nil)
}

func runQuizzesSubmissionsGet(cmd *cobra.Command, args []string) error {
	submissionID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid submission ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewQuizSubmissionsService(client)

	ctx := context.Background()
	submission, err := service.Get(ctx, quizzesCourseID, quizzesQuizID, submissionID, nil)
	if err != nil {
		return fmt.Errorf("failed to get submission: %w", err)
	}

	return formatOutput(submission, nil)
}
