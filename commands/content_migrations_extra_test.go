package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestContentMigrationsCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create course copy migration",
			Args: []string{"--course-id", "1", "--type", "course_copy_importer", "--source-course-id", "100"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/content_migrations": cmdtest.NewMockResponse(`{
					"id": 5,
					"migration_type": "course_copy_importer",
					"workflow_state": "created",
					"progress_url": "/api/v1/progress/1"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "5") && !strings.Contains(output, "course_copy_importer") {
					t.Error("expected migration ID or type in output")
				}
			},
		},
		{
			Name:        "create migration - missing course ID",
			Args:        []string{"--type", "course_copy_importer"},
			ExpectError: true,
		},
		{
			Name:        "create migration - missing type",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
		{
			Name: "create migration - invalid copy options JSON",
			Args: []string{"--course-id", "1", "--type", "course_copy_importer", "--copy-options", "notjson"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/content_migrations": cmdtest.NewMockResponse(`{}`),
			},
			ExpectError: true,
		},
		{
			Name: "create migration - invalid date shift JSON",
			Args: []string{"--course-id", "1", "--type", "course_copy_importer", "--date-shift", "notjson"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/content_migrations": cmdtest.NewMockResponse(`{}`),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newContentMigrationsCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestContentMigrationsMigratorsCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list migrators successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/content_migrations/migrators": cmdtest.NewMockResponse(`[
					{"type":"course_copy_importer","name":"Course Copy"},
					{"type":"zip_file_importer","name":"ZIP File"}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "course_copy_importer") {
					t.Error("expected migrator types in output")
				}
			},
		},
		{
			Name: "list migrators - empty",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/content_migrations/migrators": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No migrators available",
		},
		{
			Name:        "list migrators - missing course ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newContentMigrationsMigratorsCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestContentMigrationsContentCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list migration content",
			Args: []string{"5", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/content_migrations/5/content_list": cmdtest.NewMockResponse(`[
					{"type":"assignment","id":"10","label":"HW1"}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "assignment") && !strings.Contains(output, "HW1") {
					t.Error("expected content item in output")
				}
			},
		},
		{
			Name: "list migration content - empty",
			Args: []string{"5", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/content_migrations/5/content_list": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No content available for import",
		},
		{
			Name:        "list migration content - invalid ID",
			Args:        []string{"bad", "--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newContentMigrationsContentCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestContentMigrationsIssuesCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list migration issues",
			Args: []string{"5", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/content_migrations/5/migration_issues": cmdtest.NewMockResponse(`[
					{"id":1,"workflow_state":"active","description":"Could not find attachment"}
				]`),
			},
			ExpectError: false,
		},
		{
			Name: "list migration issues - empty",
			Args: []string{"5", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/content_migrations/5/migration_issues": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No issues found for this migration",
		},
		{
			Name:        "list migration issues - invalid ID",
			Args:        []string{"bad", "--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newContentMigrationsIssuesCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
