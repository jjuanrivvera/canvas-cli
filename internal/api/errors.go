package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// ParseAPIError parses an error response from the Canvas API
func ParseAPIError(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read error response: %w", err)
	}

	var apiErr APIError
	apiErr.StatusCode = resp.StatusCode

	// Try to parse as JSON error
	if err := json.Unmarshal(body, &apiErr); err != nil {
		// If parsing fails, use the raw body as error message
		apiErr.Errors = []ErrorDetail{{Message: string(body)}}
	}

	// Add suggestions based on status code
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		apiErr.Suggestion = "Your authentication token may be expired or invalid. Try running 'canvas auth login' again."
		apiErr.DocsURL = "https://canvas.instructure.com/doc/api/file.oauth.html"
	case http.StatusForbidden:
		apiErr.Suggestion = "You don't have permission to access this resource. Check your Canvas role and permissions."
	case http.StatusNotFound:
		apiErr.Suggestion = "The requested resource was not found. Verify the ID and try again."
	case http.StatusUnprocessableEntity:
		apiErr.Suggestion = "The request was invalid. Check the required fields and data format."
	case http.StatusTooManyRequests:
		apiErr.Suggestion = "Rate limit exceeded. The CLI will automatically slow down requests. Please wait a moment."
		apiErr.DocsURL = "https://canvas.instructure.com/doc/api/file.throttling.html"
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
		apiErr.Suggestion = "Canvas is experiencing issues. Please try again in a few moments."
	}

	return &apiErr
}

// IsRateLimitError checks if the error is a rate limit error
// Uses errors.As to properly handle wrapped errors
func IsRateLimitError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusTooManyRequests
	}
	return false
}

// IsAuthError checks if the error is an authentication error
// Uses errors.As to properly handle wrapped errors
func IsAuthError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusUnauthorized
	}
	return false
}

// IsNotFoundError checks if the error is a not found error
// Uses errors.As to properly handle wrapped errors
func IsNotFoundError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

// IsForbiddenError checks if the error is a forbidden error
// Uses errors.As to properly handle wrapped errors
func IsForbiddenError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusForbidden
	}
	return false
}

// IsServerError checks if the error is a server error (5xx)
// Uses errors.As to properly handle wrapped errors
func IsServerError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode >= 500 && apiErr.StatusCode < 600
	}
	return false
}
