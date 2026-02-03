// Package tailscale provides Tailscale mesh VPN helper functionality.
package tailscale

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// Status represents Tailscale connection status.
type Status struct {
	Installed    bool   `json:"installed" yaml:"installed"`
	Version      string `json:"version,omitempty" yaml:"version,omitempty"`
	Connected    bool   `json:"connected" yaml:"connected"`
	BackendState string `json:"backend_state,omitempty" yaml:"backend_state,omitempty"`
	IP           string `json:"ip,omitempty" yaml:"ip,omitempty"`
	Hostname     string `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	DNSName      string `json:"dns_name,omitempty" yaml:"dns_name,omitempty"`
	ExitNode     string `json:"exit_node,omitempty" yaml:"exit_node,omitempty"`
	ExitNodeIP   string `json:"exit_node_ip,omitempty" yaml:"exit_node_ip,omitempty"`
	Online       bool   `json:"online" yaml:"online"`
	MagicDNS     bool   `json:"magic_dns" yaml:"magic_dns"`
}

// Peer represents a Tailscale peer/device.
type Peer struct {
	ID           string   `json:"id" yaml:"id"`
	PublicKey    string   `json:"public_key,omitempty" yaml:"public_key,omitempty"`
	HostName     string   `json:"hostname" yaml:"hostname"`
	DNSName      string   `json:"dns_name,omitempty" yaml:"dns_name,omitempty"`
	OS           string   `json:"os,omitempty" yaml:"os,omitempty"`
	TailscaleIPs []string `json:"tailscale_ips" yaml:"tailscale_ips"`
	Online       bool     `json:"online" yaml:"online"`
	ExitNode     bool     `json:"exit_node" yaml:"exit_node"`
	ExitNodeOpt  bool     `json:"exit_node_option" yaml:"exit_node_option"`
	Active       bool     `json:"active" yaml:"active"`
	LastSeen     string   `json:"last_seen,omitempty" yaml:"last_seen,omitempty"`
	Relay        string   `json:"relay,omitempty" yaml:"relay,omitempty"`
}

// PingResult represents the result of a ping operation.
type PingResult struct {
	IP       string  `json:"ip" yaml:"ip"`
	Latency  float64 `json:"latency_ms" yaml:"latency_ms"`
	Via      string  `json:"via,omitempty" yaml:"via,omitempty"`
	Endpoint string  `json:"endpoint,omitempty" yaml:"endpoint,omitempty"`
}

// tailscaleStatus is the JSON structure from tailscale status --json.
type tailscaleStatus struct {
	BackendState string                    `json:"BackendState"`
	Self         tailscaleSelf             `json:"Self"`
	MagicDNSSuffix string                  `json:"MagicDNSSuffix"`
	Peer         map[string]tailscalePeer  `json:"Peer"`
	ExitNodeStatus *exitNodeStatus         `json:"ExitNodeStatus,omitempty"`
}

type tailscaleSelf struct {
	ID           string   `json:"ID"`
	PublicKey    string   `json:"PublicKey"`
	HostName     string   `json:"HostName"`
	DNSName      string   `json:"DNSName"`
	OS           string   `json:"OS"`
	TailscaleIPs []string `json:"TailscaleIPs"`
	Online       bool     `json:"Online"`
}

type tailscalePeer struct {
	ID             string   `json:"ID"`
	PublicKey      string   `json:"PublicKey"`
	HostName       string   `json:"HostName"`
	DNSName        string   `json:"DNSName"`
	OS             string   `json:"OS"`
	TailscaleIPs   []string `json:"TailscaleIPs"`
	Online         bool     `json:"Online"`
	ExitNode       bool     `json:"ExitNode"`
	ExitNodeOption bool     `json:"ExitNodeOption"`
	Active         bool     `json:"Active"`
	LastSeen       string   `json:"LastSeen"`
	Relay          string   `json:"Relay"`
}

type exitNodeStatus struct {
	ID        string `json:"ID"`
	Online    bool   `json:"Online"`
	TailscaleIPs []string `json:"TailscaleIPs"`
}

// Helper provides Tailscale helper operations.
type Helper struct {
	verbose bool
	dryRun  bool
}

// NewHelper creates a new Tailscale Helper.
func NewHelper(verbose, dryRun bool) *Helper {
	return &Helper{
		verbose: verbose,
		dryRun:  dryRun,
	}
}

// IsInstalled checks if Tailscale is installed.
func (h *Helper) IsInstalled() bool {
	_, err := exec.LookPath("tailscale")
	return err == nil
}

// GetVersion returns the Tailscale version.
func (h *Helper) GetVersion() string {
	cmd := exec.Command("tailscale", "version")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}
	return ""
}

// GetStatus returns Tailscale connection status.
func (h *Helper) GetStatus() *Status {
	status := &Status{}

	if _, err := exec.LookPath("tailscale"); err != nil {
		return status
	}
	status.Installed = true
	status.Version = h.GetVersion()

	// Get JSON status
	cmd := exec.Command("tailscale", "status", "--json")
	out, err := cmd.Output()
	if err != nil {
		return status
	}

	var tsStatus tailscaleStatus
	if err := json.Unmarshal(out, &tsStatus); err != nil {
		return status
	}

	status.BackendState = tsStatus.BackendState
	status.Connected = tsStatus.BackendState == "Running"
	status.Online = tsStatus.Self.Online
	status.Hostname = tsStatus.Self.HostName
	status.DNSName = tsStatus.Self.DNSName
	status.MagicDNS = tsStatus.MagicDNSSuffix != ""

	if len(tsStatus.Self.TailscaleIPs) > 0 {
		status.IP = tsStatus.Self.TailscaleIPs[0]
	}

	// Check for active exit node
	if tsStatus.ExitNodeStatus != nil {
		status.ExitNode = tsStatus.ExitNodeStatus.ID
		if len(tsStatus.ExitNodeStatus.TailscaleIPs) > 0 {
			status.ExitNodeIP = tsStatus.ExitNodeStatus.TailscaleIPs[0]
		}
	}

	return status
}

// GetPeers returns list of all peers.
func (h *Helper) GetPeers() ([]Peer, error) {
	cmd := exec.Command("tailscale", "status", "--json")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get tailscale status: %w", err)
	}

	var tsStatus tailscaleStatus
	if err := json.Unmarshal(out, &tsStatus); err != nil {
		return nil, fmt.Errorf("failed to parse tailscale status: %w", err)
	}

	var peers []Peer
	for _, p := range tsStatus.Peer {
		peer := Peer{
			ID:           p.ID,
			PublicKey:    p.PublicKey,
			HostName:     p.HostName,
			DNSName:      p.DNSName,
			OS:           p.OS,
			TailscaleIPs: p.TailscaleIPs,
			Online:       p.Online,
			ExitNode:     p.ExitNode,
			ExitNodeOpt:  p.ExitNodeOption,
			Active:       p.Active,
			LastSeen:     p.LastSeen,
			Relay:        p.Relay,
		}
		peers = append(peers, peer)
	}

	return peers, nil
}

// GetOnlinePeers returns list of online peers only.
func (h *Helper) GetOnlinePeers() ([]Peer, error) {
	peers, err := h.GetPeers()
	if err != nil {
		return nil, err
	}

	var online []Peer
	for _, p := range peers {
		if p.Online {
			online = append(online, p)
		}
	}

	return online, nil
}

// GetExitNodes returns list of peers that can be used as exit nodes.
func (h *Helper) GetExitNodes() ([]Peer, error) {
	peers, err := h.GetPeers()
	if err != nil {
		return nil, err
	}

	var exitNodes []Peer
	for _, p := range peers {
		if p.ExitNodeOpt {
			exitNodes = append(exitNodes, p)
		}
	}

	return exitNodes, nil
}

// Ping pings a peer and returns the result.
func (h *Helper) Ping(peer string) (*PingResult, error) {
	if h.dryRun {
		fmt.Printf("[dry-run] would run: tailscale ping %s\n", peer)
		return nil, nil
	}

	// Use --c 1 for single ping
	cmd := exec.Command("tailscale", "ping", "-c", "1", peer)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to ping %s: %w", peer, err)
	}

	// Parse output like: "pong from hostname (100.x.x.x) via DERP(xxx) in 50ms"
	// or: "pong from hostname (100.x.x.x) via 192.168.1.1:41641 in 5ms"
	output := strings.TrimSpace(string(out))
	result := &PingResult{}

	// Extract IP from parentheses
	if start := strings.Index(output, "("); start != -1 {
		if end := strings.Index(output[start:], ")"); end != -1 {
			result.IP = output[start+1 : start+end]
		}
	}

	// Extract via
	if viaIdx := strings.Index(output, "via "); viaIdx != -1 {
		rest := output[viaIdx+4:]
		if inIdx := strings.Index(rest, " in "); inIdx != -1 {
			result.Via = rest[:inIdx]
		}
	}

	// Extract latency
	if inIdx := strings.Index(output, " in "); inIdx != -1 {
		latencyStr := output[inIdx+4:]
		latencyStr = strings.TrimSuffix(latencyStr, "ms")
		var latency float64
		fmt.Sscanf(latencyStr, "%f", &latency)
		result.Latency = latency
	}

	return result, nil
}

// SetExitNode sets the exit node.
func (h *Helper) SetExitNode(node string) error {
	if h.dryRun {
		fmt.Printf("[dry-run] would run: tailscale set --exit-node=%s\n", node)
		return nil
	}

	cmd := exec.Command("tailscale", "set", "--exit-node="+node)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set exit node: %w", err)
	}

	return nil
}

// DisableExitNode disables the exit node.
func (h *Helper) DisableExitNode() error {
	if h.dryRun {
		fmt.Printf("[dry-run] would run: tailscale set --exit-node=\n")
		return nil
	}

	cmd := exec.Command("tailscale", "set", "--exit-node=")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to disable exit node: %w", err)
	}

	return nil
}

// Up brings Tailscale up.
func (h *Helper) Up() error {
	if h.dryRun {
		fmt.Println("[dry-run] would run: tailscale up")
		return nil
	}

	cmd := exec.Command("tailscale", "up")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to bring tailscale up: %w", err)
	}

	return nil
}

// Down brings Tailscale down.
func (h *Helper) Down() error {
	if h.dryRun {
		fmt.Println("[dry-run] would run: tailscale down")
		return nil
	}

	cmd := exec.Command("tailscale", "down")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to bring tailscale down: %w", err)
	}

	return nil
}

// NetCheck runs tailscale netcheck and returns raw output.
func (h *Helper) NetCheck() (string, error) {
	cmd := exec.Command("tailscale", "netcheck")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("netcheck failed: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}
