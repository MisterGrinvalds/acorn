package cmd

import (
	"github.com/mistergrinvalds/acorn/internal/components"
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/neovim"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/mistergrinvalds/acorn/internal/utils/configcmd"
	"github.com/spf13/cobra"
)

var (
	nvimVerbose bool
	nvimForce   bool
)

// nvimCmd represents the neovim command group
var nvimCmd = &cobra.Command{
	Use:   "nvim",
	Short: "Neovim configuration management",
	Long: `Neovim configuration management and helper commands.

Provides health checks, config updates, and cache cleaning.

Examples:
  acorn nvim health     # Show Neovim health status
  acorn nvim update     # Update config repo
  acorn nvim clean      # Clean cache/data directories
  acorn nvim plugin     # Show dotfiles plugin info`,
	Aliases: []string{"neovim"},
}

// nvimHealthCmd shows health status
var nvimHealthCmd = &cobra.Command{
	Use:   "health",
	Short: "Show Neovim health status",
	Long: `Display Neovim installation and configuration status.

Shows version, config location, init file, and plugin manager.

Examples:
  acorn nvim health
  acorn nvim health -o json`,
	Aliases: []string{"status"},
	RunE:    runNvimHealth,
}

// nvimUpdateCmd updates config
var nvimUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Neovim config repo",
	Long: `Pull latest changes for the Neovim configuration repository.

Works with both symlinked and direct config directories.

Examples:
  acorn nvim update`,
	RunE: runNvimUpdate,
}

// nvimCleanCmd cleans cache
var nvimCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean Neovim cache and data",
	Long: `Remove Neovim data, cache, and state directories.

This will cause plugins to be reinstalled on next launch.
Requires --force flag to actually delete files.

Examples:
  acorn nvim clean --force`,
	RunE: runNvimClean,
}

// nvimPluginCmd shows plugin info
var nvimPluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Show dotfiles.nvim plugin info",
	Long: `Display setup instructions for the dotfiles.nvim plugin.

The plugin allows running the dotfiles installer from within Neovim.

Examples:
  acorn nvim plugin`,
	RunE: runNvimPlugin,
}

func init() {

	// Add subcommands
	nvimCmd.AddCommand(nvimHealthCmd)
	nvimCmd.AddCommand(nvimUpdateCmd)
	nvimCmd.AddCommand(nvimCleanCmd)
	nvimCmd.AddCommand(nvimPluginCmd)
	nvimCmd.AddCommand(configcmd.NewConfigRouter("neovim"))

	// Persistent flags
	nvimCmd.PersistentFlags().BoolVarP(&nvimVerbose, "verbose", "v", false,
		"Show verbose output")

	// Clean command flags
	nvimCleanCmd.Flags().BoolVar(&nvimForce, "force", false,
		"Actually clean the directories (required)")
}

func runNvimHealth(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := neovim.NewHelper(nvimVerbose)
	status := helper.GetHealth()

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Neovim Health Check"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if status.Installed {
		fmt.Fprintf(os.Stdout, "%s %s\n", output.Success("✓"), status.Version)
	} else {
		fmt.Fprintf(os.Stdout, "%s Neovim not installed\n", output.Error("✗"))
		return nil
	}

	fmt.Fprintln(os.Stdout)

	switch status.ConfigType {
	case "symlink":
		fmt.Fprintf(os.Stdout, "Config: %s -> %s\n", status.ConfigDir, status.ConfigTarget)
		if status.ConfigStatus == "ok" {
			fmt.Fprintf(os.Stdout, "Status: %s\n", output.Success("OK"))
		} else {
			fmt.Fprintf(os.Stdout, "Status: %s (target doesn't exist)\n", output.Error("BROKEN LINK"))
		}
	case "directory":
		fmt.Fprintf(os.Stdout, "Config: %s (direct directory)\n", status.ConfigDir)
		fmt.Fprintf(os.Stdout, "Status: %s\n", output.Success("OK"))
	default:
		fmt.Fprintln(os.Stdout, "Config: NOT FOUND")
		fmt.Fprintln(os.Stdout, "Run 'nvim_setup' to configure")
	}

	if status.InitFile != "" {
		fmt.Fprintf(os.Stdout, "\nInit file: %s\n", status.InitFile)
	}

	if status.PluginManager != "" {
		fmt.Fprintf(os.Stdout, "Plugin manager: %s\n", status.PluginManager)
	}

	return nil
}

func runNvimUpdate(cmd *cobra.Command, args []string) error {
	helper := neovim.NewHelper(nvimVerbose)

	if err := helper.Update(); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Config updated\n", output.Success("✓"))
	return nil
}

func runNvimClean(cmd *cobra.Command, args []string) error {
	helper := neovim.NewHelper(nvimVerbose)

	// Show what will be cleaned
	fmt.Fprintln(os.Stdout, "This will remove:")
	fmt.Fprintf(os.Stdout, "  - %s\n", helper.GetDataDir())
	fmt.Fprintf(os.Stdout, "  - %s\n", helper.GetCacheDir())
	fmt.Fprintf(os.Stdout, "  - %s\n", helper.GetStateDir())
	fmt.Fprintln(os.Stdout)

	if err := helper.Clean(nvimForce); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Cleaned. Plugins will be reinstalled on next nvim launch.\n", output.Success("✓"))
	return nil
}

func runNvimPlugin(cmd *cobra.Command, args []string) error {
	helper := neovim.NewHelper(nvimVerbose)
	fmt.Fprintln(os.Stdout, helper.GetPluginInfo())
	return nil
}

func init() {
	components.Register(&components.Registration{
		Name: "neovim",
		RegisterCmd: func() *cobra.Command { return nvimCmd },
	})
}
