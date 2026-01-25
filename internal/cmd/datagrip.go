package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/data/datagrip"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	datagripDryRun  bool
	datagripVerbose bool
	datagripLine    int
)

// datagripCmd represents the datagrip command group
var datagripCmd = &cobra.Command{
	Use:   "datagrip",
	Short: "JetBrains DataGrip database IDE",
	Long: `JetBrains DataGrip database IDE commands.

DataGrip is a cross-platform IDE for databases and SQL.
Provides commands for launching, opening projects, and status.

Examples:
  acorn data datagrip status       # Show DataGrip status
  acorn data datagrip launch       # Start DataGrip
  acorn data datagrip open .       # Open current directory`,
	Aliases: []string{"dg"},
}

// datagripStatusCmd shows DataGrip status
var datagripStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show DataGrip status",
	Long: `Show DataGrip installation status and configuration info.

Examples:
  acorn data datagrip status
  acorn data datagrip status -o json`,
	RunE: runDatagripStatus,
}

// datagripLaunchCmd launches DataGrip
var datagripLaunchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch DataGrip",
	Long: `Start DataGrip IDE.

Examples:
  acorn data datagrip launch`,
	Aliases: []string{"start", "run"},
	RunE:    runDatagripLaunch,
}

// datagripOpenCmd opens a project or file
var datagripOpenCmd = &cobra.Command{
	Use:   "open <path>",
	Short: "Open project or file",
	Long: `Open a project directory or file in DataGrip.

Examples:
  acorn data datagrip open .
  acorn data datagrip open ~/projects/mydb
  acorn data datagrip open schema.sql --line 42`,
	Args: cobra.ExactArgs(1),
	RunE: runDatagripOpen,
}

// datagripCLICmd manages CLI launcher
var datagripCLICmd = &cobra.Command{
	Use:   "cli",
	Short: "CLI launcher management",
	Long: `Manage DataGrip command-line launcher.

Examples:
  acorn data datagrip cli create`,
}

var datagripCLICreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create CLI launcher",
	Long: `Create command-line launcher for DataGrip.

Creates a symlink at /usr/local/bin/datagrip.
May require sudo for permission.

Examples:
  acorn data datagrip cli create
  sudo acorn data datagrip cli create`,
	RunE: runDatagripCLICreate,
}

func init() {
	dataCmd.AddCommand(datagripCmd)

	// Add subcommands
	datagripCmd.AddCommand(datagripStatusCmd)
	datagripCmd.AddCommand(datagripLaunchCmd)
	datagripCmd.AddCommand(datagripOpenCmd)

	// CLI subcommands
	datagripCmd.AddCommand(datagripCLICmd)
	datagripCLICmd.AddCommand(datagripCLICreateCmd)

	// Flags
	datagripOpenCmd.Flags().IntVar(&datagripLine, "line", 0, "Line number to open at")

	// Persistent flags
	datagripCmd.PersistentFlags().BoolVar(&datagripDryRun, "dry-run", false,
		"Show what would be done without executing")
	datagripCmd.PersistentFlags().BoolVarP(&datagripVerbose, "verbose", "v", false,
		"Show verbose output")
}

func newDatagripHelper() *datagrip.Helper {
	return datagrip.NewHelper(datagripVerbose, datagripDryRun)
}

func runDatagripStatus(cmd *cobra.Command, args []string) error {
	helper := newDatagripHelper()
	status := helper.GetStatus()

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("DataGrip Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if status.Installed {
		fmt.Fprintf(os.Stdout, "%s DataGrip installed", output.Success("✓"))
		if status.Version != "" {
			fmt.Fprintf(os.Stdout, ": %s", status.Version)
		}
		fmt.Fprintln(os.Stdout)
	} else {
		fmt.Fprintf(os.Stdout, "%s DataGrip not found\n", output.Error("✗"))
		fmt.Fprintln(os.Stdout, "  Install from: https://www.jetbrains.com/datagrip/")
		fmt.Fprintln(os.Stdout, "  Or use JetBrains Toolbox")
		return nil
	}

	if status.AppPath != "" {
		fmt.Fprintf(os.Stdout, "App: %s\n", status.AppPath)
	}

	if status.HasCLI {
		fmt.Fprintf(os.Stdout, "%s CLI available: %s\n", output.Success("✓"), status.CLIPath)
	} else {
		fmt.Fprintf(os.Stdout, "%s CLI not configured (run 'acorn data datagrip cli create')\n", output.Warning("!"))
	}

	if status.ConfigDir != "" {
		fmt.Fprintf(os.Stdout, "Config: %s\n", status.ConfigDir)
	}

	return nil
}

func runDatagripLaunch(cmd *cobra.Command, args []string) error {
	helper := newDatagripHelper()
	if err := helper.Launch(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s DataGrip launched\n", output.Success("✓"))
	return nil
}

func runDatagripOpen(cmd *cobra.Command, args []string) error {
	helper := newDatagripHelper()
	path := args[0]

	if datagripLine > 0 {
		if err := helper.OpenFile(path, datagripLine); err != nil {
			return err
		}
	} else {
		if err := helper.OpenProject(path); err != nil {
			return err
		}
	}

	fmt.Fprintf(os.Stdout, "%s Opening in DataGrip: %s\n", output.Success("✓"), path)
	return nil
}

func runDatagripCLICreate(cmd *cobra.Command, args []string) error {
	helper := newDatagripHelper()
	return helper.CreateCLILink()
}
