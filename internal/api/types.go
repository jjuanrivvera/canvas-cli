package api

import (
	"time"
)

// Account represents a Canvas account (institution, sub-account)
type Account struct {
	ID                     int64  `json:"id"`
	Name                   string `json:"name"`
	UUID                   string `json:"uuid"`
	ParentAccountID        int64  `json:"parent_account_id"`
	RootAccountID          int64  `json:"root_account_id"`
	DefaultStorageQuotaMB  int64  `json:"default_storage_quota_mb"`
	DefaultUserStorageQuotaMB int64 `json:"default_user_storage_quota_mb"`
	DefaultGroupStorageQuotaMB int64 `json:"default_group_storage_quota_mb"`
	DefaultTimeZone        string `json:"default_time_zone"`
	SISAccountID           string `json:"sis_account_id"`
	IntegrationID          string `json:"integration_id"`
	SISImportID            int64  `json:"sis_import_id"`
	LTIGuid                string `json:"lti_guid"`
	WorkflowState          string `json:"workflow_state"`
}

// Course represents a Canvas course
type Course struct {
	ID                       int64     `json:"id"`
	Name                     string    `json:"name"`
	CourseCode               string    `json:"course_code"`
	WorkflowState            string    `json:"workflow_state"`
	AccountID                int64     `json:"account_id"`
	StartAt                  time.Time `json:"start_at"`
	EndAt                    time.Time `json:"end_at"`
	EnrollmentTermID         int64     `json:"enrollment_term_id"`
	GradingStandardID        int64     `json:"grading_standard_id"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
	DefaultView              string    `json:"default_view"`
	SyllabusBody             string    `json:"syllabus_body"`
	NeedsGradingCount        int       `json:"needs_grading_count"`
	Term                     *Term     `json:"term,omitempty"`
	CourseProgress           *Progress `json:"course_progress,omitempty"`
	ApplyAssignmentGroupWeights bool   `json:"apply_assignment_group_weights"`
	Permissions              map[string]bool `json:"permissions,omitempty"`
	IsPublic                 bool      `json:"is_public"`
	IsPublicToAuthUsers      bool      `json:"is_public_to_auth_users"`
	PublicSyllabus           bool      `json:"public_syllabus"`
	PublicSyllabusToAuth     bool      `json:"public_syllabus_to_auth"`
	PublicDescription        string    `json:"public_description"`
	StorageQuotaMB           int       `json:"storage_quota_mb"`
	StorageQuotaUsedMB       float64   `json:"storage_quota_used_mb"`
	HideFinalGrades          bool      `json:"hide_final_grades"`
	License                  string    `json:"license"`
	AllowStudentAssignmentEdits bool   `json:"allow_student_assignment_edits"`
	AllowWikiComments        bool      `json:"allow_wiki_comments"`
	AllowStudentForumAttachments bool  `json:"allow_student_forum_attachments"`
	OpenEnrollment           bool      `json:"open_enrollment"`
	SelfEnrollment           bool      `json:"self_enrollment"`
	RestrictEnrollmentsToCourseDates bool `json:"restrict_enrollments_to_course_dates"`
	CourseFormat             string    `json:"course_format"`
	AccessRestrictedByDate   bool      `json:"access_restricted_by_date"`
	TimeZone                 string    `json:"time_zone"`
	Blueprint                bool      `json:"blueprint"`
	BlueprintRestrictions    map[string]bool `json:"blueprint_restrictions,omitempty"`
	BlueprintRestrictionsByObjectType map[string]map[string]bool `json:"blueprint_restrictions_by_object_type,omitempty"`
}

// User represents a Canvas user
type User struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	SortableName   string `json:"sortable_name"`
	ShortName      string `json:"short_name"`
	SisUserID      string `json:"sis_user_id"`
	SisImportID    int64  `json:"sis_import_id"`
	IntegrationID  string `json:"integration_id"`
	LoginID        string `json:"login_id"`
	AvatarURL      string `json:"avatar_url"`
	Enrollments    []Enrollment `json:"enrollments,omitempty"`
	Email          string `json:"email"`
	Locale         string `json:"locale"`
	LastLogin      time.Time `json:"last_login"`
	TimeZone       string `json:"time_zone"`
	Bio            string `json:"bio"`
}

// Assignment represents a Canvas assignment
type Assignment struct {
	ID                      int64     `json:"id"`
	Name                    string    `json:"name"`
	Description             string    `json:"description"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
	DueAt                   time.Time `json:"due_at"`
	LockAt                  time.Time `json:"lock_at"`
	UnlockAt                time.Time `json:"unlock_at"`
	HasOverrides            bool      `json:"has_overrides"`
	CourseID                int64     `json:"course_id"`
	HTMLURL                 string    `json:"html_url"`
	SubmissionsDownloadURL  string    `json:"submissions_download_url"`
	AssignmentGroupID       int64     `json:"assignment_group_id"`
	DueDateRequired         bool      `json:"due_date_required"`
	AllowedExtensions       []string  `json:"allowed_extensions"`
	MaxNameLength           int       `json:"max_name_length"`
	TurnitinEnabled         bool      `json:"turnitin_enabled"`
	VericiteEnabled         bool      `json:"vericite_enabled"`
	TurnitinSettings        map[string]interface{} `json:"turnitin_settings,omitempty"`
	GradeGroupStudentsIndividually bool `json:"grade_group_students_individually"`
	ExternalToolTagAttributes map[string]interface{} `json:"external_tool_tag_attributes,omitempty"`
	PeerReviews             bool      `json:"peer_reviews"`
	AutomaticPeerReviews    bool      `json:"automatic_peer_reviews"`
	PeerReviewCount         int       `json:"peer_review_count"`
	PeerReviewsAssignAt     time.Time `json:"peer_reviews_assign_at"`
	IntraGroupPeerReviews   bool      `json:"intra_group_peer_reviews"`
	GroupCategoryID         int64     `json:"group_category_id"`
	NeedsGradingCount       int       `json:"needs_grading_count"`
	NeedsGradingCountBySection []SectionGradingCount `json:"needs_grading_count_by_section,omitempty"`
	Position                int       `json:"position"`
	PostToSIS               bool      `json:"post_to_sis"`
	IntegrationID           string    `json:"integration_id"`
	IntegrationData         map[string]interface{} `json:"integration_data,omitempty"`
	PointsPossible          float64   `json:"points_possible"`
	SubmissionTypes         []string  `json:"submission_types"`
	HasSubmittedSubmissions bool      `json:"has_submitted_submissions"`
	GradingType             string    `json:"grading_type"`
	GradingStandardID       int64     `json:"grading_standard_id"`
	Published               bool      `json:"published"`
	Unpublishable           bool      `json:"unpublishable"`
	OnlyVisibleToOverrides  bool      `json:"only_visible_to_overrides"`
	LockedForUser           bool      `json:"locked_for_user"`
	LockInfo                *LockInfo `json:"lock_info,omitempty"`
	LockExplanation         string    `json:"lock_explanation"`
	QuizID                  int64     `json:"quiz_id"`
	AnonymousInstructorAnnotations bool `json:"anonymous_instructor_annotations"`
	AnonymousPeerReviews    bool      `json:"anonymous_peer_reviews"`
	AnonymousMarking        bool      `json:"anonymous_marking"`
	AnonymousGrading        bool      `json:"anonymous_grading"`
	GradersAnonymousToGraders bool    `json:"graders_anonymous_to_graders"`
	GraderCount             int       `json:"grader_count"`
	GraderCommentsVisibleToGraders bool `json:"grader_comments_visible_to_graders"`
	FinalGraderID           int64     `json:"final_grader_id"`
	GraderNamesVisibleToFinalGrader bool `json:"grader_names_visible_to_final_grader"`
	AllowedAttempts         int       `json:"allowed_attempts"`
	AnnotatableAttachmentID int64     `json:"annotatable_attachment_id"`
	HideInGradebook         bool      `json:"hide_in_gradebook"`
	SecureParams            string    `json:"secure_params"`
	LTIContextID            string    `json:"lti_context_id"`
	CourseID2               int64     `json:"course_id"`
	NameHash                string    `json:"name_hash,omitempty"`
	CanDuplicate            bool      `json:"can_duplicate"`
	OriginalCourseID        int64     `json:"original_course_id"`
	OriginalAssignmentID    int64     `json:"original_assignment_id"`
	OriginalLTIResourceLinkID string  `json:"original_lti_resource_link_id"`
	OriginalAssignmentName  string    `json:"original_assignment_name"`
	OriginalQuizID          int64     `json:"original_quiz_id"`
	WorkflowState           string    `json:"workflow_state"`
	ImportantDates          bool      `json:"important_dates"`
	MutedTLN                bool      `json:"muted"`
	HTMLURL2                string    `json:"html_url"`
	HasGradableSubmissions  bool      `json:"has_gradable_submissions"`
	URL                     string    `json:"url,omitempty"`
	IsQuizAssignment        bool      `json:"is_quiz_assignment"`
	CanUpdate               bool      `json:"can_update"`
	Frozen                  bool      `json:"frozen"`
	FrozenAttributes        []string  `json:"frozen_attributes,omitempty"`
	Submission              *Submission `json:"submission,omitempty"`
	UseRubricForGrading     bool      `json:"use_rubric_for_grading"`
	RubricSettings          map[string]interface{} `json:"rubric_settings,omitempty"`
	Rubric                  []RubricCriterion `json:"rubric,omitempty"`
	AssignmentVisibility    []int64   `json:"assignment_visibility,omitempty"`
	Overrides               []AssignmentOverride `json:"overrides,omitempty"`
	OmitFromFinalGrade      bool      `json:"omit_from_final_grade"`
	ModeratedGrading        bool      `json:"moderated_grading"`
	GraderCommentsVisibleToGraders2 bool `json:"grader_comments_visible_to_graders"`
	FinalGraderID2          int64     `json:"final_grader_id"`
	GraderNamesVisibleToFinalGrader2 bool `json:"grader_names_visible_to_final_grader"`
	AllowedAttempts2        int       `json:"allowed_attempts"`
}

