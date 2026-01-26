package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/data/posting"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	postingDryRun     bool
	postingVerbose    bool
	postingCollection string
	postingEnvFile    string
)

// postingCmd represents the posting command group
var postingCmd = &cobra.Command{
	Use:   "posting",
	Short: "Terminal-based HTTP client with TUI",
	Long: `Terminal-based HTTP client with beautiful TUI interface.

Posting is a terminal HTTP client similar to Postman but runs in your terminal.

Examples:
  acorn data posting status           # Show status
  acorn data posting launch           # Launch posting
  acorn data posting locate           # Show file locations
  acorn data posting config           # Edit config file`,
	Aliases: []string{"post"},
}

// postingStatusCmd shows status
var postingStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Posting status",
	Long: `Show Posting installation status and configuration.

Examples:
  acorn data posting status
  acorn data posting status -o json`,
	RunE: runPostingStatus,
}

// postingLaunchCmd launches posting
var postingLaunchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch Posting",
	Long: `Launch Posting terminal HTTP client.

Examples:
  acorn data posting launch
  acorn data posting launch --collection ~/projects/api
  acorn data posting launch --env .env.local`,
	Aliases: []string{"open"},
	RunE:    runPostingLaunch,
}

// postingLocateCmd shows file locations
var postingLocateCmd = &cobra.Command{
	Use:   "locate [type]",
	Short: "Show file locations",
	Long: `Show Posting file and directory locations.

Types: config, themes, data, all (default)

Examples:
  acorn data posting locate
  acorn data posting locate config
  acorn data posting locate -o json`,
	RunE: runPostingLocate,
}

// postingConfigCmd opens config file
var postingConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Edit config file",
	Long: `Open Posting config file in the default editor.

Uses $EDITOR environment variable, defaults to vim.

Examples:
  acorn data posting config`,
	Aliases: []string{"edit"},
	RunE:    runPostingConfig,
}

// postingCollectionsCmd lists collections
var postingCollectionsCmd = &cobra.Command{
	Use:   "collections",
	Short: "List collections",
	Long: `List available Posting collection directories.

Examples:
  acorn data posting collections
  acorn data posting collections -o json`,
	Aliases: []string{"ls"},
	RunE:    runPostingCollections,
}

// postingInstallCmd installs posting
var postingInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Posting",
	Long: `Install Posting terminal HTTP client via uv.

Requires uv package manager. If uv is not installed, run install-uv first.

Examples:
  acorn data posting install`,
	RunE: runPostingInstall,
}

// postingInstallUVCmd installs uv
var postingInstallUVCmd = &cobra.Command{
	Use:   "install-uv",
	Short: "Install uv package manager",
	Long: `Install uv package manager (required for Posting).

Examples:
  acorn data posting install-uv`,
	RunE: runPostingInstallUV,
}

// postingUpgradeCmd upgrades posting
var postingUpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade Posting",
	Long: `Upgrade Posting to the latest version.

Examples:
  acorn data posting upgrade`,
	RunE: runPostingUpgrade,
}

// postingUninstallCmd uninstalls posting
var postingUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall Posting",
	Long: `Uninstall Posting terminal HTTP client.

Examples:
  acorn data posting uninstall`,
	Aliases: []string{"remove"},
	RunE:    runPostingUninstall,
}

func init() {
	dataCmd.AddCommand(postingCmd)

	// Add subcommands
	postingCmd.AddCommand(postingStatusCmd)
	postingCmd.AddCommand(postingLaunchCmd)
	postingCmd.AddCommand(postingLocateCmd)
	postingCmd.AddCommand(postingConfigCmd)
	postingCmd.AddCommand(postingCollectionsCmd)
	postingCmd.AddCommand(postingInstallCmd)
	postingCmd.AddCommand(postingInstallUVCmd)
	postingCmd.AddCommand(postingUpgradeCmd)
	postingCmd.AddCommand(postingUninstallCmd)

	// Launch flags
	postingLaunchCmd.Flags().StringVarP(&postingCollection, "collection", "c", "",
		"Collection directory to open")
	postingLaunchCmd.Flags().StringVarP(&postingEnvFile, "env", "e", "",
		"Environment file to load")

	// Persistent flags
	postingCmd.PersistentFlags().BoolVar(&postingDryRun, "dry-run", false,
		"Show what would be done without executing")
	postingCmd.PersistentFlags().BoolVarP(&postingVerbose, "verbose", "v", false,
		"Show verbose output")
}

