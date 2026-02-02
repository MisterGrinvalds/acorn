package cmd

import (
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/devops/k9s"
	"github.com/mistergrinvalds/acorn/internal/utils/configcmd"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/spf13/cobra"
)

var (
	k9sVerbose   bool
	k9sContext   string
	k9sNamespace string
	k9sCommand   string
	k9sReadonly  bool
	k9sHeadless  bool
)

// k9sCmd represents the k9s command group
var k9sCmd = &cobra.Command{
	Use:   "k9s",
	Short: "k9s terminal UI commands",
	Long: `k9s terminal UI for Kubernetes cluster management.

Provides quick access to k9s and related utilities.

Examples:
  acorn k9s              # Launch k9s
  acorn k9s status       # Show installation status
  acorn k9s keys         # Show keybindings
  acorn k9s config       # Edit main config
  acorn k9s info         # Show k9s info`,
	RunE: runK9sLaunch,
}

// k9sStatusCmd shows status
var k9sStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show k9s installation status",
	Long: `Display k9s installation status and cluster connection.

Examples:
  acorn k9s status
  acorn k9s status -o json`,
	RunE: runK9sStatus,
}

// k9sLaunchCmd launches k9s
var k9sLaunchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch k9s",
	Long: `Start the k9s terminal UI.

Examples:
  acorn k9s launch
  acorn k9s launch --context prod
  acorn k9s launch -n kube-system
  acorn k9s launch -c pod --readonly`,
	Aliases: []string{"start", "open"},
	RunE:    runK9sLaunch,
}

// k9sKeysCmd shows keybindings
var k9sKeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Show k9s keybindings",
	Long: `Display k9s keyboard shortcuts and commands.

Examples:
  acorn k9s keys
  acorn k9s keys -o json`,
	Aliases: []string{"keybindings", "shortcuts"},
	RunE:    runK9sKeys,
}

// k9sEditCmd opens k9s config files in editor
var k9sEditCmd = &cobra.Command{
	Use:   "edit [type]",
	Short: "Edit k9s configuration in editor",
	Long: `Open k9s configuration files in your editor.

Config types:
  main (default) - Main k9s configuration
  plugins        - Plugin definitions
  hotkeys        - Hotkey definitions
  views          - Custom column views
  aliases        - Command aliases

Creates default config files if they don't exist.

Examples:
  acorn k9s edit
  acorn k9s edit plugins
  acorn k9s edit hotkeys`,
	Args: cobra.MaximumNArgs(1),
	RunE: runK9sConfig,
}

// k9sInfoCmd shows k9s info
var k9sInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show k9s configuration info",
	Long: `Display k9s configuration and log locations.

Examples:
  acorn k9s info
  acorn k9s info -o json`,
	RunE: runK9sInfo,
}

// k9sSkinsCmd lists skins
var k9sSkinsCmd = &cobra.Command{
	Use:   "skins",
	Short: "List available k9s skins",
	Long: `List available k9s skin files.

Examples:
  acorn k9s skins
  acorn k9s skins -o json`,
	RunE: runK9sSkins,
}

// k9sContextsCmd lists contexts
var k9sContextsCmd = &cobra.Command{
	Use:   "contexts",
	Short: "List available Kubernetes contexts",
	Long: `List available kubectl contexts for k9s.

Examples:
  acorn k9s contexts
  acorn k9s contexts -o json`,
	Aliases: []string{"ctx"},
	RunE:    runK9sContexts,
}

// k9sNamespacesCmd lists namespaces
var k9sNamespacesCmd = &cobra.Command{
	Use:   "namespaces",
	Short: "List available namespaces",
	Long: `List available Kubernetes namespaces.

Examples:
  acorn k9s namespaces
  acorn k9s namespaces -o json`,
	Aliases: []string{"ns"},
	RunE:    runK9sNamespaces,
}

// k9s generate is now provided by the universal config router: acorn k9s config generate

func init() {
	devopsCmd.AddCommand(k9sCmd)

	// Add subcommands
	k9sCmd.AddCommand(k9sStatusCmd)
	k9sCmd.AddCommand(k9sLaunchCmd)
	k9sCmd.AddCommand(k9sKeysCmd)
	k9sCmd.AddCommand(configcmd.NewConfigRouter("k9s"))
	k9sCmd.AddCommand(k9sEditCmd)
	k9sCmd.AddCommand(k9sInfoCmd)
	k9sCmd.AddCommand(k9sSkinsCmd)
	k9sCmd.AddCommand(k9sContextsCmd)
	k9sCmd.AddCommand(k9sNamespacesCmd)

	// Persistent flags
	k9sCmd.PersistentFlags().BoolVarP(&k9sVerbose, "verbose", "v", false,
		"Show verbose output")

	// Launch command flags
	k9sLaunchCmd.Flags().StringVar(&k9sContext, "context", "",
		"Kubernetes context to use")
	k9sLaunchCmd.Flags().StringVarP(&k9sNamespace, "namespace", "n", "",
		"Namespace to start in")
	k9sLaunchCmd.Flags().StringVarP(&k9sCommand, "command", "c", "",
		"Initial resource command (e.g., pod, deploy)")
	k9sLaunchCmd.Flags().BoolVar(&k9sReadonly, "readonly", false,
		"Disable all modification commands")
	k9sLaunchCmd.Flags().BoolVar(&k9sHeadless, "headless", false,
		"Hide header, logo, and breadcrumbs")

	// Root command flags for default launch
	k9sCmd.Flags().StringVar(&k9sContext, "context", "",
		"Kubernetes context to use")
	k9sCmd.Flags().StringVarP(&k9sNamespace, "namespace", "n", "",
		"Namespace to start in")
	k9sCmd.Flags().StringVarP(&k9sCommand, "command", "c", "",
		"Initial resource command (e.g., pod, deploy)")
	k9sCmd.Flags().BoolVar(&k9sReadonly, "readonly", false,
		"Disable all modification commands")
	k9sCmd.Flags().BoolVar(&k9sHeadless, "headless", false,
		"Hide header, logo, and breadcrumbs")
}

