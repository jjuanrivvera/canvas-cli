package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var (
	blueprintCourseID     int64
	blueprintTemplateID   string
	blueprintCourseIDs    string
	blueprintComment      string
	blueprintNotify       bool
	blueprintCopySettings bool
	blueprintPublish      bool
	blueprintInclude      []string
)

var blueprintCmd = &cobra.Command{
	Use:   "blueprint",
	Short: "Manage blueprint courses",
	Long: `Manage Canvas blueprint courses.

Blueprint courses allow you to create a master course that can be
synced to associated courses, maintaining consistent content.

Examples:
  canvas blueprint get --course-id 1
  canvas blueprint associations list --course-id 1
  canvas blueprint sync --course-id 1 --comment "Weekly update"`,
}

var blueprintGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get blueprint details",
	Long: `Get details of a blueprint course template.

Examples:
  canvas blueprint get --course-id 1`,
	RunE: runBlueprintGet,
}

var blueprintAssociationsCmd = &cobra.Command{
	Use:   "associations",
	Short: "Manage associated courses",
	Long:  `Manage courses associated with a blueprint.`,
}

var blueprintAssociationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List associated courses",
	Long: `List courses associated with a blueprint.

Examples:
  canvas blueprint associations list --course-id 1`,
	RunE: runBlueprintAssociationsList,
}

var blueprintAssociationsAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add courses to blueprint",
	Long: `Add courses to a blueprint's associations.

Examples:
  canvas blueprint associations add --course-id 1 --course-ids-to-add 100,101,102`,
	RunE: runBlueprintAssociationsAdd,
}

var blueprintAssociationsRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove courses from blueprint",
	Long: `Remove courses from a blueprint's associations.

Examples:
  canvas blueprint associations remove --course-id 1 --course-ids-to-remove 100,101`,
	RunE: runBlueprintAssociationsRemove,
}

var blueprintSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync blueprint to associated courses",
	Long: `Begin a sync of the blueprint to all associated courses.

Examples:
  canvas blueprint sync --course-id 1
  canvas blueprint sync --course-id 1 --comment "Weekly content update"
  canvas blueprint sync --course-id 1 --send-notification --copy-settings`,
	RunE: runBlueprintSync,
}

var blueprintChangesCmd = &cobra.Command{
	Use:   "changes",
	Short: "Show unsynced changes",
	Long: `Show changes that have not been synced to associated courses.

Examples:
  canvas blueprint changes --course-id 1`,
	RunE: runBlueprintChanges,
}

var blueprintMigrationsCmd = &cobra.Command{
	Use:   "migrations",
	Short: "Manage blueprint migrations",
	Long:  `Manage blueprint sync migrations.`,
}

var blueprintMigrationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List migrations",
	Long: `List blueprint sync migrations.

Examples:
  canvas blueprint migrations list --course-id 1`,
	RunE: runBlueprintMigrationsList,
}

var blueprintMigrationsGetCmd = &cobra.Command{
	Use:   "get <migration-id>",
	Short: "Get migration details",
	Long: `Get details of a specific blueprint migration.

Examples:
  canvas blueprint migrations get 123 --course-id 1
  canvas blueprint migrations get 123 --course-id 1 --include user`,
	Args: ExactArgsWithUsage(1, "migration-id"),
	RunE: runBlueprintMigrationsGet,
}

