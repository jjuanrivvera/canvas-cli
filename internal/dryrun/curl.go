package dryrun

import (
	"fmt"
	"strings"
)

// Header represents an HTTP header key-value pair
type Header struct {
	Key   string
	Value string
}

// CurlOptions holds the options for generating a curl command
type CurlOptions struct {
	Method    string
	URL       string
	Headers   []Header
	Body      string
	ShowToken bool
}

// GenerateCurl generates a curl command string from the given options
func GenerateCurl(opts CurlOptions) string {
	var parts []string

	// Start with curl command and method
	parts = append(parts, fmt.Sprintf("curl -X %s '%s'", opts.Method, escapeURL(opts.URL)))

	// Add headers
	for _, h := range opts.Headers {
		value := h.Value
		// Redact Authorization header if not showing token
		if h.Key == "Authorization" && !opts.ShowToken {
			value = "Bearer [REDACTED]"
		}
		parts = append(parts, fmt.Sprintf("-H '%s: %s'", h.Key, escapeSingleQuotes(value)))
	}

	// Add body if present
	if opts.Body != "" {
		parts = append(parts, fmt.Sprintf("-d '%s'", escapeSingleQuotes(opts.Body)))
	}

	// Join with line continuation for readability
	return strings.Join(parts, " \\\n  ")
}

// escapeURL escapes single quotes in URLs for shell safety
func escapeURL(url string) string {
	return escapeSingleQuotes(url)
}

// escapeSingleQuotes escapes single quotes for shell safety
// Replaces ' with '\” (end quote, escaped quote, start quote)
func escapeSingleQuotes(s string) string {
	return strings.ReplaceAll(s, "'", "'\\''")
}
