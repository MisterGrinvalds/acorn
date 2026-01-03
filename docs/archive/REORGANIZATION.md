# Code Reorganization Plan

> Mapping current structure to XDG-compliant, unified shell configuration

## Summary

This plan reorganizes the codebase to:
1. Follow XDG Base Directory specification
2. Unify bash/zsh configuration with shell detection
3. Use clear, consistent naming conventions
4. Separate concerns (shell config, functions, app configs, automation)

---

## File Mapping: Current → Target

### Shell Configuration

| Current | Target | Action |
|---------|--------|--------|
| `.bash_profile` | `shell/init.sh` | **Rewrite** - becomes unified entry point |
| `.bash_profile.dir/environment.sh` | `shell/environment.sh` | Move + update XDG vars |
| `.bash_profile.dir/environment_darwin.sh` | `shell/environment.darwin.sh` | Rename (remove underscore) |
| `.bash_profile.dir/environment_linux-gnu.sh` | `shell/environment.linux.sh` | Rename (simplify) |
| `.bash_profile.dir/aliases.sh` | `shell/aliases.sh` | Move |
| `.bash_profile.dir/terminal.sh` | `shell/prompt.sh` | Rename + fix zsh compat |
| *(new)* | `shell/discovery.sh` | **Create** - shell/platform detection |
| *(new)* | `shell/xdg.sh` | **Create** - XDG directory setup |
| *(new)* | `shell/options.sh` | **Create** - shell options (shopt/setopt) |
| *(new)* | `shell/completions.sh` | **Create** - completion setup |
| *(new)* | `shell/secrets.sh` | **Create** - secrets loading |

### Functions (`.bash_tools/` → `functions/`)

| Current | Target | Action |
|---------|--------|--------|
| `.bash_tools/cd.sh` | `functions/core/cd.sh` | Move |
| `.bash_tools/history.sh` | `functions/core/history.sh` | Move |
| `.bash_tools/mktar.sh` | `functions/core/archive.sh` | Merge with mkzip |
| `.bash_tools/mkzip.sh` | `functions/core/archive.sh` | Merge with mktar |
| `.bash_tools/bash_as.sh` | `functions/core/utils.sh` | Move + rename |
| `.bash_tools/bash_completion.sh` | `shell/completions.sh` | Move to shell config |
| `.bash_tools/fzf.sh` | `shell/completions.sh` | Merge into completions |
| `.bash_tools/mkvenv.sh` | `functions/dev/python.sh` | Merge python functions |
| `.bash_tools/venv.sh` | `functions/dev/python.sh` | Merge python functions |
| `.bash_tools/rmenv.sh` | `functions/dev/python.sh` | Merge python functions |
| `.bash_tools/golang.sh` | `functions/dev/golang.sh` | Move |
| `.bash_tools/github.sh` | `functions/dev/github.sh` | Move |
| `.bash_tools/vscode.sh` | `functions/dev/vscode.sh` | Move |
| `.bash_tools/kubernetes.sh` | `functions/cloud/kubernetes.sh` | Move |
| `.bash_tools/automation.sh` | `functions/cloud/automation.sh` | Move |
| `.bash_tools/secrets.sh` | `functions/cloud/secrets.sh` | Move |
| `.bash_tools/tools.sh` | `functions/core/tools.sh` | Move |
| `.bash_tools/tmux_helpers.sh` | `functions/core/tmux.sh` | Rename |
| `.bash_tools/ollama.sh` | `functions/ai/ollama.sh` | Move |
| `.bash_tools/huggingface.sh` | `functions/ai/huggingface.sh` | Move |
| `.bash_tools/ai_aliases.sh` | `functions/ai/aliases.sh` | Move |
| *(new)* | `functions/core/inject.sh` | **Create** - dotfiles management |

### Application Configs (from `tmp-dotfiles/`)

