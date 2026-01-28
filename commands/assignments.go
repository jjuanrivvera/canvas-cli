package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
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

func init() {
	rootCmd.AddCommand(assignmentsCmd)
	assignmentsCmd.AddCommand(newAssignmentsListCmd())
	assignmentsCmd.AddCommand(newAssignmentsGetCmd())
	assignmentsCmd.AddCommand(newAssignmentsCreateCmd())
	assignmentsCmd.AddCommand(newAssignmentsUpdateCmd())
	assignmentsCmd.AddCommand(newAssignmentsDeleteCmd())
}

func newAssignmentsListCmd() *cobra.Command {
	opts := &options.AssignmentsListOptions{}

	cmd := &cobra.Command{
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
  canvas assignments list --course-id 123 --include submission,rubric

Note: If you have set a course context (canvas context set course 123),
you can omit --course-id and it will be used automatically.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Apply context values if flags not provided
			opts.CourseID = GetContextCourseID(opts.CourseID)

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAssignmentsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required, or set via 'canvas context set course')")
	cmd.Flags().StringVar(&opts.SearchTerm, "search", "", "Search by assignment name")
	cmd.Flags().StringVar(&opts.Bucket, "bucket", "", "Filter by bucket (past, overdue, undated, ungraded, unsubmitted, upcoming, future)")
	cmd.Flags().StringVar(&opts.OrderBy, "order-by", "", "Order by (position, name, due_at)")
	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Additional data to include (comma-separated)")

	return cmd
}

func newAssignmentsGetCmd() *cobra.Command {
	opts := &options.AssignmentsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <assignment-id>",
		Short: "Get details of a specific assignment",
		Long: `Get details of a specific assignment by ID.

Examples:
  canvas assignments get --course-id 123 456
  canvas assignments get --course-id 123 456 --include submission,rubric

Note: If you have set a course context, you can omit --course-id.`,
		Args: ExactArgsWithUsage(1, "assignment-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			assignmentID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid assignment ID: %s", args[0])
			}
			opts.AssignmentID = assignmentID

			// Apply context values if flags not provided
			opts.CourseID = GetContextCourseID(opts.CourseID)

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAssignmentsGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required, or set via context)")
	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Additional data to include (comma-separated)")

	return cmd
}

