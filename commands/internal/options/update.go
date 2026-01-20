package options

// UpdateCheckOptions contains options for the update check command
type UpdateCheckOptions struct {
	// Force a check even if cache is valid
	Force bool
	// Output format (table, json, yaml)
	OutputFormat string
}

// Validate validates the options
func (o *UpdateCheckOptions) Validate() error {
	if o.OutputFormat != "" {
		validFormats := []string{"table", "json", "yaml"}
		valid := false
		for _, format := range validFormats {
			if o.OutputFormat == format {
				valid = true
				break
			}
		}
		if !valid {
			return ErrInvalidValue("output-format", o.OutputFormat, validFormats...)
		}
	}
	return nil
}

// UpdateInstallOptions contains options for the update install command
type UpdateInstallOptions struct {
	// Skip confirmation prompt
	Yes bool
}

// Validate validates the options
func (o *UpdateInstallOptions) Validate() error {
	return nil
}

// UpdateEnableOptions contains options for the update enable command
type UpdateEnableOptions struct {
	// Interval in hours for automatic update checks (0 to disable)
	Interval int
}

// Validate validates the options
func (o *UpdateEnableOptions) Validate() error {
	if o.Interval < 0 {
		return ErrInvalidValue("interval", "", "positive number or 0")
	}
	return nil
}

// UpdateDisableOptions contains options for the update disable command
type UpdateDisableOptions struct {
}

// Validate validates the options
func (o *UpdateDisableOptions) Validate() error {
	return nil
}
