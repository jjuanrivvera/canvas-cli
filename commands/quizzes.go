package commands

import (
	"context"
	"encoding/json"
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

// quizzesGroupsCmd represents the quizzes question-groups command group
var quizzesGroupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "Manage quiz question groups",
	Long:  `Manage question groups within a quiz for random question selection.`,
}

// quizzesReportsCmd represents the quizzes reports command group
var quizzesReportsCmd = &cobra.Command{
	Use:   "reports",
	Short: "Manage quiz reports",
	Long:  `Generate and view quiz reports such as student analysis and item analysis.`,
}

// quizzesStatisticsCmd represents the quizzes statistics command group
var quizzesStatisticsCmd = &cobra.Command{
	Use:   "statistics",
	Short: "View quiz statistics",
	Long:  `View statistical data for a quiz including score distributions and question analysis.`,
}

// quizzesExtensionsCmd represents the quizzes extensions command group
var quizzesExtensionsCmd = &cobra.Command{
	Use:   "extensions",
	Short: "Manage quiz extensions",
	Long:  `Grant quiz time extensions or extra attempts to students.`,
}

// quizzesIPFiltersCmd represents the quizzes ip-filters command group
var quizzesIPFiltersCmd = &cobra.Command{
	Use:   "ip-filters",
	Short: "List quiz IP filters",
	Long:  `List available IP address filters for quizzes in a course.`,
}

// quizzesAssignmentOverridesCmd represents the quizzes assignment-overrides command group
var quizzesAssignmentOverridesCmd = &cobra.Command{
	Use:   "assignment-overrides",
	Short: "Manage quiz assignment overrides",
	Long:  `View and set assignment overrides for quizzes (due dates per section/student).`,
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
	quizzesCmd.AddCommand(quizzesGroupsCmd)
	quizzesCmd.AddCommand(quizzesReportsCmd)
	quizzesCmd.AddCommand(quizzesStatisticsCmd)
	quizzesCmd.AddCommand(quizzesExtensionsCmd)
	quizzesCmd.AddCommand(quizzesIPFiltersCmd)
	quizzesCmd.AddCommand(quizzesAssignmentOverridesCmd)

	// Questions subcommands
	quizzesQuestionsCmd.AddCommand(newQuizzesQuestionsListCmd())
	quizzesQuestionsCmd.AddCommand(newQuizzesQuestionsGetCmd())
	quizzesQuestionsCmd.AddCommand(newQuizzesQuestionsCreateCmd())
	quizzesQuestionsCmd.AddCommand(newQuizzesQuestionsUpdateCmd())
	quizzesQuestionsCmd.AddCommand(newQuizzesQuestionsDeleteCmd())

	// Submissions subcommands
	quizzesSubmissionsCmd.AddCommand(newQuizzesSubmissionsListCmd())
	quizzesSubmissionsCmd.AddCommand(newQuizzesSubmissionsGetCmd())
	quizzesSubmissionsCmd.AddCommand(newQuizzesSubmissionsCreateCmd())

	// Question groups subcommands
	quizzesGroupsCmd.AddCommand(newQuizzesGroupsGetCmd())
	quizzesGroupsCmd.AddCommand(newQuizzesGroupsCreateCmd())
	quizzesGroupsCmd.AddCommand(newQuizzesGroupsUpdateCmd())
	quizzesGroupsCmd.AddCommand(newQuizzesGroupsDeleteCmd())

	// Reports subcommands
	quizzesReportsCmd.AddCommand(newQuizzesReportsListCmd())
	quizzesReportsCmd.AddCommand(newQuizzesReportsGetCmd())
	quizzesReportsCmd.AddCommand(newQuizzesReportsCreateCmd())
	quizzesReportsCmd.AddCommand(newQuizzesReportsDeleteCmd())

	// Statistics subcommands
	quizzesStatisticsCmd.AddCommand(newQuizzesStatisticsListCmd())

	// Extensions subcommands
	quizzesExtensionsCmd.AddCommand(newQuizzesExtensionsCreateCmd())

	// IP filters subcommands
	quizzesIPFiltersCmd.AddCommand(newQuizzesIPFiltersListCmd())

	// Assignment overrides subcommands
	quizzesAssignmentOverridesCmd.AddCommand(newQuizzesAssignmentOverridesListCmd())
	quizzesAssignmentOverridesCmd.AddCommand(newQuizzesAssignmentOverridesSetCmd())
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
	mustMarkRequired(cmd, "course-id")

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
	mustMarkRequired(cmd, "course-id")

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
	mustMarkRequired(cmd, "course-id", "title")

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
	mustMarkRequired(cmd, "course-id")

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
	mustMarkRequired(cmd, "course-id")

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
	mustMarkRequired(cmd, "course-id", "quiz-id")

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
	mustMarkRequired(cmd, "course-id", "quiz-id")

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
	mustMarkRequired(cmd, "course-id", "quiz-id", "text")

	return cmd
}

