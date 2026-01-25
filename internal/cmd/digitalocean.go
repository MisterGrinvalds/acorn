package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/cloud/digitalocean"
	"github.com/mistergrinvalds/acorn/internal/utils/installer"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	doDryRun  bool
	doVerbose bool
	doContext string
)

// doCmd represents the digitalocean command group
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "DigitalOcean CLI helpers",
	Long: `Helpers for DigitalOcean CLI (doctl) operations.

Provides commands for managing DigitalOcean resources across Droplets, Kubernetes, Apps, and Databases.

Examples:
  acorn cloud do status          # Check doctl status and auth
  acorn cloud do whoami          # Show current account
  acorn cloud do droplet list    # List droplets
  acorn cloud do k8s list        # List Kubernetes clusters
  acorn cloud do app list        # List App Platform apps`,
	Aliases: []string{"digitalocean"},
}

// doStatusCmd shows doctl status
var doStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check doctl status and authentication",
	Long: `Check if doctl is installed and authenticated.

Shows doctl version and current account information.

Examples:
  acorn cloud do status
  acorn cloud do status -o json`,
	RunE: runDoStatus,
}

// doWhoamiCmd shows current account
var doWhoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current DigitalOcean account",
	Long: `Display the current DigitalOcean account information.

Examples:
  acorn cloud do whoami`,
	RunE: runDoWhoami,
}

// doContextsCmd lists contexts
var doContextsCmd = &cobra.Command{
	Use:   "contexts",
	Short: "List doctl contexts",
	Long: `List all configured doctl contexts.

Examples:
  acorn cloud do contexts`,
	RunE: runDoContexts,
}

// doOverviewCmd shows overview
var doOverviewCmd = &cobra.Command{
	Use:   "overview",
	Short: "Show overview of DigitalOcean resources",
	Long: `Display an overview of DigitalOcean resources including
Droplets, Kubernetes clusters, Apps, and Databases.

Examples:
  acorn cloud do overview
  acorn cloud do overview -o json`,
	RunE: runDoOverview,
}

// doLoginCmd performs login
var doLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to DigitalOcean",
	Long: `Initiate DigitalOcean authentication.

Examples:
  acorn cloud do login`,
	RunE: runDoLogin,
}

// Droplet subcommands
var doDropletCmd = &cobra.Command{
	Use:   "droplet",
	Short: "Droplet commands",
	Long: `Commands for managing DigitalOcean Droplets.

Examples:
  acorn cloud do droplet list
  acorn cloud do droplet ssh <droplet-id>`,
}

var doDropletListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Droplets",
	Long: `List all Droplets.

Examples:
  acorn cloud do droplet list`,
	RunE: runDoDropletList,
}

var doDropletSSHCmd = &cobra.Command{
	Use:   "ssh <droplet-id>",
	Short: "SSH to Droplet",
	Long: `SSH into a Droplet.

Examples:
  acorn cloud do droplet ssh 12345678`,
	Args: cobra.ExactArgs(1),
	RunE: runDoDropletSSH,
}

// Kubernetes subcommands
var doK8sCmd = &cobra.Command{
	Use:   "k8s",
	Short: "Kubernetes cluster commands",
	Long: `Commands for managing DigitalOcean Kubernetes clusters.

Examples:
  acorn cloud do k8s list
  acorn cloud do k8s credentials my-cluster`,
	Aliases: []string{"kubernetes"},
}

var doK8sListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Kubernetes clusters",
	Long: `List all Kubernetes clusters.

Examples:
  acorn cloud do k8s list`,
	RunE: runDoK8sList,
}

var doK8sCredentialsCmd = &cobra.Command{
	Use:   "credentials <cluster-name>",
	Short: "Get kubeconfig for cluster",
	Long: `Update kubeconfig to use a DigitalOcean Kubernetes cluster.

Examples:
  acorn cloud do k8s credentials my-cluster`,
	Args: cobra.ExactArgs(1),
	RunE: runDoK8sCredentials,
}

// App subcommands
var doAppCmd = &cobra.Command{
	Use:   "app",
	Short: "App Platform commands",
	Long: `Commands for managing DigitalOcean App Platform apps.

Examples:
  acorn cloud do app list`,
}

var doAppListCmd = &cobra.Command{
	Use:   "list",
	Short: "List App Platform apps",
	Long: `List all App Platform apps.

Examples:
  acorn cloud do app list`,
	RunE: runDoAppList,
}

// Database subcommands
var doDBCmd = &cobra.Command{
	Use:   "db",
	Short: "Database commands",
	Long: `Commands for managing DigitalOcean managed databases.

Examples:
  acorn cloud do db list`,
	Aliases: []string{"database"},
}

