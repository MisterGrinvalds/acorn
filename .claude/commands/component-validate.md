---
description: Validate a component against template standards
---

Validate component: $ARGUMENTS

## Instructions

Thoroughly validate the specified component against all standards.

### 1. Check Directory Structure

Verify `components/$ARGUMENTS/` exists and contains:
- [ ] `component.yaml` (REQUIRED)
- [ ] At least one of: `env.sh`, `aliases.sh`, `functions.sh`
- [ ] Files use correct naming

### 2. Validate component.yaml

**Required fields:**
- [ ] `name` - matches directory name
- [ ] `version` - valid semver (X.Y.Z)
- [ ] `description` - non-empty string
- [ ] `category` - one of: core, dev, cloud, ai, database

**Optional fields (validate if present):**
- [ ] `requires.tools` - list of strings
- [ ] `requires.components` - list of existing components
- [ ] `provides.aliases` - list of strings
- [ ] `provides.functions` - list of strings
- [ ] `platforms` - subset of [darwin, linux]
- [ ] `shells` - subset of [bash, zsh]

### 3. Validate Shell Scripts

For each .sh file:
```bash
bash -n components/$ARGUMENTS/<file>.sh
```

Check for:
- [ ] Syntax errors (bash -n)
- [ ] Proper shebang (#!/bin/sh or #!/bin/bash)
- [ ] Local variables in functions
- [ ] Consistent indentation

### 4. Check Dependencies

Verify all declared dependencies:
- [ ] Required tools exist (`command -v <tool>`)
- [ ] Required components exist in `components/`

### 5. Cross-Reference Provides

Verify that declared provides match actual content:
- [ ] All listed aliases exist in aliases.sh
- [ ] All listed functions exist in functions.sh

### 6. Check Documentation

- [ ] README.md exists
- [ ] README documents all public functions
- [ ] README includes usage examples

## Output Format

```
Component Validation: $ARGUMENTS
================================

Structure:           [PASS/FAIL]
  - component.yaml   [OK/MISSING]
  - env.sh           [OK/MISSING/N/A]
  - aliases.sh       [OK/MISSING/N/A]
  - functions.sh     [OK/MISSING/N/A]
  - README.md        [OK/MISSING]

Metadata:            [PASS/FAIL]
  - name             [OK/ERROR: <reason>]
  - version          [OK/ERROR: <reason>]
  - description      [OK/ERROR: <reason>]
  - category         [OK/ERROR: <reason>]

Syntax Check:        [PASS/FAIL]
  - env.sh           [OK/ERROR at line N]
  - aliases.sh       [OK/ERROR at line N]
  - functions.sh     [OK/ERROR at line N]

Dependencies:        [PASS/WARN/FAIL]
  - <tool>           [INSTALLED/MISSING]
  - <component>      [EXISTS/MISSING]

Provides Match:      [PASS/FAIL]
  - aliases          [N declared, N found]
  - functions        [N declared, N found]

Overall:             [VALID/INVALID]

Issues Found:
  1. <issue description>
  2. <issue description>

Recommendations:
  - <suggestion>
```
