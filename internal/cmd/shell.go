package cmd

import (
	"github.com/mistergrinvalds/acorn/internal/components"
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/shell"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	shellDryRun  bool
	shellVerbose bool
)

// shellCmd represents the shell command group
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Shell integration management",
	Long: `Manage shell integration for acorn components.

Generates shell scripts (aliases, environment, functions) for each component
and injects them into your shell configuration.

Files are stored in $XDG_CONFIG_HOME/acorn/ (typically ~/.config/acorn/).

Examples:
  acorn shell status      # Show current integration status
  acorn shell generate    # Generate shell scripts
  acorn shell inject      # Add source line to shell rc
  acorn shell install     # Generate + inject (full setup)
  acorn shell eject       # Remove from shell rc`,
}

// shellStatusCmd shows shell integration status
var shellStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show shell integration status",
	Long: `Display the current state of shell integration.

Shows:
  - Detected shell and platform
  - Acorn config directory location
  - Whether shell rc is injected
  - Generated component scripts

Examples:
  acorn shell status
  acorn shell status -o json`,
	RunE: runShellStatus,
}

// shellGenerateCmd generates shell scripts
var shellGenerateCmd = &cobra.Command{
	Use:   "generate [component...]",
	Short: "Generate shell integration scripts",
	Long: `Generate shell scripts for specific or all components.

If no component is specified, generates all components plus the entrypoint.
If specific components are specified, only generates those.

Creates files in $XDG_CONFIG_HOME/acorn/:
  - shell.sh: Main entrypoint (sources all component scripts)
  - go.sh: Go development aliases and functions
  - vscode.sh: VS Code aliases and functions
  - tools.sh: Tool management functions

Examples:
  acorn shell generate              # Generate all
  acorn shell generate go           # Generate only go.sh
  acorn shell generate go vscode    # Generate go.sh and vscode.sh
  acorn shell generate -o json      # Output as JSON (includes file content)
  acorn shell generate --dry-run    # Show what would be done`,
	Aliases: []string{"gen"},
	RunE:    runShellGenerate,
}

// shellInjectCmd injects into shell rc
var shellInjectCmd = &cobra.Command{
	Use:   "inject",
	Short: "Inject acorn into shell configuration",
	Long: `Add the acorn shell integration source line to your shell rc file.

Modifies ~/.bashrc or ~/.zshrc (based on detected shell) to source
the acorn shell.sh entrypoint.

The injection is idempotent - running multiple times is safe.

Examples:
  acorn shell inject
  acorn shell inject --dry-run`,
	RunE: runShellInject,
}

// shellEjectCmd removes from shell rc
var shellEjectCmd = &cobra.Command{
	Use:   "eject",
	Short: "Remove acorn from shell configuration",
	Long: `Remove the acorn shell integration source line from your shell rc file.

This removes only the acorn-specific lines, leaving other configuration intact.

Examples:
  acorn shell eject
  acorn shell eject --dry-run`,
	RunE: runShellEject,
}

// shellInstallCmd does full installation
var shellInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Full shell integration setup",
	Long: `Perform full shell integration setup.

This command:
  1. Generates all shell scripts
  2. Injects the source line into your shell rc

Equivalent to running: acorn shell generate && acorn shell inject

Examples:
  acorn shell install
  acorn shell install --dry-run`,
	Aliases: []string{"setup"},
	RunE:    runShellInstall,
}

// shellUninstallCmd does full uninstallation
var shellUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove shell integration",
	Long: `Remove shell integration completely.

This command:
  1. Removes the source line from shell rc
  2. Optionally removes generated scripts

Examples:
  acorn shell uninstall
  acorn shell uninstall --dry-run`,
	Aliases: []string{"remove"},
	RunE:    runShellUninstall,
}

// shellListCmd lists available components
var shellListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available shell components",
	Long: `List all components that have shell integration.

Shows the components that will generate shell scripts when
'acorn shell generate' is run.

Examples:
  acorn shell list
  acorn shell list -o json`,
	Aliases: []string{"ls"},
	RunE:    runShellList,
}

