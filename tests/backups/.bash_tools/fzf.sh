#!/bin/sh
# FZF configuration - shell portable

# Environment variables
export FZF_DEFAULT_OPS="--extended"
export FZF_BIN_PATH="$HOMEBREW_BIN_PATH/fzf"
export FZF_CELLAR_PATH="$HOMEBREW_CELLAR_PATH/fzf"

# Get FZF version in a portable way
if command -v fzf >/dev/null 2>&1; then
    FZF_VERSION=$(fzf --version | cut -d' ' -f1)
    export FZF_VERSION
fi

# Shell-specific FZF integration
if [ "$CURRENT_SHELL" = "bash" ]; then
    # Bash FZF integration
    if [ -f "$FZF_CELLAR_PATH/$FZF_VERSION/shell/completion.bash" ]; then
        [[ $- == *i* ]] && source "$FZF_CELLAR_PATH/$FZF_VERSION/shell/completion.bash" 2>/dev/null
    fi
    if [ -f "$FZF_CELLAR_PATH/$FZF_VERSION/shell/key-bindings.bash" ]; then
        source "$FZF_CELLAR_PATH/$FZF_VERSION/shell/key-bindings.bash"
    fi
elif [ "$CURRENT_SHELL" = "zsh" ]; then
    # Zsh FZF integration
    if [ -f "$FZF_CELLAR_PATH/$FZF_VERSION/shell/completion.zsh" ]; then
        source "$FZF_CELLAR_PATH/$FZF_VERSION/shell/completion.zsh" 2>/dev/null
    fi
    if [ -f "$FZF_CELLAR_PATH/$FZF_VERSION/shell/key-bindings.zsh" ]; then
        source "$FZF_CELLAR_PATH/$FZF_VERSION/shell/key-bindings.zsh"
    fi
fi

# OS Dependent FZF command configuration
if [ "$OSTYPE" = "darwin"* ]; then
    export FZF_DEFAULT_COMMAND="fd --type f"
elif [ "$OSTYPE" = "linux-gnu" ]; then
    export FZF_DEFAULT_COMMAND="fdfind --type f"
fi

# Set FZF_CTRL_T_COMMAND after FZF_DEFAULT_COMMAND is defined
export FZF_CTRL_T_COMMAND="$FZF_DEFAULT_COMMAND"

