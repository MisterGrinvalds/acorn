package componentconfig

import (
	"embed"
	"fmt"
	"io/fs"
)

// ConfigFS contains the embedded default configuration files.
// The config/ directory contains YAML configs for each component.
//
//go:embed config
var ConfigFS embed.FS

// GetEmbeddedConfig returns the embedded config.yaml for a component.
// Returns the raw YAML bytes or an error if not found.
func GetEmbeddedConfig(component string) ([]byte, error) {
	path := fmt.Sprintf("config/%s/config.yaml", component)
	data, err := ConfigFS.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("no embedded config for %s: %w", component, err)
	}
	return data, nil
}

// ListEmbeddedComponents returns all component names that have embedded configs.
func ListEmbeddedComponents() ([]string, error) {
	entries, err := fs.ReadDir(ConfigFS, "config")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded config directory: %w", err)
	}

	var components []string
	for _, entry := range entries {
		if entry.IsDir() {
			// Check if it has a config.yaml
			configPath := fmt.Sprintf("config/%s/config.yaml", entry.Name())
			if _, err := ConfigFS.ReadFile(configPath); err == nil {
				components = append(components, entry.Name())
			}
		}
	}

	return components, nil
}

// HasEmbeddedConfig checks if a component has an embedded config.
func HasEmbeddedConfig(component string) bool {
	_, err := GetEmbeddedConfig(component)
	return err == nil
}
