#!/bin/sh
# components/tmux/functions.sh - Tmux session management functions

# =============================================================================
# Development Sessions
# =============================================================================

# Create development session with multiple panes
dev_session() {
    local session_name="${1:-dev}"

    tmux new-session -d -s "$session_name"
    tmux split-window -h -t "$session_name"
    tmux split-window -v -t "$session_name:0.1"

    tmux send-keys -t "$session_name:0.0" 'echo "Editor pane"' Enter
    tmux send-keys -t "$session_name:0.1" 'echo "Main terminal"' Enter
    tmux send-keys -t "$session_name:0.2" 'echo "Logs/monitoring pane"' Enter

    tmux select-pane -t "$session_name:0.0"
    tmux attach-session -t "$session_name"
}

# Create Kubernetes development session
k8s_session() {
    local session_name="k8s"

    tmux new-session -d -s "$session_name"
    tmux new-window -t "$session_name" -n "kubectl"
    tmux new-window -t "$session_name" -n "logs"
    tmux new-window -t "$session_name" -n "k9s"

    tmux send-keys -t "$session_name:kubectl" 'echo "kubectl commands"' Enter
    tmux split-window -h -t "$session_name:logs"
    tmux send-keys -t "$session_name:logs.0" 'echo "App logs"' Enter
    tmux send-keys -t "$session_name:logs.1" 'echo "System logs"' Enter
    tmux send-keys -t "$session_name:k9s" 'k9s' Enter

    tmux select-window -t "$session_name:1"
    tmux attach-session -t "$session_name"
}

# =============================================================================
# Project Sessions
# =============================================================================

# Create project session (auto-detects project type)
project_session() {
    local project_path="${1:-.}"
    local session_name

    session_name=$(basename "$(realpath "$project_path")")
    cd "$project_path" || return 1

    tmux new-session -d -s "$session_name" -c "$project_path"

    if [ -f "go.mod" ]; then
        tmux send-keys -t "$session_name" 'echo "Go project detected"' Enter
        tmux new-window -t "$session_name" -n "test" -c "$project_path"
        tmux new-window -t "$session_name" -n "run" -c "$project_path"
    elif [ -f "requirements.txt" ] || [ -f "pyproject.toml" ]; then
        tmux send-keys -t "$session_name" 'echo "Python project detected"' Enter
        tmux new-window -t "$session_name" -n "test" -c "$project_path"
        tmux new-window -t "$session_name" -n "server" -c "$project_path"
    elif [ -f "package.json" ]; then
        tmux send-keys -t "$session_name" 'echo "Node.js project detected"' Enter
        tmux new-window -t "$session_name" -n "dev" -c "$project_path"
        tmux new-window -t "$session_name" -n "test" -c "$project_path"
    fi

    tmux new-window -t "$session_name" -n "editor" -c "$project_path"
    tmux send-keys -t "$session_name:editor" 'code .' Enter

    tmux select-window -t "$session_name:1"
    tmux attach-session -t "$session_name"
}

# =============================================================================
# Session Management (FZF Integration)
# =============================================================================

# Quick session switcher using fzf
tswitch() {
    if ! command -v fzf >/dev/null 2>&1; then
        echo "fzf is required for tswitch"
        return 1
    fi

    local session
    session=$(tmux list-sessions -F '#S' | fzf --prompt="Switch to session: ")
    if [ -n "$session" ]; then
        tmux attach-session -t "$session"
    fi
}

# Kill session with fzf selection
tkill() {
    if ! command -v fzf >/dev/null 2>&1; then
        echo "fzf is required for tkill"
        return 1
    fi

    local session
    session=$(tmux list-sessions -F '#S' | fzf --prompt="Kill session: ")
    if [ -n "$session" ]; then
        tmux kill-session -t "$session"
        echo "Killed session: $session"
    fi
}

# =============================================================================
# TPM (Tmux Plugin Manager)
# =============================================================================

# Install TPM
tmux_install_tpm() {
    local tpm_dir="${TMUX_TPM_DIR:-$HOME/.config/tmux/plugins/tpm}"

    if [ -d "$tpm_dir" ]; then
        echo "TPM already installed at: $tpm_dir"
        echo "To update, run: tmux_update_tpm"
        return 0
    fi

    echo "Installing Tmux Plugin Manager..."
    mkdir -p "$(dirname "$tpm_dir")"
    git clone https://github.com/tmux-plugins/tpm "$tpm_dir"

    if [ -d "$tpm_dir" ]; then
        echo "TPM installed successfully!"
        echo ""
        echo "Next steps:"
        echo "  1. Start tmux: tmux"
        echo "  2. Install plugins: prefix + I"
        echo "  3. Update plugins: prefix + U"
    else
        echo "Failed to install TPM"
        return 1
    fi
}

