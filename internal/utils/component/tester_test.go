package component

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestNewTester(t *testing.T) {
	tester := NewTester("/tmp/test")
	if tester == nil {
		t.Fatal("NewTester returned nil")
	}
	if tester.dotfilesRoot != "/tmp/test" {
		t.Errorf("expected dotfilesRoot /tmp/test, got %s", tester.dotfilesRoot)
	}
	if tester.skipMissing {
		t.Error("skipMissing should be false by default")
	}
	if tester.verbose {
		t.Error("verbose should be false by default")
	}
}

func TestNewTesterWithOptions(t *testing.T) {
	tester := NewTester("/tmp/test",
		WithSkipMissing(true),
		WithVerbose(true),
	)
	if !tester.skipMissing {
		t.Error("skipMissing should be true")
	}
	if !tester.verbose {
		t.Error("verbose should be true")
	}
}

func TestTestResult_Summary(t *testing.T) {
	result := &TestResult{
		Passed:  5,
		Failed:  2,
		Skipped: 1,
	}
	summary := result.Summary()
	expected := "5 passed, 2 failed, 1 skipped"
	if summary != expected {
		t.Errorf("expected %q, got %q", expected, summary)
	}
}

func TestTestResult_HasFailures(t *testing.T) {
	tests := []struct {
		name     string
		failed   int
		expected bool
	}{
		{"no failures", 0, false},
		{"has failures", 1, true},
		{"multiple failures", 5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &TestResult{Failed: tt.failed}
			if result.HasFailures() != tt.expected {
				t.Errorf("HasFailures() = %v, want %v", result.HasFailures(), tt.expected)
			}
		})
	}
}

func TestTestResult_addTest(t *testing.T) {
	result := &TestResult{}

	result.addTest("test1", TestCategoryConfig, TestStatusPassed, "")
	result.addTest("test2", TestCategoryShell, TestStatusFailed, "error message")
	result.addTest("test3", TestCategoryInstall, TestStatusSkipped, "not applicable")

	if len(result.Tests) != 3 {
		t.Errorf("expected 3 tests, got %d", len(result.Tests))
	}
	if result.Passed != 1 {
		t.Errorf("expected 1 passed, got %d", result.Passed)
	}
	if result.Failed != 1 {
		t.Errorf("expected 1 failed, got %d", result.Failed)
	}
	if result.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", result.Skipped)
	}
}

func TestMatchesCategory(t *testing.T) {
	tests := []struct {
		name     string
		category string
		expected bool
	}{
		{"python", "programming", true},
		{"node", "programming", true},
		{"go", "programming", true},
		{"claude", "ai", true},
		{"ollama", "ai", true},
		{"git", "vcs", true},
		{"github", "vcs", true},
		{"docker", "devops", true},
		{"kubernetes", "devops", true},
		{"python", "ai", false},
		{"unknown", "programming", false},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_"+tt.category, func(t *testing.T) {
			result := matchesCategory(tt.name, tt.category)
			if result != tt.expected {
				t.Errorf("matchesCategory(%q, %q) = %v, want %v",
					tt.name, tt.category, result, tt.expected)
			}
		})
	}
}

func TestGetAllCategories(t *testing.T) {
	categories := GetAllCategories()
	if len(categories) == 0 {
		t.Error("GetAllCategories returned empty slice")
	}

	// Check for expected categories
	expected := []string{"ai", "cloud", "programming", "devops", "vcs"}
	for _, exp := range expected {
		if !slices.Contains(categories, exp) {
			t.Errorf("expected category %q not found", exp)
		}
	}
}

func TestTester_TestComponent_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	os.MkdirAll(filepath.Join(tmpDir, ".sapling", "config"), 0755)

	tester := NewTester(tmpDir)
	_, err := tester.TestComponent("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent component")
	}
}

func TestTester_TestComponent_NotFound_SkipMissing(t *testing.T) {
	tmpDir := t.TempDir()
	os.MkdirAll(filepath.Join(tmpDir, ".sapling", "config"), 0755)

	tester := NewTester(tmpDir, WithSkipMissing(true))
	result, err := tester.TestComponent("nonexistent")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", result.Skipped)
	}
}

func TestTester_TestComponent_ValidConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".sapling", "config", "testcomp")
	os.MkdirAll(configDir, 0755)

	// Create a valid config.yaml
	configContent := `name: testcomp
description: Test component
version: 1.0.0

aliases:
  test: "echo test"
`
	os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte(configContent), 0644)

	tester := NewTester(tmpDir)
	result, err := tester.TestComponent("testcomp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Component != "testcomp" {
		t.Errorf("expected component testcomp, got %s", result.Component)
	}

	// Should have passed config tests
	if result.Passed == 0 {
		t.Error("expected at least one passed test")
	}
}

func TestTester_TestComponent_InvalidConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".sapling", "config", "badcomp")
	os.MkdirAll(configDir, 0755)

	// Create an invalid config.yaml
	configContent := `name:
description:
version:
`
	os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte(configContent), 0644)

	tester := NewTester(tmpDir)
	result, err := tester.TestComponent("badcomp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have failed tests for missing required fields
	if result.Failed == 0 {
		t.Error("expected at least one failed test for empty required fields")
	}
}

func TestTester_TestAll(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".sapling", "config")
	os.MkdirAll(configDir, 0755)

	// Create two valid components
	for _, name := range []string{"comp1", "comp2"} {
		compDir := filepath.Join(configDir, name)
		os.MkdirAll(compDir, 0755)
		configContent := `name: ` + name + `
description: Test component ` + name + `
version: 1.0.0

aliases:
  test: "echo test"
`
		os.WriteFile(filepath.Join(compDir, "config.yaml"), []byte(configContent), 0644)
	}

	tester := NewTester(tmpDir)
	results, err := tester.TestAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestTester_TestByCategory(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".sapling", "config")
	os.MkdirAll(configDir, 0755)

	// Create components in different categories
	components := map[string]string{
		"python":      "programming",
		"go":          "programming",
		"docker":      "devops",
		"kubernetes":  "devops",
		"git":         "vcs",
	}

	for name := range components {
		compDir := filepath.Join(configDir, name)
		os.MkdirAll(compDir, 0755)
		configContent := `name: ` + name + `
description: Test component ` + name + `
version: 1.0.0

aliases:
  test: "echo test"
`
		os.WriteFile(filepath.Join(compDir, "config.yaml"), []byte(configContent), 0644)
	}

	tester := NewTester(tmpDir)

	// Test programming category
	results, err := tester.TestByCategory("programming")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 programming components, got %d", len(results))
	}

	// Test devops category
	results, err = tester.TestByCategory("devops")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 devops components, got %d", len(results))
	}

	// Test vcs category
	results, err = tester.TestByCategory("vcs")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 vcs component, got %d", len(results))
	}
}
