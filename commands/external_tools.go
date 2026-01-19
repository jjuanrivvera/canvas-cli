package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
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

func init() {
	rootCmd.AddCommand(externalToolsCmd)
	externalToolsCmd.AddCommand(newExtToolsListCmd())
	externalToolsCmd.AddCommand(newExtToolsGetCmd())
	externalToolsCmd.AddCommand(newExtToolsCreateCmd())
	externalToolsCmd.AddCommand(newExtToolsUpdateCmd())
	externalToolsCmd.AddCommand(newExtToolsDeleteCmd())
	externalToolsCmd.AddCommand(newExtToolsLaunchCmd())
}

func newExtToolsListCmd() *cobra.Command {
	opts := &options.ExternalToolsListOptions{}

	cmd := &cobra.Command{
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
		RunE: func(cmd *cobra.Command, args []string) error {
			// Use default account ID if neither course nor account is specified
			if opts.CourseID == 0 && opts.AccountID == 0 {
				defaultID, err := getDefaultAccountID()
				if err != nil || defaultID == 0 {
					return fmt.Errorf("must specify --course-id or --account-id (no default account configured). Use 'canvas config account --detect' to set one")
				}
				opts.AccountID = defaultID
				printVerbose("Using default account ID: %d\n", defaultID)
			}

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runExtToolsList(cmd.Context(), client, cmd, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID")
	cmd.Flags().StringVar(&opts.Search, "search", "", "Search term")
	cmd.Flags().BoolVar(&opts.Selectable, "selectable", false, "Only show selectable tools")
	cmd.Flags().BoolVar(&opts.IncludeParent, "include-parents", false, "Include tools from parent contexts")

	return cmd
}

func newExtToolsGetCmd() *cobra.Command {
	opts := &options.ExternalToolsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <tool-id>",
		Short: "Get an external tool",
		Long: `Get details of a specific external tool.

Examples:
  canvas external-tools get 456 --course-id 123
  canvas external-tools get 456 --account-id 1`,
		Args: ExactArgsWithUsage(1, "tool-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			toolID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid tool ID: %w", err)
			}
			opts.ToolID = toolID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runExtToolsGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID")

	return cmd
}

func newExtToolsCreateCmd() *cobra.Command {
	opts := &options.ExternalToolsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an external tool",
		Long: `Create a new external tool in a course or account.

Examples:
  canvas external-tools create --course-id 123 --name "My Tool" --url https://tool.example.com
  canvas external-tools create --course-id 123 --name "LTI Tool" --consumer-key key123 --shared-secret secret123
  canvas external-tools create --course-id 123 --json tool-config.json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runExtToolsCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Tool name (required)")
	cmd.Flags().StringVar(&opts.URL, "url", "", "Tool URL")
	cmd.Flags().StringVar(&opts.Domain, "domain", "", "Tool domain")
	cmd.Flags().StringVar(&opts.ConsumerKey, "consumer-key", "", "OAuth consumer key")
	cmd.Flags().StringVar(&opts.SharedSecret, "shared-secret", "", "OAuth shared secret")
	cmd.Flags().StringVar(&opts.PrivacyLevel, "privacy-level", "", "Privacy level: anonymous, name_only, email_only, public")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Tool description")
	cmd.Flags().StringVar(&opts.ConfigType, "config-type", "", "Config type: url, xml")
	cmd.Flags().StringVar(&opts.ConfigURL, "config-url", "", "Configuration URL")
	cmd.Flags().StringVar(&opts.ConfigXML, "config-xml", "", "Configuration XML")
	cmd.Flags().StringVar(&opts.JSONFile, "json", "", "JSON file with full tool configuration")

	return cmd
}

func newExtToolsUpdateCmd() *cobra.Command {
	opts := &options.ExternalToolsUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <tool-id>",
		Short: "Update an external tool",
		Long: `Update an existing external tool.

Examples:
  canvas external-tools update 456 --course-id 123 --name "Updated Name"
  canvas external-tools update 456 --course-id 123 --url https://new-url.example.com`,
		Args: ExactArgsWithUsage(1, "tool-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			toolID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid tool ID: %w", err)
			}
			opts.ToolID = toolID

			// Track which fields were set
			opts.NameSet = cmd.Flags().Changed("name")
			opts.URLSet = cmd.Flags().Changed("url")
			opts.DomainSet = cmd.Flags().Changed("domain")
			opts.ConsumerKeySet = cmd.Flags().Changed("consumer-key")
			opts.SharedSecretSet = cmd.Flags().Changed("shared-secret")
			opts.PrivacyLevelSet = cmd.Flags().Changed("privacy-level")
			opts.DescriptionSet = cmd.Flags().Changed("description")

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runExtToolsUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Tool name")
	cmd.Flags().StringVar(&opts.URL, "url", "", "Tool URL")
	cmd.Flags().StringVar(&opts.Domain, "domain", "", "Tool domain")
	cmd.Flags().StringVar(&opts.ConsumerKey, "consumer-key", "", "OAuth consumer key")
	cmd.Flags().StringVar(&opts.SharedSecret, "shared-secret", "", "OAuth shared secret")
	cmd.Flags().StringVar(&opts.PrivacyLevel, "privacy-level", "", "Privacy level")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Tool description")

	return cmd
}

func newExtToolsDeleteCmd() *cobra.Command {
	opts := &options.ExternalToolsDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <tool-id>",
		Short: "Delete an external tool",
		Long: `Delete an external tool from a course or account.

Examples:
  canvas external-tools delete 456 --course-id 123
  canvas external-tools delete 456 --account-id 1`,
		Args: ExactArgsWithUsage(1, "tool-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			toolID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid tool ID: %w", err)
			}
			opts.ToolID = toolID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runExtToolsDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID")
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Skip confirmation prompt")

	return cmd
}

