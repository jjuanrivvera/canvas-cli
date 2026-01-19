package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestCalendarListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list calendar events successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/calendar_events": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"title": "Team Meeting",
						"start_at": "2024-01-15T10:00:00Z",
						"end_at": "2024-01-15T11:00:00Z",
						"context_code": "course_1"
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Team Meeting") {
					t.Error("Expected 'Team Meeting' in output")
				}
			},
		},
		{
			Name: "list calendar events - empty response",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/calendar_events": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No calendar events found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newCalendarListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestCalendarGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get calendar event successfully",
			Args: []string{"10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/calendar_events/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"title": "Office Hours",
					"start_at": "2024-01-16T14:00:00Z",
					"description": "Weekly office hours"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Office Hours") {
					t.Error("Expected 'Office Hours' in output")
				}
			},
		},
		{
			Name:        "get calendar event - missing event ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newCalendarGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestCalendarCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create calendar event successfully",
			Args: []string{"--course-id", "1", "--title", "New Event", "--start-at", "2024-02-01T10:00:00Z"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts": cmdtest.NewMockResponse(`[]`),
				"/api/v1/calendar_events": cmdtest.NewMockResponse(`{
					"id": 20,
					"title": "New Event",
					"context_code": "course_1",
					"start_at": "2024-02-01T10:00:00Z"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "New Event") {
					t.Error("Expected 'New Event' in output")
				}
			},
		},
		{
			Name:        "create calendar event - missing title",
			Args:        []string{"--course-id", "1", "--start-at", "2024-02-01T10:00:00Z"},
			ExpectError: true,
		},
		{
			Name:        "create calendar event - missing course context",
			Args:        []string{"--title", "New Event"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newCalendarCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestCalendarUpdateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "update calendar event successfully",
			Args: []string{"10", "--title", "Updated Event"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/calendar_events/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"title": "Updated Event"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Updated Event") {
					t.Error("Expected 'Updated Event' in output")
				}
			},
		},
		{
			Name:        "update calendar event - missing event ID",
			Args:        []string{"--title", "Updated"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newCalendarUpdateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestCalendarDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete calendar event with confirmation",
			Args: []string{"10", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/calendar_events/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"title": "Old Event"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete calendar event - missing event ID",
			Args:        []string{"--force"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newCalendarDeleteCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
