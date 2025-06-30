#!/bin/bash
# Secrets Management Module
# Secure handling of API keys, tokens, and sensitive configuration

# Load framework
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "$SCRIPT_DIR/framework/core.sh"

# Secrets directory
SECRETS_DIR="$AUTO_HOME/secrets"
SECRETS_CONFIG="$SECRETS_DIR/config"
ENCRYPTED_SECRETS="$SECRETS_DIR/vault.enc"
SECRETS_TEMPLATE="$SECRETS_DIR/template.env"

secrets_help() {
    cat << 'EOF'
Secrets Management

USAGE: auto secrets <command> [options]

COMMANDS:
    init                     Initialize secrets management
    setup                    Interactive setup of API keys
    validate                 Validate all configured secrets
    list                     List configured secret keys (not values)
    encrypt                  Encrypt secrets vault
    decrypt                  Decrypt secrets vault
    rotate <key>             Rotate a specific API key
    export                   Export secrets for current session
    backup                   Backup encrypted secrets
    restore <backup>         Restore from backup
    template                 Generate secrets template
    check-requirements       Check which API keys are missing

CLOUD PROVIDERS:
    aws                      Setup AWS credentials
    azure                    Setup Azure credentials  
    digitalocean             Setup DigitalOcean credentials
    github                   Setup GitHub credentials

OPTIONS:
    --force                  Force overwrite existing secrets
    --export-format <fmt>    Export format (env, json, yaml)
    --backup-dir <dir>       Custom backup directory

EXAMPLES:
    auto secrets init                    # Initialize secrets management
    auto secrets setup                   # Interactive setup wizard
    auto secrets validate                # Check all API keys
    auto secrets aws                     # Setup AWS credentials
    auto secrets export --format env     # Export as environment variables
    auto secrets rotate github-token    # Rotate GitHub token
EOF
}

# Initialize secrets management
secrets_init() {
    log_info "Initializing secrets management..."
    
    mkdir -p "$SECRETS_DIR"
    chmod 700 "$SECRETS_DIR"
    
    # Create secrets config
    if [ ! -f "$SECRETS_CONFIG" ]; then
        cat > "$SECRETS_CONFIG" << 'EOF'
# Secrets Management Configuration
SECRETS_ENCRYPTION_METHOD=openssl
SECRETS_BACKUP_RETENTION=30
SECRETS_AUTO_ROTATE=false
SECRETS_VALIDATION_INTERVAL=7d
EOF
        chmod 600 "$SECRETS_CONFIG"
        log_info "Created secrets configuration"
    fi
    
    # Create secrets template
    secrets_create_template
    
    # Create .gitignore for secrets
    cat > "$SECRETS_DIR/.gitignore" << 'EOF'
# Never commit secrets
*.env
*.key
*.pem
*.p12
*.json
vault.enc
config
!template.env
!.gitignore
EOF
    
    log_success "Secrets management initialized"
}

