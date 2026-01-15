package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var (
	sisAccountID     int64
	sisWorkflowState string
	sisCreatedSince  string
	sisCreatedBefore string

	// Create flags
	sisFilePath           string
	sisImportType         string
	sisExtension          string
	sisBatchMode          bool
	sisBatchModeTermID    int64
	sisOverrideStickiness bool
	sisAddStickiness      bool
	sisClearStickiness    bool
	sisDiffingID          string
	sisDiffingRemaster    bool
	sisChangeThreshold    float64
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

var sisListCmd = &cobra.Command{
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
	RunE: runSISList,
}

var sisGetCmd = &cobra.Command{
	Use:   "get <import-id>",
	Short: "Get a SIS import",
	Long: `Get details of a specific SIS import.

Examples:
  canvas sis-imports get 123 --account-id 1`,
	Args: ExactArgsWithUsage(1, "import-id"),
	RunE: runSISGet,
}

var sisCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a SIS import",
	Long: `Create a new SIS import from a CSV or ZIP file.

The file should follow the Canvas SIS import format:
- CSV files for individual data types
- ZIP file containing multiple CSVs

Examples:
  canvas sis-imports create --account-id 1 --file users.csv
  canvas sis-imports create --account-id 1 --file data.zip --batch-mode --batch-mode-term-id 123`,
	RunE: runSISCreate,
}

var sisAbortCmd = &cobra.Command{
	Use:   "abort <import-id>",
	Short: "Abort a SIS import",
	Long: `Abort a pending or running SIS import.

Examples:
  canvas sis-imports abort 123 --account-id 1`,
	Args: ExactArgsWithUsage(1, "import-id"),
	RunE: runSISAbort,
}

var sisRestoreCmd = &cobra.Command{
	Use:   "restore <import-id>",
	Short: "Restore workflow states",
	Long: `Restore workflow states for items from a previous import.

This will restore the workflow_state for all items that changed
their workflow_state during the import.

Examples:
  canvas sis-imports restore 123 --account-id 1`,
	Args: ExactArgsWithUsage(1, "import-id"),
	RunE: runSISRestore,
}

var sisErrorsCmd = &cobra.Command{
	Use:   "errors <import-id>",
	Short: "List import errors",
	Long: `List errors from a SIS import.

Examples:
  canvas sis-imports errors 123 --account-id 1`,
	Args: ExactArgsWithUsage(1, "import-id"),
	RunE: runSISErrors,
}

func init() {
	rootCmd.AddCommand(sisImportsCmd)
	sisImportsCmd.AddCommand(sisListCmd)
	sisImportsCmd.AddCommand(sisGetCmd)
	sisImportsCmd.AddCommand(sisCreateCmd)
	sisImportsCmd.AddCommand(sisAbortCmd)
	sisImportsCmd.AddCommand(sisRestoreCmd)
	sisImportsCmd.AddCommand(sisErrorsCmd)

	// List flags
	sisListCmd.Flags().Int64Var(&sisAccountID, "account-id", 0, "Account ID (uses default if configured)")
	sisListCmd.Flags().StringVar(&sisWorkflowState, "workflow-state", "", "Filter by workflow state")
	sisListCmd.Flags().StringVar(&sisCreatedSince, "created-since", "", "Filter imports created since date (ISO8601)")
	sisListCmd.Flags().StringVar(&sisCreatedBefore, "created-before", "", "Filter imports created before date (ISO8601)")

	// Get flags
	sisGetCmd.Flags().Int64Var(&sisAccountID, "account-id", 0, "Account ID (required)")
	sisGetCmd.MarkFlagRequired("account-id")

	// Create flags
	sisCreateCmd.Flags().Int64Var(&sisAccountID, "account-id", 0, "Account ID (required)")
	sisCreateCmd.Flags().StringVar(&sisFilePath, "file", "", "CSV or ZIP file to import (required)")
	sisCreateCmd.Flags().StringVar(&sisImportType, "import-type", "instructure_csv", "Import type")
	sisCreateCmd.Flags().StringVar(&sisExtension, "extension", "", "File extension (csv, zip)")
	sisCreateCmd.Flags().BoolVar(&sisBatchMode, "batch-mode", false, "Enable batch mode")
	sisCreateCmd.Flags().Int64Var(&sisBatchModeTermID, "batch-mode-term-id", 0, "Term ID for batch mode")
	sisCreateCmd.Flags().BoolVar(&sisOverrideStickiness, "override-sis-stickiness", false, "Override SIS stickiness")
	sisCreateCmd.Flags().BoolVar(&sisAddStickiness, "add-sis-stickiness", false, "Add SIS stickiness")
	sisCreateCmd.Flags().BoolVar(&sisClearStickiness, "clear-sis-stickiness", false, "Clear SIS stickiness")
	sisCreateCmd.Flags().StringVar(&sisDiffingID, "diffing-data-set-identifier", "", "Diffing dataset identifier")
	sisCreateCmd.Flags().BoolVar(&sisDiffingRemaster, "diffing-remaster-data-set", false, "Remaster diffing dataset")
	sisCreateCmd.Flags().Float64Var(&sisChangeThreshold, "change-threshold", 0, "Skip if changes below threshold (0.0-1.0)")
	sisCreateCmd.MarkFlagRequired("account-id")
	sisCreateCmd.MarkFlagRequired("file")

	// Abort flags
	sisAbortCmd.Flags().Int64Var(&sisAccountID, "account-id", 0, "Account ID (required)")
	sisAbortCmd.MarkFlagRequired("account-id")

	// Restore flags
	sisRestoreCmd.Flags().Int64Var(&sisAccountID, "account-id", 0, "Account ID (required)")
	sisRestoreCmd.Flags().BoolVar(&sisBatchMode, "batch-mode", false, "Use batch mode for restore")
	sisRestoreCmd.MarkFlagRequired("account-id")

	// Errors flags
	sisErrorsCmd.Flags().Int64Var(&sisAccountID, "account-id", 0, "Account ID (required)")
	sisErrorsCmd.MarkFlagRequired("account-id")
}

