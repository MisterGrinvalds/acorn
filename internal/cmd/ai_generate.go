package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mistergrinvalds/acorn/internal/utils/output"
	"github.com/spf13/cobra"
)

var (
	aiGenerateOutputFormat string
	aiGenerateList         bool
	aiGenerateDryRun       bool
	aiGenerateVerbose      bool
)

// GenerateResult contains the results of config generation.
type GenerateResult struct {
	Target          string         `json:"target"`
	SharedAgents    int            `json:"shared_agents"`
	SharedCommands  int            `json:"shared_commands"`
	ClaudeAgents    int            `json:"claude_agents,omitempty"`
	ClaudeCommands  int            `json:"claude_commands,omitempty"`
	OpenCodeAgents  int            `json:"opencode_agents,omitempty"`
	OpenCodeCommands int           `json:"opencode_commands,omitempty"`
	Generated       bool           `json:"generated"`
	Platforms       []PlatformInfo `json:"platforms,omitempty"`
}

// PlatformInfo contains info about a generated platform.
type PlatformInfo struct {
	Name     string `json:"name"`
	Agents   int    `json:"agents"`
	Commands int    `json:"commands"`
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
	aiCmd.AddCommand(aiGenerateCmd)

	aiGenerateCmd.Flags().StringVarP(&aiGenerateOutputFormat, "output", "o", "table",
		"Output format (table|json|yaml)")
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

	// If --list, just show counts
	if aiGenerateList || aiGenerateDryRun {
		// Count existing generated configs
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

		return outputGenerateResult(result, aiGenerateList)
	}

	// Run generation script
	scriptPath := filepath.Join(saplingDir, "scripts", "generate-configs.sh")
	if _, err := os.Stat(scriptPath); err != nil {
		return fmt.Errorf("generate script not found: %s", scriptPath)
	}

	execCmd := exec.Command(scriptPath, target)
	execCmd.Dir = saplingDir

	if aiGenerateVerbose {
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		if err := execCmd.Run(); err != nil {
			return fmt.Errorf("generation failed: %w", err)
		}
	} else {
		// Capture output and parse counts
		stdout, err := execCmd.StdoutPipe()
		if err != nil {
			return err
		}
		execCmd.Stderr = os.Stderr

		if err := execCmd.Start(); err != nil {
			return fmt.Errorf("failed to start generation: %w", err)
		}

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			// Parse "[INFO] Generated X Claude/OpenCode agents/commands"
			if strings.Contains(line, "Generated") {
				parseGeneratedLine(line, &result)
			}
		}

		if err := execCmd.Wait(); err != nil {
			return fmt.Errorf("generation failed: %w", err)
		}
	}

	result.Generated = true

	// Re-count after generation
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

	return outputGenerateResult(result, false)
}

func outputGenerateResult(result GenerateResult, listOnly bool) error {
	format, err := output.ParseFormat(aiGenerateOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(result)
	}

	// Table format
	if listOnly {
		fmt.Fprintln(os.Stdout, output.Info("Shared Configuration (Source of Truth)"))
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

	if !listOnly {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintf(os.Stdout, "Generated configs include: %s\n",
			output.Info("# Generated by acorn - DO NOT EDIT"))
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
