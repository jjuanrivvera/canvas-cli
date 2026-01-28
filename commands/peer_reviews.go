package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// peerReviewsCmd represents the peer-reviews command group
var peerReviewsCmd = &cobra.Command{
	Use:     "peer-reviews",
	Aliases: []string{"pr"},
	Short:   "Manage peer reviews",
	Long: `Manage Canvas peer reviews for assignments.

Peer reviews allow students to review and provide feedback
on each other's work.

Examples:
  canvas peer-reviews list --course-id 1 --assignment-id 10
  canvas peer-reviews create --course-id 1 --assignment-id 10 --submission-id 500 --user-id 300`,
}

func init() {
	rootCmd.AddCommand(peerReviewsCmd)
	peerReviewsCmd.AddCommand(newPeerReviewsListCmd())
	peerReviewsCmd.AddCommand(newPeerReviewsCreateCmd())
	peerReviewsCmd.AddCommand(newPeerReviewsDeleteCmd())
}

func newPeerReviewsListCmd() *cobra.Command {
	opts := &options.PeerReviewsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List peer reviews",
		Long: `List peer reviews for an assignment.

Examples:
  canvas peer-reviews list --course-id 1 --assignment-id 10
  canvas peer-reviews list --course-id 1 --assignment-id 10 --include user,submission_comments`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPeerReviewsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.AssignmentID, "assignment-id", 0, "Assignment ID (required)")
	cmd.Flags().StringSliceVar(&opts.Include, "include", nil, "Include options (user, submission_comments)")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("assignment-id")

	return cmd
}

func newPeerReviewsCreateCmd() *cobra.Command {
	opts := &options.PeerReviewsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a peer review",
		Long: `Assign a user as a peer reviewer for a submission.

Examples:
  canvas peer-reviews create --course-id 1 --assignment-id 10 --submission-id 500 --user-id 300`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPeerReviewsCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.AssignmentID, "assignment-id", 0, "Assignment ID (required)")
	cmd.Flags().Int64Var(&opts.SubmissionID, "submission-id", 0, "Submission ID (required)")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "Reviewer user ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("assignment-id")
	cmd.MarkFlagRequired("submission-id")
	cmd.MarkFlagRequired("user-id")

	return cmd
}

func newPeerReviewsDeleteCmd() *cobra.Command {
	opts := &options.PeerReviewsDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a peer review",
		Long: `Remove a peer review assignment.

Examples:
  canvas peer-reviews delete --course-id 1 --assignment-id 10 --submission-id 500 --user-id 300`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			// Confirm deletion
			confirmed, err := confirmDelete("peer review", opts.UserID, opts.Force)
			if err != nil {
				return err
			}
			if !confirmed {
				fmt.Println("Delete cancelled")
				return nil
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPeerReviewsDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.AssignmentID, "assignment-id", 0, "Assignment ID (required)")
	cmd.Flags().Int64Var(&opts.SubmissionID, "submission-id", 0, "Submission ID (required)")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "Reviewer user ID (required)")
	cmd.Flags().BoolVar(&opts.Force, "force", false, "Skip confirmation")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("assignment-id")
	cmd.MarkFlagRequired("submission-id")
	cmd.MarkFlagRequired("user-id")

	return cmd
}

func runPeerReviewsList(ctx context.Context, client *api.Client, opts *options.PeerReviewsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "peer_reviews.list", map[string]interface{}{
		"course_id":     opts.CourseID,
		"assignment_id": opts.AssignmentID,
	})

	var apiOpts *api.ListPeerReviewsOptions
	if len(opts.Include) > 0 {
		apiOpts = &api.ListPeerReviewsOptions{
			Include: opts.Include,
		}
	}

	service := api.NewPeerReviewsService(client)

	reviews, err := service.List(ctx, opts.CourseID, opts.AssignmentID, apiOpts)
	if err != nil {
		logger.LogCommandError(ctx, "peer_reviews.list", err, map[string]interface{}{
			"course_id":     opts.CourseID,
			"assignment_id": opts.AssignmentID,
		})
		return fmt.Errorf("failed to list peer reviews: %w", err)
	}

	if len(reviews) == 0 {
		fmt.Println("No peer reviews found")
		logger.LogCommandComplete(ctx, "peer_reviews.list", 0)
		return nil
	}

	printVerbose("Found %d peer reviews:\n\n", len(reviews))
	logger.LogCommandComplete(ctx, "peer_reviews.list", len(reviews))
	return formatOutput(reviews, nil)
}

func runPeerReviewsCreate(ctx context.Context, client *api.Client, opts *options.PeerReviewsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "peer_reviews.create", map[string]interface{}{
		"course_id":     opts.CourseID,
		"assignment_id": opts.AssignmentID,
		"submission_id": opts.SubmissionID,
		"user_id":       opts.UserID,
	})

	params := &api.CreatePeerReviewParams{
		UserID: opts.UserID,
	}

	service := api.NewPeerReviewsService(client)

	review, err := service.Create(ctx, opts.CourseID, opts.AssignmentID, opts.SubmissionID, params)
	if err != nil {
		logger.LogCommandError(ctx, "peer_reviews.create", err, map[string]interface{}{
			"course_id":     opts.CourseID,
			"assignment_id": opts.AssignmentID,
			"submission_id": opts.SubmissionID,
			"user_id":       opts.UserID,
		})
		return fmt.Errorf("failed to create peer review: %w", err)
	}

	fmt.Printf("Peer review created (ID: %d)\n", review.ID)
	fmt.Printf("Assessor ID: %d\n", review.AssessorID)
	fmt.Printf("State: %s\n", review.WorkflowState)
	logger.LogCommandComplete(ctx, "peer_reviews.create", 1)
	return nil
}

func runPeerReviewsDelete(ctx context.Context, client *api.Client, opts *options.PeerReviewsDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "peer_reviews.delete", map[string]interface{}{
		"course_id":     opts.CourseID,
		"assignment_id": opts.AssignmentID,
		"submission_id": opts.SubmissionID,
		"user_id":       opts.UserID,
	})

	service := api.NewPeerReviewsService(client)

	err := service.Delete(ctx, opts.CourseID, opts.AssignmentID, opts.SubmissionID, opts.UserID)
	if err != nil {
		logger.LogCommandError(ctx, "peer_reviews.delete", err, map[string]interface{}{
			"course_id":     opts.CourseID,
			"assignment_id": opts.AssignmentID,
			"submission_id": opts.SubmissionID,
			"user_id":       opts.UserID,
		})
		return fmt.Errorf("failed to delete peer review: %w", err)
	}

	printInfoln("Peer review deleted successfully")
	logger.LogCommandComplete(ctx, "peer_reviews.delete", 1)
	return nil
}
