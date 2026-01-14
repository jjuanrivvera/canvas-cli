package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var (
	migrationCourseID     int64
	migrationType         string
	migrationSourceCourse int64
	migrationFile         string
	migrationFileURL      string
	migrationFolderID     int64
	migrationSelective    bool
	migrationCopyOptions  string
	migrationDateShift    string
	migrationContentType  string
)

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

var migrationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List content migrations",
	Long: `List content migrations for a course.

Examples:
  canvas content-migrations list --course-id 1`,
	RunE: runMigrationsList,
}

var migrationsGetCmd = &cobra.Command{
	Use:   "get <migration-id>",
	Short: "Get a content migration",
	Long: `Get details of a specific content migration.

Examples:
  canvas content-migrations get 123 --course-id 1`,
	Args: ExactArgsWithUsage(1, "migration-id"),
	RunE: runMigrationsGet,
}

var migrationsCreateCmd = &cobra.Command{
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
	RunE: runMigrationsCreate,
}

var migrationsMigratorsCmd = &cobra.Command{
	Use:   "migrators",
	Short: "List available migration types",
	Long: `List available migrator types for a course.

This shows what types of content migrations can be performed.

Examples:
  canvas content-migrations migrators --course-id 1`,
	RunE: runMigrationsMigrators,
}

var migrationsContentCmd = &cobra.Command{
	Use:   "content <migration-id>",
	Short: "List migration content for selective import",
	Long: `List available content for selective import.

This shows what content is available to import from a migration.

Examples:
  canvas content-migrations content 123 --course-id 1
  canvas content-migrations content 123 --course-id 1 --type assignments`,
	Args: ExactArgsWithUsage(1, "migration-id"),
	RunE: runMigrationsContent,
}

var migrationsIssuesCmd = &cobra.Command{
	Use:   "issues <migration-id>",
	Short: "List migration issues",
	Long: `List issues encountered during a content migration.

Examples:
  canvas content-migrations issues 123 --course-id 1`,
	Args: ExactArgsWithUsage(1, "migration-id"),
	RunE: runMigrationsIssues,
}

func init() {
	rootCmd.AddCommand(contentMigrationsCmd)
	contentMigrationsCmd.AddCommand(migrationsListCmd)
	contentMigrationsCmd.AddCommand(migrationsGetCmd)
	contentMigrationsCmd.AddCommand(migrationsCreateCmd)
	contentMigrationsCmd.AddCommand(migrationsMigratorsCmd)
	contentMigrationsCmd.AddCommand(migrationsContentCmd)
	contentMigrationsCmd.AddCommand(migrationsIssuesCmd)

	// List flags
	migrationsListCmd.Flags().Int64Var(&migrationCourseID, "course-id", 0, "Course ID (required)")
	migrationsListCmd.MarkFlagRequired("course-id")

	// Get flags
	migrationsGetCmd.Flags().Int64Var(&migrationCourseID, "course-id", 0, "Course ID (required)")
	migrationsGetCmd.MarkFlagRequired("course-id")

	// Create flags
	migrationsCreateCmd.Flags().Int64Var(&migrationCourseID, "course-id", 0, "Course ID (required)")
	migrationsCreateCmd.Flags().StringVar(&migrationType, "type", "", "Migration type (required)")
	migrationsCreateCmd.Flags().Int64Var(&migrationSourceCourse, "source-course-id", 0, "Source course ID (for course_copy_importer)")
	migrationsCreateCmd.Flags().StringVar(&migrationFile, "file", "", "Export file to import")
	migrationsCreateCmd.Flags().StringVar(&migrationFileURL, "file-url", "", "URL to export file")
	migrationsCreateCmd.Flags().Int64Var(&migrationFolderID, "folder-id", 0, "Target folder ID")
	migrationsCreateCmd.Flags().BoolVar(&migrationSelective, "selective", false, "Enable selective import")
	migrationsCreateCmd.Flags().StringVar(&migrationCopyOptions, "copy-options", "", "JSON with copy options")
	migrationsCreateCmd.Flags().StringVar(&migrationDateShift, "date-shift", "", "JSON with date shift options")
	migrationsCreateCmd.MarkFlagRequired("course-id")
	migrationsCreateCmd.MarkFlagRequired("type")

	// Migrators flags
	migrationsMigratorsCmd.Flags().Int64Var(&migrationCourseID, "course-id", 0, "Course ID (required)")
	migrationsMigratorsCmd.MarkFlagRequired("course-id")

	// Content flags
	migrationsContentCmd.Flags().Int64Var(&migrationCourseID, "course-id", 0, "Course ID (required)")
	migrationsContentCmd.Flags().StringVar(&migrationContentType, "type", "", "Content type filter")
	migrationsContentCmd.MarkFlagRequired("course-id")

	// Issues flags
	migrationsIssuesCmd.Flags().Int64Var(&migrationCourseID, "course-id", 0, "Course ID (required)")
	migrationsIssuesCmd.MarkFlagRequired("course-id")
}

