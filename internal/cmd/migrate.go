package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/mistergrinvalds/acorn/internal/component"
	"github.com/mistergrinvalds/acorn/internal/migrate"
	"github.com/mistergrinvalds/acorn/internal/output"
	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command group
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Analyze and plan component migration to Go",
	Long: `Tools for analyzing shell components and planning their migration to Go.

The migrate commands help you understand which shell functions should be
migrated to Go commands and which should remain as shell scripts.

Migration criteria:
  - Action functions (create, modify, setup) → Go commands
  - Wrapper functions (simple passthrough) → Keep as shell
  - Environment setters → Keep as shell
  - Aliases → Keep as shell

Use 'acorn migrate analyze' to get detailed migration recommendations.`,
	Aliases: []string{"mig"},
}

// migrateAnalyzeCmd analyzes a component for migration
var migrateAnalyzeCmd = &cobra.Command{
	Use:   "analyze [component]",
	Short: "Analyze component for migration opportunities",
	Long: `Analyze a component's shell functions and recommend migration strategy.

For each function, the analyzer determines:
  - Function type (action, wrapper, env, alias)
  - Complexity (low, medium, high)
  - Migration recommendation
  - External command dependencies

Examples:
  acorn migrate analyze python      # Analyze python component
  acorn migrate analyze             # Analyze all components
  acorn migrate analyze -o json     # JSON output for scripting`,
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: completeComponentNames,
	RunE:              runMigrateAnalyze,
}

// migratePlanCmd generates a migration plan
var migratePlanCmd = &cobra.Command{
	Use:   "plan",
	Short: "Generate a prioritized migration plan",
	Long: `Generate a prioritized migration plan across all components.

The plan ranks components by migration value:
  - Number of action functions
  - Complexity of functions
  - Potential benefit from Go features

Examples:
  acorn migrate plan               # Generate plan
  acorn migrate plan -o yaml       # YAML output`,
	RunE: runMigratePlan,
}

// migrateReportCmd generates a full migration report
var migrateReportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate comprehensive migration report",
	Long: `Generate a comprehensive migration report for all components.

The report includes:
  - Summary statistics
  - Component-by-component analysis
  - Prioritized migration order
  - Estimated effort

Examples:
  acorn migrate report             # Display report
  acorn migrate report -o json     # JSON for processing`,
	RunE: runMigrateReport,
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateAnalyzeCmd)
	migrateCmd.AddCommand(migratePlanCmd)
	migrateCmd.AddCommand(migrateReportCmd)

	// Shared output flag
	migrateCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table",
		"Output format (table|json|yaml)")
}

func runMigrateAnalyze(cmd *cobra.Command, args []string) error {
	dotfilesRoot, err := getDotfilesRoot()
	if err != nil {
		return err
	}

	disco := component.NewDiscovery(dotfilesRoot)
	analyzer := migrate.NewAnalyzer(dotfilesRoot)

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

	var analyses []*migrate.ComponentAnalysis
	for _, comp := range components {
		analysis, err := analyzer.AnalyzeComponent(comp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to analyze %s: %v\n", comp.Name, err)
			continue
		}
		analyses = append(analyses, analysis)
	}

	format, err := output.ParseFormat(outputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(analyses)
	}

	// Table format
	for _, analysis := range analyses {
		printComponentAnalysis(analysis)
	}

	return nil
}

