// Package btop provides btop++ resource monitor helper functionality.
package btop

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Status represents btop installation status.
type Status struct {
	Installed    bool   `json:"installed" yaml:"installed"`
	Version      string `json:"version,omitempty" yaml:"version,omitempty"`
	ConfigDir    string `json:"config_dir" yaml:"config_dir"`
	ConfigExists bool   `json:"config_exists" yaml:"config_exists"`
	ThemesDir    string `json:"themes_dir" yaml:"themes_dir"`
	ThemeCount   int    `json:"theme_count" yaml:"theme_count"`
}

// Theme represents a btop theme.
type Theme struct {
	Name     string `json:"name" yaml:"name"`
	Path     string `json:"path" yaml:"path"`
	IsCustom bool   `json:"is_custom" yaml:"is_custom"`
}

// Helper provides btop helper operations.
type Helper struct {
	verbose   bool
	dryRun    bool
	configDir string
}

// NewHelper creates a new btop Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		home, _ := os.UserHomeDir()
		configDir = filepath.Join(home, ".config")
	}
	configDir = filepath.Join(configDir, "btop")

	return &Helper{
		verbose:   verbose,
		dryRun:    dryRun,
		configDir: configDir,
	}
}

// GetConfigDir returns the btop config directory.
func (h *Helper) GetConfigDir() string {
	return h.configDir
}

// GetStatus returns btop status information.
func (h *Helper) GetStatus() *Status {
	status := &Status{
		ConfigDir: h.configDir,
		ThemesDir: filepath.Join(h.configDir, "themes"),
	}

	// Check if btop is installed
	out, err := exec.Command("btop", "--version").Output()
	if err != nil {
		status.Installed = false
		return status
	}

	status.Installed = true
	// Parse version (format: btop version: x.x.x)
	versionStr := strings.TrimSpace(string(out))
	if strings.Contains(versionStr, "version:") {
		parts := strings.Split(versionStr, "version:")
		if len(parts) >= 2 {
			status.Version = strings.TrimSpace(parts[1])
		}
	} else {
		status.Version = versionStr
	}

	// Check if config exists
	configFile := filepath.Join(h.configDir, "btop.conf")
	if _, err := os.Stat(configFile); err == nil {
		status.ConfigExists = true
	}

	// Count themes
	themes, _ := h.ListThemes()
	status.ThemeCount = len(themes)

	return status
}

// ListThemes lists available btop themes.
func (h *Helper) ListThemes() ([]Theme, error) {
	var themes []Theme

	// Check custom themes directory
	customDir := filepath.Join(h.configDir, "themes")
	if entries, err := os.ReadDir(customDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".theme") {
				themes = append(themes, Theme{
					Name:     strings.TrimSuffix(entry.Name(), ".theme"),
					Path:     filepath.Join(customDir, entry.Name()),
					IsCustom: true,
				})
			}
		}
	}

	// Check system themes (common locations)
	systemDirs := []string{
		"/usr/share/btop/themes",
		"/usr/local/share/btop/themes",
		"/opt/homebrew/share/btop/themes",
	}

	for _, dir := range systemDirs {
		if entries, err := os.ReadDir(dir); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".theme") {
					themes = append(themes, Theme{
						Name:     strings.TrimSuffix(entry.Name(), ".theme"),
						Path:     filepath.Join(dir, entry.Name()),
						IsCustom: false,
					})
				}
			}
		}
	}

	return themes, nil
}

