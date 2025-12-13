# Enhanced Bash Profile & Automation Framework

A **comprehensive development environment** with cross-shell compatible dotfiles, automation framework, secrets management, and cloud integration. Works seamlessly with both **bash** and **zsh** on macOS and Linux with modern development workflows and enterprise-grade tooling.

## Key Features

### Cross-Shell Compatibility
- **Shell Portable**: Works with bash and zsh automatically
- **Auto-Detection**: Automatically detects and configures for your shell
- **Universal Prompts**: Beautiful git status integration with color coding

### XDG Compliant
- **Modern Standards**: Follows XDG Base Directory Specification
- **Clean Home**: Configuration organized in `~/.config/`
- **Portable**: Easy to migrate between machines

### Automation Framework
- **Unified CLI**: Single `auto` command for all development workflows
- **Multi-Cloud Support**: AWS, Azure, DigitalOcean management
- **Container Orchestration**: Kubernetes, Helm, and Docker integration
- **Project Templates**: FastAPI, Go/Cobra CLI, TypeScript scaffolding

### Enterprise Secrets Management
- **API Key Management**: Secure storage for 30+ services
- **Interactive Setup**: Guided configuration wizards
- **Credential Validation**: Real-time authentication testing

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

The enhanced installer will:
- Install cross-shell compatible dotfiles
- Setup comprehensive automation framework
- Configure secrets management for 30+ services
- Install cloud CLI tools (AWS, Azure, DigitalOcean)
- Setup development tools (Git, Go, Node.js, Python)
- Run validation tests to ensure everything works

## Architecture

```
shell/                    # Shell initialization modules
├── init.sh              # Main entry point
├── discovery.sh         # Shell/platform detection
├── xdg.sh              # XDG base directory setup
├── environment.sh       # Core environment variables
├── aliases.sh          # Shell aliases
├── completions.sh      # Completion systems + FZF
└── prompt.sh           # Git-aware prompt

functions/               # Function modules by domain
├── core/               # Essential utilities (cd, history, archive)
├── dev/                # Development (python, golang, github)
├── cloud/              # Cloud/DevOps (kubernetes, secrets)
└── ai/                 # AI/ML (ollama, huggingface)

config/                  # Application configs
├── git/                # Git configuration
├── ssh/                # SSH config
├── python/             # Python startup
└── ...                 # Other tool configs

.automation/             # Automation CLI framework
├── auto                # Main CLI entry point
├── framework/          # Core framework (logging, utilities)
├── modules/            # Feature modules (dev, k8s, cloud, secrets)
└── secrets/            # Secure API key storage
```

## Shell Experience

- **Solarized prompts** with git status integration
- **Smart completions** for bash and zsh
- **Enhanced navigation** with auto-listing and fuzzy finding
- **History management** with intelligent search

## Automation Commands

```bash
# Development workflows
auto dev init python my-fastapi-app
auto dev init go my-cli --cobra
auto dev init typescript my-app

# AI/ML workflows
make ai-setup                    # Install Ollama + Hugging Face
make ai-chat                     # Interactive AI chat
make ollama-install              # Install Ollama only
make hf-setup                    # Setup Hugging Face only

# Kubernetes operations
auto k8s deploy my-app production
auto k8s logs my-pod --follow
auto k8s monitoring

# GitHub integration
auto github repo create my-project
auto github pr create "New feature"

# Cloud management
auto cloud status
auto aws ec2 create web-server
auto azure vm create my-vm

# Secrets management
auto secrets setup
auto secrets validate
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

# API keys and authentication
make test-api-keys
make test-auth-status

# AI/ML functionality
make ai-test
make ai-status

# Comprehensive test suite
make test-comprehensive

# Check what tools are installed
make test-required-tools
```

## Supported Services

**Cloud Providers**: AWS, Azure, DigitalOcean
**Development**: GitHub, GitLab, Docker Hub
**AI/ML Platforms**: Ollama, Hugging Face, OpenAI
**Databases**: PostgreSQL, MySQL, MongoDB, Redis
**Monitoring**: DataDog, New Relic, Grafana
**Communication**: Slack, Discord
**CI/CD**: Jenkins, CircleCI, Travis CI

## Third-Party Tools

**macOS (via Homebrew):**
- **Core**: bash-completion, fd, fzf, git, gh (GitHub CLI)
- **Languages**: go, node, npm, python3
- **AI/ML**: ollama, python3-pip (for transformers, torch)
- **Kubernetes**: kubectl, helm, k9s
- **Editors**: tmux, neovim, visual-studio-code

**Linux (via apt/yum):**
- **Core**: bash-completion, fd-find, fzf, git
- **Languages**: python3, python3-pip, nodejs, npm
- **AI/ML**: Manual Ollama installation, pip packages

## Documentation

- **[Installation Guide](docs/INSTALL.md)** - Complete setup instructions
- **[Secrets Management](docs/automation/SECRETS.md)** - API keys and security
- **[Automation Framework](docs/automation/README.md)** - CLI tools and workflows
- **[Tools Management](docs/automation/TOOLS.md)** - External tools management
- **[Architecture Guide](CLAUDE.md)** - Design and development

## Getting Started

1. **[Read the Installation Guide](docs/INSTALL.md)** for platform-specific instructions
2. **[Configure Secrets](docs/automation/SECRETS.md)** for your cloud providers and services
3. **[Explore Automation](docs/automation/README.md)** to streamline your development workflow
4. **Run `make test-quick`** to validate your setup

---

**Built for modern development teams** - This comprehensive framework grows with your needs, supports enterprise security requirements, and streamlines cloud-native development workflows while maintaining the simplicity and portability of traditional dotfiles.
