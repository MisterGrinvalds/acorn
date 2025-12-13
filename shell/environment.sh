#!/bin/sh
# Shell-portable environment configuration
# Requires: shell/discovery.sh, shell/xdg.sh

# Python configuration
export IPYTHONDIR="$XDG_CONFIG_HOME/ipython"
export PYTHONSTARTUP="$DOTFILES_ROOT/.python/startup.py"
export JUPYTER_CONFIG_DIR="$XDG_CONFIG_HOME/jupyter"

# R configuration
export R_PROFILE="$DOTFILES_ROOT/.R/Rprofile.site"
export R_PROFILE_USER="$DOTFILES_ROOT/.R/.Rprofile"

# Neovim configuration
export NEOVIM_VIRTUALENV="$XDG_CONFIG_HOME/nvim/env"
export NVIM_LOG_FILE="$XDG_CACHE_HOME/nvim"
export VIM_PLUGGED="$XDG_DATA_HOME/nvim/plugged"

# Add neovim virtualenv to PATH if not already present
case ":$PATH:" in
    *":$NEOVIM_VIRTUALENV/bin:"*) ;;
    *) PATH="${PATH}:$NEOVIM_VIRTUALENV/bin" ;;
esac

# NVM (Node Version Manager)
export NVM_DIR="$XDG_DATA_HOME/nvm"
[ -s "$NVM_DIR/nvm.sh" ] && . "$NVM_DIR/nvm.sh"
[ -s "$NVM_DIR/bash_completion" ] && . "$NVM_DIR/bash_completion"

# Ensure pnpm is globally available after NVM loads
if command -v nvm >/dev/null 2>&1 && command -v node >/dev/null 2>&1; then
    if ! command -v pnpm >/dev/null 2>&1; then
        echo "Installing pnpm globally..."
        npm install -g pnpm >/dev/null 2>&1 && echo "pnpm installed successfully"
    fi
fi

# Platform-specific environment
case "$CURRENT_PLATFORM" in
    darwin)
        [ -f "$DOTFILES_ROOT/shell/environment.darwin.sh" ] && . "$DOTFILES_ROOT/shell/environment.darwin.sh"
        ;;
    linux)
        [ -f "$DOTFILES_ROOT/shell/environment.linux.sh" ] && . "$DOTFILES_ROOT/shell/environment.linux.sh"
        ;;
esac
