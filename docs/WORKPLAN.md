# Dotfiles Restructuring Work Plan

> XDG-compliant, unified bash/zsh configuration with source chain injection

## Overview

| Aspect | Decision |
|--------|----------|
| **RC Strategy** | XDG source chain (`~/.bashrc` → `~/.config/shell/init.sh` → repo) |
| **Shell Config** | Unified (one `init.sh` with shell detection) |
| **Repo Location** | `$XDG_CONFIG_HOME/dotfiles` (`~/.config/dotfiles`) |
| **Secrets** | `$XDG_CONFIG_HOME/secrets/` with encrypted `.env` files |

## Dependency Graph

See: [shell-injection-dag.excalidraw](./shell-injection-dag.excalidraw)

## Target Directory Structure

```
~/.config/dotfiles/                    # REPO ($DOTFILES_ROOT)
├── shell/                             # Shell configuration
│   ├── init.sh                        # Main entry point
│   ├── discovery.sh                   # Shell/platform detection (FIRST)
│   ├── xdg.sh                         # XDG directory setup
│   ├── environment.sh                 # Core environment (PATH, EDITOR)
│   ├── environment.darwin.sh          # macOS-specific
│   ├── environment.linux.sh           # Linux-specific
│   ├── secrets.sh                     # Secrets loading (silent)
│   ├── options.sh                     # Shell options (shopt/setopt)
│   ├── aliases.sh                     # Shell aliases
│   ├── completions.sh                 # Completion setup
│   └── prompt.sh                      # Prompt configuration
│
├── functions/                         # Shell functions
│   ├── core/                          # Core utilities
│   │   ├── cd.sh                      # Enhanced cd
│   │   ├── history.sh                 # History search (h)
│   │   ├── archive.sh                 # mktar, mkzip
│   │   ├── utils.sh                   # bash_as, etc.
│   │   └── inject.sh                  # dotfiles_inject/eject/reload
│   ├── dev/                           # Development tools
│   │   ├── python.sh                  # mkvenv, venv, rmenv
│   │   ├── golang.sh                  # Go helpers
│   │   └── github.sh                  # GitHub CLI helpers
│   ├── cloud/                         # Cloud/infrastructure
│   │   ├── kubernetes.sh              # K8s helpers
│   │   ├── providers.sh               # AWS, Azure, DO
│   │   └── automation.sh              # Automation framework
│   └── ai/                            # AI/ML tools
│       ├── ollama.sh                  # Ollama integration
│       ├── huggingface.sh             # HuggingFace integration
│       └── aliases.sh                 # AI convenience aliases
│
├── config/                            # Application configs
│   ├── git/
│   │   ├── config                     # Main gitconfig
│   │   ├── config.alias               # Git aliases
│   │   ├── config.color               # Git colors
│   │   └── config.work                # Work-specific (conditional)
│   ├── ssh/
│   │   └── config                     # SSH with identity masking
│   ├── python/
│   │   ├── startup.py                 # PYTHONSTARTUP
│   │   └── ipython/                   # IPython config
│   ├── r/
│   │   └── Rprofile                   # R startup
│   ├── conda/
│   │   └── condarc                    # Conda config
│   └── karabiner/                     # Keyboard remapping
│
├── secrets/                           # Secrets (git-ignored)
│   ├── .env                           # Main secrets file
│   ├── .env.template                  # Template (tracked)
│   └── README.md                      # Setup instructions
│
├── automation/                        # Automation framework
│   ├── auto                           # CLI entry point
│   ├── framework/                     # Core framework
│   └── modules/                       # Feature modules
│
├── docs/                              # Documentation
│   ├── WORKPLAN.md                    # This file
│   └── shell-injection-dag.excalidraw
│
├── install.sh                         # Installer
├── Makefile                           # Testing
├── CLAUDE.md                          # AI assistant guide
└── README.md                          # Main documentation
```

---

## Phase 1: Foundation (Prep Work)

### 1.1 Create New Directory Structure
- [ ] Create `shell/` directory for shell configuration
- [ ] Create `functions/` directory with subdirectories (`core/`, `dev/`, `cloud/`, `ai/`)
- [ ] Create `config/` directory for application configs
- [ ] Create `docs/` directory (done - this file is here)