# Create secrets template
secrets_create_template() {
    cat > "$SECRETS_TEMPLATE" << 'EOF'
# Secrets Template - Copy to .env and fill in your values
# DO NOT commit this file with actual secrets!

# =============================================================================
# CLOUD PROVIDERS
# =============================================================================

# AWS Credentials
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
AWS_DEFAULT_REGION=us-east-1
AWS_SESSION_TOKEN=
AWS_PROFILE=default

# Azure Credentials  
AZURE_CLIENT_ID=
AZURE_CLIENT_SECRET=
AZURE_TENANT_ID=
AZURE_SUBSCRIPTION_ID=

# DigitalOcean
DIGITALOCEAN_TOKEN=
DIGITALOCEAN_SPACES_ACCESS_KEY=
DIGITALOCEAN_SPACES_SECRET_KEY=

# =============================================================================
# SOURCE CODE MANAGEMENT
# =============================================================================

# GitHub
GITHUB_TOKEN=
GITHUB_USERNAME=
GITHUB_ORGANIZATION=

# GitLab
GITLAB_TOKEN=
GITLAB_USERNAME=

# =============================================================================
# CONTAINER REGISTRIES
# =============================================================================

# Docker Hub
DOCKER_USERNAME=
DOCKER_PASSWORD=
DOCKER_EMAIL=

# GitHub Container Registry
GHCR_TOKEN=

# Azure Container Registry
ACR_USERNAME=
ACR_PASSWORD=

# =============================================================================
# KUBERNETES
# =============================================================================

# Kubernetes Service Account Tokens
K8S_SERVICE_ACCOUNT_TOKEN=
K8S_CLUSTER_CA_CERTIFICATE=
K8S_CLUSTER_ENDPOINT=

# Helm Repository Credentials
HELM_REPO_USERNAME=
HELM_REPO_PASSWORD=

# =============================================================================
# DATABASES
# =============================================================================

# PostgreSQL
POSTGRES_USERNAME=
POSTGRES_PASSWORD=
POSTGRES_HOST=
POSTGRES_PORT=5432
POSTGRES_DATABASE=

# MySQL
MYSQL_USERNAME=
MYSQL_PASSWORD=
MYSQL_HOST=
MYSQL_PORT=3306
MYSQL_DATABASE=

# MongoDB
MONGODB_URI=
MONGODB_USERNAME=
MONGODB_PASSWORD=

# Redis
REDIS_URL=
REDIS_PASSWORD=

# =============================================================================
# MONITORING & OBSERVABILITY
# =============================================================================

# DataDog
DATADOG_API_KEY=
DATADOG_APP_KEY=

# New Relic
NEW_RELIC_LICENSE_KEY=
NEW_RELIC_API_KEY=

# Prometheus/Grafana
GRAFANA_API_TOKEN=
PROMETHEUS_USERNAME=
PROMETHEUS_PASSWORD=

# =============================================================================
# COMMUNICATION
# =============================================================================

# Slack
SLACK_BOT_TOKEN=
SLACK_APP_TOKEN=
SLACK_WEBHOOK_URL=

# Discord
DISCORD_BOT_TOKEN=
DISCORD_WEBHOOK_URL=

# =============================================================================
# CI/CD
# =============================================================================

# Jenkins
JENKINS_URL=
JENKINS_USERNAME=
JENKINS_API_TOKEN=

# CircleCI
CIRCLECI_TOKEN=

# Travis CI
TRAVIS_TOKEN=

# =============================================================================
# DEVELOPMENT TOOLS
# =============================================================================

# Sentry
SENTRY_DSN=
SENTRY_AUTH_TOKEN=

# OpenAI
OPENAI_API_KEY=

# Anthropic
ANTHROPIC_API_KEY=
EOF
    
    chmod 600 "$SECRETS_TEMPLATE"
    log_info "Created secrets template: $SECRETS_TEMPLATE"
}

# Interactive setup wizard
secrets_setup() {
    log_info "Starting interactive secrets setup..."
    
    local secrets_file="$SECRETS_DIR/.env"
    local force_overwrite=false
    
    # Check for --force flag
    if [[ "$*" =~ --force ]]; then
        force_overwrite=true
    fi
    
    # Check if secrets file exists
    if [ -f "$secrets_file" ] && [ "$force_overwrite" = false ]; then
        echo "Secrets file already exists. Use --force to overwrite."
        read -p "Do you want to add/update secrets? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            return 0
        fi
    fi
    
    # Create backup if file exists
    if [ -f "$secrets_file" ]; then
        cp "$secrets_file" "$secrets_file.backup.$(date +%Y%m%d_%H%M%S)"
        log_info "Created backup of existing secrets"
    fi
    
    # Copy template if file doesn't exist
    if [ ! -f "$secrets_file" ]; then
        cp "$SECRETS_TEMPLATE" "$secrets_file"
    fi
    
    echo "=== Secrets Setup Wizard ==="
    echo "This wizard will help you configure API keys and secrets."
    echo "You can skip any section by pressing Enter."
    echo
    
    # AWS Setup
    secrets_setup_aws
    
    # Azure Setup  
    secrets_setup_azure
    
    # DigitalOcean Setup
    secrets_setup_digitalocean
    
    # GitHub Setup
    secrets_setup_github
    
    echo
    log_success "Secrets setup completed!"
    echo "Your secrets are stored in: $secrets_file"
    echo "Remember to:"
    echo "  1. Keep this file secure (chmod 600)"
    echo "  2. Never commit it to version control"
    echo "  3. Use 'auto secrets validate' to test your keys"
    echo "  4. Use 'auto secrets encrypt' to encrypt the vault"
}

