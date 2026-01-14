package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var (
	enrollmentsCourseID int64
	enrollmentsUserID   int64
	enrollmentsType     []string
	enrollmentsState    []string
	enrollmentsInclude  []string

	// Create flags
	enrollmentsEnrollUserID int64
	enrollmentsEnrollType   string
	enrollmentsEnrollState  string
	enrollmentsSectionID    int64
	enrollmentsNotify       bool
	enrollmentsRole         string

	// Conclude/Delete flags
	enrollmentsConcludeTask string
	enrollmentsForce        bool
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

// enrollmentsListCmd represents the enrollments list command
var enrollmentsListCmd = &cobra.Command{
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
	RunE: runEnrollmentsList,
}

// enrollmentsGetCmd represents the enrollments get command
var enrollmentsGetCmd = &cobra.Command{
	Use:   "get <enrollment-id>",
	Short: "Get enrollment details",
	Long: `Get details of a specific enrollment.

Note: You must specify --course-id to indicate which course the enrollment belongs to.

Examples:
  canvas enrollments get 789 --course-id 123`,
	Args: ExactArgsWithUsage(1, "enrollment-id"),
	RunE: runEnrollmentsGet,
}

// enrollmentsCreateCmd represents the enrollments create command
var enrollmentsCreateCmd = &cobra.Command{
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
	RunE: runEnrollmentsCreate,
}

// enrollmentsConcludeCmd represents the enrollments conclude command
var enrollmentsConcludeCmd = &cobra.Command{
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
	RunE: runEnrollmentsConclude,
}

// enrollmentsReactivateCmd represents the enrollments reactivate command
var enrollmentsReactivateCmd = &cobra.Command{
	Use:   "reactivate <enrollment-id>",
	Short: "Reactivate a concluded enrollment",
	Long: `Reactivate a previously concluded or deactivated enrollment.

Examples:
  canvas enrollments reactivate 789 --course-id 123`,
	Args: ExactArgsWithUsage(1, "enrollment-id"),
	RunE: runEnrollmentsReactivate,
}

// enrollmentsAcceptCmd represents the enrollments accept command
var enrollmentsAcceptCmd = &cobra.Command{
	Use:   "accept <enrollment-id>",
	Short: "Accept a pending enrollment invitation",
	Long: `Accept a pending enrollment invitation.

Examples:
  canvas enrollments accept 789 --course-id 123`,
	Args: ExactArgsWithUsage(1, "enrollment-id"),
	RunE: runEnrollmentsAccept,
}

// enrollmentsRejectCmd represents the enrollments reject command
var enrollmentsRejectCmd = &cobra.Command{
	Use:   "reject <enrollment-id>",
	Short: "Reject a pending enrollment invitation",
	Long: `Reject a pending enrollment invitation.

Examples:
  canvas enrollments reject 789 --course-id 123`,
	Args: ExactArgsWithUsage(1, "enrollment-id"),
	RunE: runEnrollmentsReject,
}

func init() {
	rootCmd.AddCommand(enrollmentsCmd)
	enrollmentsCmd.AddCommand(enrollmentsListCmd)
	enrollmentsCmd.AddCommand(enrollmentsGetCmd)
	enrollmentsCmd.AddCommand(enrollmentsCreateCmd)
	enrollmentsCmd.AddCommand(enrollmentsConcludeCmd)
	enrollmentsCmd.AddCommand(enrollmentsReactivateCmd)
	enrollmentsCmd.AddCommand(enrollmentsAcceptCmd)
	enrollmentsCmd.AddCommand(enrollmentsRejectCmd)

	// List flags
	enrollmentsListCmd.Flags().Int64Var(&enrollmentsCourseID, "course-id", 0, "Course ID (for course enrollments)")
	enrollmentsListCmd.Flags().Int64Var(&enrollmentsUserID, "user-id", 0, "User ID (for user enrollments)")
	enrollmentsListCmd.Flags().StringSliceVar(&enrollmentsType, "type", []string{}, "Filter by enrollment type (StudentEnrollment, TeacherEnrollment, TaEnrollment, ObserverEnrollment, DesignerEnrollment)")
	enrollmentsListCmd.Flags().StringSliceVar(&enrollmentsState, "state", []string{}, "Filter by enrollment state (active, invited, creation_pending, deleted, rejected, completed, inactive)")
	enrollmentsListCmd.Flags().StringSliceVar(&enrollmentsInclude, "include", []string{}, "Additional data to include (avatar_url, group_ids, locked, observed_users, can_be_removed, uuid, current_points)")

	// Get flags
	enrollmentsGetCmd.Flags().Int64Var(&enrollmentsCourseID, "course-id", 0, "Course ID (required)")
	enrollmentsGetCmd.MarkFlagRequired("course-id")

	// Create flags
	enrollmentsCreateCmd.Flags().Int64Var(&enrollmentsCourseID, "course-id", 0, "Course ID (required)")
	enrollmentsCreateCmd.Flags().Int64Var(&enrollmentsEnrollUserID, "user-id", 0, "User ID to enroll (required)")
	enrollmentsCreateCmd.Flags().StringVar(&enrollmentsEnrollType, "type", "StudentEnrollment", "Enrollment type")
	enrollmentsCreateCmd.Flags().StringVar(&enrollmentsEnrollState, "state", "", "Initial enrollment state (active, invited)")
	enrollmentsCreateCmd.Flags().Int64Var(&enrollmentsSectionID, "section-id", 0, "Section ID")
	enrollmentsCreateCmd.Flags().BoolVar(&enrollmentsNotify, "notify", false, "Send enrollment notification email")
	enrollmentsCreateCmd.Flags().StringVar(&enrollmentsRole, "role", "", "Custom role name")
	enrollmentsCreateCmd.MarkFlagRequired("course-id")
	enrollmentsCreateCmd.MarkFlagRequired("user-id")

	// Conclude flags
	enrollmentsConcludeCmd.Flags().Int64Var(&enrollmentsCourseID, "course-id", 0, "Course ID (required)")
	enrollmentsConcludeCmd.Flags().StringVar(&enrollmentsConcludeTask, "task", "conclude", "Task: conclude, deactivate, delete")
	enrollmentsConcludeCmd.Flags().BoolVar(&enrollmentsForce, "force", false, "Skip confirmation for delete")
	enrollmentsConcludeCmd.MarkFlagRequired("course-id")

	// Reactivate flags
	enrollmentsReactivateCmd.Flags().Int64Var(&enrollmentsCourseID, "course-id", 0, "Course ID (required)")
	enrollmentsReactivateCmd.MarkFlagRequired("course-id")

	// Accept flags
	enrollmentsAcceptCmd.Flags().Int64Var(&enrollmentsCourseID, "course-id", 0, "Course ID (required)")
	enrollmentsAcceptCmd.MarkFlagRequired("course-id")

	// Reject flags
	enrollmentsRejectCmd.Flags().Int64Var(&enrollmentsCourseID, "course-id", 0, "Course ID (required)")
	enrollmentsRejectCmd.MarkFlagRequired("course-id")
}

func runEnrollmentsList(cmd *cobra.Command, args []string) error {
	// Validate that exactly one context is specified
	contextsSpecified := 0
	if enrollmentsCourseID > 0 {
		contextsSpecified++
	}
	if enrollmentsUserID > 0 {
		contextsSpecified++
	}

	if contextsSpecified == 0 {
		return fmt.Errorf("must specify one of --course-id or --user-id")
	}
	if contextsSpecified > 1 {
		return fmt.Errorf("can only specify one of --course-id or --user-id")
	}

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create enrollments service
	enrollmentsService := api.NewEnrollmentsService(client)

	// Build options
	opts := &api.ListEnrollmentsOptions{
		Type:    enrollmentsType,
		State:   enrollmentsState,
		Include: enrollmentsInclude,
	}

	// List enrollments based on context
	ctx := context.Background()
	var enrollments []api.Enrollment
	var contextName string

	if enrollmentsCourseID > 0 {
		// Course context - list all enrollments in the course
		enrollments, err = enrollmentsService.ListCourse(ctx, enrollmentsCourseID, opts)
		contextName = fmt.Sprintf("course %d", enrollmentsCourseID)
	} else {
		// User context - list all enrollments for the user
		enrollments, err = enrollmentsService.ListUser(ctx, enrollmentsUserID, opts)
		contextName = fmt.Sprintf("user %d", enrollmentsUserID)
	}

	if err != nil {
		return fmt.Errorf("failed to list enrollments: %w", err)
	}

	if len(enrollments) == 0 {
		fmt.Printf("No enrollments found for %s\n", contextName)
		return nil
	}

	// Format and display enrollments
	printVerbose("Found %d enrollments for %s:\n\n", len(enrollments), contextName)
	return formatOutput(enrollments, nil)
}

func runEnrollmentsGet(cmd *cobra.Command, args []string) error {
	// Parse enrollment ID
	enrollmentID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid enrollment ID: %w", err)
	}

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create enrollments service
	enrollmentsService := api.NewEnrollmentsService(client)

	// Get enrollment by listing course enrollments and filtering
	// Note: Canvas doesn't have a direct "get enrollment by ID" endpoint
	// so we list all enrollments for the course and find the matching one
	ctx := context.Background()
	enrollments, err := enrollmentsService.ListCourse(ctx, enrollmentsCourseID, nil)
	if err != nil {
		return fmt.Errorf("failed to list enrollments: %w", err)
	}

	// Find the enrollment by ID
	var enrollment *api.Enrollment
	for i := range enrollments {
		if enrollments[i].ID == enrollmentID {
			enrollment = &enrollments[i]
			break
		}
	}

	if enrollment == nil {
		return fmt.Errorf("enrollment %d not found in course %d", enrollmentID, enrollmentsCourseID)
	}

	// Format and display enrollment details
	return formatOutput(enrollment, nil)
}

