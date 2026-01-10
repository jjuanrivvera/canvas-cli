package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"time"
)

const (
	maxRetries     = 3
	initialBackoff = 1 * time.Second
	maxBackoff     = 8 * time.Second
)

// RetryPolicy defines retry behavior for HTTP requests
type RetryPolicy struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	Logger         *slog.Logger
}

// DefaultRetryPolicy returns the default retry policy
func DefaultRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		MaxRetries:     maxRetries,
		InitialBackoff: initialBackoff,
		MaxBackoff:     maxBackoff,
		Logger:         slog.Default(),
	}
}

// ShouldRetry determines if a request should be retried based on the response
func (p *RetryPolicy) ShouldRetry(resp *http.Response, err error) bool {
	// Don't retry on context cancellation (use errors.Is for proper wrapped error handling)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return false
		}
		// Retry on network errors
		return true
	}

	// Retry on these status codes
	switch resp.StatusCode {
	case http.StatusTooManyRequests:
		return true
	case http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	}

	return false
}

// GetBackoff calculates the backoff duration for a given attempt
func (p *RetryPolicy) GetBackoff(attempt int) time.Duration {
	// Exponential backoff: 1s, 2s, 4s
	backoff := time.Duration(math.Pow(2, float64(attempt))) * p.InitialBackoff
	if backoff > p.MaxBackoff {
		backoff = p.MaxBackoff
	}
	return backoff
}

// ExecuteWithRetry executes a function with retry logic
func (p *RetryPolicy) ExecuteWithRetry(ctx context.Context, fn func() (*http.Response, error)) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= p.MaxRetries; attempt++ {
		resp, err = fn()

		// Check if we should retry
		if !p.ShouldRetry(resp, err) {
			return resp, err
		}

		// Don't sleep after the last attempt
		if attempt == p.MaxRetries {
			break
		}

		backoff := p.GetBackoff(attempt)
		p.Logger.Warn("Request failed, retrying",
			"attempt", attempt+1,
			"max_retries", p.MaxRetries,
			"backoff", backoff,
			"error", err,
		)

		// Wait before retrying
		select {
		case <-ctx.Done():
			return resp, ctx.Err()
		case <-time.After(backoff):
			// Continue to next attempt
		}
	}

	if err != nil {
		return resp, fmt.Errorf("request failed after %d retries: %w", p.MaxRetries, err)
	}

	return resp, err
}
