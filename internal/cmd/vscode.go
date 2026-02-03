package cmd

import (
	"github.com/mistergrinvalds/acorn/internal/components"
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/vscode"
	"github.com/mistergrinvalds/acorn/internal/utils/configcmd"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	vscodeDryRun  bool
	vscodeVerbose bool
)

// vscodeCmd represents the vscode command group
var vscodeCmd = &cobra.Command{
	Use:   "vscode",
	Short: "VS Code integration and helpers",
	Long: `Manage VS Code workspaces, projects, extensions, and configuration.

Examples:
  acorn vscode workspaces           # List workspaces
  acorn vscode workspace myproject  # Open workspace
  acorn vscode project new myapp go # Create Go project
  acorn vscode ext list             # List extensions
  acorn vscode config sync          # Sync config from dotfiles`,
	Aliases: []string{"code", "vs"},
}

// vscodeWorkspacesCmd lists workspaces
var vscodeWorkspacesCmd = &cobra.Command{
	Use:   "workspaces",
	Short: "List available workspaces",
	Long: `List all VS Code workspaces in ~/.vscode/workspaces.

Examples:
  acorn vscode workspaces
  acorn vscode workspaces -o json`,
	Aliases: []string{"ws"},
	RunE:    runVscodeWorkspaces,
}

// vscodeWorkspaceCmd opens a workspace
var vscodeWorkspaceCmd = &cobra.Command{
	Use:   "workspace <name>",
	Short: "Open a VS Code workspace",
	Long: `Open a VS Code workspace by name.

Workspaces are stored in ~/.vscode/workspaces/<name>.code-workspace

Examples:
  acorn vscode workspace myproject`,
	Args: cobra.ExactArgs(1),
	RunE: runVscodeWorkspace,
}

// vscodeProjectCmd is the parent for project commands
var vscodeProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Project management commands",
	Long:  `Commands for creating and managing VS Code projects.`,
}

// vscodeProjectNewCmd creates a new project
var vscodeProjectNewCmd = &cobra.Command{
	Use:   "new <name> [language]",
	Short: "Create a new VS Code project",
	Long: `Create a new project with VS Code settings.

Supported languages: python, go, typescript, general (default)
Aliases: py, golang, ts, node, js

Examples:
  acorn vscode project new myapp
  acorn vscode project new myapi go
  acorn vscode project new frontend typescript`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runVscodeProjectNew,
}

// vscodeExtCmd is the parent for extension commands
var vscodeExtCmd = &cobra.Command{
	Use:     "ext",
	Short:   "Extension management commands",
	Long:    `Commands for managing VS Code extensions.`,
	Aliases: []string{"extensions", "extension"},
}

// vscodeExtListCmd lists extensions
var vscodeExtListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed extensions",
	Long: `List all installed VS Code extensions.

Examples:
  acorn vscode ext list
  acorn vscode ext list -o json`,
	Aliases: []string{"ls"},
	RunE:    runVscodeExtList,
}

// vscodeExtInstallCmd installs extensions
var vscodeExtInstallCmd = &cobra.Command{
	Use:   "install [file]",
	Short: "Install extensions from file",
	Long: `Install VS Code extensions from a file.

If no file is specified, uses the extensions file from dotfiles.

Examples:
  acorn vscode ext install
  acorn vscode ext install extensions.txt`,
	Args: cobra.MaximumNArgs(1),
	RunE: runVscodeExtInstall,
}

// vscodeExtExportCmd exports extensions
var vscodeExtExportCmd = &cobra.Command{
	Use:   "export [file]",
	Short: "Export installed extensions to file",
	Long: `Export all installed extensions to a file.

Default filename is vscode-extensions.txt

Examples:
  acorn vscode ext export
  acorn vscode ext export my-extensions.txt`,
	Args: cobra.MaximumNArgs(1),
	RunE: runVscodeExtExport,
}

// vscodeExtEssentialsCmd installs essential extensions
var vscodeExtEssentialsCmd = &cobra.Command{
	Use:   "essentials",
	Short: "Install essential extensions",
	Long: `Install a curated set of essential VS Code extensions.

Includes:
  - Python, Go language support
  - GitHub Pull Requests
  - GitLens
  - Kubernetes, Docker tools
  - Catppuccin theme

Examples:
  acorn vscode ext essentials`,
	RunE: runVscodeExtEssentials,
}

// vscodeConfigSyncCmd syncs config
var vscodeConfigSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync configuration from dotfiles",
	Long: `Sync VS Code settings and keybindings from dotfiles.

Backs up existing files before syncing.

Examples:
  acorn vscode config sync
  acorn vscode config sync --dry-run`,
	RunE: runVscodeConfigSync,
}

// vscodeConfigPathCmd is now provided by the universal config router

func init() {

	// Add subcommands
	vscodeCmd.AddCommand(vscodeWorkspacesCmd)
	vscodeCmd.AddCommand(vscodeWorkspaceCmd)
	vscodeCmd.AddCommand(vscodeProjectCmd)
	vscodeCmd.AddCommand(vscodeExtCmd)
	vscodeConfigRouter := configcmd.NewConfigRouter("vscode")
	vscodeConfigRouter.AddCommand(vscodeConfigSyncCmd)
	vscodeCmd.AddCommand(vscodeConfigRouter)

	// Project subcommands
	vscodeProjectCmd.AddCommand(vscodeProjectNewCmd)

	// Extension subcommands
	vscodeExtCmd.AddCommand(vscodeExtListCmd)
	vscodeExtCmd.AddCommand(vscodeExtInstallCmd)
	vscodeExtCmd.AddCommand(vscodeExtExportCmd)
	vscodeExtCmd.AddCommand(vscodeExtEssentialsCmd)

	// Persistent flags
	vscodeCmd.PersistentFlags().BoolVar(&vscodeDryRun, "dry-run", false,
		"Show what would be done without executing")
	vscodeCmd.PersistentFlags().BoolVarP(&vscodeVerbose, "verbose", "v", false,
		"Show verbose output")
}

