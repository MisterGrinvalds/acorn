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

Generate the standardized component structure directly in `components/$ARGUMENTS/`.

### 1. Create Directory Structure

Create the following directories:

```
components/$ARGUMENTS/
├── shell/
└── config/
```

Note:
- Installation is configured via `install:` section in config.yaml (not shell scripts)
- Claude agents/commands are centralized in `components/claude/config/`, not per-component

### 2. Read Existing Config (if any)

Check for existing configuration:
- `components/$ARGUMENTS/config.yaml` - current location
- `internal/componentconfig/config/$ARGUMENTS/config.yaml` - legacy Go location

### 3. Generate config.yaml

Create the component's config.yaml with this schema:

```yaml
name: $ARGUMENTS
description: <from source config or provide>
version: 1.0.0
category: <core|dev|cloud|ai|database>
platforms: [darwin, linux]
shells: [bash, zsh]

requires:
  tools: <from source>
  components: <from source>

xdg:
  config: $ARGUMENTS
  data: ""
  cache: ""
  state: ""

env: <from source env section>
paths: <from source paths section>
aliases: <from source aliases section>

# Shell functions that must stay in shell (cd, source, fzf, attach)
shell_functions: <approved functions only>

# Tool-specific configuration for Go to generate config files
# Uncomment and customize for your tool:
# tool_config:
#   setting1: value1
#   setting2: value2
```

### 4. Generate Shell Scripts

**shell/env.sh:**
```bash
#!/bin/sh
# $ARGUMENTS environment variables

<generate from config.env>
<generate PATH additions from config.paths>
```

**shell/aliases.sh:**
```bash
#!/bin/sh
# $ARGUMENTS aliases

<generate from config.aliases>
```

**shell/functions.sh:**
```bash
#!/bin/sh
# $ARGUMENTS functions
# Only functions that modify shell state (cd, source, fzf, attach)

<include approved shell_functions>
```

**shell/completions.sh:**
```bash
#!/bin/sh
# $ARGUMENTS completions

# Placeholder for completion setup
# Will be filled by /component-gen-completions
```

### 5. Create Placeholder Files

**config/.gitkeep** - placeholder for tool-specific config files

### 6. Create Claude Integration (Centralized)

Create the command subdirectory (agents and commands are centralized):
```bash
mkdir -p components/claude/config/commands/$ARGUMENTS
```

Note: Agent and commands will be created by:
- `/component-gen-agent $ARGUMENTS`
- `/component-gen-commands $ARGUMENTS`

### 7. Report Created Files

Output:
```
Generated Template: $ARGUMENTS
==============================

Created directories:
  - components/$ARGUMENTS/shell/
  - components/$ARGUMENTS/config/
  - components/claude/config/commands/$ARGUMENTS/

Created files:
  - components/$ARGUMENTS/config.yaml
  - components/$ARGUMENTS/shell/env.sh
  - components/$ARGUMENTS/shell/aliases.sh
  - components/$ARGUMENTS/shell/functions.sh
  - components/$ARGUMENTS/shell/completions.sh

Next steps:
  1. Run /component:gen-agent $ARGUMENTS
  2. Run /component:gen-commands $ARGUMENTS
  3. Run /component:gen-install $ARGUMENTS  (adds install: section to config.yaml)
  4. Run /component:gen-completions $ARGUMENTS
  5. Run /component:gen-validate $ARGUMENTS
```
