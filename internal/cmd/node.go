package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/node"
	"github.com/mistergrinvalds/acorn/internal/output"
	"github.com/spf13/cobra"
)

var (
	nodeOutputFormat string
	nodeVerbose      bool
	nodeDryRun       bool
	nodeForce        bool
)

// nodeCmd represents the node command group
var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Node.js ecosystem management",
	Long: `Node.js, NVM, and pnpm management commands.

Provides status checking, package manager detection, and cleanup utilities.

Examples:
  acorn node status       # Show Node.js ecosystem status
  acorn node detect       # Detect package manager
  acorn node clean        # Clean and reinstall node_modules
  acorn node find         # Find all node_modules
  acorn node cleanall     # Remove all node_modules`,
	Aliases: []string{"nodejs", "npm"},
}

// nodeStatusCmd shows status
var nodeStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Node.js ecosystem status",
	Long: `Display Node.js, npm, NVM, and pnpm installation status.

Shows versions and directory locations.

Examples:
  acorn node status
  acorn node status -o json`,
	RunE: runNodeStatus,
}

// nodeDetectCmd detects package manager
var nodeDetectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect package manager",
	Long: `Detect which package manager is used based on lock files.

Checks for pnpm-lock.yaml, yarn.lock, or package-lock.json.

Examples:
  acorn node detect
  acorn node detect -o json`,
	RunE: runNodeDetect,
}

// nodeCleanCmd cleans node_modules
var nodeCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean and reinstall node_modules",
	Long: `Remove node_modules and reinstall dependencies.

Automatically detects the package manager to use.

Examples:
  acorn node clean`,
	RunE: runNodeClean,
}

// nodeFindCmd finds node_modules
var nodeFindCmd = &cobra.Command{
	Use:   "find [path]",
	Short: "Find all node_modules directories",
	Long: `Find and list all node_modules directories with sizes.

Searches from current directory or specified path.

Examples:
  acorn node find
  acorn node find ~/projects`,
	RunE: runNodeFind,
}

// nodeCleanAllCmd cleans all node_modules
var nodeCleanAllCmd = &cobra.Command{
	Use:   "cleanall [path]",
	Short: "Remove all node_modules in directory tree",
	Long: `Find and remove all node_modules directories.

Requires --force flag to actually delete.

Examples:
  acorn node cleanall --force
  acorn node cleanall ~/projects --force`,
	RunE: runNodeCleanAll,
}

// nodeCacheCmd manages npm cache
var nodeCacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage npm cache",
	Long: `Show or clean npm cache.

Examples:
  acorn node cache           # Show cache info
  acorn node cache --clean   # Clean cache`,
	RunE: runNodeCache,
}

// nvmCmd represents the nvm command group
var nvmCmd = &cobra.Command{
	Use:   "nvm",
	Short: "NVM (Node Version Manager) commands",
	Long: `NVM management commands.

Provides installation, version listing, and setup utilities.

Examples:
  acorn nvm status     # Show NVM status
  acorn nvm list       # List installed versions
  acorn nvm install    # Install NVM`,
}

// nvmStatusCmd shows NVM status
var nvmStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show NVM status",
	Long: `Display NVM installation status and current Node version.

Examples:
  acorn nvm status
  acorn nvm status -o json`,
	RunE: runNvmStatus,
}

// nvmListCmd lists versions
var nvmListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed Node versions",
	Long: `List all Node versions installed via NVM.

Examples:
  acorn nvm list
  acorn nvm list -o json`,
	Aliases: []string{"ls"},
	RunE:    runNvmList,
}

// nvmInstallCmd installs NVM
var nvmInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install NVM",
	Long: `Download and install NVM (Node Version Manager).

Examples:
  acorn nvm install`,
	RunE: runNvmInstall,
}

// pnpmCmd represents the pnpm command group
var pnpmCmd = &cobra.Command{
	Use:   "pnpm",
	Short: "pnpm package manager commands",
	Long: `pnpm management commands.

Provides installation and status utilities.

Examples:
  acorn pnpm status    # Show pnpm status
  acorn pnpm install   # Install pnpm globally`,
}

