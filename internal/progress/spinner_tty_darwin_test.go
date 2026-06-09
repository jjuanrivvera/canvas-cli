//go:build darwin

package progress

// Darwin-specific tests that exercise the TTY-gated branches in Start and Stop
// by temporarily replacing os.Stderr with a pseudo-terminal slave.
//
// We use raw syscalls (stdlib only) to create the PTY:
//  - Open /dev/ptmx to get the master fd.
//  - Use TIOCPTYGNAME ioctl to retrieve the slave device path.
//  - Use TIOCPTYUNLK ioctl to unlock the slave.
//  - Open the slave; term.IsTerminal returns true for it.
//  - Swap os.Stderr, run the spinner, then restore.

import (
	"os"
	"syscall"
	"testing"
	"time"
	"unsafe"
)

// Darwin ioctl constants for pseudo-terminals (from sys/ttycom.h).
const (
	tiocPtyGname = 0x40807453 // get slave PTY name
	tiocPtyUnlk  = 0x20007452 // unlock slave PTY
)

// openSlavePTY opens /dev/ptmx and returns an *os.File for the slave PTY.
// The caller must close both the returned *os.File and the raw master fd.
// If PTY creation fails, the test is skipped.
func openSlavePTY(t *testing.T) (masterFd int, slave *os.File) {
	t.Helper()

	// Open the PTY master.
	fd, err := syscall.Open("/dev/ptmx", syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_CLOEXEC, 0)
	if err != nil {
		t.Skipf("openSlavePTY: open /dev/ptmx: %v", err)
	}

	// Retrieve the slave device name via TIOCPTYGNAME ioctl.
	var slaveName [128]byte
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(fd),
		tiocPtyGname,
		uintptr(unsafe.Pointer(&slaveName[0])),
	)
	if errno != 0 {
		syscall.Close(fd)
		t.Skipf("openSlavePTY: TIOCPTYGNAME: %v", errno)
	}

	// Unlock the slave via TIOCPTYUNLK ioctl.
	var zero int32
	_, _, errno = syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(fd),
		tiocPtyUnlk,
		uintptr(unsafe.Pointer(&zero)),
	)
	if errno != 0 {
		// Non-fatal: some systems skip the unlock step.
		_ = errno
	}

	// Determine the slave path (null-terminated C string in slaveName).
	slavePathEnd := 0
	for i, b := range slaveName {
		if b == 0 {
			slavePathEnd = i
			break
		}
	}
	slavePath := string(slaveName[:slavePathEnd])
	if slavePath == "" {
		syscall.Close(fd)
		t.Skip("openSlavePTY: empty slave path from TIOCPTYGNAME")
	}

	// Open the slave device.
	slaveFd, openErr := os.OpenFile(slavePath, os.O_RDWR|syscall.O_NOCTTY, 0)
	if openErr != nil {
		syscall.Close(fd)
		t.Skipf("openSlavePTY: open slave %s: %v", slavePath, openErr)
	}

	return fd, slaveFd
}

// TestSpinner_Start_TTYPath exercises the full TTY-gated body of Start
// (goroutine launch) and the TTY-gated clear-line write in Stop.
func TestSpinner_Start_TTYPath(t *testing.T) {
	masterFd, slaveFd := openSlavePTY(t)
	defer syscall.Close(masterFd)
	defer slaveFd.Close()

	origStderr := os.Stderr
	os.Stderr = slaveFd
	defer func() { os.Stderr = origStderr }()

	s := New("working...")
	s.Start()

	if !s.active {
		t.Error("expected active=true after Start with TTY stderr")
	}

	// Allow at least one ticker tick (80 ms) to exercise the run loop write path.
	time.Sleep(120 * time.Millisecond)

	// Stop flips active=false, closes done, waits for the goroutine, and
	// — because IsStderrTerminal() now returns true — prints the clear sequence.
	s.Stop()

	if s.active {
		t.Error("expected active=false after Stop")
	}
}

// TestSpinner_Start_DoubleStart_TTYPath verifies that a second Start() call
// while the spinner is already running is a no-op (active guard inside Start).
func TestSpinner_Start_DoubleStart_TTYPath(t *testing.T) {
	masterFd, slaveFd := openSlavePTY(t)
	defer syscall.Close(masterFd)
	defer slaveFd.Close()

	origStderr := os.Stderr
	os.Stderr = slaveFd
	defer func() { os.Stderr = origStderr }()

	s := New("loading")
	s.Start()
	s.Start() // second Start must be a no-op
	s.Stop()
}
