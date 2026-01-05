// Package secrets provides secrets management and credential checking.
package secrets

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Status represents the status of the secrets file.
type Status struct {
	FilePath  string `json:"file_path" yaml:"file_path"`
	Exists    bool   `json:"exists" yaml:"exists"`
	Readable  bool   `json:"readable" yaml:"readable"`
	KeyCount  int    `json:"key_count" yaml:"key_count"`
}

// Credential represents a credential check result.
type Credential struct {
	Name      string `json:"name" yaml:"name"`
	Available bool   `json:"available" yaml:"available"`
	EnvVars   []string `json:"env_vars,omitempty" yaml:"env_vars,omitempty"`
}

// CredentialCheck contains results of checking all credentials.
type CredentialCheck struct {
	Credentials []Credential `json:"credentials" yaml:"credentials"`
	Available   int          `json:"available" yaml:"available"`
	Missing     int          `json:"missing" yaml:"missing"`
}

// Helper provides secrets management operations.
type Helper struct {
	secretsDir string
	verbose    bool
}

// NewHelper creates a new secrets Helper.
func NewHelper(verbose bool) *Helper {
	secretsDir := os.Getenv("SECRETS_DIR")
	if secretsDir == "" {
		xdgData := os.Getenv("XDG_DATA_HOME")
		if xdgData == "" {
			home, _ := os.UserHomeDir()
			xdgData = filepath.Join(home, ".local", "share")
		}
		secretsDir = filepath.Join(xdgData, "secrets")
	}

	return &Helper{
		secretsDir: secretsDir,
		verbose:    verbose,
	}
}

// GetSecretsFile returns the path to the secrets file.
func (h *Helper) GetSecretsFile() string {
	return filepath.Join(h.secretsDir, ".env")
}

// GetStatus returns the status of the secrets file.
func (h *Helper) GetStatus() (*Status, error) {
	secretsFile := h.GetSecretsFile()
	status := &Status{
		FilePath: secretsFile,
	}

	info, err := os.Stat(secretsFile)
	if os.IsNotExist(err) {
		return status, nil
	}
	if err != nil {
		return nil, err
	}

	status.Exists = true

	// Check if readable
	file, err := os.Open(secretsFile)
	if err != nil {
		return status, nil
	}
	defer file.Close()
	status.Readable = true

	// Count keys
	keyCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 && line[0] >= 'A' && line[0] <= 'Z' && strings.Contains(line, "=") {
			keyCount++
		}
	}
	status.KeyCount = keyCount

	// Ensure permissions are secure (readable only by owner)
	if info.Mode().Perm()&0o077 != 0 {
		if h.verbose {
			fmt.Printf("Warning: secrets file has insecure permissions: %v\n", info.Mode().Perm())
		}
	}

	return status, nil
}

// ListSecrets returns a list of secret keys (not values).
func (h *Helper) ListSecrets() ([]string, error) {
	secretsFile := h.GetSecretsFile()

	file, err := os.Open(secretsFile)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("no secrets file found at: %s", secretsFile)
	}
	if err != nil {
		return nil, fmt.Errorf("cannot read secrets file: %w", err)
	}
	defer file.Close()

	var keys []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 && line[0] >= 'A' && line[0] <= 'Z' {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) >= 1 {
				keys = append(keys, parts[0])
			}
		}
	}

	sort.Strings(keys)
	return keys, scanner.Err()
}

// LoadSecrets loads secrets from the file into the current process environment.
// Note: This only affects the current process, not the parent shell.
func (h *Helper) LoadSecrets() (int, error) {
	secretsFile := h.GetSecretsFile()

	file, err := os.Open(secretsFile)
	if os.IsNotExist(err) {
		return 0, fmt.Errorf("no secrets file found at: %s", secretsFile)
	}
	if err != nil {
		return 0, fmt.Errorf("cannot read secrets file (check permissions): %w", err)
	}
	defer file.Close()

	loaded := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		// Parse KEY=value
		if idx := strings.Index(line, "="); idx > 0 {
			key := line[:idx]
			value := line[idx+1:]
			// Remove quotes if present
			value = strings.Trim(value, "\"'")
			if err := os.Setenv(key, value); err == nil {
				loaded++
				if h.verbose {
					fmt.Printf("Loaded: %s\n", key)
				}
			}
		}
	}

	return loaded, scanner.Err()
}

