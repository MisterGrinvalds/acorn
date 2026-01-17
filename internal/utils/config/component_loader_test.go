package config

import (
	"testing"
)

func TestMergeFilesEmpty(t *testing.T) {
	result := mergeFiles(nil, nil)
	if result != nil {
		t.Errorf("mergeFiles(nil, nil) = %v, want nil", result)
	}

	result = mergeFiles([]FileConfig{}, []FileConfig{})
	if result != nil {
		t.Errorf("mergeFiles(empty, empty) = %v, want nil", result)
	}
}

func TestMergeFilesBaseOnly(t *testing.T) {
	base := []FileConfig{
		{
			Target: "/path/to/config1",
			Format: "ghostty",
			Values: map[string]any{"key": "value"},
		},
	}

	result := mergeFiles(base, nil)

	if len(result) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(result))
	}

	if result[0].Target != "/path/to/config1" {
		t.Errorf("Target = %q, want %q", result[0].Target, "/path/to/config1")
	}
}

func TestMergeFilesOverrideOnly(t *testing.T) {
	override := []FileConfig{
		{
			Target: "/path/to/config2",
			Format: "json",
			Values: map[string]any{"key": "override"},
		},
	}

	result := mergeFiles(nil, override)

	if len(result) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(result))
	}

	if result[0].Target != "/path/to/config2" {
		t.Errorf("Target = %q, want %q", result[0].Target, "/path/to/config2")
	}
}

func TestMergeFilesDifferentTargets(t *testing.T) {
	base := []FileConfig{
		{Target: "/path/config1", Format: "ghostty"},
	}
	override := []FileConfig{
		{Target: "/path/config2", Format: "json"},
	}

	result := mergeFiles(base, override)

	if len(result) != 2 {
		t.Fatalf("Expected 2 files, got %d", len(result))
	}

	// Both files should be present
	targets := make(map[string]bool)
	for _, f := range result {
		targets[f.Target] = true
	}

	if !targets["/path/config1"] {
		t.Error("Missing /path/config1")
	}
	if !targets["/path/config2"] {
		t.Error("Missing /path/config2")
	}
}

func TestMergeFilesSameTargetOverrideWins(t *testing.T) {
	base := []FileConfig{
		{
			Target: "/path/config",
			Format: "ghostty",
			Values: map[string]any{"key": "base-value"},
		},
	}
	override := []FileConfig{
		{
			Target: "/path/config",
			Format: "json", // Changed format
			Values: map[string]any{"key": "override-value"},
		},
	}

	result := mergeFiles(base, override)

	if len(result) != 1 {
		t.Fatalf("Expected 1 file (merged), got %d", len(result))
	}

	// Override should win
	if result[0].Format != "json" {
		t.Errorf("Format = %q, want %q (override should win)", result[0].Format, "json")
	}

	if result[0].Values["key"] != "override-value" {
		t.Errorf("Values[key] = %v, want %q", result[0].Values["key"], "override-value")
	}
}

func TestMergeFilesMultipleMixed(t *testing.T) {
	base := []FileConfig{
		{Target: "/path/config1", Format: "ghostty"},
		{Target: "/path/config2", Format: "yaml"},
		{Target: "/path/shared", Format: "json", Values: map[string]any{"from": "base"}},
	}
	override := []FileConfig{
		{Target: "/path/config3", Format: "toml"},
		{Target: "/path/shared", Format: "ini", Values: map[string]any{"from": "override"}},
	}

	result := mergeFiles(base, override)

	if len(result) != 4 {
		t.Fatalf("Expected 4 files, got %d", len(result))
	}

	// Find the shared config and verify override won
	for _, f := range result {
		if f.Target == "/path/shared" {
			if f.Format != "ini" {
				t.Errorf("shared config Format = %q, want %q", f.Format, "ini")
			}
			if f.Values["from"] != "override" {
				t.Errorf("shared config Values[from] = %v, want %q", f.Values["from"], "override")
			}
		}
	}
}

