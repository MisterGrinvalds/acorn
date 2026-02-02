package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/cloud/cloudflare"
	"github.com/mistergrinvalds/acorn/internal/utils/installer"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/mistergrinvalds/acorn/internal/utils/configcmd"
	"github.com/spf13/cobra"
)

var (
	cfDryRun  bool
	cfVerbose bool
)

// cfCmd represents the cloudflare command group
var cfCmd = &cobra.Command{
	Use:   "cf",
	Short: "CloudFlare CLI helpers",
	Long: `Helpers for CloudFlare Workers, Pages, R2, KV, and D1.

Provides commands for managing CloudFlare resources using wrangler.

Examples:
  acorn cf status              # Check wrangler status and auth
  acorn cf workers             # List Workers deployments
  acorn cf pages               # List Pages projects
  acorn cf r2 list             # List R2 buckets
  acorn cf kv list             # List KV namespaces
  acorn cf d1 list             # List D1 databases`,
	Aliases: []string{"cloudflare"},
}

// cfStatusCmd shows CloudFlare CLI status
var cfStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check CloudFlare CLI status and authentication",
	Long: `Check if wrangler is installed and authenticated.

Shows wrangler version and current account information.

Examples:
  acorn cf status
  acorn cf status -o json`,
	RunE: runCfStatus,
}

// cfWhoamiCmd shows current account
var cfWhoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current CloudFlare account",
	Long: `Display the current CloudFlare account information.

Examples:
  acorn cf whoami`,
	RunE: runCfWhoami,
}

// cfWorkersCmd lists Workers
var cfWorkersCmd = &cobra.Command{
	Use:   "workers",
	Short: "List CloudFlare Workers deployments",
	Long: `List all Workers deployments in your account.

Examples:
  acorn cf workers`,
	RunE: runCfWorkers,
}

// cfPagesCmd lists Pages projects
var cfPagesCmd = &cobra.Command{
	Use:   "pages",
	Short: "List CloudFlare Pages projects",
	Long: `List all Pages projects in your account.

Examples:
  acorn cf pages`,
	RunE: runCfPages,
}

// cfLogsCmd tails worker logs
var cfLogsCmd = &cobra.Command{
	Use:   "logs <worker-name>",
	Short: "Tail worker logs",
	Long: `Tail logs for a CloudFlare Worker.

Examples:
  acorn cf logs my-worker`,
	Args: cobra.ExactArgs(1),
	RunE: runCfLogs,
}

// cfDeployCmd deploys current worker
var cfDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy current worker",
	Long: `Deploy the worker in the current directory.

Requires a wrangler.toml or wrangler.json configuration file.

Examples:
  acorn cf deploy`,
	RunE: runCfDeploy,
}

// cfSecretsCmd lists worker secrets
var cfSecretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "List worker secrets",
	Long: `List secrets for the current worker.

Examples:
  acorn cf secrets`,
	RunE: runCfSecrets,
}

// cfSecretPutCmd adds a secret
var cfSecretPutCmd = &cobra.Command{
	Use:   "secret-put <name>",
	Short: "Add a worker secret",
	Long: `Add a secret to the current worker.

You will be prompted to enter the secret value.

Examples:
  acorn cf secret-put MY_API_KEY`,
	Args: cobra.ExactArgs(1),
	RunE: runCfSecretPut,
}

// cfOverviewCmd shows all resources
var cfOverviewCmd = &cobra.Command{
	Use:   "overview",
	Short: "Show overview of all CloudFlare resources",
	Long: `Display an overview of all CloudFlare resources including
Workers, Pages, R2, KV, and D1.

Examples:
  acorn cf overview
  acorn cf overview -o json`,
	RunE: runCfOverview,
}

// cfLoginCmd initiates login
var cfLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to CloudFlare",
	Long: `Initiate CloudFlare login via wrangler.

Examples:
  acorn cf login`,
	RunE: runCfLogin,
}

