// Package posting provides Posting terminal HTTP client helper functionality.
package posting

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Status represents Posting installation status.
type Status struct {
	Installed  bool   `json:"installed" yaml:"installed"`
	Version    string `json:"version,omitempty" yaml:"version,omitempty"`
	ConfigPath string `json:"config_path,omitempty" yaml:"config_path,omitempty"`
	ThemesPath string `json:"themes_path,omitempty" yaml:"themes_path,omitempty"`
	DataPath   string `json:"data_path,omitempty" yaml:"data_path,omitempty"`
	Theme      string `json:"theme,omitempty" yaml:"theme,omitempty"`
	Layout     string `json:"layout,omitempty" yaml:"layout,omitempty"`
}

// Locations holds file and directory locations for posting.
type Locations struct {
	Config string `json:"config" yaml:"config"`
	Themes string `json:"themes" yaml:"themes"`
	Data   string `json:"data" yaml:"data"`
}

// Helper provides Posting helper operations.
type Helper struct {
	verbose bool
	dryRun  bool
}

// NewHelper creates a new Posting Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	return &Helper{
		verbose: verbose,
		dryRun:  dryRun,
	}
}

// IsInstalled checks if posting is installed.
func (h *Helper) IsInstalled() bool {
	// Check standard PATH first
	if _, err := exec.LookPath("posting"); err == nil {
		return true
	}

	// Check uv tools location
	home, _ := os.UserHomeDir()
	uvToolPath := filepath.Join(home, ".local", "bin", "posting")
	if _, err := os.Stat(uvToolPath); err == nil {
		return true
	}

	return false
}

// GetStatus returns Posting status information.
func (h *Helper) GetStatus() *Status {
	status := &Status{}

	if h.IsInstalled() {
		status.Installed = true
		status.Version = h.getVersion()
	}

	locations := h.GetLocations()
	status.ConfigPath = locations.Config
	status.ThemesPath = locations.Themes
	status.DataPath = locations.Data

	// Get current theme from env
	status.Theme = os.Getenv("POSTING_THEME")
	if status.Theme == "" {
		status.Theme = "galaxy"
	}

	status.Layout = os.Getenv("POSTING_LAYOUT")
	if status.Layout == "" {
		status.Layout = "horizontal"
	}

	return status
}

// getVersion returns the posting version.
func (h *Helper) getVersion() string {
	// Get version from uv tool list
	cmd := exec.Command("uv", "tool", "list")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	// Parse output like "posting v2.9.2"
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "posting ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return strings.TrimPrefix(parts[1], "v")
			}
		}
	}
	return ""
}

// GetLocations returns file and directory locations.
func (h *Helper) GetLocations() *Locations {
	home, _ := os.UserHomeDir()

	// XDG directories
	xdgConfig := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfig == "" {
		xdgConfig = filepath.Join(home, ".config")
	}

	xdgData := os.Getenv("XDG_DATA_HOME")
	if xdgData == "" {
		xdgData = filepath.Join(home, ".local", "share")
	}

	return &Locations{
		Config: filepath.Join(xdgConfig, "posting", "config.yaml"),
		Themes: filepath.Join(xdgConfig, "posting", "themes"),
		Data:   filepath.Join(xdgData, "posting"),
	}
}

// Launch starts posting with optional collection and env file.
func (h *Helper) Launch(collection, envFile string) error {
	if h.dryRun {
		args := []string{"posting"}
		if collection != "" {
			args = append(args, "--collection", collection)
		}
		if envFile != "" {
			args = append(args, "--env", envFile)
		}
		fmt.Printf("[dry-run] would run: %s\n", strings.Join(args, " "))
		return nil
	}

	if !h.IsInstalled() {
		return fmt.Errorf("posting not found - install with: uv tool install posting --python 3.13")
	}

	args := []string{}
	if collection != "" {
		args = append(args, "--collection", collection)
	}
	if envFile != "" {
		args = append(args, "--env", envFile)
	}

	cmd := exec.Command("posting", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// OpenConfig opens the config file in the default editor.
func (h *Helper) OpenConfig() error {
	locations := h.GetLocations()
	configPath := locations.Config

	if h.dryRun {
		fmt.Printf("[dry-run] would open: %s\n", configPath)
		return nil
	}

	// Check if config exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("config file not found at %s - run 'posting' once to create it", configPath)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	cmd := exec.Command(editor, configPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ListCollections returns directories in the data path.
func (h *Helper) ListCollections() ([]string, error) {
	locations := h.GetLocations()
	dataPath := locations.Data

	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		return []string{}, nil
	}

	entries, err := os.ReadDir(dataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read data directory: %w", err)
	}

	var collections []string
	for _, entry := range entries {
		if entry.IsDir() {
			collections = append(collections, filepath.Join(dataPath, entry.Name()))
		}
	}

	return collections, nil
}

// IsUVInstalled checks if uv is installed.
func (h *Helper) IsUVInstalled() bool {
	_, err := exec.LookPath("uv")
	return err == nil
}

// InstallUV installs uv package manager.
func (h *Helper) InstallUV() error {
	if h.IsUVInstalled() {
		return fmt.Errorf("uv already installed")
	}

	if h.dryRun {
		fmt.Println("[dry-run] would install uv via curl")
		return nil
	}

	fmt.Println("Installing uv...")
	cmd := exec.Command("bash", "-c", "curl -LsSf https://astral.sh/uv/install.sh | sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Install installs posting via uv.
func (h *Helper) Install() error {
	if h.IsInstalled() {
		return fmt.Errorf("posting already installed: %s", h.getVersion())
	}

	// Check if uv is installed
	if !h.IsUVInstalled() {
		return fmt.Errorf("uv not installed - run: acorn data posting install-uv")
	}

	if h.dryRun {
		fmt.Println("[dry-run] would run: uv tool install posting --python 3.13")
		return nil
	}

	fmt.Println("Installing posting...")
	cmd := exec.Command("uv", "tool", "install", "posting", "--python", "3.13")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Upgrade upgrades posting to the latest version.
func (h *Helper) Upgrade() error {
	if !h.IsInstalled() {
		return fmt.Errorf("posting not installed - run: acorn data posting install")
	}

	if !h.IsUVInstalled() {
		return fmt.Errorf("uv not installed")
	}

	if h.dryRun {
		fmt.Println("[dry-run] would run: uv tool upgrade posting")
		return nil
	}

	fmt.Println("Upgrading posting...")
	cmd := exec.Command("uv", "tool", "upgrade", "posting")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Uninstall removes posting.
func (h *Helper) Uninstall() error {
	if !h.IsInstalled() {
		return fmt.Errorf("posting not installed")
	}

	if !h.IsUVInstalled() {
		return fmt.Errorf("uv not installed")
	}

	if h.dryRun {
		fmt.Println("[dry-run] would run: uv tool uninstall posting")
		return nil
	}

	fmt.Println("Uninstalling posting...")
	cmd := exec.Command("uv", "tool", "uninstall", "posting")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
