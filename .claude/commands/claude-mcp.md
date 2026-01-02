---
description: Configure MCP (Model Context Protocol) servers for Claude Code
argument-hint: [server-name] [description]
allowed-tools: Read, Write, Edit, Glob, Bash
---

# MCP Server Configuration

Help the user configure MCP servers for Claude Code.

## What is MCP?

Model Context Protocol connects Claude Code with external tools and services:
- Databases
- APIs
- File systems
- Custom tools

## Configuration Locations

| Scope | File |
|-------|------|
| User | `~/.claude.json` |
| Project | `.mcp.json` |
| Plugin | `plugin.json` or `.mcp.json` in plugin |

## Your Task

Based on the user's request: $ARGUMENTS

1. Determine the appropriate MCP server configuration
2. Create or update the configuration file
3. Provide any required setup steps

## Configuration Format

```json
{
  "mcpServers": {
    "server-name": {
      "command": "/path/to/server",
      "args": ["--config", "config.json"],
      "env": {
        "API_KEY": "your-key"
      }
    }
  }
}
```

## Common MCP Servers

### Filesystem
```json
{
  "mcpServers": {
    "filesystem": {
      "command": "npx",
      "args": ["-y", "@anthropic/mcp-server-filesystem", "/path/to/directory"]
    }
  }
}
```

### SQLite Database
```json
{
  "mcpServers": {
    "sqlite": {
      "command": "npx",
      "args": ["-y", "@anthropic/mcp-server-sqlite", "--db", "database.db"]
    }
  }
}
```

### GitHub
```json
{
  "mcpServers": {
    "github": {
      "command": "npx",
      "args": ["-y", "@anthropic/mcp-server-github"],
      "env": {
        "GITHUB_TOKEN": "your-token"
      }
    }
  }
}
```

### PostgreSQL
```json
{
  "mcpServers": {
    "postgres": {
      "command": "npx",
      "args": ["-y", "@anthropic/mcp-server-postgres"],
      "env": {
        "DATABASE_URL": "postgresql://user:pass@host:5432/db"
      }
    }
  }
}
```

### Custom Server
```json
{
  "mcpServers": {
    "my-server": {
      "command": "${CLAUDE_PROJECT_DIR}/mcp/server.js",
      "args": ["--mode", "production"],
      "env": {
        "API_KEY": "${API_KEY}"
      }
    }
  }
}
```

## Plugin MCP Servers

In plugin's `.mcp.json`:
```json
{
  "mcpServers": {
    "plugin-db": {
      "command": "${CLAUDE_PLUGIN_ROOT}/servers/db-server",
      "args": ["--config", "${CLAUDE_PLUGIN_ROOT}/config.json"]
    }
  }
}
```

## MCP Slash Commands

MCP servers can expose prompts as slash commands:
```
/mcp__github__list_prs
/mcp__github__pr_review 456
```

## Managing MCP

- `/mcp` - View servers, check status, authenticate OAuth
- Servers start automatically when Claude Code loads

## Environment Variables

- `${CLAUDE_PROJECT_DIR}` - Project root
- `${CLAUDE_PLUGIN_ROOT}` - Plugin directory (in plugins)
- `${VAR_NAME}` - Reference environment variables

## Troubleshooting

1. **Server not starting**: Check command path and permissions
2. **Authentication issues**: Use `/mcp` to authenticate OAuth servers
3. **Connection errors**: Verify server is running with `claude --debug`

## Security

- Store secrets in environment variables, not config files
- Use `settings.local.json` for sensitive server configs
- Review MCP server permissions carefully
