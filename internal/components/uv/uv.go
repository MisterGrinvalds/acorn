// Package uv provides UV (Astral) Python package manager helper functionality.
package uv

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Status represents UV installation status.
type Status struct {
	Installed    bool   `json:"installed" yaml:"installed"`
	Version      string `json:"version,omitempty" yaml:"version,omitempty"`
	CacheDir     string `json:"cache_dir,omitempty" yaml:"cache_dir,omitempty"`
	CacheSize    string `json:"cache_size,omitempty" yaml:"cache_size,omitempty"`
	PythonPath   string `json:"python_path,omitempty" yaml:"python_path,omitempty"`
	InProject    bool   `json:"in_project" yaml:"in_project"`
	ProjectName  string `json:"project_name,omitempty" yaml:"project_name,omitempty"`
}

// PythonVersion represents an installed Python version.
type PythonVersion struct {
	Version string `json:"version" yaml:"version"`
	Path    string `json:"path" yaml:"path"`
}

// ToolInfo represents an installed UV tool.
type ToolInfo struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
}

// Helper provides UV helper operations.
type Helper struct {
	verbose bool
	dryRun  bool
}

// NewHelper creates a new UV Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	return &Helper{
		verbose: verbose,
		dryRun:  dryRun,
	}
}

// GetCacheDir returns the UV cache directory.
func (h *Helper) GetCacheDir() string {
	if cacheDir := os.Getenv("UV_CACHE_DIR"); cacheDir != "" {
		return cacheDir
	}
	cacheHome := os.Getenv("XDG_CACHE_HOME")
	if cacheHome == "" {
		home, _ := os.UserHomeDir()
		cacheHome = filepath.Join(home, ".cache")
	}
	return filepath.Join(cacheHome, "uv")
}

// GetStatus returns UV status information.
func (h *Helper) GetStatus() *Status {
	status := &Status{
		CacheDir: h.GetCacheDir(),
	}

	// Check if UV is installed
	out, err := exec.Command("uv", "--version").Output()
	if err != nil {
		status.Installed = false
		return status
	}

	status.Installed = true
	status.Version = strings.TrimSpace(string(out))

	// Get cache size
	status.CacheSize = h.getCacheSize()

	// Get Python path
	if pythonOut, err := exec.Command("uv", "python", "find").Output(); err == nil {
		status.PythonPath = strings.TrimSpace(string(pythonOut))
	}

	// Check if in a UV project
	if _, err := os.Stat("pyproject.toml"); err == nil {
		status.InProject = true
		status.ProjectName = h.getProjectName()
	}

	return status
}

// getProjectName reads the project name from pyproject.toml.
func (h *Helper) getProjectName() string {
	data, err := os.ReadFile("pyproject.toml")
	if err != nil {
		return ""
	}

	// Simple parsing for name
	for line := range strings.SplitSeq(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "name") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				name := strings.TrimSpace(parts[1])
				name = strings.Trim(name, "\"'")
				return name
			}
		}
	}
	return ""
}

// getCacheSize returns the UV cache size.
func (h *Helper) getCacheSize() string {
	cacheDir := h.GetCacheDir()
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		return "0"
	}

	cmd := exec.Command("du", "-sh", cacheDir)
	out, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	parts := strings.Fields(string(out))
	if len(parts) > 0 {
		return parts[0]
	}
	return "unknown"
}

// ListPythonVersions lists installed Python versions.
func (h *Helper) ListPythonVersions() ([]PythonVersion, error) {
	cmd := exec.Command("uv", "python", "list", "--only-installed")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list Python versions: %w", err)
	}

	var versions []PythonVersion
	for line := range strings.SplitSeq(strings.TrimSpace(string(out)), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			versions = append(versions, PythonVersion{
				Version: parts[0],
				Path:    parts[1],
			})
		}
	}

	return versions, nil
}

// InstallPython installs a Python version.
func (h *Helper) InstallPython(version string) error {
	if version == "" {
		return fmt.Errorf("version is required")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: uv python install %s\n", version)
		return nil
	}

	cmd := exec.Command("uv", "python", "install", version)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ListTools lists installed UV tools.
func (h *Helper) ListTools() ([]ToolInfo, error) {
	cmd := exec.Command("uv", "tool", "list")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}

	var tools []ToolInfo
	for line := range strings.SplitSeq(strings.TrimSpace(string(out)), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "-") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 1 {
			tool := ToolInfo{Name: parts[0]}
			if len(parts) >= 2 {
				tool.Version = parts[1]
			}
			tools = append(tools, tool)
		}
	}

	return tools, nil
}

