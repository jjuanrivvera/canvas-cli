package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
	"github.com/spf13/cobra"
)

var (
	assignmentsCourseID   int64
	assignmentsInclude    []string
	assignmentsSearchTerm string
	assignmentsBucket     string
	assignmentsOrderBy    string
)

// assignmentsCmd represents the assignments command group
var assignmentsCmd = &cobra.Command{
	Use:   "assignments",
	Short: "Manage Canvas assignments",
	Long: `Manage Canvas assignments including listing, viewing, creating, and updating assignments.

Examples:
  canvas assignments list --course-id 123
  canvas assignments get --course-id 123 456
  canvas assignments list --course-id 123 --bucket upcoming
  canvas assignments list --course-id 123 --search "quiz"`,
}

// assignmentsListCmd represents the assignments list command
var assignmentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List assignments in a course",
	Long: `List all assignments in a Canvas course.

You can filter assignments by search term, bucket, and order.

Buckets:
  - past: Assignments that are past their due date
  - overdue: Assignments that are overdue for the current user
  - undated: Assignments that have no due date
  - ungraded: Assignments that have not been graded
  - unsubmitted: Assignments that have not been submitted
  - upcoming: Assignments that are due in the future
  - future: Assignments with a due date in the future

Examples:
  canvas assignments list --course-id 123
  canvas assignments list --course-id 123 --bucket upcoming
  canvas assignments list --course-id 123 --search "quiz"
  canvas assignments list --course-id 123 --order-by due_at
  canvas assignments list --course-id 123 --include submission,rubric`,
	RunE: runAssignmentsList,
}

// assignmentsGetCmd represents the assignments get command
var assignmentsGetCmd = &cobra.Command{
	Use:   "get <assignment-id>",
	Short: "Get details of a specific assignment",
	Long: `Get details of a specific assignment by ID.

Examples:
  canvas assignments get --course-id 123 456
  canvas assignments get --course-id 123 456 --include submission,rubric`,
	Args: cobra.ExactArgs(1),
	RunE: runAssignmentsGet,
}

func init() {
	rootCmd.AddCommand(assignmentsCmd)
	assignmentsCmd.AddCommand(assignmentsListCmd)
	assignmentsCmd.AddCommand(assignmentsGetCmd)

	// List flags
	assignmentsListCmd.Flags().Int64Var(&assignmentsCourseID, "course-id", 0, "Course ID (required)")
	assignmentsListCmd.Flags().StringVar(&assignmentsSearchTerm, "search", "", "Search by assignment name")
	assignmentsListCmd.Flags().StringVar(&assignmentsBucket, "bucket", "", "Filter by bucket (past, overdue, undated, ungraded, unsubmitted, upcoming, future)")
	assignmentsListCmd.Flags().StringVar(&assignmentsOrderBy, "order-by", "", "Order by (position, name, due_at)")
	assignmentsListCmd.Flags().StringSliceVar(&assignmentsInclude, "include", []string{}, "Additional data to include (comma-separated)")
	assignmentsListCmd.MarkFlagRequired("course-id")

	// Get flags
	assignmentsGetCmd.Flags().Int64Var(&assignmentsCourseID, "course-id", 0, "Course ID (required)")
	assignmentsGetCmd.Flags().StringSliceVar(&assignmentsInclude, "include", []string{}, "Additional data to include (comma-separated)")
	assignmentsGetCmd.MarkFlagRequired("course-id")
}