func init() {
	rootCmd.AddCommand(blueprintCmd)
	blueprintCmd.AddCommand(blueprintGetCmd)
	blueprintCmd.AddCommand(blueprintAssociationsCmd)
	blueprintCmd.AddCommand(blueprintSyncCmd)
	blueprintCmd.AddCommand(blueprintChangesCmd)
	blueprintCmd.AddCommand(blueprintMigrationsCmd)

	blueprintAssociationsCmd.AddCommand(blueprintAssociationsListCmd)
	blueprintAssociationsCmd.AddCommand(blueprintAssociationsAddCmd)
	blueprintAssociationsCmd.AddCommand(blueprintAssociationsRemoveCmd)

	blueprintMigrationsCmd.AddCommand(blueprintMigrationsListCmd)
	blueprintMigrationsCmd.AddCommand(blueprintMigrationsGetCmd)

	// Get flags
	blueprintGetCmd.Flags().Int64Var(&blueprintCourseID, "course-id", 0, "Course ID (required)")
	blueprintGetCmd.Flags().StringVar(&blueprintTemplateID, "template-id", "default", "Blueprint template ID")
	blueprintGetCmd.MarkFlagRequired("course-id")

	// Associations list flags
	blueprintAssociationsListCmd.Flags().Int64Var(&blueprintCourseID, "course-id", 0, "Course ID (required)")
	blueprintAssociationsListCmd.Flags().StringVar(&blueprintTemplateID, "template-id", "default", "Blueprint template ID")
	blueprintAssociationsListCmd.MarkFlagRequired("course-id")

	// Associations add flags
	blueprintAssociationsAddCmd.Flags().Int64Var(&blueprintCourseID, "course-id", 0, "Course ID (required)")
	blueprintAssociationsAddCmd.Flags().StringVar(&blueprintTemplateID, "template-id", "default", "Blueprint template ID")
	blueprintAssociationsAddCmd.Flags().StringVar(&blueprintCourseIDs, "course-ids-to-add", "", "Comma-separated course IDs to add")
	blueprintAssociationsAddCmd.MarkFlagRequired("course-id")
	blueprintAssociationsAddCmd.MarkFlagRequired("course-ids-to-add")

	// Associations remove flags
	blueprintAssociationsRemoveCmd.Flags().Int64Var(&blueprintCourseID, "course-id", 0, "Course ID (required)")
	blueprintAssociationsRemoveCmd.Flags().StringVar(&blueprintTemplateID, "template-id", "default", "Blueprint template ID")
	blueprintAssociationsRemoveCmd.Flags().StringVar(&blueprintCourseIDs, "course-ids-to-remove", "", "Comma-separated course IDs to remove")
	blueprintAssociationsRemoveCmd.MarkFlagRequired("course-id")
	blueprintAssociationsRemoveCmd.MarkFlagRequired("course-ids-to-remove")

	// Sync flags
	blueprintSyncCmd.Flags().Int64Var(&blueprintCourseID, "course-id", 0, "Course ID (required)")
	blueprintSyncCmd.Flags().StringVar(&blueprintTemplateID, "template-id", "default", "Blueprint template ID")
	blueprintSyncCmd.Flags().StringVar(&blueprintComment, "comment", "", "Sync comment")
	blueprintSyncCmd.Flags().BoolVar(&blueprintNotify, "send-notification", false, "Send notification to users")
	blueprintSyncCmd.Flags().BoolVar(&blueprintCopySettings, "copy-settings", false, "Copy course settings")
	blueprintSyncCmd.Flags().BoolVar(&blueprintPublish, "publish", false, "Publish synced content")
	blueprintSyncCmd.MarkFlagRequired("course-id")

	// Changes flags
	blueprintChangesCmd.Flags().Int64Var(&blueprintCourseID, "course-id", 0, "Course ID (required)")
	blueprintChangesCmd.Flags().StringVar(&blueprintTemplateID, "template-id", "default", "Blueprint template ID")
	blueprintChangesCmd.MarkFlagRequired("course-id")

	// Migrations list flags
	blueprintMigrationsListCmd.Flags().Int64Var(&blueprintCourseID, "course-id", 0, "Course ID (required)")
	blueprintMigrationsListCmd.Flags().StringVar(&blueprintTemplateID, "template-id", "default", "Blueprint template ID")
	blueprintMigrationsListCmd.MarkFlagRequired("course-id")

	// Migrations get flags
	blueprintMigrationsGetCmd.Flags().Int64Var(&blueprintCourseID, "course-id", 0, "Course ID (required)")
	blueprintMigrationsGetCmd.Flags().StringVar(&blueprintTemplateID, "template-id", "default", "Blueprint template ID")
	blueprintMigrationsGetCmd.Flags().StringSliceVar(&blueprintInclude, "include", nil, "Include options (user)")
	blueprintMigrationsGetCmd.MarkFlagRequired("course-id")
}

func runBlueprintGet(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewBlueprintService(client)

	ctx := context.Background()
	template, err := service.GetTemplate(ctx, blueprintCourseID, blueprintTemplateID)
	if err != nil {
		return fmt.Errorf("failed to get blueprint: %w", err)
	}

	return formatOutput(template, nil)
}

