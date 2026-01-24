package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/artifacts/jfrog"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	jfrogOutputFormat string
	jfrogDryRun       bool
	jfrogVerbose      bool
	jfrogRecursive    bool
	jfrogFlat         bool
	jfrogQuiet        bool
	jfrogServerID     string
)

// jfrogCmd represents the jfrog command group
var jfrogCmd = &cobra.Command{
	Use:   "jfrog",
	Short: "JFrog Artifactory CLI",
	Long: `JFrog Artifactory CLI commands.

JFrog CLI for managing artifacts, Docker images, and builds.
Provides commands for server configuration, upload/download, and search.

Examples:
  acorn artifacts jfrog status       # Show status
  acorn artifacts jfrog servers      # List servers
  acorn artifacts jfrog upload       # Upload artifacts
  acorn artifacts jfrog download     # Download artifacts`,
	Aliases: []string{"jf", "artifactory"},
}

// jfrogStatusCmd shows status
var jfrogStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show JFrog CLI status",
	Long: `Show JFrog CLI installation and server configuration status.

Examples:
  acorn artifacts jfrog status
  acorn artifacts jfrog status -o json`,
	RunE: runJfrogStatus,
}

// jfrogInstallCmd installs JFrog CLI
var jfrogInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install JFrog CLI",
	Long: `Install JFrog CLI using Homebrew.

Examples:
  acorn artifacts jfrog install`,
	RunE: runJfrogInstall,
}

// jfrogPingCmd tests connectivity
var jfrogPingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Test server connectivity",
	Long: `Test connectivity to JFrog Artifactory server.

Examples:
  acorn artifacts jfrog ping
  acorn artifacts jfrog ping --server-id myserver`,
	RunE: runJfrogPing,
}

// Server management commands
var jfrogServersCmd = &cobra.Command{
	Use:   "servers",
	Short: "Server configuration",
	Long: `Manage JFrog server configurations.

Examples:
  acorn artifacts jfrog servers list
  acorn artifacts jfrog servers add`,
}

var jfrogServersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured servers",
	Long: `List all configured JFrog servers.

Examples:
  acorn artifacts jfrog servers list`,
	Aliases: []string{"ls"},
	RunE:    runJfrogServersList,
}

var jfrogServersAddCmd = &cobra.Command{
	Use:   "add <server-id>",
	Short: "Add server configuration",
	Long: `Add a new JFrog server configuration.

Runs interactive setup to configure server URL and credentials.

Examples:
  acorn artifacts jfrog servers add myserver`,
	Args: cobra.ExactArgs(1),
	RunE: runJfrogServersAdd,
}

var jfrogServersRemoveCmd = &cobra.Command{
	Use:   "remove <server-id>",
	Short: "Remove server configuration",
	Long: `Remove a JFrog server configuration.

Examples:
  acorn artifacts jfrog servers remove myserver`,
	Aliases: []string{"rm", "delete"},
	Args:    cobra.ExactArgs(1),
	RunE:    runJfrogServersRemove,
}

var jfrogServersUseCmd = &cobra.Command{
	Use:   "use <server-id>",
	Short: "Set default server",
	Long: `Set the default JFrog server.

Examples:
  acorn artifacts jfrog servers use myserver`,
	Args: cobra.ExactArgs(1),
	RunE: runJfrogServersUse,
}

// Artifact commands
var jfrogUploadCmd = &cobra.Command{
	Use:   "upload <source> <target>",
	Short: "Upload artifacts",
	Long: `Upload artifacts to Artifactory.

Examples:
  acorn artifacts jfrog upload "*.jar" libs-release-local/
  acorn artifacts jfrog upload ./dist/ builds/myapp/ --recursive`,
	Args: cobra.ExactArgs(2),
	RunE: runJfrogUpload,
}

var jfrogDownloadCmd = &cobra.Command{
	Use:   "download <source> [target]",
	Short: "Download artifacts",
	Long: `Download artifacts from Artifactory.

Examples:
  acorn artifacts jfrog download libs-release-local/mylib.jar
  acorn artifacts jfrog download "libs-release-local/*.jar" ./libs/`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runJfrogDownload,
}