// CheckCredential checks if a specific credential is available.
func (h *Helper) CheckCredential(name string, envVars ...string) *Credential {
	cred := &Credential{
		Name:    name,
		EnvVars: envVars,
	}

	// All env vars must be set for the credential to be available
	allSet := true
	for _, envVar := range envVars {
		if os.Getenv(envVar) == "" {
			allSet = false
			break
		}
	}
	cred.Available = allSet

	return cred
}

// CheckAWS checks AWS credentials.
func (h *Helper) CheckAWS() *Credential {
	return h.CheckCredential("AWS", "AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY")
}

// CheckAzure checks Azure credentials.
func (h *Helper) CheckAzure() *Credential {
	return h.CheckCredential("Azure", "AZURE_CLIENT_ID", "AZURE_CLIENT_SECRET", "AZURE_TENANT_ID")
}

// CheckGitHub checks GitHub token.
func (h *Helper) CheckGitHub() *Credential {
	return h.CheckCredential("GitHub", "GITHUB_TOKEN")
}

// CheckDigitalOcean checks DigitalOcean token.
func (h *Helper) CheckDigitalOcean() *Credential {
	return h.CheckCredential("DigitalOcean", "DIGITALOCEAN_TOKEN")
}

// CheckOpenAI checks OpenAI API key.
func (h *Helper) CheckOpenAI() *Credential {
	return h.CheckCredential("OpenAI", "OPENAI_API_KEY")
}

// CheckAnthropic checks Anthropic API key.
func (h *Helper) CheckAnthropic() *Credential {
	return h.CheckCredential("Anthropic", "ANTHROPIC_API_KEY")
}

// CheckHuggingFace checks Hugging Face token.
func (h *Helper) CheckHuggingFace() *Credential {
	return h.CheckCredential("HuggingFace", "HUGGINGFACE_TOKEN", "HF_TOKEN")
}

// CheckAllCredentials checks all known credentials.
func (h *Helper) CheckAllCredentials() *CredentialCheck {
	check := &CredentialCheck{}

	credentials := []*Credential{
		h.CheckAWS(),
		h.CheckAzure(),
		h.CheckGitHub(),
		h.CheckDigitalOcean(),
		h.CheckOpenAI(),
		h.CheckAnthropic(),
		h.CheckHuggingFace(),
	}

	for _, cred := range credentials {
		check.Credentials = append(check.Credentials, *cred)
		if cred.Available {
			check.Available++
		} else {
			check.Missing++
		}
	}

	return check
}

// ValidateSecrets validates that common credentials are configured.
func (h *Helper) ValidateSecrets() *CredentialCheck {
	return h.CheckAllCredentials()
}

// EnsureSecretsDir creates the secrets directory if it doesn't exist.
func (h *Helper) EnsureSecretsDir() error {
	return os.MkdirAll(h.secretsDir, 0o700)
}

// CreateSecretsFile creates an empty secrets file with secure permissions.
func (h *Helper) CreateSecretsFile() error {
	if err := h.EnsureSecretsDir(); err != nil {
		return err
	}

	secretsFile := h.GetSecretsFile()
	if _, err := os.Stat(secretsFile); err == nil {
		return fmt.Errorf("secrets file already exists: %s", secretsFile)
	}

	template := `# Secrets file - loaded by dotfiles
# Keep this file secure with permissions 0600

# Cloud Providers
# AWS_ACCESS_KEY_ID=
# AWS_SECRET_ACCESS_KEY=
# AZURE_CLIENT_ID=
# AZURE_CLIENT_SECRET=
# AZURE_TENANT_ID=
# DIGITALOCEAN_TOKEN=

# Developer Tools
# GITHUB_TOKEN=

# AI Services
# OPENAI_API_KEY=
# ANTHROPIC_API_KEY=
# HUGGINGFACE_TOKEN=
`

	if err := os.WriteFile(secretsFile, []byte(template), 0o600); err != nil {
		return err
	}

	return nil
}
