package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/huggingface"
	"github.com/mistergrinvalds/acorn/internal/output"
	"github.com/spf13/cobra"
)

var (
	hfOutputFormat string
	hfVerbose      bool
	hfForce        bool
)

// hfCmd represents the huggingface command group
var hfCmd = &cobra.Command{
	Use:   "hf",
	Short: "Hugging Face model management",
	Long: `Hugging Face model management and helper commands.

Provides status checking, model listings, and cache management.

Examples:
  acorn hf status      # Show HF installation status
  acorn hf models      # List popular models
  acorn hf pipelines   # List available pipelines
  acorn hf cache       # Show cache info
  acorn hf clear       # Clear model cache`,
	Aliases: []string{"huggingface"},
}

// hfStatusCmd shows status
var hfStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Hugging Face installation status",
	Long: `Display Hugging Face installation and environment status.

Shows transformers version, PyTorch version, cache location, and size.

Examples:
  acorn hf status
  acorn hf status -o json`,
	RunE: runHfStatus,
}

// hfModelsCmd lists models
var hfModelsCmd = &cobra.Command{
	Use:   "models",
	Short: "List popular Hugging Face models",
	Long: `Display a list of popular Hugging Face models.

Organized by category: text generation, language understanding, and specialized.

Examples:
  acorn hf models
  acorn hf models -o json`,
	RunE: runHfModels,
}

// hfPipelinesCmd lists pipelines
var hfPipelinesCmd = &cobra.Command{
	Use:   "pipelines",
	Short: "List available pipeline tasks",
	Long: `Display available Hugging Face pipeline tasks.

Shows common tasks like text-generation, summarization, sentiment-analysis, etc.

Examples:
  acorn hf pipelines
  acorn hf pipelines -o json`,
	RunE: runHfPipelines,
}

// hfCacheCmd shows cache info
var hfCacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Show model cache info",
	Long: `Display model cache location and size.

Examples:
  acorn hf cache
  acorn hf cache -o json`,
	RunE: runHfCache,
}

// hfClearCmd clears cache
var hfClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear model cache",
	Long: `Clear the Hugging Face model cache directory.

Requires --force flag to actually delete files.

Examples:
  acorn hf clear --force`,
	RunE: runHfClear,
}

func init() {
	rootCmd.AddCommand(hfCmd)

	// Add subcommands
	hfCmd.AddCommand(hfStatusCmd)
	hfCmd.AddCommand(hfModelsCmd)
	hfCmd.AddCommand(hfPipelinesCmd)
	hfCmd.AddCommand(hfCacheCmd)
	hfCmd.AddCommand(hfClearCmd)

	// Persistent flags
	hfCmd.PersistentFlags().StringVarP(&hfOutputFormat, "output", "o", "table",
		"Output format (table|json|yaml)")
	hfCmd.PersistentFlags().BoolVarP(&hfVerbose, "verbose", "v", false,
		"Show verbose output")

	// Clear command flags
	hfClearCmd.Flags().BoolVar(&hfForce, "force", false,
		"Actually clear the cache (required)")
}

func runHfStatus(cmd *cobra.Command, args []string) error {
	helper := huggingface.NewHelper(hfVerbose)
	status := helper.GetStatus()

	format, err := output.ParseFormat(hfOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Hugging Face Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if status.TransformersInstalled {
		fmt.Fprintf(os.Stdout, "%s Transformers: %s\n", output.Success("✓"), status.TransformersVersion)
	} else {
		fmt.Fprintf(os.Stdout, "%s Transformers: not installed\n", output.Error("✗"))
		fmt.Fprintln(os.Stdout, "  Run: hf_setup")
	}

	if status.PyTorchInstalled {
		fmt.Fprintf(os.Stdout, "%s PyTorch: %s\n", output.Success("✓"), status.PyTorchVersion)
	} else {
		fmt.Fprintf(os.Stdout, "%s PyTorch: not installed\n", output.Warning("⚠"))
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "Cache: %s\n", status.CacheDir)
	if status.CacheSize != "" {
		fmt.Fprintf(os.Stdout, "Size:  %s\n", status.CacheSize)
	}

	if status.VirtualEnv != "" {
		fmt.Fprintf(os.Stdout, "\nVirtual environment: %s\n", status.VirtualEnv)
	}

	return nil
}

func runHfModels(cmd *cobra.Command, args []string) error {
	helper := huggingface.NewHelper(hfVerbose)
	models := helper.GetModels()

	format, err := output.ParseFormat(hfOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(map[string]interface{}{"models": models})
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Popular Hugging Face Models"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	currentCategory := ""
	for _, m := range models {
		if m.Category != currentCategory {
			if currentCategory != "" {
				fmt.Fprintln(os.Stdout)
			}
			fmt.Fprintf(os.Stdout, "%s:\n", output.Info(m.Category))
			currentCategory = m.Category
		}
		fmt.Fprintf(os.Stdout, "  %-30s %s\n", m.Name, m.Description)
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "Usage: Download models using transformers library")

	return nil
}

func runHfPipelines(cmd *cobra.Command, args []string) error {
	helper := huggingface.NewHelper(hfVerbose)
	pipelines := helper.GetPipelines()

	format, err := output.ParseFormat(hfOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(map[string]interface{}{"pipelines": pipelines})
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Hugging Face Pipelines"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "Common pipeline tasks:")

	for _, p := range pipelines {
		fmt.Fprintf(os.Stdout, "  %-22s %s\n", p.Task, p.Description)
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "Example Python usage:")
	fmt.Fprintln(os.Stdout, "  from transformers import pipeline")
	fmt.Fprintln(os.Stdout, "  gen = pipeline('text-generation', model='gpt2')")
	fmt.Fprintln(os.Stdout, "  gen('Hello, I am')")

	return nil
}

func runHfCache(cmd *cobra.Command, args []string) error {
	helper := huggingface.NewHelper(hfVerbose)
	cacheDir, cacheSize, err := helper.GetCacheInfo()
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(hfOutputFormat)
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

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Hugging Face Model Cache"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintf(os.Stdout, "Location: %s\n", cacheDir)
	if cacheSize != "" {
		fmt.Fprintf(os.Stdout, "Size:     %s\n", cacheSize)
	} else {
		fmt.Fprintln(os.Stdout, "Size:     (no cache)")
	}

	return nil
}

func runHfClear(cmd *cobra.Command, args []string) error {
	helper := huggingface.NewHelper(hfVerbose)

	cacheDir, cacheSize, _ := helper.GetCacheInfo()

	if cacheSize == "" {
		fmt.Fprintln(os.Stdout, "No cache directory found")
		return nil
	}

	fmt.Fprintf(os.Stdout, "Cache: %s (%s)\n", cacheDir, cacheSize)

	if err := helper.ClearCache(hfForce); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Cache cleared\n", output.Success("✓"))
	return nil
}
