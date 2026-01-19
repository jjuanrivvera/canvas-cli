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

// groupsCmd represents the groups command group
var groupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "Manage Canvas groups",
	Long: `Manage Canvas groups and group categories.

Groups allow students and instructors to collaborate on projects and activities.
Groups can be organized into categories with different self-signup options.

Examples:
  canvas groups list --course-id 123
  canvas groups get 456
  canvas groups categories list --course-id 123`,
}

// newGroupsListCmd creates the groups list command
func newGroupsListCmd() *cobra.Command {
	opts := &options.GroupsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List groups",
		Long: `List groups for a course, account, or user.

Examples:
  canvas groups list --course-id 123
  canvas groups list --account-id 1
  canvas groups list --user-id 456
  canvas groups list  # Lists current user's groups`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGroupsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID (0 for self)")
	cmd.Flags().BoolVar(&opts.IncludeUsers, "include-users", false, "Include group users")
	cmd.Flags().BoolVar(&opts.IncludePermissions, "include-permissions", false, "Include permissions")

	return cmd
}

// newGroupsGetCmd creates the groups get command
func newGroupsGetCmd() *cobra.Command {
	opts := &options.GroupsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <group-id>",
		Short: "Get group details",
		Long: `Get details of a specific group.

Examples:
  canvas groups get 456
  canvas groups get 456 --include-users`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGroupsGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().BoolVar(&opts.IncludeUsers, "include-users", false, "Include group users")
	cmd.Flags().BoolVar(&opts.IncludePermissions, "include-permissions", false, "Include permissions")

	return cmd
}

// newGroupsCreateCmd creates the groups create command
func newGroupsCreateCmd() *cobra.Command {
	opts := &options.GroupsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new group",
		Long: `Create a new group in a group category.

Examples:
  canvas groups create --category-id 123 --name "Study Group"
  canvas groups create --category-id 123 --name "Project Team" --description "Our project team"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGroupsCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CategoryID, "category-id", 0, "Group category ID (required)")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Group name (required)")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Group description")
	cmd.Flags().BoolVar(&opts.IsPublic, "public", false, "Whether the group is public")
	cmd.Flags().StringVar(&opts.JoinLevel, "join-level", "", "Join level (parent_context_auto_join, parent_context_request, invitation_only)")
	cmd.Flags().Int64Var(&opts.StorageQuotaMb, "storage-quota-mb", 0, "Storage quota in MB")
	cmd.Flags().StringVar(&opts.SISGroupID, "sis-group-id", "", "SIS group ID")
	cmd.MarkFlagRequired("category-id")
	cmd.MarkFlagRequired("name")

	return cmd
}

// newGroupsUpdateCmd creates the groups update command
func newGroupsUpdateCmd() *cobra.Command {
	opts := &options.GroupsUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <group-id>",
		Short: "Update a group",
		Long: `Update an existing group.

Examples:
  canvas groups update 456 --name "New Name"
  canvas groups update 456 --description "Updated description"`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID

			// Track which fields were set
			opts.NameSet = cmd.Flags().Changed("name")
			opts.DescriptionSet = cmd.Flags().Changed("description")
			opts.IsPublicSet = cmd.Flags().Changed("public")
			opts.JoinLevelSet = cmd.Flags().Changed("join-level")
			opts.StorageQuotaMbSet = cmd.Flags().Changed("storage-quota-mb")
			opts.SISGroupIDSet = cmd.Flags().Changed("sis-group-id")

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGroupsUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "Group name")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Group description")
	cmd.Flags().BoolVar(&opts.IsPublic, "public", false, "Whether the group is public")
	cmd.Flags().StringVar(&opts.JoinLevel, "join-level", "", "Join level")
	cmd.Flags().Int64Var(&opts.StorageQuotaMb, "storage-quota-mb", 0, "Storage quota in MB")
	cmd.Flags().StringVar(&opts.SISGroupID, "sis-group-id", "", "SIS group ID")

	return cmd
}

// newGroupsDeleteCmd creates the groups delete command
func newGroupsDeleteCmd() *cobra.Command {
	opts := &options.GroupsDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <group-id>",
		Short: "Delete a group",
		Long: `Delete an existing group.

Examples:
  canvas groups delete 456
  canvas groups delete 456 --force`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGroupsDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Force, "force", false, "Skip confirmation prompt")

	return cmd
}

// groupsMembersCmd represents the members subcommand
var groupsMembersCmd = &cobra.Command{
	Use:   "members",
	Short: "Manage group members",
	Long:  "Commands for managing group memberships.",
}

// newGroupsMembersListCmd creates the groups members list command
func newGroupsMembersListCmd() *cobra.Command {
	opts := &options.GroupsMembersListOptions{}

	cmd := &cobra.Command{
		Use:   "list <group-id>",
		Short: "List group members",
		Long: `List all members of a group.

Examples:
  canvas groups members list 456`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGroupsMembersList(cmd.Context(), client, opts)
		},
	}

	return cmd
}

// newGroupsMembersAddCmd creates the groups members add command
func newGroupsMembersAddCmd() *cobra.Command {
	opts := &options.GroupsMembersAddOptions{}

	cmd := &cobra.Command{
		Use:   "add <group-id>",
		Short: "Add a member to a group",
		Long: `Add a user to a group.

Examples:
  canvas groups members add 456 --user-id 789`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGroupsMembersAdd(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID to add (required)")
	cmd.MarkFlagRequired("user-id")

	return cmd
}

// newGroupsMembersRemoveCmd creates the groups members remove command
func newGroupsMembersRemoveCmd() *cobra.Command {
	opts := &options.GroupsMembersRemoveOptions{}

	cmd := &cobra.Command{
		Use:   "remove <group-id>",
		Short: "Remove a member from a group",
		Long: `Remove a user from a group by membership ID.

Examples:
  canvas groups members remove 456 --membership-id 123`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %w", err)
			}
			opts.GroupID = groupID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGroupsMembersRemove(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.MembershipID, "membership-id", 0, "Membership ID to remove (required)")
	cmd.MarkFlagRequired("membership-id")

	return cmd
}

// groupsCategoriesCmd represents the categories subcommand
var groupsCategoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "Manage group categories",
	Long:  "Commands for managing group categories.",
}

// newGroupsCategoriesListCmd creates the groups categories list command
func newGroupsCategoriesListCmd() *cobra.Command {
	opts := &options.GroupsCategoriesListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List group categories",
		Long: `List group categories for a course or account.

Examples:
  canvas groups categories list --course-id 123
  canvas groups categories list --account-id 1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGroupsCategoriesList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID")

	return cmd
}

