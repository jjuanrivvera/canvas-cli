package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestFilesListCmd_FolderContext(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list files for folder successfully",
		Args: []string{"--folder-id", "5"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/folders/5/files": cmdtest.NewMockResponse(`[
				{
					"id": 3,
					"display_name": "Folder_File.docx",
					"filename": "folder_file.docx",
					"size": 2048
				}
			]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Folder_File.docx") {
				t.Error("Expected 'Folder_File.docx' in output")
			}
		},
	}
	cmd := newFilesListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestFilesListCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list files - API error",
		Args: []string{"--course-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1/files": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmd := newFilesListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestFilesGetCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get file - API error",
		Args: []string{"99"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/files/99": cmdtest.NewErrorResponse(404, "file not found"),
		},
		ExpectError: true,
	}
	cmd := newFilesGetCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestFilesDeleteCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete file - API error",
		Args: []string{"10", "--force"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/files/10": cmdtest.NewErrorResponse(404, "file not found"),
		},
		ExpectError: true,
	}
	cmd := newFilesDeleteCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestFilesQuotaCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get quota - API error",
		Args: []string{"--course-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1/files/quota": cmdtest.NewErrorResponse(403, "forbidden"),
		},
		ExpectError: true,
	}
	cmd := newFilesQuotaCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestFilesQuotaCmd_NoZeroQuota(t *testing.T) {
	// When quota is 0, the usage percentage line is skipped
	tc := cmdtest.CommandTestCase{
		Name: "get quota - zero total quota",
		Args: []string{"--user-id", "100"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/100/files/quota": cmdtest.NewMockResponse(`{
				"quota": 0,
				"quota_used": 0
			}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Storage Quota") {
				t.Error("Expected 'Storage Quota' in output")
			}
		},
	}
	cmd := newFilesQuotaCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestFormatFileSizeBytes(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
	}

	for _, tt := range tests {
		got := formatFileSize(tt.bytes)
		if got != tt.expected {
			t.Errorf("formatFileSize(%d) = %q, want %q", tt.bytes, got, tt.expected)
		}
	}
}
