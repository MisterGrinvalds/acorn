# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture Overview

This is a **shell-portable dotfiles system** compatible with both **bash** and **zsh** on macOS and Linux. The system implements a modular, XDG-compliant configuration approach with automatic shell detection.

### Key Components:
- **Main entry point**: `shell/init.sh` - unified loader that sources all modules
- **Shell modules**: `shell/` - discovery, XDG setup, environment, aliases, prompt, completions
- **Function modules**: `functions/` - organized by domain (core, dev, cloud, ai)
- **App configs**: `config/` - git, ssh, python, and other tool configurations
- **Installation script**: `install.sh` - deploys dotfiles and creates shell bootstrap files

### Shell Compatibility Features:
- **Automatic shell detection** - `shell/discovery.sh` sets `$CURRENT_SHELL` and `$CURRENT_PLATFORM`
- **Shell-specific prompts** - Bash uses `PS1`, zsh uses `PROMPT` with proper escaping
- **Portable git integration** - Uses `git symbolic-ref` instead of bash regex
- **History management** - Configures shell-appropriate history in XDG-compliant locations
- **Completion systems** - Loads bash-completion or zsh compinit as needed

The system follows the XDG Base Directory specification with `$DOTFILES_ROOT` as the repo location.

## Directory Structure

```
shell/                    # Shell initialization modules
├── init.sh              # Main entry point (load order controller)
├── discovery.sh         # Shell/platform detection
├── xdg.sh              # XDG base directory setup
├── environment.sh       # Core environment variables
├── environment.darwin.sh # macOS-specific
├── environment.linux.sh # Linux-specific
├── aliases.sh          # Shell aliases
├── options.sh          # shopt/setopt options
├── secrets.sh          # Secrets loader
├── completions.sh      # Completion systems + FZF
└── prompt.sh           # Git-aware prompt

functions/               # Function modules by domain
├── core/               # Essential utilities
│   ├── cd.sh          # Enhanced cd
│   ├── history.sh     # History functions
│   ├── archive.sh     # mktar, mkzip
│   ├── utils.sh       # bash_as, misc
│   ├── tools.sh       # Tool management
│   ├── tmux.sh        # Tmux helpers
│   └── inject.sh      # Dotfiles management
├── dev/                # Development tools
│   ├── python.sh      # mkvenv, venv, rmenv
│   ├── golang.sh      # Go helpers
│   ├── github.sh      # GitHub functions
│   └── vscode.sh      # VS Code functions
├── cloud/              # Cloud/DevOps
│   ├── kubernetes.sh  # K8s functions
│   ├── automation.sh  # Automation helpers
│   └── secrets.sh     # Secret management
└── ai/                 # AI/ML tools
    ├── ollama.sh      # Ollama integration
    ├── huggingface.sh # Hugging Face
    └── aliases.sh     # AI command aliases

config/                  # Application configs
├── git/                # Git configuration
├── ssh/                # SSH config
├── python/             # Python startup
├── r/                  # R profile
├── conda/              # Conda config
├── wget/               # Wget config
└── karabiner/          # macOS keyboard

secrets/                 # Secure storage (gitignored)
docs/                    # Documentation
.automation/             # Automation CLI framework
```

## Loading Sequence

Configuration loads in this order (see `shell/init.sh`):
1. **`discovery.sh`** - Detects shell type and platform
2. **`xdg.sh`** - Sets XDG base directories
3. **`environment.sh`** - Core environment variables (sources platform-specific)
4. **`secrets.sh`** - Loads credentials silently
5. **`options.sh`** - Shell options (shopt/setopt)
6. **`aliases.sh`** - Command aliases
7. **`functions/**/*.sh`** - All function modules
8. **`completions.sh`** - Completion systems and FZF
9. **`prompt.sh`** - Git-aware prompt
10. **`~/.config/shell/local.sh`** - Local overrides (optional)

## Installation and Setup

### Quick Install
```bash
# Clone and install
git clone <repo-url> ~/.config/dotfiles
cd ~/.config/dotfiles
./install.sh
```

The installer creates bootstrap files that source the main init.sh:
- `~/.bashrc` - Sources `$DOTFILES_ROOT/shell/init.sh`
- `~/.zshrc` - Sources `$DOTFILES_ROOT/shell/init.sh`

### Dotfiles Management Functions
```bash
dotfiles_inject     # Install bootstrap files
dotfiles_eject      # Remove all injected config
dotfiles_update     # Git pull + reload
dotfiles_reload     # Reload without restart
dotfiles_status     # Show current state
dotfiles_link_configs # Symlink app configs
```

### Third-Party Tools
**macOS (via Homebrew):**
- bash-completion, fd, fzf, xclip, xquartz

**Linux (via apt/yum):**
- bash-completion, fd-find, fzf

## Custom Functions

### Core (functions/core/)
- **`cd()`** - Enhanced cd that runs `ll` after directory change
- **`h(pattern)`** - History search with grep
- **`mktar(path)`** - Create .tar.gz archive
- **`mkzip(path)`** - Create .zip archive
- **`bash_as(user)`** - Run bash shell as another user

### Development (functions/dev/)
- **`mkvenv([name])`** - Create virtual environment
- **`venv([name])`** - Activate virtual environment
- **`rmenv([pattern])`** - Remove environment variables

### Cloud (functions/cloud/)
- **`load_secrets`** - Load secrets into environment
- **`validate_secrets`** - Validate API keys
- **Kubernetes helpers** - Pod logs, context switching

### AI (functions/ai/)
- **`ollama_setup`** - Install and configure Ollama
- **`ollama_chat`** - Interactive AI chat
- **`hf_setup`** - Setup Hugging Face environment

## Testing and Development

### Test with Makefile
```bash
make test-quick          # Syntax and basic tests
make test-dotfiles       # Dotfiles functionality
make test-syntax         # All shell script syntax
make test-comprehensive  # Full test suite
```

### Manual Testing
```bash
# Test in clean shell
bash --noprofile --norc
export DOTFILES_ROOT="$PWD"
source shell/init.sh

# Verify detection
echo "Shell: $CURRENT_SHELL"
echo "Platform: $CURRENT_PLATFORM"

# Test functions
mkvenv test-env
h git
```

### Key Test Cases
1. **Shell detection** works in bash and zsh
2. **Git prompt** shows correct colors and branch info
3. **XDG directories** created correctly
4. **Custom functions** work in both shells
5. **Completions** load without errors

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

### Development Paths
- **Python**: `PYTHONSTARTUP`, `IPYTHONDIR`
- **Node.js**: `NVM_DIR` (XDG-compliant)
- **Neovim**: `NEOVIM_VIRTUALENV`, `VIM_PLUGGED`

## Troubleshooting

### Common Issues
1. **Shell not detected**: Check `$BASH_VERSION` / `$ZSH_VERSION`
2. **Prompt broken**: Verify shell-specific color escaping
3. **Functions missing**: Ensure `functions/` is sourced
4. **XDG not set**: Check `shell/xdg.sh` sourcing

### Debug Commands
```bash
# Check detection
echo "Shell: $CURRENT_SHELL, Platform: $CURRENT_PLATFORM"

# Test git functions
git_branch; git_color

# Verify XDG
env | grep XDG

# List loaded functions
declare -F | grep -E "(mkvenv|dotfiles)"
```
