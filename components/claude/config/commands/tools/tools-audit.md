---
description: Audit installed tools, versions, and missing dependencies
argument-hint: [category: all|languages|cloud|system|dev]
allowed-tools: Bash
---

## Task

Audit the user's installed tools, check versions, and identify missing or outdated tools.

## Quick Audit Commands

```bash
# Check common tool versions
versions

# Using automation framework
tools-check    # Check all tools
tools-missing  # Show missing tools
```

## Comprehensive Audit

### System Tools
```bash
echo "=== System Tools Audit ==="

# Core utilities
for tool in git curl wget jq yq rsync; do
    whichx $tool 2>/dev/null || echo "$tool: Not installed"
done

# Search tools
for tool in fzf fd rg ag; do
    whichx $tool 2>/dev/null || echo "$tool: Not installed"
done

# Monitoring
for tool in htop btop; do
    whichx $tool 2>/dev/null || echo "$tool: Not installed"
done
```

### Programming Languages
```bash
echo "=== Languages Audit ==="

# Check versions
go version 2>/dev/null || echo "Go: Not installed"
node --version 2>/dev/null || echo "Node.js: Not installed"
python3 --version 2>/dev/null || echo "Python 3: Not installed"
ruby --version 2>/dev/null || echo "Ruby: Not installed"
rustc --version 2>/dev/null || echo "Rust: Not installed"

# Package managers
uv --version 2>/dev/null || echo "UV: Not installed"
pnpm --version 2>/dev/null || echo "pnpm: Not installed"
cargo --version 2>/dev/null || echo "Cargo: Not installed"
```

### Cloud & DevOps Tools
```bash
echo "=== Cloud Tools Audit ==="

# Container tools
docker --version 2>/dev/null || echo "Docker: Not installed"
docker-compose --version 2>/dev/null || echo "Docker Compose: Not installed"

# Kubernetes
kubectl version --client --short 2>/dev/null || echo "kubectl: Not installed"
helm version --short 2>/dev/null || echo "Helm: Not installed"
k9s version --short 2>/dev/null || echo "k9s: Not installed"

# Cloud CLIs
aws --version 2>/dev/null || echo "AWS CLI: Not installed"
gcloud --version 2>/dev/null | head -1 || echo "gcloud: Not installed"
az --version 2>/dev/null | head -1 || echo "Azure CLI: Not installed"

# CloudFlare
wrangler --version 2>/dev/null || echo "Wrangler: Not installed"
```

### Development Tools
```bash
echo "=== Development Tools Audit ==="

# Editors
nvim --version 2>/dev/null | head -1 || echo "Neovim: Not installed"
code --version 2>/dev/null | head -1 || echo "VS Code: Not installed"

# Terminal
tmux -V 2>/dev/null || echo "tmux: Not installed"

# GitHub
gh --version 2>/dev/null | head -1 || echo "GitHub CLI: Not installed"

# Build tools
make --version 2>/dev/null | head -1 || echo "Make: Not installed"
cmake --version 2>/dev/null | head -1 || echo "CMake: Not installed"
```

## Package Manager Health

### Homebrew (macOS)
```bash
echo "=== Homebrew Health ==="
brew --version
brew doctor
brew outdated
echo "Installed packages: $(brew list | wc -l)"
echo "Cask packages: $(brew list --cask | wc -l)"
```

### APT (Debian/Ubuntu)
```bash
echo "=== APT Health ==="
apt --version
apt list --upgradable 2>/dev/null | tail -n +2
echo "Installed packages: $(dpkg -l | grep '^ii' | wc -l)"
```

## Storage Analysis

```bash
echo "=== Storage Usage ==="

# Homebrew (macOS)
if command -v brew >/dev/null 2>&1; then
    echo "Homebrew: $(du -sh $(brew --prefix) 2>/dev/null | cut -f1)"
fi

# Node modules (global)
if [ -d ~/.npm ]; then
    echo "npm cache: $(du -sh ~/.npm 2>/dev/null | cut -f1)"
fi

# Go modules
if [ -d ~/go/pkg ]; then
    echo "Go packages: $(du -sh ~/go/pkg 2>/dev/null | cut -f1)"
fi

# Python
if [ -d ~/.cache/pip ]; then
    echo "pip cache: $(du -sh ~/.cache/pip 2>/dev/null | cut -f1)"
fi
```

## Outdated Packages

```bash
# Homebrew
brew outdated

# npm global
npm outdated -g

# pip
pip list --outdated 2>/dev/null
```

## Generate Report

```bash
# Save audit to file
{
    echo "Tool Audit Report - $(date)"
    echo "========================="
    versions
    echo ""
    echo "=== Missing Tools ==="
    tools-missing 2>/dev/null || echo "Automation framework not available"
} > ~/tool-audit-$(date +%Y%m%d).txt

echo "Report saved to ~/tool-audit-$(date +%Y%m%d).txt"
```