func runBlueprintAssociationsList(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewBlueprintService(client)

	ctx := context.Background()
	courses, err := service.ListAssociatedCourses(ctx, blueprintCourseID, blueprintTemplateID, nil)
	if err != nil {
		return fmt.Errorf("failed to list associated courses: %w", err)
	}

	if len(courses) == 0 {
		fmt.Println("No associated courses found")
		return nil
	}

	printVerbose("Found %d associated courses:\n\n", len(courses))
	return formatOutput(courses, nil)
}

func runBlueprintAssociationsAdd(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	courseIDs, err := parseIDList(blueprintCourseIDs)
	if err != nil {
		return fmt.Errorf("invalid course IDs: %w", err)
	}

	params := &api.UpdateAssociationsParams{
		CourseIDsToAdd: courseIDs,
	}

	service := api.NewBlueprintService(client)

	ctx := context.Background()
	err = service.UpdateAssociations(ctx, blueprintCourseID, blueprintTemplateID, params)
	if err != nil {
		return fmt.Errorf("failed to add associations: %w", err)
	}

	fmt.Printf("Added %d courses to blueprint associations\n", len(courseIDs))
	return nil
}

func runBlueprintAssociationsRemove(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	courseIDs, err := parseIDList(blueprintCourseIDs)
	if err != nil {
		return fmt.Errorf("invalid course IDs: %w", err)
	}

	params := &api.UpdateAssociationsParams{
		CourseIDsToRemove: courseIDs,
	}

	service := api.NewBlueprintService(client)

	ctx := context.Background()
	err = service.UpdateAssociations(ctx, blueprintCourseID, blueprintTemplateID, params)
	if err != nil {
		return fmt.Errorf("failed to remove associations: %w", err)
	}

	fmt.Printf("Removed %d courses from blueprint associations\n", len(courseIDs))
	return nil
}

func runBlueprintSync(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	params := &api.SyncParams{
		Comment: blueprintComment,
	}

	if cmd.Flags().Changed("send-notification") {
		params.SendNotification = &blueprintNotify
	}

	if cmd.Flags().Changed("copy-settings") {
		params.CopySettings = &blueprintCopySettings
	}

	if cmd.Flags().Changed("publish") {
		params.PublishAfterSync = &blueprintPublish
	}

	service := api.NewBlueprintService(client)

	ctx := context.Background()
	migration, err := service.BeginSync(ctx, blueprintCourseID, blueprintTemplateID, params)
	if err != nil {
		return fmt.Errorf("failed to begin sync: %w", err)
	}

	fmt.Printf("Blueprint sync started (Migration ID: %d)\n", migration.ID)
	fmt.Printf("State: %s\n", migration.WorkflowState)
	return nil
}

func runBlueprintChanges(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewBlueprintService(client)

	ctx := context.Background()
	changes, err := service.ListUnsyncedChanges(ctx, blueprintCourseID, blueprintTemplateID)
	if err != nil {
		return fmt.Errorf("failed to list changes: %w", err)
	}

	if len(changes) == 0 {
		fmt.Println("No unsynced changes")
		return nil
	}

	printVerbose("Found %d unsynced changes:\n\n", len(changes))
	return formatOutput(changes, nil)
}

func runBlueprintMigrationsList(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewBlueprintService(client)

	ctx := context.Background()
	migrations, err := service.ListMigrations(ctx, blueprintCourseID, blueprintTemplateID, nil)
	if err != nil {
		return fmt.Errorf("failed to list migrations: %w", err)
	}

	if len(migrations) == 0 {
		fmt.Println("No migrations found")
		return nil
	}

	printVerbose("Found %d migrations:\n\n", len(migrations))
	return formatOutput(migrations, nil)
}

func runBlueprintMigrationsGet(cmd *cobra.Command, args []string) error {
	migrationID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid migration ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewBlueprintService(client)

	ctx := context.Background()
	migration, err := service.GetMigration(ctx, blueprintCourseID, blueprintTemplateID, migrationID, blueprintInclude)
	if err != nil {
		return fmt.Errorf("failed to get migration: %w", err)
	}

	return formatOutput(migration, nil)
}

// parseIDList parses a comma-separated list of IDs
func parseIDList(s string) ([]int64, error) {
	if s == "" {
		return nil, nil
	}

	parts := strings.Split(s, ",")
	ids := make([]int64, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid ID '%s': %w", part, err)
		}
		ids = append(ids, id)
	}

	return ids, nil
}
