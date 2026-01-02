#!/bin/sh
# components/fzf/env.sh - FZF environment configuration

# =============================================================================
# FZF Version Detection
# =============================================================================

if command -v fzf >/dev/null 2>&1; then
    FZF_VERSION=$(fzf --version | cut -d' ' -f1)
    export FZF_VERSION
fi

# =============================================================================
# FZF Default Commands (use fd for faster search)
# =============================================================================

case "$CURRENT_PLATFORM" in
    darwin)
        export FZF_DEFAULT_COMMAND='fd --type f --hidden --follow --exclude .git'
        export FZF_ALT_C_COMMAND='fd --type d --hidden --follow --exclude .git'
        ;;
    linux)
        # fd is called fdfind on Debian/Ubuntu
        if command -v fdfind >/dev/null 2>&1; then
            export FZF_DEFAULT_COMMAND='fdfind --type f --hidden --follow --exclude .git'
            export FZF_ALT_C_COMMAND='fdfind --type d --hidden --follow --exclude .git'
        elif command -v fd >/dev/null 2>&1; then
            export FZF_DEFAULT_COMMAND='fd --type f --hidden --follow --exclude .git'
            export FZF_ALT_C_COMMAND='fd --type d --hidden --follow --exclude .git'
        fi
        ;;
esac

export FZF_CTRL_T_COMMAND="$FZF_DEFAULT_COMMAND"

# =============================================================================
# FZF Catppuccin Mocha Theme
# =============================================================================

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

# Preview options
export FZF_CTRL_T_OPTS="
  --preview '[ -d {} ] && ls -la {} || head -100 {}'
"

export FZF_ALT_C_OPTS="
  --preview 'ls -la {}'
"

# =============================================================================
# FZF Location Detection
# =============================================================================

if [ -z "$FZF_LOCATION" ]; then
    if [ -d "/opt/homebrew/opt/fzf" ]; then
        FZF_LOCATION="/opt/homebrew/opt/fzf"
    elif [ -d "/home/linuxbrew/.linuxbrew/opt/fzf" ]; then
        FZF_LOCATION="/home/linuxbrew/.linuxbrew/opt/fzf"
    elif [ -d "/usr/share/fzf" ]; then
        FZF_LOCATION="/usr/share/fzf"
    elif [ -d "$HOME/.fzf" ]; then
        FZF_LOCATION="$HOME/.fzf"
    fi
fi
export FZF_LOCATION
