// Package ghostty provides Ghostty terminal emulator configuration helpers.
package ghostty

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Info represents Ghostty information.
type Info struct {
	Installed bool   `json:"installed" yaml:"installed"`
	Version   string `json:"version,omitempty" yaml:"version,omitempty"`
	Config    string `json:"config" yaml:"config"`
	Theme     string `json:"theme,omitempty" yaml:"theme,omitempty"`
	Font      string `json:"font,omitempty" yaml:"font,omitempty"`
	FontSize  string `json:"font_size,omitempty" yaml:"font_size,omitempty"`
}

// Backup represents a config backup.
type Backup struct {
	Name      string `json:"name" yaml:"name"`
	Path      string `json:"path" yaml:"path"`
	Timestamp string `json:"timestamp" yaml:"timestamp"`
}

// Helper provides Ghostty configuration operations.
type Helper struct {
	configPath string
	backupDir  string
	verbose    bool
}

// NewHelper creates a new Ghostty Helper.
func NewHelper(verbose bool) *Helper {
	configPath := os.Getenv("GHOSTTY_CONFIG")
	if configPath == "" {
		xdgConfig := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfig == "" {
			home, _ := os.UserHomeDir()
			xdgConfig = filepath.Join(home, ".config")
		}
		configPath = filepath.Join(xdgConfig, "ghostty", "config")
	}

	xdgData := os.Getenv("XDG_DATA_HOME")
	if xdgData == "" {
		home, _ := os.UserHomeDir()
		xdgData = filepath.Join(home, ".local", "share")
	}
	backupDir := filepath.Join(xdgData, "ghostty", "backups")

	return &Helper{
		configPath: configPath,
		backupDir:  backupDir,
		verbose:    verbose,
	}
}

// GetConfigPath returns the config file path.
func (h *Helper) GetConfigPath() string {
	return h.configPath
}

// GetInfo returns Ghostty installation and config info.
func (h *Helper) GetInfo() *Info {
	info := &Info{
		Config: h.configPath,
	}

	// Check if installed
	if path, err := exec.LookPath("ghostty"); err == nil {
		info.Installed = true
		// Try to get version
		cmd := exec.Command(path, "--version")
		if out, err := cmd.Output(); err == nil {
			info.Version = strings.TrimSpace(string(out))
		} else {
			info.Version = "installed"
		}
	}

	// Parse config for theme and font
	if file, err := os.Open(h.configPath); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "theme = ") {
				info.Theme = strings.TrimPrefix(line, "theme = ")
			} else if strings.HasPrefix(line, "font-family = ") {
				info.Font = strings.Trim(strings.TrimPrefix(line, "font-family = "), "\"")
			} else if strings.HasPrefix(line, "font-size = ") {
				info.FontSize = strings.TrimPrefix(line, "font-size = ")
			}
		}
	}

	return info
}

// SetTheme sets the Ghostty theme.
func (h *Helper) SetTheme(theme string) error {
	if theme == "" {
		return fmt.Errorf("theme name is required")
	}

	if _, err := os.Stat(h.configPath); os.IsNotExist(err) {
		return fmt.Errorf("ghostty config not found: %s", h.configPath)
	}

	content, err := os.ReadFile(h.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	found := false
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "theme = ") {
			lines[i] = "theme = " + theme
			found = true
			break
		}
	}

	if !found {
		// Add theme after first non-comment line or at beginning
		for i, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
				lines = append(lines[:i+1], append([]string{"theme = " + theme}, lines[i+1:]...)...)
				break
			}
		}
	}

	if err := os.WriteFile(h.configPath, []byte(strings.Join(lines, "\n")), 0o644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// SetFont sets the Ghostty font.
func (h *Helper) SetFont(font string, size string) error {
	if font == "" {
		return fmt.Errorf("font name is required")
	}

	if _, err := os.Stat(h.configPath); os.IsNotExist(err) {
		return fmt.Errorf("ghostty config not found: %s", h.configPath)
	}

	content, err := os.ReadFile(h.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	fontPattern := regexp.MustCompile(`^font-family\s*=`)
	sizePattern := regexp.MustCompile(`^font-size\s*=`)

	for i, line := range lines {
		if fontPattern.MatchString(strings.TrimSpace(line)) {
			lines[i] = fmt.Sprintf("font-family = \"%s\"", font)
		}
		if size != "" && sizePattern.MatchString(strings.TrimSpace(line)) {
			lines[i] = fmt.Sprintf("font-size = %s", size)
		}
	}

	if err := os.WriteFile(h.configPath, []byte(strings.Join(lines, "\n")), 0o644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// CreateBackup creates a backup of the current config.
func (h *Helper) CreateBackup() (*Backup, error) {
	if _, err := os.Stat(h.configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("ghostty config not found: %s", h.configPath)
	}

	if err := os.MkdirAll(h.backupDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("config.%s", timestamp)
	backupPath := filepath.Join(h.backupDir, backupName)

	content, err := os.ReadFile(h.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	if err := os.WriteFile(backupPath, content, 0o644); err != nil {
		return nil, fmt.Errorf("failed to write backup: %w", err)
	}

	return &Backup{
		Name:      backupName,
		Path:      backupPath,
		Timestamp: timestamp,
	}, nil
}

// ListBackups lists all config backups.
func (h *Helper) ListBackups() ([]Backup, error) {
	if _, err := os.Stat(h.backupDir); os.IsNotExist(err) {
		return nil, nil
	}

	entries, err := os.ReadDir(h.backupDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	var backups []Backup
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, "config.") {
			timestamp := strings.TrimPrefix(name, "config.")
			backups = append(backups, Backup{
				Name:      name,
				Path:      filepath.Join(h.backupDir, name),
				Timestamp: timestamp,
			})
		}
	}

	return backups, nil
}

// RestoreBackup restores a config from backup.
func (h *Helper) RestoreBackup(backupName string) error {
	if backupName == "" {
		return fmt.Errorf("backup name is required")
	}

	backupPath := filepath.Join(h.backupDir, backupName)
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		// Try as absolute path
		backupPath = backupName
		if _, err := os.Stat(backupPath); os.IsNotExist(err) {
			return fmt.Errorf("backup not found: %s", backupName)
		}
	}

	// Create backup of current config first
	if _, err := h.CreateBackup(); err != nil {
		if h.verbose {
			fmt.Printf("Warning: could not backup current config: %v\n", err)
		}
	}

	content, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup: %w", err)
	}

	// Ensure config directory exists
	if err := os.MkdirAll(filepath.Dir(h.configPath), 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(h.configPath, content, 0o644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// GetAvailableThemes returns common Ghostty themes.
func (h *Helper) GetAvailableThemes() []string {
	return []string{
		"Catppuccin Mocha",
		"Catppuccin Macchiato",
		"Catppuccin Frappe",
		"Catppuccin Latte",
		"Dracula",
		"Gruvbox Dark",
		"Nord",
		"One Dark",
		"Solarized Dark",
		"Solarized Light",
		"Tokyo Night",
	}
}
