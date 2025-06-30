#!/bin/bash
# External Tools Management Module
# Check for updates and upgrade external dependencies with user prompting

# Load framework
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "$SCRIPT_DIR/framework/core.sh"

# Tool definitions with metadata
declare -A TOOLS=(
    # Core System Tools
    ["brew"]="package_manager:homebrew"
    ["git"]="development:git"
    ["curl"]="system:network"
    ["wget"]="system:network"
    ["jq"]="system:json"
    
    # Shell Enhancement Tools
    ["fzf"]="shell:fuzzy_finder"
    ["fd"]="shell:file_search"
    ["bash-completion"]="shell:completion"
    
    # Development Languages
    ["go"]="language:golang"
    ["node"]="language:nodejs"
    ["npm"]="language:nodejs"
    ["python3"]="language:python"
    ["pip3"]="language:python"
    
    # Editors and IDEs
    ["nvim"]="editor:neovim"
    ["code"]="editor:vscode"
    ["tmux"]="terminal:multiplexer"
    
    # Cloud and Container Tools
    ["docker"]="container:docker"
    ["kubectl"]="kubernetes:cli"
    ["helm"]="kubernetes:package_manager"
    ["k9s"]="kubernetes:ui"
    
    # Cloud Provider CLIs
    ["aws"]="cloud:aws"
    ["az"]="cloud:azure"
    ["doctl"]="cloud:digitalocean"
    ["gh"]="development:github"
    
    # Language Specific Tools
    ["goimports"]="golang:formatter"
    ["golangci-lint"]="golang:linter"
    ["cobra-cli"]="golang:cli_framework"
    ["black"]="python:formatter"
    ["isort"]="python:import_sorter"
    ["flake8"]="python:linter"
    ["pytest"]="python:testing"
    ["uvicorn"]="python:server"
    ["ipython"]="python:interactive"
)

# Tool installation and update methods
declare -A INSTALL_METHODS=(
    ["brew"]="curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh | bash"
    ["git"]="brew:git|apt:git|yum:git"
    ["curl"]="brew:curl|apt:curl|yum:curl"
    ["wget"]="brew:wget|apt:wget|yum:wget"
    ["jq"]="brew:jq|apt:jq|yum:jq"
    ["fzf"]="brew:fzf|apt:fzf|yum:fzf"
    ["fd"]="brew:fd|apt:fd-find|yum:fd-find"
    ["bash-completion"]="brew:bash-completion|apt:bash-completion|yum:bash-completion"
    ["go"]="brew:go|manual:https://golang.org/dl/"
    ["node"]="brew:node|apt:nodejs|yum:nodejs"
    ["npm"]="brew:npm|apt:npm|yum:npm"
    ["python3"]="brew:python3|apt:python3|yum:python3"
    ["pip3"]="brew:python3|apt:python3-pip|yum:python3-pip"
    ["nvim"]="brew:neovim|apt:neovim|yum:neovim"
    ["code"]="brew:visual-studio-code|manual:https://code.visualstudio.com/"
    ["tmux"]="brew:tmux|apt:tmux|yum:tmux"
    ["docker"]="brew:docker|manual:https://docs.docker.com/engine/install/"
    ["kubectl"]="brew:kubectl|manual:kubernetes"
    ["helm"]="brew:helm|manual:helm"
    ["k9s"]="brew:k9s|manual:k9s"
    ["aws"]="brew:awscli|manual:aws"
    ["az"]="brew:azure-cli|manual:azure"
    ["doctl"]="brew:doctl|manual:doctl"
    ["gh"]="brew:gh|manual:gh"
    ["goimports"]="go:golang.org/x/tools/cmd/goimports@latest"
    ["golangci-lint"]="manual:golangci-lint"
    ["cobra-cli"]="go:github.com/spf13/cobra-cli@latest"
    ["black"]="pip:black"
    ["isort"]="pip:isort"
    ["flake8"]="pip:flake8"
    ["pytest"]="pip:pytest"
    ["uvicorn"]="pip:uvicorn"
    ["ipython"]="pip:ipython"
)

