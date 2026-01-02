#!/bin/sh
# components/iterm2/aliases.sh - iTerm2 aliases

# Only define aliases on macOS
if [ "$CURRENT_PLATFORM" = "darwin" ]; then
    # Open iTerm2 preferences
    alias iterm-prefs='open -a iTerm "iterm2://preferences"'

    # Quick profile switching (uses iTerm2 escape sequences)
    alias iterm-default='iterm_profile "Default"'
    alias iterm-dark='iterm_profile "Dotfiles Dark"'

    # Shell integration
    alias iterm-integration='iterm_install_integration'
fi
