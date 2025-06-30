# Shell-agnostic color definitions
# Use appropriate prompt escaping based on current shell
if [ "$CURRENT_SHELL" = "zsh" ]; then
    # Zsh uses %{...%} for non-printing characters
    BLACK="%{\033[0;30m%}"
    BLACKB="%{\033[1;30m%}"
    RED="%{\033[0;31m%}"
    REDB="%{\033[1;31m%}"
    GREEN="%{\033[0;32m%}"
    GREENB="%{\033[1;32m%}"
    YELLOW="%{\033[0;33m%}"
    YELLOWB="%{\033[1;33m%}"
    BLUE="%{\033[0;34m%}"
    BLUEB="%{\033[1;34m%}"
    PURPLE="%{\033[0;35m%}"
    PURPLEB="%{\033[1;35m%}"
    CYAN="%{\033[0;36m%}"
    CYANB="%{\033[1;36m%}"
    WHITE="%{\033[0;37m%}"
    WHITEB="%{\033[1;37m%}"
    RESET="%{\033[0;0m%}"
else
    # Bash uses \001...\002 for non-printing characters
    BLACK="\001\033[0;30m\002"
    BLACKB="\001\033[1;30m\002"
    RED="\001\033[0;31m\002"
    REDB="\001\033[1;31m\002"
    GREEN="\001\033[0;32m\002"
    GREENB="\001\033[1;32m\002"
    YELLOW="\001\033[0;33m\002"
    YELLOWB="\001\033[1;33m\002"
    BLUE="\001\033[0;34m\002"
    BLUEB="\001\033[1;34m\002"
    PURPLE="\001\033[0;35m\002"
    PURPLEB="\001\033[1;35m\002"
    CYAN="\001\033[0;36m\002"
    CYANB="\001\033[1;36m\002"
    WHITE="\001\033[0;37m\002"
    WHITEB="\001\033[1;37m\002"
    RESET="\001\033[0;0m\002"
fi

# Shell-portable git branch detection
git_branch() {
    local git_status branch_line
    git_status="$(git status 2>/dev/null)" || return 0
    
    # Extract branch or commit info using portable methods
    if echo "$git_status" | grep -q "On branch"; then
        branch_line=$(echo "$git_status" | grep "On branch" | head -1)
        branch_line=${branch_line#On branch }
        echo "[$branch_line]"
    elif echo "$git_status" | grep -q "HEAD detached at"; then
        branch_line=$(echo "$git_status" | grep "HEAD detached at" | head -1)
        branch_line=${branch_line#HEAD detached at }
        echo "[$branch_line]"
    fi
}

# Shell-portable git status color
git_color() {
    local git_status
    git_status="$(git status 2>/dev/null)" || { echo -e "$WHITE"; return 0; }
    
    # Use portable string matching
    case "$git_status" in
        *"nothing added to commit but untracked files present"*)
            echo -e "$YELLOW" ;;
        *"no changes added to commit"*)
            echo -e "$YELLOW" ;;
        *"Your branch is ahead of"*)
            echo -e "$RED" ;;
        *"Changes to be committed"*)
            echo -e "$RED" ;;
        *"nothing to commit, working tree clean"*)
            echo -e "$GREEN" ;;
        *)
            echo -e "$WHITE" ;;
    esac
}

# Enable colors in ls and other tools
export CLICOLOR=1
export LSCOLORS=fxgxexcxbxegxgxbxbxfxf

# Shell-specific prompt configuration
if [ "$CURRENT_SHELL" = "zsh" ]; then
    # Zsh prompt configuration
    setopt PROMPT_SUBST
    PROMPT="${WHITE}%D{%Y-%m-%d %H:%M:%S}${WHITE} → ${CYAN}%n${WHITE} on ${BLUE}%m${WHITE} → ${PURPLE}[%~]"'$(git_color)$(git_branch)'"${WHITE}"$'\n$ '"${RESET}"
else
    # Bash prompt configuration
    PS1="${WHITE}\D{%F %T}${WHITE} → ${CYAN}\u${WHITE} on ${BLUE}\h${WHITE} → ${PURPLE}[\w]"'$(git_color)$(git_branch)'"${WHITE}"$'\n\$ '"${RESET}"
    export PS1
fi
