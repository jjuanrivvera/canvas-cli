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

// modulesCmd represents the modules command group
var modulesCmd = &cobra.Command{
	Use:   "modules",
	Short: "Manage Canvas course modules",
	Long: `Manage Canvas course modules including listing, viewing, creating, and updating modules.

Modules are collections of learning materials useful for organizing courses and optionally
providing a linear flow through them. Module items can be accessed linearly or sequentially
depending on module configuration.

Examples:
  canvas modules list --course-id 123
  canvas modules get --course-id 123 456
  canvas modules create --course-id 123 --name "Week 1"
  canvas modules items list --course-id 123 --module-id 456`,
}

// newModulesListCmd creates the modules list command
func newModulesListCmd() *cobra.Command {
	opts := &options.ModulesListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List modules in a course",
		Long: `List all modules in a Canvas course.

Examples:
  canvas modules list --course-id 123
  canvas modules list --course-id 123 --include items
  canvas modules list --course-id 123 --search "Week"
  canvas modules list --course-id 123 --student-id 789`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runModulesList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Additional data to include (items, content_details)")
	cmd.Flags().StringVar(&opts.SearchTerm, "search", "", "Search by module name")
	cmd.Flags().StringVar(&opts.StudentID, "student-id", "", "Get completion info for this student")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

// newModulesGetCmd creates the modules get command
func newModulesGetCmd() *cobra.Command {
	opts := &options.ModulesGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <module-id>",
		Short: "Get details of a specific module",
		Long: `Get details of a specific module by ID.

Examples:
  canvas modules get --course-id 123 456
  canvas modules get --course-id 123 456 --include items,content_details`,
		Args: ExactArgsWithUsage(1, "module-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			moduleID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid module ID: %s", args[0])
			}
			opts.ModuleID = moduleID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runModulesGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Additional data to include (items, content_details)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

// newModulesCreateCmd creates the modules create command
func newModulesCreateCmd() *cobra.Command {
	opts := &options.ModulesCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new module",
		Long: `Create a new module in a course.

Examples:
  canvas modules create --course-id 123 --name "Week 1"
  canvas modules create --course-id 123 --name "Week 2" --position 2
  canvas modules create --course-id 123 --name "Unit 2" --prerequisite-modules 1,2
  canvas modules create --course-id 123 --name "Final" --require-sequential-progress`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runModulesCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Module name (required)")
	cmd.Flags().StringVar(&opts.UnlockAt, "unlock-at", "", "Date to unlock the module (ISO 8601)")
	cmd.Flags().IntVar(&opts.Position, "position", 0, "Position in the course (1-based)")
	cmd.Flags().BoolVar(&opts.RequireSequentialProgress, "require-sequential-progress", false, "Require sequential progress")
	cmd.Flags().Int64SliceVar(&opts.PrerequisiteModuleIDs, "prerequisite-modules", []int64{}, "IDs of prerequisite modules")
	cmd.Flags().BoolVar(&opts.PublishFinalGrade, "publish-final-grade", false, "Publish final grade on completion")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("name")

	return cmd
}

// newModulesUpdateCmd creates the modules update command
func newModulesUpdateCmd() *cobra.Command {
	opts := &options.ModulesUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <module-id>",
		Short: "Update an existing module",
		Long: `Update an existing module.

Examples:
  canvas modules update --course-id 123 456 --name "Updated Name"
  canvas modules update --course-id 123 456 --published
  canvas modules update --course-id 123 456 --position 3`,
		Args: ExactArgsWithUsage(1, "module-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			moduleID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid module ID: %s", args[0])
			}
			opts.ModuleID = moduleID

			// Track which fields were set
			opts.NameSet = cmd.Flags().Changed("name")
			opts.UnlockAtSet = cmd.Flags().Changed("unlock-at")
			opts.PositionSet = cmd.Flags().Changed("position")
			opts.RequireSequentialProgressSet = cmd.Flags().Changed("require-sequential-progress")
			opts.PrerequisiteModuleIDsSet = cmd.Flags().Changed("prerequisite-modules")
			opts.PublishFinalGradeSet = cmd.Flags().Changed("publish-final-grade")
			opts.PublishedSet = cmd.Flags().Changed("published")

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runModulesUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.Name, "name", "", "New module name")
	cmd.Flags().StringVar(&opts.UnlockAt, "unlock-at", "", "Date to unlock the module (ISO 8601)")
	cmd.Flags().IntVar(&opts.Position, "position", 0, "New position in the course")
	cmd.Flags().BoolVar(&opts.RequireSequentialProgress, "require-sequential-progress", false, "Require sequential progress")
	cmd.Flags().Int64SliceVar(&opts.PrerequisiteModuleIDs, "prerequisite-modules", []int64{}, "IDs of prerequisite modules")
	cmd.Flags().BoolVar(&opts.PublishFinalGrade, "publish-final-grade", false, "Publish final grade on completion")
	cmd.Flags().BoolVar(&opts.Published, "published", false, "Publish the module")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

