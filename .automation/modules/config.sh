#!/bin/bash
# Configuration Management Module

# Module configuration
readonly CONFIG_PROFILES_DIR="$AUTO_CONFIG/profiles"
readonly CONFIG_TEMPLATES_DIR="$AUTO_CONFIG/templates"
readonly CONFIG_ENVIRONMENTS_DIR="$AUTO_CONFIG/environments"

# Initialize directories
mkdir -p "$CONFIG_PROFILES_DIR" "$CONFIG_TEMPLATES_DIR" "$CONFIG_ENVIRONMENTS_DIR"

# Help function for config module
config_help() {
    cat << EOF
Configuration Management

USAGE:
    auto config <command> [options]

COMMANDS:
    profile <action>              Profile management (create, switch, list, delete)
    template <action>             Template management (create, apply, list)
    environment <action>          Environment management (create, sync, list)
    backup                        Backup current configuration
    restore <backup-file>         Restore configuration from backup
    sync                          Sync configuration across machines
    validate                      Validate configuration files
    diff [profile1] [profile2]    Compare configurations

EXAMPLES:
    auto config profile create work
    auto config profile switch development
    auto config template apply python-project
    auto config environment create production
    auto config backup
    auto config sync --push origin
EOF
}

# Profile management
create_profile() {
    local profile_name="$1"
    local description="${2:-Created on $(date)}"
    
    if [ -z "$profile_name" ]; then
        log "ERROR" "Profile name is required"
        exit 1
    fi
    
    local profile_dir="$CONFIG_PROFILES_DIR/$profile_name"
    
    if [ -d "$profile_dir" ]; then
        log "ERROR" "Profile '$profile_name' already exists"
        exit 1
    fi
    
    mkdir -p "$profile_dir"
    
    # Create profile metadata
    cat > "$profile_dir/profile.json" << EOF
{
    "name": "$profile_name",
    "description": "$description",
    "created": "$(date -Iseconds)",
    "version": "1.0.0",
    "active": false
}
EOF
    
    # Copy current configuration
    cp -r "$DOTFILES/.bash_profile.dir" "$profile_dir/"
    cp -r "$DOTFILES/.bash_tools" "$profile_dir/"
    cp "$DOTFILES/.bash_profile" "$profile_dir/"
    
    # Save current environment variables
    env | grep -E '^(AUTO_|DEV_|K8S_|GITHUB_)' > "$profile_dir/environment.env"
    
    log "SUCCESS" "Profile '$profile_name' created"
}

