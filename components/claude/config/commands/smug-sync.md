---
description: Sync smug sessions with git remote
argument-hint: [action: pull|push|status|init]
allowed-tools: Bash
---

## Task

Sync smug tmux session configurations with the git remote repository.

## Actions

Based on `$ARGUMENTS`:

### (no args) or sync
Full sync - pull latest and push local changes:
```bash
smug_sync
```

### pull
Pull latest sessions from remote:
```bash
smug_pull
```

### push
Commit and push local session changes:
```bash
smug_push "commit message"
```

### status
Show git repo status and list sessions:
```bash
smug_status
```

### init
Initialize the sessions repo (first-time setup):
```bash
smug_repo_init
```

## Environment

- `SMUG_REPO` - Git repo URL (default: https://github.com/MisterGrinvalds/fmux.git)
- `SMUG_REPO_DIR` - Local repo location (default: ~/.local/share/smug-sessions)
- `SMUG_CONFIG_DIR` - Smug config dir (default: ~/.config/smug)

## Workflow

1. **New machine setup:**
   ```bash
   smug_repo_init
   ```

2. **After editing sessions:**
   ```bash
   smug_push "Add new project session"
   ```

3. **Sync on other machines:**
   ```bash
   smug_pull
   ```

4. **Regular sync:**
   ```bash
   smug_sync
   ```
