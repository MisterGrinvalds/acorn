// Package azure provides Azure CLI helper functionality.
package azure

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Status represents Azure CLI status information.
type Status struct {
	Installed        bool   `json:"installed" yaml:"installed"`
	Version          string `json:"version,omitempty" yaml:"version,omitempty"`
	Authenticated    bool   `json:"authenticated" yaml:"authenticated"`
	SubscriptionID   string `json:"subscription_id,omitempty" yaml:"subscription_id,omitempty"`
	SubscriptionName string `json:"subscription_name,omitempty" yaml:"subscription_name,omitempty"`
	TenantID         string `json:"tenant_id,omitempty" yaml:"tenant_id,omitempty"`
	UserName         string `json:"user_name,omitempty" yaml:"user_name,omitempty"`
}

// Subscription represents an Azure subscription.
type Subscription struct {
	ID        string `json:"id" yaml:"id"`
	Name      string `json:"name" yaml:"name"`
	State     string `json:"state" yaml:"state"`
	IsDefault bool   `json:"isDefault" yaml:"isDefault"`
	TenantID  string `json:"tenantId" yaml:"tenantId"`
}

// Account represents the current Azure account.
type Account struct {
	SubscriptionID   string `json:"id" yaml:"id"`
	SubscriptionName string `json:"name" yaml:"name"`
	TenantID         string `json:"tenantId" yaml:"tenantId"`
	User             struct {
		Name string `json:"name" yaml:"name"`
		Type string `json:"type" yaml:"type"`
	} `json:"user" yaml:"user"`
}

// Overview contains Azure resources summary.
type Overview struct {
	Status          *Status  `json:"status" yaml:"status"`
	VMCount         int      `json:"vm_count,omitempty" yaml:"vm_count,omitempty"`
	StorageAccounts []string `json:"storage_accounts,omitempty" yaml:"storage_accounts,omitempty"`
	AKSClusters     []string `json:"aks_clusters,omitempty" yaml:"aks_clusters,omitempty"`
	ResourceGroups  []string `json:"resource_groups,omitempty" yaml:"resource_groups,omitempty"`
}

// Helper provides Azure CLI helper operations.
type Helper struct {
	verbose      bool
	dryRun       bool
	subscription string
}

// NewHelper creates a new Azure Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	return &Helper{
		verbose: verbose,
		dryRun:  dryRun,
	}
}

// SetSubscription sets the Azure subscription.
func (h *Helper) SetSubscription(subscription string) {
	h.subscription = subscription
}

// buildArgs builds command arguments with subscription.
func (h *Helper) buildArgs(args ...string) []string {
	var cmdArgs []string
	if h.subscription != "" {
		cmdArgs = append(cmdArgs, "--subscription", h.subscription)
	}
	cmdArgs = append(cmdArgs, args...)
	return cmdArgs
}

// GetStatus returns Azure CLI status and authentication info.
func (h *Helper) GetStatus() (*Status, error) {
	status := &Status{}

	// Check if az is installed
	versionCmd := exec.Command("az", "version", "--output", "json")
	versionOut, err := versionCmd.Output()
	if err != nil {
		status.Installed = false
		return status, nil
	}

	status.Installed = true

	// Parse version
	var versionInfo map[string]interface{}
	if err := json.Unmarshal(versionOut, &versionInfo); err == nil {
		if cliVersion, ok := versionInfo["azure-cli"].(string); ok {
			status.Version = cliVersion
		}
	}

	// Check authentication
	account, err := h.GetAccount()
	if err != nil {
		status.Authenticated = false
		return status, nil
	}

	status.Authenticated = true
	status.SubscriptionID = account.SubscriptionID
	status.SubscriptionName = account.SubscriptionName
	status.TenantID = account.TenantID
	status.UserName = account.User.Name

	return status, nil
}

// GetAccount returns the current Azure account.
func (h *Helper) GetAccount() (*Account, error) {
	args := h.buildArgs("account", "show", "--output", "json")
	cmd := exec.Command("az", args...)
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

// ListSubscriptions lists all Azure subscriptions.
func (h *Helper) ListSubscriptions() ([]Subscription, error) {
	cmd := exec.Command("az", "account", "list", "--output", "json")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}

	var subscriptions []Subscription
	if err := json.Unmarshal(out, &subscriptions); err != nil {
		return nil, fmt.Errorf("failed to parse subscriptions: %w", err)
	}
	return subscriptions, nil
}

