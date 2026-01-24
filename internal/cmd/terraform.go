package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/cloud/terraform"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	tfOutputFormat string
	tfDryRun       bool
	tfVerbose      bool
	tfAutoApprove  bool
	tfUpgrade      bool
	tfPlanOut      string
	tfCheck        bool
)

// terraformCmd represents the terraform command group
var terraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Terraform Infrastructure as Code",
	Long: `Terraform Infrastructure as Code commands.

Manage infrastructure with HashiCorp Terraform.

Examples:
  acorn cloud terraform status       # Show status
  acorn cloud terraform init         # Initialize
  acorn cloud terraform plan         # Plan changes
  acorn cloud terraform apply        # Apply changes`,
	Aliases: []string{"tf"},
}

// tfStatusCmd shows status
var tfStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Terraform status",
	Long: `Show Terraform installation and workspace status.

Examples:
  acorn cloud terraform status
  acorn cloud terraform status -o json`,
	RunE: runTfStatus,
}

// tfInitCmd initializes Terraform
var tfInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Terraform",
	Long: `Initialize Terraform working directory.

Examples:
  acorn cloud terraform init
  acorn cloud terraform init --upgrade`,
	RunE: runTfInit,
}

// tfPlanCmd runs terraform plan
var tfPlanCmd = &cobra.Command{
	Use:   "plan",
	Short: "Plan changes",
	Long: `Show execution plan without applying.

Examples:
  acorn cloud terraform plan
  acorn cloud terraform plan --out plan.tfplan`,
	RunE: runTfPlan,
}

// tfApplyCmd runs terraform apply
var tfApplyCmd = &cobra.Command{
	Use:   "apply [plan-file]",
	Short: "Apply changes",
	Long: `Apply infrastructure changes.

Examples:
  acorn cloud terraform apply
  acorn cloud terraform apply --yes
  acorn cloud terraform apply plan.tfplan`,
	Args: cobra.MaximumNArgs(1),
	RunE: runTfApply,
}

// tfDestroyCmd runs terraform destroy
var tfDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy resources",
	Long: `Destroy all managed infrastructure.

Examples:
  acorn cloud terraform destroy
  acorn cloud terraform destroy --yes`,
	RunE: runTfDestroy,
}

// tfValidateCmd validates configuration
var tfValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration",
	Long: `Validate Terraform configuration files.

Examples:
  acorn cloud terraform validate`,
	RunE: runTfValidate,
}

// tfFmtCmd formats configuration
var tfFmtCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Format configuration",
	Long: `Format Terraform configuration files.

Examples:
  acorn cloud terraform fmt
  acorn cloud terraform fmt --check`,
	Aliases: []string{"format"},
	RunE:    runTfFmt,
}

// tfWorkspacesCmd lists workspaces
var tfWorkspacesCmd = &cobra.Command{
	Use:   "workspaces",
	Short: "List workspaces",
	Long: `List all Terraform workspaces.

Examples:
  acorn cloud terraform workspaces
  acorn cloud terraform workspaces -o json`,
	Aliases: []string{"ws"},
	RunE:    runTfWorkspaces,
}

// tfSelectCmd selects a workspace
var tfSelectCmd = &cobra.Command{
	Use:   "select <workspace>",
	Short: "Select workspace",
	Long: `Select a Terraform workspace.

Examples:
  acorn cloud terraform select dev
  acorn cloud terraform select prod`,
	Args: cobra.ExactArgs(1),
	RunE: runTfSelect,
}

// tfOutputsCmd shows outputs
var tfOutputsCmd = &cobra.Command{
	Use:   "outputs",
	Short: "Show outputs",
	Long: `Show Terraform outputs.

Examples:
  acorn cloud terraform outputs
  acorn cloud terraform outputs -o json`,
	RunE: runTfOutputs,
}

// tfStateCmd shows state
var tfStateCmd = &cobra.Command{
	Use:   "state",
	Short: "Show state resources",
	Long: `List resources in the Terraform state.

Examples:
  acorn cloud terraform state
  acorn cloud terraform state -o json`,
	RunE: runTfState,
}

// tfRefreshCmd refreshes state
var tfRefreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh state",
	Long: `Refresh Terraform state to match actual infrastructure.

Examples:
  acorn cloud terraform refresh`,
	RunE: runTfRefresh,
}

// tfInstallCmd installs Terraform
var tfInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Terraform",
	Long: `Install Terraform using Homebrew.

Examples:
  acorn cloud terraform install`,
	RunE: runTfInstall,
}

func init() {
	cloudCmd.AddCommand(terraformCmd)

	// Add subcommands
	terraformCmd.AddCommand(tfStatusCmd)
	terraformCmd.AddCommand(tfInitCmd)
	terraformCmd.AddCommand(tfPlanCmd)
	terraformCmd.AddCommand(tfApplyCmd)
	terraformCmd.AddCommand(tfDestroyCmd)
	terraformCmd.AddCommand(tfValidateCmd)
	terraformCmd.AddCommand(tfFmtCmd)
	terraformCmd.AddCommand(tfWorkspacesCmd)
	terraformCmd.AddCommand(tfSelectCmd)
	terraformCmd.AddCommand(tfOutputsCmd)
	terraformCmd.AddCommand(tfStateCmd)
	terraformCmd.AddCommand(tfRefreshCmd)
	terraformCmd.AddCommand(tfInstallCmd)

	// Init flags
	tfInitCmd.Flags().BoolVar(&tfUpgrade, "upgrade", false, "Upgrade providers")

	// Plan flags
	tfPlanCmd.Flags().StringVar(&tfPlanOut, "out", "", "Output plan file")

	// Apply/Destroy flags
	tfApplyCmd.Flags().BoolVarP(&tfAutoApprove, "yes", "y", false, "Skip confirmation")
	tfDestroyCmd.Flags().BoolVarP(&tfAutoApprove, "yes", "y", false, "Skip confirmation")

	// Format flags
	tfFmtCmd.Flags().BoolVar(&tfCheck, "check", false, "Check formatting without changing")

	// Persistent flags
	terraformCmd.PersistentFlags().StringVarP(&tfOutputFormat, "output", "o", "table",
		"Output format (table|json|yaml)")
	terraformCmd.PersistentFlags().BoolVar(&tfDryRun, "dry-run", false,
		"Show what would be done without executing")
	terraformCmd.PersistentFlags().BoolVarP(&tfVerbose, "verbose", "v", false,
		"Show verbose output")
}

