#!/bin/sh
# core/discovery.sh - Shell and platform detection
# MUST be sourced first - everything else depends on these variables
#
# Exports:
#   CURRENT_SHELL    - bash, zsh, or unknown
#   CURRENT_PLATFORM - darwin, linux, or unknown
#   IS_INTERACTIVE   - true or false
#   IS_LOGIN_SHELL   - true or false
#   DOTFILES_ROOT    - Location of this dotfiles repository

# =============================================================================
# Shell Detection
# =============================================================================

if [ -n "$BASH_VERSION" ]; then
    export CURRENT_SHELL="bash"
elif [ -n "$ZSH_VERSION" ]; then
    export CURRENT_SHELL="zsh"
else
    export CURRENT_SHELL="unknown"
fi

# =============================================================================
# Platform Detection
# =============================================================================

case "$OSTYPE" in
    darwin*)
        export CURRENT_PLATFORM="darwin"
        ;;
    linux*)
        export CURRENT_PLATFORM="linux"
        ;;
    *)
        # Fallback to uname
        case "$(uname -s)" in
            Darwin) export CURRENT_PLATFORM="darwin" ;;
            Linux)  export CURRENT_PLATFORM="linux" ;;
            *)      export CURRENT_PLATFORM="unknown" ;;
        esac
        ;;
esac

# =============================================================================
# Interactive Shell Detection
# =============================================================================

if [ -z "$IS_INTERACTIVE" ]; then
    case "$-" in
        *i*) export IS_INTERACTIVE="true" ;;
        *)   export IS_INTERACTIVE="false" ;;
    esac
fi

# =============================================================================
# Login Shell Detection
# =============================================================================

if [ "$CURRENT_SHELL" = "bash" ]; then
    if shopt -q login_shell 2>/dev/null; then
        export IS_LOGIN_SHELL="true"
    else
        export IS_LOGIN_SHELL="false"
    fi
elif [ "$CURRENT_SHELL" = "zsh" ]; then
    if [[ -o login ]]; then
        export IS_LOGIN_SHELL="true"
    else
        export IS_LOGIN_SHELL="false"
    fi
else
    export IS_LOGIN_SHELL="unknown"
fi

# =============================================================================
# DOTFILES_ROOT Detection
# =============================================================================

if [ -z "$DOTFILES_ROOT" ]; then
    # Detect from this file's location
    if [ -n "$BASH_SOURCE" ]; then
        _dotfiles_script="${BASH_SOURCE[0]}"
    elif [ -n "$ZSH_VERSION" ]; then
        _dotfiles_script="${(%):-%x}"
    else
        _dotfiles_script=""
    fi

    if [ -n "$_dotfiles_script" ] && [ -f "$_dotfiles_script" ]; then
        # Get the directory containing core/, which is the repo root
        DOTFILES_ROOT="$(cd "$(dirname "$_dotfiles_script")/.." && pwd)"
    else
        # Fallback
        DOTFILES_ROOT="${HOME}/Repos/personal/bash-profile"
    fi
    unset _dotfiles_script
fi
export DOTFILES_ROOT

# =============================================================================
# Early Exit for Non-Interactive
# =============================================================================

# Comment this out if you need full config in scripts
if [ "$IS_INTERACTIVE" = "false" ]; then
    return 0 2>/dev/null || exit 0
fi
