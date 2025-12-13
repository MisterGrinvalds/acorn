#!/bin/bash
# Enhanced Bash Profile & Automation Framework Installer
# Supports dotfiles, automation framework, secrets management, and development tools
# Now with XDG-compliant structure using shell/init.sh as main entry point

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Global variables
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKUP_DIR="$HOME/.dotfiles-backup-$(date +%Y%m%d_%H%M%S)"
DOTFILES_INSTALLED=false
AUTOMATION_INSTALLED=false
CLOUD_TOOLS_INSTALLED=false
RESPONSE=""
YES_TO_ALL=false

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $*"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $*"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $*"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $*"
}

# Portable function to get user response
getResponse() {
    local OPTIND MESSAGE
    while getopts ":m:" opt; do
        case "$opt" in
            m )
                MESSAGE="$OPTARG"
                ;;
        esac
    done
    shift $((OPTIND-1))

    # If YES_TO_ALL is set, automatically answer yes
    if [ "$YES_TO_ALL" = true ]; then
        printf "%s [y/n]: y (auto-yes)\n" "$MESSAGE"
        RESPONSE="y"
        return 0
    fi

    # Loop until valid response
    while true; do
        printf "%s [y/n]: " "$MESSAGE"
        read -r RESPONSE
        case "$RESPONSE" in
            y|Y|yes|Yes|YES) RESPONSE="y"; break ;;
            n|N|no|No|NO) RESPONSE="n"; break ;;
            *) echo "Please enter y or n" ;;
        esac
    done
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Install VS Code if needed
install_vscode_if_needed() {
    if [ "${SKIP_GUI_APPS:-false}" = "true" ]; then
        log_info "Skipping VS Code installation (--skip-gui flag set)"
        return 0
    fi

    if [ -d "/Applications/Visual Studio Code.app" ]; then
        log_info "VS Code is already installed"
        if [ "${NON_INTERACTIVE:-false}" = "true" ]; then
            log_info "Skipping VS Code reinstallation in non-interactive mode"
        elif [ "${YES_TO_ALL:-false}" = "true" ]; then
            log_info "Reinstalling VS Code (yes-to-all mode)..."
            brew uninstall --cask visual-studio-code 2>/dev/null || true
            if brew install --cask visual-studio-code; then
                log_success "VS Code reinstalled successfully"
            else
                log_warn "VS Code reinstallation failed - may already be installed"
            fi
        else
            getResponse -m "Reinstall Visual Studio Code?"
            if [ "$RESPONSE" = 'y' ]; then
                log_info "Reinstalling VS Code..."
                brew uninstall --cask visual-studio-code 2>/dev/null || true
                if brew install --cask visual-studio-code; then
                    log_success "VS Code reinstalled successfully"
                else
                    log_warn "VS Code reinstallation failed - may already be installed"
                fi
            fi
        fi
    else
        if [ "${NON_INTERACTIVE:-false}" = "true" ] || [ "${YES_TO_ALL:-false}" = "true" ]; then
            log_info "Installing VS Code..."
            if brew install --cask visual-studio-code; then
                log_success "VS Code installed successfully"
            else
                log_warn "VS Code installation failed"
            fi
        else
            getResponse -m "Install Visual Studio Code?"
            if [ "$RESPONSE" = 'y' ]; then
                log_info "Installing VS Code..."
                if brew install --cask visual-studio-code; then
                    log_success "VS Code installed successfully"
                else
                    log_warn "VS Code installation failed"
                fi
            fi
        fi
    fi
}

