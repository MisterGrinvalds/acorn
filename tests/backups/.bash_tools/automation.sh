# Automation Framework Integration
# Integrates the automation system with the shell environment

# Automation environment setup
export AUTO_HOME="${DOTFILES}/.automation"
export AUTO_CLI="$AUTO_HOME/auto"

# Add automation CLI to PATH if it exists
if [ -f "$AUTO_CLI" ]; then
    [[ ":$PATH:" != *":$AUTO_HOME:"* ]] && PATH="$AUTO_HOME:$PATH"
fi

# Automation aliases for quick access
alias auto='$AUTO_CLI'
alias autodev='auto dev'
alias autok8s='auto k8s' 
alias autogithub='auto github'
alias autosystem='auto system'
alias autoconfig='auto config'
alias autocloud='auto cloud'
alias autoaws='auto aws'
alias autoazure='auto azure'
alias autodo='auto digitalocean'

# Quick automation shortcuts
alias ainit='auto dev init'          # Quick project initialization
alias adeploy='auto k8s deploy'      # Quick k8s deployment
alias apr='auto github pr create'    # Quick PR creation
alias abackup='auto system backup'   # Quick backup
alias acleanup='auto system cleanup' # Quick cleanup

# Cloud shortcuts
alias acloud='auto cloud status'     # Quick cloud status
alias awsec2='auto aws ec2 list'     # Quick EC2 list
alias azvm='auto azure vm list'      # Quick VM list
alias dodrop='auto digitalocean droplets list' # Quick droplets list

# Automation helper functions
auto_status() {
    if [ -f "$AUTO_CLI" ]; then
        echo "Automation framework: Available"
        echo "CLI location: $AUTO_CLI"
        echo "Config directory: $AUTO_HOME/config"
        echo "Logs directory: $AUTO_HOME/logs"
        
        # Show recent activity
        if [ -f "$AUTO_HOME/logs/automation.log" ]; then
            echo ""
            echo "Recent activity:"
            tail -5 "$AUTO_HOME/logs/automation.log" 2>/dev/null || echo "No recent activity"
        fi
    else
        echo "Automation framework: Not installed"
        echo "Run 'install_automation' to set up the framework"
    fi
}

# Install/setup automation framework
install_automation() {
    log "INFO" "Setting up automation framework..."
    
    # Ensure framework files are executable
    chmod +x "$AUTO_HOME/auto" 2>/dev/null
    chmod +x "$AUTO_HOME/framework"/*.sh 2>/dev/null
    chmod +x "$AUTO_HOME/modules"/*.sh 2>/dev/null
    
    # Initialize framework
    if [ -f "$AUTO_CLI" ]; then
        "$AUTO_CLI" --version
        log "SUCCESS" "Automation framework ready"
        log "INFO" "Run 'auto --help' to get started"
    else
        log "ERROR" "Automation framework files not found"
        return 1
    fi
}

# Quick project setup using automation
quick_project() {
    local project_type="$1"
    local project_name="$2"
    
    if [ -z "$project_type" ] || [ -z "$project_name" ]; then
        echo "Usage: quick_project <type> <name>"
        echo "Types: python, go, typescript"
        return 1
    fi
    
    auto dev init "$project_type" "$project_name"
    cd "$HOME/projects/$project_name" 2>/dev/null || cd "$project_name"
}

# Quick k8s operations
k8s_quick_deploy() {
    local app_name="$1"
    local environment="${2:-development}"
    
    if [ -z "$app_name" ]; then
        echo "Usage: k8s_quick_deploy <app-name> [environment]"
        return 1
    fi
    
    auto k8s deploy "$app_name" "$environment"
}

# GitHub quick operations
github_quick_pr() {
    local title="$1"
    local description="${2:-}"
    
    if [ -z "$title" ]; then
        echo "Usage: github_quick_pr <title> [description]"
        return 1
    fi
    
    auto github pr create "$title" "$description"
}

# System quick operations
system_quick_backup() {
    local path="${1:-.}"
    auto system backup "$path"
}

# Auto-completion for automation commands (bash)
if [ "$CURRENT_SHELL" = "bash" ]; then
    _auto_completion() {
        local cur prev opts
        COMPREPLY=()
        cur="${COMP_WORDS[COMP_CWORD]}"
        prev="${COMP_WORDS[COMP_CWORD-1]}"
        
        case $prev in
            auto)
                opts="dev k8s github system config"
                COMPREPLY=($(compgen -W "${opts}" -- ${cur}))
                return 0
                ;;
            dev)
                opts="init test build deploy format lint"
                COMPREPLY=($(compgen -W "${opts}" -- ${cur}))
                return 0
                ;;
            k8s)
                opts="cluster deploy scale logs port-forward manifests helm monitoring backup cleanup"
                COMPREPLY=($(compgen -W "${opts}" -- ${cur}))
                return 0
                ;;
            github)
                opts="repo pr issue workflow release security"
                COMPREPLY=($(compgen -W "${opts}" -- ${cur}))
                return 0
                ;;
            system)
                opts="setup backup monitor cleanup security services update"
                COMPREPLY=($(compgen -W "${opts}" -- ${cur}))
                return 0
                ;;
            config)
                opts="profile template environment backup validate"
                COMPREPLY=($(compgen -W "${opts}" -- ${cur}))
                return 0
                ;;
        esac
    }
    
    complete -F _auto_completion auto
fi

# Initialize automation on shell startup (optional)
if [ "${AUTO_INIT_ON_STARTUP:-true}" = "true" ]; then
    # Check if automation framework exists and initialize silently
    if [ -f "$AUTO_CLI" ] && [ ! -f "$AUTO_HOME/.initialized" ]; then
        install_automation >/dev/null 2>&1 && touch "$AUTO_HOME/.initialized"
    fi
fi

# Automation update checker
check_automation_updates() {
    if [ -f "$AUTO_HOME/logs/last_update_check" ]; then
        local last_check=$(cat "$AUTO_HOME/logs/last_update_check")
        local current_time=$(date +%s)
        local week_seconds=604800
        
        if [ $((current_time - last_check)) -gt $week_seconds ]; then
            echo "💡 Automation framework update available. Run: auto system update"
            echo "$current_time" > "$AUTO_HOME/logs/last_update_check"
        fi
    else
        mkdir -p "$AUTO_HOME/logs"
        echo "$(date +%s)" > "$AUTO_HOME/logs/last_update_check"
    fi
}