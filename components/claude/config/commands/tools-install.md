---
description: Help install system tools using appropriate package manager
argument-hint: <tool-name>
allowed-tools: Read, Bash
---

## Task

Help the user install tools using the appropriate package manager for their platform.

## Quick Install

```bash
# Using dotfiles function (via automation framework)
tools-install <tool>
install_tool <tool>
```

## Platform Detection

```bash
# Check platform
echo $CURRENT_PLATFORM  # darwin or linux

# Check package manager
which brew      # macOS (Homebrew)
which apt-get   # Debian/Ubuntu
which dnf       # Fedora/RHEL
which pacman    # Arch
```

## Common Tools Installation

### System Utilities

| Tool | Homebrew | APT | DNF |
|------|----------|-----|-----|
| git | `brew install git` | `sudo apt install git` | `sudo dnf install git` |
| curl | `brew install curl` | `sudo apt install curl` | `sudo dnf install curl` |
| wget | `brew install wget` | `sudo apt install wget` | `sudo dnf install wget` |
| jq | `brew install jq` | `sudo apt install jq` | `sudo dnf install jq` |
| yq | `brew install yq` | snap install yq | `sudo dnf install yq` |
| tree | `brew install tree` | `sudo apt install tree` | `sudo dnf install tree` |
| htop | `brew install htop` | `sudo apt install htop` | `sudo dnf install htop` |

### Search Tools

| Tool | Homebrew | APT | DNF |
|------|----------|-----|-----|
| fzf | `brew install fzf` | `sudo apt install fzf` | `sudo dnf install fzf` |
| ripgrep | `brew install ripgrep` | `sudo apt install ripgrep` | `sudo dnf install ripgrep` |
| fd | `brew install fd` | `sudo apt install fd-find` | `sudo dnf install fd-find` |

### Development

| Tool | Homebrew | APT | Notes |
|------|----------|-----|-------|
| neovim | `brew install neovim` | `sudo apt install neovim` | |
| tmux | `brew install tmux` | `sudo apt install tmux` | |
| gh | `brew install gh` | See GitHub docs | |
| make | `brew install make` | `sudo apt install build-essential` | |

### Languages

```bash
# Go
brew install go                    # macOS
# Linux: Download from golang.org

# Node.js (via nvm)
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
nvm install --lts

# Python (via pyenv or system)
brew install python3               # macOS
sudo apt install python3 python3-pip  # Ubuntu

# UV (Python package manager)
curl -LsSf https://astral.sh/uv/install.sh | sh

# Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
```

### Cloud Tools

```bash
# Docker
brew install --cask docker         # macOS
# Linux: See Docker docs

# kubectl
brew install kubectl               # macOS
# Linux: See Kubernetes docs

# Helm
brew install helm                  # macOS
sudo snap install helm --classic   # Ubuntu

# AWS CLI
brew install awscli                # macOS
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"  # Linux

# GitHub CLI
brew install gh                    # macOS
# Linux: See GitHub docs

# Wrangler (CloudFlare)
npm install -g wrangler
# or: pnpm add -g wrangler
```

### Database Tools

```bash
# PostgreSQL client
brew install pgcli                 # macOS
sudo apt install pgcli             # Ubuntu

# MySQL client
brew install mycli                 # macOS
pip install mycli                  # Any platform

# Redis
brew install redis                 # macOS
sudo apt install redis-tools       # Ubuntu

# MongoDB Shell
brew install mongosh               # macOS
# Linux: See MongoDB docs
```

## Bulk Installation

### macOS (Homebrew Bundle)
```bash
# Create Brewfile
cat > Brewfile << 'EOF'
brew "git"
brew "curl"
brew "jq"
brew "fzf"
brew "ripgrep"
brew "neovim"
brew "tmux"
brew "go"
brew "node"
cask "docker"
cask "visual-studio-code"
EOF

# Install all
brew bundle
```

### Ubuntu
```bash
sudo apt update
sudo apt install -y \
    git curl wget jq \
    fzf ripgrep \
    neovim tmux \
    build-essential \
    python3 python3-pip
```

## Post-Install Verification

```bash
# Check installation
whichx <tool>

# Or check all versions
versions
```

## Troubleshooting

### Command Not Found After Install
```bash
# Refresh shell
source ~/.bashrc  # or ~/.zshrc

# Check PATH
echo $PATH

# Verify installation location
which <tool>
brew --prefix <tool>  # macOS
```

### Permission Denied
```bash
# Linux: Use sudo for system packages
sudo apt install <package>

# macOS: Fix Homebrew permissions
sudo chown -R $(whoami) $(brew --prefix)/*
```

### Outdated Package
```bash
# Update package manager first
brew update && brew upgrade <tool>  # macOS
sudo apt update && sudo apt upgrade <tool>  # Ubuntu
```
