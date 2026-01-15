package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var (
	extToolsCourseID      int64
	extToolsAccountID     int64
	extToolsSearch        string
	extToolsSelectable    bool
	extToolsIncludeParent bool

	// Create/Update flags
	extToolsName         string
	extToolsURL          string
	extToolsDomain       string
	extToolsConsumerKey  string
	extToolsSharedSecret string
	extToolsPrivacyLevel string
	extToolsDescription  string
	extToolsConfigType   string
	extToolsConfigURL    string
	extToolsConfigXML    string
	extToolsJSONFile     string

	// Launch flags
	extToolsLaunchType   string
	extToolsAssignmentID int64
	extToolsModuleItemID int64

	// Delete flags
	extToolsForce bool
)

var externalToolsCmd = &cobra.Command{
	Use:     "external-tools",
	Aliases: []string{"lti", "tools"},
	Short:   "Manage external tools (LTI)",
	Long: `Manage Canvas external tools (LTI integrations).

External tools allow you to integrate third-party tools and services with Canvas.

Examples:
  canvas external-tools list --course-id 123
  canvas external-tools get 456 --course-id 123
  canvas external-tools create --course-id 123 --name "My Tool" --url https://tool.example.com`,
}

var extToolsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List external tools",
	Long: `List external tools for a course or account.

If neither --account-id nor --course-id is specified, uses default account.

Examples:
  canvas external-tools list                   # Uses default account
  canvas external-tools list --course-id 123
  canvas external-tools list --account-id 1
  canvas external-tools list --course-id 123 --search "quiz"
  canvas external-tools list --course-id 123 --selectable --include-parents`,
	RunE: runExtToolsList,
}

var extToolsGetCmd = &cobra.Command{
	Use:   "get <tool-id>",
	Short: "Get an external tool",
	Long: `Get details of a specific external tool.

Examples:
  canvas external-tools get 456 --course-id 123
  canvas external-tools get 456 --account-id 1`,
	Args: ExactArgsWithUsage(1, "tool-id"),
	RunE: runExtToolsGet,
}

var extToolsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an external tool",
	Long: `Create a new external tool in a course or account.

Examples:
  canvas external-tools create --course-id 123 --name "My Tool" --url https://tool.example.com
  canvas external-tools create --course-id 123 --name "LTI Tool" --consumer-key key123 --shared-secret secret123
  canvas external-tools create --course-id 123 --json tool-config.json`,
	RunE: runExtToolsCreate,
}

var extToolsUpdateCmd = &cobra.Command{
	Use:   "update <tool-id>",
	Short: "Update an external tool",
	Long: `Update an existing external tool.

Examples:
  canvas external-tools update 456 --course-id 123 --name "Updated Name"
  canvas external-tools update 456 --course-id 123 --url https://new-url.example.com`,
	Args: ExactArgsWithUsage(1, "tool-id"),
	RunE: runExtToolsUpdate,
}

var extToolsDeleteCmd = &cobra.Command{
	Use:   "delete <tool-id>",
	Short: "Delete an external tool",
	Long: `Delete an external tool from a course or account.

Examples:
  canvas external-tools delete 456 --course-id 123
  canvas external-tools delete 456 --account-id 1`,
	Args: ExactArgsWithUsage(1, "tool-id"),
	RunE: runExtToolsDelete,
}

var extToolsLaunchCmd = &cobra.Command{
	Use:   "launch <tool-id>",
	Short: "Get sessionless launch URL",
	Long: `Get a sessionless launch URL for an external tool.

Launch Types:
  - course_navigation: Launch from course navigation
  - account_navigation: Launch from account navigation
  - assessment: Launch for an assessment (requires --assignment-id)
  - module_item: Launch for a module item (requires --module-item-id)

Examples:
  canvas external-tools launch 456 --course-id 123
  canvas external-tools launch 456 --course-id 123 --launch-type assessment --assignment-id 789`,
	Args: ExactArgsWithUsage(1, "tool-id"),
	RunE: runExtToolsLaunch,
}

