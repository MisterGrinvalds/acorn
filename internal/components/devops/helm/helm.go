// Package helm provides Helm helper functionality.
package helm

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Status represents Helm installation status.
type Status struct {
	Installed bool   `json:"installed" yaml:"installed"`
	Version   string `json:"version,omitempty" yaml:"version,omitempty"`
	Location  string `json:"location,omitempty" yaml:"location,omitempty"`
}

// Release represents a Helm release.
type Release struct {
	Name       string `json:"name" yaml:"name"`
	Namespace  string `json:"namespace" yaml:"namespace"`
	Revision   string `json:"revision" yaml:"revision"`
	Updated    string `json:"updated" yaml:"updated"`
	Status     string `json:"status" yaml:"status"`
	Chart      string `json:"chart" yaml:"chart"`
	AppVersion string `json:"app_version" yaml:"app_version"`
}

// Repository represents a Helm repository.
type Repository struct {
	Name string `json:"name" yaml:"name"`
	URL  string `json:"url" yaml:"url"`
}

// Chart represents a Helm chart.
type Chart struct {
	Name        string `json:"name" yaml:"name"`
	Version     string `json:"version" yaml:"version"`
	AppVersion  string `json:"app_version" yaml:"app_version"`
	Description string `json:"description" yaml:"description"`
}

// ReleaseHistory represents a release revision.
type ReleaseHistory struct {
	Revision    int    `json:"revision" yaml:"revision"`
	Updated     string `json:"updated" yaml:"updated"`
	Status      string `json:"status" yaml:"status"`
	Chart       string `json:"chart" yaml:"chart"`
	AppVersion  string `json:"app_version" yaml:"app_version"`
	Description string `json:"description" yaml:"description"`
}

// Plugin represents a Helm plugin.
type Plugin struct {
	Name        string `json:"name" yaml:"name"`
	Version     string `json:"version" yaml:"version"`
	Description string `json:"description" yaml:"description"`
}

// Helper provides Helm helper operations.
type Helper struct {
	verbose bool
	dryRun  bool
}

// NewHelper creates a new Helm Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	return &Helper{
		verbose: verbose,
		dryRun:  dryRun,
	}
}

// IsInstalled checks if Helm is installed.
func (h *Helper) IsInstalled() bool {
	_, err := exec.LookPath("helm")
	return err == nil
}

// GetStatus returns Helm installation status.
func (h *Helper) GetStatus() *Status {
	status := &Status{}

	helmPath, err := exec.LookPath("helm")
	if err != nil {
		return status
	}

	status.Installed = true
	status.Location = helmPath

	// Get version
	cmd := exec.Command("helm", "version", "--short")
	out, err := cmd.Output()
	if err == nil {
		status.Version = strings.TrimSpace(string(out))
	}

	return status
}

// ListReleases returns all releases.
func (h *Helper) ListReleases(namespace string, allNamespaces bool) ([]Release, error) {
	args := []string{"list", "-o", "json"}
	if allNamespaces {
		args = append(args, "-A")
	} else if namespace != "" {
		args = append(args, "-n", namespace)
	}

	cmd := exec.Command("helm", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list releases: %w", err)
	}

	var releases []Release
	if err := json.Unmarshal(out, &releases); err != nil {
		return nil, fmt.Errorf("failed to parse releases: %w", err)
	}

	return releases, nil
}

// GetRelease returns a specific release.
func (h *Helper) GetRelease(name, namespace string) (*Release, error) {
	releases, err := h.ListReleases(namespace, false)
	if err != nil {
		return nil, err
	}

	for _, r := range releases {
		if r.Name == name {
			return &r, nil
		}
	}

	return nil, fmt.Errorf("release %s not found", name)
}

// GetReleaseStatus returns the status of a release.
func (h *Helper) GetReleaseStatus(name, namespace string) error {
	args := []string{"status", name}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}

	cmd := exec.Command("helm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetReleaseValues returns the values of a release.
func (h *Helper) GetReleaseValues(name, namespace string, allValues bool) (string, error) {
	args := []string{"get", "values", name}
	if allValues {
		args = append(args, "--all")
	}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}

	cmd := exec.Command("helm", args...)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get values: %w", err)
	}

	return string(out), nil
}

// GetReleaseManifest returns the manifest of a release.
func (h *Helper) GetReleaseManifest(name, namespace string) (string, error) {
	args := []string{"get", "manifest", name}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}

	cmd := exec.Command("helm", args...)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get manifest: %w", err)
	}

	return string(out), nil
}

// GetReleaseHistory returns the history of a release.
func (h *Helper) GetReleaseHistory(name, namespace string) ([]ReleaseHistory, error) {
	args := []string{"history", name, "-o", "json"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}

	cmd := exec.Command("helm", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get history: %w", err)
	}

	var history []ReleaseHistory
	if err := json.Unmarshal(out, &history); err != nil {
		return nil, fmt.Errorf("failed to parse history: %w", err)
	}

	return history, nil
}

