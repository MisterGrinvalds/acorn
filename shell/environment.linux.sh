#!/bin/sh
# Linux-specific environment configuration
# Requires: shell/discovery.sh, shell/xdg.sh

# Display configuration
export DISPLAY="${DISPLAY:-:0}"

# FZF location (Linuxbrew)
export FZF_LOCATION="/home/linuxbrew/.linuxbrew/opt/fzf"

# Python virtual environments location
export ENVS_LOCATION="$HOME/envs"

# Application configuration
export LESSHISTFILE="$XDG_STATE_HOME/lesshst"
export WGETRC="$XDG_CONFIG_HOME/wgetrc"
