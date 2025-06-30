export DISPLAY=:0
export FZF_LOCATION="/home/linuxbrew/.linuxbrew/opt/fzf"

# Shell-specific history configuration
if [ "$CURRENT_SHELL" = "zsh" ]; then
    # Zsh history settings
    export HISTFILE="$HOME/.zsh_history"
    export HISTSIZE=3000
    export SAVEHIST=3000
    setopt HIST_IGNORE_DUPS
    setopt HIST_IGNORE_SPACE
    setopt HIST_EXPIRE_DUPS_FIRST
elif [ "$CURRENT_SHELL" = "bash" ]; then
    # Bash history settings
    if [[ ! "$HISTCONTROL" == *erasedups:ignoredups:ignorespace* ]]; then
        export HISTCONTROL=$HISTCONTROL:erasedups:ignoredups:ignorespace
    fi
    export HISTFILE="$HOME/.bash_history"
    export HISTFILESIZE=3000               
    export HISTSIZE=3000
fi                   
export LESSHISTFILE="$DOTFILES/.lesshst"
export WGETRC="$DOTFILES/.wgetrc"
export XDG_CACHE_HOME="$HOME/.cache"
export XDG_DATA_HOME="$HOME/.local/share"
