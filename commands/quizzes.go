package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
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

// quizzesQuestionsCmd represents the quizzes questions command group
var quizzesQuestionsCmd = &cobra.Command{
	Use:     "questions",
	Aliases: []string{"q"},
	Short:   "Manage quiz questions",
	Long:    `Manage questions within a quiz.`,
}

// quizzesSubmissionsCmd represents the quizzes submissions command group
var quizzesSubmissionsCmd = &cobra.Command{
	Use:     "submissions",
	Aliases: []string{"sub"},
	Short:   "Manage quiz submissions",
	Long:    `View and manage quiz submissions.`,
}

func init() {
	rootCmd.AddCommand(quizzesCmd)
	quizzesCmd.AddCommand(newQuizzesListCmd())
	quizzesCmd.AddCommand(newQuizzesGetCmd())
	quizzesCmd.AddCommand(newQuizzesCreateCmd())
	quizzesCmd.AddCommand(newQuizzesUpdateCmd())
	quizzesCmd.AddCommand(newQuizzesDeleteCmd())
	quizzesCmd.AddCommand(quizzesQuestionsCmd)
	quizzesCmd.AddCommand(quizzesSubmissionsCmd)

	// Questions subcommands
	quizzesQuestionsCmd.AddCommand(newQuizzesQuestionsListCmd())
	quizzesQuestionsCmd.AddCommand(newQuizzesQuestionsGetCmd())
	quizzesQuestionsCmd.AddCommand(newQuizzesQuestionsCreateCmd())
	quizzesQuestionsCmd.AddCommand(newQuizzesQuestionsDeleteCmd())

	// Submissions subcommands
	quizzesSubmissionsCmd.AddCommand(newQuizzesSubmissionsListCmd())
	quizzesSubmissionsCmd.AddCommand(newQuizzesSubmissionsGetCmd())
}