func runK9sStatus(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := k9s.NewHelper(k9sVerbose)
	status := helper.GetStatus()

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	fmt.Fprintf(os.Stdout, "Installed:         %v\n", status.Installed)
	if status.Version != "" {
		fmt.Fprintf(os.Stdout, "Version:           %s\n", status.Version)
	}
	if status.Location != "" {
		fmt.Fprintf(os.Stdout, "Location:          %s\n", status.Location)
	}
	if status.ConfigDir != "" {
		fmt.Fprintf(os.Stdout, "Config Dir:        %s\n", status.ConfigDir)
	}
	if status.DataDir != "" {
		fmt.Fprintf(os.Stdout, "Data Dir:          %s\n", status.DataDir)
	}
	fmt.Fprintf(os.Stdout, "Cluster Connected: %v\n", status.ClusterConnected)
	if status.CurrentContext != "" {
		fmt.Fprintf(os.Stdout, "Current Context:   %s\n", status.CurrentContext)
	}

	return nil
}

func runK9sLaunch(cmd *cobra.Command, args []string) error {
	helper := k9s.NewHelper(k9sVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("k9s is not installed. Install with: brew install derailed/k9s/k9s")
	}

	if !helper.IsClusterConnected() {
		return fmt.Errorf("cannot connect to Kubernetes cluster. Check your kubeconfig")
	}

	return helper.Launch(k9sContext, k9sNamespace, k9sCommand, k9sReadonly, k9sHeadless)
}

func runK9sKeys(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := k9s.NewHelper(k9sVerbose)

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(map[string]any{"keybindings": helper.GetKeybindingsList()})
	}

	fmt.Print(helper.GetKeybindings())
	return nil
}

func runK9sConfig(cmd *cobra.Command, args []string) error {
	helper := k9s.NewHelper(k9sVerbose)

	configType := "main"
	if len(args) > 0 {
		configType = args[0]
	}

	return helper.OpenConfig(configType)
}

func runK9sInfo(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := k9s.NewHelper(k9sVerbose)

	if !helper.IsInstalled() {
		return fmt.Errorf("k9s is not installed")
	}

	info, err := helper.GetInfo()
	if err != nil {
		return err
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(info)
	}

	fmt.Fprintf(os.Stdout, "Config Dir:    %s\n", info.ConfigDir)
	fmt.Fprintf(os.Stdout, "Logs Dir:      %s\n", info.LogsDir)
	fmt.Fprintf(os.Stdout, "Screen Dumps:  %s\n", info.ScreenDump)

	return nil
}

func runK9sSkins(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := k9s.NewHelper(k9sVerbose)

	skins, err := helper.ListSkins()
	if err != nil {
		return err
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(map[string]any{"skins": skins})
	}

	if len(skins) == 0 {
		fmt.Println("No skins found")
		fmt.Printf("Skins directory: %s\n", helper.GetSkinsDir())
		return nil
	}

	fmt.Println("Available skins:")
	for _, skin := range skins {
		fmt.Printf("  %s\n", skin)
	}

	return nil
}

func runK9sContexts(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := k9s.NewHelper(k9sVerbose)

	contexts, err := helper.GetContexts()
	if err != nil {
		return err
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(map[string]any{"contexts": contexts})
	}

	currentContext := helper.GetCurrentContext()
	fmt.Println("Available contexts:")
	for _, ctx := range contexts {
		if ctx == currentContext {
			fmt.Printf("  * %s (current)\n", ctx)
		} else {
			fmt.Printf("    %s\n", ctx)
		}
	}

	return nil
}

func runK9sNamespaces(cmd *cobra.Command, args []string) error {
	ioHelper := ioutils.IO(cmd)
	helper := k9s.NewHelper(k9sVerbose)

	if !helper.IsClusterConnected() {
		return fmt.Errorf("cannot connect to Kubernetes cluster")
	}

	namespaces, err := helper.GetNamespaces()
	if err != nil {
		return err
	}

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(map[string]any{"namespaces": namespaces})
	}

	fmt.Println("Available namespaces:")
	for _, ns := range namespaces {
		fmt.Printf("  %s\n", ns)
	}

	return nil
}

// runK9sGenerate has been replaced by the universal config router: acorn k9s config generate