# Version check commands
declare -A VERSION_COMMANDS=(
    ["brew"]="brew --version"
    ["git"]="git --version"
    ["curl"]="curl --version"
    ["wget"]="wget --version"
    ["jq"]="jq --version"
    ["fzf"]="fzf --version"
    ["fd"]="fd --version"
    ["bash-completion"]="echo 'bash-completion (check manually)'"
    ["go"]="go version"
    ["node"]="node --version"
    ["npm"]="npm --version"
    ["python3"]="python3 --version"
    ["pip3"]="pip3 --version"
    ["nvim"]="nvim --version | head -1"
    ["code"]="code --version | head -1"
    ["tmux"]="tmux -V"
    ["docker"]="docker --version"
    ["kubectl"]="kubectl version --client --short"
    ["helm"]="helm version --short"
    ["k9s"]="k9s version --short"
    ["aws"]="aws --version"
    ["az"]="az --version | head -1"
    ["doctl"]="doctl version"
    ["gh"]="gh --version | head -1"
    ["goimports"]="goimports -help 2>&1 | head -1"
    ["golangci-lint"]="golangci-lint --version"
    ["cobra-cli"]="cobra-cli version"
    ["black"]="black --version"
    ["isort"]="isort --version"
    ["flake8"]="flake8 --version"
    ["pytest"]="pytest --version"
    ["uvicorn"]="uvicorn --version"
    ["ipython"]="ipython --version"
)

# Update commands
declare -A UPDATE_COMMANDS=(
    ["brew"]="brew update && brew upgrade"
    ["git"]="update_via_package_manager git"
    ["curl"]="update_via_package_manager curl"
    ["wget"]="update_via_package_manager wget"
    ["jq"]="update_via_package_manager jq"
    ["fzf"]="update_via_package_manager fzf"
    ["fd"]="update_via_package_manager fd"
    ["bash-completion"]="update_via_package_manager bash-completion"
    ["go"]="update_go_manual"
    ["node"]="update_via_package_manager node"
    ["npm"]="npm update -g"
    ["python3"]="update_via_package_manager python3"
    ["pip3"]="pip3 install --upgrade pip"
    ["nvim"]="update_via_package_manager neovim"
    ["code"]="update_vscode"
    ["tmux"]="update_via_package_manager tmux"
    ["docker"]="update_docker"
    ["kubectl"]="update_kubectl"
    ["helm"]="update_helm"
    ["k9s"]="update_via_package_manager k9s"
    ["aws"]="update_aws_cli"
    ["az"]="az upgrade"
    ["doctl"]="update_doctl"
    ["gh"]="update_via_package_manager gh"
    ["goimports"]="go install golang.org/x/tools/cmd/goimports@latest"
    ["golangci-lint"]="update_golangci_lint"
    ["cobra-cli"]="go install github.com/spf13/cobra-cli@latest"
    ["black"]="pip3 install --upgrade black"
    ["isort"]="pip3 install --upgrade isort"
    ["flake8"]="pip3 install --upgrade flake8"
    ["pytest"]="pip3 install --upgrade pytest"
    ["uvicorn"]="pip3 install --upgrade uvicorn"
    ["ipython"]="pip3 install --upgrade ipython"
)

tools_help() {
    cat << 'EOF'
External Tools Management

USAGE: auto tools <command> [options]

COMMANDS:
    list                     List all tools and their status
    check [tool]             Check version of specific tool or all tools
    update [tool]            Update specific tool or prompt for all
    install <tool>           Install a specific tool
    status                   Show comprehensive status report
    outdated                 Show tools that may have updates available
    categories               List tools by category
    missing                  Show tools that are not installed

CATEGORIES:
    system                   Core system utilities
    shell                    Shell enhancement tools
    development              Development tools and languages
    cloud                    Cloud provider CLIs
    kubernetes               Kubernetes tools
    container                Container and Docker tools
    editor                   Editors and IDEs

OPTIONS:
    --force                  Skip confirmation prompts
    --yes-to-all             Answer yes to all prompts
    --category <cat>         Filter by category
    --dry-run                Show what would be done
    --verbose                Show detailed output

EXAMPLES:
    auto tools list                     # List all tools
    auto tools check                    # Check all tool versions
    auto tools check kubectl           # Check kubectl version
    auto tools update                   # Interactive update all tools
    auto tools update --force          # Update all without prompting
    auto tools update --yes-to-all     # Update all with auto-yes prompts
    auto tools install go              # Install Go language
    auto tools status                   # Comprehensive status report
    auto tools outdated                 # Show potentially outdated tools
    auto tools list --category cloud   # List only cloud tools
EOF
}

