# Reorganization Progress Tracker

> Last updated: 2025-12-12
> Status: **Phase 9 - Rename & Cleanup (in progress)**

---

## Quick Reference

```
Current branch: fix/automation-framework-compatibility
Target structure: ~/.config/dotfiles/ (XDG-compliant)
Shell strategy: Unified init.sh with discovery
```

---

## Phase 0: Planning ✅

- [x] Investigate repository structure
- [x] Create dependency DAG (shell-injection-dag.excalidraw)
- [x] Define XDG source chain strategy
- [x] Create work plan (WORKPLAN.md)
- [x] Create reorganization mapping (REORGANIZATION.md)
- [x] Commit existing AI integration work
- [x] Create this progress tracker

---

## Phase 1: Create Directory Structure ✅

**Status:** Complete (commit 6550a1f)

- [x] Create `shell/` directory
- [x] Create `functions/` with subdirectories
- [x] Create `config/` with subdirectories
- [x] Create `secrets/` directory with .gitignore
- [x] Create `docs/automation/` directory
- [x] Commit: `create new directory structure`

---

## Phase 2: Shell Discovery & XDG (Critical Path) ✅

**Status:** Complete (commit 380a02f)

### 2.1 Create shell/discovery.sh
- [x] Detect `$CURRENT_SHELL` (bash/zsh/unknown)
- [x] Detect `$CURRENT_PLATFORM` (darwin/linux/unknown)
- [x] Detect `$IS_INTERACTIVE`
- [x] Detect `$IS_LOGIN_SHELL`
- [x] Early exit for non-interactive shells

### 2.2 Create shell/xdg.sh
- [x] Set `$DOTFILES_ROOT` (repo location)
- [x] Set `$XDG_CONFIG_HOME` (default: ~/.config)
- [x] Set `$XDG_DATA_HOME` (default: ~/.local/share)
- [x] Set `$XDG_CACHE_HOME` (default: ~/.cache)
- [x] Set `$XDG_STATE_HOME` (default: ~/.local/state)
- [x] Platform-specific `$XDG_RUNTIME_DIR`
- [x] Create directories if missing
- [x] Set `$HISTFILE` to XDG-compliant location

### 2.3 Commit
- [x] Commit: `add shell discovery and XDG setup`

---

## Phase 3: Move Shell Configuration ✅

**Status:** Complete

### 3.1 Move existing files
- [x] `environment.sh` → `shell/environment.sh` (updated for XDG)
- [x] `environment_darwin.sh` → `shell/environment.darwin.sh`
- [x] `environment_linux-gnu.sh` → `shell/environment.linux.sh`
- [x] `aliases.sh` → `shell/aliases.sh` (made portable)
- [x] `terminal.sh` → `shell/prompt.sh`

### 3.2 Create new files
- [x] Create `shell/options.sh` (shopt/setopt)
- [x] Create `shell/secrets.sh` (silent loading)
- [x] Create `shell/completions.sh` (bash-completion/compinit)

### 3.3 Fix zsh compatibility in prompt.sh
- [x] Replace `BASH_REMATCH` with git symbolic-ref (portable)
- [x] Add zsh `PROMPT` alongside bash `PS1`
- [x] Fix color escape sequences for zsh

### 3.4 Commit
- [ ] Commit: `move shell configuration to shell/`

---

## Phase 4: Move Functions ✅

**Status:** Complete

### 4.1 Core functions
- [x] `cd.sh` → `functions/core/cd.sh`
- [x] `history.sh` → `functions/core/history.sh`
- [x] Merge `mktar.sh` + `mkzip.sh` → `functions/core/archive.sh`
- [x] `bash_as.sh` → `functions/core/utils.sh`
- [x] `tools.sh` → `functions/core/tools.sh`
- [x] `tmux_helpers.sh` → `functions/core/tmux.sh`

### 4.2 Dev functions
- [x] Merge `mkvenv.sh` + `venv.sh` + `rmenv.sh` → `functions/dev/python.sh`
- [x] `golang.sh` → `functions/dev/golang.sh`
- [x] `github.sh` → `functions/dev/github.sh`
- [x] `vscode.sh` → `functions/dev/vscode.sh`

### 4.3 Cloud functions
- [x] `kubernetes.sh` → `functions/cloud/kubernetes.sh`
- [x] `automation.sh` → `functions/cloud/automation.sh`
- [x] `secrets.sh` → `functions/cloud/secrets.sh`

### 4.4 AI functions
- [x] `ollama.sh` → `functions/ai/ollama.sh`
- [x] `huggingface.sh` → `functions/ai/huggingface.sh`
- [x] `ai_aliases.sh` → `functions/ai/aliases.sh`

### 4.5 Completion-related (move to shell/)
- [x] `bash_completion.sh` → merge into `shell/completions.sh`
- [x] `fzf.sh` → merge into `shell/completions.sh`

### 4.6 Commit
- [ ] Commit: `move functions to functions/`

---

## Phase 5: Import App Configs from tmp-dotfiles ✅

**Status:** Complete

### 5.1 Git config
- [x] `.gitconfig` → `config/git/config`
- [x] `.gitconfig.alias` → `config/git/config.alias`
- [x] `.gitconfig.color` → `config/git/config.color`
- [x] `.gitconfig.work` → `config/git/config.work`
- [x] Update include paths in config

