#!/bin/sh
# core/sync.sh - Dotfiles drift detection and synchronization
# Depends on: discovery.sh, xdg.sh, theme.sh
#
# Provides functions for:
#   - Git status checking
#   - Drift detection
#   - Auto-sync on shell startup
#   - Sync logging

# =============================================================================
# Configuration
# =============================================================================

SYNC_LOG="${XDG_STATE_HOME}/shell/sync.log"
SYNC_LAST_CHECK="${XDG_STATE_HOME}/shell/sync_last_check"
SYNC_AUTO_ENABLED="${DOTFILES_AUTO_SYNC:-1}"  # Enabled by default

# =============================================================================
# Logging
# =============================================================================

_sync_log() {
    local level="$1"
    local msg="$2"
    local timestamp

    timestamp=$(date '+%Y-%m-%d %H:%M:%S')

    # Log to file
    echo "[$timestamp] [$level] $msg" >> "$SYNC_LOG" 2>/dev/null

    # Also print if interactive
    case "$level" in
        info)    printf "${THEME_INFO}[sync]${THEME_RESET} %s\n" "$msg" ;;
        warn)    printf "${THEME_WARNING}[sync]${THEME_RESET} %s\n" "$msg" ;;
        error)   printf "${THEME_ERROR}[sync]${THEME_RESET} %s\n" "$msg" ;;
        success) printf "${THEME_SUCCESS}[sync]${THEME_RESET} %s\n" "$msg" ;;
    esac
}

# =============================================================================
# Git Operations
# =============================================================================

# Check if dotfiles repo is a git repo
_sync_is_git_repo() {
    git -C "$DOTFILES_ROOT" rev-parse --is-inside-work-tree >/dev/null 2>&1
}

# Get current branch
_sync_current_branch() {
    git -C "$DOTFILES_ROOT" rev-parse --abbrev-ref HEAD 2>/dev/null
}

# Get short status
_sync_git_status() {
    git -C "$DOTFILES_ROOT" status --porcelain 2>/dev/null
}

# Check if there are uncommitted changes
_sync_has_changes() {
    [ -n "$(_sync_git_status)" ]
}

# Check commits ahead/behind
_sync_commits_ahead() {
    git -C "$DOTFILES_ROOT" rev-list --count HEAD@{upstream}..HEAD 2>/dev/null || echo "0"
}

_sync_commits_behind() {
    git -C "$DOTFILES_ROOT" rev-list --count HEAD..HEAD@{upstream} 2>/dev/null || echo "0"
}

# Fetch latest from remote (quiet)
_sync_fetch() {
    git -C "$DOTFILES_ROOT" fetch --quiet 2>/dev/null
}

# =============================================================================
# Drift Detection
# =============================================================================

# Quick drift check (for shell startup)
dotfiles_check_drift() {
    local behind ahead modified untracked

    if ! _sync_is_git_repo; then
        _sync_log warn "Dotfiles directory is not a git repository"
        return 1
    fi

    # Fetch latest (background, quiet)
    _sync_fetch

    behind=$(_sync_commits_behind)
    ahead=$(_sync_commits_ahead)

    # Count modified and untracked files
    modified=$(git -C "$DOTFILES_ROOT" status --porcelain 2>/dev/null | grep -c '^ M\|^M ')
    untracked=$(git -C "$DOTFILES_ROOT" status --porcelain 2>/dev/null | grep -c '^??')

    # Build status message
    local has_drift=false
    local status_parts=""

    if [ "$behind" -gt 0 ]; then
        status_parts="${status_parts}${behind} behind, "
        has_drift=true
    fi

    if [ "$ahead" -gt 0 ]; then
        status_parts="${status_parts}${ahead} ahead, "
        has_drift=true
    fi

    if [ "$modified" -gt 0 ]; then
        status_parts="${status_parts}${modified} modified, "
        has_drift=true
    fi

    if [ "$untracked" -gt 0 ]; then
        status_parts="${status_parts}${untracked} untracked, "
        has_drift=true
    fi

    if [ "$has_drift" = "true" ]; then
        # Remove trailing comma and space
        status_parts=$(echo "$status_parts" | sed 's/, $//')
        printf "${THEME_WARNING}[dotfiles]${THEME_RESET} drift detected: %s\n" "$status_parts"
        return 1
    fi

    return 0
}

