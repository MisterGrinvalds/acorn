// Package shell provides shell integration generation and injection.
package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"text/template"

	"github.com/mistergrinvalds/acorn/internal/utils/config"
	"github.com/mistergrinvalds/acorn/internal/utils/configfile"
	rootconfig "github.com/mistergrinvalds/acorn/config"
)

// Config holds shell integration configuration.
type Config struct {
	XDGConfigHome string
	AcornDir      string
	Shell         string // bash or zsh
	Platform      string // darwin or linux
	Verbose       bool
	DryRun        bool
}

// NewConfig creates a new Config with defaults.
func NewConfig(verbose, dryRun bool) *Config {
	xdgConfig := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfig == "" {
		home, _ := os.UserHomeDir()
		xdgConfig = filepath.Join(home, ".config")
	}

	shell := detectShell()
	platform := runtime.GOOS

	return &Config{
		XDGConfigHome: xdgConfig,
		AcornDir:      filepath.Join(xdgConfig, "acorn"),
		Shell:         shell,
		Platform:      platform,
		Verbose:       verbose,
		DryRun:        dryRun,
	}
}

// detectShell detects the current shell.
func detectShell() string {
	shell := os.Getenv("SHELL")
	if strings.Contains(shell, "zsh") {
		return "zsh"
	}
	return "bash"
}

// Component represents a shell component with its scripts.
type Component struct {
	Name        string
	Description string
	Env         string // environment variable setup
	Aliases     string // shell aliases
	Functions   string // shell functions (wrappers that call acorn)
	Completions string // shell completions
}

// GeneratedScript represents a generated shell script with metadata.
type GeneratedScript struct {
	Component    string `json:"component" yaml:"component"`
	Description  string `json:"description" yaml:"description"`
	GeneratedPath string `json:"generated_path" yaml:"generated_path"` // Where file was written (generated/shell/)
	SymlinkPath  string `json:"symlink_path" yaml:"symlink_path"`       // Where symlink should point (XDG)
	Content      string `json:"content" yaml:"content"`
	Written      bool   `json:"written" yaml:"written"`
}

// GenerateResult contains the result of a generate operation.
type GenerateResult struct {
	AcornDir    string                    `json:"acorn_dir" yaml:"acorn_dir"`
	Shell       string                    `json:"shell" yaml:"shell"`
	Platform    string                    `json:"platform" yaml:"platform"`
	DryRun      bool                      `json:"dry_run" yaml:"dry_run"`
	Scripts     []*GeneratedScript        `json:"scripts" yaml:"scripts"`
	Entrypoint  *GeneratedScript          `json:"entrypoint,omitempty" yaml:"entrypoint,omitempty"`
	ConfigFiles []*configfile.GeneratedFile `json:"config_files,omitempty" yaml:"config_files,omitempty"`
}

// InjectResult contains the result of an inject/eject operation.
type InjectResult struct {
	RCFile         string `json:"rc_file" yaml:"rc_file"`
	EntrypointPath string `json:"entrypoint_path" yaml:"entrypoint_path"`
	Action         string `json:"action" yaml:"action"` // "injected", "ejected", "already_injected", "not_injected"
	DryRun         bool   `json:"dry_run" yaml:"dry_run"`
	InjectionBlock string `json:"injection_block,omitempty" yaml:"injection_block,omitempty"`
}

// Manager handles shell script generation and injection.
type Manager struct {
	config     *Config
	components map[string]*Component
	fileSpecs  map[string][]FileSpec // component name -> file specs for config file generation
}

// FileSpec holds file generation specification.
type FileSpec struct {
	Target    string                 `json:"target" yaml:"target"`
	Format    string                 `json:"format" yaml:"format"`
	Platforms []string               `json:"platforms,omitempty" yaml:"platforms,omitempty"`
	Values    map[string]interface{} `json:"values" yaml:"values"`
}

// NewManager creates a new shell Manager.
func NewManager(config *Config) *Manager {
	return &Manager{
		config:     config,
		components: make(map[string]*Component),
		fileSpecs:  make(map[string][]FileSpec),
	}
}

