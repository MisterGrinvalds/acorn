package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/identity/keycloak"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	keycloakOutputFormat  string
	keycloakDryRun        bool
	keycloakVerbose       bool
	keycloakPort          int
	keycloakDevMode       bool
	keycloakAdminUser     string
	keycloakAdminPassword string
	keycloakFollow        bool
	keycloakOutputPath    string
)

// keycloakCmd represents the keycloak command group
var keycloakCmd = &cobra.Command{
	Use:   "keycloak",
	Short: "Keycloak IAM platform",
	Long: `Keycloak Identity and Access Management commands.

Manage Keycloak instances, realms, and authentication configuration.

Examples:
  acorn identity keycloak status      # Show status
  acorn identity keycloak start       # Start Keycloak in Docker
  acorn identity keycloak stop        # Stop Keycloak
  acorn identity keycloak realms      # List exported realms`,
	Aliases: []string{"kc"},
}

// keycloakStatusCmd shows status
var keycloakStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Keycloak status",
	Long: `Show Keycloak installation and running status.

Examples:
  acorn identity keycloak status
  acorn identity keycloak status -o json`,
	RunE: runKeycloakStatus,
}

// keycloakStartCmd starts Keycloak
var keycloakStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Keycloak in Docker",
	Long: `Start Keycloak in a Docker container.

Examples:
  acorn identity keycloak start
  acorn identity keycloak start --port 8081
  acorn identity keycloak start --dev`,
	RunE: runKeycloakStart,
}

// keycloakStopCmd stops Keycloak
var keycloakStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Keycloak container",
	Long: `Stop the running Keycloak container.

Examples:
  acorn identity keycloak stop`,
	RunE: runKeycloakStop,
}

// keycloakLogsCmd shows logs
var keycloakLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Show Keycloak logs",
	Long: `Show logs from the Keycloak container.

Examples:
  acorn identity keycloak logs
  acorn identity keycloak logs -f`,
	RunE: runKeycloakLogs,
}

// keycloakOpenCmd opens admin console
var keycloakOpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Open admin console",
	Long: `Open Keycloak admin console in browser.

Examples:
  acorn identity keycloak open`,
	RunE: runKeycloakOpen,
}

// keycloakRealmsCmd lists realms
var keycloakRealmsCmd = &cobra.Command{
	Use:   "realms",
	Short: "List exported realms",
	Long: `List realm export files.

Examples:
  acorn identity keycloak realms
  acorn identity keycloak realms -o json`,
	Aliases: []string{"ls"},
	RunE:    runKeycloakRealms,
}

// keycloakExportCmd exports a realm
var keycloakExportCmd = &cobra.Command{
	Use:   "export <realm>",
	Short: "Export a realm",
	Long: `Export a realm configuration to JSON file.

Examples:
  acorn identity keycloak export master
  acorn identity keycloak export myrealm --output ./myrealm.json`,
	Args: cobra.ExactArgs(1),
	RunE: runKeycloakExport,
}

// keycloakImportCmd imports a realm
var keycloakImportCmd = &cobra.Command{
	Use:   "import <file>",
	Short: "Import a realm",
	Long: `Import a realm configuration from JSON file.

Examples:
  acorn identity keycloak import ./myrealm.json`,
	Args: cobra.ExactArgs(1),
	RunE: runKeycloakImport,
}

func init() {
	identityCmd.AddCommand(keycloakCmd)

	// Add subcommands
	keycloakCmd.AddCommand(keycloakStatusCmd)
	keycloakCmd.AddCommand(keycloakStartCmd)
	keycloakCmd.AddCommand(keycloakStopCmd)
	keycloakCmd.AddCommand(keycloakLogsCmd)
	keycloakCmd.AddCommand(keycloakOpenCmd)
	keycloakCmd.AddCommand(keycloakRealmsCmd)
	keycloakCmd.AddCommand(keycloakExportCmd)
	keycloakCmd.AddCommand(keycloakImportCmd)

	// Start flags
	keycloakStartCmd.Flags().IntVar(&keycloakPort, "port", 8080, "Port to expose")
	keycloakStartCmd.Flags().BoolVar(&keycloakDevMode, "dev", true, "Start in development mode")
	keycloakStartCmd.Flags().StringVar(&keycloakAdminUser, "admin-user", "admin", "Admin username")
	keycloakStartCmd.Flags().StringVar(&keycloakAdminPassword, "admin-password", "admin", "Admin password")

	// Logs flags
	keycloakLogsCmd.Flags().BoolVarP(&keycloakFollow, "follow", "f", false, "Follow log output")

	// Export flags
	keycloakExportCmd.Flags().StringVarP(&keycloakOutputPath, "output", "O", "", "Output file path")

	// Persistent flags
	keycloakCmd.PersistentFlags().StringVarP(&keycloakOutputFormat, "output", "o", "table",
		"Output format (table|json|yaml)")
	keycloakCmd.PersistentFlags().BoolVar(&keycloakDryRun, "dry-run", false,
		"Show what would be done without executing")
	keycloakCmd.PersistentFlags().BoolVarP(&keycloakVerbose, "verbose", "v", false,
		"Show verbose output")
}

