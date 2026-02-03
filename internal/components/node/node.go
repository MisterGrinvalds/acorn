// Package node provides Node.js, NVM, and pnpm management functionality.
package node

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Status represents Node.js ecosystem status.
type Status struct {
	NodeInstalled  bool   `json:"node_installed" yaml:"node_installed"`
	NodeVersion    string `json:"node_version,omitempty" yaml:"node_version,omitempty"`
	NpmInstalled   bool   `json:"npm_installed" yaml:"npm_installed"`
	NpmVersion     string `json:"npm_version,omitempty" yaml:"npm_version,omitempty"`
	NvmInstalled   bool   `json:"nvm_installed" yaml:"nvm_installed"`
	NvmDir         string `json:"nvm_dir,omitempty" yaml:"nvm_dir,omitempty"`
	PnpmInstalled  bool   `json:"pnpm_installed" yaml:"pnpm_installed"`
	PnpmVersion    string `json:"pnpm_version,omitempty" yaml:"pnpm_version,omitempty"`
	PnpmHome       string `json:"pnpm_home,omitempty" yaml:"pnpm_home,omitempty"`
	NpmCacheDir    string `json:"npm_cache_dir,omitempty" yaml:"npm_cache_dir,omitempty"`
}

// NodeModulesInfo represents node_modules directory info.
type NodeModulesInfo struct {
	Path string `json:"path" yaml:"path"`
	Size string `json:"size" yaml:"size"`
}

// PackageManager represents detected package manager.
type PackageManager struct {
	Name     string `json:"name" yaml:"name"`
	LockFile string `json:"lock_file,omitempty" yaml:"lock_file,omitempty"`
}

// Helper provides Node.js helper operations.
type Helper struct {
	verbose bool
	dryRun  bool
}

// NewHelper creates a new Node.js Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	return &Helper{
		verbose: verbose,
		dryRun:  dryRun,
	}
}

// GetNvmDir returns the NVM directory path.
func (h *Helper) GetNvmDir() string {
	if nvmDir := os.Getenv("NVM_DIR"); nvmDir != "" {
		return nvmDir
	}
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		home, _ := os.UserHomeDir()
		dataHome = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(dataHome, "nvm")
}

// GetPnpmHome returns the pnpm home directory path.
func (h *Helper) GetPnpmHome() string {
	if pnpmHome := os.Getenv("PNPM_HOME"); pnpmHome != "" {
		return pnpmHome
	}
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		home, _ := os.UserHomeDir()
		dataHome = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(dataHome, "pnpm")
}

// GetNpmCacheDir returns the npm cache directory path.
func (h *Helper) GetNpmCacheDir() string {
	if cacheDir := os.Getenv("npm_config_cache"); cacheDir != "" {
		return cacheDir
	}
	cacheHome := os.Getenv("XDG_CACHE_HOME")
	if cacheHome == "" {
		home, _ := os.UserHomeDir()
		cacheHome = filepath.Join(home, ".cache")
	}
	return filepath.Join(cacheHome, "npm")
}

// GetStatus returns Node.js ecosystem status.
func (h *Helper) GetStatus() *Status {
	status := &Status{
		NvmDir:      h.GetNvmDir(),
		PnpmHome:    h.GetPnpmHome(),
		NpmCacheDir: h.GetNpmCacheDir(),
	}

	// Check Node
	if out, err := exec.Command("node", "--version").Output(); err == nil {
		status.NodeInstalled = true
		status.NodeVersion = strings.TrimSpace(string(out))
	}

	// Check npm
	if out, err := exec.Command("npm", "--version").Output(); err == nil {
		status.NpmInstalled = true
		status.NpmVersion = strings.TrimSpace(string(out))
	}

	// Check NVM (check if nvm.sh exists)
	nvmScript := filepath.Join(status.NvmDir, "nvm.sh")
	if _, err := os.Stat(nvmScript); err == nil {
		status.NvmInstalled = true
	}

	// Check pnpm
	if out, err := exec.Command("pnpm", "--version").Output(); err == nil {
		status.PnpmInstalled = true
		status.PnpmVersion = strings.TrimSpace(string(out))
	}

	return status
}

// DetectPackageManager detects the package manager from lock files.
func (h *Helper) DetectPackageManager() *PackageManager {
	pm := &PackageManager{Name: "npm"}

	if _, err := os.Stat("pnpm-lock.yaml"); err == nil {
		pm.Name = "pnpm"
		pm.LockFile = "pnpm-lock.yaml"
	} else if _, err := os.Stat("yarn.lock"); err == nil {
		pm.Name = "yarn"
		pm.LockFile = "yarn.lock"
	} else if _, err := os.Stat("package-lock.json"); err == nil {
		pm.Name = "npm"
		pm.LockFile = "package-lock.json"
	}

	return pm
}

