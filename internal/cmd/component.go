package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mistergrinvalds/acorn/internal/component"
	"github.com/mistergrinvalds/acorn/internal/output"
	"github.com/spf13/cobra"
)

var (
	outputFormat string
)

// componentCmd represents the component command group
var componentCmd = &cobra.Command{
	Use:   "component",
	Short: "Manage dotfiles components",
	Long: `Manage and inspect dotfiles components.

Components are self-contained feature modules with metadata defined in
component.yaml files. Each component can provide shell functions, aliases,
environment variables, and configuration files.

Use the subcommands to list, inspect, validate, and check the health of
your dotfiles components.`,
	Aliases: []string{"comp", "components"},
}

// componentListCmd lists all components
var componentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all components",
	Long: `List all available components with their metadata.

By default, displays components in a table format showing name, version,
category, and description. Use --output to change format.

Examples:
  acorn component list
  acorn component list --output json
  acorn component list -o yaml`,
	Aliases: []string{"ls"},
	RunE:    runComponentList,
}

// componentStatusCmd shows component health status
var componentStatusCmd = &cobra.Command{
	Use:   "status [component]",
	Short: "Check health status of components",
	Long: `Check the health status of all components or a specific component.

Health checks include:
  - YAML validation
  - Required tool availability
  - Shell script syntax validation
  - Configuration file existence
  - Platform compatibility

Examples:
  acorn component status           # Check all components
  acorn component status python    # Check specific component
  acorn component status -o json   # JSON output`,
	Args:    cobra.MaximumNArgs(1),
	RunE:    runComponentStatus,
}

// componentValidateCmd validates component configurations
var componentValidateCmd = &cobra.Command{
	Use:   "validate [component]",
	Short: "Validate component configurations",
	Long: `Validate component YAML files and configurations.

Performs strict validation including:
  - Required fields (name, version, description, category)
  - Config file existence
  - Platform and shell values
  - Shell script syntax
  - Config method values (symlink/copy)

Examples:
  acorn component validate         # Validate all
  acorn component validate python  # Validate specific component`,
	Args:    cobra.MaximumNArgs(1),
	RunE:    runComponentValidate,
}

// componentInfoCmd shows detailed info about a component
var componentInfoCmd = &cobra.Command{
	Use:   "info <component>",
	Short: "Show detailed information about a component",
	Long: `Display comprehensive information about a specific component.

Shows all metadata from component.yaml including dependencies, provided
features, configuration files, and XDG directory usage.

Examples:
  acorn component info python
  acorn component info git --output yaml`,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeComponentNames,
	RunE:              runComponentInfo,
}

func init() {
	rootCmd.AddCommand(componentCmd)

	// Add subcommands
	componentCmd.AddCommand(componentListCmd)
	componentCmd.AddCommand(componentStatusCmd)
	componentCmd.AddCommand(componentValidateCmd)
	componentCmd.AddCommand(componentInfoCmd)

	// Global flags for component commands
	componentCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table",
		"Output format (table|json|yaml)")
}

// getDotfilesRoot returns the dotfiles root directory
func getDotfilesRoot() (string, error) {
	// First check environment variable
	if root := os.Getenv("DOTFILES_ROOT"); root != "" {
		return root, nil
	}

	// Check config
	if cfg != nil && cfg.DotfilesRoot != "" {
		return cfg.DotfilesRoot, nil
	}

	// Fallback to relative path (for development)
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	// Check if we're in the repo
	if _, err := os.Stat(filepath.Join(cwd, "components")); err == nil {
		return cwd, nil
	}

	return "", fmt.Errorf("DOTFILES_ROOT not set and components directory not found")
}

