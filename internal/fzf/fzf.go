// Package fzf provides FZF status and configuration helpers.
// Note: Most fzf functions remain in shell due to their interactive nature.
package fzf

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Status represents FZF installation status.
type Status struct {
	Installed   bool   `json:"installed" yaml:"installed"`
	Version     string `json:"version,omitempty" yaml:"version,omitempty"`
	Location    string `json:"location,omitempty" yaml:"location,omitempty"`
	FdInstalled bool   `json:"fd_installed" yaml:"fd_installed"`
	FdCommand   string `json:"fd_command,omitempty" yaml:"fd_command,omitempty"`
}

// Config represents FZF configuration.
type Config struct {
	DefaultCommand string `json:"default_command,omitempty" yaml:"default_command,omitempty"`
	AltCCommand    string `json:"alt_c_command,omitempty" yaml:"alt_c_command,omitempty"`
	CtrlTCommand   string `json:"ctrl_t_command,omitempty" yaml:"ctrl_t_command,omitempty"`
	DefaultOpts    string `json:"default_opts,omitempty" yaml:"default_opts,omitempty"`
}

// Helper provides FZF helper operations.
type Helper struct {
	verbose bool
}

// NewHelper creates a new FZF Helper.
func NewHelper(verbose bool) *Helper {
	return &Helper{
		verbose: verbose,
	}
}

// GetStatus returns FZF installation status.
func (h *Helper) GetStatus() *Status {
	status := &Status{}

	// Check if fzf is installed
	fzfPath, err := exec.LookPath("fzf")
	if err != nil {
		return status
	}

	status.Installed = true

	// Get version
	cmd := exec.Command("fzf", "--version")
	out, err := cmd.Output()
	if err == nil {
		parts := strings.Fields(string(out))
		if len(parts) > 0 {
			status.Version = parts[0]
		}
	}

	// Determine location
	status.Location = h.findFzfLocation(fzfPath)

	// Check for fd
	if runtime.GOOS == "linux" {
		if _, err := exec.LookPath("fdfind"); err == nil {
			status.FdInstalled = true
			status.FdCommand = "fdfind"
		} else if _, err := exec.LookPath("fd"); err == nil {
			status.FdInstalled = true
			status.FdCommand = "fd"
		}
	} else {
		if _, err := exec.LookPath("fd"); err == nil {
			status.FdInstalled = true
			status.FdCommand = "fd"
		}
	}

	return status
}

// findFzfLocation determines the FZF installation directory.
func (h *Helper) findFzfLocation(fzfPath string) string {
	// Check environment variable first
	if loc := os.Getenv("FZF_LOCATION"); loc != "" {
		return loc
	}

	// Check common locations
	locations := []string{
		"/opt/homebrew/opt/fzf",
		"/home/linuxbrew/.linuxbrew/opt/fzf",
		"/usr/share/fzf",
		filepath.Join(os.Getenv("HOME"), ".fzf"),
	}

	for _, loc := range locations {
		if info, err := os.Stat(loc); err == nil && info.IsDir() {
			return loc
		}
	}

	// Return the directory containing the fzf binary
	return filepath.Dir(fzfPath)
}

// GetConfig returns current FZF configuration from environment.
func (h *Helper) GetConfig() *Config {
	return &Config{
		DefaultCommand: os.Getenv("FZF_DEFAULT_COMMAND"),
		AltCCommand:    os.Getenv("FZF_ALT_C_COMMAND"),
		CtrlTCommand:   os.Getenv("FZF_CTRL_T_COMMAND"),
		DefaultOpts:    os.Getenv("FZF_DEFAULT_OPTS"),
	}
}

// GetThemeColors returns Catppuccin Mocha theme colors for FZF.
func (h *Helper) GetThemeColors() string {
	return `--color=bg+:#313244,bg:#1e1e2e,spinner:#f5e0dc,hl:#f38ba8
--color=fg:#cdd6f4,header:#f38ba8,info:#cba6f7,pointer:#f5e0dc
--color=marker:#f5e0dc,fg+:#cdd6f4,prompt:#cba6f7,hl+:#f38ba8`
}

// GetAvailableFunctions returns list of available fzf shell functions.
func (h *Helper) GetAvailableFunctions() []string {
	return []string{
		"fzf_files (ff)     - Interactive file finder with preview",
		"fe <query>         - Find and edit file",
		"fzf_cd (fcd)       - Interactive cd with preview",
		"fzf_git_branch     - Interactive git branch checkout",
		"fzf_git_log        - Interactive git log browser",
		"fzf_git_stash      - Interactive git stash browser",
		"fzf_git_add (fga)  - Interactive git add",
		"fzf_kill (fkill)   - Interactive process killer",
		"fzf_history (fh)   - Interactive history search",
		"fzf_env (fenv)     - Interactive env variable browser",
		"fzf_k8s_pod        - Interactive k8s pod selector",
		"fzf_k8s_logs       - Interactive k8s pod logs",
		"fzf_k8s_ns         - Interactive k8s namespace switcher",
		"fzf_docker_logs    - Interactive docker container logs",
		"fzf_docker_exec    - Interactive docker exec",
	}
}
