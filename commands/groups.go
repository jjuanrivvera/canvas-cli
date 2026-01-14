package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var (
	// Common flags
	groupsCourseID   int64
	groupsAccountID  int64
	groupsUserID     int64
	groupsCategoryID int64

	// Include flags
	groupsIncludeUsers       bool
	groupsIncludePermissions bool

	// Create/Update group flags
	groupsName           string
	groupsDescription    string
	groupsIsPublic       bool
	groupsJoinLevel      string
	groupsStorageQuotaMb int64
	groupsSISGroupID     string

	// Create/Update category flags
	groupsCategoryName     string
	groupsSelfSignup       string
	groupsAutoLeader       string
	groupsGroupLimit       int
	groupsCreateGroupCount int
	groupsSplitGroupCount  int
	groupsSISCategoryID    string

	// Member flags
	groupsMembershipID int64

	// Delete flags
	groupsForce bool
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

// groupsListCmd lists groups
var groupsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List groups",
	Long: `List groups for a course, account, or user.

Examples:
  canvas groups list --course-id 123
  canvas groups list --account-id 1
  canvas groups list --user-id 456
  canvas groups list  # Lists current user's groups`,
	RunE: runGroupsList,
}

// groupsGetCmd gets a single group
var groupsGetCmd = &cobra.Command{
	Use:   "get <group-id>",
	Short: "Get group details",
	Long: `Get details of a specific group.

Examples:
  canvas groups get 456
  canvas groups get 456 --include-users`,
	Args: ExactArgsWithUsage(1, "group-id"),
	RunE: runGroupsGet,
}

// groupsCreateCmd creates a new group
var groupsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new group",
	Long: `Create a new group in a group category.

Examples:
  canvas groups create --category-id 123 --name "Study Group"
  canvas groups create --category-id 123 --name "Project Team" --description "Our project team"`,
	RunE: runGroupsCreate,
}

// groupsUpdateCmd updates a group
var groupsUpdateCmd = &cobra.Command{
	Use:   "update <group-id>",
	Short: "Update a group",
	Long: `Update an existing group.

Examples:
  canvas groups update 456 --name "New Name"
  canvas groups update 456 --description "Updated description"`,
	Args: ExactArgsWithUsage(1, "group-id"),
	RunE: runGroupsUpdate,
}

// groupsDeleteCmd deletes a group
var groupsDeleteCmd = &cobra.Command{
	Use:   "delete <group-id>",
	Short: "Delete a group",
	Long: `Delete an existing group.

Examples:
  canvas groups delete 456
  canvas groups delete 456 --force`,
	Args: ExactArgsWithUsage(1, "group-id"),
	RunE: runGroupsDelete,
}

// groupsMembersCmd represents the members subcommand
var groupsMembersCmd = &cobra.Command{
	Use:   "members",
	Short: "Manage group members",
	Long:  "Commands for managing group memberships.",
}

// groupsMembersListCmd lists group members
var groupsMembersListCmd = &cobra.Command{
	Use:   "list <group-id>",
	Short: "List group members",
	Long: `List all members of a group.

Examples:
  canvas groups members list 456`,
	Args: ExactArgsWithUsage(1, "group-id"),
	RunE: runGroupsMembersList,
}

// groupsMembersAddCmd adds a member to a group
var groupsMembersAddCmd = &cobra.Command{
	Use:   "add <group-id>",
	Short: "Add a member to a group",
	Long: `Add a user to a group.

Examples:
  canvas groups members add 456 --user-id 789`,
	Args: ExactArgsWithUsage(1, "group-id"),
	RunE: runGroupsMembersAdd,
}

// groupsMembersRemoveCmd removes a member from a group
var groupsMembersRemoveCmd = &cobra.Command{
	Use:   "remove <group-id>",
	Short: "Remove a member from a group",
	Long: `Remove a user from a group by membership ID.

Examples:
  canvas groups members remove 456 --membership-id 123`,
	Args: ExactArgsWithUsage(1, "group-id"),
	RunE: runGroupsMembersRemove,
}

// groupsCategoriesCmd represents the categories subcommand
var groupsCategoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "Manage group categories",
	Long:  "Commands for managing group categories.",
}

// groupsCategoriesListCmd lists group categories
var groupsCategoriesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List group categories",
	Long: `List group categories for a course or account.

Examples:
  canvas groups categories list --course-id 123
  canvas groups categories list --account-id 1`,
	RunE: runGroupsCategoriesList,
}

