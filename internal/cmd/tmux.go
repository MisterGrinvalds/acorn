package cmd

import (
	"fmt"
	"os"

	tmuxpkg "github.com/mistergrinvalds/acorn/internal/components/terminal/tmux"
	"github.com/mistergrinvalds/acorn/internal/utils/installer"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	tmuxDryRun  bool
	tmuxVerbose bool
)

// tmuxCmd represents the tmux command group
var tmuxCmd = &cobra.Command{
	Use:   "tmux",
	Short: "Tmux session management and TPM helpers",
	Long: `Helpers for tmux session management, TPM, and smug.

Provides commands for TPM installation, plugin management,
smug session configs, and cross-machine session sync.

Examples:
  acorn tmux info                # Show tmux info
  acorn tmux session list        # List active sessions
  acorn tmux tpm install         # Install TPM
  acorn tmux smug list           # List smug sessions
  acorn tmux smug repo-init      # Init smug git repo`,
}

// tmuxInfoCmd shows tmux information
var tmuxInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show tmux information",
	Long: `Display tmux version, configuration paths, and session info.

Examples:
  acorn tmux info
  acorn tmux info -o json`,
	RunE: runTmuxInfo,
}

// tmuxSessionCmd is the parent for session subcommands
var tmuxSessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Session management",
	Long:  `Commands for managing tmux sessions.`,
}

// tmuxSessionListCmd lists active sessions
var tmuxSessionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List active tmux sessions",
	Long: `List all active tmux sessions.

Examples:
  acorn tmux session list
  acorn tmux session list -o json`,
	Aliases: []string{"ls"},
	RunE:    runTmuxSessionList,
}

// tmuxTPMCmd is the parent for TPM subcommands
var tmuxTPMCmd = &cobra.Command{
	Use:   "tpm",
	Short: "Tmux Plugin Manager commands",
	Long:  `Commands for managing TPM (Tmux Plugin Manager).`,
}

// tmuxTPMInstallCmd installs TPM
var tmuxTPMInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Tmux Plugin Manager",
	Long: `Install TPM from GitHub.

Examples:
  acorn tmux tpm install`,
	RunE: runTmuxTPMInstall,
}

// tmuxTPMUpdateCmd updates TPM
var tmuxTPMUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Tmux Plugin Manager",
	Long: `Update TPM to the latest version.

Examples:
  acorn tmux tpm update`,
	RunE: runTmuxTPMUpdate,
}

// tmuxTPMPluginsInstallCmd installs plugins
var tmuxTPMPluginsInstallCmd = &cobra.Command{
	Use:   "plugins-install",
	Short: "Install all tmux plugins",
	Long: `Install all plugins defined in tmux.conf via TPM.

Examples:
  acorn tmux tpm plugins-install`,
	RunE: runTmuxTPMPluginsInstall,
}

// tmuxTPMPluginsUpdateCmd updates plugins
var tmuxTPMPluginsUpdateCmd = &cobra.Command{
	Use:   "plugins-update",
	Short: "Update all tmux plugins",
	Long: `Update all installed tmux plugins via TPM.

Examples:
  acorn tmux tpm plugins-update`,
	RunE: runTmuxTPMPluginsUpdate,
}

// tmuxConfigCmd is the parent for config subcommands
var tmuxConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management",
	Long:  `Commands for managing tmux configuration.`,
}

// tmuxConfigReloadCmd reloads the config
var tmuxConfigReloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reload tmux configuration",
	Long: `Reload the tmux configuration file.

Must be run from within a tmux session.

Examples:
  acorn tmux config reload`,
	RunE: runTmuxConfigReload,
}

// tmuxSmugCmd is the parent for smug subcommands
var tmuxSmugCmd = &cobra.Command{
	Use:   "smug",
	Short: "Smug session management",
	Long: `Commands for managing smug session configurations.

Smug provides persistent, versioned session configurations
that can be synced across machines via git.`,
}

// tmuxSmugListCmd lists smug sessions
var tmuxSmugListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available smug sessions",
	Long: `List all available smug session configurations.

Examples:
  acorn tmux smug list
  acorn tmux smug list -o json`,
	Aliases: []string{"ls"},
	RunE:    runTmuxSmugList,
}

