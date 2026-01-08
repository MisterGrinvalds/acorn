// Package configfile provides config file generation in various formats.
package configfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mistergrinvalds/acorn/internal/componentconfig"
)

// Writer defines the interface for config file format writers.
type Writer interface {
	// Format returns the format identifier (e.g., "ghostty", "json", "yaml")
	Format() string

	// Write generates the config file content from values
	Write(values map[string]interface{}) ([]byte, error)
}

// Registry holds all registered format writers.
var Registry = make(map[string]Writer)

// Register adds a writer to the registry.
func Register(w Writer) {
	Registry[w.Format()] = w
}

// GetWriter returns a writer for the given format.
func GetWriter(format string) (Writer, error) {
	w, ok := Registry[format]
	if !ok {
		formats := make([]string, 0, len(Registry))
		for f := range Registry {
			formats = append(formats, f)
		}
		return nil, fmt.Errorf("unknown format %q (available: %s)", format, strings.Join(formats, ", "))
	}
	return w, nil
}

// GeneratedFile represents a generated config file with metadata.
type GeneratedFile struct {
	Target  string `json:"target" yaml:"target"`
	Format  string `json:"format" yaml:"format"`
	Content string `json:"content" yaml:"content"`
	Written bool   `json:"written" yaml:"written"`
}

// Manager handles config file generation for components.
type Manager struct {
	dryRun bool
}

// NewManager creates a new config file manager.
func NewManager(dryRun bool) *Manager {
	return &Manager{dryRun: dryRun}
}

// ExpandPath expands environment variables in a path.
func ExpandPath(path string) string {
	// Handle ${VAR} syntax
	result := os.Expand(path, func(key string) string {
		return os.Getenv(key)
	})

	// Handle ~ for home directory
	if strings.HasPrefix(result, "~/") {
		home, _ := os.UserHomeDir()
		result = filepath.Join(home, result[2:])
	}

	return result
}

// GenerateFile generates a single config file.
func (m *Manager) GenerateFile(fc componentconfig.FileConfig) (*GeneratedFile, error) {
	writer, err := GetWriter(fc.Format)
	if err != nil {
		return nil, err
	}

	content, err := writer.Write(fc.Values)
	if err != nil {
		return nil, fmt.Errorf("failed to generate %s: %w", fc.Target, err)
	}

	target := ExpandPath(fc.Target)

	result := &GeneratedFile{
		Target:  target,
		Format:  fc.Format,
		Content: string(content),
		Written: false,
	}

	if !m.dryRun {
		// Ensure parent directory exists
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return nil, fmt.Errorf("failed to create directory for %s: %w", target, err)
		}

		if err := os.WriteFile(target, content, 0o644); err != nil {
			return nil, fmt.Errorf("failed to write %s: %w", target, err)
		}
		result.Written = true
	}

	return result, nil
}

// GenerateFiles generates all config files for a component.
func (m *Manager) GenerateFiles(files []componentconfig.FileConfig) ([]*GeneratedFile, error) {
	results := make([]*GeneratedFile, 0, len(files))

	for _, fc := range files {
		result, err := m.GenerateFile(fc)
		if err != nil {
			return results, err
		}
		results = append(results, result)
	}

	return results, nil
}
