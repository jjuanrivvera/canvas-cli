package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
	"github.com/jjuanrivvera/canvas-cli/internal/batch"
)

// submissionsCmd represents the submissions command group
var submissionsCmd = &cobra.Command{
	Use:   "submissions",
	Short: "Manage Canvas submissions",
	Long: `Manage Canvas submissions including listing, viewing, and grading submissions.

Examples:
  canvas submissions list --course-id 123 --assignment-id 456
  canvas submissions get --course-id 123 --assignment-id 456 --user-id 789
  canvas submissions list --course-id 123 --assignment-id 456 --workflow-state graded`,
}

func init() {
	rootCmd.AddCommand(submissionsCmd)
	submissionsCmd.AddCommand(newSubmissionsListCmd())
	submissionsCmd.AddCommand(newSubmissionsGetCmd())
	submissionsCmd.AddCommand(newSubmissionsGradeCmd())
	submissionsCmd.AddCommand(newSubmissionsBulkGradeCmd())
	submissionsCmd.AddCommand(newSubmissionsCommentsCmd())
	submissionsCmd.AddCommand(newSubmissionsAddCommentCmd())
	submissionsCmd.AddCommand(newSubmissionsDeleteCommentCmd())
}

func newSubmissionsListCmd() *cobra.Command {
	opts := &options.SubmissionsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List submissions for an assignment",
		Long: `List all submissions for a Canvas assignment.

You can filter submissions by workflow state and graded since date.

Workflow States:
  - submitted: Submissions that have been submitted
  - unsubmitted: Submissions that have not been submitted
  - graded: Submissions that have been graded
  - pending_review: Submissions pending review

Examples:
  canvas submissions list --course-id 123 --assignment-id 456
  canvas submissions list --course-id 123 --assignment-id 456 --workflow-state graded
  canvas submissions list --course-id 123 --assignment-id 456 --include user,submission_comments
  canvas submissions list --course-id 123 --assignment-id 456 --graded-since 2024-01-01`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runSubmissionsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.AssignmentID, "assignment-id", 0, "Assignment ID (required)")
	cmd.Flags().StringVar(&opts.WorkflowState, "workflow-state", "", "Filter by workflow state (submitted, unsubmitted, graded, pending_review)")
	cmd.Flags().StringVar(&opts.GradedSince, "graded-since", "", "Filter by graded since date (ISO8601 format)")
	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Additional data to include (comma-separated)")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("assignment-id")

	return cmd
}

func newSubmissionsGetCmd() *cobra.Command {
	opts := &options.SubmissionsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a specific submission",
		Long: `Get details of a specific submission for an assignment and user.

Examples:
  canvas submissions get --course-id 123 --assignment-id 456 --user-id 789
  canvas submissions get --course-id 123 --assignment-id 456 --user-id 789 --include submission_comments,rubric_assessment`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runSubmissionsGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.AssignmentID, "assignment-id", 0, "Assignment ID (required)")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (required)")
	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Additional data to include (comma-separated)")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("assignment-id")
	cmd.MarkFlagRequired("user-id")

	return cmd
}

