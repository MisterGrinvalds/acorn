package cmd

import (
	"github.com/mistergrinvalds/acorn/internal/components"
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/n8n"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/mistergrinvalds/acorn/internal/utils/configcmd"
	"github.com/spf13/cobra"
)

var (
	n8nDryRun    bool
	n8nVerbose   bool
	n8nUseDocker bool
	n8nPort      int
	n8nOutputDir string
)

// n8nCmd represents the n8n command group
var n8nCmd = &cobra.Command{
	Use:   "n8n",
	Short: "n8n workflow automation",
	Long: `n8n workflow automation commands.

Manage n8n installation, workflows, and local instances.

Examples:
  acorn automation n8n status      # Show status
  acorn automation n8n start       # Start n8n locally
  acorn automation n8n stop        # Stop n8n
  acorn automation n8n workflows   # List workflows`,
}

// n8nStatusCmd shows status
var n8nStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show n8n status",
	Long: `Show n8n installation and running status.

Examples:
  acorn automation n8n status
  acorn automation n8n status -o json`,
	RunE: runN8nStatus,
}

// n8nStartCmd starts n8n
var n8nStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start n8n locally",
	Long: `Start n8n workflow automation locally.

By default uses npm/npx. Use --docker to run in a container.

Examples:
  acorn automation n8n start
  acorn automation n8n start --port 8080
  acorn automation n8n start --docker`,
	RunE: runN8nStart,
}

// n8nStopCmd stops n8n
var n8nStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop n8n",
	Long: `Stop running n8n instance.

Examples:
  acorn automation n8n stop`,
	RunE: runN8nStop,
}

// n8nInstallCmd installs n8n
var n8nInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install n8n globally",
	Long: `Install n8n globally via npm.

Examples:
  acorn automation n8n install`,
	RunE: runN8nInstall,
}

// n8nWorkflowsCmd lists workflows
var n8nWorkflowsCmd = &cobra.Command{
	Use:   "workflows",
	Short: "List exported workflows",
	Long: `List workflows that have been exported.

Examples:
  acorn automation n8n workflows
  acorn automation n8n workflows -o json`,
	Aliases: []string{"ls"},
	RunE:    runN8nWorkflows,
}

// n8nExportCmd exports workflows
var n8nExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export workflows",
	Long: `Export all workflows to files.

Examples:
  acorn automation n8n export
  acorn automation n8n export --output ~/backups/n8n`,
	RunE: runN8nExport,
}

// n8nImportCmd imports a workflow
var n8nImportCmd = &cobra.Command{
	Use:   "import <file>",
	Short: "Import a workflow",
	Long: `Import a workflow from a JSON file.

Examples:
  acorn automation n8n import ./my-workflow.json`,
	Args: cobra.ExactArgs(1),
	RunE: runN8nImport,
}

// n8nExportCredsCmd exports credentials
var n8nExportCredsCmd = &cobra.Command{
	Use:   "export-credentials",
	Short: "Export credentials (encrypted)",
	Long: `Export all credentials to files (encrypted).

Examples:
  acorn automation n8n export-credentials`,
	Aliases: []string{"export-creds"},
	RunE:    runN8nExportCreds,
}

func init() {

	// Add subcommands
	n8nCmd.AddCommand(n8nStatusCmd)
	n8nCmd.AddCommand(n8nStartCmd)
	n8nCmd.AddCommand(n8nStopCmd)
	n8nCmd.AddCommand(n8nInstallCmd)
	n8nCmd.AddCommand(n8nWorkflowsCmd)
	n8nCmd.AddCommand(n8nExportCmd)
	n8nCmd.AddCommand(n8nImportCmd)
	n8nCmd.AddCommand(n8nExportCredsCmd)
	n8nCmd.AddCommand(configcmd.NewConfigRouter("n8n"))

	// Start flags
	n8nStartCmd.Flags().BoolVar(&n8nUseDocker, "docker", false, "Use Docker instead of npm")
	n8nStartCmd.Flags().IntVar(&n8nPort, "port", 5678, "Port to run on")

	// Export flags
	n8nExportCmd.Flags().StringVarP(&n8nOutputDir, "output", "O", "", "Output directory")
	n8nExportCredsCmd.Flags().StringVarP(&n8nOutputDir, "output", "O", "", "Output directory")

	// Persistent flags
	n8nCmd.PersistentFlags().BoolVar(&n8nDryRun, "dry-run", false,
		"Show what would be done without executing")
	n8nCmd.PersistentFlags().BoolVarP(&n8nVerbose, "verbose", "v", false,
		"Show verbose output")
}

