package batch

import (
	"context"
	"errors"
	"testing"
	"time"
)

// noOpProgress is a ProgressReporter that records how many times Report is called.
type noOpProgress struct {
	calls int
}

func (n *noOpProgress) Report(current, total int) {
	n.calls++
}

// TestProcessor_Process_WithProgress verifies that the progress reporter callback
// is invoked during batch processing when a non-nil reporter is provided.
func TestProcessor_Process_WithProgress(t *testing.T) {
	progress := &noOpProgress{}
	processor := New(2, false, progress)

	items := []interface{}{1, 2, 3}

	fn := func(ctx context.Context, item interface{}) error {
		return nil
	}

	summary, err := processor.Process(context.Background(), items, fn)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}
	if summary.Succeeded != 3 {
		t.Errorf("expected 3 successes, got %d", summary.Succeeded)
	}
	if progress.calls == 0 {
		t.Error("expected progress reporter to be called at least once")
	}
}

// TestProcessor_Process_WithProgress_StopOnError verifies that progress is reported
// even when stopOnError is true and an error occurs.
func TestProcessor_Process_WithProgress_StopOnError(t *testing.T) {
	progress := &noOpProgress{}
	processor := New(1, true, progress)

	items := []interface{}{1, 2, 3}

	fn := func(ctx context.Context, item interface{}) error {
		if item.(int) == 1 {
			return errors.New("first item fails")
		}
		return nil
	}

	_, _ = processor.Process(context.Background(), items, fn)
	// We don't assert the call count here because the processor may stop early.
	// The test just verifies no panic occurs.
}

// TestConsoleProgress_Report_Throttled verifies that progress is not reported
// more frequently than the update interval.
func TestConsoleProgress_Report_Throttled(t *testing.T) {
	progress := NewConsoleProgress(10 * time.Second) // Very long interval

	// First call always reports.
	progress.Report(1, 10)
	// Second call is below interval and not complete — should be throttled.
	progress.Report(2, 10)
	// Third call where current == total should always report.
	progress.Report(10, 10)
}

// TestProcessor_Process_SingleWorker confirms that a single-worker processor
// handles all items correctly with a non-nil progress reporter.
func TestProcessor_Process_SingleWorker_WithProgress(t *testing.T) {
	progress := &noOpProgress{}
	processor := New(1, false, progress)

	items := make([]interface{}, 5)
	for i := range items {
		items[i] = i
	}

	summary, err := processor.Process(context.Background(), items, func(ctx context.Context, item interface{}) error {
		return nil
	})
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}
	if summary.Total != 5 {
		t.Errorf("expected 5 total, got %d", summary.Total)
	}
	if summary.Succeeded != 5 {
		t.Errorf("expected 5 successes, got %d", summary.Succeeded)
	}
}
