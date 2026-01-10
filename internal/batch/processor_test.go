package batch

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestProcessor_Process(t *testing.T) {
	processor := New(2, false, nil)

	items := []interface{}{1, 2, 3, 4, 5}

	fn := func(ctx context.Context, item interface{}) error {
		// Simple processing: just return nil
		time.Sleep(10 * time.Millisecond)
		return nil
	}

	summary, err := processor.Process(context.Background(), items, fn)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	if summary.Total != 5 {
		t.Errorf("expected total 5, got %d", summary.Total)
	}

	if summary.Succeeded != 5 {
		t.Errorf("expected 5 successes, got %d", summary.Succeeded)
	}

	if summary.Failed != 0 {
		t.Errorf("expected 0 failures, got %d", summary.Failed)
	}
}

func TestProcessor_ProcessWithErrors(t *testing.T) {
	processor := New(2, false, nil)

	items := []interface{}{1, 2, 3, 4, 5}

	fn := func(ctx context.Context, item interface{}) error {
		// Fail on even numbers
		if item.(int)%2 == 0 {
			return errors.New("even number error")
		}
		return nil
	}

	summary, err := processor.Process(context.Background(), items, fn)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	if summary.Succeeded != 3 {
		t.Errorf("expected 3 successes, got %d", summary.Succeeded)
	}

	if summary.Failed != 2 {
		t.Errorf("expected 2 failures, got %d", summary.Failed)
	}
}

func TestProcessor_StopOnError(t *testing.T) {
	processor := New(2, true, nil)

	items := []interface{}{1, 2, 3, 4, 5}

	fn := func(ctx context.Context, item interface{}) error {
		if item.(int) == 3 {
			return errors.New("error on 3")
		}
		time.Sleep(20 * time.Millisecond)
		return nil
	}

	summary, err := processor.Process(context.Background(), items, fn)

	// Should return error when stop on error is enabled
	if err == nil {
		t.Error("expected error when stopOnError is true")
	}

	if summary.Failed == 0 {
		t.Error("expected at least one failure")
	}
}

func TestProcessor_EmptyItems(t *testing.T) {
	processor := New(2, false, nil)

	items := []interface{}{}

	fn := func(ctx context.Context, item interface{}) error {
		return nil
	}

	summary, err := processor.Process(context.Background(), items, fn)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	if summary.Total != 0 {
		t.Errorf("expected total 0, got %d", summary.Total)
	}
}

func TestProcessor_ContextCancellation(t *testing.T) {
	processor := New(2, false, nil)

	items := []interface{}{1, 2, 3, 4, 5}

	ctx, cancel := context.WithCancel(context.Background())

	fn := func(ctx context.Context, item interface{}) error {
		time.Sleep(50 * time.Millisecond)
		return nil
	}

	// Cancel after a short delay
	go func() {
		time.Sleep(30 * time.Millisecond)
		cancel()
	}()

	summary, _ := processor.Process(ctx, items, fn)

	// Some items should have context cancellation errors
	if summary.Failed == 0 {
		t.Error("expected some failures due to context cancellation")
	}
}

func TestSummary_SuccessRate(t *testing.T) {
	summary := &Summary{
		Total:     10,
		Succeeded: 7,
		Failed:    3,
	}

	rate := summary.SuccessRate()
	expected := 70.0

	if rate != expected {
		t.Errorf("expected success rate %.1f%%, got %.1f%%", expected, rate)
	}
}

func TestSummary_FailedItems(t *testing.T) {
	summary := &Summary{
		Total:     5,
		Succeeded: 3,
		Failed:    2,
		Results: []Result{
			{Item: 1, Error: nil},
			{Item: 2, Error: errors.New("error")},
			{Item: 3, Error: nil},
			{Item: 4, Error: errors.New("error")},
			{Item: 5, Error: nil},
		},
	}

	failedItems := summary.FailedItems()

	if len(failedItems) != 2 {
		t.Errorf("expected 2 failed items, got %d", len(failedItems))
	}

	if failedItems[0] != 2 {
		t.Errorf("expected failed item 2, got %v", failedItems[0])
	}

	if failedItems[1] != 4 {
		t.Errorf("expected failed item 4, got %v", failedItems[1])
	}
}

func TestSummary_Errors(t *testing.T) {
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")

	summary := &Summary{
		Results: []Result{
			{Item: 1, Error: nil},
			{Item: 2, Error: err1},
			{Item: 3, Error: nil},
			{Item: 4, Error: err2},
		},
	}

	errs := summary.Errors()

	if len(errs) != 2 {
		t.Errorf("expected 2 errors, got %d", len(errs))
	}

	if errs[0] != err1 {
		t.Errorf("expected error 1, got %v", errs[0])
	}

	if errs[1] != err2 {
		t.Errorf("expected error 2, got %v", errs[1])
	}
}