var jfrogSearchCmd = &cobra.Command{
	Use:   "search <pattern>",
	Short: "Search artifacts",
	Long: `Search for artifacts in Artifactory.

Examples:
  acorn artifacts jfrog search "libs-release-local/*.jar"
  acorn artifacts jfrog search "*/myapp-*.tar.gz"`,
	Args: cobra.ExactArgs(1),
	RunE: runJfrogSearch,
}

var jfrogDeleteCmd = &cobra.Command{
	Use:   "delete <pattern>",
	Short: "Delete artifacts",
	Long: `Delete artifacts from Artifactory.

Examples:
  acorn artifacts jfrog delete "libs-snapshot-local/old/*"`,
	Args: cobra.ExactArgs(1),
	RunE: runJfrogDelete,
}

// Docker commands
var jfrogDockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Docker operations",
	Long: `Docker image operations with Artifactory.

Examples:
  acorn artifacts jfrog docker push myimage:latest docker-local
  acorn artifacts jfrog docker pull myimage:latest docker-local`,
}

var jfrogDockerPushCmd = &cobra.Command{
	Use:   "push <image> <repo>",
	Short: "Push Docker image",
	Long: `Push a Docker image to Artifactory.

Examples:
  acorn artifacts jfrog docker push myapp:latest docker-local`,
	Args: cobra.ExactArgs(2),
	RunE: runJfrogDockerPush,
}

var jfrogDockerPullCmd = &cobra.Command{
	Use:   "pull <image> <repo>",
	Short: "Pull Docker image",
	Long: `Pull a Docker image from Artifactory.

Examples:
  acorn artifacts jfrog docker pull myapp:latest docker-local`,
	Args: cobra.ExactArgs(2),
	RunE: runJfrogDockerPull,
}

func init() {
	artifactsCmd.AddCommand(jfrogCmd)

	// Add subcommands
	jfrogCmd.AddCommand(jfrogStatusCmd)
	jfrogCmd.AddCommand(jfrogInstallCmd)
	jfrogCmd.AddCommand(jfrogPingCmd)

	// Server subcommands
	jfrogCmd.AddCommand(jfrogServersCmd)
	jfrogServersCmd.AddCommand(jfrogServersListCmd)
	jfrogServersCmd.AddCommand(jfrogServersAddCmd)
	jfrogServersCmd.AddCommand(jfrogServersRemoveCmd)
	jfrogServersCmd.AddCommand(jfrogServersUseCmd)

	// Artifact subcommands
	jfrogCmd.AddCommand(jfrogUploadCmd)
	jfrogCmd.AddCommand(jfrogDownloadCmd)
	jfrogCmd.AddCommand(jfrogSearchCmd)
	jfrogCmd.AddCommand(jfrogDeleteCmd)

	// Docker subcommands
	jfrogCmd.AddCommand(jfrogDockerCmd)
	jfrogDockerCmd.AddCommand(jfrogDockerPushCmd)
	jfrogDockerCmd.AddCommand(jfrogDockerPullCmd)

	// Flags
	jfrogPingCmd.Flags().StringVar(&jfrogServerID, "server-id", "", "Server ID to ping")
	jfrogUploadCmd.Flags().BoolVar(&jfrogRecursive, "recursive", false, "Upload recursively")
	jfrogUploadCmd.Flags().BoolVar(&jfrogFlat, "flat", false, "Flatten directory structure")
	jfrogDownloadCmd.Flags().BoolVar(&jfrogRecursive, "recursive", false, "Download recursively")
	jfrogDownloadCmd.Flags().BoolVar(&jfrogFlat, "flat", false, "Flatten directory structure")
	jfrogDeleteCmd.Flags().BoolVar(&jfrogQuiet, "quiet", false, "Skip confirmation")

	// Persistent flags
	jfrogCmd.PersistentFlags().StringVarP(&jfrogOutputFormat, "output", "o", "table",
		"Output format (table|json|yaml)")
	jfrogCmd.PersistentFlags().BoolVar(&jfrogDryRun, "dry-run", false,
		"Show what would be done without executing")
	jfrogCmd.PersistentFlags().BoolVarP(&jfrogVerbose, "verbose", "v", false,
		"Show verbose output")
}

func newJfrogHelper() *jfrog.Helper {
	return jfrog.NewHelper(jfrogVerbose, jfrogDryRun)
}

