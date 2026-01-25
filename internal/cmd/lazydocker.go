package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/devops/lazydocker"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/spf13/cobra"
)

var (
	lazydockerVerbose    bool
	lazydockerConfigFile string
)

// lazydockerCmd represents the lazydocker command group
var lazydockerCmd = &cobra.Command{
	Use:   "lazydocker",
	Short: "Lazydocker terminal UI commands",
	Long: `Lazydocker terminal UI for Docker container management.

Provides quick access to lazydocker and related utilities.

Examples:
  acorn lazydocker           # Launch lazydocker
  acorn lazydocker status    # Show installation status
  acorn lazydocker keys      # Show keybindings
  acorn lazydocker config    # Edit config file`,
	Aliases: []string{"lzd"},
	RunE:    runLazydockerLaunch,
}

// lazydockerStatusCmd shows status
var lazydockerStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show lazydocker installation status",
	Long: `Display lazydocker installation status and configuration.

Examples:
  acorn lazydocker status
  acorn lazydocker status -o json`,
	RunE: runLazydockerStatus,
}

// lazydockerLaunchCmd launches lazydocker
var lazydockerLaunchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch lazydocker",
	Long: `Start the lazydocker terminal UI.

Examples:
  acorn lazydocker launch
  acorn lazydocker launch --config /path/to/config.yml`,
	Aliases: []string{"start", "open"},
	RunE:    runLazydockerLaunch,
}

// lazydockerKeysCmd shows keybindings
var lazydockerKeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Show lazydocker keybindings",
	Long: `Display lazydocker keyboard shortcuts and commands.

Examples:
  acorn lazydocker keys
  acorn lazydocker keys -o json`,
	Aliases: []string{"keybindings", "shortcuts"},
	RunE:    runLazydockerKeys,
}

// lazydockerConfigCmd edits config
var lazydockerConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Edit lazydocker configuration",
	Long: `Open lazydocker configuration file in your editor.

Creates a default config file if one doesn't exist.

Examples:
  acorn lazydocker config`,
	Aliases: []string{"edit"},
	RunE:    runLazydockerConfig,
}

func init() {
	devopsCmd.AddCommand(lazydockerCmd)

	// Add subcommands
	lazydockerCmd.AddCommand(lazydockerStatusCmd)
	lazydockerCmd.AddCommand(lazydockerLaunchCmd)
	lazydockerCmd.AddCommand(lazydockerKeysCmd)
	lazydockerCmd.AddCommand(lazydockerConfigCmd)

	// Persistent flags
	lazydockerCmd.PersistentFlags().BoolVarP(&lazydockerVerbose, "verbose", "v", false,
		"Show verbose output")

	// Command-specific flags
	lazydockerLaunchCmd.Flags().StringVar(&lazydockerConfigFile, "config", "",
		"Path to config file")
	lazydockerCmd.Flags().StringVar(&lazydockerConfigFile, "config", "",
		"Path to config file")
}

func runLazydockerStatus(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := lazydocker.NewHelper(lazydockerVerbose)
	status := helper.GetStatus()

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	fmt.Fprintf(os.Stdout, "Installed:      %v\n", status.Installed)
	if status.Version != "" {
		fmt.Fprintf(os.Stdout, "Version:        %s\n", status.Version)
	}
	if status.Location != "" {
		fmt.Fprintf(os.Stdout, "Location:       %s\n", status.Location)
	}
	if status.ConfigDir != "" {
		fmt.Fprintf(os.Stdout, "Config Dir:     %s\n", status.ConfigDir)
	}
	fmt.Fprintf(os.Stdout, "Docker Running: %v\n", status.DockerRunning)

	return nil
}

func runLazydockerLaunch(cmd *cobra.Command, args []string) error {
	helper := lazydocker.NewHelper(lazydockerVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("lazydocker is not installed. Install with: brew install lazydocker")
	}

	if !helper.IsDockerRunning() {
		return fmt.Errorf("docker is not running. Please start Docker first")
	}

	return helper.Launch(lazydockerConfigFile)
}

func runLazydockerKeys(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := lazydocker.NewHelper(lazydockerVerbose)

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(map[string]any{"keybindings": helper.GetKeybindingsList()})
	}

	fmt.Print(helper.GetKeybindings())
	return nil
}

func runLazydockerConfig(cmd *cobra.Command, args []string) error {
	helper := lazydocker.NewHelper(lazydockerVerbose)
	return helper.OpenConfig()
}