// runComponentList executes the list command
func runComponentList(cmd *cobra.Command, args []string) error {
	dotfilesRoot, err := getDotfilesRoot()
	if err != nil {
		return err
	}

	disco := component.NewDiscovery(dotfilesRoot)
	components, err := disco.DiscoverAll()
	if err != nil {
		return fmt.Errorf("failed to discover components: %w", err)
	}

	if len(components) == 0 {
		fmt.Fprintln(os.Stderr, "No components found")
		return nil
	}

	format, err := output.ParseFormat(outputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(components)
	}

	// Table format
	table := output.NewTable("NAME", "VERSION", "CATEGORY", "DESCRIPTION")
	for _, comp := range components {
		desc := comp.Description
		if len(desc) > 50 {
			desc = desc[:47] + "..."
		}
		table.AddRow(comp.Name, comp.Version, comp.Category, desc)
	}
	table.Render(os.Stdout)

	fmt.Fprintf(os.Stdout, "\nTotal: %d components\n", len(components))
	return nil
}

// runComponentStatus executes the status command
func runComponentStatus(cmd *cobra.Command, args []string) error {
	dotfilesRoot, err := getDotfilesRoot()
	if err != nil {
		return err
	}

	disco := component.NewDiscovery(dotfilesRoot)

	var components []*component.Component
	if len(args) == 1 {
		// Check specific component
		comp, err := disco.FindByName(args[0])
		if err != nil {
			return err
		}
		components = []*component.Component{comp}
	} else {
		// Check all components
		components, err = disco.DiscoverAll()
		if err != nil {
			return err
		}
	}

	// Perform health checks
	var results []*component.HealthCheck
	for _, comp := range components {
		hc := component.CheckHealth(comp)
		results = append(results, hc)
	}

	format, err := output.ParseFormat(outputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(results)
	}

	// Table format with color-coded status
	healthy := 0
	warnings := 0
	errors := 0

	for _, hc := range results {
		statusSymbol := ""
		statusColor := output.ColorGreen

		switch hc.Status {
		case component.StatusHealthy:
			statusSymbol = "✓"
			statusColor = output.ColorGreen
			healthy++
		case component.StatusWarning:
			statusSymbol = "⚠"
			statusColor = output.ColorYellow
			warnings++
		case component.StatusError:
			statusSymbol = "✗"
			statusColor = output.ColorRed
			errors++
		}

		fmt.Fprintf(os.Stdout, "%s %s - %s\n",
			output.Colorize(statusSymbol, statusColor),
			hc.Component.Name,
			hc.Component.Description)

		// Print issues
		for _, issue := range hc.Issues {
			fmt.Fprintf(os.Stdout, "  %s %s\n", output.Error("✗"), issue)
		}
		for _, warning := range hc.Warnings {
			fmt.Fprintf(os.Stdout, "  %s %s\n", output.Warning("!"), warning)
		}
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "Summary: %d components\n", len(results))
	fmt.Fprintf(os.Stdout, "  %s: %d\n", output.Success("Healthy"), healthy)
	if warnings > 0 {
		fmt.Fprintf(os.Stdout, "  %s: %d\n", output.Warning("Warnings"), warnings)
	}
	if errors > 0 {
		fmt.Fprintf(os.Stdout, "  %s: %d\n", output.Error("Errors"), errors)
	}

	return nil
}

// runComponentValidate executes the validate command
func runComponentValidate(cmd *cobra.Command, args []string) error {
	dotfilesRoot, err := getDotfilesRoot()
	if err != nil {
		return err
	}

	disco := component.NewDiscovery(dotfilesRoot)

	var components []*component.Component
	if len(args) == 1 {
		comp, err := disco.FindByName(args[0])
		if err != nil {
			return err
		}
		components = []*component.Component{comp}
	} else {
		components, err = disco.DiscoverAll()
		if err != nil {
			return err
		}
	}

	// Validate each component
	var results []*component.ValidationResult
	for _, comp := range components {
		vr := component.Validate(comp)
		results = append(results, vr)
	}

	format, err := output.ParseFormat(outputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(results)
	}

	// Table format
	validCount := 0
	invalidCount := 0

	for _, vr := range results {
		if vr.Valid {
			fmt.Fprintf(os.Stdout, "%s %s (%s)\n",
				output.Success("✓"),
				vr.Component.Name,
				vr.Component.Version)
			validCount++
		} else {
			fmt.Fprintf(os.Stdout, "%s %s\n",
				output.Error("✗"),
				vr.Component.Name)
			for _, err := range vr.Errors {
				fmt.Fprintf(os.Stdout, "  - %s\n", err)
			}
			invalidCount++
		}
	}

	fmt.Fprintln(os.Stdout)
	if invalidCount == 0 {
		fmt.Fprintln(os.Stdout, output.Success("All components are valid"))
		return nil
	}

	fmt.Fprintf(os.Stdout, "%s: %d valid, %s: %d invalid\n",
		output.Success("Valid"), validCount,
		output.Error("Invalid"), invalidCount)

	return fmt.Errorf("validation failed")
}

