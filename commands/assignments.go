package commands

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var (
	assignmentsCourseID   int64
	assignmentsInclude    []string
	assignmentsSearchTerm string
	assignmentsBucket     string
	assignmentsOrderBy    string

	// Create/Update flags
	assignmentName            string
	assignmentPoints          float64
	assignmentGradingType     string
	assignmentDueAt           string
	assignmentUnlockAt        string
	assignmentLockAt          string
	assignmentDescription     string
	assignmentPublished       bool
	assignmentSubmissionTypes []string
	assignmentGroupID         int64
	assignmentPosition        int
	assignmentJSONFile        string
	assignmentStdin           bool
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

// assignmentsCreateCmd represents the assignments create command
var assignmentsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new assignment",
	Long: `Create a new assignment in a Canvas course.

You can provide assignment data via flags or JSON file/stdin.

Examples:
  # Using flags
  canvas assignments create --course-id 123 --name "Quiz 1" --points 100
  canvas assignments create --course-id 123 --name "Essay" --points 50 --due-at "2025-02-01T23:59:00Z"

  # Using JSON file
  canvas assignments create --course-id 123 --json assignment.json

  # Using stdin
  echo '{"name":"Quiz 1","points_possible":100}' | canvas assignments create --course-id 123 --stdin`,
	RunE: runAssignmentsCreate,
}

// assignmentsUpdateCmd represents the assignments update command
var assignmentsUpdateCmd = &cobra.Command{
	Use:   "update <assignment-id>",
	Short: "Update an existing assignment",
	Long: `Update an existing assignment in a Canvas course.

You can provide assignment data via flags or JSON file/stdin.
Only specified fields will be updated.

Examples:
  # Using flags
  canvas assignments update --course-id 123 456 --name "Updated Quiz"
  canvas assignments update --course-id 123 456 --points 150 --due-at "2025-02-15T23:59:00Z"

  # Using JSON file
  canvas assignments update --course-id 123 456 --json updates.json

  # Using stdin
  echo '{"points_possible":200}' | canvas assignments update --course-id 123 456 --stdin`,
	Args: cobra.ExactArgs(1),
	RunE: runAssignmentsUpdate,
}

// assignmentsDeleteCmd represents the assignments delete command
var assignmentsDeleteCmd = &cobra.Command{
	Use:   "delete <assignment-id>",
	Short: "Delete an assignment",
	Long: `Delete an assignment from a Canvas course.

This action cannot be undone.

Examples:
  canvas assignments delete --course-id 123 456`,
	Args: cobra.ExactArgs(1),
	RunE: runAssignmentsDelete,
}

