package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mistergrinvalds/acorn/internal/claude"
	"github.com/mistergrinvalds/acorn/internal/installer"
	"github.com/mistergrinvalds/acorn/internal/output"
	"github.com/spf13/cobra"
)

var (
	claudeOutputFormat string
	claudeDryRun       bool
	claudeVerbose      bool
)

// claudeCmd represents the claude command group
var claudeCmd = &cobra.Command{
	Use:   "claude",
	Short: "Claude Code management and utilities",
	Long: `Helpers for managing Claude Code configuration, statistics, and settings.

Provides commands for viewing usage stats, managing permissions,
editing settings, and aggregating agents/commands from repositories.

Examples:
  acorn claude info                    # Show Claude Code info
  acorn claude stats                   # View usage statistics
  acorn claude stats tokens            # View token usage by model
  acorn claude permissions             # View permissions
  acorn claude settings global         # View global settings
  acorn claude projects                # List projects`,
	Aliases: []string{"cc"},
}

// claudeInfoCmd shows Claude Code information
var claudeInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show Claude Code information summary",
	Long: `Display Claude Code version, configuration status, and quick stats.

Shows:
  - Claude CLI version
  - Configuration file status
  - Quick usage statistics
  - Agent and command counts

Examples:
  acorn claude info
  acorn claude info -o json`,
	RunE: runClaudeInfo,
}

// claudeStatsCmd shows usage statistics
var claudeStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "View usage statistics",
	Long: `Display Claude Code usage statistics including session counts,
message counts, and model usage breakdown.

Examples:
  acorn claude stats
  acorn claude stats -o json`,
	RunE: runClaudeStats,
}

// claudeStatsTokensCmd shows token usage
var claudeStatsTokensCmd = &cobra.Command{
	Use:   "tokens",
	Short: "View token usage by model",
	Long: `Display token usage breakdown by model including input,
output, and cache tokens.

Examples:
  acorn claude stats tokens
  acorn claude stats tokens -o json`,
	RunE: runClaudeStatsTokens,
}

// claudeStatsDailyCmd shows daily usage
var claudeStatsDailyCmd = &cobra.Command{
	Use:   "daily [days]",
	Short: "View daily token usage",
	Long: `Display daily token usage for the last N days (default: 7).

Examples:
  acorn claude stats daily       # Last 7 days
  acorn claude stats daily 14    # Last 14 days`,
	Args: cobra.MaximumNArgs(1),
	RunE: runClaudeStatsDaily,
}

// claudePermissionsCmd shows permissions
var claudePermissionsCmd = &cobra.Command{
	Use:   "permissions",
	Short: "View and manage permissions",
	Long: `Display current permission rules from settings.local.json.

Shows allowed and denied patterns.

Examples:
  acorn claude permissions
  acorn claude permissions -o json`,
	Aliases: []string{"perms"},
	RunE:    runClaudePermissions,
}

// claudePermissionsAddCmd adds a permission
var claudePermissionsAddCmd = &cobra.Command{
	Use:   "add <rule> [allow|deny]",
	Short: "Add a permission rule",
	Long: `Add a new permission rule to the allow or deny list.

Examples:
  acorn claude permissions add "Bash(npm:*)"
  acorn claude permissions add "Bash(rm:*)" deny`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runClaudePermissionsAdd,
}

// claudePermissionsRemoveCmd removes a permission
var claudePermissionsRemoveCmd = &cobra.Command{
	Use:   "remove <rule> [allow|deny]",
	Short: "Remove a permission rule",
	Long: `Remove a permission rule from the allow or deny list.

Examples:
  acorn claude permissions remove "Bash(npm:*)"
  acorn claude permissions remove "Bash(rm:*)" deny`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runClaudePermissionsRemove,
}

// claudeSettingsCmd shows settings
var claudeSettingsCmd = &cobra.Command{
	Use:   "settings [global|local|config]",
	Short: "View settings files",
	Long: `Display Claude Code settings from various configuration files.

Files:
  global - ~/.claude/settings.json (default)
  local  - ~/.claude/settings.local.json
  config - ~/.claude.json

Examples:
  acorn claude settings           # Show global settings
  acorn claude settings local     # Show local settings
  acorn claude settings config    # Show main config`,
	Args: cobra.MaximumNArgs(1),
	RunE: runClaudeSettings,
}