// cfLogoutCmd logs out
var cfLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from CloudFlare",
	Long: `Logout from CloudFlare.

Examples:
  acorn cf logout`,
	RunE: runCfLogout,
}

// R2 subcommands
var cfR2Cmd = &cobra.Command{
	Use:   "r2",
	Short: "R2 storage commands",
	Long: `Commands for managing CloudFlare R2 storage buckets.

Examples:
  acorn cf r2 list
  acorn cf r2 create my-bucket`,
}

var cfR2ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List R2 buckets",
	Long: `List all R2 buckets in your account.

Examples:
  acorn cf r2 list`,
	RunE: runCfR2List,
}

var cfR2CreateCmd = &cobra.Command{
	Use:   "create <bucket-name>",
	Short: "Create R2 bucket",
	Long: `Create a new R2 bucket.

Examples:
  acorn cf r2 create my-bucket`,
	Args: cobra.ExactArgs(1),
	RunE: runCfR2Create,
}

// KV subcommands
var cfKVCmd = &cobra.Command{
	Use:   "kv",
	Short: "KV storage commands",
	Long: `Commands for managing CloudFlare KV namespaces.

Examples:
  acorn cf kv list
  acorn cf kv create my-namespace`,
}

var cfKVListCmd = &cobra.Command{
	Use:   "list",
	Short: "List KV namespaces",
	Long: `List all KV namespaces in your account.

Examples:
  acorn cf kv list`,
	RunE: runCfKVList,
}

var cfKVCreateCmd = &cobra.Command{
	Use:   "create <namespace-name>",
	Short: "Create KV namespace",
	Long: `Create a new KV namespace.

Examples:
  acorn cf kv create my-namespace`,
	Args: cobra.ExactArgs(1),
	RunE: runCfKVCreate,
}

// D1 subcommands
var cfD1Cmd = &cobra.Command{
	Use:   "d1",
	Short: "D1 database commands",
	Long: `Commands for managing CloudFlare D1 databases.

Examples:
  acorn cf d1 list
  acorn cf d1 create my-database`,
}

var cfD1ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List D1 databases",
	Long: `List all D1 databases in your account.

Examples:
  acorn cf d1 list`,
	RunE: runCfD1List,
}

var cfD1CreateCmd = &cobra.Command{
	Use:   "create <database-name>",
	Short: "Create D1 database",
	Long: `Create a new D1 database.

Examples:
  acorn cf d1 create my-database`,
	Args: cobra.ExactArgs(1),
	RunE: runCfD1Create,
}

// Init subcommands
var cfInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize new projects",
	Long: `Initialize new CloudFlare Worker or Pages projects.

Examples:
  acorn cf init worker my-worker
  acorn cf init pages my-site`,
}

var cfInitWorkerCmd = &cobra.Command{
	Use:   "worker [name]",
	Short: "Initialize new Worker project",
	Long: `Initialize a new CloudFlare Worker project.

Examples:
  acorn cf init worker my-worker
  acorn cf init worker`,
	Args: cobra.MaximumNArgs(1),
	RunE: runCfInitWorker,
}

var cfInitPagesCmd = &cobra.Command{
	Use:   "pages [name]",
	Short: "Initialize new Pages project",
	Long: `Initialize a new CloudFlare Pages project.

Examples:
  acorn cf init pages my-site
  acorn cf init pages`,
	Args: cobra.MaximumNArgs(1),
	RunE: runCfInitPages,
}

// cfInstallCmd installs CloudFlare CLI tools
var cfInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install CloudFlare CLI tools",
	Long: `Install wrangler and other CloudFlare CLI tools.

Automatically detects your platform and uses the appropriate
package manager (brew on macOS, apt on Linux, npm as fallback).

If prerequisites like npm are missing, they will be installed
automatically if possible.

Examples:
  acorn cf install           # Install all CloudFlare tools
  acorn cf install --dry-run # Show what would be installed
  acorn cf install -v        # Verbose output`,
	RunE: runCfInstall,
}

