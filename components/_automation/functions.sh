#!/bin/sh
# components/automation/functions.sh - Automation framework functions

# =============================================================================
# Status and Info
# =============================================================================

# Show automation framework status
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

# =============================================================================
# Installation
# =============================================================================

# Install/setup automation framework
install_automation() {
    echo "Setting up automation framework..."

    if [ ! -d "$AUTO_HOME" ]; then
        echo "Automation directory not found: $AUTO_HOME"
        return 1
    fi

    chmod +x "$AUTO_HOME/auto" 2>/dev/null
    chmod +x "$AUTO_HOME/framework"/*.sh 2>/dev/null
    chmod +x "$AUTO_HOME/modules"/*.sh 2>/dev/null

    if [ -f "$AUTO_CLI" ]; then
        "$AUTO_CLI" --version 2>/dev/null || echo "Version: unknown"
        echo "Automation framework ready"
        echo "Run 'auto --help' to get started"
    else
        echo "Automation CLI not found"
        return 1
    fi
}

# =============================================================================
# Quick Project Setup
# =============================================================================

# Quick project setup using automation
quick_project() {
    local project_type="$1"
    local project_name="$2"

    if [ -z "$project_type" ] || [ -z "$project_name" ]; then
        echo "Usage: quick_project <type> <name>"
        echo "Types: python, go, typescript"
        return 1
    fi

    if ! command -v auto >/dev/null 2>&1; then
        echo "Automation framework not available"
        return 1
    fi

    auto dev init "$project_type" "$project_name"
    cd "$HOME/projects/$project_name" 2>/dev/null || cd "$project_name" 2>/dev/null
}

# =============================================================================
# Logs and Debugging
# =============================================================================

# View automation logs
auto_logs() {
    local log_file="$AUTO_HOME/logs/automation.log"
    if [ -f "$log_file" ]; then
        if [ "$1" = "-f" ]; then
            tail -f "$log_file"
        else
            tail -50 "$log_file"
        fi
    else
        echo "No log file found at: $log_file"
    fi
}

# Clear automation logs
auto_logs_clear() {
    local log_file="$AUTO_HOME/logs/automation.log"
    if [ -f "$log_file" ]; then
        : > "$log_file"
        echo "Logs cleared"
    fi
}
