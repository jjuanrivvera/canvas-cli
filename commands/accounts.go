package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var (
	accountsInclude   []string
	accountsRecursive bool
)

// accountsCmd represents the accounts command group
var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Manage Canvas accounts",
	Long: `Manage Canvas accounts including listing accessible accounts and viewing account details.

Accounts represent institutions, sub-accounts, or organizational units within Canvas.
Most admin operations require account-level permissions.

Examples:
  canvas accounts list
  canvas accounts get 1
  canvas accounts sub 1 --recursive`,
}

// accountsListCmd represents the accounts list command
var accountsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List accessible accounts",
	Long: `List all accounts the current user has access to.

This typically returns accounts where you have admin or sub-admin permissions.

Examples:
  canvas accounts list
  canvas accounts list --include lti_guid`,
	RunE: runAccountsList,
}

// accountsGetCmd represents the accounts get command
var accountsGetCmd = &cobra.Command{
	Use:   "get <account-id>",
	Short: "Get details of a specific account",
	Long: `Get details of a specific account by ID.

Examples:
  canvas accounts get 1
  canvas accounts get 5`,
	Args: cobra.ExactArgs(1),
	RunE: runAccountsGet,
}

// accountsSubAccountsCmd represents the accounts sub command
var accountsSubAccountsCmd = &cobra.Command{
	Use:   "sub <account-id>",
	Short: "List sub-accounts of an account",
	Long: `List sub-accounts for a given parent account.

Use --recursive to get the entire account tree.

Examples:
  canvas accounts sub 1
  canvas accounts sub 1 --recursive`,
	Args: cobra.ExactArgs(1),
	RunE: runAccountsSubAccounts,
}

func init() {
	rootCmd.AddCommand(accountsCmd)
	accountsCmd.AddCommand(accountsListCmd)
	accountsCmd.AddCommand(accountsGetCmd)
	accountsCmd.AddCommand(accountsSubAccountsCmd)

	// List flags
	accountsListCmd.Flags().StringSliceVar(&accountsInclude, "include", []string{}, "Additional data to include (lti_guid, registration_settings, services)")

	// Sub-accounts flags
	accountsSubAccountsCmd.Flags().BoolVar(&accountsRecursive, "recursive", false, "List entire account tree recursively")
}

func runAccountsList(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	service := api.NewAccountsService(client)

	var opts *api.ListAccountsOptions
	if len(accountsInclude) > 0 {
		opts = &api.ListAccountsOptions{
			Include: accountsInclude,
		}
	}

	accounts, err := service.List(ctx, opts)
	if err != nil {
		return err
	}

	if len(accounts) == 0 {
		fmt.Println("No accounts found. You may not have admin access to any accounts.")
		return nil
	}

	// Display accounts
	if outputFormat == "json" {
		data, err := json.MarshalIndent(accounts, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Table output
	fmt.Printf("%-8s %-40s %-12s %-15s\n", "ID", "NAME", "PARENT_ID", "STATE")
	fmt.Println(strings.Repeat("-", 80))

	for _, account := range accounts {
		parentID := "-"
		if account.ParentAccountID != 0 {
			parentID = strconv.FormatInt(account.ParentAccountID, 10)
		}

		name := account.Name
		if len(name) > 38 {
			name = name[:35] + "..."
		}

		fmt.Printf("%-8d %-40s %-12s %-15s\n",
			account.ID,
			name,
			parentID,
			account.WorkflowState,
		)
	}

	fmt.Printf("\nTotal: %d account(s)\n", len(accounts))

	// Helpful tip for new users
	if len(accounts) > 0 {
		fmt.Println("\nTip: Use 'canvas courses list --account <id>' to see all courses in an account")
	}

	return nil
}

func runAccountsGet(cmd *cobra.Command, args []string) error {
	accountID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid account ID: %s", args[0])
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	service := api.NewAccountsService(client)

	account, err := service.Get(ctx, accountID)
	if err != nil {
		return err
	}

	// Display account
	if outputFormat == "json" {
		data, err := json.MarshalIndent(account, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Detailed output
	fmt.Printf("Account Details\n")
	fmt.Printf("===============\n\n")
	fmt.Printf("ID:                 %d\n", account.ID)
	fmt.Printf("Name:               %s\n", account.Name)
	if account.UUID != "" {
		fmt.Printf("UUID:               %s\n", account.UUID)
	}
	if account.ParentAccountID != 0 {
		fmt.Printf("Parent Account ID:  %d\n", account.ParentAccountID)
	}
	if account.RootAccountID != 0 {
		fmt.Printf("Root Account ID:    %d\n", account.RootAccountID)
	}
	fmt.Printf("Workflow State:     %s\n", account.WorkflowState)
	if account.DefaultTimeZone != "" {
		fmt.Printf("Time Zone:          %s\n", account.DefaultTimeZone)
	}
	if account.SISAccountID != "" {
		fmt.Printf("SIS Account ID:     %s\n", account.SISAccountID)
	}

	return nil
}

func runAccountsSubAccounts(cmd *cobra.Command, args []string) error {
	accountID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid account ID: %s", args[0])
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	service := api.NewAccountsService(client)

	opts := &api.ListSubAccountsOptions{
		Recursive: accountsRecursive,
	}

	accounts, err := service.ListSubAccounts(ctx, accountID, opts)
	if err != nil {
		return err
	}

	if len(accounts) == 0 {
		fmt.Printf("No sub-accounts found for account %d\n", accountID)
		return nil
	}

	// Display accounts
	if outputFormat == "json" {
		data, err := json.MarshalIndent(accounts, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Table output
	fmt.Printf("%-8s %-40s %-12s %-15s\n", "ID", "NAME", "PARENT_ID", "STATE")
	fmt.Println(strings.Repeat("-", 80))

	for _, account := range accounts {
		name := account.Name
		if len(name) > 38 {
			name = name[:35] + "..."
		}

		fmt.Printf("%-8d %-40s %-12d %-15s\n",
			account.ID,
			name,
			account.ParentAccountID,
			account.WorkflowState,
		)
	}

	fmt.Printf("\nTotal: %d sub-account(s)\n", len(accounts))

	return nil
}
