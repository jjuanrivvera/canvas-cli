package options

import "fmt"

// RubricsListOptions contains options for listing rubrics
type RubricsListOptions struct {
	CourseID  int64
	AccountID int64
	Include   []string
}

// Validate validates the options
func (o *RubricsListOptions) Validate() error {
	// No required fields - will use default account if neither is specified
	return nil
}

// RubricsGetOptions contains options for getting a rubric
type RubricsGetOptions struct {
	CourseID  int64
	AccountID int64
	RubricID  int64
	Include   []string
}

// Validate validates the options
func (o *RubricsGetOptions) Validate() error {
	if err := ValidateRequired("rubric-id", o.RubricID); err != nil {
		return err
	}
	if o.CourseID == 0 && o.AccountID == 0 {
		return fmt.Errorf("must specify either --course-id or --account-id")
	}
	return nil
}

// RubricsCreateOptions contains options for creating a rubric
type RubricsCreateOptions struct {
	CourseID                  int64
	Title                     string
	PointsPossible            float64
	FreeFormCriterionComments bool
	HideScoreTotal            bool
}

// Validate validates the options
func (o *RubricsCreateOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return ValidateRequired("title", o.Title)
}

// RubricsUpdateOptions contains options for updating a rubric
type RubricsUpdateOptions struct {
	CourseID                  int64
	RubricID                  int64
	Title                     string
	PointsPossible            float64
	FreeFormCriterionComments bool
	HideScoreTotal            bool
	// Track which fields were set
	TitleSet                     bool
	PointsPossibleSet            bool
	FreeFormCriterionCommentsSet bool
	HideScoreTotalSet            bool
}

// Validate validates the options
func (o *RubricsUpdateOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return ValidateRequired("rubric-id", o.RubricID)
}

// RubricsDeleteOptions contains options for deleting a rubric
type RubricsDeleteOptions struct {
	CourseID int64
	RubricID int64
	Force    bool
}

// Validate validates the options
func (o *RubricsDeleteOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return ValidateRequired("rubric-id", o.RubricID)
}

// RubricsAssociateOptions contains options for associating a rubric with an assignment
type RubricsAssociateOptions struct {
	CourseID       int64
	RubricID       int64
	AssignmentID   int64
	UseForGrading  bool
	HideScoreTotal bool
	HidePoints     bool
}

// Validate validates the options
func (o *RubricsAssociateOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	if err := ValidateRequired("rubric-id", o.RubricID); err != nil {
		return err
	}
	return ValidateRequired("assignment-id", o.AssignmentID)
}