// Submission represents a Canvas submission
type Submission struct {
	ID                 int64     `json:"id"`
	Body               string    `json:"body"`
	URL                string    `json:"url"`
	Grade              string    `json:"grade"`
	Score              float64   `json:"score"`
	SubmittedAt        time.Time `json:"submitted_at"`
	AssignmentID       int64     `json:"assignment_id"`
	UserID             int64     `json:"user_id"`
	SubmissionType     string    `json:"submission_type"`
	WorkflowState      string    `json:"workflow_state"`
	GradeMatchesCurrentSubmission bool `json:"grade_matches_current_submission"`
	GradedAt           time.Time `json:"graded_at"`
	GraderID           int64     `json:"grader_id"`
	Attempt            int       `json:"attempt"`
	CachedDueDate      time.Time `json:"cached_due_date"`
	ExcusedTLN         bool      `json:"excused"`
	LatePolicyStatus   string    `json:"late_policy_status"`
	PointsDeducted     float64   `json:"points_deducted"`
	GradingPeriodID    int64     `json:"grading_period_id"`
	ExtraAttempts      int       `json:"extra_attempts"`
	PostedAt           time.Time `json:"posted_at"`
	Late               bool      `json:"late"`
	Missing            bool      `json:"missing"`
	SecondsLate        int       `json:"seconds_late"`
	EnteredGrade       string    `json:"entered_grade"`
	EnteredScore       float64   `json:"entered_score"`
	PreviewURL         string    `json:"preview_url"`
	AnonymousID        string    `json:"anonymous_id"`
	User               *User     `json:"user,omitempty"`
	Attachments        []Attachment `json:"attachments,omitempty"`
	SubmissionComments []SubmissionComment `json:"submission_comments,omitempty"`
	Assignment         *Assignment `json:"assignment,omitempty"`
	Course             *Course   `json:"course,omitempty"`
	Rubric             []RubricAssessment `json:"rubric_assessment,omitempty"`
}

