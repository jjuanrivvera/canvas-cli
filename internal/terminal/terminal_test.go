package terminal

import (
	"testing"
)

func TestWidth(t *testing.T) {
	w := Width()
	if w <= 0 {
		t.Errorf("Width() returned %d, expected > 0", w)
	}
}

func TestIsTerminal(t *testing.T) {
	// Should not panic regardless of whether running in a terminal
	_ = IsStdoutTerminal()
	_ = IsStderrTerminal()
	_ = IsStdinTerminal()
}

func TestSupportsColor(t *testing.T) {
	// Should not panic
	_ = SupportsColor()
}