func newQuizzesQuestionsUpdateCmd() *cobra.Command {
	opts := &options.QuizzesQuestionsUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <question-id>",
		Short: "Update a question",
		Long: `Update an existing question in a quiz. Only the flags you pass are sent;
other fields keep their current values.

--answers-json takes a JSON array in the same shape "questions get" returns
(id, text, weight, comments, ...). For multiple-choice and true/false
questions the correct answer has weight 100 and every other answer weight 0.

Examples:
  canvas quizzes questions update 789 --course-id 123 --quiz-id 456 --points 5
  canvas quizzes questions update 789 --course-id 123 --quiz-id 456 --text "What is 2+2?" --name "Arithmetic"
  canvas quizzes questions update 789 --course-id 123 --quiz-id 456 \
    --answers-json '[{"id":1001,"text":"4","weight":100},{"id":1002,"text":"5","weight":0}]'`,
		Args: ExactArgsWithUsage(1, "question-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			questionID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid question ID: %s", args[0])
			}
			opts.QuestionID = questionID

			// Track which flags were set
			opts.QuestionNameSet = cmd.Flags().Changed("name")
			opts.QuestionTextSet = cmd.Flags().Changed("text")
			opts.QuestionTypeSet = cmd.Flags().Changed("type")
			opts.PointsPossibleSet = cmd.Flags().Changed("points")
			opts.CorrectCommentsSet = cmd.Flags().Changed("correct-comments")
			opts.IncorrectCommentsSet = cmd.Flags().Changed("incorrect-comments")
			opts.PositionSet = cmd.Flags().Changed("position")
			opts.AnswersJSONSet = cmd.Flags().Changed("answers-json")

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runQuizzesQuestionsUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.QuizID, "quiz-id", 0, "Quiz ID (required)")
	cmd.Flags().StringVar(&opts.QuestionName, "name", "", "Question name")
	cmd.Flags().StringVar(&opts.QuestionText, "text", "", "Question text")
	cmd.Flags().StringVar(&opts.QuestionType, "type", "", "Question type")
	cmd.Flags().Float64Var(&opts.PointsPossible, "points", 0, "Points possible")
	cmd.Flags().StringVar(&opts.CorrectComments, "correct-comments", "", "Comments for correct answer")
	cmd.Flags().StringVar(&opts.IncorrectComments, "incorrect-comments", "", "Comments for incorrect answer")
	cmd.Flags().IntVar(&opts.Position, "position", 0, "Position of the question in the quiz")
	cmd.Flags().StringVar(&opts.AnswersJSON, "answers-json", "", "Answers as a JSON array (replaces all answers)")
	mustMarkRequired(cmd, "course-id", "quiz-id")

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
	mustMarkRequired(cmd, "course-id", "quiz-id")

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
	mustMarkRequired(cmd, "course-id", "quiz-id")

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
	mustMarkRequired(cmd, "course-id", "quiz-id")

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
	if _, err := validateCourseID(ctx, client, opts.CourseID); err != nil {
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

	ok, confirmErr := confirmDelete("quiz", opts.QuizID, opts.Force)
	if confirmErr != nil {
		return confirmErr
	}
	if !ok {
		return nil
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

func runQuizzesQuestionsUpdate(ctx context.Context, client *api.Client, opts *options.QuizzesQuestionsUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "quizzes.questions.update", map[string]interface{}{
		"course_id":   opts.CourseID,
		"quiz_id":     opts.QuizID,
		"question_id": opts.QuestionID,
	})

	params, err := buildUpdateQuizQuestionParams(opts)
	if err != nil {
		return err
	}

	service := api.NewQuizQuestionsService(client)

	question, err := service.Update(ctx, opts.CourseID, opts.QuizID, opts.QuestionID, params)
	if err != nil {
		logger.LogCommandError(ctx, "quizzes.questions.update", err, map[string]interface{}{
			"course_id":   opts.CourseID,
			"quiz_id":     opts.QuizID,
			"question_id": opts.QuestionID,
		})
		return fmt.Errorf("failed to update question: %w", err)
	}

	// In dry-run mode the client has printed the curl command and answered
	// with an empty body; there is no updated question to report.
	if dryRun {
		return nil
	}

	printInfo("Question updated successfully (ID: %d)\n", question.ID)
	if err := formatOutput(question, nil); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "quizzes.questions.update", 1)
	return nil
}

