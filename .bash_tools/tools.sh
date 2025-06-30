#!/bin/bash
# Tools Management Integration for Shell
# Quick access to external tools management

# Quick tool status check
tools_status() {
    if command -v auto >/dev/null 2>&1; then
        auto tools status
    else
        echo "❌ Automation framework not found"
        return 1
    fi
}

# List all tools
list_tools() {
    if command -v auto >/dev/null 2>&1; then
        auto tools list
    else
        echo "❌ Automation framework not found"
        return 1
    fi
}

# Check tool versions
check_tools() {
    if command -v auto >/dev/null 2>&1; then
        auto tools check "$@"
    else
        echo "❌ Automation framework not found"
        return 1
    fi
}

# Update tools interactively
update_tools() {
    if command -v auto >/dev/null 2>&1; then
        auto tools update "$@"
    else
        echo "❌ Automation framework not found"
        return 1
    fi
}

# Show missing tools
missing_tools() {
    if command -v auto >/dev/null 2>&1; then
        auto tools missing
    else
        echo "❌ Automation framework not found"
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
        echo "❌ Automation framework not found"
        return 1
    fi
}

# Show tools by category
tools_by_category() {
    if command -v auto >/dev/null 2>&1; then
        auto tools categories
    else
        echo "❌ Automation framework not found"
        return 1
    fi
}

# Check for outdated tools
outdated_tools() {
    if command -v auto >/dev/null 2>&1; then
        auto tools outdated
    else
        echo "❌ Automation framework not found"
        return 1
    fi
}

# Quick version checks for common tools
quick_versions() {
    echo "=== Quick Version Check ==="
    
    # System tools
    echo "System Tools:"
    command -v git >/dev/null && echo "  ✅ git: $(git --version 2>/dev/null | head -1)" || echo "  ❌ git: Not installed"
    command -v curl >/dev/null && echo "  ✅ curl: $(curl --version 2>/dev/null | head -1)" || echo "  ❌ curl: Not installed"
    command -v jq >/dev/null && echo "  ✅ jq: $(jq --version 2>/dev/null)" || echo "  ❌ jq: Not installed"
    
    # Languages
    echo "Languages:"
    command -v go >/dev/null && echo "  ✅ go: $(go version 2>/dev/null)" || echo "  ❌ go: Not installed"
    command -v node >/dev/null && echo "  ✅ node: $(node --version 2>/dev/null)" || echo "  ❌ node: Not installed"
    command -v python3 >/dev/null && echo "  ✅ python3: $(python3 --version 2>/dev/null)" || echo "  ❌ python3: Not installed"
    
    # Cloud tools
    echo "Cloud Tools:"
    command -v aws >/dev/null && echo "  ✅ aws: $(aws --version 2>/dev/null | head -1)" || echo "  ❌ aws: Not installed"
    command -v az >/dev/null && echo "  ✅ az: $(az --version 2>/dev/null | head -1)" || echo "  ❌ az: Not installed"
    command -v doctl >/dev/null && echo "  ✅ doctl: $(doctl version 2>/dev/null)" || echo "  ❌ doctl: Not installed"
    command -v kubectl >/dev/null && echo "  ✅ kubectl: $(kubectl version --client --short 2>/dev/null)" || echo "  ❌ kubectl: Not installed"
    
    # Development tools
    echo "Development Tools:"
    command -v docker >/dev/null && echo "  ✅ docker: $(docker --version 2>/dev/null)" || echo "  ❌ docker: Not installed"
    command -v gh >/dev/null && echo "  ✅ gh: $(gh --version 2>/dev/null | head -1)" || echo "  ❌ gh: Not installed"
    command -v code >/dev/null && echo "  ✅ code: $(code --version 2>/dev/null | head -1)" || echo "  ❌ code: Not installed"
}

# Smart package manager detection and update
smart_update() {
    echo "🔄 Smart system update..."
    
    if command -v brew >/dev/null 2>&1; then
        echo "📦 Updating Homebrew packages..."
        brew update && brew upgrade
    elif command -v apt-get >/dev/null 2>&1; then
        echo "📦 Updating apt packages..."
        sudo apt-get update && sudo apt-get upgrade
    elif command -v yum >/dev/null 2>&1; then
        echo "📦 Updating yum packages..."
        sudo yum update
    elif command -v dnf >/dev/null 2>&1; then
        echo "📦 Updating dnf packages..."
        sudo dnf upgrade
    elif command -v pacman >/dev/null 2>&1; then
        echo "📦 Updating pacman packages..."
        sudo pacman -Syu
    else
        echo "❌ No supported package manager found"
        return 1
    fi
    
    # Update language package managers
    echo "🐍 Updating Python packages..."
    command -v pip3 >/dev/null && pip3 install --upgrade pip
    
    echo "📦 Updating npm packages..."
    command -v npm >/dev/null && npm update -g
    
    echo "🐹 Updating Go tools..."
    if command -v go >/dev/null; then
        go install golang.org/x/tools/cmd/goimports@latest
        go install github.com/spf13/cobra-cli@latest
    fi
    
    echo "✅ Smart update completed!"
}

