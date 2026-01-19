package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestEnrollmentsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list enrollments successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/enrollments": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"user_id": 100,
						"course_id": 1,
						"type": "StudentEnrollment",
						"enrollment_state": "active",
						"role": "StudentEnrollment"
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "StudentEnrollment") {
					t.Error("Expected 'StudentEnrollment' in output")
				}
			},
		},
		{
			Name: "list enrollments - empty response",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":             courseMock,
				"/api/v1/courses/1/enrollments": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No enrollments found",
		},
		{
			Name:        "list enrollments - missing course ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newEnrollmentsListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestEnrollmentsCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create enrollment successfully",
			Args: []string{"--course-id", "1", "--user-id", "100", "--type", "StudentEnrollment"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/enrollments": cmdtest.NewMockResponse(`{
					"id": 20,
					"user_id": 100,
					"course_id": 1,
					"type": "StudentEnrollment",
					"enrollment_state": "invited"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "StudentEnrollment") {
					t.Error("Expected 'StudentEnrollment' in output")
				}
			},
		},
		{
			Name:        "create enrollment - missing course ID",
			Args:        []string{"--user-id", "100", "--type", "StudentEnrollment"},
			ExpectError: true,
		},
		{
			Name:        "create enrollment - missing user ID",
			Args:        []string{"--course-id", "1", "--type", "StudentEnrollment"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newEnrollmentsCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestEnrollmentsConcludeCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete enrollment with confirmation",
			Args: []string{"--course-id", "1", "10", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                courseMock,
				"/api/v1/courses/1/enrollments/10": cmdtest.NewMockResponse(`{}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete enrollment - missing course ID",
			Args:        []string{"10", "--force"},
			ExpectError: true,
		},
		{
			Name:        "delete enrollment - missing enrollment ID",
			Args:        []string{"--course-id", "1", "--force"},
			ExpectError: true,
		},
	}
	//
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newEnrollmentsConcludeCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
