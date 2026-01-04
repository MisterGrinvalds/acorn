#!/bin/bash
# Component-Based Dotfiles Installer
# Supports dotfiles, secrets management, and development tools
# Uses component-based architecture with core/bootstrap.sh as main entry point

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

# Portable function to get user response (defaults to Y on Enter)
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
        printf "%s [Y/n]: y (auto-yes)\n" "$MESSAGE"
        RESPONSE="y"
        return 0
    fi

    # Loop until valid response (defaults to Y on empty input)
    while true; do
        printf "%s [Y/n]: " "$MESSAGE"
        read -r RESPONSE
        case "$RESPONSE" in
            y|Y|yes|Yes|YES|"") RESPONSE="y"; break ;;
            n|N|no|No|NO) RESPONSE="n"; break ;;
            *) echo "Please enter y or n (default: y)" ;;
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

    # Backup existing shell configs
    for file in .bash_profile .bashrc .zshrc .zprofile; do
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
# Dotfiles bootstrap - component-based architecture
# Generated by install.sh on $(date)
export DOTFILES_ROOT="$dotfiles_root"
[ -f "\$DOTFILES_ROOT/core/bootstrap.sh" ] && . "\$DOTFILES_ROOT/core/bootstrap.sh"
EOF
    log_success "Created ~/.bashrc"

    # Create bootstrap file for zsh
    log_info "Creating ~/.zshrc bootstrap..."
    cat > "$HOME/.zshrc" << EOF
# Dotfiles bootstrap - component-based architecture
# Generated by install.sh on $(date)
export DOTFILES_ROOT="$dotfiles_root"
[ -f "\$DOTFILES_ROOT/core/bootstrap.sh" ] && . "\$DOTFILES_ROOT/core/bootstrap.sh"
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
            brew install uv                       # Fast Python package manager

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

# Install tools for a specific component from its component.yaml
# Usage: install_component_tools <component_name>
install_component_tools() {
    local component="$1"
    local component_dir="$SCRIPT_DIR/components/$component"
    local component_yaml="$component_dir/component.yaml"

    if [ ! -f "$component_yaml" ]; then
        log_error "Component not found: $component"
        return 1
    fi

    local description
    description=$(yq -r '.description // "No description"' "$component_yaml" 2>/dev/null)
    log_info "Installing $component: $description"

    local installed_something=false

    # Platform-specific package managers
    if [[ "$OSTYPE" == "darwin"* ]]; then
        local brew_packages
        brew_packages=$(yq -r '.setup.brew // [] | join(" ")' "$component_yaml" 2>/dev/null)
        if [ -n "$brew_packages" ] && [ "$brew_packages" != "" ]; then
            log_info "Installing via Homebrew: $brew_packages"
            # shellcheck disable=SC2086
            brew install $brew_packages || log_warn "Some packages may have failed"
            installed_something=true
        fi
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        local apt_packages
        apt_packages=$(yq -r '.setup.apt // [] | join(" ")' "$component_yaml" 2>/dev/null)
        if [ -n "$apt_packages" ] && [ "$apt_packages" != "" ]; then
            log_info "Installing via apt: $apt_packages"
            # shellcheck disable=SC2086
            sudo apt-get install -y $apt_packages || log_warn "Some packages may have failed"
            installed_something=true
        fi
    fi

    # NPM packages (cross-platform)
    local npm_packages
    npm_packages=$(yq -r '.setup.npm // [] | join(" ")' "$component_yaml" 2>/dev/null)
    if [ -n "$npm_packages" ] && [ "$npm_packages" != "" ]; then
        if command_exists npm; then
            log_info "Installing via npm: $npm_packages"
            # shellcheck disable=SC2086
            npm install -g $npm_packages || log_warn "Some npm packages may have failed"
            installed_something=true
        elif command_exists pnpm; then
            log_info "Installing via pnpm: $npm_packages"
            # shellcheck disable=SC2086
            pnpm add -g $npm_packages || log_warn "Some npm packages may have failed"
            installed_something=true
        else
            log_warn "npm/pnpm not found. Install Node.js first to install: $npm_packages"
        fi
    fi

    # Show install note if nothing was installed
    if [ "$installed_something" = false ]; then
        local install_note
        install_note=$(yq -r '.setup.install_note // ""' "$component_yaml" 2>/dev/null)
        if [ -n "$install_note" ] && [ "$install_note" != "" ]; then
            log_info "Install note: $install_note"
        fi
    fi

    # Run post_install script if specified
    local post_install
    post_install=$(yq -r '.setup.post_install // ""' "$component_yaml" 2>/dev/null)
    if [ -n "$post_install" ] && [ "$post_install" != "" ] && [ -f "$component_dir/$post_install" ]; then
        log_info "Running post-install script..."
        bash "$component_dir/$post_install"
    fi

    log_success "Component $component tools installed!"
}

