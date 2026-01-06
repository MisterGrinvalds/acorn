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
├── ai/
│   └── claude/
│       ├── agents/
│       └── commands/
├── shell/
├── install/
└── config/
```

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

### 4. Validate AI Integration

**Agent file:**
- [ ] `ai/claude/agents/$ARGUMENTS-expert.md` exists
- [ ] Valid YAML frontmatter
- [ ] Required sections present:
  - [ ] Core Competencies
  - [ ] Key Concepts
  - [ ] Available Shell Functions
  - [ ] Best Practices
  - [ ] Approach

**Command files:**
- [ ] `ai/claude/commands/$ARGUMENTS-explain.md` exists
- [ ] `ai/claude/commands/$ARGUMENTS-coach.md` exists
- [ ] All commands have valid frontmatter
- [ ] All commands have clear instructions

### 5. Validate Install Scripts

**install/install.sh:**
- [ ] File exists
- [ ] Valid shell syntax
- [ ] Executable permissions (or can be set)
- [ ] Handles multiple package managers

**install/brew.yaml:**
- [ ] File exists if macOS packages needed
- [ ] Valid YAML syntax
- [ ] Package names are correct

**install/apt.yaml:**
- [ ] File exists if Linux packages needed
- [ ] Valid YAML syntax
- [ ] Package names are correct

### 6. Run Syntax Checks

```bash
# Validate YAML files
for f in components/$ARGUMENTS/config.yaml components/$ARGUMENTS/install/*.yaml; do
    yq '.' "$f" >/dev/null 2>&1 || echo "YAML error: $f"
done

# Validate shell scripts
for f in components/$ARGUMENTS/shell/*.sh components/$ARGUMENTS/install/*.sh; do
    bash -n "$f" 2>&1 || echo "Syntax error: $f"
done
```

### 7. Output Report

```
Component Validation: $ARGUMENTS
================================

Structure:               [PASS/FAIL]
  - config.yaml          [OK/MISSING]
  - shell/               [OK/MISSING]
  - ai/claude/agents/    [OK/MISSING]
  - ai/claude/commands/  [OK/MISSING]
  - install/             [OK/MISSING]
  - config/              [OK/MISSING]

Config File:             [PASS/FAIL]
  - config.yaml          [VALID/INVALID: <reason>]

Shell Scripts:           [PASS/FAIL]
  - env.sh               [VALID/ERROR: line N]
  - aliases.sh           [VALID/ERROR: line N]
  - functions.sh         [VALID/ERROR: line N]
  - completions.sh       [VALID/ERROR: line N]

AI Integration:          [PASS/FAIL]
  - Agent                [FOUND/MISSING]
  - explain command      [FOUND/MISSING]
  - coach command        [FOUND/MISSING]
  - Other commands       [N found]

Install Scripts:         [PASS/FAIL]
  - install.sh           [VALID/ERROR]
  - brew.yaml            [VALID/N/A]
  - apt.yaml             [VALID/N/A]

Overall:                 [VALID/INVALID]

Issues Found:
  1. <issue description>
  2. <issue description>

To fix issues, run:
  - /component-gen-<type> $ARGUMENTS

Ready for injection:
  acorn ai inject $ARGUMENTS
```
