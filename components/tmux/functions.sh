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

# Link smug configs from git repo (or dotfiles fallback)
smug_link_configs() {
    local repo_dir="${SMUG_REPO_DIR:-$HOME/.local/share/smug-sessions}"
    local target_dir="${SMUG_CONFIG_DIR:-$HOME/.config/smug}"

    # Prefer git repo if initialized
    if [ -d "$repo_dir/.git" ]; then
        echo "Linking smug configs from git repo..."
        rm -rf "$target_dir"
        ln -sf "$repo_dir" "$target_dir"
        echo "Linked: $target_dir -> $repo_dir"
        return 0
    fi

    # Fallback to dotfiles
    local source_dir="${DOTFILES_ROOT}/components/tmux/config/smug"
    if [ ! -d "$source_dir" ]; then
        echo "No smug configs found. Run smug_repo_init to clone the repo."
        return 1
    fi

    mkdir -p "$target_dir"

    echo "Linking smug configs from dotfiles..."
    for config in "$source_dir"/*.yml; do
        [ -f "$config" ] || continue
        local name
        name=$(basename "$config")
        ln -sf "$config" "$target_dir/$name"
        echo "  Linked: $name"
    done
    echo "Done!"
}

# =============================================================================
# Smug Git Sync - Cross-machine session portability
# =============================================================================
# Syncs smug session configs via git repo for cross-machine portability

# Initialize smug sessions repo (clone or update)
smug_repo_init() {
    local repo="${SMUG_REPO:-https://github.com/MisterGrinvalds/fmux.git}"
    local dir="${SMUG_REPO_DIR:-$HOME/.local/share/smug-sessions}"

    if [ -d "$dir/.git" ]; then
        echo "Smug repo already initialized at: $dir"
        echo "Pulling latest..."
        git -C "$dir" pull --rebase
    else
        echo "Cloning smug sessions repo..."
        mkdir -p "$(dirname "$dir")"
        git clone "$repo" "$dir"
    fi

    if [ -d "$dir/.git" ]; then
        echo ""
        echo "Smug repo initialized!"
        echo "Location: $dir"
        echo ""
        smug_link_configs
        echo ""
        echo "Commands:"
        echo "  smug_status     - Show repo status"
        echo "  smug_pull       - Pull latest sessions"
        echo "  smug_push       - Commit and push changes"
        echo "  smug_sync       - Full sync (pull + push)"
    else
        echo "Failed to initialize smug repo"
        return 1
    fi
}

# Show smug repo status
smug_status() {
    local dir="${SMUG_REPO_DIR:-$HOME/.local/share/smug-sessions}"

    if [ ! -d "$dir/.git" ]; then
        echo "Smug repo not initialized. Run: smug_repo_init"
        return 1
    fi

    echo "Smug Sessions Status"
    echo "===================="
    echo "Location: $dir"
    echo ""

    git -C "$dir" status --short --branch

    echo ""
    echo "Sessions:"
    for config in "$dir"/*.yml; do
        [ -f "$config" ] || continue
        local name
        name=$(basename "$config" .yml)
        echo "  - $name"
    done
}

# Pull latest sessions from remote
smug_pull() {
    local dir="${SMUG_REPO_DIR:-$HOME/.local/share/smug-sessions}"

    if [ ! -d "$dir/.git" ]; then
        echo "Smug repo not initialized. Run: smug_repo_init"
        return 1
    fi

    echo "Pulling latest sessions..."
    git -C "$dir" pull --rebase

    if [ $? -eq 0 ]; then
        echo "Sessions updated!"
    else
        echo "Pull failed. Check for conflicts."
        return 1
    fi
}

# Commit and push session changes
smug_push() {
    local dir="${SMUG_REPO_DIR:-$HOME/.local/share/smug-sessions}"
    local message="${1:-Update smug sessions}"

    if [ ! -d "$dir/.git" ]; then
        echo "Smug repo not initialized. Run: smug_repo_init"
        return 1
    fi

    # Check for changes
    if [ -z "$(git -C "$dir" status --porcelain)" ]; then
        echo "No changes to push"
        return 0
    fi

    echo "Changes to commit:"
    git -C "$dir" status --short
    echo ""

    # Add all session files
    git -C "$dir" add "*.yml" "*.yaml" README.md 2>/dev/null

    # Commit
    git -C "$dir" commit -m "$message"

    # Push
    echo "Pushing to remote..."
    git -C "$dir" push

    if [ $? -eq 0 ]; then
        echo "Sessions pushed!"
    else
        echo "Push failed. Check remote access."
        return 1
    fi
}

# Full sync: pull, commit local changes, push
smug_sync() {
    local dir="${SMUG_REPO_DIR:-$HOME/.local/share/smug-sessions}"

    if [ ! -d "$dir/.git" ]; then
        echo "Smug repo not initialized. Run: smug_repo_init"
        return 1
    fi

    echo "Syncing smug sessions..."
    echo ""

    local has_changes=false
    if [ -n "$(git -C "$dir" status --porcelain)" ]; then
        has_changes=true
        git -C "$dir" stash
    fi

    smug_pull || return 1

    if [ "$has_changes" = true ]; then
        git -C "$dir" stash pop
        smug_push "Sync local session changes"
    fi

    echo ""
    echo "Sync complete!"
}

# =============================================================================
# Window Alerts
# =============================================================================
# Trigger visual alerts on tmux window tabs that auto-clear when you switch to them

# Set alert on current window (default: red)
tmux_alert() {
    if [ -z "$TMUX" ]; then
        echo "Error: Not in a tmux session"
        return 1
    fi

    local window_id
    window_id=$(tmux display-message -p '#{window_id}')

    # Set alert flag
    tmux set-window-option -t "$window_id" @alert 1

    # Change window color to Catppuccin red (bold)
    tmux set-window-option -t "$window_id" window-status-style "fg=#f38ba8,bold,bg=#45475a"

    echo "Alert set for window $window_id"
}

# Set alert with custom color
tmux_alert_color() {
    if [ -z "$TMUX" ]; then
        echo "Error: Not in a tmux session"
        return 1
    fi

    local color="${1:-#f38ba8}"  # Default to red if no color specified
    local window_id
    window_id=$(tmux display-message -p '#{window_id}')

    tmux set-window-option -t "$window_id" @alert 1
    tmux set-window-option -t "$window_id" window-status-style "fg=$color,bold,bg=#45475a"

    echo "Alert set with color: $color"
}

# Priority alert levels using Catppuccin Mocha colors
tmux_alert_high() {
    tmux_alert_color "#f38ba8"  # red
}

tmux_alert_medium() {
    tmux_alert_color "#f9e2af"  # yellow
}

tmux_alert_low() {
    tmux_alert_color "#94e2d5"  # teal
}
