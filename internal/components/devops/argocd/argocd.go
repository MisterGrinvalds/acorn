// Package argocd provides ArgoCD helper functionality.
package argocd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Status represents ArgoCD installation status.
type Status struct {
	Installed       bool   `json:"installed" yaml:"installed"`
	Version         string `json:"version,omitempty" yaml:"version,omitempty"`
	Location        string `json:"location,omitempty" yaml:"location,omitempty"`
	ConfigDir       string `json:"config_dir,omitempty" yaml:"config_dir,omitempty"`
	ServerConnected bool   `json:"server_connected" yaml:"server_connected"`
	CurrentContext  string `json:"current_context,omitempty" yaml:"current_context,omitempty"`
	Server          string `json:"server,omitempty" yaml:"server,omitempty"`
}

// App represents an ArgoCD application.
type App struct {
	Name        string `json:"name" yaml:"name"`
	Namespace   string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Project     string `json:"project,omitempty" yaml:"project,omitempty"`
	Status      string `json:"status,omitempty" yaml:"status,omitempty"`
	Health      string `json:"health,omitempty" yaml:"health,omitempty"`
	SyncStatus  string `json:"sync_status,omitempty" yaml:"sync_status,omitempty"`
	Repo        string `json:"repo,omitempty" yaml:"repo,omitempty"`
	Path        string `json:"path,omitempty" yaml:"path,omitempty"`
	Destination string `json:"destination,omitempty" yaml:"destination,omitempty"`
}

// Context represents an ArgoCD context.
type Context struct {
	Name    string `json:"name" yaml:"name"`
	Server  string `json:"server,omitempty" yaml:"server,omitempty"`
	Current bool   `json:"current" yaml:"current"`
}

// Cluster represents an ArgoCD cluster.
type Cluster struct {
	Server     string `json:"server" yaml:"server"`
	Name       string `json:"name,omitempty" yaml:"name,omitempty"`
	Namespaces string `json:"namespaces,omitempty" yaml:"namespaces,omitempty"`
}

// Repo represents an ArgoCD repository.
type Repo struct {
	URL    string `json:"url" yaml:"url"`
	Type   string `json:"type,omitempty" yaml:"type,omitempty"`
	Name   string `json:"name,omitempty" yaml:"name,omitempty"`
	Status string `json:"status,omitempty" yaml:"status,omitempty"`
}

// Helper provides ArgoCD helper operations.
type Helper struct {
	verbose bool
}

// NewHelper creates a new ArgoCD Helper.
func NewHelper(verbose bool) *Helper {
	return &Helper{
		verbose: verbose,
	}
}

// IsInstalled checks if argocd CLI is installed.
func (h *Helper) IsInstalled() bool {
	_, err := exec.LookPath("argocd")
	return err == nil
}

// IsServerConnected checks if connected to an ArgoCD server.
func (h *Helper) IsServerConnected() bool {
	cmd := exec.Command("argocd", "app", "list", "--output", "name")
	return cmd.Run() == nil
}

// GetStatus returns ArgoCD installation status.
func (h *Helper) GetStatus() *Status {
	status := &Status{}

	argoPath, err := exec.LookPath("argocd")
	if err != nil {
		return status
	}

	status.Installed = true
	status.Location = argoPath

	// Get version
	cmd := exec.Command("argocd", "version", "--client")
	out, err := cmd.Output()
	if err == nil {
		// Format: argocd: v2.9.0+...
		for line := range strings.SplitSeq(string(out), "\n") {
			if version, found := strings.CutPrefix(line, "argocd:"); found {
				status.Version = strings.TrimSpace(version)
				break
			}
		}
	}

	// Get config directory
	status.ConfigDir = h.GetConfigDir()

	// Check server connection
	status.ServerConnected = h.IsServerConnected()

	// Get current context
	if ctx := h.GetCurrentContext(); ctx != nil {
		status.CurrentContext = ctx.Name
		status.Server = ctx.Server
	}

	return status
}

// GetConfigDir returns the ArgoCD config directory.
func (h *Helper) GetConfigDir() string {
	// Check environment variable first
	if dir := os.Getenv("ARGOCD_CONFIG_DIR"); dir != "" {
		return dir
	}

	// XDG config home
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, "argocd")
	}

	// Default
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "argocd")
}

// GetConfigFile returns the path to the config file.
func (h *Helper) GetConfigFile() string {
	return filepath.Join(h.GetConfigDir(), "config")
}

// GetCurrentContext returns the current ArgoCD context.
func (h *Helper) GetCurrentContext() *Context {
	cmd := exec.Command("argocd", "context")
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	// Parse context output - first line after header is current
	// Format: NAME SERVER
	var current *Context
	lineNum := 0
	for line := range strings.SplitSeq(string(out), "\n") {
		lineNum++
		if lineNum == 1 || strings.TrimSpace(line) == "" {
			continue // Skip header
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			current = &Context{
				Name:    fields[0],
				Server:  fields[1],
				Current: true,
			}
			break
		}
	}
	return current
}

// ListContexts returns all ArgoCD contexts.
func (h *Helper) ListContexts() ([]Context, error) {
	cmd := exec.Command("argocd", "context")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list contexts: %w", err)
	}

	var contexts []Context
	lineNum := 0
	for line := range strings.SplitSeq(string(out), "\n") {
		lineNum++
		if lineNum == 1 || strings.TrimSpace(line) == "" {
			continue // Skip header
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			ctx := Context{
				Name:    fields[0],
				Server:  fields[1],
				Current: lineNum == 2, // First entry after header is current
			}
			contexts = append(contexts, ctx)
		}
	}
	return contexts, nil
}

