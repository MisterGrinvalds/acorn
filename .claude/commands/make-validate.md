---
description: Validate Makefile syntax and structure
model: claude-sonnet-4-5-20250929
allowed-tools: Read, Bash
---

# Validate Makefile

Check the Makefile for syntax errors, best practices, and common issues.

## Task

1. Verify a Makefile exists in the project root
2. Run validation checks:
   - **Syntax**: Use `make -n` to dry-run and check for syntax errors
   - **PHONY targets**: Identify targets missing .PHONY declarations
   - **Documentation**: Find targets without help comments (##)
   - **Silent commands**: Suggest adding @ prefix where appropriate
   - **Variables**: Check for uninitialized or unused variables
   - **Naming**: Verify kebab-case convention for multi-word targets

3. Run `make help` to verify the help system works

4. Report findings:
   - ✓ Syntax valid / ✗ Syntax errors found
   - ✓ All targets have .PHONY / ✗ Missing .PHONY declarations
   - ✓ All targets documented / ✗ Missing help comments
   - Warnings for potential improvements

5. Offer to fix common issues automatically

## Validation Categories

- **Critical**: Syntax errors that prevent execution
- **Important**: Missing .PHONY, undocumented targets
- **Suggestions**: Style improvements, optimization opportunities

## Output Format

```
Makefile Validation Results
===========================

✓ Syntax is valid
✗ 3 targets missing .PHONY declarations: build-dev, run-watch, custom
✗ 2 targets missing help comments: internal-helper, check-deps
⚠ Consider using @ prefix for silent output in: verbose-target

Would you like me to fix these issues? (y/n)
```
