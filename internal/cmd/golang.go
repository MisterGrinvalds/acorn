package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/programming/golang"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/mistergrinvalds/acorn/internal/utils/configcmd"
	"github.com/spf13/cobra"
)

var (
	goDryRun  bool
	goVerbose bool
)

// goCmd represents the go command group
var goCmd = &cobra.Command{
	Use:   "go",
	Short: "Go development helpers",
	Long: `Helpers for Go development workflow.

Provides commands for project initialization, testing, building,
and Cobra CLI scaffolding.

Examples:
  acorn go new myapp           # Initialize new Go project
  acorn go test                # Run all tests
  acorn go cover               # Run tests with coverage
  acorn go build-all myapp     # Build for all platforms
  acorn go cobra new mycli     # Create Cobra CLI project`,
	Aliases: []string{"golang"},
}

// goNewCmd initializes a new Go project
var goNewCmd = &cobra.Command{
	Use:   "new <module-name>",
	Short: "Initialize a new Go project",
	Long: `Initialize a new Go project with go.mod and main.go.

Creates a directory with the module name, initializes go.mod,
and creates a basic main.go file.

Examples:
  acorn go new myapp
  acorn go new github.com/user/myapp`,
	Args: cobra.ExactArgs(1),
	RunE: runGoNew,
}

// goTestCmd runs tests
var goTestCmd = &cobra.Command{
	Use:   "test [pattern]",
	Short: "Run Go tests",
	Long: `Run Go tests with optional pattern filter.

If no pattern is provided, runs all tests in the project.
Use -run compatible patterns to filter specific tests.

Examples:
  acorn go test                # Run all tests
  acorn go test TestFoo        # Run tests matching "TestFoo"`,
	Args: cobra.MaximumNArgs(1),
	RunE: runGoTest,
}

// goCoverCmd runs tests with coverage
var goCoverCmd = &cobra.Command{
	Use:   "cover",
	Short: "Run tests with coverage report",
	Long: `Run Go tests with coverage profiling.

Generates coverage.out and coverage.html in the current directory,
then displays a coverage summary.

Examples:
  acorn go cover`,
	Aliases: []string{"coverage"},
	RunE:    runGoCover,
}

// goBenchCmd runs benchmarks
var goBenchCmd = &cobra.Command{
	Use:   "bench [pattern]",
	Short: "Run Go benchmarks",
	Long: `Run Go benchmarks with optional pattern filter.

Examples:
  acorn go bench                # Run all benchmarks
  acorn go bench BenchmarkFoo   # Run specific benchmark`,
	Args: cobra.MaximumNArgs(1),
	RunE: runGoBench,
}

// goBuildAllCmd builds for multiple platforms
var goBuildAllCmd = &cobra.Command{
	Use:   "build-all [name]",
	Short: "Build for multiple platforms",
	Long: `Cross-compile for Linux, macOS, and Windows.

Builds for:
  - linux/amd64
  - linux/arm64
  - darwin/amd64
  - darwin/arm64
  - windows/amd64

Output goes to dist/ directory.

Examples:
  acorn go build-all           # Build as "app"
  acorn go build-all myapp     # Build as "myapp"`,
	Args: cobra.MaximumNArgs(1),
	RunE: runGoBuildAll,
}

// goCleanCmd cleans build artifacts
var goCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean build artifacts",
	Long: `Remove build artifacts including dist/, coverage files, and go clean.

Examples:
  acorn go clean`,
	RunE: runGoClean,
}

// goEnvCmd shows Go environment
var goEnvCmd = &cobra.Command{
	Use:   "env",
	Short: "Show Go environment information",
	Long: `Display Go environment variables and configuration.

Examples:
  acorn go env
  acorn go env -o json`,
	RunE: runGoEnv,
}

// goCobraCmd is the parent for Cobra-related subcommands
var goCobraCmd = &cobra.Command{
	Use:   "cobra",
	Short: "Cobra CLI scaffolding commands",
	Long: `Commands for scaffolding Cobra CLI applications.

Examples:
  acorn go cobra new mycli     # Create new Cobra project
  acorn go cobra add serve     # Add command to existing project`,
}

// goCobraNewCmd creates a new Cobra project
var goCobraNewCmd = &cobra.Command{
	Use:   "new <app-name>",
	Short: "Create a new Cobra CLI project",
	Long: `Initialize a new Cobra CLI project.

Installs cobra-cli if not present, then creates a new CLI project
with the standard Cobra structure.

Examples:
  acorn go cobra new mycli`,
	Args: cobra.ExactArgs(1),
	RunE: runGoCobraNew,
}

// goCobraAddCmd adds a command to a Cobra project
var goCobraAddCmd = &cobra.Command{
	Use:   "add <command-name>",
	Short: "Add a command to a Cobra project",
	Long: `Add a new command to an existing Cobra project.

Must be run from within a Cobra project directory.

Examples:
  acorn go cobra add serve
  acorn go cobra add config`,
	Args: cobra.ExactArgs(1),
	RunE: runGoCobraAdd,
}

