package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mistergrinvalds/acorn/internal/utils/agentic"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

// agenticCmd represents the agentic command group
var agenticCmd = &cobra.Command{
	Use:   "agentic",
	Short: "Audit commands for JSON-first compliance",
	Long: `Check and enforce agentic/JSON-first best practices for CLI commands.

Commands should nearly always accept and output JSON data, with human
formatting as a view layer. This command helps identify compliance gaps.

Examples:
  acorn dev agentic audit                           # Audit all commands
  acorn dev agentic audit -o json                   # JSON output
  acorn dev agentic validate internal/cmd/tools.go  # Validate specific file
  acorn dev agentic patterns                        # Show expected patterns`,
}

// agenticAuditCmd audits all command files
var agenticAuditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Check all commands for JSON output compliance",
	Long: `Scan all command files in internal/cmd/ and check for JSON-first compliance.

Checks for:
  - --output/-o flag presence
  - output.ParseFormat() usage
  - output.NewPrinter() usage
  - Direct stdout bypasses

Examples:
  acorn dev agentic audit
  acorn dev agentic audit -o json
  acorn dev agentic audit -o yaml`,
	RunE: runAgenticAudit,
}

// agenticValidateCmd validates a specific file
var agenticValidateCmd = &cobra.Command{
	Use:   "validate <file>",
	Short: "Deep validation of specific command file",
	Long: `Perform deep validation of a specific command file.

Provides detailed analysis including:
  - Output flag configuration
  - Format branching patterns
  - Printer usage locations
  - All stdout bypass locations with line numbers

Examples:
  acorn dev agentic validate internal/cmd/tools.go
  acorn dev agentic validate internal/cmd/sapling.go -o json`,
	Args: cobra.ExactArgs(1),
	RunE: runAgenticValidate,
}

// agenticPatternsCmd shows expected patterns
var agenticPatternsCmd = &cobra.Command{
	Use:   "patterns",
	Short: "Show expected patterns and templates",
	Long: `Display the expected code patterns for agentic-compliant commands.

Shows templates for:
  - Output flag declaration
  - Format parsing
  - Format switch handling
  - Table output
  - Structured data types

Examples:
  acorn dev agentic patterns
  acorn dev agentic patterns -o json`,
	RunE: runAgenticPatterns,
}

func init() {
	devCmd.AddCommand(agenticCmd)

	// Add subcommands
	agenticCmd.AddCommand(agenticAuditCmd)
	agenticCmd.AddCommand(agenticValidateCmd)
	agenticCmd.AddCommand(agenticPatternsCmd)

}

func runAgenticAudit(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	// Find the internal/cmd directory
	cmdDir, err := findCmdDir()
	if err != nil {
		return err
	}

	analyzer := agentic.NewAnalyzer(cmdDir)
	result, err := analyzer.AuditAll()
	if err != nil {
		return fmt.Errorf("audit failed: %w", err)
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(result)
	}

	// Table format
	table := output.NewTable("COMMAND", "SUBS", "COMPLIANT", "-o FLAG", "BYPASSES", "ISSUES")
	for _, audit := range result.Commands {
		compliant := output.Error("NO")
		if audit.Compliant {
			compliant = output.Success("YES")
		}

		hasFlag := output.Error("NO")
		if audit.HasOutputFlag {
			hasFlag = output.Success("YES")
		}

		bypassStr := strconv.Itoa(audit.BypassCount)
		if audit.BypassCount > 10 {
			bypassStr = output.Warning(bypassStr)
		}

		issueCount := strconv.Itoa(len(audit.Issues))
		if len(audit.Issues) > 0 {
			issueCount = output.Error(issueCount)
		}

		table.AddRow(
			audit.Command,
			strconv.Itoa(len(audit.Subcommands)),
			compliant,
			hasFlag,
			bypassStr,
			issueCount,
		)
	}
	table.Render(os.Stdout)

	// Summary
	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "Summary: %d commands, %s compliant (%.1f%%)\n",
		result.Summary.Total,
		formatCompliant(result.Summary.Compliant, result.Summary.Total),
		result.Summary.Percentage)

	if result.Summary.TotalIssues > 0 {
		fmt.Fprintf(os.Stdout, "Total issues: %s\n", output.Error(strconv.Itoa(result.Summary.TotalIssues)))
	}

	return nil
}