# Install Docker Desktop if needed
install_docker_if_needed() {
    if [ "${SKIP_GUI_APPS:-false}" = "true" ]; then
        log_info "Skipping Docker Desktop installation (--skip-gui flag set)"
        return 0
    fi

    if [ -d "/Applications/Docker.app" ]; then
        log_info "Docker Desktop is already installed"
        if [ "${NON_INTERACTIVE:-false}" = "true" ]; then
            log_info "Skipping Docker Desktop reinstallation in non-interactive mode"
        elif [ "${YES_TO_ALL:-false}" = "true" ]; then
            log_info "Reinstalling Docker Desktop (yes-to-all mode)..."
            brew uninstall --cask docker 2>/dev/null || true
            if brew install --cask docker; then
                log_success "Docker Desktop reinstalled successfully"
            else
                log_warn "Docker Desktop reinstallation failed - may already be installed"
            fi
        else
            getResponse -m "Reinstall Docker Desktop?"
            if [ "$RESPONSE" = 'y' ]; then
                log_info "Reinstalling Docker Desktop..."
                brew uninstall --cask docker 2>/dev/null || true
                if brew install --cask docker; then
                    log_success "Docker Desktop reinstalled successfully"
                else
                    log_warn "Docker Desktop reinstallation failed - may already be installed"
                fi
            fi
        fi
    else
        if [ "${NON_INTERACTIVE:-false}" = "true" ] || [ "${YES_TO_ALL:-false}" = "true" ]; then
            log_info "Installing Docker Desktop..."
            if brew install --cask docker; then
                log_success "Docker Desktop installed successfully"
            else
                log_warn "Docker Desktop installation failed"
            fi
        else
            getResponse -m "Install Docker Desktop?"
            if [ "$RESPONSE" = 'y' ]; then
                log_info "Installing Docker Desktop..."
                if brew install --cask docker; then
                    log_success "Docker Desktop installed successfully"
                else
                    log_warn "Docker Desktop installation failed"
                fi
            fi
        fi
    fi
}

# Create backup of existing files
create_backup() {
    log_info "Creating backup of existing configurations..."
    mkdir -p "$BACKUP_DIR"

    # Backup existing dotfiles
    for file in .bash_profile .bashrc .zshrc .bash_profile.dir .bash_tools; do
        if [ -e "$HOME/$file" ]; then
            cp -r "$HOME/$file" "$BACKUP_DIR/" 2>/dev/null || true
            log_info "Backed up ~/$file"
        fi
    done

    # Backup XDG shell config if exists
    if [ -d "$HOME/.config/shell" ]; then
        cp -r "$HOME/.config/shell" "$BACKUP_DIR/" 2>/dev/null || true
        log_info "Backed up ~/.config/shell"
    fi

    log_success "Backup created at: $BACKUP_DIR"
}

# Install dotfiles (new XDG-compliant method)
install_dotfiles() {
    log_info "Installing dotfiles with XDG-compliant structure..."

    # Check if supported OS
    if [[ ! "$OSTYPE" =~ ^(darwin|linux-gnu) ]]; then
        log_error "Unsupported OS: $OSTYPE"
        exit 1
    fi

    # Set DOTFILES_ROOT to this repository location
    local dotfiles_root="$SCRIPT_DIR"

    # Create XDG directories if needed
    mkdir -p "$HOME/.config"
    mkdir -p "$HOME/.local/share"
    mkdir -p "$HOME/.local/state"
    mkdir -p "$HOME/.cache"

    # Create bootstrap file for bash
    log_info "Creating ~/.bashrc bootstrap..."
    cat > "$HOME/.bashrc" << EOF
# Dotfiles bootstrap - sources the main init.sh
# Generated by install.sh on $(date)
export DOTFILES_ROOT="$dotfiles_root"
[ -f "\$DOTFILES_ROOT/shell/init.sh" ] && . "\$DOTFILES_ROOT/shell/init.sh"
EOF
    log_success "Created ~/.bashrc"

    # Create bootstrap file for zsh
    log_info "Creating ~/.zshrc bootstrap..."
    cat > "$HOME/.zshrc" << EOF
# Dotfiles bootstrap - sources the main init.sh
# Generated by install.sh on $(date)
export DOTFILES_ROOT="$dotfiles_root"
[ -f "\$DOTFILES_ROOT/shell/init.sh" ] && . "\$DOTFILES_ROOT/shell/init.sh"
EOF
    log_success "Created ~/.zshrc"

    # Create ~/.bash_profile for login shells (sources .bashrc)
    log_info "Creating ~/.bash_profile..."
    cat > "$HOME/.bash_profile" << EOF
# Login shell - source .bashrc for interactive settings
# Generated by install.sh on $(date)
[ -f "\$HOME/.bashrc" ] && . "\$HOME/.bashrc"
EOF
    log_success "Created ~/.bash_profile"

    # Create ~/.zprofile for zsh login shells
    log_info "Creating ~/.zprofile..."
    cat > "$HOME/.zprofile" << EOF
# Zsh login shell - source .zshrc for interactive settings
# Generated by install.sh on $(date)
[ -f "\$HOME/.zshrc" ] && . "\$HOME/.zshrc"
EOF
    log_success "Created ~/.zprofile"

    DOTFILES_INSTALLED=true
    log_success "Dotfiles installed successfully!"
    log_info "Repository location: $dotfiles_root"
}