func TestConsoleProgress_Report(t *testing.T) {
	progress := NewConsoleProgress(100 * time.Millisecond)

	// These calls should not panic
	progress.Report(1, 10)
	progress.Report(5, 10)
	progress.Report(10, 10)
}

func TestSummary_String(t *testing.T) {
	summary := &Summary{
		Total:     10,
		Succeeded: 8,
		Failed:    2,
		Duration:  5 * time.Second,
	}

	result := summary.String()

	// Check that the string contains key information
	if result == "" {
		t.Error("expected non-empty string")
	}

	// Should contain total, succeeded, failed counts
	expectedSubstrings := []string{
		"Total: 10",
		"Succeeded: 8",
		"Failed: 2",
		"Success Rate: 80.0%",
		"5s",
	}

	for _, expected := range expectedSubstrings {
		if !strings.Contains(result, expected) {
			t.Errorf("expected string to contain %q, got %q", expected, result)
		}
	}
}

func TestNew_DefaultWorkers(t *testing.T) {
	// Test with workers <= 0 (should default to 1)
	processor := New(0, false, nil)
	if processor.workers != 1 {
		t.Errorf("expected workers to be 1, got %d", processor.workers)
	}

	processor = New(-5, false, nil)
	if processor.workers != 1 {
		t.Errorf("expected workers to be 1 with negative input, got %d", processor.workers)
	}
}

func TestNew_CustomWorkers(t *testing.T) {
	processor := New(5, true, nil)
	if processor.workers != 5 {
		t.Errorf("expected workers to be 5, got %d", processor.workers)
	}
	if !processor.stopOnError {
		t.Error("expected stopOnError to be true")
	}
}

func TestSummary_SuccessRate_ZeroTotal(t *testing.T) {
	summary := &Summary{
		Total:     0,
		Succeeded: 0,
		Failed:    0,
	}

	rate := summary.SuccessRate()
	if rate != 0 {
		t.Errorf("expected success rate 0 for zero total, got %.2f", rate)
	}
}

func TestSummary_SuccessRate_NonZero(t *testing.T) {
	tests := []struct {
		name      string
		total     int
		succeeded int
		expected  float64
	}{
		{
			name:      "100% success",
			total:     10,
			succeeded: 10,
			expected:  100.0,
		},
		{
			name:      "50% success",
			total:     10,
			succeeded: 5,
			expected:  50.0,
		},
		{
			name:      "0% success",
			total:     10,
			succeeded: 0,
			expected:  0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary := &Summary{
				Total:     tt.total,
				Succeeded: tt.succeeded,
				Failed:    tt.total - tt.succeeded,
			}

			rate := summary.SuccessRate()
			if rate != tt.expected {
				t.Errorf("expected success rate %.1f%%, got %.1f%%", tt.expected, rate)
			}
		})
	}
}

// Benchmark tests

func BenchmarkProcessor_Process_Serial(b *testing.B) {
	processor := New(1, false, nil)
	items := make([]interface{}, 100)
	for i := range items {
		items[i] = i
	}

	fn := func(ctx context.Context, item interface{}) error {
		return nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = processor.Process(context.Background(), items, fn)
	}
}

func BenchmarkProcessor_Process_Concurrent(b *testing.B) {
	processor := New(10, false, nil)
	items := make([]interface{}, 100)
	for i := range items {
		items[i] = i
	}

	fn := func(ctx context.Context, item interface{}) error {
		return nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = processor.Process(context.Background(), items, fn)
	}
}

func BenchmarkProcessor_Process_WithIO(b *testing.B) {
	processor := New(5, false, nil)
	items := make([]interface{}, 10)
	for i := range items {
		items[i] = i
	}

	fn := func(ctx context.Context, item interface{}) error {
		time.Sleep(1 * time.Millisecond) // Simulate I/O
		return nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = processor.Process(context.Background(), items, fn)
	}
}

func BenchmarkSummary_SuccessRate(b *testing.B) {
	summary := &Summary{
		Total:     1000,
		Succeeded: 750,
		Failed:    250,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = summary.SuccessRate()
	}
}

func BenchmarkSummary_Errors(b *testing.B) {
	summary := &Summary{
		Results: make([]Result, 100),
	}
	for i := range summary.Results {
		if i%2 == 0 {
			summary.Results[i] = Result{Item: i, Error: errors.New("error")}
		} else {
			summary.Results[i] = Result{Item: i, Error: nil}
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = summary.Errors()
	}
}
