# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture Overview

This is a **component-based dotfiles system** compatible with both **bash** and **zsh** on macOS and Linux. The system implements a modular, XDG-compliant configuration approach with automatic shell detection and dependency resolution.

### Key Components:

- **Main entry point**: `core/bootstrap.sh` - unified loader that initializes the component system
- **Core modules**: `core/` - discovery, XDG setup, theme, loader, sync
- **Components**: `components/` - self-contained feature modules with metadata
- **App configs**: `config/` - git, ssh, python, and other tool configurations
- **Installation script**: `install.sh` - deploys dotfiles and creates shell bootstrap files

### Component Architecture:
- **Self-describing**: Each component has a `component.yaml` with metadata
- **Dependency resolution**: Topological sort ensures correct load order
- **Optional tools**: Components gracefully degrade when tools are missing
- **XDG compliance**: All paths follow XDG Base Directory specification

## Directory Structure

```
core/                     # Core bootstrap system
├── bootstrap.sh         # Main entry point
├── discovery.sh         # Shell/platform detection
├── xdg.sh              # XDG base directory setup
├── theme.sh            # Catppuccin Mocha color definitions
├── loader.sh           # Component loader with dependency resolution
└── sync.sh             # Dotfiles management functions

components/              # Feature components (each with optional config/)
├── _template/          # Template for new components
├── shell/              # Core shell functions (cd, history, archive)
├── git/                # Git aliases, functions, and config files
├── fzf/                # FZF integration with previews
├── tools/              # Tool management utilities
├── tmux/               # Tmux session management and config
├── python/             # Python/UV virtual environments and startup
├── node/               # Node.js/NVM/pnpm setup
├── go/                 # Go development environment
├── github/             # GitHub CLI integration
├── vscode/             # VS Code helpers and settings
├── database/           # Database tool aliases
├── kubernetes/         # kubectl and helm aliases
├── secrets/            # Secret loading and validation
├── claude/             # Claude Code management and settings
├── ollama/             # Ollama local AI models
├── huggingface/        # Hugging Face integration
├── ssh/                # SSH client configuration
├── conda/              # Conda package manager config
├── karabiner/          # macOS keyboard customization
├── ghostty/            # Ghostty terminal configuration
├── iterm2/             # iTerm2 terminal configuration
├── intellij/           # IntelliJ IDE configuration
├── wget/               # Wget download utility config
└── r/                  # R statistical programming config

secrets/                 # Secure storage (gitignored)
```

### Component Structure

Each component follows this structure:
```
components/<name>/
├── component.yaml      # Metadata: name, dependencies, tools, config
├── env.sh             # Environment variables
├── aliases.sh         # Command aliases
├── functions.sh       # Shell functions
├── completions.sh     # Tab completions
└── config/            # Configuration files (optional)
    └── <config-files> # App-specific config files
```

Example `component.yaml`:
```yaml
name: python
description: Python development with UV package manager
version: 1.0.0
shells: [bash, zsh]
platforms: [darwin, linux]

dependencies:
  components:
    - shell
  tools:
    required: []
    optional:
      - uv
      - python3

# Configuration files to deploy (optional)
config:
  files:
    - source: config/startup.py
      target: ~/.config/python/startup.py
      method: symlink
  directories:
    - target: ~/.config/python
```

## Loading Sequence

Configuration loads in this order (see `core/bootstrap.sh`):
1. **`core/discovery.sh`** - Detects shell type and platform
2. **`core/xdg.sh`** - Sets XDG base directories
3. **`core/theme.sh`** - Loads Catppuccin Mocha colors
4. **`core/loader.sh`** - Discovers and loads components in dependency order
5. **`core/sync.sh`** - Dotfiles management functions
6. **`~/.config/shell/local.sh`** - Local overrides (optional)

Component loading order is determined by topological sort of dependencies.

## Installation and Setup

### Quick Install
```bash
# Clone and install
git clone <repo-url> ~/.config/dotfiles
cd ~/.config/dotfiles
./install.sh
```

The installer creates bootstrap files that source core/bootstrap.sh:
- `~/.bashrc` - Sources `$DOTFILES_ROOT/core/bootstrap.sh`
- `~/.zshrc` - Sources `$DOTFILES_ROOT/core/bootstrap.sh`
- `~/.bash_profile` - Sources ~/.bashrc for login shells

### Dotfiles Management Functions
```bash
dotfiles_inject     # Install bootstrap files
dotfiles_eject      # Remove all injected config
dotfiles_update     # Git pull + reload
dotfiles_reload     # Reload without restart
dotfiles_status     # Show current state
dotfiles_audit      # Check XDG compliance
dotfiles_link_configs   # Symlink app configs
dotfiles_unlink_configs # Remove symlinks
```

Convenience aliases: `df-inject`, `df-eject`, `df-update`, `df-reload`, `df-status`, `df-audit`

### Component Management (Makefile)
```bash
make component-list      # List all components
make component-status    # Show component health
make component-new NAME=foo  # Create from template
make component-validate  # Validate all component.yaml
make test-components     # Test component loading
```

## Key Functions by Component

### shell
- **`cd()`** - Enhanced cd that runs `ll` after directory change
- **`h(pattern)`** - History search with grep
- **`mktar(path)`** - Create .tar.gz archive
- **`mkzip(path)`** - Create .zip archive
- **`bash_as(user)`** - Run bash shell as another user

