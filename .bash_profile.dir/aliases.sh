#!/bin/bash
# bash
alias bash='bash && source ~/.bash_profile'

# git
alias g='git'
alias gs='git status'                                                           # Used as often as `ls`

# navigation
alias ..='cd ..'
alias ...='cd ../../'
alias ....='cd ../../../'
alias getsshkey="pbcopy < ~/.ssh/id_rsa.pub"                                    # Copy SSH key to the keyboard
alias ll='ls -alh'                                                              # List files
alias llr='ls -alhr'                                                            # List files (reverse)
alias lls='ls -alhS'                                                            # List files by size
alias llsr='ls -alhSr'                                                          # List files by size (reverse)
alias lld='ls -alht'                                                            # List files by date
alias lldr='ls -alhtr'                                                          # List files by date (reverse)
alias lldc='ls -alhtU'                                                          # List files by date created
alias lldcr='ls -alhtUr'                                                        # List files by date created (reverse)
alias mkdir='mkdir -pv'                                                         # Make parent directories if needed
alias perm="stat -f '%Lp'"                                                      # View the permissions of a file/dir as a number
alias resource='source ~/.bash_profile'
alias tree="\
    ls -R |\
    grep ":$" | \
    sed -e 's/:$//' -e 's/[^-][^\/]*\//--/g' -e 's/^/   /' -e 's/-/|/'\
    "                                                                           # List the file structure of the current directory

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

# tmux enhanced
alias ta='tmux attach-session -t'
alias tn='tmux new-session -s'
alias tk='tmux kill-session -t'
alias tko='tmux kill-session -a'         # Kill all other sessions
alias ti='tmux info'
alias ts='tmux list-sessions'
alias tks='tmux kill-server'             # Kill tmux server
alias td='tmux detach'                   # Detach current session

# tmux project sessions
alias twork='tmux new-session -s work -d'
alias tdev='tmux new-session -s dev -d'
alias tk8s='tmux new-session -s k8s -d'

# Xclip
alias c="xclip"                                                                 # copy to X clipboard (register *)
alias cs="xclip -selection clipboard"                                           # copy to system wide clipboard (register +)
alias v="xclip -o"                                                              # output copied content (paste)
alias vs="xclip -o -selection clipboard"                                        # paste from system wide clipboard (equivalent to `v -selection clipboard`)
