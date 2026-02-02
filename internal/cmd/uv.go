package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/programming/uv"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/mistergrinvalds/acorn/internal/utils/configcmd"
	"github.com/spf13/cobra"
)

var (
	uvDryRun  bool
	uvVerbose bool
)

// uvCmd represents the uv command group
var uvCmd = &cobra.Command{
	Use:   "uv",
	Short: "UV Python package manager",
	Long: `UV (Astral) Python package manager commands.

UV is an extremely fast Python package and project manager.
Provides commands for Python versions, tools, projects, and cache management.

Examples:
  acorn uv status             # Show UV status
  acorn uv python list        # List Python versions
  acorn uv python install 3.12
  acorn uv tool list          # List installed tools
  acorn uv tool install ruff  # Install a tool
  acorn uv cache clean        # Clean cache`,
}

// uvStatusCmd shows UV status
var uvStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show UV status",
	Long: `Show UV installation status and environment info.

Examples:
  acorn uv status
  acorn uv status -o json`,
	RunE: runUvStatus,
}

// uvInstallCmd installs UV
var uvInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install UV",
	Long: `Install UV using the official installer.

Examples:
  acorn uv install`,
	RunE: runUvInstall,
}

// uvUpdateCmd updates UV
var uvUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update UV to latest version",
	Long: `Update UV to the latest version.

Examples:
  acorn uv update`,
	Aliases: []string{"upgrade", "self-update"},
	RunE:    runUvUpdate,
}

// Python subcommands
var uvPythonCmd = &cobra.Command{
	Use:   "python",
	Short: "Python version management",
	Long: `Manage Python versions with UV.

Examples:
  acorn uv python list
  acorn uv python install 3.12`,
}

var uvPythonListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed Python versions",
	Long: `List Python versions installed by UV.

Examples:
  acorn uv python list`,
	Aliases: []string{"ls"},
	RunE:    runUvPythonList,
}

var uvPythonInstallCmd = &cobra.Command{
	Use:   "install <version>",
	Short: "Install a Python version",
	Long: `Install a Python version using UV.

Examples:
  acorn uv python install 3.12
  acorn uv python install 3.11.7`,
	Args: cobra.ExactArgs(1),
	RunE: runUvPythonInstall,
}

// Tool subcommands
var uvToolCmd = &cobra.Command{
	Use:   "tool",
	Short: "Tool management",
	Long: `Manage Python tools with UV.

Examples:
  acorn uv tool list
  acorn uv tool install ruff`,
}

var uvToolListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed tools",
	Long: `List tools installed with UV.

Examples:
  acorn uv tool list`,
	Aliases: []string{"ls"},
	RunE:    runUvToolList,
}

var uvToolInstallCmd = &cobra.Command{
	Use:   "install <tool>",
	Short: "Install a tool",
	Long: `Install a Python tool using UV.

Examples:
  acorn uv tool install ruff
  acorn uv tool install black`,
	Args: cobra.MinimumNArgs(1),
	RunE: runUvToolInstall,
}

var uvToolRunCmd = &cobra.Command{
	Use:   "run <tool> [args...]",
	Short: "Run a tool (uvx)",
	Long: `Run a Python tool without installing (uvx).

Examples:
  acorn uv tool run ruff check .
  acorn uv tool run black --check .`,
	Args: cobra.MinimumNArgs(1),
	RunE: runUvToolRun,
}

// Cache subcommands
var uvCacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Cache management",
	Long: `Manage UV cache.

Examples:
  acorn uv cache info
  acorn uv cache clean
  acorn uv cache prune`,
}

var uvCacheInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show cache information",
	Long: `Show UV cache directory and size.

Examples:
  acorn uv cache info`,
	RunE: runUvCacheInfo,
}

var uvCacheCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean cache",
	Long: `Remove all cached data.

Examples:
  acorn uv cache clean`,
	RunE: runUvCacheClean,
}

var uvCachePruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Prune cache",
	Long: `Remove unused cached data.

Examples:
  acorn uv cache prune`,
	RunE: runUvCachePrune,
}

// Project commands
var uvInitCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Initialize a new project",
	Long: `Initialize a new UV Python project.

Creates pyproject.toml and basic project structure.

Examples:
  acorn uv init
  acorn uv init myproject`,
	Args: cobra.MaximumNArgs(1),
	RunE: runUvInit,
}

var uvSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync project dependencies",
	Long: `Sync dependencies from pyproject.toml.

Examples:
  acorn uv sync`,
	RunE: runUvSync,
}

var uvAddCmd = &cobra.Command{
	Use:   "add <package> [packages...]",
	Short: "Add packages to project",
	Long: `Add packages to the project.

Examples:
  acorn uv add fastapi
  acorn uv add requests pytest`,
	Args: cobra.MinimumNArgs(1),
	RunE: runUvAdd,
}

var uvRemoveCmd = &cobra.Command{
	Use:   "remove <package> [packages...]",
	Short: "Remove packages from project",
	Long: `Remove packages from the project.

Examples:
  acorn uv remove fastapi
  acorn uv remove requests pytest`,
	Aliases: []string{"rm"},
	Args:    cobra.MinimumNArgs(1),
	RunE:    runUvRemove,
}

var uvRunCmd = &cobra.Command{
	Use:   "run <command> [args...]",
	Short: "Run command in project environment",
	Long: `Run a command in the project environment.

Examples:
  acorn uv run python main.py
  acorn uv run pytest`,
	Args: cobra.MinimumNArgs(1),
	RunE: runUvRun,
}

var uvLockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Lock project dependencies",
	Long: `Lock project dependencies to uv.lock.

Examples:
  acorn uv lock`,
	RunE: runUvLock,
}

var uvTreeCmd = &cobra.Command{
	Use:   "tree",
	Short: "Show dependency tree",
	Long: `Show the project dependency tree.

Examples:
  acorn uv tree`,
	RunE: runUvTree,
}

var uvVenvCmd = &cobra.Command{
	Use:   "venv [name]",
	Short: "Create virtual environment",
	Long: `Create a virtual environment.

Examples:
  acorn uv venv
  acorn uv venv .venv`,
	Args: cobra.MaximumNArgs(1),
	RunE: runUvVenv,
}

func init() {
	programmingCmd.AddCommand(uvCmd)

	// Add subcommands
	uvCmd.AddCommand(uvStatusCmd)
	uvCmd.AddCommand(uvInstallCmd)
	uvCmd.AddCommand(uvUpdateCmd)
	uvCmd.AddCommand(uvInitCmd)
	uvCmd.AddCommand(uvSyncCmd)
	uvCmd.AddCommand(uvAddCmd)
	uvCmd.AddCommand(uvRemoveCmd)
	uvCmd.AddCommand(uvRunCmd)
	uvCmd.AddCommand(uvLockCmd)
	uvCmd.AddCommand(uvTreeCmd)
	uvCmd.AddCommand(uvVenvCmd)

	// Python subcommands
	uvCmd.AddCommand(uvPythonCmd)
	uvPythonCmd.AddCommand(uvPythonListCmd)
	uvPythonCmd.AddCommand(uvPythonInstallCmd)

	// Tool subcommands
	uvCmd.AddCommand(uvToolCmd)
	uvToolCmd.AddCommand(uvToolListCmd)
	uvToolCmd.AddCommand(uvToolInstallCmd)
	uvToolCmd.AddCommand(uvToolRunCmd)

	// Cache subcommands
	uvCmd.AddCommand(uvCacheCmd)
	uvCmd.AddCommand(configcmd.NewConfigRouter("uv"))
	uvCacheCmd.AddCommand(uvCacheInfoCmd)
	uvCacheCmd.AddCommand(uvCacheCleanCmd)
	uvCacheCmd.AddCommand(uvCachePruneCmd)

	// Persistent flags
	uvCmd.PersistentFlags().BoolVar(&uvDryRun, "dry-run", false,
		"Show what would be done without executing")
	uvCmd.PersistentFlags().BoolVarP(&uvVerbose, "verbose", "v", false,
		"Show verbose output")
}

func newUvHelper() *uv.Helper {
	return uv.NewHelper(uvVerbose, uvDryRun)
}

