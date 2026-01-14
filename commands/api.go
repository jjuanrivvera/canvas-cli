package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var apiCmd = &cobra.Command{
	Use:   "api <METHOD> <PATH>",
	Short: "Make raw API requests to Canvas",
	Long: `Make raw API requests to any Canvas API endpoint.

This command provides direct access to the Canvas API for advanced use cases
or endpoints not yet supported by dedicated commands.

Methods: GET, POST, PUT, DELETE, PATCH, HEAD

Examples:
  # List all courses
  canvas api GET /api/v1/courses

  # Create a course (with JSON body)
  canvas api POST /api/v1/accounts/1/courses -d '{"course":{"name":"Test Course"}}'

  # Search users with query parameters
  canvas api GET /api/v1/users -q "search_term=john" -q "per_page=50"

  # Update an assignment
  canvas api PUT /api/v1/courses/123/assignments/456 -d '{"assignment":{"name":"Updated"}}'

  # Delete an assignment
  canvas api DELETE /api/v1/courses/123/assignments/456

  # Get all pages of a paginated endpoint
  canvas api GET /api/v1/courses --paginate

  # Read body from file
  canvas api POST /api/v1/accounts/1/courses --data-file course.json`,
	Args: ExactArgsWithUsage(2, "method", "path"),
	RunE: runAPICommand,
}

var (
	apiData        string
	apiDataFile    string
	apiQuery       []string
	apiHeaders     []string
	apiPaginate    bool
	apiRawOutput   bool
	apiShowHeaders bool
)

func init() {
	rootCmd.AddCommand(apiCmd)

	apiCmd.Flags().StringVarP(&apiData, "data", "d", "", "JSON data for request body")
	apiCmd.Flags().StringVar(&apiDataFile, "data-file", "", "Read JSON data from file")
	apiCmd.Flags().StringArrayVarP(&apiQuery, "query", "q", nil, "Query parameters (key=value, repeatable)")
	apiCmd.Flags().StringArrayVarP(&apiHeaders, "header", "H", nil, "Custom headers (key:value, repeatable)")
	apiCmd.Flags().BoolVar(&apiPaginate, "paginate", false, "Follow pagination links (GET only)")
	apiCmd.Flags().BoolVar(&apiRawOutput, "raw", false, "Output raw response without formatting")
	apiCmd.Flags().BoolVar(&apiShowHeaders, "show-headers", false, "Include response headers in output")
}

func runAPICommand(cmd *cobra.Command, args []string) error {
	method := strings.ToUpper(args[0])
	path := args[1]

	// Validate method
	switch method {
	case "GET", "POST", "PUT", "DELETE", "PATCH", "HEAD":
		// Valid
	default:
		return fmt.Errorf("unsupported HTTP method: %s (use GET, POST, PUT, DELETE, PATCH, or HEAD)", method)
	}

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}
	service := api.NewRawService(client)

	// Build request options
	opts := &api.RawRequestOptions{
		Paginate: apiPaginate && method == "GET",
	}

	// Parse body from --data or --data-file
	if apiData != "" && apiDataFile != "" {
		return fmt.Errorf("cannot use both --data and --data-file")
	}

	if apiData != "" {
		var body interface{}
		if err := json.Unmarshal([]byte(apiData), &body); err != nil {
			return fmt.Errorf("invalid JSON in --data: %w", err)
		}
		opts.Body = body
	}

	if apiDataFile != "" {
		var reader io.Reader
		if apiDataFile == "-" {
			reader = cmd.InOrStdin()
		} else {
			file, err := os.Open(apiDataFile)
			if err != nil {
				return fmt.Errorf("failed to open data file: %w", err)
			}
			defer file.Close()
			reader = file
		}

		data, err := io.ReadAll(reader)
		if err != nil {
			return fmt.Errorf("failed to read data file: %w", err)
		}

		var body interface{}
		if err := json.Unmarshal(data, &body); err != nil {
			return fmt.Errorf("invalid JSON in data file: %w", err)
		}
		opts.Body = body
	}

	// Parse query parameters
	if len(apiQuery) > 0 {
		query := make(map[string][]string)
		for _, q := range apiQuery {
			parts := strings.SplitN(q, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid query parameter format: %s (use key=value)", q)
			}
			key := parts[0]
			value := parts[1]
			query[key] = append(query[key], value)
		}
		opts.Query = query
	}

	// Parse custom headers
	if len(apiHeaders) > 0 {
		headers := make(map[string]string)
		for _, h := range apiHeaders {
			parts := strings.SplitN(h, ":", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid header format: %s (use key:value)", h)
			}
			headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
		opts.Headers = headers
	}

	// Make the request
	resp, err := service.Request(cmd.Context(), method, path, opts)
	if err != nil {
		return err
	}

	// Output the response
	return outputAPIResponse(cmd, resp)
}

func outputAPIResponse(cmd *cobra.Command, resp *api.RawResponse) error {
	// If raw output, just print the body
	if apiRawOutput {
		cmd.Println(string(resp.Body))
		return nil
	}

	// Build output structure
	output := make(map[string]interface{})
	output["status_code"] = resp.StatusCode

	if apiShowHeaders {
		headers := make(map[string]string)
		for key, values := range resp.Headers {
			if len(values) > 0 {
				headers[key] = values[0]
			}
		}
		output["headers"] = headers
	}

	// Parse body as JSON if possible
	if len(resp.Body) > 0 {
		var body interface{}
		if err := json.Unmarshal(resp.Body, &body); err == nil {
			output["body"] = body
		} else {
			output["body"] = string(resp.Body)
		}
	}

	// Add pagination info if available
	if resp.Pagination != nil && resp.Pagination.HasNextPage() {
		output["pagination"] = map[string]interface{}{
			"has_next": resp.Pagination.HasNextPage(),
			"next":     resp.Pagination.Next,
		}
	}

	// Format output based on output format flag
	return formatOutput(output, nil)
}
