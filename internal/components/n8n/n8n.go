// Package n8n provides n8n workflow automation helper functionality.
package n8n

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Status represents n8n installation status.
type Status struct {
	Installed       bool   `json:"installed" yaml:"installed"`
	Version         string `json:"version,omitempty" yaml:"version,omitempty"`
	Method          string `json:"method,omitempty" yaml:"method,omitempty"` // npm, docker, npx
	DockerRunning   bool   `json:"docker_running" yaml:"docker_running"`
	ContainerID     string `json:"container_id,omitempty" yaml:"container_id,omitempty"`
	WorkflowsDir    string `json:"workflows_dir,omitempty" yaml:"workflows_dir,omitempty"`
	CredentialsDir  string `json:"credentials_dir,omitempty" yaml:"credentials_dir,omitempty"`
}

// Workflow represents an n8n workflow.
type Workflow struct {
	ID     string `json:"id" yaml:"id"`
	Name   string `json:"name" yaml:"name"`
	Active bool   `json:"active" yaml:"active"`
	Path   string `json:"path,omitempty" yaml:"path,omitempty"`
}

// Helper provides n8n helper operations.
type Helper struct {
	verbose bool
	dryRun  bool
}

// NewHelper creates a new n8n Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	return &Helper{
		verbose: verbose,
		dryRun:  dryRun,
	}
}

// GetStatus returns n8n status information.
func (h *Helper) GetStatus() *Status {
	status := &Status{}

	// Check npm global installation
	if version := h.getNpmVersion(); version != "" {
		status.Installed = true
		status.Version = version
		status.Method = "npm"
	}

	// Check npx availability
	if !status.Installed {
		if h.canRunNpx() {
			status.Installed = true
			status.Method = "npx"
		}
	}

	// Check Docker
	if containerID := h.getDockerContainer(); containerID != "" {
		status.DockerRunning = true
		status.ContainerID = containerID
		if !status.Installed {
			status.Installed = true
			status.Method = "docker"
		}
	}

	// Get directories
	status.WorkflowsDir = h.getWorkflowsDir()
	status.CredentialsDir = h.getCredentialsDir()

	return status
}

// getNpmVersion returns n8n version if installed via npm.
func (h *Helper) getNpmVersion() string {
	cmd := exec.Command("n8n", "--version")
	out, err := cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(out))
	}
	return ""
}

// canRunNpx checks if npx can run n8n.
func (h *Helper) canRunNpx() bool {
	cmd := exec.Command("npx", "--yes", "n8n", "--version")
	cmd.Env = append(os.Environ(), "npm_config_yes=true")
	_, err := cmd.Output()
	return err == nil
}

// getDockerContainer returns running n8n container ID.
func (h *Helper) getDockerContainer() string {
	cmd := exec.Command("docker", "ps", "-q", "--filter", "ancestor=n8nio/n8n")
	out, err := cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(out))
	}

	// Try by name
	cmd = exec.Command("docker", "ps", "-q", "--filter", "name=n8n")
	out, err = cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(out))
	}

	return ""
}

// getWorkflowsDir returns the workflows export directory.
func (h *Helper) getWorkflowsDir() string {
	home, _ := os.UserHomeDir()

	xdgConfig := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfig == "" {
		xdgConfig = filepath.Join(home, ".config")
	}

	return filepath.Join(xdgConfig, "n8n", "workflows")
}

// getCredentialsDir returns the credentials export directory.
func (h *Helper) getCredentialsDir() string {
	home, _ := os.UserHomeDir()

	xdgConfig := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfig == "" {
		xdgConfig = filepath.Join(home, ".config")
	}

	return filepath.Join(xdgConfig, "n8n", "credentials")
}

// Start starts n8n locally.
func (h *Helper) Start(useDocker bool, port int) error {
	if port == 0 {
		port = 5678
	}

	if h.dryRun {
		if useDocker {
			fmt.Printf("[dry-run] would start n8n via Docker on port %d\n", port)
		} else {
			fmt.Printf("[dry-run] would start n8n on port %d\n", port)
		}
		return nil
	}

	if useDocker {
		return h.startDocker(port)
	}

	return h.startLocal(port)
}

