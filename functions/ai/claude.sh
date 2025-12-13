#!/bin/sh
# Claude Code management functions
# Requires: jq

# =============================================================================
# Configuration Paths
# =============================================================================

CLAUDE_DIR="${HOME}/.claude"
CLAUDE_CONFIG="${HOME}/.claude.json"
CLAUDE_SETTINGS="${CLAUDE_DIR}/settings.json"
CLAUDE_LOCAL="${CLAUDE_DIR}/settings.local.json"
CLAUDE_STATS="${CLAUDE_DIR}/stats-cache.json"
CLAUDE_PROJECTS="${CLAUDE_DIR}/projects"

# =============================================================================
# Quick Aliases
# =============================================================================

alias cc='claude'
alias ccc='claude --continue'
alias ccr='claude --resume'
alias ccp='claude -p'

# =============================================================================
# Stats and Usage
# =============================================================================

# View Claude Code usage statistics
claude_stats() {
    if [ ! -f "$CLAUDE_STATS" ]; then
        echo "No stats file found at $CLAUDE_STATS"
        return 1
    fi

    echo "Claude Code Usage Statistics"
    echo "============================"
    echo ""

    # Basic stats
    local total_sessions total_messages
    total_sessions=$(jq -r '.totalSessions // 0' "$CLAUDE_STATS")
    total_messages=$(jq -r '.totalMessages // 0' "$CLAUDE_STATS")

    echo "Total Sessions: $total_sessions"
    echo "Total Messages: $total_messages"
    echo ""

    # Model usage
    echo "Model Usage:"
    echo "------------"
    jq -r '.modelUsage | to_entries[] | "\(.key):\n  Input: \(.value.inputTokens // 0) tokens\n  Output: \(.value.outputTokens // 0) tokens\n  Cache Read: \(.value.cacheReadInputTokens // 0) tokens\n  Cache Create: \(.value.cacheCreationInputTokens // 0) tokens"' "$CLAUDE_STATS" 2>/dev/null

    echo ""

    # Longest session
    echo "Longest Session:"
    echo "----------------"
    jq -r '.longestSession | "Messages: \(.messageCount // 0)\nDuration: \((.duration // 0) / 1000 / 60 | floor) minutes"' "$CLAUDE_STATS" 2>/dev/null

    echo ""

    # Recent daily activity
    echo "Recent Activity (last 7 days):"
    echo "------------------------------"
    jq -r '.dailyActivity | .[-7:][] | "\(.date): \(.messageCount) messages, \(.toolCallCount) tool calls"' "$CLAUDE_STATS" 2>/dev/null
}

# View token usage summary
claude_tokens() {
    if [ ! -f "$CLAUDE_STATS" ]; then
        echo "No stats file found"
        return 1
    fi

    echo "Token Usage by Model"
    echo "===================="
    echo ""

    jq -r '
        .modelUsage | to_entries[] |
        "[\(.key)]",
        "  Input:        \(.value.inputTokens // 0 | tostring | . as $n | if (. | length) > 3 then ($n[:-3] + "," + $n[-3:]) else . end)",
        "  Output:       \(.value.outputTokens // 0 | tostring | . as $n | if (. | length) > 3 then ($n[:-3] + "," + $n[-3:]) else . end)",
        "  Cache Read:   \(.value.cacheReadInputTokens // 0 | tostring | . as $n | if (. | length) > 6 then ($n[:-6] + "," + $n[-6:-3] + "," + $n[-3:]) elif (. | length) > 3 then ($n[:-3] + "," + $n[-3:]) else . end)",
        "  Cache Create: \(.value.cacheCreationInputTokens // 0 | tostring | . as $n | if (. | length) > 6 then ($n[:-6] + "," + $n[-6:-3] + "," + $n[-3:]) elif (. | length) > 3 then ($n[:-3] + "," + $n[-3:]) else . end)",
        ""
    ' "$CLAUDE_STATS" 2>/dev/null
}

# View daily token usage
claude_daily() {
    if [ ! -f "$CLAUDE_STATS" ]; then
        echo "No stats file found"
        return 1
    fi

    local days="${1:-7}"

    echo "Daily Token Usage (last $days days)"
    echo "===================================="
    echo ""

    jq -r --argjson days "$days" '
        .dailyModelTokens | .[-$days:][] |
        "\(.date):",
        (.tokensByModel | to_entries[] | "  \(.key): \(.value) tokens"),
        ""
    ' "$CLAUDE_STATS" 2>/dev/null
}

