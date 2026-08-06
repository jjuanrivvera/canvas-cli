package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

// accountsMockResponse is needed because the Canvas client checks the API version on first use.
var accountsMockResponse = cmdtest.NewMockResponse(`[{"id":1,"name":"Test Account"}]`)

// prependAPI adds "api" as the first arg so RunCommandTest dispatches via rootCmd.
func prependAPI(args []string) []string {
	return append([]string{"api"}, args...)
}

// runAPICmdTest replaces the api subcommand with a fresh instance before each
// test case, then runs the case through the shared rootCmd.
//
// Cobra does not reset flag values (string/StringArray) between Execute() calls
// when the same *cobra.Command is reused.  Replacing the registered api command
// with a freshly-constructed one guarantees each test sees only its own flags,
// which is what makes these tests deterministic in CI.
func runAPICmdTest(t *testing.T, tc cmdtest.CommandTestCase) {
	t.Helper()

	// Remove the stale api subcommand (if any) and register a fresh one.
	apiSub, _, _ := rootCmd.Find([]string{"api"})
	if apiSub != nil && apiSub != rootCmd {
		rootCmd.RemoveCommand(apiSub)
	}
	rootCmd.AddCommand(newAPICmd())

	cmdtest.RunCommandTest(t, rootCmd, tc)
}

func TestAPICmd_GetRequest(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "GET request returns data",
			Args: prependAPI([]string{"GET", "/api/v1/courses"}),
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts": accountsMockResponse,
				"/api/v1/courses":  cmdtest.NewMockResponse(`[{"id":1,"name":"Test Course"}]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Test Course") && !strings.Contains(output, "status_code") {
					t.Errorf("expected 'Test Course' or 'status_code' in output, got: %s", output)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			runAPICmdTest(t, tc)
		})
	}
}

func TestAPICmd_PostRequest(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "POST request with JSON body",
			Args: prependAPI([]string{"POST", "/api/v1/courses", "--data", `{"course":{"name":"New Course"}}`}),
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts": accountsMockResponse,
				"/api/v1/courses":  cmdtest.NewMockResponse(`{"id":99,"name":"New Course"}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "status_code") && !strings.Contains(output, "99") {
					t.Errorf("expected status or ID in output, got: %s", output)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			runAPICmdTest(t, tc)
		})
	}
}

func TestAPICmd_InvalidMethod(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name:          "unsupported HTTP method",
			Args:          prependAPI([]string{"TRACE", "/api/v1/courses"}),
			MockResponses: map[string]cmdtest.MockResponse{},
			ExpectError:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			runAPICmdTest(t, tc)
		})
	}
}

func TestAPICmd_InvalidJSONBody(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "invalid JSON in --data",
			Args: prependAPI([]string{"POST", "/api/v1/courses", "--data", `not-json`}),
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts": accountsMockResponse,
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			runAPICmdTest(t, tc)
		})
	}
}

func TestAPICmd_BothDataFlags(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "cannot use both --data and --data-file",
			Args: prependAPI([]string{"POST", "/api/v1/courses", "--data", `{}`, "--data-file", "some.json"}),
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts": accountsMockResponse,
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			runAPICmdTest(t, tc)
		})
	}
}

func TestAPICmd_QueryParams(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "GET with query params",
			Args: prependAPI([]string{"GET", "/api/v1/courses", "-q", "per_page=10"}),
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts": accountsMockResponse,
				"/api/v1/courses":  cmdtest.NewMockResponse(`[{"id":1}]`),
			},
			ExpectError: false,
		},
		{
			Name: "invalid query param format",
			Args: prependAPI([]string{"GET", "/api/v1/courses", "-q", "badformat"}),
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts": accountsMockResponse,
				"/api/v1/courses":  cmdtest.NewMockResponse(`[]`),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			runAPICmdTest(t, tc)
		})
	}
}

func TestAPICmd_CustomHeaders(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "GET with custom header",
			Args: prependAPI([]string{"GET", "/api/v1/courses", "-H", "X-Custom:value"}),
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts": accountsMockResponse,
				"/api/v1/courses":  cmdtest.NewMockResponse(`[]`),
			},
			ExpectError: false,
		},
		{
			Name: "invalid header format",
			Args: prependAPI([]string{"GET", "/api/v1/courses", "-H", "badheader"}),
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts": accountsMockResponse,
				"/api/v1/courses":  cmdtest.NewMockResponse(`[]`),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			runAPICmdTest(t, tc)
		})
	}
}

func TestAPICmd_RawOutput(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "GET with raw output flag",
			Args: prependAPI([]string{"GET", "/api/v1/courses", "--raw"}),
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts": accountsMockResponse,
				"/api/v1/courses":  cmdtest.NewMockResponse(`[{"id":1}]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, `"id"`) {
					t.Errorf("expected JSON body in raw output, got: %s", output)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			runAPICmdTest(t, tc)
		})
	}
}

func TestAPICmd_ShowHeaders(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "GET with show-headers flag",
			Args: prependAPI([]string{"GET", "/api/v1/courses", "--show-headers"}),
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts": accountsMockResponse,
				"/api/v1/courses":  cmdtest.NewMockResponse(`[]`),
			},
			ExpectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			runAPICmdTest(t, tc)
		})
	}
}

func TestAPICmd_DeleteRequest(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "DELETE request",
			Args: prependAPI([]string{"DELETE", "/api/v1/courses/1"}),
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts":  accountsMockResponse,
				"/api/v1/courses/1": cmdtest.NewMockResponse(`{"deleted":true}`),
			},
			ExpectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			runAPICmdTest(t, tc)
		})
	}
}

// TestAPIGetCmd_ReadOnly covers the GET-only "canvas api get" sibling (#60):
// it dispatches, returns data, and — unlike "canvas api" — does not accept a
// request body.
func TestAPIGetCmd_ReadOnly(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "api get returns data",
			Args: prependAPI([]string{"get", "/api/v1/courses"}),
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts": accountsMockResponse,
				"/api/v1/courses":  cmdtest.NewMockResponse(`[{"id":1,"name":"Test Course"}]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Test Course") && !strings.Contains(output, "status_code") {
					t.Errorf("expected course data in output, got: %s", output)
				}
			},
		},
		{
			Name:        "api get rejects a body flag",
			Args:        prependAPI([]string{"get", "/api/v1/courses", "--data", `{"x":1}`}),
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			runAPICmdTest(t, tc)
		})
	}
}