// SwitchContext switches to a different ArgoCD context.
func (h *Helper) SwitchContext(name string) error {
	cmd := exec.Command("argocd", "context", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Login logs into an ArgoCD server.
func (h *Helper) Login(server string, sso bool, username, password string) error {
	args := []string{"login", server}
	if sso {
		args = append(args, "--sso")
	} else if username != "" && password != "" {
		args = append(args, "--username", username, "--password", password)
	}

	cmd := exec.Command("argocd", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ListApps returns all ArgoCD applications.
func (h *Helper) ListApps(project, selector string) ([]App, error) {
	args := []string{"app", "list", "--output", "json"}
	if project != "" {
		args = append(args, "--project", project)
	}
	if selector != "" {
		args = append(args, "--selector", selector)
	}

	cmd := exec.Command("argocd", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list apps: %w", err)
	}

	var apps []App
	if err := json.Unmarshal(out, &apps); err != nil {
		// Try parsing as newline-separated JSON objects
		for line := range strings.SplitSeq(string(out), "\n") {
			if strings.TrimSpace(line) == "" {
				continue
			}
			var app App
			if err := json.Unmarshal([]byte(line), &app); err == nil {
				apps = append(apps, app)
			}
		}
	}
	return apps, nil
}

// GetApp returns details for a specific application.
func (h *Helper) GetApp(name string) (*App, error) {
	cmd := exec.Command("argocd", "app", "get", name, "--output", "json")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get app %s: %w", name, err)
	}

	var app App
	if err := json.Unmarshal(out, &app); err != nil {
		return nil, fmt.Errorf("failed to parse app: %w", err)
	}
	return &app, nil
}

// SyncApp syncs an application.
func (h *Helper) SyncApp(name string, prune, force bool) error {
	args := []string{"app", "sync", name}
	if prune {
		args = append(args, "--prune")
	}
	if force {
		args = append(args, "--force")
	}

	cmd := exec.Command("argocd", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// DiffApp shows the diff for an application.
func (h *Helper) DiffApp(name string) error {
	cmd := exec.Command("argocd", "app", "diff", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RollbackApp rolls back an application to a previous revision.
func (h *Helper) RollbackApp(name string, revision int) error {
	cmd := exec.Command("argocd", "app", "rollback", name, fmt.Sprintf("%d", revision))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetAppHistory returns deployment history for an application.
func (h *Helper) GetAppHistory(name string) error {
	cmd := exec.Command("argocd", "app", "history", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// WaitForApp waits for an application to sync.
func (h *Helper) WaitForApp(name string, timeout int) error {
	args := []string{"app", "wait", name}
	if timeout > 0 {
		args = append(args, "--timeout", fmt.Sprintf("%d", timeout))
	}

	cmd := exec.Command("argocd", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ListClusters returns all registered clusters.
func (h *Helper) ListClusters() ([]Cluster, error) {
	cmd := exec.Command("argocd", "cluster", "list", "--output", "json")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list clusters: %w", err)
	}

	var clusters []Cluster
	if err := json.Unmarshal(out, &clusters); err != nil {
		// Try parsing line by line
		for line := range strings.SplitSeq(string(out), "\n") {
			if strings.TrimSpace(line) == "" {
				continue
			}
			var cluster Cluster
			if err := json.Unmarshal([]byte(line), &cluster); err == nil {
				clusters = append(clusters, cluster)
			}
		}
	}
	return clusters, nil
}

// ListRepos returns all registered repositories.
func (h *Helper) ListRepos() ([]Repo, error) {
	cmd := exec.Command("argocd", "repo", "list", "--output", "json")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list repos: %w", err)
	}

	var repos []Repo
	if err := json.Unmarshal(out, &repos); err != nil {
		// Try parsing line by line
		for line := range strings.SplitSeq(string(out), "\n") {
			if strings.TrimSpace(line) == "" {
				continue
			}
			var repo Repo
			if err := json.Unmarshal([]byte(line), &repo); err == nil {
				repos = append(repos, repo)
			}
		}
	}
	return repos, nil
}

// ListProjects returns all ArgoCD projects.
func (h *Helper) ListProjects() ([]string, error) {
	cmd := exec.Command("argocd", "proj", "list", "--output", "name")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}

	var projects []string
	for line := range strings.SplitSeq(string(out), "\n") {
		if name := strings.TrimSpace(line); name != "" {
			projects = append(projects, name)
		}
	}
	return projects, nil
}

// GetAppLogs streams logs for an application.
func (h *Helper) GetAppLogs(name string, follow bool, container, group string) error {
	args := []string{"app", "logs", name}
	if follow {
		args = append(args, "--follow")
	}
	if container != "" {
		args = append(args, "--container", container)
	}
	if group != "" {
		args = append(args, "--group", group)
	}

	cmd := exec.Command("argocd", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetAppResources returns resources for an application.
func (h *Helper) GetAppResources(name string, tree bool) error {
	args := []string{"app", "resources", name}
	if tree {
		args = append(args, "--tree")
	}

	cmd := exec.Command("argocd", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