# List all components with their install status
list_components_status() {
    log_info "Component Installation Status"
    echo "=============================="
    echo ""

    for component_dir in "$SCRIPT_DIR/components"/*/; do
        local component
        component=$(basename "$component_dir")
        [ "$component" = "_template" ] && continue

        local component_yaml="$component_dir/component.yaml"
        [ ! -f "$component_yaml" ] && continue

        local description category
        description=$(yq -r '.description // "No description"' "$component_yaml" 2>/dev/null)
        category=$(yq -r '.category // "unknown"' "$component_yaml" 2>/dev/null)

        # Check required tools
        local tools_status="OK"
        local missing_tools=""
        local required_tools
        required_tools=$(yq -r '.requires.tools // [] | .[]' "$component_yaml" 2>/dev/null)

        for tool in $required_tools; do
            if ! command_exists "$tool"; then
                tools_status="MISSING"
                missing_tools="$missing_tools $tool"
            fi
        done

        if [ "$tools_status" = "OK" ]; then
            echo -e "  ${GREEN}✓${NC} $component [$category]"
        else
            echo -e "  ${YELLOW}○${NC} $component [$category] - missing:$missing_tools"
        fi
        echo "      $description"
    done
    echo ""
}

# Interactive component installer menu
install_components_interactive() {
    log_info "Component-based Tool Installation"
    echo "==================================="
    echo ""
    echo "Available component categories:"
    echo "  1) Core (shell, fzf, git, tmux)"
    echo "  2) Development (go, node, python)"
    echo "  3) Cloud (kubernetes, cloudflare)"
    echo "  4) AI (claude, ollama, huggingface)"
    echo "  5) Database (database tools)"
    echo "  6) All components"
    echo "  7) Select individual components"
    echo "  8) Show component status"
    echo "  9) Skip"
    echo ""

    printf "Select option [1-9]: "
    read -r choice

    case "$choice" in
        1)
            for comp in shell fzf git tmux; do
                install_component_tools "$comp"
            done
            ;;
        2)
            for comp in go node python; do
                install_component_tools "$comp"
            done
            ;;
        3)
            for comp in kubernetes cloudflare; do
                install_component_tools "$comp"
            done
            ;;
        4)
            for comp in claude ollama huggingface; do
                install_component_tools "$comp"
            done
            ;;
        5)
            install_component_tools "database"
            ;;
        6)
            for component_dir in "$SCRIPT_DIR/components"/*/; do
                local comp
                comp=$(basename "$component_dir")
                [ "$comp" = "_template" ] && continue
                install_component_tools "$comp"
            done
            ;;
        7)
            install_components_selective
            ;;
        8)
            list_components_status
            install_components_interactive
            ;;
        9)
            log_info "Skipping component installation"
            ;;
        *)
            log_warn "Invalid choice"
            ;;
    esac
}

