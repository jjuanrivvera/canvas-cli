package options

import "fmt"

// ExternalToolsListOptions contains options for listing external tools
type ExternalToolsListOptions struct {
	CourseID      int64
	AccountID     int64
	Search        string
	Selectable    bool
	IncludeParent bool
}

// Validate validates the options
func (o *ExternalToolsListOptions) Validate() error {
	// At least one context should be specified (or default account will be used)
	return nil
}

// ExternalToolsGetOptions contains options for getting an external tool
type ExternalToolsGetOptions struct {
	CourseID  int64
	AccountID int64
	ToolID    int64
}

// Validate validates the options
func (o *ExternalToolsGetOptions) Validate() error {
	if o.CourseID == 0 && o.AccountID == 0 {
		return fmt.Errorf("either course-id or account-id is required")
	}
	return ValidateRequired("tool-id", o.ToolID)
}

// ExternalToolsCreateOptions contains options for creating an external tool
type ExternalToolsCreateOptions struct {
	CourseID     int64
	AccountID    int64
	Name         string
	URL          string
	Domain       string
	ConsumerKey  string
	SharedSecret string
	PrivacyLevel string
	Description  string
	ConfigType   string
	ConfigURL    string
	ConfigXML    string
	JSONFile     string
}

// Validate validates the options
func (o *ExternalToolsCreateOptions) Validate() error {
	if o.CourseID == 0 && o.AccountID == 0 {
		return fmt.Errorf("either course-id or account-id is required")
	}
	return nil
}

// ExternalToolsUpdateOptions contains options for updating an external tool
type ExternalToolsUpdateOptions struct {
	CourseID     int64
	AccountID    int64
	ToolID       int64
	Name         string
	URL          string
	Domain       string
	ConsumerKey  string
	SharedSecret string
	PrivacyLevel string
	Description  string
	ConfigType   string
	ConfigURL    string
	ConfigXML    string
	JSONFile     string
	// Track which fields were set
	NameSet         bool
	URLSet          bool
	DomainSet       bool
	ConsumerKeySet  bool
	SharedSecretSet bool
	PrivacyLevelSet bool
	DescriptionSet  bool
	ConfigTypeSet   bool
	ConfigURLSet    bool
	ConfigXMLSet    bool
}

// Validate validates the options
func (o *ExternalToolsUpdateOptions) Validate() error {
	if o.CourseID == 0 && o.AccountID == 0 {
		return fmt.Errorf("either course-id or account-id is required")
	}
	return ValidateRequired("tool-id", o.ToolID)
}

// ExternalToolsDeleteOptions contains options for deleting an external tool
type ExternalToolsDeleteOptions struct {
	CourseID  int64
	AccountID int64
	ToolID    int64
	Force     bool
}

// Validate validates the options
func (o *ExternalToolsDeleteOptions) Validate() error {
	if o.CourseID == 0 && o.AccountID == 0 {
		return fmt.Errorf("either course-id or account-id is required")
	}
	return ValidateRequired("tool-id", o.ToolID)
}

// ExternalToolsLaunchOptions contains options for launching an external tool
type ExternalToolsLaunchOptions struct {
	CourseID     int64
	ToolID       int64
	LaunchType   string
	AssignmentID int64
	ModuleItemID int64
}

// Validate validates the options
func (o *ExternalToolsLaunchOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return ValidateRequired("tool-id", o.ToolID)
}
