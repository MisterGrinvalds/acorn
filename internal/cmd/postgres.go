package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/data/postgres"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/mistergrinvalds/acorn/internal/utils/configcmd"
	"github.com/spf13/cobra"
)

var (
	pgDryRun     bool
	pgVerbose    bool
	pgHost       string
	pgPort       string
	pgUser       string
	pgUseDocker  bool
	pgPassword   string
	pgOutputPath string
)

// postgresCmd represents the postgres command group
var postgresCmd = &cobra.Command{
	Use:   "postgres",
	Short: "PostgreSQL database management",
	Long: `PostgreSQL database management commands.

Manage PostgreSQL servers, databases, and connections.

Examples:
  acorn data postgres status      # Show status
  acorn data postgres start       # Start server
  acorn data postgres databases   # List databases
  acorn data postgres connect     # Connect to database`,
	Aliases: []string{"pg", "psql"},
}

// pgStatusCmd shows status
var pgStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show PostgreSQL status",
	Long: `Show PostgreSQL installation and server status.

Examples:
  acorn data postgres status
  acorn data postgres status -o json`,
	RunE: runPgStatus,
}

// pgStartCmd starts server
var pgStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start PostgreSQL server",
	Long: `Start PostgreSQL server locally or in Docker.

Examples:
  acorn data postgres start
  acorn data postgres start --docker
  acorn data postgres start --docker --port 5433`,
	RunE: runPgStart,
}

// pgStopCmd stops server
var pgStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop PostgreSQL server",
	Long: `Stop PostgreSQL server.

Examples:
  acorn data postgres stop`,
	RunE: runPgStop,
}

// pgConnectCmd connects to database
var pgConnectCmd = &cobra.Command{
	Use:   "connect [database]",
	Short: "Connect to database",
	Long: `Connect to a PostgreSQL database using psql.

Examples:
  acorn data postgres connect
  acorn data postgres connect mydb
  acorn data postgres connect mydb -h localhost -p 5432 -U postgres`,
	Args: cobra.MaximumNArgs(1),
	RunE: runPgConnect,
}

// pgDatabasesCmd lists databases
var pgDatabasesCmd = &cobra.Command{
	Use:   "databases",
	Short: "List databases",
	Long: `List all PostgreSQL databases.

Examples:
  acorn data postgres databases
  acorn data postgres databases -o json`,
	Aliases: []string{"dbs", "ls"},
	RunE:    runPgDatabases,
}

// pgCreateCmd creates a database
var pgCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create database",
	Long: `Create a new PostgreSQL database.

Examples:
  acorn data postgres create mydb`,
	Args: cobra.ExactArgs(1),
	RunE: runPgCreate,
}

// pgDropCmd drops a database
var pgDropCmd = &cobra.Command{
	Use:   "drop <name>",
	Short: "Drop database",
	Long: `Drop a PostgreSQL database.

Examples:
  acorn data postgres drop mydb`,
	Args: cobra.ExactArgs(1),
	RunE: runPgDrop,
}

// pgDumpCmd dumps a database
var pgDumpCmd = &cobra.Command{
	Use:   "dump <database>",
	Short: "Dump database",
	Long: `Dump a PostgreSQL database to a SQL file.

Examples:
  acorn data postgres dump mydb
  acorn data postgres dump mydb --output backup.sql`,
	Args: cobra.ExactArgs(1),
	RunE: runPgDump,
}

// pgRestoreCmd restores a database
var pgRestoreCmd = &cobra.Command{
	Use:   "restore <database> <file>",
	Short: "Restore database",
	Long: `Restore a PostgreSQL database from a SQL file.

Examples:
  acorn data postgres restore mydb backup.sql`,
	Args: cobra.ExactArgs(2),
	RunE: runPgRestore,
}

// pgInstallCmd installs PostgreSQL
var pgInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install PostgreSQL",
	Long: `Install PostgreSQL using Homebrew.

Examples:
  acorn data postgres install`,
	RunE: runPgInstall,
}

func init() {
	dataCmd.AddCommand(postgresCmd)

	// Add subcommands
	postgresCmd.AddCommand(pgStatusCmd)
	postgresCmd.AddCommand(pgStartCmd)
	postgresCmd.AddCommand(pgStopCmd)
	postgresCmd.AddCommand(pgConnectCmd)
	postgresCmd.AddCommand(pgDatabasesCmd)
	postgresCmd.AddCommand(pgCreateCmd)
	postgresCmd.AddCommand(pgDropCmd)
	postgresCmd.AddCommand(pgDumpCmd)
	postgresCmd.AddCommand(pgRestoreCmd)
	postgresCmd.AddCommand(pgInstallCmd)
	postgresCmd.AddCommand(configcmd.NewConfigRouter("postgres"))

	// Start flags
	pgStartCmd.Flags().BoolVar(&pgUseDocker, "docker", false, "Use Docker instead of local")
	pgStartCmd.Flags().StringVar(&pgPassword, "password", "postgres", "Password for postgres user")
	pgStartCmd.Flags().StringVar(&pgPort, "port", "5432", "Port to run on")

	// Dump flags
	pgDumpCmd.Flags().StringVarP(&pgOutputPath, "output", "O", "", "Output file path")

	// Connection flags (persistent)
	postgresCmd.PersistentFlags().StringVarP(&pgHost, "host", "H", "localhost", "Database host")
	postgresCmd.PersistentFlags().StringVarP(&pgPort, "port", "p", "5432", "Database port")
	postgresCmd.PersistentFlags().StringVarP(&pgUser, "user", "U", "postgres", "Database user")

	// Output flags
	postgresCmd.PersistentFlags().BoolVar(&pgDryRun, "dry-run", false,
		"Show what would be done without executing")
	postgresCmd.PersistentFlags().BoolVarP(&pgVerbose, "verbose", "v", false,
		"Show verbose output")
}

