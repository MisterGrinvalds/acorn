package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mistergrinvalds/acorn/internal/components/mail/neomutt"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	neomuttDryRun  bool
	neomuttVerbose bool
)

// neomuttCmd represents the neomutt command group
var neomuttCmd = &cobra.Command{
	Use:   "neomutt",
	Short: "NeoMutt terminal email client",
	Long: `NeoMutt terminal email client commands.

NeoMutt is a command-line mail reader with vim-style keybindings.
Provides commands for status, accounts, tokens, and cache management.

Examples:
  acorn mail neomutt status        # Show NeoMutt status
  acorn mail neomutt accounts      # List email accounts
  acorn mail neomutt tokens        # Show OAuth2 token status
  acorn mail neomutt cache info    # Show cache info
  acorn mail neomutt launch        # Start NeoMutt`,
	Aliases: []string{"nm", "mutt"},
}

// neomuttStatusCmd shows NeoMutt status
var neomuttStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show NeoMutt status",
	Long: `Show NeoMutt installation status and configuration info.

Examples:
  acorn mail neomutt status
  acorn mail neomutt status -o json`,
	RunE: runNeomuttStatus,
}

// neomuttInstallCmd installs NeoMutt
var neomuttInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install NeoMutt",
	Long: `Install NeoMutt using Homebrew.

Examples:
  acorn mail neomutt install`,
	RunE: runNeomuttInstall,
}

// neomuttLaunchCmd launches NeoMutt
var neomuttLaunchCmd = &cobra.Command{
	Use:   "launch [args...]",
	Short: "Launch NeoMutt",
	Long: `Start NeoMutt terminal email client.

Examples:
  acorn mail neomutt launch
  acorn mail neomutt launch -f personal`,
	Aliases: []string{"start", "open"},
	RunE:    runNeomuttLaunch,
}

// neomuttInitCmd initializes config directories
var neomuttInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration directories",
	Long: `Create NeoMutt configuration and cache directories.

Creates:
  ~/.config/neomutt/
  ~/.config/neomutt/accounts/
  ~/.config/neomutt/signatures/
  ~/.cache/neomutt/headers/
  ~/.cache/neomutt/bodies/

Examples:
  acorn mail neomutt init`,
	RunE: runNeomuttInit,
}

// Account subcommands
var neomuttAccountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "List email accounts",
	Long: `List configured email accounts.

Shows account names, types (gmail, microsoft, imap), and OAuth2 status.

Examples:
  acorn mail neomutt accounts
  acorn mail neomutt accounts -o json`,
	Aliases: []string{"account", "ls"},
	RunE:    runNeomuttAccounts,
}

// Token subcommands
var neomuttTokensCmd = &cobra.Command{
	Use:   "tokens",
	Short: "OAuth2 token management",
	Long: `Manage OAuth2 tokens for email accounts.

Examples:
  acorn mail neomutt tokens          # Show token status
  acorn mail neomutt tokens refresh <account>`,
}

var neomuttTokensStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show OAuth2 token status",
	Long: `Show OAuth2 token status for all accounts.

Examples:
  acorn mail neomutt tokens status
  acorn mail neomutt tokens status -o json`,
	Aliases: []string{"ls", "list"},
	RunE:    runNeomuttTokensStatus,
}

var neomuttTokensRefreshCmd = &cobra.Command{
	Use:   "refresh <account>",
	Short: "Refresh OAuth2 token",
	Long: `Refresh OAuth2 token for an account.

Examples:
  acorn mail neomutt tokens refresh gmail-ross.bercot`,
	Args: cobra.ExactArgs(1),
	RunE: runNeomuttTokensRefresh,
}

var neomuttTokensAuthorizeCmd = &cobra.Command{
	Use:   "authorize <account>",
	Short: "Authorize OAuth2 for account",
	Long: `Initiate OAuth2 authorization for an account.

This will open a browser for authentication.

Examples:
  acorn mail neomutt tokens authorize gmail-ross.bercot`,
	Aliases: []string{"auth", "init"},
	Args:    cobra.ExactArgs(1),
	RunE:    runNeomuttTokensAuthorize,
}

