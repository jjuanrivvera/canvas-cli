package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// contentMigrationsCmd represents the content-migrations command group
var contentMigrationsCmd = &cobra.Command{
	Use:     "content-migrations",
	Aliases: []string{"migrations", "cm"},
	Short:   "Manage content migrations",
	Long: `Manage Canvas content migrations.

Content migrations allow you to copy content between courses,
import content from external sources, and export course content.

Examples:
  canvas content-migrations list --course-id 1
  canvas content-migrations get 123 --course-id 1
  canvas content-migrations create --course-id 1 --type course_copy_importer --source-course-id 100`,
}

func init() {
	rootCmd.AddCommand(contentMigrationsCmd)
	contentMigrationsCmd.AddCommand(newContentMigrationsListCmd())
	contentMigrationsCmd.AddCommand(newContentMigrationsGetCmd())
	contentMigrationsCmd.AddCommand(newContentMigrationsCreateCmd())
	contentMigrationsCmd.AddCommand(newContentMigrationsMigratorsCmd())
	contentMigrationsCmd.AddCommand(newContentMigrationsContentCmd())
	contentMigrationsCmd.AddCommand(newContentMigrationsIssuesCmd())
}

func newContentMigrationsListCmd() *cobra.Command {
	opts := &options.ContentMigrationsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List content migrations",
		Long: `List content migrations for a course.

Examples:
  canvas content-migrations list --course-id 1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runContentMigrationsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newContentMigrationsGetCmd() *cobra.Command {
	opts := &options.ContentMigrationsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <migration-id>",
		Short: "Get a content migration",
		Long: `Get details of a specific content migration.

Examples:
  canvas content-migrations get 123 --course-id 1`,
		Args: ExactArgsWithUsage(1, "migration-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			migrationID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid migration ID: %w", err)
			}
			opts.MigrationID = migrationID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runContentMigrationsGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newContentMigrationsCreateCmd() *cobra.Command {
	opts := &options.ContentMigrationsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a content migration",
		Long: `Create a new content migration.

Migration types:
  - course_copy_importer: Copy content from another course
  - common_cartridge_importer: Import Common Cartridge file
  - zip_file_importer: Import ZIP file
  - canvas_cartridge_importer: Import Canvas export file
  - qti_importer: Import QTI quiz file

Examples:
  canvas content-migrations create --course-id 1 --type course_copy_importer --source-course-id 100
  canvas content-migrations create --course-id 1 --type common_cartridge_importer --file export.imscc`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Track which fields were set
			opts.SourceCourseIDSet = cmd.Flags().Changed("source-course-id")
			opts.FolderIDSet = cmd.Flags().Changed("folder-id")
			opts.SelectiveSet = cmd.Flags().Changed("selective")

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runContentMigrationsCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.Type, "type", "", "Migration type (required)")
	cmd.Flags().Int64Var(&opts.SourceCourseID, "source-course-id", 0, "Source course ID (for course_copy_importer)")
	cmd.Flags().StringVar(&opts.File, "file", "", "Export file to import")
	cmd.Flags().StringVar(&opts.FileURL, "file-url", "", "URL to export file")
	cmd.Flags().Int64Var(&opts.FolderID, "folder-id", 0, "Target folder ID")
	cmd.Flags().BoolVar(&opts.Selective, "selective", false, "Enable selective import")
	cmd.Flags().StringVar(&opts.CopyOptions, "copy-options", "", "JSON with copy options")
	cmd.Flags().StringVar(&opts.DateShift, "date-shift", "", "JSON with date shift options")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("type")

	return cmd
}

func newContentMigrationsMigratorsCmd() *cobra.Command {
	opts := &options.ContentMigrationsMigratorsOptions{}

	cmd := &cobra.Command{
		Use:   "migrators",
		Short: "List available migration types",
		Long: `List available migrator types for a course.

This shows what types of content migrations can be performed.

Examples:
  canvas content-migrations migrators --course-id 1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runContentMigrationsMigrators(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newContentMigrationsContentCmd() *cobra.Command {
	opts := &options.ContentMigrationsContentOptions{}

	cmd := &cobra.Command{
		Use:   "content <migration-id>",
		Short: "List migration content for selective import",
		Long: `List available content for selective import.

This shows what content is available to import from a migration.

Examples:
  canvas content-migrations content 123 --course-id 1
  canvas content-migrations content 123 --course-id 1 --type assignments`,
		Args: ExactArgsWithUsage(1, "migration-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			migrationID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid migration ID: %w", err)
			}
			opts.MigrationID = migrationID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runContentMigrationsContent(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.ContentType, "type", "", "Content type filter")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newContentMigrationsIssuesCmd() *cobra.Command {
	opts := &options.ContentMigrationsIssuesOptions{}

	cmd := &cobra.Command{
		Use:   "issues <migration-id>",
		Short: "List migration issues",
		Long: `List issues encountered during a content migration.

Examples:
  canvas content-migrations issues 123 --course-id 1`,
		Args: ExactArgsWithUsage(1, "migration-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			migrationID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid migration ID: %w", err)
			}
			opts.MigrationID = migrationID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runContentMigrationsIssues(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func runContentMigrationsList(ctx context.Context, client *api.Client, opts *options.ContentMigrationsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "content_migrations.list", map[string]interface{}{
		"course_id": opts.CourseID,
	})

	service := api.NewContentMigrationsService(client)

	migrations, err := service.List(ctx, opts.CourseID, nil)
	if err != nil {
		logger.LogCommandError(ctx, "content_migrations.list", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to list content migrations: %w", err)
	}

	if len(migrations) == 0 {
		fmt.Println("No content migrations found")
		logger.LogCommandComplete(ctx, "content_migrations.list", 0)
		return nil
	}

	printVerbose("Found %d content migrations:\n\n", len(migrations))
	logger.LogCommandComplete(ctx, "content_migrations.list", len(migrations))
	return formatOutput(migrations, nil)
}

func runContentMigrationsGet(ctx context.Context, client *api.Client, opts *options.ContentMigrationsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "content_migrations.get", map[string]interface{}{
		"course_id":    opts.CourseID,
		"migration_id": opts.MigrationID,
	})

	service := api.NewContentMigrationsService(client)

	migration, err := service.Get(ctx, opts.CourseID, opts.MigrationID)
	if err != nil {
		logger.LogCommandError(ctx, "content_migrations.get", err, map[string]interface{}{
			"course_id":    opts.CourseID,
			"migration_id": opts.MigrationID,
		})
		return fmt.Errorf("failed to get content migration: %w", err)
	}

	logger.LogCommandComplete(ctx, "content_migrations.get", 1)
	return formatOutput(migration, nil)
}

func runContentMigrationsCreate(ctx context.Context, client *api.Client, opts *options.ContentMigrationsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "content_migrations.create", map[string]interface{}{
		"course_id": opts.CourseID,
		"type":      opts.Type,
	})

	params := &api.CreateContentMigrationParams{
		MigrationType: opts.Type,
	}

	if opts.SourceCourseIDSet {
		params.SourceCourseID = &opts.SourceCourseID
	}

	if opts.File != "" {
		params.FilePath = opts.File
	}

	if opts.FileURL != "" {
		params.FileURL = opts.FileURL
	}

	if opts.FolderIDSet {
		params.FolderID = &opts.FolderID
	}

	if opts.SelectiveSet {
		params.SelectiveImport = &opts.Selective
	}

	if opts.CopyOptions != "" {
		var copyOpts map[string]interface{}
		if err := json.Unmarshal([]byte(opts.CopyOptions), &copyOpts); err != nil {
			logger.LogCommandError(ctx, "content_migrations.create", err, map[string]interface{}{
				"course_id": opts.CourseID,
				"type":      opts.Type,
			})
			return fmt.Errorf("invalid copy options JSON: %w", err)
		}
		params.CopyOptions = copyOpts
	}

	if opts.DateShift != "" {
		var dateShift api.DateShiftOptions
		if err := json.Unmarshal([]byte(opts.DateShift), &dateShift); err != nil {
			logger.LogCommandError(ctx, "content_migrations.create", err, map[string]interface{}{
				"course_id": opts.CourseID,
				"type":      opts.Type,
			})
			return fmt.Errorf("invalid date shift JSON: %w", err)
		}
		params.DateShiftOptions = &dateShift
	}

	service := api.NewContentMigrationsService(client)

	migration, err := service.Create(ctx, opts.CourseID, params)
	if err != nil {
		logger.LogCommandError(ctx, "content_migrations.create", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"type":      opts.Type,
		})
		return fmt.Errorf("failed to create content migration: %w", err)
	}

	fmt.Printf("Content migration created (ID: %d)\n", migration.ID)
	fmt.Printf("Type: %s\n", migration.MigrationType)
	fmt.Printf("State: %s\n", migration.WorkflowState)
	logger.LogCommandComplete(ctx, "content_migrations.create", 1)
	return nil
}

func runContentMigrationsMigrators(ctx context.Context, client *api.Client, opts *options.ContentMigrationsMigratorsOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "content_migrations.migrators", map[string]interface{}{
		"course_id": opts.CourseID,
	})

	service := api.NewContentMigrationsService(client)

	migrators, err := service.ListMigrators(ctx, opts.CourseID)
	if err != nil {
		logger.LogCommandError(ctx, "content_migrations.migrators", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to list migrators: %w", err)
	}

	if len(migrators) == 0 {
		fmt.Println("No migrators available")
		logger.LogCommandComplete(ctx, "content_migrations.migrators", 0)
		return nil
	}

	printVerbose("Available migrators:\n\n")
	logger.LogCommandComplete(ctx, "content_migrations.migrators", len(migrators))
	return formatOutput(migrators, nil)
}

func runContentMigrationsContent(ctx context.Context, client *api.Client, opts *options.ContentMigrationsContentOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "content_migrations.content", map[string]interface{}{
		"course_id":    opts.CourseID,
		"migration_id": opts.MigrationID,
		"content_type": opts.ContentType,
	})

	service := api.NewContentMigrationsService(client)

	items, err := service.ListContentList(ctx, opts.CourseID, opts.MigrationID, opts.ContentType)
	if err != nil {
		logger.LogCommandError(ctx, "content_migrations.content", err, map[string]interface{}{
			"course_id":    opts.CourseID,
			"migration_id": opts.MigrationID,
		})
		return fmt.Errorf("failed to list content: %w", err)
	}

	if len(items) == 0 {
		fmt.Println("No content available for import")
		logger.LogCommandComplete(ctx, "content_migrations.content", 0)
		return nil
	}

	printVerbose("Available content:\n\n")
	logger.LogCommandComplete(ctx, "content_migrations.content", len(items))
	return formatOutput(items, nil)
}

func runContentMigrationsIssues(ctx context.Context, client *api.Client, opts *options.ContentMigrationsIssuesOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "content_migrations.issues", map[string]interface{}{
		"course_id":    opts.CourseID,
		"migration_id": opts.MigrationID,
	})

	service := api.NewContentMigrationsService(client)

	issues, err := service.ListMigrationIssues(ctx, opts.CourseID, opts.MigrationID)
	if err != nil {
		logger.LogCommandError(ctx, "content_migrations.issues", err, map[string]interface{}{
			"course_id":    opts.CourseID,
			"migration_id": opts.MigrationID,
		})
		return fmt.Errorf("failed to list migration issues: %w", err)
	}

	if len(issues) == 0 {
		fmt.Println("No issues found for this migration")
		logger.LogCommandComplete(ctx, "content_migrations.issues", 0)
		return nil
	}

	printVerbose("Found %d issues:\n\n", len(issues))
	logger.LogCommandComplete(ctx, "content_migrations.issues", len(issues))
	return formatOutput(issues, nil)
}
