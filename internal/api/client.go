package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	defaultRequestsPerSecond = 5.0
	slowRequestsPerSecond    = 2.0
	verySlowRequestsPerSecond = 1.0
	quotaWarningThreshold    = 0.5  // 50%
	quotaCriticalThreshold   = 0.2  // 20%
)

const (
	// defaultQuotaTotal is the default Canvas API quota if not detected from headers
	defaultQuotaTotal = 700.0
)

// HTTPClient defines the interface for making HTTP requests to Canvas API
// This interface allows for easier testing and mocking of the API client
type HTTPClient interface {
	// Low-level HTTP methods
	Get(ctx context.Context, path string) (*http.Response, error)
	Post(ctx context.Context, path string, body io.Reader) (*http.Response, error)
	Put(ctx context.Context, path string, body io.Reader) (*http.Response, error)
	Delete(ctx context.Context, path string) (*http.Response, error)

	// JSON convenience methods
	GetJSON(ctx context.Context, path string, result interface{}) error
	PostJSON(ctx context.Context, path string, body interface{}, result interface{}) error
	PutJSON(ctx context.Context, path string, body interface{}, result interface{}) error

	// Pagination support
	GetAllPages(ctx context.Context, path string, result interface{}) error

	// Feature detection
	SupportsFeature(feature string) bool
	GetVersion() *CanvasVersion
}

// Client is the Canvas API client
type Client struct {
	httpClient     *http.Client
	baseURL        string
	token          string
	asUserID       int64 // For admin masquerading
	rateLimiter    *AdaptiveRateLimiter
	retryPolicy    *RetryPolicy
	version        *CanvasVersion
	featureChecker *FeatureChecker
	logger         *slog.Logger
	quotaTotal     float64 // Detected or configured quota total
	mu             sync.RWMutex
}

// Ensure Client implements HTTPClient interface
var _ HTTPClient = (*Client)(nil)

// AdaptiveRateLimiter implements adaptive rate limiting based on quota
type AdaptiveRateLimiter struct {
	limiter      *rate.Limiter
	currentRate  float64
	mu           sync.RWMutex
	warningShown map[int]bool
	logger       *slog.Logger
}

// NewAdaptiveRateLimiter creates a new adaptive rate limiter
func NewAdaptiveRateLimiter(requestsPerSecond float64) *AdaptiveRateLimiter {
	return &AdaptiveRateLimiter{
		limiter:      rate.NewLimiter(rate.Limit(requestsPerSecond), 1),
		currentRate:  requestsPerSecond,
		warningShown: make(map[int]bool),
		logger:       slog.Default(),
	}
}

// Wait waits for permission to make a request
func (l *AdaptiveRateLimiter) Wait(ctx context.Context) error {
	return l.limiter.Wait(ctx)
}

// AdjustRate adjusts the rate based on remaining quota
func (l *AdaptiveRateLimiter) AdjustRate(remaining, total float64) {
	l.mu.Lock()
	defer l.mu.Unlock()

	percentage := remaining / total

	// Critical threshold - slow down to 1 req/sec
	if percentage <= quotaCriticalThreshold && l.currentRate > verySlowRequestsPerSecond {
		l.currentRate = verySlowRequestsPerSecond
		l.limiter.SetLimit(rate.Limit(verySlowRequestsPerSecond))
		if !l.warningShown[20] {
			l.logger.Warn("⚠️  API rate limit: 20% remaining, slowing to 1 req/sec")
			l.warningShown[20] = true
		}
	} else if percentage <= quotaWarningThreshold && l.currentRate > slowRequestsPerSecond {
		// Warning threshold - slow down to 2 req/sec
		l.currentRate = slowRequestsPerSecond
		l.limiter.SetLimit(rate.Limit(slowRequestsPerSecond))
		if !l.warningShown[50] {
			l.logger.Warn("⚠️  API rate limit: 50% remaining, slowing to 2 req/sec")
			l.warningShown[50] = true
		}
	} else if percentage > quotaWarningThreshold && l.currentRate < defaultRequestsPerSecond {
		// Back to normal - reset to 5 req/sec
		l.currentRate = defaultRequestsPerSecond
		l.limiter.SetLimit(rate.Limit(defaultRequestsPerSecond))
		l.warningShown = make(map[int]bool) // Reset warnings
	}
}

// GetCurrentRate returns the current rate limit
func (l *AdaptiveRateLimiter) GetCurrentRate() float64 {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.currentRate
}

// ClientConfig holds configuration for the API client
type ClientConfig struct {
	BaseURL         string
	Token           string
	RequestsPerSec  float64
	Timeout         time.Duration
	Logger          *slog.Logger
	AsUserID        int64 // For admin masquerading (appends as_user_id param)
}

// NewClient creates a new Canvas API client
func NewClient(config ClientConfig) (*Client, error) {
	if config.BaseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}

	if config.Token == "" {
		return nil, fmt.Errorf("token is required")
	}

	if config.RequestsPerSec == 0 {
		config.RequestsPerSec = defaultRequestsPerSecond
	}

	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	if config.Logger == nil {
		config.Logger = slog.Default()
	}

	// Create HTTP transport with connection pool configuration
	transport := &http.Transport{
		MaxIdleConns:          10,
		MaxIdleConnsPerHost:   5,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &Client{
		httpClient: &http.Client{
			Timeout:   config.Timeout,
			Transport: transport,
		},
		baseURL:     config.BaseURL,
		token:       config.Token,
		asUserID:    config.AsUserID,
		rateLimiter: NewAdaptiveRateLimiter(config.RequestsPerSec),
		retryPolicy: DefaultRetryPolicy(),
		logger:      config.Logger,
		quotaTotal:  defaultQuotaTotal, // Will be updated from headers if available
	}

	// Detect Canvas version
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	version, err := DetectCanvasVersion(ctx, client.httpClient, config.BaseURL)
	if err != nil {
		config.Logger.Warn("Failed to detect Canvas version", "error", err)
		// Use unknown version
		version = &CanvasVersion{Major: 999, Minor: 999, Patch: 999, Raw: "unknown"}
	}

	client.version = version
	client.featureChecker = NewFeatureChecker(version)

	return client, nil
}

