package cmd

import (
	"github.com/mistergrinvalds/acorn/internal/components"
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/azure"
	"github.com/mistergrinvalds/acorn/internal/utils/installer"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/mistergrinvalds/acorn/internal/utils/configcmd"
	"github.com/spf13/cobra"
)

var (
	azureDryRun       bool
	azureVerbose      bool
	azureSubscription string
)

// azureCmd represents the azure command group
var azureCmd = &cobra.Command{
	Use:   "azure",
	Short: "Azure CLI helpers",
	Long: `Helpers for Azure CLI operations.

Provides commands for managing Azure resources across VMs, Storage, AKS, Functions, and more.

Examples:
  acorn cloud azure status         # Check Azure CLI status and auth
  acorn cloud azure whoami         # Show current account
  acorn cloud azure subscriptions  # List subscriptions
  acorn cloud azure vm list        # List VMs
  acorn cloud azure aks list       # List AKS clusters`,
	Aliases: []string{"az"},
}

// azureStatusCmd shows Azure CLI status
var azureStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check Azure CLI status and authentication",
	Long: `Check if Azure CLI is installed and authenticated.

Shows Azure CLI version, current subscription, and account information.

Examples:
  acorn cloud azure status
  acorn cloud azure status -o json`,
	RunE: runAzureStatus,
}

// azureWhoamiCmd shows current account
var azureWhoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current Azure account",
	Long: `Display the current Azure account information.

Examples:
  acorn cloud azure whoami`,
	RunE: runAzureWhoami,
}

// azureSubscriptionsCmd lists subscriptions
var azureSubscriptionsCmd = &cobra.Command{
	Use:   "subscriptions",
	Short: "List Azure subscriptions",
	Long: `List all Azure subscriptions.

Examples:
  acorn cloud azure subscriptions`,
	Aliases: []string{"subs"},
	RunE:    runAzureSubscriptions,
}

// azureSetSubscriptionCmd sets subscription
var azureSetSubscriptionCmd = &cobra.Command{
	Use:   "set-subscription <subscription-id>",
	Short: "Set active Azure subscription",
	Long: `Set the active Azure subscription.

Examples:
  acorn cloud azure set-subscription xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`,
	Args: cobra.ExactArgs(1),
	RunE: runAzureSetSubscription,
}

// azureOverviewCmd shows overview
var azureOverviewCmd = &cobra.Command{
	Use:   "overview",
	Short: "Show overview of Azure resources",
	Long: `Display an overview of Azure resources including
VMs, Storage Accounts, AKS clusters, and Resource Groups.

Examples:
  acorn cloud azure overview
  acorn cloud azure overview -o json`,
	RunE: runAzureOverview,
}

// azureLoginCmd performs login
var azureLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Azure",
	Long: `Initiate Azure login.

Examples:
  acorn cloud azure login`,
	RunE: runAzureLogin,
}

// azureLogoutCmd performs logout
var azureLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from Azure",
	Long: `Logout from Azure.

Examples:
  acorn cloud azure logout`,
	RunE: runAzureLogout,
}

// VM subcommands
var azureVMCmd = &cobra.Command{
	Use:   "vm",
	Short: "Virtual Machine commands",
	Long: `Commands for managing Azure Virtual Machines.

Examples:
  acorn cloud azure vm list`,
}

var azureVMListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Virtual Machines",
	Long: `List all Virtual Machines.

Examples:
  acorn cloud azure vm list`,
	RunE: runAzureVMList,
}

// Storage subcommands
var azureStorageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Storage account commands",
	Long: `Commands for managing Azure Storage Accounts.

Examples:
  acorn cloud azure storage list`,
}

var azureStorageListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Storage Accounts",
	Long: `List all Storage Accounts.

Examples:
  acorn cloud azure storage list`,
	RunE: runAzureStorageList,
}

