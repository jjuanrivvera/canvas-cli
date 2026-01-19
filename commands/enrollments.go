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

// enrollmentsCmd represents the enrollments command group
var enrollmentsCmd = &cobra.Command{
	Use:   "enrollments",
	Short: "Manage Canvas enrollments",
	Long: `Manage Canvas enrollments including listing, creating, and managing course enrollments.

Examples:
  canvas enrollments list --course-id 123
  canvas enrollments list --user-id 456
  canvas enrollments get 789`,
}

func init() {
	rootCmd.AddCommand(enrollmentsCmd)
	enrollmentsCmd.AddCommand(newEnrollmentsListCmd())
	enrollmentsCmd.AddCommand(newEnrollmentsGetCmd())
	enrollmentsCmd.AddCommand(newEnrollmentsCreateCmd())
	enrollmentsCmd.AddCommand(newEnrollmentsConcludeCmd())
	enrollmentsCmd.AddCommand(newEnrollmentsReactivateCmd())
	enrollmentsCmd.AddCommand(newEnrollmentsAcceptCmd())
	enrollmentsCmd.AddCommand(newEnrollmentsRejectCmd())
}

func newEnrollmentsListCmd() *cobra.Command {
	opts := &options.EnrollmentsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List enrollments",
		Long: `List enrollments in a course or for a user.

You must specify one of --course-id or --user-id to indicate the context.

Course context:
  canvas enrollments list --course-id 123         # All enrollments in course
  canvas enrollments list --course-id 123 --type StudentEnrollment
  canvas enrollments list --course-id 123 --state active

User context:
  canvas enrollments list --user-id 456           # All enrollments for user
  canvas enrollments list --user-id 456 --state active

Examples:
  canvas enrollments list --course-id 123
  canvas enrollments list --user-id 456
  canvas enrollments list --course-id 123 --type TeacherEnrollment
  canvas enrollments list --course-id 123 --state active,invited
  canvas enrollments list --user-id 456 --include current_points`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runEnrollmentsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (for course enrollments)")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (for user enrollments)")
	cmd.Flags().StringSliceVar(&opts.Type, "type", []string{}, "Filter by enrollment type (StudentEnrollment, TeacherEnrollment, TaEnrollment, ObserverEnrollment, DesignerEnrollment)")
	cmd.Flags().StringSliceVar(&opts.State, "state", []string{}, "Filter by enrollment state (active, invited, creation_pending, deleted, rejected, completed, inactive)")
	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Additional data to include (avatar_url, group_ids, locked, observed_users, can_be_removed, uuid, current_points)")

	return cmd
}

func newEnrollmentsGetCmd() *cobra.Command {
	opts := &options.EnrollmentsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <enrollment-id>",
		Short: "Get enrollment details",
		Long: `Get details of a specific enrollment.

Note: You must specify --course-id to indicate which course the enrollment belongs to.

Examples:
  canvas enrollments get 789 --course-id 123`,
		Args: ExactArgsWithUsage(1, "enrollment-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			enrollmentID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid enrollment ID: %w", err)
			}
			opts.EnrollmentID = enrollmentID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runEnrollmentsGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newEnrollmentsCreateCmd() *cobra.Command {
	opts := &options.EnrollmentsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Enroll a user in a course",
		Long: `Enroll a user in a course.

Enrollment Types:
  StudentEnrollment   - Student role
  TeacherEnrollment   - Teacher/Instructor role
  TaEnrollment        - Teaching Assistant role
  ObserverEnrollment  - Observer role
  DesignerEnrollment  - Course Designer role

Enrollment States:
  active              - Active enrollment
  invited             - Invited (default)
  inactive            - Inactive enrollment

Examples:
  canvas enrollments create --course-id 123 --user-id 456
  canvas enrollments create --course-id 123 --user-id 456 --type TeacherEnrollment
  canvas enrollments create --course-id 123 --user-id 456 --state active --notify
  canvas enrollments create --course-id 123 --user-id 456 --section-id 789`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runEnrollmentsCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID to enroll (required)")
	cmd.Flags().StringVar(&opts.EnrollmentType, "type", "StudentEnrollment", "Enrollment type")
	cmd.Flags().StringVar(&opts.EnrollmentState, "state", "", "Initial enrollment state (active, invited)")
	cmd.Flags().Int64Var(&opts.SectionID, "section-id", 0, "Section ID")
	cmd.Flags().BoolVar(&opts.Notify, "notify", false, "Send enrollment notification email")
	cmd.Flags().StringVar(&opts.Role, "role", "", "Custom role name")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("user-id")

	return cmd
}

