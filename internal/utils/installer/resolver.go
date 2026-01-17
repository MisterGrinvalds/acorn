package installer

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/mistergrinvalds/acorn/internal/utils/config"
)

// Resolver handles prerequisite resolution for installation.
type Resolver struct {
	loader   *config.ComponentLoader
	platform *Platform
	visited  map[string]bool
	resolved []PlannedTool
}

// NewResolver creates a new Resolver.
func NewResolver(platform *Platform) *Resolver {
	return &Resolver{
		loader:   config.NewComponentLoader(),
		platform: platform,
		visited:  make(map[string]bool),
	}
}

// BuildPlan builds a complete installation plan with prerequisites resolved.
func (r *Resolver) BuildPlan(component string, cfg *config.InstallConfig) (*InstallPlan, error) {
	plan := &InstallPlan{
		Component: component,
		Platform:  *r.platform,
	}

	// Reset state
	r.visited = make(map[string]bool)
	r.resolved = nil

	// Resolve all tools
	for _, tool := range cfg.Tools {
		if err := r.resolveTool(tool, "direct"); err != nil {
			return nil, err
		}
	}

	// Separate into prerequisites and direct tools
	for _, pt := range r.resolved {
		if pt.Reason == "prerequisite" {
			plan.Prerequisites = append(plan.Prerequisites, pt)
		} else {
			plan.Tools = append(plan.Tools, pt)
		}
	}

	return plan, nil
}

// resolveTool resolves a tool and its prerequisites recursively.
func (r *Resolver) resolveTool(tool config.ToolInstall, reason string) error {
	// Check for cycles
	if r.visited[tool.Name] {
		return nil // Already processed
	}
	r.visited[tool.Name] = true

	// Check if already installed
	installed, version := r.checkInstalled(tool.Check)

	// Resolve prerequisites first (depth-first)
	for _, req := range tool.Requires {
		if err := r.resolveRequirement(req); err != nil {
			return fmt.Errorf("prerequisite %s: %w", req, err)
		}
	}

	// Get the appropriate install method for this platform
	method, found := SelectMethod(tool.Methods, r.platform)
	if !found && !installed {
		return fmt.Errorf("no install method for %s on %s", tool.Name, r.platform.OS)
	}

	// Add to resolved list
	r.resolved = append(r.resolved, PlannedTool{
		Name:             tool.Name,
		Description:      tool.Description,
		Method:           method,
		AlreadyInstalled: installed,
		Version:          version,
		Reason:           reason,
		PostInstall:      tool.PostInstall,
	})

	return nil
}

// resolveRequirement resolves a prerequisite requirement.
// Format: "command" or "component:command"
func (r *Resolver) resolveRequirement(req string) error {
	component, command := parseRequirement(req)

	// Check if command already exists
	if commandExists(command) {
		return nil // Already satisfied
	}

	// If component specified, try to load its install config
	if component != "" {
		return r.resolveFromComponent(component, command)
	}

	// Just a command name - try to find it in component configs
	// For now, we'll return an error asking to install it manually
	return fmt.Errorf("%s not found and no component specified (use component:command format)", command)
}

// resolveFromComponent resolves a tool from another component's install config.
func (r *Resolver) resolveFromComponent(component, command string) error {
	cfg, err := r.loadInstallConfig(component)
	if err != nil {
		return fmt.Errorf("cannot load %s config: %w", component, err)
	}

	// Find the tool in the component
	for _, tool := range cfg.Tools {
		if tool.Name == command {
			return r.resolveTool(tool, "prerequisite")
		}
	}

	return fmt.Errorf("tool %s not found in %s install config", command, component)
}

// loadInstallConfig loads the install config for a component.
func (r *Resolver) loadInstallConfig(component string) (*config.InstallConfig, error) {
	cfg := &config.BaseConfig{}
	if err := r.loader.Load(component, cfg); err != nil {
		return nil, err
	}
	return &cfg.Install, nil
}

// checkInstalled checks if a tool is installed using the check command.
func (r *Resolver) checkInstalled(check string) (bool, string) {
	if check == "" {
		return false, ""
	}

	// Parse the check command
	// Common formats: "command -v <name>", "<name> --version"
	parts := strings.Fields(check)
	if len(parts) == 0 {
		return false, ""
	}

	// Handle "command -v <name>" format
	if parts[0] == "command" && len(parts) >= 3 && parts[1] == "-v" {
		return commandExists(parts[2]), ""
	}

	// Try executing the check command
	cmd := exec.Command("sh", "-c", check)
	output, err := cmd.Output()
	if err != nil {
		return false, ""
	}

	// Return first line as version info
	version := strings.TrimSpace(string(output))
	if idx := strings.Index(version, "\n"); idx > 0 {
		version = version[:idx]
	}

	return true, version
}

// parseRequirement parses a requirement string into component and command.
// Format: "command" or "component:command"
func parseRequirement(req string) (component, command string) {
	if idx := strings.Index(req, ":"); idx > 0 {
		return req[:idx], req[idx+1:]
	}
	return "", req
}