func runAssignmentsList(cmd *cobra.Command, args []string) error {
	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create assignments service
	assignmentsService := api.NewAssignmentsService(client)

	// Build options
	opts := &api.ListAssignmentsOptions{
		SearchTerm: assignmentsSearchTerm,
		Bucket:     assignmentsBucket,
		OrderBy:    assignmentsOrderBy,
		Include:    assignmentsInclude,
	}

	// List assignments
	ctx := context.Background()
	assignments, err := assignmentsService.List(ctx, assignmentsCourseID, opts)
	if err != nil {
		return fmt.Errorf("failed to list assignments: %w", err)
	}

	if len(assignments) == 0 {
		fmt.Println("No assignments found")
		return nil
	}

	// Display assignments
	fmt.Printf("Found %d assignments:\n\n", len(assignments))

	for _, assignment := range assignments {
		fmt.Printf("📝 %s\n", assignment.Name)
		fmt.Printf("   ID: %d\n", assignment.ID)
		fmt.Printf("   Points: %.1f\n", assignment.PointsPossible)

		if !assignment.DueAt.IsZero() {
			fmt.Printf("   Due: %s\n", assignment.DueAt.Format("2006-01-02 15:04"))
		} else {
			fmt.Printf("   Due: No due date\n")
		}

		if len(assignment.SubmissionTypes) > 0 {
			fmt.Printf("   Types: %s\n", strings.Join(assignment.SubmissionTypes, ", "))
		}

		if assignment.Published {
			fmt.Printf("   Status: Published\n")
		} else {
			fmt.Printf("   Status: Unpublished\n")
		}

		if assignment.NeedsGradingCount > 0 {
			fmt.Printf("   Needs Grading: %d\n", assignment.NeedsGradingCount)
		}

		fmt.Println()
	}

	return nil
}

func runAssignmentsGet(cmd *cobra.Command, args []string) error {
	// Parse assignment ID
	assignmentID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid assignment ID: %w", err)
	}

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create assignments service
	assignmentsService := api.NewAssignmentsService(client)

	// Get assignment
	ctx := context.Background()
	assignment, err := assignmentsService.Get(ctx, assignmentsCourseID, assignmentID, assignmentsInclude)
	if err != nil {
		return fmt.Errorf("failed to get assignment: %w", err)
	}

	// Display assignment details
	fmt.Printf("📝 %s\n", assignment.Name)
	fmt.Printf("   ID: %d\n", assignment.ID)
	fmt.Printf("   Course ID: %d\n", assignment.CourseID)
	fmt.Printf("   Points Possible: %.1f\n", assignment.PointsPossible)

	if assignment.GradingType != "" {
		fmt.Printf("   Grading Type: %s\n", assignment.GradingType)
	}

	if len(assignment.SubmissionTypes) > 0 {
		fmt.Printf("   Submission Types: %s\n", strings.Join(assignment.SubmissionTypes, ", "))
	}

	if !assignment.DueAt.IsZero() {
		fmt.Printf("   Due: %s\n", assignment.DueAt.Format("2006-01-02 15:04"))
	} else {
		fmt.Printf("   Due: No due date\n")
	}

	if !assignment.UnlockAt.IsZero() {
		fmt.Printf("   Available From: %s\n", assignment.UnlockAt.Format("2006-01-02 15:04"))
	}

	if !assignment.LockAt.IsZero() {
		fmt.Printf("   Available Until: %s\n", assignment.LockAt.Format("2006-01-02 15:04"))
	}

	if assignment.Published {
		fmt.Printf("   Status: Published\n")
	} else {
		fmt.Printf("   Status: Unpublished\n")
	}

	if assignment.NeedsGradingCount > 0 {
		fmt.Printf("   Needs Grading: %d submissions\n", assignment.NeedsGradingCount)
	}

	if assignment.Description != "" {
		fmt.Printf("\nDescription:\n%s\n", assignment.Description)
	}

	if len(assignment.Rubric) > 0 {
		fmt.Printf("\nRubric: %d criteria\n", len(assignment.Rubric))
	}

	if len(assignment.Overrides) > 0 {
		fmt.Printf("\nOverrides: %d\n", len(assignment.Overrides))
		for i, override := range assignment.Overrides {
			fmt.Printf("  %d. %s\n", i+1, override.Title)
			if !override.DueAt.IsZero() {
				fmt.Printf("     Due: %s\n", override.DueAt.Format("2006-01-02 15:04"))
			}
		}
	}

	return nil
}