func newExtToolsLaunchCmd() *cobra.Command {
	opts := &options.ExternalToolsLaunchOptions{}

	cmd := &cobra.Command{
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
		RunE: func(cmd *cobra.Command, args []string) error {
			toolID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid tool ID: %w", err)
			}
			opts.ToolID = toolID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runExtToolsLaunch(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.Flags().StringVar(&opts.LaunchType, "launch-type", "", "Launch type: course_navigation, account_navigation, assessment, module_item")
	cmd.Flags().Int64Var(&opts.AssignmentID, "assignment-id", 0, "Assignment ID (for assessment launch)")
	cmd.Flags().Int64Var(&opts.ModuleItemID, "module-item-id", 0, "Module item ID (for module_item launch)")

	return cmd
}

func runExtToolsList(ctx context.Context, client *api.Client, cmd *cobra.Command, opts *options.ExternalToolsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "external_tools.list", map[string]interface{}{
		"course_id":      opts.CourseID,
		"account_id":     opts.AccountID,
		"search":         opts.Search,
		"selectable":     opts.Selectable,
		"include_parent": opts.IncludeParent,
	})

	service := api.NewExternalToolsService(client)

	var selectable *bool
	if cmd.Flags().Changed("selectable") {
		selectable = &opts.Selectable
	}

	apiOpts := &api.ListExternalToolsOptions{
		Search:         opts.Search,
		Selectable:     selectable,
		IncludeParents: opts.IncludeParent,
	}

	var tools []api.ExternalTool
	var err error

	if opts.CourseID > 0 {
		tools, err = service.ListByCourse(ctx, opts.CourseID, apiOpts)
	} else {
		tools, err = service.ListByAccount(ctx, opts.AccountID, apiOpts)
	}

	if err != nil {
		logger.LogCommandError(ctx, "external_tools.list", err, map[string]interface{}{
			"course_id":  opts.CourseID,
			"account_id": opts.AccountID,
		})
		return fmt.Errorf("failed to list external tools: %w", err)
	}

	if len(tools) == 0 {
		fmt.Println("No external tools found")
		logger.LogCommandComplete(ctx, "external_tools.list", 0)
		return nil
	}

	printVerbose("Found %d external tools:\n\n", len(tools))
	logger.LogCommandComplete(ctx, "external_tools.list", len(tools))
	return formatOutput(tools, nil)
}

func runExtToolsGet(ctx context.Context, client *api.Client, opts *options.ExternalToolsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "external_tools.get", map[string]interface{}{
		"course_id":  opts.CourseID,
		"account_id": opts.AccountID,
		"tool_id":    opts.ToolID,
	})

	service := api.NewExternalToolsService(client)

	var tool *api.ExternalTool
	var err error

	if opts.CourseID > 0 {
		tool, err = service.GetByCourse(ctx, opts.CourseID, opts.ToolID)
	} else {
		tool, err = service.GetByAccount(ctx, opts.AccountID, opts.ToolID)
	}

	if err != nil {
		logger.LogCommandError(ctx, "external_tools.get", err, map[string]interface{}{
			"course_id":  opts.CourseID,
			"account_id": opts.AccountID,
			"tool_id":    opts.ToolID,
		})
		return fmt.Errorf("failed to get external tool: %w", err)
	}

	logger.LogCommandComplete(ctx, "external_tools.get", 1)
	return formatOutput(tool, nil)
}