// FindNodeModules finds all node_modules directories.
func (h *Helper) FindNodeModules(root string) ([]NodeModulesInfo, error) {
	if root == "" {
		root = "."
	}

	var modules []NodeModulesInfo

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}
		if info.IsDir() && info.Name() == "node_modules" {
			// Get size
			cmd := exec.Command("du", "-sh", path)
			out, _ := cmd.Output()
			size := "unknown"
			if parts := strings.Fields(string(out)); len(parts) > 0 {
				size = parts[0]
			}

			modules = append(modules, NodeModulesInfo{
				Path: path,
				Size: size,
			})

			return filepath.SkipDir // Don't recurse into node_modules
		}
		return nil
	})

	return modules, err
}

// CleanNodeModules removes node_modules and reinstalls.
func (h *Helper) CleanNodeModules() error {
	if _, err := os.Stat("node_modules"); os.IsNotExist(err) {
		return fmt.Errorf("no node_modules directory found")
	}

	if h.dryRun {
		fmt.Println("[dry-run] would remove node_modules and reinstall")
		return nil
	}

	fmt.Println("Removing node_modules...")
	if err := os.RemoveAll("node_modules"); err != nil {
		return fmt.Errorf("failed to remove node_modules: %w", err)
	}

	fmt.Println("Reinstalling dependencies...")
	pm := h.DetectPackageManager()

	var cmd *exec.Cmd
	switch pm.Name {
	case "pnpm":
		cmd = exec.Command("pnpm", "install")
	case "yarn":
		cmd = exec.Command("yarn", "install")
	default:
		cmd = exec.Command("npm", "install")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// CleanAllNodeModules removes all node_modules in directory tree.
func (h *Helper) CleanAllNodeModules(root string, force bool) (int, error) {
	modules, err := h.FindNodeModules(root)
	if err != nil {
		return 0, err
	}

	if len(modules) == 0 {
		return 0, nil
	}

	if !force {
		return len(modules), fmt.Errorf("found %d node_modules directories, use --force to remove", len(modules))
	}

	if h.dryRun {
		for _, m := range modules {
			fmt.Printf("[dry-run] would remove %s\n", m.Path)
		}
		return len(modules), nil
	}

	count := 0
	for _, m := range modules {
		if err := os.RemoveAll(m.Path); err == nil {
			count++
		}
	}

	return count, nil
}

// GetNvmVersions returns installed NVM versions.
func (h *Helper) GetNvmVersions() ([]string, error) {
	nvmDir := h.GetNvmDir()
	versionsDir := filepath.Join(nvmDir, "versions", "node")

	if _, err := os.Stat(versionsDir); os.IsNotExist(err) {
		return nil, nil
	}

	entries, err := os.ReadDir(versionsDir)
	if err != nil {
		return nil, err
	}

	var versions []string
	for _, entry := range entries {
		if entry.IsDir() {
			versions = append(versions, entry.Name())
		}
	}

	return versions, nil
}

// GetCurrentNodeVersion returns the currently active Node version.
func (h *Helper) GetCurrentNodeVersion() string {
	out, err := exec.Command("node", "--version").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// InstallNvm installs NVM.
func (h *Helper) InstallNvm() error {
	nvmDir := h.GetNvmDir()
	nvmScript := filepath.Join(nvmDir, "nvm.sh")

	if _, err := os.Stat(nvmScript); err == nil {
		return fmt.Errorf("NVM already installed at %s", nvmDir)
	}

	if h.dryRun {
		fmt.Println("[dry-run] would install NVM")
		return nil
	}

	fmt.Println("Installing NVM...")

	// Download and run installer
	cmd := exec.Command("bash", "-c",
		`curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.1/install.sh | bash`)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), fmt.Sprintf("NVM_DIR=%s", nvmDir))

	return cmd.Run()
}

// InstallPnpm installs pnpm globally.
func (h *Helper) InstallPnpm() error {
	if _, err := exec.LookPath("pnpm"); err == nil {
		return fmt.Errorf("pnpm already installed")
	}

	if h.dryRun {
		fmt.Println("[dry-run] would install pnpm")
		return nil
	}

	fmt.Println("Installing pnpm...")
	cmd := exec.Command("npm", "install", "-g", "pnpm")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetNpmCacheSize returns the npm cache size.
func (h *Helper) GetNpmCacheSize() string {
	cacheDir := h.GetNpmCacheDir()
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		return "0"
	}

	cmd := exec.Command("du", "-sh", cacheDir)
	out, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	parts := strings.Fields(string(out))
	if len(parts) > 0 {
		return parts[0]
	}
	return "unknown"
}

// CleanNpmCache cleans the npm cache.
func (h *Helper) CleanNpmCache() error {
	if h.dryRun {
		fmt.Println("[dry-run] would clean npm cache")
		return nil
	}

	cmd := exec.Command("npm", "cache", "clean", "--force")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
