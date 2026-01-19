package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
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

var blueprintAssociationsCmd = &cobra.Command{
	Use:   "associations",
	Short: "Manage associated courses",
	Long:  `Manage courses associated with a blueprint.`,
}

var blueprintMigrationsCmd = &cobra.Command{
	Use:   "migrations",
	Short: "Manage blueprint migrations",
	Long:  `Manage blueprint sync migrations.`,
}

func init() {
	rootCmd.AddCommand(blueprintCmd)
	blueprintCmd.AddCommand(newBlueprintGetCmd())
	blueprintCmd.AddCommand(blueprintAssociationsCmd)
	blueprintCmd.AddCommand(newBlueprintSyncCmd())
	blueprintCmd.AddCommand(newBlueprintChangesCmd())
	blueprintCmd.AddCommand(blueprintMigrationsCmd)

	blueprintAssociationsCmd.AddCommand(newBlueprintAssociationsListCmd())
	blueprintAssociationsCmd.AddCommand(newBlueprintAssociationsAddCmd())
	blueprintAssociationsCmd.AddCommand(newBlueprintAssociationsRemoveCmd())

	blueprintMigrationsCmd.AddCommand(newBlueprintMigrationsListCmd())
	blueprintMigrationsCmd.AddCommand(newBlueprintMigrationsGetCmd())
}

func newBlueprintGetCmd() *cobra.Command {
	opts := &options.BlueprintGetOptions{
		TemplateID: "default",
	}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get blueprint details",
		Long: `Get details of a blueprint course template.

Examples:
  canvas blueprint get --course-id 1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runBlueprintGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.TemplateID, "template-id", "default", "Blueprint template ID")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newBlueprintAssociationsListCmd() *cobra.Command {
	opts := &options.BlueprintAssociationsListOptions{
		TemplateID: "default",
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List associated courses",
		Long: `List courses associated with a blueprint.

Examples:
  canvas blueprint associations list --course-id 1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runBlueprintAssociationsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.TemplateID, "template-id", "default", "Blueprint template ID")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newBlueprintAssociationsAddCmd() *cobra.Command {
	opts := &options.BlueprintAssociationsAddOptions{
		TemplateID: "default",
	}

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add courses to blueprint",
		Long: `Add courses to a blueprint's associations.

Examples:
  canvas blueprint associations add --course-id 1 --course-ids-to-add 100,101,102`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runBlueprintAssociationsAdd(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.TemplateID, "template-id", "default", "Blueprint template ID")
	cmd.Flags().StringVar(&opts.CourseIDsStr, "course-ids-to-add", "", "Comma-separated course IDs to add")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("course-ids-to-add")

	return cmd
}

func newBlueprintAssociationsRemoveCmd() *cobra.Command {
	opts := &options.BlueprintAssociationsRemoveOptions{
		TemplateID: "default",
	}

	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove courses from blueprint",
		Long: `Remove courses from a blueprint's associations.

Examples:
  canvas blueprint associations remove --course-id 1 --course-ids-to-remove 100,101`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runBlueprintAssociationsRemove(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.TemplateID, "template-id", "default", "Blueprint template ID")
	cmd.Flags().StringVar(&opts.CourseIDsStr, "course-ids-to-remove", "", "Comma-separated course IDs to remove")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("course-ids-to-remove")

	return cmd
}