# Full drift audit
dotfiles_audit() {
    echo ""
    echo "Dotfiles Drift Report"
    echo "====================="
    echo ""

    if ! _sync_is_git_repo; then
        echo "Error: Dotfiles directory is not a git repository"
        return 1
    fi

    # Repository info
    echo "Repository: $DOTFILES_ROOT"
    echo "Branch: $(_sync_current_branch)"
    echo ""

    # Fetch and check remote status
    _sync_fetch
    local behind=$(_sync_commits_behind)
    local ahead=$(_sync_commits_ahead)

    if [ "$behind" -gt 0 ] || [ "$ahead" -gt 0 ]; then
        printf "Remote Status: "
        [ "$behind" -gt 0 ] && printf "${THEME_WARNING}%d behind${THEME_RESET} " "$behind"
        [ "$ahead" -gt 0 ] && printf "${THEME_INFO}%d ahead${THEME_RESET}" "$ahead"
        echo ""
    else
        printf "Remote Status: ${THEME_SUCCESS}up to date${THEME_RESET}\n"
    fi
    echo ""

    # File changes
    local status
    status=$(_sync_git_status)

    if [ -n "$status" ]; then
        echo "File Changes:"
        echo "-------------"
        echo "$status" | while read -r line; do
            local prefix=$(echo "$line" | cut -c1-2)
            local file=$(echo "$line" | cut -c4-)
            case "$prefix" in
                "M "|" M") printf "  ${THEME_PEACH}M${THEME_RESET} %s\n" "$file" ;;
                "A "|" A") printf "  ${THEME_GREEN}A${THEME_RESET} %s\n" "$file" ;;
                "D "|" D") printf "  ${THEME_RED}D${THEME_RESET} %s\n" "$file" ;;
                "??")      printf "  ${THEME_MAUVE}?${THEME_RESET} %s\n" "$file" ;;
                *)         printf "  %s %s\n" "$prefix" "$file" ;;
            esac
        done
        echo ""
    else
        printf "File Changes: ${THEME_SUCCESS}none${THEME_RESET}\n\n"
    fi

    # XDG compliance check
    echo "XDG Compliance:"
    echo "---------------"

    local xdg_issues=0

    # Check for legacy dotfiles
    for legacy_file in ~/.bashrc ~/.zshrc ~/.bash_profile; do
        if [ -f "$legacy_file" ] && ! [ -L "$legacy_file" ]; then
            printf "  ${THEME_WARNING}!${THEME_RESET} %s exists (should be symlink to bootstrap)\n" "$legacy_file"
            xdg_issues=$((xdg_issues + 1))
        fi
    done

    if [ "$xdg_issues" -eq 0 ]; then
        printf "  ${THEME_SUCCESS}All checks passed${THEME_RESET}\n"
    fi
    echo ""
}

# =============================================================================
# Sync Operations
# =============================================================================

# Show detailed dotfiles status
dotfiles_status() {
    echo "Dotfiles Status"
    echo "==============="
    echo ""
    echo "Location: $DOTFILES_ROOT"
    echo "Shell: ${CURRENT_SHELL:-unknown}"
    echo "Platform: ${CURRENT_PLATFORM:-unknown}"
    echo ""

    # Check bootstrap files
    echo "Bootstrap files:"
    [ -f "$HOME/.bashrc" ] && echo "  ~/.bashrc: exists" || echo "  ~/.bashrc: missing"
    [ -f "$HOME/.zshrc" ] && echo "  ~/.zshrc: exists" || echo "  ~/.zshrc: missing"
    [ -f "$HOME/.bash_profile" ] && echo "  ~/.bash_profile: exists" || echo "  ~/.bash_profile: missing"
    echo ""

    # Check git status
    if _sync_is_git_repo; then
        echo "Git status:"
        git -C "$DOTFILES_ROOT" status --short --branch
    fi
}