# =============================================================================
# Permissions Management
# =============================================================================

# View current permissions
claude_permissions() {
    if [ ! -f "$CLAUDE_LOCAL" ]; then
        echo "No local settings file found at $CLAUDE_LOCAL"
        return 1
    fi

    echo "Claude Code Permissions"
    echo "======================="
    echo ""

    echo "Allowed:"
    echo "--------"
    jq -r '.permissions.allow[]? // empty' "$CLAUDE_LOCAL" 2>/dev/null | while read -r perm; do
        echo "  + $perm"
    done

    echo ""
    echo "Denied:"
    echo "-------"
    local denied
    denied=$(jq -r '.permissions.deny[]? // empty' "$CLAUDE_LOCAL" 2>/dev/null)
    if [ -z "$denied" ]; then
        echo "  (none)"
    else
        echo "$denied" | while read -r perm; do
            echo "  - $perm"
        done
    fi
}

# Add a permission rule
claude_permission_add() {
    local rule="$1"
    local type="${2:-allow}"

    if [ -z "$rule" ]; then
        echo "Usage: claude_permission_add <rule> [allow|deny]"
        echo "Example: claude_permission_add 'Bash(npm:*)' allow"
        return 1
    fi

    if [ ! -f "$CLAUDE_LOCAL" ]; then
        echo '{"permissions":{"allow":[],"deny":[]}}' > "$CLAUDE_LOCAL"
    fi

    local tmp_file
    tmp_file=$(mktemp)

    if [ "$type" = "deny" ]; then
        jq --arg rule "$rule" '.permissions.deny += [$rule] | .permissions.deny |= unique' "$CLAUDE_LOCAL" > "$tmp_file"
    else
        jq --arg rule "$rule" '.permissions.allow += [$rule] | .permissions.allow |= unique' "$CLAUDE_LOCAL" > "$tmp_file"
    fi

    mv "$tmp_file" "$CLAUDE_LOCAL"
    echo "Added $type rule: $rule"
}

# Remove a permission rule
claude_permission_remove() {
    local rule="$1"
    local type="${2:-allow}"

    if [ -z "$rule" ]; then
        echo "Usage: claude_permission_remove <rule> [allow|deny]"
        return 1
    fi

    if [ ! -f "$CLAUDE_LOCAL" ]; then
        echo "No local settings file found"
        return 1
    fi

    local tmp_file
    tmp_file=$(mktemp)

    if [ "$type" = "deny" ]; then
        jq --arg rule "$rule" '.permissions.deny -= [$rule]' "$CLAUDE_LOCAL" > "$tmp_file"
    else
        jq --arg rule "$rule" '.permissions.allow -= [$rule]' "$CLAUDE_LOCAL" > "$tmp_file"
    fi

    mv "$tmp_file" "$CLAUDE_LOCAL"
    echo "Removed $type rule: $rule"
}

# Clear all permissions
claude_permission_clear() {
    local type="${1:-all}"

    if [ ! -f "$CLAUDE_LOCAL" ]; then
        echo "No local settings file found"
        return 1
    fi

    local tmp_file
    tmp_file=$(mktemp)

    case "$type" in
        allow)
            jq '.permissions.allow = []' "$CLAUDE_LOCAL" > "$tmp_file"
            echo "Cleared all allow rules"
            ;;
        deny)
            jq '.permissions.deny = []' "$CLAUDE_LOCAL" > "$tmp_file"
            echo "Cleared all deny rules"
            ;;
        all)
            jq '.permissions = {"allow":[],"deny":[]}' "$CLAUDE_LOCAL" > "$tmp_file"
            echo "Cleared all permission rules"
            ;;
        *)
            echo "Usage: claude_permission_clear [allow|deny|all]"
            rm "$tmp_file"
            return 1
            ;;
    esac

    mv "$tmp_file" "$CLAUDE_LOCAL"
}

# =============================================================================
# Settings Management
# =============================================================================