// Cache subcommands
var neomuttCacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Cache management",
	Long: `Manage NeoMutt cache.

Examples:
  acorn mail neomutt cache info
  acorn mail neomutt cache clean`,
}

var neomuttCacheInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show cache information",
	Long: `Show NeoMutt cache directory and size.

Examples:
  acorn mail neomutt cache info`,
	RunE: runNeomuttCacheInfo,
}

var neomuttCacheCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean cache",
	Long: `Remove all cached headers and messages.

Examples:
  acorn mail neomutt cache clean`,
	RunE: runNeomuttCacheClean,
}

// Generate command
var neomuttGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate configuration files",
	Long: `Generate NeoMutt configuration files from component config.

Creates:
  - ~/.config/neomutt/neomuttrc (main config)

Examples:
  acorn mail neomutt generate
  acorn mail neomutt generate --dry-run`,
	Aliases: []string{"gen"},
	RunE:    runNeomuttGenerate,
}

// Account add commands
var neomuttAccountAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add email accounts",
	Long: `Add a new email account configuration.

Examples:
  acorn mail neomutt accounts add gmail
  acorn mail neomutt accounts add microsoft`,
}

var neomuttAccountAddGmailCmd = &cobra.Command{
	Use:   "gmail <email> <real-name>",
	Short: "Add Gmail account",
	Long: `Add a Gmail account with OAuth2 authentication.

Examples:
  acorn mail neomutt accounts add gmail ross.bercot@gmail.com "Ross Grinvalds"`,
	Args: cobra.ExactArgs(2),
	RunE: runNeomuttAddGmail,
}

var neomuttAccountAddMicrosoftCmd = &cobra.Command{
	Use:   "microsoft <email> <real-name>",
	Short: "Add Microsoft/Office365 account",
	Long: `Add a Microsoft/Office365 account with OAuth2 authentication.

Examples:
  acorn mail neomutt accounts add microsoft ross@lfblooms.farm "Ross Grinvalds"`,
	Aliases: []string{"ms", "office365", "outlook"},
	Args:    cobra.ExactArgs(2),
	RunE:    runNeomuttAddMicrosoft,
}

func init() {
	mailCmd.AddCommand(neomuttCmd)

	// Add subcommands
	neomuttCmd.AddCommand(neomuttStatusCmd)
	neomuttCmd.AddCommand(neomuttInstallCmd)
	neomuttCmd.AddCommand(neomuttLaunchCmd)
	neomuttCmd.AddCommand(neomuttInitCmd)
	neomuttCmd.AddCommand(neomuttAccountsCmd)

	// Token subcommands
	neomuttCmd.AddCommand(neomuttTokensCmd)
	neomuttTokensCmd.AddCommand(neomuttTokensStatusCmd)
	neomuttTokensCmd.AddCommand(neomuttTokensRefreshCmd)
	neomuttTokensCmd.AddCommand(neomuttTokensAuthorizeCmd)

	// Cache subcommands
	neomuttCmd.AddCommand(neomuttCacheCmd)
	neomuttCacheCmd.AddCommand(neomuttCacheInfoCmd)
	neomuttCacheCmd.AddCommand(neomuttCacheCleanCmd)

	// Generate command
	neomuttCmd.AddCommand(neomuttGenerateCmd)

	// Account add subcommands
	neomuttAccountsCmd.AddCommand(neomuttAccountAddCmd)
	neomuttAccountAddCmd.AddCommand(neomuttAccountAddGmailCmd)
	neomuttAccountAddCmd.AddCommand(neomuttAccountAddMicrosoftCmd)

	// Persistent flags
	neomuttCmd.PersistentFlags().BoolVar(&neomuttDryRun, "dry-run", false,
		"Show what would be done without executing")
	neomuttCmd.PersistentFlags().BoolVarP(&neomuttVerbose, "verbose", "v", false,
		"Show verbose output")
}