// AKS subcommands
var azureAKSCmd = &cobra.Command{
	Use:   "aks",
	Short: "AKS cluster commands",
	Long: `Commands for managing Azure Kubernetes Service clusters.

Examples:
  acorn cloud azure aks list
  acorn cloud azure aks credentials my-cluster my-rg`,
}

var azureAKSListCmd = &cobra.Command{
	Use:   "list",
	Short: "List AKS clusters",
	Long: `List all AKS clusters.

Examples:
  acorn cloud azure aks list`,
	RunE: runAzureAKSList,
}

var azureAKSCredentialsCmd = &cobra.Command{
	Use:   "credentials <cluster-name> <resource-group>",
	Short: "Get AKS kubeconfig",
	Long: `Update kubeconfig to use an AKS cluster.

Examples:
  acorn cloud azure aks credentials my-cluster my-rg`,
	Args: cobra.ExactArgs(2),
	RunE: runAzureAKSCredentials,
}

// ResourceGroup subcommands
var azureRGCmd = &cobra.Command{
	Use:   "rg",
	Short: "Resource Group commands",
	Long: `Commands for managing Azure Resource Groups.

Examples:
  acorn cloud azure rg list`,
}

var azureRGListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Resource Groups",
	Long: `List all Resource Groups.

Examples:
  acorn cloud azure rg list`,
	RunE: runAzureRGList,
}

// azureInstallCmd installs Azure CLI
var azureInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Azure CLI tools",
	Long: `Install Azure CLI.

Automatically detects your platform and uses the appropriate
package manager (brew on macOS, script install on Linux).

Examples:
  acorn cloud azure install           # Install Azure CLI
  acorn cloud azure install --dry-run # Show what would be installed
  acorn cloud azure install -v        # Verbose output`,
	RunE: runAzureInstall,
}

func init() {

	// Add subcommands
	azureCmd.AddCommand(azureInstallCmd)
	azureCmd.AddCommand(azureStatusCmd)
	azureCmd.AddCommand(azureWhoamiCmd)
	azureCmd.AddCommand(azureSubscriptionsCmd)
	azureCmd.AddCommand(azureSetSubscriptionCmd)
	azureCmd.AddCommand(azureOverviewCmd)
	azureCmd.AddCommand(azureLoginCmd)
	azureCmd.AddCommand(azureLogoutCmd)

	// VM subcommands
	azureCmd.AddCommand(azureVMCmd)
	azureVMCmd.AddCommand(azureVMListCmd)

	// Storage subcommands
	azureCmd.AddCommand(azureStorageCmd)
	azureStorageCmd.AddCommand(azureStorageListCmd)

	// AKS subcommands
	azureCmd.AddCommand(azureAKSCmd)
	azureAKSCmd.AddCommand(azureAKSListCmd)
	azureAKSCmd.AddCommand(azureAKSCredentialsCmd)

	// Resource Group subcommands
	azureCmd.AddCommand(azureRGCmd)
	azureCmd.AddCommand(configcmd.NewConfigRouter("azure"))
	azureRGCmd.AddCommand(azureRGListCmd)

	// Persistent flags
	azureCmd.PersistentFlags().BoolVar(&azureDryRun, "dry-run", false,
		"Show what would be done without executing")
	azureCmd.PersistentFlags().BoolVarP(&azureVerbose, "verbose", "v", false,
		"Show verbose output")
	azureCmd.PersistentFlags().StringVarP(&azureSubscription, "subscription", "s", "",
		"Azure subscription to use")
}

func newAzureHelper() *azure.Helper {
	helper := azure.NewHelper(azureVerbose, azureDryRun)
	if azureSubscription != "" {
		helper.SetSubscription(azureSubscription)
	}
	return helper
}

