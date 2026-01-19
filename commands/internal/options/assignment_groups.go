package options

// AssignmentGroupsListOptions contains options for listing assignment groups
type AssignmentGroupsListOptions struct {
	CourseID int64
	Include  []string
}

// Validate validates the options
func (o *AssignmentGroupsListOptions) Validate() error {
	return ValidateRequired("course-id", o.CourseID)
}

// AssignmentGroupsGetOptions contains options for getting an assignment group
type AssignmentGroupsGetOptions struct {
	CourseID int64
	GroupID  int64
	Include  []string
}

// Validate validates the options
func (o *AssignmentGroupsGetOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return ValidateRequired("group-id", o.GroupID)
}

// AssignmentGroupsCreateOptions contains options for creating an assignment group
type AssignmentGroupsCreateOptions struct {
	CourseID    int64
	Name        string
	Position    int
	Weight      float64
	DropLowest  int
	DropHighest int
}

// Validate validates the options
func (o *AssignmentGroupsCreateOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return ValidateRequired("name", o.Name)
}

// AssignmentGroupsUpdateOptions contains options for updating an assignment group
type AssignmentGroupsUpdateOptions struct {
	CourseID    int64
	GroupID     int64
	Name        string
	Position    int
	Weight      float64
	DropLowest  int
	DropHighest int
	// Track which fields were set
	NameSet        bool
	PositionSet    bool
	WeightSet      bool
	DropLowestSet  bool
	DropHighestSet bool
}

// Validate validates the options
func (o *AssignmentGroupsUpdateOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return ValidateRequired("group-id", o.GroupID)
}

// AssignmentGroupsDeleteOptions contains options for deleting an assignment group
type AssignmentGroupsDeleteOptions struct {
	CourseID int64
	GroupID  int64
	Force    bool
	MoveTo   int64
}

// Validate validates the options
func (o *AssignmentGroupsDeleteOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return ValidateRequired("group-id", o.GroupID)
}
