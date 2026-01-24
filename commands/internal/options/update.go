package options

// UpdateOptions holds options for the update command
type UpdateOptions struct{}

// Validate validates the update options
func (o *UpdateOptions) Validate() error {
	return nil
}

// UpdateCheckOptions holds options for the update check command
type UpdateCheckOptions struct{}

// Validate validates the update check options
func (o *UpdateCheckOptions) Validate() error {
	return nil
}

// UpdateEnableOptions holds options for the update enable command
type UpdateEnableOptions struct {
	Interval int
}

// Validate validates the update enable options
func (o *UpdateEnableOptions) Validate() error {
	if o.Interval < 1 {
		o.Interval = 60 // Default to 60 minutes
	}
	return nil
}

// UpdateDisableOptions holds options for the update disable command
type UpdateDisableOptions struct{}

// Validate validates the update disable options
func (o *UpdateDisableOptions) Validate() error {
	return nil
}

// UpdateStatusOptions holds options for the update status command
type UpdateStatusOptions struct{}

// Validate validates the update status options
func (o *UpdateStatusOptions) Validate() error {
	return nil
}
