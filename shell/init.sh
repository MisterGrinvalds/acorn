#!/bin/sh
# Shell initialization - unified entry point for bash and zsh
# This file is sourced by ~/.bashrc and ~/.zshrc
#
# Load order:
#   1. discovery.sh  - detect shell type and platform
#   2. xdg.sh        - set XDG directories
#   3. theme.sh      - Catppuccin Mocha color definitions
#   4. environment.sh - core environment variables
#   5. secrets.sh    - load secrets (silent)
#   6. options.sh    - shell options (shopt/setopt)
#   7. aliases.sh    - shell aliases
#   8. functions/**  - all function modules
#   9. completions.sh - tab completion
#  10. prompt.sh     - shell prompt
#  11. local.sh      - user overrides (optional)

# Determine DOTFILES_ROOT if not already set
if [ -z "$DOTFILES_ROOT" ]; then
    # Try to find this script's directory
    if [ -n "$BASH_SOURCE" ]; then
        DOTFILES_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
    elif [ -n "$ZSH_VERSION" ]; then
        DOTFILES_ROOT="$(cd "$(dirname "${(%):-%x}")/.." && pwd)"
    else
        DOTFILES_ROOT="${XDG_CONFIG_HOME:-$HOME/.config}/dotfiles"
    fi
    export DOTFILES_ROOT
fi

# Helper to source files safely
_source_if_exists() {
    [ -f "$1" ] && . "$1"
}

# =============================================================================
# 1. Shell and Platform Discovery (MUST be first)
# =============================================================================
_source_if_exists "$DOTFILES_ROOT/shell/discovery.sh"

# Early exit for non-interactive shells
[ "$IS_INTERACTIVE" = "false" ] && return 0

# =============================================================================
# 2. XDG Base Directory Setup
# =============================================================================
_source_if_exists "$DOTFILES_ROOT/shell/xdg.sh"

# =============================================================================
# 3. Theme (Catppuccin Mocha colors)
# =============================================================================
_source_if_exists "$DOTFILES_ROOT/shell/theme.sh"

# =============================================================================
# 4. Environment Variables
# =============================================================================
_source_if_exists "$DOTFILES_ROOT/shell/environment.sh"

# =============================================================================
# 5. Secrets (silent loading)
# =============================================================================
_source_if_exists "$DOTFILES_ROOT/shell/secrets.sh"

# =============================================================================
# 6. Shell Options
# =============================================================================
_source_if_exists "$DOTFILES_ROOT/shell/options.sh"

# =============================================================================
# 7. Aliases
# =============================================================================
_source_if_exists "$DOTFILES_ROOT/shell/aliases.sh"

# =============================================================================
# 8. Functions - load all function modules
# =============================================================================

# Core functions
for _func_file in "$DOTFILES_ROOT"/functions/core/*.sh; do
    _source_if_exists "$_func_file"
done

# Development functions
for _func_file in "$DOTFILES_ROOT"/functions/dev/*.sh; do
    _source_if_exists "$_func_file"
done

# Cloud functions
for _func_file in "$DOTFILES_ROOT"/functions/cloud/*.sh; do
    _source_if_exists "$_func_file"
done

# AI functions
for _func_file in "$DOTFILES_ROOT"/functions/ai/*.sh; do
    _source_if_exists "$_func_file"
done

unset _func_file

# =============================================================================
# 9. Completions
# =============================================================================
_source_if_exists "$DOTFILES_ROOT/shell/completions.sh"

# =============================================================================
# 10. Prompt
# =============================================================================
_source_if_exists "$DOTFILES_ROOT/shell/prompt.sh"

# =============================================================================
# 11. Local Overrides (user-specific, not in repo)
# =============================================================================
_source_if_exists "$HOME/.config/shell/local.sh"
_source_if_exists "$HOME/.shell_local"

# Cleanup helper function
unset -f _source_if_exists

# =============================================================================
# Initialization Complete
# =============================================================================
