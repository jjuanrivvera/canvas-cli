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

// adminsCmd represents the admins command
var adminsCmd = &cobra.Command{
	Use:   "admins",
	Short: "Manage account administrators",
	Long: `Manage Canvas account administrators.

Administrators have elevated privileges within an account and can manage
users, courses, and other account settings.

Examples:
  canvas admins list --account-id 1
  canvas admins add --account-id 1 --user-id 123
  canvas admins remove --account-id 1 --user-id 123`,
}

// rolesCmd represents the roles command
var rolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "Manage account roles",
	Long: `Manage Canvas account roles.

Roles define sets of permissions that can be assigned to users.
Canvas has built-in roles and allows custom roles.

Base role types:
  - AccountMembership
  - StudentEnrollment
  - TeacherEnrollment
  - TaEnrollment
  - ObserverEnrollment
  - DesignerEnrollment

Examples:
  canvas roles list --account-id 1
  canvas roles get 123 --account-id 1
  canvas roles create --account-id 1 --label "Custom Teacher" --base-type TeacherEnrollment`,
}

func init() {
	rootCmd.AddCommand(adminsCmd)
	rootCmd.AddCommand(rolesCmd)

	// Admins subcommands
	adminsCmd.AddCommand(newAdminsListCmd())
	adminsCmd.AddCommand(newAdminsAddCmd())
	adminsCmd.AddCommand(newAdminsRemoveCmd())

	// Roles subcommands
	rolesCmd.AddCommand(newRolesListCmd())
	rolesCmd.AddCommand(newRolesGetCmd())
	rolesCmd.AddCommand(newRolesCreateCmd())
	rolesCmd.AddCommand(newRolesUpdateCmd())
	rolesCmd.AddCommand(newRolesDeactivateCmd())
	rolesCmd.AddCommand(newRolesActivateCmd())
}

func newAdminsListCmd() *cobra.Command {
	opts := &options.AdminsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List administrators for an account",
		Long: `Retrieve a list of all administrators for the specified account.

If --account-id is not specified, uses the default account ID from config.
Set a default with: canvas config account --detect

Examples:
  canvas admins list                 # Uses default account
  canvas admins list --account-id 1  # Explicit account`,
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := resolveAccountID(opts.AccountID, "admins list")
			if err != nil {
				return err
			}
			opts.AccountID = accountID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAdminsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID (uses default if configured)")

	return cmd
}

func newAdminsAddCmd() *cobra.Command {
	opts := &options.AdminsAddOptions{}

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add an administrator to an account",
		Long: `Add a user as an administrator to the specified account.

Examples:
  canvas admins add --account-id 1 --user-id 123
  canvas admins add --account-id 1 --user-id 123 --role AccountAdmin`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.SendConfirmationSet = cmd.Flags().Changed("send-confirmation")

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAdminsAdd(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID (required)")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID to add as admin (required)")
	cmd.Flags().StringVar(&opts.Role, "role", "", "Role name (e.g., AccountAdmin)")
	cmd.Flags().Int64Var(&opts.RoleID, "role-id", 0, "Role ID")
	cmd.Flags().BoolVar(&opts.SendConfirmation, "send-confirmation", false, "Send confirmation email")
	cmd.MarkFlagRequired("account-id")
	cmd.MarkFlagRequired("user-id")

	return cmd
}

func newAdminsRemoveCmd() *cobra.Command {
	opts := &options.AdminsRemoveOptions{}

	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove an administrator from an account",
		Long: `Remove a user's administrator privileges from the specified account.

Examples:
  canvas admins remove --account-id 1 --user-id 123
  canvas admins remove --account-id 1 --user-id 123 --role-id 456`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.RoleIDSet = cmd.Flags().Changed("role-id")

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAdminsRemove(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID (required)")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID to remove (required)")
	cmd.Flags().Int64Var(&opts.RoleID, "role-id", 0, "Role ID (optional, removes specific role)")
	cmd.MarkFlagRequired("account-id")
	cmd.MarkFlagRequired("user-id")

	return cmd
}

