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
