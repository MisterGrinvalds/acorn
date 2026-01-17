package configfile

import (
	"encoding/json"
)

// JSONWriter implements the Writer interface for JSON format.
// Used for VS Code settings.json and similar config files.
type JSONWriter struct{}

func init() {
	Register(&JSONWriter{})
}

// Format returns the format identifier.
func (w *JSONWriter) Format() string {
	return "json"
}

// Write generates JSON content from values.
// The values map is directly marshaled to pretty-printed JSON.
func (w *JSONWriter) Write(values map[string]interface{}) ([]byte, error) {
	return json.MarshalIndent(values, "", "  ")
}

// JSONArrayWriter implements the Writer interface for JSON array format.
// Used for VS Code keybindings.json and similar array-based config files.
type JSONArrayWriter struct{}

func init() {
	Register(&JSONArrayWriter{})
}

// Format returns the format identifier.
func (w *JSONArrayWriter) Format() string {
	return "jsonarray"
}

// Write generates JSON array content from values.
// Expects values to have an "items" key containing the array.
func (w *JSONArrayWriter) Write(values map[string]interface{}) ([]byte, error) {
	items, ok := values["items"]
	if !ok {
		// If no items key, treat the entire values as a single object in an array
		return json.MarshalIndent([]interface{}{values}, "", "  ")
	}

	// Check if items is already a slice
	switch v := items.(type) {
	case []interface{}:
		return json.MarshalIndent(v, "", "  ")
	case []map[string]interface{}:
		return json.MarshalIndent(v, "", "  ")
	default:
		// Wrap in array if not already
		return json.MarshalIndent([]interface{}{items}, "", "  ")
	}
}