func newRolesListCmd() *cobra.Command {
	opts := &options.RolesListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List roles for an account",
		Long: `Retrieve a list of all roles for the specified account.

If --account-id is not specified, uses the default account ID from config.

Examples:
  canvas roles list                              # Uses default account
  canvas roles list --account-id 1
  canvas roles list --account-id 1 --state active
  canvas roles list --show-inherited`,
		RunE: func(cmd *cobra.Command, args []string) error {
			accountID, err := resolveAccountID(opts.AccountID, "roles list")
			if err != nil {
				return err
			}
			opts.AccountID = accountID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runRolesList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID (uses default if configured)")
	cmd.Flags().StringVar(&opts.State, "state", "", "Filter by state (active, inactive)")
	cmd.Flags().BoolVar(&opts.ShowInherited, "show-inherited", false, "Show inherited roles")

	return cmd
}

func newRolesGetCmd() *cobra.Command {
	opts := &options.RolesGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <role-id>",
		Short: "Get details for a specific role",
		Long: `Retrieve detailed information about a specific role.

Examples:
  canvas roles get 123 --account-id 1`,
		Args: ExactArgsWithUsage(1, "role-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			roleID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid role ID: %w", err)
			}
			opts.RoleID = roleID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runRolesGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID (required)")
	cmd.MarkFlagRequired("account-id")

	return cmd
}

func newRolesCreateCmd() *cobra.Command {
	opts := &options.RolesCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new role",
		Long: `Create a new custom role in the specified account.

Examples:
  canvas roles create --account-id 1 --label "Custom Teacher" --base-type TeacherEnrollment`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runRolesCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID (required)")
	cmd.Flags().StringVar(&opts.Label, "label", "", "Role label (required)")
	cmd.Flags().StringVar(&opts.BaseRoleType, "base-type", "", "Base role type (AccountMembership, StudentEnrollment, TeacherEnrollment, TaEnrollment, ObserverEnrollment, DesignerEnrollment)")
	cmd.MarkFlagRequired("account-id")
	cmd.MarkFlagRequired("label")

	return cmd
}

func newRolesUpdateCmd() *cobra.Command {
	opts := &options.RolesUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <role-id>",
		Short: "Update an existing role",
		Long: `Update the properties of an existing role.

Examples:
  canvas roles update 123 --account-id 1 --label "New Label"`,
		Args: ExactArgsWithUsage(1, "role-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			roleID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid role ID: %w", err)
			}
			opts.RoleID = roleID

			opts.LabelSet = cmd.Flags().Changed("label")

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runRolesUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID (required)")
	cmd.Flags().StringVar(&opts.Label, "label", "", "New role label")
	cmd.MarkFlagRequired("account-id")

	return cmd
}

func newRolesDeactivateCmd() *cobra.Command {
	opts := &options.RolesDeactivateOptions{}

	cmd := &cobra.Command{
		Use:   "deactivate <role-id>",
		Short: "Deactivate a role",
		Long: `Deactivate an existing role. Deactivated roles cannot be assigned to new users.

Examples:
  canvas roles deactivate 123 --account-id 1`,
		Args: ExactArgsWithUsage(1, "role-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			roleID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid role ID: %w", err)
			}
			opts.RoleID = roleID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runRolesDeactivate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID (required)")
	cmd.MarkFlagRequired("account-id")

	return cmd
}

func newRolesActivateCmd() *cobra.Command {
	opts := &options.RolesActivateOptions{}

	cmd := &cobra.Command{
		Use:   "activate <role-id>",
		Short: "Activate a role",
		Long: `Reactivate a previously deactivated role.

Examples:
  canvas roles activate 123 --account-id 1`,
		Args: ExactArgsWithUsage(1, "role-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			roleID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid role ID: %w", err)
			}
			opts.RoleID = roleID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runRolesActivate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID (required)")
	cmd.MarkFlagRequired("account-id")

	return cmd
}

