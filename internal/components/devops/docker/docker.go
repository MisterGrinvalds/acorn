// Package docker provides Docker helper functionality.
package docker

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Container represents a Docker container.
type Container struct {
	ID      string `json:"id" yaml:"id"`
	Names   string `json:"names" yaml:"names"`
	Image   string `json:"image" yaml:"image"`
	Status  string `json:"status" yaml:"status"`
	Ports   string `json:"ports,omitempty" yaml:"ports,omitempty"`
	Created string `json:"created,omitempty" yaml:"created,omitempty"`
}

// Image represents a Docker image.
type Image struct {
	Repository string `json:"repository" yaml:"repository"`
	Tag        string `json:"tag" yaml:"tag"`
	ImageID    string `json:"image_id" yaml:"image_id"`
	Size       string `json:"size" yaml:"size"`
	Created    string `json:"created,omitempty" yaml:"created,omitempty"`
}

// Volume represents a Docker volume.
type Volume struct {
	Name       string `json:"name" yaml:"name"`
	Driver     string `json:"driver" yaml:"driver"`
	Mountpoint string `json:"mountpoint,omitempty" yaml:"mountpoint,omitempty"`
}

// Network represents a Docker network.
type Network struct {
	Name   string `json:"name" yaml:"name"`
	ID     string `json:"id" yaml:"id"`
	Driver string `json:"driver" yaml:"driver"`
	Scope  string `json:"scope" yaml:"scope"`
}

// Status represents Docker daemon status.
type Status struct {
	Installed     bool   `json:"installed" yaml:"installed"`
	Running       bool   `json:"running" yaml:"running"`
	Version       string `json:"version,omitempty" yaml:"version,omitempty"`
	Containers    int    `json:"containers" yaml:"containers"`
	RunningCount  int    `json:"running_count" yaml:"running_count"`
	Images        int    `json:"images" yaml:"images"`
	ComposeExists bool   `json:"compose_exists" yaml:"compose_exists"`
}

// CleanupResult represents the result of a cleanup operation.
type CleanupResult struct {
	ContainersRemoved int    `json:"containers_removed" yaml:"containers_removed"`
	ImagesRemoved     int    `json:"images_removed" yaml:"images_removed"`
	VolumesRemoved    int    `json:"volumes_removed" yaml:"volumes_removed"`
	NetworksRemoved   int    `json:"networks_removed" yaml:"networks_removed"`
	SpaceReclaimed    string `json:"space_reclaimed,omitempty" yaml:"space_reclaimed,omitempty"`
}

// Helper provides Docker helper operations.
type Helper struct {
	verbose bool
	dryRun  bool
}

// NewHelper creates a new Docker Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	return &Helper{
		verbose: verbose,
		dryRun:  dryRun,
	}
}

// IsDockerInstalled checks if Docker is installed.
func (h *Helper) IsDockerInstalled() bool {
	_, err := exec.LookPath("docker")
	return err == nil
}

// IsDockerRunning checks if Docker daemon is running.
func (h *Helper) IsDockerRunning() bool {
	cmd := exec.Command("docker", "info")
	return cmd.Run() == nil
}

// GetStatus returns Docker installation and daemon status.
func (h *Helper) GetStatus() *Status {
	status := &Status{}

	// Check if docker is installed
	if _, err := exec.LookPath("docker"); err != nil {
		return status
	}
	status.Installed = true

	// Check if docker compose exists
	cmd := exec.Command("docker", "compose", "version")
	if cmd.Run() == nil {
		status.ComposeExists = true
	}

	// Check if daemon is running and get version
	cmd = exec.Command("docker", "version", "--format", "{{.Server.Version}}")
	if out, err := cmd.Output(); err == nil {
		status.Running = true
		status.Version = strings.TrimSpace(string(out))
	}

	// Get container and image counts
	if status.Running {
		cmd = exec.Command("docker", "info", "--format", "{{.Containers}} {{.ContainersRunning}} {{.Images}}")
		if out, err := cmd.Output(); err == nil {
			fmt.Sscanf(string(out), "%d %d %d", &status.Containers, &status.RunningCount, &status.Images)
		}
	}

	return status
}