func printComponentAnalysis(analysis *migrate.ComponentAnalysis) {
	comp := analysis.Component

	fmt.Fprintf(os.Stdout, "\n%s\n", output.Info(comp.Name))
	fmt.Fprintf(os.Stdout, "%s\n", strings.Repeat("=", len(comp.Name)))
	fmt.Fprintf(os.Stdout, "%s\n\n", comp.Description)

	// Summary
	s := analysis.Summary
	fmt.Fprintf(os.Stdout, "Migration Score: %s\n", formatScore(s.MigrationScore))
	fmt.Fprintf(os.Stdout, "Recommendation:  %s\n\n", s.RecommendedAction)

	// Functions
	if len(analysis.Functions) > 0 {
		fmt.Fprintln(os.Stdout, "Functions:")
		for _, f := range analysis.Functions {
			typeColor := output.ColorGray
			switch f.Type {
			case migrate.FuncTypeAction:
				typeColor = output.ColorGreen
			case migrate.FuncTypeWrapper:
				typeColor = output.ColorYellow
			case migrate.FuncTypeEnv:
				typeColor = output.ColorCyan
			}

			fmt.Fprintf(os.Stdout, "  %s %-20s %s [%s] %d lines\n",
				output.Colorize(string(f.Type[:3]), typeColor),
				f.Name,
				formatComplexity(f.Complexity),
				f.File,
				f.LineCount)

			if f.Type == migrate.FuncTypeAction {
				fmt.Fprintf(os.Stdout, "      → %s\n", f.Suggestion)
			}
		}
		fmt.Fprintln(os.Stdout)
	}

	// Aliases
	if len(analysis.Aliases) > 0 {
		fmt.Fprintf(os.Stdout, "Aliases: %d (keep as shell)\n", len(analysis.Aliases))
		for _, a := range analysis.Aliases {
			fmt.Fprintf(os.Stdout, "  %s = %s\n", a.Name, truncate(a.Command, 50))
		}
		fmt.Fprintln(os.Stdout)
	}

	// Environment variables
	if len(analysis.EnvVars) > 0 {
		fmt.Fprintf(os.Stdout, "Environment Variables: %d (keep as shell)\n", len(analysis.EnvVars))
		for _, e := range analysis.EnvVars {
			fmt.Fprintf(os.Stdout, "  %s\n", e.Name)
		}
		fmt.Fprintln(os.Stdout)
	}
}

func runMigratePlan(cmd *cobra.Command, args []string) error {
	dotfilesRoot, err := getDotfilesRoot()
	if err != nil {
		return err
	}

	disco := component.NewDiscovery(dotfilesRoot)
	analyzer := migrate.NewAnalyzer(dotfilesRoot)

	components, err := disco.DiscoverAll()
	if err != nil {
		return err
	}

	var analyses []*migrate.ComponentAnalysis
	for _, comp := range components {
		analysis, err := analyzer.AnalyzeComponent(comp)
		if err != nil {
			continue
		}
		analyses = append(analyses, analysis)
	}

	// Sort by migration score (descending)
	sort.Slice(analyses, func(i, j int) bool {
		return analyses[i].Summary.MigrationScore > analyses[j].Summary.MigrationScore
	})

	format, err := output.ParseFormat(outputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(analyses)
	}

	// Table format
	fmt.Fprintln(os.Stdout, "\nMigration Plan (Prioritized)")
	fmt.Fprintln(os.Stdout, "============================")

	table := output.NewTable("PRIORITY", "COMPONENT", "SCORE", "FUNCTIONS", "ACTIONS", "RECOMMENDATION")

	priority := 1
	for _, a := range analyses {
		if a.Summary.MigrationScore < 10 {
			continue // Skip very low-score components
		}

		rec := a.Summary.RecommendedAction
		if len(rec) > 40 {
			rec = rec[:37] + "..."
		}

		table.AddRow(
			fmt.Sprintf("%d", priority),
			a.Component.Name,
			fmt.Sprintf("%d%%", a.Summary.MigrationScore),
			fmt.Sprintf("%d", a.Summary.TotalFunctions),
			fmt.Sprintf("%d", a.Summary.ActionFunctions),
			rec,
		)
		priority++
	}

	table.Render(os.Stdout)

	// Summary
	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "Total components to migrate: %d\n", priority-1)

	return nil
}

