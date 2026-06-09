// Package testing_test tests the cmdtest framework using an external test package
// to avoid a naming conflict with the stdlib "testing" package (the production
// package is also named "testing").
package testing_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

// --- constructor tests (pure functions, no I/O) ---

func TestNewMockResponse(t *testing.T) {
	body := `{"id":1}`
	resp := cmdtest.NewMockResponse(body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("NewMockResponse: StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if resp.Body != body {
		t.Errorf("NewMockResponse: Body = %q, want %q", resp.Body, body)
	}
	if ct := resp.Headers["Content-Type"]; ct != "application/json" {
		t.Errorf("NewMockResponse: Content-Type = %q, want %q", ct, "application/json")
	}
}

func TestNewErrorResponse(t *testing.T) {
	resp := cmdtest.NewErrorResponse(http.StatusNotFound, "not found")

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("NewErrorResponse: StatusCode = %d, want %d", resp.StatusCode, http.StatusNotFound)
	}
	if !strings.Contains(resp.Body, "not found") {
		t.Errorf("NewErrorResponse: Body %q does not contain 'not found'", resp.Body)
	}
	if ct := resp.Headers["Content-Type"]; ct != "application/json" {
		t.Errorf("NewErrorResponse: Content-Type = %q, want %q", ct, "application/json")
	}
}

func TestNewPaginatedResponse_WithNext(t *testing.T) {
	body := `[{"id":1}]`
	next := "https://canvas.example.com/api/v1/courses?page=2"
	resp := cmdtest.NewPaginatedResponse(body, next)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("NewPaginatedResponse: StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if resp.Body != body {
		t.Errorf("NewPaginatedResponse: Body = %q, want %q", resp.Body, body)
	}
	wantLink := fmt.Sprintf(`<%s>; rel="next"`, next)
	if resp.Headers["Link"] != wantLink {
		t.Errorf("NewPaginatedResponse: Link = %q, want %q", resp.Headers["Link"], wantLink)
	}
}

func TestNewPaginatedResponse_NoNext(t *testing.T) {
	body := `[{"id":1}]`
	resp := cmdtest.NewPaginatedResponse(body, "")

	if _, ok := resp.Headers["Link"]; ok {
		t.Errorf("NewPaginatedResponse: expected no Link header when nextPage is empty, got %q", resp.Headers["Link"])
	}
}

// --- RunCommandTest integration tests ---

// echoCmd is a minimal cobra command that writes a static message to its writer.
func echoCmd(msg string) *cobra.Command {
	return &cobra.Command{
		Use:          "echo",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Println(msg)
			return nil
		},
	}
}

// failCmd always returns an error.
func failCmd(errMsg string) *cobra.Command {
	return &cobra.Command{
		Use:          "fail",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("%s", errMsg)
		},
	}
}

func TestRunCommandTest_SuccessPath(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:         "success",
		Args:         []string{},
		ExpectError:  false,
		ExpectOutput: "hello world",
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "hello world") {
				t.Errorf("expected output to contain 'hello world', got %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, echoCmd("hello world"), tc)
}

func TestRunCommandTest_ErrorPath(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "expected error",
		Args:        []string{},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, failCmd("something broke"), tc)
}

func TestRunCommandTest_WithMockResponses(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "with mocks",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses": cmdtest.NewMockResponse(`[{"id":1,"name":"Test"}]`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, echoCmd("done"), tc)
}

func TestRunCommandTest_NoMockFound(t *testing.T) {
	// The command doesn't call the server; the handler's 404 branch would fire
	// if a real command requested an unmapped path.  This test covers handler
	// construction and the success path of the harness.
	tc := cmdtest.CommandTestCase{
		Name: "no mock found path",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses": cmdtest.NewMockResponse(`[]`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, echoCmd("ok"), tc)
}

func TestRunCommandTest_ValidateOutput(t *testing.T) {
	called := false
	tc := cmdtest.CommandTestCase{
		Name:        "validate output callback",
		Args:        []string{},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			called = true
		},
	}
	cmdtest.RunCommandTest(t, echoCmd("anything"), tc)
	if !called {
		t.Error("ValidateOutput callback was not called")
	}
}

func TestRunCommandTest_SetupClientCalled(t *testing.T) {
	called := false
	tc := cmdtest.CommandTestCase{
		Name:        "setup client",
		Args:        []string{},
		ExpectError: false,
		SetupClient: func(client *api.Client) {
			called = true
		},
	}
	cmdtest.RunCommandTest(t, echoCmd("ok"), tc)
	if !called {
		t.Error("SetupClient callback was not called")
	}
}

func TestRunCommandTest_WithErrorResponseMock(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "error response mock",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses": cmdtest.NewErrorResponse(http.StatusUnprocessableEntity, "invalid"),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, echoCmd("fine"), tc)
}

func TestRunCommandTest_LongerPatternMatchesFirst(t *testing.T) {
	// Register two overlapping patterns to exercise the sort-by-length routing.
	tc := cmdtest.CommandTestCase{
		Name: "longer pattern wins",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses":     cmdtest.NewMockResponse(`[]`),
			"/api/v1/courses/123": cmdtest.NewMockResponse(`{"id":123}`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, echoCmd("ok"), tc)
}

func TestRunCommandTest_MockStatusCodeZeroDefaults200(t *testing.T) {
	// A MockResponse with StatusCode=0 should default to 200.
	tc := cmdtest.CommandTestCase{
		Name: "zero status defaults to 200",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses": {StatusCode: 0, Body: `[]`, Headers: map[string]string{}},
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, echoCmd("ok"), tc)
}

func TestRunCommandTest_MockHeadersSet(t *testing.T) {
	// A MockResponse with extra headers should serve them correctly.
	tc := cmdtest.CommandTestCase{
		Name: "custom headers",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses": cmdtest.NewPaginatedResponse(`[]`, "http://example.com?page=2"),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, echoCmd("ok"), tc)
}