### 1.2 Shell Discovery Module
- [ ] Create `shell/discovery.sh` with:
  - `$CURRENT_SHELL` detection (bash/zsh/unknown)
  - `$CURRENT_PLATFORM` detection (darwin/linux/unknown)
  - `$IS_INTERACTIVE` detection
  - `$IS_LOGIN_SHELL` detection
- [ ] Add early-exit for non-interactive shells (speed optimization)

### 1.3 XDG Module
- [ ] Create `shell/xdg.sh` with:
  - `$DOTFILES_ROOT` (repo location)
  - `$XDG_CONFIG_HOME` (default: `~/.config`)
  - `$XDG_DATA_HOME` (default: `~/.local/share`)
  - `$XDG_CACHE_HOME` (default: `~/.cache`)
  - `$XDG_STATE_HOME` (default: `~/.local/state`)
  - `$XDG_RUNTIME_DIR` (platform-specific)
- [ ] Create XDG directories if they don't exist

**Deliverable:** `make test-discovery` passes

---

## Phase 2: Core Shell Configuration

### 2.1 Environment Setup
- [ ] Migrate `.bash_profile.dir/environment.sh` → `shell/environment.sh`
- [ ] Migrate `.bash_profile.dir/environment_darwin.sh` → `shell/environment.darwin.sh`
- [ ] Migrate `.bash_profile.dir/environment_linux-gnu.sh` → `shell/environment.linux.sh`
- [ ] Update paths to use XDG variables
- [ ] Remove hardcoded `$HOME` references

### 2.2 Shell Options
- [ ] Create `shell/options.sh` with:
  - Bash-specific: `shopt` commands, `HISTCONTROL`, `HISTSIZE`
  - Zsh-specific: `setopt` commands, `EXTENDED_GLOB`, `AUTO_CD`
  - Shared: `HISTFILE` location (XDG-compliant)

### 2.3 Aliases
- [ ] Migrate `.bash_profile.dir/aliases.sh` → `shell/aliases.sh`
- [ ] Import git aliases from `tmp-dotfiles/.gitconfig.alias` (convert to shell aliases or keep in gitconfig)
- [ ] Add platform-aware aliases (pbcopy/xclip, open/xdg-open)

**Deliverable:** `make test-environment` passes

---

## Phase 3: Functions Migration

### 3.1 Core Functions
- [ ] Migrate `cd()` → `functions/core/cd.sh`
- [ ] Migrate `h()` (history) → `functions/core/history.sh`
- [ ] Migrate `mktar()`, `mkzip()` → `functions/core/archive.sh`
- [ ] Migrate `bash_as()` → `functions/core/utils.sh`

### 3.2 Development Functions
- [ ] Migrate `mkvenv()`, `venv()`, `rmenv()` → `functions/dev/python.sh`
- [ ] Migrate golang functions → `functions/dev/golang.sh`
- [ ] Migrate github functions → `functions/dev/github.sh`

### 3.3 Cloud/Infrastructure Functions
- [ ] Migrate kubernetes functions → `functions/cloud/kubernetes.sh`
- [ ] Migrate cloud provider functions → `functions/cloud/providers.sh`
- [ ] Keep automation framework integration → `functions/cloud/automation.sh`

### 3.4 AI Functions
- [ ] Migrate `ollama.sh` → `functions/ai/ollama.sh`
- [ ] Migrate `huggingface.sh` → `functions/ai/huggingface.sh`
- [ ] Migrate `ai_aliases.sh` → `functions/ai/aliases.sh`

### 3.5 Injection Functions (NEW)
- [ ] Create `functions/core/inject.sh` with:
  - `dotfiles_inject()` - Install bootstrap files
  - `dotfiles_eject()` - Remove all injected config
  - `dotfiles_update()` - Pull and reload
  - `dotfiles_reload()` - Reload without restart
  - `dotfiles_status()` - Show current state
  - `dotfiles_link_configs()` - Symlink app configs

**Deliverable:** `make test-functions` passes

---

## Phase 4: Application Configs (from tmp-dotfiles)

### 4.1 Git Configuration
- [ ] Move `tmp-dotfiles/.gitconfig` → `config/git/config`
- [ ] Move `tmp-dotfiles/.gitconfig.alias` → `config/git/config.alias`
- [ ] Move `tmp-dotfiles/.gitconfig.color` → `config/git/config.color`
- [ ] Move `tmp-dotfiles/.gitconfig.work` → `config/git/config.work`
- [ ] Update include paths to use `$XDG_CONFIG_HOME/dotfiles/config/git/`
- [ ] Add `dotfiles_link_configs()` logic for git

