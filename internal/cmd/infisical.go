package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/data/infisical"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	infisicalDryRun  bool
	infisicalVerbose bool
	infisicalEnv     string
	infisicalFormat  string
)

// infisicalCmd represents the infisical command group
var infisicalCmd = &cobra.Command{
	Use:   "infisical",
	Short: "Infisical secret management",
	Long: `Infisical secret management CLI commands.

Infisical is an open-source secret management platform.
Provides commands for authentication, secrets, and environment injection.

Examples:
  acorn data infisical status      # Show status
  acorn data infisical login       # Authenticate
  acorn data infisical secrets     # List secrets
  acorn data infisical run -- cmd  # Run with secrets`,
	Aliases: []string{"inf", "secrets"},
}

// infisicalStatusCmd shows status
var infisicalStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Infisical status",
	Long: `Show Infisical CLI installation and authentication status.

Examples:
  acorn data infisical status
  acorn data infisical status -o json`,
	RunE: runInfisicalStatus,
}

// infisicalInstallCmd installs Infisical
var infisicalInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Infisical CLI",
	Long: `Install Infisical CLI using Homebrew.

Examples:
  acorn data infisical install`,
	RunE: runInfisicalInstall,
}

// infisicalLoginCmd authenticates
var infisicalLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Infisical",
	Long: `Log in to Infisical.

Examples:
  acorn data infisical login`,
	RunE: runInfisicalLogin,
}

// infisicalLogoutCmd logs out
var infisicalLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out from Infisical",
	Long: `Log out from Infisical.

Examples:
  acorn data infisical logout`,
	RunE: runInfisicalLogout,
}

// infisicalInitCmd initializes project
var infisicalInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Infisical project",
	Long: `Initialize Infisical in the current directory.

Creates .infisical.json with project configuration.

Examples:
  acorn data infisical init`,
	RunE: runInfisicalInit,
}

// infisicalSecretsCmd manages secrets
var infisicalSecretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "List secrets",
	Long: `List secrets for the current project/environment.

Examples:
  acorn data infisical secrets
  acorn data infisical secrets --env production`,
	Aliases: []string{"ls", "list"},
	RunE:    runInfisicalSecrets,
}

// infisicalGetCmd gets a specific secret
var infisicalGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a secret value",
	Long: `Get the value of a specific secret.

Examples:
  acorn data infisical get DATABASE_URL
  acorn data infisical get API_KEY --env production`,
	Args: cobra.ExactArgs(1),
	RunE: runInfisicalGet,
}

// infisicalRunCmd runs command with secrets
var infisicalRunCmd = &cobra.Command{
	Use:   "run -- <command>",
	Short: "Run command with secrets",
	Long: `Execute a command with secrets injected as environment variables.

Examples:
  acorn data infisical run -- npm start
  acorn data infisical run --env production -- ./my-app`,
	Args: cobra.MinimumNArgs(1),
	RunE: runInfisicalRun,
}

// infisicalExportCmd exports secrets
var infisicalExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export secrets",
	Long: `Export secrets to various formats.

Supported formats: dotenv, json, yaml, csv

Examples:
  acorn data infisical export
  acorn data infisical export --format json
  acorn data infisical export --env production --format dotenv`,
	RunE: runInfisicalExport,
}

// infisicalScanCmd scans for secrets
var infisicalScanCmd = &cobra.Command{
	Use:   "scan [path]",
	Short: "Scan for exposed secrets",
	Long: `Scan files for accidentally committed secrets.

Examples:
  acorn data infisical scan
  acorn data infisical scan ./src`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInfisicalScan,
}

func init() {
	dataCmd.AddCommand(infisicalCmd)

	// Add subcommands
	infisicalCmd.AddCommand(infisicalStatusCmd)
	infisicalCmd.AddCommand(infisicalInstallCmd)
	infisicalCmd.AddCommand(infisicalLoginCmd)
	infisicalCmd.AddCommand(infisicalLogoutCmd)
	infisicalCmd.AddCommand(infisicalInitCmd)
	infisicalCmd.AddCommand(infisicalSecretsCmd)
	infisicalCmd.AddCommand(infisicalGetCmd)
	infisicalCmd.AddCommand(infisicalRunCmd)
	infisicalCmd.AddCommand(infisicalExportCmd)
	infisicalCmd.AddCommand(infisicalScanCmd)

	// Environment flag for multiple commands
	infisicalSecretsCmd.Flags().StringVarP(&infisicalEnv, "env", "e", "", "Environment (dev, staging, production)")
	infisicalGetCmd.Flags().StringVarP(&infisicalEnv, "env", "e", "", "Environment")
	infisicalRunCmd.Flags().StringVarP(&infisicalEnv, "env", "e", "", "Environment")
	infisicalExportCmd.Flags().StringVarP(&infisicalEnv, "env", "e", "", "Environment")
	infisicalExportCmd.Flags().StringVar(&infisicalFormat, "format", "dotenv", "Export format (dotenv|json|yaml|csv)")

	// Persistent flags
	infisicalCmd.PersistentFlags().BoolVar(&infisicalDryRun, "dry-run", false,
		"Show what would be done without executing")
	infisicalCmd.PersistentFlags().BoolVarP(&infisicalVerbose, "verbose", "v", false,
		"Show verbose output")
}

