package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mistergrinvalds/acorn/internal/components"
	"github.com/mistergrinvalds/acorn/internal/components/filesync"
	"github.com/mistergrinvalds/acorn/internal/components/opencode"
	"github.com/mistergrinvalds/acorn/internal/utils/config"
	"github.com/mistergrinvalds/acorn/internal/utils/configcmd"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	opencodeDryRun  bool
	opencodeVerbose bool
)

// opencodeCmd represents the opencode command group
var opencodeCmd = &cobra.Command{
	Use:   "opencode",
	Short: "OpenCode AI coding assistant",
	Long: `OpenCode AI coding assistant for the terminal.

OpenCode is an open-source AI coding agent that helps you write,
debug, and understand code directly in your terminal.

Examples:
  acorn ai opencode status     # Show status
  acorn ai opencode launch     # Start OpenCode TUI
  acorn ai opencode providers  # List AI providers
  acorn ai opencode upgrade    # Upgrade to latest`,
	Aliases: []string{"oc"},
}

// opencodeStatusCmd shows status
var opencodeStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show OpenCode status",
	Long: `Display OpenCode installation and configuration status.

Examples:
  acorn ai opencode status
  acorn ai opencode status -o json`,
	RunE: runOpencodeStatus,
}

// opencodeLaunchCmd launches OpenCode
var opencodeLaunchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch OpenCode TUI",
	Long: `Launch the OpenCode terminal user interface.

This starts the interactive OpenCode session in your terminal.

Examples:
  acorn ai opencode launch`,
	Aliases: []string{"start", "run"},
	RunE:    runOpencodeLaunch,
}

// opencodeProvidersCmd lists providers
var opencodeProvidersCmd = &cobra.Command{
	Use:   "providers",
	Short: "List AI providers",
	Long: `List supported AI providers for OpenCode.

OpenCode supports multiple AI providers including OpenAI,
Anthropic, Google, Groq, AWS Bedrock, and more.

Examples:
  acorn ai opencode providers
  acorn ai opencode providers -o json`,
	RunE: runOpencodeProviders,
}

// opencodeInstallCmd installs OpenCode
var opencodeInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install OpenCode",
	Long: `Install OpenCode using Homebrew.

Examples:
  acorn ai opencode install`,
	RunE: runOpencodeInstall,
}

// opencodeUpgradeCmd upgrades OpenCode
var opencodeUpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade OpenCode",
	Long: `Upgrade OpenCode to the latest version.

Examples:
  acorn ai opencode upgrade`,
	RunE: runOpencodeUpgrade,
}

// opencodeSyncCmd syncs opencode config files
var opencodeSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync OpenCode configuration files",
	Long: `Synchronize OpenCode configuration from dotfiles to ~/.config/opencode.

Syncs agents and commands directories via symlinks.

Modes:
  - symlink: Create symlinks for directories (agents, commands)

Examples:
  acorn ai opencode sync              # Sync all config files
  acorn ai opencode sync --dry-run    # Show what would be synced
  acorn ai opencode sync status       # Check sync status`,
	RunE: runOpencodeSync,
}

// opencodeSyncStatusCmd shows sync status
var opencodeSyncStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check sync status of OpenCode config files",
	Long: `Display the synchronization status of OpenCode configuration files.

Shows which files are synced, missing, or out of date.

Examples:
  acorn ai opencode sync status
  acorn ai opencode sync status -o json`,
	RunE: runOpencodeSyncStatus,
}

// opencodeAggregateCmd aggregates agents/commands from repos
var opcodeAggregateCmd = &cobra.Command{
	Use:   "aggregate [search-dir]",
	Short: "Aggregate agents and commands from repositories",
	Long: `Scan repositories for .opencode directories and aggregate their agents and commands.

By default, scans ~/Repos for repositories up to 3 levels deep.
Aggregates agents from .opencode/agents/ and commands from .opencode/commands/.

Deduplication:
  - Files with identical content are skipped
  - Files with same name but different content are renamed with repo prefix

Examples:
  acorn ai opencode aggregate              # Scan ~/Repos
  acorn ai opencode aggregate ~/Projects   # Scan custom directory
  acorn ai opencode aggregate --dry-run    # Preview without changes`,
	RunE: runOpencodeAggregate,
}

// opcodeAggregateListCmd lists aggregated items
var opcodeAggregateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all aggregated agents and commands",
	Long: `Display all agents and commands in the OpenCode configuration.

Examples:
  acorn ai opencode aggregate list
  acorn ai opencode aggregate list -o json`,
	RunE: runOpencodeAggregateList,
}

