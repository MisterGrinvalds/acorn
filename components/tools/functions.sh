#!/bin/sh
# components/tools/functions.sh - System tools management functions

# =============================================================================
# Automation Framework Integration
# =============================================================================

# Quick tool status check
tools_status() {
    if command -v auto >/dev/null 2>&1; then
        auto tools status
    else
        echo "Automation framework not found"
        return 1
    fi
}

# List all tools
list_tools() {
    if command -v auto >/dev/null 2>&1; then
        auto tools list
    else
        echo "Automation framework not found"
        return 1
    fi
}

# Check tool versions
check_tools() {
    if command -v auto >/dev/null 2>&1; then
        auto tools check "$@"
    else
        echo "Automation framework not found"
        return 1
    fi
}

# Update tools interactively
update_tools() {
    if command -v auto >/dev/null 2>&1; then
        auto tools update "$@"
    else
        echo "Automation framework not found"
        return 1
    fi
}

# Show missing tools
missing_tools() {
    if command -v auto >/dev/null 2>&1; then
        auto tools missing
    else
        echo "Automation framework not found"
        return 1
    fi
}

# Install a specific tool
install_tool() {
    local tool="$1"
    if [ -z "$tool" ]; then
        echo "Usage: install_tool <tool_name>"
        return 1
    fi

    if command -v auto >/dev/null 2>&1; then
        auto tools install "$tool"
    else
        echo "Automation framework not found"
        return 1
    fi
}

# =============================================================================
# Version Checking
# =============================================================================

# Quick version checks for common tools
quick_versions() {
    echo "=== Quick Version Check ==="

    # System tools
    echo "System Tools:"
    command -v git >/dev/null && echo "  git: $(git --version 2>/dev/null | head -1)" || echo "  git: Not installed"
    command -v curl >/dev/null && echo "  curl: $(curl --version 2>/dev/null | head -1)" || echo "  curl: Not installed"
    command -v jq >/dev/null && echo "  jq: $(jq --version 2>/dev/null)" || echo "  jq: Not installed"

    # Languages
    echo "Languages:"
    command -v go >/dev/null && echo "  go: $(go version 2>/dev/null)" || echo "  go: Not installed"
    command -v node >/dev/null && echo "  node: $(node --version 2>/dev/null)" || echo "  node: Not installed"
    command -v python3 >/dev/null && echo "  python3: $(python3 --version 2>/dev/null)" || echo "  python3: Not installed"

    # Cloud tools
    echo "Cloud Tools:"
    command -v aws >/dev/null && echo "  aws: $(aws --version 2>/dev/null | head -1)" || echo "  aws: Not installed"
    command -v kubectl >/dev/null && echo "  kubectl: $(kubectl version --client --short 2>/dev/null)" || echo "  kubectl: Not installed"

    # Development tools
    echo "Development Tools:"
    command -v docker >/dev/null && echo "  docker: $(docker --version 2>/dev/null)" || echo "  docker: Not installed"
    command -v gh >/dev/null && echo "  gh: $(gh --version 2>/dev/null | head -1)" || echo "  gh: Not installed"
}

# =============================================================================
# System Updates
# =============================================================================

# Smart package manager detection and update
smart_update() {
    echo "Smart system update..."

    if command -v brew >/dev/null 2>&1; then
        echo "Updating Homebrew packages..."
        brew update && brew upgrade
    elif command -v apt-get >/dev/null 2>&1; then
        echo "Updating apt packages..."
        sudo apt-get update && sudo apt-get upgrade
    elif command -v dnf >/dev/null 2>&1; then
        echo "Updating dnf packages..."
        sudo dnf upgrade
    elif command -v pacman >/dev/null 2>&1; then
        echo "Updating pacman packages..."
        sudo pacman -Syu
    else
        echo "No supported package manager found"
        return 1
    fi

    echo "Smart update completed!"
}

# =============================================================================
# Enhanced Utilities
# =============================================================================

# Enhanced which command with more info
which_enhanced() {
    local tool="$1"

    if [ -z "$tool" ]; then
        echo "Usage: which_enhanced <command>"
        return 1
    fi

    if command -v "$tool" >/dev/null 2>&1; then
        echo "$tool found:"
        echo "  Location: $(which "$tool")"

        # Try to get version
        if "$tool" --version >/dev/null 2>&1; then
            echo "  Version: $($tool --version 2>/dev/null | head -1)"
        elif "$tool" version >/dev/null 2>&1; then
            echo "  Version: $($tool version 2>/dev/null | head -1)"
        else
            echo "  Version: Unknown"
        fi
    else
        echo "$tool not found"
    fi
}
