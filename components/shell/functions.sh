#!/bin/sh
# components/shell/functions.sh - Core shell functions

# =============================================================================
# Git Prompt Helpers
# =============================================================================

# Get current git branch (portable)
git_branch() {
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

# =============================================================================
# Prompt Setup
# =============================================================================

# Set prompt based on shell type
_setup_prompt() {
    case "$CURRENT_SHELL" in
        bash)
            # Bash prompt using theme colors
            PS1='\['"$THEME_TEXT"'\]\n\['"$THEME_TEXT"'\]\A \['"$THEME_TEAL"'\]\u\['"$THEME_TEXT"'\] on \['"$THEME_SAPPHIRE"'\]\h\['"$THEME_TEXT"'\] \['"$THEME_BLUE"'\][\w]$(git_color)$(git_branch)\['"$THEME_RESET"'\]\n\$ '
            ;;
        zsh)
            # Zsh prompt using theme colors
            setopt PROMPT_SUBST
            PROMPT='%{'"$THEME_TEXT"'%}
%{'"$THEME_TEXT"'%}%T %{'"$THEME_TEAL"'%}%n%{'"$THEME_TEXT"'%} on %{'"$THEME_SAPPHIRE"'%}%m%{'"$THEME_TEXT"'%} %{'"$THEME_BLUE"'%}[%~]$(git_color)$(git_branch)%{'"$THEME_RESET"'%}
%# '
            ;;
    esac
}

# Initialize prompt
_setup_prompt

# =============================================================================
# Directory Operations
# =============================================================================

# Enhanced cd - automatically lists directory after changing
cd() {
    builtin cd "$@" && ll
}

# Make directory and cd into it
mkcd() {
    mkdir -p "$1" && cd "$1"
}

# Go up N directories
up() {
    local count="${1:-1}"
    local path=""
    for _ in $(seq 1 "$count"); do
        path="${path}../"
    done
    cd "$path" || return
}

# =============================================================================
# File Operations
# =============================================================================

# Extract various archive formats
extract() {
    if [ -z "$1" ]; then
        echo "Usage: extract <file>"
        return 1
    fi

    if [ ! -f "$1" ]; then
        echo "'$1' is not a valid file"
        return 1
    fi

    case "$1" in
        *.tar.bz2) tar xjf "$1" ;;
        *.tar.gz)  tar xzf "$1" ;;
        *.tar.xz)  tar xJf "$1" ;;
        *.bz2)     bunzip2 "$1" ;;
        *.gz)      gunzip "$1" ;;
        *.tar)     tar xf "$1" ;;
        *.tbz2)    tar xjf "$1" ;;
        *.tgz)     tar xzf "$1" ;;
        *.zip)     unzip "$1" ;;
        *.Z)       uncompress "$1" ;;
        *.7z)      7z x "$1" ;;
        *.rar)     unrar x "$1" ;;
        *)         echo "'$1' cannot be extracted via extract()" ;;
    esac
}

# Create a tar.gz archive (strips trailing slashes)
mktar() {
    tar -czvf "${1%%/}.tar.gz" "${1%%/}/"
}

# Create a zip archive (strips trailing slashes)
mkzip() {
    zip -r "${1%%/}.zip" "${1%%/}/"
}

# =============================================================================
# History
# =============================================================================

# Shorthand for history with grep
h() {
    history | grep "$1"
}

# =============================================================================
# Search
# =============================================================================

# Find file by name
ff() {
    find . -type f -iname "*$1*"
}

# Find directory by name
fd() {
    find . -type d -iname "*$1*"
}

# Grep recursively
rg() {
    grep -rn "$1" .
}

# =============================================================================
# System Info
# =============================================================================

# Show system info
sysinfo() {
    echo "Hostname:  $(hostname)"
    echo "OS:        $(uname -s) $(uname -r)"
    echo "Shell:     $CURRENT_SHELL"
    echo "Platform:  $CURRENT_PLATFORM"
    echo "User:      $(whoami)"
    echo "Home:      $HOME"
    echo "Dotfiles:  $DOTFILES_ROOT"
}

# =============================================================================
# Shell Utilities
# =============================================================================

# Run a bash shell as another user
bash_as() {
    sudo -u "$1" /bin/bash
}

# Remove all environment variables matching pattern
rmenv() {
    unset $(env | grep -i "${1:-prox}" | grep -oE '^[^=]+')
}
