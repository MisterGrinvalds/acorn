---
description: Check dotfiles for drift from repository
---

Run drift detection on the dotfiles repository.

## Checks to Perform

1. **Git Status**
   - Check for uncommitted changes in the repository
   - List modified, added, and deleted files
   - Report untracked files

2. **Remote Comparison**
   - Fetch latest from origin (quiet)
   - Check commits ahead/behind origin
   - Report if branch needs push or pull

3. **Component State Audit**
   - Compare loaded components vs available components
   - Check for components that failed to load
   - Report any disabled components

4. **XDG Compliance Check**
   - Verify XDG environment variables are set
   - Check for legacy dotfiles (~/.bashrc, ~/.zshrc) that should be symlinks
   - Report any config files in non-XDG locations

5. **Tool Version Check**
   - For key tools (git, fzf, yq), report installed versions
   - Note any outdated tools if version info available

## Output Format

```
Dotfiles Drift Report
=====================

Repository: /path/to/bash-profile
Branch: main

Git Status:
  Commits: 2 behind origin, 1 ahead
  Modified: 3 files
  Untracked: 1 file

Component Status:
  Loaded: 5 components
  Disabled: 0 components
  Failed: 1 component (python: uv not installed)

XDG Compliance:
  ✓ XDG directories configured
  ⚠ ~/.bashrc exists (should be symlink)

Recommendations:
  - Run `git pull --rebase` to sync with origin
  - Install uv with: brew install uv
  - Convert ~/.bashrc to symlink
```
