#!/bin/sh
# components/iterm2/env.sh - iTerm2 environment variables

# iTerm2 shell integration path
export ITERM2_SHELL_INTEGRATION="${HOME}/.iterm2_shell_integration.${CURRENT_SHELL:-bash}"

# iTerm2 utilities directory
export ITERM2_UTILITIES_DIR="${HOME}/.iterm2"

# Source iTerm2 shell integration if available and running in iTerm2
if [ -n "$ITERM_SESSION_ID" ] && [ -f "$ITERM2_SHELL_INTEGRATION" ]; then
    # shellcheck disable=SC1090
    . "$ITERM2_SHELL_INTEGRATION"
fi