// newGroupsCategoriesGetCmd creates the groups categories get command
func newGroupsCategoriesGetCmd() *cobra.Command {
	opts := &options.GroupsCategoriesGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <category-id>",
		Short: "Get group category details",
		Long: `Get details of a specific group category.

Examples:
  canvas groups categories get 456`,
		Args: ExactArgsWithUsage(1, "category-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			categoryID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid category ID: %w", err)
			}
			opts.CategoryID = categoryID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGroupsCategoriesGet(cmd.Context(), client, opts)
		},
	}

	return cmd
}

// newGroupsCategoriesCreateCmd creates the groups categories create command
func newGroupsCategoriesCreateCmd() *cobra.Command {
	opts := &options.GroupsCategoriesCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a group category",
		Long: `Create a new group category in a course or account.

Examples:
  canvas groups categories create --course-id 123 --name "Project Teams"
  canvas groups categories create --account-id 1 --name "Clubs" --self-signup enabled`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGroupsCategoriesCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Category name (required)")
	cmd.Flags().StringVar(&opts.SelfSignup, "self-signup", "", "Self signup (enabled, restricted)")
	cmd.Flags().StringVar(&opts.AutoLeader, "auto-leader", "", "Auto leader (first, random)")
	cmd.Flags().IntVar(&opts.GroupLimit, "group-limit", 0, "Group member limit")
	cmd.Flags().IntVar(&opts.CreateGroupCount, "create-group-count", 0, "Number of groups to create")
	cmd.Flags().IntVar(&opts.SplitGroupCount, "split-group-count", 0, "Number of groups to split students into")
	cmd.Flags().StringVar(&opts.SISCategoryID, "sis-category-id", "", "SIS category ID")
	cmd.MarkFlagRequired("name")

	return cmd
}