### 4.2 SSH Configuration
- [ ] Move `tmp-dotfiles/.ssh/config` → `config/ssh/config`
- [ ] Document multi-identity pattern
- [ ] Add symlink logic (SSH requires `~/.ssh/config`)

### 4.3 Python Configuration
- [ ] Move `tmp-dotfiles/.python/` → `config/python/`
- [ ] Move `tmp-dotfiles/.ipython/` → `config/ipython/`
- [ ] Update `PYTHONSTARTUP` to point to repo location
- [ ] Update `IPYTHONDIR` to XDG location

### 4.4 Other Configs
- [ ] Move `tmp-dotfiles/.R/` → `config/r/`
- [ ] Move `tmp-dotfiles/.condarc` → `config/conda/condarc`
- [ ] Move `tmp-dotfiles/.karabiner/` → `config/karabiner/`
- [ ] Move `tmp-dotfiles/.wgetrc` → `config/wget/wgetrc`

### 4.5 Cleanup
- [ ] Remove `tmp-dotfiles/` directory after migration complete
- [ ] Update `.gitignore` for new structure

**Deliverable:** `make test-configs` passes

---

## Phase 5: Secrets Management

### 5.1 Secrets Structure
- [ ] Create `secrets/` directory (git-ignored)
- [ ] Create `secrets/.env.template` (tracked, no values)
- [ ] Create `secrets/README.md` with setup instructions

### 5.2 Secrets Loading
- [ ] Create `shell/secrets.sh` with:
  - Silent loading (no echo)
  - Validation of required variables
  - Support for encrypted files (optional gpg integration)
  - Per-environment secrets (`.env.work`, `.env.personal`)

### 5.3 Secrets Functions
- [ ] Migrate existing secrets functions from `.bash_tools/secrets.sh`
- [ ] Add `secrets_init()` - Create secrets directory structure
- [ ] Add `secrets_edit()` - Securely edit secrets file
- [ ] Add `secrets_validate()` - Check required keys exist

**Deliverable:** `make test-secrets` passes

---

## Phase 6: Prompt Configuration

### 6.1 Fix Zsh Compatibility
- [ ] Rewrite `terminal.sh` → `shell/prompt.sh`
- [ ] Replace `BASH_REMATCH` with portable pattern matching
- [ ] Use `case` statements instead of bash regex
- [ ] Add zsh `PROMPT` variable alongside bash `PS1`
- [ ] Use proper color escaping:
  - Bash: `\[\e[32m\]`
  - Zsh: `%F{green}` or `%{$fg[green]%}`

### 6.2 Git Prompt Functions
- [ ] Make `git_branch()` shell-portable
- [ ] Make `git_color()` shell-portable
- [ ] Test in both bash and zsh

**Deliverable:** `make test-prompt` passes in both shells

---

## Phase 7: Completions

### 7.1 Completion Setup
- [ ] Create `shell/completions.sh` with:
  - Bash: `bash-completion` loading
  - Zsh: `compinit` setup
  - FZF completions (if available)
  - Custom completions for dotfiles functions

### 7.2 FZF Integration
- [ ] Migrate FZF config to XDG-compliant location
- [ ] Load shell-specific FZF bindings

**Deliverable:** Tab completion works in both shells

---

## Phase 8: Main Init Script

### 8.1 Create Unified Entry Point
- [ ] Create `shell/init.sh` that sources in order:
  1. `discovery.sh`
  2. `xdg.sh`
  3. `environment.sh` + platform-specific
  4. `secrets.sh`
  5. `options.sh`
  6. `aliases.sh`
  7. All `functions/**/*.sh`
  8. `completions.sh`
  9. `prompt.sh`
  10. Local overrides (`~/.config/shell/local.sh`)

### 8.2 Bootstrap Files
- [ ] Create template for `~/.bashrc` (3 lines)
- [ ] Create template for `~/.zshrc` (3 lines)
- [ ] Document the source chain

**Deliverable:** Full shell loads without errors in both bash and zsh

---

