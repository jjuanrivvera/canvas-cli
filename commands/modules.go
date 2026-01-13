package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var (
	modulesCourseID   int64
	modulesInclude    []string
	modulesSearchTerm string
	modulesStudentID  string
	modulesForce      bool
	// Create/Update flags
	modulesName                      string
	modulesUnlockAt                  string
	modulesPosition                  int
	modulesRequireSequentialProgress bool
	modulesPrerequisiteModuleIDs     []int64
	modulesPublishFinalGrade         bool
	modulesPublished                 bool
	// Module items flags
	modulesModuleID           int64
	modulesItemType           string
	modulesItemTitle          string
	modulesItemContentID      int64
	modulesItemPageURL        string
	modulesItemExternalURL    string
	modulesItemNewTab         bool
	modulesItemIndent         int
	modulesItemCompletionType string
	modulesItemMinScore       float64
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

// modulesListCmd represents the modules list command
var modulesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List modules in a course",
	Long: `List all modules in a Canvas course.

Examples:
  canvas modules list --course-id 123
  canvas modules list --course-id 123 --include items
  canvas modules list --course-id 123 --search "Week"
  canvas modules list --course-id 123 --student-id 789`,
	RunE: runModulesList,
}

// modulesGetCmd represents the modules get command
var modulesGetCmd = &cobra.Command{
	Use:   "get <module-id>",
	Short: "Get details of a specific module",
	Long: `Get details of a specific module by ID.

Examples:
  canvas modules get --course-id 123 456
  canvas modules get --course-id 123 456 --include items,content_details`,
	Args: cobra.ExactArgs(1),
	RunE: runModulesGet,
}

// modulesCreateCmd represents the modules create command
var modulesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new module",
	Long: `Create a new module in a course.

Examples:
  canvas modules create --course-id 123 --name "Week 1"
  canvas modules create --course-id 123 --name "Week 2" --position 2
  canvas modules create --course-id 123 --name "Unit 2" --prerequisite-modules 1,2
  canvas modules create --course-id 123 --name "Final" --require-sequential-progress`,
	RunE: runModulesCreate,
}

// modulesUpdateCmd represents the modules update command
var modulesUpdateCmd = &cobra.Command{
	Use:   "update <module-id>",
	Short: "Update an existing module",
	Long: `Update an existing module.

Examples:
  canvas modules update --course-id 123 456 --name "Updated Name"
  canvas modules update --course-id 123 456 --published
  canvas modules update --course-id 123 456 --position 3`,
	Args: cobra.ExactArgs(1),
	RunE: runModulesUpdate,
}

// modulesDeleteCmd represents the modules delete command
var modulesDeleteCmd = &cobra.Command{
	Use:   "delete <module-id>",
	Short: "Delete a module",
	Long: `Delete a module from a course.

Examples:
  canvas modules delete --course-id 123 456`,
	Args: cobra.ExactArgs(1),
	RunE: runModulesDelete,
}

// modulesRelockCmd represents the modules relock command
var modulesRelockCmd = &cobra.Command{
	Use:   "relock <module-id>",
	Short: "Re-lock module progressions",
	Long: `Re-lock module progressions to their default locked state.

This recalculates progressions based on current requirements. Adding progression
requirements to an active course will not lock students out of modules they have
already unlocked unless this action is called.

Examples:
  canvas modules relock --course-id 123 456`,
	Args: cobra.ExactArgs(1),
	RunE: runModulesRelock,
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

// modulesItemsListCmd represents the module items list command
var modulesItemsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List items in a module",
	Long: `List all items in a module.

Examples:
  canvas modules items list --course-id 123 --module-id 456
  canvas modules items list --course-id 123 --module-id 456 --include content_details
  canvas modules items list --course-id 123 --module-id 456 --search "Quiz"`,
	RunE: runModulesItemsList,
}

// modulesItemsGetCmd represents the module items get command
var modulesItemsGetCmd = &cobra.Command{
	Use:   "get <item-id>",
	Short: "Get details of a module item",
	Long: `Get details of a specific module item.

Examples:
  canvas modules items get --course-id 123 --module-id 456 789
  canvas modules items get --course-id 123 --module-id 456 789 --include content_details`,
	Args: cobra.ExactArgs(1),
	RunE: runModulesItemsGet,
}

// modulesItemsCreateCmd represents the module items create command
var modulesItemsCreateCmd = &cobra.Command{
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
	RunE: runModulesItemsCreate,
}