func newTfHelper() *terraform.Helper {
	return terraform.NewHelper(tfVerbose, tfDryRun)
}

func runTfStatus(cmd *cobra.Command, args []string) error {
	helper := newTfHelper()
	status := helper.GetStatus()

	format, err := output.ParseFormat(tfOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Terraform Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if status.Installed {
		fmt.Fprintf(os.Stdout, "%s Terraform installed: v%s\n", output.Success("✓"), status.Version)
	} else {
		fmt.Fprintf(os.Stdout, "%s Terraform not installed\n", output.Error("✗"))
		fmt.Fprintln(os.Stdout, "  Install: acorn cloud terraform install")
		return nil
	}

	if status.TerragruntInstalled {
		fmt.Fprintf(os.Stdout, "%s Terragrunt installed: v%s\n", output.Success("✓"), status.TerragruntVersion)
	}

	fmt.Fprintln(os.Stdout)

	if helper.HasTerraformFiles() {
		if status.Initialized {
			fmt.Fprintf(os.Stdout, "%s Working directory initialized\n", output.Success("✓"))
			if status.Workspace != "" {
				fmt.Fprintf(os.Stdout, "  Workspace: %s\n", status.Workspace)
			}
		} else {
			fmt.Fprintf(os.Stdout, "%s Working directory not initialized\n", output.Warning("!"))
			fmt.Fprintln(os.Stdout, "  Run: acorn cloud terraform init")
		}
	} else {
		fmt.Fprintf(os.Stdout, "%s No Terraform files in current directory\n", output.Info("ℹ"))
	}

	return nil
}

func runTfInit(cmd *cobra.Command, args []string) error {
	helper := newTfHelper()
	return helper.Init(tfUpgrade)
}

func runTfPlan(cmd *cobra.Command, args []string) error {
	helper := newTfHelper()
	return helper.Plan(tfPlanOut)
}

func runTfApply(cmd *cobra.Command, args []string) error {
	helper := newTfHelper()
	planFile := ""
	if len(args) > 0 {
		planFile = args[0]
	}
	return helper.Apply(planFile, tfAutoApprove)
}

func runTfDestroy(cmd *cobra.Command, args []string) error {
	helper := newTfHelper()
	return helper.Destroy(tfAutoApprove)
}

func runTfValidate(cmd *cobra.Command, args []string) error {
	helper := newTfHelper()
	if err := helper.Validate(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Configuration valid\n", output.Success("✓"))
	return nil
}

func runTfFmt(cmd *cobra.Command, args []string) error {
	helper := newTfHelper()
	return helper.Format(tfCheck)
}

func runTfWorkspaces(cmd *cobra.Command, args []string) error {
	helper := newTfHelper()
	workspaces, err := helper.ListWorkspaces()
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(tfOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(workspaces)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Workspaces"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if len(workspaces) == 0 {
		fmt.Fprintln(os.Stdout, "No workspaces found")
		return nil
	}

	for _, w := range workspaces {
		marker := "  "
		if w.Current {
			marker = output.Success("→ ")
		}
		fmt.Fprintf(os.Stdout, "%s%s\n", marker, w.Name)
	}

	return nil
}

func runTfSelect(cmd *cobra.Command, args []string) error {
	helper := newTfHelper()
	if err := helper.SelectWorkspace(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Selected workspace: %s\n", output.Success("✓"), args[0])
	return nil
}

func runTfOutputs(cmd *cobra.Command, args []string) error {
	helper := newTfHelper()
	outputs, err := helper.GetOutputs()
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(tfOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(outputs)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Outputs"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if len(outputs) == 0 {
		fmt.Fprintln(os.Stdout, "No outputs")
		return nil
	}

	for name, out := range outputs {
		if out.Sensitive {
			fmt.Fprintf(os.Stdout, "  %s: <sensitive>\n", name)
		} else {
			fmt.Fprintf(os.Stdout, "  %s: %v\n", name, out.Value)
		}
	}

	return nil
}

func runTfState(cmd *cobra.Command, args []string) error {
	helper := newTfHelper()
	resources, err := helper.GetState()
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(tfOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(resources)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("State Resources"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if len(resources) == 0 {
		fmt.Fprintln(os.Stdout, "No resources in state")
		return nil
	}

	for _, r := range resources {
		fmt.Fprintf(os.Stdout, "  %s\n", r.Address)
	}

	fmt.Fprintf(os.Stdout, "\n%d resources\n", len(resources))
	return nil
}

func runTfRefresh(cmd *cobra.Command, args []string) error {
	helper := newTfHelper()
	if err := helper.Refresh(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s State refreshed\n", output.Success("✓"))
	return nil
}

func runTfInstall(cmd *cobra.Command, args []string) error {
	helper := newTfHelper()
	if err := helper.Install(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Terraform installed\n", output.Success("✓"))
	return nil
}
