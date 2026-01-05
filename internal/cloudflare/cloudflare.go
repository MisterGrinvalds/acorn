// Package cloudflare provides CloudFlare CLI (wrangler) helper functionality.
package cloudflare

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Status represents CloudFlare CLI status information.
type Status struct {
	Installed      bool   `json:"installed" yaml:"installed"`
	Version        string `json:"version,omitempty" yaml:"version,omitempty"`
	Authenticated  bool   `json:"authenticated" yaml:"authenticated"`
	AccountName    string `json:"account_name,omitempty" yaml:"account_name,omitempty"`
	AccountID      string `json:"account_id,omitempty" yaml:"account_id,omitempty"`
	WranglerHome   string `json:"wrangler_home,omitempty" yaml:"wrangler_home,omitempty"`
}

// Worker represents a CloudFlare Worker.
type Worker struct {
	Name       string `json:"name" yaml:"name"`
	ID         string `json:"id,omitempty" yaml:"id,omitempty"`
	CreatedOn  string `json:"created_on,omitempty" yaml:"created_on,omitempty"`
	ModifiedOn string `json:"modified_on,omitempty" yaml:"modified_on,omitempty"`
}

// PagesProject represents a CloudFlare Pages project.
type PagesProject struct {
	Name       string `json:"name" yaml:"name"`
	Subdomain  string `json:"subdomain,omitempty" yaml:"subdomain,omitempty"`
	CreatedOn  string `json:"created_on,omitempty" yaml:"created_on,omitempty"`
}

// R2Bucket represents a CloudFlare R2 bucket.
type R2Bucket struct {
	Name         string `json:"name" yaml:"name"`
	CreationDate string `json:"creation_date,omitempty" yaml:"creation_date,omitempty"`
}

// KVNamespace represents a CloudFlare KV namespace.
type KVNamespace struct {
	ID    string `json:"id" yaml:"id"`
	Title string `json:"title" yaml:"title"`
}

// D1Database represents a CloudFlare D1 database.
type D1Database struct {
	UUID      string `json:"uuid" yaml:"uuid"`
	Name      string `json:"name" yaml:"name"`
	CreatedAt string `json:"created_at,omitempty" yaml:"created_at,omitempty"`
}

// Overview contains all CloudFlare resources.
type Overview struct {
	Status     *Status         `json:"status" yaml:"status"`
	Workers    []string        `json:"workers,omitempty" yaml:"workers,omitempty"`
	Pages      []string        `json:"pages,omitempty" yaml:"pages,omitempty"`
	R2Buckets  []string        `json:"r2_buckets,omitempty" yaml:"r2_buckets,omitempty"`
	KV         []string        `json:"kv_namespaces,omitempty" yaml:"kv_namespaces,omitempty"`
	D1         []string        `json:"d1_databases,omitempty" yaml:"d1_databases,omitempty"`
}

// Helper provides CloudFlare CLI helper operations.
type Helper struct {
	verbose bool
	dryRun  bool
}

// NewHelper creates a new CloudFlare Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	return &Helper{
		verbose: verbose,
		dryRun:  dryRun,
	}
}

// GetStatus returns CloudFlare CLI status and authentication info.
func (h *Helper) GetStatus() (*Status, error) {
	status := &Status{
		WranglerHome: os.Getenv("WRANGLER_HOME"),
	}

	// Check if wrangler is installed
	versionCmd := exec.Command("wrangler", "--version")
	versionOut, err := versionCmd.Output()
	if err != nil {
		status.Installed = false
		return status, nil
	}

	status.Installed = true
	status.Version = strings.TrimSpace(string(versionOut))

	// Check authentication
	whoamiCmd := exec.Command("wrangler", "whoami")
	whoamiOut, err := whoamiCmd.Output()
	if err != nil {
		status.Authenticated = false
		return status, nil
	}

	whoamiStr := string(whoamiOut)
	if strings.Contains(whoamiStr, "You are logged in") {
		status.Authenticated = true
		// Parse account info from output
		lines := strings.Split(whoamiStr, "\n")
		for _, line := range lines {
			if strings.Contains(line, "Account Name:") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					status.AccountName = strings.TrimSpace(parts[1])
				}
			}
			if strings.Contains(line, "Account ID:") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					status.AccountID = strings.TrimSpace(parts[1])
				}
			}
		}
	}

	return status, nil
}

