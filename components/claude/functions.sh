#!/bin/sh
# components/claude/functions.sh - Claude Code management functions
# Requires: jq

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

    local total_sessions total_messages
    total_sessions=$(jq -r '.totalSessions // 0' "$CLAUDE_STATS")
    total_messages=$(jq -r '.totalMessages // 0' "$CLAUDE_STATS")

    echo "Total Sessions: $total_sessions"
    echo "Total Messages: $total_messages"
    echo ""

    echo "Model Usage:"
    echo "------------"
    jq -r '.modelUsage | to_entries[] | "\(.key):\n  Input: \(.value.inputTokens // 0) tokens\n  Output: \(.value.outputTokens // 0) tokens"' "$CLAUDE_STATS" 2>/dev/null

    echo ""
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
        "  Input:  \(.value.inputTokens // 0)",
        "  Output: \(.value.outputTokens // 0)",
        "  Cache:  \(.value.cacheReadInputTokens // 0)",
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
    jq -r '.permissions.allow[]? // empty' "$CLAUDE_LOCAL" 2>/dev/null | while read -r perm; do
        echo "  + $perm"
    done

    echo ""
    echo "Denied:"
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
                jq '.' "$CLAUDE_SETTINGS"
            else
                echo "No global settings file found"
            fi
            ;;
        local|l)
            if [ -f "$CLAUDE_LOCAL" ]; then
                echo "Local Settings ($CLAUDE_LOCAL):"
                jq '.' "$CLAUDE_LOCAL"
            else
                echo "No local settings file found"
            fi
            ;;
        config|c)
            if [ -f "$CLAUDE_CONFIG" ]; then
                echo "Main Config ($CLAUDE_CONFIG):"
                jq '.' "$CLAUDE_CONFIG"
            else
                echo "No main config file found"
            fi
            ;;
        *)
            echo "Usage: claude_settings [global|local|config]"
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
        "  Cost: $\(.value.lastCost // 0 | . * 100 | floor / 100)",
        ""
    ' "$CLAUDE_CONFIG" 2>/dev/null
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

    if [ -f ".mcp.json" ]; then
        echo "Local MCP Config (.mcp.json):"
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
        return 1
    fi

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
}

# =============================================================================
# Custom Commands
# =============================================================================

# List custom commands
claude_commands() {
    echo "Custom Commands"
    echo "==============="
    echo ""

    if [ -d "${HOME}/.claude/commands" ]; then
        echo "User Commands (~/.claude/commands/):"
        ls -1 "${HOME}/.claude/commands/" 2>/dev/null | while read -r cmd; do
            echo "  /$cmd"
        done
        echo ""
    fi

    if [ -d ".claude/commands" ]; then
        echo "Project Commands (.claude/commands/):"
        ls -1 ".claude/commands/" 2>/dev/null | while read -r cmd; do
            echo "  /$cmd"
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

    if command -v claude >/dev/null 2>&1; then
        echo "Version: $(claude --version 2>/dev/null || echo 'unknown')"
    else
        echo "Version: Claude not found in PATH"
    fi
    echo ""

    echo "Configuration Files:"
    [ -f "$CLAUDE_CONFIG" ] && echo "  [x] ~/.claude.json" || echo "  [ ] ~/.claude.json"
    [ -f "$CLAUDE_SETTINGS" ] && echo "  [x] settings.json" || echo "  [ ] settings.json"
    [ -f "$CLAUDE_LOCAL" ] && echo "  [x] settings.local.json" || echo "  [ ] settings.local.json"
    [ -f "$CLAUDE_STATS" ] && echo "  [x] stats-cache.json" || echo "  [ ] stats-cache.json"
    echo ""

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
        cache)
            rm -rf "$CLAUDE_DIR/shell-snapshots" 2>/dev/null
            rm -rf "$CLAUDE_DIR/debug" 2>/dev/null
            echo "Cleared cache directories"
            ;;
        *)
            echo "Usage: claude_clear [stats|cache]"
            ;;
    esac
}

# =============================================================================
# Aggregation - Collect agents/commands from other repos
# =============================================================================

