package configfile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mistergrinvalds/acorn/internal/utils/config"
)

func TestExpandPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		envVars  map[string]string
		expected string
	}{
		{
			name:     "no expansion needed",
			input:    "/usr/local/bin",
			expected: "/usr/local/bin",
		},
		{
			name:  "expand XDG_CONFIG_HOME",
			input: "${XDG_CONFIG_HOME}/ghostty/config",
			envVars: map[string]string{
				"XDG_CONFIG_HOME": "/home/user/.config",
			},
			expected: "/home/user/.config/ghostty/config",
		},
		{
			name:  "expand HOME",
			input: "${HOME}/.config/app",
			envVars: map[string]string{
				"HOME": "/home/testuser",
			},
			expected: "/home/testuser/.config/app",
		},
		{
			name:  "expand multiple vars",
			input: "${HOME}/${APP_NAME}/config",
			envVars: map[string]string{
				"HOME":     "/home/user",
				"APP_NAME": "myapp",
			},
			expected: "/home/user/myapp/config",
		},
		{
			name:     "missing env var expands to empty",
			input:    "${NONEXISTENT_VAR}/path",
			expected: "/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set env vars
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			result := ExpandPath(tt.input)
			if result != tt.expected {
				t.Errorf("ExpandPath(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExpandPathWithTilde(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Could not get home directory")
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "expand tilde",
			input:    "~/config/app",
			expected: filepath.Join(home, "config/app"),
		},
		{
			name:     "tilde in middle not expanded",
			input:    "/path/~/config",
			expected: "/path/~/config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExpandPath(tt.input)
			if result != tt.expected {
				t.Errorf("ExpandPath(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRegistry(t *testing.T) {
	// Ghostty writer should be registered via init()
	writer, err := GetWriter("ghostty")
	if err != nil {
		t.Fatalf("GetWriter(ghostty) failed: %v", err)
	}

	if writer.Format() != "ghostty" {
		t.Errorf("writer.Format() = %q, want %q", writer.Format(), "ghostty")
	}
}

func TestGetWriterUnknownFormat(t *testing.T) {
	_, err := GetWriter("unknown_format")
	if err == nil {
		t.Error("GetWriter(unknown_format) should return error")
	}
}

func TestManagerGenerateFileDryRun(t *testing.T) {
	manager := NewManager(true) // dry run

	fc := config.FileConfig{
		Target: "/tmp/test-config-dry-run",
		Format: "ghostty",
		Values: map[string]interface{}{
			"theme":     "dark",
			"font-size": 14,
		},
	}

	result, err := manager.GenerateFile(fc)
	if err != nil {
		t.Fatalf("GenerateFile failed: %v", err)
	}

	if result.Written {
		t.Error("File should not be written in dry run mode")
	}

	// In legacy mode (no generatedDir), GeneratedPath and SymlinkTarget are the same
	if result.SymlinkTarget != "/tmp/test-config-dry-run" {
		t.Errorf("SymlinkTarget = %q, want %q", result.SymlinkTarget, "/tmp/test-config-dry-run")
	}

	if result.Format != "ghostty" {
		t.Errorf("Format = %q, want %q", result.Format, "ghostty")
	}

	if result.Content == "" {
		t.Error("Content should not be empty")
	}
}

func TestManagerGenerateFileActualWrite(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "configfile-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	targetPath := filepath.Join(tmpDir, "subdir", "config")

	manager := NewManager(false) // not dry run

	fc := config.FileConfig{
		Target: targetPath,
		Format: "ghostty",
		Values: map[string]interface{}{
			"theme":     "Catppuccin Mocha",
			"font-size": 14,
		},
	}

	result, err := manager.GenerateFile(fc)
	if err != nil {
		t.Fatalf("GenerateFile failed: %v", err)
	}

	if !result.Written {
		t.Error("File should be written")
	}

	// Verify file exists
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Verify content
	content, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	if string(content) != result.Content {
		t.Error("File content doesn't match result.Content")
	}
}

func TestManagerGenerateFiles(t *testing.T) {
	manager := NewManager(true) // dry run

	files := []config.FileConfig{
		{
			Target: "/tmp/config1",
			Format: "ghostty",
			Values: map[string]interface{}{"key1": "value1"},
		},
		{
			Target: "/tmp/config2",
			Format: "ghostty",
			Values: map[string]interface{}{"key2": "value2"},
		},
	}

	results, err := manager.GenerateFiles(files)
	if err != nil {
		t.Fatalf("GenerateFiles failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
}

func TestManagerGenerateFileUnknownFormat(t *testing.T) {
	manager := NewManager(true)

	fc := config.FileConfig{
		Target: "/tmp/test",
		Format: "unknown_format",
		Values: map[string]interface{}{"key": "value"},
	}

	_, err := manager.GenerateFile(fc)
	if err == nil {
		t.Error("GenerateFile should fail with unknown format")
	}
}