# Pull latest changes
dotfiles_pull() {
    if ! _sync_is_git_repo; then
        echo "Error: Dotfiles directory is not a git repository"
        return 1
    fi

    _sync_log info "Pulling latest changes..."
    git -C "$DOTFILES_ROOT" pull --rebase

    if [ $? -eq 0 ]; then
        _sync_log success "Pull completed successfully"
    else
        _sync_log error "Pull failed"
        return 1
    fi
}

# Commit and push changes
dotfiles_push() {
    local message="${1:-Update dotfiles}"

    if ! _sync_is_git_repo; then
        echo "Error: Dotfiles directory is not a git repository"
        return 1
    fi

    if ! _sync_has_changes; then
        _sync_log info "No changes to commit"
        return 0
    fi

    _sync_log info "Committing changes..."
    git -C "$DOTFILES_ROOT" add -A
    git -C "$DOTFILES_ROOT" commit -m "$message"

    _sync_log info "Pushing to remote..."
    git -C "$DOTFILES_ROOT" push

    if [ $? -eq 0 ]; then
        _sync_log success "Push completed successfully"
    else
        _sync_log error "Push failed"
        return 1
    fi
}

# Full sync (pull then push)
dotfiles_sync() {
    dotfiles_pull || return 1

    if _sync_has_changes; then
        dotfiles_push "Sync local changes"
    fi
}

# =============================================================================
# Auto-Sync (Shell Startup)
# =============================================================================

# Run auto-sync check on shell startup
_sync_auto_check() {
    # Skip if disabled
    [ "$SYNC_AUTO_ENABLED" != "1" ] && return 0

    # Skip if not a git repo
    _sync_is_git_repo || return 0

    # Run quick drift check (non-blocking)
    dotfiles_check_drift
}

# Enable/disable auto-sync
dotfiles_auto_sync() {
    local action="${1:-status}"

    case "$action" in
        on|enable|1)
            export DOTFILES_AUTO_SYNC=1
            _sync_log success "Auto-sync enabled"
            ;;
        off|disable|0)
            export DOTFILES_AUTO_SYNC=0
            _sync_log info "Auto-sync disabled"
            ;;
        status|*)
            if [ "$DOTFILES_AUTO_SYNC" = "1" ]; then
                echo "Auto-sync: enabled"
            else
                echo "Auto-sync: disabled"
            fi
            ;;
    esac
}

# =============================================================================
# Bootstrap Management (from inject.sh)
# =============================================================================

# Install bootstrap files to user's home
dotfiles_inject() {
    echo "Installing dotfiles bootstrap..."
    echo "Dotfiles location: $DOTFILES_ROOT"

    # Create XDG directories
    mkdir -p "${XDG_CONFIG_HOME:-$HOME/.config}"
    mkdir -p "${XDG_DATA_HOME:-$HOME/.local/share}"
    mkdir -p "${XDG_CACHE_HOME:-$HOME/.cache}"
    mkdir -p "${XDG_STATE_HOME:-$HOME/.local/state}"

    # Create .bashrc bootstrap
    if [ ! -f "$HOME/.bashrc" ] || ! grep -q "DOTFILES_ROOT" "$HOME/.bashrc" 2>/dev/null; then
        echo "Creating ~/.bashrc bootstrap..."
        cat > "$HOME/.bashrc" << EOF
# Dotfiles bootstrap - component-based architecture
export DOTFILES_ROOT="$DOTFILES_ROOT"
[ -f "\$DOTFILES_ROOT/core/bootstrap.sh" ] && . "\$DOTFILES_ROOT/core/bootstrap.sh"
EOF
        echo "Created ~/.bashrc"
    else
        echo "~/.bashrc already configured"
    fi

    # Create .zshrc bootstrap
    if [ ! -f "$HOME/.zshrc" ] || ! grep -q "DOTFILES_ROOT" "$HOME/.zshrc" 2>/dev/null; then
        echo "Creating ~/.zshrc bootstrap..."
        cat > "$HOME/.zshrc" << EOF
# Dotfiles bootstrap - component-based architecture
export DOTFILES_ROOT="$DOTFILES_ROOT"
[ -f "\$DOTFILES_ROOT/core/bootstrap.sh" ] && . "\$DOTFILES_ROOT/core/bootstrap.sh"
EOF
        echo "Created ~/.zshrc"
    else
        echo "~/.zshrc already configured"
    fi

    # Create .bash_profile that sources .bashrc
    if [ ! -f "$HOME/.bash_profile" ] || ! grep -q "bashrc" "$HOME/.bash_profile" 2>/dev/null; then
        echo "Creating ~/.bash_profile..."
        cat > "$HOME/.bash_profile" << 'EOF'
# Source .bashrc for login shells
[ -f "$HOME/.bashrc" ] && . "$HOME/.bashrc"
EOF
        echo "Created ~/.bash_profile"
    else
        echo "~/.bash_profile already configured"
    fi

    echo ""
    echo "Bootstrap installation complete!"
    echo "Restart your shell or run: source ~/.bashrc"
}

