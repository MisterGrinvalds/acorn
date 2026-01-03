#!/bin/sh
# components/shell/completions.sh - Core shell completions

# =============================================================================
# Bash Completion
# =============================================================================

if [ "$CURRENT_SHELL" = "bash" ]; then
    # Suppress errors from completion scripts using bash 4.4+ features (like nosort)
    # when running on older bash versions (macOS ships with bash 3.2)
    # macOS Homebrew bash-completion
    if [ -f "/opt/homebrew/etc/profile.d/bash_completion.sh" ]; then
        . "/opt/homebrew/etc/profile.d/bash_completion.sh" 2>/dev/null
    # Linux Homebrew bash-completion
    elif [ -f "/home/linuxbrew/.linuxbrew/etc/profile.d/bash_completion.sh" ]; then
        . "/home/linuxbrew/.linuxbrew/etc/profile.d/bash_completion.sh" 2>/dev/null
    # System bash-completion (Debian/Ubuntu)
    elif [ -f "/usr/share/bash-completion/bash_completion" ]; then
        . "/usr/share/bash-completion/bash_completion" 2>/dev/null
    # System bash-completion (older)
    elif [ -f "/etc/bash_completion" ]; then
        . "/etc/bash_completion" 2>/dev/null
    fi
fi

# =============================================================================
# Zsh Completion
# =============================================================================

if [ "$CURRENT_SHELL" = "zsh" ]; then
    # Initialize zsh completion system
    autoload -Uz compinit

    # Use cached completion dump for faster startup
    _comp_cache="${XDG_CACHE_HOME:-$HOME/.cache}/zsh/zcompdump"
    mkdir -p "$(dirname "$_comp_cache")"

    if [ -f "$_comp_cache" ] && [ "$(find "$_comp_cache" -mtime -1 2>/dev/null)" ]; then
        compinit -C -d "$_comp_cache"
    else
        compinit -d "$_comp_cache"
    fi
    unset _comp_cache

    # Completion styling
    zstyle ':completion:*' menu select
    zstyle ':completion:*' matcher-list 'm:{a-zA-Z}={A-Za-z}'
    zstyle ':completion:*' list-colors "${(s.:.)LS_COLORS}"
    zstyle ':completion:*:descriptions' format '%F{green}-- %d --%f'
    zstyle ':completion:*:warnings' format '%F{red}No matches%f'
fi