func init() {

	// Add subcommands
	shellCmd.AddCommand(shellStatusCmd)
	shellCmd.AddCommand(shellGenerateCmd)
	shellCmd.AddCommand(shellInjectCmd)
	shellCmd.AddCommand(shellEjectCmd)
	shellCmd.AddCommand(shellInstallCmd)
	shellCmd.AddCommand(shellUninstallCmd)
	shellCmd.AddCommand(shellListCmd)

	// Persistent flags
	shellCmd.PersistentFlags().BoolVar(&shellDryRun, "dry-run", false,
		"Show what would be done without executing")
	shellCmd.PersistentFlags().BoolVarP(&shellVerbose, "verbose", "v", false,
		"Show verbose output")
}

func getShellManager() *shell.Manager {
	config := shell.NewConfig(shellVerbose, shellDryRun)
	manager := shell.NewManager(config)
	shell.RegisterAllComponents(manager)
	return manager
}

func runShellStatus(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	manager := getShellManager()
	status, err := manager.GetStatus()
	if err != nil {
		return err
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n\n", output.Info("Shell Integration Status"))
	fmt.Fprintf(os.Stdout, "Shell:        %s\n", status.Shell)
	fmt.Fprintf(os.Stdout, "Platform:     %s\n", status.Platform)
	fmt.Fprintf(os.Stdout, "Config Dir:   %s\n", status.AcornDir)

	if status.AcornDirExists {
		fmt.Fprintf(os.Stdout, "Dir Exists:   %s\n", output.Success("yes"))
	} else {
		fmt.Fprintf(os.Stdout, "Dir Exists:   %s\n", output.Warning("no"))
	}

	fmt.Fprintf(os.Stdout, "RC File:      %s\n", status.RCFile)

	if status.Injected {
		fmt.Fprintf(os.Stdout, "Injected:     %s\n", output.Success("yes"))
	} else {
		fmt.Fprintf(os.Stdout, "Injected:     %s\n", output.Warning("no"))
	}

	if len(status.GeneratedFiles) > 0 {
		fmt.Fprintf(os.Stdout, "\nGenerated Files:\n")
		for _, f := range status.GeneratedFiles {
			fmt.Fprintf(os.Stdout, "  %s\n", f)
		}
	}

	if len(status.Components) > 0 {
		fmt.Fprintf(os.Stdout, "\nRegistered Components:\n")
		for _, c := range status.Components {
			fmt.Fprintf(os.Stdout, "  %s\n", c)
		}
	}

	return nil
}

func runShellGenerate(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	manager := getShellManager()

	var result *shell.GenerateResult
	var err error

	if len(args) == 0 {
		// Generate all components + entrypoint
		result, err = manager.GenerateAll()
	} else {
		// Generate specific components only
		result, err = manager.GenerateComponents(args...)
	}

	if err != nil {
		return err
	}

	// JSON/YAML output - return structured result
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(result)
	}

	// Table format - human readable output
	if shellDryRun {
		fmt.Fprintln(os.Stdout, "[dry-run] Would generate shell scripts:")
	} else {
		fmt.Fprintln(os.Stdout, "Generated shell scripts:")
	}
	fmt.Fprintln(os.Stdout)

	for _, script := range result.Scripts {
		status := output.Success("✓")
		if shellDryRun {
			status = output.Warning("○")
		}
		fmt.Fprintf(os.Stdout, "  %s %s\n", status, script.GeneratedPath)
		fmt.Fprintf(os.Stdout, "    → symlink to: %s\n", script.SymlinkPath)
		if shellVerbose {
			fmt.Fprintf(os.Stdout, "    Description: %s\n", script.Description)
		}
	}

	if result.Entrypoint != nil {
		status := output.Success("✓")
		if shellDryRun {
			status = output.Warning("○")
		}
		fmt.Fprintf(os.Stdout, "  %s %s (entrypoint)\n", status, result.Entrypoint.GeneratedPath)
		fmt.Fprintf(os.Stdout, "    → symlink to: %s\n", result.Entrypoint.SymlinkPath)
	}

	// Show config files
	if len(result.ConfigFiles) > 0 {
		fmt.Fprintln(os.Stdout)
		if shellDryRun {
			fmt.Fprintln(os.Stdout, "[dry-run] Would generate config files:")
		} else {
			fmt.Fprintln(os.Stdout, "Generated config files:")
		}
		for _, cf := range result.ConfigFiles {
			status := output.Success("✓")
			if shellDryRun {
				status = output.Warning("○")
			}
			fmt.Fprintf(os.Stdout, "  %s %s (%s)\n", status, cf.GeneratedPath, cf.Format)
			fmt.Fprintf(os.Stdout, "    → symlink to: %s\n", cf.SymlinkTarget)
		}
	}

	fmt.Fprintln(os.Stdout)
	if shellDryRun {
		fmt.Fprintf(os.Stdout, "Use without --dry-run to write files.\n")
	} else {
		totalFiles := len(result.Scripts) + len(result.ConfigFiles)
		fmt.Fprintf(os.Stdout, "%s Generated %d file(s)\n", output.Success("✓"), totalFiles)
	}

	return nil
}

