package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestAnalyticsActivityCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get course activity successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": cmdtest.NewMockResponse(`{"id":1,"name":"Test Course"}`),
				"/api/v1/courses/1/analytics/activity": cmdtest.NewMockResponse(`[
					{"date":"2024-01-15","participations":10,"views":50}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "50") && !strings.Contains(output, "10") {
					t.Error("expected activity data in output")
				}
			},
		},
		{
			Name:        "missing course ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAnalyticsActivityCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAnalyticsAssignmentsCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get assignment analytics successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": cmdtest.NewMockResponse(`{"id":1,"name":"Test Course"}`),
				"/api/v1/courses/1/analytics/assignments": cmdtest.NewMockResponse(`[
					{"assignment_id":10,"title":"Homework 1","due_at":"2024-01-20","points_possible":100}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Homework") && !strings.Contains(output, "10") {
					t.Error("expected assignment data in output")
				}
			},
		},
		{
			Name:        "missing course ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAnalyticsAssignmentsCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAnalyticsStudentsCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get student summaries successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": cmdtest.NewMockResponse(`{"id":1,"name":"Test Course"}`),
				"/api/v1/courses/1/analytics/student_summaries": cmdtest.NewMockResponse(`[
					{"id":100,"page_views":20,"participations":5}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "100") {
					t.Error("expected student ID 100 in output")
				}
			},
		},
		{
			Name: "get student summaries with sort",
			Args: []string{"--course-id", "1", "--sort", "score"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                             cmdtest.NewMockResponse(`{"id":1,"name":"Test Course"}`),
				"/api/v1/courses/1/analytics/student_summaries": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No student data found",
		},
		{
			Name:        "missing course ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAnalyticsStudentsCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAnalyticsUserCmd_Assignments(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get user assignments analytics",
			Args: []string{"100", "--course-id", "1", "--type", "assignments"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": cmdtest.NewMockResponse(`{"id":1,"name":"Test Course"}`),
				"/api/v1/courses/1/analytics/users/100/assignments": cmdtest.NewMockResponse(`[
					{"assignment_id":10,"title":"HW1","points_possible":100,"score":90}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "10") {
					t.Error("expected assignment_id 10 in output")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAnalyticsUserCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAnalyticsUserCmd_Communication(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get user communication analytics",
			Args: []string{"100", "--course-id", "1", "--type", "communication"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": cmdtest.NewMockResponse(`{"id":1,"name":"Test Course"}`),
				"/api/v1/courses/1/analytics/users/100/communication": cmdtest.NewMockResponse(`{
					"student_id":100,"discussion_topics":{"total":3}
				}`),
			},
			ExpectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAnalyticsUserCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAnalyticsUserCmd_InvalidType(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "invalid analytics type",
			Args: []string{"100", "--course-id", "1", "--type", "bogus"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": cmdtest.NewMockResponse(`{"id":1,"name":"Test Course"}`),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAnalyticsUserCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAnalyticsDepartmentCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get department statistics",
			Args: []string{"--account-id", "1", "--type", "statistics"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/analytics/current/statistics": cmdtest.NewMockResponse(`{
					"courses":10,"teachers":5,"students":200
				}`),
			},
			ExpectError: false,
		},
		{
			Name: "get department activity",
			Args: []string{"--account-id", "1", "--type", "activity"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/analytics/current/activity": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No activity data found",
		},
		{
			Name: "get department grades",
			Args: []string{"--account-id", "1", "--type", "grades"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/analytics/current/grades": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No grade data found",
		},
		{
			Name: "invalid analytics type",
			Args: []string{"--account-id", "1", "--type", "bogus"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/analytics/current/statistics": cmdtest.NewMockResponse(`{}`),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAnalyticsDepartmentCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
