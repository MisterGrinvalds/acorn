package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/mail/neomutt"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	neomuttOutputFormat string
	neomuttDryRun       bool
	neomuttVerbose      bool
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

	// Persistent flags
	neomuttCmd.PersistentFlags().StringVarP(&neomuttOutputFormat, "output", "o", "table",
		"Output format (table|json|yaml)")
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

	format, err := output.ParseFormat(neomuttOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(status)
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

	format, err := output.ParseFormat(neomuttOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(accounts)
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

	format, err := output.ParseFormat(neomuttOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(tokens)
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

	format, err := output.ParseFormat(neomuttOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(info)
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