func init() {
	rootCmd.AddCommand(assignmentsCmd)
	assignmentsCmd.AddCommand(assignmentsListCmd)
	assignmentsCmd.AddCommand(assignmentsGetCmd)
	assignmentsCmd.AddCommand(assignmentsCreateCmd)
	assignmentsCmd.AddCommand(assignmentsUpdateCmd)
	assignmentsCmd.AddCommand(assignmentsDeleteCmd)

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

	// Create flags
	assignmentsCreateCmd.Flags().Int64Var(&assignmentsCourseID, "course-id", 0, "Course ID (required)")
	assignmentsCreateCmd.Flags().StringVar(&assignmentName, "name", "", "Assignment name")
	assignmentsCreateCmd.Flags().Float64Var(&assignmentPoints, "points", 0, "Points possible")
	assignmentsCreateCmd.Flags().StringVar(&assignmentGradingType, "grading-type", "", "Grading type (points, pass_fail, percent, letter_grade, gpa_scale, not_graded)")
	assignmentsCreateCmd.Flags().StringVar(&assignmentDueAt, "due-at", "", "Due date (ISO8601 format)")
	assignmentsCreateCmd.Flags().StringVar(&assignmentUnlockAt, "unlock-at", "", "Unlock date (ISO8601 format)")
	assignmentsCreateCmd.Flags().StringVar(&assignmentLockAt, "lock-at", "", "Lock date (ISO8601 format)")
	assignmentsCreateCmd.Flags().StringVar(&assignmentDescription, "description", "", "Assignment description (HTML)")
	assignmentsCreateCmd.Flags().BoolVar(&assignmentPublished, "published", false, "Publish the assignment")
	assignmentsCreateCmd.Flags().StringSliceVar(&assignmentSubmissionTypes, "submission-types", []string{}, "Submission types (online_text_entry, online_url, online_upload, media_recording, none)")
	assignmentsCreateCmd.Flags().Int64Var(&assignmentGroupID, "group-id", 0, "Assignment group ID")
	assignmentsCreateCmd.Flags().IntVar(&assignmentPosition, "position", 0, "Position in the assignment group")
	assignmentsCreateCmd.Flags().StringVar(&assignmentJSONFile, "json", "", "JSON file with assignment data")
	assignmentsCreateCmd.Flags().BoolVar(&assignmentStdin, "stdin", false, "Read JSON from stdin")
	assignmentsCreateCmd.MarkFlagRequired("course-id")

	// Update flags
	assignmentsUpdateCmd.Flags().Int64Var(&assignmentsCourseID, "course-id", 0, "Course ID (required)")
	assignmentsUpdateCmd.Flags().StringVar(&assignmentName, "name", "", "Assignment name")
	assignmentsUpdateCmd.Flags().Float64Var(&assignmentPoints, "points", 0, "Points possible")
	assignmentsUpdateCmd.Flags().StringVar(&assignmentGradingType, "grading-type", "", "Grading type (points, pass_fail, percent, letter_grade, gpa_scale, not_graded)")
	assignmentsUpdateCmd.Flags().StringVar(&assignmentDueAt, "due-at", "", "Due date (ISO8601 format)")
	assignmentsUpdateCmd.Flags().StringVar(&assignmentUnlockAt, "unlock-at", "", "Unlock date (ISO8601 format)")
	assignmentsUpdateCmd.Flags().StringVar(&assignmentLockAt, "lock-at", "", "Lock date (ISO8601 format)")
	assignmentsUpdateCmd.Flags().StringVar(&assignmentDescription, "description", "", "Assignment description (HTML)")
	assignmentsUpdateCmd.Flags().BoolVar(&assignmentPublished, "published", false, "Publish the assignment")
	assignmentsUpdateCmd.Flags().StringSliceVar(&assignmentSubmissionTypes, "submission-types", []string{}, "Submission types")
	assignmentsUpdateCmd.Flags().Int64Var(&assignmentGroupID, "group-id", 0, "Assignment group ID")
	assignmentsUpdateCmd.Flags().IntVar(&assignmentPosition, "position", 0, "Position in the assignment group")
	assignmentsUpdateCmd.Flags().StringVar(&assignmentJSONFile, "json", "", "JSON file with assignment data")
	assignmentsUpdateCmd.Flags().BoolVar(&assignmentStdin, "stdin", false, "Read JSON from stdin")
	assignmentsUpdateCmd.MarkFlagRequired("course-id")

	// Delete flags
	assignmentsDeleteCmd.Flags().Int64Var(&assignmentsCourseID, "course-id", 0, "Course ID (required)")
	assignmentsDeleteCmd.MarkFlagRequired("course-id")
}