// Enrollment represents a Canvas enrollment
type Enrollment struct {
	ID                         int64     `json:"id"`
	CourseID                   int64     `json:"course_id"`
	CourseSectionID            int64     `json:"course_section_id"`
	EnrollmentState            string    `json:"enrollment_state"`
	LimitPrivilegesToCourseSection bool `json:"limit_privileges_to_course_section"`
	RootAccountID              int64     `json:"root_account_id"`
	Type                       string    `json:"type"`
	UserID                     int64     `json:"user_id"`
	AssociatedUserID           int64     `json:"associated_user_id"`
	Role                       string    `json:"role"`
	RoleID                     int64     `json:"role_id"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
	StartAt                    time.Time `json:"start_at"`
	EndAt                      time.Time `json:"end_at"`
	LastActivityAt             time.Time `json:"last_activity_at"`
	LastAttendedAt             time.Time `json:"last_attended_at"`
	TotalActivityTime          int       `json:"total_activity_time"`
	HTMLURL                    string    `json:"html_url"`
	Grades                     *Grades   `json:"grades,omitempty"`
	User                       *User     `json:"user,omitempty"`
	OverrideGrade              string    `json:"override_grade"`
	OverrideScore              float64   `json:"override_score"`
	UnpostedCurrentGrade       string    `json:"unposted_current_grade"`
	UnpostedCurrentScore       float64   `json:"unposted_current_score"`
	UnpostedFinalGrade         string    `json:"unposted_final_grade"`
	UnpostedFinalScore         float64   `json:"unposted_final_score"`
	HasGradingPeriods          bool      `json:"has_grading_periods"`
	TotalsForAllGradingPeriodsOption bool `json:"totals_for_all_grading_periods_option"`
	CurrentGradingPeriodTitle  string    `json:"current_grading_period_title"`
	CurrentGradingPeriodID     int64     `json:"current_grading_period_id"`
	CurrentPeriodOverrideGrade string    `json:"current_period_override_grade"`
	CurrentPeriodOverrideScore float64   `json:"current_period_override_score"`
	CurrentPeriodUnpostedCurrentGrade string `json:"current_period_unposted_current_grade"`
	CurrentPeriodUnpostedCurrentScore float64 `json:"current_period_unposted_current_score"`
	CurrentPeriodUnpostedFinalGrade string `json:"current_period_unposted_final_grade"`
	CurrentPeriodUnpostedFinalScore float64 `json:"current_period_unposted_final_score"`
}

// Term represents a Canvas enrollment term
type Term struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	StartAt   time.Time `json:"start_at"`
	EndAt     time.Time `json:"end_at"`
	CreatedAt time.Time `json:"created_at"`
	WorkflowState string `json:"workflow_state"`
	GradingPeriodGroupID int64 `json:"grading_period_group_id"`
	SISTermID string    `json:"sis_term_id"`
	SISImportID int64   `json:"sis_import_id"`
	Overrides map[string]TermOverride `json:"overrides,omitempty"`
}

// TermOverride represents date overrides for a term
type TermOverride struct {
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
}

// Progress represents course progress
type Progress struct {
	RequirementCount          int     `json:"requirement_count"`
	RequirementCompletedCount int     `json:"requirement_completed_count"`
	NextRequirementURL        string  `json:"next_requirement_url"`
	CompletedAt               time.Time `json:"completed_at"`
}

// Grades represents enrollment grades
type Grades struct {
	HTMLURL          string  `json:"html_url"`
	CurrentGrade     string  `json:"current_grade"`
	CurrentScore     float64 `json:"current_score"`
	FinalGrade       string  `json:"final_grade"`
	FinalScore       float64 `json:"final_score"`
	UnpostedCurrentGrade string `json:"unposted_current_grade"`
	UnpostedCurrentScore float64 `json:"unposted_current_score"`
	UnpostedFinalGrade string `json:"unposted_final_grade"`
	UnpostedFinalScore float64 `json:"unposted_final_score"`
}

// Attachment represents a file attachment
type Attachment struct {
	ID          int64     `json:"id"`
	UUID        string    `json:"uuid"`
	FolderID    int64     `json:"folder_id"`
	DisplayName string    `json:"display_name"`
	Filename    string    `json:"filename"`
	ContentType string    `json:"content-type"`
	URL         string    `json:"url"`
	Size        int64     `json:"size"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UnlockAt    time.Time `json:"unlock_at"`
	Locked      bool      `json:"locked"`
	Hidden      bool      `json:"hidden"`
	LockAt      time.Time `json:"lock_at"`
	HiddenForUser bool    `json:"hidden_for_user"`
	ThumbnailURL string   `json:"thumbnail_url"`
	ModifiedAt  time.Time `json:"modified_at"`
	MIMEClass   string    `json:"mime_class"`
	MediaEntryID string   `json:"media_entry_id"`
	LockedForUser bool    `json:"locked_for_user"`
	LockInfo    *LockInfo `json:"lock_info,omitempty"`
	LockExplanation string `json:"lock_explanation"`
	PreviewURL  string    `json:"preview_url"`
}

