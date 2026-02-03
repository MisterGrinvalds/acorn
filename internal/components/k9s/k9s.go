// Package k9s provides k9s helper functionality.
package k9s

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Status represents k9s installation status.
type Status struct {
	Installed        bool   `json:"installed" yaml:"installed"`
	Version          string `json:"version,omitempty" yaml:"version,omitempty"`
	Location         string `json:"location,omitempty" yaml:"location,omitempty"`
	ConfigDir        string `json:"config_dir,omitempty" yaml:"config_dir,omitempty"`
	DataDir          string `json:"data_dir,omitempty" yaml:"data_dir,omitempty"`
	ClusterConnected bool   `json:"cluster_connected" yaml:"cluster_connected"`
	CurrentContext   string `json:"current_context,omitempty" yaml:"current_context,omitempty"`
}

// Helper provides k9s helper operations.
type Helper struct {
	verbose bool
}

// NewHelper creates a new k9s Helper.
func NewHelper(verbose bool) *Helper {
	return &Helper{
		verbose: verbose,
	}
}

// IsInstalled checks if k9s is installed.
func (h *Helper) IsInstalled() bool {
	_, err := exec.LookPath("k9s")
	return err == nil
}

// IsClusterConnected checks if kubectl can connect to a cluster.
func (h *Helper) IsClusterConnected() bool {
	cmd := exec.Command("kubectl", "cluster-info")
	return cmd.Run() == nil
}

// GetCurrentContext returns the current kubectl context.
func (h *Helper) GetCurrentContext() string {
	cmd := exec.Command("kubectl", "config", "current-context")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// GetStatus returns k9s installation status.
func (h *Helper) GetStatus() *Status {
	status := &Status{}

	k9sPath, err := exec.LookPath("k9s")
	if err != nil {
		return status
	}

	status.Installed = true
	status.Location = k9sPath

	// Get version
	cmd := exec.Command("k9s", "version", "-s")
	out, err := cmd.Output()
	if err == nil {
		// Format: Version: v0.32.4
		for line := range strings.SplitSeq(string(out), "\n") {
			if version, found := strings.CutPrefix(line, "Version:"); found {
				status.Version = strings.TrimSpace(version)
				break
			}
		}
		if status.Version == "" {
			status.Version = strings.TrimSpace(string(out))
		}
	}

	// Get config directories
	status.ConfigDir = h.GetConfigDir()
	status.DataDir = h.GetDataDir()

	// Check cluster connection
	status.ClusterConnected = h.IsClusterConnected()
	if status.ClusterConnected {
		status.CurrentContext = h.GetCurrentContext()
	}

	return status
}

// GetConfigDir returns the k9s config directory.
func (h *Helper) GetConfigDir() string {
	// Check environment variable first
	if dir := os.Getenv("K9S_CONFIG_DIR"); dir != "" {
		return dir
	}

	// XDG config home
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, "k9s")
	}

	// Default
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "k9s")
}

// GetDataDir returns the k9s data directory.
func (h *Helper) GetDataDir() string {
	// Check environment variable first
	if dir := os.Getenv("K9S_DATA_DIR"); dir != "" {
		return dir
	}

	// XDG data home
	if xdgData := os.Getenv("XDG_DATA_HOME"); xdgData != "" {
		return filepath.Join(xdgData, "k9s")
	}

	// Default
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "k9s")
}

// GetConfigFile returns the path to the main config file.
func (h *Helper) GetConfigFile() string {
	return filepath.Join(h.GetConfigDir(), "config.yaml")
}

// GetPluginsFile returns the path to the plugins file.
func (h *Helper) GetPluginsFile() string {
	return filepath.Join(h.GetConfigDir(), "plugins.yaml")
}

// GetHotkeysFile returns the path to the hotkeys file.
func (h *Helper) GetHotkeysFile() string {
	return filepath.Join(h.GetConfigDir(), "hotkeys.yaml")
}

// GetViewsFile returns the path to the views file.
func (h *Helper) GetViewsFile() string {
	return filepath.Join(h.GetConfigDir(), "views.yaml")
}

// GetAliasesFile returns the path to the aliases file.
func (h *Helper) GetAliasesFile() string {
	return filepath.Join(h.GetConfigDir(), "aliases.yaml")
}

// GetSkinsDir returns the path to the skins directory.
func (h *Helper) GetSkinsDir() string {
	return filepath.Join(h.GetConfigDir(), "skins")
}