func runJfrogStatus(cmd *cobra.Command, args []string) error {
	helper := newJfrogHelper()
	status := helper.GetStatus()

	format, err := output.ParseFormat(jfrogOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("JFrog CLI Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if status.Installed {
		fmt.Fprintf(os.Stdout, "%s JFrog CLI installed: %s\n", output.Success("✓"), status.Version)
	} else {
		fmt.Fprintf(os.Stdout, "%s JFrog CLI not found\n", output.Error("✗"))
		fmt.Fprintln(os.Stdout, "  Install: brew install jfrog-cli")
		return nil
	}

	fmt.Fprintln(os.Stdout)

	if len(status.Servers) == 0 {
		fmt.Fprintf(os.Stdout, "%s No servers configured\n", output.Warning("!"))
		fmt.Fprintln(os.Stdout, "  Add server: acorn artifacts jfrog servers add <server-id>")
	} else {
		fmt.Fprintf(os.Stdout, "Servers: %d configured\n", len(status.Servers))
		for _, s := range status.Servers {
			marker := "  "
			if s.IsDefault {
				marker = output.Success("→ ")
			}
			fmt.Fprintf(os.Stdout, "%s%s: %s\n", marker, s.ID, s.URL)
		}
	}

	return nil
}

func runJfrogInstall(cmd *cobra.Command, args []string) error {
	helper := newJfrogHelper()
	return helper.Install()
}

func runJfrogPing(cmd *cobra.Command, args []string) error {
	helper := newJfrogHelper()
	if err := helper.Ping(jfrogServerID); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Server is reachable\n", output.Success("✓"))
	return nil
}

func runJfrogServersList(cmd *cobra.Command, args []string) error {
	helper := newJfrogHelper()
	status := helper.GetStatus()

	format, err := output.ParseFormat(jfrogOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(status.Servers)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Configured Servers"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if len(status.Servers) == 0 {
		fmt.Fprintln(os.Stdout, "No servers configured")
		return nil
	}

	for _, s := range status.Servers {
		marker := "  "
		if s.IsDefault {
			marker = output.Success("→ ")
		}
		fmt.Fprintf(os.Stdout, "%s%-15s %s\n", marker, s.ID, s.URL)
	}

	return nil
}

func runJfrogServersAdd(cmd *cobra.Command, args []string) error {
	helper := newJfrogHelper()
	return helper.AddServer(args[0], "", "", "", true)
}

func runJfrogServersRemove(cmd *cobra.Command, args []string) error {
	helper := newJfrogHelper()
	if err := helper.RemoveServer(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Server removed: %s\n", output.Success("✓"), args[0])
	return nil
}

func runJfrogServersUse(cmd *cobra.Command, args []string) error {
	helper := newJfrogHelper()
	if err := helper.UseServer(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Default server set: %s\n", output.Success("✓"), args[0])
	return nil
}

func runJfrogUpload(cmd *cobra.Command, args []string) error {
	helper := newJfrogHelper()
	return helper.Upload(args[0], args[1], jfrogRecursive, jfrogFlat)
}

func runJfrogDownload(cmd *cobra.Command, args []string) error {
	helper := newJfrogHelper()
	target := ""
	if len(args) > 1 {
		target = args[1]
	}
	return helper.Download(args[0], target, jfrogRecursive, jfrogFlat)
}

func runJfrogSearch(cmd *cobra.Command, args []string) error {
	helper := newJfrogHelper()
	artifacts, err := helper.Search(args[0])
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(jfrogOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(artifacts)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Search Results"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if len(artifacts) == 0 {
		fmt.Fprintln(os.Stdout, "No artifacts found")
		return nil
	}

	for _, a := range artifacts {
		fmt.Fprintf(os.Stdout, "  %s\n", a.Path)
	}

	fmt.Fprintf(os.Stdout, "\n%d artifacts found\n", len(artifacts))
	return nil
}

func runJfrogDelete(cmd *cobra.Command, args []string) error {
	helper := newJfrogHelper()
	return helper.Delete(args[0], jfrogQuiet)
}

func runJfrogDockerPush(cmd *cobra.Command, args []string) error {
	helper := newJfrogHelper()
	return helper.DockerPush(args[0], args[1])
}

func runJfrogDockerPull(cmd *cobra.Command, args []string) error {
	helper := newJfrogHelper()
	return helper.DockerPull(args[0], args[1])
}
