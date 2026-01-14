package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
	"github.com/jjuanrivvera/canvas-cli/internal/batch"
)

var (
	submissionsCourseID      int64
	submissionsAssignmentID  int64
	submissionsUserID        int64
	submissionsInclude       []string
	submissionsWorkflowState string
	submissionsGradedSince   string

	// Grade flags
	gradeScore       float64
	gradeComment     string
	gradeExcuse      bool
	gradePostedGrade string

	// Bulk grade flags
	bulkGradeCSV    string
	bulkGradeDryRun bool

	// Comment flags
	commentText       string
	commentGroupShare bool
	commentID         int64
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

// submissionsListCmd represents the submissions list command
var submissionsListCmd = &cobra.Command{
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
	RunE: runSubmissionsList,
}

// submissionsGetCmd represents the submissions get command
var submissionsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a specific submission",
	Long: `Get details of a specific submission for an assignment and user.

Examples:
  canvas submissions get --course-id 123 --assignment-id 456 --user-id 789
  canvas submissions get --course-id 123 --assignment-id 456 --user-id 789 --include submission_comments,rubric_assessment`,
	RunE: runSubmissionsGet,
}

// submissionsGradeCmd represents the submissions grade command
var submissionsGradeCmd = &cobra.Command{
	Use:   "grade",
	Short: "Grade a submission",
	Long: `Grade a specific submission for an assignment and user.

You can provide a score, comment, or excuse the submission.

Examples:
  canvas submissions grade --course-id 123 --assignment-id 456 --user-id 789 --score 95
  canvas submissions grade --course-id 123 --assignment-id 456 --user-id 789 --score 85 --comment "Good work"
  canvas submissions grade --course-id 123 --assignment-id 456 --user-id 789 --excuse`,
	RunE: runSubmissionsGrade,
}

// submissionsBulkGradeCmd represents the submissions bulk-grade command
var submissionsBulkGradeCmd = &cobra.Command{
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
	RunE: runSubmissionsBulkGrade,
}

// submissionsCommentsCmd represents the submissions comments command
var submissionsCommentsCmd = &cobra.Command{
	Use:   "comments",
	Short: "List comments for a submission",
	Long: `List all comments for a specific submission.

Examples:
  canvas submissions comments --course-id 123 --assignment-id 456 --user-id 789`,
	RunE: runSubmissionsComments,
}

// submissionsAddCommentCmd represents the submissions add-comment command
var submissionsAddCommentCmd = &cobra.Command{
	Use:   "add-comment",
	Short: "Add a comment to a submission",
	Long: `Add a comment to a specific submission.

Examples:
  canvas submissions add-comment --course-id 123 --assignment-id 456 --user-id 789 --text "Great work!"
  canvas submissions add-comment --course-id 123 --assignment-id 456 --user-id 789 --text "Feedback" --group`,
	RunE: runSubmissionsAddComment,
}

// submissionsDeleteCommentCmd represents the submissions delete-comment command
var submissionsDeleteCommentCmd = &cobra.Command{
	Use:   "delete-comment",
	Short: "Delete a submission comment",
	Long: `Delete a comment from a submission.

Examples:
  canvas submissions delete-comment --course-id 123 --assignment-id 456 --user-id 789 --comment-id 999`,
	RunE: runSubmissionsDeleteComment,
}

