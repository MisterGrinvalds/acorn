package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/kubernetes"
	"github.com/mistergrinvalds/acorn/internal/output"
	"github.com/spf13/cobra"
)

var (
	k8sOutputFormat string
	k8sVerbose      bool
	k8sDryRun       bool
)

// k8sCmd represents the kubernetes command group
var k8sCmd = &cobra.Command{
	Use:   "k8s",
	Short: "Kubernetes helper commands",
	Long: `Kubernetes helper commands for cluster management.

Provides context/namespace management, pod operations, and resource viewing.

Examples:
  acorn k8s info          # Show current context info
  acorn k8s context       # List/switch contexts
  acorn k8s namespace     # List/switch namespaces
  acorn k8s pods          # List pods
  acorn k8s all           # Show all resources
  acorn k8s clean         # Clean evicted pods`,
	Aliases: []string{"kube", "kubernetes"},
}

// k8sInfoCmd shows context info
var k8sInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show current context info",
	Long: `Display current Kubernetes context, namespace, and server.

Examples:
  acorn k8s info
  acorn k8s info -o json`,
	RunE: runK8sInfo,
}

// k8sContextCmd manages contexts
var k8sContextCmd = &cobra.Command{
	Use:   "context [name]",
	Short: "List or switch contexts",
	Long: `List available contexts or switch to a specific context.

Without arguments, lists all contexts.
With a context name, switches to that context.

Examples:
  acorn k8s context              # List contexts
  acorn k8s context minikube     # Switch to minikube`,
	Aliases: []string{"ctx"},
	RunE:    runK8sContext,
}

// k8sNamespaceCmd manages namespaces
var k8sNamespaceCmd = &cobra.Command{
	Use:   "namespace [name]",
	Short: "List or switch namespaces",
	Long: `List available namespaces or switch to a specific namespace.

Without arguments, lists all namespaces.
With a namespace name, switches to that namespace.

Examples:
  acorn k8s namespace            # List namespaces
  acorn k8s namespace kube-system  # Switch to kube-system`,
	Aliases: []string{"ns"},
	RunE:    runK8sNamespace,
}

// k8sPodsCmd lists pods
var k8sPodsCmd = &cobra.Command{
	Use:   "pods [filter]",
	Short: "List pods with optional filter",
	Long: `List pods in the current namespace.

With a filter argument, only shows pods matching the filter.

Examples:
  acorn k8s pods
  acorn k8s pods nginx`,
	RunE: runK8sPods,
}

// k8sAllCmd shows all resources
var k8sAllCmd = &cobra.Command{
	Use:   "all [namespace]",
	Short: "Show all resources in namespace",
	Long: `Display pods, services, and deployments in a namespace.

Uses current namespace if not specified.

Examples:
  acorn k8s all
  acorn k8s all kube-system`,
	RunE: runK8sAll,
}

// k8sCleanCmd cleans evicted pods
var k8sCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean evicted pods",
	Long: `Delete all evicted pods across all namespaces.

Examples:
  acorn k8s clean
  acorn k8s clean --dry-run`,
	RunE: runK8sClean,
}

func init() {
	rootCmd.AddCommand(k8sCmd)

	// Add subcommands
	k8sCmd.AddCommand(k8sInfoCmd)
	k8sCmd.AddCommand(k8sContextCmd)
	k8sCmd.AddCommand(k8sNamespaceCmd)
	k8sCmd.AddCommand(k8sPodsCmd)
	k8sCmd.AddCommand(k8sAllCmd)
	k8sCmd.AddCommand(k8sCleanCmd)

	// Persistent flags
	k8sCmd.PersistentFlags().StringVarP(&k8sOutputFormat, "output", "o", "table",
		"Output format (table|json|yaml)")
	k8sCmd.PersistentFlags().BoolVarP(&k8sVerbose, "verbose", "v", false,
		"Show verbose output")
	k8sCmd.PersistentFlags().BoolVar(&k8sDryRun, "dry-run", false,
		"Show what would be done without executing")
}