// SubmissionComment represents a comment on a submission
type SubmissionComment struct {
	ID              int64     `json:"id"`
	AuthorID        int64     `json:"author_id"`
	AuthorName      string    `json:"author_name"`
	Author          *User     `json:"author,omitempty"`
	Comment         string    `json:"comment"`
	CreatedAt       time.Time `json:"created_at"`
	EditedAt        time.Time `json:"edited_at"`
	MediaComment    *MediaComment `json:"media_comment,omitempty"`
	Attachments     []Attachment  `json:"attachments,omitempty"`
}

// MediaComment represents a media comment
type MediaComment struct {
	ContentType string `json:"content-type"`
	DisplayName string `json:"display_name"`
	MediaID     string `json:"media_id"`
	MediaType   string `json:"media_type"`
	URL         string `json:"url"`
}

// LockInfo represents lock information
type LockInfo struct {
	AssetString      string    `json:"asset_string"`
	UnlockAt         time.Time `json:"unlock_at"`
	LockAt           time.Time `json:"lock_at"`
	ContextModule    string    `json:"context_module"`
	ManuallyLocked   bool      `json:"manually_locked"`
}

// RubricCriterion represents a rubric criterion
type RubricCriterion struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	LongDescription string `json:"long_description"`
	Points      float64 `json:"points"`
	CriterionUseRange bool `json:"criterion_use_range"`
	Ratings     []RubricRating `json:"ratings"`
}

