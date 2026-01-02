#!/bin/sh
# components/tmux/env.sh - Tmux environment variables

# XDG-compliant tmux paths
export TMUX_CONFIG_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/tmux"
export TMUX_PLUGIN_DIR="${TMUX_CONFIG_DIR}/plugins"
export TMUX_TPM_DIR="${TMUX_PLUGIN_DIR}/tpm"

# Default tmux config file location (XDG)
export TMUX_CONF="${TMUX_CONFIG_DIR}/tmux.conf"
