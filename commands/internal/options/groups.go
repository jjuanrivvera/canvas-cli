package options

import "fmt"

// GroupsListOptions contains options for listing groups
type GroupsListOptions struct {
	CourseID           int64
	AccountID          int64
	UserID             int64
	IncludeUsers       bool
	IncludePermissions bool
}

// Validate validates the options
func (o *GroupsListOptions) Validate() error {
	// No required fields - defaults to user's groups
	return nil
}

// GroupsGetOptions contains options for getting a group
type GroupsGetOptions struct {
	GroupID            int64
	IncludeUsers       bool
	IncludePermissions bool
}

// Validate validates the options
func (o *GroupsGetOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	return nil
}

// GroupsCreateOptions contains options for creating a group
type GroupsCreateOptions struct {
	CategoryID     int64
	Name           string
	Description    string
	IsPublic       bool
	JoinLevel      string
	StorageQuotaMb int64
	SISGroupID     string
}

// Validate validates the options
func (o *GroupsCreateOptions) Validate() error {
	if o.CategoryID <= 0 {
		return fmt.Errorf("category-id is required and must be greater than 0")
	}
	if o.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

// GroupsUpdateOptions contains options for updating a group
type GroupsUpdateOptions struct {
	GroupID        int64
	Name           string
	Description    string
	IsPublic       bool
	JoinLevel      string
	StorageQuotaMb int64
	SISGroupID     string
	// Track which fields were actually set
	NameSet           bool
	DescriptionSet    bool
	IsPublicSet       bool
	JoinLevelSet      bool
	StorageQuotaMbSet bool
	SISGroupIDSet     bool
}

// Validate validates the options
func (o *GroupsUpdateOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	// At least one field must be set for update
	if !o.NameSet && !o.DescriptionSet && !o.IsPublicSet &&
		!o.JoinLevelSet && !o.StorageQuotaMbSet && !o.SISGroupIDSet {
		return fmt.Errorf("at least one field must be specified for update")
	}
	return nil
}

// GroupsDeleteOptions contains options for deleting a group
type GroupsDeleteOptions struct {
	GroupID int64
	Force   bool
}

// Validate validates the options
func (o *GroupsDeleteOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	return nil
}

// GroupsMembersListOptions contains options for listing group members
type GroupsMembersListOptions struct {
	GroupID int64
}

// Validate validates the options
func (o *GroupsMembersListOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	return nil
}

// GroupsMembersAddOptions contains options for adding a group member
type GroupsMembersAddOptions struct {
	GroupID int64
	UserID  int64
}

// Validate validates the options
func (o *GroupsMembersAddOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}

// GroupsMembersRemoveOptions contains options for removing a group member
type GroupsMembersRemoveOptions struct {
	GroupID      int64
	MembershipID int64
}

// Validate validates the options
func (o *GroupsMembersRemoveOptions) Validate() error {
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	if o.MembershipID <= 0 {
		return fmt.Errorf("membership-id is required and must be greater than 0")
	}
	return nil
}

// GroupsCategoriesListOptions contains options for listing group categories
type GroupsCategoriesListOptions struct {
	CourseID  int64
	AccountID int64
}

// Validate validates the options
func (o *GroupsCategoriesListOptions) Validate() error {
	if o.CourseID == 0 && o.AccountID == 0 {
		return fmt.Errorf("must specify either --course-id or --account-id")
	}
	return nil
}

// GroupsCategoriesGetOptions contains options for getting a group category
type GroupsCategoriesGetOptions struct {
	CategoryID int64
}

// Validate validates the options
func (o *GroupsCategoriesGetOptions) Validate() error {
	if o.CategoryID <= 0 {
		return fmt.Errorf("category-id is required and must be greater than 0")
	}
	return nil
}

// GroupsCategoriesCreateOptions contains options for creating a group category
type GroupsCategoriesCreateOptions struct {
	CourseID         int64
	AccountID        int64
	Name             string
	SelfSignup       string
	AutoLeader       string
	GroupLimit       int
	CreateGroupCount int
	SplitGroupCount  int
	SISCategoryID    string
}

// Validate validates the options
func (o *GroupsCategoriesCreateOptions) Validate() error {
	if o.CourseID == 0 && o.AccountID == 0 {
		return fmt.Errorf("must specify either --course-id or --account-id")
	}
	if o.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

// GroupsCategoriesUpdateOptions contains options for updating a group category
type GroupsCategoriesUpdateOptions struct {
	CategoryID    int64
	Name          string
	SelfSignup    string
	AutoLeader    string
	GroupLimit    int
	SISCategoryID string
	// Track which fields were actually set
	NameSet          bool
	SelfSignupSet    bool
	AutoLeaderSet    bool
	GroupLimitSet    bool
	SISCategoryIDSet bool
}

// Validate validates the options
func (o *GroupsCategoriesUpdateOptions) Validate() error {
	if o.CategoryID <= 0 {
		return fmt.Errorf("category-id is required and must be greater than 0")
	}
	// At least one field must be set for update
	if !o.NameSet && !o.SelfSignupSet && !o.AutoLeaderSet &&
		!o.GroupLimitSet && !o.SISCategoryIDSet {
		return fmt.Errorf("at least one field must be specified for update")
	}
	return nil
}

// GroupsCategoriesDeleteOptions contains options for deleting a group category
type GroupsCategoriesDeleteOptions struct {
	CategoryID int64
	Force      bool
}

// Validate validates the options
func (o *GroupsCategoriesDeleteOptions) Validate() error {
	if o.CategoryID <= 0 {
		return fmt.Errorf("category-id is required and must be greater than 0")
	}
	return nil
}

// GroupsCategoriesGroupsOptions contains options for listing groups in a category
type GroupsCategoriesGroupsOptions struct {
	CategoryID int64
}

// Validate validates the options
func (o *GroupsCategoriesGroupsOptions) Validate() error {
	if o.CategoryID <= 0 {
		return fmt.Errorf("category-id is required and must be greater than 0")
	}
	return nil
}
