// Package pulumi provides Pulumi Infrastructure as Code helper functionality.
package pulumi

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Status represents Pulumi installation status.
type Status struct {
	Installed    bool   `json:"installed" yaml:"installed"`
	Version      string `json:"version,omitempty" yaml:"version,omitempty"`
	LoggedIn     bool   `json:"logged_in" yaml:"logged_in"`
	Backend      string `json:"backend,omitempty" yaml:"backend,omitempty"`
	Organization string `json:"organization,omitempty" yaml:"organization,omitempty"`
	User         string `json:"user,omitempty" yaml:"user,omitempty"`
}

// Stack represents a Pulumi stack.
type Stack struct {
	Name        string `json:"name" yaml:"name"`
	Current     bool   `json:"current" yaml:"current"`
	UpdateTime  string `json:"updateTime,omitempty" yaml:"update_time,omitempty"`
	ResourceCount int  `json:"resourceCount,omitempty" yaml:"resource_count,omitempty"`
}

// Project represents a Pulumi project.
type Project struct {
	Name        string `json:"name" yaml:"name"`
	Runtime     string `json:"runtime" yaml:"runtime"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Path        string `json:"path,omitempty" yaml:"path,omitempty"`
}

// Helper provides Pulumi helper operations.
type Helper struct {
	verbose bool
	dryRun  bool
}

// NewHelper creates a new Pulumi Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	return &Helper{
		verbose: verbose,
		dryRun:  dryRun,
	}
}

// GetStatus returns Pulumi status information.
func (h *Helper) GetStatus() *Status {
	status := &Status{}

	// Check if pulumi is installed
	out, err := exec.Command("pulumi", "version").Output()
	if err != nil {
		status.Installed = false
		return status
	}

	status.Installed = true
	status.Version = strings.TrimPrefix(strings.TrimSpace(string(out)), "v")

	// Check login status
	whoami, err := exec.Command("pulumi", "whoami", "-v").Output()
	if err == nil {
		status.LoggedIn = true
		// Parse whoami output
		lines := strings.Split(string(whoami), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "User:") {
				status.User = strings.TrimSpace(strings.TrimPrefix(line, "User:"))
			} else if strings.HasPrefix(line, "Organizations:") {
				status.Organization = strings.TrimSpace(strings.TrimPrefix(line, "Organizations:"))
			} else if strings.HasPrefix(line, "Backend URL:") {
				status.Backend = strings.TrimSpace(strings.TrimPrefix(line, "Backend URL:"))
			}
		}
		// Simple user extraction if verbose output doesn't work
		if status.User == "" {
			simpleWhoami, _ := exec.Command("pulumi", "whoami").Output()
			status.User = strings.TrimSpace(string(simpleWhoami))
		}
	}

	return status
}

// Login logs in to Pulumi.
func (h *Helper) Login(backend string) error {
	args := []string{"login"}
	if backend != "" {
		args = append(args, backend)
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: pulumi %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("pulumi", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Logout logs out from Pulumi.
func (h *Helper) Logout() error {
	if h.dryRun {
		fmt.Println("[dry-run] would run: pulumi logout")
		return nil
	}

	cmd := exec.Command("pulumi", "logout")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ListStacks lists all stacks in the current project.
func (h *Helper) ListStacks() ([]Stack, error) {
	cmd := exec.Command("pulumi", "stack", "ls", "--json")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list stacks: %w", err)
	}

	var stacks []Stack
	if err := json.Unmarshal(out, &stacks); err != nil {
		return nil, fmt.Errorf("failed to parse stacks: %w", err)
	}

	return stacks, nil
}

// SelectStack selects a stack.
func (h *Helper) SelectStack(name string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would select stack: %s\n", name)
		return nil
	}

	cmd := exec.Command("pulumi", "stack", "select", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Preview runs pulumi preview.
func (h *Helper) Preview() error {
	cmd := exec.Command("pulumi", "preview")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Up runs pulumi up.
func (h *Helper) Up(autoApprove bool) error {
	args := []string{"up"}
	if autoApprove {
		args = append(args, "--yes")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: pulumi %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("pulumi", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Destroy runs pulumi destroy.
func (h *Helper) Destroy(autoApprove bool) error {
	args := []string{"destroy"}
	if autoApprove {
		args = append(args, "--yes")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: pulumi %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("pulumi", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Refresh runs pulumi refresh.
func (h *Helper) Refresh() error {
	if h.dryRun {
		fmt.Println("[dry-run] would run: pulumi refresh")
		return nil
	}

	cmd := exec.Command("pulumi", "refresh", "--yes")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetOutputs gets stack outputs.
func (h *Helper) GetOutputs() (map[string]interface{}, error) {
	cmd := exec.Command("pulumi", "stack", "output", "--json")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get outputs: %w", err)
	}

	var outputs map[string]interface{}
	if err := json.Unmarshal(out, &outputs); err != nil {
		return nil, fmt.Errorf("failed to parse outputs: %w", err)
	}

	return outputs, nil
}

// NewProject creates a new Pulumi project.
func (h *Helper) NewProject(template, name, description string) error {
	args := []string{"new", template}
	if name != "" {
		args = append(args, "--name", name)
	}
	if description != "" {
		args = append(args, "--description", description)
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: pulumi %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("pulumi", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// GetProject returns current project info.
func (h *Helper) GetProject() (*Project, error) {
	// Check for Pulumi.yaml
	cwd, _ := os.Getwd()
	pulumiFile := filepath.Join(cwd, "Pulumi.yaml")

	if _, err := os.Stat(pulumiFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("not in a Pulumi project directory")
	}

	// Read project info
	data, err := os.ReadFile(pulumiFile)
	if err != nil {
		return nil, err
	}

	// Simple YAML parsing for name and runtime
	project := &Project{Path: cwd}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "name:") {
			project.Name = strings.TrimSpace(strings.TrimPrefix(line, "name:"))
		} else if strings.HasPrefix(line, "runtime:") {
			project.Runtime = strings.TrimSpace(strings.TrimPrefix(line, "runtime:"))
		} else if strings.HasPrefix(line, "description:") {
			project.Description = strings.TrimSpace(strings.TrimPrefix(line, "description:"))
		}
	}

	return project, nil
}

// Install installs Pulumi.
func (h *Helper) Install() error {
	if _, err := exec.LookPath("pulumi"); err == nil {
		return fmt.Errorf("Pulumi is already installed")
	}

	if h.dryRun {
		fmt.Println("[dry-run] would install Pulumi via Homebrew")
		return nil
	}

	fmt.Println("Installing Pulumi...")
	cmd := exec.Command("brew", "install", "pulumi")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
