package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestGroupsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list groups for course successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/groups": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"name": "Study Group",
						"description": "Group for studying",
						"members_count": 5
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Study Group") {
					t.Error("Expected 'Study Group' in output")
				}
			},
		},
		{
			Name: "list groups - empty response",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":        courseMock,
				"/api/v1/courses/1/groups": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No groups found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newGroupsListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestGroupsGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get group successfully",
			Args: []string{"10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/groups/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"name": "Project Team",
					"description": "Final project team",
					"members_count": 4
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Project Team") {
					t.Error("Expected 'Project Team' in output")
				}
			},
		},
		{
			Name:        "get group - missing group ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newGroupsGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestGroupsCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create group successfully",
			Args: []string{"--category-id", "5", "--name", "New Group"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/group_categories/5/groups": cmdtest.NewMockResponse(`{
					"id": 20,
					"name": "New Group",
					"group_category_id": 5
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "New Group") {
					t.Error("Expected 'New Group' in output")
				}
			},
		},
		{
			Name:        "create group - missing category ID",
			Args:        []string{"--name", "New Group"},
			ExpectError: true,
		},
		{
			Name:        "create group - missing name",
			Args:        []string{"--category-id", "5"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newGroupsCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestGroupsUpdateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "update group successfully",
			Args: []string{"10", "--name", "Updated Group"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/groups/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"name": "Updated Group"
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
			Name:        "update group - missing group ID",
			Args:        []string{"--name", "Updated"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newGroupsUpdateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestGroupsDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete group with confirmation",
			Args: []string{"10", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/groups/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"name": "Old Group"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete group - missing group ID",
			Args:        []string{"--force"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newGroupsDeleteCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestGroupsCategoriesListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list group categories successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/group_categories": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"name": "Project Groups",
						"self_signup": "enabled"
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Project Groups") {
					t.Error("Expected 'Project Groups' in output")
				}
			},
		},
		{
			Name: "list group categories - empty response",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                  courseMock,
				"/api/v1/courses/1/group_categories": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No group categories found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newGroupsCategoriesListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestGroupsCategoriesGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get group category successfully",
			Args: []string{"5"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/group_categories/5": cmdtest.NewMockResponse(`{
					"id": 5,
					"name": "Study Groups",
					"self_signup": "restricted"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Study Groups") {
					t.Error("Expected 'Study Groups' in output")
				}
			},
		},
		{
			Name:        "get group category - missing category ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newGroupsCategoriesGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestGroupsCategoriesCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create group category successfully",
			Args: []string{"--course-id", "1", "--name", "New Category"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/group_categories": cmdtest.NewMockResponse(`{
					"id": 10,
					"name": "New Category",
					"course_id": 1
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "New Category") {
					t.Error("Expected 'New Category' in output")
				}
			},
		},
		{
			Name:        "create group category - missing course ID",
			Args:        []string{"--name", "New Category"},
			ExpectError: true,
		},
		{
			Name:        "create group category - missing name",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newGroupsCategoriesCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestGroupsCategoriesDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete group category with confirmation",
			Args: []string{"5", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/group_categories/5": cmdtest.NewMockResponse(`{
					"id": 5,
					"name": "Old Category"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete group category - missing category ID",
			Args:        []string{"--force"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newGroupsCategoriesDeleteCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
