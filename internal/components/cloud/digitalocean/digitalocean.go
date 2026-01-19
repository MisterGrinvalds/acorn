// Package digitalocean provides DigitalOcean CLI (doctl) helper functionality.
package digitalocean

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Status represents doctl status information.
type Status struct {
	Installed     bool   `json:"installed" yaml:"installed"`
	Version       string `json:"version,omitempty" yaml:"version,omitempty"`
	Authenticated bool   `json:"authenticated" yaml:"authenticated"`
	AccountEmail  string `json:"account_email,omitempty" yaml:"account_email,omitempty"`
	AccountUUID   string `json:"account_uuid,omitempty" yaml:"account_uuid,omitempty"`
	TeamName      string `json:"team_name,omitempty" yaml:"team_name,omitempty"`
}

// Account represents DigitalOcean account info.
type Account struct {
	Email           string `json:"email" yaml:"email"`
	UUID            string `json:"uuid" yaml:"uuid"`
	Team            string `json:"team,omitempty" yaml:"team,omitempty"`
	DropletLimit    int    `json:"droplet_limit" yaml:"droplet_limit"`
	FloatingIPLimit int    `json:"floating_ip_limit" yaml:"floating_ip_limit"`
}

// Overview contains DigitalOcean resources summary.
type Overview struct {
	Status       *Status  `json:"status" yaml:"status"`
	DropletCount int      `json:"droplet_count,omitempty" yaml:"droplet_count,omitempty"`
	K8sClusters  []string `json:"k8s_clusters,omitempty" yaml:"k8s_clusters,omitempty"`
	Apps         []string `json:"apps,omitempty" yaml:"apps,omitempty"`
	Databases    []string `json:"databases,omitempty" yaml:"databases,omitempty"`
}

// Helper provides doctl helper operations.
type Helper struct {
	verbose bool
	dryRun  bool
	context string
}

// NewHelper creates a new DigitalOcean Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	return &Helper{
		verbose: verbose,
		dryRun:  dryRun,
	}
}

// SetContext sets the doctl context.
func (h *Helper) SetContext(context string) {
	h.context = context
}

// buildArgs builds command arguments with context.
func (h *Helper) buildArgs(args ...string) []string {
	var cmdArgs []string
	if h.context != "" {
		cmdArgs = append(cmdArgs, "--context", h.context)
	}
	cmdArgs = append(cmdArgs, args...)
	return cmdArgs
}

// GetStatus returns doctl status and authentication info.
func (h *Helper) GetStatus() (*Status, error) {
	status := &Status{}

	// Check if doctl is installed
	versionCmd := exec.Command("doctl", "version")
	versionOut, err := versionCmd.Output()
	if err != nil {
		status.Installed = false
		return status, nil
	}

	status.Installed = true

	// Parse version (format: doctl version 1.x.x-release)
	versionStr := strings.TrimSpace(string(versionOut))
	if parts := strings.Fields(versionStr); len(parts) >= 3 {
		status.Version = parts[2]
	}

	// Check authentication
	account, err := h.GetAccount()
	if err != nil {
		status.Authenticated = false
		return status, nil
	}

	status.Authenticated = true
	status.AccountEmail = account.Email
	status.AccountUUID = account.UUID
	status.TeamName = account.Team

	return status, nil
}