func newAssignmentsCreateCmd() *cobra.Command {
	opts := &options.AssignmentsCreateOptions{}

	cmd := &cobra.Command{
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
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAssignmentsCreate(cmd.Context(), client, cmd, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Assignment name")
	cmd.Flags().Float64Var(&opts.Points, "points", 0, "Points possible")
	cmd.Flags().StringVar(&opts.GradingType, "grading-type", "", "Grading type (points, pass_fail, percent, letter_grade, gpa_scale, not_graded)")
	cmd.Flags().StringVar(&opts.DueAt, "due-at", "", "Due date (ISO8601 format)")
	cmd.Flags().StringVar(&opts.UnlockAt, "unlock-at", "", "Unlock date (ISO8601 format)")
	cmd.Flags().StringVar(&opts.LockAt, "lock-at", "", "Lock date (ISO8601 format)")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Assignment description (HTML)")
	cmd.Flags().BoolVar(&opts.Published, "published", false, "Publish the assignment")
	cmd.Flags().StringSliceVar(&opts.SubmissionTypes, "submission-types", []string{}, "Submission types (online_text_entry, online_url, online_upload, media_recording, none)")
	cmd.Flags().Int64Var(&opts.GroupID, "group-id", 0, "Assignment group ID")
	cmd.Flags().IntVar(&opts.Position, "position", 0, "Position in the assignment group")
	cmd.Flags().StringVar(&opts.JSONFile, "json", "", "JSON file with assignment data")
	cmd.Flags().BoolVar(&opts.Stdin, "stdin", false, "Read JSON from stdin")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newAssignmentsUpdateCmd() *cobra.Command {
	opts := &options.AssignmentsUpdateOptions{}

	cmd := &cobra.Command{
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
		Args: ExactArgsWithUsage(1, "assignment-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			assignmentID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid assignment ID: %s", args[0])
			}
			opts.AssignmentID = assignmentID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAssignmentsUpdate(cmd.Context(), client, cmd, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Assignment name")
	cmd.Flags().Float64Var(&opts.Points, "points", 0, "Points possible")
	cmd.Flags().StringVar(&opts.GradingType, "grading-type", "", "Grading type (points, pass_fail, percent, letter_grade, gpa_scale, not_graded)")
	cmd.Flags().StringVar(&opts.DueAt, "due-at", "", "Due date (ISO8601 format)")
	cmd.Flags().StringVar(&opts.UnlockAt, "unlock-at", "", "Unlock date (ISO8601 format)")
	cmd.Flags().StringVar(&opts.LockAt, "lock-at", "", "Lock date (ISO8601 format)")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Assignment description (HTML)")
	cmd.Flags().BoolVar(&opts.Published, "published", false, "Publish the assignment")
	cmd.Flags().StringSliceVar(&opts.SubmissionTypes, "submission-types", []string{}, "Submission types")
	cmd.Flags().Int64Var(&opts.GroupID, "group-id", 0, "Assignment group ID")
	cmd.Flags().IntVar(&opts.Position, "position", 0, "Position in the assignment group")
	cmd.Flags().StringVar(&opts.JSONFile, "json", "", "JSON file with assignment data")
	cmd.Flags().BoolVar(&opts.Stdin, "stdin", false, "Read JSON from stdin")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newAssignmentsDeleteCmd() *cobra.Command {
	opts := &options.AssignmentsDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <assignment-id>",
		Short: "Delete an assignment",
		Long: `Delete an assignment from a Canvas course.

This action cannot be undone.

Examples:
  canvas assignments delete --course-id 123 456`,
		Args: ExactArgsWithUsage(1, "assignment-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			assignmentID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid assignment ID: %s", args[0])
			}
			opts.AssignmentID = assignmentID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAssignmentsDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Skip confirmation prompt")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

// Run functions

func runAssignmentsList(ctx context.Context, client *api.Client, opts *options.AssignmentsListOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "assignments.list", map[string]interface{}{
		"course_id": opts.CourseID,
		"bucket":    opts.Bucket,
		"search":    opts.SearchTerm,
	})

	// Validate course ID exists
	if _, err := validateCourseID(client, opts.CourseID); err != nil {
		logger.LogCommandError(ctx, "assignments.list", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return err
	}

	assignmentsService := api.NewAssignmentsService(client)

	apiOpts := &api.ListAssignmentsOptions{
		SearchTerm: opts.SearchTerm,
		Bucket:     opts.Bucket,
		OrderBy:    opts.OrderBy,
		Include:    opts.Include,
	}

	assignments, err := assignmentsService.List(ctx, opts.CourseID, apiOpts)
	if err != nil {
		logger.LogCommandError(ctx, "assignments.list", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to list assignments: %w", err)
	}

	if err := formatEmptyOrOutput(assignments, "No assignments found"); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "assignments.list", len(assignments))
	return nil
}

func runAssignmentsGet(ctx context.Context, client *api.Client, opts *options.AssignmentsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "assignments.get", map[string]interface{}{
		"course_id":     opts.CourseID,
		"assignment_id": opts.AssignmentID,
	})

	// Validate course ID exists
	if _, err := validateCourseID(client, opts.CourseID); err != nil {
		logger.LogCommandError(ctx, "assignments.get", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return err
	}

	assignmentsService := api.NewAssignmentsService(client)

	assignment, err := assignmentsService.Get(ctx, opts.CourseID, opts.AssignmentID, opts.Include)
	if err != nil {
		logger.LogCommandError(ctx, "assignments.get", err, map[string]interface{}{
			"course_id":     opts.CourseID,
			"assignment_id": opts.AssignmentID,
		})
		return fmt.Errorf("failed to get assignment: %w", err)
	}

	if err := formatOutput(assignment, nil); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "assignments.get", 1)
	return nil
}

func runAssignmentsCreate(ctx context.Context, client *api.Client, cmd *cobra.Command, opts *options.AssignmentsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "assignments.create", map[string]interface{}{
		"course_id": opts.CourseID,
		"name":      opts.Name,
	})

	// Validate course ID exists
	if _, err := validateCourseID(client, opts.CourseID); err != nil {
		logger.LogCommandError(ctx, "assignments.create", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return err
	}

	assignmentsService := api.NewAssignmentsService(client)

	// Build params from flags or JSON
	params := &api.CreateAssignmentParams{}

	// Check for JSON input
	if opts.JSONFile != "" || opts.Stdin {
		jsonData, err := readAssignmentJSON(opts.JSONFile, opts.Stdin)
		if err != nil {
			logger.LogCommandError(ctx, "assignments.create", err, map[string]interface{}{
				"course_id": opts.CourseID,
			})
			return fmt.Errorf("failed to read JSON: %w", err)
		}
		if err := parseAssignmentCreateJSON(jsonData, params); err != nil {
			logger.LogCommandError(ctx, "assignments.create", err, map[string]interface{}{
				"course_id": opts.CourseID,
			})
			return fmt.Errorf("failed to parse JSON: %w", err)
		}
	}

	// Override with flags if provided
	if opts.Name != "" {
		params.Name = opts.Name
	}
	if opts.Points > 0 {
		params.PointsPossible = opts.Points
	}
	if opts.GradingType != "" {
		params.GradingType = opts.GradingType
	}
	if opts.DueAt != "" {
		params.DueAt = opts.DueAt
	}
	if opts.UnlockAt != "" {
		params.UnlockAt = opts.UnlockAt
	}
	if opts.LockAt != "" {
		params.LockAt = opts.LockAt
	}
	if opts.Description != "" {
		params.Description = opts.Description
	}
	if cmd.Flags().Changed("published") {
		params.Published = opts.Published
	}
	if len(opts.SubmissionTypes) > 0 {
		params.SubmissionTypes = opts.SubmissionTypes
	}
	if opts.GroupID > 0 {
		params.AssignmentGroupID = opts.GroupID
	}
	if opts.Position > 0 {
		params.Position = opts.Position
	}

	// Validate required fields
	if params.Name == "" {
		err := fmt.Errorf("assignment name is required (use --name or provide in JSON)")
		logger.LogCommandError(ctx, "assignments.create", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return err
	}

	assignment, err := assignmentsService.Create(ctx, opts.CourseID, params)
	if err != nil {
		logger.LogCommandError(ctx, "assignments.create", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"name":      params.Name,
		})
		return fmt.Errorf("failed to create assignment: %w", err)
	}

	printInfo("Assignment created successfully!\n")
	printInfo("  ID: %d\n", assignment.ID)
	printInfo("  Name: %s\n", assignment.Name)
	printInfo("  Points: %.1f\n", assignment.PointsPossible)
	if !assignment.DueAt.IsZero() {
		printInfo("  Due: %s\n", assignment.DueAt.Format("2006-01-02 15:04"))
	}
	if assignment.Published {
		printInfo("  Status: Published\n")
	} else {
		printInfo("  Status: Unpublished\n")
	}

	logger.LogCommandComplete(ctx, "assignments.create", 1)
	return nil
}

func runAssignmentsUpdate(ctx context.Context, client *api.Client, cmd *cobra.Command, opts *options.AssignmentsUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "assignments.update", map[string]interface{}{
		"course_id":     opts.CourseID,
		"assignment_id": opts.AssignmentID,
	})

	// Validate course ID exists
	if _, err := validateCourseID(client, opts.CourseID); err != nil {
		logger.LogCommandError(ctx, "assignments.update", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return err
	}

	assignmentsService := api.NewAssignmentsService(client)

	// Build params from flags or JSON
	params := &api.UpdateAssignmentParams{}

	// Check for JSON input
	if opts.JSONFile != "" || opts.Stdin {
		jsonData, err := readAssignmentJSON(opts.JSONFile, opts.Stdin)
		if err != nil {
			logger.LogCommandError(ctx, "assignments.update", err, map[string]interface{}{
				"course_id":     opts.CourseID,
				"assignment_id": opts.AssignmentID,
			})
			return fmt.Errorf("failed to read JSON: %w", err)
		}
		if err := parseAssignmentUpdateJSON(jsonData, params); err != nil {
			logger.LogCommandError(ctx, "assignments.update", err, map[string]interface{}{
				"course_id":     opts.CourseID,
				"assignment_id": opts.AssignmentID,
			})
			return fmt.Errorf("failed to parse JSON: %w", err)
		}
	}

	// Override with flags if provided
	if opts.Name != "" {
		params.Name = opts.Name
	}
	if cmd.Flags().Changed("points") {
		params.PointsPossible = &opts.Points
	}
	if opts.GradingType != "" {
		params.GradingType = opts.GradingType
	}
	if cmd.Flags().Changed("due-at") {
		params.DueAt = &opts.DueAt
	}
	if cmd.Flags().Changed("unlock-at") {
		params.UnlockAt = &opts.UnlockAt
	}
	if cmd.Flags().Changed("lock-at") {
		params.LockAt = &opts.LockAt
	}
	if opts.Description != "" {
		params.Description = opts.Description
	}
	if cmd.Flags().Changed("published") {
		params.Published = &opts.Published
	}
	if len(opts.SubmissionTypes) > 0 {
		params.SubmissionTypes = opts.SubmissionTypes
	}
	if cmd.Flags().Changed("group-id") {
		params.AssignmentGroupID = &opts.GroupID
	}
	if cmd.Flags().Changed("position") {
		params.Position = &opts.Position
	}

	// Handle dry-run: show what would be updated
	if dryRun {
		changes := make(map[string]interface{})
		if params.Name != "" {
			changes["name"] = params.Name
		}
		if params.PointsPossible != nil {
			changes["points_possible"] = *params.PointsPossible
		}
		if params.GradingType != "" {
			changes["grading_type"] = params.GradingType
		}
		if params.DueAt != nil {
			changes["due_at"] = *params.DueAt
		}
		if params.Description != "" {
			changes["description"] = "(HTML content)"
		}
		if params.Published != nil {
			changes["published"] = *params.Published
		}
		if len(params.SubmissionTypes) > 0 {
			changes["submission_types"] = params.SubmissionTypes
		}

		confirmUpdateDryRun("assignment", opts.AssignmentID, changes)
		return nil
	}

	assignment, err := assignmentsService.Update(ctx, opts.CourseID, opts.AssignmentID, params)
	if err != nil {
		logger.LogCommandError(ctx, "assignments.update", err, map[string]interface{}{
			"course_id":     opts.CourseID,
			"assignment_id": opts.AssignmentID,
		})
		return fmt.Errorf("failed to update assignment: %w", err)
	}

	printInfo("Assignment updated successfully!\n")
	printInfo("  ID: %d\n", assignment.ID)
	printInfo("  Name: %s\n", assignment.Name)
	printInfo("  Points: %.1f\n", assignment.PointsPossible)
	if !assignment.DueAt.IsZero() {
		printInfo("  Due: %s\n", assignment.DueAt.Format("2006-01-02 15:04"))
	}
	if assignment.Published {
		printInfo("  Status: Published\n")
	} else {
		printInfo("  Status: Unpublished\n")
	}

	logger.LogCommandComplete(ctx, "assignments.update", 1)
	return nil
}

func runAssignmentsDelete(ctx context.Context, client *api.Client, opts *options.AssignmentsDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "assignments.delete", map[string]interface{}{
		"course_id":     opts.CourseID,
		"assignment_id": opts.AssignmentID,
	})

	// Validate course ID exists
	if _, err := validateCourseID(client, opts.CourseID); err != nil {
		logger.LogCommandError(ctx, "assignments.delete", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return err
	}

	assignmentsService := api.NewAssignmentsService(client)

	// Fetch assignment details for preview (especially in dry-run mode)
	assignment, err := assignmentsService.Get(ctx, opts.CourseID, opts.AssignmentID, nil)
	if err != nil {
		logger.LogCommandError(ctx, "assignments.delete", err, map[string]interface{}{
			"course_id":     opts.CourseID,
			"assignment_id": opts.AssignmentID,
		})
		return fmt.Errorf("failed to get assignment: %w", err)
	}

	// Build details for confirmation
	details := map[string]interface{}{
		"Name":   assignment.Name,
		"Points": assignment.PointsPossible,
	}
	if !assignment.DueAt.IsZero() {
		details["Due"] = assignment.DueAt.Format("2006-01-02 15:04")
	}
	if assignment.Published {
		details["Status"] = "Published"
	} else {
		details["Status"] = "Unpublished"
	}

	// Confirm deletion with details (handles dry-run internally)
	confirmed, err := confirmDeleteWithDetails("assignment", opts.AssignmentID, details, opts.Force)
	if err != nil {
		return err
	}
	if !confirmed {
		if !dryRun {
			fmt.Println("Cancelled.")
		}
		return nil
	}

	if err := assignmentsService.Delete(ctx, opts.CourseID, opts.AssignmentID); err != nil {
		logger.LogCommandError(ctx, "assignments.delete", err, map[string]interface{}{
			"course_id":     opts.CourseID,
			"assignment_id": opts.AssignmentID,
		})
		return fmt.Errorf("failed to delete assignment: %w", err)
	}

	printInfo("Assignment %d deleted successfully.\n", opts.AssignmentID)

	logger.LogCommandComplete(ctx, "assignments.delete", 1)
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
