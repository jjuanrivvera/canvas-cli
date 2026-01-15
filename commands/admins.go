package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var (
	adminAccountID    int64
	adminUserID       int64
	adminRole         string
	adminRoleID       int64
	adminSendConfirm  bool
	roleAccountID     int64
	roleLabel         string
	roleBaseType      string
	roleState         string
	roleShowInherited bool
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

// adminsListCmd lists admins for an account
var adminsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List administrators for an account",
	Long: `Retrieve a list of all administrators for the specified account.

If --account-id is not specified, uses the default account ID from config.
Set a default with: canvas config account --detect

Examples:
  canvas admins list                 # Uses default account
  canvas admins list --account-id 1  # Explicit account`,
	RunE: runAdminsList,
}

// adminsAddCmd adds an admin to an account
var adminsAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an administrator to an account",
	Long: `Add a user as an administrator to the specified account.

Examples:
  canvas admins add --account-id 1 --user-id 123
  canvas admins add --account-id 1 --user-id 123 --role AccountAdmin`,
	RunE: runAdminsAdd,
}

// adminsRemoveCmd removes an admin from an account
var adminsRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove an administrator from an account",
	Long: `Remove a user's administrator privileges from the specified account.

Examples:
  canvas admins remove --account-id 1 --user-id 123
  canvas admins remove --account-id 1 --user-id 123 --role-id 456`,
	RunE: runAdminsRemove,
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

// rolesListCmd lists roles for an account
var rolesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List roles for an account",
	Long: `Retrieve a list of all roles for the specified account.

If --account-id is not specified, uses the default account ID from config.

Examples:
  canvas roles list                              # Uses default account
  canvas roles list --account-id 1
  canvas roles list --account-id 1 --state active
  canvas roles list --show-inherited`,
	RunE: runRolesList,
}

// rolesGetCmd gets a single role
var rolesGetCmd = &cobra.Command{
	Use:   "get <role-id>",
	Short: "Get details for a specific role",
	Long: `Retrieve detailed information about a specific role.

Examples:
  canvas roles get 123 --account-id 1`,
	Args: ExactArgsWithUsage(1, "role-id"),
	RunE: runRolesGet,
}

// rolesCreateCmd creates a new role
var rolesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new role",
	Long: `Create a new custom role in the specified account.

Examples:
  canvas roles create --account-id 1 --label "Custom Teacher" --base-type TeacherEnrollment`,
	RunE: runRolesCreate,
}

// rolesUpdateCmd updates a role
var rolesUpdateCmd = &cobra.Command{
	Use:   "update <role-id>",
	Short: "Update an existing role",
	Long: `Update the properties of an existing role.

Examples:
  canvas roles update 123 --account-id 1 --label "New Label"`,
	Args: ExactArgsWithUsage(1, "role-id"),
	RunE: runRolesUpdate,
}

// rolesDeactivateCmd deactivates a role
var rolesDeactivateCmd = &cobra.Command{
	Use:   "deactivate <role-id>",
	Short: "Deactivate a role",
	Long: `Deactivate an existing role. Deactivated roles cannot be assigned to new users.

Examples:
  canvas roles deactivate 123 --account-id 1`,
	Args: ExactArgsWithUsage(1, "role-id"),
	RunE: runRolesDeactivate,
}

// rolesActivateCmd activates a role
var rolesActivateCmd = &cobra.Command{
	Use:   "activate <role-id>",
	Short: "Activate a role",
	Long: `Reactivate a previously deactivated role.

Examples:
  canvas roles activate 123 --account-id 1`,
	Args: ExactArgsWithUsage(1, "role-id"),
	RunE: runRolesActivate,
}

