// Package shell provides shell integration generation and injection.
package shell

import (
	"github.com/mistergrinvalds/acorn/internal/utils/config"
)

// defaultBootstrap is the fallback bootstrap order when no scaffold is available.
// These components MUST load first — they set up the shell environment.
var defaultBootstrap = []string{
	"bootstrap", "xdg", "theme", "core",
}

// defaultOptional is the fallback optional component order when no scaffold is available.
var defaultOptional = []string{
	"sync",
	"go", "python", "r", "node", "vscode", "intellij", "neovim",
	"tmux", "fzf", "ghostty", "iterm2", "karabiner",
	"git", "github",
	"cloudflare", "kubernetes", "database", "postgres", "tailscale", "pulumi", "terraform", "vault",
	"docker", "docker-compose", "lazydocker",
	"claude", "huggingface", "ollama",
	"posting",
	"tools", "secrets", "wget",
}

// RegisterAllComponents registers all known components with the manager.
// Uses the scaffold's resolved shell order (bootstrap first, then optional).
func RegisterAllComponents(m *Manager) {
	order := GetComponentOrder()
	for _, name := range order {
		registerComponentWithFiles(m, name)
	}
}

// GetComponentOrder returns the ordered list of component names for shell generation.
//
// Resolution:
//  1. Bootstrap components always come first (from scaffold or default).
//  2. Optional components follow — either explicitly listed in the scaffold,
//     or auto-derived from all scaffold groups in declaration order.
//  3. If no scaffold is found, falls back to hardcoded defaults.
func GetComponentOrder() []string {
	s, err := config.LoadScaffold()
	if err == nil && s != nil {
		return s.ResolveShellOrder()
	}
	// No scaffold — use hardcoded defaults
	result := make([]string, 0, len(defaultBootstrap)+len(defaultOptional))
	result = append(result, defaultBootstrap...)
	result = append(result, defaultOptional...)
	return result
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
	loader := config.NewComponentLoader()
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
