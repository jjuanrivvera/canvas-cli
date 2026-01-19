package options

import "fmt"

// AuthLoginOptions contains options for auth login command
type AuthLoginOptions struct {
	InstanceURL  string
	InstanceName string
	OAuthMode    string
	ClientID     string
	ClientSecret string
}

// Validate validates the options
func (o *AuthLoginOptions) Validate() error {
	if o.InstanceURL == "" && o.InstanceName == "" {
		return fmt.Errorf("either instance URL or --instance is required")
	}
	if o.OAuthMode != "" && o.OAuthMode != "auto" && o.OAuthMode != "local" && o.OAuthMode != "oob" {
		return fmt.Errorf("invalid OAuth mode: %s (must be auto, local, or oob)", o.OAuthMode)
	}
	return nil
}

// AuthLogoutOptions contains options for auth logout command
type AuthLogoutOptions struct {
	InstanceName string
}

// Validate validates the options
func (o *AuthLogoutOptions) Validate() error {
	// Instance name is optional (uses default if not provided)
	return nil
}

// AuthStatusOptions contains options for auth status command
type AuthStatusOptions struct {
	InstanceName string
}

// Validate validates the options
func (o *AuthStatusOptions) Validate() error {
	// Instance name is optional (shows all if not provided)
	return nil
}

// AuthTokenSetOptions contains options for auth token set command
type AuthTokenSetOptions struct {
	InstanceName string
	Token        string
	URL          string
}

// Validate validates the options
func (o *AuthTokenSetOptions) Validate() error {
	if o.InstanceName == "" {
		return fmt.Errorf("instance-name is required")
	}
	return nil
}

// AuthTokenRemoveOptions contains options for auth token remove command
type AuthTokenRemoveOptions struct {
	InstanceName string
}

// Validate validates the options
func (o *AuthTokenRemoveOptions) Validate() error {
	if o.InstanceName == "" {
		return fmt.Errorf("instance-name is required")
	}
	return nil
}
