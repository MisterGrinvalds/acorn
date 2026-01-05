package tools

import (
	"os/exec"
	"strings"
)

// Checker provides tool detection capabilities.
type Checker struct{}

// NewChecker creates a new Checker.
func NewChecker() *Checker {
	return &Checker{}
}

// CheckTool checks if a specific tool is installed and gets its version.
func (c *Checker) CheckTool(name string) ToolStatus {
	def, found := FindTool(name)
	if !found {
		// Unknown tool - just check if command exists
		return c.checkUnknownTool(name)
	}

	status := ToolStatus{
		Name:     name,
		Category: def.Category,
	}

	path, err := exec.LookPath(name)
	if err != nil {
		status.Installed = false
		return status
	}

	status.Installed = true
	status.Path = path
	status.Version = c.getVersion(name, def.VersionArgs)

	return status
}

// checkUnknownTool checks a tool not in the registry.
func (c *Checker) checkUnknownTool(name string) ToolStatus {
	status := ToolStatus{
		Name:     name,
		Category: "Unknown",
	}

	path, err := exec.LookPath(name)
	if err != nil {
		status.Installed = false
		return status
	}

	status.Installed = true
	status.Path = path
	// Try common version flags
	for _, args := range [][]string{{"--version"}, {"version"}, {"-V"}, {"-v"}} {
		if v := c.getVersion(name, args); v != "" {
			status.Version = v
			break
		}
	}

	return status
}

// CheckAll checks all known tools.
func (c *Checker) CheckAll() *StatusResult {
	result := &StatusResult{}

	byCategory := ToolsByCategory()
	for _, catName := range Categories() {
		tools := byCategory[catName]
		cat := ToolCategory{Name: catName}

		for _, def := range tools {
			status := c.CheckTool(def.Name)
			cat.Tools = append(cat.Tools, status)

			result.Summary.Total++
			if status.Installed {
				result.Summary.Installed++
			} else {
				result.Summary.Missing++
			}
		}

		result.Categories = append(result.Categories, cat)
	}

	return result
}

// CheckCategory checks tools in a specific category.
func (c *Checker) CheckCategory(category string) []ToolStatus {
	var results []ToolStatus

	for _, def := range DefaultRegistry() {
		if def.Category == category {
			results = append(results, c.CheckTool(def.Name))
		}
	}

	return results
}

// CheckTools checks specific tools by name.
func (c *Checker) CheckTools(names []string) []ToolStatus {
	results := make([]ToolStatus, len(names))
	for i, name := range names {
		results[i] = c.CheckTool(name)
	}
	return results
}

// GetMissing returns tools that are not installed.
func (c *Checker) GetMissing() []ToolStatus {
	var missing []ToolStatus

	for _, def := range DefaultRegistry() {
		status := c.CheckTool(def.Name)
		if !status.Installed {
			missing = append(missing, status)
		}
	}

	return missing
}

// getVersion attempts to get the version of a tool.
func (c *Checker) getVersion(name string, args []string) string {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		// Some tools write version to stderr
		if exitErr, ok := err.(*exec.ExitError); ok {
			output = exitErr.Stderr
		} else {
			return ""
		}
	}

	version := strings.TrimSpace(string(output))
	// Take first line only
	if idx := strings.Index(version, "\n"); idx > 0 {
		version = version[:idx]
	}

	return version
}

// CommandExists checks if a command exists in PATH.
func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
