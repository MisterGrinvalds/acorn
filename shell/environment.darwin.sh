#!/bin/sh
# macOS-specific environment configuration
# Requires: shell/discovery.sh, shell/xdg.sh

# Homebrew paths
export HOMEBREW_BIN_PATH="/opt/homebrew/bin"
export HOMEBREW_CELLAR_PATH="/opt/homebrew/Cellar"

# Add Homebrew to PATH if not already present
case ":$PATH:" in
    *":$HOMEBREW_BIN_PATH:"*) ;;
    *) PATH="$HOMEBREW_BIN_PATH:$PATH" ;;
esac

# Python virtual environments location
export ENVS_LOCATION="$HOME/envs"

# Application configuration
export LESSHISTFILE="$XDG_STATE_HOME/lesshst"
export WGETRC="$XDG_CONFIG_HOME/wgetrc"

# FZF location (Homebrew)
export FZF_LOCATION="/opt/homebrew/opt/fzf"


