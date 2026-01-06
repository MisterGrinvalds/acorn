---
description: Set up and configure smug session manager
argument-hint: [action: install|init|new|list|sync]
allowed-tools: Read, Write, Edit, Bash
---

## Task

Help the user set up and use smug for persistent tmux session configurations with git sync.

## What is Smug?

Smug is a session manager for tmux that uses YAML files to define session layouts. Benefits:
- **Reproducible**: Same layout every time
- **Git synced**: Configs sync across machines via git repo
- **Parameterized**: Pass variables to sessions

## Actions

Based on `$ARGUMENTS`:

### install
Install smug:
```bash
# macOS
brew install smug

# Linux (with Go)
go install github.com/ivaaaan/smug@latest
```

### init
Initialize the smug sessions git repo (first-time setup):
```bash
smug_repo_init
```
This clones the fmux repo and links it to `~/.config/smug/`.

### new
Create a new smug session config:
```bash
smug_new <session-name>
```
Creates `~/.config/smug/<name>.yml` with a template.

### list
Show available sessions:
```bash
smug_list
```

### sync
Sync sessions with git remote:
```bash
smug_sync   # Full sync (pull + push)
smug_pull   # Pull latest from remote
smug_push   # Commit and push changes
smug_status # Show git status
```

## Config Format

Location: `~/.config/smug/<name>.yml`

```yaml
# smug session: my-project
session: my-project
root: ~/projects/my-project
attach: true

# Environment variables
env:
  PROJECT_ENV: development

# Before session starts
before_start:
  - docker-compose up -d

# Window definitions
windows:
  - name: editor
    commands:
      - nvim .

  - name: server
    commands:
      - npm run dev

  - name: terminal
    panes:
      - type: horizontal
        commands:
          - echo "Ready"
      - commands:
          - git status

# After session stops
stop:
  - docker-compose down
```

## Smug Commands

```bash
smug start <name>     # Start session
smug start <name> -a  # Start attached
smug stop <name>      # Stop session
smug print <name>     # Print tmux commands (debug)
```

## Available Functions

Session management:
- `smug_list` - List sessions with descriptions
- `smug_start [name]` - Start with fzf selection
- `smug_stop [name]` - Stop with fzf selection
- `smug_new <name>` - Create from template
- `smug_edit [name]` - Edit config

Git sync:
- `smug_repo_init` - Clone/update sessions git repo
- `smug_status` - Show git repo status
- `smug_pull` - Pull latest sessions
- `smug_push [msg]` - Commit and push changes
- `smug_sync` - Full sync (pull + push)

Setup:
- `smug_install` - Install smug
- `smug_link_configs` - Link configs (from git repo or dotfiles)

## Context

@components/tmux/config.yaml
@components/tmux/config/smug/project.yml

## Tips

1. Run `smug_repo_init` on new machines to sync sessions
2. After creating/editing sessions, run `smug_push` to sync
3. Add descriptions with `# smug session: description` at top of file
4. Use variables: `smug start myproject project_path=~/work/foo`