# Select individual components to install
install_components_selective() {
    echo ""
    echo "Available components:"
    local components=()
    local i=1

    for component_dir in "$SCRIPT_DIR/components"/*/; do
        local comp
        comp=$(basename "$component_dir")
        [ "$comp" = "_template" ] && continue
        components+=("$comp")
        local desc
        desc=$(yq -r '.description // ""' "$component_dir/component.yaml" 2>/dev/null | head -c 50)
        echo "  $i) $comp - $desc"
        ((i++))
    done
    echo "  0) Done"
    echo ""

    while true; do
        printf "Enter component number (0 to finish): "
        read -r num
        [ "$num" = "0" ] && break
        if [ "$num" -ge 1 ] && [ "$num" -le "${#components[@]}" ]; then
            local selected="${components[$((num-1))]}"
            install_component_tools "$selected"
        else
            log_warn "Invalid selection"
        fi
    done
}

# Install CloudFlare/Wrangler tools
install_cloudflare_tools() {
    log_info "Installing CloudFlare CLI (wrangler)..."

    if command_exists wrangler; then
        log_info "Wrangler is already installed: $(wrangler --version 2>/dev/null | head -1)"
        return 0
    fi

    if command_exists npm; then
        log_info "Installing wrangler via npm..."
        npm install -g wrangler
        log_success "Wrangler installed successfully!"
    elif command_exists pnpm; then
        log_info "Installing wrangler via pnpm..."
        pnpm add -g wrangler
        log_success "Wrangler installed successfully!"
    else
        log_warn "npm/pnpm not found. Install Node.js first, then run: npm install -g wrangler"
        return 1
    fi
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

            getResponse -m "Install CloudFlare CLI (wrangler)?"
            if [ "$RESPONSE" = 'y' ]; then
                install_cloudflare_tools
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

        # CloudFlare CLI
        getResponse -m "Install CloudFlare CLI (wrangler)?"
        if [ "$RESPONSE" = 'y' ]; then
            install_cloudflare_tools
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

# Setup Neovim configuration
setup_neovim() {
    log_info "Setting up Neovim configuration..."

    local config_dir="${XDG_CONFIG_HOME:-$HOME/.config}/nvim"
    local repos_dir="${HOME}/Repos/personal"

    # Check if already configured
    if [ -L "$config_dir" ]; then
        local target
        target=$(readlink "$config_dir")
        log_info "Neovim already configured: $config_dir -> $target"
        return 0
    fi

    if [ -d "$config_dir" ]; then
        log_info "Neovim config directory exists at $config_dir"
        return 0
    fi

    echo ""
    echo "Neovim configuration can be stored in a separate GitHub repository."
    echo "Popular options:"
    echo "  - https://github.com/nvim-lua/kickstart.nvim (recommended starter)"
    echo "  - Your own fork of kickstart.nvim"
    echo "  - Any other Neovim config repo"
    echo ""
    printf "Enter your Neovim config GitHub repo URL (or 'skip'): "
    read -r repo_url

    if [ "$repo_url" = "skip" ] || [ -z "$repo_url" ]; then
        log_info "Skipping Neovim config setup. Run 'nvim_setup' later to configure."
        return 0
    fi

    # Extract repo name
    local repo_name
    repo_name=$(basename "$repo_url" .git)
    local repo_path="${repos_dir}/${repo_name}"

    mkdir -p "$repos_dir"

    # Clone repo
    if [ -d "$repo_path" ]; then
        log_info "Repo already exists at $repo_path"
    else
        log_info "Cloning $repo_url..."
        git clone "$repo_url" "$repo_path"
        if [ $? -ne 0 ]; then
            log_error "Failed to clone repository"
            return 1
        fi
    fi

    # Create symlink
    ln -s "$repo_path" "$config_dir"
    log_success "Neovim config linked: $config_dir -> $repo_path"
}

# Link app configs (git, ssh, ghostty, vscode, claude, etc.)
link_app_configs() {
    log_info "Linking application configurations..."

    # Check for yq
    if ! command -v yq &>/dev/null; then
        log_error "yq is required for config linking. Install with: brew install yq"
        return 1
    fi

    # Set up environment for config linking
    export DOTFILES_ROOT="$SCRIPT_DIR"
    if [[ "$OSTYPE" == "darwin"* ]]; then
        export CURRENT_PLATFORM="darwin"
    else
        export CURRENT_PLATFORM="linux"
    fi

    # Process each component
    for comp_dir in "$SCRIPT_DIR/components"/*/; do
        [ -d "$comp_dir" ] || continue
        local component=$(basename "$comp_dir")
        [ "$component" = "_template" ] && continue

        local yaml_file="${comp_dir}component.yaml"
        [ -f "$yaml_file" ] || continue

        # Check if component has config section
        local file_count
        file_count=$(yq -r '.config.files | length // 0' "$yaml_file" 2>/dev/null)
        [ "$file_count" = "0" ] || [ -z "$file_count" ] && continue

        # Check platform support
        local platforms
        platforms=$(yq -r '.platforms // []' "$yaml_file" 2>/dev/null)
        if [ "$platforms" != "[]" ] && [ "$platforms" != "null" ]; then
            if ! echo "$platforms" | grep -q "$CURRENT_PLATFORM"; then
                continue
            fi
        fi

        # Create directories first
        local dir_count i
        dir_count=$(yq -r '.config.directories | length // 0' "$yaml_file" 2>/dev/null)
        if [ "$dir_count" -gt 0 ] 2>/dev/null; then
            for i in $(seq 0 $((dir_count - 1))); do
                local target perms
                target=$(yq -r ".config.directories[$i].target" "$yaml_file")
                perms=$(yq -r ".config.directories[$i].permissions // \"\"" "$yaml_file")
                target="${target/#\~/$HOME}"
                mkdir -p "$target"
                [ -n "$perms" ] && [ "$perms" != "null" ] && chmod "$perms" "$target" 2>/dev/null
            done
        fi

        # Process each config file
        local linked_count=0
        for i in $(seq 0 $((file_count - 1))); do
            local source target method platform perms
            source=$(yq -r ".config.files[$i].source" "$yaml_file")
            target=$(yq -r ".config.files[$i].target" "$yaml_file")
            method=$(yq -r ".config.files[$i].method // \"symlink\"" "$yaml_file")
            platform=$(yq -r ".config.files[$i].platform // \"\"" "$yaml_file")
            perms=$(yq -r ".config.files[$i].permissions // \"\"" "$yaml_file")

            # Skip if platform-specific and not matching
            if [ -n "$platform" ] && [ "$platform" != "null" ] && [ "$platform" != "$CURRENT_PLATFORM" ]; then
                continue
            fi

            # Resolve paths
            local source_path="${comp_dir}${source}"
            local target_path="${target/#\~/$HOME}"

            # Ensure source exists
            if [ ! -e "$source_path" ]; then
                continue
            fi

            # Ensure target directory exists
            mkdir -p "$(dirname "$target_path")"

            # Deploy based on method
            case "$method" in
                symlink)
                    ln -sf "$source_path" "$target_path"
                    ;;
                copy)
                    cp "$source_path" "$target_path"
                    ;;
            esac

            # Apply permissions if specified
            if [ -n "$perms" ] && [ "$perms" != "null" ]; then
                chmod "$perms" "$target_path" 2>/dev/null
            fi

            linked_count=$((linked_count + 1))
        done

        if [ "$linked_count" -gt 0 ]; then
            log_info "Linked $component: $linked_count file(s)"
        fi
    done

    log_success "Application configs linked!"
}

