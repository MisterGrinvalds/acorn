package component

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Discovery handles component discovery operations.
type Discovery struct {
	dotfilesRoot string
}

// NewDiscovery creates a new Discovery instance.
func NewDiscovery(dotfilesRoot string) *Discovery {
	return &Discovery{dotfilesRoot: dotfilesRoot}
}

// DiscoverAll finds all valid components in the components directory.
func (d *Discovery) DiscoverAll() ([]*Component, error) {
	componentsDir := filepath.Join(d.dotfilesRoot, "components")

	entries, err := os.ReadDir(componentsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read components directory: %w", err)
	}

	var components []*Component

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()

		// Skip directories starting with _ (templates, disabled)
		if strings.HasPrefix(name, "_") {
			continue
		}

		compPath := filepath.Join(componentsDir, name)
		yamlPath := filepath.Join(compPath, "component.yaml")

		// Skip if no component.yaml
		if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
			continue
		}

		comp, err := Load(compPath)
		if err != nil {
			// Return error with context but continue discovery
			fmt.Fprintf(os.Stderr, "Warning: failed to load component %s: %v\n", name, err)
			continue
		}

		components = append(components, comp)
	}

	// Sort by name for consistent output
	sort.Slice(components, func(i, j int) bool {
		return components[i].Name < components[j].Name
	})

	return components, nil
}

// FindByName finds a component by its name.
func (d *Discovery) FindByName(name string) (*Component, error) {
	compPath := filepath.Join(d.dotfilesRoot, "components", name)

	if _, err := os.Stat(compPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("component not found: %s", name)
	}

	return Load(compPath)
}

// FindByCategory returns all components in a specific category.
func (d *Discovery) FindByCategory(category string) ([]*Component, error) {
	all, err := d.DiscoverAll()
	if err != nil {
		return nil, err
	}

	var filtered []*Component
	for _, comp := range all {
		if comp.Category == category {
			filtered = append(filtered, comp)
		}
	}

	return filtered, nil
}

// Categories returns all unique categories.
func (d *Discovery) Categories() ([]string, error) {
	all, err := d.DiscoverAll()
	if err != nil {
		return nil, err
	}

	categorySet := make(map[string]struct{})
	for _, comp := range all {
		if comp.Category != "" {
			categorySet[comp.Category] = struct{}{}
		}
	}

	var categories []string
	for cat := range categorySet {
		categories = append(categories, cat)
	}

	sort.Strings(categories)
	return categories, nil
}
