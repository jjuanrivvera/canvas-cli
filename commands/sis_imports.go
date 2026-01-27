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

var sisImportsCmd = &cobra.Command{
	Use:     "sis-imports",
	Aliases: []string{"sis"},
	Short:   "Manage SIS imports",
	Long: `Manage Canvas SIS (Student Information System) imports.

SIS imports allow you to bulk import data like users, courses,
sections, and enrollments from CSV files.

Examples:
  canvas sis-imports list --account-id 1
  canvas sis-imports get 123 --account-id 1
  canvas sis-imports create --account-id 1 --file data.zip`,
}

func init() {
	rootCmd.AddCommand(sisImportsCmd)
	sisImportsCmd.AddCommand(newSISListCmd())
	sisImportsCmd.AddCommand(newSISGetCmd())
	sisImportsCmd.AddCommand(newSISCreateCmd())
	sisImportsCmd.AddCommand(newSISAbortCmd())
	sisImportsCmd.AddCommand(newSISRestoreCmd())
	sisImportsCmd.AddCommand(newSISErrorsCmd())
}

func newSISListCmd() *cobra.Command {
	opts := &options.SISImportsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List SIS imports",
		Long: `List SIS imports for an account.

If --account-id is not specified, uses the default account ID from config.

Workflow states: initializing, created, importing, cleanup_batch,
                 imported, imported_with_messages, aborted,
                 failed, failed_with_messages, restoring, partially_restored

Examples:
  canvas sis-imports list                              # Uses default account
  canvas sis-imports list --account-id 1
  canvas sis-imports list --workflow-state imported`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runSISList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID (uses default if configured)")
	cmd.Flags().StringVar(&opts.WorkflowState, "workflow-state", "", "Filter by workflow state")
	cmd.Flags().StringVar(&opts.CreatedSince, "created-since", "", "Filter imports created since date (ISO8601)")
	cmd.Flags().StringVar(&opts.CreatedBefore, "created-before", "", "Filter imports created before date (ISO8601)")

	return cmd
}

func newSISGetCmd() *cobra.Command {
	opts := &options.SISImportsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <import-id>",
		Short: "Get a SIS import",
		Long: `Get details of a specific SIS import.

Examples:
  canvas sis-imports get 123 --account-id 1`,
		Args: ExactArgsWithUsage(1, "import-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			importID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid import ID: %w", err)
			}
			opts.ImportID = importID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runSISGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID (required)")
	cmd.MarkFlagRequired("account-id")

	return cmd
}