# Install automation framework
install_automation_framework() {
    log_info "Setting up automation framework..."

    # Make automation scripts executable
    if [ -d "$SCRIPT_DIR/.automation" ]; then
        chmod +x "$SCRIPT_DIR/.automation/auto"
        chmod +x "$SCRIPT_DIR/.automation/setup.sh" 2>/dev/null || true
        find "$SCRIPT_DIR/.automation/framework" -name "*.sh" -exec chmod +x {} \; 2>/dev/null || true
        find "$SCRIPT_DIR/.automation/modules" -name "*.sh" -exec chmod +x {} \; 2>/dev/null || true

        # Run automation setup if it exists
        if [ -f "$SCRIPT_DIR/.automation/setup.sh" ]; then
            "$SCRIPT_DIR/.automation/setup.sh"
        fi

        AUTOMATION_INSTALLED=true
        log_success "Automation framework installed!"
    else
        log_warn "Automation directory not found"
    fi
}

# Install package manager tools
install_package_manager_tools() {
    log_info "Installing package manager tools..."

    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS - Install Homebrew if not present
        if ! command_exists brew; then
            log_info "Installing Homebrew..."
            /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        fi

        if command_exists brew; then
            log_info "Updating Homebrew..."
            brew update

            # Core shell tools
            log_info "Installing core shell tools..."
            brew install bash-completion fd fzf xclip jq curl wget
            brew install --cask xquartz 2>/dev/null || log_warn "XQuartz install skipped"

            return 0
        else
            log_error "Homebrew installation failed"
            return 1
        fi

    elif [[ "$OSTYPE" == "linux-gnu" ]]; then
        # Linux - Install based on available package manager
        if command_exists apt-get; then
            log_info "Using apt-get package manager..."
            sudo apt-get update
            sudo apt-get install -y bash-completion fd-find fzf git jq curl wget
            return 0
        elif command_exists yum; then
            log_info "Using yum package manager..."
            sudo yum install -y bash-completion fd-find fzf git jq curl wget
            return 0
        elif command_exists dnf; then
            log_info "Using dnf package manager..."
            sudo dnf install -y bash-completion fd-find fzf git jq curl wget
            return 0
        elif command_exists pacman; then
            log_info "Using pacman package manager..."
            sudo pacman -S bash-completion fd fzf git jq curl wget --noconfirm
            return 0
        else
            log_warn "No supported package manager found"
            return 1
        fi
    fi
}

# Install development tools
install_development_tools() {
    log_info "Installing development tools..."

    if [[ "$OSTYPE" == "darwin"* ]]; then
        if command_exists brew; then
            # Development tools
            log_info "Installing core development tools..."
            brew install git gh                    # Git and GitHub CLI
            brew install go                       # Go programming language
            brew install node npm                 # Node.js and npm
            brew install python3                  # Python 3

            # Optional but recommended
            brew install tmux neovim              # Terminal multiplexer and editor

            # VS Code (optional)
            install_vscode_if_needed

            # Docker Desktop (optional)
            install_docker_if_needed
        fi

    elif [[ "$OSTYPE" == "linux-gnu" ]]; then
        if command_exists apt-get; then
            # Development tools
            log_info "Installing development tools..."
            sudo apt-get install -y python3 python3-pip python3-venv
            sudo apt-get install -y nodejs npm
            sudo apt-get install -y tmux neovim

            # Manual installations needed
            log_info "Some tools require manual installation:"
            echo "  - Go: https://golang.org/dl/"
            echo "  - GitHub CLI: https://cli.github.com/"
            echo "  - VS Code: https://code.visualstudio.com/"
            echo "  - Docker: https://docs.docker.com/engine/install/"

        elif command_exists yum || command_exists dnf; then
            local pkg_mgr="yum"
            command_exists dnf && pkg_mgr="dnf"

            sudo $pkg_mgr install -y python3 python3-pip nodejs npm
            sudo $pkg_mgr install -y tmux neovim

            log_info "Additional tools may need manual installation"
        fi
    fi

    log_success "Development tools installation completed!"
}

