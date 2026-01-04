#!/bin/sh
# components/git/env.sh - Git environment

# Note: We don't set GIT_CONFIG_GLOBAL because:
# 1. Git 2.0+ natively supports $XDG_CONFIG_HOME/git/config
# 2. Setting it can conflict with ~/.gitconfig symlinks
# 3. Our component deploys to ~/.gitconfig for maximum compatibility