func runExtToolsCreate(ctx context.Context, client *api.Client, opts *options.ExternalToolsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "external_tools.create", map[string]interface{}{
		"course_id":  opts.CourseID,
		"account_id": opts.AccountID,
		"name":       opts.Name,
		"json_file":  opts.JSONFile,
	})

	service := api.NewExternalToolsService(client)

	var params *api.CreateExternalToolParams

	if opts.JSONFile != "" {
		data, err := os.ReadFile(opts.JSONFile)
		if err != nil {
			logger.LogCommandError(ctx, "external_tools.create", err, map[string]interface{}{
				"json_file": opts.JSONFile,
			})
			return fmt.Errorf("failed to read JSON file: %w", err)
		}
		params = &api.CreateExternalToolParams{}
		if err := json.Unmarshal(data, params); err != nil {
			logger.LogCommandError(ctx, "external_tools.create", err, map[string]interface{}{
				"json_file": opts.JSONFile,
			})
			return fmt.Errorf("failed to parse JSON file: %w", err)
		}
	} else {
		if opts.Name == "" {
			return fmt.Errorf("--name is required when not using --json")
		}

		params = &api.CreateExternalToolParams{
			Name:         opts.Name,
			URL:          opts.URL,
			Domain:       opts.Domain,
			ConsumerKey:  opts.ConsumerKey,
			SharedSecret: opts.SharedSecret,
			PrivacyLevel: opts.PrivacyLevel,
			Description:  opts.Description,
			ConfigType:   opts.ConfigType,
			ConfigURL:    opts.ConfigURL,
			ConfigXML:    opts.ConfigXML,
		}
	}

	var tool *api.ExternalTool
	var err error

	if opts.CourseID > 0 {
		tool, err = service.CreateInCourse(ctx, opts.CourseID, params)
	} else {
		tool, err = service.CreateInAccount(ctx, opts.AccountID, params)
	}

	if err != nil {
		logger.LogCommandError(ctx, "external_tools.create", err, map[string]interface{}{
			"course_id":  opts.CourseID,
			"account_id": opts.AccountID,
		})
		return fmt.Errorf("failed to create external tool: %w", err)
	}

	fmt.Printf("External tool created successfully (ID: %d)\n", tool.ID)
	logger.LogCommandComplete(ctx, "external_tools.create", 1)
	return formatOutput(tool, nil)
}