# Install cloud and container tools
install_cloud_tools() {
    log_info "Installing cloud and container tools..."

    if [[ "$OSTYPE" == "darwin"* ]]; then
        if command_exists brew; then
            # Kubernetes tools
            log_info "Installing Kubernetes tools..."
            brew install kubectl helm k9s

            # Cloud CLI tools (optional)
            getResponse -m "Install AWS CLI?"
            if [ "$RESPONSE" = 'y' ]; then
                brew install awscli
            fi

            getResponse -m "Install Azure CLI?"
            if [ "$RESPONSE" = 'y' ]; then
                brew install azure-cli
            fi

            getResponse -m "Install DigitalOcean CLI (doctl)?"
            if [ "$RESPONSE" = 'y' ]; then
                brew install doctl
            fi

            CLOUD_TOOLS_INSTALLED=true
        fi

    elif [[ "$OSTYPE" == "linux-gnu" ]]; then
        # Kubernetes tools (manual installation on Linux)
        log_info "Installing Kubernetes tools..."

        # kubectl
        if ! command_exists kubectl; then
            log_info "Installing kubectl..."
            curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
            chmod +x kubectl
            sudo mv kubectl /usr/local/bin/
        fi

        # helm
        if ! command_exists helm; then
            log_info "Installing Helm..."
            curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
        fi

        log_info "Manual cloud CLI installation needed:"
        echo "  - AWS CLI: https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html"
        echo "  - Azure CLI: https://docs.microsoft.com/en-us/cli/azure/install-azure-cli"
        echo "  - doctl: https://docs.digitalocean.com/reference/doctl/how-to/install/"
        echo "  - k9s: https://k9scli.io/topics/install/"

        CLOUD_TOOLS_INSTALLED=true
    fi

    log_success "Cloud tools installation completed!"
}

# Setup secrets management
setup_secrets_management() {
    log_info "Setting up secrets management..."

    # Initialize secrets management
    if [ -f "$SCRIPT_DIR/.automation/auto" ]; then
        "$SCRIPT_DIR/.automation/auto" secrets init

        getResponse -m "Run interactive secrets setup wizard now?"
        if [ "$RESPONSE" = 'y' ]; then
            "$SCRIPT_DIR/.automation/auto" secrets setup
        else
            log_info "You can run the secrets setup later with: auto secrets setup"
        fi
    else
        log_warn "Automation framework not found, skipping secrets setup"
    fi
}

# Link app configs (git, ssh, etc.)
link_app_configs() {
    log_info "Linking application configurations..."

    # Git config
    if [ -d "$SCRIPT_DIR/config/git" ]; then
        if [ ! -d "$HOME/.config/git" ]; then
            mkdir -p "$HOME/.config/git"
        fi

        for file in "$SCRIPT_DIR/config/git/"*; do
            local filename=$(basename "$file")
            local target="$HOME/.config/git/$filename"
            if [ ! -e "$target" ]; then
                ln -s "$file" "$target"
                log_info "Linked git/$filename"
            fi
        done
    fi

    # SSH config (be careful with this)
    if [ -f "$SCRIPT_DIR/config/ssh/config" ]; then
        mkdir -p "$HOME/.ssh"
        if [ ! -f "$HOME/.ssh/config" ]; then
            cp "$SCRIPT_DIR/config/ssh/config" "$HOME/.ssh/config"
            chmod 600 "$HOME/.ssh/config"
            log_info "Copied SSH config"
        else
            log_info "SSH config already exists, skipping"
        fi
    fi

    log_success "Application configs linked!"
}

