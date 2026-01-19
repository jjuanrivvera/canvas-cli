package options

import "fmt"

// ModulesListOptions contains options for listing modules
type ModulesListOptions struct {
	CourseID   int64
	Include    []string
	SearchTerm string
	StudentID  string
}

// Validate validates the options
func (o *ModulesListOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	return nil
}

// ModulesGetOptions contains options for getting a module
type ModulesGetOptions struct {
	CourseID int64
	ModuleID int64
	Include  []string
}

// Validate validates the options
func (o *ModulesGetOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.ModuleID <= 0 {
		return fmt.Errorf("module-id is required and must be greater than 0")
	}
	return nil
}

// ModulesCreateOptions contains options for creating a module
type ModulesCreateOptions struct {
	CourseID                  int64
	Name                      string
	UnlockAt                  string
	Position                  int
	RequireSequentialProgress bool
	PrerequisiteModuleIDs     []int64
	PublishFinalGrade         bool
	Published                 bool
}

// Validate validates the options
func (o *ModulesCreateOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

// ModulesUpdateOptions contains options for updating a module
type ModulesUpdateOptions struct {
	CourseID                  int64
	ModuleID                  int64
	Name                      string
	UnlockAt                  string
	Position                  int
	RequireSequentialProgress bool
	PrerequisiteModuleIDs     []int64
	PublishFinalGrade         bool
	Published                 bool
	// Track which fields were actually set
	NameSet                      bool
	UnlockAtSet                  bool
	PositionSet                  bool
	RequireSequentialProgressSet bool
	PrerequisiteModuleIDsSet     bool
	PublishFinalGradeSet         bool
	PublishedSet                 bool
}

// Validate validates the options
func (o *ModulesUpdateOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.ModuleID <= 0 {
		return fmt.Errorf("module-id is required and must be greater than 0")
	}
	// At least one field must be set for update
	if !o.NameSet && !o.UnlockAtSet && !o.PositionSet && !o.RequireSequentialProgressSet &&
		!o.PrerequisiteModuleIDsSet && !o.PublishFinalGradeSet && !o.PublishedSet {
		return fmt.Errorf("at least one field must be specified for update")
	}
	return nil
}

// ModulesDeleteOptions contains options for deleting a module
type ModulesDeleteOptions struct {
	CourseID int64
	ModuleID int64
	Force    bool
}

// Validate validates the options
func (o *ModulesDeleteOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.ModuleID <= 0 {
		return fmt.Errorf("module-id is required and must be greater than 0")
	}
	return nil
}

// ModulesRelockOptions contains options for relocking module progressions
type ModulesRelockOptions struct {
	CourseID int64
	ModuleID int64
}

// Validate validates the options
func (o *ModulesRelockOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.ModuleID <= 0 {
		return fmt.Errorf("module-id is required and must be greater than 0")
	}
	return nil
}

// ModulesPublishOptions contains options for publishing a module
type ModulesPublishOptions struct {
	CourseID int64
	ModuleID int64
}

// Validate validates the options
func (o *ModulesPublishOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.ModuleID <= 0 {
		return fmt.Errorf("module-id is required and must be greater than 0")
	}
	return nil
}

// ModulesUnpublishOptions contains options for unpublishing a module
type ModulesUnpublishOptions struct {
	CourseID int64
	ModuleID int64
}

// Validate validates the options
func (o *ModulesUnpublishOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.ModuleID <= 0 {
		return fmt.Errorf("module-id is required and must be greater than 0")
	}
	return nil
}

// ModulesItemsListOptions contains options for listing module items
type ModulesItemsListOptions struct {
	CourseID   int64
	ModuleID   int64
	Include    []string
	SearchTerm string
	StudentID  string
}

// Validate validates the options
func (o *ModulesItemsListOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.ModuleID <= 0 {
		return fmt.Errorf("module-id is required and must be greater than 0")
	}
	return nil
}

// ModulesItemsGetOptions contains options for getting a module item
type ModulesItemsGetOptions struct {
	CourseID int64
	ModuleID int64
	ItemID   int64
	Include  []string
}

// Validate validates the options
func (o *ModulesItemsGetOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.ModuleID <= 0 {
		return fmt.Errorf("module-id is required and must be greater than 0")
	}
	if o.ItemID <= 0 {
		return fmt.Errorf("item-id is required and must be greater than 0")
	}
	return nil
}

// ModulesItemsCreateOptions contains options for creating a module item
type ModulesItemsCreateOptions struct {
	CourseID       int64
	ModuleID       int64
	Type           string
	Title          string
	ContentID      int64
	PageURL        string
	ExternalURL    string
	NewTab         bool
	Indent         int
	CompletionType string
	MinScore       float64
}

// Validate validates the options
func (o *ModulesItemsCreateOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.ModuleID <= 0 {
		return fmt.Errorf("module-id is required and must be greater than 0")
	}
	if o.Type == "" {
		return fmt.Errorf("type is required")
	}
	if o.Title == "" {
		return fmt.Errorf("title is required")
	}

	// Validate type-specific requirements
	validTypes := map[string]bool{
		"File": true, "Page": true, "Discussion": true, "Assignment": true,
		"Quiz": true, "SubHeader": true, "ExternalUrl": true, "ExternalTool": true,
	}
	if !validTypes[o.Type] {
		return fmt.Errorf("invalid type: %s (must be one of: File, Page, Discussion, Assignment, Quiz, SubHeader, ExternalUrl, ExternalTool)", o.Type)
	}

	return nil
}

// ModulesItemsUpdateOptions contains options for updating a module item
type ModulesItemsUpdateOptions struct {
	CourseID       int64
	ModuleID       int64
	ItemID         int64
	Title          string
	Position       int
	Indent         int
	NewTab         bool
	CompletionType string
	MinScore       float64
	Published      bool
	MoveToModule   int64
	// Track which fields were actually set
	TitleSet          bool
	PositionSet       bool
	IndentSet         bool
	NewTabSet         bool
	CompletionTypeSet bool
	MinScoreSet       bool
	PublishedSet      bool
	MoveToModuleSet   bool
}

// Validate validates the options
func (o *ModulesItemsUpdateOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.ModuleID <= 0 {
		return fmt.Errorf("module-id is required and must be greater than 0")
	}
	if o.ItemID <= 0 {
		return fmt.Errorf("item-id is required and must be greater than 0")
	}
	// At least one field must be set for update
	if !o.TitleSet && !o.PositionSet && !o.IndentSet && !o.NewTabSet &&
		!o.CompletionTypeSet && !o.MinScoreSet && !o.PublishedSet && !o.MoveToModuleSet {
		return fmt.Errorf("at least one field must be specified for update")
	}
	return nil
}

// ModulesItemsDeleteOptions contains options for deleting a module item
type ModulesItemsDeleteOptions struct {
	CourseID int64
	ModuleID int64
	ItemID   int64
	Force    bool
}

// Validate validates the options
func (o *ModulesItemsDeleteOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.ModuleID <= 0 {
		return fmt.Errorf("module-id is required and must be greater than 0")
	}
	if o.ItemID <= 0 {
		return fmt.Errorf("item-id is required and must be greater than 0")
	}
	return nil
}

// ModulesItemsDoneOptions contains options for marking a module item as done
type ModulesItemsDoneOptions struct {
	CourseID int64
	ModuleID int64
	ItemID   int64
}

// Validate validates the options
func (o *ModulesItemsDoneOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.ModuleID <= 0 {
		return fmt.Errorf("module-id is required and must be greater than 0")
	}
	if o.ItemID <= 0 {
		return fmt.Errorf("item-id is required and must be greater than 0")
	}
	return nil
}
