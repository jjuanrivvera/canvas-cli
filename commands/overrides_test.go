package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestOverridesListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list overrides successfully",
			Args: []string{"--course-id", "1", "--assignment-id", "100"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts":  cmdtest.NewMockResponse(`[]`),
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/assignments/100/overrides": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"assignment_id": 100,
						"course_section_id": 10,
						"due_at": "2024-02-15T23:59:00Z"
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "1") && !strings.Contains(output, "10") {
					t.Error("Expected override data in output")
				}
			},
		},
		{
			Name: "list overrides - empty response",
			Args: []string{"--course-id", "1", "--assignment-id", "100"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                           courseMock,
				"/api/v1/courses/1/assignments/100/overrides": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No overrides found",
		},
		{
			Name:        "list overrides - missing assignment ID",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newOverridesListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestOverridesGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get override successfully",
			Args: []string{"5", "--course-id", "1", "--assignment-id", "100"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/assignments/100/overrides/5": cmdtest.NewMockResponse(`{
					"id": 5,
					"assignment_id": 100,
					"course_section_id": 10,
					"due_at": "2024-02-20T23:59:00Z"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "10") {
					t.Error("Expected '2024-02-20' in output")
				}
			},
		},
		{
			Name:        "get override - missing override ID",
			Args:        []string{"--course-id", "1", "--assignment-id", "100"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newOverridesGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestOverridesCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create override successfully",
			Args: []string{"--course-id", "1", "--assignment-id", "100", "--section-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/assignments/100/overrides": cmdtest.NewMockResponse(`{
					"id": 20,
					"assignment_id": 100,
					"course_section_id": 10
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Override created") {
					t.Error("Expected 'Override created' in output")
				}
			},
		},
		{
			Name:        "create override - missing assignment ID",
			Args:        []string{"--course-id", "1", "--section-id", "10"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newOverridesCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestOverridesUpdateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "update override successfully",
			Args: []string{"5", "--course-id", "1", "--assignment-id", "100", "--due-at", "2024-03-01T23:59:00Z"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts": cmdtest.NewMockResponse(`[]`),
				"/api/v1/courses/1/assignments/100/overrides/5": cmdtest.NewMockResponse(`{
					"id": 5,
					"assignment_id": 100,
					"due_at": "2024-03-01T23:59:00Z"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "5") && !strings.Contains(output, "100") {
					t.Error("Expected override ID or assignment ID in output")
				}
			},
		},
		{
			Name:        "update override - missing override ID",
			Args:        []string{"--course-id", "1", "--assignment-id", "100"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newOverridesUpdateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestOverridesDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete override with confirmation",
			Args: []string{"5", "--course-id", "1", "--assignment-id", "100", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/assignments/100/overrides/5": cmdtest.NewMockResponse(`{
					"id": 5
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete override - missing override ID",
			Args:        []string{"--course-id", "1", "--assignment-id", "100", "--force"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newOverridesDeleteCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