// GetContainers returns list of containers.
func (h *Helper) GetContainers(all bool) ([]Container, error) {
	args := []string{"ps", "--format", "{{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}"}
	if all {
		args = append(args, "-a")
	}

	cmd := exec.Command("docker", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get containers: %w", err)
	}

	var containers []Container
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) >= 4 {
			c := Container{
				ID:     parts[0],
				Names:  parts[1],
				Image:  parts[2],
				Status: parts[3],
			}
			if len(parts) > 4 {
				c.Ports = parts[4]
			}
			containers = append(containers, c)
		}
	}

	return containers, nil
}

// GetImages returns list of images.
func (h *Helper) GetImages() ([]Image, error) {
	cmd := exec.Command("docker", "images", "--format", "{{.Repository}}\t{{.Tag}}\t{{.ID}}\t{{.Size}}")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get images: %w", err)
	}

	var images []Image
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) >= 4 {
			images = append(images, Image{
				Repository: parts[0],
				Tag:        parts[1],
				ImageID:    parts[2],
				Size:       parts[3],
			})
		}
	}

	return images, nil
}

// GetVolumes returns list of volumes.
func (h *Helper) GetVolumes() ([]Volume, error) {
	cmd := exec.Command("docker", "volume", "ls", "--format", "{{.Name}}\t{{.Driver}}")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get volumes: %w", err)
	}

	var volumes []Volume
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) >= 2 {
			volumes = append(volumes, Volume{
				Name:   parts[0],
				Driver: parts[1],
			})
		}
	}

	return volumes, nil
}

// GetNetworks returns list of networks.
func (h *Helper) GetNetworks() ([]Network, error) {
	cmd := exec.Command("docker", "network", "ls", "--format", "{{.Name}}\t{{.ID}}\t{{.Driver}}\t{{.Scope}}")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get networks: %w", err)
	}

	var networks []Network
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) >= 4 {
			networks = append(networks, Network{
				Name:   parts[0],
				ID:     parts[1],
				Driver: parts[2],
				Scope:  parts[3],
			})
		}
	}

	return networks, nil
}