func init() {
	cloudCmd.AddCommand(cfCmd)

	// Add subcommands
	cfCmd.AddCommand(cfInstallCmd)
	cfCmd.AddCommand(cfStatusCmd)
	cfCmd.AddCommand(cfWhoamiCmd)
	cfCmd.AddCommand(cfWorkersCmd)
	cfCmd.AddCommand(cfPagesCmd)
	cfCmd.AddCommand(cfLogsCmd)
	cfCmd.AddCommand(cfDeployCmd)
	cfCmd.AddCommand(cfSecretsCmd)
	cfCmd.AddCommand(cfSecretPutCmd)
	cfCmd.AddCommand(cfOverviewCmd)
	cfCmd.AddCommand(cfLoginCmd)
	cfCmd.AddCommand(cfLogoutCmd)

	// R2 subcommands
	cfCmd.AddCommand(cfR2Cmd)
	cfR2Cmd.AddCommand(cfR2ListCmd)
	cfR2Cmd.AddCommand(cfR2CreateCmd)

	// KV subcommands
	cfCmd.AddCommand(cfKVCmd)
	cfKVCmd.AddCommand(cfKVListCmd)
	cfKVCmd.AddCommand(cfKVCreateCmd)

	// D1 subcommands
	cfCmd.AddCommand(cfD1Cmd)
	cfD1Cmd.AddCommand(cfD1ListCmd)
	cfD1Cmd.AddCommand(cfD1CreateCmd)

	// Init subcommands
	cfCmd.AddCommand(cfInitCmd)
	cfCmd.AddCommand(configcmd.NewConfigRouter("cloudflare"))
	cfInitCmd.AddCommand(cfInitWorkerCmd)
	cfInitCmd.AddCommand(cfInitPagesCmd)

	// Persistent flags
	cfCmd.PersistentFlags().BoolVar(&cfDryRun, "dry-run", false,
		"Show what would be done without executing")
	cfCmd.PersistentFlags().BoolVarP(&cfVerbose, "verbose", "v", false,
		"Show verbose output")
}

func runCfStatus(cmd *cobra.Command, args []string) error {
	helper := cloudflare.NewHelper(cfVerbose, cfDryRun)
	status, err := helper.GetStatus()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("CloudFlare CLI Status"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if status.Installed {
		fmt.Fprintf(os.Stdout, "%s wrangler installed: %s\n", output.Success("✓"), status.Version)
	} else {
		fmt.Fprintf(os.Stdout, "%s wrangler not found\n", output.Error("✗"))
		fmt.Fprintln(os.Stdout, "  Install: npm install -g wrangler")
		return nil
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("Authentication:"))
	if status.Authenticated {
		fmt.Fprintf(os.Stdout, "%s Logged in\n", output.Success("✓"))
		if status.AccountName != "" {
			fmt.Fprintf(os.Stdout, "  Account: %s\n", status.AccountName)
		}
		if status.AccountID != "" {
			fmt.Fprintf(os.Stdout, "  ID: %s\n", status.AccountID)
		}
	} else {
		fmt.Fprintf(os.Stdout, "%s Not logged in\n", output.Warning("⚠"))
		fmt.Fprintln(os.Stdout, "  Run: acorn cf login")
	}

	return nil
}

func runCfWhoami(cmd *cobra.Command, args []string) error {
	helper := cloudflare.NewHelper(cfVerbose, cfDryRun)
	info, err := helper.Whoami()
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, info)
	return nil
}

