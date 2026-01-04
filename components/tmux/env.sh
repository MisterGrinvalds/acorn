#!/bin/sh
# components/tmux/env.sh - Tmux environment variables

# XDG-compliant tmux paths
export TMUX_CONFIG_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/tmux"
export TMUX_PLUGIN_DIR="${TMUX_CONFIG_DIR}/plugins"
export TMUX_TPM_DIR="${TMUX_PLUGIN_DIR}/tpm"

# Default tmux config file location (XDG)
export TMUX_CONF="${TMUX_CONFIG_DIR}/tmux.conf"

# Smug session management with git sync
export SMUG_REPO="https://github.com/MisterGrinvalds/fmux.git"
export SMUG_REPO_DIR="${XDG_DATA_HOME:-$HOME/.local/share}/smug-sessions"
export SMUG_CONFIG_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/smug"
