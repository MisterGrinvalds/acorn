---
description: Explain system tools, package managers, and version management
argument-hint: [topic]
allowed-tools: Read, Glob, Grep
---

## Task

Explain the requested topic about system tools and package management. If no specific topic provided, give an overview of the tools component.

## Topics

### Package Managers
- **homebrew** - macOS/Linux package manager
- **apt** - Debian/Ubuntu package manager
- **dnf** - Fedora/RHEL package manager
- **pacman** - Arch Linux package manager

### Version Management
- **versions** - How version numbers work (semver)
- **nvm** - Node version manager
- **pyenv** - Python version manager
- **goenv** - Go version manager

### Tool Categories
- **system** - Core system tools (git, curl, jq)
- **languages** - Programming language runtimes
- **cloud** - Cloud CLI tools (aws, kubectl, gcloud)
- **development** - Dev tools (docker, gh, neovim)

### Operations
- **installing** - How to install tools
- **updating** - Keeping tools current
- **cleanup** - Removing old versions
- **troubleshooting** - Common issues

## Context

@components/tools/component.yaml
@components/tools/functions.sh
@components/tools/aliases.sh

## Response Format

When explaining a topic:

1. **Definition** - What it is in simple terms
2. **How it works** - Technical details
3. **Examples** - Practical usage
4. **Dotfiles integration** - Available functions/aliases
5. **Best practices** - Common patterns

## Quick Reference

### Dotfiles Functions
- `tools_status` - Check tools via automation framework
- `list_tools` - List all managed tools
- `check_tools` - Check tool versions
- `update_tools` - Update tools
- `missing_tools` - Show missing tools
- `install_tool <tool>` - Install specific tool
- `quick_versions` - Check common tool versions
- `smart_update` - Auto-detect and update system
- `which_enhanced <cmd>` - Enhanced which command

### Key Aliases
| Alias | Function |
|-------|----------|
| `tools` | tools_status |
| `tools-list` | list_tools |
| `tools-check` | check_tools |
| `tools-update` | update_tools |
| `tools-missing` | missing_tools |
| `tools-install` | install_tool |
| `versions` | quick_versions |
| `system-update` | smart_update |
| `whichx` | which_enhanced |

### Package Manager Detection
The `smart_update` function automatically detects:
- `brew` - Homebrew on macOS/Linux
- `apt-get` - Debian/Ubuntu
- `dnf` - Fedora/RHEL
- `pacman` - Arch Linux
