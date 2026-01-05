// Package tmux provides tmux session management and TPM helper functionality.
package tmux

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// SessionInfo contains tmux session information.
type SessionInfo struct {
	Name     string `json:"name" yaml:"name"`
	Windows  int    `json:"windows" yaml:"windows"`
	Created  string `json:"created,omitempty" yaml:"created,omitempty"`
	Attached bool   `json:"attached" yaml:"attached"`
}

// TmuxInfo contains tmux environment information.
type TmuxInfo struct {
	Version      string        `json:"version" yaml:"version"`
	ConfigFile   string        `json:"config_file" yaml:"config_file"`
	PluginDir    string        `json:"plugin_dir" yaml:"plugin_dir"`
	TPMInstalled bool          `json:"tpm_installed" yaml:"tpm_installed"`
	Sessions     []SessionInfo `json:"sessions,omitempty" yaml:"sessions,omitempty"`
}

// SmugSession contains smug session config information.
type SmugSession struct {
	Name        string `json:"name" yaml:"name"`
	Path        string `json:"path" yaml:"path"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

// SmugRepoStatus contains smug repo status information.
type SmugRepoStatus struct {
	Location    string        `json:"location" yaml:"location"`
	Initialized bool          `json:"initialized" yaml:"initialized"`
	Branch      string        `json:"branch,omitempty" yaml:"branch,omitempty"`
	Clean       bool          `json:"clean" yaml:"clean"`
	Sessions    []SmugSession `json:"sessions,omitempty" yaml:"sessions,omitempty"`
}

// Helper provides tmux helper operations.
type Helper struct {
	verbose bool
	dryRun  bool
}

// NewHelper creates a new Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	return &Helper{
		verbose: verbose,
		dryRun:  dryRun,
	}
}

// GetConfigDir returns the tmux config directory.
func GetConfigDir() string {
	if dir := os.Getenv("TMUX_CONFIG_DIR"); dir != "" {
		return dir
	}
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		configHome = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return filepath.Join(configHome, "tmux")
}

// GetPluginDir returns the tmux plugin directory.
func GetPluginDir() string {
	if dir := os.Getenv("TMUX_PLUGIN_DIR"); dir != "" {
		return dir
	}
	return filepath.Join(GetConfigDir(), "plugins")
}

// GetTPMDir returns the TPM directory.
func GetTPMDir() string {
	if dir := os.Getenv("TMUX_TPM_DIR"); dir != "" {
		return dir
	}
	return filepath.Join(GetPluginDir(), "tpm")
}

// GetConfigFile returns the tmux config file path.
func GetConfigFile() string {
	if file := os.Getenv("TMUX_CONF"); file != "" {
		return file
	}
	return filepath.Join(GetConfigDir(), "tmux.conf")
}

// GetSmugConfigDir returns the smug config directory.
func GetSmugConfigDir() string {
	if dir := os.Getenv("SMUG_CONFIG_DIR"); dir != "" {
		return dir
	}
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		configHome = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return filepath.Join(configHome, "smug")
}

// GetSmugRepoDir returns the smug git repo directory.
func GetSmugRepoDir() string {
	if dir := os.Getenv("SMUG_REPO_DIR"); dir != "" {
		return dir
	}
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		dataHome = filepath.Join(os.Getenv("HOME"), ".local", "share")
	}
	return filepath.Join(dataHome, "smug-sessions")
}

// GetSmugRepo returns the smug git repo URL.
func GetSmugRepo() string {
	if repo := os.Getenv("SMUG_REPO"); repo != "" {
		return repo
	}
	return "https://github.com/MisterGrinvalds/fmux.git"
}

// HasTmux checks if tmux is installed.
func (h *Helper) HasTmux() bool {
	_, err := exec.LookPath("tmux")
	return err == nil
}

// HasSmug checks if smug is installed.
func (h *Helper) HasSmug() bool {
	_, err := exec.LookPath("smug")
	return err == nil
}

// GetVersion returns the tmux version.
func (h *Helper) GetVersion() (string, error) {
	cmd := exec.Command("tmux", "-V")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// IsTPMInstalled checks if TPM is installed.
func (h *Helper) IsTPMInstalled() bool {
	_, err := os.Stat(GetTPMDir())
	return err == nil
}

// GetInfo returns tmux environment information.
func (h *Helper) GetInfo() (*TmuxInfo, error) {
	info := &TmuxInfo{
		ConfigFile:   GetConfigFile(),
		PluginDir:    GetPluginDir(),
		TPMInstalled: h.IsTPMInstalled(),
	}

	if version, err := h.GetVersion(); err == nil {
		info.Version = version
	} else {
		return nil, fmt.Errorf("tmux not installed")
	}

	// Get sessions
	sessions, _ := h.ListSessions()
	info.Sessions = sessions

	return info, nil
}

// ListSessions lists active tmux sessions.
func (h *Helper) ListSessions() ([]SessionInfo, error) {
	cmd := exec.Command("tmux", "list-sessions", "-F", "#{session_name}:#{session_windows}:#{session_created}:#{session_attached}")
	out, err := cmd.Output()
	if err != nil {
		return nil, nil // No sessions is not an error
	}

	var sessions []SessionInfo
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.Split(line, ":")
		if len(parts) >= 4 {
			windows := 0
			fmt.Sscanf(parts[1], "%d", &windows)
			sessions = append(sessions, SessionInfo{
				Name:     parts[0],
				Windows:  windows,
				Created:  parts[2],
				Attached: parts[3] == "1",
			})
		}
	}
	return sessions, nil
}

// InstallTPM installs Tmux Plugin Manager.
func (h *Helper) InstallTPM() error {
	tpmDir := GetTPMDir()

	if _, err := os.Stat(tpmDir); err == nil {
		return fmt.Errorf("TPM already installed at: %s", tpmDir)
	}

	// Create parent directory
	if err := os.MkdirAll(filepath.Dir(tpmDir), 0o755); err != nil {
		return fmt.Errorf("failed to create plugin directory: %w", err)
	}

	fmt.Println("Installing Tmux Plugin Manager...")
	if err := h.run("git", "clone", "https://github.com/tmux-plugins/tpm", tpmDir); err != nil {
		return fmt.Errorf("failed to clone TPM: %w", err)
	}

	return nil
}

// UpdateTPM updates Tmux Plugin Manager.
func (h *Helper) UpdateTPM() error {
	tpmDir := GetTPMDir()

	if _, err := os.Stat(tpmDir); os.IsNotExist(err) {
		return fmt.Errorf("TPM not installed. Run: acorn tmux tpm install")
	}

	fmt.Println("Updating TPM...")
	return h.runInDir(tpmDir, "git", "pull")
}

// InstallPlugins installs all tmux plugins via TPM.
func (h *Helper) InstallPlugins() error {
	tpmDir := GetTPMDir()

	if _, err := os.Stat(tpmDir); os.IsNotExist(err) {
		return fmt.Errorf("TPM not installed. Run: acorn tmux tpm install")
	}

	installScript := filepath.Join(tpmDir, "bin", "install_plugins")
	fmt.Println("Installing tmux plugins...")
	return h.run(installScript)
}

// UpdatePlugins updates all tmux plugins via TPM.
func (h *Helper) UpdatePlugins() error {
	tpmDir := GetTPMDir()

	if _, err := os.Stat(tpmDir); os.IsNotExist(err) {
		return fmt.Errorf("TPM not installed. Run: acorn tmux tpm install")
	}

	updateScript := filepath.Join(tpmDir, "bin", "update_plugins")
	fmt.Println("Updating tmux plugins...")
	return h.run(updateScript, "all")
}

// ReloadConfig reloads the tmux configuration.
func (h *Helper) ReloadConfig() error {
	configFile := GetConfigFile()

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return fmt.Errorf("tmux config not found: %s", configFile)
	}

	// Check if we're in tmux
	if os.Getenv("TMUX") == "" {
		return fmt.Errorf("not in a tmux session; config will load on next tmux start")
	}

	return h.run("tmux", "source-file", configFile)
}

// ListSmugSessions lists available smug session configs.
func (h *Helper) ListSmugSessions() ([]SmugSession, error) {
	smugDir := GetSmugConfigDir()

	var sessions []SmugSession

	if _, err := os.Stat(smugDir); os.IsNotExist(err) {
		return sessions, nil
	}

	entries, err := os.ReadDir(smugDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yml") {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".yml")
		path := filepath.Join(smugDir, entry.Name())

		// Try to read description from first comment line
		desc := ""
		if content, err := os.ReadFile(path); err == nil {
			for _, line := range strings.Split(string(content), "\n") {
				if strings.HasPrefix(line, "# smug session:") {
					desc = strings.TrimPrefix(line, "# smug session:")
					desc = strings.TrimSpace(desc)
					break
				}
			}
		}

		sessions = append(sessions, SmugSession{
			Name:        name,
			Path:        path,
			Description: desc,
		})
	}

	return sessions, nil
}

// CreateSmugSession creates a new smug session config.
func (h *Helper) CreateSmugSession(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("session name required")
	}

	smugDir := GetSmugConfigDir()
	if err := os.MkdirAll(smugDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create smug config directory: %w", err)
	}

	configFile := filepath.Join(smugDir, name+".yml")
	if _, err := os.Stat(configFile); err == nil {
		return configFile, fmt.Errorf("config already exists: %s", configFile)
	}

	template := fmt.Sprintf(`# smug session: %s
# Created: %s
# Usage: smug start %s

session: %s
root: ~/
attach: true

windows:
  - name: main
    commands:
      - echo "Welcome to %s session"

  - name: editor
    commands:
      - nvim

  - name: terminal
    panes:
      - type: horizontal
        commands:
          - echo "Ready"
`, name, getCurrentDate(), name, name, name)

	if h.dryRun {
		fmt.Printf("[dry-run] would create: %s\n", configFile)
		return configFile, nil
	}

	if err := os.WriteFile(configFile, []byte(template), 0o644); err != nil {
		return "", fmt.Errorf("failed to create config: %w", err)
	}

	return configFile, nil
}

// InstallSmug installs smug using brew or go.
func (h *Helper) InstallSmug() error {
	if h.HasSmug() {
		version, _ := exec.Command("smug", "--version").Output()
		return fmt.Errorf("smug already installed: %s", strings.TrimSpace(string(version)))
	}

	fmt.Println("Installing smug...")

	// Try brew first
	if _, err := exec.LookPath("brew"); err == nil {
		return h.run("brew", "install", "smug")
	}

	// Try go install
	if _, err := exec.LookPath("go"); err == nil {
		return h.run("go", "install", "github.com/ivaaaan/smug@latest")
	}

	return fmt.Errorf("no package manager found; install manually with: brew install smug")
}

// SmugRepoInit initializes the smug sessions git repo.
func (h *Helper) SmugRepoInit() error {
	repoDir := GetSmugRepoDir()
	repoURL := GetSmugRepo()

	if _, err := os.Stat(filepath.Join(repoDir, ".git")); err == nil {
		fmt.Printf("Smug repo already initialized at: %s\n", repoDir)
		fmt.Println("Pulling latest...")
		return h.runInDir(repoDir, "git", "pull", "--rebase")
	}

	fmt.Println("Cloning smug sessions repo...")
	if err := os.MkdirAll(filepath.Dir(repoDir), 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return h.run("git", "clone", repoURL, repoDir)
}

// SmugRepoStatus returns the smug repo status.
func (h *Helper) SmugRepoStatus() (*SmugRepoStatus, error) {
	repoDir := GetSmugRepoDir()

	status := &SmugRepoStatus{
		Location:    repoDir,
		Initialized: false,
		Clean:       true,
	}

	if _, err := os.Stat(filepath.Join(repoDir, ".git")); os.IsNotExist(err) {
		return status, nil
	}

	status.Initialized = true

	// Get branch
	if out, err := exec.Command("git", "-C", repoDir, "branch", "--show-current").Output(); err == nil {
		status.Branch = strings.TrimSpace(string(out))
	}

	// Check if clean
	if out, err := exec.Command("git", "-C", repoDir, "status", "--porcelain").Output(); err == nil {
		status.Clean = len(strings.TrimSpace(string(out))) == 0
	}

	// List sessions
	entries, err := os.ReadDir(repoDir)
	if err == nil {
		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".yml") {
				name := strings.TrimSuffix(entry.Name(), ".yml")
				status.Sessions = append(status.Sessions, SmugSession{
					Name: name,
					Path: filepath.Join(repoDir, entry.Name()),
				})
			}
		}
	}

	return status, nil
}

// SmugRepoPull pulls the latest from remote.
func (h *Helper) SmugRepoPull() error {
	repoDir := GetSmugRepoDir()

	if _, err := os.Stat(filepath.Join(repoDir, ".git")); os.IsNotExist(err) {
		return fmt.Errorf("smug repo not initialized. Run: acorn tmux smug repo-init")
	}

	fmt.Println("Pulling latest sessions...")
	return h.runInDir(repoDir, "git", "pull", "--rebase")
}

// SmugRepoPush commits and pushes changes.
func (h *Helper) SmugRepoPush(message string) error {
	repoDir := GetSmugRepoDir()

	if _, err := os.Stat(filepath.Join(repoDir, ".git")); os.IsNotExist(err) {
		return fmt.Errorf("smug repo not initialized. Run: acorn tmux smug repo-init")
	}

	// Check for changes
	out, err := exec.Command("git", "-C", repoDir, "status", "--porcelain").Output()
	if err != nil {
		return err
	}
	if len(strings.TrimSpace(string(out))) == 0 {
		return fmt.Errorf("no changes to push")
	}

	if message == "" {
		message = "Update smug sessions"
	}

	// Add yml files
	if err := h.runInDir(repoDir, "git", "add", "*.yml", "*.yaml"); err != nil {
		// Ignore if no yml files
	}

	// Commit
	if err := h.runInDir(repoDir, "git", "commit", "-m", message); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	// Push
	fmt.Println("Pushing to remote...")
	return h.runInDir(repoDir, "git", "push")
}

// SmugRepoSync does a full sync (pull + push).
func (h *Helper) SmugRepoSync() error {
	repoDir := GetSmugRepoDir()

	if _, err := os.Stat(filepath.Join(repoDir, ".git")); os.IsNotExist(err) {
		return fmt.Errorf("smug repo not initialized. Run: acorn tmux smug repo-init")
	}

	fmt.Println("Syncing smug sessions...")

	// Check for local changes
	out, err := exec.Command("git", "-C", repoDir, "status", "--porcelain").Output()
	if err != nil {
		return err
	}
	hasChanges := len(strings.TrimSpace(string(out))) > 0

	// Stash if needed
	if hasChanges {
		if err := h.runInDir(repoDir, "git", "stash"); err != nil {
			return fmt.Errorf("stash failed: %w", err)
		}
	}

	// Pull
	if err := h.SmugRepoPull(); err != nil {
		return err
	}

	// Pop stash and push if we had changes
	if hasChanges {
		if err := h.runInDir(repoDir, "git", "stash", "pop"); err != nil {
			return fmt.Errorf("stash pop failed: %w", err)
		}
		if err := h.SmugRepoPush("Sync local session changes"); err != nil {
			return err
		}
	}

	return nil
}

// SmugLinkConfigs links smug configs from git repo or dotfiles.
func (h *Helper) SmugLinkConfigs(dotfilesRoot string) error {
	repoDir := GetSmugRepoDir()
	targetDir := GetSmugConfigDir()

	// Prefer git repo if initialized
	if _, err := os.Stat(filepath.Join(repoDir, ".git")); err == nil {
		fmt.Println("Linking smug configs from git repo...")
		if h.dryRun {
			fmt.Printf("[dry-run] would link: %s -> %s\n", targetDir, repoDir)
			return nil
		}
		os.RemoveAll(targetDir)
		if err := os.Symlink(repoDir, targetDir); err != nil {
			return fmt.Errorf("failed to create symlink: %w", err)
		}
		fmt.Printf("Linked: %s -> %s\n", targetDir, repoDir)
		return nil
	}

	// Fallback to dotfiles
	sourceDir := filepath.Join(dotfilesRoot, "components", "tmux", "config", "smug")
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		return fmt.Errorf("no smug configs found. Run: acorn tmux smug repo-init")
	}

	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	fmt.Println("Linking smug configs from dotfiles...")
	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".yml") {
			continue
		}
		source := filepath.Join(sourceDir, entry.Name())
		target := filepath.Join(targetDir, entry.Name())

		if h.dryRun {
			fmt.Printf("[dry-run] would link: %s\n", entry.Name())
			continue
		}

		os.Remove(target)
		if err := os.Symlink(source, target); err != nil {
			return fmt.Errorf("failed to link %s: %w", entry.Name(), err)
		}
		fmt.Printf("  Linked: %s\n", entry.Name())
	}

	return nil
}

// run executes a command.
func (h *Helper) run(name string, args ...string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would run: %s %s\n", name, strings.Join(args, " "))
		return nil
	}

	if h.verbose {
		fmt.Printf("Running: %s %s\n", name, strings.Join(args, " "))
	}

	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// runInDir executes a command in a specific directory.
func (h *Helper) runInDir(dir, name string, args ...string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would run in %s: %s %s\n", dir, name, strings.Join(args, " "))
		return nil
	}

	if h.verbose {
		fmt.Printf("Running in %s: %s %s\n", dir, name, strings.Join(args, " "))
	}

	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func getCurrentDate() string {
	out, _ := exec.Command("date", "+%Y-%m-%d").Output()
	return strings.TrimSpace(string(out))
}
