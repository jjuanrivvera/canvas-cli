package batch

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Processor handles batch processing of items with concurrency control
type Processor struct {
	workers     int
	stopOnError bool
	progress    ProgressReporter
}

// New creates a new batch processor
func New(workers int, stopOnError bool, progress ProgressReporter) *Processor {
	if workers <= 0 {
		workers = 1
	}

	return &Processor{
		workers:     workers,
		stopOnError: stopOnError,
		progress:    progress,
	}
}

// ProcessFunc is a function that processes a single item
type ProcessFunc func(ctx context.Context, item interface{}) error

// Result represents the result of processing a single item
type Result struct {
	Item  interface{}
	Error error
	Index int
}

// Process processes a batch of items concurrently
func (p *Processor) Process(ctx context.Context, items []interface{}, fn ProcessFunc) (*Summary, error) {
	if len(items) == 0 {
		return &Summary{Total: 0}, nil
	}

	// Create channels
	jobs := make(chan job, len(items))
	results := make(chan Result, len(items))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < p.workers; i++ {
		wg.Add(1)
		go p.worker(ctx, &wg, jobs, results, fn)
	}

	// Send jobs
	for i, item := range items {
		jobs <- job{
			item:  item,
			index: i,
		}
	}
	close(jobs)

	// Wait for workers to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	summary := &Summary{
		Total:   len(items),
		Results: make([]Result, 0, len(items)),
	}

	start := time.Now()

	for result := range results {
		summary.Results = append(summary.Results, result)

		if result.Error != nil {
			summary.Failed++
			if p.stopOnError {
				// Cancel remaining jobs
				break
			}
		} else {
			summary.Succeeded++
		}

		// Report progress
		if p.progress != nil {
			p.progress.Report(summary.Succeeded+summary.Failed, summary.Total)
		}
	}

	summary.Duration = time.Since(start)

	// If we stopped early due to error, return error
	if p.stopOnError && summary.Failed > 0 {
		return summary, fmt.Errorf("batch processing stopped due to error")
	}

	return summary, nil
}

// job represents a single job to process
type job struct {
	item  interface{}
	index int
}

// worker processes jobs from the jobs channel
func (p *Processor) worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan job, results chan<- Result, fn ProcessFunc) {
	defer wg.Done()

	for j := range jobs {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			results <- Result{
				Item:  j.item,
				Error: ctx.Err(),
				Index: j.index,
			}
			continue
		default:
		}

		// Process the item
		err := fn(ctx, j.item)

		results <- Result{
			Item:  j.item,
			Error: err,
			Index: j.index,
		}
	}
}

// Summary represents the results of batch processing
type Summary struct {
	Total     int
	Succeeded int
	Failed    int
	Duration  time.Duration
	Results   []Result
}

// SuccessRate returns the success rate as a percentage
func (s *Summary) SuccessRate() float64 {
	if s.Total == 0 {
		return 0
	}
	return float64(s.Succeeded) / float64(s.Total) * 100
}

// FailedItems returns a slice of items that failed processing
func (s *Summary) FailedItems() []interface{} {
	items := make([]interface{}, 0, s.Failed)
	for _, result := range s.Results {
		if result.Error != nil {
			items = append(items, result.Item)
		}
	}
	return items
}

// Errors returns a slice of all errors encountered
func (s *Summary) Errors() []error {
	errors := make([]error, 0, s.Failed)
	for _, result := range s.Results {
		if result.Error != nil {
			errors = append(errors, result.Error)
		}
	}
	return errors
}

// String returns a human-readable summary
func (s *Summary) String() string {
	return fmt.Sprintf("Total: %d, Succeeded: %d, Failed: %d, Success Rate: %.1f%%, Duration: %s",
		s.Total, s.Succeeded, s.Failed, s.SuccessRate(), s.Duration)
}

// ProgressReporter is an interface for reporting progress
type ProgressReporter interface {
	Report(current, total int)
}

// ConsoleProgress reports progress to the console
type ConsoleProgress struct {
	lastUpdate     time.Time
	updateInterval time.Duration
}

// NewConsoleProgress creates a new console progress reporter
func NewConsoleProgress(updateInterval time.Duration) *ConsoleProgress {
	return &ConsoleProgress{
		updateInterval: updateInterval,
	}
}

// Report reports progress to the console
func (c *ConsoleProgress) Report(current, total int) {
	now := time.Now()
	if now.Sub(c.lastUpdate) < c.updateInterval && current < total {
		return
	}

	c.lastUpdate = now
	percentage := float64(current) / float64(total) * 100
	fmt.Printf("\rProgress: %d/%d (%.1f%%)  ", current, total, percentage)

	if current == total {
		fmt.Println() // New line when complete
	}
}