// InstallTool installs a tool using UV.
func (h *Helper) InstallTool(name string) error {
	if name == "" {
		return fmt.Errorf("tool name is required")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: uv tool install %s\n", name)
		return nil
	}

	cmd := exec.Command("uv", "tool", "install", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunTool runs a tool using uvx.
func (h *Helper) RunTool(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("tool name is required")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: uvx %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("uvx", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// CleanCache cleans the UV cache.
func (h *Helper) CleanCache() error {
	if h.dryRun {
		fmt.Println("[dry-run] would run: uv cache clean")
		return nil
	}

	cmd := exec.Command("uv", "cache", "clean")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// PruneCache prunes the UV cache.
func (h *Helper) PruneCache() error {
	if h.dryRun {
		fmt.Println("[dry-run] would run: uv cache prune")
		return nil
	}

	cmd := exec.Command("uv", "cache", "prune")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Init initializes a new UV project.
func (h *Helper) Init(name string) error {
	args := []string{"init"}
	if name != "" {
		args = append(args, name)
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: uv %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("uv", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Sync syncs project dependencies.
func (h *Helper) Sync(args ...string) error {
	cmdArgs := append([]string{"sync"}, args...)

	if h.dryRun {
		fmt.Printf("[dry-run] would run: uv %s\n", strings.Join(cmdArgs, " "))
		return nil
	}

	cmd := exec.Command("uv", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Add adds packages to the project.
func (h *Helper) Add(packages ...string) error {
	if len(packages) == 0 {
		return fmt.Errorf("no packages specified")
	}

	cmdArgs := append([]string{"add"}, packages...)

	if h.dryRun {
		fmt.Printf("[dry-run] would run: uv %s\n", strings.Join(cmdArgs, " "))
		return nil
	}

	cmd := exec.Command("uv", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Remove removes packages from the project.
func (h *Helper) Remove(packages ...string) error {
	if len(packages) == 0 {
		return fmt.Errorf("no packages specified")
	}

	cmdArgs := append([]string{"remove"}, packages...)

	if h.dryRun {
		fmt.Printf("[dry-run] would run: uv %s\n", strings.Join(cmdArgs, " "))
		return nil
	}

	cmd := exec.Command("uv", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Run executes a command in the project environment.
func (h *Helper) Run(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command specified")
	}

	cmdArgs := append([]string{"run"}, args...)

	if h.dryRun {
		fmt.Printf("[dry-run] would run: uv %s\n", strings.Join(cmdArgs, " "))
		return nil
	}

	cmd := exec.Command("uv", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Venv creates a virtual environment.
func (h *Helper) Venv(name string) error {
	if name == "" {
		name = ".venv"
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: uv venv %s\n", name)
		return nil
	}

	cmd := exec.Command("uv", "venv", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Lock locks project dependencies.
func (h *Helper) Lock() error {
	if h.dryRun {
		fmt.Println("[dry-run] would run: uv lock")
		return nil
	}

	cmd := exec.Command("uv", "lock")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Upgrade upgrades packages.
func (h *Helper) Upgrade(packages ...string) error {
	cmdArgs := []string{"lock", "--upgrade"}
	if len(packages) > 0 {
		for _, pkg := range packages {
			cmdArgs = append(cmdArgs, "--upgrade-package", pkg)
		}
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: uv %s\n", strings.Join(cmdArgs, " "))
		return nil
	}

	cmd := exec.Command("uv", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// PipInstall installs packages using uv pip.
func (h *Helper) PipInstall(packages ...string) error {
	if len(packages) == 0 {
		return fmt.Errorf("no packages specified")
	}

	cmdArgs := append([]string{"pip", "install"}, packages...)

	if h.dryRun {
		fmt.Printf("[dry-run] would run: uv %s\n", strings.Join(cmdArgs, " "))
		return nil
	}

	cmd := exec.Command("uv", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Install installs UV.
func (h *Helper) Install() error {
	if _, err := exec.LookPath("uv"); err == nil {
		return fmt.Errorf("uv is already installed")
	}

	if h.dryRun {
		fmt.Println("[dry-run] would install UV using the official installer")
		return nil
	}

	fmt.Println("Installing UV...")
	cmd := exec.Command("bash", "-c", "curl -LsSf https://astral.sh/uv/install.sh | sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// SelfUpdate updates UV to the latest version.
func (h *Helper) SelfUpdate() error {
	if h.dryRun {
		fmt.Println("[dry-run] would run: uv self update")
		return nil
	}

	cmd := exec.Command("uv", "self", "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Tree shows the dependency tree.
func (h *Helper) Tree() error {
	if h.dryRun {
		fmt.Println("[dry-run] would run: uv tree")
		return nil
	}

	cmd := exec.Command("uv", "tree")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Export exports dependencies to requirements.txt format.
func (h *Helper) Export(output string) error {
	args := []string{"export"}
	if output != "" {
		args = append(args, "-o", output)
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: uv %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("uv", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetCacheInfo returns cache information.
func (h *Helper) GetCacheInfo() (map[string]interface{}, error) {
	cmd := exec.Command("uv", "cache", "dir")
	dirOut, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	info := map[string]interface{}{
		"directory": strings.TrimSpace(string(dirOut)),
		"size":      h.getCacheSize(),
	}

	return info, nil
}