// Launch starts btop.
func (h *Helper) Launch(args ...string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would run: btop %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("btop", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Install installs btop.
func (h *Helper) Install() error {
	if _, err := exec.LookPath("btop"); err == nil {
		return fmt.Errorf("btop is already installed")
	}

	if h.dryRun {
		fmt.Println("[dry-run] would install btop via homebrew")
		return nil
	}

	fmt.Println("Installing btop...")
	cmd := exec.Command("brew", "install", "btop")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// InitConfig initializes btop configuration.
func (h *Helper) InitConfig() error {
	dirs := []string{
		h.configDir,
		filepath.Join(h.configDir, "themes"),
	}

	for _, dir := range dirs {
		if h.dryRun {
			fmt.Printf("[dry-run] would create: %s\n", dir)
			continue
		}
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create %s: %w", dir, err)
		}
	}

	return nil
}

// GenerateConfig generates the btop.conf configuration file.
func (h *Helper) GenerateConfig(dryRun bool) error {
	configFile := filepath.Join(h.configDir, "btop.conf")

	// Check if config already exists
	if _, err := os.Stat(configFile); err == nil && !dryRun {
		fmt.Printf("Config exists: %s (skipping)\n", configFile)
		return nil
	}

	content := h.generateDefaultConfig()

	if dryRun {
		fmt.Printf("Would create: %s\n", configFile)
		fmt.Println("\n--- Preview ---")
		fmt.Println(string(content))
		fmt.Println("--- End Preview ---")
		return nil
	}

	// Ensure config directory exists
	if err := os.MkdirAll(h.configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(configFile, content, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	fmt.Printf("Created: %s\n", configFile)
	return nil
}

// generateDefaultConfig generates default btop.conf content.
func (h *Helper) generateDefaultConfig() []byte {
	var b strings.Builder

	b.WriteString("#? Config file for btop v. 1.3.0\n")
	b.WriteString("# Generated by acorn - customize as needed\n\n")

	b.WriteString("#* Name of a btop++/bpytop/bashtop formatted \".theme\" file\n")
	b.WriteString("#* Themes should be placed in \"../share/btop/themes\" or \"$HOME/.config/btop/themes\"\n")
	b.WriteString("color_theme = \"Default\"\n\n")

	b.WriteString("#* If the theme set background should be shown, set to False for terminal default\n")
	b.WriteString("theme_background = True\n\n")

	b.WriteString("#* Sets if 24-bit truecolor should be used\n")
	b.WriteString("truecolor = True\n\n")

	b.WriteString("#* Set to true to force tty mode regardless of detected mode\n")
	b.WriteString("force_tty = False\n\n")

	b.WriteString("#* Define presets for the layout of the boxes\n")
	b.WriteString("presets = \"cpu:1:default,proc:0:default cpu:0:default,mem:0:default,net:0:default cpu:0:block,net:0:tty\"\n\n")

	b.WriteString("#* Set to True to enable \"h,j,k,l,g,G\" keys for navigation\n")
	b.WriteString("vim_keys = True\n\n")

	b.WriteString("#* Rounded corners on boxes\n")
	b.WriteString("rounded_corners = True\n\n")

	b.WriteString("#* Default symbols to use for graph creation\n")
	b.WriteString("graph_symbol = \"braille\"\n\n")

	b.WriteString("#* Manually set which boxes to show\n")
	b.WriteString("shown_boxes = \"cpu mem net proc\"\n\n")

	b.WriteString("#* Update time in milliseconds\n")
	b.WriteString("update_ms = 1000\n\n")

	b.WriteString("#* Processes sorting (pid, program, arguments, threads, user, memory, cpu)\n")
	b.WriteString("proc_sorting = \"cpu lazy\"\n\n")

	b.WriteString("#* Reverse sorting order (True or False)\n")
	b.WriteString("proc_reversed = False\n\n")

	b.WriteString("#* Show processes as a tree\n")
	b.WriteString("proc_tree = False\n\n")

	b.WriteString("#* Use the cpu graph colors in the process list\n")
	b.WriteString("proc_colors = True\n\n")

	b.WriteString("#* Use a darkening gradient in the process list\n")
	b.WriteString("proc_gradient = True\n\n")

	b.WriteString("#* If process cpu usage should be of the core it's running on or usage of total\n")
	b.WriteString("proc_per_core = True\n\n")

	b.WriteString("#* Show process memory as bytes instead of percent\n")
	b.WriteString("proc_mem_bytes = True\n\n")

	b.WriteString("#* Show cpu graph for each process\n")
	b.WriteString("proc_cpu_graphs = True\n\n")

	b.WriteString("#* Sets the CPU stat shown in upper half of the CPU graph\n")
	b.WriteString("cpu_graph_upper = \"total\"\n\n")

	b.WriteString("#* Sets the CPU stat shown in lower half of the CPU graph\n")
	b.WriteString("cpu_graph_lower = \"total\"\n\n")

	b.WriteString("#* If gpu info should be shown in the cpu box\n")
	b.WriteString("show_gpu_info = \"Auto\"\n\n")

	b.WriteString("#* Toggles if the lower CPU graph should be inverted\n")
	b.WriteString("cpu_invert_lower = True\n\n")

	b.WriteString("#* Set to True to completely disable the lower CPU graph\n")
	b.WriteString("cpu_single_graph = False\n\n")

	b.WriteString("#* Show cpu box at bottom of screen instead of top\n")
	b.WriteString("cpu_bottom = False\n\n")

	b.WriteString("#* Shows the system uptime in the CPU box\n")
	b.WriteString("show_uptime = True\n\n")

	b.WriteString("#* Show cpu temperature\n")
	b.WriteString("check_temp = True\n\n")

	b.WriteString("#* Which sensor to use for cpu temp (Auto = auto detect)\n")
	b.WriteString("cpu_sensor = \"Auto\"\n\n")

	b.WriteString("#* Show temps for cpu cores (only works if check_temp is True)\n")
	b.WriteString("show_coretemp = True\n\n")

	b.WriteString("#* Set a custom min value for the cpu temp graph\n")
	b.WriteString("cpu_core_map = \"\"\n\n")

	b.WriteString("#* Which temperature scale to use (celsius, fahrenheit, kelvin, rankine)\n")
	b.WriteString("temp_scale = \"celsius\"\n\n")

	b.WriteString("#* Use base 10 for bits/bytes (True) or base 2 (False)\n")
	b.WriteString("base_10_sizes = False\n\n")

	b.WriteString("#* Show battery statistics in the top right corner\n")
	b.WriteString("show_battery = True\n\n")

	b.WriteString("#* Which battery to show (default auto-detects)\n")
	b.WriteString("selected_battery = \"Auto\"\n\n")

	b.WriteString("#* Set to true to show percentage of allocated disk space that is used\n")
	b.WriteString("show_io_stat = True\n\n")

	b.WriteString("#* Set to true to show percentage of total disk space available\n")
	b.WriteString("io_mode = False\n\n")

	b.WriteString("#* Set to true to show graph of disk activity\n")
	b.WriteString("io_graph_combined = False\n\n")

	b.WriteString("#* Set the time in seconds for the minimum speed graph spans\n")
	b.WriteString("io_graph_speeds = \"\"\n\n")

	b.WriteString("#* Set fixed values for network graph max speed\n")
	b.WriteString("net_download = 100\n")
	b.WriteString("net_upload = 100\n\n")

	b.WriteString("#* Start network graph at zero instead of auto-scaling\n")
	b.WriteString("net_auto = True\n\n")

	b.WriteString("#* Network interface to show\n")
	b.WriteString("net_iface = \"\"\n\n")

	b.WriteString("#* Show network graph for upload/download combined or separate\n")
	b.WriteString("net_sync = True\n\n")

	b.WriteString("#* Show init screen at startup\n")
	b.WriteString("show_init = False\n\n")

	b.WriteString("#* Enable mouse support\n")
	b.WriteString("enable_mouse = True\n\n")

	b.WriteString("#* Set loglevel (ERROR, WARNING, INFO, DEBUG)\n")
	b.WriteString("log_level = \"WARNING\"\n")

	return []byte(b.String())
}

// SetTheme sets the btop theme.
func (h *Helper) SetTheme(themeName string) error {
	configFile := filepath.Join(h.configDir, "btop.conf")

	// Read current config
	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	// Replace theme line
	lines := strings.Split(string(data), "\n")
	var newLines []string
	themeSet := false

	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "color_theme") {
			newLines = append(newLines, fmt.Sprintf("color_theme = \"%s\"", themeName))
			themeSet = true
		} else {
			newLines = append(newLines, line)
		}
	}

	if !themeSet {
		// Add theme line if not found
		newLines = append([]string{fmt.Sprintf("color_theme = \"%s\"", themeName)}, newLines...)
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would set theme to: %s\n", themeName)
		return nil
	}

	if err := os.WriteFile(configFile, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// GetCurrentTheme returns the current theme name.
func (h *Helper) GetCurrentTheme() string {
	configFile := filepath.Join(h.configDir, "btop.conf")
	data, err := os.ReadFile(configFile)
	if err != nil {
		return "Default"
	}

	for line := range strings.SplitSeq(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "color_theme") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				theme := strings.TrimSpace(parts[1])
				theme = strings.Trim(theme, "\"'")
				return theme
			}
		}
	}

	return "Default"
}
