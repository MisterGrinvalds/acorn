package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/ai/opencode"
	"github.com/mistergrinvalds/acorn/internal/utils/configcmd"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	opencodeDryRun  bool
	opencodeVerbose bool
)

// opencodeCmd represents the opencode command group
var opencodeCmd = &cobra.Command{
	Use:   "opencode",
	Short: "OpenCode AI coding assistant",
	Long: `OpenCode AI coding assistant for the terminal.

OpenCode is an open-source AI coding agent that helps you write,
debug, and understand code directly in your terminal.

Examples:
  acorn ai opencode status     # Show status
  acorn ai opencode launch     # Start OpenCode TUI
  acorn ai opencode providers  # List AI providers
  acorn ai opencode upgrade    # Upgrade to latest`,
	Aliases: []string{"oc"},
}

// opencodeStatusCmd shows status
var opencodeStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show OpenCode status",
	Long: `Display OpenCode installation and configuration status.

Examples:
  acorn ai opencode status
  acorn ai opencode status -o json`,
	RunE: runOpencodeStatus,
}

// opencodeLaunchCmd launches OpenCode
var opencodeLaunchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch OpenCode TUI",
	Long: `Launch the OpenCode terminal user interface.

This starts the interactive OpenCode session in your terminal.

Examples:
  acorn ai opencode launch`,
	Aliases: []string{"start", "run"},
	RunE:    runOpencodeLaunch,
}

// opencodeProvidersCmd lists providers
var opencodeProvidersCmd = &cobra.Command{
	Use:   "providers",
	Short: "List AI providers",
	Long: `List supported AI providers for OpenCode.

OpenCode supports multiple AI providers including OpenAI,
Anthropic, Google, Groq, AWS Bedrock, and more.

Examples:
  acorn ai opencode providers
  acorn ai opencode providers -o json`,
	RunE: runOpencodeProviders,
}

// opencodeInstallCmd installs OpenCode
var opencodeInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install OpenCode",
	Long: `Install OpenCode using Homebrew.

Examples:
  acorn ai opencode install`,
	RunE: runOpencodeInstall,
}

// opencodeUpgradeCmd upgrades OpenCode
var opencodeUpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade OpenCode",
	Long: `Upgrade OpenCode to the latest version.

Examples:
  acorn ai opencode upgrade`,
	RunE: runOpencodeUpgrade,
}

func init() {
	aiCmd.AddCommand(opencodeCmd)

	// Add subcommands
	opencodeCmd.AddCommand(opencodeStatusCmd)
	opencodeCmd.AddCommand(opencodeLaunchCmd)
	opencodeCmd.AddCommand(opencodeProvidersCmd)
	opencodeCmd.AddCommand(opencodeInstallCmd)
	opencodeCmd.AddCommand(opencodeUpgradeCmd)
	opencodeCmd.AddCommand(configcmd.NewConfigRouter("opencode"))

	// Persistent flags
	opencodeCmd.PersistentFlags().BoolVar(&opencodeDryRun, "dry-run", false,
		"Show what would be done without executing")
	opencodeCmd.PersistentFlags().BoolVarP(&opencodeVerbose, "verbose", "v", false,
		"Show verbose output")
}

func newOpencodeHelper() *opencode.Helper {
	return opencode.NewHelper(opencodeVerbose, opencodeDryRun)
}

func runOpencodeStatus(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := newOpencodeHelper()
	status := helper.GetStatus()

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("OpenCode Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if status.Installed {
		fmt.Fprintf(os.Stdout, "%s OpenCode installed: v%s\n", output.Success("✓"), status.Version)
	} else {
		fmt.Fprintf(os.Stdout, "%s OpenCode not installed\n", output.Error("✗"))
		fmt.Fprintln(os.Stdout, "  Install: acorn ai opencode install")
		return nil
	}

	fmt.Fprintln(os.Stdout)

	if status.HasConfig {
		fmt.Fprintf(os.Stdout, "%s Configuration found\n", output.Success("✓"))
	} else {
		fmt.Fprintf(os.Stdout, "%s No configuration\n", output.Info("ℹ"))
		fmt.Fprintln(os.Stdout, "  Run 'opencode' and use /init to initialize")
	}

	fmt.Fprintf(os.Stdout, "\nConfig dir: %s\n", status.ConfigDir)

	return nil
}

func runOpencodeLaunch(cmd *cobra.Command, args []string) error {
	helper := newOpencodeHelper()
	return helper.Launch(args...)
}

func runOpencodeProviders(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := newOpencodeHelper()
	providers := helper.ListProviders()

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(providers)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Supported Providers"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for _, p := range providers {
		fmt.Fprintf(os.Stdout, "  %-15s %s\n", p.Name, p.Model)
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "Set provider via OPENCODE_PROVIDER environment variable")
	fmt.Fprintln(os.Stdout, "Set API key via ANTHROPIC_API_KEY, OPENAI_API_KEY, etc.")

	return nil
}

func runOpencodeInstall(cmd *cobra.Command, args []string) error {
	helper := newOpencodeHelper()
	if err := helper.Install(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s OpenCode installed\n", output.Success("✓"))
	return nil
}

func runOpencodeUpgrade(cmd *cobra.Command, args []string) error {
	helper := newOpencodeHelper()
	if err := helper.Upgrade(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s OpenCode upgraded\n", output.Success("✓"))
	return nil
}
