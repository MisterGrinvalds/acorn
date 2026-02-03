// Package neomutt provides NeoMutt terminal email client helper functionality.
package neomutt

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Status represents NeoMutt installation status.
type Status struct {
	Installed    bool   `json:"installed" yaml:"installed"`
	Version      string `json:"version,omitempty" yaml:"version,omitempty"`
	ConfigDir    string `json:"config_dir" yaml:"config_dir"`
	ConfigExists bool   `json:"config_exists" yaml:"config_exists"`
	CacheDir     string `json:"cache_dir" yaml:"cache_dir"`
	AccountCount int    `json:"account_count" yaml:"account_count"`
}

// Account represents an email account configuration.
type Account struct {
	Name     string `json:"name" yaml:"name"`
	Type     string `json:"type" yaml:"type"`         // gmail, microsoft, imap
	Email    string `json:"email,omitempty" yaml:"email,omitempty"`
	File     string `json:"file" yaml:"file"`
	HasToken bool   `json:"has_token" yaml:"has_token"`
}

// TokenInfo represents OAuth2 token status.
type TokenInfo struct {
	Account   string `json:"account" yaml:"account"`
	TokenFile string `json:"token_file" yaml:"token_file"`
	Exists    bool   `json:"exists" yaml:"exists"`
	Encrypted bool   `json:"encrypted" yaml:"encrypted"`
}

// CacheInfo represents cache directory information.
type CacheInfo struct {
	HeaderCache  string `json:"header_cache" yaml:"header_cache"`
	MessageCache string `json:"message_cache" yaml:"message_cache"`
	HeaderSize   string `json:"header_size,omitempty" yaml:"header_size,omitempty"`
	MessageSize  string `json:"message_size,omitempty" yaml:"message_size,omitempty"`
}

// Helper provides NeoMutt helper operations.
type Helper struct {
	verbose   bool
	dryRun    bool
	configDir string
	cacheDir  string
}

// NewHelper creates a new NeoMutt Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		home, _ := os.UserHomeDir()
		configDir = filepath.Join(home, ".config")
	}
	configDir = filepath.Join(configDir, "neomutt")

	cacheDir := os.Getenv("XDG_CACHE_HOME")
	if cacheDir == "" {
		home, _ := os.UserHomeDir()
		cacheDir = filepath.Join(home, ".cache")
	}
	cacheDir = filepath.Join(cacheDir, "neomutt")

	return &Helper{
		verbose:   verbose,
		dryRun:    dryRun,
		configDir: configDir,
		cacheDir:  cacheDir,
	}
}

// GetConfigDir returns the NeoMutt config directory.
func (h *Helper) GetConfigDir() string {
	return h.configDir
}

// GetCacheDir returns the NeoMutt cache directory.
func (h *Helper) GetCacheDir() string {
	return h.cacheDir
}

// GetStatus returns NeoMutt status information.
func (h *Helper) GetStatus() *Status {
	status := &Status{
		ConfigDir: h.configDir,
		CacheDir:  h.cacheDir,
	}

	// Check if neomutt is installed
	out, err := exec.Command("neomutt", "-v").Output()
	if err != nil {
		status.Installed = false
		return status
	}

	status.Installed = true

	// Parse version from first line (NeoMutt 20231103)
	lines := strings.Split(string(out), "\n")
	if len(lines) > 0 {
		parts := strings.Fields(lines[0])
		if len(parts) >= 2 {
			status.Version = parts[1]
		}
	}

	// Check if config exists
	mainConfig := filepath.Join(h.configDir, "neomuttrc")
	if _, err := os.Stat(mainConfig); err == nil {
		status.ConfigExists = true
	}

	// Count accounts
	accounts, _ := h.ListAccounts()
	status.AccountCount = len(accounts)

	return status
}

// ListAccounts lists configured email accounts.
func (h *Helper) ListAccounts() ([]Account, error) {
	accountsDir := filepath.Join(h.configDir, "accounts")
	entries, err := os.ReadDir(accountsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []Account{}, nil
		}
		return nil, fmt.Errorf("failed to read accounts directory: %w", err)
	}

	var accounts []Account
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".muttrc") {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".muttrc")
		account := Account{
			Name: name,
			File: filepath.Join(accountsDir, entry.Name()),
		}

		// Determine account type from filename
		if strings.HasPrefix(name, "gmail-") {
			account.Type = "gmail"
		} else if strings.HasPrefix(name, "microsoft-") {
			account.Type = "microsoft"
		} else {
			account.Type = "imap"
		}
		account.Email = extractEmailFromFile(account.File)

		// Check for OAuth token file
		tokenFile := h.getTokenFile(account.Name)
		if _, err := os.Stat(tokenFile); err == nil {
			account.HasToken = true
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

// extractEmailFromFile extracts the email address from an account config file.
func extractEmailFromFile(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}

	for line := range strings.SplitSeq(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "set from") || strings.HasPrefix(line, "set imap_user") {
			// Extract email from line like: set from = "email@example.com"
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				email := strings.TrimSpace(parts[1])
				email = strings.Trim(email, "\"' ")
				return email
			}
		}
	}
	return ""
}