func init() {
	rootCmd.AddCommand(submissionsCmd)
	submissionsCmd.AddCommand(submissionsListCmd)
	submissionsCmd.AddCommand(submissionsGetCmd)
	submissionsCmd.AddCommand(submissionsGradeCmd)
	submissionsCmd.AddCommand(submissionsBulkGradeCmd)
	submissionsCmd.AddCommand(submissionsCommentsCmd)
	submissionsCmd.AddCommand(submissionsAddCommentCmd)
	submissionsCmd.AddCommand(submissionsDeleteCommentCmd)

	// List flags
	submissionsListCmd.Flags().Int64Var(&submissionsCourseID, "course-id", 0, "Course ID (required)")
	submissionsListCmd.Flags().Int64Var(&submissionsAssignmentID, "assignment-id", 0, "Assignment ID (required)")
	submissionsListCmd.Flags().StringVar(&submissionsWorkflowState, "workflow-state", "", "Filter by workflow state (submitted, unsubmitted, graded, pending_review)")
	submissionsListCmd.Flags().StringVar(&submissionsGradedSince, "graded-since", "", "Filter by graded since date (ISO8601 format)")
	submissionsListCmd.Flags().StringSliceVar(&submissionsInclude, "include", []string{}, "Additional data to include (comma-separated)")
	submissionsListCmd.MarkFlagRequired("course-id")
	submissionsListCmd.MarkFlagRequired("assignment-id")

	// Get flags
	submissionsGetCmd.Flags().Int64Var(&submissionsCourseID, "course-id", 0, "Course ID (required)")
	submissionsGetCmd.Flags().Int64Var(&submissionsAssignmentID, "assignment-id", 0, "Assignment ID (required)")
	submissionsGetCmd.Flags().Int64Var(&submissionsUserID, "user-id", 0, "User ID (required)")
	submissionsGetCmd.Flags().StringSliceVar(&submissionsInclude, "include", []string{}, "Additional data to include (comma-separated)")
	submissionsGetCmd.MarkFlagRequired("course-id")
	submissionsGetCmd.MarkFlagRequired("assignment-id")
	submissionsGetCmd.MarkFlagRequired("user-id")

	// Grade flags
	submissionsGradeCmd.Flags().Int64Var(&submissionsCourseID, "course-id", 0, "Course ID (required)")
	submissionsGradeCmd.Flags().Int64Var(&submissionsAssignmentID, "assignment-id", 0, "Assignment ID (required)")
	submissionsGradeCmd.Flags().Int64Var(&submissionsUserID, "user-id", 0, "User ID (required)")
	submissionsGradeCmd.Flags().Float64Var(&gradeScore, "score", 0, "Score to assign")
	submissionsGradeCmd.Flags().StringVar(&gradeComment, "comment", "", "Comment to add")
	submissionsGradeCmd.Flags().BoolVar(&gradeExcuse, "excuse", false, "Excuse the submission")
	submissionsGradeCmd.Flags().StringVar(&gradePostedGrade, "posted-grade", "", "Posted grade (e.g., 'A', 'B+', 'Pass')")
	submissionsGradeCmd.MarkFlagRequired("course-id")
	submissionsGradeCmd.MarkFlagRequired("assignment-id")
	submissionsGradeCmd.MarkFlagRequired("user-id")

	// Bulk grade flags
	submissionsBulkGradeCmd.Flags().Int64Var(&submissionsCourseID, "course-id", 0, "Course ID (required)")
	submissionsBulkGradeCmd.Flags().StringVar(&bulkGradeCSV, "csv", "", "CSV file with grades (required)")
	submissionsBulkGradeCmd.Flags().BoolVar(&bulkGradeDryRun, "dry-run", false, "Preview changes without applying them")
	submissionsBulkGradeCmd.MarkFlagRequired("course-id")
	submissionsBulkGradeCmd.MarkFlagRequired("csv")

	// Comments flags
	submissionsCommentsCmd.Flags().Int64Var(&submissionsCourseID, "course-id", 0, "Course ID (required)")
	submissionsCommentsCmd.Flags().Int64Var(&submissionsAssignmentID, "assignment-id", 0, "Assignment ID (required)")
	submissionsCommentsCmd.Flags().Int64Var(&submissionsUserID, "user-id", 0, "User ID (required)")
	submissionsCommentsCmd.MarkFlagRequired("course-id")
	submissionsCommentsCmd.MarkFlagRequired("assignment-id")
	submissionsCommentsCmd.MarkFlagRequired("user-id")

	// Add comment flags
	submissionsAddCommentCmd.Flags().Int64Var(&submissionsCourseID, "course-id", 0, "Course ID (required)")
	submissionsAddCommentCmd.Flags().Int64Var(&submissionsAssignmentID, "assignment-id", 0, "Assignment ID (required)")
	submissionsAddCommentCmd.Flags().Int64Var(&submissionsUserID, "user-id", 0, "User ID (required)")
	submissionsAddCommentCmd.Flags().StringVar(&commentText, "text", "", "Comment text (required)")
	submissionsAddCommentCmd.Flags().BoolVar(&commentGroupShare, "group", false, "Share comment with group members")
	submissionsAddCommentCmd.MarkFlagRequired("course-id")
	submissionsAddCommentCmd.MarkFlagRequired("assignment-id")
	submissionsAddCommentCmd.MarkFlagRequired("user-id")
	submissionsAddCommentCmd.MarkFlagRequired("text")

	// Delete comment flags
	submissionsDeleteCommentCmd.Flags().Int64Var(&submissionsCourseID, "course-id", 0, "Course ID (required)")
	submissionsDeleteCommentCmd.Flags().Int64Var(&submissionsAssignmentID, "assignment-id", 0, "Assignment ID (required)")
	submissionsDeleteCommentCmd.Flags().Int64Var(&submissionsUserID, "user-id", 0, "User ID (required)")
	submissionsDeleteCommentCmd.Flags().Int64Var(&commentID, "comment-id", 0, "Comment ID to delete (required)")
	submissionsDeleteCommentCmd.MarkFlagRequired("course-id")
	submissionsDeleteCommentCmd.MarkFlagRequired("assignment-id")
	submissionsDeleteCommentCmd.MarkFlagRequired("user-id")
	submissionsDeleteCommentCmd.MarkFlagRequired("comment-id")
}

func runSubmissionsList(cmd *cobra.Command, args []string) error {
	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Validate course ID exists
	if _, err := validateCourseID(client, submissionsCourseID); err != nil {
		return err
	}

	// Create submissions service
	submissionsService := api.NewSubmissionsService(client)

	// Build options
	opts := &api.ListSubmissionsOptions{
		WorkflowState: submissionsWorkflowState,
		GradedSince:   submissionsGradedSince,
		Include:       submissionsInclude,
	}

	// List submissions
	ctx := context.Background()
	submissions, err := submissionsService.List(ctx, submissionsCourseID, submissionsAssignmentID, opts)
	if err != nil {
		return fmt.Errorf("failed to list submissions: %w", err)
	}

	if len(submissions) == 0 {
		fmt.Println("No submissions found")
		return nil
	}

	// Format and display submissions
	printVerbose("Found %d submissions:\n\n", len(submissions))
	return formatOutput(submissions, nil)
}