func init() {
	rootCmd.AddCommand(adminsCmd)
	rootCmd.AddCommand(rolesCmd)

	// Admins subcommands
	adminsCmd.AddCommand(adminsListCmd)
	adminsCmd.AddCommand(adminsAddCmd)
	adminsCmd.AddCommand(adminsRemoveCmd)

	// Roles subcommands
	rolesCmd.AddCommand(rolesListCmd)
	rolesCmd.AddCommand(rolesGetCmd)
	rolesCmd.AddCommand(rolesCreateCmd)
	rolesCmd.AddCommand(rolesUpdateCmd)
	rolesCmd.AddCommand(rolesDeactivateCmd)
	rolesCmd.AddCommand(rolesActivateCmd)

	// Admins list flags
	adminsListCmd.Flags().Int64Var(&adminAccountID, "account-id", 0, "Account ID (uses default if configured)")

	// Admins add flags
	adminsAddCmd.Flags().Int64Var(&adminAccountID, "account-id", 0, "Account ID (required)")
	adminsAddCmd.Flags().Int64Var(&adminUserID, "user-id", 0, "User ID to add as admin (required)")
	adminsAddCmd.Flags().StringVar(&adminRole, "role", "", "Role name (e.g., AccountAdmin)")
	adminsAddCmd.Flags().Int64Var(&adminRoleID, "role-id", 0, "Role ID")
	adminsAddCmd.Flags().BoolVar(&adminSendConfirm, "send-confirmation", false, "Send confirmation email")
	adminsAddCmd.MarkFlagRequired("account-id")
	adminsAddCmd.MarkFlagRequired("user-id")

	// Admins remove flags
	adminsRemoveCmd.Flags().Int64Var(&adminAccountID, "account-id", 0, "Account ID (required)")
	adminsRemoveCmd.Flags().Int64Var(&adminUserID, "user-id", 0, "User ID to remove (required)")
	adminsRemoveCmd.Flags().Int64Var(&adminRoleID, "role-id", 0, "Role ID (optional, removes specific role)")
	adminsRemoveCmd.MarkFlagRequired("account-id")
	adminsRemoveCmd.MarkFlagRequired("user-id")

	// Roles list flags
	rolesListCmd.Flags().Int64Var(&roleAccountID, "account-id", 0, "Account ID (uses default if configured)")
	rolesListCmd.Flags().StringVar(&roleState, "state", "", "Filter by state (active, inactive)")
	rolesListCmd.Flags().BoolVar(&roleShowInherited, "show-inherited", false, "Show inherited roles")

	// Roles get flags
	rolesGetCmd.Flags().Int64Var(&roleAccountID, "account-id", 0, "Account ID (required)")
	rolesGetCmd.MarkFlagRequired("account-id")

	// Roles create flags
	rolesCreateCmd.Flags().Int64Var(&roleAccountID, "account-id", 0, "Account ID (required)")
	rolesCreateCmd.Flags().StringVar(&roleLabel, "label", "", "Role label (required)")
	rolesCreateCmd.Flags().StringVar(&roleBaseType, "base-type", "", "Base role type (AccountMembership, StudentEnrollment, TeacherEnrollment, TaEnrollment, ObserverEnrollment, DesignerEnrollment)")
	rolesCreateCmd.MarkFlagRequired("account-id")
	rolesCreateCmd.MarkFlagRequired("label")

	// Roles update flags
	rolesUpdateCmd.Flags().Int64Var(&roleAccountID, "account-id", 0, "Account ID (required)")
	rolesUpdateCmd.Flags().StringVar(&roleLabel, "label", "", "New role label")
	rolesUpdateCmd.MarkFlagRequired("account-id")

	// Roles deactivate flags
	rolesDeactivateCmd.Flags().Int64Var(&roleAccountID, "account-id", 0, "Account ID (required)")
	rolesDeactivateCmd.MarkFlagRequired("account-id")

	// Roles activate flags
	rolesActivateCmd.Flags().Int64Var(&roleAccountID, "account-id", 0, "Account ID (required)")
	rolesActivateCmd.MarkFlagRequired("account-id")
}

