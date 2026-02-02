package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/devops/helm"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/mistergrinvalds/acorn/internal/utils/configcmd"
	"github.com/spf13/cobra"
)

var (
	helmVerbose       bool
	helmDryRun        bool
	helmNamespace     string
	helmAllNamespaces bool
	helmWait          bool
	helmAtomic        bool
	helmInstall       bool
	helmAllValues     bool
	helmStrict        bool
	helmValues        []string
)

// helmCmd represents the helm command group
var helmCmd = &cobra.Command{
	Use:   "helm",
	Short: "Helm package manager commands",
	Long: `Helm package manager commands for Kubernetes.

Provides release management, repository operations, and chart development.

Examples:
  acorn helm status           # Show Helm status
  acorn helm releases         # List releases
  acorn helm repos            # List repositories
  acorn helm search nginx     # Search charts
  acorn helm install myapp ./chart  # Install chart`,
	Aliases: []string{"hm"},
}

// helmStatusCmd shows Helm status
var helmStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Helm installation status",
	Long: `Display Helm installation status and version.

Examples:
  acorn helm status
  acorn helm status -o json`,
	RunE: runHelmStatus,
}

// helmReleasesCmd lists releases
var helmReleasesCmd = &cobra.Command{
	Use:   "releases",
	Short: "List Helm releases",
	Long: `List all Helm releases.

Examples:
  acorn helm releases
  acorn helm releases -A
  acorn helm releases -n kube-system`,
	Aliases: []string{"ls", "list"},
	RunE:    runHelmReleases,
}

// helmGetCmd gets release information
var helmGetCmd = &cobra.Command{
	Use:   "get [release]",
	Short: "Get release information",
	Long: `Get detailed information about a release.

Examples:
  acorn helm get myapp
  acorn helm get myapp -n production`,
	Args: cobra.ExactArgs(1),
	RunE: runHelmGet,
}

// helmValuesCmd gets release values
var helmValuesCmd = &cobra.Command{
	Use:   "values [release]",
	Short: "Get release values",
	Long: `Get the values for a release.

Examples:
  acorn helm values myapp
  acorn helm values myapp --all`,
	Args: cobra.ExactArgs(1),
	RunE: runHelmValues,
}

// helmHistoryCmd shows release history
var helmHistoryCmd = &cobra.Command{
	Use:   "history [release]",
	Short: "Show release history",
	Long: `Show the revision history of a release.

Examples:
  acorn helm history myapp`,
	Args: cobra.ExactArgs(1),
	RunE: runHelmHistory,
}

// helmInstallCmd installs a chart
var helmInstallCmd = &cobra.Command{
	Use:   "install [release] [chart]",
	Short: "Install a chart",
	Long: `Install a Helm chart.

Examples:
  acorn helm install myapp ./chart
  acorn helm install myapp bitnami/nginx
  acorn helm install myapp ./chart -f values.yaml --wait`,
	Args: cobra.ExactArgs(2),
	RunE: runHelmInstall,
}

// helmUpgradeCmd upgrades a release
var helmUpgradeCmd = &cobra.Command{
	Use:   "upgrade [release] [chart]",
	Short: "Upgrade a release",
	Long: `Upgrade a Helm release.

Examples:
  acorn helm upgrade myapp ./chart
  acorn helm upgrade myapp ./chart --install
  acorn helm upgrade myapp ./chart --atomic`,
	Args: cobra.ExactArgs(2),
	RunE: runHelmUpgrade,
}

// helmUninstallCmd uninstalls a release
var helmUninstallCmd = &cobra.Command{
	Use:   "uninstall [release]",
	Short: "Uninstall a release",
	Long: `Uninstall a Helm release.

Examples:
  acorn helm uninstall myapp
  acorn helm uninstall myapp -n production`,
	Aliases: []string{"delete", "del"},
	Args:    cobra.ExactArgs(1),
	RunE:    runHelmUninstall,
}

// helmRollbackCmd rolls back a release
var helmRollbackCmd = &cobra.Command{
	Use:   "rollback [release] [revision]",
	Short: "Rollback a release",
	Long: `Rollback a release to a previous revision.

Examples:
  acorn helm rollback myapp 1
  acorn helm rollback myapp 2 -n production`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runHelmRollback,
}

// helmReposCmd lists repositories
var helmReposCmd = &cobra.Command{
	Use:   "repos",
	Short: "List repositories",
	Long: `List configured Helm repositories.

Examples:
  acorn helm repos`,
	Aliases: []string{"repo"},
	RunE:    runHelmRepos,
}