// startLocal starts n8n using npm/npx.
func (h *Helper) startLocal(port int) error {
	var cmd *exec.Cmd

	// Try n8n directly first
	if _, err := exec.LookPath("n8n"); err == nil {
		cmd = exec.Command("n8n", "start", "--port", fmt.Sprintf("%d", port))
	} else {
		// Fall back to npx
		cmd = exec.Command("npx", "n8n", "start", "--port", fmt.Sprintf("%d", port))
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Printf("Starting n8n on http://localhost:%d ...\n", port)
	return cmd.Run()
}

// startDocker starts n8n using Docker.
func (h *Helper) startDocker(port int) error {
	home, _ := os.UserHomeDir()
	dataDir := filepath.Join(home, ".n8n")

	// Create data directory if needed
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	args := []string{
		"run", "-d",
		"--name", "n8n",
		"-p", fmt.Sprintf("%d:5678", port),
		"-v", fmt.Sprintf("%s:/home/node/.n8n", dataDir),
		"n8nio/n8n",
	}

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Starting n8n container on http://localhost:%d ...\n", port)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start n8n container: %w", err)
	}

	fmt.Println("n8n container started. View logs with: docker logs -f n8n")
	return nil
}

// Stop stops n8n.
func (h *Helper) Stop() error {
	if h.dryRun {
		fmt.Println("[dry-run] would stop n8n")
		return nil
	}

	// Check for Docker container
	containerID := h.getDockerContainer()
	if containerID != "" {
		cmd := exec.Command("docker", "stop", containerID)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to stop container: %w", err)
		}

		// Remove container
		exec.Command("docker", "rm", containerID).Run()
		return nil
	}

	// Try to find and kill n8n process
	cmd := exec.Command("pkill", "-f", "n8n")
	cmd.Run() // Ignore errors

	return nil
}

// ExportWorkflows exports all workflows to files.
func (h *Helper) ExportWorkflows(outputDir string) error {
	if outputDir == "" {
		outputDir = h.getWorkflowsDir()
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would export workflows to %s\n", outputDir)
		return nil
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Use n8n export command
	var cmd *exec.Cmd
	if _, err := exec.LookPath("n8n"); err == nil {
		cmd = exec.Command("n8n", "export:workflow", "--all", "--output", outputDir)
	} else {
		cmd = exec.Command("npx", "n8n", "export:workflow", "--all", "--output", outputDir)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ImportWorkflow imports a workflow from a file.
func (h *Helper) ImportWorkflow(inputFile string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would import workflow from %s\n", inputFile)
		return nil
	}

	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("workflow file not found: %s", inputFile)
	}

	var cmd *exec.Cmd
	if _, err := exec.LookPath("n8n"); err == nil {
		cmd = exec.Command("n8n", "import:workflow", "--input", inputFile)
	} else {
		cmd = exec.Command("npx", "n8n", "import:workflow", "--input", inputFile)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ExportCredentials exports credentials (encrypted).
func (h *Helper) ExportCredentials(outputDir string) error {
	if outputDir == "" {
		outputDir = h.getCredentialsDir()
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would export credentials to %s\n", outputDir)
		return nil
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	var cmd *exec.Cmd
	if _, err := exec.LookPath("n8n"); err == nil {
		cmd = exec.Command("n8n", "export:credentials", "--all", "--output", outputDir)
	} else {
		cmd = exec.Command("npx", "n8n", "export:credentials", "--all", "--output", outputDir)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ListWorkflows lists exported workflows.
func (h *Helper) ListWorkflows() ([]Workflow, error) {
	workflowsDir := h.getWorkflowsDir()
	if _, err := os.Stat(workflowsDir); os.IsNotExist(err) {
		return []Workflow{}, nil
	}

	var workflows []Workflow

	entries, err := os.ReadDir(workflowsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read workflows directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".json") {
			continue
		}

		path := filepath.Join(workflowsDir, name)
		workflow := Workflow{
			Name: strings.TrimSuffix(name, ".json"),
			Path: path,
		}

		// Try to parse workflow info
		if info := h.parseWorkflowInfo(path); info != nil {
			workflow.ID = info.ID
			workflow.Name = info.Name
			workflow.Active = info.Active
		}

		workflows = append(workflows, workflow)
	}

	return workflows, nil
}

// workflowInfo holds parsed workflow metadata.
type workflowInfo struct {
	ID     string
	Name   string
	Active bool
}

// parseWorkflowInfo parses workflow metadata from a JSON file.
func (h *Helper) parseWorkflowInfo(path string) *workflowInfo {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var parsed struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Active bool   `json:"active"`
	}

	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil
	}

	return &workflowInfo{
		ID:     parsed.ID,
		Name:   parsed.Name,
		Active: parsed.Active,
	}
}

// Install installs n8n globally via npm.
func (h *Helper) Install() error {
	if h.dryRun {
		fmt.Println("[dry-run] would install n8n via npm")
		return nil
	}

	fmt.Println("Installing n8n globally...")
	cmd := exec.Command("npm", "install", "-g", "n8n")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
