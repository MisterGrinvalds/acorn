// Package component provides types and operations for dotfiles component management.
package component

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/mistergrinvalds/acorn/internal/utils/config"
	"gopkg.in/yaml.v3"
)

// TestCategory represents a category of tests.
type TestCategory string

const (
	TestCategoryConfig  TestCategory = "config"
	TestCategoryShell   TestCategory = "shell"
	TestCategoryInstall TestCategory = "install"
	TestCategoryFiles   TestCategory = "files"
)

// TestStatus represents the status of a test.
type TestStatus string

const (
	TestStatusPassed  TestStatus = "passed"
	TestStatusFailed  TestStatus = "failed"
	TestStatusSkipped TestStatus = "skipped"
)

// TestCase represents an individual test case.
type TestCase struct {
	Name     string       `json:"name" yaml:"name"`
	Category TestCategory `json:"category" yaml:"category"`
	Status   TestStatus   `json:"status" yaml:"status"`
	Message  string       `json:"message,omitempty" yaml:"message,omitempty"`
}

// TestResult represents the result of testing a component.
type TestResult struct {
	Component string     `json:"component" yaml:"component"`
	Path      string     `json:"path" yaml:"path"`
	Tests     []TestCase `json:"tests" yaml:"tests"`
	Passed    int        `json:"passed" yaml:"passed"`
	Failed    int        `json:"failed" yaml:"failed"`
	Skipped   int        `json:"skipped" yaml:"skipped"`
}

// Summary returns a summary string.
func (tr *TestResult) Summary() string {
	return fmt.Sprintf("%d passed, %d failed, %d skipped", tr.Passed, tr.Failed, tr.Skipped)
}

// HasFailures returns true if any tests failed.
func (tr *TestResult) HasFailures() bool {
	return tr.Failed > 0
}

// Tester handles component testing operations.
type Tester struct {
	dotfilesRoot string
	skipMissing  bool
	verbose      bool
}

// TesterOption is a functional option for Tester.
type TesterOption func(*Tester)

// WithSkipMissing configures the tester to skip missing components.
func WithSkipMissing(skip bool) TesterOption {
	return func(t *Tester) {
		t.skipMissing = skip
	}
}

// WithVerbose configures verbose output.
func WithVerbose(verbose bool) TesterOption {
	return func(t *Tester) {
		t.verbose = verbose
	}
}

// NewTester creates a new Tester instance.
func NewTester(dotfilesRoot string, opts ...TesterOption) *Tester {
	t := &Tester{
		dotfilesRoot: dotfilesRoot,
		skipMissing:  false,
		verbose:      false,
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// TestComponent tests a single component by name.
func (t *Tester) TestComponent(name string) (*TestResult, error) {
	configPath := filepath.Join(t.dotfilesRoot, ".sapling", "config", name)

	// Check if component exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if t.skipMissing {
			return &TestResult{
				Component: name,
				Path:      configPath,
				Tests: []TestCase{{
					Name:     "component_exists",
					Category: TestCategoryConfig,
					Status:   TestStatusSkipped,
					Message:  "component not found",
				}},
				Skipped: 1,
			}, nil
		}
		return nil, fmt.Errorf("component not found: %s", name)
	}

	result := &TestResult{
		Component: name,
		Path:      configPath,
		Tests:     []TestCase{},
	}

	// Load and test config
	cfg, err := t.loadConfig(configPath)
	if err != nil {
		result.addTest("config_valid", TestCategoryConfig, TestStatusFailed, err.Error())
		return result, nil
	}
	result.addTest("config_valid", TestCategoryConfig, TestStatusPassed, "")

	// Run all test categories
	t.testConfigFields(result, cfg)
	t.testShellSyntax(result, configPath)
	t.testInstallConfig(result, cfg)
	t.testGeneratedFiles(result, name)

	return result, nil
}

