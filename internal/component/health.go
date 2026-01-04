package component

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// HealthStatus represents the health state of a component.
type HealthStatus string

const (
	StatusHealthy HealthStatus = "healthy"
	StatusWarning HealthStatus = "warning"
	StatusError   HealthStatus = "error"
)

// HealthCheck represents the result of a health check.
type HealthCheck struct {
	Component *Component
	Status    HealthStatus
	Issues    []string
	Warnings  []string
}

// CheckHealth performs a comprehensive health check on a component.
func CheckHealth(comp *Component) *HealthCheck {
	hc := &HealthCheck{
		Component: comp,
		Status:    StatusHealthy,
		Issues:    []string{},
		Warnings:  []string{},
	}

	// Check YAML validity (already validated if we loaded it)
	if comp.Name == "" {
		hc.addError("missing required field: name")
	}
	if comp.Version == "" {
		hc.addWarning("missing version field")
	}
	if comp.Description == "" {
		hc.addWarning("missing description field")
	}

	// Check required tools
	for _, tool := range comp.Requires.Tools {
		if !commandExists(tool) {
			hc.addWarning(fmt.Sprintf("required tool not installed: %s", tool))
		}
	}

	// Check shell file syntax
	for _, shellFile := range comp.ShellFiles() {
		if err := checkShellSyntax(shellFile); err != nil {
			hc.addError(fmt.Sprintf("syntax error in %s: %v", shellFile, err))
		}
	}

	// Check platform compatibility
	if !comp.SupportsCurrentPlatform(runtime.GOOS) {
		hc.addWarning(fmt.Sprintf("not compatible with platform: %s", runtime.GOOS))
	}

	// Check config files exist
	for _, cfg := range comp.Config.Files {
		sourcePath := comp.Path + "/" + cfg.Source
		if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
			hc.addError(fmt.Sprintf("config source missing: %s", cfg.Source))
		}
	}

	return hc
}

// addError adds an error issue and sets status to error.
func (hc *HealthCheck) addError(issue string) {
	hc.Issues = append(hc.Issues, issue)
	hc.Status = StatusError
}

// addWarning adds a warning and sets status to warning (if not already error).
func (hc *HealthCheck) addWarning(warning string) {
	hc.Warnings = append(hc.Warnings, warning)
	if hc.Status == StatusHealthy {
		hc.Status = StatusWarning
	}
}

// IsHealthy returns true if status is healthy.
func (hc *HealthCheck) IsHealthy() bool {
	return hc.Status == StatusHealthy
}

// commandExists checks if a command is available in PATH.
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// checkShellSyntax checks bash syntax of a shell file.
func checkShellSyntax(path string) error {
	cmd := exec.Command("bash", "-n", path)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("syntax check failed: %w", err)
	}
	return nil
}
