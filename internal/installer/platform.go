package installer

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// DetectPlatform detects the current platform.
func DetectPlatform() *Platform {
	p := &Platform{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	switch p.OS {
	case "darwin":
		p.PackageManager = "brew"
	case "linux":
		p.detectLinuxDistro()
	case "windows":
		p.PackageManager = "winget"
	}

	return p
}

// detectLinuxDistro detects the Linux distribution and package manager.
func (p *Platform) detectLinuxDistro() {
	// Read /etc/os-release for distro info
	data, err := os.ReadFile("/etc/os-release")
	if err == nil {
		p.parseOSRelease(string(data))
	}

	// Determine package manager based on what's available
	p.detectPackageManager()
}

// parseOSRelease parses /etc/os-release content.
func (p *Platform) parseOSRelease(content string) {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "ID=") {
			p.Distro = strings.Trim(strings.TrimPrefix(line, "ID="), "\"")
		}
		if strings.HasPrefix(line, "ID_LIKE=") {
			p.DistroFamily = strings.Trim(strings.TrimPrefix(line, "ID_LIKE="), "\"")
			// Take first word if multiple (e.g., "debian ubuntu")
			if idx := strings.Index(p.DistroFamily, " "); idx > 0 {
				p.DistroFamily = p.DistroFamily[:idx]
			}
		}
	}

	// Set distro family from distro if not set
	if p.DistroFamily == "" {
		switch p.Distro {
		case "ubuntu", "pop", "mint", "elementary":
			p.DistroFamily = "debian"
		case "fedora", "centos", "rhel", "rocky", "alma":
			p.DistroFamily = "rhel"
		case "arch", "manjaro", "endeavouros":
			p.DistroFamily = "arch"
		}
	}
}

// detectPackageManager detects the available package manager.
func (p *Platform) detectPackageManager() {
	// Check in order of preference
	managers := []string{"apt", "dnf", "yum", "pacman", "zypper", "brew"}

	for _, mgr := range managers {
		if commandExists(mgr) {
			p.PackageManager = mgr
			return
		}
	}
}

// commandExists checks if a command exists in PATH.
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// SupportsPackageManager returns true if the platform has a known package manager.
func (p *Platform) SupportsPackageManager() bool {
	return p.PackageManager != ""
}

// String returns a human-readable platform description.
func (p *Platform) String() string {
	var parts []string
	parts = append(parts, p.OS)

	if p.Distro != "" {
		parts = append(parts, p.Distro)
	}

	parts = append(parts, p.Arch)

	if p.PackageManager != "" {
		parts = append(parts, "("+p.PackageManager+")")
	}

	return strings.Join(parts, " ")
}
