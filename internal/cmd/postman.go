package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/data/postman"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	postmanOutputFormat string
	postmanDryRun       bool
	postmanVerbose      bool
	postmanEnvironment  string
	postmanReporters    []string
)

// postmanCmd represents the postman command group
var postmanCmd = &cobra.Command{
	Use:   "postman",
	Short: "Postman API development environment",
	Long: `Postman API development environment commands.

Manage Postman collections, environments, and run API tests with Newman.

Examples:
  acorn data postman status          # Show status
  acorn data postman launch          # Open Postman
  acorn data postman collections     # List collections
  acorn data postman run <collection> # Run collection with Newman`,
	Aliases: []string{"pm"},
}

// postmanStatusCmd shows status
var postmanStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Postman status",
	Long: `Show Postman and Newman installation status.

Examples:
  acorn data postman status
  acorn data postman status -o json`,
	RunE: runPostmanStatus,
}

// postmanLaunchCmd launches Postman
var postmanLaunchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch Postman",
	Long: `Launch Postman application.

Examples:
  acorn data postman launch`,
	Aliases: []string{"open"},
	RunE:    runPostmanLaunch,
}

// postmanInstallNewmanCmd installs Newman
var postmanInstallNewmanCmd = &cobra.Command{
	Use:   "install-newman",
	Short: "Install Newman CLI",
	Long: `Install Newman CLI for running collections from command line.

Examples:
  acorn data postman install-newman`,
	RunE: runPostmanInstallNewman,
}

// postmanRunCmd runs a collection
var postmanRunCmd = &cobra.Command{
	Use:   "run <collection>",
	Short: "Run collection with Newman",
	Long: `Run a Postman collection using Newman CLI.

Examples:
  acorn data postman run ./collection.json
  acorn data postman run ./collection.json -e ./env.json
  acorn data postman run ./collection.json -r cli,json`,
	Args: cobra.ExactArgs(1),
	RunE: runPostmanRun,
}

// postmanCollectionsCmd lists collections
var postmanCollectionsCmd = &cobra.Command{
	Use:   "collections",
	Short: "List collections",
	Long: `List available Postman collections.

Examples:
  acorn data postman collections
  acorn data postman collections -o json`,
	Aliases: []string{"ls"},
	RunE:    runPostmanCollections,
}

// postmanEnvironmentsCmd lists environments
var postmanEnvironmentsCmd = &cobra.Command{
	Use:   "environments",
	Short: "List environments",
	Long: `List available Postman environments.

Examples:
  acorn data postman environments`,
	Aliases: []string{"envs"},
	RunE:    runPostmanEnvironments,
}

// postmanImportCmd imports a collection
var postmanImportCmd = &cobra.Command{
	Use:   "import <file>",
	Short: "Import a collection",
	Long: `Import a Postman collection from a JSON file.

Examples:
  acorn data postman import ./my-collection.json`,
	Args: cobra.ExactArgs(1),
	RunE: runPostmanImport,
}

func init() {
	dataCmd.AddCommand(postmanCmd)

	// Add subcommands
	postmanCmd.AddCommand(postmanStatusCmd)
	postmanCmd.AddCommand(postmanLaunchCmd)
	postmanCmd.AddCommand(postmanInstallNewmanCmd)
	postmanCmd.AddCommand(postmanRunCmd)
	postmanCmd.AddCommand(postmanCollectionsCmd)
	postmanCmd.AddCommand(postmanEnvironmentsCmd)
	postmanCmd.AddCommand(postmanImportCmd)

	// Run flags
	postmanRunCmd.Flags().StringVarP(&postmanEnvironment, "environment", "e", "",
		"Environment file to use")
	postmanRunCmd.Flags().StringSliceVarP(&postmanReporters, "reporters", "r", []string{"cli"},
		"Newman reporters (cli,json,html)")

	// Persistent flags
	postmanCmd.PersistentFlags().StringVarP(&postmanOutputFormat, "output", "o", "table",
		"Output format (table|json|yaml)")
	postmanCmd.PersistentFlags().BoolVar(&postmanDryRun, "dry-run", false,
		"Show what would be done without executing")
	postmanCmd.PersistentFlags().BoolVarP(&postmanVerbose, "verbose", "v", false,
		"Show verbose output")
}

func newPostmanHelper() *postman.Helper {
	return postman.NewHelper(postmanVerbose, postmanDryRun)
}

func runPostmanStatus(cmd *cobra.Command, args []string) error {
	helper := newPostmanHelper()
	status := helper.GetStatus()

	format, err := output.ParseFormat(postmanOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Postman Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if status.Installed {
		version := status.Version
		if version == "" {
			version = "detected"
		}
		fmt.Fprintf(os.Stdout, "%s Postman installed: %s\n", output.Success("✓"), version)
		fmt.Fprintf(os.Stdout, "  Path: %s\n", status.AppPath)
	} else {
		fmt.Fprintf(os.Stdout, "%s Postman not found\n", output.Error("✗"))
		fmt.Fprintln(os.Stdout, "  Install from: https://www.postman.com/downloads/")
	}

	fmt.Fprintln(os.Stdout)

	if status.NewmanInstalled {
		fmt.Fprintf(os.Stdout, "%s Newman CLI: %s\n", output.Success("✓"), status.NewmanVersion)
	} else {
		fmt.Fprintf(os.Stdout, "%s Newman CLI not installed\n", output.Warning("!"))
		fmt.Fprintln(os.Stdout, "  Install: acorn data postman install-newman")
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "Collections dir: %s\n", status.CollectionsDir)

	return nil
}

func runPostmanLaunch(cmd *cobra.Command, args []string) error {
	helper := newPostmanHelper()
	if err := helper.Launch(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Postman launched\n", output.Success("✓"))
	return nil
}

func runPostmanInstallNewman(cmd *cobra.Command, args []string) error {
	helper := newPostmanHelper()
	if err := helper.InstallNewman(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Newman installed\n", output.Success("✓"))
	return nil
}

func runPostmanRun(cmd *cobra.Command, args []string) error {
	helper := newPostmanHelper()
	return helper.RunCollection(args[0], postmanEnvironment, postmanReporters)
}

func runPostmanCollections(cmd *cobra.Command, args []string) error {
	helper := newPostmanHelper()
	collections, err := helper.ListCollections()
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(postmanOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(collections)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Collections"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if len(collections) == 0 {
		fmt.Fprintln(os.Stdout, "No collections found")
		return nil
	}

	for _, c := range collections {
		requests := ""
		if c.Requests > 0 {
			requests = fmt.Sprintf(" (%d requests)", c.Requests)
		}
		fmt.Fprintf(os.Stdout, "  %s%s\n", c.Name, requests)
		fmt.Fprintf(os.Stdout, "    %s\n", c.Path)
	}

	return nil
}

func runPostmanEnvironments(cmd *cobra.Command, args []string) error {
	helper := newPostmanHelper()
	environments, err := helper.ListEnvironments()
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(postmanOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(environments)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Environments"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if len(environments) == 0 {
		fmt.Fprintln(os.Stdout, "No environments found")
		return nil
	}

	for _, e := range environments {
		fmt.Fprintf(os.Stdout, "  %s\n", e.Name)
		fmt.Fprintf(os.Stdout, "    %s\n", e.Path)
	}

	return nil
}

func runPostmanImport(cmd *cobra.Command, args []string) error {
	helper := newPostmanHelper()
	if err := helper.ImportCollection(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Collection imported\n", output.Success("✓"))
	return nil
}