## Phase 9: Installer Update

### 9.1 Update install.sh
- [ ] Rename `initialize.sh` → `install.sh`
- [ ] Add XDG directory creation
- [ ] Add bootstrap file creation (`~/.bashrc`, `~/.zshrc`)
- [ ] Add symlink creation for `~/.config/shell/init.sh`
- [ ] Add app config symlink logic
- [ ] Add `--uninstall` / `--eject` option
- [ ] Add `--dry-run` option

### 9.2 Installation Modes
- [ ] Interactive mode (default)
- [ ] Automatic mode (`--auto` or `-y`)
- [ ] Minimal mode (`--minimal` - just shell, no apps)
- [ ] Full mode (`--full` - everything including automation)

**Deliverable:** `./install.sh --dry-run` shows correct actions

---

## Phase 10: Testing

### 10.1 Update Makefile
- [ ] Add `test-discovery` - Shell/platform detection
- [ ] Add `test-xdg` - XDG variable setup
- [ ] Add `test-environment` - Environment loading
- [ ] Add `test-functions` - Function definitions
- [ ] Add `test-configs` - App config symlinks
- [ ] Add `test-secrets` - Secrets loading
- [ ] Add `test-prompt` - Prompt rendering
- [ ] Add `test-zsh` - Full zsh compatibility
- [ ] Add `test-bash` - Full bash compatibility
- [ ] Add `test-shell-compat` - Test both shells
- [ ] Add `verify` - Quick smoke test

### 10.2 CI Testing
- [ ] Add GitHub Actions workflow
- [ ] Test on Ubuntu (bash + zsh)
- [ ] Test on macOS (bash + zsh)

**Deliverable:** `make test-all` passes

---

## Phase 11: Documentation

### 11.1 Update Docs
- [ ] Update `README.md` with new structure
- [ ] Update `CLAUDE.md` with new architecture
- [ ] Create `docs/MIGRATION.md` for existing users
- [ ] Update `INSTALL.md`

### 11.2 Inline Documentation
- [ ] Add header comments to all shell files
- [ ] Document all exported variables
- [ ] Document all public functions

**Deliverable:** Documentation matches implementation

---

## Phase 12: Cleanup

### 12.1 Remove Old Structure
- [ ] Remove `.bash_profile` (replaced by `shell/init.sh`)
- [ ] Remove `.bash_profile.dir/` (migrated to `shell/`)
- [ ] Remove `.bash_tools/` (migrated to `functions/`)
- [ ] Remove `tmp-dotfiles/` (migrated to `config/`)

### 12.2 Final Verification
- [ ] Fresh install on clean system
- [ ] Test in Docker container (bash)
- [ ] Test in Docker container (zsh)

**Deliverable:** Clean, working system

---

## File Mapping Reference

| Old Location | New Location |
|--------------|--------------|
| `.bash_profile` | `shell/init.sh` |
| `.bash_profile.dir/environment.sh` | `shell/environment.sh` |
| `.bash_profile.dir/environment_darwin.sh` | `shell/environment.darwin.sh` |
| `.bash_profile.dir/environment_linux-gnu.sh` | `shell/environment.linux.sh` |
| `.bash_profile.dir/aliases.sh` | `shell/aliases.sh` |
| `.bash_profile.dir/terminal.sh` | `shell/prompt.sh` |
| `.bash_tools/*.sh` | `functions/{core,dev,cloud,ai}/*.sh` |
| `tmp-dotfiles/.gitconfig*` | `config/git/config*` |
| `tmp-dotfiles/.ssh/config` | `config/ssh/config` |
| `tmp-dotfiles/.python/` | `config/python/` |
| `tmp-dotfiles/.ipython/` | `config/ipython/` |
| `tmp-dotfiles/.R/` | `config/r/` |
| `tmp-dotfiles/.condarc` | `config/conda/condarc` |
| `tmp-dotfiles/.karabiner/` | `config/karabiner/` |
| `initialize.sh` | `install.sh` |

---

## Quick Start Commands

```bash
# Start work on Phase 1
make setup-phase1

# Test current phase
make test-quick

# Test everything
make test-all

# Verify in both shells
make test-shell-compat
```

---

## Timeline Estimate

This is a significant restructuring. Work through phases in order - each phase builds on the previous. Mark items complete as you go.