# Detect package manager
detect_package_manager() {
    if command -v brew >/dev/null 2>&1; then
        echo "brew"
    elif command -v apt-get >/dev/null 2>&1; then
        echo "apt"
    elif command -v yum >/dev/null 2>&1; then
        echo "yum"
    elif command -v dnf >/dev/null 2>&1; then
        echo "dnf"
    elif command -v pacman >/dev/null 2>&1; then
        echo "pacman"
    else
        echo "none"
    fi
}

# Check if tool is installed
is_tool_installed() {
    local tool="$1"
    command -v "$tool" >/dev/null 2>&1
}

# Get tool version
get_tool_version() {
    local tool="$1"
    
    if ! is_tool_installed "$tool"; then
        echo "Not installed"
        return 1
    fi
    
    local version_cmd="${VERSION_COMMANDS[$tool]}"
    if [ -n "$version_cmd" ]; then
        eval "$version_cmd" 2>/dev/null | head -1 || echo "Version unknown"
    else
        echo "Version check not available"
    fi
}

# Get tool category
get_tool_category() {
    local tool="$1"
    local metadata="${TOOLS[$tool]}"
    echo "${metadata%%:*}"
}

# Get tool subcategory
get_tool_subcategory() {
    local tool="$1"
    local metadata="${TOOLS[$tool]}"
    echo "${metadata##*:}"
}

# Update via package manager
update_via_package_manager() {
    local tool="$1"
    local pkg_mgr=$(detect_package_manager)
    
    case "$pkg_mgr" in
        "brew")
            log_info "Updating $tool via Homebrew..."
            brew upgrade "$tool" 2>/dev/null || brew install "$tool"
            ;;
        "apt")
            log_info "Updating $tool via apt..."
            sudo apt-get update >/dev/null 2>&1
            sudo apt-get install -y "$tool"
            ;;
        "yum")
            log_info "Updating $tool via yum..."
            sudo yum update -y "$tool" 2>/dev/null || sudo yum install -y "$tool"
            ;;
        "dnf")
            log_info "Updating $tool via dnf..."
            sudo dnf upgrade -y "$tool" 2>/dev/null || sudo dnf install -y "$tool"
            ;;
        "pacman")
            log_info "Updating $tool via pacman..."
            sudo pacman -S --noconfirm "$tool"
            ;;
        *)
            log_warn "No package manager found for updating $tool"
            return 1
            ;;
    esac
}

# Specialized update functions
update_go_manual() {
    log_info "Go requires manual update. Please visit: https://golang.org/dl/"
    log_info "Current version: $(go version 2>/dev/null || echo 'Not installed')"
}

update_vscode() {
    if command -v code >/dev/null 2>&1; then
        log_info "VS Code can be updated through its built-in updater"
        log_info "Or via package manager if installed that way"
        update_via_package_manager "visual-studio-code"
    else
        log_warn "VS Code not found"
    fi
}

update_docker() {
    if command -v docker >/dev/null 2>&1; then
        if [[ "$OSTYPE" == "darwin"* ]]; then
            log_info "Docker Desktop can be updated through Docker Desktop or Homebrew"
            update_via_package_manager "docker"
        else
            log_info "Docker can be updated via package manager"
            log_info "Visit: https://docs.docker.com/engine/install/ for instructions"
        fi
    else
        log_warn "Docker not found"
    fi
}

