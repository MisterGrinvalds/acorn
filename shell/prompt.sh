#!/bin/sh
# Terminal configuration with git-aware prompt (Catppuccin Mocha theme)
# Requires: shell/discovery.sh, shell/theme.sh

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

# Git status colors using theme variables
git_color() {
    local git_status
    git_status=$(git status --porcelain 2>/dev/null)

    # Check if we're in a git repo
    git rev-parse --git-dir >/dev/null 2>&1 || return

    if [ -n "$git_status" ]; then
        # Dirty - has uncommitted changes (Peach)
        printf '%b' "$THEME_GIT_DIRTY"
    elif git status 2>/dev/null | grep -q "Your branch is ahead"; then
        # Clean but ahead of remote (Red)
        printf '%b' "$THEME_GIT_AHEAD"
    else
        # Clean (Green)
        printf '%b' "$THEME_GIT_CLEAN"
    fi
}

# Set prompt based on shell type
case "$CURRENT_SHELL" in
    bash)
        # Bash prompt using theme colors
        # Note: \[ \] are bash-specific escapes for non-printing chars
        PS1='\['"$THEME_TEXT"'\]\n\['"$THEME_TEXT"'\]\A \['"$THEME_TEAL"'\]\u\['"$THEME_TEXT"'\] on \['"$THEME_SAPPHIRE"'\]\h\['"$THEME_TEXT"'\] \['"$THEME_BLUE"'\][\w]$(git_color)$(git_branch)\['"$THEME_RESET"'\]\n\$ '
        ;;
    zsh)
        # Zsh prompt using theme colors
        # Note: %{ %} are zsh-specific escapes for non-printing chars
        setopt PROMPT_SUBST
        PROMPT='%{'"$THEME_TEXT"'%}
%{'"$THEME_TEXT"'%}%T %{'"$THEME_TEAL"'%}%n%{'"$THEME_TEXT"'%} on %{'"$THEME_SAPPHIRE"'%}%m%{'"$THEME_TEXT"'%} %{'"$THEME_BLUE"'%}[%~]$(git_color)$(git_branch)%{'"$THEME_RESET"'%}
%# '
        ;;
esac

# Enable color support (LSCOLORS/LS_COLORS set in theme.sh)
case "$CURRENT_PLATFORM" in
    darwin)
        export CLICOLOR=1
        # LSCOLORS is set in theme.sh
        ;;
    linux)
        # GNU ls colors (LS_COLORS set in theme.sh)
        if command -v dircolors >/dev/null 2>&1; then
            eval "$(dircolors -b)"
        fi
        alias ls='ls --color=auto'
        ;;
esac
