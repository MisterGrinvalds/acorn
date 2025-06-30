#!/bin/bash
# Terminal configuration with git-aware prompt and Solarized colors


# Get current git branch
git_branch() {
    local git_status="$(git status 2> /dev/null)"
    local on_branch="On branch ([^${IFS}]*)"
    local on_commit="HEAD detached at ([^${IFS}]*)"

    if [[ $git_status =~ $on_branch ]]; then
        local branch=${BASH_REMATCH[1]}
        echo "($branch)"
    elif [[ $git_status =~ $on_commit ]]; then
        local commit=${BASH_REMATCH[1]}
        echo "($commit)"
    fi
}
# Git status colors
git_color() {
    local git_status="$(git status 2> /dev/null)"

    if [[ ! $git_status =~ "working directory clean" ]] && [[ ! $git_status =~ "working tree 
clean" ]]; then
        echo -e "\e[1;31m"  # red
    elif [[ $git_status =~ "Your branch is ahead of" ]]; then
        echo -e "\e[1;33m"  # yellow
    elif [[ $git_status =~ "nothing to commit" ]]; then
        echo -e "\e[1;32m"  # green
    else
        echo -e "\e[1;35m"  # purple
    fi
}


# Set prompt with git integration
PS1="\[\e[1;37m\]\n\[\e[1;37m\]\A \[\e[1;32m\]\u\[\e[1;37m\] on \[\e[1;33m\]\h\[\e[1;37m\] 
\[\e[1;34m\][\w]\$(git_color)\$(git_branch)\[\e[0m\]\n\$ "

# Enable color support
if [[ "$OSTYPE" == "darwin"* ]]; then
    export CLICOLOR=1
    export LSCOLORS=ExFxBxDxCxegedabagacad
fi