func newInfisicalHelper() *infisical.Helper {
	return infisical.NewHelper(infisicalVerbose, infisicalDryRun)
}

func runInfisicalStatus(cmd *cobra.Command, args []string) error {
	helper := newInfisicalHelper()
	status := helper.GetStatus()

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Infisical Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if status.Installed {
		fmt.Fprintf(os.Stdout, "%s Infisical CLI installed: %s\n", output.Success("✓"), status.Version)
	} else {
		fmt.Fprintf(os.Stdout, "%s Infisical CLI not found\n", output.Error("✗"))
		fmt.Fprintln(os.Stdout, "  Install: brew install infisical/get-cli/infisical")
		return nil
	}

	fmt.Fprintln(os.Stdout)

	if status.LoggedIn {
		fmt.Fprintf(os.Stdout, "%s Authenticated\n", output.Success("✓"))
	} else {
		fmt.Fprintf(os.Stdout, "%s Not logged in (run 'acorn data infisical login')\n", output.Warning("!"))
	}

	if status.InProject {
		fmt.Fprintf(os.Stdout, "%s In project\n", output.Success("✓"))
		if status.ProjectID != "" {
			fmt.Fprintf(os.Stdout, "  Project: %s\n", status.ProjectID)
		}
		if status.Environment != "" {
			fmt.Fprintf(os.Stdout, "  Default env: %s\n", status.Environment)
		}
	} else {
		fmt.Fprintf(os.Stdout, "%s No project (run 'acorn data infisical init')\n", output.Warning("!"))
	}

	return nil
}

func runInfisicalInstall(cmd *cobra.Command, args []string) error {
	helper := newInfisicalHelper()
	return helper.Install()
}

func runInfisicalLogin(cmd *cobra.Command, args []string) error {
	helper := newInfisicalHelper()
	return helper.Login()
}

func runInfisicalLogout(cmd *cobra.Command, args []string) error {
	helper := newInfisicalHelper()
	if err := helper.Logout(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Logged out\n", output.Success("✓"))
	return nil
}

func runInfisicalInit(cmd *cobra.Command, args []string) error {
	helper := newInfisicalHelper()
	return helper.Init()
}

func runInfisicalSecrets(cmd *cobra.Command, args []string) error {
	helper := newInfisicalHelper()
	secrets, err := helper.ListSecrets(infisicalEnv)
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(secrets)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Secrets"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if len(secrets) == 0 {
		fmt.Fprintln(os.Stdout, "No secrets found")
		return nil
	}

	for _, s := range secrets {
		// Mask the value for display
		maskedValue := "****"
		if len(s.Value) > 4 {
			maskedValue = s.Value[:2] + "****" + s.Value[len(s.Value)-2:]
		}
		fmt.Fprintf(os.Stdout, "  %s = %s\n", s.Key, maskedValue)
	}

	fmt.Fprintf(os.Stdout, "\n%d secrets\n", len(secrets))
	return nil
}

func runInfisicalGet(cmd *cobra.Command, args []string) error {
	helper := newInfisicalHelper()
	key := args[0]

	value, err := helper.GetSecret(key, infisicalEnv)
	if err != nil {
		return err
	}

	fmt.Println(value)
	return nil
}

func runInfisicalRun(cmd *cobra.Command, args []string) error {
	helper := newInfisicalHelper()
	return helper.Run(infisicalEnv, args)
}

func runInfisicalExport(cmd *cobra.Command, args []string) error {
	helper := newInfisicalHelper()
	exported, err := helper.Export(infisicalEnv, infisicalFormat)
	if err != nil {
		return err
	}

	fmt.Print(exported)
	return nil
}

func runInfisicalScan(cmd *cobra.Command, args []string) error {
	helper := newInfisicalHelper()
	path := ""
	if len(args) > 0 {
		path = args[0]
	}
	return helper.Scan(path)
}
