#!/bin/sh
# components/shell/env.sh - Core shell environment
#
# Shell options and history configuration

# =============================================================================
# Shell Options
# =============================================================================

case "$CURRENT_SHELL" in
    bash)
        # History options
        shopt -s histappend         # Append to history, don't overwrite
        shopt -s cmdhist            # Save multi-line commands as single entry

        # Directory navigation
        shopt -s autocd 2>/dev/null # cd by typing directory name (bash 4+)
        shopt -s cdspell            # Autocorrect typos in cd
        shopt -s dirspell 2>/dev/null  # Autocorrect directory names (bash 4+)

        # Globbing
        shopt -s globstar 2>/dev/null  # ** matches recursively (bash 4+)
        shopt -s nocaseglob         # Case-insensitive globbing
        shopt -s extglob            # Extended pattern matching

        # Misc
        shopt -s checkwinsize       # Update LINES/COLUMNS after each command
        shopt -s no_empty_cmd_completion  # Don't complete on empty line

        # History control
        export HISTCONTROL="erasedups:ignoredups:ignorespace"
        export HISTSIZE=10000
        export HISTFILESIZE=20000
        export HISTIGNORE="ls:ll:cd:pwd:bg:fg:history:clear:exit"
        ;;

    zsh)
        # History options
        setopt APPEND_HISTORY       # Append to history, don't overwrite
        setopt HIST_IGNORE_DUPS     # Don't store duplicates
        setopt HIST_IGNORE_SPACE    # Don't store commands starting with space
        setopt HIST_EXPIRE_DUPS_FIRST  # Remove duplicates first when trimming
        setopt HIST_FIND_NO_DUPS    # Don't show duplicates in search
        setopt SHARE_HISTORY        # Share history between sessions
        setopt EXTENDED_HISTORY     # Save timestamp

        # Directory navigation
        setopt AUTO_CD              # cd by typing directory name
        setopt AUTO_PUSHD           # Push directories onto stack
        setopt PUSHD_IGNORE_DUPS    # Don't push duplicates
        setopt PUSHD_SILENT         # Don't print directory stack
        setopt CDABLE_VARS          # cd to named directories

        # Globbing
        setopt EXTENDED_GLOB        # Extended pattern matching
        setopt NO_CASE_GLOB         # Case-insensitive globbing
        setopt GLOB_DOTS            # Include dotfiles in glob

        # Misc
        setopt CORRECT              # Spelling correction for commands
        setopt NO_BEEP              # Don't beep on errors
        setopt INTERACTIVE_COMMENTS # Allow comments in interactive shell

        # History settings
        export HISTSIZE=10000
        export SAVEHIST=20000
        ;;
esac

# =============================================================================
# Color Support
# =============================================================================

case "$CURRENT_PLATFORM" in
    darwin)
        export CLICOLOR=1
        ;;
    linux)
        if command -v dircolors >/dev/null 2>&1; then
            eval "$(dircolors -b)"
        fi
        ;;
esac

# =============================================================================
# Default Editor
# =============================================================================

if command -v nvim >/dev/null 2>&1; then
    export EDITOR='nvim'
    export VISUAL='nvim'
elif command -v vim >/dev/null 2>&1; then
    export EDITOR='vim'
    export VISUAL='vim'
else
    export EDITOR='vi'
    export VISUAL='vi'
fi

# =============================================================================
# Pager
# =============================================================================

export PAGER='less'
export LESS='-R -F -X'

# =============================================================================
# Claude / Anthropic
# =============================================================================

export ANTHROPIC_MODEL='claude-opus-4-5-20251101[1m]'

# =============================================================================
# FZF Keybindings
# =============================================================================

if command -v fzf >/dev/null 2>&1; then
    case "$CURRENT_SHELL" in
        bash)
            # Ctrl+r: fzf history search
            __fzf_history__() {
                local selected
                selected=$(history | sed 's/^[ ]*[0-9]*[ ]*//' | fzf --tac --no-sort --reverse --query="$READLINE_LINE")
                READLINE_LINE="$selected"
                READLINE_POINT=${#READLINE_LINE}
            }
            bind -x '"\C-r": __fzf_history__'
            ;;
        zsh)
            # Ctrl+r: fzf history search
            fzf-history-widget() {
                local selected
                selected=$(fc -l 1 | sed 's/^[ ]*[0-9]*[ ]*//' | fzf --tac --no-sort --reverse --query="$LBUFFER")
                if [[ -n "$selected" ]]; then
                    LBUFFER="$selected"
                fi
                zle reset-prompt
            }
            zle -N fzf-history-widget
            bindkey '^R' fzf-history-widget
            ;;
    esac
fi
