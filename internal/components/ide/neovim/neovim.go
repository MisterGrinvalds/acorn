// Package neovim provides Neovim configuration management functionality.
package neovim

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// HealthStatus represents Neovim health check status.
type HealthStatus struct {
	Installed     bool   `json:"installed" yaml:"installed"`
	Version       string `json:"version,omitempty" yaml:"version,omitempty"`
	ConfigDir     string `json:"config_dir" yaml:"config_dir"`
	ConfigType    string `json:"config_type" yaml:"config_type"` // symlink, directory, not_found
	ConfigTarget  string `json:"config_target,omitempty" yaml:"config_target,omitempty"`
	ConfigStatus  string `json:"config_status" yaml:"config_status"` // ok, broken, not_found
	InitFile      string `json:"init_file,omitempty" yaml:"init_file,omitempty"`
	PluginManager string `json:"plugin_manager,omitempty" yaml:"plugin_manager,omitempty"`
}

// Helper provides Neovim helper operations.
type Helper struct {
	verbose bool
}

// NewHelper creates a new Neovim Helper.
func NewHelper(verbose bool) *Helper {
	return &Helper{
		verbose: verbose,
	}
}

// GetConfigDir returns the Neovim config directory path.
func (h *Helper) GetConfigDir() string {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, _ := os.UserHomeDir()
		configHome = filepath.Join(home, ".config")
	}
	return filepath.Join(configHome, "nvim")
}

// GetDataDir returns the Neovim data directory path.
func (h *Helper) GetDataDir() string {
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		home, _ := os.UserHomeDir()
		dataHome = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(dataHome, "nvim")
}

// GetCacheDir returns the Neovim cache directory path.
func (h *Helper) GetCacheDir() string {
	cacheHome := os.Getenv("XDG_CACHE_HOME")
	if cacheHome == "" {
		home, _ := os.UserHomeDir()
		cacheHome = filepath.Join(home, ".cache")
	}
	return filepath.Join(cacheHome, "nvim")
}

// GetStateDir returns the Neovim state directory path.
func (h *Helper) GetStateDir() string {
	stateHome := os.Getenv("XDG_STATE_HOME")
	if stateHome == "" {
		home, _ := os.UserHomeDir()
		stateHome = filepath.Join(home, ".local", "state")
	}
	return filepath.Join(stateHome, "nvim")
}

// GetHealth returns Neovim health status.
func (h *Helper) GetHealth() *HealthStatus {
	status := &HealthStatus{
		ConfigDir: h.GetConfigDir(),
	}

	// Check if nvim is installed
	if _, err := exec.LookPath("nvim"); err == nil {
		status.Installed = true
		// Get version
		cmd := exec.Command("nvim", "--version")
		if out, err := cmd.Output(); err == nil {
			lines := strings.Split(string(out), "\n")
			if len(lines) > 0 {
				status.Version = strings.TrimSpace(lines[0])
			}
		}
	}

	// Check config directory
	configDir := status.ConfigDir
	info, err := os.Lstat(configDir)
	if os.IsNotExist(err) {
		status.ConfigType = "not_found"
		status.ConfigStatus = "not_found"
	} else if info.Mode()&os.ModeSymlink != 0 {
		status.ConfigType = "symlink"
		if target, err := os.Readlink(configDir); err == nil {
			status.ConfigTarget = target
			if _, err := os.Stat(configDir); err == nil {
				status.ConfigStatus = "ok"
			} else {
				status.ConfigStatus = "broken"
			}
		}
	} else if info.IsDir() {
		status.ConfigType = "directory"
		status.ConfigStatus = "ok"
	}

	// Check init file
	if status.ConfigStatus == "ok" {
		if _, err := os.Stat(filepath.Join(configDir, "init.lua")); err == nil {
			status.InitFile = "init.lua"
		} else if _, err := os.Stat(filepath.Join(configDir, "init.vim")); err == nil {
			status.InitFile = "init.vim"
		}
	}

	// Try to detect plugin manager
	if status.ConfigStatus == "ok" {
		if _, err := os.Stat(filepath.Join(configDir, "lua", "lazy")); err == nil {
			status.PluginManager = "lazy.nvim"
		} else {
			dataDir := h.GetDataDir()
			if _, err := os.Stat(filepath.Join(dataDir, "site", "pack", "packer")); err == nil {
				status.PluginManager = "packer.nvim"
			} else if _, err := os.Stat(filepath.Join(dataDir, "lazy")); err == nil {
				status.PluginManager = "lazy.nvim"
			}
		}
	}

	return status
}

// Update updates the Neovim config repo.
func (h *Helper) Update() error {
	configDir := h.GetConfigDir()

	info, err := os.Lstat(configDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("no Neovim config found. Run nvim_setup first")
	}

	var repoPath string
	if info.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(configDir)
		if err != nil {
			return fmt.Errorf("failed to read symlink: %w", err)
		}
		repoPath = target
	} else if info.IsDir() {
		repoPath = configDir
	}

	// Check if it's a git repo
	gitDir := filepath.Join(repoPath, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return fmt.Errorf("Neovim config is not a git repository")
	}

	fmt.Printf("Updating Neovim config at %s...\n", repoPath)
	cmd := exec.Command("git", "-C", repoPath, "pull")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Clean removes Neovim data, cache, and state directories.
func (h *Helper) Clean(force bool) error {
	dataDir := h.GetDataDir()
	cacheDir := h.GetCacheDir()
	stateDir := h.GetStateDir()

	if !force {
		return fmt.Errorf("use --force to actually clean the directories")
	}

	var errs []string

	if _, err := os.Stat(dataDir); err == nil {
		if err := os.RemoveAll(dataDir); err != nil {
			errs = append(errs, fmt.Sprintf("failed to remove data dir: %v", err))
		}
	}

	if _, err := os.Stat(cacheDir); err == nil {
		if err := os.RemoveAll(cacheDir); err != nil {
			errs = append(errs, fmt.Sprintf("failed to remove cache dir: %v", err))
		}
	}

	if _, err := os.Stat(stateDir); err == nil {
		if err := os.RemoveAll(stateDir); err != nil {
			errs = append(errs, fmt.Sprintf("failed to remove state dir: %v", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, "; "))
	}

	return nil
}

// GetPluginInfo returns dotfiles.nvim plugin setup info.
func (h *Helper) GetPluginInfo() string {
	dotfilesRoot := os.Getenv("DOTFILES_ROOT")
	pluginPath := filepath.Join(dotfilesRoot, "components", "neovim", "plugin")

	return fmt.Sprintf(`=== dotfiles.nvim Plugin ===

The dotfiles.nvim plugin lets you run the installer from within Neovim.

Plugin location: %s

=== Setup Instructions ===

For lazy.nvim, add to your plugins:

{
  dir = vim.env.DOTFILES_ROOT and (vim.env.DOTFILES_ROOT .. "/components/neovim/plugin") or nil,
  name = "dotfiles",
  config = function()
    require("dotfiles").setup({
      auto_setup = { enabled = true },  -- Prompt if not installed
    })
  end,
  cond = vim.env.DOTFILES_ROOT ~= nil,
}

=== Quick Setup (no config needed) ===

Just run this command in Neovim:
  :DotfilesSetup

=== Available Commands ===

  :Dotfiles           - Interactive installer
  :DotfilesMinimal    - Quick install (dotfiles + configs)
  :DotfilesComponents - Component-based installer
  :DotfilesComponent <name> - Install specific component
  :DotfilesUpdate     - Git pull dotfiles repo
  :DotfilesStatus     - Check installation status
`, pluginPath)
}