func newNeomuttHelper() *neomutt.Helper {
	return neomutt.NewHelper(neomuttVerbose, neomuttDryRun)
}

func runNeomuttStatus(cmd *cobra.Command, args []string) error {
	helper := newNeomuttHelper()
	status := helper.GetStatus()

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("NeoMutt Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if status.Installed {
		fmt.Fprintf(os.Stdout, "%s NeoMutt installed: %s\n", output.Success("✓"), status.Version)
	} else {
		fmt.Fprintf(os.Stdout, "%s NeoMutt not found\n", output.Error("✗"))
		fmt.Fprintln(os.Stdout, "  Install: brew install neomutt")
		return nil
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "Config: %s\n", status.ConfigDir)
	if status.ConfigExists {
		fmt.Fprintf(os.Stdout, "%s Configuration found\n", output.Success("✓"))
	} else {
		fmt.Fprintf(os.Stdout, "%s No configuration (run 'acorn mail neomutt init')\n", output.Warning("!"))
	}

	fmt.Fprintf(os.Stdout, "Cache: %s\n", status.CacheDir)
	fmt.Fprintf(os.Stdout, "Accounts: %d configured\n", status.AccountCount)

	return nil
}

func runNeomuttInstall(cmd *cobra.Command, args []string) error {
	helper := newNeomuttHelper()
	return helper.Install()
}

func runNeomuttLaunch(cmd *cobra.Command, args []string) error {
	helper := newNeomuttHelper()
	return helper.Launch(args...)
}

func runNeomuttInit(cmd *cobra.Command, args []string) error {
	helper := newNeomuttHelper()
	if err := helper.InitConfig(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Configuration directories created\n", output.Success("✓"))
	return nil
}

func runNeomuttAccounts(cmd *cobra.Command, args []string) error {
	helper := newNeomuttHelper()
	accounts, err := helper.ListAccounts()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(accounts)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Email Accounts"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if len(accounts) == 0 {
		fmt.Fprintln(os.Stdout, "No accounts configured")
		fmt.Fprintln(os.Stdout, "Add account configs to ~/.config/neomutt/accounts/")
		return nil
	}

	for _, a := range accounts {
		tokenStatus := ""
		if a.Type == "gmail" || a.Type == "microsoft" {
			if a.HasToken {
				tokenStatus = output.Success(" [OAuth2 ✓]")
			} else {
				tokenStatus = output.Warning(" [OAuth2 ✗]")
			}
		}

		fmt.Fprintf(os.Stdout, "  %-30s %-10s %s%s\n", a.Name, a.Type, a.Email, tokenStatus)
	}

	return nil
}

func runNeomuttTokensStatus(cmd *cobra.Command, args []string) error {
	helper := newNeomuttHelper()
	tokens, err := helper.GetTokenStatus()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(tokens)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("OAuth2 Token Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if len(tokens) == 0 {
		fmt.Fprintln(os.Stdout, "No OAuth2 accounts configured")
		return nil
	}

	for _, t := range tokens {
		status := output.Error("✗ Missing")
		if t.Exists {
			if t.Encrypted {
				status = output.Success("✓ Encrypted")
			} else {
				status = output.Warning("! Unencrypted")
			}
		}
		fmt.Fprintf(os.Stdout, "  %-30s %s\n", t.Account, status)
	}

	return nil
}

func runNeomuttTokensRefresh(cmd *cobra.Command, args []string) error {
	helper := newNeomuttHelper()
	accountName := args[0]

	fmt.Fprintf(os.Stdout, "Refreshing token for %s...\n", accountName)
	if err := helper.RefreshToken(accountName); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Token refreshed\n", output.Success("✓"))
	return nil
}

func runNeomuttTokensAuthorize(cmd *cobra.Command, args []string) error {
	helper := newNeomuttHelper()
	accountName := args[0]

	fmt.Fprintf(os.Stdout, "Authorizing OAuth2 for %s...\n", accountName)
	return helper.AuthorizeAccount(accountName)
}