// buildUpdateQuizQuestionParams maps the explicitly-set flags onto the
// pointer-typed update params so unset fields are omitted from the request.
func buildUpdateQuizQuestionParams(opts *options.QuizzesQuestionsUpdateOptions) (*api.UpdateQuizQuestionParams, error) {
	params := &api.UpdateQuizQuestionParams{}

	if opts.QuestionNameSet {
		params.QuestionName = &opts.QuestionName
	}
	if opts.QuestionTextSet {
		params.QuestionText = &opts.QuestionText
	}
	if opts.QuestionTypeSet {
		params.QuestionType = &opts.QuestionType
	}
	if opts.PointsPossibleSet {
		params.PointsPossible = &opts.PointsPossible
	}
	if opts.CorrectCommentsSet {
		params.CorrectComments = &opts.CorrectComments
	}
	if opts.IncorrectCommentsSet {
		params.IncorrectComments = &opts.IncorrectComments
	}
	if opts.PositionSet {
		params.Position = &opts.Position
	}
	if opts.AnswersJSONSet {
		var answers []api.QuizAnswer
		if err := json.Unmarshal([]byte(opts.AnswersJSON), &answers); err != nil {
			return nil, fmt.Errorf("invalid --answers-json: %w", err)
		}
		params.Answers = &answers
	}

	return params, nil
}

func runQuizzesQuestionsDelete(ctx context.Context, client *api.Client, opts *options.QuizzesQuestionsDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "quizzes.questions.delete", map[string]interface{}{
		"course_id":   opts.CourseID,
		"quiz_id":     opts.QuizID,
		"question_id": opts.QuestionID,
	})

	ok, confirmErr := confirmDeleteWithDetails("quiz question", opts.QuestionID, map[string]interface{}{
		"quiz_id": opts.QuizID,
	}, opts.Force)
	if confirmErr != nil {
		return confirmErr
	}
	if !ok {
		return nil
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

// ---- Quiz Submissions: Create (start a submission) ----

func newQuizzesSubmissionsCreateCmd() *cobra.Command {
	opts := &options.QuizzesSubmissionsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Start a new quiz submission",
		Long: `Start taking a quiz by creating a new submission.

Examples:
  canvas quizzes submissions create --course-id 123 --quiz-id 456
  canvas quizzes submissions create --course-id 123 --quiz-id 456 --access-code secret`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runQuizzesSubmissionsCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.QuizID, "quiz-id", 0, "Quiz ID (required)")
	cmd.Flags().StringVar(&opts.AccessCode, "access-code", "", "Quiz access code")
	mustMarkRequired(cmd, "course-id", "quiz-id")

	return cmd
}

func runQuizzesSubmissionsCreate(ctx context.Context, client *api.Client, opts *options.QuizzesSubmissionsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "quizzes.submissions.create", map[string]interface{}{
		"course_id": opts.CourseID,
		"quiz_id":   opts.QuizID,
	})

	service := api.NewQuizSubmissionsService(client)

	params := &api.StartQuizSubmissionParams{
		AccessCode: opts.AccessCode,
	}

	submission, err := service.Create(ctx, opts.CourseID, opts.QuizID, params)
	if err != nil {
		logger.LogCommandError(ctx, "quizzes.submissions.create", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"quiz_id":   opts.QuizID,
		})
		return fmt.Errorf("failed to start quiz submission: %w", err)
	}

	printInfo("Quiz submission started (ID: %d)\n", submission.ID)
	if err := formatOutput(submission, nil); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "quizzes.submissions.create", 1)
	return nil
}

// ---- Quiz Question Groups ----

