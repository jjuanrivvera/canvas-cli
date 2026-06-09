package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestCalendarReserveCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "reserve calendar appointment successfully",
			Args: []string{"10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/calendar_events/10/reservations": cmdtest.NewMockResponse(`{
					"id": 100,
					"title": "Reserved Slot",
					"start_at": "2024-02-15T10:00:00Z",
					"context_code": "course_1"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Reserved Slot") {
					t.Error("Expected 'Reserved Slot' in output")
				}
			},
		},
		{
			Name:        "reserve calendar appointment - missing event ID",
			Args:        []string{},
			ExpectError: true,
		},
		{
			Name: "reserve calendar appointment - API error",
			Args: []string{"10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/calendar_events/10/reservations": cmdtest.NewErrorResponse(404, "event not found"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newCalendarReserveCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestCalendarGetCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get calendar event - API error",
		Args: []string{"10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/calendar_events/10": cmdtest.NewErrorResponse(404, "event not found"),
		},
		ExpectError: true,
	}
	cmd := newCalendarGetCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestCalendarCreateCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create calendar event - API error",
		Args: []string{"--course-id", "1", "--title", "Bad Event", "--start-at", "2024-02-01T10:00:00Z"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts":        cmdtest.NewMockResponse(`[]`),
			"/api/v1/calendar_events": cmdtest.NewErrorResponse(422, "invalid event"),
		},
		ExpectError: true,
	}
	cmd := newCalendarCreateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestCalendarUpdateCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update calendar event - API error",
		Args: []string{"10", "--title", "Broken"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/calendar_events/10": cmdtest.NewErrorResponse(404, "event not found"),
		},
		ExpectError: true,
	}
	cmd := newCalendarUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestCalendarDeleteCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete calendar event - API error",
		Args: []string{"10", "--force"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/calendar_events/10": cmdtest.NewErrorResponse(404, "event not found"),
		},
		ExpectError: true,
	}
	cmd := newCalendarDeleteCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestCalendarListCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list calendar events - API error",
		Args: []string{"--course-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/calendar_events": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmd := newCalendarListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}
