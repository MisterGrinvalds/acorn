package cmd

import (
	"github.com/mistergrinvalds/acorn/internal/components"
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/components/tailscale"
	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/mistergrinvalds/acorn/internal/utils/configcmd"
	"github.com/spf13/cobra"
)

var (
	tailscaleVerbose bool
	tailscaleDryRun  bool
)

// tailscaleCmd represents the tailscale command group
var tailscaleCmd = &cobra.Command{
	Use:   "tailscale",
	Short: "Tailscale VPN helper commands",
	Long: `Tailscale mesh VPN helper commands.

Provides status, peer management, and exit node operations.

Examples:
  acorn network tailscale status       # Show connection status
  acorn network tailscale peers        # List all peers
  acorn network tailscale ping myhost  # Ping a peer
  acorn network tailscale exit-nodes   # List exit nodes
  acorn network tailscale set-exit us  # Set exit node`,
	Aliases: []string{"ts"},
}

// tailscaleStatusCmd shows Tailscale status
var tailscaleStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Tailscale connection status",
	Long: `Display Tailscale installation and connection status.

Shows version, connection state, IP address, and exit node status.

Examples:
  acorn network tailscale status
  acorn network tailscale status -o json`,
	RunE: runTailscaleStatus,
}

// tailscalePeersCmd lists peers
var tailscalePeersCmd = &cobra.Command{
	Use:   "peers",
	Short: "List Tailscale peers",
	Long: `List all Tailscale peers in your network.

Use --online to show only online peers.

Examples:
  acorn network tailscale peers
  acorn network tailscale peers --online
  acorn network tailscale peers -o json`,
	RunE: runTailscalePeers,
}

// tailscalePingCmd pings a peer
var tailscalePingCmd = &cobra.Command{
	Use:   "ping <peer>",
	Short: "Ping a Tailscale peer",
	Long: `Ping a peer in your Tailscale network.

The peer can be specified by hostname, DNS name, or IP address.

Examples:
  acorn network tailscale ping myhost
  acorn network tailscale ping 100.64.0.1`,
	Args: cobra.ExactArgs(1),
	RunE: runTailscalePing,
}

// tailscaleExitNodesCmd lists exit nodes
var tailscaleExitNodesCmd = &cobra.Command{
	Use:   "exit-nodes",
	Short: "List available exit nodes",
	Long: `List peers that can be used as exit nodes.

Exit nodes route all your internet traffic through that peer.

Examples:
  acorn network tailscale exit-nodes
  acorn network tailscale exit-nodes -o json`,
	Aliases: []string{"exits"},
	RunE:    runTailscaleExitNodes,
}

// tailscaleSetExitCmd sets exit node
var tailscaleSetExitCmd = &cobra.Command{
	Use:   "set-exit <node>",
	Short: "Set an exit node",
	Long: `Set a peer as your exit node.

All internet traffic will be routed through the specified peer.

Examples:
  acorn network tailscale set-exit us-east
  acorn network tailscale set-exit myserver`,
	Args: cobra.ExactArgs(1),
	RunE: runTailscaleSetExit,
}

// tailscaleDisableExitCmd disables exit node
var tailscaleDisableExitCmd = &cobra.Command{
	Use:   "disable-exit",
	Short: "Disable exit node",
	Long: `Disable the current exit node.

Traffic will no longer be routed through an exit node.

Examples:
  acorn network tailscale disable-exit`,
	Aliases: []string{"no-exit"},
	RunE:    runTailscaleDisableExit,
}

// tailscaleUpCmd brings Tailscale up
var tailscaleUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Connect to Tailscale",
	Long: `Bring Tailscale up and connect to the network.

Examples:
  acorn network tailscale up`,
	RunE: runTailscaleUp,
}

// tailscaleDownCmd brings Tailscale down
var tailscaleDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Disconnect from Tailscale",
	Long: `Bring Tailscale down and disconnect from the network.

Examples:
  acorn network tailscale down`,
	RunE: runTailscaleDown,
}

// tailscaleNetcheckCmd runs network check
var tailscaleNetcheckCmd = &cobra.Command{
	Use:   "netcheck",
	Short: "Run network diagnostics",
	Long: `Run Tailscale network diagnostics.

Checks connectivity to DERP servers, NAT type, and more.

Examples:
  acorn network tailscale netcheck`,
	Aliases: []string{"diag"},
	RunE:    runTailscaleNetcheck,
}