// newGroupsCategoriesUpdateCmd creates the groups categories update command
func newGroupsCategoriesUpdateCmd() *cobra.Command {
	opts := &options.GroupsCategoriesUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <category-id>",
		Short: "Update a group category",
		Long: `Update an existing group category.

Examples:
  canvas groups categories update 456 --name "New Name"
  canvas groups categories update 456 --self-signup restricted`,
		Args: ExactArgsWithUsage(1, "category-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			categoryID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid category ID: %w", err)
			}
			opts.CategoryID = categoryID

			// Track which fields were set
			opts.NameSet = cmd.Flags().Changed("name")
			opts.SelfSignupSet = cmd.Flags().Changed("self-signup")
			opts.AutoLeaderSet = cmd.Flags().Changed("auto-leader")
			opts.GroupLimitSet = cmd.Flags().Changed("group-limit")
			opts.SISCategoryIDSet = cmd.Flags().Changed("sis-category-id")

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGroupsCategoriesUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "Category name")
	cmd.Flags().StringVar(&opts.SelfSignup, "self-signup", "", "Self signup (enabled, restricted)")
	cmd.Flags().StringVar(&opts.AutoLeader, "auto-leader", "", "Auto leader (first, random)")
	cmd.Flags().IntVar(&opts.GroupLimit, "group-limit", 0, "Group member limit")
	cmd.Flags().StringVar(&opts.SISCategoryID, "sis-category-id", "", "SIS category ID")

	return cmd
}

// newGroupsCategoriesDeleteCmd creates the groups categories delete command
func newGroupsCategoriesDeleteCmd() *cobra.Command {
	opts := &options.GroupsCategoriesDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <category-id>",
		Short: "Delete a group category",
		Long: `Delete an existing group category.

Examples:
  canvas groups categories delete 456
  canvas groups categories delete 456 --force`,
		Args: ExactArgsWithUsage(1, "category-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			categoryID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid category ID: %w", err)
			}
			opts.CategoryID = categoryID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGroupsCategoriesDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Force, "force", false, "Skip confirmation prompt")

	return cmd
}

// newGroupsCategoriesGroupsCmd creates the groups categories groups command
func newGroupsCategoriesGroupsCmd() *cobra.Command {
	opts := &options.GroupsCategoriesGroupsOptions{}

	cmd := &cobra.Command{
		Use:   "groups <category-id>",
		Short: "List groups in a category",
		Long: `List all groups within a specific group category.

Examples:
  canvas groups categories groups 456`,
		Args: ExactArgsWithUsage(1, "category-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			categoryID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid category ID: %w", err)
			}
			opts.CategoryID = categoryID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runGroupsCategoriesGroups(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func init() {
	rootCmd.AddCommand(groupsCmd)
	groupsCmd.AddCommand(newGroupsListCmd())
	groupsCmd.AddCommand(newGroupsGetCmd())
	groupsCmd.AddCommand(newGroupsCreateCmd())
	groupsCmd.AddCommand(newGroupsUpdateCmd())
	groupsCmd.AddCommand(newGroupsDeleteCmd())
	groupsCmd.AddCommand(groupsMembersCmd)
	groupsCmd.AddCommand(groupsCategoriesCmd)

	// Members subcommands
	groupsMembersCmd.AddCommand(newGroupsMembersListCmd())
	groupsMembersCmd.AddCommand(newGroupsMembersAddCmd())
	groupsMembersCmd.AddCommand(newGroupsMembersRemoveCmd())

	// Categories subcommands
	groupsCategoriesCmd.AddCommand(newGroupsCategoriesListCmd())
	groupsCategoriesCmd.AddCommand(newGroupsCategoriesGetCmd())
	groupsCategoriesCmd.AddCommand(newGroupsCategoriesCreateCmd())
	groupsCategoriesCmd.AddCommand(newGroupsCategoriesUpdateCmd())
	groupsCategoriesCmd.AddCommand(newGroupsCategoriesDeleteCmd())
	groupsCategoriesCmd.AddCommand(newGroupsCategoriesGroupsCmd())
}

func runGroupsList(ctx context.Context, client *api.Client, opts *options.GroupsListOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "groups.list", map[string]interface{}{
		"course_id":  opts.CourseID,
		"account_id": opts.AccountID,
		"user_id":    opts.UserID,
	})

	service := api.NewGroupsService(client)

	apiOpts := &api.ListGroupsOptions{}
	if opts.IncludeUsers {
		apiOpts.Include = append(apiOpts.Include, "users")
	}
	if opts.IncludePermissions {
		apiOpts.Include = append(apiOpts.Include, "permissions")
	}

	var groups []api.Group
	var err error

	if opts.CourseID > 0 {
		groups, err = service.ListCourse(ctx, opts.CourseID, apiOpts)
	} else if opts.AccountID > 0 {
		groups, err = service.ListAccount(ctx, opts.AccountID, apiOpts)
	} else {
		// Default to user's groups (userID 0 means "self")
		groups, err = service.ListUser(ctx, opts.UserID, apiOpts)
	}

	if err != nil {
		logger.LogCommandError(ctx, "groups.list", err, map[string]interface{}{
			"course_id":  opts.CourseID,
			"account_id": opts.AccountID,
			"user_id":    opts.UserID,
		})
		return fmt.Errorf("failed to list groups: %w", err)
	}

	if err := formatEmptyOrOutput(groups, "No groups found"); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "groups.list", len(groups))
	return nil
}