// modulesItemsDeleteCmd represents the module items delete command
var modulesItemsDeleteCmd = &cobra.Command{
	Use:   "delete <item-id>",
	Short: "Delete a module item",
	Long: `Delete an item from a module.

Examples:
  canvas modules items delete --course-id 123 --module-id 456 789`,
	Args: cobra.ExactArgs(1),
	RunE: runModulesItemsDelete,
}

// modulesItemsDoneCmd represents the module items done command
var modulesItemsDoneCmd = &cobra.Command{
	Use:   "done <item-id>",
	Short: "Mark a module item as done",
	Long: `Mark a module item as done (for must_mark_done requirement).

Examples:
  canvas modules items done --course-id 123 --module-id 456 789`,
	Args: cobra.ExactArgs(1),
	RunE: runModulesItemsDone,
}

func init() {
	rootCmd.AddCommand(modulesCmd)
	modulesCmd.AddCommand(modulesListCmd)
	modulesCmd.AddCommand(modulesGetCmd)
	modulesCmd.AddCommand(modulesCreateCmd)
	modulesCmd.AddCommand(modulesUpdateCmd)
	modulesCmd.AddCommand(modulesDeleteCmd)
	modulesCmd.AddCommand(modulesRelockCmd)
	modulesCmd.AddCommand(modulesItemsCmd)

	// Items subcommands
	modulesItemsCmd.AddCommand(modulesItemsListCmd)
	modulesItemsCmd.AddCommand(modulesItemsGetCmd)
	modulesItemsCmd.AddCommand(modulesItemsCreateCmd)
	modulesItemsCmd.AddCommand(modulesItemsDeleteCmd)
	modulesItemsCmd.AddCommand(modulesItemsDoneCmd)

	// List flags
	modulesListCmd.Flags().Int64Var(&modulesCourseID, "course-id", 0, "Course ID (required)")
	modulesListCmd.Flags().StringSliceVar(&modulesInclude, "include", []string{}, "Additional data to include (items, content_details)")
	modulesListCmd.Flags().StringVar(&modulesSearchTerm, "search", "", "Search by module name")
	modulesListCmd.Flags().StringVar(&modulesStudentID, "student-id", "", "Get completion info for this student")
	modulesListCmd.MarkFlagRequired("course-id")

	// Get flags
	modulesGetCmd.Flags().Int64Var(&modulesCourseID, "course-id", 0, "Course ID (required)")
	modulesGetCmd.Flags().StringSliceVar(&modulesInclude, "include", []string{}, "Additional data to include (items, content_details)")
	modulesGetCmd.Flags().StringVar(&modulesStudentID, "student-id", "", "Get completion info for this student")
	modulesGetCmd.MarkFlagRequired("course-id")

	// Create flags
	modulesCreateCmd.Flags().Int64Var(&modulesCourseID, "course-id", 0, "Course ID (required)")
	modulesCreateCmd.Flags().StringVar(&modulesName, "name", "", "Module name (required)")
	modulesCreateCmd.Flags().StringVar(&modulesUnlockAt, "unlock-at", "", "Date to unlock the module (ISO 8601)")
	modulesCreateCmd.Flags().IntVar(&modulesPosition, "position", 0, "Position in the course (1-based)")
	modulesCreateCmd.Flags().BoolVar(&modulesRequireSequentialProgress, "require-sequential-progress", false, "Require sequential progress")
	modulesCreateCmd.Flags().Int64SliceVar(&modulesPrerequisiteModuleIDs, "prerequisite-modules", []int64{}, "IDs of prerequisite modules")
	modulesCreateCmd.Flags().BoolVar(&modulesPublishFinalGrade, "publish-final-grade", false, "Publish final grade on completion")
	modulesCreateCmd.MarkFlagRequired("course-id")
	modulesCreateCmd.MarkFlagRequired("name")

	// Update flags
	modulesUpdateCmd.Flags().Int64Var(&modulesCourseID, "course-id", 0, "Course ID (required)")
	modulesUpdateCmd.Flags().StringVar(&modulesName, "name", "", "New module name")
	modulesUpdateCmd.Flags().StringVar(&modulesUnlockAt, "unlock-at", "", "Date to unlock the module (ISO 8601)")
	modulesUpdateCmd.Flags().IntVar(&modulesPosition, "position", 0, "New position in the course")
	modulesUpdateCmd.Flags().BoolVar(&modulesRequireSequentialProgress, "require-sequential-progress", false, "Require sequential progress")
	modulesUpdateCmd.Flags().Int64SliceVar(&modulesPrerequisiteModuleIDs, "prerequisite-modules", []int64{}, "IDs of prerequisite modules")
	modulesUpdateCmd.Flags().BoolVar(&modulesPublishFinalGrade, "publish-final-grade", false, "Publish final grade on completion")
	modulesUpdateCmd.Flags().BoolVar(&modulesPublished, "published", false, "Publish the module")
	modulesUpdateCmd.MarkFlagRequired("course-id")

	// Delete flags
	modulesDeleteCmd.Flags().Int64Var(&modulesCourseID, "course-id", 0, "Course ID (required)")
	modulesDeleteCmd.Flags().BoolVarP(&modulesForce, "force", "f", false, "Skip confirmation prompt")
	modulesDeleteCmd.MarkFlagRequired("course-id")

	// Relock flags
	modulesRelockCmd.Flags().Int64Var(&modulesCourseID, "course-id", 0, "Course ID (required)")
	modulesRelockCmd.MarkFlagRequired("course-id")

	// Items list flags
	modulesItemsListCmd.Flags().Int64Var(&modulesCourseID, "course-id", 0, "Course ID (required)")
	modulesItemsListCmd.Flags().Int64Var(&modulesModuleID, "module-id", 0, "Module ID (required)")
	modulesItemsListCmd.Flags().StringSliceVar(&modulesInclude, "include", []string{}, "Additional data to include (content_details)")
	modulesItemsListCmd.Flags().StringVar(&modulesSearchTerm, "search", "", "Search by item title")
	modulesItemsListCmd.Flags().StringVar(&modulesStudentID, "student-id", "", "Get completion info for this student")
	modulesItemsListCmd.MarkFlagRequired("course-id")
	modulesItemsListCmd.MarkFlagRequired("module-id")

	// Items get flags
	modulesItemsGetCmd.Flags().Int64Var(&modulesCourseID, "course-id", 0, "Course ID (required)")
	modulesItemsGetCmd.Flags().Int64Var(&modulesModuleID, "module-id", 0, "Module ID (required)")
	modulesItemsGetCmd.Flags().StringSliceVar(&modulesInclude, "include", []string{}, "Additional data to include (content_details)")
	modulesItemsGetCmd.Flags().StringVar(&modulesStudentID, "student-id", "", "Get completion info for this student")
	modulesItemsGetCmd.MarkFlagRequired("course-id")
	modulesItemsGetCmd.MarkFlagRequired("module-id")

	// Items create flags
	modulesItemsCreateCmd.Flags().Int64Var(&modulesCourseID, "course-id", 0, "Course ID (required)")
	modulesItemsCreateCmd.Flags().Int64Var(&modulesModuleID, "module-id", 0, "Module ID (required)")
	modulesItemsCreateCmd.Flags().StringVar(&modulesItemType, "type", "", "Item type: File, Page, Discussion, Assignment, Quiz, SubHeader, ExternalUrl, ExternalTool (required)")
	modulesItemsCreateCmd.Flags().StringVar(&modulesItemTitle, "title", "", "Item title")
	modulesItemsCreateCmd.Flags().Int64Var(&modulesItemContentID, "content-id", 0, "Content ID (for File, Discussion, Assignment, Quiz, ExternalTool)")
	modulesItemsCreateCmd.Flags().StringVar(&modulesItemPageURL, "page-url", "", "Page URL slug (for Page type)")
	modulesItemsCreateCmd.Flags().StringVar(&modulesItemExternalURL, "external-url", "", "External URL (for ExternalUrl, ExternalTool)")
	modulesItemsCreateCmd.Flags().BoolVar(&modulesItemNewTab, "new-tab", false, "Open in new tab (for ExternalTool)")
	modulesItemsCreateCmd.Flags().IntVar(&modulesItemIndent, "indent", 0, "Indent level (0-based)")
	modulesItemsCreateCmd.Flags().IntVar(&modulesPosition, "position", 0, "Position in the module")
	modulesItemsCreateCmd.Flags().StringVar(&modulesItemCompletionType, "completion-type", "", "Completion requirement: must_view, must_contribute, must_submit, must_mark_done, min_score")
	modulesItemsCreateCmd.Flags().Float64Var(&modulesItemMinScore, "min-score", 0, "Minimum score for min_score completion type")
	modulesItemsCreateCmd.MarkFlagRequired("course-id")
	modulesItemsCreateCmd.MarkFlagRequired("module-id")
	modulesItemsCreateCmd.MarkFlagRequired("type")

	// Items delete flags
	modulesItemsDeleteCmd.Flags().Int64Var(&modulesCourseID, "course-id", 0, "Course ID (required)")
	modulesItemsDeleteCmd.Flags().Int64Var(&modulesModuleID, "module-id", 0, "Module ID (required)")
	modulesItemsDeleteCmd.Flags().BoolVarP(&modulesForce, "force", "f", false, "Skip confirmation prompt")
	modulesItemsDeleteCmd.MarkFlagRequired("course-id")
	modulesItemsDeleteCmd.MarkFlagRequired("module-id")

	// Items done flags
	modulesItemsDoneCmd.Flags().Int64Var(&modulesCourseID, "course-id", 0, "Course ID (required)")
	modulesItemsDoneCmd.Flags().Int64Var(&modulesModuleID, "module-id", 0, "Module ID (required)")
	modulesItemsDoneCmd.MarkFlagRequired("course-id")
	modulesItemsDoneCmd.MarkFlagRequired("module-id")
}

