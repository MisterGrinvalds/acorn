#!/bin/sh
# components/git/completions.sh - Git tab completion

# Git completion is usually loaded by system bash-completion
# This file ensures it's available and adds completion for aliases

if [ "$CURRENT_SHELL" = "bash" ]; then
    # Load git completion if not already loaded
    if ! type __git_complete >/dev/null 2>&1; then
        if [ -f "/opt/homebrew/etc/bash_completion.d/git-completion.bash" ]; then
            . "/opt/homebrew/etc/bash_completion.d/git-completion.bash"
        elif [ -f "/usr/share/bash-completion/completions/git" ]; then
            . "/usr/share/bash-completion/completions/git"
        fi
    fi

    # Add completion for 'g' alias
    if type __git_complete >/dev/null 2>&1; then
        __git_complete g __git_main
    fi
fi
