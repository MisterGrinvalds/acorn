package tools

// ToolDefinition defines a tool to track.
type ToolDefinition struct {
	Name        string   // Tool binary name
	Category    string   // Category (system, languages, cloud, development)
	VersionArgs []string // Command args to get version (e.g., "--version")
	Description string   // Short description
	InstallHint string   // How to install (e.g., "brew install git")
}

// Category constants for grouping tools.
const (
	CategorySystem      = "System"
	CategoryLanguages   = "Languages"
	CategoryCloud       = "Cloud"
	CategoryDevelopment = "Development"
)

// DefaultRegistry returns the built-in tool registry.
func DefaultRegistry() []ToolDefinition {
	return []ToolDefinition{
		// System tools
		{Name: "git", Category: CategorySystem, VersionArgs: []string{"--version"}, Description: "Version control system", InstallHint: "brew install git"},
		{Name: "curl", Category: CategorySystem, VersionArgs: []string{"--version"}, Description: "HTTP client", InstallHint: "brew install curl"},
		{Name: "jq", Category: CategorySystem, VersionArgs: []string{"--version"}, Description: "JSON processor", InstallHint: "brew install jq"},
		{Name: "wget", Category: CategorySystem, VersionArgs: []string{"--version"}, Description: "File downloader", InstallHint: "brew install wget"},
		{Name: "make", Category: CategorySystem, VersionArgs: []string{"--version"}, Description: "Build automation", InstallHint: "xcode-select --install"},
		{Name: "yq", Category: CategorySystem, VersionArgs: []string{"--version"}, Description: "YAML processor", InstallHint: "brew install yq"},

		// Languages
		{Name: "go", Category: CategoryLanguages, VersionArgs: []string{"version"}, Description: "Go programming language", InstallHint: "brew install go"},
		{Name: "node", Category: CategoryLanguages, VersionArgs: []string{"--version"}, Description: "Node.js runtime", InstallHint: "nvm install --lts"},
		{Name: "python3", Category: CategoryLanguages, VersionArgs: []string{"--version"}, Description: "Python 3", InstallHint: "brew install python3"},
		{Name: "cargo", Category: CategoryLanguages, VersionArgs: []string{"--version"}, Description: "Rust package manager", InstallHint: "curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh"},
		{Name: "uv", Category: CategoryLanguages, VersionArgs: []string{"--version"}, Description: "Fast Python package manager", InstallHint: "brew install uv"},

		// Cloud tools
		{Name: "aws", Category: CategoryCloud, VersionArgs: []string{"--version"}, Description: "AWS CLI", InstallHint: "brew install awscli"},
		{Name: "kubectl", Category: CategoryCloud, VersionArgs: []string{"version", "--client"}, Description: "Kubernetes CLI", InstallHint: "brew install kubectl"},
		{Name: "docker", Category: CategoryCloud, VersionArgs: []string{"--version"}, Description: "Container runtime", InstallHint: "brew install --cask docker"},
		{Name: "terraform", Category: CategoryCloud, VersionArgs: []string{"--version"}, Description: "Infrastructure as code", InstallHint: "brew install terraform"},
		{Name: "helm", Category: CategoryCloud, VersionArgs: []string{"version", "--short"}, Description: "Kubernetes package manager", InstallHint: "brew install helm"},

		// Development tools
		{Name: "gh", Category: CategoryDevelopment, VersionArgs: []string{"--version"}, Description: "GitHub CLI", InstallHint: "brew install gh"},
		{Name: "nvim", Category: CategoryDevelopment, VersionArgs: []string{"--version"}, Description: "Neovim editor", InstallHint: "brew install neovim"},
		{Name: "tmux", Category: CategoryDevelopment, VersionArgs: []string{"-V"}, Description: "Terminal multiplexer", InstallHint: "brew install tmux"},
		{Name: "fzf", Category: CategoryDevelopment, VersionArgs: []string{"--version"}, Description: "Fuzzy finder", InstallHint: "brew install fzf"},
		{Name: "rg", Category: CategoryDevelopment, VersionArgs: []string{"--version"}, Description: "ripgrep search tool", InstallHint: "brew install ripgrep"},
		{Name: "fd", Category: CategoryDevelopment, VersionArgs: []string{"--version"}, Description: "Find alternative", InstallHint: "brew install fd"},
		{Name: "bat", Category: CategoryDevelopment, VersionArgs: []string{"--version"}, Description: "Cat alternative", InstallHint: "brew install bat"},
		{Name: "eza", Category: CategoryDevelopment, VersionArgs: []string{"--version"}, Description: "ls alternative", InstallHint: "brew install eza"},
	}
}

// Categories returns all category names in display order.
func Categories() []string {
	return []string{
		CategorySystem,
		CategoryLanguages,
		CategoryCloud,
		CategoryDevelopment,
	}
}

// ToolsByCategory returns tools grouped by category.
func ToolsByCategory() map[string][]ToolDefinition {
	result := make(map[string][]ToolDefinition)
	for _, tool := range DefaultRegistry() {
		result[tool.Category] = append(result[tool.Category], tool)
	}
	return result
}

// FindTool looks up a tool by name.
func FindTool(name string) (ToolDefinition, bool) {
	for _, tool := range DefaultRegistry() {
		if tool.Name == name {
			return tool, true
		}
	}
	return ToolDefinition{}, false
}

// ToolNames returns all tool names.
func ToolNames() []string {
	registry := DefaultRegistry()
	names := make([]string, len(registry))
	for i, tool := range registry {
		names[i] = tool.Name
	}
	return names
}