# View current settings
claude_settings() {
    local file="${1:-global}"

    case "$file" in
        global|g)
            if [ -f "$CLAUDE_SETTINGS" ]; then
                echo "Global Settings ($CLAUDE_SETTINGS):"
                echo "===================================="
                jq '.' "$CLAUDE_SETTINGS"
            else
                echo "No global settings file found"
            fi
            ;;
        local|l)
            if [ -f "$CLAUDE_LOCAL" ]; then
                echo "Local Settings ($CLAUDE_LOCAL):"
                echo "================================"
                jq '.' "$CLAUDE_LOCAL"
            else
                echo "No local settings file found"
            fi
            ;;
        config|c)
            if [ -f "$CLAUDE_CONFIG" ]; then
                echo "Main Config ($CLAUDE_CONFIG):"
                echo "============================="
                jq '.' "$CLAUDE_CONFIG"
            else
                echo "No main config file found"
            fi
            ;;
        *)
            echo "Usage: claude_settings [global|local|config]"
            echo "  global (g) - Display settings (~/.claude/settings.json)"
            echo "  local (l)  - Display permissions (~/.claude/settings.local.json)"
            echo "  config (c) - Display main config (~/.claude.json)"
            ;;
    esac
}

# Edit settings file
claude_settings_edit() {
    local file="${1:-global}"
    local target

    case "$file" in
        global|g) target="$CLAUDE_SETTINGS" ;;
        local|l)  target="$CLAUDE_LOCAL" ;;
        config|c) target="$CLAUDE_CONFIG" ;;
        *)
            echo "Usage: claude_settings_edit [global|local|config]"
            return 1
            ;;
    esac

    if [ -f "$target" ]; then
        ${EDITOR:-vim} "$target"
    else
        echo "File not found: $target"
        return 1
    fi
}

# Set a specific setting value
claude_setting_set() {
    local key="$1"
    local value="$2"
    local file="${3:-global}"

    if [ -z "$key" ] || [ -z "$value" ]; then
        echo "Usage: claude_setting_set <key> <value> [global|local]"
        echo "Example: claude_setting_set '.display.show_tokens' true global"
        return 1
    fi

    local target
    case "$file" in
        global|g) target="$CLAUDE_SETTINGS" ;;
        local|l)  target="$CLAUDE_LOCAL" ;;
        *)
            echo "Invalid file type. Use 'global' or 'local'"
            return 1
            ;;
    esac

    if [ ! -f "$target" ]; then
        echo "{}" > "$target"
    fi

    local tmp_file
    tmp_file=$(mktemp)

    # Try to parse as JSON, fallback to string
    if echo "$value" | jq -e '.' >/dev/null 2>&1; then
        jq "$key = $value" "$target" > "$tmp_file"
    else
        jq --arg v "$value" "$key = \$v" "$target" > "$tmp_file"
    fi

    mv "$tmp_file" "$target"
    echo "Set $key = $value in $file settings"
}

# =============================================================================
# Project Management
# =============================================================================

# List all Claude projects
claude_projects() {
    if [ ! -f "$CLAUDE_CONFIG" ]; then
        echo "No Claude config found"
        return 1
    fi

    echo "Claude Code Projects"
    echo "===================="
    echo ""

    jq -r '
        .projects | to_entries[] |
        select(.value.hasTrustDialogAccepted == true) |
        "\(.key)",
        "  Sessions: \(.value.lastSessionId // "none" | if . != "none" then "active" else "none" end)",
        "  Cost: $\(.value.lastCost // 0 | . * 100 | floor / 100)",
        "  Lines: +\(.value.lastLinesAdded // 0) / -\(.value.lastLinesRemoved // 0)",
        ""
    ' "$CLAUDE_CONFIG" 2>/dev/null
}

# Get project-specific stats
claude_project() {
    local project_path="${1:-$(pwd)}"

    if [ ! -f "$CLAUDE_CONFIG" ]; then
        echo "No Claude config found"
        return 1
    fi

    echo "Project: $project_path"
    echo "========================"
    echo ""

    jq -r --arg path "$project_path" '
        .projects[$path] //
        (.projects | to_entries[] | select(.key | contains($path[-30:])) | .value) //
        "Project not found"
    ' "$CLAUDE_CONFIG" 2>/dev/null | jq '.' 2>/dev/null || echo "Project not found or no data"
}

# =============================================================================
# MCP Server Management
# =============================================================================

