// Package claude provides Claude Code helper functionality.
package claude

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Paths holds all Claude configuration paths.
type Paths struct {
	ClaudeDir    string `json:"claude_dir" yaml:"claude_dir"`
	Config       string `json:"config" yaml:"config"`
	Settings     string `json:"settings" yaml:"settings"`
	Local        string `json:"local" yaml:"local"`
	StatsCache   string `json:"stats_cache" yaml:"stats_cache"`
	ProjectsDir  string `json:"projects_dir" yaml:"projects_dir"`
	CommandsDir  string `json:"commands_dir" yaml:"commands_dir"`
	AgentsDir    string `json:"agents_dir" yaml:"agents_dir"`
	SubagentsDir string `json:"subagents_dir" yaml:"subagents_dir"`
}

// Info contains Claude Code information summary.
type Info struct {
	Version       string   `json:"version" yaml:"version"`
	ClaudeDir     string   `json:"claude_dir" yaml:"claude_dir"`
	ConfigExists  bool     `json:"config_exists" yaml:"config_exists"`
	SettingsExist bool     `json:"settings_exist" yaml:"settings_exist"`
	LocalExists   bool     `json:"local_exists" yaml:"local_exists"`
	StatsExist    bool     `json:"stats_exist" yaml:"stats_exist"`
	TotalSessions int      `json:"total_sessions,omitempty" yaml:"total_sessions,omitempty"`
	TotalMessages int      `json:"total_messages,omitempty" yaml:"total_messages,omitempty"`
	AgentCount    int      `json:"agent_count" yaml:"agent_count"`
	CommandCount  int      `json:"command_count" yaml:"command_count"`
	Errors        []string `json:"errors,omitempty" yaml:"errors,omitempty"`
}

// Helper provides Claude Code helper operations.
type Helper struct {
	verbose bool
	dryRun  bool
	paths   *Paths
}

// NewHelper creates a new Helper with standard paths.
func NewHelper(verbose, dryRun bool) *Helper {
	home, _ := os.UserHomeDir()
	claudeDir := filepath.Join(home, ".claude")

	return &Helper{
		verbose: verbose,
		dryRun:  dryRun,
		paths: &Paths{
			ClaudeDir:    claudeDir,
			Config:       filepath.Join(home, ".claude.json"),
			Settings:     filepath.Join(claudeDir, "settings.json"),
			Local:        filepath.Join(claudeDir, "settings.local.json"),
			StatsCache:   filepath.Join(claudeDir, "stats-cache.json"),
			ProjectsDir:  filepath.Join(claudeDir, "projects"),
			CommandsDir:  filepath.Join(claudeDir, "commands"),
			AgentsDir:    filepath.Join(claudeDir, "agents"),
			SubagentsDir: filepath.Join(claudeDir, "subagents"),
		},
	}
}

// GetPaths returns the configuration paths.
func (h *Helper) GetPaths() *Paths {
	return h.paths
}

// FileExists checks if a file exists.
func (h *Helper) FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// DirExists checks if a directory exists.
func (h *Helper) DirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// GetVersion returns the Claude CLI version.
func (h *Helper) GetVersion() (string, error) {
	cmd := exec.Command("claude", "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("claude not found in PATH: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// CountFiles counts files with given extension in a directory.
func (h *Helper) CountFiles(dir, ext string) int {
	if !h.DirExists(dir) {
		return 0
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	count := 0
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ext) {
			count++
		}
	}
	return count
}

// GetInfo returns Claude Code information summary.
func (h *Helper) GetInfo() (*Info, error) {
	info := &Info{
		ClaudeDir:     h.paths.ClaudeDir,
		ConfigExists:  h.FileExists(h.paths.Config),
		SettingsExist: h.FileExists(h.paths.Settings),
		LocalExists:   h.FileExists(h.paths.Local),
		StatsExist:    h.FileExists(h.paths.StatsCache),
		AgentCount:    h.CountFiles(h.paths.AgentsDir, ".md"),
		CommandCount:  h.CountFiles(h.paths.CommandsDir, ".md"),
	}

	// Get version
	version, err := h.GetVersion()
	if err != nil {
		info.Version = "not found"
		info.Errors = append(info.Errors, err.Error())
	} else {
		info.Version = version
	}

	// Get quick stats if stats file exists
	if info.StatsExist {
		if stats, err := h.readStatsQuick(); err == nil {
			info.TotalSessions = stats.TotalSessions
			info.TotalMessages = stats.TotalMessages
		}
	}

	return info, nil
}

// quickStats is a minimal struct for reading just session/message counts.
type quickStats struct {
	TotalSessions int `json:"totalSessions"`
	TotalMessages int `json:"totalMessages"`
}

// readStatsQuick reads just the session/message counts from stats file.
func (h *Helper) readStatsQuick() (*quickStats, error) {
	data, err := os.ReadFile(h.paths.StatsCache)
	if err != nil {
		return nil, err
	}
	var stats quickStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

// ReadJSONFile reads and parses a JSON file into the given target.
func (h *Helper) ReadJSONFile(path string, target interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

// WriteJSONFile writes the given data to a JSON file atomically.
func (h *Helper) WriteJSONFile(path string, data interface{}) error {
	if h.dryRun {
		if h.verbose {
			fmt.Printf("[dry-run] Would write to %s\n", path)
		}
		return nil
	}

	// Marshal with indentation
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Write to temp file first
	dir := filepath.Dir(path)
	tmpFile, err := os.CreateTemp(dir, "claude-*.json")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()

	if _, err := tmpFile.Write(jsonData); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return err
	}
	tmpFile.Close()

	// Atomic rename
	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return err
	}

	return nil
}
