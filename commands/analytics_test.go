package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestAnalyticsUserCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get user analytics successfully",
			Args: []string{"100", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": cmdtest.NewMockResponse(`{
					"id": 1,
					"name": "Test Course"
				}`),
				"/api/v1/courses/1/analytics/users/100/activity": cmdtest.NewMockResponse(`{
					"2024-01-15": {
						"views": 150,
						"participations": 25
					}
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "150") {
					t.Error("Expected '150' in output")
				}
			},
		},
		{
			Name:        "get user analytics - missing user ID",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAnalyticsUserCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
