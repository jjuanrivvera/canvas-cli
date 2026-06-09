package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestModulesGetCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get module - API error",
		Args: []string{"--course-id", "1", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":            courseMock,
			"/api/v1/courses/1/modules/10": cmdtest.NewErrorResponse(404, "module not found"),
		},
		ExpectError: true,
	}
	cmd := newModulesGetCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestModulesCreateCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create module - API error",
		Args: []string{"--course-id", "1", "--name", "Bad Module"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":         courseMock,
			"/api/v1/courses/1/modules": cmdtest.NewErrorResponse(422, "invalid module"),
		},
		ExpectError: true,
	}
	cmd := newModulesCreateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestModulesUpdateCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update module - API error",
		Args: []string{"--course-id", "1", "10", "--name", "Updated"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":            courseMock,
			"/api/v1/courses/1/modules/10": cmdtest.NewErrorResponse(404, "module not found"),
		},
		ExpectError: true,
	}
	cmd := newModulesUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestModulesDeleteCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete module - API error",
		Args: []string{"--course-id", "1", "10", "--force"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":            courseMock,
			"/api/v1/courses/1/modules/10": cmdtest.NewErrorResponse(404, "module not found"),
		},
		ExpectError: true,
	}
	cmd := newModulesDeleteCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestModulesRelockCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "relock module - API error",
		Args: []string{"--course-id", "1", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                   courseMock,
			"/api/v1/courses/1/modules/10/relock": cmdtest.NewErrorResponse(404, "module not found"),
		},
		ExpectError: true,
	}
	cmd := newModulesRelockCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestModulesPublishCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "publish module - API error",
		Args: []string{"--course-id", "1", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":            courseMock,
			"/api/v1/courses/1/modules/10": cmdtest.NewErrorResponse(404, "module not found"),
		},
		ExpectError: true,
	}
	cmd := newModulesPublishCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestModulesUnpublishCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "unpublish module - API error",
		Args: []string{"--course-id", "1", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":            courseMock,
			"/api/v1/courses/1/modules/10": cmdtest.NewErrorResponse(404, "module not found"),
		},
		ExpectError: true,
	}
	cmd := newModulesUnpublishCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestModulesItemsListCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list module items - API error",
		Args: []string{"--course-id", "1", "--module-id", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                  courseMock,
			"/api/v1/courses/1/modules/10/items": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmd := newModulesItemsListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestModulesItemsGetCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get module item - API error",
		Args: []string{"--course-id", "1", "--module-id", "10", "20"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                     courseMock,
			"/api/v1/courses/1/modules/10/items/20": cmdtest.NewErrorResponse(404, "item not found"),
		},
		ExpectError: true,
	}
	cmd := newModulesItemsGetCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestModulesItemsCreateCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create module item - API error",
		Args: []string{"--course-id", "1", "--module-id", "10", "--type", "Assignment", "--title", "Bad Item"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                  courseMock,
			"/api/v1/courses/1/modules/10/items": cmdtest.NewErrorResponse(422, "invalid item"),
		},
		ExpectError: true,
	}
	cmd := newModulesItemsCreateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestModulesItemsUpdateCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update module item - API error",
		Args: []string{"--course-id", "1", "--module-id", "10", "20", "--title", "Updated"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                     courseMock,
			"/api/v1/courses/1/modules/10/items/20": cmdtest.NewErrorResponse(404, "item not found"),
		},
		ExpectError: true,
	}
	cmd := newModulesItemsUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestModulesItemsDeleteCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete module item - API error",
		Args: []string{"--course-id", "1", "--module-id", "10", "20", "--force"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                     courseMock,
			"/api/v1/courses/1/modules/10/items/20": cmdtest.NewErrorResponse(404, "item not found"),
		},
		ExpectError: true,
	}
	cmd := newModulesItemsDeleteCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestModulesItemsDoneCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "mark item done - API error",
		Args: []string{"--course-id", "1", "--module-id", "10", "20"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                          courseMock,
			"/api/v1/courses/1/modules/10/items/20/done": cmdtest.NewErrorResponse(404, "item not found"),
		},
		ExpectError: true,
	}
	cmd := newModulesItemsDoneCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestModulesListCmd_CourseIDError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list modules - course not found",
		Args: []string{"--course-id", "999"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/999": cmdtest.NewErrorResponse(404, "course not found"),
		},
		ExpectError: true,
	}
	cmd := newModulesListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestModulesCreateCmd_WithSuccessOutput(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create module with success output",
		Args: []string{"--course-id", "1", "--name", "Week 1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/modules": cmdtest.NewMockResponse(`{
				"id": 11,
				"name": "Week 1",
				"position": 1,
				"workflow_state": "active"
			}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Week 1") {
				t.Error("Expected 'Week 1' in output")
			}
		},
	}
	cmd := newModulesCreateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}