// claudeSettingsEditCmd edits settings
var claudeSettingsEditCmd = &cobra.Command{
	Use:   "edit [global|local|config]",
	Short: "Edit settings file with $EDITOR",
	Long: `Open a settings file in your default editor.

Examples:
  acorn claude settings edit           # Edit global settings
  acorn claude settings edit local     # Edit local settings`,
	Args: cobra.MaximumNArgs(1),
	RunE: runClaudeSettingsEdit,
}

// claudeProjectsCmd lists projects
var claudeProjectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "List Claude projects",
	Long: `Display all Claude Code projects with their costs.

Examples:
  acorn claude projects
  acorn claude projects -o json`,
	RunE: runClaudeProjects,
}

// claudeMcpCmd shows MCP servers
var claudeMcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "View and manage MCP servers",
	Long: `Display MCP servers from configuration.

Examples:
  acorn claude mcp
  acorn claude mcp -o json`,
	RunE: runClaudeMcp,
}

// claudeMcpAddCmd adds an MCP server
var claudeMcpAddCmd = &cobra.Command{
	Use:   "add <name> <url> [type]",
	Short: "Add MCP server to .mcp.json",
	Long: `Add a new MCP server to the local .mcp.json configuration.

Examples:
  acorn claude mcp add myserver http://localhost:8080
  acorn claude mcp add myserver http://localhost:8080 http`,
	Args: cobra.RangeArgs(2, 3),
	RunE: runClaudeMcpAdd,
}

// claudeCommandsCmd lists custom commands
var claudeCommandsCmd = &cobra.Command{
	Use:   "commands",
	Short: "List custom commands",
	Long: `Display custom commands from user and project directories.

Examples:
  acorn claude commands
  acorn claude commands -o json`,
	RunE: runClaudeCommands,
}

// claudeAggregateCmd aggregates agents/commands
var claudeAggregateCmd = &cobra.Command{
	Use:   "aggregate [search-dir]",
	Short: "Aggregate agents/commands from repositories",
	Long: `Scan repositories for .claude directories and aggregate
agents, commands, and subagents to the dotfiles config.

Handles file deduplication and renames conflicting files
with repository prefixes.

Examples:
  acorn claude aggregate              # Scan ~/Repos
  acorn claude aggregate ~/Projects   # Scan custom directory`,
	Args: cobra.MaximumNArgs(1),
	RunE: runClaudeAggregate,
}

// claudeAggregateListCmd lists aggregated items
var claudeAggregateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all agents, commands, and subagents",
	Long: `Display all aggregated agents, commands, and subagents
from the dotfiles configuration.

Examples:
  acorn claude aggregate list
  acorn claude aggregate list -o json`,
	RunE: runClaudeAggregateList,
}

// claudeClearCmd clears cache/stats
var claudeClearCmd = &cobra.Command{
	Use:   "clear [cache|stats]",
	Short: "Clear cache or stats",
	Long: `Clear Claude Code cache directories or stats file.

Options:
  cache - Remove shell-snapshots and debug directories (default)
  stats - Remove stats-cache.json

Examples:
  acorn claude clear          # Clear cache
  acorn claude clear stats    # Clear stats`,
	Args: cobra.MaximumNArgs(1),
	RunE: runClaudeClear,
}

// claudeHelpCmd shows function help
var claudeHelpCmd = &cobra.Command{
	Use:   "help",
	Short: "Show Claude Code shell function help",
	Long: `Display help for all Claude Code shell functions and aliases.

Examples:
  acorn claude help`,
	RunE: runClaudeHelp,
}

// claudeInstallCmd installs Claude Code CLI
var claudeInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Claude Code CLI",
	Long: `Install the Claude Code CLI tool.

Installs claude via npm (requires Node.js).

Examples:
  acorn claude install           # Install Claude Code CLI
  acorn claude install --dry-run # Show what would be installed
  acorn claude install -v        # Verbose output`,
	RunE: runClaudeInstall,
}

