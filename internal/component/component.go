// Package component provides types and operations for dotfiles component management.
package component

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Component represents a dotfiles component with its metadata.
type Component struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Version     string   `yaml:"version"`
	Category    string   `yaml:"category"`
	Platforms   []string `yaml:"platforms"`
	Shells      []string `yaml:"shells"`
	Requires    Requires `yaml:"requires"`
	Provides    Provides `yaml:"provides"`
	XDG         XDG      `yaml:"xdg"`
	Config      Config   `yaml:"config"`

	// Metadata (not from YAML)
	Path     string `yaml:"-"` // Absolute path to component directory
	YAMLPath string `yaml:"-"` // Path to component.yaml file
}

// Requires defines component dependencies.
type Requires struct {
	Tools      []string `yaml:"tools"`
	Components []string `yaml:"components"`
}

// Provides defines what a component provides.
type Provides struct {
	Aliases     []string `yaml:"aliases"`
	Functions   []string `yaml:"functions"`
	Completions []string `yaml:"completions"`
}

// XDG defines XDG directory usage.
type XDG struct {
	Config string `yaml:"config"`
	Data   string `yaml:"data"`
	Cache  string `yaml:"cache"`
	State  string `yaml:"state"`
}

// Config defines configuration file management.
type Config struct {
	Files       []ConfigFile      `yaml:"files"`
	Directories []ConfigDirectory `yaml:"directories"`
}

// ConfigFile represents a file to be linked or copied.
type ConfigFile struct {
	Source      string `yaml:"source"`
	Target      string `yaml:"target"`
	Method      string `yaml:"method"`
	Platform    string `yaml:"platform"`
	Permissions string `yaml:"permissions"`
}

// ConfigDirectory represents a directory to create.
type ConfigDirectory struct {
	Target      string `yaml:"target"`
	Permissions string `yaml:"permissions"`
}

// Load loads a component from its component.yaml file.
func Load(path string) (*Component, error) {
	yamlPath := filepath.Join(path, "component.yaml")

	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read component.yaml: %w", err)
	}

	var comp Component
	if err := yaml.Unmarshal(data, &comp); err != nil {
		return nil, fmt.Errorf("failed to parse component.yaml: %w", err)
	}

	comp.Path = path
	comp.YAMLPath = yamlPath

	return &comp, nil
}

// ShellFiles returns paths to all shell files in the component.
func (c *Component) ShellFiles() []string {
	var files []string
	patterns := []string{"env.sh", "aliases.sh", "functions.sh", "completions.sh"}

	for _, pattern := range patterns {
		path := filepath.Join(c.Path, pattern)
		if _, err := os.Stat(path); err == nil {
			files = append(files, path)
		}
	}

	return files
}

// ConfigSourceFiles returns absolute paths to config source files.
func (c *Component) ConfigSourceFiles() []string {
	var files []string
	for _, cfg := range c.Config.Files {
		sourcePath := filepath.Join(c.Path, cfg.Source)
		files = append(files, sourcePath)
	}
	return files
}

// SupportsCurrentPlatform checks if component supports the current platform.
func (c *Component) SupportsCurrentPlatform(platform string) bool {
	// Empty means all platforms
	if len(c.Platforms) == 0 {
		return true
	}

	for _, p := range c.Platforms {
		if p == platform {
			return true
		}
	}
	return false
}

// SupportsCurrentShell checks if component supports the current shell.
func (c *Component) SupportsCurrentShell(shell string) bool {
	// Empty means all shells
	if len(c.Shells) == 0 {
		return true
	}

	for _, s := range c.Shells {
		if s == shell {
			return true
		}
	}
	return false
}
