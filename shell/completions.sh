#!/bin/sh
# Shell completion configuration
# Requires: shell/discovery.sh, shell/xdg.sh

# =============================================================================
# FZF Configuration
# =============================================================================

if command -v fzf >/dev/null 2>&1; then
    FZF_VERSION=$(fzf --version | cut -d' ' -f1)
    export FZF_VERSION

    # FZF default commands based on platform (use fd for faster search)
    case "$CURRENT_PLATFORM" in
        darwin)
            export FZF_DEFAULT_COMMAND='fd --type f --hidden --follow --exclude .git'
            ;;
        linux)
            export FZF_DEFAULT_COMMAND='fdfind --type f --hidden --follow --exclude .git'
            ;;
    esac
    export FZF_CTRL_T_COMMAND="$FZF_DEFAULT_COMMAND"

    # Alt-C: cd into directories
    case "$CURRENT_PLATFORM" in
        darwin)
            export FZF_ALT_C_COMMAND='fd --type d --hidden --follow --exclude .git'
            ;;
        linux)
            export FZF_ALT_C_COMMAND='fdfind --type d --hidden --follow --exclude .git'
            ;;
    esac

    # FZF Catppuccin Mocha theme
    export FZF_DEFAULT_OPTS="
      --extended
      --height 40%
      --layout=reverse
      --border
      --color=bg+:#313244,bg:#1e1e2e,spinner:#f5e0dc,hl:#f38ba8
      --color=fg:#cdd6f4,header:#f38ba8,info:#cba6f7,pointer:#f5e0dc
      --color=marker:#f5e0dc,fg+:#cdd6f4,prompt:#cba6f7,hl+:#f38ba8
      --bind='ctrl-/:toggle-preview'
      --preview-window='right:50%:hidden'
    "

    # CTRL-T preview (show file contents or directory tree)
    export FZF_CTRL_T_OPTS="
      --preview '[ -d {} ] && ls -la {} || head -100 {}'
    "

    # ALT-C preview (show directory contents)
    export FZF_ALT_C_OPTS="
      --preview 'ls -la {}'
    "
fi

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

# =============================================================================
# Tool-Specific Completions
# =============================================================================

# Detect Homebrew prefix
if [ -z "$HOMEBREW_PREFIX" ]; then
    if [ -d "/opt/homebrew" ]; then
        HOMEBREW_PREFIX="/opt/homebrew"
    elif [ -d "/home/linuxbrew/.linuxbrew" ]; then
        HOMEBREW_PREFIX="/home/linuxbrew/.linuxbrew"
    elif [ -d "/usr/local/Homebrew" ]; then
        HOMEBREW_PREFIX="/usr/local"
    fi
fi

if [ -n "$HOMEBREW_PREFIX" ]; then
    # Terraform completion
    if command -v terraform >/dev/null 2>&1; then
        complete -C "$HOMEBREW_PREFIX/bin/terraform" terraform 2>/dev/null
        complete -C "$HOMEBREW_PREFIX/bin/terraform" tf 2>/dev/null
    fi

    # AWS CLI completion
    if command -v aws >/dev/null 2>&1 && [ -f "$HOMEBREW_PREFIX/bin/aws_completer" ]; then
        complete -C "$HOMEBREW_PREFIX/bin/aws_completer" aws 2>/dev/null
    fi
fi

# GitHub CLI completion (dynamic)
if command -v gh >/dev/null 2>&1; then
    case "$CURRENT_SHELL" in
        bash) eval "$(gh completion -s bash 2>/dev/null)" ;;
        zsh)  eval "$(gh completion -s zsh 2>/dev/null)" ;;
    esac
fi

# kubectl completion (dynamic)
if command -v kubectl >/dev/null 2>&1; then
    case "$CURRENT_SHELL" in
        bash) eval "$(kubectl completion bash 2>/dev/null)" ;;
        zsh)  eval "$(kubectl completion zsh 2>/dev/null)" ;;
    esac
    # Also complete 'k' alias
    if [ "$CURRENT_SHELL" = "bash" ]; then
        complete -F __start_kubectl k 2>/dev/null
    fi
fi

# Helm completion (dynamic)
if command -v helm >/dev/null 2>&1; then
    case "$CURRENT_SHELL" in
        bash) eval "$(helm completion bash 2>/dev/null)" ;;
        zsh)  eval "$(helm completion zsh 2>/dev/null)" ;;
    esac
fi

# ArgoCD completion (dynamic)
if command -v argocd >/dev/null 2>&1; then
    case "$CURRENT_SHELL" in
        bash) eval "$(argocd completion bash 2>/dev/null)" ;;
        zsh)  eval "$(argocd completion zsh 2>/dev/null)" ;;
    esac
fi

# doctl completion (DigitalOcean)
if command -v doctl >/dev/null 2>&1; then
    case "$CURRENT_SHELL" in
        bash) eval "$(doctl completion bash 2>/dev/null)" ;;
        zsh)  eval "$(doctl completion zsh 2>/dev/null)" ;;
    esac
fi

# kind completion (Kubernetes in Docker)
if command -v kind >/dev/null 2>&1; then
    case "$CURRENT_SHELL" in
        bash) eval "$(kind completion bash 2>/dev/null)" ;;
        zsh)  eval "$(kind completion zsh 2>/dev/null)" ;;
    esac
fi

# kustomize completion
if command -v kustomize >/dev/null 2>&1; then
    case "$CURRENT_SHELL" in
        bash) eval "$(kustomize completion bash 2>/dev/null)" ;;
        zsh)  eval "$(kustomize completion zsh 2>/dev/null)" ;;
    esac
fi

# Vault completion
if command -v vault >/dev/null 2>&1; then
    complete -C "$HOMEBREW_PREFIX/bin/vault" vault 2>/dev/null
fi
