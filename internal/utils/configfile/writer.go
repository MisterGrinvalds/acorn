// Package configfile provides config file generation in various formats.
package configfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mistergrinvalds/acorn/internal/utils/config"
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
	// GeneratedPath is where the file was written (in the repo's generated/ directory)
	GeneratedPath string `json:"generated_path" yaml:"generated_path"`
	// SymlinkTarget is the XDG path where a symlink should point to this file
	SymlinkTarget string `json:"symlink_target" yaml:"symlink_target"`
	Format        string `json:"format" yaml:"format"`
	Content       string `json:"content" yaml:"content"`
	Written       bool   `json:"written" yaml:"written"`
	// Component is the name of the component this file belongs to
	Component string `json:"component" yaml:"component"`
}

// Manager handles config file generation for components.
type Manager struct {
	dryRun       bool
	generatedDir string // Directory where generated files are written (e.g., $DOTFILES_ROOT/generated)
}

// NewManager creates a new config file manager.
// If generatedDir is empty, files are written directly to their target paths (legacy behavior).
func NewManager(dryRun bool) *Manager {
	return &Manager{dryRun: dryRun}
}

// NewManagerWithGeneratedDir creates a config file manager that writes to a generated directory.
// Files are written to generatedDir/{component}/{filename} and symlink targets are tracked.
func NewManagerWithGeneratedDir(generatedDir string, dryRun bool) *Manager {
	return &Manager{
		dryRun:       dryRun,
		generatedDir: generatedDir,
	}
}

// SetGeneratedDir sets the directory where generated files should be written.
func (m *Manager) SetGeneratedDir(dir string) {
	m.generatedDir = dir
}

// ExpandPath expands environment variables in a path.
// Supports shell-style ${VAR:-default} syntax.
func ExpandPath(path string) string {
	// Handle ${VAR:-default} syntax first
	result := path
	for {
		start := strings.Index(result, "${")
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], "}")
		if end == -1 {
			break
		}
		end += start

		// Extract the variable expression
		expr := result[start+2 : end]

		// Check for :- default syntax
		var value string
		if idx := strings.Index(expr, ":-"); idx != -1 {
			varName := expr[:idx]
			defaultVal := expr[idx+2:]
			value = os.Getenv(varName)
			if value == "" {
				value = defaultVal
			}
		} else {
			value = os.Getenv(expr)
		}

		// Replace the variable with its value
		result = result[:start] + value + result[end+1:]
	}

	// Handle remaining $VAR syntax
	result = os.Expand(result, func(key string) string {
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
// Deprecated: Use GenerateFileForComponent which tracks symlink targets properly.
func (m *Manager) GenerateFile(fc config.FileConfig) (*GeneratedFile, error) {
	return m.GenerateFileForComponent("unknown", fc)
}

// GenerateFileForComponent generates a config file for a specific component.
// If generatedDir is set, writes to generatedDir/{component}/{filename} and tracks symlink target.
// Otherwise, writes directly to the target path (legacy behavior).
func (m *Manager) GenerateFileForComponent(component string, fc config.FileConfig) (*GeneratedFile, error) {
	writer, err := GetWriter(fc.Format)
	if err != nil {
		return nil, err
	}

	content, err := writer.Write(fc.Values)
	if err != nil {
		return nil, fmt.Errorf("failed to generate %s: %w", fc.Target, err)
	}

	symlinkTarget := ExpandPath(fc.Target)

	// Determine where to write the file
	var writePath string
	if m.generatedDir != "" {
		// Write to generated/{component}/{filename}
		filename := filepath.Base(symlinkTarget)
		writePath = filepath.Join(m.generatedDir, component, filename)
	} else {
		// Legacy: write directly to target
		writePath = symlinkTarget
	}

	result := &GeneratedFile{
		GeneratedPath: writePath,
		SymlinkTarget: symlinkTarget,
		Format:        fc.Format,
		Content:       string(content),
		Written:       false,
		Component:     component,
	}

	if !m.dryRun {
		// Ensure parent directory exists
		if err := os.MkdirAll(filepath.Dir(writePath), 0o755); err != nil {
			return nil, fmt.Errorf("failed to create directory for %s: %w", writePath, err)
		}

		if err := os.WriteFile(writePath, content, 0o644); err != nil {
			return nil, fmt.Errorf("failed to write %s: %w", writePath, err)
		}
		result.Written = true
	}

	return result, nil
}

// GenerateFiles generates all config files for a component.
func (m *Manager) GenerateFiles(files []config.FileConfig) ([]*GeneratedFile, error) {
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
