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

func init() {
	rootCmd.AddCommand(sectionsCmd)
	sectionsCmd.AddCommand(newSectionsListCmd())
	sectionsCmd.AddCommand(newSectionsGetCmd())
	sectionsCmd.AddCommand(newSectionsCreateCmd())
	sectionsCmd.AddCommand(newSectionsUpdateCmd())
	sectionsCmd.AddCommand(newSectionsDeleteCmd())
	sectionsCmd.AddCommand(newSectionsCrosslistCmd())
	sectionsCmd.AddCommand(newSectionsUncrosslistCmd())
}

func newSectionsListCmd() *cobra.Command {
	opts := &options.SectionsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List sections in a course",
		Long: `List all sections in a course.

Examples:
  canvas sections list --course-id 123
  canvas sections list --course-id 123 --include students,total_students
  canvas sections list --course-id 123 --include passback_status`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runSectionsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Include additional data (students, total_students, passback_status, permissions)")

	return cmd
}

func newSectionsGetCmd() *cobra.Command {
	opts := &options.SectionsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <section-id>",
		Short: "Get section details",
		Long: `Get details of a specific section.

Examples:
  canvas sections get 456
  canvas sections get 456 --include students,total_students`,
		Args: ExactArgsWithUsage(1, "section-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			sectionID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid section ID: %w", err)
			}
			opts.SectionID = sectionID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runSectionsGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Include additional data (students, total_students, passback_status, permissions)")

	return cmd
}

func newSectionsCreateCmd() *cobra.Command {
	opts := &options.SectionsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new section",
		Long: `Create a new section in a course.

Examples:
  canvas sections create --course-id 123 --name "Section A"
  canvas sections create --course-id 123 --name "Section B" --sis-section-id "SIS123"
  canvas sections create --course-id 123 --name "Section C" --start-at "2024-01-15" --end-at "2024-05-15" --restrict-dates`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runSectionsCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Section name (required)")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVar(&opts.SISSectionID, "sis-section-id", "", "SIS section ID")
	cmd.Flags().StringVar(&opts.IntegrationID, "integration-id", "", "Integration ID")
	cmd.Flags().StringVar(&opts.StartAt, "start-at", "", "Section start date (ISO 8601)")
	cmd.Flags().StringVar(&opts.EndAt, "end-at", "", "Section end date (ISO 8601)")
	cmd.Flags().BoolVar(&opts.RestrictDates, "restrict-dates", false, "Restrict enrollments to section dates")

	return cmd
}

func newSectionsUpdateCmd() *cobra.Command {
	opts := &options.SectionsUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <section-id>",
		Short: "Update a section",
		Long: `Update an existing section.

Examples:
  canvas sections update 456 --name "Updated Section Name"
  canvas sections update 456 --start-at "2024-02-01"
  canvas sections update 456 --restrict-dates`,
		Args: ExactArgsWithUsage(1, "section-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			sectionID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid section ID: %w", err)
			}
			opts.SectionID = sectionID

			// Track which fields were set
			opts.NameSet = cmd.Flags().Changed("name")
			opts.SISSectionIDSet = cmd.Flags().Changed("sis-section-id")
			opts.IntegrationIDSet = cmd.Flags().Changed("integration-id")
			opts.StartAtSet = cmd.Flags().Changed("start-at")
			opts.EndAtSet = cmd.Flags().Changed("end-at")
			opts.RestrictDatesSet = cmd.Flags().Changed("restrict-dates")

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runSectionsUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "Section name")
	cmd.Flags().StringVar(&opts.SISSectionID, "sis-section-id", "", "SIS section ID")
	cmd.Flags().StringVar(&opts.IntegrationID, "integration-id", "", "Integration ID")
	cmd.Flags().StringVar(&opts.StartAt, "start-at", "", "Section start date (ISO 8601)")
	cmd.Flags().StringVar(&opts.EndAt, "end-at", "", "Section end date (ISO 8601)")
	cmd.Flags().BoolVar(&opts.RestrictDates, "restrict-dates", false, "Restrict enrollments to section dates")
	cmd.Flags().BoolVar(&opts.OverrideSISStickiness, "override-sis-stickiness", false, "Override SIS stickiness")

	return cmd
}

