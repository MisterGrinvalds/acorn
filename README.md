# Component-Based Dotfiles System

A **modular, cross-shell compatible dotfiles system** for macOS and Linux. Works seamlessly with both **bash** and **zsh** using a component-based architecture with automatic dependency resolution.

## Key Features

### Cross-Shell Compatibility
- **Shell Portable**: Works with bash and zsh automatically
- **Auto-Detection**: Automatically detects and configures for your shell
- **Universal Prompts**: Beautiful git status integration with color coding

### XDG Compliant
- **Modern Standards**: Follows XDG Base Directory Specification
- **Clean Home**: Configuration organized in `~/.config/`
- **Portable**: Easy to migrate between machines

### Component Architecture
- **Self-Describing**: Each component has metadata for dependencies
- **Dependency Resolution**: Topological sort ensures correct load order
- **Optional Tools**: Components gracefully degrade when tools are missing

## Quick Start

```bash
# Clone and install
git clone <repo-url> ~/.config/dotfiles
cd ~/.config/dotfiles

# Interactive installation (recommended)
./install.sh

# Or install everything automatically
./install.sh --auto
```

**[Complete Installation Guide](docs/INSTALL.md)**

The installer will:
- Install cross-shell compatible dotfiles
- Configure secrets management
- Setup development tools (Git, Go, Node.js, Python)
- Run validation tests to ensure everything works

## Architecture

```
core/                     # Core bootstrap system
├── bootstrap.sh         # Main entry point
├── discovery.sh         # Shell/platform detection
├── xdg.sh              # XDG base directory setup
├── theme.sh            # Catppuccin Mocha color definitions
├── loader.sh           # Component loader with dependency resolution
└── sync.sh             # Dotfiles management functions

components/              # Feature components (each with optional config/)
├── shell/              # Core shell functions (cd, history, archive)
├── git/                # Git aliases, functions, and config
├── fzf/                # FZF integration with previews
├── python/             # Python/UV virtual environments and startup
├── node/               # Node.js/NVM/pnpm setup
├── go/                 # Go development environment
├── github/             # GitHub CLI integration
├── kubernetes/         # kubectl and helm aliases
├── secrets/            # Secret loading and validation
├── claude/             # Claude Code management and settings
├── ollama/             # Ollama local AI models
├── huggingface/        # Hugging Face integration
├── ssh/                # SSH client configuration
├── tmux/               # Tmux session manager and config
├── vscode/             # VS Code helpers and settings
├── ghostty/            # Ghostty terminal configuration
└── ...                 # And more (karabiner, conda, iterm2, etc.)
```

## Shell Experience

- **Catppuccin Mocha theme** with git status integration
- **Smart completions** for bash and zsh
- **Enhanced navigation** with auto-listing and fuzzy finding
- **History management** with intelligent search

## Key Functions

### Development
```bash
# Python
mkvenv myproject        # Create virtual environment
venv myproject          # Activate virtual environment
fastapi_env             # Setup FastAPI development

# Go
gonew myproject         # Initialize Go project
gotest                  # Run tests
gobuildall              # Build for multiple platforms

# Node.js
nvm_setup               # Install NVM and latest LTS
node_init               # Initialize TypeScript project
nclean                  # Clean and reinstall node_modules
```

### GitHub
```bash
quickpr                 # Push branch and create PR
newrepo myproject       # Create GitHub repository
gitcleanup              # Remove merged branches
```

### Kubernetes
```bash
kuse                    # Switch kubectl context
knsuse                  # Switch namespace
kpods                   # List pods
klf mypod               # Follow pod logs
```

### AI/ML
```bash
ollama_status           # Check installation and models
ollama_chat model msg   # Quick AI chat
ollama_code lang desc   # Generate code
```

## Dotfiles Management

```bash
dotfiles_inject     # Install bootstrap files
dotfiles_eject      # Remove all injected config
dotfiles_update     # Git pull + reload
dotfiles_reload     # Reload without restart
dotfiles_status     # Show current state
dotfiles_link_configs # Symlink app configs
```

## Shell Compatibility

The system automatically detects your shell and configures appropriately:

| Feature | Bash | Zsh | Notes |
|---------|------|-----|-------|
| Prompts | `PS1` | `PROMPT` | Shell-specific escaping |
| History | `HISTCONTROL` | `setopt` commands | Different syntax |
| Completion | bash-completion | compinit | Native systems |
| Colors | `\001...\002` | `%{...%}` | Proper escaping |

## Testing & Validation

```bash
# Quick validation tests
make test-quick

# Component tests
make test-components
make component-status

# Comprehensive test suite
make test-comprehensive

# Check what tools are installed
make test-required-tools
```

## Third-Party Tools

**macOS (via Homebrew):**
- **Core**: bash-completion, fd, fzf, git, gh (GitHub CLI)
- **Languages**: go, node, npm, python3
- **AI/ML**: ollama
- **Kubernetes**: kubectl, helm, k9s
- **Editors**: tmux, neovim

**Linux (via apt/yum):**
- **Core**: bash-completion, fd-find, fzf, git
- **Languages**: python3, python3-pip, nodejs, npm
- **AI/ML**: Manual Ollama installation

## Documentation

- **[Installation Guide](docs/INSTALL.md)** - Complete setup instructions
- **[Architecture Guide](CLAUDE.md)** - Design and development

## Getting Started

1. **[Read the Installation Guide](docs/INSTALL.md)** for platform-specific instructions
2. **Run `make test-quick`** to validate your setup
3. **Explore components** with `make component-list`

---

**Built for developers** - A lightweight, modular dotfiles system that grows with your needs while maintaining simplicity and portability.
