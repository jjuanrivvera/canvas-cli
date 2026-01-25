package commands

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// applyFiltering applies filter, columns, and sort operations to data
// Returns the filtered/sorted data
func applyFiltering(data interface{}) interface{} {
	// Only apply to slices
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return data
	}

	if v.Len() == 0 {
		return data
	}

	// Convert to []map[string]interface{} for easier manipulation
	items := toMapSlice(data)
	if items == nil {
		return data
	}

	// Apply text filter
	if filterText != "" {
		items = filterByText(items, filterText)
	}

	// Apply sorting
	if sortField != "" {
		items = sortByField(items, sortField)
	}

	// Apply column selection
	if len(filterColumns) > 0 {
		items = selectColumns(items, filterColumns)
	}

	return items
}

// toMapSlice converts a slice of structs to []map[string]interface{}
func toMapSlice(data interface{}) []map[string]interface{} {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return nil
	}

	result := make([]map[string]interface{}, 0, v.Len())
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i).Interface()
		m := structToMap(item)
		if m != nil {
			result = append(result, m)
		}
	}
	return result
}

// structToMap converts a struct to map[string]interface{}
func structToMap(data interface{}) map[string]interface{} {
	// Handle if already a map
	if m, ok := data.(map[string]interface{}); ok {
		return m
	}

	// Use JSON as intermediate format for conversion
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return nil
	}
	return result
}

// filterByText filters items that contain the text in any field
func filterByText(items []map[string]interface{}, text string) []map[string]interface{} {
	text = strings.ToLower(text)
	result := make([]map[string]interface{}, 0)

	for _, item := range items {
		if itemContainsText(item, text) {
			result = append(result, item)
		}
	}
	return result
}

// itemContainsText checks if any field in the item contains the text
func itemContainsText(item map[string]interface{}, text string) bool {
	for _, value := range item {
		if value == nil {
			continue
		}
		strVal := strings.ToLower(toString(value))
		if strings.Contains(strVal, text) {
			return true
		}
	}
	return false
}

// toString converts a value to string for searching
func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		// Format as integer if it's a whole number
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%g", val)
	case bool:
		return fmt.Sprintf("%t", val)
	case nil:
		return ""
	default:
		b, err := json.Marshal(val)
		if err != nil {
			return fmt.Sprintf("%v", val)
		}
		return string(b)
	}
}

// sortByField sorts items by a field name
// Prefix with - for descending order
func sortByField(items []map[string]interface{}, field string) []map[string]interface{} {
	descending := false
	if strings.HasPrefix(field, "-") {
		descending = true
		field = field[1:]
	}

	// Convert field name to match JSON naming (usually lowercase)
	field = strings.ToLower(field)
	fieldVariants := []string{field, strings.ReplaceAll(field, "_", ""), strings.ReplaceAll(field, "-", "_")}

	sort.SliceStable(items, func(i, j int) bool {
		var vi, vj interface{}
		for _, f := range fieldVariants {
			if v, ok := items[i][f]; ok {
				vi = v
				break
			}
		}
		for _, f := range fieldVariants {
			if v, ok := items[j][f]; ok {
				vj = v
				break
			}
		}

		result := compareValues(vi, vj)
		if descending {
			return result > 0
		}
		return result < 0
	})

	return items
}

// compareValues compares two interface values
func compareValues(a, b interface{}) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}

	// Try numeric comparison
	if na, ok := toFloat64(a); ok {
		if nb, ok := toFloat64(b); ok {
			if na < nb {
				return -1
			}
			if na > nb {
				return 1
			}
			return 0
		}
	}

	// Fall back to string comparison
	sa := strings.ToLower(toString(a))
	sb := strings.ToLower(toString(b))
	return strings.Compare(sa, sb)
}

// toFloat64 attempts to convert a value to float64
func toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case int32:
		return float64(val), true
	default:
		return 0, false
	}
}

// selectColumns filters items to only include specified columns
func selectColumns(items []map[string]interface{}, columns []string) []map[string]interface{} {
	result := make([]map[string]interface{}, len(items))

	// Normalize column names
	normalizedCols := make([]string, len(columns))
	for i, col := range columns {
		normalizedCols[i] = strings.ToLower(strings.TrimSpace(col))
	}

	for i, item := range items {
		filtered := make(map[string]interface{})
		for key, value := range item {
			keyLower := strings.ToLower(key)
			for _, col := range normalizedCols {
				if keyLower == col || strings.ReplaceAll(keyLower, "_", "") == strings.ReplaceAll(col, "_", "") {
					filtered[key] = value
					break
				}
			}
		}
		result[i] = filtered
	}

	return result
}

// hasFilteringOptions returns true if any filtering options are set
func hasFilteringOptions() bool {
	return filterText != "" || len(filterColumns) > 0 || sortField != ""
}
