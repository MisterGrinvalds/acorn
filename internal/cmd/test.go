package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/utils/component"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	testCategory    string
	testSkipMissing bool
	testVerbose     bool
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test [component]",
	Short: "Test component configurations",
	Long: `Test component configurations and validate their structure.

Tests include:
  - Config YAML validation
  - Required fields check (name, description, version)
  - Shell script syntax validation
  - Install configuration validation
  - Generated files existence check

Examples:
  acorn test              # Test all components
  acorn test python       # Test specific component
  acorn test --category programming
  acorn test --skip-missing`,
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: completeComponentNames,
	RunE:              runTest,
}

func init() {
	rootCmd.AddCommand(testCmd)

	testCmd.Flags().StringVarP(&testCategory, "category", "c", "",
		"Test only components in this category (ai, cloud, data, devops, ide, programming, terminal, vcs)")
	testCmd.Flags().BoolVar(&testSkipMissing, "skip-missing", false,
		"Skip components that don't exist instead of failing")
	testCmd.Flags().BoolVarP(&testVerbose, "verbose", "v", false,
		"Show verbose test output")
	// Output format is inherited from root command
}

func runTest(cmd *cobra.Command, args []string) error {
	dotfilesRoot, err := getDotfilesRoot()
	if err != nil {
		return err
	}

	tester := component.NewTester(dotfilesRoot,
		component.WithSkipMissing(testSkipMissing),
		component.WithVerbose(testVerbose),
	)

	var results []*component.TestResult

	if len(args) == 1 {
		// Test specific component
		result, err := tester.TestComponent(args[0])
		if err != nil {
			return err
		}
		results = []*component.TestResult{result}
	} else if testCategory != "" {
		// Test by category
		results, err = tester.TestByCategory(testCategory)
		if err != nil {
			return err
		}
	} else {
		// Test all components
		results, err = tester.TestAll()
		if err != nil {
			return err
		}
	}

	if len(results) == 0 {
		fmt.Fprintln(os.Stderr, "No components to test")
		return nil
	}

	// Output results
	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(results)
	}

	// Table format output
	return printTestResults(results)
}

func printTestResults(results []*component.TestResult) error {
	totalPassed := 0
	totalFailed := 0
	totalSkipped := 0
	failedComponents := []string{}

	for _, result := range results {
		// Component header
		statusSymbol := "✓"
		statusColor := output.ColorGreen
		if result.HasFailures() {
			statusSymbol = "✗"
			statusColor = output.ColorRed
			failedComponents = append(failedComponents, result.Component)
		} else if result.Skipped > 0 && result.Passed == 0 {
			statusSymbol = "○"
			statusColor = output.ColorYellow
		}

		fmt.Fprintf(os.Stdout, "%s %s (%s)\n",
			output.Colorize(statusSymbol, statusColor),
			result.Component,
			result.Summary())

		// Show individual tests if verbose or if there are failures
		if testVerbose || result.HasFailures() {
			for _, test := range result.Tests {
				testSymbol := ""
				switch test.Status {
				case component.TestStatusPassed:
					if testVerbose {
						testSymbol = output.Colorize("  ✓", output.ColorGreen)
					} else {
						continue
					}
				case component.TestStatusFailed:
					testSymbol = output.Colorize("  ✗", output.ColorRed)
				case component.TestStatusSkipped:
					if testVerbose {
						testSymbol = output.Colorize("  ○", output.ColorYellow)
					} else {
						continue
					}
				}

				if test.Message != "" {
					fmt.Fprintf(os.Stdout, "%s %s: %s\n", testSymbol, test.Name, test.Message)
				} else {
					fmt.Fprintf(os.Stdout, "%s %s\n", testSymbol, test.Name)
				}
			}
		}

		totalPassed += result.Passed
		totalFailed += result.Failed
		totalSkipped += result.Skipped
	}

	// Summary
	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "Components: %d tested\n", len(results))
	fmt.Fprintf(os.Stdout, "Tests:      %s passed, %s failed, %s skipped\n",
		output.Success(fmt.Sprintf("%d", totalPassed)),
		output.Error(fmt.Sprintf("%d", totalFailed)),
		output.Warning(fmt.Sprintf("%d", totalSkipped)))

	if len(failedComponents) > 0 {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintf(os.Stdout, "%s\n", output.Error("Failed components:"))
		for _, comp := range failedComponents {
			fmt.Fprintf(os.Stdout, "  - %s\n", comp)
		}
		return fmt.Errorf("tests failed")
	}

	return nil
}
