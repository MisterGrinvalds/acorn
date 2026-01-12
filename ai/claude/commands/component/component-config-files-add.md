---
description: Add config file generation support to a component
argument_hints:
  - ghostty
  - tmux
  - iterm2
  - neovim
---

Add config file support to component: $ARGUMENTS

## Overview

This command adds the `files:` section to a component's config.yaml. This is the **PRIMARY** way to manage tool config files - configs are generated from declarative values and symlinked to target locations.

**Prefer `files:` over `sync_files:`** for all tool configs except credentials/SSH.

## Quick Reference

### files: Section Structure

```yaml
files:
  - target: "${XDG_CONFIG_HOME:-$HOME/.config}/tool/config"
    format: json  # json, yaml, toml, ghostty, tmux, iterm2
    values:
      key1: "value"
      key2: true
      nested:
        key: "value"
```

### Available Formats

| Format | Output | Example Use |
|--------|--------|-------------|
| `json` | `{"key": "value"}` | VS Code settings, package.json |
| `yaml` | `key: value` | k8s configs, docker-compose |
| `toml` | `key = "value"` | Cargo.toml, pyproject.toml |
| `ghostty` | `key = value` | Ghostty terminal |
| `tmux` | `set -g key value` | tmux.conf |
| `iterm2` | Dynamic Profile JSON | iTerm2 profiles |

## Instructions

### 1. Read Existing Config

```bash
cat internal/componentconfig/config/$ARGUMENTS/config.yaml
```

### 2. Identify Config Files to Generate

- What config file does the tool use?
- Where is the target location (XDG path)?
- What format is it? (Check if format writer exists)

### 3. Add files: Section

Add to the component's config.yaml:

```yaml
files:
  - target: "${XDG_CONFIG_HOME:-$HOME/.config}/$ARGUMENTS/config"
    format: json
    values:
      setting1: "value1"
      setting2: true
```

### 4. Test Generation

```bash
acorn shell generate        # Generate all config files
ls generated/$ARGUMENTS/    # Check output
acorn sync link             # Create symlinks
```

## Examples

### JSON Config (VS Code style)

```yaml
files:
  - target: "${HOME}/Library/Application Support/Code/User/settings.json"
    format: json
    values:
      "editor.fontSize": 14
      "editor.fontFamily": "JetBrainsMono Nerd Font"
      "workbench.colorTheme": "Catppuccin Mocha"
```

### Ghostty Config

```yaml
files:
  - target: "${XDG_CONFIG_HOME:-$HOME/.config}/ghostty/config"
    format: ghostty
    values:
      theme: "Catppuccin Mocha"
      font-family: "JetBrainsMono Nerd Font"
      font-size: 14
      keybind:
        - "super+d=new_split:right"
        - "super+shift+d=new_split:down"
```

### iTerm2 Dynamic Profile

```yaml
files:
  - target: "${HOME}/Library/Application Support/iTerm2/DynamicProfiles/profile.json"
    format: iterm2
    values:
      profile:
        name: "shell-profile"
        guid: "shell-profile-001"
        parent: "Default"
      font:
        family: "JetBrainsMonoNF-Regular"
        size: 14
      colors:
        scheme: "catppuccin-mocha"
```

### tmux Config

```yaml
files:
  - target: "${XDG_CONFIG_HOME:-$HOME/.config}/tmux/tmux.conf"
    format: tmux
    values:
      set_g:
        mouse: true
        base-index: 1
        prefix: "C-a"
      bind:
        r: 'source-file ~/.config/tmux/tmux.conf'
        v: 'split-window -h'
```

## Migrating from Static Files

If the component currently uses `sync_files:` with static files:

### Before (Static)
```yaml
sync_files:
  - source: "config/tool/settings.json"
    target: "${HOME}/.config/tool/settings.json"
    mode: "symlink"
```

### After (Generated)
```yaml
files:
  - target: "${HOME}/.config/tool/settings.json"
    format: json
    values:
      # Move content from static file here
      setting1: "value1"
```

Then delete the static file in `config/tool/settings.json`.

## When NOT to Use files:

Use `sync_files:` instead for:
- SSH configs (need 600 permissions, use `mode: copy`)
- Credentials and secrets
- Git config with user-specific includes
- Files needing user customization overlay (`mode: merge`)

## Creating New Format Writers

If format doesn't exist, create `internal/configfile/<format>.go`:

```go
type MyFormatWriter struct{}

func init() {
    Register(&MyFormatWriter{})
}

func (w *MyFormatWriter) Format() string {
    return "myformat"
}

func (w *MyFormatWriter) Write(values map[string]interface{}) ([]byte, error) {
    // Generate output bytes from values
}
```

## Files Modified

- `internal/componentconfig/config/$ARGUMENTS/config.yaml` - Add files: section
- `generated/$ARGUMENTS/` - Output location after `acorn shell generate`
