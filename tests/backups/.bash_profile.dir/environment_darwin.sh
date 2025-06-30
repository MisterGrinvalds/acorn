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
export HOMEBREW_BIN_PATH="/opt/homebrew/bin"
export HOMEBREW_CELLAR_PATH="/opt/homebrew/Cellar"
export LESSHISTFILE="$DOTFILES/.lesshst"
export ENVS_LOCATION="$HOME/envs"
export PATH="$HOMEBREW_BIN_PATH:$PATH"
export WGETRC="$DOTFILES/.wgetrc"
export XDG_CACHE_HOME="$HOME/Library/Caches"
export XDG_DATA_HOME="$HOME/Library/Application Support"
if [ -f "/Library/Frameworks/R.framework/Resources/etc/Rprofile.site" ] && [ ! -L "/Library/Frameworks/R.framework/Resources/etc/Rprofile.site" ]; then
	ln -s "$DOTFILES/.R/Rprofile.site" "/Library/Frameworks/R.framework/Resources/etc/"
fi


