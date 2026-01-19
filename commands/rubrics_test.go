package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestRubricsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list rubrics for course successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts":  cmdtest.NewMockResponse(`[]`),
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/rubrics": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"title": "Essay Rubric",
						"points_possible": 100,
						"context_id": 1,
						"context_type": "Course"
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Essay Rubric") {
					t.Error("Expected 'Essay Rubric' in output")
				}
			},
		},
		{
			Name: "list rubrics - empty response",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":         courseMock,
				"/api/v1/courses/1/rubrics": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No rubrics found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newRubricsListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestRubricsGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get rubric successfully",
			Args: []string{"10", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/rubrics/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"title": "Discussion Rubric",
					"points_possible": 50,
					"context_id": 1,
					"context_type": "Course"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Discussion Rubric") {
					t.Error("Expected 'Discussion Rubric' in output")
				}
			},
		},
		{
			Name:        "get rubric - missing rubric ID",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newRubricsGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestRubricsCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create rubric successfully",
			Args: []string{"--course-id", "1", "--title", "New Rubric", "--points", "100"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts": cmdtest.NewMockResponse(`[]`),
				"/api/v1/courses/1/rubrics": cmdtest.NewMockResponse(`{
					"rubric": {
						"id": 20,
						"title": "New Rubric",
						"points_possible": 100,
						"context_id": 1,
						"context_type": "Course"
					}
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "New Rubric") && !strings.Contains(output, "20") {
					t.Error("Expected rubric title or ID in output")
				}
			},
		},
		{
			Name:        "create rubric - missing course ID",
			Args:        []string{"--title", "New Rubric"},
			ExpectError: true,
		},
		{
			Name:        "create rubric - missing title",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newRubricsCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestRubricsUpdateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "update rubric successfully",
			Args: []string{"10", "--course-id", "1", "--title", "Updated Rubric"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts": cmdtest.NewMockResponse(`[]`),
				"/api/v1/courses/1/rubrics/10": cmdtest.NewMockResponse(`{
					"rubric": {
						"id": 10,
						"title": "Updated Rubric",
						"points_possible": 100,
						"context_id": 1,
						"context_type": "Course"
					}
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Updated Rubric") && !strings.Contains(output, "10") {
					t.Error("Expected rubric title or ID in output")
				}
			},
		},
		{
			Name:        "update rubric - missing rubric ID",
			Args:        []string{"--course-id", "1", "--title", "Updated"},
			ExpectError: true,
		},
		{
			Name:        "update rubric - missing course ID",
			Args:        []string{"10", "--title", "Updated"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newRubricsUpdateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestRubricsDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete rubric with confirmation",
			Args: []string{"10", "--course-id", "1", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts": cmdtest.NewMockResponse(`[]`),
				"/api/v1/courses/1/rubrics/10": cmdtest.NewMockResponse(`{
					"rubric": {
						"id": 10,
						"title": "Old Rubric"
					}
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete rubric - missing rubric ID",
			Args:        []string{"--course-id", "1", "--force"},
			ExpectError: true,
		},
		{
			Name:        "delete rubric - missing course ID",
			Args:        []string{"10", "--force"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newRubricsDeleteCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestRubricsAssociateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "associate rubric with assignment successfully",
			Args: []string{"10", "--course-id", "1", "--assignment-id", "100"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts": cmdtest.NewMockResponse(`[]`),
				"/api/v1/courses/1/rubric_associations": cmdtest.NewMockResponse(`{
					"rubric_association": {
						"id": 50,
						"rubric_id": 10,
						"association_id": 100,
						"association_type": "Assignment",
						"use_for_grading": false
					}
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "50") && !strings.Contains(output, "10") {
					t.Error("Expected association ID or rubric ID in output")
				}
			},
		},
		{
			Name:        "associate rubric - missing rubric ID",
			Args:        []string{"--course-id", "1", "--assignment-id", "100"},
			ExpectError: true,
		},
		{
			Name:        "associate rubric - missing course ID",
			Args:        []string{"10", "--assignment-id", "100"},
			ExpectError: true,
		},
		{
			Name:        "associate rubric - missing assignment ID",
			Args:        []string{"10", "--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newRubricsAssociateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
