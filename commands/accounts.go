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

func init() {
	rootCmd.AddCommand(accountsCmd)
	accountsCmd.AddCommand(newAccountsListCmd())
	accountsCmd.AddCommand(newAccountsGetCmd())
	accountsCmd.AddCommand(newAccountsSubAccountsCmd())
}

func newAccountsListCmd() *cobra.Command {
	opts := &options.AccountsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List accessible accounts",
		Long: `List all accounts the current user has access to.

This typically returns accounts where you have admin or sub-admin permissions.

Examples:
  canvas accounts list
  canvas accounts list --include lti_guid`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAccountsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Additional data to include (lti_guid, registration_settings, services)")

	return cmd
}

func newAccountsGetCmd() *cobra.Command {
	opts := &options.AccountsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <account-id>",
		Short: "Get details of a specific account",
		Long: `Get details of a specific account by ID.

Examples:
  canvas accounts get 1
  canvas accounts get 5`,
		Args: ExactArgsWithUsage(1, "account-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			opts.AccountID = accountID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAccountsGet(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newAccountsSubAccountsCmd() *cobra.Command {
	opts := &options.AccountsSubAccountsOptions{}

	cmd := &cobra.Command{
		Use:   "sub <account-id>",
		Short: "List sub-accounts of an account",
		Long: `List sub-accounts for a given parent account.

Use --recursive to get the entire account tree.

Examples:
  canvas accounts sub 1
  canvas accounts sub 1 --recursive`,
		Args: ExactArgsWithUsage(1, "account-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid account ID: %s", args[0])
			}
			opts.AccountID = accountID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAccountsSubAccounts(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Recursive, "recursive", false, "List entire account tree recursively")

	return cmd
}

func runAccountsList(ctx context.Context, client *api.Client, opts *options.AccountsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "accounts.list", map[string]interface{}{
		"include": opts.Include,
	})

	service := api.NewAccountsService(client)

	var apiOpts *api.ListAccountsOptions
	if len(opts.Include) > 0 {
		apiOpts = &api.ListAccountsOptions{
			Include: opts.Include,
		}
	}

	accounts, err := service.List(ctx, apiOpts)
	if err != nil {
		logger.LogCommandError(ctx, "accounts.list", err, nil)
		return err
	}

	if len(accounts) == 0 {
		fmt.Println("No accounts found. You may not have admin access to any accounts.")
		logger.LogCommandComplete(ctx, "accounts.list", 0)
		return nil
	}

	logger.LogCommandComplete(ctx, "accounts.list", len(accounts))

	// Format and display accounts
	return formatOutput(accounts, func() {
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
	})
}

func runAccountsGet(ctx context.Context, client *api.Client, opts *options.AccountsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "accounts.get", map[string]interface{}{
		"account_id": opts.AccountID,
	})

	service := api.NewAccountsService(client)

	account, err := service.Get(ctx, opts.AccountID)
	if err != nil {
		logger.LogCommandError(ctx, "accounts.get", err, map[string]interface{}{
			"account_id": opts.AccountID,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "accounts.get", 1)

	// Format and display account
	return formatOutput(account, func() {
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
	})
}

func runAccountsSubAccounts(ctx context.Context, client *api.Client, opts *options.AccountsSubAccountsOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "accounts.sub", map[string]interface{}{
		"account_id": opts.AccountID,
		"recursive":  opts.Recursive,
	})

	service := api.NewAccountsService(client)

	apiOpts := &api.ListSubAccountsOptions{
		Recursive: opts.Recursive,
	}

	accounts, err := service.ListSubAccounts(ctx, opts.AccountID, apiOpts)
	if err != nil {
		logger.LogCommandError(ctx, "accounts.sub", err, map[string]interface{}{
			"account_id": opts.AccountID,
		})
		return err
	}

	if len(accounts) == 0 {
		fmt.Printf("No sub-accounts found for account %d\n", opts.AccountID)
		logger.LogCommandComplete(ctx, "accounts.sub", 0)
		return nil
	}

	logger.LogCommandComplete(ctx, "accounts.sub", len(accounts))

	// Format and display accounts
	return formatOutput(accounts, func() {
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
	})
}