func runModulesList(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Validate course ID exists
	if _, err := validateCourseID(client, modulesCourseID); err != nil {
		return err
	}

	modulesService := api.NewModulesService(client)

	opts := &api.ListModulesOptions{
		Include:    modulesInclude,
		SearchTerm: modulesSearchTerm,
		StudentID:  modulesStudentID,
	}

	ctx := context.Background()
	modules, err := modulesService.List(ctx, modulesCourseID, opts)
	if err != nil {
		return fmt.Errorf("failed to list modules: %w", err)
	}

	if len(modules) == 0 {
		fmt.Println("No modules found")
		return nil
	}

	// Format and display modules
	fmt.Printf("Found %d modules:\n\n", len(modules))
	return formatOutput(modules, nil)
}

func runModulesGet(cmd *cobra.Command, args []string) error {
	moduleID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid module ID: %s", args[0])
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Validate course ID exists
	if _, err := validateCourseID(client, modulesCourseID); err != nil {
		return err
	}

	modulesService := api.NewModulesService(client)

	ctx := context.Background()
	module, err := modulesService.Get(ctx, modulesCourseID, moduleID, modulesInclude, modulesStudentID)
	if err != nil {
		return fmt.Errorf("failed to get module: %w", err)
	}

	// Format and display module details
	return formatOutput(module, nil)
}

