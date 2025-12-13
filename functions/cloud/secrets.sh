#!/bin/sh
# Secrets Management Integration
# Requires: DOTFILES_ROOT environment variable

# Load secrets into current environment
load_secrets() {
    local secrets_file="$HOME/.automation/secrets/.env"

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
        echo "No secrets file found. Run 'auto secrets setup' first."
        return 1
    fi
}

# Quick check of secrets status
secrets_status() {
    if command -v auto >/dev/null 2>&1; then
        auto secrets check-requirements
    else
        echo "Automation framework not found"
        return 1
    fi
}

# Quick secrets validation
validate_secrets() {
    if command -v auto >/dev/null 2>&1; then
        auto secrets validate
    else
        echo "Automation framework not found"
        return 1
    fi
}

# Show configured secrets (keys only, not values)
list_secrets() {
    if command -v auto >/dev/null 2>&1; then
        auto secrets list
    else
        echo "Automation framework not found"
        return 1
    fi
}

# Check if specific API key is available
check_aws_key() {
    if [ -n "$AWS_ACCESS_KEY_ID" ] && [ -n "$AWS_SECRET_ACCESS_KEY" ]; then
        echo "AWS credentials available"
        return 0
    else
        echo "AWS credentials not found"
        return 1
    fi
}

check_azure_key() {
    if [ -n "$AZURE_CLIENT_ID" ] && [ -n "$AZURE_CLIENT_SECRET" ] && [ -n "$AZURE_TENANT_ID" ]; then
        echo "Azure credentials available"
        return 0
    else
        echo "Azure credentials not found"
        return 1
    fi
}

check_github_key() {
    if [ -n "$GITHUB_TOKEN" ]; then
        echo "GitHub token available"
        return 0
    else
        echo "GitHub token not found"
        return 1
    fi
}

check_digitalocean_key() {
    if [ -n "$DIGITALOCEAN_TOKEN" ]; then
        echo "DigitalOcean token available"
        return 0
    else
        echo "DigitalOcean token not found"
        return 1
    fi
}

# Auto-load secrets if environment variable is set
if [ "$AUTO_LOAD_SECRETS" = "true" ]; then
    load_secrets >/dev/null 2>&1
fi

# Aliases for quick access
alias secrets-load='load_secrets'
alias secrets-status='secrets_status'
alias secrets-validate='validate_secrets'
alias secrets-list='list_secrets'
alias check-aws='check_aws_key'
alias check-azure='check_azure_key'
alias check-github='check_github_key'
alias check-do='check_digitalocean_key'