func runVscodeWorkspaces(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := vscode.NewHelper(vscodeVerbose, vscodeDryRun)
	workspaces, err := helper.ListWorkspaces()
	if err != nil {
		return err
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(workspaces)
	}

	if len(workspaces) == 0 {
		fmt.Fprintln(os.Stdout, "No workspaces found")
		paths := helper.GetConfigPaths()
		fmt.Fprintf(os.Stdout, "Workspace directory: %s\n", paths.WorkspaceDir)
		return nil
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n\n", output.Info("Workspaces"))
	for _, ws := range workspaces {
		fmt.Fprintf(os.Stdout, "  %s\n", ws.Name)
	}
	fmt.Fprintf(os.Stdout, "\nTotal: %d workspaces\n", len(workspaces))

	return nil
}

func runVscodeWorkspace(cmd *cobra.Command, args []string) error {
	helper := vscode.NewHelper(vscodeVerbose, vscodeDryRun)

	if err := helper.OpenWorkspace(args[0]); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Opening workspace: %s\n", output.Success("✓"), args[0])
	return nil
}

func runVscodeProjectNew(cmd *cobra.Command, args []string) error {
	helper := vscode.NewHelper(vscodeVerbose, vscodeDryRun)

	name := args[0]
	language := ""
	if len(args) > 1 {
		language = args[1]
	}

	if err := helper.CreateProject(name, language); err != nil {
		return err
	}

	lang := language
	if lang == "" {
		lang = "general"
	}

	fmt.Fprintf(os.Stdout, "%s VS Code project created!\n", output.Success("✓"))
	fmt.Fprintf(os.Stdout, "  Name: %s\n", name)
	fmt.Fprintf(os.Stdout, "  Settings: %s\n", lang)
	fmt.Fprintf(os.Stdout, "\nNext steps:\n")
	fmt.Fprintf(os.Stdout, "  cd %s && code .\n", name)

	return nil
}

func runVscodeExtList(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	if !vscode.IsInstalled() {
		return fmt.Errorf("VS Code is not installed")
	}

	helper := vscode.NewHelper(vscodeVerbose, vscodeDryRun)
	extensions, err := helper.ListExtensions()
	if err != nil {
		return err
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(extensions)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n\n", output.Info("Installed Extensions"))
	for _, ext := range extensions {
		fmt.Fprintf(os.Stdout, "  %s\n", ext.ID)
	}
	fmt.Fprintf(os.Stdout, "\nTotal: %d extensions\n", len(extensions))

	return nil
}

func runVscodeExtInstall(cmd *cobra.Command, args []string) error {
	if !vscode.IsInstalled() {
		return fmt.Errorf("VS Code is not installed")
	}

	helper := vscode.NewHelper(vscodeVerbose, vscodeDryRun)

	var filePath string
	if len(args) > 0 {
		filePath = args[0]
	} else {
		// Default to dotfiles extensions file
		dotfilesRoot := os.Getenv("DOTFILES_ROOT")
		if dotfilesRoot == "" {
			home, _ := os.UserHomeDir()
			dotfilesRoot = home + "/.config/dotfiles"
		}
		filePath = dotfilesRoot + "/.sapling/generated/vscode/extensions.txt"
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("extensions file not found: %s", filePath)
	}

	fmt.Fprintf(os.Stdout, "Installing extensions from: %s\n\n", filePath)
	if err := helper.InstallExtensionsFromFile(filePath); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "\n%s Extensions installed!\n", output.Success("✓"))
	return nil
}

func runVscodeExtExport(cmd *cobra.Command, args []string) error {
	if !vscode.IsInstalled() {
		return fmt.Errorf("VS Code is not installed")
	}

	helper := vscode.NewHelper(vscodeVerbose, vscodeDryRun)

	filePath := "vscode-extensions.txt"
	if len(args) > 0 {
		filePath = args[0]
	}

	if err := helper.ExportExtensions(filePath); err != nil {
		return err
	}

	extensions, _ := helper.ListExtensions()
	fmt.Fprintf(os.Stdout, "%s Exported %d extensions to: %s\n",
		output.Success("✓"), len(extensions), filePath)

	return nil
}

func runVscodeExtEssentials(cmd *cobra.Command, args []string) error {
	if !vscode.IsInstalled() {
		return fmt.Errorf("VS Code is not installed")
	}

	helper := vscode.NewHelper(vscodeVerbose, vscodeDryRun)

	fmt.Fprintln(os.Stdout, "Installing essential extensions...")
	fmt.Fprintln(os.Stdout)

	if err := helper.InstallEssentialExtensions(); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "\n%s Essential extensions installed!\n", output.Success("✓"))
	return nil
}

func runVscodeConfigSync(cmd *cobra.Command, args []string) error {
	helper := vscode.NewHelper(vscodeVerbose, vscodeDryRun)

	fmt.Fprintln(os.Stdout, "Syncing VS Code configuration...")

	if err := helper.SyncConfig(); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Configuration synced!\n", output.Success("✓"))
	return nil
}

// runVscodeConfigPath has been replaced by the universal config router: acorn vscode config path

func init() {
	components.Register(&components.Registration{
		Name: "vscode",
		RegisterCmd: func() *cobra.Command { return vscodeCmd },
	})
}