// RegisterComponent registers a component for shell integration.
func (m *Manager) RegisterComponent(c *Component) {
	m.components[c.Name] = c
}

// RegisterComponentFiles registers config files for a component.
func (m *Manager) RegisterComponentFiles(name string, files []FileSpec) {
	if len(files) > 0 {
		m.fileSpecs[name] = files
	}
}

// EnsureDir ensures the acorn config directory exists.
func (m *Manager) EnsureDir() error {
	if m.config.DryRun {
		return nil
	}
	return os.MkdirAll(m.config.AcornDir, 0o755)
}

// ListComponents returns all registered component names.
func (m *Manager) ListComponents() []string {
	names := make([]string, 0, len(m.components))
	for name := range m.components {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// GetComponent returns a component by name.
func (m *Manager) GetComponent(name string) (*Component, bool) {
	c, ok := m.components[name]
	return c, ok
}

// getGeneratedShellDir returns the directory for generated shell scripts.
// Scripts are written to .sapling/generated/shell/
func (m *Manager) getGeneratedShellDir() string {
	// Try to use .sapling/generated/shell
	if genDir, err := rootconfig.GeneratedDir(); err == nil {
		return filepath.Join(genDir, "shell")
	}

	// Fallback: use DOTFILES_ROOT if set
	if dotfilesRoot := os.Getenv("DOTFILES_ROOT"); dotfilesRoot != "" {
		return filepath.Join(dotfilesRoot, "generated", "shell")
	}

	// Last resort: derive from home directory
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Repos", "personal", "tools", "generated", "shell")
}

// GenerateComponent generates a shell script for a single component.
// Returns structured result with script content and metadata.
func (m *Manager) GenerateComponent(name string) (*GenerateResult, error) {
	c, ok := m.components[name]
	if !ok {
		return nil, fmt.Errorf("component not found: %s (available: %v)", name, m.ListComponents())
	}

	generatedDir := m.getGeneratedShellDir()
	if err := os.MkdirAll(generatedDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create generated directory: %w", err)
	}

	script := m.generateComponentScript(c)
	generatedPath := filepath.Join(generatedDir, name+".sh")
	symlinkPath := filepath.Join(m.config.AcornDir, name+".sh")

	genScript := &GeneratedScript{
		Component:     name,
		Description:   c.Description,
		GeneratedPath: generatedPath,
		SymlinkPath:   symlinkPath,
		Content:       script,
		Written:       false,
	}

	// Write file if not dry-run
	if !m.config.DryRun {
		if err := os.WriteFile(generatedPath, []byte(script), 0o644); err != nil {
			return nil, fmt.Errorf("failed to write %s: %w", generatedPath, err)
		}
		genScript.Written = true
	}

	return &GenerateResult{
		AcornDir: m.config.AcornDir,
		Shell:    m.config.Shell,
		Platform: m.config.Platform,
		DryRun:   m.config.DryRun,
		Scripts:  []*GeneratedScript{genScript},
	}, nil
}

// GenerateComponents generates shell scripts for specific components.
// If names is empty, generates for all components.
// Shell scripts are written to $DOTFILES_ROOT/generated/shell/ and should be
// symlinked to $XDG_CONFIG_HOME/acorn/ via `acorn sync link`.
func (m *Manager) GenerateComponents(names ...string) (*GenerateResult, error) {
	// Ensure generated shell directory exists
	generatedShellDir := m.getGeneratedShellDir()
	if err := os.MkdirAll(generatedShellDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create generated shell directory: %w", err)
	}

	// If no names specified, use all components
	if len(names) == 0 {
		names = m.ListComponents()
	}

	result := &GenerateResult{
		AcornDir:    m.config.AcornDir,
		Shell:       m.config.Shell,
		Platform:    m.config.Platform,
		DryRun:      m.config.DryRun,
		Scripts:     make([]*GeneratedScript, 0, len(names)),
		ConfigFiles: make([]*configfile.GeneratedFile, 0),
	}

	// Create config file manager with generated directory
	// Config files are written to $DOTFILES_ROOT/generated/{component}/{filename}
	generatedDir := filepath.Dir(generatedShellDir) // parent of shell/ is generated/
	cfManager := configfile.NewManagerWithGeneratedDir(generatedDir, m.config.DryRun)

	// Generate each component script
	for _, name := range names {
		c, ok := m.components[name]
		if !ok {
			return nil, fmt.Errorf("component not found: %s (available: %v)", name, m.ListComponents())
		}

		script := m.generateComponentScript(c)
		generatedPath := filepath.Join(generatedShellDir, name+".sh")
		symlinkPath := filepath.Join(m.config.AcornDir, name+".sh")

		genScript := &GeneratedScript{
			Component:     name,
			Description:   c.Description,
			GeneratedPath: generatedPath,
			SymlinkPath:   symlinkPath,
			Content:       script,
			Written:       false,
		}

		if !m.config.DryRun {
			if err := os.WriteFile(generatedPath, []byte(script), 0o644); err != nil {
				return nil, fmt.Errorf("failed to write %s: %w", generatedPath, err)
			}
			genScript.Written = true
		}

		result.Scripts = append(result.Scripts, genScript)

		// Generate config files for this component
		if files, ok := m.fileSpecs[name]; ok {
			for _, spec := range files {
				// Skip files not for this platform
				if !m.shouldGenerateForPlatform(spec.Platforms) {
					continue
				}
				fc := componentConfigFromSpec(spec)
				genFile, err := cfManager.GenerateFileForComponent(name, fc)
				if err != nil {
					return nil, fmt.Errorf("failed to generate config for %s: %w", name, err)
				}
				result.ConfigFiles = append(result.ConfigFiles, genFile)
			}
		}
	}

	return result, nil
}

// componentConfigFromSpec converts a FileSpec to config.FileConfig.
func componentConfigFromSpec(spec FileSpec) config.FileConfig {
	return config.FileConfig{
		Target:    spec.Target,
		Format:    spec.Format,
		Platforms: spec.Platforms,
		Values:    spec.Values,
	}
}

// shouldGenerateForPlatform checks if a file should be generated for the current platform.
func (m *Manager) shouldGenerateForPlatform(platforms []string) bool {
	// If no platforms specified, generate for all
	if len(platforms) == 0 {
		return true
	}

	// Check if current platform is in the list
	for _, p := range platforms {
		if p == m.config.Platform {
			return true
		}
	}

	return false
}

// GenerateAll generates all component shell scripts and the entrypoint.
// All scripts are written to $DOTFILES_ROOT/generated/shell/ and should be
// symlinked to $XDG_CONFIG_HOME/acorn/ via `acorn sync link`.
func (m *Manager) GenerateAll() (*GenerateResult, error) {
	result, err := m.GenerateComponents() // all components
	if err != nil {
		return nil, err
	}

	// Generate the main entrypoint
	// Named "shell.sh" as the primary entrypoint sourced by rc files
	// Written to generated/shell/shell.sh, symlinked to ~/.config/acorn/shell.sh
	entrypoint := m.generateEntrypoint()
	generatedShellDir := m.getGeneratedShellDir()
	generatedPath := filepath.Join(generatedShellDir, "shell.sh")
	symlinkPath := filepath.Join(m.config.AcornDir, "shell.sh")

	result.Entrypoint = &GeneratedScript{
		Component:     "shell",
		Description:   "Main entrypoint that sources all component scripts",
		GeneratedPath: generatedPath,
		SymlinkPath:   symlinkPath,
		Content:       entrypoint,
		Written:       false,
	}

	if !m.config.DryRun {
		if err := os.WriteFile(generatedPath, []byte(entrypoint), 0o644); err != nil {
			return nil, fmt.Errorf("failed to write entrypoint: %w", err)
		}
		result.Entrypoint.Written = true
	}

	return result, nil
}

// generateComponentScript generates a shell script for a component.
func (m *Manager) generateComponentScript(c *Component) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("#!/bin/sh\n"))
	b.WriteString(fmt.Sprintf("# Acorn shell integration: %s\n", c.Name))
	b.WriteString(fmt.Sprintf("# %s\n", c.Description))
	b.WriteString("# Generated by acorn - do not edit manually\n\n")

	if c.Env != "" {
		b.WriteString("# Environment\n")
		b.WriteString(c.Env)
		b.WriteString("\n")
	}

	if c.Aliases != "" {
		b.WriteString("# Aliases\n")
		b.WriteString(c.Aliases)
		b.WriteString("\n")
	}

	if c.Functions != "" {
		b.WriteString("# Functions\n")
		b.WriteString(c.Functions)
		b.WriteString("\n")
	}

	if c.Completions != "" {
		b.WriteString("# Completions\n")
		b.WriteString(c.Completions)
		b.WriteString("\n")
	}

	return b.String()
}

