package api

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestParseAPIError_JSONError(t *testing.T) {
	jsonResponse := `{
		"errors": [
			{"message": "Invalid access token"}
		]
	}`

	resp := &http.Response{
		StatusCode: http.StatusUnauthorized,
		Body:       io.NopCloser(strings.NewReader(jsonResponse)),
	}

	err := ParseAPIError(resp)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}

	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, apiErr.StatusCode)
	}

	if len(apiErr.Errors) == 0 {
		t.Fatal("expected errors to be populated")
	}

	if apiErr.Errors[0].Message != "Invalid access token" {
		t.Errorf("expected error message 'Invalid access token', got '%s'", apiErr.Errors[0].Message)
	}

	// Check suggestion for 401
	if apiErr.Suggestion == "" {
		t.Error("expected suggestion to be set for 401 error")
	}
	if apiErr.DocsURL == "" {
		t.Error("expected docs URL to be set for 401 error")
	}
}

func TestParseAPIError_NonJSONError(t *testing.T) {
	textResponse := "Internal Server Error"

	resp := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(strings.NewReader(textResponse)),
	}

	err := ParseAPIError(resp)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}

	if apiErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, apiErr.StatusCode)
	}

	if len(apiErr.Errors) == 0 {
		t.Fatal("expected errors to be populated")
	}

	if apiErr.Errors[0].Message != textResponse {
		t.Errorf("expected error message '%s', got '%s'", textResponse, apiErr.Errors[0].Message)
	}

	// Check suggestion for 500
	if apiErr.Suggestion == "" {
		t.Error("expected suggestion to be set for 500 error")
	}
}

func TestParseAPIError_Forbidden(t *testing.T) {
	jsonResponse := `{"errors": [{"message": "Access denied"}]}`

	resp := &http.Response{
		StatusCode: http.StatusForbidden,
		Body:       io.NopCloser(strings.NewReader(jsonResponse)),
	}

	err := ParseAPIError(resp)
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}

	if apiErr.Suggestion == "" {
		t.Error("expected suggestion to be set for 403 error")
	}
}

func TestParseAPIError_NotFound(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(strings.NewReader(`{"errors": [{"message": "Not found"}]}`)),
	}

	err := ParseAPIError(resp)
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}

	if apiErr.Suggestion == "" {
		t.Error("expected suggestion to be set for 404 error")
	}
}

func TestParseAPIError_UnprocessableEntity(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusUnprocessableEntity,
		Body:       io.NopCloser(strings.NewReader(`{"errors": [{"message": "Validation failed"}]}`)),
	}

	err := ParseAPIError(resp)
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}

	if apiErr.Suggestion == "" {
		t.Error("expected suggestion to be set for 422 error")
	}
}

func TestParseAPIError_RateLimit(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Body:       io.NopCloser(strings.NewReader(`{"errors": [{"message": "Rate limit exceeded"}]}`)),
	}

	err := ParseAPIError(resp)
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}

	if apiErr.Suggestion == "" {
		t.Error("expected suggestion to be set for 429 error")
	}
	if apiErr.DocsURL == "" {
		t.Error("expected docs URL to be set for 429 error")
	}
}

func TestParseAPIError_BadGateway(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusBadGateway,
		Body:       io.NopCloser(strings.NewReader("Bad Gateway")),
	}

	err := ParseAPIError(resp)
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}

	if apiErr.Suggestion == "" {
		t.Error("expected suggestion to be set for 502 error")
	}
}

func TestParseAPIError_ServiceUnavailable(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusServiceUnavailable,
		Body:       io.NopCloser(strings.NewReader("Service Unavailable")),
	}

	err := ParseAPIError(resp)
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}

	if apiErr.Suggestion == "" {
		t.Error("expected suggestion to be set for 503 error")
	}
}

func TestIsRateLimitError_True(t *testing.T) {
	apiErr := &APIError{
		StatusCode: http.StatusTooManyRequests,
	}

	if !IsRateLimitError(apiErr) {
		t.Error("expected IsRateLimitError to return true for 429 error")
	}
}

func TestIsRateLimitError_False(t *testing.T) {
	apiErr := &APIError{
		StatusCode: http.StatusInternalServerError,
	}

	if IsRateLimitError(apiErr) {
		t.Error("expected IsRateLimitError to return false for non-429 error")
	}
}

func TestIsRateLimitError_NotAPIError(t *testing.T) {
	err := io.EOF

	if IsRateLimitError(err) {
		t.Error("expected IsRateLimitError to return false for non-APIError")
	}
}

func TestIsAuthError_True(t *testing.T) {
	apiErr := &APIError{
		StatusCode: http.StatusUnauthorized,
	}

	if !IsAuthError(apiErr) {
		t.Error("expected IsAuthError to return true for 401 error")
	}
}

func TestIsAuthError_False(t *testing.T) {
	apiErr := &APIError{
		StatusCode: http.StatusForbidden,
	}

	if IsAuthError(apiErr) {
		t.Error("expected IsAuthError to return false for non-401 error")
	}
}

func TestIsAuthError_NotAPIError(t *testing.T) {
	err := io.EOF

	if IsAuthError(err) {
		t.Error("expected IsAuthError to return false for non-APIError")
	}
}

func TestIsNotFoundError_True(t *testing.T) {
	apiErr := &APIError{
		StatusCode: http.StatusNotFound,
	}

	if !IsNotFoundError(apiErr) {
		t.Error("expected IsNotFoundError to return true for 404 error")
	}
}

func TestIsNotFoundError_False(t *testing.T) {
	apiErr := &APIError{
		StatusCode: http.StatusInternalServerError,
	}

	if IsNotFoundError(apiErr) {
		t.Error("expected IsNotFoundError to return false for non-404 error")
	}
}

func TestIsNotFoundError_NotAPIError(t *testing.T) {
	err := io.EOF

	if IsNotFoundError(err) {
		t.Error("expected IsNotFoundError to return false for non-APIError")
	}
}

func TestAPIError_Error(t *testing.T) {
	apiErr := &APIError{
		StatusCode: http.StatusNotFound,
		Errors: []ErrorDetail{
			{Message: "Resource not found"},
		},
	}

	errMsg := apiErr.Error()
	if errMsg == "" {
		t.Error("expected non-empty error message")
	}

	// Error() returns the first error message
	if !strings.Contains(errMsg, "Resource not found") {
		t.Errorf("expected error message to contain 'Resource not found', got: %s", errMsg)
	}
}
