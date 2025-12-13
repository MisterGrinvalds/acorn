# 🚀 Installation Guide

Complete installation guide for the Enhanced Bash Profile & Automation Framework with secrets management and cloud integration.

## 📋 Prerequisites

### Operating System Support
- **macOS** 10.15+ (Catalina or later)
- **Linux** (Ubuntu 18.04+, CentOS 7+, Fedora 30+, Arch Linux)

### Required Tools
- **bash** 4.0+ or **zsh** 5.0+
- **curl** or **wget**
- **git** 2.0+
- **rsync**

### Optional but Recommended
- **Homebrew** (macOS) or package manager (Linux)
- **jq** for JSON processing
- **fzf** for fuzzy finding

## 🎯 Installation Methods

### Method 1: Interactive Installation (Recommended)

```bash
# Clone the repository
git clone https://github.com/your-username/bash-profile.git
cd bash-profile

# Run interactive installer
./initialize.sh
```

The installer will guide you through:
1. **Dotfiles installation** - Core bash/zsh profile setup
2. **Automation framework** - CLI tools and modules
3. **Package management** - Install essential tools
4. **Development tools** - Git, Go, Node.js, Python, VS Code
5. **Cloud tools** - AWS CLI, Azure CLI, kubectl, Helm
6. **Secrets management** - API key configuration
7. **Validation testing** - Verify everything works

### Method 2: Non-Interactive Installation

```bash
# Install everything with defaults
./initialize.sh --auto
```

### Method 3: Selective Installation

```bash
# Install only specific components
./initialize.sh --dotfiles      # Dotfiles only
./initialize.sh --automation    # Automation framework only
./initialize.sh --dev-tools     # Development tools only
./initialize.sh --cloud-tools   # Cloud tools only
```

### Method 4: Manual Installation

```bash
# 1. Install dotfiles manually
export DOTFILES="$HOME"
rsync -avr --exclude=".git*" ./ "$DOTFILES"

# 2. Setup automation framework
chmod +x .automation/auto .automation/setup.sh
.automation/setup.sh

# 3. Initialize secrets management
auto secrets init

# 4. Create shell links
ln -s ~/.bash_profile ~/.bashrc    # for bash
ln -s ~/.bash_profile ~/.zshrc     # for zsh
```

## 🔧 Platform-Specific Setup

### macOS Setup

```bash
# Install Homebrew (if not present)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Run installer (will use Homebrew)
./initialize.sh

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
sudo apt-get install -y curl git rsync build-essential

# Run installer
./initialize.sh

# Manual cloud tools installation
# AWS CLI
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip && sudo ./aws/install

# Azure CLI
curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash

# DigitalOcean CLI
cd ~ && wget https://github.com/digitalocean/doctl/releases/download/v1.94.0/doctl-1.94.0-linux-amd64.tar.gz
tar xf doctl-1.94.0-linux-amd64.tar.gz && sudo mv doctl /usr/local/bin
```

### CentOS/RHEL Setup

```bash
# Install EPEL repository
sudo yum install -y epel-release

# Install prerequisites
sudo yum install -y curl git rsync gcc make

# Run installer
./initialize.sh

# Install development tools group
sudo yum groupinstall -y "Development Tools"
```

### Arch Linux Setup

```bash
# Install prerequisites
sudo pacman -S curl git rsync base-devel

# Run installer
./initialize.sh

# Install AUR helper (optional)
git clone https://aur.archlinux.org/yay.git
cd yay && makepkg -si
```

## 🔐 Secrets Management Setup

After installation, configure your API keys and credentials:

### Quick Setup
```bash
# Interactive setup wizard
auto secrets setup

# Check what's missing
auto secrets check-requirements

# Validate configured keys
auto secrets validate
```

### Provider-Specific Setup

#### AWS
```bash
# Method 1: Use automation wizard
auto secrets aws

# Method 2: Manual configuration
aws configure
# OR set environment variables
export AWS_ACCESS_KEY_ID="your-access-key"
export AWS_SECRET_ACCESS_KEY="your-secret-key"
export AWS_DEFAULT_REGION="us-east-1"
```

#### Azure
```bash
# Method 1: Use automation wizard
auto secrets azure

# Method 2: Azure CLI login
az login

# Method 3: Service principal
az login --service-principal -u <client-id> -p <client-secret> --tenant <tenant-id>
```

#### DigitalOcean
```bash
# Method 1: Use automation wizard
auto secrets digitalocean

# Method 2: Manual setup
doctl auth init
# OR set environment variable
export DIGITALOCEAN_TOKEN="your-token"
```

#### GitHub
```bash
# Method 1: Use automation wizard
auto secrets github

# Method 2: GitHub CLI
gh auth login

# Method 3: Manual token
export GITHUB_TOKEN="your-personal-access-token"
```

## 🧪 Testing Your Installation

### Quick Validation
```bash
# Test basic functionality
make test-quick

# Test API keys and authentication
make test-api-keys
make test-auth-status

# Test required tools
make test-required-tools
```