func runUvStatus(cmd *cobra.Command, args []string) error {
	helper := newUvHelper()
	status := helper.GetStatus()

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("UV Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if status.Installed {
		fmt.Fprintf(os.Stdout, "%s UV installed: %s\n", output.Success("✓"), status.Version)
	} else {
		fmt.Fprintf(os.Stdout, "%s UV not found\n", output.Error("✗"))
		fmt.Fprintln(os.Stdout, "  Install: curl -LsSf https://astral.sh/uv/install.sh | sh")
		return nil
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "Cache: %s (%s)\n", status.CacheDir, status.CacheSize)

	if status.PythonPath != "" {
		fmt.Fprintf(os.Stdout, "Python: %s\n", status.PythonPath)
	}

	if status.InProject {
		fmt.Fprintf(os.Stdout, "\n%s In project: %s\n", output.Success("✓"), status.ProjectName)
	}

	return nil
}

func runUvInstall(cmd *cobra.Command, args []string) error {
	helper := newUvHelper()
	return helper.Install()
}

func runUvUpdate(cmd *cobra.Command, args []string) error {
	helper := newUvHelper()
	return helper.SelfUpdate()
}

func runUvPythonList(cmd *cobra.Command, args []string) error {
	helper := newUvHelper()
	versions, err := helper.ListPythonVersions()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(versions)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Installed Python Versions"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if len(versions) == 0 {
		fmt.Fprintln(os.Stdout, "No Python versions installed by UV")
		return nil
	}

	for _, v := range versions {
		fmt.Fprintf(os.Stdout, "  %s  %s\n", v.Version, v.Path)
	}
	return nil
}

func runUvPythonInstall(cmd *cobra.Command, args []string) error {
	helper := newUvHelper()
	if err := helper.InstallPython(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Python %s installed\n", output.Success("✓"), args[0])
	return nil
}

func runUvToolList(cmd *cobra.Command, args []string) error {
	helper := newUvHelper()
	tools, err := helper.ListTools()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(tools)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Installed Tools"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if len(tools) == 0 {
		fmt.Fprintln(os.Stdout, "No tools installed")
		return nil
	}

	for _, t := range tools {
		if t.Version != "" {
			fmt.Fprintf(os.Stdout, "  %s %s\n", t.Name, t.Version)
		} else {
			fmt.Fprintf(os.Stdout, "  %s\n", t.Name)
		}
	}
	return nil
}

func runUvToolInstall(cmd *cobra.Command, args []string) error {
	helper := newUvHelper()
	for _, tool := range args {
		if err := helper.InstallTool(tool); err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "%s Installed: %s\n", output.Success("✓"), tool)
	}
	return nil
}

func runUvToolRun(cmd *cobra.Command, args []string) error {
	helper := newUvHelper()
	return helper.RunTool(args...)
}

func runUvCacheInfo(cmd *cobra.Command, args []string) error {
	helper := newUvHelper()
	info, err := helper.GetCacheInfo()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(info)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("UV Cache"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintf(os.Stdout, "Directory: %s\n", info["directory"])
	fmt.Fprintf(os.Stdout, "Size: %s\n", info["size"])
	return nil
}

func runUvCacheClean(cmd *cobra.Command, args []string) error {
	helper := newUvHelper()
	if err := helper.CleanCache(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Cache cleaned\n", output.Success("✓"))
	return nil
}

func runUvCachePrune(cmd *cobra.Command, args []string) error {
	helper := newUvHelper()
	if err := helper.PruneCache(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Cache pruned\n", output.Success("✓"))
	return nil
}

func runUvInit(cmd *cobra.Command, args []string) error {
	helper := newUvHelper()
	name := ""
	if len(args) > 0 {
		name = args[0]
	}
	return helper.Init(name)
}

func runUvSync(cmd *cobra.Command, args []string) error {
	helper := newUvHelper()
	return helper.Sync(args...)
}

func runUvAdd(cmd *cobra.Command, args []string) error {
	helper := newUvHelper()
	return helper.Add(args...)
}

func runUvRemove(cmd *cobra.Command, args []string) error {
	helper := newUvHelper()
	return helper.Remove(args...)
}

func runUvRun(cmd *cobra.Command, args []string) error {
	helper := newUvHelper()
	return helper.Run(args...)
}

func runUvLock(cmd *cobra.Command, args []string) error {
	helper := newUvHelper()
	return helper.Lock()
}

func runUvTree(cmd *cobra.Command, args []string) error {
	helper := newUvHelper()
	return helper.Tree()
}

func runUvVenv(cmd *cobra.Command, args []string) error {
	helper := newUvHelper()
	name := ""
	if len(args) > 0 {
		name = args[0]
	}
	return helper.Venv(name)
}
