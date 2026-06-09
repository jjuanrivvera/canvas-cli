package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestPlannerItemsCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list planner items successfully",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/planner/items": cmdtest.NewMockResponse(`[
					{
						"plannable_id": 1,
						"plannable_type": "assignment",
						"plannable_date": "2024-02-01T10:00:00Z",
						"plannable": {
							"id": 100,
							"title": "Assignment 1"
						}
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "1") {
					t.Error("Expected assignment data in output")
				}
			},
		},
		{
			Name: "list planner items - empty response",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/planner/items": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No planner items found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPlannerItemsCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPlannerNotesListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list planner notes successfully",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/planner_notes": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"title": "Study Session",
						"todo_date": "2024-02-01T10:00:00Z"
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Study Session") {
					t.Error("Expected 'Study Session' in output")
				}
			},
		},
		{
			Name: "list planner notes - empty response",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/planner_notes": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No planner notes found",
		},
		{
			Name: "list planner notes - API error",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/planner_notes": cmdtest.NewErrorResponse(401, "unauthorized"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPlannerNotesListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPlannerNotesGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get planner note successfully",
			Args: []string{"5"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/planner_notes/5": cmdtest.NewMockResponse(`{
					"id": 5,
					"title": "Review Notes",
					"todo_date": "2024-02-10T09:00:00Z"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Review Notes") {
					t.Error("Expected 'Review Notes' in output")
				}
			},
		},
		{
			Name:        "get planner note - missing note ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPlannerNotesGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPlannerNotesCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create planner note successfully",
			Args: []string{"--title", "Project Work"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/planner_notes": cmdtest.NewMockResponse(`{
					"id": 10,
					"title": "Project Work",
					"todo_date": null
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Project Work") {
					t.Error("Expected 'Project Work' in output")
				}
			},
		},
		{
			Name:        "create planner note - missing title",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPlannerNotesCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPlannerNotesUpdateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "update planner note successfully",
			Args: []string{"5", "--title", "Updated Title"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/planner_notes/5": cmdtest.NewMockResponse(`{
					"id": 5,
					"title": "Updated Title"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Updated Title") {
					t.Error("Expected 'Updated Title' in output")
				}
			},
		},
		{
			Name:        "update planner note - missing note ID",
			Args:        []string{"--title", "Updated"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPlannerNotesUpdateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPlannerNotesDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete planner note successfully",
			Args: []string{"5", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/planner_notes/5": cmdtest.NewMockResponse(`{}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete planner note - missing note ID",
			Args:        []string{"--force"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPlannerNotesDeleteCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPlannerCompleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "mark item complete successfully",
			Args: []string{"Assignment", "123"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/planner/overrides": cmdtest.NewMockResponse(`{
					"id": 1,
					"plannable_type": "Assignment",
					"plannable_id": 123,
					"marked_complete": true
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "complete - missing type and ID",
			Args:        []string{},
			ExpectError: true,
		},
		{
			Name:        "complete - missing ID",
			Args:        []string{"Assignment"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPlannerCompleteCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPlannerDismissCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "dismiss item successfully",
			Args: []string{"CalendarEvent", "789"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/planner/overrides": cmdtest.NewMockResponse(`{
					"id": 2,
					"plannable_type": "CalendarEvent",
					"plannable_id": 789,
					"dismissed": true
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "dismiss - missing arguments",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPlannerDismissCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPlannerOverridesCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list planner overrides successfully",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/planner/overrides": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"plannable_type": "Assignment",
						"plannable_id": 100,
						"marked_complete": true
					}
				]`),
			},
			ExpectError: false,
		},
		{
			Name: "list planner overrides - empty response",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/planner/overrides": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No planner overrides found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPlannerOverridesCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPlannerItemsCmd_WithCourseID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list planner items for course",
		Args: []string{"--course-id", "123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/planner/items": cmdtest.NewMockResponse(`[
				{
					"plannable_id": 5,
					"plannable_type": "quiz",
					"plannable_date": "2024-03-01T09:00:00Z"
				}
			]`),
		},
		ExpectError: false,
	}

	cmd := newPlannerItemsCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}
