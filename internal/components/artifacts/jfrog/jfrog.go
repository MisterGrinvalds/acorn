// Package jfrog provides JFrog Artifactory CLI helper functionality.
package jfrog

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Status represents JFrog CLI installation status.
type Status struct {
	Installed     bool     `json:"installed" yaml:"installed"`
	Version       string   `json:"version,omitempty" yaml:"version,omitempty"`
	Servers       []Server `json:"servers,omitempty" yaml:"servers,omitempty"`
	DefaultServer string   `json:"default_server,omitempty" yaml:"default_server,omitempty"`
}

// Server represents a configured JFrog server.
type Server struct {
	ID        string `json:"serverId" yaml:"id"`
	URL       string `json:"url" yaml:"url"`
	IsDefault bool   `json:"isDefault,omitempty" yaml:"is_default,omitempty"`
}

// Artifact represents an artifact in the repository.
type Artifact struct {
	Path     string `json:"path" yaml:"path"`
	Name     string `json:"name,omitempty" yaml:"name,omitempty"`
	Size     int64  `json:"size,omitempty" yaml:"size,omitempty"`
	Modified string `json:"modified,omitempty" yaml:"modified,omitempty"`
}

// Helper provides JFrog CLI helper operations.
type Helper struct {
	verbose bool
	dryRun  bool
}

// NewHelper creates a new JFrog Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	return &Helper{
		verbose: verbose,
		dryRun:  dryRun,
	}
}

// GetStatus returns JFrog CLI status information.
func (h *Helper) GetStatus() *Status {
	status := &Status{}

	// Check if jf is installed
	out, err := exec.Command("jf", "--version").Output()
	if err != nil {
		status.Installed = false
		return status
	}

	status.Installed = true
	// Parse version (format: jf version X.X.X)
	versionStr := strings.TrimSpace(string(out))
	parts := strings.Fields(versionStr)
	if len(parts) >= 3 {
		status.Version = parts[2]
	} else {
		status.Version = versionStr
	}

	// Get configured servers
	servers, defaultServer := h.getServers()
	status.Servers = servers
	status.DefaultServer = defaultServer

	return status
}

// getServers retrieves configured JFrog servers.
func (h *Helper) getServers() ([]Server, string) {
	cmd := exec.Command("jf", "config", "show")
	out, err := cmd.Output()
	if err != nil {
		return nil, ""
	}

	var servers []Server
	var defaultServer string

	// Parse JSON output
	var configs []struct {
		ServerID  string `json:"serverId"`
		URL       string `json:"url"`
		IsDefault bool   `json:"isDefault"`
	}

	if err := json.Unmarshal(out, &configs); err != nil {
		// Try parsing as single server
		var single struct {
			ServerID  string `json:"serverId"`
			URL       string `json:"url"`
			IsDefault bool   `json:"isDefault"`
		}
		if err := json.Unmarshal(out, &single); err == nil && single.ServerID != "" {
			servers = append(servers, Server{
				ID:        single.ServerID,
				URL:       single.URL,
				IsDefault: single.IsDefault,
			})
			if single.IsDefault {
				defaultServer = single.ServerID
			}
		}
		return servers, defaultServer
	}

	for _, c := range configs {
		servers = append(servers, Server{
			ID:        c.ServerID,
			URL:       c.URL,
			IsDefault: c.IsDefault,
		})
		if c.IsDefault {
			defaultServer = c.ServerID
		}
	}

	return servers, defaultServer
}