// TestAll tests all components in the .sapling/config directory.
func (t *Tester) TestAll() ([]*TestResult, error) {
	configDir := filepath.Join(t.dotfilesRoot, ".sapling", "config")

	entries, err := os.ReadDir(configDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read config directory: %w", err)
	}

	var results []*TestResult
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		result, err := t.TestComponent(name)
		if err != nil {
			// Create a failed result for this component
			result = &TestResult{
				Component: name,
				Path:      filepath.Join(configDir, name),
				Tests: []TestCase{{
					Name:     "load_error",
					Category: TestCategoryConfig,
					Status:   TestStatusFailed,
					Message:  err.Error(),
				}},
				Failed: 1,
			}
		}
		results = append(results, result)
	}

	return results, nil
}

// TestByCategory tests components in a specific category.
func (t *Tester) TestByCategory(category string) ([]*TestResult, error) {
	configDir := filepath.Join(t.dotfilesRoot, ".sapling", "config")

	entries, err := os.ReadDir(configDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read config directory: %w", err)
	}

	var results []*TestResult
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()

		// Match category by name
		if !matchesCategory(name, category) {
			continue
		}

		result, err := t.TestComponent(name)
		if err != nil {
			continue
		}
		results = append(results, result)
	}

	return results, nil
}

// loadConfig loads a component's config.yaml file.
func (t *Tester) loadConfig(configPath string) (*config.BaseConfig, error) {
	yamlPath := filepath.Join(configPath, "config.yaml")

	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config.yaml: %w", err)
	}

	var cfg config.BaseConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config.yaml: %w", err)
	}

	return &cfg, nil
}

// testConfigFields tests that required config fields are present.
func (t *Tester) testConfigFields(result *TestResult, cfg *config.BaseConfig) {
	// Required fields
	if cfg.Name == "" {
		result.addTest("field_name", TestCategoryConfig, TestStatusFailed, "missing required field: name")
	} else {
		result.addTest("field_name", TestCategoryConfig, TestStatusPassed, "")
	}

	if cfg.Description == "" {
		result.addTest("field_description", TestCategoryConfig, TestStatusFailed, "missing required field: description")
	} else {
		result.addTest("field_description", TestCategoryConfig, TestStatusPassed, "")
	}

	if cfg.Version == "" {
		result.addTest("field_version", TestCategoryConfig, TestStatusFailed, "missing required field: version")
	} else {
		result.addTest("field_version", TestCategoryConfig, TestStatusPassed, "")
	}

	// Check for at least one feature
	hasFeatures := len(cfg.Env) > 0 ||
		len(cfg.Aliases) > 0 ||
		len(cfg.ShellFunctions) > 0 ||
		len(cfg.Wrappers) > 0 ||
		len(cfg.Files) > 0 ||
		len(cfg.Paths) > 0

	if hasFeatures {
		result.addTest("has_features", TestCategoryConfig, TestStatusPassed, "")
	} else {
		result.addTest("has_features", TestCategoryConfig, TestStatusFailed, "component has no features (env, aliases, shell_functions, wrappers, files, or paths)")
	}
}

// testShellSyntax tests shell script syntax for aliases and functions.
func (t *Tester) testShellSyntax(result *TestResult, configPath string) {
	// Check for shell files in the component directory
	shellPatterns := []string{"*.sh", "**/*.sh"}

	hasShellFiles := false
	for _, pattern := range shellPatterns {
		matches, err := filepath.Glob(filepath.Join(configPath, pattern))
		if err == nil && len(matches) > 0 {
			hasShellFiles = true
			for _, match := range matches {
				if err := checkBashSyntax(match); err != nil {
					result.addTest("shell_syntax_"+filepath.Base(match), TestCategoryShell, TestStatusFailed, err.Error())
				} else {
					result.addTest("shell_syntax_"+filepath.Base(match), TestCategoryShell, TestStatusPassed, "")
				}
			}
		}
	}

	if !hasShellFiles {
		result.addTest("shell_files", TestCategoryShell, TestStatusSkipped, "no shell files found")
	}
}

