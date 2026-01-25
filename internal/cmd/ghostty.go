package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/terminal/ghostty"
	"github.com/mistergrinvalds/acorn/internal/utils/installer"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	ghosttyDryRun  bool
	ghosttyVerbose bool
)

// ghosttyCmd represents the ghostty command group
var ghosttyCmd = &cobra.Command{
	Use:   "ghostty",
	Short: "Ghostty terminal configuration",
	Long: `Manage Ghostty terminal emulator configuration.

Provides commands for theme switching, font configuration,
backup/restore, and viewing info.

Examples:
  acorn ghostty info              # Show Ghostty info
  acorn ghostty theme "Nord"      # Set theme
  acorn ghostty font "JetBrains Mono" 14
  acorn ghostty backup            # Backup current config`,
}

// ghosttyInfoCmd shows Ghostty info
var ghosttyInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show Ghostty information",
	Long: `Display Ghostty installation and configuration info.

Shows version, config path, current theme, and font settings.

Examples:
  acorn ghostty info
  acorn ghostty info -o json`,
	RunE: runGhosttyInfo,
}

// ghosttyThemeCmd sets theme
var ghosttyThemeCmd = &cobra.Command{
	Use:   "theme [name]",
	Short: "Set or list Ghostty themes",
	Long: `Set the Ghostty color theme or list available themes.

Without arguments, lists popular themes.
With a theme name, updates the config.

Examples:
  acorn ghostty theme                    # List themes
  acorn ghostty theme "Catppuccin Mocha" # Set theme`,
	Args: cobra.MaximumNArgs(1),
	RunE: runGhosttyTheme,
}

// ghosttyFontCmd sets font
var ghosttyFontCmd = &cobra.Command{
	Use:   "font <family> [size]",
	Short: "Set Ghostty font",
	Long: `Set the Ghostty font family and optionally size.

Examples:
  acorn ghostty font "JetBrains Mono"
  acorn ghostty font "Fira Code" 14`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runGhosttyFont,
}

// ghosttyBackupCmd creates backup
var ghosttyBackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup Ghostty config",
	Long: `Create a timestamped backup of the current Ghostty config.

Backups are stored in ~/.local/share/ghostty/backups/

Examples:
  acorn ghostty backup`,
	RunE: runGhosttyBackup,
}

// ghosttyBackupsCmd lists backups
var ghosttyBackupsCmd = &cobra.Command{
	Use:   "backups",
	Short: "List Ghostty config backups",
	Long: `List all available Ghostty config backups.

Examples:
  acorn ghostty backups
  acorn ghostty backups -o json`,
	Aliases: []string{"list-backups"},
	RunE:    runGhosttyBackups,
}

// ghosttyRestoreCmd restores from backup
var ghosttyRestoreCmd = &cobra.Command{
	Use:   "restore <backup>",
	Short: "Restore Ghostty config from backup",
	Long: `Restore Ghostty config from a backup file.

A backup of the current config is created before restoring.

Examples:
  acorn ghostty restore config.20240101_120000`,
	Args: cobra.ExactArgs(1),
	RunE: runGhosttyRestore,
}

// ghosttyConfigCmd shows config path
var ghosttyConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Show or edit Ghostty config",
	Long: `Show the Ghostty config file path.

Examples:
  acorn ghostty config`,
	RunE: runGhosttyConfig,
}

// ghosttyInstallCmd installs Ghostty
var ghosttyInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Ghostty terminal",
	Long: `Install the Ghostty terminal emulator.

Installs ghostty via Homebrew on macOS.

Examples:
  acorn ghostty install           # Install Ghostty
  acorn ghostty install --dry-run # Show what would be installed
  acorn ghostty install -v        # Verbose output`,
	RunE: runGhosttyInstall,
}

func init() {
	terminalCmd.AddCommand(ghosttyCmd)

	// Add subcommands
	ghosttyCmd.AddCommand(ghosttyInfoCmd)
	ghosttyCmd.AddCommand(ghosttyInstallCmd)
	ghosttyCmd.AddCommand(ghosttyThemeCmd)
	ghosttyCmd.AddCommand(ghosttyFontCmd)
	ghosttyCmd.AddCommand(ghosttyBackupCmd)
	ghosttyCmd.AddCommand(ghosttyBackupsCmd)
	ghosttyCmd.AddCommand(ghosttyRestoreCmd)
	ghosttyCmd.AddCommand(ghosttyConfigCmd)

	// Persistent flags
	ghosttyCmd.PersistentFlags().BoolVar(&ghosttyDryRun, "dry-run", false,
		"Show what would be done without executing")
	ghosttyCmd.PersistentFlags().BoolVarP(&ghosttyVerbose, "verbose", "v", false,
		"Show verbose output")
}