# AWS credentials setup
secrets_setup_aws() {
    echo
    echo "--- AWS Credentials ---"
    read -p "Enter AWS Access Key ID (or skip): " aws_key
    if [ -n "$aws_key" ]; then
        read -s -p "Enter AWS Secret Access Key: " aws_secret
        echo
        read -p "Enter AWS Default Region [us-east-1]: " aws_region
        aws_region=${aws_region:-us-east-1}
        
        secrets_update_value "AWS_ACCESS_KEY_ID" "$aws_key"
        secrets_update_value "AWS_SECRET_ACCESS_KEY" "$aws_secret"
        secrets_update_value "AWS_DEFAULT_REGION" "$aws_region"
        
        log_info "AWS credentials configured"
    fi
}

# Azure credentials setup
secrets_setup_azure() {
    echo
    echo "--- Azure Credentials ---"
    read -p "Enter Azure Client ID (or skip): " azure_client_id
    if [ -n "$azure_client_id" ]; then
        read -s -p "Enter Azure Client Secret: " azure_secret
        echo
        read -p "Enter Azure Tenant ID: " azure_tenant
        read -p "Enter Azure Subscription ID: " azure_sub
        
        secrets_update_value "AZURE_CLIENT_ID" "$azure_client_id"
        secrets_update_value "AZURE_CLIENT_SECRET" "$azure_secret"
        secrets_update_value "AZURE_TENANT_ID" "$azure_tenant"
        secrets_update_value "AZURE_SUBSCRIPTION_ID" "$azure_sub"
        
        log_info "Azure credentials configured"
    fi
}

# DigitalOcean credentials setup
secrets_setup_digitalocean() {
    echo
    echo "--- DigitalOcean Credentials ---"
    read -s -p "Enter DigitalOcean API Token (or skip): " do_token
    echo
    if [ -n "$do_token" ]; then
        secrets_update_value "DIGITALOCEAN_TOKEN" "$do_token"
        log_info "DigitalOcean credentials configured"
    fi
}

# GitHub credentials setup
secrets_setup_github() {
    echo
    echo "--- GitHub Credentials ---"
    read -p "Enter GitHub Username (or skip): " gh_user
    if [ -n "$gh_user" ]; then
        read -s -p "Enter GitHub Personal Access Token: " gh_token
        echo
        read -p "Enter GitHub Organization (optional): " gh_org
        
        secrets_update_value "GITHUB_USERNAME" "$gh_user"
        secrets_update_value "GITHUB_TOKEN" "$gh_token"
        if [ -n "$gh_org" ]; then
            secrets_update_value "GITHUB_ORGANIZATION" "$gh_org"
        fi
        
        log_info "GitHub credentials configured"
    fi
}

# Update secret value in file
secrets_update_value() {
    local key="$1"
    local value="$2"
    local secrets_file="$SECRETS_DIR/.env"
    
    if grep -q "^${key}=" "$secrets_file"; then
        # Update existing value
        sed -i.bak "s|^${key}=.*|${key}=${value}|" "$secrets_file"
        rm -f "$secrets_file.bak"
    else
        # Add new value
        echo "${key}=${value}" >> "$secrets_file"
    fi
}

