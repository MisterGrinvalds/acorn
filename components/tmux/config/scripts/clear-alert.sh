#!/bin/bash
# clear-alert.sh - Clear window alert when switching to it
# Called by tmux hook: after-select-window

WINDOW_ID="$1"

# Check if this window has an alert flag
ALERT_FLAG=$(tmux show-window-options -t "$WINDOW_ID" -v @alert 2>/dev/null)

if [ "$ALERT_FLAG" = "1" ]; then
    # Clear the alert flag
    tmux set-window-option -t "$WINDOW_ID" @alert 0

    # Reset to default Catppuccin Mocha style
    # Matches window-status-format in tmux.conf
    tmux set-window-option -t "$WINDOW_ID" window-status-style "fg=#bac2de,bg=#45475a"
fi
