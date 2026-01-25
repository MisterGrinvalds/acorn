package cmd

import (
	"fmt"
	"os"
	"strings"

	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/mistergrinvalds/acorn/internal/utils/tools"
	"github.com/spf13/cobra"
)

var (
	toolsDryRun  bool
	toolsVerbose bool
)

// toolsCmd represents the tools command group
var toolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "Manage development tools",
	Long: `Check, install, and update development tools.

This command helps you manage the various development tools installed on your
system. It can check what's installed, show versions, identify missing tools,
and help keep everything up to date.

Examples:
  acorn tools status              # Show all tool versions
  acorn tools check git go node   # Check specific tools
  acorn tools missing             # Show uninstalled tools
  acorn tools which git           # Show path and version`,
	Aliases: []string{"tool"},
}

// toolsStatusCmd shows all tool versions and status
var toolsStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show all tool versions and status",
	Long: `Display the status of all tracked development tools.

Shows each tool grouped by category (System, Languages, Cloud, Development)
with installation status and version information.

Examples:
  acorn tools status
  acorn tools status -o json
  acorn tools status -o yaml`,
	RunE: runToolsStatus,
}

// toolsListCmd lists all known tools
var toolsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all trackable tools",
	Long: `List all tools that acorn can track.

Shows the complete registry of tools including their category, description,
and installation hints.

Examples:
  acorn tools list
  acorn tools list -o json`,
	Aliases: []string{"ls"},
	RunE:    runToolsList,
}

// toolsCheckCmd checks specific tools
var toolsCheckCmd = &cobra.Command{
	Use:   "check [tool...]",
	Short: "Check if specific tools are installed",
	Long: `Check the installation status of specific tools.

If no tools are specified, checks all known tools.

Examples:
  acorn tools check git go node
  acorn tools check kubectl helm terraform`,
	RunE:              runToolsCheck,
	ValidArgsFunction: completeToolNames,
}

// toolsMissingCmd shows missing tools
var toolsMissingCmd = &cobra.Command{
	Use:   "missing",
	Short: "Show tools that are not installed",
	Long: `List all tracked tools that are not currently installed.

Shows missing tools with hints on how to install them.

Examples:
  acorn tools missing
  acorn tools missing -o json`,
	RunE: runToolsMissing,
}

// toolsWhichCmd shows tool location and version
var toolsWhichCmd = &cobra.Command{
	Use:   "which <tool>",
	Short: "Show tool location and version",
	Long: `Display the path and version of an installed tool.

Similar to the 'which' command but with additional version information.

Examples:
  acorn tools which git
  acorn tools which go`,
	Args:              cobra.ExactArgs(1),
	RunE:              runToolsWhich,
	ValidArgsFunction: completeToolNames,
}

// toolsUpdateCmd updates system packages
var toolsUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update tools via package manager",
	Long: `Update system packages using the detected package manager.

Detects and uses brew (macOS), apt-get, dnf, or pacman.

Examples:
  acorn tools update
  acorn tools update --dry-run`,
	RunE: runToolsUpdate,
}

// toolsInstallCmd installs a tool
var toolsInstallCmd = &cobra.Command{
	Use:   "install <tool>",
	Short: "Install a tool from the registry",
	Long: `Install a tool using the method specified in the registry.

Examples:
  acorn tools install bat
  acorn tools install eza --dry-run`,
	Args:              cobra.ExactArgs(1),
	RunE:              runToolsInstall,
	ValidArgsFunction: completeMissingToolNames,
}

// toolsUpgradeBashCmd upgrades bash on macOS
var toolsUpgradeBashCmd = &cobra.Command{
	Use:   "upgrade-bash",
	Short: "Upgrade to modern bash (macOS only)",
	Long: `Upgrade to a modern version of bash via Homebrew.

macOS ships with an ancient bash 3.2 due to licensing. This command installs
a modern bash via Homebrew and optionally sets it as your default shell.

Examples:
  acorn tools upgrade-bash`,
	RunE: runToolsUpgradeBash,
}

func init() {
	rootCmd.AddCommand(toolsCmd)

	// Add subcommands
	toolsCmd.AddCommand(toolsStatusCmd)
	toolsCmd.AddCommand(toolsListCmd)
	toolsCmd.AddCommand(toolsCheckCmd)
	toolsCmd.AddCommand(toolsMissingCmd)
	toolsCmd.AddCommand(toolsWhichCmd)
	toolsCmd.AddCommand(toolsUpdateCmd)
	toolsCmd.AddCommand(toolsInstallCmd)
	toolsCmd.AddCommand(toolsUpgradeBashCmd)

	// Flags for update/install commands
	toolsUpdateCmd.Flags().BoolVar(&toolsDryRun, "dry-run", false, "Show what would be done without executing")
	toolsUpdateCmd.Flags().BoolVarP(&toolsVerbose, "verbose", "v", false, "Show verbose output")
	toolsInstallCmd.Flags().BoolVar(&toolsDryRun, "dry-run", false, "Show what would be done without executing")
	toolsInstallCmd.Flags().BoolVarP(&toolsVerbose, "verbose", "v", false, "Show verbose output")
}

