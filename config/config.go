// Package config provides component configuration file loading from .sapling/config.
// This package reads config.yaml files from .sapling/config at runtime and supports
// template rendering for dynamic configuration.
package config

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"
)

// saplingConfigDir returns the path to .sapling/config directory.
// This is relative to the current working directory or can be configured
// via SAPLING_DIR environment variable.
// It searches upward from the current directory to find .sapling/config.
func saplingConfigDir() (string, error) {
	if dir := os.Getenv("SAPLING_DIR"); dir != "" {
		return filepath.Join(dir, "config"), nil
	}

	// Search for .sapling/config starting from current directory and walking up
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	dir := cwd
	for {
		candidate := filepath.Join(dir, ".sapling", "config")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root
			break
		}
		dir = parent
	}

	// Fallback to .sapling/config in current directory (even if it doesn't exist)
	return filepath.Join(cwd, ".sapling", "config"), nil
}

// GetConfig returns the config.yaml for a component from .sapling/config.
// Returns the raw YAML bytes or an error if not found.
func GetConfig(component string) ([]byte, error) {
	configDir, err := saplingConfigDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(configDir, component, "config.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("no config for %s at %s: %w", component, path, err)
	}
	return data, nil
}

// GetConfigWithTemplate returns the config.yaml for a component with template rendering.
// The config file is treated as a Go template and rendered with the provided data.
func GetConfigWithTemplate(component string, templateData map[string]any) ([]byte, error) {
	configDir, err := saplingConfigDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(configDir, component, "config.yaml")
	templateContent, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("no config for %s at %s: %w", component, path, err)
	}

	// Parse and execute template
	tmpl, err := template.New(component).Parse(string(templateContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template for %s: %w", component, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateData); err != nil {
		return nil, fmt.Errorf("failed to execute template for %s: %w", component, err)
	}

	return buf.Bytes(), nil
}

// ListComponents returns all component names that have configs in .sapling/config.
func ListComponents() ([]string, error) {
	configDir, err := saplingConfigDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(configDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read config directory %s: %w", configDir, err)
	}

	var components []string
	for _, entry := range entries {
		if entry.IsDir() {
			// Check if it has a config.yaml
			configPath := filepath.Join(configDir, entry.Name(), "config.yaml")
			if _, err := os.Stat(configPath); err == nil {
				components = append(components, entry.Name())
			}
		}
	}

	return components, nil
}

// HasConfig checks if a component has a config in .sapling/config.
func HasConfig(component string) bool {
	_, err := GetConfig(component)
	return err == nil
}

// WalkConfigs walks all config.yaml files in .sapling/config.
// For each config file found, it calls the provided function with the component name.
func WalkConfigs(fn func(component string, path string) error) error {
	configDir, err := saplingConfigDir()
	if err != nil {
		return err
	}

	return filepath.WalkDir(configDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if d.Name() == "config.yaml" {
			// Extract component name from path
			relPath, err := filepath.Rel(configDir, filepath.Dir(path))
			if err != nil {
				return err
			}

			return fn(relPath, path)
		}

		return nil
	})
}

// SaplingRoot returns the root .sapling directory path.
// This is the directory containing both config/ and generated/.
func SaplingRoot() (string, error) {
	if dir := os.Getenv("SAPLING_DIR"); dir != "" {
		return dir, nil
	}

	// Search for .sapling starting from current directory and walking up
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	dir := cwd
	for {
		candidate := filepath.Join(dir, ".sapling")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root
			break
		}
		dir = parent
	}

	// Fallback to .sapling in current directory (even if it doesn't exist)
	return filepath.Join(cwd, ".sapling"), nil
}

// GeneratedDir returns the path to .sapling/generated directory.
// This is where Acorn writes generated config files.
func GeneratedDir() (string, error) {
	root, err := SaplingRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "generated"), nil
}

// IsValidSaplingRepo checks if the .sapling directory is a real sapling repository.
// A valid sapling repo must have config/, .git, or ai/ directory.
// Just having generated/ is not enough (that gets auto-created).
func IsValidSaplingRepo() bool {
	root, err := SaplingRoot()
	if err != nil {
		return false
	}

	// Check for config/ directory (primary indicator)
	if _, err := os.Stat(filepath.Join(root, "config")); err == nil {
		return true
	}

	// Check for .git directory
	if _, err := os.Stat(filepath.Join(root, ".git")); err == nil {
		return true
	}

	// Check for ai/ directory
	if _, err := os.Stat(filepath.Join(root, "ai")); err == nil {
		return true
	}

	return false
}

// EnsureGeneratedDir creates the .sapling/generated directory if it doesn't exist.
// Returns an error if no valid sapling repository is found.
func EnsureGeneratedDir() error {
	if !IsValidSaplingRepo() {
		return fmt.Errorf("no valid .sapling repository found. Run 'acorn setup' to configure one")
	}

	genDir, err := GeneratedDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(genDir, 0o755)
}