| Current | Target | Action |
|---------|--------|--------|
| `tmp-dotfiles/.gitconfig` | `config/git/config` | Move |
| `tmp-dotfiles/.gitconfig.alias` | `config/git/config.alias` | Move |
| `tmp-dotfiles/.gitconfig.color` | `config/git/config.color` | Move |
| `tmp-dotfiles/.gitconfig.work` | `config/git/config.work` | Move |
| `tmp-dotfiles/.ssh/config` | `config/ssh/config` | Move |
| `tmp-dotfiles/.python/*` | `config/python/` | Move directory |
| `tmp-dotfiles/.ipython/*` | `config/python/ipython/` | Move + nest |
| `tmp-dotfiles/.R/*` | `config/r/` | Move directory |
| `tmp-dotfiles/.condarc` | `config/conda/condarc` | Move |
| `tmp-dotfiles/.karabiner/*` | `config/karabiner/` | Move directory |
| `tmp-dotfiles/.wgetrc` | `config/wget/wgetrc` | Move |

### Automation Framework

| Current | Target | Action |
|---------|--------|--------|
| `.automation/` | `automation/` | Keep (just verify naming) |
| `.automation/auto` | `automation/auto` | Keep |
| `.automation/framework/` | `automation/framework/` | Keep |
| `.automation/modules/` | `automation/modules/` | Keep |
| `.automation/*.md` | `docs/automation/` | Move docs to docs/ |

### Documentation

| Current | Target | Action |
|---------|--------|--------|
| `README.md` | `README.md` | Update for new structure |
| `CLAUDE.md` | `CLAUDE.md` | Update for new structure |
| `INSTALL.md` | `docs/INSTALL.md` | Move to docs/ |
| `.automation/README.md` | `docs/automation/README.md` | Move |
| `.automation/AI.md` | `docs/automation/AI.md` | Move |
| `.automation/SECRETS.md` | `docs/automation/SECRETS.md` | Move |
| `.automation/TOOLS.md` | `docs/automation/TOOLS.md` | Move |

### Root Level

| Current | Target | Action |
|---------|--------|--------|
| `initialize.sh` | `install.sh` | Rename |
| `Makefile` | `Makefile` | Update paths |
| `create_branch_and_commits.sh` | *(delete)* | Remove temp script |
| `execute_git.sh` | *(delete)* | Remove temp script |
| `git_operations.py` | *(delete)* | Remove temp script |
| `run_git_commands.py` | *(delete)* | Remove temp script |
| `package.json` | *(delete)* | Remove if not needed |
| `package-lock.json` | *(delete)* | Remove if not needed |
| `node_modules/` | *(delete)* | Remove if not needed |

### Files to Create

| File | Purpose |
|------|---------|
| `shell/init.sh` | Main entry point - sources all config in order |
| `shell/discovery.sh` | Detect $CURRENT_SHELL, $CURRENT_PLATFORM |
| `shell/xdg.sh` | Set XDG variables, create directories |
| `shell/options.sh` | Shell options (bash shopt / zsh setopt) |
| `shell/completions.sh` | Completion system setup |
| `shell/secrets.sh` | Silent secrets loading |
| `functions/core/inject.sh` | dotfiles_inject/eject/reload/status |
| `secrets/.env.template` | Template for secrets (tracked) |
| `secrets/README.md` | Secrets setup instructions |

---

## Naming Conventions

### Directories
- **Lowercase** with hyphens for multi-word: `shell/`, `functions/`, `config/`
- **No dots** at start (except hidden files that must be hidden)

### Shell Files
- **Lowercase** with dots for variants: `environment.darwin.sh`
- **No underscores** in filenames: ~~`environment_darwin.sh`~~ → `environment.darwin.sh`
- **Descriptive names**: `prompt.sh` not `terminal.sh`

### Functions
- **Grouped by domain**: `core/`, `dev/`, `cloud/`, `ai/`
- **Merged related functions**: `python.sh` contains mkvenv, venv, rmenv

---

## Directory Deletions

After reorganization, these directories should be removed:

