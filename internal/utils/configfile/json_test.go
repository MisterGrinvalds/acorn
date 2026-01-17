package configfile

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestJSONWriterFormat(t *testing.T) {
	w := &JSONWriter{}
	if w.Format() != "json" {
		t.Errorf("Format() = %q, want %q", w.Format(), "json")
	}
}

func TestJSONWriterBasicValues(t *testing.T) {
	w := &JSONWriter{}

	values := map[string]interface{}{
		"editor.fontSize":      14,
		"editor.fontFamily":    "JetBrainsMono Nerd Font",
		"editor.fontLigatures": true,
	}

	content, err := w.Write(values)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	// Check values
	if result["editor.fontSize"] != float64(14) {
		t.Errorf("editor.fontSize = %v, want 14", result["editor.fontSize"])
	}

	if result["editor.fontFamily"] != "JetBrainsMono Nerd Font" {
		t.Errorf("editor.fontFamily = %v, want JetBrainsMono Nerd Font", result["editor.fontFamily"])
	}

	if result["editor.fontLigatures"] != true {
		t.Errorf("editor.fontLigatures = %v, want true", result["editor.fontLigatures"])
	}
}

func TestJSONWriterNestedValues(t *testing.T) {
	w := &JSONWriter{}

	values := map[string]interface{}{
		"files.exclude": map[string]interface{}{
			"**/.git":         true,
			"**/node_modules": true,
		},
	}

	content, err := w.Write(values)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	// Check nested value
	exclude, ok := result["files.exclude"].(map[string]interface{})
	if !ok {
		t.Fatalf("files.exclude is not a map")
	}

	if exclude["**/.git"] != true {
		t.Errorf("files.exclude[**/.git] = %v, want true", exclude["**/.git"])
	}
}

func TestJSONWriterLanguageSpecific(t *testing.T) {
	w := &JSONWriter{}

	// VS Code language-specific settings use bracket syntax
	values := map[string]interface{}{
		"[python]": map[string]interface{}{
			"editor.tabSize":      4,
			"editor.formatOnSave": true,
		},
	}

	content, err := w.Write(values)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	// Check language-specific settings
	python, ok := result["[python]"].(map[string]interface{})
	if !ok {
		t.Fatalf("[python] is not a map")
	}

	if python["editor.tabSize"] != float64(4) {
		t.Errorf("[python].editor.tabSize = %v, want 4", python["editor.tabSize"])
	}
}

func TestJSONWriterPrettyPrinted(t *testing.T) {
	w := &JSONWriter{}

	values := map[string]interface{}{
		"key": "value",
	}

	content, err := w.Write(values)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Should be pretty-printed with indentation
	if !strings.Contains(string(content), "\n") {
		t.Error("Output should be pretty-printed with newlines")
	}

	if !strings.Contains(string(content), "  ") {
		t.Error("Output should have indentation")
	}
}

func TestJSONArrayWriterFormat(t *testing.T) {
	w := &JSONArrayWriter{}
	if w.Format() != "jsonarray" {
		t.Errorf("Format() = %q, want %q", w.Format(), "jsonarray")
	}
}

func TestJSONArrayWriterWithItems(t *testing.T) {
	w := &JSONArrayWriter{}

	values := map[string]interface{}{
		"items": []interface{}{
			map[string]interface{}{
				"key":     "cmd+p",
				"command": "workbench.action.quickOpen",
			},
			map[string]interface{}{
				"key":     "cmd+shift+p",
				"command": "workbench.action.showCommands",
			},
		},
	}

	content, err := w.Write(values)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Verify it's valid JSON array
	var result []interface{}
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatalf("Output is not valid JSON array: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 items, got %d", len(result))
	}

	// Check first item
	first, ok := result[0].(map[string]interface{})
	if !ok {
		t.Fatal("First item is not a map")
	}

	if first["key"] != "cmd+p" {
		t.Errorf("First item key = %v, want cmd+p", first["key"])
	}
}

func TestJSONArrayWriterWithWhen(t *testing.T) {
	w := &JSONArrayWriter{}

	// VS Code keybindings can have optional "when" clause
	values := map[string]interface{}{
		"items": []interface{}{
			map[string]interface{}{
				"key":     "cmd+shift+k",
				"command": "editor.action.deleteLines",
				"when":    "editorTextFocus && !editorReadonly",
			},
		},
	}

	content, err := w.Write(values)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	var result []interface{}
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatalf("Output is not valid JSON array: %v", err)
	}

	first := result[0].(map[string]interface{})
	if first["when"] != "editorTextFocus && !editorReadonly" {
		t.Errorf("when clause not preserved: %v", first["when"])
	}
}

func TestJSONArrayWriterEmptyItems(t *testing.T) {
	w := &JSONArrayWriter{}

	values := map[string]interface{}{
		"items": []interface{}{},
	}

	content, err := w.Write(values)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	var result []interface{}
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatalf("Output is not valid JSON array: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected empty array, got %d items", len(result))
	}
}

func TestJSONArrayWriterNoItemsKey(t *testing.T) {
	w := &JSONArrayWriter{}

	// If no items key, should wrap in array
	values := map[string]interface{}{
		"key":     "cmd+p",
		"command": "test",
	}

	content, err := w.Write(values)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	var result []interface{}
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatalf("Output is not valid JSON array: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 item, got %d", len(result))
	}
}