func runAdminsList(ctx context.Context, client *api.Client, opts *options.AdminsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "admins.list", map[string]interface{}{
		"account_id": opts.AccountID,
	})

	service := api.NewAdminsService(client)

	admins, err := service.List(ctx, opts.AccountID, nil)
	if err != nil {
		logger.LogCommandError(ctx, "admins.list", err, map[string]interface{}{
			"account_id": opts.AccountID,
		})
		return fmt.Errorf("failed to list admins: %w", err)
	}

	if len(admins) == 0 {
		fmt.Println("No admins found")
		logger.LogCommandComplete(ctx, "admins.list", 0)
		return nil
	}

	printVerbose("Found %d admins:\n\n", len(admins))
	logger.LogCommandComplete(ctx, "admins.list", len(admins))
	return formatOutput(admins, nil)
}

func runAdminsAdd(ctx context.Context, client *api.Client, opts *options.AdminsAddOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "admins.add", map[string]interface{}{
		"account_id": opts.AccountID,
		"user_id":    opts.UserID,
		"role":       opts.Role,
		"role_id":    opts.RoleID,
	})

	params := &api.CreateAdminParams{
		UserID: opts.UserID,
		Role:   opts.Role,
		RoleID: opts.RoleID,
	}

	if opts.SendConfirmationSet {
		params.SendConfirmation = &opts.SendConfirmation
	}

	service := api.NewAdminsService(client)

	admin, err := service.Create(ctx, opts.AccountID, params)
	if err != nil {
		logger.LogCommandError(ctx, "admins.add", err, map[string]interface{}{
			"account_id": opts.AccountID,
			"user_id":    opts.UserID,
		})
		return fmt.Errorf("failed to add admin: %w", err)
	}

	fmt.Printf("Admin added successfully (User ID: %d)\n", admin.UserID)
	logger.LogCommandComplete(ctx, "admins.add", 1)
	return nil
}

func runAdminsRemove(ctx context.Context, client *api.Client, opts *options.AdminsRemoveOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "admins.remove", map[string]interface{}{
		"account_id": opts.AccountID,
		"user_id":    opts.UserID,
		"role_id":    opts.RoleID,
	})

	var roleIDPtr *int64
	if opts.RoleIDSet {
		roleIDPtr = &opts.RoleID
	}

	service := api.NewAdminsService(client)

	admin, err := service.Delete(ctx, opts.AccountID, opts.UserID, roleIDPtr)
	if err != nil {
		logger.LogCommandError(ctx, "admins.remove", err, map[string]interface{}{
			"account_id": opts.AccountID,
			"user_id":    opts.UserID,
		})
		return fmt.Errorf("failed to remove admin: %w", err)
	}

	fmt.Printf("Admin removed successfully (User ID: %d)\n", admin.UserID)
	logger.LogCommandComplete(ctx, "admins.remove", 1)
	return nil
}

func runRolesList(ctx context.Context, client *api.Client, opts *options.RolesListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "roles.list", map[string]interface{}{
		"account_id":     opts.AccountID,
		"state":          opts.State,
		"show_inherited": opts.ShowInherited,
	})

	apiOpts := &api.ListRolesOptions{
		State:         opts.State,
		ShowInherited: opts.ShowInherited,
	}

	service := api.NewRolesService(client)

	roles, err := service.List(ctx, opts.AccountID, apiOpts)
	if err != nil {
		logger.LogCommandError(ctx, "roles.list", err, map[string]interface{}{
			"account_id": opts.AccountID,
		})
		return fmt.Errorf("failed to list roles: %w", err)
	}

	if len(roles) == 0 {
		fmt.Println("No roles found")
		logger.LogCommandComplete(ctx, "roles.list", 0)
		return nil
	}

	printVerbose("Found %d roles:\n\n", len(roles))
	logger.LogCommandComplete(ctx, "roles.list", len(roles))
	return formatOutput(roles, nil)
}