# List MCP servers
claude_mcp() {
    if [ ! -f "$CLAUDE_CONFIG" ]; then
        echo "No Claude config found"
        return 1
    fi

    echo "MCP Servers"
    echo "==========="
    echo ""

    # Check for global MCP servers
    local has_mcp=false

    jq -r '
        .projects | to_entries[] |
        select(.value.mcpServers | length > 0) |
        "Project: \(.key)",
        (.value.mcpServers | to_entries[] |
            "  [\(.key)]",
            "    Type: \(.value.type // "stdio")",
            "    URL: \(.value.url // .value.command // "N/A")"
        ),
        ""
    ' "$CLAUDE_CONFIG" 2>/dev/null

    # Check for .mcp.json in current directory
    if [ -f ".mcp.json" ]; then
        echo "Local MCP Config (.mcp.json):"
        echo "-----------------------------"
        jq '.' ".mcp.json"
    fi
}

# Add MCP server to current project
claude_mcp_add() {
    local name="$1"
    local url="$2"
    local type="${3:-http}"

    if [ -z "$name" ] || [ -z "$url" ]; then
        echo "Usage: claude_mcp_add <name> <url> [type]"
        echo "Example: claude_mcp_add context7 'https://mcp.context7.com/mcp' http"
        return 1
    fi

    # Create or update .mcp.json in current directory
    local mcp_file=".mcp.json"

    if [ ! -f "$mcp_file" ]; then
        echo '{"mcpServers":{}}' > "$mcp_file"
    fi

    local tmp_file
    tmp_file=$(mktemp)

    jq --arg name "$name" --arg url "$url" --arg type "$type" '
        .mcpServers[$name] = {type: $type, url: $url}
    ' "$mcp_file" > "$tmp_file"

    mv "$tmp_file" "$mcp_file"
    echo "Added MCP server '$name' to $mcp_file"
    echo "Restart Claude Code to apply changes"
}

# Remove MCP server
claude_mcp_remove() {
    local name="$1"

    if [ -z "$name" ]; then
        echo "Usage: claude_mcp_remove <name>"
        return 1
    fi

    local mcp_file=".mcp.json"

    if [ ! -f "$mcp_file" ]; then
        echo "No .mcp.json file found"
        return 1
    fi

    local tmp_file
    tmp_file=$(mktemp)

    jq --arg name "$name" 'del(.mcpServers[$name])' "$mcp_file" > "$tmp_file"
    mv "$tmp_file" "$mcp_file"
    echo "Removed MCP server '$name' from $mcp_file"
}

# =============================================================================
# Custom Commands & Agents
# =============================================================================

# List custom commands
claude_commands() {
    echo "Custom Commands"
    echo "==============="
    echo ""

    # Check user commands
    if [ -d "${HOME}/.claude/commands" ]; then
        echo "User Commands (~/.claude/commands/):"
        ls -1 "${HOME}/.claude/commands/" 2>/dev/null | while read -r cmd; do
            echo "  /$cmd"
        done
        echo ""
    fi

    # Check project commands
    if [ -d ".claude/commands" ]; then
        echo "Project Commands (.claude/commands/):"
        ls -1 ".claude/commands/" 2>/dev/null | while read -r cmd; do
            echo "  /$cmd"
        done
    fi
}

# Create a custom command
claude_command_create() {
    local name="$1"
    local scope="${2:-project}"

    if [ -z "$name" ]; then
        echo "Usage: claude_command_create <name> [project|user]"
        return 1
    fi

    local cmd_dir
    if [ "$scope" = "user" ]; then
        cmd_dir="${HOME}/.claude/commands"
    else
        cmd_dir=".claude/commands"
    fi

    mkdir -p "$cmd_dir"

    local cmd_file="$cmd_dir/${name}.md"

    if [ -f "$cmd_file" ]; then
        echo "Command already exists: $cmd_file"
        echo "Edit with: \$EDITOR $cmd_file"
        return 1
    fi

    cat > "$cmd_file" << 'EOF'
---
description: Description of what this command does
---

Your prompt instructions here.

You can use $ARGUMENTS to access command arguments.
EOF

    ${EDITOR:-vim} "$cmd_file"
    echo "Created command: /$name"
}

# List custom agents
claude_agents() {
    echo "Custom Agents"
    echo "============="
    echo ""

    # Check user agents
    if [ -d "${HOME}/.claude/agents" ]; then
        echo "User Agents (~/.claude/agents/):"
        ls -1 "${HOME}/.claude/agents/" 2>/dev/null | while read -r agent; do
            echo "  $agent"
        done
        echo ""
    fi

    # Check project agents
    if [ -d ".claude/agents" ]; then
        echo "Project Agents (.claude/agents/):"
        ls -1 ".claude/agents/" 2>/dev/null | while read -r agent; do
            echo "  $agent"
        done
    fi
}

