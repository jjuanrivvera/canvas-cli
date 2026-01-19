package options

// AccountsListOptions contains options for listing accounts
type AccountsListOptions struct {
	Include []string
}

// Validate validates the options
func (o *AccountsListOptions) Validate() error {
	// No required fields for list
	return nil
}

// AccountsGetOptions contains options for getting an account
type AccountsGetOptions struct {
	AccountID int64
}

// Validate validates the options
func (o *AccountsGetOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}

// AccountsSubAccountsOptions contains options for listing sub-accounts
type AccountsSubAccountsOptions struct {
	AccountID int64
	Recursive bool
}

// Validate validates the options
func (o *AccountsSubAccountsOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}