# Remove all injected configuration
dotfiles_eject() {
    echo "Removing dotfiles bootstrap..."

    local files_removed=0

    # Remove .bashrc if it's our bootstrap
    if [ -f "$HOME/.bashrc" ] && grep -q "DOTFILES_ROOT" "$HOME/.bashrc" 2>/dev/null; then
        rm "$HOME/.bashrc"
        echo "Removed ~/.bashrc"
        files_removed=$((files_removed + 1))
    fi

    # Remove .zshrc if it's our bootstrap
    if [ -f "$HOME/.zshrc" ] && grep -q "DOTFILES_ROOT" "$HOME/.zshrc" 2>/dev/null; then
        rm "$HOME/.zshrc"
        echo "Removed ~/.zshrc"
        files_removed=$((files_removed + 1))
    fi

    # Remove .bash_profile if it's our bootstrap
    if [ -f "$HOME/.bash_profile" ] && grep -q "bashrc" "$HOME/.bash_profile" 2>/dev/null; then
        rm "$HOME/.bash_profile"
        echo "Removed ~/.bash_profile"
        files_removed=$((files_removed + 1))
    fi

    if [ $files_removed -eq 0 ]; then
        echo "No bootstrap files found to remove"
    else
        echo ""
        echo "Removed $files_removed bootstrap file(s)"
        echo "Your shell configuration has been reset"
    fi
}

# Update dotfiles from git and reload
dotfiles_update() {
    echo "Updating dotfiles..."

    if _sync_is_git_repo; then
        dotfiles_pull || return 1
        echo ""
        echo "Dotfiles updated. Reloading..."
        dotfiles_reload
    else
        echo "Not a git repository: $DOTFILES_ROOT"
        return 1
    fi
}

# Reload shell configuration without restart
dotfiles_reload() {
    echo "Reloading shell configuration..."

    if [ -f "$DOTFILES_ROOT/core/bootstrap.sh" ]; then
        . "$DOTFILES_ROOT/core/bootstrap.sh"
        echo "Configuration reloaded"
    else
        echo "bootstrap.sh not found at $DOTFILES_ROOT/core/bootstrap.sh"
        return 1
    fi
}

# =============================================================================
# Config Linking
# =============================================================================