// GetAccount returns the current DigitalOcean account.
func (h *Helper) GetAccount() (*Account, error) {
	args := h.buildArgs("account", "get", "-o", "json")
	cmd := exec.Command("doctl", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	var account Account
	if err := json.Unmarshal(out, &account); err != nil {
		return nil, fmt.Errorf("failed to parse account: %w", err)
	}

	return &account, nil
}

// ListDroplets lists DigitalOcean droplets.
func (h *Helper) ListDroplets() (string, error) {
	args := h.buildArgs("compute", "droplet", "list")
	cmd := exec.Command("doctl", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if len(out) > 0 {
			return "", fmt.Errorf("%s", strings.TrimSpace(string(out)))
		}
		return "", fmt.Errorf("failed to list droplets: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// CountDroplets counts droplets.
func (h *Helper) CountDroplets() (int, error) {
	args := h.buildArgs("compute", "droplet", "list", "-o", "json")
	cmd := exec.Command("doctl", args...)
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	var droplets []any
	if err := json.Unmarshal(out, &droplets); err != nil {
		return 0, err
	}
	return len(droplets), nil
}

// ListK8sClusters lists Kubernetes clusters.
func (h *Helper) ListK8sClusters() ([]string, error) {
	args := h.buildArgs("kubernetes", "cluster", "list", "-o", "json")
	cmd := exec.Command("doctl", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list k8s clusters: %w", err)
	}

	var clusters []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(out, &clusters); err != nil {
		return nil, fmt.Errorf("failed to parse clusters: %w", err)
	}

	var names []string
	for _, c := range clusters {
		names = append(names, c.Name)
	}
	return names, nil
}

// GetK8sCredentials gets kubeconfig for a k8s cluster.
func (h *Helper) GetK8sCredentials(clusterName string) error {
	if clusterName == "" {
		return fmt.Errorf("cluster name is required")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: doctl kubernetes cluster kubeconfig save %s\n", clusterName)
		return nil
	}

	args := h.buildArgs("kubernetes", "cluster", "kubeconfig", "save", clusterName)
	cmd := exec.Command("doctl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ListApps lists App Platform apps.
func (h *Helper) ListApps() ([]string, error) {
	args := h.buildArgs("apps", "list", "-o", "json")
	cmd := exec.Command("doctl", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list apps: %w", err)
	}

	var apps []struct {
		Spec struct {
			Name string `json:"name"`
		} `json:"spec"`
	}
	if err := json.Unmarshal(out, &apps); err != nil {
		return nil, fmt.Errorf("failed to parse apps: %w", err)
	}

	var names []string
	for _, a := range apps {
		names = append(names, a.Spec.Name)
	}
	return names, nil
}

// ListDatabases lists managed databases.
func (h *Helper) ListDatabases() ([]string, error) {
	args := h.buildArgs("databases", "list", "-o", "json")
	cmd := exec.Command("doctl", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list databases: %w", err)
	}

	var databases []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(out, &databases); err != nil {
		return nil, fmt.Errorf("failed to parse databases: %w", err)
	}

	var names []string
	for _, d := range databases {
		names = append(names, d.Name)
	}
	return names, nil
}

// SSHToDroplet SSH into a droplet.
func (h *Helper) SSHToDroplet(dropletID string) error {
	if dropletID == "" {
		return fmt.Errorf("droplet ID is required")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: doctl compute ssh %s\n", dropletID)
		return nil
	}

	args := h.buildArgs("compute", "ssh", dropletID)
	cmd := exec.Command("doctl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// GetOverview returns an overview of DigitalOcean resources.
func (h *Helper) GetOverview() (*Overview, error) {
	overview := &Overview{}

	// Get status
	status, err := h.GetStatus()
	if err != nil {
		return nil, err
	}
	overview.Status = status

	if !status.Installed || !status.Authenticated {
		return overview, nil
	}

	// Get resources
	if count, err := h.CountDroplets(); err == nil {
		overview.DropletCount = count
	}

	if clusters, err := h.ListK8sClusters(); err == nil {
		overview.K8sClusters = clusters
	}

	if apps, err := h.ListApps(); err == nil {
		overview.Apps = apps
	}

	if dbs, err := h.ListDatabases(); err == nil {
		overview.Databases = dbs
	}

	return overview, nil
}

// Login initiates doctl authentication.
func (h *Helper) Login() error {
	if h.dryRun {
		fmt.Println("[dry-run] would run: doctl auth init")
		return nil
	}

	cmd := exec.Command("doctl", "auth", "init")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// ListContexts lists available doctl contexts.
func (h *Helper) ListContexts() (string, error) {
	cmd := exec.Command("doctl", "auth", "list")
	out, err := cmd.CombinedOutput()
	if err != nil {
		if len(out) > 0 {
			return "", fmt.Errorf("%s", strings.TrimSpace(string(out)))
		}
		return "", fmt.Errorf("failed to list contexts: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// SwitchContext switches the doctl context.
func (h *Helper) SwitchContext(contextName string) error {
	if contextName == "" {
		return fmt.Errorf("context name is required")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: doctl auth switch --context %s\n", contextName)
		return nil
	}

	cmd := exec.Command("doctl", "auth", "switch", "--context", contextName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
