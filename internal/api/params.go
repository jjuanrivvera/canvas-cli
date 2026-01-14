package api

import (
	"fmt"
	"net/url"
	"reflect"
)

// ParamsBuilder helps build API request parameters from structs
// It automatically handles optional fields (pointers) and zero values
type ParamsBuilder struct {
	params map[string]interface{}
	prefix string
}

// NewParamsBuilder creates a new parameter builder
func NewParamsBuilder() *ParamsBuilder {
	return &ParamsBuilder{
		params: make(map[string]interface{}),
	}
}

// WithPrefix returns a new builder with the given prefix for all keys
func (b *ParamsBuilder) WithPrefix(prefix string) *ParamsBuilder {
	return &ParamsBuilder{
		params: b.params,
		prefix: prefix,
	}
}

// Set adds a value if it's not a zero value
// For pointers, checks if pointer is non-nil and adds dereferenced value
// For strings, checks if non-empty
// For numbers, checks if non-zero
// For slices/maps, checks if non-empty
func (b *ParamsBuilder) Set(key string, value interface{}) *ParamsBuilder {
	if b.isZeroValue(value) {
		return b
	}

	fullKey := key
	if b.prefix != "" {
		fullKey = b.prefix + "[" + key + "]"
	}

	// Dereference pointer if needed
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		value = v.Elem().Interface()
	}

	b.params[fullKey] = value
	return b
}

// SetBool adds a boolean value only if it's true (for non-pointer bools)
// or if the pointer is non-nil (for pointer bools)
func (b *ParamsBuilder) SetBool(key string, value interface{}) *ParamsBuilder {
	fullKey := key
	if b.prefix != "" {
		fullKey = b.prefix + "[" + key + "]"
	}

	v := reflect.ValueOf(value)

	// Handle pointer to bool
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return b
		}
		b.params[fullKey] = v.Elem().Bool()
		return b
	}

	// Handle plain bool - only set if true
	if v.Kind() == reflect.Bool && v.Bool() {
		b.params[fullKey] = true
	}

	return b
}

// SetSlice adds a slice if it's non-empty
func (b *ParamsBuilder) SetSlice(key string, value interface{}) *ParamsBuilder {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice || v.Len() == 0 {
		return b
	}

	fullKey := key
	if b.prefix != "" {
		fullKey = b.prefix + "[" + key + "]"
	}

	b.params[fullKey] = value
	return b
}

// SetMap adds a map if it's non-empty
func (b *ParamsBuilder) SetMap(key string, value interface{}) *ParamsBuilder {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Map || v.Len() == 0 {
		return b
	}

	fullKey := key
	if b.prefix != "" {
		fullKey = b.prefix + "[" + key + "]"
	}

	b.params[fullKey] = value
	return b
}

// Build returns the accumulated parameters as a map
func (b *ParamsBuilder) Build() map[string]interface{} {
	return b.params
}

// isZeroValue checks if a value is its zero value or an empty container
func (b *ParamsBuilder) isZeroValue(value interface{}) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.Ptr:
		return v.IsNil()
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() == 0
	default:
		return false
	}
}

// URLParamsBuilder builds URL query parameters
type URLParamsBuilder struct {
	values url.Values
}

// NewURLParamsBuilder creates a new URL parameter builder
func NewURLParamsBuilder() *URLParamsBuilder {
	return &URLParamsBuilder{
		values: make(url.Values),
	}
}

// Set adds a string parameter if non-empty
func (b *URLParamsBuilder) Set(key, value string) *URLParamsBuilder {
	if value != "" {
		b.values.Set(key, value)
	}
	return b
}

// SetInt adds an integer parameter if non-zero
func (b *URLParamsBuilder) SetInt(key string, value int) *URLParamsBuilder {
	if value != 0 {
		b.values.Set(key, fmt.Sprintf("%d", value))
	}
	return b
}

// SetInt64 adds an int64 parameter if non-zero
func (b *URLParamsBuilder) SetInt64(key string, value int64) *URLParamsBuilder {
	if value != 0 {
		b.values.Set(key, fmt.Sprintf("%d", value))
	}
	return b
}

// SetBool adds a boolean parameter (always sets)
func (b *URLParamsBuilder) SetBool(key string, value bool) *URLParamsBuilder {
	b.values.Set(key, fmt.Sprintf("%t", value))
	return b
}

// SetBoolPtr adds a boolean parameter if pointer is non-nil
func (b *URLParamsBuilder) SetBoolPtr(key string, value *bool) *URLParamsBuilder {
	if value != nil {
		b.values.Set(key, fmt.Sprintf("%t", *value))
	}
	return b
}

// Add appends a value (for repeated parameters)
func (b *URLParamsBuilder) Add(key, value string) *URLParamsBuilder {
	if value != "" {
		b.values.Add(key, value)
	}
	return b
}

// AddSlice adds multiple values for the same key
func (b *URLParamsBuilder) AddSlice(key string, values []string) *URLParamsBuilder {
	for _, v := range values {
		if v != "" {
			b.values.Add(key, v)
		}
	}
	return b
}

// Build returns the url.Values
func (b *URLParamsBuilder) Build() url.Values {
	return b.values
}

// Encode returns the URL-encoded string
func (b *URLParamsBuilder) Encode() string {
	return b.values.Encode()
}
