package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestBlueprintGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get blueprint template successfully",
			Args: []string{"--course-id", "1", "--template-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts": cmdtest.NewMockResponse(`[]`),
				"/api/v1/courses/1/blueprint_templates/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"course_id": 1,
					"last_export_completed_at": "2024-01-15T10:00:00Z"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "10") {
					t.Errorf("Expected '10' (template ID) in output, got: %s", output)
				}
			},
		},
		{
			Name:        "get blueprint template - missing course ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newBlueprintGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestBlueprintAssociationsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list blueprint associations successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/blueprint_templates/default/associated_courses": cmdtest.NewMockResponse(`[
					{
						"id": 100,
						"name": "Child Course",
						"course_code": "CHILD101"
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Child Course") {
					t.Error("Expected 'Child Course' in output")
				}
			},
		},
		{
			Name: "list blueprint associations - empty response",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/blueprint_templates/default/associated_courses": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No associated courses found",
		},
		{
			Name:        "list blueprint associations - missing course ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newBlueprintAssociationsListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestBlueprintAssociationsAddCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "add blueprint association successfully",
			Args: []string{"--course-id", "1", "--course-ids-to-add", "100,101"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/blueprint_templates/default/update_associations": cmdtest.NewMockResponse(`{}`),
			},
			ExpectError: false,
		},
		{
			Name:        "add blueprint association - missing course ID",
			Args:        []string{"--course-ids-to-add", "100"},
			ExpectError: true,
		},
		{
			Name:        "add blueprint association - missing course IDs to add",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newBlueprintAssociationsAddCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestBlueprintAssociationsRemoveCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "remove blueprint association successfully",
			Args: []string{"--course-id", "1", "--course-ids-to-remove", "100"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/blueprint_templates/default/update_associations": cmdtest.NewMockResponse(`{}`),
			},
			ExpectError: false,
		},
		{
			Name:        "remove blueprint association - missing course IDs",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newBlueprintAssociationsRemoveCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestBlueprintSyncCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "sync blueprint successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/blueprint_templates/default/migrations": cmdtest.NewMockResponse(`{
					"id": 5,
					"workflow_state": "queued"
				}`),
			},
			ExpectError:  false,
			ExpectOutput: "5",
		},
		{
			Name:        "sync blueprint - missing course ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newBlueprintSyncCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestBlueprintChangesCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get blueprint changes successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/blueprint_templates/default/unsynced_changes": cmdtest.NewMockResponse(`[
					{
						"asset_type": "assignment",
						"asset_id": 50,
						"change_type": "updated"
					}
				]`),
			},
			ExpectError: false,
		},
		{
			Name: "blueprint changes - no changes",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/blueprint_templates/default/unsynced_changes": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No unsynced changes",
		},
		{
			Name:        "blueprint changes - missing course ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newBlueprintChangesCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestBlueprintMigrationsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list blueprint migrations successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/blueprint_templates/default/migrations": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"workflow_state": "completed",
						"created_at": "2024-01-15T10:00:00Z"
					}
				]`),
			},
			ExpectError: false,
		},
		{
			Name: "list blueprint migrations - empty",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/blueprint_templates/default/migrations": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No migrations found",
		},
		{
			Name:        "list blueprint migrations - missing course ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newBlueprintMigrationsListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestBlueprintMigrationsGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get blueprint migration successfully",
			Args: []string{"5", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/blueprint_templates/default/migrations/5": cmdtest.NewMockResponse(`{
					"id": 5,
					"workflow_state": "completed",
					"created_at": "2024-01-15T10:00:00Z"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "get blueprint migration - missing migration ID",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "get blueprint migration - missing course ID",
			Args:        []string{"5"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newBlueprintMigrationsGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