# Validate all secrets
secrets_validate() {
    log_info "Validating configured secrets..."
    
    local secrets_file="$SECRETS_DIR/.env"
    if [ ! -f "$secrets_file" ]; then
        log_error "No secrets file found. Run 'auto secrets setup' first."
        return 1
    fi
    
    source "$secrets_file"
    
    local validation_failed=false
    
    # Validate AWS
    if [ -n "$AWS_ACCESS_KEY_ID" ] && [ -n "$AWS_SECRET_ACCESS_KEY" ]; then
        log_info "Validating AWS credentials..."
        if command -v aws >/dev/null 2>&1; then
            if AWS_ACCESS_KEY_ID="$AWS_ACCESS_KEY_ID" AWS_SECRET_ACCESS_KEY="$AWS_SECRET_ACCESS_KEY" \
               aws sts get-caller-identity >/dev/null 2>&1; then
                log_success "✅ AWS credentials valid"
            else
                log_error "❌ AWS credentials invalid"
                validation_failed=true
            fi
        else
            log_warn "⚠️ AWS CLI not installed, skipping validation"
        fi
    fi
    
    # Validate Azure
    if [ -n "$AZURE_CLIENT_ID" ] && [ -n "$AZURE_CLIENT_SECRET" ] && [ -n "$AZURE_TENANT_ID" ]; then
        log_info "Validating Azure credentials..."
        if command -v az >/dev/null 2>&1; then
            if az login --service-principal -u "$AZURE_CLIENT_ID" -p "$AZURE_CLIENT_SECRET" --tenant "$AZURE_TENANT_ID" >/dev/null 2>&1; then
                log_success "✅ Azure credentials valid"
                az logout >/dev/null 2>&1
            else
                log_error "❌ Azure credentials invalid"
                validation_failed=true
            fi
        else
            log_warn "⚠️ Azure CLI not installed, skipping validation"
        fi
    fi
    
    # Validate DigitalOcean
    if [ -n "$DIGITALOCEAN_TOKEN" ]; then
        log_info "Validating DigitalOcean token..."
        if command -v doctl >/dev/null 2>&1; then
            if DIGITALOCEAN_ACCESS_TOKEN="$DIGITALOCEAN_TOKEN" doctl account get >/dev/null 2>&1; then
                log_success "✅ DigitalOcean token valid"
            else
                log_error "❌ DigitalOcean token invalid"
                validation_failed=true
            fi
        else
            log_warn "⚠️ doctl not installed, skipping validation"
        fi
    fi
    
    # Validate GitHub
    if [ -n "$GITHUB_TOKEN" ]; then
        log_info "Validating GitHub token..."
        if command -v gh >/dev/null 2>&1; then
            if GITHUB_TOKEN="$GITHUB_TOKEN" gh auth status >/dev/null 2>&1; then
                log_success "✅ GitHub token valid"
            else
                log_error "❌ GitHub token invalid"
                validation_failed=true
            fi
        else
            log_warn "⚠️ GitHub CLI not installed, skipping validation"
        fi
    fi
    
    if [ "$validation_failed" = true ]; then
        log_error "Some credentials failed validation"
        return 1
    else
        log_success "All configured credentials are valid"
        return 0
    fi
}

# List configured secrets (keys only, not values)
secrets_list() {
    local secrets_file="$SECRETS_DIR/.env"
    if [ ! -f "$secrets_file" ]; then
        log_error "No secrets file found"
        return 1
    fi
    
    log_info "Configured secrets:"
    grep -v '^#' "$secrets_file" | grep -v '^$' | grep '=' | while IFS='=' read -r key value; do
        if [ -n "$value" ]; then
            echo "  ✅ $key"
        else
            echo "  ❌ $key (empty)"
        fi
    done
}

# Check which API keys are missing
secrets_check_requirements() {
    log_info "Checking API key requirements..."
    
    local secrets_file="$SECRETS_DIR/.env"
    local missing_keys=()
    
    # Load existing secrets if available
    if [ -f "$secrets_file" ]; then
        source "$secrets_file"
    fi
    
    echo
    echo "=== Cloud Provider Credentials ==="
    
    # AWS
    if [ -z "$AWS_ACCESS_KEY_ID" ] || [ -z "$AWS_SECRET_ACCESS_KEY" ]; then
        echo "❌ AWS credentials missing (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY)"
        missing_keys+=("aws")
    else
        echo "✅ AWS credentials configured"
    fi
    
    # Azure
    if [ -z "$AZURE_CLIENT_ID" ] || [ -z "$AZURE_CLIENT_SECRET" ] || [ -z "$AZURE_TENANT_ID" ]; then
        echo "❌ Azure credentials missing (AZURE_CLIENT_ID, AZURE_CLIENT_SECRET, AZURE_TENANT_ID)"
        missing_keys+=("azure")
    else
        echo "✅ Azure credentials configured"
    fi
    
    # DigitalOcean
    if [ -z "$DIGITALOCEAN_TOKEN" ]; then
        echo "❌ DigitalOcean token missing (DIGITALOCEAN_TOKEN)"
        missing_keys+=("digitalocean")
    else
        echo "✅ DigitalOcean token configured"
    fi
    
    echo
    echo "=== Development Tools ==="
    
    # GitHub
    if [ -z "$GITHUB_TOKEN" ]; then
        echo "❌ GitHub token missing (GITHUB_TOKEN)"
        missing_keys+=("github")
    else
        echo "✅ GitHub token configured"
    fi
    
    # Docker Hub
    if [ -z "$DOCKER_USERNAME" ] || [ -z "$DOCKER_PASSWORD" ]; then
        echo "❌ Docker Hub credentials missing (DOCKER_USERNAME, DOCKER_PASSWORD)"
        missing_keys+=("docker")
    else
        echo "✅ Docker Hub credentials configured"
    fi
    
    echo
    echo "=== CLI Tools Status ==="
    
    # Check CLI tools
    command -v aws >/dev/null && echo "✅ AWS CLI installed" || echo "❌ AWS CLI not installed"
    command -v az >/dev/null && echo "✅ Azure CLI installed" || echo "❌ Azure CLI not installed"
    command -v doctl >/dev/null && echo "✅ doctl installed" || echo "❌ doctl not installed"
    command -v gh >/dev/null && echo "✅ GitHub CLI installed" || echo "❌ GitHub CLI not installed"
    command -v kubectl >/dev/null && echo "✅ kubectl installed" || echo "❌ kubectl not installed"
    command -v helm >/dev/null && echo "✅ Helm installed" || echo "❌ Helm not installed"
    
    echo
    if [ ${#missing_keys[@]} -eq 0 ]; then
        log_success "All required API keys are configured!"
        return 0
    else
        log_warn "Missing API keys for: ${missing_keys[*]}"
        echo
        echo "To configure missing credentials:"
        for provider in "${missing_keys[@]}"; do
            echo "  auto secrets $provider"
        done
        echo
        echo "Or run the full setup wizard:"
        echo "  auto secrets setup"
        return 1
    fi
}

# Export secrets for current session
secrets_export() {
    local format="env"
    local output_file=""
    
    # Parse options
    while [[ $# -gt 0 ]]; do
        case $1 in
            --format)
                format="$2"
                shift 2
                ;;
            --output)
                output_file="$2"
                shift 2
                ;;
            *)
                shift
                ;;
        esac
    done
    
    local secrets_file="$SECRETS_DIR/.env"
    if [ ! -f "$secrets_file" ]; then
        log_error "No secrets file found"
        return 1
    fi
    
    case "$format" in
        "env")
            if [ -n "$output_file" ]; then
                cp "$secrets_file" "$output_file"
                log_info "Secrets exported to: $output_file"
            else
                cat "$secrets_file"
            fi
            ;;
        "json")
            log_info "Converting secrets to JSON format..."
            echo "{"
            grep -v '^#' "$secrets_file" | grep -v '^$' | grep '=' | while IFS='=' read -r key value; do
                echo "  \"$key\": \"$value\","
            done | sed '$ s/,$//'
            echo "}"
            ;;
        *)
            log_error "Unsupported format: $format"
            return 1
            ;;
    esac
}

