package cmd

import (
	"github.com/mistergrinvalds/acorn/internal/components"
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/ollama"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/mistergrinvalds/acorn/internal/utils/configcmd"
	"github.com/spf13/cobra"
)

var (
	ollamaVerbose bool
	ollamaDryRun  bool
)

// ollamaCmd represents the ollama command group
var ollamaCmd = &cobra.Command{
	Use:   "ollama",
	Short: "Ollama local AI model management",
	Long: `Ollama local AI model management commands.

Provides installation, service management, and model operations.

Examples:
  acorn ollama status      # Show Ollama status
  acorn ollama models      # List installed models
  acorn ollama pull llama3.2  # Pull a model
  acorn ollama chat llama3.2 "Hello"  # Quick chat
  acorn ollama start       # Start service
  acorn ollama stop        # Stop service`,
}

// ollamaStatusCmd shows status
var ollamaStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Ollama status",
	Long: `Display Ollama installation and service status.

Shows version, running status, installed models, and storage usage.

Examples:
  acorn ollama status
  acorn ollama status -o json`,
	RunE: runOllamaStatus,
}

// ollamaModelsCmd lists models
var ollamaModelsCmd = &cobra.Command{
	Use:   "models",
	Short: "List installed models",
	Long: `List all installed Ollama models.

Examples:
  acorn ollama models
  acorn ollama models -o json`,
	Aliases: []string{"list", "ls"},
	RunE:    runOllamaModels,
}

// ollamaPullCmd pulls a model
var ollamaPullCmd = &cobra.Command{
	Use:   "pull <model>",
	Short: "Pull/download a model",
	Long: `Download an Ollama model.

Popular models: llama3.2, codellama, mistral, phi3, gemma2

Examples:
  acorn ollama pull llama3.2
  acorn ollama pull codellama`,
	Args: cobra.ExactArgs(1),
	RunE: runOllamaPull,
}

// ollamaRmCmd removes a model
var ollamaRmCmd = &cobra.Command{
	Use:   "rm <model>",
	Short: "Remove a model",
	Long: `Remove an installed Ollama model.

Examples:
  acorn ollama rm phi3`,
	Aliases: []string{"remove", "delete"},
	Args:    cobra.ExactArgs(1),
	RunE:    runOllamaRm,
}

// ollamaChatCmd chats with a model
var ollamaChatCmd = &cobra.Command{
	Use:   "chat <model> <prompt>",
	Short: "Quick chat with a model",
	Long: `Send a prompt to a model and get a response.

Examples:
  acorn ollama chat llama3.2 "Explain machine learning"
  acorn ollama chat mistral "Write a haiku about coding"`,
	Args: cobra.ExactArgs(2),
	RunE: runOllamaChat,
}

// ollamaCodeCmd generates code
var ollamaCodeCmd = &cobra.Command{
	Use:   "code <language> <description>",
	Short: "Generate code with CodeLlama",
	Long: `Generate code using CodeLlama model.

Will automatically pull CodeLlama if not installed.

Examples:
  acorn ollama code python "function to calculate fibonacci"
  acorn ollama code go "http server with health endpoint"`,
	Args: cobra.ExactArgs(2),
	RunE: runOllamaCode,
}

// ollamaStartCmd starts service
var ollamaStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Ollama service",
	Long: `Start the Ollama service in the background.

Examples:
  acorn ollama start`,
	RunE: runOllamaStart,
}

// ollamaStopCmd stops service
var ollamaStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Ollama service",
	Long: `Stop the running Ollama service.

Examples:
  acorn ollama stop`,
	RunE: runOllamaStop,
}

// ollamaInstallCmd installs Ollama
var ollamaInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Ollama",
	Long: `Download and install Ollama.

On macOS, uses Homebrew if available.
On Linux, uses the official install script.

Examples:
  acorn ollama install`,
	RunE: runOllamaInstall,
}

// ollamaExamplesCmd shows examples
var ollamaExamplesCmd = &cobra.Command{
	Use:   "examples",
	Short: "Show usage examples",
	Long: `Display Ollama usage examples and tips.

Examples:
  acorn ollama examples`,
	RunE: runOllamaExamples,
}

func init() {

	// Add subcommands
	ollamaCmd.AddCommand(ollamaStatusCmd)
	ollamaCmd.AddCommand(ollamaModelsCmd)
	ollamaCmd.AddCommand(ollamaPullCmd)
	ollamaCmd.AddCommand(ollamaRmCmd)
	ollamaCmd.AddCommand(ollamaChatCmd)
	ollamaCmd.AddCommand(ollamaCodeCmd)
	ollamaCmd.AddCommand(ollamaStartCmd)
	ollamaCmd.AddCommand(ollamaStopCmd)
	ollamaCmd.AddCommand(ollamaInstallCmd)
	ollamaCmd.AddCommand(ollamaExamplesCmd)
	ollamaCmd.AddCommand(configcmd.NewConfigRouter("ollama"))

	// Persistent flags
	ollamaCmd.PersistentFlags().BoolVarP(&ollamaVerbose, "verbose", "v", false,
		"Show verbose output")
	ollamaCmd.PersistentFlags().BoolVar(&ollamaDryRun, "dry-run", false,
		"Show what would be done without executing")
}

