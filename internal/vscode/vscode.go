// Package vscode provides VS Code integration and project helpers.
package vscode

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Workspace represents a VS Code workspace.
type Workspace struct {
	Name string `json:"name" yaml:"name"`
	Path string `json:"path" yaml:"path"`
}

// Extension represents a VS Code extension.
type Extension struct {
	ID        string `json:"id" yaml:"id"`
	Installed bool   `json:"installed,omitempty" yaml:"installed,omitempty"`
}

// ConfigPaths contains VS Code configuration paths.
type ConfigPaths struct {
	UserDir      string `json:"user_dir" yaml:"user_dir"`
	Settings     string `json:"settings" yaml:"settings"`
	Keybindings  string `json:"keybindings" yaml:"keybindings"`
	Extensions   string `json:"extensions" yaml:"extensions"`
	WorkspaceDir string `json:"workspace_dir" yaml:"workspace_dir"`
}

// ProjectSettings contains VS Code project settings by language.
var ProjectSettings = map[string]string{
	"python": `{
    "python.defaultInterpreterPath": "./.venv/bin/python",
    "python.terminal.activateEnvironment": true,
    "editor.formatOnSave": true,
    "python.testing.pytestEnabled": true,
    "python.testing.pytestArgs": ["."]
}`,
	"go": `{
    "go.formatTool": "goimports",
    "go.lintTool": "golangci-lint",
    "editor.formatOnSave": true,
    "go.testFlags": ["-v"],
    "go.testTimeout": "30s"
}`,
	"typescript": `{
    "typescript.preferences.importModuleSpecifier": "relative",
    "editor.formatOnSave": true,
    "editor.defaultFormatter": "esbenp.prettier-vscode",
    "typescript.updateImportsOnFileMove.enabled": "always"
}`,
	"general": `{
    "editor.formatOnSave": true,
    "files.trimTrailingWhitespace": true,
    "files.insertFinalNewline": true,
    "editor.rulers": [80, 120]
}`,
}

// EssentialExtensions is the list of essential extensions.
var EssentialExtensions = []string{
	"ms-python.python",
	"golang.go",
	"github.vscode-pull-request-github",
	"eamodio.gitlens",
	"ms-kubernetes-tools.vscode-kubernetes-tools",
	"ms-azuretools.vscode-docker",
	"catppuccin.catppuccin-vsc",
	"catppuccin.catppuccin-vsc-icons",
}

// Helper provides VS Code helper operations.
type Helper struct {
	verbose      bool
	dryRun       bool
	dotfilesRoot string
}

// NewHelper creates a new Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	dotfilesRoot := os.Getenv("DOTFILES_ROOT")
	if dotfilesRoot == "" {
		home, _ := os.UserHomeDir()
		dotfilesRoot = filepath.Join(home, ".config", "dotfiles")
	}

	return &Helper{
		verbose:      verbose,
		dryRun:       dryRun,
		dotfilesRoot: dotfilesRoot,
	}
}

// GetConfigPaths returns VS Code configuration paths for the current platform.
func (h *Helper) GetConfigPaths() *ConfigPaths {
	home, _ := os.UserHomeDir()

	var userDir string
	if runtime.GOOS == "darwin" {
		userDir = filepath.Join(home, "Library", "Application Support", "Code", "User")
	} else {
		userDir = filepath.Join(home, ".config", "Code", "User")
	}

	return &ConfigPaths{
		UserDir:      userDir,
		Settings:     filepath.Join(userDir, "settings.json"),
		Keybindings:  filepath.Join(userDir, "keybindings.json"),
		Extensions:   filepath.Join(home, ".vscode", "extensions"),
		WorkspaceDir: filepath.Join(home, ".vscode", "workspaces"),
	}
}

// ListWorkspaces returns available workspaces.
func (h *Helper) ListWorkspaces() ([]Workspace, error) {
	paths := h.GetConfigPaths()

	if _, err := os.Stat(paths.WorkspaceDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("workspaces directory not found: %s", paths.WorkspaceDir)
	}

	entries, err := os.ReadDir(paths.WorkspaceDir)
	if err != nil {
		return nil, err
	}

	var workspaces []Workspace
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".code-workspace") {
			name := strings.TrimSuffix(entry.Name(), ".code-workspace")
			workspaces = append(workspaces, Workspace{
				Name: name,
				Path: filepath.Join(paths.WorkspaceDir, entry.Name()),
			})
		}
	}

	return workspaces, nil
}

// OpenWorkspace opens a workspace in VS Code.
func (h *Helper) OpenWorkspace(name string) error {
	paths := h.GetConfigPaths()
	wsPath := filepath.Join(paths.WorkspaceDir, name+".code-workspace")

	if _, err := os.Stat(wsPath); os.IsNotExist(err) {
		return fmt.Errorf("workspace not found: %s", name)
	}

	return h.runCode(wsPath)
}

