#!/bin/sh
# components/secrets/env.sh - Secrets environment variables

# Default secrets location (XDG-compliant)
export SECRETS_DIR="${SECRETS_DIR:-${XDG_DATA_HOME:-$HOME/.local/share}/secrets}"

# Auto-load secrets if enabled
if [ "$AUTO_LOAD_SECRETS" = "true" ]; then
    if [ -f "$SECRETS_DIR/.env" ]; then
        set -a
        . "$SECRETS_DIR/.env"
        set +a
    fi
fi