// newModulesDeleteCmd creates the modules delete command
func newModulesDeleteCmd() *cobra.Command {
	opts := &options.ModulesDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <module-id>",
		Short: "Delete a module",
		Long: `Delete a module from a course.

Examples:
  canvas modules delete --course-id 123 456`,
		Args: ExactArgsWithUsage(1, "module-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			moduleID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid module ID: %s", args[0])
			}
			opts.ModuleID = moduleID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runModulesDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Skip confirmation prompt")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

// newModulesRelockCmd creates the modules relock command
func newModulesRelockCmd() *cobra.Command {
	opts := &options.ModulesRelockOptions{}

	cmd := &cobra.Command{
		Use:   "relock <module-id>",
		Short: "Re-lock module progressions",
		Long: `Re-lock module progressions to their default locked state.

This recalculates progressions based on current requirements. Adding progression
requirements to an active course will not lock students out of modules they have
already unlocked unless this action is called.

Examples:
  canvas modules relock --course-id 123 456`,
		Args: ExactArgsWithUsage(1, "module-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			moduleID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid module ID: %s", args[0])
			}
			opts.ModuleID = moduleID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runModulesRelock(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

// newModulesPublishCmd creates the modules publish command
func newModulesPublishCmd() *cobra.Command {
	opts := &options.ModulesPublishOptions{}

	cmd := &cobra.Command{
		Use:   "publish <module-id>",
		Short: "Publish a module",
		Long: `Publish a module to make it visible to students.

This is a convenience command equivalent to:
  canvas modules update --course-id 123 456 --published

Examples:
  canvas modules publish --course-id 123 456`,
		Args: ExactArgsWithUsage(1, "module-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			moduleID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid module ID: %s", args[0])
			}
			opts.ModuleID = moduleID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runModulesPublish(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

// newModulesUnpublishCmd creates the modules unpublish command
func newModulesUnpublishCmd() *cobra.Command {
	opts := &options.ModulesUnpublishOptions{}

	cmd := &cobra.Command{
		Use:   "unpublish <module-id>",
		Short: "Unpublish a module",
		Long: `Unpublish a module to hide it from students.

This is a convenience command that sets the module's published state to false.

Examples:
  canvas modules unpublish --course-id 123 456`,
		Args: ExactArgsWithUsage(1, "module-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			moduleID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid module ID: %s", args[0])
			}
			opts.ModuleID = moduleID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runModulesUnpublish(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

// modulesItemsCmd represents the module items subcommand group
var modulesItemsCmd = &cobra.Command{
	Use:   "items",
	Short: "Manage module items",
	Long: `Manage items within a module.

Module items can be of various types: File, Page, Discussion, Assignment,
Quiz, SubHeader, ExternalUrl, or ExternalTool.

Examples:
  canvas modules items list --course-id 123 --module-id 456
  canvas modules items get --course-id 123 --module-id 456 789
  canvas modules items create --course-id 123 --module-id 456 --type Assignment --content-id 999`,
}

// newModulesItemsListCmd creates the module items list command
func newModulesItemsListCmd() *cobra.Command {
	opts := &options.ModulesItemsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List items in a module",
		Long: `List all items in a module.

Examples:
  canvas modules items list --course-id 123 --module-id 456
  canvas modules items list --course-id 123 --module-id 456 --include content_details
  canvas modules items list --course-id 123 --module-id 456 --search "Quiz"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runModulesItemsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.ModuleID, "module-id", 0, "Module ID (required)")
	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Additional data to include (content_details)")
	cmd.Flags().StringVar(&opts.SearchTerm, "search", "", "Search by item title")
	cmd.Flags().StringVar(&opts.StudentID, "student-id", "", "Get completion info for this student")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("module-id")

	return cmd
}

// newModulesItemsGetCmd creates the module items get command
func newModulesItemsGetCmd() *cobra.Command {
	opts := &options.ModulesItemsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <item-id>",
		Short: "Get details of a module item",
		Long: `Get details of a specific module item.

Examples:
  canvas modules items get --course-id 123 --module-id 456 789
  canvas modules items get --course-id 123 --module-id 456 789 --include content_details`,
		Args: ExactArgsWithUsage(1, "item-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			itemID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid item ID: %s", args[0])
			}
			opts.ItemID = itemID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runModulesItemsGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.ModuleID, "module-id", 0, "Module ID (required)")
	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Additional data to include (content_details)")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("module-id")

	return cmd
}

// newModulesItemsCreateCmd creates the module items create command
func newModulesItemsCreateCmd() *cobra.Command {
	opts := &options.ModulesItemsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new module item",
		Long: `Create a new item in a module.

Item types:
  - File: Requires --content-id
  - Page: Requires --page-url
  - Discussion: Requires --content-id
  - Assignment: Requires --content-id
  - Quiz: Requires --content-id
  - SubHeader: Only requires --title
  - ExternalUrl: Requires --external-url
  - ExternalTool: Requires --content-id or --external-url

Completion requirement types:
  - must_view: Applies to all item types
  - must_contribute: Discussion, Assignment, Page
  - must_submit: Assignment, Quiz
  - must_mark_done: Assignment, Page
  - min_score: Assignment, Quiz (requires --min-score)

Examples:
  canvas modules items create --course-id 123 --module-id 456 --type Assignment --content-id 999
  canvas modules items create --course-id 123 --module-id 456 --type Page --page-url "intro-page"
  canvas modules items create --course-id 123 --module-id 456 --type SubHeader --title "Unit 1"
  canvas modules items create --course-id 123 --module-id 456 --type ExternalUrl --external-url "https://example.com" --title "Resource"
  canvas modules items create --course-id 123 --module-id 456 --type Assignment --content-id 999 --completion-type min_score --min-score 80`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runModulesItemsCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.ModuleID, "module-id", 0, "Module ID (required)")
	cmd.Flags().StringVar(&opts.Type, "type", "", "Item type: File, Page, Discussion, Assignment, Quiz, SubHeader, ExternalUrl, ExternalTool (required)")
	cmd.Flags().StringVar(&opts.Title, "title", "", "Item title")
	cmd.Flags().Int64Var(&opts.ContentID, "content-id", 0, "Content ID (for File, Discussion, Assignment, Quiz, ExternalTool)")
	cmd.Flags().StringVar(&opts.PageURL, "page-url", "", "Page URL slug (for Page type)")
	cmd.Flags().StringVar(&opts.ExternalURL, "external-url", "", "External URL (for ExternalUrl, ExternalTool)")
	cmd.Flags().BoolVar(&opts.NewTab, "new-tab", false, "Open in new tab (for ExternalTool)")
	cmd.Flags().IntVar(&opts.Indent, "indent", 0, "Indent level (0-based)")
	cmd.Flags().StringVar(&opts.CompletionType, "completion-type", "", "Completion requirement: must_view, must_contribute, must_submit, must_mark_done, min_score")
	cmd.Flags().Float64Var(&opts.MinScore, "min-score", 0, "Minimum score for min_score completion type")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("module-id")
	cmd.MarkFlagRequired("type")

	return cmd
}

// newModulesItemsUpdateCmd creates the module items update command
func newModulesItemsUpdateCmd() *cobra.Command {
	opts := &options.ModulesItemsUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <item-id>",
		Short: "Update a module item",
		Long: `Update an existing module item.

Examples:
  canvas modules items update --course-id 123 --module-id 456 789 --title "New Title"
  canvas modules items update --course-id 123 --module-id 456 789 --position 2
  canvas modules items update --course-id 123 --module-id 456 789 --published
  canvas modules items update --course-id 123 --module-id 456 789 --move-to-module 555`,
		Args: ExactArgsWithUsage(1, "item-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			itemID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid item ID: %s", args[0])
			}
			opts.ItemID = itemID

			// Track which fields were set
			opts.TitleSet = cmd.Flags().Changed("title")
			opts.PositionSet = cmd.Flags().Changed("position")
			opts.IndentSet = cmd.Flags().Changed("indent")
			opts.NewTabSet = cmd.Flags().Changed("new-tab")
			opts.CompletionTypeSet = cmd.Flags().Changed("completion-type")
			opts.MinScoreSet = cmd.Flags().Changed("min-score")
			opts.PublishedSet = cmd.Flags().Changed("published")
			opts.MoveToModuleSet = cmd.Flags().Changed("move-to-module")

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runModulesItemsUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.ModuleID, "module-id", 0, "Module ID (required)")
	cmd.Flags().StringVar(&opts.Title, "title", "", "New item title")
	cmd.Flags().IntVar(&opts.Position, "position", 0, "New position in the module")
	cmd.Flags().IntVar(&opts.Indent, "indent", 0, "Indent level (0-based)")
	cmd.Flags().BoolVar(&opts.NewTab, "new-tab", false, "Open in new tab (for ExternalTool)")
	cmd.Flags().StringVar(&opts.CompletionType, "completion-type", "", "Completion requirement: must_view, must_contribute, must_submit, must_mark_done, min_score")
	cmd.Flags().Float64Var(&opts.MinScore, "min-score", 0, "Minimum score for min_score completion type")
	cmd.Flags().BoolVar(&opts.Published, "published", false, "Publish the item")
	cmd.Flags().Int64Var(&opts.MoveToModule, "move-to-module", 0, "Move item to a different module")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("module-id")

	return cmd
}

// newModulesItemsDeleteCmd creates the module items delete command
func newModulesItemsDeleteCmd() *cobra.Command {
	opts := &options.ModulesItemsDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <item-id>",
		Short: "Delete a module item",
		Long: `Delete an item from a module.

Examples:
  canvas modules items delete --course-id 123 --module-id 456 789`,
		Args: ExactArgsWithUsage(1, "item-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			itemID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid item ID: %s", args[0])
			}
			opts.ItemID = itemID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runModulesItemsDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.ModuleID, "module-id", 0, "Module ID (required)")
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Skip confirmation prompt")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("module-id")

	return cmd
}

// newModulesItemsDoneCmd creates the module items done command
func newModulesItemsDoneCmd() *cobra.Command {
	opts := &options.ModulesItemsDoneOptions{}

	cmd := &cobra.Command{
		Use:   "done <item-id>",
		Short: "Mark a module item as done",
		Long: `Mark a module item as done (for must_mark_done requirement).

Examples:
  canvas modules items done --course-id 123 --module-id 456 789`,
		Args: ExactArgsWithUsage(1, "item-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			itemID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid item ID: %s", args[0])
			}
			opts.ItemID = itemID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runModulesItemsDone(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.ModuleID, "module-id", 0, "Module ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("module-id")

	return cmd
}

func init() {
	rootCmd.AddCommand(modulesCmd)
	modulesCmd.AddCommand(newModulesListCmd())
	modulesCmd.AddCommand(newModulesGetCmd())
	modulesCmd.AddCommand(newModulesCreateCmd())
	modulesCmd.AddCommand(newModulesUpdateCmd())
	modulesCmd.AddCommand(newModulesDeleteCmd())
	modulesCmd.AddCommand(newModulesRelockCmd())
	modulesCmd.AddCommand(newModulesPublishCmd())
	modulesCmd.AddCommand(newModulesUnpublishCmd())
	modulesCmd.AddCommand(modulesItemsCmd)

	// Items subcommands
	modulesItemsCmd.AddCommand(newModulesItemsListCmd())
	modulesItemsCmd.AddCommand(newModulesItemsGetCmd())
	modulesItemsCmd.AddCommand(newModulesItemsCreateCmd())
	modulesItemsCmd.AddCommand(newModulesItemsUpdateCmd())
	modulesItemsCmd.AddCommand(newModulesItemsDeleteCmd())
	modulesItemsCmd.AddCommand(newModulesItemsDoneCmd())
}

func runModulesList(ctx context.Context, client *api.Client, opts *options.ModulesListOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "modules.list", map[string]interface{}{
		"course_id":   opts.CourseID,
		"search_term": opts.SearchTerm,
		"student_id":  opts.StudentID,
	})

	// Validate course ID exists
	if _, err := validateCourseID(client, opts.CourseID); err != nil {
		logger.LogCommandError(ctx, "modules.list", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return err
	}

	modulesService := api.NewModulesService(client)

	apiOpts := &api.ListModulesOptions{
		Include:    opts.Include,
		SearchTerm: opts.SearchTerm,
		StudentID:  opts.StudentID,
	}

	modules, err := modulesService.List(ctx, opts.CourseID, apiOpts)
	if err != nil {
		logger.LogCommandError(ctx, "modules.list", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to list modules: %w", err)
	}

	if err := formatEmptyOrOutput(modules, "No modules found"); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "modules.list", len(modules))
	return nil
}

func runModulesGet(ctx context.Context, client *api.Client, opts *options.ModulesGetOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "modules.get", map[string]interface{}{
		"course_id": opts.CourseID,
		"module_id": opts.ModuleID,
	})

	// Validate course ID exists
	if _, err := validateCourseID(client, opts.CourseID); err != nil {
		logger.LogCommandError(ctx, "modules.get", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return err
	}

	modulesService := api.NewModulesService(client)

	module, err := modulesService.Get(ctx, opts.CourseID, opts.ModuleID, opts.Include, "")
	if err != nil {
		logger.LogCommandError(ctx, "modules.get", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"module_id": opts.ModuleID,
		})
		return fmt.Errorf("failed to get module: %w", err)
	}

	if err := formatOutput(module, nil); err != nil {
		return fmt.Errorf("failed to print result: %w", err)
	}

	logger.LogCommandComplete(ctx, "modules.get", 1)
	return nil
}

func runModulesCreate(ctx context.Context, client *api.Client, opts *options.ModulesCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "modules.create", map[string]interface{}{
		"course_id": opts.CourseID,
		"name":      opts.Name,
	})

	// Validate course ID exists
	if _, err := validateCourseID(client, opts.CourseID); err != nil {
		logger.LogCommandError(ctx, "modules.create", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return err
	}

	modulesService := api.NewModulesService(client)

	params := &api.CreateModuleParams{
		Name:                      opts.Name,
		UnlockAt:                  opts.UnlockAt,
		Position:                  opts.Position,
		RequireSequentialProgress: opts.RequireSequentialProgress,
		PrerequisiteModuleIDs:     opts.PrerequisiteModuleIDs,
		PublishFinalGrade:         opts.PublishFinalGrade,
	}

	module, err := modulesService.Create(ctx, opts.CourseID, params)
	if err != nil {
		logger.LogCommandError(ctx, "modules.create", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"name":      opts.Name,
		})
		return fmt.Errorf("failed to create module: %w", err)
	}

	if err := formatSuccessOutput(module, "Module created successfully!"); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "modules.create", 1)
	return nil
}

func runModulesUpdate(ctx context.Context, client *api.Client, opts *options.ModulesUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "modules.update", map[string]interface{}{
		"course_id": opts.CourseID,
		"module_id": opts.ModuleID,
	})

	// Validate course ID exists
	if _, err := validateCourseID(client, opts.CourseID); err != nil {
		logger.LogCommandError(ctx, "modules.update", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return err
	}

	modulesService := api.NewModulesService(client)

	params := &api.UpdateModuleParams{}

	// Only set fields that were explicitly provided
	if opts.NameSet {
		params.Name = &opts.Name
	}
	if opts.UnlockAtSet {
		params.UnlockAt = &opts.UnlockAt
	}
	if opts.PositionSet {
		params.Position = &opts.Position
	}
	if opts.RequireSequentialProgressSet {
		params.RequireSequentialProgress = &opts.RequireSequentialProgress
	}
	if opts.PrerequisiteModuleIDsSet {
		params.PrerequisiteModuleIDs = opts.PrerequisiteModuleIDs
	}
	if opts.PublishFinalGradeSet {
		params.PublishFinalGrade = &opts.PublishFinalGrade
	}
	if opts.PublishedSet {
		params.Published = &opts.Published
	}

	module, err := modulesService.Update(ctx, opts.CourseID, opts.ModuleID, params)
	if err != nil {
		logger.LogCommandError(ctx, "modules.update", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"module_id": opts.ModuleID,
		})
		return fmt.Errorf("failed to update module: %w", err)
	}

	if err := formatSuccessOutput(module, "Module updated successfully!"); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "modules.update", 1)
	return nil
}

func runModulesDelete(ctx context.Context, client *api.Client, opts *options.ModulesDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "modules.delete", map[string]interface{}{
		"course_id": opts.CourseID,
		"module_id": opts.ModuleID,
	})

	// Validate course ID exists
	if _, err := validateCourseID(client, opts.CourseID); err != nil {
		logger.LogCommandError(ctx, "modules.delete", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return err
	}

	// Confirm deletion
	confirmed, err := confirmDelete("module", opts.ModuleID, opts.Force)
	if err != nil {
		return err
	}
	if !confirmed {
		fmt.Println("Cancelled.")
		return nil
	}

	modulesService := api.NewModulesService(client)

	if err := modulesService.Delete(ctx, opts.CourseID, opts.ModuleID); err != nil {
		logger.LogCommandError(ctx, "modules.delete", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"module_id": opts.ModuleID,
		})
		return fmt.Errorf("failed to delete module: %w", err)
	}

	fmt.Printf("Module %d deleted successfully\n", opts.ModuleID)

	logger.LogCommandComplete(ctx, "modules.delete", 1)
	return nil
}

func runModulesRelock(ctx context.Context, client *api.Client, opts *options.ModulesRelockOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "modules.relock", map[string]interface{}{
		"course_id": opts.CourseID,
		"module_id": opts.ModuleID,
	})

	// Validate course ID exists
	if _, err := validateCourseID(client, opts.CourseID); err != nil {
		logger.LogCommandError(ctx, "modules.relock", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return err
	}

	modulesService := api.NewModulesService(client)

	module, err := modulesService.Relock(ctx, opts.CourseID, opts.ModuleID)
	if err != nil {
		logger.LogCommandError(ctx, "modules.relock", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"module_id": opts.ModuleID,
		})
		return fmt.Errorf("failed to relock module: %w", err)
	}

	if err := formatSuccessOutput(module, "Module progressions re-locked successfully!"); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "modules.relock", 1)
	return nil
}

func runModulesPublish(ctx context.Context, client *api.Client, opts *options.ModulesPublishOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "modules.publish", map[string]interface{}{
		"course_id": opts.CourseID,
		"module_id": opts.ModuleID,
	})

	// Validate course ID exists
	if _, err := validateCourseID(client, opts.CourseID); err != nil {
		logger.LogCommandError(ctx, "modules.publish", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return err
	}

	modulesService := api.NewModulesService(client)

	published := true
	params := &api.UpdateModuleParams{
		Published: &published,
	}

	module, err := modulesService.Update(ctx, opts.CourseID, opts.ModuleID, params)
	if err != nil {
		logger.LogCommandError(ctx, "modules.publish", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"module_id": opts.ModuleID,
		})
		return fmt.Errorf("failed to publish module: %w", err)
	}

	if err := formatSuccessOutput(module, "Module published successfully!"); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "modules.publish", 1)
	return nil
}

func runModulesUnpublish(ctx context.Context, client *api.Client, opts *options.ModulesUnpublishOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "modules.unpublish", map[string]interface{}{
		"course_id": opts.CourseID,
		"module_id": opts.ModuleID,
	})

	// Validate course ID exists
	if _, err := validateCourseID(client, opts.CourseID); err != nil {
		logger.LogCommandError(ctx, "modules.unpublish", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return err
	}

	modulesService := api.NewModulesService(client)

	published := false
	params := &api.UpdateModuleParams{
		Published: &published,
	}

	module, err := modulesService.Update(ctx, opts.CourseID, opts.ModuleID, params)
	if err != nil {
		logger.LogCommandError(ctx, "modules.unpublish", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"module_id": opts.ModuleID,
		})
		return fmt.Errorf("failed to unpublish module: %w", err)
	}

	if err := formatSuccessOutput(module, "Module unpublished successfully!"); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "modules.unpublish", 1)
	return nil
}

func runModulesItemsList(ctx context.Context, client *api.Client, opts *options.ModulesItemsListOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "modules.items.list", map[string]interface{}{
		"course_id":   opts.CourseID,
		"module_id":   opts.ModuleID,
		"search_term": opts.SearchTerm,
	})

	// Validate course ID exists
	if _, err := validateCourseID(client, opts.CourseID); err != nil {
		logger.LogCommandError(ctx, "modules.items.list", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return err
	}

	modulesService := api.NewModulesService(client)

	apiOpts := &api.ListModuleItemsOptions{
		Include:    opts.Include,
		SearchTerm: opts.SearchTerm,
		StudentID:  opts.StudentID,
	}

	items, err := modulesService.ListItems(ctx, opts.CourseID, opts.ModuleID, apiOpts)
	if err != nil {
		logger.LogCommandError(ctx, "modules.items.list", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"module_id": opts.ModuleID,
		})
		return fmt.Errorf("failed to list module items: %w", err)
	}

	if err := formatEmptyOrOutput(items, "No items found in this module"); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "modules.items.list", len(items))
	return nil
}

func runModulesItemsGet(ctx context.Context, client *api.Client, opts *options.ModulesItemsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "modules.items.get", map[string]interface{}{
		"course_id": opts.CourseID,
		"module_id": opts.ModuleID,
		"item_id":   opts.ItemID,
	})

	// Validate course ID exists
	if _, err := validateCourseID(client, opts.CourseID); err != nil {
		logger.LogCommandError(ctx, "modules.items.get", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return err
	}

	modulesService := api.NewModulesService(client)

	item, err := modulesService.GetItem(ctx, opts.CourseID, opts.ModuleID, opts.ItemID, opts.Include, "")
	if err != nil {
		logger.LogCommandError(ctx, "modules.items.get", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"module_id": opts.ModuleID,
			"item_id":   opts.ItemID,
		})
		return fmt.Errorf("failed to get module item: %w", err)
	}

	if err := formatOutput(item, nil); err != nil {
		return fmt.Errorf("failed to print result: %w", err)
	}

	logger.LogCommandComplete(ctx, "modules.items.get", 1)
	return nil
}

func runModulesItemsCreate(ctx context.Context, client *api.Client, opts *options.ModulesItemsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "modules.items.create", map[string]interface{}{
		"course_id": opts.CourseID,
		"module_id": opts.ModuleID,
		"type":      opts.Type,
	})

	// Validate course ID exists
	if _, err := validateCourseID(client, opts.CourseID); err != nil {
		logger.LogCommandError(ctx, "modules.items.create", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return err
	}

	modulesService := api.NewModulesService(client)

	params := &api.CreateModuleItemParams{
		Type:        opts.Type,
		Title:       opts.Title,
		ContentID:   opts.ContentID,
		PageURL:     opts.PageURL,
		ExternalURL: opts.ExternalURL,
		NewTab:      opts.NewTab,
		Indent:      opts.Indent,
	}

	if opts.CompletionType != "" {
		params.CompletionRequirement = &api.CompletionRequirementParams{
			Type:     opts.CompletionType,
			MinScore: opts.MinScore,
		}
	}

	item, err := modulesService.CreateItem(ctx, opts.CourseID, opts.ModuleID, params)
	if err != nil {
		logger.LogCommandError(ctx, "modules.items.create", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"module_id": opts.ModuleID,
			"type":      opts.Type,
		})
		return fmt.Errorf("failed to create module item: %w", err)
	}

	if err := formatSuccessOutput(item, "Module item created successfully!"); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "modules.items.create", 1)
	return nil
}

func runModulesItemsUpdate(ctx context.Context, client *api.Client, opts *options.ModulesItemsUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "modules.items.update", map[string]interface{}{
		"course_id": opts.CourseID,
		"module_id": opts.ModuleID,
		"item_id":   opts.ItemID,
	})

	// Validate course ID exists
	if _, err := validateCourseID(client, opts.CourseID); err != nil {
		logger.LogCommandError(ctx, "modules.items.update", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return err
	}

	modulesService := api.NewModulesService(client)

	params := &api.UpdateModuleItemParams{}

	// Only set fields that were explicitly provided
	if opts.TitleSet {
		params.Title = &opts.Title
	}
	if opts.PositionSet {
		params.Position = &opts.Position
	}
	if opts.IndentSet {
		params.Indent = &opts.Indent
	}
	if opts.NewTabSet {
		params.NewTab = &opts.NewTab
	}
	if opts.PublishedSet {
		params.Published = &opts.Published
	}
	if opts.MoveToModuleSet {
		params.MoveToModuleID = &opts.MoveToModule
	}
	if opts.CompletionTypeSet {
		params.CompletionRequirement = &api.CompletionRequirementParams{
			Type:     opts.CompletionType,
			MinScore: opts.MinScore,
		}
	}

	item, err := modulesService.UpdateItem(ctx, opts.CourseID, opts.ModuleID, opts.ItemID, params)
	if err != nil {
		logger.LogCommandError(ctx, "modules.items.update", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"module_id": opts.ModuleID,
			"item_id":   opts.ItemID,
		})
		return fmt.Errorf("failed to update module item: %w", err)
	}

	if err := formatSuccessOutput(item, "Module item updated successfully!"); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "modules.items.update", 1)
	return nil
}

func runModulesItemsDelete(ctx context.Context, client *api.Client, opts *options.ModulesItemsDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "modules.items.delete", map[string]interface{}{
		"course_id": opts.CourseID,
		"module_id": opts.ModuleID,
		"item_id":   opts.ItemID,
	})

	// Validate course ID exists
	if _, err := validateCourseID(client, opts.CourseID); err != nil {
		logger.LogCommandError(ctx, "modules.items.delete", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return err
	}

	// Confirm deletion
	confirmed, err := confirmDelete("module item", opts.ItemID, opts.Force)
	if err != nil {
		return err
	}
	if !confirmed {
		fmt.Println("Cancelled.")
		return nil
	}

	modulesService := api.NewModulesService(client)

	if err := modulesService.DeleteItem(ctx, opts.CourseID, opts.ModuleID, opts.ItemID); err != nil {
		logger.LogCommandError(ctx, "modules.items.delete", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"module_id": opts.ModuleID,
			"item_id":   opts.ItemID,
		})
		return fmt.Errorf("failed to delete module item: %w", err)
	}

	fmt.Printf("Module item %d deleted successfully\n", opts.ItemID)

	logger.LogCommandComplete(ctx, "modules.items.delete", 1)
	return nil
}

func runModulesItemsDone(ctx context.Context, client *api.Client, opts *options.ModulesItemsDoneOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "modules.items.done", map[string]interface{}{
		"course_id": opts.CourseID,
		"module_id": opts.ModuleID,
		"item_id":   opts.ItemID,
	})

	// Validate course ID exists
	if _, err := validateCourseID(client, opts.CourseID); err != nil {
		logger.LogCommandError(ctx, "modules.items.done", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return err
	}

	modulesService := api.NewModulesService(client)

	if err := modulesService.MarkItemDone(ctx, opts.CourseID, opts.ModuleID, opts.ItemID); err != nil {
		logger.LogCommandError(ctx, "modules.items.done", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"module_id": opts.ModuleID,
			"item_id":   opts.ItemID,
		})
		return fmt.Errorf("failed to mark item as done: %w", err)
	}

	fmt.Printf("Module item %d marked as done\n", opts.ItemID)

	logger.LogCommandComplete(ctx, "modules.items.done", 1)
	return nil
}
