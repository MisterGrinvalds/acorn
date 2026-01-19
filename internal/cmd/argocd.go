package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/devops/argocd"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	argocdOutputFormat string
	argocdVerbose      bool
	argocdProject      string
	argocdSelector     string
	argocdPrune        bool
	argocdForce        bool
	argocdFollow       bool
	argocdTree         bool
	argocdTimeout      int
	argocdSSO          bool
	argocdUsername     string
	argocdPassword     string
)

// argocdCmd represents the argocd command group
var argocdCmd = &cobra.Command{
	Use:   "argocd",
	Short: "ArgoCD GitOps commands",
	Long: `ArgoCD GitOps continuous delivery tool for Kubernetes.

Provides quick access to ArgoCD operations.

Examples:
  acorn argocd status      # Show installation status
  acorn argocd apps        # List applications
  acorn argocd sync <app>  # Sync application
  acorn argocd diff <app>  # Show app diff`,
	Aliases: []string{"argo", "acd"},
}

// argocdStatusCmd shows status
var argocdStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show ArgoCD installation status",
	Long: `Display ArgoCD CLI installation status and server connection.

Examples:
  acorn argocd status
  acorn argocd status -o json`,
	RunE: runArgocdStatus,
}

// argocdLoginCmd logs into ArgoCD
var argocdLoginCmd = &cobra.Command{
	Use:   "login [server]",
	Short: "Login to ArgoCD server",
	Long: `Login to an ArgoCD server.

Examples:
  acorn argocd login argocd.example.com
  acorn argocd login argocd.example.com --sso
  acorn argocd login argocd.example.com --username admin --password secret`,
	Args: cobra.MaximumNArgs(1),
	RunE: runArgocdLogin,
}

// argocdContextCmd manages contexts
var argocdContextCmd = &cobra.Command{
	Use:   "context [name]",
	Short: "List or switch ArgoCD contexts",
	Long: `List available ArgoCD contexts or switch to a specific one.

Examples:
  acorn argocd context
  acorn argocd context prod-cluster`,
	Aliases: []string{"ctx"},
	Args:    cobra.MaximumNArgs(1),
	RunE:    runArgocdContext,
}

// argocdAppsCmd lists applications
var argocdAppsCmd = &cobra.Command{
	Use:   "apps",
	Short: "List ArgoCD applications",
	Long: `List all ArgoCD applications.

Examples:
  acorn argocd apps
  acorn argocd apps --project myproject
  acorn argocd apps -o json`,
	Aliases: []string{"app", "list"},
	RunE:    runArgocdApps,
}

// argocdGetCmd gets application details
var argocdGetCmd = &cobra.Command{
	Use:   "get <app>",
	Short: "Get application details",
	Long: `Get detailed information about an application.

Examples:
  acorn argocd get myapp
  acorn argocd get myapp -o json`,
	Args: cobra.ExactArgs(1),
	RunE: runArgocdGet,
}

// argocdSyncCmd syncs an application
var argocdSyncCmd = &cobra.Command{
	Use:   "sync <app>",
	Short: "Sync an application",
	Long: `Sync an ArgoCD application to its target state.

Examples:
  acorn argocd sync myapp
  acorn argocd sync myapp --prune
  acorn argocd sync myapp --force`,
	Args: cobra.ExactArgs(1),
	RunE: runArgocdSync,
}

// argocdDiffCmd shows application diff
var argocdDiffCmd = &cobra.Command{
	Use:   "diff <app>",
	Short: "Show application diff",
	Long: `Show the diff between live and desired state.

Examples:
  acorn argocd diff myapp`,
	Args: cobra.ExactArgs(1),
	RunE: runArgocdDiff,
}

// argocdHistoryCmd shows application history
var argocdHistoryCmd = &cobra.Command{
	Use:   "history <app>",
	Short: "Show application deployment history",
	Long: `Show the deployment history for an application.

Examples:
  acorn argocd history myapp`,
	Args: cobra.ExactArgs(1),
	RunE: runArgocdHistory,
}

// argocdRollbackCmd rolls back an application
var argocdRollbackCmd = &cobra.Command{
	Use:   "rollback <app> <revision>",
	Short: "Rollback application to previous revision",
	Long: `Rollback an application to a previous revision.

Examples:
  acorn argocd rollback myapp 3`,
	Args: cobra.ExactArgs(2),
	RunE: runArgocdRollback,
}

