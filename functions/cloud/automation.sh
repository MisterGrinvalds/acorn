#!/bin/sh
# Automation Framework Integration
# Requires: DOTFILES_ROOT environment variable

# Automation environment setup
export AUTO_HOME="${DOTFILES_ROOT}/.automation"
export AUTO_CLI="$AUTO_HOME/auto"

# Add automation CLI to PATH if it exists
if [ -f "$AUTO_CLI" ]; then
    case ":$PATH:" in
        *":$AUTO_HOME:"*) ;;
        *) PATH="$AUTO_HOME:$PATH" ;;
    esac
fi

# Automation aliases
alias auto='$AUTO_CLI'
alias autodev='auto dev'
alias autok8s='auto k8s'
alias autogithub='auto github'
alias autosystem='auto system'
alias autoconfig='auto config'
alias autocloud='auto cloud'

# Quick automation shortcuts
alias ainit='auto dev init'
alias adeploy='auto k8s deploy'
alias apr='auto github pr create'
alias abackup='auto system backup'
alias acleanup='auto system cleanup'

# Automation helper functions
auto_status() {
    if [ -f "$AUTO_CLI" ]; then
        echo "Automation framework: Available"
        echo "CLI location: $AUTO_CLI"
        echo "Config directory: $AUTO_HOME/config"
        echo "Logs directory: $AUTO_HOME/logs"

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
    echo "Setting up automation framework..."

    chmod +x "$AUTO_HOME/auto" 2>/dev/null
    chmod +x "$AUTO_HOME/framework"/*.sh 2>/dev/null
    chmod +x "$AUTO_HOME/modules"/*.sh 2>/dev/null

    if [ -f "$AUTO_CLI" ]; then
        "$AUTO_CLI" --version
        echo "Automation framework ready"
        echo "Run 'auto --help' to get started"
    else
        echo "Automation framework files not found"
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
