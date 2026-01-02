#!/bin/sh
# components/git/env.sh - Git environment

# Use XDG-compliant git config location if exists
if [ -f "${XDG_CONFIG_HOME}/git/config" ]; then
    export GIT_CONFIG_GLOBAL="${XDG_CONFIG_HOME}/git/config"
fi
