package claude

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// MainConfig represents the ~/.claude.json structure.
type MainConfig struct {
	NumStartups   int                `json:"numStartups,omitempty" yaml:"numStartups,omitempty"`
	InstallMethod string             `json:"installMethod,omitempty" yaml:"installMethod,omitempty"`
	AutoUpdates   bool               `json:"autoUpdates,omitempty" yaml:"autoUpdates,omitempty"`
	Projects      map[string]Project `json:"projects,omitempty" yaml:"projects,omitempty"`
}

// Project represents a Claude Code project configuration.
type Project struct {
	AllowedTools           []string               `json:"allowedTools" yaml:"allowedTools"`
	MCPContextURIs         []string               `json:"mcpContextUris" yaml:"mcpContextUris"`
	MCPServers             map[string]MCPServer   `json:"mcpServers" yaml:"mcpServers"`
	EnabledMCPJSONServers  []string               `json:"enabledMcpjsonServers" yaml:"enabledMcpjsonServers"`
	DisabledMCPJSONServers []string               `json:"disabledMcpjsonServers" yaml:"disabledMcpjsonServers"`
	HasTrustDialogAccepted bool                   `json:"hasTrustDialogAccepted" yaml:"hasTrustDialogAccepted"`
	LastCost               float64                `json:"lastCost,omitempty" yaml:"lastCost,omitempty"`
	LastSessionID          string                 `json:"lastSessionId,omitempty" yaml:"lastSessionId,omitempty"`
	LastModelUsage         map[string]interface{} `json:"lastModelUsage,omitempty" yaml:"lastModelUsage,omitempty"`
}

// MCPServer represents an MCP server configuration.
type MCPServer struct {
	Type    string   `json:"type" yaml:"type"`
	URL     string   `json:"url,omitempty" yaml:"url,omitempty"`
	Command string   `json:"command,omitempty" yaml:"command,omitempty"`
	Args    []string `json:"args,omitempty" yaml:"args,omitempty"`
}

// MCPConfig represents the local .mcp.json structure.
type MCPConfig struct {
	MCPServers map[string]MCPServer `json:"mcpServers" yaml:"mcpServers"`
}

// ProjectView is a simplified view of a project for display.
type ProjectView struct {
	Path    string  `json:"path" yaml:"path"`
	Name    string  `json:"name" yaml:"name"`
	Cost    float64 `json:"cost" yaml:"cost"`
	Trusted bool    `json:"trusted" yaml:"trusted"`
}

// ProjectsView is the list of projects for display.
type ProjectsView struct {
	Projects []ProjectView `json:"projects" yaml:"projects"`
}

// MCPServerView is a simplified view of an MCP server for display.
type MCPServerView struct {
	Name    string `json:"name" yaml:"name"`
	Type    string `json:"type" yaml:"type"`
	URL     string `json:"url" yaml:"url"`
	Project string `json:"project,omitempty" yaml:"project,omitempty"`
}

// MCPView is the list of MCP servers for display.
type MCPView struct {
	Servers   []MCPServerView `json:"servers" yaml:"servers"`
	LocalFile string          `json:"local_file,omitempty" yaml:"local_file,omitempty"`
	HasLocal  bool            `json:"has_local" yaml:"has_local"`
}

// CommandView represents a custom command for display.
type CommandView struct {
	Name   string `json:"name" yaml:"name"`
	Path   string `json:"path" yaml:"path"`
	Source string `json:"source" yaml:"source"` // "user" or "project"
}

// CommandsView is the list of custom commands for display.
type CommandsView struct {
	UserCommands    []CommandView `json:"user_commands" yaml:"user_commands"`
	ProjectCommands []CommandView `json:"project_commands" yaml:"project_commands"`
}

// GetMainConfig reads the main ~/.claude.json config.
func (h *Helper) GetMainConfig() (*MainConfig, error) {
	if !h.FileExists(h.paths.Config) {
		return nil, fmt.Errorf("config file not found: %s", h.paths.Config)
	}

	var config MainConfig
	if err := h.ReadJSONFile(h.paths.Config, &config); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &config, nil
}

