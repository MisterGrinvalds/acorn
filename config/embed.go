// Package config provides embedded component configuration files.
// This package embeds all component config.yaml files from subdirectories.
package config

import (
	"embed"
	"fmt"
	"io/fs"
)

// FS contains the embedded configuration files.
// Each component has its own subdirectory with a config.yaml file.
//
//go:embed */config.yaml
var FS embed.FS

// GetConfig returns the embedded config.yaml for a component.
// Returns the raw YAML bytes or an error if not found.
func GetConfig(component string) ([]byte, error) {
	path := fmt.Sprintf("%s/config.yaml", component)
	data, err := FS.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("no embedded config for %s: %w", component, err)
	}
	return data, nil
}

// ListComponents returns all component names that have embedded configs.
func ListComponents() ([]string, error) {
	entries, err := fs.ReadDir(FS, ".")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded config directory: %w", err)
	}

	var components []string
	for _, entry := range entries {
		if entry.IsDir() {
			// Check if it has a config.yaml
			configPath := fmt.Sprintf("%s/config.yaml", entry.Name())
			if _, err := FS.ReadFile(configPath); err == nil {
				components = append(components, entry.Name())
			}
		}
	}

	return components, nil
}

// HasConfig checks if a component has an embedded config.
func HasConfig(component string) bool {
	_, err := GetConfig(component)
	return err == nil
}
