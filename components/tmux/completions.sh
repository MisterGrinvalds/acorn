#!/bin/sh
# components/tmux/completions.sh - Tab completions for tmux functions

# =============================================================================
# Helper: Get tmux session names
# =============================================================================

_tmux_sessions() {
    tmux list-sessions -F '#S' 2>/dev/null
}

# =============================================================================
# Bash Completions
# =============================================================================

if [ -n "$BASH_VERSION" ]; then
    # Completion for session-based functions
    _tmux_session_complete() {
        local cur="${COMP_WORDS[COMP_CWORD]}"
        COMPREPLY=($(compgen -W "$(_tmux_sessions)" -- "$cur"))
    }

    # tswitch - switch to session
    complete -F _tmux_session_complete tswitch

    # tkill - kill session
    complete -F _tmux_session_complete tkill

    # tmux_attach - attach to session
    complete -F _tmux_session_complete tmux_attach

    # dev_session - optional session name
    complete -F _tmux_session_complete dev_session

    # project_session - complete directories
    complete -d project_session
fi

# =============================================================================
# Zsh Completions
# =============================================================================

if [ -n "$ZSH_VERSION" ]; then
    # Completion for session-based functions
    _tmux_session_complete_zsh() {
        local sessions
        sessions=("${(@f)$(_tmux_sessions)}")
        _describe 'tmux sessions' sessions
    }

    # Register completions
    compdef _tmux_session_complete_zsh tswitch
    compdef _tmux_session_complete_zsh tkill
    compdef _tmux_session_complete_zsh tmux_attach
    compdef _tmux_session_complete_zsh dev_session
    compdef _files project_session
fi