// argocdWaitCmd waits for application sync
var argocdWaitCmd = &cobra.Command{
	Use:   "wait <app>",
	Short: "Wait for application to sync",
	Long: `Wait for an application to reach synced and healthy state.

Examples:
  acorn argocd wait myapp
  acorn argocd wait myapp --timeout 300`,
	Args: cobra.ExactArgs(1),
	RunE: runArgocdWait,
}

// argocdLogsCmd shows application logs
var argocdLogsCmd = &cobra.Command{
	Use:   "logs <app>",
	Short: "View application logs",
	Long: `Stream logs for an application.

Examples:
  acorn argocd logs myapp
  acorn argocd logs myapp --follow`,
	Args: cobra.ExactArgs(1),
	RunE: runArgocdLogs,
}

// argocdResourcesCmd shows application resources
var argocdResourcesCmd = &cobra.Command{
	Use:   "resources <app>",
	Short: "Show application resources",
	Long: `Show resources managed by an application.

Examples:
  acorn argocd resources myapp
  acorn argocd resources myapp --tree`,
	Aliases: []string{"tree"},
	Args:    cobra.ExactArgs(1),
	RunE:    runArgocdResources,
}

// argocdClustersCmd lists clusters
var argocdClustersCmd = &cobra.Command{
	Use:   "clusters",
	Short: "List registered clusters",
	Long: `List all clusters registered with ArgoCD.

Examples:
  acorn argocd clusters
  acorn argocd clusters -o json`,
	Aliases: []string{"cluster"},
	RunE:    runArgocdClusters,
}

// argocdReposCmd lists repositories
var argocdReposCmd = &cobra.Command{
	Use:   "repos",
	Short: "List registered repositories",
	Long: `List all repositories registered with ArgoCD.

Examples:
  acorn argocd repos
  acorn argocd repos -o json`,
	Aliases: []string{"repo"},
	RunE:    runArgocdRepos,
}

// argocdProjectsCmd lists projects
var argocdProjectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "List ArgoCD projects",
	Long: `List all ArgoCD projects.

Examples:
  acorn argocd projects
  acorn argocd projects -o json`,
	Aliases: []string{"proj"},
	RunE:    runArgocdProjects,
}

func init() {
	devopsCmd.AddCommand(argocdCmd)

	// Add subcommands
	argocdCmd.AddCommand(argocdStatusCmd)
	argocdCmd.AddCommand(argocdLoginCmd)
	argocdCmd.AddCommand(argocdContextCmd)
	argocdCmd.AddCommand(argocdAppsCmd)
	argocdCmd.AddCommand(argocdGetCmd)
	argocdCmd.AddCommand(argocdSyncCmd)
	argocdCmd.AddCommand(argocdDiffCmd)
	argocdCmd.AddCommand(argocdHistoryCmd)
	argocdCmd.AddCommand(argocdRollbackCmd)
	argocdCmd.AddCommand(argocdWaitCmd)
	argocdCmd.AddCommand(argocdLogsCmd)
	argocdCmd.AddCommand(argocdResourcesCmd)
	argocdCmd.AddCommand(argocdClustersCmd)
	argocdCmd.AddCommand(argocdReposCmd)
	argocdCmd.AddCommand(argocdProjectsCmd)

	// Persistent flags
	argocdCmd.PersistentFlags().StringVarP(&argocdOutputFormat, "output", "o", "table",
		"Output format (table|json|yaml)")
	argocdCmd.PersistentFlags().BoolVarP(&argocdVerbose, "verbose", "v", false,
		"Show verbose output")

	// Apps command flags
	argocdAppsCmd.Flags().StringVar(&argocdProject, "project", "",
		"Filter by project")
	argocdAppsCmd.Flags().StringVar(&argocdSelector, "selector", "",
		"Filter by label selector")

	// Sync command flags
	argocdSyncCmd.Flags().BoolVar(&argocdPrune, "prune", false,
		"Delete resources that no longer exist in git")
	argocdSyncCmd.Flags().BoolVar(&argocdForce, "force", false,
		"Force sync")

	// Wait command flags
	argocdWaitCmd.Flags().IntVar(&argocdTimeout, "timeout", 0,
		"Timeout in seconds")

	// Logs command flags
	argocdLogsCmd.Flags().BoolVarP(&argocdFollow, "follow", "f", false,
		"Follow logs")

	// Resources command flags
	argocdResourcesCmd.Flags().BoolVar(&argocdTree, "tree", false,
		"Show resources as tree")

	// Login command flags
	argocdLoginCmd.Flags().BoolVar(&argocdSSO, "sso", false,
		"Use SSO login")
	argocdLoginCmd.Flags().StringVar(&argocdUsername, "username", "",
		"Username for login")
	argocdLoginCmd.Flags().StringVar(&argocdPassword, "password", "",
		"Password for login")
}