# Update TPM
tmux_update_tpm() {
    local tpm_dir="${TMUX_TPM_DIR:-$HOME/.config/tmux/plugins/tpm}"

    if [ ! -d "$tpm_dir" ]; then
        echo "TPM not installed. Run: tmux_install_tpm"
        return 1
    fi

    echo "Updating TPM..."
    git -C "$tpm_dir" pull
    echo "TPM updated!"
}

# Install all plugins (run outside tmux)
tmux_install_plugins() {
    local tpm_dir="${TMUX_TPM_DIR:-$HOME/.config/tmux/plugins/tpm}"

    if [ ! -d "$tpm_dir" ]; then
        echo "TPM not installed. Run: tmux_install_tpm"
        return 1
    fi

    echo "Installing tmux plugins..."
    "$tpm_dir/bin/install_plugins"
}

# Update all plugins (run outside tmux)
tmux_update_plugins() {
    local tpm_dir="${TMUX_TPM_DIR:-$HOME/.config/tmux/plugins/tpm}"

    if [ ! -d "$tpm_dir" ]; then
        echo "TPM not installed. Run: tmux_install_tpm"
        return 1
    fi

    echo "Updating tmux plugins..."
    "$tpm_dir/bin/update_plugins" all
}

# =============================================================================
# Configuration Management
# =============================================================================

# Edit tmux config
tmux_config() {
    local config="${TMUX_CONF:-$HOME/.config/tmux/tmux.conf}"

    if [ ! -f "$config" ]; then
        echo "Tmux config not found: $config"
        return 1
    fi

    ${EDITOR:-vim} "$config"
    echo "Config saved. Run 'tmux source-file $config' or prefix + r to reload."
}

# Reload tmux config
tmux_reload() {
    local config="${TMUX_CONF:-$HOME/.config/tmux/tmux.conf}"

    if [ ! -f "$config" ]; then
        echo "Tmux config not found: $config"
        return 1
    fi

    if [ -n "$TMUX" ]; then
        tmux source-file "$config"
        echo "Config reloaded!"
    else
        echo "Not in a tmux session. Config will load on next tmux start."
    fi
}

# Show tmux info
tmux_info() {
    echo "Tmux Information"
    echo "================"

    if command -v tmux >/dev/null 2>&1; then
        echo "Version: $(tmux -V)"
    else
        echo "Version: not installed"
        return 1
    fi

    echo ""
    echo "Configuration:"
    echo "  Config: ${TMUX_CONF:-$HOME/.config/tmux/tmux.conf}"
    echo "  Plugins: ${TMUX_PLUGIN_DIR:-$HOME/.config/tmux/plugins}"

    local tpm_dir="${TMUX_TPM_DIR:-$HOME/.config/tmux/plugins/tpm}"
    if [ -d "$tpm_dir" ]; then
        echo "  TPM: installed"
    else
        echo "  TPM: not installed (run tmux_install_tpm)"
    fi

    echo ""
    echo "Sessions:"
    if tmux list-sessions 2>/dev/null; then
        :
    else
        echo "  No active sessions"
    fi
}

# =============================================================================
# Session Utilities
# =============================================================================

# Attach to session or create new with name
tmux_attach() {
    local session="${1:-main}"

    if tmux has-session -t "$session" 2>/dev/null; then
        tmux attach-session -t "$session"
    else
        echo "Creating new session: $session"
        tmux new-session -s "$session"
    fi
}

# =============================================================================
# Smug Session Management
# =============================================================================
# Smug provides persistent, versioned session configurations
# Install: brew install smug (or: go install github.com/ivaaaan/smug@latest)
# Configs: ~/.config/smug/*.yml

