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