func newBlueprintSyncCmd() *cobra.Command {
	opts := &options.BlueprintSyncOptions{
		TemplateID: "default",
	}

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync blueprint to associated courses",
		Long: `Begin a sync of the blueprint to all associated courses.

Examples:
  canvas blueprint sync --course-id 1
  canvas blueprint sync --course-id 1 --comment "Weekly content update"
  canvas blueprint sync --course-id 1 --send-notification --copy-settings`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.NotifySet = cmd.Flags().Changed("send-notification")
			opts.CopySettingsSet = cmd.Flags().Changed("copy-settings")
			opts.PublishSet = cmd.Flags().Changed("publish")

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runBlueprintSync(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.TemplateID, "template-id", "default", "Blueprint template ID")
	cmd.Flags().StringVar(&opts.Comment, "comment", "", "Sync comment")
	cmd.Flags().BoolVar(&opts.Notify, "send-notification", false, "Send notification to users")
	cmd.Flags().BoolVar(&opts.CopySettings, "copy-settings", false, "Copy course settings")
	cmd.Flags().BoolVar(&opts.Publish, "publish", false, "Publish synced content")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newBlueprintChangesCmd() *cobra.Command {
	opts := &options.BlueprintChangesOptions{
		TemplateID: "default",
	}

	cmd := &cobra.Command{
		Use:   "changes",
		Short: "Show unsynced changes",
		Long: `Show changes that have not been synced to associated courses.

Examples:
  canvas blueprint changes --course-id 1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runBlueprintChanges(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.TemplateID, "template-id", "default", "Blueprint template ID")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newBlueprintMigrationsListCmd() *cobra.Command {
	opts := &options.BlueprintMigrationsListOptions{
		TemplateID: "default",
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List migrations",
		Long: `List blueprint sync migrations.

Examples:
  canvas blueprint migrations list --course-id 1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runBlueprintMigrationsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.TemplateID, "template-id", "default", "Blueprint template ID")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newBlueprintMigrationsGetCmd() *cobra.Command {
	opts := &options.BlueprintMigrationsGetOptions{
		TemplateID: "default",
	}

	cmd := &cobra.Command{
		Use:   "get <migration-id>",
		Short: "Get migration details",
		Long: `Get details of a specific blueprint migration.

Examples:
  canvas blueprint migrations get 123 --course-id 1
  canvas blueprint migrations get 123 --course-id 1 --include user`,
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

			return runBlueprintMigrationsGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.TemplateID, "template-id", "default", "Blueprint template ID")
	cmd.Flags().StringSliceVar(&opts.Include, "include", nil, "Include options (user)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func runBlueprintGet(ctx context.Context, client *api.Client, opts *options.BlueprintGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "blueprint.get", map[string]interface{}{
		"course_id":   opts.CourseID,
		"template_id": opts.TemplateID,
	})

	service := api.NewBlueprintService(client)

	template, err := service.GetTemplate(ctx, opts.CourseID, opts.TemplateID)
	if err != nil {
		logger.LogCommandError(ctx, "blueprint.get", err, map[string]interface{}{
			"course_id":   opts.CourseID,
			"template_id": opts.TemplateID,
		})
		return fmt.Errorf("failed to get blueprint: %w", err)
	}

	logger.LogCommandComplete(ctx, "blueprint.get", 1)
	return formatOutput(template, nil)
}

func runBlueprintAssociationsList(ctx context.Context, client *api.Client, opts *options.BlueprintAssociationsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "blueprint.associations.list", map[string]interface{}{
		"course_id":   opts.CourseID,
		"template_id": opts.TemplateID,
	})

	service := api.NewBlueprintService(client)

	courses, err := service.ListAssociatedCourses(ctx, opts.CourseID, opts.TemplateID, nil)
	if err != nil {
		logger.LogCommandError(ctx, "blueprint.associations.list", err, map[string]interface{}{
			"course_id":   opts.CourseID,
			"template_id": opts.TemplateID,
		})
		return fmt.Errorf("failed to list associated courses: %w", err)
	}

	if len(courses) == 0 {
		fmt.Println("No associated courses found")
		logger.LogCommandComplete(ctx, "blueprint.associations.list", 0)
		return nil
	}

	printVerbose("Found %d associated courses:\n\n", len(courses))
	logger.LogCommandComplete(ctx, "blueprint.associations.list", len(courses))
	return formatOutput(courses, nil)
}

func runBlueprintAssociationsAdd(ctx context.Context, client *api.Client, opts *options.BlueprintAssociationsAddOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "blueprint.associations.add", map[string]interface{}{
		"course_id":   opts.CourseID,
		"template_id": opts.TemplateID,
	})

	courseIDs, err := parseIDList(opts.CourseIDsStr)
	if err != nil {
		logger.LogCommandError(ctx, "blueprint.associations.add", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("invalid course IDs: %w", err)
	}

	params := &api.UpdateAssociationsParams{
		CourseIDsToAdd: courseIDs,
	}

	service := api.NewBlueprintService(client)

	err = service.UpdateAssociations(ctx, opts.CourseID, opts.TemplateID, params)
	if err != nil {
		logger.LogCommandError(ctx, "blueprint.associations.add", err, map[string]interface{}{
			"course_id":   opts.CourseID,
			"template_id": opts.TemplateID,
		})
		return fmt.Errorf("failed to add associations: %w", err)
	}

	fmt.Printf("Added %d courses to blueprint associations\n", len(courseIDs))
	logger.LogCommandComplete(ctx, "blueprint.associations.add", len(courseIDs))
	return nil
}

func runBlueprintAssociationsRemove(ctx context.Context, client *api.Client, opts *options.BlueprintAssociationsRemoveOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "blueprint.associations.remove", map[string]interface{}{
		"course_id":   opts.CourseID,
		"template_id": opts.TemplateID,
	})

	courseIDs, err := parseIDList(opts.CourseIDsStr)
	if err != nil {
		logger.LogCommandError(ctx, "blueprint.associations.remove", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("invalid course IDs: %w", err)
	}

	params := &api.UpdateAssociationsParams{
		CourseIDsToRemove: courseIDs,
	}

	service := api.NewBlueprintService(client)

	err = service.UpdateAssociations(ctx, opts.CourseID, opts.TemplateID, params)
	if err != nil {
		logger.LogCommandError(ctx, "blueprint.associations.remove", err, map[string]interface{}{
			"course_id":   opts.CourseID,
			"template_id": opts.TemplateID,
		})
		return fmt.Errorf("failed to remove associations: %w", err)
	}

	fmt.Printf("Removed %d courses from blueprint associations\n", len(courseIDs))
	logger.LogCommandComplete(ctx, "blueprint.associations.remove", len(courseIDs))
	return nil
}

