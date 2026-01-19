package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestPeerReviewsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list peer reviews successfully",
			Args: []string{"--course-id", "1", "--assignment-id", "100"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts":  cmdtest.NewMockResponse(`[]`),
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/assignments/100/peer_reviews": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"asset_id": 100,
						"assessor_id": 200,
						"user_id": 300,
						"workflow_state": "assigned"
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "assigned") {
					t.Error("Expected 'assigned' in output")
				}
			},
		},
		{
			Name: "list peer reviews - empty response",
			Args: []string{"--course-id", "1", "--assignment-id", "100"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/assignments/100/peer_reviews": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No peer reviews found",
		},
		{
			Name:        "list peer reviews - missing assignment ID",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPeerReviewsListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPeerReviewsCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create peer review successfully",
			Args: []string{"--course-id", "1", "--assignment-id", "100", "--submission-id", "200", "--user-id", "300"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts": cmdtest.NewMockResponse(`[]`),
				"/api/v1/courses/1/assignments/100/submissions/200/peer_reviews": cmdtest.NewMockResponse(`{
					"id": 10,
					"asset_id": 100,
					"assessor_id": 300,
					"user_id": 200,
					"workflow_state": "assigned"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "100") && !strings.Contains(output, "300") {
					t.Error("Expected peer review data in output")
				}
			},
		},
		{
			Name:        "create peer review - missing assessor ID",
			Args:        []string{"--course-id", "1", "--assignment-id", "100", "--user-id", "200"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPeerReviewsCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestPeerReviewsDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete peer review successfully",
			Args: []string{"--course-id", "1", "--assignment-id", "100", "--submission-id", "200", "--user-id", "300", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts": cmdtest.NewMockResponse(`[]`),
				"/api/v1/courses/1/assignments/100/submissions/200/peer_reviews": cmdtest.NewMockResponse(`{
					"id": 10
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete peer review - missing user ID",
			Args:        []string{"--course-id", "1", "--assignment-id", "100", "--user-id", "300", "--force"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newPeerReviewsDeleteCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