// Ping tests connectivity to a JFrog server.
func (h *Helper) Ping(serverID string) error {
	args := []string{"rt", "ping"}
	if serverID != "" {
		args = append(args, "--server-id", serverID)
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: jf %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("jf", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// AddServer adds a new JFrog server configuration.
func (h *Helper) AddServer(serverID, url, user, password string, interactive bool) error {
	args := []string{"config", "add", serverID}

	if url != "" {
		args = append(args, "--url", url)
	}
	if user != "" {
		args = append(args, "--user", user)
	}
	if password != "" {
		args = append(args, "--password", password)
	}
	if interactive {
		args = append(args, "--interactive")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: jf %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("jf", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// RemoveServer removes a JFrog server configuration.
func (h *Helper) RemoveServer(serverID string) error {
	if serverID == "" {
		return fmt.Errorf("server ID is required")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: jf config remove %s\n", serverID)
		return nil
	}

	cmd := exec.Command("jf", "config", "remove", serverID)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// UseServer sets the default JFrog server.
func (h *Helper) UseServer(serverID string) error {
	if serverID == "" {
		return fmt.Errorf("server ID is required")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: jf config use %s\n", serverID)
		return nil
	}

	cmd := exec.Command("jf", "config", "use", serverID)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Upload uploads artifacts to Artifactory.
func (h *Helper) Upload(source, target string, recursive, flat bool) error {
	if source == "" || target == "" {
		return fmt.Errorf("source and target are required")
	}

	args := []string{"rt", "upload", source, target}
	if recursive {
		args = append(args, "--recursive")
	}
	if flat {
		args = append(args, "--flat")
	}

	if h.dryRun {
		args = append(args, "--dry-run")
	}

	cmd := exec.Command("jf", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Download downloads artifacts from Artifactory.
func (h *Helper) Download(source, target string, recursive, flat bool) error {
	if source == "" {
		return fmt.Errorf("source is required")
	}

	args := []string{"rt", "download", source}
	if target != "" {
		args = append(args, target)
	}
	if recursive {
		args = append(args, "--recursive")
	}
	if flat {
		args = append(args, "--flat")
	}

	if h.dryRun {
		args = append(args, "--dry-run")
	}

	cmd := exec.Command("jf", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Search searches for artifacts in Artifactory.
func (h *Helper) Search(pattern string) ([]Artifact, error) {
	if pattern == "" {
		return nil, fmt.Errorf("search pattern is required")
	}

	cmd := exec.Command("jf", "rt", "search", pattern)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	var results []struct {
		Path     string `json:"path"`
		Size     int64  `json:"size"`
		Modified string `json:"modified"`
	}

	if err := json.Unmarshal(out, &results); err != nil {
		return nil, fmt.Errorf("failed to parse results: %w", err)
	}

	var artifacts []Artifact
	for _, r := range results {
		artifacts = append(artifacts, Artifact{
			Path:     r.Path,
			Size:     r.Size,
			Modified: r.Modified,
		})
	}

	return artifacts, nil
}

// Delete deletes artifacts from Artifactory.
func (h *Helper) Delete(pattern string, quiet bool) error {
	if pattern == "" {
		return fmt.Errorf("pattern is required")
	}

	args := []string{"rt", "delete", pattern}
	if quiet {
		args = append(args, "--quiet")
	}

	if h.dryRun {
		args = append(args, "--dry-run")
	}

	cmd := exec.Command("jf", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// DockerPush pushes a Docker image to Artifactory.
func (h *Helper) DockerPush(image, repo string) error {
	if image == "" || repo == "" {
		return fmt.Errorf("image and repository are required")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: jf docker push %s %s\n", image, repo)
		return nil
	}

	cmd := exec.Command("jf", "docker", "push", image, repo)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// DockerPull pulls a Docker image from Artifactory.
func (h *Helper) DockerPull(image, repo string) error {
	if image == "" || repo == "" {
		return fmt.Errorf("image and repository are required")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: jf docker pull %s %s\n", image, repo)
		return nil
	}

	cmd := exec.Command("jf", "docker", "pull", image, repo)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Install installs JFrog CLI.
func (h *Helper) Install() error {
	if _, err := exec.LookPath("jf"); err == nil {
		return fmt.Errorf("jfrog CLI is already installed")
	}

	if h.dryRun {
		fmt.Println("[dry-run] would install jfrog CLI via homebrew")
		return nil
	}

	fmt.Println("Installing JFrog CLI...")
	cmd := exec.Command("brew", "install", "jfrog-cli")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