# Show tools that need manual installation
manual_install_tools() {
    echo "🔧 Tools requiring manual installation:"
    echo
    
    # Check for tools that typically need manual installation on Linux
    if [[ "$OSTYPE" == "linux-gnu" ]]; then
        echo "Linux Manual Installations:"
        
        if ! command -v go >/dev/null; then
            echo "  ❌ Go: https://golang.org/dl/"
        fi
        
        if ! command -v code >/dev/null; then
            echo "  ❌ VS Code: https://code.visualstudio.com/"
        fi
        
        if ! command -v docker >/dev/null; then
            echo "  ❌ Docker: https://docs.docker.com/engine/install/"
        fi
        
        if ! command -v gh >/dev/null; then
            echo "  ❌ GitHub CLI: https://cli.github.com/"
        fi
        
        if ! command -v helm >/dev/null; then
            echo "  ❌ Helm: https://helm.sh/docs/intro/install/"
        fi
        
        if ! command -v k9s >/dev/null; then
            echo "  ❌ k9s: https://k9scli.io/topics/install/"
        fi
    fi
    
    # Cloud CLIs that often need manual setup
    echo "Cloud CLI Manual Setup:"
    
    if ! command -v aws >/dev/null; then
        echo "  ❌ AWS CLI: https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html"
    fi
    
    if ! command -v az >/dev/null; then
        echo "  ❌ Azure CLI: https://docs.microsoft.com/en-us/cli/azure/install-azure-cli"
    fi
    
    if ! command -v doctl >/dev/null; then
        echo "  ❌ doctl: https://docs.digitalocean.com/reference/doctl/how-to/install/"
    fi
}

# Enhanced which command with more info
which_enhanced() {
    local tool="$1"
    
    if [ -z "$tool" ]; then
        echo "Usage: which_enhanced <command>"
        return 1
    fi
    
    if command -v "$tool" >/dev/null 2>&1; then
        echo "✅ $tool found:"
        echo "  Location: $(which "$tool")"
        
        # Try to get version
        if "$tool" --version >/dev/null 2>&1; then
            echo "  Version: $($tool --version 2>/dev/null | head -1)"
        elif "$tool" version >/dev/null 2>&1; then
            echo "  Version: $($tool version 2>/dev/null | head -1)"
        elif "$tool" -V >/dev/null 2>&1; then
            echo "  Version: $($tool -V 2>/dev/null | head -1)"
        else
            echo "  Version: Unknown"
        fi
        
        # Show file info
        if [ -f "$(which "$tool")" ]; then
            echo "  Type: $(file "$(which "$tool")" | cut -d: -f2-)"
            echo "  Size: $(du -h "$(which "$tool")" | cut -f1)"
            echo "  Modified: $(stat -c %y "$(which "$tool")" 2>/dev/null || stat -f %Sm "$(which "$tool")" 2>/dev/null)"
        fi
    else
        echo "❌ $tool not found"
        
        # Suggest installation
        echo "  Install with automation: auto tools install $tool"
        echo "  Or check manual installation: manual_install_tools"
    fi
}

# Aliases for convenience
alias tools='tools_status'
alias tools-list='list_tools'
alias tools-check='check_tools'
alias tools-update='update_tools'
alias tools-missing='missing_tools'
alias tools-install='install_tool'
alias tools-categories='tools_by_category'
alias tools-outdated='outdated_tools'
alias versions='quick_versions'
alias system-update='smart_update'
alias manual-tools='manual_install_tools'
alias whichx='which_enhanced'

# Help function
tools_help() {
    cat << 'EOF'
Tools Management Quick Commands

OVERVIEW COMMANDS:
    tools                   Show comprehensive tools status
    tools-list              List all tools with status
    tools-check [tool]      Check specific tool or all tool versions
    versions                Quick version check for common tools

UPDATE COMMANDS:
    tools-update [tool]     Interactive update for specific or all tools
    system-update           Smart package manager update
    tools-outdated          Show tools that may need updates

INSTALLATION:
    tools-install <tool>    Install specific tool
    tools-missing           Show tools that are not installed
    manual-tools            Show tools requiring manual installation

CATEGORIES:
    tools-categories        Show tools organized by category

UTILITIES:
    whichx <command>        Enhanced which with version and file info

AUTOMATION FRAMEWORK:
    auto tools status       Full tools status report
    auto tools update       Interactive tool updates
    auto tools list         Detailed tool listing
    auto tools install      Install tools via automation framework

EXAMPLES:
    tools                   # Show status overview
    tools-check kubectl    # Check kubectl version
    tools-update go         # Update Go language
    tools-install docker   # Install Docker
    versions               # Quick version check
    system-update          # Update all system packages
    whichx python3         # Enhanced info about python3
EOF
}

alias tools-help='tools_help'