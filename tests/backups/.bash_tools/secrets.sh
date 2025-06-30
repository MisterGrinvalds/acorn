#!/bin/bash
# Secrets Management Integration for Bash Profile
# Quick access to secrets management from any shell

# Load secrets into current environment
load_secrets() {
    local secrets_file="$HOME/.automation/secrets/.env"
    
    if [ -f "$secrets_file" ]; then
        # Check if file is readable
        if [ -r "$secrets_file" ]; then
            set -a  # Automatically export variables
            source "$secrets_file"
            set +a  # Stop automatically exporting
            echo "✅ Secrets loaded into environment"
        else
            echo "❌ Cannot read secrets file (check permissions)"
            return 1
        fi
    else
        echo "❌ No secrets file found. Run 'auto secrets setup' first."
        return 1
    fi
}

# Quick check of secrets status
secrets_status() {
    if command -v auto >/dev/null 2>&1; then
        auto secrets check-requirements
    else
        echo "❌ Automation framework not found"
        return 1
    fi
}

# Quick secrets validation
validate_secrets() {
    if command -v auto >/dev/null 2>&1; then
        auto secrets validate
    else
        echo "❌ Automation framework not found"
        return 1
    fi
}

# Show configured secrets (keys only, not values)
list_secrets() {
    if command -v auto >/dev/null 2>&1; then
        auto secrets list
    else
        echo "❌ Automation framework not found"
        return 1
    fi
}

# Quick setup functions for individual providers
setup_aws() {
    if command -v auto >/dev/null 2>&1; then
        auto secrets aws
    else
        echo "❌ Automation framework not found"
        return 1
    fi
}

setup_azure() {
    if command -v auto >/dev/null 2>&1; then
        auto secrets azure
    else
        echo "❌ Automation framework not found"
        return 1
    fi
}

setup_github() {
    if command -v auto >/dev/null 2>&1; then
        auto secrets github
    else
        echo "❌ Automation framework not found"
        return 1
    fi
}

setup_digitalocean() {
    if command -v auto >/dev/null 2>&1; then
        auto secrets digitalocean
    else
        echo "❌ Automation framework not found"
        return 1
    fi
}

# Check if specific API key is available
check_aws_key() {
    if [ -n "$AWS_ACCESS_KEY_ID" ] && [ -n "$AWS_SECRET_ACCESS_KEY" ]; then
        echo "✅ AWS credentials available"
        return 0
    else
        echo "❌ AWS credentials not found"
        return 1
    fi
}

check_azure_key() {
    if [ -n "$AZURE_CLIENT_ID" ] && [ -n "$AZURE_CLIENT_SECRET" ] && [ -n "$AZURE_TENANT_ID" ]; then
        echo "✅ Azure credentials available"
        return 0
    else
        echo "❌ Azure credentials not found"
        return 1
    fi
}

check_github_key() {
    if [ -n "$GITHUB_TOKEN" ]; then
        echo "✅ GitHub token available"
        return 0
    else
        echo "❌ GitHub token not found"
        return 1
    fi
}

check_digitalocean_key() {
    if [ -n "$DIGITALOCEAN_TOKEN" ]; then
        echo "✅ DigitalOcean token available"
        return 0
    else
        echo "❌ DigitalOcean token not found"
        return 1
    fi
}

# Auto-load secrets if they exist and environment variable is set
if [ "$AUTO_LOAD_SECRETS" = "true" ]; then
    load_secrets >/dev/null 2>&1
fi

# Aliases for quick access
alias secrets-load='load_secrets'
alias secrets-status='secrets_status'
alias secrets-validate='validate_secrets'
alias secrets-list='list_secrets'
alias setup-aws='setup_aws'
alias setup-azure='setup_azure'
alias setup-github='setup_github'
alias setup-do='setup_digitalocean'
alias check-aws='check_aws_key'
alias check-azure='check_azure_key'
alias check-github='check_github_key'
alias check-do='check_digitalocean_key'

# Help function
secrets_help() {
    cat << 'EOF'
Secrets Management Quick Commands

LOADING SECRETS:
    load_secrets       Load secrets into current environment
    secrets-load       Alias for load_secrets

STATUS & VALIDATION:
    secrets_status     Check which API keys are configured
    secrets-status     Alias for secrets_status
    validate_secrets   Validate all configured API keys
    secrets-validate   Alias for validate_secrets
    list_secrets       List configured secret keys (not values)
    secrets-list       Alias for list_secrets

SETUP COMMANDS:
    setup_aws         Setup AWS credentials
    setup-aws         Alias for setup_aws
    setup_azure       Setup Azure credentials
    setup-azure       Alias for setup_azure
    setup_github      Setup GitHub credentials
    setup-github      Alias for setup_github
    setup_digitalocean Setup DigitalOcean credentials
    setup-do          Alias for setup_digitalocean

KEY CHECKING:
    check_aws_key     Check if AWS credentials are loaded
    check-aws         Alias for check_aws_key
    check_azure_key   Check if Azure credentials are loaded
    check-azure       Alias for check_azure_key
    check_github_key  Check if GitHub token is loaded
    check-github      Alias for check_github_key
    check_digitalocean_key Check if DigitalOcean token is loaded
    check-do          Alias for check_digitalocean_key

AUTOMATION COMMANDS:
    auto secrets setup           Interactive secrets wizard
    auto secrets check-requirements Check missing API keys
    auto secrets validate        Validate all credentials
    auto secrets aws            Setup AWS credentials
    auto secrets azure          Setup Azure credentials
    auto secrets github         Setup GitHub credentials
    auto secrets digitalocean   Setup DigitalOcean credentials

ENVIRONMENT VARIABLES:
    AUTO_LOAD_SECRETS=true      Auto-load secrets on shell startup

EXAMPLES:
    load_secrets                Load all secrets into environment
    secrets-status              Quick status check
    setup-aws                   Setup AWS credentials
    check-github                Check if GitHub token is available
    auto secrets setup          Full interactive setup
EOF
}

alias secrets-help='secrets_help'