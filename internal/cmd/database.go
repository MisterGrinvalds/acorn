package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/database"
	"github.com/mistergrinvalds/acorn/internal/output"
	"github.com/spf13/cobra"
)

var (
	dbOutputFormat string
	dbDryRun       bool
	dbVerbose      bool
)

// dbCmd represents the database command group
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database service management",
	Long: `Manage database services and check their status.

Provides commands for checking database service status and managing
services via Homebrew on macOS.

Examples:
  acorn db status                # Check all database statuses
  acorn db start postgres        # Start PostgreSQL
  acorn db stop redis            # Stop Redis
  acorn db start-all             # Start common databases`,
	Aliases: []string{"database"},
}

// dbStatusCmd shows database status
var dbStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check database service status",
	Long: `Check the status of all database services.

Shows whether each database is installed and running.

Examples:
  acorn db status
  acorn db status -o json`,
	RunE: runDbStatus,
}

// dbStartCmd starts a database service
var dbStartCmd = &cobra.Command{
	Use:   "start <service>",
	Short: "Start a database service",
	Long: `Start a database service via Homebrew (macOS only).

Supported services: postgres, mysql, mongodb, redis, neo4j, kafka, zookeeper

Examples:
  acorn db start postgres
  acorn db start redis`,
	Args: cobra.ExactArgs(1),
	RunE: runDbStart,
}

// dbStopCmd stops a database service
var dbStopCmd = &cobra.Command{
	Use:   "stop <service>",
	Short: "Stop a database service",
	Long: `Stop a database service via Homebrew (macOS only).

Supported services: postgres, mysql, mongodb, redis, neo4j, kafka, zookeeper

Examples:
  acorn db stop postgres
  acorn db stop redis`,
	Args: cobra.ExactArgs(1),
	RunE: runDbStop,
}

// dbRestartCmd restarts a database service
var dbRestartCmd = &cobra.Command{
	Use:   "restart <service>",
	Short: "Restart a database service",
	Long: `Restart a database service via Homebrew (macOS only).

Supported services: postgres, mysql, mongodb, redis, neo4j, kafka, zookeeper

Examples:
  acorn db restart postgres
  acorn db restart redis`,
	Args: cobra.ExactArgs(1),
	RunE: runDbRestart,
}

// dbStartAllCmd starts all common databases
var dbStartAllCmd = &cobra.Command{
	Use:   "start-all",
	Short: "Start common database services",
	Long: `Start common database services via Homebrew (macOS only).

Starts PostgreSQL, Redis, and MongoDB.

Examples:
  acorn db start-all`,
	RunE: runDbStartAll,
}

// dbStopAllCmd stops all common databases
var dbStopAllCmd = &cobra.Command{
	Use:   "stop-all",
	Short: "Stop common database services",
	Long: `Stop common database services via Homebrew (macOS only).

Stops PostgreSQL, Redis, and MongoDB.

Examples:
  acorn db stop-all`,
	RunE: runDbStopAll,
}

// dbListCmd lists supported services
var dbListCmd = &cobra.Command{
	Use:   "list",
	Short: "List supported database services",
	Long: `List all supported database services.

Examples:
  acorn db list`,
	RunE: runDbList,
}

func init() {
	rootCmd.AddCommand(dbCmd)

	// Add subcommands
	dbCmd.AddCommand(dbStatusCmd)
	dbCmd.AddCommand(dbStartCmd)
	dbCmd.AddCommand(dbStopCmd)
	dbCmd.AddCommand(dbRestartCmd)
	dbCmd.AddCommand(dbStartAllCmd)
	dbCmd.AddCommand(dbStopAllCmd)
	dbCmd.AddCommand(dbListCmd)

	// Persistent flags
	dbCmd.PersistentFlags().StringVarP(&dbOutputFormat, "output", "o", "table",
		"Output format (table|json|yaml)")
	dbCmd.PersistentFlags().BoolVar(&dbDryRun, "dry-run", false,
		"Show what would be done without executing")
	dbCmd.PersistentFlags().BoolVarP(&dbVerbose, "verbose", "v", false,
		"Show verbose output")
}

func runDbStatus(cmd *cobra.Command, args []string) error {
	helper := database.NewHelper(dbVerbose, dbDryRun)
	status := helper.GetAllStatus()

	format, err := output.ParseFormat(dbOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Database Services Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for _, svc := range status.Services {
		var icon string
		switch {
		case !svc.Installed:
			icon = output.Warning("○")
		case svc.Running:
			icon = output.Success("●")
		default:
			icon = output.Error("○")
		}
		fmt.Fprintf(os.Stdout, "%s %-12s %s\n", icon, svc.Name+":", svc.Status)
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "Running: %d, Stopped: %d, Not installed: %d\n",
		status.Running, status.Stopped, status.Missing)

	return nil
}

func runDbStart(cmd *cobra.Command, args []string) error {
	helper := database.NewHelper(dbVerbose, dbDryRun)

	if err := helper.StartService(args[0]); err != nil {
		return err
	}

	if !dbDryRun {
		fmt.Fprintf(os.Stdout, "%s Started %s\n", output.Success("✓"), args[0])
	}
	return nil
}

func runDbStop(cmd *cobra.Command, args []string) error {
	helper := database.NewHelper(dbVerbose, dbDryRun)

	if err := helper.StopService(args[0]); err != nil {
		return err
	}

	if !dbDryRun {
		fmt.Fprintf(os.Stdout, "%s Stopped %s\n", output.Success("✓"), args[0])
	}
	return nil
}

func runDbRestart(cmd *cobra.Command, args []string) error {
	helper := database.NewHelper(dbVerbose, dbDryRun)

	if err := helper.RestartService(args[0]); err != nil {
		return err
	}

	if !dbDryRun {
		fmt.Fprintf(os.Stdout, "%s Restarted %s\n", output.Success("✓"), args[0])
	}
	return nil
}

func runDbStartAll(cmd *cobra.Command, args []string) error {
	helper := database.NewHelper(dbVerbose, dbDryRun)

	fmt.Fprintln(os.Stdout, "Starting database services...")
	if err := helper.StartAll(); err != nil {
		return err
	}

	if !dbDryRun {
		fmt.Fprintf(os.Stdout, "%s Started common database services\n", output.Success("✓"))
		fmt.Fprintln(os.Stdout, "Use 'acorn db status' to check status.")
	}
	return nil
}

func runDbStopAll(cmd *cobra.Command, args []string) error {
	helper := database.NewHelper(dbVerbose, dbDryRun)

	fmt.Fprintln(os.Stdout, "Stopping database services...")
	if err := helper.StopAll(); err != nil {
		return err
	}

	if !dbDryRun {
		fmt.Fprintf(os.Stdout, "%s Stopped common database services\n", output.Success("✓"))
	}
	return nil
}

func runDbList(cmd *cobra.Command, args []string) error {
	helper := database.NewHelper(dbVerbose, dbDryRun)
	services := helper.GetSupportedServices()

	format, err := output.ParseFormat(dbOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(map[string][]string{"services": services})
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Supported Database Services"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for _, svc := range services {
		fmt.Fprintf(os.Stdout, "  • %s\n", svc)
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "Service management requires Homebrew (macOS).")

	return nil
}
