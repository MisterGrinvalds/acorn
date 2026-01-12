---
description: Validate a generated component has all required elements
argument_hints:
  - tmux
  - go
  - python
  - node
  - kubernetes
---

Validate component: $ARGUMENTS

## Instructions

Verify the component in `components/$ARGUMENTS/` is complete and valid.

### 1. Check Directory Structure

Verify all required directories exist:
```
components/$ARGUMENTS/
├── shell/
└── config/
```

Note: Installation is configured via `install:` section in config.yaml (not shell scripts)

### 2. Validate Config

**Check config.yaml:**
- [ ] File exists at `components/$ARGUMENTS/config.yaml`
- [ ] Valid YAML syntax
- [ ] All required fields present:
  - [ ] `name` matches directory
  - [ ] `description` is non-empty
  - [ ] `version` is valid semver
  - [ ] `category` is valid (core/dev/cloud/ai/database)
  - [ ] `platforms` is valid array
  - [ ] `shells` is valid array

### 3. Validate Shell Scripts

Check each shell script in `components/$ARGUMENTS/shell/`:

**env.sh:**
- [ ] File exists
- [ ] Valid shell syntax (`bash -n`)
- [ ] Has shebang line
- [ ] Exports environment variables correctly

**aliases.sh:**
- [ ] File exists
- [ ] Valid shell syntax
- [ ] Aliases defined correctly

**functions.sh:**
- [ ] File exists
- [ ] Valid shell syntax
- [ ] Functions use `local` for variables

**completions.sh:**
- [ ] File exists
- [ ] Valid shell syntax
- [ ] Handles both bash and zsh

### 4. Validate Claude Integration (Centralized)

**Agent file:**
- [ ] `ai/agents/$ARGUMENTS-expert.md` exists
- [ ] Valid YAML frontmatter
- [ ] Required sections present:
  - [ ] Core Competencies
  - [ ] Key Concepts
  - [ ] Available Shell Functions
  - [ ] Best Practices
  - [ ] Approach

**Command files:**
- [ ] `ai/commands/$ARGUMENTS/` directory exists
- [ ] At least `explain.md` or `$ARGUMENTS-explain.md` exists
- [ ] At least `coach.md` or `$ARGUMENTS-coach.md` exists
- [ ] All commands have valid frontmatter
- [ ] All commands have clear instructions

### 5. Validate Installation Config

**Check config.yaml `install:` section:**
- [ ] `install:` section exists (if component requires tools)
- [ ] Each tool has required fields:
  - [ ] `name` - tool name
  - [ ] `check` - command to verify installation
  - [ ] `methods` - platform-specific install methods
- [ ] Each method has:
  - [ ] `type` - valid type (brew, apt, npm, pip, go, curl)
  - [ ] `package` - package name if different from tool name
- [ ] Prerequisites are valid (`requires` references existing components)

**Test installation dry-run:**
```bash
acorn $ARGUMENTS install --dry-run
```

### 6. Run Syntax Checks

```bash
# Validate YAML config
yq '.' components/$ARGUMENTS/config.yaml >/dev/null 2>&1 || echo "YAML error: config.yaml"

# Validate shell scripts
for f in components/$ARGUMENTS/shell/*.sh; do
    bash -n "$f" 2>&1 || echo "Syntax error: $f"
done

# Test installation config (if exists)
acorn $ARGUMENTS install --dry-run 2>&1 || echo "Install config error"
```

### 7. Output Report

```
Component Validation: $ARGUMENTS
================================

Structure:                    [PASS/FAIL]
  - config.yaml               [OK/MISSING]
  - shell/                    [OK/MISSING]
  - config/                   [OK/MISSING]

Config File:                  [PASS/FAIL]
  - config.yaml               [VALID/INVALID: <reason>]

Shell Scripts:                [PASS/FAIL]
  - env.sh                    [VALID/ERROR: line N]
  - aliases.sh                [VALID/ERROR: line N]
  - functions.sh              [VALID/ERROR: line N]
  - completions.sh            [VALID/ERROR: line N]

Claude Integration:           [PASS/FAIL]
  - Agent                     [FOUND/MISSING]
    Location: ai/agents/$ARGUMENTS-expert.md
  - Commands dir              [FOUND/MISSING]
    Location: ai/commands/$ARGUMENTS/
  - explain command           [FOUND/MISSING]
  - coach command             [FOUND/MISSING]
  - Other commands            [N found]

Installation Config:          [PASS/FAIL/N/A]
  - install: section          [FOUND/MISSING]
  - Tools configured          [N tools]
  - Dry-run test              [PASS/FAIL]
  Test: acorn $ARGUMENTS install --dry-run

Overall:                      [VALID/INVALID]

Issues Found:
  1. <issue description>
  2. <issue description>

To fix issues, run:
  - /component:gen-<type> $ARGUMENTS
```