func runToolsStatus(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	checker := tools.NewChecker()
	result := checker.CheckAll()

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(result)
	}

	// Table format with colored status
	for _, cat := range result.Categories {
		fmt.Fprintf(os.Stdout, "\n%s\n", output.Info(cat.Name))
		fmt.Fprintln(os.Stdout, strings.Repeat("-", len(cat.Name)))

		for _, tool := range cat.Tools {
			status := output.Success("✓")
			version := tool.Version
			if !tool.Installed {
				status = output.Error("✗")
				version = "not installed"
			}
			fmt.Fprintf(os.Stdout, "  %s %-15s %s\n", status, tool.Name, version)
		}
	}

	// Summary
	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "Total: %d tools (%s installed, %s missing)\n",
		result.Summary.Total,
		output.Success(fmt.Sprintf("%d", result.Summary.Installed)),
		output.Error(fmt.Sprintf("%d", result.Summary.Missing)))

	return nil
}

func runToolsList(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	registry := tools.DefaultRegistry()

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(registry)
	}

	// Table format
	table := output.NewTable("NAME", "CATEGORY", "DESCRIPTION", "INSTALL")
	for _, def := range registry {
		table.AddRow(def.Name, def.Category, def.Description, def.InstallHint)
	}
	table.Render(os.Stdout)

	fmt.Fprintf(os.Stdout, "\nTotal: %d tools\n", len(registry))
	return nil
}

func runToolsCheck(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	checker := tools.NewChecker()

	var results []tools.ToolStatus
	if len(args) == 0 {
		// Check all tools
		statusResult := checker.CheckAll()
		for _, cat := range statusResult.Categories {
			results = append(results, cat.Tools...)
		}
	} else {
		// Check specific tools
		results = checker.CheckTools(args)
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(results)
	}

	// Table format
	for _, tool := range results {
		status := output.Success("✓")
		info := tool.Version
		if tool.Path != "" {
			info = tool.Path
			if tool.Version != "" {
				info += " (" + tool.Version + ")"
			}
		}
		if !tool.Installed {
			status = output.Error("✗")
			info = "not installed"
		}
		fmt.Fprintf(os.Stdout, "%s %s: %s\n", status, tool.Name, info)
	}

	return nil
}

func runToolsMissing(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	checker := tools.NewChecker()
	missing := checker.GetMissing()

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(missing)
	}

	if len(missing) == 0 {
		fmt.Fprintln(os.Stdout, output.Success("All tracked tools are installed!"))
		return nil
	}

	// Table format
	fmt.Fprintf(os.Stdout, "Missing %d tools:\n\n", len(missing))
	for _, tool := range missing {
		def, found := tools.FindTool(tool.Name)
		hint := ""
		if found {
			hint = def.InstallHint
		}
		fmt.Fprintf(os.Stdout, "  %s %s\n", output.Error("✗"), tool.Name)
		if hint != "" {
			fmt.Fprintf(os.Stdout, "    Install: %s\n", output.Info(hint))
		}
	}

	return nil
}

func runToolsWhich(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	name := args[0]
	checker := tools.NewChecker()
	status := checker.CheckTool(name)

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	if !status.Installed {
		fmt.Fprintf(os.Stdout, "%s %s not found\n", output.Error("✗"), name)
		if def, found := tools.FindTool(name); found {
			fmt.Fprintf(os.Stdout, "  Install: %s\n", output.Info(def.InstallHint))
		}
		return nil
	}

	fmt.Fprintf(os.Stdout, "%s found:\n", output.Success(name))
	fmt.Fprintf(os.Stdout, "  Location: %s\n", status.Path)
	if status.Version != "" {
		fmt.Fprintf(os.Stdout, "  Version:  %s\n", status.Version)
	}
	if status.Category != "" && status.Category != "Unknown" {
		fmt.Fprintf(os.Stdout, "  Category: %s\n", status.Category)
	}

	return nil
}

func runToolsUpdate(cmd *cobra.Command, args []string) error {
	updater := tools.NewUpdater(toolsDryRun, toolsVerbose)

	pm := updater.DetectPackageManager()
	if pm == tools.PMNone {
		return fmt.Errorf("no supported package manager found")
	}

	fmt.Fprintf(os.Stdout, "Detected package manager: %s\n", output.Info(string(pm)))
	fmt.Fprintln(os.Stdout, "Updating system packages...")

	return updater.UpgradeSystem()
}

func runToolsInstall(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Check if already installed
	checker := tools.NewChecker()
	status := checker.CheckTool(name)
	if status.Installed {
		fmt.Fprintf(os.Stdout, "%s is already installed at %s\n",
			output.Success(name), status.Path)
		return nil
	}

	// Check if in registry
	def, found := tools.FindTool(name)
	if !found {
		return fmt.Errorf("tool %q not in registry", name)
	}

	fmt.Fprintf(os.Stdout, "Installing %s...\n", output.Info(name))
	fmt.Fprintf(os.Stdout, "Command: %s\n", def.InstallHint)

	updater := tools.NewUpdater(toolsDryRun, toolsVerbose)
	return updater.InstallTool(name)
}

func runToolsUpgradeBash(cmd *cobra.Command, args []string) error {
	updater := tools.NewUpdater(toolsDryRun, toolsVerbose)
	return updater.UpgradeBash()
}

// completeToolNames provides completion for tool names
func completeToolNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	names := tools.ToolNames()
	var matches []string
	for _, name := range names {
		if strings.HasPrefix(name, toComplete) {
			matches = append(matches, name)
		}
	}
	return matches, cobra.ShellCompDirectiveNoFileComp
}

// completeMissingToolNames provides completion for missing tool names
func completeMissingToolNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	checker := tools.NewChecker()
	missing := checker.GetMissing()
	var matches []string
	for _, tool := range missing {
		if strings.HasPrefix(tool.Name, toComplete) {
			matches = append(matches, tool.Name)
		}
	}
	return matches, cobra.ShellCompDirectiveNoFileComp
}
