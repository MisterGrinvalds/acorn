# shell/xdg.sh - XDG Base Directory setup
# Depends on: discovery.sh (for CURRENT_PLATFORM)
#
# Exports:
#   DOTFILES_ROOT    - Location of this dotfiles repository
#   XDG_CONFIG_HOME  - User configuration directory
#   XDG_DATA_HOME    - User data directory
#   XDG_CACHE_HOME   - User cache directory
#   XDG_STATE_HOME   - User state directory
#   XDG_RUNTIME_DIR  - Runtime directory (platform-specific)

# DOTFILES_ROOT: Where this repo is cloned
# Can be overridden before sourcing
if [ -z "$DOTFILES_ROOT" ]; then
    # Default: detect from this file's location
    # Works in bash and zsh
    if [ -n "$BASH_SOURCE" ]; then
        _dotfiles_script="${BASH_SOURCE[0]}"
    elif [ -n "$ZSH_VERSION" ]; then
        _dotfiles_script="${(%):-%x}"
    else
        _dotfiles_script=""
    fi

    if [ -n "$_dotfiles_script" ] && [ -f "$_dotfiles_script" ]; then
        # Get the directory containing shell/, which is the repo root
        DOTFILES_ROOT="$(cd "$(dirname "$_dotfiles_script")/.." && pwd)"
    else
        # Fallback to XDG default
        DOTFILES_ROOT="${XDG_CONFIG_HOME:-$HOME/.config}/dotfiles"
    fi
    unset _dotfiles_script
fi
export DOTFILES_ROOT

# XDG_CONFIG_HOME: User configuration files
# Default: ~/.config
export XDG_CONFIG_HOME="${XDG_CONFIG_HOME:-$HOME/.config}"

# XDG_DATA_HOME: User data files
# Default: ~/.local/share
export XDG_DATA_HOME="${XDG_DATA_HOME:-$HOME/.local/share}"

# XDG_CACHE_HOME: User cache files
# Default: ~/.cache
export XDG_CACHE_HOME="${XDG_CACHE_HOME:-$HOME/.cache}"

# XDG_STATE_HOME: User state files (logs, history)
# Default: ~/.local/state
export XDG_STATE_HOME="${XDG_STATE_HOME:-$HOME/.local/state}"

# XDG_RUNTIME_DIR: Runtime files (sockets, etc.)
# Platform-specific
if [ -z "$XDG_RUNTIME_DIR" ]; then
    case "$CURRENT_PLATFORM" in
        darwin)
            # macOS uses $TMPDIR which is per-user
            export XDG_RUNTIME_DIR="${TMPDIR:-/tmp}"
            ;;
        linux)
            # Linux typically has /run/user/$UID
            if [ -d "/run/user/$(id -u)" ]; then
                export XDG_RUNTIME_DIR="/run/user/$(id -u)"
            else
                export XDG_RUNTIME_DIR="/tmp/runtime-$(id -u)"
            fi
            ;;
        *)
            export XDG_RUNTIME_DIR="/tmp"
            ;;
    esac
fi

# Create XDG directories if they don't exist
# Only for directories we control (not RUNTIME_DIR)
_xdg_ensure_dirs() {
    local dir
    for dir in "$XDG_CONFIG_HOME" "$XDG_DATA_HOME" "$XDG_CACHE_HOME" "$XDG_STATE_HOME"; do
        if [ ! -d "$dir" ]; then
            mkdir -p "$dir" 2>/dev/null
        fi
    done

    # Create shell-specific directories
    mkdir -p "$XDG_CONFIG_HOME/shell" 2>/dev/null
    mkdir -p "$XDG_STATE_HOME/shell" 2>/dev/null
}
_xdg_ensure_dirs
unset -f _xdg_ensure_dirs

# Set HISTFILE to XDG-compliant location
# This is early because some shells read it immediately
if [ "$CURRENT_SHELL" = "bash" ]; then
    export HISTFILE="${XDG_STATE_HOME}/shell/bash_history"
elif [ "$CURRENT_SHELL" = "zsh" ]; then
    export HISTFILE="${XDG_STATE_HOME}/shell/zsh_history"
fi