// Install installs a chart.
func (h *Helper) Install(release, chart, namespace string, values []string, wait bool) error {
	args := []string{"install", release, chart}
	if namespace != "" {
		args = append(args, "-n", namespace, "--create-namespace")
	}
	for _, v := range values {
		args = append(args, "-f", v)
	}
	if wait {
		args = append(args, "--wait")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: helm %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("helm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Upgrade upgrades a release.
func (h *Helper) Upgrade(release, chart, namespace string, values []string, install, wait, atomic bool) error {
	args := []string{"upgrade", release, chart}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	if install {
		args = append(args, "--install")
	}
	for _, v := range values {
		args = append(args, "-f", v)
	}
	if wait {
		args = append(args, "--wait")
	}
	if atomic {
		args = append(args, "--atomic")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: helm %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("helm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Uninstall uninstalls a release.
func (h *Helper) Uninstall(release, namespace string) error {
	args := []string{"uninstall", release}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: helm %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("helm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Rollback rolls back a release to a previous revision.
func (h *Helper) Rollback(release, namespace string, revision int) error {
	args := []string{"rollback", release}
	if revision > 0 {
		args = append(args, fmt.Sprintf("%d", revision))
	}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: helm %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("helm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ListRepositories returns all configured repositories.
func (h *Helper) ListRepositories() ([]Repository, error) {
	cmd := exec.Command("helm", "repo", "list", "-o", "json")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list repos: %w", err)
	}

	var repos []Repository
	if err := json.Unmarshal(out, &repos); err != nil {
		return nil, fmt.Errorf("failed to parse repos: %w", err)
	}

	return repos, nil
}

// AddRepository adds a repository.
func (h *Helper) AddRepository(name, url string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would run: helm repo add %s %s\n", name, url)
		return nil
	}

	cmd := exec.Command("helm", "repo", "add", name, url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// UpdateRepositories updates all repositories.
func (h *Helper) UpdateRepositories() error {
	if h.dryRun {
		fmt.Println("[dry-run] would run: helm repo update")
		return nil
	}

	cmd := exec.Command("helm", "repo", "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RemoveRepository removes a repository.
func (h *Helper) RemoveRepository(name string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would run: helm repo remove %s\n", name)
		return nil
	}

	cmd := exec.Command("helm", "repo", "remove", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// SearchCharts searches for charts in repositories.
func (h *Helper) SearchCharts(query string) ([]Chart, error) {
	args := []string{"search", "repo", query, "-o", "json"}

	cmd := exec.Command("helm", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to search charts: %w", err)
	}

	var charts []Chart
	if err := json.Unmarshal(out, &charts); err != nil {
		return nil, fmt.Errorf("failed to parse charts: %w", err)
	}

	return charts, nil
}

// ShowChart shows chart information.
func (h *Helper) ShowChart(chart string) error {
	cmd := exec.Command("helm", "show", "chart", chart)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ShowValues shows chart default values.
func (h *Helper) ShowValues(chart string) error {
	cmd := exec.Command("helm", "show", "values", chart)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Template renders chart templates locally.
func (h *Helper) Template(release, chart string, values []string) error {
	args := []string{"template", release, chart}
	for _, v := range values {
		args = append(args, "-f", v)
	}

	cmd := exec.Command("helm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Lint lints a chart.
func (h *Helper) Lint(chart string, strict bool) error {
	args := []string{"lint", chart}
	if strict {
		args = append(args, "--strict")
	}

	cmd := exec.Command("helm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Package packages a chart.
func (h *Helper) Package(chart, destination string) error {
	args := []string{"package", chart}
	if destination != "" {
		args = append(args, "-d", destination)
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: helm %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("helm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Create creates a new chart.
func (h *Helper) Create(name string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would run: helm create %s\n", name)
		return nil
	}

	cmd := exec.Command("helm", "create", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ListPlugins returns installed plugins.
func (h *Helper) ListPlugins() ([]Plugin, error) {
	cmd := exec.Command("helm", "plugin", "list")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list plugins: %w", err)
	}

	var plugins []Plugin
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	// Skip header line
	for i := 1; i < len(lines); i++ {
		fields := strings.Fields(lines[i])
		if len(fields) >= 2 {
			p := Plugin{
				Name:    fields[0],
				Version: fields[1],
			}
			if len(fields) > 2 {
				p.Description = strings.Join(fields[2:], " ")
			}
			plugins = append(plugins, p)
		}
	}

	return plugins, nil
}

// InstallPlugin installs a plugin.
func (h *Helper) InstallPlugin(url string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would run: helm plugin install %s\n", url)
		return nil
	}

	cmd := exec.Command("helm", "plugin", "install", url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// UninstallPlugin uninstalls a plugin.
func (h *Helper) UninstallPlugin(name string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would run: helm plugin uninstall %s\n", name)
		return nil
	}

	cmd := exec.Command("helm", "plugin", "uninstall", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// DependencyUpdate updates chart dependencies.
func (h *Helper) DependencyUpdate(chart string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would run: helm dependency update %s\n", chart)
		return nil
	}

	cmd := exec.Command("helm", "dependency", "update", chart)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// DependencyBuild builds chart dependencies.
func (h *Helper) DependencyBuild(chart string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would run: helm dependency build %s\n", chart)
		return nil
	}

	cmd := exec.Command("helm", "dependency", "build", chart)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Test runs tests for a release.
func (h *Helper) Test(release, namespace string) error {
	args := []string{"test", release}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: helm %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("helm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