func runAdminsList(cmd *cobra.Command, args []string) error {
	accountID, err := resolveAccountID(adminAccountID, "admins list")
	if err != nil {
		return err
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewAdminsService(client)

	ctx := context.Background()
	admins, err := service.List(ctx, accountID, nil)
	if err != nil {
		return fmt.Errorf("failed to list admins: %w", err)
	}

	if len(admins) == 0 {
		fmt.Println("No admins found")
		return nil
	}

	printVerbose("Found %d admins:\n\n", len(admins))
	return formatOutput(admins, nil)
}

func runAdminsAdd(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	params := &api.CreateAdminParams{
		UserID: adminUserID,
		Role:   adminRole,
		RoleID: adminRoleID,
	}

	if cmd.Flags().Changed("send-confirmation") {
		params.SendConfirmation = &adminSendConfirm
	}

	service := api.NewAdminsService(client)

	ctx := context.Background()
	admin, err := service.Create(ctx, adminAccountID, params)
	if err != nil {
		return fmt.Errorf("failed to add admin: %w", err)
	}

	fmt.Printf("Admin added successfully (User ID: %d)\n", admin.UserID)
	return nil
}

func runAdminsRemove(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	var roleIDPtr *int64
	if cmd.Flags().Changed("role-id") {
		roleIDPtr = &adminRoleID
	}

	service := api.NewAdminsService(client)

	ctx := context.Background()
	admin, err := service.Delete(ctx, adminAccountID, adminUserID, roleIDPtr)
	if err != nil {
		return fmt.Errorf("failed to remove admin: %w", err)
	}

	fmt.Printf("Admin removed successfully (User ID: %d)\n", admin.UserID)
	return nil
}

func runRolesList(cmd *cobra.Command, args []string) error {
	accountID, err := resolveAccountID(roleAccountID, "roles list")
	if err != nil {
		return err
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	opts := &api.ListRolesOptions{
		State:         roleState,
		ShowInherited: roleShowInherited,
	}

	service := api.NewRolesService(client)

	ctx := context.Background()
	roles, err := service.List(ctx, accountID, opts)
	if err != nil {
		return fmt.Errorf("failed to list roles: %w", err)
	}

	if len(roles) == 0 {
		fmt.Println("No roles found")
		return nil
	}

	printVerbose("Found %d roles:\n\n", len(roles))
	return formatOutput(roles, nil)
}

func runRolesGet(cmd *cobra.Command, args []string) error {
	roleID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid role ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewRolesService(client)

	ctx := context.Background()
	role, err := service.Get(ctx, roleAccountID, roleID)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}

	return formatOutput(role, nil)
}

func runRolesCreate(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	params := &api.CreateRoleParams{
		Label:        roleLabel,
		BaseRoleType: roleBaseType,
	}

	service := api.NewRolesService(client)

	ctx := context.Background()
	role, err := service.Create(ctx, roleAccountID, params)
	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	fmt.Printf("Role created successfully (ID: %d)\n", role.ID)
	return nil
}

func runRolesUpdate(cmd *cobra.Command, args []string) error {
	roleID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid role ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	params := &api.UpdateRoleParams{}
	if cmd.Flags().Changed("label") {
		params.Label = &roleLabel
	}

	service := api.NewRolesService(client)

	ctx := context.Background()
	role, err := service.Update(ctx, roleAccountID, roleID, params)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	fmt.Printf("Role updated successfully (ID: %d)\n", role.ID)
	return nil
}

func runRolesDeactivate(cmd *cobra.Command, args []string) error {
	roleID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid role ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewRolesService(client)

	ctx := context.Background()
	role, err := service.Deactivate(ctx, roleAccountID, roleID)
	if err != nil {
		return fmt.Errorf("failed to deactivate role: %w", err)
	}

	fmt.Printf("Role deactivated (ID: %d, state: %s)\n", role.ID, role.WorkflowState)
	return nil
}

func runRolesActivate(cmd *cobra.Command, args []string) error {
	roleID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid role ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewRolesService(client)

	ctx := context.Background()
	role, err := service.Activate(ctx, roleAccountID, roleID)
	if err != nil {
		return fmt.Errorf("failed to activate role: %w", err)
	}

	fmt.Printf("Role activated (ID: %d, state: %s)\n", role.ID, role.WorkflowState)
	return nil
}