### Comprehensive Testing
```bash
# Full test suite
make test-comprehensive

# Generate test report
make test-report
```

### Manual Testing
```bash
# Test shell loading
source ~/.bash_profile

# Test automation CLI
auto --help
auto --version

# Test secrets management
auto secrets check-requirements
load_secrets

# Test cloud connectivity
aws sts get-caller-identity    # AWS
az account show               # Azure
doctl account get             # DigitalOcean
gh auth status               # GitHub
kubectl cluster-info         # Kubernetes
```

## 🔧 Configuration

### Environment Variables

Add to your shell profile or secrets file:

```bash
# Automation settings
export AUTO_LOG_LEVEL=INFO
export AUTO_PARALLEL_JOBS=4
export AUTO_TIMEOUT=300

# Auto-load secrets on shell startup
export AUTO_LOAD_SECRETS=true

# Default project directory
export DEV_PROJECTS_DIR=$HOME/projects

# GitHub settings
export GITHUB_DEFAULT_VISIBILITY=public

# Kubernetes settings
export K8S_DEFAULT_NAMESPACE=default
```

### Custom Configuration

```bash
# Create custom configuration
auto config profile create work
auto config template create my-stack

# Edit automation config
nano ~/.automation/config/automation.conf

# Edit secrets template
nano ~/.automation/secrets/template.env
```

## 🚨 Troubleshooting

### Common Issues

#### Permission Denied
```bash
# Fix script permissions
chmod +x .automation/auto .automation/setup.sh initialize.sh
find .automation -name "*.sh" -exec chmod +x {} \;
```

#### Command Not Found
```bash
# Add automation to PATH
export PATH="$HOME/.automation:$PATH"
echo 'export PATH="$HOME/.automation:$PATH"' >> ~/.bash_profile
source ~/.bash_profile
```

#### Secrets Not Loading
```bash
# Check secrets file
ls -la ~/.automation/secrets/.env

# Reload secrets
load_secrets

# Re-initialize secrets
auto secrets init
```

#### Cloud CLI Issues
```bash
# Check CLI installation
make test-required-tools

# Verify authentication
make test-auth-status

# Re-authenticate
auto secrets setup
```

### Debugging

```bash
# Enable debug logging
export AUTO_LOG_LEVEL=DEBUG

# Check automation logs
tail -f ~/.automation/logs/automation.log

# Run tests with verbose output
make test-quick TEST_VERBOSE=true
```

### Reset Installation

```bash
# Remove dotfiles (backup first!)
rm -rf ~/.bash_profile ~/.bash_profile.dir ~/.bash_tools ~/.automation

# Restore from backup
cp -r ~/.dotfiles-backup-YYYYMMDD_HHMMSS/* ~/

# Reinstall
./initialize.sh
```

## 🔄 Post-Installation

### Shell Configuration

Add to your shell startup files:

```bash
# ~/.bash_profile or ~/.zshrc
source ~/.bash_profile

# Auto-load secrets (optional)
export AUTO_LOAD_SECRETS=true

# Custom aliases
alias ll='ls -la'
alias la='ls -A'
alias l='ls -CF'
```

### Development Setup

```bash
# Create projects directory
mkdir -p ~/projects

# Setup common development tools
auto dev init python sample-api
auto dev init go sample-cli --cobra
auto dev init typescript sample-app

# Setup Kubernetes environment
auto k8s cluster info
auto k8s monitoring
```

### Cloud Environment

```bash
# Setup cloud resources
auto cloud status
auto aws ec2 list
auto azure vm list
auto digitalocean droplets list

# Deploy sample application
auto cloud multi-deploy sample-app
```

## 🚀 Advanced Features

### Custom Modules

Create your own automation modules:

```bash
# Create custom module
cat > ~/.automation/modules/custom.sh << 'EOF'
#!/bin/bash
custom_help() {
    echo "Custom module help"
}

custom_main() {
    case "$1" in
        "help") custom_help ;;
        *) echo "Custom action: $1" ;;
    esac
}
EOF

# Use custom module
auto custom help
```

### CI/CD Integration

```bash
# GitHub Actions example
name: Test Bash Profile
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Install
      run: ./initialize.sh --auto
    - name: Test
      run: make test-ci
```

### Container Development

```bash
# Use with dev containers
auto dev init python my-api --devcontainer
auto dev init go my-cli --docker

# Kubernetes development
auto k8s deploy my-app development
auto k8s port-forward my-app 8080:80
```

## 📚 Further Reading

- **[Automation Framework](/.automation/README.md)** - Complete automation guide
- **[Secrets Management](/.automation/SECRETS.md)** - Security and API keys
- **[Testing Guide](/tests/README.md)** - Testing and validation
- **[CLAUDE.md](/CLAUDE.md)** - Architecture and design

---

**Need Help?** Open an issue or check the troubleshooting section above.