func newN8nHelper() *n8n.Helper {
	return n8n.NewHelper(n8nVerbose, n8nDryRun)
}

func runN8nStatus(cmd *cobra.Command, args []string) error {
	helper := newN8nHelper()
	status := helper.GetStatus()

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("n8n Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if status.Installed {
		version := status.Version
		if version == "" {
			version = "available"
		}
		fmt.Fprintf(os.Stdout, "%s n8n: %s (%s)\n", output.Success("✓"), version, status.Method)
	} else {
		fmt.Fprintf(os.Stdout, "%s n8n not installed\n", output.Error("✗"))
		fmt.Fprintln(os.Stdout, "  Install: acorn automation n8n install")
		fmt.Fprintln(os.Stdout, "  Or run with: npx n8n")
	}

	fmt.Fprintln(os.Stdout)

	if status.DockerRunning {
		fmt.Fprintf(os.Stdout, "%s Docker container running: %s\n", output.Success("✓"), status.ContainerID[:12])
	} else {
		fmt.Fprintf(os.Stdout, "%s No Docker container running\n", output.Info("ℹ"))
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "Workflows dir: %s\n", status.WorkflowsDir)
	fmt.Fprintf(os.Stdout, "Credentials dir: %s\n", status.CredentialsDir)

	return nil
}

func runN8nStart(cmd *cobra.Command, args []string) error {
	helper := newN8nHelper()
	return helper.Start(n8nUseDocker, n8nPort)
}

func runN8nStop(cmd *cobra.Command, args []string) error {
	helper := newN8nHelper()
	if err := helper.Stop(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s n8n stopped\n", output.Success("✓"))
	return nil
}

func runN8nInstall(cmd *cobra.Command, args []string) error {
	helper := newN8nHelper()
	if err := helper.Install(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s n8n installed\n", output.Success("✓"))
	return nil
}

func runN8nWorkflows(cmd *cobra.Command, args []string) error {
	helper := newN8nHelper()
	workflows, err := helper.ListWorkflows()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(workflows)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Exported Workflows"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if len(workflows) == 0 {
		fmt.Fprintln(os.Stdout, "No exported workflows found")
		fmt.Fprintln(os.Stdout, "  Export: acorn automation n8n export")
		return nil
	}

	for _, w := range workflows {
		status := output.Info("○")
		if w.Active {
			status = output.Success("●")
		}
		fmt.Fprintf(os.Stdout, "  %s %s\n", status, w.Name)
		if w.ID != "" {
			fmt.Fprintf(os.Stdout, "    ID: %s\n", w.ID)
		}
	}

	return nil
}

func runN8nExport(cmd *cobra.Command, args []string) error {
	helper := newN8nHelper()
	if err := helper.ExportWorkflows(n8nOutputDir); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Workflows exported\n", output.Success("✓"))
	return nil
}

func runN8nImport(cmd *cobra.Command, args []string) error {
	helper := newN8nHelper()
	if err := helper.ImportWorkflow(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Workflow imported\n", output.Success("✓"))
	return nil
}

func runN8nExportCreds(cmd *cobra.Command, args []string) error {
	helper := newN8nHelper()
	if err := helper.ExportCredentials(n8nOutputDir); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Credentials exported\n", output.Success("✓"))
	return nil
}

func init() {
	components.Register(&components.Registration{
		Name: "n8n",
		RegisterCmd: func() *cobra.Command { return n8nCmd },
	})
}