// Whoami returns the current CloudFlare account.
func (h *Helper) Whoami() (string, error) {
	cmd := exec.Command("wrangler", "whoami")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get account info: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// ListWorkers lists all Workers deployments.
func (h *Helper) ListWorkers() (string, error) {
	cmd := exec.Command("wrangler", "deployments", "list")
	out, err := cmd.CombinedOutput()
	if err != nil {
		if len(out) > 0 {
			return "", fmt.Errorf("%s", strings.TrimSpace(string(out)))
		}
		return "", fmt.Errorf("failed to list workers: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// ListPages lists all Pages projects.
func (h *Helper) ListPages() (string, error) {
	cmd := exec.Command("wrangler", "pages", "project", "list")
	out, err := cmd.CombinedOutput()
	if err != nil {
		if len(out) > 0 {
			return "", fmt.Errorf("%s", strings.TrimSpace(string(out)))
		}
		return "", fmt.Errorf("failed to list pages: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// ListR2Buckets lists all R2 buckets.
func (h *Helper) ListR2Buckets() (string, error) {
	cmd := exec.Command("wrangler", "r2", "bucket", "list")
	out, err := cmd.CombinedOutput()
	if err != nil {
		if len(out) > 0 {
			return "", fmt.Errorf("%s", strings.TrimSpace(string(out)))
		}
		return "", fmt.Errorf("failed to list R2 buckets: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// ListKVNamespaces lists all KV namespaces.
func (h *Helper) ListKVNamespaces() (string, error) {
	cmd := exec.Command("wrangler", "kv", "namespace", "list")
	out, err := cmd.CombinedOutput()
	if err != nil {
		if len(out) > 0 {
			return "", fmt.Errorf("%s", strings.TrimSpace(string(out)))
		}
		return "", fmt.Errorf("failed to list KV namespaces: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// ListD1Databases lists all D1 databases.
func (h *Helper) ListD1Databases() (string, error) {
	cmd := exec.Command("wrangler", "d1", "list")
	out, err := cmd.CombinedOutput()
	if err != nil {
		if len(out) > 0 {
			return "", fmt.Errorf("%s", strings.TrimSpace(string(out)))
		}
		return "", fmt.Errorf("failed to list D1 databases: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// ListSecrets lists secrets for the current worker.
func (h *Helper) ListSecrets() (string, error) {
	cmd := exec.Command("wrangler", "secret", "list")
	out, err := cmd.CombinedOutput()
	if err != nil {
		if len(out) > 0 {
			return "", fmt.Errorf("%s", strings.TrimSpace(string(out)))
		}
		return "", fmt.Errorf("failed to list secrets: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// TailLogs tails logs for a worker.
func (h *Helper) TailLogs(workerName string) error {
	if workerName == "" {
		return fmt.Errorf("worker name is required")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: wrangler tail %s\n", workerName)
		return nil
	}

	cmd := exec.Command("wrangler", "tail", workerName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Deploy deploys the current worker.
func (h *Helper) Deploy(args ...string) error {
	// Check for wrangler config
	if _, err := os.Stat("wrangler.toml"); os.IsNotExist(err) {
		if _, err := os.Stat("wrangler.json"); os.IsNotExist(err) {
			return fmt.Errorf("no wrangler.toml or wrangler.json found in current directory")
		}
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: wrangler deploy %s\n", strings.Join(args, " "))
		return nil
	}

	cmdArgs := append([]string{"deploy"}, args...)
	cmd := exec.Command("wrangler", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// InitWorker initializes a new Worker project.
func (h *Helper) InitWorker(name string) error {
	if name == "" {
		name = "my-worker"
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: wrangler init %s\n", name)
		return nil
	}

	cmd := exec.Command("wrangler", "init", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// InitPages initializes a new Pages project.
func (h *Helper) InitPages(name string) error {
	if name == "" {
		name = "my-pages-site"
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: wrangler pages project create %s\n", name)
		return nil
	}

	cmd := exec.Command("wrangler", "pages", "project", "create", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// CreateR2Bucket creates a new R2 bucket.
func (h *Helper) CreateR2Bucket(name string) error {
	if name == "" {
		return fmt.Errorf("bucket name is required")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: wrangler r2 bucket create %s\n", name)
		return nil
	}

	cmd := exec.Command("wrangler", "r2", "bucket", "create", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// CreateKVNamespace creates a new KV namespace.
func (h *Helper) CreateKVNamespace(name string) error {
	if name == "" {
		return fmt.Errorf("namespace name is required")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: wrangler kv namespace create %s\n", name)
		return nil
	}

	cmd := exec.Command("wrangler", "kv", "namespace", "create", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// CreateD1Database creates a new D1 database.
func (h *Helper) CreateD1Database(name string) error {
	if name == "" {
		return fmt.Errorf("database name is required")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: wrangler d1 create %s\n", name)
		return nil
	}

	cmd := exec.Command("wrangler", "d1", "create", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// PutSecret puts a secret for the current worker.
func (h *Helper) PutSecret(name string) error {
	if name == "" {
		return fmt.Errorf("secret name is required")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: wrangler secret put %s\n", name)
		return nil
	}

	cmd := exec.Command("wrangler", "secret", "put", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// GetOverview returns an overview of all CloudFlare resources.
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

	// Get resources (capture output as strings for the overview)
	if workers, err := h.captureList("wrangler", "deployments", "list"); err == nil {
		overview.Workers = h.parseLines(workers, 10)
	}

	if pages, err := h.captureList("wrangler", "pages", "project", "list"); err == nil {
		overview.Pages = h.parseLines(pages, 10)
	}

	if r2, err := h.captureList("wrangler", "r2", "bucket", "list"); err == nil {
		overview.R2Buckets = h.parseLines(r2, 10)
	}

	if kv, err := h.captureList("wrangler", "kv", "namespace", "list"); err == nil {
		overview.KV = h.parseLines(kv, 10)
	}

	if d1, err := h.captureList("wrangler", "d1", "list"); err == nil {
		overview.D1 = h.parseLines(d1, 10)
	}

	return overview, nil
}

// captureList runs a command and captures its output.
func (h *Helper) captureList(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return stdout.String(), nil
}

// parseLines splits output into lines and limits the count.
func (h *Helper) parseLines(output string, limit int) []string {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var result []string
	for i, line := range lines {
		if i >= limit {
			break
		}
		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, line)
		}
	}
	return result
}

// Login initiates CloudFlare login.
func (h *Helper) Login() error {
	if h.dryRun {
		fmt.Println("[dry-run] would run: wrangler login")
		return nil
	}

	cmd := exec.Command("wrangler", "login")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Logout logs out from CloudFlare.
func (h *Helper) Logout() error {
	if h.dryRun {
		fmt.Println("[dry-run] would run: wrangler logout")
		return nil
	}

	cmd := exec.Command("wrangler", "logout")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
