# Session Context - Bash Profile Repository

## Repository Purpose

This is a **shell-portable dotfiles system** that provides a modular, component-isolated configuration for bash and zsh on macOS and Linux. The repository was restructured to follow a component-based architecture with XDG Base Directory compliance and drift prevention.

## Architecture Overview

### Directory Structure

```
bash-profile/
├── core/                           # Framework core (always loaded)
│   ├── bootstrap.sh                # Entry point (~/.bashrc sources this)
│   ├── discovery.sh                # Shell/platform detection
│   ├── xdg.sh                      # XDG directory management
│   ├── theme.sh                    # Catppuccin Mocha color definitions
│   ├── loader.sh                   # Component discovery & loading
│   └── sync.sh                     # Git sync & drift prevention
│
├── components/                     # Self-contained tool modules
│   ├── _template/                  # Template for new components
│   ├── shell/                      # Core shell (prompt, options, history)
│   ├── git/                        # Git configuration & aliases
│   └── fzf/                        # FZF integration
│
├── .claude/                        # Claude Code project config
│   ├── commands/                   # Custom slash commands
│   │   ├── component-new.md        # Create new component
│   │   ├── component-status.md     # Check component health
│   │   ├── sync-check.md           # Check for drift
│   │   └── xdg-audit.md            # Audit XDG compliance
│   │
│   └── agents/                     # Custom agents
│       ├── component-refactor.md   # Refactor code into component
│       ├── xdg-migrate.md          # Migrate tool to XDG compliance
│       └── drift-resolver.md       # Resolve configuration drift
│
└── (legacy files still present - .bash_profile.dir/, .bash_tools/, etc.)
```

### Component Structure

Each component has:
- `component.yaml` - Metadata (name, version, dependencies, XDG dirs)
- `env.sh` - Environment variables (loaded for all shells)
- `aliases.sh` - Shell aliases (interactive only)
- `functions.sh` - Shell functions (interactive only)
- `completions.sh` - Tab completions (interactive only)
- `setup.sh` - Installation script (optional)

### Loading Sequence

```
~/.bashrc / ~/.zshrc
    └─► core/bootstrap.sh
            ├─► core/discovery.sh     # Set CURRENT_SHELL, CURRENT_PLATFORM
            ├─► core/xdg.sh           # Initialize XDG directories
            ├─► core/theme.sh         # Load Catppuccin colors
            ├─► core/loader.sh        # Component discovery & loading
            │       ├─► Scan components/*/component.yaml
            │       ├─► Resolve dependencies (topological sort)
            │       └─► Source env.sh, aliases.sh, functions.sh, completions.sh
            ├─► core/sync.sh          # Drift detection
            └─► Local overrides (~/.config/shell/local.sh)
```

## Key Features Implemented

### 1. Component Loader (core/loader.sh)
- Discovers components by scanning `components/*/component.yaml`
- Parses YAML metadata using `yq` (required dependency)
- Resolves dependencies with topological sort
- Skips components with missing tools (shows warning once, then suppresses)
- Respects `DOTFILES_DISABLE_<COMPONENT>=1` to disable specific components
- Respects `DOTFILES_COMPONENTS="python,go,fzf"` to whitelist components

### 2. XDG Compliance (core/xdg.sh)
- Sets XDG_CONFIG_HOME, XDG_DATA_HOME, XDG_CACHE_HOME, XDG_STATE_HOME
- Provides helper functions:
  - `xdg_config_dir <component>` - Returns ~/.config/<component>
  - `xdg_data_dir <component>` - Returns ~/.local/share/<component>
  - `xdg_cache_dir <component>` - Returns ~/.cache/<component>
  - `xdg_state_dir <component>` - Returns ~/.local/state/<component>
  - `xdg_ensure_dirs <component>` - Creates all XDG dirs for component
- Warning suppression system (shows missing tool warnings once per component)

