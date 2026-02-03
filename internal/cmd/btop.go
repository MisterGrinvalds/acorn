package cmd

import (
	"github.com/mistergrinvalds/acorn/internal/components"
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/btop"
	"github.com/mistergrinvalds/acorn/internal/utils/configcmd"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	btopDryRun  bool
	btopVerbose bool
)

// btopCmd represents the btop command group
var btopCmd = &cobra.Command{
	Use:   "btop",
	Short: "btop++ resource monitor",
	Long: `btop++ terminal resource monitor commands.

btop++ is a resource monitor that shows usage and stats for processor,
memory, disks, network and processes. Similar to htop but more feature-rich.

Examples:
  acorn sysadm btop status       # Show btop status
  acorn sysadm btop launch       # Start btop
  acorn sysadm btop themes       # List available themes
  acorn sysadm btop generate     # Generate config file`,
	Aliases: []string{"top", "htop"},
}

// btopStatusCmd shows btop status
var btopStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show btop status",
	Long: `Show btop installation status and configuration info.

Examples:
  acorn sysadm btop status
  acorn sysadm btop status -o json`,
	RunE: runBtopStatus,
}

// btopInstallCmd installs btop
var btopInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install btop",
	Long: `Install btop using Homebrew.

Examples:
  acorn sysadm btop install`,
	RunE: runBtopInstall,
}

// btopLaunchCmd launches btop
var btopLaunchCmd = &cobra.Command{
	Use:   "launch [args...]",
	Short: "Launch btop",
	Long: `Start btop resource monitor.

Examples:
  acorn sysadm btop launch
  acorn sysadm btop launch --utf-force`,
	Aliases: []string{"start", "open", "run"},
	RunE:    runBtopLaunch,
}

// btopGenerateCmd generates config
var btopGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate configuration file",
	Long: `Generate btop.conf configuration file with sensible defaults.

Creates ~/.config/btop/btop.conf with:
- Vim keybindings enabled
- Rounded corners
- Truecolor support
- CPU, memory, network, and process boxes

Examples:
  acorn sysadm btop generate
  acorn sysadm btop generate --dry-run`,
	Aliases: []string{"gen", "init"},
	RunE:    runBtopGenerate,
}

// Theme subcommands
var btopThemesCmd = &cobra.Command{
	Use:   "themes",
	Short: "Theme management",
	Long: `Manage btop themes.

Examples:
  acorn sysadm btop themes list
  acorn sysadm btop themes set catppuccin`,
}

var btopThemesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available themes",
	Long: `List all available btop themes.

Shows both system themes and custom themes in ~/.config/btop/themes/

Examples:
  acorn sysadm btop themes list`,
	Aliases: []string{"ls"},
	RunE:    runBtopThemesList,
}

var btopThemesSetCmd = &cobra.Command{
	Use:   "set <theme>",
	Short: "Set active theme",
	Long: `Set the active btop theme.

Examples:
  acorn sysadm btop themes set catppuccin
  acorn sysadm btop themes set Default`,
	Args: cobra.ExactArgs(1),
	RunE: runBtopThemesSet,
}

var btopThemesCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show current theme",
	Long: `Show the currently active btop theme.

Examples:
  acorn sysadm btop themes current`,
	RunE: runBtopThemesCurrent,
}

func init() {

	// Add subcommands
	btopCmd.AddCommand(btopStatusCmd)
	btopCmd.AddCommand(btopInstallCmd)
	btopCmd.AddCommand(btopLaunchCmd)
	btopCmd.AddCommand(btopGenerateCmd)
	btopCmd.AddCommand(configcmd.NewConfigRouter("btop"))

	// Theme subcommands
	btopCmd.AddCommand(btopThemesCmd)
	btopThemesCmd.AddCommand(btopThemesListCmd)
	btopThemesCmd.AddCommand(btopThemesSetCmd)
	btopThemesCmd.AddCommand(btopThemesCurrentCmd)

	// Persistent flags
	btopCmd.PersistentFlags().BoolVar(&btopDryRun, "dry-run", false,
		"Show what would be done without executing")
	btopCmd.PersistentFlags().BoolVarP(&btopVerbose, "verbose", "v", false,
		"Show verbose output")
}