// pnpmStatusCmd shows pnpm status
var pnpmStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show pnpm status",
	Long: `Display pnpm installation status.

Examples:
  acorn pnpm status
  acorn pnpm status -o json`,
	RunE: runPnpmStatus,
}

// pnpmInstallCmd installs pnpm
var pnpmInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install pnpm globally",
	Long: `Install pnpm as a global npm package.

Examples:
  acorn pnpm install`,
	RunE: runPnpmInstall,
}

var nodeCacheClean bool

func init() {
	rootCmd.AddCommand(nodeCmd)
	rootCmd.AddCommand(nvmCmd)
	rootCmd.AddCommand(pnpmCmd)

	// Node subcommands
	nodeCmd.AddCommand(nodeStatusCmd)
	nodeCmd.AddCommand(nodeDetectCmd)
	nodeCmd.AddCommand(nodeCleanCmd)
	nodeCmd.AddCommand(nodeFindCmd)
	nodeCmd.AddCommand(nodeCleanAllCmd)
	nodeCmd.AddCommand(nodeCacheCmd)

	// NVM subcommands
	nvmCmd.AddCommand(nvmStatusCmd)
	nvmCmd.AddCommand(nvmListCmd)
	nvmCmd.AddCommand(nvmInstallCmd)

	// pnpm subcommands
	pnpmCmd.AddCommand(pnpmStatusCmd)
	pnpmCmd.AddCommand(pnpmInstallCmd)

	// Node persistent flags
	nodeCmd.PersistentFlags().StringVarP(&nodeOutputFormat, "output", "o", "table",
		"Output format (table|json|yaml)")
	nodeCmd.PersistentFlags().BoolVarP(&nodeVerbose, "verbose", "v", false,
		"Show verbose output")
	nodeCmd.PersistentFlags().BoolVar(&nodeDryRun, "dry-run", false,
		"Show what would be done without executing")

	// NVM persistent flags
	nvmCmd.PersistentFlags().StringVarP(&nodeOutputFormat, "output", "o", "table",
		"Output format (table|json|yaml)")
	nvmCmd.PersistentFlags().BoolVarP(&nodeVerbose, "verbose", "v", false,
		"Show verbose output")
	nvmCmd.PersistentFlags().BoolVar(&nodeDryRun, "dry-run", false,
		"Show what would be done without executing")

	// pnpm persistent flags
	pnpmCmd.PersistentFlags().StringVarP(&nodeOutputFormat, "output", "o", "table",
		"Output format (table|json|yaml)")
	pnpmCmd.PersistentFlags().BoolVarP(&nodeVerbose, "verbose", "v", false,
		"Show verbose output")
	pnpmCmd.PersistentFlags().BoolVar(&nodeDryRun, "dry-run", false,
		"Show what would be done without executing")

	// Clean all flags
	nodeCleanAllCmd.Flags().BoolVar(&nodeForce, "force", false,
		"Actually remove the directories (required)")

	// Cache flags
	nodeCacheCmd.Flags().BoolVar(&nodeCacheClean, "clean", false,
		"Clean the npm cache")
}

