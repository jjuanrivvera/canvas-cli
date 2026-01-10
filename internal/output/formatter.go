package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

// Formatter defines the interface for output formatters
type Formatter interface {
	// Format formats data for output
	Format(data interface{}) (string, error)
}

// FormatType represents the output format type
type FormatType string

const (
	FormatJSON  FormatType = "json"
	FormatYAML  FormatType = "yaml"
	FormatCSV   FormatType = "csv"
	FormatTable FormatType = "table"
)

// NewFormatter creates a new formatter for the specified format type
func NewFormatter(format FormatType) (Formatter, error) {
	switch format {
	case FormatJSON:
		return &JSONFormatter{}, nil
	case FormatYAML:
		return &YAMLFormatter{}, nil
	case FormatCSV:
		return &CSVFormatter{}, nil
	case FormatTable:
		return &TableFormatter{}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// JSONFormatter formats output as JSON
type JSONFormatter struct{}

// Format formats data as JSON
func (f *JSONFormatter) Format(data interface{}) (string, error) {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(bytes), nil
}

// YAMLFormatter formats output as YAML
type YAMLFormatter struct{}

// Format formats data as YAML
func (f *YAMLFormatter) Format(data interface{}) (string, error) {
	bytes, err := yaml.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal YAML: %w", err)
	}
	return string(bytes), nil
}

// CSVFormatter formats output as CSV
type CSVFormatter struct{}

// Format formats data as CSV
func (f *CSVFormatter) Format(data interface{}) (string, error) {
	var builder strings.Builder
	writer := csv.NewWriter(&builder)

	// Convert data to slice of maps or structs
	slice := toSlice(data)
	if len(slice) == 0 {
		return "", nil
	}

	// Get headers from first item
	headers := getHeaders(slice[0])
	if err := writer.Write(headers); err != nil {
		return "", fmt.Errorf("failed to write CSV headers: %w", err)
	}

	// Write rows
	for _, item := range slice {
		row := getRow(item, headers)
		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("CSV writer error: %w", err)
	}

	return builder.String(), nil
}

// TableFormatter formats output as a table
type TableFormatter struct{}

// Format formats data as a table
func (f *TableFormatter) Format(data interface{}) (string, error) {
	slice := toSlice(data)
	if len(slice) == 0 {
		return "No data", nil
	}

	// Get headers
	headers := getHeaders(slice[0])

	// Calculate column widths
	widths := make([]int, len(headers))
	for i, header := range headers {
		widths[i] = len(header)
	}

	// Get all rows and update widths
	rows := make([][]string, len(slice))
	for i, item := range slice {
		row := getRow(item, headers)
		rows[i] = row
		for j, cell := range row {
			if len(cell) > widths[j] {
				widths[j] = len(cell)
			}
		}
	}

	// Build table
	var builder strings.Builder

	// Header row
	builder.WriteString("┌")
	for i, width := range widths {
		builder.WriteString(strings.Repeat("─", width+2))
		if i < len(widths)-1 {
			builder.WriteString("┬")
		}
	}
	builder.WriteString("┐\n")

	// Header content
	builder.WriteString("│")
	for i, header := range headers {
		builder.WriteString(" ")
		builder.WriteString(header)
		builder.WriteString(strings.Repeat(" ", widths[i]-len(header)))
		builder.WriteString(" │")
	}
	builder.WriteString("\n")

	// Separator
	builder.WriteString("├")
	for i, width := range widths {
		builder.WriteString(strings.Repeat("─", width+2))
		if i < len(widths)-1 {
			builder.WriteString("┼")
		}
	}
	builder.WriteString("┤\n")

	// Data rows
	for _, row := range rows {
		builder.WriteString("│")
		for i, cell := range row {
			builder.WriteString(" ")
			builder.WriteString(cell)
			builder.WriteString(strings.Repeat(" ", widths[i]-len(cell)))
			builder.WriteString(" │")
		}
		builder.WriteString("\n")
	}

	// Footer
	builder.WriteString("└")
	for i, width := range widths {
		builder.WriteString(strings.Repeat("─", width+2))
		if i < len(widths)-1 {
			builder.WriteString("┴")
		}
	}
	builder.WriteString("┘\n")

	return builder.String(), nil
}

// toSlice converts data to a slice of interface{}
func toSlice(data interface{}) []interface{} {
	v := reflect.ValueOf(data)

	// If it's already a slice, convert to []interface{}
	if v.Kind() == reflect.Slice {
		result := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			result[i] = v.Index(i).Interface()
		}
		return result
	}

	// If it's a single item, wrap it in a slice
	return []interface{}{data}
}

// getHeaders extracts field names from a struct or map
func getHeaders(item interface{}) []string {
	v := reflect.ValueOf(item)

	// Handle pointers
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		headers := make([]string, 0, t.NumField())
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			// Use json tag if available, otherwise use field name
			name := field.Name
			if tag := field.Tag.Get("json"); tag != "" && tag != "-" {
				// Extract field name from json tag (ignore options like omitempty)
				if idx := strings.Index(tag, ","); idx != -1 {
					name = tag[:idx]
				} else {
					name = tag
				}
			}
			headers = append(headers, name)
		}
		return headers

	case reflect.Map:
		headers := make([]string, 0, v.Len())
		for _, key := range v.MapKeys() {
			headers = append(headers, fmt.Sprintf("%v", key.Interface()))
		}
		return headers
	}

	return []string{}
}

// getRow extracts values from a struct or map based on headers
func getRow(item interface{}, headers []string) []string {
	v := reflect.ValueOf(item)

	// Handle pointers
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	row := make([]string, len(headers))

	switch v.Kind() {
	case reflect.Struct:
		for i, header := range headers {
			// Find field by json tag or name
			field := findField(v, header)
			if field.IsValid() {
				row[i] = formatValue(field.Interface())
			}
		}

	case reflect.Map:
		for i, header := range headers {
			key := reflect.ValueOf(header)
			value := v.MapIndex(key)
			if value.IsValid() {
				row[i] = formatValue(value.Interface())
			}
		}
	}

	return row
}

// findField finds a struct field by json tag or name
func findField(v reflect.Value, name string) reflect.Value {
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Check json tag
		if tag := field.Tag.Get("json"); tag != "" {
			tagName := tag
			if idx := strings.Index(tag, ","); idx != -1 {
				tagName = tag[:idx]
			}
			if tagName == name {
				return v.Field(i)
			}
		}

		// Check field name
		if field.Name == name {
			return v.Field(i)
		}
	}

	return reflect.Value{}
}

// formatValue formats a value as a string for display
func formatValue(v interface{}) string {
	if v == nil {
		return ""
	}

	val := reflect.ValueOf(v)

	switch val.Kind() {
	case reflect.String:
		return val.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", val.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%.2f", val.Float())
	case reflect.Bool:
		if val.Bool() {
			return "true"
		}
		return "false"
	case reflect.Slice, reflect.Array:
		if val.Len() == 0 {
			return "[]"
		}
		items := make([]string, val.Len())
		for i := 0; i < val.Len(); i++ {
			items[i] = formatValue(val.Index(i).Interface())
		}
		return "[" + strings.Join(items, ", ") + "]"
	default:
		return fmt.Sprintf("%v", v)
	}
}

// Write writes formatted data to a writer
func Write(w io.Writer, data interface{}, format FormatType) error {
	formatter, err := NewFormatter(format)
	if err != nil {
		return err
	}

	output, err := formatter.Format(data)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(output))
	return err
}
