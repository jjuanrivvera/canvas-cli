package api

import (
	"bytes"
	"context"
	"crypto/md5" // #nosec G501 -- MD5 is used only for cache key hashing, not for any cryptographic purpose
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/time/rate"

	"github.com/jjuanrivvera/canvas-cli/internal/cache"
	"github.com/jjuanrivvera/canvas-cli/internal/dryrun"
)

const (
	defaultRequestsPerSecond  = 5.0
	slowRequestsPerSecond     = 2.0
	verySlowRequestsPerSecond = 1.0
	quotaWarningThreshold     = 0.5 // 50%
	quotaCriticalThreshold    = 0.2 // 20%
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
	token          string             // Static token (used if tokenSource is nil)
	tokenSource    oauth2.TokenSource // Auto-refreshing token source (preferred)
	asUserID       int64              // For admin masquerading
	rateLimiter    *AdaptiveRateLimiter
	retryPolicy    *RetryPolicy
	version        *CanvasVersion
	featureChecker *FeatureChecker
	logger         *slog.Logger
	quotaTotal     float64 // Detected or configured quota total
	cache          cache.CacheInterface
	cacheEnabled   bool
	userAgent      string // User-Agent header for API requests
	maxResults     int    // Max results for paginated requests (0 = unlimited)
	dryRun         bool   // Print curl commands instead of executing
	showToken      bool   // Show actual token in dry-run output (default: redacted)
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
	BaseURL        string
	Token          string             // Static token (used if TokenSource is nil)
	TokenSource    oauth2.TokenSource // Auto-refreshing token source (preferred over Token)
	RequestsPerSec float64
	Timeout        time.Duration
	Logger         *slog.Logger
	AsUserID       int64 // For admin masquerading (appends as_user_id param)
	Cache          cache.CacheInterface
	CacheEnabled   bool
	UserAgent      string // User-Agent header for API requests (required by Canvas)
	MaxResults     int    // Max results for paginated requests (0 = unlimited)
	DryRun         bool   // Print curl commands instead of executing requests
	ShowToken      bool   // Show actual token in dry-run output (default: redacted)
	// RetryInitialBackoff overrides the retry policy's initial backoff. Tests set
	// a tiny value (e.g. 1ms) so retryable-error cases don't spend ~7s sleeping
	// through the default 1s/2s/4s exponential backoff. Zero uses the default.
	RetryInitialBackoff time.Duration
}

