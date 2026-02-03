package cmd

import (
	"github.com/mistergrinvalds/acorn/internal/components"
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/aws"
	"github.com/mistergrinvalds/acorn/internal/utils/installer"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	awsDryRun  bool
	awsVerbose bool
	awsProfile string
	awsRegion  string
)

// awsCmd represents the aws command group
var awsCmd = &cobra.Command{
	Use:   "aws",
	Short: "AWS CLI helpers",
	Long: `Helpers for AWS CLI operations.

Provides commands for managing AWS resources across EC2, S3, Lambda, EKS, and more.

Examples:
  acorn cloud aws status          # Check AWS CLI status and auth
  acorn cloud aws whoami          # Show current identity
  acorn cloud aws profiles        # List configured profiles
  acorn cloud aws ec2 list        # List EC2 instances
  acorn cloud aws s3 list         # List S3 buckets`,
	Aliases: []string{"amazon"},
}

// awsStatusCmd shows AWS CLI status
var awsStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check AWS CLI status and authentication",
	Long: `Check if AWS CLI is installed and authenticated.

Shows AWS CLI version, current profile, and account information.

Examples:
  acorn cloud aws status
  acorn cloud aws status -o json`,
	RunE: runAwsStatus,
}

// awsWhoamiCmd shows current identity
var awsWhoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current AWS identity",
	Long: `Display the current AWS identity (sts get-caller-identity).

Examples:
  acorn cloud aws whoami`,
	RunE: runAwsWhoami,
}

// awsProfilesCmd lists profiles
var awsProfilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "List configured AWS profiles",
	Long: `List all configured AWS profiles.

Examples:
  acorn cloud aws profiles`,
	RunE: runAwsProfiles,
}

// awsRegionsCmd lists regions
var awsRegionsCmd = &cobra.Command{
	Use:   "regions",
	Short: "List available AWS regions",
	Long: `List all available AWS regions.

Examples:
  acorn cloud aws regions`,
	RunE: runAwsRegions,
}

// awsOverviewCmd shows overview
var awsOverviewCmd = &cobra.Command{
	Use:   "overview",
	Short: "Show overview of AWS resources",
	Long: `Display an overview of AWS resources including
EC2 instances, S3 buckets, Lambda functions, and EKS clusters.

Examples:
  acorn cloud aws overview
  acorn cloud aws overview -o json`,
	RunE: runAwsOverview,
}

// awsLoginCmd performs SSO login
var awsLoginCmd = &cobra.Command{
	Use:   "login [profile]",
	Short: "Login to AWS (SSO)",
	Long: `Perform AWS SSO login for a profile.

Examples:
  acorn cloud aws login
  acorn cloud aws login my-profile`,
	Args: cobra.MaximumNArgs(1),
	RunE: runAwsLogin,
}

// awsConfigureCmd runs aws configure
var awsConfigureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure AWS credentials",
	Long: `Run aws configure to set up credentials.

Examples:
  acorn cloud aws configure`,
	RunE: runAwsConfigure,
}

// EC2 subcommands
var awsEC2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "EC2 instance commands",
	Long: `Commands for managing AWS EC2 instances.

Examples:
  acorn cloud aws ec2 list`,
}

var awsEC2ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List EC2 instances",
	Long: `List all EC2 instances.

Examples:
  acorn cloud aws ec2 list`,
	RunE: runAwsEC2List,
}

// S3 subcommands
var awsS3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "S3 bucket commands",
	Long: `Commands for managing AWS S3 buckets.

Examples:
  acorn cloud aws s3 list`,
}

var awsS3ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List S3 buckets",
	Long: `List all S3 buckets.

Examples:
  acorn cloud aws s3 list`,
	RunE: runAwsS3List,
}

// Lambda subcommands
var awsLambdaCmd = &cobra.Command{
	Use:   "lambda",
	Short: "Lambda function commands",
	Long: `Commands for managing AWS Lambda functions.

Examples:
  acorn cloud aws lambda list`,
}

var awsLambdaListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Lambda functions",
	Long: `List all Lambda functions.

Examples:
  acorn cloud aws lambda list`,
	RunE: runAwsLambdaList,
}

// EKS subcommands
var awsEKSCmd = &cobra.Command{
	Use:   "eks",
	Short: "EKS cluster commands",
	Long: `Commands for managing AWS EKS clusters.

Examples:
  acorn cloud aws eks list
  acorn cloud aws eks kubeconfig my-cluster`,
}

