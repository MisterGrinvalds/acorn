#!/bin/sh
# components/ghostty/env.sh - Ghostty environment variables

# Ghostty config location (XDG compliant)
export GHOSTTY_CONFIG="${XDG_CONFIG_HOME:-$HOME/.config}/ghostty/config"

# Ghostty resources directory
export GHOSTTY_RESOURCES_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/ghostty"
