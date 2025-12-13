#!/bin/sh
# Shell-portable aliases
# Requires: shell/discovery.sh

# Shell reload
alias resource='source ~/.bashrc 2>/dev/null || source ~/.zshrc 2>/dev/null'

# Git shortcuts
alias g='git'
alias gs='git status'
alias ga='git add'
alias gc='git commit'
alias gp='git push'
alias gl='git pull'
alias gd='git diff'
alias gco='git checkout'
alias gb='git branch'
alias glog='git log --oneline --graph --decorate'

# Navigation
alias ..='cd ..'
alias ...='cd ../..'
alias ....='cd ../../..'
alias ll='ls -alh'
alias llr='ls -alhr'
alias lls='ls -alhS'
alias llsr='ls -alhSr'
alias lld='ls -alht'
alias lldr='ls -alhtr'
alias mkdir='mkdir -pv'

# Python development
alias pip='python -m pip'
alias pip3='python3 -m pip'
alias py='python'
alias py3='python3'
alias ipy='ipython'
alias ptest='python -m pytest'
alias ptestv='python -m pytest -v'
alias black='python -m black'
alias isort='python -m isort'
alias flake8='python -m flake8'

# FastAPI development
alias uvdev='uvicorn main:app --reload'
alias uvprod='uvicorn main:app --host 0.0.0.0 --port 8000'

# Tmux enhanced
alias ta='tmux attach-session -t'
alias tn='tmux new-session -s'
alias tk='tmux kill-session -t'
alias tko='tmux kill-session -a'
alias ti='tmux info'
alias ts='tmux list-sessions'
alias tks='tmux kill-server'
alias td='tmux detach'

# Tmux project sessions
alias twork='tmux new-session -s work -d'
alias tdev='tmux new-session -s dev -d'
alias tk8s='tmux new-session -s k8s -d'

# Platform-specific aliases
case "$CURRENT_PLATFORM" in
    darwin)
        alias getsshkey='pbcopy < ~/.ssh/id_rsa.pub'
        alias perm="stat -f '%Lp'"
        alias lldc='ls -alhtU'       # List by date created (macOS only)
        alias lldcr='ls -alhtUr'
        ;;
    linux)
        alias getsshkey='xclip -selection clipboard < ~/.ssh/id_rsa.pub'
        alias perm='stat -c "%a"'
        # xclip shortcuts
        alias c='xclip'
        alias cs='xclip -selection clipboard'
        alias v='xclip -o'
        alias vs='xclip -o -selection clipboard'
        ;;
esac

# Tree alternative using ls
alias tree='ls -R | grep ":$" | sed -e "s/:$//" -e "s/[^-][^\/]*\//--/g" -e "s/^/   /" -e "s/-/|/"'
