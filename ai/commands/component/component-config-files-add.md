---
description: Add config file generation support to a component
argument_hints:
  - ghostty
  - tmux
  - neovim
  - git
---

Add config file support to component: $ARGUMENTS

## Overview

This command adds the `files:` section to a component's config.yaml, enabling acorn to generate config files in various formats (JSON, YAML, TOML, INI, Ghostty, etc.).

## Schema Reference

### files: Array Structure

```yaml
files:
  - target: "${XDG_CONFIG_HOME}/ghostty/config"  # Output path (supports env vars)
    format: ghostty                               # Format identifier
    schema:                                       # Field definitions with types
      theme:
        type: string
        default: "Catppuccin Mocha"
      font-size:
        type: int
        default: 14
      keybind:
        type: list
        items: string
    values:                                       # Actual values to write
      theme: "Catppuccin Mocha"
      font-size: 14
      keybind:
        - "super+d=new_split:right"
```

### Supported Formats

| Format | Output Style | Use Cases |
|--------|-------------|-----------|
| `json` | `{"key": "value"}` | VS Code, package.json |
| `yaml` | `key: value` | k8s, docker-compose |
| `toml` | `key = "value"` | Cargo.toml, pyproject.toml |
| `ini` | `[section]\nkey=value` | git config, editorconfig |
| `ghostty` | `key = value` (no sections) | Ghostty terminal |
| `keyvalue` | `key=value` | Simple configs |

### Field Types

| Type | YAML Example | Output |
|------|--------------|--------|
| `string` | `"value"` | Quoted if spaces |
| `int` | `14` | Unquoted number |
| `bool` | `true` | `true`/`false` |
| `float` | `1.5` | Decimal number |
| `list` | `[a, b, c]` | Multiple lines (format-specific) |
| `map` | `{k: v}` | Nested structure |

## Instructions

### 1. Read Existing Config

```bash
cat components/$ARGUMENTS/config.yaml 2>/dev/null
cat components/$ARGUMENTS/config/* 2>/dev/null
```

### 2. Identify Config Files

List all config files the component manages:
- What format is each file?
- Where is the target location?
- What settings does it contain?

### 3. Define Schema

For each config file:
1. List all keys/settings
2. Determine the type for each
3. Set sensible defaults
4. Document any special handling (lists, etc.)

### 4. Add files: Section

Add to the component's config.yaml:

```yaml
files:
  - target: "<XDG path to config>"
    format: "<format identifier>"
    schema:
      <key>:
        type: <string|int|bool|float|list|map>
        default: <default value>
    values:
      <key>: <actual value>
```

### 5. Implement Format Writer (if new)

If the format doesn't exist yet, create:
- `internal/configfile/<format>.go` - Writer implementation
- Add to format registry in `internal/configfile/writer.go`

### 6. Test Generation

```bash
acorn $ARGUMENTS config show    # Show what would be written
acorn $ARGUMENTS config write   # Write the config file
```

## Design Decisions

- **No comment preservation**: Values only, comments lost on regeneration
- **config.yaml wins**: Direct edits to target file are overwritten
- **One format at a time**: Implement formats as needed per component

## Example: Ghostty

```yaml
name: ghostty
description: Ghostty terminal emulator

files:
  - target: "${XDG_CONFIG_HOME}/ghostty/config"
    format: ghostty
    schema:
      theme: { type: string, default: "Catppuccin Mocha" }
      font-family: { type: string, default: "JetBrainsMono Nerd Font" }
      font-size: { type: int, default: 14 }
      font-thicken: { type: bool, default: true }
      keybind: { type: list, items: string }
    values:
      theme: "Catppuccin Mocha"
      font-family: "JetBrainsMono Nerd Font"
      font-size: 14
      font-thicken: true
      keybind:
        - "super+d=new_split:right"
        - "super+shift+d=new_split:down"
```

## Files to Modify

### Schema Extension
- `internal/componentconfig/schema.go` - Add FileConfig types

### Format Writers
- `internal/configfile/writer.go` - Writer interface
- `internal/configfile/<format>.go` - Format implementation

### Integration
- `internal/componentconfig/loader.go` - Handle Files field
- `internal/shell/shell.go` - Call file generation

### Component
- `components/$ARGUMENTS/config.yaml` - Add files section