func newBtopHelper() *btop.Helper {
	return btop.NewHelper(btopVerbose, btopDryRun)
}

func runBtopStatus(cmd *cobra.Command, args []string) error {
	helper := newBtopHelper()
	status := helper.GetStatus()

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("btop++ Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if status.Installed {
		fmt.Fprintf(os.Stdout, "%s btop++ installed: %s\n", output.Success("✓"), status.Version)
	} else {
		fmt.Fprintf(os.Stdout, "%s btop++ not found\n", output.Error("✗"))
		fmt.Fprintln(os.Stdout, "  Install: brew install btop")
		return nil
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "Config: %s\n", status.ConfigDir)
	if status.ConfigExists {
		fmt.Fprintf(os.Stdout, "%s Configuration found\n", output.Success("✓"))
	} else {
		fmt.Fprintf(os.Stdout, "%s No config (run 'acorn sysadm btop generate')\n", output.Warning("!"))
	}

	fmt.Fprintf(os.Stdout, "Themes: %d available\n", status.ThemeCount)

	// Show current theme if config exists
	if status.ConfigExists {
		currentTheme := helper.GetCurrentTheme()
		fmt.Fprintf(os.Stdout, "Active theme: %s\n", currentTheme)
	}

	return nil
}

func runBtopInstall(cmd *cobra.Command, args []string) error {
	helper := newBtopHelper()
	return helper.Install()
}

func runBtopLaunch(cmd *cobra.Command, args []string) error {
	helper := newBtopHelper()
	return helper.Launch(args...)
}

func runBtopGenerate(cmd *cobra.Command, args []string) error {
	helper := newBtopHelper()

	// Initialize config directories first
	if err := helper.InitConfig(); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Generating btop++ Configuration"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if err := helper.GenerateConfig(btopDryRun); err != nil {
		return err
	}

	if btopDryRun {
		fmt.Fprintf(os.Stdout, "\n%s Dry run complete. No files written.\n", output.Info("!"))
	} else {
		fmt.Fprintf(os.Stdout, "\n%s Configuration generated successfully\n", output.Success("✓"))
		fmt.Fprintln(os.Stdout, "\nFeatures enabled:")
		fmt.Fprintln(os.Stdout, "  - Vim keybindings (h,j,k,l,g,G)")
		fmt.Fprintln(os.Stdout, "  - Rounded corners")
		fmt.Fprintln(os.Stdout, "  - Truecolor support")
		fmt.Fprintln(os.Stdout, "  - Mouse support")
	}

	return nil
}

func runBtopThemesList(cmd *cobra.Command, args []string) error {
	helper := newBtopHelper()
	themes, err := helper.ListThemes()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(themes)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Available Themes"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if len(themes) == 0 {
		fmt.Fprintln(os.Stdout, "No themes found")
		fmt.Fprintln(os.Stdout, "Add themes to ~/.config/btop/themes/")
		return nil
	}

	currentTheme := helper.GetCurrentTheme()

	for _, t := range themes {
		marker := "  "
		if t.Name == currentTheme {
			marker = output.Success("→ ")
		}

		location := "system"
		if t.IsCustom {
			location = "custom"
		}

		fmt.Fprintf(os.Stdout, "%s%-20s (%s)\n", marker, t.Name, location)
	}

	return nil
}

func runBtopThemesSet(cmd *cobra.Command, args []string) error {
	helper := newBtopHelper()
	themeName := args[0]

	if err := helper.SetTheme(themeName); err != nil {
		return err
	}

	if btopDryRun {
		return nil
	}

	fmt.Fprintf(os.Stdout, "%s Theme set to: %s\n", output.Success("✓"), themeName)
	fmt.Fprintln(os.Stdout, "Restart btop for changes to take effect.")
	return nil
}

func runBtopThemesCurrent(cmd *cobra.Command, args []string) error {
	helper := newBtopHelper()
	currentTheme := helper.GetCurrentTheme()

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(map[string]string{"theme": currentTheme})
	}

	fmt.Fprintf(os.Stdout, "Current theme: %s\n", currentTheme)
	return nil
}

func init() {
	components.Register(&components.Registration{
		Name: "btop",
		RegisterCmd: func() *cobra.Command { return btopCmd },
	})
}
