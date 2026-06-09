package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestOutcomesGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get outcome successfully",
			Args: []string{"10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/outcomes/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"title": "Critical Thinking",
					"description": "Students will demonstrate critical thinking skills"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Critical Thinking") {
					t.Error("Expected 'Critical Thinking' in output")
				}
			},
		},
		{
			Name:        "get outcome - missing outcome ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newOutcomesGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestOutcomesListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list outcomes successfully",
			Args: []string{"--account-id", "1", "--group-id", "456"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts": cmdtest.NewMockResponse(`[]`),
				"/api/v1/accounts/1/outcome_groups/456/outcomes": cmdtest.NewMockResponse(`[
					{
						"outcome": {
							"id": 1,
							"title": "Communication",
							"context_id": 1,
							"context_type": "Account"
						}
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Communication") {
					t.Error("Expected 'Communication' in output")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newOutcomesListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestOutcomesGroupsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list outcome groups successfully",
			Args: []string{"--account-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/outcome_groups": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"title": "Core Competencies",
						"description": "Essential skills"
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Core Competencies") {
					t.Error("Expected 'Core Competencies' in output")
				}
			},
		},
		{
			Name: "list outcome groups - empty response",
			Args: []string{"--account-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/outcome_groups": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No outcome groups found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newOutcomesGroupsListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestOutcomesGroupsGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get outcome group successfully",
			Args: []string{"5", "--account-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/outcome_groups/5": cmdtest.NewMockResponse(`{
					"id": 5,
					"title": "Program Outcomes",
					"description": "Program-level learning outcomes"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Program Outcomes") {
					t.Error("Expected 'Program Outcomes' in output")
				}
			},
		},
		{
			Name:        "get outcome group - missing group ID",
			Args:        []string{"--account-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newOutcomesGroupsGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestOutcomesCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create outcome in account group successfully",
			Args: []string{"--account-id", "1", "--group-id", "456", "--title", "Problem Solving"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/outcome_groups/456/outcomes": cmdtest.NewMockResponse(`{
					"outcome": {
						"id": 77,
						"title": "Problem Solving"
					},
					"context_type": "Account",
					"context_id": 1
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Problem Solving") {
					t.Error("Expected 'Problem Solving' in output")
				}
			},
		},
		{
			Name:        "create outcome - missing group ID",
			Args:        []string{"--account-id", "1", "--title", "Test"},
			ExpectError: true,
		},
		{
			Name:        "create outcome - missing title",
			Args:        []string{"--account-id", "1", "--group-id", "456"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newOutcomesCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestOutcomesUpdateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "update outcome successfully",
			Args: []string{"10", "--title", "Updated Outcome"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/outcomes/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"title": "Updated Outcome",
					"description": "Updated description"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Updated Outcome") {
					t.Error("Expected 'Updated Outcome' in output")
				}
			},
		},
		{
			Name:        "update outcome - missing outcome ID",
			Args:        []string{"--title", "Updated"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newOutcomesUpdateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestOutcomesLinkCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "link outcome to group successfully",
			Args: []string{"789", "--account-id", "1", "--group-id", "456"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/outcome_groups/456/outcomes/789": cmdtest.NewMockResponse(`{
					"outcome": {
						"id": 789,
						"title": "Linked Outcome"
					},
					"context_type": "Account",
					"context_id": 1
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "link outcome - missing group ID",
			Args:        []string{"789", "--account-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "link outcome - missing outcome ID",
			Args:        []string{"--account-id", "1", "--group-id", "456"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newOutcomesLinkCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestOutcomesUnlinkCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "unlink outcome from group successfully",
			Args: []string{"789", "--account-id", "1", "--group-id", "456", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/1/outcome_groups/456/outcomes/789": cmdtest.NewMockResponse(`{
					"outcome": {
						"id": 789,
						"title": "Unlinked Outcome"
					}
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "unlink outcome - missing group ID",
			Args:        []string{"789", "--account-id", "1", "--force"},
			ExpectError: true,
		},
		{
			Name:        "unlink outcome - missing outcome ID",
			Args:        []string{"--account-id", "1", "--group-id", "456", "--force"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newOutcomesUnlinkCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestOutcomesResultsCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get outcome results successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/outcome_results": cmdtest.NewMockResponse(`{
					"outcome_results": [
						{
							"id": 1,
							"score": 4.0,
							"possible": 5.0
						}
					]
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "get outcome results - missing course ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newOutcomesResultsCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestOutcomesAlignmentsCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get outcome alignments successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/outcome_alignments": cmdtest.NewMockResponse(`[
					{
						"id": "align_1",
						"name": "Assignment Alignment"
					}
				]`),
			},
			ExpectError: false,
		},
		{
			Name: "get outcome alignments - empty response",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/outcome_alignments": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No outcome alignments found",
		},
		{
			Name:        "get outcome alignments - missing course ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newOutcomesAlignmentsCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestOutcomesGroupsListCmd_CourseContext(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list outcome groups in course",
		Args: []string{"--course-id", "123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/123/outcome_groups": cmdtest.NewMockResponse(`[
				{
					"id": 10,
					"title": "Course Outcomes",
					"description": "Outcomes for this course"
				}
			]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Course Outcomes") {
				t.Error("Expected 'Course Outcomes' in output")
			}
		},
	}

	cmd := newOutcomesGroupsListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestOutcomesGetCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get outcome - API error",
		Args: []string{"999"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/outcomes/999": cmdtest.NewErrorResponse(404, "outcome not found"),
		},
		ExpectError: true,
	}

	cmd := newOutcomesGetCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}
