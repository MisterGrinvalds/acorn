// Package python provides Python development helper functionality.
package python

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// VenvInfo contains virtual environment information.
type VenvInfo struct {
	Name      string `json:"name" yaml:"name"`
	Path      string `json:"path" yaml:"path"`
	Python    string `json:"python,omitempty" yaml:"python,omitempty"`
	Active    bool   `json:"active" yaml:"active"`
	CreatedBy string `json:"created_by,omitempty" yaml:"created_by,omitempty"`
}

// EnvInfo contains Python environment information.
type EnvInfo struct {
	Python       string `json:"python" yaml:"python"`
	Version      string `json:"version" yaml:"version"`
	Pip          string `json:"pip,omitempty" yaml:"pip,omitempty"`
	UV           string `json:"uv,omitempty" yaml:"uv,omitempty"`
	VirtualEnv   string `json:"virtual_env,omitempty" yaml:"virtual_env,omitempty"`
	EnvsLocation string `json:"envs_location,omitempty" yaml:"envs_location,omitempty"`
}

// FastAPIDeps are the FastAPI development dependencies.
var FastAPIDeps = []string{
	"fastapi",
	"uvicorn",
	"python-multipart",
	"pytest",
	"httpx",
	"pytest-asyncio",
	"ruff",
	"python-dotenv",
}

// DevToolsDeps are common development tools.
var DevToolsDeps = []string{
	"ruff",
	"mypy",
	"pytest",
	"pytest-cov",
	"pre-commit",
}

// IPythonDeps are IPython dependencies.
var IPythonDeps = []string{
	"ipython",
	"rich",
}

// Helper provides Python development helper operations.
type Helper struct {
	verbose bool
	dryRun  bool
}

// NewHelper creates a new Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	return &Helper{
		verbose: verbose,
		dryRun:  dryRun,
	}
}

// HasUV checks if UV is installed.
func (h *Helper) HasUV() bool {
	_, err := exec.LookPath("uv")
	return err == nil
}

// HasPython checks if Python3 is installed.
func (h *Helper) HasPython() bool {
	_, err := exec.LookPath("python3")
	return err == nil
}

// CreateVenv creates a virtual environment.
func (h *Helper) CreateVenv(name string) (*VenvInfo, error) {
	if name == "" {
		name = ".venv"
	}

	absPath, err := filepath.Abs(name)
	if err != nil {
		return nil, err
	}

	var createdBy string
	if h.HasUV() {
		if err := h.run("uv", "venv", name); err != nil {
			return nil, fmt.Errorf("uv venv failed: %w", err)
		}
		createdBy = "uv"
	} else if h.HasPython() {
		if err := h.run("python3", "-m", "venv", name); err != nil {
			return nil, fmt.Errorf("python3 -m venv failed: %w", err)
		}
		createdBy = "python3"
	} else {
		return nil, fmt.Errorf("neither uv nor python3 found")
	}

	// Get Python version in the venv
	pythonPath := filepath.Join(absPath, "bin", "python")
	var pythonVersion string
	if _, err := os.Stat(pythonPath); err == nil {
		cmd := exec.Command(pythonPath, "--version")
		if out, err := cmd.Output(); err == nil {
			pythonVersion = strings.TrimSpace(string(out))
		}
	}

	return &VenvInfo{
		Name:      name,
		Path:      absPath,
		Python:    pythonVersion,
		Active:    false,
		CreatedBy: createdBy,
	}, nil
}

// InitProject initializes a new Python project with UV.
func (h *Helper) InitProject(name string) error {
	if !h.HasUV() {
		return fmt.Errorf("uv not installed; install with: curl -LsSf https://astral.sh/uv/install.sh | sh")
	}

	args := []string{"init"}
	if name != "" {
		args = append(args, name)
	}
	return h.run("uv", args...)
}

// Sync runs uv sync.
func (h *Helper) Sync(args ...string) error {
	if !h.HasUV() {
		return fmt.Errorf("uv not installed")
	}
	return h.run("uv", append([]string{"sync"}, args...)...)
}

// Add adds packages using uv.
func (h *Helper) Add(packages ...string) error {
	if !h.HasUV() {
		return fmt.Errorf("uv not installed")
	}
	if len(packages) == 0 {
		return fmt.Errorf("no packages specified")
	}
	return h.run("uv", append([]string{"add"}, packages...)...)
}

// Remove removes packages using uv.
func (h *Helper) Remove(packages ...string) error {
	if !h.HasUV() {
		return fmt.Errorf("uv not installed")
	}
	if len(packages) == 0 {
		return fmt.Errorf("no packages specified")
	}
	return h.run("uv", append([]string{"remove"}, packages...)...)
}