var (
	tailscaleOnlineOnly bool
)

func init() {

	// Add subcommands
	tailscaleCmd.AddCommand(tailscaleStatusCmd)
	tailscaleCmd.AddCommand(tailscalePeersCmd)
	tailscaleCmd.AddCommand(tailscalePingCmd)
	tailscaleCmd.AddCommand(tailscaleExitNodesCmd)
	tailscaleCmd.AddCommand(tailscaleSetExitCmd)
	tailscaleCmd.AddCommand(tailscaleDisableExitCmd)
	tailscaleCmd.AddCommand(tailscaleUpCmd)
	tailscaleCmd.AddCommand(tailscaleDownCmd)
	tailscaleCmd.AddCommand(tailscaleNetcheckCmd)
	tailscaleCmd.AddCommand(configcmd.NewConfigRouter("tailscale"))

	// Persistent flags
	tailscaleCmd.PersistentFlags().BoolVarP(&tailscaleVerbose, "verbose", "v", false,
		"Show verbose output")
	tailscaleCmd.PersistentFlags().BoolVar(&tailscaleDryRun, "dry-run", false,
		"Show what would be done without executing")

	// Command-specific flags
	tailscalePeersCmd.Flags().BoolVar(&tailscaleOnlineOnly, "online", false, "Show only online peers")
}

func runTailscaleStatus(cmd *cobra.Command, args []string) error {
	helper := tailscale.NewHelper(tailscaleVerbose, tailscaleDryRun)

	if !helper.IsInstalled() {
		return fmt.Errorf("tailscale is not installed")
	}

	status := helper.GetStatus()

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(status)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "Installed:      %v\n", status.Installed)
	if status.Version != "" {
		fmt.Fprintf(os.Stdout, "Version:        %s\n", status.Version)
	}
	fmt.Fprintf(os.Stdout, "Connected:      %v\n", status.Connected)
	fmt.Fprintf(os.Stdout, "Backend State:  %s\n", status.BackendState)
	fmt.Fprintf(os.Stdout, "Online:         %v\n", status.Online)
	if status.IP != "" {
		fmt.Fprintf(os.Stdout, "IP:             %s\n", status.IP)
	}
	if status.Hostname != "" {
		fmt.Fprintf(os.Stdout, "Hostname:       %s\n", status.Hostname)
	}
	if status.DNSName != "" {
		fmt.Fprintf(os.Stdout, "DNS Name:       %s\n", status.DNSName)
	}
	fmt.Fprintf(os.Stdout, "MagicDNS:       %v\n", status.MagicDNS)
	if status.ExitNode != "" {
		fmt.Fprintf(os.Stdout, "Exit Node:      %s (%s)\n", status.ExitNode, status.ExitNodeIP)
	}

	return nil
}

func runTailscalePeers(cmd *cobra.Command, args []string) error {
	helper := tailscale.NewHelper(tailscaleVerbose, tailscaleDryRun)

	if !helper.IsInstalled() {
		return fmt.Errorf("tailscale is not installed")
	}

	var peers []tailscale.Peer
	var err error

	if tailscaleOnlineOnly {
		peers, err = helper.GetOnlinePeers()
	} else {
		peers, err = helper.GetPeers()
	}
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(map[string]any{"peers": peers})
	}

	if len(peers) == 0 {
		fmt.Fprintln(os.Stdout, "No peers found")
		return nil
	}

	fmt.Fprintf(os.Stdout, "%-20s %-18s %-8s %-8s %-10s\n", "HOSTNAME", "IP", "ONLINE", "EXIT", "OS")
	for _, p := range peers {
		ip := ""
		if len(p.TailscaleIPs) > 0 {
			ip = p.TailscaleIPs[0]
		}
		online := "offline"
		if p.Online {
			online = "online"
		}
		exit := ""
		if p.ExitNode {
			exit = "ACTIVE"
		} else if p.ExitNodeOpt {
			exit = "yes"
		}
		fmt.Fprintf(os.Stdout, "%-20s %-18s %-8s %-8s %-10s\n",
			truncateString(p.HostName, 20), ip, online, exit, truncateString(p.OS, 10))
	}

	return nil
}

