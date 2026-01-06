---
description: Set up and manage Node.js versions with NVM
argument-hint: [action: install|use|list|default]
allowed-tools: Read, Bash
---

## Task

Help the user set up and manage Node.js versions with NVM (Node Version Manager).

## Actions

Based on `$ARGUMENTS`:

### install
Install NVM and Node:

```bash
# Using dotfiles function (recommended)
nvm_setup
# Installs NVM + latest LTS Node

# Manual NVM installation
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash

# Reload shell
source ~/.bashrc  # or ~/.zshrc

# Install latest LTS
nvm install --lts
```

### use
Switch Node versions:

```bash
# Use specific version
nvm use 20
nvmu 20  # dotfiles alias

# Use version from .nvmrc
nvm use

# Use latest LTS
nvm use --lts
nvm_latest  # dotfiles function
```

### list
List available versions:

```bash
# Installed versions
nvm ls
nvml  # dotfiles alias

# Remote versions available
nvm ls-remote --lts
nvmr  # dotfiles alias
```

### default
Set default version:

```bash
# Set default
nvm alias default 20

# Use default
nvm use default
```

## Version Management

### Installing Specific Versions
```bash
# Latest LTS
nvm install --lts

# Specific major version
nvm install 20
nvmi 20  # dotfiles alias

# Specific version
nvm install 20.10.0

# Latest of major version
nvm install 20 --latest
```

### Project-Specific Version

Create `.nvmrc` file:
```bash
# Record current version
node --version > .nvmrc

# Or specify manually
echo "20" > .nvmrc
```

Use it:
```bash
cd my-project
nvm use
# Reads version from .nvmrc
```

### Auto-switching (optional)

Add to shell config for automatic version switching:

**Bash** (`~/.bashrc`):
```bash
cd() {
  builtin cd "$@"
  if [ -f .nvmrc ]; then
    nvm use
  fi
}
```

**Zsh** (`~/.zshrc`):
```bash
autoload -U add-zsh-hook
load-nvmrc() {
  if [ -f .nvmrc ]; then
    nvm use
  fi
}
add-zsh-hook chpwd load-nvmrc
load-nvmrc
```

## Uninstalling Versions

```bash
# Uninstall specific version
nvm uninstall 18

# List to see what's installed
nvm ls
```

## Troubleshooting

### "nvm: command not found"
```bash
# Ensure NVM is loaded
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"

# Reload shell
source ~/.bashrc
```

### Global packages lost after version switch
```bash
# Reinstall globals from another version
nvm reinstall-packages <previous-version>

# Or install fresh
npm install -g typescript pnpm
```

### Slow shell startup
```bash
# Lazy load NVM (add to shell config)
# Only loads NVM when first used
```

## XDG Location

The dotfiles configure NVM at:
```
$XDG_DATA_HOME/nvm  (~/.local/share/nvm)
```

## Dotfiles Functions & Aliases

- `nvm_setup` - Install NVM and latest LTS
- `nvm_latest` - Install and use latest LTS
- `nvml` - List installed versions
- `nvmr` - List remote versions
- `nvmu` - Use version
- `nvmi` - Install version
