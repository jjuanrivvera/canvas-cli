package dryrun

import (
	"strings"
	"testing"
)

func TestGenerateCurl_BasicGET(t *testing.T) {
	opts := CurlOptions{
		Method: "GET",
		URL:    "https://canvas.example.com/api/v1/courses",
		Headers: []Header{
			{Key: "Authorization", Value: "Bearer test-token"},
			{Key: "Content-Type", Value: "application/json"},
			{Key: "Accept", Value: "application/json"},
			{Key: "User-Agent", Value: "canvas-cli/1.5.2"},
		},
		ShowToken: false,
	}

	result := GenerateCurl(opts)

	// Check method and URL
	if !strings.Contains(result, "curl -X GET") {
		t.Errorf("Expected curl command to contain 'curl -X GET', got: %s", result)
	}
	if !strings.Contains(result, "https://canvas.example.com/api/v1/courses") {
		t.Errorf("Expected URL in curl command, got: %s", result)
	}

	// Check token is redacted
	if !strings.Contains(result, "[REDACTED]") {
		t.Errorf("Expected token to be redacted, got: %s", result)
	}
	if strings.Contains(result, "test-token") {
		t.Errorf("Expected token to be hidden, but found 'test-token' in: %s", result)
	}

	// Check headers present
	if !strings.Contains(result, "Content-Type: application/json") {
		t.Errorf("Expected Content-Type header, got: %s", result)
	}
	if !strings.Contains(result, "User-Agent: canvas-cli/1.5.2") {
		t.Errorf("Expected User-Agent header, got: %s", result)
	}
}

func TestGenerateCurl_ShowToken(t *testing.T) {
	opts := CurlOptions{
		Method: "GET",
		URL:    "https://canvas.example.com/api/v1/courses",
		Headers: []Header{
			{Key: "Authorization", Value: "Bearer actual-secret-token"},
		},
		ShowToken: true,
	}

	result := GenerateCurl(opts)

	// Check token is shown when ShowToken is true
	if !strings.Contains(result, "Bearer actual-secret-token") {
		t.Errorf("Expected token to be shown when ShowToken=true, got: %s", result)
	}
	if strings.Contains(result, "[REDACTED]") {
		t.Errorf("Expected no redaction when ShowToken=true, got: %s", result)
	}
}

func TestGenerateCurl_POSTWithBody(t *testing.T) {
	opts := CurlOptions{
		Method: "POST",
		URL:    "https://canvas.example.com/api/v1/courses/123/assignments",
		Headers: []Header{
			{Key: "Authorization", Value: "Bearer test-token"},
			{Key: "Content-Type", Value: "application/json"},
		},
		Body:      `{"assignment":{"name":"Test Assignment"}}`,
		ShowToken: false,
	}

	result := GenerateCurl(opts)

	// Check method
	if !strings.Contains(result, "curl -X POST") {
		t.Errorf("Expected POST method, got: %s", result)
	}

	// Check body is present
	if !strings.Contains(result, "-d '") {
		t.Errorf("Expected -d flag for body, got: %s", result)
	}
	if !strings.Contains(result, `"assignment"`) {
		t.Errorf("Expected body content, got: %s", result)
	}
}

func TestGenerateCurl_PUTMethod(t *testing.T) {
	opts := CurlOptions{
		Method: "PUT",
		URL:    "https://canvas.example.com/api/v1/courses/123",
		Headers: []Header{
			{Key: "Authorization", Value: "Bearer test-token"},
		},
		Body:      `{"course":{"name":"Updated"}}`,
		ShowToken: false,
	}

	result := GenerateCurl(opts)

	if !strings.Contains(result, "curl -X PUT") {
		t.Errorf("Expected PUT method, got: %s", result)
	}
}

func TestGenerateCurl_DELETEMethod(t *testing.T) {
	opts := CurlOptions{
		Method: "DELETE",
		URL:    "https://canvas.example.com/api/v1/courses/123",
		Headers: []Header{
			{Key: "Authorization", Value: "Bearer test-token"},
		},
		ShowToken: false,
	}

	result := GenerateCurl(opts)

	if !strings.Contains(result, "curl -X DELETE") {
		t.Errorf("Expected DELETE method, got: %s", result)
	}

	// DELETE should not have body
	if strings.Contains(result, "-d '") {
		t.Errorf("DELETE should not have body, got: %s", result)
	}
}

func TestGenerateCurl_EmptyBody(t *testing.T) {
	opts := CurlOptions{
		Method: "POST",
		URL:    "https://canvas.example.com/api/v1/courses",
		Headers: []Header{
			{Key: "Authorization", Value: "Bearer test-token"},
		},
		Body:      "",
		ShowToken: false,
	}

	result := GenerateCurl(opts)

	// Empty body should not include -d flag
	if strings.Contains(result, "-d ''") {
		t.Errorf("Empty body should not include -d flag, got: %s", result)
	}
}

func TestGenerateCurl_URLWithQueryParams(t *testing.T) {
	opts := CurlOptions{
		Method: "GET",
		URL:    "https://canvas.example.com/api/v1/courses?as_user_id=999&include[]=term",
		Headers: []Header{
			{Key: "Authorization", Value: "Bearer test-token"},
		},
		ShowToken: false,
	}

	result := GenerateCurl(opts)

	// Check URL with query params is preserved
	if !strings.Contains(result, "as_user_id=999") {
		t.Errorf("Expected query parameter preserved, got: %s", result)
	}
	if !strings.Contains(result, "include[]=term") {
		t.Errorf("Expected include parameter preserved, got: %s", result)
	}
}

func TestEscapeSingleQuotes(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"it's a test", "it'\\''s a test"},
		{"'quoted'", "'\\''quoted'\\''"},
		{"no quotes here", "no quotes here"},
	}

	for _, test := range tests {
		result := escapeSingleQuotes(test.input)
		if result != test.expected {
			t.Errorf("escapeSingleQuotes(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestGenerateCurl_BodyWithSingleQuotes(t *testing.T) {
	opts := CurlOptions{
		Method: "POST",
		URL:    "https://canvas.example.com/api/v1/courses",
		Headers: []Header{
			{Key: "Authorization", Value: "Bearer test-token"},
		},
		Body:      `{"name":"It's a test"}`,
		ShowToken: false,
	}

	result := GenerateCurl(opts)

	// Single quote in body should be escaped
	if !strings.Contains(result, "'\\''") {
		t.Errorf("Expected escaped single quote in body, got: %s", result)
	}
}

func TestGenerateCurl_MultilineFormat(t *testing.T) {
	opts := CurlOptions{
		Method: "GET",
		URL:    "https://canvas.example.com/api/v1/courses",
		Headers: []Header{
			{Key: "Authorization", Value: "Bearer test-token"},
			{Key: "Content-Type", Value: "application/json"},
		},
		ShowToken: false,
	}

	result := GenerateCurl(opts)

	// Check for line continuation characters
	if !strings.Contains(result, "\\\n") {
		t.Errorf("Expected multi-line format with continuation, got: %s", result)
	}
}
