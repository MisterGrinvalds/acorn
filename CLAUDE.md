# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture Overview

This is a **shell-portable dotfiles system** compatible with both **bash** and **zsh** on macOS and Linux. The system implements a modular configuration approach with automatic shell detection.

### Key Components:
- **Main entry point**: `.bash_profile` - detects shell type and sources all configuration files
- **Configuration modules**: `.bash_profile.dir/` - environment, aliases, terminal customization
- **Custom utilities**: `.bash_tools/` - 10 custom shell functions for development workflow
- **Installation script**: `initialize.sh` - deploys dotfiles and creates appropriate shell links

### Shell Compatibility Features:
- **Automatic shell detection** - Sets `$CURRENT_SHELL` (bash/zsh/unknown)
- **Shell-specific prompts** - Bash uses `PS1`, zsh uses `PROMPT` with proper escaping
- **Portable git integration** - Uses case statements instead of bash regex
- **History management** - Configures shell-appropriate history settings
- **Completion systems** - Loads bash-completion or zsh compinit as needed

The system honors `$DOTFILES` (defaults to `$HOME`) and XDG Base Directory specification.

## Loading Sequence

Configuration loads in this order:
1. **Shell detection** - Sets `$CURRENT_SHELL` variable
2. **`environment.sh`** - Core environment variables and paths
3. **`aliases.sh`** - Shell-agnostic command aliases
4. **`terminal.sh`** - Shell-specific prompt and color configuration
5. **`.bash_tools/*`** - Custom utility functions (all files)

## Installation and Setup

### Quick Install
```bash
# Clone and install
git clone <repo-url> ~/.dotfiles
cd ~/.dotfiles
./initialize.sh
```

The installer will:
- Copy files to `$HOME` using rsync
- Create `~/.bashrc` → `~/.bash_profile` symlink
- Create `~/.zshrc` → `~/.bash_profile` symlink  
- Optionally install third-party tools via package manager

### Manual Setup
```bash
# For bash users
echo 'source ~/.bash_profile' >> ~/.bashrc

# For zsh users  
echo 'source ~/.bash_profile' >> ~/.zshrc
```

### Third-Party Tools
**macOS (via Homebrew):**
- bash-completion, fd, fzf, xclip, xquartz

**Linux (via apt/yum):**
- bash-completion, fd-find, fzf

## Configuration Structure

### Environment Files
- **`environment.sh`** - Base environment variables
- **`environment_darwin.sh`** - macOS-specific settings
- **`environment_linux-gnu.sh`** - Linux-specific settings

### Shell-Specific Features
- **History**: Bash uses `HISTCONTROL`, zsh uses `setopt` commands
- **Completion**: Loads appropriate completion system per shell
- **Prompts**: Shell-aware color escaping and variable substitution
- **FZF Integration**: Loads shell-specific key bindings and completion

## Custom Functions (.bash_tools/)

### Python Development
- **`mkvenv([name])`** - Create virtual environment (local .venv or named in $ENVS_LOCATION)
- **`venv([name])`** - Activate virtual environment
- **`rmenv([pattern])`** - Remove environment variables matching pattern

### Navigation & Files  
- **`cd()`** - Enhanced cd that runs `ll` after directory change
- **`h(pattern)`** - History search with grep
- **`mktar(path)`** - Create .tar.gz archive
- **`mkzip(path)`** - Create .zip archive

### System Tools
- **`bash-as(user)`** - Run bash shell as another user

## Testing and Development

### Manual Testing Approach
```bash
# Test in clean shell
bash --noprofile --norc
source ~/.bash_profile

# Verify shell detection
echo "Current shell: $CURRENT_SHELL"

# Test git prompt in different repo states
cd /path/to/git/repo
# Verify colors: clean (green), dirty (yellow), ahead (red)

# Test custom functions
mkvenv test-env
venv test-env
cd /some/path  # Should auto-run ll
h git          # Should search history
```

### Key Test Cases
1. **Shell detection** works in bash and zsh
2. **Git prompt** shows correct colors and branch info
3. **History** settings work per shell
4. **Completion** loads without errors
5. **Custom functions** work in both shells
6. **FZF integration** provides key bindings

## Environment Variables

### Development Paths
- **Python**: `PYTHONSTARTUP`, `IPYTHONDIR`, `JUPYTER_CONFIG_DIR`
- **Node.js**: `NVM_DIR` with automatic loading
- **R**: `R_PROFILE`, `R_PROFILE_USER` 
- **Neovim**: `NEOVIM_VIRTUALENV`, `VIM_PLUGGED`

### XDG Compliance
- **`XDG_CONFIG_HOME`** - Configuration directory
- **`XDG_CACHE_HOME`** - Cache directory  
- **`XDG_DATA_HOME`** - Data directory

### Security
- **Secrets loading**: Automatically sources files from `~/.secrets/`
- **No hardcoded credentials**: All sensitive data externalized

## Troubleshooting

### Common Issues
1. **Shell not detected**: Check `$BASH_VERSION` / `$ZSH_VERSION` variables
2. **Prompt broken**: Verify color escaping for your shell
3. **Completion not working**: Check if completion files exist and are sourced
4. **Git colors wrong**: Test `git status` output format changes

### Debug Commands
```bash
# Check shell detection
echo "Shell: $CURRENT_SHELL"

# Test git functions
git_branch; git_color

# Verify environment loading
env | grep -E "(HIST|XDG|DOTFILES)"
```