func init() {
	rootCmd.AddCommand(externalToolsCmd)
	externalToolsCmd.AddCommand(extToolsListCmd)
	externalToolsCmd.AddCommand(extToolsGetCmd)
	externalToolsCmd.AddCommand(extToolsCreateCmd)
	externalToolsCmd.AddCommand(extToolsUpdateCmd)
	externalToolsCmd.AddCommand(extToolsDeleteCmd)
	externalToolsCmd.AddCommand(extToolsLaunchCmd)

	// List flags
	extToolsListCmd.Flags().Int64Var(&extToolsCourseID, "course-id", 0, "Course ID")
	extToolsListCmd.Flags().Int64Var(&extToolsAccountID, "account-id", 0, "Account ID")
	extToolsListCmd.Flags().StringVar(&extToolsSearch, "search", "", "Search term")
	extToolsListCmd.Flags().BoolVar(&extToolsSelectable, "selectable", false, "Only show selectable tools")
	extToolsListCmd.Flags().BoolVar(&extToolsIncludeParent, "include-parents", false, "Include tools from parent contexts")

	// Get flags
	extToolsGetCmd.Flags().Int64Var(&extToolsCourseID, "course-id", 0, "Course ID")
	extToolsGetCmd.Flags().Int64Var(&extToolsAccountID, "account-id", 0, "Account ID")

	// Create flags
	extToolsCreateCmd.Flags().Int64Var(&extToolsCourseID, "course-id", 0, "Course ID")
	extToolsCreateCmd.Flags().Int64Var(&extToolsAccountID, "account-id", 0, "Account ID")
	extToolsCreateCmd.Flags().StringVar(&extToolsName, "name", "", "Tool name (required)")
	extToolsCreateCmd.Flags().StringVar(&extToolsURL, "url", "", "Tool URL")
	extToolsCreateCmd.Flags().StringVar(&extToolsDomain, "domain", "", "Tool domain")
	extToolsCreateCmd.Flags().StringVar(&extToolsConsumerKey, "consumer-key", "", "OAuth consumer key")
	extToolsCreateCmd.Flags().StringVar(&extToolsSharedSecret, "shared-secret", "", "OAuth shared secret")
	extToolsCreateCmd.Flags().StringVar(&extToolsPrivacyLevel, "privacy-level", "", "Privacy level: anonymous, name_only, email_only, public")
	extToolsCreateCmd.Flags().StringVar(&extToolsDescription, "description", "", "Tool description")
	extToolsCreateCmd.Flags().StringVar(&extToolsConfigType, "config-type", "", "Config type: url, xml")
	extToolsCreateCmd.Flags().StringVar(&extToolsConfigURL, "config-url", "", "Configuration URL")
	extToolsCreateCmd.Flags().StringVar(&extToolsConfigXML, "config-xml", "", "Configuration XML")
	extToolsCreateCmd.Flags().StringVar(&extToolsJSONFile, "json", "", "JSON file with full tool configuration")

	// Update flags
	extToolsUpdateCmd.Flags().Int64Var(&extToolsCourseID, "course-id", 0, "Course ID")
	extToolsUpdateCmd.Flags().Int64Var(&extToolsAccountID, "account-id", 0, "Account ID")
	extToolsUpdateCmd.Flags().StringVar(&extToolsName, "name", "", "Tool name")
	extToolsUpdateCmd.Flags().StringVar(&extToolsURL, "url", "", "Tool URL")
	extToolsUpdateCmd.Flags().StringVar(&extToolsDomain, "domain", "", "Tool domain")
	extToolsUpdateCmd.Flags().StringVar(&extToolsConsumerKey, "consumer-key", "", "OAuth consumer key")
	extToolsUpdateCmd.Flags().StringVar(&extToolsSharedSecret, "shared-secret", "", "OAuth shared secret")
	extToolsUpdateCmd.Flags().StringVar(&extToolsPrivacyLevel, "privacy-level", "", "Privacy level")
	extToolsUpdateCmd.Flags().StringVar(&extToolsDescription, "description", "", "Tool description")

	// Delete flags
	extToolsDeleteCmd.Flags().Int64Var(&extToolsCourseID, "course-id", 0, "Course ID")
	extToolsDeleteCmd.Flags().Int64Var(&extToolsAccountID, "account-id", 0, "Account ID")
	extToolsDeleteCmd.Flags().BoolVarP(&extToolsForce, "force", "f", false, "Skip confirmation prompt")

	// Launch flags
	extToolsLaunchCmd.Flags().Int64Var(&extToolsCourseID, "course-id", 0, "Course ID")
	extToolsLaunchCmd.Flags().Int64Var(&extToolsAccountID, "account-id", 0, "Account ID")
	extToolsLaunchCmd.Flags().StringVar(&extToolsLaunchType, "launch-type", "", "Launch type: course_navigation, account_navigation, assessment, module_item")
	extToolsLaunchCmd.Flags().Int64Var(&extToolsAssignmentID, "assignment-id", 0, "Assignment ID (for assessment launch)")
	extToolsLaunchCmd.Flags().Int64Var(&extToolsModuleItemID, "module-item-id", 0, "Module item ID (for module_item launch)")
}