func newEnrollmentsConcludeCmd() *cobra.Command {
	opts := &options.EnrollmentsConcludeOptions{}

	cmd := &cobra.Command{
		Use:   "conclude <enrollment-id>",
		Short: "Conclude an enrollment",
		Long: `Conclude an enrollment (soft delete).

Tasks:
  conclude  - Conclude the enrollment (default)
  deactivate - Mark enrollment as inactive
  delete    - Permanently delete the enrollment

Examples:
  canvas enrollments conclude 789 --course-id 123
  canvas enrollments conclude 789 --course-id 123 --task deactivate
  canvas enrollments conclude 789 --course-id 123 --task delete --force`,
		Args: ExactArgsWithUsage(1, "enrollment-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			enrollmentID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid enrollment ID: %w", err)
			}
			opts.EnrollmentID = enrollmentID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runEnrollmentsConclude(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.Task, "task", "conclude", "Task: conclude, deactivate, delete")
	cmd.Flags().BoolVar(&opts.Force, "force", false, "Skip confirmation for delete")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newEnrollmentsReactivateCmd() *cobra.Command {
	opts := &options.EnrollmentsReactivateOptions{}

	cmd := &cobra.Command{
		Use:   "reactivate <enrollment-id>",
		Short: "Reactivate a concluded enrollment",
		Long: `Reactivate a previously concluded or deactivated enrollment.

Examples:
  canvas enrollments reactivate 789 --course-id 123`,
		Args: ExactArgsWithUsage(1, "enrollment-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			enrollmentID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid enrollment ID: %w", err)
			}
			opts.EnrollmentID = enrollmentID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runEnrollmentsReactivate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newEnrollmentsAcceptCmd() *cobra.Command {
	opts := &options.EnrollmentsAcceptOptions{}

	cmd := &cobra.Command{
		Use:   "accept <enrollment-id>",
		Short: "Accept a pending enrollment invitation",
		Long: `Accept a pending enrollment invitation.

Examples:
  canvas enrollments accept 789 --course-id 123`,
		Args: ExactArgsWithUsage(1, "enrollment-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			enrollmentID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid enrollment ID: %w", err)
			}
			opts.EnrollmentID = enrollmentID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runEnrollmentsAccept(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newEnrollmentsRejectCmd() *cobra.Command {
	opts := &options.EnrollmentsRejectOptions{}

	cmd := &cobra.Command{
		Use:   "reject <enrollment-id>",
		Short: "Reject a pending enrollment invitation",
		Long: `Reject a pending enrollment invitation.

Examples:
  canvas enrollments reject 789 --course-id 123`,
		Args: ExactArgsWithUsage(1, "enrollment-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			enrollmentID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid enrollment ID: %w", err)
			}
			opts.EnrollmentID = enrollmentID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runEnrollmentsReject(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func runEnrollmentsList(ctx context.Context, client *api.Client, opts *options.EnrollmentsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "enrollments.list", map[string]interface{}{
		"course_id": opts.CourseID,
		"user_id":   opts.UserID,
	})

	enrollmentsService := api.NewEnrollmentsService(client)

	// Build options
	apiOpts := &api.ListEnrollmentsOptions{
		Type:    opts.Type,
		State:   opts.State,
		Include: opts.Include,
	}

	// List enrollments based on context
	var enrollments []api.Enrollment
	var contextName string
	var err error

	if opts.CourseID > 0 {
		// Course context - list all enrollments in the course
		enrollments, err = enrollmentsService.ListCourse(ctx, opts.CourseID, apiOpts)
		contextName = fmt.Sprintf("course %d", opts.CourseID)
	} else {
		// User context - list all enrollments for the user
		enrollments, err = enrollmentsService.ListUser(ctx, opts.UserID, apiOpts)
		contextName = fmt.Sprintf("user %d", opts.UserID)
	}

	if err != nil {
		logger.LogCommandError(ctx, "enrollments.list", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"user_id":   opts.UserID,
		})
		return fmt.Errorf("failed to list enrollments: %w", err)
	}

	if len(enrollments) == 0 {
		logger.LogCommandComplete(ctx, "enrollments.list", 0)
		fmt.Printf("No enrollments found for %s\n", contextName)
		return nil
	}

	// Format and display enrollments
	printVerbose("Found %d enrollments for %s:\n\n", len(enrollments), contextName)
	logger.LogCommandComplete(ctx, "enrollments.list", len(enrollments))
	return formatOutput(enrollments, nil)
}

func runEnrollmentsGet(ctx context.Context, client *api.Client, opts *options.EnrollmentsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "enrollments.get", map[string]interface{}{
		"course_id":     opts.CourseID,
		"enrollment_id": opts.EnrollmentID,
	})

	enrollmentsService := api.NewEnrollmentsService(client)

	// Get enrollment by listing course enrollments and filtering
	// Note: Canvas doesn't have a direct "get enrollment by ID" endpoint
	// so we list all enrollments for the course and find the matching one
	enrollments, err := enrollmentsService.ListCourse(ctx, opts.CourseID, nil)
	if err != nil {
		logger.LogCommandError(ctx, "enrollments.get", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to list enrollments: %w", err)
	}

	// Find the enrollment by ID
	var enrollment *api.Enrollment
	for i := range enrollments {
		if enrollments[i].ID == opts.EnrollmentID {
			enrollment = &enrollments[i]
			break
		}
	}

	if enrollment == nil {
		err := fmt.Errorf("enrollment %d not found in course %d", opts.EnrollmentID, opts.CourseID)
		logger.LogCommandError(ctx, "enrollments.get", err, map[string]interface{}{
			"course_id":     opts.CourseID,
			"enrollment_id": opts.EnrollmentID,
		})
		return err
	}

	// Format and display enrollment details
	logger.LogCommandComplete(ctx, "enrollments.get", 1)
	return formatOutput(enrollment, nil)
}

func runEnrollmentsCreate(ctx context.Context, client *api.Client, opts *options.EnrollmentsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "enrollments.create", map[string]interface{}{
		"course_id": opts.CourseID,
		"user_id":   opts.UserID,
		"type":      opts.EnrollmentType,
	})

	enrollmentsService := api.NewEnrollmentsService(client)

	// Build params
	params := &api.EnrollUserParams{
		UserID:          opts.UserID,
		Type:            opts.EnrollmentType,
		EnrollmentState: opts.EnrollmentState,
		CourseSectionID: opts.SectionID,
		Notify:          opts.Notify,
	}

	// Note: role flag maps to a custom role name, not RoleID
	// For custom roles, you'd need to look up the role ID first
	// For now, we support the standard enrollment types

	// Create enrollment
	enrollment, err := enrollmentsService.EnrollUser(ctx, opts.CourseID, params)
	if err != nil {
		logger.LogCommandError(ctx, "enrollments.create", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"user_id":   opts.UserID,
		})
		return fmt.Errorf("failed to create enrollment: %w", err)
	}

	fmt.Printf("Enrollment created successfully (ID: %d)\n", enrollment.ID)
	logger.LogCommandComplete(ctx, "enrollments.create", 1)
	return formatOutput(enrollment, nil)
}

func runEnrollmentsConclude(ctx context.Context, client *api.Client, opts *options.EnrollmentsConcludeOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "enrollments.conclude", map[string]interface{}{
		"course_id":     opts.CourseID,
		"enrollment_id": opts.EnrollmentID,
		"task":          opts.Task,
	})

	// Confirmation for delete
	if opts.Task == "delete" && !opts.Force {
		fmt.Printf("WARNING: This will permanently delete enrollment %d.\n", opts.EnrollmentID)
		fmt.Print("Type 'yes' to confirm: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			logger.LogCommandComplete(ctx, "enrollments.conclude", 0)
			fmt.Println("Delete cancelled")
			return nil
		}
	}

	enrollmentsService := api.NewEnrollmentsService(client)

	// Conclude/deactivate/delete enrollment
	enrollment, err := enrollmentsService.Conclude(ctx, opts.CourseID, opts.EnrollmentID, opts.Task)
	if err != nil {
		logger.LogCommandError(ctx, "enrollments.conclude", err, map[string]interface{}{
			"course_id":     opts.CourseID,
			"enrollment_id": opts.EnrollmentID,
			"task":          opts.Task,
		})
		return fmt.Errorf("failed to %s enrollment: %w", opts.Task, err)
	}

	switch opts.Task {
	case "conclude":
		fmt.Printf("Enrollment %d concluded\n", opts.EnrollmentID)
	case "deactivate":
		fmt.Printf("Enrollment %d deactivated\n", opts.EnrollmentID)
	case "delete":
		fmt.Printf("Enrollment %d deleted\n", opts.EnrollmentID)
	}

	logger.LogCommandComplete(ctx, "enrollments.conclude", 1)
	return formatOutput(enrollment, nil)
}

func runEnrollmentsReactivate(ctx context.Context, client *api.Client, opts *options.EnrollmentsReactivateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "enrollments.reactivate", map[string]interface{}{
		"course_id":     opts.CourseID,
		"enrollment_id": opts.EnrollmentID,
	})

	enrollmentsService := api.NewEnrollmentsService(client)

	// Reactivate enrollment
	enrollment, err := enrollmentsService.Reactivate(ctx, opts.CourseID, opts.EnrollmentID)
	if err != nil {
		logger.LogCommandError(ctx, "enrollments.reactivate", err, map[string]interface{}{
			"course_id":     opts.CourseID,
			"enrollment_id": opts.EnrollmentID,
		})
		return fmt.Errorf("failed to reactivate enrollment: %w", err)
	}

	fmt.Printf("Enrollment %d reactivated\n", opts.EnrollmentID)
	logger.LogCommandComplete(ctx, "enrollments.reactivate", 1)
	return formatOutput(enrollment, nil)
}