func newPostingHelper() *posting.Helper {
	return posting.NewHelper(postingVerbose, postingDryRun)
}

func runPostingStatus(cmd *cobra.Command, args []string) error {
	helper := newPostingHelper()
	status := helper.GetStatus()

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Posting Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if status.Installed {
		version := status.Version
		if version == "" {
			version = "detected"
		}
		fmt.Fprintf(os.Stdout, "%s Posting installed: %s\n", output.Success("✓"), version)
	} else {
		fmt.Fprintf(os.Stdout, "%s Posting not found\n", output.Error("✗"))
		fmt.Fprintln(os.Stdout, "  Install: uv tool install posting --python 3.13")
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "Theme:  %s\n", status.Theme)
	fmt.Fprintf(os.Stdout, "Layout: %s\n", status.Layout)

	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "Config: %s\n", status.ConfigPath)
	fmt.Fprintf(os.Stdout, "Data:   %s\n", status.DataPath)

	return nil
}

func runPostingLaunch(cmd *cobra.Command, args []string) error {
	helper := newPostingHelper()
	return helper.Launch(postingCollection, postingEnvFile)
}

func runPostingLocate(cmd *cobra.Command, args []string) error {
	helper := newPostingHelper()
	locations := helper.GetLocations()

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(locations)
	}

	// Filter by type if specified
	locationType := "all"
	if len(args) > 0 {
		locationType = args[0]
	}

	switch locationType {
	case "config":
		fmt.Fprintln(os.Stdout, locations.Config)
	case "themes":
		fmt.Fprintln(os.Stdout, locations.Themes)
	case "data":
		fmt.Fprintln(os.Stdout, locations.Data)
	default:
		fmt.Fprintf(os.Stdout, "%s\n", output.Info("Posting Locations"))
		fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Fprintf(os.Stdout, "Config: %s\n", locations.Config)
		fmt.Fprintf(os.Stdout, "Themes: %s\n", locations.Themes)
		fmt.Fprintf(os.Stdout, "Data:   %s\n", locations.Data)
	}

	return nil
}

func runPostingConfig(cmd *cobra.Command, args []string) error {
	helper := newPostingHelper()
	return helper.OpenConfig()
}

func runPostingCollections(cmd *cobra.Command, args []string) error {
	helper := newPostingHelper()
	collections, err := helper.ListCollections()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(collections)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Collections"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if len(collections) == 0 {
		fmt.Fprintln(os.Stdout, "No collections found")
		return nil
	}

	for _, c := range collections {
		fmt.Fprintf(os.Stdout, "  %s\n", c)
	}

	return nil
}

func runPostingInstall(cmd *cobra.Command, args []string) error {
	helper := newPostingHelper()

	if err := helper.Install(); err != nil {
		return err
	}

	if !postingDryRun {
		fmt.Fprintf(os.Stdout, "%s Posting installed\n", output.Success("✓"))
		fmt.Fprintln(os.Stdout, "Run 'posting' to start the HTTP client")
	}

	return nil
}

func runPostingInstallUV(cmd *cobra.Command, args []string) error {
	helper := newPostingHelper()

	if err := helper.InstallUV(); err != nil {
		return err
	}

	if !postingDryRun {
		fmt.Fprintf(os.Stdout, "%s uv installed\n", output.Success("✓"))
		fmt.Fprintln(os.Stdout, "Restart your shell or run: source ~/.local/bin/env")
	}

	return nil
}

func runPostingUpgrade(cmd *cobra.Command, args []string) error {
	helper := newPostingHelper()

	if err := helper.Upgrade(); err != nil {
		return err
	}

	if !postingDryRun {
		fmt.Fprintf(os.Stdout, "%s Posting upgraded\n", output.Success("✓"))
	}

	return nil
}

func runPostingUninstall(cmd *cobra.Command, args []string) error {
	helper := newPostingHelper()

	if err := helper.Uninstall(); err != nil {
		return err
	}

	if !postingDryRun {
		fmt.Fprintf(os.Stdout, "%s Posting uninstalled\n", output.Success("✓"))
	}

	return nil
}