### 5.2 SSH config
- [x] `.ssh/config` → `config/ssh/config`

### 5.3 Python config
- [x] `.python/*` → `config/python/`
- [ ] `.ipython/*` → `config/python/ipython/` (skipped - empty)

### 5.4 Other configs
- [x] `.R/*` → `config/r/`
- [x] `.condarc` → `config/conda/condarc`
- [x] `.karabiner/*` → `config/karabiner/`
- [x] `.wgetrc` → `config/wget/wgetrc`

### 5.5 Commit
- [ ] Commit: `import app configs from tmp-dotfiles`

---

## Phase 6: Create Injection Functions ✅

**Status:** Complete

### 6.1 Create functions/core/inject.sh
- [x] `dotfiles_inject()` - install bootstrap files
- [x] `dotfiles_eject()` - remove all injected config
- [x] `dotfiles_update()` - git pull + reload
- [x] `dotfiles_reload()` - reload without restart
- [x] `dotfiles_status()` - show current state
- [x] `dotfiles_link_configs()` - symlink app configs

### 6.2 Commit
- [ ] Commit: `add dotfiles injection functions`

---

## Phase 7: Create Main Entry Point ✅

**Status:** Complete

### 7.1 Create shell/init.sh
- [x] Source in correct order:
  1. discovery.sh
  2. xdg.sh
  3. environment.sh + platform-specific
  4. secrets.sh
  5. options.sh
  6. aliases.sh
  7. All functions/**/*.sh
  8. completions.sh
  9. prompt.sh
  10. Local overrides (~/.config/shell/local.sh)

### 7.2 Create bootstrap templates
- [x] Template for ~/.bashrc (in inject.sh)
- [x] Template for ~/.zshrc (in inject.sh)

### 7.3 Commit
- [ ] Commit: `create unified shell/init.sh entry point`

---

## Phase 8: Move Documentation ✅

**Status:** Complete

- [x] `INSTALL.md` → `docs/INSTALL.md`
- [x] `.automation/README.md` → `docs/automation/README.md`
- [x] `.automation/AI.md` → `docs/automation/AI.md`
- [x] `.automation/SECRETS.md` → `docs/automation/SECRETS.md`
- [x] `.automation/TOOLS.md` → `docs/automation/TOOLS.md`
- [ ] Commit: `move documentation to docs/`

---

## Phase 9: Rename & Cleanup ✅

**Status:** Complete

### 9.1 Rename
- [x] `initialize.sh` → `install.sh`

### 9.2 Delete temp files
- [x] Remove `create_branch_and_commits.sh`
- [x] Remove `execute_git.sh`
- [x] Remove `git_operations.py`
- [x] Remove `run_git_commands.py`

### 9.3 Delete node_modules (if not needed)
- [x] Check if package.json is needed (empty, not needed)
- [x] Remove `node_modules/`, `package.json`, `package-lock.json`

### 9.4 Delete legacy directories
- [ ] Remove `.bash_profile.dir/` (after testing)
- [ ] Remove `.bash_tools/` (after testing)
- [ ] Remove `tmp-dotfiles/` (after testing)

### 9.5 Commit
- [ ] Commit: `rename install.sh and remove legacy files`

---

## Phase 10: Update References ✅

**Status:** Complete

- [x] Update `Makefile` paths for new structure
- [x] Update `CLAUDE.md` with new architecture
- [x] Update `README.md` with new structure
- [x] Update `install.sh` for new paths
- [ ] Commit: `update all references for new structure`

---

## Phase 11: Testing

**Status:** Not started

### 11.1 Add new test targets
- [ ] `test-discovery` - shell/platform detection
- [ ] `test-xdg` - XDG variable setup
- [ ] `test-zsh` - full zsh compatibility
- [ ] `test-shell-compat` - both shells
- [ ] `verify` - quick smoke test

### 11.2 Run tests
- [ ] `make test-syntax` passes
- [ ] `make test-quick` passes
- [ ] Shell loads in bash without errors
- [ ] Shell loads in zsh without errors
- [ ] All functions are defined
- [ ] Git prompt works in both shells
- [ ] Completions load in both shells

### 11.3 Commit
- [ ] Commit: `add shell compatibility tests`

---

## Phase 12: Final Documentation

**Status:** Not started

- [ ] Update `README.md` with final structure
- [ ] Update `CLAUDE.md` for AI assistant
- [ ] Create `docs/MIGRATION.md` for existing users
- [ ] Remove or archive planning docs (WORKPLAN.md, REORGANIZATION.md)
- [ ] Commit: `finalize documentation`

---

## Blocking Issues

*Track any issues that block progress here*

| Issue | Status | Resolution |
|-------|--------|------------|
| *(none yet)* | | |

---

## Notes

*Running notes and decisions made during implementation*

- **2025-12-12**: Created planning documents, committed AI integration work
- XDG source chain chosen: `~/.bashrc` → `~/.config/shell/init.sh` → repo
- Repo target location: `~/.config/dotfiles`

---

## Commands Cheatsheet

```bash
# Test current syntax
make test-syntax

# Test in bash
bash --noprofile --norc -c 'source shell/init.sh'

# Test in zsh
zsh -f -c 'source shell/init.sh'

# Check what's changed
git status

# Commit current phase
git add -A && git commit -m "phase X: description"
```