// helmRepoAddCmd adds a repository
var helmRepoAddCmd = &cobra.Command{
	Use:   "repo-add [name] [url]",
	Short: "Add a repository",
	Long: `Add a Helm repository.

Examples:
  acorn helm repo-add bitnami https://charts.bitnami.com/bitnami`,
	Args: cobra.ExactArgs(2),
	RunE: runHelmRepoAdd,
}

// helmRepoUpdateCmd updates repositories
var helmRepoUpdateCmd = &cobra.Command{
	Use:   "repo-update",
	Short: "Update repositories",
	Long: `Update all Helm repositories.

Examples:
  acorn helm repo-update`,
	Aliases: []string{"update"},
	RunE:    runHelmRepoUpdate,
}

// helmSearchCmd searches for charts
var helmSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for charts",
	Long: `Search for charts in repositories.

Examples:
  acorn helm search nginx
  acorn helm search ""  # list all`,
	Args: cobra.MaximumNArgs(1),
	RunE: runHelmSearch,
}

// helmShowCmd shows chart info
var helmShowCmd = &cobra.Command{
	Use:   "show [chart]",
	Short: "Show chart information",
	Long: `Show information about a chart.

Examples:
  acorn helm show bitnami/nginx`,
	Args: cobra.ExactArgs(1),
	RunE: runHelmShow,
}

// helmTemplateCmd renders chart templates
var helmTemplateCmd = &cobra.Command{
	Use:   "template [release] [chart]",
	Short: "Render chart templates",
	Long: `Render chart templates locally without installing.

Examples:
  acorn helm template myapp ./chart
  acorn helm template myapp ./chart -f values.yaml`,
	Args: cobra.ExactArgs(2),
	RunE: runHelmTemplate,
}

// helmLintCmd lints a chart
var helmLintCmd = &cobra.Command{
	Use:   "lint [chart]",
	Short: "Lint a chart",
	Long: `Lint a Helm chart for issues.

Examples:
  acorn helm lint ./chart
  acorn helm lint ./chart --strict`,
	Args: cobra.MaximumNArgs(1),
	RunE: runHelmLint,
}

// helmCreateCmd creates a new chart
var helmCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new chart",
	Long: `Create a new Helm chart with the given name.

Examples:
  acorn helm create myapp`,
	Args: cobra.ExactArgs(1),
	RunE: runHelmCreate,
}

// helmPluginsCmd lists plugins
var helmPluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "List plugins",
	Long: `List installed Helm plugins.

Examples:
  acorn helm plugins`,
	RunE: runHelmPlugins,
}

func init() {
	devopsCmd.AddCommand(helmCmd)

	// Add subcommands
	helmCmd.AddCommand(helmStatusCmd)
	helmCmd.AddCommand(helmReleasesCmd)
	helmCmd.AddCommand(helmGetCmd)
	helmCmd.AddCommand(helmValuesCmd)
	helmCmd.AddCommand(helmHistoryCmd)
	helmCmd.AddCommand(helmInstallCmd)
	helmCmd.AddCommand(helmUpgradeCmd)
	helmCmd.AddCommand(helmUninstallCmd)
	helmCmd.AddCommand(helmRollbackCmd)
	helmCmd.AddCommand(helmReposCmd)
	helmCmd.AddCommand(helmRepoAddCmd)
	helmCmd.AddCommand(helmRepoUpdateCmd)
	helmCmd.AddCommand(helmSearchCmd)
	helmCmd.AddCommand(helmShowCmd)
	helmCmd.AddCommand(helmTemplateCmd)
	helmCmd.AddCommand(helmLintCmd)
	helmCmd.AddCommand(helmCreateCmd)
	helmCmd.AddCommand(helmPluginsCmd)
	helmCmd.AddCommand(configcmd.NewConfigRouter("helm"))

	// Persistent flags (output format is inherited from root command)
	helmCmd.PersistentFlags().BoolVarP(&helmVerbose, "verbose", "v", false,
		"Show verbose output")
	helmCmd.PersistentFlags().BoolVar(&helmDryRun, "dry-run", false,
		"Show what would be done without executing")
	helmCmd.PersistentFlags().StringVarP(&helmNamespace, "namespace", "n", "",
		"Kubernetes namespace")
	helmCmd.PersistentFlags().BoolVarP(&helmAllNamespaces, "all-namespaces", "A", false,
		"All namespaces")

	// Command-specific flags
	helmValuesCmd.Flags().BoolVar(&helmAllValues, "all", false, "Show all values including defaults")
	helmInstallCmd.Flags().StringArrayVarP(&helmValues, "values", "f", nil, "Values files")
	helmInstallCmd.Flags().BoolVar(&helmWait, "wait", false, "Wait for resources to be ready")
	helmUpgradeCmd.Flags().StringArrayVarP(&helmValues, "values", "f", nil, "Values files")
	helmUpgradeCmd.Flags().BoolVar(&helmWait, "wait", false, "Wait for resources to be ready")
	helmUpgradeCmd.Flags().BoolVar(&helmAtomic, "atomic", false, "Rollback on failure")
	helmUpgradeCmd.Flags().BoolVar(&helmInstall, "install", false, "Install if not exists")
	helmTemplateCmd.Flags().StringArrayVarP(&helmValues, "values", "f", nil, "Values files")
	helmLintCmd.Flags().BoolVar(&helmStrict, "strict", false, "Strict mode")
}

