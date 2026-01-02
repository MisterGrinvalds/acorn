#!/bin/sh
# components/automation/env.sh - Automation framework environment

# Automation framework home
export AUTO_HOME="${DOTFILES_ROOT:-.}/.automation"
export AUTO_CLI="$AUTO_HOME/auto"

# Add automation CLI to PATH if it exists
if [ -f "$AUTO_CLI" ]; then
    case ":$PATH:" in
        *":$AUTO_HOME:"*) ;;
        *) export PATH="$AUTO_HOME:$PATH" ;;
    esac
fi
