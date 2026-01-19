package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	toolDocs     []string
	toolCategory string
	toolModel    string
)

// toolCmd represents the tool command group
var toolCmd = &cobra.Command{
	Use:   "tool",
	Short: "Manage tool agents and commands",
	Long: `Manage tool-specific agents and slash commands.

This command helps you create Claude Code agents and slash commands
for various tools based on their documentation.

Examples:
  acorn tool generate docker --docs https://docs.docker.com/reference/cli/docker/
  acorn tool generate kubectl --category devops
  acorn tool list`,
}

// toolGenerateCmd generates agent and commands for a tool
var toolGenerateCmd = &cobra.Command{
	Use:   "generate <tool-name>",
	Short: "Generate agent and commands for a tool",
	Long: `Generate a Claude Code agent and slash commands for a tool.

Provide documentation URLs to help generate comprehensive commands.
The command creates:
  - An expert agent in .sapling/config/claude/agents/<tool>-expert.md
  - Slash commands in .sapling/config/claude/commands/<tool>/

Examples:
  acorn tool generate docker --docs https://docs.docker.com/reference/cli/docker/
  acorn tool generate helm --docs https://helm.sh/docs/ --category devops`,
	Args: cobra.ExactArgs(1),
	RunE: runToolGenerate,
}

// toolListCmd lists existing tool agents
var toolListCmd = &cobra.Command{
	Use:   "list",
	Short: "List existing tool agents and commands",
	RunE:  runToolList,
}

func init() {
	rootCmd.AddCommand(toolCmd)

	toolCmd.AddCommand(toolGenerateCmd)
	toolCmd.AddCommand(toolListCmd)

	toolGenerateCmd.Flags().StringSliceVar(&toolDocs, "docs", []string{},
		"Documentation URLs to reference (can specify multiple)")
	toolGenerateCmd.Flags().StringVarP(&toolCategory, "category", "c", "",
		"Tool category (devops, cloud, programming, terminal, etc.)")
	toolGenerateCmd.Flags().StringVarP(&toolModel, "model", "m", "sonnet",
		"Claude model for the agent (sonnet, opus, haiku)")
}

func runToolGenerate(cmd *cobra.Command, args []string) error {
	toolName := args[0]

	dotfilesRoot, err := getDotfilesRoot()
	if err != nil {
		return err
	}

	// Create directories
	agentsDir := filepath.Join(dotfilesRoot, ".sapling", "config", "claude", "agents")
	commandsDir := filepath.Join(dotfilesRoot, ".sapling", "config", "claude", "commands", toolName)

	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		return fmt.Errorf("failed to create agents directory: %w", err)
	}
	if err := os.MkdirAll(commandsDir, 0755); err != nil {
		return fmt.Errorf("failed to create commands directory: %w", err)
	}

	// Generate agent
	agentPath := filepath.Join(agentsDir, toolName+"-expert.md")
	if err := generateAgent(agentPath, toolName, toolDocs, toolCategory, toolModel); err != nil {
		return fmt.Errorf("failed to generate agent: %w", err)
	}
	fmt.Printf("Created agent: %s\n", agentPath)

	// Generate standard commands
	commands := getStandardCommands(toolName)
	for _, c := range commands {
		cmdPath := filepath.Join(commandsDir, c.Filename)
		if err := generateCommand(cmdPath, c); err != nil {
			return fmt.Errorf("failed to generate command %s: %w", c.Filename, err)
		}
		fmt.Printf("Created command: %s\n", cmdPath)
	}

	fmt.Printf("\nGenerated %d files for %s\n", len(commands)+1, toolName)
	fmt.Println("\nNext steps:")
	fmt.Println("1. Review and customize the generated files")
	fmt.Println("2. Add tool-specific commands based on documentation")
	fmt.Println("3. Run 'acorn setup' to sync to ~/.claude/")

	if len(toolDocs) > 0 {
		fmt.Println("\nDocumentation references:")
		for _, doc := range toolDocs {
			fmt.Printf("  - %s\n", doc)
		}
	}

	return nil
}

func runToolList(cmd *cobra.Command, args []string) error {
	dotfilesRoot, err := getDotfilesRoot()
	if err != nil {
		return err
	}

	agentsDir := filepath.Join(dotfilesRoot, ".sapling", "config", "claude", "agents")
	commandsDir := filepath.Join(dotfilesRoot, ".sapling", "config", "claude", "commands")

	// List agents
	fmt.Println("Agents:")
	agents, err := os.ReadDir(agentsDir)
	if err == nil {
		for _, a := range agents {
			if name, found := strings.CutSuffix(a.Name(), "-expert.md"); found {
				fmt.Printf("  - %s\n", name)
			}
		}
	}

	// List command categories
	fmt.Println("\nCommand Categories:")
	categories, err := os.ReadDir(commandsDir)
	if err == nil {
		for _, c := range categories {
			if c.IsDir() {
				// Count commands in category
				cmdDir := filepath.Join(commandsDir, c.Name())
				cmds, _ := os.ReadDir(cmdDir)
				count := 0
				for _, cmd := range cmds {
					if strings.HasSuffix(cmd.Name(), ".md") {
						count++
					}
				}
				fmt.Printf("  - %s (%d commands)\n", c.Name(), count)
			}
		}
	}

	return nil
}