// testInstallConfig tests that install configuration is valid.
func (t *Tester) testInstallConfig(result *TestResult, cfg *config.BaseConfig) {
	if !cfg.HasInstall() {
		result.addTest("install_config", TestCategoryInstall, TestStatusSkipped, "no install configuration")
		return
	}

	for _, tool := range cfg.Install.Tools {
		// Check that tool has a check command
		if tool.Check == "" {
			result.addTest("install_check_"+tool.Name, TestCategoryInstall, TestStatusFailed,
				fmt.Sprintf("tool %s missing check command", tool.Name))
		} else {
			result.addTest("install_check_"+tool.Name, TestCategoryInstall, TestStatusPassed, "")
		}

		// Check that tool has at least one install method
		if len(tool.Methods) == 0 {
			result.addTest("install_methods_"+tool.Name, TestCategoryInstall, TestStatusFailed,
				fmt.Sprintf("tool %s has no install methods", tool.Name))
		} else {
			result.addTest("install_methods_"+tool.Name, TestCategoryInstall, TestStatusPassed, "")
		}
	}
}

// testGeneratedFiles tests that generated files exist at expected locations.
func (t *Tester) testGeneratedFiles(result *TestResult, name string) {
	generatedDir := filepath.Join(t.dotfilesRoot, ".sapling", "generated", name)

	// Check if generated directory exists
	if _, err := os.Stat(generatedDir); os.IsNotExist(err) {
		result.addTest("generated_dir", TestCategoryFiles, TestStatusSkipped, "no generated files directory")
		return
	}

	// Check for expected generated files
	expectedFiles := []string{"env.sh", "aliases.sh"}
	for _, file := range expectedFiles {
		filePath := filepath.Join(generatedDir, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			result.addTest("generated_"+file, TestCategoryFiles, TestStatusSkipped, "file not generated")
		} else {
			result.addTest("generated_"+file, TestCategoryFiles, TestStatusPassed, "")
		}
	}
}

// addTest adds a test case to the result.
func (result *TestResult) addTest(name string, category TestCategory, status TestStatus, message string) {
	result.Tests = append(result.Tests, TestCase{
		Name:     name,
		Category: category,
		Status:   status,
		Message:  message,
	})

	switch status {
	case TestStatusPassed:
		result.Passed++
	case TestStatusFailed:
		result.Failed++
	case TestStatusSkipped:
		result.Skipped++
	}
}

// checkBashSyntax checks bash syntax of a shell file.
func checkBashSyntax(path string) error {
	cmd := exec.Command("bash", "-n", path)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("syntax error: %s", strings.TrimSpace(string(output)))
	}
	return nil
}

// matchesCategory checks if a component matches a category.
func matchesCategory(name string, category string) bool {
	// Category mappings
	categories := map[string][]string{
		"ai":          {"claude", "huggingface", "ollama"},
		"cloud":       {"cloudflare", "aws", "azure", "digitalocean"},
		"data":        {"database", "secrets", "posting", "datagrip", "infisical", "openmetadata", "postman"},
		"devops":      {"kubernetes", "docker", "lazydocker", "helm", "kubectl", "k9s", "argocd", "n8n"},
		"ide":         {"neovim", "vscode", "goland", "intellij"},
		"programming": {"python", "node", "go", "uv", "nvm", "pnpm"},
		"terminal":    {"fzf", "ghostty", "tmux", "jq", "yq"},
		"vcs":         {"git", "github"},
		"sysadm":      {"btop"},
		"identity":    {"keycloak"},
		"artifacts":   {"jfrog"},
	}

	if components, ok := categories[category]; ok {
		return slices.Contains(components, name)
	}

	return false
}

// GetAllCategories returns all available test categories.
func GetAllCategories() []string {
	return []string{
		"ai", "cloud", "data", "devops", "ide",
		"programming", "terminal", "vcs", "sysadm",
		"identity", "artifacts",
	}
}
