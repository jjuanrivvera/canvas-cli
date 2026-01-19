package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

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
	return NewFormatterWithOptions(format, false)
}

// NewFormatterWithOptions creates a new formatter with additional options
func NewFormatterWithOptions(format FormatType, verbose bool) (Formatter, error) {
	switch format {
	case FormatJSON:
		return &JSONFormatter{}, nil
	case FormatYAML:
		return &YAMLFormatter{}, nil
	case FormatCSV:
		return &CSVFormatter{}, nil
	case FormatTable:
		return &TableFormatter{Verbose: verbose}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// JSONFormatter formats output as JSON
type JSONFormatter struct{}

// Format formats data as JSON
func (f *JSONFormatter) Format(data interface{}) (string, error) {
	// Handle nil by returning null
	if data == nil {
		return "null", nil
	}

	// Handle empty slices - ensure we output [] instead of null
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Slice && v.IsNil() {
		return "[]", nil
	}

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
	// Handle nil by returning empty list
	if data == nil {
		return "[]\n", nil
	}

	// Handle empty slices - ensure we output [] instead of null
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Slice && v.IsNil() {
		return "[]\n", nil
	}

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
type TableFormatter struct {
	Verbose bool // When false, only show key fields
}

// keyFieldsMap defines which fields to show by default for each struct type
// These are the most commonly needed fields for quick viewing
var keyFieldsMap = map[string][]string{
	// Course fields - expanded for better overview
	"Course": {"id", "name", "course_code", "workflow_state", "account_id", "enrollment_term_id", "start_at", "end_at", "default_view"},
	// User fields - expanded with sortable name
	"User": {"id", "name", "sortable_name", "email", "login_id", "created_at"},
	// Assignment fields - expanded with grading info
	"Assignment": {"id", "name", "due_at", "points_possible", "grading_type", "published", "submission_types", "unlock_at", "lock_at"},
	// Submission fields - expanded with more context
	"Submission": {"id", "user_id", "assignment_id", "score", "grade", "workflow_state", "submitted_at", "graded_at", "late"},
	// Module fields - good as is
	"Module": {"id", "name", "position", "workflow_state", "published", "items_count", "unlock_at"},
	// ModuleItem fields - expanded
	"ModuleItem": {"id", "title", "type", "position", "module_id", "content_id", "indent", "published"},
	// Page fields - good as is
	"Page": {"url", "title", "created_at", "updated_at", "published", "front_page", "editing_roles"},
	// Attachment (File) fields - expanded
	"Attachment": {"id", "display_name", "filename", "content-type", "size", "created_at", "updated_at", "folder_id", "hidden"},
	// DiscussionTopic fields - removed long message field
	"DiscussionTopic": {"id", "title", "posted_at", "published", "discussion_type", "discussion_subentry_count", "user_name", "locked"},
	// DiscussionEntry fields - condensed message
	"DiscussionEntry": {"id", "user_id", "user_name", "created_at", "read_state", "rating_count"},
	// CalendarEvent fields - expanded with location
	"CalendarEvent": {"id", "title", "start_at", "end_at", "all_day", "location_name", "context_code", "workflow_state"},
	// Enrollment fields - expanded with activity
	"Enrollment": {"id", "user_id", "course_id", "type", "enrollment_state", "role", "created_at", "last_activity_at"},
	// PlannerItem fields
	"PlannerItem": {"plannable_id", "plannable_type", "plannable_date", "context_name", "context_type", "planner_override"},
	// PlannerNote fields
	"PlannerNote": {"id", "title", "todo_date", "course_id", "workflow_state"},
	// PlannerOverride fields
	"PlannerOverride": {"id", "plannable_id", "plannable_type", "marked_complete", "dismissed"},
	// Folder fields - expanded
	"Folder": {"id", "name", "full_name", "parent_folder_id", "files_count", "folders_count", "created_at", "hidden"},
	// Account fields - expanded
	"Account": {"id", "name", "workflow_state", "parent_account_id", "root_account_id", "default_time_zone"},
	// PageRevision fields
	"PageRevision": {"revision_id", "updated_at", "title", "edited_by"},
}

// Format formats data as a table
func (f *TableFormatter) Format(data interface{}) (string, error) {
	slice := toSlice(data)
	if len(slice) == 0 {
		return "No data", nil
	}

	// Get headers
	headers := getHeaders(slice[0])

	// Filter headers if not verbose
	if !f.Verbose {
		headers = f.filterKeyFields(slice[0], headers)
	}

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
	_, _ = builder.WriteString("┌")
	for i, width := range widths {
		_, _ = builder.WriteString(strings.Repeat("─", width+2))
		if i < len(widths)-1 {
			_, _ = builder.WriteString("┬")
		}
	}
	_, _ = builder.WriteString("┐\n")

	// Header content
	_, _ = builder.WriteString("│")
	for i, header := range headers {
		_, _ = builder.WriteString(" ")
		_, _ = builder.WriteString(header)
		_, _ = builder.WriteString(strings.Repeat(" ", widths[i]-len(header)))
		_, _ = builder.WriteString(" │")
	}
	_, _ = builder.WriteString("\n")

	// Separator
	_, _ = builder.WriteString("├")
	for i, width := range widths {
		_, _ = builder.WriteString(strings.Repeat("─", width+2))
		if i < len(widths)-1 {
			_, _ = builder.WriteString("┼")
		}
	}
	_, _ = builder.WriteString("┤\n")

	// Data rows
	for _, row := range rows {
		_, _ = builder.WriteString("│")
		for i, cell := range row {
			_, _ = builder.WriteString(" ")
			_, _ = builder.WriteString(cell)
			_, _ = builder.WriteString(strings.Repeat(" ", widths[i]-len(cell)))
			_, _ = builder.WriteString(" │")
		}
		_, _ = builder.WriteString("\n")
	}

	// Footer
	_, _ = builder.WriteString("└")
	for i, width := range widths {
		_, _ = builder.WriteString(strings.Repeat("─", width+2))
		if i < len(widths)-1 {
			_, _ = builder.WriteString("┴")
		}
	}
	_, _ = builder.WriteString("┘\n")

	return builder.String(), nil
}

// filterKeyFields filters headers to only include key fields for the given item type
func (f *TableFormatter) filterKeyFields(item interface{}, allHeaders []string) []string {
	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Get struct type name
	typeName := ""
	if v.Kind() == reflect.Struct {
		typeName = v.Type().Name()
	}

	// Look up key fields for this type
	keyFields, found := keyFieldsMap[typeName]
	if !found || len(keyFields) == 0 {
		// If no key fields defined, show first 6 fields as default
		if len(allHeaders) > 6 {
			return allHeaders[:6]
		}
		return allHeaders
	}

	// Filter to only include key fields that exist in the struct
	keyFieldSet := make(map[string]bool)
	for _, field := range keyFields {
		keyFieldSet[field] = true
	}

	filtered := make([]string, 0, len(keyFields))
	for _, header := range allHeaders {
		if keyFieldSet[header] {
			filtered = append(filtered, header)
		}
	}

	// If no matching fields found, return original headers limited to 6
	if len(filtered) == 0 {
		if len(allHeaders) > 6 {
			return allHeaders[:6]
		}
		return allHeaders
	}

	return filtered
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

	// Handle time.Time specially - show "Not set" for zero dates
	if t, ok := v.(time.Time); ok {
		if t.IsZero() {
			return "Not set"
		}
		return t.Format("2006-01-02 15:04")
	}

	// Handle *time.Time
	if t, ok := v.(*time.Time); ok {
		if t == nil || t.IsZero() {
			return "Not set"
		}
		return t.Format("2006-01-02 15:04")
	}

	val := reflect.ValueOf(v)

	// Handle nil pointers
	if val.Kind() == reflect.Ptr && val.IsNil() {
		return ""
	}

	// Handle map type - show empty string for empty maps instead of "map[]"
	if val.Kind() == reflect.Map {
		if val.Len() == 0 {
			return ""
		}
		// For non-empty maps, format as key:value pairs
		items := make([]string, 0, val.Len())
		for _, key := range val.MapKeys() {
			items = append(items, fmt.Sprintf("%v:%v", key.Interface(), val.MapIndex(key).Interface()))
		}
		return "{" + strings.Join(items, ", ") + "}"
	}

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
			return ""
		}
		items := make([]string, val.Len())
		for i := 0; i < val.Len(); i++ {
			items[i] = formatValue(val.Index(i).Interface())
		}
		return "[" + strings.Join(items, ", ") + "]"
	case reflect.Struct:
		// For structs, try to find a name/title/id field to display
		return formatStructCompact(val)
	case reflect.Ptr:
		// Handle pointers to structs
		if !val.IsNil() && val.Elem().Kind() == reflect.Struct {
			return formatStructCompact(val.Elem())
		}
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// formatStructCompact formats a struct as a compact string by showing key identifying fields
func formatStructCompact(val reflect.Value) string {
	t := val.Type()

	// Try common identifying field names in priority order
	identifyingFields := []string{"name", "Name", "title", "Title", "id", "ID", "Id"}

	for _, fieldName := range identifyingFields {
		field := val.FieldByName(fieldName)
		if field.IsValid() && !field.IsZero() {
			// For ID fields, show "Type(ID)" format
			if strings.ToLower(fieldName) == "id" {
				return fmt.Sprintf("%s(%v)", t.Name(), field.Interface())
			}
			// For name/title fields, just show the value
			return fmt.Sprintf("%v", field.Interface())
		}
	}

	// If we found an ID and a name, show both
	idField := val.FieldByName("ID")
	if !idField.IsValid() {
		idField = val.FieldByName("Id")
	}
	nameField := val.FieldByName("Name")
	if !nameField.IsValid() {
		nameField = val.FieldByName("name")
	}

	if idField.IsValid() && !idField.IsZero() && nameField.IsValid() && !nameField.IsZero() {
		return fmt.Sprintf("%v (%v)", nameField.Interface(), idField.Interface())
	}

	// Fallback: show type name with first non-zero field
	for i := 0; i < t.NumField(); i++ {
		field := val.Field(i)
		if !field.IsZero() {
			return fmt.Sprintf("%s{%s: %v}", t.Name(), t.Field(i).Name, field.Interface())
		}
	}

	return t.Name() + "{}"
}

// Write writes formatted data to a writer
func Write(w io.Writer, data interface{}, format FormatType) error {
	return WriteWithOptions(w, data, format, false)
}

// WriteWithOptions writes formatted data to a writer with verbose option
func WriteWithOptions(w io.Writer, data interface{}, format FormatType, verbose bool) error {
	formatter, err := NewFormatterWithOptions(format, verbose)
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