# Run tests to validate installation
run_installation_tests() {
    log_info "Running installation validation tests..."

    # Test core bootstrap
    if [ -f "$SCRIPT_DIR/core/bootstrap.sh" ]; then
        log_success "✅ Core bootstrap.sh found"
    else
        log_error "❌ Core bootstrap.sh not found"
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

    # Test function modules
    if [ -d "$SCRIPT_DIR/functions" ]; then
        log_success "✅ Function modules found"
    else
        log_error "❌ Function modules not found"
    fi

    # Test component loader
    if [ -f "$SCRIPT_DIR/core/loader.sh" ]; then
        log_success "✅ Component loader found"
    else
        log_error "❌ Component loader not found"
    fi

    # Test components directory
    if [ -d "$SCRIPT_DIR/components" ]; then
        local component_count=$(find "$SCRIPT_DIR/components" -name "component.yaml" | wc -l | tr -d ' ')
        log_success "✅ Components directory found ($component_count components)"
    else
        log_warn "⚠️ Components directory not found"
    fi

    # Test shell loading
    log_info "Testing bootstrap loading..."
    if bash -c "export DOTFILES_ROOT='$SCRIPT_DIR'; source core/bootstrap.sh; [ -n \"\$CURRENT_SHELL\" ]" 2>/dev/null; then
        log_success "✅ Bootstrap loads successfully"
    else
        log_warn "⚠️ Bootstrap may have issues"
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

    if [ "$CLOUD_TOOLS_INSTALLED" = true ]; then
        echo "2. ☁️ Configure cloud providers:"
        echo "   Configure AWS: aws configure"
        echo "   Configure Azure: az login"
        echo "   Configure GitHub: gh auth login"
        echo
    fi

    echo "3. 🔧 Dotfiles management functions:"
    echo "   dotfiles_status      # Show current configuration"
    echo "   dotfiles_reload      # Reload without restart"
    echo "   dotfiles_update      # Git pull + reload"
    echo "   dotfiles_link_configs # Link app configs (git, ssh)"
    echo
    echo "4. 📦 Component System:"
    echo "   Components are loaded from: \$DOTFILES_ROOT/components/"
    echo "   Disable a component:  export DOTFILES_DISABLE_<NAME>=1"
    echo "   Whitelist components: export DOTFILES_COMPONENTS=\"git,fzf\""
    echo

    echo "5. 🧪 Test your setup:"
    echo "   make test-quick           # Quick validation tests"
    echo "   make test-auth-status     # Check authentication"
    echo

    echo "6. 📚 Documentation:"
    echo "   📖 docs/INSTALL.md - Installation guide"
    echo "   📖 CLAUDE.md - Architecture guide"
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
    echo "🚀 Component-Based Dotfiles Installer"
    echo "======================================"
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

    # Setup Neovim
    getResponse -m "Setup Neovim configuration (external repo)?"
    if [ "$RESPONSE" = 'y' ]; then
        setup_neovim
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
        echo "Component-Based Dotfiles Installer"
        echo
        echo "Usage: $0 [options]"
        echo
        echo "Options:"
        echo "  --help, -h          Show this help message"
        echo "  --auto              Run with all defaults (non-interactive)"
        echo "  --yes-to-all        Answer yes to all prompts (interactive but automatic)"
        echo "  --dotfiles          Install only dotfiles"
        echo "  --dev-tools         Install only development tools"
        echo "  --cloud-tools       Install only cloud tools"
        echo "  --skip-gui          Skip GUI applications (VS Code, Docker Desktop)"
        echo "  --minimal           Minimal install (dotfiles + app configs only)"
        echo "  --test              Run tests only"
        echo ""
        echo "Component-based installation:"
        echo "  --components        Interactive component-based tool installer"
        echo "  --component <name>  Install tools for a specific component"
        echo "  --list-components   List all components and their status"
        echo
        echo "Examples:"
        echo "  $0 --components              # Interactive component menu"
        echo "  $0 --component cloudflare    # Install CloudFlare/wrangler tools"
        echo "  $0 --component python        # Install Python tools (python3, uv)"
        echo "  $0 --list-components         # Show which tools are installed"
        echo
        exit 0
        ;;
    "--auto")
        # Non-interactive mode - install everything
        log_info "Running in non-interactive mode..."
        export NON_INTERACTIVE=true
        create_backup
        install_dotfiles
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
    "--minimal")
        # Minimal install for foreign/remote environments
        log_info "Running minimal install (dotfiles + app configs only)..."
        create_backup
        install_dotfiles
        link_app_configs
        log_success "Minimal installation complete!"
        echo ""
        echo "Restart your shell or run: source ~/.bashrc"
        ;;
    "--test")
        run_installation_tests
        ;;
    "--components")
        # Interactive component-based installer
        install_components_interactive
        ;;
    "--component")
        # Install specific component
        if [ -z "${2:-}" ]; then
            log_error "Usage: $0 --component <component_name>"
            echo "Available components:"
            for d in "$SCRIPT_DIR/components"/*/; do
                comp=$(basename "$d")
                [ "$comp" != "_template" ] && echo "  - $comp"
            done
            exit 1
        fi
        install_component_tools "$2"
        ;;
    "--list-components")
        list_components_status
        ;;
    *)
        # Interactive mode
        main
        ;;
esac
