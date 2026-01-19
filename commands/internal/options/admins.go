package options

import "fmt"

// AdminsListOptions contains options for listing admins
type AdminsListOptions struct {
	AccountID int64
}

// Validate validates the options
func (o *AdminsListOptions) Validate() error {
	// AccountID is resolved by resolveAccountID, so no validation needed here
	return nil
}

// AdminsAddOptions contains options for adding an admin
type AdminsAddOptions struct {
	AccountID        int64
	UserID           int64
	Role             string
	RoleID           int64
	SendConfirmation bool
	// Track which fields were set
	SendConfirmationSet bool
}

// Validate validates the options
func (o *AdminsAddOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("user-id", o.UserID)
}

// AdminsRemoveOptions contains options for removing an admin
type AdminsRemoveOptions struct {
	AccountID int64
	UserID    int64
	RoleID    int64
	// Track which fields were set
	RoleIDSet bool
}

// Validate validates the options
func (o *AdminsRemoveOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("user-id", o.UserID)
}

// RolesListOptions contains options for listing roles
type RolesListOptions struct {
	AccountID     int64
	State         string
	ShowInherited bool
}

// Validate validates the options
func (o *RolesListOptions) Validate() error {
	// AccountID is resolved by resolveAccountID, so no validation needed here
	return nil
}

// RolesGetOptions contains options for getting a role
type RolesGetOptions struct {
	AccountID int64
	RoleID    int64
}

// Validate validates the options
func (o *RolesGetOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("role-id", o.RoleID)
}

// RolesCreateOptions contains options for creating a role
type RolesCreateOptions struct {
	AccountID    int64
	Label        string
	BaseRoleType string
}

// Validate validates the options
func (o *RolesCreateOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	if o.Label == "" {
		return fmt.Errorf("label is required")
	}
	return nil
}

// RolesUpdateOptions contains options for updating a role
type RolesUpdateOptions struct {
	AccountID int64
	RoleID    int64
	Label     string
	// Track which fields were set
	LabelSet bool
}

// Validate validates the options
func (o *RolesUpdateOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("role-id", o.RoleID)
}

// RolesDeactivateOptions contains options for deactivating a role
type RolesDeactivateOptions struct {
	AccountID int64
	RoleID    int64
}

// Validate validates the options
func (o *RolesDeactivateOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("role-id", o.RoleID)
}

// RolesActivateOptions contains options for activating a role
type RolesActivateOptions struct {
	AccountID int64
	RoleID    int64
}

// Validate validates the options
func (o *RolesActivateOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("role-id", o.RoleID)
}
