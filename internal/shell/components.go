// Package shell provides shell integration generation and injection.
package shell

import (
	"github.com/mistergrinvalds/acorn/internal/componentconfig"
)

// RegisterAllComponents registers all known components with the manager.
func RegisterAllComponents(m *Manager) {
	m.RegisterComponent(GoComponent())
	m.RegisterComponent(VSCodeComponent())
	m.RegisterComponent(ToolsComponent())
	m.RegisterComponent(PythonComponent())
	m.RegisterComponent(TmuxComponent())
	m.RegisterComponent(ClaudeComponent())
	m.RegisterComponent(CloudFlareComponent())
	m.RegisterComponent(SecretsComponent())
	m.RegisterComponent(DatabaseComponent())
	m.RegisterComponent(FzfComponent())
	m.RegisterComponent(GhosttyComponent())
	m.RegisterComponent(GitComponent())
	m.RegisterComponent(GitHubComponent())
	m.RegisterComponent(HuggingFaceComponent())
	m.RegisterComponent(KubernetesComponent())
	m.RegisterComponent(NeovimComponent())
	m.RegisterComponent(NodeComponent())
	m.RegisterComponent(OllamaComponent())
}

// GoComponent returns the Go shell integration component.
// Loads configuration from embedded YAML with optional user overrides.
func GoComponent() *Component {
	return loadComponentFromConfig("go")
}

// loadComponentFromConfig loads a component from YAML config.
// Falls back to a minimal component if config loading fails.
func loadComponentFromConfig(name string) *Component {
	loader := componentconfig.NewLoader()
	cfg, err := loader.LoadBase(name)
	if err != nil {
		// Return minimal component on error
		return &Component{
			Name:        name,
			Description: "Component (config error: " + err.Error() + ")",
		}
	}

	gen := NewGenerator()
	return gen.GenerateComponent(cfg)
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
