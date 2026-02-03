// Package lazydocker provides Lazydocker helper functionality.
package lazydocker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Status represents Lazydocker installation status.
type Status struct {
	Installed     bool   `json:"installed" yaml:"installed"`
	Version       string `json:"version,omitempty" yaml:"version,omitempty"`
	Location      string `json:"location,omitempty" yaml:"location,omitempty"`
	ConfigDir     string `json:"config_dir,omitempty" yaml:"config_dir,omitempty"`
	DockerRunning bool   `json:"docker_running" yaml:"docker_running"`
}

// Helper provides Lazydocker helper operations.
type Helper struct {
	verbose bool
}

// NewHelper creates a new Lazydocker Helper.
func NewHelper(verbose bool) *Helper {
	return &Helper{
		verbose: verbose,
	}
}

// IsInstalled checks if Lazydocker is installed.
func (h *Helper) IsInstalled() bool {
	_, err := exec.LookPath("lazydocker")
	return err == nil
}

// IsDockerRunning checks if Docker daemon is running.
func (h *Helper) IsDockerRunning() bool {
	cmd := exec.Command("docker", "info")
	return cmd.Run() == nil
}

// GetStatus returns Lazydocker installation status.
func (h *Helper) GetStatus() *Status {
	status := &Status{}

	lzdPath, err := exec.LookPath("lazydocker")
	if err != nil {
		return status
	}

	status.Installed = true
	status.Location = lzdPath

	// Get version
	cmd := exec.Command("lazydocker", "--version")
	out, err := cmd.Output()
	if err == nil {
		// Format: "Version: 0.23.1\nCommit: ..."
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "Version:") {
				status.Version = strings.TrimSpace(strings.TrimPrefix(line, "Version:"))
				break
			}
		}
		if status.Version == "" {
			status.Version = strings.TrimSpace(string(out))
		}
	}

	// Get config directory
	status.ConfigDir = h.GetConfigDir()

	// Check Docker status
	status.DockerRunning = h.IsDockerRunning()

	return status
}

// GetConfigDir returns the Lazydocker config directory.
func (h *Helper) GetConfigDir() string {
	// Check environment variable first
	if dir := os.Getenv("LAZYDOCKER_CONFIG_DIR"); dir != "" {
		return dir
	}

	// XDG config home
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, "lazydocker")
	}

	// Default
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "lazydocker")
}

// GetConfigFile returns the path to the config file.
func (h *Helper) GetConfigFile() string {
	return filepath.Join(h.GetConfigDir(), "config.yml")
}

// ConfigExists checks if the config file exists.
func (h *Helper) ConfigExists() bool {
	_, err := os.Stat(h.GetConfigFile())
	return err == nil
}

// Launch starts Lazydocker.
func (h *Helper) Launch(configFile string) error {
	if !h.IsInstalled() {
		return fmt.Errorf("lazydocker is not installed")
	}

	if !h.IsDockerRunning() {
		return fmt.Errorf("docker is not running")
	}

	args := []string{}
	if configFile != "" {
		args = append(args, "--config", configFile)
	}

	cmd := exec.Command("lazydocker", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// OpenConfig opens the config file in the default editor.
func (h *Helper) OpenConfig() error {
	configFile := h.GetConfigFile()
	configDir := h.GetConfigDir()

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create config file if it doesn't exist
	if !h.ConfigExists() {
		defaultConfig := `# Lazydocker configuration
# See: https://github.com/jesseduffield/lazydocker/blob/master/docs/Config.md

gui:
  # Scroll height in the main panels
  scrollHeight: 2

  # Show bottom line (docker stats)
  showBottomLine: true

  # Whether to show all containers (including stopped ones)
  showAllContainers: false

logs:
  # Enable timestamps in logs
  timestamps: false

  # Tail initial number of lines
  tail: "300"

commandTemplates:
  # Custom commands that appear in the custom command menu
  # dockerCompose: docker-compose
  # restartService: docker-compose restart {{ .Service.Name }}
`
		if err := os.WriteFile(configFile, []byte(defaultConfig), 0o644); err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	cmd := exec.Command(editor, configFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetKeybindings returns the default keybindings.
func (h *Helper) GetKeybindings() string {
	return `Lazydocker Keybindings
======================

Navigation:
  ↑/↓, j/k    Navigate items
  ←/→, h/l    Switch panels
  Tab         Next panel
  [/]         Previous/next tab
  g/G         Top/bottom of list
  /           Filter

Actions:
  Enter       Focus panel / Execute
  Space       Select item
  d           Remove (with confirmation)
  D           Remove (force, no confirm)
  s           Stop container
  r           Restart container
  a           Attach to container
  m           View logs
  E           Exec shell in container
  b           View bulk commands
  c           Run custom command

Views:
  e           Hide/show stopped containers
  x           Open menu
  ?           Help
  q, Esc      Quit / Back

Containers Panel:
  Enter       Enter container menu
  [           Previous container tab
  ]           Next container tab

Logs:
  Ctrl+s      Save logs to file
  f           Follow logs

Other:
  +           Increase log window
  -           Decrease log window
  Ctrl+r      Refresh
`
}

// Keybinding represents a lazydocker keybinding.
type Keybinding struct {
	Key         string `json:"key" yaml:"key"`
	Description string `json:"description" yaml:"description"`
	Category    string `json:"category" yaml:"category"`
}

// GetKeybindingsList returns keybindings as a structured list.
func (h *Helper) GetKeybindingsList() []Keybinding {
	return []Keybinding{
		// Navigation
		{Key: "↑/↓, j/k", Description: "Navigate items", Category: "Navigation"},
		{Key: "←/→, h/l", Description: "Switch panels", Category: "Navigation"},
		{Key: "Tab", Description: "Next panel", Category: "Navigation"},
		{Key: "[/]", Description: "Previous/next tab", Category: "Navigation"},
		{Key: "g/G", Description: "Top/bottom of list", Category: "Navigation"},
		{Key: "/", Description: "Filter", Category: "Navigation"},

		// Actions
		{Key: "Enter", Description: "Focus panel / Execute", Category: "Actions"},
		{Key: "Space", Description: "Select item", Category: "Actions"},
		{Key: "d", Description: "Remove (with confirmation)", Category: "Actions"},
		{Key: "D", Description: "Remove (force)", Category: "Actions"},
		{Key: "s", Description: "Stop container", Category: "Actions"},
		{Key: "r", Description: "Restart container", Category: "Actions"},
		{Key: "a", Description: "Attach to container", Category: "Actions"},
		{Key: "m", Description: "View logs", Category: "Actions"},
		{Key: "E", Description: "Exec shell in container", Category: "Actions"},
		{Key: "b", Description: "View bulk commands", Category: "Actions"},
		{Key: "c", Description: "Run custom command", Category: "Actions"},

		// Views
		{Key: "e", Description: "Hide/show stopped containers", Category: "Views"},
		{Key: "x", Description: "Open menu", Category: "Views"},
		{Key: "?", Description: "Help", Category: "Views"},
		{Key: "q, Esc", Description: "Quit / Back", Category: "Views"},

		// Logs
		{Key: "Ctrl+s", Description: "Save logs to file", Category: "Logs"},
		{Key: "f", Description: "Follow logs", Category: "Logs"},
		{Key: "+/-", Description: "Resize log window", Category: "Logs"},
		{Key: "Ctrl+r", Description: "Refresh", Category: "Logs"},
	}
}