func init() {
	rootCmd.AddCommand(claudeCmd)

	// Add subcommands
	claudeCmd.AddCommand(claudeInfoCmd)
	claudeCmd.AddCommand(claudeInstallCmd)
	claudeCmd.AddCommand(claudeStatsCmd)
	claudeCmd.AddCommand(claudePermissionsCmd)
	claudeCmd.AddCommand(claudeSettingsCmd)
	claudeCmd.AddCommand(claudeProjectsCmd)
	claudeCmd.AddCommand(claudeMcpCmd)
	claudeCmd.AddCommand(claudeCommandsCmd)
	claudeCmd.AddCommand(claudeAggregateCmd)
	claudeCmd.AddCommand(claudeClearCmd)
	claudeCmd.AddCommand(claudeHelpCmd)

	// Stats subcommands
	claudeStatsCmd.AddCommand(claudeStatsTokensCmd)
	claudeStatsCmd.AddCommand(claudeStatsDailyCmd)

	// Permissions subcommands
	claudePermissionsCmd.AddCommand(claudePermissionsAddCmd)
	claudePermissionsCmd.AddCommand(claudePermissionsRemoveCmd)

	// Settings subcommands
	claudeSettingsCmd.AddCommand(claudeSettingsEditCmd)

	// MCP subcommands
	claudeMcpCmd.AddCommand(claudeMcpAddCmd)

	// Aggregate subcommands
	claudeAggregateCmd.AddCommand(claudeAggregateListCmd)

	// Persistent flags
	claudeCmd.PersistentFlags().StringVarP(&claudeOutputFormat, "output", "o", "table",
		"Output format (table|json|yaml)")
	claudeCmd.PersistentFlags().BoolVar(&claudeDryRun, "dry-run", false,
		"Show what would be done without executing")
	claudeCmd.PersistentFlags().BoolVarP(&claudeVerbose, "verbose", "v", false,
		"Show verbose output")
}