func runGhosttyInfo(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := ghostty.NewHelper(ghosttyVerbose)
	info := helper.GetInfo()

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(info)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Ghostty Terminal Information"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if info.Installed {
		fmt.Fprintf(os.Stdout, "%s Installed: %s\n", output.Success("✓"), info.Version)
	} else {
		fmt.Fprintf(os.Stdout, "%s Not installed\n", output.Error("✗"))
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Configuration:"))
	fmt.Fprintf(os.Stdout, "  Config: %s\n", info.Config)

	if info.Theme != "" {
		fmt.Fprintf(os.Stdout, "  Theme:  %s\n", info.Theme)
	} else {
		fmt.Fprintln(os.Stdout, "  Theme:  (manual palette)")
	}

	if info.Font != "" {
		fmt.Fprintf(os.Stdout, "  Font:   %s", info.Font)
		if info.FontSize != "" {
			fmt.Fprintf(os.Stdout, " (%s)", info.FontSize)
		}
		fmt.Fprintln(os.Stdout)
	}

	return nil
}

func runGhosttyTheme(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := ghostty.NewHelper(ghosttyVerbose)

	if len(args) == 0 {
		// List available themes
		themes := helper.GetAvailableThemes()

		if ioHelper.IsStructured() {
			return ioHelper.WriteOutput(map[string][]string{"themes": themes})
		}

		fmt.Fprintf(os.Stdout, "%s\n", output.Info("Available Themes"))
		fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Fprintln(os.Stdout, "Run 'ghostty +list-themes' for full list")
		fmt.Fprintln(os.Stdout)
		for _, theme := range themes {
			fmt.Fprintf(os.Stdout, "  • %s\n", theme)
		}
		return nil
	}

	// Set theme
	if err := helper.SetTheme(args[0]); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Theme set to: %s\n", output.Success("✓"), args[0])
	fmt.Fprintln(os.Stdout, "Press Cmd+Shift+, (macOS) or Ctrl+Shift+, (Linux) to reload.")

	return nil
}

func runGhosttyFont(cmd *cobra.Command, args []string) error {
	helper := ghostty.NewHelper(ghosttyVerbose)

	font := args[0]
	size := ""
	if len(args) > 1 {
		size = args[1]
	}

	if err := helper.SetFont(font, size); err != nil {
		return err
	}

	msg := fmt.Sprintf("Font set to: %s", font)
	if size != "" {
		msg += fmt.Sprintf(" (size: %s)", size)
	}
	fmt.Fprintf(os.Stdout, "%s %s\n", output.Success("✓"), msg)
	fmt.Fprintln(os.Stdout, "Press Cmd+Shift+, (macOS) or Ctrl+Shift+, (Linux) to reload.")

	return nil
}

func runGhosttyBackup(cmd *cobra.Command, args []string) error {
	helper := ghostty.NewHelper(ghosttyVerbose)

	backup, err := helper.CreateBackup()
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Backup created: %s\n", output.Success("✓"), backup.Path)
	return nil
}

func runGhosttyBackups(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := ghostty.NewHelper(ghosttyVerbose)

	backups, err := helper.ListBackups()
	if err != nil {
		return err
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(map[string]interface{}{"backups": backups})
	}

	if len(backups) == 0 {
		fmt.Fprintln(os.Stdout, "No backups found")
		return nil
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Ghostty Config Backups"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for _, b := range backups {
		fmt.Fprintf(os.Stdout, "  %s\n", b.Name)
	}

	fmt.Fprintf(os.Stdout, "\nTotal: %d backups\n", len(backups))
	return nil
}

func runGhosttyRestore(cmd *cobra.Command, args []string) error {
	helper := ghostty.NewHelper(ghosttyVerbose)

	if err := helper.RestoreBackup(args[0]); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Config restored from: %s\n", output.Success("✓"), args[0])
	fmt.Fprintln(os.Stdout, "Press Cmd+Shift+, (macOS) or Ctrl+Shift+, (Linux) to reload.")

	return nil
}

func runGhosttyConfig(cmd *cobra.Command, args []string) error {
	helper := ghostty.NewHelper(ghosttyVerbose)
	fmt.Fprintln(os.Stdout, helper.GetConfigPath())
	return nil
}

func runGhosttyInstall(cmd *cobra.Command, args []string) error {
	inst := installer.NewInstaller(
		installer.WithDryRun(ghosttyDryRun),
		installer.WithVerbose(ghosttyVerbose),
	)

	// Show platform info
	platform := inst.GetPlatform()
	if ghosttyVerbose {
		fmt.Fprintf(os.Stdout, "Platform: %s (%s)\n\n", platform.OS, platform.PackageManager)
	}

	// Get the plan first
	plan, err := inst.Plan(cmd.Context(), "ghostty")
	if err != nil {
		return err
	}

	// Show what will be installed
	if ghosttyDryRun {
		fmt.Fprintf(os.Stdout, "%s\n", output.Info("Ghostty Installation Plan"))
		fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Fprintf(os.Stdout, "Platform: %s (%s)\n\n", platform.OS, platform.PackageManager)
	}

	pending := plan.PendingTools()
	if len(pending) == 0 {
		fmt.Fprintf(os.Stdout, "%s All tools already installed\n", output.Success("✓"))
		return nil
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

	if ghosttyDryRun {
		fmt.Fprintln(os.Stdout, "\nRun without --dry-run to install.")
		return nil
	}

	// Execute installation
	fmt.Fprintln(os.Stdout, "\nInstalling...")
	result, err := inst.Install(cmd.Context(), "ghostty")
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
