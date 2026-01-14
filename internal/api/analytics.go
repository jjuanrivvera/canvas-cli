package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// AnalyticsService handles analytics-related API calls
type AnalyticsService struct {
	client *Client
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(client *Client) *AnalyticsService {
	return &AnalyticsService{client: client}
}

// CourseActivity represents course participation data over time
type CourseActivity struct {
	Date           string `json:"date"`
	Views          int    `json:"views"`
	Participations int    `json:"participations"`
}

// AssignmentAnalytics represents assignment statistics
type AssignmentAnalytics struct {
	AssignmentID         int64           `json:"assignment_id"`
	Title                string          `json:"title"`
	DueAt                string          `json:"due_at,omitempty"`
	UnlockedAt           string          `json:"unlock_at,omitempty"`
	PointsPossible       float64         `json:"points_possible"`
	NonDigitalSubmission bool            `json:"non_digital_submission"`
	Muted                bool            `json:"muted"`
	MinScore             float64         `json:"min_score"`
	MaxScore             float64         `json:"max_score"`
	MedianScore          float64         `json:"median"`
	FirstQuartile        float64         `json:"first_quartile"`
	ThirdQuartile        float64         `json:"third_quartile"`
	Tardiness            *TardinessStats `json:"tardiness,omitempty"`
}

// TardinessStats represents tardiness breakdown
type TardinessStats struct {
	Missing  int `json:"missing"`
	Late     int `json:"late"`
	OnTime   int `json:"on_time"`
	Floating int `json:"floating"`
	Total    int `json:"total"`
}

// StudentSummary represents a student's course summary
type StudentSummary struct {
	ID                  int64           `json:"id"`
	PageViews           int             `json:"page_views"`
	MaxPageViews        int             `json:"max_page_views,omitempty"`
	PageViewsLevel      int             `json:"page_views_level,omitempty"`
	Participations      int             `json:"participations"`
	MaxParticipations   int             `json:"max_participations,omitempty"`
	ParticipationsLevel int             `json:"participations_level,omitempty"`
	Tardiness           *TardinessStats `json:"tardiness,omitempty"`
	CurrentScore        float64         `json:"current_score,omitempty"`
	FinalScore          float64         `json:"final_score,omitempty"`
	CurrentGrade        string          `json:"current_grade,omitempty"`
	FinalGrade          string          `json:"final_grade,omitempty"`
}

// UserActivity represents a user's participation in a course
type UserActivity struct {
	Date           string `json:"date"`
	Views          int    `json:"views"`
	Participations int    `json:"participations"`
}

// UserAssignmentAnalytics represents a user's assignment data
type UserAssignmentAnalytics struct {
	AssignmentID   int64   `json:"assignment_id"`
	Title          string  `json:"title"`
	DueAt          string  `json:"due_at,omitempty"`
	UnlockedAt     string  `json:"unlock_at,omitempty"`
	PointsPossible float64 `json:"points_possible"`
	Muted          bool    `json:"muted"`
	Score          float64 `json:"score,omitempty"`
	Submission     *struct {
		SubmittedAt string `json:"submitted_at,omitempty"`
	} `json:"submission,omitempty"`
}

// UserCommunication represents a user's messaging data
type UserCommunication struct {
	InstructorMessages int `json:"instructorMessages"`
	StudentMessages    int `json:"studentMessages"`
}

// DepartmentStatistics represents department-level statistics
type DepartmentStatistics struct {
	Subaccounts       int `json:"subaccounts"`
	Teachers          int `json:"teachers"`
	Students          int `json:"students"`
	DiscussionTopics  int `json:"discussion_topics"`
	DiscussionReplies int `json:"discussion_replies"`
	MediaObjects      int `json:"media_objects"`
	Attachments       int `json:"attachments"`
	Assignments       int `json:"assignments"`
}

// DepartmentActivity represents department participation data
type DepartmentActivity struct {
	Date           string `json:"date"`
	Views          int    `json:"views"`
	Participations int    `json:"participations"`
}

// DepartmentGrades represents department grade distribution
type DepartmentGrades struct {
	Score float64 `json:"score"`
	Count int     `json:"count"`
}

// GetCourseActivity retrieves course participation over time
func (s *AnalyticsService) GetCourseActivity(ctx context.Context, courseID int64) ([]CourseActivity, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/analytics/activity", courseID)

	var activity []CourseActivity
	if err := s.client.GetJSON(ctx, path, &activity); err != nil {
		return nil, err
	}

	return activity, nil
}

// GetCourseAssignments retrieves assignment analytics for a course
func (s *AnalyticsService) GetCourseAssignments(ctx context.Context, courseID int64) ([]AssignmentAnalytics, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/analytics/assignments", courseID)

	var assignments []AssignmentAnalytics
	if err := s.client.GetJSON(ctx, path, &assignments); err != nil {
		return nil, err
	}

	return assignments, nil
}

// ListStudentSummariesOptions holds options for listing student summaries
type ListStudentSummariesOptions struct {
	SortColumn string // name, name_descending, score, score_descending, participations, etc.
	StudentID  int64  // Filter to single student
	Page       int
	PerPage    int
}

// GetStudentSummaries retrieves student summary analytics for a course
func (s *AnalyticsService) GetStudentSummaries(ctx context.Context, courseID int64, opts *ListStudentSummariesOptions) ([]StudentSummary, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/analytics/student_summaries", courseID)

	if opts != nil {
		query := url.Values{}

		if opts.SortColumn != "" {
			query.Add("sort_column", opts.SortColumn)
		}

		if opts.StudentID > 0 {
			query.Add("student_id", strconv.FormatInt(opts.StudentID, 10))
		}

		if opts.Page > 0 {
			query.Add("page", strconv.Itoa(opts.Page))
		}

		if opts.PerPage > 0 {
			query.Add("per_page", strconv.Itoa(opts.PerPage))
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var summaries []StudentSummary
	if err := s.client.GetAllPages(ctx, path, &summaries); err != nil {
		return nil, err
	}

	return summaries, nil
}

// GetUserActivity retrieves a user's participation in a course
func (s *AnalyticsService) GetUserActivity(ctx context.Context, courseID, userID int64) ([]UserActivity, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/analytics/users/%d/activity", courseID, userID)

	var activity []UserActivity
	if err := s.client.GetJSON(ctx, path, &activity); err != nil {
		return nil, err
	}

	return activity, nil
}

// GetUserAssignments retrieves a user's assignment data for a course
func (s *AnalyticsService) GetUserAssignments(ctx context.Context, courseID, userID int64) ([]UserAssignmentAnalytics, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/analytics/users/%d/assignments", courseID, userID)

	var assignments []UserAssignmentAnalytics
	if err := s.client.GetJSON(ctx, path, &assignments); err != nil {
		return nil, err
	}

	return assignments, nil
}

// GetUserCommunication retrieves a user's messaging data for a course
func (s *AnalyticsService) GetUserCommunication(ctx context.Context, courseID, userID int64) (*UserCommunication, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/analytics/users/%d/communication", courseID, userID)

	var communication UserCommunication
	if err := s.client.GetJSON(ctx, path, &communication); err != nil {
		return nil, err
	}

	return &communication, nil
}

// DepartmentAnalyticsOptions holds options for department analytics
type DepartmentAnalyticsOptions struct {
	TermID    int64  // Filter by term
	StartDate string // ISO8601 date
	EndDate   string // ISO8601 date
}

// GetDepartmentStatistics retrieves department-level statistics
func (s *AnalyticsService) GetDepartmentStatistics(ctx context.Context, accountID int64, opts *DepartmentAnalyticsOptions) (*DepartmentStatistics, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/analytics/current/statistics", accountID)
	path = s.addDepartmentOptions(path, opts)

	var stats DepartmentStatistics
	if err := s.client.GetJSON(ctx, path, &stats); err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetDepartmentActivity retrieves department participation data
func (s *AnalyticsService) GetDepartmentActivity(ctx context.Context, accountID int64, opts *DepartmentAnalyticsOptions) ([]DepartmentActivity, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/analytics/current/activity", accountID)
	path = s.addDepartmentOptions(path, opts)

	var activity []DepartmentActivity
	if err := s.client.GetJSON(ctx, path, &activity); err != nil {
		return nil, err
	}

	return activity, nil
}

// GetDepartmentGrades retrieves department grade distribution
func (s *AnalyticsService) GetDepartmentGrades(ctx context.Context, accountID int64, opts *DepartmentAnalyticsOptions) ([]DepartmentGrades, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/analytics/current/grades", accountID)
	path = s.addDepartmentOptions(path, opts)

	var grades []DepartmentGrades
	if err := s.client.GetJSON(ctx, path, &grades); err != nil {
		return nil, err
	}

	return grades, nil
}

func (s *AnalyticsService) addDepartmentOptions(path string, opts *DepartmentAnalyticsOptions) string {
	if opts == nil {
		return path
	}

	query := url.Values{}

	if opts.TermID > 0 {
		query.Add("term_id", strconv.FormatInt(opts.TermID, 10))
	}

	if opts.StartDate != "" {
		query.Add("start_date", opts.StartDate)
	}

	if opts.EndDate != "" {
		query.Add("end_date", opts.EndDate)
	}

	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	return path
}
