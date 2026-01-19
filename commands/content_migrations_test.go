package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestContentMigrationsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list content migrations successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts": cmdtest.NewMockResponse(`[]`),
				"/api/v1/courses/1/content_migrations": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"migration_type": "course_copy_importer",
						"workflow_state": "completed",
						"progress_url": "/api/v1/progress/100"
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "1") && !strings.Contains(output, "course_copy_importer") {
					t.Error("Expected migration ID or type in output")
				}
			},
		},
		{
			Name: "list content migrations - empty response",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts":                     cmdtest.NewMockResponse(`[]`),
				"/api/v1/courses/1/content_migrations": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No content migrations found",
		},
		{
			Name:        "list content migrations - missing course ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newContentMigrationsListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestContentMigrationsGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get content migration successfully",
			Args: []string{"10", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts": cmdtest.NewMockResponse(`[]`),
				"/api/v1/courses/1/content_migrations/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"migration_type": "common_cartridge_importer",
					"workflow_state": "running",
					"progress_url": "/api/v1/progress/200"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "10") && !strings.Contains(output, "common_cartridge_importer") {
					t.Error("Expected migration ID or type in output")
				}
			},
		},
		{
			Name:        "get content migration - missing migration ID",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newContentMigrationsGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
