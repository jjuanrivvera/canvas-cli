package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

// courseMock provides a mock response for course validation
var courseMock = cmdtest.NewMockResponse(`{
	"id": 1,
	"name": "Test Course",
	"course_code": "TEST101",
	"workflow_state": "available"
}`)

func TestModulesListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list modules successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts":  cmdtest.NewMockResponse(`[]`),
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/modules": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"name": "Module 1",
						"position": 1,
						"published": true
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Module 1") {
					t.Error("Expected 'Module 1' in output")
				}
			},
		},
		{
			Name: "list modules - empty response",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":         courseMock,
				"/api/v1/courses/1/modules": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No modules found",
		},
		{
			Name:        "list modules - missing course ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newModulesListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestModulesGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get module successfully",
			Args: []string{"--course-id", "1", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/modules/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"name": "Test Module",
					"position": 1,
					"published": true
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Test Module") {
					t.Error("Expected 'Test Module' in output")
				}
			},
		},
		{
			Name:        "get module - missing course ID",
			Args:        []string{"10"},
			ExpectError: true,
		},
		{
			Name:        "get module - missing module ID",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newModulesGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestModulesCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create module successfully",
			Args: []string{"--course-id", "1", "--name", "New Module"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/modules": cmdtest.NewMockResponse(`{
					"id": 20,
					"name": "New Module",
					"position": 2,
					"published": false
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "New Module") {
					t.Error("Expected 'New Module' in output")
				}
			},
		},
		{
			Name:        "create module - missing course ID",
			Args:        []string{"--name", "New Module"},
			ExpectError: true,
		},
		{
			Name:        "create module - missing name",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newModulesCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestModulesUpdateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "update module name",
			Args: []string{"--course-id", "1", "10", "--name", "Updated Module"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/modules/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"name": "Updated Module",
					"position": 1,
					"published": false
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Updated Module") {
					t.Error("Expected 'Updated Module' in output")
				}
			},
		},
		{
			Name:        "update module - missing course ID",
			Args:        []string{"10", "--name", "Updated"},
			ExpectError: true,
		},
		{
			Name:        "update module - missing module ID",
			Args:        []string{"--course-id", "1", "--name", "Updated"},
			ExpectError: true,
		},
		{
			Name: "update module - no changes",
			Args: []string{"--course-id", "1", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newModulesUpdateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestModulesDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete module with confirmation",
			Args: []string{"--course-id", "1", "10", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/modules/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"name": "Deleted Module"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete module - missing course ID",
			Args:        []string{"10", "--force"},
			ExpectError: true,
		},
		{
			Name:        "delete module - missing module ID",
			Args:        []string{"--course-id", "1", "--force"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newModulesDeleteCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestModulesRelockCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "relock module successfully",
			Args: []string{"--course-id", "1", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/modules/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"name": "Relocked Module",
					"position": 1
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "relock module - missing course ID",
			Args:        []string{"10"},
			ExpectError: true,
		},
		{
			Name:        "relock module - missing module ID",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newModulesRelockCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestModulesPublishCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "publish module successfully",
			Args: []string{"--course-id", "1", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/modules/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"name": "Published Module",
					"published": true
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Published Module") {
					t.Error("Expected 'Published Module' in output")
				}
			},
		},
		{
			Name:        "publish module - missing course ID",
			Args:        []string{"10"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newModulesPublishCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestModulesUnpublishCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "unpublish module successfully",
			Args: []string{"--course-id", "1", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/modules/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"name": "Unpublished Module",
					"published": false
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "unpublish module - missing course ID",
			Args:        []string{"10"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newModulesUnpublishCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestModulesItemsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list module items successfully",
			Args: []string{"--course-id", "1", "--module-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/modules/10/items": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"title": "Introduction Page",
						"position": 1,
						"type": "Page",
						"published": true
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Introduction Page") {
					t.Error("Expected 'Introduction Page' in output")
				}
			},
		},
		{
			Name: "list module items - empty response",
			Args: []string{"--course-id", "1", "--module-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts":                   cmdtest.NewMockResponse(`[]`),
				"/api/v1/courses/1":                  courseMock,
				"/api/v1/courses/1/modules/10/items": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No items found in this module",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newModulesItemsListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestModulesItemsGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get module item successfully",
			Args: []string{"--course-id", "1", "--module-id", "10", "789"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/modules/10/items/789": cmdtest.NewMockResponse(`{
					"id": 789,
					"title": "Quiz 1",
					"type": "Quiz",
					"published": true
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Quiz 1") {
					t.Error("Expected 'Quiz 1' in output")
				}
			},
		},
		{
			Name:        "get module item - missing course ID",
			Args:        []string{"--module-id", "10", "789"},
			ExpectError: true,
		},
		{
			Name:        "get module item - missing module ID",
			Args:        []string{"--course-id", "1", "789"},
			ExpectError: true,
		},
		{
			Name:        "get module item - missing item ID",
			Args:        []string{"--course-id", "1", "--module-id", "10"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newModulesItemsGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestModulesItemsCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create module item successfully",
			Args: []string{"--course-id", "1", "--module-id", "10", "--type", "Assignment", "--content-id", "99", "--title", "Assignment Item"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/modules/10/items": cmdtest.NewMockResponse(`{
					"id": 200,
					"title": "Assignment Item",
					"type": "Assignment",
					"content_id": 99
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Assignment Item") {
					t.Error("Expected 'Assignment Item' in output")
				}
			},
		},
		{
			Name:        "create module item - missing course ID",
			Args:        []string{"--module-id", "10", "--type", "Assignment", "--title", "My Item"},
			ExpectError: true,
		},
		{
			Name:        "create module item - missing module ID",
			Args:        []string{"--course-id", "1", "--type", "Assignment", "--title", "My Item"},
			ExpectError: true,
		},
		{
			Name:        "create module item - missing type",
			Args:        []string{"--course-id", "1", "--module-id", "10", "--title", "My Item"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newModulesItemsCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestModulesItemsUpdateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "update module item successfully",
			Args: []string{"--course-id", "1", "--module-id", "10", "789", "--title", "Updated Item"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/modules/10/items/789": cmdtest.NewMockResponse(`{
					"id": 789,
					"title": "Updated Item",
					"type": "Page",
					"published": true
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Updated Item") {
					t.Error("Expected 'Updated Item' in output")
				}
			},
		},
		{
			Name:        "update module item - missing course ID",
			Args:        []string{"--module-id", "10", "789", "--title", "Updated"},
			ExpectError: true,
		},
		{
			Name:        "update module item - missing module ID",
			Args:        []string{"--course-id", "1", "789", "--title", "Updated"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newModulesItemsUpdateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestModulesItemsDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete module item successfully",
			Args: []string{"--course-id", "1", "--module-id", "10", "789", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/modules/10/items/789": cmdtest.NewMockResponse(`{
					"id": 789
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete module item - missing item ID",
			Args:        []string{"--course-id", "1", "--module-id", "10", "--force"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newModulesItemsDeleteCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestModulesItemsDoneCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "mark module item done successfully",
			Args: []string{"--course-id", "1", "--module-id", "10", "789"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                      courseMock,
				"/api/v1/courses/1/modules/10/items/789": cmdtest.NewMockResponse(`{}`),
			},
			ExpectError: false,
		},
		{
			Name:        "mark done - missing item ID",
			Args:        []string{"--course-id", "1", "--module-id", "10"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newModulesItemsDoneCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestModulesListCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list modules - API error",
		Args: []string{"--course-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":         courseMock,
			"/api/v1/courses/1/modules": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}

	cmd := newModulesListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}
