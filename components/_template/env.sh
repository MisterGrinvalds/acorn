#!/bin/sh
# components/_template/env.sh - Environment variables
#
# This file is sourced for ALL shells (interactive and non-interactive).
# Use it for:
#   - PATH modifications
#   - Environment variable exports
#   - Early initialization that scripts need
#
# Available variables:
#   DOTFILES_ROOT      - Root of the dotfiles repository
#   CURRENT_SHELL      - bash, zsh, or unknown
#   CURRENT_PLATFORM   - darwin, linux, or unknown
#   XDG_CONFIG_HOME    - User config directory (~/.config)
#   XDG_DATA_HOME      - User data directory (~/.local/share)
#   XDG_CACHE_HOME     - User cache directory (~/.cache)
#   XDG_STATE_HOME     - User state directory (~/.local/state)

# Example: Add tool to PATH
# export PATH="${XDG_DATA_HOME}/template/bin:${PATH}"

# Example: Set tool configuration
# export TEMPLATE_CONFIG="${XDG_CONFIG_HOME}/template"