func newSISCreateCmd() *cobra.Command {
	opts := &options.SISImportsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a SIS import",
		Long: `Create a new SIS import from a CSV or ZIP file.

The file should follow the Canvas SIS import format:
- CSV files for individual data types
- ZIP file containing multiple CSVs

Examples:
  canvas sis-imports create --account-id 1 --file users.csv
  canvas sis-imports create --account-id 1 --file data.zip --batch-mode --batch-mode-term-id 123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Track which fields were set
			opts.BatchModeSet = cmd.Flags().Changed("batch-mode")
			opts.BatchModeTermIDSet = cmd.Flags().Changed("batch-mode-term-id")
			opts.OverrideStickinessSet = cmd.Flags().Changed("override-sis-stickiness")
			opts.AddStickinessSet = cmd.Flags().Changed("add-sis-stickiness")
			opts.ClearStickinessSet = cmd.Flags().Changed("clear-sis-stickiness")
			opts.DiffingRemasterSet = cmd.Flags().Changed("diffing-remaster-data-set")
			opts.ChangeThresholdSet = cmd.Flags().Changed("change-threshold")

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runSISCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID (required)")
	cmd.Flags().StringVar(&opts.FilePath, "file", "", "CSV or ZIP file to import (required)")
	cmd.Flags().StringVar(&opts.ImportType, "import-type", "instructure_csv", "Import type")
	cmd.Flags().StringVar(&opts.Extension, "extension", "", "File extension (csv, zip)")
	cmd.Flags().BoolVar(&opts.BatchMode, "batch-mode", false, "Enable batch mode")
	cmd.Flags().Int64Var(&opts.BatchModeTermID, "batch-mode-term-id", 0, "Term ID for batch mode")
	cmd.Flags().BoolVar(&opts.OverrideStickiness, "override-sis-stickiness", false, "Override SIS stickiness")
	cmd.Flags().BoolVar(&opts.AddStickiness, "add-sis-stickiness", false, "Add SIS stickiness")
	cmd.Flags().BoolVar(&opts.ClearStickiness, "clear-sis-stickiness", false, "Clear SIS stickiness")
	cmd.Flags().StringVar(&opts.DiffingID, "diffing-data-set-identifier", "", "Diffing dataset identifier")
	cmd.Flags().BoolVar(&opts.DiffingRemaster, "diffing-remaster-data-set", false, "Remaster diffing dataset")
	cmd.Flags().Float64Var(&opts.ChangeThreshold, "change-threshold", 0, "Skip if changes below threshold (0.0-1.0)")
	cmd.MarkFlagRequired("account-id")
	cmd.MarkFlagRequired("file")

	return cmd
}

func newSISAbortCmd() *cobra.Command {
	opts := &options.SISImportsAbortOptions{}

	cmd := &cobra.Command{
		Use:   "abort <import-id>",
		Short: "Abort a SIS import",
		Long: `Abort a pending or running SIS import.

Examples:
  canvas sis-imports abort 123 --account-id 1`,
		Args: ExactArgsWithUsage(1, "import-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			importID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid import ID: %w", err)
			}
			opts.ImportID = importID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runSISAbort(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID (required)")
	cmd.MarkFlagRequired("account-id")

	return cmd
}

func newSISRestoreCmd() *cobra.Command {
	opts := &options.SISImportsRestoreOptions{}

	cmd := &cobra.Command{
		Use:   "restore <import-id>",
		Short: "Restore workflow states",
		Long: `Restore workflow states for items from a previous import.

This will restore the workflow_state for all items that changed
their workflow_state during the import.

Examples:
  canvas sis-imports restore 123 --account-id 1`,
		Args: ExactArgsWithUsage(1, "import-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			importID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid import ID: %w", err)
			}
			opts.ImportID = importID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runSISRestore(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID (required)")
	cmd.Flags().BoolVar(&opts.BatchMode, "batch-mode", false, "Use batch mode for restore")
	cmd.MarkFlagRequired("account-id")

	return cmd
}

func newSISErrorsCmd() *cobra.Command {
	opts := &options.SISImportsErrorsOptions{}

	cmd := &cobra.Command{
		Use:   "errors <import-id>",
		Short: "List import errors",
		Long: `List errors from a SIS import.

Examples:
  canvas sis-imports errors 123 --account-id 1`,
		Args: ExactArgsWithUsage(1, "import-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			importID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid import ID: %w", err)
			}
			opts.ImportID = importID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runSISErrors(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID (required)")
	cmd.MarkFlagRequired("account-id")

	return cmd
}

func runSISList(ctx context.Context, client *api.Client, opts *options.SISImportsListOptions) error {
	logger := logging.NewCommandLogger(verbose)

	accountID, err := resolveAccountID(opts.AccountID, "sis-imports list")
	if err != nil {
		return err
	}

	logger.LogCommandStart(ctx, "sis_imports.list", map[string]interface{}{
		"account_id":     accountID,
		"workflow_state": opts.WorkflowState,
	})

	service := api.NewSISImportsService(client)

	apiOpts := &api.ListSISImportsOptions{
		WorkflowState: opts.WorkflowState,
		CreatedSince:  opts.CreatedSince,
		CreatedBefore: opts.CreatedBefore,
	}

	imports, err := service.List(ctx, accountID, apiOpts)
	if err != nil {
		logger.LogCommandError(ctx, "sis_imports.list", err, map[string]interface{}{
			"account_id": accountID,
		})
		return fmt.Errorf("failed to list SIS imports: %w", err)
	}

	if len(imports) == 0 {
		fmt.Println("No SIS imports found")
		logger.LogCommandComplete(ctx, "sis_imports.list", 0)
		return nil
	}

	printVerbose("Found %d SIS imports:\n\n", len(imports))
	logger.LogCommandComplete(ctx, "sis_imports.list", len(imports))
	return formatOutput(imports, nil)
}

func runSISGet(ctx context.Context, client *api.Client, opts *options.SISImportsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "sis_imports.get", map[string]interface{}{
		"account_id": opts.AccountID,
		"import_id":  opts.ImportID,
	})

	service := api.NewSISImportsService(client)

	sisImport, err := service.Get(ctx, opts.AccountID, opts.ImportID)
	if err != nil {
		logger.LogCommandError(ctx, "sis_imports.get", err, map[string]interface{}{
			"account_id": opts.AccountID,
			"import_id":  opts.ImportID,
		})
		return fmt.Errorf("failed to get SIS import: %w", err)
	}

	logger.LogCommandComplete(ctx, "sis_imports.get", 1)
	return formatOutput(sisImport, nil)
}

func runSISCreate(ctx context.Context, client *api.Client, opts *options.SISImportsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "sis_imports.create", map[string]interface{}{
		"account_id": opts.AccountID,
		"file_path":  opts.FilePath,
		"batch_mode": opts.BatchMode,
	})

	service := api.NewSISImportsService(client)

	params := &api.CreateSISImportParams{
		FilePath:                 opts.FilePath,
		ImportType:               opts.ImportType,
		Extension:                opts.Extension,
		DiffingDataSetIdentifier: opts.DiffingID,
	}

	if opts.BatchModeSet {
		params.BatchMode = &opts.BatchMode
	}

	if opts.BatchModeTermIDSet {
		params.BatchModeTermID = &opts.BatchModeTermID
	}

	if opts.OverrideStickinessSet {
		params.OverrideSISStickiness = &opts.OverrideStickiness
	}

	if opts.AddStickinessSet {
		params.AddSISStickiness = &opts.AddStickiness
	}

	if opts.ClearStickinessSet {
		params.ClearSISStickiness = &opts.ClearStickiness
	}

	if opts.DiffingRemasterSet {
		params.DiffingRemasterDataSet = &opts.DiffingRemaster
	}

	if opts.ChangeThresholdSet {
		params.ChangeThreshold = &opts.ChangeThreshold
	}

	sisImport, err := service.Create(ctx, opts.AccountID, params)
	if err != nil {
		logger.LogCommandError(ctx, "sis_imports.create", err, map[string]interface{}{
			"account_id": opts.AccountID,
			"file_path":  opts.FilePath,
		})
		return fmt.Errorf("failed to create SIS import: %w", err)
	}

	printInfo("SIS import created successfully (ID: %d)\n", sisImport.ID)
	fmt.Printf("Workflow state: %s\n", sisImport.WorkflowState)
	logger.LogCommandComplete(ctx, "sis_imports.create", 1)
	return nil
}

func runSISAbort(ctx context.Context, client *api.Client, opts *options.SISImportsAbortOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "sis_imports.abort", map[string]interface{}{
		"account_id": opts.AccountID,
		"import_id":  opts.ImportID,
	})

	service := api.NewSISImportsService(client)

	sisImport, err := service.Abort(ctx, opts.AccountID, opts.ImportID)
	if err != nil {
		logger.LogCommandError(ctx, "sis_imports.abort", err, map[string]interface{}{
			"account_id": opts.AccountID,
			"import_id":  opts.ImportID,
		})
		return fmt.Errorf("failed to abort SIS import: %w", err)
	}

	fmt.Printf("SIS import %d aborted\n", sisImport.ID)
	fmt.Printf("Workflow state: %s\n", sisImport.WorkflowState)
	logger.LogCommandComplete(ctx, "sis_imports.abort", 1)
	return nil
}

func runSISRestore(ctx context.Context, client *api.Client, opts *options.SISImportsRestoreOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "sis_imports.restore", map[string]interface{}{
		"account_id": opts.AccountID,
		"import_id":  opts.ImportID,
		"batch_mode": opts.BatchMode,
	})

	service := api.NewSISImportsService(client)

	progress, err := service.RestoreStates(ctx, opts.AccountID, opts.ImportID, opts.BatchMode)
	if err != nil {
		logger.LogCommandError(ctx, "sis_imports.restore", err, map[string]interface{}{
			"account_id": opts.AccountID,
			"import_id":  opts.ImportID,
		})
		return fmt.Errorf("failed to restore SIS import states: %w", err)
	}

	fmt.Printf("Restore initiated (Progress ID: %d)\n", progress.ID)
	fmt.Printf("Workflow state: %s\n", progress.WorkflowState)
	logger.LogCommandComplete(ctx, "sis_imports.restore", 1)
	return nil
}

func runSISErrors(ctx context.Context, client *api.Client, opts *options.SISImportsErrorsOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "sis_imports.errors", map[string]interface{}{
		"account_id": opts.AccountID,
		"import_id":  opts.ImportID,
	})

	service := api.NewSISImportsService(client)

	errors, err := service.ListErrors(ctx, opts.AccountID, opts.ImportID)
	if err != nil {
		logger.LogCommandError(ctx, "sis_imports.errors", err, map[string]interface{}{
			"account_id": opts.AccountID,
			"import_id":  opts.ImportID,
		})
		return fmt.Errorf("failed to get SIS import errors: %w", err)
	}

	if len(errors) == 0 {
		fmt.Println("No errors found for this import")
		logger.LogCommandComplete(ctx, "sis_imports.errors", 0)
		return nil
	}

	printVerbose("Found %d errors:\n\n", len(errors))
	logger.LogCommandComplete(ctx, "sis_imports.errors", len(errors))
	return formatOutput(errors, nil)
}
