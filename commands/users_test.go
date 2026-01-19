package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestUsersListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list users successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/users": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"name": "John Doe",
						"email": "john@example.com",
						"login_id": "john",
						"enrollments": []
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "John Doe") {
					t.Error("Expected 'John Doe' in output")
				}
			},
		},
		{
			Name: "list users - empty response",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":       courseMock,
				"/api/v1/courses/1/users": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No users found",
		},
		{
			Name:        "list users - missing course ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newUsersListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestUsersGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get user successfully",
			Args: []string{"10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/users/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"name": "Jane Smith",
					"email": "jane@example.com",
					"login_id": "jane"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Jane Smith") {
					t.Error("Expected 'Jane Smith' in output")
				}
			},
		},
		{
			Name:        "get user - missing user ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newUsersGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestUsersMeCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get current user successfully",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/users/self": cmdtest.NewMockResponse(`{
					"id": 100,
					"name": "Current User",
					"email": "current@example.com",
					"login_id": "current"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Current User") {
					t.Error("Expected 'Current User' in output")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newUsersMeCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
