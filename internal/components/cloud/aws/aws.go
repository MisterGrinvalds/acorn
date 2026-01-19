// Package aws provides AWS CLI helper functionality.
package aws

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Status represents AWS CLI status information.
type Status struct {
	Installed     bool   `json:"installed" yaml:"installed"`
	Version       string `json:"version,omitempty" yaml:"version,omitempty"`
	Authenticated bool   `json:"authenticated" yaml:"authenticated"`
	AccountID     string `json:"account_id,omitempty" yaml:"account_id,omitempty"`
	UserARN       string `json:"user_arn,omitempty" yaml:"user_arn,omitempty"`
	Profile       string `json:"profile,omitempty" yaml:"profile,omitempty"`
	Region        string `json:"region,omitempty" yaml:"region,omitempty"`
}

// CallerIdentity represents AWS STS get-caller-identity response.
type CallerIdentity struct {
	Account string `json:"Account"`
	Arn     string `json:"Arn"`
	UserID  string `json:"UserId"`
}

// Overview contains AWS resources summary.
type Overview struct {
	Status      *Status  `json:"status" yaml:"status"`
	EC2Count    int      `json:"ec2_count,omitempty" yaml:"ec2_count,omitempty"`
	S3Buckets   []string `json:"s3_buckets,omitempty" yaml:"s3_buckets,omitempty"`
	Lambdas     []string `json:"lambdas,omitempty" yaml:"lambdas,omitempty"`
	EKSClusters []string `json:"eks_clusters,omitempty" yaml:"eks_clusters,omitempty"`
}

// Helper provides AWS CLI helper operations.
type Helper struct {
	verbose bool
	dryRun  bool
	profile string
	region  string
}

// NewHelper creates a new AWS Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	return &Helper{
		verbose: verbose,
		dryRun:  dryRun,
		profile: os.Getenv("AWS_PROFILE"),
		region:  os.Getenv("AWS_DEFAULT_REGION"),
	}
}

// SetProfile sets the AWS profile.
func (h *Helper) SetProfile(profile string) {
	h.profile = profile
}

// SetRegion sets the AWS region.
func (h *Helper) SetRegion(region string) {
	h.region = region
}

// buildArgs builds command arguments with profile and region.
func (h *Helper) buildArgs(args ...string) []string {
	var cmdArgs []string
	if h.profile != "" {
		cmdArgs = append(cmdArgs, "--profile", h.profile)
	}
	if h.region != "" {
		cmdArgs = append(cmdArgs, "--region", h.region)
	}
	cmdArgs = append(cmdArgs, args...)
	return cmdArgs
}

// GetStatus returns AWS CLI status and authentication info.
func (h *Helper) GetStatus() (*Status, error) {
	status := &Status{
		Profile: h.profile,
		Region:  h.region,
	}

	// Check if aws is installed
	versionCmd := exec.Command("aws", "--version")
	versionOut, err := versionCmd.Output()
	if err != nil {
		status.Installed = false
		return status, nil
	}

	status.Installed = true
	// Parse version (format: aws-cli/2.x.x Python/3.x.x ...)
	versionStr := strings.TrimSpace(string(versionOut))
	if parts := strings.Fields(versionStr); len(parts) > 0 {
		status.Version = parts[0]
	}

	// Check authentication
	identity, err := h.GetCallerIdentity()
	if err != nil {
		status.Authenticated = false
		return status, nil
	}

	status.Authenticated = true
	status.AccountID = identity.Account
	status.UserARN = identity.Arn

	// Get current region if not set
	if status.Region == "" {
		if region := h.getConfiguredRegion(); region != "" {
			status.Region = region
		}
	}

	return status, nil
}

// GetCallerIdentity returns the current AWS identity.
func (h *Helper) GetCallerIdentity() (*CallerIdentity, error) {
	args := h.buildArgs("sts", "get-caller-identity", "--output", "json")
	cmd := exec.Command("aws", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get caller identity: %w", err)
	}

	var identity CallerIdentity
	if err := json.Unmarshal(out, &identity); err != nil {
		return nil, fmt.Errorf("failed to parse identity: %w", err)
	}

	return &identity, nil
}

// getConfiguredRegion gets the configured region from aws config.
func (h *Helper) getConfiguredRegion() string {
	args := h.buildArgs("configure", "get", "region")
	cmd := exec.Command("aws", args...)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// ListProfiles lists all configured AWS profiles.
func (h *Helper) ListProfiles() ([]string, error) {
	cmd := exec.Command("aws", "configure", "list-profiles")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list profiles: %w", err)
	}

	var profiles []string
	for line := range strings.SplitSeq(strings.TrimSpace(string(out)), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			profiles = append(profiles, line)
		}
	}
	return profiles, nil
}

