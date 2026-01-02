#!/bin/sh
# components/fzf/completions.sh - FZF shell integrations

# Load FZF key bindings and completions
if [ -n "$FZF_LOCATION" ] && [ -d "$FZF_LOCATION" ]; then
    case "$CURRENT_SHELL" in
        bash)
            [ -f "$FZF_LOCATION/shell/completion.bash" ] && . "$FZF_LOCATION/shell/completion.bash"
            [ -f "$FZF_LOCATION/shell/key-bindings.bash" ] && . "$FZF_LOCATION/shell/key-bindings.bash"
            ;;
        zsh)
            [ -f "$FZF_LOCATION/shell/completion.zsh" ] && . "$FZF_LOCATION/shell/completion.zsh"
            [ -f "$FZF_LOCATION/shell/key-bindings.zsh" ] && . "$FZF_LOCATION/shell/key-bindings.zsh"
            ;;
    esac
fi