func newQuizzesGroupsGetCmd() *cobra.Command {
	opts := &options.QuizzesGroupsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <group-id>",
		Short: "Get a quiz question group",
		Long: `Get details of a specific quiz question group.

Examples:
  canvas quizzes groups get 3 --course-id 123 --quiz-id 456`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %s", args[0])
			}
			opts.GroupID = groupID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runQuizzesGroupsGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.QuizID, "quiz-id", 0, "Quiz ID (required)")
	mustMarkRequired(cmd, "course-id", "quiz-id")

	return cmd
}

func runQuizzesGroupsGet(ctx context.Context, client *api.Client, opts *options.QuizzesGroupsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "quizzes.groups.get", map[string]interface{}{
		"course_id": opts.CourseID,
		"quiz_id":   opts.QuizID,
		"group_id":  opts.GroupID,
	})

	service := api.NewQuizQuestionGroupsService(client)

	group, err := service.Get(ctx, opts.CourseID, opts.QuizID, opts.GroupID)
	if err != nil {
		logger.LogCommandError(ctx, "quizzes.groups.get", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"quiz_id":   opts.QuizID,
			"group_id":  opts.GroupID,
		})
		return fmt.Errorf("failed to get quiz question group: %w", err)
	}

	if err := formatOutput(group, nil); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "quizzes.groups.get", 1)
	return nil
}

func newQuizzesGroupsCreateCmd() *cobra.Command {
	opts := &options.QuizzesGroupsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a quiz question group",
		Long: `Create a new question group in a quiz for random question selection.

Examples:
  canvas quizzes groups create --course-id 123 --quiz-id 456 --name "Chapter 5 Questions" --pick-count 5 --points 2`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runQuizzesGroupsCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.QuizID, "quiz-id", 0, "Quiz ID (required)")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Group name")
	cmd.Flags().IntVar(&opts.PickCount, "pick-count", 0, "Number of questions to pick from group")
	cmd.Flags().Float64Var(&opts.QuestionPoints, "points", 0, "Points per question in group")
	mustMarkRequired(cmd, "course-id", "quiz-id")

	return cmd
}

func runQuizzesGroupsCreate(ctx context.Context, client *api.Client, opts *options.QuizzesGroupsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "quizzes.groups.create", map[string]interface{}{
		"course_id": opts.CourseID,
		"quiz_id":   opts.QuizID,
		"name":      opts.Name,
	})

	service := api.NewQuizQuestionGroupsService(client)

	params := &api.CreateQuizQuestionGroupParams{
		Name:           opts.Name,
		PickCount:      opts.PickCount,
		QuestionPoints: opts.QuestionPoints,
	}

	group, err := service.Create(ctx, opts.CourseID, opts.QuizID, params)
	if err != nil {
		logger.LogCommandError(ctx, "quizzes.groups.create", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"quiz_id":   opts.QuizID,
		})
		return fmt.Errorf("failed to create quiz question group: %w", err)
	}

	printInfo("Quiz question group created (ID: %d)\n", group.ID)
	if err := formatOutput(group, nil); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "quizzes.groups.create", 1)
	return nil
}

func newQuizzesGroupsUpdateCmd() *cobra.Command {
	opts := &options.QuizzesGroupsUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <group-id>",
		Short: "Update a quiz question group",
		Long: `Update an existing quiz question group.

Examples:
  canvas quizzes groups update 3 --course-id 123 --quiz-id 456 --pick-count 10`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %s", args[0])
			}
			opts.GroupID = groupID
			opts.NameSet = cmd.Flags().Changed("name")
			opts.PickCountSet = cmd.Flags().Changed("pick-count")
			opts.PointsSet = cmd.Flags().Changed("points")

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runQuizzesGroupsUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.QuizID, "quiz-id", 0, "Quiz ID (required)")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Group name")
	cmd.Flags().IntVar(&opts.PickCount, "pick-count", 0, "Number of questions to pick from group")
	cmd.Flags().Float64Var(&opts.QuestionPoints, "points", 0, "Points per question in group")
	mustMarkRequired(cmd, "course-id", "quiz-id")

	return cmd
}