func runEnrollmentsCreate(cmd *cobra.Command, args []string) error {
	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create enrollments service
	enrollmentsService := api.NewEnrollmentsService(client)

	// Build params
	params := &api.EnrollUserParams{
		UserID:          enrollmentsEnrollUserID,
		Type:            enrollmentsEnrollType,
		EnrollmentState: enrollmentsEnrollState,
		CourseSectionID: enrollmentsSectionID,
		Notify:          enrollmentsNotify,
	}

	// Note: role flag maps to a custom role name, not RoleID
	// For custom roles, you'd need to look up the role ID first
	// For now, we support the standard enrollment types

	// Create enrollment
	ctx := context.Background()
	enrollment, err := enrollmentsService.EnrollUser(ctx, enrollmentsCourseID, params)
	if err != nil {
		return fmt.Errorf("failed to create enrollment: %w", err)
	}

	fmt.Printf("Enrollment created successfully (ID: %d)\n", enrollment.ID)
	return formatOutput(enrollment, nil)
}

func runEnrollmentsConclude(cmd *cobra.Command, args []string) error {
	// Parse enrollment ID
	enrollmentID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid enrollment ID: %w", err)
	}

	// Validate task
	switch enrollmentsConcludeTask {
	case "conclude", "deactivate", "delete":
		// Valid
	default:
		return fmt.Errorf("invalid task: %s (use 'conclude', 'deactivate', or 'delete')", enrollmentsConcludeTask)
	}

	// Confirmation for delete
	if enrollmentsConcludeTask == "delete" && !enrollmentsForce {
		fmt.Printf("WARNING: This will permanently delete enrollment %d.\n", enrollmentID)
		fmt.Print("Type 'yes' to confirm: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Delete cancelled")
			return nil
		}
	}

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create enrollments service
	enrollmentsService := api.NewEnrollmentsService(client)

	// Conclude/deactivate/delete enrollment
	ctx := context.Background()
	enrollment, err := enrollmentsService.Conclude(ctx, enrollmentsCourseID, enrollmentID, enrollmentsConcludeTask)
	if err != nil {
		return fmt.Errorf("failed to %s enrollment: %w", enrollmentsConcludeTask, err)
	}

	switch enrollmentsConcludeTask {
	case "conclude":
		fmt.Printf("Enrollment %d concluded\n", enrollmentID)
	case "deactivate":
		fmt.Printf("Enrollment %d deactivated\n", enrollmentID)
	case "delete":
		fmt.Printf("Enrollment %d deleted\n", enrollmentID)
	}

	return formatOutput(enrollment, nil)
}

