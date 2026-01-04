---
description: Create a new smug session configuration
argument-hint: <session-name> [root-path]
allowed-tools: Bash, Read, Write, Edit
---

## Task

Create a new smug session configuration with an interactive or templated approach.

## Execution

### Using Shell Function

```bash
smug_new <session-name>
```

### Manual Creation

1. Determine session name from `$ARGUMENTS` (first word)
2. Determine root path (second argument or current directory)
3. Create config file at `${SMUG_CONFIG_DIR:-$HOME/.config/smug}/<name>.yml`

## Template

```yaml
# smug session: <description>
session: <name>
root: <path>
attach: true

# Environment variables (optional)
# env:
#   PROJECT_ENV: development

# Commands to run before session starts (optional)
# before_start:
#   - docker-compose up -d

windows:
  - name: editor
    commands:
      - nvim .

  - name: terminal
    # Single pane with shell

  # - name: split
  #   panes:
  #     - type: horizontal
  #       commands:
  #         - echo "top pane"
  #     - commands:
  #         - echo "bottom pane"

# Commands to run when session stops (optional)
# stop:
#   - docker-compose down
```

## Project Type Detection

Detect project type and suggest appropriate windows:

| Indicator | Type | Suggested Windows |
|-----------|------|-------------------|
| `go.mod` | Go | editor, server, test |
| `package.json` | Node | editor, dev-server, terminal |
| `pyproject.toml` | Python | editor, server, terminal |
| `Cargo.toml` | Rust | editor, cargo-watch, terminal |
| `Makefile` | General | editor, make, terminal |

## After Creation

Remind user to:
1. Edit the config: `smug_edit <name>`
2. Test it: `smug start <name>`
3. Sync to git: `smug_push "Add <name> session"`

## Context

@components/tmux/functions.sh