func runQuizzesGroupsUpdate(ctx context.Context, client *api.Client, opts *options.QuizzesGroupsUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "quizzes.groups.update", map[string]interface{}{
		"course_id": opts.CourseID,
		"quiz_id":   opts.QuizID,
		"group_id":  opts.GroupID,
	})

	service := api.NewQuizQuestionGroupsService(client)

	params := &api.UpdateQuizQuestionGroupParams{}
	if opts.NameSet {
		params.Name = &opts.Name
	}
	if opts.PickCountSet {
		params.PickCount = &opts.PickCount
	}
	if opts.PointsSet {
		params.QuestionPoints = &opts.QuestionPoints
	}

	group, err := service.Update(ctx, opts.CourseID, opts.QuizID, opts.GroupID, params)
	if err != nil {
		logger.LogCommandError(ctx, "quizzes.groups.update", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"quiz_id":   opts.QuizID,
			"group_id":  opts.GroupID,
		})
		return fmt.Errorf("failed to update quiz question group: %w", err)
	}

	printInfo("Quiz question group updated (ID: %d)\n", group.ID)
	if err := formatOutput(group, nil); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "quizzes.groups.update", 1)
	return nil
}

func newQuizzesGroupsDeleteCmd() *cobra.Command {
	opts := &options.QuizzesGroupsDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <group-id>",
		Short: "Delete a quiz question group",
		Long: `Delete a quiz question group.

Examples:
  canvas quizzes groups delete 3 --course-id 123 --quiz-id 456 --force`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %s", args[0])
			}
			opts.GroupID = groupID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runQuizzesGroupsDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.QuizID, "quiz-id", 0, "Quiz ID (required)")
	cmd.Flags().BoolVar(&opts.Force, "force", false, "Skip confirmation prompt")
	mustMarkRequired(cmd, "course-id", "quiz-id")

	return cmd
}

func runQuizzesGroupsDelete(ctx context.Context, client *api.Client, opts *options.QuizzesGroupsDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "quizzes.groups.delete", map[string]interface{}{
		"course_id": opts.CourseID,
		"quiz_id":   opts.QuizID,
		"group_id":  opts.GroupID,
	})

	ok, confirmErr := confirmDeleteWithDetails("quiz question group", opts.GroupID, map[string]interface{}{
		"quiz_id": opts.QuizID,
	}, opts.Force)
	if confirmErr != nil {
		return confirmErr
	}
	if !ok {
		return nil
	}

	service := api.NewQuizQuestionGroupsService(client)

	if err := service.Delete(ctx, opts.CourseID, opts.QuizID, opts.GroupID); err != nil {
		logger.LogCommandError(ctx, "quizzes.groups.delete", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"quiz_id":   opts.QuizID,
			"group_id":  opts.GroupID,
		})
		return fmt.Errorf("failed to delete quiz question group: %w", err)
	}

	printInfo("Quiz question group %d deleted\n", opts.GroupID)
	logger.LogCommandComplete(ctx, "quizzes.groups.delete", 1)
	return nil
}

// ---- Quiz Reports ----

func newQuizzesReportsListCmd() *cobra.Command {
	opts := &options.QuizzesReportsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List quiz reports",
		Long: `List all reports for a quiz.

Examples:
  canvas quizzes reports list --course-id 123 --quiz-id 456`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runQuizzesReportsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.QuizID, "quiz-id", 0, "Quiz ID (required)")
	cmd.Flags().BoolVar(&opts.IncludesAllVersions, "all-versions", false, "Include all quiz versions")
	mustMarkRequired(cmd, "course-id", "quiz-id")

	return cmd
}

func runQuizzesReportsList(ctx context.Context, client *api.Client, opts *options.QuizzesReportsListOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "quizzes.reports.list", map[string]interface{}{
		"course_id": opts.CourseID,
		"quiz_id":   opts.QuizID,
	})

	service := api.NewQuizReportsService(client)

	apiOpts := &api.ListQuizReportsOptions{IncludesAllVersions: opts.IncludesAllVersions}
	reports, err := service.List(ctx, opts.CourseID, opts.QuizID, apiOpts)
	if err != nil {
		logger.LogCommandError(ctx, "quizzes.reports.list", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"quiz_id":   opts.QuizID,
		})
		return fmt.Errorf("failed to list quiz reports: %w", err)
	}

	if err := formatEmptyOrOutput(reports, "No reports found"); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "quizzes.reports.list", len(reports))
	return nil
}

func newQuizzesReportsGetCmd() *cobra.Command {
	opts := &options.QuizzesReportsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <report-id>",
		Short: "Get a quiz report",
		Long: `Get details of a specific quiz report.

Examples:
  canvas quizzes reports get 5 --course-id 123 --quiz-id 456`,
		Args: ExactArgsWithUsage(1, "report-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			reportID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid report ID: %s", args[0])
			}
			opts.ReportID = reportID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runQuizzesReportsGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.QuizID, "quiz-id", 0, "Quiz ID (required)")
	mustMarkRequired(cmd, "course-id", "quiz-id")

	return cmd
}

