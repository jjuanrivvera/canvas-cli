package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var (
	// List flags
	sectionsCourseID int64
	sectionsInclude  []string

	// Create/Update flags
	sectionsName          string
	sectionsSISSectionID  string
	sectionsIntegrationID string
	sectionsStartAt       string
	sectionsEndAt         string
	sectionsRestrictDates bool

	// Crosslist flags
	sectionsNewCourseID           int64
	sectionsOverrideSISStickiness bool

	// Delete flags
	sectionsForce bool
)

// sectionsCmd represents the sections command group
var sectionsCmd = &cobra.Command{
	Use:   "sections",
	Short: "Manage Canvas course sections",
	Long: `Manage Canvas course sections including listing, creating, updating, and deleting sections.

Sections allow you to organize students within a course into groups that can have
different due dates, grade visibility settings, or be crosslisted to other courses.

Examples:
  canvas sections list --course-id 123
  canvas sections get 456
  canvas sections create --course-id 123 --name "Section A"`,
}

// sectionsListCmd represents the sections list command
var sectionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List sections in a course",
	Long: `List all sections in a course.

Examples:
  canvas sections list --course-id 123
  canvas sections list --course-id 123 --include students,total_students
  canvas sections list --course-id 123 --include passback_status`,
	RunE: runSectionsList,
}

// sectionsGetCmd represents the sections get command
var sectionsGetCmd = &cobra.Command{
	Use:   "get <section-id>",
	Short: "Get section details",
	Long: `Get details of a specific section.

Examples:
  canvas sections get 456
  canvas sections get 456 --include students,total_students`,
	Args: ExactArgsWithUsage(1, "section-id"),
	RunE: runSectionsGet,
}

// sectionsCreateCmd represents the sections create command
var sectionsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new section",
	Long: `Create a new section in a course.

Examples:
  canvas sections create --course-id 123 --name "Section A"
  canvas sections create --course-id 123 --name "Section B" --sis-section-id "SIS123"
  canvas sections create --course-id 123 --name "Section C" --start-at "2024-01-15" --end-at "2024-05-15" --restrict-dates`,
	RunE: runSectionsCreate,
}

// sectionsUpdateCmd represents the sections update command
var sectionsUpdateCmd = &cobra.Command{
	Use:   "update <section-id>",
	Short: "Update a section",
	Long: `Update an existing section.

Examples:
  canvas sections update 456 --name "Updated Section Name"
  canvas sections update 456 --start-at "2024-02-01"
  canvas sections update 456 --restrict-dates`,
	Args: ExactArgsWithUsage(1, "section-id"),
	RunE: runSectionsUpdate,
}

// sectionsDeleteCmd represents the sections delete command
var sectionsDeleteCmd = &cobra.Command{
	Use:   "delete <section-id>",
	Short: "Delete a section",
	Long: `Delete a section.

WARNING: This action cannot be undone. All students in the section will be
removed from the course unless they are also enrolled in another section.

Examples:
  canvas sections delete 456
  canvas sections delete 456 --force`,
	Args: ExactArgsWithUsage(1, "section-id"),
	RunE: runSectionsDelete,
}

// sectionsCrosslistCmd represents the sections crosslist command
var sectionsCrosslistCmd = &cobra.Command{
	Use:   "crosslist <section-id>",
	Short: "Crosslist a section to another course",
	Long: `Move a section to a different course (crosslist).

When you crosslist a section, it is moved from its original course to a new course.
Students in the section will be enrolled in both courses.

Examples:
  canvas sections crosslist 456 --new-course-id 789
  canvas sections crosslist 456 --new-course-id 789 --override-sis-stickiness`,
	Args: ExactArgsWithUsage(1, "section-id"),
	RunE: runSectionsCrosslist,
}

// sectionsUncrosslistCmd represents the sections uncrosslist command
var sectionsUncrosslistCmd = &cobra.Command{
	Use:   "uncrosslist <section-id>",
	Short: "Return a crosslisted section to its original course",
	Long: `Return a crosslisted section to its original course.

Examples:
  canvas sections uncrosslist 456
  canvas sections uncrosslist 456 --override-sis-stickiness`,
	Args: ExactArgsWithUsage(1, "section-id"),
	RunE: runSectionsUncrosslist,
}