func runMigrationsList(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewContentMigrationsService(client)

	ctx := context.Background()
	migrations, err := service.List(ctx, migrationCourseID, nil)
	if err != nil {
		return fmt.Errorf("failed to list content migrations: %w", err)
	}

	if len(migrations) == 0 {
		fmt.Println("No content migrations found")
		return nil
	}

	printVerbose("Found %d content migrations:\n\n", len(migrations))
	return formatOutput(migrations, nil)
}

func runMigrationsGet(cmd *cobra.Command, args []string) error {
	migrationID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid migration ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewContentMigrationsService(client)

	ctx := context.Background()
	migration, err := service.Get(ctx, migrationCourseID, migrationID)
	if err != nil {
		return fmt.Errorf("failed to get content migration: %w", err)
	}

	return formatOutput(migration, nil)
}

func runMigrationsCreate(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	params := &api.CreateContentMigrationParams{
		MigrationType: migrationType,
	}

	if cmd.Flags().Changed("source-course-id") {
		params.SourceCourseID = &migrationSourceCourse
	}

	if migrationFile != "" {
		params.FilePath = migrationFile
	}

	if migrationFileURL != "" {
		params.FileURL = migrationFileURL
	}

	if cmd.Flags().Changed("folder-id") {
		params.FolderID = &migrationFolderID
	}

	if cmd.Flags().Changed("selective") {
		params.SelectiveImport = &migrationSelective
	}

	if migrationCopyOptions != "" {
		var copyOpts map[string]interface{}
		if err := json.Unmarshal([]byte(migrationCopyOptions), &copyOpts); err != nil {
			return fmt.Errorf("invalid copy options JSON: %w", err)
		}
		params.CopyOptions = copyOpts
	}

	if migrationDateShift != "" {
		var dateShift api.DateShiftOptions
		if err := json.Unmarshal([]byte(migrationDateShift), &dateShift); err != nil {
			return fmt.Errorf("invalid date shift JSON: %w", err)
		}
		params.DateShiftOptions = &dateShift
	}

	service := api.NewContentMigrationsService(client)

	ctx := context.Background()
	migration, err := service.Create(ctx, migrationCourseID, params)
	if err != nil {
		return fmt.Errorf("failed to create content migration: %w", err)
	}

	fmt.Printf("Content migration created (ID: %d)\n", migration.ID)
	fmt.Printf("Type: %s\n", migration.MigrationType)
	fmt.Printf("State: %s\n", migration.WorkflowState)
	return nil
}

func runMigrationsMigrators(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewContentMigrationsService(client)

	ctx := context.Background()
	migrators, err := service.ListMigrators(ctx, migrationCourseID)
	if err != nil {
		return fmt.Errorf("failed to list migrators: %w", err)
	}

	if len(migrators) == 0 {
		fmt.Println("No migrators available")
		return nil
	}

	printVerbose("Available migrators:\n\n")
	return formatOutput(migrators, nil)
}

func runMigrationsContent(cmd *cobra.Command, args []string) error {
	migrationID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid migration ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewContentMigrationsService(client)

	ctx := context.Background()
	items, err := service.ListContentList(ctx, migrationCourseID, migrationID, migrationContentType)
	if err != nil {
		return fmt.Errorf("failed to list content: %w", err)
	}

	if len(items) == 0 {
		fmt.Println("No content available for import")
		return nil
	}

	printVerbose("Available content:\n\n")
	return formatOutput(items, nil)
}

func runMigrationsIssues(cmd *cobra.Command, args []string) error {
	migrationID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid migration ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewContentMigrationsService(client)

	ctx := context.Background()
	issues, err := service.ListMigrationIssues(ctx, migrationCourseID, migrationID)
	if err != nil {
		return fmt.Errorf("failed to list migration issues: %w", err)
	}

	if len(issues) == 0 {
		fmt.Println("No issues found for this migration")
		return nil
	}

	printVerbose("Found %d issues:\n\n", len(issues))
	return formatOutput(issues, nil)
}