func runOllamaStatus(cmd *cobra.Command, args []string) error {
	helper := ollama.NewHelper(ollamaVerbose, ollamaDryRun)
	status := helper.GetStatus()

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Ollama Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if status.Installed {
		fmt.Fprintf(os.Stdout, "%s Ollama installed: %s\n", output.Success("✓"), status.Version)
	} else {
		fmt.Fprintf(os.Stdout, "%s Ollama not installed\n", output.Error("✗"))
		fmt.Fprintln(os.Stdout, "  Run: acorn ollama install")
		return nil
	}

	if status.ServiceRunning {
		fmt.Fprintf(os.Stdout, "%s Service: running\n", output.Success("✓"))
	} else {
		fmt.Fprintf(os.Stdout, "%s Service: not running\n", output.Warning("⚠"))
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "Installed models:")
	if len(status.Models) == 0 {
		fmt.Fprintln(os.Stdout, "  (none)")
	} else {
		for _, m := range status.Models {
			fmt.Fprintf(os.Stdout, "  %s", m.Name)
			if m.Size != "" {
				fmt.Fprintf(os.Stdout, " (%s)", m.Size)
			}
			fmt.Fprintln(os.Stdout)
		}
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "Storage: %s (%s)\n", status.HomeDir, status.StorageSize)

	return nil
}

func runOllamaModels(cmd *cobra.Command, args []string) error {
	helper := ollama.NewHelper(ollamaVerbose, ollamaDryRun)

	models, err := helper.ListModels()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(map[string]interface{}{"models": models})
	}

	if len(models) == 0 {
		fmt.Fprintln(os.Stdout, "No models installed")
		fmt.Fprintln(os.Stdout, "Pull a model: acorn ollama pull llama3.2")
		return nil
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Installed Models"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for _, m := range models {
		fmt.Fprintf(os.Stdout, "%-30s %s\n", m.Name, m.Size)
	}

	return nil
}

func runOllamaPull(cmd *cobra.Command, args []string) error {
	helper := ollama.NewHelper(ollamaVerbose, ollamaDryRun)

	if err := helper.Pull(args[0]); err != nil {
		return err
	}

	if !ollamaDryRun {
		fmt.Fprintf(os.Stdout, "%s Model pulled: %s\n", output.Success("✓"), args[0])
	}

	return nil
}

func runOllamaRm(cmd *cobra.Command, args []string) error {
	helper := ollama.NewHelper(ollamaVerbose, ollamaDryRun)

	if err := helper.Remove(args[0]); err != nil {
		return err
	}

	if !ollamaDryRun {
		fmt.Fprintf(os.Stdout, "%s Model removed: %s\n", output.Success("✓"), args[0])
	}

	return nil
}

func runOllamaChat(cmd *cobra.Command, args []string) error {
	helper := ollama.NewHelper(ollamaVerbose, ollamaDryRun)
	return helper.Chat(args[0], args[1])
}

func runOllamaCode(cmd *cobra.Command, args []string) error {
	helper := ollama.NewHelper(ollamaVerbose, ollamaDryRun)
	return helper.Code(args[0], args[1])
}

func runOllamaStart(cmd *cobra.Command, args []string) error {
	helper := ollama.NewHelper(ollamaVerbose, ollamaDryRun)

	if err := helper.Start(); err != nil {
		return err
	}

	if !ollamaDryRun {
		fmt.Fprintf(os.Stdout, "%s Ollama service started\n", output.Success("✓"))
	}

	return nil
}

func runOllamaStop(cmd *cobra.Command, args []string) error {
	helper := ollama.NewHelper(ollamaVerbose, ollamaDryRun)

	if err := helper.Stop(); err != nil {
		return err
	}

	if !ollamaDryRun {
		fmt.Fprintf(os.Stdout, "%s Ollama service stopped\n", output.Success("✓"))
	}

	return nil
}

func runOllamaInstall(cmd *cobra.Command, args []string) error {
	helper := ollama.NewHelper(ollamaVerbose, ollamaDryRun)

	if err := helper.Install(); err != nil {
		return err
	}

	if !ollamaDryRun {
		fmt.Fprintf(os.Stdout, "%s Ollama installed\n", output.Success("✓"))
	}

	return nil
}

func runOllamaExamples(cmd *cobra.Command, args []string) error {
	helper := ollama.NewHelper(ollamaVerbose, ollamaDryRun)
	fmt.Fprintln(os.Stdout, helper.GetExamples())
	return nil
}

func init() {
	components.Register(&components.Registration{
		Name: "ollama",
		RegisterCmd: func() *cobra.Command { return ollamaCmd },
	})
}