func runSISList(cmd *cobra.Command, args []string) error {
	accountID, err := resolveAccountID(sisAccountID, "sis-imports list")
	if err != nil {
		return err
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewSISImportsService(client)

	opts := &api.ListSISImportsOptions{
		WorkflowState: sisWorkflowState,
		CreatedSince:  sisCreatedSince,
		CreatedBefore: sisCreatedBefore,
	}

	ctx := context.Background()
	imports, err := service.List(ctx, accountID, opts)
	if err != nil {
		return fmt.Errorf("failed to list SIS imports: %w", err)
	}

	if len(imports) == 0 {
		fmt.Println("No SIS imports found")
		return nil
	}

	printVerbose("Found %d SIS imports:\n\n", len(imports))
	return formatOutput(imports, nil)
}

func runSISGet(cmd *cobra.Command, args []string) error {
	importID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid import ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewSISImportsService(client)

	ctx := context.Background()
	sisImport, err := service.Get(ctx, sisAccountID, importID)
	if err != nil {
		return fmt.Errorf("failed to get SIS import: %w", err)
	}

	return formatOutput(sisImport, nil)
}

func runSISCreate(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewSISImportsService(client)

	params := &api.CreateSISImportParams{
		FilePath:                 sisFilePath,
		ImportType:               sisImportType,
		Extension:                sisExtension,
		DiffingDataSetIdentifier: sisDiffingID,
	}

	if cmd.Flags().Changed("batch-mode") {
		params.BatchMode = &sisBatchMode
	}

	if cmd.Flags().Changed("batch-mode-term-id") {
		params.BatchModeTermID = &sisBatchModeTermID
	}

	if cmd.Flags().Changed("override-sis-stickiness") {
		params.OverrideSISStickiness = &sisOverrideStickiness
	}

	if cmd.Flags().Changed("add-sis-stickiness") {
		params.AddSISStickiness = &sisAddStickiness
	}

	if cmd.Flags().Changed("clear-sis-stickiness") {
		params.ClearSISStickiness = &sisClearStickiness
	}

	if cmd.Flags().Changed("diffing-remaster-data-set") {
		params.DiffingRemasterDataSet = &sisDiffingRemaster
	}

	if cmd.Flags().Changed("change-threshold") {
		params.ChangeThreshold = &sisChangeThreshold
	}

	ctx := context.Background()
	sisImport, err := service.Create(ctx, sisAccountID, params)
	if err != nil {
		return fmt.Errorf("failed to create SIS import: %w", err)
	}

	fmt.Printf("SIS import created successfully (ID: %d)\n", sisImport.ID)
	fmt.Printf("Workflow state: %s\n", sisImport.WorkflowState)
	return nil
}

func runSISAbort(cmd *cobra.Command, args []string) error {
	importID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid import ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewSISImportsService(client)

	ctx := context.Background()
	sisImport, err := service.Abort(ctx, sisAccountID, importID)
	if err != nil {
		return fmt.Errorf("failed to abort SIS import: %w", err)
	}

	fmt.Printf("SIS import %d aborted\n", sisImport.ID)
	fmt.Printf("Workflow state: %s\n", sisImport.WorkflowState)
	return nil
}

func runSISRestore(cmd *cobra.Command, args []string) error {
	importID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid import ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewSISImportsService(client)

	ctx := context.Background()
	progress, err := service.RestoreStates(ctx, sisAccountID, importID, sisBatchMode)
	if err != nil {
		return fmt.Errorf("failed to restore SIS import states: %w", err)
	}

	fmt.Printf("Restore initiated (Progress ID: %d)\n", progress.ID)
	fmt.Printf("Workflow state: %s\n", progress.WorkflowState)
	return nil
}

func runSISErrors(cmd *cobra.Command, args []string) error {
	importID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid import ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewSISImportsService(client)

	ctx := context.Background()
	errors, err := service.ListErrors(ctx, sisAccountID, importID)
	if err != nil {
		return fmt.Errorf("failed to get SIS import errors: %w", err)
	}

	if len(errors) == 0 {
		fmt.Println("No errors found for this import")
		return nil
	}

	printVerbose("Found %d errors:\n\n", len(errors))
	return formatOutput(errors, nil)
}
