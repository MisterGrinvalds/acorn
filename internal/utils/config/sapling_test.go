package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetComponentConfig(t *testing.T) {
	// Test loading a component config
	data, err := GetComponentConfig("git")
	if err != nil {
		t.Fatalf("Failed to load git config: %v", err)
	}

	if len(data) == 0 {
		t.Error("Expected non-empty config data")
	}
}

func TestListComponentConfigs(t *testing.T) {
	components, err := ListComponentConfigs()
	if err != nil {
		t.Fatalf("Failed to list components: %v", err)
	}

	if len(components) == 0 {
		t.Error("Expected at least one component")
	}

	t.Logf("Found %d components: %v", len(components), components)
}

func TestGetComponentConfigWithTemplate(t *testing.T) {
	// Create a temporary test config with template
	tempDir := t.TempDir()
	testComponent := "test-template"
	// SAPLING_DIR should point to the directory containing the config/ directory
	saplingDir := filepath.Join(tempDir, ".sapling")
	configDir := filepath.Join(saplingDir, "config")
	componentDir := filepath.Join(configDir, testComponent)

	if err := os.MkdirAll(componentDir, 0o755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Write a test config with template variables
	testConfig := `name: {{ .Name }}
version: {{ .Version }}
description: Test component for {{ .Name }}`

	configPath := filepath.Join(componentDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(testConfig), 0o644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Set SAPLING_DIR to use our test directory (it will append /config)
	oldSaplingDir := os.Getenv("SAPLING_DIR")
	os.Setenv("SAPLING_DIR", saplingDir)
	defer os.Setenv("SAPLING_DIR", oldSaplingDir)

	// Test template rendering
	templateData := map[string]any{
		"Name":    "TestApp",
		"Version": "1.0.0",
	}

	rendered, err := GetComponentConfigWithTemplate(testComponent, templateData)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	expected := `name: TestApp
version: 1.0.0
description: Test component for TestApp`

	if string(rendered) != expected {
		t.Errorf("Template rendering mismatch.\nExpected:\n%s\n\nGot:\n%s", expected, string(rendered))
	}
}

func TestHasComponentConfig(t *testing.T) {
	// Test with existing component
	if !HasComponentConfig("git") {
		t.Error("Expected git component to exist")
	}

	// Test with non-existing component
	if HasComponentConfig("nonexistent-component-xyz") {
		t.Error("Expected nonexistent component to not exist")
	}
}

func TestSaplingRoot(t *testing.T) {
	// Should find .sapling in parent directory
	root, err := SaplingRoot()
	if err != nil {
		t.Fatalf("Failed to get sapling root: %v", err)
	}

	if !filepath.IsAbs(root) {
		t.Error("Expected absolute path")
	}

	// Should contain "sapling" in the path
	if !strings.Contains(root, "sapling") {
		t.Errorf("Expected path to contain 'sapling', got: %s", root)
	}
}

func TestGeneratedDir(t *testing.T) {
	genDir, err := GeneratedDir()
	if err != nil {
		t.Fatalf("Failed to get generated dir: %v", err)
	}

	if !filepath.IsAbs(genDir) {
		t.Error("Expected absolute path")
	}

	// Should end with generated
	if !strings.HasSuffix(genDir, "generated") {
		t.Errorf("Expected path to end with 'generated', got: %s", genDir)
	}
}

func TestEnsureGeneratedDir(t *testing.T) {
	// Create temp test environment
	tempDir := t.TempDir()
	saplingDir := filepath.Join(tempDir, ".sapling")

	// Create config directory to make it a valid sapling repo
	configDir := filepath.Join(saplingDir, "config")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	oldSaplingDir := os.Getenv("SAPLING_DIR")
	os.Setenv("SAPLING_DIR", saplingDir)
	defer os.Setenv("SAPLING_DIR", oldSaplingDir)

	// Generated dir should not exist yet
	genDir := filepath.Join(saplingDir, "generated")
	if _, err := os.Stat(genDir); !os.IsNotExist(err) {
		t.Error("Generated dir should not exist yet")
	}

	// Create it
	if err := EnsureGeneratedDir(); err != nil {
		t.Fatalf("Failed to ensure generated dir: %v", err)
	}

	// Should exist now
	if info, err := os.Stat(genDir); err != nil || !info.IsDir() {
		t.Error("Generated dir should exist and be a directory")
	}
}