update_kubectl() {
    if command -v kubectl >/dev/null 2>&1; then
        local pkg_mgr=$(detect_package_manager)
        if [ "$pkg_mgr" = "brew" ]; then
            brew upgrade kubectl
        else
            log_info "Updating kubectl manually..."
            local version=$(curl -L -s https://dl.k8s.io/release/stable.txt)
            curl -LO "https://dl.k8s.io/release/$version/bin/linux/amd64/kubectl"
            chmod +x kubectl
            sudo mv kubectl /usr/local/bin/
        fi
    else
        log_warn "kubectl not found"
    fi
}

update_helm() {
    if command -v helm >/dev/null 2>&1; then
        local pkg_mgr=$(detect_package_manager)
        if [ "$pkg_mgr" = "brew" ]; then
            brew upgrade helm
        else
            log_info "Updating Helm via installer script..."
            curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
        fi
    else
        log_warn "Helm not found"
    fi
}

update_aws_cli() {
    if command -v aws >/dev/null 2>&1; then
        local pkg_mgr=$(detect_package_manager)
        if [ "$pkg_mgr" = "brew" ]; then
            brew upgrade awscli
        else
            log_info "AWS CLI manual update required"
            log_info "Visit: https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html"
        fi
    else
        log_warn "AWS CLI not found"
    fi
}

update_doctl() {
    if command -v doctl >/dev/null 2>&1; then
        local pkg_mgr=$(detect_package_manager)
        if [ "$pkg_mgr" = "brew" ]; then
            brew upgrade doctl
        else
            log_info "doctl manual update required"
            log_info "Visit: https://docs.digitalocean.com/reference/doctl/how-to/install/"
        fi
    else
        log_warn "doctl not found"
    fi
}

update_golangci_lint() {
    if command -v golangci-lint >/dev/null 2>&1; then
        log_info "Updating golangci-lint..."
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
    else
        log_warn "golangci-lint not found"
    fi
}

# List all tools
tools_list() {
    local category_filter=""
    
    # Parse options
    while [[ $# -gt 0 ]]; do
        case $1 in
            --category)
                category_filter="$2"
                shift 2
                ;;
            *)
                shift
                ;;
        esac
    done
    
    log_info "External Tools Inventory"
    echo
    
    local current_category=""
    
    for tool in $(printf '%s\n' "${!TOOLS[@]}" | sort); do
        local category=$(get_tool_category "$tool")
        local subcategory=$(get_tool_subcategory "$tool")
        
        # Skip if category filter doesn't match
        if [ -n "$category_filter" ] && [ "$category" != "$category_filter" ]; then
            continue
        fi
        
        # Print category header
        if [ "$category" != "$current_category" ]; then
            echo "=== ${category^^} TOOLS ==="
            current_category="$category"
        fi
        
        local status="❌"
        local version="Not installed"
        
        if is_tool_installed "$tool"; then
            status="✅"
            version=$(get_tool_version "$tool")
        fi
        
        printf "  %-20s %s %-15s %s\n" "$tool" "$status" "($subcategory)" "$version"
    done
    
    echo
}

# Check tool versions
tools_check() {
    local specific_tool="$1"
    
    if [ -n "$specific_tool" ]; then
        if [[ ! "${!TOOLS[@]}" =~ $specific_tool ]]; then
            log_error "Unknown tool: $specific_tool"
            return 1
        fi
        
        log_info "Checking $specific_tool..."
        if is_tool_installed "$specific_tool"; then
            local version=$(get_tool_version "$specific_tool")
            log_success "✅ $specific_tool: $version"
        else
            log_warn "❌ $specific_tool: Not installed"
        fi
    else
        log_info "Checking all tool versions..."
        echo
        
        local installed=0
        local missing=0
        
        for tool in $(printf '%s\n' "${!TOOLS[@]}" | sort); do
            if is_tool_installed "$tool"; then
                local version=$(get_tool_version "$tool")
                echo "✅ $tool: $version"
                ((installed++))
            else
                echo "❌ $tool: Not installed"
                ((missing++))
            fi
        done
        
        echo
        log_info "Summary: $installed installed, $missing missing"
    fi
}

# Update tools with user prompting
tools_update() {
    local specific_tool=""
    local force_update=false
    
    # Parse options
    while [[ $# -gt 0 ]]; do
        case $1 in
            --force)
                force_update=true
                shift
                ;;
            --yes-to-all)
                export YES_TO_ALL=true
                shift
                ;;
            *)
                if [ -z "$specific_tool" ]; then
                    specific_tool="$1"
                fi
                shift
                ;;
        esac
    done
    
    if [ -n "$specific_tool" ]; then
        if [[ ! "${!TOOLS[@]}" =~ $specific_tool ]]; then
            log_error "Unknown tool: $specific_tool"
            return 1
        fi
        
        update_single_tool "$specific_tool" "$force_update"
    else
        update_all_tools "$force_update"
    fi
}