// SetActiveSubscription sets the active subscription.
func (h *Helper) SetActiveSubscription(subscriptionID string) error {
	if subscriptionID == "" {
		return fmt.Errorf("subscription ID is required")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: az account set --subscription %s\n", subscriptionID)
		return nil
	}

	cmd := exec.Command("az", "account", "set", "--subscription", subscriptionID)
	return cmd.Run()
}

// ListVMs lists Azure VMs.
func (h *Helper) ListVMs() (string, error) {
	args := h.buildArgs("vm", "list", "-o", "table")
	cmd := exec.Command("az", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if len(out) > 0 {
			return "", fmt.Errorf("%s", strings.TrimSpace(string(out)))
		}
		return "", fmt.Errorf("failed to list VMs: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// ListStorageAccounts lists storage accounts.
func (h *Helper) ListStorageAccounts() ([]string, error) {
	args := h.buildArgs("storage", "account", "list", "--query", "[].name", "-o", "json")
	cmd := exec.Command("az", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list storage accounts: %w", err)
	}

	var accounts []string
	if err := json.Unmarshal(out, &accounts); err != nil {
		return nil, fmt.Errorf("failed to parse storage accounts: %w", err)
	}
	return accounts, nil
}

// ListAKSClusters lists AKS clusters.
func (h *Helper) ListAKSClusters() ([]string, error) {
	args := h.buildArgs("aks", "list", "--query", "[].name", "-o", "json")
	cmd := exec.Command("az", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list AKS clusters: %w", err)
	}

	var clusters []string
	if err := json.Unmarshal(out, &clusters); err != nil {
		return nil, fmt.Errorf("failed to parse AKS clusters: %w", err)
	}
	return clusters, nil
}

// ListResourceGroups lists resource groups.
func (h *Helper) ListResourceGroups() ([]string, error) {
	args := h.buildArgs("group", "list", "--query", "[].name", "-o", "json")
	cmd := exec.Command("az", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list resource groups: %w", err)
	}

	var groups []string
	if err := json.Unmarshal(out, &groups); err != nil {
		return nil, fmt.Errorf("failed to parse resource groups: %w", err)
	}
	return groups, nil
}

// GetAKSCredentials gets kubeconfig for an AKS cluster.
func (h *Helper) GetAKSCredentials(clusterName, resourceGroup string) error {
	if clusterName == "" || resourceGroup == "" {
		return fmt.Errorf("cluster name and resource group are required")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: az aks get-credentials --name %s --resource-group %s\n", clusterName, resourceGroup)
		return nil
	}

	args := h.buildArgs("aks", "get-credentials", "--name", clusterName, "--resource-group", resourceGroup, "--overwrite-existing")
	cmd := exec.Command("az", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetOverview returns an overview of Azure resources.
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
	if storage, err := h.ListStorageAccounts(); err == nil {
		overview.StorageAccounts = storage
	}

	if clusters, err := h.ListAKSClusters(); err == nil {
		overview.AKSClusters = clusters
	}

	if groups, err := h.ListResourceGroups(); err == nil {
		overview.ResourceGroups = groups
	}

	// Count VMs
	if count, err := h.countVMs(); err == nil {
		overview.VMCount = count
	}

	return overview, nil
}

// countVMs counts running VMs.
func (h *Helper) countVMs() (int, error) {
	args := h.buildArgs("vm", "list", "--query", "length(@)", "-o", "json")
	cmd := exec.Command("az", args...)
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	var count int
	if err := json.Unmarshal(out, &count); err != nil {
		return 0, err
	}
	return count, nil
}

// Login initiates Azure login.
func (h *Helper) Login() error {
	if h.dryRun {
		fmt.Println("[dry-run] would run: az login")
		return nil
	}

	cmd := exec.Command("az", "login")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Logout logs out from Azure.
func (h *Helper) Logout() error {
	if h.dryRun {
		fmt.Println("[dry-run] would run: az logout")
		return nil
	}

	cmd := exec.Command("az", "logout")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
