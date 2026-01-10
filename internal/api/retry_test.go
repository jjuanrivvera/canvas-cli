package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"testing"
	"time"
)

func TestDefaultRetryPolicy(t *testing.T) {
	policy := DefaultRetryPolicy()
	if policy == nil {
		t.Fatal("expected non-nil retry policy")
	}

	if policy.MaxRetries != 3 {
		t.Errorf("expected MaxRetries = 3, got %d", policy.MaxRetries)
	}

	if policy.InitialBackoff != 1*time.Second {
		t.Errorf("expected InitialBackoff = 1s, got %v", policy.InitialBackoff)
	}

	if policy.MaxBackoff != 8*time.Second {
		t.Errorf("expected MaxBackoff = 8s, got %v", policy.MaxBackoff)
	}

	if policy.Logger == nil {
		t.Error("expected non-nil logger")
	}
}

func TestRetryPolicy_ShouldRetry(t *testing.T) {
	policy := DefaultRetryPolicy()

	tests := []struct {
		name       string
		statusCode int
		err        error
		want       bool
	}{
		{
			name:       "429 Too Many Requests",
			statusCode: http.StatusTooManyRequests,
			err:        nil,
			want:       true,
		},
		{
			name:       "500 Internal Server Error",
			statusCode: http.StatusInternalServerError,
			err:        nil,
			want:       true,
		},
		{
			name:       "502 Bad Gateway",
			statusCode: http.StatusBadGateway,
			err:        nil,
			want:       true,
		},
		{
			name:       "503 Service Unavailable",
			statusCode: http.StatusServiceUnavailable,
			err:        nil,
			want:       true,
		},
		{
			name:       "504 Gateway Timeout",
			statusCode: http.StatusGatewayTimeout,
			err:        nil,
			want:       true,
		},
		{
			name:       "200 OK",
			statusCode: http.StatusOK,
			err:        nil,
			want:       false,
		},
		{
			name:       "400 Bad Request",
			statusCode: http.StatusBadRequest,
			err:        nil,
			want:       false,
		},
		{
			name:       "404 Not Found",
			statusCode: http.StatusNotFound,
			err:        nil,
			want:       false,
		},
		{
			name:       "network error",
			statusCode: 0,
			err:        errors.New("network error"),
			want:       true,
		},
		{
			name:       "context canceled",
			statusCode: 0,
			err:        context.Canceled,
			want:       false,
		},
		{
			name:       "context deadline exceeded",
			statusCode: 0,
			err:        context.DeadlineExceeded,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp *http.Response
			if tt.statusCode != 0 {
				resp = &http.Response{
					StatusCode: tt.statusCode,
				}
			}

			got := policy.ShouldRetry(resp, tt.err)
			if got != tt.want {
				t.Errorf("ShouldRetry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRetryPolicy_GetBackoff(t *testing.T) {
	policy := &RetryPolicy{
		InitialBackoff: 1 * time.Second,
		MaxBackoff:     8 * time.Second,
	}

	tests := []struct {
		name    string
		attempt int
		want    time.Duration
	}{
		{
			name:    "attempt 0",
			attempt: 0,
			want:    1 * time.Second, // 2^0 * 1s = 1s
		},
		{
			name:    "attempt 1",
			attempt: 1,
			want:    2 * time.Second, // 2^1 * 1s = 2s
		},
		{
			name:    "attempt 2",
			attempt: 2,
			want:    4 * time.Second, // 2^2 * 1s = 4s
		},
		{
			name:    "attempt 3",
			attempt: 3,
			want:    8 * time.Second, // 2^3 * 1s = 8s (capped at MaxBackoff)
		},
		{
			name:    "attempt 4",
			attempt: 4,
			want:    8 * time.Second, // Would be 16s, but capped at 8s
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := policy.GetBackoff(tt.attempt)
			if got != tt.want {
				t.Errorf("GetBackoff(%d) = %v, want %v", tt.attempt, got, tt.want)
			}
		})
	}
}

func TestRetryPolicy_ExecuteWithRetry_Success(t *testing.T) {
	policy := &RetryPolicy{
		MaxRetries:     3,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     50 * time.Millisecond,
		Logger:         slog.Default(),
	}

	ctx := context.Background()
	callCount := 0

	fn := func() (*http.Response, error) {
		callCount++
		return &http.Response{StatusCode: http.StatusOK}, nil
	}

	resp, err := policy.ExecuteWithRetry(ctx, fn)
	if err != nil {
		t.Fatalf("ExecuteWithRetry() error = %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if callCount != 1 {
		t.Errorf("expected 1 call, got %d", callCount)
	}
}

func TestRetryPolicy_ExecuteWithRetry_RetriesOnError(t *testing.T) {
	policy := &RetryPolicy{
		MaxRetries:     2,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     50 * time.Millisecond,
		Logger:         slog.Default(),
	}

	ctx := context.Background()
	callCount := 0

	fn := func() (*http.Response, error) {
		callCount++
		if callCount < 3 {
			// Fail first 2 attempts
			return nil, errors.New("temporary error")
		}
		// Succeed on 3rd attempt
		return &http.Response{StatusCode: http.StatusOK}, nil
	}

	resp, err := policy.ExecuteWithRetry(ctx, fn)
	if err != nil {
		t.Fatalf("ExecuteWithRetry() error = %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if callCount != 3 {
		t.Errorf("expected 3 calls (1 initial + 2 retries), got %d", callCount)
	}
}

func TestRetryPolicy_ExecuteWithRetry_MaxRetriesExceeded(t *testing.T) {
	policy := &RetryPolicy{
		MaxRetries:     2,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     50 * time.Millisecond,
		Logger:         slog.Default(),
	}

	ctx := context.Background()
	callCount := 0

	fn := func() (*http.Response, error) {
		callCount++
		return nil, errors.New("persistent error")
	}

	_, err := policy.ExecuteWithRetry(ctx, fn)
	if err == nil {
		t.Error("expected error after max retries")
	}

	// Should call 1 initial + 2 retries = 3 times
	if callCount != 3 {
		t.Errorf("expected 3 calls, got %d", callCount)
	}
}

func TestRetryPolicy_ExecuteWithRetry_ContextCanceled(t *testing.T) {
	policy := &RetryPolicy{
		MaxRetries:     3,
		InitialBackoff: 100 * time.Millisecond,
		MaxBackoff:     500 * time.Millisecond,
		Logger:         slog.Default(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	callCount := 0

	fn := func() (*http.Response, error) {
		callCount++
		if callCount == 1 {
			// Cancel context after first call
			cancel()
		}
		return nil, errors.New("error")
	}

	_, err := policy.ExecuteWithRetry(ctx, fn)
	if err == nil {
		t.Error("expected error from canceled context")
	}

	// Should have made first call, then detected context cancellation
	if callCount > 2 {
		t.Errorf("expected at most 2 calls, got %d", callCount)
	}
}

func TestRetryPolicy_ExecuteWithRetry_NoRetryOnNonRetryableStatus(t *testing.T) {
	policy := &RetryPolicy{
		MaxRetries:     3,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     50 * time.Millisecond,
		Logger:         slog.Default(),
	}

	ctx := context.Background()
	callCount := 0

	fn := func() (*http.Response, error) {
		callCount++
		return &http.Response{StatusCode: http.StatusBadRequest}, nil
	}

	resp, err := policy.ExecuteWithRetry(ctx, fn)
	if err != nil {
		t.Fatalf("ExecuteWithRetry() error = %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	// Should only call once (no retries for 400)
	if callCount != 1 {
		t.Errorf("expected 1 call, got %d", callCount)
	}
}