func runAssignmentsList(cmd *cobra.Command, args []string) error {
	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Validate course ID exists
	if _, err := validateCourseID(client, assignmentsCourseID); err != nil {
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

	// Format and display assignments
	return formatOutput(assignments, func() {
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
	})
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

	// Validate course ID exists
	if _, err := validateCourseID(client, assignmentsCourseID); err != nil {
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

	// Format and display assignment details
	return formatOutput(assignment, func() {
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
	})
}

func runAssignmentsCreate(cmd *cobra.Command, args []string) error {
	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Validate course ID exists
	if _, err := validateCourseID(client, assignmentsCourseID); err != nil {
		return err
	}

	// Create assignments service
	assignmentsService := api.NewAssignmentsService(client)

	// Build params from flags or JSON
	params := &api.CreateAssignmentParams{}

	// Check for JSON input
	if assignmentJSONFile != "" || assignmentStdin {
		jsonData, err := readAssignmentJSON(assignmentJSONFile, assignmentStdin)
		if err != nil {
			return fmt.Errorf("failed to read JSON: %w", err)
		}
		if err := parseAssignmentCreateJSON(jsonData, params); err != nil {
			return fmt.Errorf("failed to parse JSON: %w", err)
		}
	}

	// Override with flags if provided
	if assignmentName != "" {
		params.Name = assignmentName
	}
	if assignmentPoints > 0 {
		params.PointsPossible = assignmentPoints
	}
	if assignmentGradingType != "" {
		params.GradingType = assignmentGradingType
	}
	if assignmentDueAt != "" {
		params.DueAt = assignmentDueAt
	}
	if assignmentUnlockAt != "" {
		params.UnlockAt = assignmentUnlockAt
	}
	if assignmentLockAt != "" {
		params.LockAt = assignmentLockAt
	}
	if assignmentDescription != "" {
		params.Description = assignmentDescription
	}
	if cmd.Flags().Changed("published") {
		params.Published = assignmentPublished
	}
	if len(assignmentSubmissionTypes) > 0 {
		params.SubmissionTypes = assignmentSubmissionTypes
	}
	if assignmentGroupID > 0 {
		params.AssignmentGroupID = assignmentGroupID
	}
	if assignmentPosition > 0 {
		params.Position = assignmentPosition
	}

	// Validate required fields
	if params.Name == "" {
		return fmt.Errorf("assignment name is required (use --name or provide in JSON)")
	}

	// Create assignment
	ctx := context.Background()
	assignment, err := assignmentsService.Create(ctx, assignmentsCourseID, params)
	if err != nil {
		return fmt.Errorf("failed to create assignment: %w", err)
	}

	fmt.Printf("Assignment created successfully!\n")
	fmt.Printf("  ID: %d\n", assignment.ID)
	fmt.Printf("  Name: %s\n", assignment.Name)
	fmt.Printf("  Points: %.1f\n", assignment.PointsPossible)
	if !assignment.DueAt.IsZero() {
		fmt.Printf("  Due: %s\n", assignment.DueAt.Format("2006-01-02 15:04"))
	}
	if assignment.Published {
		fmt.Printf("  Status: Published\n")
	} else {
		fmt.Printf("  Status: Unpublished\n")
	}

	return nil
}

func runAssignmentsUpdate(cmd *cobra.Command, args []string) error {
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

	// Validate course ID exists
	if _, err := validateCourseID(client, assignmentsCourseID); err != nil {
		return err
	}

	// Create assignments service
	assignmentsService := api.NewAssignmentsService(client)

	// Build params from flags or JSON
	params := &api.UpdateAssignmentParams{}

	// Check for JSON input
	if assignmentJSONFile != "" || assignmentStdin {
		jsonData, err := readAssignmentJSON(assignmentJSONFile, assignmentStdin)
		if err != nil {
			return fmt.Errorf("failed to read JSON: %w", err)
		}
		if err := parseAssignmentUpdateJSON(jsonData, params); err != nil {
			return fmt.Errorf("failed to parse JSON: %w", err)
		}
	}

	// Override with flags if provided
	if assignmentName != "" {
		params.Name = assignmentName
	}
	if cmd.Flags().Changed("points") {
		params.PointsPossible = &assignmentPoints
	}
	if assignmentGradingType != "" {
		params.GradingType = assignmentGradingType
	}
	if cmd.Flags().Changed("due-at") {
		params.DueAt = &assignmentDueAt
	}
	if cmd.Flags().Changed("unlock-at") {
		params.UnlockAt = &assignmentUnlockAt
	}
	if cmd.Flags().Changed("lock-at") {
		params.LockAt = &assignmentLockAt
	}
	if assignmentDescription != "" {
		params.Description = assignmentDescription
	}
	if cmd.Flags().Changed("published") {
		params.Published = &assignmentPublished
	}
	if len(assignmentSubmissionTypes) > 0 {
		params.SubmissionTypes = assignmentSubmissionTypes
	}
	if cmd.Flags().Changed("group-id") {
		params.AssignmentGroupID = &assignmentGroupID
	}
	if cmd.Flags().Changed("position") {
		params.Position = &assignmentPosition
	}

	// Update assignment
	ctx := context.Background()
	assignment, err := assignmentsService.Update(ctx, assignmentsCourseID, assignmentID, params)
	if err != nil {
		return fmt.Errorf("failed to update assignment: %w", err)
	}

	fmt.Printf("Assignment updated successfully!\n")
	fmt.Printf("  ID: %d\n", assignment.ID)
	fmt.Printf("  Name: %s\n", assignment.Name)
	fmt.Printf("  Points: %.1f\n", assignment.PointsPossible)
	if !assignment.DueAt.IsZero() {
		fmt.Printf("  Due: %s\n", assignment.DueAt.Format("2006-01-02 15:04"))
	}
	if assignment.Published {
		fmt.Printf("  Status: Published\n")
	} else {
		fmt.Printf("  Status: Unpublished\n")
	}

	return nil
}

func runAssignmentsDelete(cmd *cobra.Command, args []string) error {
	// Parse assignment ID
	assignmentID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid assignment ID: %w", err)
	}

	// Get API client first to validate course ID before asking for confirmation
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Validate course ID exists
	if _, err := validateCourseID(client, assignmentsCourseID); err != nil {
		return err
	}

	// Confirm deletion
	fmt.Printf("Are you sure you want to delete assignment %d? [y/N]: ", assignmentID)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	if response != "y" && response != "yes" {
		fmt.Println("Cancelled.")
		return nil
	}

	// Create assignments service
	assignmentsService := api.NewAssignmentsService(client)

	// Delete assignment
	ctx := context.Background()
	if err := assignmentsService.Delete(ctx, assignmentsCourseID, assignmentID); err != nil {
		return fmt.Errorf("failed to delete assignment: %w", err)
	}

	fmt.Printf("Assignment %d deleted successfully.\n", assignmentID)
	return nil
}

// Helper functions for JSON input

func readAssignmentJSON(filePath string, useStdin bool) ([]byte, error) {
	if filePath != "" {
		return os.ReadFile(filePath)
	}
	if useStdin {
		return io.ReadAll(os.Stdin)
	}
	return nil, nil
}

type assignmentJSONInput struct {
	Name              string   `json:"name"`
	PointsPossible    float64  `json:"points_possible"`
	GradingType       string   `json:"grading_type"`
	DueAt             string   `json:"due_at"`
	UnlockAt          string   `json:"unlock_at"`
	LockAt            string   `json:"lock_at"`
	Description       string   `json:"description"`
	Published         *bool    `json:"published"`
	SubmissionTypes   []string `json:"submission_types"`
	AssignmentGroupID int64    `json:"assignment_group_id"`
	Position          int      `json:"position"`
}

func parseAssignmentCreateJSON(data []byte, params *api.CreateAssignmentParams) error {
	var input assignmentJSONInput
	if err := json.Unmarshal(data, &input); err != nil {
		return err
	}

	if input.Name != "" {
		params.Name = input.Name
	}
	if input.PointsPossible > 0 {
		params.PointsPossible = input.PointsPossible
	}
	if input.GradingType != "" {
		params.GradingType = input.GradingType
	}
	if input.DueAt != "" {
		params.DueAt = input.DueAt
	}
	if input.UnlockAt != "" {
		params.UnlockAt = input.UnlockAt
	}
	if input.LockAt != "" {
		params.LockAt = input.LockAt
	}
	if input.Description != "" {
		params.Description = input.Description
	}
	if input.Published != nil {
		params.Published = *input.Published
	}
	if len(input.SubmissionTypes) > 0 {
		params.SubmissionTypes = input.SubmissionTypes
	}
	if input.AssignmentGroupID > 0 {
		params.AssignmentGroupID = input.AssignmentGroupID
	}
	if input.Position > 0 {
		params.Position = input.Position
	}

	return nil
}

func parseAssignmentUpdateJSON(data []byte, params *api.UpdateAssignmentParams) error {
	var input assignmentJSONInput
	if err := json.Unmarshal(data, &input); err != nil {
		return err
	}

	if input.Name != "" {
		params.Name = input.Name
	}
	if input.PointsPossible > 0 {
		params.PointsPossible = &input.PointsPossible
	}
	if input.GradingType != "" {
		params.GradingType = input.GradingType
	}
	if input.DueAt != "" {
		params.DueAt = &input.DueAt
	}
	if input.UnlockAt != "" {
		params.UnlockAt = &input.UnlockAt
	}
	if input.LockAt != "" {
		params.LockAt = &input.LockAt
	}
	if input.Description != "" {
		params.Description = input.Description
	}
	if input.Published != nil {
		params.Published = input.Published
	}
	if len(input.SubmissionTypes) > 0 {
		params.SubmissionTypes = input.SubmissionTypes
	}
	if input.AssignmentGroupID > 0 {
		params.AssignmentGroupID = &input.AssignmentGroupID
	}
	if input.Position > 0 {
		params.Position = &input.Position
	}

	return nil
}