### 3. Drift Prevention (core/sync.sh)
- `dotfiles_status()` - Show git status of dotfiles repo
- `dotfiles_pull()` - Pull latest from remote
- `dotfiles_push()` - Commit and push changes
- `dotfiles_sync()` - Pull, then push any local changes
- `dotfiles_check_drift()` - Quick drift check (modified/untracked count)
- `dotfiles_audit()` - Full drift report with file details
- `dotfiles_auto_sync()` - Enable/disable auto-sync on startup
- Auto-sync runs on shell startup (~100ms latency acceptable)

### 4. Theme System (core/theme.sh)
- Catppuccin Mocha color palette
- Truecolor detection and fallback to 256-color
- Semantic color aliases (THEME_GIT_CLEAN, THEME_WARNING, THEME_ERROR, etc.)

## Components Created

### shell/
- Core shell aliases (ll, c, cp, df, du, etc.)
- Shell options (histappend, autocd, globstar)
- History configuration (XDG-compliant location)
- Helper functions (mkcd, up, extract, sysinfo)
- Git-aware prompt with Catppuccin colors

### git/
- Git aliases (g, gs, ga, gc, gco, gb, gd, gp, gl, etc.)
- Git functions (gclone, gcob, gpush, gpull, ginfo, gundo, gamend)
- Completion loading

### fzf/
- FZF environment with Catppuccin theme colors
- FZF_DEFAULT_COMMAND using fd
- FZF functions (fzf_files, fzf_cd, fzf_git_branch, fzf_kill, fzf_history)
- K8s and Docker integrations (fzf_k8s_pods, fzf_docker_logs, etc.)

## Remaining Work

### Phase 3-4: Additional Components (Not Yet Migrated)
Components that exist in legacy .bash_tools/ but need migration:
- python (UV, venv management)
- node (NVM, pnpm)
- go
- claude (Claude Code management functions)
- ollama (AI)
- huggingface
- kubernetes
- terraform, aws, azure
- docker
- database (pg, mysql, redis, mongo, sqlite)
- secrets management
- tmux

### Phase 6: Documentation & Polish
- Update main CLAUDE.md with new architecture
- Create component development guide
- Update install.sh for new structure
- Update Makefile with new targets
- Remove legacy files after full migration

## How to Activate

To use the new framework, update ~/.bashrc or ~/.zshrc:

```bash
# Replace existing dotfiles loading with:
source ~/path/to/bash-profile/core/bootstrap.sh
```

## Testing the Framework

```bash
# Test in clean shell
/bin/bash --noprofile --norc -c '
export DOTFILES_ROOT="/Users/mistergrinvalds/Repos/personal/bash-profile"
unset XDG_CONFIG_HOME XDG_DATA_HOME XDG_CACHE_HOME XDG_STATE_HOME
source "$DOTFILES_ROOT/core/discovery.sh"
IS_INTERACTIVE=true
source "$DOTFILES_ROOT/core/xdg.sh"
source "$DOTFILES_ROOT/core/theme.sh"
source "$DOTFILES_ROOT/core/loader.sh"
loader_run
alias  # Should show all loaded aliases
'
```

## Dependencies

- **yq** (v4+) - Required for YAML parsing. Install with: `brew install yq`
- **fd** - Optional, used by FZF for file finding
- **fzf** - Optional, for fuzzy finding features

## Known Issues

- User's system has `XDG_CONFIG_HOME` pre-set to just `$HOME` instead of `~/.config`. The framework respects this per XDG spec, but it may cause issues.
- `/opt/homebrew/etc/bash_completion.d/pgcli` has a syntax error (unrelated to this framework)

## Git Status at End of Session

Branch: `fix/automation-framework-compatibility`

New untracked directories:
- `core/` - New framework core
- `components/` - New component structure
- `.claude/commands/` - Claude Code commands
- `.claude/agents/` - Claude Code agents

Modified files:
- `.claude/settings.local.json` - Updated with heightened permissions

## User Preferences (from session)

- **Metadata format**: YAML (component.yaml) with yq dependency
- **Missing tool handling**: Skip with warning (show on first startup, then suppress)
- **Auto-sync behavior**: Always run on startup (~100ms latency acceptable)
- **Migration priority**: shell → git → fzf first (completed)
