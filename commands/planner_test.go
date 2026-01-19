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