func runTailscalePing(cmd *cobra.Command, args []string) error {
	helper := tailscale.NewHelper(tailscaleVerbose, tailscaleDryRun)

	if !helper.IsInstalled() {
		return fmt.Errorf("tailscale is not installed")
	}

	result, err := helper.Ping(args[0])
	if err != nil {
		return err
	}

	if tailscaleDryRun {
		return nil
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(result)
	}

	fmt.Fprintf(os.Stdout, "%s Pong from %s\n", output.Success(""), args[0])
	fmt.Fprintf(os.Stdout, "  IP:      %s\n", result.IP)
	fmt.Fprintf(os.Stdout, "  Via:     %s\n", result.Via)
	fmt.Fprintf(os.Stdout, "  Latency: %.2f ms\n", result.Latency)

	return nil
}

func runTailscaleExitNodes(cmd *cobra.Command, args []string) error {
	helper := tailscale.NewHelper(tailscaleVerbose, tailscaleDryRun)

	if !helper.IsInstalled() {
		return fmt.Errorf("tailscale is not installed")
	}

	exitNodes, err := helper.GetExitNodes()
	if err != nil {
		return err
	}

	ioHelper := ioutils.IO(cmd)
	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(map[string]any{"exit_nodes": exitNodes})
	}

	if len(exitNodes) == 0 {
		fmt.Fprintln(os.Stdout, "No exit nodes available")
		return nil
	}

	fmt.Fprintf(os.Stdout, "%-20s %-18s %-8s %-8s\n", "HOSTNAME", "IP", "ONLINE", "ACTIVE")
	for _, n := range exitNodes {
		ip := ""
		if len(n.TailscaleIPs) > 0 {
			ip = n.TailscaleIPs[0]
		}
		online := "offline"
		if n.Online {
			online = "online"
		}
		active := ""
		if n.ExitNode {
			active = "ACTIVE"
		}
		fmt.Fprintf(os.Stdout, "%-20s %-18s %-8s %-8s\n",
			truncateString(n.HostName, 20), ip, online, active)
	}

	return nil
}

func runTailscaleSetExit(cmd *cobra.Command, args []string) error {
	helper := tailscale.NewHelper(tailscaleVerbose, tailscaleDryRun)

	if !helper.IsInstalled() {
		return fmt.Errorf("tailscale is not installed")
	}

	if err := helper.SetExitNode(args[0]); err != nil {
		return err
	}

	if !tailscaleDryRun {
		fmt.Fprintf(os.Stdout, "%s Set exit node to %s\n", output.Success(""), args[0])
	}
	return nil
}

func runTailscaleDisableExit(cmd *cobra.Command, args []string) error {
	helper := tailscale.NewHelper(tailscaleVerbose, tailscaleDryRun)

	if !helper.IsInstalled() {
		return fmt.Errorf("tailscale is not installed")
	}

	if err := helper.DisableExitNode(); err != nil {
		return err
	}

	if !tailscaleDryRun {
		fmt.Fprintf(os.Stdout, "%s Disabled exit node\n", output.Success(""))
	}
	return nil
}

func runTailscaleUp(cmd *cobra.Command, args []string) error {
	helper := tailscale.NewHelper(tailscaleVerbose, tailscaleDryRun)

	if !helper.IsInstalled() {
		return fmt.Errorf("tailscale is not installed")
	}

	if err := helper.Up(); err != nil {
		return err
	}

	if !tailscaleDryRun {
		fmt.Fprintf(os.Stdout, "%s Tailscale connected\n", output.Success(""))
	}
	return nil
}

func runTailscaleDown(cmd *cobra.Command, args []string) error {
	helper := tailscale.NewHelper(tailscaleVerbose, tailscaleDryRun)

	if !helper.IsInstalled() {
		return fmt.Errorf("tailscale is not installed")
	}

	if err := helper.Down(); err != nil {
		return err
	}

	if !tailscaleDryRun {
		fmt.Fprintf(os.Stdout, "%s Tailscale disconnected\n", output.Success(""))
	}
	return nil
}

func runTailscaleNetcheck(cmd *cobra.Command, args []string) error {
	helper := tailscale.NewHelper(tailscaleVerbose, tailscaleDryRun)

	if !helper.IsInstalled() {
		return fmt.Errorf("tailscale is not installed")
	}

	result, err := helper.NetCheck()
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout, result)
	return nil
}

// truncateString truncates a string to max length
func truncateString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func init() {
	components.Register(&components.Registration{
		Name: "tailscale",
		RegisterCmd: func() *cobra.Command { return tailscaleCmd },
	})
}
