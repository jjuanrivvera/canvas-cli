package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestAssignmentsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list assignments successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/assignments": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"name": "Assignment 1",
						"due_at": "2024-12-31T23:59:59Z",
						"points_possible": 100,
						"published": true
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Assignment 1") {
					t.Error("Expected 'Assignment 1' in output")
				}
			},
		},
		{
			Name: "list assignments - empty response",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":             courseMock,
				"/api/v1/courses/1/assignments": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No assignments found",
		},
		{
			Name:        "list assignments - missing course ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAssignmentsListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAssignmentsGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get assignment successfully",
			Args: []string{"--course-id", "1", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/assignments/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"name": "Test Assignment",
					"due_at": "2024-12-31T23:59:59Z",
					"points_possible": 100,
					"published": true
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Test Assignment") {
					t.Error("Expected 'Test Assignment' in output")
				}
			},
		},
		{
			Name:        "get assignment - missing course ID",
			Args:        []string{"10"},
			ExpectError: true,
		},
		{
			Name:        "get assignment - missing assignment ID",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAssignmentsGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAssignmentsCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create assignment successfully",
			Args: []string{"--course-id", "1", "--name", "New Assignment"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/assignments": cmdtest.NewMockResponse(`{
					"id": 20,
					"name": "New Assignment",
					"due_at": null,
					"points_possible": 0,
					"published": false
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "New Assignment") {
					t.Error("Expected 'New Assignment' in output")
				}
			},
		},
		{
			Name:        "create assignment - missing course ID",
			Args:        []string{"--name", "New Assignment"},
			ExpectError: true,
		},
		{
			Name:        "create assignment - missing name",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAssignmentsCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAssignmentsDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete assignment with confirmation",
			Args: []string{"--course-id", "1", "10", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/assignments/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"name": "Deleted Assignment"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete assignment - missing course ID",
			Args:        []string{"10", "--force"},
			ExpectError: true,
		},
		{
			Name:        "delete assignment - missing assignment ID",
			Args:        []string{"--course-id", "1", "--force"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAssignmentsDeleteCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
