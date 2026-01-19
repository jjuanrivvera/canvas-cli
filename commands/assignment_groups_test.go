package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestAssignmentGroupsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list assignment groups successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/assignment_groups": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"name": "Homework",
						"position": 1,
						"group_weight": 25
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Homework") {
					t.Error("Expected 'Homework' in output")
				}
			},
		},
		{
			Name: "list assignment groups - empty response",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/assignment_groups": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No assignment groups found",
		},
		{
			Name:        "list assignment groups - missing course ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAssignmentGroupsListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAssignmentGroupsGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get assignment group successfully",
			Args: []string{"10", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/assignment_groups/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"name": "Exams",
					"position": 2,
					"group_weight": 40
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Exams") {
					t.Error("Expected 'Exams' in output")
				}
			},
		},
		{
			Name:        "get assignment group - missing group ID",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "get assignment group - missing course ID",
			Args:        []string{"10"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAssignmentGroupsGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAssignmentGroupsCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create assignment group successfully",
			Args: []string{"--course-id", "1", "--name", "Projects"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/assignment_groups": cmdtest.NewMockResponse(`{
					"id": 20,
					"name": "Projects",
					"course_id": 1,
					"group_weight": 0
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Projects") {
					t.Error("Expected 'Projects' in output")
				}
			},
		},
		{
			Name:        "create assignment group - missing course ID",
			Args:        []string{"--name", "Projects"},
			ExpectError: true,
		},
		{
			Name:        "create assignment group - missing name",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAssignmentGroupsCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAssignmentGroupsUpdateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "update assignment group successfully",
			Args: []string{"10", "--course-id", "1", "--name", "Updated Group"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/assignment_groups/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"name": "Updated Group",
					"course_id": 1
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Updated Group") {
					t.Error("Expected 'Updated Group' in output")
				}
			},
		},
		{
			Name:        "update assignment group - missing group ID",
			Args:        []string{"--course-id", "1", "--name", "Updated"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAssignmentGroupsUpdateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAssignmentGroupsDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete assignment group with confirmation",
			Args: []string{"10", "--course-id", "1", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/assignment_groups/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"name": "Old Group"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete assignment group - missing group ID",
			Args:        []string{"--course-id", "1", "--force"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAssignmentGroupsDeleteCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
