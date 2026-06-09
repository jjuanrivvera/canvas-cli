package commands

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// captureCompletionOutput captures os.Stdout during execution of fn.
// GenBashCompletionV2 and friends write to os.Stdout (the real fd), not
// the cobra writer, so we must intercept the fd.
//
// The pipe is drained concurrently: completion scripts are large (the bash
// script is ~16KB) and reading only after fn() returns would deadlock once the
// output exceeds the OS pipe buffer, which is far smaller on Windows than the
// 64KB on Linux/macOS — the write blocks forever with no reader. This was the
// cause of a ~22-minute Windows CI hang.
func captureCompletionOutput(fn func()) string {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	outCh := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(r)
		outCh <- buf.String()
	}()

	fn()
	w.Close()
	os.Stdout = old
	return <-outCh
}

func TestCompletionCmd_Bash(t *testing.T) {
	out := captureCompletionOutput(func() {
		err := runCompletion(completionCmd, []string{"bash"})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if len(out) == 0 {
		t.Error("expected non-empty bash completion output")
	}
	if !strings.Contains(out, "canvas") {
		t.Errorf("expected 'canvas' in bash completion output, got length %d", len(out))
	}
}

func TestCompletionCmd_Zsh(t *testing.T) {
	out := captureCompletionOutput(func() {
		err := runCompletion(completionCmd, []string{"zsh"})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if len(out) == 0 {
		t.Error("expected non-empty zsh completion output")
	}
}

func TestCompletionCmd_Fish(t *testing.T) {
	out := captureCompletionOutput(func() {
		err := runCompletion(completionCmd, []string{"fish"})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if len(out) == 0 {
		t.Error("expected non-empty fish completion output")
	}
}

func TestCompletionCmd_PowerShell(t *testing.T) {
	out := captureCompletionOutput(func() {
		err := runCompletion(completionCmd, []string{"powershell"})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if len(out) == 0 {
		t.Error("expected non-empty PowerShell completion output")
	}
}

func TestRunCompletion_InvalidShell(t *testing.T) {
	err := runCompletion(completionCmd, []string{"tcsh"})
	if err == nil {
		t.Fatal("expected error for unsupported shell")
	}
	if !strings.Contains(err.Error(), "invalid shell") {
		t.Errorf("unexpected error message: %v", err)
	}
}
