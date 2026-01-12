---
description: Analyze a component and show current functions, aliases, and config
argument_hints:
  - tmux
  - go
  - python
  - node
  - kubernetes
  - github
  - fzf
  - ollama
---

Analyze component: $ARGUMENTS

## Instructions

Analyze the specified component to understand its current state before refactoring.

### 1. Read Component Config

Read the component's config.yaml:

```bash
cat components/$ARGUMENTS/config.yaml 2>/dev/null
```

If not found, the component needs to be created from scratch.

### 2. Gather All Elements

From the config.yaml, extract and list:

**Environment Variables:**
- List all entries in the `env:` section
- Note any XDG path usage

**PATH Entries:**
- List all entries in the `paths:` section
- Note platform conditions (darwin/linux)

**Aliases:**
- List all entries in the `aliases:` section
- Format: `alias_name` -> `command`

**Shell Functions:**
- List all entries in the `shell_functions:` section
- These contain raw shell code (for cd, source, fzf, attach)

**Tool Config (if present):**
- List any `tool_config:` section entries
- These define typed config for Go to generate files

### 3. Check Existing AI Integration

Look for existing agents and commands in the component:

```bash
# Check for agent in component
ls components/$ARGUMENTS/ai/claude/agents/$ARGUMENTS-expert.md 2>/dev/null

# Check for commands in component
ls components/$ARGUMENTS/ai/claude/commands/$ARGUMENTS-*.md 2>/dev/null

# Also check legacy location (ai/)
ls ai/agents/$ARGUMENTS-*.md 2>/dev/null
ls ai/commands/$ARGUMENTS-*.md 2>/dev/null
```

### 4. Check Existing Shell Scripts

```bash
ls components/$ARGUMENTS/shell/*.sh 2>/dev/null
```

### 5. Check Installation Config

Check if config.yaml has an `install:` section:

```bash
# Check for install: section in config.yaml
grep -A5 "^install:" internal/componentconfig/config/$ARGUMENTS/config.yaml 2>/dev/null

# Or test with acorn CLI
acorn $ARGUMENTS install --dry-run 2>/dev/null
```

### 6. Present Summary Table

Output a summary in this format:

```
Component Analysis: $ARGUMENTS
================================

Config Source: components/$ARGUMENTS/config.yaml
Status: [EXISTS/MISSING]

Environment Variables: N entries
---------------------------------
  VARNAME: value

PATH Entries: N entries
-----------------------
  $PATH_ENTRY [condition]

Aliases: N entries
------------------
  alias -> command

Shell Functions: N entries
--------------------------
  function_name: [KEEP/REMOVE]
    Reason: <why keep or remove>

Tool Config: [EXISTS/MISSING]
-----------------------------
  <key settings if exists>

Existing Files:
---------------
  Shell:
    - env.sh        [EXISTS/MISSING]
    - aliases.sh    [EXISTS/MISSING]
    - functions.sh  [EXISTS/MISSING]
    - completions.sh [EXISTS/MISSING]

  AI Integration:
    - Agent         [EXISTS/MISSING/LEGACY]
    - Commands      N found [COMPONENT/LEGACY]

  Installation:
    - install: section in config.yaml  [EXISTS/MISSING]
    - Tools configured: N
    - Test: acorn $ARGUMENTS install --dry-run

  Config:
    - <tool-specific files>

Recommendations:
----------------
  - Keep: function1, function2 (used for cd/source/fzf/attach)
  - Remove: function3 (simple wrapper, use acorn instead)
  - Migrate: move agents/commands from legacy to component
```

### 7. Categorize Shell Functions

For each shell function, recommend:

**KEEP** - Functions that:
- Use `cd` to change directory
- Use `source` or `.` to activate environments
- Use `tmux attach` or similar session operations
- Integrate with fzf for interactive selection
- Open `$EDITOR` or modify shell state

**REMOVE** - Functions that:
- Are simple wrappers calling `acorn <component> <action>`
- Conflict with common CLI tools (fd, rg, bat)
- Are simple aliases disguised as functions
- Can be replaced by `acorn` commands

Ask the user to confirm which functions to keep before proceeding.
