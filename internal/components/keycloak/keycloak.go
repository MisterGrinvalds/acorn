// Package keycloak provides Keycloak IAM helper functionality.
package keycloak

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Status represents Keycloak installation status.
type Status struct {
	DockerInstalled bool   `json:"docker_installed" yaml:"docker_installed"`
	ContainerRunning bool   `json:"container_running" yaml:"container_running"`
	ContainerID     string `json:"container_id,omitempty" yaml:"container_id,omitempty"`
	ContainerName   string `json:"container_name,omitempty" yaml:"container_name,omitempty"`
	AdminCLI        bool   `json:"admin_cli" yaml:"admin_cli"`
	AdminCLIPath    string `json:"admin_cli_path,omitempty" yaml:"admin_cli_path,omitempty"`
	RealmsDir       string `json:"realms_dir,omitempty" yaml:"realms_dir,omitempty"`
	Port            string `json:"port,omitempty" yaml:"port,omitempty"`
}

// Realm represents a Keycloak realm.
type Realm struct {
	Name    string `json:"realm" yaml:"realm"`
	Enabled bool   `json:"enabled" yaml:"enabled"`
	Path    string `json:"path,omitempty" yaml:"path,omitempty"`
}

// Client represents a Keycloak client.
type Client struct {
	ID       string `json:"id" yaml:"id"`
	ClientID string `json:"clientId" yaml:"client_id"`
	Name     string `json:"name,omitempty" yaml:"name,omitempty"`
	Enabled  bool   `json:"enabled" yaml:"enabled"`
}

// Helper provides Keycloak helper operations.
type Helper struct {
	verbose bool
	dryRun  bool
}

// NewHelper creates a new Keycloak Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	return &Helper{
		verbose: verbose,
		dryRun:  dryRun,
	}
}

// GetStatus returns Keycloak status information.
func (h *Helper) GetStatus() *Status {
	status := &Status{}

	// Check Docker
	if _, err := exec.LookPath("docker"); err == nil {
		status.DockerInstalled = true
	}

	// Check for running container
	containerID, containerName := h.getRunningContainer()
	if containerID != "" {
		status.ContainerRunning = true
		status.ContainerID = containerID
		status.ContainerName = containerName
		status.Port = h.getContainerPort(containerID)
	}

	// Check for admin CLI
	if cli := h.findAdminCLI(); cli != "" {
		status.AdminCLI = true
		status.AdminCLIPath = cli
	}

	// Get realms directory
	status.RealmsDir = h.getRealmsDir()

	return status
}

// getRunningContainer returns the running Keycloak container.
func (h *Helper) getRunningContainer() (string, string) {
	// Try by image
	cmd := exec.Command("docker", "ps", "--filter", "ancestor=quay.io/keycloak/keycloak",
		"--format", "{{.ID}}\t{{.Names}}")
	out, err := cmd.Output()
	if err == nil && strings.TrimSpace(string(out)) != "" {
		parts := strings.Split(strings.TrimSpace(string(out)), "\t")
		if len(parts) >= 2 {
			return parts[0], parts[1]
		}
		return parts[0], ""
	}

	// Try by name
	cmd = exec.Command("docker", "ps", "--filter", "name=keycloak",
		"--format", "{{.ID}}\t{{.Names}}")
	out, err = cmd.Output()
	if err == nil && strings.TrimSpace(string(out)) != "" {
		parts := strings.Split(strings.TrimSpace(string(out)), "\t")
		if len(parts) >= 2 {
			return parts[0], parts[1]
		}
		return parts[0], ""
	}

	return "", ""
}

// getContainerPort returns the mapped port for a container.
func (h *Helper) getContainerPort(containerID string) string {
	cmd := exec.Command("docker", "port", containerID, "8080")
	out, err := cmd.Output()
	if err == nil {
		// Format: 0.0.0.0:8080
		parts := strings.Split(strings.TrimSpace(string(out)), ":")
		if len(parts) >= 2 {
			return parts[len(parts)-1]
		}
	}
	return "8080"
}

// findAdminCLI finds the Keycloak admin CLI.
func (h *Helper) findAdminCLI() string {
	// Check if kcadm.sh is in PATH
	if path, err := exec.LookPath("kcadm.sh"); err == nil {
		return path
	}

	// Check common locations
	paths := []string{
		"/opt/keycloak/bin/kcadm.sh",
		"/usr/local/keycloak/bin/kcadm.sh",
	}

	home, _ := os.UserHomeDir()
	paths = append(paths, filepath.Join(home, ".keycloak/bin/kcadm.sh"))

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

// getRealmsDir returns the realms export directory.
func (h *Helper) getRealmsDir() string {
	home, _ := os.UserHomeDir()

	xdgConfig := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfig == "" {
		xdgConfig = filepath.Join(home, ".config")
	}

	return filepath.Join(xdgConfig, "keycloak", "realms")
}

// Start starts Keycloak in Docker.
func (h *Helper) Start(port int, devMode bool, adminUser, adminPassword string) error {
	if port == 0 {
		port = 8080
	}
	if adminUser == "" {
		adminUser = "admin"
	}
	if adminPassword == "" {
		adminPassword = "admin"
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would start Keycloak on port %d\n", port)
		return nil
	}

	// Check if already running
	if id, _ := h.getRunningContainer(); id != "" {
		return fmt.Errorf("Keycloak container already running: %s", id)
	}

	args := []string{
		"run", "-d",
		"--name", "keycloak",
		"-p", fmt.Sprintf("%d:8080", port),
		"-e", fmt.Sprintf("KEYCLOAK_ADMIN=%s", adminUser),
		"-e", fmt.Sprintf("KEYCLOAK_ADMIN_PASSWORD=%s", adminPassword),
		"quay.io/keycloak/keycloak:latest",
	}

	if devMode {
		args = append(args, "start-dev")
	} else {
		args = append(args, "start", "--optimized")
	}

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Starting Keycloak on http://localhost:%d ...\n", port)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start Keycloak: %w", err)
	}

	fmt.Printf("Keycloak started. Admin console: http://localhost:%d/admin\n", port)
	fmt.Printf("Admin credentials: %s / %s\n", adminUser, adminPassword)
	return nil
}

