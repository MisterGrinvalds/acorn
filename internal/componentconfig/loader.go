package componentconfig

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Loader handles config loading with XDG overlay support.
type Loader struct {
	xdgConfigHome string
	acornDir      string
}

// NewLoader creates a loader with XDG-compliant paths.
func NewLoader() *Loader {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, _ := os.UserHomeDir()
		configHome = filepath.Join(home, ".config")
	}

	return &Loader{
		xdgConfigHome: configHome,
		acornDir:      filepath.Join(configHome, "acorn"),
	}
}

// Load loads a component config with user overlay.
// First loads embedded defaults, then merges user overrides if present.
func (l *Loader) Load(component string, target interface{}) error {
	// 1. Load embedded default
	defaultData, err := GetEmbeddedConfig(component)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(defaultData, target); err != nil {
		return fmt.Errorf("failed to parse embedded config for %s: %w", component, err)
	}

	// 2. Load user override if exists and merge
	userPath := l.UserConfigPath(component)
	userData, err := os.ReadFile(userPath)
	if err == nil {
		// User config exists - merge over defaults
		if err := yaml.Unmarshal(userData, target); err != nil {
			return fmt.Errorf("failed to parse user config at %s: %w", userPath, err)
		}
	}

	return nil
}

// LoadBase loads a component config into a BaseConfig struct.
func (l *Loader) LoadBase(component string) (*BaseConfig, error) {
	cfg := &BaseConfig{}
	if err := l.Load(component, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// UserConfigPath returns the path for user override config.
func (l *Loader) UserConfigPath(component string) string {
	return filepath.Join(l.acornDir, component+".yaml")
}

// HasUserOverride checks if user has a config override for a component.
func (l *Loader) HasUserOverride(component string) bool {
	_, err := os.Stat(l.UserConfigPath(component))
	return err == nil
}

// EnsureUserConfigDir creates the user config directory if it doesn't exist.
func (l *Loader) EnsureUserConfigDir() error {
	return os.MkdirAll(l.acornDir, 0o755)
}

// GetAcornDir returns the acorn config directory path.
func (l *Loader) GetAcornDir() string {
	return l.acornDir
}

// CreateUserOverride creates a user override file with the component's default config.
func (l *Loader) CreateUserOverride(component string) error {
	// Get embedded default
	defaultData, err := GetEmbeddedConfig(component)
	if err != nil {
		return err
	}

	// Ensure directory exists
	if err := l.EnsureUserConfigDir(); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write user override file
	userPath := l.UserConfigPath(component)
	if err := os.WriteFile(userPath, defaultData, 0o644); err != nil {
		return fmt.Errorf("failed to write user config: %w", err)
	}

	return nil
}

// MergeConfigs merges two BaseConfig structs.
// Values from 'override' take precedence over 'base'.
// Maps are merged (override keys win), slices are replaced.
func MergeConfigs(base, override *BaseConfig) *BaseConfig {
	result := &BaseConfig{
		Name:        coalesce(override.Name, base.Name),
		Description: coalesce(override.Description, base.Description),
		Version:     coalesce(override.Version, base.Version),
	}

	// Merge env maps
	result.Env = mergeMaps(base.Env, override.Env)

	// Merge aliases maps
	result.Aliases = mergeMaps(base.Aliases, override.Aliases)

	// Merge shell functions maps
	result.ShellFunctions = mergeMaps(base.ShellFunctions, override.ShellFunctions)

	// Slices are replaced, not merged
	if len(override.Paths) > 0 {
		result.Paths = override.Paths
	} else {
		result.Paths = base.Paths
	}

	if len(override.Wrappers) > 0 {
		result.Wrappers = override.Wrappers
	} else {
		result.Wrappers = base.Wrappers
	}

	return result
}

// coalesce returns the first non-empty string.
func coalesce(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

// mergeMaps merges two string maps, with 'b' values taking precedence.
func mergeMaps(a, b map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range a {
		result[k] = v
	}
	for k, v := range b {
		result[k] = v
	}
	return result
}