func newSectionsDeleteCmd() *cobra.Command {
	opts := &options.SectionsDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <section-id>",
		Short: "Delete a section",
		Long: `Delete a section.

WARNING: This action cannot be undone. All students in the section will be
removed from the course unless they are also enrolled in another section.

Examples:
  canvas sections delete 456
  canvas sections delete 456 --force`,
		Args: ExactArgsWithUsage(1, "section-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			sectionID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid section ID: %w", err)
			}
			opts.SectionID = sectionID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runSectionsDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Force, "force", false, "Skip confirmation prompt")

	return cmd
}

func newSectionsCrosslistCmd() *cobra.Command {
	opts := &options.SectionsCrosslistOptions{}

	cmd := &cobra.Command{
		Use:   "crosslist <section-id>",
		Short: "Crosslist a section to another course",
		Long: `Move a section to a different course (crosslist).

When you crosslist a section, it is moved from its original course to a new course.
Students in the section will be enrolled in both courses.

Examples:
  canvas sections crosslist 456 --new-course-id 789
  canvas sections crosslist 456 --new-course-id 789 --override-sis-stickiness`,
		Args: ExactArgsWithUsage(1, "section-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			sectionID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid section ID: %w", err)
			}
			opts.SectionID = sectionID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runSectionsCrosslist(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.NewCourseID, "new-course-id", 0, "Target course ID (required)")
	cmd.MarkFlagRequired("new-course-id")
	cmd.Flags().BoolVar(&opts.OverrideSISStickiness, "override-sis-stickiness", false, "Override SIS stickiness")

	return cmd
}

func newSectionsUncrosslistCmd() *cobra.Command {
	opts := &options.SectionsUncrosslistOptions{}

	cmd := &cobra.Command{
		Use:   "uncrosslist <section-id>",
		Short: "Return a crosslisted section to its original course",
		Long: `Return a crosslisted section to its original course.

Examples:
  canvas sections uncrosslist 456
  canvas sections uncrosslist 456 --override-sis-stickiness`,
		Args: ExactArgsWithUsage(1, "section-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			sectionID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid section ID: %w", err)
			}
			opts.SectionID = sectionID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runSectionsUncrosslist(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().BoolVar(&opts.OverrideSISStickiness, "override-sis-stickiness", false, "Override SIS stickiness")

	return cmd
}

func runSectionsList(ctx context.Context, client *api.Client, opts *options.SectionsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "sections.list", map[string]interface{}{
		"course_id": opts.CourseID,
		"include":   opts.Include,
	})

	sectionsService := api.NewSectionsService(client)

	apiOpts := &api.ListSectionsOptions{
		Include: opts.Include,
	}

	sections, err := sectionsService.ListCourse(ctx, opts.CourseID, apiOpts)
	if err != nil {
		logger.LogCommandError(ctx, "sections.list", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to list sections: %w", err)
	}

	if len(sections) == 0 {
		fmt.Printf("No sections found in course %d\n", opts.CourseID)
		logger.LogCommandComplete(ctx, "sections.list", 0)
		return nil
	}

	printVerbose("Found %d sections in course %d:\n\n", len(sections), opts.CourseID)
	logger.LogCommandComplete(ctx, "sections.list", len(sections))
	return formatOutput(sections, nil)
}

func runSectionsGet(ctx context.Context, client *api.Client, opts *options.SectionsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "sections.get", map[string]interface{}{
		"section_id": opts.SectionID,
		"include":    opts.Include,
	})

	sectionsService := api.NewSectionsService(client)

	section, err := sectionsService.Get(ctx, opts.SectionID, opts.Include)
	if err != nil {
		logger.LogCommandError(ctx, "sections.get", err, map[string]interface{}{
			"section_id": opts.SectionID,
		})
		return fmt.Errorf("failed to get section: %w", err)
	}

	logger.LogCommandComplete(ctx, "sections.get", 1)
	return formatOutput(section, nil)
}

func runSectionsCreate(ctx context.Context, client *api.Client, opts *options.SectionsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "sections.create", map[string]interface{}{
		"course_id": opts.CourseID,
		"name":      opts.Name,
	})

	sectionsService := api.NewSectionsService(client)

	params := &api.CreateSectionParams{
		Name:                              opts.Name,
		SISSectionID:                      opts.SISSectionID,
		IntegrationID:                     opts.IntegrationID,
		StartAt:                           opts.StartAt,
		EndAt:                             opts.EndAt,
		RestrictEnrollmentsToSectionDates: opts.RestrictDates,
	}

	section, err := sectionsService.Create(ctx, opts.CourseID, params)
	if err != nil {
		logger.LogCommandError(ctx, "sections.create", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"name":      opts.Name,
		})
		return fmt.Errorf("failed to create section: %w", err)
	}

	fmt.Printf("Section created successfully (ID: %d)\n", section.ID)
	logger.LogCommandComplete(ctx, "sections.create", 1)
	return formatOutput(section, nil)
}

