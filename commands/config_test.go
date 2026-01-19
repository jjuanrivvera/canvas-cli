package commands

import (
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestConfigShowCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name:          "get config successfully",
			Args:          []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				// Config doesn't make API calls, it reads local config
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				// Output will show config file path or values
				// Just check it doesn't error
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newConfigShowCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
