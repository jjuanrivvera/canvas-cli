package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestUsersSearchCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "search users - API error",
		// search-term is a positional argument
		Args: []string{"John"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/search/recipients": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmd := newUsersSearchCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestUsersSearchCmd_EmptyResults(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "search users - no results",
		Args: []string{"nonexistent"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/search/recipients": cmdtest.NewMockResponse(`[]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "No users found") {
				t.Error("Expected 'No users found' in output")
			}
		},
	}
	cmd := newUsersSearchCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestUsersMeCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get current user - API error",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/profile": cmdtest.NewErrorResponse(401, "unauthorized"),
		},
		ExpectError: true,
	}
	cmd := newUsersMeCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestUsersCreateCmd_NoName(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create user - no name and no JSON",
		Args: []string{"--account-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/users": cmdtest.NewMockResponse(`{}`),
		},
		ExpectError: true,
	}
	cmd := newUsersCreateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestUsersCreateCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create user - API error",
		Args: []string{"--account-id", "1", "--name", "Bad User"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/users": cmdtest.NewErrorResponse(422, "invalid user"),
		},
		ExpectError: true,
	}
	cmd := newUsersCreateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestUsersCreateCmd_WithAllFlags(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create user with all flags",
		Args: []string{
			"--account-id", "1",
			"--name", "John Doe",
			"--email", "john@example.com",
			"--login-id", "john.doe",
			"--timezone", "America/New_York",
		},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/users": cmdtest.NewMockResponse(`{
				"id": 100,
				"name": "John Doe",
				"login_id": "john.doe",
				"email": "john@example.com"
			}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "User created successfully") {
				t.Error("Expected 'User created successfully' in output")
			}
		},
	}
	cmd := newUsersCreateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestUsersUpdateCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update user - API error",
		Args: []string{"99", "--name", "Updated Name"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/99": cmdtest.NewErrorResponse(404, "user not found"),
		},
		ExpectError: true,
	}
	cmd := newUsersUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestUsersUpdateCmd_WithAllFlags(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update user with all flags",
		Args: []string{
			"100",
			"--name", "Updated Name",
			"--email", "updated@example.com",
			"--timezone", "UTC",
			"--locale", "en",
		},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/100": cmdtest.NewMockResponse(`{
				"id": 100,
				"name": "Updated Name",
				"email": "updated@example.com"
			}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "User updated successfully") {
				t.Error("Expected 'User updated successfully' in output")
			}
		},
	}
	cmd := newUsersUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestUsersListCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list users - API error",
		Args: []string{"--account-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/users": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmd := newUsersListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestUsersCreateCmd_WithEmail(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create user displays email in success output",
		Args: []string{"--account-id", "1", "--name", "Jane Smith", "--email", "jane@example.com"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/users": cmdtest.NewMockResponse(`{
				"id": 101,
				"name": "Jane Smith",
				"email": "jane@example.com",
				"login_id": "jane.smith"
			}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "jane@example.com") {
				t.Error("Expected email in output")
			}
		},
	}
	cmd := newUsersCreateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestUsersCreateCmd_FromJSONFile(t *testing.T) {
	tmpDir := t.TempDir()
	jsonPath := filepath.Join(tmpDir, "user.json")
	jsonContent := `{"name":"JSON User","email":"json@example.com","login_id":"json.user"}`
	if err := os.WriteFile(jsonPath, []byte(jsonContent), 0644); err != nil {
		t.Fatal(err)
	}

	tc := cmdtest.CommandTestCase{
		Name: "create user from JSON file",
		Args: []string{"--account-id", "1", "--json", jsonPath},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/users": cmdtest.NewMockResponse(`{
				"id": 200,
				"name": "JSON User",
				"email": "json@example.com",
				"login_id": "json.user"
			}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "User created successfully") {
				t.Error("Expected 'User created successfully' in output")
			}
		},
	}
	cmd := newUsersCreateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestUsersCreateCmd_InvalidJSONFile(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create user from nonexistent JSON file",
		Args: []string{"--account-id", "1", "--json", "/nonexistent/user.json"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/users": cmdtest.NewMockResponse(`{}`),
		},
		ExpectError: true,
	}
	cmd := newUsersCreateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestUsersUpdateCmd_FromJSONFile(t *testing.T) {
	tmpDir := t.TempDir()
	jsonPath := filepath.Join(tmpDir, "update.json")
	jsonContent := `{"name":"Updated JSON User","email":"updated@example.com"}`
	if err := os.WriteFile(jsonPath, []byte(jsonContent), 0644); err != nil {
		t.Fatal(err)
	}

	tc := cmdtest.CommandTestCase{
		Name: "update user from JSON file",
		Args: []string{"100", "--json", jsonPath},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/100": cmdtest.NewMockResponse(`{
				"id": 100,
				"name": "Updated JSON User",
				"email": "updated@example.com"
			}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "User updated successfully") {
				t.Error("Expected 'User updated successfully' in output")
			}
		},
	}
	cmd := newUsersUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestUsersUpdateCmd_InvalidJSONFile(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update user from nonexistent JSON file",
		Args: []string{"100", "--json", "/nonexistent/update.json"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/100": cmdtest.NewMockResponse(`{}`),
		},
		ExpectError: true,
	}
	cmd := newUsersUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}
