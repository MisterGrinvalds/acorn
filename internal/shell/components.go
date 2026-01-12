// Package shell provides shell integration generation and injection.
package shell

import (
	"github.com/mistergrinvalds/acorn/internal/componentconfig"
)

// componentOrder defines the order components should be loaded.
// This is critical - bootstrap and xdg must be first as other components depend on them.
var componentOrder = []string{
	// Bootstrap/core (MUST be first - sets up environment)
	"bootstrap", // Shell detection, platform detection, XDG vars, Homebrew
	"xdg",       // XDG helper functions, shell history setup
	"theme",     // Colors (depends on CURRENT_PLATFORM from bootstrap)
	"core",      // Shell options, prompt, keybindings (depends on theme)

	// Infrastructure
	"sync", // Dotfiles sync aliases and functions

	// Development tools
	"go",
	"python",
	"node",
	"vscode",
	"intellij",
	"neovim",

	// Terminal and shell
	"tmux",
	"fzf",
	"ghostty",
	"iterm2",

	// Version control
	"git",
	"github",

	// Cloud and infrastructure
	"cloudflare",
	"kubernetes",
	"database",

	// AI and ML
	"claude",
	"huggingface",
	"ollama",

	// Utilities
	"tools",
	"secrets",
}

// RegisterAllComponents registers all known components with the manager.
// Components are registered in a specific order to ensure dependencies are met.
func RegisterAllComponents(m *Manager) {
	for _, name := range componentOrder {
		registerComponentWithFiles(m, name)
	}
}

// GetComponentOrder returns the ordered list of component names.
// Used by shell.go to generate the entrypoint with correct sourcing order.
func GetComponentOrder() []string {
	return componentOrder
}

// registerComponentWithFiles loads and registers a component with its file specs.
func registerComponentWithFiles(m *Manager, name string) {
	result := loadComponentWithFiles(name)
	m.RegisterComponent(result.Component)
	if len(result.Files) > 0 {
		m.RegisterComponentFiles(name, result.Files)
	}
}

// GoComponent returns the Go shell integration component.
// Loads configuration from embedded YAML with optional user overrides.
func GoComponent() *Component {
	return loadComponentFromConfig("go")
}

// ComponentWithFiles holds both the shell component and its file specs.
type ComponentWithFiles struct {
	Component *Component
	Files     []FileSpec
}

// loadComponentFromConfig loads a component from YAML config.
// Falls back to a minimal component if config loading fails.
func loadComponentFromConfig(name string) *Component {
	result := loadComponentWithFiles(name)
	return result.Component
}

// loadComponentWithFiles loads a component with its file specs.
func loadComponentWithFiles(name string) *ComponentWithFiles {
	loader := componentconfig.NewLoader()
	cfg, err := loader.LoadBase(name)
	if err != nil {
		// Return minimal component on error
		return &ComponentWithFiles{
			Component: &Component{
				Name:        name,
				Description: "Component (config error: " + err.Error() + ")",
			},
		}
	}

	gen := NewGenerator()
	result := &ComponentWithFiles{
		Component: gen.GenerateComponent(cfg),
	}

	// Convert file configs to file specs
	if len(cfg.Files) > 0 {
		result.Files = make([]FileSpec, len(cfg.Files))
		for i, fc := range cfg.Files {
			result.Files[i] = FileSpec{
				Target:    fc.Target,
				Format:    fc.Format,
				Platforms: fc.Platforms,
				Values:    fc.Values,
			}
		}
	}

	return result
}

// VSCodeComponent returns the VS Code shell integration component.
func VSCodeComponent() *Component {
	return loadComponentFromConfig("vscode")
}

// ToolsComponent returns the tools shell integration component.
func ToolsComponent() *Component {
	return loadComponentFromConfig("tools")
}

// PythonComponent returns the Python shell integration component.
func PythonComponent() *Component {
	return loadComponentFromConfig("python")
}

// TmuxComponent returns the tmux shell integration component.
func TmuxComponent() *Component {
	return loadComponentFromConfig("tmux")
}

// ClaudeComponent returns the Claude Code shell integration component.
func ClaudeComponent() *Component {
	return loadComponentFromConfig("claude")
}

// CloudFlareComponent returns the CloudFlare shell integration component.
func CloudFlareComponent() *Component {
	return loadComponentFromConfig("cloudflare")
}

// SecretsComponent returns the secrets shell integration component.
func SecretsComponent() *Component {
	return loadComponentFromConfig("secrets")
}

// CoreComponent returns the core shell integration component.
// Includes shell options, prompt, keybindings, and utility functions.
func CoreComponent() *Component {
	return loadComponentFromConfig("core")
}

// DatabaseComponent returns the database shell integration component.
func DatabaseComponent() *Component {
	return loadComponentFromConfig("database")
}

// FzfComponent returns the FZF shell integration component.
func FzfComponent() *Component {
	return loadComponentFromConfig("fzf")
}

// GhosttyComponent returns the Ghostty shell integration component.
func GhosttyComponent() *Component {
	return loadComponentFromConfig("ghostty")
}

// GitComponent returns the Git shell integration component.
func GitComponent() *Component {
	return loadComponentFromConfig("git")
}

// GitHubComponent returns GitHub CLI integration shell functions.
func GitHubComponent() *Component {
	return loadComponentFromConfig("github")
}

// HuggingFaceComponent returns Hugging Face shell integration.
func HuggingFaceComponent() *Component {
	return loadComponentFromConfig("huggingface")
}

// KubernetesComponent returns Kubernetes shell integration.
func KubernetesComponent() *Component {
	return loadComponentFromConfig("kubernetes")
}

// NeovimComponent returns Neovim shell integration.
func NeovimComponent() *Component {
	return loadComponentFromConfig("neovim")
}

// NodeComponent returns Node.js shell integration.
func NodeComponent() *Component {
	return loadComponentFromConfig("node")
}

// OllamaComponent returns Ollama shell integration.
func OllamaComponent() *Component {
	return loadComponentFromConfig("ollama")
}
