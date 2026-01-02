#!/bin/sh
# components/node/env.sh - Node.js environment variables

# NVM (Node Version Manager) - XDG compliant location
export NVM_DIR="${XDG_DATA_HOME:-$HOME/.local/share}/nvm"

# Load NVM if available
[ -s "$NVM_DIR/nvm.sh" ] && . "$NVM_DIR/nvm.sh"

# Load NVM bash completion
[ -s "$NVM_DIR/bash_completion" ] && . "$NVM_DIR/bash_completion"

# Ensure pnpm is globally available after NVM loads
if command -v nvm >/dev/null 2>&1 && command -v node >/dev/null 2>&1; then
    if ! command -v pnpm >/dev/null 2>&1; then
        npm install -g pnpm >/dev/null 2>&1
    fi
fi

# pnpm home directory
export PNPM_HOME="${XDG_DATA_HOME:-$HOME/.local/share}/pnpm"
case ":$PATH:" in
    *":$PNPM_HOME:"*) ;;
    *) export PATH="$PNPM_HOME:$PATH" ;;
esac

# npm cache location
export npm_config_cache="${XDG_CACHE_HOME:-$HOME/.cache}/npm"