var awsEKSListCmd = &cobra.Command{
	Use:   "list",
	Short: "List EKS clusters",
	Long: `List all EKS clusters.

Examples:
  acorn cloud aws eks list`,
	RunE: runAwsEKSList,
}

var awsEKSKubeconfigCmd = &cobra.Command{
	Use:   "kubeconfig <cluster-name>",
	Short: "Update kubeconfig for EKS cluster",
	Long: `Update kubeconfig to use an EKS cluster.

Examples:
  acorn cloud aws eks kubeconfig my-cluster`,
	Args: cobra.ExactArgs(1),
	RunE: runAwsEKSKubeconfig,
}

// SSM subcommand
var awsSSMCmd = &cobra.Command{
	Use:   "ssm",
	Short: "SSM session commands",
	Long: `Commands for AWS Systems Manager sessions.

Examples:
  acorn cloud aws ssm connect i-1234567890abcdef`,
}

var awsSSMConnectCmd = &cobra.Command{
	Use:   "connect <instance-id>",
	Short: "Start SSM session to instance",
	Long: `Start an SSM session to an EC2 instance.

Examples:
  acorn cloud aws ssm connect i-1234567890abcdef`,
	Args: cobra.ExactArgs(1),
	RunE: runAwsSSMConnect,
}

// awsInstallCmd installs AWS CLI
var awsInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install AWS CLI tools",
	Long: `Install AWS CLI and related tools.

Automatically detects your platform and uses the appropriate
package manager (brew on macOS, script install on Linux).

Examples:
  acorn cloud aws install           # Install AWS CLI
  acorn cloud aws install --dry-run # Show what would be installed
  acorn cloud aws install -v        # Verbose output`,
	RunE: runAwsInstall,
}

func init() {

	// Add subcommands
	awsCmd.AddCommand(awsInstallCmd)
	awsCmd.AddCommand(awsStatusCmd)
	awsCmd.AddCommand(awsWhoamiCmd)
	awsCmd.AddCommand(awsProfilesCmd)
	awsCmd.AddCommand(awsRegionsCmd)
	awsCmd.AddCommand(awsOverviewCmd)
	awsCmd.AddCommand(awsLoginCmd)
	awsCmd.AddCommand(awsConfigureCmd)

	// EC2 subcommands
	awsCmd.AddCommand(awsEC2Cmd)
	awsEC2Cmd.AddCommand(awsEC2ListCmd)

	// S3 subcommands
	awsCmd.AddCommand(awsS3Cmd)
	awsS3Cmd.AddCommand(awsS3ListCmd)

	// Lambda subcommands
	awsCmd.AddCommand(awsLambdaCmd)
	awsLambdaCmd.AddCommand(awsLambdaListCmd)

	// EKS subcommands
	awsCmd.AddCommand(awsEKSCmd)
	awsEKSCmd.AddCommand(awsEKSListCmd)
	awsEKSCmd.AddCommand(awsEKSKubeconfigCmd)

	// SSM subcommands
	awsCmd.AddCommand(awsSSMCmd)
	awsSSMCmd.AddCommand(awsSSMConnectCmd)

	// Persistent flags
	awsCmd.PersistentFlags().BoolVar(&awsDryRun, "dry-run", false,
		"Show what would be done without executing")
	awsCmd.PersistentFlags().BoolVarP(&awsVerbose, "verbose", "v", false,
		"Show verbose output")
	awsCmd.PersistentFlags().StringVarP(&awsProfile, "profile", "p", "",
		"AWS profile to use")
	awsCmd.PersistentFlags().StringVarP(&awsRegion, "region", "r", "",
		"AWS region to use")
}

func newAwsHelper() *aws.Helper {
	helper := aws.NewHelper(awsVerbose, awsDryRun)
	if awsProfile != "" {
		helper.SetProfile(awsProfile)
	}
	if awsRegion != "" {
		helper.SetRegion(awsRegion)
	}
	return helper
}