// getTokenFile returns the token file path for an account.
// Parses the account config file to find the token file path from imap_oauth_refresh_command.
func (h *Helper) getTokenFile(accountName string) string {
	accountFile := filepath.Join(h.configDir, "accounts", accountName+".muttrc")
	return extractTokenFileFromConfig(accountFile)
}

// extractTokenFileFromConfig extracts the token file path from an account config.
func extractTokenFileFromConfig(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}

	for line := range strings.SplitSeq(string(data), "\n") {
		line = strings.TrimSpace(line)
		// Look for: set imap_oauth_refresh_command = "python3 ~/.config/neomutt/mutt_oauth2.py ~/.config/neomutt/gmail.ross.bercot.tokens"
		if strings.Contains(line, "imap_oauth_refresh_command") {
			// Extract the token file path (last argument)
			for part := range strings.FieldsSeq(line) {
				if strings.HasSuffix(part, ".tokens") || strings.HasSuffix(part, ".tokens\"") {
					// Clean up the path
					tokenFile := strings.Trim(part, "\"'")
					// Expand ~ to home directory
					if strings.HasPrefix(tokenFile, "~") {
						home, _ := os.UserHomeDir()
						tokenFile = strings.Replace(tokenFile, "~", home, 1)
					}
					return tokenFile
				}
			}
		}
	}
	return ""
}

// GetTokenStatus returns OAuth2 token status for accounts.
func (h *Helper) GetTokenStatus() ([]TokenInfo, error) {
	accounts, err := h.ListAccounts()
	if err != nil {
		return nil, err
	}

	var tokens []TokenInfo
	for _, account := range accounts {
		if account.Type != "gmail" && account.Type != "microsoft" {
			continue // Only OAuth accounts
		}

		tokenFile := h.getTokenFile(account.Name)
		info := TokenInfo{
			Account:   account.Name,
			TokenFile: tokenFile,
		}

		if _, err := os.Stat(tokenFile); err == nil {
			info.Exists = true
			// Check if file is GPG encrypted (has binary content or .gpg extension)
			data, _ := os.ReadFile(tokenFile)
			if len(data) > 0 && !strings.HasPrefix(string(data), "{") {
				info.Encrypted = true
			}
		}

		tokens = append(tokens, info)
	}

	return tokens, nil
}

// GetCacheInfo returns cache directory information.
func (h *Helper) GetCacheInfo() *CacheInfo {
	info := &CacheInfo{
		HeaderCache:  filepath.Join(h.cacheDir, "headers"),
		MessageCache: filepath.Join(h.cacheDir, "bodies"),
	}

	info.HeaderSize = h.getDirSize(info.HeaderCache)
	info.MessageSize = h.getDirSize(info.MessageCache)

	return info
}

// getDirSize returns the size of a directory.
func (h *Helper) getDirSize(dir string) string {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return "0"
	}

	cmd := exec.Command("du", "-sh", dir)
	out, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	parts := strings.Fields(string(out))
	if len(parts) > 0 {
		return parts[0]
	}
	return "unknown"
}

// CleanCache cleans the NeoMutt cache.
func (h *Helper) CleanCache() error {
	if h.dryRun {
		fmt.Printf("[dry-run] would remove: %s\n", h.cacheDir)
		return nil
	}

	// Remove header cache
	headerCache := filepath.Join(h.cacheDir, "headers")
	if err := os.RemoveAll(headerCache); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to clean header cache: %w", err)
	}

	// Remove message cache
	messageCache := filepath.Join(h.cacheDir, "bodies")
	if err := os.RemoveAll(messageCache); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to clean message cache: %w", err)
	}

	// Recreate directories
	if err := os.MkdirAll(headerCache, 0755); err != nil {
		return fmt.Errorf("failed to recreate header cache: %w", err)
	}
	if err := os.MkdirAll(messageCache, 0755); err != nil {
		return fmt.Errorf("failed to recreate message cache: %w", err)
	}

	return nil
}