# Run tests to validate installation
run_installation_tests() {
    log_info "Running installation validation tests..."

    # Test shell init
    if [ -f "$SCRIPT_DIR/shell/init.sh" ]; then
        log_success "✅ Shell init.sh found"
    else
        log_error "❌ Shell init.sh not found"
    fi

    # Test bootstrap files
    if [ -f "$HOME/.bashrc" ] && grep -q "DOTFILES_ROOT" "$HOME/.bashrc"; then
        log_success "✅ ~/.bashrc bootstrap configured"
    else
        log_error "❌ ~/.bashrc not properly configured"
    fi

    if [ -f "$HOME/.zshrc" ] && grep -q "DOTFILES_ROOT" "$HOME/.zshrc"; then
        log_success "✅ ~/.zshrc bootstrap configured"
    else
        log_error "❌ ~/.zshrc not properly configured"
    fi

    # Test automation framework
    if [ -f "$SCRIPT_DIR/.automation/auto" ] && [ -x "$SCRIPT_DIR/.automation/auto" ]; then
        log_success "✅ Automation framework installed"

        # Test basic commands
        if "$SCRIPT_DIR/.automation/auto" --version >/dev/null 2>&1; then
            log_success "✅ Automation CLI working"
        else
            log_warn "⚠️ Automation CLI not responding"
        fi
    else
        log_error "❌ Automation framework not found or not executable"
    fi

    # Test function modules
    if [ -d "$SCRIPT_DIR/functions" ]; then
        log_success "✅ Function modules found"
    else
        log_error "❌ Function modules not found"
    fi

    # Test shell loading
    log_info "Testing shell init sourcing..."
    if bash -c "export DOTFILES_ROOT='$SCRIPT_DIR'; source shell/init.sh; [ -n \"\$CURRENT_SHELL\" ]" 2>/dev/null; then
        log_success "✅ Shell init loads successfully"
    else
        log_warn "⚠️ Shell init may have issues"
    fi

    # Run comprehensive tests if Makefile exists
    if [ -f "$SCRIPT_DIR/Makefile" ]; then
        getResponse -m "Run comprehensive test suite?"
        if [ "$RESPONSE" = 'y' ]; then
            log_info "Running test suite..."
            cd "$SCRIPT_DIR"
            make test-quick || log_warn "Some tests failed - check logs"
        fi
    fi
}

# Show post-installation information
show_post_install_info() {
    log_success "🎉 Installation completed successfully!"
    echo
    echo "=== Next Steps ==="
    echo

    if [ "$DOTFILES_INSTALLED" = true ]; then
        echo "1. 🔄 Restart your shell or run:"
        echo "   source ~/.bashrc     # for bash"
        echo "   source ~/.zshrc      # for zsh"
        echo
    fi

    if [ "$AUTOMATION_INSTALLED" = true ]; then
        echo "2. 🤖 Try automation commands:"
        echo "   auto --help              # Show all commands"
        echo "   auto secrets setup       # Setup API keys"
        echo "   auto dev init python my-app # Create new project"
        echo
    fi

    if [ "$CLOUD_TOOLS_INSTALLED" = true ]; then
        echo "3. ☁️ Configure cloud providers:"
        echo "   auto secrets aws         # Setup AWS credentials"
        echo "   auto secrets azure       # Setup Azure credentials"
        echo "   auto secrets github      # Setup GitHub credentials"
        echo
    fi

    echo "4. 🔧 Dotfiles management functions:"
    echo "   dotfiles_status      # Show current configuration"
    echo "   dotfiles_reload      # Reload without restart"
    echo "   dotfiles_update      # Git pull + reload"
    echo "   dotfiles_link_configs # Link app configs (git, ssh)"
    echo

    echo "5. 🧪 Test your setup:"
    echo "   make test-quick           # Quick validation tests"
    echo "   make test-auth-status     # Check authentication"
    echo

    echo "6. 📚 Documentation:"
    echo "   📖 docs/INSTALL.md - Installation guide"
    echo "   🔐 docs/automation/SECRETS.md - Secrets management"
    echo "   🤖 docs/automation/README.md - Automation framework"
    echo "   🏠 Repository: $SCRIPT_DIR"
    echo

    if [ -d "$BACKUP_DIR" ]; then
        echo "💾 Backup location: $BACKUP_DIR"
        echo
    fi

    echo "ℹ️ Installer Options:"
    echo "   ./install.sh --auto        # Non-interactive mode"
    echo "   ./install.sh --yes-to-all  # Answer yes to all prompts"
    echo "   ./install.sh --skip-gui    # Skip GUI apps (VS Code, Docker)"
    echo "   ./install.sh --dev-tools   # Only install development tools"
    echo "   ./install.sh --help        # Show all options"
    echo

    echo "🚀 Your enhanced development environment is ready!"
}

