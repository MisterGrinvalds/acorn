package cmd

import (
	"github.com/mistergrinvalds/acorn/internal/components"
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/goland"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/mistergrinvalds/acorn/internal/utils/configcmd"
	"github.com/spf13/cobra"
)

var (
	golandDryRun  bool
	golandVerbose bool
	golandLine    int
)

// golandCmd represents the goland command group
var golandCmd = &cobra.Command{
	Use:   "goland",
	Short: "JetBrains GoLand IDE",
	Long: `JetBrains GoLand IDE commands.

Manage GoLand installation, launch, and command-line integration.

Examples:
  acorn ide goland status      # Show status
  acorn ide goland launch      # Open GoLand
  acorn ide goland open .      # Open current directory
  acorn ide goland cli-link    # Create CLI launcher`,
	Aliases: []string{"gl"},
}

// golandStatusCmd shows status
var golandStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show GoLand status",
	Long: `Show GoLand installation and configuration status.

Examples:
  acorn ide goland status
  acorn ide goland status -o json`,
	RunE: runGolandStatus,
}

// golandLaunchCmd launches GoLand
var golandLaunchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch GoLand",
	Long: `Launch GoLand application.

Examples:
  acorn ide goland launch`,
	RunE: runGolandLaunch,
}

// golandOpenCmd opens a project or file
var golandOpenCmd = &cobra.Command{
	Use:   "open <path>",
	Short: "Open project or file",
	Long: `Open a project directory or file in GoLand.

Examples:
  acorn ide goland open .
  acorn ide goland open ~/projects/myapp
  acorn ide goland open main.go --line 42`,
	Args: cobra.ExactArgs(1),
	RunE: runGolandOpen,
}

// golandCLILinkCmd creates CLI launcher
var golandCLILinkCmd = &cobra.Command{
	Use:   "cli-link",
	Short: "Create CLI launcher",
	Long: `Create a command-line launcher for GoLand.

Creates a symlink at /usr/local/bin/goland pointing to the app bundle.

Examples:
  acorn ide goland cli-link
  sudo acorn ide goland cli-link`,
	RunE: runGolandCLILink,
}

func init() {

	// Add subcommands
	golandCmd.AddCommand(golandStatusCmd)
	golandCmd.AddCommand(golandLaunchCmd)
	golandCmd.AddCommand(golandOpenCmd)
	golandCmd.AddCommand(golandCLILinkCmd)
	golandCmd.AddCommand(configcmd.NewConfigRouter("goland"))

	// Open flags
	golandOpenCmd.Flags().IntVar(&golandLine, "line", 0, "Line number to jump to")

	// Persistent flags
	golandCmd.PersistentFlags().BoolVar(&golandDryRun, "dry-run", false,
		"Show what would be done without executing")
	golandCmd.PersistentFlags().BoolVarP(&golandVerbose, "verbose", "v", false,
		"Show verbose output")
}

func newGolandHelper() *goland.Helper {
	return goland.NewHelper(golandVerbose, golandDryRun)
}

func runGolandStatus(cmd *cobra.Command, args []string) error {
	helper := newGolandHelper()
	status := helper.GetStatus()

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("GoLand Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if status.Installed {
		version := status.Version
		if version == "" {
			version = "detected"
		}
		fmt.Fprintf(os.Stdout, "%s GoLand installed: %s\n", output.Success("✓"), version)
		fmt.Fprintf(os.Stdout, "  Path: %s\n", status.AppPath)
	} else {
		fmt.Fprintf(os.Stdout, "%s GoLand not found\n", output.Error("✗"))
		fmt.Fprintln(os.Stdout, "  Install via JetBrains Toolbox or: brew install --cask goland")
		return nil
	}

	fmt.Fprintln(os.Stdout)

	if status.HasCLI {
		fmt.Fprintf(os.Stdout, "%s CLI launcher: %s\n", output.Success("✓"), status.CLIPath)
	} else {
		fmt.Fprintf(os.Stdout, "%s CLI launcher not configured\n", output.Warning("!"))
		fmt.Fprintln(os.Stdout, "  Create: acorn ide goland cli-link")
	}

	if status.ConfigDir != "" {
		fmt.Fprintf(os.Stdout, "\nConfig: %s\n", status.ConfigDir)
	}

	return nil
}

func runGolandLaunch(cmd *cobra.Command, args []string) error {
	helper := newGolandHelper()
	if err := helper.Launch(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s GoLand launched\n", output.Success("✓"))
	return nil
}

func runGolandOpen(cmd *cobra.Command, args []string) error {
	helper := newGolandHelper()

	path := args[0]

	// Check if it's a file with line number
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("path not found: %s", path)
	}

	if info.IsDir() {
		if err := helper.OpenProject(path); err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "%s Opening project in GoLand\n", output.Success("✓"))
	} else {
		if err := helper.OpenFile(path, golandLine); err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "%s Opening file in GoLand\n", output.Success("✓"))
	}

	return nil
}

func runGolandCLILink(cmd *cobra.Command, args []string) error {
	helper := newGolandHelper()
	if err := helper.CreateCLILink(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s CLI launcher created\n", output.Success("✓"))
	return nil
}

func init() {
	components.Register(&components.Registration{
		Name: "goland",
		RegisterCmd: func() *cobra.Command { return golandCmd },
	})
}
