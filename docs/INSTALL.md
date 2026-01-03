# Installation Guide

Complete installation guide for the Component-Based Dotfiles System with secrets management and cloud integration.

## Prerequisites

### Operating System Support
- **macOS** 10.15+ (Catalina or later)
- **Linux** (Ubuntu 18.04+, CentOS 7+, Fedora 30+, Arch Linux)

### Required Tools
- **bash** 4.0+ or **zsh** 5.0+
- **curl** or **wget**
- **git** 2.0+

### Optional but Recommended
- **Homebrew** (macOS) or package manager (Linux)
- **jq** for JSON processing
- **fzf** for fuzzy finding

## Installation Methods

### Method 1: Interactive Installation (Recommended)

```bash
# Clone the repository
git clone https://github.com/your-username/bash-profile.git ~/.config/dotfiles
cd ~/.config/dotfiles

# Run interactive installer
./install.sh
```

The installer will guide you through:
1. **Dotfiles installation** - Core bash/zsh profile setup
2. **Application configs** - Link git, ssh, and other app configurations
3. **Package management** - Install essential tools
4. **Development tools** - Git, Go, Node.js, Python, VS Code
5. **Cloud tools** - AWS CLI, Azure CLI, kubectl, Helm
6. **Neovim setup** - External config repository linking
7. **Validation testing** - Verify everything works

### Method 2: Non-Interactive Installation

```bash
# Install everything with defaults
./install.sh --auto
```

### Method 3: Selective Installation

```bash
# Install only specific components
./install.sh --dotfiles      # Dotfiles only
./install.sh --dev-tools     # Development tools only
./install.sh --cloud-tools   # Cloud tools only
```

## Platform-Specific Setup

### macOS Setup

```bash
# Install Homebrew (if not present)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Run installer (will use Homebrew)
./install.sh

# Recommended additional tools
brew install --cask iterm2        # Better terminal
brew install --cask rectangle     # Window management
brew install tree htop            # System utilities
```

### Ubuntu/Debian Setup

```bash
# Update package list
sudo apt-get update

# Install prerequisites
sudo apt-get install -y curl git build-essential

# Run installer
./install.sh

# Manual cloud tools installation
# AWS CLI
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip && sudo ./aws/install

# Azure CLI
curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash
```

### CentOS/RHEL Setup

```bash
# Install EPEL repository
sudo yum install -y epel-release

# Install prerequisites
sudo yum install -y curl git gcc make

# Run installer
./install.sh
```

### Arch Linux Setup

```bash
# Install prerequisites
sudo pacman -S curl git base-devel

# Run installer
./install.sh
```

## Cloud Provider Setup

After installation, configure your cloud providers:

### AWS
```bash
aws configure
# OR set environment variables
export AWS_ACCESS_KEY_ID="your-access-key"
export AWS_SECRET_ACCESS_KEY="your-secret-key"
export AWS_DEFAULT_REGION="us-east-1"
```

### Azure
```bash
az login
# OR with service principal
az login --service-principal -u <client-id> -p <client-secret> --tenant <tenant-id>
```

### DigitalOcean
```bash
doctl auth init
# OR set environment variable
export DIGITALOCEAN_TOKEN="your-token"
```

### GitHub
```bash
gh auth login
# OR set token
export GITHUB_TOKEN="your-personal-access-token"
```

## Testing Your Installation

### Quick Validation
```bash
# Test basic functionality
make test-quick

# Test authentication status
make test-auth-status

# Test required tools
make test-required-tools
```

### Comprehensive Testing
```bash
# Full test suite
make test-comprehensive

# Test components
make test-components
make component-status
```

### Manual Testing
```bash
# Reload shell
source ~/.bashrc     # for bash
source ~/.zshrc      # for zsh

# Check dotfiles status
dotfiles_status

# Test cloud connectivity
aws sts get-caller-identity    # AWS
az account show               # Azure
doctl account get             # DigitalOcean
gh auth status               # GitHub
kubectl cluster-info         # Kubernetes
```

## Configuration

### Environment Variables

The system uses XDG Base Directory specification:

```bash
# XDG directories (set automatically)
export XDG_CONFIG_HOME="$HOME/.config"
export XDG_DATA_HOME="$HOME/.local/share"
export XDG_CACHE_HOME="$HOME/.cache"
export XDG_STATE_HOME="$HOME/.local/state"

# Dotfiles location
export DOTFILES_ROOT="$HOME/.config/dotfiles"
```

### Component Configuration

Disable specific components:
```bash
export DOTFILES_DISABLE_KUBERNETES=1
export DOTFILES_DISABLE_OLLAMA=1
```

Whitelist specific components:
```bash
export DOTFILES_COMPONENTS="git,fzf,python,shell"
```

## Troubleshooting

### Common Issues

#### Permission Denied
```bash
# Fix script permissions
chmod +x install.sh
```

#### Components Not Loading
```bash
# Check component status
make component-status

# Validate component.yaml files
make component-validate
```

#### Shell Not Detected
```bash
# Check shell detection
echo "Shell: $CURRENT_SHELL"
echo "Platform: $CURRENT_PLATFORM"
```

### Debugging

```bash
# Reload dotfiles
dotfiles_reload

# Check XDG compliance
dotfiles_audit

# View loaded components
make component-list
```

### Reset Installation

```bash
# Remove bootstrap files
dotfiles_eject

# Restore from backup (if available)
cp -r ~/.dotfiles-backup-*/* ~/

# Reinstall
./install.sh
```

## Post-Installation

### Dotfiles Management

```bash
dotfiles_status         # Show current state
dotfiles_reload         # Reload without restart
dotfiles_update         # Git pull + reload
dotfiles_link_configs   # Link app configs
dotfiles_unlink_configs # Remove app config links
```

### Development Setup

```bash
# Python virtual environment
mkvenv myproject
venv myproject

# Go project
gonew myproject

# Node.js setup
nvm_setup
node_init
```

## Further Reading

- **[CLAUDE.md](/CLAUDE.md)** - Architecture and design guide
- **[README.md](/README.md)** - Project overview

---

**Need Help?** Open an issue or check the troubleshooting section above.