func runBlueprintSync(ctx context.Context, client *api.Client, opts *options.BlueprintSyncOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "blueprint.sync", map[string]interface{}{
		"course_id":   opts.CourseID,
		"template_id": opts.TemplateID,
		"comment":     opts.Comment,
	})

	params := &api.SyncParams{
		Comment: opts.Comment,
	}

	if opts.NotifySet {
		params.SendNotification = &opts.Notify
	}

	if opts.CopySettingsSet {
		params.CopySettings = &opts.CopySettings
	}

	if opts.PublishSet {
		params.PublishAfterSync = &opts.Publish
	}

	service := api.NewBlueprintService(client)

	migration, err := service.BeginSync(ctx, opts.CourseID, opts.TemplateID, params)
	if err != nil {
		logger.LogCommandError(ctx, "blueprint.sync", err, map[string]interface{}{
			"course_id":   opts.CourseID,
			"template_id": opts.TemplateID,
		})
		return fmt.Errorf("failed to begin sync: %w", err)
	}

	fmt.Printf("Blueprint sync started (Migration ID: %d)\n", migration.ID)
	fmt.Printf("State: %s\n", migration.WorkflowState)
	logger.LogCommandComplete(ctx, "blueprint.sync", 1)
	return nil
}

func runBlueprintChanges(ctx context.Context, client *api.Client, opts *options.BlueprintChangesOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "blueprint.changes", map[string]interface{}{
		"course_id":   opts.CourseID,
		"template_id": opts.TemplateID,
	})

	service := api.NewBlueprintService(client)

	changes, err := service.ListUnsyncedChanges(ctx, opts.CourseID, opts.TemplateID)
	if err != nil {
		logger.LogCommandError(ctx, "blueprint.changes", err, map[string]interface{}{
			"course_id":   opts.CourseID,
			"template_id": opts.TemplateID,
		})
		return fmt.Errorf("failed to list changes: %w", err)
	}

	if len(changes) == 0 {
		fmt.Println("No unsynced changes")
		logger.LogCommandComplete(ctx, "blueprint.changes", 0)
		return nil
	}

	printVerbose("Found %d unsynced changes:\n\n", len(changes))
	logger.LogCommandComplete(ctx, "blueprint.changes", len(changes))
	return formatOutput(changes, nil)
}

func runBlueprintMigrationsList(ctx context.Context, client *api.Client, opts *options.BlueprintMigrationsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "blueprint.migrations.list", map[string]interface{}{
		"course_id":   opts.CourseID,
		"template_id": opts.TemplateID,
	})

	service := api.NewBlueprintService(client)

	migrations, err := service.ListMigrations(ctx, opts.CourseID, opts.TemplateID, nil)
	if err != nil {
		logger.LogCommandError(ctx, "blueprint.migrations.list", err, map[string]interface{}{
			"course_id":   opts.CourseID,
			"template_id": opts.TemplateID,
		})
		return fmt.Errorf("failed to list migrations: %w", err)
	}

	if len(migrations) == 0 {
		fmt.Println("No migrations found")
		logger.LogCommandComplete(ctx, "blueprint.migrations.list", 0)
		return nil
	}

	printVerbose("Found %d migrations:\n\n", len(migrations))
	logger.LogCommandComplete(ctx, "blueprint.migrations.list", len(migrations))
	return formatOutput(migrations, nil)
}

func runBlueprintMigrationsGet(ctx context.Context, client *api.Client, opts *options.BlueprintMigrationsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "blueprint.migrations.get", map[string]interface{}{
		"course_id":    opts.CourseID,
		"template_id":  opts.TemplateID,
		"migration_id": opts.MigrationID,
	})

	service := api.NewBlueprintService(client)

	migration, err := service.GetMigration(ctx, opts.CourseID, opts.TemplateID, opts.MigrationID, opts.Include)
	if err != nil {
		logger.LogCommandError(ctx, "blueprint.migrations.get", err, map[string]interface{}{
			"course_id":    opts.CourseID,
			"template_id":  opts.TemplateID,
			"migration_id": opts.MigrationID,
		})
		return fmt.Errorf("failed to get migration: %w", err)
	}

	logger.LogCommandComplete(ctx, "blueprint.migrations.get", 1)
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