func newPgHelper() *postgres.Helper {
	return postgres.NewHelper(pgVerbose, pgDryRun)
}

func runPgStatus(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := newPgHelper()
	status := helper.GetStatus()

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("PostgreSQL Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if status.Installed {
		fmt.Fprintf(os.Stdout, "%s PostgreSQL client: %s\n", output.Success("✓"), status.PsqlVersion)
		if status.Version != "" {
			fmt.Fprintf(os.Stdout, "%s PostgreSQL server: %s\n", output.Success("✓"), status.Version)
		}
	} else {
		fmt.Fprintf(os.Stdout, "%s PostgreSQL not installed\n", output.Error("✗"))
		fmt.Fprintln(os.Stdout, "  Install: acorn data postgres install")
		return nil
	}

	fmt.Fprintln(os.Stdout)

	if status.ServerRunning {
		fmt.Fprintf(os.Stdout, "%s Server running (port %s)\n", output.Success("✓"), status.Port)
	} else if status.DockerRunning {
		fmt.Fprintf(os.Stdout, "%s Docker container running (port %s)\n", output.Success("✓"), status.Port)
		fmt.Fprintf(os.Stdout, "  Container: %s\n", status.ContainerID[:12])
	} else {
		fmt.Fprintf(os.Stdout, "%s Server not running\n", output.Warning("!"))
		fmt.Fprintln(os.Stdout, "  Start: acorn data postgres start")
	}

	if status.DataDir != "" {
		fmt.Fprintf(os.Stdout, "\nData dir: %s\n", status.DataDir)
	}

	return nil
}

func runPgStart(cmd *cobra.Command, args []string) error {
	helper := newPgHelper()

	port := 5432
	if pgPort != "" {
		fmt.Sscanf(pgPort, "%d", &port)
	}

	if err := helper.Start(pgUseDocker, port, pgPassword); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s PostgreSQL started\n", output.Success("✓"))
	return nil
}

func runPgStop(cmd *cobra.Command, args []string) error {
	helper := newPgHelper()
	if err := helper.Stop(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s PostgreSQL stopped\n", output.Success("✓"))
	return nil
}

func runPgConnect(cmd *cobra.Command, args []string) error {
	helper := newPgHelper()

	database := "postgres"
	if len(args) > 0 {
		database = args[0]
	}

	return helper.Connect(database, pgHost, pgPort, pgUser)
}

func runPgDatabases(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := newPgHelper()
	databases, err := helper.ListDatabases(pgHost, pgPort, pgUser)
	if err != nil {
		return err
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(databases)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Databases"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if len(databases) == 0 {
		fmt.Fprintln(os.Stdout, "No databases found")
		return nil
	}

	for _, db := range databases {
		fmt.Fprintf(os.Stdout, "  %-20s %-15s %s\n", db.Name, db.Owner, db.Size)
	}

	return nil
}

func runPgCreate(cmd *cobra.Command, args []string) error {
	helper := newPgHelper()
	if err := helper.CreateDatabase(args[0], pgHost, pgPort, pgUser); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Database created: %s\n", output.Success("✓"), args[0])
	return nil
}

func runPgDrop(cmd *cobra.Command, args []string) error {
	helper := newPgHelper()
	if err := helper.DropDatabase(args[0], pgHost, pgPort, pgUser); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Database dropped: %s\n", output.Success("✓"), args[0])
	return nil
}

func runPgDump(cmd *cobra.Command, args []string) error {
	helper := newPgHelper()
	outputPath := pgOutputPath
	if outputPath == "" {
		outputPath = args[0] + ".sql"
	}
	if err := helper.Dump(args[0], outputPath, pgHost, pgPort, pgUser); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Database dumped to: %s\n", output.Success("✓"), outputPath)
	return nil
}

func runPgRestore(cmd *cobra.Command, args []string) error {
	helper := newPgHelper()
	if err := helper.Restore(args[0], args[1], pgHost, pgPort, pgUser); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Database restored: %s\n", output.Success("✓"), args[0])
	return nil
}

func runPgInstall(cmd *cobra.Command, args []string) error {
	helper := newPgHelper()
	if err := helper.Install(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s PostgreSQL installed\n", output.Success("✓"))
	return nil
}
