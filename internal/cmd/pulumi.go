package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/cloud/pulumi"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/mistergrinvalds/acorn/internal/utils/configcmd"
	"github.com/spf13/cobra"
)

var (
	pulumiDryRun  bool
	pulumiVerbose bool
	pulumiYes     bool
	pulumiBackend string
)

// pulumiCmd represents the pulumi command group
var pulumiCmd = &cobra.Command{
	Use:   "pulumi",
	Short: "Pulumi Infrastructure as Code",
	Long: `Pulumi Infrastructure as Code commands.

Manage infrastructure with Pulumi using your favorite programming languages.

Examples:
  acorn cloud pulumi status      # Show status
  acorn cloud pulumi stacks      # List stacks
  acorn cloud pulumi preview     # Preview changes
  acorn cloud pulumi up          # Deploy changes`,
	Aliases: []string{"pu"},
}

// pulumiStatusCmd shows status
var pulumiStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Pulumi status",
	Long: `Show Pulumi installation and authentication status.

Examples:
  acorn cloud pulumi status
  acorn cloud pulumi status -o json`,
	RunE: runPulumiStatus,
}

// pulumiLoginCmd logs in
var pulumiLoginCmd = &cobra.Command{
	Use:   "login [backend]",
	Short: "Log in to Pulumi",
	Long: `Log in to Pulumi Cloud or a self-hosted backend.

Examples:
  acorn cloud pulumi login                    # Pulumi Cloud
  acorn cloud pulumi login --local            # Local backend
  acorn cloud pulumi login s3://my-bucket     # S3 backend`,
	RunE: runPulumiLogin,
}

// pulumiLogoutCmd logs out
var pulumiLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out from Pulumi",
	Long: `Log out from the current Pulumi backend.

Examples:
  acorn cloud pulumi logout`,
	RunE: runPulumiLogout,
}

// pulumiStacksCmd lists stacks
var pulumiStacksCmd = &cobra.Command{
	Use:   "stacks",
	Short: "List stacks",
	Long: `List all stacks in the current project.

Examples:
  acorn cloud pulumi stacks
  acorn cloud pulumi stacks -o json`,
	Aliases: []string{"ls"},
	RunE:    runPulumiStacks,
}

// pulumiSelectCmd selects a stack
var pulumiSelectCmd = &cobra.Command{
	Use:   "select <stack>",
	Short: "Select a stack",
	Long: `Select a stack to use for subsequent operations.

Examples:
  acorn cloud pulumi select dev
  acorn cloud pulumi select prod`,
	Args: cobra.ExactArgs(1),
	RunE: runPulumiSelect,
}

// pulumiPreviewCmd previews changes
var pulumiPreviewCmd = &cobra.Command{
	Use:   "preview",
	Short: "Preview changes",
	Long: `Preview infrastructure changes without applying them.

Examples:
  acorn cloud pulumi preview`,
	RunE: runPulumiPreview,
}

// pulumiUpCmd deploys changes
var pulumiUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Deploy changes",
	Long: `Deploy infrastructure changes.

Examples:
  acorn cloud pulumi up
  acorn cloud pulumi up --yes`,
	Aliases: []string{"deploy"},
	RunE:    runPulumiUp,
}

// pulumiDestroyCmd destroys resources
var pulumiDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy resources",
	Long: `Destroy all resources in the current stack.

Examples:
  acorn cloud pulumi destroy
  acorn cloud pulumi destroy --yes`,
	RunE: runPulumiDestroy,
}

// pulumiRefreshCmd refreshes state
var pulumiRefreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh state",
	Long: `Refresh the state to match actual infrastructure.

Examples:
  acorn cloud pulumi refresh`,
	RunE: runPulumiRefresh,
}

// pulumiOutputsCmd shows outputs
var pulumiOutputsCmd = &cobra.Command{
	Use:   "outputs",
	Short: "Show stack outputs",
	Long: `Show outputs from the current stack.

Examples:
  acorn cloud pulumi outputs
  acorn cloud pulumi outputs -o json`,
	RunE: runPulumiOutputs,
}

// pulumiNewCmd creates a new project
var pulumiNewCmd = &cobra.Command{
	Use:   "new <template>",
	Short: "Create new project",
	Long: `Create a new Pulumi project from a template.

Common templates:
  aws-go, aws-typescript, aws-python
  azure-go, azure-typescript, azure-python
  gcp-go, gcp-typescript, gcp-python
  kubernetes-go, kubernetes-typescript

Examples:
  acorn cloud pulumi new aws-go
  acorn cloud pulumi new kubernetes-typescript`,
	Args: cobra.ExactArgs(1),
	RunE: runPulumiNew,
}

// pulumiInstallCmd installs Pulumi
var pulumiInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Pulumi CLI",
	Long: `Install Pulumi CLI using Homebrew.

Examples:
  acorn cloud pulumi install`,
	RunE: runPulumiInstall,
}