func runModulesCreate(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Validate course ID exists
	if _, err := validateCourseID(client, modulesCourseID); err != nil {
		return err
	}

	modulesService := api.NewModulesService(client)

	params := &api.CreateModuleParams{
		Name:                      modulesName,
		UnlockAt:                  modulesUnlockAt,
		Position:                  modulesPosition,
		RequireSequentialProgress: modulesRequireSequentialProgress,
		PrerequisiteModuleIDs:     modulesPrerequisiteModuleIDs,
		PublishFinalGrade:         modulesPublishFinalGrade,
	}

	ctx := context.Background()
	module, err := modulesService.Create(ctx, modulesCourseID, params)
	if err != nil {
		return fmt.Errorf("failed to create module: %w", err)
	}

	fmt.Println("Module created successfully!")
	displayModule(module)

	return nil
}

func runModulesUpdate(cmd *cobra.Command, args []string) error {
	moduleID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid module ID: %s", args[0])
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Validate course ID exists
	if _, err := validateCourseID(client, modulesCourseID); err != nil {
		return err
	}

	modulesService := api.NewModulesService(client)

	params := &api.UpdateModuleParams{}

	// Only set fields that were explicitly provided
	if cmd.Flags().Changed("name") {
		params.Name = &modulesName
	}
	if cmd.Flags().Changed("unlock-at") {
		params.UnlockAt = &modulesUnlockAt
	}
	if cmd.Flags().Changed("position") {
		params.Position = &modulesPosition
	}
	if cmd.Flags().Changed("require-sequential-progress") {
		params.RequireSequentialProgress = &modulesRequireSequentialProgress
	}
	if cmd.Flags().Changed("prerequisite-modules") {
		params.PrerequisiteModuleIDs = modulesPrerequisiteModuleIDs
	}
	if cmd.Flags().Changed("publish-final-grade") {
		params.PublishFinalGrade = &modulesPublishFinalGrade
	}
	if cmd.Flags().Changed("published") {
		params.Published = &modulesPublished
	}

	ctx := context.Background()
	module, err := modulesService.Update(ctx, modulesCourseID, moduleID, params)
	if err != nil {
		return fmt.Errorf("failed to update module: %w", err)
	}

	fmt.Println("Module updated successfully!")
	displayModule(module)

	return nil
}

