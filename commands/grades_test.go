package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestGradesHistoryCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get gradebook history successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts": cmdtest.NewMockResponse(`[]`),
				"/api/v1/courses/1/gradebook_history/days": cmdtest.NewMockResponse(`[
					{
						"date": "2024-01-15",
						"graders": []
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "2024-01-15") {
					t.Error("Expected '2024-01-15' in output")
				}
			},
		},
		{
			Name: "get gradebook history - empty response",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts":                         cmdtest.NewMockResponse(`[]`),
				"/api/v1/courses/1/gradebook_history/days": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No gradebook history found",
		},
		{
			Name:        "get gradebook history - missing course ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newGradesHistoryCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestGradesFeedCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get gradebook feed successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/gradebook_history/feed": cmdtest.NewMockResponse(`[
					{
						"assignment_id": 100,
						"grader_id": 10,
						"student_id": 20,
						"score": 95
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "100") {
					t.Error("Expected '100' in output")
				}
			},
		},
		{
			Name: "get gradebook feed - empty response",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                        courseMock,
				"/api/v1/courses/1/gradebook_history/feed": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No gradebook entries found",
		},
		{
			Name:        "get gradebook feed - missing course ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newGradesFeedCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestGradesColumnsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list custom columns successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/custom_gradebook_columns": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"title": "Notes",
						"position": 1,
						"teacher_notes": false,
						"read_only": false,
						"hidden": false
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Notes") {
					t.Error("Expected 'Notes' in output")
				}
			},
		},
		{
			Name: "list custom columns - empty response",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                          courseMock,
				"/api/v1/courses/1/custom_gradebook_columns": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No custom columns found",
		},
		{
			Name:        "list custom columns - missing course ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newGradesColumnsListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestGradesColumnsGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get custom column successfully",
			Args: []string{"10", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/custom_gradebook_columns/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"title": "Attendance",
					"position": 2,
					"teacher_notes": false
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Attendance") {
					t.Error("Expected 'Attendance' in output")
				}
			},
		},
		{
			Name:        "get custom column - missing column ID",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "get custom column - missing course ID",
			Args:        []string{"10"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newGradesColumnsGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestGradesColumnsCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create custom column successfully",
			Args: []string{"--course-id", "1", "--title", "Participation"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/custom_gradebook_columns": cmdtest.NewMockResponse(`{
					"id": 20,
					"title": "Participation",
					"position": 1,
					"teacher_notes": false
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Participation") {
					t.Error("Expected 'Participation' in output")
				}
			},
		},
		{
			Name:        "create custom column - missing course ID",
			Args:        []string{"--title", "New Column"},
			ExpectError: true,
		},
		{
			Name:        "create custom column - missing title",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newGradesColumnsCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestGradesColumnsDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete custom column with confirmation",
			Args: []string{"10", "--course-id", "1", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/custom_gradebook_columns/10": cmdtest.NewMockResponse(`{
					"id": 10,
					"title": "Old Column"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete custom column - missing column ID",
			Args:        []string{"--course-id", "1", "--force"},
			ExpectError: true,
		},
		{
			Name:        "delete custom column - missing course ID",
			Args:        []string{"10", "--force"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newGradesColumnsDeleteCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestGradesColumnsDataListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list column data successfully",
			Args: []string{"10", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/custom_gradebook_columns/10/data": cmdtest.NewMockResponse(`[
					{
						"user_id": 100,
						"content": "Present"
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Present") {
					t.Error("Expected 'Present' in output")
				}
			},
		},
		{
			Name: "list column data - empty response",
			Args: []string{"10", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/custom_gradebook_columns/10/data": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No column data found",
		},
		{
			Name:        "list column data - missing column ID",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newGradesColumnsDataListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestGradesColumnsDataSetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "set column data successfully",
			Args: []string{"10", "--course-id", "1", "--user-id", "100", "--content", "Present"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/custom_gradebook_columns/10/data/100": cmdtest.NewMockResponse(`{
					"user_id": 100,
					"content": "Present"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Present") {
					t.Error("Expected 'Present' in output")
				}
			},
		},
		{
			Name:        "set column data - missing column ID",
			Args:        []string{"--course-id", "1", "--user-id", "100", "--content", "Present"},
			ExpectError: true,
		},
		{
			Name:        "set column data - missing user ID",
			Args:        []string{"10", "--course-id", "1", "--content", "Present"},
			ExpectError: true,
		},
		{
			Name:        "set column data - missing content",
			Args:        []string{"10", "--course-id", "1", "--user-id", "100"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newGradesColumnsDataSetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