// Launch starts NeoMutt.
func (h *Helper) Launch(args ...string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would run: neomutt %s\n", strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command("neomutt", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// RefreshToken refreshes OAuth2 token for an account.
func (h *Helper) RefreshToken(accountName string) error {
	accounts, err := h.ListAccounts()
	if err != nil {
		return err
	}

	var account *Account
	for _, a := range accounts {
		if a.Name == accountName {
			account = &a
			break
		}
	}

	if account == nil {
		return fmt.Errorf("account not found: %s", accountName)
	}

	if account.Type != "gmail" && account.Type != "microsoft" {
		return fmt.Errorf("account %s does not use OAuth2", accountName)
	}

	tokenFile := h.getTokenFile(account.Name)
	oauthScript := filepath.Join(h.configDir, "mutt_oauth2.py")

	if _, err := os.Stat(oauthScript); os.IsNotExist(err) {
		return fmt.Errorf("OAuth2 script not found: %s", oauthScript)
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: python3 %s %s\n", oauthScript, tokenFile)
		return nil
	}

	cmd := exec.Command("python3", oauthScript, tokenFile, "--verbose")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// AuthorizeAccount initiates OAuth2 authorization for an account.
func (h *Helper) AuthorizeAccount(accountName string) error {
	accounts, err := h.ListAccounts()
	if err != nil {
		return err
	}

	var account *Account
	for _, a := range accounts {
		if a.Name == accountName {
			account = &a
			break
		}
	}

	if account == nil {
		return fmt.Errorf("account not found: %s", accountName)
	}

	if account.Type != "gmail" && account.Type != "microsoft" {
		return fmt.Errorf("account %s does not use OAuth2", accountName)
	}

	tokenFile := h.getTokenFile(account.Name)
	oauthScript := filepath.Join(h.configDir, "mutt_oauth2.py")

	if _, err := os.Stat(oauthScript); os.IsNotExist(err) {
		return fmt.Errorf("OAuth2 script not found: %s", oauthScript)
	}

	if h.dryRun {
		fmt.Printf("[dry-run] would run: python3 %s %s --authorize\n", oauthScript, tokenFile)
		return nil
	}

	cmd := exec.Command("python3", oauthScript, tokenFile, "--verbose", "--authorize")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Install installs NeoMutt.
func (h *Helper) Install() error {
	if _, err := exec.LookPath("neomutt"); err == nil {
		return fmt.Errorf("neomutt is already installed")
	}

	if h.dryRun {
		fmt.Println("[dry-run] would install neomutt via homebrew")
		return nil
	}

	fmt.Println("Installing NeoMutt...")
	cmd := exec.Command("brew", "install", "neomutt")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// InitConfig initializes NeoMutt configuration directories.
func (h *Helper) InitConfig() error {
	dirs := []string{
		h.configDir,
		filepath.Join(h.configDir, "accounts"),
		filepath.Join(h.configDir, "signatures"),
		filepath.Join(h.cacheDir, "headers"),
		filepath.Join(h.cacheDir, "bodies"),
	}

	for _, dir := range dirs {
		if h.dryRun {
			fmt.Printf("[dry-run] would create: %s\n", dir)
			continue
		}
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create %s: %w", dir, err)
		}
	}

	return nil
}

// GenerateMainConfig generates the main neomuttrc configuration file.
func (h *Helper) GenerateMainConfig(dryRun bool) error {
	configFile := filepath.Join(h.configDir, "neomuttrc")

	// Check if config already exists
	if _, err := os.Stat(configFile); err == nil && !dryRun {
		fmt.Printf("Config exists: %s (skipping, use --force to overwrite)\n", configFile)
		return nil
	}

	// Generate config content
	content := h.generateDefaultConfig()

	if dryRun {
		fmt.Printf("Would create: %s\n", configFile)
		fmt.Println("\n--- Preview ---")
		fmt.Println(string(content))
		fmt.Println("--- End Preview ---")
		return nil
	}

	// Write the file
	if err := os.WriteFile(configFile, content, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	fmt.Printf("Created: %s\n", configFile)
	return nil
}

// generateDefaultConfig generates default neomuttrc content.
func (h *Helper) generateDefaultConfig() []byte {
	var b strings.Builder

	b.WriteString("# NeoMutt Configuration\n")
	b.WriteString("# Generated by acorn - customize as needed\n\n")

	b.WriteString("# ===== Basic Settings =====\n")
	b.WriteString("set editor = \"vim\"\n")
	b.WriteString("set mail_check = 60\n")
	b.WriteString("set timeout = 10\n")
	b.WriteString("set sort = reverse-date\n")
	b.WriteString("set ssl_force_tls = yes\n\n")

	b.WriteString("# Caching (shared across accounts)\n")
	fmt.Fprintf(&b, "set header_cache = \"%s/headers\"\n", h.cacheDir)
	fmt.Fprintf(&b, "set message_cachedir = \"%s/bodies\"\n\n", h.cacheDir)

	b.WriteString("# Sidebar\n")
	b.WriteString("set sidebar_visible = yes\n")
	b.WriteString("set sidebar_width = 25\n")
	b.WriteString("set sidebar_format = \"%B%?F? [%F]?%* %?N?%N/?%S\"\n")
	b.WriteString("set sidebar_short_path = yes\n\n")

	b.WriteString("# ===== Load Default Account =====\n")
	b.WriteString("# Uncomment and modify after adding accounts:\n")
	fmt.Fprintf(&b, "# source %s/accounts/default.muttrc\n\n", h.configDir)

	b.WriteString("# ===== Keybindings =====\n")
	b.WriteString("# Sidebar navigation\n")
	b.WriteString("bind index,pager \\Ck sidebar-prev\n")
	b.WriteString("bind index,pager \\Cj sidebar-next\n")
	b.WriteString("bind index,pager \\Co sidebar-open\n")
	b.WriteString("bind index,pager B sidebar-toggle-visible\n\n")

	b.WriteString("# Vim-like navigation\n")
	b.WriteString("bind index g noop\n")
	b.WriteString("bind index gg first-entry\n")
	b.WriteString("bind index G last-entry\n")
	b.WriteString("bind pager g noop\n")
	b.WriteString("bind pager gg top\n")
	b.WriteString("bind pager G bottom\n")
	b.WriteString("bind index,pager \\Cu half-up\n")
	b.WriteString("bind index,pager \\Cd half-down\n\n")

	b.WriteString("# ===== Colors =====\n")
	b.WriteString("color sidebar_new yellow default\n")
	b.WriteString("color indicator black cyan\n")

	return []byte(b.String())
}

// AddGmailAccount adds a Gmail account configuration.
func (h *Helper) AddGmailAccount(email, realName string, dryRun bool) (string, error) {
	// Create account name from email
	username := strings.Split(email, "@")[0]
	accountName := "gmail-" + username
	accountFile := filepath.Join(h.configDir, "accounts", accountName+".muttrc")
	tokenFile := filepath.Join(h.configDir, "gmail."+username+".tokens")

	// Check if account already exists
	if _, err := os.Stat(accountFile); err == nil && !dryRun {
		return "", fmt.Errorf("account already exists: %s", accountFile)
	}

	// Generate account config using the writer
	writer := NewAccountWriter()
	content := writer.WriteGmailAccount(email, realName, tokenFile)

	if dryRun {
		fmt.Printf("Would create: %s\n", accountFile)
		fmt.Println("\n--- Preview ---")
		fmt.Println(string(content))
		fmt.Println("--- End Preview ---")
		return accountFile, nil
	}

	// Ensure accounts directory exists
	if err := os.MkdirAll(filepath.Dir(accountFile), 0755); err != nil {
		return "", fmt.Errorf("failed to create accounts directory: %w", err)
	}

	// Write account file
	if err := os.WriteFile(accountFile, content, 0600); err != nil {
		return "", fmt.Errorf("failed to write account config: %w", err)
	}

	// Create empty signature file
	sigFile := filepath.Join(h.configDir, "signatures", accountName)
	if _, err := os.Stat(sigFile); os.IsNotExist(err) {
		sigContent := fmt.Sprintf("--\n%s\n%s\n", realName, email)
		_ = os.WriteFile(sigFile, []byte(sigContent), 0644)
	}

	return accountFile, nil
}

// AddMicrosoftAccount adds a Microsoft/Office365 account configuration.
func (h *Helper) AddMicrosoftAccount(email, realName string, dryRun bool) (string, error) {
	// Create account name from email (domain-user format)
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid email format: %s", email)
	}
	accountName := "microsoft-" + parts[1] + "-" + parts[0]
	accountFile := filepath.Join(h.configDir, "accounts", accountName+".muttrc")
	tokenFile := filepath.Join(h.configDir, "microsoft."+email+".tokens")

	// Check if account already exists
	if _, err := os.Stat(accountFile); err == nil && !dryRun {
		return "", fmt.Errorf("account already exists: %s", accountFile)
	}

	// Generate account config using the writer
	writer := NewAccountWriter()
	content := writer.WriteMicrosoftAccount(email, realName, tokenFile)

	if dryRun {
		fmt.Printf("Would create: %s\n", accountFile)
		fmt.Println("\n--- Preview ---")
		fmt.Println(string(content))
		fmt.Println("--- End Preview ---")
		return accountFile, nil
	}

	// Ensure accounts directory exists
	if err := os.MkdirAll(filepath.Dir(accountFile), 0755); err != nil {
		return "", fmt.Errorf("failed to create accounts directory: %w", err)
	}

	// Write account file
	if err := os.WriteFile(accountFile, content, 0600); err != nil {
		return "", fmt.Errorf("failed to write account config: %w", err)
	}

	// Create empty signature file
	sigFile := filepath.Join(h.configDir, "signatures", accountName)
	if _, err := os.Stat(sigFile); os.IsNotExist(err) {
		sigContent := fmt.Sprintf("--\n%s\n%s\n", realName, email)
		_ = os.WriteFile(sigFile, []byte(sigContent), 0644)
	}

	return accountFile, nil
}