// tmuxSmugNewCmd creates a new smug session
var tmuxSmugNewCmd = &cobra.Command{
	Use:   "new <name>",
	Short: "Create a new smug session config",
	Long: `Create a new smug session configuration from template.

Examples:
  acorn tmux smug new myproject`,
	Args: cobra.ExactArgs(1),
	RunE: runTmuxSmugNew,
}

// tmuxSmugInstallCmd installs smug
var tmuxSmugInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install smug",
	Long: `Install smug using brew or go install.

Examples:
  acorn tmux smug install`,
	RunE: runTmuxSmugInstall,
}

// tmuxSmugLinkCmd links smug configs
var tmuxSmugLinkCmd = &cobra.Command{
	Use:   "link",
	Short: "Link smug configs from git repo or dotfiles",
	Long: `Link smug session configs from the git repo or dotfiles.

Prefers the git repo if initialized, otherwise uses dotfiles.

Examples:
  acorn tmux smug link`,
	RunE: runTmuxSmugLink,
}

// tmuxSmugRepoInitCmd initializes the smug repo
var tmuxSmugRepoInitCmd = &cobra.Command{
	Use:   "repo-init",
	Short: "Initialize smug sessions git repo",
	Long: `Clone or update the smug sessions git repository.

This enables cross-machine session portability.

Examples:
  acorn tmux smug repo-init`,
	RunE: runTmuxSmugRepoInit,
}

// tmuxSmugStatusCmd shows repo status
var tmuxSmugStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show smug repo status",
	Long: `Show the status of the smug sessions git repository.

Examples:
  acorn tmux smug status
  acorn tmux smug status -o json`,
	RunE: runTmuxSmugStatus,
}

// tmuxSmugPullCmd pulls from remote
var tmuxSmugPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull latest sessions from remote",
	Long: `Pull the latest smug session configs from the remote repository.

Examples:
  acorn tmux smug pull`,
	RunE: runTmuxSmugPull,
}

// tmuxSmugPushCmd pushes to remote
var tmuxSmugPushCmd = &cobra.Command{
	Use:   "push [message]",
	Short: "Commit and push session changes",
	Long: `Commit local changes and push to the remote repository.

Examples:
  acorn tmux smug push
  acorn tmux smug push "Add new project session"`,
	Args: cobra.MaximumNArgs(1),
	RunE: runTmuxSmugPush,
}

// tmuxSmugSyncCmd does a full sync
var tmuxSmugSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Full sync (pull + push)",
	Long: `Perform a full sync: pull remote changes, then push local changes.

Examples:
  acorn tmux smug sync`,
	RunE: runTmuxSmugSync,
}

// tmuxInstallCmd installs tmux component tools
var tmuxInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install tmux tools",
	Long: `Install required tools for the tmux component.

Installs tmux, smug (session manager), and fzf (fuzzy finder).
Uses brew on macOS and apt on Linux.

Examples:
  acorn tmux install           # Install all tmux tools
  acorn tmux install --dry-run # Show what would be installed
  acorn tmux install -v        # Verbose output`,
	RunE: runTmuxInstall,
}

func init() {
	terminalCmd.AddCommand(tmuxCmd)

	// Main subcommands
	tmuxCmd.AddCommand(tmuxInfoCmd)
	tmuxCmd.AddCommand(tmuxInstallCmd)
	tmuxCmd.AddCommand(tmuxSessionCmd)
	tmuxCmd.AddCommand(tmuxTPMCmd)
	tmuxCmd.AddCommand(tmuxConfigCmd)
	tmuxCmd.AddCommand(tmuxSmugCmd)

	// Session subcommands
	tmuxSessionCmd.AddCommand(tmuxSessionListCmd)

	// TPM subcommands
	tmuxTPMCmd.AddCommand(tmuxTPMInstallCmd)
	tmuxTPMCmd.AddCommand(tmuxTPMUpdateCmd)
	tmuxTPMCmd.AddCommand(tmuxTPMPluginsInstallCmd)
	tmuxTPMCmd.AddCommand(tmuxTPMPluginsUpdateCmd)

	// Config subcommands
	tmuxConfigCmd.AddCommand(tmuxConfigReloadCmd)

	// Smug subcommands
	tmuxSmugCmd.AddCommand(tmuxSmugListCmd)
	tmuxSmugCmd.AddCommand(tmuxSmugNewCmd)
	tmuxSmugCmd.AddCommand(tmuxSmugInstallCmd)
	tmuxSmugCmd.AddCommand(tmuxSmugLinkCmd)
	tmuxSmugCmd.AddCommand(tmuxSmugRepoInitCmd)
	tmuxSmugCmd.AddCommand(tmuxSmugStatusCmd)
	tmuxSmugCmd.AddCommand(tmuxSmugPullCmd)
	tmuxSmugCmd.AddCommand(tmuxSmugPushCmd)
	tmuxSmugCmd.AddCommand(tmuxSmugSyncCmd)

	// Persistent flags
	tmuxCmd.PersistentFlags().BoolVar(&tmuxDryRun, "dry-run", false,
		"Show what would be done without executing")
	tmuxCmd.PersistentFlags().BoolVarP(&tmuxVerbose, "verbose", "v", false,
		"Show verbose output")
}