func TestMergeConfigsWithFiles(t *testing.T) {
	base := &BaseConfig{
		Name: "test",
		Files: []FileConfig{
			{Target: "/base/config", Format: "ghostty"},
		},
	}

	override := &BaseConfig{
		Name: "override",
		Files: []FileConfig{
			{Target: "/override/config", Format: "json"},
		},
	}

	result := MergeConfigs(base, override)

	if result.Name != "override" {
		t.Errorf("Name = %q, want %q", result.Name, "override")
	}

	if len(result.Files) != 2 {
		t.Fatalf("Expected 2 files, got %d", len(result.Files))
	}
}

func TestMergeConfigsFilesOverrideSameTarget(t *testing.T) {
	base := &BaseConfig{
		Name: "test",
		Files: []FileConfig{
			{
				Target: "/shared/config",
				Format: "ghostty",
				Values: map[string]any{"theme": "dark"},
			},
		},
	}

	override := &BaseConfig{
		Files: []FileConfig{
			{
				Target: "/shared/config",
				Format: "ghostty",
				Values: map[string]any{"theme": "light"},
			},
		},
	}

	result := MergeConfigs(base, override)

	if len(result.Files) != 1 {
		t.Fatalf("Expected 1 file (same target), got %d", len(result.Files))
	}

	// Override values should win
	if result.Files[0].Values["theme"] != "light" {
		t.Errorf("theme = %v, want %q", result.Files[0].Values["theme"], "light")
	}
}

func TestFileConfigStruct(t *testing.T) {
	fc := FileConfig{
		Target: "${XDG_CONFIG_HOME}/app/config",
		Format: "ghostty",
		Schema: map[string]FieldSchema{
			"theme": {
				Type:    "string",
				Default: "dark",
			},
			"font-size": {
				Type:    "int",
				Default: 14,
			},
		},
		Values: map[string]any{
			"theme":     "light",
			"font-size": 16,
		},
	}

	if fc.Target != "${XDG_CONFIG_HOME}/app/config" {
		t.Errorf("Target = %q", fc.Target)
	}

	if fc.Format != "ghostty" {
		t.Errorf("Format = %q", fc.Format)
	}

	if len(fc.Schema) != 2 {
		t.Errorf("Schema len = %d, want 2", len(fc.Schema))
	}

	if fc.Schema["theme"].Type != "string" {
		t.Errorf("Schema[theme].Type = %q", fc.Schema["theme"].Type)
	}

	if fc.Schema["font-size"].Default != 14 {
		t.Errorf("Schema[font-size].Default = %v", fc.Schema["font-size"].Default)
	}
}

func TestFieldSchemaStruct(t *testing.T) {
	fs := FieldSchema{
		Type:        "list",
		Default:     []string{"a", "b"},
		Items:       "string",
		Description: "A list of strings",
	}

	if fs.Type != "list" {
		t.Errorf("Type = %q", fs.Type)
	}

	if fs.Items != "string" {
		t.Errorf("Items = %q", fs.Items)
	}

	if fs.Description != "A list of strings" {
		t.Errorf("Description = %q", fs.Description)
	}
}

func TestCoalesce(t *testing.T) {
	tests := []struct {
		a, b, expected string
	}{
		{"first", "second", "first"},
		{"", "second", "second"},
		{"first", "", "first"},
		{"", "", ""},
	}

	for _, tt := range tests {
		result := coalesce(tt.a, tt.b)
		if result != tt.expected {
			t.Errorf("coalesce(%q, %q) = %q, want %q", tt.a, tt.b, result, tt.expected)
		}
	}
}

func TestMergeMaps(t *testing.T) {
	a := map[string]string{"key1": "a1", "key2": "a2"}
	b := map[string]string{"key2": "b2", "key3": "b3"}

	result := mergeMaps(a, b)

	if result["key1"] != "a1" {
		t.Errorf("key1 = %q, want %q", result["key1"], "a1")
	}

	// b should win for key2
	if result["key2"] != "b2" {
		t.Errorf("key2 = %q, want %q (b should win)", result["key2"], "b2")
	}

	if result["key3"] != "b3" {
		t.Errorf("key3 = %q, want %q", result["key3"], "b3")
	}
}

func TestMergeMapsNil(t *testing.T) {
	result := mergeMaps(nil, nil)
	if result == nil {
		t.Error("mergeMaps should return empty map, not nil")
	}
	if len(result) != 0 {
		t.Errorf("Expected empty map, got %v", result)
	}
}