// Run executes a command using uv run.
func (h *Helper) Run(args ...string) error {
	if !h.HasUV() {
		return fmt.Errorf("uv not installed")
	}
	if len(args) == 0 {
		return fmt.Errorf("no command specified")
	}
	return h.run("uv", append([]string{"run"}, args...)...)
}

// InstallPackages installs packages using uv pip or pip.
func (h *Helper) InstallPackages(packages []string) error {
	if len(packages) == 0 {
		return nil
	}

	if h.HasUV() {
		return h.run("uv", append([]string{"pip", "install"}, packages...)...)
	}
	return h.run("pip", append([]string{"install"}, packages...)...)
}

// SetupFastAPI sets up a FastAPI development environment.
func (h *Helper) SetupFastAPI(venvName string) (*VenvInfo, error) {
	// Create venv first
	info, err := h.CreateVenv(venvName)
	if err != nil {
		return nil, err
	}

	// Install FastAPI deps
	fmt.Println("Installing FastAPI development dependencies...")
	if err := h.InstallPackages(FastAPIDeps); err != nil {
		return info, fmt.Errorf("failed to install FastAPI deps: %w", err)
	}

	return info, nil
}

// SetupIPython installs IPython with rich.
func (h *Helper) SetupIPython() error {
	fmt.Println("Installing IPython with rich output...")
	return h.InstallPackages(IPythonDeps)
}

// SetupDevTools installs common development tools.
func (h *Helper) SetupDevTools() error {
	fmt.Println("Installing development tools...")
	return h.InstallPackages(DevToolsDeps)
}

// GetEnvInfo returns Python environment information.
func (h *Helper) GetEnvInfo() *EnvInfo {
	info := &EnvInfo{}

	// Python version
	if out, err := exec.Command("python3", "--version").Output(); err == nil {
		info.Python = "python3"
		info.Version = strings.TrimSpace(string(out))
	}

	// Pip version
	if out, err := exec.Command("pip3", "--version").Output(); err == nil {
		parts := strings.Fields(string(out))
		if len(parts) >= 2 {
			info.Pip = parts[1]
		}
	}

	// UV version
	if out, err := exec.Command("uv", "--version").Output(); err == nil {
		info.UV = strings.TrimSpace(string(out))
	}

	// Active virtual environment
	info.VirtualEnv = os.Getenv("VIRTUAL_ENV")

	// Envs location
	info.EnvsLocation = os.Getenv("ENVS_LOCATION")
	if info.EnvsLocation == "" {
		info.EnvsLocation = filepath.Join(os.Getenv("HOME"), ".virtualenvs")
	}

	return info
}

// ListVenvs lists virtual environments in the default location.
func (h *Helper) ListVenvs() ([]VenvInfo, error) {
	envsLocation := os.Getenv("ENVS_LOCATION")
	if envsLocation == "" {
		envsLocation = filepath.Join(os.Getenv("HOME"), ".virtualenvs")
	}

	var venvs []VenvInfo

	// Check if envs location exists
	if _, err := os.Stat(envsLocation); os.IsNotExist(err) {
		return venvs, nil
	}

	entries, err := os.ReadDir(envsLocation)
	if err != nil {
		return nil, err
	}

	activeVenv := os.Getenv("VIRTUAL_ENV")

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		venvPath := filepath.Join(envsLocation, entry.Name())
		activatePath := filepath.Join(venvPath, "bin", "activate")

		// Check if it's actually a venv
		if _, err := os.Stat(activatePath); os.IsNotExist(err) {
			continue
		}

		info := VenvInfo{
			Name:   entry.Name(),
			Path:   venvPath,
			Active: venvPath == activeVenv,
		}

		// Get Python version
		pythonPath := filepath.Join(venvPath, "bin", "python")
		if _, err := os.Stat(pythonPath); err == nil {
			cmd := exec.Command(pythonPath, "--version")
			if out, err := cmd.Output(); err == nil {
				info.Python = strings.TrimSpace(string(out))
			}
		}

		venvs = append(venvs, info)
	}

	// Also check current directory for .venv
	if _, err := os.Stat(".venv/bin/activate"); err == nil {
		absPath, _ := filepath.Abs(".venv")
		info := VenvInfo{
			Name:   ".venv",
			Path:   absPath,
			Active: absPath == activeVenv,
		}
		pythonPath := filepath.Join(absPath, "bin", "python")
		if _, err := os.Stat(pythonPath); err == nil {
			cmd := exec.Command(pythonPath, "--version")
			if out, err := cmd.Output(); err == nil {
				info.Python = strings.TrimSpace(string(out))
			}
		}
		venvs = append([]VenvInfo{info}, venvs...)
	}

	return venvs, nil
}

// run executes a command.
func (h *Helper) run(name string, args ...string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would run: %s %s\n", name, strings.Join(args, " "))
		return nil
	}

	if h.verbose {
		fmt.Printf("Running: %s %s\n", name, strings.Join(args, " "))
	}

	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