# Main installation flow
main() {
    echo
    echo "🚀 Enhanced Bash Profile & Automation Framework Installer"
    echo "=========================================================="
    echo

    # Show mode information
    if [ "$YES_TO_ALL" = true ]; then
        log_info "Running in YES-TO-ALL mode - will automatically answer 'yes' to all prompts"
        echo
    fi

    # Create backup
    create_backup

    # Install dotfiles
    getResponse -m "Install bash profile and dotfiles?"
    if [ "$RESPONSE" = 'y' ]; then
        install_dotfiles
    fi

    # Install automation framework
    getResponse -m "Setup automation framework?"
    if [ "$RESPONSE" = 'y' ]; then
        install_automation_framework
    fi

    # Link app configs
    getResponse -m "Link application configs (git, ssh)?"
    if [ "$RESPONSE" = 'y' ]; then
        link_app_configs
    fi

    # Install package manager tools
    getResponse -m "Install package manager and core tools?"
    if [ "$RESPONSE" = 'y' ]; then
        install_package_manager_tools
    fi

    # Install development tools
    getResponse -m "Install development tools (Git, Go, Node.js, Python)?"
    if [ "$RESPONSE" = 'y' ]; then
        install_development_tools
    fi

    # Install cloud tools
    getResponse -m "Install cloud and Kubernetes tools?"
    if [ "$RESPONSE" = 'y' ]; then
        install_cloud_tools
    fi

    # Setup secrets management
    getResponse -m "Setup secrets management?"
    if [ "$RESPONSE" = 'y' ]; then
        setup_secrets_management
    fi

    # Run validation tests
    getResponse -m "Run installation validation tests?"
    if [ "$RESPONSE" = 'y' ]; then
        run_installation_tests
    fi

    # Show post-installation info
    show_post_install_info
}

# Handle command line arguments
case "${1:-}" in
    "--help"|"-h")
        echo "Enhanced Bash Profile & Automation Framework Installer"
        echo
        echo "Usage: $0 [options]"
        echo
        echo "Options:"
        echo "  --help, -h     Show this help message"
        echo "  --auto         Run with all defaults (non-interactive)"
        echo "  --yes-to-all   Answer yes to all prompts (interactive but automatic)"
        echo "  --dotfiles     Install only dotfiles"
        echo "  --automation   Install only automation framework"
        echo "  --dev-tools    Install only development tools"
        echo "  --cloud-tools  Install only cloud tools"
        echo "  --skip-gui     Skip GUI applications (VS Code, Docker Desktop)"
        echo "  --test         Run tests only"
        echo
        exit 0
        ;;
    "--auto")
        # Non-interactive mode - install everything
        log_info "Running in non-interactive mode..."
        export NON_INTERACTIVE=true
        create_backup
        install_dotfiles
        install_automation_framework
        link_app_configs
        install_package_manager_tools
        install_development_tools
        install_cloud_tools
        run_installation_tests
        show_post_install_info
        ;;
    "--yes-to-all")
        # Yes to all mode - interactive but automatic
        log_info "Running in yes-to-all mode..."
        export YES_TO_ALL=true
        main
        ;;
    "--dotfiles")
        create_backup
        install_dotfiles
        show_post_install_info
        ;;
    "--automation")
        install_automation_framework
        show_post_install_info
        ;;
    "--dev-tools")
        install_package_manager_tools
        install_development_tools
        show_post_install_info
        ;;
    "--cloud-tools")
        install_cloud_tools
        show_post_install_info
        ;;
    "--skip-gui")
        export SKIP_GUI_APPS=true
        main
        ;;
    "--test")
        run_installation_tests
        ;;
    *)
        # Interactive mode
        main
        ;;
esac