func init() {
	rootCmd.AddCommand(sectionsCmd)
	sectionsCmd.AddCommand(sectionsListCmd)
	sectionsCmd.AddCommand(sectionsGetCmd)
	sectionsCmd.AddCommand(sectionsCreateCmd)
	sectionsCmd.AddCommand(sectionsUpdateCmd)
	sectionsCmd.AddCommand(sectionsDeleteCmd)
	sectionsCmd.AddCommand(sectionsCrosslistCmd)
	sectionsCmd.AddCommand(sectionsUncrosslistCmd)

	// List flags
	sectionsListCmd.Flags().Int64Var(&sectionsCourseID, "course-id", 0, "Course ID (required)")
	sectionsListCmd.MarkFlagRequired("course-id")
	sectionsListCmd.Flags().StringSliceVar(&sectionsInclude, "include", []string{}, "Include additional data (students, total_students, passback_status, permissions)")

	// Get flags
	sectionsGetCmd.Flags().StringSliceVar(&sectionsInclude, "include", []string{}, "Include additional data (students, total_students, passback_status, permissions)")

	// Create flags
	sectionsCreateCmd.Flags().Int64Var(&sectionsCourseID, "course-id", 0, "Course ID (required)")
	sectionsCreateCmd.MarkFlagRequired("course-id")
	sectionsCreateCmd.Flags().StringVar(&sectionsName, "name", "", "Section name (required)")
	sectionsCreateCmd.MarkFlagRequired("name")
	sectionsCreateCmd.Flags().StringVar(&sectionsSISSectionID, "sis-section-id", "", "SIS section ID")
	sectionsCreateCmd.Flags().StringVar(&sectionsIntegrationID, "integration-id", "", "Integration ID")
	sectionsCreateCmd.Flags().StringVar(&sectionsStartAt, "start-at", "", "Section start date (ISO 8601)")
	sectionsCreateCmd.Flags().StringVar(&sectionsEndAt, "end-at", "", "Section end date (ISO 8601)")
	sectionsCreateCmd.Flags().BoolVar(&sectionsRestrictDates, "restrict-dates", false, "Restrict enrollments to section dates")

	// Update flags
	sectionsUpdateCmd.Flags().StringVar(&sectionsName, "name", "", "Section name")
	sectionsUpdateCmd.Flags().StringVar(&sectionsSISSectionID, "sis-section-id", "", "SIS section ID")
	sectionsUpdateCmd.Flags().StringVar(&sectionsIntegrationID, "integration-id", "", "Integration ID")
	sectionsUpdateCmd.Flags().StringVar(&sectionsStartAt, "start-at", "", "Section start date (ISO 8601)")
	sectionsUpdateCmd.Flags().StringVar(&sectionsEndAt, "end-at", "", "Section end date (ISO 8601)")
	sectionsUpdateCmd.Flags().BoolVar(&sectionsRestrictDates, "restrict-dates", false, "Restrict enrollments to section dates")
	sectionsUpdateCmd.Flags().BoolVar(&sectionsOverrideSISStickiness, "override-sis-stickiness", false, "Override SIS stickiness")

	// Delete flags
	sectionsDeleteCmd.Flags().BoolVar(&sectionsForce, "force", false, "Skip confirmation prompt")

	// Crosslist flags
	sectionsCrosslistCmd.Flags().Int64Var(&sectionsNewCourseID, "new-course-id", 0, "Target course ID (required)")
	sectionsCrosslistCmd.MarkFlagRequired("new-course-id")
	sectionsCrosslistCmd.Flags().BoolVar(&sectionsOverrideSISStickiness, "override-sis-stickiness", false, "Override SIS stickiness")

	// Uncrosslist flags
	sectionsUncrosslistCmd.Flags().BoolVar(&sectionsOverrideSISStickiness, "override-sis-stickiness", false, "Override SIS stickiness")
}

func runSectionsList(cmd *cobra.Command, args []string) error {
	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create sections service
	sectionsService := api.NewSectionsService(client)

	// Build options
	opts := &api.ListSectionsOptions{
		Include: sectionsInclude,
	}

	// List sections
	ctx := context.Background()
	sections, err := sectionsService.ListCourse(ctx, sectionsCourseID, opts)
	if err != nil {
		return fmt.Errorf("failed to list sections: %w", err)
	}

	if len(sections) == 0 {
		fmt.Printf("No sections found in course %d\n", sectionsCourseID)
		return nil
	}

	printVerbose("Found %d sections in course %d:\n\n", len(sections), sectionsCourseID)
	return formatOutput(sections, nil)
}