func init() {
	programmingCmd.AddCommand(goCmd)

	// Add subcommands
	goCmd.AddCommand(goNewCmd)
	goCmd.AddCommand(goTestCmd)
	goCmd.AddCommand(goCoverCmd)
	goCmd.AddCommand(goBenchCmd)
	goCmd.AddCommand(goBuildAllCmd)
	goCmd.AddCommand(goCleanCmd)
	goCmd.AddCommand(goEnvCmd)
	goCmd.AddCommand(goCobraCmd)
	goCmd.AddCommand(configcmd.NewConfigRouter("go"))

	// Cobra subcommands
	goCobraCmd.AddCommand(goCobraNewCmd)
	goCobraCmd.AddCommand(goCobraAddCmd)

	// Persistent flags
	goCmd.PersistentFlags().BoolVar(&goDryRun, "dry-run", false,
		"Show what would be done without executing")
	goCmd.PersistentFlags().BoolVarP(&goVerbose, "verbose", "v", false,
		"Show verbose output")
}

func runGoNew(cmd *cobra.Command, args []string) error {
	helper := golang.NewHelper(goVerbose, goDryRun)
	project, err := helper.InitProject(args[0])
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Go project initialized!\n", output.Success("✓"))
	fmt.Fprintf(os.Stdout, "  Path: %s\n", project.Path)
	fmt.Fprintf(os.Stdout, "  Module: %s\n", project.Module)
	fmt.Fprintf(os.Stdout, "\nNext steps:\n")
	fmt.Fprintf(os.Stdout, "  cd %s\n", project.Name)
	fmt.Fprintf(os.Stdout, "  go run .\n")

	return nil
}

func runGoTest(cmd *cobra.Command, args []string) error {
	helper := golang.NewHelper(goVerbose, goDryRun)

	pattern := ""
	if len(args) > 0 {
		pattern = args[0]
	}

	fmt.Fprintln(os.Stdout, "Running tests...")
	return helper.RunTests(pattern)
}

func runGoCover(cmd *cobra.Command, args []string) error {
	helper := golang.NewHelper(goVerbose, goDryRun)

	fmt.Fprintln(os.Stdout, "Running tests with coverage...")
	if err := helper.RunTestsWithCoverage(); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "\n%s Coverage report generated: coverage.html\n", output.Success("✓"))
	return nil
}

func runGoBench(cmd *cobra.Command, args []string) error {
	helper := golang.NewHelper(goVerbose, goDryRun)

	pattern := ""
	if len(args) > 0 {
		pattern = args[0]
	}

	fmt.Fprintln(os.Stdout, "Running benchmarks...")
	return helper.RunBenchmarks(pattern)
}

func runGoBuildAll(cmd *cobra.Command, args []string) error {
	helper := golang.NewHelper(goVerbose, goDryRun)

	name := "app"
	if len(args) > 0 {
		name = args[0]
	}

	fmt.Fprintf(os.Stdout, "Building %s for multiple platforms...\n\n", output.Info(name))
	if err := helper.BuildAll(name); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "\n%s Builds complete in dist/\n", output.Success("✓"))
	return nil
}

func runGoClean(cmd *cobra.Command, args []string) error {
	helper := golang.NewHelper(goVerbose, goDryRun)

	fmt.Fprintln(os.Stdout, "Cleaning build artifacts...")
	if err := helper.Clean(); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Build artifacts cleaned\n", output.Success("✓"))
	return nil
}

func runGoEnv(cmd *cobra.Command, args []string) error {
	helper := golang.NewHelper(goVerbose, goDryRun)

	version, err := helper.GetGoVersion()
	if err != nil {
		return fmt.Errorf("go not installed: %w", err)
	}

	env := helper.GetGoEnv()

	// Add version to env for JSON/YAML output
	env["VERSION"] = version

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(env)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n\n", output.Info("Go Environment"))
	fmt.Fprintf(os.Stdout, "Version:    %s\n", version)
	fmt.Fprintf(os.Stdout, "GOOS:       %s\n", env["GOOS"])
	fmt.Fprintf(os.Stdout, "GOARCH:     %s\n", env["GOARCH"])
	fmt.Fprintf(os.Stdout, "GOPATH:     %s\n", env["GOPATH"])
	fmt.Fprintf(os.Stdout, "GOROOT:     %s\n", env["GOROOT"])
	fmt.Fprintf(os.Stdout, "GOBIN:      %s\n", env["GOBIN"])
	fmt.Fprintf(os.Stdout, "GOPROXY:    %s\n", env["GOPROXY"])
	fmt.Fprintf(os.Stdout, "GOMODCACHE: %s\n", env["GOMODCACHE"])

	return nil
}

func runGoCobraNew(cmd *cobra.Command, args []string) error {
	helper := golang.NewHelper(goVerbose, goDryRun)

	project, err := helper.InitCobraProject(args[0])
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Cobra CLI project initialized!\n", output.Success("✓"))
	fmt.Fprintf(os.Stdout, "  Path: %s\n", project.Path)
	fmt.Fprintf(os.Stdout, "\nNext steps:\n")
	fmt.Fprintf(os.Stdout, "  cd %s\n", project.Name)
	fmt.Fprintf(os.Stdout, "  go run . --help\n")

	return nil
}

func runGoCobraAdd(cmd *cobra.Command, args []string) error {
	helper := golang.NewHelper(goVerbose, goDryRun)

	if err := helper.AddCobraCommand(args[0]); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Command '%s' added\n", output.Success("✓"), args[0])
	return nil
}