func runAgenticValidate(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	filePath := args[0]

	// Make path absolute if relative
	if !filepath.IsAbs(filePath) {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		filePath = filepath.Join(cwd, filePath)
	}

	// Check file exists
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("file not found: %s", filePath)
	}

	cmdDir := filepath.Dir(filePath)
	analyzer := agentic.NewAnalyzer(cmdDir)
	result, err := analyzer.ValidateFile(filePath)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(result)
	}

	// Table format - detailed view
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("File: "+result.File))
	fmt.Fprintf(os.Stdout, "%s\n\n", strings.Repeat("=", len("File: "+result.File)))

	// Validity status
	if result.Valid {
		fmt.Fprintf(os.Stdout, "Status: %s\n\n", output.Success("COMPLIANT"))
	} else {
		fmt.Fprintf(os.Stdout, "Status: %s\n\n", output.Error("NON-COMPLIANT"))
	}

	// Output flag
	fmt.Fprintln(os.Stdout, output.Info("Output Flag:"))
	if result.OutputFlag != nil {
		fmt.Fprintf(os.Stdout, "  %s Found at line %d\n", output.Success("✓"), result.OutputFlag.Line)
	} else {
		fmt.Fprintf(os.Stdout, "  %s Not found\n", output.Error("✗"))
	}

	// Format branches
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, output.Info("Format Branches:"))
	if len(result.FormatBranches) > 0 {
		for _, branch := range result.FormatBranches {
			fmt.Fprintf(os.Stdout, "  %s Line %d\n", output.Success("✓"), branch.Line)
		}
	} else {
		fmt.Fprintf(os.Stdout, "  %s No format branches found\n", output.Warning("!"))
	}

	// Printer usage
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, output.Info("Printer Usage:"))
	if len(result.PrinterUsage) > 0 {
		for _, usage := range result.PrinterUsage {
			fmt.Fprintf(os.Stdout, "  %s Line %d: %s\n", output.Success("✓"), usage.Line, usage.Function)
		}
	} else {
		fmt.Fprintf(os.Stdout, "  %s No printer usage found\n", output.Warning("!"))
	}

	// Bypasses
	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "%s (%d total)\n", output.Info("Stdout Bypasses:"), len(result.Bypasses))
	if len(result.Bypasses) > 0 {
		// Show first 10
		limit := min(len(result.Bypasses), 10)
		for i := range limit {
			bypass := result.Bypasses[i]
			inBranch := ""
			if bypass.InFormat {
				inBranch = " (in format branch)"
			}
			fmt.Fprintf(os.Stdout, "  Line %d%s\n", bypass.Line, inBranch)
		}
		if len(result.Bypasses) > 10 {
			fmt.Fprintf(os.Stdout, "  ... and %d more\n", len(result.Bypasses)-10)
		}
	}

	// Issues
	if len(result.Issues) > 0 {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, output.Error("Issues:"))
		for _, issue := range result.Issues {
			fmt.Fprintf(os.Stdout, "  %s %s\n", output.Error("✗"), issue)
		}
	}

	// Recommendations
	if len(result.Recommendations) > 0 {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, output.Info("Recommendations:"))
		for _, rec := range result.Recommendations {
			fmt.Fprintf(os.Stdout, "  %s %s\n", output.Info("→"), rec)
		}
	}

	return nil
}

func runAgenticPatterns(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	patterns := agentic.GetPatterns()

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(patterns)
	}

	// Table format - show each pattern
	printPattern := func(p agentic.PatternExample) {
		fmt.Fprintf(os.Stdout, "%s\n", output.Info(p.Name))
		fmt.Fprintf(os.Stdout, "%s\n\n", p.Description)
		fmt.Fprintf(os.Stdout, "%s\n\n", p.Code)
	}

	fmt.Fprintln(os.Stdout, "Expected Patterns for Agentic-Compliant Commands")
	fmt.Fprintln(os.Stdout, strings.Repeat("=", 50))
	fmt.Fprintln(os.Stdout)

	printPattern(patterns.OutputFlag)
	printPattern(patterns.ParseFormat)
	printPattern(patterns.FormatSwitch)
	printPattern(patterns.TableOutput)
	printPattern(patterns.StructuredData)

	return nil
}

// findCmdDir locates the internal/cmd directory
func findCmdDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Check if we're in the repo root
	cmdDir := filepath.Join(cwd, "internal", "cmd")
	if _, err := os.Stat(cmdDir); err == nil {
		return cmdDir, nil
	}

	// Check parent directories
	dir := cwd
	for range 5 {
		cmdDir = filepath.Join(dir, "internal", "cmd")
		if _, err := os.Stat(cmdDir); err == nil {
			return cmdDir, nil
		}
		dir = filepath.Dir(dir)
	}

	return "", fmt.Errorf("could not find internal/cmd directory")
}

// formatCompliant formats the compliant count with color
func formatCompliant(compliant, total int) string {
	if total == 0 {
		return "0"
	}
	pct := float64(compliant) / float64(total) * 100
	s := strconv.Itoa(compliant)
	if pct >= 80 {
		return output.Success(s)
	} else if pct >= 50 {
		return output.Warning(s)
	}
	return output.Error(s)
}
