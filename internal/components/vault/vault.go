// Package vault provides HashiCorp Vault CLI helper functionality.
package vault

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Status represents Vault installation and connection status.
type Status struct {
	Installed       bool   `json:"installed" yaml:"installed"`
	Version         string `json:"version,omitempty" yaml:"version,omitempty"`
	Location        string `json:"location,omitempty" yaml:"location,omitempty"`
	ServerAddress   string `json:"server_address,omitempty" yaml:"server_address,omitempty"`
	Connected       bool   `json:"connected" yaml:"connected"`
	Authenticated   bool   `json:"authenticated" yaml:"authenticated"`
	Sealed          bool   `json:"sealed,omitempty" yaml:"sealed,omitempty"`
	TokenExpiry     string `json:"token_expiry,omitempty" yaml:"token_expiry,omitempty"`
}

// Secret represents a vault secret.
type Secret struct {
	Path string            `json:"path" yaml:"path"`
	Data map[string]string `json:"data" yaml:"data"`
}

// Helper provides Vault helper operations.
type Helper struct {
	verbose bool
	dryRun  bool
}

// NewHelper creates a new Vault Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	return &Helper{
		verbose: verbose,
		dryRun:  dryRun,
	}
}

// IsInstalled checks if vault is installed.
func (h *Helper) IsInstalled() bool {
	_, err := exec.LookPath("vault")
	return err == nil
}

// GetStatus returns Vault installation and connection status.
func (h *Helper) GetStatus() *Status {
	status := &Status{
		ServerAddress: os.Getenv("VAULT_ADDR"),
	}

	// Check if vault CLI is installed
	vaultPath, err := exec.LookPath("vault")
	if err != nil {
		return status
	}

	status.Installed = true
	status.Location = vaultPath

	// Get version
	cmd := exec.Command("vault", "version")
	out, err := cmd.Output()
	if err == nil {
		version := strings.TrimSpace(string(out))
		// Extract version number (format: "Vault v1.15.0 (...")
		if strings.HasPrefix(version, "Vault v") {
			parts := strings.Fields(version)
			if len(parts) >= 2 {
				status.Version = strings.TrimPrefix(parts[1], "v")
			}
		} else {
			status.Version = version
		}
	}

	// Check server connection and status
	if status.ServerAddress != "" {
		cmd = exec.Command("vault", "status", "-format=json")
		out, err = cmd.Output()
		if err == nil {
			status.Connected = true
			// Parse JSON to get seal status
			if strings.Contains(string(out), `"sealed":true`) {
				status.Sealed = true
			}
		}

		// Check authentication
		cmd = exec.Command("vault", "token", "lookup", "-format=json")
		out, err = cmd.Output()
		if err == nil {
			status.Authenticated = true
			// Could parse expire_time from JSON here if needed
		}
	}

	return status
}

// Login authenticates to Vault with a token.
func (h *Helper) Login(token string) error {
	if token == "" {
		return fmt.Errorf("token is required")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would login to Vault with token\n")
		return nil
	}

	cmd := exec.Command("vault", "login", "-method=token", fmt.Sprintf("token=%s", token))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Read reads a secret from Vault.
func (h *Helper) Read(path string) (*Secret, error) {
	if path == "" {
		return nil, fmt.Errorf("path is required")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would read secret from: %s\n", path)
		return &Secret{Path: path, Data: map[string]string{}}, nil
	}

	cmd := exec.Command("vault", "kv", "get", "-format=json", path)
	_, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to read secret: %w", err)
	}

	// Simple parsing - in real implementation, would properly parse JSON
	secret := &Secret{
		Path: path,
		Data: make(map[string]string),
	}

	return secret, nil
}

// Write writes a secret to Vault.
func (h *Helper) Write(path string, data map[string]string) error {
	if path == "" {
		return fmt.Errorf("path is required")
	}
	if len(data) == 0 {
		return fmt.Errorf("data is required")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would write secret to: %s\n", path)
		return nil
	}

	args := []string{"kv", "put", path}
	for k, v := range data {
		args = append(args, fmt.Sprintf("%s=%s", k, v))
	}

	cmd := exec.Command("vault", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Delete deletes a secret from Vault.
func (h *Helper) Delete(path string) error {
	if path == "" {
		return fmt.Errorf("path is required")
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would delete secret at: %s\n", path)
		return nil
	}

	cmd := exec.Command("vault", "kv", "delete", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// List lists secrets at a path.
func (h *Helper) List(path string) ([]string, error) {
	if path == "" {
		path = "secret/"
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would list secrets at: %s\n", path)
		return []string{}, nil
	}

	cmd := exec.Command("vault", "kv", "list", "-format=json", path)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}

	// Simple parsing - in real implementation, would properly parse JSON
	var secrets []string
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "[") && !strings.HasPrefix(line, "]") {
			line = strings.Trim(line, `"`)
			line = strings.TrimSuffix(line, ",")
			if line != "" {
				secrets = append(secrets, line)
			}
		}
	}

	return secrets, nil
}

// RenewToken renews the current Vault token.
func (h *Helper) RenewToken(increment string) error {
	if increment == "" {
		increment = "1h"
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would renew token for: %s\n", increment)
		return nil
	}

	cmd := exec.Command("vault", "token", "renew", fmt.Sprintf("-increment=%s", increment))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// SetEnvironment sets the Vault environment (server address).
func (h *Helper) SetEnvironment(env string) (string, error) {
	var addr string
	switch strings.ToLower(env) {
	case "prod", "production":
		addr = "https://vault.prod.example.com:8200"
	case "staging", "stage":
		addr = "https://vault.staging.example.com:8200"
	case "dev", "development":
		addr = "http://localhost:8200"
	default:
		return "", fmt.Errorf("unknown environment: %s (use: prod, staging, dev)", env)
	}

	if h.verbose {
		fmt.Printf("Vault environment set to: %s (%s)\n", env, addr)
	}

	return addr, nil
}

// CheckConnection verifies connection to Vault server.
func (h *Helper) CheckConnection() error {
	cmd := exec.Command("vault", "status")
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("cannot connect to Vault server: %w", err)
	}
	return nil
}
