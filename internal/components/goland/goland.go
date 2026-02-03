// Package goland provides JetBrains GoLand IDE helper functionality.
package goland

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Status represents GoLand installation status.
type Status struct {
	Installed bool   `json:"installed" yaml:"installed"`
	Version   string `json:"version,omitempty" yaml:"version,omitempty"`
	AppPath   string `json:"app_path,omitempty" yaml:"app_path,omitempty"`
	ConfigDir string `json:"config_dir,omitempty" yaml:"config_dir,omitempty"`
	HasCLI    bool   `json:"has_cli" yaml:"has_cli"`
	CLIPath   string `json:"cli_path,omitempty" yaml:"cli_path,omitempty"`
}

// Project represents a recent GoLand project.
type Project struct {
	Name   string `json:"name" yaml:"name"`
	Path   string `json:"path" yaml:"path"`
	GoMod  bool   `json:"go_mod" yaml:"go_mod"`
}

// Helper provides GoLand helper operations.
type Helper struct {
	verbose bool
	dryRun  bool
}

// NewHelper creates a new GoLand Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	return &Helper{
		verbose: verbose,
		dryRun:  dryRun,
	}
}

// GetStatus returns GoLand status information.
func (h *Helper) GetStatus() *Status {
	status := &Status{}

	// Find GoLand installation
	appPath := h.findGoLand()
	if appPath != "" {
		status.Installed = true
		status.AppPath = appPath
		status.Version = h.getVersion(appPath)
	}

	// Check for CLI tool
	cliPath := h.findCLI()
	if cliPath != "" {
		status.HasCLI = true
		status.CLIPath = cliPath
	}

	// Get config directory
	status.ConfigDir = h.getConfigDir()

	return status
}

// findGoLand finds the GoLand application.
func (h *Helper) findGoLand() string {
	if runtime.GOOS == "darwin" {
		// Check common macOS locations
		paths := []string{
			"/Applications/GoLand.app",
			filepath.Join(os.Getenv("HOME"), "Applications/GoLand.app"),
			"/Applications/JetBrains Toolbox/GoLand.app",
		}

		// Also check for versioned installations via JetBrains Toolbox
		toolboxApps := filepath.Join(os.Getenv("HOME"), "Library/Application Support/JetBrains/Toolbox/apps/goland")
		if entries, err := os.ReadDir(toolboxApps); err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					channelPath := filepath.Join(toolboxApps, entry.Name())
					if subEntries, err := os.ReadDir(channelPath); err == nil {
						for _, subEntry := range subEntries {
							appPath := filepath.Join(channelPath, subEntry.Name(), "GoLand.app")
							if _, err := os.Stat(appPath); err == nil {
								return appPath
							}
						}
					}
				}
			}
		}

		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}
	}

	return ""
}

// findCLI finds the GoLand CLI tool.
func (h *Helper) findCLI() string {
	// Check if 'goland' command is in PATH
	if path, err := exec.LookPath("goland"); err == nil {
		return path
	}

	// Check common CLI locations on macOS
	if runtime.GOOS == "darwin" {
		paths := []string{
			"/usr/local/bin/goland",
			filepath.Join(os.Getenv("HOME"), ".local/bin/goland"),
		}
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}
	}

	return ""
}

// getVersion extracts GoLand version from the app.
func (h *Helper) getVersion(appPath string) string {
	if runtime.GOOS == "darwin" {
		// Read version from Info.plist
		plistPath := filepath.Join(appPath, "Contents/Info.plist")
		cmd := exec.Command("defaults", "read", plistPath, "CFBundleShortVersionString")
		out, err := cmd.Output()
		if err == nil {
			return strings.TrimSpace(string(out))
		}
	}
	return ""
}

// getConfigDir returns the GoLand config directory.
func (h *Helper) getConfigDir() string {
	home, _ := os.UserHomeDir()

	if runtime.GOOS == "darwin" {
		// Check for JetBrains Toolbox managed config
		configBase := filepath.Join(home, "Library/Application Support/JetBrains")
		if entries, err := os.ReadDir(configBase); err == nil {
			for _, entry := range entries {
				if strings.HasPrefix(entry.Name(), "GoLand") {
					return filepath.Join(configBase, entry.Name())
				}
			}
		}
	}

	return ""
}

// Launch starts GoLand.
func (h *Helper) Launch(args ...string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would launch GoLand with args: %v\n", args)
		return nil
	}

	// Try CLI first
	if cli := h.findCLI(); cli != "" {
		cmd := exec.Command(cli, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Start() // Don't wait
	}

	// Fall back to open command on macOS
	if runtime.GOOS == "darwin" {
		appPath := h.findGoLand()
		if appPath == "" {
			return fmt.Errorf("GoLand not found")
		}

		openArgs := []string{"-a", appPath}
		if len(args) > 0 {
			openArgs = append(openArgs, args...)
		}

		cmd := exec.Command("open", openArgs...)
		return cmd.Start()
	}

	return fmt.Errorf("GoLand not found")
}

// OpenProject opens a project in GoLand.
func (h *Helper) OpenProject(path string) error {
	// Resolve path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", absPath)
	}

	return h.Launch(absPath)
}

// OpenFile opens a specific file in GoLand.
func (h *Helper) OpenFile(filePath string, line int) error {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	args := []string{absPath}
	if line > 0 {
		args = []string{"--line", fmt.Sprintf("%d", line), absPath}
	}

	return h.Launch(args...)
}

// CreateCLILink creates the command-line launcher.
func (h *Helper) CreateCLILink() error {
	if h.dryRun {
		fmt.Println("[dry-run] would create GoLand CLI launcher")
		return nil
	}

	appPath := h.findGoLand()
	if appPath == "" {
		return fmt.Errorf("GoLand not found")
	}

	if runtime.GOOS == "darwin" {
		// The CLI script is inside the app bundle
		cliScript := filepath.Join(appPath, "Contents/MacOS/goland")
		if _, err := os.Stat(cliScript); os.IsNotExist(err) {
			return fmt.Errorf("CLI script not found in app bundle")
		}

		// Create symlink in /usr/local/bin
		linkPath := "/usr/local/bin/goland"

		// Remove existing link if present
		os.Remove(linkPath)

		if err := os.Symlink(cliScript, linkPath); err != nil {
			return fmt.Errorf("failed to create symlink (try with sudo): %w", err)
		}

		fmt.Printf("Created: %s -> %s\n", linkPath, cliScript)
		return nil
	}

	return fmt.Errorf("CLI link creation not supported on this platform")
}

// IsGoProject checks if a directory is a Go project.
func (h *Helper) IsGoProject(path string) bool {
	goMod := filepath.Join(path, "go.mod")
	if _, err := os.Stat(goMod); err == nil {
		return true
	}

	// Check for .go files
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
			return true
		}
	}

	return false
}

// ListRecentProjects lists recent GoLand projects.
func (h *Helper) ListRecentProjects() ([]Project, error) {
	configDir := h.getConfigDir()
	if configDir == "" {
		return nil, fmt.Errorf("config directory not found")
	}

	// Recent projects are stored in recentProjects.xml
	recentFile := filepath.Join(configDir, "options/recentProjects.xml")
	if _, err := os.Stat(recentFile); os.IsNotExist(err) {
		return []Project{}, nil
	}

	// Parse would require XML parsing - return empty for now
	return []Project{}, nil
}
