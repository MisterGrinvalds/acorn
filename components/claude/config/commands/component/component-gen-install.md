---
description: Generate installation scripts for a component
argument_hints:
  - tmux
  - go
  - python
  - node
  - kubernetes
---

Generate install scripts for: $ARGUMENTS

## Instructions

Create installation scripts and package manifests for the component.

### 1. Research Package Names

Identify the correct package names for each package manager:

| Manager | Example |
|---------|---------|
| brew | `tmux`, `go`, `python@3.11`, `node` |
| apt | `tmux`, `golang-go`, `python3`, `nodejs` |
| dnf | Similar to apt, may vary |
| pacman | `tmux`, `go`, `python`, `nodejs` |

### 2. Generate install/brew.yaml

```yaml
# Homebrew packages for $ARGUMENTS
formulas:
  - name: <package>
    description: <what it provides>

casks: []  # Desktop apps if needed

taps: []   # Custom taps if needed
```

### 3. Generate install/apt.yaml

```yaml
# APT packages for $ARGUMENTS
packages:
  - name: <package>
    description: <what it provides>

ppas: []  # Additional PPAs if needed
```

### 4. Generate install/install.sh

```bash
#!/bin/bash
# Installation script for $ARGUMENTS
# Supports: macOS (brew), Ubuntu/Debian (apt), Fedora (dnf), Arch (pacman)

set -e

COMPONENT="$ARGUMENTS"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info() { echo -e "${GREEN}[INFO]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

# Detect package manager
detect_package_manager() {
    if command -v brew &>/dev/null; then
        echo "brew"
    elif command -v apt-get &>/dev/null; then
        echo "apt"
    elif command -v dnf &>/dev/null; then
        echo "dnf"
    elif command -v pacman &>/dev/null; then
        echo "pacman"
    else
        error "No supported package manager found"
    fi
}

# Install with Homebrew
install_brew() {
    info "Installing with Homebrew..."
    # Read from brew.yaml if yq is available
    if command -v yq &>/dev/null && [ -f "$SCRIPT_DIR/brew.yaml" ]; then
        yq -r '.formulas[].name' "$SCRIPT_DIR/brew.yaml" | while read -r pkg; do
            brew install "$pkg" || warn "Failed to install $pkg"
        done
    else
        # Fallback to hardcoded list
        brew install <package> || warn "Failed to install <package>"
    fi
}

# Install with APT
install_apt() {
    info "Installing with APT..."
    sudo apt-get update
    if command -v yq &>/dev/null && [ -f "$SCRIPT_DIR/apt.yaml" ]; then
        yq -r '.packages[].name' "$SCRIPT_DIR/apt.yaml" | while read -r pkg; do
            sudo apt-get install -y "$pkg" || warn "Failed to install $pkg"
        done
    else
        sudo apt-get install -y <package> || warn "Failed to install <package>"
    fi
}

# Install with DNF
install_dnf() {
    info "Installing with DNF..."
    if command -v yq &>/dev/null && [ -f "$SCRIPT_DIR/apt.yaml" ]; then
        yq -r '.packages[].name' "$SCRIPT_DIR/apt.yaml" | while read -r pkg; do
            sudo dnf install -y "$pkg" || warn "Failed to install $pkg"
        done
    else
        sudo dnf install -y <package> || warn "Failed to install <package>"
    fi
}

# Install with Pacman
install_pacman() {
    info "Installing with Pacman..."
    sudo pacman -S --noconfirm <package> || warn "Failed to install <package>"
}

# Post-install configuration
post_install() {
    info "Running post-install configuration..."
    # Add any post-install steps here
    # e.g., creating config directories, linking files
}

# Main
main() {
    info "Installing $COMPONENT..."

    PM=$(detect_package_manager)
    info "Detected package manager: $PM"

    case "$PM" in
        brew)   install_brew ;;
        apt)    install_apt ;;
        dnf)    install_dnf ;;
        pacman) install_pacman ;;
    esac

    post_install

    info "$COMPONENT installation complete!"
}

main "$@"
```

### 5. Add Tool-Specific Post-Install

For each component, add appropriate post-install steps:

| Component | Post-Install Steps |
|-----------|-------------------|
| tmux | Install TPM, create config directory |
| go | Set GOPATH, create go directories |
| python | Create virtualenv directory, install pip |
| node | Install NVM, configure npmrc |

### 6. Report

Output:
```
Generated Install Scripts: $ARGUMENTS
=====================================

Created:
  - components/$ARGUMENTS/install/brew.yaml
  - components/$ARGUMENTS/install/apt.yaml
  - components/$ARGUMENTS/install/install.sh

Packages:
  - brew: <package1>, <package2>
  - apt: <package1>, <package2>

Post-install steps:
  - <step1>
  - <step2>
```
