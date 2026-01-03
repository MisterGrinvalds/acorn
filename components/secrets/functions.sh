#!/bin/sh
# components/secrets/functions.sh - Secrets management functions

# =============================================================================
# Core Secret Functions
# =============================================================================

# Load secrets into current environment
load_secrets() {
    local secrets_file="${SECRETS_DIR:-${XDG_DATA_HOME:-$HOME/.local/share}/secrets}/.env"

    if [ -f "$secrets_file" ]; then
        if [ -r "$secrets_file" ]; then
            set -a
            . "$secrets_file"
            set +a
            echo "Secrets loaded into environment"
        else
            echo "Cannot read secrets file (check permissions)"
            return 1
        fi
    else
        echo "No secrets file found at: $secrets_file"
        return 1
    fi
}

# Quick check of secrets status
secrets_status() {
    local secrets_file="${SECRETS_DIR:-${XDG_DATA_HOME:-$HOME/.local/share}/secrets}/.env"
    if [ -f "$secrets_file" ]; then
        echo "Secrets file: $secrets_file"
        echo "Keys defined: $(grep -c '^[A-Z]' "$secrets_file" 2>/dev/null || echo 0)"
    else
        echo "No secrets file found"
    fi
}

# Quick secrets validation
validate_secrets() {
    echo "Checking common credentials..."
    check_aws_key
    check_azure_key
    check_github_key
    check_digitalocean_key
}

# Show configured secrets (keys only, not values)
list_secrets() {
    local secrets_file="${SECRETS_DIR:-${XDG_DATA_HOME:-$HOME/.local/share}/secrets}/.env"
    if [ -f "$secrets_file" ]; then
        echo "Configured secrets:"
        grep '^[A-Z]' "$secrets_file" 2>/dev/null | cut -d'=' -f1 | sort
    else
        echo "No secrets file found"
    fi
}

# =============================================================================
# Cloud Credential Checks
# =============================================================================

check_aws_key() {
    if [ -n "$AWS_ACCESS_KEY_ID" ] && [ -n "$AWS_SECRET_ACCESS_KEY" ]; then
        echo "AWS credentials: available"
        return 0
    else
        echo "AWS credentials: not found"
        return 1
    fi
}

check_azure_key() {
    if [ -n "$AZURE_CLIENT_ID" ] && [ -n "$AZURE_CLIENT_SECRET" ] && [ -n "$AZURE_TENANT_ID" ]; then
        echo "Azure credentials: available"
        return 0
    else
        echo "Azure credentials: not found"
        return 1
    fi
}

check_github_key() {
    if [ -n "$GITHUB_TOKEN" ]; then
        echo "GitHub token: available"
        return 0
    else
        echo "GitHub token: not found"
        return 1
    fi
}

check_digitalocean_key() {
    if [ -n "$DIGITALOCEAN_TOKEN" ]; then
        echo "DigitalOcean token: available"
        return 0
    else
        echo "DigitalOcean token: not found"
        return 1
    fi
}

check_openai_key() {
    if [ -n "$OPENAI_API_KEY" ]; then
        echo "OpenAI API key: available"
        return 0
    else
        echo "OpenAI API key: not found"
        return 1
    fi
}

check_anthropic_key() {
    if [ -n "$ANTHROPIC_API_KEY" ]; then
        echo "Anthropic API key: available"
        return 0
    else
        echo "Anthropic API key: not found"
        return 1
    fi
}

# =============================================================================
# All Credentials Check
# =============================================================================

check_all_keys() {
    echo "Credential Status"
    echo "================="
    check_aws_key
    check_azure_key
    check_github_key
    check_digitalocean_key
    check_openai_key
    check_anthropic_key
}