func runQuizzesReportsGet(ctx context.Context, client *api.Client, opts *options.QuizzesReportsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "quizzes.reports.get", map[string]interface{}{
		"course_id": opts.CourseID,
		"quiz_id":   opts.QuizID,
		"report_id": opts.ReportID,
	})

	service := api.NewQuizReportsService(client)

	report, err := service.Get(ctx, opts.CourseID, opts.QuizID, opts.ReportID, false)
	if err != nil {
		logger.LogCommandError(ctx, "quizzes.reports.get", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"quiz_id":   opts.QuizID,
			"report_id": opts.ReportID,
		})
		return fmt.Errorf("failed to get quiz report: %w", err)
	}

	if err := formatOutput(report, nil); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "quizzes.reports.get", 1)
	return nil
}

func newQuizzesReportsCreateCmd() *cobra.Command {
	opts := &options.QuizzesReportsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create (or retrieve) a quiz report",
		Long: `Generate a quiz report. Report types: student_analysis, item_analysis.

Examples:
  canvas quizzes reports create --course-id 123 --quiz-id 456 --report-type student_analysis`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runQuizzesReportsCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.QuizID, "quiz-id", 0, "Quiz ID (required)")
	cmd.Flags().StringVar(&opts.ReportType, "report-type", "", "Report type: student_analysis, item_analysis (required)")
	cmd.Flags().BoolVar(&opts.IncludesAllVersions, "all-versions", false, "Include all quiz versions")
	mustMarkRequired(cmd, "course-id", "quiz-id", "report-type")

	return cmd
}

func runQuizzesReportsCreate(ctx context.Context, client *api.Client, opts *options.QuizzesReportsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "quizzes.reports.create", map[string]interface{}{
		"course_id":   opts.CourseID,
		"quiz_id":     opts.QuizID,
		"report_type": opts.ReportType,
	})

	service := api.NewQuizReportsService(client)

	params := &api.CreateQuizReportParams{
		ReportType:          opts.ReportType,
		IncludesAllVersions: opts.IncludesAllVersions,
	}

	report, err := service.Create(ctx, opts.CourseID, opts.QuizID, params)
	if err != nil {
		logger.LogCommandError(ctx, "quizzes.reports.create", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"quiz_id":   opts.QuizID,
		})
		return fmt.Errorf("failed to create quiz report: %w", err)
	}

	printInfo("Quiz report created (ID: %d)\n", report.ID)
	if err := formatOutput(report, nil); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "quizzes.reports.create", 1)
	return nil
}

func newQuizzesReportsDeleteCmd() *cobra.Command {
	opts := &options.QuizzesReportsDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <report-id>",
		Short: "Delete a quiz report",
		Long: `Delete a quiz report.

Examples:
  canvas quizzes reports delete 5 --course-id 123 --quiz-id 456 --force`,
		Args: ExactArgsWithUsage(1, "report-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			reportID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid report ID: %s", args[0])
			}
			opts.ReportID = reportID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runQuizzesReportsDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.QuizID, "quiz-id", 0, "Quiz ID (required)")
	cmd.Flags().BoolVar(&opts.Force, "force", false, "Skip confirmation prompt")
	mustMarkRequired(cmd, "course-id", "quiz-id")

	return cmd
}

func runQuizzesReportsDelete(ctx context.Context, client *api.Client, opts *options.QuizzesReportsDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "quizzes.reports.delete", map[string]interface{}{
		"course_id": opts.CourseID,
		"quiz_id":   opts.QuizID,
		"report_id": opts.ReportID,
	})

	ok, confirmErr := confirmDeleteWithDetails("quiz report", opts.ReportID, map[string]interface{}{
		"quiz_id": opts.QuizID,
	}, opts.Force)
	if confirmErr != nil {
		return confirmErr
	}
	if !ok {
		return nil
	}

	service := api.NewQuizReportsService(client)

	if err := service.Delete(ctx, opts.CourseID, opts.QuizID, opts.ReportID); err != nil {
		logger.LogCommandError(ctx, "quizzes.reports.delete", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"quiz_id":   opts.QuizID,
			"report_id": opts.ReportID,
		})
		return fmt.Errorf("failed to delete quiz report: %w", err)
	}

	printInfo("Quiz report %d deleted\n", opts.ReportID)
	logger.LogCommandComplete(ctx, "quizzes.reports.delete", 1)
	return nil
}