var doDBListCmd = &cobra.Command{
	Use:   "list",
	Short: "List managed databases",
	Long: `List all managed databases.

Examples:
  acorn cloud do db list`,
	RunE: runDoDBList,
}

// doInstallCmd installs doctl
var doInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install doctl CLI",
	Long: `Install DigitalOcean CLI (doctl).

Automatically detects your platform and uses the appropriate
package manager (brew on macOS, script install on Linux).

Examples:
  acorn cloud do install           # Install doctl
  acorn cloud do install --dry-run # Show what would be installed
  acorn cloud do install -v        # Verbose output`,
	RunE: runDoInstall,
}

func init() {
	cloudCmd.AddCommand(doCmd)

	// Add subcommands
	doCmd.AddCommand(doInstallCmd)
	doCmd.AddCommand(doStatusCmd)
	doCmd.AddCommand(doWhoamiCmd)
	doCmd.AddCommand(doContextsCmd)
	doCmd.AddCommand(doOverviewCmd)
	doCmd.AddCommand(doLoginCmd)

	// Droplet subcommands
	doCmd.AddCommand(doDropletCmd)
	doDropletCmd.AddCommand(doDropletListCmd)
	doDropletCmd.AddCommand(doDropletSSHCmd)

	// Kubernetes subcommands
	doCmd.AddCommand(doK8sCmd)
	doK8sCmd.AddCommand(doK8sListCmd)
	doK8sCmd.AddCommand(doK8sCredentialsCmd)

	// App subcommands
	doCmd.AddCommand(doAppCmd)
	doAppCmd.AddCommand(doAppListCmd)

	// Database subcommands
	doCmd.AddCommand(doDBCmd)
	doDBCmd.AddCommand(doDBListCmd)

	// Persistent flags
	doCmd.PersistentFlags().BoolVar(&doDryRun, "dry-run", false,
		"Show what would be done without executing")
	doCmd.PersistentFlags().BoolVarP(&doVerbose, "verbose", "v", false,
		"Show verbose output")
	doCmd.PersistentFlags().StringVarP(&doContext, "context", "c", "",
		"doctl context to use")
}

func newDoHelper() *digitalocean.Helper {
	helper := digitalocean.NewHelper(doVerbose, doDryRun)
	if doContext != "" {
		helper.SetContext(doContext)
	}
	return helper
}

func runDoStatus(cmd *cobra.Command, args []string) error {
	helper := newDoHelper()
	status, err := helper.GetStatus()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("DigitalOcean CLI Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if status.Installed {
		fmt.Fprintf(os.Stdout, "%s doctl installed: %s\n", output.Success("✓"), status.Version)
	} else {
		fmt.Fprintf(os.Stdout, "%s doctl not found\n", output.Error("✗"))
		fmt.Fprintln(os.Stdout, "  Install: brew install doctl (macOS)")
		return nil
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Authentication:"))
	if status.Authenticated {
		fmt.Fprintf(os.Stdout, "%s Authenticated\n", output.Success("✓"))
		if status.AccountEmail != "" {
			fmt.Fprintf(os.Stdout, "  Email: %s\n", status.AccountEmail)
		}
		if status.AccountUUID != "" {
			fmt.Fprintf(os.Stdout, "  UUID: %s\n", status.AccountUUID)
		}
		if status.TeamName != "" {
			fmt.Fprintf(os.Stdout, "  Team: %s\n", status.TeamName)
		}
	} else {
		fmt.Fprintf(os.Stdout, "%s Not authenticated\n", output.Warning("⚠"))
		fmt.Fprintln(os.Stdout, "  Run: doctl auth init")
	}

	return nil
}

func runDoWhoami(cmd *cobra.Command, args []string) error {
	helper := newDoHelper()
	account, err := helper.GetAccount()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(account)
	}

	fmt.Fprintf(os.Stdout, "Email:            %s\n", account.Email)
	fmt.Fprintf(os.Stdout, "UUID:             %s\n", account.UUID)
	fmt.Fprintf(os.Stdout, "Droplet Limit:    %d\n", account.DropletLimit)
	fmt.Fprintf(os.Stdout, "Floating IP Limit: %d\n", account.FloatingIPLimit)
	if account.Team != "" {
		fmt.Fprintf(os.Stdout, "Team:             %s\n", account.Team)
	}
	return nil
}

