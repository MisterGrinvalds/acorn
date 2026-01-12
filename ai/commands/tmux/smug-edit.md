---
description: Edit a smug session configuration
argument-hint: [session-name]
allowed-tools: Bash, Read, Write, Edit
---

## Task

Edit an existing smug session configuration file.

## Execution

### Using Shell Function

```bash
smug_edit <session-name>
# Or with fzf selection:
smug_edit
```

### Direct Edit

```bash
${EDITOR:-nvim} "${SMUG_CONFIG_DIR:-$HOME/.config/smug}/$ARGUMENTS.yml"
```

## Validation

After editing, validate the YAML:
```bash
yq eval '.' "${SMUG_CONFIG_DIR:-$HOME/.config/smug}/$ARGUMENTS.yml" > /dev/null && echo "Valid YAML" || echo "Invalid YAML"
```

## Common Modifications

### Add a Window
```yaml
windows:
  # ... existing windows ...
  - name: new-window
    commands:
      - echo "New window"
```

### Add Split Panes
```yaml
  - name: split-window
    panes:
      - type: horizontal  # or vertical
        commands:
          - top
      - commands:
          - htop
```

### Add Environment Variables
```yaml
env:
  NODE_ENV: development
  DEBUG: "true"
```

### Add Startup Commands
```yaml
before_start:
  - docker-compose up -d
  - sleep 2
```

## After Editing

Remind user to sync changes:
```bash
smug_push "Update <session-name> configuration"
```

## Context

@components/tmux/config.yaml