// ---- Quiz Statistics ----

func newQuizzesStatisticsListCmd() *cobra.Command {
	opts := &options.QuizzesStatisticsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "View quiz statistics",
		Long: `View statistical data for a quiz.

Examples:
  canvas quizzes statistics list --course-id 123 --quiz-id 456`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runQuizzesStatisticsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.QuizID, "quiz-id", 0, "Quiz ID (required)")
	mustMarkRequired(cmd, "course-id", "quiz-id")

	return cmd
}

func runQuizzesStatisticsList(ctx context.Context, client *api.Client, opts *options.QuizzesStatisticsListOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "quizzes.statistics.list", map[string]interface{}{
		"course_id": opts.CourseID,
		"quiz_id":   opts.QuizID,
	})

	service := api.NewQuizStatisticsService(client)

	stats, err := service.List(ctx, opts.CourseID, opts.QuizID)
	if err != nil {
		logger.LogCommandError(ctx, "quizzes.statistics.list", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"quiz_id":   opts.QuizID,
		})
		return fmt.Errorf("failed to get quiz statistics: %w", err)
	}

	if err := formatEmptyOrOutput(stats, "No statistics found"); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "quizzes.statistics.list", len(stats))
	return nil
}

// ---- Quiz Extensions ----

func newQuizzesExtensionsCreateCmd() *cobra.Command {
	opts := &options.QuizzesExtensionsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Grant a quiz extension to a student",
		Long: `Grant extra time or attempts on a quiz to a specific student.

Examples:
  canvas quizzes extensions create --course-id 123 --quiz-id 456 --user-id 789 --extra-time 30
  canvas quizzes extensions create --course-id 123 --quiz-id 456 --user-id 789 --extra-attempts 2`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runQuizzesExtensionsCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.QuizID, "quiz-id", 0, "Quiz ID (required)")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID to extend (required)")
	cmd.Flags().IntVar(&opts.ExtraAttempts, "extra-attempts", 0, "Extra attempts to grant")
	cmd.Flags().IntVar(&opts.ExtraTime, "extra-time", 0, "Extra time in minutes")
	cmd.Flags().BoolVar(&opts.ManuallyUnlocked, "manually-unlocked", false, "Manually unlock for student")
	cmd.Flags().IntVar(&opts.ExtendFromNow, "extend-from-now", 0, "Extend from now by N minutes")
	mustMarkRequired(cmd, "course-id", "quiz-id", "user-id")

	return cmd
}

func runQuizzesExtensionsCreate(ctx context.Context, client *api.Client, opts *options.QuizzesExtensionsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "quizzes.extensions.create", map[string]interface{}{
		"course_id": opts.CourseID,
		"quiz_id":   opts.QuizID,
		"user_id":   opts.UserID,
	})

	service := api.NewQuizExtensionsService(client)

	entries := []api.QuizExtensionEntry{
		{
			UserID:           opts.UserID,
			ExtraAttempts:    opts.ExtraAttempts,
			ExtraTime:        opts.ExtraTime,
			ManuallyUnlocked: opts.ManuallyUnlocked,
			ExtendFromNow:    opts.ExtendFromNow,
		},
	}

	exts, err := service.Create(ctx, opts.CourseID, opts.QuizID, entries)
	if err != nil {
		logger.LogCommandError(ctx, "quizzes.extensions.create", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"quiz_id":   opts.QuizID,
			"user_id":   opts.UserID,
		})
		return fmt.Errorf("failed to create quiz extension: %w", err)
	}

	if err := formatEmptyOrOutput(exts, "No extension data returned"); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "quizzes.extensions.create", len(exts))
	return nil
}

// ---- Quiz IP Filters ----

func newQuizzesIPFiltersListCmd() *cobra.Command {
	opts := &options.QuizzesIPFiltersListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List quiz IP filters",
		Long: `List IP address filters available for quizzes in a course.

Examples:
  canvas quizzes ip-filters list --course-id 123 --quiz-id 456`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runQuizzesIPFiltersList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.QuizID, "quiz-id", 0, "Quiz ID (required)")
	mustMarkRequired(cmd, "course-id", "quiz-id")

	return cmd
}