# =============================================================================
# Utilities
# =============================================================================

# Show Claude Code info summary
claude_info() {
    echo "Claude Code Information"
    echo "======================="
    echo ""

    # Version
    if command -v claude >/dev/null 2>&1; then
        echo "Version: $(claude --version 2>/dev/null || echo 'unknown')"
    else
        echo "Version: Claude not found in PATH"
    fi
    echo ""

    # Directories
    echo "Directories:"
    echo "  Config:   $CLAUDE_DIR"
    echo "  Projects: $CLAUDE_PROJECTS"
    echo ""

    # File status
    echo "Configuration Files:"
    [ -f "$CLAUDE_CONFIG" ] && echo "  [x] ~/.claude.json" || echo "  [ ] ~/.claude.json"
    [ -f "$CLAUDE_SETTINGS" ] && echo "  [x] settings.json" || echo "  [ ] settings.json"
    [ -f "$CLAUDE_LOCAL" ] && echo "  [x] settings.local.json" || echo "  [ ] settings.local.json"
    [ -f "$CLAUDE_STATS" ] && echo "  [x] stats-cache.json" || echo "  [ ] stats-cache.json"
    echo ""

    # Quick stats
    if [ -f "$CLAUDE_STATS" ]; then
        echo "Quick Stats:"
        echo "  Sessions: $(jq -r '.totalSessions // 0' "$CLAUDE_STATS")"
        echo "  Messages: $(jq -r '.totalMessages // 0' "$CLAUDE_STATS")"
    fi
}

# Open Claude config directory
claude_dir() {
    if [ -d "$CLAUDE_DIR" ]; then
        cd "$CLAUDE_DIR" || return
        echo "Changed to $CLAUDE_DIR"
        ls -la
    else
        echo "Claude directory not found: $CLAUDE_DIR"
    fi
}

# Clear Claude cache/stats
claude_clear() {
    local what="${1:-cache}"

    case "$what" in
        stats)
            if [ -f "$CLAUDE_STATS" ]; then
                rm "$CLAUDE_STATS"
                echo "Cleared stats cache"
            fi
            ;;
        history)
            if [ -f "$CLAUDE_DIR/history.jsonl" ]; then
                rm "$CLAUDE_DIR/history.jsonl"
                echo "Cleared history"
            fi
            ;;
        cache)
            rm -rf "$CLAUDE_DIR/shell-snapshots" 2>/dev/null
            rm -rf "$CLAUDE_DIR/debug" 2>/dev/null
            echo "Cleared cache directories"
            ;;
        all)
            claude_clear stats
            claude_clear history
            claude_clear cache
            ;;
        *)
            echo "Usage: claude_clear [stats|history|cache|all]"
            ;;
    esac
}

# =============================================================================
# Help
# =============================================================================

claude_help() {
    cat << 'EOF'
Claude Code Shell Functions
===========================

Aliases:
  cc          - claude
  ccc         - claude --continue
  ccr         - claude --resume
  ccp         - claude -p (piped mode)

Stats & Usage:
  claude_stats         - View usage statistics
  claude_tokens        - View token usage by model
  claude_daily [n]     - View daily token usage (last n days)

Permissions:
  claude_permissions          - View all permissions
  claude_permission_add       - Add a permission rule
  claude_permission_remove    - Remove a permission rule
  claude_permission_clear     - Clear permissions [allow|deny|all]

Settings:
  claude_settings [type]      - View settings [global|local|config]
  claude_settings_edit [type] - Edit settings file
  claude_setting_set          - Set a specific value

Projects:
  claude_projects      - List all projects
  claude_project [path] - View project stats

MCP Servers:
  claude_mcp           - List MCP servers
  claude_mcp_add       - Add MCP server
  claude_mcp_remove    - Remove MCP server

Commands & Agents:
  claude_commands      - List custom commands
  claude_command_create - Create custom command
  claude_agents        - List custom agents

Utilities:
  claude_info          - Show Claude Code info
  claude_dir           - Open Claude config directory
  claude_clear         - Clear cache/stats
  claude_help          - Show this help
EOF
}