### python
- **`mkvenv([name])`** - Create virtual environment (uses UV if available)
- **`venv([name])`** - Activate virtual environment
- **`dvenv()`** - Deactivate virtual environment
- **`fastapi_env()`** - Setup FastAPI development environment

### node
- **`nvm_setup()`** - Install NVM and latest LTS Node
- **`node_init()`** - Initialize TypeScript project
- **`nclean()`** - Remove and reinstall node_modules

### go
- **`gonew(name)`** - Initialize new Go project
- **`gotest([pattern])`** - Run tests with optional filter
- **`gotestcover()`** - Generate coverage report
- **`gobuildall()`** - Build for multiple platforms

### github
- **`quickpr()`** - Push branch and create PR
- **`newrepo(name)`** - Create GitHub repository
- **`gitcleanup()`** - Remove merged branches

### kubernetes
- **`kuse([context])`** - Switch kubectl context
- **`knsuse([ns])`** - Switch namespace
- **`kpods([filter])`** - List pods with optional filter
- **`klf(pod)`** - Follow pod logs

### claude
- **`claude_stats()`** - View usage statistics
- **`claude_tokens()`** - View token usage by model
- **`claude_permissions()`** - View/manage permissions
- **`claude_mcp()`** - List MCP servers
- **`claude_settings()`** - View settings (global/local/config)
- **`claude_info()`** - Show Claude Code info summary
- **`claude_help()`** - Show all available functions

### ollama
- **`ollama_status()`** - Check installation and models
- **`ollama_chat(model, prompt)`** - Quick AI chat
- **`ollama_code(lang, desc)`** - Generate code

## Testing and Development

### Test with Makefile
```bash
make test-quick          # Syntax and basic tests
make test-dotfiles       # Dotfiles functionality
make test-syntax         # All shell script syntax
make test-components     # Component loading tests
make test-comprehensive  # Full test suite
```

### Manual Testing
```bash
# Test in clean shell
bash --noprofile --norc
export DOTFILES_ROOT="$PWD"
source core/bootstrap.sh

# Verify detection
echo "Shell: $CURRENT_SHELL"
echo "Platform: $CURRENT_PLATFORM"

# Test functions
mkvenv test-env
h git
```

### Creating New Components
```bash
make component-new NAME=mycomponent
# Edit components/mycomponent/component.yaml
# Add functions to components/mycomponent/functions.sh
```

## Environment Variables

### XDG Directories
- **`DOTFILES_ROOT`** - Repository location
- **`XDG_CONFIG_HOME`** - ~/.config
- **`XDG_DATA_HOME`** - ~/.local/share
- **`XDG_CACHE_HOME`** - ~/.cache
- **`XDG_STATE_HOME`** - ~/.local/state

### Shell Detection
- **`CURRENT_SHELL`** - bash/zsh/unknown
- **`CURRENT_PLATFORM`** - darwin/linux/unknown
- **`IS_INTERACTIVE`** - true/false
- **`IS_LOGIN_SHELL`** - true/false

### Theme Colors (Catppuccin Mocha)
- **`CLR_*`** - Color codes (CLR_RED, CLR_GREEN, CLR_BLUE, etc.)
- **`CLR_RESET`** - Reset color

## Troubleshooting

### Common Issues
1. **Shell not detected**: Check `$BASH_VERSION` / `$ZSH_VERSION`
2. **Component not loading**: Check `component.yaml` syntax with `yq`
3. **Functions missing**: Run `make component-status` to check health
4. **XDG not set**: Check `core/xdg.sh` sourcing

### Debug Commands
```bash
# Check detection
echo "Shell: $CURRENT_SHELL, Platform: $CURRENT_PLATFORM"

# List loaded components
make component-status

# Verify XDG
env | grep XDG

# List loaded functions
declare -F | grep -E "(mkvenv|dotfiles)"

# Check component dependencies
yq '.dependencies' components/python/component.yaml
```

## Claude Code Configuration

The `components/claude/config/settings.json` contains global Claude Code settings that are symlinked to `~/.claude/settings.json` during installation.

### Settings Location
- **Global settings**: `~/.claude/settings.json` (symlinked from `components/claude/config/settings.json`)
- **Project settings**: `.claude/settings.json` in any project directory
- **Local overrides**: `~/.claude/settings.local.json`

### Configuration Features
- **Permissions**: Pre-approved tool patterns (e.g., `Bash(git:*)`, `Bash(make:*)`)
- **Hooks**: Custom shell scripts for various lifecycle events
- **Status line**: Custom status bar integration
- **Preferences**: Auto-approve read operations, preferred shell

### Claude Component Functions
The `components/claude/` component provides shell functions for managing Claude Code:
```bash
claude_info          # Quick overview
claude_stats         # Usage statistics
claude_permissions   # View/manage permissions
claude_settings      # View settings files
claude_mcp           # MCP server management
```

## Legacy Migration

The system was migrated from a legacy `shell/` + `functions/` structure to the current component-based architecture. Legacy directories may still exist but are no longer used:

- `shell/` → Replaced by `core/`
- `functions/` → Replaced by `components/`

If you find references to legacy paths, update them to use the new structure.