func runTmuxInfo(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := tmuxpkg.NewHelper(tmuxVerbose, tmuxDryRun)

	info, err := helper.GetInfo()
	if err != nil {
		return err
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(info)
	}

	fmt.Fprintf(os.Stdout, "%s\n\n", output.Info("Tmux Information"))
	fmt.Fprintf(os.Stdout, "Version:     %s\n", info.Version)
	fmt.Fprintf(os.Stdout, "Config:      %s\n", info.ConfigFile)
	fmt.Fprintf(os.Stdout, "Plugins:     %s\n", info.PluginDir)
	if info.TPMInstalled {
		fmt.Fprintf(os.Stdout, "TPM:         %s\n", output.Success("installed"))
	} else {
		fmt.Fprintf(os.Stdout, "TPM:         %s (run: acorn tmux tpm install)\n", output.Warning("not installed"))
	}

	if len(info.Sessions) > 0 {
		fmt.Fprintf(os.Stdout, "\n%s\n", output.Info("Active Sessions"))
		for _, s := range info.Sessions {
			status := ""
			if s.Attached {
				status = " " + output.Success("(attached)")
			}
			fmt.Fprintf(os.Stdout, "  %s: %d windows%s\n", s.Name, s.Windows, status)
		}
	} else {
		fmt.Fprintln(os.Stdout, "\nNo active sessions")
	}

	return nil
}

func runTmuxSessionList(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := tmuxpkg.NewHelper(tmuxVerbose, tmuxDryRun)

	sessions, err := helper.ListSessions()
	if err != nil {
		return err
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(sessions)
	}

	if len(sessions) == 0 {
		fmt.Fprintln(os.Stdout, "No active tmux sessions")
		return nil
	}

	fmt.Fprintf(os.Stdout, "%s\n\n", output.Info("Active Sessions"))
	for _, s := range sessions {
		status := ""
		if s.Attached {
			status = " " + output.Success("(attached)")
		}
		fmt.Fprintf(os.Stdout, "  %s: %d windows%s\n", s.Name, s.Windows, status)
	}

	return nil
}

func runTmuxTPMInstall(cmd *cobra.Command, args []string) error {
	helper := tmuxpkg.NewHelper(tmuxVerbose, tmuxDryRun)

	if err := helper.InstallTPM(); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "\n%s TPM installed successfully!\n", output.Success("✓"))
	fmt.Fprintln(os.Stdout, "\nNext steps:")
	fmt.Fprintln(os.Stdout, "  1. Start tmux: tmux")
	fmt.Fprintln(os.Stdout, "  2. Install plugins: prefix + I")
	fmt.Fprintln(os.Stdout, "  3. Update plugins: prefix + U")
	return nil
}

func runTmuxTPMUpdate(cmd *cobra.Command, args []string) error {
	helper := tmuxpkg.NewHelper(tmuxVerbose, tmuxDryRun)

	if err := helper.UpdateTPM(); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s TPM updated!\n", output.Success("✓"))
	return nil
}

func runTmuxTPMPluginsInstall(cmd *cobra.Command, args []string) error {
	helper := tmuxpkg.NewHelper(tmuxVerbose, tmuxDryRun)

	if err := helper.InstallPlugins(); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Plugins installed!\n", output.Success("✓"))
	return nil
}

