package options

// DoctorOptions contains options for running diagnostics
type DoctorOptions struct {
	Verbose bool
	JSON    bool
}

// Validate validates the options
func (o *DoctorOptions) Validate() error {
	// No required fields - all options are optional
	return nil
}