func runAwsStatus(cmd *cobra.Command, args []string) error {
	helper := newAwsHelper()
	status, err := helper.GetStatus()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("AWS CLI Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if status.Installed {
		fmt.Fprintf(os.Stdout, "%s AWS CLI installed: %s\n", output.Success("✓"), status.Version)
	} else {
		fmt.Fprintf(os.Stdout, "%s AWS CLI not found\n", output.Error("✗"))
		fmt.Fprintln(os.Stdout, "  Install: brew install awscli (macOS)")
		return nil
	}

	fmt.Fprintln(os.Stdout)
	if status.Profile != "" {
		fmt.Fprintf(os.Stdout, "Profile: %s\n", status.Profile)
	}
	if status.Region != "" {
		fmt.Fprintf(os.Stdout, "Region: %s\n", status.Region)
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Authentication:"))
	if status.Authenticated {
		fmt.Fprintf(os.Stdout, "%s Authenticated\n", output.Success("✓"))
		if status.AccountID != "" {
			fmt.Fprintf(os.Stdout, "  Account: %s\n", status.AccountID)
		}
		if status.UserARN != "" {
			fmt.Fprintf(os.Stdout, "  ARN: %s\n", status.UserARN)
		}
	} else {
		fmt.Fprintf(os.Stdout, "%s Not authenticated\n", output.Warning("⚠"))
		fmt.Fprintln(os.Stdout, "  Run: aws configure")
	}

	return nil
}

func runAwsWhoami(cmd *cobra.Command, args []string) error {
	helper := newAwsHelper()
	identity, err := helper.GetCallerIdentity()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(identity)
	}

	fmt.Fprintf(os.Stdout, "Account: %s\n", identity.Account)
	fmt.Fprintf(os.Stdout, "User ID: %s\n", identity.UserID)
	fmt.Fprintf(os.Stdout, "ARN:     %s\n", identity.Arn)
	return nil
}

