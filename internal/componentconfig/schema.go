// Package componentconfig provides typed configuration for acorn shell components.
package componentconfig

// BaseConfig is the common structure for all component configurations.
// Each component's config.yaml should follow this schema.
type BaseConfig struct {
	// Metadata
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Version     string `yaml:"version,omitempty"`

	// Environment variables to export (key-value)
	Env map[string]string `yaml:"env,omitempty"`

	// Paths to add to PATH
	Paths []PathEntry `yaml:"paths,omitempty"`

	// Shell aliases (name -> command)
	Aliases map[string]string `yaml:"aliases,omitempty"`

	// Wrapper functions that call acorn commands
	Wrappers []Wrapper `yaml:"wrappers,omitempty"`

	// Raw shell functions (for interactive code that can't be generated)
	ShellFunctions map[string]string `yaml:"shell_functions,omitempty"`
}

// PathEntry represents a path to add to PATH with optional conditions.
type PathEntry struct {
	// Path is the path to add (supports shell variable expansion)
	Path string `yaml:"path"`

	// Condition is an optional platform filter: "darwin", "linux", or empty for all
	Condition string `yaml:"condition,omitempty"`
}

// Wrapper defines an acorn command wrapper function.
type Wrapper struct {
	// Name is the shell function name
	Name string `yaml:"name"`

	// Command is the acorn command to call (e.g., "acorn go new")
	Command string `yaml:"command"`

	// Usage is the usage hint shown when called without args (optional)
	Usage string `yaml:"usage,omitempty"`

	// DefaultArg is a default argument if none provided (optional)
	DefaultArg string `yaml:"default_arg,omitempty"`

	// PostAction is a special action after command: "cd" to change directory
	PostAction string `yaml:"post_action,omitempty"`

	// RequiresArg if true, shows usage and returns 1 if no arg provided
	RequiresArg bool `yaml:"requires_arg,omitempty"`
}

// GetEnv returns environment variables as a map, never nil.
func (c *BaseConfig) GetEnv() map[string]string {
	if c.Env == nil {
		return make(map[string]string)
	}
	return c.Env
}

// GetAliases returns aliases as a map, never nil.
func (c *BaseConfig) GetAliases() map[string]string {
	if c.Aliases == nil {
		return make(map[string]string)
	}
	return c.Aliases
}

// GetShellFunctions returns shell functions as a map, never nil.
func (c *BaseConfig) GetShellFunctions() map[string]string {
	if c.ShellFunctions == nil {
		return make(map[string]string)
	}
	return c.ShellFunctions
}