func runEnrollmentsReactivate(cmd *cobra.Command, args []string) error {
	// Parse enrollment ID
	enrollmentID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid enrollment ID: %w", err)
	}

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create enrollments service
	enrollmentsService := api.NewEnrollmentsService(client)

	// Reactivate enrollment
	ctx := context.Background()
	enrollment, err := enrollmentsService.Reactivate(ctx, enrollmentsCourseID, enrollmentID)
	if err != nil {
		return fmt.Errorf("failed to reactivate enrollment: %w", err)
	}

	fmt.Printf("Enrollment %d reactivated\n", enrollmentID)
	return formatOutput(enrollment, nil)
}

func runEnrollmentsAccept(cmd *cobra.Command, args []string) error {
	// Parse enrollment ID
	enrollmentID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid enrollment ID: %w", err)
	}

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create enrollments service
	enrollmentsService := api.NewEnrollmentsService(client)

	// Accept enrollment invitation
	ctx := context.Background()
	err = enrollmentsService.Accept(ctx, enrollmentsCourseID, enrollmentID)
	if err != nil {
		return fmt.Errorf("failed to accept enrollment: %w", err)
	}

	fmt.Printf("Enrollment invitation %d accepted\n", enrollmentID)
	return nil
}

func runEnrollmentsReject(cmd *cobra.Command, args []string) error {
	// Parse enrollment ID
	enrollmentID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid enrollment ID: %w", err)
	}

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create enrollments service
	enrollmentsService := api.NewEnrollmentsService(client)

	// Reject enrollment invitation
	ctx := context.Background()
	err = enrollmentsService.Reject(ctx, enrollmentsCourseID, enrollmentID)
	if err != nil {
		return fmt.Errorf("failed to reject enrollment: %w", err)
	}

	fmt.Printf("Enrollment invitation %d rejected\n", enrollmentID)
	return nil
}