// groupsCategoriesGetCmd gets a single category
var groupsCategoriesGetCmd = &cobra.Command{
	Use:   "get <category-id>",
	Short: "Get group category details",
	Long: `Get details of a specific group category.

Examples:
  canvas groups categories get 456`,
	Args: ExactArgsWithUsage(1, "category-id"),
	RunE: runGroupsCategoriesGet,
}

// groupsCategoriesCreateCmd creates a new category
var groupsCategoriesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a group category",
	Long: `Create a new group category in a course or account.

Examples:
  canvas groups categories create --course-id 123 --name "Project Teams"
  canvas groups categories create --account-id 1 --name "Clubs" --self-signup enabled`,
	RunE: runGroupsCategoriesCreate,
}

// groupsCategoriesUpdateCmd updates a category
var groupsCategoriesUpdateCmd = &cobra.Command{
	Use:   "update <category-id>",
	Short: "Update a group category",
	Long: `Update an existing group category.

Examples:
  canvas groups categories update 456 --name "New Name"
  canvas groups categories update 456 --self-signup restricted`,
	Args: ExactArgsWithUsage(1, "category-id"),
	RunE: runGroupsCategoriesUpdate,
}

// groupsCategoriesDeleteCmd deletes a category
var groupsCategoriesDeleteCmd = &cobra.Command{
	Use:   "delete <category-id>",
	Short: "Delete a group category",
	Long: `Delete an existing group category.

Examples:
  canvas groups categories delete 456
  canvas groups categories delete 456 --force`,
	Args: ExactArgsWithUsage(1, "category-id"),
	RunE: runGroupsCategoriesDelete,
}

// groupsCategoriesGroupsCmd lists groups in a category
var groupsCategoriesGroupsCmd = &cobra.Command{
	Use:   "groups <category-id>",
	Short: "List groups in a category",
	Long: `List all groups within a specific group category.

Examples:
  canvas groups categories groups 456`,
	Args: ExactArgsWithUsage(1, "category-id"),
	RunE: runGroupsCategoriesGroups,
}

