// Package terraform provides Terraform Infrastructure as Code helper functionality.
package terraform

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Status represents Terraform installation status.
type Status struct {
	Installed    bool   `json:"installed" yaml:"installed"`
	Version      string `json:"version,omitempty" yaml:"version,omitempty"`
	TerragruntInstalled bool `json:"terragrunt_installed" yaml:"terragrunt_installed"`
	TerragruntVersion   string `json:"terragrunt_version,omitempty" yaml:"terragrunt_version,omitempty"`
	Initialized  bool   `json:"initialized" yaml:"initialized"`
	Backend      string `json:"backend,omitempty" yaml:"backend,omitempty"`
	Workspace    string `json:"workspace,omitempty" yaml:"workspace,omitempty"`
}

// Workspace represents a Terraform workspace.
type Workspace struct {
	Name    string `json:"name" yaml:"name"`
	Current bool   `json:"current" yaml:"current"`
}

// Resource represents a Terraform resource.
type Resource struct {
	Type    string `json:"type" yaml:"type"`
	Name    string `json:"name" yaml:"name"`
	Module  string `json:"module,omitempty" yaml:"module,omitempty"`
	Address string `json:"address" yaml:"address"`
}

// Output represents a Terraform output.
type Output struct {
	Name      string      `json:"name" yaml:"name"`
	Value     interface{} `json:"value" yaml:"value"`
	Sensitive bool        `json:"sensitive" yaml:"sensitive"`
}

// Helper provides Terraform helper operations.
type Helper struct {
	verbose bool
	dryRun  bool
}

// NewHelper creates a new Terraform Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	return &Helper{
		verbose: verbose,
		dryRun:  dryRun,
	}
}

// GetStatus returns Terraform status information.
func (h *Helper) GetStatus() *Status {
	status := &Status{}

	// Check if terraform is installed
	out, err := exec.Command("terraform", "version", "-json").Output()
	if err != nil {
		// Try without -json
		out, err = exec.Command("terraform", "version").Output()
		if err != nil {
			status.Installed = false
			return status
		}
		status.Installed = true
		lines := strings.Split(string(out), "\n")
		if len(lines) > 0 {
			parts := strings.Fields(lines[0])
			if len(parts) >= 2 {
				status.Version = strings.TrimPrefix(parts[1], "v")
			}
		}
	} else {
		status.Installed = true
		var versionInfo struct {
			TerraformVersion string `json:"terraform_version"`
		}
		if json.Unmarshal(out, &versionInfo) == nil {
			status.Version = versionInfo.TerraformVersion
		}
	}

	// Check terragrunt
	out, err = exec.Command("terragrunt", "--version").Output()
	if err == nil {
		status.TerragruntInstalled = true
		parts := strings.Fields(string(out))
		if len(parts) >= 2 {
			status.TerragruntVersion = strings.TrimPrefix(parts[1], "v")
		}
	}

	// Check if initialized
	if _, err := os.Stat(".terraform"); err == nil {
		status.Initialized = true
	}

	// Get workspace
	out, _ = exec.Command("terraform", "workspace", "show").Output()
	if len(out) > 0 {
		status.Workspace = strings.TrimSpace(string(out))
	}

	return status
}