func runArgocdStatus(cmd *cobra.Command, args []string) error {
	helper := argocd.NewHelper(argocdVerbose)
	status := helper.GetStatus()

	format, err := output.ParseFormat(argocdOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(status)
	}

	fmt.Fprintf(os.Stdout, "Installed:        %v\n", status.Installed)
	if status.Version != "" {
		fmt.Fprintf(os.Stdout, "Version:          %s\n", status.Version)
	}
	if status.Location != "" {
		fmt.Fprintf(os.Stdout, "Location:         %s\n", status.Location)
	}
	if status.ConfigDir != "" {
		fmt.Fprintf(os.Stdout, "Config Dir:       %s\n", status.ConfigDir)
	}
	fmt.Fprintf(os.Stdout, "Server Connected: %v\n", status.ServerConnected)
	if status.CurrentContext != "" {
		fmt.Fprintf(os.Stdout, "Current Context:  %s\n", status.CurrentContext)
	}
	if status.Server != "" {
		fmt.Fprintf(os.Stdout, "Server:           %s\n", status.Server)
	}

	return nil
}

func runArgocdLogin(cmd *cobra.Command, args []string) error {
	helper := argocd.NewHelper(argocdVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("argocd CLI is not installed. Install with: brew install argocd")
	}

	server := os.Getenv("ARGOCD_SERVER")
	if len(args) > 0 {
		server = args[0]
	}
	if server == "" {
		return fmt.Errorf("server address required. Provide as argument or set ARGOCD_SERVER")
	}

	return helper.Login(server, argocdSSO, argocdUsername, argocdPassword)
}

func runArgocdContext(cmd *cobra.Command, args []string) error {
	helper := argocd.NewHelper(argocdVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("argocd CLI is not installed")
	}

	// If name provided, switch context
	if len(args) > 0 {
		return helper.SwitchContext(args[0])
	}

	// List contexts
	contexts, err := helper.ListContexts()
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(argocdOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(map[string]any{"contexts": contexts})
	}

	fmt.Println("ArgoCD Contexts:")
	for _, ctx := range contexts {
		marker := "  "
		if ctx.Current {
			marker = "* "
		}
		fmt.Printf("%s%s (%s)\n", marker, ctx.Name, ctx.Server)
	}

	return nil
}

func runArgocdApps(cmd *cobra.Command, args []string) error {
	helper := argocd.NewHelper(argocdVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("argocd CLI is not installed")
	}

	if !helper.IsServerConnected() {
		return fmt.Errorf("not connected to ArgoCD server. Run: argocd login")
	}

	apps, err := helper.ListApps(argocdProject, argocdSelector)
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(argocdOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(map[string]any{"apps": apps})
	}

	if len(apps) == 0 {
		fmt.Println("No applications found")
		return nil
	}

	fmt.Printf("%-30s %-15s %-15s %-15s\n", "NAME", "PROJECT", "SYNC", "HEALTH")
	for _, app := range apps {
		fmt.Printf("%-30s %-15s %-15s %-15s\n",
			truncateArgo(app.Name, 30),
			truncateArgo(app.Project, 15),
			truncateArgo(app.SyncStatus, 15),
			truncateArgo(app.Health, 15))
	}

	return nil
}