func runMigrateReport(cmd *cobra.Command, args []string) error {
	dotfilesRoot, err := getDotfilesRoot()
	if err != nil {
		return err
	}

	disco := component.NewDiscovery(dotfilesRoot)
	analyzer := migrate.NewAnalyzer(dotfilesRoot)

	components, err := disco.DiscoverAll()
	if err != nil {
		return err
	}

	// Gather all analyses
	var analyses []*migrate.ComponentAnalysis
	totalFunctions := 0
	totalActions := 0
	totalAliases := 0
	totalEnvVars := 0

	for _, comp := range components {
		analysis, err := analyzer.AnalyzeComponent(comp)
		if err != nil {
			continue
		}
		analyses = append(analyses, analysis)
		totalFunctions += analysis.Summary.TotalFunctions
		totalActions += analysis.Summary.ActionFunctions
		totalAliases += analysis.Summary.TotalAliases
		totalEnvVars += analysis.Summary.TotalEnvVars
	}

	format, err := output.ParseFormat(outputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		report := map[string]interface{}{
			"summary": map[string]interface{}{
				"total_components":   len(analyses),
				"total_functions":    totalFunctions,
				"action_functions":   totalActions,
				"total_aliases":      totalAliases,
				"total_env_vars":     totalEnvVars,
				"migration_coverage": fmt.Sprintf("%.1f%%", float64(totalActions)/float64(totalFunctions)*100),
			},
			"components": analyses,
		}
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(report)
	}

	// Table format report
	fmt.Fprintln(os.Stdout, "\n╔════════════════════════════════════════════════════════════╗")
	fmt.Fprintln(os.Stdout, "║           ACORN MIGRATION REPORT                           ║")
	fmt.Fprintln(os.Stdout, "╚════════════════════════════════════════════════════════════╝")

	fmt.Fprintln(os.Stdout, "SUMMARY")
	fmt.Fprintln(os.Stdout, "-------")
	fmt.Fprintf(os.Stdout, "  Components analyzed:    %d\n", len(analyses))
	fmt.Fprintf(os.Stdout, "  Total functions:        %d\n", totalFunctions)
	fmt.Fprintf(os.Stdout, "  Action functions:       %d (migrate to Go)\n", totalActions)
	fmt.Fprintf(os.Stdout, "  Wrapper functions:      %d (keep as shell)\n", totalFunctions-totalActions)
	fmt.Fprintf(os.Stdout, "  Aliases:                %d (keep as shell)\n", totalAliases)
	fmt.Fprintf(os.Stdout, "  Environment variables:  %d (keep as shell)\n", totalEnvVars)
	fmt.Fprintln(os.Stdout)

	migrationCoverage := float64(totalActions) / float64(totalFunctions) * 100
	fmt.Fprintf(os.Stdout, "  Migration coverage: %.1f%% of functions should move to Go\n\n", migrationCoverage)

	// By category
	fmt.Fprintln(os.Stdout, "BY CATEGORY")
	fmt.Fprintln(os.Stdout, "-----------")

	categoryStats := make(map[string]struct{ total, actions int })
	for _, a := range analyses {
		cat := a.Component.Category
		stats := categoryStats[cat]
		stats.total += a.Summary.TotalFunctions
		stats.actions += a.Summary.ActionFunctions
		categoryStats[cat] = stats
	}

	for cat, stats := range categoryStats {
		pct := 0.0
		if stats.total > 0 {
			pct = float64(stats.actions) / float64(stats.total) * 100
		}
		fmt.Fprintf(os.Stdout, "  %-10s %3d functions, %3d actions (%.0f%% migrate)\n",
			cat, stats.total, stats.actions, pct)
	}
	fmt.Fprintln(os.Stdout)

	// Top candidates
	sort.Slice(analyses, func(i, j int) bool {
		return analyses[i].Summary.MigrationScore > analyses[j].Summary.MigrationScore
	})

	fmt.Fprintln(os.Stdout, "TOP MIGRATION CANDIDATES")
	fmt.Fprintln(os.Stdout, "------------------------")

	for i, a := range analyses {
		if i >= 10 || a.Summary.MigrationScore < 20 {
			break
		}
		fmt.Fprintf(os.Stdout, "  %2d. %-15s %3d%% score, %d action functions\n",
			i+1, a.Component.Name, a.Summary.MigrationScore, a.Summary.ActionFunctions)
	}
	fmt.Fprintln(os.Stdout)

	fmt.Fprintln(os.Stdout, "NEXT STEPS")
	fmt.Fprintln(os.Stdout, "----------")
	fmt.Fprintln(os.Stdout, "  1. Run 'acorn migrate analyze <component>' for detailed analysis")
	fmt.Fprintln(os.Stdout, "  2. Start with high-score components")
	fmt.Fprintln(os.Stdout, "  3. Create Go commands for action functions")
	fmt.Fprintln(os.Stdout, "  4. Keep wrappers/aliases in shell files")
	fmt.Fprintln(os.Stdout)

	return nil
}

func formatScore(score int) string {
	if score >= 70 {
		return output.Success(fmt.Sprintf("%d%% (High)", score))
	} else if score >= 40 {
		return output.Warning(fmt.Sprintf("%d%% (Medium)", score))
	}
	return output.Colorize(fmt.Sprintf("%d%% (Low)", score), output.ColorGray)
}

func formatComplexity(complexity string) string {
	switch complexity {
	case "high":
		return output.Error("HIGH")
	case "medium":
		return output.Warning("MED")
	default:
		return output.Colorize("LOW", output.ColorGray)
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
