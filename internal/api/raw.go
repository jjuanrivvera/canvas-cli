package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// RawService provides raw access to any Canvas API endpoint
type RawService struct {
	client *Client
}

// NewRawService creates a new RawService
func NewRawService(client *Client) *RawService {
	return &RawService{client: client}
}

// RawRequestOptions contains options for a raw API request
type RawRequestOptions struct {
	Body     interface{}         // Request body (will be JSON-encoded if not nil)
	Query    map[string][]string // Query parameters
	Headers  map[string]string   // Custom headers
	Paginate bool                // Follow pagination links (GET only)
}

// RawResponse contains the response from a raw API request
type RawResponse struct {
	StatusCode int              `json:"status_code"`
	Headers    http.Header      `json:"headers,omitempty"`
	Body       json.RawMessage  `json:"body"`
	Pagination *PaginationLinks `json:"pagination,omitempty"`
}

// Request makes a raw API request
func (s *RawService) Request(ctx context.Context, method, path string, opts *RawRequestOptions) (*RawResponse, error) {
	if opts == nil {
		opts = &RawRequestOptions{}
	}

	// Normalize method to uppercase
	method = strings.ToUpper(method)

	// Validate method
	switch method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodHead:
		// Valid
	default:
		return nil, fmt.Errorf("unsupported HTTP method: %s", method)
	}

	// Build URL with query parameters
	fullPath := path
	if len(opts.Query) > 0 {
		params := url.Values(opts.Query)
		if strings.Contains(path, "?") {
			fullPath = path + "&" + params.Encode()
		} else {
			fullPath = path + "?" + params.Encode()
		}
	}

	// Handle pagination for GET requests
	if opts.Paginate && method == http.MethodGet {
		return s.requestWithPagination(ctx, fullPath, opts)
	}

	return s.doRequest(ctx, method, fullPath, opts)
}

// doRequest performs a single request
func (s *RawService) doRequest(ctx context.Context, method, path string, opts *RawRequestOptions) (*RawResponse, error) {
	var bodyReader io.Reader
	if opts.Body != nil {
		jsonBody, err := json.Marshal(opts.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	resp, err := s.client.doRequest(ctx, method, path, bodyReader)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	result := &RawResponse{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
	}

	// Only set body if it's valid JSON or non-empty
	if len(body) > 0 {
		// Try to compact JSON if it's valid
		var compacted bytes.Buffer
		if err := json.Compact(&compacted, body); err == nil {
			result.Body = compacted.Bytes()
		} else {
			// Not valid JSON, wrap as a string
			result.Body, _ = json.Marshal(string(body))
		}
	}

	// Parse pagination links
	result.Pagination = ParsePaginationLinks(resp)

	return result, nil
}

// requestWithPagination follows pagination links and returns all results
func (s *RawService) requestWithPagination(ctx context.Context, path string, opts *RawRequestOptions) (*RawResponse, error) {
	var allResults []json.RawMessage
	currentPath := path
	var lastResponse *RawResponse

	for currentPath != "" {
		resp, err := s.doRequest(ctx, http.MethodGet, currentPath, &RawRequestOptions{
			Headers: opts.Headers,
			// Don't pass body or query for paginated requests after the first one
		})
		if err != nil {
			return nil, err
		}
		lastResponse = resp

		// Try to unmarshal as an array
		var pageResults []json.RawMessage
		if err := json.Unmarshal(resp.Body, &pageResults); err != nil {
			// Not an array, just return the single response
			return resp, nil
		}

		allResults = append(allResults, pageResults...)

		// Check for next page
		if resp.Pagination != nil && resp.Pagination.HasNextPage() {
			// Extract path from full URL
			nextURL, err := url.Parse(resp.Pagination.Next)
			if err != nil {
				return nil, fmt.Errorf("failed to parse next URL: %w", err)
			}
			if nextURL.RawQuery != "" {
				currentPath = nextURL.Path + "?" + nextURL.RawQuery
			} else {
				currentPath = nextURL.Path
			}
		} else {
			currentPath = ""
		}
	}

	// Combine all results
	combinedBody, err := json.Marshal(allResults)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal combined results: %w", err)
	}

	result := &RawResponse{
		StatusCode: lastResponse.StatusCode,
		Headers:    lastResponse.Headers,
		Body:       combinedBody,
	}

	return result, nil
}

// GetRequestStats returns statistics about requests made
func (s *RawService) GetRequestStats() map[string]interface{} {
	return map[string]interface{}{
		"rate_limit": s.client.rateLimiter.GetCurrentRate(),
		"base_url":   s.client.baseURL,
	}
}
