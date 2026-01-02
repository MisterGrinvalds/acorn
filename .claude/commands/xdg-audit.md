---
description: Audit XDG Base Directory compliance
---

Audit XDG compliance across all components and the system.

## Checks to Perform

1. **XDG Environment Variables**
   - Verify all XDG variables are set:
     - XDG_CONFIG_HOME (~/.config)
     - XDG_DATA_HOME (~/.local/share)
     - XDG_CACHE_HOME (~/.cache)
     - XDG_STATE_HOME (~/.local/state)
     - XDG_RUNTIME_DIR
   - Check they have correct default values

2. **Directory Existence**
   - Verify each XDG directory exists
   - Check permissions are correct (especially XDG_RUNTIME_DIR should be 0700)

3. **Legacy Dotfile Detection**
   - Scan home directory for common dotfiles that should be in XDG locations:
     - ~/.bashrc, ~/.bash_profile, ~/.zshrc (should be bootstrap symlinks)
     - ~/.gitconfig (should be XDG_CONFIG_HOME/git/config)
     - ~/.vimrc (should be XDG_CONFIG_HOME/vim/vimrc)
     - ~/.tmux.conf (should be XDG_CONFIG_HOME/tmux/tmux.conf)
     - ~/.ssh/config (not XDG, but check permissions)
   - Report which can be migrated to XDG locations

4. **Component XDG Usage**
   - For each component, check if it declares XDG directories in component.yaml
   - Verify declared directories exist
   - Check component-specific config files are in correct locations

5. **Migration Recommendations**
   - List files that could be moved to XDG locations
   - Provide environment variable exports needed
   - Suggest symlinks for backwards compatibility

## Output Format

```
XDG Compliance Audit
====================

Environment Variables:
  ✓ XDG_CONFIG_HOME = ~/.config
  ✓ XDG_DATA_HOME = ~/.local/share
  ✓ XDG_CACHE_HOME = ~/.cache
  ✓ XDG_STATE_HOME = ~/.local/state
  ✓ XDG_RUNTIME_DIR = /var/run/user/1000

Directory Status:
  ✓ ~/.config exists (drwxr-xr-x)
  ✓ ~/.local/share exists (drwxr-xr-x)
  ✓ ~/.cache exists (drwxr-xr-x)
  ✓ ~/.local/state exists (drwxr-xr-x)

Legacy Dotfiles Found:
  ⚠ ~/.gitconfig - can migrate to ~/.config/git/config
  ⚠ ~/.vimrc - can migrate to ~/.config/vim/vimrc
  ✓ ~/.bashrc - is bootstrap symlink (correct)

Component XDG:
  shell: ~/.config/shell, ~/.local/state/shell
  git: ~/.config/git
  fzf: ~/.config/fzf

Recommendations:
  1. Move ~/.gitconfig to ~/.config/git/config
     Set: export GIT_CONFIG_GLOBAL="$XDG_CONFIG_HOME/git/config"
  2. Move ~/.vimrc to ~/.config/vim/vimrc
     Set: export VIMINIT="source $XDG_CONFIG_HOME/vim/vimrc"
```
