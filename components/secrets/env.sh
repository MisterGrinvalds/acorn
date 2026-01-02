#!/bin/sh
# components/secrets/env.sh - Secrets environment variables

# Default secrets location
export SECRETS_DIR="${SECRETS_DIR:-$HOME/.automation/secrets}"

# Auto-load secrets if enabled
if [ "$AUTO_LOAD_SECRETS" = "true" ]; then
    if [ -f "$SECRETS_DIR/.env" ]; then
        set -a
        . "$SECRETS_DIR/.env"
        set +a
    fi
fi