func runShellInject(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	manager := getShellManager()

	result, err := manager.Inject()
	if err != nil {
		return err
	}

	// JSON/YAML output
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(result)
	}

	// Table format
	switch result.Action {
	case "already_injected":
		fmt.Fprintf(os.Stdout, "%s Already injected in %s\n", output.Info("ℹ"), result.RCFile)
	case "would_inject":
		fmt.Fprintf(os.Stdout, "[dry-run] Would inject into: %s\n", result.RCFile)
		if shellVerbose {
			fmt.Fprintf(os.Stdout, "Injection block:\n%s\n", result.InjectionBlock)
		}
	case "injected":
		fmt.Fprintf(os.Stdout, "%s Injected into %s\n", output.Success("✓"), result.RCFile)
		fmt.Fprintf(os.Stdout, "Restart your shell or run: source %s\n", result.RCFile)
	}

	return nil
}

func runShellEject(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	manager := getShellManager()

	result, err := manager.Eject()
	if err != nil {
		return err
	}

	// JSON/YAML output
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(result)
	}

	// Table format
	switch result.Action {
	case "not_injected":
		fmt.Fprintf(os.Stdout, "%s No acorn injection found in %s\n", output.Info("ℹ"), result.RCFile)
	case "would_eject":
		fmt.Fprintf(os.Stdout, "[dry-run] Would remove acorn from: %s\n", result.RCFile)
	case "ejected":
		fmt.Fprintf(os.Stdout, "%s Removed acorn from %s\n", output.Success("✓"), result.RCFile)
	}

	return nil
}

// InstallResult combines generate and inject results.
type InstallResult struct {
	Generate *shell.GenerateResult `json:"generate" yaml:"generate"`
	Inject   *shell.InjectResult   `json:"inject" yaml:"inject"`
}

