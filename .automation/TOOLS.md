# 🛠️ External Tools Management

Comprehensive management system for external dependencies, CLI tools, and development environments with automatic update checking and selective upgrades.

## 🎯 Overview

The tools management system provides:
- **Comprehensive inventory** of 40+ external tools
- **Version checking** and update detection
- **Interactive updates** with user prompting
- **Category-based organization** (system, development, cloud, etc.)
- **Multiple installation methods** (package managers, manual, language-specific)
- **Cross-platform support** (macOS, Linux distributions)

## 🚀 Quick Start

```bash
# Show comprehensive tools status
auto tools status

# List all tools with current status
auto tools list

# Check versions of all tools
auto tools check

# Interactive updates with prompting
auto tools update

# Update specific tool
auto tools update kubectl

# Show missing tools
auto tools missing

# Install specific tool
auto tools install go
```

## 📦 Tool Categories

### 🖥️ System Tools
- **Package Managers**: brew (macOS), apt/yum/dnf/pacman (Linux)
- **Network Tools**: curl, wget
- **JSON Processing**: jq
- **Shell Enhancement**: bash-completion

### 🔍 Shell Enhancement
- **Fuzzy Finder**: fzf - Interactive file/command finder
- **File Search**: fd - Modern find replacement
- **Search Integration**: ripgrep support for enhanced searching

### 💻 Development Languages
- **Go**: golang.org compiler and tools
- **Node.js**: JavaScript runtime and npm package manager
- **Python**: Python 3 interpreter and pip package manager

### 📝 Editors & IDEs
- **Neovim**: Modern vim-based editor
- **VS Code**: Microsoft Visual Studio Code
- **Terminal**: tmux terminal multiplexer

### ☁️ Cloud Provider CLIs
- **AWS CLI**: Amazon Web Services command line
- **Azure CLI**: Microsoft Azure command line
- **doctl**: DigitalOcean command line tool

### ☸️ Kubernetes Ecosystem
- **kubectl**: Kubernetes command line tool
- **Helm**: Kubernetes package manager
- **k9s**: Terminal-based Kubernetes UI

### 🐳 Container Tools
- **Docker**: Container platform and CLI
- **Container registries**: Integration with Docker Hub, GitHub Container Registry

### 🔗 Development Integration
- **GitHub CLI**: gh - GitHub command line tool
- **Git**: Version control system

### 🐹 Go Ecosystem
- **goimports**: Import management tool
- **golangci-lint**: Go linting suite
- **cobra-cli**: CLI application framework

### 🐍 Python Ecosystem
- **black**: Code formatter
- **isort**: Import sorter
- **flake8**: Linting tool
- **pytest**: Testing framework
- **uvicorn**: ASGI server for FastAPI
- **ipython**: Interactive Python shell

## 🔧 Installation Methods

### Package Manager Installation
```bash
# macOS (Homebrew)
brew install <tool>

# Ubuntu/Debian
sudo apt-get install <tool>

# CentOS/RHEL
sudo yum install <tool>

# Fedora
sudo dnf install <tool>

# Arch Linux
sudo pacman -S <tool>
```

### Language Package Managers
```bash
# Python tools
pip3 install <package>

# Go tools
go install <package>@latest

# Node.js tools
npm install -g <package>
```

### Manual Installation
Some tools require manual installation:
- **Go**: Download from https://golang.org/dl/
- **VS Code**: Download from https://code.visualstudio.com/
- **Docker**: Follow instructions at https://docs.docker.com/engine/install/

## 📊 Tool Status and Information

### Comprehensive Status Report
```bash
auto tools status
```
Shows:
- System information (OS, package manager, shell)
- Tools organized by category
- Installation status and versions
- Overall coverage statistics

### Version Checking
```bash
# Check all tool versions
auto tools check

# Check specific tool
auto tools check kubectl

# Quick version check for common tools
versions
```

### Missing Tools Analysis
```bash
# Show missing tools
auto tools missing

# Show tools by category
auto tools categories

# Show tools requiring manual installation
manual-tools
```

## 🔄 Update Management

### Interactive Updates
```bash
# Update all tools with prompting
auto tools update

# Update specific tool
auto tools update go

# Force update all without prompting
auto tools update --force
```

### Update Strategies by Tool Type

#### Package Manager Tools
- **Homebrew**: `brew upgrade <tool>`
- **APT**: `sudo apt-get install <tool>`
- **YUM/DNF**: `sudo yum/dnf update <tool>`

#### Language-Specific Tools
- **Go tools**: `go install <package>@latest`
- **Python tools**: `pip3 install --upgrade <package>`
- **Node.js tools**: `npm update -g`

#### Manual Update Tools
- **Go language**: Manual download and installation
- **VS Code**: Built-in updater or package manager
- **Docker Desktop**: Built-in updater (macOS) or manual (Linux)

#### Cloud CLI Tools
- **AWS CLI**: Package manager or manual download
- **Azure CLI**: `az upgrade` or package manager
- **doctl**: Package manager or manual download

### Specialized Update Functions

#### Kubernetes Tools
```bash
# kubectl - Latest stable release
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"

# Helm - Install script
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
```

#### Cloud Tools
```bash
# AWS CLI (Linux)
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip && sudo ./aws/install

# Azure CLI (Linux)
curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash
```

## 🧪 Testing and Validation