func runAwsProfiles(cmd *cobra.Command, args []string) error {
	helper := newAwsHelper()
	profiles, err := helper.ListProfiles()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(profiles)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("AWS Profiles"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	if len(profiles) == 0 {
		fmt.Fprintln(os.Stdout, "No profiles configured")
		return nil
	}
	for _, p := range profiles {
		fmt.Fprintf(os.Stdout, "  %s\n", p)
	}
	return nil
}

func runAwsRegions(cmd *cobra.Command, args []string) error {
	helper := newAwsHelper()
	regions, err := helper.ListRegions()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(regions)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("AWS Regions"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	for _, r := range regions {
		fmt.Fprintf(os.Stdout, "  %s\n", r)
	}
	return nil
}

func runAwsOverview(cmd *cobra.Command, args []string) error {
	helper := newAwsHelper()
	overview, err := helper.GetOverview()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(overview)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("AWS Overview"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintln(os.Stdout)

	if !overview.Status.Installed {
		fmt.Fprintf(os.Stdout, "%s AWS CLI not installed\n", output.Error("✗"))
		return nil
	}

	if !overview.Status.Authenticated {
		fmt.Fprintf(os.Stdout, "%s Not authenticated. Run: aws configure\n", output.Warning("⚠"))
		return nil
	}

	if overview.Status.AccountID != "" {
		fmt.Fprintf(os.Stdout, "Account: %s\n", overview.Status.AccountID)
	}
	if overview.Status.Region != "" {
		fmt.Fprintf(os.Stdout, "Region: %s\n", overview.Status.Region)
	}
	fmt.Fprintln(os.Stdout)

	fmt.Fprintf(os.Stdout, "EC2 Instances (running): %d\n", overview.EC2Count)

	fmt.Fprintf(os.Stdout, "\n%s:\n", output.Info("S3 Buckets"))
	if len(overview.S3Buckets) == 0 {
		fmt.Fprintln(os.Stdout, "  None found")
	} else {
		for _, b := range overview.S3Buckets {
			fmt.Fprintf(os.Stdout, "  %s\n", b)
		}
	}

	fmt.Fprintf(os.Stdout, "\n%s:\n", output.Info("Lambda Functions"))
	if len(overview.Lambdas) == 0 {
		fmt.Fprintln(os.Stdout, "  None found")
	} else {
		for _, l := range overview.Lambdas {
			fmt.Fprintf(os.Stdout, "  %s\n", l)
		}
	}

	fmt.Fprintf(os.Stdout, "\n%s:\n", output.Info("EKS Clusters"))
	if len(overview.EKSClusters) == 0 {
		fmt.Fprintln(os.Stdout, "  None found")
	} else {
		for _, c := range overview.EKSClusters {
			fmt.Fprintf(os.Stdout, "  %s\n", c)
		}
	}

	return nil
}

func runAwsLogin(cmd *cobra.Command, args []string) error {
	helper := newAwsHelper()
	profile := ""
	if len(args) > 0 {
		profile = args[0]
	} else if awsProfile != "" {
		profile = awsProfile
	}
	return helper.SSOLogin(profile)
}

func runAwsConfigure(cmd *cobra.Command, args []string) error {
	helper := newAwsHelper()
	return helper.Configure()
}

func runAwsEC2List(cmd *cobra.Command, args []string) error {
	helper := newAwsHelper()

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("EC2 Instances"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	instances, err := helper.ListEC2Instances()
	if err != nil {
		fmt.Fprintln(os.Stdout, "No instances found or not authenticated")
		return nil
	}
	fmt.Fprintln(os.Stdout, instances)
	return nil
}

func runAwsS3List(cmd *cobra.Command, args []string) error {
	helper := newAwsHelper()
	buckets, err := helper.ListS3Buckets()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(buckets)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("S3 Buckets"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	if len(buckets) == 0 {
		fmt.Fprintln(os.Stdout, "No buckets found")
		return nil
	}
	for _, b := range buckets {
		fmt.Fprintf(os.Stdout, "  %s\n", b)
	}
	return nil
}

func runAwsLambdaList(cmd *cobra.Command, args []string) error {
	helper := newAwsHelper()
	functions, err := helper.ListLambdaFunctions()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(functions)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Lambda Functions"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	if len(functions) == 0 {
		fmt.Fprintln(os.Stdout, "No functions found")
		return nil
	}
	for _, f := range functions {
		fmt.Fprintf(os.Stdout, "  %s\n", f)
	}
	return nil
}

func runAwsEKSList(cmd *cobra.Command, args []string) error {
	helper := newAwsHelper()
	clusters, err := helper.ListEKSClusters()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(clusters)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("EKS Clusters"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	if len(clusters) == 0 {
		fmt.Fprintln(os.Stdout, "No clusters found")
		return nil
	}
	for _, c := range clusters {
		fmt.Fprintf(os.Stdout, "  %s\n", c)
	}
	return nil
}

func runAwsEKSKubeconfig(cmd *cobra.Command, args []string) error {
	helper := newAwsHelper()
	return helper.UpdateKubeconfig(args[0])
}

func runAwsSSMConnect(cmd *cobra.Command, args []string) error {
	helper := newAwsHelper()
	return helper.SSMConnect(args[0])
}

func runAwsInstall(cmd *cobra.Command, args []string) error {
	inst := installer.NewInstaller(
		installer.WithDryRun(awsDryRun),
		installer.WithVerbose(awsVerbose),
	)

	// Show platform info
	platform := inst.GetPlatform()
	if awsVerbose {
		fmt.Fprintf(os.Stdout, "Platform: %s\n\n", platform)
	}

	// Get the plan first
	plan, err := inst.Plan(cmd.Context(), "aws")
	if err != nil {
		return err
	}

	// Show what will be installed
	if awsDryRun {
		fmt.Fprintf(os.Stdout, "%s\n", output.Info("AWS Installation Plan"))
		fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Fprintf(os.Stdout, "Platform: %s\n\n", platform)
	}

	pending := plan.PendingTools()
	if len(pending) == 0 {
		fmt.Fprintf(os.Stdout, "%s All tools already installed\n", output.Success("✓"))
		return nil
	}

	// Show tools
	fmt.Fprintln(os.Stdout, "Tools:")
	for _, t := range plan.Tools {
		status := output.Warning("○")
		suffix := ""
		if t.AlreadyInstalled {
			status = output.Success("✓")
			suffix = " (installed)"
		} else if awsDryRun {
			suffix = fmt.Sprintf(" (via %s)", t.Method.Type)
		}
		fmt.Fprintf(os.Stdout, "  %s %s%s\n", status, t.Name, suffix)
	}

	if awsDryRun {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Run without --dry-run to install.")
		return nil
	}

	// Execute installation
	fmt.Fprintln(os.Stdout)
	result, err := inst.Install(cmd.Context(), "aws")
	if err != nil {
		return err
	}

	// Show results
	fmt.Fprintln(os.Stdout)
	installed, skipped, failed := result.Summary()

	if result.Success {
		fmt.Fprintf(os.Stdout, "%s Installation complete (%d installed, %d skipped)\n",
			output.Success("✓"), installed, skipped)
	} else {
		fmt.Fprintf(os.Stdout, "%s Installation failed (%d installed, %d skipped, %d failed)\n",
			output.Error("✗"), installed, skipped, failed)

		// Show errors
		for _, t := range result.Tools {
			if t.Error != nil {
				fmt.Fprintf(os.Stdout, "  %s: %s\n", t.Name, t.Error)
			}
		}
	}

	return nil
}

func init() {
	components.Register(&components.Registration{
		Name: "aws",
		RegisterCmd: func() *cobra.Command { return awsCmd },
	})
}
