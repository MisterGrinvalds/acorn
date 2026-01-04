#!/bin/sh
# components/cloudflare/env.sh - CloudFlare CLI environment variables
#
# This file is sourced for ALL shells (interactive and non-interactive).
# Configures wrangler to use XDG-compliant directories.

# Configure wrangler to use XDG directories
# Wrangler respects WRANGLER_HOME for configuration storage
export WRANGLER_HOME="${XDG_CONFIG_HOME}/wrangler"

# Ensure wrangler config directory exists
if [ ! -d "${WRANGLER_HOME}" ]; then
    mkdir -p "${WRANGLER_HOME}" 2>/dev/null
fi