func runModulesDelete(cmd *cobra.Command, args []string) error {
	moduleID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid module ID: %s", args[0])
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Validate course ID exists
	if _, err := validateCourseID(client, modulesCourseID); err != nil {
		return err
	}

	// Confirm deletion
	confirmed, err := confirmDelete("module", moduleID, modulesForce)
	if err != nil {
		return err
	}
	if !confirmed {
		fmt.Println("Cancelled.")
		return nil
	}

	modulesService := api.NewModulesService(client)

	ctx := context.Background()
	if err := modulesService.Delete(ctx, modulesCourseID, moduleID); err != nil {
		return fmt.Errorf("failed to delete module: %w", err)
	}

	fmt.Printf("Module %d deleted successfully\n", moduleID)
	return nil
}

func runModulesRelock(cmd *cobra.Command, args []string) error {
	moduleID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid module ID: %s", args[0])
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Validate course ID exists
	if _, err := validateCourseID(client, modulesCourseID); err != nil {
		return err
	}

	modulesService := api.NewModulesService(client)

	ctx := context.Background()
	module, err := modulesService.Relock(ctx, modulesCourseID, moduleID)
	if err != nil {
		return fmt.Errorf("failed to relock module: %w", err)
	}

	fmt.Println("Module progressions re-locked successfully!")
	displayModule(module)

	return nil
}

func runModulesItemsList(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Validate course ID exists
	if _, err := validateCourseID(client, modulesCourseID); err != nil {
		return err
	}

	modulesService := api.NewModulesService(client)

	opts := &api.ListModuleItemsOptions{
		Include:    modulesInclude,
		SearchTerm: modulesSearchTerm,
		StudentID:  modulesStudentID,
	}

	ctx := context.Background()
	items, err := modulesService.ListItems(ctx, modulesCourseID, modulesModuleID, opts)
	if err != nil {
		return fmt.Errorf("failed to list module items: %w", err)
	}

	if len(items) == 0 {
		fmt.Println("No items found in this module")
		return nil
	}

	// Format and display items
	fmt.Printf("Found %d items:\n\n", len(items))
	return formatOutput(items, nil)
}

func runModulesItemsGet(cmd *cobra.Command, args []string) error {
	itemID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid item ID: %s", args[0])
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Validate course ID exists
	if _, err := validateCourseID(client, modulesCourseID); err != nil {
		return err
	}

	modulesService := api.NewModulesService(client)

	ctx := context.Background()
	item, err := modulesService.GetItem(ctx, modulesCourseID, modulesModuleID, itemID, modulesInclude, modulesStudentID)
	if err != nil {
		return fmt.Errorf("failed to get module item: %w", err)
	}

	// Format and display item details
	return formatOutput(item, nil)
}

func runModulesItemsCreate(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Validate course ID exists
	if _, err := validateCourseID(client, modulesCourseID); err != nil {
		return err
	}

	modulesService := api.NewModulesService(client)

	params := &api.CreateModuleItemParams{
		Type:        modulesItemType,
		Title:       modulesItemTitle,
		ContentID:   modulesItemContentID,
		PageURL:     modulesItemPageURL,
		ExternalURL: modulesItemExternalURL,
		NewTab:      modulesItemNewTab,
		Indent:      modulesItemIndent,
		Position:    modulesPosition,
	}

	if modulesItemCompletionType != "" {
		params.CompletionRequirement = &api.CompletionRequirementParams{
			Type:     modulesItemCompletionType,
			MinScore: modulesItemMinScore,
		}
	}

	ctx := context.Background()
	item, err := modulesService.CreateItem(ctx, modulesCourseID, modulesModuleID, params)
	if err != nil {
		return fmt.Errorf("failed to create module item: %w", err)
	}

	fmt.Println("Module item created successfully!")
	displayModuleItem(item, "")

	return nil
}