func runHelmStatus(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := helm.NewHelper(helmVerbose, helmDryRun)
	status := helper.GetStatus()

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	fmt.Fprintf(os.Stdout, "Installed: %v\n", status.Installed)
	if status.Version != "" {
		fmt.Fprintf(os.Stdout, "Version:   %s\n", status.Version)
	}
	if status.Location != "" {
		fmt.Fprintf(os.Stdout, "Location:  %s\n", status.Location)
	}

	return nil
}

func runHelmReleases(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := helm.NewHelper(helmVerbose, helmDryRun)

	if !helper.IsInstalled() {
		return fmt.Errorf("helm is not installed")
	}

	releases, err := helper.ListReleases(helmNamespace, helmAllNamespaces)
	if err != nil {
		return err
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(map[string]any{"releases": releases})
	}

	if len(releases) == 0 {
		fmt.Fprintln(os.Stdout, "No releases found")
		return nil
	}

	fmt.Fprintf(os.Stdout, "%-20s %-15s %-8s %-12s %-30s\n", "NAME", "NAMESPACE", "REVISION", "STATUS", "CHART")
	for _, r := range releases {
		fmt.Fprintf(os.Stdout, "%-20s %-15s %-8s %-12s %-30s\n",
			truncateStr(r.Name, 20),
			truncateStr(r.Namespace, 15),
			r.Revision,
			r.Status,
			truncateStr(r.Chart, 30))
	}

	return nil
}

func runHelmGet(cmd *cobra.Command, args []string) error {
	helper := helm.NewHelper(helmVerbose, helmDryRun)

	if !helper.IsInstalled() {
		return fmt.Errorf("helm is not installed")
	}

	return helper.GetReleaseStatus(args[0], helmNamespace)
}

func runHelmValues(cmd *cobra.Command, args []string) error {
	helper := helm.NewHelper(helmVerbose, helmDryRun)

	if !helper.IsInstalled() {
		return fmt.Errorf("helm is not installed")
	}

	values, err := helper.GetReleaseValues(args[0], helmNamespace, helmAllValues)
	if err != nil {
		return err
	}

	fmt.Print(values)
	return nil
}

func runHelmHistory(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := helm.NewHelper(helmVerbose, helmDryRun)

	if !helper.IsInstalled() {
		return fmt.Errorf("helm is not installed")
	}

	history, err := helper.GetReleaseHistory(args[0], helmNamespace)
	if err != nil {
		return err
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(map[string]any{"history": history})
	}

	fmt.Fprintf(os.Stdout, "%-10s %-25s %-12s %-30s\n", "REVISION", "UPDATED", "STATUS", "CHART")
	for _, h := range history {
		fmt.Fprintf(os.Stdout, "%-10d %-25s %-12s %-30s\n",
			h.Revision,
			truncateStr(h.Updated, 25),
			h.Status,
			truncateStr(h.Chart, 30))
	}

	return nil
}

func runHelmInstall(cmd *cobra.Command, args []string) error {
	helper := helm.NewHelper(helmVerbose, helmDryRun)

	if !helper.IsInstalled() {
		return fmt.Errorf("helm is not installed")
	}

	return helper.Install(args[0], args[1], helmNamespace, helmValues, helmWait)
}

func runHelmUpgrade(cmd *cobra.Command, args []string) error {
	helper := helm.NewHelper(helmVerbose, helmDryRun)

	if !helper.IsInstalled() {
		return fmt.Errorf("helm is not installed")
	}

	return helper.Upgrade(args[0], args[1], helmNamespace, helmValues, helmInstall, helmWait, helmAtomic)
}

func runHelmUninstall(cmd *cobra.Command, args []string) error {
	helper := helm.NewHelper(helmVerbose, helmDryRun)

	if !helper.IsInstalled() {
		return fmt.Errorf("helm is not installed")
	}

	if err := helper.Uninstall(args[0], helmNamespace); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Uninstalled release %s\n", output.Success("✓"), args[0])
	return nil
}

