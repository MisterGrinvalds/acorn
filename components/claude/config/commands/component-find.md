---
description: Find existing components that may overlap with a proposed new component
---

Find existing components that might overlap with: $ARGUMENTS

## Instructions

Search the codebase for existing functionality that overlaps with the proposed component.

### 1. Search Component Names

Look for components with similar names in `components/` directory.

### 2. Search Function Names

Search for functions that might provide similar functionality:
```bash
grep -rh "^[a-z_]*() {" components/*/functions.sh | sort -u
```

### 3. Search Aliases

Look for aliases that might conflict:
```bash
grep -rh "^alias " components/*/aliases.sh | sort -u
```

### 4. Search Tool Integrations

If the proposed component integrates with a specific tool, search for existing integrations:
```bash
grep -rl "$ARGUMENTS" components/*/component.yaml
grep -rl "$ARGUMENTS" components/*/functions.sh
```

### 5. Report Findings

Report:
- Any components with similar names
- Overlapping functions or aliases
- Existing integrations with the same tools
- Recommendation: extend existing component vs create new one

## Output Format

```
Overlap Analysis for: $ARGUMENTS
================================

Similar Components:
  - <component> - <description>

Overlapping Functions:
  - <function> in <component>

Overlapping Aliases:
  - <alias> in <component>

Tool Integrations Found:
  - <tool> integrated in <component>

Recommendation:
  [ ] Create new component
  [ ] Extend existing component: <name>
  [ ] Merge with: <name>
```