// generateEntrypoint generates the main shell.sh entrypoint.
// Components are sourced in the order defined by GetComponentOrder() to ensure
// dependencies are met (e.g., shell before theme, xdg before everything else).
func (m *Manager) generateEntrypoint() string {
	var b strings.Builder

	b.WriteString("#!/bin/sh\n")
	b.WriteString("# Acorn shell integration entrypoint\n")
	b.WriteString("# Generated by acorn - do not edit manually\n")
	b.WriteString("# Source this file from your shell rc file\n\n")

	b.WriteString("# Acorn configuration directory\n")
	b.WriteString(fmt.Sprintf("ACORN_CONFIG_DIR=\"%s\"\n", m.config.AcornDir))
	b.WriteString("export ACORN_CONFIG_DIR\n\n")

	b.WriteString("# Source all component scripts in dependency order\n")
	// Use GetComponentOrder() to maintain correct loading order
	for _, name := range GetComponentOrder() {
		// Only include components that are registered
		if _, ok := m.components[name]; ok {
			b.WriteString(fmt.Sprintf("[ -f \"$ACORN_CONFIG_DIR/%s.sh\" ] && . \"$ACORN_CONFIG_DIR/%s.sh\"\n", name, name))
		}
	}

	b.WriteString("\n# Acorn CLI completions\n")
	b.WriteString("if command -v acorn >/dev/null 2>&1; then\n")
	if m.config.Shell == "zsh" {
		b.WriteString("    eval \"$(acorn completion zsh)\"\n")
	} else {
		b.WriteString("    eval \"$(acorn completion bash)\"\n")
	}
	b.WriteString("fi\n")

	return b.String()
}

