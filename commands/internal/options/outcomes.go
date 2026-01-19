package options

import "fmt"

// OutcomesGetOptions contains options for getting an outcome
type OutcomesGetOptions struct {
	OutcomeID int64
}

// Validate validates the options
func (o *OutcomesGetOptions) Validate() error {
	if o.OutcomeID <= 0 {
		return fmt.Errorf("outcome-id is required and must be greater than 0")
	}
	return nil
}

// OutcomesCreateOptions contains options for creating an outcome
type OutcomesCreateOptions struct {
	CourseID          int64
	AccountID         int64
	GroupID           int64
	Title             string
	DisplayName       string
	Description       string
	MasteryPoints     float64
	CalculationMethod string
	CalculationInt    int
}

// Validate validates the options
func (o *OutcomesCreateOptions) Validate() error {
	if o.CourseID == 0 && o.AccountID == 0 {
		return fmt.Errorf("must specify either --course-id or --account-id")
	}
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	if o.Title == "" {
		return fmt.Errorf("title is required")
	}
	return nil
}

// OutcomesUpdateOptions contains options for updating an outcome
type OutcomesUpdateOptions struct {
	OutcomeID         int64
	Title             string
	DisplayName       string
	Description       string
	MasteryPoints     float64
	CalculationMethod string
	CalculationInt    int
	// Track which fields were actually set
	TitleSet             bool
	DisplayNameSet       bool
	DescriptionSet       bool
	MasteryPointsSet     bool
	CalculationMethodSet bool
	CalculationIntSet    bool
}

// Validate validates the options
func (o *OutcomesUpdateOptions) Validate() error {
	if o.OutcomeID <= 0 {
		return fmt.Errorf("outcome-id is required and must be greater than 0")
	}
	// At least one field must be set for update
	if !o.TitleSet && !o.DisplayNameSet && !o.DescriptionSet &&
		!o.MasteryPointsSet && !o.CalculationMethodSet && !o.CalculationIntSet {
		return fmt.Errorf("at least one field must be specified for update")
	}
	return nil
}

// OutcomesListOptions contains options for listing outcomes in a group
type OutcomesListOptions struct {
	CourseID  int64
	AccountID int64
	GroupID   int64
}

// Validate validates the options
func (o *OutcomesListOptions) Validate() error {
	if o.CourseID == 0 && o.AccountID == 0 {
		return fmt.Errorf("must specify either --course-id or --account-id")
	}
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	return nil
}

// OutcomesLinkOptions contains options for linking an outcome to a group
type OutcomesLinkOptions struct {
	CourseID  int64
	AccountID int64
	GroupID   int64
	OutcomeID int64
}

// Validate validates the options
func (o *OutcomesLinkOptions) Validate() error {
	if o.CourseID == 0 && o.AccountID == 0 {
		return fmt.Errorf("must specify either --course-id or --account-id")
	}
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	if o.OutcomeID <= 0 {
		return fmt.Errorf("outcome-id is required and must be greater than 0")
	}
	return nil
}

// OutcomesUnlinkOptions contains options for unlinking an outcome from a group
type OutcomesUnlinkOptions struct {
	CourseID  int64
	AccountID int64
	GroupID   int64
	OutcomeID int64
	Force     bool
}

// Validate validates the options
func (o *OutcomesUnlinkOptions) Validate() error {
	if o.CourseID == 0 && o.AccountID == 0 {
		return fmt.Errorf("must specify either --course-id or --account-id")
	}
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	if o.OutcomeID <= 0 {
		return fmt.Errorf("outcome-id is required and must be greater than 0")
	}
	return nil
}

// OutcomesGroupsListOptions contains options for listing outcome groups
type OutcomesGroupsListOptions struct {
	CourseID  int64
	AccountID int64
}

// Validate validates the options
func (o *OutcomesGroupsListOptions) Validate() error {
	// Either course or account can be specified, but not required
	// Will use default account if neither is provided
	return nil
}

// OutcomesGroupsGetOptions contains options for getting an outcome group
type OutcomesGroupsGetOptions struct {
	CourseID  int64
	AccountID int64
	GroupID   int64
}

// Validate validates the options
func (o *OutcomesGroupsGetOptions) Validate() error {
	if o.CourseID == 0 && o.AccountID == 0 {
		return fmt.Errorf("must specify either --course-id or --account-id")
	}
	if o.GroupID <= 0 {
		return fmt.Errorf("group-id is required and must be greater than 0")
	}
	return nil
}

// OutcomesResultsOptions contains options for getting outcome results
type OutcomesResultsOptions struct {
	CourseID      int64
	UserIDs       []int64
	OutcomeIDs    []int64
	Include       []string
	IncludeHidden bool
}

// Validate validates the options
func (o *OutcomesResultsOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	return nil
}

// OutcomesAlignmentsOptions contains options for getting outcome alignments
type OutcomesAlignmentsOptions struct {
	CourseID int64
}

// Validate validates the options
func (o *OutcomesAlignmentsOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	return nil
}
