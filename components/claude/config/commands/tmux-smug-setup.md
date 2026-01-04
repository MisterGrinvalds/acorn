---
description: Set up and configure smug session manager
argument-hint: [action: install|new|list]
allowed-tools: Read, Write, Edit, Bash
---

## Task

Help the user set up and use smug for persistent tmux session configurations.

## What is Smug?

Smug is a session manager for tmux that uses YAML files to define session layouts. Benefits:
- **Reproducible**: Same layout every time
- **Version controlled**: Configs live in dotfiles
- **Parameterized**: Pass variables to sessions

## Actions

Based on `$ARGUMENTS`:

### install
Install smug and link dotfiles configs:
```bash
# macOS
brew install smug

# Linux (with Go)
go install github.com/ivaaaan/smug@latest

# Link configs from dotfiles
smug_link_configs
```

### new
Create a new smug session config:
```bash
smug_new <session-name>
```

This creates `~/.config/smug/<name>.yml` with a template.

### list
Show available sessions:
```bash
smug_list
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

## Dotfiles Integration

Available functions:
- `smug_list` - List sessions with descriptions
- `smug_start [name]` - Start with fzf selection
- `smug_stop [name]` - Stop with fzf selection
- `smug_new <name>` - Create from template
- `smug_edit [name]` - Edit config
- `smug_install` - Install smug
- `smug_link_configs` - Link configs from dotfiles

## Context

@components/tmux/functions.sh
@components/tmux/config/smug

## Tips

1. Store smug configs in `components/tmux/config/smug/` for version control
2. Use `smug_link_configs` to sync to `~/.config/smug/`
3. Add descriptions with `# smug session: description` at top of file
4. Use variables: `smug start myproject project_path=~/work/foo`