// GetRCFile returns the shell rc file path.
// For bash on macOS, returns .bash_profile (login shell default).
// For bash on Linux, returns .bashrc (interactive shell default).
// For zsh, returns .zshrc (works for both).
func (m *Manager) GetRCFile() string {
	home, _ := os.UserHomeDir()

	if m.config.Shell == "zsh" {
		return filepath.Join(home, ".zshrc")
	}

	// Bash: use .bash_profile on macOS (login shells), .bashrc on Linux
	if m.config.Platform == "darwin" {
		return filepath.Join(home, ".bash_profile")
	}
	return filepath.Join(home, ".bashrc")
}

// InjectMarker is the comment used to identify acorn injections.
const InjectMarker = "# >>> acorn shell integration >>>"
const InjectMarkerEnd = "# <<< acorn shell integration <<<"

// Inject adds the acorn source line to the shell rc file.
func (m *Manager) Inject() (*InjectResult, error) {
	rcFile := m.GetRCFile()
	entrypointPath := filepath.Join(m.config.AcornDir, "shell.sh")

	result := &InjectResult{
		RCFile:         rcFile,
		EntrypointPath: entrypointPath,
		DryRun:         m.config.DryRun,
	}

	// Read existing rc file
	content, err := os.ReadFile(rcFile)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to read %s: %w", rcFile, err)
	}

	// Check if already injected
	if strings.Contains(string(content), InjectMarker) {
		result.Action = "already_injected"
		return result, nil
	}

	// Create injection block with ACORN_CONFIG_DIR variable
	injection := fmt.Sprintf("\n%s\nexport ACORN_CONFIG_DIR=\"%s\"\n[ -f \"$ACORN_CONFIG_DIR/shell.sh\" ] && . \"$ACORN_CONFIG_DIR/shell.sh\"\n%s\n",
		InjectMarker, m.config.AcornDir, InjectMarkerEnd)
	result.InjectionBlock = injection

	if m.config.DryRun {
		result.Action = "would_inject"
		return result, nil
	}

	// Append to rc file
	f, err := os.OpenFile(rcFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", rcFile, err)
	}
	defer f.Close()

	if _, err := f.WriteString(injection); err != nil {
		return nil, fmt.Errorf("failed to write to %s: %w", rcFile, err)
	}

	result.Action = "injected"
	return result, nil
}

