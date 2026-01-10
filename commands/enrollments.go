package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
	"github.com/spf13/cobra"
)

var (
	enrollmentsCourseID int64
	enrollmentsUserID   int64
	enrollmentsType     []string
	enrollmentsState    []string
	enrollmentsInclude  []string
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
	Args: cobra.ExactArgs(1),
	RunE: runEnrollmentsGet,
}

func init() {
	rootCmd.AddCommand(enrollmentsCmd)
	enrollmentsCmd.AddCommand(enrollmentsListCmd)
	enrollmentsCmd.AddCommand(enrollmentsGetCmd)

	// List flags
	enrollmentsListCmd.Flags().Int64Var(&enrollmentsCourseID, "course-id", 0, "Course ID (for course enrollments)")
	enrollmentsListCmd.Flags().Int64Var(&enrollmentsUserID, "user-id", 0, "User ID (for user enrollments)")
	enrollmentsListCmd.Flags().StringSliceVar(&enrollmentsType, "type", []string{}, "Filter by enrollment type (StudentEnrollment, TeacherEnrollment, TaEnrollment, ObserverEnrollment, DesignerEnrollment)")
	enrollmentsListCmd.Flags().StringSliceVar(&enrollmentsState, "state", []string{}, "Filter by enrollment state (active, invited, creation_pending, deleted, rejected, completed, inactive)")
	enrollmentsListCmd.Flags().StringSliceVar(&enrollmentsInclude, "include", []string{}, "Additional data to include (avatar_url, group_ids, locked, observed_users, can_be_removed, uuid, current_points)")

	// Get flags
	enrollmentsGetCmd.Flags().Int64Var(&enrollmentsCourseID, "course-id", 0, "Course ID (required)")
	enrollmentsGetCmd.MarkFlagRequired("course-id")
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

	// Display enrollments
	fmt.Printf("Found %d enrollments for %s:\n\n", len(enrollments), contextName)

	for _, enrollment := range enrollments {
		// Get enrollment type emoji
		typeEmoji := "📋"
		switch enrollment.Type {
		case "StudentEnrollment":
			typeEmoji = "🎓"
		case "TeacherEnrollment":
			typeEmoji = "👨‍🏫"
		case "TaEnrollment":
			typeEmoji = "👤"
		case "ObserverEnrollment":
			typeEmoji = "👁️"
		case "DesignerEnrollment":
			typeEmoji = "🎨"
		}

		fmt.Printf("%s Enrollment %d\n", typeEmoji, enrollment.ID)
		fmt.Printf("   Type: %s\n", enrollment.Type)
		fmt.Printf("   State: %s\n", enrollment.EnrollmentState)

		if enrollment.User != nil {
			fmt.Printf("   User: %s (ID: %d)\n", enrollment.User.Name, enrollment.User.ID)
		}

		if enrollment.CourseID > 0 {
			fmt.Printf("   Course ID: %d\n", enrollment.CourseID)
		}

		if enrollment.CourseSectionID > 0 {
			fmt.Printf("   Section ID: %d\n", enrollment.CourseSectionID)
		}

		if enrollment.Role != "" {
			fmt.Printf("   Role: %s\n", enrollment.Role)
		}

		fmt.Println()
	}

	return nil
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

	// Display enrollment details
	fmt.Printf("Enrollment Details\n")
	fmt.Printf("==================\n\n")
	fmt.Printf("ID:             %d\n", enrollment.ID)
	fmt.Printf("Type:           %s\n", enrollment.Type)
	fmt.Printf("State:          %s\n", enrollment.EnrollmentState)
	fmt.Printf("Course ID:      %d\n", enrollment.CourseID)

	if enrollment.CourseSectionID > 0 {
		fmt.Printf("Section ID:     %d\n", enrollment.CourseSectionID)
	}

	if enrollment.User != nil {
		fmt.Printf("User:           %s (ID: %d)\n", enrollment.User.Name, enrollment.User.ID)
	}

	if enrollment.Role != "" {
		fmt.Printf("Role:           %s\n", enrollment.Role)
	}

	if enrollment.RoleID > 0 {
		fmt.Printf("Role ID:        %d\n", enrollment.RoleID)
	}

	if !enrollment.CreatedAt.IsZero() {
		fmt.Printf("Created:        %s\n", enrollment.CreatedAt.Format("2006-01-02 15:04"))
	}

	if !enrollment.UpdatedAt.IsZero() {
		fmt.Printf("Updated:        %s\n", enrollment.UpdatedAt.Format("2006-01-02 15:04"))
	}

	if !enrollment.StartAt.IsZero() {
		fmt.Printf("Starts:         %s\n", enrollment.StartAt.Format("2006-01-02 15:04"))
	}

	if !enrollment.EndAt.IsZero() {
		fmt.Printf("Ends:           %s\n", enrollment.EndAt.Format("2006-01-02 15:04"))
	}

	if enrollment.Grades != nil {
		fmt.Printf("\nGrades:\n")
		if enrollment.Grades.CurrentScore > 0 {
			fmt.Printf("   Current Score: %.2f\n", enrollment.Grades.CurrentScore)
		}
		if enrollment.Grades.FinalScore > 0 {
			fmt.Printf("   Final Score:   %.2f\n", enrollment.Grades.FinalScore)
		}
		if enrollment.Grades.CurrentGrade != "" {
			fmt.Printf("   Current Grade: %s\n", enrollment.Grades.CurrentGrade)
		}
		if enrollment.Grades.FinalGrade != "" {
			fmt.Printf("   Final Grade:   %s\n", enrollment.Grades.FinalGrade)
		}
	}

	return nil
}
