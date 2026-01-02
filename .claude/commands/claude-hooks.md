---
description: Create or configure Claude Code hooks for tool events
argument-hint: [hook-type] [description]
allowed-tools: Read, Write, Edit, Glob, Bash
---

# Claude Code Hooks Configuration

Help the user create or configure hooks for Claude Code.

## Hook Events Reference

| Event | Purpose | Matcher |
|-------|---------|---------|
| **PreToolUse** | Before tool execution | Tool names |
| **PostToolUse** | After tool completes | Tool names |
| **PermissionRequest** | Permission dialog shown | Tool names |
| **UserPromptSubmit** | Before processing prompt | N/A |
| **Stop** | Main agent finishes | N/A |
| **SubagentStop** | Subagent finishes | N/A |
| **SessionStart** | Session begins | startup, resume, clear, compact |
| **SessionEnd** | Session ends | N/A |
| **Notification** | During notifications | permission_prompt, idle_prompt |
| **PreCompact** | Before compacting | manual, auto |

## Hook Types

1. **Command hooks** (`type: "command"`) - Execute bash scripts
2. **Prompt hooks** (`type: "prompt"`) - Use LLM for decisions

## Configuration Location

- User: `~/.claude/settings.json`
- Project: `.claude/settings.json`
- Local: `.claude/settings.local.json`

## Your Task

Based on the user's request: $ARGUMENTS

1. Determine which hook event is needed
2. Create the hook configuration JSON
3. Create any required scripts in `.claude/hooks/`
4. Make scripts executable
5. Show how to add to settings.json

## Template

```json
{
  "hooks": {
    "EventName": [
      {
        "matcher": "ToolPattern",
        "hooks": [
          {
            "type": "command",
            "command": "${CLAUDE_PROJECT_DIR}/.claude/hooks/script.sh",
            "timeout": 60
          }
        ]
      }
    ]
  }
}
```

## Script Template

```bash
#!/usr/bin/env bash
# Hook script - receives JSON via stdin

input=$(cat)
tool_name=$(echo "$input" | jq -r '.tool_name // ""')

# Your logic here

# Exit codes:
# 0 = success (stdout added to context)
# 2 = block action (stderr shown as error)
# other = warning (stderr shown in verbose mode)

exit 0
```

## Security Reminders

- Always validate/sanitize inputs
- Quote all variables: `"$VAR"` not `$VAR`
- Use absolute paths for scripts
- Make scripts executable: `chmod +x`