func runNodeStatus(cmd *cobra.Command, args []string) error {
	helper := node.NewHelper(nodeVerbose, nodeDryRun)
	status := helper.GetStatus()

	format, err := output.ParseFormat(nodeOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Node.js Ecosystem Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if status.NodeInstalled {
		fmt.Fprintf(os.Stdout, "%s Node.js: %s\n", output.Success("✓"), status.NodeVersion)
	} else {
		fmt.Fprintf(os.Stdout, "%s Node.js: not installed\n", output.Error("✗"))
	}

	if status.NpmInstalled {
		fmt.Fprintf(os.Stdout, "%s npm: %s\n", output.Success("✓"), status.NpmVersion)
	} else {
		fmt.Fprintf(os.Stdout, "%s npm: not installed\n", output.Error("✗"))
	}

	fmt.Fprintln(os.Stdout)

	if status.NvmInstalled {
		fmt.Fprintf(os.Stdout, "%s NVM: installed\n", output.Success("✓"))
		fmt.Fprintf(os.Stdout, "  Dir: %s\n", status.NvmDir)
	} else {
		fmt.Fprintf(os.Stdout, "%s NVM: not installed\n", output.Warning("⚠"))
		fmt.Fprintln(os.Stdout, "  Run: acorn nvm install")
	}

	if status.PnpmInstalled {
		fmt.Fprintf(os.Stdout, "%s pnpm: %s\n", output.Success("✓"), status.PnpmVersion)
		fmt.Fprintf(os.Stdout, "  Home: %s\n", status.PnpmHome)
	} else {
		fmt.Fprintf(os.Stdout, "%s pnpm: not installed\n", output.Warning("⚠"))
		fmt.Fprintln(os.Stdout, "  Run: acorn pnpm install")
	}

	fmt.Fprintln(os.Stdout)
	cacheSize := helper.GetNpmCacheSize()
	fmt.Fprintf(os.Stdout, "npm cache: %s (%s)\n", status.NpmCacheDir, cacheSize)

	return nil
}

func runNodeDetect(cmd *cobra.Command, args []string) error {
	helper := node.NewHelper(nodeVerbose, nodeDryRun)
	pm := helper.DetectPackageManager()

	format, err := output.ParseFormat(nodeOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(pm)
	}

	fmt.Fprintf(os.Stdout, "Package manager: %s\n", pm.Name)
	if pm.LockFile != "" {
		fmt.Fprintf(os.Stdout, "Lock file: %s\n", pm.LockFile)
	}

	return nil
}

func runNodeClean(cmd *cobra.Command, args []string) error {
	helper := node.NewHelper(nodeVerbose, nodeDryRun)

	if err := helper.CleanNodeModules(); err != nil {
		return err
	}

	if !nodeDryRun {
		fmt.Fprintf(os.Stdout, "%s Dependencies reinstalled\n", output.Success("✓"))
	}

	return nil
}

func runNodeFind(cmd *cobra.Command, args []string) error {
	helper := node.NewHelper(nodeVerbose, nodeDryRun)

	root := "."
	if len(args) > 0 {
		root = args[0]
	}

	modules, err := helper.FindNodeModules(root)
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(nodeOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(map[string]interface{}{"node_modules": modules})
	}

	if len(modules) == 0 {
		fmt.Fprintln(os.Stdout, "No node_modules directories found")
		return nil
	}

	var totalSize int64
	for _, m := range modules {
		fmt.Fprintf(os.Stdout, "%s\t%s\n", m.Size, m.Path)
	}

	fmt.Fprintf(os.Stdout, "\nFound: %d node_modules directories\n", len(modules))
	if totalSize > 0 {
		fmt.Fprintf(os.Stdout, "Total: %d bytes\n", totalSize)
	}

	return nil
}

func runNodeCleanAll(cmd *cobra.Command, args []string) error {
	helper := node.NewHelper(nodeVerbose, nodeDryRun)

	root := "."
	if len(args) > 0 {
		root = args[0]
	}

	// First show what we found
	modules, err := helper.FindNodeModules(root)
	if err != nil {
		return err
	}

	if len(modules) == 0 {
		fmt.Fprintln(os.Stdout, "No node_modules directories found")
		return nil
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Found node_modules:"))
	for _, m := range modules {
		fmt.Fprintf(os.Stdout, "  %s\t%s\n", m.Size, m.Path)
	}
	fmt.Fprintln(os.Stdout)

	count, err := helper.CleanAllNodeModules(root, nodeForce)
	if err != nil {
		return err
	}

	if nodeDryRun {
		fmt.Fprintf(os.Stdout, "Would remove %d directories\n", count)
	} else {
		fmt.Fprintf(os.Stdout, "%s Removed %d node_modules directories\n", output.Success("✓"), count)
	}

	return nil
}

func runNodeCache(cmd *cobra.Command, args []string) error {
	helper := node.NewHelper(nodeVerbose, nodeDryRun)

	if nodeCacheClean {
		if err := helper.CleanNpmCache(); err != nil {
			return err
		}
		if !nodeDryRun {
			fmt.Fprintf(os.Stdout, "%s npm cache cleaned\n", output.Success("✓"))
		}
		return nil
	}

	// Show cache info
	cacheDir := helper.GetNpmCacheDir()
	cacheSize := helper.GetNpmCacheSize()

	format, err := output.ParseFormat(nodeOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(map[string]string{
			"cache_dir":  cacheDir,
			"cache_size": cacheSize,
		})
	}

	fmt.Fprintf(os.Stdout, "Cache dir:  %s\n", cacheDir)
	fmt.Fprintf(os.Stdout, "Cache size: %s\n", cacheSize)

	return nil
}

func runNvmStatus(cmd *cobra.Command, args []string) error {
	helper := node.NewHelper(nodeVerbose, nodeDryRun)
	status := helper.GetStatus()

	format, err := output.ParseFormat(nodeOutputFormat)
	if err != nil {
		return err
	}

	nvmStatus := map[string]interface{}{
		"installed":      status.NvmInstalled,
		"nvm_dir":        status.NvmDir,
		"current_node":   status.NodeVersion,
		"node_installed": status.NodeInstalled,
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(nvmStatus)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("NVM Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if status.NvmInstalled {
		fmt.Fprintf(os.Stdout, "%s NVM installed\n", output.Success("✓"))
		fmt.Fprintf(os.Stdout, "  Dir: %s\n", status.NvmDir)
	} else {
		fmt.Fprintf(os.Stdout, "%s NVM not installed\n", output.Error("✗"))
		fmt.Fprintln(os.Stdout, "  Run: acorn nvm install")
		return nil
	}

	fmt.Fprintln(os.Stdout)
	if status.NodeInstalled {
		fmt.Fprintf(os.Stdout, "Current Node: %s\n", status.NodeVersion)
	} else {
		fmt.Fprintln(os.Stdout, "Current Node: (none)")
	}

	return nil
}

func runNvmList(cmd *cobra.Command, args []string) error {
	helper := node.NewHelper(nodeVerbose, nodeDryRun)

	versions, err := helper.GetNvmVersions()
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(nodeOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(map[string]interface{}{
			"versions":        versions,
			"current_version": helper.GetCurrentNodeVersion(),
		})
	}

	if len(versions) == 0 {
		fmt.Fprintln(os.Stdout, "No Node versions installed via NVM")
		return nil
	}

	currentVersion := helper.GetCurrentNodeVersion()

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Installed Node Versions"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for _, v := range versions {
		marker := " "
		if v == currentVersion {
			marker = "*"
		}
		fmt.Fprintf(os.Stdout, "%s %s\n", marker, v)
	}

	return nil
}

func runNvmInstall(cmd *cobra.Command, args []string) error {
	helper := node.NewHelper(nodeVerbose, nodeDryRun)

	if err := helper.InstallNvm(); err != nil {
		return err
	}

	if !nodeDryRun {
		fmt.Fprintf(os.Stdout, "%s NVM installed\n", output.Success("✓"))
		fmt.Fprintln(os.Stdout, "Restart your shell or run: source ~/.bashrc")
	}

	return nil
}

func runPnpmStatus(cmd *cobra.Command, args []string) error {
	helper := node.NewHelper(nodeVerbose, nodeDryRun)
	status := helper.GetStatus()

	format, err := output.ParseFormat(nodeOutputFormat)
	if err != nil {
		return err
	}

	pnpmStatus := map[string]interface{}{
		"installed":  status.PnpmInstalled,
		"version":    status.PnpmVersion,
		"pnpm_home":  status.PnpmHome,
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(pnpmStatus)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("pnpm Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if status.PnpmInstalled {
		fmt.Fprintf(os.Stdout, "%s pnpm %s\n", output.Success("✓"), status.PnpmVersion)
		fmt.Fprintf(os.Stdout, "  Home: %s\n", status.PnpmHome)
	} else {
		fmt.Fprintf(os.Stdout, "%s pnpm not installed\n", output.Error("✗"))
		fmt.Fprintln(os.Stdout, "  Run: acorn pnpm install")
	}

	return nil
}

func runPnpmInstall(cmd *cobra.Command, args []string) error {
	helper := node.NewHelper(nodeVerbose, nodeDryRun)

	if err := helper.InstallPnpm(); err != nil {
		return err
	}

	if !nodeDryRun {
		fmt.Fprintf(os.Stdout, "%s pnpm installed globally\n", output.Success("✓"))
	}

	return nil
}