func init() {

	// Add subcommands
	opencodeCmd.AddCommand(opencodeStatusCmd)
	opencodeCmd.AddCommand(opencodeLaunchCmd)
	opencodeCmd.AddCommand(opencodeProvidersCmd)
	opencodeCmd.AddCommand(opencodeInstallCmd)
	opencodeCmd.AddCommand(opencodeUpgradeCmd)
	opencodeCmd.AddCommand(opencodeSyncCmd)
	opencodeCmd.AddCommand(opcodeAggregateCmd)
	opencodeCmd.AddCommand(configcmd.NewConfigRouter("opencode"))

	// Add sync subcommands
	opencodeSyncCmd.AddCommand(opencodeSyncStatusCmd)

	// Add aggregate subcommands
	opcodeAggregateCmd.AddCommand(opcodeAggregateListCmd)

	// Persistent flags
	opencodeCmd.PersistentFlags().BoolVar(&opencodeDryRun, "dry-run", false,
		"Show what would be done without executing")
	opencodeCmd.PersistentFlags().BoolVarP(&opencodeVerbose, "verbose", "v", false,
		"Show verbose output")
}

func newOpencodeHelper() *opencode.Helper {
	return opencode.NewHelper(opencodeVerbose, opencodeDryRun)
}

func runOpencodeStatus(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := newOpencodeHelper()
	status := helper.GetStatus()

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("OpenCode Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if status.Installed {
		fmt.Fprintf(os.Stdout, "%s OpenCode installed: v%s\n", output.Success("✓"), status.Version)
	} else {
		fmt.Fprintf(os.Stdout, "%s OpenCode not installed\n", output.Error("✗"))
		fmt.Fprintln(os.Stdout, "  Install: acorn ai opencode install")
		return nil
	}

	fmt.Fprintln(os.Stdout)

	if status.HasConfig {
		fmt.Fprintf(os.Stdout, "%s Configuration found\n", output.Success("✓"))
	} else {
		fmt.Fprintf(os.Stdout, "%s No configuration\n", output.Info("ℹ"))
		fmt.Fprintln(os.Stdout, "  Run 'opencode' and use /init to initialize")
	}

	fmt.Fprintf(os.Stdout, "\nConfig dir: %s\n", status.ConfigDir)

	return nil
}

func runOpencodeLaunch(cmd *cobra.Command, args []string) error {
	helper := newOpencodeHelper()
	return helper.Launch(args...)
}

func runOpencodeProviders(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := newOpencodeHelper()
	providers := helper.ListProviders()

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(providers)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Supported Providers"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for _, p := range providers {
		fmt.Fprintf(os.Stdout, "  %-15s %s\n", p.Name, p.Model)
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "Set provider via OPENCODE_PROVIDER environment variable")
	fmt.Fprintln(os.Stdout, "Set API key via ANTHROPIC_API_KEY, OPENAI_API_KEY, etc.")

	return nil
}

