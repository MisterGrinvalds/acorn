package cmd

import (
	"github.com/mistergrinvalds/acorn/internal/components"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	ioutils "github.com/mistergrinvalds/acorn/internal/utils/io"
	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	aiGenerateList    bool
	aiGenerateDryRun  bool
	aiGenerateVerbose bool
)

// GenerateResult contains the results of config generation.
type GenerateResult struct {
	Target           string         `json:"target"`
	SharedAgents     int            `json:"shared_agents"`
	SharedCommands   int            `json:"shared_commands"`
	ClaudeAgents     int            `json:"claude_agents,omitempty"`
	ClaudeCommands   int            `json:"claude_commands,omitempty"`
	OpenCodeAgents   int            `json:"opencode_agents,omitempty"`
	OpenCodeCommands int            `json:"opencode_commands,omitempty"`
	Generated        bool           `json:"generated"`
	DryRun           bool           `json:"dry_run,omitempty"`
	Platforms        []PlatformInfo `json:"platforms,omitempty"`
	Files            []string       `json:"files,omitempty"`
}

// PlatformInfo contains info about a generated platform.
type PlatformInfo struct {
	Name     string   `json:"name"`
	Agents   int      `json:"agents"`
	Commands int      `json:"commands"`
	Files    []string `json:"files,omitempty"`
}

// aiGenerateCmd generates platform configs from shared sources
var aiGenerateCmd = &cobra.Command{
	Use:   "generate [claude|opencode|all]",
	Short: "Generate platform configs from shared sources",
	Long: `Generate Claude Code and OpenCode configurations from shared sources.

The shared configuration in .sapling/config/shared/ is the single source
of truth. This command generates platform-specific configs for Claude Code
and OpenCode from that shared source.

Targets:
  claude    Generate Claude Code configurations only
  opencode  Generate OpenCode configurations only
  all       Generate both (default)

Examples:
  acorn ai generate                    # Generate all configs
  acorn ai generate claude             # Generate Claude configs only
  acorn ai generate opencode           # Generate OpenCode configs only
  acorn ai generate --list             # Show counts without generating
  acorn ai generate -o json            # JSON output`,
	Args:      cobra.MaximumNArgs(1),
	ValidArgs: []string{"claude", "opencode", "all"},
	RunE:      runAIGenerate,
}

func init() {

	aiGenerateCmd.Flags().BoolVarP(&aiGenerateList, "list", "l", false,
		"Show counts without generating")
	aiGenerateCmd.Flags().BoolVar(&aiGenerateDryRun, "dry-run", false,
		"Show what would be done without executing")
	aiGenerateCmd.Flags().BoolVarP(&aiGenerateVerbose, "verbose", "v", false,
		"Show verbose output")
}

func runAIGenerate(cmd *cobra.Command, args []string) error {
	target := "all"
	if len(args) > 0 {
		target = args[0]
		if target != "claude" && target != "opencode" && target != "all" {
			return fmt.Errorf("invalid target: %s (must be claude, opencode, or all)", target)
		}
	}

	// Find .sapling directory
	saplingDir, err := findSaplingDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(saplingDir, "config")
	sharedDir := filepath.Join(configDir, "shared")
	claudeDir := filepath.Join(configDir, "claude")
	openCodeDir := filepath.Join(configDir, "opencode")

	// Count shared configs
	sharedAgents, err := countDirs(filepath.Join(sharedDir, "agents"))
	if err != nil {
		return fmt.Errorf("failed to count shared agents: %w", err)
	}

	sharedCommands, err := countFiles(filepath.Join(sharedDir, "commands"), "meta.yaml")
	if err != nil {
		return fmt.Errorf("failed to count shared commands: %w", err)
	}

	result := GenerateResult{
		Target:         target,
		SharedAgents:   sharedAgents,
		SharedCommands: sharedCommands,
		Generated:      false,
	}

	// If --list, just show counts (no script run)
	if aiGenerateList {
		if target == "claude" || target == "all" {
			result.ClaudeAgents, _ = countFiles(filepath.Join(claudeDir, "agents"), "*.md")
			result.ClaudeCommands, _ = countFiles(filepath.Join(claudeDir, "commands"), "*.md")
			result.Platforms = append(result.Platforms, PlatformInfo{
				Name:     "claude",
				Agents:   result.ClaudeAgents,
				Commands: result.ClaudeCommands,
			})
		}
		if target == "opencode" || target == "all" {
			result.OpenCodeAgents, _ = countFiles(filepath.Join(openCodeDir, "agents"), "*.md")
			result.OpenCodeCommands, _ = countFiles(filepath.Join(openCodeDir, "commands"), "*.md")
			result.Platforms = append(result.Platforms, PlatformInfo{
				Name:     "opencode",
				Agents:   result.OpenCodeAgents,
				Commands: result.OpenCodeCommands,
			})
		}

		return outputGenerateResult(cmd, result, true)
	}

	// If --dry-run, run script with ACORN_DRY_RUN=1 to get file list without writing
	if aiGenerateDryRun {
		result.DryRun = true
		files, err := runGenerateScript(saplingDir, target, true)
		if err != nil {
			return err
		}
		result.Files = files
		result.Platforms = buildPlatformInfo(files, &result)
		return outputGenerateResult(cmd, result, false)
	}

	// Run generation script
	files, err := runGenerateScript(saplingDir, target, false)
	if err != nil {
		return err
	}

	result.Generated = true
	result.Files = files
	result.Platforms = buildPlatformInfo(files, &result)

	return outputGenerateResult(cmd, result, false)
}

