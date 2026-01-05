package tools

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// PackageManager represents a system package manager.
type PackageManager string

const (
	PMBrew   PackageManager = "brew"
	PMApt    PackageManager = "apt"
	PMDnf    PackageManager = "dnf"
	PMPacman PackageManager = "pacman"
	PMNone   PackageManager = "none"
)

// Updater handles tool updates.
type Updater struct {
	dryRun  bool
	verbose bool
}

// NewUpdater creates a new Updater.
func NewUpdater(dryRun, verbose bool) *Updater {
	return &Updater{
		dryRun:  dryRun,
		verbose: verbose,
	}
}

// DetectPackageManager detects the system package manager.
func (u *Updater) DetectPackageManager() PackageManager {
	if CommandExists("brew") {
		return PMBrew
	}
	if CommandExists("apt-get") {
		return PMApt
	}
	if CommandExists("dnf") {
		return PMDnf
	}
	if CommandExists("pacman") {
		return PMPacman
	}
	return PMNone
}

// UpdateSystem runs the system package manager update.
func (u *Updater) UpdateSystem() error {
	pm := u.DetectPackageManager()

	switch pm {
	case PMBrew:
		return u.runCmd("brew", "update")
	case PMApt:
		return u.runCmd("sudo", "apt-get", "update")
	case PMDnf:
		return u.runCmd("sudo", "dnf", "check-update")
	case PMPacman:
		return u.runCmd("sudo", "pacman", "-Sy")
	default:
		return fmt.Errorf("no supported package manager found")
	}
}

// UpgradeSystem upgrades all packages via the system package manager.
func (u *Updater) UpgradeSystem() error {
	pm := u.DetectPackageManager()

	switch pm {
	case PMBrew:
		if err := u.runCmd("brew", "update"); err != nil {
			return err
		}
		return u.runCmd("brew", "upgrade")
	case PMApt:
		if err := u.runCmd("sudo", "apt-get", "update"); err != nil {
			return err
		}
		return u.runCmd("sudo", "apt-get", "upgrade", "-y")
	case PMDnf:
		return u.runCmd("sudo", "dnf", "upgrade", "-y")
	case PMPacman:
		return u.runCmd("sudo", "pacman", "-Syu")
	default:
		return fmt.Errorf("no supported package manager found")
	}
}

// InstallTool installs a tool using the appropriate method.
func (u *Updater) InstallTool(name string) error {
	def, found := FindTool(name)
	if !found {
		return fmt.Errorf("tool %q not in registry, cannot determine install method", name)
	}

	// Parse install hint to determine method
	hint := def.InstallHint
	if hint == "" {
		return fmt.Errorf("no install hint for tool %q", name)
	}

	// Execute the install command
	parts := strings.Fields(hint)
	if len(parts) == 0 {
		return fmt.Errorf("invalid install hint for tool %q", name)
	}

	return u.runCmd(parts[0], parts[1:]...)
}

// UpgradeBash upgrades to modern bash on macOS.
func (u *Updater) UpgradeBash() error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("upgrade-bash is only supported on macOS")
	}

	brewBash := "/opt/homebrew/bin/bash"
	if !fileExists(brewBash) {
		brewBash = "/usr/local/bin/bash"
	}

	// Install bash if not present
	if !CommandExists(brewBash) {
		fmt.Println("Installing modern bash via Homebrew...")
		if err := u.runCmd("brew", "install", "bash"); err != nil {
			return err
		}
	}

	// Get versions
	brewVersion, _ := getCommandOutput(brewBash, "--version")
	sysVersion, _ := getCommandOutput("/bin/bash", "--version")

	fmt.Printf("Homebrew bash: %s\n", firstLine(brewVersion))
	fmt.Printf("System bash:   %s\n", firstLine(sysVersion))

	// Check if already in /etc/shells
	shells, _ := os.ReadFile("/etc/shells")
	if !strings.Contains(string(shells), brewBash) {
		fmt.Printf("\nAdding %s to /etc/shells...\n", brewBash)
		if !u.dryRun {
			if err := u.runCmd("sudo", "tee", "-a", "/etc/shells"); err != nil {
				return fmt.Errorf("failed to add to /etc/shells: %w", err)
			}
		}
	}

	// Check current shell
	currentShell := os.Getenv("SHELL")
	if currentShell != brewBash {
		fmt.Printf("\nCurrent shell: %s\n", currentShell)
		fmt.Printf("To set %s as default: chsh -s %s\n", brewBash, brewBash)
	} else {
		fmt.Println("\nAlready using Homebrew bash as default shell.")
	}

	return nil
}

// runCmd executes a command, optionally as dry-run.
func (u *Updater) runCmd(name string, args ...string) error {
	if u.dryRun {
		fmt.Printf("[dry-run] would run: %s %s\n", name, strings.Join(args, " "))
		return nil
	}

	if u.verbose {
		fmt.Printf("Running: %s %s\n", name, strings.Join(args, " "))
	}

	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// getCommandOutput runs a command and returns its output.
func getCommandOutput(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}

// firstLine returns the first line of a string.
func firstLine(s string) string {
	if idx := strings.Index(s, "\n"); idx > 0 {
		return s[:idx]
	}
	return s
}

// fileExists checks if a file exists.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
