package options

// SISImportsListOptions contains options for listing SIS imports
type SISImportsListOptions struct {
	AccountID     int64
	WorkflowState string
	CreatedSince  string
	CreatedBefore string
}

// Validate validates the options
func (o *SISImportsListOptions) Validate() error {
	// Account ID is optional - will use default if not specified
	return nil
}

// SISImportsGetOptions contains options for getting a SIS import
type SISImportsGetOptions struct {
	AccountID int64
	ImportID  int64
}

// Validate validates the options
func (o *SISImportsGetOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("import-id", o.ImportID)
}

// SISImportsCreateOptions contains options for creating a SIS import
type SISImportsCreateOptions struct {
	AccountID          int64
	FilePath           string
	ImportType         string
	Extension          string
	BatchMode          bool
	BatchModeTermID    int64
	OverrideStickiness bool
	AddStickiness      bool
	ClearStickiness    bool
	DiffingID          string
	DiffingRemaster    bool
	ChangeThreshold    float64
	// Track which fields were set
	BatchModeSet          bool
	BatchModeTermIDSet    bool
	OverrideStickinessSet bool
	AddStickinessSet      bool
	ClearStickinessSet    bool
	DiffingRemasterSet    bool
	ChangeThresholdSet    bool
}

// Validate validates the options
func (o *SISImportsCreateOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("file", o.FilePath)
}

// SISImportsAbortOptions contains options for aborting a SIS import
type SISImportsAbortOptions struct {
	AccountID int64
	ImportID  int64
}

// Validate validates the options
func (o *SISImportsAbortOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("import-id", o.ImportID)
}

// SISImportsRestoreOptions contains options for restoring SIS import states
type SISImportsRestoreOptions struct {
	AccountID int64
	ImportID  int64
	BatchMode bool
}

// Validate validates the options
func (o *SISImportsRestoreOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("import-id", o.ImportID)
}

// SISImportsErrorsOptions contains options for listing SIS import errors
type SISImportsErrorsOptions struct {
	AccountID int64
	ImportID  int64
}

// Validate validates the options
func (o *SISImportsErrorsOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("import-id", o.ImportID)
}
