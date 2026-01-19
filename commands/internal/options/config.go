package options

import "fmt"

// ConfigListOptions contains options for listing instances
type ConfigListOptions struct {
	// No options needed
}

// Validate validates the options
func (o *ConfigListOptions) Validate() error {
	return nil
}

// ConfigAddOptions contains options for adding an instance
type ConfigAddOptions struct {
	Name        string
	URL         string
	Description string
	ClientID    string
}

// Validate validates the options
func (o *ConfigAddOptions) Validate() error {
	if o.Name == "" {
		return fmt.Errorf("instance name is required")
	}
	if o.URL == "" {
		return fmt.Errorf("url is required")
	}
	return nil
}

// ConfigUseOptions contains options for switching instances
type ConfigUseOptions struct {
	InstanceName string
}

// Validate validates the options
func (o *ConfigUseOptions) Validate() error {
	if o.InstanceName == "" {
		return fmt.Errorf("instance name is required")
	}
	return nil
}

// ConfigRemoveOptions contains options for removing an instance
type ConfigRemoveOptions struct {
	InstanceName string
	Force        bool
}

// Validate validates the options
func (o *ConfigRemoveOptions) Validate() error {
	if o.InstanceName == "" {
		return fmt.Errorf("instance name is required")
	}
	return nil
}

// ConfigShowOptions contains options for showing configuration
type ConfigShowOptions struct {
	// No options needed
}

// Validate validates the options
func (o *ConfigShowOptions) Validate() error {
	return nil
}

// ConfigAccountOptions contains options for setting account ID
type ConfigAccountOptions struct {
	InstanceName string
	AccountID    int64
	Detect       bool
}

// Validate validates the options
func (o *ConfigAccountOptions) Validate() error {
	// Instance name is optional, will use default if not provided
	// AccountID validation is done in the run function
	return nil
}
