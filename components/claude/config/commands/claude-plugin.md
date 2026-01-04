---
description: Create or configure Claude Code plugins
argument-hint: [plugin-name] [description]
allowed-tools: Read, Write, Edit, Glob, Bash
---

# Claude Code Plugin Development

Help the user create or configure Claude Code plugins.

## Plugin Components

1. **Commands** - Custom slash commands (`commands/`)
2. **Agents** - Specialized subagents (`agents/`)
3. **Skills** - Model-invoked capabilities (`skills/`)
4. **Hooks** - Event handlers (`hooks/hooks.json`)
5. **MCP Servers** - External tool connections (`.mcp.json`)
6. **LSP Servers** - Language intelligence (`.lsp.json`)

## Directory Structure

```
my-plugin/
├── .claude-plugin/
│   └── plugin.json          # Required manifest
├── commands/                 # Slash commands
│   └── my-command.md
├── agents/                   # Subagents
│   └── my-agent.md
├── skills/                   # Skills
│   └── my-skill/
│       └── SKILL.md
├── hooks/
│   └── hooks.json
├── .mcp.json                # MCP servers
├── .lsp.json                # LSP servers
└── scripts/                 # Utility scripts
```

## Your Task

Based on the user's request: $ARGUMENTS

1. Create the plugin directory structure
2. Generate `plugin.json` manifest
3. Create requested components
4. Provide installation instructions

## Plugin Manifest Template

```json
{
  "name": "my-plugin",
  "version": "1.0.0",
  "description": "Plugin description",
  "author": {
    "name": "Author Name"
  },
  "commands": "./commands/",
  "agents": "./agents/",
  "skills": "./skills/",
  "hooks": "./hooks/hooks.json",
  "mcpServers": "./.mcp.json"
}
```

## Installation

```bash
# Install to user scope
claude plugin install ./my-plugin --scope user

# Install to project scope
claude plugin install ./my-plugin --scope project

# Enable/disable
claude plugin enable my-plugin
claude plugin disable my-plugin
```

## Environment Variables

- `${CLAUDE_PLUGIN_ROOT}` - Absolute path to plugin directory (use in all paths)

## Best Practices

- Use `${CLAUDE_PLUGIN_ROOT}` for all file references
- Keep plugins self-contained (no external references)
- Use semantic versioning
- Include LICENSE and CHANGELOG.md