// Eject removes the acorn source line from the shell rc file.
func (m *Manager) Eject() (*InjectResult, error) {
	rcFile := m.GetRCFile()
	entrypointPath := filepath.Join(m.config.AcornDir, "shell.sh")

	result := &InjectResult{
		RCFile:         rcFile,
		EntrypointPath: entrypointPath,
		DryRun:         m.config.DryRun,
	}

	content, err := os.ReadFile(rcFile)
	if err != nil {
		if os.IsNotExist(err) {
			result.Action = "not_injected"
			return result, nil
		}
		return nil, fmt.Errorf("failed to read %s: %w", rcFile, err)
	}

	if !strings.Contains(string(content), InjectMarker) {
		result.Action = "not_injected"
		return result, nil
	}

	// Remove the injection block
	lines := strings.Split(string(content), "\n")
	var newLines []string
	inBlock := false

	for _, line := range lines {
		if strings.Contains(line, InjectMarker) {
			inBlock = true
			continue
		}
		if strings.Contains(line, InjectMarkerEnd) {
			inBlock = false
			continue
		}
		if !inBlock {
			newLines = append(newLines, line)
		}
	}

	newContent := strings.Join(newLines, "\n")

	if m.config.DryRun {
		result.Action = "would_eject"
		return result, nil
	}

	if err := os.WriteFile(rcFile, []byte(newContent), 0o644); err != nil {
		return nil, fmt.Errorf("failed to write %s: %w", rcFile, err)
	}

	result.Action = "ejected"
	return result, nil
}

// Status returns the current shell integration status.
type Status struct {
	Shell          string   `json:"shell" yaml:"shell"`
	Platform       string   `json:"platform" yaml:"platform"`
	AcornDir       string   `json:"acorn_dir" yaml:"acorn_dir"`
	AcornDirExists bool     `json:"acorn_dir_exists" yaml:"acorn_dir_exists"`
	RCFile         string   `json:"rc_file" yaml:"rc_file"`
	Injected       bool     `json:"injected" yaml:"injected"`
	Components     []string `json:"components" yaml:"components"`
	GeneratedFiles []string `json:"generated_files" yaml:"generated_files"`
}

// GetStatus returns the current shell integration status.
func (m *Manager) GetStatus() (*Status, error) {
	status := &Status{
		Shell:    m.config.Shell,
		Platform: m.config.Platform,
		AcornDir: m.config.AcornDir,
		RCFile:   m.GetRCFile(),
	}

	// Check acorn dir
	if _, err := os.Stat(m.config.AcornDir); err == nil {
		status.AcornDirExists = true

		// List generated files
		entries, _ := os.ReadDir(m.config.AcornDir)
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".sh") {
				status.GeneratedFiles = append(status.GeneratedFiles, e.Name())
			}
		}
	}

	// Check injection
	content, err := os.ReadFile(status.RCFile)
	if err == nil {
		status.Injected = strings.Contains(string(content), InjectMarker)
	}

	// List registered components
	for name := range m.components {
		status.Components = append(status.Components, name)
	}

	return status, nil
}

// TemplateData provides data for shell script templates.
type TemplateData struct {
	Shell    string
	Platform string
	AcornBin string
}

// ExecuteTemplate executes a template string with the given data.
func ExecuteTemplate(tmpl string, data TemplateData) (string, error) {
	t, err := template.New("shell").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	if err := t.Execute(&b, data); err != nil {
		return "", err
	}

	return b.String(), nil
}
