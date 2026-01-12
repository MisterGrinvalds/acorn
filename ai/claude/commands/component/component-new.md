---
description: Create a new component from template
---

Create a new component in the acorn dotfiles system.

Component name: $ARGUMENTS

## Instructions

1. Create the component config directory:
   - Location: `internal/componentconfig/config/$ARGUMENTS/`
   - Create `config.yaml` with the schema below

2. Define the component in `config.yaml`:

```yaml
name: $ARGUMENTS
description: <description>
version: 1.0.0
platforms: [darwin, linux]

# Environment variables
env:
  EXAMPLE_VAR: "${HOME}/.example"

# Shell aliases
aliases:
  ex: "example-command"

# Shell functions (only for cd, source, fzf, attach operations)
shell_functions:
  example_func: |
    echo "This function needs shell state"

# GENERATED CONFIG FILES (primary pattern)
files:
  - target: "${XDG_CONFIG_HOME:-$HOME/.config}/$ARGUMENTS/config"
    format: json  # or yaml, toml, ghostty, tmux, iterm2
    values:
      setting1: "value1"
      setting2: true

# Installation
install:
  tools:
    - name: example-tool
      check: "command -v example-tool"
      methods:
        darwin:
          type: brew
          package: example-tool
        linux:
          type: apt
          package: example-tool
```

3. Choose config strategy:
   - **`files:` section** - For tool configs that should be generated (tmux, ghostty, vscode, iterm2)
   - **`sync_files:` section** - Only for credentials, SSH keys, or files needing user overlay

4. Test the component:
   - Run `go build ./...` to verify config.yaml is valid
   - Run `acorn shell generate` to generate config files
   - Check `generated/$ARGUMENTS/` for output

5. Optionally create Claude integration:
   - Agent: `ai/claude/agents/$ARGUMENTS-expert.md`
   - Commands: `ai/claude/commands/$ARGUMENTS/`

## Available Format Writers

| Format | Use For |
|--------|---------|
| `json` | JSON config files |
| `yaml` | YAML config files |
| `toml` | TOML config files |
| `ghostty` | Ghostty terminal config |
| `tmux` | tmux.conf |
| `iterm2` | iTerm2 dynamic profiles |

## Checklist

- [ ] `internal/componentconfig/config/$ARGUMENTS/config.yaml` created
- [ ] `files:` section defined for tool configs (preferred over static files)
- [ ] `go build ./...` succeeds
- [ ] `acorn shell generate` produces expected output
- [ ] Generated files appear in `generated/$ARGUMENTS/`