func runQuizzesIPFiltersList(ctx context.Context, client *api.Client, opts *options.QuizzesIPFiltersListOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "quizzes.ip-filters.list", map[string]interface{}{
		"course_id": opts.CourseID,
		"quiz_id":   opts.QuizID,
	})

	service := api.NewQuizIPFiltersService(client)

	filters, err := service.List(ctx, opts.CourseID, opts.QuizID)
	if err != nil {
		logger.LogCommandError(ctx, "quizzes.ip-filters.list", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to list quiz IP filters: %w", err)
	}

	if err := formatEmptyOrOutput(filters, "No IP filters found"); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "quizzes.ip-filters.list", len(filters))
	return nil
}

// ---- Quiz Assignment Overrides ----

func newQuizzesAssignmentOverridesListCmd() *cobra.Command {
	opts := &options.QuizzesAssignmentOverridesListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List quiz assignment overrides",
		Long: `List assignment overrides for quizzes in a course.

Examples:
  canvas quizzes assignment-overrides list --course-id 123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runQuizzesAssignmentOverridesList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	mustMarkRequired(cmd, "course-id")

	return cmd
}

func runQuizzesAssignmentOverridesList(ctx context.Context, client *api.Client, opts *options.QuizzesAssignmentOverridesListOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "quizzes.assignment-overrides.list", map[string]interface{}{
		"course_id": opts.CourseID,
	})

	service := api.NewQuizAssignmentOverridesService(client)

	overrides, err := service.List(ctx, opts.CourseID, opts.QuizIDs)
	if err != nil {
		logger.LogCommandError(ctx, "quizzes.assignment-overrides.list", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to list quiz assignment overrides: %w", err)
	}

	if err := formatEmptyOrOutput(overrides, "No assignment overrides found"); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "quizzes.assignment-overrides.list", len(overrides))
	return nil
}

func newQuizzesAssignmentOverridesSetCmd() *cobra.Command {
	opts := &options.QuizzesAssignmentOverridesSetOptions{}

	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set quiz assignment overrides",
		Long: `Create or update assignment overrides for a quiz.

Examples:
  canvas quizzes assignment-overrides set --course-id 123 --quiz-id 456 --section-id 789 --due-at "2026-08-01T23:59:00Z"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runQuizzesAssignmentOverridesSet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.QuizID, "quiz-id", 0, "Quiz ID (required)")
	cmd.Flags().Int64Var(&opts.CourseSectionID, "section-id", 0, "Course section ID for override")
	cmd.Flags().StringVar(&opts.DueAt, "due-at", "", "Due date (ISO 8601)")
	cmd.Flags().StringVar(&opts.UnlockAt, "unlock-at", "", "Unlock date (ISO 8601)")
	cmd.Flags().StringVar(&opts.LockAt, "lock-at", "", "Lock date (ISO 8601)")
	mustMarkRequired(cmd, "course-id", "quiz-id")

	return cmd
}

func runQuizzesAssignmentOverridesSet(ctx context.Context, client *api.Client, opts *options.QuizzesAssignmentOverridesSetOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "quizzes.assignment-overrides.set", map[string]interface{}{
		"course_id": opts.CourseID,
		"quiz_id":   opts.QuizID,
	})

	service := api.NewQuizAssignmentOverridesService(client)

	override := api.AssignmentOverrideEntry{}
	if opts.CourseSectionID > 0 {
		override.CourseSectionID = opts.CourseSectionID
	}
	if opts.DueAt != "" {
		override.DueAt = opts.DueAt
	}
	if opts.UnlockAt != "" {
		override.UnlockAt = opts.UnlockAt
	}
	if opts.LockAt != "" {
		override.LockAt = opts.LockAt
	}

	quizIDStr := strconv.FormatInt(opts.QuizID, 10)
	params := &api.SetQuizAssignmentOverridesParams{
		QuizAssignmentOverrides: []api.QuizAssignmentOverrideSetInput{
			{
				QuizID:    quizIDStr,
				Overrides: []api.AssignmentOverrideEntry{override},
			},
		},
	}

	overrides, err := service.Set(ctx, opts.CourseID, params)
	if err != nil {
		logger.LogCommandError(ctx, "quizzes.assignment-overrides.set", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"quiz_id":   opts.QuizID,
		})
		return fmt.Errorf("failed to set quiz assignment overrides: %w", err)
	}

	if err := formatEmptyOrOutput(overrides, "No overrides returned"); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "quizzes.assignment-overrides.set", len(overrides))
	return nil
}