func runTmuxTPMPluginsUpdate(cmd *cobra.Command, args []string) error {
	helper := tmuxpkg.NewHelper(tmuxVerbose, tmuxDryRun)

	if err := helper.UpdatePlugins(); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Plugins updated!\n", output.Success("✓"))
	return nil
}

func runTmuxConfigReload(cmd *cobra.Command, args []string) error {
	helper := tmuxpkg.NewHelper(tmuxVerbose, tmuxDryRun)

	if err := helper.ReloadConfig(); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Config reloaded!\n", output.Success("✓"))
	return nil
}

func runTmuxSmugList(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := tmuxpkg.NewHelper(tmuxVerbose, tmuxDryRun)

	sessions, err := helper.ListSmugSessions()
	if err != nil {
		return err
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(sessions)
	}

	if len(sessions) == 0 {
		fmt.Fprintln(os.Stdout, "No smug sessions found.")
		fmt.Fprintf(os.Stdout, "  Run: acorn tmux smug new <name>\n")
		return nil
	}

	fmt.Fprintf(os.Stdout, "%s\n\n", output.Info("Available Smug Sessions"))
	for _, s := range sessions {
		if s.Description != "" {
			fmt.Fprintf(os.Stdout, "  %s - %s\n", s.Name, s.Description)
		} else {
			fmt.Fprintf(os.Stdout, "  %s\n", s.Name)
		}
	}

	fmt.Fprintln(os.Stdout, "\nCommands:")
	fmt.Fprintln(os.Stdout, "  smug start <name>    - Start a session")
	fmt.Fprintln(os.Stdout, "  smug stop <name>     - Stop a session")

	return nil
}

func runTmuxSmugNew(cmd *cobra.Command, args []string) error {
	helper := tmuxpkg.NewHelper(tmuxVerbose, tmuxDryRun)

	configFile, err := helper.CreateSmugSession(args[0])
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Created: %s\n", output.Success("✓"), configFile)
	fmt.Fprintln(os.Stdout, "\nNext steps:")
	fmt.Fprintf(os.Stdout, "  Edit: $EDITOR %s\n", configFile)
	fmt.Fprintf(os.Stdout, "  Start: smug start %s\n", args[0])
	return nil
}

func runTmuxSmugInstall(cmd *cobra.Command, args []string) error {
	helper := tmuxpkg.NewHelper(tmuxVerbose, tmuxDryRun)

	if err := helper.InstallSmug(); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Smug installed!\n", output.Success("✓"))
	return nil
}

func runTmuxSmugLink(cmd *cobra.Command, args []string) error {
	helper := tmuxpkg.NewHelper(tmuxVerbose, tmuxDryRun)

	dotfilesRoot := os.Getenv("DOTFILES_ROOT")
	if dotfilesRoot == "" {
		dotfilesRoot = "."
	}

	if err := helper.SmugLinkConfigs(dotfilesRoot); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Smug configs linked!\n", output.Success("✓"))
	return nil
}

func runTmuxSmugRepoInit(cmd *cobra.Command, args []string) error {
	helper := tmuxpkg.NewHelper(tmuxVerbose, tmuxDryRun)

	if err := helper.SmugRepoInit(); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "\n%s Smug repo initialized!\n", output.Success("✓"))
	fmt.Fprintf(os.Stdout, "Location: %s\n", tmuxpkg.GetSmugRepoDir())
	fmt.Fprintln(os.Stdout, "\nCommands:")
	fmt.Fprintln(os.Stdout, "  acorn tmux smug status   - Show repo status")
	fmt.Fprintln(os.Stdout, "  acorn tmux smug pull     - Pull latest sessions")
	fmt.Fprintln(os.Stdout, "  acorn tmux smug push     - Commit and push changes")
	fmt.Fprintln(os.Stdout, "  acorn tmux smug sync     - Full sync")
	return nil
}

