---
description: Configure Claude Code settings (permissions, env, model, sandbox)
argument-hint: [setting-area]
allowed-tools: Read, Write, Edit, Glob
---

# Claude Code Settings Configuration

Help the user configure Claude Code settings.

## Settings Hierarchy (Highest to Lowest)

1. **Managed** (Enterprise) - Cannot be overridden
2. **CLI arguments** - Temporary overrides
3. **Local project** - `.claude/settings.local.json`
4. **Shared project** - `.claude/settings.json`
5. **User** - `~/.claude/settings.json`

## File Locations

| Scope | Location |
|-------|----------|
| User | `~/.claude/settings.json` |
| Project | `.claude/settings.json` |
| Local | `.claude/settings.local.json` |
| Enterprise (macOS) | `/Library/Application Support/ClaudeCode/managed-settings.json` |
| Enterprise (Linux) | `/etc/claude-code/managed-settings.json` |

## Your Task

Based on the user's request: $ARGUMENTS

1. Identify the setting area
2. Determine appropriate scope
3. Create or update settings.json

## Common Settings

### Permissions
```json
{
  "permissions": {
    "allow": [
      "Bash(npm:*)",
      "Bash(git:*)",
      "Read(~/.zshrc)"
    ],
    "ask": [
      "Bash(git push:*)"
    ],
    "deny": [
      "Read(./.env)",
      "Read(./secrets/**)",
      "WebFetch"
    ],
    "defaultMode": "acceptEdits"
  }
}
```

### Environment Variables
```json
{
  "env": {
    "NODE_ENV": "development",
    "ANTHROPIC_MODEL": "claude-sonnet-4-5-20250929"
  }
}
```

### Model Configuration
```json
{
  "model": "claude-sonnet-4-5-20250929"
}
```

### Sandbox
```json
{
  "sandbox": {
    "enabled": true,
    "excludedCommands": ["docker", "git"],
    "network": {
      "allowLocalBinding": true
    }
  }
}
```

### Status Line
```json
{
  "statusLine": {
    "type": "command",
    "command": "~/.claude/statusline.sh",
    "padding": 0
  }
}
```

### Attribution
```json
{
  "attribution": {
    "commit": "Generated with Claude Code",
    "pr": ""
  }
}
```

## Permission Modes

- `acceptEdits` - Auto-approve file edits
- `askOnEdit` - Ask for each edit
- `readOnly` - No file modifications
- `bypassPermissions` - Skip all permission prompts

## Key Environment Variables

| Variable | Purpose |
|----------|---------|
| `ANTHROPIC_API_KEY` | API key |
| `ANTHROPIC_MODEL` | Override model |
| `DISABLE_TELEMETRY` | Opt out of telemetry |
| `BASH_DEFAULT_TIMEOUT_MS` | Command timeout |
| `MAX_THINKING_TOKENS` | Extended thinking |

## Commands

- `/config` - Open settings interface
- `/permissions` - View/update permissions