### Testing Framework Integration
```bash
# Test tools module
make test-tools-module

# Test specific tool availability
make test-required-tools

# Show tools status via Makefile
make tools-status

# Interactive updates via Makefile
make tools-update
```

### Shell Integration Testing
```bash
# Quick version check
quick_versions

# Enhanced which command
whichx python3

# System-wide package update
smart_update
```

## 📱 Shell Integration

### Available Commands
```bash
# Status and listing
tools                   # Show comprehensive status
tools-list             # List all tools
tools-check [tool]     # Check versions
versions               # Quick version check

# Updates and installation
tools-update [tool]    # Interactive updates
tools-install <tool>   # Install specific tool
system-update          # Smart package manager update

# Analysis
tools-missing          # Show missing tools
tools-categories       # Show by category
tools-outdated         # Show potentially outdated
manual-tools           # Show manual installation needed

# Utilities
whichx <command>       # Enhanced which command
```

### Aliases Available
- `tools` → `tools_status`
- `tools-list` → `list_tools`
- `tools-check` → `check_tools`
- `tools-update` → `update_tools`
- `versions` → `quick_versions`
- `system-update` → `smart_update`

## 🔍 Advanced Features

### Category Filtering
```bash
# List only cloud tools
auto tools list --category cloud

# Check only development tools
auto tools check --category development
```

### Batch Operations
```bash
# Update all Go tools
for tool in goimports golangci-lint cobra-cli; do
    auto tools update "$tool"
done

# Install all missing cloud tools
auto tools missing | grep cloud | while read tool; do
    auto tools install "$tool"
done
```

### Custom Tool Addition
To add new tools to the system:

1. **Add to tool definitions** in `.automation/modules/tools.sh`:
```bash
TOOLS["newtool"]="category:subcategory"
INSTALL_METHODS["newtool"]="brew:newtool|apt:newtool"
VERSION_COMMANDS["newtool"]="newtool --version"
UPDATE_COMMANDS["newtool"]="update_via_package_manager newtool"
```

2. **Test the addition**:
```bash
auto tools check newtool
auto tools install newtool
```

## 🌍 Cross-Platform Considerations

### macOS (Darwin)
- **Primary package manager**: Homebrew
- **Advantages**: Consistent installation experience
- **Cask support**: GUI applications via `brew install --cask`

### Linux Distributions

#### Ubuntu/Debian
- **Package manager**: apt-get
- **Special handling**: Some tools use different package names (fd-find vs fd)

#### CentOS/RHEL
- **Package managers**: yum (older) or dnf (newer)
- **EPEL repository**: May be required for additional tools

#### Fedora
- **Package manager**: dnf
- **Generally up-to-date**: Often has latest versions

#### Arch Linux
- **Package manager**: pacman
- **AUR support**: Additional tools available via AUR helpers

### Platform-Specific Tools
- **XQuartz**: macOS-only X11 server
- **Linuxbrew**: Alternative Homebrew for Linux
- **Windows Subsystem for Linux**: Partial support consideration

## 🔒 Security Considerations

### Tool Verification
- **Package signatures**: Verify package manager signatures
- **Checksums**: Verify download checksums for manual installations
- **Official sources**: Only install from official repositories

### Update Security
- **Automatic updates**: Consider security implications
- **Staging**: Test updates in non-production environments
- **Rollback capability**: Maintain ability to rollback problematic updates

## 🚨 Troubleshooting

### Common Issues

#### Tool Not Found After Installation
```bash
# Check PATH
echo $PATH

# Reload shell
source ~/.bash_profile

# Check installation location
which <tool>
```

#### Version Check Failures
```bash
# Check if tool supports --version
<tool> --help

# Try alternative version commands
<tool> version
<tool> -V
```

#### Package Manager Issues
```bash
# Update package manager
sudo apt-get update  # Ubuntu/Debian
brew update          # macOS

# Check for broken packages
sudo apt-get check   # Ubuntu/Debian
brew doctor          # macOS
```

#### Permission Issues
```bash
# Fix homebrew permissions (macOS)
sudo chown -R $(whoami) /opt/homebrew

# Fix pip permissions
pip3 install --user <package>
```

### Debug Mode
```bash
# Enable verbose logging
export AUTO_LOG_LEVEL=DEBUG

# Run with debug output
auto tools check --verbose

# Check automation logs
tail -f ~/.automation/logs/automation.log
```

## 📈 Integration with CI/CD

### GitHub Actions
```yaml
name: Tool Management Test
on: [push, pull_request]
jobs:
  test-tools:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Test Tools
      run: |
        ./initialize.sh --auto
        make test-tools-module
        auto tools check
```

### Pre-commit Hooks
```bash
# .pre-commit-config.yaml
repos:
- repo: local
  hooks:
  - id: tools-check
    name: Check tool versions
    entry: auto tools check
    language: system
```

## 🔮 Future Enhancements

### Planned Features
- **Automatic update detection**: Check for newer versions automatically
- **Update notifications**: Alert users when updates are available
- **Dependency resolution**: Handle tool dependencies automatically
- **Configuration backup**: Backup tool configurations before updates
- **Rollback capability**: Rollback problematic updates
- **Custom tool repositories**: Support for custom tool sources

### Version Comparison
- **Semantic version parsing**: Compare version numbers intelligently
- **Update recommendations**: Suggest which updates are critical
- **Change logs**: Show what's changed in new versions

---

**The tools management system provides enterprise-grade dependency management while maintaining the simplicity and flexibility needed for individual developer workflows.**