func runEnrollmentsAccept(ctx context.Context, client *api.Client, opts *options.EnrollmentsAcceptOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "enrollments.accept", map[string]interface{}{
		"course_id":     opts.CourseID,
		"enrollment_id": opts.EnrollmentID,
	})

	enrollmentsService := api.NewEnrollmentsService(client)

	// Accept enrollment invitation
	err := enrollmentsService.Accept(ctx, opts.CourseID, opts.EnrollmentID)
	if err != nil {
		logger.LogCommandError(ctx, "enrollments.accept", err, map[string]interface{}{
			"course_id":     opts.CourseID,
			"enrollment_id": opts.EnrollmentID,
		})
		return fmt.Errorf("failed to accept enrollment: %w", err)
	}

	fmt.Printf("Enrollment invitation %d accepted\n", opts.EnrollmentID)
	logger.LogCommandComplete(ctx, "enrollments.accept", 1)
	return nil
}

func runEnrollmentsReject(ctx context.Context, client *api.Client, opts *options.EnrollmentsRejectOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "enrollments.reject", map[string]interface{}{
		"course_id":     opts.CourseID,
		"enrollment_id": opts.EnrollmentID,
	})

	enrollmentsService := api.NewEnrollmentsService(client)

	// Reject enrollment invitation
	err := enrollmentsService.Reject(ctx, opts.CourseID, opts.EnrollmentID)
	if err != nil {
		logger.LogCommandError(ctx, "enrollments.reject", err, map[string]interface{}{
			"course_id":     opts.CourseID,
			"enrollment_id": opts.EnrollmentID,
		})
		return fmt.Errorf("failed to reject enrollment: %w", err)
	}

	fmt.Printf("Enrollment invitation %d rejected\n", opts.EnrollmentID)
	logger.LogCommandComplete(ctx, "enrollments.reject", 1)
	return nil
}