func runGroupsGet(ctx context.Context, client *api.Client, opts *options.GroupsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "groups.get", map[string]interface{}{
		"group_id": opts.GroupID,
	})

	service := api.NewGroupsService(client)

	var include []string
	if opts.IncludeUsers {
		include = append(include, "users")
	}
	if opts.IncludePermissions {
		include = append(include, "permissions")
	}

	group, err := service.Get(ctx, opts.GroupID, include)
	if err != nil {
		logger.LogCommandError(ctx, "groups.get", err, map[string]interface{}{
			"group_id": opts.GroupID,
		})
		return fmt.Errorf("failed to get group: %w", err)
	}

	if err := formatOutput(group, nil); err != nil {
		return fmt.Errorf("failed to print result: %w", err)
	}

	logger.LogCommandComplete(ctx, "groups.get", 1)
	return nil
}

func runGroupsCreate(ctx context.Context, client *api.Client, opts *options.GroupsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "groups.create", map[string]interface{}{
		"category_id": opts.CategoryID,
		"name":        opts.Name,
	})

	service := api.NewGroupsService(client)

	params := &api.CreateGroupParams{
		Name:           opts.Name,
		Description:    opts.Description,
		IsPublic:       opts.IsPublic,
		JoinLevel:      opts.JoinLevel,
		StorageQuotaMb: opts.StorageQuotaMb,
		SISGroupID:     opts.SISGroupID,
	}

	group, err := service.Create(ctx, opts.CategoryID, params)
	if err != nil {
		logger.LogCommandError(ctx, "groups.create", err, map[string]interface{}{
			"category_id": opts.CategoryID,
			"name":        opts.Name,
		})
		return fmt.Errorf("failed to create group: %w", err)
	}

	fmt.Printf("Group created successfully (ID: %d)\n", group.ID)
	if err := formatOutput(group, nil); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.create", 1)
	return nil
}

func runGroupsUpdate(ctx context.Context, client *api.Client, opts *options.GroupsUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "groups.update", map[string]interface{}{
		"group_id": opts.GroupID,
	})

	service := api.NewGroupsService(client)

	params := &api.UpdateGroupParams{}

	if opts.NameSet {
		params.Name = &opts.Name
	}
	if opts.DescriptionSet {
		params.Description = &opts.Description
	}
	if opts.IsPublicSet {
		params.IsPublic = &opts.IsPublic
	}
	if opts.JoinLevelSet {
		params.JoinLevel = &opts.JoinLevel
	}
	if opts.StorageQuotaMbSet {
		params.StorageQuotaMb = &opts.StorageQuotaMb
	}
	if opts.SISGroupIDSet {
		params.SISGroupID = &opts.SISGroupID
	}

	group, err := service.Update(ctx, opts.GroupID, params)
	if err != nil {
		logger.LogCommandError(ctx, "groups.update", err, map[string]interface{}{
			"group_id": opts.GroupID,
		})
		return fmt.Errorf("failed to update group: %w", err)
	}

	fmt.Printf("Group updated successfully (ID: %d)\n", group.ID)
	if err := formatOutput(group, nil); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.update", 1)
	return nil
}

func runGroupsDelete(ctx context.Context, client *api.Client, opts *options.GroupsDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "groups.delete", map[string]interface{}{
		"group_id": opts.GroupID,
	})

	// Confirm deletion
	confirmed, err := confirmDelete("group", opts.GroupID, opts.Force)
	if err != nil {
		return err
	}
	if !confirmed {
		fmt.Println("Cancelled.")
		return nil
	}

	service := api.NewGroupsService(client)

	group, err := service.Delete(ctx, opts.GroupID)
	if err != nil {
		logger.LogCommandError(ctx, "groups.delete", err, map[string]interface{}{
			"group_id": opts.GroupID,
		})
		return fmt.Errorf("failed to delete group: %w", err)
	}

	fmt.Printf("Group %d deleted successfully\n", group.ID)

	logger.LogCommandComplete(ctx, "groups.delete", 1)
	return nil
}