func newQuizzesListCmd() *cobra.Command {
	opts := &options.QuizzesListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List quizzes in a course",
		Long: `List all quizzes in a course.

Examples:
  canvas quizzes list --course-id 123
  canvas quizzes list --course-id 123 --search "midterm"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runQuizzesList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.SearchTerm, "search", "", "Search term")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newQuizzesGetCmd() *cobra.Command {
	opts := &options.QuizzesGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <quiz-id>",
		Short: "Get quiz details",
		Long: `Get details of a specific quiz.

Examples:
  canvas quizzes get 456 --course-id 123`,
		Args: ExactArgsWithUsage(1, "quiz-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			quizID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid quiz ID: %s", args[0])
			}
			opts.QuizID = quizID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runQuizzesGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newQuizzesCreateCmd() *cobra.Command {
	opts := &options.QuizzesCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new quiz",
		Long: `Create a new quiz in a course.

Examples:
  canvas quizzes create --course-id 123 --title "Midterm Exam" --quiz-type assignment
  canvas quizzes create --course-id 123 --title "Practice Quiz" --quiz-type practice_quiz --time-limit 30
  canvas quizzes create --course-id 123 --title "Survey" --quiz-type survey --anonymous`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runQuizzesCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.Title, "title", "", "Quiz title (required)")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Quiz description")
	cmd.Flags().StringVar(&opts.QuizType, "quiz-type", "assignment", "Quiz type: practice_quiz, assignment, graded_survey, survey")
	cmd.Flags().Int64Var(&opts.AssignmentGroupID, "assignment-group-id", 0, "Assignment group ID")
	cmd.Flags().IntVar(&opts.TimeLimit, "time-limit", 0, "Time limit in minutes")
	cmd.Flags().BoolVar(&opts.ShuffleAnswers, "shuffle-answers", false, "Shuffle answer choices")
	cmd.Flags().StringVar(&opts.HideResults, "hide-results", "", "When to hide results: always, until_after_last_attempt")
	cmd.Flags().BoolVar(&opts.ShowCorrectAnswers, "show-correct", false, "Show correct answers")
	cmd.Flags().StringVar(&opts.ScoringPolicy, "scoring-policy", "", "Scoring policy: keep_highest, keep_latest")
	cmd.Flags().IntVar(&opts.AllowedAttempts, "attempts", 1, "Number of allowed attempts (-1 = unlimited)")
	cmd.Flags().BoolVar(&opts.OneQuestionAtATime, "one-at-a-time", false, "Show one question at a time")
	cmd.Flags().BoolVar(&opts.CantGoBack, "cant-go-back", false, "Prevent going back to previous questions")
	cmd.Flags().StringVar(&opts.AccessCode, "access-code", "", "Quiz access code")
	cmd.Flags().StringVar(&opts.IPFilter, "ip-filter", "", "IP address filter")
	cmd.Flags().StringVar(&opts.DueAt, "due-at", "", "Due date (ISO 8601)")
	cmd.Flags().StringVar(&opts.LockAt, "lock-at", "", "Lock date (ISO 8601)")
	cmd.Flags().StringVar(&opts.UnlockAt, "unlock-at", "", "Unlock date (ISO 8601)")
	cmd.Flags().BoolVar(&opts.Published, "published", false, "Publish immediately")
	cmd.Flags().BoolVar(&opts.AnonymousSubmissions, "anonymous", false, "Anonymous submissions")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("title")

	return cmd
}

func newQuizzesUpdateCmd() *cobra.Command {
	opts := &options.QuizzesUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <quiz-id>",
		Short: "Update a quiz",
		Long: `Update an existing quiz.

Examples:
  canvas quizzes update 456 --course-id 123 --title "Updated Title"
  canvas quizzes update 456 --course-id 123 --time-limit 60
  canvas quizzes update 456 --course-id 123 --published`,
		Args: ExactArgsWithUsage(1, "quiz-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			quizID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid quiz ID: %s", args[0])
			}
			opts.QuizID = quizID

			// Track which flags were set
			opts.TitleSet = cmd.Flags().Changed("title")
			opts.DescriptionSet = cmd.Flags().Changed("description")
			opts.QuizTypeSet = cmd.Flags().Changed("quiz-type")
			opts.AssignmentGroupIDSet = cmd.Flags().Changed("assignment-group-id")
			opts.TimeLimitSet = cmd.Flags().Changed("time-limit")
			opts.ShuffleAnswersSet = cmd.Flags().Changed("shuffle-answers")
			opts.HideResultsSet = cmd.Flags().Changed("hide-results")
			opts.ShowCorrectAnswersSet = cmd.Flags().Changed("show-correct")
			opts.ScoringPolicySet = cmd.Flags().Changed("scoring-policy")
			opts.AllowedAttemptsSet = cmd.Flags().Changed("attempts")
			opts.OneQuestionAtATimeSet = cmd.Flags().Changed("one-at-a-time")
			opts.CantGoBackSet = cmd.Flags().Changed("cant-go-back")
			opts.AccessCodeSet = cmd.Flags().Changed("access-code")
			opts.IPFilterSet = cmd.Flags().Changed("ip-filter")
			opts.DueAtSet = cmd.Flags().Changed("due-at")
			opts.LockAtSet = cmd.Flags().Changed("lock-at")
			opts.UnlockAtSet = cmd.Flags().Changed("unlock-at")
			opts.PublishedSet = cmd.Flags().Changed("published")
			opts.AnonymousSubmissionsSet = cmd.Flags().Changed("anonymous")

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runQuizzesUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.Title, "title", "", "Quiz title")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Quiz description")
	cmd.Flags().StringVar(&opts.QuizType, "quiz-type", "", "Quiz type")
	cmd.Flags().Int64Var(&opts.AssignmentGroupID, "assignment-group-id", 0, "Assignment group ID")
	cmd.Flags().IntVar(&opts.TimeLimit, "time-limit", 0, "Time limit in minutes")
	cmd.Flags().BoolVar(&opts.ShuffleAnswers, "shuffle-answers", false, "Shuffle answer choices")
	cmd.Flags().StringVar(&opts.HideResults, "hide-results", "", "When to hide results")
	cmd.Flags().BoolVar(&opts.ShowCorrectAnswers, "show-correct", false, "Show correct answers")
	cmd.Flags().StringVar(&opts.ScoringPolicy, "scoring-policy", "", "Scoring policy")
	cmd.Flags().IntVar(&opts.AllowedAttempts, "attempts", 0, "Number of allowed attempts")
	cmd.Flags().BoolVar(&opts.OneQuestionAtATime, "one-at-a-time", false, "Show one question at a time")
	cmd.Flags().BoolVar(&opts.CantGoBack, "cant-go-back", false, "Prevent going back")
	cmd.Flags().StringVar(&opts.AccessCode, "access-code", "", "Quiz access code")
	cmd.Flags().StringVar(&opts.IPFilter, "ip-filter", "", "IP address filter")
	cmd.Flags().StringVar(&opts.DueAt, "due-at", "", "Due date (ISO 8601)")
	cmd.Flags().StringVar(&opts.LockAt, "lock-at", "", "Lock date (ISO 8601)")
	cmd.Flags().StringVar(&opts.UnlockAt, "unlock-at", "", "Unlock date (ISO 8601)")
	cmd.Flags().BoolVar(&opts.Published, "published", false, "Publish quiz")
	cmd.Flags().BoolVar(&opts.AnonymousSubmissions, "anonymous", false, "Anonymous submissions")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newQuizzesDeleteCmd() *cobra.Command {
	opts := &options.QuizzesDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <quiz-id>",
		Short: "Delete a quiz",
		Long: `Delete a quiz.

Examples:
  canvas quizzes delete 456 --course-id 123
  canvas quizzes delete 456 --course-id 123 --force`,
		Args: ExactArgsWithUsage(1, "quiz-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			quizID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid quiz ID: %s", args[0])
			}
			opts.QuizID = quizID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runQuizzesDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().BoolVar(&opts.Force, "force", false, "Skip confirmation prompt")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newQuizzesQuestionsListCmd() *cobra.Command {
	opts := &options.QuizzesQuestionsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List questions in a quiz",
		Long: `List all questions in a quiz.

Examples:
  canvas quizzes questions list --course-id 123 --quiz-id 456`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runQuizzesQuestionsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.QuizID, "quiz-id", 0, "Quiz ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("quiz-id")

	return cmd
}

func newQuizzesQuestionsGetCmd() *cobra.Command {
	opts := &options.QuizzesQuestionsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <question-id>",
		Short: "Get question details",
		Long: `Get details of a specific question.

Examples:
  canvas quizzes questions get 789 --course-id 123 --quiz-id 456`,
		Args: ExactArgsWithUsage(1, "question-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			questionID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid question ID: %s", args[0])
			}
			opts.QuestionID = questionID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runQuizzesQuestionsGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.QuizID, "quiz-id", 0, "Quiz ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("quiz-id")

	return cmd
}

func newQuizzesQuestionsCreateCmd() *cobra.Command {
	opts := &options.QuizzesQuestionsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new question",
		Long: `Create a new question in a quiz.

Examples:
  canvas quizzes questions create --course-id 123 --quiz-id 456 --text "What is 2+2?" --type multiple_choice_question --points 10`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runQuizzesQuestionsCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.QuizID, "quiz-id", 0, "Quiz ID (required)")
	cmd.Flags().StringVar(&opts.QuestionName, "name", "", "Question name")
	cmd.Flags().StringVar(&opts.QuestionText, "text", "", "Question text (required)")
	cmd.Flags().StringVar(&opts.QuestionType, "type", "multiple_choice_question", "Question type")
	cmd.Flags().Float64Var(&opts.PointsPossible, "points", 0, "Points possible")
	cmd.Flags().StringVar(&opts.CorrectComments, "correct-comments", "", "Comments for correct answer")
	cmd.Flags().StringVar(&opts.IncorrectComments, "incorrect-comments", "", "Comments for incorrect answer")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("quiz-id")
	cmd.MarkFlagRequired("text")

	return cmd
}

func newQuizzesQuestionsDeleteCmd() *cobra.Command {
	opts := &options.QuizzesQuestionsDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <question-id>",
		Short: "Delete a question",
		Long: `Delete a question from a quiz.

Examples:
  canvas quizzes questions delete 789 --course-id 123 --quiz-id 456 --force`,
		Args: ExactArgsWithUsage(1, "question-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			questionID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid question ID: %s", args[0])
			}
			opts.QuestionID = questionID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runQuizzesQuestionsDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.QuizID, "quiz-id", 0, "Quiz ID (required)")
	cmd.Flags().BoolVar(&opts.Force, "force", false, "Skip confirmation prompt")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("quiz-id")

	return cmd
}

func newQuizzesSubmissionsListCmd() *cobra.Command {
	opts := &options.QuizzesSubmissionsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List quiz submissions",
		Long: `List all submissions for a quiz.

Examples:
  canvas quizzes submissions list --course-id 123 --quiz-id 456`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runQuizzesSubmissionsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.QuizID, "quiz-id", 0, "Quiz ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("quiz-id")

	return cmd
}

func newQuizzesSubmissionsGetCmd() *cobra.Command {
	opts := &options.QuizzesSubmissionsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <submission-id>",
		Short: "Get submission details",
		Long: `Get details of a specific quiz submission.

Examples:
  canvas quizzes submissions get 789 --course-id 123 --quiz-id 456`,
		Args: ExactArgsWithUsage(1, "submission-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			submissionID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid submission ID: %s", args[0])
			}
			opts.SubmissionID = submissionID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runQuizzesSubmissionsGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.QuizID, "quiz-id", 0, "Quiz ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("quiz-id")

	return cmd
}

// Run functions

func runQuizzesList(ctx context.Context, client *api.Client, opts *options.QuizzesListOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "quizzes.list", map[string]interface{}{
		"course_id":   opts.CourseID,
		"search_term": opts.SearchTerm,
	})

	// Validate course ID exists
	if _, err := validateCourseID(client, opts.CourseID); err != nil {
		logger.LogCommandError(ctx, "quizzes.list", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return err
	}

	service := api.NewQuizzesService(client)

	apiOpts := &api.ListQuizzesOptions{}
	if opts.SearchTerm != "" {
		apiOpts.SearchTerm = opts.SearchTerm
	}

	quizzes, err := service.List(ctx, opts.CourseID, apiOpts)
	if err != nil {
		logger.LogCommandError(ctx, "quizzes.list", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to list quizzes: %w", err)
	}

	if err := formatEmptyOrOutput(quizzes, "No quizzes found"); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "quizzes.list", len(quizzes))
	return nil
}

func runQuizzesGet(ctx context.Context, client *api.Client, opts *options.QuizzesGetOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "quizzes.get", map[string]interface{}{
		"course_id": opts.CourseID,
		"quiz_id":   opts.QuizID,
	})

	service := api.NewQuizzesService(client)

	quiz, err := service.Get(ctx, opts.CourseID, opts.QuizID)
	if err != nil {
		logger.LogCommandError(ctx, "quizzes.get", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"quiz_id":   opts.QuizID,
		})
		return fmt.Errorf("failed to get quiz: %w", err)
	}

	if err := formatOutput(quiz, nil); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "quizzes.get", 1)
	return nil
}

func runQuizzesCreate(ctx context.Context, client *api.Client, opts *options.QuizzesCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "quizzes.create", map[string]interface{}{
		"course_id": opts.CourseID,
		"title":     opts.Title,
		"quiz_type": opts.QuizType,
	})

	service := api.NewQuizzesService(client)

	params := &api.CreateQuizParams{
		Title:                opts.Title,
		Description:          opts.Description,
		QuizType:             opts.QuizType,
		AssignmentGroupID:    opts.AssignmentGroupID,
		TimeLimit:            opts.TimeLimit,
		ShuffleAnswers:       opts.ShuffleAnswers,
		HideResults:          opts.HideResults,
		ShowCorrectAnswers:   opts.ShowCorrectAnswers,
		ScoringPolicy:        opts.ScoringPolicy,
		AllowedAttempts:      opts.AllowedAttempts,
		OneQuestionAtATime:   opts.OneQuestionAtATime,
		CantGoBack:           opts.CantGoBack,
		AccessCode:           opts.AccessCode,
		IPFilter:             opts.IPFilter,
		DueAt:                opts.DueAt,
		LockAt:               opts.LockAt,
		UnlockAt:             opts.UnlockAt,
		Published:            opts.Published,
		AnonymousSubmissions: opts.AnonymousSubmissions,
	}

	quiz, err := service.Create(ctx, opts.CourseID, params)
	if err != nil {
		logger.LogCommandError(ctx, "quizzes.create", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"title":     opts.Title,
		})
		return fmt.Errorf("failed to create quiz: %w", err)
	}

	printInfo("Quiz created successfully (ID: %d)\n", quiz.ID)
	if err := formatOutput(quiz, nil); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "quizzes.create", 1)
	return nil
}

func runQuizzesUpdate(ctx context.Context, client *api.Client, opts *options.QuizzesUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "quizzes.update", map[string]interface{}{
		"course_id": opts.CourseID,
		"quiz_id":   opts.QuizID,
	})

	service := api.NewQuizzesService(client)

	params := &api.UpdateQuizParams{}

	if opts.TitleSet {
		params.Title = &opts.Title
	}
	if opts.DescriptionSet {
		params.Description = &opts.Description
	}
	if opts.QuizTypeSet {
		params.QuizType = &opts.QuizType
	}
	if opts.AssignmentGroupIDSet {
		params.AssignmentGroupID = &opts.AssignmentGroupID
	}
	if opts.TimeLimitSet {
		params.TimeLimit = &opts.TimeLimit
	}
	if opts.ShuffleAnswersSet {
		params.ShuffleAnswers = &opts.ShuffleAnswers
	}
	if opts.HideResultsSet {
		params.HideResults = &opts.HideResults
	}
	if opts.ShowCorrectAnswersSet {
		params.ShowCorrectAnswers = &opts.ShowCorrectAnswers
	}
	if opts.ScoringPolicySet {
		params.ScoringPolicy = &opts.ScoringPolicy
	}
	if opts.AllowedAttemptsSet {
		params.AllowedAttempts = &opts.AllowedAttempts
	}
	if opts.OneQuestionAtATimeSet {
		params.OneQuestionAtATime = &opts.OneQuestionAtATime
	}
	if opts.CantGoBackSet {
		params.CantGoBack = &opts.CantGoBack
	}
	if opts.AccessCodeSet {
		params.AccessCode = &opts.AccessCode
	}
	if opts.IPFilterSet {
		params.IPFilter = &opts.IPFilter
	}
	if opts.DueAtSet {
		params.DueAt = &opts.DueAt
	}
	if opts.LockAtSet {
		params.LockAt = &opts.LockAt
	}
	if opts.UnlockAtSet {
		params.UnlockAt = &opts.UnlockAt
	}
	if opts.PublishedSet {
		params.Published = &opts.Published
	}
	if opts.AnonymousSubmissionsSet {
		params.AnonymousSubmissions = &opts.AnonymousSubmissions
	}

	quiz, err := service.Update(ctx, opts.CourseID, opts.QuizID, params)
	if err != nil {
		logger.LogCommandError(ctx, "quizzes.update", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"quiz_id":   opts.QuizID,
		})
		return fmt.Errorf("failed to update quiz: %w", err)
	}

	printInfo("Quiz updated successfully (ID: %d)\n", quiz.ID)
	if err := formatOutput(quiz, nil); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "quizzes.update", 1)
	return nil
}

func runQuizzesDelete(ctx context.Context, client *api.Client, opts *options.QuizzesDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "quizzes.delete", map[string]interface{}{
		"course_id": opts.CourseID,
		"quiz_id":   opts.QuizID,
	})

	if !opts.Force {
		fmt.Printf("WARNING: This will delete quiz %d.\n", opts.QuizID)
		fmt.Print("Type 'yes' to confirm: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Delete cancelled")
			return nil
		}
	}

	service := api.NewQuizzesService(client)

	quiz, err := service.Delete(ctx, opts.CourseID, opts.QuizID)
	if err != nil {
		logger.LogCommandError(ctx, "quizzes.delete", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"quiz_id":   opts.QuizID,
		})
		return fmt.Errorf("failed to delete quiz: %w", err)
	}

	printInfo("Quiz %d deleted\n", quiz.ID)

	logger.LogCommandComplete(ctx, "quizzes.delete", 1)
	return nil
}

func runQuizzesQuestionsList(ctx context.Context, client *api.Client, opts *options.QuizzesQuestionsListOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "quizzes.questions.list", map[string]interface{}{
		"course_id": opts.CourseID,
		"quiz_id":   opts.QuizID,
	})

	service := api.NewQuizQuestionsService(client)

	questions, err := service.List(ctx, opts.CourseID, opts.QuizID, nil)
	if err != nil {
		logger.LogCommandError(ctx, "quizzes.questions.list", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"quiz_id":   opts.QuizID,
		})
		return fmt.Errorf("failed to list questions: %w", err)
	}

	if err := formatEmptyOrOutput(questions, "No questions found"); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "quizzes.questions.list", len(questions))
	return nil
}

func runQuizzesQuestionsGet(ctx context.Context, client *api.Client, opts *options.QuizzesQuestionsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "quizzes.questions.get", map[string]interface{}{
		"course_id":   opts.CourseID,
		"quiz_id":     opts.QuizID,
		"question_id": opts.QuestionID,
	})

	service := api.NewQuizQuestionsService(client)

	question, err := service.Get(ctx, opts.CourseID, opts.QuizID, opts.QuestionID)
	if err != nil {
		logger.LogCommandError(ctx, "quizzes.questions.get", err, map[string]interface{}{
			"course_id":   opts.CourseID,
			"quiz_id":     opts.QuizID,
			"question_id": opts.QuestionID,
		})
		return fmt.Errorf("failed to get question: %w", err)
	}

	if err := formatOutput(question, nil); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "quizzes.questions.get", 1)
	return nil
}

func runQuizzesQuestionsCreate(ctx context.Context, client *api.Client, opts *options.QuizzesQuestionsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "quizzes.questions.create", map[string]interface{}{
		"course_id":     opts.CourseID,
		"quiz_id":       opts.QuizID,
		"question_type": opts.QuestionType,
	})

	service := api.NewQuizQuestionsService(client)

	params := &api.CreateQuizQuestionParams{
		QuestionName:      opts.QuestionName,
		QuestionText:      opts.QuestionText,
		QuestionType:      opts.QuestionType,
		PointsPossible:    opts.PointsPossible,
		CorrectComments:   opts.CorrectComments,
		IncorrectComments: opts.IncorrectComments,
	}

	question, err := service.Create(ctx, opts.CourseID, opts.QuizID, params)
	if err != nil {
		logger.LogCommandError(ctx, "quizzes.questions.create", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"quiz_id":   opts.QuizID,
		})
		return fmt.Errorf("failed to create question: %w", err)
	}

	printInfo("Question created successfully (ID: %d)\n", question.ID)
	if err := formatOutput(question, nil); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "quizzes.questions.create", 1)
	return nil
}

func runQuizzesQuestionsDelete(ctx context.Context, client *api.Client, opts *options.QuizzesQuestionsDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "quizzes.questions.delete", map[string]interface{}{
		"course_id":   opts.CourseID,
		"quiz_id":     opts.QuizID,
		"question_id": opts.QuestionID,
	})

	if !opts.Force {
		fmt.Printf("WARNING: This will delete question %d from quiz %d.\n", opts.QuestionID, opts.QuizID)
		fmt.Print("Type 'yes' to confirm: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Delete cancelled")
			return nil
		}
	}

	service := api.NewQuizQuestionsService(client)

	if err := service.Delete(ctx, opts.CourseID, opts.QuizID, opts.QuestionID); err != nil {
		logger.LogCommandError(ctx, "quizzes.questions.delete", err, map[string]interface{}{
			"course_id":   opts.CourseID,
			"quiz_id":     opts.QuizID,
			"question_id": opts.QuestionID,
		})
		return fmt.Errorf("failed to delete question: %w", err)
	}

	printInfo("Question %d deleted\n", opts.QuestionID)

	logger.LogCommandComplete(ctx, "quizzes.questions.delete", 1)
	return nil
}

func runQuizzesSubmissionsList(ctx context.Context, client *api.Client, opts *options.QuizzesSubmissionsListOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "quizzes.submissions.list", map[string]interface{}{
		"course_id": opts.CourseID,
		"quiz_id":   opts.QuizID,
	})

	service := api.NewQuizSubmissionsService(client)

	submissions, err := service.List(ctx, opts.CourseID, opts.QuizID, nil)
	if err != nil {
		logger.LogCommandError(ctx, "quizzes.submissions.list", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"quiz_id":   opts.QuizID,
		})
		return fmt.Errorf("failed to list submissions: %w", err)
	}

	if err := formatEmptyOrOutput(submissions, "No submissions found"); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "quizzes.submissions.list", len(submissions))
	return nil
}

func runQuizzesSubmissionsGet(ctx context.Context, client *api.Client, opts *options.QuizzesSubmissionsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "quizzes.submissions.get", map[string]interface{}{
		"course_id":     opts.CourseID,
		"quiz_id":       opts.QuizID,
		"submission_id": opts.SubmissionID,
	})

	service := api.NewQuizSubmissionsService(client)

	submission, err := service.Get(ctx, opts.CourseID, opts.QuizID, opts.SubmissionID, nil)
	if err != nil {
		logger.LogCommandError(ctx, "quizzes.submissions.get", err, map[string]interface{}{
			"course_id":     opts.CourseID,
			"quiz_id":       opts.QuizID,
			"submission_id": opts.SubmissionID,
		})
		return fmt.Errorf("failed to get submission: %w", err)
	}

	if err := formatOutput(submission, nil); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "quizzes.submissions.get", 1)
	return nil
}