func runAzureStatus(cmd *cobra.Command, args []string) error {
	helper := newAzureHelper()
	status, err := helper.GetStatus()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Azure CLI Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if status.Installed {
		fmt.Fprintf(os.Stdout, "%s Azure CLI installed: %s\n", output.Success("✓"), status.Version)
	} else {
		fmt.Fprintf(os.Stdout, "%s Azure CLI not found\n", output.Error("✗"))
		fmt.Fprintln(os.Stdout, "  Install: brew install azure-cli (macOS)")
		return nil
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Authentication:"))
	if status.Authenticated {
		fmt.Fprintf(os.Stdout, "%s Authenticated\n", output.Success("✓"))
		if status.UserName != "" {
			fmt.Fprintf(os.Stdout, "  User: %s\n", status.UserName)
		}
		if status.SubscriptionName != "" {
			fmt.Fprintf(os.Stdout, "  Subscription: %s\n", status.SubscriptionName)
		}
		if status.SubscriptionID != "" {
			fmt.Fprintf(os.Stdout, "  ID: %s\n", status.SubscriptionID)
		}
		if status.TenantID != "" {
			fmt.Fprintf(os.Stdout, "  Tenant: %s\n", status.TenantID)
		}
	} else {
		fmt.Fprintf(os.Stdout, "%s Not authenticated\n", output.Warning("⚠"))
		fmt.Fprintln(os.Stdout, "  Run: az login")
	}

	return nil
}

func runAzureWhoami(cmd *cobra.Command, args []string) error {
	helper := newAzureHelper()
	account, err := helper.GetAccount()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(account)
	}

	fmt.Fprintf(os.Stdout, "Subscription: %s\n", account.SubscriptionName)
	fmt.Fprintf(os.Stdout, "ID:           %s\n", account.SubscriptionID)
	fmt.Fprintf(os.Stdout, "Tenant:       %s\n", account.TenantID)
	fmt.Fprintf(os.Stdout, "User:         %s (%s)\n", account.User.Name, account.User.Type)
	return nil
}

func runAzureSubscriptions(cmd *cobra.Command, args []string) error {
	helper := newAzureHelper()
	subscriptions, err := helper.ListSubscriptions()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(subscriptions)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Azure Subscriptions"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	if len(subscriptions) == 0 {
		fmt.Fprintln(os.Stdout, "No subscriptions found")
		return nil
	}
	for _, s := range subscriptions {
		marker := "  "
		if s.IsDefault {
			marker = "* "
		}
		fmt.Fprintf(os.Stdout, "%s%s (%s)\n", marker, s.Name, s.ID)
	}
	return nil
}