func init() {
	cloudCmd.AddCommand(pulumiCmd)

	// Add subcommands
	pulumiCmd.AddCommand(pulumiStatusCmd)
	pulumiCmd.AddCommand(pulumiLoginCmd)
	pulumiCmd.AddCommand(pulumiLogoutCmd)
	pulumiCmd.AddCommand(pulumiStacksCmd)
	pulumiCmd.AddCommand(pulumiSelectCmd)
	pulumiCmd.AddCommand(pulumiPreviewCmd)
	pulumiCmd.AddCommand(pulumiUpCmd)
	pulumiCmd.AddCommand(pulumiDestroyCmd)
	pulumiCmd.AddCommand(pulumiRefreshCmd)
	pulumiCmd.AddCommand(pulumiOutputsCmd)
	pulumiCmd.AddCommand(pulumiNewCmd)
	pulumiCmd.AddCommand(pulumiInstallCmd)
	pulumiCmd.AddCommand(configcmd.NewConfigRouter("pulumi"))

	// Login flags
	pulumiLoginCmd.Flags().StringVar(&pulumiBackend, "backend", "", "Backend URL")
	pulumiLoginCmd.Flags().Bool("local", false, "Use local backend")

	// Up/Destroy flags
	pulumiUpCmd.Flags().BoolVarP(&pulumiYes, "yes", "y", false, "Skip confirmation")
	pulumiDestroyCmd.Flags().BoolVarP(&pulumiYes, "yes", "y", false, "Skip confirmation")

	// Persistent flags
	pulumiCmd.PersistentFlags().BoolVar(&pulumiDryRun, "dry-run", false,
		"Show what would be done without executing")
	pulumiCmd.PersistentFlags().BoolVarP(&pulumiVerbose, "verbose", "v", false,
		"Show verbose output")
}

func newPulumiHelper() *pulumi.Helper {
	return pulumi.NewHelper(pulumiVerbose, pulumiDryRun)
}

func runPulumiStatus(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := newPulumiHelper()
	status := helper.GetStatus()

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Pulumi Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if status.Installed {
		fmt.Fprintf(os.Stdout, "%s Pulumi installed: v%s\n", output.Success("✓"), status.Version)
	} else {
		fmt.Fprintf(os.Stdout, "%s Pulumi not installed\n", output.Error("✗"))
		fmt.Fprintln(os.Stdout, "  Install: acorn cloud pulumi install")
		return nil
	}

	fmt.Fprintln(os.Stdout)

	if status.LoggedIn {
		fmt.Fprintf(os.Stdout, "%s Logged in as: %s\n", output.Success("✓"), status.User)
		if status.Backend != "" {
			fmt.Fprintf(os.Stdout, "  Backend: %s\n", status.Backend)
		}
		if status.Organization != "" {
			fmt.Fprintf(os.Stdout, "  Organization: %s\n", status.Organization)
		}
	} else {
		fmt.Fprintf(os.Stdout, "%s Not logged in\n", output.Warning("!"))
		fmt.Fprintln(os.Stdout, "  Login: acorn cloud pulumi login")
	}

	// Check for project
	if project, err := helper.GetProject(); err == nil {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintf(os.Stdout, "Project: %s (%s)\n", project.Name, project.Runtime)
	}

	return nil
}

func runPulumiLogin(cmd *cobra.Command, args []string) error {
	helper := newPulumiHelper()

	backend := pulumiBackend
	if len(args) > 0 {
		backend = args[0]
	}

	local, _ := cmd.Flags().GetBool("local")
	if local {
		backend = "--local"
	}

	return helper.Login(backend)
}

func runPulumiLogout(cmd *cobra.Command, args []string) error {
	helper := newPulumiHelper()
	if err := helper.Logout(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Logged out\n", output.Success("✓"))
	return nil
}

func runPulumiStacks(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := newPulumiHelper()
	stacks, err := helper.ListStacks()
	if err != nil {
		return err
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(stacks)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Stacks"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if len(stacks) == 0 {
		fmt.Fprintln(os.Stdout, "No stacks found")
		return nil
	}

	for _, s := range stacks {
		marker := "  "
		if s.Current {
			marker = output.Success("→ ")
		}
		fmt.Fprintf(os.Stdout, "%s%s\n", marker, s.Name)
	}

	return nil
}

func runPulumiSelect(cmd *cobra.Command, args []string) error {
	helper := newPulumiHelper()
	if err := helper.SelectStack(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Selected stack: %s\n", output.Success("✓"), args[0])
	return nil
}

func runPulumiPreview(cmd *cobra.Command, args []string) error {
	helper := newPulumiHelper()
	return helper.Preview()
}

func runPulumiUp(cmd *cobra.Command, args []string) error {
	helper := newPulumiHelper()
	return helper.Up(pulumiYes)
}

func runPulumiDestroy(cmd *cobra.Command, args []string) error {
	helper := newPulumiHelper()
	return helper.Destroy(pulumiYes)
}

func runPulumiRefresh(cmd *cobra.Command, args []string) error {
	helper := newPulumiHelper()
	return helper.Refresh()
}

func runPulumiOutputs(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := newPulumiHelper()
	outputs, err := helper.GetOutputs()
	if err != nil {
		return err
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(outputs)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Stack Outputs"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if len(outputs) == 0 {
		fmt.Fprintln(os.Stdout, "No outputs")
		return nil
	}

	for k, v := range outputs {
		fmt.Fprintf(os.Stdout, "  %s: %v\n", k, v)
	}

	return nil
}

func runPulumiNew(cmd *cobra.Command, args []string) error {
	helper := newPulumiHelper()
	return helper.NewProject(args[0], "", "")
}

func runPulumiInstall(cmd *cobra.Command, args []string) error {
	helper := newPulumiHelper()
	if err := helper.Install(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Pulumi installed\n", output.Success("✓"))
	return nil
}