func runShellInstall(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	manager := getShellManager()

	// Generate all
	genResult, err := manager.GenerateAll()
	if err != nil {
		return err
	}

	// Inject
	injectResult, err := manager.Inject()
	if err != nil {
		return err
	}

	installResult := &InstallResult{
		Generate: genResult,
		Inject:   injectResult,
	}

	// JSON/YAML output
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(installResult)
	}

	// Table format
	fmt.Fprintln(os.Stdout, "Shell integration install:")
	fmt.Fprintln(os.Stdout)

	// Show generated scripts
	fmt.Fprintln(os.Stdout, "Generated scripts:")
	for _, script := range genResult.Scripts {
		status := output.Success("✓")
		if shellDryRun {
			status = output.Warning("○")
		}
		fmt.Fprintf(os.Stdout, "  %s %s\n", status, script.GeneratedPath)
		fmt.Fprintf(os.Stdout, "    → symlink to: %s\n", script.SymlinkPath)
	}
	if genResult.Entrypoint != nil {
		status := output.Success("✓")
		if shellDryRun {
			status = output.Warning("○")
		}
		fmt.Fprintf(os.Stdout, "  %s %s (entrypoint)\n", status, genResult.Entrypoint.GeneratedPath)
		fmt.Fprintf(os.Stdout, "    → symlink to: %s\n", genResult.Entrypoint.SymlinkPath)
	}

	// Show config files
	if len(genResult.ConfigFiles) > 0 {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Generated config files:")
		for _, cf := range genResult.ConfigFiles {
			status := output.Success("✓")
			if shellDryRun {
				status = output.Warning("○")
			}
			fmt.Fprintf(os.Stdout, "  %s %s (%s)\n", status, cf.GeneratedPath, cf.Format)
			fmt.Fprintf(os.Stdout, "    → symlink to: %s\n", cf.SymlinkTarget)
		}
	}

	fmt.Fprintln(os.Stdout)

	// Show injection status
	switch injectResult.Action {
	case "already_injected":
		fmt.Fprintf(os.Stdout, "%s Already injected in %s\n", output.Info("ℹ"), injectResult.RCFile)
	case "would_inject":
		fmt.Fprintf(os.Stdout, "[dry-run] Would inject into: %s\n", injectResult.RCFile)
	case "injected":
		fmt.Fprintf(os.Stdout, "%s Injected into %s\n", output.Success("✓"), injectResult.RCFile)
	}

	if !shellDryRun && injectResult.Action == "injected" {
		fmt.Fprintf(os.Stdout, "\nRestart your shell or run: source %s\n", injectResult.RCFile)
	}

	return nil
}

// UninstallResult contains the eject result.
type UninstallResult struct {
	Eject    *shell.InjectResult `json:"eject" yaml:"eject"`
	AcornDir string              `json:"acorn_dir" yaml:"acorn_dir"`
	Note     string              `json:"note" yaml:"note"`
}

func runShellUninstall(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	manager := getShellManager()

	result, err := manager.Eject()
	if err != nil {
		return err
	}

	config := shell.NewConfig(false, false)
	uninstallResult := &UninstallResult{
		Eject:    result,
		AcornDir: config.AcornDir,
		Note:     "Generated scripts were not removed. Delete " + config.AcornDir + " manually if desired.",
	}

	// JSON/YAML output
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(uninstallResult)
	}

	// Table format
	switch result.Action {
	case "not_injected":
		fmt.Fprintf(os.Stdout, "%s No acorn injection found in %s\n", output.Info("ℹ"), result.RCFile)
	case "would_eject":
		fmt.Fprintf(os.Stdout, "[dry-run] Would remove acorn from: %s\n", result.RCFile)
	case "ejected":
		fmt.Fprintf(os.Stdout, "%s Removed acorn from %s\n", output.Success("✓"), result.RCFile)
		fmt.Fprintf(os.Stdout, "\nNote: Generated scripts in %s were not removed.\n", config.AcornDir)
	}

	return nil
}

func runShellList(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	manager := getShellManager()
	status, err := manager.GetStatus()
	if err != nil {
		return err
	}

	type ComponentInfo struct {
		Name string `json:"name" yaml:"name"`
	}

	components := make([]ComponentInfo, 0, len(status.Components))
	for _, c := range status.Components {
		components = append(components, ComponentInfo{Name: c})
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(components)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n\n", output.Info("Available Shell Components"))
	for _, c := range status.Components {
		fmt.Fprintf(os.Stdout, "  %s\n", c)
	}
	fmt.Fprintf(os.Stdout, "\nTotal: %d components\n", len(status.Components))

	return nil
}

func init() {
	components.Register(&components.Registration{
		Name: "shell",
		RegisterCmd: func() *cobra.Command { return shellCmd },
	})
}
