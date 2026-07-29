package commands

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
	"github.com/jjuanrivvera/canvas-cli/internal/config"
	"github.com/jjuanrivvera/canvas-cli/internal/diagnostics"
)

// fakeRunner is a deterministic diagnostics.Runner used to drive the doctor
// command's output/formatting and failure-handling paths without touching the
// network, host filesystem, or real credentials.
type fakeRunner struct {
	report *diagnostics.Report
	err    error
}

func (f *fakeRunner) Run(_ context.Context) (*diagnostics.Report, error) {
	return f.report, f.err
}

// withDoctorRunner swaps the package-level newDoctorRunner seam for one that
// returns the given runner, restoring the original when the test finishes.
func withDoctorRunner(t *testing.T, runner diagnostics.Runner) {
	t.Helper()
	orig := newDoctorRunner
	newDoctorRunner = func(_ *config.Config, _ *api.Client) diagnostics.Runner {
		return runner
	}
	t.Cleanup(func() { newDoctorRunner = orig })
}

// healthyReport builds a deterministic all-green report covering every status
// the human formatter renders (PASS/WARN/SKIP), with no failures.
func healthyReport() *diagnostics.Report {
	return &diagnostics.Report{
		Checks: []diagnostics.Check{
			{Name: "Environment", Description: "System environment and runtime", Status: diagnostics.StatusPass, Message: "OS: linux", Duration: time.Millisecond},
			{Name: "Configuration", Description: "Configuration file and settings", Status: diagnostics.StatusPass, Message: "Instance: test", Duration: time.Millisecond},
			{Name: "Connectivity", Description: "Network connectivity to Canvas", Status: diagnostics.StatusSkipped, Message: "Configuration not available", Duration: time.Millisecond},
			{Name: "Permissions", Description: "File and directory permissions", Status: diagnostics.StatusWarning, Message: "insecure permissions", Duration: time.Millisecond},
		},
		Duration:  5 * time.Millisecond,
		PassCount: 2,
		WarnCount: 1,
		SkipCount: 1,
	}
}

// unhealthyReport builds a deterministic report with a failing check.
func unhealthyReport() *diagnostics.Report {
	return &diagnostics.Report{
		Checks: []diagnostics.Check{
			{Name: "Environment", Description: "System environment and runtime", Status: diagnostics.StatusPass, Message: "OS: linux", Duration: time.Millisecond},
			{
				Name:        "Authentication",
				Description: "API token authentication",
				Status:      diagnostics.StatusFail,
				Message:     "Authentication failed: invalid token",
				Error:       context.DeadlineExceeded,
				Duration:    time.Millisecond,
			},
		},
		Duration:  3 * time.Millisecond,
		PassCount: 1,
		FailCount: 1,
	}
}

