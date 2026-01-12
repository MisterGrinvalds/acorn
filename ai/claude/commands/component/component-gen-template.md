---
description: Generate the standardized template structure for a component
argument_hints:
  - tmux
  - go
  - python
  - node
  - kubernetes
---

Generate template structure for: $ARGUMENTS

## Instructions

Generate the standardized component structure in `internal/componentconfig/config/$ARGUMENTS/`.

### 1. Create Component Directory

```bash
mkdir -p internal/componentconfig/config/$ARGUMENTS
```

### 2. Generate config.yaml

Create `internal/componentconfig/config/$ARGUMENTS/config.yaml`:

```yaml
name: $ARGUMENTS
description: <description>
version: 1.0.0
platforms: [darwin, linux]

# Environment variables
env:
  # Example: TOOL_HOME: "${XDG_CONFIG_HOME:-$HOME/.config}/$ARGUMENTS"

# Shell aliases
aliases:
  # Example: alias: "command"

# Shell functions (ONLY for operations requiring shell state)
# Use these sparingly - prefer acorn commands instead
shell_functions:
  # Example functions that MUST stay in shell:
  # - cd wrappers (change directory)
  # - source/activation (modify shell environment)
  # - fzf integration (interactive selection)
  # - tmux attach (session attachment)

# GENERATED CONFIG FILES (primary pattern for tool configs)
# Configs are generated to generated/$ARGUMENTS/ and symlinked
files:
  - target: "${XDG_CONFIG_HOME:-$HOME/.config}/$ARGUMENTS/config"
    format: json  # Options: json, yaml, toml, ghostty, tmux, iterm2
    values:
      # Tool-specific settings
      setting1: "value1"
      setting2: true

# STATIC FILE SYNC (use sparingly)
# Only for: SSH configs, credentials, files needing 600 permissions
# sync_files:
#   - source: "config/$ARGUMENTS/file"
#     target: "${HOME}/.file"
#     mode: symlink  # or copy, merge

# Installation configuration
install:
  tools:
    - name: $ARGUMENTS
      description: <what it does>
      check: "command -v $ARGUMENTS"
      methods:
        darwin:
          type: brew
          package: $ARGUMENTS
        linux:
          type: apt
          package: $ARGUMENTS
```

### 3. Create Claude Integration (Optional)

```bash
# Create command directory
mkdir -p ai/claude/commands/$ARGUMENTS

# Create agent (optional)
# ai/claude/agents/$ARGUMENTS-expert.md
```

### 4. Report Created Files

Output:
```
Generated Template: $ARGUMENTS
==============================

Created:
  - internal/componentconfig/config/$ARGUMENTS/config.yaml

Config Strategy:
  - files: section for GENERATED configs (preferred)
  - sync_files: section only for credentials/SSH

Next steps:
  1. Edit config.yaml to add component-specific values
  2. Run: acorn shell generate
  3. Check: generated/$ARGUMENTS/ for output
  4. Run: acorn sync link to create symlinks
  5. Optionally run:
     - /component-gen-agent $ARGUMENTS
     - /component-gen-commands $ARGUMENTS
```

## Config Strategy Decision Tree

```
Does the tool need a config file?
├── Yes → Is it a standard format (JSON, YAML, TOML)?
│   ├── Yes → Use files: section with appropriate format
│   └── No → Is there a format writer? (ghostty, tmux, iterm2)
│       ├── Yes → Use files: section with that format
│       └── No → Create new format writer in internal/configfile/
└── No → Only use env:, aliases:, shell_functions:
```

## Available Format Writers

| Format | File | Use For |
|--------|------|---------|
| `json` | writer.go | Generic JSON |
| `yaml` | writer.go | Generic YAML |
| `toml` | writer.go | Generic TOML |
| `ghostty` | ghostty.go | Ghostty terminal |
| `tmux` | tmux.go | tmux.conf |
| `iterm2` | iterm2.go | iTerm2 profiles |