// runClaudeInfo displays Claude Code information
func runClaudeInfo(cmd *cobra.Command, args []string) error {
	helper := claude.NewHelper(claudeVerbose, claudeDryRun)
	info, err := helper.GetInfo()
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(claudeOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(info)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n\n", output.Info("Claude Code Information"))

	fmt.Fprintf(os.Stdout, "Version: %s\n\n", output.Success(info.Version))

	fmt.Fprintf(os.Stdout, "Configuration Files:\n")
	printFileStatus("~/.claude.json", info.ConfigExists)
	printFileStatus("settings.json", info.SettingsExist)
	printFileStatus("settings.local.json", info.LocalExists)
	printFileStatus("stats-cache.json", info.StatsExist)
	fmt.Println()

	if info.StatsExist {
		fmt.Fprintf(os.Stdout, "Quick Stats:\n")
		fmt.Fprintf(os.Stdout, "  Sessions: %d\n", info.TotalSessions)
		fmt.Fprintf(os.Stdout, "  Messages: %d\n", info.TotalMessages)
		fmt.Println()
	}

	fmt.Fprintf(os.Stdout, "Assets:\n")
	fmt.Fprintf(os.Stdout, "  Agents:   %d\n", info.AgentCount)
	fmt.Fprintf(os.Stdout, "  Commands: %d\n", info.CommandCount)

	return nil
}

func printFileStatus(name string, exists bool) {
	if exists {
		fmt.Fprintf(os.Stdout, "  [%s] %s\n", output.Success("x"), name)
	} else {
		fmt.Fprintf(os.Stdout, "  [ ] %s\n", name)
	}
}

// Stats command handlers

func runClaudeStats(cmd *cobra.Command, args []string) error {
	helper := claude.NewHelper(claudeVerbose, claudeDryRun)
	summary, err := helper.GetStatsSummary()
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(claudeOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(summary)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n\n", output.Info("Claude Code Usage Statistics"))

	fmt.Fprintf(os.Stdout, "Total Sessions: %s\n", output.Success(fmt.Sprintf("%d", summary.TotalSessions)))
	fmt.Fprintf(os.Stdout, "Total Messages: %s\n\n", output.Success(fmt.Sprintf("%d", summary.TotalMessages)))

	fmt.Fprintf(os.Stdout, "Model Usage:\n")
	fmt.Fprintf(os.Stdout, "------------\n")
	for _, model := range summary.ModelBreakdown {
		fmt.Fprintf(os.Stdout, "%s:\n", model.Model)
		fmt.Fprintf(os.Stdout, "  Input:  %d tokens\n", model.InputTokens)
		fmt.Fprintf(os.Stdout, "  Output: %d tokens\n", model.OutputTokens)
		fmt.Println()
	}

	fmt.Fprintf(os.Stdout, "Recent Activity (last 7 days):\n")
	fmt.Fprintf(os.Stdout, "------------------------------\n")
	for _, day := range summary.RecentActivity {
		fmt.Fprintf(os.Stdout, "%s: %d messages, %d tool calls\n",
			day.Date, day.MessageCount, day.ToolCallCount)
	}

	return nil
}

func runClaudeStatsTokens(cmd *cobra.Command, args []string) error {
	helper := claude.NewHelper(claudeVerbose, claudeDryRun)
	usage, err := helper.GetTokenUsage()
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(claudeOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(usage)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n\n", output.Info("Token Usage by Model"))

	for _, model := range usage.Models {
		fmt.Fprintf(os.Stdout, "[%s]\n", output.Success(model.Model))
		fmt.Fprintf(os.Stdout, "  Input:         %d\n", model.InputTokens)
		fmt.Fprintf(os.Stdout, "  Output:        %d\n", model.OutputTokens)
		fmt.Fprintf(os.Stdout, "  Cache Read:    %d\n", model.CacheReadInputTokens)
		fmt.Fprintf(os.Stdout, "  Cache Create:  %d\n", model.CacheCreationInputTokens)
		fmt.Println()
	}

	fmt.Fprintf(os.Stdout, "Totals:\n")
	fmt.Fprintf(os.Stdout, "  Input:         %d\n", usage.Total.Input)
	fmt.Fprintf(os.Stdout, "  Output:        %d\n", usage.Total.Output)
	fmt.Fprintf(os.Stdout, "  Cache Read:    %d\n", usage.Total.CacheRead)
	fmt.Fprintf(os.Stdout, "  Cache Create:  %d\n", usage.Total.CacheCreation)

	return nil
}

func runClaudeStatsDaily(cmd *cobra.Command, args []string) error {
	days := 7
	if len(args) > 0 {
		var err error
		_, err = fmt.Sscanf(args[0], "%d", &days)
		if err != nil {
			return fmt.Errorf("invalid days argument: %s", args[0])
		}
	}

	helper := claude.NewHelper(claudeVerbose, claudeDryRun)
	usage, err := helper.GetDailyUsage(days)
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(claudeOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(usage)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n\n", output.Info(fmt.Sprintf("Daily Token Usage (last %d days)", days)))

	for _, day := range usage.Days {
		fmt.Fprintf(os.Stdout, "%s:\n", output.Success(day.Date))
		for _, model := range day.Models {
			fmt.Fprintf(os.Stdout, "  %s: %d tokens\n", model.Model, model.Tokens)
		}
		fmt.Println()
	}

	return nil
}

func runClaudePermissions(cmd *cobra.Command, args []string) error {
	helper := claude.NewHelper(claudeVerbose, claudeDryRun)
	perms, err := helper.GetPermissions()
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(claudeOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(perms)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n\n", output.Info("Claude Code Permissions"))

	fmt.Fprintf(os.Stdout, "Allowed:\n")
	if len(perms.Allow) == 0 {
		fmt.Fprintf(os.Stdout, "  (none)\n")
	} else {
		for _, rule := range perms.Allow {
			fmt.Fprintf(os.Stdout, "  + %s\n", output.Success(rule))
		}
	}

	fmt.Println()
	fmt.Fprintf(os.Stdout, "Denied:\n")
	if len(perms.Deny) == 0 {
		fmt.Fprintf(os.Stdout, "  (none)\n")
	} else {
		for _, rule := range perms.Deny {
			fmt.Fprintf(os.Stdout, "  - %s\n", output.Error(rule))
		}
	}

	return nil
}

func runClaudePermissionsAdd(cmd *cobra.Command, args []string) error {
	rule := args[0]
	permType := "allow"
	if len(args) > 1 {
		permType = args[1]
		if permType != "allow" && permType != "deny" {
			return fmt.Errorf("invalid permission type: %s (use allow or deny)", permType)
		}
	}

	helper := claude.NewHelper(claudeVerbose, claudeDryRun)
	if err := helper.AddPermission(rule, permType); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Added %s rule: %s\n", permType, output.Success(rule))
	return nil
}

func runClaudePermissionsRemove(cmd *cobra.Command, args []string) error {
	rule := args[0]
	permType := "allow"
	if len(args) > 1 {
		permType = args[1]
		if permType != "allow" && permType != "deny" {
			return fmt.Errorf("invalid permission type: %s (use allow or deny)", permType)
		}
	}

	helper := claude.NewHelper(claudeVerbose, claudeDryRun)
	if err := helper.RemovePermission(rule, permType); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Removed %s rule: %s\n", permType, rule)
	return nil
}

func runClaudeSettings(cmd *cobra.Command, args []string) error {
	typeArg := ""
	if len(args) > 0 {
		typeArg = args[0]
	}

	st, err := claude.ParseSettingsType(typeArg)
	if err != nil {
		return err
	}

	helper := claude.NewHelper(claudeVerbose, claudeDryRun)
	path := helper.GetSettingsPath(st)

	settings, err := helper.GetSettingsRaw(st)
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(claudeOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(settings)
	}

	// Table format - pretty print JSON
	fmt.Fprintf(os.Stdout, "%s Settings (%s):\n\n", output.Info(string(st)), path)
	jsonData, _ := json.MarshalIndent(settings, "", "  ")
	fmt.Println(string(jsonData))

	return nil
}

func runClaudeSettingsEdit(cmd *cobra.Command, args []string) error {
	typeArg := ""
	if len(args) > 0 {
		typeArg = args[0]
	}

	st, err := claude.ParseSettingsType(typeArg)
	if err != nil {
		return err
	}

	helper := claude.NewHelper(claudeVerbose, claudeDryRun)
	return helper.EditSettings(st)
}

func runClaudeProjects(cmd *cobra.Command, args []string) error {
	helper := claude.NewHelper(claudeVerbose, claudeDryRun)
	projects, err := helper.GetProjects()
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(claudeOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(projects)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n\n", output.Info("Claude Code Projects"))

	if len(projects.Projects) == 0 {
		fmt.Println("No trusted projects found.")
		return nil
	}

	for _, p := range projects.Projects {
		fmt.Fprintf(os.Stdout, "%s\n", output.Success(p.Path))
		fmt.Fprintf(os.Stdout, "  Cost: $%.2f\n\n", p.Cost)
	}

	return nil
}

func runClaudeMcp(cmd *cobra.Command, args []string) error {
	helper := claude.NewHelper(claudeVerbose, claudeDryRun)
	mcpView, err := helper.GetMCPServers()
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(claudeOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(mcpView)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n\n", output.Info("MCP Servers"))

	if len(mcpView.Servers) == 0 {
		fmt.Println("No MCP servers configured.")
	} else {
		for _, server := range mcpView.Servers {
			fmt.Fprintf(os.Stdout, "[%s] %s\n", output.Success(server.Name), server.Project)
			fmt.Fprintf(os.Stdout, "  Type: %s\n", server.Type)
			fmt.Fprintf(os.Stdout, "  URL:  %s\n\n", server.URL)
		}
	}

	if mcpView.HasLocal {
		fmt.Fprintf(os.Stdout, "Local MCP Config: %s\n", output.Success(mcpView.LocalFile))
	}

	return nil
}

func runClaudeMcpAdd(cmd *cobra.Command, args []string) error {
	name := args[0]
	url := args[1]
	serverType := "http"
	if len(args) > 2 {
		serverType = args[2]
	}

	helper := claude.NewHelper(claudeVerbose, claudeDryRun)
	if err := helper.AddMCPServer(name, url, serverType); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Added MCP server '%s' to .mcp.json\n", output.Success(name))
	return nil
}

func runClaudeCommands(cmd *cobra.Command, args []string) error {
	helper := claude.NewHelper(claudeVerbose, claudeDryRun)
	commands, err := helper.GetCommands()
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(claudeOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(commands)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n\n", output.Info("Custom Commands"))

	if len(commands.UserCommands) > 0 {
		fmt.Fprintf(os.Stdout, "User Commands (~/.claude/commands/):\n")
		for _, cmd := range commands.UserCommands {
			fmt.Fprintf(os.Stdout, "  %s\n", output.Success(cmd.Name))
		}
		fmt.Println()
	}

	if len(commands.ProjectCommands) > 0 {
		fmt.Fprintf(os.Stdout, "Project Commands (.claude/commands/):\n")
		for _, cmd := range commands.ProjectCommands {
			fmt.Fprintf(os.Stdout, "  %s\n", output.Success(cmd.Name))
		}
	}

	if len(commands.UserCommands) == 0 && len(commands.ProjectCommands) == 0 {
		fmt.Println("No custom commands found.")
	}

	return nil
}

func runClaudeAggregate(cmd *cobra.Command, args []string) error {
	searchDir := os.Getenv("HOME") + "/Repos"
	if len(args) > 0 {
		searchDir = args[0]
	}

	helper := claude.NewHelper(claudeVerbose, claudeDryRun)
	result, err := helper.Aggregate(searchDir)
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(claudeOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(result)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "Scanning %s for Claude Code configurations...\n\n", searchDir)

	for _, item := range result.Items {
		switch item.Action {
		case "added":
			fmt.Fprintf(os.Stdout, "  Added: %s/%s (from %s)\n",
				item.Type+"s", item.FileName, item.SourceRepo)
		case "renamed":
			fmt.Fprintf(os.Stdout, "  Added (renamed): %s/%s (from %s)\n",
				item.Type+"s", item.FileName, item.SourceRepo)
		case "skipped":
			if claudeVerbose {
				fmt.Fprintf(os.Stdout, "  Skipped (duplicate): %s/%s\n",
					item.Type+"s", item.FileName)
			}
		}
	}

	fmt.Println()
	fmt.Fprintf(os.Stdout, "Summary:\n")
	fmt.Fprintf(os.Stdout, "  Repos with configs: %d\n", result.ReposScanned)
	fmt.Fprintf(os.Stdout, "  Agents added:       %d\n", result.AgentsAdded)
	fmt.Fprintf(os.Stdout, "  Commands added:     %d\n", result.CommandsAdded)
	fmt.Fprintf(os.Stdout, "  Subagents added:    %d\n", result.SubagentsAdded)
	fmt.Fprintf(os.Stdout, "  Skipped (dups):     %d\n", result.Skipped)
	fmt.Fprintf(os.Stdout, "  Renamed (conflicts):%d\n", result.Renamed)

	return nil
}

func runClaudeAggregateList(cmd *cobra.Command, args []string) error {
	helper := claude.NewHelper(claudeVerbose, claudeDryRun)
	result, err := helper.List()
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(claudeOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(result)
	}

	// Table format
	fmt.Fprintf(os.Stdout, "%s\n\n", output.Info("Claude Code Agents & Commands"))

	if len(result.Agents) > 0 {
		fmt.Fprintf(os.Stdout, "Agents (%d):\n", len(result.Agents))
		for _, agent := range result.Agents {
			fmt.Fprintf(os.Stdout, "  %s\n", agent)
		}
		fmt.Println()
	}

	if len(result.Commands) > 0 {
		fmt.Fprintf(os.Stdout, "Commands (%d):\n", len(result.Commands))
		for _, cmd := range result.Commands {
			fmt.Fprintf(os.Stdout, "  %s\n", cmd)
		}
		fmt.Println()
	}

	if len(result.Subagents) > 0 {
		fmt.Fprintf(os.Stdout, "Subagents (%d):\n", len(result.Subagents))
		for _, sub := range result.Subagents {
			fmt.Fprintf(os.Stdout, "  %s\n", sub)
		}
	}

	return nil
}

func runClaudeClear(cmd *cobra.Command, args []string) error {
	what := "cache"
	if len(args) > 0 {
		what = args[0]
	}

	helper := claude.NewHelper(claudeVerbose, claudeDryRun)
	result, err := helper.Clear(what)
	if err != nil {
		return err
	}

	format, err := output.ParseFormat(claudeOutputFormat)
	if err != nil {
		return err
	}

	if format != output.FormatTable {
		printer := output.NewPrinter(os.Stdout, format)
		return printer.Print(result)
	}

	// Table format
	if len(result.Cleared) == 0 {
		fmt.Printf("Nothing to clear for %s\n", what)
	} else {
		fmt.Printf("Cleared %s:\n", what)
		for _, path := range result.Cleared {
			fmt.Printf("  %s\n", output.Success(path))
		}
	}

	return nil
}

func runClaudeHelp(cmd *cobra.Command, args []string) error {
	helpText := `Claude Code Shell Functions
===========================

Aliases:
  cc          - claude
  ccc         - claude --continue
  ccr         - claude --resume
  ccp         - claude -p (piped mode)

Stats & Usage:
  acorn claude stats         - View usage statistics
  acorn claude stats tokens  - View token usage by model
  acorn claude stats daily   - View daily token usage

Permissions:
  acorn claude permissions        - View all permissions
  acorn claude permissions add    - Add a permission rule
  acorn claude permissions remove - Remove a permission rule

Settings:
  acorn claude settings [type]      - View settings [global|local|config]
  acorn claude settings edit [type] - Edit settings file

Projects & MCP:
  acorn claude projects      - List all projects
  acorn claude mcp           - List MCP servers
  acorn claude mcp add       - Add MCP server
  acorn claude commands      - List custom commands

Utilities:
  acorn claude info          - Show Claude Code info
  acorn claude clear         - Clear cache/stats
  acorn claude help          - Show this help

Aggregation:
  acorn claude aggregate [dir] - Collect agents/commands from repos
  acorn claude aggregate list  - List all agents and commands
`
	fmt.Print(helpText)
	return nil
}

func runClaudeInstall(cmd *cobra.Command, args []string) error {
	inst := installer.NewInstaller(
		installer.WithDryRun(claudeDryRun),
		installer.WithVerbose(claudeVerbose),
	)

	// Show platform info
	platform := inst.GetPlatform()
	if claudeVerbose {
		fmt.Fprintf(os.Stdout, "Platform: %s (%s)\n\n", platform.OS, platform.PackageManager)
	}

	// Get the plan first
	plan, err := inst.Plan(cmd.Context(), "claude")
	if err != nil {
		return err
	}

	// Show what will be installed
	if claudeDryRun {
		fmt.Fprintf(os.Stdout, "%s\n", output.Info("Claude Installation Plan"))
		fmt.Fprintln(os.Stdout, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Fprintf(os.Stdout, "Platform: %s (%s)\n\n", platform.OS, platform.PackageManager)
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
		}
		fmt.Fprintf(os.Stdout, "  %s %s - %s%s\n", status, t.Name, t.Description, suffix)
	}

	if claudeDryRun {
		fmt.Fprintln(os.Stdout, "\nRun without --dry-run to install.")
		return nil
	}

	// Execute installation
	fmt.Fprintln(os.Stdout, "\nInstalling...")
	result, err := inst.Install(cmd.Context(), "claude")
	if err != nil {
		return err
	}

	installed, skipped, failed := result.Summary()
	if failed == 0 {
		fmt.Fprintf(os.Stdout, "%s Installation complete (%d installed, %d skipped)\n",
			output.Success("✓"), installed, skipped)
	} else {
		fmt.Fprintf(os.Stdout, "%s Installation failed (%d installed, %d skipped, %d failed)\n",
			output.Error("✗"), installed, skipped, failed)
	}

	return nil
}