func runHelmRollback(cmd *cobra.Command, args []string) error {
	helper := helm.NewHelper(helmVerbose, helmDryRun)

	if !helper.IsInstalled() {
		return fmt.Errorf("helm is not installed")
	}

	revision := 0
	if len(args) > 1 {
		fmt.Sscanf(args[1], "%d", &revision)
	}

	return helper.Rollback(args[0], helmNamespace, revision)
}

func runHelmRepos(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := helm.NewHelper(helmVerbose, helmDryRun)

	if !helper.IsInstalled() {
		return fmt.Errorf("helm is not installed")
	}

	repos, err := helper.ListRepositories()
	if err != nil {
		return err
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(map[string]any{"repositories": repos})
	}

	if len(repos) == 0 {
		fmt.Fprintln(os.Stdout, "No repositories found")
		return nil
	}

	fmt.Fprintf(os.Stdout, "%-20s %s\n", "NAME", "URL")
	for _, r := range repos {
		fmt.Fprintf(os.Stdout, "%-20s %s\n", r.Name, r.URL)
	}

	return nil
}

func runHelmRepoAdd(cmd *cobra.Command, args []string) error {
	helper := helm.NewHelper(helmVerbose, helmDryRun)

	if !helper.IsInstalled() {
		return fmt.Errorf("helm is not installed")
	}

	if err := helper.AddRepository(args[0], args[1]); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Added repository %s\n", output.Success("✓"), args[0])
	return nil
}

func runHelmRepoUpdate(cmd *cobra.Command, args []string) error {
	helper := helm.NewHelper(helmVerbose, helmDryRun)

	if !helper.IsInstalled() {
		return fmt.Errorf("helm is not installed")
	}

	return helper.UpdateRepositories()
}

func runHelmSearch(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := helm.NewHelper(helmVerbose, helmDryRun)

	if !helper.IsInstalled() {
		return fmt.Errorf("helm is not installed")
	}

	query := ""
	if len(args) > 0 {
		query = args[0]
	}

	charts, err := helper.SearchCharts(query)
	if err != nil {
		return err
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(map[string]any{"charts": charts})
	}

	if len(charts) == 0 {
		fmt.Fprintln(os.Stdout, "No charts found")
		return nil
	}

	fmt.Fprintf(os.Stdout, "%-40s %-15s %-15s %s\n", "NAME", "CHART VERSION", "APP VERSION", "DESCRIPTION")
	for _, c := range charts {
		fmt.Fprintf(os.Stdout, "%-40s %-15s %-15s %s\n",
			truncateStr(c.Name, 40),
			c.Version,
			c.AppVersion,
			truncateStr(c.Description, 50))
	}

	return nil
}

func runHelmShow(cmd *cobra.Command, args []string) error {
	helper := helm.NewHelper(helmVerbose, helmDryRun)

	if !helper.IsInstalled() {
		return fmt.Errorf("helm is not installed")
	}

	return helper.ShowChart(args[0])
}

func runHelmTemplate(cmd *cobra.Command, args []string) error {
	helper := helm.NewHelper(helmVerbose, helmDryRun)

	if !helper.IsInstalled() {
		return fmt.Errorf("helm is not installed")
	}

	return helper.Template(args[0], args[1], helmValues)
}

func runHelmLint(cmd *cobra.Command, args []string) error {
	helper := helm.NewHelper(helmVerbose, helmDryRun)

	if !helper.IsInstalled() {
		return fmt.Errorf("helm is not installed")
	}

	chart := "."
	if len(args) > 0 {
		chart = args[0]
	}

	return helper.Lint(chart, helmStrict)
}

func runHelmCreate(cmd *cobra.Command, args []string) error {
	helper := helm.NewHelper(helmVerbose, helmDryRun)

	if !helper.IsInstalled() {
		return fmt.Errorf("helm is not installed")
	}

	if err := helper.Create(args[0]); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s Created chart %s\n", output.Success("✓"), args[0])
	return nil
}

func runHelmPlugins(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := helm.NewHelper(helmVerbose, helmDryRun)

	if !helper.IsInstalled() {
		return fmt.Errorf("helm is not installed")
	}

	plugins, err := helper.ListPlugins()
	if err != nil {
		return err
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(map[string]any{"plugins": plugins})
	}

	if len(plugins) == 0 {
		fmt.Fprintln(os.Stdout, "No plugins installed")
		return nil
	}

	fmt.Fprintf(os.Stdout, "%-20s %-15s %s\n", "NAME", "VERSION", "DESCRIPTION")
	for _, p := range plugins {
		fmt.Fprintf(os.Stdout, "%-20s %-15s %s\n", p.Name, p.Version, p.Description)
	}

	return nil
}