// runGenerateScript executes the generate-configs.sh script and returns the list of files
// it wrote (or would write in dry-run mode).
func runGenerateScript(saplingDir, target string, dryRun bool) ([]string, error) {
	scriptPath := filepath.Join(saplingDir, "scripts", "generate-configs.sh")
	if _, err := os.Stat(scriptPath); err != nil {
		return nil, fmt.Errorf("generate script not found: %s", scriptPath)
	}

	execCmd := exec.Command(scriptPath, target)
	execCmd.Dir = saplingDir

	if dryRun {
		execCmd.Env = append(os.Environ(), "ACORN_DRY_RUN=1")
	}

	stdout, err := execCmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	execCmd.Stderr = os.Stderr

	if err := execCmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start generation: %w", err)
	}

	var files []string
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		// Parse [FILE] lines for file paths
		if strings.Contains(line, "[FILE] would-write:") {
			path := strings.TrimSpace(strings.SplitN(line, "would-write:", 2)[1])
			files = append(files, path)
		} else if strings.Contains(line, "[FILE] wrote:") {
			path := strings.TrimSpace(strings.SplitN(line, "wrote:", 2)[1])
			files = append(files, path)
		}
		// Pass through verbose output
		if aiGenerateVerbose {
			fmt.Fprintln(os.Stdout, line)
		}
	}

	if err := execCmd.Wait(); err != nil {
		return nil, fmt.Errorf("generation failed: %w", err)
	}

	sort.Strings(files)
	return files, nil
}

// buildPlatformInfo groups file paths by platform and populates result counts.
func buildPlatformInfo(files []string, result *GenerateResult) []PlatformInfo {
	platformFiles := make(map[string][]string)
	platformAgents := make(map[string]int)
	platformCommands := make(map[string]int)

	for _, f := range files {
		parts := strings.SplitN(f, "/", 3)
		if len(parts) < 2 {
			continue
		}
		platform := parts[0] // "claude" or "opencode"
		kind := parts[1]     // "agents" or "commands"

		platformFiles[platform] = append(platformFiles[platform], f)
		switch kind {
		case "agents":
			platformAgents[platform]++
		case "commands":
			platformCommands[platform]++
		}
	}

	// Set result counts
	result.ClaudeAgents = platformAgents["claude"]
	result.ClaudeCommands = platformCommands["claude"]
	result.OpenCodeAgents = platformAgents["opencode"]
	result.OpenCodeCommands = platformCommands["opencode"]

	var platforms []PlatformInfo
	for _, name := range []string{"claude", "opencode"} {
		if f, ok := platformFiles[name]; ok {
			platforms = append(platforms, PlatformInfo{
				Name:     name,
				Agents:   platformAgents[name],
				Commands: platformCommands[name],
				Files:    f,
			})
		}
	}
	return platforms
}