func runAzureSetSubscription(cmd *cobra.Command, args []string) error {
	helper := newAzureHelper()
	if err := helper.SetActiveSubscription(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s Set active subscription to: %s\n", output.Success("✓"), args[0])
	return nil
}

func runAzureOverview(cmd *cobra.Command, args []string) error {
	helper := newAzureHelper()
	overview, err := helper.GetOverview()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(overview)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Azure Overview"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintln(os.Stdout)

	if !overview.Status.Installed {
		fmt.Fprintf(os.Stdout, "%s Azure CLI not installed\n", output.Error("✗"))
		return nil
	}

	if !overview.Status.Authenticated {
		fmt.Fprintf(os.Stdout, "%s Not authenticated. Run: az login\n", output.Warning("⚠"))
		return nil
	}

	if overview.Status.SubscriptionName != "" {
		fmt.Fprintf(os.Stdout, "Subscription: %s\n\n", overview.Status.SubscriptionName)
	}

	fmt.Fprintf(os.Stdout, "Virtual Machines: %d\n", overview.VMCount)

	fmt.Fprintf(os.Stdout, "\n%s:\n", output.Info("Storage Accounts"))
	if len(overview.StorageAccounts) == 0 {
		fmt.Fprintln(os.Stdout, "  None found")
	} else {
		for _, s := range overview.StorageAccounts {
			fmt.Fprintf(os.Stdout, "  %s\n", s)
		}
	}

	fmt.Fprintf(os.Stdout, "\n%s:\n", output.Info("AKS Clusters"))
	if len(overview.AKSClusters) == 0 {
		fmt.Fprintln(os.Stdout, "  None found")
	} else {
		for _, c := range overview.AKSClusters {
			fmt.Fprintf(os.Stdout, "  %s\n", c)
		}
	}

	fmt.Fprintf(os.Stdout, "\n%s:\n", output.Info("Resource Groups"))
	if len(overview.ResourceGroups) == 0 {
		fmt.Fprintln(os.Stdout, "  None found")
	} else {
		for _, g := range overview.ResourceGroups {
			fmt.Fprintf(os.Stdout, "  %s\n", g)
		}
	}

	return nil
}

func runAzureLogin(cmd *cobra.Command, args []string) error {
	helper := newAzureHelper()
	return helper.Login()
}

func runAzureLogout(cmd *cobra.Command, args []string) error {
	helper := newAzureHelper()
	return helper.Logout()
}

func runAzureVMList(cmd *cobra.Command, args []string) error {
	helper := newAzureHelper()

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Azure Virtual Machines"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	vms, err := helper.ListVMs()
	if err != nil {
		fmt.Fprintln(os.Stdout, "No VMs found or not authenticated")
		return nil
	}
	fmt.Fprintln(os.Stdout, vms)
	return nil
}

func runAzureStorageList(cmd *cobra.Command, args []string) error {
	helper := newAzureHelper()
	accounts, err := helper.ListStorageAccounts()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(accounts)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Azure Storage Accounts"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	if len(accounts) == 0 {
		fmt.Fprintln(os.Stdout, "No storage accounts found")
		return nil
	}
	for _, a := range accounts {
		fmt.Fprintf(os.Stdout, "  %s\n", a)
	}
	return nil
}

func runAzureAKSList(cmd *cobra.Command, args []string) error {
	helper := newAzureHelper()
	clusters, err := helper.ListAKSClusters()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(clusters)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Azure AKS Clusters"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	if len(clusters) == 0 {
		fmt.Fprintln(os.Stdout, "No AKS clusters found")
		return nil
	}
	for _, c := range clusters {
		fmt.Fprintf(os.Stdout, "  %s\n", c)
	}
	return nil
}

func runAzureAKSCredentials(cmd *cobra.Command, args []string) error {
	helper := newAzureHelper()
	return helper.GetAKSCredentials(args[0], args[1])
}

func runAzureRGList(cmd *cobra.Command, args []string) error {
	helper := newAzureHelper()
	groups, err := helper.ListResourceGroups()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(groups)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Azure Resource Groups"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	if len(groups) == 0 {
		fmt.Fprintln(os.Stdout, "No resource groups found")
		return nil
	}
	for _, g := range groups {
		fmt.Fprintf(os.Stdout, "  %s\n", g)
	}
	return nil
}

func runAzureInstall(cmd *cobra.Command, args []string) error {
	inst := installer.NewInstaller(
		installer.WithDryRun(azureDryRun),
		installer.WithVerbose(azureVerbose),
	)

	// Show platform info
	platform := inst.GetPlatform()
	if azureVerbose {
		fmt.Fprintf(os.Stdout, "Platform: %s\n\n", platform)
	}

	// Get the plan first
	plan, err := inst.Plan(cmd.Context(), "azure")
	if err != nil {
		return err
	}

	// Show what will be installed
	if azureDryRun {
		fmt.Fprintf(os.Stdout, "%s\n", output.Info("Azure Installation Plan"))
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
		} else if azureDryRun {
			suffix = fmt.Sprintf(" (via %s)", t.Method.Type)
		}
		fmt.Fprintf(os.Stdout, "  %s %s%s\n", status, t.Name, suffix)
	}

	if azureDryRun {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Run without --dry-run to install.")
		return nil
	}

	// Execute installation
	fmt.Fprintln(os.Stdout)
	result, err := inst.Install(cmd.Context(), "azure")
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
		Name: "azure",
		RegisterCmd: func() *cobra.Command { return azureCmd },
	})
}