func runDoContexts(cmd *cobra.Command, args []string) error {
	helper := newDoHelper()
	contexts, err := helper.ListContexts()
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("doctl Contexts"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintln(os.Stdout, contexts)
	return nil
}

func runDoOverview(cmd *cobra.Command, args []string) error {
	helper := newDoHelper()
	overview, err := helper.GetOverview()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(overview)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("DigitalOcean Overview"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintln(os.Stdout)

	if !overview.Status.Installed {
		fmt.Fprintf(os.Stdout, "%s doctl not installed\n", output.Error("✗"))
		return nil
	}

	if !overview.Status.Authenticated {
		fmt.Fprintf(os.Stdout, "%s Not authenticated. Run: doctl auth init\n", output.Warning("⚠"))
		return nil
	}

	if overview.Status.AccountEmail != "" {
		fmt.Fprintf(os.Stdout, "Account: %s\n\n", overview.Status.AccountEmail)
	}

	fmt.Fprintf(os.Stdout, "Droplets: %d\n", overview.DropletCount)

	fmt.Fprintf(os.Stdout, "\n%s:\n", output.Info("Kubernetes Clusters"))
	if len(overview.K8sClusters) == 0 {
		fmt.Fprintln(os.Stdout, "  None found")
	} else {
		for _, c := range overview.K8sClusters {
			fmt.Fprintf(os.Stdout, "  %s\n", c)
		}
	}

	fmt.Fprintf(os.Stdout, "\n%s:\n", output.Info("App Platform Apps"))
	if len(overview.Apps) == 0 {
		fmt.Fprintln(os.Stdout, "  None found")
	} else {
		for _, a := range overview.Apps {
			fmt.Fprintf(os.Stdout, "  %s\n", a)
		}
	}

	fmt.Fprintf(os.Stdout, "\n%s:\n", output.Info("Managed Databases"))
	if len(overview.Databases) == 0 {
		fmt.Fprintln(os.Stdout, "  None found")
	} else {
		for _, d := range overview.Databases {
			fmt.Fprintf(os.Stdout, "  %s\n", d)
		}
	}

	return nil
}

func runDoLogin(cmd *cobra.Command, args []string) error {
	helper := newDoHelper()
	return helper.Login()
}

func runDoDropletList(cmd *cobra.Command, args []string) error {
	helper := newDoHelper()

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("DigitalOcean Droplets"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	droplets, err := helper.ListDroplets()
	if err != nil {
		fmt.Fprintln(os.Stdout, "No droplets found or not authenticated")
		return nil
	}
	fmt.Fprintln(os.Stdout, droplets)
	return nil
}

func runDoDropletSSH(cmd *cobra.Command, args []string) error {
	helper := newDoHelper()
	return helper.SSHToDroplet(args[0])
}

func runDoK8sList(cmd *cobra.Command, args []string) error {
	helper := newDoHelper()
	clusters, err := helper.ListK8sClusters()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(clusters)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("DigitalOcean Kubernetes Clusters"))
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

func runDoK8sCredentials(cmd *cobra.Command, args []string) error {
	helper := newDoHelper()
	return helper.GetK8sCredentials(args[0])
}

func runDoAppList(cmd *cobra.Command, args []string) error {
	helper := newDoHelper()
	apps, err := helper.ListApps()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(apps)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("DigitalOcean App Platform Apps"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	if len(apps) == 0 {
		fmt.Fprintln(os.Stdout, "No apps found")
		return nil
	}
	for _, a := range apps {
		fmt.Fprintf(os.Stdout, "  %s\n", a)
	}
	return nil
}

func runDoDBList(cmd *cobra.Command, args []string) error {
	helper := newDoHelper()
	databases, err := helper.ListDatabases()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(databases)
	}

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("DigitalOcean Managed Databases"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	if len(databases) == 0 {
		fmt.Fprintln(os.Stdout, "No databases found")
		return nil
	}
	for _, d := range databases {
		fmt.Fprintf(os.Stdout, "  %s\n", d)
	}
	return nil
}

func runDoInstall(cmd *cobra.Command, args []string) error {
	inst := installer.NewInstaller(
		installer.WithDryRun(doDryRun),
		installer.WithVerbose(doVerbose),
	)

	// Show platform info
	platform := inst.GetPlatform()
	if doVerbose {
		fmt.Fprintf(os.Stdout, "Platform: %s\n\n", platform)
	}

	// Get the plan first
	plan, err := inst.Plan(cmd.Context(), "digitalocean")
	if err != nil {
		return err
	}

	// Show what will be installed
	if doDryRun {
		fmt.Fprintf(os.Stdout, "%s\n", output.Info("DigitalOcean Installation Plan"))
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
		} else if doDryRun {
			suffix = fmt.Sprintf(" (via %s)", t.Method.Type)
		}
		fmt.Fprintf(os.Stdout, "  %s %s%s\n", status, t.Name, suffix)
	}

	if doDryRun {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Run without --dry-run to install.")
		return nil
	}

	// Execute installation
	fmt.Fprintln(os.Stdout)
	result, err := inst.Install(cmd.Context(), "digitalocean")
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
