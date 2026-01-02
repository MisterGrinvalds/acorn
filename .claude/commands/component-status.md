---
description: Check health status of all components
---

Audit all components in the bash-profile repository for health and correctness.

## Checks to Perform

1. **Discover all components**
   - List all directories in `components/` (excluding `_template`)
   - Report which have valid `component.yaml` files

2. **Validate component.yaml files**
   - Use `yq` to parse each component.yaml
   - Verify required fields: name, version, description, category
   - Check that referenced tools exist (via `command -v`)

3. **Check shell file syntax**
   - Run `bash -n` on all .sh files in each component
   - Report any syntax errors

4. **Verify dependencies**
   - For each component's `requires.tools`, check if installed
   - For each component's `requires.components`, verify the component exists

5. **Report summary**
   - Total components found
   - Components with issues
   - Missing tools
   - Syntax errors

## Output Format

```
Component Status Report
=======================

Components Found: N

✓ shell (1.0.0) - Core shell configuration
  Tools: all present
  Dependencies: none

✓ git (1.0.0) - Git version control
  Tools: git present
  Dependencies: shell (ok)

⚠ python (1.0.0) - Python development
  Tools: python3 present, uv MISSING
  Dependencies: shell (ok)

Summary:
- N components checked
- N healthy
- N with warnings
- N with errors
```