func runExtToolsUpdate(ctx context.Context, client *api.Client, opts *options.ExternalToolsUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "external_tools.update", map[string]interface{}{
		"course_id":  opts.CourseID,
		"account_id": opts.AccountID,
		"tool_id":    opts.ToolID,
	})

	service := api.NewExternalToolsService(client)

	params := &api.UpdateExternalToolParams{}

	if opts.NameSet {
		params.Name = &opts.Name
	}
	if opts.URLSet {
		params.URL = &opts.URL
	}
	if opts.DomainSet {
		params.Domain = &opts.Domain
	}
	if opts.ConsumerKeySet {
		params.ConsumerKey = &opts.ConsumerKey
	}
	if opts.SharedSecretSet {
		params.SharedSecret = &opts.SharedSecret
	}
	if opts.PrivacyLevelSet {
		params.PrivacyLevel = &opts.PrivacyLevel
	}
	if opts.DescriptionSet {
		params.Description = &opts.Description
	}

	var tool *api.ExternalTool
	var err error

	if opts.CourseID > 0 {
		tool, err = service.UpdateInCourse(ctx, opts.CourseID, opts.ToolID, params)
	} else {
		tool, err = service.UpdateInAccount(ctx, opts.AccountID, opts.ToolID, params)
	}

	if err != nil {
		logger.LogCommandError(ctx, "external_tools.update", err, map[string]interface{}{
			"course_id":  opts.CourseID,
			"account_id": opts.AccountID,
			"tool_id":    opts.ToolID,
		})
		return fmt.Errorf("failed to update external tool: %w", err)
	}

	fmt.Printf("External tool updated successfully (ID: %d)\n", tool.ID)
	logger.LogCommandComplete(ctx, "external_tools.update", 1)
	return formatOutput(tool, nil)
}

func runExtToolsDelete(ctx context.Context, client *api.Client, opts *options.ExternalToolsDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "external_tools.delete", map[string]interface{}{
		"course_id":  opts.CourseID,
		"account_id": opts.AccountID,
		"tool_id":    opts.ToolID,
		"force":      opts.Force,
	})

	// Confirm deletion
	confirmed, err := confirmDelete("external tool", opts.ToolID, opts.Force)
	if err != nil {
		logger.LogCommandError(ctx, "external_tools.delete", err, map[string]interface{}{
			"tool_id": opts.ToolID,
		})
		return err
	}
	if !confirmed {
		fmt.Println("Cancelled.")
		return nil
	}

	service := api.NewExternalToolsService(client)

	if opts.CourseID > 0 {
		_, err = service.DeleteFromCourse(ctx, opts.CourseID, opts.ToolID)
	} else {
		_, err = service.DeleteFromAccount(ctx, opts.AccountID, opts.ToolID)
	}

	if err != nil {
		logger.LogCommandError(ctx, "external_tools.delete", err, map[string]interface{}{
			"course_id":  opts.CourseID,
			"account_id": opts.AccountID,
			"tool_id":    opts.ToolID,
		})
		return fmt.Errorf("failed to delete external tool: %w", err)
	}

	fmt.Printf("External tool %d deleted successfully\n", opts.ToolID)
	logger.LogCommandComplete(ctx, "external_tools.delete", 1)
	return nil
}

func runExtToolsLaunch(ctx context.Context, client *api.Client, opts *options.ExternalToolsLaunchOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "external_tools.launch", map[string]interface{}{
		"course_id":      opts.CourseID,
		"tool_id":        opts.ToolID,
		"launch_type":    opts.LaunchType,
		"assignment_id":  opts.AssignmentID,
		"module_item_id": opts.ModuleItemID,
	})

	service := api.NewExternalToolsService(client)

	params := &api.SessionlessLaunchParams{
		ID:           opts.ToolID,
		LaunchType:   opts.LaunchType,
		AssignmentID: opts.AssignmentID,
		ModuleItemID: opts.ModuleItemID,
	}

	result, err := service.GetSessionlessLaunchURLForCourse(ctx, opts.CourseID, params)
	if err != nil {
		logger.LogCommandError(ctx, "external_tools.launch", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"tool_id":   opts.ToolID,
		})
		return fmt.Errorf("failed to get launch URL: %w", err)
	}

	fmt.Printf("Launch URL for %s:\n%s\n", result.Name, result.URL)
	logger.LogCommandComplete(ctx, "external_tools.launch", 1)
	return nil
}
