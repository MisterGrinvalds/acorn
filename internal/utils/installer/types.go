// Package installer provides component installation functionality.
package installer

import (
	"time"

	"github.com/mistergrinvalds/acorn/internal/utils/config"
)

// Platform represents the detected system platform.
type Platform struct {
	OS             string // darwin, linux, windows
	Distro         string // debian, ubuntu, fedora, arch (linux only)
	DistroFamily   string // debian, rhel, arch
	Arch           string // amd64, arm64
	PackageManager string // brew, apt, dnf, pacman
}

// GetMethodKeys returns method keys to try in order of specificity.
// Most specific first: linux/ubuntu, linux/debian, linux, then fallback.
func (p *Platform) GetMethodKeys() []string {
	keys := []string{}

	// Most specific first
	if p.Distro != "" {
		keys = append(keys, p.OS+"/"+p.Distro)
	}
	if p.DistroFamily != "" && p.DistroFamily != p.Distro {
		keys = append(keys, p.OS+"/"+p.DistroFamily)
	}
	keys = append(keys, p.OS)

	return keys
}

// InstallPlan represents the planned installation steps.
type InstallPlan struct {
	Component     string
	Platform      Platform
	Prerequisites []PlannedTool // Tools to install first (dependencies)
	Tools         []PlannedTool // Direct tools to install
	DryRun        bool
}

// PlannedTool represents a tool to be installed.
type PlannedTool struct {
	Name             string
	Description      string
	Method           config.InstallMethod
	AlreadyInstalled bool
	Version          string // Current version if installed
	Reason           string // "direct" or "prerequisite"
	PostInstall      config.PostInstallConfig
}

// TotalTools returns the total number of tools in the plan.
func (p *InstallPlan) TotalTools() int {
	return len(p.Prerequisites) + len(p.Tools)
}

// PendingTools returns tools that need to be installed.
func (p *InstallPlan) PendingTools() []PlannedTool {
	var pending []PlannedTool
	for _, t := range p.Prerequisites {
		if !t.AlreadyInstalled {
			pending = append(pending, t)
		}
	}
	for _, t := range p.Tools {
		if !t.AlreadyInstalled {
			pending = append(pending, t)
		}
	}
	return pending
}

// InstallResult represents the result of an installation.
type InstallResult struct {
	Component string
	Success   bool
	Tools     []ToolResult
	Duration  time.Duration
	DryRun    bool
}

// ToolResult represents the result for a single tool installation.
type ToolResult struct {
	Name       string
	Success    bool
	Skipped    bool
	SkipReason string
	Version    string
	Error      error
	Duration   time.Duration
}

// Summary returns a summary of the installation result.
func (r *InstallResult) Summary() (installed, skipped, failed int) {
	for _, t := range r.Tools {
		switch {
		case t.Skipped:
			skipped++
		case t.Success:
			installed++
		default:
			failed++
		}
	}
	return
}

// InstallType constants for installation methods.
const (
	InstallTypeBrew   = "brew"
	InstallTypeApt    = "apt"
	InstallTypeDnf    = "dnf"
	InstallTypePacman = "pacman"
	InstallTypeNpm    = "npm"
	InstallTypePip    = "pip"
	InstallTypeCargo  = "cargo"
	InstallTypeGo     = "go"
	InstallTypeCurl   = "curl"
	InstallTypeBinary = "binary"
)