// GetProjects returns a list of projects for display.
func (h *Helper) GetProjects() (*ProjectsView, error) {
	config, err := h.GetMainConfig()
	if err != nil {
		return nil, err
	}

	view := &ProjectsView{
		Projects: []ProjectView{},
	}

	for path, project := range config.Projects {
		if !project.HasTrustDialogAccepted {
			continue
		}
		view.Projects = append(view.Projects, ProjectView{
			Path:    path,
			Name:    filepath.Base(path),
			Cost:    project.LastCost,
			Trusted: project.HasTrustDialogAccepted,
		})
	}

	// Sort by cost descending
	sort.Slice(view.Projects, func(i, j int) bool {
		return view.Projects[i].Cost > view.Projects[j].Cost
	})

	return view, nil
}

// GetMCPServers returns a list of MCP servers for display.
func (h *Helper) GetMCPServers() (*MCPView, error) {
	config, err := h.GetMainConfig()
	if err != nil {
		return nil, err
	}

	view := &MCPView{
		Servers: []MCPServerView{},
	}

	// Get servers from projects
	for path, project := range config.Projects {
		for name, server := range project.MCPServers {
			url := server.URL
			if url == "" {
				url = server.Command
			}
			view.Servers = append(view.Servers, MCPServerView{
				Name:    name,
				Type:    server.Type,
				URL:     url,
				Project: filepath.Base(path),
			})
		}
	}

	// Check for local .mcp.json
	localMCPPath := ".mcp.json"
	if h.FileExists(localMCPPath) {
		view.HasLocal = true
		view.LocalFile = localMCPPath

		var mcpConfig MCPConfig
		if err := h.ReadJSONFile(localMCPPath, &mcpConfig); err == nil {
			for name, server := range mcpConfig.MCPServers {
				url := server.URL
				if url == "" {
					url = server.Command
				}
				view.Servers = append(view.Servers, MCPServerView{
					Name:    name,
					Type:    server.Type,
					URL:     url,
					Project: "(local)",
				})
			}
		}
	}

	// Sort by name
	sort.Slice(view.Servers, func(i, j int) bool {
		return view.Servers[i].Name < view.Servers[j].Name
	})

	return view, nil
}

// AddMCPServer adds an MCP server to the local .mcp.json.
func (h *Helper) AddMCPServer(name, url, serverType string) error {
	mcpPath := ".mcp.json"

	var mcpConfig MCPConfig
	if h.FileExists(mcpPath) {
		if err := h.ReadJSONFile(mcpPath, &mcpConfig); err != nil {
			return fmt.Errorf("failed to read .mcp.json: %w", err)
		}
	} else {
		mcpConfig = MCPConfig{
			MCPServers: make(map[string]MCPServer),
		}
	}

	if mcpConfig.MCPServers == nil {
		mcpConfig.MCPServers = make(map[string]MCPServer)
	}

	mcpConfig.MCPServers[name] = MCPServer{
		Type: serverType,
		URL:  url,
	}

	if err := h.WriteJSONFile(mcpPath, mcpConfig); err != nil {
		return fmt.Errorf("failed to write .mcp.json: %w", err)
	}

	return nil
}

// GetCommands returns a list of custom commands for display.
func (h *Helper) GetCommands() (*CommandsView, error) {
	view := &CommandsView{
		UserCommands:    []CommandView{},
		ProjectCommands: []CommandView{},
	}

	// Get user commands
	userCmdDir := h.paths.CommandsDir
	if h.DirExists(userCmdDir) {
		entries, err := os.ReadDir(userCmdDir)
		if err == nil {
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
					name := strings.TrimSuffix(entry.Name(), ".md")
					view.UserCommands = append(view.UserCommands, CommandView{
						Name:   "/" + name,
						Path:   filepath.Join(userCmdDir, entry.Name()),
						Source: "user",
					})
				}
			}
		}
	}

	// Get project commands
	projectCmdDir := ".claude/commands"
	if h.DirExists(projectCmdDir) {
		entries, err := os.ReadDir(projectCmdDir)
		if err == nil {
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
					name := strings.TrimSuffix(entry.Name(), ".md")
					view.ProjectCommands = append(view.ProjectCommands, CommandView{
						Name:   "/" + name,
						Path:   filepath.Join(projectCmdDir, entry.Name()),
						Source: "project",
					})
				}
			}
		}
	}

	// Sort by name
	sort.Slice(view.UserCommands, func(i, j int) bool {
		return view.UserCommands[i].Name < view.UserCommands[j].Name
	})
	sort.Slice(view.ProjectCommands, func(i, j int) bool {
		return view.ProjectCommands[i].Name < view.ProjectCommands[j].Name
	})

	return view, nil
}
