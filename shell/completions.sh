#!/bin/sh
# Shell completion configuration
# Requires: shell/discovery.sh, shell/xdg.sh

# FZF environment configuration
export FZF_DEFAULT_OPS="--extended"

if command -v fzf >/dev/null 2>&1; then
    FZF_VERSION=$(fzf --version | cut -d' ' -f1)
    export FZF_VERSION
fi

# FZF default commands based on platform
case "$CURRENT_PLATFORM" in
    darwin)
        export FZF_DEFAULT_COMMAND="fd --type f"
        ;;
    linux)
        export FZF_DEFAULT_COMMAND="fdfind --type f"
        ;;
esac
export FZF_CTRL_T_COMMAND="$FZF_DEFAULT_COMMAND"

case "$CURRENT_SHELL" in
    bash)
        # macOS Homebrew bash-completion
        if [ -f "/opt/homebrew/etc/profile.d/bash_completion.sh" ]; then
            . "/opt/homebrew/etc/profile.d/bash_completion.sh"
        # Linux Homebrew bash-completion
        elif [ -f "/home/linuxbrew/.linuxbrew/etc/profile.d/bash_completion.sh" ]; then
            . "/home/linuxbrew/.linuxbrew/etc/profile.d/bash_completion.sh"
        # System bash-completion (Debian/Ubuntu)
        elif [ -f "/usr/share/bash-completion/bash_completion" ]; then
            . "/usr/share/bash-completion/bash_completion"
        # System bash-completion (older)
        elif [ -f "/etc/bash_completion" ]; then
            . "/etc/bash_completion"
        fi

        # FZF bash integration
        if [ -n "$FZF_LOCATION" ] && [ -d "$FZF_LOCATION" ]; then
            [ -f "$FZF_LOCATION/shell/completion.bash" ] && . "$FZF_LOCATION/shell/completion.bash"
            [ -f "$FZF_LOCATION/shell/key-bindings.bash" ] && . "$FZF_LOCATION/shell/key-bindings.bash"
        fi
        ;;

    zsh)
        # Initialize zsh completion system
        autoload -Uz compinit

        # Use cached completion dump for faster startup
        # Rebuild once a day
        _comp_cache="$XDG_CACHE_HOME/zsh/zcompdump"
        mkdir -p "$(dirname "$_comp_cache")"

        if [ -f "$_comp_cache" ] && [ "$(find "$_comp_cache" -mtime -1 2>/dev/null)" ]; then
            compinit -C -d "$_comp_cache"
        else
            compinit -d "$_comp_cache"
        fi
        unset _comp_cache

        # Completion styling
        zstyle ':completion:*' menu select
        zstyle ':completion:*' matcher-list 'm:{a-zA-Z}={A-Za-z}'  # Case insensitive
        zstyle ':completion:*' list-colors "${(s.:.)LS_COLORS}"
        zstyle ':completion:*:descriptions' format '%F{green}-- %d --%f'
        zstyle ':completion:*:warnings' format '%F{red}No matches%f'

        # FZF zsh integration
        if [ -n "$FZF_LOCATION" ] && [ -d "$FZF_LOCATION" ]; then
            [ -f "$FZF_LOCATION/shell/completion.zsh" ] && . "$FZF_LOCATION/shell/completion.zsh"
            [ -f "$FZF_LOCATION/shell/key-bindings.zsh" ] && . "$FZF_LOCATION/shell/key-bindings.zsh"
        fi
        ;;
esac

# Git completion (if not loaded by system completion)
if [ "$CURRENT_SHELL" = "bash" ]; then
    if ! type __git_complete >/dev/null 2>&1; then
        if [ -f "/opt/homebrew/etc/bash_completion.d/git-completion.bash" ]; then
            . "/opt/homebrew/etc/bash_completion.d/git-completion.bash"
        elif [ -f "/usr/share/bash-completion/completions/git" ]; then
            . "/usr/share/bash-completion/completions/git"
        fi
    fi
fi