func runExtToolsList(cmd *cobra.Command, args []string) error {
	// Use default account ID if neither course nor account is specified
	if extToolsCourseID == 0 && extToolsAccountID == 0 {
		defaultID, err := getDefaultAccountID()
		if err != nil || defaultID == 0 {
			return fmt.Errorf("must specify --course-id or --account-id (no default account configured). Use 'canvas config account --detect' to set one")
		}
		extToolsAccountID = defaultID
		printVerbose("Using default account ID: %d\n", defaultID)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewExternalToolsService(client)

	var selectable *bool
	if cmd.Flags().Changed("selectable") {
		selectable = &extToolsSelectable
	}

	opts := &api.ListExternalToolsOptions{
		Search:         extToolsSearch,
		Selectable:     selectable,
		IncludeParents: extToolsIncludeParent,
	}

	ctx := context.Background()
	var tools []api.ExternalTool

	if extToolsCourseID > 0 {
		tools, err = service.ListByCourse(ctx, extToolsCourseID, opts)
	} else {
		tools, err = service.ListByAccount(ctx, extToolsAccountID, opts)
	}

	if err != nil {
		return fmt.Errorf("failed to list external tools: %w", err)
	}

	if len(tools) == 0 {
		fmt.Println("No external tools found")
		return nil
	}

	printVerbose("Found %d external tools:\n\n", len(tools))
	return formatOutput(tools, nil)
}

func runExtToolsGet(cmd *cobra.Command, args []string) error {
	if extToolsCourseID == 0 && extToolsAccountID == 0 {
		return fmt.Errorf("either --course-id or --account-id is required")
	}

	toolID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid tool ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewExternalToolsService(client)

	ctx := context.Background()
	var tool *api.ExternalTool

	if extToolsCourseID > 0 {
		tool, err = service.GetByCourse(ctx, extToolsCourseID, toolID)
	} else {
		tool, err = service.GetByAccount(ctx, extToolsAccountID, toolID)
	}

	if err != nil {
		return fmt.Errorf("failed to get external tool: %w", err)
	}

	return formatOutput(tool, nil)
}

func runExtToolsCreate(cmd *cobra.Command, args []string) error {
	if extToolsCourseID == 0 && extToolsAccountID == 0 {
		return fmt.Errorf("either --course-id or --account-id is required")
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewExternalToolsService(client)

	var params *api.CreateExternalToolParams

	if extToolsJSONFile != "" {
		data, err := os.ReadFile(extToolsJSONFile)
		if err != nil {
			return fmt.Errorf("failed to read JSON file: %w", err)
		}
		params = &api.CreateExternalToolParams{}
		if err := json.Unmarshal(data, params); err != nil {
			return fmt.Errorf("failed to parse JSON file: %w", err)
		}
	} else {
		if extToolsName == "" {
			return fmt.Errorf("--name is required when not using --json")
		}

		params = &api.CreateExternalToolParams{
			Name:         extToolsName,
			URL:          extToolsURL,
			Domain:       extToolsDomain,
			ConsumerKey:  extToolsConsumerKey,
			SharedSecret: extToolsSharedSecret,
			PrivacyLevel: extToolsPrivacyLevel,
			Description:  extToolsDescription,
			ConfigType:   extToolsConfigType,
			ConfigURL:    extToolsConfigURL,
			ConfigXML:    extToolsConfigXML,
		}
	}

	ctx := context.Background()
	var tool *api.ExternalTool

	if extToolsCourseID > 0 {
		tool, err = service.CreateInCourse(ctx, extToolsCourseID, params)
	} else {
		tool, err = service.CreateInAccount(ctx, extToolsAccountID, params)
	}

	if err != nil {
		return fmt.Errorf("failed to create external tool: %w", err)
	}

	fmt.Printf("External tool created successfully (ID: %d)\n", tool.ID)
	return formatOutput(tool, nil)
}

func runExtToolsUpdate(cmd *cobra.Command, args []string) error {
	if extToolsCourseID == 0 && extToolsAccountID == 0 {
		return fmt.Errorf("either --course-id or --account-id is required")
	}

	toolID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid tool ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewExternalToolsService(client)

	params := &api.UpdateExternalToolParams{}

	if cmd.Flags().Changed("name") {
		params.Name = &extToolsName
	}
	if cmd.Flags().Changed("url") {
		params.URL = &extToolsURL
	}
	if cmd.Flags().Changed("domain") {
		params.Domain = &extToolsDomain
	}
	if cmd.Flags().Changed("consumer-key") {
		params.ConsumerKey = &extToolsConsumerKey
	}
	if cmd.Flags().Changed("shared-secret") {
		params.SharedSecret = &extToolsSharedSecret
	}
	if cmd.Flags().Changed("privacy-level") {
		params.PrivacyLevel = &extToolsPrivacyLevel
	}
	if cmd.Flags().Changed("description") {
		params.Description = &extToolsDescription
	}

	ctx := context.Background()
	var tool *api.ExternalTool

	if extToolsCourseID > 0 {
		tool, err = service.UpdateInCourse(ctx, extToolsCourseID, toolID, params)
	} else {
		tool, err = service.UpdateInAccount(ctx, extToolsAccountID, toolID, params)
	}

	if err != nil {
		return fmt.Errorf("failed to update external tool: %w", err)
	}

	fmt.Printf("External tool updated successfully (ID: %d)\n", tool.ID)
	return formatOutput(tool, nil)
}

func runExtToolsDelete(cmd *cobra.Command, args []string) error {
	if extToolsCourseID == 0 && extToolsAccountID == 0 {
		return fmt.Errorf("either --course-id or --account-id is required")
	}

	toolID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid tool ID: %w", err)
	}

	// Confirm deletion
	confirmed, err := confirmDelete("external tool", toolID, extToolsForce)
	if err != nil {
		return err
	}
	if !confirmed {
		fmt.Println("Cancelled.")
		return nil
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewExternalToolsService(client)

	ctx := context.Background()

	if extToolsCourseID > 0 {
		_, err = service.DeleteFromCourse(ctx, extToolsCourseID, toolID)
	} else {
		_, err = service.DeleteFromAccount(ctx, extToolsAccountID, toolID)
	}

	if err != nil {
		return fmt.Errorf("failed to delete external tool: %w", err)
	}

	fmt.Printf("External tool %d deleted successfully\n", toolID)
	return nil
}

func runExtToolsLaunch(cmd *cobra.Command, args []string) error {
	if extToolsCourseID == 0 && extToolsAccountID == 0 {
		return fmt.Errorf("either --course-id or --account-id is required")
	}

	toolID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid tool ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewExternalToolsService(client)

	params := &api.SessionlessLaunchParams{
		ID:           toolID,
		LaunchType:   extToolsLaunchType,
		AssignmentID: extToolsAssignmentID,
		ModuleItemID: extToolsModuleItemID,
	}

	ctx := context.Background()
	var result *api.SessionlessLaunchURL

	if extToolsCourseID > 0 {
		result, err = service.GetSessionlessLaunchURLForCourse(ctx, extToolsCourseID, params)
	} else {
		result, err = service.GetSessionlessLaunchURLForAccount(ctx, extToolsAccountID, params)
	}

	if err != nil {
		return fmt.Errorf("failed to get launch URL: %w", err)
	}

	fmt.Printf("Launch URL for %s:\n%s\n", result.Name, result.URL)
	return nil
}
