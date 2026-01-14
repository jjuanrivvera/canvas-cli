package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var (
	peerReviewCourseID     int64
	peerReviewAssignmentID int64
	peerReviewSubmissionID int64
	peerReviewUserID       int64
	peerReviewInclude      []string
	peerReviewForce        bool
)

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

var peerReviewsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List peer reviews",
	Long: `List peer reviews for an assignment.

Examples:
  canvas peer-reviews list --course-id 1 --assignment-id 10
  canvas peer-reviews list --course-id 1 --assignment-id 10 --include user,submission_comments`,
	RunE: runPeerReviewsList,
}

var peerReviewsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a peer review",
	Long: `Assign a user as a peer reviewer for a submission.

Examples:
  canvas peer-reviews create --course-id 1 --assignment-id 10 --submission-id 500 --user-id 300`,
	RunE: runPeerReviewsCreate,
}

var peerReviewsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a peer review",
	Long: `Remove a peer review assignment.

Examples:
  canvas peer-reviews delete --course-id 1 --assignment-id 10 --submission-id 500 --user-id 300`,
	RunE: runPeerReviewsDelete,
}

func init() {
	rootCmd.AddCommand(peerReviewsCmd)
	peerReviewsCmd.AddCommand(peerReviewsListCmd)
	peerReviewsCmd.AddCommand(peerReviewsCreateCmd)
	peerReviewsCmd.AddCommand(peerReviewsDeleteCmd)

	// List flags
	peerReviewsListCmd.Flags().Int64Var(&peerReviewCourseID, "course-id", 0, "Course ID (required)")
	peerReviewsListCmd.Flags().Int64Var(&peerReviewAssignmentID, "assignment-id", 0, "Assignment ID (required)")
	peerReviewsListCmd.Flags().StringSliceVar(&peerReviewInclude, "include", nil, "Include options (user, submission_comments)")
	peerReviewsListCmd.MarkFlagRequired("course-id")
	peerReviewsListCmd.MarkFlagRequired("assignment-id")

	// Create flags
	peerReviewsCreateCmd.Flags().Int64Var(&peerReviewCourseID, "course-id", 0, "Course ID (required)")
	peerReviewsCreateCmd.Flags().Int64Var(&peerReviewAssignmentID, "assignment-id", 0, "Assignment ID (required)")
	peerReviewsCreateCmd.Flags().Int64Var(&peerReviewSubmissionID, "submission-id", 0, "Submission ID (required)")
	peerReviewsCreateCmd.Flags().Int64Var(&peerReviewUserID, "user-id", 0, "Reviewer user ID (required)")
	peerReviewsCreateCmd.MarkFlagRequired("course-id")
	peerReviewsCreateCmd.MarkFlagRequired("assignment-id")
	peerReviewsCreateCmd.MarkFlagRequired("submission-id")
	peerReviewsCreateCmd.MarkFlagRequired("user-id")

	// Delete flags
	peerReviewsDeleteCmd.Flags().Int64Var(&peerReviewCourseID, "course-id", 0, "Course ID (required)")
	peerReviewsDeleteCmd.Flags().Int64Var(&peerReviewAssignmentID, "assignment-id", 0, "Assignment ID (required)")
	peerReviewsDeleteCmd.Flags().Int64Var(&peerReviewSubmissionID, "submission-id", 0, "Submission ID (required)")
	peerReviewsDeleteCmd.Flags().Int64Var(&peerReviewUserID, "user-id", 0, "Reviewer user ID (required)")
	peerReviewsDeleteCmd.Flags().BoolVar(&peerReviewForce, "force", false, "Skip confirmation")
	peerReviewsDeleteCmd.MarkFlagRequired("course-id")
	peerReviewsDeleteCmd.MarkFlagRequired("assignment-id")
	peerReviewsDeleteCmd.MarkFlagRequired("submission-id")
	peerReviewsDeleteCmd.MarkFlagRequired("user-id")
}

func runPeerReviewsList(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	var opts *api.ListPeerReviewsOptions
	if len(peerReviewInclude) > 0 {
		opts = &api.ListPeerReviewsOptions{
			Include: peerReviewInclude,
		}
	}

	service := api.NewPeerReviewsService(client)

	ctx := context.Background()
	reviews, err := service.List(ctx, peerReviewCourseID, peerReviewAssignmentID, opts)
	if err != nil {
		return fmt.Errorf("failed to list peer reviews: %w", err)
	}

	if len(reviews) == 0 {
		fmt.Println("No peer reviews found")
		return nil
	}

	printVerbose("Found %d peer reviews:\n\n", len(reviews))
	return formatOutput(reviews, nil)
}

func runPeerReviewsCreate(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	params := &api.CreatePeerReviewParams{
		UserID: peerReviewUserID,
	}

	service := api.NewPeerReviewsService(client)

	ctx := context.Background()
	review, err := service.Create(ctx, peerReviewCourseID, peerReviewAssignmentID, peerReviewSubmissionID, params)
	if err != nil {
		return fmt.Errorf("failed to create peer review: %w", err)
	}

	fmt.Printf("Peer review created (ID: %d)\n", review.ID)
	fmt.Printf("Assessor ID: %d\n", review.AssessorID)
	fmt.Printf("State: %s\n", review.WorkflowState)
	return nil
}

func runPeerReviewsDelete(cmd *cobra.Command, args []string) error {
	if !peerReviewForce {
		fmt.Printf("Are you sure you want to delete peer review for user %d? (y/N): ", peerReviewUserID)
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Cancelled")
			return nil
		}
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewPeerReviewsService(client)

	ctx := context.Background()
	err = service.Delete(ctx, peerReviewCourseID, peerReviewAssignmentID, peerReviewSubmissionID, peerReviewUserID)
	if err != nil {
		return fmt.Errorf("failed to delete peer review: %w", err)
	}

	fmt.Println("Peer review deleted successfully")
	return nil
}