func runOpencodeInstall(cmd *cobra.Command, args []string) error {
	helper := newOpencodeHelper()
	if err := helper.Install(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s OpenCode installed\n", output.Success("✓"))
	return nil
}

func runOpencodeUpgrade(cmd *cobra.Command, args []string) error {
	helper := newOpencodeHelper()
	if err := helper.Upgrade(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s OpenCode upgraded\n", output.Success("✓"))
	return nil
}

func runOpencodeSync(cmd *cobra.Command, args []string) error {
	// Get dotfiles root
	dotfilesRoot, err := getDotfilesRoot()
	if err != nil {
		return fmt.Errorf("failed to get dotfiles root: %w", err)
	}

	// Load opencode component config
	loader := config.NewComponentLoader()
	cfg, err := loader.LoadBase("opencode")
	if err != nil {
		return fmt.Errorf("failed to load opencode config: %w", err)
	}

	if !cfg.HasSyncFiles() {
		fmt.Fprintf(os.Stdout, "%s No sync files configured\n", output.Info("ℹ"))
		return nil
	}

	// Create syncer and sync files
	syncer := filesync.NewSyncer(dotfilesRoot, opencodeDryRun, opencodeVerbose)
	result, err := syncer.Sync(cfg.GetSyncFiles())
	if err != nil {
		return fmt.Errorf("sync failed: %w", err)
	}

	// Output results
	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(result)
	}

	// Table format
	if result.DryRun {
		fmt.Fprintf(os.Stdout, "%s Sync Preview (dry-run)\n", output.Info("ℹ"))
	} else {
		fmt.Fprintf(os.Stdout, "%s OpenCode Config Sync\n", output.Info("ℹ"))
	}
	fmt.Fprintln(os.Stdout)

	if len(result.Synced) > 0 {
		for _, f := range result.Synced {
			fmt.Fprintf(os.Stdout, "  %s %s → %s (%s)\n",
				output.Success("✓"), f.Source, f.Target, f.Mode)
		}
	}

	if len(result.Skipped) > 0 && opencodeVerbose {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Unchanged:")
		for _, f := range result.Skipped {
			fmt.Fprintf(os.Stdout, "  %s %s\n", output.Info("○"), f.Target)
		}
	}

	if len(result.Errors) > 0 {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Errors:")
		for _, e := range result.Errors {
			fmt.Fprintf(os.Stdout, "  %s %s: %s\n", output.Error("✗"), e.Source, e.Error)
		}
	}

	if len(result.Synced) > 0 || len(result.Skipped) > 0 {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintf(os.Stdout, "Summary: %d synced, %d unchanged, %d errors\n",
			len(result.Synced), len(result.Skipped), len(result.Errors))
	}

	return nil
}

func runOpencodeSyncStatus(cmd *cobra.Command, args []string) error {
	// Get dotfiles root
	dotfilesRoot, err := getDotfilesRoot()
	if err != nil {
		return fmt.Errorf("failed to get dotfiles root: %w", err)
	}

	// Load opencode component config
	loader := config.NewComponentLoader()
	cfg, err := loader.LoadBase("opencode")
	if err != nil {
		return fmt.Errorf("failed to load opencode config: %w", err)
	}

	if !cfg.HasSyncFiles() {
		fmt.Fprintf(os.Stdout, "%s No sync files configured\n", output.Info("ℹ"))
		return nil
	}

	// Create syncer and check status
	syncer := filesync.NewSyncer(dotfilesRoot, false, opencodeVerbose)
	result, err := syncer.Status(cfg.GetSyncFiles())
	if err != nil {
		return fmt.Errorf("status check failed: %w", err)
	}

	// Output results
	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(result)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s OpenCode Sync Status\n", output.Info("ℹ"))
	fmt.Fprintln(os.Stdout)

	// Show synced files
	for _, f := range result.Synced {
		fmt.Fprintf(os.Stdout, "  %s %s (%s)\n", output.Success("✓"), f.Target, f.Action)
	}

	// Show files needing sync
	for _, f := range result.Skipped {
		var statusIcon string
		switch f.Action {
		case "missing":
			statusIcon = output.Error("✗")
		case "not_symlink", "wrong_target":
			statusIcon = output.Warning("!")
		default:
			statusIcon = output.Info("?")
		}
		fmt.Fprintf(os.Stdout, "  %s %s (%s)\n", statusIcon, f.Target, f.Action)
	}

	return nil
}

func runOpencodeAggregate(cmd *cobra.Command, args []string) error {
	helper := newOpencodeHelper()

	// Default search directory
	searchDir := ""
	if len(args) > 0 {
		searchDir = args[0]
	} else {
		home, _ := os.UserHomeDir()
		searchDir = filepath.Join(home, "Repos")
	}

	result, err := helper.Aggregate(searchDir)
	if err != nil {
		return fmt.Errorf("aggregate failed: %w", err)
	}

	// Output results
	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(result)
	}

	// Table format
	if opencodeDryRun {
		fmt.Fprintf(os.Stdout, "%s Aggregate Preview (dry-run)\n", output.Info("ℹ"))
	} else {
		fmt.Fprintf(os.Stdout, "%s OpenCode Aggregate\n", output.Info("ℹ"))
	}
	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "Search directory: %s\n", result.SearchDir)
	fmt.Fprintf(os.Stdout, "Target directory: %s\n", result.TargetDir)
	fmt.Fprintln(os.Stdout)

	if len(result.Items) > 0 && opencodeVerbose {
		for _, item := range result.Items {
			var icon string
			switch item.Action {
			case "added":
				icon = output.Success("✓")
			case "renamed":
				icon = output.Warning("!")
			case "skipped":
				icon = output.Info("○")
			}
			fmt.Fprintf(os.Stdout, "  %s [%s] %s from %s\n",
				icon, item.Type, item.FileName, item.SourceRepo)
		}
		fmt.Fprintln(os.Stdout)
	}

	fmt.Fprintf(os.Stdout, "Summary: %d repos scanned, %d agents, %d commands, %d skipped, %d renamed\n",
		result.ReposScanned, result.AgentsAdded, result.CommandsAdded, result.Skipped, result.Renamed)

	return nil
}

func runOpencodeAggregateList(cmd *cobra.Command, args []string) error {
	helper := newOpencodeHelper()

	result, err := helper.List()
	if err != nil {
		return fmt.Errorf("list failed: %w", err)
	}

	// Output results
	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(result)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s OpenCode Agents & Commands\n", output.Info("ℹ"))
	fmt.Fprintln(os.Stdout)

	if len(result.Agents) > 0 {
		fmt.Fprintln(os.Stdout, "Agents:")
		for _, agent := range result.Agents {
			fmt.Fprintf(os.Stdout, "  @%s\n", agent)
		}
		fmt.Fprintln(os.Stdout)
	}

	if len(result.Commands) > 0 {
		fmt.Fprintln(os.Stdout, "Commands:")
		for _, cmd := range result.Commands {
			fmt.Fprintf(os.Stdout, "  %s\n", cmd)
		}
		fmt.Fprintln(os.Stdout)
	}

	fmt.Fprintf(os.Stdout, "Total: %d agents, %d commands\n",
		len(result.Agents), len(result.Commands))

	return nil
}

func init() {
	components.Register(&components.Registration{
		Name:        "opencode",
		RegisterCmd: func() *cobra.Command { return opencodeCmd },
	})
}