// NewClient creates a new Canvas API client
func NewClient(config ClientConfig) (*Client, error) {
	if config.BaseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}

	// Require either Token or TokenSource
	if config.Token == "" && config.TokenSource == nil {
		return nil, fmt.Errorf("token or token source is required")
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

	// Set default User-Agent if not provided (required by Canvas API)
	if config.UserAgent == "" {
		config.UserAgent = "canvas-cli"
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

	retryPolicy := DefaultRetryPolicy()
	if config.RetryInitialBackoff > 0 {
		retryPolicy.InitialBackoff = config.RetryInitialBackoff
		if maxBackoff := config.RetryInitialBackoff * 8; maxBackoff < retryPolicy.MaxBackoff {
			retryPolicy.MaxBackoff = maxBackoff
		}
	}

	client := &Client{
		httpClient: &http.Client{
			Timeout:   config.Timeout,
			Transport: transport,
		},
		baseURL:      config.BaseURL,
		token:        config.Token,
		tokenSource:  config.TokenSource,
		asUserID:     config.AsUserID,
		rateLimiter:  NewAdaptiveRateLimiter(config.RequestsPerSec),
		retryPolicy:  retryPolicy,
		logger:       config.Logger,
		quotaTotal:   defaultQuotaTotal, // Will be updated from headers if available
		cache:        config.Cache,
		cacheEnabled: config.CacheEnabled && config.Cache != nil,
		userAgent:    config.UserAgent,
		maxResults:   config.MaxResults,
		dryRun:       config.DryRun,
		showToken:    config.ShowToken,
	}

	// Skip version detection in dry-run mode (no actual requests)
	if config.DryRun {
		client.version = &CanvasVersion{Major: 999, Minor: 999, Patch: 999, Raw: "dry-run"}
		client.featureChecker = NewFeatureChecker(client.version)
		return client, nil
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

// GetMaxResults returns the configured maximum results limit (0 = unlimited).
// This value is set during client construction and is immutable, so no lock is needed.
func (c *Client) GetMaxResults() int {
	return c.maxResults
}

// SupportsFeature checks if a feature is supported
func (c *Client) SupportsFeature(feature string) bool {
	return c.featureChecker.SupportsFeature(feature)
}

// getToken retrieves the current access token, refreshing if necessary
func (c *Client) getToken() (string, error) {
	// Prefer token source (supports auto-refresh)
	if c.tokenSource != nil {
		token, err := c.tokenSource.Token()
		if err != nil {
			return "", err
		}
		return token.AccessToken, nil
	}
	// Fall back to static token
	return c.token, nil
}

// doRequest performs an HTTP request with rate limiting and retry logic
func (c *Client) doRequest(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	// Get current token (may trigger refresh if using token source)
	token, err := c.getToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
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

	// Handle dry-run mode: print curl command and return mock response
	if c.dryRun {
		return c.handleDryRun(method, fullURL, token, body)
	}

	// Buffer the body once so each retry attempt gets a fresh reader.
	// Retrying with the same io.Reader would send an empty body on the second
	// attempt because the reader is already at EOF.
	var bodyBytes []byte
	if body != nil {
		var readErr error
		bodyBytes, readErr = io.ReadAll(body)
		if readErr != nil {
			return nil, fmt.Errorf("failed to read request body: %w", readErr)
		}
	}

	// Wait for rate limiter
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	// Execute with retry
	return c.retryPolicy.ExecuteWithRetry(ctx, func() (*http.Response, error) {
		var reqBody io.Reader
		if bodyBytes != nil {
			reqBody = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Set headers
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", c.userAgent)

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
			resp.Body.Close() // #nosec G104 -- closing body on error path; close error is irrelevant when returning an API error
			return resp, err
		}

		return resp, nil
	})
}

// handleDryRun prints the curl command and returns a mock response
func (c *Client) handleDryRun(method, fullURL, token string, body io.Reader) (*http.Response, error) {
	// Read body if present
	var bodyStr string
	if body != nil {
		bodyBytes, err := io.ReadAll(body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}
		bodyStr = string(bodyBytes)
	}

	// Build headers
	headers := []dryrun.Header{
		{Key: "Authorization", Value: "Bearer " + token},
		{Key: "Content-Type", Value: "application/json"},
		{Key: "Accept", Value: "application/json"},
		{Key: "User-Agent", Value: c.userAgent},
	}

	// Generate and print curl command
	curlCmd := dryrun.GenerateCurl(dryrun.CurlOptions{
		Method:    method,
		URL:       fullURL,
		Headers:   headers,
		Body:      bodyStr,
		ShowToken: c.showToken,
	})

	fmt.Println(curlCmd)

	// Return mock response with empty JSON array body
	// This works well for list operations (most common dry-run use case)
	// For single-object operations, the curl is still printed before any unmarshal errors
	mockBody := io.NopCloser(strings.NewReader("[]"))
	return &http.Response{
		StatusCode: 200,
		Body:       mockBody,
		Header:     make(http.Header),
	}, nil
}

// updateRateLimitFromHeaders updates the rate limiter based on response headers
func (c *Client) updateRateLimitFromHeaders(resp *http.Response) {
	// Parse X-Rate-Limit-Remaining header
	remaining := resp.Header.Get("X-Rate-Limit-Remaining")
	if remaining != "" {
		if remainingFloat, err := strconv.ParseFloat(remaining, 64); err == nil {
			c.mu.RLock()
			total := c.quotaTotal
			c.mu.RUnlock()
			c.rateLimiter.AdjustRate(remainingFloat, total)
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

// cacheKey generates a unique cache key for the given path
func (c *Client) cacheKey(path string) string {
	// Include base URL and masquerade user to ensure unique keys per instance/user
	key := c.baseURL + path
	if c.asUserID > 0 {
		key += fmt.Sprintf(":as_user:%d", c.asUserID)
	}

	// Hash the key to avoid issues with special characters — MD5 is used as a
	// fast hash to create a cache key, not for any cryptographic purpose.
	hash := md5.Sum([]byte(key)) // #nosec G401 -- non-cryptographic use: converts URL path to a cache key
	return hex.EncodeToString(hash[:])
}

// IsCacheEnabled returns whether caching is enabled
func (c *Client) IsCacheEnabled() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cacheEnabled
}

// SetCacheEnabled enables or disables caching
func (c *Client) SetCacheEnabled(enabled bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cacheEnabled = enabled && c.cache != nil
}

// ClearCache clears all cached responses
func (c *Client) ClearCache() {
	if c.cache != nil {
		c.cache.Clear()
	}
}

// CacheStats returns cache statistics
func (c *Client) CacheStats() cache.Stats {
	if c.cache != nil {
		return c.cache.Stats()
	}
	return cache.Stats{}
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

// Delete performs a DELETE request.
// It drains and closes the response body before returning, which keeps the
// underlying connection eligible for reuse. Callers that need to decode the
// response body should use DeleteJSON instead.
func (c *Client) Delete(ctx context.Context, path string) (*http.Response, error) {
	resp, err := c.doRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return resp, err
	}
	// Drain and close so the connection can be reused; both errors are
	// irrelevant to the caller (the DELETE itself already succeeded).
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
	return resp, nil
}

// unmarshalTolerant decodes body into result, tolerating Canvas endpoints that
// occasionally return a single object wrapped in an array (e.g. some deployments
// answer an enrollment create with [] or [obj] instead of obj). When the body is
// a JSON array and result targets a single (non-slice) value, the first element
// is used; an empty array leaves result zero-valued and is treated as success.
// Any other decode error is returned unchanged, and slice targets keep normal
// semantics.
func unmarshalTolerant(body []byte, result interface{}) error {
	err := json.Unmarshal(body, result)
	if err == nil {
		return nil
	}

	// Only attempt the array fallback when the body actually is a JSON array.
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 || trimmed[0] != '[' {
		return err
	}

	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return err
	}
	elemType := rv.Elem().Type()
	if elemType.Kind() == reflect.Slice {
		// Target already expects a slice; the original error stands.
		return err
	}

	slicePtr := reflect.New(reflect.SliceOf(elemType))
	if e := json.Unmarshal(body, slicePtr.Interface()); e != nil {
		// Not an array of the target type — surface the original error.
		return err
	}
	slice := slicePtr.Elem()
	if slice.Len() > 0 {
		rv.Elem().Set(slice.Index(0))
	} else {
		// Empty array: reset result to its zero value and treat as success.
		rv.Elem().Set(reflect.Zero(elemType))
	}
	return nil
}

// DeleteJSON performs a DELETE request and decodes the JSON response body into
// result. If result is nil the body is discarded. This is the preferred method
// when the Canvas DELETE endpoint returns a meaningful JSON body.
func (c *Client) DeleteJSON(ctx context.Context, path string, result interface{}) error {
	resp, err := c.doRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read delete response body: %w", err)
	}

	if result != nil {
		return unmarshalTolerant(bodyBytes, result)
	}
	return nil
}

// GetJSON performs a GET request and decodes JSON response
// If caching is enabled, cached responses will be returned when available
func (c *Client) GetJSON(ctx context.Context, path string, result interface{}) error {
	// Check cache first if enabled
	if c.cacheEnabled && c.cache != nil {
		key := c.cacheKey(path)
		if err := c.cache.GetJSON(key, result); err == nil {
			return nil // Cache hit
		}
	}

	resp, err := c.Get(ctx, path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the response body so we can cache it
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for Canvas soft errors (HTTP 200 with error body)
	if err := checkSoftError(body); err != nil {
		return err
	}

	// Decode the response
	if err := unmarshalTolerant(body, result); err != nil {
		return err
	}

	// Cache the response if caching is enabled
	if c.cacheEnabled && c.cache != nil {
		key := c.cacheKey(path)
		c.cache.Set(key, body)
	}

	return nil
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

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := checkSoftError(bodyBytes); err != nil {
		return err
	}

	if result != nil {
		return unmarshalTolerant(bodyBytes, result)
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

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := checkSoftError(bodyBytes); err != nil {
		return err
	}

	if result != nil {
		return unmarshalTolerant(bodyBytes, result)
	}

	return nil
}

// PatchJSON performs a PATCH request with JSON body and decodes JSON response.
// Canvas uses PATCH for partial updates on some endpoints (e.g. grading_periods/batch_update).
func (c *Client) PatchJSON(ctx context.Context, path string, body interface{}, result interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPatch, path, bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := checkSoftError(bodyBytes); err != nil {
		return err
	}

	if result != nil {
		return unmarshalTolerant(bodyBytes, result)
	}

	return nil
}

// maxPaginationPages is an upper bound on the number of pages fetched in a
// single GetAllPages / GetAllPagesGeneric call. A misbehaving server that keeps
// returning the same Link: next header would otherwise loop forever.
const maxPaginationPages = 10_000

// GetAllPagesGeneric fetches all pages of a paginated endpoint using generics
// This is the preferred method for type-safe pagination with better performance
// If caching is enabled, cached responses will be returned when available
// If maxResults is set, stops fetching when limit is reached
func GetAllPagesGeneric[T any](c *Client, ctx context.Context, path string) ([]T, error) {
	// Check cache first if enabled (only if no limit set, as cached results might exceed limit)
	if c.cacheEnabled && c.cache != nil && c.maxResults == 0 {
		key := c.cacheKey("pages:" + path)
		var cached []T
		if err := c.cache.GetJSON(key, &cached); err == nil {
			return cached, nil // Cache hit
		}
	}

	var allResults []T
	currentURL := path
	pageCount := 0

	for currentURL != "" {
		pageCount++
		if pageCount > maxPaginationPages {
			return nil, fmt.Errorf("pagination aborted: exceeded %d pages (possible infinite loop from repeated Link header)", maxPaginationPages)
		}

		resp, err := c.Get(ctx, currentURL)
		if err != nil {
			return nil, err
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close() // #nosec G104 -- body already fully read; close error is irrelevant
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		if err := checkSoftError(bodyBytes); err != nil {
			return nil, err
		}

		var pageResults []T
		if err := json.Unmarshal(bodyBytes, &pageResults); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		allResults = append(allResults, pageResults...)

		// Check if we've reached the limit
		if c.maxResults > 0 && len(allResults) >= c.maxResults {
			// Truncate to exact limit
			allResults = allResults[:c.maxResults]
			break
		}

		// Check for next page
		links := ParsePaginationLinks(resp)
		if links.HasNextPage() {
			// Extract path from full URL
			nextURL, err := url.Parse(links.Next)
			if err != nil {
				return nil, fmt.Errorf("failed to parse next URL: %w", err)
			}
			// Detect same-URL cycle before following the link.
			var next string
			if nextURL.RawQuery != "" {
				next = nextURL.Path + "?" + nextURL.RawQuery
			} else {
				next = nextURL.Path
			}
			if next == currentURL {
				return nil, fmt.Errorf("pagination aborted: server returned identical next page URL %q", next)
			}
			currentURL = next
		} else {
			currentURL = ""
		}
	}

	// Cache the combined result if caching is enabled
	if c.cacheEnabled && c.cache != nil {
		key := c.cacheKey("pages:" + path)
		// Marshal for caching
		allJSON, err := json.Marshal(allResults)
		if err == nil {
			c.cache.Set(key, allJSON)
		}
	}

	return allResults, nil
}

// GetAllPages fetches all pages of a paginated endpoint
// Deprecated: Use GetAllPagesGeneric for better performance and type safety
// This method will be removed in v2.0.0
// If caching is enabled, cached responses will be returned when available
// If maxResults is set, stops fetching when limit is reached
func (c *Client) GetAllPages(ctx context.Context, path string, result interface{}) error {
	// Check cache first if enabled (only if no limit set, as cached results might exceed limit)
	if c.cacheEnabled && c.cache != nil && c.maxResults == 0 {
		key := c.cacheKey("pages:" + path)
		if err := c.cache.GetJSON(key, result); err == nil {
			return nil // Cache hit
		}
	}

	var allResults []json.RawMessage
	currentURL := path
	pageCount := 0

	for currentURL != "" {
		pageCount++
		if pageCount > maxPaginationPages {
			return fmt.Errorf("pagination aborted: exceeded %d pages (possible infinite loop from repeated Link header)", maxPaginationPages)
		}

		resp, err := c.Get(ctx, currentURL)
		if err != nil {
			return err
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close() // #nosec G104 -- body already fully read; close error is irrelevant
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		if err := checkSoftError(bodyBytes); err != nil {
			return err
		}

		var pageResults []json.RawMessage
		if err := json.Unmarshal(bodyBytes, &pageResults); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}

		allResults = append(allResults, pageResults...)

		// Check if we've reached the limit
		if c.maxResults > 0 && len(allResults) >= c.maxResults {
			// Truncate to exact limit
			allResults = allResults[:c.maxResults]
			break
		}

		// Check for next page
		links := ParsePaginationLinks(resp)
		if links.HasNextPage() {
			// Extract path from full URL
			nextURL, err := url.Parse(links.Next)
			if err != nil {
				return fmt.Errorf("failed to parse next URL: %w", err)
			}
			var next string
			if nextURL.RawQuery != "" {
				next = nextURL.Path + "?" + nextURL.RawQuery
			} else {
				next = nextURL.Path
			}
			if next == currentURL {
				return fmt.Errorf("pagination aborted: server returned identical next page URL %q", next)
			}
			currentURL = next
		} else {
			currentURL = ""
		}
	}

	// Use reflection to directly append results to the slice
	// This avoids unnecessary marshal/unmarshal round-trip
	resultValue := reflect.ValueOf(result)
	if resultValue.Kind() != reflect.Pointer || resultValue.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("result must be a pointer to a slice")
	}

	sliceValue := resultValue.Elem()
	elemType := sliceValue.Type().Elem()

	for i, raw := range allResults {
		// Create a new element of the slice's element type
		elem := reflect.New(elemType)
		if err := json.Unmarshal(raw, elem.Interface()); err != nil {
			// Degrade gracefully: one malformed element shouldn't discard the
			// whole page. Warn and skip it so the command stays useful.
			c.logger.Warn("skipping list element that failed to decode",
				"index", i, "error", err)
			continue
		}
		sliceValue = reflect.Append(sliceValue, elem.Elem())
	}

	// Set the slice back to the result pointer
	resultValue.Elem().Set(sliceValue)

	// Cache the combined result if caching is enabled
	if c.cacheEnabled && c.cache != nil {
		key := c.cacheKey("pages:" + path)
		// Marshal once for caching only
		allJSON, err := json.Marshal(allResults)
		if err == nil {
			c.cache.Set(key, allJSON)
		}
	}

	return nil
}