func runSubmissionsGet(cmd *cobra.Command, args []string) error {
	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Validate course ID exists
	if _, err := validateCourseID(client, submissionsCourseID); err != nil {
		return err
	}

	// Create submissions service
	submissionsService := api.NewSubmissionsService(client)

	// Get submission
	ctx := context.Background()
	submission, err := submissionsService.Get(ctx, submissionsCourseID, submissionsAssignmentID, submissionsUserID, submissionsInclude)
	if err != nil {
		return fmt.Errorf("failed to get submission: %w", err)
	}

	// Format and display submission details
	return formatOutput(submission, nil)
}

func runSubmissionsGrade(cmd *cobra.Command, args []string) error {
	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Validate course ID exists
	if _, err := validateCourseID(client, submissionsCourseID); err != nil {
		return err
	}

	// Create submissions service
	submissionsService := api.NewSubmissionsService(client)

	// Build grade params
	params := &api.GradeSubmissionParams{}

	// Handle score (convert to string for PostedGrade)
	if gradeScore > 0 {
		params.PostedGrade = fmt.Sprintf("%.2f", gradeScore)
	} else if gradePostedGrade != "" {
		params.PostedGrade = gradePostedGrade
	}

	// Handle comment
	if gradeComment != "" {
		params.Comment = &api.SubmissionCommentParams{
			TextComment: gradeComment,
		}
	}

	// Handle excuse
	if gradeExcuse {
		params.Excuse = true
	}

	// Validate at least one grading parameter is provided
	if params.PostedGrade == "" && params.Comment == nil && !params.Excuse {
		return fmt.Errorf("at least one grading parameter is required: --score, --comment, --excuse, or --posted-grade")
	}

	// Grade submission
	ctx := context.Background()
	submission, err := submissionsService.Grade(ctx, submissionsCourseID, submissionsAssignmentID, submissionsUserID, params)
	if err != nil {
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

	return nil
}

func runSubmissionsBulkGrade(cmd *cobra.Command, args []string) error {
	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Validate course ID exists
	if _, err := validateCourseID(client, submissionsCourseID); err != nil {
		return err
	}

	// Create submissions service
	submissionsService := api.NewSubmissionsService(client)

	// Read grades from CSV
	grades, err := batch.ReadGradesCSV(bulkGradeCSV)
	if err != nil {
		return fmt.Errorf("failed to read CSV file: %w", err)
	}

	if len(grades) == 0 {
		return fmt.Errorf("no grades found in CSV file")
	}

	printVerbose("Found %d grades in CSV file\n\n", len(grades))

	if bulkGradeDryRun {
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
		return nil
	}

	// Process grades
	ctx := context.Background()
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
		_, err = submissionsService.Grade(ctx, submissionsCourseID, grade.AssignmentID, grade.UserID, params)
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

	if errorCount > 0 {
		return fmt.Errorf("bulk grading completed with %d errors", errorCount)
	}

	return nil
}

func runSubmissionsComments(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	if _, err := validateCourseID(client, submissionsCourseID); err != nil {
		return err
	}

	submissionsService := api.NewSubmissionsService(client)

	ctx := context.Background()
	submission, err := submissionsService.Get(ctx, submissionsCourseID, submissionsAssignmentID, submissionsUserID, []string{"submission_comments"})
	if err != nil {
		return fmt.Errorf("failed to get submission: %w", err)
	}

	if len(submission.SubmissionComments) == 0 {
		fmt.Println("No comments found for this submission")
		return nil
	}

	printVerbose("Found %d comments:\n\n", len(submission.SubmissionComments))
	return formatOutput(submission.SubmissionComments, nil)
}

func runSubmissionsAddComment(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	if _, err := validateCourseID(client, submissionsCourseID); err != nil {
		return err
	}

	submissionsService := api.NewSubmissionsService(client)

	params := &api.GradeSubmissionParams{
		Comment: &api.SubmissionCommentParams{
			TextComment:  commentText,
			GroupComment: commentGroupShare,
		},
	}

	ctx := context.Background()
	submission, err := submissionsService.Grade(ctx, submissionsCourseID, submissionsAssignmentID, submissionsUserID, params)
	if err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}

	fmt.Printf("Comment added successfully to submission for user %d\n", submission.UserID)
	return nil
}

func runSubmissionsDeleteComment(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	if _, err := validateCourseID(client, submissionsCourseID); err != nil {
		return err
	}

	submissionsService := api.NewSubmissionsService(client)

	ctx := context.Background()
	_, err = submissionsService.DeleteComment(ctx, submissionsCourseID, submissionsAssignmentID, submissionsUserID, commentID)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	fmt.Printf("Comment %d deleted successfully\n", commentID)
	return nil
}
