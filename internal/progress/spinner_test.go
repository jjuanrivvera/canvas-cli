package progress

import (
	"testing"
	"time"
)

func TestSpinner_StartStop(t *testing.T) {
	s := New("Loading...")
	// In CI/test environments stderr is not a TTY, so Start is a no-op.
	// This verifies there are no panics or races.
	s.Start()
	time.Sleep(50 * time.Millisecond)
	s.Stop()
}

func TestSpinner_DoubleStart(t *testing.T) {
	s := New("Loading...")
	s.Start()
	s.Start() // second start should be a no-op
	s.Stop()
}

func TestSpinner_DoubleStop(t *testing.T) {
	s := New("Loading...")
	s.Start()
	s.Stop()
	s.Stop() // second stop should be a no-op
}

func TestSpinner_UpdateMessage(t *testing.T) {
	s := New("Loading...")
	s.Start()
	s.UpdateMessage("Still loading...")
	s.Stop()
}

func TestSpinner_StopWithoutStart(t *testing.T) {
	s := New("Loading...")
	s.Stop() // should not panic
}