# Symlink app configs to their expected locations
dotfiles_link_configs() {
    echo "Linking app configurations..."

    # Git config
    if [ -f "$DOTFILES_ROOT/config/git/config" ]; then
        ln -sf "$DOTFILES_ROOT/config/git/config" "$HOME/.gitconfig"
        echo "Linked: ~/.gitconfig"
    fi

    # SSH config (be careful with permissions)
    if [ -f "$DOTFILES_ROOT/config/ssh/config" ]; then
        mkdir -p "$HOME/.ssh"
        chmod 700 "$HOME/.ssh"
        ln -sf "$DOTFILES_ROOT/config/ssh/config" "$HOME/.ssh/config"
        chmod 600 "$HOME/.ssh/config"
        echo "Linked: ~/.ssh/config"
    fi

    # Conda config
    if [ -f "$DOTFILES_ROOT/config/conda/condarc" ]; then
        ln -sf "$DOTFILES_ROOT/config/conda/condarc" "$HOME/.condarc"
        echo "Linked: ~/.condarc"
    fi

    # Karabiner config (macOS only)
    if [ "$CURRENT_PLATFORM" = "darwin" ] && [ -f "$DOTFILES_ROOT/config/karabiner/karabiner.json" ]; then
        mkdir -p "$HOME/.config/karabiner"
        ln -sf "$DOTFILES_ROOT/config/karabiner/karabiner.json" "$HOME/.config/karabiner/karabiner.json"
        echo "Linked: ~/.config/karabiner/karabiner.json"
    fi

    # Ghostty config
    if [ -f "$DOTFILES_ROOT/config/ghostty/config" ]; then
        mkdir -p "$HOME/.config/ghostty"
        ln -sf "$DOTFILES_ROOT/config/ghostty/config" "$HOME/.config/ghostty/config"
        echo "Linked: ~/.config/ghostty/config"
    fi

    # Tmux config
    if [ -f "$DOTFILES_ROOT/config/tmux/tmux.conf" ]; then
        mkdir -p "$HOME/.config/tmux"
        ln -sf "$DOTFILES_ROOT/config/tmux/tmux.conf" "$HOME/.config/tmux/tmux.conf"
        echo "Linked: ~/.config/tmux/tmux.conf"
    fi

    # iTerm2 config (macOS only)
    if [ "$CURRENT_PLATFORM" = "darwin" ] && [ -d "$DOTFILES_ROOT/config/iterm2" ]; then
        local iterm_profiles_dir="$HOME/Library/Application Support/iTerm2/DynamicProfiles"
        mkdir -p "$iterm_profiles_dir"

        # Link Dynamic Profiles
        if [ -d "$DOTFILES_ROOT/config/iterm2/DynamicProfiles" ]; then
            for profile in "$DOTFILES_ROOT/config/iterm2/DynamicProfiles"/*.json; do
                [ -f "$profile" ] || continue
                local profile_name
                profile_name=$(basename "$profile")
                ln -sf "$profile" "$iterm_profiles_dir/$profile_name"
                echo "Linked: iTerm2 DynamicProfiles/$profile_name"
            done
        fi

        # Note about color scheme
        if [ -f "$DOTFILES_ROOT/config/iterm2/catppuccin-mocha.itermcolors" ]; then
            echo "iTerm2 color scheme available: $DOTFILES_ROOT/config/iterm2/catppuccin-mocha.itermcolors"
            echo "  Import via: iTerm2 > Preferences > Profiles > Colors > Color Presets > Import"
        fi
    fi

    # VS Code config
    if [ -d "$DOTFILES_ROOT/config/vscode" ]; then
        local vscode_dir
        if [ "$CURRENT_PLATFORM" = "darwin" ]; then
            vscode_dir="$HOME/Library/Application Support/Code/User"
        else
            vscode_dir="$HOME/.config/Code/User"
        fi
        mkdir -p "$vscode_dir"

        if [ -f "$DOTFILES_ROOT/config/vscode/settings.json" ]; then
            ln -sf "$DOTFILES_ROOT/config/vscode/settings.json" "$vscode_dir/settings.json"
            echo "Linked: VS Code settings.json"
        fi
        if [ -f "$DOTFILES_ROOT/config/vscode/keybindings.json" ]; then
            ln -sf "$DOTFILES_ROOT/config/vscode/keybindings.json" "$vscode_dir/keybindings.json"
            echo "Linked: VS Code keybindings.json"
        fi
    fi

    # IntelliJ config (reference only - requires manual setup due to version directories)
    if [ -d "$DOTFILES_ROOT/config/intellij" ]; then
        echo "IntelliJ configs available at: $DOTFILES_ROOT/config/intellij/"
        echo "  - Manual linking required due to version-specific paths"
    fi

    # Claude Code config
    if [ -f "$DOTFILES_ROOT/config/claude/settings.json" ]; then
        mkdir -p "$HOME/.claude"
        ln -sf "$DOTFILES_ROOT/config/claude/settings.json" "$HOME/.claude/settings.json"
        echo "Linked: ~/.claude/settings.json"
    fi

    echo ""
    echo "Config linking complete!"
}

# Unlink app configs
dotfiles_unlink_configs() {
    echo "Unlinking app configurations..."

    [ -L "$HOME/.gitconfig" ] && rm "$HOME/.gitconfig" && echo "Unlinked: ~/.gitconfig"
    [ -L "$HOME/.ssh/config" ] && rm "$HOME/.ssh/config" && echo "Unlinked: ~/.ssh/config"
    [ -L "$HOME/.condarc" ] && rm "$HOME/.condarc" && echo "Unlinked: ~/.condarc"
    [ -L "$HOME/.config/karabiner/karabiner.json" ] && rm "$HOME/.config/karabiner/karabiner.json" && echo "Unlinked: ~/.config/karabiner/karabiner.json"
    [ -L "$HOME/.config/ghostty/config" ] && rm "$HOME/.config/ghostty/config" && echo "Unlinked: ~/.config/ghostty/config"
    [ -L "$HOME/.config/tmux/tmux.conf" ] && rm "$HOME/.config/tmux/tmux.conf" && echo "Unlinked: ~/.config/tmux/tmux.conf"
    [ -L "$HOME/.claude/settings.json" ] && rm "$HOME/.claude/settings.json" && echo "Unlinked: ~/.claude/settings.json"

    # iTerm2 DynamicProfiles (macOS only)
    if [ "$CURRENT_PLATFORM" = "darwin" ]; then
        local iterm_profiles_dir="$HOME/Library/Application Support/iTerm2/DynamicProfiles"
        if [ -d "$iterm_profiles_dir" ]; then
            for profile in "$iterm_profiles_dir"/*.json; do
                [ -L "$profile" ] && rm "$profile" && echo "Unlinked: iTerm2 $(basename "$profile")"
            done
        fi
    fi

    # VS Code configs
    local vscode_dir
    if [ "$CURRENT_PLATFORM" = "darwin" ]; then
        vscode_dir="$HOME/Library/Application Support/Code/User"
    else
        vscode_dir="$HOME/.config/Code/User"
    fi
    [ -L "$vscode_dir/settings.json" ] && rm "$vscode_dir/settings.json" && echo "Unlinked: VS Code settings.json"
    [ -L "$vscode_dir/keybindings.json" ] && rm "$vscode_dir/keybindings.json" && echo "Unlinked: VS Code keybindings.json"

    echo ""
    echo "Config unlinking complete!"
}

# =============================================================================
# Help
# =============================================================================

# Help function
dotfiles_help() {
    cat << 'EOF'
Dotfiles Management Commands
============================

INSTALLATION:
  dotfiles_inject         Install bootstrap files (~/.bashrc, ~/.zshrc)
  dotfiles_eject          Remove all bootstrap files
  dotfiles_link_configs   Symlink app configs (git, ssh, etc.)
  dotfiles_unlink_configs Remove app config symlinks

MANAGEMENT:
  dotfiles_update         Git pull and reload configuration
  dotfiles_reload         Reload without restart
  dotfiles_status         Show current status
  dotfiles_audit          Full drift report

SYNC:
  dotfiles_pull           Pull latest changes from remote
  dotfiles_push           Commit and push local changes
  dotfiles_sync           Pull then push (full sync)
  dotfiles_check_drift    Quick drift check
  dotfiles_auto_sync      Enable/disable auto-sync on startup

EXAMPLES:
  dotfiles_inject         # First-time setup
  dotfiles_link_configs   # Link git, ssh configs
  dotfiles_update         # Update from git
  dotfiles_status         # Check current state
  dotfiles_audit          # Full drift report
EOF
}

# Convenience aliases
alias df-inject='dotfiles_inject'
alias df-eject='dotfiles_eject'
alias df-update='dotfiles_update'
alias df-reload='dotfiles_reload'
alias df-status='dotfiles_status'
alias df-audit='dotfiles_audit'
alias df-link='dotfiles_link_configs'
alias df-unlink='dotfiles_unlink_configs'
alias df-help='dotfiles_help'

# =============================================================================
# Initialize
# =============================================================================

# Create sync log directory
mkdir -p "$(dirname "$SYNC_LOG")" 2>/dev/null
