package api

import (
	"strings"
	"testing"
)

func TestAPIError_Error_AllStatusCodes(t *testing.T) {
	cases := []struct {
		code    int
		wantStr string
	}{
		{409, "Conflict"},
		{422, "Unprocessable entity"},
		{429, "Rate limit exceeded"},
		{500, "Internal server error"},
		{502, "Bad gateway"},
		{503, "Service unavailable"},
		// default: client error range
		{410, "Client error (HTTP 410)"},
		// default: server error range
		{504, "Server error (HTTP 504)"},
		// default: completely unknown
		{200, "Unknown API error"},
	}

	for _, tc := range cases {
		apiErr := &APIError{StatusCode: tc.code}
		msg := apiErr.Error()
		if msg == "" {
			t.Errorf("status %d: expected non-empty message", tc.code)
			continue
		}
		if !strings.Contains(msg, tc.wantStr) {
			t.Errorf("status %d: expected %q in message, got: %q", tc.code, tc.wantStr, msg)
		}
	}
}

func TestAPIError_Error_WithSuggestionAndDocs(t *testing.T) {
	apiErr := &APIError{
		StatusCode: 403,
		Suggestion: "Check your permissions",
		DocsURL:    "https://docs.example.com/api",
	}
	msg := apiErr.Error()
	if !strings.Contains(msg, "Suggestion: Check your permissions") {
		t.Errorf("expected suggestion in message, got: %q", msg)
	}
	if !strings.Contains(msg, "Docs: https://docs.example.com/api") {
		t.Errorf("expected docs URL in message, got: %q", msg)
	}
}

func TestAPIError_Error_WithErrorsAndSuggestion(t *testing.T) {
	apiErr := &APIError{
		StatusCode: 422,
		Errors:     []ErrorDetail{{Message: "field is required"}},
		Suggestion: "Provide all required fields",
	}
	msg := apiErr.Error()
	if !strings.Contains(msg, "field is required") {
		t.Errorf("expected 'field is required' in message, got: %q", msg)
	}
	if !strings.Contains(msg, "Suggestion: Provide all required fields") {
		t.Errorf("expected suggestion appended to message, got: %q", msg)
	}
}