func runNeomuttCacheInfo(cmd *cobra.Command, args []string) error {
	helper := newNeomuttHelper()
	info := helper.GetCacheInfo()

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(info)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("NeoMutt Cache"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintf(os.Stdout, "Headers:  %s (%s)\n", info.HeaderCache, info.HeaderSize)
	fmt.Fprintf(os.Stdout, "Messages: %s (%s)\n", info.MessageCache, info.MessageSize)

	return nil
}

func runNeomuttCacheClean(cmd *cobra.Command, args []string) error {
	helper := newNeomuttHelper()
	if err := helper.CleanCache(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Cache cleaned\n", output.Success("✓"))
	return nil
}

func runNeomuttGenerate(cmd *cobra.Command, args []string) error {
	helper := newNeomuttHelper()

	// Initialize config directories first
	if err := helper.InitConfig(); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Generating NeoMutt Configuration"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Generate main neomuttrc
	if err := helper.GenerateMainConfig(neomuttDryRun); err != nil {
		return err
	}

	if neomuttDryRun {
		fmt.Fprintf(os.Stdout, "\n%s Dry run complete. No files written.\n", output.Info("!"))
	} else {
		fmt.Fprintf(os.Stdout, "\n%s Configuration generated successfully\n", output.Success("✓"))
		fmt.Fprintln(os.Stdout, "\nNext steps:")
		fmt.Fprintln(os.Stdout, "  1. Add accounts: acorn mail neomutt accounts add gmail <email> <name>")
		fmt.Fprintln(os.Stdout, "  2. Authorize OAuth2: acorn mail neomutt tokens authorize <account>")
		fmt.Fprintln(os.Stdout, "  3. Launch: neomutt")
	}

	return nil
}

func runNeomuttAddGmail(cmd *cobra.Command, args []string) error {
	helper := newNeomuttHelper()
	email := args[0]
	realName := args[1]

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Adding Gmail Account"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	accountFile, err := helper.AddGmailAccount(email, realName, neomuttDryRun)
	if err != nil {
		return err
	}

	if neomuttDryRun {
		fmt.Fprintf(os.Stdout, "Would create: %s\n", accountFile)
	} else {
		fmt.Fprintf(os.Stdout, "%s Created: %s\n", output.Success("✓"), accountFile)
		fmt.Fprintln(os.Stdout, "\nNext steps:")
		fmt.Fprintln(os.Stdout, "  1. Authorize OAuth2: acorn mail neomutt tokens authorize gmail-"+extractUsername(email))
		fmt.Fprintln(os.Stdout, "  2. Update neomuttrc to source this account")
	}

	return nil
}

func runNeomuttAddMicrosoft(cmd *cobra.Command, args []string) error {
	helper := newNeomuttHelper()
	email := args[0]
	realName := args[1]

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Adding Microsoft Account"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	accountFile, err := helper.AddMicrosoftAccount(email, realName, neomuttDryRun)
	if err != nil {
		return err
	}

	if neomuttDryRun {
		fmt.Fprintf(os.Stdout, "Would create: %s\n", accountFile)
	} else {
		fmt.Fprintf(os.Stdout, "%s Created: %s\n", output.Success("✓"), accountFile)
		fmt.Fprintln(os.Stdout, "\nNext steps:")
		fmt.Fprintln(os.Stdout, "  1. Authorize OAuth2: acorn mail neomutt tokens authorize microsoft-"+extractAccountName(email))
		fmt.Fprintln(os.Stdout, "  2. Update neomuttrc to source this account")
	}

	return nil
}

// extractUsername extracts username from email (e.g., ross.bercot from ross.bercot@gmail.com)
func extractUsername(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 0 {
		return parts[0]
	}
	return email
}

// extractAccountName creates account name from email (e.g., domain.com-user from user@domain.com)
func extractAccountName(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) == 2 {
		return parts[1] + "-" + parts[0]
	}
	return email
}
