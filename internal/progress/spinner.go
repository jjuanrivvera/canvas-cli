package progress

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/jjuanrivvera/canvas-cli/internal/terminal"
)

// Spinner displays an animated spinner on stderr while a long operation runs.
// It is TTY-aware: when stderr is not a terminal (piped output), all
// operations become no-ops so machine-readable output is never polluted.
type Spinner struct {
	mu      sync.Mutex
	wg      sync.WaitGroup
	message string
	active  bool
	done    chan struct{}
	frames  []string
}

// New creates a new Spinner with the given initial message.
func New(message string) *Spinner {
	return &Spinner{
		message: message,
		frames:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	}
}

// Start begins the spinner animation. Safe to call multiple times;
// subsequent calls are no-ops while the spinner is already running.
func (s *Spinner) Start() {
	if !terminal.IsStderrTerminal() {
		return
	}

	s.mu.Lock()
	if s.active {
		s.mu.Unlock()
		return
	}
	s.active = true
	s.done = make(chan struct{})
	s.wg.Add(1)
	s.mu.Unlock()

	go s.run()
}

// Stop halts the spinner and clears the line.
func (s *Spinner) Stop() {
	s.mu.Lock()
	if !s.active {
		s.mu.Unlock()
		return
	}
	s.active = false
	close(s.done)
	s.mu.Unlock()

	// Wait for goroutine to exit before clearing the line
	s.wg.Wait()

	// Clear the spinner line
	if terminal.IsStderrTerminal() {
		fmt.Fprintf(os.Stderr, "\r\033[K")
	}
}

// UpdateMessage changes the spinner message while it is running.
func (s *Spinner) UpdateMessage(msg string) {
	s.mu.Lock()
	s.message = msg
	s.mu.Unlock()
}

func (s *Spinner) run() {
	defer s.wg.Done()

	ticker := time.NewTicker(80 * time.Millisecond)
	defer ticker.Stop()

	i := 0
	for {
		select {
		case <-s.done:
			return
		case <-ticker.C:
			s.mu.Lock()
			msg := s.message
			s.mu.Unlock()

			frame := s.frames[i%len(s.frames)]
			fmt.Fprintf(os.Stderr, "\r\033[K%s %s", frame, msg)
			i++
		}
	}
}