switch_profile() {
    local profile_name="$1"
    
    if [ -z "$profile_name" ]; then
        log "INFO" "Available profiles:"
        list_profiles
        return 0
    fi
    
    local profile_dir="$CONFIG_PROFILES_DIR/$profile_name"
    
    if [ ! -d "$profile_dir" ]; then
        log "ERROR" "Profile '$profile_name' does not exist"
        exit 1
    fi
    
    # Backup current configuration
    backup_current_config
    
    # Apply profile configuration
    cp -r "$profile_dir/.bash_profile.dir" "$DOTFILES/"
    cp -r "$profile_dir/.bash_tools" "$DOTFILES/"
    cp "$profile_dir/.bash_profile" "$DOTFILES/"
    
    # Load environment variables
    if [ -f "$profile_dir/environment.env" ]; then
        source "$profile_dir/environment.env"
    fi
    
    # Update profile metadata
    jq '.active = true' "$profile_dir/profile.json" > "$profile_dir/profile.json.tmp"
    mv "$profile_dir/profile.json.tmp" "$profile_dir/profile.json"
    
    # Mark other profiles as inactive
    for other_profile in "$CONFIG_PROFILES_DIR"/*; do
        if [ -d "$other_profile" ] && [ "$other_profile" != "$profile_dir" ]; then
            if [ -f "$other_profile/profile.json" ]; then
                jq '.active = false' "$other_profile/profile.json" > "$other_profile/profile.json.tmp"
                mv "$other_profile/profile.json.tmp" "$other_profile/profile.json"
            fi
        fi
    done
    
    log "SUCCESS" "Switched to profile: $profile_name"
    log "INFO" "Please restart your shell or run: source ~/.bash_profile"
}

list_profiles() {
    log "INFO" "Available profiles:"
    
    for profile_dir in "$CONFIG_PROFILES_DIR"/*; do
        if [ -d "$profile_dir" ] && [ -f "$profile_dir/profile.json" ]; then
            local name=$(jq -r '.name' "$profile_dir/profile.json")
            local description=$(jq -r '.description' "$profile_dir/profile.json")
            local active=$(jq -r '.active' "$profile_dir/profile.json")
            local created=$(jq -r '.created' "$profile_dir/profile.json")
            
            local status=""
            [ "$active" = "true" ] && status=" (active)"
            
            printf "  %-15s %s%s\n" "$name" "$description" "$status"
            printf "  %-15s Created: %s\n" "" "$created"
            echo ""
        fi
    done
}

delete_profile() {
    local profile_name="$1"
    
    if [ -z "$profile_name" ]; then
        log "ERROR" "Profile name is required"
        exit 1
    fi
    
    local profile_dir="$CONFIG_PROFILES_DIR/$profile_name"
    
    if [ ! -d "$profile_dir" ]; then
        log "ERROR" "Profile '$profile_name' does not exist"
        exit 1
    fi
    
    # Check if profile is active
    if [ -f "$profile_dir/profile.json" ]; then
        local active=$(jq -r '.active' "$profile_dir/profile.json")
        if [ "$active" = "true" ]; then
            log "ERROR" "Cannot delete active profile. Switch to another profile first."
            exit 1
        fi
    fi
    
    if confirm "Delete profile '$profile_name'? This cannot be undone."; then
        rm -rf "$profile_dir"
        log "SUCCESS" "Profile '$profile_name' deleted"
    fi
}

# Template management
create_template() {
    local template_name="$1"
    local template_type="${2:-custom}"
    
    if [ -z "$template_name" ]; then
        log "ERROR" "Template name is required"
        exit 1
    fi
    
    local template_dir="$CONFIG_TEMPLATES_DIR/$template_name"
    mkdir -p "$template_dir"
    
    case "$template_type" in
        "python-dev")
            create_python_dev_template "$template_dir"
            ;;
        "go-dev")
            create_go_dev_template "$template_dir"
            ;;
        "k8s-ops")
            create_k8s_ops_template "$template_dir"
            ;;
        "custom")
            create_custom_template "$template_dir"
            ;;
        *)
            log "ERROR" "Unknown template type: $template_type"
            exit 1
            ;;
    esac
    
    log "SUCCESS" "Template '$template_name' created"
}

create_python_dev_template() {
    local template_dir="$1"
    
    cat > "$template_dir/template.json" << 'EOF'
{
    "name": "Python Development",
    "description": "Configuration for Python development with FastAPI",
    "type": "python-dev",
    "variables": {
        "PYTHON_VERSION": "3.11",
        "DEFAULT_VENV_PATH": "$HOME/envs",
        "FASTAPI_HOST": "0.0.0.0",
        "FASTAPI_PORT": "8000"
    }
}
EOF
    
    cat > "$template_dir/aliases.sh" << 'EOF'
# Python Development Aliases
alias py='python3'
alias pym='python3 -m'
alias pyv='python3 --version'
alias pip='python3 -m pip'
alias venv='source venv/bin/activate'
alias deactivate='deactivate'

# FastAPI specific
alias fastdev='uvicorn main:app --reload --host {{FASTAPI_HOST}} --port {{FASTAPI_PORT}}'
alias fastprod='uvicorn main:app --host {{FASTAPI_HOST}} --port {{FASTAPI_PORT}}'

# Testing
alias pytest='python -m pytest'
alias pytestv='python -m pytest -v'
alias pytestcov='python -m pytest --cov=.'
EOF
    
    cat > "$template_dir/environment.sh" << 'EOF'
# Python Development Environment
export PYTHON_VERSION={{PYTHON_VERSION}}
export PYTHONPATH="${PYTHONPATH}:."
export DEFAULT_VENV_PATH={{DEFAULT_VENV_PATH}}
export FASTAPI_HOST={{FASTAPI_HOST}}
export FASTAPI_PORT={{FASTAPI_PORT}}

# Add Python tools to PATH
[[ ":$PATH:" != *":$HOME/.local/bin:"* ]] && PATH="$HOME/.local/bin:$PATH"
EOF
}

apply_template() {
    local template_name="$1"
    local template_dir="$CONFIG_TEMPLATES_DIR/$template_name"
    
    if [ ! -d "$template_dir" ]; then
        log "ERROR" "Template '$template_name' does not exist"
        exit 1
    fi
    
    if [ ! -f "$template_dir/template.json" ]; then
        log "ERROR" "Invalid template: missing template.json"
        exit 1
    fi
    
    log "INFO" "Applying template: $template_name"
    
    # Read template variables
    local variables=$(jq -r '.variables | to_entries[] | "\(.key)=\(.value)"' "$template_dir/template.json")
    
    # Process and apply template files
    for template_file in "$template_dir"/*.sh; do
        if [ -f "$template_file" ]; then
            local filename=$(basename "$template_file")
            local target_file="$DOTFILES/.bash_profile.dir/$filename"
            
            # Process template variables
            process_template "$template_file" "$target_file" $variables
            
            log "DEBUG" "Applied template file: $filename"
        fi
    done
    
    log "SUCCESS" "Template '$template_name' applied"
    log "INFO" "Please restart your shell or run: source ~/.bash_profile"
}

# Environment management
create_environment() {
    local env_name="$1"
    local env_type="${2:-development}"
    
    if [ -z "$env_name" ]; then
        log "ERROR" "Environment name is required"
        exit 1
    fi
    
    local env_dir="$CONFIG_ENVIRONMENTS_DIR/$env_name"
    mkdir -p "$env_dir"
    
    cat > "$env_dir/environment.json" << EOF
{
    "name": "$env_name",
    "type": "$env_type",
    "created": "$(date -Iseconds)",
    "variables": {}
}
EOF
    
    cat > "$env_dir/variables.env" << EOF
# Environment: $env_name
# Type: $env_type
# Created: $(date)

# Add your environment variables here
# Example:
# DATABASE_URL=postgresql://user:pass@localhost:5432/db
# API_KEY=your-api-key-here
# DEBUG=true
EOF
    
    log "SUCCESS" "Environment '$env_name' created"
    log "INFO" "Edit variables: $env_dir/variables.env"
}

load_environment() {
    local env_name="$1"
    local env_dir="$CONFIG_ENVIRONMENTS_DIR/$env_name"
    
    if [ ! -d "$env_dir" ]; then
        log "ERROR" "Environment '$env_name' does not exist"
        exit 1
    fi
    
    if [ -f "$env_dir/variables.env" ]; then
        source "$env_dir/variables.env"
        log "SUCCESS" "Environment '$env_name' loaded"
    else
        log "WARN" "No variables file found for environment '$env_name'"
    fi
}

# Configuration backup and restore
backup_current_config() {
    local backup_name="config_backup_$(date +%Y%m%d_%H%M%S)"
    local backup_dir="$AUTO_CACHE/config-backups/$backup_name"
    
    mkdir -p "$backup_dir"
    
    # Backup dotfiles
    cp -r "$DOTFILES/.bash_profile.dir" "$backup_dir/"
    cp -r "$DOTFILES/.bash_tools" "$backup_dir/"
    cp "$DOTFILES/.bash_profile" "$backup_dir/"
    
    # Backup automation config
    cp -r "$AUTO_CONFIG" "$backup_dir/automation-config"
    
    # Create metadata
    cat > "$backup_dir/backup.json" << EOF
{
    "name": "$backup_name",
    "created": "$(date -Iseconds)",
    "type": "full-config",
    "hostname": "$(hostname)",
    "user": "$(whoami)"
}
EOF
    
    # Create archive
    tar -czf "$backup_dir.tar.gz" -C "$(dirname "$backup_dir")" "$(basename "$backup_dir")"
    rm -rf "$backup_dir"
    
    log "SUCCESS" "Configuration backed up: $backup_dir.tar.gz"
    echo "$backup_dir.tar.gz"
}

validate_config() {
    log "INFO" "Validating configuration files..."
    
    local errors=0
    
    # Check main profile file
    if [ -f "$DOTFILES/.bash_profile" ]; then
        bash -n "$DOTFILES/.bash_profile" || ((errors++))
    else
        log "ERROR" "Main profile file missing: .bash_profile"
        ((errors++))
    fi
    
    # Check module files
    for module_file in "$DOTFILES/.bash_profile.dir"/*.sh; do
        if [ -f "$module_file" ]; then
            bash -n "$module_file" || ((errors++))
        fi
    done
    
    # Check tool files
    for tool_file in "$DOTFILES/.bash_tools"/*.sh; do
        if [ -f "$tool_file" ]; then
            bash -n "$tool_file" || ((errors++))
        fi
    done
    
    if [ $errors -eq 0 ]; then
        log "SUCCESS" "Configuration validation passed"
    else
        log "ERROR" "Configuration validation failed with $errors errors"
        exit 1
    fi
}

# Main config module function
config_main() {
    local command="${1:-help}"
    shift || true
    
    case "$command" in
        "help"|"-h"|"--help")
            config_help
            ;;
        "profile")
            local action="$1"
            case "$action" in
                "create") create_profile "$2" "$3" ;;
                "switch") switch_profile "$2" ;;
                "list") list_profiles ;;
                "delete") delete_profile "$2" ;;
                *) log "ERROR" "Unknown profile action: $action" ;;
            esac
            ;;
        "template")
            local action="$1"
            case "$action" in
                "create") create_template "$2" "$3" ;;
                "apply") apply_template "$2" ;;
                "list") ls -1 "$CONFIG_TEMPLATES_DIR" ;;
                *) log "ERROR" "Unknown template action: $action" ;;
            esac
            ;;
        "environment")
            local action="$1"
            case "$action" in
                "create") create_environment "$2" "$3" ;;
                "load") load_environment "$2" ;;
                "list") ls -1 "$CONFIG_ENVIRONMENTS_DIR" ;;
                *) log "ERROR" "Unknown environment action: $action" ;;
            esac
            ;;
        "backup")
            backup_current_config
            ;;
        "validate")
            validate_config
            ;;
        *)
            log "ERROR" "Unknown command: $command"
            config_help
            exit 1
            ;;
    esac
}