func runGroupsMembersList(ctx context.Context, client *api.Client, opts *options.GroupsMembersListOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "groups.members.list", map[string]interface{}{
		"group_id": opts.GroupID,
	})

	service := api.NewGroupsService(client)

	users, err := service.ListMembers(ctx, opts.GroupID)
	if err != nil {
		logger.LogCommandError(ctx, "groups.members.list", err, map[string]interface{}{
			"group_id": opts.GroupID,
		})
		return fmt.Errorf("failed to list group members: %w", err)
	}

	if err := formatEmptyOrOutput(users, "No members found in group"); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "groups.members.list", len(users))
	return nil
}

func runGroupsMembersAdd(ctx context.Context, client *api.Client, opts *options.GroupsMembersAddOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "groups.members.add", map[string]interface{}{
		"group_id": opts.GroupID,
		"user_id":  opts.UserID,
	})

	service := api.NewGroupsService(client)

	membership, err := service.AddMember(ctx, opts.GroupID, opts.UserID)
	if err != nil {
		logger.LogCommandError(ctx, "groups.members.add", err, map[string]interface{}{
			"group_id": opts.GroupID,
			"user_id":  opts.UserID,
		})
		return fmt.Errorf("failed to add member: %w", err)
	}

	fmt.Printf("User %d added to group %d\n", opts.UserID, opts.GroupID)
	if err := formatOutput(membership, nil); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.members.add", 1)
	return nil
}

func runGroupsMembersRemove(ctx context.Context, client *api.Client, opts *options.GroupsMembersRemoveOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "groups.members.remove", map[string]interface{}{
		"group_id":      opts.GroupID,
		"membership_id": opts.MembershipID,
	})

	service := api.NewGroupsService(client)

	err := service.RemoveMember(ctx, opts.GroupID, opts.MembershipID)
	if err != nil {
		logger.LogCommandError(ctx, "groups.members.remove", err, map[string]interface{}{
			"group_id":      opts.GroupID,
			"membership_id": opts.MembershipID,
		})
		return fmt.Errorf("failed to remove member: %w", err)
	}

	fmt.Printf("Membership %d removed from group %d\n", opts.MembershipID, opts.GroupID)

	logger.LogCommandComplete(ctx, "groups.members.remove", 1)
	return nil
}

func runGroupsCategoriesList(ctx context.Context, client *api.Client, opts *options.GroupsCategoriesListOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "groups.categories.list", map[string]interface{}{
		"course_id":  opts.CourseID,
		"account_id": opts.AccountID,
	})

	service := api.NewGroupsService(client)

	var categories []api.GroupCategory
	var err error

	if opts.CourseID > 0 {
		categories, err = service.ListCategoriesCourse(ctx, opts.CourseID, nil)
	} else {
		categories, err = service.ListCategoriesAccount(ctx, opts.AccountID, nil)
	}

	if err != nil {
		logger.LogCommandError(ctx, "groups.categories.list", err, map[string]interface{}{
			"course_id":  opts.CourseID,
			"account_id": opts.AccountID,
		})
		return fmt.Errorf("failed to list group categories: %w", err)
	}

	if err := formatEmptyOrOutput(categories, "No group categories found"); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "groups.categories.list", len(categories))
	return nil
}

func runGroupsCategoriesGet(ctx context.Context, client *api.Client, opts *options.GroupsCategoriesGetOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "groups.categories.get", map[string]interface{}{
		"category_id": opts.CategoryID,
	})

	service := api.NewGroupsService(client)

	category, err := service.GetCategory(ctx, opts.CategoryID)
	if err != nil {
		logger.LogCommandError(ctx, "groups.categories.get", err, map[string]interface{}{
			"category_id": opts.CategoryID,
		})
		return fmt.Errorf("failed to get category: %w", err)
	}

	if err := formatOutput(category, nil); err != nil {
		return fmt.Errorf("failed to print result: %w", err)
	}

	logger.LogCommandComplete(ctx, "groups.categories.get", 1)
	return nil
}

