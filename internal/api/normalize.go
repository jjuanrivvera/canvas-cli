package api

// NormalizeCourse ensures consistent data structure for a course
func NormalizeCourse(course *Course) *Course {
	if course == nil {
		return nil
	}

	// Ensure permissions map exists
	if course.Permissions == nil {
		course.Permissions = make(map[string]bool)
	}

	// Ensure blueprint restrictions map exists
	if course.BlueprintRestrictions == nil {
		course.BlueprintRestrictions = make(map[string]bool)
	}

	if course.BlueprintRestrictionsByObjectType == nil {
		course.BlueprintRestrictionsByObjectType = make(map[string]map[string]bool)
	}

	return course
}

// NormalizeUser ensures consistent data structure for a user
func NormalizeUser(user *User) *User {
	if user == nil {
		return nil
	}

	// Ensure enrollments slice is not nil
	if user.Enrollments == nil {
		user.Enrollments = []Enrollment{}
	}

	return user
}

// NormalizeAssignment ensures consistent data structure for an assignment
func NormalizeAssignment(assignment *Assignment) *Assignment {
	if assignment == nil {
		return nil
	}

	// Ensure slices are not nil
	if assignment.AllowedExtensions == nil {
		assignment.AllowedExtensions = []string{}
	}

	if assignment.SubmissionTypes == nil {
		assignment.SubmissionTypes = []string{}
	}

	if assignment.NeedsGradingCountBySection == nil {
		assignment.NeedsGradingCountBySection = []SectionGradingCount{}
	}

	if assignment.FrozenAttributes == nil {
		assignment.FrozenAttributes = []string{}
	}

	if assignment.Rubric == nil {
		assignment.Rubric = []RubricCriterion{}
	}

	if assignment.AssignmentVisibility == nil {
		assignment.AssignmentVisibility = []int64{}
	}

	if assignment.Overrides == nil {
		assignment.Overrides = []AssignmentOverride{}
	}

	// Ensure maps are not nil
	if assignment.TurnitinSettings == nil {
		assignment.TurnitinSettings = make(map[string]interface{})
	}

	if assignment.ExternalToolTagAttributes == nil {
		assignment.ExternalToolTagAttributes = make(map[string]interface{})
	}

	if assignment.IntegrationData == nil {
		assignment.IntegrationData = make(map[string]interface{})
	}

	if assignment.RubricSettings == nil {
		assignment.RubricSettings = make(map[string]interface{})
	}

	return assignment
}

// NormalizeSubmission ensures consistent data structure for a submission
func NormalizeSubmission(submission *Submission) *Submission {
	if submission == nil {
		return nil
	}

	// Ensure slices are not nil
	if submission.Attachments == nil {
		submission.Attachments = []Attachment{}
	}

	if submission.SubmissionComments == nil {
		submission.SubmissionComments = []SubmissionComment{}
	}

	if submission.Rubric == nil {
		submission.Rubric = []RubricAssessment{}
	}

	return submission
}

// NormalizeTerm ensures consistent data structure for a term
func NormalizeTerm(term *Term) *Term {
	if term == nil {
		return nil
	}

	// Ensure overrides map exists
	if term.Overrides == nil {
		term.Overrides = make(map[string]TermOverride)
	}

	return term
}

// NormalizeCourses normalizes a slice of courses
func NormalizeCourses(courses []Course) []Course {
	if courses == nil {
		return []Course{}
	}

	for i := range courses {
		NormalizeCourse(&courses[i])
	}

	return courses
}

// NormalizeUsers normalizes a slice of users
func NormalizeUsers(users []User) []User {
	if users == nil {
		return []User{}
	}

	for i := range users {
		NormalizeUser(&users[i])
	}

	return users
}

// NormalizeAssignments normalizes a slice of assignments
func NormalizeAssignments(assignments []Assignment) []Assignment {
	if assignments == nil {
		return []Assignment{}
	}

	for i := range assignments {
		NormalizeAssignment(&assignments[i])
	}

	return assignments
}

// NormalizeSubmissions normalizes a slice of submissions
func NormalizeSubmissions(submissions []Submission) []Submission {
	if submissions == nil {
		return []Submission{}
	}

	for i := range submissions {
		NormalizeSubmission(&submissions[i])
	}

	return submissions
}

// NormalizeEnrollment ensures consistent data structure for an enrollment
func NormalizeEnrollment(enrollment *Enrollment) *Enrollment {
	if enrollment == nil {
		return nil
	}

	// Ensure nested objects are initialized if present
	if enrollment.Grades == nil {
		enrollment.Grades = &Grades{}
	}

	return enrollment
}

// NormalizeEnrollments normalizes a slice of enrollments
func NormalizeEnrollments(enrollments []Enrollment) []Enrollment {
	if enrollments == nil {
		return []Enrollment{}
	}

	for i := range enrollments {
		NormalizeEnrollment(&enrollments[i])
	}

	return enrollments
}