// RubricRating represents a rubric rating
type RubricRating struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	LongDescription string `json:"long_description"`
	Points      float64 `json:"points"`
}

// RubricAssessment represents a rubric assessment
type RubricAssessment struct {
	CriterionID string  `json:"criterion_id"`
	Points      float64 `json:"points"`
	Comments    string  `json:"comments"`
	RatingID    string  `json:"rating_id"`
}

// AssignmentOverride represents an assignment override
type AssignmentOverride struct {
	ID           int64     `json:"id"`
	AssignmentID int64     `json:"assignment_id"`
	StudentIDs   []int64   `json:"student_ids,omitempty"`
	GroupID      int64     `json:"group_id,omitempty"`
	CourseSectionID int64  `json:"course_section_id,omitempty"`
	Title        string    `json:"title"`
	DueAt        time.Time `json:"due_at"`
	AllDay       bool      `json:"all_day"`
	AllDayDate   string    `json:"all_day_date"`
	UnlockAt     time.Time `json:"unlock_at"`
	LockAt       time.Time `json:"lock_at"`
}

// SectionGradingCount represents grading count by section
type SectionGradingCount struct {
	SectionID         int64 `json:"section_id"`
	NeedsGradingCount int   `json:"needs_grading_count"`
}

// PaginationLinks represents pagination links from Link header
type PaginationLinks struct {
	Current string
	Next    string
	Prev    string
	First   string
	Last    string
}

// APIError represents an error from the Canvas API
type APIError struct {
	StatusCode    int           `json:"-"`
	Errors        []ErrorDetail `json:"errors"`
	ErrorReportID int64         `json:"error_report_id,omitempty"`
	Suggestion    string        `json:"-"`
	DocsURL       string        `json:"-"`
}

// ErrorDetail represents detailed error information
type ErrorDetail struct {
	Message string `json:"message"`
	ErrorCode string `json:"error_code,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	if len(e.Errors) > 0 {
		msg := e.Errors[0].Message
		if e.Suggestion != "" {
			msg += "\n\nSuggestion: " + e.Suggestion
		}
		if e.DocsURL != "" {
			msg += "\nDocs: " + e.DocsURL
		}
		return msg
	}
	return "Unknown API error"
}

// RateLimitInfo represents rate limit information from response headers
type RateLimitInfo struct {
	Limit     float64
	Remaining float64
	Reset     time.Time
}
