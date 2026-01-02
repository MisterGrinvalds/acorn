#!/bin/sh
# core/bootstrap.sh - Main entry point for dotfiles
#
# This is the only file that needs to be sourced from ~/.bashrc or ~/.zshrc:
#   source ~/path/to/bash-profile/core/bootstrap.sh
#
# Loading sequence:
#   1. discovery.sh  - Shell/platform detection
#   2. xdg.sh        - XDG Base Directory setup
#   3. theme.sh      - Color definitions
#   4. loader.sh     - Component discovery and loading
#   5. sync.sh       - Drift detection
#   6. Local overrides

# =============================================================================
# Bootstrap Initialization
# =============================================================================

# Detect DOTFILES_ROOT from this file's location
if [ -z "$DOTFILES_ROOT" ]; then
    if [ -n "$BASH_SOURCE" ]; then
        _bootstrap_script="${BASH_SOURCE[0]}"
    elif [ -n "$ZSH_VERSION" ]; then
        _bootstrap_script="${(%):-%x}"
    else
        _bootstrap_script=""
    fi

    if [ -n "$_bootstrap_script" ] && [ -f "$_bootstrap_script" ]; then
        DOTFILES_ROOT="$(cd "$(dirname "$_bootstrap_script")/.." && pwd)"
    fi
    unset _bootstrap_script
fi

if [ -z "$DOTFILES_ROOT" ] || [ ! -d "$DOTFILES_ROOT" ]; then
    echo "Error: Could not determine DOTFILES_ROOT"
    return 1 2>/dev/null || exit 1
fi

export DOTFILES_ROOT

# =============================================================================
# Helper Function
# =============================================================================

_source_core() {
    local file="${DOTFILES_ROOT}/core/$1"
    if [ -f "$file" ]; then
        . "$file"
    else
        echo "Warning: Core module not found: $file"
    fi
}

# =============================================================================
# Core Loading Sequence
# =============================================================================

# 1. Shell and platform detection (MUST be first)
_source_core "discovery.sh"

# Exit early for non-interactive shells
[ "$IS_INTERACTIVE" = "false" ] && return 0 2>/dev/null

# 2. XDG Base Directory setup
_source_core "xdg.sh"

# 3. Theme/colors
_source_core "theme.sh"

# 4. Component loader
_source_core "loader.sh"

# 5. Sync/drift detection
_source_core "sync.sh"

# =============================================================================
# Load Components
# =============================================================================

# Run the component loader
loader_run

# =============================================================================
# Auto-Sync Check
# =============================================================================

# Run drift check on startup (if enabled)
_sync_auto_check

# =============================================================================
# Local Overrides
# =============================================================================

# Source user's local configuration if it exists
# This allows machine-specific customizations
if [ -f "${XDG_CONFIG_HOME}/shell/local.sh" ]; then
    . "${XDG_CONFIG_HOME}/shell/local.sh"
elif [ -f "${HOME}/.shell_local" ]; then
    . "${HOME}/.shell_local"
fi

# =============================================================================
# Cleanup
# =============================================================================

unset -f _source_core 2>/dev/null

# =============================================================================
# Export Key Variables
# =============================================================================

export DOTFILES_ROOT
export DOTFILES_COMPONENTS_LOADED=1