# Aggregate Claude agents/commands from multiple repos
claude_aggregate() {
    local search_dir="${1:-$HOME/Repos}"
    local target_dir="${DOTFILES_ROOT}/components/claude/config"
    local verbose="${2:-false}"

    if [ ! -d "$search_dir" ]; then
        echo "Directory not found: $search_dir"
        return 1
    fi

    if [ ! -d "$target_dir" ]; then
        echo "Target directory not found: $target_dir"
        echo "Ensure DOTFILES_ROOT is set correctly"
        return 1
    fi

    echo "Scanning $search_dir for Claude Code configurations..."
    echo ""

    local repos_scanned=0
    local agents_added=0
    local commands_added=0
    local subagents_added=0
    local skipped=0
    local renamed=0

    # Find all .claude directories
    find "$search_dir" -maxdepth 3 -type d -name ".claude" 2>/dev/null | while read -r claude_dir; do
        local repo_dir
        repo_dir=$(dirname "$claude_dir")
        local repo_name
        repo_name=$(basename "$repo_dir")

        # Skip the dotfiles repo itself
        if [ "$repo_dir" = "$DOTFILES_ROOT" ]; then
            continue
        fi

        local found_something=false

        # Process agents
        if [ -d "$claude_dir/agents" ]; then
            for agent_file in "$claude_dir/agents"/*.md; do
                [ -f "$agent_file" ] || continue
                local filename
                filename=$(basename "$agent_file")

                # Skip session files
                case "$filename" in
                    SESSION_CONTEXT.md|session-notes.md|*.local.md) continue ;;
                esac

                local target_file="$target_dir/agents/$filename"

                if [ -f "$target_file" ]; then
                    # Compare contents
                    if diff -q "$agent_file" "$target_file" >/dev/null 2>&1; then
                        [ "$verbose" = "true" ] && echo "  Skipped (duplicate): agents/$filename"
                        skipped=$((skipped + 1))
                    else
                        # Rename with repo prefix
                        local new_name="${repo_name}-${filename}"
                        cp "$agent_file" "$target_dir/agents/$new_name"
                        echo "  Added (renamed): agents/$new_name (from $repo_name)"
                        renamed=$((renamed + 1))
                        agents_added=$((agents_added + 1))
                    fi
                else
                    cp "$agent_file" "$target_file"
                    echo "  Added: agents/$filename (from $repo_name)"
                    agents_added=$((agents_added + 1))
                fi
                found_something=true
            done
        fi

        # Process commands
        if [ -d "$claude_dir/commands" ]; then
            for cmd_file in "$claude_dir/commands"/*.md; do
                [ -f "$cmd_file" ] || continue
                local filename
                filename=$(basename "$cmd_file")

                local target_file="$target_dir/commands/$filename"

                if [ -f "$target_file" ]; then
                    if diff -q "$cmd_file" "$target_file" >/dev/null 2>&1; then
                        [ "$verbose" = "true" ] && echo "  Skipped (duplicate): commands/$filename"
                        skipped=$((skipped + 1))
                    else
                        local new_name="${repo_name}-${filename}"
                        cp "$cmd_file" "$target_dir/commands/$new_name"
                        echo "  Added (renamed): commands/$new_name (from $repo_name)"
                        renamed=$((renamed + 1))
                        commands_added=$((commands_added + 1))
                    fi
                else
                    cp "$cmd_file" "$target_file"
                    echo "  Added: commands/$filename (from $repo_name)"
                    commands_added=$((commands_added + 1))
                fi
                found_something=true
            done
        fi

        # Process subagents
        if [ -d "$claude_dir/subagents" ]; then
            for sub_file in "$claude_dir/subagents"/*.md; do
                [ -f "$sub_file" ] || continue
                local filename
                filename=$(basename "$sub_file")

                local target_file="$target_dir/subagents/$filename"

                if [ -f "$target_file" ]; then
                    if diff -q "$sub_file" "$target_file" >/dev/null 2>&1; then
                        [ "$verbose" = "true" ] && echo "  Skipped (duplicate): subagents/$filename"
                        skipped=$((skipped + 1))
                    else
                        local new_name="${repo_name}-${filename}"
                        cp "$sub_file" "$target_dir/subagents/$new_name"
                        echo "  Added (renamed): subagents/$new_name (from $repo_name)"
                        renamed=$((renamed + 1))
                        subagents_added=$((subagents_added + 1))
                    fi
                else
                    cp "$sub_file" "$target_file"
                    echo "  Added: subagents/$filename (from $repo_name)"
                    subagents_added=$((subagents_added + 1))
                fi
                found_something=true
            done
        fi

        if [ "$found_something" = true ]; then
            repos_scanned=$((repos_scanned + 1))
        fi
    done

    echo ""
    echo "Summary:"
    echo "  Repos with configs: $repos_scanned"
    echo "  Agents added: $agents_added"
    echo "  Commands added: $commands_added"
    echo "  Subagents added: $subagents_added"
    echo "  Skipped (duplicates): $skipped"
    echo "  Renamed (conflicts): $renamed"
}

# List all agents and commands
claude_list() {
    local target_dir="${DOTFILES_ROOT}/components/claude/config"

    echo "Claude Code Agents & Commands"
    echo "=============================="
    echo ""

    if [ -d "$target_dir/agents" ]; then
        echo "Agents:"
        ls -1 "$target_dir/agents/"*.md 2>/dev/null | while read -r f; do
            echo "  $(basename "$f" .md)"
        done
        echo ""
    fi

    if [ -d "$target_dir/commands" ]; then
        echo "Commands:"
        ls -1 "$target_dir/commands/"*.md 2>/dev/null | while read -r f; do
            echo "  /$(basename "$f" .md)"
        done
        echo ""
    fi

    if [ -d "$target_dir/subagents" ]; then
        echo "Subagents:"
        ls -1 "$target_dir/subagents/"*.md 2>/dev/null | while read -r f; do
            echo "  $(basename "$f" .md)"
        done
    fi
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
  claude_daily [n]     - View daily token usage

Permissions:
  claude_permissions      - View all permissions
  claude_permission_add   - Add a permission rule
  claude_permission_remove - Remove a permission rule

Settings:
  claude_settings [type]     - View settings [global|local|config]
  claude_settings_edit [type] - Edit settings file

Projects & MCP:
  claude_projects      - List all projects
  claude_mcp           - List MCP servers
  claude_mcp_add       - Add MCP server
  claude_commands      - List custom commands

Utilities:
  claude_info          - Show Claude Code info
  claude_dir           - Open Claude config directory
  claude_clear         - Clear cache/stats
  claude_help          - Show this help

Aggregation:
  claude_aggregate [dir] - Collect agents/commands from repos
  claude_list            - List all agents and commands
EOF
}