func runTmuxSmugStatus(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := tmuxpkg.NewHelper(tmuxVerbose, tmuxDryRun)

	status, err := helper.SmugRepoStatus()
	if err != nil {
		return err
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	fmt.Fprintf(os.Stdout, "%s\n\n", output.Info("Smug Sessions Status"))
	fmt.Fprintf(os.Stdout, "Location: %s\n", status.Location)

	if !status.Initialized {
		fmt.Fprintf(os.Stdout, "Status:   %s (run: acorn tmux smug repo-init)\n", output.Warning("not initialized"))
		return nil
	}

	fmt.Fprintf(os.Stdout, "Branch:   %s\n", status.Branch)
	if status.Clean {
		fmt.Fprintf(os.Stdout, "Status:   %s\n", output.Success("clean"))
	} else {
		fmt.Fprintf(os.Stdout, "Status:   %s\n", output.Warning("uncommitted changes"))
	}

	if len(status.Sessions) > 0 {
		fmt.Fprintln(os.Stdout, "\nSessions:")
		for _, s := range status.Sessions {
			fmt.Fprintf(os.Stdout, "  - %s\n", s.Name)
		}
	}

	return nil
}

func runTmuxSmugPull(cmd *cobra.Command, args []string) error {
	helper := tmuxpkg.NewHelper(tmuxVerbose, tmuxDryRun)

	if err := helper.SmugRepoPull(); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Sessions updated!\n", output.Success("✓"))
	return nil
}

func runTmuxSmugPush(cmd *cobra.Command, args []string) error {
	helper := tmuxpkg.NewHelper(tmuxVerbose, tmuxDryRun)

	message := ""
	if len(args) > 0 {
		message = args[0]
	}

	if err := helper.SmugRepoPush(message); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Sessions pushed!\n", output.Success("✓"))
	return nil
}

func runTmuxSmugSync(cmd *cobra.Command, args []string) error {
	helper := tmuxpkg.NewHelper(tmuxVerbose, tmuxDryRun)

	if err := helper.SmugRepoSync(); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "\n%s Sync complete!\n", output.Success("✓"))
	return nil
}

func runTmuxInstall(cmd *cobra.Command, args []string) error {
	inst := installer.NewInstaller(
		installer.WithDryRun(tmuxDryRun),
		installer.WithVerbose(tmuxVerbose),
	)

	// Show platform info
	platform := inst.GetPlatform()
	if tmuxVerbose {
		fmt.Fprintf(os.Stdout, "Platform: %s (%s)\n\n", platform.OS, platform.PackageManager)
	}

	// Get the plan first
	plan, err := inst.Plan(cmd.Context(), "tmux")
	if err != nil {
		return err
	}

	// Show what will be installed
	if tmuxDryRun {
		fmt.Fprintf(os.Stdout, "%s\n", output.Info("Tmux Installation Plan"))
		fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Fprintf(os.Stdout, "Platform: %s (%s)\n\n", platform.OS, platform.PackageManager)
	}

	pending := plan.PendingTools()
	if len(pending) == 0 {
		fmt.Fprintf(os.Stdout, "%s All tools already installed\n", output.Success("✓"))
		return nil
	}

	// Show prerequisites
	if len(plan.Prerequisites) > 0 {
		fmt.Fprintln(os.Stdout, "Prerequisites:")
		for _, t := range plan.Prerequisites {
			status := output.Warning("○")
			suffix := ""
			if t.AlreadyInstalled {
				status = output.Success("✓")
				suffix = " (installed)"
			}
			fmt.Fprintf(os.Stdout, "  %s %s%s\n", status, t.Name, suffix)
		}
		fmt.Fprintln(os.Stdout)
	}

	// Show tools
	fmt.Fprintln(os.Stdout, "Tools:")
	for _, t := range plan.Tools {
		status := output.Warning("○")
		suffix := ""
		if t.AlreadyInstalled {
			status = output.Success("✓")
			suffix = " (installed)"
		}
		fmt.Fprintf(os.Stdout, "  %s %s - %s%s\n", status, t.Name, t.Description, suffix)
	}

	if tmuxDryRun {
		fmt.Fprintln(os.Stdout, "\nRun without --dry-run to install.")
		return nil
	}

	// Execute installation
	fmt.Fprintln(os.Stdout, "\nInstalling...")
	result, err := inst.Install(cmd.Context(), "tmux")
	if err != nil {
		return err
	}

	installed, skipped, failed := result.Summary()
	if failed == 0 {
		fmt.Fprintf(os.Stdout, "%s Installation complete (%d installed, %d skipped)\n",
			output.Success("✓"), installed, skipped)
	} else {
		fmt.Fprintf(os.Stdout, "%s Installation failed (%d installed, %d skipped, %d failed)\n",
			output.Error("✗"), installed, skipped, failed)
	}

	return nil
}
