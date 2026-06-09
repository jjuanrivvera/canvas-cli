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

func TestEnrollmentsGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get enrollment successfully",
			Args: []string{"--course-id", "1", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				// Canvas has no direct "get enrollment by ID" endpoint, so the command
				// lists all course enrollments and filters by ID
				"/api/v1/courses/1/enrollments": cmdtest.NewMockResponse(`[
					{
						"id": 10,
						"user_id": 100,
						"course_id": 1,
						"type": "StudentEnrollment",
						"enrollment_state": "active"
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
			Name:        "get enrollment - missing course ID",
			Args:        []string{"10"},
			ExpectError: true,
		},
		{
			Name:        "get enrollment - missing enrollment ID",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newEnrollmentsGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestEnrollmentsReactivateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "reactivate enrollment successfully",
			Args: []string{"--course-id", "1", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/enrollments/10/reactivate": cmdtest.NewMockResponse(`{
					"id": 10,
					"enrollment_state": "active"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "reactivate - missing course ID",
			Args:        []string{"10"},
			ExpectError: true,
		},
		{
			Name:        "reactivate - missing enrollment ID",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newEnrollmentsReactivateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestEnrollmentsAcceptCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "accept enrollment successfully",
			Args: []string{"--course-id", "1", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/enrollments/10/accept": cmdtest.NewMockResponse(`{
					"success": true
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "accept - missing course ID",
			Args:        []string{"10"},
			ExpectError: true,
		},
		{
			Name:        "accept - missing enrollment ID",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newEnrollmentsAcceptCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestEnrollmentsRejectCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "reject enrollment successfully",
			Args: []string{"--course-id", "1", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/enrollments/10/reject": cmdtest.NewMockResponse(`{
					"success": true
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "reject - missing course ID",
			Args:        []string{"10"},
			ExpectError: true,
		},
		{
			Name:        "reject - missing enrollment ID",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newEnrollmentsRejectCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestEnrollmentsListCmd_UserContext(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list enrollments for user",
		Args: []string{"--user-id", "100"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/100/enrollments": cmdtest.NewMockResponse(`[
				{
					"id": 1,
					"user_id": 100,
					"course_id": 5,
					"type": "StudentEnrollment",
					"enrollment_state": "active"
				}
			]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "StudentEnrollment") {
				t.Error("Expected 'StudentEnrollment' in output")
			}
		},
	}

	cmd := newEnrollmentsListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}