// Stop stops Keycloak container.
func (h *Helper) Stop() error {
	if h.dryRun {
		fmt.Println("[dry-run] would stop Keycloak container")
		return nil
	}

	containerID, _ := h.getRunningContainer()
	if containerID == "" {
		return fmt.Errorf("no Keycloak container running")
	}

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

// Logs shows Keycloak container logs.
func (h *Helper) Logs(follow bool) error {
	containerID, _ := h.getRunningContainer()
	if containerID == "" {
		return fmt.Errorf("no Keycloak container running")
	}

	args := []string{"logs"}
	if follow {
		args = append(args, "-f")
	}
	args = append(args, containerID)

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ExportRealm exports a realm to a JSON file.
func (h *Helper) ExportRealm(realmName, outputPath string) error {
	if realmName == "" {
		return fmt.Errorf("realm name is required")
	}

	if outputPath == "" {
		outputPath = filepath.Join(h.getRealmsDir(), realmName+".json")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would export realm %s to %s\n", realmName, outputPath)
		return nil
	}

	containerID, _ := h.getRunningContainer()
	if containerID == "" {
		return fmt.Errorf("Keycloak container not running")
	}

	// Create output directory
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Export realm using container's kc.sh
	cmd := exec.Command("docker", "exec", containerID,
		"/opt/keycloak/bin/kc.sh", "export",
		"--realm", realmName,
		"--file", "/tmp/realm-export.json")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to export realm: %w", err)
	}

	// Copy file from container
	cmd = exec.Command("docker", "cp",
		fmt.Sprintf("%s:/tmp/realm-export.json", containerID),
		outputPath)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to copy export file: %w", err)
	}

	fmt.Printf("Exported realm to: %s\n", outputPath)
	return nil
}

// ImportRealm imports a realm from a JSON file.
func (h *Helper) ImportRealm(inputPath string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would import realm from %s\n", inputPath)
		return nil
	}

	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return fmt.Errorf("realm file not found: %s", inputPath)
	}

	containerID, _ := h.getRunningContainer()
	if containerID == "" {
		return fmt.Errorf("Keycloak container not running")
	}

	// Copy file to container
	cmd := exec.Command("docker", "cp", inputPath,
		fmt.Sprintf("%s:/tmp/realm-import.json", containerID))
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to copy import file: %w", err)
	}

	// Import realm
	cmd = exec.Command("docker", "exec", containerID,
		"/opt/keycloak/bin/kc.sh", "import",
		"--file", "/tmp/realm-import.json")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to import realm: %w", err)
	}

	return nil
}

// ListRealms lists exported realm files.
func (h *Helper) ListRealms() ([]Realm, error) {
	realmsDir := h.getRealmsDir()
	if _, err := os.Stat(realmsDir); os.IsNotExist(err) {
		return []Realm{}, nil
	}

	var realms []Realm

	entries, err := os.ReadDir(realmsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read realms directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".json") {
			continue
		}

		path := filepath.Join(realmsDir, name)
		realm := Realm{
			Name: strings.TrimSuffix(name, ".json"),
			Path: path,
		}

		// Try to parse realm info
		if info := h.parseRealmInfo(path); info != nil {
			realm.Name = info.Name
			realm.Enabled = info.Enabled
		}

		realms = append(realms, realm)
	}

	return realms, nil
}

// realmInfo holds parsed realm metadata.
type realmInfo struct {
	Name    string
	Enabled bool
}

// parseRealmInfo parses realm metadata from a JSON file.
func (h *Helper) parseRealmInfo(path string) *realmInfo {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var parsed struct {
		Realm   string `json:"realm"`
		Enabled bool   `json:"enabled"`
	}

	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil
	}

	return &realmInfo{
		Name:    parsed.Realm,
		Enabled: parsed.Enabled,
	}
}

// Open opens Keycloak admin console in browser.
func (h *Helper) Open() error {
	status := h.GetStatus()
	if !status.ContainerRunning {
		return fmt.Errorf("Keycloak container not running")
	}

	url := fmt.Sprintf("http://localhost:%s/admin", status.Port)

	cmd := exec.Command("open", url)
	if err := cmd.Run(); err != nil {
		// Try xdg-open for Linux
		cmd = exec.Command("xdg-open", url)
		return cmd.Run()
	}

	return nil
}