func runRolesGet(ctx context.Context, client *api.Client, opts *options.RolesGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "roles.get", map[string]interface{}{
		"account_id": opts.AccountID,
		"role_id":    opts.RoleID,
	})

	service := api.NewRolesService(client)

	role, err := service.Get(ctx, opts.AccountID, opts.RoleID)
	if err != nil {
		logger.LogCommandError(ctx, "roles.get", err, map[string]interface{}{
			"account_id": opts.AccountID,
			"role_id":    opts.RoleID,
		})
		return fmt.Errorf("failed to get role: %w", err)
	}

	logger.LogCommandComplete(ctx, "roles.get", 1)
	return formatOutput(role, nil)
}

func runRolesCreate(ctx context.Context, client *api.Client, opts *options.RolesCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "roles.create", map[string]interface{}{
		"account_id":     opts.AccountID,
		"label":          opts.Label,
		"base_role_type": opts.BaseRoleType,
	})

	params := &api.CreateRoleParams{
		Label:        opts.Label,
		BaseRoleType: opts.BaseRoleType,
	}

	service := api.NewRolesService(client)

	role, err := service.Create(ctx, opts.AccountID, params)
	if err != nil {
		logger.LogCommandError(ctx, "roles.create", err, map[string]interface{}{
			"account_id": opts.AccountID,
			"label":      opts.Label,
		})
		return fmt.Errorf("failed to create role: %w", err)
	}

	fmt.Printf("Role created successfully (ID: %d)\n", role.ID)
	logger.LogCommandComplete(ctx, "roles.create", 1)
	return nil
}

func runRolesUpdate(ctx context.Context, client *api.Client, opts *options.RolesUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "roles.update", map[string]interface{}{
		"account_id": opts.AccountID,
		"role_id":    opts.RoleID,
		"label":      opts.Label,
	})

	params := &api.UpdateRoleParams{}
	if opts.LabelSet {
		params.Label = &opts.Label
	}

	service := api.NewRolesService(client)

	role, err := service.Update(ctx, opts.AccountID, opts.RoleID, params)
	if err != nil {
		logger.LogCommandError(ctx, "roles.update", err, map[string]interface{}{
			"account_id": opts.AccountID,
			"role_id":    opts.RoleID,
		})
		return fmt.Errorf("failed to update role: %w", err)
	}

	fmt.Printf("Role updated successfully (ID: %d)\n", role.ID)
	logger.LogCommandComplete(ctx, "roles.update", 1)
	return nil
}

func runRolesDeactivate(ctx context.Context, client *api.Client, opts *options.RolesDeactivateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "roles.deactivate", map[string]interface{}{
		"account_id": opts.AccountID,
		"role_id":    opts.RoleID,
	})

	service := api.NewRolesService(client)

	role, err := service.Deactivate(ctx, opts.AccountID, opts.RoleID)
	if err != nil {
		logger.LogCommandError(ctx, "roles.deactivate", err, map[string]interface{}{
			"account_id": opts.AccountID,
			"role_id":    opts.RoleID,
		})
		return fmt.Errorf("failed to deactivate role: %w", err)
	}

	fmt.Printf("Role deactivated (ID: %d, state: %s)\n", role.ID, role.WorkflowState)
	logger.LogCommandComplete(ctx, "roles.deactivate", 1)
	return nil
}

func runRolesActivate(ctx context.Context, client *api.Client, opts *options.RolesActivateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "roles.activate", map[string]interface{}{
		"account_id": opts.AccountID,
		"role_id":    opts.RoleID,
	})

	service := api.NewRolesService(client)

	role, err := service.Activate(ctx, opts.AccountID, opts.RoleID)
	if err != nil {
		logger.LogCommandError(ctx, "roles.activate", err, map[string]interface{}{
			"account_id": opts.AccountID,
			"role_id":    opts.RoleID,
		})
		return fmt.Errorf("failed to activate role: %w", err)
	}

	fmt.Printf("Role activated (ID: %d, state: %s)\n", role.ID, role.WorkflowState)
	logger.LogCommandComplete(ctx, "roles.activate", 1)
	return nil
}