func init() {
	rootCmd.AddCommand(groupsCmd)
	groupsCmd.AddCommand(groupsListCmd)
	groupsCmd.AddCommand(groupsGetCmd)
	groupsCmd.AddCommand(groupsCreateCmd)
	groupsCmd.AddCommand(groupsUpdateCmd)
	groupsCmd.AddCommand(groupsDeleteCmd)
	groupsCmd.AddCommand(groupsMembersCmd)
	groupsCmd.AddCommand(groupsCategoriesCmd)

	// Members subcommands
	groupsMembersCmd.AddCommand(groupsMembersListCmd)
	groupsMembersCmd.AddCommand(groupsMembersAddCmd)
	groupsMembersCmd.AddCommand(groupsMembersRemoveCmd)

	// Categories subcommands
	groupsCategoriesCmd.AddCommand(groupsCategoriesListCmd)
	groupsCategoriesCmd.AddCommand(groupsCategoriesGetCmd)
	groupsCategoriesCmd.AddCommand(groupsCategoriesCreateCmd)
	groupsCategoriesCmd.AddCommand(groupsCategoriesUpdateCmd)
	groupsCategoriesCmd.AddCommand(groupsCategoriesDeleteCmd)
	groupsCategoriesCmd.AddCommand(groupsCategoriesGroupsCmd)

	// List flags
	groupsListCmd.Flags().Int64Var(&groupsCourseID, "course-id", 0, "Course ID")
	groupsListCmd.Flags().Int64Var(&groupsAccountID, "account-id", 0, "Account ID")
	groupsListCmd.Flags().Int64Var(&groupsUserID, "user-id", 0, "User ID (0 for self)")
	groupsListCmd.Flags().BoolVar(&groupsIncludeUsers, "include-users", false, "Include group users")
	groupsListCmd.Flags().BoolVar(&groupsIncludePermissions, "include-permissions", false, "Include permissions")

	// Get flags
	groupsGetCmd.Flags().BoolVar(&groupsIncludeUsers, "include-users", false, "Include group users")
	groupsGetCmd.Flags().BoolVar(&groupsIncludePermissions, "include-permissions", false, "Include permissions")

	// Create flags
	groupsCreateCmd.Flags().Int64Var(&groupsCategoryID, "category-id", 0, "Group category ID (required)")
	groupsCreateCmd.Flags().StringVar(&groupsName, "name", "", "Group name (required)")
	groupsCreateCmd.Flags().StringVar(&groupsDescription, "description", "", "Group description")
	groupsCreateCmd.Flags().BoolVar(&groupsIsPublic, "public", false, "Whether the group is public")
	groupsCreateCmd.Flags().StringVar(&groupsJoinLevel, "join-level", "", "Join level (parent_context_auto_join, parent_context_request, invitation_only)")
	groupsCreateCmd.Flags().Int64Var(&groupsStorageQuotaMb, "storage-quota-mb", 0, "Storage quota in MB")
	groupsCreateCmd.Flags().StringVar(&groupsSISGroupID, "sis-group-id", "", "SIS group ID")
	groupsCreateCmd.MarkFlagRequired("category-id")
	groupsCreateCmd.MarkFlagRequired("name")

	// Update flags
	groupsUpdateCmd.Flags().StringVar(&groupsName, "name", "", "Group name")
	groupsUpdateCmd.Flags().StringVar(&groupsDescription, "description", "", "Group description")
	groupsUpdateCmd.Flags().BoolVar(&groupsIsPublic, "public", false, "Whether the group is public")
	groupsUpdateCmd.Flags().StringVar(&groupsJoinLevel, "join-level", "", "Join level")
	groupsUpdateCmd.Flags().Int64Var(&groupsStorageQuotaMb, "storage-quota-mb", 0, "Storage quota in MB")
	groupsUpdateCmd.Flags().StringVar(&groupsSISGroupID, "sis-group-id", "", "SIS group ID")

	// Delete flags
	groupsDeleteCmd.Flags().BoolVar(&groupsForce, "force", false, "Skip confirmation prompt")

	// Members add flags
	groupsMembersAddCmd.Flags().Int64Var(&groupsUserID, "user-id", 0, "User ID to add (required)")
	groupsMembersAddCmd.MarkFlagRequired("user-id")

	// Members remove flags
	groupsMembersRemoveCmd.Flags().Int64Var(&groupsMembershipID, "membership-id", 0, "Membership ID to remove (required)")
	groupsMembersRemoveCmd.MarkFlagRequired("membership-id")

	// Categories list flags
	groupsCategoriesListCmd.Flags().Int64Var(&groupsCourseID, "course-id", 0, "Course ID")
	groupsCategoriesListCmd.Flags().Int64Var(&groupsAccountID, "account-id", 0, "Account ID")

	// Categories create flags
	groupsCategoriesCreateCmd.Flags().Int64Var(&groupsCourseID, "course-id", 0, "Course ID")
	groupsCategoriesCreateCmd.Flags().Int64Var(&groupsAccountID, "account-id", 0, "Account ID")
	groupsCategoriesCreateCmd.Flags().StringVar(&groupsCategoryName, "name", "", "Category name (required)")
	groupsCategoriesCreateCmd.Flags().StringVar(&groupsSelfSignup, "self-signup", "", "Self signup (enabled, restricted)")
	groupsCategoriesCreateCmd.Flags().StringVar(&groupsAutoLeader, "auto-leader", "", "Auto leader (first, random)")
	groupsCategoriesCreateCmd.Flags().IntVar(&groupsGroupLimit, "group-limit", 0, "Group member limit")
	groupsCategoriesCreateCmd.Flags().IntVar(&groupsCreateGroupCount, "create-group-count", 0, "Number of groups to create")
	groupsCategoriesCreateCmd.Flags().IntVar(&groupsSplitGroupCount, "split-group-count", 0, "Number of groups to split students into")
	groupsCategoriesCreateCmd.Flags().StringVar(&groupsSISCategoryID, "sis-category-id", "", "SIS category ID")
	groupsCategoriesCreateCmd.MarkFlagRequired("name")

	// Categories update flags
	groupsCategoriesUpdateCmd.Flags().StringVar(&groupsCategoryName, "name", "", "Category name")
	groupsCategoriesUpdateCmd.Flags().StringVar(&groupsSelfSignup, "self-signup", "", "Self signup (enabled, restricted)")
	groupsCategoriesUpdateCmd.Flags().StringVar(&groupsAutoLeader, "auto-leader", "", "Auto leader (first, random)")
	groupsCategoriesUpdateCmd.Flags().IntVar(&groupsGroupLimit, "group-limit", 0, "Group member limit")
	groupsCategoriesUpdateCmd.Flags().StringVar(&groupsSISCategoryID, "sis-category-id", "", "SIS category ID")

	// Categories delete flags
	groupsCategoriesDeleteCmd.Flags().BoolVar(&groupsForce, "force", false, "Skip confirmation prompt")
}