func newKeycloakHelper() *keycloak.Helper {
	return keycloak.NewHelper(keycloakVerbose, keycloakDryRun)
}

func runKeycloakStatus(cmd *cobra.Command, args []string) error {
	helper := newKeycloakHelper()
	status := helper.GetStatus()

	format, err := output.ParseFormat(keycloakOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Keycloak Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if status.DockerInstalled {
		fmt.Fprintf(os.Stdout, "%s Docker available\n", output.Success("✓"))
	} else {
		fmt.Fprintf(os.Stdout, "%s Docker not found\n", output.Error("✗"))
		fmt.Fprintln(os.Stdout, "  Keycloak requires Docker to run locally")
		return nil
	}

	fmt.Fprintln(os.Stdout)

	if status.ContainerRunning {
		fmt.Fprintf(os.Stdout, "%s Keycloak running\n", output.Success("✓"))
		fmt.Fprintf(os.Stdout, "  Container: %s (%s)\n", status.ContainerName, status.ContainerID[:12])
		fmt.Fprintf(os.Stdout, "  URL: http://localhost:%s\n", status.Port)
		fmt.Fprintf(os.Stdout, "  Admin: http://localhost:%s/admin\n", status.Port)
	} else {
		fmt.Fprintf(os.Stdout, "%s Keycloak not running\n", output.Info("ℹ"))
		fmt.Fprintln(os.Stdout, "  Start: acorn identity keycloak start")
	}

	if status.AdminCLI {
		fmt.Fprintf(os.Stdout, "\n%s Admin CLI: %s\n", output.Success("✓"), status.AdminCLIPath)
	}

	fmt.Fprintf(os.Stdout, "\nRealms dir: %s\n", status.RealmsDir)

	return nil
}

func runKeycloakStart(cmd *cobra.Command, args []string) error {
	helper := newKeycloakHelper()
	return helper.Start(keycloakPort, keycloakDevMode, keycloakAdminUser, keycloakAdminPassword)
}

func runKeycloakStop(cmd *cobra.Command, args []string) error {
	helper := newKeycloakHelper()
	if err := helper.Stop(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Keycloak stopped\n", output.Success("✓"))
	return nil
}

func runKeycloakLogs(cmd *cobra.Command, args []string) error {
	helper := newKeycloakHelper()
	return helper.Logs(keycloakFollow)
}

func runKeycloakOpen(cmd *cobra.Command, args []string) error {
	helper := newKeycloakHelper()
	if err := helper.Open(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Opening admin console\n", output.Success("✓"))
	return nil
}

func runKeycloakRealms(cmd *cobra.Command, args []string) error {
	helper := newKeycloakHelper()
	realms, err := helper.ListRealms()
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(keycloakOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(realms)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Exported Realms"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if len(realms) == 0 {
		fmt.Fprintln(os.Stdout, "No exported realms found")
		fmt.Fprintln(os.Stdout, "  Export: acorn identity keycloak export <realm>")
		return nil
	}

	for _, r := range realms {
		status := output.Info("○")
		if r.Enabled {
			status = output.Success("●")
		}
		fmt.Fprintf(os.Stdout, "  %s %s\n", status, r.Name)
		fmt.Fprintf(os.Stdout, "    %s\n", r.Path)
	}

	return nil
}

func runKeycloakExport(cmd *cobra.Command, args []string) error {
	helper := newKeycloakHelper()
	if err := helper.ExportRealm(args[0], keycloakOutputPath); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Realm exported\n", output.Success("✓"))
	return nil
}

func runKeycloakImport(cmd *cobra.Command, args []string) error {
	helper := newKeycloakHelper()
	if err := helper.ImportRealm(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Realm imported\n", output.Success("✓"))
	return nil
}