// Init initializes Terraform.
func (h *Helper) Init(upgrade bool) error {
	args := []string{"init"}
	if upgrade {
		args = append(args, "-upgrade")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: terraform %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("terraform", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Plan runs terraform plan.
func (h *Helper) Plan(outFile string) error {
	args := []string{"plan"}
	if outFile != "" {
		args = append(args, "-out", outFile)
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: terraform %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("terraform", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Apply runs terraform apply.
func (h *Helper) Apply(planFile string, autoApprove bool) error {
	args := []string{"apply"}
	if autoApprove {
		args = append(args, "-auto-approve")
	}
	if planFile != "" {
		args = append(args, planFile)
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: terraform %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("terraform", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Destroy runs terraform destroy.
func (h *Helper) Destroy(autoApprove bool) error {
	args := []string{"destroy"}
	if autoApprove {
		args = append(args, "-auto-approve")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: terraform %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("terraform", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Validate validates configuration.
func (h *Helper) Validate() error {
	cmd := exec.Command("terraform", "validate")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Format formats configuration.
func (h *Helper) Format(check bool) error {
	args := []string{"fmt"}
	if check {
		args = append(args, "-check")
	} else {
		args = append(args, "-recursive")
	}

	cmd := exec.Command("terraform", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ListWorkspaces lists all workspaces.
func (h *Helper) ListWorkspaces() ([]Workspace, error) {
	cmd := exec.Command("terraform", "workspace", "list")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list workspaces: %w", err)
	}

	var workspaces []Workspace
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		current := false
		if strings.HasPrefix(line, "* ") {
			current = true
			line = strings.TrimPrefix(line, "* ")
		}

		workspaces = append(workspaces, Workspace{
			Name:    line,
			Current: current,
		})
	}

	return workspaces, nil
}

// SelectWorkspace selects a workspace.
func (h *Helper) SelectWorkspace(name string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would select workspace: %s\n", name)
		return nil
	}

	cmd := exec.Command("terraform", "workspace", "select", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// NewWorkspace creates a new workspace.
func (h *Helper) NewWorkspace(name string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would create workspace: %s\n", name)
		return nil
	}

	cmd := exec.Command("terraform", "workspace", "new", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetOutputs gets all outputs.
func (h *Helper) GetOutputs() (map[string]Output, error) {
	cmd := exec.Command("terraform", "output", "-json")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get outputs: %w", err)
	}

	var rawOutputs map[string]struct {
		Value     interface{} `json:"value"`
		Sensitive bool        `json:"sensitive"`
	}
	if err := json.Unmarshal(out, &rawOutputs); err != nil {
		return nil, fmt.Errorf("failed to parse outputs: %w", err)
	}

	outputs := make(map[string]Output)
	for name, raw := range rawOutputs {
		outputs[name] = Output{
			Name:      name,
			Value:     raw.Value,
			Sensitive: raw.Sensitive,
		}
	}

	return outputs, nil
}

// GetState gets the current state.
func (h *Helper) GetState() ([]Resource, error) {
	cmd := exec.Command("terraform", "state", "list")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list state: %w", err)
	}

	var resources []Resource
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		resource := Resource{Address: line}

		// Parse address to get type and name
		parts := strings.Split(line, ".")
		if len(parts) >= 2 {
			// Check for module prefix
			if strings.HasPrefix(parts[0], "module") && len(parts) >= 4 {
				resource.Module = parts[0] + "." + parts[1]
				resource.Type = parts[2]
				resource.Name = strings.Join(parts[3:], ".")
			} else {
				resource.Type = parts[0]
				resource.Name = strings.Join(parts[1:], ".")
			}
		}

		resources = append(resources, resource)
	}

	return resources, nil
}

// Refresh refreshes state.
func (h *Helper) Refresh() error {
	if h.dryRun {
		fmt.Println("[dry-run] would run: terraform refresh")
		return nil
	}

	cmd := exec.Command("terraform", "refresh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Taint taints a resource.
func (h *Helper) Taint(address string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would taint: %s\n", address)
		return nil
	}

	cmd := exec.Command("terraform", "taint", address)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Untaint untaints a resource.
func (h *Helper) Untaint(address string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would untaint: %s\n", address)
		return nil
	}

	cmd := exec.Command("terraform", "untaint", address)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// IsInitialized checks if current directory is initialized.
func (h *Helper) IsInitialized() bool {
	_, err := os.Stat(".terraform")
	return err == nil
}

// HasTerraformFiles checks if current directory has Terraform files.
func (h *Helper) HasTerraformFiles() bool {
	matches, _ := filepath.Glob("*.tf")
	return len(matches) > 0
}

// Install installs Terraform.
func (h *Helper) Install() error {
	if _, err := exec.LookPath("terraform"); err == nil {
		return fmt.Errorf("Terraform is already installed")
	}

	if h.dryRun {
		fmt.Println("[dry-run] would install Terraform via Homebrew")
		return nil
	}

	fmt.Println("Installing Terraform...")
	cmd := exec.Command("brew", "install", "terraform")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
