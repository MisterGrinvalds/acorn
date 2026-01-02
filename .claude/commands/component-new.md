---
description: Create a new component from template
---

Create a new component in the bash-profile repository.

Component name: $ARGUMENTS

## Instructions

1. Copy the template directory to create a new component:
   - Source: `components/_template/`
   - Target: `components/$ARGUMENTS/`

2. Update `component.yaml` with:
   - Correct name (the component name)
   - Appropriate version (start with 1.0.0)
   - Description of what the component provides
   - Category: core, dev, cloud, ai, or database
   - Required tools (CLI commands that must be installed)
   - Required components (other components that must load first)
   - What the component provides (aliases, functions, completions)
   - XDG directories it uses (if any)
   - Platform and shell support

3. Implement the component files:
   - `env.sh`: Environment variables and PATH modifications
   - `aliases.sh`: Shell aliases
   - `functions.sh`: Helper functions
   - `completions.sh`: Tab completion setup
   - `setup.sh`: Installation/configuration script (optional)

4. Test the component:
   - Run `bash -n` on all .sh files to check syntax
   - Source the bootstrap and verify the component loads

5. Update README.md with documentation for the new component

## Checklist

- [ ] Template copied to `components/$ARGUMENTS/`
- [ ] `component.yaml` updated with correct metadata
- [ ] At least one of env.sh, aliases.sh, or functions.sh implemented
- [ ] All .sh files pass `bash -n` syntax check
- [ ] Component loads without errors