func runGroupsList(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewGroupsService(client)

	opts := &api.ListGroupsOptions{}
	if groupsIncludeUsers {
		opts.Include = append(opts.Include, "users")
	}
	if groupsIncludePermissions {
		opts.Include = append(opts.Include, "permissions")
	}

	var groups []api.Group
	ctx := context.Background()

	if groupsCourseID > 0 {
		groups, err = service.ListCourse(ctx, groupsCourseID, opts)
	} else if groupsAccountID > 0 {
		groups, err = service.ListAccount(ctx, groupsAccountID, opts)
	} else {
		// Default to user's groups (userID 0 means "self")
		groups, err = service.ListUser(ctx, groupsUserID, opts)
	}

	if err != nil {
		return fmt.Errorf("failed to list groups: %w", err)
	}

	if len(groups) == 0 {
		fmt.Println("No groups found")
		return nil
	}

	printVerbose("Found %d groups:\n\n", len(groups))
	return formatOutput(groups, nil)
}

func runGroupsGet(cmd *cobra.Command, args []string) error {
	groupID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid group ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewGroupsService(client)

	var include []string
	if groupsIncludeUsers {
		include = append(include, "users")
	}
	if groupsIncludePermissions {
		include = append(include, "permissions")
	}

	group, err := service.Get(context.Background(), groupID, include)
	if err != nil {
		return fmt.Errorf("failed to get group: %w", err)
	}

	return formatOutput(group, nil)
}

func runGroupsCreate(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewGroupsService(client)

	params := &api.CreateGroupParams{
		Name:           groupsName,
		Description:    groupsDescription,
		IsPublic:       groupsIsPublic,
		JoinLevel:      groupsJoinLevel,
		StorageQuotaMb: groupsStorageQuotaMb,
		SISGroupID:     groupsSISGroupID,
	}

	group, err := service.Create(context.Background(), groupsCategoryID, params)
	if err != nil {
		return fmt.Errorf("failed to create group: %w", err)
	}

	fmt.Printf("Group created successfully (ID: %d)\n", group.ID)
	return formatOutput(group, nil)
}

func runGroupsUpdate(cmd *cobra.Command, args []string) error {
	groupID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid group ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewGroupsService(client)

	params := &api.UpdateGroupParams{}

	if cmd.Flags().Changed("name") {
		params.Name = &groupsName
	}
	if cmd.Flags().Changed("description") {
		params.Description = &groupsDescription
	}
	if cmd.Flags().Changed("public") {
		params.IsPublic = &groupsIsPublic
	}
	if cmd.Flags().Changed("join-level") {
		params.JoinLevel = &groupsJoinLevel
	}
	if cmd.Flags().Changed("storage-quota-mb") {
		params.StorageQuotaMb = &groupsStorageQuotaMb
	}
	if cmd.Flags().Changed("sis-group-id") {
		params.SISGroupID = &groupsSISGroupID
	}

	group, err := service.Update(context.Background(), groupID, params)
	if err != nil {
		return fmt.Errorf("failed to update group: %w", err)
	}

	fmt.Printf("Group updated successfully (ID: %d)\n", group.ID)
	return formatOutput(group, nil)
}

func runGroupsDelete(cmd *cobra.Command, args []string) error {
	groupID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid group ID: %w", err)
	}

	if !groupsForce {
		fmt.Printf("WARNING: This will delete group %d.\n", groupID)
		fmt.Print("Type 'yes' to confirm: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Delete cancelled")
			return nil
		}
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewGroupsService(client)

	group, err := service.Delete(context.Background(), groupID)
	if err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}

	fmt.Printf("Group %d deleted successfully\n", group.ID)
	return nil
}

func runGroupsMembersList(cmd *cobra.Command, args []string) error {
	groupID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid group ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewGroupsService(client)

	users, err := service.ListMembers(context.Background(), groupID)
	if err != nil {
		return fmt.Errorf("failed to list group members: %w", err)
	}

	if len(users) == 0 {
		fmt.Println("No members found in group")
		return nil
	}

	printVerbose("Found %d members:\n\n", len(users))
	return formatOutput(users, nil)
}

func runGroupsMembersAdd(cmd *cobra.Command, args []string) error {
	groupID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid group ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewGroupsService(client)

	membership, err := service.AddMember(context.Background(), groupID, groupsUserID)
	if err != nil {
		return fmt.Errorf("failed to add member: %w", err)
	}

	fmt.Printf("User %d added to group %d\n", groupsUserID, groupID)
	return formatOutput(membership, nil)
}