# List available smug sessions
smug_list() {
    local smug_dir="${XDG_CONFIG_HOME:-$HOME/.config}/smug"

    if ! command -v smug >/dev/null 2>&1; then
        echo "smug not installed. Install with: brew install smug"
        return 1
    fi

    echo "Available Smug Sessions"
    echo "======================="
    echo ""

    if [ -d "$smug_dir" ]; then
        for config in "$smug_dir"/*.yml; do
            [ -f "$config" ] || continue
            local name
            name=$(basename "$config" .yml)
            local desc
            desc=$(grep -m1 "^# smug session:" "$config" 2>/dev/null | sed 's/^# smug session: //')
            echo "  $name - ${desc:-No description}"
        done
    else
        echo "  No sessions found in $smug_dir"
    fi

    echo ""
    echo "Commands:"
    echo "  smug start <name>    - Start a session"
    echo "  smug stop <name>     - Stop a session"
    echo "  smug new <name>      - Create new config"
    echo "  smug edit <name>     - Edit config"
}

# Start smug session with fzf selection
smug_start() {
    local session="$1"
    local smug_dir="${XDG_CONFIG_HOME:-$HOME/.config}/smug"

    if ! command -v smug >/dev/null 2>&1; then
        echo "smug not installed. Install with: brew install smug"
        return 1
    fi

    # If no session given, use fzf to select
    if [ -z "$session" ]; then
        if ! command -v fzf >/dev/null 2>&1; then
            echo "Usage: smug_start <session_name>"
            smug_list
            return 1
        fi

        session=$(ls "$smug_dir"/*.yml 2>/dev/null | xargs -I {} basename {} .yml | fzf --prompt="Start session: ")
        [ -z "$session" ] && return 0
    fi

    shift 2>/dev/null  # Remove first arg, rest are variables
    smug start "$session" "$@"
}

# Stop smug session with fzf selection
smug_stop() {
    local session="$1"

    if ! command -v smug >/dev/null 2>&1; then
        echo "smug not installed"
        return 1
    fi

    # If no session given, use fzf to select from active sessions
    if [ -z "$session" ]; then
        if ! command -v fzf >/dev/null 2>&1; then
            echo "Usage: smug_stop <session_name>"
            return 1
        fi

        session=$(tmux list-sessions -F '#S' 2>/dev/null | fzf --prompt="Stop session: ")
        [ -z "$session" ] && return 0
    fi

    smug stop "$session"
}

# Create a new smug session config
smug_new() {
    local name="$1"
    local smug_dir="${XDG_CONFIG_HOME:-$HOME/.config}/smug"

    if [ -z "$name" ]; then
        echo "Usage: smug_new <session_name>"
        return 1
    fi

    if ! command -v smug >/dev/null 2>&1; then
        echo "smug not installed. Install with: brew install smug"
        return 1
    fi

    mkdir -p "$smug_dir"

    local config_file="$smug_dir/$name.yml"
    if [ -f "$config_file" ]; then
        echo "Config already exists: $config_file"
        printf "Edit it? [Y/n] "
        read -r response
        if [ "$response" != "n" ] && [ "$response" != "N" ]; then
            ${EDITOR:-vim} "$config_file"
        fi
        return 0
    fi

    # Create from template
    cat > "$config_file" << EOF
# smug session: $name
# Created: $(date +%Y-%m-%d)
# Usage: smug start $name

session: $name
root: ~/
attach: true

windows:
  - name: main
    commands:
      - echo "Welcome to $name session"

  - name: editor
    commands:
      - nvim

  - name: terminal
    panes:
      - type: horizontal
        commands:
          - echo "Ready"
EOF

    echo "Created: $config_file"
    ${EDITOR:-vim} "$config_file"
}

# Edit existing smug config
smug_edit() {
    local name="$1"
    local smug_dir="${XDG_CONFIG_HOME:-$HOME/.config}/smug"

    if ! command -v smug >/dev/null 2>&1; then
        echo "smug not installed"
        return 1
    fi

    # If no name given, use fzf to select
    if [ -z "$name" ]; then
        if command -v fzf >/dev/null 2>&1; then
            name=$(ls "$smug_dir"/*.yml 2>/dev/null | xargs -I {} basename {} .yml | fzf --prompt="Edit config: ")
            [ -z "$name" ] && return 0
        else
            echo "Usage: smug_edit <session_name>"
            return 1
        fi
    fi

    local config_file="$smug_dir/$name.yml"
    if [ ! -f "$config_file" ]; then
        echo "Config not found: $config_file"
        return 1
    fi

    ${EDITOR:-vim} "$config_file"
}

# Install smug
smug_install() {
    if command -v smug >/dev/null 2>&1; then
        echo "smug already installed: $(smug --version 2>/dev/null || echo 'version unknown')"
        return 0
    fi

    echo "Installing smug..."

    if command -v brew >/dev/null 2>&1; then
        brew install smug
    elif command -v go >/dev/null 2>&1; then
        go install github.com/ivaaaan/smug@latest
    else
        echo "Install manually:"
        echo "  brew install smug"
        echo "  go install github.com/ivaaaan/smug@latest"
        return 1
    fi

    # Link config directory
    smug_link_configs
}

# Link smug configs from dotfiles
smug_link_configs() {
    local source_dir="${DOTFILES_ROOT}/components/tmux/config/smug"
    local target_dir="${XDG_CONFIG_HOME:-$HOME/.config}/smug"

    if [ ! -d "$source_dir" ]; then
        echo "No smug configs in dotfiles"
        return 1
    fi

    mkdir -p "$target_dir"

    echo "Linking smug configs..."
    for config in "$source_dir"/*.yml; do
        [ -f "$config" ] || continue
        local name
        name=$(basename "$config")
        ln -sf "$config" "$target_dir/$name"
        echo "  Linked: $name"
    done
    echo "Done!"
}
