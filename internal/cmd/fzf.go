package cmd

import (
	"github.com/mistergrinvalds/acorn/internal/components"
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/fzf"
	"github.com/mistergrinvalds/acorn/internal/utils/configcmd"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	fzfVerbose bool
)

// fzfCmd represents the fzf command group
var fzfCmd = &cobra.Command{
	Use:   "fzf",
	Short: "FZF fuzzy finder helpers",
	Long: `FZF fuzzy finder status and configuration.

Note: Most fzf functions are interactive and remain as shell functions.
Use this command to check status and view available functions.

Examples:
  acorn fzf status               # Check fzf installation
  acorn fzf config               # Show current configuration
  acorn fzf functions            # List available functions`,
}

// fzfStatusCmd shows fzf status
var fzfStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check FZF installation status",
	Long: `Check if FZF and related tools are installed.

Shows FZF version, location, and whether fd is available.

Examples:
  acorn fzf status
  acorn fzf status -o json`,
	RunE: runFzfStatus,
}

// fzfEnvCmd shows fzf runtime environment configuration
var fzfEnvCmd = &cobra.Command{
	Use:   "env",
	Short: "Show FZF runtime environment",
	Long: `Show current FZF configuration from environment variables.

Displays FZF_DEFAULT_COMMAND, FZF_DEFAULT_OPTS, and other
runtime environment variables that control FZF behavior.

Examples:
  acorn fzf env
  acorn fzf env -o json`,
	RunE: runFzfConfig,
}

// fzfFunctionsCmd lists available functions
var fzfFunctionsCmd = &cobra.Command{
	Use:   "functions",
	Short: "List available FZF shell functions",
	Long: `List all available FZF-powered shell functions.

These functions are interactive and run in your shell.

Examples:
  acorn fzf functions`,
	Aliases: []string{"funcs", "list"},
	RunE:    runFzfFunctions,
}

// fzfThemeCmd shows theme colors
var fzfThemeCmd = &cobra.Command{
	Use:   "theme",
	Short: "Show Catppuccin Mocha theme colors",
	Long: `Show the Catppuccin Mocha color scheme for FZF.

These colors are automatically applied via FZF_DEFAULT_OPTS.

Examples:
  acorn fzf theme`,
	RunE: runFzfTheme,
}

func init() {

	// Add subcommands
	fzfCmd.AddCommand(fzfStatusCmd)
	fzfCmd.AddCommand(fzfEnvCmd)
	fzfCmd.AddCommand(fzfFunctionsCmd)
	fzfCmd.AddCommand(configcmd.NewConfigRouter("fzf"))
	fzfCmd.AddCommand(fzfThemeCmd)

	// Persistent flags
	fzfCmd.PersistentFlags().BoolVarP(&fzfVerbose, "verbose", "v", false,
		"Show verbose output")
}

func runFzfStatus(cmd *cobra.Command, args []string) error {
	helper := fzf.NewHelper(fzfVerbose)
	status := helper.GetStatus()

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("FZF Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if !status.Installed {
		fmt.Fprintf(os.Stdout, "%s FZF not installed\n", output.Error("✗"))
		fmt.Fprintln(os.Stdout, "  Install: brew install fzf")
		return nil
	}

	fmt.Fprintf(os.Stdout, "%s FZF installed\n", output.Success("✓"))
	fmt.Fprintf(os.Stdout, "  Version:  %s\n", status.Version)
	fmt.Fprintf(os.Stdout, "  Location: %s\n", status.Location)

	fmt.Fprintln(os.Stdout)
	if status.FdInstalled {
		fmt.Fprintf(os.Stdout, "%s fd installed (%s)\n", output.Success("✓"), status.FdCommand)
	} else {
		fmt.Fprintf(os.Stdout, "%s fd not installed (using find fallback)\n", output.Warning("⚠"))
		fmt.Fprintln(os.Stdout, "  Install: brew install fd")
	}

	return nil
}

func runFzfConfig(cmd *cobra.Command, args []string) error {
	helper := fzf.NewHelper(fzfVerbose)
	config := helper.GetConfig()

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(config)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("FZF Configuration"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if config.DefaultCommand != "" {
		fmt.Fprintf(os.Stdout, "FZF_DEFAULT_COMMAND:\n  %s\n\n", config.DefaultCommand)
	} else {
		fmt.Fprintln(os.Stdout, "FZF_DEFAULT_COMMAND: (not set)")
	}

	if config.AltCCommand != "" {
		fmt.Fprintf(os.Stdout, "FZF_ALT_C_COMMAND:\n  %s\n\n", config.AltCCommand)
	}

	if config.CtrlTCommand != "" {
		fmt.Fprintf(os.Stdout, "FZF_CTRL_T_COMMAND:\n  %s\n\n", config.CtrlTCommand)
	}

	if config.DefaultOpts != "" {
		fmt.Fprintf(os.Stdout, "FZF_DEFAULT_OPTS:\n%s\n", config.DefaultOpts)
	}

	return nil
}

func runFzfFunctions(cmd *cobra.Command, args []string) error {
	helper := fzf.NewHelper(fzfVerbose)
	functions := helper.GetAvailableFunctions()

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(map[string][]string{"functions": functions})
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Available FZF Functions"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintln(os.Stdout)

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("File & Directory:"))
	fmt.Fprintln(os.Stdout, "  fzf_files (ff)     Interactive file finder with preview")
	fmt.Fprintln(os.Stdout, "  fe <query>         Find and edit file")
	fmt.Fprintln(os.Stdout, "  fzf_cd (fcd)       Interactive cd with preview")
	fmt.Fprintln(os.Stdout)

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Git:"))
	fmt.Fprintln(os.Stdout, "  fzf_git_branch     Interactive branch checkout")
	fmt.Fprintln(os.Stdout, "  fzf_git_log        Interactive log browser")
	fmt.Fprintln(os.Stdout, "  fzf_git_stash      Interactive stash browser")
	fmt.Fprintln(os.Stdout, "  fzf_git_add (fga)  Interactive git add")
	fmt.Fprintln(os.Stdout)

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("System:"))
	fmt.Fprintln(os.Stdout, "  fzf_kill (fkill)   Interactive process killer")
	fmt.Fprintln(os.Stdout, "  fzf_history (fh)   Interactive history search")
	fmt.Fprintln(os.Stdout, "  fzf_env (fenv)     Interactive env variable browser")
	fmt.Fprintln(os.Stdout)

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Kubernetes:"))
	fmt.Fprintln(os.Stdout, "  fzf_k8s_pod        Interactive pod selector")
	fmt.Fprintln(os.Stdout, "  fzf_k8s_logs       Interactive pod logs")
	fmt.Fprintln(os.Stdout, "  fzf_k8s_ns         Interactive namespace switcher")
	fmt.Fprintln(os.Stdout)

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Docker:"))
	fmt.Fprintln(os.Stdout, "  fzf_docker_logs    Interactive container logs")
	fmt.Fprintln(os.Stdout, "  fzf_docker_exec    Interactive container exec")

	return nil
}

func runFzfTheme(cmd *cobra.Command, args []string) error {
	helper := fzf.NewHelper(fzfVerbose)

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Catppuccin Mocha Theme"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, helper.GetThemeColors())

	return nil
}

func init() {
	components.Register(&components.Registration{
		Name: "fzf",
		RegisterCmd: func() *cobra.Command { return fzfCmd },
	})
}