func runSectionsGet(cmd *cobra.Command, args []string) error {
	// Parse section ID
	sectionID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid section ID: %w", err)
	}

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create sections service
	sectionsService := api.NewSectionsService(client)

	// Get section
	ctx := context.Background()
	section, err := sectionsService.Get(ctx, sectionID, sectionsInclude)
	if err != nil {
		return fmt.Errorf("failed to get section: %w", err)
	}

	return formatOutput(section, nil)
}

func runSectionsCreate(cmd *cobra.Command, args []string) error {
	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create sections service
	sectionsService := api.NewSectionsService(client)

	// Build params
	params := &api.CreateSectionParams{
		Name:                              sectionsName,
		SISSectionID:                      sectionsSISSectionID,
		IntegrationID:                     sectionsIntegrationID,
		StartAt:                           sectionsStartAt,
		EndAt:                             sectionsEndAt,
		RestrictEnrollmentsToSectionDates: sectionsRestrictDates,
	}

	// Create section
	ctx := context.Background()
	section, err := sectionsService.Create(ctx, sectionsCourseID, params)
	if err != nil {
		return fmt.Errorf("failed to create section: %w", err)
	}

	fmt.Printf("Section created successfully (ID: %d)\n", section.ID)
	return formatOutput(section, nil)
}

func runSectionsUpdate(cmd *cobra.Command, args []string) error {
	// Parse section ID
	sectionID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid section ID: %w", err)
	}

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create sections service
	sectionsService := api.NewSectionsService(client)

	// Build params - only include changed flags
	params := &api.UpdateSectionParams{
		OverrideSISStickiness: sectionsOverrideSISStickiness,
	}

	if cmd.Flags().Changed("name") {
		params.Name = &sectionsName
	}
	if cmd.Flags().Changed("sis-section-id") {
		params.SISSectionID = &sectionsSISSectionID
	}
	if cmd.Flags().Changed("integration-id") {
		params.IntegrationID = &sectionsIntegrationID
	}
	if cmd.Flags().Changed("start-at") {
		params.StartAt = &sectionsStartAt
	}
	if cmd.Flags().Changed("end-at") {
		params.EndAt = &sectionsEndAt
	}
	if cmd.Flags().Changed("restrict-dates") {
		params.RestrictEnrollmentsToSectionDates = &sectionsRestrictDates
	}

	// Update section
	ctx := context.Background()
	section, err := sectionsService.Update(ctx, sectionID, params)
	if err != nil {
		return fmt.Errorf("failed to update section: %w", err)
	}

	fmt.Printf("Section updated successfully (ID: %d)\n", section.ID)
	return formatOutput(section, nil)
}

func runSectionsDelete(cmd *cobra.Command, args []string) error {
	// Parse section ID
	sectionID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid section ID: %w", err)
	}

	// Confirmation
	if !sectionsForce {
		fmt.Printf("WARNING: This will delete section %d and may remove students from the course.\n", sectionID)
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

	// Create sections service
	sectionsService := api.NewSectionsService(client)

	// Delete section
	ctx := context.Background()
	section, err := sectionsService.Delete(ctx, sectionID)
	if err != nil {
		return fmt.Errorf("failed to delete section: %w", err)
	}

	fmt.Printf("Section %d deleted\n", section.ID)
	return nil
}

func runSectionsCrosslist(cmd *cobra.Command, args []string) error {
	// Parse section ID
	sectionID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid section ID: %w", err)
	}

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create sections service
	sectionsService := api.NewSectionsService(client)

	// Crosslist section
	ctx := context.Background()
	section, err := sectionsService.Crosslist(ctx, sectionID, sectionsNewCourseID, sectionsOverrideSISStickiness)
	if err != nil {
		return fmt.Errorf("failed to crosslist section: %w", err)
	}

	fmt.Printf("Section %d crosslisted to course %d\n", section.ID, section.CourseID)
	return formatOutput(section, nil)
}

func runSectionsUncrosslist(cmd *cobra.Command, args []string) error {
	// Parse section ID
	sectionID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid section ID: %w", err)
	}

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create sections service
	sectionsService := api.NewSectionsService(client)

	// Uncrosslist section
	ctx := context.Background()
	section, err := sectionsService.Uncrosslist(ctx, sectionID, sectionsOverrideSISStickiness)
	if err != nil {
		return fmt.Errorf("failed to uncrosslist section: %w", err)
	}

	fmt.Printf("Section %d returned to course %d\n", section.ID, section.CourseID)
	return formatOutput(section, nil)
}
