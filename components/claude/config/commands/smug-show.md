---
description: Show a smug session configuration
argument-hint: <session-name>
allowed-tools: Bash, Read
---

## Task

Display the contents of a smug session configuration file.

## Execution

If `$ARGUMENTS` provided:
```bash
cat "${SMUG_CONFIG_DIR:-$HOME/.config/smug}/$ARGUMENTS.yml"
```

If no arguments, use fzf selection:
```bash
smug_edit  # Opens in editor, but shows content
# Or list and let user choose
smug_list
```

## Output

Show the full YAML configuration with syntax highlighting if available:
```bash
cat "${SMUG_CONFIG_DIR:-$HOME/.config/smug}/$ARGUMENTS.yml" | bat -l yaml 2>/dev/null || cat "${SMUG_CONFIG_DIR:-$HOME/.config/smug}/$ARGUMENTS.yml"
```

## Config Structure Reference

```yaml
session: name
root: ~/path
attach: true
env:
  VAR: value
before_start:
  - command
windows:
  - name: window-name
    commands:
      - command
    panes:
      - commands:
          - command
stop:
  - cleanup-command
```

## Context

@components/tmux/functions.sh
