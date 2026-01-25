package options

import (
	"fmt"
	"strings"
)

// AliasSetOptions contains options for alias set command
type AliasSetOptions struct {
	Name      string
	Expansion string
}

// Validate validates the options
func (o *AliasSetOptions) Validate() error {
	if o.Name == "" {
		return fmt.Errorf("alias name is required")
	}
	if o.Expansion == "" {
		return fmt.Errorf("alias expansion is required")
	}
	if strings.Contains(o.Name, " ") {
		return fmt.Errorf("alias name cannot contain spaces")
	}
	return nil
}

// AliasListOptions contains options for alias list command
type AliasListOptions struct {
	// No options needed for list
}

// Validate validates the options
func (o *AliasListOptions) Validate() error {
	return nil
}

// AliasDeleteOptions contains options for alias delete command
type AliasDeleteOptions struct {
	Name string
}

// Validate validates the options
func (o *AliasDeleteOptions) Validate() error {
	if o.Name == "" {
		return fmt.Errorf("alias name is required")
	}
	return nil
}
