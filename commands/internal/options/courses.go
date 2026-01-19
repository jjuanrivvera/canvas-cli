package options

// CoursesListOptions encapsulates all flags for courses list command
type CoursesListOptions struct {
	// User context flags
	EnrollmentType  string
	EnrollmentState string
	Include         []string
	State           []string

	// Account context flags (admin mode)
	AccountID  int64
	SearchTerm string
	Sort       string
	Order      string

	// Pagination
	PerPage int
}

// Validate performs option validation
func (o *CoursesListOptions) Validate() error {
	// No required fields for listing courses
	return nil
}

// CoursesGetOptions encapsulates all flags for courses get command
type CoursesGetOptions struct {
	CourseID int64
	Include  []string
}

// Validate performs option validation
func (o *CoursesGetOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return nil
}

// CoursesCreateOptions encapsulates all flags for courses create command
type CoursesCreateOptions struct {
	AccountID                 int64
	Name                      string
	CourseCode                string
	StartAt                   string
	EndAt                     string
	License                   string
	IsPublic                  bool
	IsPublicToAuthUsers       bool
	PublicSyllabus            bool
	PublicSyllabusToAuth      bool
	PublicDescription         string
	AllowStudentWikiEdits     bool
	AllowWikiComments         bool
	AllowStudentForumAttach   bool
	OpenEnrollment            bool
	SelfEnrollment            bool
	RestrictEnrollmentsToDate bool
	TermID                    int64
	SISCourseID               string
	IntegrationID             string
	HideFinalGrades           bool
	ApplyAssignmentGroupWts   bool
	TimeZone                  string
	Offer                     bool
	EnrollMe                  bool
	DefaultView               string
	SyllabusCourseBody        string
	GradingStandardID         int64
	CourseFormat              string
}

// Validate performs option validation
func (o *CoursesCreateOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	if err := ValidateRequired("name", o.Name); err != nil {
		return err
	}
	return nil
}

// CoursesUpdateOptions encapsulates all flags for courses update command
type CoursesUpdateOptions struct {
	CourseID    int64
	Name        string
	CourseCode  string
	StartAt     string
	EndAt       string
	License     string
	IsPublic    *bool // Pointer to differentiate between not set and false
	DefaultView string
}

// Validate performs option validation
func (o *CoursesUpdateOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	return nil
}

// CoursesDeleteOptions encapsulates all flags for courses delete command
type CoursesDeleteOptions struct {
	CourseID int64
	Event    string
	Force    bool
}

// Validate performs option validation
func (o *CoursesDeleteOptions) Validate() error {
	if err := ValidateRequired("course-id", o.CourseID); err != nil {
		return err
	}
	// Validate event type
	if o.Event != "conclude" && o.Event != "delete" {
		return ErrInvalidValue("event", o.Event, "conclude", "delete")
	}
	return nil
}
