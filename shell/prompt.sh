#!/bin/sh
# Terminal configuration with git-aware prompt
# Requires: shell/discovery.sh

# Get current git branch (portable)
git_branch() {
    # Use git symbolic-ref for branch, git rev-parse for detached HEAD
    local branch
    branch=$(git symbolic-ref --short HEAD 2>/dev/null)
    if [ -n "$branch" ]; then
        printf "(%s)" "$branch"
        return
    fi

    # Check for detached HEAD
    local commit
    commit=$(git rev-parse --short HEAD 2>/dev/null)
    if [ -n "$commit" ]; then
        printf "(%s)" "$commit"
    fi
}

# Git status colors (portable)
git_color() {
    local git_status
    git_status=$(git status --porcelain 2>/dev/null)

    # Check if we're in a git repo
    git rev-parse --git-dir >/dev/null 2>&1 || return

    if [ -n "$git_status" ]; then
        # Dirty - has uncommitted changes
        printf '\033[1;31m'  # red
    elif git status 2>/dev/null | grep -q "Your branch is ahead"; then
        # Clean but ahead of remote
        printf '\033[1;33m'  # yellow
    else
        # Clean
        printf '\033[1;32m'  # green
    fi
}

# Set prompt based on shell type
case "$CURRENT_SHELL" in
    bash)
        # Bash prompt with \[ \] escape sequences
        PS1='\[\e[1;37m\]\n\[\e[1;37m\]\A \[\e[1;32m\]\u\[\e[1;37m\] on \[\e[1;33m\]\h\[\e[1;37m\] \[\e[1;34m\][\w]$(git_color)$(git_branch)\[\e[0m\]\n\$ '
        ;;
    zsh)
        # Zsh prompt with %{ %} escape sequences
        setopt PROMPT_SUBST
        PROMPT='%{%F{white}%}
%{%F{white}%}%T %{%F{green}%}%n%{%F{white}%} on %{%F{yellow}%}%m%{%F{white}%} %{%F{blue}%}[%~]$(git_color)$(git_branch)%{%f%}
%# '
        ;;
esac

# Enable color support
case "$CURRENT_PLATFORM" in
    darwin)
        export CLICOLOR=1
        export LSCOLORS=ExFxBxDxCxegedabagacad
        ;;
    linux)
        # GNU ls colors
        if command -v dircolors >/dev/null 2>&1; then
            eval "$(dircolors -b)"
        fi
        alias ls='ls --color=auto'
        ;;
esac