# Encrypt secrets vault
secrets_encrypt() {
    local secrets_file="$SECRETS_DIR/.env"
    if [ ! -f "$secrets_file" ]; then
        log_error "No secrets file found"
        return 1
    fi
    
    read -s -p "Enter encryption password: " password
    echo
    
    if openssl enc -aes-256-cbc -salt -in "$secrets_file" -out "$ENCRYPTED_SECRETS" -pass pass:"$password"; then
        log_success "Secrets encrypted successfully"
        log_info "Encrypted vault: $ENCRYPTED_SECRETS"
        
        # Optionally remove plaintext file
        read -p "Remove plaintext secrets file? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            rm "$secrets_file"
            log_info "Plaintext secrets file removed"
        fi
    else
        log_error "Failed to encrypt secrets"
        return 1
    fi
}

# Decrypt secrets vault
secrets_decrypt() {
    if [ ! -f "$ENCRYPTED_SECRETS" ]; then
        log_error "No encrypted vault found"
        return 1
    fi
    
    read -s -p "Enter decryption password: " password
    echo
    
    local secrets_file="$SECRETS_DIR/.env"
    if openssl enc -aes-256-cbc -d -in "$ENCRYPTED_SECRETS" -out "$secrets_file" -pass pass:"$password"; then
        chmod 600 "$secrets_file"
        log_success "Secrets decrypted successfully"
        log_info "Plaintext secrets: $secrets_file"
    else
        log_error "Failed to decrypt secrets (wrong password?)"
        return 1
    fi
}

# Main secrets function
secrets_main() {
    case "${1:-help}" in
        "init") secrets_init ;;
        "setup") shift; secrets_setup "$@" ;;
        "validate") secrets_validate ;;
        "list") secrets_list ;;
        "encrypt") secrets_encrypt ;;
        "decrypt") secrets_decrypt ;;
        "export") shift; secrets_export "$@" ;;
        "check-requirements"|"check") secrets_check_requirements ;;
        "template") cat "$SECRETS_TEMPLATE" ;;
        "aws") secrets_setup_aws ;;
        "azure") secrets_setup_azure ;;
        "digitalocean") secrets_setup_digitalocean ;;
        "github") secrets_setup_github ;;
        "help"|*) secrets_help ;;
    esac
}