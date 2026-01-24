// Package opencode provides OpenCode AI coding assistant helper functionality.
package opencode

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Status represents OpenCode installation status.
type Status struct {
	Installed   bool   `json:"installed" yaml:"installed"`
	Version     string `json:"version,omitempty" yaml:"version,omitempty"`
	ConfigDir   string `json:"config_dir,omitempty" yaml:"config_dir,omitempty"`
	HasConfig   bool   `json:"has_config" yaml:"has_config"`
	Provider    string `json:"provider,omitempty" yaml:"provider,omitempty"`
}

// Provider represents an AI provider configuration.
type Provider struct {
	Name    string `json:"name" yaml:"name"`
	Model   string `json:"model,omitempty" yaml:"model,omitempty"`
	APIKey  bool   `json:"api_key" yaml:"api_key"`
}

// Helper provides OpenCode helper operations.
type Helper struct {
	verbose bool
	dryRun  bool
}

// NewHelper creates a new OpenCode Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	return &Helper{
		verbose: verbose,
		dryRun:  dryRun,
	}
}

// GetStatus returns OpenCode status information.
func (h *Helper) GetStatus() *Status {
	status := &Status{}

	// Check if opencode is installed
	out, err := exec.Command("opencode", "--version").Output()
	if err != nil {
		status.Installed = false
		return status
	}

	status.Installed = true
	status.Version = strings.TrimSpace(string(out))

	// Get config directory
	status.ConfigDir = h.getConfigDir()

	// Check if config exists
	configFile := filepath.Join(status.ConfigDir, "config.json")
	if _, err := os.Stat(configFile); err == nil {
		status.HasConfig = true
	}

	return status
}

// getConfigDir returns the OpenCode config directory.
func (h *Helper) getConfigDir() string {
	// OpenCode uses XDG config
	xdgConfig := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfig == "" {
		home, _ := os.UserHomeDir()
		xdgConfig = filepath.Join(home, ".config")
	}
	return filepath.Join(xdgConfig, "opencode")
}

// Launch starts OpenCode in the current directory.
func (h *Helper) Launch(args ...string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would run: opencode %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("opencode", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Init initializes OpenCode in the current directory.
func (h *Helper) Init() error {
	if h.dryRun {
		fmt.Println("[dry-run] would initialize OpenCode in current directory")
		return nil
	}

	// OpenCode init is done interactively via /init command
	// We'll just launch it
	return h.Launch()
}

// RunCommand runs an OpenCode command non-interactively.
func (h *Helper) RunCommand(prompt string) error {
	if prompt == "" {
		return fmt.Errorf("prompt is required")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run opencode with prompt: %s\n", prompt)
		return nil
	}

	cmd := exec.Command("opencode", "-p", prompt)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// SetProvider sets the AI provider.
func (h *Helper) SetProvider(provider string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would set provider to: %s\n", provider)
		return nil
	}

	// OpenCode provider is set via environment or config
	fmt.Printf("Set provider via OPENCODE_PROVIDER=%s or in config\n", provider)
	return nil
}

// ListProviders returns supported providers.
func (h *Helper) ListProviders() []Provider {
	return []Provider{
		{Name: "anthropic", Model: "claude-sonnet-4-20250514"},
		{Name: "openai", Model: "gpt-4o"},
		{Name: "google", Model: "gemini-2.0-flash"},
		{Name: "groq", Model: "llama-3.3-70b-versatile"},
		{Name: "aws-bedrock", Model: "anthropic.claude-sonnet-4-20250514-v1:0"},
		{Name: "azure-openai", Model: "gpt-4o"},
		{Name: "openrouter", Model: "anthropic/claude-sonnet-4"},
	}
}

// Install installs OpenCode.
func (h *Helper) Install() error {
	if _, err := exec.LookPath("opencode"); err == nil {
		return fmt.Errorf("OpenCode is already installed")
	}

	if h.dryRun {
		fmt.Println("[dry-run] would install OpenCode via Homebrew")
		return nil
	}

	fmt.Println("Installing OpenCode...")
	cmd := exec.Command("brew", "install", "opencode")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Upgrade upgrades OpenCode to the latest version.
func (h *Helper) Upgrade() error {
	if h.dryRun {
		fmt.Println("[dry-run] would upgrade OpenCode")
		return nil
	}

	fmt.Println("Upgrading OpenCode...")
	cmd := exec.Command("brew", "upgrade", "opencode")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