```
.bash_profile.dir/     → merged into shell/
.bash_tools/           → merged into functions/
tmp-dotfiles/          → merged into config/
node_modules/          → not needed
```

---

## Execution Order

### Phase 1: Create New Structure
```bash
mkdir -p shell functions/{core,dev,cloud,ai} config/{git,ssh,python,r,conda,karabiner,wget} secrets docs/automation
```

### Phase 2: Move Shell Config
```bash
# Create new files first, then move content
touch shell/{init,discovery,xdg,options,completions,secrets}.sh
mv .bash_profile.dir/environment.sh shell/
mv .bash_profile.dir/environment_darwin.sh shell/environment.darwin.sh
mv .bash_profile.dir/environment_linux-gnu.sh shell/environment.linux.sh
mv .bash_profile.dir/aliases.sh shell/
mv .bash_profile.dir/terminal.sh shell/prompt.sh
```

### Phase 3: Move Functions
```bash
# Core
mv .bash_tools/cd.sh functions/core/
mv .bash_tools/history.sh functions/core/
cat .bash_tools/mktar.sh .bash_tools/mkzip.sh > functions/core/archive.sh
mv .bash_tools/bash_as.sh functions/core/utils.sh
mv .bash_tools/tools.sh functions/core/
mv .bash_tools/tmux_helpers.sh functions/core/tmux.sh

# Dev
cat .bash_tools/{mkvenv,venv,rmenv}.sh > functions/dev/python.sh
mv .bash_tools/golang.sh functions/dev/
mv .bash_tools/github.sh functions/dev/
mv .bash_tools/vscode.sh functions/dev/

# Cloud
mv .bash_tools/kubernetes.sh functions/cloud/
mv .bash_tools/automation.sh functions/cloud/
mv .bash_tools/secrets.sh functions/cloud/

# AI
mv .bash_tools/ollama.sh functions/ai/
mv .bash_tools/huggingface.sh functions/ai/
mv .bash_tools/ai_aliases.sh functions/ai/aliases.sh
```

### Phase 4: Move App Configs (from tmp-dotfiles)
```bash
mv tmp-dotfiles/.gitconfig* config/git/
mv tmp-dotfiles/.ssh/config config/ssh/
mv tmp-dotfiles/.python/* config/python/
mv tmp-dotfiles/.ipython config/python/ipython
mv tmp-dotfiles/.R/* config/r/
mv tmp-dotfiles/.condarc config/conda/condarc
mv tmp-dotfiles/.karabiner/* config/karabiner/
mv tmp-dotfiles/.wgetrc config/wget/wgetrc
```

### Phase 5: Move Documentation
```bash
mv INSTALL.md docs/
mv .automation/*.md docs/automation/
```

### Phase 6: Rename & Cleanup
```bash
mv initialize.sh install.sh
rm -rf .bash_profile.dir .bash_tools tmp-dotfiles
rm -f create_branch_and_commits.sh execute_git.sh git_operations.py run_git_commands.py
rm -rf node_modules package.json package-lock.json
```

### Phase 7: Update References
- Update `Makefile` paths
- Update `CLAUDE.md` documentation
- Update `README.md` structure
- Create `shell/init.sh` as new entry point

---

## Verification Checklist

After reorganization:

- [ ] `make test-syntax` passes
- [ ] `make test-quick` passes
- [ ] Shell loads in bash without errors
- [ ] Shell loads in zsh without errors
- [ ] All functions are defined
- [ ] Git prompt works
- [ ] Completions load
- [ ] `dotfiles_status` command works

---

## Commit Strategy

Reorganize in logical commits:

1. `create new directory structure`
2. `move shell configuration to shell/`
3. `move functions to functions/`
4. `import app configs from tmp-dotfiles`
5. `move documentation to docs/`
6. `rename initialize.sh to install.sh`
7. `remove legacy files and directories`
8. `create shell/init.sh entry point`
9. `create injection functions`
10. `update Makefile for new structure`
11. `update documentation`