# Update single tool
update_single_tool() {
    local tool="$1"
    local force="$2"
    
    if ! is_tool_installed "$tool"; then
        log_warn "$tool is not installed"
        if [ "$force" = true ] || prompt_user "Install $tool?"; then
            install_tool "$tool"
        fi
        return
    fi
    
    local current_version=$(get_tool_version "$tool")
    log_info "Current $tool version: $current_version"
    
    if [ "$force" = true ] || prompt_user "Update $tool?"; then
        local update_cmd="${UPDATE_COMMANDS[$tool]}"
        if [ -n "$update_cmd" ]; then
            log_info "Updating $tool..."
            eval "$update_cmd"
            
            # Check new version
            local new_version=$(get_tool_version "$tool")
            if [ "$new_version" != "$current_version" ]; then
                log_success "✅ $tool updated: $current_version → $new_version"
            else
                log_info "ℹ️ $tool already up to date: $new_version"
            fi
        else
            log_warn "No update method defined for $tool"
        fi
    else
        log_info "Skipping $tool update"
    fi
}

# Update all tools
update_all_tools() {
    local force="$1"
    
    log_info "Checking all tools for updates..."
    echo
    
    local tools_to_update=()
    
    # First pass: identify tools that need updates
    for tool in $(printf '%s\n' "${!TOOLS[@]}" | sort); do
        if is_tool_installed "$tool"; then
            tools_to_update+=("$tool")
        else
            log_warn "❌ $tool: Not installed"
        fi
    done
    
    if [ ${#tools_to_update[@]} -eq 0 ]; then
        log_warn "No tools available for update"
        return
    fi
    
    log_info "Found ${#tools_to_update[@]} tools to check for updates"
    echo
    
    # Second pass: update each tool with prompting
    for tool in "${tools_to_update[@]}"; do
        echo "--- Updating $tool ---"
        update_single_tool "$tool" "$force"
        echo
    done
    
    log_success "Tool update process completed!"
}

# Install a specific tool
install_tool() {
    local tool="$1"
    
    if [[ ! "${!TOOLS[@]}" =~ $tool ]]; then
        log_error "Unknown tool: $tool"
        return 1
    fi
    
    if is_tool_installed "$tool"; then
        log_warn "$tool is already installed"
        return 0
    fi
    
    local install_method="${INSTALL_METHODS[$tool]}"
    if [ -z "$install_method" ]; then
        log_error "No installation method defined for $tool"
        return 1
    fi
    
    log_info "Installing $tool..."
    
    # Parse installation method
    if [[ "$install_method" =~ ^brew: ]]; then
        local package="${install_method#brew:}"
        brew install "$package"
    elif [[ "$install_method" =~ ^apt: ]]; then
        local package="${install_method#apt:}"
        sudo apt-get update && sudo apt-get install -y "$package"
    elif [[ "$install_method" =~ ^yum: ]]; then
        local package="${install_method#yum:}"
        sudo yum install -y "$package"
    elif [[ "$install_method" =~ ^pip: ]]; then
        local package="${install_method#pip:}"
        pip3 install "$package"
    elif [[ "$install_method" =~ ^go: ]]; then
        local package="${install_method#go:}"
        go install "$package"
    elif [[ "$install_method" =~ ^manual: ]]; then
        local info="${install_method#manual:}"
        log_info "Manual installation required for $tool"
        log_info "Please visit: $info"
    else
        log_info "Running: $install_method"
        eval "$install_method"
    fi
    
    # Verify installation
    if is_tool_installed "$tool"; then
        local version=$(get_tool_version "$tool")
        log_success "✅ $tool installed successfully: $version"
    else
        log_error "❌ $tool installation failed"
        return 1
    fi
}

# Show comprehensive status
tools_status() {
    log_info "External Tools Status Report"
    echo "Generated: $(date)"
    echo
    
    # System info
    echo "=== SYSTEM INFORMATION ==="
    echo "OS: $(uname -s) $(uname -r)"
    echo "Package Manager: $(detect_package_manager)"
    echo "Shell: $SHELL"
    echo
    
    # Category breakdown
    local categories=($(printf '%s\n' "${TOOLS[@]}" | cut -d: -f1 | sort -u))
    
    for category in "${categories[@]}"; do
        echo "=== ${category^^} TOOLS ==="
        
        local installed=0
        local missing=0
        
        for tool in $(printf '%s\n' "${!TOOLS[@]}" | sort); do
            if [ "$(get_tool_category "$tool")" = "$category" ]; then
                if is_tool_installed "$tool"; then
                    local version=$(get_tool_version "$tool")
                    echo "  ✅ $tool: $version"
                    ((installed++))
                else
                    echo "  ❌ $tool: Not installed"
                    ((missing++))
                fi
            fi
        done
        
        echo "  Total: $installed installed, $missing missing"
        echo
    done
    
    # Overall summary
    local total_installed=0
    local total_missing=0
    
    for tool in "${!TOOLS[@]}"; do
        if is_tool_installed "$tool"; then
            ((total_installed++))
        else
            ((total_missing++))
        fi
    done
    
    echo "=== OVERALL SUMMARY ==="
    echo "Total tools: $((total_installed + total_missing))"
    echo "Installed: $total_installed"
    echo "Missing: $total_missing"
    echo "Coverage: $((total_installed * 100 / (total_installed + total_missing)))%"
}

# Show missing tools
tools_missing() {
    log_info "Missing Tools"
    echo
    
    local missing_tools=()
    
    for tool in $(printf '%s\n' "${!TOOLS[@]}" | sort); do
        if ! is_tool_installed "$tool"; then
            missing_tools+=("$tool")
        fi
    done
    
    if [ ${#missing_tools[@]} -eq 0 ]; then
        log_success "🎉 All tools are installed!"
        return 0
    fi
    
    echo "The following tools are not installed:"
    echo
    
    for tool in "${missing_tools[@]}"; do
        local category=$(get_tool_category "$tool")
        local subcategory=$(get_tool_subcategory "$tool")
        echo "  ❌ $tool ($category:$subcategory)"
    done
    
    echo
    log_info "To install missing tools:"
    echo "  auto tools install <tool_name>"
    echo "  auto tools update --force    # Install all missing tools"
}

# Show categories
tools_categories() {
    log_info "Tool Categories"
    echo
    
    local categories=($(printf '%s\n' "${TOOLS[@]}" | cut -d: -f1 | sort -u))
    
    for category in "${categories[@]}"; do
        echo "=== ${category^^} ==="
        
        for tool in $(printf '%s\n' "${!TOOLS[@]}" | sort); do
            if [ "$(get_tool_category "$tool")" = "$category" ]; then
                local subcategory=$(get_tool_subcategory "$tool")
                local status="❌"
                [ is_tool_installed "$tool" ] && status="✅"
                echo "  $status $tool ($subcategory)"
            fi
        done
        echo
    done
}

# Show potentially outdated tools
tools_outdated() {
    log_info "Checking for potentially outdated tools..."
    echo
    
    # This is a simple implementation - could be enhanced with actual version comparison
    log_info "Note: This shows installed tools that might have updates available."
    log_info "Run 'auto tools update' to check and update individual tools."
    echo
    
    for tool in $(printf '%s\n' "${!TOOLS[@]}" | sort); do
        if is_tool_installed "$tool"; then
            local version=$(get_tool_version "$tool")
            echo "🔍 $tool: $version"
        fi
    done
    
    echo
    log_info "To update all tools: auto tools update"
    log_info "To update specific tool: auto tools update <tool_name>"
}

# Utility function for user prompts
prompt_user() {
    local message="$1"
    local response
    
    # Check for YES_TO_ALL environment variable
    if [ "${YES_TO_ALL:-false}" = "true" ]; then
        printf "%s [y/n]: y (auto-yes)\n" "$message"
        return 0
    fi
    
    while true; do
        printf "%s [y/n]: " "$message"
        read -r response
        case "$response" in
            y|Y|yes|Yes|YES) return 0 ;;
            n|N|no|No|NO) return 1 ;;
            *) echo "Please enter y or n" ;;
        esac
    done
}

# Main tools function
tools_main() {
    case "${1:-help}" in
        "list") shift; tools_list "$@" ;;
        "check") shift; tools_check "$@" ;;
        "update") shift; tools_update "$@" ;;
        "install") shift; install_tool "$@" ;;
        "status") tools_status ;;
        "missing") tools_missing ;;
        "categories") tools_categories ;;
        "outdated") tools_outdated ;;
        "help"|*) tools_help ;;
    esac
}