// ListRegions lists available AWS regions.
func (h *Helper) ListRegions() ([]string, error) {
	args := h.buildArgs("ec2", "describe-regions", "--query", "Regions[].RegionName", "--output", "json")
	cmd := exec.Command("aws", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list regions: %w", err)
	}

	var regions []string
	if err := json.Unmarshal(out, &regions); err != nil {
		return nil, fmt.Errorf("failed to parse regions: %w", err)
	}
	return regions, nil
}

// ListEC2Instances lists EC2 instances.
func (h *Helper) ListEC2Instances() (string, error) {
	args := h.buildArgs("ec2", "describe-instances",
		"--query", "Reservations[].Instances[].[InstanceId,State.Name,InstanceType,PrivateIpAddress,Tags[?Key==`Name`].Value|[0]]",
		"--output", "table")
	cmd := exec.Command("aws", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if len(out) > 0 {
			return "", fmt.Errorf("%s", strings.TrimSpace(string(out)))
		}
		return "", fmt.Errorf("failed to list EC2 instances: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// ListS3Buckets lists S3 buckets.
func (h *Helper) ListS3Buckets() ([]string, error) {
	args := h.buildArgs("s3api", "list-buckets", "--query", "Buckets[].Name", "--output", "json")
	cmd := exec.Command("aws", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list S3 buckets: %w", err)
	}

	var buckets []string
	if err := json.Unmarshal(out, &buckets); err != nil {
		return nil, fmt.Errorf("failed to parse buckets: %w", err)
	}
	return buckets, nil
}

// ListLambdaFunctions lists Lambda functions.
func (h *Helper) ListLambdaFunctions() ([]string, error) {
	args := h.buildArgs("lambda", "list-functions", "--query", "Functions[].FunctionName", "--output", "json")
	cmd := exec.Command("aws", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list Lambda functions: %w", err)
	}

	var functions []string
	if err := json.Unmarshal(out, &functions); err != nil {
		return nil, fmt.Errorf("failed to parse functions: %w", err)
	}
	return functions, nil
}

// ListEKSClusters lists EKS clusters.
func (h *Helper) ListEKSClusters() ([]string, error) {
	args := h.buildArgs("eks", "list-clusters", "--query", "clusters", "--output", "json")
	cmd := exec.Command("aws", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list EKS clusters: %w", err)
	}

	var clusters []string
	if err := json.Unmarshal(out, &clusters); err != nil {
		return nil, fmt.Errorf("failed to parse clusters: %w", err)
	}
	return clusters, nil
}

// UpdateKubeconfig updates kubeconfig for an EKS cluster.
func (h *Helper) UpdateKubeconfig(clusterName string) error {
	if clusterName == "" {
		return fmt.Errorf("cluster name is required")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: aws eks update-kubeconfig --name %s\n", clusterName)
		return nil
	}

	args := h.buildArgs("eks", "update-kubeconfig", "--name", clusterName)
	cmd := exec.Command("aws", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// SSMConnect starts an SSM session to an instance.
func (h *Helper) SSMConnect(instanceID string) error {
	if instanceID == "" {
		return fmt.Errorf("instance ID is required")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: aws ssm start-session --target %s\n", instanceID)
		return nil
	}

	args := h.buildArgs("ssm", "start-session", "--target", instanceID)
	cmd := exec.Command("aws", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// GetOverview returns an overview of AWS resources.
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
	if buckets, err := h.ListS3Buckets(); err == nil {
		overview.S3Buckets = buckets
	}

	if lambdas, err := h.ListLambdaFunctions(); err == nil {
		overview.Lambdas = lambdas
	}

	if clusters, err := h.ListEKSClusters(); err == nil {
		overview.EKSClusters = clusters
	}

	// Count EC2 instances
	if count, err := h.countEC2Instances(); err == nil {
		overview.EC2Count = count
	}

	return overview, nil
}

// countEC2Instances counts running EC2 instances.
func (h *Helper) countEC2Instances() (int, error) {
	args := h.buildArgs("ec2", "describe-instances",
		"--filters", "Name=instance-state-name,Values=running",
		"--query", "length(Reservations[].Instances[])",
		"--output", "json")
	cmd := exec.Command("aws", args...)
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

// Configure runs aws configure.
func (h *Helper) Configure() error {
	if h.dryRun {
		fmt.Println("[dry-run] would run: aws configure")
		return nil
	}

	cmd := exec.Command("aws", "configure")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// SSOLogin performs SSO login.
func (h *Helper) SSOLogin(profile string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would run: aws sso login --profile %s\n", profile)
		return nil
	}

	args := []string{"sso", "login"}
	if profile != "" {
		args = append(args, "--profile", profile)
	}

	cmd := exec.Command("aws", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

