package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestFilesListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list files for course successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/files": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"display_name": "Syllabus.pdf",
						"filename": "syllabus.pdf",
						"size": 102400,
						"content-type": "application/pdf"
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Syllabus.pdf") {
					t.Error("Expected 'Syllabus.pdf' in output")
				}
			},
		},
		{
			Name: "list files - empty response",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/files": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No files found",
		},
		{
			Name: "list files for user successfully",
			Args: []string{"--user-id", "100"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/users/100/files": cmdtest.NewMockResponse(`[
					{
						"id": 2,
						"display_name": "Notes.txt",
						"filename": "notes.txt",
						"size": 1024
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Notes.txt") {
					t.Error("Expected 'Notes.txt' in output")
				}
			},
		},
		{
			Name:        "list files - no context specified",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newFilesListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestFilesGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get file successfully",
			Args: []string{"10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/files/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"display_name": "Lecture1.pdf",
					"filename": "lecture1.pdf",
					"size": 524288,
					"content-type": "application/pdf",
					"url": "https://example.com/files/10"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Lecture1.pdf") {
					t.Error("Expected 'Lecture1.pdf' in output")
				}
			},
		},
		{
			Name:        "get file - missing file ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newFilesGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestFilesDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete file with confirmation",
			Args: []string{"10", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/files/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"display_name": "OldFile.pdf"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "deleted successfully") {
					t.Error("Expected 'deleted successfully' in output")
				}
			},
		},
		{
			Name:        "delete file - missing file ID",
			Args:        []string{"--force"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newFilesDeleteCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestFilesQuotaCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get quota for course successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/files/quota": cmdtest.NewMockResponse(`{
					"quota": 524288000,
					"quota_used": 104857600
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Storage Quota") {
					t.Error("Expected 'Storage Quota' in output")
				}
			},
		},
		{
			Name: "get quota for user successfully",
			Args: []string{"--user-id", "100"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/users/100/files/quota": cmdtest.NewMockResponse(`{
					"quota": 1073741824,
					"quota_used": 268435456
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Storage Quota") {
					t.Error("Expected 'Storage Quota' in output")
				}
			},
		},
		{
			Name:        "get quota - no context specified",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newFilesQuotaCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