func runGroupsCategoriesCreate(ctx context.Context, client *api.Client, opts *options.GroupsCategoriesCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "groups.categories.create", map[string]interface{}{
		"course_id":  opts.CourseID,
		"account_id": opts.AccountID,
		"name":       opts.Name,
	})

	service := api.NewGroupsService(client)

	params := &api.CreateCategoryParams{
		Name:               opts.Name,
		SelfSignup:         opts.SelfSignup,
		AutoLeader:         opts.AutoLeader,
		GroupLimit:         opts.GroupLimit,
		CreateGroupCount:   opts.CreateGroupCount,
		SplitGroupCount:    opts.SplitGroupCount,
		SISGroupCategoryID: opts.SISCategoryID,
	}

	var category *api.GroupCategory
	var err error

	if opts.CourseID > 0 {
		category, err = service.CreateCategoryCourse(ctx, opts.CourseID, params)
	} else {
		category, err = service.CreateCategoryAccount(ctx, opts.AccountID, params)
	}

	if err != nil {
		logger.LogCommandError(ctx, "groups.categories.create", err, map[string]interface{}{
			"course_id":  opts.CourseID,
			"account_id": opts.AccountID,
			"name":       opts.Name,
		})
		return fmt.Errorf("failed to create category: %w", err)
	}

	fmt.Printf("Category created successfully (ID: %d)\n", category.ID)
	if err := formatOutput(category, nil); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.categories.create", 1)
	return nil
}

func runGroupsCategoriesUpdate(ctx context.Context, client *api.Client, opts *options.GroupsCategoriesUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "groups.categories.update", map[string]interface{}{
		"category_id": opts.CategoryID,
	})

	service := api.NewGroupsService(client)

	params := &api.UpdateCategoryParams{}

	if opts.NameSet {
		params.Name = &opts.Name
	}
	if opts.SelfSignupSet {
		params.SelfSignup = &opts.SelfSignup
	}
	if opts.AutoLeaderSet {
		params.AutoLeader = &opts.AutoLeader
	}
	if opts.GroupLimitSet {
		params.GroupLimit = &opts.GroupLimit
	}
	if opts.SISCategoryIDSet {
		params.SISGroupCategoryID = &opts.SISCategoryID
	}

	category, err := service.UpdateCategory(ctx, opts.CategoryID, params)
	if err != nil {
		logger.LogCommandError(ctx, "groups.categories.update", err, map[string]interface{}{
			"category_id": opts.CategoryID,
		})
		return fmt.Errorf("failed to update category: %w", err)
	}

	fmt.Printf("Category updated successfully (ID: %d)\n", category.ID)
	if err := formatOutput(category, nil); err != nil {
		return err
	}

	logger.LogCommandComplete(ctx, "groups.categories.update", 1)
	return nil
}

func runGroupsCategoriesDelete(ctx context.Context, client *api.Client, opts *options.GroupsCategoriesDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "groups.categories.delete", map[string]interface{}{
		"category_id": opts.CategoryID,
	})

	// Confirm deletion
	if !opts.Force {
		fmt.Printf("WARNING: This will delete category %d and all groups in it.\n", opts.CategoryID)
		fmt.Print("Type 'yes' to confirm: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	service := api.NewGroupsService(client)

	category, err := service.DeleteCategory(ctx, opts.CategoryID)
	if err != nil {
		logger.LogCommandError(ctx, "groups.categories.delete", err, map[string]interface{}{
			"category_id": opts.CategoryID,
		})
		return fmt.Errorf("failed to delete category: %w", err)
	}

	fmt.Printf("Category %d deleted successfully\n", category.ID)

	logger.LogCommandComplete(ctx, "groups.categories.delete", 1)
	return nil
}

func runGroupsCategoriesGroups(ctx context.Context, client *api.Client, opts *options.GroupsCategoriesGroupsOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "groups.categories.groups", map[string]interface{}{
		"category_id": opts.CategoryID,
	})

	service := api.NewGroupsService(client)

	groups, err := service.ListGroupsInCategory(ctx, opts.CategoryID)
	if err != nil {
		logger.LogCommandError(ctx, "groups.categories.groups", err, map[string]interface{}{
			"category_id": opts.CategoryID,
		})
		return fmt.Errorf("failed to list groups in category: %w", err)
	}

	if err := formatEmptyOrOutput(groups, "No groups found in category"); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "groups.categories.groups", len(groups))
	return nil
}