func outputGenerateResult(cmd *cobra.Command, result GenerateResult, listOnly bool) error {
	ioHelper := ioutils.IO(cmd)

	if ioHelper.IsStructured() {
		return ioHelper.WriteOutput(result)
	}

	// Table format
	if listOnly {
		fmt.Fprintln(os.Stdout, output.Info("Shared Configuration (Source of Truth)"))
	} else if result.DryRun {
		fmt.Fprintln(os.Stdout, output.Warning("Dry Run — would generate:"))
	} else {
		fmt.Fprintln(os.Stdout, output.Success("Generation Complete"))
	}
	fmt.Fprintln(os.Stdout)

	table := output.NewTable("SOURCE", "AGENTS", "COMMANDS")
	table.AddRow("shared", strconv.Itoa(result.SharedAgents), strconv.Itoa(result.SharedCommands))
	table.Render(os.Stdout)

	if len(result.Platforms) > 0 {
		fmt.Fprintln(os.Stdout)
		if listOnly {
			fmt.Fprintln(os.Stdout, output.Info("Generated Platforms"))
		} else {
			fmt.Fprintln(os.Stdout, output.Info("Generated Output"))
		}
		fmt.Fprintln(os.Stdout)

		table2 := output.NewTable("PLATFORM", "AGENTS", "COMMANDS", "STATUS")
		for _, p := range result.Platforms {
			status := output.Success("synced")
			if p.Agents != result.SharedAgents || p.Commands != result.SharedCommands {
				status = output.Warning("drift")
			}
			table2.AddRow(p.Name, strconv.Itoa(p.Agents), strconv.Itoa(p.Commands), status)
		}
		table2.Render(os.Stdout)
	}

	// Print file list for dry-run and normal generation
	if !listOnly && len(result.Platforms) > 0 {
		for _, p := range result.Platforms {
			if len(p.Files) > 0 {
				fmt.Fprintln(os.Stdout)
				if result.DryRun {
					fmt.Fprintf(os.Stdout, "  %s files in scope:\n", output.Info(p.Name))
				} else {
					fmt.Fprintf(os.Stdout, "  %s files changed:\n", output.Info(p.Name))
				}
				for _, f := range p.Files {
					fmt.Fprintf(os.Stdout, "    %s\n", f)
				}
			}
		}
	}

	return nil
}

func parseGeneratedLine(line string, result *GenerateResult) {
	// Parse lines like "[INFO] Generated 33 Claude agents"
	parts := strings.Fields(line)
	for i, part := range parts {
		if part == "Generated" && i+2 < len(parts) {
			count, err := strconv.Atoi(parts[i+1])
			if err != nil {
				continue
			}
			what := strings.ToLower(parts[i+2])
			if strings.Contains(what, "claude") {
				if strings.Contains(line, "agents") {
					result.ClaudeAgents = count
				} else if strings.Contains(line, "commands") {
					result.ClaudeCommands = count
				}
			} else if strings.Contains(what, "opencode") {
				if strings.Contains(line, "agents") {
					result.OpenCodeAgents = count
				} else if strings.Contains(line, "commands") {
					result.OpenCodeCommands = count
				}
			}
		}
	}
}

func findSaplingDir() (string, error) {
	// Check DOTFILES_ROOT first
	if root := os.Getenv("DOTFILES_ROOT"); root != "" {
		saplingDir := filepath.Join(root, ".sapling")
		if _, err := os.Stat(saplingDir); err == nil {
			return saplingDir, nil
		}
	}

	// Check config
	if cfg != nil && cfg.DotfilesRoot != "" {
		saplingDir := filepath.Join(cfg.DotfilesRoot, ".sapling")
		if _, err := os.Stat(saplingDir); err == nil {
			return saplingDir, nil
		}
	}

	// Check current directory and parents
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir := cwd
	for range 10 {
		saplingDir := filepath.Join(dir, ".sapling")
		if _, err := os.Stat(saplingDir); err == nil {
			return saplingDir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf(".sapling directory not found")
}

func countDirs(path string) (int, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			count++
		}
	}
	return count, nil
}

func countFiles(path string, pattern string) (int, error) {
	// Recursive search for pattern
	count := 0
	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}
		if info.IsDir() {
			return nil
		}

		// Match pattern
		if pattern == "*.md" && strings.HasSuffix(info.Name(), ".md") {
			count++
		} else if info.Name() == pattern {
			count++
		}
		return nil
	})
	return count, err
}

func init() {
	components.Register(&components.Registration{
		Name: "ai-generate",
		RegisterCmd: func() *cobra.Command { return aiGenerateCmd },
	})
}
