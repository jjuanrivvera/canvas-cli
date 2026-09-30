package commands

import (
	"log/slog"
	"os"
	"testing"

	"github.com/spf13/cobra"
)

// --quiet says it suppresses informational messages. The packages under
// internal/ log through slog's default logger, so unless the flag raises that
// logger's floor a "quiet" run still narrates itself on stderr — which is what
// put a "Detected Canvas version" line in the middle of the README recording.
func TestQuietSilencesTheDefaultLogger(t *testing.T) {
	original := slog.Default()
	t.Cleanup(func() { slog.SetDefault(original) })
	originalQuiet := quiet
	t.Cleanup(func() { quiet = originalQuiet })

	// Start from a logger that would print everything, so that anything the
	// flag changes is the flag's doing and not the default's.
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
	if !slog.Default().Enabled(t.Context(), slog.LevelInfo) {
		t.Fatal("the starting logger is already quiet; the test would prove nothing")
	}

	quiet = true
	rootCmd.PersistentPreRun(&cobra.Command{}, nil)

	// Quiet is not silent: informational records are dropped, warnings are not.
	if slog.Default().Enabled(t.Context(), slog.LevelInfo) {
		t.Error("with --quiet, an Info record should no longer be enabled")
	}
	if !slog.Default().Enabled(t.Context(), slog.LevelWarn) {
		t.Error("--quiet must not swallow warnings")
	}
}

// Without the flag nothing changes: the default logger keeps whatever level it
// had, so a normal run is as talkative as it ever was.
func TestWithoutQuietTheLoggerIsLeftAlone(t *testing.T) {
	original := slog.Default()
	t.Cleanup(func() { slog.SetDefault(original) })
	originalQuiet := quiet
	t.Cleanup(func() { quiet = originalQuiet })

	mine := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(mine)
	quiet = false
	rootCmd.PersistentPreRun(&cobra.Command{}, nil)

	if !slog.Default().Enabled(t.Context(), slog.LevelDebug) {
		t.Error("a run that did not ask to be quiet should keep its logger")
	}
}
