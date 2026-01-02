#!/bin/sh
# components/git/aliases.sh - Git aliases

# =============================================================================
# Core Git Commands
# =============================================================================

alias g='git'
alias gs='git status'
alias ga='git add'
alias gaa='git add --all'
alias gc='git commit'
alias gcm='git commit -m'
alias gca='git commit --amend'
alias gco='git checkout'
alias gb='git branch'
alias gba='git branch -a'
alias gbd='git branch -d'
alias gbD='git branch -D'
alias gd='git diff'
alias gds='git diff --staged'
alias gp='git push'
alias gpf='git push --force-with-lease'
alias gpl='git pull'
alias gplr='git pull --rebase'
alias gf='git fetch'
alias gfa='git fetch --all'
alias gm='git merge'
alias gr='git rebase'
alias gri='git rebase -i'
alias grc='git rebase --continue'
alias gra='git rebase --abort'
alias gst='git stash'
alias gstp='git stash pop'
alias gstl='git stash list'

# =============================================================================
# Git Log
# =============================================================================

alias gl='git log --oneline -20'
alias glog='git log --oneline --graph --decorate'
alias glg='git log --graph --pretty=format:"%Cred%h%Creset -%C(yellow)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset" --abbrev-commit'
alias gla='git log --oneline --all --graph --decorate'

# =============================================================================
# Git Remote
# =============================================================================

alias gra='git remote add'
alias grv='git remote -v'
alias gru='git remote update'

# =============================================================================
# Git Reset
# =============================================================================

alias grh='git reset HEAD'
alias grhh='git reset HEAD --hard'
alias grhs='git reset HEAD --soft'

# =============================================================================
# Git Clean
# =============================================================================

alias gclean='git clean -fd'
alias gpristine='git reset --hard && git clean -fdx'
