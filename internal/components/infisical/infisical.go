// Package infisical provides Infisical secret management CLI helper functionality.
package infisical

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Status represents Infisical CLI installation status.
type Status struct {
	Installed     bool   `json:"installed" yaml:"installed"`
	Version       string `json:"version,omitempty" yaml:"version,omitempty"`
	LoggedIn      bool   `json:"logged_in" yaml:"logged_in"`
	ConfigDir     string `json:"config_dir" yaml:"config_dir"`
	InProject     bool   `json:"in_project" yaml:"in_project"`
	ProjectID     string `json:"project_id,omitempty" yaml:"project_id,omitempty"`
	Environment   string `json:"environment,omitempty" yaml:"environment,omitempty"`
}

// Secret represents a secret from Infisical.
type Secret struct {
	Key   string `json:"secretKey" yaml:"key"`
	Value string `json:"secretValue,omitempty" yaml:"value,omitempty"`
	Type  string `json:"type,omitempty" yaml:"type,omitempty"`
}

// Project represents an Infisical project.
type Project struct {
	ID          string `json:"id" yaml:"id"`
	Name        string `json:"name" yaml:"name"`
	Slug        string `json:"slug,omitempty" yaml:"slug,omitempty"`
}

// Helper provides Infisical helper operations.
type Helper struct {
	verbose bool
	dryRun  bool
}

// NewHelper creates a new Infisical Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	return &Helper{
		verbose: verbose,
		dryRun:  dryRun,
	}
}

// GetConfigDir returns the Infisical config directory.
func (h *Helper) GetConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".infisical")
}

// GetStatus returns Infisical CLI status information.
func (h *Helper) GetStatus() *Status {
	status := &Status{
		ConfigDir: h.GetConfigDir(),
	}

	// Check if infisical is installed
	out, err := exec.Command("infisical", "--version").Output()
	if err != nil {
		status.Installed = false
		return status
	}

	status.Installed = true
	status.Version = strings.TrimSpace(string(out))

	// Check if logged in by trying to get user info
	if h.isLoggedIn() {
		status.LoggedIn = true
	}

	// Check if in a project (has .infisical.json)
	if _, err := os.Stat(".infisical.json"); err == nil {
		status.InProject = true
		status.ProjectID, status.Environment = h.getProjectInfo()
	}

	return status
}

// isLoggedIn checks if user is logged in.
func (h *Helper) isLoggedIn() bool {
	cmd := exec.Command("infisical", "user")
	err := cmd.Run()
	return err == nil
}

// getProjectInfo reads project info from .infisical.json.
func (h *Helper) getProjectInfo() (string, string) {
	data, err := os.ReadFile(".infisical.json")
	if err != nil {
		return "", ""
	}

	var config struct {
		WorkspaceID string `json:"workspaceId"`
		Environment string `json:"defaultEnvironment"`
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return "", ""
	}

	return config.WorkspaceID, config.Environment
}

// Login authenticates with Infisical.
func (h *Helper) Login() error {
	if h.dryRun {
		fmt.Println("[dry-run] would run: infisical login")
		return nil
	}

	cmd := exec.Command("infisical", "login")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Logout logs out from Infisical.
func (h *Helper) Logout() error {
	if h.dryRun {
		fmt.Println("[dry-run] would run: infisical logout")
		return nil
	}

	cmd := exec.Command("infisical", "logout")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Init initializes an Infisical project in the current directory.
func (h *Helper) Init() error {
	if h.dryRun {
		fmt.Println("[dry-run] would run: infisical init")
		return nil
	}

	cmd := exec.Command("infisical", "init")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// ListSecrets lists secrets for the current project/environment.
func (h *Helper) ListSecrets(env string) ([]Secret, error) {
	args := []string{"secrets", "--plain"}
	if env != "" {
		args = append(args, "--env", env)
	}

	cmd := exec.Command("infisical", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}

	// Parse output - infisical secrets --plain outputs KEY=VALUE format
	var secrets []Secret
	for line := range strings.SplitSeq(strings.TrimSpace(string(out)), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			secrets = append(secrets, Secret{
				Key:   parts[0],
				Value: parts[1],
			})
		}
	}

	return secrets, nil
}

// GetSecret retrieves a specific secret.
func (h *Helper) GetSecret(key, env string) (string, error) {
	args := []string{"secrets", "get", key, "--plain"}
	if env != "" {
		args = append(args, "--env", env)
	}

	cmd := exec.Command("infisical", args...)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get secret %s: %w", key, err)
	}

	return strings.TrimSpace(string(out)), nil
}

// Run executes a command with secrets injected.
func (h *Helper) Run(env string, command []string) error {
	if len(command) == 0 {
		return fmt.Errorf("no command specified")
	}

	args := []string{"run"}
	if env != "" {
		args = append(args, "--env", env)
	}
	args = append(args, "--")
	args = append(args, command...)

	if h.dryRun {
		fmt.Printf("[dry-run] would run: infisical %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("infisical", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Export exports secrets to various formats.
func (h *Helper) Export(env, format string) (string, error) {
	args := []string{"export"}
	if env != "" {
		args = append(args, "--env", env)
	}
	if format != "" {
		args = append(args, "--format", format)
	}

	cmd := exec.Command("infisical", args...)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to export secrets: %w", err)
	}

	return string(out), nil
}

// Install installs Infisical CLI.
func (h *Helper) Install() error {
	if _, err := exec.LookPath("infisical"); err == nil {
		return fmt.Errorf("infisical is already installed")
	}

	if h.dryRun {
		fmt.Println("[dry-run] would install infisical via homebrew")
		return nil
	}

	fmt.Println("Installing Infisical CLI...")
	cmd := exec.Command("brew", "install", "infisical/get-cli/infisical")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Scan scans for secrets in files.
func (h *Helper) Scan(path string) error {
	args := []string{"scan"}
	if path != "" {
		args = append(args, path)
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: infisical %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("infisical", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