func runModulesItemsDelete(cmd *cobra.Command, args []string) error {
	itemID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid item ID: %s", args[0])
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Validate course ID exists
	if _, err := validateCourseID(client, modulesCourseID); err != nil {
		return err
	}

	// Confirm deletion
	confirmed, err := confirmDelete("module item", itemID, modulesForce)
	if err != nil {
		return err
	}
	if !confirmed {
		fmt.Println("Cancelled.")
		return nil
	}

	modulesService := api.NewModulesService(client)

	ctx := context.Background()
	if err := modulesService.DeleteItem(ctx, modulesCourseID, modulesModuleID, itemID); err != nil {
		return fmt.Errorf("failed to delete module item: %w", err)
	}

	fmt.Printf("Module item %d deleted successfully\n", itemID)
	return nil
}

func runModulesItemsDone(cmd *cobra.Command, args []string) error {
	itemID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid item ID: %s", args[0])
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Validate course ID exists
	if _, err := validateCourseID(client, modulesCourseID); err != nil {
		return err
	}

	modulesService := api.NewModulesService(client)

	ctx := context.Background()
	if err := modulesService.MarkItemDone(ctx, modulesCourseID, modulesModuleID, itemID); err != nil {
		return fmt.Errorf("failed to mark item as done: %w", err)
	}

	fmt.Printf("Module item %d marked as done\n", itemID)
	return nil
}

func displayModule(module *api.Module) {
	stateIcon := "📦"
	if module.Published {
		stateIcon = "📗"
	}

	fmt.Printf("%s %s\n", stateIcon, module.Name)
	fmt.Printf("   ID: %d\n", module.ID)
	fmt.Printf("   Position: %d\n", module.Position)
	fmt.Printf("   State: %s\n", module.WorkflowState)
	fmt.Printf("   Items: %d\n", module.ItemsCount)

	if module.Published {
		fmt.Printf("   Published: Yes\n")
	} else {
		fmt.Printf("   Published: No\n")
	}

	if module.RequireSequentialProgress {
		fmt.Printf("   Sequential Progress: Required\n")
	}

	if len(module.PrerequisiteModuleIDs) > 0 {
		fmt.Printf("   Prerequisites: %v\n", module.PrerequisiteModuleIDs)
	}

	if module.UnlockAt != nil {
		fmt.Printf("   Unlock At: %s\n", module.UnlockAt.Format("2006-01-02 15:04"))
	}

	if module.State != "" {
		fmt.Printf("   Student State: %s\n", module.State)
	}

	fmt.Println()
}

func displayModuleItem(item *api.ModuleItem, indent string) {
	typeIcon := getItemTypeIcon(item.Type)

	fmt.Printf("%s%s %s\n", indent, typeIcon, item.Title)
	fmt.Printf("%s   ID: %d\n", indent, item.ID)
	fmt.Printf("%s   Type: %s\n", indent, item.Type)
	fmt.Printf("%s   Position: %d\n", indent, item.Position)

	if item.Indent > 0 {
		fmt.Printf("%s   Indent: %d\n", indent, item.Indent)
	}

	if item.Published {
		fmt.Printf("%s   Published: Yes\n", indent)
	} else {
		fmt.Printf("%s   Published: No\n", indent)
	}

	if item.CompletionRequirement != nil {
		fmt.Printf("%s   Completion: %s", indent, item.CompletionRequirement.Type)
		if item.CompletionRequirement.MinScore > 0 {
			fmt.Printf(" (min: %.0f)", item.CompletionRequirement.MinScore)
		}
		if item.CompletionRequirement.Completed {
			fmt.Printf(" [Completed]")
		}
		fmt.Println()
	}

	if item.ExternalURL != "" {
		fmt.Printf("%s   URL: %s\n", indent, item.ExternalURL)
	}

	if item.PageURL != "" {
		fmt.Printf("%s   Page: %s\n", indent, item.PageURL)
	}

	fmt.Println()
}

func getItemTypeIcon(itemType string) string {
	switch itemType {
	case "File":
		return "📄"
	case "Page":
		return "📝"
	case "Discussion":
		return "💬"
	case "Assignment":
		return "📋"
	case "Quiz":
		return "❓"
	case "SubHeader":
		return "📑"
	case "ExternalUrl":
		return "🔗"
	case "ExternalTool":
		return "🔧"
	default:
		return "📌"
	}
}
