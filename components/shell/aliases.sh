#!/bin/sh
# components/shell/aliases.sh - Core shell aliases

# =============================================================================
# Directory Listing
# =============================================================================

case "$CURRENT_PLATFORM" in
    darwin)
        alias ls='ls -G'
        alias ll='ls -lhG'
        alias la='ls -lahG'
        alias l='ls -CF'
        ;;
    linux)
        alias ls='ls --color=auto'
        alias ll='ls -lh --color=auto'
        alias la='ls -lah --color=auto'
        alias l='ls -CF --color=auto'
        ;;
esac

# =============================================================================
# Navigation
# =============================================================================

alias ..='cd ..'
alias ...='cd ../..'
alias ....='cd ../../..'
alias .....='cd ../../../..'
alias ~='cd ~'
alias -- -='cd -'

# =============================================================================
# Safety
# =============================================================================

alias rm='rm -i'
alias cp='cp -i'
alias mv='mv -i'
alias mkdir='mkdir -p'

# =============================================================================
# Grep
# =============================================================================

alias grep='grep --color=auto'
alias fgrep='fgrep --color=auto'
alias egrep='egrep --color=auto'

# =============================================================================
# Disk Usage
# =============================================================================

alias df='df -h'
alias du='du -h'
alias dud='du -d 1 -h'
alias duf='du -sh *'

# =============================================================================
# Process
# =============================================================================

alias psg='ps aux | grep -v grep | grep -i -e VSZ -e'
alias top='top -o cpu'

# =============================================================================
# Date/Time
# =============================================================================

alias now='date +"%Y-%m-%d %H:%M:%S"'
alias nowdate='date +"%Y-%m-%d"'
alias nowtime='date +"%H:%M:%S"'
alias week='date +%V'

# =============================================================================
# Misc
# =============================================================================

alias c='clear'
# h is a function in functions.sh (history | grep)
alias j='jobs -l'
alias path='echo -e ${PATH//:/\\n}'
alias reload='source ~/.bashrc 2>/dev/null || source ~/.zshrc'
