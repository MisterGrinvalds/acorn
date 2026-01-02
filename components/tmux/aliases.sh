#!/bin/sh
# components/tmux/aliases.sh - Tmux aliases

# Basic tmux shortcuts
alias tm='tmux'
alias tma='tmux attach-session'
alias tmat='tmux attach-session -t'
alias tmn='tmux new-session'
alias tmns='tmux new-session -s'
alias tml='tmux list-sessions'
alias tmk='tmux kill-session -t'
alias tmka='tmux kill-server'

# Attach to last session or create new
alias tmx='tmux attach-session 2>/dev/null || tmux new-session'

# Quick session access
alias tm0='tmux attach-session -t 0'
alias tm1='tmux attach-session -t 1'
alias tmdev='tmux attach-session -t dev 2>/dev/null || dev_session'