// CreateProject creates a new VS Code project with language-specific settings.
func (h *Helper) CreateProject(name, language string) error {
	if name == "" {
		return fmt.Errorf("project name is required")
	}

	// Normalize language
	switch language {
	case "py":
		language = "python"
	case "golang":
		language = "go"
	case "ts", "node", "js":
		language = "typescript"
	case "":
		language = "general"
	}

	// Get settings template
	settings, ok := ProjectSettings[language]
	if !ok {
		settings = ProjectSettings["general"]
	}

	// Create project directory
	if err := os.MkdirAll(name, 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create .vscode directory
	vscodeDir := filepath.Join(name, ".vscode")
	if err := os.MkdirAll(vscodeDir, 0o755); err != nil {
		return fmt.Errorf("failed to create .vscode directory: %w", err)
	}

	// Write settings.json
	settingsPath := filepath.Join(vscodeDir, "settings.json")
	if err := os.WriteFile(settingsPath, []byte(settings), 0o644); err != nil {
		return fmt.Errorf("failed to write settings.json: %w", err)
	}

	return nil
}

// OpenProject opens a project in VS Code.
func (h *Helper) OpenProject(path string) error {
	return h.runCode(path)
}

// ListExtensions returns installed extensions.
func (h *Helper) ListExtensions() ([]Extension, error) {
	cmd := exec.Command("code", "--list-extensions")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list extensions: %w", err)
	}

	var extensions []Extension
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		id := strings.TrimSpace(scanner.Text())
		if id != "" {
			extensions = append(extensions, Extension{ID: id, Installed: true})
		}
	}

	return extensions, nil
}

// InstallExtension installs a VS Code extension.
func (h *Helper) InstallExtension(id string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would install extension: %s\n", id)
		return nil
	}

	if h.verbose {
		fmt.Printf("Installing: %s\n", id)
	}

	cmd := exec.Command("code", "--install-extension", id, "--force")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// InstallExtensionsFromFile installs extensions from a file.
func (h *Helper) InstallExtensionsFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		id := strings.TrimSpace(scanner.Text())
		if id == "" || strings.HasPrefix(id, "#") {
			continue
		}
		if err := h.InstallExtension(id); err != nil {
			fmt.Printf("Warning: failed to install %s: %v\n", id, err)
		}
	}

	return scanner.Err()
}

// InstallEssentialExtensions installs the essential extensions.
func (h *Helper) InstallEssentialExtensions() error {
	for _, id := range EssentialExtensions {
		if err := h.InstallExtension(id); err != nil {
			fmt.Printf("Warning: failed to install %s: %v\n", id, err)
		}
	}
	return nil
}

// ExportExtensions exports installed extensions to a file.
func (h *Helper) ExportExtensions(filePath string) error {
	extensions, err := h.ListExtensions()
	if err != nil {
		return err
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would export %d extensions to: %s\n", len(extensions), filePath)
		return nil
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	for _, ext := range extensions {
		fmt.Fprintln(file, ext.ID)
	}

	return nil
}

// SyncConfig syncs VS Code config from dotfiles.
func (h *Helper) SyncConfig() error {
	srcDir := filepath.Join(h.dotfilesRoot, "config", "vscode")
	paths := h.GetConfigPaths()

	// Check source directory
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		return fmt.Errorf("source config not found: %s", srcDir)
	}

	// Check destination directory
	if _, err := os.Stat(paths.UserDir); os.IsNotExist(err) {
		return fmt.Errorf("VS Code user directory not found: %s (is VS Code installed?)", paths.UserDir)
	}

	// Backup and sync settings.json
	if err := h.syncFile(
		filepath.Join(srcDir, "settings.json"),
		paths.Settings,
	); err != nil {
		return err
	}

	// Backup and sync keybindings.json if exists
	keySrc := filepath.Join(srcDir, "keybindings.json")
	if _, err := os.Stat(keySrc); err == nil {
		if err := h.syncFile(keySrc, paths.Keybindings); err != nil {
			return err
		}
	}

	return nil
}

// syncFile backs up and copies a file.
func (h *Helper) syncFile(src, dst string) error {
	// Backup existing
	if _, err := os.Stat(dst); err == nil {
		backupPath := dst + ".backup"
		if h.dryRun {
			fmt.Printf("[dry-run] would backup: %s -> %s\n", dst, backupPath)
		} else {
			data, err := os.ReadFile(dst)
			if err != nil {
				return fmt.Errorf("failed to read existing file: %w", err)
			}
			if err := os.WriteFile(backupPath, data, 0o644); err != nil {
				return fmt.Errorf("failed to backup: %w", err)
			}
			if h.verbose {
				fmt.Printf("Backed up: %s\n", filepath.Base(dst))
			}
		}
	}

	// Copy new file
	if h.dryRun {
		fmt.Printf("[dry-run] would sync: %s -> %s\n", src, dst)
		return nil
	}

	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read source: %w", err)
	}

	if err := os.WriteFile(dst, data, 0o644); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	if h.verbose {
		fmt.Printf("Synced: %s\n", filepath.Base(dst))
	}

	return nil
}

// runCode runs the VS Code command.
func (h *Helper) runCode(args ...string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would run: code %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("code", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// IsInstalled checks if VS Code is installed.
func IsInstalled() bool {
	_, err := exec.LookPath("code")
	return err == nil
}
