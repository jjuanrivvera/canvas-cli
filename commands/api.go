package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

func init() {
	rootCmd.AddCommand(newAPICmd())
}

func newAPICmd() *cobra.Command {
	opts := &options.APIOptions{}

	cmd := &cobra.Command{
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
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runAPICommand(cmd, args, client, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Data, "data", "d", "", "JSON data for request body")
	cmd.Flags().StringVar(&opts.DataFile, "data-file", "", "Read JSON data from file")
	cmd.Flags().StringArrayVarP(&opts.Query, "query", "q", nil, "Query parameters (key=value, repeatable)")
	cmd.Flags().StringArrayVarP(&opts.Headers, "header", "H", nil, "Custom headers (key:value, repeatable)")
	cmd.Flags().BoolVar(&opts.Paginate, "paginate", false, "Follow pagination links (GET only)")
	cmd.Flags().BoolVar(&opts.RawOutput, "raw", false, "Output raw response without formatting")
	cmd.Flags().BoolVar(&opts.ShowHeaders, "show-headers", false, "Include response headers in output")

	cmd.AddCommand(newAPIGetCmd())

	return cmd
}

// newAPIGetCmd is a GET-only sibling of "canvas api". Because it can never
// mutate state, it is safe to advertise to read-only MCP clients: the shared
// classifier (classifyCanvasCommand) buckets "api get" as a read, so it carries
// readOnlyHint=true while the general "canvas api" escape hatch (any HTTP verb)
// stays unannotated. It gives broad Canvas read coverage from a single tool
// schema instead of allowlisting every typed read tool. See issue #60.
func newAPIGetCmd() *cobra.Command {
	opts := &options.APIOptions{}

	cmd := &cobra.Command{
		Use:   "get <PATH>",
		Short: "Make a read-only GET request to Canvas",
		Long: `Make a raw GET request to any Canvas API endpoint.

A GET-only sibling of "canvas api": it never mutates state, so it is exposed to
read-only MCP clients (it carries the readOnlyHint annotation). For other verbs,
use "canvas api <METHOD> <PATH>".

Examples:
  # List all courses
  canvas api get /api/v1/courses

  # Search users with query parameters
  canvas api get /api/v1/users -q "search_term=john" -q "per_page=50"

  # Follow pagination
  canvas api get /api/v1/courses --paginate`,
		Args: ExactArgsWithUsage(1, "path"),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			// Force the method to GET; reuse the shared runner.
			return runAPICommand(cmd, []string{"GET", args[0]}, client, opts)
		},
	}

	cmd.Flags().StringArrayVarP(&opts.Query, "query", "q", nil, "Query parameters (key=value, repeatable)")
	cmd.Flags().StringArrayVarP(&opts.Headers, "header", "H", nil, "Custom headers (key:value, repeatable)")
	cmd.Flags().BoolVar(&opts.Paginate, "paginate", false, "Follow pagination links")
	cmd.Flags().BoolVar(&opts.RawOutput, "raw", false, "Output raw response without formatting")
	cmd.Flags().BoolVar(&opts.ShowHeaders, "show-headers", false, "Include response headers in output")

	return cmd
}

func runAPICommand(cmd *cobra.Command, args []string, client *api.Client, opts *options.APIOptions) error {
	logger := logging.NewCommandLogger(verbose)
	ctx := cmd.Context()

	method := strings.ToUpper(args[0])
	path := args[1]

	// Validate method
	switch method {
	case "GET", "POST", "PUT", "DELETE", "PATCH", "HEAD":
		// Valid
	default:
		return fmt.Errorf("unsupported HTTP method: %s (use GET, POST, PUT, DELETE, PATCH, or HEAD)", method)
	}

	logger.LogCommandStart(ctx, "api.request", map[string]interface{}{
		"method": method,
		"path":   path,
	})

	service := api.NewRawService(client)

	// Build request options
	reqOpts := &api.RawRequestOptions{
		Paginate: opts.Paginate && method == "GET",
	}

	// Parse body from --data or --data-file.
	// Use Flags().Changed() rather than empty-string checks so that stale
	// values from a previous cobra Execute() call do not leak between tests.
	dataChanged := cmd.Flags().Changed("data")
	dataFileChanged := cmd.Flags().Changed("data-file")

	if dataChanged && dataFileChanged {
		return fmt.Errorf("cannot use both --data and --data-file")
	}

	if dataChanged {
		var body interface{}
		if err := json.Unmarshal([]byte(opts.Data), &body); err != nil {
			return fmt.Errorf("invalid JSON in --data: %w", err)
		}
		reqOpts.Body = body
	}

	if dataFileChanged {
		var reader io.Reader
		if opts.DataFile == "-" {
			reader = cmd.InOrStdin()
		} else {
			file, err := os.Open(opts.DataFile)
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
		reqOpts.Body = body
	}

	// Parse query parameters.
	// Guard with Changed() to avoid accumulating values from previous Execute() calls
	// when the cobra command is reused (e.g., in tests).
	if cmd.Flags().Changed("query") {
		query := make(map[string][]string)
		for _, q := range opts.Query {
			parts := strings.SplitN(q, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid query parameter format: %s (use key=value)", q)
			}
			key := parts[0]
			value := parts[1]
			query[key] = append(query[key], value)
		}
		reqOpts.Query = query
	}

	// Parse custom headers.
	// Same Changed() guard as for query params.
	if cmd.Flags().Changed("header") {
		headers := make(map[string]string)
		for _, h := range opts.Headers {
			parts := strings.SplitN(h, ":", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid header format: %s (use key:value)", h)
			}
			headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
		reqOpts.Headers = headers
	}

	// Make the request
	resp, err := service.Request(ctx, method, path, reqOpts)
	if err != nil {
		logger.LogCommandError(ctx, "api.request", err, map[string]interface{}{
			"method": method,
			"path":   path,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "api.request", 1)
	return outputAPIResponse(cmd, resp, opts)
}

func outputAPIResponse(cmd *cobra.Command, resp *api.RawResponse, opts *options.APIOptions) error {
	// If raw output, just print the body
	if opts.RawOutput {
		cmd.Println(string(resp.Body))
		return nil
	}

	// Build output structure
	out := make(map[string]interface{})
	out["status_code"] = resp.StatusCode

	if opts.ShowHeaders {
		headers := make(map[string]string)
		for key, values := range resp.Headers {
			if len(values) > 0 {
				headers[key] = values[0]
			}
		}
		out["headers"] = headers
	}

	// Parse body as JSON if possible
	if len(resp.Body) > 0 {
		var body interface{}
		if err := json.Unmarshal(resp.Body, &body); err == nil {
			out["body"] = body
		} else {
			out["body"] = string(resp.Body)
		}
	}

	// Add pagination info if available
	if resp.Pagination != nil && resp.Pagination.HasNextPage() {
		out["pagination"] = map[string]interface{}{
			"has_next": resp.Pagination.HasNextPage(),
			"next":     resp.Pagination.Next,
		}
	}

	// Format output based on output format flag
	return formatOutput(out, nil)
}
