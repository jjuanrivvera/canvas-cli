package testing

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestMockJSONResponse(t *testing.T) {
	type TestData struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	response := MockJSONResponse(TestData{ID: 1, Name: "test"})
	expected := `{"id":1,"name":"test"}`

	if response != expected {
		t.Errorf("MockJSONResponse() = %q, want %q", response, expected)
	}
}

func TestMockErrorResponse(t *testing.T) {
	response := MockErrorResponse("not found")
	expected := `{"errors":[{"message":"not found"}]}`

	if response != expected {
		t.Errorf("MockErrorResponse() = %q, want %q", response, expected)
	}
}

func TestRunCommandTest(t *testing.T) {
	// Create a simple test command
	testCmd := &cobra.Command{
		Use:   "test",
		Short: "Test command",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Println("Test output")
			return nil
		},
	}

	tc := CommandTestCase{
		Name: "simple command",
		Args: []string{},
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Test output") {
				t.Errorf("Expected output to contain 'Test output', got %q", output)
			}
		},
	}

	RunCommandTest(t, testCmd, tc)
}