func runGroupsMembersRemove(cmd *cobra.Command, args []string) error {
	groupID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid group ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewGroupsService(client)

	err = service.RemoveMember(context.Background(), groupID, groupsMembershipID)
	if err != nil {
		return fmt.Errorf("failed to remove member: %w", err)
	}

	fmt.Printf("Membership %d removed from group %d\n", groupsMembershipID, groupID)
	return nil
}

func runGroupsCategoriesList(cmd *cobra.Command, args []string) error {
	if groupsCourseID == 0 && groupsAccountID == 0 {
		return fmt.Errorf("must specify either --course-id or --account-id")
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewGroupsService(client)

	var categories []api.GroupCategory
	ctx := context.Background()

	if groupsCourseID > 0 {
		categories, err = service.ListCategoriesCourse(ctx, groupsCourseID, nil)
	} else {
		categories, err = service.ListCategoriesAccount(ctx, groupsAccountID, nil)
	}

	if err != nil {
		return fmt.Errorf("failed to list group categories: %w", err)
	}

	if len(categories) == 0 {
		fmt.Println("No group categories found")
		return nil
	}

	printVerbose("Found %d categories:\n\n", len(categories))
	return formatOutput(categories, nil)
}

func runGroupsCategoriesGet(cmd *cobra.Command, args []string) error {
	categoryID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid category ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewGroupsService(client)

	category, err := service.GetCategory(context.Background(), categoryID)
	if err != nil {
		return fmt.Errorf("failed to get category: %w", err)
	}

	return formatOutput(category, nil)
}

func runGroupsCategoriesCreate(cmd *cobra.Command, args []string) error {
	if groupsCourseID == 0 && groupsAccountID == 0 {
		return fmt.Errorf("must specify either --course-id or --account-id")
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewGroupsService(client)

	params := &api.CreateCategoryParams{
		Name:               groupsCategoryName,
		SelfSignup:         groupsSelfSignup,
		AutoLeader:         groupsAutoLeader,
		GroupLimit:         groupsGroupLimit,
		CreateGroupCount:   groupsCreateGroupCount,
		SplitGroupCount:    groupsSplitGroupCount,
		SISGroupCategoryID: groupsSISCategoryID,
	}

	var category *api.GroupCategory
	ctx := context.Background()

	if groupsCourseID > 0 {
		category, err = service.CreateCategoryCourse(ctx, groupsCourseID, params)
	} else {
		category, err = service.CreateCategoryAccount(ctx, groupsAccountID, params)
	}

	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}

	fmt.Printf("Category created successfully (ID: %d)\n", category.ID)
	return formatOutput(category, nil)
}

func runGroupsCategoriesUpdate(cmd *cobra.Command, args []string) error {
	categoryID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid category ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewGroupsService(client)

	params := &api.UpdateCategoryParams{}

	if cmd.Flags().Changed("name") {
		params.Name = &groupsCategoryName
	}
	if cmd.Flags().Changed("self-signup") {
		params.SelfSignup = &groupsSelfSignup
	}
	if cmd.Flags().Changed("auto-leader") {
		params.AutoLeader = &groupsAutoLeader
	}
	if cmd.Flags().Changed("group-limit") {
		params.GroupLimit = &groupsGroupLimit
	}
	if cmd.Flags().Changed("sis-category-id") {
		params.SISGroupCategoryID = &groupsSISCategoryID
	}

	category, err := service.UpdateCategory(context.Background(), categoryID, params)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}

	fmt.Printf("Category updated successfully (ID: %d)\n", category.ID)
	return formatOutput(category, nil)
}

func runGroupsCategoriesDelete(cmd *cobra.Command, args []string) error {
	categoryID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid category ID: %w", err)
	}

	if !groupsForce {
		fmt.Printf("WARNING: This will delete category %d and all groups in it.\n", categoryID)
		fmt.Print("Type 'yes' to confirm: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Delete cancelled")
			return nil
		}
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewGroupsService(client)

	category, err := service.DeleteCategory(context.Background(), categoryID)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	fmt.Printf("Category %d deleted successfully\n", category.ID)
	return nil
}

func runGroupsCategoriesGroups(cmd *cobra.Command, args []string) error {
	categoryID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid category ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewGroupsService(client)

	groups, err := service.ListGroupsInCategory(context.Background(), categoryID)
	if err != nil {
		return fmt.Errorf("failed to list groups in category: %w", err)
	}

	if len(groups) == 0 {
		fmt.Println("No groups found in category")
		return nil
	}

	printVerbose("Found %d groups:\n\n", len(groups))
	return formatOutput(groups, nil)
}
