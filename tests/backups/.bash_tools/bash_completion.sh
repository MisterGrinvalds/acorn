# Shell-specific completion loading
if [ "$CURRENT_SHELL" = "bash" ]; then
    # Load bash completion
    if [ -r "/opt/homebrew/etc/profile.d/bash_completion.sh" ]; then
        . "/opt/homebrew/etc/profile.d/bash_completion.sh"
    elif [ -r "/usr/share/bash-completion/bash_completion" ]; then
        . "/usr/share/bash-completion/bash_completion"
    fi
elif [ "$CURRENT_SHELL" = "zsh" ]; then
    # Enable zsh completion system
    autoload -Uz compinit
    compinit
    
    # Load additional zsh completions if available
    if [ -d "/opt/homebrew/share/zsh/site-functions" ]; then
        fpath=("/opt/homebrew/share/zsh/site-functions" $fpath)
    fi
fi
