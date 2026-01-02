#!/bin/sh
# core/xdg.sh - XDG Base Directory Specification implementation
# Depends on: discovery.sh
#
# Exports:
#   XDG_CONFIG_HOME  - User configuration directory (~/.config)
#   XDG_DATA_HOME    - User data directory (~/.local/share)
#   XDG_CACHE_HOME   - User cache directory (~/.cache)
#   XDG_STATE_HOME   - User state directory (~/.local/state)
#   XDG_RUNTIME_DIR  - Runtime directory (platform-specific)
#
# Provides functions:
#   xdg_config_dir <component>  - Get config dir for component
#   xdg_data_dir <component>    - Get data dir for component
#   xdg_cache_dir <component>   - Get cache dir for component
#   xdg_state_dir <component>   - Get state dir for component
#   xdg_ensure_dirs <component> - Create all XDG dirs for component

# =============================================================================
# XDG Base Directories
# =============================================================================

export XDG_CONFIG_HOME="${XDG_CONFIG_HOME:-$HOME/.config}"
export XDG_DATA_HOME="${XDG_DATA_HOME:-$HOME/.local/share}"
export XDG_CACHE_HOME="${XDG_CACHE_HOME:-$HOME/.cache}"
export XDG_STATE_HOME="${XDG_STATE_HOME:-$HOME/.local/state}"

# Runtime directory (platform-specific)
if [ -z "$XDG_RUNTIME_DIR" ]; then
    case "$CURRENT_PLATFORM" in
        darwin)
            export XDG_RUNTIME_DIR="${TMPDIR:-/tmp}"
            ;;
        linux)
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

# =============================================================================
# Component XDG Helper Functions
# =============================================================================

# Get XDG config directory for a component
# Usage: xdg_config_dir python  # Returns ~/.config/python
xdg_config_dir() {
    local component="$1"
    if [ -n "$component" ]; then
        echo "${XDG_CONFIG_HOME}/${component}"
    else
        echo "${XDG_CONFIG_HOME}"
    fi
}

# Get XDG data directory for a component
xdg_data_dir() {
    local component="$1"
    if [ -n "$component" ]; then
        echo "${XDG_DATA_HOME}/${component}"
    else
        echo "${XDG_DATA_HOME}"
    fi
}

# Get XDG cache directory for a component
xdg_cache_dir() {
    local component="$1"
    if [ -n "$component" ]; then
        echo "${XDG_CACHE_HOME}/${component}"
    else
        echo "${XDG_CACHE_HOME}"
    fi
}

# Get XDG state directory for a component
xdg_state_dir() {
    local component="$1"
    if [ -n "$component" ]; then
        echo "${XDG_STATE_HOME}/${component}"
    else
        echo "${XDG_STATE_HOME}"
    fi
}

# Create all XDG directories for a component
# Usage: xdg_ensure_dirs python
xdg_ensure_dirs() {
    local component="$1"

    if [ -n "$component" ]; then
        mkdir -p "${XDG_CONFIG_HOME}/${component}" 2>/dev/null
        mkdir -p "${XDG_DATA_HOME}/${component}" 2>/dev/null
        mkdir -p "${XDG_CACHE_HOME}/${component}" 2>/dev/null
        mkdir -p "${XDG_STATE_HOME}/${component}" 2>/dev/null
    else
        # Create base XDG directories
        mkdir -p "$XDG_CONFIG_HOME" 2>/dev/null
        mkdir -p "$XDG_DATA_HOME" 2>/dev/null
        mkdir -p "$XDG_CACHE_HOME" 2>/dev/null
        mkdir -p "$XDG_STATE_HOME" 2>/dev/null
    fi
}

# =============================================================================
# Initialize Base XDG Directories
# =============================================================================

xdg_ensure_dirs

# Create shell-specific directories
mkdir -p "${XDG_CONFIG_HOME}/shell" 2>/dev/null
mkdir -p "${XDG_STATE_HOME}/shell" 2>/dev/null
mkdir -p "${XDG_DATA_HOME}/shell" 2>/dev/null

# =============================================================================
# Shell History (XDG-compliant location)
# =============================================================================

if [ "$CURRENT_SHELL" = "bash" ]; then
    export HISTFILE="${XDG_STATE_HOME}/shell/bash_history"
elif [ "$CURRENT_SHELL" = "zsh" ]; then
    export HISTFILE="${XDG_STATE_HOME}/shell/zsh_history"
fi

# =============================================================================
# Warning Suppression State File
# =============================================================================

# File to track which component warnings have been shown
DOTFILES_WARNING_FILE="${XDG_STATE_HOME}/shell/component_warnings"
export DOTFILES_WARNING_FILE

# Check if warning was already shown for a component
# Usage: if xdg_warning_shown python; then ...
xdg_warning_shown() {
    local component="$1"
    [ -f "$DOTFILES_WARNING_FILE" ] && grep -q "^${component}$" "$DOTFILES_WARNING_FILE" 2>/dev/null
}

# Mark warning as shown for a component
# Usage: xdg_warning_mark python
xdg_warning_mark() {
    local component="$1"
    echo "$component" >> "$DOTFILES_WARNING_FILE"
}