// GetVersion returns the detected Canvas version
func (c *Client) GetVersion() *CanvasVersion {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.version
}

// SupportsFeature checks if a feature is supported
func (c *Client) SupportsFeature(feature string) bool {
	return c.featureChecker.SupportsFeature(feature)
}

// doRequest performs an HTTP request with rate limiting and retry logic
func (c *Client) doRequest(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	// Wait for rate limiter
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	// Build full URL with masquerade parameter if set
	fullURL := c.baseURL + path
	if c.asUserID > 0 {
		// Append as_user_id for masquerading
		parsedURL, err := url.Parse(fullURL)
		if err == nil {
			query := parsedURL.Query()
			query.Set("as_user_id", strconv.FormatInt(c.asUserID, 10))
			parsedURL.RawQuery = query.Encode()
			fullURL = parsedURL.String()
		}
	}

	// Execute with retry
	return c.retryPolicy.ExecuteWithRetry(ctx, func() (*http.Response, error) {
		req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Set headers
		req.Header.Set("Authorization", "Bearer "+c.token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		// Make request
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		// Update rate limiter based on response headers
		c.updateRateLimitFromHeaders(resp)

		// Check for API errors
		if resp.StatusCode >= 400 {
			err := ParseAPIError(resp)
			resp.Body.Close()
			return resp, err
		}

		return resp, nil
	})
}

// updateRateLimitFromHeaders updates the rate limiter based on response headers
func (c *Client) updateRateLimitFromHeaders(resp *http.Response) {
	// Parse X-Rate-Limit-Remaining header
	remaining := resp.Header.Get("X-Rate-Limit-Remaining")
	if remaining != "" {
		if remainingFloat, err := strconv.ParseFloat(remaining, 64); err == nil {
			c.rateLimiter.AdjustRate(remainingFloat, c.quotaTotal)
		}
	}
}

// SetQuotaTotal allows configuring the Canvas API quota total
// This is useful when the actual quota differs from the default (700)
func (c *Client) SetQuotaTotal(quota float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.quotaTotal = quota
}

// GetQuotaTotal returns the current quota total setting
func (c *Client) GetQuotaTotal() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.quotaTotal
}

// Get performs a GET request
func (c *Client) Get(ctx context.Context, path string) (*http.Response, error) {
	return c.doRequest(ctx, http.MethodGet, path, nil)
}

// Post performs a POST request
func (c *Client) Post(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	return c.doRequest(ctx, http.MethodPost, path, body)
}

// Put performs a PUT request
func (c *Client) Put(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	return c.doRequest(ctx, http.MethodPut, path, body)
}

// Delete performs a DELETE request
func (c *Client) Delete(ctx context.Context, path string) (*http.Response, error) {
	return c.doRequest(ctx, http.MethodDelete, path, nil)
}

// GetJSON performs a GET request and decodes JSON response
func (c *Client) GetJSON(ctx context.Context, path string, result interface{}) error {
	resp, err := c.Get(ctx, path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(result)
}

// PostJSON performs a POST request with JSON body and decodes JSON response
func (c *Client) PostJSON(ctx context.Context, path string, body interface{}, result interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	resp, err := c.Post(ctx, path, bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}

	return nil
}

// PutJSON performs a PUT request with JSON body and decodes JSON response
func (c *Client) PutJSON(ctx context.Context, path string, body interface{}, result interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	resp, err := c.Put(ctx, path, bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}

	return nil
}

// GetAllPages fetches all pages of a paginated endpoint
func (c *Client) GetAllPages(ctx context.Context, path string, result interface{}) error {
	var allResults []json.RawMessage
	currentURL := path

	for currentURL != "" {
		resp, err := c.Get(ctx, currentURL)
		if err != nil {
			return err
		}

		var pageResults []json.RawMessage
		if err := json.NewDecoder(resp.Body).Decode(&pageResults); err != nil {
			resp.Body.Close()
			return fmt.Errorf("failed to decode response: %w", err)
		}
		resp.Body.Close()

		allResults = append(allResults, pageResults...)

		// Check for next page
		links := ParsePaginationLinks(resp)
		if links.HasNextPage() {
			// Extract path from full URL
			nextURL, err := url.Parse(links.Next)
			if err != nil {
				return fmt.Errorf("failed to parse next URL: %w", err)
			}
			// Handle empty query string properly to avoid trailing '?'
			if nextURL.RawQuery != "" {
				currentURL = nextURL.Path + "?" + nextURL.RawQuery
			} else {
				currentURL = nextURL.Path
			}
		} else {
			currentURL = ""
		}
	}

	// Marshal all results and unmarshal into result
	allJSON, err := json.Marshal(allResults)
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}

	return json.Unmarshal(allJSON, result)
}