func TestDoctorCmd(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		runner      diagnostics.Runner
		expectError bool
		validate    func(t *testing.T, output string)
	}{
		{
			name:   "human output, healthy",
			args:   []string{},
			runner: &fakeRunner{report: healthyReport()},
			validate: func(t *testing.T, output string) {
				assertContains(t, output, "Running diagnostics")
				assertContains(t, output, "Diagnostic Report")
				assertContains(t, output, "Environment")
				assertContains(t, output, "Permissions")
				assertContains(t, output, "All checks passed")
				// Non-verbose must not print per-check descriptions.
				if strings.Contains(output, "Description:") {
					t.Errorf("non-verbose output should not contain per-check Description, got:\n%s", output)
				}
			},
		},
		{
			name:   "human output, verbose",
			args:   []string{"--verbose"},
			runner: &fakeRunner{report: healthyReport()},
			validate: func(t *testing.T, output string) {
				assertContains(t, output, "Description: System environment and runtime")
				assertContains(t, output, "Duration:")
			},
		},
		{
			name:        "human output, failing check",
			args:        []string{},
			runner:      &fakeRunner{report: unhealthyReport()},
			expectError: true,
			validate: func(t *testing.T, output string) {
				assertContains(t, output, "Authentication")
				assertContains(t, output, "check(s) failed")
				assertContains(t, output, "1 check(s) failed")
			},
		},
		{
			name:        "verbose failing check shows error line",
			args:        []string{"--verbose"},
			runner:      &fakeRunner{report: unhealthyReport()},
			expectError: true,
			validate: func(t *testing.T, output string) {
				assertContains(t, output, "Error:")
			},
		},
		{
			name:   "json output via --json flag",
			args:   []string{"--json"},
			runner: &fakeRunner{report: healthyReport()},
			validate: func(t *testing.T, output string) {
				var r struct {
					Duration string `json:"duration"`
					Summary  string `json:"summary"`
					Healthy  bool   `json:"healthy"`
					Checks   []struct {
						Name   string `json:"name"`
						Status string `json:"status"`
					} `json:"checks"`
				}
				// The deprecated --json flag makes cobra append a deprecation
				// notice after the JSON; a streaming decoder reads just the
				// first value and ignores the trailing text.
				if err := json.NewDecoder(strings.NewReader(output)).Decode(&r); err != nil {
					t.Fatalf("output is not valid JSON: %v\noutput:\n%s", err, output)
				}
				if !r.Healthy {
					t.Errorf("expected healthy=true, got false")
				}
				if len(r.Checks) != 4 {
					t.Errorf("expected 4 checks in JSON, got %d", len(r.Checks))
				}
				// JSON path must not emit the human "Running diagnostics" banner.
				if strings.Contains(output, "Running diagnostics") {
					t.Errorf("JSON output should not contain human banner")
				}
			},
		},
		{
			name:        "json output, failing check reports healthy=false",
			args:        []string{"--json"},
			runner:      &fakeRunner{report: unhealthyReport()},
			expectError: true,
			validate: func(t *testing.T, output string) {
				var r struct {
					Healthy bool `json:"healthy"`
				}
				if err := json.NewDecoder(strings.NewReader(output)).Decode(&r); err != nil {
					t.Fatalf("output is not valid JSON: %v", err)
				}
				if r.Healthy {
					t.Errorf("expected healthy=false for failing report")
				}
			},
		},
		{
			name:        "runner error surfaces as diagnostic error",
			args:        []string{},
			runner:      &fakeRunner{err: context.DeadlineExceeded},
			expectError: true,
			validate:    func(t *testing.T, _ string) {},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			withDoctorRunner(t, tc.runner)
			cmd := newDoctorCmd()
			cmdtest.RunCommandTest(t, cmd, cmdtest.CommandTestCase{
				Name:           tc.name,
				Args:           tc.args,
				ExpectError:    tc.expectError,
				ValidateOutput: tc.validate,
			})
		})
	}
}

// TestDoctorCmd_OutputJSONViaGlobalFlag covers the -o/--output json code path,
// which maps the global outputFormat onto opts.JSON.
func TestDoctorCmd_OutputJSONViaGlobalFlag(t *testing.T) {
	withDoctorRunner(t, &fakeRunner{report: healthyReport()})

	origFormat := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = origFormat })

	cmd := newDoctorCmd()
	cmdtest.RunCommandTest(t, cmd, cmdtest.CommandTestCase{
		Name: "global -o json",
		Args: []string{},
		ValidateOutput: func(t *testing.T, output string) {
			assertContains(t, output, `"healthy"`)
			if strings.Contains(output, "Running diagnostics") {
				t.Errorf("global -o json should suppress the human banner")
			}
		},
	})
}

// TestDoctorCmd_Live is the single opt-in end-to-end smoke test that exercises
// the real diagnostics runner (real environment/network/filesystem checks). It
// is skipped by default and only runs when CANVAS_DOCTOR_LIVE=1, so the default
// `go test ./commands -run TestDoctorCmd` stays deterministic and host-free.
func TestDoctorCmd_Live(t *testing.T) {
	if os.Getenv("CANVAS_DOCTOR_LIVE") != "1" {
		t.Skip("set CANVAS_DOCTOR_LIVE=1 to run the live doctor smoke test")
	}

	cmd := newDoctorCmd()
	cmdtest.RunCommandTest(t, cmd, cmdtest.CommandTestCase{
		Name: "live doctor",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self": cmdtest.NewMockResponse(`{"id":1,"name":"Test User"}`),
			"/api/v1/courses":    cmdtest.NewMockResponse(`[]`),
		},
		ValidateOutput: func(t *testing.T, output string) {
			assertContains(t, output, "Diagnostic Report")
			assertContains(t, output, "Environment")
		},
	})
}

func assertContains(t *testing.T, output, want string) {
	t.Helper()
	if !strings.Contains(output, want) {
		t.Errorf("expected output to contain %q, got:\n%s", want, output)
	}
}