func runArgocdGet(cmd *cobra.Command, args []string) error {
	helper := argocd.NewHelper(argocdVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("argocd CLI is not installed")
	}

	app, err := helper.GetApp(args[0])
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(argocdOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(app)
	}

	fmt.Fprintf(os.Stdout, "Name:        %s\n", app.Name)
	fmt.Fprintf(os.Stdout, "Project:     %s\n", app.Project)
	fmt.Fprintf(os.Stdout, "Sync Status: %s\n", app.SyncStatus)
	fmt.Fprintf(os.Stdout, "Health:      %s\n", app.Health)
	if app.Repo != "" {
		fmt.Fprintf(os.Stdout, "Repo:        %s\n", app.Repo)
	}
	if app.Path != "" {
		fmt.Fprintf(os.Stdout, "Path:        %s\n", app.Path)
	}
	if app.Destination != "" {
		fmt.Fprintf(os.Stdout, "Destination: %s\n", app.Destination)
	}

	return nil
}

func runArgocdSync(cmd *cobra.Command, args []string) error {
	helper := argocd.NewHelper(argocdVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("argocd CLI is not installed")
	}

	return helper.SyncApp(args[0], argocdPrune, argocdForce)
}

func runArgocdDiff(cmd *cobra.Command, args []string) error {
	helper := argocd.NewHelper(argocdVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("argocd CLI is not installed")
	}

	return helper.DiffApp(args[0])
}

func runArgocdHistory(cmd *cobra.Command, args []string) error {
	helper := argocd.NewHelper(argocdVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("argocd CLI is not installed")
	}

	return helper.GetAppHistory(args[0])
}

func runArgocdRollback(cmd *cobra.Command, args []string) error {
	helper := argocd.NewHelper(argocdVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("argocd CLI is not installed")
	}

	var revision int
	if _, err := fmt.Sscanf(args[1], "%d", &revision); err != nil {
		return fmt.Errorf("invalid revision: %s", args[1])
	}

	return helper.RollbackApp(args[0], revision)
}

func runArgocdWait(cmd *cobra.Command, args []string) error {
	helper := argocd.NewHelper(argocdVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("argocd CLI is not installed")
	}

	return helper.WaitForApp(args[0], argocdTimeout)
}

func runArgocdLogs(cmd *cobra.Command, args []string) error {
	helper := argocd.NewHelper(argocdVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("argocd CLI is not installed")
	}

	return helper.GetAppLogs(args[0], argocdFollow, "", "")
}

func runArgocdResources(cmd *cobra.Command, args []string) error {
	helper := argocd.NewHelper(argocdVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("argocd CLI is not installed")
	}

	return helper.GetAppResources(args[0], argocdTree)
}

func runArgocdClusters(cmd *cobra.Command, args []string) error {
	helper := argocd.NewHelper(argocdVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("argocd CLI is not installed")
	}

	clusters, err := helper.ListClusters()
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(argocdOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(map[string]any{"clusters": clusters})
	}

	if len(clusters) == 0 {
		fmt.Println("No clusters registered")
		return nil
	}

	fmt.Printf("%-40s %-30s\n", "SERVER", "NAME")
	for _, cluster := range clusters {
		fmt.Printf("%-40s %-30s\n",
			truncateArgo(cluster.Server, 40),
			truncateArgo(cluster.Name, 30))
	}

	return nil
}

func runArgocdRepos(cmd *cobra.Command, args []string) error {
	helper := argocd.NewHelper(argocdVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("argocd CLI is not installed")
	}

	repos, err := helper.ListRepos()
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(argocdOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(map[string]any{"repos": repos})
	}

	if len(repos) == 0 {
		fmt.Println("No repositories registered")
		return nil
	}

	fmt.Printf("%-60s %-10s %-10s\n", "URL", "TYPE", "STATUS")
	for _, repo := range repos {
		fmt.Printf("%-60s %-10s %-10s\n",
			truncateArgo(repo.URL, 60),
			truncateArgo(repo.Type, 10),
			truncateArgo(repo.Status, 10))
	}

	return nil
}

func runArgocdProjects(cmd *cobra.Command, args []string) error {
	helper := argocd.NewHelper(argocdVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("argocd CLI is not installed")
	}

	projects, err := helper.ListProjects()
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(argocdOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(map[string]any{"projects": projects})
	}

	if len(projects) == 0 {
		fmt.Println("No projects found")
		return nil
	}

	fmt.Println("Projects:")
	for _, proj := range projects {
		fmt.Printf("  %s\n", proj)
	}

	return nil
}

func truncateArgo(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