func runCfWorkers(cmd *cobra.Command, args []string) error {
	helper := cloudflare.NewHelper(cfVerbose, cfDryRun)

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("CloudFlare Workers"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	workers, err := helper.ListWorkers()
	if err != nil {
		fmt.Fprintln(os.Stdout, "No workers found or not authenticated")
		return nil
	}
	fmt.Fprintln(os.Stdout, workers)
	return nil
}

func runCfPages(cmd *cobra.Command, args []string) error {
	helper := cloudflare.NewHelper(cfVerbose, cfDryRun)

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("CloudFlare Pages Projects"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	pages, err := helper.ListPages()
	if err != nil {
		fmt.Fprintln(os.Stdout, "No pages projects found or not authenticated")
		return nil
	}
	fmt.Fprintln(os.Stdout, pages)
	return nil
}

func runCfLogs(cmd *cobra.Command, args []string) error {
	helper := cloudflare.NewHelper(cfVerbose, cfDryRun)
	return helper.TailLogs(args[0])
}

func runCfDeploy(cmd *cobra.Command, args []string) error {
	helper := cloudflare.NewHelper(cfVerbose, cfDryRun)
	return helper.Deploy(args...)
}

func runCfSecrets(cmd *cobra.Command, args []string) error {
	helper := cloudflare.NewHelper(cfVerbose, cfDryRun)
	secrets, err := helper.ListSecrets()
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, secrets)
	return nil
}

func runCfSecretPut(cmd *cobra.Command, args []string) error {
	helper := cloudflare.NewHelper(cfVerbose, cfDryRun)
	return helper.PutSecret(args[0])
}

func runCfOverview(cmd *cobra.Command, args []string) error {
	helper := cloudflare.NewHelper(cfVerbose, cfDryRun)
	overview, err := helper.GetOverview()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(overview)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n", output.Info("CloudFlare Overview"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintln(os.Stdout)

	if !overview.Status.Installed {
		fmt.Fprintf(os.Stdout, "%s wrangler not installed\n", output.Error("✗"))
		return nil
	}

	if !overview.Status.Authenticated {
		fmt.Fprintf(os.Stdout, "%s Not authenticated. Run: acorn cf login\n", output.Warning("⚠"))
		return nil
	}

	if overview.Status.AccountName != "" {
		fmt.Fprintf(os.Stdout, "Account: %s\n\n", overview.Status.AccountName)
	}

	printResourceList("Workers", overview.Workers)
	printResourceList("Pages", overview.Pages)
	printResourceList("R2 Buckets", overview.R2Buckets)
	printResourceList("KV Namespaces", overview.KV)
	printResourceList("D1 Databases", overview.D1)

	return nil
}

func printResourceList(name string, items []string) {
	fmt.Fprintf(os.Stdout, "%s:\n", output.Info(name))
	if len(items) == 0 {
		fmt.Fprintln(os.Stdout, "  None found")
	} else {
		for _, item := range items {
			fmt.Fprintf(os.Stdout, "  %s\n", item)
		}
	}
	fmt.Fprintln(os.Stdout)
}

func runCfLogin(cmd *cobra.Command, args []string) error {
	helper := cloudflare.NewHelper(cfVerbose, cfDryRun)
	return helper.Login()
}

func runCfLogout(cmd *cobra.Command, args []string) error {
	helper := cloudflare.NewHelper(cfVerbose, cfDryRun)
	return helper.Logout()
}

func runCfR2List(cmd *cobra.Command, args []string) error {
	helper := cloudflare.NewHelper(cfVerbose, cfDryRun)

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("CloudFlare R2 Buckets"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	buckets, err := helper.ListR2Buckets()
	if err != nil {
		fmt.Fprintln(os.Stdout, "No R2 buckets found or not authenticated")
		return nil
	}
	fmt.Fprintln(os.Stdout, buckets)
	return nil
}

func runCfR2Create(cmd *cobra.Command, args []string) error {
	helper := cloudflare.NewHelper(cfVerbose, cfDryRun)
	if err := helper.CreateR2Bucket(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s R2 bucket '%s' created\n", output.Success("✓"), args[0])
	return nil
}

func runCfKVList(cmd *cobra.Command, args []string) error {
	helper := cloudflare.NewHelper(cfVerbose, cfDryRun)

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("CloudFlare KV Namespaces"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	namespaces, err := helper.ListKVNamespaces()
	if err != nil {
		fmt.Fprintln(os.Stdout, "No KV namespaces found or not authenticated")
		return nil
	}
	fmt.Fprintln(os.Stdout, namespaces)
	return nil
}

func runCfKVCreate(cmd *cobra.Command, args []string) error {
	helper := cloudflare.NewHelper(cfVerbose, cfDryRun)
	if err := helper.CreateKVNamespace(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s KV namespace '%s' created\n", output.Success("✓"), args[0])
	return nil
}

func runCfD1List(cmd *cobra.Command, args []string) error {
	helper := cloudflare.NewHelper(cfVerbose, cfDryRun)

	fmt.Fprintf(os.Stdout, "%s\n", output.Info("CloudFlare D1 Databases"))
	fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	databases, err := helper.ListD1Databases()
	if err != nil {
		fmt.Fprintln(os.Stdout, "No D1 databases found or not authenticated")
		return nil
	}
	fmt.Fprintln(os.Stdout, databases)
	return nil
}

func runCfD1Create(cmd *cobra.Command, args []string) error {
	helper := cloudflare.NewHelper(cfVerbose, cfDryRun)
	if err := helper.CreateD1Database(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s D1 database '%s' created\n", output.Success("✓"), args[0])
	return nil
}

func runCfInitWorker(cmd *cobra.Command, args []string) error {
	helper := cloudflare.NewHelper(cfVerbose, cfDryRun)
	name := ""
	if len(args) > 0 {
		name = args[0]
	}
	return helper.InitWorker(name)
}

func runCfInitPages(cmd *cobra.Command, args []string) error {
	helper := cloudflare.NewHelper(cfVerbose, cfDryRun)
	name := ""
	if len(args) > 0 {
		name = args[0]
	}
	return helper.InitPages(name)
}

func runCfInstall(cmd *cobra.Command, args []string) error {
	inst := installer.NewInstaller(
		installer.WithDryRun(cfDryRun),
		installer.WithVerbose(cfVerbose),
	)

	// Show platform info
	platform := inst.GetPlatform()
	if cfVerbose {
		fmt.Fprintf(os.Stdout, "Platform: %s\n\n", platform)
	}

	// Get the plan first
	plan, err := inst.Plan(cmd.Context(), "cloudflare")
	if err != nil {
		return err
	}

	// Show what will be installed
	if cfDryRun {
		fmt.Fprintf(os.Stdout, "%s\n", output.Info("CloudFlare Installation Plan"))
		fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Fprintf(os.Stdout, "Platform: %s\n\n", platform)
	}

	pending := plan.PendingTools()
	if len(pending) == 0 {
		fmt.Fprintf(os.Stdout, "%s All tools already installed\n", output.Success("✓"))
		return nil
	}

	// Show prerequisites
	if len(plan.Prerequisites) > 0 {
		fmt.Fprintln(os.Stdout, "Prerequisites:")
		for _, t := range plan.Prerequisites {
			status := output.Warning("○")
			suffix := ""
			if t.AlreadyInstalled {
				status = output.Success("✓")
				suffix = " (installed)"
			}
			fmt.Fprintf(os.Stdout, "  %s %s%s\n", status, t.Name, suffix)
		}
		fmt.Fprintln(os.Stdout)
	}

	// Show tools
	fmt.Fprintln(os.Stdout, "Tools:")
	for _, t := range plan.Tools {
		status := output.Warning("○")
		suffix := ""
		if t.AlreadyInstalled {
			status = output.Success("✓")
			suffix = " (installed)"
		} else if cfDryRun {
			suffix = fmt.Sprintf(" (via %s)", t.Method.Type)
		}
		fmt.Fprintf(os.Stdout, "  %s %s%s\n", status, t.Name, suffix)
	}

	if cfDryRun {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Run without --dry-run to install.")
		return nil
	}

	// Execute installation
	fmt.Fprintln(os.Stdout)
	result, err := inst.Install(cmd.Context(), "cloudflare")
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