func newSubmissionsGradeCmd() *cobra.Command {
	opts := &options.SubmissionsGradeOptions{}

	cmd := &cobra.Command{
		Use:   "grade",
		Short: "Grade a submission",
		Long: `Grade a specific submission for an assignment and user.

You can provide a score, comment, or excuse the submission.

Examples:
  canvas submissions grade --course-id 123 --assignment-id 456 --user-id 789 --score 95
  canvas submissions grade --course-id 123 --assignment-id 456 --user-id 789 --score 85 --comment "Good work"
  canvas submissions grade --course-id 123 --assignment-id 456 --user-id 789 --excuse`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runSubmissionsGrade(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.AssignmentID, "assignment-id", 0, "Assignment ID (required)")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (required)")
	cmd.Flags().Float64Var(&opts.Score, "score", 0, "Score to assign")
	cmd.Flags().StringVar(&opts.Comment, "comment", "", "Comment to add")
	cmd.Flags().BoolVar(&opts.Excuse, "excuse", false, "Excuse the submission")
	cmd.Flags().StringVar(&opts.PostedGrade, "posted-grade", "", "Posted grade (e.g., 'A', 'B+', 'Pass')")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("assignment-id")
	cmd.MarkFlagRequired("user-id")

	return cmd
}

func newSubmissionsBulkGradeCmd() *cobra.Command {
	opts := &options.SubmissionsBulkGradeOptions{}

	cmd := &cobra.Command{
		Use:   "bulk-grade",
		Short: "Grade multiple submissions from CSV",
		Long: `Grade multiple submissions at once by importing grades from a CSV file.

The CSV file should have the following format:
  user_id,assignment_id,score,comment

Example CSV:
  123,456,95,"Excellent work"
  124,456,87,"Good job"
  125,456,92,"Great effort"

Examples:
  canvas submissions bulk-grade --course-id 123 --csv grades.csv
  canvas submissions bulk-grade --course-id 123 --csv grades.csv --dry-run`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runSubmissionsBulkGrade(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.CSV, "csv", "", "CSV file with grades (required)")
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Preview changes without applying them")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("csv")

	return cmd
}

func newSubmissionsCommentsCmd() *cobra.Command {
	opts := &options.SubmissionsCommentsOptions{}

	cmd := &cobra.Command{
		Use:   "comments",
		Short: "List comments for a submission",
		Long: `List all comments for a specific submission.

Examples:
  canvas submissions comments --course-id 123 --assignment-id 456 --user-id 789`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runSubmissionsComments(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.AssignmentID, "assignment-id", 0, "Assignment ID (required)")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("assignment-id")
	cmd.MarkFlagRequired("user-id")

	return cmd
}

func newSubmissionsAddCommentCmd() *cobra.Command {
	opts := &options.SubmissionsAddCommentOptions{}

	cmd := &cobra.Command{
		Use:   "add-comment",
		Short: "Add a comment to a submission",
		Long: `Add a comment to a specific submission.

Examples:
  canvas submissions add-comment --course-id 123 --assignment-id 456 --user-id 789 --text "Great work!"
  canvas submissions add-comment --course-id 123 --assignment-id 456 --user-id 789 --text "Feedback" --group`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runSubmissionsAddComment(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.AssignmentID, "assignment-id", 0, "Assignment ID (required)")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (required)")
	cmd.Flags().StringVar(&opts.Text, "text", "", "Comment text (required)")
	cmd.Flags().BoolVar(&opts.GroupShare, "group", false, "Share comment with group members")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("assignment-id")
	cmd.MarkFlagRequired("user-id")
	cmd.MarkFlagRequired("text")

	return cmd
}

func newSubmissionsDeleteCommentCmd() *cobra.Command {
	opts := &options.SubmissionsDeleteCommentOptions{}

	cmd := &cobra.Command{
		Use:   "delete-comment",
		Short: "Delete a submission comment",
		Long: `Delete a comment from a submission.

Examples:
  canvas submissions delete-comment --course-id 123 --assignment-id 456 --user-id 789 --comment-id 999`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runSubmissionsDeleteComment(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.AssignmentID, "assignment-id", 0, "Assignment ID (required)")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (required)")
	cmd.Flags().Int64Var(&opts.CommentID, "comment-id", 0, "Comment ID to delete (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("assignment-id")
	cmd.MarkFlagRequired("user-id")
	cmd.MarkFlagRequired("comment-id")

	return cmd
}

func runSubmissionsList(ctx context.Context, client *api.Client, opts *options.SubmissionsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "submissions.list", map[string]interface{}{
		"course_id":      opts.CourseID,
		"assignment_id":  opts.AssignmentID,
		"workflow_state": opts.WorkflowState,
	})

	// Validate course ID exists
	if _, err := validateCourseID(client, opts.CourseID); err != nil {
		logger.LogCommandError(ctx, "submissions.list", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return err
	}

	// Create submissions service
	submissionsService := api.NewSubmissionsService(client)

	// Build options
	apiOpts := &api.ListSubmissionsOptions{
		WorkflowState: opts.WorkflowState,
		GradedSince:   opts.GradedSince,
		Include:       opts.Include,
	}

	// List submissions
	submissions, err := submissionsService.List(ctx, opts.CourseID, opts.AssignmentID, apiOpts)
	if err != nil {
		logger.LogCommandError(ctx, "submissions.list", err, map[string]interface{}{
			"course_id":     opts.CourseID,
			"assignment_id": opts.AssignmentID,
		})
		return fmt.Errorf("failed to list submissions: %w", err)
	}

	if len(submissions) == 0 {
		logger.LogCommandComplete(ctx, "submissions.list", 0)
		fmt.Println("No submissions found")
		return nil
	}

	// Format and display submissions
	printVerbose("Found %d submissions:\n\n", len(submissions))
	logger.LogCommandComplete(ctx, "submissions.list", len(submissions))
	return formatOutput(submissions, nil)
}

func runSubmissionsGet(ctx context.Context, client *api.Client, opts *options.SubmissionsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "submissions.get", map[string]interface{}{
		"course_id":     opts.CourseID,
		"assignment_id": opts.AssignmentID,
		"user_id":       opts.UserID,
	})

	// Validate course ID exists
	if _, err := validateCourseID(client, opts.CourseID); err != nil {
		logger.LogCommandError(ctx, "submissions.get", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return err
	}

	// Create submissions service
	submissionsService := api.NewSubmissionsService(client)

	// Get submission
	submission, err := submissionsService.Get(ctx, opts.CourseID, opts.AssignmentID, opts.UserID, opts.Include)
	if err != nil {
		logger.LogCommandError(ctx, "submissions.get", err, map[string]interface{}{
			"course_id":     opts.CourseID,
			"assignment_id": opts.AssignmentID,
			"user_id":       opts.UserID,
		})
		return fmt.Errorf("failed to get submission: %w", err)
	}

	// Format and display submission details
	logger.LogCommandComplete(ctx, "submissions.get", 1)
	return formatOutput(submission, nil)
}

func runSubmissionsGrade(ctx context.Context, client *api.Client, opts *options.SubmissionsGradeOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "submissions.grade", map[string]interface{}{
		"course_id":     opts.CourseID,
		"assignment_id": opts.AssignmentID,
		"user_id":       opts.UserID,
		"has_score":     opts.Score > 0,
		"has_comment":   opts.Comment != "",
		"excuse":        opts.Excuse,
	})

	// Validate course ID exists
	if _, err := validateCourseID(client, opts.CourseID); err != nil {
		logger.LogCommandError(ctx, "submissions.grade", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return err
	}

	// Create submissions service
	submissionsService := api.NewSubmissionsService(client)

	// Build grade params
	params := &api.GradeSubmissionParams{}

	// Handle score (convert to string for PostedGrade)
	if opts.Score > 0 {
		params.PostedGrade = fmt.Sprintf("%.2f", opts.Score)
	} else if opts.PostedGrade != "" {
		params.PostedGrade = opts.PostedGrade
	}

	// Handle comment
	if opts.Comment != "" {
		params.Comment = &api.SubmissionCommentParams{
			TextComment: opts.Comment,
		}
	}

	// Handle excuse
	if opts.Excuse {
		params.Excuse = true
	}

	// Grade submission
	submission, err := submissionsService.Grade(ctx, opts.CourseID, opts.AssignmentID, opts.UserID, params)
	if err != nil {
		logger.LogCommandError(ctx, "submissions.grade", err, map[string]interface{}{
			"course_id":     opts.CourseID,
			"assignment_id": opts.AssignmentID,
			"user_id":       opts.UserID,
		})
		return fmt.Errorf("failed to grade submission: %w", err)
	}

	// Display success message
	userName := "Unknown"
	if submission.User != nil {
		userName = submission.User.Name
	}

	fmt.Printf("✅ Successfully graded submission for %s\n", userName)
	fmt.Printf("   User ID: %d\n", submission.UserID)
	fmt.Printf("   Assignment ID: %d\n", submission.AssignmentID)

	if submission.Score > 0 {
		fmt.Printf("   Score: %.1f\n", submission.Score)
	}

	if submission.Grade != "" {
		fmt.Printf("   Grade: %s\n", submission.Grade)
	}

	if submission.ExcusedTLN {
		fmt.Printf("   ✓ Excused\n")
	}

	if !submission.GradedAt.IsZero() {
		fmt.Printf("   Graded: %s\n", submission.GradedAt.Format("2006-01-02 15:04"))
	}

	logger.LogCommandComplete(ctx, "submissions.grade", 1)
	return nil
}

func runSubmissionsBulkGrade(ctx context.Context, client *api.Client, opts *options.SubmissionsBulkGradeOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "submissions.bulk-grade", map[string]interface{}{
		"course_id": opts.CourseID,
		"csv_file":  opts.CSV,
		"dry_run":   opts.DryRun,
	})

	// Validate course ID exists
	if _, err := validateCourseID(client, opts.CourseID); err != nil {
		logger.LogCommandError(ctx, "submissions.bulk-grade", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return err
	}

	// Create submissions service
	submissionsService := api.NewSubmissionsService(client)

	// Read grades from CSV
	grades, err := batch.ReadGradesCSV(opts.CSV)
	if err != nil {
		logger.LogCommandError(ctx, "submissions.bulk-grade", err, map[string]interface{}{
			"csv_file": opts.CSV,
		})
		return fmt.Errorf("failed to read CSV file: %w", err)
	}

	if len(grades) == 0 {
		err := fmt.Errorf("no grades found in CSV file")
		logger.LogCommandError(ctx, "submissions.bulk-grade", err, map[string]interface{}{
			"csv_file": opts.CSV,
		})
		return err
	}

	printVerbose("Found %d grades in CSV file\n\n", len(grades))

	if opts.DryRun {
		fmt.Println("DRY RUN - No changes will be applied")
		fmt.Println()
		fmt.Println("The following grades would be applied:")
		for i, grade := range grades {
			fmt.Printf("%d. User %d, Assignment %d: Score=%s", i+1, grade.UserID, grade.AssignmentID, grade.Grade)
			if grade.Comment != "" {
				fmt.Printf(", Comment=%s", grade.Comment)
			}
			fmt.Println()
		}
		logger.LogCommandComplete(ctx, "submissions.bulk-grade", 0)
		return nil
	}

	// Process grades
	successCount := 0
	errorCount := 0
	var errors []string

	for i, grade := range grades {
		fmt.Printf("Processing %d/%d: User %d, Assignment %d...", i+1, len(grades), grade.UserID, grade.AssignmentID)

		// Build params
		params := &api.GradeSubmissionParams{
			PostedGrade: grade.Grade,
		}

		if grade.Comment != "" {
			params.Comment = &api.SubmissionCommentParams{
				TextComment: grade.Comment,
			}
		}

		// Grade submission
		_, err = submissionsService.Grade(ctx, opts.CourseID, grade.AssignmentID, grade.UserID, params)
		if err != nil {
			fmt.Printf(" ❌ Error: %v\n", err)
			errorCount++
			errors = append(errors, fmt.Sprintf("Row %d: %v", grade.Row, err))
			continue
		}

		fmt.Printf(" ✅\n")
		successCount++
	}

	// Print summary
	fmt.Printf("\n═══════════════════════════════════════\n")
	fmt.Printf("Bulk Grading Complete\n")
	fmt.Printf("═══════════════════════════════════════\n")
	fmt.Printf("Total: %d\n", len(grades))
	fmt.Printf("Success: %d\n", successCount)
	fmt.Printf("Errors: %d\n", errorCount)

	if len(errors) > 0 {
		fmt.Printf("\nErrors:\n")
		for _, errMsg := range errors {
			fmt.Printf("  - %s\n", errMsg)
		}
	}

	logger.LogCommandComplete(ctx, "submissions.bulk-grade", successCount)

	if errorCount > 0 {
		return fmt.Errorf("bulk grading completed with %d errors", errorCount)
	}

	return nil
}

func runSubmissionsComments(ctx context.Context, client *api.Client, opts *options.SubmissionsCommentsOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "submissions.comments", map[string]interface{}{
		"course_id":     opts.CourseID,
		"assignment_id": opts.AssignmentID,
		"user_id":       opts.UserID,
	})

	if _, err := validateCourseID(client, opts.CourseID); err != nil {
		logger.LogCommandError(ctx, "submissions.comments", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return err
	}

	submissionsService := api.NewSubmissionsService(client)

	submission, err := submissionsService.Get(ctx, opts.CourseID, opts.AssignmentID, opts.UserID, []string{"submission_comments"})
	if err != nil {
		logger.LogCommandError(ctx, "submissions.comments", err, map[string]interface{}{
			"course_id":     opts.CourseID,
			"assignment_id": opts.AssignmentID,
			"user_id":       opts.UserID,
		})
		return fmt.Errorf("failed to get submission: %w", err)
	}

	if len(submission.SubmissionComments) == 0 {
		logger.LogCommandComplete(ctx, "submissions.comments", 0)
		fmt.Println("No comments found for this submission")
		return nil
	}

	printVerbose("Found %d comments:\n\n", len(submission.SubmissionComments))
	logger.LogCommandComplete(ctx, "submissions.comments", len(submission.SubmissionComments))
	return formatOutput(submission.SubmissionComments, nil)
}

func runSubmissionsAddComment(ctx context.Context, client *api.Client, opts *options.SubmissionsAddCommentOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "submissions.add-comment", map[string]interface{}{
		"course_id":     opts.CourseID,
		"assignment_id": opts.AssignmentID,
		"user_id":       opts.UserID,
	})

	if _, err := validateCourseID(client, opts.CourseID); err != nil {
		logger.LogCommandError(ctx, "submissions.add-comment", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return err
	}

	submissionsService := api.NewSubmissionsService(client)

	params := &api.GradeSubmissionParams{
		Comment: &api.SubmissionCommentParams{
			TextComment:  opts.Text,
			GroupComment: opts.GroupShare,
		},
	}

	submission, err := submissionsService.Grade(ctx, opts.CourseID, opts.AssignmentID, opts.UserID, params)
	if err != nil {
		logger.LogCommandError(ctx, "submissions.add-comment", err, map[string]interface{}{
			"course_id":     opts.CourseID,
			"assignment_id": opts.AssignmentID,
			"user_id":       opts.UserID,
		})
		return fmt.Errorf("failed to add comment: %w", err)
	}

	logger.LogCommandComplete(ctx, "submissions.add-comment", 1)
	fmt.Printf("Comment added successfully to submission for user %d\n", submission.UserID)
	return nil
}

func runSubmissionsDeleteComment(ctx context.Context, client *api.Client, opts *options.SubmissionsDeleteCommentOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "submissions.delete-comment", map[string]interface{}{
		"course_id":     opts.CourseID,
		"assignment_id": opts.AssignmentID,
		"user_id":       opts.UserID,
		"comment_id":    opts.CommentID,
	})

	if _, err := validateCourseID(client, opts.CourseID); err != nil {
		logger.LogCommandError(ctx, "submissions.delete-comment", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return err
	}

	submissionsService := api.NewSubmissionsService(client)

	_, err := submissionsService.DeleteComment(ctx, opts.CourseID, opts.AssignmentID, opts.UserID, opts.CommentID)
	if err != nil {
		logger.LogCommandError(ctx, "submissions.delete-comment", err, map[string]interface{}{
			"course_id":     opts.CourseID,
			"assignment_id": opts.AssignmentID,
			"user_id":       opts.UserID,
			"comment_id":    opts.CommentID,
		})
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	logger.LogCommandComplete(ctx, "submissions.delete-comment", 1)
	fmt.Printf("Comment %d deleted successfully\n", opts.CommentID)
	return nil
}