func runK8sInfo(cmd *cobra.Command, args []string) error {
	helper := kubernetes.NewHelper(k8sVerbose, k8sDryRun)

	if !helper.IsKubectlInstalled() {
		return fmt.Errorf("kubectl is not installed")
	}

	info, err := helper.GetContextInfo()
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(k8sOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(info)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "Context:   %s\n", info.Context)
	fmt.Fprintf(os.Stdout, "Namespace: %s\n", info.Namespace)
	fmt.Fprintf(os.Stdout, "Server:    %s\n", info.Server)

	return nil
}

func runK8sContext(cmd *cobra.Command, args []string) error {
	helper := kubernetes.NewHelper(k8sVerbose, k8sDryRun)

	if !helper.IsKubectlInstalled() {
		return fmt.Errorf("kubectl is not installed")
	}

	if len(args) == 0 {
		// List contexts
		contexts, err := helper.GetContexts()
		if err != nil {
			return err
		}

		format, err := output.ParseFormat(k8sOutputFormat)
		if err != nil {
			return err
		}

		if format != output.FormatTable {
			printer := output.NewPrinter(os.Stdout, format)
			return printer.Print(map[string]interface{}{"contexts": contexts})
		}

		// Table format
		for _, ctx := range contexts {
			marker := " "
			if ctx.Current {
				marker = "*"
			}
			fmt.Fprintf(os.Stdout, "%s %s\n", marker, ctx.Name)
		}

		return nil
	}

	// Switch context
	return helper.UseContext(args[0])
}

func runK8sNamespace(cmd *cobra.Command, args []string) error {
	helper := kubernetes.NewHelper(k8sVerbose, k8sDryRun)

	if !helper.IsKubectlInstalled() {
		return fmt.Errorf("kubectl is not installed")
	}

	if len(args) == 0 {
		// List namespaces
		namespaces, err := helper.GetNamespaces()
		if err != nil {
			return err
		}

		format, err := output.ParseFormat(k8sOutputFormat)
		if err != nil {
			return err
		}

		if format != output.FormatTable {
			printer := output.NewPrinter(os.Stdout, format)
			return printer.Print(map[string]interface{}{"namespaces": namespaces})
		}

		// Table format
		for _, ns := range namespaces {
			fmt.Fprintf(os.Stdout, "%-30s %s\n", ns.Name, ns.Status)
		}

		return nil
	}

	// Switch namespace
	return helper.UseNamespace(args[0])
}

func runK8sPods(cmd *cobra.Command, args []string) error {
	helper := kubernetes.NewHelper(k8sVerbose, k8sDryRun)

	if !helper.IsKubectlInstalled() {
		return fmt.Errorf("kubectl is not installed")
	}

	filter := ""
	if len(args) > 0 {
		filter = args[0]
	}

	pods, err := helper.GetPods(filter)
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(k8sOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(map[string]interface{}{"pods": pods})
	}

	// Table format
	if len(pods) == 0 {
		fmt.Fprintln(os.Stdout, "No pods found")
		return nil
	}

	fmt.Fprintf(os.Stdout, "%-50s %-10s %-12s %-10s\n", "NAME", "READY", "STATUS", "RESTARTS")
	for _, pod := range pods {
		fmt.Fprintf(os.Stdout, "%-50s %-10s %-12s %-10s\n", pod.Name, pod.Ready, pod.Status, pod.Restarts)
	}

	return nil
}

func runK8sAll(cmd *cobra.Command, args []string) error {
	helper := kubernetes.NewHelper(k8sVerbose, k8sDryRun)

	if !helper.IsKubectlInstalled() {
		return fmt.Errorf("kubectl is not installed")
	}

	namespace := ""
	if len(args) > 0 {
		namespace = args[0]
	}

	return helper.GetAllResources(namespace)
}

func runK8sClean(cmd *cobra.Command, args []string) error {
	helper := kubernetes.NewHelper(k8sVerbose, k8sDryRun)

	if !helper.IsKubectlInstalled() {
		return fmt.Errorf("kubectl is not installed")
	}

	count, err := helper.CleanEvictedPods()
	if err != nil {
		return err
	}

	if count == 0 {
		fmt.Fprintln(os.Stdout, "No evicted pods found")
	} else if k8sDryRun {
		fmt.Fprintf(os.Stdout, "Would delete %d evicted pods\n", count)
	} else {
		fmt.Fprintf(os.Stdout, "%s Deleted %d evicted pods\n", output.Success("✓"), count)
	}

	return nil
}