// ConfigExists checks if the main config file exists.
func (h *Helper) ConfigExists() bool {
	_, err := os.Stat(h.GetConfigFile())
	return err == nil
}

// Launch starts k9s with optional arguments.
func (h *Helper) Launch(context, namespace, command string, readonly, headless bool) error {
	if !h.IsInstalled() {
		return fmt.Errorf("k9s is not installed")
	}

	if !h.IsClusterConnected() {
		return fmt.Errorf("cannot connect to Kubernetes cluster")
	}

	args := []string{}
	if context != "" {
		args = append(args, "--context", context)
	}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	if command != "" {
		args = append(args, "-c", command)
	}
	if readonly {
		args = append(args, "--readonly")
	}
	if headless {
		args = append(args, "--headless", "--logoless", "--crumbsless")
	}

	cmd := exec.Command("k9s", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// OpenConfig opens a k9s config file in the default editor.
func (h *Helper) OpenConfig(configType string) error {
	var configFile string
	var defaultContent string

	switch configType {
	case "main", "config":
		configFile = h.GetConfigFile()
		defaultContent = `# k9s configuration
# See: https://k9scli.io/topics/config/
k9s:
  refreshRate: 2
  headless: false
  logoless: false
  crumbsless: false
  readOnly: false
  noExitOnCtrlC: false
  ui:
    skin: ""
    enableMouse: false
  skipLatestRevCheck: false
`
	case "plugins":
		configFile = h.GetPluginsFile()
		defaultContent = `# k9s plugins
# See: https://k9scli.io/topics/plugins/
plugins: {}
`
	case "hotkeys":
		configFile = h.GetHotkeysFile()
		defaultContent = `# k9s hotkeys
# See: https://k9scli.io/topics/hotkeys/
hotKeys: {}
`
	case "views":
		configFile = h.GetViewsFile()
		defaultContent = `# k9s custom views
# See: https://k9scli.io/topics/columns/
views: {}
`
	case "aliases":
		configFile = h.GetAliasesFile()
		defaultContent = `# k9s command aliases
# See: https://k9scli.io/topics/aliases/
aliases: {}
`
	default:
		return fmt.Errorf("unknown config type: %s (valid: main, plugins, hotkeys, views, aliases)", configType)
	}

	configDir := h.GetConfigDir()

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create config file if it doesn't exist
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := os.WriteFile(configFile, []byte(defaultContent), 0o644); err != nil {
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

// ListSkins lists available skins.
func (h *Helper) ListSkins() ([]string, error) {
	skinsDir := h.GetSkinsDir()

	if _, err := os.Stat(skinsDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	entries, err := os.ReadDir(skinsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read skins directory: %w", err)
	}

	var skins []string
	for _, entry := range entries {
		if !entry.IsDir() && (strings.HasSuffix(entry.Name(), ".yaml") || strings.HasSuffix(entry.Name(), ".yml")) {
			name := strings.TrimSuffix(strings.TrimSuffix(entry.Name(), ".yaml"), ".yml")
			skins = append(skins, name)
		}
	}

	return skins, nil
}

// GetContexts returns available kubectl contexts.
func (h *Helper) GetContexts() ([]string, error) {
	cmd := exec.Command("kubectl", "config", "get-contexts", "-o", "name")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get contexts: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var contexts []string
	for _, line := range lines {
		if line = strings.TrimSpace(line); line != "" {
			contexts = append(contexts, line)
		}
	}
	return contexts, nil
}

// GetNamespaces returns available namespaces.
func (h *Helper) GetNamespaces() ([]string, error) {
	cmd := exec.Command("kubectl", "get", "namespaces", "-o", "jsonpath={.items[*].metadata.name}")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get namespaces: %w", err)
	}

	parts := strings.Split(strings.TrimSpace(string(out)), " ")
	var namespaces []string
	for _, ns := range parts {
		if ns = strings.TrimSpace(ns); ns != "" {
			namespaces = append(namespaces, ns)
		}
	}
	return namespaces, nil
}

// GetKeybindings returns the default keybindings.
func (h *Helper) GetKeybindings() string {
	return `k9s Keybindings
===============

Navigation:
  ?           Help/keybindings
  :           Command mode
  /           Filter mode
  esc         Exit view/mode
  ctrl-a      All resource aliases
  -           Go back
  [/]         Navigate command history

Commands:
  :pod        View pods
  :dp         View deployments
  :svc        View services
  :ctx        Switch context
  :ns         Switch namespace
  :pulses     Cluster health
  :xray       Resource graph
  :pf         Port-forwards

Actions:
  d           Describe resource
  v/y         View YAML
  e           Edit resource
  l           View logs
  s           Shell into container
  f           Port-forward
  ctrl-d      Delete (confirm)
  ctrl-k      Kill (no confirm)

Logs:
  w           Toggle line wrap
  t           Toggle timestamps
  a           Toggle autoscroll
  p           Previous container logs
  c           Clear log buffer
  ctrl-s      Save logs to file

Filters:
  /text       Regex filter
  /!text      Inverse filter
  /-l label   Label selector
  /-f text    Fuzzy find
`
}

// Keybinding represents a k9s keybinding.
type Keybinding struct {
	Key         string `json:"key" yaml:"key"`
	Description string `json:"description" yaml:"description"`
	Category    string `json:"category" yaml:"category"`
}

// GetKeybindingsList returns keybindings as a structured list.
func (h *Helper) GetKeybindingsList() []Keybinding {
	return []Keybinding{
		// Navigation
		{Key: "?", Description: "Help/keybindings", Category: "Navigation"},
		{Key: ":", Description: "Command mode", Category: "Navigation"},
		{Key: "/", Description: "Filter mode", Category: "Navigation"},
		{Key: "esc", Description: "Exit view/mode", Category: "Navigation"},
		{Key: "ctrl-a", Description: "All resource aliases", Category: "Navigation"},
		{Key: "-", Description: "Go back", Category: "Navigation"},
		{Key: "[/]", Description: "Navigate command history", Category: "Navigation"},

		// Commands
		{Key: ":pod", Description: "View pods", Category: "Commands"},
		{Key: ":dp", Description: "View deployments", Category: "Commands"},
		{Key: ":svc", Description: "View services", Category: "Commands"},
		{Key: ":ctx", Description: "Switch context", Category: "Commands"},
		{Key: ":ns", Description: "Switch namespace", Category: "Commands"},
		{Key: ":pulses", Description: "Cluster health", Category: "Commands"},
		{Key: ":xray", Description: "Resource graph", Category: "Commands"},
		{Key: ":pf", Description: "Port-forwards", Category: "Commands"},

		// Actions
		{Key: "d", Description: "Describe resource", Category: "Actions"},
		{Key: "v/y", Description: "View YAML", Category: "Actions"},
		{Key: "e", Description: "Edit resource", Category: "Actions"},
		{Key: "l", Description: "View logs", Category: "Actions"},
		{Key: "s", Description: "Shell into container", Category: "Actions"},
		{Key: "f", Description: "Port-forward", Category: "Actions"},
		{Key: "ctrl-d", Description: "Delete (confirm)", Category: "Actions"},
		{Key: "ctrl-k", Description: "Kill (no confirm)", Category: "Actions"},

		// Logs
		{Key: "w", Description: "Toggle line wrap", Category: "Logs"},
		{Key: "t", Description: "Toggle timestamps", Category: "Logs"},
		{Key: "a", Description: "Toggle autoscroll", Category: "Logs"},
		{Key: "p", Description: "Previous container logs", Category: "Logs"},
		{Key: "c", Description: "Clear log buffer", Category: "Logs"},
		{Key: "ctrl-s", Description: "Save logs to file", Category: "Logs"},
	}
}

// Info represents k9s info output.
type Info struct {
	ConfigDir  string `json:"config_dir" yaml:"config_dir"`
	LogsDir    string `json:"logs_dir" yaml:"logs_dir"`
	ScreenDump string `json:"screen_dump" yaml:"screen_dump"`
}

// GetInfo returns k9s info.
func (h *Helper) GetInfo() (*Info, error) {
	cmd := exec.Command("k9s", "info")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get k9s info: %w", err)
	}

	info := &Info{}
	for line := range strings.SplitSeq(string(out), "\n") {
		key, value, found := strings.Cut(line, ":")
		if !found {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		switch key {
		case "Configuration":
			info.ConfigDir = value
		case "Logs":
			info.LogsDir = value
		case "Screen Dumps":
			info.ScreenDump = value
		}
	}

	return info, nil
}