// ExecInContainer executes a command in a container.
func (h *Helper) ExecInContainer(container string, command []string, interactive bool) error {
	args := []string{"exec"}
	if interactive {
		args = append(args, "-it")
	}
	args = append(args, container)
	args = append(args, command...)

	if h.dryRun {
		fmt.Printf("[dry-run] would run: docker %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("docker", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetLogs gets logs from a container.
func (h *Helper) GetLogs(container string, follow bool, tail int) error {
	args := []string{"logs"}
	if follow {
		args = append(args, "-f")
	}
	if tail > 0 {
		args = append(args, "--tail", fmt.Sprintf("%d", tail))
	}
	args = append(args, container)

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// StopContainer stops a container.
func (h *Helper) StopContainer(container string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would run: docker stop %s\n", container)
		return nil
	}

	cmd := exec.Command("docker", "stop", container)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RemoveContainer removes a container.
func (h *Helper) RemoveContainer(container string, force bool) error {
	args := []string{"rm"}
	if force {
		args = append(args, "-f")
	}
	args = append(args, container)

	if h.dryRun {
		fmt.Printf("[dry-run] would run: docker %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Cleanup removes unused Docker resources.
func (h *Helper) Cleanup(all bool) (*CleanupResult, error) {
	result := &CleanupResult{}

	// Prune stopped containers
	args := []string{"container", "prune", "-f"}
	if h.dryRun {
		fmt.Println("[dry-run] would run: docker container prune -f")
	} else {
		cmd := exec.Command("docker", args...)
		if out, err := cmd.Output(); err == nil {
			result.ContainersRemoved = strings.Count(string(out), "Deleted")
		}
	}

	// Prune dangling images (or all unused if all=true)
	args = []string{"image", "prune", "-f"}
	if all {
		args = append(args, "-a")
	}
	if h.dryRun {
		fmt.Printf("[dry-run] would run: docker %s\n", strings.Join(args, " "))
	} else {
		cmd := exec.Command("docker", args...)
		if out, err := cmd.Output(); err == nil {
			result.ImagesRemoved = strings.Count(string(out), "Deleted")
		}
	}

	// Prune unused volumes
	args = []string{"volume", "prune", "-f"}
	if h.dryRun {
		fmt.Println("[dry-run] would run: docker volume prune -f")
	} else {
		cmd := exec.Command("docker", args...)
		if out, err := cmd.Output(); err == nil {
			result.VolumesRemoved = strings.Count(string(out), "Deleted")
		}
	}

	// Prune unused networks
	args = []string{"network", "prune", "-f"}
	if h.dryRun {
		fmt.Println("[dry-run] would run: docker network prune -f")
	} else {
		cmd := exec.Command("docker", args...)
		if out, err := cmd.Output(); err == nil {
			result.NetworksRemoved = strings.Count(string(out), "Deleted")
		}
	}

	return result, nil
}

// SystemPrune runs docker system prune.
func (h *Helper) SystemPrune(all, volumes bool) (*CleanupResult, error) {
	args := []string{"system", "prune", "-f"}
	if all {
		args = append(args, "-a")
	}
	if volumes {
		args = append(args, "--volumes")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: docker %s\n", strings.Join(args, " "))
		return &CleanupResult{}, nil
	}

	cmd := exec.Command("docker", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("docker system prune failed: %w", err)
	}

	result := &CleanupResult{}
	// Parse output for space reclaimed
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "reclaimed") {
			result.SpaceReclaimed = strings.TrimSpace(line)
		}
	}

	return result, nil
}

// ComposeUp runs docker compose up.
func (h *Helper) ComposeUp(file string, detach, build bool, services []string) error {
	args := []string{"compose"}
	if file != "" {
		args = append(args, "-f", file)
	}
	args = append(args, "up")
	if detach {
		args = append(args, "-d")
	}
	if build {
		args = append(args, "--build")
	}
	args = append(args, services...)

	if h.dryRun {
		fmt.Printf("[dry-run] would run: docker %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ComposeDown runs docker compose down.
func (h *Helper) ComposeDown(file string, removeVolumes, removeOrphans bool) error {
	args := []string{"compose"}
	if file != "" {
		args = append(args, "-f", file)
	}
	args = append(args, "down")
	if removeVolumes {
		args = append(args, "-v")
	}
	if removeOrphans {
		args = append(args, "--remove-orphans")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: docker %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetComposeLogs gets logs from compose services.
func (h *Helper) GetComposeLogs(file string, follow bool, services []string) error {
	args := []string{"compose"}
	if file != "" {
		args = append(args, "-f", file)
	}
	args = append(args, "logs")
	if follow {
		args = append(args, "-f")
	}
	args = append(args, services...)

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// BuildImage builds a Docker image.
func (h *Helper) BuildImage(dockerfile, tag, context string) error {
	args := []string{"build"}
	if dockerfile != "" {
		args = append(args, "-f", dockerfile)
	}
	if tag != "" {
		args = append(args, "-t", tag)
	}
	if context == "" {
		context = "."
	}
	args = append(args, context)

	if h.dryRun {
		fmt.Printf("[dry-run] would run: docker %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// InspectContainer returns detailed container info as JSON.
func (h *Helper) InspectContainer(container string) (map[string]interface{}, error) {
	cmd := exec.Command("docker", "inspect", container)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to inspect container: %w", err)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(out, &result); err != nil {
		return nil, fmt.Errorf("failed to parse inspect output: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("container not found: %s", container)
	}

	return result[0], nil
}