// CommandTemplate represents a command to generate
type CommandTemplate struct {
	Filename    string
	Name        string
	Description string
	ArgHint     string
	Tools       string
	Content     string
}

func getStandardCommands(toolName string) []CommandTemplate {
	title := cases.Title(language.English).String(toolName)

	return []CommandTemplate{
		{
			Filename:    toolName + "-explain.md",
			Name:        toolName + "-explain",
			Description: fmt.Sprintf("Explain %s concepts, commands, and workflows", title),
			ArgHint:     "[topic]",
			Tools:       "Read, Glob, Grep, WebFetch, WebSearch",
			Content: fmt.Sprintf(`## Task

Explain %s concepts, commands, or workflows to the user.

## Approach

1. **Understand the question** - What aspect of %s are they asking about?
2. **Search for context** - Check local configs, documentation
3. **Provide clear explanation** - Use examples and practical advice
4. **Reference documentation** - Link to official docs when helpful

## Topics to Cover

- Basic concepts and terminology
- Common commands and usage
- Configuration options
- Best practices
- Troubleshooting tips

## Context

Check the user's %s configuration:
@.sapling/config/%s/config.yaml

## User's Question

$ARGUMENTS
`, title, title, toolName, toolName),
		},
		{
			Filename:    toolName + "-coach.md",
			Name:        toolName + "-coach",
			Description: fmt.Sprintf("Interactive coaching session to learn %s step by step", title),
			ArgHint:     "[skill-level: beginner|intermediate|advanced]",
			Tools:       "Read, Glob, Grep, Bash",
			Content: fmt.Sprintf(`## Task

Guide the user through learning %s interactively.

## Approach

1. **Assess level** - Ask about their %s experience
2. **Set goals** - Identify what they want to accomplish
3. **Progressive exercises** - Start simple, build complexity
4. **Real-time practice** - Have them run commands
5. **Reinforce** - Summarize and suggest next steps

## Skill Levels

### Beginner
- Installation and basic setup
- Core concepts and terminology
- Basic commands and workflows
- Reading help and documentation

### Intermediate
- Advanced configuration
- Common patterns and best practices
- Troubleshooting
- Integration with other tools

### Advanced
- Performance optimization
- Custom workflows
- Scripting and automation
- Advanced features

## Context

Check available aliases and functions:
@.sapling/config/%s/config.yaml

## Coaching Style

- Start with fundamentals
- Use practical examples
- Reference available aliases/functions
- Build toward real-world workflows

## Skill Level

$ARGUMENTS
`, title, title, toolName),
		},
	}
}

const agentTemplate = `---
name: {{.Name}}-expert
description: Expert in {{.Title}} {{.CategoryDesc}}
tools: Read, Write, Edit, Glob, Grep, Bash
model: {{.Model}}
---

You are a **{{.Title}} Expert** specializing in {{.Title}} {{.CategoryDesc}}.

## Your Core Competencies

- {{.Title}} installation and configuration
- Common workflows and best practices
- Troubleshooting and debugging
- Integration with other tools
- Performance optimization

## Documentation References
{{range .Docs}}
- {{.}}
{{end}}
{{if not .Docs}}
- Official {{.Title}} documentation
{{end}}

## Available Configuration

Check the component configuration:
@.sapling/config/{{.Name}}/config.yaml

## Key Aliases

| Alias | Command |
|-------|---------|
| (see config.yaml for available aliases) |

## Shell Functions

| Function | Description |
|----------|-------------|
| (see config.yaml for available functions) |

## Best Practices

1. **Configuration** - Use XDG-compliant paths
2. **Security** - Follow principle of least privilege
3. **Automation** - Use shell functions for common workflows
4. **Documentation** - Keep configs well-commented

## Your Approach

When providing {{.Title}} guidance:
1. **Check** existing configuration and setup
2. **Understand** the user's goal
3. **Suggest** best practices
4. **Provide** practical examples
5. **Verify** the solution works
`

const commandTemplate = `---
description: {{.Description}}
argument-hint: {{.ArgHint}}
allowed-tools: {{.Tools}}
---

{{.Content}}
`

type AgentData struct {
	Name         string
	Title        string
	CategoryDesc string
	Model        string
	Docs         []string
}

func generateAgent(path, toolName string, docs []string, category, model string) error {
	categoryDesc := ""
	if category != "" {
		categoryDesc = fmt.Sprintf("for %s workflows", category)
	}

	data := AgentData{
		Name:         toolName,
		Title:        cases.Title(language.English).String(toolName),
		CategoryDesc: categoryDesc,
		Model:        model,
		Docs:         docs,
	}

	tmpl, err := template.New("agent").Parse(agentTemplate)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}

func generateCommand(path string, cmd CommandTemplate) error {
	tmpl, err := template.New("command").Parse(commandTemplate)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, cmd)
}