// runComponentInfo executes the info command
func runComponentInfo(cmd *cobra.Command, args []string) error {
	dotfilesRoot, err := getDotfilesRoot()
	if err != nil {
		return err
	}

	disco := component.NewDiscovery(dotfilesRoot)
	comp, err := disco.FindByName(args[0])
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(outputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(comp)
	}

	// Detailed table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info(comp.Name))
	fmt.Fprintf(os.Stdout, "%s\n\n", strings.Repeat("=", len(comp.Name)))

	fmt.Fprintf(os.Stdout, "Description: %s\n", comp.Description)
	fmt.Fprintf(os.Stdout, "Version:     %s\n", comp.Version)
	fmt.Fprintf(os.Stdout, "Category:    %s\n", comp.Category)
	fmt.Fprintf(os.Stdout, "Path:        %s\n", comp.Path)

	if len(comp.Platforms) > 0 {
		fmt.Fprintf(os.Stdout, "Platforms:   %v\n", comp.Platforms)
	}
	if len(comp.Shells) > 0 {
		fmt.Fprintf(os.Stdout, "Shells:      %v\n", comp.Shells)
	}

	if len(comp.Requires.Tools) > 0 {
		fmt.Fprintln(os.Stdout, "\nRequired Tools:")
		for _, tool := range comp.Requires.Tools {
			exists := commandExists(tool)
			if exists {
				fmt.Fprintf(os.Stdout, "  %s %s\n", output.Success("✓"), tool)
			} else {
				fmt.Fprintf(os.Stdout, "  %s %s (not installed)\n", output.Warning("✗"), tool)
			}
		}
	}

	if len(comp.Requires.Components) > 0 {
		fmt.Fprintf(os.Stdout, "\nDependencies: %v\n", comp.Requires.Components)
	}

	if len(comp.Provides.Functions) > 0 {
		fmt.Fprintf(os.Stdout, "\nProvides Functions: %v\n", comp.Provides.Functions)
	}
	if len(comp.Provides.Aliases) > 0 {
		fmt.Fprintf(os.Stdout, "Provides Aliases:   %v\n", comp.Provides.Aliases)
	}

	if len(comp.Config.Files) > 0 {
		fmt.Fprintln(os.Stdout, "\nConfiguration Files:")
		for _, cfg := range comp.Config.Files {
			method := cfg.Method
			if method == "" {
				method = "symlink"
			}
			fmt.Fprintf(os.Stdout, "  %s → %s (%s)\n", cfg.Source, cfg.Target, method)
		}
	}

	// Show compatibility for current platform
	currentPlatform := runtime.GOOS
	compatible := comp.SupportsCurrentPlatform(currentPlatform)
	fmt.Fprintln(os.Stdout)
	if compatible {
		fmt.Fprintf(os.Stdout, "Platform: %s %s\n",
			output.Success("Compatible with"), currentPlatform)
	} else {
		fmt.Fprintf(os.Stdout, "Platform: %s %s\n",
			output.Warning("Not compatible with"), currentPlatform)
	}

	return nil
}

// completeComponentNames provides completion for component names
func completeComponentNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	dotfilesRoot, err := getDotfilesRoot()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	disco := component.NewDiscovery(dotfilesRoot)
	components, err := disco.DiscoverAll()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string
	for _, comp := range components {
		names = append(names, comp.Name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

// commandExists checks if a command exists in PATH
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