func runSectionsUpdate(ctx context.Context, client *api.Client, opts *options.SectionsUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "sections.update", map[string]interface{}{
		"section_id": opts.SectionID,
	})

	sectionsService := api.NewSectionsService(client)

	params := &api.UpdateSectionParams{
		OverrideSISStickiness: opts.OverrideSISStickiness,
	}

	if opts.NameSet {
		params.Name = &opts.Name
	}
	if opts.SISSectionIDSet {
		params.SISSectionID = &opts.SISSectionID
	}
	if opts.IntegrationIDSet {
		params.IntegrationID = &opts.IntegrationID
	}
	if opts.StartAtSet {
		params.StartAt = &opts.StartAt
	}
	if opts.EndAtSet {
		params.EndAt = &opts.EndAt
	}
	if opts.RestrictDatesSet {
		params.RestrictEnrollmentsToSectionDates = &opts.RestrictDates
	}

	section, err := sectionsService.Update(ctx, opts.SectionID, params)
	if err != nil {
		logger.LogCommandError(ctx, "sections.update", err, map[string]interface{}{
			"section_id": opts.SectionID,
		})
		return fmt.Errorf("failed to update section: %w", err)
	}

	fmt.Printf("Section updated successfully (ID: %d)\n", section.ID)
	logger.LogCommandComplete(ctx, "sections.update", 1)
	return formatOutput(section, nil)
}

func runSectionsDelete(ctx context.Context, client *api.Client, opts *options.SectionsDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "sections.delete", map[string]interface{}{
		"section_id": opts.SectionID,
		"force":      opts.Force,
	})

	// Confirmation
	if !opts.Force {
		fmt.Printf("WARNING: This will delete section %d and may remove students from the course.\n", opts.SectionID)
		fmt.Print("Type 'yes' to confirm: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Delete cancelled")
			return nil
		}
	}

	sectionsService := api.NewSectionsService(client)

	section, err := sectionsService.Delete(ctx, opts.SectionID)
	if err != nil {
		logger.LogCommandError(ctx, "sections.delete", err, map[string]interface{}{
			"section_id": opts.SectionID,
		})
		return fmt.Errorf("failed to delete section: %w", err)
	}

	fmt.Printf("Section %d deleted\n", section.ID)
	logger.LogCommandComplete(ctx, "sections.delete", 1)
	return nil
}

func runSectionsCrosslist(ctx context.Context, client *api.Client, opts *options.SectionsCrosslistOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "sections.crosslist", map[string]interface{}{
		"section_id":              opts.SectionID,
		"new_course_id":           opts.NewCourseID,
		"override_sis_stickiness": opts.OverrideSISStickiness,
	})

	sectionsService := api.NewSectionsService(client)

	section, err := sectionsService.Crosslist(ctx, opts.SectionID, opts.NewCourseID, opts.OverrideSISStickiness)
	if err != nil {
		logger.LogCommandError(ctx, "sections.crosslist", err, map[string]interface{}{
			"section_id":    opts.SectionID,
			"new_course_id": opts.NewCourseID,
		})
		return fmt.Errorf("failed to crosslist section: %w", err)
	}

	fmt.Printf("Section %d crosslisted to course %d\n", section.ID, section.CourseID)
	logger.LogCommandComplete(ctx, "sections.crosslist", 1)
	return formatOutput(section, nil)
}

func runSectionsUncrosslist(ctx context.Context, client *api.Client, opts *options.SectionsUncrosslistOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "sections.uncrosslist", map[string]interface{}{
		"section_id":              opts.SectionID,
		"override_sis_stickiness": opts.OverrideSISStickiness,
	})

	sectionsService := api.NewSectionsService(client)

	section, err := sectionsService.Uncrosslist(ctx, opts.SectionID, opts.OverrideSISStickiness)
	if err != nil {
		logger.LogCommandError(ctx, "sections.uncrosslist", err, map[string]interface{}{
			"section_id": opts.SectionID,
		})
		return fmt.Errorf("failed to uncrosslist section: %w", err)
	}

	fmt.Printf("Section %d returned to course %d\n", section.ID, section.CourseID)
	logger.LogCommandComplete(ctx, "sections.uncrosslist", 1)
	return formatOutput(section, nil)
}
